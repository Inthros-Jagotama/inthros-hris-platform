<template>
  <div class="space-y-4">
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import DateInput from '@/components/DateInput.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(false)
const movementReport = ref({ total: 0, by_type: {}, by_status: {} })
const contractReport = ref({ total: 0, by_status: {}, expiring: 0 })

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
})
</script>
