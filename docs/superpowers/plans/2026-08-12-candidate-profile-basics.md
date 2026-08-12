# Candidate Profile Basics (G-6 sub-project 1) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `candidates.candidate_number` (auto-generated) and two new sub-tables — `candidate_educations` and `candidate_work_experiences` — to the recruitment module, so a candidate's education and employment history can be recorded structurally.

**Architecture:** Two new GORM entities following the recruitment module's existing per-entity CRUD pattern (repository → service → handler → routes), plus one column addition to `Candidate`. `candidate_educations.education_id`/`education_major_id` reference the existing seeded `setting.Education`/`setting.EducationMajor` master tables, mirroring the already-shipped `employee.EmployeeEducation` pattern exactly (same dual FK + free-text-fallback shape, same response-flattening approach).

**Tech Stack:** Go, GORM, Gin, MySQL + PostgreSQL dual migrations, existing `backend/internal/modules/recruitment` module, cross-module read-only relation to `backend/internal/modules/setting`.

## Global Constraints

- New migration number is **098** (last existing is `097_recruitment_stage_history`).
- Migrations must be idempotent, dual-dialect (mysql+postgres), with matching `.down.sql`, following the exact style of `097_recruitment_stage_history.sql` (CREATE TABLE IF NOT EXISTS / ADD COLUMN IF NOT EXISTS pattern already used in this migration series).
- `candidate_educations.education_id` → FK to `educations.id` (table already exists, created by an earlier platform/tenant migration — do NOT recreate it). `candidate_educations.education_major_id` → FK to `education_majors.id` (same — already exists).
- `candidate_number` format: `CAND-YYYYMM-XXXXXXXX`, generated via a `generateCandidateNumber()` helper mirroring `generateRequisitionNumber()` (`service.go:1938`) and `generateOfferNumber()` (`service.go:2042`) exactly — same `fmt.Sprintf` shape, just prefix `CAND` and 8 uppercased hex chars from a fresh UUID.
- No FE changes in this plan (backend only).
- `candidate_work_experiences` has no master-table reference (no employer/job-title master exists in this codebase) — plain free-text fields.
- Follow this module's existing response-DTO flattening pattern for the education relation, copied from `employee.dto.go:385-406` (`toEducationResponse`) — read that function before writing the recruitment equivalent, do not redesign the shape.

---

## Task 1: Migration 098 — `candidate_number` column + 2 new tables

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.down.sql`

**Interfaces:**
- Produces: `candidates.candidate_number` column; tables `candidate_educations (id, candidate_id, education_id, institution_name, education_major_id, major, gpa, start_year, end_year, is_highest, notes, created_at, updated_at)` and `candidate_work_experiences (id, candidate_id, company_name, job_title, employment_type, start_date, end_date, is_current, description, created_at, updated_at)` that Task 2's GORM models must match exactly.

- [ ] **Step 1: Confirm the `educations`/`education_majors` tables' actual column types before writing FKs**

Run: `grep -n "CREATE TABLE" -A 15 backend/internal/pkg/migrator/migrations/**/*.sql | grep -B2 -A15 "TABLE.*educations\b"` (or search the platform/tenant migration that created `educations`/`education_majors` — likely an early-numbered migration) to confirm the `id` column type (expected `CHAR(36)`, matching every other master table in this codebase, but verify before assuming).

- [ ] **Step 2: Write postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 098_candidate_profile_basics (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 1: candidate_number + candidate_educations +
-- candidate_work_experiences
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-profile-basics-design.md)
--
-- candidates.candidate_number VARCHAR(50) NULL — auto-generated
-- CAND-YYYYMM-XXXXXXXX (pola requisition_number/offer_number).
--
-- candidate_educations — riwayat pendidikan kandidat, education_id/
-- education_major_id merujuk ke master setting.Education/EducationMajor
-- (pola employee_educations existing) dengan fallback teks bebas.
--
-- candidate_work_experiences — riwayat pekerjaan kandidat, tanpa master
-- (tidak ada master employer/job-title di sistem).
--
-- Idempotent: ADD COLUMN / CREATE TABLE IF NOT EXISTS.

ALTER TABLE candidates
    ADD COLUMN IF NOT EXISTS candidate_number VARCHAR(50) NULL;

CREATE TABLE IF NOT EXISTS candidate_educations (
    id                  CHAR(36) PRIMARY KEY,
    candidate_id        CHAR(36) NOT NULL,
    education_id        CHAR(36) NULL,
    institution_name    VARCHAR(255) NOT NULL,
    education_major_id  CHAR(36) NULL,
    major               VARCHAR(255) NULL,
    gpa                 DECIMAL(3,2) NULL,
    start_year          INT NULL,
    end_year            INT NULL,
    is_highest          BOOLEAN NOT NULL DEFAULT false,
    notes               TEXT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_edu_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_edu_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL,
    CONSTRAINT fk_cand_edu_major FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_cand_edu_candidate ON candidate_educations (candidate_id);

CREATE TABLE IF NOT EXISTS candidate_work_experiences (
    id               CHAR(36) PRIMARY KEY,
    candidate_id     CHAR(36) NOT NULL,
    company_name     VARCHAR(255) NOT NULL,
    job_title        VARCHAR(255) NOT NULL,
    employment_type  VARCHAR(50) NULL,
    start_date       DATE NOT NULL,
    end_date         DATE NULL,
    is_current       BOOLEAN NOT NULL DEFAULT false,
    description      TEXT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_exp_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_exp_candidate ON candidate_work_experiences (candidate_id);
```

If Step 1 found `educations`/`education_majors` use a different `id` type than `CHAR(36)`, adjust `education_id`/`education_major_id` to match exactly before proceeding — a type mismatch breaks the FK constraint.

- [ ] **Step 3: Write postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 098_candidate_profile_basics (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_work_experiences;
DROP TABLE IF EXISTS candidate_educations;
ALTER TABLE candidates DROP COLUMN IF EXISTS candidate_number;
```

- [ ] **Step 4: Write mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.sql` — same schema as postgres, MySQL idempotency idiom for the `ALTER TABLE ... ADD COLUMN` (mirror the exact `information_schema.COLUMNS` + `PREPARE`/`EXECUTE` pattern from `096_recruitment_employee_handoff.sql` mysql version), plain `CREATE TABLE IF NOT EXISTS` + `CREATE INDEX` (no `IF NOT EXISTS` on index, matching `095`/`097` mysql convention) for the two new tables:

```sql
-- =============================================================================
-- Tenant Migration: 098_candidate_profile_basics (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: ALTER via information_schema + PREPARE/EXECUTE; CREATE TABLE IF NOT EXISTS.

SET @add_candidate_number = IF(
  EXISTS(
    SELECT 1 FROM information_schema.COLUMNS
    WHERE TABLE_SCHEMA = DATABASE()
      AND TABLE_NAME = 'candidates'
      AND COLUMN_NAME = 'candidate_number'
  ),
  'DO 0',
  'ALTER TABLE candidates ADD COLUMN candidate_number VARCHAR(50) NULL'
);
PREPARE stmt FROM @add_candidate_number;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS candidate_educations (
    id                  CHAR(36) PRIMARY KEY,
    candidate_id        CHAR(36) NOT NULL,
    education_id        CHAR(36) NULL,
    institution_name    VARCHAR(255) NOT NULL,
    education_major_id  CHAR(36) NULL,
    major               VARCHAR(255) NULL,
    gpa                 DECIMAL(3,2) NULL,
    start_year          INT NULL,
    end_year            INT NULL,
    is_highest          BOOLEAN NOT NULL DEFAULT false,
    notes               TEXT NULL,
    created_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_edu_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_edu_education FOREIGN KEY (education_id) REFERENCES educations(id) ON DELETE SET NULL,
    CONSTRAINT fk_cand_edu_major FOREIGN KEY (education_major_id) REFERENCES education_majors(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_edu_candidate ON candidate_educations (candidate_id);

CREATE TABLE IF NOT EXISTS candidate_work_experiences (
    id               CHAR(36) PRIMARY KEY,
    candidate_id     CHAR(36) NOT NULL,
    company_name     VARCHAR(255) NOT NULL,
    job_title        VARCHAR(255) NOT NULL,
    employment_type  VARCHAR(50) NULL,
    start_date       DATE NOT NULL,
    end_date         DATE NULL,
    is_current       BOOLEAN NOT NULL DEFAULT false,
    description      TEXT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_exp_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_exp_candidate ON candidate_work_experiences (candidate_id);
```

- [ ] **Step 5: Write mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 098_candidate_profile_basics (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_work_experiences;
DROP TABLE IF EXISTS candidate_educations;
ALTER TABLE candidates DROP COLUMN candidate_number;
```

(MySQL `DROP COLUMN` has no `IF EXISTS` in older MySQL versions used by this repo's other down migrations — check `096_recruitment_employee_handoff.down.sql` mysql version for the exact idiom this repo uses and mirror it; if that file also drops columns unconditionally, match it.)

- [ ] **Step 6: Verify no migration number collision, then commit**

Run: `ls backend/internal/pkg/migrator/migrations/tenant/postgres/ | grep "^098"` and same for mysql — expect exactly the 2 files each just created.

```bash
git add backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/098_candidate_profile_basics.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/098_candidate_profile_basics.down.sql
git commit -m "feat: migration 098 candidate_number + candidate_educations + candidate_work_experiences (G-6)"
```

---

## Task 2: GORM Models — `CandidateEducation`, `CandidateWorkExperience`, `Candidate.CandidateNumber`

**Files:**
- Modify: `backend/internal/modules/recruitment/model.go`

**Interfaces:**
- Consumes: `setting.Education`, `setting.EducationMajor` (read-only relation types, already exist in `backend/internal/modules/setting/model.go`).
- Produces: `CandidateEducation` struct, `CandidateWorkExperience` struct, `Candidate.CandidateNumber` field — Task 3 (repository) and Task 4 (service) depend on these exact field names.

- [ ] **Step 1: Add the setting-module import**

In `model.go`'s import block, add `"github.com/inthros/hris-platform/internal/modules/setting"` (mirror exactly how `backend/internal/modules/employee/model.go` imports it).

- [ ] **Step 2: Add `CandidateNumber` to the `Candidate` struct**

In the existing `Candidate` struct (`model.go:208-229`), add one field (placed near `Source`, before `Notes` — matches the column's logical grouping):

```go
	CandidateNumber string         `gorm:"type:varchar(50)" json:"candidate_number,omitempty"`
```

- [ ] **Step 3: Add `CandidateEducation` and `CandidateWorkExperience` structs**

Append after the `Candidate` block (after its `BeforeCreate`, before `JobApplication`), mirroring `employee.EmployeeEducation`/`EmployeeExperience` field naming exactly where the concepts overlap:

```go
// =========================================================================
// CandidateEducation (G-6 — riwayat pendidikan kandidat)
// =========================================================================
// education_id/education_major_id merujuk ke master setting.Education/
// EducationMajor (pola sama employee.EmployeeEducation) — nullable dengan
// fallback teks bebas (Major) untuk nilai yang belum ada di master.

type CandidateEducation struct {
	ID               uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID      uuid.UUID  `gorm:"type:char(36);not null;index:idx_cand_edu_candidate" json:"candidate_id"`
	EducationID      *uuid.UUID `gorm:"type:char(36)" json:"education_id,omitempty"`
	InstitutionName  string     `gorm:"type:varchar(255);not null" json:"institution_name"`
	EducationMajorID *uuid.UUID `gorm:"type:char(36)" json:"education_major_id,omitempty"`
	Major            *string    `gorm:"type:varchar(255)" json:"major,omitempty"`
	GPA              *float64   `gorm:"type:decimal(3,2)" json:"gpa,omitempty"`
	StartYear        *int       `gorm:"type:int" json:"start_year,omitempty"`
	EndYear          *int       `gorm:"type:int" json:"end_year,omitempty"`
	IsHighest        bool       `gorm:"type:boolean;default:false" json:"is_highest"`
	Notes            string     `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// Relasi read-only (settings module) — dipakai untuk expand MajorName di
	// response, tidak diserialisasi langsung (pola employee.EmployeeEducation).
	Education      *setting.Education      `gorm:"foreignKey:EducationID" json:"-"`
	EducationMajor *setting.EducationMajor `gorm:"foreignKey:EducationMajorID" json:"-"`
}

func (CandidateEducation) TableName() string {
	return "candidate_educations"
}

func (e *CandidateEducation) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// CandidateWorkExperience (G-6 — riwayat pekerjaan kandidat)
// =========================================================================
// Tanpa master employer/job-title (tidak ada di sistem) — field teks bebas.

type CandidateWorkExperience struct {
	ID             uuid.UUID  `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID    uuid.UUID  `gorm:"type:char(36);not null;index:idx_cand_exp_candidate" json:"candidate_id"`
	CompanyName    string     `gorm:"type:varchar(255);not null" json:"company_name"`
	JobTitle       string     `gorm:"type:varchar(255);not null" json:"job_title"`
	EmploymentType *string    `gorm:"type:varchar(50)" json:"employment_type,omitempty"`
	StartDate      string     `gorm:"type:date;not null" json:"start_date"`
	EndDate        *string    `gorm:"type:date" json:"end_date,omitempty"`
	IsCurrent      bool       `gorm:"type:boolean;default:false" json:"is_current"`
	Description    string     `gorm:"type:text" json:"description,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (CandidateWorkExperience) TableName() string {
	return "candidate_work_experiences"
}

func (e *CandidateWorkExperience) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}
```

Note: `StartDate`/`EndDate` are modeled as `string` (not `time.Time`) with GORM type `date`, matching this module's existing convention for date-only fields (see `JobRequisition.TargetStartDate`/`JobOffer.StartDate`/`ExpiryDate` — all `string` with `gorm:"type:varchar(10)"` or similar date-as-string pattern). Read one of those existing fields first to confirm the exact convention before deciding final type — do not introduce `time.Time` for dates if the rest of the module avoids it.

- [ ] **Step 4: Compile check**

Run: `cd backend && go build ./internal/modules/recruitment/...`
Expected: no errors. If `setting` package import creates a circular dependency (it shouldn't — `employee` already imports it the same way, and `recruitment` doesn't appear in `setting`'s own imports), this would surface here.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/recruitment/model.go
git commit -m "feat: model CandidateEducation + CandidateWorkExperience + Candidate.CandidateNumber (G-6)"
```

---

## Task 3: Repository — CRUD for both entities, with `EducationMajor` preload

**Files:**
- Modify: `backend/internal/modules/recruitment/repository.go`
- Test: `backend/internal/modules/recruitment/repository_test.go`

**Interfaces:**
- Consumes: `CandidateEducation`, `CandidateWorkExperience` (Task 2).
- Produces (used by Task 4 service):
  - `func (r *Repository) CreateCandidateEducation(ctx, e *CandidateEducation) error`
  - `func (r *Repository) FindCandidateEducationByID(ctx, id uuid.UUID) (*CandidateEducation, error)` — must `Preload("EducationMajor")` so the response can flatten `MajorName`.
  - `func (r *Repository) ListCandidateEducations(ctx, candidateID uuid.UUID) ([]CandidateEducation, error)` — must `Preload("EducationMajor")`.
  - `func (r *Repository) UpdateCandidateEducation(ctx, e *CandidateEducation) error`
  - `func (r *Repository) DeleteCandidateEducation(ctx, id uuid.UUID) error`
  - Same 5-method shape for `CandidateWorkExperience` (no preload needed — no relation).

- [ ] **Step 1: Write failing repository tests**

Append to `repository_test.go`, following the exact structure of the G-5 work's `TestRepository_CreateAndFindStageByCode` (create dependency chain: requisition → candidate → education/experience row → find/list):

```go
func TestRepository_CreateAndFindCandidateEducation(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Edu", LastName: "Test", Email: "edutest@test.com"}
	repo.CreateCandidate(ctx, cand)

	edu := &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Universitas Test"}
	if err := repo.CreateCandidateEducation(ctx, edu); err != nil {
		t.Fatalf("CreateCandidateEducation failed: %v", err)
	}

	found, err := repo.FindCandidateEducationByID(ctx, edu.ID)
	if err != nil {
		t.Fatalf("FindCandidateEducationByID failed: %v", err)
	}
	if found.InstitutionName != "Universitas Test" {
		t.Errorf("expected institution 'Universitas Test', got %s", found.InstitutionName)
	}
}

func TestRepository_ListCandidateEducations(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Edu", Email: "listedu@test.com"}
	repo.CreateCandidate(ctx, cand)
	repo.CreateCandidateEducation(ctx, &CandidateEducation{CandidateID: cand.ID, InstitutionName: "SMA 1"})
	repo.CreateCandidateEducation(ctx, &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Universitas A"})

	list, err := repo.ListCandidateEducations(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateEducations failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 education rows, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateEducation(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Edu", Email: "updedu@test.com"}
	repo.CreateCandidate(ctx, cand)
	edu := &CandidateEducation{CandidateID: cand.ID, InstitutionName: "Original"}
	repo.CreateCandidateEducation(ctx, edu)

	edu.InstitutionName = "Updated"
	if err := repo.UpdateCandidateEducation(ctx, edu); err != nil {
		t.Fatalf("UpdateCandidateEducation failed: %v", err)
	}
	found, _ := repo.FindCandidateEducationByID(ctx, edu.ID)
	if found.InstitutionName != "Updated" {
		t.Errorf("expected 'Updated', got %s", found.InstitutionName)
	}

	if err := repo.DeleteCandidateEducation(ctx, edu.ID); err != nil {
		t.Fatalf("DeleteCandidateEducation failed: %v", err)
	}
	if _, err := repo.FindCandidateEducationByID(ctx, edu.ID); err == nil {
		t.Error("expected error finding deleted education, got nil")
	}
}

func TestRepository_CreateAndListCandidateWorkExperience(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Exp", LastName: "Test", Email: "exptest@test.com"}
	repo.CreateCandidate(ctx, cand)

	exp := &CandidateWorkExperience{CandidateID: cand.ID, CompanyName: "Acme Corp", JobTitle: "Engineer", StartDate: "2020-01-01"}
	if err := repo.CreateCandidateWorkExperience(ctx, exp); err != nil {
		t.Fatalf("CreateCandidateWorkExperience failed: %v", err)
	}

	list, err := repo.ListCandidateWorkExperiences(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateWorkExperiences failed: %v", err)
	}
	if len(list) != 1 || list[0].CompanyName != "Acme Corp" {
		t.Errorf("expected 1 row with company 'Acme Corp', got %+v", list)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestRepository_CreateAndFindCandidateEducation -v`
Expected: compile error (methods undefined).

- [ ] **Step 3: Implement repository methods**

Append to `repository.go` (after the Candidates section, before Job Applications — or at end of file, matching existing section-comment style):

```go
// =========================================================================
// Candidate Educations (G-6)
// =========================================================================

func (r *Repository) CreateCandidateEducation(ctx context.Context, e *CandidateEducation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindCandidateEducationByID(ctx context.Context, id uuid.UUID) (*CandidateEducation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e CandidateEducation
	if err := db.WithContext(ctx).Preload("EducationMajor").First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate education not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListCandidateEducations(ctx context.Context, candidateID uuid.UUID) ([]CandidateEducation, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateEducation
	if err := db.WithContext(ctx).Preload("EducationMajor").Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateEducation(ctx context.Context, e *CandidateEducation) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(e).Error
}

func (r *Repository) DeleteCandidateEducation(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateEducation{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate education not found")
	}
	return result.Error
}

// =========================================================================
// Candidate Work Experiences (G-6)
// =========================================================================

func (r *Repository) CreateCandidateWorkExperience(ctx context.Context, e *CandidateWorkExperience) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(e).Error
}

func (r *Repository) FindCandidateWorkExperienceByID(ctx context.Context, id uuid.UUID) (*CandidateWorkExperience, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var e CandidateWorkExperience
	if err := db.WithContext(ctx).First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate work experience not found")
		}
		return nil, err
	}
	return &e, nil
}

func (r *Repository) ListCandidateWorkExperiences(ctx context.Context, candidateID uuid.UUID) ([]CandidateWorkExperience, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateWorkExperience
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("start_date DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateWorkExperience(ctx context.Context, e *CandidateWorkExperience) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(e).Error
}

func (r *Repository) DeleteCandidateWorkExperience(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateWorkExperience{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate work experience not found")
	}
	return result.Error
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestRepository_.*Candidate(Education|WorkExperience)" -v`
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/recruitment/repository.go backend/internal/modules/recruitment/repository_test.go
git commit -m "feat: repository CRUD CandidateEducation + CandidateWorkExperience (G-6)"
```

---

## Task 4: Service — CRUD, `generateCandidateNumber`, wire into `CreateCandidate`

**Files:**
- Modify: `backend/internal/modules/recruitment/service.go`
- Test: `backend/internal/modules/recruitment/service_test.go`

**Interfaces:**
- Consumes: repository methods from Task 3.
- Produces:
  - `func generateCandidateNumber() string`
  - `func (s *Service) CreateCandidateEducation(ctx, candidateID string, req CreateCandidateEducationRequest) (*CandidateEducationResponse, error)` (validates candidate exists → 404-equivalent error if not)
  - `func (s *Service) ListCandidateEducations(ctx, candidateID string) ([]CandidateEducationResponse, error)`
  - `func (s *Service) UpdateCandidateEducation(ctx, id string, req UpdateCandidateEducationRequest) (*CandidateEducationResponse, error)`
  - `func (s *Service) DeleteCandidateEducation(ctx, id string) error`
  - Same 4-method shape for `CandidateWorkExperience`.
  - `CreateCandidate` sets `c.CandidateNumber` via `generateCandidateNumber()` unless `req.CandidateNumber` is explicitly provided.

- [ ] **Step 1: Add DTOs to `dto.go`**

```go
// =========================================================================
// Candidate Education DTOs (G-6)
// =========================================================================

type CreateCandidateEducationRequest struct {
	EducationID      *string  `json:"education_id" binding:"omitempty"`
	InstitutionName  string   `json:"institution_name" binding:"required,max=255"`
	EducationMajorID *string  `json:"education_major_id" binding:"omitempty"`
	Major            *string  `json:"major" binding:"omitempty,max=255"`
	GPA              *float64 `json:"gpa"`
	StartYear        *int     `json:"start_year"`
	EndYear          *int     `json:"end_year"`
	IsHighest        bool     `json:"is_highest"`
	Notes            string   `json:"notes"`
}

type UpdateCandidateEducationRequest struct {
	EducationID      *string  `json:"education_id"`
	InstitutionName  *string  `json:"institution_name" binding:"omitempty,max=255"`
	EducationMajorID *string  `json:"education_major_id"`
	Major            *string  `json:"major" binding:"omitempty,max=255"`
	GPA              *float64 `json:"gpa"`
	StartYear        *int     `json:"start_year"`
	EndYear          *int     `json:"end_year"`
	IsHighest        *bool    `json:"is_highest"`
	Notes            *string  `json:"notes"`
}

type CandidateEducationResponse struct {
	ID               string  `json:"id"`
	CandidateID      string  `json:"candidate_id"`
	EducationID      string  `json:"education_id,omitempty"`
	EducationMajorID string  `json:"education_major_id,omitempty"`
	MajorName        string  `json:"major_name,omitempty"`
	InstitutionName  string  `json:"institution_name"`
	Major            string  `json:"major,omitempty"`
	GPA              float64 `json:"gpa,omitempty"`
	StartYear        int     `json:"start_year,omitempty"`
	EndYear          int     `json:"end_year,omitempty"`
	IsHighest        bool    `json:"is_highest"`
	Notes            string  `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Work Experience DTOs (G-6)
// =========================================================================

type CreateCandidateWorkExperienceRequest struct {
	CompanyName    string  `json:"company_name" binding:"required,max=255"`
	JobTitle       string  `json:"job_title" binding:"required,max=255"`
	EmploymentType *string `json:"employment_type" binding:"omitempty,max=50"`
	StartDate      string  `json:"start_date" binding:"required,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	IsCurrent      bool    `json:"is_current"`
	Description    string  `json:"description"`
}

type UpdateCandidateWorkExperienceRequest struct {
	CompanyName    *string `json:"company_name" binding:"omitempty,max=255"`
	JobTitle       *string `json:"job_title" binding:"omitempty,max=255"`
	EmploymentType *string `json:"employment_type" binding:"omitempty,max=50"`
	StartDate      *string `json:"start_date" binding:"omitempty,max=10"`
	EndDate        *string `json:"end_date" binding:"omitempty,max=10"`
	IsCurrent      *bool   `json:"is_current"`
	Description    *string `json:"description"`
}

type CandidateWorkExperienceResponse struct {
	ID             string `json:"id"`
	CandidateID    string `json:"candidate_id"`
	CompanyName    string `json:"company_name"`
	JobTitle       string `json:"job_title"`
	EmploymentType string `json:"employment_type,omitempty"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date,omitempty"`
	IsCurrent      bool   `json:"is_current"`
	Description    string `json:"description,omitempty"`
}
```

Also add `CandidateNumber string \`json:"candidate_number,omitempty"\`` to `CandidateResponse` (`dto.go:182-201`), and `CandidateNumber *string \`json:"candidate_number" binding:"omitempty,max=50"\`` to `CreateCandidateRequest`/`UpdateCandidateRequest`.

- [ ] **Step 2: Write failing service tests**

Append to `service_test.go`:

```go
func TestService_GenerateCandidateNumber_Format(t *testing.T) {
	num := generateCandidateNumber()
	if !strings.HasPrefix(num, "CAND-") {
		t.Errorf("expected prefix CAND-, got %s", num)
	}
	if len(num) != len("CAND-200601-XXXXXXXX") {
		t.Errorf("expected length %d, got %d (%s)", len("CAND-200601-XXXXXXXX"), len(num), num)
	}
}

func TestService_CreateCandidate_AutoGeneratesCandidateNumber(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, err := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Auto", LastName: "Num", Email: "autonum@test.com"})
	if err != nil {
		t.Fatalf("CreateCandidate failed: %v", err)
	}
	if !strings.HasPrefix(cand.CandidateNumber, "CAND-") {
		t.Errorf("expected auto-generated candidate_number, got %q", cand.CandidateNumber)
	}
}

func TestService_CreateCandidateEducation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Edu", LastName: "Svc", Email: "edusvc@test.com"})

	resp, err := svc.CreateCandidateEducation(ctx, cand.ID, CreateCandidateEducationRequest{
		InstitutionName: "Universitas Indonesia",
	})
	if err != nil {
		t.Fatalf("CreateCandidateEducation failed: %v", err)
	}
	if resp.InstitutionName != "Universitas Indonesia" {
		t.Errorf("expected institution 'Universitas Indonesia', got %s", resp.InstitutionName)
	}
}

func TestService_CreateCandidateEducation_UnknownCandidate(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	_, err := svc.CreateCandidateEducation(ctx, uuid.New().String(), CreateCandidateEducationRequest{InstitutionName: "X"})
	if err == nil {
		t.Fatal("expected error for unknown candidate, got nil")
	}
}

func TestService_ListCandidateEducations(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "List", LastName: "Svc", Email: "listsvc@test.com"})
	svc.CreateCandidateEducation(ctx, cand.ID, CreateCandidateEducationRequest{InstitutionName: "SMA 1"})
	svc.CreateCandidateEducation(ctx, cand.ID, CreateCandidateEducationRequest{InstitutionName: "Universitas A"})

	list, err := svc.ListCandidateEducations(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateEducations failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
}

func TestService_UpdateAndDeleteCandidateEducation(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Upd", LastName: "Svc", Email: "updsvc@test.com"})
	created, _ := svc.CreateCandidateEducation(ctx, cand.ID, CreateCandidateEducationRequest{InstitutionName: "Original"})

	newName := "Updated"
	updated, err := svc.UpdateCandidateEducation(ctx, created.ID, UpdateCandidateEducationRequest{InstitutionName: &newName})
	if err != nil {
		t.Fatalf("UpdateCandidateEducation failed: %v", err)
	}
	if updated.InstitutionName != "Updated" {
		t.Errorf("expected 'Updated', got %s", updated.InstitutionName)
	}

	if err := svc.DeleteCandidateEducation(ctx, created.ID); err != nil {
		t.Fatalf("DeleteCandidateEducation failed: %v", err)
	}
}

func TestService_CreateCandidateWorkExperience(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Exp", LastName: "Svc", Email: "expsvc@test.com"})

	resp, err := svc.CreateCandidateWorkExperience(ctx, cand.ID, CreateCandidateWorkExperienceRequest{
		CompanyName: "Acme", JobTitle: "Engineer", StartDate: "2020-01-01",
	})
	if err != nil {
		t.Fatalf("CreateCandidateWorkExperience failed: %v", err)
	}
	if resp.CompanyName != "Acme" {
		t.Errorf("expected company 'Acme', got %s", resp.CompanyName)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestService_GenerateCandidateNumber_Format -v`
Expected: compile error.

- [ ] **Step 4: Implement `generateCandidateNumber` and wire into `CreateCandidate`**

Add near `generateOfferNumber` (`service.go:2042`):

```go
// generateCandidateNumber (G-6) membuat nomor kandidat otomatis dengan
// format CAND-YYYYMM-XXXXXXXX (pola sama generateRequisitionNumber G-2 /
// generateOfferNumber G-3).
func generateCandidateNumber() string {
	return fmt.Sprintf("CAND-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))
}
```

In `CreateCandidate` (`service.go:952-1000`), after the `c := &Candidate{...}` literal and before `applyCandidateTypeFields`, add:

```go
	if req.CandidateNumber != nil && *req.CandidateNumber != "" {
		c.CandidateNumber = *req.CandidateNumber
	} else {
		c.CandidateNumber = generateCandidateNumber()
	}
```

Update `candidateToResponse` (search for it in `service.go`, near `applicationToResponse`) to include `CandidateNumber: c.CandidateNumber`.

- [ ] **Step 5: Implement CRUD service methods for both entities**

Add near `CreateCandidate`/`GetCandidateByID` (or in a new section after the existing Candidates block):

```go
// =========================================================================
// Candidate Educations (G-6)
// =========================================================================

func (s *Service) CreateCandidateEducation(ctx context.Context, candidateID string, req CreateCandidateEducationRequest) (*CandidateEducationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}

	e := &CandidateEducation{
		CandidateID:     candUUID,
		InstitutionName: req.InstitutionName,
		GPA:             req.GPA,
		StartYear:       req.StartYear,
		EndYear:         req.EndYear,
		IsHighest:       req.IsHighest,
		Notes:           req.Notes,
	}
	if req.EducationID != nil && *req.EducationID != "" {
		eduID, err := uuid.Parse(*req.EducationID)
		if err != nil {
			return nil, fmt.Errorf("invalid education_id: %w", err)
		}
		e.EducationID = &eduID
	}
	if req.EducationMajorID != nil && *req.EducationMajorID != "" {
		majorID, err := uuid.Parse(*req.EducationMajorID)
		if err != nil {
			return nil, fmt.Errorf("invalid education_major_id: %w", err)
		}
		e.EducationMajorID = &majorID
	}
	if req.Major != nil {
		e.Major = req.Major
	}

	if err := s.repo.CreateCandidateEducation(ctx, e); err != nil {
		return nil, err
	}
	created, err := s.repo.FindCandidateEducationByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	return candidateEducationToResponse(created), nil
}

func (s *Service) ListCandidateEducations(ctx context.Context, candidateID string) ([]CandidateEducationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateEducations(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateEducationResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateEducationToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateEducation(ctx context.Context, id string, req UpdateCandidateEducationRequest) (*CandidateEducationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindCandidateEducationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.EducationID != nil {
		if *req.EducationID == "" {
			e.EducationID = nil
		} else {
			eduID, err := uuid.Parse(*req.EducationID)
			if err != nil {
				return nil, fmt.Errorf("invalid education_id: %w", err)
			}
			e.EducationID = &eduID
		}
	}
	if req.InstitutionName != nil {
		e.InstitutionName = *req.InstitutionName
	}
	if req.EducationMajorID != nil {
		if *req.EducationMajorID == "" {
			e.EducationMajorID = nil
		} else {
			majorID, err := uuid.Parse(*req.EducationMajorID)
			if err != nil {
				return nil, fmt.Errorf("invalid education_major_id: %w", err)
			}
			e.EducationMajorID = &majorID
		}
	}
	if req.Major != nil {
		e.Major = req.Major
	}
	if req.GPA != nil {
		e.GPA = req.GPA
	}
	if req.StartYear != nil {
		e.StartYear = req.StartYear
	}
	if req.EndYear != nil {
		e.EndYear = req.EndYear
	}
	if req.IsHighest != nil {
		e.IsHighest = *req.IsHighest
	}
	if req.Notes != nil {
		e.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateEducation(ctx, e); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindCandidateEducationByID(ctx, e.ID)
	if err != nil {
		return nil, err
	}
	return candidateEducationToResponse(updated), nil
}

func (s *Service) DeleteCandidateEducation(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateEducation(ctx, uid)
}

func candidateEducationToResponse(e *CandidateEducation) *CandidateEducationResponse {
	resp := &CandidateEducationResponse{
		ID:              e.ID.String(),
		CandidateID:     e.CandidateID.String(),
		InstitutionName: e.InstitutionName,
		IsHighest:       e.IsHighest,
		Notes:           e.Notes,
	}
	if e.EducationID != nil {
		resp.EducationID = e.EducationID.String()
	}
	if e.EducationMajorID != nil {
		resp.EducationMajorID = e.EducationMajorID.String()
	}
	if e.EducationMajor != nil {
		resp.MajorName = e.EducationMajor.Name
	}
	if e.Major != nil {
		resp.Major = *e.Major
	}
	if e.GPA != nil {
		resp.GPA = *e.GPA
	}
	if e.StartYear != nil {
		resp.StartYear = *e.StartYear
	}
	if e.EndYear != nil {
		resp.EndYear = *e.EndYear
	}
	return resp
}

// =========================================================================
// Candidate Work Experiences (G-6)
// =========================================================================

func (s *Service) CreateCandidateWorkExperience(ctx context.Context, candidateID string, req CreateCandidateWorkExperienceRequest) (*CandidateWorkExperienceResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	e := &CandidateWorkExperience{
		CandidateID:    candUUID,
		CompanyName:    req.CompanyName,
		JobTitle:       req.JobTitle,
		EmploymentType: req.EmploymentType,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsCurrent:      req.IsCurrent,
		Description:    req.Description,
	}
	if err := s.repo.CreateCandidateWorkExperience(ctx, e); err != nil {
		return nil, err
	}
	return candidateWorkExperienceToResponse(e), nil
}

func (s *Service) ListCandidateWorkExperiences(ctx context.Context, candidateID string) ([]CandidateWorkExperienceResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateWorkExperiences(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateWorkExperienceResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateWorkExperienceToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateWorkExperience(ctx context.Context, id string, req UpdateCandidateWorkExperienceRequest) (*CandidateWorkExperienceResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	e, err := s.repo.FindCandidateWorkExperienceByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.CompanyName != nil {
		e.CompanyName = *req.CompanyName
	}
	if req.JobTitle != nil {
		e.JobTitle = *req.JobTitle
	}
	if req.EmploymentType != nil {
		e.EmploymentType = req.EmploymentType
	}
	if req.StartDate != nil {
		e.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		e.EndDate = req.EndDate
	}
	if req.IsCurrent != nil {
		e.IsCurrent = *req.IsCurrent
	}
	if req.Description != nil {
		e.Description = *req.Description
	}
	if err := s.repo.UpdateCandidateWorkExperience(ctx, e); err != nil {
		return nil, err
	}
	return candidateWorkExperienceToResponse(e), nil
}

func (s *Service) DeleteCandidateWorkExperience(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateWorkExperience(ctx, uid)
}

func candidateWorkExperienceToResponse(e *CandidateWorkExperience) *CandidateWorkExperienceResponse {
	resp := &CandidateWorkExperienceResponse{
		ID:          e.ID.String(),
		CandidateID: e.CandidateID.String(),
		CompanyName: e.CompanyName,
		JobTitle:    e.JobTitle,
		StartDate:   e.StartDate,
		IsCurrent:   e.IsCurrent,
		Description: e.Description,
	}
	if e.EmploymentType != nil {
		resp.EmploymentType = *e.EmploymentType
	}
	if e.EndDate != nil {
		resp.EndDate = *e.EndDate
	}
	return resp
}
```

- [ ] **Step 6: Run all recruitment service tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v 2>&1 | tail -100`
Expected: all pass, including every pre-existing test (no regressions from the `CandidateResponse`/`CreateCandidateRequest` field additions — check `candidateToResponse` update didn't break any test asserting exact response shape via struct comparison rather than field-by-field checks; if any does, that test needs updating to account for the new field, not the feature scaled back).

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/recruitment/service.go backend/internal/modules/recruitment/service_test.go backend/internal/modules/recruitment/dto.go
git commit -m "feat: service CRUD CandidateEducation + CandidateWorkExperience + candidate_number generation (G-6)"
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
  POST   /recruitment/candidates/:candidate_id/educations
  GET    /recruitment/candidates/:candidate_id/educations
  PUT    /recruitment/educations/:id
  DELETE /recruitment/educations/:id
  POST   /recruitment/candidates/:candidate_id/work-experiences
  GET    /recruitment/candidates/:candidate_id/work-experiences
  PUT    /recruitment/work-experiences/:id
  DELETE /recruitment/work-experiences/:id
  ```

- [ ] **Step 1: Write failing handler tests**

Append to `handler_test.go`, following the exact request/response-parsing style already used throughout that file (see `TestHandler_CreateInterview` for the create→assert-201 pattern, `TestHandler_ListOnboardingTaskItems` for the create-parent→create-child→list pattern):

```go
func TestHandler_CreateCandidateEducation(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Edu", LastName: "Handler", Email: "eduhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{
		InstitutionName: "Universitas Test",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ListCandidateEducations(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "List", LastName: "EduH", Email: "listeduh@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{InstitutionName: "SMA 1"})

	w := performRequest(r, "GET", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateAndDeleteCandidateEducation(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Upd", LastName: "EduH", Email: "updeduh@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	eW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/educations", CreateCandidateEducationRequest{InstitutionName: "Original"})
	var eResp map[string]interface{}
	json.Unmarshal(eW.Body.Bytes(), &eResp)
	eid := eResp["data"].(map[string]interface{})["id"].(string)

	newName := "Updated"
	w := performRequest(r, "PUT", "/api/v1/tenant/recruitment/educations/"+eid, UpdateCandidateEducationRequest{InstitutionName: &newName})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	w2 := performRequest(r, "DELETE", "/api/v1/tenant/recruitment/educations/"+eid, nil)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

func TestHandler_CreateCandidateWorkExperience(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Exp", LastName: "Handler", Email: "exphandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/work-experiences", CreateCandidateWorkExperienceRequest{
		CompanyName: "Acme", JobTitle: "Engineer", StartDate: "2020-01-01",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
```

- [ ] **Step 2: Run to verify failure**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestHandler_CreateCandidateEducation -v`
Expected: 404 (route not registered) or compile error.

- [ ] **Step 3: Implement handlers**

Add to `handler.go` near the existing Candidates handlers:

```go
func (h *Handler) CreateCandidateEducation(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	var req CreateCandidateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateEducation(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateEducations(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	resp, err := h.svc.ListCandidateEducations(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateEducation(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateEducationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateEducation(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateEducation(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateEducation(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateWorkExperience(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	var req CreateCandidateWorkExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateWorkExperience(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateWorkExperiences(c *gin.Context) {
	candidateID := c.Param("candidate_id")
	resp, err := h.svc.ListCandidateWorkExperiences(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateWorkExperience(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateWorkExperienceRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateWorkExperience(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateWorkExperience(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateWorkExperience(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
```

- [ ] **Step 4: Register routes**

In `routes.go`, add after the existing `candidates/:id` routes (`routes.go:32-36`):

```go
	rec.POST("/candidates/:candidate_id/educations", handler.CreateCandidateEducation)
	rec.GET("/candidates/:candidate_id/educations", handler.ListCandidateEducations)
	rec.PUT("/educations/:id", handler.UpdateCandidateEducation)
	rec.DELETE("/educations/:id", handler.DeleteCandidateEducation)
	rec.POST("/candidates/:candidate_id/work-experiences", handler.CreateCandidateWorkExperience)
	rec.GET("/candidates/:candidate_id/work-experiences", handler.ListCandidateWorkExperiences)
	rec.PUT("/work-experiences/:id", handler.UpdateCandidateWorkExperience)
	rec.DELETE("/work-experiences/:id", handler.DeleteCandidateWorkExperience)
```

- [ ] **Step 5: Run handler tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v 2>&1 | tail -60`
Expected: all pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/handler.go backend/internal/modules/recruitment/routes.go backend/internal/modules/recruitment/handler_test.go
git commit -m "feat: handler+routes CandidateEducation/CandidateWorkExperience CRUD (G-6)"
```

---

## Task 6: Module wiring — `AutoMigrate` (test-DB convenience)

**Files:**
- Modify: `backend/internal/modules/recruitment/module.go`
- Modify: `backend/internal/modules/recruitment/helpers_test.go` (if the test-DB setup there needs the new tables — check first; if `setupTestDB`/`newTestService`/`setupTestRouter` already `AutoMigrate` every model dynamically or via a shared list that Task 2's structs get picked up by automatically, no change needed here — verify before assuming a change is required)

**Interfaces:**
- Consumes: `CandidateEducation`, `CandidateWorkExperience` (Task 2).

- [ ] **Step 1: Add both entities to `Migrate()`**

In `module.go`'s `Migrate()` method, add `&CandidateEducation{}` and `&CandidateWorkExperience{}` to the existing `db.AutoMigrate(...)` call (alongside `&RecruitmentStage{}`, `&ApplicationStageHistory{}` from the prior G-5 work).

> Note: per the G-5 work's final review, this `Migrate()`/`Seed()` wiring is **not actually invoked for tenant databases in production** (a pre-existing, separately-tracked bug — `cmd/server/main.go` only calls `.Migrate()`/`.Seed()` for platform modules). This step is still worth doing for consistency with the module's existing pattern and because the in-package unit tests' `setupTestDB` may reuse this exact `AutoMigrate` list (check `helpers_test.go` — if it calls `recModule.Migrate(db)` or duplicates the model list separately, update whichever path the tests actually exercise). The REAL schema source of truth for production remains migration 098 from Task 1.

- [ ] **Step 2: Add the new models to the test-DB AutoMigrate list**

`backend/internal/modules/recruitment/helpers_test.go`'s `setupTestDB()` has its own explicit `db.AutoMigrate(...)` model list (confirmed separate from `module.go`'s — it currently lists `&JobRequisition{}, &Candidate{}, &JobApplication{}, &Interview{}, &JobOffer{}, &OnboardingTaskTemplate{}, &EmployeeOnboarding{}, &OnboardingTaskItem{}, &RecruitmentStage{}, &ApplicationStageHistory{}`). Add `&CandidateEducation{}` and `&CandidateWorkExperience{}` to that same list.

Do **not** add `&setting.Education{}`/`&setting.EducationMajor{}` to this list — confirmed by precedent in `backend/internal/modules/employee/helpers_test.go` (which has the identical `EmployeeEducation.EducationMajor` relation and also does NOT migrate `setting` models in its test DB): GORM's `Preload` skips issuing a query against the related table entirely when the collected foreign-key values are all nil, which is the case for every test in this plan (none of Task 3/4's tests set `EducationMajorID`). Adding the `setting` import here would be unnecessary scope.

- [ ] **Step 3: Full build + test run**

Run: `cd backend && go build ./... && go test ./internal/modules/recruitment/... -v 2>&1 | tail -60`
Expected: clean build, all tests pass.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/recruitment/module.go backend/internal/modules/recruitment/helpers_test.go
git commit -m "feat: AutoMigrate CandidateEducation + CandidateWorkExperience (G-6)"
```

(If Step 2 found no change was needed in `helpers_test.go`, skip staging that file — only commit what actually changed.)

---

## Task 7: Update plan doc status

**Files:**
- Modify: `docs/module-recruitment-development-plan.md`

- [ ] **Step 1: Update §G-6 status**

The G-6 section currently covers the full original scope (4 candidate columns + 5 sub-tables). Update it to reflect partial completion: mark `candidate_number` + `candidate_educations` + `candidate_work_experiences` as ✅ done (migration 098, following the G-1..G-5 status-line style), and explicitly list what's still open for G-6 (candidates.status — skipped as not well-defined, source_id/source master — deferred, candidate_skills/candidate_certifications — deferred pending Skill/Certification master decisions, candidate_documents/candidate_consents — deferred, own compliance concerns). Reference the spec file `docs/superpowers/specs/2026-08-12-candidate-profile-basics-design.md` for the full scoping rationale.

- [ ] **Step 2: Update §8.1/§8.2 API Plan**

Add the 8 new endpoints to §8.1 (existing), incrementing the endpoint count.

- [ ] **Step 3: Commit**

```bash
git add docs/module-recruitment-development-plan.md
git commit -m "docs: update G-6 status - candidate profile basics selesai (sub-project 1/3)"
```
