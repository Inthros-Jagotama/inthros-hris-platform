package documenttemplate

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}
	if err := db.AutoMigrate(&DocumentTemplate{}, &DocumentTemplateVersion{}, &DocumentTemplateAudit{}, &GeneratedDocument{}); err != nil {
		panic(fmt.Sprintf("failed to automigrate: %v", err))
	}
	sqlDB, _ := db.DB()
	return db, func() { sqlDB.Close() }
}

func testDBResolver(db *gorm.DB) func(ctx context.Context) (*gorm.DB, error) {
	return func(ctx context.Context) (*gorm.DB, error) { return db, nil }
}

func newTestRepo(db *gorm.DB) *Repository {
	return NewRepository(testDBResolver(db))
}

func uuidStr() string { return uuid.New().String() }

// makeDocx membangun file .docx minimal (zip OOXML) dengan isi dokumen berisi
// konten yang diberikan — dipakai test upload agar placeholder detection punya
// file yang benar-benar valid sebagai zip/XML.
func makeDocx(t *testing.T, body string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatalf("create docx part: %v", err)
	}
	xml := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:body><w:p><w:r><w:t>` + body + `</w:t></w:r></w:p></w:body></w:document>`
	if _, err := w.Write([]byte(xml)); err != nil {
		t.Fatalf("write docx part: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close docx zip: %v", err)
	}
	return buf.Bytes()
}

func createTestTemplate(db *gorm.DB, code, documentType, status string) *DocumentTemplate {
	tpl := &DocumentTemplate{
		ID:           uuidStr(),
		Name:         code,
		Code:         code,
		DocumentType: documentType,
		Status:       status,
		IsActive:     true,
	}
	if err := db.Create(tpl).Error; err != nil {
		panic(fmt.Sprintf("failed to create test template: %v", err))
	}
	return tpl
}
