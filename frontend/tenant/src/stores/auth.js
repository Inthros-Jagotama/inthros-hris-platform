import { reactive } from 'vue'
import api from '@/services/api'

const TOKEN_KEY = 'tenant_token'
const REFRESH_KEY = 'tenant_refresh'
const USER_KEY = 'tenant_user'

// decodeJwtPayload mendekode payload JWT (base64url) tanpa dependency eksternal.
// Mengembalikan objek claims, atau null jika token tidak valid.
function decodeJwtPayload(token) {
  if (!token) return null
  try {
    const part = token.split('.')[1]
    if (!part) return null
    const base64 = part.replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64 + '='.repeat((4 - (base64.length % 4)) % 4)
    const json = decodeURIComponent(
      atob(padded)
        .split('')
        .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
        .join('')
    )
    return JSON.parse(json)
  } catch {
    return null
  }
}

// extractPermissions mengambil daftar permission dari JWT claims.
// JWT berisi field `permissions` (array "resource.action") — di-set backend saat login.
function extractPermissions(token) {
  const claims = decodeJwtPayload(token)
  return Array.isArray(claims?.permissions) ? claims.permissions : []
}

const state = reactive({
  user: JSON.parse(localStorage.getItem(USER_KEY) || 'null'),
  accessToken: localStorage.getItem(TOKEN_KEY) || null,
  refreshToken: localStorage.getItem(REFRESH_KEY) || null,
  permissions: extractPermissions(localStorage.getItem(TOKEN_KEY) || null),
  isAuthenticated: !!localStorage.getItem(TOKEN_KEY)
})

export function useAuth() {
  function applyToken(accessToken) {
    state.accessToken = accessToken
    state.permissions = extractPermissions(accessToken)
    localStorage.setItem(TOKEN_KEY, accessToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
  }

  async function login(email, password) {
    const res = await api.post('/api/v1/platform/login', { email, password })
    const data = res.data?.data || res.data
    // Backend returns snake_case: access_token, refresh_token
    const {
      access_token: accessToken,
      refresh_token: refreshToken,
      user
    } = data
    applyToken(accessToken)
    state.refreshToken = refreshToken
    state.user = user
    state.isAuthenticated = true
    localStorage.setItem(REFRESH_KEY, refreshToken)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    return user
  }

  async function refreshToken() {
    if (!state.refreshToken) throw new Error('No refresh token')
    const res = await api.post('/api/v1/platform/refresh', {
      refresh_token: state.refreshToken
    })
    const data = res.data?.data || res.data
    applyToken(data.access_token)
    return data.access_token
  }

  function logout() {
    state.user = null
    state.accessToken = null
    state.refreshToken = null
    state.permissions = []
    state.isAuthenticated = false
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
    delete api.defaults.headers.common['Authorization']
  }

  // hasPermission — filter Level 2 (Dynamic Menu Rendering):
  // true jika user memiliki permission yang diminta (format "resource.action").
  // Mendukung wildcard: "*" (semua) dan "resource.*" (semua action resource).
  // Saat permissions belum ter-decode (null/empty), default true agar tidak
  // menyembunyikan UI sebelum state siap.
  function hasPermission(required) {
    if (!required) return true
    const perms = state.permissions
    if (!Array.isArray(perms) || perms.length === 0) return true
    if (perms.includes('*')) return true
    if (perms.includes(required)) return true

    // wildcard per resource: resource.*
    const [resource] = String(required).split('.')
    if (resource && perms.includes(resource + '.*')) return true
    return false
  }

  function initAuth() {
    if (state.accessToken) {
      api.defaults.headers.common['Authorization'] = `Bearer ${state.accessToken}`
    }
  }

  initAuth()

  return { state, login, refreshToken, logout, hasPermission }
}
