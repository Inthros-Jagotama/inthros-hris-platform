---
name: database-migration-consolidation
description: Use when a module has accumulated many historical migration files and someone asks to consolidate/squash/baseline them into one migration, when migration folders feel cluttered, or before writing a new "baseline" schema migration for a module in backend/internal/pkg/migrator/migrations/tenant/{postgres,mysql}/.
---

# Database Migration Consolidation

## Overview

Convert a module's historical migrations into one final baseline migration that reproduces the exact same end schema, with seed/reference data (`INSERT`) kept separate from schema definition (`CREATE`/`ALTER`/index/constraint). This repo's migrator (`backend/internal/pkg/migrator/`) tracks applied versions by numeric filename prefix per tenant, with no real checksum (`Checksum = len(content)`) and no squash/supersede concept — that constrains what "consolidation" can safely mean here.

## The Hard Constraint

**Never delete, renumber, or edit an already-shipped migration file's content.** Tenants are migrated independently and may be sitting at any version. The tracking table only knows "version N applied y/n" — it cannot detect that a file's content changed, and a deleted/renumbered file breaks every tenant that already ran it (or never will).

So a "baseline" is additive: a **new**, higher-numbered migration containing the consolidated schema, or a **non-executed reference file** (kept outside the migrator's scanned folder, e.g. `docs/schema/`) if the goal is documentation/readability rather than changing what fresh tenants run. Ask which one is wanted before writing SQL — don't assume.

## Steps

1. **Inventory** — list every up/down pair in the module's numeric range, for both postgres and mysql, confirm filename parity (`diff` the two directory listings).
2. **Read every file fully, in order** — not just the obvious ones. For each `CREATE`/`ALTER` note its target table. Two things to catch here decide the whole plan:
   - **Intra-range mutation**: does a later file in the range `ALTER`/rename/drop a column the range itself created earlier? The baseline's `CREATE TABLE` must reflect the *final* shape after all such mutations — never just concatenate the original `CREATE` statements.
   - **Cross-module boundary**: does any file touch a table owned by a different module? Flag it explicitly and get a decision on which module's baseline it belongs to — don't silently fold it in or silently drop it.
3. **Separate seed from schema** — grep the range for `INSERT INTO`. Anything inserting lookup/master/reference data goes in its own seed file (follow this repo's existing `NNN_seed_<name>.sql` convention), never mixed into the schema file. If there's no seed data, say so explicitly in a comment rather than leaving it ambiguous.
4. **Write the consolidated schema file** — `CREATE TABLE` in dependency order (parents before children), each in its final column shape, then indexes/constraints. Write the paired `.down.sql` as `DROP` in reverse dependency order.
5. **Verify byte-for-byte schema equivalence** — this is not optional and not "looks right to me":
   - DB A: run the full original sequence unmodified.
   - DB B: run everything before the range, then the new baseline, then everything after.
   - Diff structurally: table list, columns (type/default/nullable), indexes, constraints — via `pg_dump --schema-only` / `SHOW CREATE TABLE` (mysql), and cross-check against `information_schema` directly since dump formatting can mask diffs.
   - Run the down-migration path in both and confirm convergence too.
6. **Report open questions instead of guessing** — module boundary calls, whether old files stay live forever vs. become a reference doc, and any file you didn't fully read before concluding "safe to concatenate."

## Common Mistakes

| Mistake | Why it breaks |
|---|---|
| Deleting/renumbering old migration files | Breaks tenants already at those versions; migrator has no drift detection to catch it |
| Concatenating original `CREATE` statements without checking for later `ALTER`s in-range | Produces a baseline that doesn't match the real final schema |
| Mixing `INSERT` seed data into the schema file | Violates schema/data separation; makes the baseline non-reusable for fresh vs. re-seeded installs |
| Treating "looks equivalent" as verified | Only a structural schema diff across two real migrated databases actually proves it |
| Assuming a cross-module `ALTER` belongs to the module being consolidated | Silently changes another module's ownership boundary |
