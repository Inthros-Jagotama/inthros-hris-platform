#!/usr/bin/env python3
"""Inject organization-summaries CRUD endpoints + schemas into openapi.json."""

import json
import re

JSON_PATH = "backend/internal/pkg/docs/openapi.json"

with open(JSON_PATH, "r") as f:
    content = f.read()

spec = json.loads(content)

# ── 1. Add paths ──
# Insert after the last settings/salary-grades path entry (before components)
summary_paths = {
    "/api/v1/tenant/organization-summaries": {
        "post": {
            "tags": ["Tenant: Organizations"],
            "summary": "Create organization summary",
            "operationId": "createOrganizationSummary",
            "security": [{"bearerAuth": []}],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/CreateOrganizationSummaryRequest"}
                    }
                }
            },
            "responses": {
                "201": {
                    "description": "Organization summary created",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {"$ref": "#/components/schemas/OrganizationSummaryResponse"},
                                    "message": {"type": "string", "example": "Created successfully"}
                                }
                            }
                        }
                    }
                },
                "400": {"description": "Validation error"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Create a new organization summary record. Requires code, decree_no, and decree_date."
        },
        "get": {
            "tags": ["Tenant: Organizations"],
            "summary": "List organization summaries (paginated)",
            "operationId": "listOrganizationSummaries",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "page", "in": "query", "schema": {"type": "integer", "default": 1}},
                {"name": "per_page", "in": "query", "schema": {"type": "integer", "default": 20}}
            ],
            "responses": {
                "200": {
                    "description": "Paginated list of organization summaries",
                    "content": {
                        "application/json": {
                            "schema": {"$ref": "#/components/schemas/PaginatedSummaryResponse"}
                        }
                    }
                },
                "401": {"description": "Unauthorized"}
            },
            "description": "Retrieve a paginated list of organization summaries sorted by creation date (descending)."
        }
    },
    "/api/v1/tenant/organization-summaries/stats": {
        "get": {
            "tags": ["Tenant: Organizations"],
            "summary": "Get organization summary statistics",
            "operationId": "getOrganizationSummaryStats",
            "security": [{"bearerAuth": []}],
            "responses": {
                "200": {
                    "description": "Summary statistics",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {
                                        "type": "object",
                                        "properties": {
                                            "total_summaries": {"type": "integer"},
                                            "total_orgs": {"type": "integer"},
                                            "max_depth": {"type": "integer"},
                                            "updated_at": {"type": "string", "format": "date-time"}
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "401": {"description": "Unauthorized"}
            },
            "description": "Get aggregate statistics about organization summaries and organizations."
        }
    },
    "/api/v1/tenant/organization-summaries/{id}": {
        "get": {
            "tags": ["Tenant: Organizations"],
            "summary": "Get organization summary by ID",
            "operationId": "getOrganizationSummaryById",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}}
            ],
            "responses": {
                "200": {
                    "description": "Organization summary details",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {"$ref": "#/components/schemas/OrganizationSummaryResponse"}
                                }
                            }
                        }
                    }
                },
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Get detailed information about a specific organization summary by its ID."
        },
        "put": {
            "tags": ["Tenant: Organizations"],
            "summary": "Update organization summary",
            "operationId": "updateOrganizationSummary",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}}
            ],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {"$ref": "#/components/schemas/UpdateOrganizationSummaryRequest"}
                    }
                }
            },
            "responses": {
                "200": {
                    "description": "Organization summary updated",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "data": {"$ref": "#/components/schemas/OrganizationSummaryResponse"}
                                }
                            }
                        }
                    }
                },
                "400": {"description": "Validation error"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Update an existing organization summary's code, decree_no, decree_date, or status."
        },
        "delete": {
            "tags": ["Tenant: Organizations"],
            "summary": "Delete organization summary",
            "operationId": "deleteOrganizationSummary",
            "security": [{"bearerAuth": []}],
            "parameters": [
                {"name": "id", "in": "path", "required": True, "schema": {"type": "string"}}
            ],
            "responses": {
                "200": {
                    "description": "Organization summary deleted",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "object",
                                "properties": {
                                    "success": {"type": "boolean", "example": True},
                                    "message": {"type": "string", "example": "Deleted successfully"}
                                }
                            }
                        }
                    }
                },
                "409": {"description": "Conflict — summary has organizations attached"},
                "404": {"description": "Not found"},
                "401": {"description": "Unauthorized"}
            },
            "description": "Soft-delete an organization summary. Cannot delete if organizations are still attached to this summary."
        }
    }
}

# Insert paths before the closing "}" of paths (before components)
# Find the insertion point - before the last SalaryGrade DELETE entry
spec["paths"].update(summary_paths)

# ── 2. Add schemas ──
summary_schemas = {
    "OrganizationSummaryResponse": {
        "type": "object",
        "properties": {
            "id": {"type": "string"},
            "code": {"type": "string"},
            "decree_no": {"type": "string"},
            "decree_date": {"type": "string"},
            "status": {"type": "string"},
            "org_count": {"type": "integer"},
            "created_at": {"type": "string", "format": "date-time"},
            "updated_at": {"type": "string", "format": "date-time"}
        }
    },
    "PaginatedSummaryResponse": {
        "type": "object",
        "properties": {
            "success": {"type": "boolean"},
            "data": {
                "type": "array",
                "items": {"$ref": "#/components/schemas/OrganizationSummaryResponse"}
            },
            "page": {"type": "integer"},
            "per_page": {"type": "integer"},
            "total": {"type": "integer"},
            "total_pages": {"type": "integer"}
        }
    },
    "CreateOrganizationSummaryRequest": {
        "type": "object",
        "required": ["code", "decree_no", "decree_date"],
        "properties": {
            "code": {"type": "string", "maxLength": 7},
            "decree_no": {"type": "string", "maxLength": 20},
            "decree_date": {"type": "string", "format": "date"}
        }
    },
    "UpdateOrganizationSummaryRequest": {
        "type": "object",
        "properties": {
            "code": {"type": "string", "maxLength": 7},
            "decree_no": {"type": "string", "maxLength": 20},
            "decree_date": {"type": "string", "format": "date"},
            "status": {"type": "string"}
        }
    }
}

spec["components"]["schemas"].update(summary_schemas)

# ── 3. Write back ──
with open(JSON_PATH, "w") as f:
    json.dump(spec, f, indent=2)

print("[OK] organization-summaries endpoints and schemas injected successfully!")
print(f"   Paths added: {len(summary_paths)}")
print(f"   Endpoints added: 6")
print(f"   Schemas added: {len(summary_schemas)}")
