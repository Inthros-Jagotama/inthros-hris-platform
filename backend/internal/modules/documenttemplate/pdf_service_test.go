package documenttemplate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveLibreOfficeBinaryExplicitPath: binPath eksplisit yang ada di disk
// harus dipakai apa adanya.
func TestResolveLibreOfficeBinaryExplicitPath(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "soffice.exe")
	if err := os.WriteFile(bin, []byte("dummy"), 0755); err != nil {
		t.Fatalf("create fake binary: %v", err)
	}

	got, err := resolveLibreOfficeBinary(bin)
	if err != nil {
		t.Fatalf("expected resolve to succeed with existing explicit path, got %v", err)
	}
	if got != bin {
		t.Fatalf("expected %q, got %q", bin, got)
	}
}

// TestResolveLibreOfficeBinaryExplicitPathMissing: binPath eksplisit yang tidak
// ada → error menyebut path yang dikonfigurasi.
func TestResolveLibreOfficeBinaryExplicitPathMissing(t *testing.T) {
	_, err := resolveLibreOfficeBinary(filepath.Join(t.TempDir(), "nonexistent-soffice.exe"))
	if err == nil {
		t.Fatal("expected error for missing explicit binary path")
	}
	if !strings.Contains(err.Error(), "configured libreoffice binary") {
		t.Fatalf("expected error to mention configured binary, got: %v", err)
	}
}

// TestResolveLibreOfficeBinaryNotFound: tidak ada binPath & tidak ada kandidat
// → error jelas berisi panduan install/set env.
func TestResolveLibreOfficeBinaryNotFound(t *testing.T) {
	_, err := resolveLibreOfficeBinary("")
	if err == nil {
		t.Fatal("expected error when no binary is available")
	}
	if !strings.Contains(err.Error(), "LibreOffice not installed or not found") {
		t.Fatalf("expected friendly missing-install error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "HRIS_STORAGE_LIBREOFFICE_PATH") {
		t.Fatalf("expected error to mention env var, got: %v", err)
	}
}

// TestLibreOfficePDFServiceConvertMissingBinary: ConvertDOCXToPDF tanpa binary
// → error jelas (bukan exec error mentah).
func TestLibreOfficePDFServiceConvertMissingBinary(t *testing.T) {
	svc := NewLibreOfficePDFService(filepath.Join(t.TempDir(), "does-not-exist.exe"))
	err := svc.ConvertDOCXToPDF(t.Context(), "input.docx", filepath.Join(t.TempDir(), "out.pdf"))
	if err == nil {
		t.Fatal("expected error when binary is missing")
	}
	if !strings.Contains(err.Error(), "configured libreoffice binary") {
		t.Fatalf("expected clear configured-binary error, got: %v", err)
	}
}
