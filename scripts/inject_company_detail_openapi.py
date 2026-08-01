#!/usr/bin/env python3
"""Inject tenant self-service company endpoint (GET /companies/me) into openapi.json."""

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

if "UpdateCurrentCompanyRequest" not in schemas:
    schemas["UpdateCurrentCompanyRequest"] = {
        "type": "object",
        "properties": {
            "email": {"type": "string", "format": "email", "nullable": True},
            "phone": {"type": "string", "maxLength": 20, "nullable": True},
            "address": {"type": "string", "nullable": True},
            "npwp": {"type": "string", "maxLength": 16, "nullable": True},
            "nib": {"type": "string", "maxLength": 25, "nullable": True}
        }
    }

if "TenantCompanyDetailResponse" not in schemas:
    schemas["TenantCompanyDetailResponse"] = {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "name": {"type": "string"},
            "slug": {"type": "string"},
            "subdomain": {"type": "string", "nullable": True},
            "domain": {"type": "string", "nullable": True},
            "npwp": {"type": "string", "nullable": True},
            "nib": {"type": "string", "nullable": True},
            "address": {"type": "string", "nullable": True},
            "email": {"type": "string", "format": "email", "nullable": True},
            "phone": {"type": "string", "nullable": True},
            "status": {"type": "string"},
            "admin_user": {
                "type": "object",
                "nullable": True,
                "properties": {
                    "id": {"type": "string", "format": "uuid"},
                    "name": {"type": "string"},
                    "email": {"type": "string", "format": "email"},
                    "role": {"type": "string"}
                }
            },
            "license_info": {
                "type": "object",
                "nullable": True,
                "properties": {
                    "id": {"type": "string", "format": "uuid"},
                    "license_key": {"type": "string"},
                    "plan_type": {"type": "string"},
                    "package_id": {"type": "string", "format": "uuid"},
                    "package_name": {"type": "string"},
                    "start_date": {"type": "string", "format": "date-time"},
                    "end_date": {"type": "string", "format": "date-time"},
                    "max_employees": {"type": "integer", "format": "int32", "description": "Employee quota from license (0 = unlimited)"},
                    "employee_total": {"type": "integer", "format": "int64", "description": "Current active employee count in tenant database"}
                }
            },
            "provisioning_info": {
                "type": "object",
                "nullable": True,
                "properties": {
                    "provisioned": {"type": "boolean"},
                    "is_active": {"type": "boolean"},
                    "driver": {"type": "string"},
                    "db_name": {"type": "string"}
                }
            },
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    }
else:
    # Update idempotent: pastikan license_info punya field baru (package_name,
    # start_date, end_date) meski schema sudah ada dari run sebelumnya.
    lic_props = schemas["TenantCompanyDetailResponse"]["properties"]["license_info"]["properties"]
    lic_props.setdefault("package_name", {"type": "string"})
    lic_props.setdefault("start_date", {"type": "string", "format": "date-time"})
    lic_props.setdefault("end_date", {"type": "string", "format": "date-time"})
    lic_props.setdefault("max_employees", {"type": "integer", "format": "int32", "description": "Employee quota from license (0 = unlimited)"})
    lic_props.setdefault("employee_total", {"type": "integer", "format": "int64", "description": "Current active employee count in tenant database"})

# Ensure tag exists
tags = spec.setdefault("tags", [])
if not any(t["name"] == "Tenant: Company" for t in tags):
    tags.append({"name": "Tenant: Company", "description": "Self-service company endpoints for the authenticated tenant user"})

# =========================================================================
# Endpoints
# =========================================================================
paths = spec["paths"]
if "/api/v1/tenant/companies/me" not in paths:
    paths["/api/v1/tenant/companies/me"] = {
        "get": {
            "tags": ["Tenant: Company"],
            "summary": "Get current company detail",
            "operationId": "getCurrentCompany",
            "security": [{"bearerAuth": []}],
            "responses": {
                "200": {
                    "description": "Current company detail",
                    "content": {"application/json": {"schema": {"$ref": "#/components/schemas/TenantCompanyDetailResponse"}}}
                },
                "401": {"description": "Unauthorized"},
                "404": {"description": "Company not found"}
            },
            "description": "Retrieve the profile of the company the authenticated tenant user belongs to. Company is resolved from the tenant context (X-Tenant-ID / JWT claims)."
        },
        "put": {
            "tags": ["Tenant: Company"],
            "summary": "Update current company information",
            "operationId": "updateCurrentCompany",
            "security": [{"bearerAuth": []}],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/UpdateCurrentCompanyRequest"}
                    }
                }
            },
            "responses": {
                "200": {
                    "description": "Company updated",
                    "content": {"application/json": {"schema": {"$ref": "#/components/schemas/TenantCompanyDetailResponse"}}}
                },
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Company not found"}
            },
            "description": "Update the tenant's own company profile (email, phone, address, NPWP, NIB). Company is resolved from the tenant context; name/subdomain/domain are managed by platform admin."
        }
    }
    added_paths += 2
else:
    # Update idempotent: pastikan PUT ada bila path sudah ada dari run sebelumnya.
    me_path = paths["/api/v1/tenant/companies/me"]
    if "put" not in me_path:
        me_path["put"] = {
            "tags": ["Tenant: Company"],
            "summary": "Update current company information",
            "operationId": "updateCurrentCompany",
            "security": [{"bearerAuth": []}],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/UpdateCurrentCompanyRequest"}
                    }
                }
            },
            "responses": {
                "200": {
                    "description": "Company updated",
                    "content": {"application/json": {"schema": {"$ref": "#/components/schemas/TenantCompanyDetailResponse"}}}
                },
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Company not found"}
            },
            "description": "Update the tenant's own company profile (email, phone, address, NPWP, NIB). Company is resolved from the tenant context; name/subdomain/domain are managed by platform admin."
        }
        added_paths += 1

# =========================================================================
# Write back
# =========================================================================
with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2)

print(f"[OK] Injected {added_paths} company endpoint into {JSON_PATH}")
print("   - GET /api/v1/tenant/companies/me (Tenant: Company)")
print("   - Schema added: TenantCompanyDetailResponse")
