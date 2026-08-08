#!/usr/bin/env python3
"""Check documentation gaps across the docs/ tree.

Validates three things:
  1) Broken internal markdown links (excluding ignored folders).
  2) Inline `docs/xxx.md` references pointing to non-existent files.
  3) Markdown files in docs/ root that may be missing from the docs/README.md index.

Run from the project root:
    python scripts/check_doc_gaps.py

Ignored folders (tmp/, archive/, backlog/, seeder/) are excluded from scanning —
see docs/README.md → "Folder yang Diabaikan".
"""
import os
import re

root = 'docs'
IGNORED = {'tmp', 'archive', 'backlog', 'seeder', 'node_modules', 'vendor'}

print('=' * 70)
print('1. BROKEN INTERNAL LINKS (excl. ignored folders)')
print('=' * 70)
broken = []
link_re = re.compile(r'\[([^\]]*)\]\(([^)]+)\)')
for dirpath, dirnames, filenames in os.walk(root):
    dirnames[:] = [d for d in dirnames if d not in IGNORED]
    for fn in filenames:
        if not fn.endswith('.md'):
            continue
        path = os.path.join(dirpath, fn)
        with open(path, encoding='utf-8') as f:
            content = f.read()
        for m in link_re.finditer(content):
            target = m.group(2)
            if target.startswith(('http://', 'https://', 'mailto:', '#')):
                continue
            if target.startswith('data:'):
                continue
            filepart = target.split('#')[0]
            if not filepart:
                continue
            parts = filepart.replace('\\', '/').split('/')
            if any(p in IGNORED for p in parts[:-1]):
                broken.append(f'{path}: -> {target} (folder diabaikan)')
                continue
            resolved = os.path.normpath(os.path.join(dirpath, filepart))
            if not os.path.exists(resolved):
                broken.append(f'{path}: -> {target}')
if broken:
    for b in broken:
        print(' ', b)
    print(f'TOTAL: {len(broken)}')
else:
    print('OK: no broken links')

print()
print('=' * 70)
print('2. FILE YANG DIRUJUK (inline, bukan link) TAPI TIDAK ADA')
print('=' * 70)
# scan for `docs/xxx.md` inline mentions pointing to non-existent files (outside ignored folders)
inline_re = re.compile(r'`(docs/[A-Za-z0-9_\-./]+\.md)`')
missing = set()
for dirpath, dirnames, filenames in os.walk(root):
    dirnames[:] = [d for d in dirnames if d not in IGNORED]
    for fn in filenames:
        if not fn.endswith('.md'):
            continue
        path = os.path.join(dirpath, fn)
        with open(path, encoding='utf-8') as f:
            content = f.read()
        for m in inline_re.finditer(content):
            ref = m.group(1)
            if ref.startswith('docs/archive') or ref.startswith('docs/tmp') or ref.startswith('docs/backlog') or ref.startswith('docs/seeder'):
                continue
            if not os.path.exists(ref):
                missing.add(f'{path}: -> {ref}')
if missing:
    for x in sorted(missing):
        print(' ', x)
    print(f'TOTAL: {len(missing)}')
else:
    print('OK: no missing inline doc references')

print()
print('=' * 70)
print('3. DAFTAR DOKUMEN: docs/ root (untuk dicek vs docs/README.md navigasi)')
print('=' * 70)
actual = sorted(fn for fn in os.listdir(root) if fn.endswith('.md') and fn != 'README.md')
print('File markdown di docs/ root:', len(actual))
for a in actual:
    print('  ', a)
