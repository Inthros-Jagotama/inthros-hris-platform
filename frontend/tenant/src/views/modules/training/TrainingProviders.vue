<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.provider_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-building text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.providers_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('training.code')" style="width:110px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="type" :header="t('training.type')" style="width:110px">
        <template #body="{data}"><Tag :value="typeLabel(data.type)" :severity="data.type === 'EXTERNAL' ? 'warning' : 'info'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="contact_name" :header="t('training.contact_name')" style="width:150px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.contact_name || '-' }}</span></template>
      </Column>
      <Column field="email" :header="t('common.email')" style="width:200px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.email || '-' }}</span></template>
      </Column>
      <Column field="phone" :header="t('training.phone')" style="width:140px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.phone || '-' }}</span></template>
      </Column>
      <Column field="is_active" :header="t('training.is_active')" style="width:90px">
        <template #body="{data}"><Tag :value="data.is_active ? t('common.yes') : t('common.no')" :severity="data.is_active ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" /></template>
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

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.provider_edit') : t('training.provider_new')" modal :style="{ width: '600px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('training.type')">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <RadioLabel v-for="opt in typeOptions" :key="opt.value" v-model="form.type" :value="opt.value" :label="opt.label" :id="'provider-type-' + opt.value" />
          </div>
        </FormRow>
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="200" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('training.contact_name')">
          <TextInput v-model="form.contact_name" maxlength="150" />
        </FormRow>
        <FormRow :label="t('common.email')">
          <TextInput v-model="form.email" maxlength="150" />
        </FormRow>
        <FormRow :label="t('training.phone')">
          <TextInput v-model="form.phone" maxlength="30" />
        </FormRow>
        <FormRow :label="t('training.website')">
          <TextInput v-model="form.website" maxlength="200" />
        </FormRow>
        <FormRow :label="t('training.address')">
          <TextInput v-model="form.address" textarea :rows="2" />
        </FormRow>
        <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
          <div>
            <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ t('training.is_active') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.provider_is_active_desc') }}</p>
          </div>
          <ToggleSwitch v-model="form.is_active" />
        </div>
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
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_provider', { name: deleteTarget?.name || '' })"
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
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import RadioLabel from '@/components/RadioLabel.vue'

const { t } = useI18n()
const toast = useToast()

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
const form = ref(defaultForm())

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'tag', width: 'w-14', headerWidth: 'w-16' },
  { type: 'text', width: 'w-36', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function typeLabel(type) {
  const key = `training.type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
const typeOptions = computed(() => ['INTERNAL', 'EXTERNAL'].map(v => ({ label: typeLabel(v), value: v })))

function defaultForm() {
  return { name: '', type: 'EXTERNAL', contact_name: '', email: '', phone: '', address: '', website: '', is_active: true }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/trainings/providers', { params })
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
        name: item.name || '',
        type: item.type || 'EXTERNAL',
        contact_name: item.contact_name || '',
        email: item.email || '',
        phone: item.phone || '',
        address: item.address || '',
        website: item.website || '',
        is_active: !!item.is_active
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
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      type: form.value.type,
      contact_name: form.value.contact_name?.trim() || '',
      email: form.value.email?.trim() || '',
      phone: form.value.phone?.trim() || '',
      address: form.value.address?.trim() || '',
      website: form.value.website?.trim() || '',
      is_active: form.value.is_active
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/providers/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/providers', payload)
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
    await api.delete(`/api/v1/tenant/trainings/providers/${deleteTarget.value.id}`)
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
