#!/usr/bin/env python3
"""Inject missing Attendance & Notification endpoints into openapi.json (idempotent).

Documents the 10 endpoints registered in backend routes that were missing from
the OpenAPI spec:

Attendance (Tenant: Time & Attendance):
    GET   /api/v1/tenant/attendance/calendar          -> [SessionResponse]
    GET   /api/v1/tenant/attendance/summary           -> SummaryResponse
    GET   /api/v1/tenant/attendance/reports/sessions  -> [SessionResponse]
    POST  /api/v1/tenant/attendance/corrections       -> CorrectionResponse (201)
    GET   /api/v1/tenant/attendance/corrections       -> CorrectionPaginatedResponse
    GET   /api/v1/tenant/attendance/corrections/{id}  -> CorrectionResponse

Notification (Tenant: Notifications):
    GET   /api/v1/tenant/notifications                -> NotificationPaginatedResponse
    GET   /api/v1/tenant/notifications/unread-count   -> NotificationUnreadCountResponse
    PATCH /api/v1/tenant/notifications/{id}/read      -> GenericSuccessResponse
    POST  /api/v1/tenant/notifications/read-all       -> GenericSuccessResponse

User Accounts (Tenant: User Accounts):
    GET   /api/v1/tenant/user-accounts/me             -> AccountResponse

Usage:
    python scripts/inject_attendance_notification_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]

TAG_ATT = "Tenant: Time & Attendance"
TAG_NOTIF = "Tenant: Notifications"
TAG_UACC = "Tenant: User Accounts"


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


def responses_array(item_ref, desc="List of resources"):
    """Plain non-paginated list: {success, data: [item]}."""
    return {
        "200": {
            "description": desc,
            "content": {
                "application/json": {
                    "schema": {
                        "type": "object",
                        "properties": {
                            "success": {"type": "boolean", "example": True},
                            "data": {"type": "array", "items": ref(item_ref)},
                        },
                    },
                },
            },
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
tags = spec.setdefault("tags", [])
added = 0

A = "/api/v1/tenant/attendance"
N = "/api/v1/tenant/notifications"


def add_path(path, ops):
    global added
    if path not in paths:
        paths[path] = {}
    for method, o in ops.items():
        if method not in paths[path]:
            paths[path][method] = o
            added += 1


def ensure_tag(name):
    if not any(t.get("name") == name for t in tags):
        tags.append({"name": name})


ensure_tag(TAG_NOTIF)

# =========================================================================
# Attendance — Employee Calendar & Summary
# =========================================================================
add_path(f"{A}/calendar", {
    "get": op("getEmployeeAttendanceCalendar", "Get employee attendance calendar",
              "Ambil sesi kerja harian satu karyawan dalam rentang tanggal (untuk tampilan kalender). "
              "Query employee_id, from, dan to wajib diisi.",
              TAG_ATT, parameters=[
                  qparam("employee_id", {"type": "string", "format": "uuid"}),
                  qparam("from", {"type": "string", "format": "date", "example": "2026-01-01"}),
                  qparam("to", {"type": "string", "format": "date", "example": "2026-01-31"}),
              ], responses=responses_array("SessionResponse", "Employee daily work sessions in range")),
})

add_path(f"{A}/summary", {
    "get": op("getEmployeeAttendanceSummary", "Get employee attendance summary",
              "Rekap kehadiran satu karyawan dalam rentang tanggal (present, late, missing check-in/out, "
              "day off, leave, total work & overtime minutes). Query employee_id, from, dan to wajib diisi.",
              TAG_ATT, parameters=[
                  qparam("employee_id", {"type": "string", "format": "uuid"}),
                  qparam("from", {"type": "string", "format": "date", "example": "2026-01-01"}),
                  qparam("to", {"type": "string", "format": "date", "example": "2026-01-31"}),
              ], responses=responses_ok("SummaryResponse")),
})

add_path(f"{A}/reports/sessions", {
    "get": op("getAttendanceReport", "Get attendance report (all employees)",
              "Laporan sesi kehadiran semua karyawan dalam rentang tanggal (non-paginated). "
              "Query from dan to wajib diisi.",
              TAG_ATT, parameters=[
                  qparam("from", {"type": "string", "format": "date", "example": "2026-01-01"}),
                  qparam("to", {"type": "string", "format": "date", "example": "2026-01-31"}),
              ], responses=responses_array("SessionResponse", "Attendance sessions report")),
})

# =========================================================================
# Attendance — Correction Requests
# =========================================================================
add_path(f"{A}/corrections", {
    "post": op("createAttendanceCorrectionRequest", "Create attendance correction request",
               "Ajukan koreksi kehadiran (check-in/check-out salah atau tidak tercatat). "
               "Request diproses lewat alur persetujuan bila flow_id diberikan.",
               TAG_ATT, request_body=content("CreateCorrectionRequest"),
               responses=responses_created("CorrectionResponse")),
    "get": op("listAttendanceCorrectionRequests", "List attendance correction requests",
              "Daftar pengajuan koreksi kehadiran, dapat difilter per karyawan (paginated).",
              TAG_ATT, parameters=[
                  qparam("employee_id", {"type": "string", "format": "uuid"}),
                  qparam("page", {"type": "integer", "default": 1}),
                  qparam("per_page", {"type": "integer", "default": 20}),
              ], responses=responses_ok("CorrectionPaginatedResponse")),
})

add_path(f"{A}/corrections/{{id}}", {
    "get": op("getAttendanceCorrectionRequest", "Get attendance correction request by ID",
              "Ambil detail satu pengajuan koreksi kehadiran berdasarkan ID.",
              TAG_ATT, parameters=[param("id")],
              responses=responses_ok("CorrectionResponse")),
})

# =========================================================================
# Notification — Feed, Unread Count, Read/Read-all
# =========================================================================
add_path(f"{N}", {
    "get": op("listNotifications", "List my notifications",
              "Feed notifikasi pengguna yang sedang login (paginated). Bisa difilter status baca via is_read.",
              TAG_NOTIF, parameters=[
                  qparam("is_read", {"type": "boolean"}),
                  qparam("page", {"type": "integer", "default": 1}),
                  qparam("per_page", {"type": "integer", "default": 20}),
              ], responses=responses_ok("NotificationPaginatedResponse")),
})

add_path(f"{N}/unread-count", {
    "get": op("getNotificationUnreadCount", "Get unread notification count",
              "Jumlah notifikasi belum dibaca untuk badge milik pengguna yang sedang login.",
              TAG_NOTIF, responses=responses_ok("NotificationUnreadCountResponse")),
})

add_path(f"{N}/{{id}}/read", {
    "patch": op("markNotificationAsRead", "Mark a notification as read",
                "Tandai satu notifikasi (milik pengguna yang sedang login) sebagai sudah dibaca.",
                TAG_NOTIF, parameters=[param("id")],
                responses=responses_ok("GenericSuccessResponse")),
})

add_path(f"{N}/read-all", {
    "post": op("markAllNotificationsAsRead", "Mark all notifications as read",
               "Tandai semua notifikasi belum dibaca milik pengguna yang sedang login sebagai sudah dibaca.",
               TAG_NOTIF, responses=responses_ok("GenericSuccessResponse")),
})

# =========================================================================
# User Accounts — Get My Account
# =========================================================================
ensure_tag(TAG_UACC)
add_path("/api/v1/tenant/user-accounts/me", {
    "get": op("getMyUserAccount", "Get my user account",
               "Ambil akun login employee milik user yang sedang login (employee_id, email, status password_set).",
               TAG_UACC, responses=responses_ok("AccountResponse")),
})

# =========================================================================
# Schemas
# =========================================================================
new_schemas = {
    "SummaryResponse": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "from_date": {"type": "string", "format": "date"},
            "to_date": {"type": "string", "format": "date"},
            "total_sessions": {"type": "integer"},
            "present_days": {"type": "integer"},
            "late_days": {"type": "integer"},
            "missing_checkin_days": {"type": "integer"},
            "missing_checkout_days": {"type": "integer"},
            "day_off_days": {"type": "integer"},
            "leave_days": {"type": "number"},
            "total_work_minutes": {"type": "integer"},
            "total_overtime_minutes": {"type": "integer"},
        },
    },
    "CreateCorrectionRequest": {
        "type": "object",
        "required": ["employee_id", "attendance_session_id", "correction_type", "reason"],
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "attendance_session_id": {"type": "string", "format": "uuid"},
            "correction_type": {"type": "string", "example": "checkin_missing"},
            "requested_checkin": {"type": "string", "format": "date-time"},
            "requested_checkout": {"type": "string", "format": "date-time"},
            "reason": {"type": "string", "example": "Lupa check-in saat masuk kantor"},
            "flow_id": {"type": "string", "format": "uuid"},
        },
    },
    "CorrectionResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "attendance_session_id": {"type": "string", "format": "uuid"},
            "correction_type": {"type": "string"},
            "requested_checkin": {"type": "string", "format": "date-time"},
            "requested_checkout": {"type": "string", "format": "date-time"},
            "reason": {"type": "string"},
            "status": {"type": "string", "example": "pending"},
            "approval_instance_id": {"type": "string", "format": "uuid"},
            "approved_at": {"type": "string", "format": "date-time"},
            "created_at": {"type": "string", "format": "date-time"},
        },
    },
    "CorrectionPaginatedResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean", "example": True},
            "data": {"type": "array", "items": ref("CorrectionResponse")},
            "page": {"type": "integer"},
            "per_page": {"type": "integer"},
            "total": {"type": "integer", "format": "int64"},
            "total_pages": {"type": "integer"},
        },
    },
    "NotificationResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "recipient_user_id": {"type": "string", "format": "uuid"},
            "type": {"type": "string", "example": "approval"},
            "title": {"type": "string", "example": "Approval baru menunggu Anda"},
            "body": {"type": "string", "example": "Permintaan cuti Andi menunggu persetujuan Anda."},
            "reference_type": {"type": "string"},
            "reference_id": {"type": "string", "format": "uuid"},
            "is_read": {"type": "boolean"},
            "read_at": {"type": "string", "format": "date-time"},
            "created_at": {"type": "string", "format": "date-time"},
        },
    },
    "NotificationPaginatedResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean", "example": True},
            "data": {
                "type": "object",
                "properties": {
                    "data": {"type": "array", "items": ref("NotificationResponse")},
                    "total": {"type": "integer", "format": "int64"},
                    "page": {"type": "integer"},
                    "per_page": {"type": "integer"},
                },
            },
        },
    },
    "NotificationUnreadCountResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean", "example": True},
            "data": {
                "type": "object",
                "properties": {
                    "unread_count": {"type": "integer", "example": 3},
                },
            },
        },
    },
    "GenericSuccessResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean", "example": True},
            "data": {
                "type": "object",
                "properties": {
                    "success": {"type": "boolean", "example": True},
                },
            },
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
total_tags = len(spec["tags"])
print(f"[OK] Injected/updated {added} items into {JSON_PATH}")
print(f"     Total paths:   {total_paths}")
print(f"     Total schemas: {total_schemas}")
print(f"     Total tags:    {total_tags}")
