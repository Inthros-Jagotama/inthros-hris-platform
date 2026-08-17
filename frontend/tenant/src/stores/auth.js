import { reactive } from 'vue'
import api from '@/services/api'
import { resetMyEmployeeCache } from '@/composables/useMyEmployee'

const TOKEN_KEY = 'tenant_token'
const REFRESH_KEY = 'tenant_refresh'
const USER_KEY = 'tenant_user'
const COMPANY_KEY = 'tenant_company'

// safeJsonParse mem-parse nilai localStorage yang seharusnya JSON, tapi
// bisa saja rusak/bukan JSON (data lama dari versi app sebelumnya, dsb).
// Mengembalikan fallback alih-alih melempar error yang bisa meng-crash app.
function safeJsonParse(value, fallback = null) {
  if (value == null) return fallback
  try {
    return JSON.parse(value)
  } catch {
    return fallback
  }
}

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
  user: safeJsonParse(localStorage.getItem(USER_KEY)),
  accessToken: localStorage.getItem(TOKEN_KEY) || null,
  refreshToken: localStorage.getItem(REFRESH_KEY) || null,
  permissions: extractPermissions(localStorage.getItem(TOKEN_KEY) || null),
  // company: { id, slug, name } — di-sync otomatis dari response header X-Tenant-ID
  // (backend TenantResolver) sehingga state company selalu sinkron tanpa isi manual.
  company: safeJsonParse(localStorage.getItem(COMPANY_KEY)),
  isAuthenticated: !!localStorage.getItem(TOKEN_KEY)
})

export function useAuth() {
  function applyToken(accessToken) {
    state.accessToken = accessToken
    state.permissions = extractPermissions(accessToken)
    localStorage.setItem(TOKEN_KEY, accessToken)
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
  }

  // setCompany memperbarui state company + persist ke localStorage.
  // Dipanggil dari: login (response X-Tenant-ID), api.js response interceptor
  // (agar selalu sinkron), dan resolve subdomain (prefill sebelum login).
  function setCompany(company) {
    if (!company || !company.id) return
    state.company = {
      ...(state.company || {}),
      ...company
    }
    localStorage.setItem(COMPANY_KEY, JSON.stringify(state.company))
    // Backend TenantResolver membaca header X-Tenant-ID bila ada — kirim agar
    // konsisten meski Host tidak ter-resolve (mis. dev tanpa subdomain).
    if (state.company.id) {
      api.defaults.headers.common['X-Tenant-ID'] = state.company.id
    }
  }

  async function login(email, password, companySlug, companyId, companyName) {
    // Tenant login: user tenant (employee) login. Company diidentifikasi via
    // company_slug/company_id (opsional) — backend juga bisa auto-resolve dari
    // Host header (mode SaaS, tiap company punya URL sendiri), sehingga FE tidak
    // wajib mengirim company manual.
    const payload = {
      email,
      password
    }
    if (companyId) {
      payload.company_id = companyId
    } else if (companySlug) {
      payload.company_slug = companySlug
    }
    const res = await api.post('/api/v1/tenant/auth/login', payload)
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

    // Sync company dari response header X-Tenant-ID (di-set backend TenantResolver).
    // Fallback: user?.company_id (selalu ada di TenantLoginResponse) — dipakai saat
    // response header tidak ada (mis. dev tanpa subdomain, company dikirim manual/env).
    const resolvedId = res.headers?.['x-tenant-id'] || user?.company_id || companyId || ''
    if (resolvedId) {
      // Jika id berubah (login ke company lain), slug/name lama di-reset agar
      // tidak menampilkan company yang salah (id adalah sumber kebenaran).
      const companyChanged = state.company?.id && state.company.id !== resolvedId
      setCompany({
        id: resolvedId,
        slug: companySlug || (companyChanged ? '' : state.company?.slug) || '',
        name: companyName || user?.company_name || (companyChanged ? '' : state.company?.name) || ''
      })
    }
    // Reset cache employee_id level-modul: menutup skenario ganti user tanpa
    // logout (login langsung sebagai user lain) agar halaman berbasis
    // "my employee" tidak memakai employee_id user sebelumnya.
    resetMyEmployeeCache()
    return user
  }

  async function refreshToken() {
    if (!state.refreshToken) throw new Error('No refresh token')
    const res = await api.post('/api/v1/tenant/auth/refresh', {
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
    state.company = null
    state.isAuthenticated = false
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(REFRESH_KEY)
    localStorage.removeItem(USER_KEY)
    localStorage.removeItem(COMPANY_KEY)
    delete api.defaults.headers.common['Authorization']
    delete api.defaults.headers.common['X-Tenant-ID']
    // Reset cache employee_id level-modul agar user berikutnya yang login
    // tidak mewarisi data user sebelumnya di halaman berbasis "my employee"
    // (mis. Overtime Requests, Attendance, Leave).
    resetMyEmployeeCache()
  }

  // hasPermission — filter Level 2 (Dynamic Menu Rendering):
  // true jika user memiliki permission yang diminta (format "resource.action"
  // atau "resource.submenu.action"). Mendukung wildcard: "*" (semua) dan
  // "resource.*" (semua action resource). Untuk permission level-submenu
  // ("resource.submenu.action"), permission level-module dengan action yang
  // sama ("resource.action") juga dianggap memenuhi — supaya role yang hanya
  // punya "competency.view" tetap melihat semua submenu competency.
  // Saat permissions belum ter-decode (null/empty), default true agar tidak
  // menyembunyikan UI sebelum state siap.
  function hasPermission(required) {
    if (!required) return true
    const perms = state.permissions
    if (!Array.isArray(perms) || perms.length === 0) return true
    if (perms.includes('*')) return true
    if (perms.includes(required)) return true

    const parts = String(required).split('.')
    const [resource] = parts
    // wildcard per resource: resource.*
    if (resource && perms.includes(resource + '.*')) return true

    // level-submenu: resource.submenu.action → module-level resource.action
    // juga memenuhi (mis. competency.view → competency.settings.view).
    if (parts.length >= 3 && resource) {
      const action = parts[parts.length - 1]
      if (perms.includes(resource + '.' + action)) return true
    }

    // Arah sebaliknya: role yang hanya diberi permission level-submenu
    // (mis. competency.settings.view tanpa competency.view) tetap harus bisa
    // melihat menu module-nya — resource.action dianggap terpenuhi bila ada
    // permission resource.<submenu>.action mana pun.
    if (parts.length === 2 && resource) {
      const action = parts[1]
      const prefix = resource + '.'
      if (perms.some(p => p.startsWith(prefix) && p.endsWith('.' + action))) return true
    }
    return false
  }

  // hasExactPermission — sama seperti hasPermission tapi TANPA fallback
  // module-covers-submenu. Dipakai untuk submenu admin-config yang sengaja
  // harus terpisah dari permission module-levelnya (mis. "approval.settings"
  // tidak boleh otomatis terpenuhi oleh "approval.view", karena hampir semua
  // role — termasuk Employee default — punya "approval.view" untuk melihat
  // tugas approval sehari-hari, bukan untuk mengelola alur persetujuan).
  function hasExactPermission(required) {
    if (!required) return true
    const perms = state.permissions
    if (!Array.isArray(perms) || perms.length === 0) return true
    if (perms.includes('*')) return true
    return perms.includes(required)
  }

  function initAuth() {
    if (state.accessToken) {
      api.defaults.headers.common['Authorization'] = `Bearer ${state.accessToken}`
    }
  }

  initAuth()

  return { state, login, refreshToken, logout, hasPermission, hasExactPermission, setCompany }
}
