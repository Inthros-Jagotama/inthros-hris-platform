# Candidate Consents (G-6 sub-project 3b) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `candidate_consents` — an append-only audit log recording GRANTED/REVOKED data-processing consent entries for candidates, closing out the last of the originally-planned G-6 sub-tables.

**Architecture:** One new GORM entity, but a NARROWER CRUD surface than every other G-6 sub-resource so far — only `Create` (append) and `List` (read history), no `Update`/`Delete`, since the table is append-only by design. `changed_by` is populated from the authenticated user's id in the request context, following the exact pattern already established for `job_application_stage_histories.changed_by` (G-5).

**Tech Stack:** Go, GORM, Gin, MySQL + PostgreSQL dual migrations, existing `backend/internal/modules/recruitment` module.

**Spec:** `docs/superpowers/specs/2026-08-12-candidate-consents-design.md`

## Global Constraints

- New migration number is **101** (last existing is `100_candidate_documents`).
- Migrations must be idempotent, dual-dialect (mysql+postgres), with matching `.down.sql`, following the exact style of `100_candidate_documents.sql`.
- `candidate_consents.candidate_id` FK → `candidates(id)` `ON DELETE CASCADE`.
- `action` is `VARCHAR(20) NOT NULL`, no DB default (unlike `candidate_documents.document_type`'s `DEFAULT 'OTHER'`) — there's no sensible default for "did they consent or not", so the request field is REQUIRED, not optional-with-default.
- **No `Update`/`Delete` repository methods, service methods, handler methods, or routes** — this is a deliberate scope narrowing versus every other G-6 sub-resource; do not add them "for consistency," that would violate the append-only design.
- `changed_by` populated via `c.GetString("user_id")` in the handler, parsed to `*uuid.UUID` (nil on empty/parse failure — never fail the request over this), passed down to the service — exact pattern already in `handler.go`'s `UpdateApplicationStatus` (see `handler.go:612-624`).
- `changed_at` is unix nano, set server-side (`time.Now().UnixNano()`) in the service, not accepted from the request.
- Route parameter is `:id` (not `:candidate_id`) for `/candidates/:id/consents`, matching the module's established Gin wildcard-name convention.
- No new permission — reuses existing `recruitment.view/create/update/delete`.
- **`docs/database-schema.md` MUST be updated in Task 7** — named explicitly here per the convention established after two prior sub-projects missed it.

---

## Task 1: Migration 101 — `candidate_consents`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.down.sql`

**Interfaces:**
- Produces: table `candidate_consents (id, candidate_id, action, notes, changed_by, changed_at, created_at)` that Task 2's GORM model must match exactly.

- [ ] **Step 1: Write postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 101_candidate_consents (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 3b: candidate_consents (LAST of the originally-planned G-6
-- sub-tables — closes out education/experience/skills/certifications/
-- documents/consents; only candidates.status/source_id remain open, both
-- explicitly skipped by earlier design decisions, not deferred-to-later)
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-consents-design.md)
--
-- candidate_consents — append-only audit log of data-processing consent
-- (GDPR-style). action GRANTED|REVOKED enforced at Gin binding layer, not
-- DB constraint (module convention). No updated_at — rows are never
-- updated, only inserted. changed_by = staff user who recorded the entry
-- (this system has no public candidate-facing portal, so consent is
-- documented by HR/recruiter, not self-service).
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_consents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    action         VARCHAR(20) NOT NULL,
    notes          TEXT NULL,
    changed_by     CHAR(36) NULL,
    changed_at     BIGINT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_consent_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_consent_candidate ON candidate_consents (candidate_id);
```

- [ ] **Step 2: Write postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 101_candidate_consents (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_consents;
```

- [ ] **Step 3: Write mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 101_candidate_consents (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_consents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    action         VARCHAR(20) NOT NULL,
    notes          TEXT NULL,
    changed_by     CHAR(36) NULL,
    changed_at     BIGINT NOT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_consent_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_consent_candidate ON candidate_consents (candidate_id);
```

- [ ] **Step 4: Write mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 101_candidate_consents (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_consents;
```

- [ ] **Step 5: Verify no collision, then commit**

Run: `ls backend/internal/pkg/migrator/migrations/tenant/postgres/ | grep "^101"` and same for mysql — expect exactly the 2 files each just created.

```bash
git add backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/101_candidate_consents.down.sql
git commit -m "feat: migration 101 candidate_consents (G-6)"
```

---

## Task 2: GORM Model — `CandidateConsent`

**Files:**
- Modify: `backend/internal/modules/recruitment/model.go`

**Interfaces:**
- Produces: `CandidateConsent` struct — Task 3/4 depend on these exact field names.

- [ ] **Step 1: Add the `CandidateConsent` struct**

Append after the `CandidateDocument` block:

```go
// =========================================================================
// CandidateConsent (G-6 — audit log consent pemrosesan data pribadi)
// =========================================================================
// Append-only — tidak ada Update/Delete. changed_by = staff yang mencatat
// (sistem ini tidak punya portal publik untuk kandidat, jadi consent
// didokumentasikan HR/recruiter, bukan self-service). action GRANTED/
// REVOKED di-enforce di layer Gin binding, bukan DB constraint.

type CandidateConsent struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID uuid.UUID  `gorm:"type:char(36);not null;index:idx_cand_consent_candidate" json:"candidate_id"`
	Action      string     `gorm:"type:varchar(20);not null" json:"action"`
	Notes       string     `gorm:"type:text" json:"notes,omitempty"`
	ChangedBy   *uuid.UUID `gorm:"type:char(36)" json:"changed_by,omitempty"`
	ChangedAt   int64      `gorm:"type:bigint;not null" json:"changed_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func (CandidateConsent) TableName() string {
	return "candidate_consents"
}

func (c *CandidateConsent) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
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
git commit -m "feat: model CandidateConsent (G-6)"
```

---

## Task 3: Repository — Create + List only

**Files:**
- Modify: `backend/internal/modules/recruitment/repository.go`
- Test: `backend/internal/modules/recruitment/repository_test.go`

**Interfaces:**
- Consumes: `CandidateConsent` (Task 2).
- Produces (used by Task 4 service) — ONLY 2 methods, not the usual 5-method CRUD shape:
  - `func (r *Repository) CreateCandidateConsent(ctx, c *CandidateConsent) error`
  - `func (r *Repository) ListCandidateConsents(ctx, candidateID uuid.UUID) ([]CandidateConsent, error)` — ordered `changed_at ASC`.

Do NOT add `FindCandidateConsentByID`, `UpdateCandidateConsent`, or `DeleteCandidateConsent` — there is no update/delete API surface for this entity per the append-only design (Global Constraints).

- [ ] **Step 1: Write failing repository tests**

Append to `repository_test.go`:

```go
func TestRepository_CreateAndListCandidateConsents(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Consent", LastName: "Test", Email: "consenttest@test.com"}
	repo.CreateCandidate(ctx, cand)

	consent := &CandidateConsent{CandidateID: cand.ID, Action: "GRANTED", ChangedAt: 1000}
	if err := repo.CreateCandidateConsent(ctx, consent); err != nil {
		t.Fatalf("CreateCandidateConsent failed: %v", err)
	}

	list, err := repo.ListCandidateConsents(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateConsents failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 consent entry, got %d", len(list))
	}
	if list[0].Action != "GRANTED" {
		t.Errorf("expected action 'GRANTED', got %s", list[0].Action)
	}
}

func TestRepository_ListCandidateConsents_OrderedByChangedAt(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Order", LastName: "Consent", Email: "orderconsent@test.com"}
	repo.CreateCandidate(ctx, cand)

	repo.CreateCandidateConsent(ctx, &CandidateConsent{CandidateID: cand.ID, Action: "GRANTED", ChangedAt: 2000})
	repo.CreateCandidateConsent(ctx, &CandidateConsent{CandidateID: cand.ID, Action: "REVOKED", ChangedAt: 1000})

	list, err := repo.ListCandidateConsents(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateConsents failed: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	if list[0].ChangedAt != 1000 || list[1].ChangedAt != 2000 {
		t.Errorf("expected ascending changed_at order [1000, 2000], got [%d, %d]", list[0].ChangedAt, list[1].ChangedAt)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestRepository_CreateAndListCandidateConsents -v`
Expected: compile error.

- [ ] **Step 3: Implement repository methods**

Append to `repository.go` (after the Candidate Documents section):

```go
// =========================================================================
// Candidate Consents (G-6) — append-only, no Update/Delete
// =========================================================================

func (r *Repository) CreateCandidateConsent(ctx context.Context, c *CandidateConsent) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) ListCandidateConsents(ctx context.Context, candidateID uuid.UUID) ([]CandidateConsent, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateConsent
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("changed_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
```

- [ ] **Step 4: Add `CandidateConsent` to `helpers_test.go`'s AutoMigrate list**

`backend/internal/modules/recruitment/helpers_test.go`'s `setupTestDB()` needs `&CandidateConsent{}` added to its `AutoMigrate(...)` call.

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestRepository_.*CandidateConsent" -v`
Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/repository.go backend/internal/modules/recruitment/repository_test.go backend/internal/modules/recruitment/helpers_test.go
git commit -m "feat: repository create+list CandidateConsent (G-6)"
```

---

## Task 4: Service — Create (with `changed_by`/`changed_at`) + List

**Files:**
- Modify: `backend/internal/modules/recruitment/service.go`
- Modify: `backend/internal/modules/recruitment/dto.go`
- Test: `backend/internal/modules/recruitment/service_test.go`

**Interfaces:**
- Consumes: repository methods from Task 3.
- Produces — ONLY 2 methods:
  - `func (s *Service) CreateCandidateConsent(ctx, candidateID string, req CreateCandidateConsentRequest, changedBy *uuid.UUID) (*CandidateConsentResponse, error)` — candidate-existence guard, sets `ChangedAt` server-side.
  - `func (s *Service) ListCandidateConsents(ctx, candidateID string) ([]CandidateConsentResponse, error)`

Note the extra `changedBy *uuid.UUID` parameter on `CreateCandidateConsent` — this is different from every other G-6 sub-resource's create method signature, because the actor id must flow from the HTTP layer (Task 5's handler) down through the service into the model, exactly mirroring `UpdateApplicationStatus`'s existing `changedBy` parameter (`service.go`, check its signature for the exact parameter-passing convention before writing this).

- [ ] **Step 1: Add DTOs to `dto.go`**

```go
// =========================================================================
// Candidate Consent DTOs (G-6) — append-only, no Update
// =========================================================================

type CreateCandidateConsentRequest struct {
	Action string `json:"action" binding:"required,oneof=GRANTED REVOKED"`
	Notes  string `json:"notes"`
}

type CandidateConsentResponse struct {
	ID          string `json:"id"`
	CandidateID string `json:"candidate_id"`
	Action      string `json:"action"`
	Notes       string `json:"notes,omitempty"`
	ChangedBy   string `json:"changed_by,omitempty"`
	ChangedAt   int64  `json:"changed_at"`
}
```

- [ ] **Step 2: Write failing service tests**

Append to `service_test.go`:

```go
func TestService_CreateCandidateConsent(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Consent", LastName: "Svc", Email: "consentsvc@test.com"})

	resp, err := svc.CreateCandidateConsent(ctx, cand.ID, CreateCandidateConsentRequest{Action: "GRANTED"}, nil)
	if err != nil {
		t.Fatalf("CreateCandidateConsent failed: %v", err)
	}
	if resp.Action != "GRANTED" {
		t.Errorf("expected action 'GRANTED', got %s", resp.Action)
	}
	if resp.ChangedAt == 0 {
		t.Error("expected changed_at to be set server-side, got 0")
	}
}

func TestService_CreateCandidateConsent_WithChangedBy(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Actor", LastName: "Consent", Email: "actorconsent@test.com"})
	actorID := uuid.New()

	resp, err := svc.CreateCandidateConsent(ctx, cand.ID, CreateCandidateConsentRequest{Action: "GRANTED"}, &actorID)
	if err != nil {
		t.Fatalf("CreateCandidateConsent failed: %v", err)
	}
	if resp.ChangedBy != actorID.String() {
		t.Errorf("expected changed_by %s, got %s", actorID.String(), resp.ChangedBy)
	}
}

func TestService_CreateCandidateConsent_UnknownCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateCandidateConsent(ctx, uuid.New().String(), CreateCandidateConsentRequest{Action: "GRANTED"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown candidate, got nil")
	}
}

func TestService_ListCandidateConsents(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "List", LastName: "Consent", Email: "listconsent@test.com"})
	svc.CreateCandidateConsent(ctx, cand.ID, CreateCandidateConsentRequest{Action: "GRANTED"}, nil)
	svc.CreateCandidateConsent(ctx, cand.ID, CreateCandidateConsentRequest{Action: "REVOKED"}, nil)

	list, err := svc.ListCandidateConsents(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateConsents failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestService_CreateCandidateConsent -v`
Expected: compile error.

- [ ] **Step 4: Implement service methods**

```go
// =========================================================================
// Candidate Consents (G-6) — append-only, no Update/Delete
// =========================================================================

func (s *Service) CreateCandidateConsent(ctx context.Context, candidateID string, req CreateCandidateConsentRequest, changedBy *uuid.UUID) (*CandidateConsentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	c := &CandidateConsent{
		CandidateID: candUUID,
		Action:      req.Action,
		Notes:       req.Notes,
		ChangedBy:   changedBy,
		ChangedAt:   time.Now().UnixNano(),
	}
	if err := s.repo.CreateCandidateConsent(ctx, c); err != nil {
		return nil, err
	}
	return candidateConsentToResponse(c), nil
}

func (s *Service) ListCandidateConsents(ctx context.Context, candidateID string) ([]CandidateConsentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateConsents(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateConsentResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateConsentToResponse(&list[i]))
	}
	return out, nil
}

func candidateConsentToResponse(c *CandidateConsent) *CandidateConsentResponse {
	resp := &CandidateConsentResponse{
		ID:          c.ID.String(),
		CandidateID: c.CandidateID.String(),
		Action:      c.Action,
		Notes:       c.Notes,
		ChangedAt:   c.ChangedAt,
	}
	if c.ChangedBy != nil {
		resp.ChangedBy = c.ChangedBy.String()
	}
	return resp
}
```

- [ ] **Step 5: Run all recruitment service tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | tail -100`
Expected: all pass, including every pre-existing test.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/service.go backend/internal/modules/recruitment/service_test.go backend/internal/modules/recruitment/dto.go
git commit -m "feat: service create+list CandidateConsent with changed_by (G-6)"
```

---

## Task 5: Handler + Routes

**Files:**
- Modify: `backend/internal/modules/recruitment/handler.go`
- Modify: `backend/internal/modules/recruitment/routes.go`
- Test: `backend/internal/modules/recruitment/handler_test.go`

**Interfaces:**
- Consumes: service methods from Task 4.
- Produces — ONLY 2 endpoints:
  ```
  POST   /recruitment/candidates/:id/consents
  GET    /recruitment/candidates/:id/consents
  ```

- [ ] **Step 1: Write failing handler tests**

Append to `handler_test.go`:

```go
func TestHandler_CreateCandidateConsent(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Consent", LastName: "Handler", Email: "consenthandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "GRANTED",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateConsent_InvalidAction(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Action", Email: "badaction@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "MAYBE",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid action, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateConsent_ChangedByFromContext(t *testing.T) {
	userID := uuid.New().String()
	r, _, cleanup := setupTestRouterWithUserID(userID)
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Actor", LastName: "Handler", Email: "actorhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{
		Action: "GRANTED",
	})
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	if data["changed_by"] != userID {
		t.Errorf("expected changed_by %s, got %v", userID, data["changed_by"])
	}
}

func TestHandler_ListCandidateConsents(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "ConsentH", Email: "listconsenth@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", CreateCandidateConsentRequest{Action: "GRANTED"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/consents", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
```

`setupTestRouterWithUserID(userID string)` already exists in `handler_test.go` (added during the G-5 `changed_by` work) — confirm it before writing this test, don't redefine it.

- [ ] **Step 2: Run to verify failure**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestHandler_CreateCandidateConsent -v`
Expected: 404 (route not registered) or compile error.

- [ ] **Step 3: Implement handlers**

Add to `handler.go`, mirroring `UpdateApplicationStatus`'s `changed_by` extraction pattern (`handler.go:612-624`) combined with the create/list shape from sibling candidate sub-resources:

```go
func (h *Handler) CreateCandidateConsent(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateConsentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	var changedBy *uuid.UUID
	if userIDStr := c.GetString("user_id"); userIDStr != "" {
		if uid, err := uuid.Parse(userIDStr); err == nil {
			changedBy = &uid
		}
	}
	resp, err := h.svc.CreateCandidateConsent(c.Request.Context(), candidateID, req, changedBy)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateConsents(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateConsents(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}
```

- [ ] **Step 4: Register routes**

In `routes.go`, add after the existing `documents/:id` routes:

```go
	rec.POST("/candidates/:id/consents", handler.CreateCandidateConsent)
	rec.GET("/candidates/:id/consents", handler.ListCandidateConsents)
```

- [ ] **Step 5: Run handler tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | tail -60`
Expected: all pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/handler.go backend/internal/modules/recruitment/routes.go backend/internal/modules/recruitment/handler_test.go
git commit -m "feat: handler+routes CandidateConsent create+list (G-6)"
```

---

## Task 6: Module wiring — `AutoMigrate`

**Files:**
- Modify: `backend/internal/modules/recruitment/module.go`

**Interfaces:**
- Consumes: `CandidateConsent` (Task 2).

- [ ] **Step 1: Add to `Migrate()`**

Add `&CandidateConsent{}` to the existing `db.AutoMigrate(...)` call in `module.go`'s `Migrate()` method.

> Note: as established in prior sub-projects' final reviews, this `Migrate()`/`Seed()` wiring is not actually invoked for tenant databases in production. This step is for consistency; the real schema source of truth is migration 101 from Task 1.

Do NOT touch `helpers_test.go` — Task 3 already added `&CandidateConsent{}` there.

- [ ] **Step 2: Build + test run, get exact counts**

Run:
```bash
cd backend && go build ./...
go test ./internal/modules/recruitment/... -v -count=1 2>&1 | grep -c "^--- PASS:"
go test ./internal/modules/recruitment/... -v -count=1 2>&1 | grep -c "^--- FAIL:"
```
Expected: clean build, FAIL count 0. Record the exact PASS count for use in Task 7 (do not guess it in Task 7 — re-verify there too).

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/recruitment/module.go
git commit -m "feat: AutoMigrate CandidateConsent (G-6)"
```

---

## Task 7: Update plan doc + `docs/database-schema.md`

**Files:**
- Modify: `docs/module-recruitment-development-plan.md`
- Modify: `docs/database-schema.md`

**This task's scope explicitly includes `docs/database-schema.md` — do not skip it.** This is the third G-6 sub-project in a row to name this file explicitly (sub-project 3a was the first to get it right, after two earlier sub-projects missed it and needed final-review fixes).

**Before writing ANY column name, type, or count into either doc, read the actual current code** — `backend/internal/pkg/migrator/migrations/tenant/postgres/101_candidate_consents.sql`, `backend/internal/modules/recruitment/model.go`'s `CandidateConsent` struct, `dto.go`'s `CandidateConsent*` types. Do not reconstruct from memory of this plan.

- [ ] **Step 1: Update §G-6 status in `module-recruitment-development-plan.md`**

Add a "Sub-project 3b" entry to the G-6 section (search for "sub-project 3a" to find the current state), following the exact style of the prior sub-project entries: migration number 101, model/service/handler/DTO summary, note the append-only design (no Update/Delete) explicitly since it's a deliberate deviation from every prior sub-project's 5-method CRUD shape, and an accurate test count with running total. Verify the test count yourself:

```bash
cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | grep -c "^--- PASS:"
grep -c "^func Test" internal/modules/recruitment/handler_test.go internal/modules/recruitment/repository_test.go internal/modules/recruitment/service_test.go
```

Update the "Rencana (sisa G-6)" deferred list: remove `candidate_consents`, leaving only `candidates.status` and `source_id` as the sole remaining open G-6 items (both explicitly skipped by earlier design decisions during sub-project 1's brainstorming, not "deferred to a future sub-project" — this G-6 gap-analysis item's sub-table work is now fully closed out). Consider whether the G-6 section header itself should now read differently (e.g. no longer "🔶 partial" for the sub-table work specifically) given all planned sub-tables are done — use judgment, but don't overclaim: `status`/`source_id` are still open, so G-6 as a whole is not 100% complete.

- [ ] **Step 2: Update §8.1 API Plan**

Add the 2 new endpoints to §8.1, increment the endpoint count (read the doc first to get the current count, don't assume).

- [ ] **Step 3: Sweep for other stale references (learned from sub-project 3a's final review finding)**

Search the whole document for other places that might reference "consents" as still-open/deferred and could now be stale — check §23 Definition of Done specifically (search for the DoD line about "profile terstruktur" that was already updated in sub-project 3a to include documents — it may still say something like "consents ❌" and need updating to ✅, or may not mention consents at all, in which case no change is needed there). Do not leave a contradiction between this update and any other section.

- [ ] **Step 4: Update `docs/database-schema.md`**

Read the existing `candidate_documents` entry first (table-inventory row + Mermaid ER block, both in the "## Recruitment" section) to match its exact format. Add:
- A table-inventory row: `| \`candidate_consents\` | 7 | candidates.id |` (verify the column count — 7 — against the actual migration file, don't trust this number blindly).
- A Mermaid ER entity block for `candidate_consents` listing all 7 columns with types, matching the style of the `candidate_documents` block.

- [ ] **Step 5: Commit**

```bash
git add docs/module-recruitment-development-plan.md docs/database-schema.md
git commit -m "docs: update G-6 status - candidate consents selesai (sub-project 3b, G-6 sub-tables complete)"
```
