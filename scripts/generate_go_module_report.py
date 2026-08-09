#!/usr/bin/env python3
"""Regenerate docs/go-module-architecture-report.md from the Go source tree.

Scans backend/internal/modules, backend/internal/platform and
backend/internal/pkg and computes, per module:

  - Entities        : `type X struct` declared in model.go
  - Service Methods : `func (s *Service)` receiver methods (all non-test files)
  - Repo Methods    : `func (r *Repository)` receiver methods
  - Handler Funcs   : `func (h *Handler)` receiver methods
  - Route Regs      : `.GET(/.POST(/.PUT(/.PATCH(/.DELETE(` calls in routes.go
  - Tests           : `func Test` in *_test.go (with repo/service/handler split)

Output mirrors the hand-maintained report structure (Sections 1-7).

Usage:
    python scripts/generate_go_module_report.py
"""
import datetime
import os
import re

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
BACKEND = os.path.join(PROJECT_ROOT, "backend")
OUT_PATH = os.path.join(PROJECT_ROOT, "docs", "go-module-architecture-report.md")

MODULES_DIR = os.path.join(BACKEND, "internal", "modules")
PLATFORM_DIR = os.path.join(BACKEND, "internal", "platform")
PKG_DIR = os.path.join(BACKEND, "internal", "pkg")

RE_STRUCT = re.compile(r"^type\s+(\w+)\s+struct\s*\{", re.M)
RE_METHOD = re.compile(r"^func\s*\(\s*(\w+)\s*\*\s*(\w+)\s*\)\s*(\w+)\s*\(", re.M)
RE_ROUTE = re.compile(r"\.(?:GET|POST|PUT|PATCH|DELETE)\(")
RE_TEST = re.compile(r"^func\s+Test\w+\s*\(", re.M)


def walk_go_files(directory):
    for root, dirs, files in os.walk(directory):
        dirs[:] = [d for d in dirs if d not in (".git", "vendor")]
        for fn in sorted(files):
            if fn.endswith(".go"):
                yield os.path.join(root, fn)


def read(path):
    with open(path, "r", encoding="utf-8", errors="ignore") as f:
        return f.read()


def count(regex, text):
    return len(regex.findall(text))


def module_stats(mod_dir, name):
    """Return dict of stats for one module directory."""
    entities = []
    service, repo, handler, routes, tests = [], [], [], 0, 0
    repo_tests = service_tests = handler_tests = other_tests = 0
    test_files = []

    for fp in walk_go_files(mod_dir):
        base = os.path.basename(fp)
        if base.endswith("_test.go"):
            txt = read(fp)
            n = count(RE_TEST, txt)
            tests += n
            test_files.append((os.path.relpath(fp, BACKEND), n))
            if "repository_test" in base:
                repo_tests += n
            elif "service_test" in base:
                service_tests += n
            elif "handler_test" in base:
                handler_tests += n
            else:
                other_tests += n
            continue

        txt = read(fp)
        if base == "model.go":
            entities.extend(m.group(1) for m in RE_STRUCT.finditer(txt))
        for m in RE_METHOD.finditer(txt):
            recv, typ, fn = m.group(1), m.group(2), m.group(3)
            if typ == "Service":
                service.append(fn)
            elif typ == "Repository":
                repo.append(fn)
            elif typ == "Handler":
                handler.append(fn)
        if base == "routes.go":
            routes += count(RE_ROUTE, txt)

    return {
        "name": name,
        "entities": entities,
        "service": service,
        "repo": repo,
        "handler": handler,
        "routes": routes,
        "tests": tests,
        "repo_tests": repo_tests,
        "service_tests": service_tests,
        "handler_tests": handler_tests,
        "other_tests": other_tests,
        "test_files": sorted(test_files, key=lambda x: x[0]),
    }


def pkg_stats(pkg_dir, name):
    """Return dict of stats for one shared kernel package."""
    go_files = []
    tests = 0
    for fp in walk_go_files(pkg_dir):
        txt = read(fp)
        go_files.append(os.path.basename(fp))
        tests += count(RE_TEST, txt)
    desc = pkg_description(pkg_dir, name)
    return {"name": name, "files": len(go_files), "tests": tests, "desc": desc}


PKG_DESC = {
    "auth": "JWT authentication",
    "authctx": "Auth context helpers (user/company from gin context)",
    "authz": "Casbin RBAC authorization",
    "cache": "Distributed cache + Pub/Sub",
    "config": "Viper configuration loader",
    "crypto": "AES-256-GCM encryption",
    "database": "Multi-tenant DB connection manager",
    "docs": "OpenAPI/Scalar documentation",
    "driver": "DB driver detection",
    "errors": "Shared error helpers",
    "httputil": "Bilingual response helpers (SuccessJSON, CreatedJSON, ErrorJSON, NotFound) + locale message catalog (80+ EN/ID pairs) + custom Indonesian validators (NIK, NPWP, KK, Passport, SIM, No Rekening)",
    "logger": "Structured logging",
    "mailer": "Email sending (SMTP/template)",
    "middleware": "HTTP middleware: Auth, CORS, Logger, Recovery, Tenant, Localize (auto-detect Accept-Language)",
    "migrator": "Database migration engine",
    "module": "Module SDK",
    "onpremise": "On-premise license enforcement",
    "router": "Router setup",
    "telemetry": "Telemetry/metrics",
    "tenant": "Tenant resolution helpers",
    "tenantseed": "Tenant seed data (nationalities, competencies, RBAC)",
    "validator": "Validator helpers",
}


def pkg_description(pkg_dir, name):
    if name in PKG_DESC:
        return PKG_DESC[name]
    for fp in walk_go_files(pkg_dir):
        txt = read(fp)
        m = re.search(r"^// Package \w+ (.*)$", txt, re.M)
        if m:
            return m.group(1).strip()
    return name


def fmt_table(stats):
    header = "| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |"
    sep = "|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|"
    rows = []
    for s in sorted(stats, key=lambda x: x["name"]):
        rows.append(
            f"| {s['name']} | {len(s['entities'])} | {len(s['service'])} | {len(s['repo'])} | "
            f"{len(s['handler'])} | {s['routes']} | {s['tests']} |"
        )
    totals = {
        "entities": sum(len(s["entities"]) for s in stats),
        "service": sum(len(s["service"]) for s in stats),
        "repo": sum(len(s["repo"]) for s in stats),
        "handler": sum(len(s["handler"]) for s in stats),
        "routes": sum(s["routes"] for s in stats),
        "tests": sum(s["tests"] for s in stats),
    }
    rows.append(
        f"| **TOTAL** | **{totals['entities']}** | **{totals['service']}** | **{totals['repo']}** | "
        f"**{totals['handler']}** | **{totals['routes']}** | **{totals['tests']}** |"
    )
    return "\n".join([header, sep] + rows), totals


def section4(stats):
    lines = ["## SECTION 4: ENTITY DETAIL PER MODULE"]
    for s in sorted(stats, key=lambda x: x["name"]):
        lines.append(f"\n### {s['name']}")
        for e in s["entities"]:
            lines.append(f"- {e}")
    return "\n".join(lines)


def section5(stats):
    lines = ["## SECTION 5: SERVICE METHOD DETAIL"]
    for s in sorted(stats, key=lambda x: x["name"]):
        if not s["service"]:
            continue
        lines.append(f"\n### {s['name']}.Service — {len(s['service'])} methods")
        for m in s["service"]:
            lines.append(f"- `{m}()`")
    return "\n".join(lines)


def section7(all_test_files):
    lines = ["## SECTION 7: TEST FILE INVENTORY", "", "| File | Test Funcs |", "|------|:----------:|"]
    for fp, n in all_test_files:
        lines.append(f"| `{fp}` | {n} |")
    return "\n".join(lines)


def collect_stats():
    """Scan the source tree once and return (mod_stats, plat_stats, pkg_stats_list)."""
    mod_stats = []
    for name in sorted(os.listdir(MODULES_DIR)):
        d = os.path.join(MODULES_DIR, name)
        if os.path.isdir(d):
            mod_stats.append(module_stats(d, name))

    plat_stats = []
    for name in sorted(os.listdir(PLATFORM_DIR)):
        d = os.path.join(PLATFORM_DIR, name)
        if os.path.isdir(d):
            plat_stats.append(module_stats(d, name))

    pkg_stats_list = []
    for name in sorted(os.listdir(PKG_DIR)):
        d = os.path.join(PKG_DIR, name)
        if os.path.isdir(d):
            pkg_stats_list.append(pkg_stats(d, name))

    return mod_stats, plat_stats, pkg_stats_list


def render(mod_stats, plat_stats, pkg_stats_list):
    """Build the full report body from precomputed stats (does NOT write)."""
    mod_table, mod_totals = fmt_table(mod_stats)
    plat_table, plat_totals = fmt_table(plat_stats)

    # ---- SECTION 1 body ----
    s1 = ["## SECTION 1: TENANT MODULES (internal/modules/)", "", mod_table]

    # Test breakdown per module (only modules with any test file)
    tbd = ["", "### Test Breakdown per Module", "",
           "| Module | Repo Tests | Service Tests | Handler Tests | Other | Total |",
           "|--------|:----------:|:-------------:|:-------------:|:-----:|:-----:|"]
    for s in sorted(mod_stats, key=lambda x: x["name"]):
        if s["tests"] == 0 and not s["test_files"]:
            continue
        tbd.append(
            f"| {s['name']} | {s['repo_tests']} | {s['service_tests']} | {s['handler_tests']} | "
            f"{s['other_tests']} | {s['tests']} |"
        )
    s1.extend(tbd)

    # ---- SECTION 2 ----
    s2 = ["", "## SECTION 2: PLATFORM MODULES (internal/platform/)", "", plat_table]

    # ---- SECTION 3 ----
    pkg_header = "| Package | Go Files | Test Funcs | Description |"
    pkg_sep = "|---------|:--------:|:----------:|-------------|"
    pkg_rows = []
    for p in sorted(pkg_stats_list, key=lambda x: x["name"]):
        pkg_rows.append(f"| {p['name']} | {p['files']} | {p['tests']} | {p['desc']} |")
    pkg_total_files = sum(p["files"] for p in pkg_stats_list)
    pkg_total_tests = sum(p["tests"] for p in pkg_stats_list)
    pkg_rows.append(f"| **TOTAL** | **{pkg_total_files}** | **{pkg_total_tests}** | |")
    s3 = ["", "## SECTION 3: SHARED KERNEL PACKAGES (internal/pkg/)", "",
          "\n".join([pkg_header, pkg_sep] + pkg_rows)]

    # ---- SECTION 4 & 5 ----
    s4 = ["", section4(mod_stats + plat_stats)]
    s5 = ["", section5(mod_stats + plat_stats)]

    # ---- SECTION 6: GRAND TOTALS ----
    all_test_files = []
    for s in mod_stats + plat_stats:
        all_test_files.extend(s["test_files"])
    all_test_files.sort(key=lambda x: x[0])

    total_src = 0
    total_test = 0
    for root, dirs, files in os.walk(BACKEND):
        dirs[:] = [d for d in dirs if d not in (".git", "vendor")]
        for fn in files:
            if fn.endswith(".go"):
                total_src += 1
                if fn.endswith("_test.go"):
                    total_test += 1

    entities_all = mod_totals["entities"] + plat_totals["entities"]
    service_all = mod_totals["service"] + plat_totals["service"]
    repo_all = mod_totals["repo"] + plat_totals["repo"]
    handler_all = mod_totals["handler"] + plat_totals["handler"]
    routes_all = mod_totals["routes"] + plat_totals["routes"]
    tests_all = mod_totals["tests"] + plat_totals["tests"] + pkg_total_tests
    layers = len(mod_stats) + len(plat_stats) + len(pkg_stats_list)

    s6 = [
        "", "## SECTION 6: GRAND TOTALS", "",
        "| Category | Count |", "|----------|:-----:|",
        f"| Tenant Modules | {len(mod_stats)} |",
        f"| Platform Modules | {len(plat_stats)} |",
        f"| Shared Kernel Packages | {len(pkg_stats_list)} |",
        f"| **Total Architecture Layers** | **{layers}** |",
        f"| Total GORM Entities (tenant) | {mod_totals['entities']} |",
        f"| Total GORM Entities (platform) | {plat_totals['entities']} |",
        f"| **Total Entities (combined)** | **{entities_all}** |",
        f"| Total Service Methods | {service_all} |",
        f"| Total Repository Methods | {repo_all} |",
        f"| Total Handler Functions | {handler_all} |",
        f"| Total Route Registrations | {routes_all} |",
        f"| **Total Unit Tests (all)** | **{tests_all}** |",
        f"| Total Go Source Files | {total_src - total_test} |",
        f"| Total Test Files (_test.go) | {total_test} |",
        f"| **Total Go Files** | **{total_src}** |",
    ]

    # ---- SECTION 7 ----
    s7 = ["", section7(all_test_files)]

    header = [
        "=" * 100,
        "  HRIS PLATFORM — GO MODULE ARCHITECTURE REPORT",
        f"  Generated: {datetime.date.today().strftime('%d %b %Y')}",
        "",
        "  Index dokumentasi: docs/README.md  |  Terkait: platform-architecture-design.md, openapi-report.md",
        "=" * 100,
        "",
    ]

    body = "\n".join(header + s1 + s2 + s3 + s4 + s5 + s6 + s7)
    body += "\n\n====================================================================================================\n  END OF REPORT\n"
    return body


def build_doc():
    """Return the full report content as a string (does NOT write to disk)."""
    return render(*collect_stats())


def main():
    mod_stats, plat_stats, pkg_stats_list = collect_stats()
    body = render(mod_stats, plat_stats, pkg_stats_list)

    with open(OUT_PATH, "w", encoding="utf-8") as f:
        f.write(body)

    _, mod_totals = fmt_table(mod_stats)
    _, plat_totals = fmt_table(plat_stats)
    pkg_total_tests = sum(p["tests"] for p in pkg_stats_list)

    entities_all = mod_totals["entities"] + plat_totals["entities"]
    service_all = mod_totals["service"] + plat_totals["service"]
    repo_all = mod_totals["repo"] + plat_totals["repo"]
    handler_all = mod_totals["handler"] + plat_totals["handler"]
    routes_all = mod_totals["routes"] + plat_totals["routes"]
    tests_all = mod_totals["tests"] + plat_totals["tests"] + pkg_total_tests

    total_src = 0
    total_test = 0
    for root, dirs, files in os.walk(BACKEND):
        dirs[:] = [d for d in dirs if d not in (".git", "vendor")]
        for fn in files:
            if fn.endswith(".go"):
                total_src += 1
                if fn.endswith("_test.go"):
                    total_test += 1

    print(f"Generated: {OUT_PATH}")
    print(f"  Tenant modules: {len(mod_stats)}, Platform: {len(plat_stats)}, Shared pkg: {len(pkg_stats_list)}")
    print(f"  Entities: {entities_all}, Service: {service_all}, Repo: {repo_all}, Handler: {handler_all}, Routes: {routes_all}")
    print(f"  Unit tests (all): {tests_all}, Go files: {total_src} ({total_src - total_test} source + {total_test} test)")


if __name__ == "__main__":
    main()
