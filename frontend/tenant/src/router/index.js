import { createRouter, createWebHistory } from 'vue-router'
import { useActiveModules } from '@/stores/activeModules'

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
        component: () => import('@/views/modules/OrganizationSummary.vue'),
        meta: { title: 'Organization', titleKey: 'organization.title', descKey: 'org_summary.description', icon: 'pi pi-building', module: 'organization' }
      },
      // Organization Management
      {
        path: 'organizations',
        name: 'Organizations',
        component: () => import('@/views/modules/Organizations.vue'),
        meta: { title: 'Organization', titleKey: 'organization.title', descKey: 'organization.description', icon: 'pi pi-sitemap', module: 'organization' }
      },
      // Employee Management
      {
        path: 'employees',
        name: 'Employees',
        component: () => import('@/views/modules/Employees.vue'),
        meta: { title: 'Employees', titleKey: 'employee.title', descKey: 'employee.description', icon: 'pi pi-users', module: 'employee' }
      },
      {
        path: 'employees/new',
        name: 'EmployeeNew',
        component: () => import('@/views/modules/EmployeeForm.vue'),
        meta: { title: 'New Employee', titleKey: 'employee.new', descKey: 'employee.description', icon: 'pi pi-user-plus', module: 'employee', backRoute: '/employees', backLabelKey: 'nav.employees' }
      },
      {
        path: 'employees/:id/edit',
        name: 'EmployeeEdit',
        component: () => import('@/views/modules/EmployeeForm.vue'),
        meta: { title: 'Edit Employee', titleKey: 'employee.edit', descKey: 'employee.description', icon: 'pi pi-user-edit', module: 'employee', backRoute: '/employees', backLabelKey: 'nav.employees' }
      },
      // Job Management
      {
        path: 'job-management',
        name: 'JobManagement',
        component: () => import('@/views/modules/JobManagement.vue'),
        meta: { title: 'Job Management', titleKey: 'job_management.title', descKey: 'job_management.page_description', icon: 'pi pi-briefcase', module: 'jobmanagement' }
      },
      {
        path: 'job-management/values',
        name: 'JobValues',
        component: () => import('@/views/modules/JobValuesIndex.vue'),
        meta: { title: 'Job Values', titleKey: 'job_management.values', descKey: 'job_values.description', icon: 'pi pi-sliders-h', module: 'jobmanagement' }
      },
      {
        path: 'job-management/values/:type',
        name: 'JobValuesType',
        component: () => import('@/views/modules/JobValuesForm.vue'),
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
        component: () => import('@/views/modules/Competencies.vue'),
        meta: { title: 'Competency', titleKey: 'competency.title', descKey: 'competency.description', icon: 'pi pi-star', module: 'competency' }
      },
      // Employee Movement
      {
        path: 'employee-movements',
        name: 'EmployeeMovements',
        component: () => import('@/views/modules/EmployeeMovements.vue'),
        meta: { title: 'Movement', titleKey: 'employee_movement.title', descKey: 'employee_movement.description', icon: 'pi pi-arrows-alt', module: 'employee-movement' }
      },
      // Time & Attendance
      {
        path: 'attendance',
        name: 'Attendance',
        component: () => import('@/views/modules/Attendance.vue'),
        meta: { title: 'Attendance', titleKey: 'attendance.title', descKey: 'attendance.description', icon: 'pi pi-clock', module: 'attendance' }
      },
      {
        path: 'attendance/admin',
        name: 'AttendanceAdmin',
        component: () => import('@/views/modules/AttendanceAdmin.vue'),
        meta: { title: 'Attendance Admin', titleKey: 'attendance.admin', descKey: 'attendance.admin_description', icon: 'pi pi-cog', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/settings',
        name: 'AttendanceSettings',
        component: () => import('@/views/modules/AttendanceSettings.vue'),
        meta: { title: 'Attendance Settings', titleKey: 'attendance.settings', descKey: 'attendance.settings_description', icon: 'pi pi-cog', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/shifts',
        name: 'AttendanceShifts',
        component: () => import('@/views/modules/AttendanceShifts.vue'),
        meta: { title: 'Shifts', titleKey: 'attendance.shifts', descKey: 'attendance.shifts_description', icon: 'pi pi-clock', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/employee-shifts',
        name: 'AttendanceEmployeeShifts',
        component: () => import('@/views/modules/AttendanceEmployeeShifts.vue'),
        meta: { title: 'Employee Shifts', titleKey: 'attendance.employee_shifts', descKey: 'attendance.employee_shifts_description', icon: 'pi pi-users', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/locations',
        name: 'AttendanceLocations',
        component: () => import('@/views/modules/AttendanceLocations.vue'),
        meta: { title: 'Locations', titleKey: 'attendance.locations', descKey: 'attendance.locations_description', icon: 'pi pi-map-marker', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/exempt-positions',
        name: 'AttendanceExemptPositions',
        component: () => import('@/views/modules/AttendanceExemptPositions.vue'),
        meta: { title: 'Exempt Positions', titleKey: 'attendance.exempt_positions', descKey: 'attendance.exempt_positions_description', icon: 'pi pi-shield', module: 'attendance', backRoute: '/attendance/admin', backLabelKey: 'attendance.admin' }
      },
      {
        path: 'attendance/overtime',
        name: 'AttendanceOvertime',
        component: () => import('@/views/modules/AttendanceOvertime.vue'),
        meta: { title: 'Overtime', titleKey: 'attendance.overtime', descKey: 'attendance.overtime_description', icon: 'pi pi-clock', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      {
        path: 'attendance/corrections',
        name: 'AttendanceCorrections',
        component: () => import('@/views/modules/AttendanceCorrections.vue'),
        meta: { title: 'Corrections', titleKey: 'attendance.corrections', descKey: 'attendance.corrections_description', icon: 'pi pi-pencil', module: 'attendance', backRoute: '/attendance', backLabelKey: 'attendance.title' }
      },
      // Approval Engine
      {
        path: 'approvals',
        name: 'Approvals',
        component: () => import('@/views/modules/Approvals.vue'),
        meta: { title: 'Approval', titleKey: 'approval.title', descKey: 'approval.description', icon: 'pi pi-check-square', module: 'approval' }
      },
      {
        path: 'approvals/flows',
        name: 'ApprovalFlows',
        component: () => import('@/views/modules/ApprovalFlows.vue'),
        meta: { title: 'Approval Flows', titleKey: 'approval.flows', descKey: 'approval.description', icon: 'pi pi-sitemap', module: 'approval' }
      },
      // Payroll
      {
        path: 'payroll',
        name: 'Payroll',
        component: () => import('@/views/modules/Payroll.vue'),
        meta: { title: 'Payroll', titleKey: 'payroll.title', descKey: 'payroll.description', icon: 'pi pi-dollar', module: 'payroll' }
      },
      // Leave & Time Off
      {
        path: 'leave',
        name: 'Leave',
        component: () => import('@/views/modules/Leave.vue'),
        meta: { title: 'Leave', titleKey: 'leave.title', descKey: 'leave.description', icon: 'pi pi-calendar', module: 'leave' }
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
        component: () => import('@/views/modules/Recruitment.vue'),
        meta: { title: 'Recruitment', titleKey: 'recruitment.title', descKey: 'recruitment.description', icon: 'pi pi-user-plus', module: 'recruitment' }
      },
      // Reimbursement
      {
        path: 'reimbursements',
        name: 'Reimbursements',
        component: () => import('@/views/modules/Reimbursements.vue'),
        meta: { title: 'Reimbursement', titleKey: 'reimbursement.title', descKey: 'reimbursement.description', icon: 'pi pi-credit-card', module: 'reimbursement' }
      },
      // Training & Development
      {
        path: 'training',
        name: 'Training',
        component: () => import('@/views/modules/Training.vue'),
        meta: { title: 'Training', titleKey: 'training.title', descKey: 'training.description', icon: 'pi pi-book', module: 'training' }
      },
      // Workforce Intelligence
      {
        path: 'workforce-intelligence',
        name: 'WorkforceIntelligence',
        component: () => import('@/views/modules/WorkforceIntelligence.vue'),
        meta: { title: 'Workforce Intel', titleKey: 'workforce_intel.title', descKey: 'workforce_intel.description', icon: 'pi pi-chart-bar', module: 'workforce-intelligence' }
      },
      // Career Intelligence
      {
        path: 'career-intelligence',
        name: 'CareerIntelligence',
        component: () => import('@/views/modules/CareerIntelligence.vue'),
        meta: { title: 'Career Intel', titleKey: 'career_intel.title', descKey: 'career_intel.description', icon: 'pi pi-road', module: 'career-intelligence' }
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
      { path: 'settings/competencies', name: 'SettingsCompetencies', component: () => import('@/views/settings/CompetenciesView.vue'), meta: { title: 'Competencies', titleKey: 'competencies.title', descKey: 'competencies.description', icon: 'pi pi-star', module: 'setting' } },
      { path: 'settings/ters', name: 'SettingsTERs', component: () => import('@/views/settings/TersView.vue'), meta: { title: 'TER', titleKey: 'ters.title', descKey: 'ters.description', icon: 'pi pi-calculator', module: 'setting' } },
      { path: 'settings/ptkps', name: 'SettingsPTKPs', component: () => import('@/views/settings/PtkpsView.vue'), meta: { title: 'PTKP', titleKey: 'ptkps.title', descKey: 'ptkps.description', icon: 'pi pi-receipt', module: 'setting' } },
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
