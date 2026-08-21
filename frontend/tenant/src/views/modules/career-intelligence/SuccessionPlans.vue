<template>
  <div class="space-y-1">
    <!-- ── Toolbar ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('succession_plans.add_plan')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-user-plus text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('succession_plans.empty_plans') }}</p>
        </div>
      </template>
      <Column :header="t('succession_plans.position')" style="width:200px">
        <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ positionName(data.position_id) }}</span></template>
      </Column>
      <Column :header="t('succession_plans.successor')" style="width:200px">
        <template #body="{data}"><span class="text-gray-700 dark:text-gray-200">{{ employeeName(data.successor_id) }}</span></template>
      </Column>
      <Column field="readiness_level" :header="t('succession_plans.readiness_level')" style="width:130px">
        <template #body="{data}">
          <Tag :value="t('succession_plans.readiness_' + data.readiness_level)" :severity="readinessSeverity(data.readiness_level)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="priority_order" :header="t('succession_plans.priority_order')" style="width:90px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.priority_order }}</span></template>
      </Column>
      <Column field="target_date" :header="t('succession_plans.target_date')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.target_date ? formatDate(data.target_date, locale) : '—' }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}">
          <Tag :value="data.status" severity="success" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:110px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button icon="pi pi-eye" size="small" severity="info" text v-tooltip.left="t('common.view')" @click="openDetail(data)" />
            <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Detail ── -->
    <Dialog v-model:visible="detailVisible" :header="t('succession_plans.plan_detail_title')" modal :style="{ width: '480px' }">
      <div v-if="detailItem" class="space-y-3">
        <div>
          <p class="text-base font-semibold text-navy-800 dark:text-gray-100">{{ positionName(detailItem.position_id) }}</p>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ t('succession_plans.successor') }}: {{ employeeName(detailItem.successor_id) }}</p>
        </div>
        <div class="flex items-center gap-2 flex-wrap">
          <Tag :value="t('succession_plans.readiness_' + detailItem.readiness_level)" :severity="readinessSeverity(detailItem.readiness_level)" class="!text-xs !px-1.5 !py-0.5" />
          <Tag :value="detailItem.status" severity="success" class="!text-xs !px-1.5 !py-0.5" />
          <span class="text-xs text-gray-400">{{ t('succession_plans.priority_order') }}: {{ detailItem.priority_order }}</span>
        </div>
        <div v-if="detailItem.target_date" class="text-sm">
          <span class="text-gray-500 dark:text-gray-400">{{ t('succession_plans.target_date') }}: </span>
          <span class="text-gray-700 dark:text-gray-200">{{ formatDate(detailItem.target_date, locale) }}</span>
        </div>
        <div v-if="detailItem.development_plan">
          <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-1">{{ t('succession_plans.development_plan') }}</p>
          <p class="text-sm text-gray-700 dark:text-gray-200 whitespace-pre-line">{{ detailItem.development_plan }}</p>
        </div>
        <div v-if="detailItem.notes">
          <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-1">{{ t('succession_plans.notes') }}</p>
          <p class="text-sm text-gray-700 dark:text-gray-200 whitespace-pre-line">{{ detailItem.notes }}</p>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="detailVisible = false" />
          <Button :label="t('common.edit')" icon="pi pi-pencil" size="small" @click="openDialog(detailItem); detailVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Tambah / Ubah ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingId ? t('succession_plans.edit_plan') : t('succession_plans.add_plan')"
      modal
      :style="{ width: '560px' }"
      @hide="resetForm"
    >
      <div class="space-y-3">
        <FormRow :label="t('succession_plans.position')" required :errors="errors?.position_id">
          <Select
            v-model="form.position_id"
            :options="positionOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            :disabled="!!editingId"
            class="w-full !text-sm"
            :placeholder="t('career_paths.select_position')"
          />
        </FormRow>
        <FormRow :label="t('succession_plans.successor')" required :errors="errors?.successor_id">
          <Select
            v-model="form.successor_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            :disabled="!!editingId"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('succession_plans.readiness_level')" required :errors="errors?.readiness_level">
            <Select
              v-model="form.readiness_level"
              :options="readinessOptions"
              optionLabel="label"
              optionValue="value"
              class="w-full !text-sm"
              :placeholder="t('common.select')"
            />
          </FormRow>
          <FormRow :label="t('succession_plans.priority_order')" :errors="errors?.priority_order">
            <InputNumber v-model="form.priority_order" :min="1" :useGrouping="false" class="w-full" />
          </FormRow>
        </div>
        <FormRow :label="t('succession_plans.target_date')" :errors="errors?.target_date">
          <DatePicker v-model="form.target_date" dateFormat="yy-mm-dd" :placeholder="t('common.select_date')" class="w-full" showIcon />
        </FormRow>
        <FormRow :label="t('succession_plans.development_plan')" :errors="errors?.development_plan">
          <TextInput v-model="form.development_plan" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('succession_plans.notes')" :errors="errors?.notes">
          <TextInput v-model="form.notes" textarea :rows="2" />
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
      :title="t('succession_plans.confirm_delete_title')"
      :message="t('succession_plans.confirm_delete_msg')"
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
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()

// ── Daftar ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

// ── Referensi ──
const organizations = ref([])
const employees = ref([])

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())

// ── Detail ──
const detailVisible = ref(false)
const detailItem = ref(null)

// ── Konfirmasi hapus ──
const deleteTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const deleteConfirmVisible = ref(false)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// Referensi posisi memakai organizations (konsep organization = position).
const positionOptions = computed(() =>
  organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id }))
)
const employeeOptions = computed(() =>
  employees.value.map(e => ({ label: e.name, value: e.id }))
)
const readinessOptions = computed(() => [
  { label: t('succession_plans.readiness_READY_NOW'), value: 'READY_NOW' },
  { label: t('succession_plans.readiness_READY_1YR'), value: 'READY_1YR' },
  { label: t('succession_plans.readiness_READY_2YR'), value: 'READY_2YR' },
  { label: t('succession_plans.readiness_POTENTIAL'), value: 'POTENTIAL' }
])

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

function emptyForm() {
  return {
    position_id: null,
    successor_id: null,
    readiness_level: null,
    priority_order: 1,
    target_date: null,
    development_plan: '',
    notes: ''
  }
}

function positionName(id) {
  return organizations.value.find(o => o.id === id)?.nomenclature || '—'
}
function employeeName(id) {
  const e = employees.value.find(e => e.id === id)
  return e ? (e.name || e.full_name) : '—'
}
function readinessSeverity(level) {
  if (level === 'READY_NOW') return 'success'
  if (level === 'READY_1YR') return 'info'
  if (level === 'READY_2YR') return 'warning'
  return 'secondary'
}

// ── Load data ──
async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/successions', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
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

async function loadReferences() {
  const [orgRes, empRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } }),
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  ])
  organizations.value = orgRes.value?.data?.data || []
  employees.value = empRes.value?.data?.data || []
}

// ── Dialog create/edit ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      position_id: data.position_id,
      successor_id: data.successor_id,
      readiness_level: data.readiness_level,
      priority_order: data.priority_order ?? 1,
      target_date: data.target_date ? new Date(data.target_date) : null,
      development_plan: data.development_plan || '',
      notes: data.notes || ''
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

// ── Detail ──
function openDetail(data) {
  detailItem.value = data
  detailVisible.value = true
}

// ── Simpan ──
function validateForm() {
  const e = {}
  if (!form.value.position_id) e.position_id = t('career_paths.field_required')
  if (!form.value.successor_id) e.successor_id = t('career_paths.field_required')
  if (!form.value.readiness_level) e.readiness_level = t('career_paths.field_required')
  return e
}

function formatDateForApi(date) {
  if (!date) return null
  if (typeof date === 'string') return date
  return new Date(date).toISOString().split('T')[0]
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    if (editingId.value) {
      // Update: position_id/successor_id tidak diikutsertakan (immutable
      // setelah dibuat sesuai kontrak backend UpdateSuccessionPlanRequest).
      const payload = {
        readiness_level: form.value.readiness_level,
        priority_order: form.value.priority_order ?? 1,
        target_date: formatDateForApi(form.value.target_date),
        development_plan: form.value.development_plan?.trim() || '',
        notes: form.value.notes?.trim() || ''
      }
      await api.put(`/api/v1/tenant/career-intelligence/successions/${editingId.value}`, payload)
    } else {
      const payload = {
        position_id: form.value.position_id,
        successor_id: form.value.successor_id,
        readiness_level: form.value.readiness_level,
        priority_order: form.value.priority_order ?? 1,
        target_date: formatDateForApi(form.value.target_date) || undefined,
        development_plan: form.value.development_plan?.trim() || undefined,
        notes: form.value.notes?.trim() || undefined
      }
      await api.post('/api/v1/tenant/career-intelligence/successions', payload)
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
  api.delete(`/api/v1/tenant/career-intelligence/successions/${deleteTarget.value.id}`)
    .then(async () => {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      deleteTarget.value = null
      await loadData()
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
