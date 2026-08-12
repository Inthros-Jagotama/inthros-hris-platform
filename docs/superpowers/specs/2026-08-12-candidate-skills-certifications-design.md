# Candidate Skills & Certifications (G-6 sub-project 2) — Design

> Ref: `docs/module-recruitment-development-plan.md` §G-6 (partial — second of 3 planned sub-projects; sub-project 1 was `candidate_number` + `candidate_educations` + `candidate_work_experiences`, migration 098).

## Goal

Add two more candidate profile sub-tables: `candidate_skills` and `candidate_certifications`, following the same narrow-scope pattern as sub-project 1 — backend only, reusing existing masters where they exist, and not inventing a new master where none does.

## Explicitly out of scope (deferred to sub-project 3 or later)

- `candidates.status`, `source_id` — same as sub-project 1's decision, still not built.
- `candidate_documents`, `candidate_consents` — deferred to sub-project 3 (own compliance/file-storage concerns, plan §17).
- No new "Certification Master" catalog table — none exists in the codebase today (`training.TrainingCertificate` is a participant-issued record from a specific training session, not a reusable catalog of certification types) and building one is out of scope here.
- No frontend changes.

## Data Model

### `candidate_skills` (new table)

Unlike `candidate_educations`' nullable master references (sub-project 1), `competency_id` here is **NOT NULL** — a skill entry without a link to the `competencies` master isn't actionable for the candidate-matching feature planned in G-9 (`docs/module-recruitment-development-plan.md` §G-9, "Candidate Matching memakai Job Requirement + Competency + ... → candidate_match_score"), so every skill row must resolve to a real competency.

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, `ON DELETE CASCADE`, index |
| `competency_id` | CHAR(36) NOT NULL | FK → `competencies` (existing master, `backend/internal/modules/competency/model.go`), `ON DELETE CASCADE` (a skill row referencing a deleted competency is meaningless, so cascade rather than orphan it) |
| `level` | SMALLINT NULL | 1-5, matching the level-scale convention already used by `competency.CompetencyValue.Level`/`CompetencyScoreDetail.StandardLevel` in the same module |
| `notes` | TEXT NULL | |
| `created_at`, `updated_at` | TIMESTAMP | |

No uniqueness constraint on `(candidate_id, competency_id)` in this iteration — a candidate could theoretically have the same competency listed twice (e.g. self-reported vs. verified); not worth enforcing now, can be revisited if it becomes a real data-quality problem.

### `candidate_certifications` (new table)

No master reference — mirrors `candidate_work_experiences`' free-text-only design from sub-project 1.

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, `ON DELETE CASCADE`, index |
| `name` | VARCHAR(255) NOT NULL | free text, e.g. "AWS Certified Solutions Architect" |
| `issuing_organization` | VARCHAR(255) NULL | |
| `issue_date` | DATE NULL | |
| `expiry_date` | DATE NULL | NULL = does not expire / unknown |
| `credential_url` | TEXT NULL | verification link, if provided |
| `notes` | TEXT NULL | |
| `created_at`, `updated_at` | TIMESTAMP | |

No expiry validation logic in this iteration (e.g. no "is this certification expired" computed flag) — just storage; a future iteration can add that if a real use case needs it.

## API

Same nested-CRUD shape as sub-project 1, with the same `:id` route-parameter naming (not `:candidate_id`) — confirmed necessary due to the Gin wildcard-name constraint discovered and fixed during sub-project 1 (`/candidates/:id` already exists as a top-level route, so any route nested under `/candidates/...` must reuse the same param name `:id`, not a differently-named one like `:candidate_id`).

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

## Error Handling

- Creating a skill/certification row for a non-existent `candidate_id` → error (existing pattern: `s.repo.FindCandidateByID` guard, mirrored from sub-project 1's `CreateCandidateEducation`/`CreateCandidateWorkExperience`).
- Creating a skill row for a non-existent `competency_id` → error, validated via a `FindCompetencyByID`-equivalent lookup before insert (reusing the `competency` module's repository if accessible, or a lightweight cross-module read — follow whatever pattern `employee` module already uses to reference `setting.Education`, i.e. import the `competency` package's model type for the GORM relation, and add a repository method in `recruitment`'s own repository that queries the `competencies` table directly via the shared DB handle, since `recruitment` doesn't have a dependency-injected reference to `competency`'s service/repository — this mirrors how `candidate_educations` reads `setting.Education` directly via GORM relation rather than calling into the `setting` module's service).
- Standard Gin binding validation on required fields (`competency_id`, `name`).

## Testing Plan

- Repository: create + list round-trip for both new tables; verify the `Competency` relation preloads correctly for skill responses (mirroring sub-project 1's `EducationMajor` preload test pattern), including a real `competency.Competency` row created via GORM in the test (added to the shared test-DB `AutoMigrate` list, matching how `setting.EducationMajor` was added in sub-project 1's Task 4 fix).
- Service: CRUD for both entities; candidate-existence guard test; competency-existence guard test for skills; level range is NOT validated at this layer (no business rule requiring 1-5 enforcement in this iteration — if a caller sends `level: 99`, it's stored as-is; add validation later if this becomes a real problem).
- Handler: create/list/update/delete for both entities, error case for unknown candidate and unknown competency.

## Files Touched (summary — exact paths/code in the implementation plan)

- New migration **099** (mysql+postgres, up+down, idempotent) — next after `098_candidate_profile_basics`.
- `backend/internal/modules/recruitment/model.go` — `CandidateSkill`, `CandidateCertification` structs; `CandidateSkill.Competency *competency.Competency` relation (import `backend/internal/modules/competency`).
- `backend/internal/modules/recruitment/repository.go` — CRUD for both, `FindCompetencyByID` helper query (direct GORM query against `competencies` table, no cross-module service call).
- `backend/internal/modules/recruitment/service.go` — CRUD for both, candidate + competency existence guards.
- `backend/internal/modules/recruitment/dto.go` — request/response DTOs for both entities.
- `backend/internal/modules/recruitment/handler.go` + `routes.go` — 8 new endpoints.
- `backend/internal/modules/recruitment/module.go` — `AutoMigrate` additions (test-DB/production consistency, not the real migration mechanism).
- `docs/module-recruitment-development-plan.md` — update G-6 status (sub-project 2 of 3 done; documents/consents remain open).
