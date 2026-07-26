#!/usr/bin/env python3
"""Add POST /api/v1/tenant/packages/{id}/unsubscribe to OpenAPI JSON.
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

tenant_pkg_sub_path = "/api/v1/tenant/packages/{id}/subscribe"
tenant_pkg_unsub_path = "/api/v1/tenant/packages/{id}/unsubscribe"

# =============================================================================
# Add unsubscribe path
# =============================================================================
if tenant_pkg_unsub_path not in paths:
    new_path = {
        "post": {
            "tags": ["Tenant: Packages"],
            "summary": "Unsubscribe from a package (deactivate modules + suspend license)",
            "operationId": "unsubscribePackage",
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    },
                    "description": "Package ID to unsubscribe from"
                }
            ],
            "security": [{"bearerAuth": []}],
            "responses": {
                "200": {
                    "description": "Successfully unsubscribed from package",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "message": {"type": "string", "example": "Successfully unsubscribed from package"}
                                }
                            }
                        }
                    }
                },
                "401": {
                    "description": "Unauthorized: missing or invalid JWT"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Unsubscribe the authenticated company from a package. "
                           "Deactivates all modules included in the package and suspends "
                           "the active license associated with this company and package. "
                           "Requires JWT Bearer Token with company context."
        }
    }

    # Insert after the subscribe path
    if tenant_pkg_sub_path in paths:
        new_paths = {}
        for p, v in paths.items():
            new_paths[p] = v
            if p == tenant_pkg_sub_path:
                new_paths[tenant_pkg_unsub_path] = new_path
        spec["paths"] = new_paths
        print("[OK] Added '{}' path after '{}'".format(tenant_pkg_unsub_path, tenant_pkg_sub_path))
    else:
        spec["paths"][tenant_pkg_unsub_path] = new_path
        print("[OK] Added '{}' path at end".format(tenant_pkg_unsub_path))
else:
    print("[SKIP] {} path already exists".format(tenant_pkg_unsub_path))

with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print("\n[OK] OpenAPI JSON updated successfully!")
