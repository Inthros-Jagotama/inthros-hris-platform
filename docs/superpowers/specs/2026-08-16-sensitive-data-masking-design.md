# Sensitive Employee Data Masking & Encryption — Design

## Problem

Sensitive employee data (NIK, KK/family NIK, passport, phone number, email,
bank account number/name, emergency contact phone) is currently stored as
plain text in tenant databases and returned unmasked to any caller with
`employee.view`-level access. There is no way to (a) encrypt these values at
rest, or (b) restrict which roles can see the plain value versus a masked
one.

## Goals

- Let an admin toggle, per field, whether new writes to that field are
  encrypted at rest.
- Let an admin control, per role and per field, who can see the real value
  versus a masked value.
- Reuse the existing AES-256-GCM crypto utility and the existing
  resource.action RBAC model — no new crypto or permission engine.
- Support old plaintext rows and newly-encrypted rows coexisting (no
  backfill migration).

## Non-goals

- Backfilling/re-encrypting historical plaintext data.
- Arbitrary/free-text field selection — the set of maskable fields is a
  fixed, developer-maintained list; admins toggle within that list, they
  don't define new ones.
- Field-level encryption for modules outside employee/family/bank
  account/emergency contact in this iteration.

## Approach

Explicit encrypt/decrypt/mask calls in the service layer (not GORM hooks),
matching this codebase's existing preference for explicit, non-magic data
handling (tenant schema is already managed via versioned SQL migrations
rather than GORM AutoMigrate — see [[tenant-schema-migration-requirement]]).

Two independent axes:

1. **Encryption at rest** — controlled by a per-field on/off setting.
   Governs whether the *stored* value is ciphertext or plaintext.
2. **Masking on read** — controlled by per-field RBAC permissions. Governs
   whether the *value returned to a given caller* is the real value or a
   masked one. Applies regardless of whether the value happens to be
   encrypted at rest — a caller without permission gets a masked value even
   if the underlying storage is still plaintext.

These are deliberately decoupled: an admin can turn on masking-by-role
before ever enabling encryption, and enabling encryption doesn't
automatically grant/revoke view access.

## 1. Field registry

A fixed Go constant list of eligible fields, each with a machine key,
owning module, and bilingual label (ID/EN), consistent with the existing
bilingual RBAC module naming:

| field_key | module | model.field |
|---|---|---|
| `employee.nik` | employee | `Employee.NIK` |
| `employee.passport` | employee | `Employee.Passport` |
| `employee.phone_number` | employee | `Employee.PhoneNumber` |
| `employee.email` | employee | `Employee.Email` |
| `employee_family.nik` | employee | `EmployeeFamily.NIK` |
| `employee_bank_account.account_number` | employee | `EmployeeBankAccount.AccountNumber` |
| `employee_bank_account.account_name` | employee | `EmployeeBankAccount.AccountName` |
| `emergency_contact.phone_number` | employee | `EmergencyContact.PhoneNumber` |

This list lives in code (e.g. `backend/internal/modules/employee/sensitivefield/registry.go`), not the database — it's the catalog of *what can be
configured*, not the configuration itself.

## 2. Settings storage + admin UI

New tenant table `sensitive_field_settings`:

```sql
CREATE TABLE sensitive_field_settings (
    id BIGINT PRIMARY KEY,
    field_key VARCHAR(100) NOT NULL UNIQUE,
    is_encryption_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    updated_by BIGINT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

Seeded (via the same migration) with one row per registry entry, all
`is_encryption_enabled = false`.

New settings page, "Pengaturan Data Sensitif" / "Sensitive Data Settings",
lists fields grouped by module with a toggle per field. Toggling only
affects writes made after the change — existing rows are left untouched
until next edited (no backfill job).

## 3. Per-field RBAC permissions

For each registry entry, seed a permission `resource.action` pair —
`{module}.view_{field_suffix}` — e.g. `employee.view_nik`,
`employee_bank_account.view_account_number`. These are ordinary rows in
the existing `Permission` table, shown on the existing RBAC permissions
page (the same "one switch column per action" UI already shipped),
grouped under their module like any other permission.

Default grants at seed time: `super_admin` and `company_admin` get all
field-view permissions; other roles get none. Admins adjust from there
through the existing RBAC UI — no new permission-management UI is needed.

These permissions gate masking only, independent of the encryption toggle
in §2.

## 4. Column width migration

AES-256-GCM output, base64-encoded, is longer than the plaintext it
replaces (nonce + tag overhead, ~1.37x base64 expansion). Existing columns
sized for plaintext (e.g. `NIK varchar(16)`) can't hold ciphertext.

A versioned tenant SQL migration pair (mysql + postgres, under
`backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/`,
following the existing numeric-prefix pattern) widens every column backing
a registry entry to accommodate ciphertext (`varchar(255)` or `text`,
matched to worst-case field length) — applied up front for all eligible
columns, regardless of each field's current toggle state, so enabling
encryption later never requires a follow-up schema change. Per
[[tenant-schema-migration-requirement]], this must be a versioned SQL
migration; GORM AutoMigrate does not run for tenant tables.

## 5. Write path

In each affected service's create/update method (employee, employee
family, employee bank account, emergency contact), before calling the
repository:

```
for each registry field on this model:
    if sensitive_field_settings[field_key].is_encryption_enabled:
        value = crypto.EncryptString(value)
```

Settings are read once per request (small helper, short-TTL cache) to
avoid a DB round-trip per field per request.

## 6. Read path

After the repository fetch, for every registry field on the model:

```
if crypto.LooksEncrypted(value):
    value = crypto.DecryptString(value)
// else: leave as-is — old plaintext row
```

Then, at the DTO-mapping step, for each registry field:

```
if caller has permission {module}.view_{field}:
    dto.Field = value
else:
    dto.Field = mask.PartialMask(value)
```

`mask.PartialMask` (new package `backend/internal/pkg/mask/mask.go`)
returns the value with all but the last 4 characters replaced by `*`
(e.g. `3201xxxxxxxx5678` → `************5678`, `1234567890` →
`******7890`). Values of length ≤ 4 are fully masked (all `*`, same
length). The permission check uses the same JWT `permissions` claim /
`authz.Enforcer` fallback the existing middleware already uses — no new
auth mechanism.

The full decrypted value must never be serialized into a response for a
caller lacking the corresponding permission, even transiently — masking
happens before the DTO is returned, not as a post-processing step on the
client side.

## Data flow summary

```
Write:  request → service → [encrypt if enabled] → repository → DB
Read:   DB → repository → service → [decrypt if LooksEncrypted]
                                   → DTO mapping → [mask if no permission] → response
```

## Testing

- Unit: encrypt/decrypt round-trip; mixed plaintext+ciphertext rows decode
  correctly; `PartialMask` on various lengths including ≤4 chars.
- Unit: DTO mapping masks/unmasks correctly per permission, per field,
  independent of encryption toggle state.
- Migration test: widened columns accept both a plaintext-length and a
  ciphertext-length value on both mysql and postgres dialects.
- Settings toggle: enabling a field after some rows already exist leaves
  old rows readable (plaintext) and new writes encrypted.

## Open items for implementation planning

- Exact column width per field (compute from AES-GCM overhead + base64
  expansion for the longest current value in each column).
- Whether `updated_by`/`updated_at` audit trail on `sensitive_field_settings`
  needs to surface in an activity log (existing pattern elsewhere in the
  app, if any, should be followed for consistency).
