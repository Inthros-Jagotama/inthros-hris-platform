package documenttemplate

import "testing"

func TestFormatIndonesianDate(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2026-08-14", "14 Agustus 2026"},
		{"2026-01-01", "1 Januari 2026"},
		{"2026-12-31", "31 Desember 2026"},
		{"1990-02-05", "5 Februari 1990"},
		{"2026-08-14T00:00:00Z", "14 Agustus 2026"},
		{"2026-08-14 07:30:00", "14 Agustus 2026"},
		{"", ""},
		{"   ", ""},
		// Sudah terformat / bukan tanggal ISO → dikembalikan apa adanya.
		{"14 Agustus 2026", "14 Agustus 2026"},
		{"belum diisi", "belum diisi"},
	}
	for _, c := range cases {
		if got := FormatIndonesianDate(c.in); got != c.want {
			t.Errorf("FormatIndonesianDate(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
