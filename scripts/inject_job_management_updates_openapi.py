#!/usr/bin/env python3
"""Inject new job-management endpoints into openapi.json (05 Aug 2026):
- GET  /job-management/values/tree
- GET/PUT /job-management/values/clusters/{type}
- POST/GET/GET/PUT/DELETE /job-management/relationships/{id}/details(/{detailId})
"""

import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

TAG = "Tenant: Job Management"
SEC = [{"bearerAuth": []}]


def param(name, where="path", required=True, schema=None):
    return {"name": name, "in": where, "required": required, "schema": schema or {"type": "string", "format": "uuid"}}


def qparam(name, schema=None):
    return param(name, where="query", required=False, schema=schema or {"type": "string"})


def resp(desc):
    return {"200": {"description": desc}, "400": {"description": "Validation error"}, "401": {"description": "Unauthorized"}, "500": {"description": "Internal server error"}}


with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

paths = spec["paths"]
added = 0

tree_path = "/api/v1/tenant/job-management/values/tree"
if tree_path not in paths:
    paths[tree_path] = {
        "get": {
            "tags": [TAG],
            "summary": "Get job values tree",
            "operationId": "listJobValuesTree",
            "security": SEC,
            "parameters": [qparam("type_group")],
            "responses": resp("Job values tree (type_group → types → level options)"),
            "description": "Mengembalikan hierarki type_group → daftar tipe (label = description_group) → options per tipe (level + deskripsi) dengan urutan grup tetap. Dipakai form potensi (filter type_group → multi-select tipe).",
        }
    }
    added += 1

clusters_path = "/api/v1/tenant/job-management/values/clusters/{type}"
if clusters_path not in paths:
    paths[clusters_path] = {
        "get": {
            "tags": [TAG],
            "summary": "List cluster mapping for job value type",
            "operationId": "listJobValueClusters",
            "security": SEC,
            "parameters": [param("type")],
            "responses": resp("Cluster mapping (technical/managerial → competencies)"),
            "description": "Mapping tipe technical/managerial → cluster kompetensi dari tabel job_management_value_clusters.",
        },
        "put": {
            "tags": [TAG],
            "summary": "Update cluster mapping for job value type",
            "operationId": "updateJobValueClusters",
            "security": SEC,
            "parameters": [param("type")],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "clusters": {"type": "array", "items": {"type": "string"}}
                            },
                        }
                    }
                },
            },
            "responses": resp("Cluster mapping updated"),
            "description": "Simpan mapping cluster untuk tipe tertentu (technical/managerial).",
        },
    }
    added += 2

# Relationship details (nested under relationships/{id}/details)
details_base = "/api/v1/tenant/job-management/relationships/{id}/details"
details_item = "/api/v1/tenant/job-management/relationships/{id}/details/{detailId}"
if details_base not in paths:
    paths[details_base] = {
        "get": {
            "tags": [TAG],
            "summary": "List relationship details",
            "operationId": "listJobRelationshipDetails",
            "security": SEC,
            "parameters": [param("id")],
            "responses": resp("List of relationship details"),
            "description": "Detail banyak-per-relationship (migration 048).",
        },
        "post": {
            "tags": [TAG],
            "summary": "Create relationship detail",
            "operationId": "createJobRelationshipDetail",
            "security": SEC,
            "parameters": [param("id")],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "work_relations_id": {"type": "string", "format": "uuid"},
                                "activity_in_connection": {"type": "string"},
                            },
                        }
                    }
                },
            },
            "responses": resp("Relationship detail created"),
            "description": "Tambah detail hubungan kerja (work relations + activity in connection).",
        },
    }
    added += 2
if details_item not in paths:
    paths[details_item] = {
        "get": {
            "tags": [TAG],
            "summary": "Get relationship detail by ID",
            "operationId": "getJobRelationshipDetail",
            "security": SEC,
            "parameters": [param("id"), param("detailId")],
            "responses": resp("Relationship detail"),
            "description": "Ambil satu detail hubungan kerja.",
        },
        "put": {
            "tags": [TAG],
            "summary": "Update relationship detail",
            "operationId": "updateJobRelationshipDetail",
            "security": SEC,
            "parameters": [param("id"), param("detailId")],
            "requestBody": {
                "required": True,
                "content": {
                    "application/json": {
                        "schema": {
                            "type": "object",
                            "properties": {
                                "work_relations_id": {"type": "string", "format": "uuid"},
                                "activity_in_connection": {"type": "string"},
                            },
                        }
                    }
                },
            },
            "responses": resp("Relationship detail updated"),
            "description": "Update satu detail hubungan kerja.",
        },
        "delete": {
            "tags": [TAG],
            "summary": "Delete relationship detail",
            "operationId": "deleteJobRelationshipDetail",
            "security": SEC,
            "parameters": [param("id"), param("detailId")],
            "responses": resp("Relationship detail deleted"),
            "description": "Hapus satu detail hubungan kerja.",
        },
    }
    added += 3

# Update scores PUT description to reflect auto-recalc
scores_path = "/api/v1/tenant/job-management/scores/org/{orgId}"
if scores_path in paths and "put" in paths[scores_path]:
    paths[scores_path]["put"]["description"] = (
        "Hitung ulang skor otomatis (body kosong) lalu simpan. Menghasilkan components, "
        "sub_component_points, is_complete, completed_at."
    )

with open(JSON_PATH, "w", encoding="utf-8") as f:
    json.dump(spec, f, indent=2, ensure_ascii=False)

print(f"Injected {added} new endpoints into {JSON_PATH}")
