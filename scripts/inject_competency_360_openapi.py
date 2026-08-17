#!/usr/bin/env python3
"""Inject missing competency endpoints into openapi.json (idempotent).

Covers endpoints registered in backend/internal/modules/competency/routes.go
(Competency 360 module) that are missing from the OpenAPI spec:
- Rating scales, assessment templates (+ template indicators), indicators
- Rater assignment (event-targets raters, suggested-raters, delete rater)
- Manager assessments, my assessments (detail, responses, submit, approval)
- Employee result/gap/report, manager report, HR report

Usage:
    python scripts/inject_competency_360_openapi.py
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


def responses_plain(desc_200="OK"):
    return {
        "200": {"description": desc_200, "content": {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}}}}}},
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
        "tags": ["Competency"],
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

    C = "/api/v1/tenant/competency"

    # ------------------------------------------------------------------
    # Schemas
    # ------------------------------------------------------------------

    add_schema(spec, "RatingScaleItemRequest", obj({
        "value": prop("integer"),
        "label": prop("string"),
        "description": prop("string"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
    }, required=["value", "label"]))

    add_schema(spec, "CreateRatingScaleRequest", obj({
        "name": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "items": arr(ref("RatingScaleItemRequest")),
    }, required=["name"]))

    add_schema(spec, "UpdateRatingScaleRequest", obj({
        "name": prop("string"),
        "code": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "items": arr(ref("RatingScaleItemRequest")),
    }))

    add_schema(spec, "RatingScaleItemResponse", obj({
        "id": prop("string", fmt="uuid"),
        "scale_id": prop("string", fmt="uuid"),
        "value": prop("integer"),
        "label": prop("string"),
        "description": prop("string"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "RatingScaleResponse", obj({
        "id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "code": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "items": arr(ref("RatingScaleItemResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "TemplateCompetencyRequest", obj({
        "competency_id": prop("string", fmt="uuid"),
        "required_level": prop("integer"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
    }, required=["competency_id"]))

    add_schema(spec, "TemplateRaterTypeRequest", obj({
        "rater_type": prop("string", enum=["self", "superior", "peer", "subordinate", "other"]),
        "weight": prop("number"),
        "min_rater": prop("integer"),
        "max_rater": prop("integer"),
        "required": prop("boolean"),
        "anonymous": prop("boolean"),
    }, required=["rater_type"]))

    add_schema(spec, "CreateAssessmentTemplateRequest", obj({
        "name": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "scale_id": prop("string", fmt="uuid"),
        "competencies": arr(ref("TemplateCompetencyRequest")),
        "rater_types": arr(ref("TemplateRaterTypeRequest")),
    }, required=["name"]))

    add_schema(spec, "UpdateAssessmentTemplateRequest", obj({
        "name": prop("string"),
        "code": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "scale_id": prop("string", fmt="uuid"),
        "competencies": arr(ref("TemplateCompetencyRequest")),
        "rater_types": arr(ref("TemplateRaterTypeRequest")),
    }))

    add_schema(spec, "TemplateCompetencyResponse", obj({
        "id": prop("string", fmt="uuid"),
        "template_id": prop("string", fmt="uuid"),
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("integer"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
    }))

    add_schema(spec, "TemplateRaterTypeResponse", obj({
        "id": prop("string", fmt="uuid"),
        "template_id": prop("string", fmt="uuid"),
        "rater_type": prop("string"),
        "weight": prop("number"),
        "min_rater": prop("integer"),
        "max_rater": prop("integer"),
        "required": prop("boolean"),
        "anonymous": prop("boolean"),
    }))

    add_schema(spec, "AssessmentTemplateResponse", obj({
        "id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "code": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "scale_id": prop("string", fmt="uuid"),
        "competencies": arr(ref("TemplateCompetencyResponse")),
        "rater_types": arr(ref("TemplateRaterTypeResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateIndicatorRequest", obj({
        "competency_id": prop("string", fmt="uuid"),
        "code": prop("string"),
        "statement": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "sort_order": prop("integer"),
    }, required=["competency_id", "statement"]))

    add_schema(spec, "UpdateIndicatorRequest", obj({
        "competency_id": prop("string", fmt="uuid"),
        "code": prop("string"),
        "statement": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "sort_order": prop("integer"),
    }))

    add_schema(spec, "IndicatorResponse", obj({
        "id": prop("string", fmt="uuid"),
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "code": prop("string"),
        "statement": prop("string"),
        "description": prop("string"),
        "status": prop("string", enum=["active", "inactive"]),
        "sort_order": prop("integer"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "TemplateIndicatorRequest", obj({
        "indicator_id": prop("string", fmt="uuid"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
    }, required=["indicator_id"]))

    add_schema(spec, "TemplateIndicatorResponse", obj({
        "id": prop("string", fmt="uuid"),
        "template_id": prop("string", fmt="uuid"),
        "indicator_id": prop("string", fmt="uuid"),
        "statement": prop("string"),
        "weight": prop("number"),
        "sort_order": prop("integer"),
    }))

    add_schema(spec, "RaterAssignmentRequest", obj({
        "rater_employee_id": prop("string", fmt="uuid"),
        "rater_type": prop("string", enum=["self", "superior", "peer", "subordinate", "other"]),
        "weight": prop("number"),
    }, required=["rater_employee_id", "rater_type"]))

    add_schema(spec, "AssignRatersRequest", obj({
        "raters": arr(ref("RaterAssignmentRequest")),
    }, required=["raters"]))

    add_schema(spec, "RaterResponse", obj({
        "id": prop("string", fmt="uuid"),
        "competency_event_target_id": prop("string", fmt="uuid"),
        "competency_event_id": prop("string", fmt="uuid"),
        "rater_employee_id": prop("string", fmt="uuid"),
        "rater_employee_name": prop("string"),
        "subject_employee_id": prop("string", fmt="uuid"),
        "subject_employee_name": prop("string"),
        "rater_type": prop("string", enum=["self", "superior", "peer", "subordinate", "other"]),
        "weight": prop("number"),
        "status": prop("string"),
        "assigned_at": prop("string", fmt="date-time"),
        "submitted_at": prop("string", fmt="date-time"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "EmployeeBriefDTO", obj({
        "id": prop("string", fmt="uuid"),
        "name": prop("string"),
    }))

    add_schema(spec, "SuggestedRatersDTO", obj({
        "self": ref("EmployeeBriefDTO"),
        "superior": arr(ref("EmployeeBriefDTO")),
        "subordinates": arr(ref("EmployeeBriefDTO")),
    }))

    add_schema(spec, "SaveResponseRequest", obj({
        "indicator_id": prop("string", fmt="uuid"),
        "rating_value": prop("integer"),
        "comment": prop("string"),
    }, required=["indicator_id", "rating_value"]))

    add_schema(spec, "SaveResponsesRequest", obj({
        "responses": arr(ref("SaveResponseRequest")),
    }, required=["responses"]))

    add_schema(spec, "AssessmentResponseDTO", obj({
        "id": prop("string", fmt="uuid"),
        "rater_id": prop("string", fmt="uuid"),
        "indicator_id": prop("string", fmt="uuid"),
        "statement": prop("string"),
        "rating_value": prop("integer"),
        "comment": prop("string"),
        "submitted_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "AssessmentDetailDTO", obj({
        "rater": ref("RaterResponse"),
        "target": ref("CompetencyEventTargetResponse"),
        "indicators": arr(ref("TemplateIndicatorResponse")),
        "responses": arr(ref("AssessmentResponseDTO")),
        "scale": ref("RatingScaleResponse"),
    }))

    add_schema(spec, "ManagerAssessmentItem", obj({
        "employee_id": prop("string", fmt="uuid"),
        "employee_name": prop("string"),
        "target_id": prop("string", fmt="uuid"),
        "competency_event_id": prop("string", fmt="uuid"),
        "rater_id": prop("string", fmt="uuid"),
        "rater_status": prop("string"),
        "assigned_at": prop("string", fmt="date-time"),
        "submitted_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "GapItem", obj({
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("number"),
        "score": prop("number"),
        "gap": prop("number"),
        "weighted_gap": prop("number"),
    }))

    add_schema(spec, "GapAnalysisResponse", obj({
        "target_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "overall_score": prop("number"),
        "total_gap": prop("number"),
        "self_score": prop("number"),
        "others_score": prop("number"),
        "perception_gap": prop("number"),
        "strengths": arr(ref("GapItem")),
        "development_areas": arr(ref("GapItem")),
    }))

    add_schema(spec, "CompetencyScoreResult", obj({
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("number"),
        "score": prop("number"),
        "gap": prop("number"),
        "weighted_gap": prop("number"),
        "weight": prop("number"),
        "rater_scores": obj({}, ),
    }))

    add_schema(spec, "EmployeeReportDTO", obj({
        "employee_id": prop("string", fmt="uuid"),
        "target_id": prop("string", fmt="uuid"),
        "event_id": prop("string", fmt="uuid"),
        "overall_score": prop("number"),
        "total_gap": prop("number"),
        "self_score": prop("number"),
        "others_score": prop("number"),
        "perception_gap": prop("number"),
        "competencies": arr(ref("CompetencyScoreResult")),
        "strengths": arr(ref("GapItem")),
        "development_areas": arr(ref("GapItem")),
        "comments": arr(prop("string")),
    }))

    add_schema(spec, "EmployeeSummaryItem", obj({
        "employee_id": prop("string", fmt="uuid"),
        "target_id": prop("string", fmt="uuid"),
        "overall_score": prop("number"),
        "total_gap": prop("number"),
        "status": prop("string"),
        "rater_completion": prop("integer", desc="Persen (0-100)"),
    }))

    add_schema(spec, "ManagerReportDTO", obj({
        "event_id": prop("string", fmt="uuid"),
        "employees": arr(ref("EmployeeSummaryItem")),
        "total_employees": prop("integer"),
        "avg_score": prop("number"),
    }))

    add_schema(spec, "HRReportDTO", obj({
        "event_id": prop("string", fmt="uuid"),
        "total_targets": prop("integer", fmt="int64"),
        "finalized_targets": prop("integer", fmt="int64"),
        "rater_completion": prop("integer", desc="Persen"),
        "avg_score": prop("number"),
        "top_strengths": arr(ref("GapItem")),
        "top_development_gaps": arr(ref("GapItem")),
    }))

    # ------------------------------------------------------------------
    # Endpoints
    # ------------------------------------------------------------------

    # Rating scales
    add_endpoint(spec, "POST", f"{C}/rating-scales", "Buat skala penilaian (rating scale)",
                 request_body="CreateRatingScaleRequest", responses=responses_created("RatingScaleResponse"))
    add_endpoint(spec, "GET", f"{C}/rating-scales", "Daftar skala penilaian (paginated)",
                 responses=responses_list("CompetencyPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{C}/rating-scales/{{id}}", "Detail skala penilaian",
                 responses=responses_ok("RatingScaleResponse"))
    add_endpoint(spec, "PUT", f"{C}/rating-scales/{{id}}", "Perbarui skala penilaian",
                 request_body="UpdateRatingScaleRequest", responses=responses_ok("RatingScaleResponse"))
    add_endpoint(spec, "DELETE", f"{C}/rating-scales/{{id}}", "Hapus skala penilaian",
                 responses=responses_plain("Rating scale deleted"))

    # Assessment templates
    add_endpoint(spec, "POST", f"{C}/templates", "Buat template assessment 360",
                 request_body="CreateAssessmentTemplateRequest", responses=responses_created("AssessmentTemplateResponse"))
    add_endpoint(spec, "GET", f"{C}/templates", "Daftar template assessment (paginated)",
                 responses=responses_list("CompetencyPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{C}/templates/{{id}}", "Detail template assessment",
                 responses=responses_ok("AssessmentTemplateResponse"))
    add_endpoint(spec, "PUT", f"{C}/templates/{{id}}", "Perbarui template assessment",
                 request_body="UpdateAssessmentTemplateRequest", responses=responses_ok("AssessmentTemplateResponse"))
    add_endpoint(spec, "DELETE", f"{C}/templates/{{id}}", "Hapus template assessment",
                 responses=responses_plain("Template deleted"))

    # Template indicators
    add_endpoint(spec, "GET", f"{C}/templates/{{id}}/indicators", "Daftar indikator dalam template",
                 responses=responses_array("TemplateIndicatorResponse", "List of template indicators"))
    add_endpoint(spec, "PUT", f"{C}/templates/{{id}}/indicators", "Set daftar indikator template",
                 request_body={"type": "object", "properties": {"indicators": arr(ref("TemplateIndicatorRequest"))}, "required": ["indicators"]},
                 responses=responses_array("TemplateIndicatorResponse", "Updated template indicators"))

    # Indicators
    add_endpoint(spec, "POST", f"{C}/indicators", "Buat indikator kompetensi",
                 request_body="CreateIndicatorRequest", responses=responses_created("IndicatorResponse"))
    add_endpoint(spec, "GET", f"{C}/indicators", "Daftar indikator kompetensi (paginated)",
                 responses=responses_list("CompetencyPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{C}/indicators/{{id}}", "Detail indikator kompetensi",
                 responses=responses_ok("IndicatorResponse"))
    add_endpoint(spec, "PUT", f"{C}/indicators/{{id}}", "Perbarui indikator kompetensi",
                 request_body="UpdateIndicatorRequest", responses=responses_ok("IndicatorResponse"))
    add_endpoint(spec, "DELETE", f"{C}/indicators/{{id}}", "Hapus indikator kompetensi",
                 responses=responses_plain("Indicator deleted"))

    # Rater assignment
    add_endpoint(spec, "POST", f"{C}/event-targets/{{id}}/raters", "Assign rater ke event target",
                 request_body="AssignRatersRequest", responses=responses_array("RaterResponse", "Assigned raters"))
    add_endpoint(spec, "GET", f"{C}/event-targets/{{id}}/raters", "Daftar rater event target",
                 responses=responses_array("RaterResponse", "List of raters"))
    add_endpoint(spec, "GET", f"{C}/event-targets/{{id}}/suggested-raters", "Saran rater (self/superior/subordinates)",
                 responses=responses_ok("SuggestedRatersDTO"))
    add_endpoint(spec, "DELETE", f"{C}/raters/{{id}}", "Hapus rater",
                 responses=responses_plain("Rater deleted"))

    # Manager assessments
    add_endpoint(spec, "GET", f"{C}/manager-assessments", "Daftar assessment bawahan (untuk manager)",
                 responses=responses_array("ManagerAssessmentItem", "List of manager assessment items"))

    # My assessments
    add_endpoint(spec, "GET", f"{C}/my-assessments", "Daftar assessment milik saya (sebagai rater)",
                 responses=responses_array("RaterResponse", "List of my assessments"))
    add_endpoint(spec, "GET", f"{C}/my-assessments/{{id}}", "Detail assessment saya (indikator + respons)",
                 responses=responses_ok("AssessmentDetailDTO"))
    add_endpoint(spec, "POST", f"{C}/my-assessments/{{id}}/responses", "Simpan respons assessment (draft)",
                 request_body="SaveResponsesRequest", responses=responses_ok("AssessmentDetailDTO"))
    add_endpoint(spec, "POST", f"{C}/my-assessments/{{id}}/submit", "Submit assessment",
                 responses=responses_ok("RaterResponse"))

    # Approval
    add_endpoint(spec, "POST", f"{C}/event-targets/{{id}}/submit-approval", "Submit assessment event target untuk approval",
                 responses=responses_ok("CompetencyEventTargetResponse"))

    # Result, gap & reports
    add_endpoint(spec, "GET", f"{C}/employees/{{employee}}/result", "Hasil assessment seorang employee",
                 responses=responses_ok("CompetencyScoreResult"))
    add_endpoint(spec, "GET", f"{C}/employees/{{employee}}/gap", "Analisis gap kompetensi employee",
                 responses=responses_ok("GapAnalysisResponse"))
    add_endpoint(spec, "GET", f"{C}/employees/{{employee}}/report", "Laporan assessment employee",
                 responses=responses_ok("EmployeeReportDTO"))
    add_endpoint(spec, "GET", f"{C}/reports/manager", "Laporan manager (overview employee per event)",
                 responses=responses_ok("ManagerReportDTO"),
                 query=[qparam("event_id", {"type": "string", "format": "uuid"})])
    add_endpoint(spec, "GET", f"{C}/reports/hr", "Laporan HR per event (distribusi & completion)",
                 responses=responses_ok("HRReportDTO"),
                 query=[qparam("event_id", {"type": "string", "format": "uuid"})])

    with open(JSON_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)
    print("OK: competency 360 endpoints injected into openapi.json")


if __name__ == "__main__":
    main()
