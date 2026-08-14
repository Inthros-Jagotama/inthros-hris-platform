package numbering

import (
	"context"
	"testing"

	sqlite "github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) *Service {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&DocumentNumberingSetting{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	seed := []DocumentNumberingSetting{
		{ID: "11111111-1111-1111-1111-111111111111", DocumentType: DocumentTypeEmployeeMovement, FormatTemplate: "SK/{sequence:3}/{year}", ResetPeriod: ResetPeriodYearly},
		{ID: "22222222-2222-2222-2222-222222222222", DocumentType: DocumentTypeEmployeeContract, FormatTemplate: "CTR/{sequence:3}/{year}", ResetPeriod: ResetPeriodNever},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	resolver := func(ctx context.Context) (*gorm.DB, error) { return db, nil }
	return NewService(resolver, zap.NewNop())
}

func TestServiceGenerateIncrementsSequence(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	first, err := svc.Generate(ctx, DocumentTypeEmployeeMovement)
	if err != nil {
		t.Fatalf("Generate #1: %v", err)
	}
	second, err := svc.Generate(ctx, DocumentTypeEmployeeMovement)
	if err != nil {
		t.Fatalf("Generate #2: %v", err)
	}
	if first == second {
		t.Fatalf("expected distinct numbers, got %q twice", first)
	}
}

func TestServicePreviewDoesNotMutateSequence(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	preview, err := svc.Preview(ctx, DocumentTypeEmployeeContract)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	generated, err := svc.Generate(ctx, DocumentTypeEmployeeContract)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if preview != generated {
		t.Fatalf("Preview() = %q, first Generate() = %q — preview should predict the next number exactly", preview, generated)
	}
}

func TestServiceUpdateValidatesInput(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	if _, err := svc.Update(ctx, "not_a_real_type", "X/{sequence}", ResetPeriodYearly); err != ErrInvalidDocumentType {
		t.Fatalf("expected ErrInvalidDocumentType, got %v", err)
	}
	if _, err := svc.Update(ctx, DocumentTypeEmployeeMovement, "X/{sequence}", "weekly"); err != ErrInvalidResetPeriod {
		t.Fatalf("expected ErrInvalidResetPeriod, got %v", err)
	}

	updated, err := svc.Update(ctx, DocumentTypeEmployeeMovement, "X/{sequence}", ResetPeriodMonthly)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.FormatTemplate != "X/{sequence}" || updated.ResetPeriod != ResetPeriodMonthly {
		t.Fatalf("Update did not persist changes: %+v", updated)
	}
}

func TestServiceListReturnsBothTypes(t *testing.T) {
	svc := newTestService(t)
	items, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 settings, got %d", len(items))
	}
}
