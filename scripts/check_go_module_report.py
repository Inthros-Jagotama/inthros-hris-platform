#!/usr/bin/env python3
"""Verify docs/go-module-architecture-report.md is in sync with the Go source.

Does NOT write to the doc — it regenerates the content in memory (via
generate_go_module_report.build_doc) and compares against the committed file.
Exit code 0 = in sync, 1 = out of sync.

The `Generated:` date line is excluded from comparison (the report header
carries the generation date, which is not part of the code-derived content).
"""
import difflib
import io
import os
import re
import sys

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, SCRIPT_DIR)

import generate_go_module_report as gen  # noqa: E402

OUTPUT = gen.OUT_PATH

RE_GENERATED = re.compile(r"^  Generated: .*$", re.M)


def main():
    if not os.path.exists(OUTPUT):
        print(f"ERROR: {OUTPUT} not found.")
        print("Run `python scripts/generate_go_module_report.py` (or `make arch-report`) first.")
        return 1

    with io.open(OUTPUT, encoding="utf-8", errors="ignore") as f:
        actual = f.read()

    expected = gen.build_doc()

    # Date-agnostic: the header date changes on every regenerate even when the
    # code-derived content is identical, so strip it from both sides.
    actual_norm = RE_GENERATED.sub("  Generated: <date>", actual)
    expected_norm = RE_GENERATED.sub("  Generated: <date>", expected)

    if actual_norm != expected_norm:
        print("MISMATCH: docs/go-module-architecture-report.md is OUT OF SYNC with the Go source.")
        print("Regenerate with: `python scripts/generate_go_module_report.py` or `make arch-report`.")
        print()
        diff = difflib.unified_diff(
            actual_norm.splitlines(), expected_norm.splitlines(),
            fromfile="docs/go-module-architecture-report.md (current)",
            tofile="docs/go-module-architecture-report.md (generated)",
            lineterm="",
        )
        lines = list(diff)
        for line in lines[:60]:
            print(line)
        if len(lines) > 60:
            print(f"... ({len(lines) - 60} more diff lines)")
        return 1

    print("OK: docs/go-module-architecture-report.md is in sync with the Go source.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
