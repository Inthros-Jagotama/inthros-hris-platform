<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('dashboard.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ t('dashboard.welcome') }}</p>
      </div>
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
import { ref, computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import SelectButton from 'primevue/selectbutton'

const { t } = useI18n()

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
</script>
