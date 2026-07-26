import { ref, watch } from 'vue'

/**
 * Composable untuk auto-generate slug dari name field.
 * 
 * @param {() => string} getName  - Getter untuk nilai name (reactive)
 * @param {(v: string) => void} setSlug - Setter untuk nilai slug
 * @returns {{ slugManuallyEdited, slugHighlighted, slugify, resetSlug, disableAutoSlug }}
 * 
 * @example
 * const { slugManuallyEdited, slugHighlighted, resetSlug, disableAutoSlug } = useSlugify(
 *   () => form.value.name,
 *   (v) => { form.value.slug = v }
 * )
 */
export function useSlugify(getName, setSlug) {
  const slugManuallyEdited = ref(false)
  const slugHighlighted = ref(false)
  let highlightTimer

  /**
   * Convert text to URL-friendly slug.
   * @param {string} text
   * @returns {string}
   */
  function slugify(text) {
    return text
      .toLowerCase()
      .trim()
      .replace(/[^a-z0-9\s-]/g, '')
      .replace(/[\s_]+/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-|-$/g, '')
  }

  // Auto-generate slug when name changes (unless manually edited)
  watch(getName, (newName) => {
    if (!slugManuallyEdited.value && newName) {
      setSlug(slugify(newName))
      // Trigger highlight animation (clear any pending timer first)
      clearTimeout(highlightTimer)
      slugHighlighted.value = true
      highlightTimer = setTimeout(() => {
        slugHighlighted.value = false
      }, 600)
    }
  })

  /**
   * Reset slug state for a fresh create form.
   * Clears the manually-edited flag, cancels pending animation, resets highlight.
   */
  function resetSlug() {
    slugManuallyEdited.value = false
    clearTimeout(highlightTimer)
    slugHighlighted.value = false
  }

  /**
   * Disable auto-generation (e.g. when editing existing record).
   */
  function disableAutoSlug() {
    slugManuallyEdited.value = true
  }

  return {
    slugManuallyEdited,
    slugHighlighted,
    slugify,
    resetSlug,
    disableAutoSlug,
  }
}
