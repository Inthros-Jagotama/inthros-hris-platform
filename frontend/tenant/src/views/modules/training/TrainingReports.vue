<template>
  <div class="space-y-1">
    <!-- ── Tabs: Dashboard | Participation | Cost | Compliance ── -->
    <div class="flex items-center gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
        :class="activeTab === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab(tab.key)"
      >
        {{ t(tab.labelKey) }}
      </button>
    </div>

    <!-- ══ Dashboard analytics ══ -->
    <div v-if="activeTab === 'dashboard'" class="space-y-4">
      <div v-if="dashboardLoading" class="space-y-3">
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
          <div v-for="i in 12" :key="i" class="h-24 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        </div>
      </div>
      <template v-else>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
          <div v-for="card in dashboardCards" :key="card.key" class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 hover:shadow-md transition-shadow">
            <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t(card.labelKey) }}</p>
            <p class="text-xl font-bold text-gray-800 dark:text-gray-100 mt-1">{{ formatStat(card.value, card.money) }}</p>
          </div>
        </div>

        <!-- KPI cards: completion & pass rate + cost -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
            <div class="flex items-center justify-between mb-2">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('training.completion_rate') }}</p>
              <i class="pi pi-check-circle text-emerald-500 text-base"></i>
            </div>
            <div class="flex items-end justify-between gap-2">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ dashboard.completion_rate }}%</p>
            </div>
            <div class="mt-3 h-2 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
              <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${dashboard.completion_rate}%` }"></div>
            </div>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
            <div class="flex items-center justify-between mb-2">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('training.pass_rate') }}</p>
              <i class="pi pi-flag-fill text-sky-500 text-base"></i>
            </div>
            <div class="flex items-end justify-between gap-2">
              <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ dashboard.pass_rate }}%</p>
            </div>
            <div class="mt-3 h-2 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
              <div class="h-full rounded-full bg-sky-500 transition-all" :style="{ width: `${dashboard.pass_rate}%` }"></div>
            </div>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
            <div class="flex items-center justify-between mb-2">
              <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('training.total_training_cost') }}</p>
              <i class="pi pi-dollar text-amber-500 text-base"></i>
            </div>
            <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ formatMoney(dashboard.total_training_cost) }}</p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ t('training.cost_hint') }}</p>
          </div>
        </div>
      </template>
    </div>

    <!-- ══ Participation report ══ -->
    <div v-if="activeTab === 'participation'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2 flex-wrap">
          <SelectLabel v-model="sessionStatusFilter" :options="sessionStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_status')" class="!w-48" showClear @update:modelValue="onFilterChange" />
          <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" outlined :loading="reportLoading" @click="loadParticipation" />
          <span v-if="participationRows.length" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ participationRows.length }} {{ t('common.items') }}</span>
        </div>
      </div>
      <SkeletonTable v-if="reportLoading" :columns="participationSkeletonColumns" :rows="8" />
      <DataTable v-else :value="participationRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.reports_empty') }}</p>
          </div>
        </template>
        <Column field="employee_name" :header="t('training.employee')" sortable>
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template>
        </Column>
        <Column field="organization_name" :header="t('training.organization')" sortable style="width:180px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.organization_name || '-' }}</span></template>
        </Column>
        <Column field="course_name" :header="t('training.course')" sortable style="width:220px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.course_name }}</span></template>
        </Column>
        <Column field="session_code" :header="t('training.session_code')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.session_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="session_status" :header="t('training.session_status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="statusLabel(data.session_status)" :severity="statusSeverity(data.session_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="attendance_status" :header="t('training.attendance_status')" sortable style="width:120px">
          <template #body="{data}"><Tag :value="attStatusLabel(data.attendance_status)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="score" :header="t('training.score')" sortable style="width:90px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.score ?? '-' }}</span></template>
        </Column>
        <Column field="completion_status" :header="t('training.completion_status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="compStatusLabel(data.completion_status)" :severity="compStatusSeverity(data.completion_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>
    </div>

    <!-- ══ Cost report ══ -->
    <div v-if="activeTab === 'cost'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" outlined :loading="reportLoading" @click="loadCostReport" />
          <span v-if="costRows.length" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ costRows.length }} {{ t('common.items') }}</span>
        </div>
        <div v-if="totalReportCost > 0" class="flex items-center gap-2 text-xs">
          <span class="text-gray-400">{{ t('training.total_cost') }}</span>
          <span class="font-semibold text-emerald-600 dark:text-emerald-400">{{ formatMoney(totalReportCost) }}</span>
        </div>
      </div>
      <SkeletonTable v-if="reportLoading" :columns="costSkeletonColumns" :rows="8" />
      <DataTable v-else :value="costRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-dollar text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.reports_empty') }}</p>
          </div>
        </template>
        <Column field="session_code" :header="t('training.session_code')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.session_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="course_name" :header="t('training.course')" sortable style="width:240px">
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.course_name }}</span></template>
        </Column>
        <Column field="provider_name" :header="t('training.provider')" sortable style="width:180px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.provider_name || '-' }}</span></template>
        </Column>
        <Column field="total_cost" :header="t('training.total_cost')" sortable style="width:140px">
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.total_cost ? formatMoney(data.total_cost) : '-' }}</span></template>
        </Column>
        <Column field="participant_count" :header="t('training.participant_count')" sortable style="width:120px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.participant_count }}</span></template>
        </Column>
        <Column field="cost_per_participant" :header="t('training.cost_per_participant')" sortable style="width:160px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.cost_per_participant ? formatMoney(data.cost_per_participant) : '-' }}</span></template>
        </Column>
      </DataTable>
    </div>

    <!-- ══ Compliance report ══ -->
    <div v-if="activeTab === 'compliance'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" outlined :loading="reportLoading" @click="loadComplianceReport" />
          <span v-if="complianceRows.length" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ complianceRows.length }} {{ t('common.items') }}</span>
        </div>
      </div>
      <SkeletonTable v-if="reportLoading" :columns="complianceSkeletonColumns" :rows="8" />
      <DataTable v-else :value="complianceRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-shield text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.reports_empty') }}</p>
          </div>
        </template>
        <Column field="employee_name" :header="t('training.employee')" sortable>
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template>
        </Column>
        <Column field="organization_name" :header="t('training.organization')" sortable style="width:180px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.organization_name || '-' }}</span></template>
        </Column>
        <Column field="course_name" :header="t('training.course')" sortable style="width:240px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.course_name }}</span></template>
        </Column>
        <Column field="due_date" :header="t('training.due_date')" sortable style="width:120px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.due_date || '-' }}</span></template>
        </Column>
        <Column field="completion_status" :header="t('training.completion_status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="compStatusLabel(data.completion_status)" :severity="compStatusSeverity(data.completion_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="status" :header="t('common.status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="complianceStatusLabel(data.status)" :severity="complianceStatusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()

const tabs = [
  { key: 'dashboard', labelKey: 'training.tab_dashboard' },
  { key: 'participation', labelKey: 'training.tab_participation' },
  { key: 'cost', labelKey: 'training.tab_cost' },
  { key: 'compliance', labelKey: 'training.tab_compliance' }
]
const activeTab = ref('dashboard')

const dashboardLoading = ref(false)
const dashboard = ref({
  total_courses: 0,
  total_sessions: 0,
  total_participants: 0,
  total_providers: 0,
  total_requests: 0,
  approved_requests: 0,
  pending_requests: 0,
  completion_rate: 0,
  pass_rate: 0,
  total_training_cost: 0,
  certificates_issued: 0
})
const reportLoading = ref(false)
const sessionStatusFilter = ref(null)

const participationRows = ref([])
const costRows = ref([])
const complianceRows = ref([])

const dashboardCards = computed(() => [
  { key: 'courses', labelKey: 'training.stat_total_courses', value: dashboard.value.total_courses, money: false },
  { key: 'sessions', labelKey: 'training.stat_total_sessions', value: dashboard.value.total_sessions, money: false },
  { key: 'participants', labelKey: 'training.stat_total_participants', value: dashboard.value.total_participants, money: false },
  { key: 'providers', labelKey: 'training.stat_total_providers', value: dashboard.value.total_providers, money: false },
  { key: 'requests', labelKey: 'training.total_requests', value: dashboard.value.total_requests, money: false },
  { key: 'approved', labelKey: 'training.approved_requests', value: dashboard.value.approved_requests, money: false },
  { key: 'pending', labelKey: 'training.pending_requests', value: dashboard.value.pending_requests, money: false },
  { key: 'certificates', labelKey: 'training.certificates_issued', value: dashboard.value.certificates_issued, money: false }
])

const totalReportCost = computed(() => costRows.value.reduce((s, r) => s + (Number(r.total_cost) || 0), 0))

const sessionStatusOptions = computed(() => ['DRAFT', 'SCHEDULED', 'REGISTRATION_OPEN', 'FULL', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'].map(v => ({ label: statusLabel(v), value: v })))

const participationSkeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-14', headerWidth: 'w-14' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' }
]
const costSkeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' }
]
const complianceSkeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' }
]

function formatStat(v, money) {
  return money ? formatMoney(v) : (v ?? 0)
}
function formatMoney(v) {
  try { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v) } catch { return v }
}
function statusLabel(s) {
  const key = `training.status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusSeverity(s) {
  switch (s) {
    case 'COMPLETED': return 'success'
    case 'IN_PROGRESS': return 'info'
    case 'REGISTRATION_OPEN': return 'success'
    case 'SCHEDULED': return 'info'
    case 'DRAFT': return 'secondary'
    case 'FULL': return 'warning'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function attStatusLabel(s) {
  const key = `training.att_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusLabel(s) {
  const key = `training.comp_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusSeverity(s) {
  switch (s) {
    case 'COMPLETED': return 'success'
    case 'FAILED': return 'danger'
    case 'IN_PROGRESS': return 'info'
    default: return 'secondary'
  }
}
function complianceStatusLabel(s) {
  const key = `training.compliance_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function complianceStatusSeverity(s) {
  switch (String(s || '').toUpperCase()) {
    case 'COMPLIANT': return 'success'
    case 'OVERDUE': return 'danger'
    case 'PENDING': return 'warning'
    case 'IN_PROGRESS': return 'info'
    default: return 'secondary'
  }
}

async function loadDashboard() {
  dashboardLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/reports/dashboard')
    dashboard.value = res.data?.data || dashboard.value
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    dashboardLoading.value = false
  }
}

async function loadParticipation() {
  reportLoading.value = true
  try {
    const params = {}
    if (sessionStatusFilter.value) params.session_status = sessionStatusFilter.value
    const res = await api.get('/api/v1/tenant/trainings/reports/participation', { params })
    participationRows.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    reportLoading.value = false
  }
}

async function loadCostReport() {
  reportLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/reports/cost')
    costRows.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    reportLoading.value = false
  }
}

async function loadComplianceReport() {
  reportLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/reports/compliance')
    complianceRows.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    reportLoading.value = false
  }
}

function onFilterChange() {
  loadParticipation()
}

function switchTab(key) {
  activeTab.value = key
  if (key === 'participation' && !participationRows.value.length && !reportLoading.value) loadParticipation()
  if (key === 'cost' && !costRows.value.length && !reportLoading.value) loadCostReport()
  if (key === 'compliance' && !complianceRows.value.length && !reportLoading.value) loadComplianceReport()
}

onMounted(() => {
  loadDashboard()
})
</script>
