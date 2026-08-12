====================================================================================================
  HRIS PLATFORM — GO MODULE ARCHITECTURE REPORT
  Generated: 12 Aug 2026

  Index dokumentasi: docs/README.md  |  Terkait: platform-architecture-design.md, openapi-report.md
====================================================================================================

## SECTION 1: TENANT MODULES (internal/modules/)

| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |
|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|
| approval | 6 | 28 | 43 | 17 | 17 | 105 |
| attendance | 11 | 54 | 56 | 40 | 40 | 146 |
| careerintelligence | 5 | 23 | 28 | 21 | 21 | 66 |
| competency | 7 | 35 | 37 | 35 | 35 | 60 |
| employee | 10 | 38 | 51 | 36 | 36 | 40 |
| employeemovement | 6 | 46 | 50 | 25 | 25 | 150 |
| jobmanagement | 23 | 100 | 98 | 96 | 96 | 83 |
| leave | 7 | 38 | 33 | 25 | 25 | 79 |
| notification | 1 | 5 | 5 | 5 | 4 | 8 |
| organization | 3 | 20 | 26 | 18 | 18 | 9 |
| payroll | 21 | 49 | 70 | 47 | 47 | 43 |
| performance | 17 | 109 | 126 | 99 | 147 | 125 |
| rbac | 5 | 8 | 15 | 8 | 8 | 0 |
| recruitment | 7 | 33 | 38 | 33 | 33 | 75 |
| reimbursement | 3 | 17 | 17 | 15 | 15 | 67 |
| setting | 21 | 119 | 126 | 107 | 107 | 52 |
| training | 28 | 129 | 139 | 123 | 123 | 71 |
| useraccount | 3 | 10 | 20 | 7 | 4 | 0 |
| workforceintelligence | 7 | 47 | 38 | 71 | 69 | 112 |
| **TOTAL** | **191** | **908** | **1016** | **828** | **870** | **1291** |

### Test Breakdown per Module

| Module | Repo Tests | Service Tests | Handler Tests | Other | Total |
|--------|:----------:|:-------------:|:-------------:|:-----:|:-----:|
| approval | 25 | 32 | 22 | 26 | 105 |
| attendance | 34 | 39 | 26 | 47 | 146 |
| careerintelligence | 24 | 20 | 22 | 0 | 66 |
| competency | 12 | 35 | 13 | 0 | 60 |
| employee | 18 | 22 | 0 | 0 | 40 |
| employeemovement | 27 | 49 | 16 | 58 | 150 |
| jobmanagement | 27 | 29 | 10 | 17 | 83 |
| leave | 17 | 18 | 11 | 33 | 79 |
| notification | 2 | 6 | 0 | 0 | 8 |
| organization | 0 | 9 | 0 | 0 | 9 |
| payroll | 17 | 26 | 0 | 0 | 43 |
| performance | 19 | 27 | 14 | 65 | 125 |
| recruitment | 27 | 20 | 28 | 0 | 75 |
| reimbursement | 19 | 24 | 17 | 7 | 67 |
| setting | 18 | 26 | 0 | 8 | 52 |
| training | 14 | 17 | 9 | 31 | 71 |
| workforceintelligence | 33 | 43 | 36 | 0 | 112 |

## SECTION 2: PLATFORM MODULES (internal/platform/)

| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |
|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|
| company | 2 | 15 | 11 | 14 | 11 | 30 |
| license | 1 | 5 | 10 | 5 | 5 | 0 |
| modulemgmt | 2 | 9 | 7 | 7 | 7 | 3 |
| monitoring | 0 | 0 | 0 | 5 | 5 | 0 |
| package | 2 | 12 | 9 | 9 | 9 | 27 |
| user | 1 | 12 | 8 | 8 | 8 | 0 |
| **TOTAL** | **8** | **53** | **45** | **48** | **45** | **60** |

## SECTION 3: SHARED KERNEL PACKAGES (internal/pkg/)

| Package | Go Files | Test Funcs | Description |
|---------|:--------:|:----------:|-------------|
| auth | 1 | 0 | JWT authentication |
| authctx | 1 | 0 | Auth context helpers (user/company from gin context) |
| authz | 13 | 104 | Casbin RBAC authorization |
| cache | 6 | 51 | Distributed cache + Pub/Sub |
| config | 1 | 0 | Viper configuration loader |
| crypto | 2 | 8 | AES-256-GCM encryption |
| database | 1 | 0 | Multi-tenant DB connection manager |
| docs | 1 | 0 | OpenAPI/Scalar documentation |
| driver | 1 | 0 | DB driver detection |
| errors | 0 | 0 | Shared error helpers |
| httputil | 4 | 0 | Bilingual response helpers (SuccessJSON, CreatedJSON, ErrorJSON, NotFound) + locale message catalog (80+ EN/ID pairs) + custom Indonesian validators (NIK, NPWP, KK, Passport, SIM, No Rekening) |
| logger | 1 | 0 | Structured logging |
| mailer | 1 | 0 | Email sending (SMTP/template) |
| middleware | 10 | 16 | HTTP middleware: Auth, CORS, Logger, Recovery, Tenant, Localize (auto-detect Accept-Language) |
| migrator | 3 | 3 | Database migration engine |
| module | 1 | 0 | Module SDK |
| onpremise | 2 | 6 | On-premise license enforcement |
| router | 1 | 0 | Router setup |
| telemetry | 0 | 0 | Telemetry/metrics |
| tenant | 0 | 0 | Tenant resolution helpers |
| tenantseed | 7 | 8 | Tenant seed data (nationalities, competencies, RBAC) |
| upload | 1 | 0 | menyediakan endpoint upload file generik untuk lampiran |
| validator | 0 | 0 | Validator helpers |
| **TOTAL** | **58** | **196** | |

## SECTION 4: ENTITY DETAIL PER MODULE

### approval
- ApprovalFlow
- ApprovalFlowStep
- ApprovalFlowStepOrganization
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
- AttendanceCorrectionRequest

### careerintelligence
- CareerTalentMap
- CareerInterest
- CareerPath
- CareerPathStep
- CareerSuccessionPlan

### company
- LicenseRef
- Company

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
- EmployeeBankAccount
- Employment

### employeemovement
- EmployeeMovement
- EmployeeMovementAudit
- EmployeeMovementDocument
- CareerPath
- CareerPathStep
- EmployeeContract

### jobmanagement
- JobTitle
- JobTitleSub
- JobValue
- JobManagementValueCluster
- JobObjective
- JobIdentification
- JobResponsibility
- JobEducationExperience
- JobManagementMajor
- JobManagementJobFamily
- JobHRAuthority
- JobOperationalAuthority
- JobWorkingActivity
- JobWorkingRisk
- JobRelationship
- JobManagementRelationshipDetail
- OrganizationRef
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
- LeaveBalanceTransaction

### license
- License

### modulemgmt
- PlatformModule
- CompanyModule

### monitoring

### notification
- Notification

### organization
- Organization
- OrganizationHistory
- OrganizationVersion

### package
- Package
- PackageModule

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

### performance
- PerformancePeriod
- PerformancePerspective
- PerformanceTemplate
- PerformanceIndicator
- PerformanceEvaluation
- PerformanceEvaluationDetail
- PerformanceEvaluationProgramItem
- PerformanceTarget
- PerformanceProgress
- PerformanceComment
- PerformanceAttachment
- PerformanceRating
- PerformanceIndicatorFormula
- PerformanceLog
- PerformanceComponent
- PerformanceOrganizationComponent
- PerformanceEvaluationComponent

### rbac
- Role
- Permission
- RoleHasPermission
- ModelHasRole
- TenantUser

### recruitment
- JobRequisition
- Candidate
- JobApplication
- Interview
- OnboardingTaskTemplate
- EmployeeOnboarding
- OnboardingTaskItem

### reimbursement
- ReimbursementType
- ReimbursementRequest
- ReimbursementItem

### setting
- Zone
- Province
- Regency
- District
- Village
- Education
- EducationMajor
- Religion
- MaritalStatus
- RelationshipType
- EmploymentStatus
- Bank
- Nationality
- JobFamily
- Grading
- Insurance
- CompanyHoliday
- SalaryGrade
- TER
- Competency
- PTKP

### training
- TrainingCategory
- TrainingCourse
- TrainingSession
- TrainingParticipant
- TrainingMaterial
- TrainingEvaluation
- TrainingCertificate
- TrainingProvider
- TrainingTrainer
- TrainingSessionTrainer
- TrainingAttendance
- TrainingPlan
- TrainingPlanItem
- TrainingNeed
- TrainingRequest
- TrainingCourseObjective
- TrainingCourseCompetency
- TrainingCoursePrerequisite
- TrainingMandatory
- TrainingSessionCost
- TrainingDocument
- TrainingAssessment
- TrainingAssessmentResult
- TrainingEvaluationForm
- TrainingEvaluationQuestion
- TrainingEvaluationAnswer
- TrainingEffectivenessAssessment
- TrainingCertification

### user
- PlatformUser

### useraccount
- EmployeeAccount
- TenantUser
- ModelHasRole

### workforceintelligence
- WorkforcePlanningHeadcount
- WorkforceForecast
- WorkforceKPI
- WorkforceAnalyticsCache
- WorkforceScenario
- WorkforceRiskIndicator
- WorkforceHealthScore

## SECTION 5: SERVICE METHOD DETAIL

### approval.Service — 28 methods
- `SetNotifier()`
- `RegisterStatusHandler()`
- `notifyStatusChange()`
- `notifyNewTasks()`
- `resolveNotifyRecipients()`
- `SetModuleChecker()`
- `ListAvailableModules()`
- `ensureModuleSubscribed()`
- `CreateFlow()`
- `GetFlowByID()`
- `GetActiveFlowByModule()`
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
- `resolveSupervisorAssignees()`
- `resolveOrganizationAssignees()`
- `advanceThroughWatcherSteps()`

### attendance.Service — 54 methods
- `SetApprovalEngine()`
- `SetNotifier()`
- `notifyRequestOutcome()`
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
- `checkEventSequence()`
- `applyEventValidation()`
- `GetEventByID()`
- `ListEvents()`
- `GetSession()`
- `ListSessions()`
- `GetEmployeeCalendar()`
- `GetAttendanceReport()`
- `GetEmployeeSummary()`
- `CreateOvertimeRequest()`
- `HandleApprovalStatusChange()`
- `handleOvertimeApprovalStatusChange()`
- `AssignOvertimeRequest()`
- `notifyOvertimeAssigned()`
- `SubmitActualOvertime()`
- `ListAssignableEmployees()`
- `CancelOvertimeRequest()`
- `handleActualApprovalStatusChange()`
- `applyOvertimeCalculation()`
- `GetOvertimeRequestByID()`
- `ListOvertimeRequests()`
- `CreateCorrectionRequest()`
- `handleCorrectionApprovalStatusChange()`
- `applyCorrectionToSession()`
- `GetCorrectionRequestByID()`
- `ListCorrectionRequests()`
- `CreateExemptPosition()`
- `GetExemptPositionByID()`
- `ListExemptPositions()`
- `UpdateExemptPosition()`
- `DeleteExemptPosition()`
- `recalculateSession()`
- `ApplyApprovedLeave()`

### careerintelligence.Service — 23 methods
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
- `CreateCareerPathLadder()`
- `GetCareerPathByID()`
- `UpdateCareerPath()`
- `buildCareerPathSteps()`
- `ListCareerPaths()`
- `DeleteCareerPath()`
- `GetGapAnalysis()`
- `CreateSuccessionPlan()`
- `ListSuccessionPlans()`
- `GetSuccessionPlanByID()`
- `UpdateSuccessionPlan()`
- `DeleteSuccessionPlan()`
- `careerPathToResponse()`

### company.Service — 15 methods
- `Create()`
- `provisionTenant()`
- `RotateCredentials()`
- `MigrateTenantDB()`
- `ResolveByHost()`
- `GetByID()`
- `countTenantEmployees()`
- `List()`
- `Update()`
- `UpdateCurrent()`
- `Delete()`
- `UpdateStatus()`
- `Suspend()`
- `Activate()`
- `Terminate()`

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

### employee.Service — 38 methods
- `SetQuotaChecker()`
- `checkQuota()`
- `Create()`
- `GetByID()`
- `List()`
- `Update()`
- `Delete()`
- `UpdatePhoto()`
- `DeletePhoto()`
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
- `CreateDocumentRecord()`
- `UpdateDocumentFile()`
- `CreateInsurance()`
- `UpdateInsurance()`
- `DeleteInsurance()`
- `CreateBank()`
- `UpdateBank()`
- `DeleteBank()`
- `CreateEmployment()`
- `UpdateEmployment()`
- `DeleteEmployment()`

### employeemovement.Service — 46 methods
- `SetApprovalEngine()`
- `SetCareerExecutor()`
- `SetNotifier()`
- `SetPerformanceProvider()`
- `SetCompetencyProvider()`
- `SetOKRProvider()`
- `recordAudit()`
- `notifyMovementOutcome()`
- `validateMovement()`
- `CreateMovement()`
- `fillMovementSnapshot()`
- `enrichMovementResponses()`
- `enrichContractResponses()`
- `GetMovementByID()`
- `ListMovementsByEmployee()`
- `ListMovements()`
- `UpdateMovement()`
- `DeleteMovement()`
- `SubmitMovement()`
- `HandleApprovalStatusChange()`
- `ExecuteMovement()`
- `CancelMovement()`
- `HandleCancellationStatusChange()`
- `ListMovementAudits()`
- `notifyContractEvent()`
- `ProcessContractExpiration()`
- `CreateMovementDocument()`
- `ListMovementDocuments()`
- `DeleteMovementDocument()`
- `GetCareerHistory()`
- `currentPosition()`
- `CreateContract()`
- `GetContractByID()`
- `ListContractsByEmployee()`
- `ListContracts()`
- `GetMovementReport()`
- `GetContractReport()`
- `GetHRDashboard()`
- `UpdateContract()`
- `DeleteContract()`
- `GetMovementEligibility()`
- `GetPromotionEligibility()`
- `buildEligibility()`
- `evaluateDefaultRules()`
- `evaluatePromotionRules()`
- `findPromotionNextStep()`

### jobmanagement.Service — 100 methods
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
- `ListJobValuesTree()`
- `ListJobValues()`
- `UpdateJobValue()`
- `DeleteJobValue()`
- `ListJobValueClusters()`
- `UpdateJobValueClusters()`
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
- `CreateJobRelationshipDetail()`
- `GetJobRelationshipDetailByID()`
- `ListJobRelationshipDetails()`
- `UpdateJobRelationshipDetail()`
- `DeleteJobRelationshipDetail()`
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
- `recalculateScore()`
- `UpsertJobScore()`
- `GetJobScoreByOrganization()`
- `RecalculateJobScore()`
- `RecalculateJobScores()`
- `scoreFromResult()`
- `ListJobScores()`
- `CreateJobCompetencyGroup()`
- `GetJobCompetencyGroupByID()`
- `ListJobCompetencyGroups()`
- `UpdateJobCompetencyGroup()`
- `DeleteJobCompetencyGroup()`

### leave.Service — 38 methods
- `getOrCreateLeaveBalance()`
- `applyLeaveUsage()`
- `reverseLeaveUsage()`
- `writeLeaveBalanceTransaction()`
- `SetApprovalEngine()`
- `SetHolidayProvider()`
- `SetNotifier()`
- `SetAttendanceSessionUpdater()`
- `holidayDatesInRange()`
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
- `applyBalanceEffectOnStatusChange()`
- `HandleApprovalStatusChange()`
- `notifyLeaveOutcome()`
- `applyAttendanceIntegration()`
- `DeleteLeaveRequest()`
- `ListLeaveRequestDetails()`
- `GetEmployeeCalendar()`
- `GetLeaveUsageReport()`
- `GetLeaveBalance()`
- `ListLeaveBalances()`

### license.Service — 5 methods
- `CreateLicense()`
- `GetLicense()`
- `ListLicenses()`
- `UpdateLicense()`
- `DeleteLicense()`

### modulemgmt.Service — 9 methods
- `SetCacheManager()`
- `invalidateLicenseCache()`
- `CreateModule()`
- `GetModule()`
- `ListModules()`
- `UpdateModule()`
- `ActivateModule()`
- `DeactivateModule()`
- `ListCompanyModules()`

### notification.Service — 5 methods
- `Notify()`
- `ListNotifications()`
- `MarkAsRead()`
- `MarkAllAsRead()`
- `GetUnreadCount()`

### organization.Service — 20 methods
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
- `captureHistory()`
- `CreateSummary()`
- `GetSummaryByID()`
- `ListSummaries()`
- `UpdateSummary()`
- `DeleteSummary()`
- `GetSummaryStats()`

### package.Service — 12 methods
- `CreatePackage()`
- `GetPackage()`
- `GetPackageBySlug()`
- `ListPackages()`
- `ListPublishedPackages()`
- `UpdatePackage()`
- `DeletePackage()`
- `PublishPackage()`
- `UnpublishPackage()`
- `ValidatePackageDependencies()`
- `updateStatus()`
- `validateModuleDependencies()`

### payroll.Service — 49 methods
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
- `HandleApprovalStatusChange()`
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

### performance.Service — 109 methods
- `SetApprovalEngine()`
- `SetNotifier()`
- `notifyEvaluationOutcome()`
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
- `notifyTemplateCreated()`
- `GetPerformanceTemplateByID()`
- `enrichTemplateResponses()`
- `ListPerformanceTemplates()`
- `GetMyKPIContext()`
- `ListTemplateOrganizationScope()`
- `UpdatePerformanceTemplate()`
- `DeletePerformanceTemplate()`
- `authorizeTemplateOrg()`
- `DuplicatePerformanceTemplate()`
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
- `CreateProgramItem()`
- `ListProgramItems()`
- `UpdateProgramItemTarget()`
- `UpdateProgramItemActual()`
- `DeleteProgramItem()`
- `CreatePerformanceTarget()`
- `ListPerformanceTargets()`
- `UpdatePerformanceTarget()`
- `DeletePerformanceTarget()`
- `CreatePerformanceProgress()`
- `GetPerformanceProgressByID()`
- `ListPerformanceProgressByDetailID()`
- `UpdatePerformanceProgress()`
- `DeletePerformanceProgress()`
- `CreatePerformanceComment()`
- `GetPerformanceCommentByID()`
- `ListPerformanceCommentsByEvaluationID()`
- `UpdatePerformanceComment()`
- `DeletePerformanceComment()`
- `CreatePerformanceAttachment()`
- `GetPerformanceAttachmentByID()`
- `ListPerformanceAttachmentsByDetailID()`
- `UpdatePerformanceAttachment()`
- `DeletePerformanceAttachment()`
- `CreatePerformanceRating()`
- `GetPerformanceRatingByID()`
- `ListPerformanceRatings()`
- `UpdatePerformanceRating()`
- `DeletePerformanceRating()`
- `CreatePerformanceIndicatorFormula()`
- `GetPerformanceIndicatorFormulaByID()`
- `ListPerformanceIndicatorFormulas()`
- `UpdatePerformanceIndicatorFormula()`
- `DeletePerformanceIndicatorFormula()`
- `GetPerformanceLogByID()`
- `ListPerformanceLogs()`
- `ListPerformanceLogsByEvaluationID()`
- `CreateEvaluationWithSnapshot()`
- `GetEvaluationWithDetails()`
- `UpdateEvaluationTarget()`
- `UpdateEvaluationActual()`
- `BulkUpdateEvaluationActuals()`
- `RecalculateEvaluationScore()`
- `GetEvaluationProgressSummary()`
- `SubmitTarget()`
- `HandleTargetApprovalStatusChange()`
- `ApproveTarget()`
- `RejectTarget()`
- `SubmitEvaluation()`
- `HandleRealizationApprovalStatusChange()`
- `ApproveEvaluation()`
- `RejectEvaluation()`
- `CompleteEvaluation()`
- `propagateSubordinateScoreUpwardBestEffort()`
- `propagateSubordinateScoreUpward()`
- `RecalculatePeriodScoring()`
- `GetEmployeeDashboard()`
- `GetManagerDashboard()`
- `GetHRDashboard()`
- `CreatePerformanceComponent()`
- `GetPerformanceComponentByID()`
- `ListPerformanceComponents()`
- `UpdatePerformanceComponent()`
- `DeletePerformanceComponent()`
- `UpsertOrganizationComponent()`
- `ListOrganizationComponents()`
- `DeleteOrganizationComponent()`
- `ListEvaluationComponents()`
- `UpdateEvaluationComponentScore()`
- `CalculateEvaluationComponentScoring()`

### rbac.Service — 8 methods
- `ListRoles()`
- `CreateRole()`
- `UpdateRole()`
- `DeleteRole()`
- `ListPermissions()`
- `AssignRolePermissions()`
- `ListUsers()`
- `AssignUserRoles()`

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

### reimbursement.Service — 17 methods
- `SetApprovalEngine()`
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
- `HandleApprovalStatusChange()`
- `DeleteReimbursementRequest()`
- `CreateReimbursementItem()`
- `ListReimbursementItems()`
- `UpdateReimbursementItem()`
- `DeleteReimbursementItem()`

### setting.Service — 119 methods
- `validateUniqueCode()`
- `validateUniqueCodeExcludeSelf()`
- `CreateZone()`
- `GetZoneByID()`
- `ListZones()`
- `ListZonesActive()`
- `UpdateZone()`
- `DeleteZone()`
- `CreateProvince()`
- `GetProvinceByID()`
- `ListProvinces()`
- `ListAllProvinces()`
- `UpdateProvince()`
- `DeleteProvince()`
- `CreateRegency()`
- `GetRegencyByID()`
- `ListRegencies()`
- `ListRegenciesByProvince()`
- `UpdateRegency()`
- `DeleteRegency()`
- `CreateDistrict()`
- `GetDistrictByID()`
- `ListDistricts()`
- `ListDistrictsByRegency()`
- `UpdateDistrict()`
- `DeleteDistrict()`
- `CreateVillage()`
- `GetVillageByID()`
- `ListVillages()`
- `ListVillagesByDistrict()`
- `UpdateVillage()`
- `DeleteVillage()`
- `GetVillageDetail()`
- `SearchVillages()`
- `CreateEducation()`
- `GetEducationByID()`
- `ListEducations()`
- `UpdateEducation()`
- `DeleteEducation()`
- `CreateEducationMajor()`
- `GetEducationMajorByID()`
- `ListEducationMajors()`
- `UpdateEducationMajor()`
- `DeleteEducationMajor()`
- `CreateReligion()`
- `GetReligionByID()`
- `ListReligions()`
- `UpdateReligion()`
- `DeleteReligion()`
- `CreateMaritalStatus()`
- `GetMaritalStatusByID()`
- `ListMaritalStatuses()`
- `UpdateMaritalStatus()`
- `DeleteMaritalStatus()`
- `CreateBank()`
- `GetBankByID()`
- `ListBanks()`
- `UpdateBank()`
- `DeleteBank()`
- `CreateNationality()`
- `GetNationalityByID()`
- `ListNationalities()`
- `UpdateNationality()`
- `DeleteNationality()`
- `CreateRelationshipType()`
- `GetRelationshipTypeByID()`
- `ListRelationshipTypes()`
- `UpdateRelationshipType()`
- `DeleteRelationshipType()`
- `CreateEmploymentStatus()`
- `GetEmploymentStatusByID()`
- `ListEmploymentStatuses()`
- `UpdateEmploymentStatus()`
- `DeleteEmploymentStatus()`
- `CreateJobFamily()`
- `GetJobFamilyByID()`
- `ListJobFamilies()`
- `UpdateJobFamily()`
- `DeleteJobFamily()`
- `CreateGrading()`
- `GetGradingByID()`
- `ListGradings()`
- `UpdateGrading()`
- `DeleteGrading()`
- `CreateSalaryGrade()`
- `GetSalaryGradeByID()`
- `ListSalaryGrades()`
- `UpdateSalaryGrade()`
- `DeleteSalaryGrade()`
- `CreateTER()`
- `GetTERByID()`
- `ListTERs()`
- `UpdateTER()`
- `DeleteTER()`
- `CreatePTKP()`
- `GetPTKPByID()`
- `ListPTKPs()`
- `CreateInsurance()`
- `GetInsuranceByID()`
- `ListInsurances()`
- `UpdateInsurance()`
- `DeleteInsurance()`
- `CreateCompanyHoliday()`
- `GetCompanyHolidayByID()`
- `ListCompanyHolidays()`
- `UpdateCompanyHoliday()`
- `DeleteCompanyHoliday()`
- `ListHolidayDatesInRange()`
- `validateUniqueName()`
- `validateUniqueNameExcludeSelf()`
- `validateUniqueHolidayDate()`
- `validateUniqueHolidayDateExcludeSelf()`
- `CreateCompetency()`
- `GetCompetencyByID()`
- `ListCompetencies()`
- `UpdateCompetency()`
- `DeleteCompetency()`
- `UpdatePTKP()`
- `DeletePTKP()`

### training.Service — 129 methods
- `SetApprovalEngine()`
- `SetNotifier()`
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
- `CreateProvider()`
- `GetProviderByID()`
- `ListProviders()`
- `UpdateProvider()`
- `DeleteProvider()`
- `CreateTrainer()`
- `GetTrainerByID()`
- `ListTrainers()`
- `UpdateTrainer()`
- `DeleteTrainer()`
- `AddSessionTrainer()`
- `ListSessionTrainers()`
- `RemoveSessionTrainer()`
- `MarkAttendance()`
- `markOneAttendance()`
- `UpdateAttendance()`
- `ListAttendanceBySession()`
- `CreateAssessment()`
- `ListAssessmentsBySession()`
- `SubmitAssessmentResult()`
- `CreatePlan()`
- `GetPlanByID()`
- `ListPlans()`
- `UpdatePlan()`
- `DeletePlan()`
- `CreatePlanItem()`
- `ListPlanItems()`
- `UpdatePlanItem()`
- `DeletePlanItem()`
- `CreateNeed()`
- `GetNeedByID()`
- `ListNeeds()`
- `DeleteNeed()`
- `UpdateNeed()`
- `CreateRequest()`
- `GetRequestByID()`
- `ListRequests()`
- `SubmitRequest()`
- `CancelRequest()`
- `HandleApprovalStatusChange()`
- `autoEnrollApprovedRequest()`
- `CreateCourseObjective()`
- `ListCourseObjectives()`
- `UpdateCourseObjective()`
- `DeleteCourseObjective()`
- `CreateCourseCompetency()`
- `ListCourseCompetencies()`
- `DeleteCourseCompetency()`
- `CreateCoursePrerequisite()`
- `ListCoursePrerequisites()`
- `DeleteCoursePrerequisite()`
- `CreateMandatory()`
- `GetMandatoryByID()`
- `ListMandatories()`
- `UpdateMandatory()`
- `DeleteMandatory()`
- `CreateSessionCost()`
- `ListSessionCosts()`
- `UpdateSessionCost()`
- `DeleteSessionCost()`
- `CreateDocument()`
- `ListDocuments()`
- `DeleteDocument()`
- `CreateEvaluationForm()`
- `GetEvaluationFormByID()`
- `ListEvaluationForms()`
- `UpdateEvaluationForm()`
- `DeleteEvaluationForm()`
- `GetEvaluationFormBySession()`
- `CreateEvaluationQuestion()`
- `ListEvaluationQuestions()`
- `UpdateEvaluationQuestion()`
- `DeleteEvaluationQuestion()`
- `SubmitEvaluationAnswers()`
- `ListEvaluationAnswers()`
- `CreateEffectivenessAssessment()`
- `GetEffectivenessAssessmentByID()`
- `ListEffectivenessAssessments()`
- `UpdateEffectivenessAssessment()`
- `DeleteEffectivenessAssessment()`
- `CreateCertification()`
- `GetCertificationByID()`
- `ListCertifications()`
- `UpdateCertification()`
- `DeleteCertification()`
- `GenerateCertificate()`
- `UpdateCertificateFile()`
- `GetTrainingHistory()`
- `GetParticipationReport()`
- `GetCostReport()`
- `GetComplianceReport()`
- `GetDashboardReport()`

### user.Service — 12 methods
- `Login()`
- `RefreshToken()`
- `CreateUser()`
- `GetUser()`
- `ListUsers()`
- `enrichWithCompany()`
- `CreateCompanyUser()`
- `DeleteUser()`
- `UpdateUser()`
- `ChangePassword()`
- `EnsureMigrate()`
- `EnsureSeed()`

### useraccount.Service — 10 methods
- `CreateAccount()`
- `ResendSetupEmail()`
- `GetAccountStatus()`
- `GetMyAccount()`
- `SetPassword()`
- `Login()`
- `loginTenantUser()`
- `loginPlatformUser()`
- `Refresh()`
- `tryLinkEmployeeUser()`

### workforceintelligence.Service — 47 methods
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
- `GetKPISummary()`
- `ListKPIs()`
- `GetHeadcountAnalytics()`
- `GetMovementAnalytics()`
- `CreateScenario()`
- `RunScenario()`
- `executeSimulation()`
- `GetRiskDashboard()`
- `GetExecutiveSummary()`
- `GetProjections()`
- `ListRiskIndicators()`
- `GetRiskIndicatorByID()`
- `UpdateRiskIndicator()`
- `ListScenarios()`
- `GetScenarioByID()`
- `UpdateScenario()`
- `DeleteScenario()`
- `CloneScenario()`
- `GetPeopleAnalytics()`
- `GetCapacityForecast()`
- `GetPayrollCostBreakdown()`
- `GetCostPerEmployee()`
- `GetExecutiveGrowth()`
- `GetExecutiveCostTrend()`
- `GetExecutiveAttritionTrend()`
- `GetExecutiveCapacity()`
- `GetExecutiveRiskOverview()`
- `GetExecutiveHealthScore()`
- `GetHealthScoreByID()`
- `GetSpanOfControl()`
- `GetSuccessionReadiness()`
- `GetRiskDetail()`
- `GetKPIByCode()`
- `GetHealthDashboard()`
- `ListHealthScores()`
- `CandidateSearch()`

## SECTION 6: GRAND TOTALS

| Category | Count |
|----------|:-----:|
| Tenant Modules | 19 |
| Platform Modules | 6 |
| Shared Kernel Packages | 23 |
| **Total Architecture Layers** | **48** |
| Total GORM Entities (tenant) | 191 |
| Total GORM Entities (platform) | 8 |
| **Total Entities (combined)** | **199** |
| Total Service Methods | 961 |
| Total Repository Methods | 1061 |
| Total Handler Functions | 876 |
| Total Route Registrations | 915 |
| **Total Unit Tests (all)** | **1547** |
| Total Go Source Files | 232 |
| Total Test Files (_test.go) | 128 |
| **Total Go Files** | **360** |

## SECTION 7: TEST FILE INVENTORY

| File | Test Funcs |
|------|:----------:|
| `internal\modules\approval\errors_test.go` | 4 |
| `internal\modules\approval\handler_test.go` | 14 |
| `internal\modules\approval\helpers_test.go` | 0 |
| `internal\modules\approval\hierarchy_test.go` | 7 |
| `internal\modules\approval\module_subscription_test.go` | 10 |
| `internal\modules\approval\notifier_test.go` | 5 |
| `internal\modules\approval\repository_test.go` | 25 |
| `internal\modules\approval\service_test.go` | 32 |
| `internal\modules\approval\status_handler_test.go` | 8 |
| `internal\modules\attendance\approval_integration_test.go` | 14 |
| `internal\modules\attendance\correction_test.go` | 3 |
| `internal\modules\attendance\handler_test.go` | 26 |
| `internal\modules\attendance\helpers_test.go` | 0 |
| `internal\modules\attendance\notifier_integration_test.go` | 4 |
| `internal\modules\attendance\overtime_two_flow_test.go` | 15 |
| `internal\modules\attendance\repository_test.go` | 34 |
| `internal\modules\attendance\service_test.go` | 39 |
| `internal\modules\attendance\session_test.go` | 8 |
| `internal\modules\attendance\summary_test.go` | 3 |
| `internal\modules\careerintelligence\handler_test.go` | 22 |
| `internal\modules\careerintelligence\helpers_test.go` | 0 |
| `internal\modules\careerintelligence\repository_test.go` | 24 |
| `internal\modules\careerintelligence\service_test.go` | 20 |
| `internal\modules\competency\handler_test.go` | 13 |
| `internal\modules\competency\helpers_test.go` | 0 |
| `internal\modules\competency\repository_test.go` | 12 |
| `internal\modules\competency\service_test.go` | 35 |
| `internal\modules\employee\helpers_test.go` | 0 |
| `internal\modules\employee\repository_test.go` | 18 |
| `internal\modules\employee\service_test.go` | 22 |
| `internal\modules\employeemovement\approval_integration_test.go` | 10 |
| `internal\modules\employeemovement\audit_test.go` | 10 |
| `internal\modules\employeemovement\career_test.go` | 4 |
| `internal\modules\employeemovement\document_test.go` | 6 |
| `internal\modules\employeemovement\eligibility_test.go` | 7 |
| `internal\modules\employeemovement\enrichment_test.go` | 6 |
| `internal\modules\employeemovement\expiry_test.go` | 3 |
| `internal\modules\employeemovement\handler_test.go` | 16 |
| `internal\modules\employeemovement\helpers_test.go` | 0 |
| `internal\modules\employeemovement\report_test.go` | 12 |
| `internal\modules\employeemovement\repository_test.go` | 27 |
| `internal\modules\employeemovement\service_test.go` | 49 |
| `internal\modules\jobmanagement\calculator_test.go` | 13 |
| `internal\modules\jobmanagement\education_update_test.go` | 1 |
| `internal\modules\jobmanagement\handler_test.go` | 10 |
| `internal\modules\jobmanagement\helpers_test.go` | 0 |
| `internal\modules\jobmanagement\repository_test.go` | 27 |
| `internal\modules\jobmanagement\router_test.go` | 3 |
| `internal\modules\jobmanagement\service_test.go` | 29 |
| `internal\modules\leave\approval_integration_test.go` | 10 |
| `internal\modules\leave\attendance_integration_test.go` | 2 |
| `internal\modules\leave\balance_test.go` | 3 |
| `internal\modules\leave\calculation_test.go` | 8 |
| `internal\modules\leave\calendar_test.go` | 2 |
| `internal\modules\leave\handler_test.go` | 11 |
| `internal\modules\leave\helpers_test.go` | 0 |
| `internal\modules\leave\notifier_integration_test.go` | 6 |
| `internal\modules\leave\report_test.go` | 2 |
| `internal\modules\leave\repository_test.go` | 17 |
| `internal\modules\leave\service_test.go` | 18 |
| `internal\modules\notification\helpers_test.go` | 0 |
| `internal\modules\notification\repository_test.go` | 2 |
| `internal\modules\notification\service_test.go` | 6 |
| `internal\modules\organization\helpers_test.go` | 0 |
| `internal\modules\organization\service_test.go` | 9 |
| `internal\modules\payroll\helpers_test.go` | 0 |
| `internal\modules\payroll\repository_test.go` | 17 |
| `internal\modules\payroll\service_test.go` | 26 |
| `internal\modules\performance\evaluation_enrichment_test.go` | 3 |
| `internal\modules\performance\handler_test.go` | 14 |
| `internal\modules\performance\helpers_test.go` | 0 |
| `internal\modules\performance\kpi_approval_routing_test.go` | 8 |
| `internal\modules\performance\kpi_two_phase_test.go` | 5 |
| `internal\modules\performance\my_kpi_context_test.go` | 3 |
| `internal\modules\performance\notifier_integration_test.go` | 28 |
| `internal\modules\performance\okr_approval_routing_test.go` | 7 |
| `internal\modules\performance\okr_eligibility_test.go` | 2 |
| `internal\modules\performance\okr_objective_scope_test.go` | 6 |
| `internal\modules\performance\okr_two_phase_test.go` | 3 |
| `internal\modules\performance\repository_test.go` | 19 |
| `internal\modules\performance\service_test.go` | 27 |
| `internal\modules\recruitment\handler_test.go` | 28 |
| `internal\modules\recruitment\helpers_test.go` | 0 |
| `internal\modules\recruitment\repository_test.go` | 27 |
| `internal\modules\recruitment\service_test.go` | 20 |
| `internal\modules\reimbursement\approval_integration_test.go` | 7 |
| `internal\modules\reimbursement\handler_test.go` | 17 |
| `internal\modules\reimbursement\helpers_test.go` | 0 |
| `internal\modules\reimbursement\repository_test.go` | 19 |
| `internal\modules\reimbursement\service_test.go` | 24 |
| `internal\modules\setting\competency_test.go` | 8 |
| `internal\modules\setting\helpers_test.go` | 0 |
| `internal\modules\setting\repository_test.go` | 18 |
| `internal\modules\setting\service_test.go` | 26 |
| `internal\modules\training\advanced_p2_test.go` | 7 |
| `internal\modules\training\handler_test.go` | 9 |
| `internal\modules\training\helpers_test.go` | 1 |
| `internal\modules\training\plan_p1_test.go` | 23 |
| `internal\modules\training\repository_test.go` | 14 |
| `internal\modules\training\service_test.go` | 17 |
| `internal\modules\workforceintelligence\handler_test.go` | 36 |
| `internal\modules\workforceintelligence\helpers_test.go` | 0 |
| `internal\modules\workforceintelligence\repository_test.go` | 33 |
| `internal\modules\workforceintelligence\service_test.go` | 43 |
| `internal\platform\company\handler_test.go` | 8 |
| `internal\platform\company\helpers_test.go` | 0 |
| `internal\platform\company\service_test.go` | 22 |
| `internal\platform\modulemgmt\repository_test.go` | 3 |
| `internal\platform\package\helpers_test.go` | 0 |
| `internal\platform\package\service_test.go` | 27 |

====================================================================================================
  END OF REPORT
