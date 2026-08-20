// Package timezone menyediakan fungsi untuk resolusi dan validasi zona waktu IANA
// untuk tenant HRIS, dengan dukungan override per zona dan fallback ke default.
package timezone

import (
	"errors"
	"time"
)

var ErrInvalidTimezone = errors.New("invalid timezone")

var allowed = map[string]bool{
	"Asia/Jakarta":  true,
	"Asia/Makassar": true,
	"Asia/Jayapura": true,
}

var cache = map[string]*time.Location{}

// Allowed mengembalikan daftar timezone IANA yang diizinkan.
func Allowed() []string {
	return []string{"Asia/Jakarta", "Asia/Makassar", "Asia/Jayapura"}
}

// IsValid memeriksa apakah timezone yang diberikan ada dalam daftar whitelist.
func IsValid(tz string) bool {
	return allowed[tz]
}

// loadCached memuat timezone dari cache atau dari time.LoadLocation jika belum ada.
func loadCached(tz string) (*time.Location, error) {
	if loc, ok := cache[tz]; ok {
		return loc, nil
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	cache[tz] = loc
	return loc, nil
}

// Resolve menentukan zona waktu efektif: zoneTimezone (jika tidak nil) menang
// atas companyTimezone; jika companyTimezone kosong, fallback ke Asia/Jakarta.
func Resolve(companyTimezone string, zoneTimezone *string) (*time.Location, error) {
	tz := companyTimezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	if zoneTimezone != nil && *zoneTimezone != "" {
		tz = *zoneTimezone
	}
	if !IsValid(tz) {
		return nil, ErrInvalidTimezone
	}
	return loadCached(tz)
}
