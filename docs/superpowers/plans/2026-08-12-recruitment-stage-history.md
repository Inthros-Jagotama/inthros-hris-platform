# Recruitment Pipeline Stage History (G-5) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an auditable status-transition history for `job_applications` (recruitment pipeline), enforced by a single validated state machine, plus a `GET` endpoint to read it.

**Architecture:** New master table `recruitment_stages` (seeded from the existing 8 `CandidateStatus` values) and a `job_application_stage_histories` table recording every transition. A single private service method `transitionApplicationStatus` validates the transition and writes history; it's called by both `UpdateApplicationStatus` (manual, via HTTP endpoint) and `AcceptOffer` (automatic, when an offer is accepted) so no status change can bypass validation/history.

**Tech Stack:** Go, GORM, Gin, MySQL + PostgreSQL (dual migration files), existing `backend/internal/modules/recruitment` module.

## Global Constraints

- Migrations must be idempotent and written for **both** `mysql` and `postgres` under `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/`, plus a matching `.down.sql` — follow the exact style of `095_recruitment_offer.sql` (CREATE TABLE) since this task creates two new tables, not ALTER COLUMN.
- `CandidateStatus` taxonomy is **not** changed (8 existing values only: `NEW, SCREENED, SHORTLISTED, INTERVIEWED, OFFERED, ACCEPTED, REJECTED, WITHDRAWN`). Do not introduce new status strings.
- State machine (see spec `docs/superpowers/specs/2026-08-12-recruitment-stage-history-design.md` for full rationale):
  - From any non-terminal status → `ACCEPTED`, `REJECTED`, `WITHDRAWN`: always allowed.
  - Between non-terminal statuses: allowed only if the destination's order is **≥** the source's order (order: `NEW=1, SCREENED=2, SHORTLISTED=3, INTERVIEWED=4, OFFERED=5`) — forward jumps allowed (e.g. `NEW→OFFERED`), backward moves rejected (e.g. `SHORTLISTED→NEW`).
  - From a terminal status (`ACCEPTED`, `REJECTED`, `WITHDRAWN`) to anything else: rejected.
  - `from == to`: no-op (no error, no new history row) — required for `AcceptOffer` idempotency on a second offer.
- Existing tests `TestService_UpdateApplicationStatus` (`service_test.go:638`, `NEW→SHORTLISTED`) and `TestHandler_UpdateApplicationStatus` (`handler_test.go:342`, same) must keep passing unmodified — they exercise a forward jump, which the state machine above allows.
- Invalid transitions return a sentinel error `ErrInvalidStatusTransition` (wrapped with `%w` and the transition detail) so the handler can map it to `400` via `errors.Is` — follow this repo's existing `httputil.BadRequest` / `httputil.InternalError` split (see `handler.go:317` for the one existing `BadRequest` usage).
- New migration number is **097** (last existing is `096_recruitment_employee_handoff`).

---

## Task 1: Migration 097 — `recruitment_stages` + `job_application_stage_histories`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.down.sql`

**Interfaces:**
- Produces: tables `recruitment_stages (id, code, name, sort_order, created_at, updated_at)` and `job_application_stage_histories (id, application_id, from_stage_id, to_stage_id, changed_by, notes, changed_at, created_at)` that Task 2's GORM models must match exactly (column names/types).

- [ ] **Step 1: Write postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 097_recruitment_stage_history (PostgreSQL)
-- =============================================================================
-- G-5 Pipeline Stage History
-- (docs/module-recruitment-development-plan.md §G-5)
--
-- recruitment_stages (master, seeded dari 8 CandidateStatus existing):
--   id          CHAR(36) PK
--   code        VARCHAR(20) NN UNIQUE  → NEW | SCREENED | SHORTLISTED |
--                                        INTERVIEWED | OFFERED | ACCEPTED |
--                                        REJECTED | WITHDRAWN
--   name        VARCHAR(100) NN        → label tampilan
--   sort_order  INT NN DEFAULT 0
--
-- job_application_stage_histories (audit trail transisi status aplikasi):
--   id              CHAR(36) PK
--   application_id  CHAR(36) NN  → job_applications
--   from_stage_id   CHAR(36) NULL → recruitment_stages (NULL untuk histori
--                                  pertama saat aplikasi dibuat berstatus NEW)
--   to_stage_id     CHAR(36) NN  → recruitment_stages
--   changed_by      CHAR(36) NULL → user id aktor (NULL bila sistem)
--   notes           TEXT NULL
--   changed_at      BIGINT NN    → unix nano
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS recruitment_stages (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_recruitment_stages_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS job_application_stage_histories (
    id              CHAR(36) PRIMARY KEY,
    application_id  CHAR(36) NOT NULL,
    from_stage_id   CHAR(36) NULL,
    to_stage_id     CHAR(36) NOT NULL,
    changed_by      CHAR(36) NULL,
    notes           TEXT NULL,
    changed_at      BIGINT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ash_application FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT fk_ash_from_stage FOREIGN KEY (from_stage_id) REFERENCES recruitment_stages(id),
    CONSTRAINT fk_ash_to_stage FOREIGN KEY (to_stage_id) REFERENCES recruitment_stages(id)
);

CREATE INDEX IF NOT EXISTS idx_ash_app ON job_application_stage_histories (application_id);
CREATE INDEX IF NOT EXISTS idx_ash_changed_at ON job_application_stage_histories (changed_at);
```

- [ ] **Step 2: Write postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 097_recruitment_stage_history (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS job_application_stage_histories;
DROP TABLE IF EXISTS recruitment_stages;
```

- [ ] **Step 3: Write mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 097_recruitment_stage_history (MySQL)
-- =============================================================================
-- G-5 Pipeline Stage History — lihat versi postgres untuk penjelasan kolom.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS recruitment_stages (
    id          CHAR(36) PRIMARY KEY,
    code        VARCHAR(20) NOT NULL,
    name        VARCHAR(100) NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_recruitment_stages_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS job_application_stage_histories (
    id              CHAR(36) PRIMARY KEY,
    application_id  CHAR(36) NOT NULL,
    from_stage_id   CHAR(36) NULL,
    to_stage_id     CHAR(36) NOT NULL,
    changed_by      CHAR(36) NULL,
    notes           TEXT NULL,
    changed_at      BIGINT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ash_application FOREIGN KEY (application_id) REFERENCES job_applications(id) ON DELETE CASCADE,
    CONSTRAINT fk_ash_from_stage FOREIGN KEY (from_stage_id) REFERENCES recruitment_stages(id),
    CONSTRAINT fk_ash_to_stage FOREIGN KEY (to_stage_id) REFERENCES recruitment_stages(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_ash_app ON job_application_stage_histories (application_id);
CREATE INDEX idx_ash_changed_at ON job_application_stage_histories (changed_at);
```

- [ ] **Step 4: Write mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 097_recruitment_stage_history (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS job_application_stage_histories;
DROP TABLE IF EXISTS recruitment_stages;
```

- [ ] **Step 5: Verify migration numbering doesn't collide**

Run: `ls backend/internal/pkg/migrator/migrations/tenant/postgres/ | grep "^097"` and same for mysql — expect exactly the 2 files each just created (4 total including down). Also confirm no other `097_*` already exists before writing (re-check `096` is still the highest prior to this task).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/097_recruitment_stage_history.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/097_recruitment_stage_history.down.sql
git commit -m "feat: migration 097 recruitment_stages + job_application_stage_histories (G-5)"
```

---

## Task 2: GORM Models — `RecruitmentStage`, `ApplicationStageHistory`

**Files:**
- Modify: `backend/internal/modules/recruitment/model.go`

**Interfaces:**
- Consumes: nothing new.
- Produces: `RecruitmentStage` struct (fields `ID, Code, Name, SortOrder, CreatedAt, UpdatedAt`), `ApplicationStageHistory` struct (fields `ID, ApplicationID, FromStageID *uuid.UUID, ToStageID uuid.UUID, ChangedBy *uuid.UUID, Notes string, ChangedAt int64, CreatedAt time.Time`) — Task 3 (repository) and Task 4 (service) depend on these exact field names.

- [ ] **Step 1: Add the two structs to `model.go`**

Append after the `JobApplication` block (after line ~273, before the `Interview` section), matching the existing style (`BeforeCreate` UUID generator, `TableName()`):

```go
// =========================================================================
// RecruitmentStage (G-5 — master stage, seeded dari CandidateStatus)
// =========================================================================
// Seeded 1:1 dari 8 CandidateStatus existing (bukan taxonomy baru) — lihat
// docs/superpowers/specs/2026-08-12-recruitment-stage-history-design.md.

type RecruitmentStage struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(20);not null;uniqueIndex:uq_recruitment_stages_code" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	SortOrder int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RecruitmentStage) TableName() string {
	return "recruitment_stages"
}

func (s *RecruitmentStage) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// ApplicationStageHistory (G-5 — audit trail transisi status aplikasi)
// =========================================================================

type ApplicationStageHistory struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	ApplicationID uuid.UUID  `gorm:"type:char(36);not null;index:idx_ash_app" json:"application_id"`
	FromStageID   *uuid.UUID `gorm:"type:char(36)" json:"from_stage_id,omitempty"`
	ToStageID     uuid.UUID  `gorm:"type:char(36);not null" json:"to_stage_id"`
	ChangedBy     *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	Notes         string     `gorm:"type:text" json:"notes"`
	ChangedAt     int64      `gorm:"type:bigint;not null;index:idx_ash_changed_at" json:"changed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (ApplicationStageHistory) TableName() string {
	return "job_application_stage_histories"
}

func (h *ApplicationStageHistory) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}
```

- [ ] **Step 2: Compile check**

Run: `cd backend && go build ./internal/modules/recruitment/...`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/recruitment/model.go
git commit -m "feat: model RecruitmentStage + ApplicationStageHistory (G-5)"
```

---

## Task 3: Repository — stage lookup + history CRUD

**Files:**
- Modify: `backend/internal/modules/recruitment/repository.go`
- Test: `backend/internal/modules/recruitment/repository_test.go`

**Interfaces:**
- Consumes: `RecruitmentStage`, `ApplicationStageHistory` (Task 2).
- Produces (used by Task 4 service):
  - `func (r *Repository) FindStageByCode(ctx context.Context, code string) (*RecruitmentStage, error)`
  - `func (r *Repository) ListStages(ctx context.Context) ([]RecruitmentStage, error)`
  - `func (r *Repository) CreateStage(ctx context.Context, s *RecruitmentStage) error`
  - `func (r *Repository) CreateStageHistory(ctx context.Context, h *ApplicationStageHistory) error`
  - `func (r *Repository) ListStageHistoryByApplication(ctx context.Context, applicationID uuid.UUID) ([]ApplicationStageHistory, error)`

- [ ] **Step 1: Write failing repository tests**

Append to `backend/internal/modules/recruitment/repository_test.go` (mirror the style of existing `TestRepository_CreateOffer`/`TestRepository_FindOfferByID` — check those first for exact `newTestRepo`/setup helper name used in that file):

```go
func TestRepository_CreateAndFindStageByCode(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	stage := &RecruitmentStage{Code: "NEW", Name: "New Application", SortOrder: 1}
	if err := repo.CreateStage(ctx, stage); err != nil {
		t.Fatalf("CreateStage failed: %v", err)
	}

	found, err := repo.FindStageByCode(ctx, "NEW")
	if err != nil {
		t.Fatalf("FindStageByCode failed: %v", err)
	}
	if found.ID != stage.ID {
		t.Errorf("expected id %s, got %s", stage.ID, found.ID)
	}
}

func TestRepository_ListStages(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	repo.CreateStage(ctx, &RecruitmentStage{Code: "NEW", Name: "New", SortOrder: 1})
	repo.CreateStage(ctx, &RecruitmentStage{Code: "SCREENED", Name: "Screened", SortOrder: 2})

	list, err := repo.ListStages(ctx)
	if err != nil {
		t.Fatalf("ListStages failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 stages, got %d", len(list))
	}
}

func TestRepository_CreateAndListStageHistory(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	req := &JobRequisition{OrganizationID: uuid.New(), Title: "Engineer"}
	repo.CreateRequisition(ctx, req)
	cand := &Candidate{FirstName: "A", LastName: "B", Email: "ab@test.com"}
	repo.CreateCandidate(ctx, cand)
	app := &JobApplication{RequisitionID: req.ID, CandidateID: cand.ID, Status: CandStatusNew}
	repo.CreateApplication(ctx, app)

	newStage := &RecruitmentStage{Code: "NEW", Name: "New", SortOrder: 1}
	repo.CreateStage(ctx, newStage)
	screenedStage := &RecruitmentStage{Code: "SCREENED", Name: "Screened", SortOrder: 2}
	repo.CreateStage(ctx, screenedStage)

	h := &ApplicationStageHistory{
		ApplicationID: app.ID,
		FromStageID:   &newStage.ID,
		ToStageID:     screenedStage.ID,
		ChangedAt:     1000,
	}
	if err := repo.CreateStageHistory(ctx, h); err != nil {
		t.Fatalf("CreateStageHistory failed: %v", err)
	}

	list, err := repo.ListStageHistoryByApplication(ctx, app.ID)
	if err != nil {
		t.Fatalf("ListStageHistoryByApplication failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 history row, got %d", len(list))
	}
	if list[0].ToStageID != screenedStage.ID {
		t.Errorf("expected to_stage_id %s, got %s", screenedStage.ID, list[0].ToStageID)
	}
}
```

Add these tests to `repository_test.go` after the existing Job Offers repository test section, matching the file's existing plain (no `*testing.T` passed into helper) style.

- [ ] **Step 2: Run tests to verify they fail (missing methods)**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestRepository_CreateAndFindStageByCode -v`
Expected: compile error, `CreateStage`/`FindStageByCode` undefined.

- [ ] **Step 3: Implement repository methods**

Append to `backend/internal/modules/recruitment/repository.go` (after the Job Offers section, before Candidates — or at the end of the file, matching existing section-comment style `// === Job Offers (G-3) ===`):

```go
// =========================================================================
// Recruitment Stages (G-5 — master, seeded)
// =========================================================================

func (r *Repository) CreateStage(ctx context.Context, s *RecruitmentStage) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindStageByCode(ctx context.Context, code string) (*RecruitmentStage, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s RecruitmentStage
	if err := db.WithContext(ctx).Where("code = ?", code).First(&s).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("recruitment stage not found: %s", code)
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListStages(ctx context.Context) ([]RecruitmentStage, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []RecruitmentStage
	if err := db.WithContext(ctx).Order("sort_order ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// =========================================================================
// Application Stage History (G-5)
// =========================================================================

func (r *Repository) CreateStageHistory(ctx context.Context, h *ApplicationStageHistory) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(h).Error
}

func (r *Repository) ListStageHistoryByApplication(ctx context.Context, applicationID uuid.UUID) ([]ApplicationStageHistory, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []ApplicationStageHistory
	if err := db.WithContext(ctx).Where("application_id = ?", applicationID).Order("changed_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestRepository_CreateAndFindStageByCode|TestRepository_ListStages|TestRepository_CreateAndListStageHistory" -v`
Expected: PASS (all 3).

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/recruitment/repository.go backend/internal/modules/recruitment/repository_test.go
git commit -m "feat: repository CRUD RecruitmentStage + ApplicationStageHistory (G-5)"
```

---

## Task 4: Service — state machine, `transitionApplicationStatus`, wire into existing flows

**Files:**
- Modify: `backend/internal/modules/recruitment/service.go`
- Test: `backend/internal/modules/recruitment/service_test.go`

**Interfaces:**
- Consumes: `Repository.FindStageByCode`, `CreateStageHistory`, `ListStageHistoryByApplication` (Task 3).
- Produces:
  - `var ErrInvalidStatusTransition = errors.New("invalid status transition")` (package-level, in `service.go`).
  - `func (s *Service) transitionApplicationStatus(ctx context.Context, a *JobApplication, newStatus CandidateStatus, changedBy *uuid.UUID, notes string) error` — validates + writes history; does **not** call `repo.UpdateApplication` (caller's responsibility, matching the existing `UpdateApplicationStatus`/`AcceptOffer` pattern of setting fields then saving once).
  - `func (s *Service) GetApplicationHistory(ctx context.Context, applicationID string) ([]StageHistoryResponse, error)` — used by Task 5 handler.

- [ ] **Step 1: Write failing service tests for the state machine**

Append to `backend/internal/modules/recruitment/service_test.go` (near `TestService_UpdateApplicationStatus` at line 623):

```go
func TestService_UpdateApplicationStatus_ForwardJumpAllowed(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "fwd@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	// NEW -> OFFERED direct jump must remain allowed (state machine allows
	// forward jumps between non-terminal stages).
	updated, err := svc.UpdateApplicationStatus(ctx, app.ID, "OFFERED", "", "")
	if err != nil {
		t.Fatalf("expected forward jump to succeed, got error: %v", err)
	}
	if updated.Status != "OFFERED" {
		t.Errorf("expected OFFERED, got %s", updated.Status)
	}
}

func TestService_UpdateApplicationStatus_BackwardRejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "back@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "SHORTLISTED", "", ""); err != nil {
		t.Fatalf("setup transition failed: %v", err)
	}

	_, err := svc.UpdateApplicationStatus(ctx, app.ID, "NEW", "", "")
	if err == nil {
		t.Fatal("expected error for backward transition SHORTLISTED -> NEW, got nil")
	}
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition, got: %v", err)
	}
}

func TestService_UpdateApplicationStatus_FromTerminalRejected(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "term@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "REJECTED", "", ""); err != nil {
		t.Fatalf("setup transition failed: %v", err)
	}

	_, err := svc.UpdateApplicationStatus(ctx, app.ID, "SCREENED", "", "")
	if !errors.Is(err, ErrInvalidStatusTransition) {
		t.Errorf("expected ErrInvalidStatusTransition from terminal status, got: %v", err)
	}
}

func TestService_UpdateApplicationStatus_SameStatusNoop(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "noop@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "NEW", "", ""); err != nil {
		t.Fatalf("same-status transition should be a no-op, got error: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	// Only the initial NEW history row from CreateApplication — the no-op
	// NEW->NEW call must not add a second row.
	if len(hist) != 1 {
		t.Errorf("expected 1 history row (initial only), got %d", len(hist))
	}
}

func TestService_CreateApplication_WritesInitialHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "init@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	if len(hist) != 1 {
		t.Fatalf("expected 1 initial history row, got %d", len(hist))
	}
	if hist[0].FromStage != nil {
		t.Errorf("expected initial history from_stage nil, got %v", hist[0].FromStage)
	}
	if hist[0].ToStage.Code != "NEW" {
		t.Errorf("expected initial history to_stage NEW, got %s", hist[0].ToStage.Code)
	}
}

func TestService_UpdateApplicationStatus_WritesHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	req, _ := svc.CreateRequisition(ctx, CreateRequisitionRequest{OrganizationID: createTestOrgID(), Title: "Engineer"})
	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "A", LastName: "B", Email: "wh@test.com"})
	app, _ := svc.CreateApplication(ctx, CreateApplicationRequest{RequisitionID: req.ID, CandidateID: cand.ID})

	if _, err := svc.UpdateApplicationStatus(ctx, app.ID, "SCREENED", "", "moved to screening"); err != nil {
		t.Fatalf("UpdateApplicationStatus failed: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	if len(hist) != 2 {
		t.Fatalf("expected 2 history rows (initial + transition), got %d", len(hist))
	}
	last := hist[len(hist)-1]
	if last.FromStage == nil || last.FromStage.Code != "NEW" || last.ToStage.Code != "SCREENED" {
		t.Errorf("expected NEW->SCREENED, got from=%v to=%s", last.FromStage, last.ToStage.Code)
	}
	if last.Notes != "moved to screening" {
		t.Errorf("expected notes preserved, got %q", last.Notes)
	}
}
```

Add `AcceptOffer` history coverage next to the existing offer tests (search
`TestService_AcceptOffer` in `service_test.go` for the exact existing test
name/location and add alongside it):

```go
func TestService_AcceptOffer_WritesHistory(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	svc.SetApprovalEngine(&fakeApprovalEngine{instanceID: uuid.New().String(), flowID: uuid.New().String()})

	// seedDraftOffer + approveAndSendOffer are existing test helpers in this
	// file (used by TestService_AcceptOffer_NoDoubleIncrementSlotsFilled at
	// service_test.go:1675) that create a requisition+candidate+application+
	// draft offer and drive it to SENT via the approval flow — reuse them
	// instead of re-deriving the submit/approve/send sequence.
	offer := seedDraftOffer(t, svc, ctx)
	app, err := svc.GetApplicationByID(ctx, offer.ApplicationID)
	if err != nil {
		t.Fatalf("GetApplicationByID failed: %v", err)
	}
	offerID := approveAndSendOffer(t, svc, ctx, offer.ID)

	if _, err := svc.AcceptOffer(ctx, offerID); err != nil {
		t.Fatalf("AcceptOffer failed: %v", err)
	}

	hist, err := svc.GetApplicationHistory(ctx, app.ID)
	if err != nil {
		t.Fatalf("GetApplicationHistory failed: %v", err)
	}
	last := hist[len(hist)-1]
	if last.ToStage.Code != "ACCEPTED" {
		t.Errorf("expected last history to_stage ACCEPTED, got %s", last.ToStage.Code)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestService_UpdateApplicationStatus_ForwardJumpAllowed -v`
Expected: compile error (`ErrInvalidStatusTransition`, `GetApplicationHistory`, `StageHistoryResponse` undefined) or runtime failure.

- [ ] **Step 3: Implement the state machine + `transitionApplicationStatus`**

Add near the top of `service.go` (package-level, after imports):

```go
var ErrInvalidStatusTransition = errors.New("invalid status transition")

// applicationStageOrder define urutan progresi status non-terminal (G-5).
// Status terminal (ACCEPTED/REJECTED/WITHDRAWN) sengaja tidak masuk sini —
// diperlakukan khusus di isValidStatusTransition.
var applicationStageOrder = map[CandidateStatus]int{
	CandStatusNew:         1,
	CandStatusScreened:    2,
	CandStatusShortlisted: 3,
	CandStatusInterviewed: 4,
	CandStatusOffered:     5,
}

func isTerminalStatus(s CandidateStatus) bool {
	return s == CandStatusAccepted || s == CandStatusRejected || s == CandStatusWithdrawn
}

// isValidStatusTransition menegakkan state machine G-5:
//   - from == to                         → true (no-op, ditangani caller)
//   - from terminal                      → false (tidak ada transisi keluar)
//   - to ∈ {ACCEPTED, REJECTED, WITHDRAWN} dan from non-terminal → true
//   - from & to sama-sama non-terminal    → true hanya jika order[to] >= order[from]
func isValidStatusTransition(from, to CandidateStatus) bool {
	if from == to {
		return true
	}
	if isTerminalStatus(from) {
		return false
	}
	if isTerminalStatus(to) {
		return true
	}
	fromOrder, fromOK := applicationStageOrder[from]
	toOrder, toOK := applicationStageOrder[to]
	if !fromOK || !toOK {
		return false
	}
	return toOrder >= fromOrder
}
```

Add the helper method (place near `UpdateApplicationStatus`, e.g. just
above it at `service.go:1126`):

```go
// transitionApplicationStatus memvalidasi transisi (state machine G-5),
// menulis baris job_application_stage_histories, dan meng-update field
// status + timestamp stage pada a (in-memory) — caller bertanggung jawab
// memanggil repo.UpdateApplication untuk menyimpannya. Dipakai oleh
// UpdateApplicationStatus (manual) dan AcceptOffer (otomatis, G-3/G-4) agar
// tidak ada perubahan status yang lolos tanpa histori.
func (s *Service) transitionApplicationStatus(ctx context.Context, a *JobApplication, newStatus CandidateStatus, changedBy *uuid.UUID, notes string) error {
	from := a.Status
	if !isValidStatusTransition(from, newStatus) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, from, newStatus)
	}
	if from == newStatus {
		return nil // no-op — idempotent, tidak menulis history baru
	}

	fromStage, err := s.repo.FindStageByCode(ctx, string(from))
	if err != nil {
		return fmt.Errorf("recruitment stage lookup failed for %s: %w", from, err)
	}
	toStage, err := s.repo.FindStageByCode(ctx, string(newStatus))
	if err != nil {
		return fmt.Errorf("recruitment stage lookup failed for %s: %w", newStatus, err)
	}

	now := time.Now().UnixNano()
	hist := &ApplicationStageHistory{
		ApplicationID: a.ID,
		FromStageID:   &fromStage.ID,
		ToStageID:     toStage.ID,
		ChangedBy:     changedBy,
		Notes:         notes,
		ChangedAt:     now,
	}
	if err := s.repo.CreateStageHistory(ctx, hist); err != nil {
		return fmt.Errorf("failed to write stage history: %w", err)
	}

	a.Status = newStatus
	switch newStatus {
	case CandStatusScreened:
		a.ScreenedAt = &now
	case CandStatusShortlisted:
		a.ShortlistedAt = &now
	case CandStatusOffered:
		a.OfferedAt = &now
	case CandStatusAccepted:
		a.AcceptedAt = &now
	case CandStatusRejected:
		a.RejectedAt = &now
	case CandStatusWithdrawn:
		a.WithdrawnAt = &now
	}
	return nil
}
```

Add `"errors"` to the `service.go` import block if not already present.

- [ ] **Step 4: Rewire `UpdateApplicationStatus` to use the helper**

Replace the body of `UpdateApplicationStatus` (`service.go:1126-1176`) — keep
the signature identical (handler in Task 5 stays unchanged for this
method):

```go
func (s *Service) UpdateApplicationStatus(ctx context.Context, id, status, reason, notes string) (*ApplicationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	a, err := s.repo.FindApplicationByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	if err := s.transitionApplicationStatus(ctx, a, CandidateStatus(status), nil, notes); err != nil {
		return nil, err
	}
	if reason != "" {
		a.RejectionReason = reason
	}
	if notes != "" {
		a.Notes = notes
	}

	// ACCEPTED: pertahankan efek samping slots_filled existing (di luar
	// tanggung jawab transitionApplicationStatus, yang hanya urus status +
	// history).
	if CandidateStatus(status) == CandStatusAccepted {
		req, findErr := s.repo.FindRequisitionByID(ctx, a.RequisitionID)
		if findErr == nil && req != nil {
			req.SlotsFilled++
			if req.SlotsFilled >= req.SlotsAvailable {
				req.Status = ReqStatusFilled
			}
			if err := s.repo.UpdateRequisition(ctx, req); err != nil {
				s.logger.Warn("failed to update requisition slots_filled", zap.String("requisition_id", req.ID.String()), zap.Error(err))
			}
		}
	}

	if err := s.repo.UpdateApplication(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Application status updated", zap.String("id", a.ID.String()), zap.String("status", string(a.Status)))
	return applicationToResponse(a), nil
}
```

> ⚠️ Note: this changes existing behavior slightly — manually calling
> `PUT /applications/:id/status` with `ACCEPTED` will now increment
> `slots_filled` **every time it's called with the same already-ACCEPTED
> status** is NOT the case (transition is a no-op when `from == to`, so
> repeated `ACCEPTED` calls after the first do nothing — the
> `if CandidateStatus(status) == CandStatusAccepted` slots_filled block
> above is guarded implicitly because `transitionApplicationStatus` returns
> nil without changing `a.Status` on no-op, but the slots_filled block
> still runs unconditionally on every ACCEPTED call). **Fix this in the same
> step**: only run the slots_filled block when the transition actually
> happened, e.g. capture `wasAlreadyAccepted := a.Status == CandStatusAccepted`
> before calling `transitionApplicationStatus`, and skip the slots_filled
> block if `wasAlreadyAccepted` — mirroring the exact guard `AcceptOffer`
> already uses (`service.go:726`, `wasAccepted`). Write a test for this
> (repeated manual `ACCEPTED` call must not double-increment `slots_filled`)
> before considering this step done.

- [ ] **Step 5: Rewire `AcceptOffer` to use the helper**

In `AcceptOffer` (`service.go:725-732`), replace:

```go
	if a, findErr := s.repo.FindApplicationByID(ctx, o.ApplicationID); findErr == nil && a != nil {
		wasAccepted := a.Status == CandStatusAccepted
		a.Status = CandStatusAccepted
		a.AcceptedAt = &now
		if err := s.repo.UpdateApplication(ctx, a); err != nil {
```

with:

```go
	if a, findErr := s.repo.FindApplicationByID(ctx, o.ApplicationID); findErr == nil && a != nil {
		wasAccepted := a.Status == CandStatusAccepted
		if err := s.transitionApplicationStatus(ctx, a, CandStatusAccepted, nil, ""); err != nil {
			s.logger.Warn("offer accepted but application transition failed",
				zap.String("offer_id", o.ID.String()), zap.String("application_id", a.ID.String()), zap.Error(err))
		}
		if err := s.repo.UpdateApplication(ctx, a); err != nil {
```

Keep the rest of the function (the `if !wasAccepted && req != nil { ... slots_filled ... }` block right after) unchanged — it already has the correct idempotency guard.

- [ ] **Step 6: Write initial history row in `CreateApplication`**

In `CreateApplication` (`service.go:1049-1082`), replace:

```go
	if err := s.repo.CreateApplication(ctx, a); err != nil {
		return nil, err
	}
	s.logger.Info("Job application created", zap.String("id", a.ID.String()))
	return applicationToResponse(a), nil
}
```

with:

```go
	if err := s.repo.CreateApplication(ctx, a); err != nil {
		return nil, err
	}
	if newStage, stageErr := s.repo.FindStageByCode(ctx, string(CandStatusNew)); stageErr == nil {
		if err := s.repo.CreateStageHistory(ctx, &ApplicationStageHistory{
			ApplicationID: a.ID,
			FromStageID:   nil,
			ToStageID:     newStage.ID,
			ChangedAt:     time.Now().UnixNano(),
		}); err != nil {
			s.logger.Warn("failed to write initial stage history", zap.String("application_id", a.ID.String()), zap.Error(err))
		}
	} else {
		s.logger.Warn("failed to look up NEW stage for initial history", zap.String("application_id", a.ID.String()), zap.Error(stageErr))
	}
	s.logger.Info("Job application created", zap.String("id", a.ID.String()))
	return applicationToResponse(a), nil
}
```

- [ ] **Step 7: Add `GetApplicationHistory` + `StageHistoryResponse`/`StageRef` DTOs**

Add to `dto.go` (near `ApplicationResponse`, `service.go:219` region):

```go
type StageRef struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type StageHistoryResponse struct {
	ID        string    `json:"id"`
	FromStage *StageRef `json:"from_stage"`
	ToStage   StageRef  `json:"to_stage"`
	ChangedBy string    `json:"changed_by,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	ChangedAt int64     `json:"changed_at"`
}
```

Add to `service.go`:

```go
func (s *Service) GetApplicationHistory(ctx context.Context, applicationID string) ([]StageHistoryResponse, error) {
	appUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, fmt.Errorf("invalid application id: %w", err)
	}
	if _, err := s.repo.FindApplicationByID(ctx, appUID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListStageHistoryByApplication(ctx, appUID)
	if err != nil {
		return nil, err
	}
	stages, err := s.repo.ListStages(ctx)
	if err != nil {
		return nil, err
	}
	byID := make(map[uuid.UUID]RecruitmentStage, len(stages))
	for _, st := range stages {
		byID[st.ID] = st
	}

	out := make([]StageHistoryResponse, 0, len(rows))
	for _, r := range rows {
		resp := StageHistoryResponse{
			ID:        r.ID.String(),
			ToStage:   StageRef{Code: byID[r.ToStageID].Code, Name: byID[r.ToStageID].Name},
			Notes:     r.Notes,
			ChangedAt: r.ChangedAt,
		}
		if r.FromStageID != nil {
			if st, ok := byID[*r.FromStageID]; ok {
				resp.FromStage = &StageRef{Code: st.Code, Name: st.Name}
			}
		}
		if r.ChangedBy != nil {
			resp.ChangedBy = r.ChangedBy.String()
		}
		out = append(out, resp)
	}
	return out, nil
}
```

- [ ] **Step 8: Run all recruitment service tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v 2>&1 | tail -100`
Expected: all tests pass, including the pre-existing
`TestService_UpdateApplicationStatus` (`NEW→SHORTLISTED`) and the new tests
from Step 1. If any pre-existing test that relies on `ACCEPTED` slots_filled
counting breaks, fix per the Step 4 note (idempotency guard) — do not weaken
the state machine to make it pass.

- [ ] **Step 9: Commit**

```bash
git add backend/internal/modules/recruitment/service.go backend/internal/modules/recruitment/service_test.go backend/internal/modules/recruitment/dto.go
git commit -m "feat: state machine + stage history wiring in UpdateApplicationStatus/AcceptOffer/CreateApplication (G-5)"
```

---

## Task 5: Handler + Routes — `GET /applications/:id/history`

**Files:**
- Modify: `backend/internal/modules/recruitment/handler.go`
- Modify: `backend/internal/modules/recruitment/routes.go`
- Test: `backend/internal/modules/recruitment/handler_test.go`

**Interfaces:**
- Consumes: `Service.GetApplicationHistory`, `ErrInvalidStatusTransition` (Task 4).
- Produces: `GET /api/v1/tenant/recruitment/applications/:id/history` (200 with list, 404 if application not found).

- [ ] **Step 1: Write failing handler tests**

Append to `handler_test.go`:

```go
func TestHandler_GetApplicationHistory(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Hist", LastName: "Test", Email: "hist@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+appID+"/history", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetApplicationHistory_NotFound(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/applications/"+uuid.New().String()+"/history", nil)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestHandler_UpdateApplicationStatus_InvalidTransitionReturns400(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	reqW := performRequest(r, "POST", "/api/v1/tenant/recruitment/requisitions", CreateRequisitionRequest{
		OrganizationID: createTestOrgID(), Title: "My Req",
	})
	var reqResp map[string]interface{}
	json.Unmarshal(reqW.Body.Bytes(), &reqResp)
	rid := reqResp["data"].(map[string]interface{})["id"].(string)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Transition", Email: "badtrans@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	appW := performRequest(r, "POST", "/api/v1/tenant/recruitment/applications", CreateApplicationRequest{
		RequisitionID: rid, CandidateID: cid,
	})
	var appResp map[string]interface{}
	json.Unmarshal(appW.Body.Bytes(), &appResp)
	appID := appResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{Status: "REJECTED"})

	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/applications/"+appID+"/status", UpdateApplicationStatusRequest{Status: "SCREENED"})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for transition out of terminal status, got %d: %s", w.Code, w.Body.String())
	}
}
```

- [ ] **Step 2: Run to verify failure**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestHandler_GetApplicationHistory -v`
Expected: 404 (route not registered) or compile error.

- [ ] **Step 3: Implement handler**

Add to `handler.go` near `GetApplicationByID`/`UpdateApplicationStatus`
(`handler.go:365-387`):

```go
func (h *Handler) GetApplicationHistory(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.svc.GetApplicationHistory(c.Request.Context(), id)
	if err != nil {
		httputil.NotFound(c, "")
		return
	}
	httputil.SuccessJSON(c, resp)
}
```

Update `UpdateApplicationStatus` to map `ErrInvalidStatusTransition` to 400
(replace the existing error handling at `handler.go:381-385`):

```go
	resp, err := h.svc.UpdateApplicationStatus(c.Request.Context(), id, req.Status, req.RejectionReason, req.Notes)
	if err != nil {
		if errors.Is(err, ErrInvalidStatusTransition) {
			httputil.BadRequest(c, err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
```

Add `"errors"` to `handler.go`'s import block if not already present.

- [ ] **Step 4: Register the route**

In `routes.go`, add after the existing `applications/:id/status` line
(`routes.go:45`):

```go
	rec.GET("/applications/:id/history", handler.GetApplicationHistory)
```

- [ ] **Step 5: Run handler tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestHandler_GetApplicationHistory|TestHandler_UpdateApplicationStatus" -v`
Expected: all PASS, including the pre-existing `TestHandler_UpdateApplicationStatus` (`NEW→SHORTLISTED`, still 200).

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/handler.go backend/internal/modules/recruitment/routes.go backend/internal/modules/recruitment/handler_test.go
git commit -m "feat: GET /applications/:id/history endpoint + 400 on invalid transition (G-5)"
```

---

## Task 6: Module wiring — `AutoMigrate` + seed 8 stages

**Files:**
- Modify: `backend/internal/modules/recruitment/module.go`

**Interfaces:**
- Consumes: `RecruitmentStage`, `ApplicationStageHistory` (Task 2).

- [ ] **Step 1: Add both entities to `Migrate()`**

In `module.go:97-107`, add to the `AutoMigrate` call:

```go
func (m *recModule) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&JobRequisition{},
		&Candidate{},
		&JobApplication{},
		&Interview{},
		&OnboardingTaskTemplate{},
		&EmployeeOnboarding{},
		&OnboardingTaskItem{},
		&RecruitmentStage{},
		&ApplicationStageHistory{},
	)
}
```

> Note: per project convention, `AutoMigrate` here is not the actual
> migration mechanism for tenant DBs (the real migration is the SQL file
> from Task 1) — this update is for interface consistency with the existing
> module pattern, not a substitute for Task 1.

- [ ] **Step 2: Seed the 8 stages idempotently**

In `Seed(db *gorm.DB) error` (`module.go:109-135`), add before the
`return nil`:

```go
	var stageCount int64
	if err := db.Model(&RecruitmentStage{}).Count(&stageCount).Error; err != nil {
		return err
	}
	if stageCount == 0 {
		stages := []RecruitmentStage{
			{Code: "NEW", Name: "New Application", SortOrder: 1},
			{Code: "SCREENED", Name: "Screened", SortOrder: 2},
			{Code: "SHORTLISTED", Name: "Shortlisted", SortOrder: 3},
			{Code: "INTERVIEWED", Name: "Interviewed", SortOrder: 4},
			{Code: "OFFERED", Name: "Offered", SortOrder: 5},
			{Code: "ACCEPTED", Name: "Accepted", SortOrder: 6},
			{Code: "REJECTED", Name: "Rejected", SortOrder: 7},
			{Code: "WITHDRAWN", Name: "Withdrawn", SortOrder: 8},
		}
		for _, st := range stages {
			if err := db.Create(&st).Error; err != nil {
				return err
			}
		}
	}
```

- [ ] **Step 2: Compile + full module test run**

Run: `cd backend && go build ./... && go test ./internal/modules/recruitment/... -v 2>&1 | tail -50`
Expected: build succeeds, all recruitment tests pass.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/recruitment/module.go
git commit -m "feat: AutoMigrate + seed 8 recruitment stages (G-5)"
```

---

## Task 7: Update plan doc status

**Files:**
- Modify: `docs/module-recruitment-development-plan.md`

- [ ] **Step 1: Update §G-5 status**

Change the G-5 section (`🔴 PIPELINE STAGE HISTORY`) header and status line
to `✅ Selesai (2026-08-12)`, following the exact style of the G-1..G-4
sections above it (bullet list of "Yang diimplementasikan", migration
number `097`, test count delta, `Ref:` line kept as-is).

- [ ] **Step 2: Update §8.2 API Plan**

Move `GET /recruitment/applications/{id}/history` from "target tambahan"
(§8.2) to §8.1 (existing, now 34 endpoints).

- [ ] **Step 3: Update §23 Definition of Done**

Check the box: `- [x] Stage transition memiliki history.`

- [ ] **Step 4: Commit**

```bash
git add docs/module-recruitment-development-plan.md
git commit -m "docs: update G-5 status to selesai (stage history)"
```
