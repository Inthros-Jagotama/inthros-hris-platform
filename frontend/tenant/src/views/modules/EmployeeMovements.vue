<template>
  <div class="space-y-1">
    <!-- ── Toolbar: filter + tombol ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Select
          v-model="filterType"
          :options="typeOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.filter_all_types')"
          class="!w-44"
          size="small"
          showClear
          @change="onFilterChange"
        />
        <Select
          v-model="filterStatus"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('employee_movement.filter_all_status')"
          class="!w-40"
          size="small"
          showClear
          @change="onFilterChange"
        />
        <InputText
          v-model="searchTerm"
          :placeholder="t('employee_movement.search_placeholder')"
          class="!w-64"
          size="small"
          @keyup.enter="onFilterChange"
        />
        <Button
          v-if="searchTerm || filterType || filterStatus"
          icon="pi pi-times"
          severity="secondary"
          outlined
          size="small"
          v-tooltip.left="t('common.reset')"
          @click="resetFilters"
        />
        <Button :label="t('employee_movement.add_movement')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-arrows-alt text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee_movement.empty') }}</p>
        </div>
      </template>
      <Column field="employee_name" :header="t('employee_movement.employee')" style="width:180px">
        <template #body="{data}">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-200 font-medium">{{ data.employee_name || '-' }}</span>
            <span v-if="data.employee_code" class="text-xs text-gray-400">{{ data.employee_code }}</span>
          </div>
        </template>
      </Column>
      <Column field="movement_type" :header="t('employee_movement.movement_type')" style="width:170px">
        <template #body="{data}">
          <Tag :value="typeLabel(data.movement_type)" :severity="typeSeverity(data.movement_type)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('employee_movement.to_position')" style="width:180px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.to_position_name || data.to_organization_name || '-' }}</span>
        </template>
      </Column>
      <Column :header="t('employee_movement.to_employment_status')" style="width:150px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.to_employment_status_name || '-' }}</span>
        </template>
      </Column>
      <Column field="decision_letter_number" :header="t('employee_movement.decision_letter_number')" style="width:130px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ data.decision_letter_number || '-' }}</span></template>
      </Column>
      <Column field="effective_date" :header="t('employee_movement.effective_date')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.effective_date, locale) }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:140px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:270px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button
              icon="pi pi-eye"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('common.view')"
              @click="openDetail(data)"
            />
            <Button
              v-if="data.status === 'draft'"
              :label="t('employee_movement.submit')"
              icon="pi pi-send"
              size="small"
              severity="primary"
              @click="openSubmitConfirm(data)"
            />
            <Button
              v-if="data.status === 'approved'"
              :label="t('employee_movement.execute')"
              icon="pi pi-check"
              size="small"
              severity="success"
              @click="openExecuteConfirm(data)"
            />
            <Button
              v-if="data.status === 'draft'"
              icon="pi pi-trash"
              size="small"
              severity="danger"
              text
              v-tooltip.left="t('common.delete')"
              @click="openDeleteConfirm(data)"
            />
            <Button
              v-if="data.status === 'pending_approval' || data.status === 'approved'"
              icon="pi pi-times"
              size="small"
              severity="danger"
              text
              v-tooltip.left="t('employee_movement.cancel')"
              @click="openCancelConfirm(data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Buat Mutasi ── -->
    <Dialog v-model:visible="dialogVisible" :header="t('employee_movement.add_movement')" modal :style="{ width: '560px' }" @hide="resetForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('employee_movement.hint_per_type') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('employee_movement.employee')" required :errors="errors?.employee_id">
          <Select v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :placeholder="t('employee_movement.select_employee')" />
        </FormRow>
        <FormRow :label="t('employee_movement.movement_type')" required :errors="errors?.movement_type">
          <Select v-model="form.movement_type" :options="typeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>

        <!-- to_* fields: tampil sesuai tipe -->
        <FormRow v-if="requiresOrganization(form.movement_type) || requiresPosition(form.movement_type)" :label="t('employee_movement.to_organization')" :required="requiresOrganization(form.movement_type)" :errors="errors?.to_organization_id">
          <Select v-model="form.to_organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" />
        </FormRow>
        <FormRow v-if="requiresPosition(form.movement_type)" :label="t('employee_movement.to_position')" required :errors="errors?.to_position_id">
          <Select v-model="form.to_position_id" :options="positionOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" />
        </FormRow>
        <FormRow v-if="form.movement_type === 'status_change'" :label="t('employee_movement.to_employment_status')" required :errors="errors?.to_employment_status_id">
          <Select v-model="form.to_employment_status_id" :options="statusOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" />
        </FormRow>

        <FormRow :label="t('employee_movement.decision_letter_number')" required :errors="errors?.decision_letter_number">
          <TextInput v-model="form.decision_letter_number" />
        </FormRow>
        <FormRow :label="t('employee_movement.decision_letter_date')" required :errors="errors?.decision_letter_date">
          <DateInput v-model="form.decision_letter_date" />
        </FormRow>
        <FormRow :label="t('employee_movement.effective_date')" required :errors="errors?.effective_date">
          <DateInput v-model="form.effective_date" />
        </FormRow>
        <FormRow :label="t('employee_movement.reason')" :errors="errors?.reason">
          <TextInput v-model="form.reason" textarea :rows="2" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Detail Movement ── -->
    <Dialog v-model:visible="detailVisible" :header="t('employee_movement.detail_title')" modal :style="{ width: '680px' }" @hide="detailItem = null">
      <div v-if="detailItem" class="space-y-4">
        <!-- Ringkasan tipe & status -->
        <div class="flex items-center justify-between gap-2 flex-wrap">
          <Tag :value="typeLabel(detailItem.movement_type)" :severity="typeSeverity(detailItem.movement_type)" class="!text-xs !px-1.5 !py-0.5" />
          <Tag :value="statusLabel(detailItem.status)" :severity="statusSeverity(detailItem.status)" class="!text-xs !px-1.5 !py-0.5" />
        </div>

        <!-- Karyawan & SK -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <ViewLabel :label="t('employee_movement.employee')">
            <span class="font-medium">{{ detailItem.employee_name || '-' }}</span>
            <span v-if="detailItem.employee_code" class="ml-1 text-xs text-gray-400">({{ detailItem.employee_code }})</span>
          </ViewLabel>
          <ViewLabel :label="t('employee_movement.decision_letter_number')" :value="detailItem.decision_letter_number || '-'" mono breakAll />
        </div>

        <!-- Dari → Ke -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
          <div v-if="hasAnyField(detailItem, 'from')" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-2">
            <p class="text-xs uppercase tracking-wide text-gray-400 font-medium">{{ t('employee_movement.from') }}</p>
            <ViewLabel :label="t('employee_movement.from_organization')" :value="detailItem.from_organization_name || '-'" />
            <ViewLabel :label="t('employee_movement.from_position')" :value="detailItem.from_position_name || '-'" />
            <ViewLabel :label="t('employee_movement.from_employment_status')" :value="detailItem.from_employment_status_name || '-'" />
          </div>
          <div v-if="hasAnyField(detailItem, 'to')" class="rounded-lg border border-emerald-200 dark:border-emerald-900/40 p-3 space-y-2">
            <p class="text-xs uppercase tracking-wide text-emerald-500 font-medium">{{ t('employee_movement.to') }}</p>
            <ViewLabel :label="t('employee_movement.to_organization')" :value="detailItem.to_organization_name || '-'" />
            <ViewLabel :label="t('employee_movement.to_position')" :value="detailItem.to_position_name || '-'" />
            <ViewLabel :label="t('employee_movement.to_employment_status')" :value="detailItem.to_employment_status_name || '-'" />
          </div>
        </div>

        <!-- Tanggal -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <ViewLabel :label="t('employee_movement.decision_letter_date')" :value="formatDate(detailItem.decision_letter_date, locale)" />
          <ViewLabel :label="t('employee_movement.effective_date')" :value="formatDate(detailItem.effective_date, locale)" />
          <ViewLabel :label="t('employee_movement.created_at')" :value="formatDateTime(detailItem.created_at)" />
        </div>

        <!-- Alasan & catatan -->
        <div v-if="detailItem.reason || detailItem.notes" class="grid grid-cols-1 gap-3">
          <ViewLabel v-if="detailItem.reason" :label="t('employee_movement.reason')" :value="detailItem.reason" />
          <ViewLabel v-if="detailItem.notes" :label="t('employee_movement.notes')" :value="detailItem.notes" />
        </div>

        <!-- Riwayat approval & eksekusi -->
        <div v-if="detailItem.approved_at || detailItem.executed_at" class="grid grid-cols-1 sm:grid-cols-2 gap-3 rounded-lg bg-gray-50 dark:bg-gray-800/40 p-3">
          <ViewLabel v-if="detailItem.approved_at" :label="t('employee_movement.approved_at')" :value="formatDateTime(detailItem.approved_at)" />
          <ViewLabel v-if="detailItem.executed_at" :label="t('employee_movement.executed_at')" :value="formatDateTime(detailItem.executed_at)" />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2 flex-wrap">
          <Button
            v-if="detailItem?.status === 'draft'"
            :label="t('employee_movement.submit')"
            icon="pi pi-send"
            size="small"
            severity="primary"
            @click="actionFromDetail('submit')"
          />
          <Button
            v-if="detailItem?.status === 'approved'"
            :label="t('employee_movement.execute')"
            icon="pi pi-check"
            size="small"
            severity="success"
            @click="actionFromDetail('execute')"
          />
          <Button
            v-if="detailItem?.status === 'pending_approval' || detailItem?.status === 'approved'"
            :label="t('employee_movement.cancel')"
            icon="pi pi-times"
            size="small"
            severity="danger"
            text
            @click="actionFromDetail('cancel')"
          />
          <Button
            v-if="detailItem?.status === 'draft'"
            :label="t('common.delete')"
            icon="pi pi-trash"
            size="small"
            severity="danger"
            text
            @click="actionFromDetail('delete')"
          />
          <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="detailVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── Konfirmasi aksi ── -->
    <ConfirmActionDialog
      v-model:visible="submitConfirmVisible"
      :title="t('employee_movement.confirm_submit_title')"
      :message="t('employee_movement.confirm_submit_msg')"
      :loading="actionLoading"
      :error-msg="actionError"
      :cancel-label="t('common.no')"
      :confirm-label="t('employee_movement.submit')"
      icon="pi pi-send"
      @confirm="handleSubmitConfirm"
      @cancel="submitConfirmVisible = false"
    />
    <ConfirmActionDialog
      v-model:visible="executeConfirmVisible"
      :title="t('employee_movement.confirm_execute_title')"
      :message="t('employee_movement.confirm_execute_msg')"
      :loading="actionLoading"
      :error-msg="actionError"
      :cancel-label="t('common.no')"
      :confirm-label="t('employee_movement.execute')"
      icon="pi pi-check"
      severity="success"
      @confirm="handleExecuteConfirm"
      @cancel="executeConfirmVisible = false"
    />
    <ConfirmActionDialog
      v-model:visible="cancelConfirmVisible"
      :title="t('employee_movement.confirm_cancel_title')"
      :message="t('employee_movement.confirm_cancel_msg')"
      :loading="actionLoading"
      :error-msg="actionError"
      :cancel-label="t('common.no')"
      :confirm-label="t('employee_movement.cancel')"
      icon="pi pi-times"
      severity="danger"
      @confirm="handleCancelConfirm"
      @cancel="cancelConfirmVisible = false"
    />
    <ConfirmDeleteDialog
      v-model:visible="deleteConfirmVisible"
      :title="t('employee_movement.confirm_delete_title')"
      :message="t('employee_movement.confirm_delete_msg')"
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
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import ViewLabel from '@/components/ViewLabel.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()

// ── Daftar ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const searchTerm = ref('')
const filterType = ref(null)
const filterStatus = ref(null)

// ── Referensi dropdown ──
const employees = ref([])
const organizations = ref([])
const employmentStatuses = ref([])

// ── Dialog create ──
const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref(emptyForm())

// ── Konfirmasi aksi ──
const actionTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const submitConfirmVisible = ref(false)
const executeConfirmVisible = ref(false)
const cancelConfirmVisible = ref(false)
const deleteConfirmVisible = ref(false)

// ── Detail ──
const detailVisible = ref(false)
const detailItem = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const typeOptions = computed(() => [
  'promotion', 'demotion', 'mutation', 'contract_extension', 'status_change', 'retirement', 'offboarding', 'other'
].map(v => ({ label: typeLabel(v), value: v })))

const statusOptions = computed(() => [
  'draft', 'pending_approval', 'approved', 'rejected', 'executed', 'cancelled', 'cancellation_pending'
].map(v => ({ label: statusLabel(v), value: v })))

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.employee_id })))
const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id })))
const positionOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id })))

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

function emptyForm() {
  return {
    employee_id: '',
    movement_type: 'promotion',
    to_organization_id: null,
    to_position_id: null,
    to_employment_status_id: null,
    decision_letter_number: '',
    decision_letter_date: '',
    effective_date: '',
    reason: ''
  }
}

// ── Label & severity ──
function typeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function typeSeverity(type) {
  switch (type) {
    case 'promotion': return 'success'
    case 'demotion': return 'danger'
    case 'mutation': return 'info'
    case 'contract_extension': return 'warning'
    case 'status_change': return 'info'
    case 'retirement': return 'secondary'
    case 'offboarding': return 'danger'
    default: return 'secondary'
  }
}

function statusLabel(status) {
  const key = `employee_movement.status_${status}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'rejected': return 'danger'
    case 'executed': return 'success'
    case 'cancelled': return 'secondary'
    case 'cancellation_pending': return 'warning'
    default: return 'secondary'
  }
}

// ── Field per tipe (mirip validasi G-7 backend) ──
function requiresOrganization(type) {
  return type === 'mutation'
}

function requiresPosition(type) {
  return type === 'promotion' || type === 'demotion' || type === 'mutation'
}

// ── Load data ──
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (filterType.value) params.movement_type = filterType.value
    if (filterStatus.value) params.status = filterStatus.value
    if (searchTerm.value) params.search = searchTerm.value
    const res = await api.get('/api/v1/tenant/employee-movements/movements', { params })
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

function resetFilters() {
  searchTerm.value = ''
  filterType.value = null
  filterStatus.value = null
  currentPage.value = 1
  loadData()
}

// Referensi: employees, organizations aktif (summary status active — instruksi
// user: position = organization dari summary aktif), employment statuses.
async function loadReferences() {
  const [empRes, orgRes, statusRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } }),
    api.get('/api/v1/tenant/settings/employment-statuses', { params: { per_page: 500 } })
  ])
  employees.value = empRes.value?.data?.data || []
  organizations.value = orgRes.value?.data?.data || []
  employmentStatuses.value = statusRes.value?.data?.data || []
}

// ── Create ──
function openDialog() {
  errors.value = {}
  form.value = emptyForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = emptyForm()
  errors.value = {}
}

function validateForm() {
  const e = {}
  if (!form.value.employee_id) e.employee_id = t('employee_movement.field_required')
  if (!form.value.movement_type) e.movement_type = t('employee_movement.field_required')
  if (requiresOrganization(form.value.movement_type) && !form.value.to_organization_id && !form.value.to_position_id) {
    e.to_organization_id = t('employee_movement.field_required')
  }
  if (requiresPosition(form.value.movement_type) && !form.value.to_position_id) {
    e.to_position_id = t('employee_movement.field_required')
  }
  if (form.value.movement_type === 'status_change' && !form.value.to_employment_status_id) {
    e.to_employment_status_id = t('employee_movement.field_required')
  }
  if (!form.value.decision_letter_number?.trim()) e.decision_letter_number = t('employee_movement.field_required')
  if (!form.value.decision_letter_date) e.decision_letter_date = t('employee_movement.field_required')
  if (!form.value.effective_date) e.effective_date = t('employee_movement.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    await api.post('/api/v1/tenant/employee-movements/movements', {
      employee_id: form.value.employee_id,
      movement_type: form.value.movement_type,
      to_organization_id: form.value.to_organization_id || undefined,
      to_position_id: form.value.to_position_id || undefined,
      to_employment_status_id: form.value.to_employment_status_id || undefined,
      decision_letter_number: form.value.decision_letter_number,
      decision_letter_date: form.value.decision_letter_date,
      effective_date: form.value.effective_date,
      reason: form.value.reason || undefined
    })
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

// ── Aksi: submit / execute / cancel / delete ──
function runAction(fn) {
  actionError.value = ''
  actionLoading.value = true
  fn().finally(() => { actionLoading.value = false })
}

function openSubmitConfirm(data) {
  actionTarget.value = data
  actionError.value = ''
  submitConfirmVisible.value = true
}

function handleSubmitConfirm() {
  if (!actionTarget.value?.id) return
  runAction(async () => {
    try {
      await api.post(`/api/v1/tenant/employee-movements/movements/${actionTarget.value.id}/submit`, {})
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
      submitConfirmVisible.value = false
      actionTarget.value = null
      await loadData()
    } catch (e) {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    }
  })
}

function openExecuteConfirm(data) {
  actionTarget.value = data
  actionError.value = ''
  executeConfirmVisible.value = true
}

function handleExecuteConfirm() {
  if (!actionTarget.value?.id) return
  runAction(async () => {
    try {
      await api.post(`/api/v1/tenant/employee-movements/movements/${actionTarget.value.id}/execute`, {})
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
      executeConfirmVisible.value = false
      actionTarget.value = null
      await loadData()
    } catch (e) {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    }
  })
}

function openCancelConfirm(data) {
  actionTarget.value = data
  actionError.value = ''
  cancelConfirmVisible.value = true
}

function handleCancelConfirm() {
  if (!actionTarget.value?.id) return
  runAction(async () => {
    try {
      await api.post(`/api/v1/tenant/employee-movements/movements/${actionTarget.value.id}/cancel`, {})
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
      cancelConfirmVisible.value = false
      actionTarget.value = null
      await loadData()
    } catch (e) {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    }
  })
}

function openDeleteConfirm(data) {
  actionTarget.value = data
  actionError.value = ''
  deleteConfirmVisible.value = true
}

// ── Detail ──
function openDetail(data) {
  detailItem.value = data
  detailVisible.value = true
}

// Aksi dari dalam dialog detail: tutup detail lalu buka konfirmasi yang sama
// dengan yang dipakai di tabel (reuse actionTarget). Item disalin dulu karena
// `@hide` dialog men-null-kan detailItem.
function actionFromDetail(action) {
  const item = detailItem.value
  if (!item) return
  detailVisible.value = false
  if (action === 'submit') openSubmitConfirm(item)
  else if (action === 'execute') openExecuteConfirm(item)
  else if (action === 'cancel') openCancelConfirm(item)
  else if (action === 'delete') openDeleteConfirm(item)
}

function formatDateTime(v) {
  if (!v) return '-'
  const datePart = formatDate(v, locale.value)
  if (!datePart) return '-'
  const time = new Date(v).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

function hasAnyField(item, prefix) {
  return ['organization', 'position', 'employment_status'].some(suffix => !!item?.[`${prefix}_${suffix}_name`])
}

function handleDeleteConfirm() {
  if (!actionTarget.value?.id) return
  runAction(async () => {
    try {
      await api.delete(`/api/v1/tenant/employee-movements/movements/${actionTarget.value.id}`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      actionTarget.value = null
      await loadData()
    } catch (e) {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    }
  })
}

onMounted(() => {
  loadData()
  loadReferences()
})
</script>
