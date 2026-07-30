#!/usr/bin/env python3
"""Inject insurance settings endpoints into openapi.json."""

import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

added_paths = 0

# =========================================================================
# Schemas
# =========================================================================
schemas = spec.setdefault("components", {}).setdefault("schemas", {})

if "CreateInsuranceRequest" not in schemas:
    schemas["CreateInsuranceRequest"] = {
        "type": "object",
        "required": ["code", "name"],
        "properties": {
            "code": {"type": "string", "maxLength": 20, "description": "Insurance code"},
            "name": {"type": "string", "maxLength": 255, "description": "Insurance type name"},
            "sort_order": {"type": "integer", "description": "Sort order"}
        }
    }

if "UpdateInsuranceRequest" not in schemas:
    schemas["UpdateInsuranceRequest"] = {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 20, "description": "Insurance code"},
            "name": {"type": "string", "maxLength": 255, "description": "Insurance type name"},
            "sort_order": {"type": "integer", "description": "Sort order"}
        }
    }

if "InsuranceResponse" not in schemas:
    schemas["InsuranceResponse"] = {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "code": {"type": "string"},
            "name": {"type": "string"},
            "sort_order": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    }

# =========================================================================
# Endpoints
# =========================================================================
insurance_paths = {
    "/api/v1/tenant/settings/insurances": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "List all Insurances",
            "operationId": "settingInsuranceList",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}
            ],
            "responses": {
                "200": {"description": "Paginated list of Insurance entries", "content": {"application/json": {"schema": {"type": "object", "properties": {
                    "success": {"type": "boolean"},
                    "data": {"type": "array", "items": {"$ref": "#/components/schemas/InsuranceResponse"}},
                    "page": {"type": "integer"},
                    "per_page": {"type": "integer"},
                    "total": {"type": "integer"},
                    "total_pages": {"type": "integer"}
                }}}}},
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a paginated list of Insurance types (reference data for employee insurance)."
        },
        "post": {
            "tags": ["Tenant: Settings"],
            "summary": "Create a new insurance type",
            "operationId": "settingInsuranceCreate",
            "security": [{"bearerAuth": []}],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/CreateInsuranceRequest"}
                    }
                }
            },
            "responses": {
                "201": {"description": "Insurance type created", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/InsuranceResponse"}}}},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"},
                "409": {"description": "Duplicate code"}
            },
            "description": "Create a new insurance type entry (e.g. BPJS Kesehatan, BPJS Ketenagakerjaan, etc.)."
        }
    },
    "/api/v1/tenant/settings/insurances/{id}": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "Get insurance type by ID",
            "operationId": "settingInsuranceGet",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Insurance type details", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/InsuranceResponse"}}}},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Get detailed information about a specific insurance type by its ID."
        },
        "put": {
            "tags": ["Tenant: Settings"],
            "summary": "Update insurance type",
            "operationId": "settingInsuranceUpdate",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/UpdateInsuranceRequest"}
                    }
                }
            },
            "responses": {
                "200": {"description": "Insurance type updated", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/InsuranceResponse"}}}},
                "400": {"description": "Validation error"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"},
                "409": {"description": "Duplicate code"}
            },
            "description": "Update an existing insurance type's code, name, or sort order."
        },
        "delete": {
            "tags": ["Tenant: Settings"],
            "summary": "Delete insurance type",
            "operationId": "settingInsuranceDelete",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Insurance type deleted"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Not found"}
            },
            "description": "Soft-delete an insurance type record."
        }
    }
}

spec["paths"].update(insurance_paths)
added_paths += 5

# =========================================================================
# Write back
# =========================================================================
with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2)

print(f"[OK] Injected {added_paths} insurance endpoints into {JSON_PATH}")
print(f"   - Settings Insurances: GET list + POST create + GET by id + PUT update + DELETE (5)")
print(f"   - Schemas added: CreateInsuranceRequest, UpdateInsuranceRequest, InsuranceResponse")
