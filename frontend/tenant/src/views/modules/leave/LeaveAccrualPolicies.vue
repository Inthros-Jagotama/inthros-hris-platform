<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
        {{ totalRecords }} {{ t('common.items') }}
      </span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('leave.accrual_policies_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-percentage text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('leave.accrual_policies_empty') }}</p>
        </div>
      </template>
      <Column :header="t('leave.leave_type')" style="width:160px">
        <template #body="{data}">
          <span class="text-navy-800 dark:text-gray-100 font-medium">{{ leaveTypeName(data.leave_type_id) }}</span>
        </template>
      </Column>
      <Column field="base_quota_days" :header="t('leave.base_quota_days')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDays(data.base_quota_days) }}</span></template>
      </Column>
      <Column field="extra_every_years" :header="t('leave.extra_every_years')" style="width:100px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.extra_every_years || '-' }}</span></template>
      </Column>
      <Column field="extra_days" :header="t('leave.extra_days')" style="width:90px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.extra_days || '-' }}</span></template>
      </Column>
      <Column field="max_extra_days" :header="t('leave.max_extra_days')" style="width:110px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.max_extra_days ?? '-' }}</span></template>
      </Column>
      <Column field="effective_from" :header="t('leave.effective_from')" style="width:150px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.effective_from, locale) }}</span></template>
      </Column>
      <Column field="effective_to" :header="t('leave.effective_to')" style="width:150px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.effective_to ? formatDate(data.effective_to, locale) : '-' }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('leave.accrual_policies_edit') : t('leave.accrual_policies_new')" modal :style="{ width: '560px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('leave.leave_type')" required :errors="errors?.leave_type_id">
          <SelectLabel v-model="form.leave_type_id" :options="leaveTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.leave_type_id }" />
        </FormRow>
        <FormRow :label="t('leave.base_quota_days')" required :errors="errors?.base_quota_days">
          <InputNumber v-model="form.base_quota_days" class="!w-full" :min="0" :maxFractionDigits="2" size="small" />
        </FormRow>
        <FormRow :label="t('leave.extra_every_years')">
          <InputNumber v-model="form.extra_every_years" class="!w-full" :min="0" size="small" />
        </FormRow>
        <FormRow :label="t('leave.extra_days')">
          <InputNumber v-model="form.extra_days" class="!w-full" :min="0" :maxFractionDigits="2" size="small" />
        </FormRow>
        <FormRow :label="t('leave.max_extra_days')">
          <InputNumber v-model="form.max_extra_days" class="!w-full" :min="0" :maxFractionDigits="2" size="small" />
        </FormRow>
        <FormRow :label="t('leave.effective_from')" required :errors="errors?.effective_from">
          <DateInput v-model="form.effective_from" :maxDate="form.effective_to || null" />
        </FormRow>
        <FormRow :label="t('leave.effective_to')">
          <DateInput v-model="form.effective_to" :minDate="form.effective_from || null" showClear />
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
      :title="t('leave.confirm_delete_policy_title')"
      :message="t('leave.confirm_delete_policy')"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t, locale } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const leaveTypes = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'text', width: 'w-10', headerWidth: 'w-16' },
  { type: 'text', width: 'w-10', headerWidth: 'w-12' },
  { type: 'text', width: 'w-10', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const leaveTypeOptions = computed(() =>
  leaveTypes.value.map(lt => ({ label: lt.name, value: lt.id }))
)

function defaultForm() {
  return {
    leave_type_id: '',
    base_quota_days: null,
    extra_every_years: null,
    extra_days: null,
    max_extra_days: null,
    effective_from: '',
    effective_to: ''
  }
}

function formatDays(n) {
  const num = Number(n) || 0
  return Number.isInteger(num) ? String(num) : num.toFixed(2).replace(/0+$/, '').replace(/\.$/, '')
}

function leaveTypeName(id) {
  return leaveTypes.value.find(x => x.id === id)?.name || '—'
}

async function loadLeaveTypes() {
  try {
    const res = await api.get('/api/v1/tenant/leave/types', { params: { page: 1, per_page: 100 } })
    leaveTypes.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/leave/accrual-policies', { params })
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
  form.value = item
    ? {
        leave_type_id: item.leave_type_id || '',
        base_quota_days: item.base_quota_days ?? null,
        extra_every_years: item.extra_every_years ?? null,
        extra_days: item.extra_days ?? null,
        max_extra_days: item.max_extra_days ?? null,
        effective_from: item.effective_from || '',
        effective_to: item.effective_to || ''
      }
    : defaultForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = defaultForm()
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.leave_type_id) { errors.value = { leave_type_id: t('form.required') }; return }
  if (form.value.base_quota_days === null || form.value.base_quota_days === undefined) { errors.value = { base_quota_days: t('form.required') }; return }
  if (!form.value.effective_from) { errors.value = { effective_from: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      leave_type_id: form.value.leave_type_id,
      base_quota_days: form.value.base_quota_days,
      extra_every_years: form.value.extra_every_years,
      extra_days: form.value.extra_days,
      max_extra_days: form.value.max_extra_days,
      effective_from: form.value.effective_from,
      effective_to: form.value.effective_to || null
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/leave/accrual-policies/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/leave/accrual-policies', payload)
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
    await api.delete(`/api/v1/tenant/leave/accrual-policies/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await loadLeaveTypes()
  await loadData()
})
</script>
