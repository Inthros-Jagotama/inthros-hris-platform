#!/usr/bin/env python3
"""Inject Performance KPI/OKR endpoints into openapi.json (idempotent).

Refactor Performance module menempatkan endpoint KPI di bawah prefix
`/performance/kpi/*` dan menambahkan sub-modul OKR di `/performance/okr/*`.
Script ini menyinkronkan openapi.json dengan routes.go:

1. MIGRASI — pindahkan dokumen endpoint performance lama (tanpa prefix) ke
   `/performance/kpi/*` (operationId di-rename: Performance → Kpi).
2. KPI BARU — tambahkan endpoint KPI yang belum terdokumentasi (dashboard,
   workflow submit/approve/reject/complete, recalculate, snapshot, dll).
3. KPI COMPONENTS (Phase 5) — tambahkan endpoint components,
   organization-components, dan scoring engine + schemas.
4. OKR BARU — tambahkan seluruh endpoint sub-modul OKR + schemas.
5. HAPUS PHANTOM — hapus path lama yang sudah dipindah (tidak lagi terdaftar
   di routes.go).

CATATAN: skema yang didefinisikan di `new_schemas` bersifat kanonik dan
selalu disinkronkan (overwrite jika definisi berubah). Jangan edit skema
KPI/OKR langsung di openapi.json — ubah definisi di `new_schemas`.

Usage:
    python scripts/inject_performance_kpi_okr_openapi.py

Usage:
    python scripts/inject_performance_kpi_okr_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]
TAG = "Tenant: Performance Management"

# Path performance yang MASIH valid (terdaftar di routes.go tanpa prefix /kpi/)
KEEP_SUFFIXES = {
    "periods", "periods/{id}",
    "ratings", "ratings/{id}",
    "indicator-formulas", "indicator-formulas/{id}",
    "logs", "logs/{id}",
}


def param(name, where="path", required=True, schema=None):
    return {"name": name, "in": where, "required": required,
            "schema": schema or {"type": "string", "format": "uuid"}}


def qparam(name, schema=None):
    return param(name, where="query", required=False, schema=schema or {"type": "string"})


def ref(name):
    return {"$ref": f"#/components/schemas/{name}"}


def content(schema_ref, required=True):
    return {"required": required, "content": {"application/json": {"schema": ref(schema_ref)}}}


def responses_created(schema_ref):
    return {
        "400": {"description": "Bad request / Validation error"},
        "201": {"description": "Resource created successfully", "content": {"application/json": {"schema": ref(schema_ref)}}},
    }


def responses_ok(schema_ref):
    return {
        "200": {"description": "OK", "content": {"application/json": {"schema": ref(schema_ref)}}},
        "400": {"description": "Bad request / Validation error"},
        "404": {"description": "Resource not found"},
    }


def responses_list(schema_ref):
    return {
        "200": {"description": "Paginated list of resources", "content": {"application/json": {"schema": ref(schema_ref)}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_array(item_ref, desc="List of resources"):
    return {
        "200": {"description": desc, "content": {"application/json": {"schema": {"type": "object", "properties": {
            "success": {"type": "boolean", "example": True},
            "data": {"type": "array", "items": ref(item_ref)}}}}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_plain(desc_200="OK"):
    return {
        "200": {"description": desc_200},
        "400": {"description": "Bad request / Validation error"},
        "404": {"description": "Resource not found"},
    }


def op(oid, summary, description, tag, security=True, parameters=None, request_body=None, responses=None):
    o = {
        "operationId": oid,
        "summary": summary,
        "description": description,
        "tags": [tag],
        "parameters": parameters or [],
        "responses": responses or responses_plain(),
    }
    if security:
        o["security"] = list(SEC)
    if request_body is not None:
        o["requestBody"] = request_body
    return o


def migrate_operation_id(oid):
    """Rename operationId saat migrasi performance → kpi."""
    if "Performance" in oid:
        return oid.replace("Performance", "Kpi")
    if "Evaluation" in oid:
        return oid.replace("Evaluation", "KpiEvaluation", 1)
    return oid


with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

paths = spec["paths"]
schemas = spec["components"]["schemas"]
added = 0

# =========================================================================
# 1. MIGRASI — performance/* → performance/kpi/*
# =========================================================================
phantom_paths = []
for p in list(paths.keys()):
    if not p.startswith("/api/v1/tenant/performance/"):
        continue
    suffix = p[len("/api/v1/tenant/performance/"):]
    if suffix in KEEP_SUFFIXES:
        continue  # masih valid, jangan dipindah
    if suffix.startswith(("kpi/", "okr/")):
        continue  # sudah ber-prefix (struktur baru), jangan dipindah lagi
    phantom_paths.append(p)

for old in phantom_paths:
    new = old.replace("/performance/", "/performance/kpi/", 1)
    if new in paths:
        # Path tujuan sudah ada — pindahkan method yang belum ada
        for method, o in paths[old].items():
            if method == "parameters":
                continue
            if method not in paths[new]:
                paths[new][method] = json.loads(json.dumps(o))
                if "operationId" in paths[new][method]:
                    paths[new][method]["operationId"] = migrate_operation_id(paths[new][method]["operationId"])
                added += 1
    else:
        migrated = {}
        for method, o in paths[old].items():
            if method == "parameters":
                migrated[method] = json.loads(json.dumps(o))
                continue
            entry = json.loads(json.dumps(o))
            if "operationId" in entry:
                entry["operationId"] = migrate_operation_id(entry["operationId"])
            migrated[method] = entry
            added += 1
        paths[new] = migrated
    del paths[old]
    print(f"  migrated {old}  ->  {new}")

# =========================================================================
# 2. KPI BARU — endpoint yang belum pernah didokumentasikan
# =========================================================================
B = "/api/v1/tenant/performance/kpi"


def add_path(path, ops):
    global added
    if path not in paths:
        paths[path] = {}
    for method, o in ops.items():
        if method not in paths[path]:
            paths[path][method] = o
            added += 1


# --- Dashboards ---
add_path(f"{B}/dashboard/employee/{{employee_id}}", {
    "get": op("getKpiEmployeeDashboard", "Get KPI employee dashboard",
              "Dashboard KPI untuk employee: daftar KPI, progress, achievement, dan skor evaluasi periode aktif.",
              TAG, parameters=[param("employee_id")], responses=responses_plain("Employee KPI dashboard data")),
})
add_path(f"{B}/dashboard/manager/{{manager_id}}", {
    "get": op("getKpiManagerDashboard", "Get KPI manager dashboard",
              "Dashboard KPI untuk manager: ringkasan tim, anggota tim beserta skor & status KPI masing-masing.",
              TAG, parameters=[param("manager_id")], responses=responses_plain("Manager KPI dashboard data")),
})
add_path(f"{B}/dashboard/hr", {
    "get": op("getKpiHRDashboard", "Get KPI HR dashboard",
              "Dashboard KPI untuk HR: ringkasan evaluasi seluruh organisasi (total, status, skor rata-rata, distribusi).",
              TAG, responses=responses_plain("HR KPI dashboard data")),
})

# --- KPI Templates (organization scope) ---
add_path(f"{B}/templates/organization-scope", {
    "get": op("listKpiTemplateOrganizationScope", "List KPI template organization scope",
              "Ambil daftar organisasi yang boleh dipilih saat membuat/mengedit KPI template — hanya organisasi turunan dari organisasi user (hierarki org), organisasi milik user sendiri tidak disertakan.",
              TAG, responses=responses_array("OrganizationOptionResponse", "List of selectable organizations for KPI template")),
})

# --- Evaluation snapshot & full ---
add_path(f"{B}/evaluations/snapshot", {
    "post": op("createKpiEvaluationWithSnapshot", "Create KPI evaluation with snapshot",
               "Buat evaluasi KPI baru sekaligus snapshot KPI dari template ke evaluation details (nilai target terkunci).",
               TAG, request_body=content("CreatePerformanceEvaluationRequest"),
               responses=responses_created("PerformanceEvaluationResponse")),
})
add_path(f"{B}/evaluations/{{id}}/full", {
    "get": op("getKpiEvaluationWithDetails", "Get KPI evaluation with details",
              "Ambil evaluasi KPI lengkap: detail perspektif, target, progress, komentar, dan lampiran.",
              TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})
add_path(f"{B}/evaluations/{{id}}/progress-summary", {
    "get": op("getKpiEvaluationProgressSummary", "Get KPI evaluation progress summary",
              "Ringkasan progres realisasi KPI untuk sebuah evaluasi (nilai aktual terbaru & achievement per detail).",
              TAG, parameters=[param("id")], responses=responses_plain("Evaluation progress summary")),
})

# --- Score calculation ---
add_path(f"{B}/evaluations/{{id}}/recalculate", {
    "post": op("recalculateKpiEvaluationScore", "Recalculate KPI evaluation score",
               "Hitung ulang skor evaluasi KPI dari nilai aktual & formula tiap indikator, lalu simpan hasilnya.",
               TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})

# --- Workflow status transitions ---
add_path(f"{B}/evaluations/{{id}}/submit", {
    "post": op("submitKpiEvaluation", "Submit KPI evaluation",
               "Ajukan evaluasi KPI (status → SUBMITTED) untuk direview atasan.",
               TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})
add_path(f"{B}/evaluations/{{id}}/approve", {
    "post": op("approveKpiEvaluation", "Approve KPI evaluation",
               "Setujui evaluasi KPI (status → APPROVED).",
               TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})
add_path(f"{B}/evaluations/{{id}}/reject", {
    "post": op("rejectKpiEvaluation", "Reject KPI evaluation",
               "Tolak evaluasi KPI (status → REJECTED) untuk direvisi.",
               TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})
add_path(f"{B}/evaluations/{{id}}/complete", {
    "post": op("completeKpiEvaluation", "Complete KPI evaluation",
               "Selesaikan evaluasi KPI (status → COMPLETED) — hasil akhir terkunci.",
               TAG, parameters=[param("id")], responses=responses_ok("PerformanceEvaluationResponse")),
})

# --- Actuals input ---
add_path(f"{B}/evaluation-details/{{id}}/actual", {
    "put": op("updateKpiEvaluationDetailActual", "Update KPI evaluation detail actual",
              "Input nilai aktual (actual) untuk satu detail evaluasi — skor & achievement dihitung ulang.",
              TAG, parameters=[param("id")], request_body=content("UpdateEvaluationDetailRequest"),
              responses=responses_ok("PerformanceEvaluationDetailResponse")),
})
add_path(f"{B}/evaluations/{{id}}/actuals", {
    "put": op("bulkUpdateKpiEvaluationActuals", "Bulk update KPI evaluation actuals",
              "Input nilai aktual sekaligus untuk banyak evaluation detail dalam satu request.",
              TAG, parameters=[param("id")], request_body=content("BulkUpdateEvaluationActualsRequest"),
              responses=responses_plain("Evaluation actuals updated")),
})

# =========================================================================
# 3. KPI COMPONENTS (Phase 5 — Scoring Configuration)
# =========================================================================

# --- Performance Components (master data) ---
# CATATAN: sejak KPI Phase 1 (commit 6b84dce, lock components) routes.go hanya
# mendaftarkan GET /components, GET/PUT /components/{id} — POST & DELETE sudah
# tidak ada di backend (komponen terkunci setelah dipakai evaluasi).
add_path(f"{B}/components", {
    "get": op("listKpiComponents", "List performance components",
              "Ambil daftar komponen scoring KPI dengan pagination.",
              TAG, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{B}/components/{{id}}", {
    "get": op("getKpiComponentByID", "Get performance component by ID",
              "Ambil detail satu komponen scoring KPI.",
              TAG, parameters=[param("id")], responses=responses_ok("PerformanceComponentResponse")),
    "put": op("updateKpiComponent", "Update performance component",
              "Perbarui kode, nama, deskripsi, urutan, atau status aktif komponen scoring.",
              TAG, parameters=[param("id")], request_body=content("UpdatePerformanceComponentRequest"),
              responses=responses_ok("PerformanceComponentResponse")),
})
# Hapus phantom method POST/DELETE yang sempat terdokumentasi (tidak terdaftar
# di routes.go) — script tetap idempoten.
for _p, _m in ((f"{B}/components", "post"), (f"{B}/components/{{id}}", "delete")):
    if _p in paths and _m in paths[_p]:
        del paths[_p][_m]
        print(f"  removed phantom {_m.upper()} {_p}")

# --- Organization component weight configuration ---
add_path(f"{B}/organization-components", {
    "post": op("upsertOrganizationComponent", "Upsert organization component weight",
               "Atur/upsert konfigurasi bobot komponen scoring untuk sebuah organisasi (enable/disable, weight, sort order).",
               TAG, request_body=content("UpsertOrganizationComponentRequest"),
               responses=responses_created("OrganizationComponentResponse")),
})
add_path(f"{B}/organizations/{{organization_id}}/components", {
    "get": op("listOrganizationComponents", "List organization components",
              "Ambil daftar komponen scoring yang diaktifkan untuk sebuah organisasi beserta bobotnya.",
              TAG, parameters=[param("organization_id")],
              responses=responses_array("OrganizationComponentResponse", "List of organization components")),
})
add_path(f"{B}/organization-components/{{id}}", {
    "delete": op("deleteOrganizationComponent", "Delete organization component",
                 "Hapus konfigurasi komponen scoring dari organisasi.",
                 TAG, parameters=[param("id")], responses=responses_plain("Organization component deleted")),
})

# --- Evaluation scoring engine ---
add_path(f"{B}/evaluations/{{id}}/components", {
    "get": op("listEvaluationComponents", "List evaluation component scores",
              "Ambil skor per komponen untuk sebuah evaluasi KPI (hasil scoring engine).",
              TAG, parameters=[param("id")],
              responses=responses_array("PerformanceEvaluationComponentResponse", "List of evaluation component scores")),
})
add_path(f"{B}/evaluations/{{id}}/calculate-scoring", {
    "post": op("calculateEvaluationComponentScoring", "Calculate evaluation component scoring",
               "Jalankan scoring engine: hitung skor tiap komponen dari data terkait (KPI, competency, dll) lalu simpan hasilnya.",
               TAG, parameters=[param("id")],
               responses=responses_array("PerformanceEvaluationComponentResponse", "Calculated component scores")),
})
add_path(f"{B}/evaluations/{{id}}/components/{{component_id}}", {
    "put": op("updateEvaluationComponentScore", "Update evaluation component score",
              "Isi skor komponen secara manual (mis. Work Program yang tidak bisa dihitung otomatis — wajib diisi reviewer).",
              TAG, parameters=[param("id"), param("component_id")],
              request_body=content("UpdateEvaluationComponentScoreRequest"),
              responses=responses_ok("PerformanceEvaluationComponentResponse")),
})
add_path(f"{B}/periods/{{period_id}}/recalculate-scoring", {
    "post": op("recalculatePeriodScoring", "Recalculate period scoring (batch)",
               "Jalankan batch recalculation scoring untuk seluruh evaluasi dalam sebuah periode KPI (bottom-up subordinate scoring / akhir periode) lalu simpan hasilnya.",
               TAG, parameters=[param("period_id")],
               responses=responses_ok("RecalculatePeriodScoringResponse")),
})

# =========================================================================
# 4. OKR BARU — seluruh endpoint sub-modul OKR
# =========================================================================
O = "/api/v1/tenant/performance/okr"

# --- OKR Context ---
add_path(f"{O}/my-context", {
    "get": op("getMyOkrContext", "Get my OKR context",
              "Konteks OKR user saat ini: apakah punya posisi, employee/organization terkait, dan daftar template OKR yang tersedia (dipakai self-assessment).",
              TAG, responses=responses_ok("MyOKRContextResponse")),
})
add_path(f"{O}/templates/objective-scope", {
    "get": op("getOkrObjectiveScope", "Get OKR objective creation scope",
              "Resolusi scope pembuatan objective cascading: apakah employee saat ini boleh membuat objective untuk organisasi bawahan, dan organisasi mana yang memenuhi syarat sebagai subordinate (walking down melalui organisasi kosong).",
              TAG, responses=responses_ok("OKRObjectiveScopeResponse")),
})

# --- OKR Templates ---
add_path(f"{O}/templates", {
    "post": op("createOkrTemplate", "Create OKR template",
               "Buat template OKR untuk sebuah organisasi & periode (dengan objective & key results).",
               TAG, request_body=content("CreateOKRTemplateRequest"), responses=responses_created("OKRTemplateResponse")),
    "get": op("listOkrTemplates", "List OKR templates",
              "Ambil daftar template OKR dengan pagination, filter per organisasi/periode.",
              TAG, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{O}/templates/{{id}}", {
    "get": op("getOkrTemplateByID", "Get OKR template by ID",
              "Ambil detail template OKR termasuk daftar objective & key results.",
              TAG, parameters=[param("id")], responses=responses_ok("OKRTemplateResponse")),
    "put": op("updateOkrTemplate", "Update OKR template",
              "Perbarui nama, periode, status, atau tanggal efektif template OKR.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKRTemplateRequest"),
              responses=responses_ok("OKRTemplateResponse")),
    "delete": op("deleteOkrTemplate", "Delete OKR template",
                 "Hapus template OKR beserta objective & key results terkait.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR template deleted")),
})
add_path(f"{O}/templates/{{id}}/duplicate", {
    "post": op("duplicateOkrTemplate", "Duplicate OKR template",
               "Duplikat template OKR (beserta objective & key results) sebagai template baru.",
               TAG, parameters=[param("id")], responses=responses_created("OKRTemplateResponse")),
})
add_path(f"{O}/templates/{{id}}/objectives", {
    "get": op("listOkrObjectivesByTemplateID", "List OKR objectives by template",
              "Ambil daftar objective dalam sebuah template OKR (non-paginated).",
              TAG, parameters=[param("id")], responses=responses_array("OKRObjectiveResponse", "List of OKR objectives")),
})

# --- OKR Objectives ---
add_path(f"{O}/objectives", {
    "post": op("createOkrObjective", "Create OKR objective",
               "Buat objective baru di dalam sebuah template OKR.",
               TAG, request_body=content("CreateOKRObjectiveRequest"), responses=responses_created("OKRObjectiveResponse")),
})
add_path(f"{O}/objectives/{{id}}", {
    "get": op("getOkrObjectiveByID", "Get OKR objective by ID",
              "Ambil detail objective OKR beserta key results.",
              TAG, parameters=[param("id")], responses=responses_ok("OKRObjectiveResponse")),
    "put": op("updateOkrObjective", "Update OKR objective",
              "Perbarui judul, deskripsi, bobot, atau urutan objective.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKRObjectiveRequest"),
              responses=responses_ok("OKRObjectiveResponse")),
    "delete": op("deleteOkrObjective", "Delete OKR objective",
                 "Hapus objective beserta key results terkait.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR objective deleted")),
})
add_path(f"{O}/objectives/{{id}}/key-results", {
    "get": op("listOkrKeyResultsByObjectiveID", "List OKR key results by objective",
              "Ambil daftar key results dalam sebuah objective (non-paginated).",
              TAG, parameters=[param("id")], responses=responses_array("OKRKeyResultResponse", "List of OKR key results")),
})

# --- OKR Key Results ---
add_path(f"{O}/key-results", {
    "post": op("createOkrKeyResult", "Create OKR key result",
               "Buat key result terukur di dalam sebuah objective (target, unit, formula, bobot).",
               TAG, request_body=content("CreateOKRKeyResultRequest"), responses=responses_created("OKRKeyResultResponse")),
})
add_path(f"{O}/key-results/{{id}}", {
    "get": op("getOkrKeyResultByID", "Get OKR key result by ID",
              "Ambil detail key result OKR.",
              TAG, parameters=[param("id")], responses=responses_ok("OKRKeyResultResponse")),
    "put": op("updateOkrKeyResult", "Update OKR key result",
              "Perbarui target, unit, formula, bobot, atau status key result.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKRKeyResultRequest"),
              responses=responses_ok("OKRKeyResultResponse")),
    "delete": op("deleteOkrKeyResult", "Delete OKR key result",
                 "Hapus key result dari objective.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR key result deleted")),
})

# --- OKR Evaluations ---
add_path(f"{O}/evaluations", {
    "post": op("createOkrEvaluation", "Create OKR evaluation with snapshot",
               "Buat evaluasi OKR employee dari template (objective & key results di-snapshot ke evaluation details).",
               TAG, request_body=content("CreateOKREvaluationRequest"), responses=responses_created("OKREvaluationResponse")),
    "get": op("listOkrEvaluations", "List OKR evaluations",
              "Ambil daftar evaluasi OKR dengan pagination, filter per employee/organisasi/periode/status.",
              TAG, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{O}/evaluations/{{id}}", {
    "get": op("getOkrEvaluationByID", "Get OKR evaluation by ID",
              "Ambil detail evaluasi OKR termasuk details, skor, dan status workflow.",
              TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
    "put": op("updateOkrEvaluation", "Update OKR evaluation",
              "Perbarui status atau catatan reviewer evaluasi OKR.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKREvaluationRequest"),
              responses=responses_ok("OKREvaluationResponse")),
    "delete": op("deleteOkrEvaluation", "Delete OKR evaluation",
                 "Hapus evaluasi OKR (hanya status DRAFT).",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR evaluation deleted")),
})
add_path(f"{O}/evaluations/{{id}}/details", {
    "get": op("getOkrEvaluationWithDetails", "Get OKR evaluation with details",
              "Ambil evaluasi OKR lengkap dengan seluruh evaluation details.",
              TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/actuals", {
    "put": op("bulkUpdateOkrEvaluationActuals", "Bulk update OKR evaluation actuals",
              "Input nilai aktual sekaligus untuk banyak evaluation detail OKR.",
              TAG, parameters=[param("id")], request_body=content("OKRBulkUpdateActualsRequest"),
              responses=responses_plain("OKR evaluation actuals updated")),
})
add_path(f"{O}/evaluations/{{id}}/recalculate", {
    "post": op("recalculateOkrEvaluationScore", "Recalculate OKR evaluation score",
               "Hitung ulang skor & achievement evaluasi OKR dari nilai aktual dan formula tiap key result.",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/key-results", {
    "post": op("createOkrEvaluationKeyResult", "Create employee-proposed OKR evaluation key result",
               "Tambahkan Key Result yang diajukan employee di bawah Objective hasil snapshot, hanya saat evaluasi berstatus DRAFT (fase KR proposal).",
               TAG, parameters=[param("id")], request_body=content("CreateOKREvaluationKeyResultRequest"),
               responses=responses_created("OKREvaluationDetailResponse")),
})
add_path(f"{O}/evaluations/{{id}}/submit-key-results", {
    "post": op("submitOkrKeyResults", "Submit OKR key results for approval",
               "Ajukan proposal Key Results karyawan untuk persetujuan (status → KR_SUBMITTED) melalui approval flow modul performance.",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/approve-key-results", {
    "post": op("approveOkrKeyResults", "Approve OKR key results",
               "Setujui proposal Key Results (status → KR_APPROVED, \"OKR Active\") — dapat diresolusi otomatis oleh approval flow jika dikonfigurasi.",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/reject-key-results", {
    "post": op("rejectOkrKeyResults", "Reject OKR key results",
               "Tolak proposal Key Results kembali ke DRAFT (status → DRAFT) dengan catatan penolakan, agar karyawan dapat memperbaiki proposal.",
               TAG, parameters=[param("id")], request_body=content("RejectOkrKeyResultsRequest"),
               responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/submit", {
    "post": op("submitOkrEvaluation", "Submit OKR evaluation",
               "Ajukan evaluasi OKR (status → SUBMITTED).",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/approve", {
    "post": op("approveOkrEvaluation", "Approve OKR evaluation",
               "Setujui evaluasi OKR (status → APPROVED).",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/reject", {
    "post": op("rejectOkrEvaluation", "Reject OKR evaluation",
               "Tolak evaluasi OKR (status → REJECTED).",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/complete", {
    "post": op("completeOkrEvaluation", "Complete OKR evaluation",
               "Selesaikan evaluasi OKR (status → COMPLETED).",
               TAG, parameters=[param("id")], responses=responses_ok("OKREvaluationResponse")),
})
add_path(f"{O}/evaluations/{{id}}/comments", {
    "get": op("listOkrCommentsByEvaluationID", "List OKR comments by evaluation",
              "Ambil semua komentar (dengan replies) pada sebuah evaluasi OKR.",
              TAG, parameters=[param("id")], responses=responses_array("OKRCommentResponse", "List of OKR comments")),
})

# --- OKR Evaluation Details ---
add_path(f"{O}/evaluation-details/{{id}}", {
    "put": op("updateOkrEvaluationDetailActual", "Update OKR evaluation detail actual",
              "Input nilai aktual satu evaluation detail OKR.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKREvaluationDetailRequest"),
              responses=responses_ok("OKREvaluationDetailResponse")),
})
add_path(f"{O}/evaluation-details/{{id}}/progress", {
    "get": op("listOkrProgressByDetailID", "List OKR progress by evaluation detail",
              "Ambil riwayat progres (check-in) untuk satu evaluation detail OKR.",
              TAG, parameters=[param("id")], responses=responses_array("OKRProgressResponse", "List of OKR progress entries")),
})
add_path(f"{O}/evaluation-details/{{id}}/attachments", {
    "get": op("listOkrAttachmentsByDetailID", "List OKR attachments by evaluation detail",
              "Ambil daftar lampiran bukti pada satu evaluation detail OKR.",
              TAG, parameters=[param("id")], responses=responses_array("OKRAttachmentResponse", "List of OKR attachments")),
})

# --- OKR Evaluation Key Results (employee-proposed, DRAFT phase) ---
add_path(f"{O}/evaluation-key-results/{{id}}/target", {
    "put": op("updateOkrEvaluationKeyResultTarget", "Update OKR evaluation key result target",
              "Perbarui target Key Result yang diajukan karyawan (judul, target, unit, formula, bobot) sebelum disubmit.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKREvaluationKeyResultTargetRequest"),
              responses=responses_ok("OKREvaluationDetailResponse")),
})
add_path(f"{O}/evaluation-key-results/{{id}}", {
    "delete": op("deleteOkrEvaluationKeyResult", "Delete OKR evaluation key result",
                 "Hapus Key Result yang diajukan karyawan, hanya saat evaluasi masih DRAFT.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR evaluation key result deleted")),
})

# --- OKR Progress ---
add_path(f"{O}/progress", {
    "post": op("createOkrProgress", "Create OKR progress",
               "Catat progres (check-in) nilai aktual untuk satu evaluation detail OKR.",
               TAG, request_body=content("CreateOKRProgressRequest"), responses=responses_created("OKRProgressResponse")),
})
add_path(f"{O}/progress/{{id}}", {
    "get": op("getOkrProgressByID", "Get OKR progress by ID",
              "Ambil satu catatan progres OKR.",
              TAG, parameters=[param("id")], responses=responses_ok("OKRProgressResponse")),
    "put": op("updateOkrProgress", "Update OKR progress",
              "Perbarui tanggal, nilai aktual, atau catatan progres OKR.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKRProgressRequest"),
              responses=responses_ok("OKRProgressResponse")),
    "delete": op("deleteOkrProgress", "Delete OKR progress",
                 "Hapus satu catatan progres OKR.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR progress deleted")),
})

# --- OKR Comments ---
add_path(f"{O}/comments", {
    "post": op("createOkrComment", "Create OKR comment",
               "Tambahkan komentar/review pada evaluasi OKR (mendukung reply via parent_id).",
               TAG, request_body=content("CreateOKRCommentRequest"), responses=responses_created("OKRCommentResponse")),
})
add_path(f"{O}/comments/{{id}}", {
    "put": op("updateOkrComment", "Update OKR comment",
              "Perbarui isi komentar OKR.",
              TAG, parameters=[param("id")], request_body=content("UpdateOKRCommentRequest"),
              responses=responses_ok("OKRCommentResponse")),
    "delete": op("deleteOkrComment", "Delete OKR comment",
                 "Hapus komentar OKR.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR comment deleted")),
})

# --- OKR Attachments ---
add_path(f"{O}/attachments", {
    "post": op("createOkrAttachment", "Create OKR attachment",
               "Lampirkan file bukti/dokumen pendukung pada evaluation detail OKR.",
               TAG, request_body=content("CreateOKRAttachmentRequest"), responses=responses_created("OKRAttachmentResponse")),
})
add_path(f"{O}/attachments/{{id}}", {
    "delete": op("deleteOkrAttachment", "Delete OKR attachment",
                 "Hapus satu lampiran OKR.",
                 TAG, parameters=[param("id")], responses=responses_plain("OKR attachment deleted")),
})

# --- OKR Dashboard ---
add_path(f"{O}/dashboard/hr", {
    "get": op("getOkrHRDashboard", "Get OKR HR dashboard",
              "Dashboard OKR untuk HR: total evaluasi, sebaran status, skor & achievement rata-rata, distribusi rating.",
              TAG, responses=responses_plain("OKR HR dashboard data")),
})

# =========================================================================
# 5. Schemas (KPI Components + OKR)
# =========================================================================
new_schemas = {
    # --- Performance Component (Phase 5) ---
    # CATATAN: CreatePerformanceComponentRequest sengaja TIDAK dibuat — sejak
    # KPI Phase 1 (lock components) tidak ada endpoint create komponen.
    "UpdatePerformanceComponentRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 50},
            "name": {"type": "string", "maxLength": 100},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
            "is_active": {"type": "boolean"},
        },
    },
    "PerformanceComponentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "name": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
            "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "UpsertOrganizationComponentRequest": {
        "type": "object", "required": ["organization_id", "component_id"],
        "properties": {
            "organization_id": {"type": "string", "format": "uuid"},
            "component_id": {"type": "string", "format": "uuid"},
            "weight": {"type": "number", "example": 30},
            "is_enabled": {"type": "boolean", "example": True},
            "sort_order": {"type": "integer"},
        },
    },
    "OrganizationComponentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"},
            "component_id": {"type": "string", "format": "uuid"},
            "component_code": {"type": "string"},
            "component_name": {"type": "string"},
            "weight": {"type": "number"},
            "is_enabled": {"type": "boolean"},
            "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "PerformanceEvaluationComponentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_id": {"type": "string", "format": "uuid"},
            "component_id": {"type": "string", "format": "uuid"},
            "component_name": {"type": "string"},
            "score": {"type": "number"},
            "weight": {"type": "number"},
            "final_score": {"type": "number"},
            "calculated_at": {"type": "string", "format": "date-time"},
        },
    },
    "UpdateEvaluationComponentScoreRequest": {
        "type": "object", "required": ["score"],
        "properties": {
            "score": {"type": "number", "minimum": 0, "maximum": 100},
        },
    },
    "OrganizationOptionResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
        },
    },
    "RecalculatePeriodScoringResponse": {
        "type": "object",
        "properties": {
            "period_id": {"type": "string", "format": "uuid"},
            "total": {"type": "integer", "description": "Total evaluasi di periode", "example": 42},
            "processed": {"type": "integer", "description": "Jumlah evaluasi berhasil dihitung ulang", "example": 41},
            "failed": {"type": "integer", "description": "Jumlah evaluasi gagal", "example": 1},
        },
    },
    "MyOKRContextResponse": {
        "type": "object",
        "properties": {
            "has_position": {"type": "boolean", "example": True},
            "employee_id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"},
            "organization_name": {"type": "string"},
            "templates": {"type": "array", "items": ref("OKRTemplateResponse")},
        },
    },
    "OKRObjectiveScopeResponse": {
        "type": "object",
        "properties": {
            "organization_id": {"type": "string", "format": "uuid"},
            "organization_name": {"type": "string"},
            "eligible": {"type": "boolean", "description": "Apakah employee saat ini boleh membuat objective untuk organisasi bawahan", "example": True},
            "ineligible_reason_key": {"type": "string", "description": "Kunci pesan alasan tidak boleh membuat objective (saat eligible=false)", "example": "okr.no_employee_position"},
            "subordinate_organizations": {"type": "array", "items": ref("OrganizationOptionResponse")},
        },
    },
    # --- OKR Template ---
    "CreateOKRTemplateRequest": {
        "type": "object", "required": ["organization_id", "name"],
        "properties": {
            "organization_id": {"type": "string", "format": "uuid"},
            "period_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "status": {"type": "integer", "example": 0},
            "effective_date": {"type": "string", "format": "date"},
            "expired_date": {"type": "string", "format": "date"},
        },
    },
    "UpdateOKRTemplateRequest": {
        "type": "object",
        "properties": {
            "period_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "status": {"type": "integer"},
            "effective_date": {"type": "string", "format": "date"},
            "expired_date": {"type": "string", "format": "date"},
        },
    },
    "OKRTemplateResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"},
            "organization_name": {"type": "string"},
            "period_id": {"type": "string", "format": "uuid"},
            "period_code": {"type": "string"},
            "name": {"type": "string"},
            "description": {"type": "string"},
            "status": {"type": "integer"},
            "effective_date": {"type": "string", "format": "date"},
            "expired_date": {"type": "string", "format": "date"},
            "objective_count": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
            "objectives": {"type": "array", "items": ref("OKRObjectiveResponse")},
        },
    },
    # --- OKR Objective ---
    "CreateOKRObjectiveRequest": {
        "type": "object", "required": ["template_id", "title"],
        "properties": {
            "template_id": {"type": "string", "format": "uuid"},
            "code": {"type": "string", "maxLength": 50},
            "title": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "weight": {"type": "number", "example": 100},
            "sort_order": {"type": "integer"},
        },
    },
    "UpdateOKRObjectiveRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 50},
            "title": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "weight": {"type": "number"},
            "sort_order": {"type": "integer"},
        },
    },
    "OKRObjectiveResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "template_id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "title": {"type": "string"},
            "description": {"type": "string"},
            "weight": {"type": "number"},
            "sort_order": {"type": "integer"},
            "key_results": {"type": "array", "items": ref("OKRKeyResultResponse")},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Key Result ---
    "CreateOKRKeyResultRequest": {
        "type": "object", "required": ["objective_id", "title"],
        "properties": {
            "objective_id": {"type": "string", "format": "uuid"},
            "code": {"type": "string", "maxLength": 50},
            "title": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "target_type": {"type": "string", "enum": ["NUMBER", "CURRENCY", "PERCENTAGE", "DURATION", "BOOLEAN"]},
            "target_value": {"type": "number"},
            "unit": {"type": "string", "maxLength": 50},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE", "BOOLEAN", "PERCENTAGE"]},
            "weight": {"type": "number"},
            "minimum_score": {"type": "number"},
            "maximum_score": {"type": "number"},
            "sort_order": {"type": "integer"},
            "is_required": {"type": "boolean"},
        },
    },
    "UpdateOKRKeyResultRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 50},
            "title": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "target_type": {"type": "string", "enum": ["NUMBER", "CURRENCY", "PERCENTAGE", "DURATION", "BOOLEAN"]},
            "target_value": {"type": "number"},
            "unit": {"type": "string", "maxLength": 50},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE", "BOOLEAN", "PERCENTAGE"]},
            "weight": {"type": "number"},
            "minimum_score": {"type": "number"},
            "maximum_score": {"type": "number"},
            "sort_order": {"type": "integer"},
            "is_required": {"type": "boolean"},
        },
    },
    "OKRKeyResultResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "objective_id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "title": {"type": "string"},
            "description": {"type": "string"},
            "target_type": {"type": "string"},
            "target_value": {"type": "number"},
            "unit": {"type": "string"},
            "formula_type": {"type": "string"},
            "weight": {"type": "number"},
            "minimum_score": {"type": "number"},
            "maximum_score": {"type": "number"},
            "sort_order": {"type": "integer"},
            "is_required": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Evaluation ---
    "CreateOKREvaluationRequest": {
        "type": "object", "required": ["employee_id", "organization_id", "period_id", "template_id"],
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"},
            "period_id": {"type": "string", "format": "uuid"},
            "template_id": {"type": "string", "format": "uuid"},
        },
    },
    "UpdateOKREvaluationRequest": {
        "type": "object",
        "properties": {
            "status": {"type": "string", "enum": ["DRAFT", "SUBMITTED", "APPROVED", "COMPLETED", "REJECTED"]},
            "reviewer_notes": {"type": "string"},
        },
    },
    "OKREvaluationResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "employee_name": {"type": "string"},
            "organization_id": {"type": "string", "format": "uuid"},
            "organization_name": {"type": "string"},
            "period_id": {"type": "string", "format": "uuid"},
            "period_code": {"type": "string"},
            "template_id": {"type": "string", "format": "uuid"},
            "template_name": {"type": "string"},
            "status": {"type": "string", "description": "Status workflow dua fase: DRAFT → KR_SUBMITTED → KR_APPROVED → SUBMITTED → COMPLETED (atau REJECTED)", "example": "KR_SUBMITTED"},
            "submitted_at": {"type": "string", "format": "date-time"},
            "approved_at": {"type": "string", "format": "date-time"},
            "kr_approval_instance_id": {"type": "string", "format": "uuid", "description": "ID approval instance untuk checkpoint persetujuan Key Results"},
            "assessment_approval_instance_id": {"type": "string", "format": "uuid", "description": "ID approval instance untuk checkpoint persetujuan penilaian akhir"},
            "final_score": {"type": "number"},
            "rating_id": {"type": "string", "format": "uuid"},
            "rating_name": {"type": "string"},
            "rating_color": {"type": "string"},
            "reviewer_notes": {"type": "string"},
            "details": {"type": "array", "items": ref("OKREvaluationDetailResponse")},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Evaluation Detail ---
    "UpdateOKREvaluationDetailRequest": {
        "type": "object", "required": ["actual_value"],
        "properties": {
            "actual_value": {"type": "number"},
            "remarks": {"type": "string"},
        },
    },
    "OKRBulkUpdateActualsRequest": {
        "type": "object", "required": ["details"],
        "properties": {
            "details": {"type": "array", "items": {
                "type": "object", "required": ["id"],
                "properties": {"id": {"type": "string", "format": "uuid"}, "actual_value": {"type": "number"}},
            }},
        },
    },
    "CreateOKREvaluationKeyResultRequest": {
        "type": "object",
        "required": ["evaluation_id", "objective_id", "objective_title", "title"],
        "properties": {
            "evaluation_id": {"type": "string", "format": "uuid"},
            "objective_id": {"type": "string", "format": "uuid"},
            "objective_title": {"type": "string", "maxLength": 255, "example": "Meningkatkan Kepuasan Pelanggan"},
            "objective_weight": {"type": "number", "example": 40},
            "title": {"type": "string", "maxLength": 255, "example": "Skor CSAT ≥ 85"},
            "target_type": {"type": "string", "enum": ["NUMBER", "PERCENTAGE", "CURRENCY", "BOOLEAN", "DURATION", "SCORE"], "example": "PERCENTAGE"},
            "target_value": {"type": "number", "example": 85},
            "unit": {"type": "string", "maxLength": 50, "example": "%"},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE", "BOOLEAN", "PERCENTAGE"], "example": "HIGHER_BETTER"},
            "weight": {"type": "number", "example": 100},
        },
    },
    "UpdateOKREvaluationKeyResultTargetRequest": {
        "type": "object",
        "properties": {
            "title": {"type": "string", "maxLength": 255},
            "target_type": {"type": "string", "enum": ["NUMBER", "PERCENTAGE", "CURRENCY", "BOOLEAN", "DURATION", "SCORE"]},
            "target_value": {"type": "number"},
            "unit": {"type": "string", "maxLength": 50},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE", "BOOLEAN", "PERCENTAGE"]},
            "weight": {"type": "number"},
        },
    },
    "RejectOkrKeyResultsRequest": {
        "type": "object",
        "properties": {
            "notes": {"type": "string", "description": "Catatan alasan penolakan proposal Key Results", "example": "Target perlu direvisi agar lebih realistis"},
        },
    },
    "OKREvaluationDetailResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_id": {"type": "string", "format": "uuid"},
            "objective_id": {"type": "string", "format": "uuid"},
            "key_result_id": {"type": "string", "format": "uuid"},
            "objective_title": {"type": "string"},
            "key_result_title": {"type": "string"},
            "objective_weight": {"type": "number"},
            "key_result_weight": {"type": "number"},
            "target_value": {"type": "number"},
            "target_type": {"type": "string"},
            "unit": {"type": "string"},
            "formula_type": {"type": "string"},
            "actual_value": {"type": "number"},
            "achievement": {"type": "number"},
            "score": {"type": "number"},
            "remarks": {"type": "string"},
            "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Progress ---
    "CreateOKRProgressRequest": {
        "type": "object", "required": ["evaluation_detail_id", "progress_date"],
        "properties": {
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "progress_date": {"type": "string", "format": "date"},
            "actual_value": {"type": "number"},
            "notes": {"type": "string"},
        },
    },
    "UpdateOKRProgressRequest": {
        "type": "object",
        "properties": {
            "progress_date": {"type": "string", "format": "date"},
            "actual_value": {"type": "number"},
            "notes": {"type": "string"},
        },
    },
    "OKRProgressResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "progress_date": {"type": "string", "format": "date"},
            "actual_value": {"type": "number"},
            "achievement": {"type": "number"},
            "notes": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Comment ---
    "CreateOKRCommentRequest": {
        "type": "object", "required": ["evaluation_id", "comment"],
        "properties": {
            "evaluation_id": {"type": "string", "format": "uuid"},
            "parent_id": {"type": "string", "format": "uuid"},
            "comment": {"type": "string"},
        },
    },
    "UpdateOKRCommentRequest": {
        "type": "object", "required": ["comment"],
        "properties": {"comment": {"type": "string"}},
    },
    "OKRCommentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_id": {"type": "string", "format": "uuid"},
            "parent_id": {"type": "string", "format": "uuid"},
            "comment": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
            "created_by_name": {"type": "string"},
            "replies": {"type": "array", "items": ref("OKRCommentResponse")},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- OKR Attachment ---
    "CreateOKRAttachmentRequest": {
        "type": "object", "required": ["evaluation_detail_id", "file_path", "file_name"],
        "properties": {
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "file_path": {"type": "string", "example": "/uploads/evidence/xxx.pdf"},
            "file_name": {"type": "string", "example": "bukti-okr.pdf"},
            "file_type": {"type": "string"},
            "file_size": {"type": "integer", "format": "int64"},
            "description": {"type": "string"},
        },
    },
    "OKRAttachmentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "file_path": {"type": "string"},
            "file_name": {"type": "string"},
            "file_type": {"type": "string"},
            "file_size": {"type": "integer", "format": "int64"},
            "description": {"type": "string"},
            "uploaded_by": {"type": "string", "format": "uuid"},
            "uploaded_by_name": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
}

for name, schema in new_schemas.items():
    # Sinkronkan selalu: overwrite jika definisi berubah, tambah jika belum ada.
    if name not in schemas or schemas[name] != schema:
        schemas[name] = schema
        added += 1

with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

total_paths = len(spec["paths"])
total_schemas = len(spec["components"]["schemas"])
print(f"[OK] Injected/updated {added} items into {JSON_PATH}")
print(f"     Total paths:   {total_paths}")
print(f"     Total schemas: {total_schemas}")
