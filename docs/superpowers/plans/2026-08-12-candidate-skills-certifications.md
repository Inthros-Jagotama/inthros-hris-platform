# Candidate Skills & Certifications (G-6 sub-project 2) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `candidate_skills` (referencing the existing `competency.Competency` master, with a 1-5 level) and `candidate_certifications` (free-text, no master) to the recruitment module.

**Architecture:** Two new GORM entities following the exact same per-entity CRUD pattern established in G-6 sub-project 1 (`CandidateEducation`/`CandidateWorkExperience`): repository → service → handler → routes, wired the same way. `CandidateSkill.Competency` is a read-only cross-module relation to `competency.Competency`, mirroring how `CandidateEducation.Education`/`EducationMajor` already reference the `setting` module.

**Tech Stack:** Go, GORM, Gin, MySQL + PostgreSQL dual migrations, existing `backend/internal/modules/recruitment` module, cross-module read-only relation to `backend/internal/modules/competency`.

## Global Constraints

- New migration number is **099** (last existing is `098_candidate_profile_basics`).
- Migrations must be idempotent, dual-dialect (mysql+postgres), with matching `.down.sql`, following the exact style of `098_candidate_profile_basics.sql`.
- `candidate_skills.competency_id` is **NOT NULL**, FK → `competencies(id)`, `ON DELETE CASCADE` — unlike sub-project 1's nullable education FKs, a skill without a competency link isn't actionable for the G-9 candidate-matching feature planned later, and `ON DELETE CASCADE` is the only valid option for a NOT NULL FK anyway (confirmed as the existing convention: `competency` module's own child tables use `ON DELETE CASCADE` for their NOT NULL `competency_id` FKs in `008_competency.sql`).
- `candidate_certifications` has no master-table reference (none exists in the codebase — `training.TrainingCertificate` is a participant-issued record, not a catalog).
- No FE changes in this plan (backend only).
- Route parameter is `:id` (not `:candidate_id`) for the nested `/candidates/:id/...` routes — same Gin wildcard-name constraint discovered in sub-project 1 (a route can't introduce a differently-named wildcard at a path segment position where `/candidates/:id` already registered `:id`).
- `level` on `candidate_skills` is stored as-is, no range validation (1-5 is a convention, not enforced) — matches the spec's explicit no-validation decision.

---

## Task 1: Migration 099 — `candidate_skills` + `candidate_certifications`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.down.sql`

**Interfaces:**
- Produces: tables `candidate_skills (id, candidate_id, competency_id, level, notes, created_at, updated_at)` and `candidate_certifications (id, candidate_id, name, issuing_organization, issue_date, expiry_date, credential_url, notes, created_at, updated_at)` that Task 2's GORM models must match exactly.

- [ ] **Step 1: Confirm `competencies.id` column type**

Run: `grep -n "CREATE TABLE IF NOT EXISTS competencies" -A 10 backend/internal/pkg/migrator/migrations/tenant/postgres/008_competency.sql` — expected `CHAR(36)` (matches every other master table in this codebase), confirm before writing the FK.

- [ ] **Step 2: Write postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 099_candidate_skills_certifications (PostgreSQL)
-- =============================================================================
-- G-6 sub-project 2: candidate_skills + candidate_certifications
-- (docs/module-recruitment-development-plan.md §G-6;
--  docs/superpowers/specs/2026-08-12-candidate-skills-certifications-design.md)
--
-- candidate_skills — kandidat + skill (referensi competency.Competency
-- master existing, bukan Skill Master baru) + level proficiency 1-5.
-- competency_id NOT NULL (beda dari candidate_educations yang nullable) —
-- skill tanpa referensi competency tidak actionable untuk candidate
-- matching (G-9). ON DELETE CASCADE (satu-satunya opsi valid untuk FK
-- NOT NULL, konsisten dengan child table competency module sendiri).
--
-- candidate_certifications — tanpa master (tidak ada Certification Master
-- di sistem), field bebas + info penerbit/masa berlaku.
--
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_skills (
    id            CHAR(36) PRIMARY KEY,
    candidate_id  CHAR(36) NOT NULL,
    competency_id CHAR(36) NOT NULL,
    level         SMALLINT NULL,
    notes         TEXT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_skill_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_skill_competency FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_skill_candidate ON candidate_skills (candidate_id);
CREATE INDEX IF NOT EXISTS idx_cand_skill_competency ON candidate_skills (competency_id);

CREATE TABLE IF NOT EXISTS candidate_certifications (
    id                   CHAR(36) PRIMARY KEY,
    candidate_id         CHAR(36) NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    issuing_organization VARCHAR(255) NULL,
    issue_date           DATE NULL,
    expiry_date          DATE NULL,
    credential_url       TEXT NULL,
    notes                TEXT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_cert_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_cand_cert_candidate ON candidate_certifications (candidate_id);
```

If Step 1 found `competencies.id` isn't `CHAR(36)`, adjust `competency_id` to match before proceeding.

- [ ] **Step 3: Write postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 099_candidate_skills_certifications (PostgreSQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_certifications;
DROP TABLE IF EXISTS candidate_skills;
```

- [ ] **Step 4: Write mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.sql`:

```sql
-- =============================================================================
-- Tenant Migration: 099_candidate_skills_certifications (MySQL)
-- =============================================================================
-- See postgres version for full column documentation.
-- Idempotent: CREATE TABLE IF NOT EXISTS.

CREATE TABLE IF NOT EXISTS candidate_skills (
    id            CHAR(36) PRIMARY KEY,
    candidate_id  CHAR(36) NOT NULL,
    competency_id CHAR(36) NOT NULL,
    level         SMALLINT NULL,
    notes         TEXT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_skill_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE,
    CONSTRAINT fk_cand_skill_competency FOREIGN KEY (competency_id) REFERENCES competencies(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_skill_candidate ON candidate_skills (candidate_id);
CREATE INDEX idx_cand_skill_competency ON candidate_skills (competency_id);

CREATE TABLE IF NOT EXISTS candidate_certifications (
    id                   CHAR(36) PRIMARY KEY,
    candidate_id         CHAR(36) NOT NULL,
    name                 VARCHAR(255) NOT NULL,
    issuing_organization VARCHAR(255) NULL,
    issue_date           DATE NULL,
    expiry_date          DATE NULL,
    credential_url       TEXT NULL,
    notes                TEXT NULL,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_cand_cert_candidate FOREIGN KEY (candidate_id) REFERENCES candidates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_cand_cert_candidate ON candidate_certifications (candidate_id);
```

- [ ] **Step 5: Write mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.down.sql`:

```sql
-- =============================================================================
-- Tenant Migration Down: 099_candidate_skills_certifications (MySQL)
-- =============================================================================

DROP TABLE IF EXISTS candidate_certifications;
DROP TABLE IF EXISTS candidate_skills;
```

- [ ] **Step 6: Verify no collision, then commit**

Run: `ls backend/internal/pkg/migrator/migrations/tenant/postgres/ | grep "^099"` and same for mysql — expect exactly the 2 files each just created.

```bash
git add backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/099_candidate_skills_certifications.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/099_candidate_skills_certifications.down.sql
git commit -m "feat: migration 099 candidate_skills + candidate_certifications (G-6)"
```

---

## Task 2: GORM Models — `CandidateSkill`, `CandidateCertification`

**Files:**
- Modify: `backend/internal/modules/recruitment/model.go`

**Interfaces:**
- Consumes: `competency.Competency` (read-only relation type, `backend/internal/modules/competency/model.go`).
- Produces: `CandidateSkill` struct, `CandidateCertification` struct — Task 3/4 depend on these exact field names.

- [ ] **Step 1: Add the competency-module import**

In `model.go`'s import block, add `"github.com/inthros/hris-platform/internal/modules/competency"` alongside the existing `"github.com/inthros/hris-platform/internal/modules/setting"` import.

- [ ] **Step 2: Add `CandidateSkill` and `CandidateCertification` structs**

Append after the `CandidateWorkExperience` block (same section of the file as the other G-6 sub-project 1 entities):

```go
// =========================================================================
// CandidateSkill (G-6 — skill kandidat, referensi competency.Competency)
// =========================================================================
// competency_id NOT NULL — beda dari candidate_educations yang nullable;
// skill tanpa referensi competency tidak actionable untuk candidate
// matching (G-9). Tidak ada Skill Master terpisah — reuse competencies
// master yang sudah ada (job_management/performance juga memakainya).

type CandidateSkill struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID  uuid.UUID `gorm:"type:char(36);not null;index:idx_cand_skill_candidate" json:"candidate_id"`
	CompetencyID uuid.UUID `gorm:"type:char(36);not null;index:idx_cand_skill_competency" json:"competency_id"`
	Level        *int      `gorm:"type:smallint" json:"level,omitempty"`
	Notes        string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relasi read-only (competency module) — dipakai untuk expand
	// CompetencyName di response, tidak diserialisasi langsung (pola
	// CandidateEducation.EducationMajor).
	Competency *competency.Competency `gorm:"foreignKey:CompetencyID" json:"-"`
}

func (CandidateSkill) TableName() string {
	return "candidate_skills"
}

func (s *CandidateSkill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// =========================================================================
// CandidateCertification (G-6 — sertifikasi kandidat, tanpa master)
// =========================================================================
// Tidak ada Certification Master di sistem — field bebas.

type CandidateCertification struct {
	ID                  uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	CandidateID         uuid.UUID `gorm:"type:char(36);not null;index:idx_cand_cert_candidate" json:"candidate_id"`
	Name                string    `gorm:"type:varchar(255);not null" json:"name"`
	IssuingOrganization *string   `gorm:"type:varchar(255)" json:"issuing_organization,omitempty"`
	IssueDate           *string   `gorm:"type:date" json:"issue_date,omitempty"`
	ExpiryDate          *string   `gorm:"type:date" json:"expiry_date,omitempty"`
	CredentialURL       *string   `gorm:"type:text" json:"credential_url,omitempty"`
	Notes               string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (CandidateCertification) TableName() string {
	return "candidate_certifications"
}

func (c *CandidateCertification) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
```

Note: `Level` is `*int` (not `*int16`) even though the SQL column is `SMALLINT` — confirmed by checking `competency.CompetencyValue.Level` (`*int` `gorm:"type:int"`) and `competency.CompetencyScoreDetail.StandardLevel`/`EmployeeLevel` (`*int` `gorm:"type:smallint"`) in `backend/internal/modules/competency/model.go:44,190,192` — this codebase's established convention for `smallint`-backed level fields is Go `*int`, and consistency with that sibling convention matters more than technically-precise Go integer width.

- [ ] **Step 3: Compile check**

Run: `cd backend && go build ./internal/modules/recruitment/...`
Expected: no errors. If the `competency` package import creates any issue, verify `competency`'s own imports don't reference `recruitment` (`grep -rn "modules/recruitment" backend/internal/modules/competency/`).

- [ ] **Step 4: Commit**

```bash
git add backend/internal/modules/recruitment/model.go
git commit -m "feat: model CandidateSkill + CandidateCertification (G-6)"
```

---

## Task 3: Repository — CRUD for both entities, with `Competency` preload

**Files:**
- Modify: `backend/internal/modules/recruitment/repository.go`
- Test: `backend/internal/modules/recruitment/repository_test.go`

**Interfaces:**
- Consumes: `CandidateSkill`, `CandidateCertification` (Task 2).
- Produces (used by Task 4 service):
  - `func (r *Repository) CreateCandidateSkill(ctx, s *CandidateSkill) error`
  - `func (r *Repository) FindCandidateSkillByID(ctx, id uuid.UUID) (*CandidateSkill, error)` — must `.Preload("Competency")`.
  - `func (r *Repository) ListCandidateSkills(ctx, candidateID uuid.UUID) ([]CandidateSkill, error)` — must `.Preload("Competency")`, order by `created_at ASC`.
  - `func (r *Repository) UpdateCandidateSkill(ctx, s *CandidateSkill) error`
  - `func (r *Repository) DeleteCandidateSkill(ctx, id uuid.UUID) error`
  - `func (r *Repository) FindCompetencyByID(ctx, id uuid.UUID) (*competency.Competency, error)` — direct query against the `competencies` table (no cross-module service call; `recruitment` doesn't have a dependency-injected reference to `competency`'s service), used by Task 4's existence-guard before creating a skill.
  - Same 5-method CRUD shape for `CandidateCertification` (no preload — no relation), ordered by `created_at ASC`.

- [ ] **Step 1: Write failing repository tests**

Append to `repository_test.go`, following the exact structure of `TestRepository_CreateAndFindCandidateEducation` (G-6 sub-project 1):

```go
func TestRepository_CreateAndFindCandidateSkill(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "Skill", LastName: "Test", Email: "skilltest@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp := &competency.Competency{Name: "Go Programming"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	skill := &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp.ID}
	if err := repo.CreateCandidateSkill(ctx, skill); err != nil {
		t.Fatalf("CreateCandidateSkill failed: %v", err)
	}

	found, err := repo.FindCandidateSkillByID(ctx, skill.ID)
	if err != nil {
		t.Fatalf("FindCandidateSkillByID failed: %v", err)
	}
	if found.CompetencyID != comp.ID {
		t.Errorf("expected competency_id %s, got %s", comp.ID, found.CompetencyID)
	}
}

func TestRepository_ListCandidateSkills(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "List", LastName: "Skill", Email: "listskill@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp1 := &competency.Competency{Name: "Go"}
	comp2 := &competency.Competency{Name: "SQL"}
	if err := db.Create(comp1).Error; err != nil {
		t.Fatalf("failed to seed competency comp1: %v", err)
	}
	if err := db.Create(comp2).Error; err != nil {
		t.Fatalf("failed to seed competency comp2: %v", err)
	}

	repo.CreateCandidateSkill(ctx, &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp1.ID})
	repo.CreateCandidateSkill(ctx, &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp2.ID})

	list, err := repo.ListCandidateSkills(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateSkills failed: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2 skills, got %d", len(list))
	}
}

func TestRepository_UpdateAndDeleteCandidateSkill(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	ctx := context.Background()

	cand := &Candidate{FirstName: "Upd", LastName: "Skill", Email: "updskill@test.com"}
	repo.CreateCandidate(ctx, cand)
	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}
	skill := &CandidateSkill{CandidateID: cand.ID, CompetencyID: comp.ID}
	repo.CreateCandidateSkill(ctx, skill)

	level := 4
	skill.Level = &level
	if err := repo.UpdateCandidateSkill(ctx, skill); err != nil {
		t.Fatalf("UpdateCandidateSkill failed: %v", err)
	}
	found, _ := repo.FindCandidateSkillByID(ctx, skill.ID)
	if found.Level == nil || *found.Level != level {
		t.Errorf("expected level %d, got %v", level, found.Level)
	}

	if err := repo.DeleteCandidateSkill(ctx, skill.ID); err != nil {
		t.Fatalf("DeleteCandidateSkill failed: %v", err)
	}
	if _, err := repo.FindCandidateSkillByID(ctx, skill.ID); err == nil {
		t.Error("expected error finding deleted skill, got nil")
	}
}

func TestRepository_CreateAndListCandidateCertification(t *testing.T) {
	repo, cleanup := newTestRepository()
	defer cleanup()
	ctx := context.Background()

	cand := &Candidate{FirstName: "Cert", LastName: "Test", Email: "certtest@test.com"}
	repo.CreateCandidate(ctx, cand)

	cert := &CandidateCertification{CandidateID: cand.ID, Name: "AWS Certified Solutions Architect"}
	if err := repo.CreateCandidateCertification(ctx, cert); err != nil {
		t.Fatalf("CreateCandidateCertification failed: %v", err)
	}

	list, err := repo.ListCandidateCertifications(ctx, cand.ID)
	if err != nil {
		t.Fatalf("ListCandidateCertifications failed: %v", err)
	}
	if len(list) != 1 || list[0].Name != "AWS Certified Solutions Architect" {
		t.Errorf("expected 1 cert named 'AWS Certified Solutions Architect', got %+v", list)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestRepository_CreateAndFindCandidateSkill -v`
Expected: compile error (methods undefined) or, once compiling, a runtime error if `competency.Competency` isn't yet in the test-DB `AutoMigrate` list (it needs to be added — see Task 3's helpers_test.go note below, or handle it in this same step since these tests need it to pass).

- [ ] **Step 3: Add `competency.Competency` to the test-DB AutoMigrate list**

`helpers_test.go`'s `setupTestDB()` needs `&competency.Competency{}` added to its `AutoMigrate(...)` call (alongside the existing `&setting.EducationMajor{}` from sub-project 1), and the `"github.com/inthros/hris-platform/internal/modules/competency"` import added to that file. Do this now so Step 2's tests can actually run against a real `competencies` table.

- [ ] **Step 4: Implement repository methods**

Append to `repository.go` (after the Candidate Work Experiences section):

```go
// =========================================================================
// Candidate Skills (G-6)
// =========================================================================

func (r *Repository) CreateCandidateSkill(ctx context.Context, s *CandidateSkill) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(s).Error
}

func (r *Repository) FindCandidateSkillByID(ctx context.Context, id uuid.UUID) (*CandidateSkill, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var s CandidateSkill
	if err := db.WithContext(ctx).Preload("Competency").First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate skill not found")
		}
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListCandidateSkills(ctx context.Context, candidateID uuid.UUID) ([]CandidateSkill, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateSkill
	if err := db.WithContext(ctx).Preload("Competency").Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateSkill(ctx context.Context, s *CandidateSkill) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Omit(clause.Associations).Save(s).Error
}

func (r *Repository) DeleteCandidateSkill(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateSkill{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate skill not found")
	}
	return result.Error
}

func (r *Repository) FindCompetencyByID(ctx context.Context, id uuid.UUID) (*competency.Competency, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c competency.Competency
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("competency not found")
		}
		return nil, err
	}
	return &c, nil
}

// =========================================================================
// Candidate Certifications (G-6)
// =========================================================================

func (r *Repository) CreateCandidateCertification(ctx context.Context, c *CandidateCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Create(c).Error
}

func (r *Repository) FindCandidateCertificationByID(ctx context.Context, id uuid.UUID) (*CandidateCertification, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var c CandidateCertification
	if err := db.WithContext(ctx).First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("candidate certification not found")
		}
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCandidateCertifications(ctx context.Context, candidateID uuid.UUID) ([]CandidateCertification, error) {
	db, err := r.db(ctx)
	if err != nil {
		return nil, err
	}
	var list []CandidateCertification
	if err := db.WithContext(ctx).Where("candidate_id = ?", candidateID).Order("created_at ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repository) UpdateCandidateCertification(ctx context.Context, c *CandidateCertification) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCandidateCertification(ctx context.Context, id uuid.UUID) error {
	db, err := r.db(ctx)
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Delete(&CandidateCertification{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("candidate certification not found")
	}
	return result.Error
}
```

> ⚠️ **Important — learned from G-6 sub-project 1's Critical fix-round finding:** `UpdateCandidateSkill` uses `Omit(clause.Associations)` on `Save` **from the start** (not retrofitted after a bug report this time) — this prevents GORM's `SaveBeforeAssociations` callback from re-deriving/resurrecting `CompetencyID` from a stale preloaded `Competency` pointer when a `CandidateSkill` loaded via `FindCandidateSkillByID` (which preloads `Competency`) gets updated and saved. Make sure `"gorm.io/gorm/clause"` is imported in `repository.go` if not already (it was added there during sub-project 1's fix — check first, it may already be imported).

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/recruitment/... -run "TestRepository_.*Candidate(Skill|Certification)" -v`
Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/repository.go backend/internal/modules/recruitment/repository_test.go backend/internal/modules/recruitment/helpers_test.go
git commit -m "feat: repository CRUD CandidateSkill + CandidateCertification (G-6)"
```

---

## Task 4: Service — CRUD, candidate + competency existence guards

**Files:**
- Modify: `backend/internal/modules/recruitment/service.go`
- Modify: `backend/internal/modules/recruitment/dto.go`
- Test: `backend/internal/modules/recruitment/service_test.go`

**Interfaces:**
- Consumes: repository methods from Task 3.
- Produces:
  - `func (s *Service) CreateCandidateSkill(ctx, candidateID string, req CreateCandidateSkillRequest) (*CandidateSkillResponse, error)` — validates candidate exists AND competency exists (both 404-equivalent errors if not).
  - `func (s *Service) ListCandidateSkills(ctx, candidateID string) ([]CandidateSkillResponse, error)`
  - `func (s *Service) UpdateCandidateSkill(ctx, id string, req UpdateCandidateSkillRequest) (*CandidateSkillResponse, error)`
  - `func (s *Service) DeleteCandidateSkill(ctx, id string) error`
  - Same 4-method shape for `CandidateCertification` (candidate-existence guard only, no competency to check).

- [ ] **Step 1: Add DTOs to `dto.go`**

```go
// =========================================================================
// Candidate Skill DTOs (G-6)
// =========================================================================

type CreateCandidateSkillRequest struct {
	CompetencyID string `json:"competency_id" binding:"required"`
	Level        *int   `json:"level"`
	Notes        string `json:"notes"`
}

type UpdateCandidateSkillRequest struct {
	Level *int    `json:"level"`
	Notes *string `json:"notes"`
}

type CandidateSkillResponse struct {
	ID             string `json:"id"`
	CandidateID    string `json:"candidate_id"`
	CompetencyID   string `json:"competency_id"`
	CompetencyName string `json:"competency_name,omitempty"`
	Level          int    `json:"level,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

// =========================================================================
// Candidate Certification DTOs (G-6)
// =========================================================================

type CreateCandidateCertificationRequest struct {
	Name                string  `json:"name" binding:"required,max=255"`
	IssuingOrganization *string `json:"issuing_organization" binding:"omitempty,max=255"`
	IssueDate           *string `json:"issue_date" binding:"omitempty,max=10"`
	ExpiryDate          *string `json:"expiry_date" binding:"omitempty,max=10"`
	CredentialURL       *string `json:"credential_url"`
	Notes               string  `json:"notes"`
}

type UpdateCandidateCertificationRequest struct {
	Name                *string `json:"name" binding:"omitempty,max=255"`
	IssuingOrganization *string `json:"issuing_organization" binding:"omitempty,max=255"`
	IssueDate           *string `json:"issue_date" binding:"omitempty,max=10"`
	ExpiryDate          *string `json:"expiry_date" binding:"omitempty,max=10"`
	CredentialURL       *string `json:"credential_url"`
	Notes               *string `json:"notes"`
}

type CandidateCertificationResponse struct {
	ID                  string `json:"id"`
	CandidateID         string `json:"candidate_id"`
	Name                string `json:"name"`
	IssuingOrganization string `json:"issuing_organization,omitempty"`
	IssueDate           string `json:"issue_date,omitempty"`
	ExpiryDate          string `json:"expiry_date,omitempty"`
	CredentialURL       string `json:"credential_url,omitempty"`
	Notes               string `json:"notes,omitempty"`
}
```

- [ ] **Step 2: Write failing service tests**

Append to `service_test.go` (create a `competency.Competency` fixture the same way Task 3's repository tests do):

```go
func TestService_CreateCandidateSkill(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Skill", LastName: "Svc", Email: "skillsvc@test.com"})
	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	resp, err := svc.CreateCandidateSkill(ctx, cand.ID, CreateCandidateSkillRequest{CompetencyID: comp.ID.String()})
	if err != nil {
		t.Fatalf("CreateCandidateSkill failed: %v", err)
	}
	if resp.CompetencyName != "Go" {
		t.Errorf("expected competency_name 'Go', got %s", resp.CompetencyName)
	}
}

func TestService_CreateCandidateSkill_UnknownCandidate(t *testing.T) {
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	ctx := context.Background()
	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	_, err := svc.CreateCandidateSkill(ctx, uuid.New().String(), CreateCandidateSkillRequest{CompetencyID: comp.ID.String()})
	if err == nil {
		t.Fatal("expected error for unknown candidate, got nil")
	}
}

func TestService_CreateCandidateSkill_UnknownCompetency(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Skill2", LastName: "Svc", Email: "skillsvc2@test.com"})

	_, err := svc.CreateCandidateSkill(ctx, cand.ID, CreateCandidateSkillRequest{CompetencyID: uuid.New().String()})
	if err == nil {
		t.Fatal("expected error for unknown competency, got nil")
	}
}

func TestService_CreateCandidateCertification(t *testing.T) {
	svc, cleanup := newTestService()
	defer cleanup()
	ctx := context.Background()

	cand, _ := svc.CreateCandidate(ctx, CreateCandidateRequest{FirstName: "Cert", LastName: "Svc", Email: "certsvc@test.com"})

	resp, err := svc.CreateCandidateCertification(ctx, cand.ID, CreateCandidateCertificationRequest{Name: "AWS SAA"})
	if err != nil {
		t.Fatalf("CreateCandidateCertification failed: %v", err)
	}
	if resp.Name != "AWS SAA" {
		t.Errorf("expected name 'AWS SAA', got %s", resp.Name)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestService_CreateCandidateSkill -v`
Expected: compile error.

- [ ] **Step 4: Implement service CRUD methods**

```go
// =========================================================================
// Candidate Skills (G-6)
// =========================================================================

func (s *Service) CreateCandidateSkill(ctx context.Context, candidateID string, req CreateCandidateSkillRequest) (*CandidateSkillResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	compUUID, err := uuid.Parse(req.CompetencyID)
	if err != nil {
		return nil, fmt.Errorf("invalid competency_id: %w", err)
	}
	if _, err := s.repo.FindCompetencyByID(ctx, compUUID); err != nil {
		return nil, fmt.Errorf("competency not found: %w", err)
	}

	sk := &CandidateSkill{
		CandidateID:  candUUID,
		CompetencyID: compUUID,
		Level:        req.Level,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateCandidateSkill(ctx, sk); err != nil {
		return nil, err
	}
	created, err := s.repo.FindCandidateSkillByID(ctx, sk.ID)
	if err != nil {
		return nil, err
	}
	return candidateSkillToResponse(created), nil
}

func (s *Service) ListCandidateSkills(ctx context.Context, candidateID string) ([]CandidateSkillResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateSkills(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateSkillResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateSkillToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateSkill(ctx context.Context, id string, req UpdateCandidateSkillRequest) (*CandidateSkillResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	sk, err := s.repo.FindCandidateSkillByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Level != nil {
		sk.Level = req.Level
	}
	if req.Notes != nil {
		sk.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateSkill(ctx, sk); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindCandidateSkillByID(ctx, sk.ID)
	if err != nil {
		return nil, err
	}
	return candidateSkillToResponse(updated), nil
}

func (s *Service) DeleteCandidateSkill(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateSkill(ctx, uid)
}

func candidateSkillToResponse(s *CandidateSkill) *CandidateSkillResponse {
	resp := &CandidateSkillResponse{
		ID:           s.ID.String(),
		CandidateID:  s.CandidateID.String(),
		CompetencyID: s.CompetencyID.String(),
		Notes:        s.Notes,
	}
	if s.Competency != nil {
		resp.CompetencyName = s.Competency.Name
	}
	if s.Level != nil {
		resp.Level = *s.Level
	}
	return resp
}

// =========================================================================
// Candidate Certifications (G-6)
// =========================================================================

func (s *Service) CreateCandidateCertification(ctx context.Context, candidateID string, req CreateCandidateCertificationRequest) (*CandidateCertificationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	if _, err := s.repo.FindCandidateByID(ctx, candUUID); err != nil {
		return nil, fmt.Errorf("candidate not found: %w", err)
	}
	c := &CandidateCertification{
		CandidateID:         candUUID,
		Name:                req.Name,
		IssuingOrganization: req.IssuingOrganization,
		IssueDate:           req.IssueDate,
		ExpiryDate:          req.ExpiryDate,
		CredentialURL:       req.CredentialURL,
		Notes:               req.Notes,
	}
	if err := s.repo.CreateCandidateCertification(ctx, c); err != nil {
		return nil, err
	}
	return candidateCertificationToResponse(c), nil
}

func (s *Service) ListCandidateCertifications(ctx context.Context, candidateID string) ([]CandidateCertificationResponse, error) {
	candUUID, err := uuid.Parse(candidateID)
	if err != nil {
		return nil, fmt.Errorf("invalid candidate_id: %w", err)
	}
	list, err := s.repo.ListCandidateCertifications(ctx, candUUID)
	if err != nil {
		return nil, err
	}
	out := make([]CandidateCertificationResponse, 0, len(list))
	for i := range list {
		out = append(out, *candidateCertificationToResponse(&list[i]))
	}
	return out, nil
}

func (s *Service) UpdateCandidateCertification(ctx context.Context, id string, req UpdateCandidateCertificationRequest) (*CandidateCertificationResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}
	c, err := s.repo.FindCandidateCertificationByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.IssuingOrganization != nil {
		c.IssuingOrganization = req.IssuingOrganization
	}
	if req.IssueDate != nil {
		c.IssueDate = req.IssueDate
	}
	if req.ExpiryDate != nil {
		c.ExpiryDate = req.ExpiryDate
	}
	if req.CredentialURL != nil {
		c.CredentialURL = req.CredentialURL
	}
	if req.Notes != nil {
		c.Notes = *req.Notes
	}
	if err := s.repo.UpdateCandidateCertification(ctx, c); err != nil {
		return nil, err
	}
	return candidateCertificationToResponse(c), nil
}

func (s *Service) DeleteCandidateCertification(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}
	return s.repo.DeleteCandidateCertification(ctx, uid)
}

func candidateCertificationToResponse(c *CandidateCertification) *CandidateCertificationResponse {
	resp := &CandidateCertificationResponse{
		ID:          c.ID.String(),
		CandidateID: c.CandidateID.String(),
		Name:        c.Name,
		Notes:       c.Notes,
	}
	if c.IssuingOrganization != nil {
		resp.IssuingOrganization = *c.IssuingOrganization
	}
	if c.IssueDate != nil {
		resp.IssueDate = *c.IssueDate
	}
	if c.ExpiryDate != nil {
		resp.ExpiryDate = *c.ExpiryDate
	}
	if c.CredentialURL != nil {
		resp.CredentialURL = *c.CredentialURL
	}
	return resp
}
```

- [ ] **Step 5: Run all recruitment service tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v 2>&1 | tail -100`
Expected: all pass, including every pre-existing test.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/service.go backend/internal/modules/recruitment/service_test.go backend/internal/modules/recruitment/dto.go
git commit -m "feat: service CRUD CandidateSkill + CandidateCertification (G-6)"
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
  POST   /recruitment/candidates/:id/skills
  GET    /recruitment/candidates/:id/skills
  PUT    /recruitment/skills/:id
  DELETE /recruitment/skills/:id
  POST   /recruitment/candidates/:id/certifications
  GET    /recruitment/candidates/:id/certifications
  PUT    /recruitment/certifications/:id
  DELETE /recruitment/certifications/:id
  ```

- [ ] **Step 1: Write failing handler tests**

`setupTestRouter()` does not expose its underlying `*gorm.DB` (it only returns `(*gin.Engine, *Handler, func())`), so the skill handler test — which needs a real `competency.Competency` row to reference since `competency_id` is required — must build the router manually via `setupTestDB()` instead, mirroring `setupTestRouter()`'s own body (`handler_test.go:15-29`) plus the extra `db.Create(comp)` step:

```go
func TestHandler_CreateCandidateSkill(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, dbResolver, cleanup := setupTestDB()
	defer cleanup()
	repo := NewRepository(dbResolver)
	svc := NewService(repo, zap.NewNop())
	seedDefaultRecruitmentStages(db)
	handler := NewHandler(svc)
	r := gin.New()
	rg := r.Group("/api/v1/tenant")
	RegisterRoutes(rg, handler)

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Skill", LastName: "Handler", Email: "skillhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	comp := &competency.Competency{Name: "Go"}
	if err := db.Create(comp).Error; err != nil {
		t.Fatalf("failed to seed competency: %v", err)
	}

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/skills", CreateCandidateSkillRequest{
		CompetencyID: comp.ID.String(),
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandler_CreateCandidateCertification(t *testing.T) {
	r, _, cleanup := setupTestRouter()
	defer cleanup()

	cW := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates", CreateCandidateRequest{
		FirstName: "Cert", LastName: "Handler", Email: "certhandler@test.com",
	})
	var cResp map[string]interface{}
	json.Unmarshal(cW.Body.Bytes(), &cResp)
	cid := cResp["data"].(map[string]interface{})["id"].(string)

	w := performRequest(r, "POST", "/api/v1/tenant/recruitment/candidates/"+cid+"/certifications", CreateCandidateCertificationRequest{
		Name: "AWS SAA",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}
```

- [ ] **Step 2: Run to verify failure**

Run: `cd backend && go test ./internal/modules/recruitment/... -run TestHandler_CreateCandidateSkill -v`
Expected: 404 (route not registered) or compile error.

- [ ] **Step 3: Implement handlers**

Add to `handler.go`, mirroring `CreateCandidateEducation`/`ListCandidateEducations`/`UpdateCandidateEducation`/`DeleteCandidateEducation` exactly, substituting the skill/certification service methods:

```go
func (h *Handler) CreateCandidateSkill(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateSkillRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateSkill(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateSkills(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateSkills(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateSkill(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateSkillRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateSkill(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateSkill(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateSkill(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}

func (h *Handler) CreateCandidateCertification(c *gin.Context) {
	candidateID := c.Param("id")
	var req CreateCandidateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.CreateCandidateCertification(c.Request.Context(), candidateID, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.CreatedJSON(c, resp, "success.created")
}

func (h *Handler) ListCandidateCertifications(c *gin.Context) {
	candidateID := c.Param("id")
	resp, err := h.svc.ListCandidateCertifications(c.Request.Context(), candidateID)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) UpdateCandidateCertification(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCandidateCertificationRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	resp, err := h.svc.UpdateCandidateCertification(c.Request.Context(), id, req)
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.SuccessJSON(c, resp)
}

func (h *Handler) DeleteCandidateCertification(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteCandidateCertification(c.Request.Context(), id); err != nil {
		httputil.NotFound(c, err.Error())
		return
	}
	httputil.DeletedJSON(c, "success.deleted")
}
```

- [ ] **Step 4: Register routes**

In `routes.go`, add after the existing `work-experiences/:id` routes:

```go
	rec.POST("/candidates/:id/skills", handler.CreateCandidateSkill)
	rec.GET("/candidates/:id/skills", handler.ListCandidateSkills)
	rec.PUT("/skills/:id", handler.UpdateCandidateSkill)
	rec.DELETE("/skills/:id", handler.DeleteCandidateSkill)
	rec.POST("/candidates/:id/certifications", handler.CreateCandidateCertification)
	rec.GET("/candidates/:id/certifications", handler.ListCandidateCertifications)
	rec.PUT("/certifications/:id", handler.UpdateCandidateCertification)
	rec.DELETE("/certifications/:id", handler.DeleteCandidateCertification)
```

- [ ] **Step 5: Run handler tests**

Run: `cd backend && go test ./internal/modules/recruitment/... -v 2>&1 | tail -60`
Expected: all pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/recruitment/handler.go backend/internal/modules/recruitment/routes.go backend/internal/modules/recruitment/handler_test.go
git commit -m "feat: handler+routes CandidateSkill/CandidateCertification CRUD (G-6)"
```

---

## Task 6: Module wiring — `AutoMigrate` (production consistency)

**Files:**
- Modify: `backend/internal/modules/recruitment/module.go`

**Interfaces:**
- Consumes: `CandidateSkill`, `CandidateCertification` (Task 2).

- [ ] **Step 1: Add both entities to `Migrate()`**

Add `&CandidateSkill{}` and `&CandidateCertification{}` to the existing `db.AutoMigrate(...)` call in `module.go`'s `Migrate()` method, alongside `&CandidateEducation{}`, `&CandidateWorkExperience{}`.

> Note: as established in prior plans' final reviews, this `Migrate()`/`Seed()` wiring is not actually invoked for tenant databases in production (a separately-tracked, pre-existing bug in `cmd/server/main.go`). This step is for consistency with the module's existing pattern; the real schema source of truth is migration 099 from Task 1. Do NOT attempt to fix the broader wiring bug here — out of scope.

Do NOT touch `helpers_test.go` in this task — Task 3 already added `&CandidateSkill{}`/`&CandidateCertification{}`/`&competency.Competency{}` to the test-DB `AutoMigrate` list out of necessity for its own tests. Verify this by reading the file first; if for some reason it's missing, add it, but expect it to already be there.

- [ ] **Step 2: Build + test run**

Run: `cd backend && go build ./... && go test ./internal/modules/recruitment/... -v 2>&1 | tail -60`
Expected: clean build, all tests pass.

- [ ] **Step 3: Commit**

```bash
git add backend/internal/modules/recruitment/module.go
git commit -m "feat: AutoMigrate CandidateSkill + CandidateCertification (G-6)"
```

---

## Task 7: Update plan doc status

**Files:**
- Modify: `docs/archive/module-recruitment-development-plan.md`

- [ ] **Step 1: Update §G-6 status**

Update the G-6 section (currently marked partial after sub-project 1) to reflect sub-project 2 done: add `candidate_skills` + `candidate_certifications` to the "Yang diimplementasikan" list, following the exact style already established for sub-project 1 (migration number 099, test count delta — count actual new tests added by this plan, do not guess). Update the "Rencana (sisa G-6)" list to remove skills/certifications from the deferred items, leaving only `candidates.status`, `source_id`, `candidate_documents`, `candidate_consents` as still open (sub-project 3).

- [ ] **Step 2: Update §8.1/§8.2 API Plan**

Add the 8 new endpoints to §8.1 (existing), increment the endpoint count.

- [ ] **Step 3: Verify no other doc sections need updating**

Check §23 Definition of Done — if there's a line about "profile terstruktur (education/experience/skills/certification/document)" already annotated with partial-progress markers (from sub-project 1's final-review fix), update the annotation to reflect skills/certifications now also done, e.g. `🔶 sebagian: education/experience/skills/certifications ✅ (G-6 sub-project 1+2), documents ❌`.

- [ ] **Step 4: Commit**

```bash
git add docs/module-recruitment-development-plan.md
git commit -m "docs: update G-6 status - skills+certifications selesai (sub-project 2/3)"
```
