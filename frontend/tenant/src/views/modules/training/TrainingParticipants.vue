<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="sessionFilter" :options="sessionOptions" optionLabel="label" optionValue="value" filter showClear class="!w-64" :placeholder="t('training.filter_all_sessions')" @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable
      v-else
      :value="items"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.participants_empty') }}</p>
        </div>
      </template>
      <Column field="employee_name" :header="t('training.employee')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="session_name" :header="t('training.session')" style="width:200px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ sessionName(data.session_id) }}</span></template>
      </Column>
      <Column field="registration_status" :header="t('training.registration_status')" style="width:140px">
        <template #body="{data}"><Tag :value="regStatusLabel(data.registration_status)" :severity="regStatusSeverity(data.registration_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="attendance_status" :header="t('training.attendance_status')" style="width:120px">
        <template #body="{data}"><Tag :value="attStatusLabel(data.attendance_status)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="completion_status" :header="t('training.completion_status')" style="width:130px">
        <template #body="{data}"><Tag :value="compStatusLabel(data.completion_status)" :severity="compStatusSeverity(data.completion_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="score" :header="t('training.score')" style="width:90px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.score ?? '-' }}</span></template>
      </Column>
      <Column field="passed" :header="t('training.passed')" style="width:90px">
        <template #body="{data}">
          <Tag v-if="data.passed !== null && data.passed !== undefined" :value="data.passed ? t('common.yes') : t('common.no')" :severity="data.passed ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.view')" @click="router.push(`/training/sessions/${data.session_id}`)" />
        </template>
      </Column>
    </DataTable>
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
import SkeletonTable from '@/components/SkeletonTable.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const sessions = ref([])
const employees = ref([])
const courses = ref([])
const sessionFilter = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const sessionOptions = computed(() => sessions.value.map(s => ({ label: `${s.session_code} — ${courseName(s.course_id)}`, value: s.id })))

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-36', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || id
}
function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || id
}
function sessionName(id) {
  const s = sessions.value.find(x => x.id === id)
  return s ? `${s.session_code} — ${courseName(s.course_id)}` : id
}

function regStatusLabel(s) {
  const key = `training.reg_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function regStatusSeverity(s) {
  switch (s) {
    case 'REGISTERED': return 'success'
    case 'APPROVED': return 'info'
    case 'NOMINATED': return 'warning'
    case 'REQUESTED': return 'info'
    case 'WAITLISTED': return 'warning'
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

async function loadReferences() {
  const [sRes, eRes, cRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/trainings/sessions', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } })
  ])
  sessions.value = sRes.status === 'fulfilled' ? (sRes.value.data?.data || []) : []
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (sessionFilter.value) params.session_id = sessionFilter.value
    const res = await api.get('/api/v1/tenant/trainings/participants', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
