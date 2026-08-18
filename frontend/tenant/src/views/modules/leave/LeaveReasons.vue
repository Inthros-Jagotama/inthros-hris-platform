<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="items.length > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
        {{ items.length }} {{ t('common.items') }}
      </span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('leave.reasons_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />
    <DataTable
      v-else
      :value="items"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      sortField="sort_order"
      :sortOrder="1"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-list text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('leave.reasons_empty') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="is_other" :header="t('leave.is_other')" style="width:120px">
        <template #body="{data}"><Tag :value="data.is_other ? t('common.yes') : t('common.no')" :severity="data.is_other ? 'info' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="sort_order" :header="t('leave.sort_order')" sortable style="width:110px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span></template>
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

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('leave.reasons_edit') : t('leave.reasons_new')" modal :style="{ width: '460px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="150" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('leave.sort_order')" :errors="errors?.sort_order">
          <InputNumber v-model="form.sort_order" class="!w-full" :min="0" size="small" />
        </FormRow>
        <FormRow :label="t('leave.is_other')">
          <ToggleSwitch v-model="form.is_other" />
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
      :title="t('leave.confirm_delete_reason_title')"
      :message="t('leave.confirm_delete_reason', { name: deleteTarget?.name || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)

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
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'text', width: 'w-10', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

function defaultForm() {
  return { name: '', is_other: false, sort_order: 0 }
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/leave/reasons')
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? { name: item.name || '', is_other: !!item.is_other, sort_order: item.sort_order || 0 }
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
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      is_other: form.value.is_other,
      sort_order: form.value.sort_order || 0
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/leave/reasons/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/leave/reasons', payload)
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
    await api.delete(`/api/v1/tenant/leave/reasons/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
