"""
Script to add Package Management endpoints, schemas, and tag to openapi.json.
Reads the existing JSON, injects the new content at the correct positions, and writes back.
"""

import json
from pathlib import Path

OPENAPI_PATH = Path(__file__).resolve().parent.parent / "internal/pkg/docs/openapi.json"

with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

# =========================================================================
# 1. Add Tag: "Platform: Packages" after "Platform: Licenses"
# =========================================================================
new_tag = {
    "name": "Platform: Packages",
    "description": "Package management — bundle tenant modules with pricing, dependency validation, and publishing"
}

# Cari dan hapus tag Package yang sudah ada (untuk menghindari duplikat)
tags = spec["tags"]
tags[:] = [t for t in tags if t["name"] != "Platform: Packages"]

# Insert after Licenses tag
licenses_idx = None
for i, t in enumerate(tags):
    if t["name"] == "Platform: Licenses":
        licenses_idx = i
        break

if licenses_idx is not None:
    tags.insert(licenses_idx + 1, new_tag)
    print(f"Added 'Platform: Packages' tag at index {licenses_idx + 1}")
else:
    tags.append(new_tag)
    print("Added 'Platform: Packages' tag at end (Licenses tag not found)")

# =========================================================================
# 2. Add Paths (insert after Licenses paths, before Monitoring paths)
# =========================================================================
package_paths = {
    # --- Public (no auth) ---
    "/api/v1/public/packages": {
        "get": {
            "tags": ["Platform: Packages"],
            "summary": "List published packages (public)",
            "operationId": "listPublishedPackages",
            "responses": {
                "200": {
                    "description": "List of published packages with modules",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/components/schemas/PublicPackageResponse"
                                }
                            }
                        }
                    }
                }
            },
            "description": "Retrieve a list of published packages for public display. No authentication required. Returns package name, description, price, and included module list."
        }
    },
    # --- Admin: List ---
    "/api/v1/platform/packages": {
        "get": {
            "tags": ["Platform: Packages"],
            "summary": "List all packages (admin)",
            "operationId": "listPackages",
            "x-rbac": {
                "permission": "package.view",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "page",
                    "in": "query",
                    "schema": {"type": "integer", "default": 1}
                },
                {
                    "name": "per_page",
                    "in": "query",
                    "schema": {"type": "integer", "default": 20}
                }
            ],
            "responses": {
                "200": {
                    "description": "Paginated list of packages"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                }
            },
            "description": "Retrieve a paginated list of all packages. Includes draft, published, and archived packages with their module associations."
        },
        "post": {
            "tags": ["Platform: Packages"],
            "summary": "Create a new package",
            "operationId": "createPackage",
            "x-rbac": {
                "permission": "package.create",
                "roles": ["super_admin"]
            },
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "$ref": "#/components/schemas/CreatePackageRequest"
                        }
                    }
                }
            },
            "responses": {
                "201": {
                    "description": "Package created"
                },
                "400": {
                    "description": "Validation error or dependency validation failed"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                },
                "409": {
                    "description": "Package slug already exists"
                }
            },
            "description": "Create a new package that bundles tenant modules with pricing. Validates module dependencies before creation — all required dependencies (depends_on) must be included in the package."
        }
    },
    # --- Admin: Get by ID ---
    "/api/v1/platform/packages/{id}": {
        "get": {
            "tags": ["Platform: Packages"],
            "summary": "Get package by ID",
            "operationId": "getPackage",
            "x-rbac": {
                "permission": "package.view",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "responses": {
                "200": {
                    "description": "Package details with modules",
                    "content": {
                        "application/json": {
                            "schema": {
                                "$ref": "#/components/schemas/PackageResponse"
                            }
                        }
                    }
                },
                "403": {
                    "description": "Forbidden"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Get detailed information about a specific package including its modules, pricing, status, and sort order."
        },
        "put": {
            "tags": ["Platform: Packages"],
            "summary": "Update package",
            "operationId": "updatePackage",
            "x-rbac": {
                "permission": "package.update",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "requestBody": {
                "content": {
                    "application/json": {
                        "schema": {
                            "$ref": "#/components/schemas/UpdatePackageRequest"
                        }
                    }
                }
            },
            "responses": {
                "200": {
                    "description": "Package updated"
                },
                "400": {
                    "description": "Validation error or dependency validation failed"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                },
                "409": {
                    "description": "Conflict (e.g., slug already exists)"
                }
            },
            "description": "Update a package's metadata, pricing, status, or module associations. When updating modules, dependency validation is re-run."
        },
        "delete": {
            "tags": ["Platform: Packages"],
            "summary": "Delete package (soft-delete)",
            "operationId": "deletePackage",
            "x-rbac": {
                "permission": "package.delete",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "responses": {
                "200": {
                    "description": "Package soft-deleted"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Soft-delete a package. Sets the deleted_at timestamp and removes module associations. The package is hidden from standard queries."
        }
    },
    # --- Admin: Publish ---
    "/api/v1/platform/packages/{id}/publish": {
        "post": {
            "tags": ["Platform: Packages"],
            "summary": "Publish package",
            "operationId": "publishPackage",
            "x-rbac": {
                "permission": "package.update",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "responses": {
                "200": {
                    "description": "Package published — status=published, is_public=true"
                },
                "400": {
                    "description": "Validation error: no modules or missing dependencies"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Publish a package to make it visible on public endpoints. Validates that the package has at least one module and all module dependencies are fulfilled."
        }
    },
    # --- Admin: Unpublish ---
    "/api/v1/platform/packages/{id}/unpublish": {
        "post": {
            "tags": ["Platform: Packages"],
            "summary": "Unpublish package",
            "operationId": "unpublishPackage",
            "x-rbac": {
                "permission": "package.update",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "responses": {
                "200": {
                    "description": "Package unpublished — status=draft"
                },
                "403": {
                    "description": "Forbidden: only super_admin"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Unpublish a package. Sets status back to 'draft' and removes it from public endpoints."
        }
    },
    # --- Admin: Validate dependencies ---
    "/api/v1/platform/packages/{id}/validate": {
        "get": {
            "tags": ["Platform: Packages"],
            "summary": "Validate package module dependencies",
            "operationId": "validatePackageDependencies",
            "x-rbac": {
                "permission": "package.view",
                "roles": ["super_admin"]
            },
            "parameters": [
                {
                    "name": "id",
                    "in": "path",
                    "required": True,
                    "schema": {
                        "type": "string",
                        "format": "uuid"
                    }
                }
            ],
            "responses": {
                "200": {
                    "description": "Dependency validation results",
                    "content": {
                        "application/json": {
                            "schema": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/components/schemas/ModuleDependency"
                                }
                            }
                        }
                    }
                },
                "403": {
                    "description": "Forbidden"
                },
                "404": {
                    "description": "Package not found"
                }
            },
            "description": "Validate that all module dependencies within a package are fulfilled. Returns a detailed report showing each module's dependency status (resolved/unresolved)."
        }
    }
}

# Insert paths after the last License path, before Monitoring
paths = spec["paths"]
new_paths = {}
inserted = False
for key, value in paths.items():
    new_paths[key] = value
    if key == "/api/v1/platform/licenses/{id}":
        # Insert package paths after Licenses
        new_paths.update(package_paths)
        inserted = True
        print(f"Inserted package paths after Licenses (before Monitoring)")
    elif not inserted and key.startswith("/api/v1/platform/monitoring"):
        # If we somehow missed the Licenses insertion point, insert before Monitoring
        new_paths.update(package_paths)
        inserted = True
        print(f"Inserted package paths before Monitoring (fallback)")

# Fallback: if still not inserted, append at end
if not inserted:
    new_paths.update(package_paths)
    print("Appended package paths at end (fallback)")

spec["paths"] = new_paths

# =========================================================================
# 3. Add Schemas
# =========================================================================
package_schemas = {
    "CreatePackageRequest": {
        "type": "object",
        "required": ["name", "slug"],
        "properties": {
            "name": {
                "type": "string",
                "minLength": 3,
                "maxLength": 255,
                "description": "Package display name"
            },
            "slug": {
                "type": "string",
                "minLength": 2,
                "maxLength": 100,
                "description": "Unique URL-friendly identifier"
            },
            "description": {
                "type": "string",
                "description": "Package description"
            },
            "price": {
                "type": "number",
                "minimum": 0,
                "description": "Package price"
            },
            "sort_order": {
                "type": "integer",
                "description": "Display order (lower = first)"
            },
            "modules": {
                "type": "array",
                "items": {
                    "$ref": "#/components/schemas/PackageModuleInput"
                },
                "description": "Modules included in this package"
            }
        }
    },
    "PackageModuleInput": {
        "type": "object",
        "required": ["module_id"],
        "properties": {
            "module_id": {
                "type": "string",
                "format": "uuid",
                "description": "ID of the module from module registry"
            },
            "is_mandatory": {
                "type": "boolean",
                "description": "Whether this module is mandatory (cannot be disabled)",
                "default": False
            },
            "sort_order": {
                "type": "integer",
                "description": "Display order within the package",
                "default": 0
            }
        }
    },
    "UpdatePackageRequest": {
        "type": "object",
        "properties": {
            "name": {
                "type": "string",
                "minLength": 3,
                "description": "Package display name"
            },
            "slug": {
                "type": "string",
                "minLength": 2,
                "description": "Unique URL-friendly identifier"
            },
            "description": {
                "type": "string",
                "description": "Package description"
            },
            "price": {
                "type": "number",
                "minimum": 0,
                "description": "Package price"
            },
            "sort_order": {
                "type": "integer",
                "description": "Display order"
            },
            "status": {
                "type": "string",
                "enum": ["draft", "published", "archived"],
                "description": "Package status"
            },
            "is_public": {
                "type": "boolean",
                "description": "Whether visible on public endpoints"
            },
            "modules": {
                "type": "array",
                "items": {
                    "$ref": "#/components/schemas/PackageModuleInput"
                },
                "description": "Replace all module associations"
            }
        }
    },
    "PackageResponse": {
        "type": "object",
        "properties": {
            "id": {
                "type": "string",
                "format": "uuid"
            },
            "name": {
                "type": "string"
            },
            "slug": {
                "type": "string"
            },
            "description": {
                "type": "string"
            },
            "price": {
                "type": "number"
            },
            "status": {
                "type": "string",
                "enum": ["draft", "published", "archived"]
            },
            "is_public": {
                "type": "boolean"
            },
            "sort_order": {
                "type": "integer"
            },
            "module_count": {
                "type": "integer",
                "description": "Number of modules in this package"
            },
            "modules": {
                "type": "array",
                "items": {
                    "$ref": "#/components/schemas/PackageModuleResponse"
                }
            },
            "created_at": {
                "type": "string",
                "format": "date-time"
            },
            "updated_at": {
                "type": "string",
                "format": "date-time"
            }
        }
    },
    "PackageModuleResponse": {
        "type": "object",
        "properties": {
            "module_id": {
                "type": "string",
                "format": "uuid"
            },
            "module_name": {
                "type": "string"
            },
            "module_slug": {
                "type": "string"
            },
            "is_mandatory": {
                "type": "boolean"
            },
            "sort_order": {
                "type": "integer"
            }
        }
    },
    "PublicPackageResponse": {
        "type": "object",
        "properties": {
            "id": {
                "type": "string",
                "format": "uuid"
            },
            "name": {
                "type": "string"
            },
            "slug": {
                "type": "string"
            },
            "description": {
                "type": "string"
            },
            "price": {
                "type": "number"
            },
            "sort_order": {
                "type": "integer"
            },
            "module_count": {
                "type": "integer"
            },
            "modules": {
                "type": "array",
                "items": {
                    "$ref": "#/components/schemas/PackageModuleResponse"
                }
            }
        }
    },
    "ModuleDependency": {
        "type": "object",
        "properties": {
            "module_id": {
                "type": "string",
                "format": "uuid"
            },
            "module_name": {
                "type": "string"
            },
            "depends_on": {
                "type": "string",
                "description": "Dependency description (e.g. 'needs: employee,organization' or '(none)')"
            },
            "resolved": {
                "type": "boolean",
                "description": "Whether all dependencies are satisfied within this package"
            }
        }
    }
}

# Schemas are inside components
spec.setdefault("components", {})
spec["components"]["schemas"].update(package_schemas)
print(f"Added {len(package_schemas)} package schemas")

# =========================================================================
# 4. Write back
# =========================================================================
with open(OPENAPI_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print(f"\nDone! Updated {OPENAPI_PATH}")
print(f"Tags: {len(spec['tags'])}")
print(f"Paths: {len(spec['paths'])}")
endpoints = sum(
    len([m for m in methods if m != 'parameters'])
    for path, methods in spec['paths'].items()
)
print(f"Endpoints: {endpoints}")
print(f"Schemas: {len(spec['components']['schemas'])}")
