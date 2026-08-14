# Document Template — Phase 1 (DB & Backend Foundation) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the backend foundation for the Document Template feature (Settings module): schema, GORM models, service logic (one-active-template-per-document-type, template versioning, "use default" copy flow), and a full REST API — no PDF generation, no editor UI, no Contract/Movement integration yet (those are later phases).

**Architecture:** New flat Go module `backend/internal/modules/documenttemplate` (same shape as `setting`/`employeemovement`): model, repository, service, handler, routes, module, errors, dto. The existing `document_templates` table (created in migration `011_settings.sql`) is extended via `ALTER TABLE` in a new versioned migration `110_document_templates.sql`; three new tables (`document_template_versions`, `document_template_audits`, `generated_documents`) are created in the same migration file. Table creation/alteration happens **only** through the versioned SQL migration — never through GORM `AutoMigrate` — because tenant-DB schema changes in this codebase must go through `migrator/migrations/tenant/{mysql,postgres}/`; a module's `Migrate(db)` method is a required interface stub but is not the mechanism that actually changes tenant schema.

**Tech Stack:** Go, Gin, GORM (MySQL + Postgres tenant DBs).

**Spec:** `docs/module-settngs-fitur-template-dokumen-plan.md` (sections 1–14, 18–19 apply to this phase; sections 6–9, 15–17, 20 are later phases). This plan implements only the "Phase 1 — Database & Backend Foundation" checklist from that document's §21, expanded into concrete tasks.

## Global Constraints

- Tenant DB schema changes go **only** through versioned SQL migrations in `backend/internal/pkg/migrator/migrations/tenant/{mysql,postgres}/`. Next available migration number: **110**. Every mysql file has a matching postgres file, and both have `.down.sql` counterparts.
- Every response struct that crosses the JSON wire (anything returned from a handler) **must** carry explicit `json:"snake_case_name"` tags on every field — do not rely on Go's default PascalCase marshaling. (A prior feature in this codebase shipped with GORM-only tags and no JSON tags, which silently broke on the frontend; do not repeat that mistake.)
- Document type is a closed, validated set for this phase: `CONTRACT_AGREEMENT`, `MOVEMENT_SK`. Reject anything else with a 400, but keep the underlying column a plain string (not a DB enum) so future document types (per spec §2.2) don't require a migration.
- Only one template may have `status = 'ACTIVE'` per document type at a time. Activating a template must, in the same DB transaction, deactivate any other template of the same document type. Enforce with a Postgres partial unique index; on MySQL (no partial index support) enforce the same rule in the service layer inside a transaction.
- Exactly one template per document type has `is_default = true` (the seeded reference template). It can never be activated, never deleted, and is not counted toward the "one active" rule.
- `company_id` tenant context key is the raw string literal `"company_id"` (no typed constant exists in this codebase — match the existing convention, don't introduce a new one).
- Permission strings follow `documenttemplate.<resource>.<action>` (lowercase, underscore), registered in the module's `Info().Permissions`.

---

### Task 1: Migration 110 — extend `document_templates`, add version/audit/generated tables, seed defaults

**Files:**
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.down.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.sql`
- Create: `backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.down.sql`

**Interfaces:**
- Produces: `document_templates` gains columns `code, description, active_version_id, status, is_default, deleted_at`; existing `type` column continues to hold the document type value (mapped in Go as `DocumentType`, gorm column `type` — no column rename). New tables `document_template_versions`, `document_template_audits`, `generated_documents`. Seed: one `REFERENCE`/`is_default=true` row per document type (`CONTRACT_AGREEMENT`, `MOVEMENT_SK`) in `document_templates`, each with one corresponding `v1` row in `document_template_versions` referenced by `active_version_id`... except default/reference templates do **not** use `active_version_id` (per spec §2.3/§10 — reference templates don't use the active-version mechanism); store the reference content directly in `document_templates.content` instead, leave `active_version_id` NULL for default rows.

- [ ] **Step 1: Write the mysql up migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.sql`:
```sql
ALTER TABLE document_templates
    ADD COLUMN code VARCHAR(100) NULL AFTER name,
    ADD COLUMN description TEXT NULL AFTER code,
    ADD COLUMN active_version_id CHAR(36) NULL AFTER content,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'INACTIVE' AFTER active_version_id,
    ADD COLUMN is_default TINYINT(1) NOT NULL DEFAULT 0 AFTER status,
    ADD COLUMN deleted_at TIMESTAMP NULL AFTER updated_at;

CREATE UNIQUE INDEX uk_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type_status ON document_templates (type, status);

CREATE TABLE IF NOT EXISTS document_template_versions (
    id            CHAR(36) PRIMARY KEY,
    template_id   CHAR(36) NOT NULL,
    version       INT NOT NULL,
    content       LONGTEXT NOT NULL,
    paper_size    VARCHAR(20) NOT NULL DEFAULT 'A4',
    orientation   VARCHAR(20) NOT NULL DEFAULT 'portrait',
    margin_top    INT NOT NULL DEFAULT 20,
    margin_right  INT NOT NULL DEFAULT 20,
    margin_bottom INT NOT NULL DEFAULT 20,
    margin_left   INT NOT NULL DEFAULT 20,
    created_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_doctplver_template FOREIGN KEY (template_id) REFERENCES document_templates(id) ON DELETE CASCADE,
    CONSTRAINT uq_doctplver_template_version UNIQUE (template_id, version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS document_template_audits (
    id          CHAR(36) PRIMARY KEY,
    template_id CHAR(36) NOT NULL,
    version_id  CHAR(36) NULL,
    action      VARCHAR(50) NOT NULL,
    actor_id    CHAR(36) NULL,
    payload     JSON NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_doctplaudit_template FOREIGN KEY (template_id) REFERENCES document_templates(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_doctplaudit_template ON document_template_audits (template_id);

CREATE TABLE IF NOT EXISTS generated_documents (
    id                   CHAR(36) PRIMARY KEY,
    template_id          CHAR(36) NOT NULL,
    template_version_id  CHAR(36) NOT NULL,
    document_type        VARCHAR(50) NOT NULL,
    reference_type       VARCHAR(50) NOT NULL,
    reference_id         CHAR(36) NOT NULL,
    file_name            VARCHAR(255) NOT NULL,
    file_path            VARCHAR(500) NOT NULL,
    mime_type            VARCHAR(100) NOT NULL DEFAULT 'application/pdf',
    generated_by         CHAR(36) NULL,
    generated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_gendoc_template FOREIGN KEY (template_id) REFERENCES document_templates(id),
    CONSTRAINT fk_gendoc_version FOREIGN KEY (template_version_id) REFERENCES document_template_versions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_gendoc_reference ON generated_documents (reference_type, reference_id);

-- Seed: one REFERENCE (default) template per document type. Reference
-- templates store content directly and never use active_version_id.
INSERT INTO document_templates (id, name, code, type, description, content, active_version_id, status, is_default, is_active, created_at, updated_at)
VALUES
    (UUID(), 'Perjanjian Kerja Waktu Tertentu (Default)', 'CONTRACT_AGREEMENT_DEFAULT', 'CONTRACT_AGREEMENT', 'Template referensi bawaan untuk Perjanjian Kontrak', '<h2>PERJANJIAN KERJA WAKTU TERTENTU</h2><p>Nomor: {{contract.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan: {{employee.position}}</p>', NULL, 'REFERENCE', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (UUID(), 'SK Movement (Default)', 'MOVEMENT_SK_DEFAULT', 'MOVEMENT_SK', 'Template referensi bawaan untuk SK Movement', '<h2>SURAT KEPUTUSAN MUTASI/PROMOSI/DEMOSI</h2><p>Nomor: {{movement.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan Baru: {{movement.new_position}}</p>', NULL, 'REFERENCE', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON DUPLICATE KEY UPDATE id = id;
```

- [ ] **Step 2: Write the mysql down migration**

`backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.down.sql`:
```sql
DROP TABLE IF EXISTS generated_documents;
DROP TABLE IF EXISTS document_template_audits;
DROP TABLE IF EXISTS document_template_versions;

DELETE FROM document_templates WHERE code IN ('CONTRACT_AGREEMENT_DEFAULT', 'MOVEMENT_SK_DEFAULT');

DROP INDEX idx_document_templates_type_status ON document_templates;
DROP INDEX uk_document_templates_code ON document_templates;

ALTER TABLE document_templates
    DROP COLUMN deleted_at,
    DROP COLUMN is_default,
    DROP COLUMN status,
    DROP COLUMN active_version_id,
    DROP COLUMN description,
    DROP COLUMN code;
```

- [ ] **Step 3: Write the postgres up migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.sql`:
```sql
ALTER TABLE document_templates
    ADD COLUMN code VARCHAR(100) NULL,
    ADD COLUMN description TEXT NULL,
    ADD COLUMN active_version_id CHAR(36) NULL,
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'INACTIVE',
    ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN deleted_at TIMESTAMP NULL;

CREATE UNIQUE INDEX uk_document_templates_code ON document_templates (code);
CREATE INDEX idx_document_templates_type_status ON document_templates (type, status);

-- Only one ACTIVE template per document type (partial unique index, Postgres only)
CREATE UNIQUE INDEX uq_document_templates_active_per_type
    ON document_templates (type)
    WHERE status = 'ACTIVE';

-- Only one default (reference) template per document type
CREATE UNIQUE INDEX uq_document_templates_default_per_type
    ON document_templates (type)
    WHERE is_default = TRUE;

CREATE TABLE IF NOT EXISTS document_template_versions (
    id            CHAR(36) PRIMARY KEY,
    template_id   CHAR(36) NOT NULL REFERENCES document_templates(id) ON DELETE CASCADE,
    version       INT NOT NULL,
    content       TEXT NOT NULL,
    paper_size    VARCHAR(20) NOT NULL DEFAULT 'A4',
    orientation   VARCHAR(20) NOT NULL DEFAULT 'portrait',
    margin_top    INT NOT NULL DEFAULT 20,
    margin_right  INT NOT NULL DEFAULT 20,
    margin_bottom INT NOT NULL DEFAULT 20,
    margin_left   INT NOT NULL DEFAULT 20,
    created_by    CHAR(36) NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_doctplver_template_version UNIQUE (template_id, version)
);

CREATE TABLE IF NOT EXISTS document_template_audits (
    id          CHAR(36) PRIMARY KEY,
    template_id CHAR(36) NOT NULL REFERENCES document_templates(id) ON DELETE CASCADE,
    version_id  CHAR(36) NULL,
    action      VARCHAR(50) NOT NULL,
    actor_id    CHAR(36) NULL,
    payload     JSONB NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_doctplaudit_template ON document_template_audits (template_id);

CREATE TABLE IF NOT EXISTS generated_documents (
    id                   CHAR(36) PRIMARY KEY,
    template_id          CHAR(36) NOT NULL REFERENCES document_templates(id),
    template_version_id  CHAR(36) NOT NULL REFERENCES document_template_versions(id),
    document_type        VARCHAR(50) NOT NULL,
    reference_type       VARCHAR(50) NOT NULL,
    reference_id         CHAR(36) NOT NULL,
    file_name            VARCHAR(255) NOT NULL,
    file_path            VARCHAR(500) NOT NULL,
    mime_type            VARCHAR(100) NOT NULL DEFAULT 'application/pdf',
    generated_by         CHAR(36) NULL,
    generated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_gendoc_reference ON generated_documents (reference_type, reference_id);

INSERT INTO document_templates (id, name, code, type, description, content, active_version_id, status, is_default, is_active, created_at, updated_at)
VALUES
    (gen_random_uuid(), 'Perjanjian Kerja Waktu Tertentu (Default)', 'CONTRACT_AGREEMENT_DEFAULT', 'CONTRACT_AGREEMENT', 'Template referensi bawaan untuk Perjanjian Kontrak', '<h2>PERJANJIAN KERJA WAKTU TERTENTU</h2><p>Nomor: {{contract.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan: {{employee.position}}</p>', NULL, 'REFERENCE', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'SK Movement (Default)', 'MOVEMENT_SK_DEFAULT', 'MOVEMENT_SK', 'Template referensi bawaan untuk SK Movement', '<h2>SURAT KEPUTUSAN MUTASI/PROMOSI/DEMOSI</h2><p>Nomor: {{movement.number}}</p><p>Nama: {{employee.name}}</p><p>Jabatan Baru: {{movement.new_position}}</p>', NULL, 'REFERENCE', TRUE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (code) DO NOTHING;
```

- [ ] **Step 4: Write the postgres down migration**

`backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.down.sql`:
```sql
DROP TABLE IF EXISTS generated_documents;
DROP TABLE IF EXISTS document_template_audits;
DROP TABLE IF EXISTS document_template_versions;

DELETE FROM document_templates WHERE code IN ('CONTRACT_AGREEMENT_DEFAULT', 'MOVEMENT_SK_DEFAULT');

DROP INDEX IF EXISTS uq_document_templates_default_per_type;
DROP INDEX IF EXISTS uq_document_templates_active_per_type;
DROP INDEX IF EXISTS idx_document_templates_type_status;
DROP INDEX IF EXISTS uk_document_templates_code;

ALTER TABLE document_templates
    DROP COLUMN deleted_at,
    DROP COLUMN is_default,
    DROP COLUMN status,
    DROP COLUMN active_version_id,
    DROP COLUMN description,
    DROP COLUMN code;
```

- [ ] **Step 5: Sanity-check the SQL**

Re-read all four files. Confirm: mysql `code` uses `ON DUPLICATE KEY UPDATE id = id` (idempotent replay), postgres uses `ON CONFLICT (code) DO NOTHING`; both `.down.sql` files reverse every `ADD COLUMN`/`CREATE TABLE`/`CREATE INDEX` the up-migration made, in reverse dependency order (drop tables with FKs to `document_templates` before altering it). If a local/test tenant DB is reachable, apply the migration and confirm `document_templates` has 2 rows with `is_default = true`; if no DB is reachable in this environment, note that in your report and rely on the visual re-read instead.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.sql backend/internal/pkg/migrator/migrations/tenant/mysql/110_document_templates.down.sql backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.sql backend/internal/pkg/migrator/migrations/tenant/postgres/110_document_templates.down.sql
git commit -m "feat: add document_templates extension migration + version/audit/generated tables"
```

---

### Task 2: `documenttemplate` package — models + domain errors

**Files:**
- Create: `backend/internal/modules/documenttemplate/model.go`
- Create: `backend/internal/modules/documenttemplate/errors.go`
- Test: `backend/internal/modules/documenttemplate/model_test.go`

**Interfaces:**
- Produces:
  - `type DocumentTemplate struct` (table `document_templates`) with fields: `ID, Name, Code, DocumentType (gorm column "type"), Description, Content, ActiveVersionID, Status, IsDefault, IsActive, CreatedAt, UpdatedAt, DeletedAt` — every field carries both `gorm:` and `json:"snake_case"` tags (`DeletedAt` uses `json:"-"`).
  - `type DocumentTemplateVersion struct` (table `document_template_versions`): `ID, TemplateID, Version, Content, PaperSize, Orientation, MarginTop, MarginRight, MarginBottom, MarginLeft, CreatedBy, CreatedAt` — same json-tag rule.
  - `type DocumentTemplateAudit struct` (table `document_template_audits`): `ID, TemplateID, VersionID, Action, ActorID, Payload (json.RawMessage), CreatedAt`.
  - `type GeneratedDocument struct` (table `generated_documents`): full field set per Task 1's schema — this phase only needs the struct to exist (no handler uses it yet), but it must compile and carry json tags for forward compatibility.
  - Constants: `DocumentTypeContractAgreement = "CONTRACT_AGREEMENT"`, `DocumentTypeMovementSK = "MOVEMENT_SK"`; `StatusActive = "ACTIVE"`, `StatusInactive = "INACTIVE"`, `StatusReference = "REFERENCE"`.
  - `errors.go`: `type DuplicateActiveTemplateError struct { DocumentType string }`, `type DuplicateCodeError struct { Code string }`, `type ReferenceTemplateImmutableError struct { Action string }` (returned when someone tries to activate/delete/directly-generate-from a reference template), `ErrTemplateNotFound = errors.New(...)`, `ErrVersionNotFound = errors.New(...)` — each error type has an `Error() string` method, following the exact style of `setting/errors.go`'s `DuplicateCodeError`.

- [ ] **Step 1: Write the failing JSON-key regression test**

`backend/internal/modules/documenttemplate/model_test.go` — this test exists specifically because a prior feature in this codebase shipped a model with GORM tags but no JSON tags, silently breaking the frontend; guard against repeating that:
```go
package documenttemplate

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDocumentTemplateJSONKeysAreSnakeCase(t *testing.T) {
	tpl := DocumentTemplate{
		ID:           "id-1",
		Name:         "Test",
		Code:         "TEST_CODE",
		DocumentType: DocumentTypeContractAgreement,
		Status:       StatusActive,
	}
	b, err := json.Marshal(tpl)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{`"id"`, `"name"`, `"code"`, `"document_type"`, `"status"`, `"is_default"`, `"created_at"`} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON output to contain %s, got: %s", key, s)
		}
	}
	for _, badKey := range []string{`"ID"`, `"DocumentType"`, `"Status"`} {
		if strings.Contains(s, badKey) {
			t.Errorf("JSON output should not contain PascalCase key %s, got: %s", badKey, s)
		}
	}
}

func TestDocumentTemplateVersionJSONKeysAreSnakeCase(t *testing.T) {
	v := DocumentTemplateVersion{ID: "v-1", TemplateID: "t-1", Version: 1, Content: "<p>x</p>"}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(b)
	for _, key := range []string{`"id"`, `"template_id"`, `"version"`, `"content"`, `"paper_size"`, `"orientation"`} {
		if !strings.Contains(s, key) {
			t.Errorf("expected JSON output to contain %s, got: %s", key, s)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -run TestDocumentTemplateJSONKeysAreSnakeCase -v`
Expected: FAIL — build error, package/types don't exist yet.

- [ ] **Step 3: Write `model.go`**

```go
package documenttemplate

import (
	"encoding/json"
	"time"
)

const (
	DocumentTypeContractAgreement = "CONTRACT_AGREEMENT"
	DocumentTypeMovementSK        = "MOVEMENT_SK"
)

const (
	StatusActive    = "ACTIVE"
	StatusInactive  = "INACTIVE"
	StatusReference = "REFERENCE"
)

// ValidDocumentTypes is the closed set accepted by this phase; extend when
// spec §2.2's future document types are added.
var ValidDocumentTypes = map[string]bool{
	DocumentTypeContractAgreement: true,
	DocumentTypeMovementSK:        true,
}

type DocumentTemplate struct {
	ID              string     `gorm:"column:id;primaryKey" json:"id"`
	Name            string     `gorm:"column:name" json:"name"`
	Code            string     `gorm:"column:code" json:"code"`
	DocumentType    string     `gorm:"column:type" json:"document_type"`
	Description     *string    `gorm:"column:description" json:"description,omitempty"`
	Content         *string    `gorm:"column:content" json:"content,omitempty"`
	ActiveVersionID *string    `gorm:"column:active_version_id" json:"active_version_id,omitempty"`
	Status          string     `gorm:"column:status" json:"status"`
	IsDefault       bool       `gorm:"column:is_default" json:"is_default"`
	IsActive        bool       `gorm:"column:is_active" json:"is_active"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at" json:"-"`
}

func (DocumentTemplate) TableName() string { return "document_templates" }

type DocumentTemplateVersion struct {
	ID           string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateID   string    `gorm:"column:template_id" json:"template_id"`
	Version      int       `gorm:"column:version" json:"version"`
	Content      string    `gorm:"column:content" json:"content"`
	PaperSize    string    `gorm:"column:paper_size" json:"paper_size"`
	Orientation  string    `gorm:"column:orientation" json:"orientation"`
	MarginTop    int       `gorm:"column:margin_top" json:"margin_top"`
	MarginRight  int       `gorm:"column:margin_right" json:"margin_right"`
	MarginBottom int       `gorm:"column:margin_bottom" json:"margin_bottom"`
	MarginLeft   int       `gorm:"column:margin_left" json:"margin_left"`
	CreatedBy    *string   `gorm:"column:created_by" json:"created_by,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (DocumentTemplateVersion) TableName() string { return "document_template_versions" }

type DocumentTemplateAudit struct {
	ID         string          `gorm:"column:id;primaryKey" json:"id"`
	TemplateID string          `gorm:"column:template_id" json:"template_id"`
	VersionID  *string         `gorm:"column:version_id" json:"version_id,omitempty"`
	Action     string          `gorm:"column:action" json:"action"`
	ActorID    *string         `gorm:"column:actor_id" json:"actor_id,omitempty"`
	Payload    json.RawMessage `gorm:"column:payload" json:"payload,omitempty"`
	CreatedAt  time.Time       `gorm:"column:created_at" json:"created_at"`
}

func (DocumentTemplateAudit) TableName() string { return "document_template_audits" }

type GeneratedDocument struct {
	ID                 string    `gorm:"column:id;primaryKey" json:"id"`
	TemplateID         string    `gorm:"column:template_id" json:"template_id"`
	TemplateVersionID  string    `gorm:"column:template_version_id" json:"template_version_id"`
	DocumentType       string    `gorm:"column:document_type" json:"document_type"`
	ReferenceType      string    `gorm:"column:reference_type" json:"reference_type"`
	ReferenceID        string    `gorm:"column:reference_id" json:"reference_id"`
	FileName           string    `gorm:"column:file_name" json:"file_name"`
	FilePath           string    `gorm:"column:file_path" json:"file_path"`
	MimeType           string    `gorm:"column:mime_type" json:"mime_type"`
	GeneratedBy        *string   `gorm:"column:generated_by" json:"generated_by,omitempty"`
	GeneratedAt        time.Time `gorm:"column:generated_at" json:"generated_at"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"created_at"`
}

func (GeneratedDocument) TableName() string { return "generated_documents" }
```

- [ ] **Step 4: Write `errors.go`**

```go
package documenttemplate

import (
	"errors"
	"fmt"
)

var (
	ErrTemplateNotFound = errors.New("document template not found")
	ErrVersionNotFound  = errors.New("document template version not found")
)

type DuplicateActiveTemplateError struct {
	DocumentType string
}

func (e *DuplicateActiveTemplateError) Error() string {
	return fmt.Sprintf("an active template already exists for document type '%s'", e.DocumentType)
}

type DuplicateCodeError struct {
	Code string
}

func (e *DuplicateCodeError) Error() string {
	return fmt.Sprintf("template code '%s' already exists", e.Code)
}

type ReferenceTemplateImmutableError struct {
	Action string
}

func (e *ReferenceTemplateImmutableError) Error() string {
	return fmt.Sprintf("reference (default) templates cannot be %s", e.Action)
}

type InvalidDocumentTypeError struct {
	DocumentType string
}

func (e *InvalidDocumentTypeError) Error() string {
	return fmt.Sprintf("invalid document type '%s'", e.DocumentType)
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -v`
Expected: PASS — both JSON-key tests green.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/documenttemplate/model.go backend/internal/modules/documenttemplate/errors.go backend/internal/modules/documenttemplate/model_test.go
git commit -m "feat: add documenttemplate models and domain errors"
```

---

### Task 3: Repository — CRUD, pagination, versions, audit log (TDD with in-memory sqlite)

**Files:**
- Create: `backend/internal/modules/documenttemplate/repository.go`
- Create: `backend/internal/modules/documenttemplate/helpers_test.go`
- Test: `backend/internal/modules/documenttemplate/repository_test.go`

**Interfaces:**
- Consumes: `DocumentTemplate`, `DocumentTemplateVersion`, `DocumentTemplateAudit` from Task 2.
- Produces:
  - `type Repository struct { dbResolver func(ctx context.Context) (*gorm.DB, error) }`
  - `func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository`
  - `func NewTenantDBResolver(dbManager *database.Manager) func(ctx context.Context) (*gorm.DB, error)` — reads `ctx.Value("company_id").(string)`, calls `dbManager.TenantDB(companyID)`, matching `setting/module.go`'s existing pattern exactly (including its exact error message on missing `company_id`).
  - `func (r *Repository) List(ctx context.Context, page, perPage int, documentType, status, search string) ([]DocumentTemplate, int64, error)` — paginated, `search` matches `name` or `code` (case-insensitive LIKE, wildcards escaped), `documentType`/`status` are optional exact-match filters, excludes soft-deleted rows (`deleted_at IS NULL`).
  - `func (r *Repository) GetByID(ctx context.Context, id string) (*DocumentTemplate, error)` — returns `ErrTemplateNotFound` when absent or soft-deleted.
  - `func (r *Repository) GetByCode(ctx context.Context, code string) (*DocumentTemplate, error)`
  - `func (r *Repository) FindActiveByType(ctx context.Context, documentType string) (*DocumentTemplate, error)` — returns `ErrTemplateNotFound` if none.
  - `func (r *Repository) FindDefaultByType(ctx context.Context, documentType string) (*DocumentTemplate, error)`
  - `func (r *Repository) Create(ctx context.Context, tpl *DocumentTemplate) error`
  - `func (r *Repository) Update(ctx context.Context, tpl *DocumentTemplate) error`
  - `func (r *Repository) SoftDelete(ctx context.Context, id string) error` — sets `deleted_at`.
  - `func (r *Repository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error` — thin wrapper over `db.Transaction(fn)`, used by the service for the activate/deactivate transaction in Task 4.
  - `func (r *Repository) CreateVersion(ctx context.Context, tx *gorm.DB, v *DocumentTemplateVersion) error` and `func (r *Repository) ListVersions(ctx context.Context, templateID string) ([]DocumentTemplateVersion, error)` and `func (r *Repository) GetVersion(ctx context.Context, templateID, versionID string) (*DocumentTemplateVersion, error)` and `func (r *Repository) NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID string) (int, error)` (returns `MAX(version)+1`, or `1` if none exist — call inside the same tx as the insert to avoid a race).
  - `func (r *Repository) CreateAudit(ctx context.Context, tx *gorm.DB, a *DocumentTemplateAudit) error` — `tx` may be `nil`, in which case it resolves and uses a plain (non-transactional) `*gorm.DB` via `r.getDB(ctx)`.

- [ ] **Step 1: Write `helpers_test.go`**

Copy the exact structure of `backend/internal/modules/setting/helpers_test.go` (sqlite in-memory setup, `AutoMigrate`, cleanup, `testDBResolver`) but migrate `documenttemplate`'s four models instead:
```go
package documenttemplate

import (
	"context"
	"fmt"

	sqlite "github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open test db: %v", err))
	}
	if err := db.AutoMigrate(&DocumentTemplate{}, &DocumentTemplateVersion{}, &DocumentTemplateAudit{}, &GeneratedDocument{}); err != nil {
		panic(fmt.Sprintf("failed to automigrate: %v", err))
	}
	sqlDB, _ := db.DB()
	return db, func() { sqlDB.Close() }
}

func testDBResolver(db *gorm.DB) func(ctx context.Context) (*gorm.DB, error) {
	return func(ctx context.Context) (*gorm.DB, error) { return db, nil }
}

func newTestRepo(db *gorm.DB) *Repository {
	return NewRepository(testDBResolver(db))
}

func uuidStr() string { return uuid.New().String() }

func createTestTemplate(db *gorm.DB, code, documentType, status string, isDefault bool) *DocumentTemplate {
	tpl := &DocumentTemplate{
		ID:           uuidStr(),
		Name:         code,
		Code:         code,
		DocumentType: documentType,
		Status:       status,
		IsDefault:    isDefault,
		IsActive:     true,
	}
	if err := db.Create(tpl).Error; err != nil {
		panic(fmt.Sprintf("failed to create test template: %v", err))
	}
	return tpl
}
```

- [ ] **Step 2: Write the failing repository tests**

`backend/internal/modules/documenttemplate/repository_test.go`:
```go
package documenttemplate

import (
	"context"
	"testing"
)

func TestRepositoryCreateAndGetByID(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()

	tpl := &DocumentTemplate{ID: uuidStr(), Name: "PKWT", Code: "PKWT01", DocumentType: DocumentTypeContractAgreement, Status: StatusInactive}
	if err := repo.Create(ctx, tpl); err != nil {
		t.Fatalf("Create: %v", err)
	}
	got, err := repo.GetByID(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Code != "PKWT01" {
		t.Fatalf("expected code PKWT01, got %s", got.Code)
	}
}

func TestRepositoryGetByIDNotFound(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)

	_, err := repo.GetByID(context.Background(), uuidStr())
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound, got %v", err)
	}
}

func TestRepositoryFindActiveByType(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	createTestTemplate(db, "INACTIVE1", DocumentTypeContractAgreement, StatusInactive, false)
	active := createTestTemplate(db, "ACTIVE1", DocumentTypeContractAgreement, StatusActive, false)

	got, err := repo.FindActiveByType(context.Background(), DocumentTypeContractAgreement)
	if err != nil {
		t.Fatalf("FindActiveByType: %v", err)
	}
	if got.ID != active.ID {
		t.Fatalf("expected active template %s, got %s", active.ID, got.ID)
	}
}

func TestRepositoryListPaginationAndSearch(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	for i := 0; i < 5; i++ {
		createTestTemplate(db, uuidStr()[:8], DocumentTypeContractAgreement, StatusInactive, false)
	}
	createTestTemplate(db, "FINDME", DocumentTypeMovementSK, StatusInactive, false)

	items, total, err := repo.List(context.Background(), 1, 3, "", "", "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 6 || len(items) != 3 {
		t.Fatalf("expected total=6 len=3, got total=%d len=%d", total, len(items))
	}

	filtered, ftotal, err := repo.List(context.Background(), 1, 10, DocumentTypeMovementSK, "", "")
	if err != nil {
		t.Fatalf("List filtered: %v", err)
	}
	if ftotal != 1 || filtered[0].Code != "FINDME" {
		t.Fatalf("expected 1 result FINDME, got total=%d", ftotal)
	}

	searched, stotal, err := repo.List(context.Background(), 1, 10, "", "", "findme")
	if err != nil {
		t.Fatalf("List search: %v", err)
	}
	if stotal != 1 || searched[0].Code != "FINDME" {
		t.Fatalf("expected search to find FINDME case-insensitively, got total=%d", stotal)
	}
}

func TestRepositorySoftDeleteExcludesFromList(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()
	tpl := createTestTemplate(db, "TOBEDEL", DocumentTypeContractAgreement, StatusInactive, false)

	if err := repo.SoftDelete(ctx, tpl.ID); err != nil {
		t.Fatalf("SoftDelete: %v", err)
	}
	_, err := repo.GetByID(ctx, tpl.ID)
	if err != ErrTemplateNotFound {
		t.Fatalf("expected ErrTemplateNotFound after soft delete, got %v", err)
	}
	items, total, _ := repo.List(ctx, 1, 10, "", "", "")
	if total != 0 || len(items) != 0 {
		t.Fatalf("expected soft-deleted row excluded from List, got total=%d", total)
	}
}

func TestRepositoryVersionsCreateListNextNumber(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	ctx := context.Background()
	tpl := createTestTemplate(db, "VERTEST", DocumentTypeContractAgreement, StatusInactive, false)

	err := repo.WithTx(ctx, func(tx interface{ Error }) error { return nil }) // placeholder removed below
	_ = err

	err = repo.WithTx(ctx, func(tx *gormDBAlias) error {
		next, nerr := repo.NextVersionNumber(ctx, tx, tpl.ID)
		if nerr != nil {
			return nerr
		}
		if next != 1 {
			t.Fatalf("expected first version number 1, got %d", next)
		}
		v := &DocumentTemplateVersion{ID: uuidStr(), TemplateID: tpl.ID, Version: next, Content: "<p>v1</p>", PaperSize: "A4", Orientation: "portrait"}
		return repo.CreateVersion(ctx, tx, v)
	})
	if err != nil {
		t.Fatalf("WithTx create version 1: %v", err)
	}

	err = repo.WithTx(ctx, func(tx *gormDBAlias) error {
		next, nerr := repo.NextVersionNumber(ctx, tx, tpl.ID)
		if nerr != nil {
			return nerr
		}
		if next != 2 {
			t.Fatalf("expected second version number 2, got %d", next)
		}
		v := &DocumentTemplateVersion{ID: uuidStr(), TemplateID: tpl.ID, Version: next, Content: "<p>v2</p>", PaperSize: "A4", Orientation: "portrait"}
		return repo.CreateVersion(ctx, tx, v)
	})
	if err != nil {
		t.Fatalf("WithTx create version 2: %v", err)
	}

	versions, err := repo.ListVersions(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("ListVersions: %v", err)
	}
	if len(versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(versions))
	}
}
```
**Note for the implementer:** the `gormDBAlias` placeholder and the first throwaway `repo.WithTx(ctx, func(tx interface{ Error }) error { return nil })` line above are WRONG Go — they won't compile. Replace `gormDBAlias` with the real `*gorm.DB` type (add `"gorm.io/gorm"` to the imports) and delete the throwaway placeholder line entirely; it was left in this brief only to flag that `WithTx`'s callback signature is `func(tx *gorm.DB) error`, matching the interface described above. Write the real, compiling version of this test using `*gorm.DB` directly.

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -v`
Expected: FAIL — `Repository`/`NewRepository`/etc. undefined (and/or compile errors from the test file, which is expected before `repository.go` exists).

- [ ] **Step 4: Write `repository.go`**

```go
package documenttemplate

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
)

type Repository struct {
	dbResolver func(ctx context.Context) (*gorm.DB, error)
}

func NewRepository(dbResolver func(ctx context.Context) (*gorm.DB, error)) *Repository {
	return &Repository{dbResolver: dbResolver}
}

func NewTenantDBResolver(dbManager *database.Manager) func(ctx context.Context) (*gorm.DB, error) {
	return func(ctx context.Context) (*gorm.DB, error) {
		companyID, ok := ctx.Value("company_id").(string)
		if !ok || companyID == "" {
			return nil, fmt.Errorf("tenant context not found in request: company_id is required")
		}
		return dbManager.TenantDB(companyID)
	}
}

func (r *Repository) getDB(ctx context.Context) (*gorm.DB, error) {
	return r.dbResolver(ctx)
}

func escapeLike(s string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return replacer.Replace(s)
}

func (r *Repository) List(ctx context.Context, page, perPage int, documentType, status, search string) ([]DocumentTemplate, int64, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, 0, err
	}
	query := db.Model(&DocumentTemplate{}).Where("deleted_at IS NULL")
	if documentType != "" {
		query = query.Where("type = ?", documentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		like := "%" + escapeLike(search) + "%"
		query = query.Where("(name LIKE ? OR code LIKE ?)", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count document templates: %w", err)
	}

	var items []DocumentTemplate
	offset := (page - 1) * perPage
	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list document templates: %w", err)
	}
	return items, total, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("id = ? AND deleted_at IS NULL", id).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get document template: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("code = ? AND deleted_at IS NULL", code).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to get document template by code: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) FindActiveByType(ctx context.Context, documentType string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("type = ? AND status = ? AND deleted_at IS NULL", documentType, StatusActive).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to find active document template: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) FindDefaultByType(ctx context.Context, documentType string) (*DocumentTemplate, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var tpl DocumentTemplate
	err = db.Where("type = ? AND is_default = ? AND deleted_at IS NULL", documentType, true).First(&tpl).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("failed to find default document template: %w", err)
	}
	return &tpl, nil
}

func (r *Repository) Create(ctx context.Context, tpl *DocumentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(tpl).Error; err != nil {
		return fmt.Errorf("failed to create document template: %w", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, tpl *DocumentTemplate) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Save(tpl).Error; err != nil {
		return fmt.Errorf("failed to update document template: %w", err)
	}
	return nil
}

func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Model(&DocumentTemplate{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP")).Error; err != nil {
		return fmt.Errorf("failed to soft delete document template: %w", err)
	}
	return nil
}

func (r *Repository) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	return db.Transaction(fn)
}

func (r *Repository) CreateVersion(ctx context.Context, tx *gorm.DB, v *DocumentTemplateVersion) error {
	if err := tx.Create(v).Error; err != nil {
		return fmt.Errorf("failed to create document template version: %w", err)
	}
	return nil
}

func (r *Repository) ListVersions(ctx context.Context, templateID string) ([]DocumentTemplateVersion, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var versions []DocumentTemplateVersion
	if err := db.Where("template_id = ?", templateID).Order("version DESC").Find(&versions).Error; err != nil {
		return nil, fmt.Errorf("failed to list document template versions: %w", err)
	}
	return versions, nil
}

func (r *Repository) GetVersion(ctx context.Context, templateID, versionID string) (*DocumentTemplateVersion, error) {
	db, err := r.getDB(ctx)
	if err != nil {
		return nil, err
	}
	var v DocumentTemplateVersion
	err = db.Where("id = ? AND template_id = ?", versionID, templateID).First(&v).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrVersionNotFound
		}
		return nil, fmt.Errorf("failed to get document template version: %w", err)
	}
	return &v, nil
}

func (r *Repository) NextVersionNumber(ctx context.Context, tx *gorm.DB, templateID string) (int, error) {
	var max int
	err := tx.Model(&DocumentTemplateVersion{}).
		Where("template_id = ?", templateID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&max).Error
	if err != nil {
		return 0, fmt.Errorf("failed to compute next version number: %w", err)
	}
	return max + 1, nil
}

func (r *Repository) CreateAudit(ctx context.Context, tx *gorm.DB, a *DocumentTemplateAudit) error {
	if tx != nil {
		if err := tx.Create(a).Error; err != nil {
			return fmt.Errorf("failed to create document template audit: %w", err)
		}
		return nil
	}
	db, err := r.getDB(ctx)
	if err != nil {
		return err
	}
	if err := db.Create(a).Error; err != nil {
		return fmt.Errorf("failed to create document template audit: %w", err)
	}
	return nil
}
```
Adjust the internal import path prefix (`github.com/inthros/hris-platform/...`) if the actual module path differs — confirm against any neighboring file's import block before finalizing.

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -v`
Expected: PASS — all repository tests green, including the version/next-number test you fixed in Step 2.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/modules/documenttemplate/repository.go backend/internal/modules/documenttemplate/helpers_test.go backend/internal/modules/documenttemplate/repository_test.go
git commit -m "feat: add documenttemplate repository with pagination, versions, audit log"
```

---

### Task 4: Service — one-active-per-type, from-default copy flow, versioning (TDD)

**Files:**
- Create: `backend/internal/modules/documenttemplate/service.go`
- Test: `backend/internal/modules/documenttemplate/service_test.go`

**Interfaces:**
- Consumes: `Repository` and all its methods from Task 3; models/errors from Task 2.
- Produces:
  - `type Service struct { repo *Repository; logger *zap.Logger }`
  - `func NewService(repo *Repository, logger *zap.Logger) *Service`
  - `func (s *Service) List(ctx, page, perPage int, documentType, status, search string) ([]DocumentTemplate, int64, error)`
  - `func (s *Service) GetByID(ctx, id string) (*DocumentTemplate, error)`
  - `func (s *Service) Create(ctx, name, code, documentType, description string, actorID string) (*DocumentTemplate, error)` — validates `documentType` against `ValidDocumentTypes` (else `*InvalidDocumentTypeError`), checks code uniqueness (else `*DuplicateCodeError`), creates with `status = INACTIVE`, `is_default = false`, writes a `CREATED` audit row.
  - `func (s *Service) CreateFromDefault(ctx, documentType, name, code, actorID string) (*DocumentTemplate, error)` — loads the default/reference template for `documentType` via `FindDefaultByType`, copies its `Content` into a new non-default `INACTIVE` template (new `name`/`code` supplied by caller), writes a `CREATED_FROM_DEFAULT` audit row. Returns `ErrTemplateNotFound` if no default exists for that type.
  - `func (s *Service) Update(ctx, id, name, description string, actorID string) (*DocumentTemplate, error)` — returns `*ReferenceTemplateImmutableError{Action: "edited"}` if the target `IsDefault`. Writes an `UPDATED` audit row.
  - `func (s *Service) UpdateDefaultContent(ctx, documentType, content, actorID string) (*DocumentTemplate, error)` — the one path allowed to change a reference template's content; writes a `DEFAULT_UPDATED` audit row.
  - `func (s *Service) Delete(ctx, id, actorID string) error` — returns `*ReferenceTemplateImmutableError{Action: "deleted"}` if `IsDefault`; soft-deletes; writes a `DELETED` audit row. (Spec §2.1 says only templates "not yet used" can be deleted — this phase has no `generated_documents` writer yet, so every template is deletable for now; leave a one-line comment noting the future check belongs here once Phase 5 exists.)
  - `func (s *Service) Activate(ctx, id, actorID string) (*DocumentTemplate, error)` — returns `*ReferenceTemplateImmutableError{Action: "activated"}` if `IsDefault`. Inside `repo.WithTx`: finds any other `ACTIVE` template of the same `DocumentType` (excluding itself) and sets it to `INACTIVE`, then sets the target to `ACTIVE`, writes an `ACTIVATED` audit row (and a `DEACTIVATED` row for the template that got bumped, if any).
  - `func (s *Service) Deactivate(ctx, id, actorID string) (*DocumentTemplate, error)` — sets `status = INACTIVE`, writes a `DEACTIVATED` audit row.
  - `func (s *Service) CreateVersion(ctx, templateID, content, paperSize, orientation string, margins [4]int, actorID string) (*DocumentTemplateVersion, error)` — inside `repo.WithTx`: computes `NextVersionNumber`, inserts the version, sets the template's `ActiveVersionID` to the new version's ID (`repo.Update` inside the same tx — note: `Repository.Update` currently takes no `tx` parameter; either add a `tx`-aware overload or have this method call `tx.Save(tpl)` directly on the model inside the transaction rather than going through `r.repo.Update`, since that method resolves its own non-transactional DB handle. Pick whichever is less invasive to Task 3's `Repository` and document the choice in your report), writes a `VERSION_CREATED` audit row.
  - `func (s *Service) ListVersions(ctx, templateID string) ([]DocumentTemplateVersion, error)`

- [ ] **Step 1: Write the failing service tests**

`backend/internal/modules/documenttemplate/service_test.go`:
```go
package documenttemplate

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

func newTestService(db interface{ AutoMigrate(...interface{}) error }) {} // placeholder — see note below

func TestServiceCreateRejectsInvalidDocumentType(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())

	_, err := svc.Create(context.Background(), "X", "CODEX", "NOT_A_TYPE", "", "actor-1")
	var invalidErr *InvalidDocumentTypeError
	if !errorsAs(err, &invalidErr) {
		t.Fatalf("expected InvalidDocumentTypeError, got %v", err)
	}
}

func TestServiceCreateRejectsDuplicateCode(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())
	ctx := context.Background()

	if _, err := svc.Create(ctx, "First", "DUPCODE", DocumentTypeContractAgreement, "", "actor-1"); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	_, err := svc.Create(ctx, "Second", "DUPCODE", DocumentTypeContractAgreement, "", "actor-1")
	var dupErr *DuplicateCodeError
	if !errorsAs(err, &dupErr) {
		t.Fatalf("expected DuplicateCodeError, got %v", err)
	}
}

func TestServiceActivateDeactivatesPreviousActive(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	svc := NewService(newTestRepo(db), zap.NewNop())
	ctx := context.Background()

	first, err := svc.Create(ctx, "First", "ACT1", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := svc.Activate(ctx, first.ID, "actor-1"); err != nil {
		t.Fatalf("activate first: %v", err)
	}

	second, err := svc.Create(ctx, "Second", "ACT2", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create second: %v", err)
	}
	if _, err := svc.Activate(ctx, second.ID, "actor-1"); err != nil {
		t.Fatalf("activate second: %v", err)
	}

	gotFirst, err := svc.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("get first: %v", err)
	}
	if gotFirst.Status != StatusInactive {
		t.Fatalf("expected first template to be deactivated, got status=%s", gotFirst.Status)
	}
	gotSecond, err := svc.GetByID(ctx, second.ID)
	if err != nil {
		t.Fatalf("get second: %v", err)
	}
	if gotSecond.Status != StatusActive {
		t.Fatalf("expected second template to be active, got status=%s", gotSecond.Status)
	}
}

func TestServiceActivateRejectsDefaultTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	def := createTestTemplate(db, "DEF1", DocumentTypeContractAgreement, StatusReference, true)

	_, err := svc.Activate(context.Background(), def.ID, "actor-1")
	var immErr *ReferenceTemplateImmutableError
	if !errorsAs(err, &immErr) {
		t.Fatalf("expected ReferenceTemplateImmutableError, got %v", err)
	}
}

func TestServiceCreateFromDefaultCopiesContent(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	content := "<p>default content</p>"
	def := createTestTemplate(db, "DEFCOPY", DocumentTypeContractAgreement, StatusReference, true)
	def.Content = &content
	if err := db.Save(def).Error; err != nil {
		t.Fatalf("seed default content: %v", err)
	}

	copied, err := svc.CreateFromDefault(context.Background(), DocumentTypeContractAgreement, "New From Default", "NEWCODE", "actor-1")
	if err != nil {
		t.Fatalf("CreateFromDefault: %v", err)
	}
	if copied.Content == nil || *copied.Content != content {
		t.Fatalf("expected copied content %q, got %v", content, copied.Content)
	}
	if copied.IsDefault {
		t.Fatalf("copied template must not itself be default")
	}
	if copied.Status != StatusInactive {
		t.Fatalf("copied template must start INACTIVE, got %s", copied.Status)
	}
}

func TestServiceCreateVersionIncrementsAndSetsActiveVersion(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	ctx := context.Background()
	tpl, err := svc.Create(ctx, "Vers", "VERSVC", DocumentTypeContractAgreement, "", "actor-1")
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	v1, err := svc.CreateVersion(ctx, tpl.ID, "<p>v1</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "actor-1")
	if err != nil {
		t.Fatalf("CreateVersion v1: %v", err)
	}
	if v1.Version != 1 {
		t.Fatalf("expected version 1, got %d", v1.Version)
	}
	updated, err := svc.GetByID(ctx, tpl.ID)
	if err != nil {
		t.Fatalf("GetByID after version: %v", err)
	}
	if updated.ActiveVersionID == nil || *updated.ActiveVersionID != v1.ID {
		t.Fatalf("expected template.active_version_id to point at v1, got %v", updated.ActiveVersionID)
	}

	v2, err := svc.CreateVersion(ctx, tpl.ID, "<p>v2</p>", "A4", "portrait", [4]int{20, 20, 20, 20}, "actor-1")
	if err != nil {
		t.Fatalf("CreateVersion v2: %v", err)
	}
	if v2.Version != 2 {
		t.Fatalf("expected version 2, got %d", v2.Version)
	}
}

func TestServiceDeleteRejectsDefaultTemplate(t *testing.T) {
	db, cleanup := setupTestDB()
	defer cleanup()
	repo := newTestRepo(db)
	svc := NewService(repo, zap.NewNop())
	def := createTestTemplate(db, "DEFDEL", DocumentTypeMovementSK, StatusReference, true)

	err := svc.Delete(context.Background(), def.ID, "actor-1")
	var immErr *ReferenceTemplateImmutableError
	if !errorsAs(err, &immErr) {
		t.Fatalf("expected ReferenceTemplateImmutableError, got %v", err)
	}
}
```
**Note for the implementer:** delete the `newTestService` placeholder line at the top (`func newTestService(db interface{...}) {}`) — it's invalid/unused Go left in this brief by mistake, not part of the real test file. Add an `errorsAs` test helper (or just use `errors.As` directly from the standard `errors` package, which is simpler — prefer that over adding a helper) to check error types in the tests above.

- [ ] **Step 2: Run tests to verify they fail**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -run TestService -v`
Expected: FAIL — `Service`/`NewService`/etc. undefined, and/or the placeholder line fails to compile (remove it, per the note above, before running).

- [ ] **Step 3: Write `service.go`**

Implement exactly the methods listed in Interfaces above. Key logic notes:
- `Create`/`CreateFromDefault`: validate `documentType` with `if !ValidDocumentTypes[documentType] { return nil, &InvalidDocumentTypeError{DocumentType: documentType} }` before any DB call. Check code uniqueness via `s.repo.GetByCode`; if it returns a template (no error), return `&DuplicateCodeError{Code: code}`; if it returns `ErrTemplateNotFound`, proceed; any other error propagates.
- `Activate`: use `s.repo.WithTx(ctx, func(tx *gorm.DB) error { ... })`. Inside: `var previous DocumentTemplate; tx.Where("type = ? AND status = ? AND id != ? AND deleted_at IS NULL", target.DocumentType, StatusActive, target.ID).First(&previous)` — if found (no `gorm.ErrRecordNotFound`), set `previous.Status = StatusInactive` and `tx.Save(&previous)`, then write a `DEACTIVATED` audit row for it via `s.repo.CreateAudit(ctx, tx, ...)`. Then set the target's `Status = StatusActive`, `tx.Save(target)`, write an `ACTIVATED` audit row. Return the updated target (re-fetch or just return the in-memory struct after `Save`).
- `CreateVersion`: use `s.repo.WithTx`. Inside: `next, err := s.repo.NextVersionNumber(ctx, tx, templateID)`; build the `DocumentTemplateVersion`; `s.repo.CreateVersion(ctx, tx, v)`; then update the template's `ActiveVersionID` — per the Interfaces note above, do this via `tx.Model(&DocumentTemplate{}).Where("id = ?", templateID).Update("active_version_id", v.ID)` directly (simplest, avoids touching `Repository.Update`'s signature); write a `VERSION_CREATED` audit row referencing `VersionID: &v.ID`.
- Every audit-writing call uses `s.repo.CreateAudit(ctx, tx, &DocumentTemplateAudit{ID: uuid.New().String(), TemplateID: ..., Action: "...", ActorID: &actorID})` — import `github.com/google/uuid` for ID generation (matching the rest of this codebase's ID-generation convention, confirmed in Task 2's models via `BeforeCreate`-free plain string IDs — since these models don't have GORM `BeforeCreate` hooks, IDs must be generated explicitly in the service before `Create`/`CreateVersion`/`CreateAudit`, not left to the DB or a hook).
- Add `uuid.New().String()` ID generation to `Create`/`CreateFromDefault` too — the DTO/handler layer (Task 5) will not generate IDs, so the service must.

- [ ] **Step 4: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -v`
Expected: PASS — all service tests green, including version/audit/activate-deactivate behavior.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/modules/documenttemplate/service.go backend/internal/modules/documenttemplate/service_test.go
git commit -m "feat: add documenttemplate service — activation, from-default copy, versioning"
```

---

### Task 5: DTO, handler, routes, module — wire the API

**Files:**
- Create: `backend/internal/modules/documenttemplate/dto.go`
- Create: `backend/internal/modules/documenttemplate/handler.go`
- Create: `backend/internal/modules/documenttemplate/routes.go`
- Create: `backend/internal/modules/documenttemplate/module.go`
- Create: `backend/internal/pkg/httputil/locale.go` — **modify**, not create (add new keys to the existing `localeMessages` map)
- Test: `backend/internal/modules/documenttemplate/handler_test.go`

**Interfaces:**
- Consumes: `Service` and all its methods from Task 4; `httputil.SuccessJSON/CreatedJSON/UpdatedJSON/DeletedJSON/ErrorJSON/ErrorSimple/NotFound/InternalError/BadRequest/BindAndValidate` (existing package, confirmed signatures — do not guess, read `backend/internal/pkg/httputil/response.go` and `validator.go` before writing handler.go if any signature is unclear).
- Produces: routes under `/api/v1/tenant/document-templates` (mounted via the tenant router the same way `setting`/`employeemovement` are — see `NewModule`'s `RegisterRoutes`).

- [ ] **Step 1: Write `dto.go`**

```go
package documenttemplate

type CreateTemplateRequest struct {
	Name         string `json:"name" binding:"required,max=255"`
	Code         string `json:"code" binding:"required,max=100"`
	DocumentType string `json:"document_type" binding:"required"`
	Description  string `json:"description,omitempty" binding:"max=1000"`
}

type CreateFromDefaultRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	Name         string `json:"name" binding:"required,max=255"`
	Code         string `json:"code" binding:"required,max=100"`
}

type UpdateTemplateRequest struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,max=255"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
}

type UpdateDefaultContentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CreateVersionRequest struct {
	Content      string `json:"content" binding:"required"`
	PaperSize    string `json:"paper_size,omitempty" binding:"omitempty,oneof=A4 A5 Letter Legal"`
	Orientation  string `json:"orientation,omitempty" binding:"omitempty,oneof=portrait landscape"`
	MarginTop    int    `json:"margin_top,omitempty"`
	MarginRight  int    `json:"margin_right,omitempty"`
	MarginBottom int    `json:"margin_bottom,omitempty"`
	MarginLeft   int    `json:"margin_left,omitempty"`
}

type TemplateListResponse struct {
	Data  []DocumentTemplate `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
}
```
Default `PaperSize`/`Orientation`/margins (when omitted) to `"A4"`/`"portrait"`/`20` each in the handler before calling the service, not via zero-value binding tags.

- [ ] **Step 2: Write `handler.go`**

```go
package documenttemplate

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/inthros/hris-platform/internal/pkg/httputil"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func actorID(c *gin.Context) string {
	if v, ok := c.Get("user_id"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (h *Handler) handleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *InvalidDocumentTypeError:
		httputil.ErrorSimple(c, http.StatusBadRequest, e.Error())
	case *DuplicateCodeError:
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_CODE", "documenttemplate.duplicate_code", e.Code)
	case *DuplicateActiveTemplateError:
		httputil.ErrorJSON(c, http.StatusConflict, "DUPLICATE_ACTIVE", "documenttemplate.duplicate_active", e.DocumentType)
	case *ReferenceTemplateImmutableError:
		httputil.ErrorSimple(c, http.StatusForbidden, e.Error())
	default:
		if err == ErrTemplateNotFound || err == ErrVersionNotFound {
			httputil.NotFound(c, err.Error())
			return
		}
		httputil.InternalError(c, err.Error())
	}
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	items, total, err := h.svc.List(c.Request.Context(), page, perPage, c.Query("document_type"), c.Query("status"), c.Query("search"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, TemplateListResponse{Data: items, Total: total, Page: page})
}

func (h *Handler) GetByID(c *gin.Context) {
	tpl, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, tpl)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.Create(c.Request.Context(), req.Name, req.Code, req.DocumentType, req.Description, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, tpl, "documenttemplate.created")
}

func (h *Handler) CreateFromDefault(c *gin.Context) {
	var req CreateFromDefaultRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	tpl, err := h.svc.CreateFromDefault(c.Request.Context(), req.DocumentType, req.Name, req.Code, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, tpl, "documenttemplate.created_from_default")
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateTemplateRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	name, desc := "", ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		desc = *req.Description
	}
	tpl, err := h.svc.Update(c.Request.Context(), c.Param("id"), name, desc, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.updated")
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), actorID(c)); err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.DeletedJSON(c, "documenttemplate.deleted")
}

func (h *Handler) Activate(c *gin.Context) {
	tpl, err := h.svc.Activate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.activated")
}

func (h *Handler) Deactivate(c *gin.Context) {
	tpl, err := h.svc.Deactivate(c.Request.Context(), c.Param("id"), actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.UpdatedJSON(c, tpl, "documenttemplate.deactivated")
}

func (h *Handler) ListVersions(c *gin.Context) {
	versions, err := h.svc.ListVersions(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.SuccessJSON(c, versions)
}

func (h *Handler) CreateVersion(c *gin.Context) {
	var req CreateVersionRequest
	if !httputil.BindAndValidate(c, &req) {
		return
	}
	paperSize, orientation := req.PaperSize, req.Orientation
	if paperSize == "" {
		paperSize = "A4"
	}
	if orientation == "" {
		orientation = "portrait"
	}
	margins := [4]int{req.MarginTop, req.MarginRight, req.MarginBottom, req.MarginLeft}
	for i, m := range margins {
		if m == 0 {
			margins[i] = 20
		}
	}
	v, err := h.svc.CreateVersion(c.Request.Context(), c.Param("id"), req.Content, paperSize, orientation, margins, actorID(c))
	if err != nil {
		h.handleServiceError(c, err)
		return
	}
	httputil.CreatedJSON(c, v, "documenttemplate.version_created")
}
```

- [ ] **Step 3: Write `routes.go`**

```go
package documenttemplate

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	templates := rg.Group("/document-templates")
	{
		templates.GET("", handler.List)
		templates.POST("", handler.Create)
		templates.POST("/from-default", handler.CreateFromDefault)
		templates.GET("/:id", handler.GetByID)
		templates.PUT("/:id", handler.Update)
		templates.DELETE("/:id", handler.Delete)
		templates.POST("/:id/activate", handler.Activate)
		templates.POST("/:id/deactivate", handler.Deactivate)
		templates.GET("/:id/versions", handler.ListVersions)
		templates.POST("/:id/versions", handler.CreateVersion)
	}
}
```
Note the static route `/from-default` is registered before `/:id` — required so Gin doesn't treat `from-default` as an `:id` value (matches this codebase's existing convention, e.g. villages' `/search`).

- [ ] **Step 4: Write `module.go`**

```go
package documenttemplate

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/inthros/hris-platform/internal/pkg/database"
	"github.com/inthros/hris-platform/internal/pkg/module"
)

type documentTemplateModule struct {
	handler *Handler
}

func NewModule(dbManager *database.Manager, logger *zap.Logger) module.Module {
	repo := NewRepository(NewTenantDBResolver(dbManager))
	svc := NewService(repo, logger)
	return &documentTemplateModule{handler: NewHandler(svc)}
}

func (m *documentTemplateModule) Info() module.ModuleInfo {
	return module.ModuleInfo{
		Name:        "Document Template",
		Slug:        "documenttemplate",
		Version:     "1.0.0",
		Description: "Manage reusable document templates for contracts and movement SKs",
		IsCore:      false,
		Permissions: []string{
			"documenttemplate.template.view",
			"documenttemplate.template.create",
			"documenttemplate.template.update",
			"documenttemplate.template.delete",
			"documenttemplate.template.activate",
			"documenttemplate.template.deactivate",
			"documenttemplate.template.set_default",
			"documenttemplate.template.version",
			"documenttemplate.generated.view",
		},
	}
}

func (m *documentTemplateModule) RegisterRoutes(rg *gin.RouterGroup) {
	RegisterRoutes(rg, m.handler)
}

func (m *documentTemplateModule) Migrate(db *gorm.DB) error {
	// Tenant schema for this module is owned entirely by the versioned SQL
	// migration (110_document_templates.sql) — AutoMigrate does not run
	// against tenant databases in this codebase. This method is a required
	// module.Module interface stub, not the real migration path.
	return nil
}

func (m *documentTemplateModule) Seed(db *gorm.DB) error {
	// Default/reference templates are seeded by the SQL migration itself.
	return nil
}

func (m *documentTemplateModule) Permissions() []string {
	return m.Info().Permissions
}
```
Confirm the exact import path prefix (`github.com/inthros/hris-platform/...`) and the `module.Module`/`module.ModuleInfo` field names against `backend/internal/pkg/module/module.go` before finalizing — adjust if anything differs from what Task-gathering research reported.

- [ ] **Step 5: Add locale keys**

In `backend/internal/pkg/httputil/locale.go`, add to the existing `localeMessages` map (following the exact `"setting.duplicate_code": {LangEN: ..., LangID: ...}` shape already in that file):
```go
"documenttemplate.created": {
	LangEN: "Document template created successfully",
	LangID: "Template dokumen berhasil dibuat",
},
"documenttemplate.created_from_default": {
	LangEN: "Document template created from default successfully",
	LangID: "Template dokumen berhasil dibuat dari default",
},
"documenttemplate.updated": {
	LangEN: "Document template updated successfully",
	LangID: "Template dokumen berhasil diperbarui",
},
"documenttemplate.deleted": {
	LangEN: "Document template deleted successfully",
	LangID: "Template dokumen berhasil dihapus",
},
"documenttemplate.activated": {
	LangEN: "Document template activated successfully",
	LangID: "Template dokumen berhasil diaktifkan",
},
"documenttemplate.deactivated": {
	LangEN: "Document template deactivated successfully",
	LangID: "Template dokumen berhasil dinonaktifkan",
},
"documenttemplate.version_created": {
	LangEN: "New template version created successfully",
	LangID: "Versi template baru berhasil dibuat",
},
"documenttemplate.duplicate_code": {
	LangEN: "Template code '%s' is already in use",
	LangID: "Kode template '%s' sudah digunakan",
},
"documenttemplate.duplicate_active": {
	LangEN: "An active template already exists for document type '%s'",
	LangID: "Template aktif sudah ada untuk jenis dokumen '%s'",
},
```

- [ ] **Step 6: Write `handler_test.go`**

Follow Task 4 in the prior document-numbering-settings plan's pattern for a standalone-router handler test (private `:memory:` sqlite, not shared-cache — that DSN caused real cross-test pollution in a prior feature in this codebase, confirmed by direct investigation; use `sqlite.Open(":memory:")`):
```go
package documenttemplate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	db, cleanup := setupTestDB()
	t.Cleanup(cleanup)
	svc := NewService(newTestRepo(db), zap.NewNop())
	handler := NewHandler(svc)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	group := r.Group("/api/v1/tenant")
	RegisterRoutes(group, handler)
	return r
}

func TestHandlerCreateAndList(t *testing.T) {
	r := setupTestRouter(t)

	payload := `{"name":"PKWT","code":"PKWT_H1","document_type":"CONTRACT_AGREEMENT"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/tenant/document-templates", nil)
	listW := httptest.NewRecorder()
	r.ServeHTTP(listW, listReq)
	if listW.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", listW.Code, listW.Body.String())
	}
	var body struct {
		Data struct {
			Data  []DocumentTemplate `json:"data"`
			Total int64              `json:"total"`
		} `json:"data"`
	}
	if err := json.Unmarshal(listW.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.Total != 1 {
		t.Fatalf("expected 1 template, got %d", body.Data.Total)
	}
}

func TestHandlerCreateRejectsInvalidDocumentType(t *testing.T) {
	r := setupTestRouter(t)
	payload := `{"name":"X","code":"XCODE","document_type":"NOT_REAL"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestHandlerFromDefaultRouteNotShadowedByIDParam(t *testing.T) {
	r := setupTestRouter(t)
	payload := `{"document_type":"NOT_REAL","name":"X","code":"XCODE2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenant/document-templates/from-default", bytes.NewReader([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Must reach the from-default handler (which validates document_type and
	// 400s) rather than being swallowed by the GET /:id route as id="from-default".
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 (from-default handler reached, invalid doc type), got %d: %s", w.Code, w.Body.String())
	}
}
```
Check `httputil.SuccessJSON`/`CreatedJSON`'s actual response envelope shape (read `response.go` if the assumed `{"data": {...}}` / `{"success": true, "data": ...}` shape in the test above doesn't match) and adjust the test's unmarshal target accordingly.

- [ ] **Step 7: Run tests to verify they pass**

Run: `cd backend && go build ./... && go test ./internal/modules/documenttemplate/... -v`
Expected: `go build` succeeds (this task doesn't yet register the module in `main.go`, so nothing else should break); all `documenttemplate` package tests pass.

- [ ] **Step 8: Commit**

```bash
git add backend/internal/modules/documenttemplate/dto.go backend/internal/modules/documenttemplate/handler.go backend/internal/modules/documenttemplate/routes.go backend/internal/modules/documenttemplate/module.go backend/internal/modules/documenttemplate/handler_test.go backend/internal/pkg/httputil/locale.go
git commit -m "feat: add documenttemplate DTO, handler, routes, module wiring"
```

---

### Task 6: Register the module in `main.go`, add a static variables registry endpoint

**Files:**
- Modify: `backend/cmd/server/main.go`
- Create: `backend/internal/modules/documenttemplate/variables.go`
- Modify: `backend/internal/modules/documenttemplate/routes.go`
- Modify: `backend/internal/modules/documenttemplate/handler.go`
- Test: `backend/internal/modules/documenttemplate/variables_test.go`

**Interfaces:**
- Produces: `GET /api/v1/tenant/document-templates/variables` returning the static variable registry from spec §7 (category → list of `{key, label}`), registered **before** the `/:id` routes for the same reason as `/from-default`.

- [ ] **Step 1: Write the failing test**

`backend/internal/modules/documenttemplate/variables_test.go`:
```go
package documenttemplate

import "testing"

func TestVariableRegistryHasExpectedCategories(t *testing.T) {
	reg := VariableRegistry()
	wantCategories := []string{"employee", "contract", "movement", "company"}
	for _, cat := range wantCategories {
		found := false
		for _, group := range reg {
			if group.Category == cat {
				found = true
				if len(group.Variables) == 0 {
					t.Errorf("category %q has no variables", cat)
				}
				break
			}
		}
		if !found {
			t.Errorf("expected category %q in registry, not found", cat)
		}
	}
}

func TestVariableRegistryKeysAreDotted(t *testing.T) {
	reg := VariableRegistry()
	for _, group := range reg {
		for _, v := range group.Variables {
			if v.Key == "" || v.Label == "" {
				t.Errorf("variable in category %q has empty key or label: %+v", group.Category, v)
			}
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -run TestVariableRegistry -v`
Expected: FAIL — `VariableRegistry` undefined.

- [ ] **Step 3: Write `variables.go`**

Per spec §7's exact variable list:
```go
package documenttemplate

type TemplateVariable struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type VariableGroup struct {
	Category  string             `json:"category"`
	Variables []TemplateVariable `json:"variables"`
}

// VariableRegistry is the static list of placeholders available to the
// template editor's "Insert Variable" picker (spec §7). Backend-owned so the
// frontend never has to keep this list in sync by hand.
func VariableRegistry() []VariableGroup {
	return []VariableGroup{
		{
			Category: "employee",
			Variables: []TemplateVariable{
				{Key: "employee.name", Label: "Name"},
				{Key: "employee.nik", Label: "NIK"},
				{Key: "employee.position", Label: "Position"},
				{Key: "employee.organization", Label: "Organization"},
			},
		},
		{
			Category: "contract",
			Variables: []TemplateVariable{
				{Key: "contract.number", Label: "Contract Number"},
				{Key: "contract.start_date", Label: "Start Date"},
				{Key: "contract.end_date", Label: "End Date"},
			},
		},
		{
			Category: "movement",
			Variables: []TemplateVariable{
				{Key: "movement.number", Label: "Movement Number"},
				{Key: "movement.effective_date", Label: "Effective Date"},
				{Key: "movement.previous_position", Label: "Previous Position"},
				{Key: "movement.new_position", Label: "New Position"},
			},
		},
		{
			Category: "company",
			Variables: []TemplateVariable{
				{Key: "company.name", Label: "Name"},
				{Key: "company.address", Label: "Address"},
			},
		},
	}
}
```

- [ ] **Step 4: Add the handler and route**

In `handler.go`, add:
```go
func (h *Handler) ListVariables(c *gin.Context) {
	httputil.SuccessJSON(c, VariableRegistry())
}
```
In `routes.go`, add `templates.GET("/variables", handler.ListVariables)` **before** `templates.GET("/:id", handler.GetByID)` in the registration block.

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/modules/documenttemplate/... -v`
Expected: PASS — all package tests green, including the new variable-registry tests.

- [ ] **Step 6: Register the module in `main.go`**

Read `backend/cmd/server/main.go` around line 1210 (where `setting.NewModule(dbManager, l, numberingSvc)` is registered at `Priority: 16`) to confirm current exact surrounding code, then add immediately after it (or wherever fits without renumbering existing entries):
```go
module.ModuleRegistration{
	Module:   documenttemplate.NewModule(dbManager, l),
	TargetDB: module.TargetTenant,
	Priority: 20,
},
```
Add the `"github.com/inthros/hris-platform/internal/modules/documenttemplate"` import (matching this file's existing import path convention — copy the exact prefix from a neighboring module import rather than assuming).

- [ ] **Step 7: Verify the whole backend builds and all relevant tests pass**

Run: `cd backend && go build ./...`
Expected: clean, zero errors.

Run: `cd backend && go test ./internal/modules/documenttemplate/... ./internal/pkg/httputil/... -v`
Expected: all PASS, no regressions in `httputil` from the locale-key addition.

- [ ] **Step 8: Commit**

```bash
git add backend/cmd/server/main.go backend/internal/modules/documenttemplate/variables.go backend/internal/modules/documenttemplate/routes.go backend/internal/modules/documenttemplate/handler.go backend/internal/modules/documenttemplate/variables_test.go
git commit -m "feat: register documenttemplate module and add variables registry endpoint"
```

---

## Post-Implementation Checklist

- [ ] All backend tests pass: `cd backend && go test ./...` (the pre-existing `internal/pkg/migrator` integration test failure requiring a live DB connection is a known, unrelated environment limitation — not a regression from this work; confirm no other package fails)
- [ ] Backend builds cleanly: `cd backend && go build ./...`
- [ ] Migration 110 applied and verified on a real/test tenant DB (both mysql and postgres) if a DB is reachable in the execution environment; otherwise confirmed via careful re-read
- [ ] Spec §21 "Phase 1" checklist items all covered: module + registration ✓, models + tables ✓, SQL migrator with partial index + seed ✓, tenant resolver + repository ✓, service (one-active-per-type, default copy, versioning) ✓, DTO + validation ✓, API handler + routes ✓, permissions registered ✓, httputil locale keys ✓
- [ ] Out of scope for this plan (confirm nothing here was accidentally started): PrimeVue Editor UI, chromedp/PDF generation, preview endpoints, Contract/Movement integration, generated-document write path — all later phases
