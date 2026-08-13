<template>
  <div class="space-y-4">
    <!-- ── Hub navigation cards: Movements & Contracts (highlight) ── -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <div
        class="bg-gradient-to-r from-emerald-50 to-teal-50 dark:from-emerald-950/50 dark:to-teal-950/30 border border-emerald-200 dark:border-emerald-800/60 rounded-lg p-4 flex items-center gap-3 cursor-pointer hover:border-emerald-400 hover:shadow-lg dark:hover:shadow-gray-900/50 hover:-translate-y-0.5 transition-all group"
        @click="router.push('/admin/career/movements')"
      >
        <div class="w-11 h-11 rounded-xl bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center shrink-0 shadow-sm group-hover:shadow-md group-hover:scale-105 transition-all">
          <i class="pi pi-arrows-alt text-white"></i>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('employee_movement.movements') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ t('employee_movement.card_movements_desc') }}</p>
        </div>
        <div class="shrink-0 flex flex-col items-center gap-1">
          <span class="text-[10px] font-semibold uppercase tracking-wider text-emerald-600 dark:text-emerald-400">{{ t('common.open') }}</span>
          <i class="pi pi-arrow-right text-emerald-500 group-hover:translate-x-0.5 transition-all"></i>
        </div>
      </div>
      <div
        class="bg-gradient-to-r from-indigo-50 to-blue-50 dark:from-indigo-950/50 dark:to-blue-950/30 border border-indigo-200 dark:border-indigo-800/60 rounded-lg p-4 flex items-center gap-3 cursor-pointer hover:border-indigo-400 hover:shadow-lg dark:hover:shadow-gray-900/50 hover:-translate-y-0.5 transition-all group"
        @click="router.push('/admin/career/contracts')"
      >
        <div class="w-11 h-11 rounded-xl bg-gradient-to-br from-indigo-500 to-blue-600 flex items-center justify-center shrink-0 shadow-sm group-hover:shadow-md group-hover:scale-105 transition-all">
          <i class="pi pi-file-edit text-white"></i>
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('employee_movement.contracts') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ t('employee_movement.card_contracts_desc') }}</p>
        </div>
        <div class="shrink-0 flex flex-col items-center gap-1">
          <span class="text-[10px] font-semibold uppercase tracking-wider text-indigo-600 dark:text-indigo-400">{{ t('common.open') }}</span>
          <i class="pi pi-arrow-right text-indigo-500 group-hover:translate-x-0.5 transition-all"></i>
        </div>
      </div>
    </div>

    <!-- ── Filter bar ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
      <div class="flex items-center gap-2 flex-wrap">
        <DateInput v-model="dateFrom" :placeholder="t('employee_movement.report_date_from')" class="!w-36" />
        <span class="text-gray-400">-</span>
        <DateInput v-model="dateTo" :placeholder="t('employee_movement.report_date_to')" class="!w-36" />
        <Select
          v-model="filterOrg"
          :options="organizationOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.report_all_orgs')"
          class="!w-48"
          size="small"
          showClear
        />
        <Select
          v-model="filterPosition"
          :options="positionOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.report_all_positions')"
          class="!w-48"
          size="small"
          showClear
        />
        <Select
          v-model="filterType"
          :options="typeOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.filter_all_types')"
          class="!w-40"
          size="small"
          showClear
        />
        <Select
          v-model="filterStatus"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.filter_all_status')"
          class="!w-36"
          size="small"
          showClear
        />
        <Button :label="t('common.refresh')" icon="pi pi-search" size="small" :loading="loading" class="!whitespace-nowrap shrink-0" @click="loadData" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <template v-else>
      <!-- ── Movement Report ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-arrows-alt text-sm text-emerald-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('employee_movement.movement_report') }}</h2>
          </div>
          <Tag :value="`${movementReport.total} ${t('common.items')}`" severity="info" class="!text-xs !px-2 !py-0.5" />
        </div>

        <!-- Stat cards: by type -->
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
          <div
            v-for="t in movementTypeList"
            :key="t.value"
            class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 flex items-center justify-between hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow"
          >
            <div class="min-w-0">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t.label }}</p>
              <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ movementReport.by_type?.[t.value] || 0 }}</p>
            </div>
            <i :class="[typeIcon(t.value), typeIconColor(t.value)]" class="text-lg shrink-0"></i>
          </div>
          <div class="rounded-lg border border-emerald-200 dark:border-emerald-900/40 bg-emerald-50/50 dark:bg-emerald-900/10 p-3 flex items-center justify-between">
            <div>
              <p class="text-xs font-medium text-emerald-600 dark:text-emerald-400 uppercase tracking-wider">{{ t('common.total') }}</p>
              <p class="text-xl font-bold text-emerald-700 dark:text-emerald-300">{{ movementReport.total }}</p>
            </div>
            <i class="pi pi-chart-bar text-lg text-emerald-500"></i>
          </div>
        </div>

        <!-- Status breakdown -->
        <div class="flex items-center gap-2 flex-wrap mt-3 pt-3 border-t border-gray-100 dark:border-gray-800">
          <span class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ t('employee_movement.report_by_status') }}</span>
          <Tag
            v-for="s in statusList"
            :key="s.value"
            :value="`${s.label}: ${movementReport.by_status?.[s.value] || 0}`"
            :severity="statusSeverity(s.value)"
            class="!text-xs !px-2 !py-0.5"
          />
        </div>
      </div>

      <!-- ── Contract Report ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center gap-2 mb-3">
          <i class="pi pi-file-edit text-sm text-emerald-500"></i>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('employee_movement.contract_report') }}</h2>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-3">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('employee_movement.status_active') }}</p>
            <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ contractReport.by_status?.active || 0 }}</p>
          </div>
          <div class="rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-3">
            <p class="text-xs font-medium text-amber-600 dark:text-amber-400 uppercase tracking-wider">{{ t('employee_movement.report_expiring') }}</p>
            <p class="text-xl font-bold text-amber-700 dark:text-amber-300">{{ contractReport.expiring || 0 }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('employee_movement.status_expired') }}</p>
            <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ contractReport.by_status?.expired || 0 }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('employee_movement.status_extended') }}</p>
            <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ contractReport.by_status?.extended || 0 }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('employee_movement.status_terminated') }}</p>
            <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ contractReport.by_status?.terminated || 0 }}</p>
          </div>
        </div>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-3">{{ t('employee_movement.report_contract_hint') }}</p>
      </div>
    </template>

    <!-- ── Career Timeline (plan §12.8 — read model dari movements + employments + contracts) ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
      <div class="flex items-center gap-2 mb-3">
        <i class="pi pi-history text-sm text-emerald-500"></i>
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('employee_movement.career_timeline') }}</h2>
      </div>

      <!-- Employee selector -->
      <div class="flex items-center gap-2 flex-wrap mb-4">
        <Select
          v-model="timelineEmployeeId"
          :options="employeeOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.timeline_select_employee')"
          filter
          class="!w-80"
          size="small"
          :loading="employeesLoading"
        />
        <Button
          :label="t('common.view')"
          icon="pi pi-history"
          size="small"
          :loading="timelineLoading"
          :disabled="!timelineEmployeeId"
          class="!whitespace-nowrap shrink-0"
          @click="loadTimeline"
        />
      </div>

      <template v-if="timelineLoading">
        <div class="space-y-4">
          <div v-for="i in 3" :key="i" class="flex gap-3">
            <div class="w-6 h-6 rounded-full bg-gray-100 dark:bg-gray-700/50 shrink-0"></div>
            <div class="flex-1 space-y-1.5">
              <div class="h-3.5 w-40 rounded bg-gray-100 dark:bg-gray-700/50"></div>
              <div class="h-3 w-64 max-w-full rounded bg-gray-100 dark:bg-gray-700/50"></div>
            </div>
          </div>
        </div>
      </template>

      <template v-else-if="timelineData">
        <!-- Current position summary -->
        <div
          v-if="timelineData.current_position"
          class="mb-4 rounded-lg border border-emerald-200 dark:border-emerald-800/60 bg-emerald-50/50 dark:bg-emerald-900/10 p-3 flex items-center gap-3 flex-wrap"
        >
          <div class="w-9 h-9 rounded-lg bg-emerald-500 flex items-center justify-center shrink-0">
            <i class="pi pi-briefcase text-white text-sm"></i>
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-xs font-medium text-emerald-600 dark:text-emerald-400 uppercase tracking-wider">{{ t('employee_movement.timeline_current_position') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">
              {{ timelineData.employee_name || '-' }}
              <span v-if="timelineData.employee_code" class="text-gray-400 font-normal"> ({{ timelineData.employee_code }})</span>
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              <span>{{ timelineData.current_position.position_name || '-' }}</span>
              <span v-if="timelineData.current_position.organization_name"> · {{ timelineData.current_position.organization_name }}</span>
              <span v-if="timelineData.current_position.employment_status_name"> · {{ timelineData.current_position.employment_status_name }}</span>
            </p>
          </div>
        </div>

        <!-- Timeline events — grouped by year, ASC kronologis (plan §12.8) -->
        <ol v-if="timelineData.timeline?.length" class="relative border-l border-gray-200 dark:border-gray-700 ml-3 space-y-4">
          <template v-for="(group, gi) in timelineByYear" :key="gi">
            <li class="relative pl-6">
              <p class="text-[11px] font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500">{{ group.year }}</p>
            </li>
            <li v-for="(ev, i) in group.items" :key="`${gi}-${i}`" class="relative pl-6">
              <span
                class="absolute -left-[13px] top-0 w-6 h-6 rounded-full flex items-center justify-center ring-4 ring-white dark:ring-gray-800"
                :class="timelineDotClass(ev)"
              >
                <i :class="timelineEventIcon(ev)" class="text-[10px] text-white"></i>
              </span>
              <div class="flex items-center gap-2 flex-wrap">
                <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ timelineEventLabel(ev) }}</p>
                <Tag v-if="ev.movement_type" :value="typeLabel(ev.movement_type)" :severity="typeSeverity(ev.movement_type)" class="!text-[10px] !px-1.5 !py-0" />
                <Tag v-if="ev.contract_type" :value="contractTypeLabel(ev.contract_type)" severity="info" class="!text-[10px] !px-1.5 !py-0" />
              </div>
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ formatDate(ev.date, locale) }}</p>
              <p v-if="timelineEventTitle(ev)" class="text-sm text-gray-600 dark:text-gray-300 mt-0.5">{{ timelineEventTitle(ev) }}</p>
              <p v-if="ev.description" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ ev.description }}</p>
            </li>
          </template>
        </ol>
        <div v-else class="text-center py-8 text-gray-400 dark:text-gray-500">
          <i class="pi pi-history text-2xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('employee_movement.timeline_empty') }}</p>
        </div>
      </template>

      <div v-else class="text-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-user text-2xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee_movement.timeline_hint') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import DateInput from '@/components/DateInput.vue'

const { t, locale } = useI18n()
const router = useRouter()
const toast = useToast()

const loading = ref(false)
const movementReport = ref({ total: 0, by_type: {}, by_status: {} })
const contractReport = ref({ total: 0, by_status: {}, expiring: 0 })

// ── Career Timeline (plan §12.8) ──
const timelineEmployeeId = ref(null)
const employees = ref([])
const employeesLoading = ref(false)
const timelineLoading = ref(false)
const timelineData = ref(null)

const dateFrom = ref('')
const dateTo = ref('')
const filterOrg = ref(null)
const filterPosition = ref(null)
const filterType = ref(null)
const filterStatus = ref(null)

const organizations = ref([])

const typeOptions = computed(() => [
  'promotion', 'demotion', 'mutation', 'contract_extension', 'status_change', 'retirement', 'offboarding', 'other'
].map(v => ({ label: typeLabel(v), value: v })))

const statusOptions = computed(() => [
  'draft', 'pending_approval', 'approved', 'rejected', 'executed', 'cancelled', 'cancellation_pending'
].map(v => ({ label: statusLabel(v), value: v })))

// Movement types for stat cards — all known types so empty buckets show 0.
const movementTypeList = computed(() => [
  'promotion', 'demotion', 'mutation', 'contract_extension', 'status_change', 'retirement', 'offboarding', 'other'
].map(v => ({ label: typeLabel(v), value: v })))

const statusList = computed(() => [
  'draft', 'pending_approval', 'approved', 'rejected', 'executed', 'cancelled', 'cancellation_pending'
].map(v => ({ label: statusLabel(v), value: v })))

const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id })))
const positionOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id })))
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.id })))

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' }
]

function typeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function statusLabel(status) {
  const key = `employee_movement.status_${status}`
  return t(key) !== key ? t(key) : status
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

function typeSeverity(type) {
  switch (type) {
    case 'promotion': return 'success'
    case 'demotion': return 'danger'
    case 'mutation': return 'info'
    case 'contract_extension': return 'warning'
    case 'status_change': return 'info'
    case 'retirement': return 'secondary'
    case 'offboarding': return 'danger'
    default: return 'secondary'
  }
}

function contractTypeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

// ── Career Timeline helpers ──
const timelineByYear = computed(() => {
  const entries = timelineData.value?.timeline || []
  const groups = []
  let currentYear = null
  for (const ev of entries) {
    const year = String(ev.date || '').slice(0, 4)
    if (year !== currentYear) {
      currentYear = year
      groups.push({ year, items: [] })
    }
    groups[groups.length - 1].items.push(ev)
  }
  return groups
})

function timelineEventLabel(ev) {
  if (ev.event_type === 'JOINED') return t('employee_movement.timeline_joined')
  if (ev.event_type === 'CONTRACT') return t('employee_movement.timeline_contract')
  if (ev.movement_type) return typeLabel(ev.movement_type)
  return ev.event_type || ''
}

function timelineEventTitle(ev) {
  // JOINED: title = nama posisi/org; CONTRACT: title = nomor kontrak;
  // MOVEMENT: badge tipe sudah menampilkan label, title mentah di-skip.
  if (ev.event_type === 'MOVEMENT') return ''
  return ev.title || ''
}

function timelineEventIcon(ev) {
  if (ev.event_type === 'JOINED') return 'pi pi-user-plus'
  if (ev.event_type === 'CONTRACT') return 'pi pi-file-edit'
  return typeIcon(ev.movement_type)
}

function timelineDotClass(ev) {
  if (ev.event_type === 'JOINED') return 'bg-emerald-500'
  if (ev.event_type === 'CONTRACT') return 'bg-amber-500'
  switch (ev.movement_type) {
    case 'promotion': return 'bg-emerald-500'
    case 'demotion': return 'bg-red-500'
    case 'mutation': return 'bg-sky-500'
    case 'contract_extension': return 'bg-amber-500'
    case 'status_change': return 'bg-indigo-500'
    case 'retirement': return 'bg-gray-400'
    case 'offboarding': return 'bg-rose-500'
    default: return 'bg-slate-400'
  }
}

async function loadEmployees() {
  employeesLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch (err) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: getErrorMessage(err), life: 4000 })
  } finally {
    employeesLoading.value = false
  }
}

async function loadTimeline() {
  if (!timelineEmployeeId.value) return
  timelineLoading.value = true
  timelineData.value = null
  try {
    const res = await api.get(`/api/v1/tenant/employee-movements/employees/${timelineEmployeeId.value}/career-history`)
    timelineData.value = res.data?.data || null
  } catch (err) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: getErrorMessage(err), life: 4000 })
  } finally {
    timelineLoading.value = false
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

function statusSeverity(status) {
  switch (status) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'rejected': return 'danger'
    case 'executed': return 'success'
    case 'cancelled': return 'secondary'
    case 'cancellation_pending': return 'warning'
    default: return 'secondary'
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = {}
    if (dateFrom.value) params.date_from = dateFrom.value
    if (dateTo.value) params.date_to = dateTo.value
    if (filterOrg.value) params.organization_id = filterOrg.value
    if (filterPosition.value) params.position_id = filterPosition.value
    if (filterType.value) params.movement_type = filterType.value
    if (filterStatus.value) params.status = filterStatus.value

    const [moveRes, contractRes] = await Promise.all([
      api.get('/api/v1/tenant/employee-movements/reports/movements', { params }),
      api.get('/api/v1/tenant/employee-movements/reports/contracts')
    ])
    movementReport.value = moveRes.data?.data || { total: 0, by_type: {}, by_status: {} }
    contractReport.value = contractRes.data?.data || { total: 0, by_status: {}, expiring: 0 }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadReferences() {
  try {
    const res = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
    organizations.value = res.data?.data || []
  } catch {
    organizations.value = []
  }
}

onMounted(() => {
  loadReferences()
  loadData()
  loadEmployees()
})
</script>
