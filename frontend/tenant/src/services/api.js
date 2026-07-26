import axios from 'axios'
import { useLanguage } from '@/stores/language'

const api = axios.create({
  baseURL: '',
  timeout: 30000,
  headers: { 'Content-Type': 'application/json' }
})

// Request interceptor — auto-attach language header
api.interceptors.request.use((config) => {
  const { state } = useLanguage()
  config.headers['Accept-Language'] = state.lang
  return config
})

// Response interceptor — auto-refresh on 401
api.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config
    if (err.response?.status === 401 && !original._retry) {
      original._retry = true
    }
    return Promise.reject(err)
  }
)

export default api
