#!/usr/bin/env python3
"""Inject missing endpoints into openapi.json — photo, documents/upload, gradings, villages cascading, user password."""

import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

added_paths = 0
added_schemas = 0

# =========================================================================
# 1. Employee — Photo Upload/Delete
# =========================================================================
photo_paths = {
    "/api/v1/tenant/employees/{id}/photo": {
        "put": {
            "tags": ["Tenant: Employees"],
            "summary": "Upload employee profile photo",
            "operationId": "uploadEmployeePhoto",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "multipart/form-data": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "photo": {"type": "string", "format": "binary", "description": "Photo file (JPG, PNG, GIF, WebP, max 2MB)"}
                            },
                            "required": ["photo"]
                        }
                    }
                }
            },
            "responses": {
                "200": {"description": "Photo uploaded successfully"},
                "400": {"description": "Validation error — invalid file type or size"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Upload a profile photo for an employee. Accepts JPG, PNG, GIF, WebP files up to 2MB. The photo is stored on the server and the employee's profile_picture field is updated."
        },
        "delete": {
            "tags": ["Tenant: Employees"],
            "summary": "Delete employee profile photo",
            "operationId": "deleteEmployeePhoto",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Photo deleted successfully"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Employee not found"}
            },
            "description": "Remove the profile photo from an employee. Deletes the file from the server and clears the profile_picture field."
        }
    },
    "/api/v1/tenant/employees/{id}/documents/upload": {
        "post": {
            "tags": ["Tenant: Employees"],
            "summary": "Upload a document file for employee",
            "operationId": "uploadEmployeeDocument",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "multipart/form-data": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string", "description": "Document name"},
                                "file": {"type": "string", "format": "binary", "description": "Document file (PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, GIF, TXT, max 10MB)"},
                                "note": {"type": "string", "description": "Optional note"}
                            },
                            "required": ["name", "file"]
                        }
                    }
                }
            },
            "responses": {
                "201": {"description": "Document created"},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Upload a document file for an employee. Creates a document record with the uploaded file path. Supports PDF, DOC, DOCX, XLS, XLSX, JPG, PNG, GIF, TXT files up to 10MB."
        }
    },
    "/api/v1/tenant/employees/{id}/documents/{documentId}/upload": {
        "put": {
            "tags": ["Tenant: Employees"],
            "summary": "Replace document file",
            "operationId": "replaceEmployeeDocument",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}},
                {"name": "documentId", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "content": {
                    "multipart/form-data": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "name": {"type": "string", "description": "Document name"},
                                "file": {"type": "string", "format": "binary", "description": "New document file (optional — if omitted, only metadata is updated)"},
                                "note": {"type": "string", "description": "Optional note"}
                            }
                        }
                    }
                }
            },
            "responses": {
                "200": {"description": "Document updated"},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Document not found"}
            },
            "description": "Replace a document file or update document metadata. If a file is provided, it replaces the existing file on disk. If no file is provided, only name/note fields are updated via JSON."
        }
    }
}
spec["paths"].update(photo_paths)
added_paths += 3

# =========================================================================
# 2. Settings — Gradings CRUD + Village Search/Detail + Cascading
# =========================================================================
settings_extra_paths = {
    # Gradings CRUD
    "/api/v1/tenant/settings/gradings": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "List all Gradings",
            "operationId": "settingGradingList",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}
            ],
            "responses": {
                "200": {"description": "Paginated list of Gradings"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a paginated list of Gradings (job grading levels)."
        },
        "post": {
            "tags": ["Tenant: Settings"],
            "summary": "Create a new grading",
            "operationId": "settingGradingCreate",
            "security": [{"bearerAuth": []}],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "required": ["code", "name"],
                            "properties": {
                                "code": {"type": "string"},
                                "name": {"type": "string"},
                                "sort_order": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "responses": {
                "201": {"description": "Grading created"},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Create a new grading record (job grading level)."
        }
    },
    "/api/v1/tenant/settings/gradings/{id}": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "Get grading by ID",
            "operationId": "settingGradingGet",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Grading details"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Get detailed information about a specific grading by its ID."
        },
        "put": {
            "tags": ["Tenant: Settings"],
            "summary": "Update grading",
            "operationId": "settingGradingUpdate",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "code": {"type": "string"},
                                "name": {"type": "string"},
                                "sort_order": {"type": "integer"}
                            }
                        }
                    }
                }
            },
            "responses": {
                "200": {"description": "Grading updated"},
                "400": {"description": "Validation error"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Update a grading record's code, name, or sort order."
        },
        "delete": {
            "tags": ["Tenant: Settings"],
            "summary": "Delete grading",
            "operationId": "settingGradingDelete",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Grading deleted"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Not found"}
            },
            "description": "Soft-delete a grading record."
        }
    },
    # Villages search & detail
    "/api/v1/tenant/settings/villages/search": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "Search villages by name",
            "operationId": "searchVillages",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "q", "in": "query", "required": True, "schema": {"type": "string"}, "description": "Search query for village name"}
            ],
            "responses": {
                "200": {"description": "Search results with hierarchy"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Search villages by name for autocomplete. Returns matching villages with their province, regency, and district hierarchy names."
        }
    },
    "/api/v1/tenant/settings/villages/{id}/detail": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "Get village detail with hierarchy",
            "operationId": "getVillageDetail",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}, "description": "Village code (10 chars)"}
            ],
            "responses": {
                "200": {"description": "Village with district/regency/province names"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Get a village by its code including the full hierarchy: village name, district name, regency name, and province name."
        }
    },
    # Cascading: provinces/{id}/regencies
    "/api/v1/tenant/settings/provinces/{id}/regencies": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "List regencies by province",
            "operationId": "listRegenciesByProvince",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}, "description": "Province ID (code)"},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 200}}
            ],
            "responses": {
                "200": {"description": "List of regencies"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a list of regencies/cities within a specific province for cascading dropdown selection."
        }
    },
    # Cascading: regencies/{id}/districts
    "/api/v1/tenant/settings/regencies/{id}/districts": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "List districts by regency",
            "operationId": "listDistrictsByRegency",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}, "description": "Regency ID (code)"},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 200}}
            ],
            "responses": {
                "200": {"description": "List of districts"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a list of districts within a specific regency for cascading dropdown selection."
        }
    },
    # Cascading: districts/{id}/villages
    "/api/v1/tenant/settings/districts/{id}/villages": {
        "get": {
            "tags": ["Tenant: Settings"],
            "summary": "List villages by district",
            "operationId": "listVillagesByDistrict",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}, "description": "District ID (code)"},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 200}}
            ],
            "responses": {
                "200": {"description": "List of villages"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a list of villages within a specific district for cascading dropdown selection."
        }
    }
}
spec["paths"].update(settings_extra_paths)
added_paths += 10

# =========================================================================
# 3. Platform — User Password Change
# =========================================================================
user_extra_paths = {
    "/api/v1/platform/users/{id}/password": {
        "put": {
            "tags": ["Platform: Users"],
            "summary": "Change user password",
            "operationId": "changeUserPassword",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "required": ["current_password", "new_password", "confirm_password"],
                            "properties": {
                                "current_password": {"type": "string", "description": "Current password"},
                                "new_password": {"type": "string", "minLength": 6, "description": "New password (min 6 characters)"},
                                "confirm_password": {"type": "string", "description": "Confirm new password"}
                            }
                        }
                    }
                }
            },
            "responses": {
                "200": {"description": "Password changed successfully"},
                "400": {"description": "Validation error — current password wrong or passwords don't match"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "User not found"}
            },
            "description": "Change a user's password. Requires current password verification, new password (min 6 characters), and confirmation."
        }
    },
    "/api/v1/platform/users/{id}": {
        "delete": {
            "tags": ["Platform: Users"],
            "summary": "Delete platform user",
            "operationId": "deletePlatformUser",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "User soft-deleted"},
                "403": {"description": "Forbidden: super_admin cannot be deleted"},
                "404": {"description": "User not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Soft-delete a platform user account. Super admin users cannot be deleted."
        }
    }
}
spec["paths"].update(user_extra_paths)
added_paths += 2

# =========================================================================
# Write back
# =========================================================================
with open(JSON_PATH, "w") as f:
    json.dump(spec, f, indent=2)

print(f"[OK] Injected {added_paths} endpoint paths into {JSON_PATH}")
print(f"   - Employee: photo upload/delete + documents upload (3)")
print(f"   - Settings: gradings + villages search/detail + cascading (10)")
print(f"   - Platform: user password + user delete (2)")
