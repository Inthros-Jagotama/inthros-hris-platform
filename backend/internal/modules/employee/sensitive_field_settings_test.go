package employee

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func setupSensitiveFieldTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&SensitiveFieldSetting{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	for _, d := range SensitiveFieldRegistry {
		db.Create(&SensitiveFieldSetting{ID: uuid.New().String(), FieldKey: d.Key, IsEncryptionEnabled: false})
	}
	return db
}

func TestListSensitiveFieldSettings(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	settings, err := repo.ListSensitiveFieldSettings(context.Background())
	if err != nil {
		t.Fatalf("ListSensitiveFieldSettings() error = %v", err)
	}
	if len(settings) != len(SensitiveFieldRegistry) {
		t.Fatalf("got %d settings, want %d", len(settings), len(SensitiveFieldRegistry))
	}
}

func TestSetSensitiveFieldEnabled(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	ctx := context.Background()

	if err := repo.SetSensitiveFieldEnabled(ctx, "employee.nik", true, nil); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	var setting SensitiveFieldSetting
	db.Where("field_key = ?", "employee.nik").First(&setting)
	if !setting.IsEncryptionEnabled {
		t.Error("expected employee.nik encryption to be enabled")
	}
}

func TestService_IsFieldEncryptionEnabled(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)
	ctx := context.Background()

	enabled, err := svc.IsFieldEncryptionEnabled(ctx, "employee.nik")
	if err != nil {
		t.Fatalf("IsFieldEncryptionEnabled() error = %v", err)
	}
	if enabled {
		t.Error("expected employee.nik to default to disabled")
	}

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	enabled, err = svc.IsFieldEncryptionEnabled(ctx, "employee.nik")
	if err != nil {
		t.Fatalf("IsFieldEncryptionEnabled() error = %v", err)
	}
	if !enabled {
		t.Error("expected employee.nik to be enabled after toggling")
	}
}

func TestService_SetSensitiveFieldEnabled_UnknownKey(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	logger, _ := zap.NewDevelopment()
	svc := NewService(repo, logger)

	err := svc.SetSensitiveFieldEnabled(context.Background(), "not.a.real.field", true)
	if err == nil {
		t.Error("expected error for unknown field key")
	}
}
