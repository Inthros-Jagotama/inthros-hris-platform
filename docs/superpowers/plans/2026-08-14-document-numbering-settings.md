# Document Numbering Settings (Movements & Contracts) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let tenants configure a document-numbering format (per document type) in Settings, and auto-generate `decision_letter_number` (movements) / `contract_number` (contracts) when a record is created with that field left blank, while still allowing manual entry.

**Architecture:** A new standalone Go package `backend/internal/pkg/numbering` owns the `document_numbering_settings` table, the token-formatting logic, and the `Generate`/`Preview`/`List`/`Update` operations (with row-level locking for safe concurrent sequence increments). The existing `setting` module exposes HTTP routes that delegate to `numbering.Service`. The existing `employeemovement` module gets a `SetNumberingService` setter (same pattern as `SetApprovalEngine`) and calls `Generate` inside `CreateMovement`/`CreateContract` when the number field is blank. Frontend adds one new Settings page and tweaks two existing forms.

**Tech Stack:** Go, Gin, GORM (MySQL + Postgres tenant DBs), Vue 3 + PrimeVue (frontend/tenant).

**Spec:** `docs/superpowers/specs/2026-08-14-document-numbering-settings-design.md`

## Global Constraints

- Tenant DB schema changes MUST go through versioned SQL migrations in `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/` — `AutoMigrate` never runs for tenant modules. New migration id: `109`.
- Table columns: `id CHAR(36) PRIMARY KEY` (UUID as string, set in Go), `created_at`/`updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP`, `ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci` (mysql file); mirror in the postgres file using postgres syntax.
- Every mysql migration file has a matching postgres file and both have `.down.sql` counterparts.
- Document types: exactly two — `employee_movement`, `employee_contract`. No per-subtype (promotion/mutation/demotion) split.
- Reset periods: exactly three — `yearly`, `monthly`, `never`.
- Manual override: if the incoming create request already has a non-blank number field, it is used as-is — `Generate` is only called when blank.

---

### Task 1: Tenant migration for `document_numbering_settings`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.down.sql`

**Interfaces:**
- Produces: table `document_numbering_settings` with columns `id, document_type, format_template, reset_period, last_sequence, last_reset_key, created_at, updated_at`, seeded with 2 rows (`employee_movement`, `employee_contract`). Later Go tasks map to this exact table/column set.

- [ ] **Step 1: Write the mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.sql`:
```sql
CREATE TABLE IF NOT EXISTS document_numbering_settings (
    id              CHAR(36) PRIMARY KEY,
    document_type   VARCHAR(50) NOT NULL,
    format_template VARCHAR(255) NOT NULL,
    reset_period    VARCHAR(20) NOT NULL DEFAULT 'yearly',
    last_sequence   INT NOT NULL DEFAULT 0,
    last_reset_key  VARCHAR(16) NOT NULL DEFAULT '',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE KEY uk_doc_numbering_type (document_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO document_numbering_settings (id, document_type, format_template, reset_period, last_sequence, last_reset_key)
VALUES
    (UUID(), 'employee_movement', 'SK/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, ''),
    (UUID(), 'employee_contract', 'CTR/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, '');
```

- [ ] **Step 2: Write the mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.down.sql`:
```sql
DROP TABLE IF EXISTS document_numbering_settings;
```

- [ ] **Step 3: Write the postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.sql`:
```sql
CREATE TABLE IF NOT EXISTS document_numbering_settings (
    id              CHAR(36) PRIMARY KEY,
    document_type   VARCHAR(50) NOT NULL,
    format_template VARCHAR(255) NOT NULL,
    reset_period    VARCHAR(20) NOT NULL DEFAULT 'yearly',
    last_sequence   INT NOT NULL DEFAULT 0,
    last_reset_key  VARCHAR(16) NOT NULL DEFAULT '',
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uk_doc_numbering_type UNIQUE (document_type)
);

INSERT INTO document_numbering_settings (id, document_type, format_template, reset_period, last_sequence, last_reset_key)
VALUES
    (gen_random_uuid(), 'employee_movement', 'SK/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, ''),
    (gen_random_uuid(), 'employee_contract', 'CTR/{sequence:3}/HRIS/{month_roman}/{year}', 'yearly', 0, '');
```

- [ ] **Step 4: Write the postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.down.sql`:
```sql
DROP TABLE IF EXISTS document_numbering_settings;
```

- [ ] **Step 5: Run the tenant migrator against a local/test tenant DB and confirm the table exists**

Run whatever command this repo uses to apply tenant migrations (check `backend/internal/pkg/migrator` for a CLI entrypoint, e.g. `go run ./cmd/migrate tenant up` or similar — confirm the exact command by checking `backend/cmd/` for a migrate command before running). Expected: migration `109` applies cleanly on both a mysql and postgres tenant DB with no errors, and `SELECT * FROM document_numbering_settings;` returns the 2 seeded rows.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.sql backend/internal/pkg/migrator/migrations/tenant/mysql/109_document_numbering_settings.down.sql backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.sql backend/internal/pkg/migrator/migrations/tenant/postgres/109_document_numbering_settings.down.sql
git commit -m "feat: add document_numbering_settings tenant migration"
```

---

### Task 2: `numbering` package — model + token formatter (TDD)

**Files:**
- Create: `backend/internal/pkg/numbering/model.go`
- Create: `backend/internal/pkg/numbering/format.go`
- Test: `backend/internal/pkg/numbering/format_test.go`

**Interfaces:**
- Produces:
  - `type DocumentNumberingSetting struct { ID string; DocumentType string; FormatTemplate string; ResetPeriod string; LastSequence int; LastResetKey string; CreatedAt time.Time; UpdatedAt time.Time }` with `TableName() string { return "document_numbering_settings" }`.
  - Constants: `DocumentTypeEmployeeMovement = "employee_movement"`, `DocumentTypeEmployeeContract = "employee_contract"`.
  - Constants: `ResetPeriodYearly = "yearly"`, `ResetPeriodMonthly = "monthly"`, `ResetPeriodNever = "never"`.
  - `func FormatTemplate(template string, sequence int, now time.Time) string`
  - `func ResetKeyFor(period string, now time.Time) string`

- [ ] **Step 1: Write the failing tests for `FormatTemplate` and `ResetKeyFor`**

`backend/internal/pkg/numbering/format_test.go`:
```go
package numbering

import (
	"testing"
	"time"
)

func TestFormatTemplate(t *testing.T) {
	now := time.Date(2026, time.August, 14, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		template string
		sequence int
		want     string
	}{
		{
			name:     "sequence with padding, month roman, year",
			template: "SK/{sequence:3}/HRIS/{month_roman}/{year}",
			sequence: 7,
			want:     "SK/007/HRIS/VIII/2026",
		},
		{
			name:     "plain sequence and two digit year",
			template: "CTR-{sequence}-{yy}{month}",
			sequence: 42,
			want:     "CTR-42-2608",
		},
		{
			name:     "sequence padding wider than value",
			template: "{sequence:5}",
			sequence: 3,
			want:     "00003",
		},
		{
			name:     "unknown token left literal",
			template: "{prefix}-{sequence:2}",
			sequence: 1,
			want:     "{prefix}-01",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatTemplate(tc.template, tc.sequence, now)
			if got != tc.want {
				t.Errorf("FormatTemplate(%q, %d) = %q, want %q", tc.template, tc.sequence, got, tc.want)
			}
		})
	}
}

func TestResetKeyFor(t *testing.T) {
	now := time.Date(2026, time.August, 14, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		period string
		want   string
	}{
		{ResetPeriodYearly, "2026"},
		{ResetPeriodMonthly, "2026-08"},
		{ResetPeriodNever, ""},
	}

	for _, tc := range cases {
		t.Run(tc.period, func(t *testing.T) {
			got := ResetKeyFor(tc.period, now)
			if got != tc.want {
				t.Errorf("ResetKeyFor(%q) = %q, want %q", tc.period, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail (package doesn't compile yet)**

Run: `cd backend && go test ./internal/pkg/numbering/... -run TestFormatTemplate -v`
Expected: FAIL — build error, `FormatTemplate`/`ResetKeyFor`/`ResetPeriodYearly` etc. undefined.

- [ ] **Step 3: Write `model.go`**

`backend/internal/pkg/numbering/model.go`:
```go
package numbering

import "time"

const (
	DocumentTypeEmployeeMovement = "employee_movement"
	DocumentTypeEmployeeContract = "employee_contract"
)

const (
	ResetPeriodYearly  = "yearly"
	ResetPeriodMonthly = "monthly"
	ResetPeriodNever   = "never"
)

// DocumentNumberingSetting stores the numbering format and running sequence
// for one document type (employee_movement / employee_contract).
type DocumentNumberingSetting struct {
	ID             string    `gorm:"column:id;primaryKey"`
	DocumentType   string    `gorm:"column:document_type"`
	FormatTemplate string    `gorm:"column:format_template"`
	ResetPeriod    string    `gorm:"column:reset_period"`
	LastSequence   int       `gorm:"column:last_sequence"`
	LastResetKey   string    `gorm:"column:last_reset_key"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (DocumentNumberingSetting) TableName() string {
	return "document_numbering_settings"
}
```

- [ ] **Step 4: Write `format.go`**

`backend/internal/pkg/numbering/format.go`:
```go
package numbering

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var romanMonths = [...]string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}

var sequenceTokenRe = regexp.MustCompile(`\{sequence(?::(\d+))?\}`)

// FormatTemplate substitutes numbering tokens in template with concrete
// values for the given sequence number and point in time. Unknown tokens
// are left as literal text.
func FormatTemplate(template string, sequence int, now time.Time) string {
	result := sequenceTokenRe.ReplaceAllStringFunc(template, func(match string) string {
		sub := sequenceTokenRe.FindStringSubmatch(match)
		if sub[1] == "" {
			return strconv.Itoa(sequence)
		}
		width, err := strconv.Atoi(sub[1])
		if err != nil {
			return strconv.Itoa(sequence)
		}
		return fmt.Sprintf("%0*d", width, sequence)
	})

	replacer := strings.NewReplacer(
		"{year}", strconv.Itoa(now.Year()),
		"{yy}", fmt.Sprintf("%02d", now.Year()%100),
		"{month}", fmt.Sprintf("%02d", int(now.Month())),
		"{month_roman}", romanMonths[int(now.Month())-1],
	)
	return replacer.Replace(result)
}

// ResetKeyFor returns the period key used to detect when a sequence should
// reset: the year for "yearly", "YYYY-MM" for "monthly", and a constant
// empty string for "never" (which therefore never resets).
func ResetKeyFor(period string, now time.Time) string {
	switch period {
	case ResetPeriodYearly:
		return strconv.Itoa(now.Year())
	case ResetPeriodMonthly:
		return fmt.Sprintf("%04d-%02d", now.Year(), int(now.Month()))
	default:
		return ""
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/pkg/numbering/... -v`
Expected: PASS — all `TestFormatTemplate` and `TestResetKeyFor` subtests green.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/pkg/numbering/model.go backend/internal/pkg/numbering/format.go backend/internal/pkg/numbering/format_test.go
git commit -m "feat: add numbering package model and token formatter"
```

---

### Task 3: `numbering.Service` — List, Update, Preview, Generate (TDD with sqlite in-memory DB)

**Files:**
- Create: `backend/internal/pkg/numbering/service.go`
- Test: `backend/internal/pkg/numbering/service_test.go`

**Interfaces:**
- Consumes: `DocumentNumberingSetting`, `FormatTemplate`, `ResetKeyFor`, `ResetPeriodYearly`/`Monthly`/`Never`, `DocumentTypeEmployeeMovement`/`EmployeeContract` from Task 2.
- Produces:
  - `type Service struct { ... }` (unexported fields)
  - `func NewService(dbResolver func(ctx context.Context) (*gorm.DB, error), logger *zap.Logger) *Service`
  - `func (s *Service) List(ctx context.Context) ([]DocumentNumberingSetting, error)`
  - `func (s *Service) Update(ctx context.Context, documentType, formatTemplate, resetPeriod string) (*DocumentNumberingSetting, error)` — returns `ErrInvalidDocumentType` / `ErrInvalidResetPeriod` (sentinel errors defined in this file) on bad input.
  - `func (s *Service) Preview(ctx context.Context, documentType string) (string, error)` — read-only, does not mutate `last_sequence`.
  - `func (s *Service) Generate(ctx context.Context, documentType string) (string, error)` — locks the row, increments/resets sequence, persists, returns formatted number. This is what `employeemovement` calls in Task 4.

- [ ] **Step 1: Write the failing tests using an in-memory sqlite DB**

Check `go.mod` / existing tests for the sqlite driver already used in this repo (search for `gorm.io/driver/sqlite` usage in other `*_test.go` files under `backend/internal/modules/`) and reuse the same driver import. Write `backend/internal/pkg/numbering/service_test.go`:
```go
package numbering

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
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
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/pkg/numbering/... -run TestService -v`
Expected: FAIL — `Service`/`NewService`/`ErrInvalidDocumentType`/etc. undefined.

- [ ] **Step 3: Write `service.go`**

`backend/internal/pkg/numbering/service.go`:
```go
package numbering

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrInvalidResetPeriod  = errors.New("invalid reset period")
	ErrSettingNotFound     = errors.New("numbering setting not found")
)

var validDocumentTypes = map[string]bool{
	DocumentTypeEmployeeMovement: true,
	DocumentTypeEmployeeContract: true,
}

var validResetPeriods = map[string]bool{
	ResetPeriodYearly:  true,
	ResetPeriodMonthly: true,
	ResetPeriodNever:   true,
}

type Service struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
	logger     *zap.Logger
	now        func() time.Time
}

func NewService(dbResolver func(ctx context.Context) (*gorm.DB, error), logger *zap.Logger) *Service {
	return &Service{dbResolver: dbResolver, logger: logger, now: time.Now}
}

func (s *Service) getDB(ctx context.Context) (*gorm.DB, error) {
	return s.dbResolver(ctx)
}

func (s *Service) List(ctx context.Context) ([]DocumentNumberingSetting, error) {
	db, err := s.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var items []DocumentNumberingSetting
	if err := db.Order("document_type").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list numbering settings: %w", err)
	}
	return items, nil
}

func (s *Service) Update(ctx context.Context, documentType, formatTemplate, resetPeriod string) (*DocumentNumberingSetting, error) {
	if !validDocumentTypes[documentType] {
		return nil, ErrInvalidDocumentType
	}
	if !validResetPeriods[resetPeriod] {
		return nil, ErrInvalidResetPeriod
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var setting DocumentNumberingSetting
	if err := db.Where("document_type = ?", documentType).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSettingNotFound
		}
		return nil, fmt.Errorf("failed to load numbering setting: %w", err)
	}
	setting.FormatTemplate = formatTemplate
	setting.ResetPeriod = resetPeriod
	if err := db.Save(&setting).Error; err != nil {
		return nil, fmt.Errorf("failed to update numbering setting: %w", err)
	}
	return &setting, nil
}

// Preview formats what the next Generate() call would return, without
// mutating last_sequence — used by the settings UI to show a live example.
func (s *Service) Preview(ctx context.Context, documentType string) (string, error) {
	if !validDocumentTypes[documentType] {
		return "", ErrInvalidDocumentType
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return "", err
	}
	var setting DocumentNumberingSetting
	if err := db.Where("document_type = ?", documentType).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrSettingNotFound
		}
		return "", fmt.Errorf("failed to load numbering setting: %w", err)
	}
	now := s.now()
	nextSeq := setting.LastSequence + 1
	if ResetKeyFor(setting.ResetPeriod, now) != setting.LastResetKey {
		nextSeq = 1
	}
	return FormatTemplate(setting.FormatTemplate, nextSeq, now), nil
}

// Generate atomically increments (and resets, if the period rolled over)
// the sequence for documentType and returns the formatted number. Must be
// safe under concurrent callers, hence the row lock.
func (s *Service) Generate(ctx context.Context, documentType string) (string, error) {
	if !validDocumentTypes[documentType] {
		return "", ErrInvalidDocumentType
	}
	db, err := s.getDB(ctx)
	if err != nil {
		return "", err
	}

	var result string
	txErr := db.Transaction(func(tx *gorm.DB) error {
		var setting DocumentNumberingSetting
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("document_type = ?", documentType).
			First(&setting).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSettingNotFound
			}
			return fmt.Errorf("failed to load numbering setting: %w", err)
		}

		now := s.now()
		resetKey := ResetKeyFor(setting.ResetPeriod, now)
		if resetKey != setting.LastResetKey {
			setting.LastSequence = 0
			setting.LastResetKey = resetKey
		}
		setting.LastSequence++

		if err := tx.Save(&setting).Error; err != nil {
			return fmt.Errorf("failed to persist numbering sequence: %w", err)
		}

		result = FormatTemplate(setting.FormatTemplate, setting.LastSequence, now)
		return nil
	})
	if txErr != nil {
		return "", txErr
	}
	return result, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd backend && go test ./internal/pkg/numbering/... -v`
Expected: PASS — all tests in the package green (format + service).

- [ ] **Step 5: Commit**

```bash
git add backend/internal/pkg/numbering/service.go backend/internal/pkg/numbering/service_test.go
git commit -m "feat: add numbering.Service with List/Update/Preview/Generate"
```

---

### Task 4: Expose settings routes for document numbering

**Files:**
- Modify: `backend/internal/modules/setting/handler.go`
- Modify: `backend/internal/modules/setting/routes.go`
- Modify: `backend/internal/modules/setting/module.go`
- Test: `backend/internal/modules/setting/numbering_handler_test.go`

**Interfaces:**
- Consumes: `numbering.Service` (`List`, `Update`, `Preview` from Task 3), `numbering.ErrInvalidDocumentType`, `numbering.ErrInvalidResetPeriod`, `numbering.ErrSettingNotFound`.
- Produces: HTTP routes `GET /api/v1/tenant/settings/document-numbering`, `PUT /api/v1/tenant/settings/document-numbering/:document_type`, `GET /api/v1/tenant/settings/document-numbering/:document_type/preview`.

- [ ] **Step 1: Read the existing `Handler` struct and `RegisterRoutes` signature**

Open `backend/internal/modules/setting/handler.go` and `backend/internal/modules/setting/routes.go` to see the exact `Handler` struct fields and constructor (`NewHandler(...)`), and confirm the `settings := rg.Group("/settings")` grouping used for `company-holidays`. Match that exact style for the new sub-resource.

- [ ] **Step 2: Write the failing handler test**

`backend/internal/modules/setting/numbering_handler_test.go` (adjust imports/helpers to match whatever test-setup helpers this package's existing `*_test.go` files use, e.g. `service_test.go` — reuse them rather than inventing new ones):
```go
package setting

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hris/backend/internal/pkg/numbering"
)

func setupNumberingTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&numbering.DocumentNumberingSetting{}); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	seed := numbering.DocumentNumberingSetting{
		ID: "11111111-1111-1111-1111-111111111111", DocumentType: numbering.DocumentTypeEmployeeMovement,
		FormatTemplate: "SK/{sequence:3}/{year}", ResetPeriod: numbering.ResetPeriodYearly,
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed: %v", err)
	}
	resolver := func(ctx context.Context) (*gorm.DB, error) { return db, nil }
	numberingSvc := numbering.NewService(resolver, zap.NewNop())

	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewNumberingHandler(numberingSvc)
	group := r.Group("/api/v1/tenant/settings")
	RegisterNumberingRoutes(group, handler)
	return r
}

func TestListNumberingSettings(t *testing.T) {
	r := setupNumberingTestRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/settings/document-numbering", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var body struct {
		Data []numbering.DocumentNumberingSetting `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(body.Data) != 1 {
		t.Fatalf("expected 1 setting, got %d", len(body.Data))
	}
}

func TestUpdateNumberingSettingRejectsBadDocumentType(t *testing.T) {
	r := setupNumberingTestRouter(t)
	payload := `{"format_template":"X/{sequence}","reset_period":"yearly"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tenant/settings/document-numbering/not_a_type", stringsReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}
```

Add this small helper at the bottom of the same test file (avoids pulling in `strings` just for one call site inconsistently with the rest of the file — check whether `strings` is already imported elsewhere in the package's tests and just use `strings.NewReader` directly if so, dropping this helper):
```go
func stringsReader(s string) *bytes.Reader {
	return bytes.NewReader([]byte(s))
}
```
(adjust the import block to add `"bytes"` if you keep this helper).

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/setting/... -run TestListNumberingSettings -v`
Expected: FAIL — `NewNumberingHandler`/`RegisterNumberingRoutes` undefined.

- [ ] **Step 4: Add handler code**

Append to `backend/internal/modules/setting/handler.go` (or create `backend/internal/modules/setting/numbering_handler.go` if this package splits handlers into multiple files — check the existing file first and follow its convention):
```go
type NumberingHandler struct {
	svc *numbering.Service
}

func NewNumberingHandler(svc *numbering.Service) *NumberingHandler {
	return &NumberingHandler{svc: svc}
}

func (h *NumberingHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

type updateNumberingRequest struct {
	FormatTemplate string `json:"format_template" binding:"required"`
	ResetPeriod    string `json:"reset_period" binding:"required"`
}

func (h *NumberingHandler) Update(c *gin.Context) {
	documentType := c.Param("document_type")
	var req updateNumberingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.svc.Update(c.Request.Context(), documentType, req.FormatTemplate, req.ResetPeriod)
	if err != nil {
		switch err {
		case numbering.ErrInvalidDocumentType, numbering.ErrInvalidResetPeriod:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case numbering.ErrSettingNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

func (h *NumberingHandler) Preview(c *gin.Context) {
	documentType := c.Param("document_type")
	preview, err := h.svc.Preview(c.Request.Context(), documentType)
	if err != nil {
		switch err {
		case numbering.ErrInvalidDocumentType:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case numbering.ErrSettingNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"preview": preview}})
}
```
Add `"hris/backend/internal/pkg/numbering"` to the import block (use whatever module path prefix the rest of `backend/internal/modules/setting/handler.go` already uses for internal imports — copy it exactly rather than guessing `hris/backend/...`).

- [ ] **Step 5: Add route registration**

In `backend/internal/modules/setting/routes.go`, inside the existing `settings := rg.Group("/settings")` block (same place `company-holidays` is registered), add:
```go
// Document Numbering
docNumbering := settings.Group("/document-numbering")
{
	docNumbering.GET("", numberingHandler.List)
	docNumbering.PUT("/:document_type", numberingHandler.Update)
	docNumbering.GET("/:document_type/preview", numberingHandler.Preview)
}
```
Also add a standalone `RegisterNumberingRoutes(rg *gin.RouterGroup, handler *NumberingHandler)` function (used directly by the Task 4 test, and callable from the main `RegisterRoutes` too) so the routes are testable in isolation:
```go
func RegisterNumberingRoutes(rg *gin.RouterGroup, handler *NumberingHandler) {
	docNumbering := rg.Group("/document-numbering")
	{
		docNumbering.GET("", handler.List)
		docNumbering.PUT("/:document_type", handler.Update)
		docNumbering.GET("/:document_type/preview", handler.Preview)
	}
}
```
Then simplify the block added to `RegisterRoutes` to just call `RegisterNumberingRoutes(settings, numberingHandler)`.

- [ ] **Step 6: Wire `NumberingHandler` into the module's construction**

In `backend/internal/modules/setting/module.go`, find where `Handler` is constructed (`NewHandler(...)`) inside `NewModule`/`NewModuleWithService` and add a `numbering.Service` field: construct it as `numberingSvc := numbering.NewService(NewTenantDBResolver(dbManager), logger)` (reusing the module's existing `dbManager`/`logger` already in scope there) and pass `numberingHandler := NewNumberingHandler(numberingSvc)` through to wherever `RegisterRoutes` is invoked so it can call the new `RegisterNumberingRoutes`. Store `numberingSvc` on the module struct if `employeemovement` will need to reach it externally (see Task 5 — likely simplest to instead construct `numbering.Service` once in `main.go` and pass it to both modules; if so, skip constructing it again here and accept it as a constructor parameter of `NewModule`/`NewModuleWithService` instead. Pick whichever requires the smaller diff to `main.go` and note the choice in the commit message).

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/setting/... -v`
Expected: PASS — including `TestListNumberingSettings` and `TestUpdateNumberingSettingRejectsBadDocumentType`, and no regressions in existing `setting` package tests.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/modules/setting/
git commit -m "feat: add document numbering settings routes to setting module"
```

---

### Task 5: Auto-generate numbers in `employeemovement.Service`

**Files:**
- Modify: `backend/internal/modules/employeemovement/service.go`
- Modify: `backend/cmd/server/main.go`
- Test: `backend/internal/modules/employeemovement/numbering_test.go`

**Interfaces:**
- Consumes: `numbering.Service.Generate(ctx, documentType string) (string, error)`, `numbering.DocumentTypeEmployeeMovement`, `numbering.DocumentTypeEmployeeContract` (Task 3).
- Produces: `func (s *Service) SetNumberingService(ns NumberingGenerator)` where `NumberingGenerator` is a small interface `{ Generate(ctx context.Context, documentType string) (string, error) }` defined in `employeemovement` (so this module doesn't import `numbering` directly, mirroring how `ApprovalEngine`/`Notifier` are interfaces owned by this package) — `CreateMovement` and `CreateContract` use it when the incoming number field is blank.

- [ ] **Step 1: Write the failing tests**

`backend/internal/modules/employeemovement/numbering_test.go` — check the existing `service_test.go` in this package for how a `*Service` is constructed with a fake repo/logger in tests, and reuse that exact setup helper rather than duplicating it. Then add:
```go
package employeemovement

import (
	"context"
	"testing"
)

type fakeNumberingGenerator struct {
	calls  []string
	number string
}

func (f *fakeNumberingGenerator) Generate(ctx context.Context, documentType string) (string, error) {
	f.calls = append(f.calls, documentType)
	return f.number, nil
}

func TestCreateMovementGeneratesNumberWhenBlank(t *testing.T) {
	svc, repo := newTestServiceWithRepo(t) // reuse whatever helper service_test.go already provides
	gen := &fakeNumberingGenerator{number: "SK/001/HRIS/VIII/2026"}
	svc.SetNumberingService(gen)

	req := CreateMovementRequest{
		// fill required fields exactly as an existing passing test in
		// service_test.go does for CreateMovement, but leave DecisionLetterNumber blank
		DecisionLetterNumber: "",
	}
	resp, err := svc.CreateMovement(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMovement: %v", err)
	}
	if resp.DecisionLetterNumber != "SK/001/HRIS/VIII/2026" {
		t.Fatalf("expected generated number, got %q", resp.DecisionLetterNumber)
	}
	if len(gen.calls) != 1 || gen.calls[0] != "employee_movement" {
		t.Fatalf("expected one Generate call for employee_movement, got %v", gen.calls)
	}
	_ = repo
}

func TestCreateMovementKeepsManualNumber(t *testing.T) {
	svc, _ := newTestServiceWithRepo(t)
	gen := &fakeNumberingGenerator{number: "SK/999/HRIS/VIII/2026"}
	svc.SetNumberingService(gen)

	req := CreateMovementRequest{
		DecisionLetterNumber: "SK/MANUAL/2026",
	}
	resp, err := svc.CreateMovement(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateMovement: %v", err)
	}
	if resp.DecisionLetterNumber != "SK/MANUAL/2026" {
		t.Fatalf("expected manual number preserved, got %q", resp.DecisionLetterNumber)
	}
	if len(gen.calls) != 0 {
		t.Fatalf("expected Generate not called when number provided, got %v", gen.calls)
	}
}
```
Note: `newTestServiceWithRepo` is a placeholder name — open `backend/internal/modules/employeemovement/service_test.go` first and use whatever constructor helper (or inline `NewService(fakeRepo, zap.NewNop())`) it already uses, filling in whatever other required `CreateMovementRequest` fields its existing `TestCreateMovement*` tests use so the request passes validation.

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/employeemovement/... -run TestCreateMovementGeneratesNumberWhenBlank -v`
Expected: FAIL — `SetNumberingService`/`NumberingGenerator` undefined.

- [ ] **Step 3: Add the interface, field, setter, and call sites**

In `backend/internal/modules/employeemovement/service.go`, near the other provider interfaces (`ApprovalEngine`, `Notifier`, etc., around line 150), add:
```go
// NumberingGenerator generates the next document number for a document type
// (see backend/internal/pkg/numbering for the concrete implementation).
type NumberingGenerator interface {
	Generate(ctx context.Context, documentType string) (string, error)
}
```
In the `Service` struct (around line 165), add a field:
```go
	numberingService NumberingGenerator
```
Add a setter next to `SetApprovalEngine` (around line 185):
```go
// SetNumberingService wires the document numbering package so
// CreateMovement/CreateContract can auto-generate a number when the caller
// leaves it blank.
func (s *Service) SetNumberingService(ns NumberingGenerator) {
	s.numberingService = ns
}
```
In `CreateMovement` (service.go:332), right before the `EmployeeMovement{...}` struct is built (or right after, before `s.repo.CreateMovement` is called — wherever `req.DecisionLetterNumber` is currently read into the model), add:
```go
	decisionLetterNumber := req.DecisionLetterNumber
	if decisionLetterNumber == "" && s.numberingService != nil {
		generated, err := s.numberingService.Generate(ctx, "employee_movement")
		if err != nil {
			return nil, fmt.Errorf("failed to generate decision letter number: %w", err)
		}
		decisionLetterNumber = generated
	}
```
then use `decisionLetterNumber` (instead of `req.DecisionLetterNumber`) when constructing the `EmployeeMovement` model.

In `CreateContract` (service.go:1878), apply the same pattern for `req.ContractNumber`, calling `s.numberingService.Generate(ctx, "employee_contract")` when blank, using the constant string `"employee_contract"` — this package intentionally uses raw string literals here (not `numbering.DocumentTypeEmployeeContract`) to avoid importing `numbering` directly, per the `NumberingGenerator` interface design; add a short comment noting the two literals must stay in sync with `numbering.DocumentTypeEmployeeMovement`/`DocumentTypeEmployeeContract`.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/employeemovement/... -v`
Expected: PASS — new tests green, no regressions in existing `employeemovement` tests.

- [ ] **Step 5: Wire the concrete `numbering.Service` in `main.go`**

In `backend/cmd/server/main.go`, near line 931-940 (right after `employeeMovementSvc.SetNotifier(notificationSvc)`), add:
```go
employeeMovementSvc.SetNumberingService(numberingSvc)
```
where `numberingSvc` is the same `*numbering.Service` instance constructed for/by the `setting` module in Task 4 Step 6 — if Task 4 ended up constructing it inside `setting.NewModule`, change that to instead construct it once here in `main.go` (`numberingSvc := numbering.NewService(setting.NewTenantDBResolver(dbManager), l.Named("numbering"))`, reusing the exported `setting.NewTenantDBResolver` helper already used at main.go:912) and pass it into both `setting.NewModule(...)` (as a new parameter) and this `SetNumberingService` call, so there is exactly one `numbering.Service` per process. Add the `"hris/backend/internal/pkg/numbering"` import (matching this file's existing import path convention) to `main.go`.

- [ ] **Step 6: Build the whole backend to confirm wiring compiles**

Run: `cd backend && go build ./...`
Expected: builds with no errors.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/employeemovement/service.go backend/internal/modules/employeemovement/numbering_test.go backend/cmd/server/main.go backend/internal/modules/setting/module.go
git commit -m "feat: auto-generate movement/contract numbers via numbering service"
```

---

### Task 6: Frontend — Numbering Settings page

**Files:**
- Create: `frontend/tenant/src/views/settings/NumberingSettingsView.vue`
- Modify: `frontend/tenant/src/router/index.js`
- Modify: `frontend/tenant/src/layouts/Sidebar.vue`
- Modify: `frontend/tenant/src/locales/en.json`
- Modify: `frontend/tenant/src/locales/id.json`

**Interfaces:**
- Consumes: `GET /api/v1/tenant/settings/document-numbering`, `PUT /api/v1/tenant/settings/document-numbering/:document_type`, `GET /api/v1/tenant/settings/document-numbering/:document_type/preview` (Task 4).

- [ ] **Step 1: Read `CompanyHolidaysView.vue` fully to copy its exact structural conventions**

Open `frontend/tenant/src/views/settings/CompanyHolidaysView.vue` and note: how `api` is imported, the `useI18n()`/`useToast()` setup, how `FormRow`/`TextInput` are imported and used, how `getValidationErrors` is used to populate an `errors` ref, and the overall `<template>`/`<script setup>`/style structure. Match this file's conventions exactly in the new file (do not introduce a different state-management or API-calling style).

- [ ] **Step 2: Write `NumberingSettingsView.vue`**

`frontend/tenant/src/views/settings/NumberingSettingsView.vue`:
```vue
<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import FormRow from '@/components/form/FormRow.vue'
import TextInput from '@/components/form/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

const DOCUMENT_TYPES = [
  { key: 'employee_movement', labelKey: 'numbering_settings.employee_movement' },
  { key: 'employee_contract', labelKey: 'numbering_settings.employee_contract' },
]

const RESET_PERIOD_OPTIONS = [
  { value: 'yearly', labelKey: 'numbering_settings.reset_yearly' },
  { value: 'monthly', labelKey: 'numbering_settings.reset_monthly' },
  { value: 'never', labelKey: 'numbering_settings.reset_never' },
]

const loading = ref(true)
const forms = ref({}) // { [document_type]: { format_template, reset_period } }
const previews = ref({}) // { [document_type]: string }
const saving = ref({}) // { [document_type]: boolean }
const errors = ref({}) // { [document_type]: {...} }

async function loadSettings() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/document-numbering')
    const items = res.data?.data || []
    for (const item of items) {
      forms.value[item.document_type] = {
        format_template: item.format_template,
        reset_period: item.reset_period,
      }
    }
    await Promise.all(DOCUMENT_TYPES.map((dt) => refreshPreview(dt.key)))
  } finally {
    loading.value = false
  }
}

async function refreshPreview(documentType) {
  try {
    const res = await api.get(`/api/v1/tenant/settings/document-numbering/${documentType}/preview`)
    previews.value[documentType] = res.data?.data?.preview || ''
  } catch {
    previews.value[documentType] = ''
  }
}

async function save(documentType) {
  saving.value[documentType] = true
  errors.value[documentType] = {}
  try {
    const payload = forms.value[documentType]
    await api.put(`/api/v1/tenant/settings/document-numbering/${documentType}`, payload)
    toast.add({ severity: 'success', summary: t('common.saved'), life: 3000 })
    await refreshPreview(documentType)
  } catch (e) {
    errors.value[documentType] = getValidationErrors(e)
  } finally {
    saving.value[documentType] = false
  }
}

onMounted(loadSettings)
</script>

<template>
  <div class="numbering-settings">
    <h2>{{ t('numbering_settings.title') }}</h2>

    <div v-if="loading">{{ t('common.loading') }}</div>

    <div v-else v-for="dt in DOCUMENT_TYPES" :key="dt.key" class="numbering-settings__section">
      <h3>{{ t(dt.labelKey) }}</h3>

      <FormRow :label="t('numbering_settings.format_template')" :errors="errors[dt.key]?.format_template">
        <TextInput v-model="forms[dt.key].format_template" @update:modelValue="refreshPreview(dt.key)" />
      </FormRow>

      <FormRow :label="t('numbering_settings.reset_period')" :errors="errors[dt.key]?.reset_period">
        <select v-model="forms[dt.key].reset_period" @change="refreshPreview(dt.key)">
          <option v-for="opt in RESET_PERIOD_OPTIONS" :key="opt.value" :value="opt.value">
            {{ t(opt.labelKey) }}
          </option>
        </select>
      </FormRow>

      <p class="numbering-settings__help">{{ t('numbering_settings.tokens_help') }}</p>
      <p class="numbering-settings__preview">
        {{ t('numbering_settings.preview_label') }}: <strong>{{ previews[dt.key] }}</strong>
      </p>

      <button :disabled="saving[dt.key]" @click="save(dt.key)">
        {{ t('common.save') }}
      </button>
    </div>
  </div>
</template>
```
Adjust the `select`/`button` markup to use whichever PrimeVue components (`Select`, `Button`) `CompanyHolidaysView.vue` uses instead of raw HTML, once Step 1's read confirms the exact component names/imports.

- [ ] **Step 3: Add i18n strings**

In `frontend/tenant/src/locales/en.json`, add a `numbering_settings` block:
```json
"numbering_settings": {
  "title": "Document Numbering",
  "employee_movement": "Movements (Promotion / Mutation / Demotion SK)",
  "employee_contract": "Contracts",
  "format_template": "Number Format",
  "reset_period": "Reset Sequence",
  "reset_yearly": "Every Year",
  "reset_monthly": "Every Month",
  "reset_never": "Never",
  "tokens_help": "Available tokens: {sequence}, {sequence:3} (zero-padded), {year}, {yy}, {month}, {month_roman}",
  "preview_label": "Next number preview"
}
```
In `frontend/tenant/src/locales/id.json`, add the matching Indonesian block:
```json
"numbering_settings": {
  "title": "Penomoran Dokumen",
  "employee_movement": "Movements (SK Promosi / Mutasi / Demosi)",
  "employee_contract": "Kontrak",
  "format_template": "Format Nomor",
  "reset_period": "Reset Urutan",
  "reset_yearly": "Setiap Tahun",
  "reset_monthly": "Setiap Bulan",
  "reset_never": "Tidak Pernah",
  "tokens_help": "Token yang tersedia: {sequence}, {sequence:3} (diberi angka nol di depan), {year}, {yy}, {month}, {month_roman}",
  "preview_label": "Pratinjau nomor berikutnya"
}
```

- [ ] **Step 4: Register the route**

In `frontend/tenant/src/router/index.js`, find the settings children array (around lines 560-588, next to the `company-holidays` route) and add:
```js
{
  path: 'settings/numbering',
  name: 'SettingsNumbering',
  component: () => import('@/views/settings/NumberingSettingsView.vue'),
  meta: {
    title: 'Document Numbering',
    titleKey: 'numbering_settings.title',
    descKey: 'numbering_settings.title',
    icon: 'pi pi-hashtag',
    module: 'setting',
  },
},
```
(match indentation/field set exactly to the neighboring `company-holidays` route entry.)

- [ ] **Step 5: Add the sidebar entry**

In `frontend/tenant/src/layouts/Sidebar.vue`, find where the Settings section lists its children (next to the `company-holidays` entry) and add a matching entry pointing to `SettingsNumbering`, copying the exact structure of the neighboring entry (icon, label i18n key, route name).

- [ ] **Step 6: Manually verify in the browser**

Start the frontend dev server, log in as a tenant user with settings access, navigate to Settings → Document Numbering. Confirm: both sections load with the seeded default templates, editing the template updates the preview after save, and Save persists (reload the page and confirm values stuck).

- [ ] **Step 7: Commit**

```bash
git add frontend/tenant/src/views/settings/NumberingSettingsView.vue frontend/tenant/src/router/index.js frontend/tenant/src/layouts/Sidebar.vue frontend/tenant/src/locales/en.json frontend/tenant/src/locales/id.json
git commit -m "feat: add document numbering settings page"
```

---

### Task 7: Frontend — auto-fill movement & contract number fields

**Files:**
- Modify: `frontend/tenant/src/views/modules/employeemovement/EmployeeMovements.vue`
- Modify: `frontend/tenant/src/views/modules/employeemovement/EmployeeContracts.vue`
- Modify: `backend/internal/modules/employeemovement/dto.go`

**Interfaces:**
- Consumes: Task 5's backend behavior (blank number field → auto-generated in the create response).

- [ ] **Step 1: Make `DecisionLetterNumber` optional in `CreateMovementRequest`**

In `backend/internal/modules/employeemovement/dto.go`, change the `DecisionLetterNumber` field's binding tag from `binding:"required"` to `binding:"omitempty"` so a blank value is accepted by validation (the numbering service fills it in during `CreateMovement`). Leave `ContractNumber` on `CreateContractRequest` as-is — check its current binding tag; if it's `binding:"required"`, change it the same way so blank submissions pass validation.

- [ ] **Step 2: Rebuild and run the existing DTO/service tests to confirm no regressions**

Run: `cd backend && go build ./... && go test ./internal/modules/employeemovement/... -v`
Expected: PASS, including the Task 5 tests (which rely on blank values reaching the service layer).

- [ ] **Step 3: Update `EmployeeMovements.vue` create form**

Open `frontend/tenant/src/views/modules/employeemovement/EmployeeMovements.vue`, find the `FormRow` for `decision_letter_number` (client-side required validation and `TextInput` binding). Remove the client-side "required" check for `decision_letter_number` in the create flow (find the validation function that currently does `if (!form.value.decision_letter_number?.trim()) e.decision_letter_number = ...` and either drop that check entirely for create mode, or gate it so it only applies on edit — follow whichever existing conditional pattern the file already uses to distinguish create vs edit, e.g. an `editing` ref). Add a `:placeholder` to the `TextInput` for this field: `:placeholder="t('employee_movement.number_auto_placeholder')"`.

- [ ] **Step 4: Update `EmployeeContracts.vue` create form**

Same change in `frontend/tenant/src/views/modules/employeemovement/EmployeeContracts.vue` for `contract_number`: drop/gate the required client-side check for create mode, add the placeholder.

- [ ] **Step 5: Add the placeholder i18n key**

Add to both `frontend/tenant/src/locales/en.json` and `id.json` under the existing `employee_movement` block:
- en: `"number_auto_placeholder": "Auto-generated if left blank"`
- id: `"number_auto_placeholder": "Akan digenerate otomatis jika dikosongkan"`

- [ ] **Step 6: Manually verify in the browser**

Create a new employee movement leaving the decision letter number blank → confirm it saves and the created record shows a generated number matching the configured template. Repeat for a new contract leaving contract number blank. Then create one more of each while typing a manual number → confirm the manual value is kept, not overwritten.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/employeemovement/dto.go frontend/tenant/src/views/modules/employeemovement/EmployeeMovements.vue frontend/tenant/src/views/modules/employeemovement/EmployeeContracts.vue frontend/tenant/src/locales/en.json frontend/tenant/src/locales/id.json
git commit -m "feat: auto-fill movement and contract numbers in create forms"
```

---

## Post-Implementation Checklist

- [ ] All backend tests pass: `cd backend && go test ./...`
- [ ] Backend builds cleanly: `cd backend && go build ./...`
- [ ] Both mysql and postgres tenant migrations applied and verified on a real/test tenant DB
- [ ] Manual browser walkthrough from Task 6 Step 6 and Task 7 Step 6 completed
- [ ] Spec (`docs/superpowers/specs/2026-08-14-document-numbering-settings-design.md`) requirements all covered: format tokens ✓, per-document-type reset period ✓, single counter per document type (not per movement subtype) ✓, auto-generate on create ✓, manual override allowed ✓, settings CRUD + preview ✓
