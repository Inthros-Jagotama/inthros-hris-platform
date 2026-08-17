# Candidate Profile Basics — Education & Work Experience (G-6 sub-project 1) — Design

> Ref: `docs/archive/module-recruitment-development-plan.md` §G-6 (partial — first of 3 planned sub-projects).

## Goal

G-6 in the plan bundles four `candidates` column additions and five new sub-tables into one large item. This sub-project narrows scope to the foundational, self-contained piece: an auto-generated `candidate_number`, and two profile sub-tables — `candidate_educations` and `candidate_work_experiences` — so a candidate's education and employment history can be recorded structurally instead of only via free-text `resume_url`/`current_company`/`current_title`.

## Explicitly out of scope (deferred to later sub-projects)

- `candidates.status` (candidate-level availability status like ACTIVE/BLACKLISTED) — no clear requirement yet, skipped.
- `candidates.source_id` / a recruitment-source master table — `source` stays free text; building a source master is a separate, larger piece of work.
- `candidates.consent_status`, `candidate_consents`, `candidate_documents` — compliance/file-storage concerns with their own security considerations (plan §17), handled in a later sub-project.
- `candidate_skills`, `candidate_certifications` — blocked on a decision about whether/how to reference a Skill Master or Certification Master, neither of which exists in the codebase today (verified: no `Skill`/`Certification` master entity in `competency` or elsewhere). Deferred.
- No frontend changes in this iteration.

## Data Model

### `candidates.candidate_number` (new column)

Added via migration to the existing `candidates` table: `candidate_number VARCHAR(50) NULL`.

Auto-generated on create, format `CAND-YYYYMM-XXXXXXXX`, following the exact existing pattern of `generateRequisitionNumber`/`generateOfferNumber` (`service.go:1938`, `service.go:2042`): `fmt.Sprintf("CAND-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.New().String()[:8]))`. Can be overridden by an explicit value in the create/update request, same as `requisition_number`/`offer_number`.

### `candidate_educations` (new table)

Mirrors the existing `EmployeeEducation` pattern (`backend/internal/modules/employee/model.go:148-163`) rather than inventing a new shape: `level` references the existing seeded `setting.Education` master (rows are exactly `SD/SMP/SMA/D3/S1/S2/S3/...`, see `tenantseed/seed_data.go`), and `major` references the existing seeded `setting.EducationMajor` master, both nullable with a free-text fallback for values not present in the master (same dual `*ID` + free-text-string pattern `EmployeeEducation` already uses for major).

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, index |
| `education_id` | CHAR(36) NULL | FK → `setting.educations` (the level master: SD/SMP/SMA/D3/S1/S2/S3/...) — nullable since not every historical record will match a master row |
| `institution_name` | VARCHAR(255) NOT NULL | |
| `education_major_id` | CHAR(36) NULL | FK → `setting.education_majors` |
| `major` | VARCHAR(255) NULL | free-text fallback when the major isn't in the master (mirrors `EmployeeEducation.Major`) |
| `gpa` | DECIMAL(3,2) NULL | |
| `start_year` | INT NULL | |
| `end_year` | INT NULL | |
| `is_highest` | BOOLEAN NOT NULL DEFAULT false | marks the candidate's highest completed education among their rows (no DB-level single-row enforcement — service layer may optionally unset others on write, but this is not required for this iteration; a candidate may have zero or more than one marked `true` if the caller sets it inconsistently) |
| `notes` | TEXT NULL | |
| `created_at`, `updated_at` | TIMESTAMP | |

GORM model should include the same read-only relation fields `EmployeeEducation` uses: `Education *setting.Education` and `EducationMajor *setting.EducationMajor` (both `gorm:"foreignKey:..."`, `json:"-"` — not serialized directly). The response DTO instead exposes flattened expanded names, mirroring `employee.EducationResponse` (`dto.go:228-236`) exactly:

```go
type CandidateEducationResponse struct {
	ID               string `json:"id"`
	EducationID      string `json:"education_id,omitempty"`
	EducationMajorID string `json:"education_major_id,omitempty"`
	MajorName        string `json:"major_name,omitempty"` // resolved from EducationMajor.Name when education_major_id is set
	InstitutionName  string `json:"institution_name"`
	Major            string `json:"major,omitempty"`       // free-text fallback
	GPA              string `json:"gpa,omitempty"`
	StartYear        int    `json:"start_year,omitempty"`
	EndYear          int    `json:"end_year,omitempty"`
	IsHighest        bool   `json:"is_highest"`
	Notes            string `json:"notes,omitempty"`
}
```

`candidate_work_experiences` has no analogous master reference (employer/job-title have no master table in this codebase) — its design stays as specified below, unrelated to this change.

### `candidate_work_experiences` (new table)

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, index |
| `company_name` | VARCHAR(255) NOT NULL | |
| `job_title` | VARCHAR(255) NOT NULL | |
| `employment_type` | VARCHAR(50) NULL | free text, e.g. `FULL_TIME/PART_TIME/CONTRACT/INTERNSHIP` — not a master reference |
| `start_date` | DATE NOT NULL | |
| `end_date` | DATE NULL | NULL means still employed there |
| `is_current` | BOOLEAN NOT NULL DEFAULT false | |
| `description` | TEXT NULL | |
| `created_at`, `updated_at` | TIMESTAMP | |

No business-rule validation is required between `end_date`/`is_current` consistency in this iteration (e.g. not enforcing `is_current=true` implies `end_date=NULL`) — both fields are simply stored as given; a future iteration can tighten this if it becomes a real problem.

## API

Standard CRUD, following this module's existing per-entity pattern (see `Interview` or `OnboardingTaskTemplate` CRUD in `handler.go`/`routes.go` for the exact shape to mirror):

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

`candidate_number` requires no new endpoint — it's returned as part of the existing `CandidateResponse` (add the field) and settable via the existing `CreateCandidateRequest`/`UpdateCandidateRequest`.

## Error Handling

- Creating an education/work-experience row for a non-existent `candidate_id` → 404 (mirror `CreateApplication`'s requisition/candidate-exists check pattern).
- Standard validation via Gin `binding` tags on required fields (`level`, `institution_name` for educations; `company_name`, `job_title`, `start_date` for work experiences).

## Testing Plan

- Repository: create + list round-trip for both new tables (mirror Task 3's pattern from the G-5 work).
- Service: CRUD for both entities; `candidate_number` auto-generation format test; explicit-override test.
- Handler: create/list/update/delete for both entities, 404 on unknown candidate.

## Files Touched (summary — exact paths/code in the implementation plan)

- New migration (next available number after 097 — confirm at implementation time) for `candidates.candidate_number` + the two new tables, mysql+postgres, up+down, idempotent.
- `backend/internal/modules/recruitment/model.go` — `CandidateEducation`, `CandidateWorkExperience` structs; add `CandidateNumber` field to `Candidate`.
- `backend/internal/modules/recruitment/repository.go` — CRUD methods for both.
- `backend/internal/modules/recruitment/service.go` — CRUD methods, `generateCandidateNumber()` helper (mirrors `generateRequisitionNumber`), wire into `CreateCandidate`.
- `backend/internal/modules/recruitment/dto.go` — request/response DTOs for both entities; add `candidate_number` to `CandidateResponse`.
- `backend/internal/modules/recruitment/handler.go` + `routes.go` — new endpoints.
- `backend/internal/modules/recruitment/module.go` — `AutoMigrate` additions (test-DB convenience only, not the real migration).
- `docs/archive/module-recruitment-development-plan.md` — update G-6 status to reflect partial completion (this sub-project only; skills/certifications/documents/consents/status/source_id remain open).
