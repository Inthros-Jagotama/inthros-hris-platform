<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.reports') }}</h2>
      </div>
    </div>

    <!-- Event selector -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center gap-2 flex-wrap">
        <Select v-model="eventFilter" :options="eventOptions" optionLabel="label" optionValue="value" showClear filter class="w-80" :placeholder="t('competency_360.select_event')" @change="reload" />
        <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" outlined @click="reload" />
      </div>
    </div>

    <SkeletonCard v-if="loading" type="stat" :count="6" />
    <template v-else-if="error">
      <Message severity="warn" :closable="false">{{ error }}</Message>
    </template>
    <template v-else>
      <!-- HR summary -->
      <div v-if="hrReport" class="grid grid-cols-2 md:grid-cols-4 gap-2">
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.total_targets') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ hrReport.total_targets ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.finalized_targets') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ hrReport.finalized_targets ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.rater_completion') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ hrReport.rater_completion ?? 0 }}%</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.avg_score') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ hrReport.avg_score ?? 0 }}</p>
        </div>
      </div>

      <!-- Manager report -->
      <div v-if="managerReport" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.manager_report') }}</h3>
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.total_employees') }}: {{ managerReport.total_employees ?? 0 }} · {{ t('competency_360.avg_score') }}: {{ managerReport.avg_score ?? 0 }}</span>
        </div>
        <DataTable :value="managerReport.employees || []" size="small" class="!text-sm p-datatable-sm">
          <template #empty>
            <div class="text-center py-6 text-gray-400 dark:text-gray-500 text-sm">{{ t('competency_360.no_data') }}</div>
          </template>
          <Column :header="t('competency_360.employee')">
            <template #body="{data}">
              <span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span>
            </template>
          </Column>
          <Column field="overall_score" :header="t('competency_360.overall_score')" style="width:100px">
            <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ data.overall_score ?? '-' }}</span></template>
          </Column>
          <Column field="total_gap" :header="t('competency_360.total_gap')" style="width:100px">
            <template #body="{data}">
              <span :class="[(data.total_gap ?? 0) <= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400']">{{ data.total_gap ?? '-' }}</span>
            </template>
          </Column>
          <Column field="rater_completion" :header="t('competency_360.rater_completion')" style="width:130px">
            <template #body="{data}">
              <Tag :value="`${data.rater_completion ?? 0}%`" :severity="(data.rater_completion ?? 0) >= 100 ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column field="status" :header="t('common.status')" style="width:110px">
            <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
          </Column>
        </DataTable>
      </div>

      <!-- HR strengths / gaps -->
      <div v-if="hrReport" class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <h3 class="text-sm font-semibold text-emerald-600 dark:text-emerald-400 mb-3 flex items-center gap-2"><i class="pi pi-arrow-up"></i>{{ t('competency_360.top_strengths') }}</h3>
          <div v-if="hrReport.top_strengths?.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_data') }}</div>
          <div v-for="s in hrReport.top_strengths" :key="s.competency_id" class="flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0">
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ s.competency_name || s.competency_id }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ s.score }} ({{ s.gap }})</span>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <h3 class="text-sm font-semibold text-rose-600 dark:text-rose-400 mb-3 flex items-center gap-2"><i class="pi pi-arrow-down"></i>{{ t('competency_360.top_development_gaps') }}</h3>
          <div v-if="hrReport.top_development_gaps?.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_data') }}</div>
          <div v-for="d in hrReport.top_development_gaps" :key="d.competency_id" class="flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0">
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ d.competency_name || d.competency_id }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ d.score }} ({{ d.gap }})</span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Message from 'primevue/message'
import SkeletonCard from '@/components/SkeletonCard.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const events = ref([])
const employees = ref([])
const eventFilter = ref(null)
const loading = ref(false)
const error = ref('')
const managerReport = ref(null)
const hrReport = ref(null)

const eventOptions = computed(() => events.value.map(e => ({ label: eventLabel(e), value: e.id })))

function eventLabel(e) {
  const parts = [e.period_year]
  if (e.period_number) parts.unshift(`${e.period_type === 'quarter' ? 'Q' : e.period_type === 'semester' ? 'S' : 'P'}${e.period_number}`)
  return parts.join(' ')
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'finalized': return 'success'
    case 'submitted': return 'info'
    case 'pending_approval': return 'warn'
    case 'rejected': return 'danger'
    case 'closed': return 'secondary'
    default: return 'secondary'
  }
}

function employeeName(id) {
  return employees.value.find(e => e.employee_id === id)?.name || id?.slice(0, 8) || '-'
}

async function loadReferences() {
  try {
    const [evRes, empRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/events', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    ])
    events.value = evRes.status === 'fulfilled' ? (evRes.value.data?.data || []) : []
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
    if (events.value.length > 0) eventFilter.value = events.value[0].id
  } catch {
    // fail-silent
  }
}

async function reload() {
  error.value = ''
  managerReport.value = null
  hrReport.value = null
  if (!eventFilter.value) return
  loading.value = true
  try {
    const [mgrRes, hrRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/reports/manager', { params: { event_id: eventFilter.value } }),
      api.get('/api/v1/tenant/competency/reports/hr', { params: { event_id: eventFilter.value } })
    ])
    if (mgrRes.status === 'fulfilled') managerReport.value = mgrRes.value.data?.data || mgrRes.value.data
    if (hrRes.status === 'fulfilled') hrReport.value = hrRes.value.data?.data || hrRes.value.data
    if (mgrRes.status === 'rejected' && hrRes.status === 'rejected') {
      error.value = getErrorMessage(mgrRes.reason, t('competency_360.no_data'))
    }
  } catch (e) {
    error.value = getErrorMessage(e, t('competency_360.no_data'))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadReferences()
  await reload()
})
</script>
