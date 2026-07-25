= HRIS Platform — OpenAPI Comprehensive Report (v8) =

**Generated:** 25 Juli 2026
**Spec Version:** 1.6.9
**Total Paths:** 300
**Total Endpoints (methods):** 544
**Total Schemas:** 352
**Total Tags:** 24

## Coverage Summary

| Metric | Coverage | % |
|---|---|---|
| Endpoints with `summary` | 544/544 | 100% |
| Endpoints with `description` | 544/544 | 100% |
| Endpoints with `operationId` | 544/544 | 100% |

## 1. Endpoints per Module (Tag)

| # | Tag | Endpoints | Paths |
|---|---|---|---|
| 1 | Tenant: Job Management | 88 | 36 |
| 2 | Tenant: Payroll & Compensation Engine | 47 | 24 |
| 3 | Tenant: Competency Management | 36 | 15 |
| 4 | Tenant: Performance Management | 34 | 17 |
| 5 | Tenant: Recruitment & Onboarding (ATS) | 33 | 16 |
| 6 | Tenant: Time & Attendance | 30 | 15 |
| 7 | Tenant: Employees | 29 | 18 |
| 8 | Tenant: Leave & Time Off | 23 | 12 |
| 9 | Tenant: Approval | 15 | 9 |
| 10 | Tenant: Employee Movement & Career Management | 15 | 9 |
| 11 | Tenant: Reimbursement & Claim | 15 | 7 |
| 12 | Platform: Companies | 10 | 7 |
| 13 | Platform: RBAC Management | 10 | 6 |
| 14 | Platform: Modules | 7 | 5 |
| 15 | Tenant: Organizations | 12 | 8 |
| 16 | Tenant: Training & Development Management | 35 | 15 |
| 17 | Health | 4 | 4 |
| 18 | Platform: Licenses | 4 | 2 |
| 19 | Platform: Users | 4 | 2 |
| 20 | Platform: Monitoring | 3 | 3 |
| 21 | Platform: Auth | 2 | 2 |
| 22 | Tenant: Approval Engine | 1 | 1 |
| 23 | Tenant: Workforce Intelligence & Strategic Planning | 68 | 58 |
| 24 | Tenant: Career Intelligence | 19 | 10 |
| | **TOTAL** | **544** | **300** |

## 2. Module Detail

### Tenant: Job Management
**Endpoints:** 88 | **Paths:** 36
**Methods:** DELETE=17 GET=36 POST=17 PUT=18

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/job-management/assets` | List job assets with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/assets` | Create job asset | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/assets/{id}` | Delete job asset | Delete a assets record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/job-management/assets/{id}` | Get job asset by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/assets/{id}` | Update job asset | Update an existing assets record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/job-management/competency-groups` | List competency groups by organization | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/competency-groups` | Create competency group | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/competency-groups/{id}` | Delete competency group | Delete a competency groups record by its unique ID. This action may be reversible depending on sy... |
| `GET` | `/api/v1/tenant/job-management/competency-groups/{id}` | Get competency group by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/competency-groups/{id}` | Update competency group | Update an existing competency groups record by its unique ID. Accepts partial updates; only provi... |
| `GET` | `/api/v1/tenant/job-management/education-experiences` | List job education experiences with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/education-experiences` | Create job education experience | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/education-experiences/{id}` | Delete job education experience | Delete a education experiences record by its unique ID. This action may be reversible depending o... |
| `GET` | `/api/v1/tenant/job-management/education-experiences/{id}` | Get job education experience by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/education-experiences/{id}` | Update job education experience | Update an existing education experiences record by its unique ID. Accepts partial updates; only p... |
| `GET` | `/api/v1/tenant/job-management/financials` | List job financials with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/financials` | Create job financial | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/financials/{id}` | Delete job financial | Delete a financials record by its unique ID. This action may be reversible depending on system co... |
| `GET` | `/api/v1/tenant/job-management/financials/{id}` | Get job financial by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/financials/{id}` | Update job financial | Update an existing financials record by its unique ID. Accepts partial updates; only provided fie... |
| `GET` | `/api/v1/tenant/job-management/hr-authorities` | List HR authorities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/hr-authorities` | Create HR authority | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Delete HR authority | Delete a hr authorities record by its unique ID. This action may be reversible depending on syste... |
| `GET` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Get HR authority by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/hr-authorities/{id}` | Update HR authority | Update an existing hr authorities record by its unique ID. Accepts partial updates; only provided... |
| `GET` | `/api/v1/tenant/job-management/identifications` | List job identifications with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/identifications` | Create a new job identification | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/identifications/{id}` | Delete job identification | Delete a identifications record by its unique ID. This action may be reversible depending on syst... |
| `GET` | `/api/v1/tenant/job-management/identifications/{id}` | Get job identification by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/identifications/{id}` | Update job identification | Update an existing identifications record by its unique ID. Accepts partial updates; only provide... |
| `GET` | `/api/v1/tenant/job-management/objectives` | List job objectives with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/objectives` | Create a new job objective | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/objectives/{id}` | Delete job objective | Delete a objectives record by its unique ID. This action may be reversible depending on system co... |
| `GET` | `/api/v1/tenant/job-management/objectives/{id}` | Get job objective by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/objectives/{id}` | Update job objective | Update an existing objectives record by its unique ID. Accepts partial updates; only provided fie... |
| `GET` | `/api/v1/tenant/job-management/operational-authorities` | List operational authorities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/operational-authorities` | Create operational authority | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Delete operational authority | Delete a operational authorities record by its unique ID. This action may be reversible depending... |
| `GET` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Get operational authority by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/operational-authorities/{id}` | Update operational authority | Update an existing operational authorities record by its unique ID. Accepts partial updates; only... |
| `GET` | `/api/v1/tenant/job-management/potency-competencies` | List potency competencies with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/potency-competencies` | Create potency competency | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Delete potency competency | Delete a potency competencies record by its unique ID. This action may be reversible depending on... |
| `GET` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Get potency competency by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/potency-competencies/{id}` | Update potency competency | Update an existing potency competencies record by its unique ID. Accepts partial updates; only pr... |
| `GET` | `/api/v1/tenant/job-management/relationships` | List relationships with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/relationships` | Create job relationship | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/relationships/{id}` | Delete relationship | Delete a relationships record by its unique ID. This action may be reversible depending on system... |
| `GET` | `/api/v1/tenant/job-management/relationships/{id}` | Get relationship by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/relationships/{id}` | Update relationship | Update an existing relationships record by its unique ID. Accepts partial updates; only provided ... |
| `GET` | `/api/v1/tenant/job-management/responsibilities` | List job responsibilities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/responsibilities` | Create a new job responsibility | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/responsibilities/{id}` | Delete job responsibility | Delete a responsibilities record by its unique ID. This action may be reversible depending on sys... |
| `GET` | `/api/v1/tenant/job-management/responsibilities/{id}` | Get job responsibility by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/responsibilities/{id}` | Update job responsibility | Update an existing responsibilities record by its unique ID. Accepts partial updates; only provid... |
| `GET` | `/api/v1/tenant/job-management/scores` | List job scores with pagination | Retrieve a paginated list of job management resources. |
| `GET` | `/api/v1/tenant/job-management/scores/org/{orgId}` | Get job score by organization | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/scores/org/{orgId}` | Upsert job score for organization | Update an existing {orgId} record by its unique ID. Accepts partial updates; only provided fields... |
| `GET` | `/api/v1/tenant/job-management/subordinate-controls` | List subordinate controls with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/subordinate-controls` | Create subordinate control | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Delete subordinate control | Delete a subordinate controls record by its unique ID. This action may be reversible depending on... |
| `GET` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Get subordinate control by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/subordinate-controls/{id}` | Update subordinate control | Update an existing subordinate controls record by its unique ID. Accepts partial updates; only pr... |
| `GET` | `/api/v1/tenant/job-management/titles` | List job titles with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/titles` | Create a new job title | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/titles/{id}` | Delete job title | Delete a titles record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/job-management/titles/{id}` | Get job title by ID with subs | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/titles/{id}` | Update job title | Update an existing titles record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/job-management/titles/{titleId}/subs` | List subs under a job title | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/titles/{titleId}/subs` | Create a sub under a job title | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Delete job title sub | Delete a {subId} record by its unique ID. This action may be reversible depending on system confi... |
| `GET` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Get job title sub by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/titles/{titleId}/subs/{subId}` | Update job title sub | Update an existing {subId} record by its unique ID. Accepts partial updates; only provided fields... |
| `GET` | `/api/v1/tenant/job-management/values` | List job values with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/values` | Create a new job value | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/values/{id}` | Delete job value | Delete a values record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/job-management/values/{id}` | Get job value by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/values/{id}` | Update job value | Update an existing values record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/job-management/working-activities` | List working activities with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/working-activities` | Create working activity | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/working-activities/{id}` | Delete working activity | Delete a working activities record by its unique ID. This action may be reversible depending on s... |
| `GET` | `/api/v1/tenant/job-management/working-activities/{id}` | Get working activity by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/working-activities/{id}` | Update working activity | Update an existing working activities record by its unique ID. Accepts partial updates; only prov... |
| `GET` | `/api/v1/tenant/job-management/working-risks` | List working risks with pagination | Retrieve a paginated list of job management resources. |
| `POST` | `/api/v1/tenant/job-management/working-risks` | Create working risk | Create a new job management resource. |
| `DELETE` | `/api/v1/tenant/job-management/working-risks/{id}` | Delete working risk | Delete a working risks record by its unique ID. This action may be reversible depending on system... |
| `GET` | `/api/v1/tenant/job-management/working-risks/{id}` | Get working risk by ID | Retrieve a paginated list of job management resources. |
| `PUT` | `/api/v1/tenant/job-management/working-risks/{id}` | Update working risk | Update an existing working risks record by its unique ID. Accepts partial updates; only provided ... |

### Tenant: Payroll & Compensation Engine
**Endpoints:** 47 | **Paths:** 24
**Methods:** DELETE=8 GET=18 POST=12 PUT=9

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/payroll/bpjs-rate-components` | Create BPJS rate component | Create a new bpjs rate components record. Validates required fields and returns the created resou... |
| `DELETE` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Delete BPJS rate component | Delete a bpjs rate components record by its unique ID. This action may be reversible depending on... |
| `GET` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Get BPJS rate component by ID | Retrieve a paginated list of bpjs rate components records. Supports filtering, sorting, and pagin... |
| `PUT` | `/api/v1/tenant/payroll/bpjs-rate-components/{id}` | Update BPJS rate component | Update an existing bpjs rate components record by its unique ID. Accepts partial updates; only pr... |
| `GET` | `/api/v1/tenant/payroll/bpjs-settings` | List BPJS settings | Retrieve a paginated list of bpjs settings records. Supports filtering, sorting, and pagination p... |
| `POST` | `/api/v1/tenant/payroll/bpjs-settings` | Create BPJS setting | Create a new bpjs settings record. Validates required fields and returns the created resource wit... |
| `DELETE` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Delete BPJS setting | Delete a bpjs settings record by its unique ID. This action may be reversible depending on system... |
| `GET` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Get BPJS setting by ID | Retrieve a paginated list of bpjs settings records. Supports filtering, sorting, and pagination p... |
| `PUT` | `/api/v1/tenant/payroll/bpjs-settings/{id}` | Update BPJS setting | Update an existing bpjs settings record by its unique ID. Accepts partial updates; only provided ... |
| `POST` | `/api/v1/tenant/payroll/employee-bank-profiles` | Create employee bank profile | Create a new employee bank profiles record. Validates required fields and returns the created res... |
| `DELETE` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Delete employee bank profile | Delete a employee bank profiles record by its unique ID. This action may be reversible depending ... |
| `GET` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Get employee bank profile by ID | Retrieve a paginated list of employee bank profiles records. Supports filtering, sorting, and pag... |
| `PUT` | `/api/v1/tenant/payroll/employee-bank-profiles/{id}` | Update employee bank profile | Update an existing employee bank profiles record by its unique ID. Accepts partial updates; only ... |
| `POST` | `/api/v1/tenant/payroll/employee-bpjs-profiles` | Create employee BPJS profile | Create a new employee bpjs profiles record. Validates required fields and returns the created res... |
| `DELETE` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Delete employee BPJS profile | Delete a employee bpjs profiles record by its unique ID. This action may be reversible depending ... |
| `GET` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Get employee BPJS profile by ID | Retrieve a paginated list of employee bpjs profiles records. Supports filtering, sorting, and pag... |
| `PUT` | `/api/v1/tenant/payroll/employee-bpjs-profiles/{id}` | Update employee BPJS profile | Update an existing employee bpjs profiles record by its unique ID. Accepts partial updates; only ... |
| `GET` | `/api/v1/tenant/payroll/employee-payroll-profiles` | List employee payroll profiles | Retrieve a paginated list of employee payroll profiles records. Supports filtering, sorting, and ... |
| `POST` | `/api/v1/tenant/payroll/employee-payroll-profiles` | Create employee payroll profile | Create a new employee payroll profiles record. Validates required fields and returns the created ... |
| `DELETE` | `/api/v1/tenant/payroll/employee-payroll-profiles/{id}` | Delete employee payroll profile | Delete a employee payroll profiles record by its unique ID. This action may be reversible dependi... |
| `GET` | `/api/v1/tenant/payroll/employee-payroll-profiles/{id}` | Get employee payroll profile by ID | Retrieve a paginated list of employee payroll profiles records. Supports filtering, sorting, and ... |
| `POST` | `/api/v1/tenant/payroll/employee-tax-profiles` | Create employee tax profile | Create a new employee tax profiles record. Validates required fields and returns the created reso... |
| `DELETE` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Delete employee tax profile | Delete a employee tax profiles record by its unique ID. This action may be reversible depending o... |
| `GET` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Get employee tax profile by ID | Retrieve a paginated list of employee tax profiles records. Supports filtering, sorting, and pagi... |
| `PUT` | `/api/v1/tenant/payroll/employee-tax-profiles/{id}` | Update employee tax profile | Update an existing employee tax profiles record by its unique ID. Accepts partial updates; only p... |
| `GET` | `/api/v1/tenant/payroll/periods` | List payroll periods | Retrieve a paginated list of periods records. Supports filtering, sorting, and pagination paramet... |
| `POST` | `/api/v1/tenant/payroll/periods` | Create payroll period | Create a new periods record. Validates required fields and returns the created resource with its ... |
| `PUT` | `/api/v1/tenant/payroll/periods/{id}` | Update payroll period | Update an existing periods record by its unique ID. Accepts partial updates; only provided fields... |
| `GET` | `/api/v1/tenant/payroll/pph21-ptkp-rates` | List PPh21 PTKP rates | Retrieve a paginated list of pph21 ptkp rates records. Supports filtering, sorting, and paginatio... |
| `POST` | `/api/v1/tenant/payroll/pph21-ptkp-rates` | Create PPh21 PTKP rate | Create a new pph21 ptkp rates record. Validates required fields and returns the created resource ... |
| `GET` | `/api/v1/tenant/payroll/pph21-settings` | List PPh21 settings | Retrieve a paginated list of pph21 settings records. Supports filtering, sorting, and pagination ... |
| `POST` | `/api/v1/tenant/payroll/pph21-settings` | Create PPh21 setting | Create a new pph21 settings record. Validates required fields and returns the created resource wi... |
| `DELETE` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Delete PPh21 setting | Delete a pph21 settings record by its unique ID. This action may be reversible depending on syste... |
| `GET` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Get PPh21 setting by ID | Retrieve a paginated list of pph21 settings records. Supports filtering, sorting, and pagination ... |
| `PUT` | `/api/v1/tenant/payroll/pph21-settings/{id}` | Update PPh21 setting | Update an existing pph21 settings record by its unique ID. Accepts partial updates; only provided... |
| `GET` | `/api/v1/tenant/payroll/pph21-tax-brackets` | List PPh21 tax brackets | Retrieve a paginated list of pph21 tax brackets records. Supports filtering, sorting, and paginat... |
| `POST` | `/api/v1/tenant/payroll/pph21-tax-brackets` | Create PPh21 tax bracket | Create a new pph21 tax brackets record. Validates required fields and returns the created resourc... |
| `GET` | `/api/v1/tenant/payroll/runs` | List payroll runs | Retrieve a paginated list of runs records. Supports filtering, sorting, and pagination parameters. |
| `POST` | `/api/v1/tenant/payroll/runs` | Create payroll run | Create a new runs record. Validates required fields and returns the created resource with its ass... |
| `GET` | `/api/v1/tenant/payroll/runs/{id}` | Get payroll run by ID | Retrieve a paginated list of runs records. Supports filtering, sorting, and pagination parameters. |
| `GET` | `/api/v1/tenant/payroll/runs/{id}/approval` | Check payroll run approval status | Retrieve a paginated list of approval records. Supports filtering, sorting, and pagination parame... |
| `PUT` | `/api/v1/tenant/payroll/runs/{id}/status` | Update payroll run status | Update an existing status record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/payroll/salary-components` | List salary components with pagination | Retrieve a paginated list of salary components records. Supports filtering, sorting, and paginati... |
| `POST` | `/api/v1/tenant/payroll/salary-components` | Create salary component | Create a new salary components record. Validates required fields and returns the created resource... |
| `DELETE` | `/api/v1/tenant/payroll/salary-components/{id}` | Delete salary component | Delete a salary components record by its unique ID. This action may be reversible depending on sy... |
| `GET` | `/api/v1/tenant/payroll/salary-components/{id}` | Get salary component by ID | Retrieve a paginated list of salary components records. Supports filtering, sorting, and paginati... |
| `PUT` | `/api/v1/tenant/payroll/salary-components/{id}` | Update salary component | Update an existing salary components record by its unique ID. Accepts partial updates; only provi... |

### Tenant: Competency Management
**Endpoints:** 36 | **Paths:** 15
**Methods:** DELETE=7 GET=15 POST=7 PUT=7

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/competency/competence-values` | List competence values (legacy) | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/competence-values` | Create competence value (legacy) | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/competence-values/{id}` | Delete competence value | Delete a competence values record by its unique ID. This action may be reversible depending on sy... |
| `GET` | `/api/v1/tenant/competency/competence-values/{id}` | Get competence value by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/competence-values/{id}` | Update competence value | Update an existing competence values record by its unique ID. Accepts partial updates; only provi... |
| `GET` | `/api/v1/tenant/competency/competencies` | List competencies | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/competencies` | Create a new competency | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/competencies/{id}` | Delete competency | Delete a competencies record by its unique ID. This action may be reversible depending on system ... |
| `GET` | `/api/v1/tenant/competency/competencies/{id}` | Get competency by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/competencies/{id}` | Update competency | Update an existing competencies record by its unique ID. Accepts partial updates; only provided f... |
| `GET` | `/api/v1/tenant/competency/event-targets` | List competency event targets | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/event-targets` | Create competency event target | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/event-targets/{id}` | Delete event target | Delete a event targets record by its unique ID. This action may be reversible depending on system... |
| `GET` | `/api/v1/tenant/competency/event-targets/{id}` | Get event target by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/event-targets/{id}` | Update event target | Update an existing event targets record by its unique ID. Accepts partial updates; only provided ... |
| `GET` | `/api/v1/tenant/competency/events` | List competency events | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/events` | Create competency event | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/events/{id}` | Delete competency event | Delete a events record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/competency/events/{id}` | Get competency event by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/events/{id}` | Update competency event | Update an existing events record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/competency/score-details` | Retrieve a list of score details records. | Retrieve a paginated list of score details records. Supports filtering, sorting, and pagination p... |
| `POST` | `/api/v1/tenant/competency/score-details` | Create score detail | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/score-details/{id}` | Delete score detail | Delete a score details record by its unique ID. This action may be reversible depending on system... |
| `GET` | `/api/v1/tenant/competency/score-details/{id}` | Get score detail by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/score-details/{id}` | Update score detail | Update an existing score details record by its unique ID. Accepts partial updates; only provided ... |
| `GET` | `/api/v1/tenant/competency/scores` | List competency scores | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/scores` | Create competency score | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/scores/{id}` | Delete competency score | Delete a scores record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/competency/scores/{id}` | Get competency score by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/scores/{id}` | Update competency score | Update an existing scores record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/competency/scores/{scoreId}/details` | List competency score details | Retrieve a paginated list of competency resources. |
| `GET` | `/api/v1/tenant/competency/values` | List competency values (structured) | Retrieve a paginated list of competency resources. |
| `POST` | `/api/v1/tenant/competency/values` | Create competency value | Create a new competency resource. |
| `DELETE` | `/api/v1/tenant/competency/values/{id}` | Delete competency value | Delete a values record by its unique ID. This action may be reversible depending on system config... |
| `GET` | `/api/v1/tenant/competency/values/{id}` | Get competency value by ID | Retrieve a paginated list of competency resources. |
| `PUT` | `/api/v1/tenant/competency/values/{id}` | Update competency value | Update an existing values record by its unique ID. Accepts partial updates; only provided fields ... |

### Tenant: Performance Management
**Endpoints:** 34 | **Paths:** 17
**Methods:** DELETE=7 GET=12 POST=7 PUT=8

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/performance/evaluation-details` | Create evaluation detail | Add a BSC perspective detail to a performance evaluation, including achievement percentage, weigh... |
| `DELETE` | `/api/v1/tenant/performance/evaluation-details/{id}` | Delete evaluation detail | Permanently delete a BSC perspective detail from the evaluation. |
| `PUT` | `/api/v1/tenant/performance/evaluation-details/{id}` | Update evaluation detail | Update a BSC perspective detail's achievement percentage, weight, score, or description. |
| `GET` | `/api/v1/tenant/performance/evaluations` | List performance evaluations | Retrieve a paginated list of performance evaluations, optionally filtered by employee, organizati... |
| `POST` | `/api/v1/tenant/performance/evaluations` | Create performance evaluation | Start a new performance evaluation for an employee. Links the employee to a performance period an... |
| `DELETE` | `/api/v1/tenant/performance/evaluations/{id}` | Delete performance evaluation | Permanently delete a performance evaluation. Only evaluations in DRAFT status can be deleted. |
| `GET` | `/api/v1/tenant/performance/evaluations/{id}` | Get performance evaluation by ID | Retrieve detailed information about a specific performance evaluation, including its status, scor... |
| `PUT` | `/api/v1/tenant/performance/evaluations/{id}` | Update performance evaluation | Update evaluation metadata such as supervisor assignment or notes. Only provided fields will be u... |
| `GET` | `/api/v1/tenant/performance/evaluations/{id}/details` | List evaluation details by evaluation ID | Retrieve all BSC perspective detail records for a specific performance evaluation, showing achiev... |
| `PUT` | `/api/v1/tenant/performance/evaluations/{id}/status` | Update evaluation status | Transition a performance evaluation through its workflow: DRAFT -> PLAN_SUBMITTED -> PLAN_APPROVE... |
| `GET` | `/api/v1/tenant/performance/evaluations/{id}/targets` | List performance targets by evaluation ID | Retrieve all KPI targets for a specific performance evaluation, showing planned vs actual achieve... |
| `GET` | `/api/v1/tenant/performance/indicators` | List KPI indicators | Retrieve a paginated list of KPI indicators, optionally filtered by template or perspective. |
| `POST` | `/api/v1/tenant/performance/indicators` | Create KPI indicator | Create a new KPI indicator linked to a template and BSC perspective. Defines target value, weight... |
| `DELETE` | `/api/v1/tenant/performance/indicators/{id}` | Delete KPI indicator | Permanently delete a KPI indicator from its template. |
| `GET` | `/api/v1/tenant/performance/indicators/{id}` | Get KPI indicator by ID | Retrieve a specific KPI indicator by its unique ID, including target value and measurement settings. |
| `PUT` | `/api/v1/tenant/performance/indicators/{id}` | Update KPI indicator | Update a KPI indicator's title, weight, target value, or measurement unit. Only provided fields w... |
| `GET` | `/api/v1/tenant/performance/periods` | List performance periods | Retrieve a paginated list of performance periods, optionally filtered by year or status. |
| `POST` | `/api/v1/tenant/performance/periods` | Create performance period | Create a new performance evaluation period (e.g. Q1 2026). Period type must be one of: MONTHLY, Q... |
| `DELETE` | `/api/v1/tenant/performance/periods/{id}` | Delete performance period | Permanently delete a performance period by its unique ID. |
| `GET` | `/api/v1/tenant/performance/periods/{id}` | Get performance period by ID | Retrieve detailed information about a specific performance period by its unique UUID. |
| `PUT` | `/api/v1/tenant/performance/periods/{id}` | Update performance period | Update an existing performance period's details. Only provided fields will be updated. |
| `GET` | `/api/v1/tenant/performance/perspectives` | List BSC perspectives | Retrieve a paginated list of BSC perspectives used in performance templates. Ordered by sort_orde... |
| `POST` | `/api/v1/tenant/performance/perspectives` | Create BSC perspective | Create a new Balanced Scorecard perspective (e.g. Financial, Customer, Internal Process, Learning... |
| `DELETE` | `/api/v1/tenant/performance/perspectives/{id}` | Delete BSC perspective | Permanently delete a BSC perspective from the system. |
| `GET` | `/api/v1/tenant/performance/perspectives/{id}` | Get BSC perspective by ID | Retrieve a specific BSC perspective by its unique ID. |
| `PUT` | `/api/v1/tenant/performance/perspectives/{id}` | Update BSC perspective | Update a BSC perspective's name, description, or sort order. Only provided fields will be updated. |
| `POST` | `/api/v1/tenant/performance/targets` | Create performance target | Add an individual KPI target to a performance evaluation, setting the target value and weight for... |
| `DELETE` | `/api/v1/tenant/performance/targets/{id}` | Delete performance target | Permanently delete a KPI target from the evaluation. |
| `PUT` | `/api/v1/tenant/performance/targets/{id}` | Update performance target | Update a KPI target's planned value, actual achievement, or weight. Setting actual_value triggers... |
| `GET` | `/api/v1/tenant/performance/templates` | List KPI templates | Retrieve a paginated list of KPI templates, optionally filtered by organization. |
| `POST` | `/api/v1/tenant/performance/templates` | Create KPI template | Create a new KPI template for an organization. Templates define the structure of performance eval... |
| `DELETE` | `/api/v1/tenant/performance/templates/{id}` | Delete KPI template | Permanently delete a KPI template. Indicators linked to this template may also be removed. |
| `GET` | `/api/v1/tenant/performance/templates/{id}` | Get KPI template by ID | Retrieve a specific KPI template by its unique ID, including associated indicators. |
| `PUT` | `/api/v1/tenant/performance/templates/{id}` | Update KPI template | Update a KPI template's name, description, or status. Status can be transitioned between DRAFT, P... |

### Tenant: Recruitment & Onboarding (ATS)
**Endpoints:** 33 | **Paths:** 16
**Methods:** DELETE=7 GET=12 POST=7 PUT=7

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/recruitment/applications` | List job applications | Retrieve paginated list of applications, optionally filtered by requisition, candidate, or status |
| `POST` | `/api/v1/tenant/recruitment/applications` | Create job application | Submit a candidate's application to a job requisition. Candidate and requisition must exist |
| `DELETE` | `/api/v1/tenant/recruitment/applications/{id}` | Delete job application | Permanently delete an application record |
| `GET` | `/api/v1/tenant/recruitment/applications/{id}` | Get job application by ID | Retrieve application details including current status and notes |
| `PUT` | `/api/v1/tenant/recruitment/applications/{id}/status` | Update application status | Update application status throughout the recruitment pipeline. Automatically updates requisition ... |
| `GET` | `/api/v1/tenant/recruitment/candidates` | List candidates | Retrieve paginated list of candidates with optional search by name or email |
| `POST` | `/api/v1/tenant/recruitment/candidates` | Create candidate | Register a new candidate. Email must be unique across the system |
| `DELETE` | `/api/v1/tenant/recruitment/candidates/{id}` | Delete candidate | Permanently delete a candidate record |
| `GET` | `/api/v1/tenant/recruitment/candidates/{id}` | Get candidate by ID | Retrieve detailed candidate information including contact details and resume links |
| `PUT` | `/api/v1/tenant/recruitment/candidates/{id}` | Update candidate | Update candidate profile fields. Only provided fields will be updated |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings` | List employee onboardings | Retrieve paginated list of employee onboardings, optionally filtered by status |
| `POST` | `/api/v1/tenant/recruitment/employee-onboardings` | Create employee onboarding | Start onboarding for an accepted candidate. Automatically creates task items from active templates |
| `DELETE` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Delete employee onboarding | Permanently delete an employee onboarding record and its task items |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Get employee onboarding by ID | Retrieve onboarding details including start date, buddy, and current status |
| `PUT` | `/api/v1/tenant/recruitment/employee-onboardings/{id}` | Update employee onboarding | Update onboarding details. Setting status to COMPLETED automatically records completion timestamp |
| `GET` | `/api/v1/tenant/recruitment/employee-onboardings/{id}/task-items` | List onboarding task items | Retrieve all task items for a specific employee onboarding, ordered by due date |
| `GET` | `/api/v1/tenant/recruitment/interviews` | List interviews | Retrieve paginated list of interviews, optionally filtered by application or interviewer |
| `POST` | `/api/v1/tenant/recruitment/interviews` | Create interview | Schedule a new interview for a job application with interviewer, stage, and time slot |
| `DELETE` | `/api/v1/tenant/recruitment/interviews/{id}` | Delete interview | Permanently delete an interview record |
| `GET` | `/api/v1/tenant/recruitment/interviews/{id}` | Get interview by ID | Retrieve interview details including score, feedback, and status |
| `PUT` | `/api/v1/tenant/recruitment/interviews/{id}` | Update interview | Update interview schedule, score, feedback, or status. Setting status to COMPLETED automatically ... |
| `POST` | `/api/v1/tenant/recruitment/onboarding-task-items` | Create onboarding task item | Add a custom task item to an employee onboarding. Can optionally link to a template |
| `DELETE` | `/api/v1/tenant/recruitment/onboarding-task-items/{id}` | Delete onboarding task item | Permanently delete a task item |
| `PUT` | `/api/v1/tenant/recruitment/onboarding-task-items/{id}` | Update onboarding task item | Update task item details. Setting is_completed to true automatically records completion timestamp |
| `GET` | `/api/v1/tenant/recruitment/onboarding-task-templates` | List onboarding task templates | Retrieve paginated list of task templates, optionally filtered by category |
| `POST` | `/api/v1/tenant/recruitment/onboarding-task-templates` | Create onboarding task template | Create a reusable task template for employee onboarding (e.g., IT Setup, Contract Signing) |
| `DELETE` | `/api/v1/tenant/recruitment/onboarding-task-templates/{id}` | Delete onboarding task template | Permanently delete a task template |
| `PUT` | `/api/v1/tenant/recruitment/onboarding-task-templates/{id}` | Update onboarding task template | Update a task template properties |
| `GET` | `/api/v1/tenant/recruitment/requisitions` | List job requisitions | Retrieve paginated list of job requisitions, optionally filtered by organization and status |
| `POST` | `/api/v1/tenant/recruitment/requisitions` | Create job requisition | Create a new job requisition with position details, salary range, and number of slots available |
| `DELETE` | `/api/v1/tenant/recruitment/requisitions/{id}` | Delete job requisition | Permanently delete a job requisition |
| `GET` | `/api/v1/tenant/recruitment/requisitions/{id}` | Get job requisition by ID | Retrieve detailed job requisition information by UUID |
| `PUT` | `/api/v1/tenant/recruitment/requisitions/{id}` | Update job requisition | Update job requisition fields. Only provided fields will be updated |

### Tenant: Time & Attendance
**Endpoints:** 30 | **Paths:** 15
**Methods:** DELETE=4 GET=15 POST=6 PUT=5

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/attendance/employee-shifts` | List employee shift assignments | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/employee-shifts` | Assign a shift to an employee | Create a new employee shifts record. Validates required fields and returns the created resource w... |
| `DELETE` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Delete employee shift assignment | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Get employee shift assignment by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/employee-shifts/{id}` | Update employee shift assignment | Update an attendance record. |
| `GET` | `/api/v1/tenant/attendance/events` | List attendance events (check-in/out) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/events` | Create an attendance event (check-in/out) | Create a new events record. Validates required fields and returns the created resource with its a... |
| `GET` | `/api/v1/tenant/attendance/events/{id}` | Get event by ID | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/exempt-positions` | List exempt positions (positions not requiring attendance) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/exempt-positions` | Create an exempt position | Create a new exempt positions record. Validates required fields and returns the created resource ... |
| `DELETE` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Delete an exempt position | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Get exempt position by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/exempt-positions/{id}` | Update an exempt position | Update an attendance record. |
| `GET` | `/api/v1/tenant/attendance/locations` | List attendance locations (geofence) | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/locations` | Create an attendance location (geofence) | Create a new locations record. Validates required fields and returns the created resource with it... |
| `DELETE` | `/api/v1/tenant/attendance/locations/{id}` | Delete a location | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/locations/{id}` | Get location by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/locations/{id}` | Update a location | Update an attendance record. |
| `GET` | `/api/v1/tenant/attendance/overtime-requests` | List overtime requests | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/overtime-requests` | Create an overtime request | Create a new overtime requests record. Validates required fields and returns the created resource... |
| `GET` | `/api/v1/tenant/attendance/overtime-requests/{id}` | Get overtime request by ID | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/sessions` | List daily work sessions | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/sessions/detail` | Get session detail for an employee on a specific date | Get details of a specific attendance record. |
| `GET` | `/api/v1/tenant/attendance/settings` | Get company attendance settings | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/settings` | Upsert company attendance settings | Update an attendance record. |
| `GET` | `/api/v1/tenant/attendance/shifts` | List company shifts | Get details of a specific attendance record. |
| `POST` | `/api/v1/tenant/attendance/shifts` | Create a company shift | Create a new shifts record. Validates required fields and returns the created resource with its a... |
| `DELETE` | `/api/v1/tenant/attendance/shifts/{id}` | Delete a shift | Remove an attendance record. |
| `GET` | `/api/v1/tenant/attendance/shifts/{id}` | Get shift by ID | Get details of a specific attendance record. |
| `PUT` | `/api/v1/tenant/attendance/shifts/{id}` | Update a shift | Update an attendance record. |

### Tenant: Employees
**Endpoints:** 29 | **Paths:** 18
**Methods:** DELETE=9 GET=2 POST=9 PUT=9

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/employees` | List employees with pagination | Retrieve a paginated list of employees records. Supports filtering, sorting, and pagination param... |
| `POST` | `/api/v1/tenant/employees` | Create a new employee | Create a new employees record. Validates required fields and returns the created resource with it... |
| `DELETE` | `/api/v1/tenant/employees/{id}` | Delete employee (hard delete) | Delete a employees record by its unique ID. This action may be reversible depending on system con... |
| `GET` | `/api/v1/tenant/employees/{id}` | Get employee by ID with all sub-modules | Retrieve a paginated list of employees records. Supports filtering, sorting, and pagination param... |
| `PUT` | `/api/v1/tenant/employees/{id}` | Update employee | Update an existing employees record by its unique ID. Accepts partial updates; only provided fiel... |
| `POST` | `/api/v1/tenant/employees/{id}/addresses` | Add employee address | Add a new address for an employee. Supports address types: MAIN (primary residence) and DOMICILE ... |
| `DELETE` | `/api/v1/tenant/employees/{id}/addresses/{addressId}` | Delete employee address | Remove an employee's address record from the system. |
| `PUT` | `/api/v1/tenant/employees/{id}/addresses/{addressId}` | Update employee address | Update an existing employee address record. Can modify address type, full address details, RT/RW,... |
| `POST` | `/api/v1/tenant/employees/{id}/documents` | Upload employee document | Upload and attach a document to an employee's profile. Supports document types such as ID card (K... |
| `DELETE` | `/api/v1/tenant/employees/{id}/documents/{documentId}` | Delete employee document | Remove a document from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/documents/{documentId}` | Update document metadata | Update an employee document's metadata such as document name, type, or description. |
| `POST` | `/api/v1/tenant/employees/{id}/educations` | Add education record | Add an educational background record for an employee. Includes education level, institution name,... |
| `DELETE` | `/api/v1/tenant/employees/{id}/educations/{educationId}` | Delete education record | Remove an educational background record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/educations/{educationId}` | Update education record | Update an employee's educational background record including institution, degree, graduation date... |
| `POST` | `/api/v1/tenant/employees/{id}/emergency-contacts` | Add emergency contact | Register an emergency contact person for an employee. Includes name, relationship, phone number, ... |
| `DELETE` | `/api/v1/tenant/employees/{id}/emergency-contacts/{contactId}` | Delete emergency contact | Remove an emergency contact record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/emergency-contacts/{contactId}` | Update emergency contact | Update an employee's emergency contact details such as name, relationship, or phone number. |
| `POST` | `/api/v1/tenant/employees/{id}/employments` | Add employment record | Add an employment assignment for an employee. Associates the employee with an organization unit, ... |
| `DELETE` | `/api/v1/tenant/employees/{id}/employments/{employmentId}` | Delete employment record | Remove an employment assignment from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/employments/{employmentId}` | Update employment record | Update an employee's employment assignment including organization, position, status, and decision... |
| `POST` | `/api/v1/tenant/employees/{id}/experiences` | Add work experience | Add a work experience record for an employee. Includes company name, position, start and end date... |
| `DELETE` | `/api/v1/tenant/employees/{id}/experiences/{experienceId}` | Delete work experience | Remove a work experience record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/experiences/{experienceId}` | Update work experience | Update an employee's work experience record including company, position, employment period, and r... |
| `POST` | `/api/v1/tenant/employees/{id}/families` | Add family member | Add a family member record for an employee. Includes name, relationship (spouse, child, parent, s... |
| `DELETE` | `/api/v1/tenant/employees/{id}/families/{familyId}` | Delete family member | Remove a family member record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/families/{familyId}` | Update family member | Update an employee's family member details including relationship, personal information, and tax ... |
| `POST` | `/api/v1/tenant/employees/{id}/insurances` | Add insurance (BPJS) | Register an insurance record for an employee. Typically used for BPJS Kesehatan (health) and BPJS... |
| `DELETE` | `/api/v1/tenant/employees/{id}/insurances/{insuranceId}` | Delete insurance record | Remove an insurance record from the employee's profile. |
| `PUT` | `/api/v1/tenant/employees/{id}/insurances/{insuranceId}` | Update insurance record | Update an employee's insurance record including BPJS participation number, coverage type, and con... |

### Tenant: Leave & Time Off
**Endpoints:** 23 | **Paths:** 12
**Methods:** DELETE=4 GET=11 POST=4 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/leave/accrual-policies` | List accrual policies with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/accrual-policies` | Create an accrual policy | Create a new leave resource. |
| `DELETE` | `/api/v1/tenant/leave/accrual-policies/{id}` | Delete accrual policy | Delete a accrual policies record by its unique ID. This action may be reversible depending on sys... |
| `GET` | `/api/v1/tenant/leave/accrual-policies/{id}` | Get accrual policy by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/accrual-policies/{id}` | Update accrual policy | Update an existing accrual policies record by its unique ID. Accepts partial updates; only provid... |
| `GET` | `/api/v1/tenant/leave/balances` | List leave balances with pagination | Retrieve a paginated list of leave resources. |
| `GET` | `/api/v1/tenant/leave/balances/employees/{employeeId}/types/{leaveTypeId}` | Get leave balance for specific employee and leave type | Retrieve a paginated list of leave resources. |
| `GET` | `/api/v1/tenant/leave/reasons` | List leave reasons | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/reasons` | Create a leave reason | Create a new leave resource. |
| `DELETE` | `/api/v1/tenant/leave/reasons/{id}` | Delete leave reason | Delete a reasons record by its unique ID. This action may be reversible depending on system confi... |
| `GET` | `/api/v1/tenant/leave/reasons/{id}` | Get leave reason by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/reasons/{id}` | Update leave reason | Update an existing reasons record by its unique ID. Accepts partial updates; only provided fields... |
| `GET` | `/api/v1/tenant/leave/requests` | List leave requests with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/requests` | Create a leave request | Create a new leave resource. |
| `DELETE` | `/api/v1/tenant/leave/requests/{id}` | Delete leave request | Delete a requests record by its unique ID. This action may be reversible depending on system conf... |
| `GET` | `/api/v1/tenant/leave/requests/{id}` | Get leave request by ID | Retrieve a paginated list of leave resources. |
| `GET` | `/api/v1/tenant/leave/requests/{id}/details` | List leave request details (daily breakdown) | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/requests/{id}/status` | Update leave request status (approve/reject/cancel) | Update an existing status record by its unique ID. Accepts partial updates; only provided fields ... |
| `GET` | `/api/v1/tenant/leave/types` | List leave types with pagination | Retrieve a paginated list of leave resources. |
| `POST` | `/api/v1/tenant/leave/types` | Create a new leave type | Create a new leave resource. |
| `DELETE` | `/api/v1/tenant/leave/types/{id}` | Delete leave type | Delete a types record by its unique ID. This action may be reversible depending on system configu... |
| `GET` | `/api/v1/tenant/leave/types/{id}` | Get leave type by ID | Retrieve a paginated list of leave resources. |
| `PUT` | `/api/v1/tenant/leave/types/{id}` | Update leave type | Update an existing types record by its unique ID. Accepts partial updates; only provided fields w... |

### Tenant: Approval
**Endpoints:** 15 | **Paths:** 9
**Methods:** DELETE=2 GET=6 POST=5 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/approval/flows` | List approval flows | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/flows` | Create approval flow | Create a new approval resource. |
| `DELETE` | `/api/v1/tenant/approval/flows/{flowId}` | Delete approval flow | Delete a flows record by its unique ID. This action may be reversible depending on system configu... |
| `GET` | `/api/v1/tenant/approval/flows/{flowId}` | Get approval flow by ID | Retrieve a paginated list of approval resources. |
| `PUT` | `/api/v1/tenant/approval/flows/{flowId}` | Update approval flow | Update an existing flows record by its unique ID. Accepts partial updates; only provided fields w... |
| `GET` | `/api/v1/tenant/approval/flows/{flowId}/steps` | List approval flow steps | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/flows/{flowId}/steps` | Create approval flow step | Create a new approval resource. |
| `DELETE` | `/api/v1/tenant/approval/flows/{flowId}/steps/{stepId}` | Delete approval flow step | Delete a steps record by its unique ID. This action may be reversible depending on system configu... |
| `PUT` | `/api/v1/tenant/approval/flows/{flowId}/steps/{stepId}` | Update approval flow step | Update an existing steps record by its unique ID. Accepts partial updates; only provided fields w... |
| `GET` | `/api/v1/tenant/approval/instances` | List approval instances | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/instances` | Create approval instance | Create a new approval resource. |
| `GET` | `/api/v1/tenant/approval/instances/{id}` | Get approval instance by ID | Retrieve a paginated list of approval resources. |
| `POST` | `/api/v1/tenant/approval/instances/{id}/actions` | Submit approval action (approve/reject) | Create a new approval resource. |
| `POST` | `/api/v1/tenant/approval/instances/{id}/cancel` | Cancel approval instance | Cancel an active approval instance. This will void all pending tasks and mark the instance as CAN... |
| `GET` | `/api/v1/tenant/approval/tasks/pending` | List my pending approval tasks | Retrieve a paginated list of approval resources. |

### Tenant: Employee Movement & Career Management
**Endpoints:** 15 | **Paths:** 9
**Methods:** DELETE=2 GET=6 POST=5 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/employee-movements/contracts` | List employee contracts | Retrieve a paginated list of employee movement resources. |
| `POST` | `/api/v1/tenant/employee-movements/contracts` | Create employee contract | Create a new employee movement resource. |
| `DELETE` | `/api/v1/tenant/employee-movements/contracts/{id}` | Delete contract | Delete a contracts record by its unique ID. This action may be reversible depending on system con... |
| `GET` | `/api/v1/tenant/employee-movements/contracts/{id}` | Get contract by ID | Retrieve a paginated list of employee movement resources. |
| `PUT` | `/api/v1/tenant/employee-movements/contracts/{id}` | Update contract | Update an existing contracts record by its unique ID. Accepts partial updates; only provided fiel... |
| `GET` | `/api/v1/tenant/employee-movements/employees/{employeeId}/contracts` | List contracts by employee | Retrieve a paginated list of employee movement resources. |
| `GET` | `/api/v1/tenant/employee-movements/employees/{employeeId}/movements` | List movements by employee | Retrieve a paginated list of employee movement resources. |
| `GET` | `/api/v1/tenant/employee-movements/movements` | List employee movements | Retrieve a paginated list of employee movement resources. |
| `POST` | `/api/v1/tenant/employee-movements/movements` | Create employee movement | Create a new employee movement resource. |
| `DELETE` | `/api/v1/tenant/employee-movements/movements/{id}` | Delete movement | Delete a movements record by its unique ID. This action may be reversible depending on system con... |
| `GET` | `/api/v1/tenant/employee-movements/movements/{id}` | Get movement by ID | Retrieve a paginated list of employee movement resources. |
| `PUT` | `/api/v1/tenant/employee-movements/movements/{id}` | Update movement | Update an existing movements record by its unique ID. Accepts partial updates; only provided fiel... |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/approve` | Approve movement | Create a new employee movement resource. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/cancel` | Cancel movement | Create a new employee movement resource. |
| `POST` | `/api/v1/tenant/employee-movements/movements/{id}/execute` | Execute movement | Create a new employee movement resource. |

### Tenant: Reimbursement & Claim
**Endpoints:** 15 | **Paths:** 7
**Methods:** DELETE=3 GET=5 POST=3 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/reimbursements/requests` | List reimbursement requests | Retrieve a paginated list of reimbursement requests. Supports filtering by employee, status (DRAF... |
| `POST` | `/api/v1/tenant/reimbursements/requests` | Create reimbursement request | Submit a new reimbursement request for approval. Includes the reimbursement type, title, descript... |
| `DELETE` | `/api/v1/tenant/reimbursements/requests/{id}` | Delete reimbursement request | Delete a reimbursement request. Only requests in DRAFT or CANCELLED status can be deleted. Finali... |
| `GET` | `/api/v1/tenant/reimbursements/requests/{id}` | Get reimbursement request | Get detailed information about a specific reimbursement request, including its items, status hist... |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{id}` | Update reimbursement request | Update a reimbursement request. Only requests in DRAFT status can be modified. Changes to title, ... |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{id}/status` | Update reimbursement request status | Transition a reimbursement request through its approval workflow. Valid transitions: DRAFT->SUBMI... |
| `GET` | `/api/v1/tenant/reimbursements/requests/{requestId}/items` | List reimbursement items | Retrieve all expense items attached to a specific reimbursement request. Returns item details inc... |
| `POST` | `/api/v1/tenant/reimbursements/requests/{requestId}/items` | Add reimbursement item | Add a new expense item to a reimbursement request. Each item represents a single expense with dat... |
| `DELETE` | `/api/v1/tenant/reimbursements/requests/{requestId}/items/{itemId}` | Delete reimbursement item | Remove an expense item from a reimbursement request. The total request amount will be recalculate... |
| `PUT` | `/api/v1/tenant/reimbursements/requests/{requestId}/items/{itemId}` | Update reimbursement item | Update an expense item's details including expense date, type, amount, description, or receipt UR... |
| `GET` | `/api/v1/tenant/reimbursements/types` | List reimbursement types | Retrieve a paginated list of reimbursement type configurations. Supports filtering by name or cat... |
| `POST` | `/api/v1/tenant/reimbursements/types` | Create reimbursement type | Create a new reimbursement type for the company. Defines the category name, maximum claimable amo... |
| `DELETE` | `/api/v1/tenant/reimbursements/types/{id}` | Delete reimbursement type | Delete a reimbursement type from the system. Existing requests using this type will retain their ... |
| `GET` | `/api/v1/tenant/reimbursements/types/{id}` | Get reimbursement type | Get detailed information about a specific reimbursement type, including its name, maximum amount,... |
| `PUT` | `/api/v1/tenant/reimbursements/types/{id}` | Update reimbursement type | Update a reimbursement type's name, maximum amount, or description. Changes apply to new requests... |

### Platform: Companies
**Endpoints:** 10 | **Paths:** 7
**Methods:** DELETE=1 GET=2 POST=6 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/companies` | List all companies | Retrieve a paginated list of all registered companies (tenants) in the platform. Includes company... |
| `POST` | `/api/v1/platform/companies` | Create a new company/tenant | Register a new company tenant. This triggers the full tenant provisioning flow: creates the compa... |
| `DELETE` | `/api/v1/platform/companies/{id}` | Soft delete company (deactivate connection + deleted_at) | Soft-delete a company tenant. Deactivates the tenant database connection and sets the deleted_at ... |
| `GET` | `/api/v1/platform/companies/{id}` | Get company by ID | Get detailed information about a specific company/tenant including its status, contact details, s... |
| `PUT` | `/api/v1/platform/companies/{id}` | Update company | Update a company's profile information including name, email, phone, address, and other contact d... |
| `POST` | `/api/v1/platform/companies/{id}/activate` | Activate a company/tenant (reactivate connection) | Reactivate a previously suspended company tenant. Re-establishes the database connection and sets... |
| `POST` | `/api/v1/platform/companies/{id}/backup` | Trigger tenant backup (Phase 2) | Trigger an on-demand database backup for the specified company tenant. The backup is stored accor... |
| `POST` | `/api/v1/platform/companies/{id}/restore` | Trigger tenant restore (Phase 2) | Restore a company tenant's database from a previously created backup. Requires a valid backup ref... |
| `POST` | `/api/v1/platform/companies/{id}/suspend` | Suspend a company/tenant (deactivate connection) | Suspend a company tenant — deactivates the database connection, clears cached connections, and se... |
| `POST` | `/api/v1/platform/companies/{id}/terminate` | Terminate a company/tenant (drop database + remove connection) | Permanently terminate a company tenant. This drops the tenant database entirely, removes the conn... |

### Platform: RBAC Management
**Endpoints:** 10 | **Paths:** 6
**Methods:** DELETE=3 GET=3 POST=3 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/rbac/permissions` | List all permissions (resource + action) | Retrieve all available permissions in the system. Permissions are defined as resource.action pair... |
| `POST` | `/api/v1/platform/rbac/permissions` | Create a new permission (resource.action) | Create a new permission in the format resource.action (e.g., report.view). Permissions can be ass... |
| `DELETE` | `/api/v1/platform/rbac/permissions/{id}` | Delete a permission (non-system only) | Delete a non-system permission. The permission will be removed from all role assignments. System ... |
| `GET` | `/api/v1/platform/rbac/roles` | List all roles with their permissions | Retrieve all RBAC roles with their associated permissions. Roles are organized in a hierarchy wit... |
| `POST` | `/api/v1/platform/rbac/roles` | Create a new role | Create a new RBAC role with a name, description, and optional parent role. New roles inherit perm... |
| `DELETE` | `/api/v1/platform/rbac/roles/{id}` | Delete a role (non-system roles only) | Delete a non-system role. Users assigned to this role will lose their associated permissions unti... |
| `GET` | `/api/v1/platform/rbac/roles/{id}` | Get role by ID with permissions | Get detailed information about a specific role including its name, description, parent role (if a... |
| `PUT` | `/api/v1/platform/rbac/roles/{id}` | Update role (name, description, parent) | Update a role's name, description, or parent role assignment. Changes to system roles may be rest... |
| `POST` | `/api/v1/platform/rbac/roles/{id}/permissions` | Assign a permission to a role | Assign a permission to a role. The permission becomes available to all users with that role (and ... |
| `DELETE` | `/api/v1/platform/rbac/roles/{id}/permissions/{permissionId}` | Revoke a permission from a role | Revoke a permission from a role. The permission will no longer be available to users with that ro... |

### Platform: Modules
**Endpoints:** 7 | **Paths:** 5
**Methods:** GET=3 POST=3 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/modules` | List all registered modules | Retrieve a paginated list of all registered system modules. Each module represents a functional a... |
| `POST` | `/api/v1/platform/modules` | Register a new module | Register a new system module with its name, slug, version, and optional dependencies. Modules can... |
| `GET` | `/api/v1/platform/modules/{id}` | Get module by ID | Get detailed information about a specific module including its version, dependencies, and activat... |
| `PUT` | `/api/v1/platform/modules/{id}` | Update module | Update a module's configuration, metadata, version, or feature flags. Changes apply globally acro... |
| `POST` | `/api/v1/platform/modules/{id}/activate` | Activate module for a company | Activate a module for a specific company tenant. The module's features become available in that t... |
| `GET` | `/api/v1/platform/modules/{id}/companies` | List companies using this module | Retrieve a list of companies that have this module activated. Shows activation date and module-sp... |
| `POST` | `/api/v1/platform/modules/{id}/deactivate` | Deactivate module for a company | Deactivate a module for a specific company tenant. The module's features are hidden from that ten... |

### Tenant: Organizations
**Endpoints:** 12 | **Paths:** 8
**Methods:** DELETE=1 GET=6 POST=4 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/organizations` | List organizations or get tree | Get detailed information about an organizational unit and its children. |
| `POST` | `/api/v1/tenant/organizations` | Create organization | Create a new organizations record. Validates required fields and returns the created resource wit... |
| `DELETE` | `/api/v1/tenant/organizations/{id}` | Delete organization | Remove an organizational unit from the hierarchy. |
| `GET` | `/api/v1/tenant/organizations/{id}` | Get organization by ID | Get detailed information about an organizational unit and its children. |
| `PUT` | `/api/v1/tenant/organizations/{id}` | Update organization | Update organizational unit code, name, or parent assignment. |
| `GET` | `/api/v1/tenant/organizations/history` | List organization change history | Retrieve a paginated audit trail of all changes made to the organization structure. Supports filt... |
| `POST` | `/api/v1/tenant/organizations/versions` | Create organization version snapshot | Capture a point-in-time snapshot of the full organization tree as a version for comparison or res... |
| `GET` | `/api/v1/tenant/organizations/versions` | List organization versions | Retrieve a paginated list of organization version snapshots with their status and timestamp. |
| `GET` | `/api/v1/tenant/organizations/versions/{id}` | Get organization version by ID | Retrieve a specific version snapshot's metadata. Use ?snapshot=true to include the full tree data. |
| `GET` | `/api/v1/tenant/organizations/versions/{id}/diff/{targetId}` | Compare two organization versions | Compare two version snapshots and identify ADDED, MODIFIED, and REMOVED organizational nodes between them. |
| `POST` | `/api/v1/tenant/organizations/versions/{id}/restore` | Restore organization from version | Restore the entire organization tree to the state captured in a version snapshot. This clears all... |
| `POST` | `/api/v1/tenant/organizations/clone` | Clone current organization tree to a draft version | Create a DRAFT version snapshot of the current organization tree for restructuring simulation. The... |

### Tenant: Training & Development Management
**Endpoints:** 35 | **Paths:** 15
**Methods:** DELETE=7 GET=13 POST=7 PUT=8

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/trainings/categories` | Create training category | Create a new training category (e.g. Technical, Soft Skill, Leadership, Compliance). Categories are used to group related training courses. |
| `GET` | `/api/v1/tenant/trainings/categories` | List training categories | Retrieve a paginated list of training categories, ordered by code. |
| `GET` | `/api/v1/tenant/trainings/categories/{id}` | Get training category by ID | Retrieve detailed information about a specific training category including its code, name, description, and active status. |
| `PUT` | `/api/v1/tenant/trainings/categories/{id}` | Update training category | Update an existing training category. Fields that are not provided will remain unchanged. |
| `DELETE` | `/api/v1/tenant/trainings/categories/{id}` | Delete training category | Soft-delete a training category. The category will be marked as deleted but retained in the database for historical purposes. |
| `POST` | `/api/v1/tenant/trainings/courses` | Create training course | Create a new training course under a specific category. Each course has a unique code and can include duration, cost, and minimum score requirements. |
| `GET` | `/api/v1/tenant/trainings/courses` | List training courses | Retrieve a paginated list of training courses. Optionally filter by category_id. |
| `GET` | `/api/v1/tenant/trainings/courses/{id}` | Get training course by ID | Retrieve detailed information about a specific training course including category, duration, cost, and certification settings. |
| `PUT` | `/api/v1/tenant/trainings/courses/{id}` | Update training course | Update an existing training course. Only provided fields will be updated. |
| `DELETE` | `/api/v1/tenant/trainings/courses/{id}` | Delete training course | Soft-delete a training course. Associated sessions and materials are not deleted. |
| `POST` | `/api/v1/tenant/trainings/sessions` | Create training session | Schedule a new training session/class for a course. Defines the trainer, date range, location, and maximum participant quota. |
| `GET` | `/api/v1/tenant/trainings/sessions` | List training sessions | Retrieve a paginated list of training sessions. Supports filtering by course_id and status. |
| `GET` | `/api/v1/tenant/trainings/sessions/{id}` | Get training session by ID | Retrieve detailed information about a specific training session. |
| `PUT` | `/api/v1/tenant/trainings/sessions/{id}` | Update training session | Update an existing training session's schedule, trainer, location, or quota. |
| `PUT` | `/api/v1/tenant/trainings/sessions/{id}/status` | Update training session status | Transition a training session through its lifecycle: SCHEDULED -> IN_PROGRESS -> COMPLETED or CANCELLED. |
| `DELETE` | `/api/v1/tenant/trainings/sessions/{id}` | Delete training session | Soft-delete a training session. |
| `POST` | `/api/v1/tenant/trainings/participants` | Register training participant | Register an employee as a participant in a training session. Validates that the session is not cancelled and has available quota. |
| `GET` | `/api/v1/tenant/trainings/participants` | List training participants | Retrieve a paginated list of training participants. Filter by session_id or employee_id. |
| `GET` | `/api/v1/tenant/trainings/participants/{id}` | Get training participant by ID | Retrieve participant details including attendance status, score, and completion date. |
| `PUT` | `/api/v1/tenant/trainings/participants/{id}` | Update training participant | Update a participant's attendance status or score. Setting a score auto-marks completion. |
| `DELETE` | `/api/v1/tenant/trainings/participants/{id}` | Delete training participant | Remove a participant from a training session. |
| `POST` | `/api/v1/tenant/trainings/materials` | Create training material | Add a new material (file, document, or resource) to a training session. |
| `GET` | `/api/v1/tenant/trainings/materials` | List training materials by session | List all materials attached to a training session. Requires session_id as a query parameter. |
| `PUT` | `/api/v1/tenant/trainings/materials/{id}` | Update training material | Update a training material's title, file URL, file type, or sort order. |
| `DELETE` | `/api/v1/tenant/trainings/materials/{id}` | Delete training material | Remove a training material from the session. |
| `POST` | `/api/v1/tenant/trainings/evaluations` | Create training evaluation | Submit a training evaluation/feedback for a session. Rating must be between 1 and 5. |
| `GET` | `/api/v1/tenant/trainings/evaluations` | List training evaluations | Retrieve a paginated list of training evaluations. Filter by session_id or employee_id. |
| `GET` | `/api/v1/tenant/trainings/evaluations/{id}` | Get training evaluation by ID | Retrieve a specific training evaluation including rating and feedback details. |
| `PUT` | `/api/v1/tenant/trainings/evaluations/{id}` | Update training evaluation | Update a training evaluation's rating or feedback text. |
| `DELETE` | `/api/v1/tenant/trainings/evaluations/{id}` | Delete training evaluation | Remove a training evaluation from the system. |
| `POST` | `/api/v1/tenant/trainings/certificates` | Issue training certificate | Issue a certificate to a training participant. Requires a unique certificate number. |
| `GET` | `/api/v1/tenant/trainings/certificates` | List training certificates | Retrieve a paginated list of issued certificates. Filter by participant_id. |
| `GET` | `/api/v1/tenant/trainings/certificates/{id}` | Get training certificate by ID | Retrieve a specific training certificate by its unique ID. |
| `PUT` | `/api/v1/tenant/trainings/certificates/{id}` | Update training certificate | Update a training certificate's number, issued date, or expiry date. |
| `DELETE` | `/api/v1/tenant/trainings/certificates/{id}` | Delete training certificate | Delete a training certificate record. |

### Health
**Endpoints:** 4 | **Paths:** 4
**Methods:** GET=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/docs` | Scalar API Documentation UI | Interactive API documentation powered by Scalar UI. Browse all available endpoints, test requests... |
| `GET` | `/healthz` | Server health check | Simple health check endpoint for load balancer probes. Returns HTTP 200 with service status when ... |
| `GET` | `/openapi.json` | OpenAPI 3.0 Specification | Download the complete OpenAPI 3.0 specification as JSON. Compatible with tools like Postman, Inso... |
| `GET` | `/readyz` | Readiness check | Readiness check endpoint for Kubernetes or container orchestration probes. Returns HTTP 200 when ... |

### Platform: Licenses
**Endpoints:** 4 | **Paths:** 2
**Methods:** GET=2 POST=1 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/licenses` | List all licenses | Retrieve a paginated list of all license records. Licenses define the plan type, feature entitlem... |
| `POST` | `/api/v1/platform/licenses` | Create a new license for company | Issue a new software license to a company tenant. Specifies the plan type (trial, basic, professi... |
| `GET` | `/api/v1/platform/licenses/{id}` | Get license by ID | Get detailed information about a specific license including plan type, validity period, seat usag... |
| `PUT` | `/api/v1/platform/licenses/{id}` | Update license | Update license terms including plan upgrade/downgrade, extension of validity period, seat count a... |

### Platform: Users
**Endpoints:** 4 | **Paths:** 2
**Methods:** GET=2 POST=1 PUT=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/users` | List all platform users | Retrieve a paginated list of platform user accounts. Supports filtering by company, role, and sea... |
| `POST` | `/api/v1/platform/users` | Create a new platform user | Register a new platform user account. Assigns the user to a company with a specific role (super_a... |
| `GET` | `/api/v1/platform/users/{id}` | Get platform user by ID | Get detailed information about a specific platform user including their assigned company, role, a... |
| `PUT` | `/api/v1/platform/users/{id}` | Update platform user | Update a platform user's profile information such as name, email, or role assignment. Cannot chan... |

### Platform: Monitoring
**Endpoints:** 3 | **Paths:** 3
**Methods:** GET=3

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/platform/monitoring/health` | Platform health check with database status | Platform health check endpoint providing detailed status of database connectivity, Redis cache he... |
| `GET` | `/api/v1/platform/monitoring/tenants` | List all active tenant connections health | List all active tenant database connections with their health status, pool statistics (open/idle ... |
| `GET` | `/api/v1/platform/monitoring/tenants/{id}` | Get tenant connection health detail | Get detailed health information for a specific tenant database connection, including connection p... |

### Platform: Auth
**Endpoints:** 2 | **Paths:** 2
**Methods:** POST=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/platform/login` | Platform admin login | Authenticate a platform admin user with email and password credentials. Returns a JWT access toke... |
| `POST` | `/api/v1/platform/refresh` | Refresh access token | Exchange a valid refresh token for a new access token. Use this endpoint to maintain session cont... |

### Tenant: Approval Engine
**Endpoints:** 1 | **Paths:** 1
**Methods:** POST=1

| Method | Path | Summary | Description |
|---|---|---|---|
| `POST` | `/api/v1/tenant/approval/instances/{id}` | Cancel approval instance | Create a new instances record. Validates required fields and returns the created resource with it... |

## 3. Recent Improvements

### Tenant: Workforce Intelligence & Strategic Planning
**Endpoints:** 68 | **Paths:** 58
**Methods:** DELETE=3 GET=56 POST=5 PUT=4

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/headcounts` | List headcount plans | Retrieve a paginated list of headcount plans. Filter by period and organization unit. |
| `POST` | `/api/v1/tenant/workforce-intelligence/planning/headcounts` | Create headcount plan | Create a new headcount plan for a specific period and organization. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Get headcount plan by ID | Get detailed information about a specific headcount plan. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Update headcount plan | Update an existing headcount plan. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/planning/headcounts/{id}` | Delete headcount plan | Remove a headcount plan record. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/forecasts` | List workforce forecasts | Retrieve a paginated list of workforce forecasts. |
| `POST` | `/api/v1/tenant/workforce-intelligence/planning/forecasts` | Create workforce forecast | Create a new workforce forecast. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Get forecast by ID | Get detailed information about a specific workforce forecast. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Update forecast | Update an existing workforce forecast. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/planning/forecasts/{id}` | Delete forecast | Remove a workforce forecast record. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/gap-analysis` | Workforce gap analysis | Analyze workforce gaps by comparing supply vs demand. |
| `GET` | `/api/v1/tenant/workforce-intelligence/planning/projections` | Workforce projections | Get workforce projections including hiring needs and retirement counts. |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi` | List KPIs | Retrieve a paginated list of workforce KPIs. |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi/summary` | KPI summary | Get a summary of KPIs showing on-target vs below-target counts. |
| `GET` | `/api/v1/tenant/workforce-intelligence/kpi/{code}` | Get KPI by code | Get a specific KPI by its unique code. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/headcount` | Headcount analytics dashboard | Analyze workforce composition. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/attendance` | Attendance analytics dashboard | Analyze attendance metrics. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/leave` | Leave analytics dashboard | Analyze leave utilization. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/overtime` | Overtime analytics dashboard | Analyze overtime patterns. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/payroll` | Payroll analytics dashboard | Analyze payroll metrics. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/performance` | Performance analytics dashboard | Analyze employee performance. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/learning` | Learning analytics dashboard | Analyze learning and development. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/recruitment` | Recruitment analytics dashboard | Analyze recruitment efficiency. |
| `GET` | `/api/v1/tenant/workforce-intelligence/analytics/movement` | Movement analytics dashboard | Analyze employee movement. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/dashboard` | Capacity dashboard | Get workforce capacity dashboard. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/utilization` | Resource utilization rate | Get workforce utilization rate. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/forecast` | Capacity forecast | Get projected capacity forecast. |
| `GET` | `/api/v1/tenant/workforce-intelligence/capacity/bottlenecks` | Bottleneck analysis | Identify capacity bottlenecks across departments. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/summary` | Cost summary dashboard | Get workforce cost summary. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/payroll` | Payroll cost breakdown | Get detailed payroll cost breakdown. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/per-employee` | Cost per employee analysis | Get cost per employee metrics. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/per-department` | Cost by department | Get workforce cost by department. |
| `GET` | `/api/v1/tenant/workforce-intelligence/cost/budget-vs-actual` | Budget vs actual cost comparison | Get budget vs actual cost comparison. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/dashboard` | Risk dashboard | Get risk dashboard overview. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/indicators` | List risk indicators | Retrieve a paginated list of risk indicators. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/indicators/{id}` | Get risk indicator by ID | Get detailed information about a specific risk indicator. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/risk/indicators/{id}` | Update risk indicator | Update a risk indicator level and recommendation. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/high-turnover` | High turnover risk analysis | Analyze high turnover risk. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/retirement` | Retirement risk analysis | Analyze retirement risk. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/contract-expiry` | Contract expiration risk analysis | Analyze contract expiry risk. |
| `GET` | `/api/v1/tenant/workforce-intelligence/risk/high-absenteeism` | High absenteeism risk analysis | Analyze high absenteeism risk. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/summary` | Executive workforce summary | Executive dashboard summary. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/growth` | Executive workforce growth trend | Executive-level growth trend analysis. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/cost-trend` | Executive cost trend | Executive-level cost trend analysis. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/attrition-trend` | Executive attrition trend | Executive-level attrition rate trend. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/capacity` | Executive capacity overview | Executive-level capacity overview. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/hiring-progress` | Hiring progress tracker | Track hiring progress. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/risk-overview` | Executive risk overview | Executive-level risk overview. |
| `GET` | `/api/v1/tenant/workforce-intelligence/executive/health-score` | Executive health score | Executive-level organization health score. |
| `GET` | `/api/v1/tenant/workforce-intelligence/scenarios` | List scenarios | Retrieve a paginated list of saved simulation scenarios. |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios` | Create scenario | Create a new scenario for workforce simulation. |
| `GET` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Get scenario by ID | Get detailed information about a specific scenario. |
| `PUT` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Update scenario | Update an existing scenario. |
| `DELETE` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}` | Delete scenario | Soft-delete a scenario by ID. |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}/run` | Run scenario simulation | Execute a scenario simulation. |
| `POST` | `/api/v1/tenant/workforce-intelligence/scenarios/{id}/clone` | Clone scenario | Clone an existing scenario. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/dashboard` | Organization health dashboard | Get organization health dashboard. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/scores` | List health scores | Retrieve a paginated list of health scores. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/scores/{id}` | Get health score by ID | Get detailed health score components. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/span-of-control` | Span of control analysis | Analyze span of control. |
| `GET` | `/api/v1/tenant/workforce-intelligence/health/succession` | Succession readiness analysis | Analyze succession readiness. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/training-vs-performance` | Training vs performance correlation | Analyze correlation between training and performance. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/overtime-vs-productivity` | Overtime vs productivity correlation | Analyze correlation between overtime and productivity. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/attendance-vs-performance` | Attendance vs performance correlation | Analyze correlation between attendance and performance. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/compensation-vs-turnover` | Compensation vs turnover correlation | Analyze correlation between compensation and turnover. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/source-vs-retention` | Recruitment source vs retention correlation | Analyze correlation between recruitment source and retention. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/career-progression` | Career progression vs performance correlation | Analyze correlation between career progression and performance. |
| `GET` | `/api/v1/tenant/workforce-intelligence/people-analytics/learning-effectiveness` | Learning effectiveness analysis | Analyze learning program effectiveness. |

### Tenant: Career Intelligence
**Endpoints:** 19 | **Paths:** 10
**Methods:** DELETE=3 GET=10 POST=4 PUT=2

| Method | Path | Summary | Description |
|---|---|---|---|
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps` | List talent maps | Retrieve a paginated list of talent mapping records showing employee placement in the 9-box grid. Filters by organization, position, or competency cluster. |
| `POST` | `/api/v1/tenant/career-intelligence/talent-maps` | Create talent mapping entry | Create a new talent mapping assessment for an employee, placing them in the performance-potential 9-box grid. |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/grid` | Get talent map grid overview | Retrieve the 9-box talent grid summary showing distribution of employees across performance-potential quadrants. |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/employee/{employeeId}` | Get employee talent profile | Retrieve a specific employee's talent mapping profile including current grid position, historical changes, and recommended career paths. |
| `GET` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Get talent map by ID | Get detailed information about a specific talent mapping entry. |
| `PUT` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Update talent map entry | Update an existing talent map entry's performance rating, potential rating, and/or notes. Changes automatically recalculate the grid position. |
| `DELETE` | `/api/v1/tenant/career-intelligence/talent-maps/{id}` | Delete talent map entry | Soft-delete a talent map entry. |
| `GET` | `/api/v1/tenant/career-intelligence/interests` | List career interests | Retrieve a paginated list of career interest records, optionally filtered by employee or interest category. |
| `POST` | `/api/v1/tenant/career-intelligence/interests` | Record career interest | Create a new career interest record for an employee (e.g., leadership track, specialist track, international assignment). |
| `GET` | `/api/v1/tenant/career-intelligence/interests/employee/{employeeId}` | Get employee career interests | Retrieve all recorded career interests and aspirations for a specific employee. |
| `GET` | `/api/v1/tenant/career-intelligence/paths` | List career paths | Retrieve a paginated list of career path templates defining possible progression routes between positions. |
| `POST` | `/api/v1/tenant/career-intelligence/paths` | Create career path | Create a new career path template linking source and target positions with competency requirements and typical duration. |
| `DELETE` | `/api/v1/tenant/career-intelligence/paths/{id}` | Delete career path | Remove a career path template by its unique ID. |
| `GET` | `/api/v1/tenant/career-intelligence/paths/gap-analysis` | Run career gap analysis | Analyze competency gaps for an employee against a target position or career path. Returns skill gaps, recommended training, and estimated readiness timeline. |
| `GET` | `/api/v1/tenant/career-intelligence/successions` | List succession plans | Retrieve a paginated list of succession plans, optionally filtered by position or readiness status. |
| `POST` | `/api/v1/tenant/career-intelligence/successions` | Create succession plan | Create a succession plan designating a successor (employee) for a critical position with readiness level and target date. |
| `GET` | `/api/v1/tenant/career-intelligence/successions/{id}` | Get succession plan by ID | Get detailed information about a specific succession plan including successor details, readiness level, and development plan. |
| `PUT` | `/api/v1/tenant/career-intelligence/successions/{id}` | Update succession plan | Update the readiness level, target date, or notes for an existing succession plan entry. |
| `DELETE` | `/api/v1/tenant/career-intelligence/successions/{id}` | Delete succession plan | Remove a succession plan record by its unique ID. |

### v11 (Current — Career Intelligence Expansion)
- **Career Intelligence** — Expanded from 15 to 19 endpoints (4 new CRUD operations added)
- **New Endpoints (4):** GET/PUT/DELETE /talent-maps/{id} and GET /successions/{id}
- **New Schemas:** CreateTalentMapRequest, UpdateTalentMapRequest, CreateCareerPathRequest, CreateSuccessionPlanRequest, UpdateSuccessionPlanRequest, CareerGapAnalysisResponse, CareerGapRecommendation, PaginatedResponseCI
- **Stats update:** 540 to 544 endpoints (+4), 348 to 352 schemas (+4), 290 to 300 paths (+10), 24 tags (unchanged)

### v10 (25 Juli 2026)
- **Career Intelligence** — New standalone module (15 endpoints, 14 schemas, 7 paths)
- **New Endpoints:** 15 endpoints across 4 resource groups under `/career-intelligence/` prefix: Talent Maps (4), Career Interests (3), Career Paths (3), Succession Plans (5)
- **New Schemas:** 14 new schemas: TalentMappingRequest, TalentMappingResponse, TalentGridResponse, EmployeeTalentProfileResponse, CareerInterestRequest, CareerInterestResponse, CareerPathRequest, CareerPathResponse, GapAnalysisRequest, CareerGapAnalysisResponse, SuccessionPlanRequest, SuccessionPlanResponse, PaginatedCareerResponse, CareerPathListResponse
- **Stats update:** 525 to 540 endpoints, 334 to 348 schemas, 290 to 290 paths, 23 to 24 tags
- **New Tag:** Tenant: Career Intelligence

### v9 (25 Juli 2026)
- **Workforce Intelligence & Strategic Planning** — New strategic analytics module (68 endpoints, 44 schemas, 58 paths)
- **New Endpoints:** 68 endpoints across 10 groups: Planning (12), Analytics (9), Risk (8), Executive (8), Scenarios (7), People Analytics (7), Cost (5), Capacity (4), Health (5), KPI (3)
- **New Schemas:** 44 new schemas: HeadcountPlan, Forecast, GapAnalysis/Projection, KPI, 9 Analytics types, Capacity, Cost, Risk Dashboard/Detail, Executive Summary/Health/Capacity, Scenario, Health Score, People Analytics Correlation, Span of Control, Succession Readiness
- **Stats update:** 457 to 525 endpoints, 290 to 334 schemas, 232 to 290 paths, 22 to 23 tags
- **New Tag:** Tenant: Workforce Intelligence & Strategic Planning

### v8 (1 Agustus 2026)
- **Cross-documentation sync** — Training & Development Management module now fully documented in:
  - ✅ README.md (Module Completion Status, Phase 4, Testing table)
  - ✅ ARCHITECTURE_DESIGN_v1.6_Updated.md (Section 2.0, 3.1, 3.2, 4, Changelog v8)
  - ✅ docs/go-module-architecture-report.md (7 entities, 35 service methods, 31 tests)
  - ✅ Migration files 016_training.sql (MySQL + Postgres, 7 new tables)
- **Stats unchanged:** 457 endpoints, 290 schemas, 232 paths, 22 tags — final v7 numbers stable

### v7 (31 Juli 2026)
- **Training & Development Management** — Complete module with 7 resource groups (categories, courses, sessions, participants, materials, evaluations, certificates)
- **New Endpoints:** 35 new endpoints across 15 paths
- **New Schemas:** 28 new schemas for create/update/response/paginated across all 7 resource groups
- **Stats update:** 422 → 457 endpoints, 262 → 290 schemas, 217 → 232 paths, 21 → 22 tags

### v6 (31 Juli 2026)
- **Organization History, Versioning & Cloning** — Added 7 new endpoints (history audit trail, version snapshots, diff comparison, version restore, tree cloning)
- **New Schemas:** 9 new schemas: HistoryResponse, CreateVersionRequest, VersionResponse, PaginatedHistoryResponse, PaginatedVersionResponse, DiffEntry, DiffResponse, CloneRequest, CloneResponse
- **Stats update:** 415 → 422 endpoints, 253 → 262 schemas, 211 → 217 paths

### v5 (31 Juli 2026)
- **Performance Management & Recruitment ATS** — Full OpenAPI docs (67 combined endpoints)
- **Total:** 415 endpoints, 253 schemas, 211 paths, 21 tags| Area | Changes | Count |
|---|---|---|
| Performance Management module | New complete module — BSC-based KPI periods, perspectives, templates, indicators, evaluations, details, targets with full status workflow (DRAFT->COMPLETED) | 34 |
| Recruitment & Onboarding (ATS) module | Complete module — job requisitions, candidates, applications, interviews, onboarding | 33 |
| OpenAPI tag fixes | Fixed missing tag registrations and untagged endpoints | 3 |
| **Total improvements** | | **70** |

---
*Generated from OpenAPI spec v1.6.5*