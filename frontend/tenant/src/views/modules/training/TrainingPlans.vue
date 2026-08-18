<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <div class="w-44">
          <InputYear v-model="yearFilter" :placeholder="t('training.filter_all_years')" showClear @update:modelValue="onFilterChange" />
        </div>
        <SelectLabel v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_status')" class="!w-44" showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.plan_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-calendar-plus text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.plans_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('training.plan_code')" style="width:130px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('training.plan_name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="year" :header="t('training.plan_year')" style="width:90px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.year }}</span></template>
      </Column>
      <Column field="status" :header="t('training.plan_status')" style="width:120px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Dialog: plan -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.plan_edit') : t('training.plan_new')" modal :style="{ width: '520px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('training.plan_year')" required :errors="errors?.year">
          <InputYear v-model="form.year" :min-year="2000" :max-year="2100" />
        </FormRow>
        <FormRow :label="t('training.plan_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="200" :placeholder="t('training.plan_name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('training.plan_status')">
          <div class="flex flex-wrap gap-2">
            <RadioLabel v-for="opt in statusOptions" :key="opt.value" v-model="form.status" :value="opt.value" :label="opt.label" :id="'plan-status-' + opt.value" />
          </div>
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
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_plan', { name: deleteTarget?.name || '' })"
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
import InputYear from '@/components/InputYear.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import RadioLabel from '@/components/RadioLabel.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const yearFilter = ref(null)
const statusFilter = ref(null)
const currentYear = new Date().getFullYear()

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
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const statusOptions = computed(() => ['DRAFT', 'ACTIVE', 'ARCHIVED'].map(v => ({ label: statusLabel(v), value: v })))

function statusLabel(s) {
  const key = `training.plan_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusSeverity(s) {
  switch (s) {
    case 'ACTIVE': return 'success'
    case 'DRAFT': return 'secondary'
    case 'ARCHIVED': return 'warning'
    default: return 'secondary'
  }
}
function defaultForm() {
  return { name: '', year: currentYear, description: '', status: 'DRAFT' }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (yearFilter.value) params.year = yearFilter.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/trainings/plans', { params })
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? { name: item.name || '', year: item.year || currentYear, description: item.description || '', status: item.status || 'DRAFT' }
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
  if (!form.value.year) { errors.value = { year: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      year: form.value.year,
      description: form.value.description?.trim() || '',
      status: form.value.status || 'DRAFT'
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/plans/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/plans', payload)
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
    await api.delete(`/api/v1/tenant/trainings/plans/${deleteTarget.value.id}`)
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
