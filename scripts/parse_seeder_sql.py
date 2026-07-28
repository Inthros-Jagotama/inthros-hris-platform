#!/usr/bin/env python3
"""
Parse PHP seeder files (DistrictSeeder.php, VillageSeeder.php) and
generate MySQL multi-row INSERT SQL files.

The PHP files split data across multiple insert() calls (15 for districts, 167 for villages)
to avoid PHP memory limits.
"""

import re
import os
import sys
import time

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
OUTPUT_DIR = os.path.join(BASE_DIR, "backend", "migrations", "seeders")


def find_all_insert_calls(content):
    """Find all DB::table(...)->insert( array( positions."""
    inserts = []
    for m in re.finditer(r"DB::table\([^)]+\)->insert\s*\(\s*array\s*\(", content, re.DOTALL):
        inserts.append(m.end())
    return inserts


def find_matching_close(content, start_pos):
    """
    Find the position of the matching closing paren for content[start_pos],
    where start_pos is right after the opening '('.
    """
    depth = 1
    i = start_pos
    while i < len(content) and depth > 0:
        c = content[i]
        if c == '(':
            depth += 1
        elif c == ')':
            depth -= 1
        i += 1
    return i - 1  # position of closing paren


def parse_records_from_body(body):
    """Parse individual array records from an array body."""
    records = []
    pos = 0
    
    while pos < len(body):
        # Find 'array(' or 'array ('
        idx1 = body.find('array(', pos)
        idx2 = body.find('array (', pos)
        
        idx = -1
        if idx1 != -1 and idx2 != -1:
            idx = min(idx1, idx2)
        elif idx1 != -1:
            idx = idx1
        elif idx2 != -1:
            idx = idx2
        else:
            break
        
        # Move past 'array' to find '('
        j = idx + 5
        while j < len(body) and body[j] in ' \n\r\t':
            j += 1
        
        if j >= len(body) or body[j] != '(':
            pos = idx + 5
            continue
        
        # Find matching closing paren
        inner_start = j + 1
        close = find_matching_close(body, inner_start)
        
        if close == -1 or close <= inner_start:
            pos = idx + 5
            continue
        
        inner = body[inner_start:close]
        
        # Extract key-value pairs
        # Handle PHP-escaped single quotes: \' -> '
        kv_pattern = r"'([^']*?)'\s*=>\s*'([^']*?)'"
        # Try with escaped quote handling first
        escaped_kv_pattern = r"'([^']*?)'\s*=>\s*'((?:[^'\\]|\\.)*?)'"
        
        pairs = re.findall(escaped_kv_pattern, inner)
        
        if pairs and len(pairs) >= 2:
            record = {}
            for key, val in pairs:
                val = val.replace("\\'", "'").replace("\\\\", "\\")
                record[key] = val
            records.append(record)
        
        pos = close + 1  # After this array
    
    return records


def parse_php_records(filepath):
    """
    Parse ALL records from a PHP seeder file by iterating over all insert() calls.
    """
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()
    
    inserts = find_all_insert_calls(content)
    print(f"  Found {len(inserts)} insert() calls", file=sys.stderr)
    
    all_records = []
    total_body_chars = 0
    
    for nth, body_start in enumerate(inserts, 1):
        close = find_matching_close(content, body_start)
        body = content[body_start:close]
        records = parse_records_from_body(body)
        all_records.extend(records)
        total_body_chars += len(body)
        
        if nth % 10 == 0 or nth == len(inserts):
            print(f"  ... processed insert {nth}/{len(inserts)} ({len(records)} records in this batch, total: {len(all_records)})", file=sys.stderr)
    
    return all_records


def escape_mysql(val):
    """Escape a string for MySQL INSERT."""
    if val is None:
        return "NULL"
    escaped = val.replace("\\", "\\\\").replace("'", "\\'")
    return f"'{escaped}'"


def build_insert_sql(table, records, columns, batch_size=500):
    """Generate multi-row INSERT statements."""
    lines = []
    total = len(records)
    
    for i in range(0, total, batch_size):
        batch = records[i:i + batch_size]
        values_list = []
        
        for rec in batch:
            vals = ", ".join(escape_mysql(rec.get(col, "")) for col in columns)
            values_list.append(f"({vals})")
        
        cols = ", ".join(f"`{col}`" for col in columns)
        insert = f"INSERT INTO `{table}` ({cols}) VALUES\n  " + ",\n  ".join(values_list) + ";"
        lines.append(insert)
    
    return "\n\n".join(lines)


def generate_districts_sql():
    """Generate districts SQL from DistrictSeeder.php."""
    php_file = os.path.join(BASE_DIR, "docs", "seeder", "DistrictSeeder.php")
    sql_file = os.path.join(OUTPUT_DIR, "002_seed_districts.sql")
    
    print(f"Reading {php_file}...")
    records = parse_php_records(php_file)
    print(f"  TOTAL: {len(records)} district records")
    
    if not records:
        print("  ERROR: No records found!")
        return False
    
    enriched = []
    for rec in records:
        enriched.append({
            "id": rec.get("id", ""),
            "code": rec.get("id", ""),
            "regency_id": rec.get("regency_id", ""),
            "name": rec.get("name", ""),
        })
    
    sql = f"""
-- ============================================================
-- Districts Seeder (Kecamatan) - {len(enriched)} records
-- Generated from docs/seeder/DistrictSeeder.php
-- ============================================================

DELETE FROM `districts` WHERE 1=1;

"""
    sql += build_insert_sql("districts", enriched, ["id", "code", "regency_id", "name"], batch_size=500)
    
    os.makedirs(os.path.dirname(sql_file), exist_ok=True)
    with open(sql_file, "w", encoding="utf-8") as f:
        f.write(sql)
    
    file_size = os.path.getsize(sql_file)
    print(f"  Written to {sql_file} ({file_size:,} bytes)")
    return True


def generate_villages_sql():
    """Generate villages SQL from VillageSeeder.php."""
    php_file = os.path.join(BASE_DIR, "docs", "seeder", "VillageSeeder.php")
    sql_file = os.path.join(OUTPUT_DIR, "003_seed_villages.sql")
    
    print(f"Reading {php_file}...")
    records = parse_php_records(php_file)
    print(f"  TOTAL: {len(records)} village records")
    
    if not records:
        print("  ERROR: No records found!")
        return False
    
    enriched = []
    for rec in records:
        enriched.append({
            "id": rec.get("id", ""),
            "code": rec.get("id", ""),
            "district_id": rec.get("district_id", ""),
            "name": rec.get("name", ""),
        })
    
    sql = f"""
-- ============================================================
-- Villages Seeder (Desa/Kelurahan) - {len(enriched)} records
-- Generated from docs/seeder/VillageSeeder.php
-- ============================================================

DELETE FROM `villages` WHERE 1=1;

"""
    sql += build_insert_sql("villages", enriched, ["id", "code", "district_id", "name"], batch_size=500)
    
    os.makedirs(os.path.dirname(sql_file), exist_ok=True)
    with open(sql_file, "w", encoding="utf-8") as f:
        f.write(sql)
    
    file_size = os.path.getsize(sql_file)
    print(f"  Written to {sql_file} ({file_size:,} bytes)")
    return True


def main():
    start = time.time()
    
    print("=" * 60)
    print("PHP Seeder to SQL Converter")
    print("=" * 60)
    print()
    
    print("[1/2] Generating districts SQL...")
    d_ok = generate_districts_sql()
    
    print()
    print("[2/2] Generating villages SQL...")
    v_ok = generate_villages_sql()
    
    elapsed = time.time() - start
    print()
    print(f"Done in {elapsed:.1f}s")
    
    if d_ok and v_ok:
        print("OK - Both SQL files generated successfully!")
        return 0
    else:
        print("ERROR - Some files had errors")
        return 1


if __name__ == "__main__":
    sys.exit(main())
