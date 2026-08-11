<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_request_status')" class="!w-48" showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.request_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-send text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.requests_empty') }}</p>
        </div>
      </template>
      <Column field="employee_id" :header="t('training.request_employee')" style="width:200px">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="course_id" :header="t('training.request_course')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ courseName(data.course_id) }}</span></template>
      </Column>
      <Column field="requested_date" :header="t('training.request_requested_date')" style="width:130px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.requested_date || '-' }}</span></template>
      </Column>
      <Column field="priority" :header="t('training.request_priority')" style="width:100px">
        <template #body="{data}"><Tag :value="priorityLabel(data.priority)" :severity="prioritySeverity(data.priority)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="status" :header="t('training.request_status')" style="width:150px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="supervisor_note" :header="t('training.request_supervisor_note')">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400 line-clamp-1" v-tooltip.top="data.supervisor_note || ''">{{ data.supervisor_note || '-' }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:180px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button v-if="['DRAFT', 'REJECTED'].includes(data.status)" icon="pi pi-send" size="small" text severity="success" v-tooltip.left="t('training.request_submit')" @click="confirmSubmit(data)" />
            <Button v-if="['DRAFT', 'SUBMITTED', 'PENDING_APPROVAL'].includes(data.status)" icon="pi pi-times" size="small" text severity="danger" v-tooltip.left="t('training.request_cancel')" @click="confirmCancel(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Dialog: request -->
    <Dialog v-model:visible="dialogVisible" :header="t('training.request_new')" modal :style="{ width: '560px' }" @hide="resetForm">
      <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.request_employee')" required :errors="errors?.employee_id">
            <SelectLabel v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.employee_id }" />
          </FormRow>
          <FormRow :label="t('training.request_course')" required :errors="errors?.course_id">
            <SelectLabel v-model="form.course_id" :options="courseOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.course_id }" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.request_session')">
            <SelectLabel v-model="form.session_id" :options="sessionOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
          </FormRow>
          <FormRow :label="t('training.request_requested_date')" required :errors="errors?.requested_date">
            <DateInput v-model="form.requested_date" :class="{ 'p-invalid': errors?.requested_date }" />
          </FormRow>
        </div>
        <FormRow :label="t('training.request_reason')">
          <TextInput v-model="form.reason" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('training.request_priority')">
          <SelectLabel v-model="form.priority" :options="priorityOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- Dialog: submit for approval -->
    <Dialog v-model:visible="submitDialogVisible" :header="t('training.request_submit')" modal :style="{ width: '440px' }">
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('training.request_submit_confirm') }}</p>
      <p class="text-xs text-gray-400 mt-2">{{ t('training.request_flow_hint') }}</p>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="submitDialogVisible = false" />
          <Button :label="t('training.request_submit')" icon="pi pi-send" size="small" :loading="submitting" :disabled="submitting" @click="handleSubmit" />
        </div>
      </template>
    </Dialog>

    <!-- Dialog: cancel confirmation -->
    <Dialog v-model:visible="cancelDialogVisible" :header="t('training.request_cancel')" modal :style="{ width: '440px' }">
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('training.request_cancel_confirm') }}</p>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.no')" severity="secondary" outlined size="small" @click="cancelDialogVisible = false" />
          <Button :label="t('common.yes')" severity="danger" size="small" :loading="cancelling" :disabled="cancelling" @click="handleCancel" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const statusFilter = ref(null)

const employees = ref([])
const courses = ref([])
const sessions = ref([])

const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const submitDialogVisible = ref(false)
const submitting = ref(false)
const submitTarget = ref(null)
const cancelDialogVisible = ref(false)
const cancelling = ref(false)
const cancelTarget = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-28', headerWidth: 'w-24' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'icons', count: 3, headerWidth: 'w-24' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const courseOptions = computed(() => courses.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))
const sessionOptions = computed(() => sessions.value.map(s => ({ label: `${s.session_code} — ${courseName(s.course_id)}`, value: s.id })))
const priorityOptions = computed(() => ['LOW', 'MEDIUM', 'HIGH', 'URGENT'].map(v => ({ label: priorityLabel(v), value: v })))
const statusOptions = computed(() => ['DRAFT', 'SUBMITTED', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED', 'CANCELLED'].map(v => ({ label: statusLabel(v), value: v })))

function priorityLabel(p) {
  const key = `training.priority_${String(p || '').toLowerCase()}`
  return t(key) !== key ? t(key) : p
}
function prioritySeverity(p) {
  switch (p) {
    case 'URGENT': return 'danger'
    case 'HIGH': return 'warning'
    case 'MEDIUM': return 'info'
    default: return 'secondary'
  }
}
function statusLabel(s) {
  const key = `training.request_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusSeverity(s) {
  switch (s) {
    case 'APPROVED': return 'success'
    case 'PENDING_APPROVAL': return 'warning'
    case 'SUBMITTED': return 'info'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || (id ? id : '-')
}
function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || (id ? id : '-')
}

function defaultForm() {
  return { employee_id: null, course_id: null, session_id: null, requested_date: null, reason: '', priority: 'MEDIUM' }
}

async function loadReferences() {
  const [eRes, cRes, sRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/sessions', { params: { per_page: 500 } })
  ])
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
  sessions.value = sRes.status === 'fulfilled' ? (sRes.value.data?.data || []) : []
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/trainings/requests', { params })
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

function openDialog() {
  errors.value = {}
  form.value = defaultForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = defaultForm()
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  if (!form.value.employee_id) { errors.value = { employee_id: t('form.required') }; return }
  if (!form.value.course_id) { errors.value = { course_id: t('form.required') }; return }
  if (!form.value.requested_date) { errors.value = { requested_date: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      employee_id: form.value.employee_id,
      course_id: form.value.course_id,
      session_id: form.value.session_id || null,
      requested_date: form.value.requested_date,
      reason: form.value.reason?.trim() || '',
      priority: form.value.priority || 'MEDIUM'
    }
    await api.post('/api/v1/tenant/trainings/requests', payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmSubmit(item) {
  submitTarget.value = item
  submitDialogVisible.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    await api.post(`/api/v1/tenant/trainings/requests/${submitTarget.value.id}/submit`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('training.request_submit'), life: 3000 })
    submitDialogVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    submitting.value = false
  }
}

function confirmCancel(item) {
  cancelTarget.value = item
  cancelDialogVisible.value = true
}

async function handleCancel() {
  cancelling.value = true
  try {
    await api.post(`/api/v1/tenant/trainings/requests/${cancelTarget.value.id}/cancel`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    cancelDialogVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    cancelling.value = false
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
