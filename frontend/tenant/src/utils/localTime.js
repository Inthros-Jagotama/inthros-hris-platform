/**
 * toLocalISOString — builds an ISO string with the browser's local timezone
 * offset (e.g. "2026-08-08T14:30:00+07:00"), since several attendance
 * backend endpoints (CreateEvent, CreateOvertimeRequest) require an explicit
 * offset rather than plain UTC or a bare local timestamp.
 *
 * @param {Date} date
 * @returns {string}
 */
export function toLocalISOString(date) {
  const pad = (n) => String(n).padStart(2, '0')
  const offsetMin = -date.getTimezoneOffset()
  const sign = offsetMin >= 0 ? '+' : '-'
  const offH = pad(Math.floor(Math.abs(offsetMin) / 60))
  const offM = pad(Math.abs(offsetMin) % 60)
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}${sign}${offH}:${offM}`
}

/**
 * localDateTimeISOString — combines a date-only string (YYYY-MM-DD) with a
 * "HH:mm:ss" time string into an ISO string carrying the browser's current
 * local timezone offset. Used for backend fields like
 * overtime start_time_local/end_time_local, which are parsed as RFC3339.
 *
 * @param {string} dateStr — "YYYY-MM-DD"
 * @param {string} timeStr — "HH:mm:ss" or "HH:mm"
 * @returns {string}
 */
export function localDateTimeISOString(dateStr, timeStr) {
  const [h = '00', m = '00', s = '00'] = String(timeStr).split(':')
  const [y, mo, d] = String(dateStr).split('-').map(Number)
  const date = new Date(y, mo - 1, d, Number(h), Number(m), Number(s))
  return toLocalISOString(date)
}
