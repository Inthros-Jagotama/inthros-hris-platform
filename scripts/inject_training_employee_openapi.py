#!/usr/bin/env python3
"""Inject missing endpoints into openapi.json (idempotent).

Covers endpoints registered in routes.go but missing from the OpenAPI spec:

* Training & Development P0-P2 (88 endpoints): course sub-resources, providers,
  trainers, session trainers/attendance/assessments/costs/documents, plans +
  plan items, needs, requests (+ submit/cancel), mandatories, evaluation
  forms/questions/answers, effectiveness, certifications, history & reports.
* Employee Movement (10): career-history, movement/promotion eligibility,
  movement audits & documents, movement/contract reports, HR dashboard.
* Attendance overtime (4): assignable-employees, assign, actual, cancel.
* Career Intelligence (2): GET/PUT paths/{id}.
* Performance KPI (1): duplicate template.
* Workforce Intelligence (1): candidate-search.

Idempotent: existing operations are left untouched.

Usage:
    python scripts/inject_training_employee_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]

TAG_TRN = "Tenant: Training & Development Management"
TAG_EMOV = "Tenant: Employee Movement & Career Management"
TAG_ATT = "Tenant: Time & Attendance"
TAG_CI = "Tenant: Career Intelligence"
TAG_WFI = "Tenant: Workforce Intelligence & Strategic Planning"
TAG_PERF = "Tenant: Performance Management"


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
        "200": {
            "description": desc,
            "content": {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}, "data": {"type": "array", "items": ref(item_ref)}}}}},
        },
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

T = "/api/v1/tenant/trainings"


def add_path(path, ops):
    """Idempotent: add operations only when the method is not yet documented."""
    global added
    if path not in paths:
        paths[path] = {}
    for method, o in ops.items():
        if method not in paths[path]:
            paths[path][method] = o
            added += 1


def paginated_schema(name, item_schema):
    """Register a Paginated<Name>Response schema if not present; return its name.
    Naming follows the existing convention (e.g. PaginatedTrainingCategoryResponse)."""
    sname = f"Paginated{name}Response"
    if sname not in schemas:
        schemas[sname] = {
            "type": "object",
            "properties": {
                "success": {"type": "boolean"},
                "data": {"type": "array", "items": ref(item_schema)},
                "page": {"type": "integer"},
                "per_page": {"type": "integer"},
                "total": {"type": "integer", "format": "int64"},
                "total_pages": {"type": "integer"},
            },
        }
        global added
        added += 1
    return sname


def add_crud(prefix, oid, label, tag, item_schema, create_req=None, update_req=None,
             desc_create=None, desc_get=None, desc_update=None, desc_delete=None,
             id_param="id"):
    """Standard resource: POST/GET collection + GET/PUT/DELETE item."""
    add_path(prefix, {
        "post": op(f"create{oid}", f"Create {label}",
                   desc_create or f"Buat data {label} baru.",
                   tag, request_body=content(create_req), responses=responses_created(item_schema)),
        "get": op(f"list{oid}", f"List {label}",
                  f"Ambil daftar {label} dengan pagination.",
                  tag, parameters=[qparam("page", {"type": "integer", "example": 1}),
                                   qparam("per_page", {"type": "integer", "example": 20})],
                  responses=responses_list(paginated_schema(oid, item_schema))),
    })
    add_path(f"{prefix}/{{{id_param}}}", {
        "get": op(f"get{oid}", f"Get {label} by ID",
                  desc_get or f"Ambil satu {label} berdasarkan ID.",
                  tag, parameters=[param(id_param)], responses=responses_ok(item_schema)),
        "put": op(f"update{oid}", f"Update {label}",
                  desc_update or f"Perbarui data {label}.",
                  tag, parameters=[param(id_param)], request_body=content(update_req),
                  responses=responses_ok(item_schema)),
        "delete": op(f"delete{oid}", f"Delete {label}",
                     desc_delete or f"Hapus satu {label}.",
                     tag, parameters=[param(id_param)], responses=responses_plain(f"{label} deleted")),
    })


# =========================================================================
# 1. TRAINING & DEVELOPMENT (88 endpoints)
# =========================================================================

# --- Course sub-resources: objectives ---
add_path(f"{T}/courses/{{id}}/objectives", {
    "get": op("listTrainingCourseObjectives", "List course objectives",
              "Ambil daftar objective sebuah course (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("CourseObjectiveResponse", "List of course objectives")),
    "post": op("createTrainingCourseObjective", "Create course objective",
               "Tambah objective (tujuan pembelajaran) pada sebuah course.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateCourseObjectiveRequest"),
               responses=responses_created("CourseObjectiveResponse")),
})
add_path(f"{T}/course-objectives/{{id}}", {
    "put": op("updateTrainingCourseObjective", "Update course objective",
              "Perbarui teks atau urutan objective course.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateCourseObjectiveRequest"),
              responses=responses_ok("CourseObjectiveResponse")),
    "delete": op("deleteTrainingCourseObjective", "Delete course objective",
                 "Hapus satu objective dari course.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Course objective deleted")),
})

# --- Course sub-resources: competencies ---
add_path(f"{T}/courses/{{id}}/competencies", {
    "get": op("listTrainingCourseCompetencies", "List course competencies",
              "Ambil daftar kompetensi yang dipetakan ke course (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("CourseCompetencyResponse", "List of course competencies")),
    "post": op("createTrainingCourseCompetency", "Create course competency",
               "Petakan kompetensi (target level) ke sebuah course.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateCourseCompetencyRequest"),
               responses=responses_created("CourseCompetencyResponse")),
})
add_path(f"{T}/course-competencies/{{id}}", {
    "delete": op("deleteTrainingCourseCompetency", "Delete course competency",
                 "Lepas pemetaan kompetensi dari course.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Course competency deleted")),
})

# --- Course sub-resources: prerequisites ---
add_path(f"{T}/courses/{{id}}/prerequisites", {
    "get": op("listTrainingCoursePrerequisites", "List course prerequisites",
              "Ambil daftar prasyarat sebuah course (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("CoursePrerequisiteResponse", "List of course prerequisites")),
    "post": op("createTrainingCoursePrerequisite", "Create course prerequisite",
               "Tambah prasyarat (COURSE/COMPETENCY/CERTIFICATION/EXPERIENCE) pada course.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateCoursePrerequisiteRequest"),
               responses=responses_created("CoursePrerequisiteResponse")),
})
add_path(f"{T}/course-prerequisites/{{id}}", {
    "delete": op("deleteTrainingCoursePrerequisite", "Delete course prerequisite",
                 "Hapus satu prasyarat dari course.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Course prerequisite deleted")),
})

# --- Providers ---
add_crud(f"{T}/providers", "TrainingProvider", "training provider", TAG_TRN,
         "TrainingProviderResponse", "CreateTrainingProviderRequest", "UpdateTrainingProviderRequest",
         desc_get="Ambil satu provider pelatihan (vendor internal/eksternal) berdasarkan ID.")

# --- Trainers ---
add_crud(f"{T}/trainers", "TrainingTrainer", "training trainer", TAG_TRN,
         "TrainingTrainerResponse", "CreateTrainingTrainerRequest", "UpdateTrainingTrainerRequest",
         desc_get="Ambil satu trainer pelatihan (internal/eksternal) berdasarkan ID.")

# --- Session trainers ---
add_path(f"{T}/sessions/{{id}}/trainers", {
    "get": op("listTrainingSessionTrainers", "List session trainers",
              "Ambil daftar trainer yang ditugaskan pada session (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("TrainingSessionTrainerResponse", "List of session trainers")),
    "post": op("addTrainingSessionTrainer", "Add session trainer",
               "Tugaskan trainer (role MAIN/ASSISTANT) ke session.",
               TAG_TRN, parameters=[param("id")], request_body=content("AddSessionTrainerRequest"),
               responses=responses_created("TrainingSessionTrainerResponse")),
})
add_path(f"{T}/session-trainers/{{id}}", {
    "delete": op("removeTrainingSessionTrainer", "Remove session trainer",
                 "Lepas penugasan trainer dari session.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Session trainer removed")),
})

# --- Session attendance ---
add_path(f"{T}/sessions/{{id}}/attendance", {
    "get": op("listTrainingAttendanceBySession", "List session attendance",
              "Ambil daftar kehadiran peserta dalam session (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("SessionAttendanceRow", "List of session attendance rows")),
    "post": op("markTrainingAttendance", "Mark training attendance",
               "Catat kehadiran peserta (check-in/check-out, status PRESENT/ABSENT/LATE/EXCUSED) pada tanggal tertentu.",
               TAG_TRN, parameters=[param("id")], request_body=content("MarkTrainingAttendanceRequest"),
               responses=responses_created("TrainingAttendanceResponse")),
})
add_path(f"{T}/attendances/{{id}}", {
    "put": op("updateTrainingAttendance", "Update training attendance",
              "Perbarui waktu check-in/check-out atau status kehadiran peserta.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateTrainingAttendanceRequest"),
              responses=responses_ok("TrainingAttendanceResponse")),
})

# --- Session assessments ---
add_path(f"{T}/sessions/{{id}}/assessments", {
    "get": op("listTrainingAssessmentsBySession", "List session assessments",
              "Ambil daftar assessment (pre-test/post-test/final/practical) sebuah session (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("TrainingAssessmentResponse", "List of session assessments")),
    "post": op("createTrainingAssessment", "Create session assessment",
               "Buat assessment baru untuk session (tipe, skor maksimum, passing score, dll).",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateTrainingAssessmentRequest"),
               responses=responses_created("TrainingAssessmentResponse")),
})
add_path(f"{T}/assessments/{{id}}/results", {
    "post": op("submitTrainingAssessmentResult", "Submit assessment result",
               "Input nilai peserta untuk satu assessment (penanda passed dihitung dari passing score).",
               TAG_TRN, parameters=[param("id")], request_body=content("SubmitAssessmentResultRequest"),
               responses=responses_created("TrainingAssessmentResultResponse")),
})

# --- Session costs ---
add_path(f"{T}/sessions/{{id}}/costs", {
    "get": op("listTrainingSessionCosts", "List session costs",
              "Ambil daftar biaya session (TRAINER/PROVIDER/VENUE/MATERIAL/dll) (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("TrainingSessionCostResponse", "List of session costs")),
    "post": op("createTrainingSessionCost", "Create session cost",
               "Catat komponen biaya pada sebuah session pelatihan.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateTrainingSessionCostRequest"),
               responses=responses_created("TrainingSessionCostResponse")),
})
add_path(f"{T}/session-costs/{{id}}", {
    "put": op("updateTrainingSessionCost", "Update session cost",
              "Perbarui tipe/deskripsi/nominal biaya session.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateTrainingSessionCostRequest"),
              responses=responses_ok("TrainingSessionCostResponse")),
    "delete": op("deleteTrainingSessionCost", "Delete session cost",
                 "Hapus satu komponen biaya session.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Session cost deleted")),
})

# --- Session documents ---
add_path(f"{T}/sessions/{{id}}/documents", {
    "get": op("listTrainingDocuments", "List session documents",
              "Ambil daftar dokumen session (proposal, quotation, invoice, dll) (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("TrainingDocumentResponse", "List of session documents")),
    "post": op("createTrainingDocument", "Create session document",
               "Lampirkan dokumen (metadata + file_url) ke sebuah session.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateTrainingDocumentRequest"),
               responses=responses_created("TrainingDocumentResponse")),
})
add_path(f"{T}/documents/{{id}}", {
    "delete": op("deleteTrainingDocument", "Delete session document",
                 "Hapus satu dokumen session.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Document deleted")),
})

# --- Plans + plan items ---
add_crud(f"{T}/plans", "TrainingPlan", "training plan", TAG_TRN,
         "TrainingPlanResponse", "CreateTrainingPlanRequest", "UpdateTrainingPlanRequest",
         desc_get="Ambil satu rencana pelatihan (tahunan) berdasarkan ID.")
add_path(f"{T}/plans/{{id}}/items", {
    "get": op("listTrainingPlanItems", "List plan items",
              "Ambil daftar item (course target) dalam sebuah training plan (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_array("TrainingPlanItemResponse", "List of plan items")),
    "post": op("createTrainingPlanItem", "Create plan item",
               "Tambah item rencana (course, target peserta, estimasi biaya, prioritas) ke plan.",
               TAG_TRN, parameters=[param("id")], request_body=content("CreateTrainingPlanItemRequest"),
               responses=responses_created("TrainingPlanItemResponse")),
})
add_path(f"{T}/plan-items/{{id}}", {
    "put": op("updateTrainingPlanItem", "Update plan item",
              "Perbarui item rencana pelatihan.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateTrainingPlanItemRequest"),
              responses=responses_ok("TrainingPlanItemResponse")),
    "delete": op("deleteTrainingPlanItem", "Delete plan item",
                 "Hapus satu item dari rencana pelatihan.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Plan item deleted")),
})

# --- Needs ---
add_crud(f"{T}/needs", "TrainingNeed", "training need", TAG_TRN,
         "TrainingNeedResponse", "CreateTrainingNeedRequest", "UpdateTrainingNeedRequest",
         desc_get="Ambil satu kebutuhan pelatihan berdasarkan ID.")

# --- Requests (+ submit / cancel) ---
add_path(f"{T}/requests", {
    "post": op("createTrainingRequest", "Create training request",
               "Buat permintaan pelatihan karyawan (diarahkan ke Central Approval).",
               TAG_TRN, request_body=content("CreateTrainingRequestRequest"), responses=responses_created("TrainingRequestResponse")),
    "get": op("listTrainingRequests", "List training requests",
              "Ambil daftar permintaan pelatihan dengan pagination.",
              TAG_TRN, parameters=[qparam("page", {"type": "integer", "example": 1}),
                                   qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list(paginated_schema("TrainingRequest", "TrainingRequestResponse"))),
})
add_path(f"{T}/requests/{{id}}", {
    "get": op("getTrainingRequest", "Get training request by ID",
              "Ambil satu permintaan pelatihan berdasarkan ID.",
              TAG_TRN, parameters=[param("id")], responses=responses_ok("TrainingRequestResponse")),
})
add_path(f"{T}/requests/{{id}}/submit", {
    "post": op("submitTrainingRequest", "Submit training request for approval",
               "Kirim permintaan pelatihan ke alur persetujuan (Central Approval).",
               TAG_TRN, parameters=[param("id")], request_body=content("SubmitTrainingRequestRequest"),
               responses=responses_ok("TrainingRequestResponse")),
})
add_path(f"{T}/requests/{{id}}/cancel", {
    "post": op("cancelTrainingRequest", "Cancel training request",
               "Batalkan permintaan pelatihan (status draft/approved).",
               TAG_TRN, parameters=[param("id")], responses=responses_ok("TrainingRequestResponse")),
})

# --- Mandatories ---
add_crud(f"{T}/mandatories", "TrainingMandatory", "training mandatory", TAG_TRN,
         "TrainingMandatoryResponse", "CreateTrainingMandatoryRequest", "UpdateTrainingMandatoryRequest",
         desc_get="Ambil satu kebijakan pelatihan wajib berdasarkan ID.")

# --- Evaluation forms ---
add_path(f"{T}/evaluation-forms", {
    "get": op("listTrainingEvaluationForms", "List evaluation forms",
              "Ambil daftar form evaluasi pelatihan dengan pagination.",
              TAG_TRN, parameters=[qparam("page", {"type": "integer", "example": 1}),
                                   qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list(paginated_schema("TrainingEvaluationForm", "EvaluationFormResponse"))),
    "post": op("createTrainingEvaluationForm", "Create evaluation form",
               "Buat form evaluasi untuk session (nama + status aktif).",
               TAG_TRN, request_body=content("CreateEvaluationFormRequest"), responses=responses_created("EvaluationFormResponse")),
})
add_path(f"{T}/evaluation-forms/{{form_id}}", {
    "get": op("getTrainingEvaluationForm", "Get evaluation form by ID",
              "Ambil satu form evaluasi berikut daftar pertanyaannya.",
              TAG_TRN, parameters=[param("form_id")], responses=responses_ok("EvaluationFormWithQuestionsResponse")),
    "put": op("updateTrainingEvaluationForm", "Update evaluation form",
              "Perbarui nama/status form evaluasi.",
              TAG_TRN, parameters=[param("form_id")], request_body=content("UpdateEvaluationFormRequest"),
              responses=responses_ok("EvaluationFormResponse")),
    "delete": op("deleteTrainingEvaluationForm", "Delete evaluation form",
                 "Hapus form evaluasi beserta pertanyaannya.",
                 TAG_TRN, parameters=[param("form_id")], responses=responses_plain("Evaluation form deleted")),
})
add_path(f"{T}/sessions/{{id}}/evaluation-form", {
    "get": op("getTrainingEvaluationFormBySession", "Get evaluation form by session",
              "Ambil form evaluasi aktif yang dipakai sebuah session (non-paginated).",
              TAG_TRN, parameters=[param("id")], responses=responses_ok("EvaluationFormWithQuestionsResponse")),
})

# --- Evaluation questions ---
add_path(f"{T}/evaluation-forms/{{form_id}}/questions", {
    "get": op("listTrainingEvaluationQuestions", "List evaluation questions",
              "Ambil daftar pertanyaan sebuah form evaluasi (non-paginated).",
              TAG_TRN, parameters=[param("form_id")], responses=responses_array("EvaluationQuestionResponse", "List of evaluation questions")),
    "post": op("createTrainingEvaluationQuestion", "Create evaluation question",
               "Tambah pertanyaan (RATING/TEXT/SINGLE_CHOICE/MULTIPLE_CHOICE) ke form.",
               TAG_TRN, parameters=[param("form_id")], request_body=content("CreateEvaluationQuestionRequest"),
               responses=responses_created("EvaluationQuestionResponse")),
})
add_path(f"{T}/evaluation-questions/{{id}}", {
    "put": op("updateTrainingEvaluationQuestion", "Update evaluation question",
              "Perbarui pertanyaan, tipe, urutan, atau flag wajib.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateEvaluationQuestionRequest"),
              responses=responses_ok("EvaluationQuestionResponse")),
    "delete": op("deleteTrainingEvaluationQuestion", "Delete evaluation question",
                 "Hapus satu pertanyaan dari form evaluasi.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Evaluation question deleted")),
})

# --- Evaluation answers ---
add_path(f"{T}/evaluation-answers", {
    "get": op("listTrainingEvaluationAnswers", "List evaluation answers",
              "Ambil daftar jawaban evaluasi (filter ?form_id / ?question_id / ?participant_id) — non-paginated.",
              TAG_TRN, parameters=[qparam("form_id"), qparam("question_id"), qparam("participant_id")],
              responses=responses_array("EvaluationAnswerResponse", "List of evaluation answers")),
})
add_path(f"{T}/evaluation-forms/{{form_id}}/participants/{{participant_id}}/answers", {
    "post": op("submitTrainingEvaluationAnswers", "Submit evaluation answers",
               "Kirim jawaban evaluasi peserta untuk satu form evaluasi.",
               TAG_TRN, parameters=[param("form_id"), param("participant_id")],
               request_body=content("SubmitEvaluationAnswersRequest"), responses=responses_created("EvaluationAnswerResponse")),
})

# --- Effectiveness ---
add_path(f"{T}/effectiveness", {
    "post": op("createTrainingEffectivenessAssessment", "Create effectiveness assessment",
               "Catat penilaian efektivitas pelatihan (before/after score) untuk peserta.",
               TAG_TRN, request_body=content("CreateEffectivenessAssessmentRequest"), responses=responses_created("EffectivenessAssessmentResponse")),
    "get": op("listTrainingEffectivenessAssessments", "List effectiveness assessments",
              "Ambil daftar penilaian efektivitas (filter ?participant_id) dengan pagination.",
              TAG_TRN, parameters=[qparam("participant_id"),
                                   qparam("page", {"type": "integer", "example": 1}),
                                   qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list(paginated_schema("TrainingEffectivenessAssessment", "EffectivenessAssessmentResponse"))),
})
add_path(f"{T}/effectiveness/{{id}}", {
    "get": op("getTrainingEffectivenessAssessment", "Get effectiveness assessment by ID",
              "Ambil satu penilaian efektivitas pelatihan (before/after score) berdasarkan ID.",
              TAG_TRN, parameters=[param("id")], responses=responses_ok("EffectivenessAssessmentResponse")),
    "put": op("updateTrainingEffectivenessAssessment", "Update effectiveness assessment",
              "Perbarui penilaian efektivitas pelatihan.",
              TAG_TRN, parameters=[param("id")], request_body=content("UpdateEffectivenessAssessmentRequest"),
              responses=responses_ok("EffectivenessAssessmentResponse")),
    "delete": op("deleteTrainingEffectivenessAssessment", "Delete effectiveness assessment",
                 "Hapus satu penilaian efektivitas.",
                 TAG_TRN, parameters=[param("id")], responses=responses_plain("Effectiveness assessment deleted")),
})

# --- Certifications master ---
add_crud(f"{T}/certifications", "TrainingCertification", "training certification", TAG_TRN,
         "CertificationResponse", "CreateCertificationRequest", "UpdateCertificationRequest",
         desc_get="Ambil satu master sertifikasi (badan penerbit, masa berlaku) berdasarkan ID.")

# --- Generate certificate for participant ---
add_path(f"{T}/participants/{{id}}/certificate", {
    "post": op("generateTrainingCertificate", "Generate participant certificate",
               "Generate/isi sertifikat untuk satu peserta (nomor, file_url, tanggal terbit).",
               TAG_TRN, parameters=[param("id")], request_body=content("GenerateCertificateRequest"),
               responses=responses_created("TrainingCertificateResponse")),
})

# --- History & reports ---
add_path(f"{T}/history", {
    "get": op("getTrainingHistory", "Get training history",
              "Riwayat pelatihan karyawan — filter ?employee_id (non-paginated).",
              TAG_TRN, parameters=[qparam("employee_id")],
              responses=responses_array("TrainingHistoryResponse", "List of training history entries")),
})
add_path(f"{T}/reports/participation", {
    "get": op("getTrainingParticipationReport", "Get participation report",
              "Laporan partisipasi pelatihan (per karyawan) — filter ?session_id / ?course_id.",
              TAG_TRN, parameters=[qparam("session_id"), qparam("course_id")],
              responses=responses_array("TrainingParticipationReportRow", "Participation report rows")),
})
add_path(f"{T}/reports/cost", {
    "get": op("getTrainingCostReport", "Get cost report",
              "Laporan biaya pelatihan per session — filter ?year / ?course_id.",
              TAG_TRN, parameters=[qparam("year", {"type": "integer", "example": 2026}), qparam("course_id")],
              responses=responses_array("TrainingCostReportRow", "Cost report rows")),
})
add_path(f"{T}/reports/compliance", {
    "get": op("getTrainingComplianceReport", "Get compliance report",
              "Laporan kepatuhan pelatihan wajib (mandatory) per karyawan.",
              TAG_TRN, parameters=[qparam("organization_id"), qparam("course_id")],
              responses=responses_array("TrainingComplianceReportRow", "Compliance report rows")),
})
add_path(f"{T}/reports/dashboard", {
    "get": op("getTrainingDashboardReport", "Get training dashboard report",
              "Ringkasan kartu dashboard training: total kursus/sesi/peserta, request approval, completion & pass rate, total biaya, sertifikat terbit.",
              TAG_TRN, responses=responses_ok("TrainingDashboardReport")),
})

# =========================================================================
# 2. EMPLOYEE MOVEMENT (10 endpoints)
# =========================================================================
EM = "/api/v1/tenant/employee-movements"

add_path(f"{EM}/employees/{{employeeId}}/career-history", {
    "get": op("getEmployeeCareerHistory", "Get employee career history",
              "Timeline karier karyawan (JOINED/MOVEMENT/CONTRACT) + posisi saat ini.",
              TAG_EMOV, parameters=[param("employeeId")], responses=responses_ok("CareerHistoryResponse")),
})
add_path(f"{EM}/employees/{{employeeId}}/movement-eligibility", {
    "get": op("getEmployeeMovementEligibility", "Get movement eligibility",
              "Evaluasi kelayakan perpindahan umum: masa kerja, skor performa/KPI/OKR/kompetensi, posisi dalam career path, dan hasil rule.",
              TAG_EMOV, parameters=[param("employeeId")], responses=responses_ok("MovementEligibilityResponse")),
})
add_path(f"{EM}/employees/{{employeeId}}/promotion-eligibility", {
    "get": op("getEmployeePromotionEligibility", "Get promotion eligibility",
              "Evaluasi kelayakan promosi: posisi target berikutnya dalam career path, minimum service months, skor, dan hasil rule.",
              TAG_EMOV, parameters=[param("employeeId")], responses=responses_ok("PromotionEligibilityResponse")),
})
add_path(f"{EM}/movements/{{id}}/audits", {
    "get": op("listMovementAudits", "List movement audits",
              "Audit trail perubahan sebuah movement (create/update/submit/execute/cancel) dengan pagination.",
              TAG_EMOV, parameters=[param("id"),
                                    qparam("page", {"type": "integer", "example": 1}),
                                    qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list(paginated_schema("MovementAudit", "MovementAuditResponse"))),
})
add_path(f"{EM}/movements/{{id}}/documents", {
    "get": op("listMovementDocuments", "List movement documents",
              "Ambil daftar dokumen (SK, surat) sebuah movement dengan pagination.",
              TAG_EMOV, parameters=[param("id"),
                                    qparam("page", {"type": "integer", "example": 1}),
                                    qparam("per_page", {"type": "integer", "example": 20})],
              responses=responses_list(paginated_schema("MovementDocument", "MovementDocumentResponse"))),
    "post": op("createMovementDocument", "Create movement document",
               "Simpan metadata dokumen movement (file fisik di-upload via endpoint upload generik).",
               TAG_EMOV, parameters=[param("id")], request_body=content("CreateMovementDocumentRequest"),
               responses=responses_created("MovementDocumentResponse")),
})
add_path(f"{EM}/movements/{{id}}/documents/{{documentId}}", {
    "delete": op("deleteMovementDocument", "Delete movement document",
                 "Hapus satu dokumen dari movement.",
                 TAG_EMOV, parameters=[param("id"), param("documentId")], responses=responses_plain("Movement document deleted")),
})
add_path(f"{EM}/reports/movements", {
    "get": op("getMovementReport", "Get movement report",
              "Agregasi laporan movement: total per tipe dan per status — filter ?from_date / ?to_date.",
              TAG_EMOV, parameters=[qparam("from_date", {"type": "string", "format": "date"}),
                                    qparam("to_date", {"type": "string", "format": "date"})],
              responses=responses_ok("MovementReportResponse")),
})
add_path(f"{EM}/reports/contracts", {
    "get": op("getContractReport", "Get contract report",
              "Agregasi laporan kontrak: total per status + jumlah kontrak akan berakhir dalam 30 hari.",
              TAG_EMOV, parameters=[qparam("from_date", {"type": "string", "format": "date"}),
                                    qparam("to_date", {"type": "string", "format": "date"})],
              responses=responses_ok("ContractReportResponse")),
})
add_path(f"{EM}/dashboard", {
    "get": op("getMovementHRDashboard", "Get employee movement HR dashboard",
              "Ringkasan kartu HR: movement per tipe, pending approval, effective bulan ini, dan ringkasan kontrak.",
              TAG_EMOV, responses=responses_ok("HRDashboardResponse")),
})

# =========================================================================
# 3. ATTENDANCE — OVERTIME REQUESTS (4 endpoints)
# =========================================================================
ATT = "/api/v1/tenant/attendance"

add_path(f"{ATT}/overtime-requests/assignable-employees", {
    "get": op("listAssignableOvertimeEmployees", "List assignable employees for overtime",
              "Daftar karyawan yang dapat diberi overtime (dua-alur ASSIGNED) — filter ?date (non-paginated).",
              TAG_ATT, parameters=[qparam("date", {"type": "string", "format": "date"})],
              responses=responses_array("AssignableEmployeeResponse", "List of assignable employees")),
})
add_path(f"{ATT}/overtime-requests/assign", {
    "post": op("assignOvertimeRequest", "Assign overtime request",
               "Buat overtime alur ASSIGNED: tetapkan karyawan + jadwal kerja lembur.",
               TAG_ATT, request_body=content("AssignOvertimeRequest"), responses=responses_created("OvertimeResponse")),
})
add_path(f"{ATT}/overtime-requests/{{id}}/actual", {
    "post": op("submitOvertimeActual", "Submit overtime actual",
               "Input jam lembur aktual (start/end aktual + catatan) untuk overtime yang disetujui.",
               TAG_ATT, parameters=[param("id")], request_body=content("SubmitOvertimeActualRequest"),
               responses=responses_ok("OvertimeResponse")),
})
add_path(f"{ATT}/overtime-requests/{{id}}/cancel", {
    "post": op("cancelOvertimeRequest", "Cancel overtime request",
               "Batalkan overtime request (draft/approved).",
               TAG_ATT, parameters=[param("id")], responses=responses_ok("OvertimeResponse")),
})

# =========================================================================
# 4. CAREER INTELLIGENCE — paths/{id} GET & PUT (2 endpoints)
# =========================================================================
CI = "/api/v1/tenant/career-intelligence"
add_path(f"{CI}/paths/{{id}}", {
    "get": op("getCareerPathByID", "Get career path by ID",
              "Ambil satu jalur karier (ladder-style: name + steps[]) berdasarkan ID.",
              TAG_CI, parameters=[param("id")], responses=responses_ok("CareerPathResponse")),
    "put": op("updateCareerPath", "Update career path",
              "Perbarui jalur karier (nama, steps, path type, masa kerja, requirements).",
              TAG_CI, parameters=[param("id")], request_body=content("UpdateCareerPathRequest"),
              responses=responses_ok("CareerPathResponse")),
})

# =========================================================================
# 5. PERFORMANCE — KPI duplicate template (1 endpoint)
# =========================================================================
add_path("/api/v1/tenant/performance/kpi/templates/{id}/duplicate", {
    "post": op("duplicatePerformanceTemplate", "Duplicate KPI template",
               "Salin template KPI beserta seluruh indikator & bobotnya menjadi template baru (draft).",
               TAG_PERF, parameters=[param("id")], responses=responses_created("PerformanceTemplateResponse")),
})

# =========================================================================
# 6. WORKFORCE INTELLIGENCE — candidate search (1 endpoint)
# =========================================================================
add_path("/api/v1/tenant/workforce-intelligence/candidate-search", {
    "get": op("candidateSearch", "Search candidates for open positions",
              "Cari kandidat recruitment untuk posisi kosong — filter ?position_id / ?status / ?query (non-paginated).",
              TAG_WFI, parameters=[qparam("position_id"), qparam("status"), qparam("query")],
              responses=responses_array("CandidateResponse", "List of matching candidates")),
})

# =========================================================================
# Schemas
# =========================================================================
new_schemas = {
    # --- Course sub-resources ---
    "CreateCourseObjectiveRequest": {
        "type": "object", "required": ["objective"],
        "properties": {"objective": {"type": "string"}, "sort_order": {"type": "integer"}},
    },
    "UpdateCourseObjectiveRequest": {
        "type": "object",
        "properties": {"objective": {"type": "string"}, "sort_order": {"type": "integer"}},
    },
    "CourseObjectiveResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "objective": {"type": "string"}, "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "CreateCourseCompetencyRequest": {
        "type": "object", "required": ["competency_id"],
        "properties": {"competency_id": {"type": "string", "format": "uuid"}, "target_level": {"type": "integer"}},
    },
    "CourseCompetencyResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "competency_id": {"type": "string", "format": "uuid"}, "target_level": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "CreateCoursePrerequisiteRequest": {
        "type": "object", "required": ["prerequisite_type"],
        "properties": {
            "prerequisite_type": {"type": "string", "enum": ["COURSE", "COMPETENCY", "CERTIFICATION", "EXPERIENCE"]},
            "prerequisite_id": {"type": "string", "format": "uuid"}, "is_required": {"type": "boolean"},
        },
    },
    "CoursePrerequisiteResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "prerequisite_type": {"type": "string"}, "prerequisite_id": {"type": "string", "format": "uuid"},
            "is_required": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Providers ---
    "CreateTrainingProviderRequest": {
        "type": "object", "required": ["code", "name"],
        "properties": {
            "code": {"type": "string", "maxLength": 20}, "name": {"type": "string", "maxLength": 200},
            "type": {"type": "string", "enum": ["INTERNAL", "EXTERNAL"]},
            "contact_name": {"type": "string"}, "email": {"type": "string", "format": "email"},
            "phone": {"type": "string"}, "address": {"type": "string"}, "website": {"type": "string"},
            "is_active": {"type": "boolean"},
        },
    },
    "UpdateTrainingProviderRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 20}, "name": {"type": "string", "maxLength": 200},
            "type": {"type": "string", "enum": ["INTERNAL", "EXTERNAL"]},
            "contact_name": {"type": "string"}, "email": {"type": "string", "format": "email"},
            "phone": {"type": "string"}, "address": {"type": "string"}, "website": {"type": "string"},
            "is_active": {"type": "boolean"},
        },
    },
    "TrainingProviderResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "code": {"type": "string"}, "name": {"type": "string"},
            "type": {"type": "string"}, "contact_name": {"type": "string"}, "email": {"type": "string"},
            "phone": {"type": "string"}, "address": {"type": "string"}, "website": {"type": "string"},
            "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Trainers ---
    "CreateTrainingTrainerRequest": {
        "type": "object", "required": ["type", "name"],
        "properties": {
            "type": {"type": "string", "enum": ["INTERNAL", "EXTERNAL"]},
            "employee_id": {"type": "string", "format": "uuid"}, "provider_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string", "maxLength": 200}, "email": {"type": "string", "format": "email"},
            "phone": {"type": "string"}, "bio": {"type": "string"}, "is_active": {"type": "boolean"},
        },
    },
    "UpdateTrainingTrainerRequest": {
        "type": "object",
        "properties": {
            "type": {"type": "string", "enum": ["INTERNAL", "EXTERNAL"]},
            "employee_id": {"type": "string", "format": "uuid"}, "provider_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string", "maxLength": 200}, "email": {"type": "string", "format": "email"},
            "phone": {"type": "string"}, "bio": {"type": "string"}, "is_active": {"type": "boolean"},
        },
    },
    "TrainingTrainerResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "type": {"type": "string"},
            "employee_id": {"type": "string", "format": "uuid"}, "provider_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"}, "email": {"type": "string"}, "phone": {"type": "string"},
            "bio": {"type": "string"}, "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Session trainers ---
    "AddSessionTrainerRequest": {
        "type": "object", "required": ["trainer_id"],
        "properties": {"trainer_id": {"type": "string", "format": "uuid"}, "role": {"type": "string", "enum": ["MAIN", "ASSISTANT"]}},
    },
    "TrainingSessionTrainerResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "trainer_id": {"type": "string", "format": "uuid"}, "role": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Attendance ---
    "MarkTrainingAttendanceRequest": {
        "type": "object", "required": ["participant_id", "attendance_date"],
        "properties": {
            "participant_id": {"type": "string", "format": "uuid"},
            "attendance_date": {"type": "string", "format": "date"},
            "check_in": {"type": "string"}, "check_out": {"type": "string"},
            "status": {"type": "string", "enum": ["PRESENT", "ABSENT", "LATE", "EXCUSED"]},
            "remarks": {"type": "string"},
        },
    },
    "UpdateTrainingAttendanceRequest": {
        "type": "object",
        "properties": {
            "check_in": {"type": "string"}, "check_out": {"type": "string"},
            "status": {"type": "string", "enum": ["PRESENT", "ABSENT", "LATE", "EXCUSED"]},
            "remarks": {"type": "string"},
        },
    },
    "TrainingAttendanceResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "participant_id": {"type": "string", "format": "uuid"},
            "attendance_date": {"type": "string", "format": "date"}, "check_in": {"type": "string"},
            "check_out": {"type": "string"}, "status": {"type": "string"}, "remarks": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "SessionAttendanceRow": {
        "type": "object",
        "properties": {
            "attendance_id": {"type": "string", "format": "uuid"}, "participant_id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"}, "attendance_date": {"type": "string", "format": "date"},
            "status": {"type": "string"}, "check_in": {"type": "string"}, "check_out": {"type": "string"},
            "remarks": {"type": "string"},
        },
    },

    # --- Assessments ---
    "CreateTrainingAssessmentRequest": {
        "type": "object", "required": ["session_id", "name"],
        "properties": {
            "session_id": {"type": "string", "format": "uuid"}, "name": {"type": "string", "maxLength": 200},
            "type": {"type": "string", "enum": ["PRE_TEST", "POST_TEST", "FINAL", "PRACTICAL", "OTHER"]},
            "max_score": {"type": "number"}, "passing_score": {"type": "number"},
            "attempt_limit": {"type": "integer"}, "is_required": {"type": "boolean"},
        },
    },
    "TrainingAssessmentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"}, "type": {"type": "string"}, "max_score": {"type": "number"},
            "passing_score": {"type": "number"}, "attempt_limit": {"type": "integer"}, "is_required": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "SubmitAssessmentResultRequest": {
        "type": "object", "required": ["participant_id", "score"],
        "properties": {
            "participant_id": {"type": "string", "format": "uuid"}, "score": {"type": "number"},
            "completed_at": {"type": "string", "format": "date-time"},
        },
    },
    "TrainingAssessmentResultResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "assessment_id": {"type": "string", "format": "uuid"},
            "participant_id": {"type": "string", "format": "uuid"}, "score": {"type": "number"},
            "passed": {"type": "boolean"}, "attempt": {"type": "integer"},
            "completed_at": {"type": "string"}, "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Session costs ---
    "CreateTrainingSessionCostRequest": {
        "type": "object", "required": ["cost_type"],
        "properties": {
            "cost_type": {"type": "string", "enum": ["TRAINER", "PROVIDER", "VENUE", "MATERIAL", "CERTIFICATION", "TRAVEL", "ACCOMMODATION", "OTHER"]},
            "description": {"type": "string"}, "amount": {"type": "number", "minimum": 0},
        },
    },
    "UpdateTrainingSessionCostRequest": {
        "type": "object",
        "properties": {
            "cost_type": {"type": "string", "enum": ["TRAINER", "PROVIDER", "VENUE", "MATERIAL", "CERTIFICATION", "TRAVEL", "ACCOMMODATION", "OTHER"]},
            "description": {"type": "string"}, "amount": {"type": "number", "minimum": 0},
        },
    },
    "TrainingSessionCostResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "cost_type": {"type": "string"}, "description": {"type": "string"}, "amount": {"type": "number"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Session documents ---
    "CreateTrainingDocumentRequest": {
        "type": "object", "required": ["document_type", "file_url"],
        "properties": {
            "document_type": {"type": "string", "enum": ["PROPOSAL", "QUOTATION", "ATTENDANCE_SHEET", "INVOICE", "CONTRACT", "TRAINING_REPORT", "OTHER"]},
            "file_name": {"type": "string", "maxLength": 255}, "file_url": {"type": "string"},
        },
    },
    "TrainingDocumentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "document_type": {"type": "string"}, "file_name": {"type": "string"}, "file_url": {"type": "string"},
            "uploaded_by": {"type": "string", "format": "uuid"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Plans ---
    "CreateTrainingPlanRequest": {
        "type": "object", "required": ["code", "name", "year"],
        "properties": {
            "code": {"type": "string", "maxLength": 30}, "name": {"type": "string", "maxLength": 200},
            "year": {"type": "integer", "minimum": 2000, "maximum": 2100},
            "description": {"type": "string"}, "status": {"type": "string", "enum": ["DRAFT", "ACTIVE", "ARCHIVED"]},
        },
    },
    "UpdateTrainingPlanRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 30}, "name": {"type": "string", "maxLength": 200},
            "year": {"type": "integer", "minimum": 2000, "maximum": 2100},
            "description": {"type": "string"}, "status": {"type": "string", "enum": ["DRAFT", "ACTIVE", "ARCHIVED"]},
        },
    },
    "TrainingPlanResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "code": {"type": "string"}, "name": {"type": "string"},
            "year": {"type": "integer"}, "description": {"type": "string"}, "status": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "CreateTrainingPlanItemRequest": {
        "type": "object", "required": ["course_id"],
        "properties": {
            "course_id": {"type": "string", "format": "uuid"}, "target_date": {"type": "string", "format": "date"},
            "target_participants": {"type": "integer"}, "estimated_cost": {"type": "number"},
            "priority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "URGENT"]},
        },
    },
    "UpdateTrainingPlanItemRequest": {
        "type": "object",
        "properties": {
            "course_id": {"type": "string", "format": "uuid"}, "target_date": {"type": "string", "format": "date"},
            "target_participants": {"type": "integer"}, "estimated_cost": {"type": "number"},
            "priority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "URGENT"]},
        },
    },
    "TrainingPlanItemResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "training_plan_id": {"type": "string", "format": "uuid"},
            "course_id": {"type": "string", "format": "uuid"}, "target_date": {"type": "string"},
            "target_participants": {"type": "integer"}, "estimated_cost": {"type": "number"},
            "priority": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Needs ---
    "CreateTrainingNeedRequest": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "organization_id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "reason": {"type": "string"}, "priority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "URGENT"]},
            "source_type": {"type": "string", "enum": ["MANUAL", "PERFORMANCE", "COMPETENCY", "CAREER", "SUCCESSION", "COMPLIANCE", "WORKFORCE"]},
            "source_id": {"type": "string", "format": "uuid"},
            "status": {"type": "string", "enum": ["OPEN", "PLANNED", "FULFILLED", "CANCELLED"]},
        },
    },
    "UpdateTrainingNeedRequest": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "organization_id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "reason": {"type": "string"}, "priority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "URGENT"]},
            "source_type": {"type": "string", "enum": ["MANUAL", "PERFORMANCE", "COMPETENCY", "CAREER", "SUCCESSION", "COMPLIANCE", "WORKFORCE"]},
            "source_id": {"type": "string", "format": "uuid"},
            "status": {"type": "string", "enum": ["OPEN", "PLANNED", "FULFILLED", "CANCELLED"]},
        },
    },
    "TrainingNeedResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "employee_id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"}, "position_id": {"type": "string", "format": "uuid"},
            "course_id": {"type": "string", "format": "uuid"}, "reason": {"type": "string"},
            "priority": {"type": "string"}, "source_type": {"type": "string"}, "source_id": {"type": "string", "format": "uuid"},
            "status": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Requests ---
    "CreateTrainingRequestRequest": {
        "type": "object", "required": ["employee_id", "course_id", "requested_date"],
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "session_id": {"type": "string", "format": "uuid"}, "requested_date": {"type": "string", "format": "date"},
            "reason": {"type": "string"}, "priority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH", "URGENT"]},
        },
    },
    "SubmitTrainingRequestRequest": {
        "type": "object",
        "properties": {"flow_id": {"type": "string", "format": "uuid", "nullable": True, "description": "ID approval flow; kosong = auto-resolve flow aktif module training_request"}},
    },
    "TrainingRequestResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "employee_id": {"type": "string", "format": "uuid"},
            "course_id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "requested_date": {"type": "string"}, "reason": {"type": "string"}, "priority": {"type": "string"},
            "status": {"type": "string"}, "approval_instance_id": {"type": "string", "format": "uuid"},
            "approved_at": {"type": "string"}, "rejected_at": {"type": "string"}, "supervisor_note": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Mandatories ---
    "CreateTrainingMandatoryRequest": {
        "type": "object", "required": ["course_id"],
        "properties": {
            "course_id": {"type": "string", "format": "uuid"}, "organization_id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"}, "employment_status_id": {"type": "string", "format": "uuid"},
            "due_days": {"type": "integer"}, "validity_period_month": {"type": "integer"}, "is_active": {"type": "boolean"},
        },
    },
    "UpdateTrainingMandatoryRequest": {
        "type": "object",
        "properties": {
            "course_id": {"type": "string", "format": "uuid"}, "organization_id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"}, "employment_status_id": {"type": "string", "format": "uuid"},
            "due_days": {"type": "integer"}, "validity_period_month": {"type": "integer"}, "is_active": {"type": "boolean"},
        },
    },
    "TrainingMandatoryResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "course_id": {"type": "string", "format": "uuid"},
            "organization_id": {"type": "string", "format": "uuid"}, "position_id": {"type": "string", "format": "uuid"},
            "employment_status_id": {"type": "string", "format": "uuid"}, "due_days": {"type": "integer"},
            "validity_period_month": {"type": "integer"}, "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Evaluation forms ---
    "CreateEvaluationFormRequest": {
        "type": "object", "required": ["session_id", "name"],
        "properties": {"session_id": {"type": "string", "format": "uuid"}, "name": {"type": "string", "maxLength": 200}, "is_active": {"type": "boolean"}},
    },
    "UpdateEvaluationFormRequest": {
        "type": "object",
        "properties": {"name": {"type": "string", "maxLength": 200}, "is_active": {"type": "boolean"}},
    },
    "EvaluationFormResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "session_id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"}, "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "EvaluationFormWithQuestionsResponse": {
        "type": "object",
        "properties": {
            "form": {"$ref": "#/components/schemas/EvaluationFormResponse"},
            "questions": {"type": "array", "items": {"$ref": "#/components/schemas/EvaluationQuestionResponse"}},
        },
    },
    "CreateEvaluationQuestionRequest": {
        "type": "object", "required": ["question", "question_type"],
        "properties": {
            "question": {"type": "string"},
            "question_type": {"type": "string", "enum": ["RATING", "TEXT", "SINGLE_CHOICE", "MULTIPLE_CHOICE"]},
            "sort_order": {"type": "integer"}, "is_required": {"type": "boolean"},
        },
    },
    "UpdateEvaluationQuestionRequest": {
        "type": "object",
        "properties": {
            "question": {"type": "string"},
            "question_type": {"type": "string", "enum": ["RATING", "TEXT", "SINGLE_CHOICE", "MULTIPLE_CHOICE"]},
            "sort_order": {"type": "integer"}, "is_required": {"type": "boolean"},
        },
    },
    "EvaluationQuestionResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "form_id": {"type": "string", "format": "uuid"},
            "question": {"type": "string"}, "question_type": {"type": "string"}, "sort_order": {"type": "integer"},
            "is_required": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "SubmitEvaluationAnswersRequest": {
        "type": "object", "required": ["answers"],
        "properties": {
            "answers": {"type": "array", "minItems": 1, "items": {
                "type": "object", "required": ["question_id"],
                "properties": {"question_id": {"type": "string", "format": "uuid"}, "answer": {"type": "string"}},
            }},
        },
    },
    "EvaluationAnswerResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "question_id": {"type": "string", "format": "uuid"},
            "participant_id": {"type": "string", "format": "uuid"}, "answer": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Effectiveness ---
    "CreateEffectivenessAssessmentRequest": {
        "type": "object", "required": ["participant_id", "assessment_date"],
        "properties": {
            "participant_id": {"type": "string", "format": "uuid"}, "assessment_date": {"type": "string", "format": "date"},
            "assessor_employee_id": {"type": "string", "format": "uuid"},
            "before_score": {"type": "number"}, "after_score": {"type": "number"},
            "effectiveness_score": {"type": "number"}, "remarks": {"type": "string"},
        },
    },
    "UpdateEffectivenessAssessmentRequest": {
        "type": "object",
        "properties": {
            "assessment_date": {"type": "string", "format": "date"}, "assessor_employee_id": {"type": "string", "format": "uuid"},
            "before_score": {"type": "number"}, "after_score": {"type": "number"},
            "effectiveness_score": {"type": "number"}, "remarks": {"type": "string"},
        },
    },
    "EffectivenessAssessmentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "participant_id": {"type": "string", "format": "uuid"},
            "assessment_date": {"type": "string"}, "assessor_employee_id": {"type": "string", "format": "uuid"},
            "before_score": {"type": "number"}, "after_score": {"type": "number"},
            "effectiveness_score": {"type": "number"}, "remarks": {"type": "string"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Certifications master ---
    "CreateCertificationRequest": {
        "type": "object", "required": ["code", "name"],
        "properties": {
            "code": {"type": "string", "maxLength": 30}, "name": {"type": "string", "maxLength": 200},
            "issuing_body": {"type": "string"}, "validity_period_month": {"type": "integer"},
            "renewal_required": {"type": "boolean"}, "is_active": {"type": "boolean"},
        },
    },
    "UpdateCertificationRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 30}, "name": {"type": "string", "maxLength": 200},
            "issuing_body": {"type": "string"}, "validity_period_month": {"type": "integer"},
            "renewal_required": {"type": "boolean"}, "is_active": {"type": "boolean"},
        },
    },
    "CertificationResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "code": {"type": "string"}, "name": {"type": "string"},
            "issuing_body": {"type": "string"}, "validity_period_month": {"type": "integer"},
            "renewal_required": {"type": "boolean"}, "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"}, "updated_at": {"type": "string", "format": "date-time"},
        },
    },
    "GenerateCertificateRequest": {
        "type": "object",
        "properties": {
            "certification_id": {"type": "string", "format": "uuid"},
            "certificate_file_url": {"type": "string"}, "expiry_date": {"type": "string", "format": "date"},
        },
    },

    # --- History & reports ---
    "TrainingHistoryResponse": {
        "type": "object",
        "properties": {
            "participant_id": {"type": "string", "format": "uuid"}, "employee_id": {"type": "string", "format": "uuid"},
            "course_id": {"type": "string", "format": "uuid"}, "course_name": {"type": "string"},
            "session_id": {"type": "string", "format": "uuid"}, "session_code": {"type": "string"},
            "start_date": {"type": "string"}, "end_date": {"type": "string"},
            "attendance_status": {"type": "string"}, "score": {"type": "number"},
            "completion_status": {"type": "string"}, "completion_date": {"type": "string"},
            "certificate_no": {"type": "string"}, "certificate_id": {"type": "string", "format": "uuid"},
        },
    },
    "TrainingParticipationReportRow": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_name": {"type": "string"},
            "organization_name": {"type": "string"}, "course_id": {"type": "string", "format": "uuid"},
            "course_name": {"type": "string"}, "session_code": {"type": "string"},
            "session_status": {"type": "string"}, "attendance_status": {"type": "string"},
            "score": {"type": "number"}, "completion_status": {"type": "string"},
        },
    },
    "TrainingCostReportRow": {
        "type": "object",
        "properties": {
            "session_id": {"type": "string", "format": "uuid"}, "session_code": {"type": "string"},
            "course_id": {"type": "string", "format": "uuid"}, "course_name": {"type": "string"},
            "provider_name": {"type": "string"}, "total_cost": {"type": "number"},
            "participant_count": {"type": "integer", "format": "int64"}, "cost_per_participant": {"type": "number"},
        },
    },
    "TrainingComplianceReportRow": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_name": {"type": "string"},
            "organization_name": {"type": "string"}, "course_id": {"type": "string", "format": "uuid"},
            "course_name": {"type": "string"}, "due_date": {"type": "string"},
            "completion_status": {"type": "string"}, "status": {"type": "string"},
        },
    },
    "TrainingDashboardReport": {
        "type": "object",
        "properties": {
            "total_courses": {"type": "integer", "format": "int64"}, "total_sessions": {"type": "integer", "format": "int64"},
            "total_participants": {"type": "integer", "format": "int64"}, "total_providers": {"type": "integer", "format": "int64"},
            "total_requests": {"type": "integer", "format": "int64"}, "approved_requests": {"type": "integer", "format": "int64"},
            "pending_requests": {"type": "integer", "format": "int64"}, "completion_rate": {"type": "number"},
            "pass_rate": {"type": "number"}, "total_training_cost": {"type": "number"},
            "certificates_issued": {"type": "integer", "format": "int64"},
        },
    },

    # --- Employee Movement: career history ---
    "CareerHistoryResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {"$ref": "#/components/schemas/CareerHistoryData"},
        },
    },
    "CareerHistoryData": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_name": {"type": "string"},
            "employee_code": {"type": "string"},
            "current_position": {"$ref": "#/components/schemas/CareerPositionInfo"},
            "timeline": {"type": "array", "items": {"$ref": "#/components/schemas/CareerTimelineEntry"}},
        },
    },
    "CareerPositionInfo": {
        "type": "object",
        "properties": {
            "employment_id": {"type": "string", "format": "uuid"}, "effective_date": {"type": "string"},
            "organization_id": {"type": "string", "format": "uuid"}, "organization_name": {"type": "string"},
            "position_id": {"type": "string", "format": "uuid"}, "position_name": {"type": "string"},
            "employment_status_id": {"type": "string", "format": "uuid"}, "employment_status_name": {"type": "string"},
        },
    },
    "CareerTimelineEntry": {
        "type": "object",
        "properties": {
            "date": {"type": "string"}, "event_type": {"type": "string", "enum": ["JOINED", "MOVEMENT", "CONTRACT"]},
            "title": {"type": "string"}, "description": {"type": "string"},
            "movement_type": {"type": "string"}, "contract_type": {"type": "string"},
            "employment_id": {"type": "string", "format": "uuid"}, "movement_id": {"type": "string", "format": "uuid"},
            "contract_id": {"type": "string", "format": "uuid"},
        },
    },

    # --- Employee Movement: audits & documents ---
    "MovementAuditResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "movement_id": {"type": "string", "format": "uuid"},
            "action": {"type": "string"}, "old_status": {"type": "string"}, "new_status": {"type": "string"},
            "old_data": {"type": "string"}, "new_data": {"type": "string"}, "reason": {"type": "string"},
            "acted_by": {"type": "string", "format": "uuid"}, "acted_at": {"type": "string", "format": "date-time"},
        },
    },
    "CreateMovementDocumentRequest": {
        "type": "object", "required": ["document_type", "file_name", "file_url"],
        "properties": {
            "document_type": {"type": "string", "enum": ["PROMOTION_SK", "MUTATION_SK", "DEMOTION_SK", "RETIREMENT_LETTER", "OFFBOARDING_LETTER", "OTHER"]},
            "file_name": {"type": "string"}, "file_url": {"type": "string", "description": "Path hasil upload generik (diawali /)"},
        },
    },
    "MovementDocumentResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"}, "movement_id": {"type": "string", "format": "uuid"},
            "document_type": {"type": "string"}, "file_name": {"type": "string"}, "file_url": {"type": "string"},
            "uploaded_by": {"type": "string", "format": "uuid"}, "created_at": {"type": "string", "format": "date-time"},
        },
    },

    # --- Employee Movement: eligibility ---
    "EligibilityRuleResult": {
        "type": "object",
        "properties": {
            "code": {"type": "string"}, "label": {"type": "string"}, "met": {"type": "boolean"},
            "actual": {"type": "string"}, "required": {"type": "string"}, "detail": {"type": "string"},
            "available": {"type": "boolean"},
        },
    },
    "MovementEligibilityResponse": {
        "type": "object",
        "properties": {"success": {"type": "boolean"}, "data": {"$ref": "#/components/schemas/MovementEligibilityData"}},
    },
    "MovementEligibilityData": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_name": {"type": "string"},
            "employee_code": {"type": "string"}, "tenure_months": {"type": "integer"},
            "current_position": {"$ref": "#/components/schemas/CareerPositionInfo"},
            "performance_score": {"type": "number"}, "competency_score": {"type": "number"},
            "okr_score": {"type": "number"},
            "career_path_id": {"type": "string", "format": "uuid"}, "career_path_name": {"type": "string"},
            "rules": {"type": "array", "items": {"$ref": "#/components/schemas/EligibilityRuleResult"}},
            "eligible": {"type": "boolean"},
        },
    },
    "PromotionEligibilityResponse": {
        "type": "object",
        "properties": {"success": {"type": "boolean"}, "data": {"$ref": "#/components/schemas/PromotionEligibilityData"}},
    },
    "PromotionEligibilityData": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_name": {"type": "string"},
            "employee_code": {"type": "string"}, "tenure_months": {"type": "integer"},
            "current_position": {"$ref": "#/components/schemas/CareerPositionInfo"},
            "next_position_id": {"type": "string", "format": "uuid"}, "next_position_name": {"type": "string"},
            "next_position_sequence": {"type": "integer"}, "minimum_service_months": {"type": "integer"},
            "performance_score": {"type": "number"}, "competency_score": {"type": "number"},
            "okr_score": {"type": "number"},
            "rules": {"type": "array", "items": {"$ref": "#/components/schemas/EligibilityRuleResult"}},
            "eligible": {"type": "boolean"},
        },
    },

    # --- Employee Movement: reports & dashboard ---
    "MovementReportResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {
                "type": "object",
                "properties": {
                    "total": {"type": "integer", "format": "int64"},
                    "by_type": {"type": "object", "additionalProperties": {"type": "integer", "format": "int64"}},
                    "by_status": {"type": "object", "additionalProperties": {"type": "integer", "format": "int64"}},
                },
            },
        },
    },
    "ContractReportResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {
                "type": "object",
                "properties": {
                    "total": {"type": "integer", "format": "int64"},
                    "by_status": {"type": "object", "additionalProperties": {"type": "integer", "format": "int64"}},
                    "expiring": {"type": "integer", "format": "int64"},
                },
            },
        },
    },
    "HRDashboardResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {
                "type": "object",
                "properties": {
                    "movement_by_type": {"type": "object", "additionalProperties": {"type": "integer", "format": "int64"}},
                    "pending_approval": {"type": "integer", "format": "int64"},
                    "effective_this_month": {"type": "integer", "format": "int64"},
                    "contracts": {
                        "type": "object",
                        "properties": {
                            "active": {"type": "integer", "format": "int64"},
                            "expiring": {"type": "integer", "format": "int64"},
                            "expired": {"type": "integer", "format": "int64"},
                        },
                    },
                },
            },
        },
    },

    # --- Attendance overtime ---
    "AssignableEmployeeResponse": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"}, "employee_code": {"type": "string"},
            "name": {"type": "string"}, "organization_id": {"type": "string", "format": "uuid"},
            "organization_name": {"type": "string"},
        },
    },
    "AssignOvertimeRequest": {
        "type": "object", "required": ["assigned_employee_id", "work_date", "start_time_local", "end_time_local", "requested_minutes"],
        "properties": {
            "assigned_employee_id": {"type": "string", "format": "uuid"},
            "work_date": {"type": "string", "format": "date"},
            "start_time_local": {"type": "string", "format": "date-time"},
            "end_time_local": {"type": "string", "format": "date-time"},
            "requested_minutes": {"type": "integer"}, "reason": {"type": "string"},
        },
    },
    "SubmitOvertimeActualRequest": {
        "type": "object", "required": ["actual_start_time_local", "actual_end_time_local"],
        "properties": {
            "actual_start_time_local": {"type": "string", "format": "date-time"},
            "actual_end_time_local": {"type": "string", "format": "date-time"},
            "actual_note": {"type": "string"}, "attachment_url": {"type": "string"},
        },
    },

    # --- Career Intelligence: update path (ladder-style) ---
    "CreateCareerPathStepRequest": {
        "type": "object", "required": ["position_id", "sequence"],
        "properties": {
            "position_id": {"type": "string", "format": "uuid"},
            "sequence": {"type": "integer", "minimum": 1},
            "minimum_service_months": {"type": "integer", "minimum": 0},
            "requirements": {"type": "string"},
        },
    },
    "CareerPathStepResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"},
            "position_name": {"type": "string"},
            "sequence": {"type": "integer"},
            "minimum_service_months": {"type": "integer"},
            "requirements": {"type": "string"},
            "path_type": {"type": "string"},
            "typical_tenure": {"type": "integer"},
            "competencies": {"type": "string"},
            "certifications": {"type": "string"},
        },
    },
    "UpdateCareerPathRequest": {
        "type": "object", "required": ["steps"],
        "properties": {
            "name": {"type": "string"},
            "description": {"type": "string"},
            "is_active": {"type": "boolean"},
            "steps": {"type": "array", "minItems": 1, "items": {"$ref": "#/components/schemas/CreateCareerPathStepRequest"}},
        },
    },
}

# CareerPathResponse (existing) diperbarui agar mencakup field ladder-style
# (name, description, steps) sesuai response aktual GET/PUT /paths/{id}.
if "CareerPathResponse" in schemas:
    cp = schemas["CareerPathResponse"]
    props = cp.setdefault("properties", {})
    if "name" not in props:
        props["name"] = {"type": "string"}
        added += 1
    if "description" not in props:
        props["description"] = {"type": "string"}
        added += 1
    if "steps" not in props:
        props["steps"] = {"type": "array", "items": {"$ref": "#/components/schemas/CareerPathStepResponse"}}
        added += 1

for name, schema in new_schemas.items():
    if name not in schemas:
        schemas[name] = schema
        added += 1

# Bersihkan schema paginated bernama lama (suffix) dari run sebelumnya —
# konvensi sekarang prefix: Paginated<Name>Response (lihat paginated_schema).
legacy_paginated = [
    "TrainingProviderPaginatedResponse", "TrainingTrainerPaginatedResponse",
    "TrainingPlanPaginatedResponse", "TrainingNeedPaginatedResponse",
    "TrainingRequestPaginatedResponse", "TrainingMandatoryPaginatedResponse",
    "TrainingEvaluationFormPaginatedResponse", "TrainingEvaluationAnswerPaginatedResponse",
    "TrainingEffectivenessAssessmentPaginatedResponse", "TrainingCertificationPaginatedResponse",
    "MovementAuditPaginatedResponse", "MovementDocumentPaginatedResponse",
]
for _s in legacy_paginated:
    if _s in schemas:
        del schemas[_s]
        added += 1
        print(f"  removed legacy schema {_s}")

# =========================================================================
# Force-fixup — idempotent add tidak menimpa operasi/schema yang sudah ada,
# jadi perbaikan berikut ditulis langsung (overwrite).
# =========================================================================

# 1) UpdateCareerPathRequest: versi lama non-ladder ditimpa dengan versi
#    ladder-style (name + steps) — hanya bila konten berbeda (idempotent).
if "UpdateCareerPathRequest" in new_schemas:
    if schemas.get("UpdateCareerPathRequest") != new_schemas["UpdateCareerPathRequest"]:
        schemas["UpdateCareerPathRequest"] = new_schemas["UpdateCareerPathRequest"]
        added += 1
        print("  replaced schema UpdateCareerPathRequest (ladder-style)")

# 2) GET operations yang mereferensikan schema paginated bernama lama
#    (suffix) diarahkan ulang ke nama baru (prefix Paginated<Name>Response).
_get_fixups = {
    f"{T}/providers": "PaginatedTrainingProviderResponse",
    f"{T}/trainers": "PaginatedTrainingTrainerResponse",
    f"{T}/plans": "PaginatedTrainingPlanResponse",
    f"{T}/needs": "PaginatedTrainingNeedResponse",
    f"{T}/requests": "PaginatedTrainingRequestResponse",
    f"{T}/mandatories": "PaginatedTrainingMandatoryResponse",
    f"{T}/evaluation-forms": "PaginatedTrainingEvaluationFormResponse",
    f"{T}/effectiveness": "PaginatedTrainingEffectivenessAssessmentResponse",
    f"{T}/certifications": "PaginatedTrainingCertificationResponse",
    f"{EM}/movements/{{id}}/audits": "PaginatedMovementAuditResponse",
    f"{EM}/movements/{{id}}/documents": "PaginatedMovementDocumentResponse",
}
for _path, _sname in _get_fixups.items():
    _entry = paths.get(_path)
    if _entry and "get" in _entry:
        _new = responses_list(_sname)
        if _entry["get"].get("responses") != _new:
            _entry["get"]["responses"] = _new
            added += 1
            print(f"  fixed GET {_path} -> {_sname}")

# 3) evaluation-answers: non-paginated (handler tanpa parsePagination).
_path = f"{T}/evaluation-answers"
if _path in paths and "get" in paths[_path]:
    _new = responses_array("EvaluationAnswerResponse", "List of evaluation answers")
    if paths[_path]["get"].get("responses") != _new:
        paths[_path]["get"]["responses"] = _new
        paths[_path]["get"]["parameters"] = [qparam("form_id"), qparam("question_id"), qparam("participant_id")]
        added += 1
        print("  fixed GET evaluation-answers -> non-paginated array")

# =========================================================================
# Tags (all referenced tags already exist; add only if missing)
# =========================================================================
existing_tags = {t["name"] for t in spec.get("tags", [])}
for tag in (TAG_TRN, TAG_EMOV, TAG_ATT, TAG_CI, TAG_WFI, TAG_PERF):
    if tag not in existing_tags:
        spec.setdefault("tags", []).append({"name": tag, "description": tag})
        added += 1

with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

total_paths = len(spec["paths"])
total_schemas = len(spec["components"]["schemas"])
print(f"[OK] Injected/updated {added} items into {JSON_PATH}")
print(f"     Total paths:   {total_paths}")
print(f"     Total schemas: {total_schemas}")
