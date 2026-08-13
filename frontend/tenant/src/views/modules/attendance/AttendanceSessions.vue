<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Select v-model="employeeFilter" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-64" :placeholder="t('attendance.filter_all_employees')" @update:modelValue="onFilterChange" />
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
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.sessions_empty') }}</p>
        </div>
      </template>
      <Column field="employee_id" :header="t('employee.title')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="work_date" :header="t('attendance.work_date')" style="width:120px" />
      <Column field="status" :header="t('common.status')" style="width:140px">
        <template #body="{data}"><Tag :value="t('attendance.status_' + data.status.toLowerCase())" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="lateness_minutes" :header="t('attendance.lateness_minutes')" style="width:110px">
        <template #body="{data}"><span :class="data.lateness_minutes > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400'">{{ data.lateness_minutes }} min</span></template>
      </Column>
      <Column field="early_leave_minutes" :header="t('attendance.early_leave_minutes')" style="width:110px">
        <template #body="{data}"><span :class="data.early_leave_minutes > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400'">{{ data.early_leave_minutes }} min</span></template>
      </Column>
      <Column field="work_minutes" :header="t('attendance.work_minutes')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.work_minutes }} min</span></template>
      </Column>
      <Column field="overtime_minutes" :header="t('attendance.overtime_minutes')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.overtime_minutes || 0 }} min</span></template>
      </Column>
    </DataTable>
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
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const employees = ref([])
const employeeFilter = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' }
]

function employeeName(id) {
  const e = employees.value.find(x => x.id === id)
  return e ? `${e.name} (${e.employee_id})` : id
}

function statusSeverity(status) {
  switch (status) {
    case 'CLOSED': return 'success'
    case 'OPEN': return 'info'
    case 'MISSING_CHECKIN':
    case 'MISSING_CHECKOUT': return 'warn'
    case 'LEAVE': return 'help'
    case 'DAY_OFF': return 'secondary'
    default: return 'secondary'
  }
}

async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (employeeFilter.value) params.employee_id = employeeFilter.value
    const res = await api.get('/api/v1/tenant/attendance/sessions', { params })
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
  loadEmployees()
  loadData()
})
</script>
