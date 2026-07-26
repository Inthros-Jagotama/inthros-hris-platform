#!/usr/bin/env python3
"""Update OpenAPI JSON for company endpoint with package/license info.

1. Add package_id field to CreateCompanyRequest schema
2. Update createCompany endpoint description to mention license_info
3. Update getCompany and listCompanies descriptions to mention license_info
"""

import json
import os

OPENAPI_PATH = os.path.join(
    os.path.dirname(os.path.dirname(os.path.abspath(__file__))),
    "internal/pkg/docs/openapi.json",
)

with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

schemas = spec["components"]["schemas"]
paths = spec["paths"]

# =============================================================================
# 1. Update CreateCompanyRequest -- add package_id field
# =============================================================================
ccr = schemas.get("CreateCompanyRequest")
if ccr:
    props = ccr.get("properties", {})
    if "package_id" not in props:
        props["package_id"] = {
            "type": "string",
            "format": "uuid",
            "description": "Optional: associate this company with a published package (auto-creates license on signup)",
        }
        print("[OK] Added 'package_id' to CreateCompanyRequest")
    else:
        print("[SKIP] CreateCompanyRequest already has package_id")
else:
    print("[WARN] CreateCompanyRequest schema not found")

# =============================================================================
# 2. Update createCompany endpoint description
# =============================================================================
create_ep = paths.get("/api/v1/platform/companies", {}).get("post", {})
if create_ep:
    desc = create_ep.get("description", "")
    if "license_info" not in desc:
        create_ep["description"] = (
            "Register a new company tenant. Also creates a company_admin user automatically. "
            "Request now includes admin_name, admin_email, admin_password fields. "
            "If package_id is provided, a license is auto-created from the published package "
            "and the response includes license_info with id, license_key, plan_type, package_id. "
            "This triggers the full tenant provisioning flow: creates the company record, "
            "provisions a new database, runs all migrations (13 files, ~111 tables), "
            "activates the tenant, and creates the admin user. "
            "The process is atomic -- if provisioning fails, the company is marked as suspended."
        )
        print("[OK] Updated createCompany description with license_info")
    else:
        print("[SKIP] createCompany description already has license_info")
else:
    print("[WARN] createCompany endpoint not found")

# =============================================================================
# 3. Update 201 response description for createCompany
# =============================================================================
responses_201 = create_ep.get("responses", {}).get("201", {})
if responses_201:
    resp_desc = responses_201.get("description", "")
    if "license_info" not in resp_desc:
        responses_201["description"] = (
            "Company created successfully. Response includes admin_user "
            "(id, name, email, role) and if package_id was provided, license_info "
            "(id, license_key, plan_type, package_id)."
        )
        print("[OK] Updated 201 response description for createCompany")
    else:
        print("[SKIP] 201 response already has license_info")

# =============================================================================
# 4. Update getCompany endpoint description
# =============================================================================
get_ep = paths.get("/api/v1/platform/companies/{id}", {}).get("get", {})
if get_ep:
    desc = get_ep.get("description", "")
    if "license_info" not in desc:
        get_ep["description"] = (
            "Get detailed information about a specific company/tenant including its status, "
            "contact details, subscription plan, database connection health, "
            "and associated license_info (if any)."
        )
        print("[OK] Updated getCompany description with license_info")
    else:
        print("[SKIP] getCompany description already has license_info")

    # Update 200 response description
    resp_200 = get_ep.get("responses", {}).get("200", {})
    if resp_200:
        resp_desc = resp_200.get("description", "")
        if "license_info" not in resp_desc:
            resp_200["description"] = (
                "Company details with admin_user and license_info (if applicable)"
            )
            print("[OK] Updated getCompany 200 response description")
else:
    print("[WARN] getCompany endpoint not found")

# =============================================================================
# 5. Update listCompanies endpoint description
# =============================================================================
list_ep = paths.get("/api/v1/platform/companies", {}).get("get", {})
if list_ep:
    desc = list_ep.get("description", "")
    if "license_info" not in desc:
        list_ep["description"] = (
            "Retrieve a paginated list of all registered companies (tenants) in the platform. "
            "Includes company status, contact information, subscription details, "
            "and associated license_info per company."
        )
        print("[OK] Updated listCompanies description with license_info")
    else:
        print("[SKIP] listCompanies description already has license_info")

    # Update 200 response
    resp_200 = list_ep.get("responses", {}).get("200", {})
    if resp_200:
        resp_desc = resp_200.get("description", "")
        if "license_info" not in resp_desc:
            resp_200["description"] = (
                "Paginated list of companies with admin_user and license_info"
            )
            print("[OK] Updated listCompanies 200 response description")
else:
    print("[WARN] listCompanies endpoint not found")

# =============================================================================
# Write updated JSON
# =============================================================================
with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print("\n[OK] OpenAPI JSON updated successfully!")
print("   File: {}".format(OPENAPI_PATH))
