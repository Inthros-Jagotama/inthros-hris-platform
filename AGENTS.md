# AGENTS.md

This file provides guidance to AI coding agents (Claude Code, Codex, etc.) when working with code in this repository.

## Project Overview

Multi-tenant HRIS (HR Information System) SaaS platform. Go backend (modular monolith, one MySQL/PostgreSQL database per tenant + a shared platform database) with two separate Vue 3 frontends: a tenant app and a platform admin app.

## Commands

All backend commands run from `backend/`.

```bash
# Build
make build                    # builds ./bin/server
make build-installer          # builds ./bin/installer (CLI provisioning tool)

# Run
make run                      # go run ./cmd/server --config ./config/config.yaml
make run-hot                  # hot reload via air

# Test
make test                     # go test -v -race -count=1 ./...  (needs cgo/gcc — often fails on Windows without one)
make test-short               # go test -short -count=1 ./...    (skips integration tests, no cgo needed — use this on Windows)
make test-pkg pkg=modules/attendance   # run tests for one package only
go test ./internal/modules/attendance/... -run TestName -v      # single test
make race-test                # full suite with -race inside a golang:1.24 Docker container (Windows-friendly alternative to `make test`)
make coverage                 # tests with coverage report

# Lint / vet
make lint                     # golangci-lint run ./...
make vet                      # go vet ./...
make tidy                     # go mod tidy && go mod verify

# Docs
make docs                     # regenerate docs/openapi-report.md (checks all routes are documented first)
make db-docs                  # regenerate docs/database-schema.md from SQL migrations

# Migrations & seed data
make migrate                  # run platform migrations
make seed                     # run seeders
./bin/installer migrate --company=<id> --config=./config/config.yaml   # run pending migrations for one tenant
```

Frontend commands run from `frontend/tenant/` or `frontend/platform-admin/` (separate Vite apps, separate `package.json`):

```bash
npm run dev        # vite dev server
npm run build       # vite build
npm run preview
```

## Architecture

### Modular monolith, module SDK contract

Every feature lives under `backend/internal/modules/<name>/` (tenant-scoped) or `backend/internal/platform/<name>/` (platform-scoped) and implements the `module.Module` interface (`internal/pkg/module/module.go`): `Info()`, `RegisterRoutes()`, `Migrate()`, `Seed()`, `Permissions()`. A module is self-contained: its own routes, its own migration files, its own seed data, its own RBAC permission list — all declared directly in that module's `module.go` (Go structs, not a config file) and wired up at startup in `cmd/server/main.go`.

Standard file layout inside a module package: `model.go` (GORM structs), `dto.go` (request/response types + `ToResponse()` converters), `repository.go` (DB access — GORM for typed queries, raw `db.Table().Joins()` for cross-module reads), `service.go` (business logic), `handler.go` (Gin handlers), `routes.go` (route registration), `module.go` (the `Module` interface implementation + constructor wiring), plus `*_test.go` files.

**Cross-module reads never import another module's package** — this avoids circular dependencies between tenant modules. Instead, a module queries another module's tables directly via `db.Table("other_table").Joins(...)` raw GORM queries (e.g. `attendance` reading `employees`/`organizations`/`zones` owned by other modules). When you need to read data from a foreign module's table, follow this raw-query pattern rather than adding an import.

### Multi-tenant database topology

One **platform database** (companies, licenses, platform users, module registry — `backend/internal/platform/`) plus one **tenant database per company** (all `backend/internal/modules/` data), provisioned on signup. `internal/pkg/database.Manager` resolves the right `*gorm.DB` per request based on the JWT/tenant context; `middleware.TenantRequired` populates `company_id` (and `authctx.GetCompanyID(ctx)`) into the request context. A module that needs platform data from a tenant-scoped repository takes an explicit `platformDB *gorm.DB` field via a `NewRepositoryWithPlatformDB(...)` constructor (see `internal/modules/setting/repository.go` or `internal/modules/attendance/repository.go` for the pattern) rather than assuming the primary `getDB(ctx)` is the platform DB.

Migrations live in `backend/internal/pkg/migrator/migrations/`, split into `platform/` (no dialect split — the same `.sql` file runs against MySQL and PostgreSQL, so avoid dialect-specific syntax like MySQL's `AFTER column`) and `tenant/{mysql,postgres}/` (dialect-duplicated). Every migration is a sequentially-numbered pair: `NNN_description.sql` + `NNN_description.down.sql`. **AutoMigrate is never used for tenant modules** — every schema change needs a versioned SQL migration in the correct dialect folder(s).

### RBAC

Database-backed, enforced via middleware reading `resource.action` permission strings from JWT claims (`authctx.GetPermissions`/`HasPermission`). Each module declares its own permissions in `manifest.yaml`; routes that need authorization wrap handlers with permission-check middleware (see `requireAttendanceSettings`/`requireAttendanceReport`-style helpers in each module's `routes.go`).

### Testing conventions

Tests use `github.com/glebarez/sqlite` (pure-Go, no cgo) in-memory databases, not a real MySQL/PostgreSQL connection — see each module's `helpers_test.go` for `setupTestDB()` / `newTestService()` / `newTestRepository()` patterns, including hand-rolled minimal schemas for raw-queried foreign tables (e.g. attendance's test DB defines a minimal `organizations`/`zones`/`employments` table set to satisfy its raw JOINs against those other modules' data). Reuse the existing helpers in a module's test file rather than inventing new DB setup — the pattern (fake platform DB via `setupTestPlatformDB()`, seeded company/zone/employment rows, context carrying `user_id`/`company_id`) is consistent across modules. The one integration test suite that needs real databases (`internal/pkg/migrator/migrator_integration_test.go`) requires `make db-race-up` (Docker MySQL + PostgreSQL) first — it will fail with connection-refused errors otherwise; that failure is expected in an environment without those containers running, not a regression.

### Timezone handling

Timestamps are stored as UTC everywhere. `internal/pkg/timezone` resolves a tenant's effective IANA timezone (Zone override → Company default → `Asia/Jakarta` fallback, limited to 3 Indonesian zones) for interpreting date boundaries (e.g. attendance "work date", lateness) and for display — never for writing timestamps. See `docs/flow/module-settings-flow.md` §9 for the full design.

### Frontend

Two independent Vue 3 + PrimeVue + Tailwind apps (`frontend/tenant/`, `frontend/platform-admin/`), each its own Vite project — no shared build, no monorepo tooling between them (`frontend/shared/` exists but is currently empty). Bilingual (ID/EN) via a custom `useI18n()` composable (`src/composables/useI18n.js`) backed by `src/locales/{en,id}.json` and a `useLanguage()` Pinia-less store (`src/stores/language.js`) — not vue-i18n. API calls go through `src/services/api.js` (axios instance); most views call `api.get/post/...` directly rather than a generated client.

## Documentation

- `README.md` — setup, full API endpoint reference, environment variables, module list.
- `docs/platform-architecture-design.md` — the primary architecture document (multi-tenancy, module SDK, provisioning, security).
- `docs/flow/module-*-flow.md` — one runbook per module explaining its business flow, entities, and endpoints. Read the relevant one before making non-trivial changes to a module.
- `docs/openapi-report.md` — generated; run `make docs` to refresh after adding/removing routes.

## Recommended Skills/Workflow (Claude Code)

For non-trivial changes, follow this sequence rather than editing ad hoc:

1. **`superpowers:brainstorming`** — required before any implementation. Classifies scope (bounded vs architectural) and gets explicit approval on the design before touching code.
2. **`superpowers:writing-plans`** + **`superpowers:subagent-driven-development`** — for architectural-scale features (new subsystem, cross-module changes): spec → implementation plan → per-task execution with review loops.
3. **`superpowers:using-git-worktrees`** — isolate implementation work for a plan so `main` isn't touched until the work is reviewed.
4. **`superpowers:test-driven-development`** — mandatory for Go changes in this repo. This codebase's test suite (in-memory sqlite, no cgo) is built around TDD-sized units; write the failing test first, using the existing `helpers_test.go` patterns in the module you're touching.
5. **`superpowers:finishing-a-development-branch`** — merge/PR/keep decision after work is verified, never skip the test-suite check first.
6. **`superpowers:systematic-debugging`** — for bugs/unexpected behavior, before proposing a fix.
7. **`code-review`** / **`security-review`** — worth running before merging changes that touch sensitive data (payroll, employee PII) or cross-tenant boundaries, given this is a multi-tenant SaaS handling HR/payroll data.

Not generally relevant to this repo: `frontend-design`/`design`/`dataviz` (UI already has an established PrimeVue + Tailwind pattern), `claude-api` (not an LLM-integrated app).
