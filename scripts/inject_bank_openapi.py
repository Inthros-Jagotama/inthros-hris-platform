#!/usr/bin/env python3
"""Inject employee bank account endpoints into openapi.json."""

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

if "CreateBankAccountRequest" not in schemas:
    schemas["CreateBankAccountRequest"] = {
        "type": "object",
        "required": ["account_number", "account_name"],
        "properties": {
            "bank_id": {"type": "string", "format": "uuid", "description": "Bank ID from settings"},
            "account_number": {"type": "string", "maxLength": 50, "description": "Bank account number"},
            "account_name": {"type": "string", "maxLength": 255, "description": "Account holder name"}
        }
    }

if "UpdateBankAccountRequest" not in schemas:
    schemas["UpdateBankAccountRequest"] = {
        "type": "object",
        "properties": {
            "bank_id": {"type": "string", "format": "uuid", "description": "Bank ID from settings"},
            "account_number": {"type": "string", "maxLength": 50, "description": "Bank account number"},
            "account_name": {"type": "string", "maxLength": 255, "description": "Account holder name"}
        }
    }

if "BankAccountResponse" not in schemas:
    schemas["BankAccountResponse"] = {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "bank_id": {"type": "string", "format": "uuid", "description": "Bank ID from settings"},
            "account_number": {"type": "string"},
            "account_name": {"type": "string"}
        }
    }

# =========================================================================
# Endpoints
# =========================================================================
bank_paths = {
    "/api/v1/tenant/employees/{id}/banks": {
        "post": {
            "tags": ["Tenant: Employees"],
            "summary": "Create employee bank account",
            "operationId": "createEmployeeBank",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/CreateBankAccountRequest"}
                    }
                }
            },
            "responses": {
                "201": {"description": "Bank account created", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/BankAccountResponse"}}}},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Create a new bank account record for an employee. References a Bank from settings/banks."
        }
    },
    "/api/v1/tenant/employees/{id}/banks/{bankId}": {
        "put": {
            "tags": ["Tenant: Employees"],
            "summary": "Update employee bank account",
            "operationId": "updateEmployeeBank",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}},
                {"name": "bankId", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/UpdateBankAccountRequest"}
                    }
                }
            },
            "responses": {
                "200": {"description": "Bank account updated", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/BankAccountResponse"}}}},
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Bank account not found"}
            },
            "description": "Update an existing employee bank account record."
        },
        "delete": {
            "tags": ["Tenant: Employees"],
            "summary": "Delete employee bank account",
            "operationId": "deleteEmployeeBank",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}},
                {"name": "bankId", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}
            ],
            "responses": {
                "200": {"description": "Bank account deleted"},
                "401": {"description": "Unauthorized"},
                "404": {"description": "Bank account not found"}
            },
            "description": "Delete an employee bank account record."
        }
    }
}

spec["paths"].update(bank_paths)
added_paths += 3

# =========================================================================
# Write back
# =========================================================================
with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2)

print(f"[OK] Injected {added_paths} bank endpoints into {JSON_PATH}")
print(f"   - Employee Bank Accounts: POST + PUT + DELETE (3)")
print(f"   - Schemas added: CreateBankAccountRequest, UpdateBankAccountRequest, BankAccountResponse")
