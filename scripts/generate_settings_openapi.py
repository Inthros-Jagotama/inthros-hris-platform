"""
Generate OpenAPI paths and schemas for the Settings module (14 CRUD entities).
Inserts them into backend/internal/pkg/docs/openapi.json.

Usage: python scripts/generate_settings_openapi.py
"""

import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OPENAPI_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SETTINGS_TAG = {
    "name": "Tenant: Settings",
    "description": "Settings & Master Data Reference -- manage zones, provinces, regencies, districts, villages, educations, religions, marital statuses, relationship types, banks, employment statuses, nationalities, job families, and salary grades. CRUD operations for all tenant reference data."
}

ENTITIES = [
    {"slug": "zones", "name": "Zone", "display": "Zones", "op_prefix": "settingZone",
     "extra_fields": [{"name": "region", "type": "string", "desc": "Region name"}, {"name": "is_active", "type": "boolean", "desc": "Whether the zone is active"}],
     "parent_field": None},
    {"slug": "provinces", "name": "Province", "display": "Provinces", "op_prefix": "settingProvince", "extra_fields": [], "parent_field": None},
    {"slug": "regencies", "name": "Regency", "display": "Regencies", "op_prefix": "settingRegency", "extra_fields": [], "parent_field": {"name": "province_id", "type": "string", "format": "uuid", "desc": "Parent province ID"}},
    {"slug": "districts", "name": "District", "display": "Districts", "op_prefix": "settingDistrict", "extra_fields": [], "parent_field": {"name": "regency_id", "type": "string", "format": "uuid", "desc": "Parent regency ID"}},
    {"slug": "villages", "name": "Village", "display": "Villages", "op_prefix": "settingVillage", "extra_fields": [], "parent_field": {"name": "district_id", "type": "string", "format": "uuid", "desc": "Parent district ID"}},
    {"slug": "educations", "name": "Education", "display": "Educations", "op_prefix": "settingEducation", "extra_fields": [], "parent_field": None},
    {"slug": "religions", "name": "Religion", "display": "Religions", "op_prefix": "settingReligion", "extra_fields": [], "parent_field": None},
    {"slug": "marital-statuses", "name": "MaritalStatus", "display": "Marital Statuses", "op_prefix": "settingMaritalStatus", "extra_fields": [], "parent_field": None},
    {"slug": "relationship-types", "name": "RelationshipType", "display": "Relationship Types", "op_prefix": "settingRelationshipType", "extra_fields": [], "parent_field": None},
    {"slug": "banks", "name": "Bank", "display": "Banks", "op_prefix": "settingBank", "extra_fields": [], "parent_field": None},
    {"slug": "employment-statuses", "name": "EmploymentStatus", "display": "Employment Statuses", "op_prefix": "settingEmploymentStatus", "extra_fields": [], "parent_field": None},
    {"slug": "nationalities", "name": "Nationality", "display": "Nationalities", "op_prefix": "settingNationality", "extra_fields": [], "parent_field": None},
    {"slug": "job-families", "name": "JobFamily", "display": "Job Families", "op_prefix": "settingJobFamily",
     "extra_fields": [{"name": "description", "type": "string", "desc": "Job family description"}], "parent_field": None},
    {"slug": "salary-grades", "name": "SalaryGrade", "display": "Salary Grades", "op_prefix": "settingSalaryGrade",
     "extra_fields": [{"name": "description", "type": "string", "desc": "Salary grade description"}, {"name": "min_amount", "type": "number", "desc": "Minimum salary amount"}, {"name": "max_amount", "type": "number", "desc": "Maximum salary amount"}], "parent_field": None},
]


def build_properties(extra_fields, parent_field):
    props = {"id": {"type": "string", "format": "uuid"}, "code": {"type": "string"}, "name": {"type": "string"}}
    if parent_field:
        p = {"type": parent_field["type"], "description": parent_field["desc"]}
        if parent_field.get("format"): p["format"] = parent_field["format"]
        props[parent_field["name"]] = p
    for ef in extra_fields:
        props[ef["name"]] = {"type": ef["type"], "description": ef["desc"]}
    props["sort_order"] = {"type": "integer"}
    props["created_at"] = {"type": "string", "format": "date-time"}
    props["updated_at"] = {"type": "string", "format": "date-time"}
    return props


def gen_schemas():
    out = {}
    for e in ENTITIES:
        en = e["name"]
        resp_name = f"{en}Response"
        out[resp_name] = {"type": "object", "properties": build_properties(e["extra_fields"], e["parent_field"])}
        out[f"{en}PaginatedResponse"] = {
            "type": "object",
            "properties": {
                "success": {"type": "boolean"},
                "data": {"type": "array", "items": {"$ref": f"#/components/schemas/{resp_name}"}},
                "page": {"type": "integer"}, "per_page": {"type": "integer"}, "total": {"type": "integer"}, "total_pages": {"type": "integer"}
            }
        }
        # Create request
        c = {"type": "object", "required": ["code", "name"], "properties": {"code": {"type": "string"}, "name": {"type": "string"}}}
        if e["parent_field"]:
            pf = e["parent_field"]; c["required"].append(pf["name"])
            p = {"type": pf["type"], "description": pf["desc"]}
            if pf.get("format"): p["format"] = pf["format"]
            c["properties"][pf["name"]] = p
        for ef in e["extra_fields"]:
            c["properties"][ef["name"]] = {"type": ef["type"], "description": ef["desc"]}
        c["properties"]["sort_order"] = {"type": "integer"}
        out[f"Create{en}Request"] = c
        # Update request
        u = {"type": "object", "properties": {"code": {"type": "string"}, "name": {"type": "string"}}}
        if e["parent_field"]:
            pf = e["parent_field"]; p = {"type": pf["type"], "description": pf["desc"]}
            if pf.get("format"): p["format"] = pf["format"]
            u["properties"][pf["name"]] = p
        for ef in e["extra_fields"]:
            u["properties"][ef["name"]] = {"type": ef["type"], "description": ef["desc"]}
        u["properties"]["sort_order"] = {"type": "integer"}
        out[f"Update{en}Request"] = u
    return out


def gen_paths():
    out = {}
    for e in ENTITIES:
        slug, en, disp, op = e["slug"], e["name"], e["display"], e["op_prefix"]
        tag, resp, pag = ["Tenant: Settings"], f"{en}Response", f"{en}PaginatedResponse"
        pk, pkid = f"/api/v1/tenant/settings/{slug}", f"/api/v1/tenant/settings/{slug}/{{id}}"
        lp = [{"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
              {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}]
        out[pk] = {
            "get": {"tags": tag, "summary": f"List all {disp}", "operationId": f"{op}List", "parameters": lp,
                    "security": [{"bearerAuth": []}],
                    "responses": {"200": {"description": f"Paginated list of {disp}", "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{pag}"}}}},
                                  "401": {"description": "Unauthorized: missing or invalid JWT"}},
                    "description": f"Retrieve a paginated list of {disp}. Supports pagination parameters."},
            "post": {"tags": tag, "summary": f"Create a new {en.lower()}", "operationId": f"{op}Create",
                     "security": [{"bearerAuth": []}],
                     "requestBody": {"required": True, "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/Create{en}Request"}}}},
                     "responses": {"201": {"description": f"{en} created successfully", "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{resp}"}}}},
                                   "400": {"description": "Validation error"}, "401": {"description": "Unauthorized: missing or invalid JWT"}},
                     "description": f"Create a new {disp.lower()} record. Validates required fields and returns the created resource with its assigned ID."}
        }
        out[pkid] = {
            "get": {"tags": tag, "summary": f"Get {en.lower()} by ID", "operationId": f"{op}Get",
                    "parameters": [{"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}],
                    "security": [{"bearerAuth": []}],
                    "responses": {"200": {"description": f"{en} details", "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{resp}"}}}},
                                  "401": {"description": "Unauthorized"}, "404": {"description": f"{en} not found"}},
                    "description": f"Get detailed information about a specific {en.lower()} by its ID."},
            "put": {"tags": tag, "summary": f"Update {en.lower()}", "operationId": f"{op}Update",
                    "parameters": [{"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}],
                    "security": [{"bearerAuth": []}],
                    "requestBody": {"required": True, "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/Update{en}Request"}}}},
                    "responses": {"200": {"description": f"{en} updated", "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{resp}"}}}},
                                  "400": {"description": "Validation error"}, "401": {"description": "Unauthorized"}, "404": {"description": f"{en} not found"}},
                    "description": f"Update a {en.lower()} record's details including code, name, and other attributes."},
            "delete": {"tags": tag, "summary": f"Delete {en.lower()}", "operationId": f"{op}Delete",
                       "parameters": [{"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}],
                       "security": [{"bearerAuth": []}],
                       "responses": {"200": {"description": f"{en} deleted"}, "401": {"description": "Unauthorized"}, "404": {"description": f"{en} not found"}},
                       "description": f"Soft-delete a {en.lower()} record. Sets the deleted_at timestamp and hides it from standard queries."}
        }
    return out


def inner_json(obj, indent):
    """Serialize obj to JSON, removing outer braces, indenting by `indent` spaces."""
    raw = json.dumps(obj, indent=2, ensure_ascii=False)
    # Remove outer { and }
    lines = raw.split("\n")
    if lines[0].strip() == "{" and lines[-1].strip() == "}":
        lines = lines[1:-1]
    # Re-indent
    result = []
    for line in lines:
        if line.strip():
            result.append(" " * indent + line)
        else:
            result.append(line)
    return "\n".join(result)


def main():
    print(f"Reading {OPENAPI_PATH}...", flush=True)
    with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
        content = f.read()

    # ── 1. Insert tag ──
    print("[1/3] Adding Tenant: Settings tag...", flush=True)
    tag_json = json.dumps(SETTINGS_TAG, indent=2, ensure_ascii=False)
    tag_lines = tag_json.split("\n")
    tag_indented = "\n".join("    " + l for l in tag_lines)  # 4-space indent to match other tags
    old_tag = 'leadership pipeline readiness."\n    }\n  ],'
    new_tag = 'leadership pipeline readiness."\n    },\n' + tag_indented + "\n  ],"
    if old_tag in content:
        content = content.replace(old_tag, new_tag, 1)
        print("  OK", flush=True)
    else:
        print("  FAIL - Could not find tag marker", flush=True)
        err_idx = content.find("leadership pipeline readiness")
        if err_idx >= 0:
            print(f"  Found marker at position {err_idx}", flush=True)
            print(f"  Context: {repr(content[err_idx:err_idx+100])}", flush=True)
        return

    # ── 2. Insert paths ──
    print("[2/3] Adding 70 Settings paths...", flush=True)
    paths = gen_paths()
    paths_inner = inner_json(paths, 4)  # 4-space indent to match existing paths
    old_paths = 'Soft-delete a succession plan."\n      }\n    }\n  },'
    new_paths = 'Soft-delete a succession plan."\n      }\n    },\n' + paths_inner + "\n  },"
    if old_paths in content:
        content = content.replace(old_paths, new_paths, 1)
        print("  OK", flush=True)
    else:
        print("  FAIL - Could not find paths marker", flush=True)
        return

    # ── 3. Insert schemas ──
    print("[3/3] Adding 42 Settings schemas...", flush=True)
    schemas = gen_schemas()
    schemas_inner = inner_json(schemas, 6)  # 6-space indent to match existing schemas
    old_schemas = 'satisfied within this package"\n          }\n        }\n      }\n    }\n  }\n}'
    new_schemas = 'satisfied within this package"\n          }\n        }\n      },\n' + schemas_inner + "\n    }\n  }\n}"
    if old_schemas in content:
        content = content.replace(old_schemas, new_schemas, 1)
        print("  OK", flush=True)
    else:
        print("  FAIL - Could not find schemas marker", flush=True)
        return

    # ── Write ──
    print(f"\nWriting to {OPENAPI_PATH}...", flush=True)
    with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
        f.write(content)

    # ── Validate ──
    print("Validating JSON...", flush=True)
    try:
        with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
            json.load(f)
        print("OK - JSON is valid!", flush=True)
    except json.JSONDecodeError as exc:
        print(f"FAIL - JSON is invalid: {exc}", flush=True)

    print("\nDone!", flush=True)


if __name__ == "__main__":
    main()
