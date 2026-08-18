<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('performance_periods.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable v-else :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="year" :sortOrder="-1">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('performance_periods.empty_title') }}</p>
        </div>
      </template>
      <Column field="period_code" :header="t('performance_periods.code')" sortable style="width:120px">
        <template #body="{data}"><Tag :value="data.period_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="period_type" :header="t('performance_periods.type')" sortable style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.period_type }}</span></template>
      </Column>
      <Column field="year" :header="t('performance_periods.year')" sortable style="width:80px">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-semibold">{{ data.year }}</span></template>
      </Column>
      <Column field="start_date" :header="t('performance_periods.start_date')" sortable style="width:120px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.start_date || '-' }}</span></template>
      </Column>
      <Column field="end_date" :header="t('performance_periods.end_date')" sortable style="width:120px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.end_date || '-' }}</span></template>
      </Column>
      <Column field="status" :header="t('performance_periods.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" class="!text-xs" />
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:140px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-calculator" size="small" text severity="info" v-tooltip.left="t('performance_periods.recalculate_scoring')" @click="confirmRecalculate(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('performance_periods.edit') : t('performance_periods.new')" modal :style="{ width: '420px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('performance_periods.code')" required :errors="errors?.period_code">
          <TextInput v-model="form.period_code" maxlength="10" autofocus :placeholder="t('performance_periods.code_placeholder')" :class="{'p-invalid':errors?.period_code}" />
        </FormRow>
        <FormRow :label="t('performance_periods.type')" required :errors="errors?.period_type">
          <Select v-model="form.period_type" :options="periodTypes" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" :class="{'p-invalid':errors?.period_type}" />
        </FormRow>
        <FormRow :label="t('performance_periods.year')" required :errors="errors?.year">
          <InputNumber v-model="form.year" class="!w-full" :min="2020" :max="2100" :useGrouping="false" :maxFractionDigits="0" size="small" />
        </FormRow>
        <FormRow :label="t('performance_periods.start_date')" :errors="errors?.start_date">
          <DatePicker v-model="form.start_date" dateFormat="yy-mm-dd" :placeholder="t('common.select_date')" class="w-full" showIcon />
        </FormRow>
        <FormRow :label="t('performance_periods.end_date')" :errors="errors?.end_date">
          <DatePicker v-model="form.end_date" dateFormat="yy-mm-dd" :placeholder="t('common.select_date')" class="w-full" showIcon />
        </FormRow>
        <FormRow :label="t('performance_periods.status')" :errors="errors?.status">
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('performance_periods.delete_title')"
      :message="t('performance_periods.delete_message', { name: deleteTarget?.period_code })"
      :loading="deleting"
      :error="deleteError"
      @confirm="handleDelete"
    />

    <ConfirmActionDialog
      v-model:visible="recalculateDialogVisible"
      :title="t('performance_periods.recalculate_scoring')"
      :message="t('performance_periods.recalculate_scoring_confirm', { code: recalculateTarget?.period_code })"
      :loading="recalculating"
      icon="pi pi-calculator"
      severity="info"
      :confirm-label="t('performance_periods.recalculate_scoring')"
      :cancel-label="t('common.cancel')"
      @confirm="handleRecalculate"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()
const recalculateDialogVisible = ref(false)
const recalculating = ref(false)
const recalculateTarget = ref(null)
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
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const form = ref({
  period_code: '',
  period_type: 'ANNUAL',
  year: new Date().getFullYear(),
  start_date: null,
  end_date: null,
  status: 'draft'
})

const periodTypes = [
  { label: 'Monthly', value: 'MONTHLY' },
  { label: 'Quarterly', value: 'QUARTERLY' },
  { label: 'Semester', value: 'SEMESTER' },
  { label: 'Annual', value: 'ANNUAL' }
]

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Active', value: 'active' },
  { label: 'Closed', value: 'closed' }
]

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function getStatusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'closed': return 'secondary'
    default: return 'warn'
  }
}

function formatDate(date) {
  if (!date) return null
  if (typeof date === 'string') return date
  const d = new Date(date)
  return d.toISOString().split('T')[0]
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/periods', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
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
    period_code: item?.period_code || '',
    period_type: item?.period_type || 'ANNUAL',
    year: item?.year || new Date().getFullYear(),
    start_date: item?.start_date ? new Date(item.start_date) : null,
    end_date: item?.end_date ? new Date(item.end_date) : null,
    status: item?.status || 'draft'
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = {
    period_code: '',
    period_type: 'ANNUAL',
    year: new Date().getFullYear(),
    start_date: null,
    end_date: null,
    status: 'draft'
  }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.period_code?.trim()) {
    errors.value = { period_code: [t('form.required')] }
    return
  }
  if (!form.value.period_type) {
    errors.value = { period_type: [t('form.required')] }
    return
  }
  if (!form.value.year) {
    errors.value = { year: [t('form.required')] }
    return
  }

  saving.value = true
  try {
    const payload = {
      period_code: form.value.period_code,
      period_type: form.value.period_type,
      year: form.value.year,
      start_date: formatDate(form.value.start_date),
      end_date: formatDate(form.value.end_date),
      status: form.value.status
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/performance/periods/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_periods.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/performance/periods', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_periods.created'), life: 3000 })
    }
    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmRecalculate(item) {
  recalculateTarget.value = item
  recalculateDialogVisible.value = true
}

async function handleRecalculate() {
  recalculating.value = true
  try {
    const res = await api.post(`/api/v1/tenant/performance/kpi/periods/${recalculateTarget.value.id}/recalculate-scoring`)
    const body = res.data?.data || res.data
    recalculateDialogVisible.value = false
    toast.add({
      severity: 'success',
      summary: t('message.success'),
      detail: t('performance_periods.recalculate_scoring_done', { processed: body?.processed ?? 0, total: body?.total ?? 0 }),
      life: 4000
    })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    recalculating.value = false
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
    await api.delete(`/api/v1/tenant/performance/periods/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_periods.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
