package documenttemplate

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// PDFService mengabstraksi konversi DOCX → PDF (spec §15.1). Abstraksi ini
// memungkinkan preview dan generate document memakai pipeline yang sama, dan
// memudahkan test dengan implementasi mock.
type PDFService interface {
	ConvertDOCXToPDF(ctx context.Context, inputPath, outputPath string) error
}

// LibreOfficePDFService mengonversi DOCX → PDF via LibreOffice headless:
//
//	<bin> --headless --convert-to pdf --outdir <dir> <file.docx>
//
// LibreOffice menulis hasil di --outdir dengan nama = basename input; hasilnya
// lalu di-rename ke outputPath yang diminta.
type LibreOfficePDFService struct {
	binPath string
	timeout time.Duration
}

// commonLibreOfficeBinaries mengembalikan kandidat nama binary per platform.
// Di Windows executable-nya bernama soffice.exe (bukan libreoffice) dan
// installer-nya TIDAK menambahkan ke PATH — jadi path penuh dari lokasi instal
// standar dicoba dulu, baru nama di PATH.
func commonLibreOfficeBinaries() []string {
	if runtime.GOOS == "windows" {
		return []string{
			`C:\Program Files\LibreOffice\program\soffice.exe`,
			`C:\Program Files (x86)\LibreOffice\program\soffice.exe`,
			`soffice.exe`,
			`soffice`,
			`libreoffice`,
		}
	}
	if runtime.GOOS == "darwin" {
		return []string{
			"/Applications/LibreOffice.app/Contents/MacOS/soffice",
			"soffice",
			"libreoffice",
		}
	}
	return []string{
		"/usr/bin/libreoffice",
		"/usr/bin/soffice",
		"/usr/local/bin/libreoffice",
		"/usr/local/bin/soffice",
		"libreoffice",
		"soffice",
	}
}

// resolveLibreOfficeBinary menentukan binary yang dipakai:
//   1. binPath eksplisit dari config (storage.libreoffice_path) bila ada;
//   2. kandidat umum per platform yang benar-benar ada (file di disk atau
//      ter-resolve via PATH).
// Kembalian kedua = error bila tidak ada satupun yang ditemukan.
func resolveLibreOfficeBinary(binPath string) (string, error) {
	if binPath != "" {
		if _, err := exec.LookPath(binPath); err == nil {
			return binPath, nil
		}
		if _, err := os.Stat(binPath); err == nil {
			return binPath, nil
		}
		return "", fmt.Errorf("configured libreoffice binary %q not found", binPath)
	}
	for _, cand := range commonLibreOfficeBinaries() {
		if _, err := exec.LookPath(cand); err == nil {
			return cand, nil
		}
		if _, err := os.Stat(cand); err == nil {
			return cand, nil
		}
	}
	return "", fmt.Errorf(
		"LibreOffice not installed or not found. Install LibreOffice, then set HRIS_STORAGE_LIBREOFFICE_PATH (or storage.libreoffice_path) to the soffice binary (e.g. C:\\Program Files\\LibreOffice\\program\\soffice.exe)")
}

// NewLibreOfficePDFService membuat service. binPath kosong → auto-detect
// (PATH + lokasi instal umum per platform).
func NewLibreOfficePDFService(binPath string) *LibreOfficePDFService {
	return &LibreOfficePDFService{
		binPath: binPath,
		timeout: 60 * time.Second,
	}
}

// BinaryPath mengembalikan path binary yang akan dipakai (hasil resolve, atau
// binPath eksplisit bila tidak bisa di-resolve). Dipakai test & logging.
func (s *LibreOfficePDFService) BinaryPath() string {
	if s.binPath == "" {
		return "(auto-detect)"
	}
	return s.binPath
}

func (s *LibreOfficePDFService) ConvertDOCXToPDF(ctx context.Context, inputPath, outputPath string) error {
	bin, err := resolveLibreOfficeBinary(s.binPath)
	if err != nil {
		return err
	}

	outDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create pdf output dir: %w", err)
	}

	cmdCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, bin,
		"--headless",
		"--norestore",
		"--convert-to", "pdf",
		"--outdir", outDir,
		inputPath,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		if cmdCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("libreoffice conversion timed out after %s: %s", s.timeout, out)
		}
		return fmt.Errorf("libreoffice conversion failed (%s): %w: %s", bin, err, out)
	}

	// LibreOffice menghasilkan <outdir>/<basename-input>.pdf
	generated := filepath.Join(outDir, filepath.Base(inputPath[:len(inputPath)-len(filepath.Ext(inputPath))])+".pdf")
	if _, err := os.Stat(generated); err != nil {
		return fmt.Errorf("libreoffice did not produce expected pdf %s: %w", generated, err)
	}
	if generated != outputPath {
		if err := os.Rename(generated, outputPath); err != nil {
			return fmt.Errorf("failed to move generated pdf to %s: %w", outputPath, err)
		}
	}
	return nil
}
