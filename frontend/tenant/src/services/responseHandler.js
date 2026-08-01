/**
 * Response Handler — Bilingual response parser
 *
 * Backend mengembalikan response dengan format bilingual:
 *   success response: { "message": "Created successfully" }  ← sudah dalam bahasa yg diminta
 *   error response:   { "error": { "code": "...", "message": "..." } }
 *   validation error: { "error": { "code": "VALIDATION_ERROR", "errors": { "field": ["msg"] } } }
 */

export function getMessage(response, fallback = '') {
  if (!response) return fallback
  if (typeof response.message === 'string') return response.message
  if (response.data?.message) return response.data.message
  return fallback
}

export function getErrorMessage(error, fallback = 'An error occurred') {
  if (!error) return fallback
  const errData = error.response?.data || error
  if (!errData) return fallback
  if (errData.error?.message) {
    return typeof errData.error.message === 'string'
      ? errData.error.message
      : errData.error.message
  }
  if (typeof errData.error === 'string') return errData.error
  if (typeof errData.message === 'string') return errData.message
  return fallback
}

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

export function getStatus(error) {
  return error?.response?.status || 0
}

export function getErrorCode(error) {
  return error?.response?.data?.error?.code || error?.error?.code || ''
}

export function isValidationError(error) {
  return getStatus(error) === 422 || error?.response?.data?.error?.code === 'VALIDATION_ERROR'
}
