package documenttemplate

import (
	"context"
	"fmt"
	"time"

	docx2pdf "github.com/bobyeoh/docx2pdf-go"
)

// Docx2pdfPDFService mengonversi DOCX → PDF menggunakan library pure-Go
// github.com/bobyeoh/docx2pdf-go (MIT, tanpa LibreOffice/MS Word/JVM/CGO).
//
// Ini adalah implementasi ALTERNATIF dari PDFService — opsi kedua selain
// LibreOfficePDFService. Keuntungan: tanpa dependency eksternal (tidak perlu
// install LibreOffice di server), binary kecil, sangat cepat (~13–21 ms per
// dokumen). Kekurangan: bukan pixel-perfect Word (fidelity menurun pada
// floating-image wrap, SmartArt, math) dan library masih pre-1.0 (API bisa
// berubah antar minor release — pin tag di go.mod).
//
// Pilih engine mana yang dipakai via config storage.pdf_engine
// ("libreoffice" default | "docx2pdf"). Keduanya tetap tersedia; service
// dibuat di main.go sesuai pilihan.
type Docx2pdfPDFService struct {
	timeout time.Duration
}

// NewDocx2pdfPDFService membuat service konversi pure-Go. timeout membatasi
// durasi konversi per dokumen (default 60s bila 0).
func NewDocx2pdfPDFService(timeout time.Duration) *Docx2pdfPDFService {
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	return &Docx2pdfPDFService{timeout: timeout}
}

// ConvertDOCXToPDF mengonversi inputPath (docx) → outputPath (pdf) murni di Go.
// Output ditulis langsung ke path tujuan (tidak perlu rename seperti LibreOffice
// yang menamai hasil berdasarkan basename input).
func (s *Docx2pdfPDFService) ConvertDOCXToPDF(ctx context.Context, inputPath, outputPath string) error {
	convCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := docx2pdf.ConvertContext(convCtx, inputPath, outputPath, docx2pdf.Options{
		// Author diisi agar metadata PDF konsisten (kolom AUTHOR di Word).
		Author: "HRIS Platform",
	}); err != nil {
		if convCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("docx2pdf conversion timed out after %s", s.timeout)
		}
		return fmt.Errorf("docx2pdf conversion failed: %w", err)
	}
	return nil
}
