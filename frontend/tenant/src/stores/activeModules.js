import { reactive } from 'vue'
import api from '@/services/api'

const state = reactive({
  activeModules: [],
  isLoading: false,
  loaded: false
})

export function useActiveModules() {
  async function fetchActiveModules() {
    if (state.loaded) return
    state.isLoading = true
    try {
      const res = await api.get('/api/v1/tenant/company-modules')
      const data = res.data?.data || res.data
      state.activeModules = data.modules || []
      state.loaded = true
    } catch (err) {
      console.error('Failed to fetch active modules:', err)
      state.activeModules = []
    } finally {
      state.isLoading = false
    }
  }

  function hasModule(slug) {
    if (!state.loaded) return true // Default to showing all until loaded
    return state.activeModules.includes(slug)
  }

  function isModuleActive(slug) {
    return state.activeModules.includes(slug)
  }

  function reset() {
    state.activeModules = []
    state.isLoading = false
    state.loaded = false
  }

  return {
    state,
    fetchActiveModules,
    hasModule,
    isModuleActive,
    reset
  }
}
