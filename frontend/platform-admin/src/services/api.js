import axios from 'axios'
import { useAuth } from '@/stores/auth'

const api = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})

// Response interceptor — auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config
    if (err.response?.status === 401 && !original._retry) {
      original._retry = true
      try {
        const { refreshToken, logout } = useAuth()
        await refreshToken()
        original.headers['Authorization'] = `Bearer ${localStorage.getItem('platform_admin_token')}`
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
