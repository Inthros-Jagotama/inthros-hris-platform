package numbering

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

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

// TestDocumentNumberingSettingJSONKeysAreSnakeCase verifies the actual
// wire-format JSON keys, not just a Go-struct round-trip (which would pass
// even without json tags on the struct fields).
func TestDocumentNumberingSettingJSONKeysAreSnakeCase(t *testing.T) {
	setting := DocumentNumberingSetting{
		ID:             "11111111-1111-1111-1111-111111111111",
		DocumentType:   DocumentTypeEmployeeMovement,
		FormatTemplate: "SK/{sequence:3}/{year}",
		ResetPeriod:    ResetPeriodYearly,
		LastSequence:   1,
		LastResetKey:   "2026",
	}
	b, err := json.Marshal(setting)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	out := string(b)

	wantKeys := []string{
		`"id"`, `"document_type"`, `"format_template"`, `"reset_period"`,
		`"last_sequence"`, `"last_reset_key"`, `"created_at"`, `"updated_at"`,
	}
	for _, key := range wantKeys {
		if !strings.Contains(out, key) {
			t.Errorf("expected JSON output to contain %s, got: %s", key, out)
		}
	}

	badKeys := []string{`"DocumentType"`, `"FormatTemplate"`, `"ResetPeriod"`}
	for _, key := range badKeys {
		if strings.Contains(out, key) {
			t.Errorf("JSON output should not contain PascalCase key %s, got: %s", key, out)
		}
	}
}

// TestServiceGenerateResetsSequenceOnPeriodRollover verifies that when the
// stored last_reset_key no longer matches the current period's key (e.g.
// the year has rolled over for a yearly-reset setting), Generate resets
// LastSequence to 0 before incrementing rather than continuing the old
// count.
func TestServiceGenerateResetsSequenceOnPeriodRollover(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	// Simulate "now" as being in 2025 and generate a few numbers so
	// last_sequence and last_reset_key get populated for that period.
	svc.now = func() time.Time { return time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC) }
	if _, err := svc.Generate(ctx, DocumentTypeEmployeeMovement); err != nil {
		t.Fatalf("Generate (2025) #1: %v", err)
	}
	if _, err := svc.Generate(ctx, DocumentTypeEmployeeMovement); err != nil {
		t.Fatalf("Generate (2025) #2: %v", err)
	}

	// Now simulate the clock rolling over into 2026 — the yearly reset
	// period should kick in and sequence should restart at 1.
	svc.now = func() time.Time { return time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC) }
	result, err := svc.Generate(ctx, DocumentTypeEmployeeMovement)
	if err != nil {
		t.Fatalf("Generate (2026): %v", err)
	}
	if !strings.Contains(result, "001") {
		t.Fatalf("expected sequence to reset to 1 after period rollover, got %q", result)
	}
	if !strings.Contains(result, "2026") {
		t.Fatalf("expected formatted number to reflect the new year, got %q", result)
	}

	var setting DocumentNumberingSetting
	db, _ := svc.getDB(ctx)
	if err := db.Where("document_type = ?", DocumentTypeEmployeeMovement).First(&setting).Error; err != nil {
		t.Fatalf("reload setting: %v", err)
	}
	if setting.LastSequence != 1 {
		t.Fatalf("expected LastSequence to be reset to 1, got %d", setting.LastSequence)
	}
	if setting.LastResetKey != "2026" {
		t.Fatalf("expected LastResetKey to be updated to 2026, got %q", setting.LastResetKey)
	}
}
