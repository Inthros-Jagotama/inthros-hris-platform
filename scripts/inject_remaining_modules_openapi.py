#!/usr/bin/env python3
"""Inject missing endpoints for remaining modules into openapi.json (idempotent).

Covers endpoints registered in routes.go that are missing from the OpenAPI spec:
- document-templates (documenttemplate module): CRUD, versions, activate/deactivate, preview
- employee-movements (employeemovement): generate-document + generated-documents (movement & contract)
- employees (employee module): stats/gender, stats/employment-status, settings/sensitive-fields
- document-numbering (setting module): list, update, preview
- approval: GET /approval/tasks/done
- job-management: GET /job-management/dashboard
- leave: GET /leave/reports/on-leave-today

Usage:
    python scripts/inject_remaining_modules_openapi.py
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
            "content": {"application/json": {"schema": {"type": "object", "properties": {
                "success": {"type": "boolean", "example": True},
                "data": {"type": "array", "items": ref(item_ref)}}}}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_plain(desc_200="OK", content_type="application/json"):
    c = {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}}}}}
    if content_type != "application/json":
        c = {content_type: {"schema": {"type": "string"}}}
    return {
        "200": {"description": desc_200, "content": c},
        "400": {"description": "Bad request / Validation error"},
    }


def prop(type_, fmt=None, enum=None, desc=None, example=None):
    p = {"type": type_}
    if fmt:
        p["format"] = fmt
    if enum:
        p["enum"] = enum
    if desc:
        p["description"] = desc
    if example is not None:
        p["example"] = example
    return p


def obj(props, required=None):
    o = {"type": "object", "properties": props}
    if required:
        o["required"] = required
    return o


def arr(items):
    return {"type": "array", "items": items}


def add_schema(spec, name, schema):
    spec["components"]["schemas"][name] = schema


def add_endpoint(spec, method, path, summary, request_body=None, responses=None, query=None):
    op = {
        "tags": ["Tenant: Others"],
        "summary": summary,
        "security": SEC,
        "responses": responses or responses_plain(),
    }
    params = []
    for p in path.split("/"):
        if p.startswith("{") and p.endswith("}"):
            params.append(param(p[1:-1]))
    if query:
        params.extend(query)
    if params:
        op["parameters"] = params
    if request_body:
        op["requestBody"] = content(request_body)
    spec["paths"].setdefault(path, {})[method.lower()] = op


def main():
    with open(JSON_PATH, encoding="utf-8") as f:
        spec = json.load(f)

    # ------------------------------------------------------------------
    # Schemas — Document Templates
    # ------------------------------------------------------------------

    add_schema(spec, "CreateTemplateRequest", obj({
        "name": prop("string"),
        "code": prop("string", desc="Kosong → digenerate otomatis"),
        "document_type": prop("string", enum=["CONTRACT_AGREEMENT", "MOVEMENT_SK"]),
        "movement_type": prop("string", desc="Hanya untuk MOVEMENT_SK; kosong = template umum"),
        "description": prop("string"),
    }, required=["name", "document_type"]))

    add_schema(spec, "UpdateTemplateRequest", obj({
        "name": prop("string"),
        "description": prop("string"),
    }))

    add_schema(spec, "DocumentTemplateResponse", obj({
        "id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "code": prop("string"),
        "document_type": prop("string"),
        "movement_type": prop("string"),
        "description": prop("string"),
        "content": prop("string"),
        "active_version_id": prop("string", fmt="uuid"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "is_active": prop("boolean"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "TemplateListResponse", obj({
        "data": arr(ref("DocumentTemplateResponse")),
        "total": prop("integer", fmt="int64"),
        "page": prop("integer"),
    }))

    add_schema(spec, "MovementTypeOption", obj({
        "value": prop("string"),
        "label": prop("string"),
    }))

    add_schema(spec, "CreateVersionRequest", obj({
        "content": prop("string", desc="Konten template (DOCX/HTML)"),
        "paper_size": prop("string", enum=["A4", "A5", "Letter", "Legal"]),
        "orientation": prop("string", enum=["portrait", "landscape"]),
        "margin_top": prop("integer"),
        "margin_right": prop("integer"),
        "margin_bottom": prop("integer"),
        "margin_left": prop("integer"),
    }, required=["content"]))

    add_schema(spec, "DocumentTemplateVersionResponse", obj({
        "id": prop("string", fmt="uuid"),
        "template_id": prop("string", fmt="uuid"),
        "version": prop("integer"),
        "content": prop("string"),
        "file_name": prop("string"),
        "file_url": prop("string"),
        "paper_size": prop("string"),
        "orientation": prop("string"),
        "margin_top": prop("integer"),
        "margin_right": prop("integer"),
        "margin_bottom": prop("integer"),
        "margin_left": prop("integer"),
        "created_by": prop("string", fmt="uuid"),
        "created_at": prop("string", fmt="date-time"),
    }))

    # ------------------------------------------------------------------
    # Schemas — Generated documents (employeemovement)
    # ------------------------------------------------------------------

    add_schema(spec, "GeneratedDocumentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "template_id": prop("string", fmt="uuid"),
        "template_version_id": prop("string", fmt="uuid"),
        "document_type": prop("string"),
        "reference_type": prop("string", enum=["movement", "contract"]),
        "reference_id": prop("string", fmt="uuid"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "generated_by": prop("string", fmt="uuid"),
        "generated_at": prop("string", fmt="date-time"),
        "created_at": prop("string", fmt="date-time"),
        "file_url": prop("string"),
    }))

    add_schema(spec, "PaginatedGeneratedDocumentResponse", obj({
        "success": prop("boolean"),
        "data": arr(ref("GeneratedDocumentResponse")),
        "page": prop("integer"),
        "per_page": prop("integer"),
        "total": prop("integer", fmt="int64"),
        "total_pages": prop("integer"),
    }))

    # ------------------------------------------------------------------
    # Schemas — Employees (stats & sensitive fields)
    # ------------------------------------------------------------------

    add_schema(spec, "GenderStatsResponse", obj({
        "male": prop("integer", fmt="int64"),
        "female": prop("integer", fmt="int64"),
        "other": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "EmploymentStatusCount", obj({
        "name": prop("string"),
        "count": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "EmploymentStatusStatsResponse", obj({
        "groups": arr(ref("EmploymentStatusCount")),
        "unclassified": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "SensitiveFieldSettingResponse", obj({
        "id": prop("string", fmt="uuid"),
        "field_key": prop("string"),
        "is_encryption_enabled": prop("boolean"),
        "updated_by": prop("string", fmt="uuid"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "UpdateSensitiveFieldRequest", obj({
        "is_encryption_enabled": prop("boolean"),
    }, required=["is_encryption_enabled"]))

    # ------------------------------------------------------------------
    # Schemas — Document numbering
    # ------------------------------------------------------------------

    add_schema(spec, "DocumentNumberingSettingResponse", obj({
        "id": prop("string", fmt="uuid"),
        "document_type": prop("string"),
        "format_template": prop("string"),
        "reset_period": prop("string", enum=["DAILY", "MONTHLY", "YEARLY", "NEVER"]),
        "last_sequence": prop("integer"),
        "last_reset_key": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "UpdateDocumentNumberingRequest", obj({
        "format_template": prop("string"),
        "reset_period": prop("string", enum=["DAILY", "MONTHLY", "YEARLY", "NEVER"]),
    }, required=["format_template", "reset_period"]))

    # ------------------------------------------------------------------
    # Schemas — Job management dashboard & leave
    # ------------------------------------------------------------------

    add_schema(spec, "JobManagementDashboardResponse", obj({
        "summary": obj({}, ),
        "total_organizations": prop("integer"),
        "with_employees": prop("integer"),
        "without_employees": prop("integer"),
        "value_not_started": prop("integer"),
        "value_on_progress": prop("integer"),
        "value_completed": prop("integer"),
        "with_financial_authority": prop("integer"),
        "without_financial_authority": prop("integer"),
    }))

    add_schema(spec, "OnLeaveTodayResponse", obj({
        "count": prop("integer"),
    }))

    # ------------------------------------------------------------------
    # Endpoints — Document Templates
    # ------------------------------------------------------------------
    DT = "/api/v1/tenant/document-templates"

    add_endpoint(spec, "GET", DT, "Daftar template dokumen (paginated)",
                 responses=responses_ok("TemplateListResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"}),
                        qparam("document_type", {"type": "string"}), qparam("movement_type", {"type": "string"})])
    add_endpoint(spec, "POST", DT, "Buat template dokumen",
                 request_body="CreateTemplateRequest", responses=responses_created("DocumentTemplateResponse"))
    add_endpoint(spec, "GET", f"{DT}/movement-types", "Daftar jenis movement untuk template SK",
                 responses=responses_array("MovementTypeOption", "List of movement type options"))
    add_endpoint(spec, "GET", f"{DT}/variables", "Daftar variabel yang tersedia untuk template",
                 responses=responses_array({"type": "object"}, "List of variable groups"))
    add_endpoint(spec, "GET", f"{DT}/{{id}}", "Detail template dokumen",
                 responses=responses_ok("DocumentTemplateResponse"))
    add_endpoint(spec, "PUT", f"{DT}/{{id}}", "Perbarui template dokumen",
                 request_body="UpdateTemplateRequest", responses=responses_ok("DocumentTemplateResponse"))
    add_endpoint(spec, "DELETE", f"{DT}/{{id}}", "Hapus template dokumen",
                 responses=responses_plain("Template deleted"))
    add_endpoint(spec, "POST", f"{DT}/{{id}}/preview", "Preview template (render PDF)",
                 responses=responses_plain("Rendered preview file", content_type="application/pdf"))
    add_endpoint(spec, "POST", f"{DT}/{{id}}/activate", "Aktifkan template",
                 responses=responses_ok("DocumentTemplateResponse"))
    add_endpoint(spec, "POST", f"{DT}/{{id}}/deactivate", "Nonaktifkan template",
                 responses=responses_ok("DocumentTemplateResponse"))
    add_endpoint(spec, "GET", f"{DT}/{{id}}/versions", "Daftar versi template",
                 responses=responses_array("DocumentTemplateVersionResponse", "List of template versions"))
    add_endpoint(spec, "POST", f"{DT}/{{id}}/versions", "Buat versi template baru",
                 request_body="CreateVersionRequest", responses=responses_created("DocumentTemplateVersionResponse"))
    add_endpoint(spec, "GET", f"{DT}/{{id}}/versions/{{versionId}}", "Detail versi template",
                 responses=responses_ok("DocumentTemplateVersionResponse"))

    # ------------------------------------------------------------------
    # Endpoints — Employee Movements (generated documents)
    # ------------------------------------------------------------------
    EM = "/api/v1/tenant/employee-movements"

    add_endpoint(spec, "POST", f"{EM}/movements/{{id}}/generate-document", "Generate PDF SK Movement dari template aktif",
                 responses=responses_created("GeneratedDocumentResponse"))
    add_endpoint(spec, "GET", f"{EM}/movements/{{id}}/generated-documents", "Daftar dokumen SK yang sudah digenerate",
                 responses=responses_ok("PaginatedGeneratedDocumentResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "POST", f"{EM}/contracts/{{id}}/generate-document", "Generate PDF Perjanjian Kerja dari template aktif",
                 responses=responses_created("GeneratedDocumentResponse"))
    add_endpoint(spec, "GET", f"{EM}/contracts/{{id}}/generated-documents", "Daftar dokumen kontrak yang sudah digenerate",
                 responses=responses_ok("PaginatedGeneratedDocumentResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])

    # ------------------------------------------------------------------
    # Endpoints — Employees (stats & sensitive fields)
    # ------------------------------------------------------------------
    EP = "/api/v1/tenant/employees"

    add_endpoint(spec, "GET", f"{EP}/stats/gender", "Statistik jumlah karyawan per jenis kelamin",
                 responses=responses_ok("GenderStatsResponse"))
    add_endpoint(spec, "GET", f"{EP}/stats/employment-status", "Statistik jumlah karyawan per status kepegawaian",
                 responses=responses_ok("EmploymentStatusStatsResponse"))
    add_endpoint(spec, "GET", f"{EP}/settings/sensitive-fields", "Daftar setelan field sensitif (enkripsi at-rest)",
                 responses=responses_array("SensitiveFieldSettingResponse", "List of sensitive field settings"))
    add_endpoint(spec, "PUT", f"{EP}/settings/sensitive-fields/{{fieldKey}}", "Aktifkan/nonaktifkan enkripsi field sensitif",
                 request_body="UpdateSensitiveFieldRequest", responses=responses_ok("SensitiveFieldSettingResponse"))

    # ------------------------------------------------------------------
    # Endpoints — Document Numbering
    # ------------------------------------------------------------------
    DN = "/api/v1/tenant/document-numbering"

    add_endpoint(spec, "GET", DN, "Daftar setelan penomoran dokumen",
                 responses=responses_array("DocumentNumberingSettingResponse", "List of numbering settings"))
    add_endpoint(spec, "PUT", f"{DN}/{{document_type}}", "Perbarui format penomoran dokumen",
                 request_body="UpdateDocumentNumberingRequest", responses=responses_ok("DocumentNumberingSettingResponse"))
    add_endpoint(spec, "GET", f"{DN}/{{document_type}}/preview", "Preview nomor dokumen berikutnya",
                 responses=responses_plain("Preview nomor dokumen", content_type="application/json"))

    # ------------------------------------------------------------------
    # Endpoints — Approval, Job Management, Leave
    # ------------------------------------------------------------------

    add_endpoint(spec, "GET", "/api/v1/tenant/approval/tasks/done", "Daftar task approval yang sudah diselesaikan (paginated)",
                 responses=responses_list("ApprovalPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"}),
                        qparam("flow_id", {"type": "string", "format": "uuid"})])

    add_endpoint(spec, "GET", "/api/v1/tenant/job-management/dashboard", "Dashboard ringkasan master data Job Management",
                 responses=responses_ok("JobManagementDashboardResponse"))

    add_endpoint(spec, "GET", "/api/v1/tenant/leave/reports/on-leave-today", "Jumlah karyawan yang sedang cuti hari ini",
                 responses=responses_ok("OnLeaveTodayResponse"))

    with open(JSON_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)
    print("OK: remaining module endpoints injected into openapi.json")


if __name__ == "__main__":
    main()
