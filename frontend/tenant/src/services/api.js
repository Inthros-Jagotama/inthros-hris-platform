import axios from 'axios'
import { useAuth } from '@/stores/auth'
import { useLanguage } from '@/stores/language'

const api = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})

// Request interceptor — auto-attach language header + auth token + X-Tenant-ID
api.interceptors.request.use((config) => {
  const { state } = useLanguage()
  config.headers['Accept-Language'] = state.lang

  const { state: authState } = useAuth()
  if (authState.accessToken) {
    config.headers['Authorization'] = `Bearer ${authState.accessToken}`
  }
  // Backend TenantResolver membaca header X-Tenant-ID bila ada — kirim agar
  // company selalu sinkron meski Host tidak ter-resolve (mis. dev tanpa subdomain).
  if (authState.company?.id) {
    config.headers['X-Tenant-ID'] = authState.company.id
  }
  return config
})

// Response interceptor — sync company dari X-Tenant-ID header + auto-refresh on 401
api.interceptors.response.use(
  (res) => {
    // Backend TenantResolver meng-set response header X-Tenant-ID (company_id)
    // pada semua response /api/v1/tenant/** — sync ke store agar selalu aktual.
    const tenantId = res.headers?.['x-tenant-id']
    if (tenantId) {
      const { setCompany } = useAuth()
      setCompany({ id: tenantId })
    }
    return res
  },
  async (err) => {
    const original = err.config
    const isRefreshCall = original?.url?.includes('/auth/refresh')

    if (err.response?.status === 401) {
      const { refreshToken, logout } = useAuth()

      // Endpoint refresh sendiri ditolak (refresh token expired/invalid), atau
      // request ini sudah pernah retry sekali — jangan coba refresh lagi
      // (hindari infinite loop), langsung logout.
      if (isRefreshCall || original._retry) {
        logout()
        window.location.href = '/login'
        return Promise.reject(err)
      }

      original._retry = true
      try {
        await refreshToken()
        original.headers['Authorization'] = `Bearer ${localStorage.getItem('tenant_token')}`
        return api(original)
      } catch {
        logout()
        window.location.href = '/login'
      }
    }
    return Promise.reject(err)
  }
)

export default api
