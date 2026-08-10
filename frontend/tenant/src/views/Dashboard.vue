<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-end">
      <div class="flex items-center gap-2">
        <SelectButton
          v-model="periodFilter"
          :options="periodOptions"
          optionLabel="label"
          optionValue="value"
          size="small"
        />
      </div>
    </div>
    <!-- KPI Cards Row -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
      <div
        v-for="kpi in kpiCards"
        :key="kpi.label"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow"
      >
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ kpi.label }}</span>
          <i :class="[kpi.icon, kpi.iconColor]" class="text-lg"></i>
        </div>
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ kpi.value }}</div>
        <div class="flex items-center gap-1 mt-1">
          <i
            :class="kpi.trend >= 0 ? 'pi pi-arrow-up text-emerald-500' : 'pi pi-arrow-down text-rose-500'"
            class="text-sm"
          ></i>
          <span :class="kpi.trend >= 0 ? 'text-emerald-600' : 'text-rose-600'" class="text-sm font-medium">
            {{ Math.abs(kpi.trend) }}%
          </span>
          <span class="text-sm text-gray-400 dark:text-gray-500">{{ t('dashboard.vs_last_month') }}</span>
        </div>
      </div>
    </div>

    <!-- ── Employee Movement & Contracts (plan §12.18 — HR Dashboard) ── -->
    <div
      v-if="movementModuleActive && movementLoading"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse"
    >
      <div class="flex items-center gap-2 mb-3">
        <div class="w-4 h-4 rounded bg-gray-200 dark:bg-gray-700"></div>
        <div class="h-4 w-40 rounded bg-gray-200 dark:bg-gray-700"></div>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
        <div v-for="i in 8" :key="i" class="h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50"></div>
      </div>
    </div>
    <div
      v-if="movementModuleActive && !movementLoading"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3"
    >
      <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
        <div class="flex items-center gap-2">
          <i class="pi pi-arrows-alt text-sm text-emerald-500"></i>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.movement_title') }}</h2>
        </div>
        <Button
          :label="t('dashboard.view_reports')"
          icon="pi pi-chart-bar"
          size="small"
          text
          class="!text-xs"
          @click="$router.push('/admin/career/reports')"
        />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <!-- Movement by type -->
        <div class="lg:col-span-2">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
            <div
              v-for="mt in movementTypeList"
              :key="mt.value"
              class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5 flex items-center justify-between gap-2 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow"
            >
              <div class="min-w-0">
                <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ mt.label }}</p>
                <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ movementData.movement_by_type?.[mt.value] || 0 }}</p>
              </div>
              <i :class="[typeIcon(mt.value), typeIconColor(mt.value)]" class="text-base shrink-0"></i>
            </div>
          </div>
        </div>

        <!-- Pending approval + effective this month -->
        <div class="space-y-3">
          <div class="rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-3">
            <p class="text-xs font-medium text-amber-600 dark:text-amber-400 uppercase tracking-wider flex items-center gap-1.5">
              <i class="pi pi-clock text-xs"></i>{{ t('dashboard.movement_pending') }}
            </p>
            <p class="text-2xl font-bold text-amber-700 dark:text-amber-300">{{ movementData.pending_approval || 0 }}</p>
          </div>
          <div class="rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-3">
            <p class="text-xs font-medium text-sky-600 dark:text-sky-400 uppercase tracking-wider flex items-center gap-1.5">
              <i class="pi pi-calendar text-xs"></i>{{ t('dashboard.movement_effective_month') }}
            </p>
            <p class="text-2xl font-bold text-sky-700 dark:text-sky-300">{{ movementData.effective_this_month || 0 }}</p>
          </div>
        </div>
      </div>

      <!-- Contract summary -->
      <div class="mt-3 pt-3 border-t border-gray-100 dark:border-gray-800">
        <div class="flex items-center gap-2 mb-2">
          <i class="pi pi-file-edit text-xs text-gray-400"></i>
          <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.contract_title') }}</span>
        </div>
        <div class="grid grid-cols-3 gap-2">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.contract_active') }}</p>
            <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ movementData.contracts?.active || 0 }}</p>
          </div>
          <div class="rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-2.5">
            <p class="text-[11px] text-amber-600 dark:text-amber-400">{{ t('dashboard.contract_expiring') }}</p>
            <p class="text-lg font-bold text-amber-700 dark:text-amber-300">{{ movementData.contracts?.expiring || 0 }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.contract_expired') }}</p>
            <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ movementData.contracts?.expired || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Quick Access Modules -->
      <div class="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.quick_access') }}</h2>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-2">
          <div
            v-for="mod in quickModules"
            :key="mod.name"
            class="flex flex-col items-center gap-1.5 p-2.5 rounded-lg cursor-pointer hover:bg-emerald-50 dark:hover:bg-emerald-900/20 hover:border-emerald-200 dark:hover:border-emerald-700 border border-transparent transition-all"
            @click="$router.push(mod.route)"
          >
            <div :class="mod.bg" class="w-9 h-9 rounded-lg flex items-center justify-center">
              <i :class="[mod.icon, mod.color]" class="text-sm"></i>
            </div>
            <span class="text-sm text-gray-600 dark:text-gray-300 text-center leading-tight">{{ mod.name }}</span>
          </div>
        </div>
      </div>
      <!-- Recent Activity -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.recent_activity') }}</h2>
        <div class="space-y-3">
          <div v-for="(activity, i) in recentActivities" :key="i" class="flex items-start gap-2.5">
            <div :class="activity.dotColor" class="w-2 h-2 rounded-full mt-1.5 shrink-0"></div>
            <div class="min-w-0">
              <p class="text-sm text-gray-700 dark:text-gray-200">{{ activity.text }}</p>
              <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-0.5">{{ activity.time }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import { useActiveModules } from '@/stores/activeModules'
import api from '@/services/api'

import SelectButton from 'primevue/selectbutton'
import Button from 'primevue/button'

const { t } = useI18n()
const toast = useToast()
const activeMod = useActiveModules()

const periodFilter = ref('this-month')
const periodOptions = computed(() => [
  { label: t('dashboard.this_month'), value: 'this-month' },
  { label: t('dashboard.this_quarter'), value: 'this-quarter' },
  { label: t('dashboard.this_year'), value: 'this-year' }
])
const kpiCards = computed(() => [
  { label: t('dashboard.kpi_total_employees'), value: '1,247', icon: 'pi pi-users', iconColor: 'text-emerald-500', trend: 3.2 },
  { label: t('dashboard.kpi_active_today'), value: '1,183', icon: 'pi pi-check-circle', iconColor: 'text-blue-500', trend: 1.5 },
  { label: t('dashboard.kpi_on_leave'), value: '42', icon: 'pi pi-calendar', iconColor: 'text-amber-500', trend: -2.1 },
  { label: t('dashboard.kpi_pending_approvals'), value: '28', icon: 'pi pi-clock', iconColor: 'text-rose-500', trend: 12.5 }
])
const quickModules = computed(() => [
  { name: t('dashboard.employees'), icon: 'pi pi-users', route: '/employees', bg: 'bg-blue-50', color: 'text-blue-600' },
  { name: t('dashboard.attendance'), icon: 'pi pi-clock', route: '/attendance', bg: 'bg-emerald-50', color: 'text-emerald-600' },
  { name: t('dashboard.leave'), icon: 'pi pi-calendar', route: '/leave', bg: 'bg-amber-50', color: 'text-amber-600' },
  { name: t('dashboard.payroll'), icon: 'pi pi-dollar', route: '/payroll', bg: 'bg-indigo-50', color: 'text-indigo-600' },
  { name: t('dashboard.approvals'), icon: 'pi pi-check-square', route: '/approvals', bg: 'bg-violet-50', color: 'text-violet-600' },
  { name: t('dashboard.performance'), icon: 'pi pi-chart-line', route: '/performance', bg: 'bg-cyan-50', color: 'text-cyan-600' },
  { name: t('dashboard.training'), icon: 'pi pi-book', route: '/training', bg: 'bg-orange-50', color: 'text-orange-600' },
  { name: t('dashboard.recruitment'), icon: 'pi pi-user-plus', route: '/recruitment', bg: 'bg-rose-50', color: 'text-rose-600' },
  { name: t('dashboard.organization'), icon: 'pi pi-sitemap', route: '/organizations', bg: 'bg-teal-50', color: 'text-teal-600' },
  { name: t('dashboard.reimbursement'), icon: 'pi pi-credit-card', route: '/reimbursements', bg: 'bg-sky-50', color: 'text-sky-600' },
  { name: t('dashboard.workforce_intel'), icon: 'pi pi-chart-bar', route: '/workforce-intelligence', bg: 'bg-slate-50', color: 'text-slate-600' },
  { name: t('dashboard.career_intel'), icon: 'pi pi-road', route: '/career-intelligence', bg: 'bg-pink-50', color: 'text-pink-600' }
])
const recentActivities = [
  { text: '15 new employees added this week', time: '2 hours ago', dotColor: 'bg-emerald-400' },
  { text: 'Payroll run for August completed', time: '5 hours ago', dotColor: 'bg-blue-400' },
  { text: '3 leave requests pending approval', time: '1 day ago', dotColor: 'bg-amber-400' },
  { text: 'Performance reviews Q3 initiated', time: '2 days ago', dotColor: 'bg-violet-400' },
  { text: 'Training session "Leadership 101" scheduled', time: '3 days ago', dotColor: 'bg-orange-400' }
]

// ── HR Dashboard: Employee Movement & Contracts (plan §12.18) ──
const movementModuleActive = ref(false)
const movementLoading = ref(false)
const movementData = ref({ movement_by_type: {}, pending_approval: 0, effective_this_month: 0, contracts: {} })

const movementTypeList = computed(() => [
  'promotion', 'demotion', 'mutation', 'contract_extension', 'status_change', 'retirement', 'offboarding', 'other'
].map(v => ({ label: typeLabel(v), value: v })))

function typeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function typeIcon(type) {
  switch (type) {
    case 'promotion': return 'pi pi-arrow-up'
    case 'demotion': return 'pi pi-arrow-down'
    case 'mutation': return 'pi pi-shuffle'
    case 'contract_extension': return 'pi pi-file-edit'
    case 'status_change': return 'pi pi-id-card'
    case 'retirement': return 'pi pi-sun'
    case 'offboarding': return 'pi pi-sign-out'
    default: return 'pi pi-circle'
  }
}

function typeIconColor(type) {
  switch (type) {
    case 'promotion': return 'text-emerald-500'
    case 'demotion': return 'text-red-500'
    case 'mutation': return 'text-sky-500'
    case 'contract_extension': return 'text-amber-500'
    case 'status_change': return 'text-indigo-500'
    case 'retirement': return 'text-gray-400'
    case 'offboarding': return 'text-red-400'
    default: return 'text-gray-400'
  }
}

async function loadMovementDashboard() {
  movementLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/employee-movements/dashboard')
    movementData.value = res.data?.data || movementData.value
  } catch (e) {
    // Jangan ganggu dashboard utama — kartu HR hanya menampilkan 0 bila gagal.
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    movementLoading.value = false
  }
}

onMounted(async () => {
  // Hanya tampilkan kartu HR bila module employeemovement aktif.
  await activeMod.fetchActiveModules()
  movementModuleActive.value = activeMod.hasModule('employeemovement')
  if (movementModuleActive.value) {
    loadMovementDashboard()
  }
})
</script>
