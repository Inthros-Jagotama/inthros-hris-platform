package documenttemplate

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractDocxPlaceholders(t *testing.T) {
	docx := makeDocx(t, "Nama: {{employee.name}} NIK: {{employee.nik}} {{employee.name}}")
	placeholders, err := extractDocxPlaceholders(docx)
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if len(placeholders) != 2 {
		t.Fatalf("expected 2 unique placeholders, got %v", placeholders)
	}
	if placeholders[0] != "employee.name" || placeholders[1] != "employee.nik" {
		t.Fatalf("unexpected placeholder list: %v", placeholders)
	}
}

func TestExtractDocxPlaceholdersRejectsNonZip(t *testing.T) {
	if _, err := extractDocxPlaceholders([]byte("not a zip")); err == nil {
		t.Fatal("expected error for non-zip input")
	}
}

func TestUnknownPlaceholders(t *testing.T) {
	unknown := unknownPlaceholders([]string{"employee.name", "company.name", "signature.author", "custom.field"})
	if len(unknown) != 2 || unknown[0] != "custom.field" || unknown[1] != "signature.author" {
		t.Fatalf("expected custom.field & signature.author unknown, got %v", unknown)
	}
}

func TestResolveDocxVariables(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src.docx")
	if err := os.WriteFile(src, makeDocx(t, "Nama: {{employee.name}} Alamat: {{company.address}}"), 0644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	dst := filepath.Join(dir, "dst.docx")
	values := map[string]string{
		"employee.name":  "Asep Ruswanda",
		"company.address": "PT & Co <Jl. Merdeka>",
	}
	if err := resolveDocxVariables(src, dst, values); err != nil {
		t.Fatalf("resolve: %v", err)
	}

	// Baca ulang hasil → placeholder ter-replace, nilai di-escape XML.
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	// Baca ulang document.xml dari hasil untuk memastikan placeholder ter-resolve.
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open resolved docx: %v", err)
	}
	found := ""
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open part: %v", err)
		}
		buf, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		found = string(buf)
	}
	if found == "" {
		t.Fatal("document.xml not found in resolved docx")
	}
	if !strings.Contains(found, "Asep Ruswanda") {
		t.Fatalf("expected employee.name resolved, got: %s", found)
	}
	if !strings.Contains(found, "PT &amp; Co &lt;Jl. Merdeka&gt;") {
		t.Fatalf("expected company.address resolved with XML escape, got: %s", found)
	}
	if strings.Contains(found, "{{") {
		t.Fatalf("expected all placeholders resolved, still contains {{ in: %s", found)
	}
}

// makeDocxWithXML membangun .docx mini dengan document.xml sesuai yang diberikan
// (untuk test kasus run terpecah ala Word).
func makeDocxWithXML(t *testing.T, documentXML string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create docx part: %v", err)
	}
	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:body>` + documentXML + `</w:body></w:document>`
	if _, err := w.Write([]byte(xml)); err != nil {
		t.Fatalf("write docx part: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close docx zip: %v", err)
	}
	return buf.Bytes()
}

// resolveDocxXMLText menjalankan resolveDocxVariables pada .docx yang dibangun
// dari documentXML lalu mengembalikan isi word/document.xml hasilnya.
func resolveDocxXMLText(t *testing.T, documentXML string, values map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	src := filepath.Join(dir, "src.docx")
	if err := os.WriteFile(src, makeDocxWithXML(t, documentXML), 0644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	dst := filepath.Join(dir, "dst.docx")
	if err := resolveDocxVariables(src, dst, values); err != nil {
		t.Fatalf("resolve: %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("open resolved docx: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open part: %v", err)
		}
		buf, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		return string(buf)
	}
	t.Fatal("document.xml not found")
	return ""
}

// TestResolveDocxVariablesSplitRuns: Word memecah {{contract.number}} menjadi
// beberapa <w:t> ({{ di satu run, key di run lain, }} di run lain — pola asli
// dari template pengguna). Semua variasi harus ter-resolve.
func TestResolveDocxVariablesSplitRuns(t *testing.T) {
	// Pola asli dari template Word: {{ | contract.number | }} di run terpisah
	// dengan atribut/format berbeda per run.
	body := `<w:p><w:r><w:t xml:space="preserve">Nomor: {{</w:t></w:r>` +
		`<w:proofErr w:type="spellStart"/><w:r><w:rPr><w:b/></w:rPr><w:t>contract.number</w:t></w:r>` +
		`<w:proofErr w:type="spellEnd"/><w:r><w:t>}} </w:t></w:r><w:r><w:t>NIK: {{employee.nik}}</w:t></w:r></w:p>`
	out := resolveDocxXMLText(t, body, map[string]string{
		"contract.number": "CTR-2026-001",
		"employee.nik":    "199001012015011001",
	})
	if strings.Contains(out, "{{") {
		t.Fatalf("expected all placeholders resolved, still contains {{ in: %s", out)
	}
	if !strings.Contains(out, "CTR-2026-001") {
		t.Fatalf("expected contract.number resolved, got: %s", out)
	}
	if !strings.Contains(out, "199001012015011001") {
		t.Fatalf("expected employee.nik resolved, got: %s", out)
	}
	// Nilai harus berada di run pertama dari placeholder (format run itu dipakai)
	// dan run tengah dikosongkan — struktur XML lain dipertahankan.
	if !strings.Contains(out, `<w:t xml:space="preserve">Nomor: CTR-2026-001</w:t>`) {
		t.Fatalf("expected value placed in first run, got: %s", out)
	}
	if !strings.Contains(out, `<w:t></w:t>`) {
		t.Fatalf("expected middle runs emptied, got: %s", out)
	}
}

// TestResolveDocxVariablesSplitPerChar: kasus terparah — Word memecah per
// karakter. Hasil akhir tetap satu nilai utuh di run pertama.
func TestResolveDocxVariablesSplitPerChar(t *testing.T) {
	var sb strings.Builder
	sb.WriteString(`<w:p>`)
	for _, ch := range "{{contract.start_date}}" {
		sb.WriteString(`<w:r><w:t>`)
		sb.WriteString(string(ch))
		sb.WriteString(`</w:t></w:r>`)
	}
	sb.WriteString(`</w:p>`)
	out := resolveDocxXMLText(t, sb.String(), map[string]string{"contract.start_date": "2026-01-01"})
	if strings.Contains(out, "{{") || strings.Contains(out, "start_date") {
		t.Fatalf("expected placeholder fully resolved, got: %s", out)
	}
	if !strings.Contains(out, "2026-01-01") {
		t.Fatalf("expected start_date resolved, got: %s", out)
	}
}

// TestResolveDocxVariablesMixed: placeholder utuh + terpecah + tanpa nilai
// (tidak ada di values) dalam satu paragraf. Tanpa nilai → dibiarkan utuh.
func TestResolveDocxVariablesMixed(t *testing.T) {
	body := `<w:p><w:r><w:t>Nama: {{employee.name}} | </w:t></w:r>` +
		`<w:r><w:t>{{contract.num</w:t></w:r><w:r><w:t>ber}} | {{unknown.key}}</w:t></w:r></w:p>`
	out := resolveDocxXMLText(t, body, map[string]string{
		"employee.name":  "Asep",
		"contract.number": "CTR-1",
	})
	// Teks tampil (gabungan semua <w:t>) harus ter-resolve; nilai placeholder
	// terpecah masuk ke run pertamanya sehingga di XML mentah teks terpisah run.
	if joined := paragraphJoinedText(out); joined != "Nama: Asep | CTR-1 | {{unknown.key}}" {
		t.Fatalf("expected intact+split resolved, got: %q", joined)
	}
	if !strings.Contains(out, "{{unknown.key}}") {
		t.Fatalf("expected unknown placeholder left intact, got: %s", out)
	}
}

// TestExtractDocxPlaceholdersSplitRuns: deteksi placeholder yang terpecah antar
// run — konsisten dengan resolusi (sebelumnya hanya mendeteksi yang utuh).
func TestExtractDocxPlaceholdersSplitRuns(t *testing.T) {
	body := `<w:p><w:r><w:t>Nomor: {{</w:t></w:r><w:r><w:t>contract.number</w:t></w:r>` +
		`<w:r><w:t>}}</w:t></w:r></w:p><w:p><w:r><w:t>NIK: {{employee.nik}}</w:t></w:r></w:p>`
	placeholders, err := extractDocxPlaceholders(makeDocxWithXML(t, body))
	if err != nil {
		t.Fatalf("extract: %v", err)
	}
	if len(placeholders) != 2 {
		t.Fatalf("expected 2 placeholders (incl. split one), got %v", placeholders)
	}
	if placeholders[0] != "contract.number" || placeholders[1] != "employee.nik" {
		t.Fatalf("unexpected placeholder list: %v", placeholders)
	}
}

func TestSampleDataCoversRegistry(t *testing.T) {
	keys := registryKeys()
	samples := sampleData()
	for k := range keys {
		if _, ok := samples[k]; !ok {
			t.Fatalf("sample data missing key %q", k)
		}
	}
}

// makeDocxWithBadCRC membangun .docx mini di mana entry [Content_Types].xml
// ditulis dengan CRC32 salah (zip.Writer.CreateRaw + FileHeader.CRC32 tidak
// akurat) — mensimulasikan .docx dari penulis (WPS/Word versi tertentu) yang
// menulis CRC tidak valid namun isi file tetap benar dan bisa dibuka Word.
func makeDocxWithBadCRC(t *testing.T, documentXML string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	// [Content_Types].xml — CRC sengaja salah.
	ct := []byte(`<?xml version="1.0"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"></Types>`)
	fh := &zip.FileHeader{Name: "[Content_Types].xml", Method: zip.Store}
	fh.CRC32 = 0xDEADBEEF // nilai salah
	fh.UncompressedSize64 = uint64(len(ct))
	fh.CompressedSize64 = uint64(len(ct))
	w, err := zw.CreateRaw(fh)
	if err != nil {
		t.Fatalf("create raw entry: %v", err)
	}
	if _, err := w.Write(ct); err != nil {
		t.Fatalf("write raw entry: %v", err)
	}

	// word/document.xml — normal (isi yang memuat placeholder).
	dw, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create document.xml: %v", err)
	}
	if _, err := dw.Write([]byte(documentXML)); err != nil {
		t.Fatalf("write document.xml: %v", err)
	}

	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func TestResolveDocxVariablesToleratesBadCRC(t *testing.T) {
	dir := t.TempDir()
	// Document XML realistis (teks di dalam paragraf) — poin test ini adalah
	// toleransi CRC, bukan struktur XML.
	bad := makeDocxWithBadCRC(t, `<w:p><w:r><w:t>Nama: {{employee.name}}</w:t></w:r></w:p>`)

	// Sanity: membaca [Content_Types].xml dari zip ini memang memunculkan
	// checksum error lewat jalur lama (zip.File.Open + io.ReadAll).
	zr, err := zip.NewReader(bytes.NewReader(bad), int64(len(bad)))
	if err != nil {
		t.Fatalf("open zip: %v", err)
	}
	crcFailed := false
	for _, f := range zr.File {
		if f.Name != "[Content_Types].xml" {
			continue
		}
		rc, err := f.Open()
		if err == nil {
			_, rerr := io.ReadAll(rc)
			rc.Close()
			crcFailed = rerr != nil
		}
	}
	if !crcFailed {
		t.Fatal("sanity: expected checksum error when reading [Content_Types].xml")
	}

	// Jalur yang tadinya gagal (zip: checksum error) sekarang harus sukses.
	src := filepath.Join(dir, "src-badcrc.docx")
	if err := os.WriteFile(src, bad, 0644); err != nil {
		t.Fatalf("write src: %v", err)
	}
	dst := filepath.Join(dir, "dst.docx")
	if err := resolveDocxVariables(src, dst, map[string]string{"employee.name": "Asep"}); err != nil {
		t.Fatalf("resolve with bad CRC should succeed, got: %v", err)
	}

	// Hasil harus valid zip dan placeholder ter-resolve.
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("read dst: %v", err)
	}
	zr2, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("resolved docx not a valid zip: %v", err)
	}
	found := ""
	for _, f := range zr2.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open resolved part: %v", err)
		}
		buf, rerr := io.ReadAll(rc)
		rc.Close()
		if rerr != nil {
			t.Fatalf("read resolved part: %v", rerr)
		}
		found = string(buf)
	}
	if !strings.Contains(found, "Asep") {
		t.Fatalf("expected placeholder resolved in output, got: %s", found)
	}
}
