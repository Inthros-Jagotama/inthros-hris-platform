import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/AppLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: 'Dashboard', icon: 'pi pi-home' }
      },
      // Organization Management
      {
        path: 'organizations',
        name: 'Organizations',
        component: () => import('@/views/modules/Organizations.vue'),
        meta: { title: 'Organization', icon: 'pi pi-sitemap', module: 'organization' }
      },
      // Employee Management
      {
        path: 'employees',
        name: 'Employees',
        component: () => import('@/views/modules/Employees.vue'),
        meta: { title: 'Employees', icon: 'pi pi-users', module: 'employee' }
      },
      // Job Management
      {
        path: 'job-management',
        name: 'JobManagement',
        component: () => import('@/views/modules/JobManagement.vue'),
        meta: { title: 'Job Management', icon: 'pi pi-briefcase', module: 'job-management' }
      },
      // Competency Management
      {
        path: 'competencies',
        name: 'Competencies',
        component: () => import('@/views/modules/Competencies.vue'),
        meta: { title: 'Competency', icon: 'pi pi-star', module: 'competency' }
      },
      // Employee Movement
      {
        path: 'employee-movements',
        name: 'EmployeeMovements',
        component: () => import('@/views/modules/EmployeeMovements.vue'),
        meta: { title: 'Movement', icon: 'pi pi-arrows-alt', module: 'employee-movement' }
      },
      // Time & Attendance
      {
        path: 'attendance',
        name: 'Attendance',
        component: () => import('@/views/modules/Attendance.vue'),
        meta: { title: 'Attendance', icon: 'pi pi-clock', module: 'attendance' }
      },
      // Approval Engine
      {
        path: 'approvals',
        name: 'Approvals',
        component: () => import('@/views/modules/Approvals.vue'),
        meta: { title: 'Approval', icon: 'pi pi-check-square', module: 'approval' }
      },
      // Payroll
      {
        path: 'payroll',
        name: 'Payroll',
        component: () => import('@/views/modules/Payroll.vue'),
        meta: { title: 'Payroll', icon: 'pi pi-dollar', module: 'payroll' }
      },
      // Leave & Time Off
      {
        path: 'leave',
        name: 'Leave',
        component: () => import('@/views/modules/Leave.vue'),
        meta: { title: 'Leave', icon: 'pi pi-calendar', module: 'leave' }
      },
      // Performance Management
      {
        path: 'performance',
        name: 'Performance',
        component: () => import('@/views/modules/Performance.vue'),
        meta: { title: 'Performance', icon: 'pi pi-chart-line', module: 'performance' }
      },
      // Recruitment & Onboarding
      {
        path: 'recruitment',
        name: 'Recruitment',
        component: () => import('@/views/modules/Recruitment.vue'),
        meta: { title: 'Recruitment', icon: 'pi pi-user-plus', module: 'recruitment' }
      },
      // Reimbursement
      {
        path: 'reimbursements',
        name: 'Reimbursements',
        component: () => import('@/views/modules/Reimbursements.vue'),
        meta: { title: 'Reimbursement', icon: 'pi pi-credit-card', module: 'reimbursement' }
      },
      // Training & Development
      {
        path: 'training',
        name: 'Training',
        component: () => import('@/views/modules/Training.vue'),
        meta: { title: 'Training', icon: 'pi pi-book', module: 'training' }
      },
      // Workforce Intelligence
      {
        path: 'workforce-intelligence',
        name: 'WorkforceIntelligence',
        component: () => import('@/views/modules/WorkforceIntelligence.vue'),
        meta: { title: 'Workforce Intel', icon: 'pi pi-chart-bar', module: 'workforce-intelligence' }
      },
      // Career Intelligence
      {
        path: 'career-intelligence',
        name: 'CareerIntelligence',
        component: () => import('@/views/modules/CareerIntelligence.vue'),
        meta: { title: 'Career Intel', icon: 'pi pi-road', module: 'career-intelligence' }
      }
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

export default router
