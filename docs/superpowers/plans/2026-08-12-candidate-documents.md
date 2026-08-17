# Candidate Documents (G-6 sub-project 3a) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `candidate_documents` — a reference-only table storing file URLs (resume, cover letter, certificate, portfolio, identity, other) for candidates, reusing the existing generic upload endpoint rather than building new file storage.

**Architecture:** One new GORM entity following the exact per-entity CRUD pattern established across G-6 sub-projects 1-2 (repository → service → handler → routes). No relations to preload (unlike `CandidateSkill`), no new master tables, no new upload mechanism, no new permission.

**Tech Stack:** Go, GORM, Gin, MySQL + PostgreSQL dual migrations, existing `backend/internal/modules/recruitment` module.

## Global Constraints

- New migration number is **100** (last existing is `099_candidate_skills_certifications`).
- Migrations must be idempotent, dual-dialect (mysql+postgres), with matching `.down.sql`, following the exact style of `099_candidate_skills_certifications.sql`.
- `candidate_documents.candidate_id` FK → `candidates(id)` `ON DELETE CASCADE`.
- `document_type` is `VARCHAR(20) NOT NULL DEFAULT 'OTHER'` — enum enforced at the Gin `binding` layer only (`oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER`), no DB-level CHECK/ENUM constraint — matches this module's existing convention (e.g. `CandidateType`, `OfferStatus` are plain VARCHAR).
- No new upload endpoint — `file_url` stores whatever URL string the caller provides (expected to come from the existing `POST /api/v1/tenant/uploads` endpoint, but this table doesn't validate that).
- No new permission — reuses existing `recruitment.view/create/update/delete`.
- No `created_by`/`updated_by` columns — no table in this module tracks that (confirmed: `grep -n "CreatedBy" backend/internal/modules/recruitment/model.go` returns nothing).
- Route parameter is `:id` (not `:candidate_id`) for the nested `/candidates/:id/...` route — same established Gin wildcard-name convention as every other candidate sub-resource in this module.
- **`docs/database-schema.md` MUST be updated in this plan** (Task 7) — a prior sub-project's final review found this file had been missed for two sub-projects in a row because no single task owned it. Do not repeat that gap.

---

## Task 1: Migration 100 — `candidate_documents`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.down.sql`

**Interfaces:**
- Produces: table `candidate_documents (id, candidate_id, document_type, name, file_url, notes, created_at, updated_at)` that Task 2's GORM model must match exactly.

- [ ] **Step 1: Write postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 100_candidate_documents (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 3a: candidate_documents
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-documents-design.md)
--
-- candidate_documents — referensi dokumen kandidat (bukan binary). File
-- sesungguhnya diupload lewat endpoint generik POST /api/v1/tenant/uploads
-- (backend/internal/pkg/upload) yang mengembalikan URL; tabel ini hanya
-- menyimpan referensi URL tersebut. document_type enum di-enforce di layer
-- Gin binding (oneof=...), bukan DB constraint — pola sama dengan
-- CandidateType/OfferStatus di modul ini.
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_documents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    document_type  VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    name           VARCHAR(255) NOT NULL,
    file_url       TEXT NOT NULL,
    notes          TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_doc_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_doc_candidate ON candidate_documents (candidate_id);
```

- [ ] **Step 2: Write postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 100_candidate_documents (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_documents;
```

- [ ] **Step 3: Write mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 100_candidate_documents (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_documents (
    id             CHAR(36) PRIMARY KEY,
    candidate_id   CHAR(36) NOT NULL,
    document_type  VARCHAR(20) NOT NULL DEFAULT 'OTHER',
    name           VARCHAR(255) NOT NULL,
    file_url       TEXT NOT NULL,
    notes          TEXT NULL,
    created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_doc_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_doc_candidate ON candidate_documents (candidate_id);
```

- [ ] **Step 4: Write mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 100_candidate_documents (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_documents;
```

- [ ] **Step 5: Verify no collision, then commit**

Run: `ls backend/internal/pkg/migrator/migrations/tenant/postgres/ | grep "^100"` and same for mysql — expect exactly the 2 files each just created.

```bash
git add backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/100_candidate_documents.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/100_candidate_documents.down.sql
git commit -m "feat: migration 100 candidate_documents (G-6)"
```

---

## Task 2: GORM Model — `CandidateDocument`

**Files:**
- Modify: `backend/internal/modules/recruitment/model.go`

**Interfaces:**
- Produces: `CandidateDocument` struct — Task 3/4 depend on these exact field names.

- [ ] **Step 1: Add the `CandidateDocument` struct**

Append after the `CandidateCertification` block (same section of the file as the other G-6 sub-project entities):

```go
// =========================================================================
// CandidateDocument (G-6 — referensi dokumen kandidat)
// =========================================================================
// Referensi saja, bukan binary — file sesungguhnya diupload lewat endpoint
// generik POST /api/v1/tenant/uploads (backend/internal/pkg/upload), yang
// mengembalikan URL; tabel ini menyimpan URL tersebut. document_type
// (RESUME/COVER_LETTER/CERTIFICATE/PORTFOLIO/IDENTITY/OTHER) di-enforce di
// layer request (Gin binding), bukan constraint DB.

type CandidateDocument struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID  uuid.UUID `gorm:"type:char(36);not null;index:idx_cand_doc_candidate" json:"candidate_id"`
	DocumentType string    `gorm:"type:varchar(20);not null;default:OTHER" json:"document_type"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	FileURL      string    `gorm:"type:text;not null" json:"file_url"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CandidateDocument) TableName() string {
	return "candidate_documents"
}

func (d *CandidateDocument) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
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
git commit -m "feat: model CandidateDocument (G-6)"
```

---

## Task 3: Repository — CRUD

**Files:**
- Modify: `backend/internal/modules/recruitment/repository.go`
- Test: `backend/internal/modules/recruitment/repository_test.go`

**Interfaces:**
- Consumes: `CandidateDocument` (Task 2).
- Produces (used by Task 4 service):
  - `func (r *Repository) CreateCandidateDocument(ctx, d *CandidateDocument) error`
  - `func (r *Repository) FindCandidateDocumentByID(ctx, id uuid.UUID) (*CandidateDocument, error)`
  - `func (r *Repository) ListCandidateDocuments(ctx, candidateID uuid.UUID) ([]CandidateDocument, error)` — ordered `created_at ASC`.
  - `func (r *Repository) UpdateCandidateDocument(ctx, d *CandidateDocument) error`
  - `func (r *Repository) DeleteCandidateDocument(ctx, id uuid.UUID) error`

No `Preload` needed anywhere — `CandidateDocument` has no relation fields, so no `Omit(clause.Associations)` protection is needed on update either (unlike `CandidateSkill`'s `UpdateCandidateSkill`, which needed it because it preloads `Competency`).

- [ ] **Step 1: Write failing repository tests**

Append to `repository_test.go`:

```go
func TestRepository_CreateAndFindCandidateDocument(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Doc", LastName: "Test", Email: "doctest@test.com"}
	repo.CreateCandidate(ctx, cand)

	doc := &CandidateDocument{CandidateID: cand.ID, DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/uploads/attachments/abc.pdf"}
	if err := repo.CreateCandidateDocument(ctx, doc); err != nil {
		t.Fatalf("CreateCandidateDocument failed: %v", err)
	}

	found, err := repo.FindCandidateDocumentByID(ctx, doc.ID)
	if err != nil {
		t.Fatalf("FindCandidateDocumentByID failed: %v", err)
	}
	if found.Name != "resume.pdf" {
		t.Errorf("expected name 'resume.pdf', got %s", found.Name)
	}
}

func TestRepository_ListCandidateDocuments(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Doc", Email: "listdoc@test.com"}
	repo.CreateCandidate(ctx, cand)
	repo.CreateCandidateDocument(ctx, &CandidateDocument{CandidateID: cand.ID, DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/u/a.pdf"})
	repo.CreateCandidateDocument(ctx, &CandidateDocument{CandidateID: cand.ID, DocumentType: "PORTFOLIO", Name: "portfolio.pdf", FileURL: "/u/b.pdf"})

	list, err := repo.ListCandidateDocuments(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateDocuments failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 documents, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateDocument(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Doc", Email: "upddoc@test.com"}
	repo.CreateCandidate(ctx, cand)
	doc := &CandidateDocument{CandidateID: cand.ID, DocumentType: "OTHER", Name: "Original", FileURL: "/u/c.pdf"}
	repo.CreateCandidateDocument(ctx, doc)

	doc.Name = "Updated"
	if err := repo.UpdateCandidateDocument(ctx, doc); err != nil {
		t.Fatalf("UpdateCandidateDocument failed: %v", err)
	}
	found, _ := repo.FindCandidateDocumentByID(ctx, doc.ID)
	if found.Name != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.Name)
	}

	if err := repo.DeleteCandidateDocument(ctx, doc.ID); err != nil {
		t.Fatalf("DeleteCandidateDocument failed: %v", err)
	}
	if _, err := repo.FindCandidateDocumentByID(ctx, doc.ID); err == nil {
		t.Error("expected error finding deleted document, got nil")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestRepository_CreateAndFindCandidateDocument -v`
Expected: compile error.

- [ ] **Step 3: Implement repository methods**

Append to `repository.go` (after the Candidate Certifications section):

```go
// =========================================================================
// Candidate Documents (G-6)
// =========================================================================

func (r *Repository) CreateCandidateDocument(ctx context.Context, d *CandidateDocument) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(d).Error
}

func (r *Repository) FindCandidateDocumentByID(ctx context.Context, id uuid.UUID) (*CandidateDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var d CandidateDocument
	if err := db.WithContext(ctx).First(&d, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate document not found")
		}
		return nil, err
	}
	return &d, nil
}

func (r *Repository) ListCandidateDocuments(ctx context.Context, candidateID uuid.UUID) ([]CandidateDocument, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateDocument
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateDocument(ctx context.Context, d *CandidateDocument) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(d).Error
}

func (r *Repository) DeleteCandidateDocument(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateDocument{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate document not found")
	}
	return result.Error
}
```

- [ ] **Step 4: Add `CandidateDocument` to `helpers_test.go`'s AutoMigrate list**

`backend/internal/modules/recruitment/helpers_test.go`'s `setupTestDB()` needs `&CandidateDocument{}` added to its `AutoMigrate(...)` call.

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestRepository_.*CandidateDocument" -v`
Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/repository.go backend/internal/modules/recruitment/repository_test.go backend/internal/modules/recruitment/helpers_test.go
git commit -m "feat: repository CRUD CandidateDocument (G-6)"
```

---

## Task 4: Service — CRUD, candidate-existence guard

**Files:**
- Modify: `backend/internal/modules/recruitment/service.go`
- Modify: `backend/internal/modules/recruitment/dto.go`
- Test: `backend/internal/modules/recruitment/service_test.go`

**Interfaces:**
- Consumes: repository methods from Task 3.
- Produces:
  - `func (s *Service) CreateCandidateDocument(ctx, candidateID string, req CreateCandidateDocumentRequest) (*CandidateDocumentResponse, error)` — candidate-existence guard.
  - `func (s *Service) ListCandidateDocuments(ctx, candidateID string) ([]CandidateDocumentResponse, error)`
  - `func (s *Service) UpdateCandidateDocument(ctx, id string, req UpdateCandidateDocumentRequest) (*CandidateDocumentResponse, error)`
  - `func (s *Service) DeleteCandidateDocument(ctx, id string) error`

- [ ] **Step 1: Add DTOs to `dto.go`**

```go
// =========================================================================
// Candidate Document DTOs (G-6)
// =========================================================================

type CreateCandidateDocumentRequest struct {
	DocumentType string `json:"document_type" binding:"omitempty,oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER"`
	Name         string `json:"name" binding:"required,max=255"`
	FileURL      string `json:"file_url" binding:"required"`
	Notes        string `json:"notes"`
}

type UpdateCandidateDocumentRequest struct {
	DocumentType *string `json:"document_type" binding:"omitempty,oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER"`
	Name         *string `json:"name" binding:"omitempty,max=255"`
	FileURL      *string `json:"file_url" binding:"omitempty"`
	Notes        *string `json:"notes"`
}

type CandidateDocumentResponse struct {
	ID           string `json:"id"`
	CandidateID  string `json:"candidate_id"`
	DocumentType string `json:"document_type"`
	Name         string `json:"name"`
	FileURL      string `json:"file_url"`
	Notes        string `json:"notes,omitempty"`
}
```

- [ ] **Step 2: Write failing service tests**

Append to `service_test.go`:

```go
func TestService_CreateCandidateDocument(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Doc", LastName: "Svc", Email: "docsvc@test.com"})

	resp, err := svc.CreateCandidateDocument(ctx, cand.ID, CreateCandidateDocumentRequest{
		DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/uploads/attachments/x.pdf",
	})
	if err != nil {
		t.Fatalf("CreateCandidateDocument failed: %v", err)
	}
	if resp.Name != "resume.pdf" {
		t.Errorf("expected name 'resume.pdf', got %s", resp.Name)
	}
	if resp.DocumentType != "RESUME" {
		t.Errorf("expected document_type 'RESUME', got %s", resp.DocumentType)
	}
}

func TestService_CreateCandidateDocument_UnknownCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateCandidateDocument(ctx, uuid.New().String(), CreateCandidateDocumentRequest{
		Name: "x.pdf", FileURL: "/u/x.pdf",
	})
	if err == nil {
		t.Fatal("expected error for unknown candidate, got nil")
	}
}

func TestService_ListCandidateDocuments(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "List", LastName: "Doc", Email: "listdocsvc@test.com"})
	svc.CreateCandidateDocument(ctx, cand.ID, CreateCandidateDocumentRequest{Name: "a.pdf", FileURL: "/u/a.pdf"})
	svc.CreateCandidateDocument(ctx, cand.ID, CreateCandidateDocumentRequest{Name: "b.pdf", FileURL: "/u/b.pdf"})

	list, err := svc.ListCandidateDocuments(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateDocuments failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestService_CreateCandidateDocument -v`
Expected: compile error.

- [ ] **Step 4: Implement service CRUD methods**

```go
// =========================================================================
// Candidate Documents (G-6)
// =========================================================================

func (s *Service) CreateCandidateDocument(ctx context.Context, candidateID string, req CreateCandidateDocumentRequest) (*CandidateDocumentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	docType := req.DocumentType
	if docType == "" {
		docType = "OTHER"
	}
	d := &CandidateDocument{
		CandidateID:  candUUID,
		DocumentType: docType,
		Name:         req.Name,
		FileURL:      req.FileURL,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateCandidateDocument(ctx, d); err != nil {
		return nil, err
	}
	return candidateDocumentToResponse(d), nil
}

func (s *Service) ListCandidateDocuments(ctx context.Context, candidateID string) ([]CandidateDocumentResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateDocuments(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateDocumentResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateDocumentToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateDocument(ctx context.Context, id string, req UpdateCandidateDocumentRequest) (*CandidateDocumentResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	d, err := s.repo.FindCandidateDocumentByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.DocumentType != nil {
		d.DocumentType = *req.DocumentType
	}
	if req.Name != nil {
		d.Name = *req.Name
	}
	if req.FileURL != nil {
		d.FileURL = *req.FileURL
	}
	if req.Notes != nil {
		d.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateDocument(ctx, d); err != nil {
		return nil, err
	}
	return candidateDocumentToResponse(d), nil
}

func (s *Service) DeleteCandidateDocument(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateDocument(ctx, uid)
}

func candidateDocumentToResponse(d *CandidateDocument) *CandidateDocumentResponse {
	return &CandidateDocumentResponse{
		ID:           d.ID.String(),
		CandidateID:  d.CandidateID.String(),
		DocumentType: d.DocumentType,
		Name:         d.Name,
		FileURL:      d.FileURL,
		Notes:        d.Notes,
	}
}
```

- [ ] **Step 5: Run all recruitment service tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | tail -100`
Expected: all pass, including every pre-existing test.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/service.go backend/internal/modules/recruitment/service_test.go backend/internal/modules/recruitment/dto.go
git commit -m "feat: service CRUD CandidateDocument (G-6)"
```

---

## Task 5: Handler + Routes

**Files:**
- Modify: `backend/internal/modules/recruitment/handler.go`
- Modify: `backend/internal/modules/recruitment/routes.go`
- Test: `backend/internal/modules/recruitment/handler_test.go`

**Interfaces:**
- Consumes: service methods from Task 4.
- Produces:
  ```
  POST   /recruitment/candidates/:id/documents
  GET    /recruitment/candidates/:id/documents
  PUT    /recruitment/documents/:id
  DELETE /recruitment/documents/:id
  ```

- [ ] **Step 1: Write failing handler tests**

Append to `handler_test.go`:

```go
func TestHandler_CreateCandidateDocument(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Doc", LastName: "Handler", Email: "dochandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{
		DocumentType: "RESUME", Name: "resume.pdf", FileURL: "/uploads/attachments/x.pdf",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateDocument_InvalidDocumentType(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Bad", LastName: "Type", Email: "badtype@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{
		DocumentType: "NOT_A_REAL_TYPE", Name: "x.pdf", FileURL: "/u/x.pdf",
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid document_type, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCandidateDocuments(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "DocH", Email: "listdoch@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", CreateCandidateDocumentRequest{Name: "a.pdf", FileURL: "/u/a.pdf"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/documents", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
```

- [ ] **Step 2: Run to verify failure**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestHandler_CreateCandidateDocument -v`
Expected: 404 (route not registered) or compile error.

- [ ] **Step 3: Implement handlers**

Add to `handler.go`, mirroring `CreateCandidateSkill`/`ListCandidateSkills`/`UpdateCandidateSkill`/`DeleteCandidateSkill`:

```go
func (h *Handler) CreateCandidateDocument(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateDocument(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateDocuments(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateDocuments(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateDocument(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateDocumentRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateDocument(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateDocument(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateDocument(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
```

- [ ] **Step 4: Register routes**

In `routes.go`, add after the existing `certifications/:id` routes:

```go
	rec.POST("/candidates/:id/documents", handler.CreateCandidateDocument)
	rec.GET("/candidates/:id/documents", handler.ListCandidateDocuments)
	rec.PUT("/documents/:id", handler.UpdateCandidateDocument)
	rec.DELETE("/documents/:id", handler.DeleteCandidateDocument)
```

- [ ] **Step 5: Run handler tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | tail -60`
Expected: all pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/handler.go backend/internal/modules/recruitment/routes.go backend/internal/modules/recruitment/handler_test.go
git commit -m "feat: handler+routes CandidateDocument CRUD (G-6)"
```

---

## Task 6: Module wiring — `AutoMigrate`

**Files:**
- Modify: `backend/internal/modules/recruitment/module.go`

**Interfaces:**
- Consumes: `CandidateDocument` (Task 2).

- [ ] **Step 1: Add to `Migrate()`**

Add `&CandidateDocument{}` to the existing `db.AutoMigrate(...)` call in `module.go`'s `Migrate()` method, alongside the other G-6 entities.

> Note: as established in prior sub-projects' final reviews, this `Migrate()`/`Seed()` wiring is not actually invoked for tenant databases in production (a separately-tracked, pre-existing bug). This step is for consistency; the real schema source of truth is migration 100 from Task 1.

Do NOT touch `helpers_test.go` — Task 3 already added `&CandidateDocument{}` there.

- [ ] **Step 2: Build + test run**

Run: `cd backend && go build ./... && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | tail -60`
Expected: clean build, all tests pass.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/recruitment/module.go
git commit -m "feat: AutoMigrate CandidateDocument (G-6)"
```

---

## Task 7: Update plan doc status + `docs/database-schema.md`

**Files:**
- Modify: `docs/archive/module-recruitment-development-plan.md`
- Modify: `docs/database-schema.md`

**This task's scope explicitly includes `docs/database-schema.md` — do not skip it.** A prior sub-project's final review found this file had been missed for two sub-projects running because no task in those plans named it explicitly. This plan names it here specifically to close that gap.

- [ ] **Step 1: Update §G-6 status in `module-recruitment-development-plan.md`**

Add a "Sub-project 3a" entry to the G-6 section (search for "sub-project 2/3" to find the current state), following the exact style of sub-projects 1 and 2 (migration number, model/service/handler summary, DTO names, accurate test count with running total). Before writing the test count, run these yourself and use the real numbers — do not guess or copy from this plan:

```bash
cd backend && go test ./internal/modules/recruitment/... -v -count=1 2>&1 | grep -c "^--- PASS:"
grep -c "^func Test" internal/modules/recruitment/handler_test.go internal/modules/recruitment/repository_test.go internal/modules/recruitment/service_test.go
```

Update the "Rencana (sisa G-6)" deferred list: remove `candidate_documents`, leaving only `candidates.status`, `source_id`, `candidate_consents` as open.

Double-check every schema/type claim you write (column names, types, DTO names) against the actual current code in `model.go`/the migration file/`dto.go` — this exact category of error (wrong module name, wrong column name, wrong type) was found and had to be fixed in a prior sub-project's final review. Read the real code, don't reconstruct from memory of the plan.

- [ ] **Step 2: Update §8.1 API Plan**

Add the 4 new endpoints to §8.1 (existing), increment the endpoint count (verify the current count first by reading the doc, don't assume a number).

- [ ] **Step 3: Update `docs/database-schema.md`**

Read the existing `candidate_skills`/`candidate_certifications` entries in this file first (search for "candidate_skills" — both the table-inventory row around the "## Recruitment" section and the Mermaid ER block entry) to match their exact format. Add:
- A table-inventory row: `| \`candidate_documents\` | 8 | candidates.id |` (verify the column count — 8 — against the actual migration file, don't trust this number blindly).
- A Mermaid ER entity block for `candidate_documents` listing all 8 columns with their types, matching the style of the `candidate_certifications` block immediately above it (e.g. `CHAR id`, `CHAR candidate_id`, `VARCHAR document_type`, `VARCHAR name`, `TEXT file_url`, `TEXT notes`, `TIMESTAMP created_at`, `TIMESTAMP updated_at`).

- [ ] **Step 4: Commit**

```bash
git add docs/module-recruitment-development-plan.md docs/database-schema.md
git commit -m "docs: update G-6 status - candidate documents selesai (sub-project 3a)"
```
