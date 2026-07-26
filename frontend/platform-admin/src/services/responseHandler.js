/**
 * Response Handler — Bilingual response parser
 *
 * Backend mengembalikan response dengan format bilingual:
 *   success response: { "message": "Created successfully" }  ← sudah dalam bahasa yg diminta
 *   error response:   { "error": { "code": "...", "message": "..." } }
 *   validation error: { "error": { "code": "VALIDATION_ERROR", "errors": { "field": ["msg"] } } }
 */

/**
 * Extract message dari success response.
 * Backend sudah mengembalikan message dalam bahasa yg sesuai Accept-Language header.
 */
export function getMessage(response, fallback = '') {
  if (!response) return fallback
  if (typeof response.message === 'string') return response.message
  if (response.data?.message) return response.data.message
  return fallback
}

/**
 * Extract error message dari error response.
 * Handle format: { error: { message: "..." } } atau { error: "string" }
 */
export function getErrorMessage(error, fallback = 'An error occurred') {
  if (!error) return fallback
  const errData = error.response?.data || error
  if (!errData) return fallback

  // Format: { error: { message: "..." } }
  if (errData.error?.message) {
    return typeof errData.error.message === 'string'
      ? errData.error.message
      : errData.error.message
  }
  // Format: { error: "string" }
  if (typeof errData.error === 'string') return errData.error
  // Format: { message: "..." }
  if (typeof errData.message === 'string') return errData.message

  return fallback
}

/**
 * Extract validation errors dari error response.
 * Backend bisa mengirim dengan key "errors" atau "fields":
 *   Format 1: { error: { code: "VALIDATION_ERROR", errors: { field: ["msg"] } } }
 *   Format 2: { error: { code: "VALIDATION_ERROR", fields: { field: ["msg"] } } }
 *
 * Selalu mengembalikan object { field: "string" } — array values di-implode dengan koma.
 * Pesan error sudah dalam bahasa yang sesuai dengan Accept-Language header.
 */
export function getValidationErrors(error) {
  if (!error) return {}
  const errData = error.response?.data || error
  const raw = errData?.error?.errors || errData?.error?.fields || {}

  // Implode array values menjadi comma-separated string
  const result = {}
  for (const [key, value] of Object.entries(raw)) {
    result[key] = Array.isArray(value) ? value.join(', ') : String(value)
  }
  return result
}

/**
 * Get HTTP status code from error
 */
export function getStatus(error) {
  return error?.response?.status || 0
}

/**
 * Check if error is validation error
 */
export function isValidationError(error) {
  return getStatus(error) === 422 || error?.response?.data?.error?.code === 'VALIDATION_ERROR'
}
