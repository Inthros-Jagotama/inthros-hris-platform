<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('attendance.view_in_approvals')" icon="pi pi-external-link" size="small" severity="secondary" outlined @click="router.push('/approvals')" />
        <Button :label="t('attendance.overtime_assign')" icon="pi pi-user-plus" size="small" severity="secondary" outlined :disabled="!employeeId" @click="openAssignDialog()" />
        <Button :label="t('common.add')" icon="pi pi-plus" size="small" :disabled="!employeeId" @click="openDialog()" />
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
          <p class="text-sm font-medium">{{ t('attendance.overtime_empty') }}</p>
        </div>
      </template>
      <Column field="employee_name" :header="t('attendance.submitted_by')" style="width:160px">
        <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ data.employee_name || '-' }}</span></template>
      </Column>
      <Column field="organization_name" :header="t('attendance.organization')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.organization_name || '-' }}</span></template>
      </Column>
      <Column field="flow_type" :header="t('attendance.flow_type')" style="width:100px">
        <template #body="{data}">
          <Tag :value="data.flow_type === 'ASSIGNED' ? t('attendance.flow_assigned') : t('attendance.flow_self')" :severity="data.flow_type === 'ASSIGNED' ? 'info' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="assigned_by" :header="t('attendance.assigned_by')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ assignerName(data.assigned_by) }}</span></template>
      </Column>
      <Column field="work_date" :header="t('attendance.work_date')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.work_date, locale) }}</span></template>
      </Column>
      <Column :header="t('attendance.start_time')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatTime(data.start_time_local) }}</span></template>
      </Column>
      <Column :header="t('attendance.end_time')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatTime(data.end_time_local) }}</span></template>
      </Column>
      <Column field="requested_minutes" :header="t('attendance.requested_minutes')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.requested_minutes }} min</span></template>
      </Column>
      <Column field="calculated_minutes" :header="t('attendance.calculated_minutes')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.calculated_minutes ?? '-' }}</span></template>
      </Column>
      <Column :header="t('attendance.actual_start_time')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatTime(data.actual_start_time_local) }}</span></template>
      </Column>
      <Column :header="t('attendance.actual_end_time')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatTime(data.actual_end_time_local) }}</span></template>
      </Column>
      <Column field="actual_note" :header="t('attendance.actual_note')">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.actual_note || '-' }}</span>
          <a v-if="data.attachment_url" :href="data.attachment_url" target="_blank" class="block text-xs text-emerald-600 dark:text-emerald-400 hover:underline mt-0.5">
            <i class="pi pi-paperclip mr-1"></i>{{ t('attendance.attachment') }}
          </a>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:140px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:200px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button
              v-if="data.status === 'WAITING_ACTUAL'"
              :label="t('attendance.fill_actual')"
              icon="pi pi-pencil"
              size="small"
              severity="primary"
              @click="openActualDialog(data)"
            />
            <Button
              v-if="data.status === 'PENDING_APPROVAL' || data.status === 'WAITING_ACTUAL'"
              icon="pi pi-times"
              size="small"
              severity="danger"
              text
              v-tooltip.left="t('attendance.cancel_request')"
              @click="openCancelConfirm(data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Ajukan Lembur (SELF) ── -->
    <Dialog v-model:visible="dialogVisible" :header="t('attendance.overtime_new')" modal :style="{ width: '480px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('attendance.work_date')" required :errors="errors?.work_date">
          <DateInput v-model="form.work_date" />
        </FormRow>
        <FormRow :label="t('attendance.start_time')" required :errors="errors?.start_time_local">
          <TimeInput v-model="form.start_time" />
        </FormRow>
        <FormRow :label="t('attendance.end_time')" required :errors="errors?.end_time_local">
          <TimeInput v-model="form.end_time" />
        </FormRow>
        <FormRow :label="t('attendance.requested_minutes')" required :errors="errors?.requested_minutes">
          <InputNumber v-model="form.requested_minutes" class="!w-full" :min="1" size="small" suffix=" min" />
          <p v-if="isCrossDay" class="text-xs text-amber-500 mt-1 flex items-center gap-1">
            <i class="pi pi-moon"></i>
            {{ t('attendance.overtime_cross_day') }}
          </p>
        </FormRow>
        <FormRow :label="t('attendance.reason')" :errors="errors?.reason">
          <TextInput v-model="form.reason" textarea :rows="3" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Tugaskan Lembur (ASSIGNED) ── -->
    <Dialog v-model:visible="assignDialogVisible" :header="t('attendance.overtime_assign')" modal :style="{ width: '520px' }" @hide="resetAssignForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('attendance.overtime_assign_description') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('attendance.assigned_employee')" required :errors="assignErrors?.assigned_employee_id">
          <Select v-model="assignForm.assigned_employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :placeholder="t('attendance.select_employee')" />
        </FormRow>
        <FormRow :label="t('attendance.work_date')" required :errors="assignErrors?.work_date">
          <DateInput v-model="assignForm.work_date" />
        </FormRow>
        <FormRow :label="t('attendance.start_time')" required :errors="assignErrors?.start_time_local">
          <TimeInput v-model="assignForm.start_time" />
        </FormRow>
        <FormRow :label="t('attendance.end_time')" required :errors="assignErrors?.end_time_local">
          <TimeInput v-model="assignForm.end_time" />
        </FormRow>
        <FormRow :label="t('attendance.requested_minutes')" required :errors="assignErrors?.requested_minutes">
          <InputNumber v-model="assignForm.requested_minutes" class="!w-full" :min="1" size="small" suffix=" min" />
          <p v-if="isAssignCrossDay" class="text-xs text-amber-500 mt-1 flex items-center gap-1">
            <i class="pi pi-moon"></i>
            {{ t('attendance.overtime_cross_day') }}
          </p>
        </FormRow>
        <FormRow :label="t('attendance.reason')" :errors="assignErrors?.reason">
          <TextInput v-model="assignForm.reason" textarea :rows="3" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="assignDialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="assignSaving" :disabled="assignSaving" @click="handleAssign" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Isi Aktual ── -->
    <Dialog v-model:visible="actualDialogVisible" :header="t('attendance.fill_actual_title')" modal :style="{ width: '520px' }" @hide="resetActualForm">
      <div class="space-y-3">
        <FormRow :label="t('attendance.actual_start_time')" required :errors="actualErrors?.actual_start_time_local">
          <TimeInput v-model="actualForm.actual_start_time" />
        </FormRow>
        <FormRow :label="t('attendance.actual_end_time')" required :errors="actualErrors?.actual_end_time_local">
          <TimeInput v-model="actualForm.actual_end_time" />
        </FormRow>
        <FormRow :label="t('attendance.actual_note')" :errors="actualErrors?.actual_note">
          <TextInput v-model="actualForm.actual_note" textarea :rows="3" />
        </FormRow>
        <FormRow :label="t('attendance.attachment')" :errors="actualErrors?.attachment_url">
          <div class="flex items-center gap-2">
            <input ref="fileInputRef" type="file" class="hidden" @change="onFileSelected" />
            <Button
              icon="pi pi-upload"
              size="small"
              severity="secondary"
              outlined
              :label="selectedFile ? selectedFile.name : t('attendance.choose_file')"
              @click="fileInputRef?.click()"
              class="!justify-start max-w-full overflow-hidden"
              :loading="uploading"
              :disabled="uploading"
            />
            <Button v-if="selectedFile" icon="pi pi-times" size="small" text severity="danger" @click="clearSelectedFile" />
            <a v-if="actualForm.attachment_url" :href="actualForm.attachment_url" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline">
              <i class="pi pi-paperclip mr-1"></i>{{ t('attendance.attachment') }}
            </a>
          </div>
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="actualDialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="actualSaving" :disabled="actualSaving" @click="handleSubmitActual" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Batal ── -->
    <ConfirmDeleteDialog
      v-model:visible="cancelVisible"
      :title="t('attendance.cancel_request')"
      :message="t('attendance.confirm_cancel_overtime')"
      :loading="cancelLoading"
      :error-msg="cancelError"
      :cancel-label="t('common.no')"
      :confirm-label="t('attendance.cancel_request')"
      @confirm="handleCancelConfirm"
      @cancel="cancelVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { localDateTimeISOString } from '@/utils/localTime'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import TimeInput from '@/components/TimeInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

// ── Dialog SELF ──
const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref({ work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' })

// ── Dialog ASSIGNED ──
const assignDialogVisible = ref(false)
const assignSaving = ref(false)
const assignErrors = ref({})
const employees = ref([])
const assignForm = ref({ assigned_employee_id: '', work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' })

// ── Dialog Isi Aktual ──
const actualDialogVisible = ref(false)
const actualSaving = ref(false)
const actualErrors = ref({})
const actualTarget = ref(null)
const actualForm = ref({ actual_start_time: '', actual_end_time: '', actual_note: '', attachment_url: '' })
const fileInputRef = ref(null)
const selectedFile = ref(null)
const uploading = ref(false)

// ── Dialog Batal ──
const cancelVisible = ref(false)
const cancelTarget = ref(null)
const cancelLoading = ref(false)
const cancelError = ref('')

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
// Label option: nama + kode karyawan (employee_code, bukan UUID internal).
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.employee_id })))

// assigned_by dari backend berupa UUID employee — resolve ke nama via daftar
// karyawan yang sudah dimuat (untuk dropdown assign).
function assignerName(id) {
  if (!id) return '-'
  const emp = employees.value.find(e => e.employee_id === id)
  return emp?.name || id.slice(0, 8)
}

// ── Perhitungan menit otomatis + dukungan lintas hari (SELF) ──
function timeToMinutes(t) {
  const m = String(t || '').match(/^(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?$/)
  if (!m) return null
  return Number(m[1]) * 60 + Number(m[2])
}

function addDays(dateStr, days) {
  const [y, mo, d] = String(dateStr).split('-').map(Number)
  const date = new Date(y, mo - 1, d + days)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function isCrossDayOf(f) {
  const s = timeToMinutes(f.start_time)
  const e = timeToMinutes(f.end_time)
  return s !== null && e !== null && e <= s
}

const isCrossDay = computed(() => isCrossDayOf(form.value))
const isAssignCrossDay = computed(() => isCrossDayOf(assignForm.value))

function autoMinutesOf(f) {
  const s = timeToMinutes(f.start_time)
  const e = timeToMinutes(f.end_time)
  if (s === null || e === null) return null
  let diff = e - s
  if (diff <= 0) diff += 1440
  return diff
}

const autoMinutes = computed(() => autoMinutesOf(form.value))
const assignAutoMinutes = computed(() => autoMinutesOf(assignForm.value))

// Isi otomatis saat start/end lengkap (tetap bisa diedit manual)
watch(autoMinutes, (m) => {
  if (m !== null && dialogVisible.value) form.value.requested_minutes = m
})
watch(assignAutoMinutes, (m) => {
  if (m !== null && assignDialogVisible.value) assignForm.value.requested_minutes = m
})

const skeletonColumns = [
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-20' }
]

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'PENDING_APPROVAL': return 'info'
    case 'WAITING_ACTUAL': return 'warning'
    case 'ACTUAL_SUBMITTED': return 'info'
    case 'CANCELLED': return 'secondary'
    default: return 'secondary'
  }
}

function statusLabel(status) {
  const key = `attendance.status_${status.toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

// formatTime — HH:MM lokal dari timestamp RFC3339 (pola sama dgn Approvals/Notifications).
function formatTime(v) {
  if (!v) return '-'
  const d = new Date(v)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

async function loadData() {
  if (!employeeId.value) return
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value, employee_id: employeeId.value }
    const res = await api.get('/api/v1/tenant/attendance/overtime-requests', { params })
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

// ═══════════ SELF ═══════════
function openDialog() {
  errors.value = {}
  form.value = { work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' }
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  if (!form.value.work_date) { errors.value = { work_date: t('form.required') }; return }
  if (!form.value.start_time?.trim()) { errors.value = { start_time_local: t('form.required') }; return }
  if (!form.value.end_time?.trim()) { errors.value = { end_time_local: t('form.required') }; return }
  if (!form.value.requested_minutes) { errors.value = { requested_minutes: t('form.required') }; return }
  saving.value = true
  try {
    await api.post('/api/v1/tenant/attendance/overtime-requests', {
      employee_id: employeeId.value,
      work_date: form.value.work_date,
      start_time_local: localDateTimeISOString(form.value.work_date, form.value.start_time),
      // Lintas hari: end_time <= start_time → tanggal end = hari berikutnya
      end_time_local: localDateTimeISOString(
        isCrossDay.value ? addDays(form.value.work_date, 1) : form.value.work_date,
        form.value.end_time
      ),
      requested_minutes: form.value.requested_minutes,
      reason: form.value.reason
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

// ═══════════ ASSIGNED ═══════════
function openAssignDialog() {
  assignErrors.value = {}
  assignForm.value = { assigned_employee_id: '', work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' }
  assignDialogVisible.value = true
}

function resetAssignForm() {
  assignForm.value = { assigned_employee_id: '', work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' }
  assignErrors.value = {}
}

async function handleAssign() {
  assignErrors.value = {}
  if (!assignForm.value.assigned_employee_id) { assignErrors.value = { assigned_employee_id: t('form.required') }; return }
  if (!assignForm.value.work_date) { assignErrors.value = { work_date: t('form.required') }; return }
  if (!assignForm.value.start_time?.trim()) { assignErrors.value = { start_time_local: t('form.required') }; return }
  if (!assignForm.value.end_time?.trim()) { assignErrors.value = { end_time_local: t('form.required') }; return }
  if (!assignForm.value.requested_minutes) { assignErrors.value = { requested_minutes: t('form.required') }; return }
  assignSaving.value = true
  try {
    await api.post('/api/v1/tenant/attendance/overtime-requests/assign', {
      assigned_employee_id: assignForm.value.assigned_employee_id,
      work_date: assignForm.value.work_date,
      start_time_local: localDateTimeISOString(assignForm.value.work_date, assignForm.value.start_time),
      end_time_local: localDateTimeISOString(
        isAssignCrossDay.value ? addDays(assignForm.value.work_date, 1) : assignForm.value.work_date,
        assignForm.value.end_time
      ),
      requested_minutes: assignForm.value.requested_minutes,
      reason: assignForm.value.reason
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    assignDialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      assignErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    assignSaving.value = false
  }
}

// ═══════════ Isi Aktual ═══════════
function openActualDialog(data) {
  actualErrors.value = {}
  actualTarget.value = data
  actualForm.value = {
    actual_start_time: data.actual_start_time_local ? formatTime(data.actual_start_time_local) : '',
    actual_end_time: data.actual_end_time_local ? formatTime(data.actual_end_time_local) : '',
    actual_note: data.actual_note || '',
    attachment_url: data.attachment_url || ''
  }
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
  actualDialogVisible.value = true
}

function resetActualForm() {
  actualTarget.value = null
  actualForm.value = { actual_start_time: '', actual_end_time: '', actual_note: '', attachment_url: '' }
  actualErrors.value = {}
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

function onFileSelected(e) {
  const file = e.target.files?.[0]
  if (file) {
    selectedFile.value = file
    actualErrors.value.attachment_url = null
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
    actualForm.value.attachment_url = res.data?.data?.url || ''
    selectedFile.value = null
    if (fileInputRef.value) fileInputRef.value.value = ''
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    throw e
  } finally {
    uploading.value = false
  }
}

async function handleSubmitActual() {
  actualErrors.value = {}
  if (!actualForm.value.actual_start_time?.trim()) { actualErrors.value = { actual_start_time_local: t('form.required') }; return }
  if (!actualForm.value.actual_end_time?.trim()) { actualErrors.value = { actual_end_time_local: t('form.required') }; return }
  if (!actualTarget.value?.id) return
  actualSaving.value = true
  try {
    if (selectedFile.value) await uploadAttachment()
    await api.post(`/api/v1/tenant/attendance/overtime-requests/${actualTarget.value.id}/actual`, {
      actual_start_time_local: localDateTimeISOString(actualTarget.value.work_date, actualForm.value.actual_start_time),
      actual_end_time_local: localDateTimeISOString(
        timeToMinutes(actualForm.value.actual_end_time) <= timeToMinutes(actualForm.value.actual_start_time)
          ? addDays(actualTarget.value.work_date, 1)
          : actualTarget.value.work_date,
        actualForm.value.actual_end_time
      ),
      actual_note: actualForm.value.actual_note,
      attachment_url: actualForm.value.attachment_url
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    actualDialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      actualErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    actualSaving.value = false
  }
}

// ═══════════ Batal ═══════════
function openCancelConfirm(data) {
  cancelTarget.value = data
  cancelError.value = ''
  cancelVisible.value = true
}

async function handleCancelConfirm() {
  if (!cancelTarget.value?.id) return
  cancelLoading.value = true
  cancelError.value = ''
  try {
    await api.post(`/api/v1/tenant/attendance/overtime-requests/${cancelTarget.value.id}/cancel`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    cancelVisible.value = false
    cancelTarget.value = null
    await loadData()
  } catch (e) {
    cancelError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    cancelLoading.value = false
  }
}

onMounted(async () => {
  await loadMyEmployeeId()
  await Promise.all([loadData(), loadEmployees()])
})

// Bawahan yang bisa ditugaskan: hanya karyawan di organisasi bawahan efektif
// dari user yang login (anak langsung terisi, atau turun level sampai ada
// karyawannya) — backend attendance/overtime-requests/assignable-employees.
async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/attendance/overtime-requests/assignable-employees')
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}
</script>
