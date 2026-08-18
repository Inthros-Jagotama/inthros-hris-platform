<template>
  <div class="space-y-3">
    <!-- KPI karyawan (data real) -->
    <div v-if="empLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
      <div v-for="i in 4" :key="i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
        <div class="h-4 w-28 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
        <div class="h-8 w-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.kpi_total_employees') }}</span>
          <i class="pi pi-users text-lg text-emerald-500"></i>
        </div>
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ empStats.total ?? '—' }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.kpi_active') }}</span>
          <i class="pi pi-check-circle text-lg text-blue-500"></i>
        </div>
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ empStats.active ?? '—' }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.on_leave_today') }}</span>
          <i class="pi pi-calendar text-lg text-amber-500"></i>
        </div>
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ empStats.onLeave ?? '—' }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.kpi_pending_approvals') }}</span>
          <i class="pi pi-clock text-lg text-rose-500"></i>
        </div>
        <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ empStats.pending ?? '—' }}</div>
      </div>
    </div>

    <!-- ── Pie charts: jenis kelamin & status kepegawaian ── -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <!-- Komposisi karyawan per jenis kelamin -->
      <div v-if="genderTotal > 0" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.gender_distribution') }}</h2>
        <div class="flex items-center gap-6 flex-wrap">
          <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
            <circle
              v-for="seg in donutSegments"
              :key="seg.label"
              cx="60" cy="60" r="45" fill="none" stroke-width="18"
              :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
              transform="rotate(-90 60 60)"
            />
            <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ genderTotal }}</text>
            <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('employee.gender') }}</text>
          </svg>
          <div class="space-y-2 flex-1 min-w-0">
            <div v-for="seg in donutSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
              <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
              <span class="font-semibold text-gray-800 dark:text-gray-100">{{ seg.value }}</span>
              <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Komposisi karyawan per status kepegawaian -->
      <div v-if="empStatusTotal > 0" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.employment_status_distribution') }}</h2>
        <div class="flex items-center gap-6 flex-wrap">
          <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
            <circle
              v-for="seg in empStatusSegments"
              :key="seg.label"
              cx="60" cy="60" r="45" fill="none" stroke-width="18"
              :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
              transform="rotate(-90 60 60)"
            />
            <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ empStatusTotal }}</text>
            <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('employee.employment_status') }}</text>
          </svg>
          <div class="space-y-2 flex-1 min-w-0">
            <div v-for="seg in empStatusSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
              <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
              <span class="font-semibold text-gray-800 dark:text-gray-100">{{ seg.value }}</span>
              <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
            </div>
          </div>
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

    <!-- ── Quality of Hire (S-6 — Workforce Intelligence) ── -->
    <div
      v-if="wiModuleActive && qohLoading"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse"
    >
      <div class="flex items-center gap-2 mb-3">
        <div class="w-4 h-4 rounded bg-gray-200 dark:bg-gray-700"></div>
        <div class="h-4 w-44 rounded bg-gray-200 dark:bg-gray-700"></div>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2">
        <div v-for="i in 5" :key="i" class="h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50"></div>
      </div>
    </div>
    <div
      v-if="wiModuleActive && !qohLoading && qohData.hires_analyzed > 0"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3"
    >
      <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
        <div class="flex items-center gap-2">
          <i class="pi pi-bullseye text-sm text-emerald-500"></i>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.quality_of_hire_title') }}</h2>
          <span class="text-[11px] text-gray-400 dark:text-gray-500 hidden sm:inline">{{ t('dashboard.quality_of_hire_desc') }}</span>
        </div>
        <Button
          :label="t('dashboard.view_analytics')"
          icon="pi pi-chart-bar"
          size="small"
          text
          class="!text-xs"
          @click="$router.push('/workforce-intelligence/quality-of-hire')"
        />
      </div>

      <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2">
        <div class="rounded-lg border border-emerald-200 dark:border-emerald-700/50 bg-emerald-50/50 dark:bg-emerald-900/10 p-2.5">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.overall_score') }}</p>
          <p class="text-lg font-bold text-emerald-600 dark:text-emerald-400">{{ fmtScore(qohData.overall_score) }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.hires_analyzed') }}</p>
          <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ qohData.hires_analyzed }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.interview_score') }}</p>
          <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ fmtScore(qohData.interview_score) }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.onboarding_completion') }}</p>
          <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ fmtPct(qohData.onboarding_completion_rate) }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.retention_rate') }}</p>
          <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ fmtPct(qohData.retention_rate) }}</p>
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
import { fmtScore, fmtPct, buildDonutSegments } from '@/utils/dashboard'

import Button from 'primevue/button'

const { t } = useI18n()
const toast = useToast()
const activeMod = useActiveModules()

const empLoading = ref(false)
const empStats = ref({ total: null, active: null, onLeave: null, pending: null })
const genderStats = ref({ male: 0, female: 0, other: 0 })
const empStatusStats = ref({ groups: [], unclassified: 0 })
const EMP_STATUS_COLORS = ['#3b82f6', '#ec4899', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#ef4444', '#84cc16', '#f97316', '#6366f1']

// Segmen donut per status kepegawaian — nama status tenant-configurable;
// karyawan tanpa employment berjalan masuk "Tanpa Status".
const empStatusSegments = computed(() => {
  const items = (empStatusStats.value.groups || []).map(g => ({ label: g.name || t('dashboard.without_status'), value: Number(g.count) || 0 }))
  const unclassified = Number(empStatusStats.value.unclassified) || 0
  if (unclassified > 0) items.push({ label: t('dashboard.without_status'), value: unclassified })
  return buildDonutSegments(items.map((i, idx) => ({ ...i, color: EMP_STATUS_COLORS[idx % EMP_STATUS_COLORS.length] })))
})
const empStatusTotal = computed(() => empStatusSegments.value.reduce((s, i) => s + i.value, 0))

// Segmen donut jumlah per jenis kelamin.
const donutSegments = computed(() => buildDonutSegments([
  { label: t('employee.gender_m'), value: genderStats.value.male, color: '#3b82f6' },
  { label: t('employee.gender_f'), value: genderStats.value.female, color: '#ec4899' },
  { label: t('dashboard.gender_other'), value: genderStats.value.other, color: '#9ca3af' }
]))
const genderTotal = computed(() => genderStats.value.male + genderStats.value.female + genderStats.value.other)

// ── Quality of Hire (S-6), gated WI module ──
const wiModuleActive = ref(false)
const qohLoading = ref(false)
const qohData = ref({ overall_score: 0, hires_analyzed: 0, interview_score: 0, onboarding_completion_rate: 0, retention_rate: 0 })

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

async function loadEmploymentDashboard() {
  if (empLoading.value) return
  empLoading.value = true
  qohLoading.value = true
  movementLoading.value = true
  try {
    const [empRes, activeRes, leaveRes, apprRes, genderRes, empStatusRes, mvRes, qohRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/employees', { params: { page: 1, per_page: 1 } }),
      api.get('/api/v1/tenant/employees', { params: { page: 1, per_page: 1, status: 'active' } }),
      api.get('/api/v1/tenant/leave/reports/on-leave-today'),
      api.get('/api/v1/tenant/approval/tasks/pending', { params: { page: 1, per_page: 1 } }),
      api.get('/api/v1/tenant/employees/stats/gender'),
      api.get('/api/v1/tenant/employees/stats/employment-status'),
      movementModuleActive.value ? api.get('/api/v1/tenant/employee-movements/dashboard') : Promise.resolve({ data: {} }),
      wiModuleActive.value ? api.get('/api/v1/tenant/workforce-intelligence/analytics/quality-of-hire') : Promise.resolve({ data: {} })
    ])
    empStats.value = {
      total: empRes.status === 'fulfilled' ? (empRes.value.data?.total ?? null) : null,
      active: activeRes.status === 'fulfilled' ? (activeRes.value.data?.total ?? null) : null,
      onLeave: leaveRes.status === 'fulfilled' ? (leaveRes.value.data?.data?.count ?? null) : null,
      pending: apprRes.status === 'fulfilled' ? (apprRes.value.data?.total ?? null) : null
    }
    if (genderRes.status === 'fulfilled') genderStats.value = genderRes.value.data?.data || genderStats.value
    if (empStatusRes.status === 'fulfilled') empStatusStats.value = empStatusRes.value.data?.data || empStatusStats.value
    if (mvRes.status === 'fulfilled') movementData.value = mvRes.value.data?.data || movementData.value
    if (qohRes.status === 'fulfilled') qohData.value = qohRes.value.data?.data || qohData.value
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    empLoading.value = false
    qohLoading.value = false
    movementLoading.value = false
  }
}

onMounted(async () => {
  await activeMod.fetchActiveModules()
  movementModuleActive.value = activeMod.hasModule('employeemovement')
  wiModuleActive.value = activeMod.hasModule('workforce-intelligence')
  await loadEmploymentDashboard()
})
</script>
