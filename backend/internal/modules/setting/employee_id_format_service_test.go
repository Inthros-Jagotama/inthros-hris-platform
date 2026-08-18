package setting

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	sqlite "github.com/glebarez/sqlite"

	"github.com/inthros/hris-platform/internal/pkg/numbering"
)

// newEmployeeIDFormatTestService sets up an isolated in-memory SQLite DB
// (private, not the shared-cache DSN used elsewhere in this package) so
// this singleton-row table doesn't leak across tests — same pattern as
// setupNumberingTestRouter in numbering_handler_test.go.
func newEmployeeIDFormatTestService(t *testing.T) *EmployeeIDFormatService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	t.Cleanup(func() {
		if sqlDB, dbErr := db.DB(); dbErr == nil {
			_ = sqlDB.Close()
		}
	})
	if err := db.AutoMigrate(&EmployeeIDFormatSetting{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	resolver := func(ctx context.Context) (*gorm.DB, error) { return db, nil }
	repo := NewEmployeeIDFormatRepository(resolver)
	return NewEmployeeIDFormatService(repo, resolver, zap.NewNop())
}

func TestEmployeeIDFormatService_GetSetting_CreatesDefaultRow(t *testing.T) {
	svc := newEmployeeIDFormatTestService(t)
	ctx := context.Background()

	s1, err := svc.GetSetting(ctx)
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if s1.GenerationMode != EmployeeIDGenerationModeManual {
		t.Fatalf("expected default mode MANUAL, got %q", s1.GenerationMode)
	}
	if s1.FormatTemplate != defaultEmployeeIDFormatTemplate {
		t.Fatalf("expected default template %q, got %q", defaultEmployeeIDFormatTemplate, s1.FormatTemplate)
	}

	// Second read must return the SAME row (singleton), not create another.
	s2, err := svc.GetSetting(ctx)
	if err != nil {
		t.Fatalf("GetSetting (2nd): %v", err)
	}
	if s1.ID != s2.ID {
		t.Fatalf("expected singleton row, got different IDs %q vs %q", s1.ID, s2.ID)
	}
}

func TestEmployeeIDFormatService_UpdateSetting_ValidatesMode(t *testing.T) {
	svc := newEmployeeIDFormatTestService(t)
	ctx := context.Background()

	if _, err := svc.UpdateSetting(ctx, "BOGUS", "EMP-{sequence:4}", numbering.ResetPeriodYearly); !errors.Is(err, ErrInvalidEmployeeIDGenerationMode) {
		t.Fatalf("expected ErrInvalidEmployeeIDGenerationMode, got %v", err)
	}
	if _, err := svc.UpdateSetting(ctx, EmployeeIDGenerationModeAuto, "EMP-{sequence:4}", "bogus-period"); !errors.Is(err, ErrInvalidEmployeeIDResetPeriod) {
		t.Fatalf("expected ErrInvalidEmployeeIDResetPeriod, got %v", err)
	}

	updated, err := svc.UpdateSetting(ctx, EmployeeIDGenerationModeAuto, "EMP-{sequence:4}", numbering.ResetPeriodNever)
	if err != nil {
		t.Fatalf("UpdateSetting: %v", err)
	}
	if updated.GenerationMode != EmployeeIDGenerationModeAuto || updated.ResetPeriod != numbering.ResetPeriodNever {
		t.Fatalf("update did not persist: %+v", updated)
	}
}

func TestEmployeeIDFormatService_Preview_DoesNotMutateSequence(t *testing.T) {
	svc := newEmployeeIDFormatTestService(t)
	svc.now = func() time.Time { return time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC) }
	ctx := context.Background()

	if _, err := svc.UpdateSetting(ctx, EmployeeIDGenerationModeAuto, "EMP-{year}-{sequence:3}", numbering.ResetPeriodYearly); err != nil {
		t.Fatalf("UpdateSetting: %v", err)
	}

	preview1, err := svc.Preview(ctx)
	if err != nil {
		t.Fatalf("Preview: %v", err)
	}
	if preview1 != "EMP-2026-001" {
		t.Fatalf("expected EMP-2026-001, got %q", preview1)
	}
	// Calling Preview again must return the identical result — no mutation.
	preview2, err := svc.Preview(ctx)
	if err != nil {
		t.Fatalf("Preview (2nd): %v", err)
	}
	if preview2 != preview1 {
		t.Fatalf("Preview must be idempotent, got %q then %q", preview1, preview2)
	}

	s, err := svc.GetSetting(ctx)
	if err != nil {
		t.Fatalf("GetSetting: %v", err)
	}
	if s.LastSequence != 0 {
		t.Fatalf("Preview must not mutate last_sequence, got %d", s.LastSequence)
	}
}

func TestEmployeeIDFormatService_Generate_IncrementsAndResetsOnPeriodRollover(t *testing.T) {
	svc := newEmployeeIDFormatTestService(t)
	current := time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return current }
	ctx := context.Background()

	if _, err := svc.UpdateSetting(ctx, EmployeeIDGenerationModeAuto, "EMP-{year}-{sequence:3}", numbering.ResetPeriodYearly); err != nil {
		t.Fatalf("UpdateSetting: %v", err)
	}

	id1, err := svc.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate #1: %v", err)
	}
	if id1 != "EMP-2026-001" {
		t.Fatalf("expected EMP-2026-001, got %q", id1)
	}

	id2, err := svc.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate #2: %v", err)
	}
	if id2 != "EMP-2026-002" {
		t.Fatalf("expected EMP-2026-002, got %q", id2)
	}

	// Roll over to next year — sequence must reset to 1.
	current = time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	id3, err := svc.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate #3: %v", err)
	}
	if id3 != "EMP-2027-001" {
		t.Fatalf("expected reset to EMP-2027-001, got %q", id3)
	}
}
