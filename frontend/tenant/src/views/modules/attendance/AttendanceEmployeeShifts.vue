<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button v-if="hasPermission('attendance.create')" :label="t('common.add')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <p class="text-sm font-medium">{{ t('attendance.employee_shifts_empty') }}</p>
        </div>
      </template>
      <Column field="employee_id" :header="t('employee.title')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="attendance_shift_id" :header="t('attendance.shift')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ shiftName(data.attendance_shift_id) }}</span></template>
      </Column>
      <Column field="effective_date_from" :header="t('attendance.effective_date_from')" style="width:130px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.effective_date_from }}</span></template>
      </Column>
      <Column field="effective_date_to" :header="t('attendance.effective_date_to')" style="width:130px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.effective_date_to || '-' }}</span></template>
      </Column>
      <Column field="is_day_off" :header="t('attendance.is_day_off')" style="width:100px">
        <template #body="{data}"><Tag v-if="data.is_day_off" :value="t('common.yes')" severity="warn" class="!text-xs !px-1.5 !py-0.5" /><span v-else class="text-gray-400">-</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button v-if="hasPermission('attendance.update')" icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button v-if="hasPermission('attendance.delete')" icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('common.edit') : t('common.add')" modal :style="{ width: '520px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('employee.title')" required :errors="errors?.employee_id">
          <Select v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :disabled="editing" :placeholder="t('attendance.select_employee')" />
        </FormRow>
        <FormRow :label="t('attendance.shift')" required :errors="errors?.attendance_shift_id">
          <Select v-model="form.attendance_shift_id" :options="shiftOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :placeholder="t('attendance.select_shift')" />
        </FormRow>
        <FormRow :label="t('attendance.effective_date_from')" required :errors="errors?.effective_date_from">
          <DateInput v-model="form.effective_date_from" />
        </FormRow>
        <FormRow :label="t('attendance.effective_date_to')" :errors="errors?.effective_date_to">
          <DateInput v-model="form.effective_date_to" showClear />
        </FormRow>
        <FormRow :label="t('attendance.is_day_off')">
          <ToggleSwitch v-model="form.is_day_off" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="deleting"
      :error="deleteError"
      :title="t('attendance.confirm_delete_title')"
      :message="t('attendance.confirm_delete_employee_shift')"
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const employees = ref([])
const shifts = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ employee_id: '', attendance_shift_id: '', effective_date_from: '', effective_date_to: '', is_day_off: false })

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const shiftOptions = computed(() => shifts.value.map(s => ({ label: s.shift_name, value: s.id })))
const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-14', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

function employeeName(id) {
  const e = employees.value.find(x => x.id === id)
  return e ? `${e.name} (${e.employee_id})` : id
}
function shiftName(id) {
  return shifts.value.find(s => s.id === id)?.shift_name || id
}

async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}

async function loadShifts() {
  try {
    const res = await api.get('/api/v1/tenant/attendance/shifts', { params: { per_page: 500 } })
    shifts.value = res.data?.data || []
  } catch {
    shifts.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/attendance/employee-shifts', { params })
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = {
    employee_id: item?.employee_id || '',
    attendance_shift_id: item?.attendance_shift_id || '',
    effective_date_from: item?.effective_date_from || '',
    effective_date_to: item?.effective_date_to || '',
    is_day_off: !!item?.is_day_off
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { employee_id: '', attendance_shift_id: '', effective_date_from: '', effective_date_to: '', is_day_off: false }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.employee_id) { errors.value = { employee_id: t('form.required') }; return }
  if (!form.value.attendance_shift_id) { errors.value = { attendance_shift_id: t('form.required') }; return }
  if (!form.value.effective_date_from) { errors.value = { effective_date_from: t('form.required') }; return }
  saving.value = true
  try {
    if (editing.value) {
      await api.put(`/api/v1/tenant/attendance/employee-shifts/${editingId.value}`, {
        attendance_shift_id: form.value.attendance_shift_id,
        effective_date_from: form.value.effective_date_from,
        effective_date_to: form.value.effective_date_to || null,
        is_day_off: form.value.is_day_off
      })
    } else {
      await api.post('/api/v1/tenant/attendance/employee-shifts', {
        employee_id: form.value.employee_id,
        attendance_shift_id: form.value.attendance_shift_id,
        effective_date_from: form.value.effective_date_from,
        effective_date_to: form.value.effective_date_to || null,
        is_day_off: form.value.is_day_off
      })
    }
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

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/attendance/employee-shifts/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadEmployees()
  loadShifts()
  loadData()
})
</script>
