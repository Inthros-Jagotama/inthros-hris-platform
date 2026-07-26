#!/usr/bin/env python3
"""Add activated_modules field to subscribe endpoint response in OpenAPI JSON.
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

tenant_pkg_id_path = "/api/v1/tenant/packages/{id}/subscribe"

ep = paths.get(tenant_pkg_id_path, {}).get("post", {})
resp_201 = ep.get("responses", {}).get("201", {})
schema = resp_201.get("content", {}).get("application/json", {}).get("schema", {})
data_props = schema.get("properties", {}).get("data", {}).get("properties", {})

if "activated_modules" not in data_props:
    data_props["activated_modules"] = {
        "type": "array",
        "items": {"type": "string"},
        "description": "List of module names that were auto-activated for this company"
    }
    print("[OK] Added 'activated_modules' array to subscribe 201 response")
else:
    print("[SKIP] activated_modules already exists")

# Also update description to mention auto-activation
desc = ep.get("description", "")
if "auto-activates" not in desc:
    ep["description"] = (
        "Subscribe the authenticated company to a published package. "
        "Creates a new license for the company associated with the specified package "
        "and auto-activates all modules included in the package. "
        "Requires JWT Bearer Token with company context. "
        "The company_id is extracted from the JWT claims."
    )
    print("[OK] Updated subscribe description with auto-activation info")

with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print("\n[OK] OpenAPI JSON updated successfully!")
