#!/usr/bin/env python3
"""Generate docs/database-schema.md (database structure + ERD) from migrations.

Can also be imported by scripts/check_database_schema_doc.py — build_doc()
returns the document content without writing any file.
"""

import os
import re
import io
from collections import OrderedDict

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
BASE_DIR = os.path.dirname(SCRIPT_DIR)  # root repo
TENANT_DIR = os.path.join(BASE_DIR, "backend/internal/pkg/migrator/migrations/tenant/postgres")
MYSQL_TENANT_DIR = os.path.join(BASE_DIR, "backend/internal/pkg/migrator/migrations/tenant/mysql")
PLATFORM_DIR = os.path.join(BASE_DIR, "backend/internal/pkg/migrator/migrations/platform")
OUTPUT = os.path.join(BASE_DIR, "docs/database-schema.md")


_SQL_CLAUSE_KEYWORDS = (
    "CONSTRAINT", "PRIMARY", "UNIQUE", "FOREIGN", "INDEX",
    "KEY", "REFERENCES", "CHECK", "FULLTEXT", "SPATIAL",
)

# Cocok dengan keyword klausa SQL di awal part. Batas kata (\(\w|\s|\() mencegah
# kolom legitimate yang diawali kata serupa (mis. index_fingerprint, key_label)
# ikut ter-skip.
_SQL_CLAUSE_START = re.compile(r"(?:%s)[\s(]" % "|".join(_SQL_CLAUSE_KEYWORDS), re.I)


def split_column_line(line):
    """Pecah baris yang memuat beberapa definisi kolom dalam satu baris.

    MySQL terkadang menulis beberapa kolom dalam satu baris (mis. migrasi 001:
    `created_at ... , updated_at ...`). Koma di dalam tipe ber-parens
    (mis. DECIMAL(18,2), ENUM('a','b')) tidak boleh dipecah.
    """
    parts = []
    depth = 0
    cur = []
    in_str = None
    for ch in line:
        if in_str:
            if ch == in_str:
                in_str = None
            cur.append(ch)
            continue
        if ch in ("'", '"'):
            in_str = ch
            cur.append(ch)
        elif ch == "(":
            depth += 1
            cur.append(ch)
        elif ch == ")":
            depth = max(0, depth - 1)
            cur.append(ch)
        elif ch == "," and depth == 0:
            parts.append("".join(cur))
            cur = []
        else:
            cur.append(ch)
    parts.append("".join(cur))
    return parts


def parse_create_tables(path):
    out = {}
    try:
        txt = io.open(path, encoding="utf-8", errors="ignore").read()
    except OSError:
        return out
    # Menangani dua format: postgres `) ;` dan mysql `) ENGINE=InnoDB ... ;`
    # Regex juga menangkap CREATE TABLE tanpa IF NOT EXISTS (mis. 054_create_job_management_value_clusters)
    for m in re.finditer(r"CREATE TABLE (?:IF NOT EXISTS )?(\w+) \((.*?)\)\s*(?:ENGINE[^;]*)?;", txt, re.S):
        name = m.group(1)
        body = m.group(2)
        cols = []
        fks = []
        for line in body.splitlines():
            line = line.strip()
            if not line:
                continue
            if line.startswith(("CONSTRAINT", "PRIMARY", "UNIQUE", "FOREIGN", "INDEX")):
                fm = re.search(r"FOREIGN KEY \((\w+)\)\s+REFERENCES\s+(\w+)\((\w+)\)", line)
                if fm:
                    fks.append((fm.group(1), fm.group(2), fm.group(3)))
                continue
            # Baris kolom (mungkin berisi beberapa kolom sekaligus di mysql)
            for part in split_column_line(line):
                part = part.strip()
                if not part:
                    continue
                if _SQL_CLAUSE_START.match(part):
                    continue
                # Nama kolom boleh di-backtick (mysql `x`) atau double-quote
                # (postgres "group"), tipe boleh ber-parens/string
                cm = re.match(r"^[`\"]?(\w+)[`\"]?\s+([\w(),']+)", part)
                if cm:
                    cols.append((cm.group(1), cm.group(2).strip()))
        out[name] = {"cols": cols, "fks": fks}

    # Inline column-level REFERENCES di dalam body CREATE TABLE
    # (mis. migrasi 066: performance_evaluation_id ... REFERENCES performance_evaluations(id))
    for tname, info in out.items():
        mb = re.search(
            r"CREATE TABLE (?:IF NOT EXISTS )?" + re.escape(tname) + r" \((.*?)\)\s*(?:ENGINE[^;]*)?;",
            txt, re.S,
        )
        if not mb:
            continue
        body = mb.group(1)
        for col, _ in info["cols"]:
            fm = re.search(
                r"\b" + re.escape(col) + r"\b[^,()]*(?:\([^,()]*\)[^,()]*)*?\bREFERENCES\s+(\w+)\s*\((\w+)\)",
                body, re.I,
            )
            if fm:
                fk = (col, fm.group(1), fm.group(2))
                if fk not in info["fks"]:
                    info["fks"].append(fk)
    return out


_ALTER_FK_RE = re.compile(
    r"ALTER TABLE\s+(\w+)\s+ADD\s+(?:CONSTRAINT\s+\w+\s+)?FOREIGN KEY\s*\((\w+)\)\s*REFERENCES\s+(\w+)\s*\((\w+)\)",
    re.I,
)


def parse_alter_fks(txt):
    """FK yang dideklarasikan via ALTER TABLE ... ADD CONSTRAINT FOREIGN KEY
    (mis. migrasi 058/059 memakai blok DO $$ ... END $$). Dipanggil di load_all
    setelah semua file di-merge agar FK antar-file tidak terlewat."""
    return [(m.group(1), m.group(2), m.group(3), m.group(4)) for m in _ALTER_FK_RE.finditer(txt)]


def parse_alter_add_columns(txt):
    """Kolom yang ditambahkan via ALTER TABLE ... ADD COLUMN.

    Banyak migrasi (055, 060-069, dll.) menambah kolom dengan ALTER TABLE, bukan
    di CREATE TABLE. Karena load_all memanggil ini per file setelah tabel di-merge,
    kolom hasil ALTER ikut tercatat di dokumen. Mendukung postgres (multi-baris,
    IF NOT EXISTS) dan mysql (ADD COLUMN ... AFTER col).
    """
    out = []  # (tabel, nama_kolom, tipe)
    for m in re.finditer(r"ALTER TABLE\s+(\w+)\s*([^;]*);", txt, re.S | re.I):
        tname = m.group(1)
        body = m.group(2)
        # 'COLUMN' wajib: menolak false positive dari 'ADD CONSTRAINT ... FOREIGN KEY'
        # (postgres: 'ADD COLUMN IF NOT EXISTS x TYPE', mysql: 'ADD COLUMN x TYPE AFTER y')
        for cm in re.finditer(r"ADD\s+COLUMN\s+(?:IF NOT EXISTS\s+)?(\w+)\s+([\w(),]+)", body, re.I):
            out.append((tname, cm.group(1), cm.group(2)))
    return out


def load_all(directory):
    """Load tables in migration order, honoring DROP TABLE statements."""
    tables = {}
    for fn in sorted(os.listdir(directory)):
        if fn.endswith(".down.sql"):
            continue
        path = os.path.join(directory, fn)
        parsed = parse_create_tables(path)
        tables.update(parsed)
        # Tabel yang di-DROP di migrasi ini harus dihapus dari hasil akhir
        try:
            txt = io.open(path, encoding="utf-8", errors="ignore").read()
        except OSError:
            continue
        for m in re.finditer(r"DROP TABLE IF EXISTS (\w+)", txt):
            tables.pop(m.group(1), None)
        # FK via ALTER TABLE (mungkin di file yang berbeda dari CREATE TABLE)
        for tname, col, ref_t, ref_c in parse_alter_fks(txt):
            if tname in tables:
                tables[tname]["fks"].append((col, ref_t, ref_c))
        # Kolom via ALTER TABLE ... ADD COLUMN
        for tname, col, ctype in parse_alter_add_columns(txt):
            if tname in tables and col not in [c for c, _ in tables[tname]["cols"]]:
                tables[tname]["cols"].append((col, ctype))
    # Buang FK yang menunjuk ke tabel yang sudah tidak ada (mis. countries), dan dedup
    for name, info in tables.items():
        fks = []
        for fk in info["fks"]:
            if fk[1] in tables and fk not in fks:
                fks.append(fk)
        info["fks"] = fks
    return tables


def clean_type(t):
    base = re.sub(r"\(.*$", "", t.upper())
    return base or "TEXT"


def render_erd(tables, table_names, title):
    lines = [f"### {title}", "", "```mermaid", "erDiagram"]
    for t in table_names:
        if t not in tables:
            continue
        lines.append(f"    {t} {{")
        for col, ctype in tables[t]["cols"]:
            lines.append(f"        {clean_type(ctype)} {col}")
        lines.append("    }")
    for t in table_names:
        if t not in tables:
            continue
        for col, ref_t, ref_c in tables[t]["fks"]:
            if ref_t in table_names:
                lines.append(f'    {ref_t} ||--o{{ {t} : "{col}"')
    lines.append("```")
    lines.append("")
    return "\n".join(lines)


def table_markdown(tables, table_names):
    """Markdown table: tabel | kolom | relasi FK utama"""
    rows = ["| Tabel | Jumlah Kolom | FK Utama |", "|---|---|---|"]
    for t in table_names:
        if t not in tables:
            continue
        info = tables[t]
        fks = ", ".join(f"{c}->{rt}" for c, rt, _ in info["fks"]) or "-"
        rows.append(f"| `{t}` | {len(info['cols'])} | {fks} |")
    return "\n".join(rows)


# ---------------------------------------------------------------------------
# Load
# ---------------------------------------------------------------------------
tenant = load_all(TENANT_DIR)
platform = load_all(PLATFORM_DIR)

# Tabel platform yang dibuat via GORM AutoMigrate (bukan file SQL migration)
platform["packages"] = {
    "cols": [
        ("id", "char(36)"), ("name", "varchar(255)"), ("slug", "varchar(100)"),
        ("description", "text"), ("price", "decimal(15,2)"), ("status", "varchar(20)"),
        ("is_public", "boolean"), ("sort_order", "int"),
        ("created_at", "timestamp"), ("updated_at", "timestamp"), ("deleted_at", "timestamp"),
    ],
    "fks": [],
}
platform["package_modules"] = {
    "cols": [
        ("package_id", "char(36)"), ("module_id", "char(36)"), ("is_mandatory", "boolean"),
        ("sort_order", "int"), ("created_at", "timestamp"),
        ("module_name", "varchar(255)"), ("module_slug", "varchar(100)"),
    ],
    "fks": [("package_id", "packages", "id"), ("module_id", "modules", "id")],
}

# ---------------------------------------------------------------------------
# Module groups (order & table membership)
# ---------------------------------------------------------------------------
MODULES = OrderedDict([
    ("Master Data & Settings", [
        "religions", "educations", "education_majors", "marital_statuses",
        "provinces", "regencies", "districts", "villages",
        "relationship_types", "employment_statuses", "gradings", "job_families",
        "banks", "nationalities", "salary_grades", "ptkps", "ters",
        "insurances", "company_holidays", "document_templates",
    ]),
    ("Organization", [
        "organization_summaries", "organization_levels", "zones",
        "organizations", "positions", "job_family_competencies",
    ]),
    ("Employee", [
        "employees", "employments", "employee_addresses", "emergency_contacts",
        "employee_families", "employee_educations", "employee_experiences",
        "employee_documents", "employee_insurances", "employee_bank_accounts",
    ]),
    ("Attendance", [
        "attendance_company_settings", "attendance_company_shifts",
        "attendance_employee_shifts", "attendance_locations",
        "attendance_device_captures", "attendance_face_captures",
        "attendance_events", "attendance_sessions",
        "attendance_overtime_requests", "attendance_exempt_positions",
    ]),
    ("Leave", [
        "leave_types", "leave_accrual_policies", "leave_reasons", "leave_requests",
        "leave_request_details", "employee_leave_balances",
    ]),
    ("Payroll", [
        "salary_components", "salary_grade_components", "salary_employee_components",
        "salary_change_logs", "salary_employee_adjustments", "payroll_periods",
        "employee_payroll_profiles", "employee_bank_profiles",
        "employee_bpjs_profiles", "employee_tax_profiles",
        "bpjs_settings", "bpjs_rate_components", "pph21_settings",
        "pph21_ptkp_rates", "pph21_tax_brackets",
        "payroll_runs", "payroll_run_employees", "payroll_run_items",
        "payroll_payslips", "pph21_calculation_logs", "payroll_profile_change_logs",
    ]),
    ("Competency", [
        "competencies", "competence_values", "competency_values",
        "competency_events", "competency_event_targets",
        "competency_scores", "competency_score_details",
    ]),
    ("Job Management", [
        "job_management_titles", "job_management_title_subs",
        "job_management_values", "job_management_objectives",
        "job_management_identifications", "job_management_responsibilities",
        "job_management_education_experiences", "job_management_hr_authorities",
        "job_management_operational_authorities", "job_management_working_activities",
        "job_management_working_risks", "job_management_relationships",
        "job_management_subordinate_controls", "job_management_assets",
        "job_management_financials", "job_management_potency_competencies",
        "job_management_scores", "job_management_competency_groups",
        "job_management_value_clusters", "job_management_majors",
        "job_management_job_family", "job_management_relationship_details",
    ]),
    ("Approval Engine", [
        "approval_flows", "approval_flow_steps", "approval_instances",
        "approval_actions", "approval_tasks", "approval_flow_step_organizations",
    ]),
    ("RBAC & Auth", [
        "features", "permissions", "roles", "feature_permission",
        "model_has_roles", "model_has_permissions", "role_has_permissions",
        "users", "employee_accounts",
    ]),
    ("Employee Movement", [
        "employee_movements", "employee_contracts",
    ]),
    ("Reimbursement", [
        "reimbursement_types", "reimbursement_requests", "reimbursement_items",
    ]),
    ("Performance — KPI", [
        "performance_periods", "performance_perspectives",
        "performance_templates", "performance_indicators",
        "performance_evaluations", "performance_evaluation_details",
        "performance_components", "performance_organization_components",
        "performance_evaluation_components", "performance_evaluation_program_items",
        "performance_targets", "performance_ratings",
        "performance_indicator_formulas", "performance_progress",
        "performance_comments", "performance_attachments", "performance_logs",
    ]),
    ("Performance — OKR", [
        "okr_templates", "okr_objectives", "okr_key_results", "okr_evaluations",
        "okr_evaluation_details", "okr_progress", "okr_comments", "okr_attachments",
    ]),
    ("Recruitment", [
        "job_requisitions", "candidates", "job_applications", "interviews",
        "onboarding_task_templates", "employee_onboardings", "onboarding_task_items",
    ]),
    ("Training", [
        "training_categories", "training_courses", "training_sessions",
        "training_participants", "training_materials",
        "training_evaluations", "training_certificates",
    ]),
    ("Workforce Intelligence", [
        "workforce_planning_headcounts", "workforce_forecasts", "workforce_kpis",
        "workforce_analytics_cache", "workforce_scenarios",
        "workforce_risk_indicators", "workforce_health_scores",
    ]),
    ("Career Intelligence", [
        "career_talent_maps", "career_interests", "career_paths",
        "career_succession_plans",
    ]),
])

PLATFORM_TABLES = [
    "companies", "tenant_connections", "platform_users", "modules",
    "company_modules", "licenses", "rbac_roles", "rbac_permissions",
    "rbac_role_permissions", "packages", "package_modules",
]


# ---------------------------------------------------------------------------
# Build document
# ---------------------------------------------------------------------------
def build_doc():
    doc = []
    doc.append("# HRIS Platform — Struktur Database & ERD")
    doc.append("")
    doc.append("> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md)  ")
    doc.append("> **Terkait:** [`platform-architecture-design.md`](platform-architecture-design.md) · [`api/api-usage-guide.md`](api/api-usage-guide.md) · [`go-module-architecture-report.md`](go-module-architecture-report.md)")
    doc.append("")
    doc.append("Dokumen ini menjelaskan struktur database HRIS Platform: arsitektur dua-database (Platform & Tenant), skema tabel, relasi (ERD), dan konvensi kolom.")
    doc.append("")
    doc.append("## Ringkasan")
    doc.append("")
    doc.append("| Database | Isi | Jumlah Tabel |")
    doc.append("|---|---|---|")
    doc.append(f"| **Platform DB** | Data multi-tenant: companies, modul, lisensi, paket, RBAC platform | {len(PLATFORM_TABLES)} |")
    doc.append(f"| **Tenant DB** (1 per company) | Seluruh data HRIS milik satu company | {len(tenant)} |")
    doc.append("")
    doc.append("> Sumber kebenaran: file migrasi SQL di `backend/internal/pkg/migrator/migrations/` (dialect `postgres/` & `mysql/` identik).")
    doc.append("")

    # ---- Platform section ----
    doc.append("## Platform Database")
    doc.append("")
    doc.append("Database pusat penyedia SaaS (Platform Central DB). Menyimpan data **perusahaan/tenant**, kredensial koneksi tenant, katalog modul, lisensi, paket, dan RBAC platform.")
    doc.append("")
    doc.append("### Tabel Platform")
    doc.append("")
    doc.append(table_markdown(platform, PLATFORM_TABLES))
    doc.append("")
    doc.append(render_erd(platform, PLATFORM_TABLES, "ERD Platform Database"))
    doc.append("")

    # ---- Tenant sections ----
    doc.append("## Tenant Database")
    doc.append("")
    doc.append("Satu database terisolasi **per company** (database per tenant). Struktur identik untuk semua tenant (dibuat saat provisioning company).")
    doc.append("")
    doc.append(f"Total **{len(tenant)} tabel** dikelompokkan dalam {len(MODULES)} modul:")
    doc.append("")

    # module summary table
    doc.append("| # | Modul | Jumlah Tabel |")
    doc.append("|---|---|---|")
    for i, (mod, tnames) in enumerate(MODULES.items(), 1):
        present = [t for t in tnames if t in tenant]
        doc.append(f"| {i} | {mod} | {len(present)} |")
    doc.append("")
    doc.append("---")
    doc.append("")

    for mod, tnames in MODULES.items():
        present = [t for t in tnames if t in tenant]
        doc.append(f"## {mod}")
        doc.append("")
        doc.append(table_markdown(tenant, present))
        doc.append("")
        doc.append(render_erd(tenant, present, f"ERD — {mod}"))
        doc.append("")

    # ---- conventions ----
    doc.append("## Konvensi Kolom")
    doc.append("")
    doc.append("| Konvensi | Keterangan |")
    doc.append("|---|---|")
    doc.append("| **Primary Key** | `id CHAR(36)` — UUID v4 untuk semua tabel |")
    doc.append("| **Foreign Key** | Kolom berakhiran `_id` merujuk PK tabel lain (relasi ditampilkan di ERD) |")
    doc.append("| **Soft Delete** | `deleted_at TIMESTAMP NULL` — record tidak dihapus fisik (kecuali tabel tertentu) |")
    doc.append("| **Timestamps** | `created_at` & `updated_at TIMESTAMP` pada hampir semua tabel |")
    doc.append("| **Audit** | `created_by` / `updated_by CHAR(36)` pada tabel master data |")
    doc.append("| **Enum via VARCHAR** | Status & tipe disimpan sebagai string (mis. `status VARCHAR(20)`) |")
    doc.append("")
    doc.append("## Migrasi & Dialect")
    doc.append("")
    up_files = len([f for f in os.listdir(TENANT_DIR) if f.endswith(".sql") and not f.endswith(".down.sql")])
    down_files = len([f for f in os.listdir(TENANT_DIR) if f.endswith(".down.sql")])
    doc.append(f"- Migrasi tenant tersedia untuk **PostgreSQL** (`postgres/`) dan **MySQL** (`mysql/`) — {up_files} file up + {down_files} file down per dialect ({up_files * 2 + down_files * 2} total).")
    doc.append("- Migrasi **platform** bersifat cross-dialect di `migrations/platform/`.")
    doc.append("- Tabel tambahan platform (`packages`, `package_modules`) dibuat via GORM `AutoMigrate`.")
    doc.append("- Dijalankan otomatis saat **provisioning company** (tenant DB dibuat + migrasi + seed).")
    doc.append("")
    doc.append("### Regenerasi Dokumen")
    doc.append("")
    doc.append("Dokumen ini di-generate dari migrasi SQL. Untuk regenerasi atau verifikasi sinkronisasi:")
    doc.append("")
    doc.append("```bash")
    doc.append("python scripts/generate_database_schema_doc.py   # regenerasi (tulis file)")
    doc.append("python scripts/check_database_schema_doc.py     # verifikasi sinkron (tanpa menimpa)")
    doc.append("```")
    doc.append("")
    doc.append("> Verifikasi sinkronisasi mencakup **kedua dialect** (PostgreSQL & MySQL) — `check-db-docs` memastikan jumlah & daftar tabel di dokumen sama dengan migrasi `postgres/` **dan** `mysql/`.")
    doc.append("")

    return "\n".join(doc)


def main():
    content = build_doc()
    with io.open(OUTPUT, "w", encoding="utf-8", newline="") as f:
        f.write(content)
    print(f"Generated: {OUTPUT}")
    print(f"  Platform tables: {len(PLATFORM_TABLES)}")
    print(f"  Tenant tables: {len(tenant)}")
    print(f"  Modules: {len(MODULES)}")


if __name__ == "__main__":
    main()
