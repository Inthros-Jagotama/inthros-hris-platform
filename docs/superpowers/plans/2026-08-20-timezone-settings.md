# Timezone Settings (Company Default & Zone Override) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a company-level default timezone and an optional per-Zone override, with a resolver helper that other modules can use to determine "today" and format timestamps correctly for a tenant, instead of assuming server/UTC time everywhere.

**Architecture:** Two new nullable/required string columns (`companies.timezone`, `zones.timezone`) store IANA timezone names limited to 3 Indonesian zones. A new `backend/internal/pkg/timezone` package resolves the effective `*time.Location` for an organization (Zone override → Company default → `Asia/Jakarta` fallback), with an in-memory cache since only 3 values exist. Settings API/UI let admins configure both levels. As the first real consumer, the Attendance module's "today's attendance" dashboard query (which currently has no client-supplied local time to lean on) uses the resolver to determine the tenant's current calendar date server-side.

**Tech Stack:** Go, Gin, GORM (raw `db.Table().Joins()` queries for cross-module reads — no GORM preload relations exist between Organization/Zone), MySQL + PostgreSQL dual-dialect migrations, Vue 3 + PrimeVue frontend.

**Spec:** `docs/superpowers/specs/2026-08-20-timezone-settings-design.md`

## Global Constraints

- Only 3 timezone values are ever valid: `Asia/Jakarta` (WIB), `Asia/Makassar` (WITA), `Asia/Jayapura` (WIT). Enforced server-side; no free-text IANA input accepted yet.
- All timestamps remain stored as UTC in the DB — this plan adds interpretation/display logic only, never changes how timestamps are written.
- Tenant-DB migrations need both MySQL and PostgreSQL variants plus `.down.sql`, following the existing sequential-numbering convention (next number after `155` in `backend/internal/pkg/migrator/migrations/tenant/mysql/`, next number after `013` in `backend/internal/pkg/migrator/migrations/platform/`). No AutoMigrate.
- Follow the `numbering` package's validation pattern: package-level `map[string]bool` of allowed values + sentinel error + `switch` in the handler mapping to HTTP status (this codebase's established pattern for small fixed-value string fields, per `backend/internal/pkg/numbering/service.go:16-29,57-63`).

---

## File Structure

**New files:**
- `backend/internal/pkg/timezone/timezone.go` — `Resolve()` function + in-memory `*time.Location` cache + allowed-values map + sentinel errors.
- `backend/internal/pkg/timezone/timezone_test.go` — unit tests.
- `backend/internal/pkg/migrator/migrations/platform/014_add_companies_timezone.sql` + `.down.sql`
- `backend/internal/pkg/migrator/migrations/tenant/mysql/156_zone_timezone.sql` + `.down.sql`
- `backend/internal/pkg/migrator/migrations/tenant/postgres/156_zone_timezone.sql` + `.down.sql`
- `frontend/tenant/src/views/settings/CompanyTimezoneView.vue` (or a section added to an existing company-profile settings view — see Task 7)

**Modified files:**
- `backend/internal/platform/company/model.go` — add `Timezone` field to `Company` struct.
- `backend/internal/modules/setting/model.go` — add `Timezone *string` field to `Zone` struct.
- `backend/internal/modules/setting/dto.go` — add `Timezone` to Create/Update Zone DTOs and `ZoneResponse`.
- `backend/internal/modules/setting/service.go` — validate + persist `Timezone` in `CreateZone`/`UpdateZone`.
- `backend/internal/modules/setting/handler.go` + `routes.go` — new company-timezone endpoints.
- `backend/internal/modules/attendance/repository.go` — new method to resolve effective timezone for an organization.
- `backend/internal/modules/attendance/service.go` — "today's attendance" query uses resolved timezone instead of raw server date.
- `frontend/tenant/src/router/index.js`, `frontend/tenant/src/layouts/Sidebar.vue` — only if a new page/route is added (Task 7 decides this at implementation time based on where company profile settings currently live).

---

## Task 1: `timezone` package — resolver + validation

**Files:**
- Create: `backend/internal/pkg/timezone/timezone.go`
- Test: `backend/internal/pkg/timezone/timezone_test.go`

**Interfaces:**
- Produces: `timezone.Resolve(companyTimezone string, zoneTimezone *string) (*time.Location, error)`, `timezone.IsValid(tz string) bool`, `timezone.Allowed() []string`, `timezone.ErrInvalidTimezone`.

- [ ] **Step 1: Write the failing tests**

```go
package timezone

import "testing"

func TestIsValid(t *testing.T) {
	cases := map[string]bool{
		"Asia/Jakarta":  true,
		"Asia/Makassar": true,
		"Asia/Jayapura":  true,
		"Asia/Singapore": false,
		"":                false,
		"WIB":             false,
	}
	for tz, want := range cases {
		if got := IsValid(tz); got != want {
			t.Errorf("IsValid(%q) = %v, want %v", tz, got, want)
		}
	}
}

func TestResolve_ZoneOverrideWins(t *testing.T) {
	zoneTz := "Asia/Jayapura"
	loc, err := Resolve("Asia/Jakarta", &zoneTz)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jayapura" {
		t.Errorf("got %s, want Asia/Jayapura", loc.String())
	}
}

func TestResolve_FallsBackToCompanyWhenZoneNil(t *testing.T) {
	loc, err := Resolve("Asia/Makassar", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Makassar" {
		t.Errorf("got %s, want Asia/Makassar", loc.String())
	}
}

func TestResolve_FallsBackToJakartaWhenCompanyEmpty(t *testing.T) {
	loc, err := Resolve("", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("got %s, want Asia/Jakarta", loc.String())
	}
}

func TestResolve_InvalidCompanyTimezone(t *testing.T) {
	_, err := Resolve("Not/AZone", nil)
	if err != ErrInvalidTimezone {
		t.Errorf("got %v, want ErrInvalidTimezone", err)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/pkg/timezone/... -v`
Expected: FAIL — package `timezone` does not exist yet (build error).

- [ ] **Step 3: Implement the package**

```go
package timezone

import (
	"errors"
	"time"
)

var ErrInvalidTimezone = errors.New("invalid timezone")

var allowed = map[string]bool{
	"Asia/Jakarta":  true,
	"Asia/Makassar": true,
	"Asia/Jayapura":  true,
}

var cache = map[string]*time.Location{}

func Allowed() []string {
	return []string{"Asia/Jakarta", "Asia/Makassar", "Asia/Jayapura"}
}

func IsValid(tz string) bool {
	return allowed[tz]
}

func loadCached(tz string) (*time.Location, error) {
	if loc, ok := cache[tz]; ok {
		return loc, nil
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return nil, err
	}
	cache[tz] = loc
	return loc, nil
}

// Resolve determina la zona efectiva: zoneTimezone (jika tidak nil) menang
// atas companyTimezone; jika companyTimezone kosong, fallback ke Asia/Jakarta.
func Resolve(companyTimezone string, zoneTimezone *string) (*time.Location, error) {
	tz := companyTimezone
	if tz == "" {
		tz = "Asia/Jakarta"
	}
	if zoneTimezone != nil && *zoneTimezone != "" {
		tz = *zoneTimezone
	}
	if !IsValid(tz) {
		return nil, ErrInvalidTimezone
	}
	return loadCached(tz)
}
```

(Fix any stray non-English comment artifacts before committing — write comments in English or Indonesian consistently with the surrounding file; this codebase's comments are Indonesian, so keep it Indonesian.)

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./backend/internal/pkg/timezone/... -v`
Expected: PASS (all 5 tests)

- [ ] **Step 5: Commit**

```bash
git add backend/internal/pkg/timezone/
git commit -m "feat: tambah package timezone untuk resolusi zona waktu tenant"
```

---

## Task 2: Migration — `companies.timezone`

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/platform/014_add_companies_timezone.sql`
- Create: `backend/internal/pkg/migrator/migrations/platform/014_add_companies_timezone.down.sql`

**Interfaces:**
- Consumes: nothing (pure DB migration).
- Produces: `companies.timezone` column, `NOT NULL DEFAULT 'Asia/Jakarta'`, used by Task 4.

- [ ] **Step 1: Write the up migration**

```sql
-- 014_add_companies_timezone.sql
ALTER TABLE companies
    ADD COLUMN timezone VARCHAR(64) NOT NULL DEFAULT 'Asia/Jakarta' AFTER phone;
```

- [ ] **Step 2: Write the down migration**

```sql
-- 014_add_companies_timezone.down.sql
ALTER TABLE companies DROP COLUMN timezone;
```

- [ ] **Step 3: Run the migration locally and verify**

Run: `go run ./backend/cmd/migrate up` (or the project's existing migrate command — check `backend/cmd/` for the exact binary name used elsewhere in this repo before running)
Expected: migration `014_add_companies_timezone` applies without error; `DESCRIBE companies;` shows the new `timezone` column with default `Asia/Jakarta`.

- [ ] **Step 4: Run the down migration and verify rollback**

Run: `go run ./backend/cmd/migrate down 1`
Expected: `timezone` column removed; re-run `up` to leave the DB in the migrated state before continuing.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/platform/014_add_companies_timezone.sql backend/internal/pkg/migrator/migrations/platform/014_add_companies_timezone.down.sql
git commit -m "feat(db): tambah kolom companies.timezone"
```

---

## Task 3: Migration — `zones.timezone` (mysql + postgres)

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/156_zone_timezone.sql` + `.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/156_zone_timezone.sql` + `.down.sql`

**Interfaces:**
- Produces: `zones.timezone` column, nullable, used by Task 5.

- [ ] **Step 1: Write MySQL up/down**

```sql
-- tenant/mysql/156_zone_timezone.sql
ALTER TABLE zones
    ADD COLUMN timezone VARCHAR(64) NULL AFTER region;
```

```sql
-- tenant/mysql/156_zone_timezone.down.sql
ALTER TABLE zones DROP COLUMN timezone;
```

- [ ] **Step 2: Write PostgreSQL up/down**

```sql
-- tenant/postgres/156_zone_timezone.sql
ALTER TABLE zones
    ADD COLUMN timezone VARCHAR(64) NULL;
```

```sql
-- tenant/postgres/156_zone_timezone.down.sql
ALTER TABLE zones DROP COLUMN timezone;
```

- [ ] **Step 3: Run migrations on a tenant DB (both dialects if both are configured locally) and verify**

Run: whatever tenant-migration command this repo uses (check `backend/cmd/` — likely `go run ./backend/cmd/migrate tenant up` or similar; confirm exact invocation from existing docs/README before running)
Expected: `zones.timezone` column present, nullable, no default.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/156_zone_timezone.sql backend/internal/pkg/migrator/migrations/tenant/mysql/156_zone_timezone.down.sql backend/internal/pkg/migrator/migrations/tenant/postgres/156_zone_timezone.sql backend/internal/pkg/migrator/migrations/tenant/postgres/156_zone_timezone.down.sql
git commit -m "feat(db): tambah kolom zones.timezone (override opsional)"
```

---

## Task 4: `Company.Timezone` field + company timezone API

**Files:**
- Modify: `backend/internal/platform/company/model.go` (add field near `Phone`, line ~44)
- Modify: `backend/internal/modules/setting/dto.go` (new DTOs)
- Modify: `backend/internal/modules/setting/service.go` (new service methods)
- Modify: `backend/internal/modules/setting/handler.go` (new handlers)
- Modify: `backend/internal/modules/setting/routes.go` (new routes)
- Test: `backend/internal/modules/setting/service_test.go` (create if it doesn't exist, following existing test file conventions in that package)

**Interfaces:**
- Consumes: `timezone.IsValid(tz string) bool` from Task 1.
- Produces: `GET /api/v1/tenant/settings/company/timezone`, `PUT /api/v1/tenant/settings/company/timezone`; `Service.GetCompanyTimezone(ctx) (string, error)`, `Service.UpdateCompanyTimezone(ctx, tz string) error`.

- [ ] **Step 1: Add `Timezone` field to `Company` struct**

Edit `backend/internal/platform/company/model.go`, add after the `Phone` field:

```go
	Timezone  string         `gorm:"type:varchar(64);not null;default:Asia/Jakarta" json:"timezone"`
```

- [ ] **Step 2: Write failing service test**

```go
// backend/internal/modules/setting/service_test.go
package setting

import (
	"context"
	"testing"
)

func TestUpdateCompanyTimezone_RejectsInvalidValue(t *testing.T) {
	s := newTestService(t) // reuse whatever test-service constructor already exists in this package's other _test.go files
	err := s.UpdateCompanyTimezone(context.Background(), "Asia/Singapore")
	if err != ErrInvalidCompanyTimezone {
		t.Errorf("got %v, want ErrInvalidCompanyTimezone", err)
	}
}
```

(If no `newTestService` helper exists yet in this package, check `backend/internal/modules/setting/*_test.go` for the actual existing pattern — e.g. sqlite in-memory GORM setup — and use that exact pattern instead of inventing a new one.)

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./backend/internal/modules/setting/... -run TestUpdateCompanyTimezone -v`
Expected: FAIL — `ErrInvalidCompanyTimezone` / `UpdateCompanyTimezone` undefined.

- [ ] **Step 4: Implement service methods**

Add to `backend/internal/modules/setting/service.go`:

```go
var ErrInvalidCompanyTimezone = errors.New("invalid company timezone")

func (s *Service) GetCompanyTimezone(ctx context.Context) (string, error) {
	companyID := tenantctx.CompanyID(ctx) // use whatever existing helper this codebase already uses elsewhere in this package to get the current company ID from context — check other Service methods in service.go for the exact call
	var tz string
	if err := s.platformDB.Model(&company.Company{}).
		Where("id = ?", companyID).
		Pluck("timezone", &tz).Error; err != nil {
		return "", err
	}
	return tz, nil
}

func (s *Service) UpdateCompanyTimezone(ctx context.Context, tz string) error {
	if !timezone.IsValid(tz) {
		return ErrInvalidCompanyTimezone
	}
	companyID := tenantctx.CompanyID(ctx)
	return s.platformDB.Model(&company.Company{}).
		Where("id = ?", companyID).
		Update("timezone", tz).Error
}
```

(Adjust `s.platformDB` and `tenantctx.CompanyID(ctx)` to match this package's actual existing field name and context-helper — grep `backend/internal/modules/setting/service.go` for how it currently accesses the platform DB / company ID, since `setting.Service` today only touches the tenant DB for Zone records; wiring a platform-DB handle into this service may require adding a field to the `Service` struct and its constructor — check `NewService(...)` and update its signature plus the module wiring in `backend/internal/modules/setting/module.go` accordingly.)

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./backend/internal/modules/setting/... -run TestUpdateCompanyTimezone -v`
Expected: PASS

- [ ] **Step 6: Add handlers and routes**

Add to `backend/internal/modules/setting/handler.go`:

```go
func (h *Handler) GetCompanyTimezone(c *gin.Context) {
	tz, err := h.svc.GetCompanyTimezone(c.Request.Context())
	if err != nil {
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.Success(c, gin.H{"timezone": tz})
}

type updateCompanyTimezoneRequest struct {
	Timezone string `json:"timezone" binding:"required"`
}

func (h *Handler) UpdateCompanyTimezone(c *gin.Context) {
	var req updateCompanyTimezoneRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	if err := h.svc.UpdateCompanyTimezone(c.Request.Context(), req.Timezone); err != nil {
		if errors.Is(err, ErrInvalidCompanyTimezone) {
			httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
		return
	}
	httputil.Success(c, gin.H{"timezone": req.Timezone})
}
```

(Match `httputil.Success`/`httputil.ErrorSimple`/`httputil.BindAndValidate` call shapes exactly to what's already used in the surrounding handlers in this file — e.g. `CreateZone` — rather than the illustrative shape above if signatures differ.)

Add to `backend/internal/modules/setting/routes.go`, inside the `settings` group:

```go
	settings.GET("/company/timezone", handler.GetCompanyTimezone)
	settings.PUT("/company/timezone", handler.UpdateCompanyTimezone)
```

- [ ] **Step 7: Manual verification**

Run the backend, then:
```bash
curl -X PUT http://localhost:8080/api/v1/tenant/settings/company/timezone -H "Authorization: Bearer <token>" -H "Content-Type: application/json" -d '{"timezone":"Asia/Makassar"}'
curl http://localhost:8080/api/v1/tenant/settings/company/timezone -H "Authorization: Bearer <token>"
```
Expected: PUT returns 200 with `{"timezone":"Asia/Makassar"}`; subsequent GET returns the same value; PUT with `"Asia/Singapore"` returns 400.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/platform/company/model.go backend/internal/modules/setting/
git commit -m "feat: tambah API get/update timezone default company"
```

---

## Task 5: `Zone.Timezone` override field + API

**Files:**
- Modify: `backend/internal/modules/setting/model.go` (add `Timezone *string` to `Zone` struct)
- Modify: `backend/internal/modules/setting/dto.go` (`CreateZoneRequest`, `UpdateZoneRequest`, `ZoneResponse`, `ToResponse()`)
- Modify: `backend/internal/modules/setting/service.go` (`CreateZone`, `UpdateZone`)
- Test: `backend/internal/modules/setting/service_test.go`

**Interfaces:**
- Consumes: `timezone.IsValid(tz string) bool` from Task 1.
- Produces: `Zone.Timezone *string` field consumed by Task 6's repository query.

- [ ] **Step 1: Write failing test**

```go
func TestCreateZone_RejectsInvalidTimezone(t *testing.T) {
	s := newTestService(t)
	tz := "Asia/Singapore"
	_, err := s.CreateZone(context.Background(), CreateZoneRequest{
		Code: "Z1", Name: "Zone 1", Timezone: &tz,
	})
	if err != ErrInvalidZoneTimezone {
		t.Errorf("got %v, want ErrInvalidZoneTimezone", err)
	}
}

func TestCreateZone_AllowsNilTimezone(t *testing.T) {
	s := newTestService(t)
	resp, err := s.CreateZone(context.Background(), CreateZoneRequest{Code: "Z2", Name: "Zone 2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Timezone != nil {
		t.Errorf("got %v, want nil", resp.Timezone)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/setting/... -run TestCreateZone -v`
Expected: FAIL — `Timezone` field / `ErrInvalidZoneTimezone` undefined.

- [ ] **Step 3: Add field to `Zone` model**

Edit `backend/internal/modules/setting/model.go`, add to the `Zone` struct after `Region`:

```go
	Timezone    *string        `gorm:"type:varchar(64)" json:"timezone,omitempty"`
```

- [ ] **Step 4: Add field to DTOs**

Edit `backend/internal/modules/setting/dto.go`:

```go
type CreateZoneRequest struct {
	Code      string  `json:"code" binding:"required,max=20"`
	Name      string  `json:"name" binding:"required,max=255"`
	Region    string  `json:"region,omitempty" binding:"max=100"`
	Timezone  *string `json:"timezone,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	SortOrder int     `json:"sort_order,omitempty"`
}

type UpdateZoneRequest struct {
	Code      *string `json:"code,omitempty" binding:"omitempty,max=20"`
	Name      *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Region    *string `json:"region,omitempty" binding:"omitempty,max=100"`
	Timezone  *string `json:"timezone,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	SortOrder *int    `json:"sort_order,omitempty"`
}

type ZoneResponse struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Region    string    `json:"region,omitempty"`
	Timezone  *string   `json:"timezone,omitempty"`
	IsActive  bool      `json:"is_active"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (z *Zone) ToResponse() ZoneResponse {
	return ZoneResponse{ID: z.ID.String(), Code: z.Code, Name: z.Name, Region: z.Region, Timezone: z.Timezone, IsActive: z.IsActive, SortOrder: z.SortOrder, CreatedAt: z.CreatedAt, UpdatedAt: z.UpdatedAt}
}
```

- [ ] **Step 5: Add validation + persistence in service**

Edit `backend/internal/modules/setting/service.go`:

```go
var ErrInvalidZoneTimezone = errors.New("invalid zone timezone")

func validateZoneTimezone(tz *string) error {
	if tz == nil || *tz == "" {
		return nil
	}
	if !timezone.IsValid(*tz) {
		return ErrInvalidZoneTimezone
	}
	return nil
}
```

In `CreateZone`, before `s.repo.CreateZone`:
```go
	if err := validateZoneTimezone(req.Timezone); err != nil {
		return nil, err
	}
	zone := &Zone{Code: req.Code, Name: req.Name, Zone: req.Name, Region: req.Region, Timezone: req.Timezone, IsActive: isActive, SortOrder: req.SortOrder}
```

In `UpdateZone`, alongside the other `if req.X != nil` blocks:
```go
	if req.Timezone != nil {
		if err := validateZoneTimezone(req.Timezone); err != nil {
			return nil, err
		}
		zone.Timezone = req.Timezone
	}
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/setting/... -run TestCreateZone -v`
Expected: PASS

- [ ] **Step 7: Wire validation errors to HTTP 400 in handler**

In `backend/internal/modules/setting/handler.go`, in `CreateZone` and `UpdateZone` error handling (wherever the existing `err != nil` branch calls `handleDupErr` or similar), add a check before the generic fallback:
```go
	if errors.Is(err, ErrInvalidZoneTimezone) {
		httputil.ErrorSimple(c, http.StatusBadRequest, err.Error())
		return
	}
```

- [ ] **Step 8: Commit**

```bash
git add backend/internal/modules/setting/
git commit -m "feat: tambah field timezone override pada Zone"
```

---

## Task 6: Attendance "today" resolution via effective timezone

**Files:**
- Modify: `backend/internal/modules/attendance/repository.go` (new method)
- Modify: `backend/internal/modules/attendance/service.go` (wherever "today's attendance" is queried — locate via grep for the dashboard/summary handler that filters by current date)
- Test: `backend/internal/modules/attendance/repository_test.go` or `service_test.go` (follow existing test file for this package)

**Interfaces:**
- Consumes: `timezone.Resolve(companyTimezone string, zoneTimezone *string) (*time.Location, error)` from Task 1.
- Produces: `Repository.ResolveOrganizationTimezone(ctx, organizationID uuid.UUID) (*time.Location, error)`.

- [ ] **Step 1: Write failing repository test**

```go
func TestResolveOrganizationTimezone_UsesZoneOverride(t *testing.T) {
	repo := newTestRepository(t) // reuse this package's existing test-DB setup helper
	// seed: a company with timezone "Asia/Jakarta", a zone with timezone "Asia/Jayapura",
	// an organization with that zone_id — follow this package's existing seeding helpers/fixtures
	loc, err := repo.ResolveOrganizationTimezone(context.Background(), orgID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loc.String() != "Asia/Jayapura" {
		t.Errorf("got %s, want Asia/Jayapura", loc.String())
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/internal/modules/attendance/... -run TestResolveOrganizationTimezone -v`
Expected: FAIL — method undefined.

- [ ] **Step 3: Implement repository method**

Add to `backend/internal/modules/attendance/repository.go`, following the existing raw `db.Table().Joins()` style used by `FindOrganizationIDByUserID` (repository.go:872):

```go
// ResolveOrganizationTimezone mengembalikan zona waktu efektif untuk sebuah
// organization: Zone.timezone (jika di-set) mengalahkan Company.timezone.
func (r *Repository) ResolveOrganizationTimezone(ctx context.Context, organizationID uuid.UUID) (*time.Location, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}

	var zoneTz *string
	err = db.Table("organizations AS o").
		Joins("LEFT JOIN zones AS z ON z.id = o.zone_id").
		Where("o.id = ?", organizationID).
		Select("z.timezone AS zone_tz").
		Scan(&zoneTz).Error
	if err != nil {
		return nil, err
	}

	companyTz, err := r.getCompanyTimezone(ctx) // new small helper: reads current tenant's company timezone from the platform DB — implement using whatever cross-DB access pattern this module already uses elsewhere to reach platform data (check how other modules read Company fields, if any; otherwise this may need to be passed in from a context value already populated by tenant-resolution middleware — check `backend/internal/middleware/` or similar for how company info reaches request context today, and use that instead of a fresh platform DB query if it's already available)
	if err != nil {
		return nil, err
	}

	return timezone.Resolve(companyTz, zoneTz)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/internal/modules/attendance/... -run TestResolveOrganizationTimezone -v`
Expected: PASS

- [ ] **Step 5: Wire into "today's attendance" query**

Locate the service method that computes "today" for a dashboard/summary listing (grep `backend/internal/modules/attendance/service.go` for a function filtering by the current date without a client-supplied date parameter — this is the one place in Attendance that has no client-local-time input to lean on, unlike `CreateEvent`'s `EventTimeLocal`). Replace its `time.Now()` date computation with:

```go
	loc, err := s.repo.ResolveOrganizationTimezone(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	today := time.Now().In(loc).Format("2006-01-02")
```

- [ ] **Step 6: Manual verification**

Set company timezone to `Asia/Jayapura` via the Task 4 API. Trigger the "today's attendance" dashboard endpoint close to midnight WIB (e.g. 23:30 WIB = 01:30 WIT next day) and confirm it reports the WIT calendar date, not the WIB/server date.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/attendance/
git commit -m "feat: attendance hari-ini mengikuti zona waktu efektif organization"
```

---

## Task 7: Frontend — company timezone & zone override settings UI

**Files:**
- Modify or create: a view under `frontend/tenant/src/views/settings/` (check whether a company-profile settings page already exists — if yes, add the timezone field there instead of creating a new page; if no, create `CompanyTimezoneView.vue` following the `CompanyHolidaysView.vue` pattern for API calls, i18n, and toast usage)
- Modify: the existing Zone create/edit form component under `frontend/tenant/src/views/settings/` (locate via grep for `CreateZoneRequest`/`zones` usage in the frontend) to add the override dropdown
- Modify: `frontend/tenant/src/router/index.js`, `frontend/tenant/src/layouts/Sidebar.vue` — only if Task 7 creates a new page

**Interfaces:**
- Consumes: `GET/PUT /api/v1/tenant/settings/company/timezone` (Task 4), Zone create/update endpoints with `timezone` field (Task 5).

- [ ] **Step 1: Add company timezone dropdown**

In the company settings view, add a PrimeVue `Dropdown` bound to a `timezone` ref with options:
```js
const timezoneOptions = [
  { label: 'WIB (Asia/Jakarta)', value: 'Asia/Jakarta' },
  { label: 'WITA (Asia/Makassar)', value: 'Asia/Makassar' },
  { label: 'WIT (Asia/Jayapura)', value: 'Asia/Jayapura' },
]
```
On mount, `GET` the current value; on save, `PUT` the selected value; show a toast on success/failure via `useToast()`, following the exact call pattern already used in `CompanyHolidaysView.vue` for its own save action.

- [ ] **Step 2: Add zone override dropdown**

In the Zone form, add the same `timezoneOptions` list plus a leading `{ label: 'Ikut default perusahaan', value: null }` option, bound to the zone's `timezone` field, submitted alongside `code`/`name`/`region` in the existing create/update payload.

- [ ] **Step 3: Manual verification in browser**

Start the frontend dev server. Navigate to company settings, change the timezone, save, reload the page, and confirm the saved value persists. Create/edit a Zone, set an override, save, and confirm the value round-trips (edit again and see the same value pre-selected). Set it back to "Ikut default perusahaan" and confirm the field clears (`null`) on the next `GET`.

- [ ] **Step 4: Commit**

```bash
git add frontend/tenant/src/
git commit -m "feat: UI setting timezone company & zone override"
```

---

## Not In This Plan (future phases, per spec's Rollout section)

- Payroll cutoff timezone integration.
- Leave/cuti date validation timezone integration.
- Broader audit of the ~60 files calling `time.Now()` across the backend.
