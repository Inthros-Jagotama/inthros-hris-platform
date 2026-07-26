#!/usr/bin/env python3
"""Add tenant /packages endpoint to OpenAPI JSON.

1. Add "Tenant: Packages" tag after "Tenant: Time & Attendance"
2. Add GET /api/v1/tenant/packages path with auth bearer
"""

import json
import os

OPENAPI_PATH = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "internal/pkg/docs/openapi.json",
)

with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

tags = spec["tags"]
paths = spec["paths"]

# =============================================================================
# 1. Add Tenant: Packages tag (after Tenant: Time & Attendance)
# =============================================================================
tag_exists = any(t["name"] == "Tenant: Packages" for t in tags)
if not tag_exists:
    # Find insertion index: after "Tenant: Time & Attendance"
    insert_idx = None
    for i, t in enumerate(tags):
        if t["name"] == "Tenant: Time & Attendance":
            insert_idx = i + 1
            break
    if insert_idx is None:
        insert_idx = len(tags)  # append at end if not found

    tags.insert(insert_idx, {
        "name": "Tenant: Packages",
        "description": "Published package browsing for authenticated tenant users"
    })
    print("[OK] Added 'Tenant: Packages' tag at index {}".format(insert_idx))
else:
    print("[SKIP] Tenant: Packages tag already exists")

# =============================================================================
# 2. Add GET /api/v1/tenant/packages path
# =============================================================================
tenant_pkg_path = "/api/v1/tenant/packages"
if tenant_pkg_path not in paths:
    # Find insertion point: after /api/v1/tenant/organizations/{id} (last org path)
    # Or after the last tenant path before employee paths
    # For simplicity, insert before /api/v1/tenant/employees
    insert_before = None
    for p in paths:
        if p.startswith("/api/v1/tenant/employees"):
            insert_before = p
            break

    new_path = {
        "get": {
            "tags": ["Tenant: Packages"],
            "summary": "List published packages (tenant)",
            "operationId": "listTenantPackages",
            "security": [{"bearerAuth": []}],
            "responses": {
                "200": {
                    "description": "List of published packages with modules",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {
                                        "type": "array",
                                        "items": {
                                            "$ref": "#/components/schemas/PublicPackageResponse"
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "401": {
                    "description": "Unauthorized: missing or invalid JWT"
                }
            },
            "description": "Retrieve a list of published packages for authenticated tenant users. "
                           "Requires JWT Bearer Token. Returns the same data as the public "
                           "endpoint but within the authenticated tenant context. "
                           "Useful for company admins to browse available packages for subscription or upgrade."
        }
    }

    # Insert before the employees path
    if insert_before:
        new_paths = {}
        for p, v in paths.items():
            if p == insert_before:
                new_paths[tenant_pkg_path] = new_path
            new_paths[p] = v
        spec["paths"] = new_paths
        print("[OK] Added '{}' path before '{}'".format(tenant_pkg_path, insert_before))
    else:
        spec["paths"][tenant_pkg_path] = new_path
        print("[OK] Added '{}' path at end".format(tenant_pkg_path))
else:
    print("[SKIP] {} path already exists".format(tenant_pkg_path))

# =============================================================================
# Write updated JSON
# =============================================================================
with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print("\n[OK] OpenAPI JSON updated successfully!")
print("   File: {}".format(OPENAPI_PATH))
