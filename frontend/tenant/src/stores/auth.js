import { reactive } from 'vue'
import api from '@/services/api'

const TOKEN_KEY = 'tenant_token'
const REFRESH_KEY = 'tenant_refresh'
const USER_KEY = 'tenant_user'

const state = reactive({
  user: JSON.parse(localStorage.getItem(USER_KEY) || 'null'),
  accessToken: localStorage.getItem(TOKEN_KEY) || null,
  refreshToken: localStorage.getItem(REFRESH_KEY) || null,
  isAuthenticated: !!localStorage.getItem(TOKEN_KEY)
})

export function useAuth() {
  async function login(email, password) {
    const res = await api.post('/api/v1/platform/login', { email, password })
    const data = res.data?.data || res.data
    // Backend returns snake_case: access_token, refresh_token
    const {
      access_token: accessToken,
      refresh_token: refreshToken,
      user
    } = data
    state.accessToken = accessToken
    state.refreshToken = refreshToken
    state.user = user
    state.isAuthenticated = true
    localStorage.setItem(TOKEN_KEY, accessToken)
    localStorage.setItem(REFRESH_KEY, refreshToken)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
    return user
  }

  async function refreshToken() {
    if (!state.refreshToken) throw new Error('No refresh token')
    const res = await api.post('/api/v1/platform/refresh', {
      refresh_token: state.refreshToken
    })
    const data = res.data?.data || res.data
    const { access_token: accessToken } = data
    state.accessToken = accessToken
    localStorage.setItem(TOKEN_KEY, accessToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
    return accessToken
  }

  function logout() {
    state.user = null
    state.accessToken = null
    state.refreshToken = null
    state.isAuthenticated = false
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
    delete api.defaults.headers.common['Authorization']
  }

  function initAuth() {
    if (state.accessToken) {
      api.defaults.headers.common['Authorization'] = `Bearer ${state.accessToken}`
    }
  }

  initAuth()

  return { state, login, refreshToken, logout }
}
