#!/usr/bin/env python3
"""Inject missing Leave Calendar & Usage Report endpoints into openapi.json (idempotent).

Documents the 2 endpoints registered in backend routes that were missing from
the OpenAPI spec:

Leave (Tenant: Leave & Time Off):
    GET /api/v1/tenant/leave/calendar        -> [CalendarEntryResponse]
    GET /api/v1/tenant/leave/reports/usage   -> [LeaveRequestResponse]

Query parameters follow the handlers (backend/internal/modules/leave/handler.go):
    calendar      : employee_id, from, to (all required)
    reports/usage : from, to (both required)

Usage:
    python scripts/inject_leave_calendar_report_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]

TAG_LEAVE = "Tenant: Leave & Time Off"


def param(name, where="path", required=True, schema=None):
    return {"name": name, "in": where, "required": required,
            "schema": schema or {"type": "string", "format": "uuid"}}


def qparam(name, schema=None):
    return param(name, where="query", required=False, schema=schema or {"type": "string"})


def ref(name):
    return {"$ref": f"#/components/schemas/{name}"}


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


def op(oid, summary, description, parameters=None, responses=None):
    return {
        "operationId": oid,
        "summary": summary,
        "description": description,
        "tags": [TAG_LEAVE],
        "security": list(SEC),
        "parameters": parameters or [],
        "responses": responses or responses_plain(),
    }


def responses_plain(desc_200="OK"):
    return {
        "200": {"description": desc_200},
        "400": {"description": "Bad request / Validation error"},
    }


with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

paths = spec["paths"]
schemas = spec["components"]["schemas"]
tags = spec.setdefault("tags", [])
added = 0

L = "/api/v1/tenant/leave"


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


ensure_tag(TAG_LEAVE)

# =========================================================================
# Leave — Employee Calendar
# =========================================================================
add_path(f"{L}/calendar", {
    "get": op("getEmployeeLeaveCalendar", "Get employee leave calendar",
              "Entri cuti harian satu karyawan dalam rentang tanggal (untuk tampilan kalender). "
              "Query employee_id, from, dan to wajib diisi.",
              parameters=[
                  qparam("employee_id", {"type": "string", "format": "uuid"}),
                  qparam("from", {"type": "string", "format": "date", "example": "2026-01-01"}),
                  qparam("to", {"type": "string", "format": "date", "example": "2026-01-31"}),
              ], responses=responses_array("CalendarEntryResponse", "Employee leave calendar entries in range")),
})

# =========================================================================
# Leave — Usage Report
# =========================================================================
add_path(f"{L}/reports/usage", {
    "get": op("getLeaveUsageReport", "Get leave usage report (all employees)",
              "Semua permintaan cuti karyawan yang rentang tanggalnya beririsan dengan [from, to] — "
              "bentuk item sama dengan ListLeaveRequests, non-paginated. Data sumber untuk laporan "
              "penggunaan cuti. Query from dan to wajib diisi.",
              parameters=[
                  qparam("from", {"type": "string", "format": "date", "example": "2026-01-01"}),
                  qparam("to", {"type": "string", "format": "date", "example": "2026-01-31"}),
              ], responses=responses_array("LeaveRequestResponse", "Leave usage report (all employees)")),
})

# =========================================================================
# Schemas
# =========================================================================
new_schemas = {
    "CalendarEntryResponse": {
        "type": "object",
        "properties": {
            "leave_request_id": {"type": "string", "format": "uuid"},
            "leave_date": {"type": "string", "format": "date", "example": "2026-01-05"},
            "day_fraction": {"type": "number", "format": "float", "example": 1.0},
            "leave_type_id": {"type": "string", "format": "uuid"},
            "status": {"type": "string", "example": "approved"},
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
