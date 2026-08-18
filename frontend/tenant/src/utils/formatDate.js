/**
 * formatDate — Bilingual date formatting utility for Tenant app.
 *
 * Provides date formatting functions that respect the current language
 * (English / Bahasa Indonesia) for month names and display conventions.
 *
 * @example
 * import { formatDate, formatDateShort } from '@/utils/formatDate'
 * import { useI18n } from '@/composables/useI18n'
 *
 * const { locale } = useI18n()
 * formatDate('2026-07-30', locale.value)   // → "30 July 2026" / "30 Juli 2026"
 * formatDateShort('2026-07-30', locale.value) // → "30/07/2026"
 */

const MONTHS = {
  en: [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ],
  id: [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ]
}

/**
 * Parse a date string (YYYY-MM-DD) into its components.
 * Returns null if parsing fails.
 */
function parseDate(value) {
  if (!value) return null
  // Handle Date objects
  if (value instanceof Date) {
    return { year: value.getFullYear(), month: value.getMonth() + 1, day: value.getDate() }
  }
  // Handle strings like "2026-07-30"
  const match = String(value).match(/^(\d{4})-(\d{2})-(\d{2})/)
  if (match) {
    return { year: parseInt(match[1], 10), month: parseInt(match[2], 10), day: parseInt(match[3], 10) }
  }
  // Handle ISO strings like "2026-07-30T00:00:00Z"
  const isoMatch = String(value).match(/^(\d{4})-(\d{2})-(\d{2})T/)
  if (isoMatch) {
    return { year: parseInt(isoMatch[1], 10), month: parseInt(isoMatch[2], 10), day: parseInt(isoMatch[3], 10) }
  }
  return null
}

/**
 * Format a date value to long form: "30 July 2026" (EN) or "30 Juli 2026" (ID).
 *
 * @param {string|Date} value — Date string (YYYY-MM-DD), ISO string, or Date object
 * @param {string} lang — Language code: 'en' or 'id' (default: 'en')
 * @returns {string} Formatted date string, or empty string if invalid
 */
export function formatDate(value, lang = 'en') {
  const parsed = parseDate(value)
  if (!parsed) return ''
  const months = MONTHS[lang] || MONTHS.en
  return `${String(parsed.day).padStart(2, '0')} ${months[parsed.month - 1]} ${parsed.year}`
}

/**
 * Format a date value to short form: "30/07/2026".
 *
 * @param {string|Date} value — Date string (YYYY-MM-DD), ISO string, or Date object
 * @returns {string} Formatted date string (dd/mm/yyyy), or empty string if invalid
 */
export function formatDateShort(value) {
  const parsed = parseDate(value)
  if (!parsed) return ''
  return `${String(parsed.day).padStart(2, '0')}/${String(parsed.month).padStart(2, '0')}/${parsed.year}`
}

/**
 * Format month name only: "July" (EN) or "Juli" (ID).
 *
 * @param {number} monthIndex — Month index (1-12)
 * @param {string} lang — Language code: 'en' or 'id'
 * @returns {string} Month name, or empty string if invalid
 */
export function formatMonth(monthIndex, lang = 'en') {
  const months = MONTHS[lang] || MONTHS.en
  if (monthIndex < 1 || monthIndex > 12) return ''
  return months[monthIndex - 1]
}

/**
 * Format a date+time value: "30 July 2026 14:30" (EN) or "30 Juli 2026 14:30" (ID).
 *
 * @param {string|Date} value — ISO timestamp or date string
 * @param {string} lang — Language code: 'en' or 'id'
 * @returns {string} Formatted date+time string, or '-'
 */
export function formatDateDateTime(value, lang = 'en') {
  if (!value) return '-'
  const datePart = formatDate(value, lang)
  if (!datePart) return '-'
  const time = new Date(value).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

/**
 * Format time only: "14:30" from an RFC3339 timestamp.
 *
 * @param {string} value — RFC3339 timestamp
 * @returns {string} Formatted time string, or '-'
 */
export function formatTime(value) {
  if (!value) return '-'
  const d = new Date(value)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

/**
 * Format currency in IDR without decimals: "Rp 1.000.000".
 *
 * @param {number|string|null} value
 * @returns {string} Formatted currency string, or '-'
 */
export function formatCurrency(value) {
  if (value === null || value === undefined || value === '') return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(value))
}
