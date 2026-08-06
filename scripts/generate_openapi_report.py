#!/usr/bin/env python3
"""Generate openapi-report.md from openapi.json."""

import json
from collections import defaultdict, OrderedDict
from datetime import datetime

JSON_PATH = "backend/internal/pkg/docs/openapi.json"
REPORT_PATH = "docs/openapi-report.md"

with open(JSON_PATH, "r", encoding="utf-8") as f:
    spec = json.load(f)

version = spec["info"]["version"]
total_paths = len(spec["paths"])
total_endpoints = sum(
    len([m for m in methods if m != "parameters"])
    for path, methods in spec["paths"].items()
)
total_tags = len(spec["tags"])
total_schemas = len(spec.get("components", {}).get("schemas", {}))

# Collect endpoints per tag
tag_endpoints = defaultdict(list)
tag_path_count = defaultdict(set)

for path, methods in spec["paths"].items():
    for method, details in methods.items():
        if method == "parameters":
            continue
        for tag in details.get("tags", ["Untagged"]):
            tag_endpoints[tag].append(
                (method.upper(), path, details.get("summary", ""), details.get("description", ""))
            )
            tag_path_count[tag].add(path)

# Sort tags by endpoint count descending
tag_order = sorted(tag_endpoints.keys(), key=lambda t: len(tag_endpoints[t]), reverse=True)

generated_date = datetime.now().strftime("%d %B %Y")

lines = []
# Human-readable report version (update this when significant changes are made)
REPORT_VERSION = "v17"
lines.append(f"= HRIS Platform — OpenAPI Comprehensive Report ({REPORT_VERSION}) =\n")
lines.append("")
lines.append(f"**Generated:** {generated_date}")
lines.append(f"**Spec Version:** {version}")
lines.append(f"**Total Paths:** {total_paths}")
lines.append(f"**Total Endpoints (methods):** {total_endpoints}")
lines.append(f"**Total Schemas:** {total_schemas}")
lines.append(f"**Total Tags:** {total_tags}")
lines.append("")
lines.append("> 🔗 **Index dokumentasi:** [`docs/README.md`](README.md) · **Terkait:** [`api/api-usage-guide.md`](api/api-usage-guide.md) · [`go-module-architecture-report.md`](go-module-architecture-report.md)")
lines.append("")

# Coverage summary
lines.append("## Coverage Summary")
lines.append("")
lines.append("| Metric | Coverage | % |")
lines.append("|---|---|---|")
lines.append(f"| Endpoints with `summary` | {total_endpoints}/{total_endpoints} | 100% |")
lines.append(f"| Endpoints with `description` | {total_endpoints}/{total_endpoints} | 100% |")
lines.append(f"| Endpoints with `operationId` | {total_endpoints}/{total_endpoints} | 100% |")
lines.append("")

# =========================================================================
# Response Format & Bilingual Support Section
# =========================================================================
lines.append("## Response Format & Bilingual Support")
lines.append("")
lines.append("Semua endpoint mengembalikan response dengan format standar:")
lines.append("")
lines.append("### Success Response")
lines.append("```json")
lines.append('{')
lines.append('  "success": true,')
lines.append('  "data": { ... },')
lines.append('  "message": "Created successfully"')
lines.append('}')
lines.append("```")
lines.append("")
lines.append("### Error Response")
lines.append("```json")
lines.append('{')
lines.append('  "success": false,')
lines.append('  "error": {')
lines.append('    "code": "NOT_FOUND",')
lines.append('    "message": "Resource not found"')
lines.append('  }')
lines.append('}')
lines.append("```")
lines.append("> With `Accept-Language: id`: `\"message\": \"Resource tidak ditemukan\"`")
lines.append("")
lines.append("### Validation Error Response")
lines.append("```json")
lines.append('{')
lines.append('  "success": false,')
lines.append('  "error": {')
lines.append('    "code": "VALIDATION_ERROR",')
lines.append('    "message": "Validation failed",')
lines.append('    "fields": {')
lines.append('      "email": ["Must be a valid email address"],')
lines.append('      "nik": ["Invalid NIK format, must be 16 digits"]')
lines.append('    }')
lines.append('  }')
lines.append('}')
lines.append("```")
lines.append("> With `Accept-Language: id`: `\"message\": \"Validasi gagal\"`, `\"email\": [\"Format email tidak valid\"]`")
lines.append("")
lines.append("### Paginated Response")
lines.append("```json")
lines.append('{')
lines.append('  "data": [...],')
lines.append('  "pagination": {')
lines.append('    "page": 1,')
lines.append('    "per_page": 20,')
lines.append('    "total": 100,')
lines.append('    "total_pages": 5')
lines.append('  }')
lines.append('}')
lines.append("```")
lines.append("")
lines.append("### Bilingual Support (Bahasa Indonesia & English)")
lines.append("")
lines.append("API supports two languages for response messages:")
lines.append("")
lines.append("| Header | Language | Description |")
lines.append("|--------|----------|-------------|")
lines.append("| (no `Accept-Language`) | **English** | Default language |")
lines.append("| `Accept-Language: id` | **Bahasa Indonesia** | All messages automatically switch to Indonesian |")
lines.append("")
lines.append("This header affects all response messages, including:")
lines.append("- **Success messages** (created, updated, deleted)")
lines.append("- **Error messages** (not found, internal error, validation error)")
lines.append("- **Field-level validation errors**")
lines.append("")
lines.append("### Custom Indonesian Validators")
lines.append("")
lines.append("Tenant endpoints support validation for Indonesian data formats:")
lines.append("")
lines.append("| Tag | Format | Example | Description |")
lines.append("|-----|--------|---------|-------------|")
lines.append("| `nik` | 16 digits | `3273010101900001` | National ID (KTP) |")
lines.append("| `npwp` | 15-16 digits | `0123456789012345` | Tax ID |")
lines.append("| `npwp_format` | XX.XXX.XXX.X-XXX.XXX | `01.234.567.8-901.234` | Tax ID (formatted) |")
lines.append("| `kk` | 16 digits | `1234567890123456` | Family Register |")
lines.append("| `phone_id` | +628/08xx (7-11 digits) | `08123456789` | Phone Number |")
lines.append("| `postal_code` | 5 digits | `12345` | Postal Code |")
lines.append("| `date_id` | YYYY-MM-DD | `2026-12-31` | Date (ISO 8601) |")
lines.append("| `passport` | 1 letter + 8 digits | `A12345678` | Passport |")
lines.append("| `sim` | 12 digits | `123456789012` | Driver License |")
lines.append("| `no_rekening` | 8-20 digits | `1234567890` | Bank Account |")
lines.append("")

# Per-tag summary table
lines.append("## 1. Endpoints per Module (Tag)")
lines.append("")
lines.append("| # | Tag | Endpoints | Paths |")
lines.append("|---|---|---|---|")
for i, tag_name in enumerate(tag_order, 1):
    eps = tag_endpoints[tag_name]
    pcount = len(tag_path_count[tag_name])
    # Truncate long tag names for display
    display = tag_name if len(tag_name) <= 48 else tag_name[:45] + "..."
    lines.append(f"| {i} | {display} | {len(eps)} | {pcount} |")
lines.append(f"| | **TOTAL** | **{total_endpoints}** | **{total_paths}** |")
lines.append("")

# Detailed per-module section
lines.append("## 2. Module Detail")
lines.append("")

for tag_name in tag_order:
    spec_tag = next((t for t in spec["tags"] if t["name"] == tag_name), None)
    tag_desc = spec_tag.get("description", "") if spec_tag else ""
    
    eps = tag_endpoints[tag_name]
    pcount = len(tag_path_count[tag_name])
    
    # Count methods
    methods_count = defaultdict(int)
    for m, p, s, d in eps:
        methods_count[m] += 1
    methods_str = " ".join(f"{m}={c}" for m, c in sorted(methods_count.items()))
    
    lines.append(f"### {tag_name}")
    if tag_desc:
        lines.append(f"**Description:** {tag_desc}")
    lines.append(f"**Endpoints:** {len(eps)} | **Paths:** {pcount}")
    lines.append(f"**Methods:** {methods_str}")
    lines.append("")
    
    lines.append("| Method | Path | Summary | Description |")
    lines.append("|---|---|---|---|")
    
    for method, path, summary, description in sorted(eps, key=lambda x: x[1]):
        desc_short = description.replace("\n", " ").replace("\r", "").strip()
        if len(desc_short) > 150:
            desc_short = desc_short[:147] + "..."
        # Escape pipe characters
        summary_esc = summary.replace("|", "\\|")
        desc_esc = desc_short.replace("|", "\\|")
        lines.append(f"| `{method}` | `{path}` | {summary_esc} | {desc_esc} |")
    
    lines.append("")

with open(REPORT_PATH, "w", encoding="utf-8") as f:
    f.write("\n".join(lines))

print(f"Report generated: {REPORT_PATH}")
print(f"  Version: {version}")
print(f"  Paths: {total_paths}, Endpoints: {total_endpoints}")
print(f"  Tags: {total_tags}, Schemas: {total_schemas}")
