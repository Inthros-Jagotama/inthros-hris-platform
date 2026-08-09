import { ref } from 'vue'
import api from '@/services/api'

let cachedEmployeeId = null
let inFlight = null

/**
 * resetMyEmployeeCache — menghapus cache employee_id level-modul.
 * Dipanggil saat logout (auth store) agar user berikutnya yang login tidak
 * mewarisi employee_id user sebelumnya (SPA tidak me-reload halaman).
 */
export function resetMyEmployeeCache() {
  cachedEmployeeId = null
  inFlight = null
}

/**
 * useMyEmployee — resolves the logged-in user's own employee_id via
 * GET /api/v1/tenant/user-accounts/me. Cached module-wide for the session
 * since the mapping never changes without a fresh login.
 */
export function useMyEmployee() {
  const employeeId = ref(cachedEmployeeId)

  async function loadMyEmployeeId() {
    if (cachedEmployeeId) {
      employeeId.value = cachedEmployeeId
      return cachedEmployeeId
    }
    if (!inFlight) {
      inFlight = api.get('/api/v1/tenant/user-accounts/me')
        .then(res => {
          cachedEmployeeId = res.data?.data?.employee_id || ''
          return cachedEmployeeId
        })
        .finally(() => { inFlight = null })
    }
    employeeId.value = await inFlight
    return employeeId.value
  }

  return { employeeId, loadMyEmployeeId }
}
