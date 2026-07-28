import { ref, computed } from 'vue'

/**
 * useSkeletonPage — Standardizes loading/loaded/error states for any page.
 *
 * @param {Object} opts
 * @param {boolean} opts.startLoading — Start with loading=true so skeleton shows immediately (default: true)
 * @returns
 */
export function useSkeletonPage(opts = {}) {
  const { startLoading = true } = opts

  const loading = ref(startLoading)
  const loaded = ref(false)
  const error = ref(null)

  const showInitialSkeleton = computed(() => loading.value && !loaded.value)
  const showRefreshSkeleton = computed(() => loading.value && loaded.value)
  const hasLoaded = computed(() => loaded.value)
  const isReady = computed(() => !loading.value && loaded.value)

  async function wrapLoad(fn) {
    if (loading.value && loaded.value) return
    loading.value = true
    error.value = null
    try {
      await fn()
      loaded.value = true
    } catch (e) {
      error.value = e
      throw e
    } finally {
      loading.value = false
    }
  }

  function reset() {
    loading.value = startLoading
    loaded.value = false
    error.value = null
  }

  return {
    loading, loaded, error,
    showInitialSkeleton, showRefreshSkeleton,
    hasLoaded, isReady,
    wrapLoad, reset
  }
}
