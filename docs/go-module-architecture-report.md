====================================================================================================
  HRIS PLATFORM — GO MODULE ARCHITECTURE REPORT
  Generated: 19 Aug 2026

  Index dokumentasi: docs/README.md  |  Terkait: platform-architecture-design.md, openapi-report.md
====================================================================================================

## SECTION 1: TENANT MODULES (internal/modules/)

| Module | Entities | Service Methods | Repo Methods | Handler Funcs | Route Regs | Tests |
|--------|:--------:|:--------------:|:------------:|:-------------:|:----------:|:-----:|
| approval | 6 | 30 | 45 | 18 | 18 | 119 |
| attendance | 11 | 122 | 123 | 97 | 98 | 151 |
| careerintelligence | 5 | 27 | 32 | 23 | 23 | 82 |
| competency | 7 | 77 | 93 | 67 | 67 | 81 |
| documenttemplate | 4 | 12 | 18 | 14 | 13 | 51 |
| employee | 10 | 50 | 57 | 40 | 40 | 90 |
| employeemovement | 6 | 51 | 51 | 30 | 29 | 158 |
| jobmanagement | 23 | 101 | 99 | 97 | 97 | 85 |
| leave | 7 | 39 | 34 | 26 | 26 | 81 |
| notification | 1 | 5 | 5 | 5 | 4 | 8 |
| organization | 3 | 20 | 26 | 18 | 18 | 9 |
| payroll | 30 | 108 | 110 | 85 | 83 | 120 |
| performance | 17 | 109 | 126 | 99 | 147 | 125 |
| rbac | 5 | 8 | 15 | 8 | 8 | 2 |
| recruitment | 24 | 122 | 119 | 101 | 101 | 272 |
| reimbursement | 3 | 19 | 20 | 15 | 15 | 78 |
| setting | 21 | 119 | 126 | 107 | 113 | 60 |
| training | 28 | 130 | 145 | 123 | 123 | 79 |
| useraccount | 3 | 11 | 20 | 7 | 4 | 0 |
| workforceintelligence | 7 | 52 | 48 | 72 | 70 | 137 |
| **TOTAL** | **221** | **1212** | **1312** | **1052** | **1097** | **1788** |

### Test Breakdown per Module

| Module | Repo Tests | Service Tests | Handler Tests | Other | Total |
|--------|:----------:|:-------------:|:-------------:|:-----:|:-----:|
| approval | 25 | 35 | 22 | 37 | 119 |
| attendance | 34 | 39 | 26 | 52 | 151 |
| careerintelligence | 31 | 29 | 22 | 0 | 82 |
| competency | 12 | 37 | 13 | 19 | 81 |
| documenttemplate | 7 | 15 | 10 | 19 | 51 |
| employee | 22 | 37 | 3 | 28 | 90 |
| employeemovement | 27 | 49 | 16 | 66 | 158 |
| jobmanagement | 29 | 29 | 10 | 17 | 85 |
| leave | 17 | 18 | 11 | 35 | 81 |
| notification | 2 | 6 | 0 | 0 | 8 |
| organization | 0 | 9 | 0 | 0 | 9 |
| payroll | 23 | 37 | 0 | 60 | 120 |
| performance | 19 | 27 | 14 | 65 | 125 |
| rbac | 0 | 2 | 0 | 0 | 2 |
| recruitment | 62 | 153 | 57 | 0 | 272 |
| reimbursement | 19 | 27 | 17 | 15 | 78 |
| setting | 18 | 30 | 2 | 10 | 60 |
| training | 14 | 21 | 10 | 34 | 79 |
| workforceintelligence | 46 | 55 | 36 | 0 | 137 |

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
| authctx | 2 | 4 | Auth context helpers (user/company from gin context) |
| authz | 17 | 120 | Casbin RBAC authorization |
| cache | 6 | 51 | Distributed cache + Pub/Sub |
| config | 1 | 0 | Viper configuration loader |
| crypto | 2 | 9 | AES-256-GCM encryption |
| database | 1 | 0 | Multi-tenant DB connection manager |
| docs | 1 | 0 | OpenAPI/Scalar documentation |
| driver | 1 | 0 | DB driver detection |
| errors | 0 | 0 | Shared error helpers |
| httputil | 4 | 0 | Bilingual response helpers (SuccessJSON, CreatedJSON, ErrorJSON, NotFound) + locale message catalog (80+ EN/ID pairs) + custom Indonesian validators (NIK, NPWP, KK, Passport, SIM, No Rekening) |
| logger | 1 | 0 | Structured logging |
| mailer | 1 | 0 | Email sending (SMTP/template) |
| mask | 2 | 1 | menyediakan utilitas untuk menyamarkan sebagian nilai |
| middleware | 10 | 16 | HTTP middleware: Auth, CORS, Logger, Recovery, Tenant, Localize (auto-detect Accept-Language) |
| migrator | 3 | 3 | Database migration engine |
| module | 1 | 0 | Module SDK |
| numbering | 5 | 8 | numbering |
| onpremise | 2 | 6 | On-premise license enforcement |
| router | 1 | 0 | Router setup |
| telemetry | 0 | 0 | Telemetry/metrics |
| tenant | 0 | 0 | Tenant resolution helpers |
| tenantseed | 8 | 11 | Tenant seed data (nationalities, competencies, RBAC) |
| upload | 1 | 0 | menyediakan endpoint upload file generik untuk lampiran |
| validator | 0 | 0 | Validator helpers |
| **TOTAL** | **71** | **229** | |

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

### documenttemplate
- DocumentTemplate
- DocumentTemplateVersion
- DocumentTemplateAudit
- GeneratedDocument

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
- TerRate
- Ptkp
- Pph21TaxBracket
- PayrollRun
- PayrollRunEmployee
- PayrollRunItem
- PayrollPayslip
- PayrollPayment
- Pph21CalculationLog
- PayrollProfileChangeLog
- EmployeeRead
- EmploymentRead
- PositionRead
- GradingRead
- AttendanceSessionRead
- LeaveRequestRead
- LeaveRequestDetailRead

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
- JobOffer
- Candidate
- CandidateEducation
- CandidateWorkExperience
- CandidateSkill
- CandidateCertification
- CandidateDocument
- CandidateConsent
- JobApplication
- JobRequisitionRequirement
- JobRequisitionCompetency
- RecruitmentStage
- ApplicationStageHistory
- ApplicationScreening
- ApplicationAssessment
- RecruitmentAssessment
- AssessmentParticipant
- Interview
- Interviewer
- InterviewScorecardItem
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

### approval.Service — 30 methods
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
- `ListMyDoneTasks()`
- `listMyTasks()`
- `resolveStepAssignees()`
- `resolveSupervisorAssignees()`
- `resolveOrganizationAssignees()`
- `advanceThroughWatcherSteps()`

### attendance.Service — 122 methods
- `SetApprovalEngine()`
- `SetModuleChecker()`
- `SetPayrollAdjuster()`
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
- `GetAttendanceStats()`
- `GetOvertimeTrend()`
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
- `CreateBusinessTravel()`
- `GetBusinessTravelByID()`
- `ListBusinessTravels()`
- `UpdateBusinessTravel()`
- `SubmitBusinessTravel()`
- `CancelBusinessTravel()`
- `AddBusinessTravelParticipant()`
- `requireDraftTravel()`
- `UpdateBusinessTravelParticipant()`
- `DeleteBusinessTravelParticipant()`
- `ListBusinessTravelParticipants()`
- `AddBusinessTravelDestination()`
- `UpdateBusinessTravelDestination()`
- `DeleteBusinessTravelDestination()`
- `ListBusinessTravelDestinations()`
- `AddBusinessTravelActivity()`
- `UpdateBusinessTravelActivity()`
- `DeleteBusinessTravelActivity()`
- `ListBusinessTravelActivities()`
- `AddBusinessTravelSchedule()`
- `UpdateBusinessTravelSchedule()`
- `DeleteBusinessTravelSchedule()`
- `ListBusinessTravelSchedules()`
- `CreateFundingMethod()`
- `ListFundingMethods()`
- `CreateFunding()`
- `ListFundingsByTravel()`
- `UpdateFunding()`
- `ConfirmFunding()`
- `AddFundingDocument()`
- `CreateExpenseCategory()`
- `ListExpenseCategories()`
- `CreateExpense()`
- `ListExpensesByTravel()`
- `UpdateExpense()`
- `DeleteExpense()`
- `AddExpenseDocument()`
- `CreateSettlement()`
- `GetSettlementByID()`
- `ListSettlementsByTravel()`
- `SubmitSettlement()`
- `HandleSettlementApprovalStatusChange()`
- `ListRefundsByTravel()`
- `ConfirmRefund()`
- `ListTravelReimbursements()`
- `ApproveTravelReimbursement()`
- `ProcessTravelReimbursement()`
- `PayTravelReimbursement()`
- `maybeSettleAndCloseTravel()`
- `AddTravelDocument()`
- `ListTravelDocuments()`
- `DeleteTravelDocument()`
- `pushBusinessTravelPayrollAdjustments()`
- `HandleBusinessTravelApprovalStatusChange()`
- `notifyBusinessTravelUser()`
- `resolveBusinessTravelRecipientUserID()`
- `pushBusinessTravelAttendance()`
- `GetBusinessTravelReport()`
- `GetBusinessTravelFundingReport()`
- `GetBusinessTravelAdvanceReport()`
- `GetBusinessTravelReimbursementReport()`
- `GetBusinessTravelRefundReport()`
- `GetBusinessTravelCostReport()`
- `recalculateSession()`
- `ApplyApprovedLeave()`
- `ApplyApprovedBusinessTravel()`

### careerintelligence.Service — 27 methods
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
- `GetEligibleEmployeesForPath()`
- `GetEligibleEmployeesByPosition()`
- `CreateSuccessionPlan()`
- `ListSuccessionPlans()`
- `GetSuccessionPlanByID()`
- `UpdateSuccessionPlan()`
- `DeleteSuccessionPlan()`
- `GetSuccessionGaps()`
- `CheckSuccessionGapByPosition()`
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

### competency.Service — 77 methods
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
- `ensureSelfRater()`
- `GetCompetencyEventTargetByID()`
- `ListCompetencyEventTargets()`
- `buildTargetRaterSummaries()`
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
- `CreateRatingScale()`
- `GetRatingScaleByID()`
- `ListRatingScales()`
- `UpdateRatingScale()`
- `DeleteRatingScale()`
- `CreateAssessmentTemplate()`
- `GetAssessmentTemplateByID()`
- `ListAssessmentTemplates()`
- `UpdateAssessmentTemplate()`
- `DeleteAssessmentTemplate()`
- `replaceTemplateChildren()`
- `CreateIndicator()`
- `GetIndicatorByID()`
- `ListIndicators()`
- `UpdateIndicator()`
- `DeleteIndicator()`
- `SetTemplateIndicators()`
- `ListTemplateIndicators()`
- `SetApprovalEngine()`
- `SubmitAssessmentForApproval()`
- `HandleAssessmentApprovalStatusChange()`
- `AssignRaters()`
- `SuggestedRaters()`
- `templateRequiresSelf()`
- `ListRatersByTarget()`
- `DeleteRater()`
- `MyAssessments()`
- `ManagerAssessments()`
- `GetAssessmentDetail()`
- `SaveResponses()`
- `templateIndicatorSet()`
- `SubmitAssessment()`
- `enrichRaterNames()`
- `CalculateTarget()`
- `FinalizeTarget()`
- `GetEmployeeResult()`
- `GetEmployeeGap()`
- `GetEmployeeReport()`
- `GetManagerReport()`
- `GetHRReport()`

### documenttemplate.Service — 12 methods
- `List()`
- `GetByID()`
- `checkCodeUnique()`
- `validateMovementType()`
- `Create()`
- `Update()`
- `Delete()`
- `Activate()`
- `Deactivate()`
- `CreateVersion()`
- `ListVersions()`
- `GetVersion()`

### employee.Service — 50 methods
- `resolveEmployeeID()`
- `isSelfEmployee()`
- `isSelfEmployeeStr()`
- `ListSensitiveFieldSettings()`
- `SetSensitiveFieldEnabled()`
- `IsFieldEncryptionEnabled()`
- `SetQuotaChecker()`
- `SetEmployeeIDFormatProvider()`
- `checkQuota()`
- `encryptIfEnabled()`
- `Create()`
- `GetEmploymentStatusStats()`
- `GetGenderStats()`
- `fillRefNames()`
- `GetByID()`
- `GetByIDUnmasked()`
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

### employeemovement.Service — 51 methods
- `SetApprovalEngine()`
- `SetCareerExecutor()`
- `SetNotifier()`
- `SetPerformanceProvider()`
- `SetCompetencyProvider()`
- `SetOKRProvider()`
- `SetDocumentGenerator()`
- `SetNumberingService()`
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
- `GenerateMovementDocument()`
- `GenerateContractDocument()`
- `ListGeneratedDocuments()`
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

### jobmanagement.Service — 101 methods
- `GetDashboard()`
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

### leave.Service — 39 methods
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
- `GetOnLeaveToday()`
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

### payroll.Service — 108 methods
- `calculateBpjsContributions()`
- `buildBpjsRunItem()`
- `CalculatePayrollRun()`
- `selectRunEmployees()`
- `findActivePayrollProfiles()`
- `loadAllActiveComponents()`
- `validateComponentFormulas()`
- `calculateEmployee()`
- `employmentProrationFactor()`
- `buildComponentInputs()`
- `evaluateComponents()`
- `tryEvaluateComponent()`
- `evaluateFormulaWithValues()`
- `buildRunItem()`
- `persistRunSnapshot()`
- `CreatePaymentBatch()`
- `UpdatePaymentStatus()`
- `ListPaymentsByRun()`
- `GetPaymentByID()`
- `ExportPaymentsCSV()`
- `GeneratePayslips()`
- `PublishPayslip()`
- `CancelPayslip()`
- `ListPayslipsByRun()`
- `GetPayslipByID()`
- `GetPayslipHTML()`
- `calculatePph21()`
- `calculatePph21TerMonthly()`
- `buildPph21RunItem()`
- `GetPayrollSummaryReport()`
- `GetPayrollDashboard()`
- `GetPayrollDetailReport()`
- `GetBpjsReport()`
- `GetTaxReport()`
- `GetBankTransferReport()`
- `SetFormulaEngine()`
- `SetApprovalEngine()`
- `SetNotifier()`
- `CreateSalaryComponent()`
- `GetSalaryComponentByID()`
- `ListSalaryComponents()`
- `UpdateSalaryComponent()`
- `DeleteSalaryComponent()`
- `validateSalaryComponentCalculation()`
- `ValidateFormula()`
- `ListFormulaVariables()`
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
- `ListEmployeeBankProfiles()`
- `ListEmployeeBpjsProfiles()`
- `ListEmployeeTaxProfiles()`
- `buildPaginated()`
- `ListSalaryGradeComponents()`
- `GetSalaryGradeComponentByID()`
- `CreateSalaryGradeComponent()`
- `UpdateSalaryGradeComponent()`
- `DeleteSalaryGradeComponent()`
- `ListSalaryEmployeeComponents()`
- `GetSalaryEmployeeComponentByID()`
- `CreateSalaryEmployeeComponent()`
- `UpdateSalaryEmployeeComponent()`
- `DeleteSalaryEmployeeComponent()`
- `CreateBpjsSetting()`
- `CreatePph21Setting()`
- `CreatePayrollRun()`
- `ListPayrollRuns()`
- `GetPayrollRunByID()`
- `ListPayrollRunEmployees()`
- `ListPayrollRunItems()`
- `UpdatePayrollRunStatus()`
- `CheckPayrollRunApproval()`
- `HandleApprovalStatusChange()`
- `notifyRunOutcome()`
- `UpdateBpjsSetting()`
- `GetBpjsSettingByID()`
- `ListBpjsSettings()`
- `DeleteBpjsSetting()`
- `GetBpjsRateComponentByID()`
- `ListBpjsRateComponentsBySettingID()`
- `UpdateBpjsRateComponent()`
- `DeleteBpjsRateComponent()`
- `CreateBpjsRateComponent()`
- `GetPph21SettingByID()`
- `ListPph21Settings()`
- `UpdatePph21Setting()`
- `DeletePph21Setting()`
- `CreatePph21TaxBracket()`
- `UpdatePph21TaxBracket()`
- `DeletePph21TaxBracket()`
- `ListPph21TaxBrackets()`
- `loadWorkforceSummary()`

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

### recruitment.Service — 122 methods
- `SetWorkforceGapProvider()`
- `SetInternalCandidateProvider()`
- `SetSuccessionGapProvider()`
- `SetTrainingHandoffProvider()`
- `SetApprovalEngine()`
- `SetNotifier()`
- `SetEmployeeProvider()`
- `SetMovementProvider()`
- `SetJobManagementProvider()`
- `CreateRequisition()`
- `GetRequisitionByID()`
- `ListRequisitions()`
- `UpdateRequisition()`
- `DeleteRequisition()`
- `CreateOffer()`
- `GetOfferByID()`
- `ListOffers()`
- `UpdateOffer()`
- `DeleteOffer()`
- `SubmitOffer()`
- `HandleOfferApprovalStatusChange()`
- `SendOffer()`
- `AcceptOffer()`
- `RejectOffer()`
- `WithdrawOffer()`
- `SubmitRequisition()`
- `HandleApprovalStatusChange()`
- `notifyRequisitionOutcome()`
- `CreateCandidate()`
- `GetCandidateByID()`
- `ListCandidates()`
- `UpdateCandidate()`
- `DeleteCandidate()`
- `CreateApplication()`
- `GetApplicationByID()`
- `ListApplications()`
- `transitionApplicationStatus()`
- `UpdateApplicationStatus()`
- `GetApplicationHistory()`
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
- `handoffOnboardingTraining()`
- `DeleteEmployeeOnboarding()`
- `CreateOnboardingTaskItem()`
- `ListOnboardingTaskItems()`
- `UpdateOnboardingTaskItem()`
- `DeleteOnboardingTaskItem()`
- `resolveWorkforceGapSlots()`
- `validateSuccessionGapFallback()`
- `GetEligibleInternalCandidates()`
- `handoffHiredEmployee()`
- `CreateCandidateEducation()`
- `ListCandidateEducations()`
- `UpdateCandidateEducation()`
- `DeleteCandidateEducation()`
- `CreateCandidateWorkExperience()`
- `ListCandidateWorkExperiences()`
- `UpdateCandidateWorkExperience()`
- `DeleteCandidateWorkExperience()`
- `CreateCandidateSkill()`
- `ListCandidateSkills()`
- `UpdateCandidateSkill()`
- `DeleteCandidateSkill()`
- `CreateCandidateCertification()`
- `ListCandidateCertifications()`
- `UpdateCandidateCertification()`
- `DeleteCandidateCertification()`
- `CreateCandidateDocument()`
- `ListCandidateDocuments()`
- `UpdateCandidateDocument()`
- `DeleteCandidateDocument()`
- `CreateCandidateConsent()`
- `ListCandidateConsents()`
- `filterTemplatesForApplication()`
- `CreateApplicationScreening()`
- `ListApplicationScreenings()`
- `UpdateApplicationScreening()`
- `DeleteApplicationScreening()`
- `assessmentRequirementStrings()`
- `effectiveAssessmentCompetencies()`
- `GetApplicationAssessment()`
- `SaveApplicationAssessment()`
- `CreateAssessment()`
- `ListAssessments()`
- `GetAssessmentByID()`
- `UpdateAssessment()`
- `DeleteAssessment()`
- `AddAssessmentParticipant()`
- `ListAssessmentParticipants()`
- `UpdateAssessmentParticipant()`
- `DeleteAssessmentParticipant()`
- `GetRecruitmentAnalyticsSummary()`
- `CreateRequisitionRequirement()`
- `ListRequisitionRequirements()`
- `UpdateRequisitionRequirement()`
- `DeleteRequisitionRequirement()`
- `CreateRequisitionCompetency()`
- `ListRequisitionCompetencies()`
- `UpdateRequisitionCompetency()`
- `DeleteRequisitionCompetency()`
- `GetCandidateMatchScore()`
- `assessmentCandidateLevels()`
- `AddInterviewer()`
- `ListInterviewers()`
- `RemoveInterviewer()`
- `AddScorecardItem()`
- `ListScorecardItems()`
- `UpdateScorecardItem()`
- `DeleteScorecardItem()`
- `CompleteInterview()`

### reimbursement.Service — 19 methods
- `SetApprovalEngine()`
- `SetNotifier()`
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
- `notifyReimbursementOutcome()`
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

### training.Service — 130 methods
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
- `CreateOnboardingNeed()`
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

### useraccount.Service — 11 methods
- `SetEmployeeProfileProvider()`
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

### workforceintelligence.Service — 52 methods
- `SetInternalEligibilityProvider()`
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
- `GetRecruitmentAnalytics()`
- `GetQualityOfHire()`
- `computeTimeToFill()`
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
- `resolveEligibleInternalCandidates()`

## SECTION 6: GRAND TOTALS

| Category | Count |
|----------|:-----:|
| Tenant Modules | 20 |
| Platform Modules | 6 |
| Shared Kernel Packages | 25 |
| **Total Architecture Layers** | **51** |
| Total GORM Entities (tenant) | 221 |
| Total GORM Entities (platform) | 8 |
| **Total Entities (combined)** | **229** |
| Total Service Methods | 1265 |
| Total Repository Methods | 1357 |
| Total Handler Functions | 1100 |
| Total Route Registrations | 1142 |
| **Total Unit Tests (all)** | **2077** |
| Total Go Source Files | 290 |
| Total Test Files (_test.go) | 175 |
| **Total Go Files** | **465** |

## SECTION 7: TEST FILE INVENTORY

| File | Test Funcs |
|------|:----------:|
| `internal\modules\approval\errors_test.go` | 4 |
| `internal\modules\approval\handler_test.go` | 14 |
| `internal\modules\approval\helpers_test.go` | 0 |
| `internal\modules\approval\hierarchy_test.go` | 7 |
| `internal\modules\approval\module_subscription_test.go` | 20 |
| `internal\modules\approval\notifier_test.go` | 5 |
| `internal\modules\approval\repository_test.go` | 25 |
| `internal\modules\approval\routes_test.go` | 1 |
| `internal\modules\approval\service_test.go` | 35 |
| `internal\modules\approval\status_handler_test.go` | 8 |
| `internal\modules\attendance\approval_integration_test.go` | 14 |
| `internal\modules\attendance\correction_test.go` | 3 |
| `internal\modules\attendance\handler_test.go` | 26 |
| `internal\modules\attendance\helpers_test.go` | 0 |
| `internal\modules\attendance\notifier_integration_test.go` | 4 |
| `internal\modules\attendance\overtime_two_flow_test.go` | 15 |
| `internal\modules\attendance\repository_test.go` | 34 |
| `internal\modules\attendance\routes_test.go` | 2 |
| `internal\modules\attendance\service_test.go` | 39 |
| `internal\modules\attendance\session_test.go` | 8 |
| `internal\modules\attendance\summary_test.go` | 6 |
| `internal\modules\careerintelligence\handler_test.go` | 22 |
| `internal\modules\careerintelligence\helpers_test.go` | 0 |
| `internal\modules\careerintelligence\repository_test.go` | 31 |
| `internal\modules\careerintelligence\service_test.go` | 29 |
| `internal\modules\competency\handler_test.go` | 13 |
| `internal\modules\competency\helpers_test.go` | 0 |
| `internal\modules\competency\repository_test.go` | 12 |
| `internal\modules\competency\service_calculation_test.go` | 19 |
| `internal\modules\competency\service_test.go` | 37 |
| `internal\modules\documenttemplate\dateformat_test.go` | 1 |
| `internal\modules\documenttemplate\docx2pdf_service_test.go` | 1 |
| `internal\modules\documenttemplate\docx_test.go` | 10 |
| `internal\modules\documenttemplate\generator_test.go` | 4 |
| `internal\modules\documenttemplate\handler_test.go` | 10 |
| `internal\modules\documenttemplate\helpers_test.go` | 0 |
| `internal\modules\documenttemplate\model_test.go` | 2 |
| `internal\modules\documenttemplate\pdf_service_test.go` | 4 |
| `internal\modules\documenttemplate\repository_test.go` | 7 |
| `internal\modules\documenttemplate\service_test.go` | 10 |
| `internal\modules\documenttemplate\variables_test.go` | 2 |
| `internal\modules\employee\dto_test.go` | 7 |
| `internal\modules\employee\employee_id_format_test.go` | 4 |
| `internal\modules\employee\helpers_test.go` | 0 |
| `internal\modules\employee\repository_test.go` | 22 |
| `internal\modules\employee\sensitive_field_registry_test.go` | 3 |
| `internal\modules\employee\sensitive_field_settings_authz_test.go` | 2 |
| `internal\modules\employee\sensitive_field_settings_handler_test.go` | 3 |
| `internal\modules\employee\sensitive_field_settings_test.go` | 4 |
| `internal\modules\employee\sensitive_field_update_test.go` | 8 |
| `internal\modules\employee\service_test.go` | 37 |
| `internal\modules\employeemovement\approval_integration_test.go` | 10 |
| `internal\modules\employeemovement\audit_test.go` | 10 |
| `internal\modules\employeemovement\career_test.go` | 4 |
| `internal\modules\employeemovement\document_sensitive_test.go` | 2 |
| `internal\modules\employeemovement\document_test.go` | 10 |
| `internal\modules\employeemovement\eligibility_test.go` | 7 |
| `internal\modules\employeemovement\enrichment_test.go` | 6 |
| `internal\modules\employeemovement\expiry_test.go` | 3 |
| `internal\modules\employeemovement\handler_test.go` | 16 |
| `internal\modules\employeemovement\helpers_test.go` | 0 |
| `internal\modules\employeemovement\numbering_test.go` | 2 |
| `internal\modules\employeemovement\report_test.go` | 12 |
| `internal\modules\employeemovement\repository_test.go` | 27 |
| `internal\modules\employeemovement\service_test.go` | 49 |
| `internal\modules\jobmanagement\calculator_test.go` | 13 |
| `internal\modules\jobmanagement\education_update_test.go` | 1 |
| `internal\modules\jobmanagement\handler_test.go` | 10 |
| `internal\modules\jobmanagement\helpers_test.go` | 0 |
| `internal\modules\jobmanagement\repository_test.go` | 29 |
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
| `internal\modules\leave\report_test.go` | 3 |
| `internal\modules\leave\repository_test.go` | 17 |
| `internal\modules\leave\routes_test.go` | 1 |
| `internal\modules\leave\service_test.go` | 18 |
| `internal\modules\notification\helpers_test.go` | 0 |
| `internal\modules\notification\repository_test.go` | 2 |
| `internal\modules\notification\service_test.go` | 6 |
| `internal\modules\organization\helpers_test.go` | 0 |
| `internal\modules\organization\service_test.go` | 9 |
| `internal\modules\payroll\bpjs_test.go` | 5 |
| `internal\modules\payroll\calculation_test.go` | 8 |
| `internal\modules\payroll\calculator\engine_test.go` | 18 |
| `internal\modules\payroll\calculator\proration_test.go` | 7 |
| `internal\modules\payroll\helpers_test.go` | 0 |
| `internal\modules\payroll\payslip_payment_test.go` | 8 |
| `internal\modules\payroll\pph21_test.go` | 6 |
| `internal\modules\payroll\report_test.go` | 2 |
| `internal\modules\payroll\repository_test.go` | 23 |
| `internal\modules\payroll\service_test.go` | 37 |
| `internal\modules\payroll\workforce_test.go` | 6 |
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
| `internal\modules\rbac\service_test.go` | 2 |
| `internal\modules\recruitment\handler_test.go` | 57 |
| `internal\modules\recruitment\helpers_test.go` | 0 |
| `internal\modules\recruitment\repository_test.go` | 62 |
| `internal\modules\recruitment\service_test.go` | 153 |
| `internal\modules\reimbursement\approval_integration_test.go` | 9 |
| `internal\modules\reimbursement\handler_test.go` | 17 |
| `internal\modules\reimbursement\helpers_test.go` | 0 |
| `internal\modules\reimbursement\notifier_integration_test.go` | 6 |
| `internal\modules\reimbursement\repository_test.go` | 19 |
| `internal\modules\reimbursement\service_test.go` | 27 |
| `internal\modules\setting\competency_test.go` | 8 |
| `internal\modules\setting\employee_id_format_service_test.go` | 4 |
| `internal\modules\setting\helpers_test.go` | 0 |
| `internal\modules\setting\numbering_handler_test.go` | 2 |
| `internal\modules\setting\permission_test.go` | 2 |
| `internal\modules\setting\repository_test.go` | 18 |
| `internal\modules\setting\service_test.go` | 26 |
| `internal\modules\training\advanced_p2_test.go` | 7 |
| `internal\modules\training\handler_test.go` | 10 |
| `internal\modules\training\helpers_test.go` | 1 |
| `internal\modules\training\plan_p1_test.go` | 25 |
| `internal\modules\training\repository_test.go` | 14 |
| `internal\modules\training\routes_test.go` | 1 |
| `internal\modules\training\service_test.go` | 21 |
| `internal\modules\workforceintelligence\handler_test.go` | 36 |
| `internal\modules\workforceintelligence\helpers_test.go` | 0 |
| `internal\modules\workforceintelligence\repository_test.go` | 46 |
| `internal\modules\workforceintelligence\service_test.go` | 55 |
| `internal\platform\company\handler_test.go` | 8 |
| `internal\platform\company\helpers_test.go` | 0 |
| `internal\platform\company\service_test.go` | 22 |
| `internal\platform\modulemgmt\repository_test.go` | 3 |
| `internal\platform\package\helpers_test.go` | 0 |
| `internal\platform\package\service_test.go` | 27 |

====================================================================================================
  END OF REPORT
