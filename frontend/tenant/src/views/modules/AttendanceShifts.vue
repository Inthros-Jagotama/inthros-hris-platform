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
          <i class="pi pi-clock text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.shifts_empty') }}</p>
        </div>
      </template>
      <Column field="shift_name" :header="t('attendance.shift_name')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.shift_name }}</span></template>
      </Column>
      <Column field="check_in_time" :header="t('attendance.check_in_time')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.check_in_time }}</span></template>
      </Column>
      <Column field="check_out_time" :header="t('attendance.check_out_time')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.check_out_time }}</span></template>
      </Column>
      <Column field="is_cross_midnight" :header="t('attendance.is_cross_midnight')" style="width:140px">
        <template #body="{data}"><Tag v-if="data.is_cross_midnight" :value="t('common.yes')" severity="info" class="!text-xs !px-1.5 !py-0.5" /><span v-else class="text-gray-400">-</span></template>
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

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('common.edit') : t('common.add')" modal :style="{ width: '480px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('attendance.shift_name')" required :errors="errors?.shift_name">
          <TextInput v-model="form.shift_name" maxlength="255" autofocus />
        </FormRow>
        <FormRow :label="t('attendance.check_in_time')" required :errors="errors?.check_in_time">
          <InputText v-model="form.check_in_time" placeholder="08:00:00" size="small" class="w-full" />
        </FormRow>
        <FormRow :label="t('attendance.check_out_time')" required :errors="errors?.check_out_time">
          <InputText v-model="form.check_out_time" placeholder="17:00:00" size="small" class="w-full" />
        </FormRow>
        <FormRow :label="t('attendance.is_cross_midnight')">
          <ToggleSwitch v-model="form.is_cross_midnight" />
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
      :message="t('attendance.confirm_delete_shift', { name: deleteTarget?.shift_name || '' })"
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
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ shift_name: '', check_in_time: '', check_out_time: '', is_cross_midnight: false })

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/attendance/shifts', { params })
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
    shift_name: item?.shift_name || '',
    check_in_time: item?.check_in_time || '',
    check_out_time: item?.check_out_time || '',
    is_cross_midnight: !!item?.is_cross_midnight
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { shift_name: '', check_in_time: '', check_out_time: '', is_cross_midnight: false }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.shift_name?.trim()) { errors.value = { shift_name: t('form.required') }; return }
  if (!form.value.check_in_time?.trim()) { errors.value = { check_in_time: t('form.required') }; return }
  if (!form.value.check_out_time?.trim()) { errors.value = { check_out_time: t('form.required') }; return }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (editing.value) {
      await api.put(`/api/v1/tenant/attendance/shifts/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/attendance/shifts', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    }
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
    await api.delete(`/api/v1/tenant/attendance/shifts/${deleteTarget.value.id}`)
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
