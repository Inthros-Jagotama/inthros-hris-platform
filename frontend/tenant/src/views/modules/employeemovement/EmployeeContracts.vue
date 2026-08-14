<template>
  <div class="space-y-1">
    <!-- ── Toolbar: filter + tombol ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
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
          :placeholder="t('employee_movement.search_contracts_placeholder')"
          class="!w-64"
          size="small"
          @keyup.enter="onFilterChange"
        />
        <Button
          v-if="searchTerm || filterStatus"
          icon="pi pi-times"
          severity="secondary"
          outlined
          size="small"
          v-tooltip.left="t('common.reset')"
          @click="resetFilters"
        />
        <Button :label="t('employee_movement.add_contract')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-file-edit text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee_movement.empty_contracts') }}</p>
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
      <Column field="contract_number" :header="t('employee_movement.contract_number')" style="width:150px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ data.contract_number || '-' }}</span></template>
      </Column>
      <Column field="contract_type" :header="t('employee_movement.contract_type')" style="width:130px">
        <template #body="{data}">
          <Tag :value="typeLabel(data.contract_type)" :severity="typeSeverity(data.contract_type)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="start_date" :header="t('employee_movement.start_date')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.start_date, locale) }}</span></template>
      </Column>
      <Column field="end_date" :header="t('employee_movement.end_date')" style="width:110px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">
            {{ data.end_date ? formatDate(data.end_date, locale) : '-' }}
          </span>
        </template>
      </Column>
      <Column field="extension_count" :header="t('employee_movement.extension_count')" style="width:100px">
        <template #body="{data}">
          <span v-if="data.extension_count > 0" class="inline-flex items-center gap-1">
            <Tag value="x" severity="warning" class="!text-xs !px-1 !py-0.5" />
            <span class="text-gray-600 dark:text-gray-300 text-sm font-medium">{{ data.extension_count }}</span>
          </span>
          <span v-else class="text-gray-400 dark:text-gray-500">0</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:120px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="document_url" :header="t('employee_movement.document')" style="width:100px">
        <template #body="{data}">
          <a v-if="data.document_url" :href="data.document_url" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline">
            <i class="pi pi-paperclip mr-1"></i>{{ t('employee_movement.attachment') }}
          </a>
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:110px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Tambah / Ubah Kontrak ── -->
    <Dialog v-model:visible="dialogVisible" :header="editingId ? t('employee_movement.edit_contract') : t('employee_movement.add_contract')" modal :style="{ width: '600px' }" @hide="resetForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('employee_movement.contract_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('employee_movement.employee')" required :errors="errors?.employee_id">
          <Select
            v-model="form.employee_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            class="w-full"
            :disabled="!!editingId"
            :placeholder="t('employee_movement.select_employee')"
            @change="onEmployeeChange"
          />
        </FormRow>
        <FormRow :label="t('employee_movement.contract_number')" :required="!!editingId" :errors="errors?.contract_number">
          <TextInput v-model="form.contract_number" :placeholder="t('employee_movement.number_auto_placeholder')" />
        </FormRow>
        <FormRow :label="t('employee_movement.contract_type')" required :errors="errors?.contract_type">
          <Select v-model="form.contract_type" :options="typeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
        <FormRow :label="t('employee_movement.start_date')" required :errors="errors?.start_date">
          <DateInput v-model="form.start_date" :disabled="!!editingId" />
        </FormRow>
        <FormRow :label="t('employee_movement.end_date')" :errors="errors?.end_date">
          <DateInput v-model="form.end_date" :min-date="form.start_date || null" />
        </FormRow>
        <FormRow :label="t('employee_movement.previous_contract')" :errors="errors?.previous_contract_id">
          <Select
            v-model="form.previous_contract_id"
            :options="previousContractOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            class="w-full"
            :placeholder="t('employee_movement.select_previous_contract')"
          />
        </FormRow>
        <FormRow :label="t('employee_movement.decision_letter_number')" :errors="errors?.decision_letter_number">
          <TextInput v-model="form.decision_letter_number" />
        </FormRow>
        <FormRow :label="t('employee_movement.notes')" :errors="errors?.notes">
          <TextInput v-model="form.notes" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('employee_movement.document')" :errors="errors?.document_url">
          <div class="flex items-center gap-2">
            <input ref="fileInputRef" type="file" class="hidden" @change="onFileSelected" />
            <Button
              icon="pi pi-upload"
              size="small"
              severity="secondary"
              outlined
              :label="selectedFile ? selectedFile.name : t('employee_movement.choose_file')"
              @click="fileInputRef?.click()"
              class="!justify-start max-w-full overflow-hidden"
              :loading="uploading"
              :disabled="uploading"
            />
            <Button v-if="selectedFile" icon="pi pi-times" size="small" text severity="danger" @click="clearSelectedFile" />
            <a v-if="form.document_url" :href="form.document_url" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline">
              <i class="pi pi-paperclip mr-1"></i>{{ t('employee_movement.attachment') }}
            </a>
          </div>
        </FormRow>
        <FormRow v-if="editingId" :label="t('common.status')" :errors="errors?.status">
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
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
      :title="t('employee_movement.confirm_delete_contract_title')"
      :message="t('employee_movement.confirm_delete_contract_msg')"
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
const filterStatus = ref(null)

// ── Referensi ──
const employees = ref([])
const contracts = ref([]) // semua kontrak (untuk dropdown previous contract)

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())
const fileInputRef = ref(null)
const selectedFile = ref(null)
const uploading = ref(false)

// ── Konfirmasi hapus ──
const deleteTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const deleteConfirmVisible = ref(false)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const typeOptions = computed(() => [
  'pkwt', 'pkwtt', 'daily', 'other'
].map(v => ({ label: typeLabel(v), value: v })))

const statusOptions = computed(() => [
  'active', 'expired', 'extended', 'terminated'
].map(v => ({ label: statusLabel(v), value: v })))

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.id })))

// Kontrak milik employee terpilih (untuk previous_contract), tanpa dirinya sendiri saat edit.
const previousContractOptions = computed(() => {
  const empId = form.value.employee_id
  if (!empId) return []
  return contracts.value
    .filter(c => c.employee_id === empId && c.id !== editingId.value)
    .map(c => ({ label: `${c.contract_number} (${statusLabel(c.status)})`, value: c.id }))
})

const skeletonColumns = [
  { type: 'compound', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-20' }
]

function emptyForm() {
  return {
    employee_id: '',
    contract_number: '',
    contract_type: 'pkwt',
    start_date: '',
    end_date: '',
    previous_contract_id: null,
    decision_letter_number: '',
    notes: '',
    document_url: '',
    status: 'active'
  }
}

// ── Label & severity ──
function typeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function typeSeverity(type) {
  switch (type) {
    case 'pkwt': return 'primary'
    case 'pkwtt': return 'info'
    case 'daily': return 'warning'
    default: return 'secondary'
  }
}

function statusLabel(status) {
  const key = `employee_movement.status_${status}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'expired': return 'secondary'
    case 'extended': return 'warning'
    case 'terminated': return 'danger'
    default: return 'secondary'
  }
}

// ── Load data ──
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (filterStatus.value) params.status = filterStatus.value
    if (searchTerm.value) params.search = searchTerm.value
    const res = await api.get('/api/v1/tenant/employee-movements/contracts', { params })
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
  filterStatus.value = null
  currentPage.value = 1
  loadData()
}

// Referensi: employees + semua kontrak (untuk dropdown previous contract).
async function loadReferences() {
  const [empRes, contractRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/employee-movements/contracts', { params: { per_page: 500 } })
  ])
  employees.value = empRes.value?.data?.data || []
  contracts.value = contractRes.value?.data?.data || []
}

function onEmployeeChange() {
  form.value.previous_contract_id = null
}

// ── Dialog ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      employee_id: data.employee_id,
      contract_number: data.contract_number,
      contract_type: data.contract_type,
      start_date: data.start_date || '',
      end_date: data.end_date || '',
      previous_contract_id: data.previous_contract_id || null,
      decision_letter_number: data.decision_letter_number || '',
      notes: data.notes || '',
      document_url: data.document_url || '',
      status: data.status || 'active'
    }
  } else {
    form.value = emptyForm()
  }
  selectedFile.value = null
  dialogVisible.value = true
}

function resetForm() {
  form.value = emptyForm()
  errors.value = {}
  editingId.value = null
  selectedFile.value = null
}

// ── Upload dokumen ──
function onFileSelected(e) {
  const file = e.target.files?.[0]
  if (file) {
    selectedFile.value = file
    errors.value.document_url = null
  }
}

function clearSelectedFile() {
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

async function uploadAttachment() {
  if (!selectedFile.value) return
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', selectedFile.value)
    const res = await api.post('/api/v1/tenant/uploads', fd, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    form.value.document_url = res.data?.data?.url || ''
    selectedFile.value = null
    if (fileInputRef.value) fileInputRef.value.value = ''
  } catch (e) {
    // Biarkan handleSave yang menampilkan error (hindari double toast).
    throw e
  } finally {
    uploading.value = false
  }
}

// ── Simpan ──
function validateForm() {
  const e = {}
  if (!form.value.employee_id) e.employee_id = t('employee_movement.field_required')
  if (editingId.value && !form.value.contract_number?.trim()) e.contract_number = t('employee_movement.field_required')
  if (!form.value.contract_type) e.contract_type = t('employee_movement.field_required')
  if (!form.value.start_date) e.start_date = t('employee_movement.field_required')
  if (form.value.contract_type === 'pkwt' && !form.value.end_date) e.end_date = t('employee_movement.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    if (selectedFile.value) {
      await uploadAttachment()
    }
    if (editingId.value) {
      await api.put(`/api/v1/tenant/employee-movements/contracts/${editingId.value}`, {
        contract_number: form.value.contract_number,
        contract_type: form.value.contract_type,
        end_date: form.value.end_date || undefined,
        decision_letter_number: form.value.decision_letter_number || undefined,
        notes: form.value.notes || undefined,
        document_url: form.value.document_url || undefined,
        status: form.value.status || undefined
      })
    } else {
      await api.post('/api/v1/tenant/employee-movements/contracts', {
        employee_id: form.value.employee_id,
        contract_number: form.value.contract_number,
        contract_type: form.value.contract_type,
        start_date: form.value.start_date,
        end_date: form.value.end_date || undefined,
        previous_contract_id: form.value.previous_contract_id || undefined,
        decision_letter_number: form.value.decision_letter_number || undefined,
        notes: form.value.notes || undefined,
        document_url: form.value.document_url || undefined
      })
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadData()
    await loadReferences()
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
  api.delete(`/api/v1/tenant/employee-movements/contracts/${deleteTarget.value.id}`)
    .then(async () => {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      deleteTarget.value = null
      await loadData()
      await loadReferences()
    })
    .catch(e => {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    })
    .finally(() => { actionLoading.value = false })
}

onMounted(() => {
  loadData()
  loadReferences()
})
</script>
