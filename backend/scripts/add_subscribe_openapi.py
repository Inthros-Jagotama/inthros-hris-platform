#!/usr/bin/env python3
"""Add POST /api/v1/tenant/packages/{id}/subscribe to OpenAPI JSON.

1. Add POST method to the existing /api/v1/tenant/packages path
2. Add SubscribePackageRequest schema
"""

import json
import os

OPENAPI_PATH = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "internal/pkg/docs/openapi.json",
)

with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

paths = spec["paths"]
schemas = spec["components"]["schemas"]

tenant_pkg_path = "/api/v1/tenant/packages"
tenant_pkg_id_path = "/api/v1/tenant/packages/{id}/subscribe"

# =============================================================================
# 1. Add subscribe path
# =============================================================================
if tenant_pkg_id_path not in paths:
    new_path = {
        "post": {
            "tags": ["Tenant: Packages"],
            "summary": "Subscribe to a package (create/renew license)",
            "operationId": "subscribePackage",
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "Package ID to subscribe to"
                }
            ],
            "security": [{"bearerAuth": []}],
            "responses": {
                "201": {
                    "description": "Successfully subscribed to package",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {
                                        "type": "object",
                                        "properties": {
                                            "license_id": {"type": "string", "format": "uuid"},
                                            "license_key": {"type": "string"},
                                            "plan_type": {"type": "string", "example": "pro"},
                                            "package_id": {"type": "string", "format": "uuid"},
                                            "package_name": {"type": "string"}
                                        }
                                    },
                                    "message": {"type": "string", "example": "Successfully subscribed to package"}
                                }
                            }
                        }
                    }
                },
                "400": {
                    "description": "Validation error: package not published"
                },
                "401": {
                    "description": "Unauthorized: missing or invalid JWT"
                },
                "404": {
                    "description": "Package not found"
                },
                "500": {
                    "description": "Internal server error: license creation failed"
                }
            },
            "description": "Subscribe the authenticated company to a published package. "
                           "Creates a new license for the company associated with the specified package. "
                           "Requires JWT Bearer Token with company context. "
                           "The company_id is extracted from the JWT claims."
        }
    }

    # Insert after /api/v1/tenant/packages (the GET list endpoint)
    if tenant_pkg_path in paths:
        new_paths = {}
        for p, v in paths.items():
            new_paths[p] = v
            if p == tenant_pkg_path:
                new_paths[tenant_pkg_id_path] = new_path
        spec["paths"] = new_paths
        print("[OK] Added '{}' path after '{}'".format(tenant_pkg_id_path, tenant_pkg_path))
    else:
        spec["paths"][tenant_pkg_id_path] = new_path
        print("[OK] Added '{}' path at end".format(tenant_pkg_id_path))
else:
    print("[SKIP] {} path already exists".format(tenant_pkg_id_path))

# =============================================================================
# Write updated JSON
# =============================================================================
with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print("\n[OK] OpenAPI JSON updated successfully!")
print("   File: {}".format(OPENAPI_PATH))
