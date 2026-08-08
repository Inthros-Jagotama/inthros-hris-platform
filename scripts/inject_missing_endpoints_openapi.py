#!/usr/bin/env python3
"""Inject missing endpoints into openapi.json (idempotent).

Finds endpoints registered in backend routes (performance Phase 2 KPI,
tenant RBAC, settings competencies, user accounts, platform users/licenses,
career intelligence paths delete, public setup-password,approval available-modules & active-flow, employee movement submit) that are
missing from the OpenAPI spec and documents them with request/response schemas.

Usage:
    python scripts/inject_missing_endpoints_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]


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
    """Plain non-paginated list: {success, data: [item]}."""
    return {
        "200": {
            "description": desc,
            "content": {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}, "data": {"type": "array", "items": ref(item_ref)}}}}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_string_array(desc="List of strings"):
    """Plain non-paginated list of strings: {success, data: ["slug", ...]}."""
    return {
        "200": {
            "description": desc,
            "content": {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}, "data": {"type": "array", "items": {"type": "string", "example": "leave"}}}}}}},
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


with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

paths = spec["paths"]
schemas = spec["components"]["schemas"]
added = 0
TAG_PERF = "Tenant: Performance Management"

# =========================================================================
# Performance Phase 2 — Progress
# (KPI sub-module: routes terdaftar di /performance/kpi/*)
# =========================================================================
B = "/api/v1/tenant/performance"
K = f"{B}/kpi"

def add_path(path, ops):
    global added
    if path not in paths:
        paths[path] = {}
    for method, o in ops.items():
        if method not in paths[path]:
            paths[path][method] = o
            added += 1


# --- Progress ---
add_path(f"{K}/progress", {
    "post": op("createPerformanceProgress", "Create performance progress",
               "Catat progres realisasi KPI untuk satu evaluation detail (nilai aktual per tanggal).",
               TAG_PERF, request_body=content("CreatePerformanceProgressRequest"), responses=responses_created("PerformanceProgressResponse")),
})
add_path(f"{K}/evaluation-details/{{id}}/progress", {
    "get": op("listPerformanceProgressByDetailID", "List progress by evaluation detail",
              "Ambil semua progres KPI yang tercatat untuk evaluation detail tertentu (non-paginated).",
              TAG_PERF, parameters=[param("id")], responses=responses_array("PerformanceProgressResponse", "List of performance progress entries")),
})
add_path(f"{K}/progress/{{id}}", {
    "get": op("getPerformanceProgress", "Get performance progress by ID",
              "Ambil satu catatan progres KPI berdasarkan ID.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceProgressResponse")),
    "put": op("updatePerformanceProgress", "Update performance progress",
              "Perbarui tanggal, nilai aktual, achievement, atau catatan progres KPI.",
              TAG_PERF, parameters=[param("id")], request_body=content("UpdatePerformanceProgressRequest"), responses=responses_ok("PerformanceProgressResponse")),
    "delete": op("deletePerformanceProgress", "Delete performance progress",
                 "Hapus satu catatan progres KPI.",
                 TAG_PERF, parameters=[param("id")], responses=responses_plain("Progress deleted")),
})

# --- Comments ---
add_path(f"{K}/comments", {
    "post": op("createPerformanceComment", "Create performance comment",
               "Tambahkan komentar/review pada sebuah performance evaluation.",
               TAG_PERF, request_body=content("CreatePerformanceCommentRequest"), responses=responses_created("PerformanceCommentResponse")),
})
add_path(f"{K}/evaluations/{{id}}/comments", {
    "get": op("listPerformanceCommentsByEvaluationID", "List comments by evaluation",
              "Ambil semua komentar pada sebuah performance evaluation (non-paginated).",
              TAG_PERF, parameters=[param("id")], responses=responses_array("PerformanceCommentResponse", "List of performance comments")),
})
add_path(f"{K}/comments/{{id}}", {
    "get": op("getPerformanceComment", "Get performance comment by ID",
              "Ambil satu komentar performance berdasarkan ID.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceCommentResponse")),
    "put": op("updatePerformanceComment", "Update performance comment",
              "Perbarui isi komentar performance.",
              TAG_PERF, parameters=[param("id")], request_body=content("UpdatePerformanceCommentRequest"), responses=responses_ok("PerformanceCommentResponse")),
    "delete": op("deletePerformanceComment", "Delete performance comment",
                 "Hapus satu komentar performance.",
                 TAG_PERF, parameters=[param("id")], responses=responses_plain("Comment deleted")),
})

# --- Attachments ---
add_path(f"{K}/attachments", {
    "post": op("createPerformanceAttachment", "Create performance attachment",
               "Lampirkan file bukti/dokumen pendukung pada evaluation detail.",
               TAG_PERF, request_body=content("CreatePerformanceAttachmentRequest"), responses=responses_created("PerformanceAttachmentResponse")),
})
add_path(f"{K}/evaluation-details/{{id}}/attachments", {
    "get": op("listPerformanceAttachmentsByDetailID", "List attachments by evaluation detail",
              "Ambil semua lampiran pada evaluation detail tertentu (non-paginated).",
              TAG_PERF, parameters=[param("id")], responses=responses_array("PerformanceAttachmentResponse", "List of performance attachments")),
})
add_path(f"{K}/attachments/{{id}}", {
    "get": op("getPerformanceAttachment", "Get performance attachment by ID",
              "Ambil satu lampiran performance berdasarkan ID.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceAttachmentResponse")),
    "put": op("updatePerformanceAttachment", "Update performance attachment",
              "Perbarui deskripsi lampiran performance.",
              TAG_PERF, parameters=[param("id")], request_body=content("UpdatePerformanceAttachmentRequest"), responses=responses_ok("PerformanceAttachmentResponse")),
    "delete": op("deletePerformanceAttachment", "Delete performance attachment",
                 "Hapus satu lampiran performance.",
                 TAG_PERF, parameters=[param("id")], responses=responses_plain("Attachment deleted")),
})

# --- Ratings (master data) ---
add_path(f"{B}/ratings", {
    "post": op("createPerformanceRating", "Create performance rating",
               "Buat skala rating penilaian (mis. A=90-100, B=80-89) untuk konversi skor akhir.",
               TAG_PERF, request_body=content("CreatePerformanceRatingRequest"), responses=responses_created("PerformanceRatingResponse")),
    "get": op("listPerformanceRatings", "List performance ratings",
              "Ambil daftar rating skala penilaian dengan pagination.",
              TAG_PERF, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{B}/ratings/{{id}}", {
    "get": op("getPerformanceRating", "Get performance rating by ID",
              "Ambil satu rating skala penilaian.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceRatingResponse")),
    "put": op("updatePerformanceRating", "Update performance rating",
              "Perbarui rating skala penilaian.",
              TAG_PERF, parameters=[param("id")], request_body=content("UpdatePerformanceRatingRequest"), responses=responses_ok("PerformanceRatingResponse")),
    "delete": op("deletePerformanceRating", "Delete performance rating",
                 "Hapus satu rating skala penilaian.",
                 TAG_PERF, parameters=[param("id")], responses=responses_plain("Rating deleted")),
})

# --- Indicator Formulas (master data) ---
add_path(f"{B}/indicator-formulas", {
    "post": op("createPerformanceIndicatorFormula", "Create indicator formula",
               "Buat formula kalkulasi skor KPI (MANUAL/HIGHER_BETTER/LOWER_BETTER/RANGE).",
               TAG_PERF, request_body=content("CreatePerformanceIndicatorFormulaRequest"), responses=responses_created("PerformanceIndicatorFormulaResponse")),
    "get": op("listPerformanceIndicatorFormulas", "List indicator formulas",
              "Ambil daftar formula kalkulasi skor KPI dengan pagination.",
              TAG_PERF, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{B}/indicator-formulas/{{id}}", {
    "get": op("getPerformanceIndicatorFormula", "Get indicator formula by ID",
              "Ambil satu formula kalkulasi skor KPI.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceIndicatorFormulaResponse")),
    "put": op("updatePerformanceIndicatorFormula", "Update indicator formula",
              "Perbarui formula kalkulasi skor KPI.",
              TAG_PERF, parameters=[param("id")], request_body=content("UpdatePerformanceIndicatorFormulaRequest"), responses=responses_ok("PerformanceIndicatorFormulaResponse")),
    "delete": op("deletePerformanceIndicatorFormula", "Delete indicator formula",
                 "Hapus satu formula kalkulasi skor KPI.",
                 TAG_PERF, parameters=[param("id")], responses=responses_plain("Formula deleted")),
})

# --- Logs (read-only) ---
add_path(f"{B}/logs", {
    "get": op("listPerformanceLogs", "List performance audit logs",
              "Ambil daftar audit trail perubahan data performance dengan pagination.",
              TAG_PERF, parameters=[qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{K}/evaluations/{{id}}/logs", {
    "get": op("listPerformanceLogsByEvaluationID", "List audit logs by evaluation",
              "Ambil audit trail perubahan pada sebuah performance evaluation.",
              TAG_PERF, parameters=[param("id"), qparam("page", {"type": "integer", "example": 1}), qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list("PerformancePaginatedResponse")),
})
add_path(f"{B}/logs/{{id}}", {
    "get": op("getPerformanceLog", "Get performance log by ID",
              "Ambil satu audit trail berdasarkan ID.",
              TAG_PERF, parameters=[param("id")], responses=responses_ok("PerformanceLogResponse")),
})

# =========================================================================
# Tenant RBAC
# =========================================================================
TAG_RBAC = "Tenant: RBAC Management"
R = "/api/v1/tenant/rbac"

add_path(f"{R}/roles", {
    "get": op("listTenantRbacRoles", "List tenant RBAC roles",
              "Ambil daftar role RBAC tenant beserta permission-nya (non-paginated).",
              TAG_RBAC, responses=responses_array("TenantRoleResponse", "List of tenant RBAC roles")),
    "post": op("createTenantRbacRole", "Create tenant RBAC role",
               "Buat role RBAC tenant baru.",
               TAG_RBAC, request_body=content("CreateTenantRoleRequest"), responses=responses_created("TenantRoleResponse")),
})
add_path(f"{R}/roles/{{id}}", {
    "put": op("updateTenantRbacRole", "Update tenant RBAC role",
              "Perbarui nama/deskripsi role RBAC tenant.",
              TAG_RBAC, parameters=[param("id")], request_body=content("UpdateTenantRoleRequest"), responses=responses_ok("TenantRoleResponse")),
    "delete": op("deleteTenantRbacRole", "Delete tenant RBAC role",
                 "Hapus role RBAC tenant (role system tidak dapat dihapus).",
                 TAG_RBAC, parameters=[param("id")], responses=responses_plain("Role deleted")),
})
add_path(f"{R}/permissions", {
    "get": op("listTenantRbacPermissions", "List tenant RBAC permissions",
              "Ambil daftar permission RBAC tenant (non-paginated).",
              TAG_RBAC, responses=responses_array("TenantPermissionResponse", "List of tenant RBAC permissions")),
})
add_path(f"{R}/roles/{{id}}/permissions", {
    "put": op("assignTenantRolePermissions", "Assign permissions to tenant role",
              "Ganti (replace) daftar permission milik role RBAC tenant.",
              TAG_RBAC, parameters=[param("id")], request_body=content("TenantAssignPermissionsRequest"), responses=responses_plain("Permissions updated")),
})
add_path(f"{R}/users", {
    "get": op("listTenantRbacUsers", "List tenant RBAC users",
              "Ambil daftar user tenant beserta role-nya (non-paginated).",
              TAG_RBAC, responses=responses_array("TenantUserResponse", "List of tenant RBAC users")),
})
add_path(f"{R}/users/{{id}}/roles", {
    "put": op("assignTenantUserRoles", "Assign roles to tenant user",
              "Ganti (replace) daftar role milik user tenant.",
              TAG_RBAC, parameters=[param("id")], request_body=content("TenantAssignUserRolesRequest"), responses=responses_plain("User roles updated")),
})

# =========================================================================
# Settings — Competencies
# =========================================================================
TAG_SET = "Tenant: Settings"
S = "/api/v1/tenant/settings/competencies"

add_path(S, {
    "post": op("createSettingCompetency", "Create setting competency",
               "Buat data kompetensi di master data settings.",
               TAG_SET, request_body=content("CreateCompetencyRequest"), responses=responses_created("CompetencyResponse")),
})
add_path(f"{S}/{{id}}", {
    "get": op("getSettingCompetency", "Get setting competency by ID",
              "Ambil satu data kompetensi di settings.",
              TAG_SET, parameters=[param("id")], responses=responses_ok("CompetencyResponse")),
    "put": op("updateSettingCompetency", "Update setting competency",
              "Perbarui data kompetensi di settings.",
              TAG_SET, parameters=[param("id")], request_body=content("UpdateCompetencyRequest"), responses=responses_ok("CompetencyResponse")),
    "delete": op("deleteSettingCompetency", "Delete setting competency",
                 "Hapus data kompetensi di settings.",
                 TAG_SET, parameters=[param("id")], responses=responses_plain("Competency deleted")),
})

# =========================================================================
# User Accounts (tenant)
# =========================================================================
TAG_UA = "Tenant: User Accounts"
U = "/api/v1/tenant/user-accounts"

add_path(f"{U}/employees/{{employeeId}}", {
    "get": op("getUserAccountStatus", "Get employee account status",
              "Ambil status akun login employee (email, password_set, setup token).",
              TAG_UA, parameters=[param("employeeId")], responses=responses_ok("AccountResponse")),
    "post": op("createUserAccount", "Create employee account",
               "Buat akun login employee (kirim email setup password).",
               TAG_UA, parameters=[param("employeeId")], request_body=content("CreateAccountRequest"), responses=responses_created("AccountResponse")),
})
add_path(f"{U}/employees/{{employeeId}}/resend", {
    "post": op("resendSetupEmail", "Resend account setup email",
               "Kirim ulang email link setup password akun employee.",
               TAG_UA, parameters=[param("employeeId")], responses=responses_plain("Setup email resent")),
})

# =========================================================================
# Platform — Users (GET/PUT by id)
# =========================================================================
TAG_PU = "Platform: Users"
PU = "/api/v1/platform/users/{id}"

add_path(PU, {
    "get": op("getPlatformUser", "Get platform user by ID",
              "Ambil detail platform user berdasarkan ID.",
              TAG_PU, parameters=[param("id")], responses=responses_ok("PlatformUserResponse")),
    "put": op("updatePlatformUser", "Update platform user",
              "Perbarui nama/email/role/status platform user.",
              TAG_PU, parameters=[param("id")], request_body=content("UpdateUserRequest"), responses=responses_ok("PlatformUserResponse")),
})

# =========================================================================
# Platform — Licenses (DELETE)
# =========================================================================
add_path("/api/v1/platform/licenses/{id}", {
    "delete": op("deleteLicense", "Delete license",
                 "Hapus lisensi perusahaan.",
                 "Platform: Licenses", parameters=[param("id")], responses=responses_plain("License deleted")),
})

# =========================================================================
# Career Intelligence — DELETE paths/{id}
# (existing spec wrongly documents DELETE on collection path — move it)
# =========================================================================
CI = "/api/v1/tenant/career-intelligence"
if "delete" in paths.get(CI + "/paths", {}):
    del paths[CI + "/paths"]["delete"]
    added -= 1  # will re-count when adding below
add_path(CI + "/paths/{id}", {
    "delete": op("deleteCareerPath", "Delete career path",
                 "Hapus satu jalur karier (career path).",
                 "Tenant: Career Intelligence", parameters=[param("id")], responses=responses_plain("Career path deleted")),
})

# =========================================================================
# Public — setup password
# =========================================================================
add_path("/api/v1/public/account/setup-password", {
    "post": op("publicSetPassword", "Set account password via email link",
               "Atur password akun login employee melalui link email (tanpa autentikasi).",
               "Public", security=False, request_body=content("SetPasswordRequest"), responses=responses_plain("Password set successfully")),
})

# =========================================================================
# Approval — available modules
# =========================================================================
TAG_APPROVAL = "Tenant: Approval"
add_path("/api/v1/tenant/approval/available-modules", {
    "get": op("listAvailableModules", "List available approval modules",
               "Ambil slug module yang aktif/disubscribe tenant — dipakai flow builder agar hanya menampilkan module yang tersedia.",
               TAG_APPROVAL, responses=responses_string_array("List of active module slugs")),
})

# =========================================================================
# Approval — active flow by module
# =========================================================================
add_path("/api/v1/tenant/approval/active-flow", {
    "get": op("getActiveFlowByModule", "Get active approval flow by module",
               "Resolusi otomatis alur persetujuan aktif untuk sebuah module (?module=xxx) — dipakai consumer yang ingin auto-resolve flow tanpa memilih flow_id manual (mis. KPI self-assessment dua fase). Kembalikan flow aktif versi tertinggi; bila tidak ada flow khusus, fallback ke flow module dasar (mis. performance untuk performance_kpi_target).",
               TAG_APPROVAL, parameters=[param("module", where="query", required=True, schema={"type": "string"})],
               responses=responses_ok("ApprovalFlowResponse")),
})

# =========================================================================
# Employee Movement — submit to approval
# =========================================================================
TAG_EMOV = "Tenant: Employee Movement & Career Management"
add_path("/api/v1/tenant/employee-movements/movements/{id}/submit", {
    "post": op("submitMovement", "Submit employee movement for approval",
               "Kirim movement berstatus draft ke alur persetujuan (approval flow) terpusat. Setelah disetujui, approval engine akan mengeksekusi perpindahan.",
               TAG_EMOV, parameters=[param("id")], request_body=content("SubmitMovementRequest"), responses=responses_ok("MovementResponse")),
})

# =========================================================================
# Schemas
# =========================================================================
new_schemas = {
    # --- Progress ---
    "CreatePerformanceProgressRequest": {
        "type": "object", "required": ["evaluation_detail_id", "progress_date", "created_by"],
        "properties": {
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "progress_date": {"type": "string", "format": "date", "example": "2026-08-01"},
            "actual_value": {"type": "number", "example": 12.5},
            "achievement": {"type": "number", "example": 96.15},
            "notes": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
        },
    },
    "UpdatePerformanceProgressRequest": {
        "type": "object",
        "properties": {
            "progress_date": {"type": "string", "format": "date"},
            "actual_value": {"type": "number"},
            "achievement": {"type": "number"},
            "notes": {"type": "string"},
        },
    },
    "PerformanceProgressResponse": {
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
    # --- Comments ---
    "CreatePerformanceCommentRequest": {
        "type": "object", "required": ["evaluation_id", "employee_id", "comment", "created_by"],
        "properties": {
            "evaluation_id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "comment": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
        },
    },
    "UpdatePerformanceCommentRequest": {
        "type": "object",
        "properties": {"comment": {"type": "string"}},
    },
    "PerformanceCommentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "comment": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- Attachments ---
    "CreatePerformanceAttachmentRequest": {
        "type": "object", "required": ["evaluation_detail_id", "file_path", "file_name", "uploaded_by"],
        "properties": {
            "evaluation_detail_id": {"type": "string", "format": "uuid"},
            "file_path": {"type": "string", "example": "/uploads/evidence/xxx.pdf"},
            "file_name": {"type": "string", "example": "bukti-pencapaian.pdf"},
            "file_type": {"type": "string", "example": "application/pdf"},
            "file_size": {"type": "integer", "format": "int64"},
            "description": {"type": "string"},
            "uploaded_by": {"type": "string", "format": "uuid"},
        },
    },
    "UpdatePerformanceAttachmentRequest": {
        "type": "object",
        "properties": {"description": {"type": "string"}},
    },
    "PerformanceAttachmentResponse": {
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
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- Ratings ---
    "CreatePerformanceRatingRequest": {
        "type": "object", "required": ["code", "name", "min_score", "max_score"],
        "properties": {
            "code": {"type": "string", "example": "A"},
            "name": {"type": "string", "example": "Excellent"},
            "min_score": {"type": "number", "example": 90},
            "max_score": {"type": "number", "example": 100},
            "color": {"type": "string", "example": "#22c55e"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer", "example": 1},
        },
    },
    "UpdatePerformanceRatingRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string"},
            "name": {"type": "string"},
            "min_score": {"type": "number"},
            "max_score": {"type": "number"},
            "color": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
        },
    },
    "PerformanceRatingResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "name": {"type": "string"},
            "min_score": {"type": "number"},
            "max_score": {"type": "number"},
            "color": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- Indicator Formulas ---
    "CreatePerformanceIndicatorFormulaRequest": {
        "type": "object", "required": ["code", "name", "formula_type"],
        "properties": {
            "code": {"type": "string", "example": "HB"},
            "name": {"type": "string", "example": "Higher is Better"},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE"]},
            "expression": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer", "example": 1},
        },
    },
    "UpdatePerformanceIndicatorFormulaRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string"},
            "name": {"type": "string"},
            "formula_type": {"type": "string", "enum": ["MANUAL", "HIGHER_BETTER", "LOWER_BETTER", "RANGE"]},
            "expression": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
        },
    },
    "PerformanceIndicatorFormulaResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "name": {"type": "string"},
            "formula_type": {"type": "string"},
            "expression": {"type": "string"},
            "description": {"type": "string"},
            "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- Logs ---
    "PerformanceLogResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "evaluation_id": {"type": "string", "format": "uuid"},
            "entity_type": {"type": "string", "example": "evaluation"},
            "entity_id": {"type": "string", "format": "uuid"},
            "action": {"type": "string", "example": "UPDATE"},
            "old_values": {"type": "string"},
            "new_values": {"type": "string"},
            "created_by": {"type": "string", "format": "uuid"},
            "created_at": {"type": "string", "format": "date-time"},
        },
    },
    # --- Tenant RBAC ---
    "CreateTenantRoleRequest": {
        "type": "object", "required": ["name"],
        "properties": {
            "name": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "is_default": {"type": "boolean", "default": False},
        },
    },
    "UpdateTenantRoleRequest": {
        "type": "object",
        "properties": {
            "name": {"type": "string", "maxLength": 255},
            "description": {"type": "string"},
            "is_default": {"type": "boolean"},
        },
    },
    "TenantRoleResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "guard_name": {"type": "string"},
            "description": {"type": "string"},
            "is_default": {"type": "boolean"},
            "is_system": {"type": "boolean"},
            "user_count": {"type": "integer"},
            "permission_ids": {"type": "array", "items": {"type": "string"}},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "TenantPermissionResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "guard_name": {"type": "string"},
            "resource": {"type": "string"},
            "action": {"type": "string"},
        },
    },
    "TenantAssignPermissionsRequest": {
        "type": "object", "required": ["permission_ids"],
        "properties": {"permission_ids": {"type": "array", "items": {"type": "string", "format": "uuid"}}},
    },
    "TenantUserResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "email": {"type": "string", "format": "email"},
            "is_active": {"type": "boolean"},
            "role_ids": {"type": "array", "items": {"type": "string"}},
            "role_names": {"type": "array", "items": {"type": "string"}},
            "created_at": {"type": "string", "format": "date-time"},
        },
    },
    "TenantAssignUserRolesRequest": {
        "type": "object", "required": ["role_ids"],
        "properties": {"role_ids": {"type": "array", "items": {"type": "string", "format": "uuid"}}},
    },
    # --- User Accounts ---
    "CreateAccountRequest": {
        "type": "object", "required": ["email"],
        "properties": {"email": {"type": "string", "format": "email", "maxLength": 255}},
    },
    "AccountResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "company_id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "user_id": {"type": "string", "format": "uuid"},
            "email": {"type": "string", "format": "email"},
            "role_name": {"type": "string"},
            "setup_token": {"type": "string"},
            "setup_token_expires": {"type": "string", "format": "date-time"},
            "password_set": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "SetPasswordRequest": {
        "type": "object", "required": ["token", "new_password"],
        "properties": {
            "token": {"type": "string", "description": "Setup token dari email link"},
            "new_password": {"type": "string", "format": "password", "minLength": 8, "maxLength": 72},
        },
    },
    # --- Employee Movement — Submit ---
    "SubmitMovementRequest": {
        "type": "object",
        "properties": {
            "flow_id": {"type": "string", "format": "uuid", "nullable": True, "description": "ID approval flow yang dipakai; wajib saat routing ke approval engine"},
        },
    },
    # --- Platform User ---
    "PlatformUserResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "email": {"type": "string", "format": "email"},
            "name": {"type": "string"},
            "role": {"type": "string", "enum": ["super_admin", "company_admin"]},
            "is_active": {"type": "boolean"},
            "company_id": {"type": "string", "format": "uuid"},
            "company_name": {"type": "string"},
            "last_login": {"type": "string", "format": "date-time"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },
}

for name, schema in new_schemas.items():
    if name not in schemas:
        schemas[name] = schema
        added += 1

# =========================================================================
# Fix previously-injected list responses to non-paginated arrays
# (handlers return {success, data: [...]} without pagination)
# =========================================================================
list_fixups = {
    f"{K}/evaluation-details/{{id}}/progress": ("listPerformanceProgressByDetailID", "PerformanceProgressResponse", "List of performance progress entries"),
    f"{K}/evaluations/{{id}}/comments": ("listPerformanceCommentsByEvaluationID", "PerformanceCommentResponse", "List of performance comments"),
    f"{K}/evaluation-details/{{id}}/attachments": ("listPerformanceAttachmentsByDetailID", "PerformanceAttachmentResponse", "List of performance attachments"),
    f"{R}/roles": ("listTenantRbacRoles", "TenantRoleResponse", "List of tenant RBAC roles"),
    f"{R}/permissions": ("listTenantRbacPermissions", "TenantPermissionResponse", "List of tenant RBAC permissions"),
    f"{R}/users": ("listTenantRbacUsers", "TenantUserResponse", "List of tenant RBAC users"),
}
for _path, (_oid, _schema, _desc) in list_fixups.items():
    _entry = paths.get(_path)
    if _entry and "get" in _entry:
        _new_resp = responses_array(_schema, _desc)
        if _entry["get"].get("responses") != _new_resp:
            _entry["get"]["responses"] = _new_resp
            _entry["get"].pop("parameters", None)
            added += 1

# =========================================================================
# Remove phantom endpoints (documented but NOT registered in routes.go)
# =========================================================================
phantom_paths = [
    "/api/v1/tenant/payroll/pph21-ptkp-rates/{id}",
    "/api/v1/tenant/payroll/pph21-tax-brackets/{id}",
    "/api/v1/tenant/settings/districts/{id}/villages",
    "/api/v1/tenant/settings/provinces/{id}/regencies",
    "/api/v1/tenant/settings/regencies/{id}/districts",
]
for _p in phantom_paths:
    if _p in paths:
        del paths[_p]
        added += 1
        print(f"  removed phantom path {_p}")

# =========================================================================
# Remove schemas orphaned by phantom endpoint removal
# =========================================================================
orphan_schemas = [
    "UpdatePph21PtkpRateRequest",
    "UpdatePph21TaxBracketRequest",
]
for _s in orphan_schemas:
    if _s in schemas:
        del schemas[_s]
        added += 1
        print(f"  removed orphaned schema {_s}")

# =========================================================================
# Tags
# =========================================================================
existing_tags = {t["name"] for t in spec.get("tags", [])}
new_tags = [
    {"name": "Tenant: RBAC Management", "description": "Tenant RBAC role & permission management"},
    {"name": "Tenant: User Accounts", "description": "Employee login account management"},
]
for tag in new_tags:
    if tag["name"] not in existing_tags:
        spec.setdefault("tags", []).append(tag)
        added += 1

with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

total_paths = len(spec["paths"])
total_schemas = len(spec["components"]["schemas"])
print(f"[OK] Injected/updated {added} items into {JSON_PATH}")
print(f"     Total paths:   {total_paths}")
print(f"     Total schemas: {total_schemas}")
