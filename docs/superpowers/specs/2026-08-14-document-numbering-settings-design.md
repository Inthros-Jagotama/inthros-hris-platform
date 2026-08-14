# Document Numbering Settings (Movements & Contracts)

## Purpose
HR admins currently type the SK/decision-letter number (`decision_letter_number`) and contract number (`contract_number`) by hand when creating an employee movement (promotion/mutation/demotion) or a contract. This is error-prone and inconsistent. This feature lets tenants configure a numbering format per document type in Settings, and the number is generated automatically when a movement or contract is created — while still allowing manual override.

## Scope
- Two document types: `employee_movement` (covers promotion/mutation/demotion SKs) and `employee_contract`.
- One shared format/counter per document type (not per movement sub-type).
- Tenant-scoped (each tenant configures its own format).

## Data Model
New tenant-DB table `document_numbering_settings`, migration `109_document_numbering_settings` (mysql + postgres, with `.down.sql`):

| Column | Type | Notes |
|---|---|---|
| `id` | PK | |
| `document_type` | varchar/enum | `employee_movement` \| `employee_contract`, unique |
| `format_template` | varchar(255) | free-text with tokens, e.g. `SK/{sequence:3}/HRIS/{month_roman}/{year}` |
| `reset_period` | varchar/enum | `yearly` \| `monthly` \| `never` |
| `last_sequence` | int | current counter value |
| `last_reset_key` | varchar(16) | last period key seen (`"2026"` for yearly, `"2026-08"` for monthly, `""` for never) |
| `created_at`, `updated_at` | timestamp | |

Seed rows on migration: one row per document type with a sensible default template and `reset_period = yearly`.

### Supported tokens
- `{sequence}` — plain integer sequence
- `{sequence:N}` — zero-padded to N digits (e.g. `{sequence:3}` → `007`)
- `{year}` — 4-digit year
- `{yy}` — 2-digit year
- `{month}` — 2-digit month
- `{month_roman}` — month in Roman numerals (I–XII)

Unknown/malformed tokens are left as literal text (no error) — validated client-side in the settings form with a live preview instead.

## Generation Logic
New function in a small `numbering` package (or within `employeemovement` service — implementation plan decides placement) :

```
GenerateDocumentNumber(tx, documentType) (string, error)
```

- Runs inside the same DB transaction as the movement/contract `Create`.
- `SELECT ... FOR UPDATE` on the matching `document_numbering_settings` row (prevents race conditions on concurrent creates).
- Computes `resetKey` for `reset_period`:
  - `yearly` → current year as string
  - `monthly` → `YYYY-MM`
  - `never` → constant `""`
- If `resetKey != last_reset_key`: set `last_sequence = 0`, `last_reset_key = resetKey`.
- Increment `last_sequence`, persist the row, format `format_template` substituting tokens, return the resulting string.

### Integration point
In `backend/internal/modules/employeemovement/service.go`, in the `Create` methods for movement and contract:
- If the incoming DTO's number field (`decision_letter_number` for movement, `contract_number` for contract) is empty/blank → call `GenerateDocumentNumber` and use the result.
- If non-empty → use the user-provided value as-is (manual override, no validation against the sequence).
- This preserves current behavior for `decision_letter_number` on contracts, which stays optional.

## Backend API
New handlers under `backend/internal/modules/settings` (or extend existing settings routes if a package already exists by implementation time):

- `GET /api/v1/tenant/settings/document-numbering` → returns both config rows (movement, contract).
- `PUT /api/v1/tenant/settings/document-numbering/:document_type` → update `format_template` and `reset_period` for one type. Validates `document_type` is one of the two allowed values and `reset_period` is one of the three allowed values.
- `GET /api/v1/tenant/settings/document-numbering/:document_type/preview` → returns what the *next* number would look like, computed read-only (does not increment `last_sequence`), for live preview in the settings UI.

## Frontend
- New page: `frontend/tenant/src/views/settings/NumberingSettingsView.vue`, following the `CompanyHolidaysView.vue` pattern (axios calls directly via `@/services/api`, PrimeVue form controls, `FormRow`/`TextInput`, `useI18n()`, `useToast()`, `getValidationErrors`).
- Route added in `frontend/tenant/src/router/index.js` under the settings children: `settings/numbering` → `SettingsNumbering`, plus a sidebar entry in `frontend/tenant/src/layouts/Sidebar.vue`.
- Page layout: two sections ("Movements" / "Contracts"), each with:
  - Template text input + inline help listing supported tokens
  - Reset period dropdown (Yearly / Monthly / Never)
  - Read-only "next number preview" that re-fetches the preview endpoint (debounced) as the template changes
  - Save button per section
- `frontend/tenant/src/views/modules/employeemovement/EmployeeMovements.vue` and `EmployeeContracts.vue`:
  - `decision_letter_number` / `contract_number` `TextInput` fields remain editable, but default to empty on create, placeholder text indicates "Auto-generated if left blank", and the client-side "required" validation on `contract_number` is removed for the create flow (still required check keeps working for edit if you want to keep it non-empty, but since it'll already be populated after first save this is moot).
  - After a successful create response, the returned number is shown in the form/toast so the user sees what was generated.

## Testing
- Backend: unit tests for the numbering generator covering token formatting, zero-padding, and reset transitions (year rollover, month rollover, `never` never resets). Table-driven test with a fixed clock/injected "now".
- Backend: integration test that two concurrent `Create` calls do not produce duplicate numbers (exercises the `FOR UPDATE` lock).
- Manual verification in the browser: configure a template in Settings, create a movement/contract with the number field blank and confirm the generated value matches the template; edit the template and confirm the preview and next generated number reflect the change; manually type a number and confirm it's respected without being overwritten.

## Out of Scope
- Per-movement-subtype (promotion vs mutation vs demotion) separate numbering — explicitly deferred; a single shared counter per document type is sufficient per requirements.
- Numbering for other document types (offer letters, requisitions, etc.) — existing hardcoded generators are untouched.
