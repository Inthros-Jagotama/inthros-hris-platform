import { createRouter, createWebHistory } from 'vue-router'
import { useActiveModules } from '@/stores/activeModules'
import { useNotifications } from '@/stores/notifications'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/set-password',
    name: 'SetPassword',
    component: () => import('@/views/SetPassword.vue'),
    meta: { guest: true }
  },
  {
    path: '/',
    component: () => import('@/layouts/AppLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/Profile.vue'),
        meta: { title: 'Profile', titleKey: 'profile.title', descKey: 'profile.description', icon: 'pi pi-user' }
      },
      // Company Profile (detail perusahaan)
      {
        path: 'company',
        name: 'CompanyDetail',
        component: () => import('@/views/CompanyDetail.vue'),
        meta: { title: 'Company Profile', titleKey: 'company_detail.title', descKey: 'company_detail.description', icon: 'pi pi-building' }
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: 'Dashboard', titleKey: 'dashboard.title', descKey: 'dashboard.description', icon: 'pi pi-home' }
      },
      // Organization Summary (masuk via menu Organization)
      {
        path: 'organization-summary',
        name: 'OrganizationSummary',
        component: () => import('@/views/modules/organization/OrganizationSummary.vue'),
        meta: { title: 'Organization', titleKey: 'organization.title', descKey: 'org_summary.description', icon: 'pi pi-building', module: 'organization' }
      },
      // Organization Management
      {
        path: 'organizations',
        name: 'Organizations',
        component: () => import('@/views/modules/organization/Organizations.vue'),
        meta: { title: 'Organization', titleKey: 'organization.title', descKey: 'organization.description', icon: 'pi pi-sitemap', module: 'organization' }
      },
      // Employee Management
      {
        path: 'employees',
        name: 'Employees',
        component: () => import('@/views/modules/employee/Employees.vue'),
        meta: { title: 'Employees', titleKey: 'employee.title', descKey: 'employee.description', icon: 'pi pi-users', module: 'employee' }
      },
      {
        path: 'employees/new',
        name: 'EmployeeNew',
        component: () => import('@/views/modules/employee/EmployeeForm.vue'),
        meta: { title: 'New Employee', titleKey: 'employee.new', descKey: 'employee.description', icon: 'pi pi-user-plus', module: 'employee', backRoute: '/employees', backLabelKey: 'nav.employees' }
      },
      {
        path: 'employees/:id',
        name: 'EmployeeDetail',
        component: () => import('@/views/modules/employee/EmployeeDetail.vue'),
        meta: { title: 'View Employee', titleKey: 'employee.view', descKey: 'employee.description', icon: 'pi pi-eye', module: 'employee', backRoute: '/employees', backLabelKey: 'nav.employees' }
      },
      {
        path: 'employees/:id/edit',
        name: 'EmployeeEdit',
        component: () => import('@/views/modules/employee/EmployeeForm.vue'),
        meta: { title: 'Edit Employee', titleKey: 'employee.edit', descKey: 'employee.description', icon: 'pi pi-user-edit', module: 'employee', backRoute: '/employees', backLabelKey: 'nav.employees' }
      },
      // Job Management
      {
        path: 'job-management',
        name: 'JobManagement',
        component: () => import('@/views/modules/job/JobManagement.vue'),
        meta: { title: 'Job Management', titleKey: 'job_management.title', descKey: 'job_management.page_description', icon: 'pi pi-briefcase', module: 'jobmanagement' }
      },
      {
        path: 'job-management/values',
        name: 'JobValues',
        component: () => import('@/views/modules/jobvalues/JobValuesIndex.vue'),
        meta: { title: 'Job Values', titleKey: 'job_management.values', descKey: 'job_values.description', icon: 'pi pi-sliders-h', module: 'jobmanagement' }
      },
      {
        path: 'job-management/values/:type',
        name: 'JobValuesType',
        component: () => import('@/views/modules/jobvalues/JobValuesForm.vue'),
        meta: { title: 'Job Value Type', titleKey: 'job_management.values', descKey: 'job_values.description', icon: 'pi pi-sliders-h', module: 'jobmanagement' }
      },
      {
        path: 'job-management/form',
        name: 'JobManagementForm',
        component: () => import('@/views/modules/job/JobManagementForm.vue'),
        meta: { title: 'Manage Job Data', titleKey: 'job_management.manage', descKey: 'job_management.manage_description', icon: 'pi pi-briefcase', module: 'jobmanagement', backRoute: '/job-management', backLabelKey: 'nav.job_management' }
      },
      // Competency Management
      {
        path: 'competencies',
        name: 'Competencies',
        component: () => import('@/views/modules/competency/Competencies.vue'),
        meta: { title: 'Competency', titleKey: 'competency.title', descKey: 'competency.description', icon: 'pi pi-star', module: 'competency' }
      },
      {
        path: 'competencies/values',
        name: 'CompetencyValues',
        component: () => import('@/views/modules/competency/CompetencyValues.vue'),
        meta: { title: 'Rating Scales', titleKey: 'competency_360.rating_scales', descKey: 'competency_360.rating_scales_desc', icon: 'pi pi-sliders-h', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/templates',
        name: 'AssessmentTemplates',
        component: () => import('@/views/modules/competency/AssessmentTemplates.vue'),
        meta: { title: 'Assessment Templates', titleKey: 'competency_360.templates', descKey: 'competency_360.templates_desc', icon: 'pi pi-clone', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/events',
        name: 'CompetencyEvents',
        component: () => import('@/views/modules/competency/CompetencyEvents.vue'),
        meta: { title: 'Competency Events', titleKey: 'competency_360.events', descKey: 'competency_360.events_desc', icon: 'pi pi-calendar', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/raters',
        name: 'RaterAssignment',
        component: () => import('@/views/modules/competency/RaterAssignment.vue'),
        meta: { title: 'Rater Assignment', titleKey: 'competency_360.rater_assignment', descKey: 'competency_360.rater_assignment_desc', icon: 'pi pi-users', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/my-assessments',
        name: 'MyAssessments',
        component: () => import('@/views/modules/competency/MyAssessments.vue'),
        meta: { title: 'My Assessment', titleKey: 'competency_360.my_assessments', descKey: 'competency_360.my_assessments_desc', icon: 'pi pi-clipboard-list', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/results',
        name: 'AssessmentResult',
        component: () => import('@/views/modules/competency/AssessmentResult.vue'),
        meta: { title: 'Assessment Result', titleKey: 'competency_360.results', descKey: 'competency_360.results_desc', icon: 'pi pi-chart-bar', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      {
        path: 'competencies/reports',
        name: 'CompetencyReports',
        component: () => import('@/views/modules/competency/CompetencyReports.vue'),
        meta: { title: 'Competency Reports', titleKey: 'competency_360.reports', descKey: 'competency_360.reports_desc', icon: 'pi pi-file-chart', module: 'competency', backRoute: '/competencies', backLabelKey: 'nav.competency' }
      },
      // Employee Movement & Career Management (G-8: route menyamakan menu server)
      {
        path: 'admin/career/movements',
        name: 'EmployeeMovements',
        component: () => import('@/views/modules/employeemovement/EmployeeMovements.vue'),
        meta: { title: 'Movements', titleKey: 'employee_movement.movements', descKey: 'employee_movement.description', icon: 'pi pi-arrows-alt', module: 'employeemovement', backRoute: '/admin/career/reports', backLabelKey: 'employee_movement.movements_contracts' }
      },
      {
        path: 'admin/career/contracts',
        name: 'EmployeeContracts',
        component: () => import('@/views/modules/employeemovement/EmployeeContracts.vue'),
        meta: { title: 'Contracts', titleKey: 'employee_movement.contracts', descKey: 'employee_movement.description', icon: 'pi pi-file-edit', module: 'employeemovement', backRoute: '/admin/career/reports', backLabelKey: 'employee_movement.movements_contracts' }
      },
      {
        path: 'admin/career/reports',
        name: 'EmployeeMovementReports',
        component: () => import('@/views/modules/employeemovement/EmployeeMovementReports.vue'),
        meta: { title: 'Movements & Contracts', titleKey: 'employee_movement.movements_contracts', descKey: 'employee_movement.reports_description', icon: 'pi pi-chart-bar', module: 'employeemovement' }
      },
      {
        path: 'career-intelligence/paths',
        name: 'CareerPaths',
        component: () => import('@/views/modules/career-intelligence/CareerPaths.vue'),
        meta: { title: 'Career Paths', titleKey: 'career_paths.title', descKey: 'career_paths.description', icon: 'pi pi-sitemap', module: 'career-intelligence', backRoute: '/career-intelligence', backLabelKey: 'career_intel.title' }
      },
      // Time & Attendance
      {
        path: 'attendance',
        name: 'Attendance',
        component: () => import('@/views/modules/attendance/Attendance.vue'),
        meta: { title: 'Attendance', titleKey: 'attendance.title', descKey: 'attendance.description', icon: 'pi pi-clock', module: 'attendance' }
      },
      {
        path: 'attendance/admin',
        name: 'AttendanceAdmin',
        component: () => import('@/views/modules/attendance/AttendanceAdmin.vue'),
        meta: { title: 'Attendance Admin', titleKey: 'attendance.admin', descKey: 'attendance.admin_description', icon: 'pi pi-cog', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/settings',
        name: 'AttendanceSettings',
        component: () => import('@/views/modules/attendance/AttendanceSettings.vue'),
        meta: { title: 'Attendance Settings', titleKey: 'attendance.settings', descKey: 'attendance.settings_description', icon: 'pi pi-cog', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/shifts',
        name: 'AttendanceShifts',
        component: () => import('@/views/modules/attendance/AttendanceShifts.vue'),
        meta: { title: 'Shifts', titleKey: 'attendance.shifts', descKey: 'attendance.shifts_description', icon: 'pi pi-clock', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/employee-shifts',
        name: 'AttendanceEmployeeShifts',
        component: () => import('@/views/modules/attendance/AttendanceEmployeeShifts.vue'),
        meta: { title: 'Employee Shifts', titleKey: 'attendance.employee_shifts', descKey: 'attendance.employee_shifts_description', icon: 'pi pi-users', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/locations',
        name: 'AttendanceLocations',
        component: () => import('@/views/modules/attendance/AttendanceLocations.vue'),
        meta: { title: 'Locations', titleKey: 'attendance.locations', descKey: 'attendance.locations_description', icon: 'pi pi-map-marker', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/exempt-positions',
        name: 'AttendanceExemptPositions',
        component: () => import('@/views/modules/attendance/AttendanceExemptPositions.vue'),
        meta: { title: 'Exempt Positions', titleKey: 'attendance.exempt_positions', descKey: 'attendance.exempt_positions_description', icon: 'pi pi-shield', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/overtime',
        name: 'AttendanceOvertime',
        component: () => import('@/views/modules/attendance/AttendanceOvertime.vue'),
        meta: { title: 'Overtime', titleKey: 'attendance.overtime', descKey: 'attendance.overtime_description', icon: 'pi pi-clock', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/corrections',
        name: 'AttendanceCorrections',
        component: () => import('@/views/modules/attendance/AttendanceCorrections.vue'),
        meta: { title: 'Corrections', titleKey: 'attendance.corrections', descKey: 'attendance.corrections_description', icon: 'pi pi-pencil', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/business-travel',
        name: 'BusinessTravelList',
        component: () => import('@/views/modules/attendance/business-travel/BusinessTravelList.vue'),
        meta: { title: 'Business Travel', titleKey: 'business_travel.title', descKey: 'business_travel.description', icon: 'pi pi-briefcase', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/business-travel/:id',
        name: 'BusinessTravelDetail',
        component: () => import('@/views/modules/attendance/business-travel/BusinessTravelDetail.vue'),
        meta: { title: 'Business Travel', titleKey: 'business_travel.view', descKey: 'business_travel.description', icon: 'pi pi-briefcase', module: 'attendance', backRoute: '/attendance/business-travel', backLabelKey: 'business_travel.title' }
      },
      {
        path: 'attendance/events',
        name: 'AttendanceEvents',
        component: () => import('@/views/modules/attendance/AttendanceEvents.vue'),
        meta: { title: 'Events', titleKey: 'attendance.events', descKey: 'attendance.events_description', icon: 'pi pi-list', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/sessions',
        name: 'AttendanceSessions',
        component: () => import('@/views/modules/attendance/AttendanceSessions.vue'),
        meta: { title: 'Sessions', titleKey: 'attendance.sessions', descKey: 'attendance.sessions_description', icon: 'pi pi-calendar', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/reports',
        name: 'AttendanceReports',
        component: () => import('@/views/modules/attendance/AttendanceReports.vue'),
        meta: { title: 'Reports', titleKey: 'attendance.reports', descKey: 'attendance.reports_description', icon: 'pi pi-chart-bar', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      // Approval Engine
      {
        path: 'approvals',
        name: 'Approvals',
        component: () => import('@/views/modules/approval/Approvals.vue'),
        meta: { title: 'Approval', titleKey: 'approval.title', descKey: 'approval.description', icon: 'pi pi-check-square', module: 'approval' }
      },
      {
        path: 'approvals/flows',
        name: 'ApprovalFlows',
        component: () => import('@/views/modules/approval/ApprovalFlows.vue'),
        meta: { title: 'Approval Flows', titleKey: 'approval.flows', descKey: 'approval.description', icon: 'pi pi-sitemap', module: 'approval' }
      },
      // Notifications
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/views/modules/notification/Notifications.vue'),
        meta: { title: 'Notifications', titleKey: 'notification.title', descKey: 'notification.description', icon: 'pi pi-bell', module: 'notification' }
      },
      // Payroll
      {
        path: 'payroll',
        name: 'Payroll',
        component: () => import('@/views/modules/payroll/Payroll.vue'),
        meta: { title: 'Payroll', titleKey: 'payroll.title', descKey: 'payroll.description', icon: 'pi pi-dollar', module: 'payroll' }
      },
      {
        path: 'payroll/runs/:id',
        name: 'PayrollRunDetail',
        component: () => import('@/views/modules/payroll/PayrollRunDetail.vue'),
        meta: { title: 'Run Detail', titleKey: 'payroll.run_detail', descKey: 'payroll.description', icon: 'pi pi-dollar', module: 'payroll', backRoute: '/payroll', backLabelKey: 'payroll.payroll_runs' }
      },
      // Leave & Time Off
      {
        path: 'leave',
        name: 'Leave',
        component: () => import('@/views/modules/leave/Leave.vue'),
        meta: { title: 'Leave', titleKey: 'leave.title', descKey: 'leave.description', icon: 'pi pi-calendar', module: 'leave' }
      },
      {
        path: 'leave/admin',
        name: 'LeaveAdmin',
        component: () => import('@/views/modules/leave/LeaveAdmin.vue'),
        meta: { title: 'Leave Admin', titleKey: 'leave.admin', descKey: 'leave.admin_description', icon: 'pi pi-cog', module: 'leave', backRoute: '/leave', backLabelKey: 'leave.title' }
      },
      {
        path: 'leave/types',
        name: 'LeaveTypes',
        component: () => import('@/views/modules/leave/LeaveTypes.vue'),
        meta: { title: 'Leave Types', titleKey: 'leave.types', descKey: 'leave.types_description', icon: 'pi pi-tags', module: 'leave', backRoute: '/leave/admin', backLabelKey: 'leave.admin' }
      },
      {
        path: 'leave/accrual-policies',
        name: 'LeaveAccrualPolicies',
        component: () => import('@/views/modules/leave/LeaveAccrualPolicies.vue'),
        meta: { title: 'Accrual Policies', titleKey: 'leave.accrual_policies', descKey: 'leave.accrual_policies_description', icon: 'pi pi-percentage', module: 'leave', backRoute: '/leave/admin', backLabelKey: 'leave.admin' }
      },
      {
        path: 'leave/reasons',
        name: 'LeaveReasons',
        component: () => import('@/views/modules/leave/LeaveReasons.vue'),
        meta: { title: 'Leave Reasons', titleKey: 'leave.reasons', descKey: 'leave.reasons_description', icon: 'pi pi-list', module: 'leave', backRoute: '/leave/admin', backLabelKey: 'leave.admin' }
      },
      // Performance Management
      {
        path: 'performance',
        name: 'Performance',
        component: () => import('@/views/modules/performance/PerformanceIndex.vue'),
        meta: { title: 'Performance', titleKey: 'performance.title', descKey: 'performance.description', icon: 'pi pi-chart-line', module: 'performance' }
      },
      {
        path: 'performance/kpi',
        name: 'PerformanceKPI',
        component: () => import('@/views/modules/performance/kpi/KPIIndex.vue'),
        meta: { title: 'KPI', titleKey: 'kpi.title', descKey: 'kpi.description', icon: 'pi pi-chart-bar', module: 'performance' }
      },
      {
        path: 'performance/kpi/templates',
        name: 'KPITemplates',
        component: () => import('@/views/modules/performance/kpi/KPITemplates.vue'),
        meta: { title: 'KPI Templates', titleKey: 'kpi.templates', descKey: 'kpi.templates_desc', icon: 'pi pi-file', module: 'performance', backRoute: '/performance', backLabelKey: 'performance.title' }
      },
      {
        path: 'performance/kpi/templates/new',
        name: 'KPITemplateNew',
        component: () => import('@/views/modules/performance/kpi/KPITemplateForm.vue'),
        meta: { title: 'New Template', titleKey: 'kpi.template_new', descKey: 'kpi.template_desc', icon: 'pi pi-plus', module: 'performance', backRoute: '/performance/kpi/templates', backLabelKey: 'kpi.templates' }
      },
      {
        path: 'performance/kpi/templates/:id/edit',
        name: 'KPITemplateEdit',
        component: () => import('@/views/modules/performance/kpi/KPITemplateForm.vue'),
        meta: { title: 'Edit Template', titleKey: 'kpi.template_edit', descKey: 'kpi.template_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/kpi/templates', backLabelKey: 'kpi.templates' }
      },
      {
        path: 'performance/kpi/periods',
        name: 'KPIPeriods',
        component: () => import('@/views/settings/PerformancePeriodsView.vue'),
        meta: { title: 'Periods', titleKey: 'performance_periods.title', descKey: 'performance_periods.description', icon: 'pi pi-calendar', module: 'performance', backRoute: '/performance', backLabelKey: 'performance.title' }
      },
      {
        path: 'performance/kpi/my-evaluation',
        name: 'KPISelfAssessment',
        component: () => import('@/views/modules/performance/kpi/KPISelfAssessment.vue'),
        meta: { title: 'My Evaluation', titleKey: 'kpi.my_evaluation', descKey: 'kpi.my_evaluation_desc', icon: 'pi pi-user-edit', module: 'performance', backRoute: '/performance', backLabelKey: 'performance.title' }
      },
      {
        path: 'performance/kpi/evaluation/:id',
        name: 'KPIEvaluationDetail',
        component: () => import('@/views/modules/performance/kpi/KPIEvaluationDetail.vue'),
        meta: { title: 'Evaluation Detail', titleKey: 'kpi.evaluation_detail', descKey: 'kpi.evaluation_detail_desc', icon: 'pi pi-chart-bar', module: 'performance', backRoute: '/performance/kpi', backLabelKey: 'kpi.evaluations' }
      },
      {
        path: 'performance/okr/my-evaluation',
        name: 'OKRSelfAssessment',
        component: () => import('@/views/modules/performance/okr/OKRSelfAssessment.vue'),
        meta: { title: 'My OKR', titleKey: 'okr.my_evaluation', descKey: 'okr.my_evaluation_desc', icon: 'pi pi-user-edit', module: 'performance' }
      },
      {
        path: 'performance/okr',
        name: 'PerformanceOKR',
        component: () => import('@/views/modules/performance/okr/OKRIndex.vue'),
        meta: { title: 'OKR', titleKey: 'okr.title', descKey: 'okr.description', icon: 'pi pi-bullseye', module: 'performance' }
      },
      {
        path: 'performance/okr/templates',
        name: 'OKRTemplates',
        component: () => import('@/views/modules/performance/okr/OKRTemplates.vue'),
        meta: { title: 'OKR Templates', titleKey: 'okr.templates', descKey: 'okr.templates_desc', icon: 'pi pi-file', module: 'performance', backRoute: '/performance', backLabelKey: 'performance.title' }
      },
      {
        path: 'performance/okr/templates/new',
        name: 'OKRTemplateNew',
        component: () => import('@/views/modules/performance/okr/OKRTemplateForm.vue'),
        meta: { title: 'New Template', titleKey: 'okr.template_new', descKey: 'okr.template_desc', icon: 'pi pi-plus', module: 'performance', backRoute: '/performance/okr/templates', backLabelKey: 'okr.templates' }
      },
      {
        path: 'performance/okr/templates/:id/edit',
        name: 'OKRTemplateEdit',
        component: () => import('@/views/modules/performance/okr/OKRTemplateForm.vue'),
        meta: { title: 'Edit Template', titleKey: 'okr.template_edit', descKey: 'okr.template_desc', icon: 'pi pi-pencil', module: 'performance', backRoute: '/performance/okr/templates', backLabelKey: 'okr.templates' }
      },
      {
        path: 'performance/okr/evaluation/:id',
        name: 'OKREvaluationDetail',
        component: () => import('@/views/modules/performance/okr/OKREvaluationDetail.vue'),
        meta: { title: 'Evaluation Detail', titleKey: 'okr.evaluation_detail', descKey: 'okr.evaluation_detail_desc', icon: 'pi pi-bullseye', module: 'performance', backRoute: '/performance/okr', backLabelKey: 'okr.evaluations' }
      },
      // Recruitment & Onboarding
      {
        path: 'recruitment',
        name: 'Recruitment',
        component: () => import('@/views/modules/recruitment/Recruitment.vue'),
        meta: { title: 'Recruitment', titleKey: 'recruitment.title', descKey: 'recruitment.description', icon: 'pi pi-user-plus', module: 'recruitment' }
      },
      // Job Requisitions — S-1 (workforce gap) / S-5 (succession gap) reason_type
      {
        path: 'recruitment/requisitions',
        name: 'Requisitions',
        component: () => import('@/views/modules/recruitment/Requisitions.vue'),
        meta: { title: 'Job Requisitions', titleKey: 'recruitment.requisitions', descKey: 'requisitions.description', icon: 'pi pi-briefcase', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      // Requisition Requirements & Competencies — G-9 sub-1 (+ Job Management fallback, G-9 sub-3)
      {
        path: 'recruitment/requisitions/:id/requirements',
        name: 'RequisitionRequirements',
        component: () => import('@/views/modules/recruitment/RequisitionRequirements.vue'),
        meta: { title: 'Requirements & Competencies', titleKey: 'requisitions.requirements_competencies', module: 'recruitment', backRoute: '/recruitment/requisitions', backLabelKey: 'recruitment.requisitions' }
      },
      // Applications — pipeline (list + detail: history/screening/assessment/interviews/match score)
      {
        path: 'recruitment/applications',
        name: 'Applications',
        component: () => import('@/views/modules/recruitment/Applications.vue'),
        meta: { title: 'Applications', titleKey: 'recruitment.applications', descKey: 'applications.description', icon: 'pi pi-send', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      {
        path: 'recruitment/applications/:id',
        name: 'ApplicationDetail',
        component: () => import('@/views/modules/recruitment/ApplicationDetail.vue'),
        meta: { title: 'Application Detail', titleKey: 'applications.detail', module: 'recruitment', backRoute: '/recruitment/applications', backLabelKey: 'recruitment.applications' }
      },
      // Candidates — G-6 profile (educations/experiences/skills/certifications/documents/consents)
      {
        path: 'recruitment/candidates',
        name: 'Candidates',
        component: () => import('@/views/modules/recruitment/Candidates.vue'),
        meta: { title: 'Candidates', titleKey: 'recruitment.candidates', descKey: 'candidates.description', icon: 'pi pi-users', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      {
        path: 'recruitment/candidates/:id',
        name: 'CandidateDetail',
        component: () => import('@/views/modules/recruitment/CandidateDetail.vue'),
        meta: { title: 'Candidate Profile', titleKey: 'candidates.profile', module: 'recruitment', backRoute: '/recruitment/candidates', backLabelKey: 'recruitment.candidates' }
      },
      // Internal Candidates — S-4 (eligible via career path dari Career Intelligence)
      {
        path: 'recruitment/internal-candidates',
        name: 'InternalCandidates',
        component: () => import('@/views/modules/recruitment/InternalCandidates.vue'),
        meta: { title: 'Internal Candidates', titleKey: 'internal_candidates.title', descKey: 'internal_candidates.description', icon: 'pi pi-user-plus', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      // Job Offers — G-3 (offer management + approval workflow)
      {
        path: 'recruitment/offers',
        name: 'Offers',
        component: () => import('@/views/modules/recruitment/Offers.vue'),
        meta: { title: 'Job Offers', titleKey: 'recruitment.offers', descKey: 'offers.description', icon: 'pi pi-file-edit', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      // Assessments — G-7 sub-project 2 (batch session + peserta)
      {
        path: 'recruitment/assessments',
        name: 'Assessments',
        component: () => import('@/views/modules/recruitment/Assessments.vue'),
        meta: { title: 'Assessments', titleKey: 'recruitment.assessments', descKey: 'assessments.description', icon: 'pi pi-clipboard', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      // Onboarding — G-4: employee hasil offer (recruited_from_application_id) + status (S-7 handoff)
      {
        path: 'recruitment/onboarding',
        name: 'Onboarding',
        component: () => import('@/views/modules/recruitment/Onboarding.vue'),
        meta: { title: 'Onboarding', titleKey: 'recruitment.onboarding', descKey: 'onboarding.description', icon: 'pi pi-rocket', module: 'recruitment', backRoute: '/recruitment', backLabelKey: 'recruitment.title' }
      },
      // Reimbursement — index page (menu cards + summary)
      {
        path: 'reimbursements',
        name: 'Reimbursements',
        component: () => import('@/views/modules/reimbursement/Reimbursements.vue'),
        meta: { title: 'Reimbursement', titleKey: 'reimbursement.title', descKey: 'reimbursement.description', icon: 'pi pi-credit-card', module: 'reimbursement' }
      },
      {
        path: 'reimbursements/all',
        name: 'ReimbursementAllRequests',
        component: () => import('@/views/modules/reimbursement/ReimbursementRequests.vue'),
        meta: { title: 'Reimbursement', titleKey: 'reimbursement.requests', descKey: 'reimbursement.description', icon: 'pi pi-briefcase', module: 'reimbursement', backRoute: '/reimbursements', backLabelKey: 'reimbursement.title' }
      },
      {
        path: 'reimbursements/my-requests',
        name: 'ReimbursementMyRequests',
        component: () => import('@/views/modules/reimbursement/ReimbursementRequests.vue'),
        meta: { title: 'My Reimbursement Requests', titleKey: 'reimbursement.my_requests', descKey: 'reimbursement.description', icon: 'pi pi-credit-card', module: 'reimbursement', backRoute: '/reimbursements', backLabelKey: 'reimbursement.title' }
      },
      {
        path: 'reimbursements/types',
        name: 'ReimbursementTypes',
        component: () => import('@/views/modules/reimbursement/ReimbursementTypes.vue'),
        meta: { title: 'Reimbursement Types', titleKey: 'reimbursement.types', descKey: 'reimbursement.types_description', icon: 'pi pi-tags', module: 'reimbursement', backRoute: '/reimbursements', backLabelKey: 'reimbursement.title' }
      },
      {
        path: 'reimbursements/:id',
        name: 'ReimbursementRequestDetail',
        component: () => import('@/views/modules/reimbursement/ReimbursementRequestDetail.vue'),
        meta: { title: 'View Reimbursement Request', titleKey: 'reimbursement.view', descKey: 'reimbursement.description', icon: 'pi pi-eye', module: 'reimbursement', backRoute: '/reimbursements', backLabelKey: 'reimbursement.title' }
      },
      // Training & Development
      {
        path: 'training',
        name: 'Training',
        component: () => import('@/views/modules/training/Training.vue'),
        meta: { title: 'Training', titleKey: 'training.title', descKey: 'training.description', icon: 'pi pi-book', module: 'training' }
      },
      {
        path: 'training/courses',
        name: 'TrainingCourses',
        component: () => import('@/views/modules/training/TrainingCourses.vue'),
        meta: { title: 'Training Courses', titleKey: 'training.courses', descKey: 'training.courses_desc', icon: 'pi pi-book', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/categories',
        name: 'TrainingCategories',
        component: () => import('@/views/modules/training/TrainingCategories.vue'),
        meta: { title: 'Training Categories', titleKey: 'training.categories', descKey: 'training.categories_desc', icon: 'pi pi-tags', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/providers',
        name: 'TrainingProviders',
        component: () => import('@/views/modules/training/TrainingProviders.vue'),
        meta: { title: 'Training Providers', titleKey: 'training.providers', descKey: 'training.providers_desc', icon: 'pi pi-building', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/trainers',
        name: 'TrainingTrainers',
        component: () => import('@/views/modules/training/TrainingTrainers.vue'),
        meta: { title: 'Training Trainers', titleKey: 'training.trainers', descKey: 'training.trainers_desc', icon: 'pi pi-user', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/sessions',
        name: 'TrainingSessions',
        component: () => import('@/views/modules/training/TrainingSessions.vue'),
        meta: { title: 'Training Sessions', titleKey: 'training.sessions', descKey: 'training.sessions_desc', icon: 'pi pi-calendar', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/sessions/:id',
        name: 'TrainingSessionDetail',
        component: () => import('@/views/modules/training/TrainingSessionDetail.vue'),
        meta: { title: 'Session Detail', titleKey: 'training.session_detail', descKey: 'training.sessions_desc', icon: 'pi pi-calendar', module: 'training', backRoute: '/training/sessions', backLabelKey: 'training.sessions' }
      },
      {
        path: 'training/participants',
        name: 'TrainingParticipants',
        component: () => import('@/views/modules/training/TrainingParticipants.vue'),
        meta: { title: 'Training Participants', titleKey: 'training.participants', descKey: 'training.participants_desc', icon: 'pi pi-users', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/plans',
        name: 'TrainingPlans',
        component: () => import('@/views/modules/training/TrainingPlans.vue'),
        meta: { title: 'Training Planning', titleKey: 'training.planning', descKey: 'training.planning_desc', icon: 'pi pi-calendar-plus', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/requests',
        name: 'TrainingRequests',
        component: () => import('@/views/modules/training/TrainingRequests.vue'),
        meta: { title: 'Training Requests', titleKey: 'training.requests', descKey: 'training.requests_desc', icon: 'pi pi-send', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/needs',
        name: 'TrainingNeeds',
        component: () => import('@/views/modules/training/TrainingNeeds.vue'),
        meta: { title: 'Training Needs', titleKey: 'training.needs', descKey: 'training.needs_desc', icon: 'pi pi-bullseye', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/certificates',
        name: 'TrainingCertificates',
        component: () => import('@/views/modules/training/TrainingCertificates.vue'),
        meta: { title: 'Training Certificates', titleKey: 'training.certificates', descKey: 'training.certificates_desc', icon: 'pi pi-id-card', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/history',
        name: 'TrainingHistory',
        component: () => import('@/views/modules/training/TrainingHistory.vue'),
        meta: { title: 'Training History', titleKey: 'training.history', descKey: 'training.history_desc', icon: 'pi pi-history', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      {
        path: 'training/reports',
        name: 'TrainingReports',
        component: () => import('@/views/modules/training/TrainingReports.vue'),
        meta: { title: 'Training Reports', titleKey: 'training.reports', descKey: 'training.reports_desc', icon: 'pi pi-chart-bar', module: 'training', backRoute: '/training', backLabelKey: 'training.title' }
      },
      // Workforce Intelligence
      {
        path: 'workforce-intelligence',
        name: 'WorkforceIntelligence',
        component: () => import('@/views/modules/workforce-intelligence/WorkforceIntelligence.vue'),
        meta: { title: 'Workforce Intel', titleKey: 'workforce_intel.title', descKey: 'workforce_intel.description', icon: 'pi pi-chart-bar', module: 'workforce-intelligence' }
      },
      // Candidate Search — posisi kosong + kandidat recruitment
      {
        path: 'workforce-intelligence/candidate-search',
        name: 'CandidateSearch',
        component: () => import('@/views/modules/workforce-intelligence/CandidateSearch.vue'),
        meta: { title: 'Candidate Search', titleKey: 'candidate_search.title', descKey: 'candidate_search.description', icon: 'pi pi-user-plus', module: 'workforce-intelligence', backRoute: '/workforce-intelligence', backLabelKey: 'workforce_intel.title' }
      },
      // Quality of Hire — metrik agregat kualitas hire (S-6)
      {
        path: 'workforce-intelligence/quality-of-hire',
        name: 'QualityOfHire',
        component: () => import('@/views/modules/workforce-intelligence/QualityOfHire.vue'),
        meta: { title: 'Quality of Hire', titleKey: 'quality_of_hire.title', descKey: 'quality_of_hire.description', icon: 'pi pi-bullseye', module: 'workforce-intelligence', backRoute: '/workforce-intelligence', backLabelKey: 'workforce_intel.title' }
      },
      // Recruitment Analytics — S-2 (remaining gap, expected hires) + S-3 (time to hire/fill, OAR, source conversion)
      {
        path: 'workforce-intelligence/recruitment-analytics',
        name: 'RecruitmentAnalytics',
        component: () => import('@/views/modules/workforce-intelligence/RecruitmentAnalytics.vue'),
        meta: { title: 'Recruitment Analytics', titleKey: 'recruitment_analytics.title', descKey: 'recruitment_analytics.description', icon: 'pi pi-chart-line', module: 'workforce-intelligence', backRoute: '/workforce-intelligence', backLabelKey: 'workforce_intel.title' }
      },
      // Career Intelligence
      {
        path: 'career-intelligence',
        name: 'CareerIntelligence',
        component: () => import('@/views/modules/career-intelligence/CareerIntelligence.vue'),
        meta: { title: 'Career Intel', titleKey: 'career_intel.title', descKey: 'career_intel.description', icon: 'pi pi-chart-line', module: 'career-intelligence' }
      },
      // Succession Gaps — S-5 (posisi kunci tanpa successor siap → external recruitment)
      {
        path: 'career-intelligence/successions',
        name: 'SuccessionGaps',
        component: () => import('@/views/modules/career-intelligence/SuccessionGaps.vue'),
        meta: { title: 'Succession Gaps', titleKey: 'succession_gaps.title', descKey: 'succession_gaps.description', icon: 'pi pi-users', module: 'career-intelligence', backRoute: '/career-intelligence', backLabelKey: 'career_intel.title' }
      },
      // Settings — index page menampilkan card sub-menu settings
      { path: 'settings', name: 'SettingsIndex', component: () => import('@/views/settings/SettingsIndex.vue'), meta: { title: 'Settings', titleKey: 'settings.title', descKey: 'settings.description', icon: 'pi pi-cog', module: 'setting' } },
      { path: 'settings/zones', name: 'SettingsZones', component: () => import('@/views/settings/ZonesView.vue'), meta: { title: 'Zones', titleKey: 'zones.title', descKey: 'zones.description', icon: 'pi pi-map-marker', module: 'setting' } },
      { path: 'settings/provinces', name: 'SettingsProvinces', component: () => import('@/views/settings/ProvincesView.vue'), meta: { title: 'Provinces', titleKey: 'provinces.title', descKey: 'provinces.description', icon: 'pi pi-globe', module: 'setting' } },
      { path: 'settings/regencies', name: 'SettingsRegencies', component: () => import('@/views/settings/RegenciesView.vue'), meta: { title: 'Regencies', titleKey: 'regencies.title', descKey: 'regencies.description', icon: 'pi pi-map', module: 'setting' } },
      { path: 'settings/districts', name: 'SettingsDistricts', component: () => import('@/views/settings/DistrictsView.vue'), meta: { title: 'Districts', titleKey: 'districts.title', descKey: 'districts.description', icon: 'pi pi-building', module: 'setting' } },
      { path: 'settings/villages', name: 'SettingsVillages', component: () => import('@/views/settings/VillagesView.vue'), meta: { title: 'Villages', titleKey: 'villages.title', descKey: 'villages.description', icon: 'pi pi-home', module: 'setting' } },
      { path: 'settings/educations', name: 'SettingsEducations', component: () => import('@/views/settings/EducationsView.vue'), meta: { title: 'Educations', titleKey: 'educations.title', descKey: 'educations.description', icon: 'pi pi-graduation-cap', module: 'setting' } },
      { path: 'settings/education-majors', name: 'SettingsEducationMajors', component: () => import('@/views/settings/EducationMajorsView.vue'), meta: { title: 'Education Majors', titleKey: 'education_majors.title', descKey: 'education_majors.description', icon: 'pi pi-graduation-cap', module: 'setting' } },
      { path: 'settings/religions', name: 'SettingsReligions', component: () => import('@/views/settings/ReligionsView.vue'), meta: { title: 'Religions', titleKey: 'religions.title', descKey: 'religions.description', icon: 'pi pi-globe', module: 'setting' } },
      { path: 'settings/marital-statuses', name: 'SettingsMaritalStatuses', component: () => import('@/views/settings/MaritalStatusesView.vue'), meta: { title: 'Marital Statuses', titleKey: 'marital_status.title', descKey: 'marital_status.description', icon: 'pi pi-heart', module: 'setting' } },
      { path: 'settings/relationship-types', name: 'SettingsRelationshipTypes', component: () => import('@/views/settings/RelationshipTypesView.vue'), meta: { title: 'Relationship Types', titleKey: 'relationship_types.title', descKey: 'relationship_types.description', icon: 'pi pi-users', module: 'setting' } },
      { path: 'settings/banks', name: 'SettingsBanks', component: () => import('@/views/settings/BanksView.vue'), meta: { title: 'Banks', titleKey: 'banks.title', descKey: 'banks.description', icon: 'pi pi-home', module: 'setting' } },
      { path: 'settings/employment-statuses', name: 'SettingsEmploymentStatuses', component: () => import('@/views/settings/EmploymentStatusesView.vue'), meta: { title: 'Employment Statuses', titleKey: 'employment_statuses.title', descKey: 'employment_statuses.description', icon: 'pi pi-briefcase', module: 'setting' } },
      { path: 'settings/nationalities', name: 'SettingsNationalities', component: () => import('@/views/settings/NationalitiesView.vue'), meta: { title: 'Nationalities', titleKey: 'nationalities.title', descKey: 'nationalities.description', icon: 'pi pi-globe', module: 'setting' } },
      { path: 'settings/job-families', name: 'SettingsJobFamilies', component: () => import('@/views/settings/JobFamiliesView.vue'), meta: { title: 'Job Families', titleKey: 'job_families.title', descKey: 'job_families.description', icon: 'pi pi-briefcase', module: 'setting' } },
      { path: 'settings/gradings', name: 'SettingsGradings', component: () => import('@/views/settings/GradingsView.vue'), meta: { title: 'Gradings', titleKey: 'gradings.title', descKey: 'gradings.description', icon: 'pi pi-chart-bar', module: 'setting' } },
      { path: 'settings/salary-grades', name: 'SettingsSalaryGrades', component: () => import('@/views/settings/SalaryGradesView.vue'), meta: { title: 'Salary Grades', titleKey: 'salary_grades.title', descKey: 'salary_grades.description', icon: 'pi pi-chart-bar', module: 'setting' } },
      { path: 'settings/insurances', name: 'SettingsInsurances', component: () => import('@/views/settings/InsurancesView.vue'), meta: { title: 'Insurances', titleKey: 'insurances.title', descKey: 'insurances.description', icon: 'pi pi-shield', module: 'setting' } },
      { path: 'settings/company-holidays', name: 'SettingsCompanyHolidays', component: () => import('@/views/settings/CompanyHolidaysView.vue'), meta: { title: 'Company Holidays', titleKey: 'company_holidays.title', descKey: 'company_holidays.description', icon: 'pi pi-calendar', module: 'setting' } },
      { path: 'settings/numbering', name: 'SettingsNumbering', component: () => import('@/views/settings/NumberingSettingsView.vue'), meta: { title: 'Document Numbering', titleKey: 'numbering_settings.title', descKey: 'numbering_settings.description', icon: 'pi pi-hashtag', module: 'setting' } },
      { path: 'settings/document-templates', name: 'SettingsDocumentTemplates', component: () => import('@/views/settings/DocumentTemplatesView.vue'), meta: { title: 'Document Templates', titleKey: 'document_templates.title', descKey: 'document_templates.description', icon: 'pi pi-file-edit', module: 'setting' } },
      { path: 'settings/document-templates/new', name: 'SettingsDocumentTemplateNew', component: () => import('@/views/settings/DocumentTemplateForm.vue'), meta: { title: 'New Document Template', titleKey: 'document_templates.new', descKey: 'document_templates.description', icon: 'pi pi-plus', module: 'setting', backRoute: '/settings/document-templates', backLabelKey: 'document_templates.title' } },
      { path: 'settings/document-templates/:id/edit', name: 'SettingsDocumentTemplateEdit', component: () => import('@/views/settings/DocumentTemplateForm.vue'), meta: { title: 'Edit Document Template', titleKey: 'document_templates.edit', descKey: 'document_templates.description', icon: 'pi pi-pencil', module: 'setting', backRoute: '/settings/document-templates', backLabelKey: 'document_templates.title' } },
      { path: 'settings/competencies', name: 'SettingsCompetencies', component: () => import('@/views/settings/CompetenciesView.vue'), meta: { title: 'Competencies', titleKey: 'competencies.title', descKey: 'competencies.description', icon: 'pi pi-star', module: 'setting' } },
      { path: 'settings/salary-components', name: 'SettingsSalaryComponents', component: () => import('@/views/settings/SalaryComponentsView.vue'), meta: { title: 'Salary Components', titleKey: 'payroll.salary_components', descKey: 'payroll.salary_components_desc', icon: 'pi pi-list', module: 'payroll' } },
      { path: 'settings/payroll-periods', name: 'SettingsPayrollPeriods', component: () => import('@/views/settings/PayrollPeriodsView.vue'), meta: { title: 'Payroll Periods', titleKey: 'payroll.payroll_periods', descKey: 'payroll.payroll_periods_desc', icon: 'pi pi-calendar', module: 'payroll' } },
      { path: 'settings/bpjs-settings', name: 'SettingsBpjsSettings', component: () => import('@/views/settings/BpjsSettingsView.vue'), meta: { title: 'BPJS Settings', titleKey: 'payroll.bpjs', descKey: 'payroll.bpjs_desc', icon: 'pi pi-shield', module: 'payroll' } },
      { path: 'settings/pph21-settings', name: 'SettingsPph21Settings', component: () => import('@/views/settings/Pph21SettingsView.vue'), meta: { title: 'PPh21 Settings', titleKey: 'payroll.pph21', descKey: 'payroll.pph21_desc', icon: 'pi pi-percentage', module: 'payroll' } },
      { path: 'settings/salary-structure', name: 'SettingsSalaryStructure', component: () => import('@/views/settings/SalaryStructureView.vue'), meta: { title: 'Salary Structure', titleKey: 'payroll.salary_structure', descKey: 'payroll.salary_structure_desc', icon: 'pi pi-sitemap', module: 'payroll' } },
      // Performance Settings
      { path: 'settings/performance-perspectives', name: 'SettingsPerformancePerspectives', component: () => import('@/views/settings/PerformancePerspectivesView.vue'), meta: { title: 'BSC Perspectives', titleKey: 'performance_perspectives.title', descKey: 'performance_perspectives.description', icon: 'pi pi-th-large', module: 'setting' } },
      { path: 'settings/performance-ratings', name: 'SettingsPerformanceRatings', component: () => import('@/views/settings/PerformanceRatingsView.vue'), meta: { title: 'Performance Ratings', titleKey: 'performance_ratings.title', descKey: 'performance_ratings.description', icon: 'pi pi-star', module: 'setting' } },
      { path: 'settings/performance-formulas', name: 'SettingsPerformanceFormulas', component: () => import('@/views/settings/PerformanceFormulasView.vue'), meta: { title: 'KPI Formulas', titleKey: 'performance_formulas.title', descKey: 'performance_formulas.description', icon: 'pi pi-calculator', module: 'setting' } },
      { path: 'settings/performance-components', name: 'SettingsPerformanceComponents', component: () => import('@/views/settings/PerformanceComponentsView.vue'), meta: { title: 'Performance Components', titleKey: 'performance_components.title', descKey: 'performance_components.description', icon: 'pi pi-sliders-h', module: 'setting' } },
      { path: 'settings/performance-scoring', name: 'SettingsPerformanceScoring', component: () => import('@/views/settings/PerformanceScoringConfigView.vue'), meta: { title: 'Performance Scoring', titleKey: 'performance_scoring.title', descKey: 'performance_scoring.description', icon: 'pi pi-percentage', module: 'setting' } },
      { path: 'settings/rbac', name: 'SettingsRbac', component: () => import('@/views/settings/RolesPermissions.vue'), meta: { title: 'RBAC', titleKey: 'rbac.title', descKey: 'rbac.description', icon: 'pi pi-shield', module: 'setting' } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

// Navigation guard — redirect to /login if not authenticated
// Also check module access based on active modules
router.beforeEach(async (to, from, next) => {
  // Reset active modules cache on logout (navigating to login)
  if (to.name === 'Login') {
    useActiveModules().reset()
    useNotifications().reset()
  }
  const token = localStorage.getItem('tenant_token')
  const isAuthenticated = !!token

  if (to.meta.requiresAuth && !isAuthenticated) {
    next({ name: 'Login' })
    return
  }
  
  if (to.meta.guest && isAuthenticated) {
    next({ name: 'Dashboard' })
    return
  }

  // Check module access for routes with meta.module
  if (to.meta.module && isAuthenticated) {
    try {
      const activeMod = useActiveModules()
      // Ensure modules are loaded
      if (!activeMod.state.loaded) {
        await activeMod.fetchActiveModules()
      }
      if (!activeMod.hasModule(to.meta.module)) {
        // Module not active — redirect to dashboard
        next({ name: 'Dashboard' })
        return
      }
    } catch {
      // On error, allow access (fail open)
    }
  }

  next()
})

export default router
