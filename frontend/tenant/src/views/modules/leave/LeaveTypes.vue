<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
        {{ totalRecords }} {{ t('common.items') }}
      </span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('leave.types_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-tags text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('leave.types_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('leave.code')" style="width:120px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="quota_period" :header="t('leave.quota_period')" style="width:110px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ t('leave.quota_period_' + (data.quota_period || 'none').toLowerCase()) }}</span></template>
      </Column>
      <Column field="default_quota_days" :header="t('leave.default_quota_days')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.default_quota_days ?? '-' }}</span></template>
      </Column>
      <Column field="is_paid" :header="t('leave.is_paid')" style="width:90px">
        <template #body="{data}"><Tag :value="data.is_paid ? t('common.yes') : t('common.no')" :severity="data.is_paid ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="is_active" :header="t('leave.is_active')" style="width:90px">
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

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('leave.types_edit') : t('leave.types_new')" modal :style="{ width: '560px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('leave.code')" :errors="errors?.code">
          <TextInput v-model="form.code" maxlength="50" readonly :placeholder="t('leave.code')" :class="{ 'p-invalid': errors?.code }" />
          <p class="text-xs text-gray-400 mt-1">{{ t('leave.code_auto_hint') }}</p>
        </FormRow>
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="150" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('leave.quota_period')">
          <SelectLabel v-model="form.quota_period" :options="quotaPeriodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('leave.default_quota_days')" :errors="errors?.default_quota_days">
          <InputNumber v-model="form.default_quota_days" class="!w-full" :min="0" size="small" />
        </FormRow>
        <FormRow :label="t('leave.is_paid')">
          <ToggleSwitch v-model="form.is_paid" />
        </FormRow>
        <FormRow :label="t('leave.counts_against_quota')">
          <ToggleSwitch v-model="form.counts_against_quota" />
        </FormRow>
        <FormRow :label="t('leave.requires_attachment')">
          <ToggleSwitch v-model="form.requires_attachment" />
        </FormRow>
        <FormRow :label="t('leave.is_active')">
          <ToggleSwitch v-model="form.is_active" />
        </FormRow>
        <FormRow :label="t('leave.allow_half_day')">
          <ToggleSwitch v-model="form.allow_half_day" />
        </FormRow>
        <FormRow :label="t('leave.allow_hourly')">
          <ToggleSwitch v-model="form.allow_hourly" />
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
      :title="t('leave.confirm_delete_title')"
      :message="t('leave.confirm_delete', { name: deleteTarget?.name || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
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
import SelectLabel from '@/components/SelectLabel.vue'

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
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-36', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const quotaPeriodOptions = computed(() => [
  { label: t('leave.quota_period_year'), value: 'YEAR' },
  { label: t('leave.quota_period_month'), value: 'MONTH' },
  { label: t('leave.quota_period_none'), value: 'NONE' }
])

function defaultForm() {
  return {
    code: '',
    name: '',
    description: '',
    is_paid: true,
    requires_attachment: false,
    allow_half_day: true,
    allow_hourly: false,
    counts_against_quota: true,
    default_quota_days: null,
    quota_period: 'YEAR',
    is_active: true
  }
}

// Slug code otomatis dari nama (hanya saat create; saat edit code dipertahankan)
function slugify(text) {
  return (text || '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/[\s_-]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

watch(
  () => form.value.name,
  (name) => {
    if (!editing.value) form.value.code = slugify(name)
  }
)

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/leave/types', { params })
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
        code: item.code || '',
        name: item.name || '',
        description: item.description || '',
        is_paid: !!item.is_paid,
        requires_attachment: !!item.requires_attachment,
        allow_half_day: !!item.allow_half_day,
        allow_hourly: !!item.allow_hourly,
        counts_against_quota: !!item.counts_against_quota,
        default_quota_days: item.default_quota_days ?? null,
        quota_period: item.quota_period || 'YEAR',
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
      code: form.value.code?.trim() || '',
      name: form.value.name.trim(),
      description: form.value.description?.trim() || '',
      is_paid: form.value.is_paid,
      requires_attachment: form.value.requires_attachment,
      allow_half_day: form.value.allow_half_day,
      allow_hourly: form.value.allow_hourly,
      counts_against_quota: form.value.counts_against_quota,
      default_quota_days: form.value.default_quota_days,
      quota_period: form.value.quota_period,
      is_active: form.value.is_active
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/leave/types/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/leave/types', payload)
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
    await api.delete(`/api/v1/tenant/leave/types/${deleteTarget.value.id}`)
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
