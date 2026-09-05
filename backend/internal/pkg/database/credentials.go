package database

import (
	"crypto/rand"
	"fmt"
	"strings"
)

// GenerateStrongPassword membuat password acak kuat untuk user database tenant.
//
// Charset sengaja tidak memuat kutip tunggal maupun backslash: password dipakai
// langsung sebagai SQL string literal saat CREATE/ALTER USER, dan cara escape
// karakter tersebut berbeda antara MySQL dan PostgreSQL — lebih aman tidak
// pernah menghasilkannya daripada meng-escape per dialect.
func GenerateStrongPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*-_=+"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

// TenantUsername menyusun nama user database milik satu tenant, berbentuk
// "hris_<label dipotong>_<6 hex pertama company id>".
//
// Suffix company id dipakai supaya dua company dengan label mirip (atau label
// panjang yang terpotong di titik yang sama) tidak pernah memakai user yang
// sama. Panjang total dijaga <= 32 karakter mengikuti batas terketat antar
// dialect (username MySQL maksimal 32; identifier PostgreSQL 63), dan karakter
// selain huruf/angka diubah jadi underscore agar aman dipakai sebagai
// identifier tanpa escaping tambahan.
func TenantUsername(label, companyID string) string {
	sanitized := make([]rune, 0, len(label))
	for _, r := range strings.ToLower(label) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sanitized = append(sanitized, r)
		} else {
			sanitized = append(sanitized, '_')
		}
	}
	// 32 total = "hris_" (5) + label (maks 20) + "_" (1) + 6 hex company id.
	if len(sanitized) > 20 {
		sanitized = sanitized[:20]
	}

	idHex := strings.ReplaceAll(companyID, "-", "")
	if len(idHex) > 6 {
		idHex = idHex[:6]
	}

	return fmt.Sprintf("hris_%s_%s", string(sanitized), idHex)
}
