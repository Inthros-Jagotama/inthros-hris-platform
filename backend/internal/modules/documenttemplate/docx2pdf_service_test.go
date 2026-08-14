package documenttemplate

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// TestDocx2pdfPDFServiceConvertsDocx verifies the pure-Go engine actually
// produces a PDF from our minimal test docx (no LibreOffice needed).
func TestDocx2pdfPDFServiceConvertsDocx(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.docx")
	if err := os.WriteFile(in, makeDocx(t, "Nama: {{employee.name}}"), 0644); err != nil {
		t.Fatalf("write docx: %v", err)
	}
	out := filepath.Join(dir, "out.pdf")

	svc := NewDocx2pdfPDFService(0)
	if err := svc.ConvertDOCXToPDF(context.Background(), in, out); err != nil {
		t.Fatalf("convert with docx2pdf failed: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(data) < 5 || string(data[:5]) != "%PDF-" {
		t.Fatalf("expected PDF output, got %d bytes starting with %q", len(data), data[:min(len(data), 10)])
	}
}
