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
          <i class="pi pi-shield text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.exempt_positions_empty') }}</p>
        </div>
      </template>
      <Column field="organization_id" :header="t('attendance.organization')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ organizationName(data.organization_id) }}</span></template>
      </Column>
      <Column field="is_exempt" :header="t('attendance.is_exempt')" style="width:120px">
        <template #body="{data}"><Tag :value="data.is_exempt ? t('common.yes') : t('common.no')" :severity="data.is_exempt ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="note" :header="t('attendance.note')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.note || '-' }}</span></template>
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
        <FormRow :label="t('attendance.organization')" required :errors="errors?.organization_id">
          <Select v-model="form.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :disabled="editing" :placeholder="t('attendance.select_organization')" />
        </FormRow>
        <FormRow :label="t('attendance.is_exempt')">
          <ToggleSwitch v-model="form.is_exempt" />
        </FormRow>
        <FormRow :label="t('attendance.note')" :errors="errors?.note">
          <TextInput v-model="form.note" textarea :rows="3" />
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
      :message="t('attendance.confirm_delete_exempt_position')"
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
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const organizations = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ organization_id: '', is_exempt: true, note: '' })

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.name, value: o.id })))
const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

function organizationName(id) {
  return organizations.value.find(o => o.id === id)?.name || id
}

async function loadOrganizations() {
  try {
    const res = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500 } })
    organizations.value = res.data?.data || []
  } catch {
    organizations.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/attendance/exempt-positions', { params })
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
    organization_id: item?.organization_id || '',
    is_exempt: item ? !!item.is_exempt : true,
    note: item?.note || ''
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { organization_id: '', is_exempt: true, note: '' }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.organization_id) { errors.value = { organization_id: t('form.required') }; return }
  saving.value = true
  try {
    if (editing.value) {
      await api.put(`/api/v1/tenant/attendance/exempt-positions/${editingId.value}`, { is_exempt: form.value.is_exempt, note: form.value.note })
    } else {
      await api.post('/api/v1/tenant/attendance/exempt-positions', { organization_id: form.value.organization_id, is_exempt: form.value.is_exempt, note: form.value.note })
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
    await api.delete(`/api/v1/tenant/attendance/exempt-positions/${deleteTarget.value.id}`)
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
  loadOrganizations()
  loadData()
})
</script>
