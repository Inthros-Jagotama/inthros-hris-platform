package documenttemplate

import (
	"strconv"
	"strings"
	"time"
)

// indonesianMonths adalah nama bulan lengkap Bahasa Indonesia, dipakai oleh
// FormatIndonesianDate untuk memformat variable tanggal di dokumen.
var indonesianMonths = [12]string{
	"Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// FormatIndonesianDate memformat tanggal ISO "2006-01-02" menjadi format
// Bahasa Indonesia "tanggal nama-bulan-lengkap tahun" (mis. "2026-08-14" →
// "14 Agustus 2026"). Nilai kosong atau yang tidak bisa di-parse dikembalikan
// apa adanya (idempotent — tanggal yang sudah terformat tidak berubah), dan
// format datetime (mis. "2006-01-02T00:00:00Z") tetap diterima dengan hanya
// mengambil bagian tanggalnya.
func FormatIndonesianDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	datePart := s
	if len(s) > 10 && s[4] == '-' && s[7] == '-' {
		datePart = s[:10]
	}
	t, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		return s
	}
	return strconv.Itoa(t.Day()) + " " + indonesianMonths[t.Month()-1] + " " + strconv.Itoa(t.Year())
}
