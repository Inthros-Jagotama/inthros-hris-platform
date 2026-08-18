import { useI18n } from '@/composables/useI18n'

/**
 * genderLabel — bilingual gender label resolver.
 * Usage: genderLabel('M') → "Laki-laki" (id) / "Male" (en)
 *
 * @param {string|null} value - 'M' or 'F'
 * @returns {string}
 */
export function genderLabel(value) {
  if (!value) return ''
  const { t } = useI18n()
  return t('employee.gender_' + value.toLowerCase())
}
