#!/usr/bin/env python3
"""Inject missing payroll endpoints into openapi.json (idempotent).

Covers endpoints registered in backend/internal/modules/payroll/routes.go that
are missing from the OpenAPI spec:
- Formula engine (validate, variables)
- Salary grade components + salary employee components
- Employee bank/BPJS/tax profiles (GET lists)
- BPJS rate components (GET)
- PPh21 tax brackets (PUT/DELETE {id})
- Payroll runs: dashboard, employees, items, payslips, payments (+export),
  reports (summary/detail/bpjs/tax/bank), calculate, payslips/payments generation
- Payslips: get, html, publish, cancel
- Payments: get, status update

Usage:
    python scripts/inject_payroll_updates_openapi.py
"""
import json
import os

PROJECT_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
JSON_PATH = os.path.join(PROJECT_ROOT, "backend", "internal", "pkg", "docs", "openapi.json")

SEC = [{"bearerAuth": []}]


def param(name, where="path", required=True, schema=None):
    return {"name": name, "in": where, "required": required,
            "schema": schema or {"type": "string", "format": "uuid"}}


def qparam(name, schema=None):
    return param(name, where="query", required=False, schema=schema or {"type": "string"})


def ref(name):
    return {"$ref": f"#/components/schemas/{name}"}


def content(schema_ref, required=True):
    return {"required": required, "content": {"application/json": {"schema": ref(schema_ref)}}}


def responses_created(schema_ref):
    return {
        "400": {"description": "Bad request / Validation error"},
        "201": {"description": "Resource created successfully", "content": {"application/json": {"schema": ref(schema_ref)}}},
    }


def responses_ok(schema_ref):
    return {
        "200": {"description": "OK", "content": {"application/json": {"schema": ref(schema_ref)}}},
        "400": {"description": "Bad request / Validation error"},
        "404": {"description": "Resource not found"},
    }


def responses_list(schema_ref):
    return {
        "200": {"description": "Paginated list of resources", "content": {"application/json": {"schema": ref(schema_ref)}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_array(item_ref, desc="List of resources"):
    """Plain non-paginated list: {success, data: [item]}."""
    return {
        "200": {
            "description": desc,
            "content": {"application/json": {"schema": {"type": "object", "properties": {
                "success": {"type": "boolean", "example": True},
                "data": {"type": "array", "items": ref(item_ref)}}}}}},
        "400": {"description": "Bad request / Validation error"},
    }


def responses_plain(desc_200="OK", content_type="application/json"):
    c = {"application/json": {"schema": {"type": "object", "properties": {"success": {"type": "boolean", "example": True}}}}}
    if content_type != "application/json":
        c = {content_type: {"schema": {"type": "string"}}}
    return {
        "200": {"description": desc_200, "content": c},
        "400": {"description": "Bad request / Validation error"},
    }


def prop(type_, fmt=None, enum=None, desc=None, example=None):
    p = {"type": type_}
    if fmt:
        p["format"] = fmt
    if enum:
        p["enum"] = enum
    if desc:
        p["description"] = desc
    if example is not None:
        p["example"] = example
    return p


def obj(props, required=None):
    o = {"type": "object", "properties": props}
    if required:
        o["required"] = required
    return o


def arr(items):
    return {"type": "array", "items": items}


def add_schema(spec, name, schema):
    spec["components"]["schemas"][name] = schema


def add_endpoint(spec, method, path, summary, request_body=None, responses=None, query=None):
    op = {
        "tags": ["Payroll"],
        "summary": summary,
        "security": SEC,
        "responses": responses or responses_plain(),
    }
    params = []
    for p in path.split("/"):
        if p.startswith("{") and p.endswith("}"):
            params.append(param(p[1:-1]))
    if query:
        params.extend(query)
    if params:
        op["parameters"] = params
    if request_body:
        op["requestBody"] = content(request_body)
    spec["paths"].setdefault(path, {})[method.lower()] = op


def main():
    with open(JSON_PATH, encoding="utf-8") as f:
        spec = json.load(f)

    P = "/api/v1/tenant/payroll"

    # ------------------------------------------------------------------
    # Schemas
    # ------------------------------------------------------------------

    add_schema(spec, "ValidateFormulaRequest", obj({
        "formula": prop("string", desc="Formula yang divalidasi, mis. BASIC_SALARY * 0.5"),
    }, required=["formula"]))

    add_schema(spec, "ValidateFormulaResponse", obj({
        "valid": prop("boolean"),
        "formula": prop("string"),
        "variables": arr(prop("string")),
    }))

    add_schema(spec, "CreateSalaryGradeComponentRequest", obj({
        "grading_id": prop("string", fmt="uuid"),
        "salary_component_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "is_mandatory": prop("boolean"),
        "is_default": prop("boolean"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "notes": prop("string"),
    }, required=["salary_component_id", "amount", "effective_start_date"]))

    add_schema(spec, "UpdateSalaryGradeComponentRequest", obj({
        "grading_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "is_mandatory": prop("boolean"),
        "is_default": prop("boolean"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "notes": prop("string"),
    }))

    add_schema(spec, "SalaryGradeComponentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "grading_id": prop("string", fmt="uuid"),
        "salary_component_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "is_mandatory": prop("boolean"),
        "is_default": prop("boolean"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "notes": prop("string"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "CreateSalaryEmployeeComponentRequest", obj({
        "employee_id": prop("string", fmt="uuid"),
        "employment_id": prop("string", fmt="uuid"),
        "position_id": prop("string", fmt="uuid"),
        "grading_id": prop("string", fmt="uuid"),
        "salary_component_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "source_type": prop("string", enum=["MANUAL", "GRADE_INHERIT", "FORMULA", "ADJUSTMENT"]),
        "source_ref_id": prop("string", fmt="uuid"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "notes": prop("string"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
    }, required=["employee_id", "salary_component_id", "amount", "effective_start_date"]))

    add_schema(spec, "UpdateSalaryEmployeeComponentRequest", obj({
        "amount": prop("number"),
        "currency_code": prop("string"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "source_type": prop("string", enum=["MANUAL", "GRADE_INHERIT", "FORMULA", "ADJUSTMENT"]),
        "notes": prop("string"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
    }))

    add_schema(spec, "SalaryEmployeeComponentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "employment_id": prop("string", fmt="uuid"),
        "position_id": prop("string", fmt="uuid"),
        "grading_id": prop("string", fmt="uuid"),
        "salary_component_id": prop("string", fmt="uuid"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "source_type": prop("string", enum=["MANUAL", "GRADE_INHERIT", "FORMULA", "ADJUSTMENT"]),
        "source_ref_id": prop("string", fmt="uuid"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "notes": prop("string"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "UpdatePph21TaxBracketRequest", obj({
        "bracket_order": prop("integer"),
        "lower_bound": prop("number"),
        "upper_bound": prop("number"),
        "rate_percent": prop("number"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
    }))

    add_schema(spec, "Pph21TaxBracketResponse", obj({
        "id": prop("string", fmt="uuid"),
        "bracket_order": prop("integer"),
        "lower_bound": prop("number"),
        "upper_bound": prop("number"),
        "rate_percent": prop("number"),
        "effective_start_date": prop("string", fmt="date"),
        "effective_end_date": prop("string", fmt="date"),
        "status": prop("string", enum=["ACTIVE", "INACTIVE"]),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "PayrollRunResponse", obj({
        "id": prop("string", fmt="uuid"),
        "payroll_period_id": prop("string", fmt="uuid"),
        "run_code": prop("string"),
        "run_type": prop("string", enum=["REGULAR", "OFF_CYCLE", "THR", "BONUS"]),
        "proration_method": prop("string", enum=["CALENDAR_DAYS", "WORKING_DAYS", "FIXED_30_DAYS", "ATTENDANCE_DAYS"]),
        "status": prop("string", enum=["DRAFT", "CALCULATED", "REVIEWED", "APPROVED", "LOCKED", "CANCELLED"]),
        "total_employees": prop("integer"),
        "total_earning": prop("number"),
        "total_deduction": prop("number"),
        "total_employer_contribution": prop("number"),
        "total_net": prop("number"),
        "total_company_cost": prop("number"),
        "calculated_at": prop("string", fmt="date-time"),
        "reviewed_at": prop("string", fmt="date-time"),
        "approved_at": prop("string", fmt="date-time"),
        "locked_at": prop("string", fmt="date-time"),
        "approval_instance_id": prop("string", fmt="uuid"),
        "created_at": prop("string", fmt="date-time"),
        "updated_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "PayrollRunEmployeeResponse", obj({
        "id": prop("string", fmt="uuid"),
        "payroll_run_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "employment_id": prop("string", fmt="uuid"),
        "position_id": prop("string", fmt="uuid"),
        "grading_id": prop("string", fmt="uuid"),
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "position_title": prop("string"),
        "grading_name": prop("string"),
        "total_earning": prop("number"),
        "total_deduction": prop("number"),
        "total_employer_contribution": prop("number"),
        "net_amount": prop("number"),
        "total_company_cost": prop("number"),
        "status": prop("string"),
    }))

    add_schema(spec, "PayrollRunItemResponse", obj({
        "id": prop("string", fmt="uuid"),
        "payroll_run_id": prop("string", fmt="uuid"),
        "payroll_run_employee_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "salary_component_id": prop("string", fmt="uuid"),
        "component_code": prop("string"),
        "component_name": prop("string"),
        "component_type": prop("string"),
        "calculation_type": prop("string"),
        "item_category": prop("string"),
        "paid_by": prop("string", enum=["EMPLOYEE", "EMPLOYER"]),
        "affects_gross_pay": prop("boolean"),
        "affects_net_pay": prop("boolean"),
        "affects_company_cost": prop("boolean"),
        "print_on_payslip": prop("boolean"),
        "amount": prop("number"),
        "base_amount": prop("number"),
        "rate": prop("number"),
        "formula": prop("string"),
        "formula_result": prop("number"),
        "currency_code": prop("string"),
        "source_group": prop("string"),
        "notes": prop("string"),
    }))

    add_schema(spec, "PayslipItemResponse", obj({
        "component_code": prop("string"),
        "component_name": prop("string"),
        "item_category": prop("string"),
        "amount": prop("number"),
    }))

    add_schema(spec, "PayrollPayslipResponse", obj({
        "id": prop("string", fmt="uuid"),
        "payroll_run_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "payslip_number": prop("string"),
        "period_code": prop("string"),
        "period_year": prop("integer"),
        "period_month": prop("integer"),
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "position_title": prop("string"),
        "grading_name": prop("string"),
        "total_earning": prop("number"),
        "total_deduction": prop("number"),
        "total_employer_contribution": prop("number"),
        "net_amount": prop("number"),
        "status": prop("string", enum=["DRAFT", "PUBLISHED", "CANCELLED"]),
        "generated_at": prop("string", fmt="date-time"),
        "published_at": prop("string", fmt="date-time"),
        "cancelled_at": prop("string", fmt="date-time"),
        "items": arr(ref("PayslipItemResponse")),
    }))

    add_schema(spec, "PaymentBatchResponse", obj({
        "run_id": prop("string", fmt="uuid"),
        "total": prop("integer"),
        "total_amount": prop("number"),
        "skipped_no_bank_profile": prop("integer"),
        "status": prop("string"),
    }))

    add_schema(spec, "PayrollPaymentResponse", obj({
        "id": prop("string", fmt="uuid"),
        "payroll_run_id": prop("string", fmt="uuid"),
        "employee_id": prop("string", fmt="uuid"),
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "amount": prop("number"),
        "currency_code": prop("string"),
        "payment_date": prop("string", fmt="date"),
        "bank_name": prop("string"),
        "bank_account_number": prop("string"),
        "bank_account_holder_name": prop("string"),
        "status": prop("string", enum=["PENDING", "PROCESSING", "PAID", "FAILED", "REVERSED"]),
        "reference": prop("string"),
        "failed_reason": prop("string"),
        "processed_at": prop("string", fmt="date-time"),
        "paid_at": prop("string", fmt="date-time"),
        "failed_at": prop("string", fmt="date-time"),
        "reversed_at": prop("string", fmt="date-time"),
        "created_at": prop("string", fmt="date-time"),
    }))

    add_schema(spec, "UpdatePaymentStatusRequest", obj({
        "status": prop("string", enum=["PENDING", "PROCESSING", "PAID", "FAILED", "REVERSED"]),
        "reason": prop("string", desc="Dipakai saat FAILED"),
        "reference": prop("string", desc="Dipakai saat PAID/PROCESSING"),
    }, required=["status"]))

    add_schema(spec, "PayrollSummaryReport", obj({
        "run_id": prop("string", fmt="uuid"),
        "run_code": prop("string"),
        "period_code": prop("string"),
        "status": prop("string"),
        "total_employees": prop("integer"),
        "gross_salary": prop("number"),
        "employee_deduction": prop("number"),
        "employer_contribution": prop("number"),
        "net_salary": prop("number"),
        "total_company_cost": prop("number"),
    }))

    add_schema(spec, "PayrollDetailReportRow", obj({
        "employee_id": prop("string", fmt="uuid"),
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "component_code": prop("string"),
        "component_name": prop("string"),
        "item_category": prop("string"),
        "component_type": prop("string"),
        "paid_by": prop("string"),
        "source_group": prop("string"),
        "base_amount": prop("number"),
        "rate": prop("number"),
        "amount": prop("number"),
    }))

    add_schema(spec, "BpjsReportRow", obj({
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "bpjs_number": prop("string"),
        "wage_basis": prop("number"),
        "employee_contribution": prop("number"),
        "employer_contribution": prop("number"),
        "total_contribution": prop("number"),
    }))

    add_schema(spec, "TaxReportRow", obj({
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "taxable_income": prop("number"),
        "pph21": prop("number"),
    }))

    add_schema(spec, "BankTransferReportRow", obj({
        "employee_code": prop("string"),
        "employee_name": prop("string"),
        "bank_name": prop("string"),
        "account_number": prop("string"),
        "account_holder_name": prop("string"),
        "amount": prop("number"),
        "status": prop("string"),
    }))

    # ------------------------------------------------------------------
    # Endpoints
    # ------------------------------------------------------------------

    # Formula engine
    add_endpoint(spec, "GET", f"{P}/formula/variables", "Daftar variabel built-in formula engine",
                 responses=responses_array({"type": "string"}, "List of formula variables"))
    add_endpoint(spec, "POST", f"{P}/formula/validate", "Validasi sintaks formula",
                 request_body="ValidateFormulaRequest", responses=responses_ok("ValidateFormulaResponse"))

    # Salary grade components
    add_endpoint(spec, "GET", f"{P}/salary-grade-components", "Daftar komponen gaji per grade (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{P}/salary-grade-components/{{id}}", "Detail komponen gaji per grade",
                 responses=responses_ok("SalaryGradeComponentResponse"))
    add_endpoint(spec, "POST", f"{P}/salary-grade-components", "Tambah komponen gaji per grade",
                 request_body="CreateSalaryGradeComponentRequest", responses=responses_created("SalaryGradeComponentResponse"))
    add_endpoint(spec, "PUT", f"{P}/salary-grade-components/{{id}}", "Perbarui komponen gaji per grade",
                 request_body="UpdateSalaryGradeComponentRequest", responses=responses_ok("SalaryGradeComponentResponse"))
    add_endpoint(spec, "DELETE", f"{P}/salary-grade-components/{{id}}", "Hapus komponen gaji per grade",
                 responses=responses_plain("Component deleted"))

    # Salary employee components
    add_endpoint(spec, "GET", f"{P}/salary-employee-components", "Daftar komponen gaji per employee (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{P}/salary-employee-components/{{id}}", "Detail komponen gaji per employee",
                 responses=responses_ok("SalaryEmployeeComponentResponse"))
    add_endpoint(spec, "POST", f"{P}/salary-employee-components", "Tambah komponen gaji per employee (override)",
                 request_body="CreateSalaryEmployeeComponentRequest", responses=responses_created("SalaryEmployeeComponentResponse"))
    add_endpoint(spec, "PUT", f"{P}/salary-employee-components/{{id}}", "Perbarui komponen gaji per employee",
                 request_body="UpdateSalaryEmployeeComponentRequest", responses=responses_ok("SalaryEmployeeComponentResponse"))
    add_endpoint(spec, "DELETE", f"{P}/salary-employee-components/{{id}}", "Hapus komponen gaji per employee",
                 responses=responses_plain("Component deleted"))

    # Employee bank/BPJS/tax profiles + bpjs rate components (GET list)
    add_endpoint(spec, "GET", f"{P}/employee-bank-profiles", "Daftar profil bank karyawan (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{P}/employee-bpjs-profiles", "Daftar profil BPJS karyawan (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{P}/employee-tax-profiles", "Daftar profil pajak karyawan (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])
    add_endpoint(spec, "GET", f"{P}/bpjs-rate-components", "Daftar komponen rate BPJS (paginated)",
                 responses=responses_list("PayrollPaginatedResponse"),
                 query=[qparam("page", {"type": "integer"}), qparam("per_page", {"type": "integer"})])

    # PPh21 tax brackets
    add_endpoint(spec, "PUT", f"{P}/pph21-tax-brackets/{{id}}", "Perbarui bracket PPh21",
                 request_body="UpdatePph21TaxBracketRequest", responses=responses_ok("Pph21TaxBracketResponse"))
    add_endpoint(spec, "DELETE", f"{P}/pph21-tax-brackets/{{id}}", "Hapus bracket PPh21",
                 responses=responses_plain("Bracket deleted"))

    # Payroll runs
    add_endpoint(spec, "POST", f"{P}/runs/{{id}}/calculate", "Hitung payroll run",
                 responses=responses_ok("PayrollRunResponse"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/dashboard", "Dashboard agregat payroll run",
                 responses=responses_ok("PayrollSummaryReport"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/employees", "Daftar employee dalam payroll run",
                 responses=responses_array("PayrollRunEmployeeResponse", "List of payroll run employees"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/items", "Daftar item kalkulasi payroll run",
                 responses=responses_array("PayrollRunItemResponse", "List of payroll run items"))
    add_endpoint(spec, "POST", f"{P}/runs/{{id}}/payslips", "Generate payslips untuk run",
                 responses=responses_array("PayrollPayslipResponse", "Generated payslips"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/payslips", "Daftar payslip dalam run",
                 responses=responses_array("PayrollPayslipResponse", "List of run payslips"))
    add_endpoint(spec, "POST", f"{P}/runs/{{id}}/payments", "Buat batch pembayaran run",
                 responses=responses_ok("PaymentBatchResponse"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/payments", "Daftar pembayaran run",
                 responses=responses_array("PayrollPaymentResponse", "List of run payments"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/payments/export", "Export pembayaran run sebagai CSV",
                 responses=responses_plain("CSV file export", content_type="text/csv"))

    # Run reports
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/reports/summary", "Laporan ringkasan payroll run",
                 responses=responses_ok("PayrollSummaryReport"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/reports/detail", "Laporan detail payroll per employee per komponen",
                 responses=responses_array("PayrollDetailReportRow", "Payroll detail report rows"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/reports/bpjs", "Laporan BPJS run",
                 responses=responses_array("BpjsReportRow", "BPJS report rows"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/reports/tax", "Laporan pajak (PPh21) run",
                 responses=responses_array("TaxReportRow", "Tax report rows"))
    add_endpoint(spec, "GET", f"{P}/runs/{{id}}/reports/bank", "Laporan transfer bank run",
                 responses=responses_array("BankTransferReportRow", "Bank transfer report rows"))

    # Payslips
    add_endpoint(spec, "GET", f"{P}/payslips/{{id}}", "Detail payslip",
                 responses=responses_ok("PayrollPayslipResponse"))
    add_endpoint(spec, "GET", f"{P}/payslips/{{id}}/html", "HTML payslip (untuk cetak)",
                 responses=responses_plain("HTML payslip", content_type="text/html"))
    add_endpoint(spec, "POST", f"{P}/payslips/{{id}}/publish", "Publish payslip",
                 responses=responses_ok("PayrollPayslipResponse"))
    add_endpoint(spec, "POST", f"{P}/payslips/{{id}}/cancel", "Batalkan payslip",
                 responses=responses_ok("PayrollPayslipResponse"))

    # Payments
    add_endpoint(spec, "GET", f"{P}/payments/{{id}}", "Detail pembayaran",
                 responses=responses_ok("PayrollPaymentResponse"))
    add_endpoint(spec, "POST", f"{P}/payments/{{id}}/status", "Perbarui status pembayaran",
                 request_body="UpdatePaymentStatusRequest", responses=responses_ok("PayrollPaymentResponse"))

    with open(JSON_PATH, "w", encoding="utf-8") as f:
        json.dump(spec, f, indent=2, ensure_ascii=False)
    print("OK: payroll endpoints injected into openapi.json")


if __name__ == "__main__":
    main()
