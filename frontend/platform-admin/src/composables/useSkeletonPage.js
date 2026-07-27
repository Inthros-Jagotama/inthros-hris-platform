import { ref, computed } from 'vue'

/**
 * useSkeletonPage — standardizes the loading/loaded/error pattern used across all platform-admin pages.
 *
 * Provides:
 *  - loading, loaded, error refs
 *  - showInitialSkeleton (loading && !loaded) — true on first load
 *  - showRefreshSkeleton (loading && loaded) — true on auto-refresh / re-fetch
 *  - wrapLoad(fn) — wraps an async function with loading/lifecycle
 *  - reset() — resets all states back to initial
 *
 * @param {Object} options
 * @param {boolean} [options.startLoading=true] — whether loading starts as true (for pages that load on mount)
 * @returns
 */
export function useSkeletonPage(options = {}) {
  const { startLoading = true } = options

  const loading = ref(startLoading)
  const loaded = ref(false)
  const error = ref(null)

  /** True during the very first load (show full skeleton). */
  const showInitialSkeleton = computed(() => loading.value && !loaded.value)

  /** True during subsequent refreshes (show inline skeleton / stale indicator). */
  const showRefreshSkeleton = computed(() => loading.value && loaded.value)

  /** True when data has been loaded at least once. */
  const hasLoaded = computed(() => loaded.value)

  /** True only when NOT loading (data is ready and stable). */
  const isReady = computed(() => !loading.value && loaded.value)

  /**
   * Wraps an async load function with loading/error/lifecycle management.
   *
   * Usage:
   *   async function loadData() {
   *     await wrapLoad(async () => {
   *       const res = await api.get('/endpoint')
   *       data.value = res.data
   *     })
   *   }
   */
  async function wrapLoad(fn) {
    // Prevent concurrent loads
    // Uses `loading && loaded` so the first load (loaded=false) can proceed even though
    // `startLoading: true` has already set loading=true on mount.
    if (loading.value && loaded.value) return
    loading.value = true
    error.value = null

    try {
      await fn()
      loaded.value = true
    } catch (e) {
      error.value = e
      throw e // Re-throw so caller can handle toast / UI feedback
    } finally {
      loading.value = false
    }
  }

  /**
   * Resets all states back to initial (useful for page refresh / retry).
   */
  function reset() {
    loading.value = startLoading
    loaded.value = false
    error.value = null
  }

  return {
    // State
    loading,
    loaded,
    error,

    // Derived
    showInitialSkeleton,
    showRefreshSkeleton,
    hasLoaded,
    isReady,

    // Methods
    wrapLoad,
    reset
  }
}
