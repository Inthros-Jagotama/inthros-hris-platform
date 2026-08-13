# Candidate Documents (G-6 sub-project 3a) — Design

> Ref: `docs/module-recruitment-development-plan.md` §G-6 (partial — third of the originally-planned sub-projects; `candidate_consents` is split out to its own separate sub-project, not part of this one).

## Goal

Add `candidate_documents`, a table storing references to files a candidate has attached (resume, cover letter, certificate scan, portfolio, identity document, other) — a reference table, not a file-storage implementation.

## Explicitly out of scope

- `candidate_consents` — split into its own future sub-project (compliance/GDPR concerns are a distinct problem from document references).
- `candidates.status`, `source_id` — still not built, unchanged from earlier sub-projects' decisions.
- **No new upload endpoint.** The codebase already has a generic file-upload endpoint, `POST /api/v1/tenant/uploads` (`backend/internal/pkg/upload/handler.go`), which accepts a multipart `file` field, validates extension/size (10MB cap, images/pdf/office-doc/text types), stores it under `{uploadDir}/attachments/{uuid}{ext}`, and returns `{"url": "/uploads/attachments/{uuid}.pdf"}`. `candidate_documents` only stores a reference to a URL already produced by that endpoint — this sub-project does not touch the upload/storage mechanism at all.
- **No new permission.** Reuses the existing `recruitment.view/create/update/delete` permissions already registered for the module — no new granular permission (e.g. a hypothetical `recruitment.candidate.document.view`) in this iteration, per explicit decision during design.
- No file content validation beyond what the shared `/uploads` endpoint already does (extension allowlist, size cap) — `candidate_documents` just stores whatever URL string it's given.
- No frontend changes.

## Data Model

### `candidate_documents` (new table)

Mirrors `employee.EmployeeDocument`'s reference-only shape (`Name`/`File` fields), but omits `CreatedBy`/`UpdatedBy` — no table in the `recruitment` module tracks actor-on-create/update (confirmed: `grep -n "CreatedBy" backend/internal/modules/recruitment/model.go` returns nothing), so `candidate_documents` follows the recruitment module's own convention rather than importing employee's.

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, `ON DELETE CASCADE`, index |
| `document_type` | VARCHAR(20) NOT NULL DEFAULT 'OTHER' | `RESUME, COVER_LETTER, CERTIFICATE, PORTFOLIO, IDENTITY, OTHER` — enforced via Gin `binding:"oneof=..."` at the request layer, not a DB constraint (matches this module's existing convention — e.g. `CandidateType`, `OfferStatus` etc. are plain `VARCHAR` with app-layer enum enforcement, not SQL `CHECK`/`ENUM`) |
| `name` | VARCHAR(255) NOT NULL | display label / original filename |
| `file_url` | TEXT NOT NULL | URL returned by `POST /uploads` (or any other absolute/relative URL the caller supplies — this table does not validate the URL points at a real uploaded file) |
| `notes` | TEXT NULL | |
| `created_at`, `updated_at` | TIMESTAMP | |

## API

Same nested-CRUD shape and `:id` route-parameter convention as the module's other candidate sub-resources (per the established Gin wildcard-name constraint):

```
POST   /recruitment/candidates/:id/documents
GET    /recruitment/candidates/:id/documents
PUT    /recruitment/documents/:id
DELETE /recruitment/documents/:id
```

No `PATCH`-only partial file replace — updating `file_url` on an existing document row is just a normal `PUT` field update (caller re-uploads via `/uploads` first, then `PUT`s the new URL in). No separate "replace file" endpoint in this iteration.

## Error Handling

- Creating a document for a non-existent `candidate_id` → error (existing pattern: `s.repo.FindCandidateByID` guard, same as every other candidate sub-resource in this module).
- `document_type` validated via `binding:"omitempty,oneof=RESUME COVER_LETTER CERTIFICATE PORTFOLIO IDENTITY OTHER"` on create (defaults to `OTHER` if omitted, matching the SQL column default) and `binding:"omitempty,oneof=..."` on update.
- No validation that `file_url` actually resolves to an existing file — out of scope, matches the "reference only" design.

## Testing Plan

- Repository: create + list round-trip.
- Service: CRUD + candidate-existence guard test.
- Handler: create/list/update/delete, 8-document-type binding validation test (reject an invalid `document_type` value with 400).

## Files Touched (summary — exact paths/code in the implementation plan)

- New migration **100** (mysql+postgres, up+down, idempotent).
- `backend/internal/modules/recruitment/model.go` — `CandidateDocument` struct.
- `backend/internal/modules/recruitment/repository.go` — CRUD.
- `backend/internal/modules/recruitment/service.go` — CRUD + candidate-existence guard.
- `backend/internal/modules/recruitment/dto.go` — request/response DTOs.
- `backend/internal/modules/recruitment/handler.go` + `routes.go` — 4 new endpoints.
- `backend/internal/modules/recruitment/module.go` — `AutoMigrate` addition (test-DB/production consistency, not the real migration mechanism).
- `docs/module-recruitment-development-plan.md` — update G-6 status (documents done; consents still open, now the sole remaining G-6 item besides `status`/`source_id`).
- `docs/database-schema.md` — add `candidate_documents` to the table inventory and Mermaid ER block, following the convention re-established during sub-project 2's final review fix (a gap that slipped through in earlier sub-projects — this plan's Task list must explicitly include this file so no task/review boundary misses it again).
