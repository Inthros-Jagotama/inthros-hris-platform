#!/usr/bin/env python3
"""Verify docs/database-schema.md is in sync with the SQL migrations.

Does NOT write to the doc — it regenerates the content in memory (via
generate_database_schema_doc.build_doc) and compares against the committed
file. Exit code 0 = in sync, 1 = out of sync.
"""

import io
import os
import sys
import difflib

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, SCRIPT_DIR)

import generate_database_schema_doc as gen  # noqa: E402

OUTPUT = gen.OUTPUT


def dialect_mismatch(mysql_tables):
    """Return None if both dialects expose the same table set, else a message."""
    pg = set(gen.tenant.keys())
    my = set(mysql_tables.keys())
    if pg == my:
        return None
    only_pg = sorted(pg - my)
    only_my = sorted(my - pg)
    msg = [
        f"DIALECT MISMATCH: postgres has {len(pg)} tables but mysql has {len(my)} tables.",
        "Both dialects must expose the same table set (postgres/ & mysql/ are mirrors).",
    ]
    if only_pg:
        msg.append(f"  only in postgres ({len(only_pg)}): {', '.join(only_pg)}")
    if only_my:
        msg.append(f"  only in mysql ({len(only_my)}): {', '.join(only_my)}")
    return "\n".join(msg)


def main():
    expected = gen.build_doc()
    mysql_tables = gen.load_all(gen.MYSQL_TENANT_DIR)

    if not os.path.exists(OUTPUT):
        print(f"ERROR: {OUTPUT} not found.")
        print("Run `python scripts/generate_database_schema_doc.py` (or `make db-docs`) first.")
        return 1

    with io.open(OUTPUT, encoding="utf-8", errors="ignore") as f:
        actual = f.read()

    # 1) Postgres (source of truth for the doc content) must match the file.
    if actual != expected:
        print("MISMATCH: docs/database-schema.md is OUT OF SYNC with SQL migrations (postgres).")
        print("Regenerate with: `python scripts/generate_database_schema_doc.py` or `make db-docs`.")
        print()
        diff = difflib.unified_diff(
            actual.splitlines(), expected.splitlines(),
            fromfile="docs/database-schema.md (current)", tofile="docs/database-schema.md (generated)",
            lineterm="",
        )
        lines = list(diff)
        for line in lines[:60]:
            print(line)
        if len(lines) > 60:
            print(f"... ({len(lines) - 60} more diff lines)")
        return 1

    # 2) MySQL must expose the same table set as postgres (doc table counts must hold for both dialects).
    msg = dialect_mismatch(mysql_tables)
    if msg:
        print(msg)
        print("Fix the migrations (or regenerate the doc) so both dialects match.")
        return 1

    print(f"OK: docs/database-schema.md is in sync with SQL migrations (postgres {len(gen.tenant)} tables = mysql {len(mysql_tables)} tables).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
