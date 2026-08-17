#!/usr/bin/env python3
"""Inject missing recruitment endpoints into openapi.json (idempotent).

Covers endpoints registered in backend/internal/modules/recruitment/routes.go that
are missing from the OpenAPI spec:
- Requisition requirements + competencies (G-9)
- Job offers workflow (G-3: submit/send/accept/reject/withdraw)
- Candidate consents
- Recruitment analytics summary (G-11)
- Application history, match-score (G-9), screening (G-7), assessment (G-12)
- Assessments + participants (G-7 sub-project 2)
- Interviewers + scorecard items (G-8), interview complete

Usage:
    python scripts/inject_recruitment_updates_openapi.py
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
        "tags": ["Recruitment"],
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

    R = "/api/v1/tenant/recruitment"

    # ------------------------------------------------------------------
    # Schemas
    # ------------------------------------------------------------------

    add_schema(spec, "CreateRequisitionRequirementRequest", obj({
        "requirement_type": prop("string", desc="Jenis requirement (education/experience/dll)"),
        "name": prop("string"),
        "description": prop("string"),
        "minimum_value": prop("number"),
        "maximum_value": prop("number"),
        "is_required": prop("boolean"),
        "sort_order": prop("integer"),
    }, required=["requirement_type", "name"]))

    add_schema(spec, "UpdateRequisitionRequirementRequest", obj({
        "requirement_type": prop("string"),
        "name": prop("string"),
        "description": prop("string"),
        "minimum_value": prop("number"),
        "maximum_value": prop("number"),
        "is_required": prop("boolean"),
        "sort_order": prop("integer"),
    }))

    add_schema(spec, "RequisitionRequirementResponse", obj({
        "id": prop("string", fmt="uuid"),
        "requisition_id": prop("string", fmt="uuid"),
        "requirement_type": prop("string"),
        "name": prop("string"),
        "description": prop("string"),
        "minimum_value": prop("number"),
        "maximum_value": prop("number"),
        "is_required": prop("boolean"),
        "sort_order": prop("integer"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateRequisitionCompetencyRequest", obj({
        "competency_id": prop("string", fmt="uuid"),
        "required_level": prop("integer"),
        "is_required": prop("boolean"),
        "weight": prop("number"),
    }, required=["competency_id"]))

    add_schema(spec, "UpdateRequisitionCompetencyRequest", obj({
        "required_level": prop("integer"),
        "is_required": prop("boolean"),
        "weight": prop("number"),
    }))

    add_schema(spec, "RequisitionCompetencyResponse", obj({
        "id": prop("string", fmt="uuid"),
        "requisition_id": prop("string", fmt="uuid"),
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("integer"),
        "is_required": prop("boolean"),
        "weight": prop("number"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateOfferRequest", obj({
        "application_id": prop("string", fmt="uuid"),
        "employment_type": prop("string", desc="mis. FULL_TIME / CONTRACT"),
        "salary": prop("number"),
        "allowances": prop("number"),
        "benefits": prop("string"),
        "start_date": prop("string", fmt="date"),
        "expiry_date": prop("string", fmt="date"),
    }, required=["application_id"]))

    add_schema(spec, "UpdateOfferRequest", obj({
        "employment_type": prop("string"),
        "salary": prop("number"),
        "allowances": prop("number"),
        "benefits": prop("string"),
        "start_date": prop("string", fmt="date"),
        "expiry_date": prop("string", fmt="date"),
    }))

    add_schema(spec, "SubmitOfferRequest", obj({
        "flow_id": prop("string", fmt="uuid", desc="Opsional; auto-resolve dari active flow modul recruitment_offer"),
    }))

    add_schema(spec, "SubmitRequisitionRequest", obj({
        "flow_id": prop("string", fmt="uuid", desc="Opsional; auto-resolve dari active flow modul recruitment_requisition"),
    }))

    add_schema(spec, "OfferResponse", obj({
        "id": prop("string", fmt="uuid"),
        "application_id": prop("string", fmt="uuid"),
        "offer_number": prop("string"),
        "employment_type": prop("string"),
        "salary": prop("number"),
        "allowances": prop("number"),
        "benefits": prop("string"),
        "start_date": prop("string", fmt="date"),
        "expiry_date": prop("string", fmt="date"),
        "status": prop("string", enum=["DRAFT", "PENDING_APPROVAL", "APPROVED", "REJECTED", "SENT", "ACCEPTED", "DECLINED", "WITHDRAWN"]),
        "approval_instance_id": prop("string", fmt="uuid"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateCandidateConsentRequest", obj({
        "action": prop("string", enum=["GRANTED", "REVOKED"]),
        "notes": prop("string"),
    }, required=["action"]))

    add_schema(spec, "CandidateConsentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "candidate_id": prop("string", fmt="uuid"),
        "action": prop("string", enum=["GRANTED", "REVOKED"]),
        "notes": prop("string"),
        "changed_by": prop("string", fmt="uuid"),
        "changed_at": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "StageRef", obj({
        "code": prop("string"),
        "name": prop("string"),
    }))

    add_schema(spec, "StageHistoryResponse", obj({
        "id": prop("string", fmt="uuid"),
        "from_stage": ref("StageRef"),
        "to_stage": ref("StageRef"),
        "changed_by": prop("string", fmt="uuid"),
        "notes": prop("string"),
        "changed_at": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "CreateApplicationScreeningRequest", obj({
        "screened_by": prop("string", fmt="uuid"),
        "screened_at": prop("integer", fmt="int64"),
        "score": prop("number"),
        "result": prop("string", enum=["PASS", "FAIL", "HOLD"]),
        "notes": prop("string"),
    }))

    add_schema(spec, "UpdateApplicationScreeningRequest", obj({
        "screened_by": prop("string", fmt="uuid"),
        "screened_at": prop("integer", fmt="int64"),
        "score": prop("number"),
        "result": prop("string", enum=["PASS", "FAIL", "HOLD"]),
        "notes": prop("string"),
    }))

    add_schema(spec, "ApplicationScreeningResponse", obj({
        "id": prop("string", fmt="uuid"),
        "application_id": prop("string", fmt="uuid"),
        "screened_by": prop("string", fmt="uuid"),
        "screened_at": prop("integer", fmt="int64"),
        "score": prop("number"),
        "result": prop("string", enum=["PASS", "FAIL", "HOLD"]),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "AssessmentCompetencyLevel", obj({
        "competency_id": prop("string", fmt="uuid"),
        "level": prop("integer"),
    }))

    add_schema(spec, "SaveApplicationAssessmentRequest", obj({
        "education_match": prop("boolean"),
        "education_note": prop("string"),
        "experience_match": prop("boolean"),
        "experience_note": prop("string"),
        "competency_levels": arr(ref("AssessmentCompetencyLevel")),
    }))

    add_schema(spec, "AssessmentRequirementCompetency", obj({
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("integer"),
        "weight": prop("number"),
    }))

    add_schema(spec, "AssessmentRequirement", obj({
        "education": prop("string"),
        "experience": prop("string"),
        "competencies": arr(ref("AssessmentRequirementCompetency")),
    }))

    add_schema(spec, "AssessmentCompetencyBreakdown", obj({
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("integer"),
        "candidate_level": prop("integer"),
        "weight": prop("number"),
        "contribution": prop("number"),
    }))

    add_schema(spec, "ApplicationAssessmentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "application_id": prop("string", fmt="uuid"),
        "education_match": prop("boolean"),
        "education_note": prop("string"),
        "experience_match": prop("boolean"),
        "experience_note": prop("string"),
        "competency_levels": arr(ref("AssessmentCompetencyLevel")),
        "score": prop("number"),
        "breakdown": arr(ref("AssessmentCompetencyBreakdown")),
        "assessed_by": prop("string", fmt="uuid"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "ApplicationAssessmentDetailResponse", obj({
        "application_id": prop("string", fmt="uuid"),
        "requisition_id": prop("string", fmt="uuid"),
        "candidate_id": prop("string", fmt="uuid"),
        "requirements": ref("AssessmentRequirement"),
        "assessment": ref("ApplicationAssessmentResponse"),
    }))

    add_schema(spec, "CreateAssessmentRequest", obj({
        "requisition_id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "type": prop("string", enum=["TECHNICAL", "PSYCHOLOGICAL", "COGNITIVE", "PERSONALITY", "CASE_STUDY", "CODING", "LANGUAGE", "OTHER"]),
        "scheduled_at": prop("integer", fmt="int64"),
        "location": prop("string"),
        "meeting_link": prop("string"),
        "notes": prop("string"),
    }, required=["name"]))

    add_schema(spec, "UpdateAssessmentRequest", obj({
        "name": prop("string"),
        "type": prop("string", enum=["TECHNICAL", "PSYCHOLOGICAL", "COGNITIVE", "PERSONALITY", "CASE_STUDY", "CODING", "LANGUAGE", "OTHER"]),
        "scheduled_at": prop("integer", fmt="int64"),
        "location": prop("string"),
        "meeting_link": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "AssessmentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "requisition_id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "type": prop("string"),
        "scheduled_at": prop("integer", fmt="int64"),
        "location": prop("string"),
        "meeting_link": prop("string"),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "AddAssessmentParticipantRequest", obj({
        "application_id": prop("string", fmt="uuid"),
    }, required=["application_id"]))

    add_schema(spec, "UpdateAssessmentParticipantRequest", obj({
        "status": prop("string", enum=["INVITED", "COMPLETED", "NO_SHOW"]),
        "score": prop("number"),
        "result": prop("string", enum=["PASS", "FAIL", "HOLD"]),
        "recommendation": prop("string"),
    }))

    add_schema(spec, "AssessmentParticipantResponse", obj({
        "id": prop("string", fmt="uuid"),
        "assessment_id": prop("string", fmt="uuid"),
        "application_id": prop("string", fmt="uuid"),
        "status": prop("string", enum=["INVITED", "COMPLETED", "NO_SHOW"]),
        "score": prop("number"),
        "result": prop("string", enum=["PASS", "FAIL", "HOLD"]),
        "recommendation": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "MatchScoreCompetencyBreakdown", obj({
        "competency_id": prop("string", fmt="uuid"),
        "competency_name": prop("string"),
        "required_level": prop("integer"),
        "candidate_level": prop("integer"),
        "weight": prop("number"),
        "contribution": prop("number"),
    }))

    add_schema(spec, "MatchScoreResponse", obj({
        "application_id": prop("string", fmt="uuid"),
        "candidate_id": prop("string", fmt="uuid"),
        "requisition_id": prop("string", fmt="uuid"),
        "score": prop("number"),
        "breakdown": arr(ref("MatchScoreCompetencyBreakdown")),
        "note": prop("string"),
    }))

    add_schema(spec, "SourceConversionEntry", obj({
        "source": prop("string"),
        "applications": prop("integer", fmt="int64"),
        "hires": prop("integer", fmt="int64"),
        "conversion_rate": prop("number"),
    }))

    add_schema(spec, "RecruitmentAnalyticsSummaryResponse", obj({
        "open_requisitions": prop("integer", fmt="int64"),
        "candidates": prop("integer", fmt="int64"),
        "applications": prop("integer", fmt="int64"),
        "shortlisted": prop("integer", fmt="int64"),
        "interviews": prop("integer", fmt="int64"),
        "offers": prop("integer", fmt="int64"),
        "hires": prop("integer", fmt="int64"),
        "rejected": prop("integer", fmt="int64"),
        "withdrawn": prop("integer", fmt="int64"),
        "time_to_hire_days": prop("number"),
        "offer_acceptance_rate": prop("number"),
        "application_conversion_rate": prop("number"),
        "source_conversion": arr(ref("SourceConversionEntry")),
        "note": prop("string"),
    }))

    add_schema(spec, "AddInterviewerRequest", obj({
        "employee_id": prop("string", fmt="uuid"),
        "role": prop("string"),
    }, required=["employee_id"]))

    add_schema(spec, "InterviewerResponse", obj({
        "id": prop("string", fmt="uuid"),
        "interview_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "role": prop("string"),
        "created_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "AddScorecardItemRequest", obj({
        "criterion": prop("string"),
        "weight": prop("number"),
        "score": prop("number"),
        "notes": prop("string"),
    }, required=["criterion"]))

    add_schema(spec, "UpdateScorecardItemRequest", obj({
        "criterion": prop("string"),
        "weight": prop("number"),
        "score": prop("number"),
        "notes": prop("string"),
    }))

    add_schema(spec, "ScorecardItemResponse", obj({
        "id": prop("string", fmt="uuid"),
        "interview_id": prop("string", fmt="uuid"),
        "criterion": prop("string"),
        "weight": prop("number"),
        "score": prop("number"),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    # ------------------------------------------------------------------
    # Endpoints
    # ------------------------------------------------------------------

    # Requisition submit (G-1)
    add_endpoint(spec, "POST", f"{R}/requisitions/{{id}}/submit", "Submit requisition draft ke Central Approval",
                 request_body="SubmitRequisitionRequest", responses=responses_ok("RequisitionResponse"))

    # Requisition requirements (G-9)
    add_endpoint(spec, "POST", f"{R}/requisitions/{{id}}/requirements", "Tambah requirement ke job requisition",
                 request_body="CreateRequisitionRequirementRequest", responses=responses_created("RequisitionRequirementResponse"))
    add_endpoint(spec, "GET", f"{R}/requisitions/{{id}}/requirements", "Daftar requirement job requisition",
                 responses=responses_array("RequisitionRequirementResponse", "List of requisition requirements"))
    add_endpoint(spec, "PUT", f"{R}/requirements/{{id}}", "Perbarui requirement job requisition",
                 request_body="UpdateRequisitionRequirementRequest", responses=responses_ok("RequisitionRequirementResponse"))
    add_endpoint(spec, "DELETE", f"{R}/requirements/{{id}}", "Hapus requirement job requisition",
                 responses=responses_plain("Requirement deleted"))

    # Requisition competencies (G-9)
    add_endpoint(spec, "POST", f"{R}/requisitions/{{id}}/competencies", "Tambah kompetensi ke job requisition",
                 request_body="CreateRequisitionCompetencyRequest", responses=responses_created("RequisitionCompetencyResponse"))
    add_endpoint(spec, "GET", f"{R}/requisitions/{{id}}/competencies", "Daftar kompetensi job requisition",
                 responses=responses_array("RequisitionCompetencyResponse", "List of requisition competencies"))
    add_endpoint(spec, "PUT", f"{R}/requisition-competencies/{{id}}", "Perbarui kompetensi job requisition",
                 request_body="UpdateRequisitionCompetencyRequest", responses=responses_ok("RequisitionCompetencyResponse"))
    add_endpoint(spec, "DELETE", f"{R}/requisition-competencies/{{id}}", "Hapus kompetensi job requisition",
                 responses=responses_plain("Competency deleted"))

    # Job offers (G-3)
    add_endpoint(spec, "POST", f"{R}/offers", "Buat penawaran kerja (offer)",
                 request_body="CreateOfferRequest", responses=responses_created("OfferResponse"))
    add_endpoint(spec, "GET", f"{R}/offers", "Daftar penawaran kerja (paginated)",
                 responses=responses_list("RecruitmentPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{R}/offers/{{id}}", "Detail penawaran kerja",
                 responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "PUT", f"{R}/offers/{{id}}", "Perbarui penawaran kerja",
                 request_body="UpdateOfferRequest", responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "DELETE", f"{R}/offers/{{id}}", "Hapus penawaran kerja",
                 responses=responses_plain("Offer deleted"))
    add_endpoint(spec, "POST", f"{R}/offers/{{id}}/submit", "Submit offer draft ke Central Approval",
                 request_body="SubmitOfferRequest", responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "POST", f"{R}/offers/{{id}}/send", "Kirim offer ke kandidat",
                 responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "POST", f"{R}/offers/{{id}}/accept", "Terima offer (kandidat menerima)",
                 responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "POST", f"{R}/offers/{{id}}/reject", "Tolak offer",
                 responses=responses_ok("OfferResponse"))
    add_endpoint(spec, "POST", f"{R}/offers/{{id}}/withdraw", "Tarik offer",
                 responses=responses_ok("OfferResponse"))

    # Candidate consents
    add_endpoint(spec, "POST", f"{R}/candidates/{{id}}/consents", "Catat persetujuan/penolakan data kandidat",
                 request_body="CreateCandidateConsentRequest", responses=responses_created("CandidateConsentResponse"))
    add_endpoint(spec, "GET", f"{R}/candidates/{{id}}/consents", "Riwayat persetujuan data kandidat",
                 responses=responses_array("CandidateConsentResponse", "List of candidate consents"))

    # Analytics summary (G-11)
    add_endpoint(spec, "GET", f"{R}/analytics/summary", "Ringkasan analitik rekrutmen",
                 responses=responses_ok("RecruitmentAnalyticsSummaryResponse"))

    # Application history & match score
    add_endpoint(spec, "GET", f"{R}/applications/{{id}}/history", "Riwayat perubahan stage aplikasi",
                 responses=responses_array("StageHistoryResponse", "List of stage history entries"))
    add_endpoint(spec, "GET", f"{R}/applications/{{id}}/match-score", "Skor kecocokan kandidat dengan requisition (advisory)",
                 responses=responses_ok("MatchScoreResponse"))

    # Application screenings (G-7)
    add_endpoint(spec, "POST", f"{R}/applications/{{id}}/screenings", "Tambah screening aplikasi",
                 request_body="CreateApplicationScreeningRequest", responses=responses_created("ApplicationScreeningResponse"))
    add_endpoint(spec, "GET", f"{R}/applications/{{id}}/screenings", "Daftar screening aplikasi",
                 responses=responses_array("ApplicationScreeningResponse", "List of application screenings"))
    add_endpoint(spec, "PUT", f"{R}/screenings/{{id}}", "Perbarui screening aplikasi",
                 request_body="UpdateApplicationScreeningRequest", responses=responses_ok("ApplicationScreeningResponse"))
    add_endpoint(spec, "DELETE", f"{R}/screenings/{{id}}", "Hapus screening aplikasi",
                 responses=responses_plain("Screening deleted"))

    # Application assessment (G-12)
    add_endpoint(spec, "GET", f"{R}/applications/{{id}}/assessment", "Detail penilaian kandidat (requirement + assessment)",
                 responses=responses_ok("ApplicationAssessmentDetailResponse"))
    add_endpoint(spec, "PUT", f"{R}/applications/{{id}}/assessment", "Simpan penilaian kandidat (one-per-application)",
                 request_body="SaveApplicationAssessmentRequest", responses=responses_ok("ApplicationAssessmentResponse"))

    # Assessments (G-7 sub-project 2)
    add_endpoint(spec, "POST", f"{R}/assessments", "Buat sesi assessment (batch)",
                 request_body="CreateAssessmentRequest", responses=responses_created("AssessmentResponse"))
    add_endpoint(spec, "GET", f"{R}/assessments", "Daftar sesi assessment (paginated)",
                 responses=responses_list("RecruitmentPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{R}/assessments/{{id}}", "Detail sesi assessment",
                 responses=responses_ok("AssessmentResponse"))
    add_endpoint(spec, "PUT", f"{R}/assessments/{{id}}", "Perbarui sesi assessment",
                 request_body="UpdateAssessmentRequest", responses=responses_ok("AssessmentResponse"))
    add_endpoint(spec, "DELETE", f"{R}/assessments/{{id}}", "Hapus sesi assessment",
                 responses=responses_plain("Assessment deleted"))
    add_endpoint(spec, "POST", f"{R}/assessments/{{id}}/participants", "Tambah peserta assessment",
                 request_body="AddAssessmentParticipantRequest", responses=responses_created("AssessmentParticipantResponse"))
    add_endpoint(spec, "GET", f"{R}/assessments/{{id}}/participants", "Daftar peserta assessment",
                 responses=responses_array("AssessmentParticipantResponse", "List of assessment participants"))
    add_endpoint(spec, "PUT", f"{R}/assessment-participants/{{id}}", "Perbarui peserta assessment (status/skor/hasil)",
                 request_body="UpdateAssessmentParticipantRequest", responses=responses_ok("AssessmentParticipantResponse"))
    add_endpoint(spec, "DELETE", f"{R}/assessment-participants/{{id}}", "Hapus peserta assessment",
                 responses=responses_plain("Participant deleted"))

    # Interviewers (G-8)
    add_endpoint(spec, "POST", f"{R}/interviews/{{id}}/interviewers", "Tambah interviewer ke interview",
                 request_body="AddInterviewerRequest", responses=responses_created("InterviewerResponse"))
    add_endpoint(spec, "GET", f"{R}/interviews/{{id}}/interviewers", "Daftar interviewer interview",
                 responses=responses_array("InterviewerResponse", "List of interviewers"))
    add_endpoint(spec, "DELETE", f"{R}/interviewers/{{id}}", "Hapus interviewer dari interview",
                 responses=responses_plain("Interviewer removed"))

    # Scorecard items (G-8)
    add_endpoint(spec, "POST", f"{R}/interviews/{{id}}/scorecard-items", "Tambah item scorecard interview",
                 request_body="AddScorecardItemRequest", responses=responses_created("ScorecardItemResponse"))
    add_endpoint(spec, "GET", f"{R}/interviews/{{id}}/scorecard-items", "Daftar item scorecard interview",
                 responses=responses_array("ScorecardItemResponse", "List of scorecard items"))
    add_endpoint(spec, "PUT", f"{R}/scorecard-items/{{id}}", "Perbarui item scorecard interview",
                 request_body="UpdateScorecardItemRequest", responses=responses_ok("ScorecardItemResponse"))
    add_endpoint(spec, "DELETE", f"{R}/scorecard-items/{{id}}", "Hapus item scorecard interview",
                 responses=responses_plain("Scorecard item deleted"))

    # Interview complete
    add_endpoint(spec, "POST", f"{R}/interviews/{{id}}/complete", "Tandai interview selesai",
                 responses=responses_ok("InterviewResponse"))

    with open(JSON_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)
    print("OK: recruitment endpoints injected into openapi.json")


if __name__ == "__main__":
    main()
