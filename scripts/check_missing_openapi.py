#!/usr/bin/env python3
"""Check which registered API endpoints are missing from openapi.json.

Scans all backend routes.go files (module + platform) plus the extra routes
registered directly in cmd/server/main.go, reconstructs the full URL for each
endpoint (converting gin ":param" to "{param}"), and reports endpoints that
are not documented in backend/internal/pkg/docs/openapi.json.

Usage:
    python scripts/check_missing_openapi.py
"""
import json
import os
import re
import sys

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
BACKEND = os.path.join(PROJECT_ROOT, "backend")
OPENAPI_PATH = os.path.join(BACKEND, "internal", "pkg", "docs", "openapi.json")

GROUP_RE = re.compile(r'(\w+)\s*:=\s*(\w+)\.Group\("([^"]*)"\)')
ROUTE_RE = re.compile(r'(\w+)\.(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)\("([^"]*)"\s*,')


def gin_to_openapi(path: str) -> str:
    """Convert gin path (/:id) to openapi path (/{id})."""
    return re.sub(r":(\w+)", r"{\1}", path)


def normalize_path(path: str) -> str:
    """Normalize path param names for comparison ({id} -> {})."""
    return re.sub(r"\{[^}]*\}", "{}", path)


def collect_routes() -> dict:
    """Return {method: set(paths)} for all registered endpoints."""
    routes = {}

    def add(method: str, path: str):
        path = gin_to_openapi(path)
        if not path.startswith("/"):
            path = "/" + path
        routes.setdefault(method.upper(), set()).add(path)

    # Walk routes.go files
    for root, _dirs, files in os.walk(BACKEND):
        for fname in files:
            if not fname.endswith("routes.go"):
                continue
            fpath = os.path.join(root, fname)
            rel = os.path.relpath(fpath, BACKEND).replace(os.sep, "/")
            # Determine base prefix for root group var
            if rel.startswith("internal/modules"):
                base = "/api/v1/tenant"
            elif rel.startswith("internal/platform") or "pkg/authz" in rel:
                base = "/api/v1/platform"
            else:
                base = ""
            # Root vars in these files (rg / r / router / publicRG)
            prefixes = {"rg": base, "r": base, "router": base, "publicRG": "/api/v1/public"}
            with open(fpath, "r", encoding="utf-8") as f:
                content = f.read()

            # Discover group declarations in order
            for m in GROUP_RE.finditer(content):
                var, parent, path = m.group(1), m.group(2), m.group(3)
                if parent in prefixes:
                    child = (prefixes[parent] + path).rstrip("/")
                    prefixes[var] = child

            for m in ROUTE_RE.finditer(content):
                var, method, path = m.group(1), m.group(2), m.group(3)
                if var in prefixes:
                    full = (prefixes[var] + path).rstrip("/") or prefixes[var]
                    add(method, full)

    # Routes registered directly in cmd/server/main.go
    main_path = os.path.join(BACKEND, "cmd", "server", "main.go")
    with open(main_path, "r", encoding="utf-8") as f:
        main_content = f.read()

    main_routes = [
        ("GET", "/api/v1/tenant/companies/me"),
        ("PUT", "/api/v1/tenant/companies/me"),
        ("GET", "/api/v1/tenant/company-modules"),
        ("GET", "/api/v1/tenant/packages"),
        ("POST", "/api/v1/tenant/packages/{id}/subscribe"),
        ("POST", "/api/v1/tenant/packages/{id}/unsubscribe"),
        ("GET", "/api/v1/public/packages"),
        ("GET", "/api/v1/public/companies/resolve"),
        ("POST", "/api/v1/public/account/setup-password"),
        ("POST", "/api/v1/tenant/auth/login"),
        ("POST", "/api/v1/tenant/auth/refresh"),
        ("GET", "/healthz"),
        ("GET", "/readyz"),
        ("GET", "/docs"),
        ("GET", "/openapi.json"),
    ]
    for method, path in main_routes:
        add(method, path)

    return routes


def main():
    with open(OPENAPI_PATH, "r", encoding="utf-8") as f:
        spec = json.load(f)

    doc_paths = spec.get("paths", {})
    # Build lookup of documented (method, normalized-path)
    documented = set()
    for path, ops in doc_paths.items():
        for method in ops:
            documented.add((method.upper(), normalize_path(path)))

    routes = collect_routes()

    missing = []
    total = 0
    registered_keys = set()
    for method in sorted(routes):
        for path in sorted(routes[method]):
            total += 1
            registered_keys.add((method, normalize_path(path)))
            if (method, normalize_path(path)) not in documented:
                missing.append((method, path))

    # Reverse direction: documented endpoints that are NOT registered in code
    # (phantom docs — e.g. endpoints removed from routes.go but left in the spec)
    phantom = []
    for path, ops in doc_paths.items():
        for method in ops:
            if method == "parameters":
                continue
            if (method.upper(), normalize_path(path)) not in registered_keys:
                phantom.append((method.upper(), path))

    print(f"Total registered endpoints: {total}")
    print(f"Missing from openapi.json:  {len(missing)}")
    print(f"Phantom (documented but not registered): {len(phantom)}")
    print()
    if missing:
        for method, path in missing:
            print(f"  {method:7s} {path}")
        print()
        # Summary per module
        print("=== Summary per prefix ===")
        from collections import Counter
        c = Counter()
        for _m, p in missing:
            seg = p.split("/")
            if len(seg) >= 5:
                key = "/".join(seg[:5])
            else:
                key = p
            c[key] += 1
        for key, n in c.most_common():
            print(f"  {n:3d}  {key}")
    else:
        print("All registered endpoints are documented.")

    if phantom:
        print()
        print("=== Phantom endpoints (documented but NOT registered in routes.go) ===")
        for method, path in sorted(phantom):
            print(f"  {method:7s} {path}")

    return 1 if (missing or phantom) else 0


if __name__ == "__main__":
    sys.exit(main())
