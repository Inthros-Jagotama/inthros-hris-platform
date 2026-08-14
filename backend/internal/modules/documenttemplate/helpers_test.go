package documenttemplate

import (
	"context"
	"fmt"

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

func createTestTemplate(db *gorm.DB, code, documentType, status string, isDefault bool) *DocumentTemplate {
	tpl := &DocumentTemplate{
		ID:           uuidStr(),
		Name:         code,
		Code:         code,
		DocumentType: documentType,
		Status:       status,
		IsDefault:    isDefault,
		IsActive:     true,
	}
	if err := db.Create(tpl).Error; err != nil {
		panic(fmt.Sprintf("failed to create test template: %v", err))
	}
	return tpl
}
