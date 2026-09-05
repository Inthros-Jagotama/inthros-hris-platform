package database

import (
	"strings"
	"testing"
)

func TestTenantUsername(t *testing.T) {
	const companyID = "52964b24-26d0-4402-980d-cd1070104e30"

	tests := []struct {
		name     string
		label    string
		expected string
	}{
		{
			name:     "slug pendek dipakai apa adanya",
			label:    "acme",
			expected: "hris_acme_52964b",
		},
		{
			name:     "karakter non-alfanumerik jadi underscore",
			label:    "PT. Wok D Tok",
			expected: "hris_pt__wok_d_tok_52964b",
		},
		{
			// 20 karakter pertama = "perusahaan_dengan_na"; total tepat 32.
			name:     "slug panjang dipotong 20 karakter",
			label:    "perusahaan-dengan-nama-yang-sangat-panjang",
			expected: "hris_perusahaan_dengan_na_52964b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TenantUsername(tt.label, companyID)
			if got != tt.expected {
				t.Errorf("TenantUsername(%q) = %q, want %q", tt.label, got, tt.expected)
			}
		})
	}
}

// Username MySQL maksimal 32 karakter — batas terketat antar dialect yang
// didukung. Label sepanjang apa pun tidak boleh melewatinya.
func TestTenantUsername_TidakMelebihiBatasMySQL(t *testing.T) {
	got := TenantUsername(strings.Repeat("a", 200), "52964b24-26d0-4402-980d-cd1070104e30")
	if len(got) > 32 {
		t.Errorf("username %q panjangnya %d, melebihi batas 32 karakter MySQL", got, len(got))
	}
}

func TestGenerateStrongPassword(t *testing.T) {
	pw, err := GenerateStrongPassword(24)
	if err != nil {
		t.Fatalf("GenerateStrongPassword failed: %v", err)
	}
	if len(pw) != 24 {
		t.Errorf("panjang password = %d, want 24", len(pw))
	}
	// Password dipakai langsung sebagai SQL string literal saat CREATE/ALTER
	// USER, jadi karakter ini tidak boleh pernah muncul.
	if strings.ContainsAny(pw, "'\\") {
		t.Errorf("password %q mengandung kutip tunggal atau backslash", pw)
	}
}
