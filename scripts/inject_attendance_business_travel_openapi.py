#!/usr/bin/env python3
"""Inject missing attendance endpoints (Business Travels + stats) into openapi.json (idempotent).

Covers endpoints registered in backend/internal/modules/attendance/routes.go that are
missing from the OpenAPI spec:
- GET /attendance/stats/summary, GET /attendance/stats/overtime-trend
- Business Travel CRUD + nested participants/destinations/activities/schedules
- Funding methods + fundings + funding documents
- Expense categories + expenses + expense documents
- Travel documents
- Settlements + settlement submit
- Refunds + refund confirm
- Reimbursements (approve/process/pay)

Usage:
    python scripts/inject_attendance_business_travel_openapi.py
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


# ---------------------------------------------------------------------------
# Schema helpers
# ---------------------------------------------------------------------------

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
        "tags": ["Attendance"],
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
    # Schemas — Business Travel (from dto_businesstravel.go)
    # ------------------------------------------------------------------

    add_schema(spec, "CreateBusinessTravelRequest", obj({
        "title": prop("string", desc="Judul perjalanan dinas"),
        "purpose": prop("string"),
        "description": prop("string"),
        "start_date": prop("string", fmt="date", desc="Tanggal mulai"),
        "end_date": prop("string", fmt="date", desc="Tanggal selesai"),
        "origin": prop("string"),
        "participants": arr(ref("CreateBusinessTravelParticipantRequest")),
        "destinations": arr(ref("CreateBusinessTravelDestinationRequest")),
    }, required=["title", "start_date", "end_date"]))

    add_schema(spec, "UpdateBusinessTravelRequest", obj({
        "title": prop("string"),
        "purpose": prop("string"),
        "description": prop("string"),
        "start_date": prop("string", fmt="date"),
        "end_date": prop("string", fmt="date"),
        "origin": prop("string"),
    }))

    add_schema(spec, "SubmitBusinessTravelRequest", obj({
        "flow_id": prop("string", fmt="uuid", desc="Opsional; auto-resolve dari active flow module business_travel jika kosong"),
    }))

    add_schema(spec, "CreateBusinessTravelParticipantRequest", obj({
        "participant_type": prop("string", desc="EMPLOYEE / EXTERNAL"),
        "employee_id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "organization": prop("string"),
        "position": prop("string"),
        "identity_number": prop("string"),
        "email": prop("string", fmt="email"),
        "phone": prop("string"),
        "role": prop("string"),
        "notes": prop("string"),
    }, required=["participant_type"]))

    add_schema(spec, "CreateBusinessTravelDestinationRequest", obj({
        "sequence": prop("integer"),
        "country": prop("string"),
        "province": prop("string"),
        "city": prop("string"),
        "location": prop("string"),
        "arrival_date": prop("string", fmt="date"),
        "departure_date": prop("string", fmt="date"),
        "purpose": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "CreateBusinessTravelActivityRequest", obj({
        "activity_date": prop("string", fmt="date"),
        "start_time": prop("string"),
        "end_time": prop("string"),
        "title": prop("string"),
        "description": prop("string"),
        "location": prop("string"),
        "organizer": prop("string"),
        "notes": prop("string"),
    }, required=["activity_date", "title"]))

    add_schema(spec, "CreateBusinessTravelScheduleRequest", obj({
        "schedule_type": prop("string", desc="TRANSPORT / HOTEL / MEETING / OTHER"),
        "departure_datetime": prop("string", fmt="date-time"),
        "arrival_datetime": prop("string", fmt="date-time"),
        "origin": prop("string"),
        "destination": prop("string"),
        "transportation_type": prop("string"),
        "provider": prop("string"),
        "booking_reference": prop("string"),
        "notes": prop("string"),
    }, required=["schedule_type"]))

    add_schema(spec, "BusinessTravelResponse", obj({
        "id": prop("string", fmt="uuid"),
        "request_number": prop("string"),
        "requester_id": prop("string", fmt="uuid"),
        "title": prop("string"),
        "purpose": prop("string"),
        "description": prop("string"),
        "start_date": prop("string", fmt="date"),
        "end_date": prop("string", fmt="date"),
        "origin": prop("string"),
        "status": prop("string", enum=["DRAFT", "PENDING_APPROVAL", "APPROVED", "REJECTED", "IN_PROGRESS", "COMPLETED", "CANCELLED"]),
        "approval_status": prop("string"),
        "approval_instance_id": prop("string", fmt="uuid"),
        "participants": arr(ref("BusinessTravelParticipantResponse")),
        "destinations": arr(ref("BusinessTravelDestinationResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "BusinessTravelParticipantResponse", obj({
        "id": prop("string", fmt="uuid"),
        "participant_type": prop("string"),
        "employee_id": prop("string", fmt="uuid"),
        "name": prop("string"),
        "organization": prop("string"),
        "position": prop("string"),
        "identity_number": prop("string"),
        "email": prop("string", fmt="email"),
        "phone": prop("string"),
        "role": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "BusinessTravelDestinationResponse", obj({
        "id": prop("string", fmt="uuid"),
        "sequence": prop("integer"),
        "country": prop("string"),
        "province": prop("string"),
        "city": prop("string"),
        "location": prop("string"),
        "arrival_date": prop("string", fmt="date"),
        "departure_date": prop("string", fmt="date"),
        "purpose": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "BusinessTravelActivityResponse", obj({
        "id": prop("string", fmt="uuid"),
        "activity_date": prop("string", fmt="date"),
        "start_time": prop("string"),
        "end_time": prop("string"),
        "title": prop("string"),
        "description": prop("string"),
        "location": prop("string"),
        "organizer": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "BusinessTravelScheduleResponse", obj({
        "id": prop("string", fmt="uuid"),
        "schedule_type": prop("string"),
        "departure_datetime": prop("string", fmt="date-time"),
        "arrival_datetime": prop("string", fmt="date-time"),
        "origin": prop("string"),
        "destination": prop("string"),
        "transportation_type": prop("string"),
        "provider": prop("string"),
        "booking_reference": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "CreateFundingMethodRequest", obj({
        "name": prop("string"),
        "description": prop("string"),
    }, required=["name"]))

    add_schema(spec, "FundingMethodResponse", obj({
        "id": prop("string", fmt="uuid"),
        "code": prop("string"),
        "name": prop("string"),
        "description": prop("string"),
        "active": prop("boolean"),
    }))

    add_schema(spec, "CreateFundingRequest", obj({
        "funding_method_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "funding_date": prop("string", fmt="date"),
        "payment_method": prop("string"),
        "payment_reference": prop("string"),
        "notes": prop("string"),
    }, required=["funding_method_id", "amount"]))

    add_schema(spec, "UpdateFundingRequest", obj({
        "amount": prop("number"),
        "funding_date": prop("string", fmt="date"),
        "payment_method": prop("string"),
        "payment_reference": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "AddFundingDocumentRequest", obj({
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }, required=["document_type", "file_name", "file_path"]))

    add_schema(spec, "FundingDocumentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "FundingResponse", obj({
        "id": prop("string", fmt="uuid"),
        "business_travel_id": prop("string", fmt="uuid"),
        "funding_method_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "funding_date": prop("string", fmt="date"),
        "payment_method": prop("string"),
        "payment_reference": prop("string"),
        "funded_by": prop("string", fmt="uuid"),
        "status": prop("string", enum=["PENDING", "CONFIRMED", "CANCELLED"]),
        "notes": prop("string"),
        "documents": arr(ref("FundingDocumentResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateExpenseCategoryRequest", obj({
        "name": prop("string"),
        "description": prop("string"),
        "requires_receipt": prop("boolean"),
        "reimbursable": prop("boolean"),
        "payroll_treatment": prop("string"),
        "account_code": prop("string"),
    }, required=["name"]))

    add_schema(spec, "ExpenseCategoryResponse", obj({
        "id": prop("string", fmt="uuid"),
        "code": prop("string"),
        "name": prop("string"),
        "description": prop("string"),
        "requires_receipt": prop("boolean"),
        "reimbursable": prop("boolean"),
        "payroll_treatment": prop("string"),
        "account_code": prop("string"),
        "active": prop("boolean"),
    }))

    add_schema(spec, "CreateExpenseRequest", obj({
        "participant_id": prop("string", fmt="uuid"),
        "expense_category_id": prop("string", fmt="uuid"),
        "expense_date": prop("string", fmt="date"),
        "description": prop("string"),
        "quantity": prop("number"),
        "unit": prop("string"),
        "amount": prop("number"),
        "funding_method_id": prop("string", fmt="uuid"),
        "vendor": prop("string"),
        "receipt_number": prop("string"),
        "notes": prop("string"),
    }, required=["expense_category_id", "expense_date", "amount"]))

    add_schema(spec, "UpdateExpenseRequest", obj({
        "expense_category_id": prop("string", fmt="uuid"),
        "expense_date": prop("string", fmt="date"),
        "description": prop("string"),
        "quantity": prop("number"),
        "unit": prop("string"),
        "amount": prop("number"),
        "funding_method_id": prop("string", fmt="uuid"),
        "vendor": prop("string"),
        "receipt_number": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "AddExpenseDocumentRequest", obj({
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }, required=["document_type", "file_name", "file_path"]))

    add_schema(spec, "ExpenseDocumentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }))

    add_schema(spec, "ExpenseResponse", obj({
        "id": prop("string", fmt="uuid"),
        "business_travel_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "expense_category_id": prop("string", fmt="uuid"),
        "expense_date": prop("string", fmt="date"),
        "description": prop("string"),
        "quantity": prop("number"),
        "unit": prop("string"),
        "amount": prop("number"),
        "funding_method_id": prop("string", fmt="uuid"),
        "vendor": prop("string"),
        "receipt_number": prop("string"),
        "status": prop("string"),
        "notes": prop("string"),
        "documents": arr(ref("ExpenseDocumentResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateSettlementRequest", obj({
        "participant_id": prop("string", fmt="uuid", desc="Kosongkan untuk settlement gabungan seluruh peserta"),
        "notes": prop("string"),
    }))

    add_schema(spec, "SubmitSettlementRequest", obj({
        "flow_id": prop("string", fmt="uuid", desc="Opsional; auto-resolve dari active flow module business_travel_settlement"),
    }))

    add_schema(spec, "SettlementItemResponse", obj({
        "id": prop("string", fmt="uuid"),
        "expense_id": prop("string", fmt="uuid"),
        "funding_method_id": prop("string", fmt="uuid"),
        "item_type": prop("string"),
        "category": prop("string"),
        "amount": prop("number"),
        "notes": prop("string"),
    }))

    add_schema(spec, "SettlementResponse", obj({
        "id": prop("string", fmt="uuid"),
        "business_travel_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "total_advance": prop("number"),
        "total_actual_expense": prop("number"),
        "total_company_paid": prop("number"),
        "total_reimbursement": prop("number"),
        "total_refund": prop("number"),
        "balance": prop("number"),
        "status": prop("string", enum=["DRAFT", "PENDING_APPROVAL", "APPROVED", "REJECTED", "SETTLED"]),
        "approval_instance_id": prop("string", fmt="uuid"),
        "submitted_at": prop("string", fmt="date-time"),
        "approved_at": prop("string", fmt="date-time"),
        "settled_at": prop("string", fmt="date-time"),
        "notes": prop("string"),
        "items": arr(ref("SettlementItemResponse")),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "ConfirmRefundRequest", obj({
        "refund_reference": prop("string"),
        "refund_document": prop("string"),
    }))

    add_schema(spec, "RefundResponse", obj({
        "id": prop("string", fmt="uuid"),
        "business_travel_id": prop("string", fmt="uuid"),
        "settlement_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "refund_amount": prop("number"),
        "refund_date": prop("string", fmt="date"),
        "refund_reference": prop("string"),
        "refunded_by": prop("string", fmt="uuid"),
        "refund_document": prop("string"),
        "status": prop("string", enum=["PENDING", "CONFIRMED", "CANCELLED"]),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "PayReimbursementRequest", obj({
        "payment_reference": prop("string"),
    }))

    add_schema(spec, "TravelReimbursementResponse", obj({
        "id": prop("string", fmt="uuid"),
        "business_travel_id": prop("string", fmt="uuid"),
        "participant_id": prop("string", fmt="uuid"),
        "settlement_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "status": prop("string", enum=["REQUESTED", "APPROVED", "PAID", "REJECTED"]),
        "requested_at": prop("string", fmt="date-time"),
        "approved_at": prop("string", fmt="date-time"),
        "paid_at": prop("string", fmt="date-time"),
        "payment_reference": prop("string"),
        "paid_by": prop("string", fmt="uuid"),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "AddTravelDocumentRequest", obj({
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }, required=["document_type", "file_name", "file_path"]))

    add_schema(spec, "TravelDocumentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "document_type": prop("string"),
        "file_name": prop("string"),
        "file_path": prop("string"),
        "mime_type": prop("string"),
        "file_size": prop("integer", fmt="int64"),
    }))

    # Stats (from dto.go)
    add_schema(spec, "AttendanceStatsResponse", obj({
        "from_date": prop("string", fmt="date"),
        "to_date": prop("string", fmt="date"),
        "total_sessions": prop("integer"),
        "present": prop("integer"),
        "late": prop("integer"),
        "missing_checkin": prop("integer"),
        "missing_checkout": prop("integer"),
        "absent": prop("integer"),
        "leave_days": prop("number"),
        "total_work_minutes": prop("integer"),
        "total_overtime_minutes": prop("integer"),
        "overtime_total": prop("integer"),
        "overtime_pending": prop("integer"),
        "overtime_approved": prop("integer"),
        "overtime_minutes": prop("integer"),
        "travel_total": prop("integer"),
        "travel_approved": prop("integer"),
        "travel_in_progress": prop("integer"),
        "travel_completed": prop("integer"),
    }))

    add_schema(spec, "OvertimeWeek", obj({
        "week_start": prop("string", fmt="date", desc="YYYY-MM-DD (Senin)"),
        "count": prop("integer"),
        "approved": prop("integer"),
        "minutes": prop("integer"),
    }))

    add_schema(spec, "OvertimeTrendResponse", obj({
        "from": prop("string", fmt="date"),
        "to": prop("string", fmt="date"),
        "weeks": arr(ref("OvertimeWeek")),
    }))

    # ------------------------------------------------------------------
    # Endpoints
    # ------------------------------------------------------------------
    BT = "/api/v1/tenant/attendance/business-travels"

    # Stats
    add_endpoint(spec, "GET", "/api/v1/tenant/attendance/stats/summary", "Ringkasan statistik absensi seluruh karyawan (HR dashboard)",
                 responses=responses_ok("AttendanceStatsResponse"),
                 query=[qparam("from_date", {"type": "string", "format": "date"}), qparam("to_date", {"type": "string", "format": "date"})])
    add_endpoint(spec, "GET", "/api/v1/tenant/attendance/stats/overtime-trend", "Tren lembur per minggu (chart HR dashboard)",
                 responses=responses_ok("OvertimeTrendResponse"),
                 query=[qparam("from", {"type": "string", "format": "date"}), qparam("to", {"type": "string", "format": "date"})])

    # Business Travel CRUD
    add_endpoint(spec, "POST", BT, "Buat pengajuan perjalanan dinas baru",
                 request_body="CreateBusinessTravelRequest", responses=responses_created("BusinessTravelResponse"))
    add_endpoint(spec, "GET", BT, "Daftar perjalanan dinas (paginated)",
                 responses=responses_list("AttendancePaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"}),
                        qparam("requester_id", {"type": "string", "format": "uuid"}), qparam("status", {"type": "string"})])
    add_endpoint(spec, "GET", f"{BT}/{{id}}", "Detail perjalanan dinas",
                 responses=responses_ok("BusinessTravelResponse"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}", "Perbarui perjalanan dinas",
                 request_body="UpdateBusinessTravelRequest", responses=responses_ok("BusinessTravelResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/submit", "Submit perjalanan dinas ke Central Approval",
                 request_body="SubmitBusinessTravelRequest", responses=responses_ok("BusinessTravelResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/cancel", "Batalkan perjalanan dinas",
                 responses=responses_ok("BusinessTravelResponse"))

    # Participants
    add_endpoint(spec, "POST", f"{BT}/{{id}}/participants", "Tambah peserta perjalanan dinas",
                 request_body="CreateBusinessTravelParticipantRequest", responses=responses_created("BusinessTravelParticipantResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/participants", "Daftar peserta perjalanan dinas",
                 responses=responses_array("BusinessTravelParticipantResponse", "List of travel participants"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/participants/{{participantId}}", "Perbarui peserta perjalanan dinas",
                 request_body="CreateBusinessTravelParticipantRequest", responses=responses_ok("BusinessTravelParticipantResponse"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/participants/{{participantId}}", "Hapus peserta perjalanan dinas",
                 responses=responses_plain("Participant deleted"))

    # Destinations
    add_endpoint(spec, "POST", f"{BT}/{{id}}/destinations", "Tambah destinasi perjalanan dinas",
                 request_body="CreateBusinessTravelDestinationRequest", responses=responses_created("BusinessTravelDestinationResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/destinations", "Daftar destinasi perjalanan dinas",
                 responses=responses_array("BusinessTravelDestinationResponse", "List of travel destinations"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/destinations/{{destinationId}}", "Perbarui destinasi perjalanan dinas",
                 request_body="CreateBusinessTravelDestinationRequest", responses=responses_ok("BusinessTravelDestinationResponse"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/destinations/{{destinationId}}", "Hapus destinasi perjalanan dinas",
                 responses=responses_plain("Destination deleted"))

    # Activities
    add_endpoint(spec, "POST", f"{BT}/{{id}}/activities", "Tambah aktivitas perjalanan dinas",
                 request_body="CreateBusinessTravelActivityRequest", responses=responses_created("BusinessTravelActivityResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/activities", "Daftar aktivitas perjalanan dinas",
                 responses=responses_array("BusinessTravelActivityResponse", "List of travel activities"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/activities/{{activityId}}", "Perbarui aktivitas perjalanan dinas",
                 request_body="CreateBusinessTravelActivityRequest", responses=responses_ok("BusinessTravelActivityResponse"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/activities/{{activityId}}", "Hapus aktivitas perjalanan dinas",
                 responses=responses_plain("Activity deleted"))

    # Schedules
    add_endpoint(spec, "POST", f"{BT}/{{id}}/schedules", "Tambah jadwal perjalanan dinas",
                 request_body="CreateBusinessTravelScheduleRequest", responses=responses_created("BusinessTravelScheduleResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/schedules", "Daftar jadwal perjalanan dinas",
                 responses=responses_array("BusinessTravelScheduleResponse", "List of travel schedules"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/schedules/{{scheduleId}}", "Perbarui jadwal perjalanan dinas",
                 request_body="CreateBusinessTravelScheduleRequest", responses=responses_ok("BusinessTravelScheduleResponse"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/schedules/{{scheduleId}}", "Hapus jadwal perjalanan dinas",
                 responses=responses_plain("Schedule deleted"))

    # Funding methods & fundings
    add_endpoint(spec, "POST", "/api/v1/tenant/attendance/business-travel-funding-methods", "Tambah metode pendanaan",
                 request_body="CreateFundingMethodRequest", responses=responses_created("FundingMethodResponse"))
    add_endpoint(spec, "GET", "/api/v1/tenant/attendance/business-travel-funding-methods", "Daftar metode pendanaan",
                 responses=responses_array("FundingMethodResponse", "List of funding methods"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/fundings", "Tambah pendanaan perjalanan dinas (setelah travel APPROVED)",
                 request_body="CreateFundingRequest", responses=responses_created("FundingResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/fundings", "Daftar pendanaan perjalanan dinas",
                 responses=responses_array("FundingResponse", "List of fundings"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/fundings/{{fundingId}}", "Perbarui pendanaan perjalanan dinas",
                 request_body="UpdateFundingRequest", responses=responses_ok("FundingResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/fundings/{{fundingId}}/confirm", "Konfirmasi pendanaan",
                 responses=responses_ok("FundingResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/fundings/{{fundingId}}/documents", "Tambah dokumen pendanaan",
                 request_body="AddFundingDocumentRequest", responses=responses_created("FundingDocumentResponse"))

    # Expense categories & expenses
    add_endpoint(spec, "POST", "/api/v1/tenant/attendance/business-travel-expense-categories", "Tambah kategori biaya",
                 request_body="CreateExpenseCategoryRequest", responses=responses_created("ExpenseCategoryResponse"))
    add_endpoint(spec, "GET", "/api/v1/tenant/attendance/business-travel-expense-categories", "Daftar kategori biaya",
                 responses=responses_array("ExpenseCategoryResponse", "List of expense categories"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/expenses", "Tambah biaya aktual perjalanan dinas",
                 request_body="CreateExpenseRequest", responses=responses_created("ExpenseResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/expenses", "Daftar biaya perjalanan dinas",
                 responses=responses_array("ExpenseResponse", "List of expenses"))
    add_endpoint(spec, "PUT", f"{BT}/{{id}}/expenses/{{expenseId}}", "Perbarui biaya perjalanan dinas",
                 request_body="UpdateExpenseRequest", responses=responses_ok("ExpenseResponse"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/expenses/{{expenseId}}", "Hapus biaya perjalanan dinas",
                 responses=responses_plain("Expense deleted"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/expenses/{{expenseId}}/documents", "Tambah dokumen biaya",
                 request_body="AddExpenseDocumentRequest", responses=responses_created("ExpenseDocumentResponse"))

    # Travel documents
    add_endpoint(spec, "POST", f"{BT}/{{id}}/documents", "Tambah dokumen perjalanan (travel order, tiket, hotel, dst)",
                 request_body="AddTravelDocumentRequest", responses=responses_created("TravelDocumentResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/documents", "Daftar dokumen perjalanan dinas",
                 responses=responses_array("TravelDocumentResponse", "List of travel documents"))
    add_endpoint(spec, "DELETE", f"{BT}/{{id}}/documents/{{documentId}}", "Hapus dokumen perjalanan dinas",
                 responses=responses_plain("Document deleted"))

    # Settlements
    add_endpoint(spec, "POST", f"{BT}/{{id}}/settlements", "Buat settlement perjalanan dinas (setelah travel COMPLETED)",
                 request_body="CreateSettlementRequest", responses=responses_created("SettlementResponse"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/settlements", "Daftar settlement perjalanan dinas",
                 responses=responses_array("SettlementResponse", "List of settlements"))
    add_endpoint(spec, "GET", f"{BT}/{{id}}/settlements/{{settlementId}}", "Detail settlement perjalanan dinas",
                 responses=responses_ok("SettlementResponse"))
    add_endpoint(spec, "GET", "/api/v1/tenant/attendance/business-travel-settlements/{settlementId}", "Detail settlement (flat lookup by settlement ID)",
                 responses=responses_ok("SettlementResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/settlements/{{settlementId}}/submit", "Submit settlement ke Central Approval",
                 request_body="SubmitSettlementRequest", responses=responses_ok("SettlementResponse"))

    # Refunds
    add_endpoint(spec, "GET", f"{BT}/{{id}}/refunds", "Daftar refund perjalanan dinas",
                 responses=responses_array("RefundResponse", "List of refunds"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/refunds/{{refundId}}/confirm", "Konfirmasi refund",
                 request_body="ConfirmRefundRequest", responses=responses_ok("RefundResponse"))

    # Reimbursements
    add_endpoint(spec, "GET", f"{BT}/{{id}}/reimbursements", "Daftar reimbursement perjalanan dinas",
                 responses=responses_array("TravelReimbursementResponse", "List of travel reimbursements"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/reimbursements/{{reimbursementId}}/approve", "Setujui reimbursement",
                 responses=responses_ok("TravelReimbursementResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/reimbursements/{{reimbursementId}}/process", "Proses reimbursement (cek subscription modul Reimbursement)",
                 responses=responses_ok("TravelReimbursementResponse"))
    add_endpoint(spec, "POST", f"{BT}/{{id}}/reimbursements/{{reimbursementId}}/pay", "Bayar reimbursement",
                 request_body="PayReimbursementRequest", responses=responses_ok("TravelReimbursementResponse"))

    with open(JSON_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)
    print("OK: attendance business-travel endpoints injected into openapi.json")


if __name__ == "__main__":
    main()
