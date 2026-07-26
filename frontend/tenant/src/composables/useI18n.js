import { computed } from 'vue'
import { useLanguage } from '@/stores/language'
import en from '@/locales/en.json'
import id from '@/locales/id.json'

const messages = { en, id }

/**
 * useI18n — Composable for translating static UI text in Tenant app.
 *
 * Returns a `t(key, params?)` function that looks up the given dot-notation key
 * in the current locale (EN/ID) and falls back to English if not found.
 *
 * @example
 * const { t } = useI18n()
 * t('nav.dashboard')          // → "Dashboard" or "Dasbor"
 * t('time.minutes_ago', '5')  // → "5 minutes ago" or "5 menit yang lalu"
 */
export function useI18n() {
  const { state } = useLanguage()

  function deepLookup(obj, path) {
    return path.split('.').reduce((current, key) => {
      return current?.[key] !== undefined ? current[key] : undefined
    }, obj)
  }

  function t(key, ...params) {
    let text = deepLookup(messages[state.lang], key)
    if (text === undefined) {
      text = deepLookup(messages.en, key)
    }
    if (text === undefined) {
      return key
    }

    if (params && params.length > 0) {
      let paramIndex = 0
      return text.replace(/\{(\w+)\}|%s/g, (match, named) => {
        if (named && params[0]?.[named] !== undefined) {
          return params[0][named]
        }
        const val = params[paramIndex]
        paramIndex++
        return val !== undefined ? String(val) : match
      })
    }

    return text
  }

  const isID = computed(() => state.lang === 'id')
  const isEN = computed(() => state.lang === 'en')
  const locale = computed(() => state.lang)

  return { t, isID, isEN, locale }
}
