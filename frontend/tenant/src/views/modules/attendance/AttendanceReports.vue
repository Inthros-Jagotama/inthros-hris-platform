<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <DateInput v-model="fromDate" :placeholder="t('attendance.from_date')" class="!w-36" />
        <span class="text-gray-400">-</span>
        <DateInput v-model="toDate" :placeholder="t('attendance.to_date')" class="!w-36" />
        <Button :label="t('common.refresh')" icon="pi pi-search" size="small" :loading="loading" class="!whitespace-nowrap shrink-0" @click="loadData" />
      </div>
      <span v-if="items.length > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ items.length }} {{ t('common.items') }}</span>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable
      v-else
      :value="items"
      paginator
      :rows="20"
      :rowsPerPageOptions="[10, 20, 50, 100]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-chart-bar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.reports_empty') }}</p>
        </div>
      </template>
      <Column field="employee_id" :header="t('employee.title')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="work_date" :header="t('attendance.work_date')" sortable style="width:150px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.work_date, locale) }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable style="width:140px">
        <template #body="{data}"><Tag :value="t('attendance.status_' + data.status.toLowerCase())" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="lateness_minutes" :header="t('attendance.lateness_minutes')" sortable style="width:110px">
        <template #body="{data}"><span :class="data.lateness_minutes > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400'">{{ data.lateness_minutes }} min</span></template>
      </Column>
      <Column field="early_leave_minutes" :header="t('attendance.early_leave_minutes')" sortable style="width:110px">
        <template #body="{data}"><span :class="data.early_leave_minutes > 0 ? 'text-amber-600 dark:text-amber-400' : 'text-gray-400'">{{ data.early_leave_minutes }} min</span></template>
      </Column>
      <Column field="work_minutes" :header="t('attendance.work_minutes')" sortable style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.work_minutes }} min</span></template>
      </Column>
      <Column field="overtime_minutes" :header="t('attendance.overtime_minutes')" sortable style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.overtime_minutes || 0 }} min</span></template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { formatDate } from '@/utils/formatDate'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import DateInput from '@/components/DateInput.vue'

const { t, locale } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const employees = ref([])

function toDateOnly(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const today = new Date()
const fromDate = ref(toDateOnly(new Date(today.getFullYear(), today.getMonth(), 1)))
const toDate = ref(toDateOnly(today))

const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
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
  if (!fromDate.value || !toDate.value) return
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/attendance/reports/sessions', { params: { from: fromDate.value, to: toDate.value } })
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadEmployees()
  loadData()
})
</script>
