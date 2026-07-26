package httputil

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Custom validation tag constants — use these in binding tags:
//
//	type CreateEmployeeRequest struct {
//	    NIK  string `json:"nik"  binding:"omitempty,nik"`
//	    NPWP string `json:"npwp" binding:"omitempty,npwp"`
//	}
const (
	TagNIK              = "nik"
	TagNPWP             = "npwp"
	TagPhoneID          = "phone_id"
	TagPostalCode       = "postal_code"
	TagDateID           = "date_id"
	TagKodePos          = "kode_pos"       // alias for postal_code
	TagPhoneIndonesia   = "phone_indonesia" // alias for phone_id

	// Additional Indonesian data formats
	TagKK               = "kk"          // Kartu Keluarga (16 digit)
	TagNoRekening       = "no_rekening" // Nomor rekening bank
	TagPassport         = "passport"    // Nomor passport Indonesia
	TagSIM              = "sim"         // Nomor SIM
	TagNPWPFormatted    = "npwp_format" // NPWP dengan format XX.XXX.XXX.X-XXX.XXX
)

// Compiled regex patterns (lazily initialised).
var (
	reNIK        *regexp.Regexp
	reNPWP       *regexp.Regexp
	rePhoneID    *regexp.Regexp
	rePostalCode *regexp.Regexp
	reDateID     *regexp.Regexp
	reKK         *regexp.Regexp
	reNoRekening *regexp.Regexp
	rePassport   *regexp.Regexp
	reSIM        *regexp.Regexp
	reNPWPFormat *regexp.Regexp
)

func init() {
	reNIK = regexp.MustCompile(`^[0-9]{16}$`)
	reNPWP = regexp.MustCompile(`^[0-9]{15,16}$`)
	rePhoneID = regexp.MustCompile(`^(?:\+62|0)8[0-9]{7,11}$`)
	rePostalCode = regexp.MustCompile(`^[0-9]{5}$`)
	reDateID = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	reKK = regexp.MustCompile(`^[0-9]{16}$`)        // 16 digit
	reNoRekening = regexp.MustCompile(`^[0-9]{8,20}$`) // 8-20 digit (variasi antar bank)
	rePassport = regexp.MustCompile(`^[A-Za-z][0-9]{8}$`) // e.g. A12345678
	reSIM = regexp.MustCompile(`^[0-9]{12}$`)        // 12 digit
	reNPWPFormat = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}\.\d{1}-\d{3}\.\d{3}$`) // XX.XXX.XXX.X-XXX.XXX

	// Auto-register custom validators with Gin's validator engine.
	RegisterCustomValidators()
}

// RegisterCustomValidators registers all Indonesian custom validators with
// Gin's underlying go-playground/validator engine.
//
// This function is called automatically via init() when the httputil package
// is imported. You can also call it explicitly in main() if needed.
func RegisterCustomValidators() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	validators := map[string]validator.Func{
		TagNIK:            validateNIK,
		TagNPWP:           validateNPWP,
		TagPhoneID:        validatePhoneID,
		TagPhoneIndonesia: validatePhoneID,
		TagPostalCode:     validatePostalCode,
		TagKodePos:        validatePostalCode,
		TagDateID:         validateDateID,
		TagKK:             validateKK,
		TagNoRekening:     validateNoRekening,
		TagPassport:       validatePassport,
		TagSIM:            validateSIM,
		TagNPWPFormatted:  validateNPWPFormatted,
	}

	for tag, fn := range validators {
		if err := v.RegisterValidation(tag, fn); err != nil {
			// Log only — validation will fall through to the default handler.
			// In production you may want a real logger here.
			println("httputil: failed to register validator '" + tag + "': " + err.Error())
		}
	}
}

// --- Validation functions ---------------------------------------------------

// validateNIK validates a 16-digit NIK (Nomor Induk Kependudukan).
func validateNIK(fl validator.FieldLevel) bool {
	return reNIK.MatchString(fl.Field().String())
}

// validateNPWP validates a 15/16-digit NPWP (Nomor Pokok Wajib Pajak).
// Accepts plain digits only (formatted versions with dots/dashes are not
// accepted at binding level — strip formatting before sending).
func validateNPWP(fl validator.FieldLevel) bool {
	return reNPWP.MatchString(fl.Field().String())
}

// validatePhoneID validates an Indonesian phone number.
// Accepted formats:
//   - 08xxxxxxxxxx (10-13 digits after 08)
//   - +628xxxxxxxxxx (10-13 digits after +62)
func validatePhoneID(fl validator.FieldLevel) bool {
	return rePhoneID.MatchString(fl.Field().String())
}

// validatePostalCode validates a 5-digit Indonesian postal code.
func validatePostalCode(fl validator.FieldLevel) bool {
	return rePostalCode.MatchString(fl.Field().String())
}

// validateDateID validates a date in YYYY-MM-DD format.
func validateDateID(fl validator.FieldLevel) bool {
	return reDateID.MatchString(fl.Field().String())
}

// validateKK validates a 16-digit KK (Kartu Keluarga).
func validateKK(fl validator.FieldLevel) bool {
	return reKK.MatchString(fl.Field().String())
}

// validateNoRekening validates an Indonesian bank account number.
// Length varies by bank (typically 8-20 digits).
func validateNoRekening(fl validator.FieldLevel) bool {
	return reNoRekening.MatchString(fl.Field().String())
}

// validatePassport validates an Indonesian passport number.
// Format: 1 letter followed by 8 digits (e.g. A12345678).
func validatePassport(fl validator.FieldLevel) bool {
	return rePassport.MatchString(fl.Field().String())
}

// validateSIM validates a 12-digit SIM (Surat Izin Mengemudi) number.
func validateSIM(fl validator.FieldLevel) bool {
	return reSIM.MatchString(fl.Field().String())
}

// validateNPWPFormatted validates NPWP with formatting XX.XXX.XXX.X-XXX.XXX.
func validateNPWPFormatted(fl validator.FieldLevel) bool {
	return reNPWPFormat.MatchString(fl.Field().String())
}

// Messages for custom validators are now defined in locale.go.
// See the localeMessages map in that file.
