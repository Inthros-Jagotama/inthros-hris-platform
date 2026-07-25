====================================================================================================
  HRIS PLATFORM — GO MODULE ARCHITECTURE REPORT
  Generated: 31 Juli 2026
====================================================================================================

## SECTION 1: TENANT MODULES (internal/modules/)

| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |
|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|
| approval | 5 | 16 | 30 | 15 | 15 | 64 |
| attendance | 10 | 30 | 34 | 30 | 30 | 88 |
| careerintelligence | 4 | 19 | 24 | 19 | 19 | 65 |
| competency | 7 | 35 | 36 | 35 | 35 | 60 |
| employee | 9 | 29 | 39 | 29 | 29 | 35 |
| employeemovement | 2 | 15 | 17 | 15 | 15 | 62 |
| jobmanagement | 18 | 88 | 90 | 88 | 88 | 67 |
| leave | 6 | 23 | 27 | 23 | 23 | 39 |
| organization | 3 | 13 | 16 | 12 | 12 | 0 |
| training | 7 | 35 | 38 | 35 | 35 | 31 |
| payroll | 21 | 48 | 70 | 47 | 47 | 40 |
| reimbursement | 3 | 15 | 17 | 15 | 15 | 60 |
| recruitment | 7 | 33 | 37 | 33 | 33 | 66 |
| performance | 7 | 35 | 36 | 34 | 34 | 55 |
| workforceintelligence | 7 | 46 | 36 | 70 | 68 | 108 |
| **TOTAL** | **116** | **480** | **547** | **500** | **498** | **840** |

### Test Breakdown per Module

| Module | Repo Tests | Service Tests | Handler Tests | Total |
|--------|:----------:|:-------------:|:-------------:|:-----:|
| approval | 25 | 25 | 14 | 64 |
| attendance | 33 | 29 | 26 | 88 |
| careerintelligence | 23 | 20 | 22 | 65 |
| competency | 12 | 35 | 13 | 60 |
| employee | 18 | 17 | 0 | 35 |
| employeemovement | 20 | 26 | 16 | 62 |
| jobmanagement | 26 | 28 | 10 | 64 |
| leave | 15 | 13 | 11 | 39 |
| organization | 0 | 0 | 0 | 0 |
| workforceintelligence | 31 | 41 | 36 | 108 |
| training | 14 | 9 | 8 | 31 |
| payroll | 17 | 23 | 0 | 40 |
| reimbursement | 19 | 24 | 17 | 60 |
| recruitment | 15 | 17 | 34 | 66 |
| performance | 14 | 24 | 17 | 55 |

## SECTION 2: PLATFORM MODULES (internal/platform/)

| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |
|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|
| company | 1 | 11 | 6 | 10 | 10 | 11 |
| license | 1 | 4 | 6 | 4 | 4 | 0 |
| modulemgmt | 2 | 7 | 7 | 7 | 7 | 0 |
| monitoring | 0 | 0 | 0 | 4 | 4 | 0 |
| user | 1 | 8 | 6 | 6 | 6 | 0 |
| **TOTAL** | **5** | **30** | **25** | **31** | **31** | **11** |

## SECTION 3: SHARED KERNEL PACKAGES (internal/pkg/)

| Package | Go Files | Test Funcs | Description |
|---------|:--------:|:----------:|-------------|
| auth | 1 | 0 | JWT authentication |
| authz | 12 | 91 | Casbin RBAC authorization |
| cache | 6 | 51 | Distributed cache + Pub/Sub |
| config | 1 | 0 | Viper configuration loader |
| crypto | 2 | 8 | AES-256-GCM encryption |
| database | 1 | 0 | Multi-tenant DB connection manager |
| docs | 1 | 0 | OpenAPI/Scalar documentation |
| driver | 1 | 0 | DB driver detection |
| errors | 0 | 0 | errors |
| logger | 1 | 0 | Structured logging |
| middleware | 5 | 0 | HTTP middleware (auth, cors, logger, recovery, tenant) |
| migrator | 3 | 3 | Database migration engine |
| module | 1 | 0 | Module SDK |
| router | 1 | 0 | Router setup |
| telemetry | 0 | 0 | telemetry |
| tenant | 0 | 0 | tenant |
| validator | 0 | 0 | validator |
| **TOTAL** | **36** | **153** | |

## SECTION 4: ENTITY DETAIL PER MODULE

### approval
- ApprovalFlow
- ApprovalFlowStep
- ApprovalInstance
- ApprovalAction
- ApprovalTask

### attendance
- AttendanceCompanySetting
- AttendanceCompanyShift
- AttendanceEmployeeShift
- AttendanceLocation
- AttendanceDeviceCapture
- AttendanceFaceCapture
- AttendanceEvent
- AttendanceSession
- AttendanceOvertimeRequest
- AttendanceExemptPosition

### competency
- Competency
- CompetenceValue
- CompetencyValue
- CompetencyEvent
- CompetencyEventTarget
- CompetencyScore
- CompetencyScoreDetail

### employee
- Employee
- EmployeeAddress
- EmergencyContact
- EmployeeFamily
- EmployeeEducation
- EmployeeExperience
- EmployeeDocument
- EmployeeInsurance
- Employment

### employeemovement
- EmployeeMovement
- EmployeeContract

### jobmanagement
- JobTitle
- JobTitleSub
- JobValue
- JobObjective
- JobIdentification
- JobResponsibility
- JobEducationExperience
- JobHRAuthority
- JobOperationalAuthority
- JobWorkingActivity
- JobWorkingRisk
- JobRelationship
- JobSubordinateControl
- JobAsset
- JobFinancial
- JobPotencyCompetency
- JobScore
- JobCompetencyGroup

### leave
- LeaveType
- LeaveAccrualPolicy
- LeaveReason
- LeaveRequest
- LeaveRequestDetail
- EmployeeLeaveBalance

### organization
- Organization
- OrganizationHistory (audit log)
- OrganizationVersion (snapshot)

### payroll
- SalaryComponent
- SalaryGradeComponent
- SalaryEmployeeComponent
- SalaryChangeLog
- SalaryEmployeeAdjustment
- PayrollPeriod
- EmployeePayrollProfile
- EmployeeBankProfile
- EmployeeBpjsProfile
- EmployeeTaxProfile
- BpjsSetting
- BpjsRateComponent
- Pph21Setting
- Pph21PtkpRate
- Pph21TaxBracket
- PayrollRun
- PayrollRunEmployee
- PayrollRunItem
- PayrollPayslip
- Pph21CalculationLog
- PayrollProfileChangeLog

### reimbursement
- ReimbursementType
- ReimbursementRequest
- ReimbursementItem

### recruitment
- JobRequisition
- Candidate
- JobApplication
- Interview
- OnboardingTaskTemplate
- EmployeeOnboarding
- OnboardingTaskItem

### performance
- PerformancePeriod
- PerformancePerspective
- PerformanceTemplate
- PerformanceIndicator
- PerformanceEvaluation
- PerformanceEvaluationDetail
- PerformanceTarget

### training
- TrainingCategory
- TrainingCourse
- TrainingSession
- TrainingParticipant
- TrainingMaterial
- TrainingEvaluation
- TrainingCertificate

### workforceintelligence
- WorkforcePlanningHeadcount
- WorkforceForecast
- WorkforceKPI
- WorkforceAnalyticsCache
- WorkforceScenario
- WorkforceRiskIndicator
- WorkforceHealthScore

### careerintelligence
- CareerTalentMap
- CareerInterest
- CareerPath
- CareerSuccessionPlan

### company
- Company

### license
- License

### modulemgmt
- PlatformModule
- CompanyModule

### user
- PlatformUser

## SECTION 5: SERVICE METHOD DETAIL

### approval.Service — 16 methods
- `CreateFlow()`
- `GetFlowByID()`
- `ListFlows()`
- `UpdateFlow()`
- `DeleteFlow()`
- `CreateStep()`
- `ListStepsByFlow()`
- `UpdateStep()`
- `DeleteStep()`
- `CreateInstance()`
- `GetInstanceByID()`
- `ListInstances()`
- `CancelInstance()`
- `SubmitAction()`
- `ListMyPendingTasks()`
- `resolveStepAssignees()`

### attendance.Service — 30 methods
- `UpsertCompanySetting()`
- `GetCompanySetting()`
- `CreateShift()`
- `GetShiftByID()`
- `ListShifts()`
- `UpdateShift()`
- `DeleteShift()`
- `CreateEmployeeShift()`
- `GetEmployeeShiftByID()`
- `ListEmployeeShifts()`
- `UpdateEmployeeShift()`
- `DeleteEmployeeShift()`
- `CreateLocation()`
- `GetLocationByID()`
- `ListLocations()`
- `UpdateLocation()`
- `DeleteLocation()`
- `CreateEvent()`
- `GetEventByID()`
- `ListEvents()`
- `GetSession()`
- `ListSessions()`
- `CreateOvertimeRequest()`
- `GetOvertimeRequestByID()`
- `ListOvertimeRequests()`
- `CreateExemptPosition()`
- `GetExemptPositionByID()`
- `ListExemptPositions()`
- `UpdateExemptPosition()`
- `DeleteExemptPosition()`

### competency.Service — 35 methods
- `CreateCompetency()`
- `GetCompetencyByID()`
- `ListCompetencies()`
- `UpdateCompetency()`
- `DeleteCompetency()`
- `CreateCompetenceValue()`
- `GetCompetenceValueByID()`
- `ListCompetenceValues()`
- `UpdateCompetenceValue()`
- `DeleteCompetenceValue()`
- `CreateCompetencyValue()`
- `GetCompetencyValueByID()`
- `ListCompetencyValues()`
- `UpdateCompetencyValue()`
- `DeleteCompetencyValue()`
- `CreateCompetencyEvent()`
- `GetCompetencyEventByID()`
- `ListCompetencyEvents()`
- `UpdateCompetencyEvent()`
- `DeleteCompetencyEvent()`
- `CreateCompetencyEventTarget()`
- `GetCompetencyEventTargetByID()`
- `ListCompetencyEventTargets()`
- `UpdateCompetencyEventTarget()`
- `DeleteCompetencyEventTarget()`
- `CreateCompetencyScore()`
- `GetCompetencyScoreByID()`
- `ListCompetencyScores()`
- `UpdateCompetencyScore()`
- `DeleteCompetencyScore()`
- `CreateCompetencyScoreDetail()`
- `GetCompetencyScoreDetailByID()`
- `ListCompetencyScoreDetails()`
- `UpdateCompetencyScoreDetail()`
- `DeleteCompetencyScoreDetail()`

### employee.Service — 29 methods
- `Create()`
- `GetByID()`
- `List()`
- `Update()`
- `Delete()`
- `CreateAddress()`
- `UpdateAddress()`
- `DeleteAddress()`
- `CreateEmergencyContact()`
- `UpdateEmergencyContact()`
- `DeleteEmergencyContact()`
- `CreateFamily()`
- `UpdateFamily()`
- `DeleteFamily()`
- `CreateEducation()`
- `UpdateEducation()`
- `DeleteEducation()`
- `CreateExperience()`
- `UpdateExperience()`
- `DeleteExperience()`
- `CreateDocument()`
- `UpdateDocument()`
- `DeleteDocument()`
- `CreateInsurance()`
- `UpdateInsurance()`
- `DeleteInsurance()`
- `CreateEmployment()`
- `UpdateEmployment()`
- `DeleteEmployment()`

### employeemovement.Service — 15 methods
- `CreateMovement()`
- `GetMovementByID()`
- `ListMovementsByEmployee()`
- `ListMovements()`
- `UpdateMovement()`
- `DeleteMovement()`
- `ApproveMovement()`
- `ExecuteMovement()`
- `CancelMovement()`
- `CreateContract()`
- `GetContractByID()`
- `ListContractsByEmployee()`
- `ListContracts()`
- `UpdateContract()`
- `DeleteContract()`

### jobmanagement.Service — 88 methods
- `CreateJobTitle()`
- `GetJobTitleByID()`
- `ListJobTitles()`
- `UpdateJobTitle()`
- `DeleteJobTitle()`
- `CreateJobTitleSub()`
- `GetJobTitleSubByID()`
- `ListJobTitleSubs()`
- `UpdateJobTitleSub()`
- `DeleteJobTitleSub()`
- `CreateJobValue()`
- `GetJobValueByID()`
- `ListJobValues()`
- `UpdateJobValue()`
- `DeleteJobValue()`
- `CreateJobObjective()`
- `GetJobObjectiveByID()`
- `ListJobObjectives()`
- `UpdateJobObjective()`
- `DeleteJobObjective()`
- `CreateJobIdentification()`
- `GetJobIdentificationByID()`
- `ListJobIdentifications()`
- `UpdateJobIdentification()`
- `DeleteJobIdentification()`
- `CreateJobResponsibility()`
- `GetJobResponsibilityByID()`
- `ListJobResponsibilities()`
- `UpdateJobResponsibility()`
- `DeleteJobResponsibility()`
- `CreateJobEducationExperience()`
- `GetJobEducationExperienceByID()`
- `ListJobEducationExperiences()`
- `UpdateJobEducationExperience()`
- `DeleteJobEducationExperience()`
- `CreateJobHRAuthority()`
- `GetJobHRAuthorityByID()`
- `ListJobHRAuthorities()`
- `UpdateJobHRAuthority()`
- `DeleteJobHRAuthority()`
- `CreateJobOperationalAuthority()`
- `GetJobOperationalAuthorityByID()`
- `ListJobOperationalAuthorities()`
- `UpdateJobOperationalAuthority()`
- `DeleteJobOperationalAuthority()`
- `CreateJobWorkingActivity()`
- `GetJobWorkingActivityByID()`
- `ListJobWorkingActivities()`
- `UpdateJobWorkingActivity()`
- `DeleteJobWorkingActivity()`
- `CreateJobWorkingRisk()`
- `GetJobWorkingRiskByID()`
- `ListJobWorkingRisks()`
- `UpdateJobWorkingRisk()`
- `DeleteJobWorkingRisk()`
- `CreateJobRelationship()`
- `GetJobRelationshipByID()`
- `ListJobRelationships()`
- `UpdateJobRelationship()`
- `DeleteJobRelationship()`
- `CreateJobSubordinateControl()`
- `GetJobSubordinateControlByID()`
- `ListJobSubordinateControls()`
- `UpdateJobSubordinateControl()`
- `DeleteJobSubordinateControl()`
- `CreateJobAsset()`
- `GetJobAssetByID()`
- `ListJobAssets()`
- `UpdateJobAsset()`
- `DeleteJobAsset()`
- `CreateJobFinancial()`
- `GetJobFinancialByID()`
- `ListJobFinancials()`
- `UpdateJobFinancial()`
- `DeleteJobFinancial()`
- `CreateJobPotencyCompetency()`
- `GetJobPotencyCompetencyByID()`
- `ListJobPotencyCompetencies()`
- `UpdateJobPotencyCompetency()`
- `DeleteJobPotencyCompetency()`
- `UpsertJobScore()`
- `GetJobScoreByOrganization()`
- `ListJobScores()`
- `CreateJobCompetencyGroup()`
- `GetJobCompetencyGroupByID()`
- `ListJobCompetencyGroups()`
- `UpdateJobCompetencyGroup()`
- `DeleteJobCompetencyGroup()`

### leave.Service — 23 methods
- `CreateLeaveType()`
- `GetLeaveTypeByID()`
- `ListLeaveTypes()`
- `UpdateLeaveType()`
- `DeleteLeaveType()`
- `CreateAccrualPolicy()`
- `GetAccrualPolicyByID()`
- `ListAccrualPolicies()`
- `UpdateAccrualPolicy()`
- `DeleteAccrualPolicy()`
- `CreateLeaveReason()`
- `GetLeaveReasonByID()`
- `ListLeaveReasons()`
- `UpdateLeaveReason()`
- `DeleteLeaveReason()`
- `CreateLeaveRequest()`
- `GetLeaveRequestByID()`
- `ListLeaveRequests()`
- `UpdateLeaveRequestStatus()`
- `DeleteLeaveRequest()`
- `ListLeaveRequestDetails()`
- `GetLeaveBalance()`
- `ListLeaveBalances()`

### training.Service — 35 methods
- `CreateCategory()`
- `GetCategoryByID()`
- `ListCategories()`
- `UpdateCategory()`
- `DeleteCategory()`
- `CreateCourse()`
- `GetCourseByID()`
- `ListCourses()`
- `UpdateCourse()`
- `DeleteCourse()`
- `CreateSession()`
- `GetSessionByID()`
- `ListSessions()`
- `UpdateSession()`
- `UpdateSessionStatus()`
- `DeleteSession()`
- `CreateParticipant()`
- `GetParticipantByID()`
- `ListParticipants()`
- `UpdateParticipant()`
- `DeleteParticipant()`
- `CreateMaterial()`
- `ListMaterials()`
- `UpdateMaterial()`
- `DeleteMaterial()`
- `CreateEvaluation()`
- `GetEvaluationByID()`
- `ListEvaluations()`
- `UpdateEvaluation()`
- `DeleteEvaluation()`
- `CreateCertificate()`
- `GetCertificateByID()`
- `ListCertificates()`
- `UpdateCertificate()`
- `DeleteCertificate()`

### organization.Service — 13 methods
- `Create()`
- `GetByID()`
- `List()`
- `GetTree()`
- `Update()`
- `Delete()`
- `GetHistory()`
- `CreateVersion()`
- `ListVersions()`
- `GetVersion()`
- `DiffVersions()`
- `RestoreVersion()`
- `CloneVersion()`

### payroll.Service — 48 methods
- `SetApprovalEngine()`
- `CreateSalaryComponent()`
- `GetSalaryComponentByID()`
- `ListSalaryComponents()`
- `UpdateSalaryComponent()`
- `DeleteSalaryComponent()`
- `CreatePayrollPeriod()`
- `ListPayrollPeriods()`
- `UpdatePayrollPeriod()`
- `CreateEmployeePayrollProfile()`
- `GetEmployeePayrollProfileByID()`
- `ListEmployeePayrollProfiles()`
- `DeleteEmployeePayrollProfile()`
- `GetEmployeeBankProfileByID()`
- `UpdateEmployeeBankProfile()`
- `DeleteEmployeeBankProfile()`
- `CreateEmployeeBankProfile()`
- `GetEmployeeBpjsProfileByID()`
- `UpdateEmployeeBpjsProfile()`
- `DeleteEmployeeBpjsProfile()`
- `CreateEmployeeBpjsProfile()`
- `GetEmployeeTaxProfileByID()`
- `UpdateEmployeeTaxProfile()`
- `DeleteEmployeeTaxProfile()`
- `CreateEmployeeTaxProfile()`
- `CreateBpjsSetting()`
- `CreatePph21Setting()`
- `CreatePayrollRun()`
- `ListPayrollRuns()`
- `GetPayrollRunByID()`
- `UpdatePayrollRunStatus()`
- `CheckPayrollRunApproval()`
- `UpdateBpjsSetting()`
- `GetBpjsSettingByID()`
- `ListBpjsSettings()`
- `DeleteBpjsSetting()`
- `GetBpjsRateComponentByID()`
- `UpdateBpjsRateComponent()`
- `DeleteBpjsRateComponent()`
- `CreateBpjsRateComponent()`
- `GetPph21SettingByID()`
- `ListPph21Settings()`
- `UpdatePph21Setting()`
- `DeletePph21Setting()`
- `CreatePph21PtkpRate()`
- `ListPph21PtkpRates()`
- `CreatePph21TaxBracket()`
- `ListPph21TaxBrackets()`

### reimbursement.Service — 15 methods
- `CreateReimbursementType()`
- `GetReimbursementTypeByID()`
- `ListReimbursementTypes()`
- `UpdateReimbursementType()`
- `DeleteReimbursementType()`
- `CreateReimbursementRequest()`
- `GetReimbursementRequestByID()`
- `ListReimbursementRequests()`
- `UpdateReimbursementRequest()`
- `UpdateReimbursementRequestStatus()`
- `DeleteReimbursementRequest()`
- `CreateReimbursementItem()`
- `ListReimbursementItems()`
- `UpdateReimbursementItem()`
- `DeleteReimbursementItem()`

### recruitment.Service — 33 methods
- `CreateRequisition()`
- `GetRequisitionByID()`
- `ListRequisitions()`
- `UpdateRequisition()`
- `DeleteRequisition()`
- `CreateCandidate()`
- `GetCandidateByID()`
- `ListCandidates()`
- `UpdateCandidate()`
- `DeleteCandidate()`
- `CreateApplication()`
- `GetApplicationByID()`
- `ListApplications()`
- `UpdateApplicationStatus()`
- `DeleteApplication()`
- `CreateInterview()`
- `GetInterviewByID()`
- `ListInterviews()`
- `UpdateInterview()`
- `DeleteInterview()`
- `CreateOnboardingTaskTemplate()`
- `ListOnboardingTaskTemplates()`
- `UpdateOnboardingTaskTemplate()`
- `DeleteOnboardingTaskTemplate()`
- `CreateEmployeeOnboarding()`
- `GetEmployeeOnboardingByID()`
- `ListEmployeeOnboardings()`
- `UpdateEmployeeOnboarding()`
- `DeleteEmployeeOnboarding()`
- `CreateOnboardingTaskItem()`
- `ListOnboardingTaskItems()`
- `UpdateOnboardingTaskItem()`
- `DeleteOnboardingTaskItem()`
- `calcTotalPages()`

### performance.Service — 35 methods
- `CreatePerformancePeriod()`
- `GetPerformancePeriodByID()`
- `ListPerformancePeriods()`
- `UpdatePerformancePeriod()`
- `DeletePerformancePeriod()`
- `CreatePerformancePerspective()`
- `GetPerformancePerspectiveByID()`
- `ListPerformancePerspectives()`
- `UpdatePerformancePerspective()`
- `DeletePerformancePerspective()`
- `CreatePerformanceTemplate()`
- `GetPerformanceTemplateByID()`
- `ListPerformanceTemplates()`
- `UpdatePerformanceTemplate()`
- `DeletePerformanceTemplate()`
- `CreatePerformanceIndicator()`
- `GetPerformanceIndicatorByID()`
- `ListPerformanceIndicators()`
- `UpdatePerformanceIndicator()`
- `DeletePerformanceIndicator()`
- `CreatePerformanceEvaluation()`
- `GetPerformanceEvaluationByID()`
- `ListPerformanceEvaluations()`
- `UpdatePerformanceEvaluation()`
- `UpdateEvaluationStatus()`
- `DeletePerformanceEvaluation()`
- `CreateEvaluationDetail()`
- `ListEvaluationDetails()`
- `UpdateEvaluationDetail()`
- `DeleteEvaluationDetail()`
- `CreatePerformanceTarget()`
- `ListPerformanceTargets()`
- `UpdatePerformanceTarget()`
- `DeletePerformanceTarget()`
- `calcTotalPages()`

### workforceintelligence.Service — 46 methods
- `CreateHeadcountPlan()`
- `GetHeadcountPlanByID()`
- `ListHeadcountPlans()`
- `UpdateHeadcountPlan()`
- `DeleteHeadcountPlan()`
- `CreateForecast()`
- `GetForecastByID()`
- `ListForecasts()`
- `UpdateForecast()`
- `DeleteForecast()`
- `GetGapAnalysis()`
- `GetProjections()`
- `CreateKPI()`
- `ListKPIs()`
- `GetKPISummary()`
- `GetKPIByCode()`
- `GetHeadcountAnalytics()`
- `GetMovementAnalytics()`
- `GetRiskDashboard()`
- `GetRiskDetail()`
- `ListRiskIndicators()`
- `GetRiskIndicatorByID()`
- `UpdateRiskIndicator()`
- `GetExecutiveSummary()`
- `GetExecutiveGrowth()`
- `GetExecutiveCostTrend()`
- `GetExecutiveAttritionTrend()`
- `GetExecutiveCapacity()`
- `GetExecutiveRiskOverview()`
- `GetExecutiveHealthScore()`
- `ListScenarios()`
- `CreateScenario()`
- `GetScenarioByID()`
- `UpdateScenario()`
- `DeleteScenario()`
- `RunScenario()`
- `CloneScenario()`
- `GetHealthDashboard()`
- `ListHealthScores()`
- `GetHealthScoreByID()`
- `GetSpanOfControl()`
- `GetSuccessionReadiness()`
- `GetPeopleAnalytics()`
- `GetCapacityForecast()`
- `GetPayrollCostBreakdown()`
- `GetCostPerEmployee()`
- `executeSimulation()` (private)

### careerintelligence.Service — 19 methods
- `CreateTalentMap()`
- `GetTalentMapByID()`
- `ListTalentMaps()`
- `UpdateTalentMap()`
- `DeleteTalentMap()`
- `GetTalentGrid()`
- `GetEmployeeTalentProfile()`
- `CreateCareerInterest()`
- `ListCareerInterests()`
- `GetEmployeeCareerInterests()`
- `CreateCareerPath()`
- `ListCareerPaths()`
- `DeleteCareerPath()`
- `GetGapAnalysis()`
- `CreateSuccessionPlan()`
- `ListSuccessionPlans()`
- `GetSuccessionPlanByID()`
- `UpdateSuccessionPlan()`
- `DeleteSuccessionPlan()`

## SECTION 6: GRAND TOTALS

| Category | Count |
|----------|:-----:|
| Tenant Modules | 14 |
| Platform Modules | 5 |
| Shared Kernel Packages | 17 |
| **Total Architecture Layers** | **36** |
| Total GORM Entities (tenant) | 109 |
| Total GORM Entities (platform) | 5 |
| **Total Entities (combined)** | **121** |
| Total Service Methods | 510 |
| Total Repository Methods | 572 |
| Total Handler Functions | 531 |
| Total Route Registrations | 529 |
| **Total Unit Tests (all)** | **1004** |
| Total Go Source Files | 155 |
| Total Test Files (_test.go) | 69 |
| **Total Go Files** | **224** |

## SECTION 7: TEST FILE INVENTORY

| File | Test Funcs |
|------|:----------:|
| `modules\approval\handler_test.go` | 14 |
| `modules\approval\helpers_test.go` | 0 |
| `modules\approval\repository_test.go` | 25 |
| `modules\approval\service_test.go` | 25 |
| `modules\attendance\handler_test.go` | 26 |
| `modules\attendance\helpers_test.go` | 0 |
| `modules\attendance\repository_test.go` | 33 |
| `modules\attendance\service_test.go` | 29 |
| `modules\competency\handler_test.go` | 13 |
| `modules\competency\helpers_test.go` | 0 |
| `modules\competency\repository_test.go` | 12 |
| `modules\competency\service_test.go` | 35 |
| `modules\employee\helpers_test.go` | 0 |
| `modules\employee\repository_test.go` | 18 |
| `modules\employee\service_test.go` | 17 |
| `modules\employeemovement\handler_test.go` | 16 |
| `modules\employeemovement\helpers_test.go` | 0 |
| `modules\employeemovement\repository_test.go` | 20 |
| `modules\employeemovement\service_test.go` | 26 |
| `modules\recruitment\handler_test.go` | 34 |
| `modules\recruitment\helpers_test.go` | 0 |
| `modules\recruitment\repository_test.go` | 15 |
| `modules\recruitment\service_test.go` | 17 |
| `modules\performance\handler_test.go` | 17 |
| `modules\performance\helpers_test.go` | 0 |
| `modules\performance\repository_test.go` | 14 |
| `modules\performance\service_test.go` | 24 |
| `modules\jobmanagement\handler_test.go` | 10 |
| `modules\jobmanagement\helpers_test.go` | 0 |
| `modules\jobmanagement\repository_test.go` | 26 |
| `modules\jobmanagement\router_test.go` | 3 |
| `modules\jobmanagement\service_test.go` | 28 |
| `modules\leave\handler_test.go` | 11 |
| `modules\leave\helpers_test.go` | 0 |
| `modules\leave\repository_test.go` | 15 |
| `modules\leave\service_test.go` | 13 |
| `modules\payroll\helpers_test.go` | 0 |
| `modules\payroll\repository_test.go` | 17 |
| `modules\payroll\service_test.go` | 23 |
| `modules\reimbursement\handler_test.go` | 17 |
| `modules\reimbursement\helpers_test.go` | 0 |
| `modules\reimbursement\repository_test.go` | 19 |
| `modules\reimbursement\service_test.go` | 24 |
| `modules\workforceintelligence\helpers_test.go` | 0 |
| `modules\workforceintelligence\repository_test.go` | 31 |
| `modules\workforceintelligence\service_test.go` | 41 |
| `modules\workforceintelligence\handler_test.go` | 36 |
| `modules\careerintelligence\helpers_test.go` | 0 |
| `modules\careerintelligence\repository_test.go` | 23 |
| `modules\careerintelligence\service_test.go` | 20 |
| `modules\careerintelligence\handler_test.go` | 22 |
| `modules\training\handler_test.go` | 8 |
| `modules\training\helpers_test.go` | 0 |
| `modules\training\repository_test.go` | 14 |
| `modules\training\service_test.go` | 9 |
| `pkg\authz\handler_test.go` | 16 |
| `pkg\authz\helpers_test.go` | 0 |
| `pkg\authz\rbac_test.go` | 26 |
| `pkg\authz\repository_test.go` | 24 |
| `pkg\authz\service_test.go` | 25 |
| `pkg\cache\bench_test.go` | 0 |
| `pkg\cache\cache_integration_test.go` | 8 |
| `pkg\cache\cache_test.go` | 24 |
| `pkg\cache\pubsub_test.go` | 19 |
| `pkg\crypto\crypto_test.go` | 8 |
| `pkg\migrator\migrator_integration_test.go` | 3 |
| `platform\company\handler_test.go` | 4 |
| `platform\company\helpers_test.go` | 0 |
| `platform\company\service_test.go` | 7 |

====================================================================================================
  END OF REPORT
====================================================================================================