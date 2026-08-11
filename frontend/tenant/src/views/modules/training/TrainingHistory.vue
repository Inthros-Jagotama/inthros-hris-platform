<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="employeeFilter" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('training.select_employee')" class="!w-72" @update:modelValue="onEmployeeChange" />
        <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" outlined :loading="loading" @click="loadHistory" />
        <span v-if="items.length" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ items.length }} {{ t('common.items') }}</span>
      </div>
    </div>

    <div v-if="!employeeFilter" class="rounded-lg border border-dashed border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800/50 py-14 text-center">
      <i class="pi pi-history text-3xl text-gray-300 dark:text-gray-600 mb-3"></i>
      <p class="text-sm text-gray-400 dark:text-gray-500">{{ t('training.history_select_hint') }}</p>
    </div>

    <template v-else>
      <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
      <DataTable v-else :value="items" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-history text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.history_empty') }}</p>
          </div>
        </template>
        <Column field="session_code" :header="t('training.session_code')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.session_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="course_name" :header="t('training.course')" sortable style="width:240px">
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.course_name }}</span></template>
        </Column>
        <Column field="start_date" :header="t('training.start_date')" sortable style="width:130px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.start_date || '-' }}</span></template>
        </Column>
        <Column field="end_date" :header="t('training.end_date')" sortable style="width:130px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.end_date || '-' }}</span></template>
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
        <Column field="completion_date" :header="t('training.completion_date')" sortable style="width:130px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.completion_date || '-' }}</span></template>
        </Column>
        <Column field="certificate_no" :header="t('training.certificate_no')" sortable style="width:160px">
          <template #body="{data}">
            <Tag v-if="data.certificate_no" :value="data.certificate_no" severity="success" class="!text-xs !px-1.5 !py-0.5" />
            <span v-else class="text-gray-400">-</span>
          </template>
        </Column>
      </DataTable>
    </template>
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

const employees = ref([])
const employeeFilter = ref(null)
const items = ref([])
const loading = ref(false)

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-14', headerWidth: 'w-14' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-28', headerWidth: 'w-24' }
]

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

async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}

async function loadHistory() {
  if (!employeeFilter.value) return
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/history', { params: { employee_id: employeeFilter.value } })
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onEmployeeChange() {
  items.value = []
  loadHistory()
}

onMounted(() => {
  loadEmployees()
})
</script>
