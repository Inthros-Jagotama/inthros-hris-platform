#!/usr/bin/env python3
"""
Add Career Intelligence module (19 endpoints, ~15 schemas) to the existing openapi.json document.
"""

import json, sys, os

OPENAPI_PATH = os.path.join(os.path.dirname(__file__), '..', 'internal', 'pkg', 'docs', 'openapi.json')

with open(OPENAPI_PATH, 'r', encoding='utf-8') as f:
    spec = json.load(f)

TAG = "Tenant: Career Intelligence"
TAG_DESC = "Career Intelligence & Talent Management — strategic talent analytics for 9-box talent mapping, career interests tracking, career path gap analysis, and succession planning. Provides talent review data to identify high-potential employees, plan career development, and ensure leadership pipeline readiness."

# ─────────────────────────────────────────────────────────────────────────────
# 1. Add tag (if not already present)
# ─────────────────────────────────────────────────────────────────────────────
existing_tags = [t['name'] for t in spec['tags']]
if TAG not in existing_tags:
    spec['tags'].append({"name": TAG, "description": TAG_DESC})

# ─────────────────────────────────────────────────────────────────────────────
# 2. Schemas
# ─────────────────────────────────────────────────────────────────────────────
if 'components' not in spec:
    spec['components'] = {}
if 'schemas' not in spec['components']:
    spec['components']['schemas'] = {}

schemas = {
    # ── Talent Maps ──
    "CreateTalentMapRequest": {
        "type": "object",
        "required": ["employee_id", "period", "performance", "potential"],
        "properties": {
            "employee_id": {"type": "string", "format": "uuid", "description": "ID of the employee being assessed"},
            "period": {"type": "string", "maxLength": 7, "example": "2026-Q1", "description": "Assessment period (e.g., 2026-Q1, 2026-H1)"},
            "performance": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"], "description": "Performance rating"},
            "potential": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"], "description": "Potential rating"},
            "notes": {"type": "string", "description": "Additional assessment notes"}
        }
    },
    "UpdateTalentMapRequest": {
        "type": "object",
        "properties": {
            "performance": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"], "description": "Updated performance rating"},
            "potential": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"], "description": "Updated potential rating"},
            "notes": {"type": "string", "description": "Updated assessment notes"}
        }
    },
    "TalentMapResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "period": {"type": "string"},
            "performance": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"]},
            "potential": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"]},
            "grid_position": {"type": "string", "description": "9-Box grid position (e.g., 9-BOX-1, 9-BOX-5, 9-BOX-9)"},
            "notes": {"type": "string"},
            "assessor_id": {"type": "string", "format": "uuid", "description": "ID of the assessor (manager/HR)"},
            "assessed_at": {"type": "string", "format": "date-time"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "TalentGridResponse": {
        "type": "object",
        "properties": {
            "period": {"type": "string"},
            "quadrants": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/TalentQuadrant"}
            },
            "total": {"type": "integer"}
        }
    },
    "TalentQuadrant": {
        "type": "object",
        "properties": {
            "label": {"type": "string", "description": "Quadrant label (e.g., 'High Performer - High Potential')"},
            "position": {"type": "string", "description": "Grid position code (e.g., '9-BOX-1')"},
            "count": {"type": "integer", "description": "Number of employees in this quadrant"},
            "description": {"type": "string", "description": "Quadrant description"}
        }
    },
    "EmployeeTalentProfileResponse": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "current_map": {"$ref": "#/components/schemas/TalentMapResponse", "description": "Current period's talent map entry"},
            "history": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/TalentMapResponse"},
                "description": "Historical talent map entries"
            },
            "interests": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/CareerInterestResponse"},
                "description": "Employee's career interests"
            },
            "ready_for": {
                "type": "array",
                "items": {"type": "string"},
                "description": "Recommended next positions"
            }
        }
    },

    # ── Career Interests ──
    "CreateCareerInterestRequest": {
        "type": "object",
        "required": ["employee_id", "interest_type"],
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "interest_type": {"type": "string", "enum": ["LEADERSHIP", "SPECIALIST", "INTERNATIONAL", "ENTREPRENEUR"]},
            "target_position": {"type": "string", "description": "Desired position title"},
            "target_department": {"type": "string", "description": "Desired department"},
            "motivation": {"type": "string", "description": "Career motivation or reason"},
            "readiness_level": {"type": "string", "enum": ["NOW", "1_YEAR", "2_3_YEARS", "3_PLUS"], "description": "Readiness to move"}
        }
    },
    "CareerInterestResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "employee_id": {"type": "string", "format": "uuid"},
            "interest_type": {"type": "string"},
            "target_position": {"type": "string"},
            "target_department": {"type": "string"},
            "motivation": {"type": "string"},
            "readiness_level": {"type": "string"},
            "is_active": {"type": "boolean"},
            "recorded_at": {"type": "string", "format": "date"},
            "created_at": {"type": "string", "format": "date-time"}
        }
    },

    # ── Career Paths ──
    "CreateCareerPathRequest": {
        "type": "object",
        "required": ["source_title_id", "target_title_id", "path_type"],
        "properties": {
            "source_title_id": {"type": "string", "format": "uuid", "description": "Source position title ID"},
            "target_title_id": {"type": "string", "format": "uuid", "description": "Target position title ID"},
            "path_type": {"type": "string", "enum": ["PROMOTION", "LATERAL", "DEMOTION", "CROSSFUNCTIONAL"]},
            "typical_tenure": {"type": "integer", "description": "Typical tenure in months before moving"},
            "requirements": {"type": "string", "description": "Requirements for this career path"},
            "competencies": {"type": "string", "description": "Required competencies"}
        }
    },
    "CareerPathResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "source_title_id": {"type": "string", "format": "uuid"},
            "target_title_id": {"type": "string", "format": "uuid"},
            "path_type": {"type": "string"},
            "typical_tenure": {"type": "integer"},
            "requirements": {"type": "string"},
            "competencies": {"type": "string"},
            "certifications": {"type": "string"},
            "is_active": {"type": "boolean"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "CareerGapAnalysisResponse": {
        "type": "object",
        "properties": {
            "employee_id": {"type": "string", "format": "uuid"},
            "target_title": {"type": "string", "description": "Target position title name"},
            "matched_skills": {"type": "integer"},
            "total_required": {"type": "integer"},
            "gap_percentage": {"type": "number", "format": "float"},
            "recommendations": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/CareerGapRecommendation"}
            },
            "estimated_timeline": {"type": "string", "description": "Estimated timeline (e.g., '12-18 months')"}
        }
    },
    "CareerGapRecommendation": {
        "type": "object",
        "properties": {
            "category": {"type": "string", "enum": ["TRAINING", "EXPERIENCE", "CERTIFICATION"]},
            "description": {"type": "string"},
            "priority": {"type": "string", "enum": ["HIGH", "MEDIUM", "LOW"]}
        }
    },

    # ── Succession Plans ──
    "CreateSuccessionPlanRequest": {
        "type": "object",
        "required": ["position_id", "successor_id", "readiness_level"],
        "properties": {
            "position_id": {"type": "string", "format": "uuid", "description": "Position being planned for succession"},
            "successor_id": {"type": "string", "format": "uuid", "description": "Employee ID of the potential successor"},
            "readiness_level": {"type": "string", "enum": ["READY_NOW", "READY_1YR", "READY_2YR", "POTENTIAL"]},
            "priority_order": {"type": "integer", "description": "Priority order among multiple successors"},
            "target_date": {"type": "string", "format": "date", "description": "Target date for succession"},
            "development_plan": {"type": "string", "description": "Development plan for the successor"},
            "notes": {"type": "string"}
        }
    },
    "UpdateSuccessionPlanRequest": {
        "type": "object",
        "properties": {
            "readiness_level": {"type": "string", "enum": ["READY_NOW", "READY_1YR", "READY_2YR", "POTENTIAL"]},
            "priority_order": {"type": "integer"},
            "target_date": {"type": "string", "format": "date"},
            "development_plan": {"type": "string"},
            "notes": {"type": "string"}
        }
    },
    "SuccessionPlanResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string", "format": "uuid"},
            "position_id": {"type": "string", "format": "uuid"},
            "successor_id": {"type": "string", "format": "uuid"},
            "readiness_level": {"type": "string"},
            "priority_order": {"type": "integer"},
            "target_date": {"type": "string", "format": "date"},
            "development_plan": {"type": "string"},
            "notes": {"type": "string"},
            "status": {"type": "string", "enum": ["ACTIVE", "INACTIVE", "FILLED"]},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },

    # ── Pagination ──
    "PaginatedResponseCI": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {"type": "array", "items": {"type": "object"}},
            "page": {"type": "integer"},
            "per_page": {"type": "integer"},
            "total": {"type": "integer"},
            "total_pages": {"type": "integer"}
        }
    }
}

# Add only schemas that don't already exist
for name, schema in schemas.items():
    if name not in spec['components']['schemas']:
        spec['components']['schemas'][name] = schema

# ─────────────────────────────────────────────────────────────────────────────
# 3. Paths
# ─────────────────────────────────────────────────────────────────────────────
PREFIX = "/api/v1/tenant/career-intelligence"
TAGS = [TAG]
SECURITY = [{"bearerAuth": []}]
QUERY_PAGE = {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}}
QUERY_PER_PAGE = {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}
QUERY_PERIOD = {"name": "period", "in": "query", "schema": {"type": "string"}, "description": "Assessment period (e.g., 2026-Q1)"}
Q_EMPLOYEE_ID = {"name": "employee_id", "in": "query", "schema": {"type": "string", "format": "uuid"}}

def path_id():
    return {"name": "id", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}

def path_employee_id():
    return {"name": "employeeId", "in": "path", "required": True, "schema": {"type": "string", "format": "uuid"}}

def std_list_params(extra=None):
    params = [QUERY_PAGE, QUERY_PER_PAGE]
    if extra:
        params.extend(extra)
    return params

def request_body(ref):
    return {"required": True, "content": {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}}

def ok_wrapper(ref=None):
    r = {"200": {"description": "Success"}}
    if ref:
        r["200"]["content"] = {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}
    return r

def created_response(ref=None):
    r = {"201": {"description": "Created"}}
    if ref:
        r["201"]["content"] = {"application/json": {"schema": {"$ref": f"#/components/schemas/{ref}"}}}
    return r

# Build all paths
paths = {}

# ── Talent Maps (7 routes) ──
paths[f"{PREFIX}/talent-maps"] = {
    "get": {
        "tags": TAGS, "summary": "List talent maps",
        "operationId": "ciListTalentMaps",
        "parameters": std_list_params([QUERY_PERIOD, Q_EMPLOYEE_ID]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseCI"),
        "description": "Retrieve a paginated list of talent mapping entries for 9-box assessment. Filter by assessment period (e.g., 2026-Q1) and/or employee."
    },
    "post": {
        "tags": TAGS, "summary": "Create talent map entry",
        "operationId": "ciCreateTalentMap",
        "security": SECURITY,
        "requestBody": request_body("CreateTalentMapRequest"),
        "responses": created_response("TalentMapResponse"),
        "description": "Create a new 9-box talent mapping entry for an employee. Combines performance rating (LOW/MEDIUM/HIGH) and potential rating (LOW/MEDIUM/HIGH) to determine grid position (9-BOX-1 through 9-BOX-9)."
    }
}

paths[f"{PREFIX}/talent-maps/grid"] = {
    "get": {
        "tags": TAGS, "summary": "Get talent map grid (9-box view)",
        "operationId": "ciGetTalentGrid",
        "parameters": [QUERY_PERIOD],
        "security": SECURITY,
        "responses": ok_wrapper("TalentGridResponse"),
        "description": "Get the 9-box talent grid overview for a given period. Returns employee counts per quadrant (9-BOX-1 through 9-BOX-9) with labels and descriptions for talent review purposes."
    }
}

paths[f"{PREFIX}/talent-maps/employee/{{employeeId}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get employee talent profile",
        "operationId": "ciGetEmployeeTalentProfile",
        "parameters": [path_employee_id()],
        "security": SECURITY,
        "responses": ok_wrapper("EmployeeTalentProfileResponse"),
        "description": "Get a comprehensive talent profile for an employee: current 9-box assessment, historical talent map entries, career interests, and recommended next positions."
    }
}

paths[f"{PREFIX}/talent-maps/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get talent map by ID",
        "operationId": "ciGetTalentMap",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("TalentMapResponse"),
        "description": "Get detailed information about a specific talent mapping entry."
    },
    "put": {
        "tags": TAGS, "summary": "Update talent map entry",
        "operationId": "ciUpdateTalentMap",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateTalentMapRequest"),
        "responses": ok_wrapper("TalentMapResponse"),
        "description": "Update an existing talent map entry's performance rating, potential rating, and/or notes. Changes to performance/potential automatically recalculate the grid position."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete talent map entry",
        "operationId": "ciDeleteTalentMap",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Talent map entry deleted"}},
        "description": "Soft-delete a talent map entry."
    }
}

# ── Career Interests (3 routes) ──
paths[f"{PREFIX}/interests"] = {
    "get": {
        "tags": TAGS, "summary": "List career interests",
        "operationId": "ciListCareerInterests",
        "parameters": std_list_params([Q_EMPLOYEE_ID]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseCI"),
        "description": "Retrieve a paginated list of career interests. Filter by employee to see an individual's career aspirations."
    },
    "post": {
        "tags": TAGS, "summary": "Record career interest",
        "operationId": "ciCreateCareerInterest",
        "security": SECURITY,
        "requestBody": request_body("CreateCareerInterestRequest"),
        "responses": created_response("CareerInterestResponse"),
        "description": "Record a career interest for an employee. Supports interest types: LEADERSHIP, SPECIALIST, INTERNATIONAL, and ENTREPRENEUR with optional target position, department, and readiness level."
    }
}

paths[f"{PREFIX}/interests/employee/{{employeeId}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get employee career interests",
        "operationId": "ciGetEmployeeCareerInterests",
        "parameters": [path_employee_id()],
        "security": SECURITY,
        "responses": ok_wrapper("CareerInterestResponse"),
        "description": "Get all active career interests for a specific employee. Used in talent review and career development planning."
    }
}

# ── Career Paths (4 routes) ──
paths[f"{PREFIX}/paths"] = {
    "get": {
        "tags": TAGS, "summary": "List career paths",
        "operationId": "ciListCareerPaths",
        "parameters": std_list_params([]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseCI"),
        "description": "Retrieve a paginated list of defined career paths between position titles. Includes PROMOTION, LATERAL, DEMOTION, and CROSSFUNCTIONAL path types."
    },
    "post": {
        "tags": TAGS, "summary": "Create career path",
        "operationId": "ciCreateCareerPath",
        "security": SECURITY,
        "requestBody": request_body("CreateCareerPathRequest"),
        "responses": created_response("CareerPathResponse"),
        "description": "Define a career path between two position titles. Specifies the path type (PROMOTION/LATERAL/DEMOTION/CROSSFUNCTIONAL), typical tenure, requirements, and required competencies."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete career path",
        "operationId": "ciDeleteCareerPath",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Career path deleted"}},
        "description": "Soft-delete a career path definition."
    }
}

paths[f"{PREFIX}/paths/gap-analysis"] = {
    "get": {
        "tags": TAGS, "summary": "Career gap analysis",
        "operationId": "ciGetGapAnalysis",
        "parameters": [
            {"name": "employee_id", "in": "query", "required": True, "schema": {"type": "string", "format": "uuid"}},
            {"name": "target_title_id", "in": "query", "required": True, "schema": {"type": "string", "format": "uuid"}}
        ],
        "security": SECURITY,
        "responses": ok_wrapper("CareerGapAnalysisResponse"),
        "description": "Analyze the gap between an employee's current qualifications and the requirements of a target position title. Returns matched skills, total requirements, gap percentage, and actionable recommendations (TRAINING/EXPERIENCE/CERTIFICATION) with priority levels."
    }
}

# ── Succession Plans (5 routes) ──
paths[f"{PREFIX}/successions"] = {
    "get": {
        "tags": TAGS, "summary": "List succession plans",
        "operationId": "ciListSuccessionPlans",
        "parameters": std_list_params([]),
        "security": SECURITY,
        "responses": ok_wrapper("PaginatedResponseCI"),
        "description": "Retrieve a paginated list of succession plans. Shows which positions have identified successors with readiness levels (READY_NOW, READY_1YR, READY_2YR, POTENTIAL)."
    },
    "post": {
        "tags": TAGS, "summary": "Create succession plan",
        "operationId": "ciCreateSuccessionPlan",
        "security": SECURITY,
        "requestBody": request_body("CreateSuccessionPlanRequest"),
        "responses": created_response("SuccessionPlanResponse"),
        "description": "Create a succession plan for a key position. Identifies a potential successor with readiness level, priority order, target date, and development plan."
    }
}

paths[f"{PREFIX}/successions/{{id}}"] = {
    "get": {
        "tags": TAGS, "summary": "Get succession plan by ID",
        "operationId": "ciGetSuccessionPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": ok_wrapper("SuccessionPlanResponse"),
        "description": "Get detailed information about a specific succession plan including successor details, readiness level, and development plan."
    },
    "put": {
        "tags": TAGS, "summary": "Update succession plan",
        "operationId": "ciUpdateSuccessionPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "requestBody": request_body("UpdateSuccessionPlanRequest"),
        "responses": ok_wrapper("SuccessionPlanResponse"),
        "description": "Update an existing succession plan's readiness level, priority order, target date, development plan, and/or notes."
    },
    "delete": {
        "tags": TAGS, "summary": "Delete succession plan",
        "operationId": "ciDeleteSuccessionPlan",
        "parameters": [path_id()],
        "security": SECURITY,
        "responses": {"200": {"description": "Succession plan deleted"}},
        "description": "Soft-delete a succession plan."
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# 4. Merge paths into the spec
# ─────────────────────────────────────────────────────────────────────────────
for path, methods in paths.items():
    if path not in spec['paths']:
        spec['paths'][path] = methods
    else:
        for method, endpoint in methods.items():
            if method not in spec['paths'][path]:
                spec['paths'][path][method] = endpoint

# ─────────────────────────────────────────────────────────────────────────────
# 5. Write output
# ─────────────────────────────────────────────────────────────────────────────
with open(OPENAPI_PATH, 'w', encoding='utf-8') as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

# ─────────────────────────────────────────────────────────────────────────────
# 6. Summary
# ─────────────────────────────────────────────────────────────────────────────
ci_paths = {p: list(v.keys()) for p, v in spec['paths'].items() if 'career-intelligence' in p}
total_endpoints = sum(len(m) for m in ci_paths.values())
total_schemas = len(spec['components']['schemas'])
total_openapi_paths = len(spec['paths'])
total_tags = len(spec['tags'])

new_schemas = [n for n in schemas if n in spec['components']['schemas']]

print(f"=== Career Intelligence OpenAPI Results ===")
print(f"Endpoints added: {total_endpoints}")
print(f"Total unique paths: {len(ci_paths)}")
print(f"New schemas added: {len(new_schemas)}")
print(f"")
print(f"UPDATED STATE:")
print(f"Total paths: {total_openapi_paths}")
print(f"Total schemas: {total_schemas}")
print(f"Total tags: {total_tags}")
print(f"")
print(f"New tag: {TAG}")
print(f"")
print(f"=== Path groups ===")
for p in sorted(ci_paths.keys()):
    methods = ", ".join(m.upper() for m in sorted(ci_paths[p]))
    print(f"  {methods} {p}")
