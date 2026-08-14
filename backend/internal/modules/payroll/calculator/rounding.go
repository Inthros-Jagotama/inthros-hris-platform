package calculator

import "math"

// RoundingMode menentukan cara pembulatan nominal payroll.
type RoundingMode string

const (
	RoundingNone  RoundingMode = "NONE"
	RoundingRound RoundingMode = "ROUND"
	RoundingCeil  RoundingMode = "CEIL"
	RoundingFloor RoundingMode = "FLOOR"
)

// IsValidRoundingMode memvalidasi mode pembulatan.
func IsValidRoundingMode(mode string) bool {
	switch RoundingMode(mode) {
	case RoundingNone, RoundingRound, RoundingCeil, RoundingFloor:
		return true
	default:
		return false
	}
}

// Round membulatkan value sesuai mode. NONE mengembalikan nilai apa adanya.
func Round(value float64, mode RoundingMode) float64 {
	switch mode {
	case RoundingRound:
		return math.Round(value)
	case RoundingCeil:
		return math.Ceil(value)
	case RoundingFloor:
		return math.Floor(value)
	default:
		return value
	}
}

// RoundToUnit membulatkan value ke kelipatan unit (mis. 1000 = ribuan).
// NONE tetap mengembalikan nilai asli.
func RoundToUnit(value float64, unit float64, mode RoundingMode) float64 {
	if unit <= 0 {
		return Round(value, mode)
	}
	switch mode {
	case RoundingRound:
		return math.Round(value/unit) * unit
	case RoundingCeil:
		return math.Ceil(value/unit) * unit
	case RoundingFloor:
		return math.Floor(value/unit) * unit
	default:
		return value
	}
}

// NormalizeRoundingMode menormalkan string mode ke nilai yang dikenal,
// fallback ke ROUND untuk kompatibilitas dengan default skema.
func NormalizeRoundingMode(mode string) RoundingMode {
	if mode == "" {
		return RoundingRound
	}
	if !IsValidRoundingMode(mode) {
		return RoundingRound
	}
	return RoundingMode(mode)
}
