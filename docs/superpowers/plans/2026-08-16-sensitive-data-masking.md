# Sensitive Employee Data Masking & Encryption Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let admins toggle per-field encryption-at-rest for sensitive employee data (NIK, passport, phone, email, family NIK, bank account, emergency contact phone), and let admins control per-role/per-field whether a caller sees the real value or a partially masked one (last 4 chars visible).

**Architecture:** Explicit encrypt/decrypt/mask calls in the employee module's service layer (no GORM hooks), reusing the existing AES-256-GCM `crypto` package and the existing `resource.action` RBAC model. Two independent axes: a `sensitive_field_settings` table controls whether a field is encrypted at write time; per-field RBAC permissions (`employee.view_<field>` etc.) control whether a read-time value is unmasked for the caller.

**Tech Stack:** Go, Gin, GORM, MySQL/Postgres (dual-dialect tenant DBs), Vue 3 + PrimeVue (tenant frontend).

**Spec:** [docs/superpowers/specs/2026-08-16-sensitive-data-masking-design.md](../specs/2026-08-16-sensitive-data-masking-design.md)

## Global Constraints

- Masking format: last 4 characters visible, rest replaced with `*`; values of length ≤ 4 are fully masked to the same length (spec §6, revised in chat to partial mask).
- Encryption is **encrypt-on-write only** — no backfill of existing plaintext rows (spec Non-goals).
- The set of maskable fields is a fixed, developer-maintained registry — admins toggle within it, they don't define new fields (spec Non-goals).
- Tenant schema changes MUST be versioned SQL migrations under `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/` with matching numeric-prefix filenames — GORM AutoMigrate must not be relied on for these changes, per [[tenant-schema-migration-requirement]].
- Encryption key: env var `HRIS_ENCRYPTION_KEY` (64 hex chars), read via existing `backend/internal/pkg/crypto` package — do not add a second key-loading mechanism.
- Next available migration number is `150` (last existing is `149_rbac_submenu_permissions.sql`).

---

### Task 1: Partial-mask utility

**Files:**
- Create: `backend/internal/pkg/mask/mask.go`
- Test: `backend/internal/pkg/mask/mask_test.go`

**Interfaces:**
- Produces: `func PartialMask(value string) string` — used by Task 12.

- [ ] **Step 1: Write the failing test**

```go
package mask

import "testing"

func TestPartialMask(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"long value shows last 4", "3201010101985678", "************5678"},
		{"exact 5 chars", "12345", "*2345"},
		{"exact 4 chars fully masked", "1234", "****"},
		{"3 chars fully masked", "123", "***"},
		{"empty stays empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := PartialMask(c.input)
			if got != c.want {
				t.Errorf("PartialMask(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/internal/pkg/mask/... -run TestPartialMask -v`
Expected: FAIL — `undefined: PartialMask` (package doesn't exist yet)

- [ ] **Step 3: Write minimal implementation**

```go
// Package mask menyediakan utilitas untuk menyamarkan sebagian nilai
// data sensitif sebelum dikirim ke caller yang tidak punya izin melihat
// nilai asli.
package mask

// PartialMask mengganti semua karakter kecuali 4 karakter terakhir
// dengan '*'. Nilai dengan panjang <= 4 disamarkan penuh (semua '*'),
// dengan panjang yang sama seperti input.
func PartialMask(value string) string {
	runes := []rune(value)
	n := len(runes)
	if n == 0 {
		return ""
	}
	if n <= 4 {
		return repeatStar(n)
	}
	visibleStart := n - 4
	masked := make([]rune, n)
	for i := 0; i < visibleStart; i++ {
		masked[i] = '*'
	}
	copy(masked[visibleStart:], runes[visibleStart:])
	return string(masked)
}

func repeatStar(n int) string {
	stars := make([]rune, n)
	for i := range stars {
		stars[i] = '*'
	}
	return string(stars)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/internal/pkg/mask/... -run TestPartialMask -v`
Expected: PASS (all 5 subtests)

- [ ] **Step 5: Commit**

```bash
git add backend/internal/pkg/mask/mask.go backend/internal/pkg/mask/mask_test.go
git commit -m "feat(mask): add PartialMask utility for sensitive field display"
```

---

### Task 2: Propagate permission claims into request context + authctx helper

**Files:**
- Modify: `backend/internal/pkg/middleware/tenant.go:33-37`
- Modify: `backend/internal/pkg/authctx/authctx.go`
- Test: `backend/internal/pkg/authctx/authctx_test.go`

**Interfaces:**
- Consumes: nothing new — mirrors the existing `user_id`/`company_id` propagation pattern already in `tenant.go`.
- Produces: `func GetPermissions(ctx context.Context) []string` and `func HasPermission(ctx context.Context, resource, action string) bool` — used by Task 12.

**Context:** `middleware/auth.go:83` already does `c.Set("permissions", claims.Permissions)` on the gin.Context, but `tenant.go:33-37` (which builds the `context.Context` that flows into services via `authctx.GetUserID`/`GetCompanyID`) does not currently copy it. Service-layer code only ever sees `context.Context`, never `*gin.Context`, so without this the masking check in Task 12 has no way to know the caller's permissions.

- [ ] **Step 1: Write the failing test for authctx.GetPermissions/HasPermission**

```go
// backend/internal/pkg/authctx/authctx_test.go
package authctx

import (
	"context"
	"testing"
)

func TestGetPermissions(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee.view", "employee.view_nik"})
	got := GetPermissions(ctx)
	if len(got) != 2 || got[0] != "employee.view" {
		t.Fatalf("GetPermissions() = %v, want [employee.view employee.view_nik]", got)
	}
}

func TestGetPermissions_Missing(t *testing.T) {
	got := GetPermissions(context.Background())
	if len(got) != 0 {
		t.Fatalf("GetPermissions() = %v, want empty", got)
	}
}

func TestHasPermission(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee.view_nik"})
	if !HasPermission(ctx, "employee", "view_nik") {
		t.Error("HasPermission(employee, view_nik) = false, want true")
	}
	if HasPermission(ctx, "employee", "view_account_number") {
		t.Error("HasPermission(employee, view_account_number) = true, want false")
	}
}

func TestHasPermission_Wildcard(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"*"})
	if !HasPermission(ctx, "employee", "view_nik") {
		t.Error("HasPermission with wildcard = false, want true")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/internal/pkg/authctx/... -v`
Expected: FAIL — `undefined: GetPermissions`, `undefined: HasPermission`

- [ ] **Step 3: Implement GetPermissions and HasPermission**

Append to `backend/internal/pkg/authctx/authctx.go`:

```go
// GetPermissions extracts the caller's permission claims ("resource.action"
// strings) from the request context. Returns an empty slice if not found.
// permissions diset oleh middleware AuthJWT dan di-propagate oleh
// middleware.TenantRequired ke request context.
func GetPermissions(ctx context.Context) []string {
	if perms, ok := ctx.Value("permissions").([]string); ok {
		return perms
	}
	return []string{}
}

// HasPermission mengecek apakah caller punya permission "resource.action"
// tertentu, termasuk dukungan wildcard "*".
func HasPermission(ctx context.Context, resource, action string) bool {
	required := resource + "." + action
	for _, p := range GetPermissions(ctx) {
		if p == "*" || p == required {
			return true
		}
	}
	return false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/internal/pkg/authctx/... -v`
Expected: PASS (all 4 tests)

- [ ] **Step 5: Propagate permissions in tenant.go**

Read `backend/internal/pkg/middleware/tenant.go` around lines 30-38 first to confirm exact surrounding code, then modify:

```go
// Before (existing):
ctx := context.WithValue(c.Request.Context(), "company_id", companyID)
if userID, ok := c.Get("user_id"); ok {
	ctx = context.WithValue(ctx, "user_id", userID)
}
c.Request = c.Request.WithContext(ctx)

// After:
ctx := context.WithValue(c.Request.Context(), "company_id", companyID)
if userID, ok := c.Get("user_id"); ok {
	ctx = context.WithValue(ctx, "user_id", userID)
}
if perms, ok := c.Get("permissions"); ok {
	if permsSlice, ok := perms.([]string); ok {
		ctx = context.WithValue(ctx, "permissions", permsSlice)
	}
}
c.Request = c.Request.WithContext(ctx)
```

Keep the exact surrounding structure/formatting of the file as read — only add the new `if perms, ok := ...` block in the same style as the existing `user_id` block.

- [ ] **Step 6: Manually verify no existing tests break**

Run: `go test ./backend/internal/pkg/middleware/... -v`
Expected: PASS (no regressions — this only adds a new context key, doesn't change existing behavior)

- [ ] **Step 7: Commit**

```bash
git add backend/internal/pkg/authctx/authctx.go backend/internal/pkg/authctx/authctx_test.go backend/internal/pkg/middleware/tenant.go
git commit -m "feat(authctx): propagate permission claims into request context"
```

---

### Task 3: Migration 150 — create and seed `sensitive_field_settings` table

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/150_sensitive_field_settings.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/150_sensitive_field_settings.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/150_sensitive_field_settings.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/150_sensitive_field_settings.down.sql`

**Interfaces:**
- Produces: table `sensitive_field_settings(id, field_key, is_encryption_enabled, updated_by, updated_at)`, one seeded row per registry key below — used by Task 6's repository.

The registry (fixed field keys, matches Task 5's model struct fields):
`employee.nik`, `employee.passport`, `employee.phone_number`, `employee.email`, `employee_family.nik`, `employee_bank_account.account_number`, `employee_bank_account.account_name`, `emergency_contact.phone_number`.

- [ ] **Step 1: Write mysql up migration**

```sql
-- =============================================================================
-- 150_sensitive_field_settings.sql
-- Sensitive Data Masking — tabel setting toggle enkripsi per field.
-- Setiap baris mewakili satu field sensitif yang bisa di-enkripsi saat
-- ditulis (encrypt-on-write). Toggle ini independen dari permission
-- view per-field (lihat migration 152).
-- =============================================================================

CREATE TABLE sensitive_field_settings (
    id CHAR(36) PRIMARY KEY,
    field_key VARCHAR(100) NOT NULL,
    is_encryption_enabled TINYINT(1) NOT NULL DEFAULT 0,
    updated_by CHAR(36) NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uq_sensitive_field_settings_key (field_key)
);

INSERT IGNORE INTO sensitive_field_settings (id, field_key, is_encryption_enabled, updated_at) VALUES
    ('a3f1b2c4-0001-5f1a-9c1e-000000000001', 'employee.nik', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000002', 'employee.passport', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000003', 'employee.phone_number', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000004', 'employee.email', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000005', 'employee_family.nik', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000006', 'employee_bank_account.account_number', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000007', 'employee_bank_account.account_name', 0, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000008', 'emergency_contact.phone_number', 0, CURRENT_TIMESTAMP);
```

- [ ] **Step 2: Write mysql down migration**

```sql
-- 150_sensitive_field_settings.down.sql
DROP TABLE IF EXISTS sensitive_field_settings;
```

- [ ] **Step 3: Write postgres up migration**

```sql
-- =============================================================================
-- 150_sensitive_field_settings.sql (postgres)
-- Sensitive Data Masking — tabel setting toggle enkripsi per field.
-- =============================================================================

CREATE TABLE sensitive_field_settings (
    id CHAR(36) PRIMARY KEY,
    field_key VARCHAR(100) NOT NULL UNIQUE,
    is_encryption_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_by CHAR(36) NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO sensitive_field_settings (id, field_key, is_encryption_enabled, updated_at) VALUES
    ('a3f1b2c4-0001-5f1a-9c1e-000000000001', 'employee.nik', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000002', 'employee.passport', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000003', 'employee.phone_number', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000004', 'employee.email', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000005', 'employee_family.nik', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000006', 'employee_bank_account.account_number', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000007', 'employee_bank_account.account_name', FALSE, CURRENT_TIMESTAMP),
    ('a3f1b2c4-0001-5f1a-9c1e-000000000008', 'emergency_contact.phone_number', FALSE, CURRENT_TIMESTAMP)
ON CONFLICT (field_key) DO NOTHING;
```

- [ ] **Step 4: Write postgres down migration**

```sql
-- 150_sensitive_field_settings.down.sql (postgres)
DROP TABLE IF EXISTS sensitive_field_settings;
```

- [ ] **Step 5: Run migrator tests / a local tenant migration dry run**

Run: `go test ./backend/internal/pkg/migrator/... -v`
Expected: PASS. If the migrator test suite runs migrations against a real/test tenant DB, also confirm `sensitive_field_settings` has 8 rows after running: `SELECT COUNT(*) FROM sensitive_field_settings;` → `8`.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/150_sensitive_field_settings.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/150_sensitive_field_settings.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/150_sensitive_field_settings.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/150_sensitive_field_settings.down.sql
git commit -m "feat(migration): add sensitive_field_settings table (150)"
```

---

### Task 4: Migration 151 — widen sensitive columns for ciphertext

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/151_widen_sensitive_columns.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/151_widen_sensitive_columns.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/151_widen_sensitive_columns.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/151_widen_sensitive_columns.down.sql`

**Context:** AES-256-GCM output is hex-encoded `[12-byte nonce][ciphertext][16-byte tag]`, i.e. `(12 + len(plaintext) + 16) * 2` hex characters. For a 16-char NIK that's `(12+16+16)*2 = 88` chars — `varchar(16)` cannot hold it. Widen every column backing a registry field to `varchar(255)` up front (per spec §4), so enabling encryption later never needs a follow-up schema change.

- [ ] **Step 1: Write mysql up migration**

```sql
-- =============================================================================
-- 151_widen_sensitive_columns.sql
-- Sensitive Data Masking — perbesar kolom yang berpotensi diisi ciphertext
-- (AES-256-GCM hex-encoded lebih panjang dari plaintext aslinya).
-- Dijalankan terlepas dari status toggle enkripsi saat ini, supaya
-- mengaktifkan enkripsi nanti tidak perlu migrasi skema tambahan.
-- =============================================================================

ALTER TABLE employees
    MODIFY COLUMN nik VARCHAR(255) NULL,
    MODIFY COLUMN passport VARCHAR(255) NULL;
-- phone_number and email already varchar(255), no change needed.

ALTER TABLE employee_families
    MODIFY COLUMN nik VARCHAR(255) NULL;

ALTER TABLE employee_bank_accounts
    MODIFY COLUMN account_number VARCHAR(255) NOT NULL;
-- account_name already varchar(255), no change needed.

-- emergency_contacts.phone_number already varchar(50) -> widen for ciphertext.
ALTER TABLE emergency_contacts
    MODIFY COLUMN phone_number VARCHAR(255) NOT NULL;
```

- [ ] **Step 2: Write mysql down migration**

```sql
-- 151_widen_sensitive_columns.down.sql
-- NB: only safe to run down if no encrypted (longer) values were written
-- while the wider columns were in place — this truncates on rollback.

ALTER TABLE employees
    MODIFY COLUMN nik VARCHAR(16) NULL,
    MODIFY COLUMN passport VARCHAR(50) NULL;

ALTER TABLE employee_families
    MODIFY COLUMN nik VARCHAR(16) NULL;

ALTER TABLE employee_bank_accounts
    MODIFY COLUMN account_number VARCHAR(50) NOT NULL;

ALTER TABLE emergency_contacts
    MODIFY COLUMN phone_number VARCHAR(50) NOT NULL;
```

- [ ] **Step 3: Write postgres up migration**

```sql
-- 151_widen_sensitive_columns.sql (postgres)

ALTER TABLE employees
    ALTER COLUMN nik TYPE VARCHAR(255),
    ALTER COLUMN passport TYPE VARCHAR(255);

ALTER TABLE employee_families
    ALTER COLUMN nik TYPE VARCHAR(255);

ALTER TABLE employee_bank_accounts
    ALTER COLUMN account_number TYPE VARCHAR(255);

ALTER TABLE emergency_contacts
    ALTER COLUMN phone_number TYPE VARCHAR(255);
```

- [ ] **Step 4: Write postgres down migration**

```sql
-- 151_widen_sensitive_columns.down.sql (postgres)

ALTER TABLE employees
    ALTER COLUMN nik TYPE VARCHAR(16),
    ALTER COLUMN passport TYPE VARCHAR(50);

ALTER TABLE employee_families
    ALTER COLUMN nik TYPE VARCHAR(16);

ALTER TABLE employee_bank_accounts
    ALTER COLUMN account_number TYPE VARCHAR(50);

ALTER TABLE emergency_contacts
    ALTER COLUMN phone_number TYPE VARCHAR(50);
```

- [ ] **Step 5: Confirm exact current column definitions before applying**

Before running, re-read `backend/internal/modules/employee/model.go` lines 15-59, 97-119, 122-145, 253-270 (already captured in this plan's context) to double check no column name/type has drifted from what's listed above. Also grep the tenant migration history for the actual `CREATE TABLE` statements for `employees`, `employee_families`, `employee_bank_accounts`, `emergency_contacts` to confirm column names match exactly (e.g. confirm `account_number` and not `bank_account_number`):

Run: `grep -rn "CREATE TABLE employees\b" backend/internal/pkg/migrator/migrations/tenant/mysql/`

Adjust the ALTER statements above if any name differs from what was found.

- [ ] **Step 6: Run migrator tests**

Run: `go test ./backend/internal/pkg/migrator/... -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/151_widen_sensitive_columns.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/151_widen_sensitive_columns.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/151_widen_sensitive_columns.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/151_widen_sensitive_columns.down.sql
git commit -m "feat(migration): widen sensitive employee columns for ciphertext (151)"
```

---

### Task 5: Migration 152 — seed per-field RBAC permissions

**Files:**
- Create: `backend/internal/scripts/print_permission_uuids.go` (throwaway helper, deleted in Step 3)
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/152_sensitive_field_view_permissions.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/152_sensitive_field_view_permissions.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/152_sensitive_field_view_permissions.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/152_sensitive_field_view_permissions.down.sql`

**Interfaces:**
- Produces: 8 new rows in `permissions` (`employee.view_nik`, `employee.view_passport`, `employee.view_phone_number`, `employee.view_email`, `employee_family.view_nik`, `employee_bank_account.view_account_number`, `employee_bank_account.view_account_name`, `emergency_contact.view_phone_number`), each linked to the Admin role via `role_has_permissions` — consumed by Task 12's `authctx.HasPermission` checks and shown on the existing RBAC permissions page.

**Context:** Existing permission IDs are deterministic: `uuid.NewSHA1(uuid.NameSpaceDNS, []byte("hris-permission-"+code))` (see `backend/internal/pkg/tenantseed/seed_data.go:32-34`, `codeToUUID("permission", name)`). The Admin role ID is also deterministic: `codeToUUID("role", "ADMIN")`. Compute both before writing the SQL so migration IDs match what `SeedTenantRBAC` would generate for new tenants (avoiding duplicate-but-different-ID rows).

- [ ] **Step 1: Write a throwaway script to compute the deterministic UUIDs**

```go
// backend/internal/scripts/print_permission_uuids.go
package main

import (
	"fmt"

	"github.com/google/uuid"
)

func codeToUUID(table, code string) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte("hris-"+table+"-"+code)).String()
}

func main() {
	codes := []string{
		"employee.view_nik",
		"employee.view_passport",
		"employee.view_phone_number",
		"employee.view_email",
		"employee_family.view_nik",
		"employee_bank_account.view_account_number",
		"employee_bank_account.view_account_name",
		"emergency_contact.view_phone_number",
	}
	for _, c := range codes {
		fmt.Printf("%-50s %s\n", c, codeToUUID("permission", c))
	}
	fmt.Println("role.ADMIN:", codeToUUID("role", "ADMIN"))
}
```

- [ ] **Step 2: Run it and record the output**

Run: `go run backend/internal/scripts/print_permission_uuids.go`

Copy the 8 printed permission UUIDs and the Admin role UUID — use these exact values in Step 4/5 below (do not invent placeholder UUIDs).

- [ ] **Step 3: Delete the throwaway script**

```bash
rm backend/internal/scripts/print_permission_uuids.go
```

- [ ] **Step 4: Write mysql up/down migrations using the computed UUIDs**

```sql
-- =============================================================================
-- 152_sensitive_field_view_permissions.sql
-- Sensitive Data Masking — permission per-field untuk melihat nilai asli
-- (bukan hasil masking). Default: hanya role Admin yang diberi akses;
-- role lain diatur manual lewat halaman RBAC.
-- ID deterministik sama persis dengan codeToUUID di SeedTenantRBAC
-- (uuid.NewSHA1), jadi aman: migrasi & re-seed tidak duplikat.
-- =============================================================================

INSERT IGNORE INTO permissions (id, name, guard_name, created_at, updated_at) VALUES
    ('<UUID_view_nik>', 'employee.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_passport>', 'employee.view_passport', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_phone_number>', 'employee.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_email>', 'employee.view_email', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_family_view_nik>', 'employee_family.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_bank_view_account_number>', 'employee_bank_account.view_account_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_bank_view_account_name>', 'employee_bank_account.view_account_name', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_emergency_view_phone_number>', 'emergency_contact.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT IGNORE INTO role_has_permissions (permission_id, role_id) VALUES
    ('<UUID_view_nik>', '<UUID_role_ADMIN>'),
    ('<UUID_view_passport>', '<UUID_role_ADMIN>'),
    ('<UUID_view_phone_number>', '<UUID_role_ADMIN>'),
    ('<UUID_view_email>', '<UUID_role_ADMIN>'),
    ('<UUID_family_view_nik>', '<UUID_role_ADMIN>'),
    ('<UUID_bank_view_account_number>', '<UUID_role_ADMIN>'),
    ('<UUID_bank_view_account_name>', '<UUID_role_ADMIN>'),
    ('<UUID_emergency_view_phone_number>', '<UUID_role_ADMIN>');
```

Replace every `<UUID_...>` placeholder with the exact value printed in Step 2 before saving the file — this migration must not be committed with placeholders still in it.

```sql
-- 152_sensitive_field_view_permissions.down.sql
DELETE FROM role_has_permissions WHERE permission_id IN (
    '<UUID_view_nik>', '<UUID_view_passport>', '<UUID_view_phone_number>', '<UUID_view_email>',
    '<UUID_family_view_nik>', '<UUID_bank_view_account_number>', '<UUID_bank_view_account_name>',
    '<UUID_emergency_view_phone_number>'
);
DELETE FROM permissions WHERE id IN (
    '<UUID_view_nik>', '<UUID_view_passport>', '<UUID_view_phone_number>', '<UUID_view_email>',
    '<UUID_family_view_nik>', '<UUID_bank_view_account_number>', '<UUID_bank_view_account_name>',
    '<UUID_emergency_view_phone_number>'
);
```

- [ ] **Step 5: Write postgres up/down migrations (same UUIDs, `ON CONFLICT DO NOTHING`)**

```sql
-- 152_sensitive_field_view_permissions.sql (postgres)

INSERT INTO permissions (id, name, guard_name, created_at, updated_at) VALUES
    ('<UUID_view_nik>', 'employee.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_passport>', 'employee.view_passport', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_phone_number>', 'employee.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_view_email>', 'employee.view_email', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_family_view_nik>', 'employee_family.view_nik', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_bank_view_account_number>', 'employee_bank_account.view_account_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_bank_view_account_name>', 'employee_bank_account.view_account_name', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('<UUID_emergency_view_phone_number>', 'emergency_contact.view_phone_number', 'web', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;

INSERT INTO role_has_permissions (permission_id, role_id) VALUES
    ('<UUID_view_nik>', '<UUID_role_ADMIN>'),
    ('<UUID_view_passport>', '<UUID_role_ADMIN>'),
    ('<UUID_view_phone_number>', '<UUID_role_ADMIN>'),
    ('<UUID_view_email>', '<UUID_role_ADMIN>'),
    ('<UUID_family_view_nik>', '<UUID_role_ADMIN>'),
    ('<UUID_bank_view_account_number>', '<UUID_role_ADMIN>'),
    ('<UUID_bank_view_account_name>', '<UUID_role_ADMIN>'),
    ('<UUID_emergency_view_phone_number>', '<UUID_role_ADMIN>')
ON CONFLICT (permission_id, role_id) DO NOTHING;
```

```sql
-- 152_sensitive_field_view_permissions.down.sql (postgres)
DELETE FROM role_has_permissions WHERE permission_id IN (
    '<UUID_view_nik>', '<UUID_view_passport>', '<UUID_view_phone_number>', '<UUID_view_email>',
    '<UUID_family_view_nik>', '<UUID_bank_view_account_number>', '<UUID_bank_view_account_name>',
    '<UUID_emergency_view_phone_number>'
);
DELETE FROM permissions WHERE id IN (
    '<UUID_view_nik>', '<UUID_view_passport>', '<UUID_view_phone_number>', '<UUID_view_email>',
    '<UUID_family_view_nik>', '<UUID_bank_view_account_number>', '<UUID_bank_view_account_name>',
    '<UUID_emergency_view_phone_number>'
);
```

- [ ] **Step 6: Run migrator tests**

Run: `go test ./backend/internal/pkg/migrator/... -v`
Expected: PASS. If a local tenant test DB is available, verify: `SELECT name FROM permissions WHERE name LIKE '%view_nik%' OR name LIKE '%view_account%' OR name LIKE '%view_passport%' OR name LIKE '%view_email%' OR name LIKE '%view_phone_number%';` returns 8 rows.

- [ ] **Step 7: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/152_sensitive_field_view_permissions.sql \
        backend/internal/pkg/migrator/migrations/tenant/mysql/152_sensitive_field_view_permissions.down.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/152_sensitive_field_view_permissions.sql \
        backend/internal/pkg/migrator/migrations/tenant/postgres/152_sensitive_field_view_permissions.down.sql
git commit -m "feat(migration): seed per-field sensitive-data view permissions (152)"
```

---

### Task 6: Field registry constants

**Files:**
- Create: `backend/internal/modules/employee/sensitive_field_registry.go`
- Test: `backend/internal/modules/employee/sensitive_field_registry_test.go`

**Interfaces:**
- Produces: `type SensitiveFieldDef struct { Key, Resource, Action string }`, `var SensitiveFieldRegistry []SensitiveFieldDef`, `func FieldDef(key string) (SensitiveFieldDef, bool)` — used by Tasks 7, 9, 10, 11, 12.

- [ ] **Step 1: Write the failing test**

```go
// backend/internal/modules/employee/sensitive_field_registry_test.go
package employee

import "testing"

func TestSensitiveFieldRegistry_HasEightEntries(t *testing.T) {
	if len(SensitiveFieldRegistry) != 8 {
		t.Fatalf("len(SensitiveFieldRegistry) = %d, want 8", len(SensitiveFieldRegistry))
	}
}

func TestFieldDef_Found(t *testing.T) {
	def, ok := FieldDef("employee.nik")
	if !ok {
		t.Fatal("FieldDef(employee.nik) not found")
	}
	if def.Resource != "employee" || def.Action != "view_nik" {
		t.Errorf("FieldDef(employee.nik) = %+v, want Resource=employee Action=view_nik", def)
	}
}

func TestFieldDef_NotFound(t *testing.T) {
	_, ok := FieldDef("employee.does_not_exist")
	if ok {
		t.Error("FieldDef(employee.does_not_exist) should not be found")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./backend/internal/modules/employee/... -run TestSensitiveFieldRegistry -v`
Expected: FAIL — `undefined: SensitiveFieldRegistry`

- [ ] **Step 3: Write the implementation**

```go
// backend/internal/modules/employee/sensitive_field_registry.go
package employee

// SensitiveFieldDef mendefinisikan satu field sensitif yang bisa
// di-toggle enkripsinya (lewat sensitive_field_settings) dan dibatasi
// akses lihatnya (lewat permission RBAC "Resource.Action").
type SensitiveFieldDef struct {
	// Key adalah field_key di tabel sensitive_field_settings,
	// format "<model>.<field>", contoh "employee.nik".
	Key string
	// Resource dan Action membentuk permission RBAC "Resource.Action"
	// yang mengontrol apakah caller boleh melihat nilai asli (bukan masked).
	Resource string
	Action   string
}

// SensitiveFieldRegistry adalah daftar tetap field sensitif yang didukung
// sistem. Admin hanya bisa toggle enkripsi field yang ada di daftar ini —
// daftar ini tidak bisa diubah lewat UI/API, hanya lewat kode.
var SensitiveFieldRegistry = []SensitiveFieldDef{
	{Key: "employee.nik", Resource: "employee", Action: "view_nik"},
	{Key: "employee.passport", Resource: "employee", Action: "view_passport"},
	{Key: "employee.phone_number", Resource: "employee", Action: "view_phone_number"},
	{Key: "employee.email", Resource: "employee", Action: "view_email"},
	{Key: "employee_family.nik", Resource: "employee_family", Action: "view_nik"},
	{Key: "employee_bank_account.account_number", Resource: "employee_bank_account", Action: "view_account_number"},
	{Key: "employee_bank_account.account_name", Resource: "employee_bank_account", Action: "view_account_name"},
	{Key: "emergency_contact.phone_number", Resource: "emergency_contact", Action: "view_phone_number"},
}

// FieldDef mencari definisi field sensitif berdasarkan field_key.
func FieldDef(key string) (SensitiveFieldDef, bool) {
	for _, d := range SensitiveFieldRegistry {
		if d.Key == key {
			return d, true
		}
	}
	return SensitiveFieldDef{}, false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./backend/internal/modules/employee/... -run TestSensitiveFieldRegistry -v -run TestFieldDef`
Expected: PASS (all 3 tests)

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/employee/sensitive_field_registry.go backend/internal/modules/employee/sensitive_field_registry_test.go
git commit -m "feat(employee): add sensitive field registry"
```

---

### Task 7: Sensitive field settings repository + service (encryption toggle)

**Files:**
- Create: `backend/internal/modules/employee/sensitive_field_settings.go`
- Test: `backend/internal/modules/employee/sensitive_field_settings_test.go`

**Interfaces:**
- Consumes: `SensitiveFieldRegistry`, `FieldDef` (Task 6); `Repository` (existing, has `dbResolver func(ctx context.Context) (*gorm.DB, error)` per `repository.go:11-15`).
- Produces: `type SensitiveFieldSetting struct{...}` (GORM model), `func (r *Repository) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error)`, `func (r *Repository) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool, updatedBy *uuid.UUID) error`, `func (s *Service) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error)`, `func (s *Service) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool) error`, `func (s *Service) IsFieldEncryptionEnabled(ctx context.Context, fieldKey string) (bool, error)` — used by Tasks 8, 9, 10.

- [ ] **Step 1: Write the failing tests**

```go
// backend/internal/modules/employee/sensitive_field_settings_test.go
package employee

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupSensitiveFieldTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	if err := db.AutoMigrate(&SensitiveFieldSetting{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	for _, d := range SensitiveFieldRegistry {
		db.Create(&SensitiveFieldSetting{ID: mustNewUUID(), FieldKey: d.Key, IsEncryptionEnabled: false})
	}
	return db
}

func TestListSensitiveFieldSettings(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	settings, err := repo.ListSensitiveFieldSettings(context.Background())
	if err != nil {
		t.Fatalf("ListSensitiveFieldSettings() error = %v", err)
	}
	if len(settings) != len(SensitiveFieldRegistry) {
		t.Fatalf("got %d settings, want %d", len(settings), len(SensitiveFieldRegistry))
	}
}

func TestSetSensitiveFieldEnabled(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	ctx := context.Background()

	if err := repo.SetSensitiveFieldEnabled(ctx, "employee.nik", true, nil); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	var setting SensitiveFieldSetting
	db.Where("field_key = ?", "employee.nik").First(&setting)
	if !setting.IsEncryptionEnabled {
		t.Error("expected employee.nik encryption to be enabled")
	}
}

func TestService_IsFieldEncryptionEnabled(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	enabled, err := svc.IsFieldEncryptionEnabled(ctx, "employee.nik")
	if err != nil {
		t.Fatalf("IsFieldEncryptionEnabled() error = %v", err)
	}
	if enabled {
		t.Error("expected employee.nik to default to disabled")
	}

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	enabled, err = svc.IsFieldEncryptionEnabled(ctx, "employee.nik")
	if err != nil {
		t.Fatalf("IsFieldEncryptionEnabled() error = %v", err)
	}
	if !enabled {
		t.Error("expected employee.nik to be enabled after toggling")
	}
}

func TestService_SetSensitiveFieldEnabled_UnknownKey(t *testing.T) {
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())

	err := svc.SetSensitiveFieldEnabled(context.Background(), "not.a.real.field", true)
	if err == nil {
		t.Error("expected error for unknown field key")
	}
}
```

Read `backend/internal/modules/employee/service_test.go` (or any existing `*_test.go` in this package) first to find how `testLogger()`-equivalent and `mustNewUUID()`-equivalent helpers are already defined in this package (a `zap.NewNop()` helper and a uuid helper likely already exist) — reuse them instead of redefining, adjusting the test code above to call the existing helper names if they differ.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestListSensitiveFieldSettings -run TestSetSensitiveFieldEnabled -run TestService_IsFieldEncryptionEnabled -run TestService_SetSensitiveFieldEnabled_UnknownKey -v`
Expected: FAIL — `undefined: SensitiveFieldSetting`

- [ ] **Step 3: Implement the model, repository methods, and service methods**

```go
// backend/internal/modules/employee/sensitive_field_settings.go
package employee

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SensitiveFieldSetting adalah baris toggle enkripsi-at-rest untuk satu
// field sensitif. Dibuat oleh migration SQL (150_sensitive_field_settings),
// bukan AutoMigrate — lihat catatan Migrate() di module.go.
type SensitiveFieldSetting struct {
	ID                  string    `gorm:"type:char(36);primaryKey" json:"id"`
	FieldKey            string    `gorm:"type:varchar(100);uniqueIndex" json:"field_key"`
	IsEncryptionEnabled bool      `gorm:"column:is_encryption_enabled" json:"is_encryption_enabled"`
	UpdatedBy           *string   `gorm:"type:char(36)" json:"updated_by,omitempty"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (SensitiveFieldSetting) TableName() string { return "sensitive_field_settings" }

// ListSensitiveFieldSettings mengembalikan seluruh baris setting field sensitif.
func (r *Repository) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error) {
	db, err := r.dbResolver(ctx)
	if err != nil {
		return nil, err
	}
	var settings []SensitiveFieldSetting
	if err := db.WithContext(ctx).Order("field_key").Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// SetSensitiveFieldEnabled meng-update toggle enkripsi satu field.
func (r *Repository) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool, updatedBy *uuid.UUID) error {
	db, err := r.dbResolver(ctx)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"is_encryption_enabled": enabled,
		"updated_at":            time.Now(),
	}
	if updatedBy != nil {
		id := updatedBy.String()
		updates["updated_by"] = id
	}
	result := db.WithContext(ctx).Model(&SensitiveFieldSetting{}).
		Where("field_key = ?", fieldKey).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("sensitive field setting not found: %s", fieldKey)
	}
	return nil
}

// ListSensitiveFieldSettings (Service) — passthrough untuk handler.
func (s *Service) ListSensitiveFieldSettings(ctx context.Context) ([]SensitiveFieldSetting, error) {
	return s.repo.ListSensitiveFieldSettings(ctx)
}

// SetSensitiveFieldEnabled mengubah toggle enkripsi untuk satu field,
// setelah memvalidasi field_key ada di SensitiveFieldRegistry.
func (s *Service) SetSensitiveFieldEnabled(ctx context.Context, fieldKey string, enabled bool) error {
	if _, ok := FieldDef(fieldKey); !ok {
		return fmt.Errorf("unknown sensitive field key: %s", fieldKey)
	}
	var updatedBy *uuid.UUID
	if uid := authctxGetUserID(ctx); uid != nil {
		updatedBy = uid
	}
	return s.repo.SetSensitiveFieldEnabled(ctx, fieldKey, enabled, updatedBy)
}

// IsFieldEncryptionEnabled mengecek apakah field tertentu sedang di-toggle
// aktif enkripsinya. Dipanggil sebelum menulis nilai field sensitif.
func (s *Service) IsFieldEncryptionEnabled(ctx context.Context, fieldKey string) (bool, error) {
	settings, err := s.repo.ListSensitiveFieldSettings(ctx)
	if err != nil {
		return false, err
	}
	for _, st := range settings {
		if st.FieldKey == fieldKey {
			return st.IsEncryptionEnabled, nil
		}
	}
	return false, fmt.Errorf("unknown sensitive field key: %s", fieldKey)
}
```

Read `backend/internal/modules/employee/service.go:1-30` first to see the exact existing import alias used for `authctx` (e.g. it may already be imported as `authctx "github.com/inthros/hris-platform/internal/pkg/authctx"` given `authctx.GetUserID(ctx)` is used at lines 75/213 per the plan's research). Replace the placeholder `authctxGetUserID(ctx)` call above with the real `authctx.GetUserID(ctx)` call using whatever import alias the file already uses — do not add a second import of the same package under a different name.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestListSensitiveFieldSettings -run TestSetSensitiveFieldEnabled -run TestService_IsFieldEncryptionEnabled -run TestService_SetSensitiveFieldEnabled_UnknownKey -v`
Expected: PASS (all 4 tests)

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/employee/sensitive_field_settings.go backend/internal/modules/employee/sensitive_field_settings_test.go
git commit -m "feat(employee): add sensitive field settings repository and service"
```

---

### Task 8: Sensitive field settings HTTP endpoints

**Files:**
- Modify: `backend/internal/modules/employee/handler.go` (append methods)
- Modify: `backend/internal/modules/employee/routes.go:7-66`
- Test: `backend/internal/modules/employee/handler_test.go` (or new `sensitive_field_settings_handler_test.go` if `handler_test.go` doesn't already exist — check first)

**Interfaces:**
- Consumes: `Service.ListSensitiveFieldSettings`, `Service.SetSensitiveFieldEnabled` (Task 7).
- Produces: `GET /api/v1/tenant/employees/settings/sensitive-fields`, `PUT /api/v1/tenant/employees/settings/sensitive-fields/:fieldKey` — used by Task 13 (frontend).

- [ ] **Step 1: Check existing handler test conventions**

Run: `ls backend/internal/modules/employee/*_test.go` and read one existing handler test (if any) to match the request-building/assertion style used in this package (e.g. `httptest.NewRecorder()` + `gin.CreateTestContext`).

- [ ] **Step 2: Write the failing handler tests**

```go
// Append to an existing employee handler test file, or create
// backend/internal/modules/employee/sensitive_field_settings_handler_test.go
// matching whatever pattern Step 1 found. Example assuming gin test context style:

package employee

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHandler_ListSensitiveFieldSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupSensitiveFieldTestDB(t)
	repo := NewRepository(func(ctx gin.Context) (interface{}, error) { return nil, nil }) // placeholder, replace with real resolver signature used elsewhere in this file
	_ = repo
	_ = db
	// NOTE: replace this test body with the exact handler-test scaffolding
	// found in Step 1 (gin.CreateTestContext + w := httptest.NewRecorder()),
	// wired to a Service backed by setupSensitiveFieldTestDB(t), asserting:
	// GET /settings/sensitive-fields -> 200, body is a JSON array of 8 objects
	// each with "field_key" and "is_encryption_enabled".
}

func TestHandler_SetSensitiveFieldEnabled(t *testing.T) {
	// Same scaffolding: PUT /settings/sensitive-fields/employee.nik
	// with body {"is_encryption_enabled": true} -> 200, then GET confirms
	// employee.nik now has is_encryption_enabled=true.
}

func TestHandler_SetSensitiveFieldEnabled_UnknownKey(t *testing.T) {
	// PUT /settings/sensitive-fields/not.a.field -> 400 or 404 with an
	// error message, not a 500.
}
```

Because this package's exact handler-test scaffolding (context builder, response assertion helpers) wasn't captured verbatim during planning, the implementer must replace the placeholder bodies above with real assertions using the pattern found in Step 1 before proceeding — do not leave the placeholder/no-op bodies in the committed test file.

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestHandler_ListSensitiveFieldSettings -run TestHandler_SetSensitiveFieldEnabled -v`
Expected: FAIL — handler methods don't exist yet

- [ ] **Step 4: Implement the handler methods**

Append to `backend/internal/modules/employee/handler.go` (matching this file's existing method style — read `handler.go:1-40` first for the exact receiver name and response-helper functions used, e.g. `response.OK(c, ...)` / `response.Error(c, ...)` conventions already in the file):

```go
// ListSensitiveFieldSettings menampilkan daftar field sensitif beserta
// status toggle enkripsinya. GET /employees/settings/sensitive-fields
func (h *Handler) ListSensitiveFieldSettings(c *gin.Context) {
	settings, err := h.service.ListSensitiveFieldSettings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

type setSensitiveFieldEnabledRequest struct {
	IsEncryptionEnabled bool `json:"is_encryption_enabled"`
}

// SetSensitiveFieldEnabled mengubah toggle enkripsi satu field.
// PUT /employees/settings/sensitive-fields/:fieldKey
func (h *Handler) SetSensitiveFieldEnabled(c *gin.Context) {
	fieldKey := c.Param("fieldKey")
	var req setSensitiveFieldEnabledRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SetSensitiveFieldEnabled(c.Request.Context(), fieldKey, req.IsEncryptionEnabled); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"field_key": fieldKey, "is_encryption_enabled": req.IsEncryptionEnabled})
}
```

Adjust the error-response calls (`c.JSON(http.StatusInternalServerError, ...)` etc.) to match whatever shared response helper (e.g. `response.Error(c, ...)`) the rest of `handler.go` actually uses, instead of raw `c.JSON`, if such a helper exists in this file.

- [ ] **Step 5: Register the routes**

Modify `backend/internal/modules/employee/routes.go`, adding inside the existing `emps := rg.Group("/employees")` block (after the Employment routes, before the closing brace at line 65):

```go
		// Sensitive field settings (encryption toggle, admin only)
		emps.GET("/settings/sensitive-fields", handler.ListSensitiveFieldSettings)
		emps.PUT("/settings/sensitive-fields/:fieldKey", handler.SetSensitiveFieldEnabled)
```

- [ ] **Step 6: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestHandler_ListSensitiveFieldSettings -run TestHandler_SetSensitiveFieldEnabled -v`
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add backend/internal/modules/employee/handler.go backend/internal/modules/employee/routes.go backend/internal/modules/employee/*sensitive_field_settings_handler_test.go
git commit -m "feat(employee): add sensitive field settings HTTP endpoints"
```

---

### Task 9: Encrypt-on-write for Employee fields (NIK, Passport, PhoneNumber, Email)

**Files:**
- Modify: `backend/internal/modules/employee/service.go` (Create at line 64, Update at line 203)
- Test: `backend/internal/modules/employee/service_test.go` (append)

**Interfaces:**
- Consumes: `Service.IsFieldEncryptionEnabled` (Task 7), `crypto.EncryptString` (`backend/internal/pkg/crypto/crypto.go:115`).
- Produces: `func (s *Service) encryptIfEnabled(ctx context.Context, fieldKey string, value *string) error` — a helper reused by Task 10 for the other three models.

- [ ] **Step 1: Write the failing test**

```go
// Append to service_test.go — set HRIS_ENCRYPTION_KEY env var for the test.
func TestCreate_EncryptsNIKWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa") // 64 hex chars
	db := setupServiceTestDB(t) // reuse whatever existing helper this file already has for Employee service tests
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	nik := "3201010101985678"
	resp, err := svc.Create(ctx, CreateEmployeeRequest{ /* fill required fields per this file's existing Create tests */ NIK: &nik })
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var stored Employee
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK == nik {
		t.Fatal("expected NIK to be stored encrypted, got plaintext or nil")
	}
	if !crypto.LooksEncrypted(*stored.NIK) {
		t.Errorf("stored NIK %q does not look encrypted", *stored.NIK)
	}
}

func TestCreate_StoresPlaintextWhenDisabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	db := setupServiceTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()
	// employee.nik defaults to disabled — no toggle call.

	nik := "3201010101985678"
	resp, err := svc.Create(ctx, CreateEmployeeRequest{NIK: &nik})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var stored Employee
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK != nik {
		t.Fatalf("expected NIK stored as plaintext %q, got %v", nik, stored.NIK)
	}
}
```

Read `backend/internal/modules/employee/service_test.go` first (specifically whatever `setupServiceTestDB`-equivalent helper and `CreateEmployeeRequest` required-fields already exist for existing Create tests) and adjust field names/required-fields in the test above to match — do not invent a request shape that doesn't compile against the real `CreateEmployeeRequest` struct in `dto.go`.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestCreate_EncryptsNIKWhenEnabled -run TestCreate_StoresPlaintextWhenDisabled -v`
Expected: FAIL — NIK is stored as plaintext regardless of the toggle (current behavior, no encryption wired up yet)

- [ ] **Step 3: Implement encryptIfEnabled helper and wire it into Create/Update**

Add near the top of `service.go` (or in `sensitive_field_settings.go` next to the other Service methods from Task 7):

```go
// encryptIfEnabled meng-enkripsi *value in-place jika toggle enkripsi
// untuk fieldKey aktif. Nilai nil/kosong tidak diproses.
func (s *Service) encryptIfEnabled(ctx context.Context, fieldKey string, value *string) error {
	if value == nil || *value == "" {
		return nil
	}
	enabled, err := s.IsFieldEncryptionEnabled(ctx, fieldKey)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	encrypted, err := crypto.EncryptString(*value)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", fieldKey, err)
	}
	*value = encrypted
	return nil
}
```

Read `backend/internal/modules/employee/service.go:64-110` (the `Create` method) exactly as it stands, then insert calls to `encryptIfEnabled` right before the repository call, for each of the four fields:

```go
	if err := s.encryptIfEnabled(ctx, "employee.nik", emp.NIK); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.passport", emp.Passport); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.phone_number", emp.PhoneNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee.email", emp.Email); err != nil {
		return nil, err
	}
```

Place this block after `emp.NIK`, `emp.Passport`, `emp.PhoneNumber`, `emp.Email` are assigned from the request (around the existing lines 78-79, 96-97, 105-110) and before the call to `s.repo.CreateEmployee(ctx, emp)`. Add the identical block (same four calls) into `Update` (line 203) after the fields are re-assigned from the update request and before `s.repo.UpdateEmployee(ctx, emp)`.

Add the import for `crypto`: `"github.com/inthros/hris-platform/internal/pkg/crypto"` to `service.go`'s import block if not already present.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestCreate_EncryptsNIKWhenEnabled -run TestCreate_StoresPlaintextWhenDisabled -v`
Expected: PASS

- [ ] **Step 5: Run the full employee package test suite to check for regressions**

Run: `go test ./backend/internal/modules/employee/... -v`
Expected: PASS (no existing Create/Update tests broken by the new calls, since encryption defaults to disabled and `encryptIfEnabled` no-ops on nil/empty values)

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/employee/service.go backend/internal/modules/employee/service_test.go
git commit -m "feat(employee): encrypt-on-write for NIK, passport, phone, email"
```

---

### Task 10: Encrypt-on-write for Family, Bank Account, Emergency Contact

**Files:**
- Modify: `backend/internal/modules/employee/service.go` (CreateFamily:493, UpdateFamily:528, CreateBank:914, UpdateBank:940, CreateEmergencyContact:418, UpdateEmergencyContact:447)
- Test: `backend/internal/modules/employee/service_test.go` (append)

**Interfaces:**
- Consumes: `encryptIfEnabled` (Task 9).

- [ ] **Step 1: Write the failing tests**

```go
// Append to service_test.go
func TestCreateBank_EncryptsAccountNumberWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	db := setupServiceTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee_bank_account.account_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	// Reuse whatever fixture this file already uses to get a valid employee ID
	// and CreateBankRequest shape (check existing CreateBank tests in this file).
	empID := createTestEmployee(t, svc, ctx) // adjust to existing helper name
	resp, err := svc.CreateBank(ctx, empID, CreateBankRequest{AccountNumber: "1234567890", AccountName: "Budi"})
	if err != nil {
		t.Fatalf("CreateBank() error = %v", err)
	}

	var stored EmployeeBankAccount
	db.First(&stored, "id = ?", resp.ID)
	if stored.AccountNumber == "1234567890" {
		t.Fatal("expected account_number to be stored encrypted")
	}
	if !crypto.LooksEncrypted(stored.AccountNumber) {
		t.Errorf("stored account_number %q does not look encrypted", stored.AccountNumber)
	}
}

func TestCreateFamily_EncryptsNIKWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	db := setupServiceTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee_family.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	empID := createTestEmployee(t, svc, ctx)
	nik := "3201010101985678"
	resp, err := svc.CreateFamily(ctx, empID, CreateFamilyRequest{NIK: &nik, Name: "Anak Pertama"})
	if err != nil {
		t.Fatalf("CreateFamily() error = %v", err)
	}

	var stored EmployeeFamily
	db.First(&stored, "id = ?", resp.ID)
	if stored.NIK == nil || *stored.NIK == nik {
		t.Fatal("expected family NIK to be stored encrypted")
	}
}

func TestCreateEmergencyContact_EncryptsPhoneWhenEnabled(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	db := setupServiceTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "emergency_contact.phone_number", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}

	empID := createTestEmployee(t, svc, ctx)
	resp, err := svc.CreateEmergencyContact(ctx, empID, CreateEmergencyContactRequest{Name: "Ibu", PhoneNumber: "081234567890"})
	if err != nil {
		t.Fatalf("CreateEmergencyContact() error = %v", err)
	}

	var stored EmergencyContact
	db.First(&stored, "id = ?", resp.ID)
	if stored.PhoneNumber == "081234567890" {
		t.Fatal("expected emergency contact phone to be stored encrypted")
	}
}
```

Adjust `CreateBankRequest`/`CreateFamilyRequest`/`CreateEmergencyContactRequest` field names and the `createTestEmployee` helper name to match what's actually in `dto.go` and any existing test helpers in this package — read them first.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestCreateBank_EncryptsAccountNumberWhenEnabled -run TestCreateFamily_EncryptsNIKWhenEnabled -run TestCreateEmergencyContact_EncryptsPhoneWhenEnabled -v`
Expected: FAIL

- [ ] **Step 3: Wire encryptIfEnabled into the six service methods**

In `CreateFamily` (line 493) and `UpdateFamily` (line 528), before the repository call:
```go
	if err := s.encryptIfEnabled(ctx, "employee_family.nik", family.NIK); err != nil {
		return nil, err
	}
```

In `CreateBank` (line 914) and `UpdateBank` (line 940), before the repository call — note `AccountNumber`/`AccountName` are non-pointer `string`, not `*string`, per the model (`repository.go` findings), so `encryptIfEnabled` (which takes `*string`) needs `&bank.AccountNumber`:
```go
	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_number", &bank.AccountNumber); err != nil {
		return nil, err
	}
	if err := s.encryptIfEnabled(ctx, "employee_bank_account.account_name", &bank.AccountName); err != nil {
		return nil, err
	}
```

In `CreateEmergencyContact` (line 418) and `UpdateEmergencyContact` (line 447), before the repository call (also non-pointer `string`):
```go
	if err := s.encryptIfEnabled(ctx, "emergency_contact.phone_number", &contact.PhoneNumber); err != nil {
		return nil, err
	}
```

Read each method's exact current body first (variable names may differ from `family`/`bank`/`contact` assumed above — match whatever local variable the method actually uses for the model instance being saved).

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestCreateBank_EncryptsAccountNumberWhenEnabled -run TestCreateFamily_EncryptsNIKWhenEnabled -run TestCreateEmergencyContact_EncryptsPhoneWhenEnabled -v`
Expected: PASS

- [ ] **Step 5: Run the full employee package test suite**

Run: `go test ./backend/internal/modules/employee/... -v`
Expected: PASS, no regressions

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/employee/service.go backend/internal/modules/employee/service_test.go
git commit -m "feat(employee): encrypt-on-write for family NIK, bank account, emergency contact"
```

---

### Task 11: Decrypt-on-read in DTO converters

**Files:**
- Modify: `backend/internal/modules/employee/dto.go` (lines 350, 365, 441, 501)
- Test: `backend/internal/modules/employee/dto_test.go` (create if it doesn't exist)

**Interfaces:**
- Consumes: `crypto.LooksEncrypted`, `crypto.DecryptString` (`backend/internal/pkg/crypto/crypto.go:132`, `:120`).
- Produces: converters now return decrypted plaintext regardless of storage state — consumed directly by Task 12's masking step, which runs after these converters.

**Context:** Old rows may still be plaintext (encrypt-on-write only, no backfill); new rows may be ciphertext. `crypto.LooksEncrypted` disambiguates so both are handled uniformly.

- [ ] **Step 1: Write the failing tests**

```go
// backend/internal/modules/employee/dto_test.go
package employee

import (
	"testing"

	"github.com/inthros/hris-platform/internal/pkg/crypto"
)

func TestToFamilyResponse_DecryptsEncryptedNIK(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	fam := &EmployeeFamily{NIK: &encrypted}

	resp := toFamilyResponse(fam)

	if resp.NIK != plain {
		t.Errorf("toFamilyResponse().NIK = %q, want %q", resp.NIK, plain)
	}
}

func TestToFamilyResponse_LeavesPlaintextNIKAsIs(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	fam := &EmployeeFamily{NIK: &plain}

	resp := toFamilyResponse(fam)

	if resp.NIK != plain {
		t.Errorf("toFamilyResponse().NIK = %q, want %q (unchanged plaintext)", resp.NIK, plain)
	}
}

func TestToBankResponse_DecryptsEncryptedAccountNumber(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	plain := "1234567890"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	bank := &EmployeeBankAccount{AccountNumber: encrypted, AccountName: "Budi"}

	resp := toBankResponse(bank)

	if resp.AccountNumber != plain {
		t.Errorf("toBankResponse().AccountNumber = %q, want %q", resp.AccountNumber, plain)
	}
}

func TestEmployeeToResponse_DecryptsEncryptedNIK(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	plain := "3201010101985678"
	encrypted, err := crypto.EncryptString(plain)
	if err != nil {
		t.Fatalf("EncryptString() error = %v", err)
	}
	emp := &Employee{NIK: &encrypted}

	resp := emp.ToResponse()

	if resp.NIK != plain {
		t.Errorf("ToResponse().NIK = %q, want %q", resp.NIK, plain)
	}
}
```

Read `dto.go` lines 350-380, 441-450, 501-520 first to confirm the exact response struct field names (`FamilyResponse.NIK`, `BankResponse.AccountNumber`, `EmployeeResponse.NIK` — assumed `string` not `*string` per the existing "if not nil, deref" pattern already seen at lines 370-372 and 510-512) before finalizing these tests.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestToFamilyResponse -run TestToBankResponse -run TestEmployeeToResponse -v`
Expected: FAIL — encrypted values pass through undecrypted (garbled ciphertext instead of plaintext)

- [ ] **Step 3: Add a decrypt helper and wire it into the four converters**

Add near the top of `dto.go`:

```go
// decryptIfLooksEncrypted mengembalikan nilai asli jika value terlihat
// seperti hasil enkripsi (crypto.LooksEncrypted), atau value apa adanya
// jika masih plaintext (data lama, encrypt-on-write belum menyentuhnya).
// Kegagalan decrypt tidak fatal — fallback ke value asli supaya response
// tidak error karena satu baris data lama yang rusak.
func decryptIfLooksEncrypted(value string) string {
	if value == "" || !crypto.LooksEncrypted(value) {
		return value
	}
	decrypted, err := crypto.DecryptString(value)
	if err != nil {
		return value
	}
	return decrypted
}
```

Then wrap each sensitive field assignment in the four converters:

`toEmergencyContactResponse` (line ~354): `PhoneNumber: decryptIfLooksEncrypted(c.PhoneNumber)`

`toFamilyResponse` (lines ~370-372): change
```go
if f.NIK != nil { r.NIK = *f.NIK }
```
to
```go
if f.NIK != nil { r.NIK = decryptIfLooksEncrypted(*f.NIK) }
```

`toBankResponse` (lines ~444-445): change `AccountNumber`/`AccountName` assignments to `decryptIfLooksEncrypted(b.AccountNumber)` / `decryptIfLooksEncrypted(b.AccountName)`.

`Employee.ToResponse()` (lines ~510-512 and onward for Passport/PhoneNumber/Email — read the rest of the method to find their exact assignment lines): wrap each with `decryptIfLooksEncrypted(...)` the same way as NIK.

Add `"github.com/inthros/hris-platform/internal/pkg/crypto"` to `dto.go`'s imports if not already present.

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestToFamilyResponse -run TestToBankResponse -run TestEmployeeToResponse -v`
Expected: PASS

- [ ] **Step 5: Run the full employee package test suite**

Run: `go test ./backend/internal/modules/employee/... -v`
Expected: PASS, no regressions

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/employee/dto.go backend/internal/modules/employee/dto_test.go
git commit -m "feat(employee): decrypt-on-read in DTO converters, tolerant of plaintext legacy rows"
```

---

### Task 12: Role-based masking on responses

**Files:**
- Modify: `backend/internal/modules/employee/service.go` (wrap every method that returns `EmployeeResponse`, `FamilyResponse`, `BankResponse`, `EmergencyContactResponse` to the caller)
- Test: `backend/internal/modules/employee/service_test.go` (append)

**Interfaces:**
- Consumes: `authctx.HasPermission` (Task 2), `mask.PartialMask` (Task 1), `FieldDef`/`SensitiveFieldRegistry` (Task 6).
- Produces: `func maskEmployeeResponse(ctx context.Context, r *EmployeeResponse)`, `func maskFamilyResponse(ctx context.Context, r *FamilyResponse)`, `func maskBankResponse(ctx context.Context, r *BankResponse)`, `func maskEmergencyContactResponse(ctx context.Context, r *EmergencyContactResponse)` — called at the end of every Create/Update/Get/List method in `service.go` that returns these DTOs.

- [ ] **Step 1: Write the failing tests**

```go
// Append to service_test.go
func TestMaskFamilyResponse_MasksWithoutPermission(t *testing.T) {
	ctx := context.Background() // no permissions in context
	resp := &FamilyResponse{NIK: "3201010101985678"}

	maskFamilyResponse(ctx, resp)

	if resp.NIK != "************5678" {
		t.Errorf("NIK = %q, want masked", resp.NIK)
	}
}

func TestMaskFamilyResponse_UnmaskedWithPermission(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee_family.view_nik"})
	resp := &FamilyResponse{NIK: "3201010101985678"}

	maskFamilyResponse(ctx, resp)

	if resp.NIK != "3201010101985678" {
		t.Errorf("NIK = %q, want unmasked plaintext", resp.NIK)
	}
}

func TestMaskBankResponse_MasksAccountNumberAndName(t *testing.T) {
	ctx := context.Background()
	resp := &BankResponse{AccountNumber: "1234567890", AccountName: "Budi Santoso"}

	maskBankResponse(ctx, resp)

	if resp.AccountNumber != "******7890" {
		t.Errorf("AccountNumber = %q, want masked", resp.AccountNumber)
	}
	if resp.AccountName != "**********oso" && resp.AccountName != "**********oso"[:len(resp.AccountName)] {
		// AccountName masking follows the same PartialMask rule as any other field.
	}
}

func TestMaskEmployeeResponse_PerFieldGranularity(t *testing.T) {
	ctx := context.WithValue(context.Background(), "permissions", []string{"employee.view_nik"}) // NIK only, not passport/phone/email
	resp := &EmployeeResponse{NIK: "3201010101985678", Passport: "A1234567", PhoneNumber: "081234567890", Email: "budi@example.com"}

	maskEmployeeResponse(ctx, resp)

	if resp.NIK != "3201010101985678" {
		t.Errorf("NIK should be unmasked, got %q", resp.NIK)
	}
	if resp.Passport == "A1234567" {
		t.Error("Passport should be masked")
	}
	if resp.PhoneNumber == "081234567890" {
		t.Error("PhoneNumber should be masked")
	}
}
```

Read `EmployeeResponse`, `FamilyResponse`, `BankResponse`, `EmergencyContactResponse` struct definitions in `dto.go` first to confirm field names/types are `string` (not `*string`) as assumed — adjust the test literals if any field is actually a pointer.

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./backend/internal/modules/employee/... -run TestMaskFamilyResponse -run TestMaskBankResponse -run TestMaskEmployeeResponse -v`
Expected: FAIL — `undefined: maskFamilyResponse`

- [ ] **Step 3: Implement the four masking functions**

Add to `service.go` (or a new `sensitive_field_masking.go` in the same package — prefer the new file, since `service.go` is already large per the repository's file-size norms):

```go
// backend/internal/modules/employee/sensitive_field_masking.go
package employee

import (
	"context"

	"github.com/inthros/hris-platform/internal/pkg/authctx"
	"github.com/inthros/hris-platform/internal/pkg/mask"
)

// maskField menyamarkan value jika caller tidak punya permission
// "resource.action" untuk field terkait.
func maskField(ctx context.Context, resource, action string, value *string) {
	if *value == "" {
		return
	}
	if authctx.HasPermission(ctx, resource, action) {
		return
	}
	*value = mask.PartialMask(*value)
}

func maskEmployeeResponse(ctx context.Context, r *EmployeeResponse) {
	maskField(ctx, "employee", "view_nik", &r.NIK)
	maskField(ctx, "employee", "view_passport", &r.Passport)
	maskField(ctx, "employee", "view_phone_number", &r.PhoneNumber)
	maskField(ctx, "employee", "view_email", &r.Email)
}

func maskFamilyResponse(ctx context.Context, r *FamilyResponse) {
	maskField(ctx, "employee_family", "view_nik", &r.NIK)
}

func maskBankResponse(ctx context.Context, r *BankResponse) {
	maskField(ctx, "employee_bank_account", "view_account_number", &r.AccountNumber)
	maskField(ctx, "employee_bank_account", "view_account_name", &r.AccountName)
}

func maskEmergencyContactResponse(ctx context.Context, r *EmergencyContactResponse) {
	maskField(ctx, "emergency_contact", "view_phone_number", &r.PhoneNumber)
}
```

- [ ] **Step 4: Run the new masking-function tests to verify they pass**

Run: `go test ./backend/internal/modules/employee/... -run TestMaskFamilyResponse -run TestMaskBankResponse -run TestMaskEmployeeResponse -v`
Expected: PASS

- [ ] **Step 5: Wire masking into every service method returning these DTOs**

Read `service.go` fully to enumerate every method returning `*EmployeeResponse`, `*FamilyResponse`, `*BankResponse`, `*EmergencyContactResponse` (at minimum: `Create`, `Update`, `GetByID`, `List` for Employee; `CreateFamily`, `UpdateFamily`, and any `GetFamily`/`ListFamilies`; `CreateBank`, `UpdateBank`, and any list/get variant; `CreateEmergencyContact`, `UpdateEmergencyContact`, and any list/get variant). For each, call the matching mask function on the response immediately before `return resp, nil` (or before appending to a list, for List-style methods):

```go
	maskEmployeeResponse(ctx, resp)
	return resp, nil
```

For `List`-style methods returning `[]EmployeeResponse` or similar, mask each element in the loop that builds the response slice.

- [ ] **Step 6: Add an integration-style test proving the full pipeline (encrypt → decrypt → mask)**

```go
// Append to service_test.go
func TestGetByID_FullPipeline_EncryptDecryptMask(t *testing.T) {
	t.Setenv("HRIS_ENCRYPTION_KEY", "0000000000000000000000000000000000000000000000000000000000aa")
	db := setupServiceTestDB(t)
	repo := NewRepository(func(ctx context.Context) (*gorm.DB, error) { return db, nil })
	svc := NewService(repo, testLogger())
	ctx := context.Background()

	if err := svc.SetSensitiveFieldEnabled(ctx, "employee.nik", true); err != nil {
		t.Fatalf("SetSensitiveFieldEnabled() error = %v", err)
	}
	nik := "3201010101985678"
	created, err := svc.Create(ctx, CreateEmployeeRequest{NIK: &nik})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Caller without view_nik permission: masked.
	noPermCtx := context.Background()
	got, err := svc.GetByID(noPermCtx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.NIK != "************5678" {
		t.Errorf("GetByID without permission: NIK = %q, want masked", got.NIK)
	}

	// Caller with view_nik permission: real value, even though it's
	// encrypted at rest.
	permCtx := context.WithValue(context.Background(), "permissions", []string{"employee.view_nik"})
	got, err = svc.GetByID(permCtx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.NIK != nik {
		t.Errorf("GetByID with permission: NIK = %q, want plaintext %q", got.NIK, nik)
	}
}
```

- [ ] **Step 7: Run the full employee package test suite**

Run: `go test ./backend/internal/modules/employee/... -v`
Expected: PASS, including the new integration test and no regressions in any pre-existing test

- [ ] **Step 8: Commit**

```bash
git add backend/internal/modules/employee/sensitive_field_masking.go backend/internal/modules/employee/service.go backend/internal/modules/employee/service_test.go
git commit -m "feat(employee): apply per-field role-based masking to sensitive field responses"
```

---

### Task 13: Frontend settings page — "Pengaturan Data Sensitif"

**Files:**
- Create: `frontend/tenant/src/views/settings/SensitiveFieldSettings.vue`
- Modify: `frontend/tenant/src/router/index.js` (add route)
- Modify: `frontend/tenant/src/views/settings/SettingsIndex.vue` (add nav entry, matching how the existing RBAC entry is wired there)
- Modify: `frontend/tenant/src/locales/en.json`, `frontend/tenant/src/locales/id.json` (add `sensitive_field.*` keys)

**Interfaces:**
- Consumes: `GET /api/v1/tenant/employees/settings/sensitive-fields`, `PUT /api/v1/tenant/employees/settings/sensitive-fields/:fieldKey` (Task 8).

- [ ] **Step 1: Read existing conventions**

Read `frontend/tenant/src/views/settings/RolesPermissions.vue` in full (already partially read: uses `api.get`/`api.put` from an `api` module, PrimeVue `ToggleSwitch`, `DataTable`, `SkeletonTable`, `t()` i18n) and `frontend/tenant/src/router/index.js` around the existing `/admin/settings/rbac` route registration, and `SettingsIndex.vue`'s nav-item list, to match exact import paths and route/nav registration patterns before writing new code.

- [ ] **Step 2: Add locale keys**

Add to `frontend/tenant/src/locales/en.json` (find the top-level key used for `rbac.*`, e.g. `"rbac": {...}`, and add a sibling):
```json
"sensitive_field": {
  "title": "Sensitive Data Settings",
  "description": "Control which sensitive employee fields are encrypted when saved.",
  "field_name": "Field",
  "module": "Module",
  "encryption_enabled": "Encrypt at rest",
  "updated_toast": "Setting updated successfully"
}
```

Add to `frontend/tenant/src/locales/id.json`:
```json
"sensitive_field": {
  "title": "Pengaturan Data Sensitif",
  "description": "Atur field data karyawan mana yang dienkripsi saat disimpan.",
  "field_name": "Field",
  "module": "Modul",
  "encryption_enabled": "Enkripsi saat disimpan",
  "updated_toast": "Setting berhasil diperbarui"
}
```

- [ ] **Step 3: Write the Vue component**

```vue
<!-- frontend/tenant/src/views/settings/SensitiveFieldSettings.vue -->
<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('sensitive_field.title') }}</h2>
      <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('sensitive_field.description') }}</p>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="settings"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      sortField="field_key"
      :sortOrder="1"
    >
      <Column field="field_key" :header="t('sensitive_field.field_name')" sortable />
      <Column :header="t('sensitive_field.encryption_enabled')" style="width:160px">
        <template #body="{ data }">
          <ToggleSwitch
            :modelValue="data.is_encryption_enabled"
            @update:modelValue="(val) => toggleField(data.field_key, val)"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ToggleSwitch from 'primevue/toggleswitch'
import { useToast } from 'primevue/usetoast'
import SkeletonTable from '@/components/SkeletonTable.vue'
import api from '@/services/api'

const { t } = useI18n()
const toast = useToast()

const settings = ref([])
const loading = ref(true)
const skeletonColumns = [{ field: 'field_key' }, { field: 'is_encryption_enabled' }]

async function loadSettings() {
  loading.value = true
  try {
    const { data } = await api.get('/api/v1/tenant/employees/settings/sensitive-fields')
    settings.value = data
  } finally {
    loading.value = false
  }
}

async function toggleField(fieldKey, enabled) {
  const row = settings.value.find((s) => s.field_key === fieldKey)
  const previous = row?.is_encryption_enabled
  if (row) row.is_encryption_enabled = enabled
  try {
    await api.put(`/api/v1/tenant/employees/settings/sensitive-fields/${fieldKey}`, {
      is_encryption_enabled: enabled,
    })
    toast.add({ severity: 'success', summary: t('sensitive_field.updated_toast'), life: 2000 })
  } catch (err) {
    if (row) row.is_encryption_enabled = previous
    toast.add({ severity: 'error', summary: err?.response?.data?.error || 'Error', life: 3000 })
  }
}

onMounted(loadSettings)
</script>
```

Adjust the `SkeletonTable` and `api` import paths, and the `toast` API (`useToast`/`toast.add`), to match whatever the exact existing paths/APIs are in `RolesPermissions.vue` — read it in Step 1 before finalizing, since the import paths above are inferred from convention, not verbatim-confirmed for this exact file location.

- [ ] **Step 4: Register the route**

In `frontend/tenant/src/router/index.js`, find the existing route object for `/admin/settings/rbac` (component: `RolesPermissions.vue` or similar) and add a sibling route:

```js
{
  path: '/admin/settings/sensitive-fields',
  name: 'sensitive-field-settings',
  component: () => import('@/views/settings/SensitiveFieldSettings.vue'),
  meta: { requiresAuth: true, permission: 'employee.view' },
},
```

Match the exact `meta` shape (auth guard keys) used by the neighboring RBAC route — read it first rather than assuming `requiresAuth`/`permission` are the real key names.

- [ ] **Step 5: Add nav entry in SettingsIndex.vue**

Read `frontend/tenant/src/views/settings/SettingsIndex.vue` to find the list/array of settings nav items (likely includes an RBAC entry with icon/label/route), and add a new entry following the exact same shape:

```js
{ label: t('sensitive_field.title'), icon: 'pi pi-lock', route: '/admin/settings/sensitive-fields' }
```

- [ ] **Step 6: Manual verification in the browser**

Run the frontend dev server (check `frontend/tenant/package.json` for the dev script, typically `npm run dev`), log in as a tenant admin, navigate to Settings → the new "Pengaturan Data Sensitif" entry, and confirm:
- The 8 fields load with their current toggle state.
- Toggling a field calls the PUT endpoint (check network tab) and shows the success toast.
- Toggling a field, then refreshing the page, shows the new state persisted (confirms the GET reflects the PUT).

- [ ] **Step 7: Commit**

```bash
git add frontend/tenant/src/views/settings/SensitiveFieldSettings.vue \
        frontend/tenant/src/router/index.js \
        frontend/tenant/src/views/settings/SettingsIndex.vue \
        frontend/tenant/src/locales/en.json \
        frontend/tenant/src/locales/id.json
git commit -m "feat(frontend): add sensitive data settings page"
```

---

## Post-plan verification

After all 13 tasks are complete:

- [ ] Run the full backend test suite: `go test ./backend/... -v` — expect PASS, no regressions.
- [ ] Run migrations against a local/test tenant DB (both mysql and postgres if both are available) and confirm no errors, `sensitive_field_settings` has 8 rows, and 8 new permissions exist.
- [ ] In the RBAC permissions page (already-shipped UI), confirm the 8 new `employee*.view_*` / `employee_bank_account.view_*` / `emergency_contact.view_*` permissions appear grouped correctly and are togglable per role — no frontend change should be needed there since that page reads permissions dynamically from the API, but confirm this assumption holds by checking the page after running migration 152 against a test tenant.
- [ ] Manually verify end-to-end: enable encryption for `employee.nik`, create an employee with a NIK, confirm the DB row is ciphertext, confirm a user without `employee.view_nik` sees a masked NIK in the UI/API response, and a user with the permission sees the real value.
