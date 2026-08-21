<template>
  <div class="space-y-4">
    <!-- ── Toolbar: filter period + organisasi ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <InputText
          v-model="period"
          :placeholder="t('headcount_planning.period_placeholder')"
          class="!w-32 !text-sm"
          @keyup.enter="loadAll"
        />
        <Select
          v-model="orgFilter"
          :options="organizationOptions"
          optionLabel="label"
          optionValue="value"
          filter
          showClear
          class="!w-56 !text-sm"
          :placeholder="t('headcount_planning.filter_by_organization')"
          @change="loadPlans"
        />
        <Button icon="pi pi-refresh" size="small" severity="secondary" outlined v-tooltip.bottom="t('common.refresh')" @click="loadAll" />
      </div>
      <Button :label="t('headcount_planning.add_plan')" icon="pi pi-plus" size="small" @click="openDialog()" />
    </div>

    <!-- ── Gap Analysis ── -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-3">
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('headcount_planning.supply') }}</p>
        <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ gap?.supply ?? '—' }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('headcount_planning.demand') }}</p>
        <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ gap?.demand ?? '—' }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('headcount_planning.gap') }}</p>
        <p class="text-lg font-bold" :class="gapColorClass">{{ gap?.gap ?? '—' }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('headcount_planning.status') }}</p>
        <Tag v-if="gap?.status" :value="gap.status" :severity="statusSeverity(gap.status)" class="!text-xs !px-1.5 !py-0.5" />
        <p v-else class="text-lg font-bold text-gray-300">—</p>
      </div>
    </div>

    <div v-if="gap?.departments?.length" class="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div class="px-4 py-2 border-b border-gray-200 dark:border-gray-700">
        <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">{{ t('headcount_planning.department_breakdown') }}</p>
      </div>
      <DataTable :value="gap.departments" size="small" class="!text-sm">
        <Column field="organization_name" :header="t('headcount_planning.organization')" />
        <Column field="supply" :header="t('headcount_planning.supply')" style="width:100px" />
        <Column field="demand" :header="t('headcount_planning.demand')" style="width:100px" />
        <Column field="gap" :header="t('headcount_planning.gap')" style="width:100px" />
        <Column :header="t('headcount_planning.status')" style="width:120px">
          <template #body="{data}"><Tag :value="data.status" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>
    </div>

    <!-- ── Daftar Headcount Plan ── -->
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
          <i class="pi pi-arrows-alt text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('headcount_planning.empty_plans') }}</p>
        </div>
      </template>
      <Column field="period" :header="t('headcount_planning.period')" style="width:100px" />
      <Column :header="t('headcount_planning.organization')" style="width:220px">
        <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ organizationName(data.organization_id) }}</span></template>
      </Column>
      <Column field="planned_hc" :header="t('headcount_planning.planned_hc')" style="width:110px" />
      <Column field="actual_hc" :header="t('headcount_planning.actual_hc')" style="width:110px" />
      <Column :header="t('headcount_planning.variance')" style="width:100px">
        <template #body="{data}">
          <span :class="data.variance < 0 ? 'text-rose-600 dark:text-rose-400' : data.variance > 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-500'">{{ data.variance }}</span>
        </template>
      </Column>
      <Column field="snapshot_date" :header="t('headcount_planning.snapshot_date')" style="width:120px" />
      <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Tambah / Ubah ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingId ? t('headcount_planning.edit_plan') : t('headcount_planning.add_plan')"
      modal
      :style="{ width: '440px' }"
      @hide="resetForm"
    >
      <div class="space-y-3">
        <FormRow :label="t('headcount_planning.period')" required :errors="errors?.period">
          <InputText v-model="form.period" :disabled="!!editingId" :placeholder="t('headcount_planning.period_placeholder')" class="w-full !text-sm" />
        </FormRow>
        <FormRow :label="t('headcount_planning.organization')" required :errors="errors?.organization_id">
          <Select
            v-model="form.organization_id"
            :options="organizationOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            :disabled="!!editingId"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('headcount_planning.planned_hc')" required :errors="errors?.planned_hc">
          <InputNumber v-model="form.planned_hc" :min="0" :useGrouping="false" inputClass="!text-sm" class="w-full" />
        </FormRow>
        <FormRow :label="t('headcount_planning.snapshot_date')" :errors="errors?.snapshot_date">
          <DatePicker v-model="form.snapshot_date" dateFormat="yy-mm-dd" :placeholder="t('common.select_date')" class="w-full" showIcon />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Konfirmasi hapus ── -->
    <ConfirmDeleteDialog
      v-model:visible="deleteConfirmVisible"
      :title="t('headcount_planning.confirm_delete_title')"
      :message="t('headcount_planning.confirm_delete_msg')"
      :loading="actionLoading"
      :error-msg="actionError"
      :cancel-label="t('common.no')"
      :confirm-label="t('common.delete')"
      @confirm="handleDeleteConfirm"
      @cancel="deleteConfirmVisible = false"
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
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t } = useI18n()
const toast = useToast()

function currentPeriod() {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
}

// ── State ──
const period = ref(currentPeriod())
const orgFilter = ref(null)
const loading = ref(false)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const gap = ref(null)
const organizations = ref([])

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())

// ── Konfirmasi hapus ──
const deleteTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const deleteConfirmVisible = ref(false)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const organizationOptions = computed(() =>
  organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id }))
)

const gapColorClass = computed(() => {
  const g = gap.value?.gap ?? 0
  if (g > 0) return 'text-rose-600 dark:text-rose-400'
  if (g < 0) return 'text-emerald-600 dark:text-emerald-400'
  return 'text-navy-800 dark:text-gray-100'
})

const skeletonColumns = [
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-24' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

function emptyForm() {
  return { period: period.value, organization_id: null, planned_hc: null, snapshot_date: null }
}

function organizationName(id) {
  return organizations.value.find(o => o.id === id)?.nomenclature || '—'
}
function statusSeverity(status) {
  if (status === 'SHORTAGE') return 'danger'
  if (status === 'SURPLUS') return 'info'
  return 'success'
}

// ── Load data ──
async function loadPlans() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (period.value) params.period = period.value
    if (orgFilter.value) params.organization_id = orgFilter.value
    const res = await api.get('/api/v1/tenant/workforce-intelligence/planning/headcounts', { params })
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

async function loadGapAnalysis() {
  try {
    const res = await api.get('/api/v1/tenant/workforce-intelligence/planning/gap-analysis', { params: { period: period.value } })
    gap.value = res.data?.data || null
  } catch (e) {
    gap.value = null
  }
}

async function loadAll() {
  currentPage.value = 1
  await Promise.all([loadPlans(), loadGapAnalysis()])
}

async function loadReferences() {
  const orgRes = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
  organizations.value = orgRes.data?.data || []
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadPlans()
}

// ── Dialog create/edit ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      period: data.period,
      organization_id: data.organization_id,
      planned_hc: data.planned_hc,
      snapshot_date: data.snapshot_date ? new Date(data.snapshot_date) : null
    }
  } else {
    form.value = emptyForm()
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = emptyForm()
  errors.value = {}
  editingId.value = null
}

function validateForm() {
  const e = {}
  if (!editingId.value && !form.value.period?.trim()) e.period = t('career_paths.field_required')
  if (!editingId.value && !form.value.organization_id) e.organization_id = t('career_paths.field_required')
  if (form.value.planned_hc === null || form.value.planned_hc === undefined) e.planned_hc = t('career_paths.field_required')
  return e
}

function formatDateForApi(date) {
  if (!date) return undefined
  if (typeof date === 'string') return date
  return new Date(date).toISOString().split('T')[0]
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    if (editingId.value) {
      const payload = {
        planned_hc: form.value.planned_hc,
        snapshot_date: formatDateForApi(form.value.snapshot_date)
      }
      await api.put(`/api/v1/tenant/workforce-intelligence/planning/headcounts/${editingId.value}`, payload)
    } else {
      const payload = {
        period: form.value.period.trim(),
        organization_id: form.value.organization_id,
        planned_hc: form.value.planned_hc,
        snapshot_date: formatDateForApi(form.value.snapshot_date)
      }
      await api.post('/api/v1/tenant/workforce-intelligence/planning/headcounts', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadAll()
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

// ── Hapus ──
function openDeleteConfirm(data) {
  deleteTarget.value = data
  actionError.value = ''
  deleteConfirmVisible.value = true
}

function handleDeleteConfirm() {
  if (!deleteTarget.value?.id) return
  actionError.value = ''
  actionLoading.value = true
  api.delete(`/api/v1/tenant/workforce-intelligence/planning/headcounts/${deleteTarget.value.id}`)
    .then(async () => {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      deleteTarget.value = null
      await loadAll()
    })
    .catch(e => {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    })
    .finally(() => { actionLoading.value = false })
}

onMounted(() => {
  loadAll()
  loadReferences()
})
</script>
