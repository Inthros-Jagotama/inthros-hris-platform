"""Inject education-majors settings endpoints + schemas into openapi.json (idempotent).

Adds:
- Paths: /api/v1/tenant/settings/education-majors (GET list + POST create),
         /api/v1/tenant/settings/education-majors/{id} (GET + PUT + DELETE)
- Schemas: EducationMajorResponse, EducationMajorPaginatedResponse,
           CreateEducationMajorRequest, UpdateEducationMajorRequest
- Tag: Tenant: Settings (adds education majors mention)

Usage: python scripts/inject_education_majors_openapi.py
"""

import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OPENAPI_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

TAG = "Tenant: Settings"
BASE = "/api/v1/tenant/settings/education-majors"
PK = f"{BASE}/{{id}}"

RESP_SCHEMA = {
    "type": "object",
    "properties": {
        "id": {"type": "string", "format": "uuid"},
        "code": {"type": "string"},
        "name": {"type": "string"},
        "sort_order": {"type": "integer"},
        "created_at": {"type": "string", "format": "date-time"},
        "updated_at": {"type": "string", "format": "date-time"},
    },
}
PAG_SCHEMA = {
    "type": "object",
    "properties": {
        "success": {"type": "boolean"},
        "data": {"type": "array", "items": {"$ref": "#/components/schemas/EducationMajorResponse"}},
        "page": {"type": "integer"},
        "per_page": {"type": "integer"},
        "total": {"type": "integer"},
        "total_pages": {"type": "integer"},
    },
}
CREATE_SCHEMA = {
    "type": "object",
    "required": ["code", "name"],
    "properties": {
        "code": {"type": "string", "maxLength": 20},
        "name": {"type": "string", "maxLength": 255},
        "sort_order": {"type": "integer"},
    },
}
UPDATE_SCHEMA = {
    "type": "object",
    "properties": {
        "code": {"type": "string", "maxLength": 20},
        "name": {"type": "string", "maxLength": 255},
        "sort_order": {"type": "integer"},
    },
}


def main():
    with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
        spec = json.load(f)

    added = 0

    # ── Schemas ──
    schemas = spec.setdefault("components", {}).setdefault("schemas", {})
    for name, body in {
        "EducationMajorResponse": RESP_SCHEMA,
        "EducationMajorPaginatedResponse": PAG_SCHEMA,
        "CreateEducationMajorRequest": CREATE_SCHEMA,
        "UpdateEducationMajorRequest": UPDATE_SCHEMA,
    }.items():
        if name not in schemas:
            schemas[name] = body
            added += 1

    # ── Paths ──
    paths = spec["paths"]
    list_params = [
        {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
        {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}},
    ]
    id_param = [{"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}]
    auth = [{"bearerAuth": []}]

    if BASE not in paths:
        paths[BASE] = {
            "get": {
                "tags": [TAG],
                "summary": "List all Education Majors",
                "operationId": "settingEducationMajorList",
                "parameters": list_params,
                "security": auth,
                "responses": {
                    "200": {
                        "description": "Paginated list of Education Majors",
                        "content": {"application/json": {"schema": {"$ref": "#/components/schemas/EducationMajorPaginatedResponse"}}},
                    },
                    "401": {"description": "Unauthorized: missing or invalid JWT"},
                },
                "description": "Retrieve a paginated list of education majors (fields of study). Supports pagination parameters.",
            },
            "post": {
                "tags": [TAG],
                "summary": "Create a new education major",
                "operationId": "settingEducationMajorCreate",
                "security": auth,
                "requestBody": {
                    "required": True,
                    "content": {"application/json": {"schema": {"$ref": "#/components/schemas/CreateEducationMajorRequest"}}},
                },
                "responses": {
                    "201": {
                        "description": "Education major created successfully",
                        "content": {"application/json": {"schema": {"$ref": "#/components/schemas/EducationMajorResponse"}}},
                    },
                    "400": {"description": "Validation error"},
                    "401": {"description": "Unauthorized: missing or invalid JWT"},
                },
                "description": "Create a new education major record. Validates required fields and returns the created resource with its assigned ID.",
            },
        }
        added += 2

    if PK not in paths:
        paths[PK] = {
            "get": {
                "tags": [TAG],
                "summary": "Get education major by ID",
                "operationId": "settingEducationMajorGet",
                "parameters": id_param,
                "security": auth,
                "responses": {
                    "200": {
                        "description": "Education major details",
                        "content": {"application/json": {"schema": {"$ref": "#/components/schemas/EducationMajorResponse"}}},
                    },
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Education major not found"},
                },
                "description": "Get detailed information about a specific education major by its ID.",
            },
            "put": {
                "tags": [TAG],
                "summary": "Update education major",
                "operationId": "settingEducationMajorUpdate",
                "parameters": id_param,
                "security": auth,
                "requestBody": {
                    "required": True,
                    "content": {"application/json": {"schema": {"$ref": "#/components/schemas/UpdateEducationMajorRequest"}}},
                },
                "responses": {
                    "200": {
                        "description": "Education major updated",
                        "content": {"application/json": {"schema": {"$ref": "#/components/schemas/EducationMajorResponse"}}},
                    },
                    "400": {"description": "Validation error"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Education major not found"},
                },
                "description": "Update an education major record's details including code, name, and sort order.",
            },
            "delete": {
                "tags": [TAG],
                "summary": "Delete education major",
                "operationId": "settingEducationMajorDelete",
                "parameters": id_param,
                "security": auth,
                "responses": {
                    "200": {"description": "Education major deleted"},
                    "401": {"description": "Unauthorized"},
                    "404": {"description": "Education major not found"},
                },
                "description": "Soft-delete an education major record. Sets the deleted_at timestamp and hides it from standard queries.",
            },
        }
        added += 3

    # ── Tag description mention ──
    for tag in spec.get("tags", []):
        if tag.get("name") == TAG and "education majors" not in tag.get("description", ""):
            tag["description"] = tag["description"].rstrip(".") + ", education majors."
            break

    with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)

    print(f"[OK] Injected {added} education-majors items into {OPENAPI_PATH}")
    print("   - Paths: GET list + POST create, GET/PUT/DELETE by id (5)")


if __name__ == "__main__":
    main()
