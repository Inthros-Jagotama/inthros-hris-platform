# Candidate Consents (G-6 sub-project 3b) — Design

> Ref: `docs/module-recruitment-development-plan.md` §G-6 (partial — the last of the originally-planned G-6 sub-tables; `candidates.status`/`source_id` remain separately deferred, out of scope here).

## Goal

Add `candidate_consents`, an append-only audit log recording when a candidate's consent to personal-data processing (GDPR-style) was granted or revoked, and by whom (the recruiter/HR user who recorded it — this system has no public candidate-facing portal, so consent is documented by staff, not self-service).

## Explicitly out of scope

- Only ONE consent type: data-processing consent. No `consent_type` field, no support for separate marketing/background-check/etc. consent categories — YAGNI, nothing in the current scope needs more than one type, and adding a type enum now would be speculative.
- No IP address / user-agent capture — meaningless here since staff record consent on the candidate's behalf, not the candidate themselves via a self-service form (confirmed: no public candidate application portal exists anywhere in this codebase).
- No UPDATE or DELETE endpoints — this table is append-only by design (audit log). "Current consent status" is derived by reading the latest row, not stored as mutable state.
- No `candidates.consent_status` column (mentioned in the original, now-superseded G-6 plan text) — redundant with the append-only log; the log itself is the source of truth, avoiding the dual-source-of-truth problem the G-1 stage-history work already decided against for a similar reason.
- `candidates.status`, `source_id` — still not built, unrelated to this sub-project.
- No frontend changes.

## Data Model

### `candidate_consents` (new table, append-only)

| Column | Type | Notes |
|---|---|---|
| `id` | CHAR(36) PK | |
| `candidate_id` | CHAR(36) NOT NULL | FK → `candidates`, `ON DELETE CASCADE`, index |
| `action` | VARCHAR(20) NOT NULL | `GRANTED` or `REVOKED` — enforced via Gin `binding:"oneof=GRANTED REVOKED"`, not a DB constraint (matches module convention) |
| `notes` | TEXT NULL | free text, e.g. "signed consent form on file", "verbal during phone screen" |
| `changed_by` | CHAR(36) NULL | user id of the staff member who recorded this entry, read from `c.GetString("user_id")` in the handler — same pattern already used for `job_application_stage_histories.changed_by` (G-5) |
| `changed_at` | BIGINT NOT NULL | unix nano, set server-side at creation |
| `created_at` | TIMESTAMP | |

No `updated_at` — rows are never updated (append-only).

## API

Deliberately narrower than every other G-6 sub-resource — no update, no delete, matching the append-only design:

```
POST   /recruitment/candidates/:id/consents   (append a new GRANTED/REVOKED entry)
GET    /recruitment/candidates/:id/consents   (full history, ordered changed_at ASC)
```

No `PUT`/`DELETE /consents/:id` routes — there is nothing to update or delete in an audit log. If a mistaken entry needs correcting, the correction is itself a new row (e.g. `REVOKED` with a note explaining the entry above it was made in error), never a mutation of history.

The "current" consent state for a candidate is derivable by the caller as the last item in the `GET` response (ordered ascending) — no separate "current status" endpoint in this iteration; if that becomes a real UI need later, it's a cheap read-only addition on top of the same table.

## Error Handling

- Creating a consent entry for a non-existent `candidate_id` → error (existing pattern: `s.repo.FindCandidateByID` guard).
- `action` validated via `binding:"required,oneof=GRANTED REVOKED"` (unlike other G-6 sub-resources' optional-with-default enums, this one is required — there's no sensible default for "did they consent or not").
- No validation preventing consecutive duplicate actions (e.g. two `GRANTED` rows in a row) — the log records what staff entered; deduplication/state-machine validation is not required for an audit trail and would just get in the way of correcting mistakes via a new row.

## Testing Plan

- Repository: create + list (ordered) round-trip.
- Service: create with candidate-existence guard; list.
- Handler: create (valid action), create (invalid action → 400), list.

## Files Touched (summary — exact paths/code in the implementation plan)

- New migration **101** (mysql+postgres, up+down, idempotent).
- `backend/internal/modules/recruitment/model.go` — `CandidateConsent` struct.
- `backend/internal/modules/recruitment/repository.go` — `CreateCandidateConsent`, `ListCandidateConsents` (only 2 methods — no Find-by-id/Update/Delete needed since there's no update/delete API).
- `backend/internal/modules/recruitment/service.go` — `CreateCandidateConsent` (candidate-existence guard + `changed_by`/`changed_at` population), `ListCandidateConsents`.
- `backend/internal/modules/recruitment/dto.go` — request/response DTOs.
- `backend/internal/modules/recruitment/handler.go` + `routes.go` — 2 new endpoints; the create handler reads `c.GetString("user_id")` the same way the G-5 `changed_by` fix does.
- `backend/internal/modules/recruitment/module.go` — `AutoMigrate` addition.
- `docs/module-recruitment-development-plan.md` — update G-6 status; this closes out the LAST of the originally-planned G-6 sub-tables (only `candidates.status`/`source_id` remain open, both explicitly skipped by earlier design decisions rather than deferred-to-later).
- `docs/database-schema.md` — add `candidate_consents` to the table inventory and Mermaid ER block (named explicitly per the now-established convention from sub-project 3a).
