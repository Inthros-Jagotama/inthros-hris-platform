package documenttemplate

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

// placeholderRe mencocokkan {{key}} — key boleh berisi huruf/angka/underscore/dot.
var placeholderRe = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.]+)\s*\}\}`)

// wTRe mencocokkan elemen teks <w:t>...</w:t> (termasuk atribut seperti
// xml:space="preserve"). Dipakai untuk menggabungkan teks per paragraf agar
// placeholder yang terpecah antar run tetap terdeteksi & ter-resolve.
var wTRe = regexp.MustCompile(`(?s)<w:t\b[^>]*>(.*?)</w:t>`)

// wPRe mencocokkan satu paragraf <w:p>...</w:p>. Paragraf tidak bersarang di
// WordprocessingML (paragraf di dalam sel tabel tetap dipisah oleh w:tc), jadi
// pencocokan non-greedy dari tag buka ke tag tutup terdekat sudah benar.
var wPRe = regexp.MustCompile(`(?s)<w:p\b[^>]*>.*?</w:p>`)

// readZipEntry membaca isi sebuah entry zip (data mentah, belum di-decompress).
// Bila verifikasi CRC32 gagal (zip: checksum error — beberapa penulis .docx
// seperti WPS/Word versi tertentu menulis CRC yang tidak akurat namun file
// tetap valid dan bisa dibuka Word), fallback membaca stream mentah dan
// mendekompresi manual tanpa cek CRC. Word sendiri toleran terhadap CRC salah,
// jadi kita ikut toleran agar template seperti itu tetap bisa diproses.
func readZipEntry(f *zip.File) ([]byte, error) {
	rc, err := f.Open()
	if err == nil {
		data, rerr := io.ReadAll(rc)
		rc.Close()
		if rerr == nil {
			return data, nil
		}
		// CRC mismatch — jatuh ke jalur tanpa verifikasi CRC di bawah.
	}
	raw, err := f.OpenRaw()
	if err != nil {
		return nil, fmt.Errorf("failed to open docx part %s: %w", f.Name, err)
	}
	// OpenRaw mengembalikan io.Reader (bukan ReadCloser) — tidak perlu Close.
	switch f.Method {
	case zip.Store:
		return io.ReadAll(raw)
	case zip.Deflate:
		zr := flate.NewReader(raw)
		defer zr.Close()
		return io.ReadAll(zr)
	default:
		return nil, fmt.Errorf("unsupported compression method %d for docx part %s", f.Method, f.Name)
	}
}

// extractDocxPlaceholders membuka file .docx (format zip OOXML), membaca semua
// bagian XML di folder word/ (document.xml, header/footer, dll.), lalu
// mengembalikan daftar unik placeholder {{key}} yang ditemukan (diurutkan).
// Placeholder detection (Phase 3) — memakai stdlib archive/zip, tanpa library
// eksternal; pada fase ini cukup deteksi teks, belum parsing layout Word.
func extractDocxPlaceholders(data []byte) ([]string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid docx (not a valid zip/OOXML): %w", err)
	}
	seen := map[string]bool{}
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, "word/") || !strings.HasSuffix(f.Name, ".xml") {
			continue
		}
		buf, err := readZipEntry(f)
		if err != nil {
			continue
		}
		// Deteksi dilakukan pada teks gabungan per paragraf (bukan mentah XML)
		// agar placeholder yang dipecah Word menjadi beberapa <w:t>/run tetap
		// terdeteksi — konsisten dengan resolusi di resolveDocxVariables.
		for _, pm := range wPRe.FindAllString(string(buf), -1) {
			for _, m := range placeholderRe.FindAllStringSubmatch(paragraphJoinedText(pm), -1) {
				seen[m[1]] = true
			}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out, nil
}

// registryKeys membangun set key terdaftar dari VariableRegistry (satu sumber
// kebenaran placeholder; spec §7 — registry di backend, bukan tabel DB).
func registryKeys() map[string]bool {
	keys := map[string]bool{}
	for _, g := range VariableRegistry() {
		for _, v := range g.Variables {
			keys[v.Key] = true
		}
	}
	return keys
}

// unknownPlaceholders mengembalikan placeholder yang dipakai di template tetapi
// tidak terdaftar di registry (variable validation, Phase 3).
func unknownPlaceholders(placeholders []string) []string {
	keys := registryKeys()
	var out []string
	for _, p := range placeholders {
		if !keys[p] {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

// sampleData mengembalikan data contoh untuk preview (spec §8.3) — preview
// tidak memakai data asli, hanya contoh agar hasil variabel terlihat.
func sampleData() map[string]string {
	return map[string]string{
		"employee.employee_id":       "EMP-2026-001",
		"employee.name":              "Asep Ruswanda",
		"employee.nik":               "199001012015011001",
		"employee.family_id":         "3201010101010001",
		"employee.mother_name":       "Siti Aminah",
		"employee.gender":            "Laki-laki",
		"employee.dob":               "1990-01-01",
		"employee.pob":               "Bandung",
		"employee.nationality_type":  "WNI",
		"employee.nationality_id":    "ID",
		"employee.passport":          "A1234567",
		"employee.phone_number":      "081234567890",
		"employee.email":             "asep@example.com",
		"employee.linkedin":          "linkedin.com/in/asep",
		"employee.instagram":         "@asep",
		"employee.religion":          "Islam",
		"employee.marital_status":    "Menikah",
		"employee.status":            "active",
		"employee.join_date":         "2018-06-01",
		"employee.position":          "HR Staff",
		"employee.organization":      "HR Division",
		"contract.number":            "CTR-2026-001",
		"contract.start_date":        "2026-01-01",
		"contract.end_date":          "2027-01-01",
		"movement.number":            "SK-MOV-2026-001",
		"movement.effective_date":    "2026-01-15",
		"movement.previous_position": "HR Staff",
		"movement.new_position":      "Senior HR Staff",
		"company.name":               "PT Maju Bersama",
		"company.address":            "Jl. Merdeka No. 1, Jakarta",
	}
}

// xmlEscape meng-escape karakter XML di dalam nilai variabel agar aman saat
// disisipkan ke document.xml (mis. & < > " ').
func xmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}

// paragraphJoinedText menggabungkan teks semua <w:t> dalam satu paragraf sesuai
// urutan kemunculannya. Word sering memecah teks (termasuk placeholder {{key}})
// menjadi beberapa run — mis. karena spell-check (w:proofErr), paste dari
// sumber lain, atau format berbeda per bagian — sehingga teks yang tampil utuh
// di Word tersimpan terpecah di beberapa <w:t>. Penggabungan ini membuat
// placeholder kembali utuh untuk deteksi & resolusi.
func paragraphJoinedText(p string) string {
	var sb strings.Builder
	for _, m := range wTRe.FindAllStringSubmatch(p, -1) {
		sb.WriteString(m[1])
	}
	return sb.String()
}

// resolveDocxXML mengganti {{key}} pada seluruh paragraf dokumen XML (document,
// header, footer, dll.) menggunakan values. Placeholder yang tidak ada di values
// dibiarkan apa adanya (tetap tampil sebagai {{key}}).
func resolveDocxXML(xml string, values map[string]string) string {
	if len(values) == 0 {
		return xml
	}
	return wPRe.ReplaceAllStringFunc(xml, func(p string) string {
		return resolveParagraph(p, values)
	})
}

// resolveParagraph mengganti {{key}} di dalam satu paragraf, termasuk
// placeholder yang terpecah antar beberapa <w:t> (kasus umum dari Word). Nilai
// dimasukkan ke run pertama yang memuat awal placeholder (format run tersebut
// yang dipakai); run di tengah dikosongkan; run terakhir hanya menyisakan teks
// setelah placeholder. Struktur XML di luar <w:t> tidak diubah.
func resolveParagraph(p string, values map[string]string) string {
	ms := wTRe.FindAllStringSubmatchIndex(p, -1)
	if len(ms) == 0 {
		return p
	}
	runs := make([]struct{ full, text string }, len(ms))
	for i, m := range ms {
		runs[i].full = p[m[0]:m[1]]
		runs[i].text = p[m[2]:m[3]]
	}

	// Teks gabungan + offset kumulatif per run (di teks gabungan).
	var joined strings.Builder
	offs := make([]int, len(runs)+1)
	for i, r := range runs {
		offs[i] = joined.Len()
		joined.WriteString(r.text)
	}
	offs[len(runs)] = joined.Len()
	js := joined.String()

	// Kumpulkan span {{key}} yang punya nilai di values.
	type repl struct{ start, end int; value string }
	var repls []repl
	for _, m := range placeholderRe.FindAllStringSubmatchIndex(js, -1) {
		val, ok := values[js[m[2]:m[3]]]
		if !ok {
			continue
		}
		repls = append(repls, repl{start: m[0], end: m[1], value: xmlEscape(val)})
	}
	if len(repls) == 0 {
		return p
	}

	// Rekonstruksi teks per run: teks asli di luar span disalin ke run asalnya,
	// nilai placeholder ditulis ke run yang memuat awal span, dan run di tengah
	// span (jika placeholder terpecah) tidak menerima teks apa pun.
	for i := range runs {
		runs[i].text = "" // akumulator baru
	}
	runAt := func(pos int) int {
		for k := 0; k < len(runs); k++ {
			if offs[k+1] > pos {
				return k
			}
		}
		return len(runs) - 1
	}
	copyRange := func(from, to int) {
		if from >= to {
			return
		}
		for k := 0; k < len(runs); k++ {
			start, end := from, to
			if start < offs[k] {
				start = offs[k]
			}
			if end > offs[k+1] {
				end = offs[k+1]
			}
			if start < end {
				runs[k].text += js[start:end]
			}
		}
	}

	pos := 0
	for _, r := range repls {
		copyRange(pos, r.start)
		runs[runAt(r.start)].text += r.value
		pos = r.end
	}
	copyRange(pos, len(js))

	// Rebuild paragraf: semua bagian di luar <w:t> dipertahankan apa adanya.
	var out strings.Builder
	last := 0
	for i, m := range ms {
		out.WriteString(p[last:m[0]])
		if ci := strings.IndexByte(runs[i].full, '>'); ci >= 0 {
			out.WriteString(runs[i].full[:ci+1])
			out.WriteString(runs[i].text)
			out.WriteString("</w:t>")
		} else { // defensif — tak seharusnya terjadi untuk <w:t>...</w:t>
			out.WriteString(runs[i].full)
		}
		last = m[1]
	}
	out.WriteString(p[last:])
	return out.String()
}

// resolveDocxVariables menyalin srcPath (docx) ke dstPath dengan mengganti
// setiap {{key}} pada bagian XML word/ menggunakan values. Placeholder yang
// tidak ada di values dibiarkan apa adanya (tetap tampil sebagai {{key}}).
//
// Placeholder ditangani baik yang utuh dalam satu <w:t> maupun yang dipecah
// Word menjadi beberapa run — nilai dimasukkan ke run pertama, run di tengah
// dikosongkan (format placeholder mengikuti run pertama).
func resolveDocxVariables(srcPath, dstPath string, values map[string]string) error {
	zr, err := zip.OpenReader(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open template docx: %w", err)
	}
	defer zr.Close()

	out, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create resolved docx: %w", err)
	}
	defer out.Close()
	zw := zip.NewWriter(out)

	for _, f := range zr.File {
		modified := strings.HasPrefix(f.Name, "word/") && strings.HasSuffix(f.Name, ".xml")
		if !modified {
			// Salin bagian yang tidak diubah secara mentah (raw) — tanpa
			// decompress/re-compress dan tanpa verifikasi CRC. Bagian seperti
			// [Content_Types].xml kadang punya CRC salah dari penulisnya;
			// menyalin mentah membuat output tetap valid tanpa menolak file.
			if err := zw.Copy(f); err != nil {
				return fmt.Errorf("failed to copy docx part %s: %w", f.Name, err)
			}
			continue
		}

		data, err := readZipEntry(f)
		if err != nil {
			return fmt.Errorf("failed to read docx part %s: %w", f.Name, err)
		}
		s := resolveDocxXML(string(data), values)
		w, err := zw.CreateHeader(&f.FileHeader)
		if err != nil {
			return fmt.Errorf("failed to create docx part %s: %w", f.Name, err)
		}
		if _, err := w.Write([]byte(s)); err != nil {
			return fmt.Errorf("failed to write docx part %s: %w", f.Name, err)
		}
	}
	if err := zw.Close(); err != nil {
		return fmt.Errorf("failed to finalize resolved docx: %w", err)
	}
	return nil
}
