// Package mask menyediakan utilitas untuk menyamarkan sebagian nilai
// data sensitif sebelum dikirim ke caller yang tidak punya izin melihat
// nilai asli.
package mask

// PartialMask menyamarkan value, menyisakan sejumlah karakter terakhir
// tetap terlihat tergantung panjangnya:
//   - panjang >= 10: 4 karakter terakhir terlihat
//   - panjang 6-9:   3 karakter terakhir terlihat
//   - panjang 1-5:   disamarkan penuh (semua '*')
//   - panjang 0:     dikembalikan apa adanya ("")
func PartialMask(value string) string {
	runes := []rune(value)
	n := len(runes)
	if n == 0 {
		return ""
	}

	visible := 0
	switch {
	case n >= 10:
		visible = 4
	case n >= 6:
		visible = 3
	}

	if visible == 0 {
		return repeatStar(n)
	}
	visibleStart := n - visible
	masked := make([]rune, n)
	for i := 0; i < visibleStart; i++ {
		masked[i] = '*'
	}
	copy(masked[visibleStart:], runes[visibleStart:])
	return string(masked)
}

func repeatStar(n int) string {
	stars := make([]rune, n)
	for i := range stars {
		stars[i] = '*'
	}
	return string(stars)
}
