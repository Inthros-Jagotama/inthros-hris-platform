/**
 * primevueLocale — Bilingual PrimeVue locale configuration for Tenant app.
 *
 * Provides locale objects for PrimeVue components (DatePicker, Calendar, etc.)
 * in English and Bahasa Indonesia.
 *
 * @example
 * import { getPrimeLocale } from '@/utils/primevueLocale'
 * const localeObj = getPrimeLocale('id')  // Indonesian locale for PrimeVue
 */

const EN_LOCALE = {
  firstDayOfWeek: 0,
  dayNames: ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
  dayNamesShort: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
  dayNamesMin: ['Su', 'Mo', 'Tu', 'We', 'Th', 'Fr', 'Sa'],
  monthNames: [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ],
  monthNamesShort: [
    'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
    'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
  ],
  today: 'Today',
  clear: 'Clear'
}

const ID_LOCALE = {
  firstDayOfWeek: 1,
  dayNames: ['Minggu', 'Senin', 'Selasa', 'Rabu', 'Kamis', 'Jumat', 'Sabtu'],
  dayNamesShort: ['Min', 'Sen', 'Sel', 'Rab', 'Kam', 'Jum', 'Sab'],    dayNamesMin: ['Mg', 'Sn', 'Sl', 'Rb', 'Km', 'Jm', 'Sb'],
  monthNames: [
    'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
    'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
  ],
  monthNamesShort: [
    'Jan', 'Feb', 'Mar', 'Apr', 'Mei', 'Jun',
    'Jul', 'Agu', 'Sep', 'Okt', 'Nov', 'Des'
  ],
  today: 'Hari Ini',
  clear: 'Hapus'
}

/**
 * Get PrimeVue locale object for the given language code.
 *
 * @param {string} lang — Language code: 'en' or 'id'
 * @returns {object} PrimeVue-compatible locale object
 */
export function getPrimeLocale(lang = 'en') {
  return lang === 'id' ? ID_LOCALE : EN_LOCALE
}
