<template>
  <div class="space-y-4">
    <div v-if="loading" class="flex items-center justify-center h-40">
      <i class="pi pi-spinner pi-spin text-2xl text-emerald-500"></i>
    </div>

    <template v-else-if="!employeeId">
      <Message severity="warn" :closable="false">{{ t('leave.no_employee_linked') }}</Message>
    </template>

    <template v-else>
      <!-- Top actions -->
      <div class="flex justify-end gap-2">
        <Button v-if="hasPermission('leave.update')" :label="t('leave.admin')" icon="pi pi-cog" size="small" severity="secondary" outlined @click="router.push('/leave/admin')" />
      </div>

      <!-- Balance cards -->
      <div v-if="balancesLoading" class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="n in 4" :key="n" class="h-24 bg-gray-100 dark:bg-gray-700 rounded-lg animate-pulse"></div>
      </div>
      <div v-else-if="balances.length === 0" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6 text-center text-gray-400 dark:text-gray-500">
        <i class="pi pi-wallet text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('leave.balances_empty') }}</p>
      </div>
      <div v-else class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="b in balances" :key="b.id" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-1">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ leaveTypeName(b.leave_type_id) }}</span>
            <Tag :value="b.period_year" severity="secondary" class="!text-[10px] !px-1.5 !py-0.5" />
          </div>
          <div class="text-xl font-bold text-emerald-600 dark:text-emerald-400">{{ formatDays(b.remaining_days) }}</div>
          <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">
            {{ t('leave.quota_days') }} {{ formatDays(b.quota_days) }} · {{ t('leave.used_days') }} {{ formatDays(b.used_days) }}
          </div>
        </div>
      </div>

      <!-- My Requests -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between flex-wrap gap-2 mb-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('leave.my_requests') }}</h2>
          <Button :label="t('leave.new_request')" icon="pi pi-plus" size="small" :disabled="!canCreate" @click="openDialog()" />
        </div>
        <SkeletonTable v-if="requestsLoading" :columns="skeletonColumns" :rows="6" />
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
              <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
              <p class="text-sm font-medium">{{ t('leave.requests_empty') }}</p>
            </div>
          </template>
          <Column :header="t('leave.leave_type')" style="width:180px">
            <template #body="{data}">
              <span class="text-gray-800 dark:text-gray-100 font-medium">{{ leaveTypeName(data.leave_type_id) }}</span>
            </template>
          </Column>
          <Column :header="t('leave.request_start_date')" style="width:150px">
            <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.request_start_date, locale) }}</span></template>
          </Column>
          <Column :header="t('leave.request_end_date')" style="width:150px">
            <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.request_end_date, locale) }}</span></template>
          </Column>
          <Column :header="t('leave.requested_days')" style="width:100px">
            <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDays(data.requested_days) }}</span></template>
          </Column>
          <Column field="status" :header="t('common.status')" style="width:160px">
            <template #body="{data}">
              <Tag :value="t('leave.status_' + data.status.toLowerCase())" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column field="supervisor_note" :header="t('leave.approval_note')" style="width:220px">
            <template #body="{data}">
              <span v-if="data.supervisor_note" class="text-xs text-gray-600 dark:text-gray-300">{{ data.supervisor_note }}</span>
              <span v-else class="text-xs text-gray-300 dark:text-gray-600">-</span>
            </template>
          </Column>
          <Column :header="t('common.actions')" style="width:80px" frozen alignFrozen="right">
            <template #body="{data}">
              <div class="flex items-center gap-1 justify-end">
                <Button
                  v-if="canCancel(data.status)"
                  icon="pi pi-times"
                  size="small"
                  text
                  severity="danger"
                  v-tooltip.left="t('leave.cancel_request')"
                  @click="confirmCancel(data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- This-month calendar -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('leave.calendar_this_month') }}</h2>
        <div v-if="calendarEntries.length === 0" class="text-center py-8 text-gray-400 dark:text-gray-500">
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('leave.calendar_empty') }}</p>
        </div>
        <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
          <div v-for="entry in calendarEntries" :key="entry.leave_request_id + entry.leave_date" class="flex items-center justify-between py-2 text-sm">
            <span class="text-gray-600 dark:text-gray-300">{{ formatDate(entry.leave_date, locale) }}</span>
            <div class="flex items-center gap-3">
              <span v-if="entry.day_fraction < 1" class="text-xs text-gray-400 dark:text-gray-500">{{ formatDays(entry.day_fraction) }} {{ t('leave.requested_days') }}</span>
              <Tag :value="leaveTypeName(entry.leave_type_id)" :severity="statusSeverity(entry.status)" class="!text-xs !px-1.5 !py-0.5" />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- New Request dialog -->
    <Dialog v-model:visible="dialogVisible" :header="t('leave.new_request')" modal :style="{ width: '720px', maxWidth: '95vw' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('leave.leave_type')" required :errors="errors?.leave_type_id">
          <SelectLabel v-model="form.leave_type_id" :options="leaveTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.leave_type_id }" />
        </FormRow>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('leave.request_start_date')" required :errors="errors?.request_start_date">
            <DateInput v-model="form.request_start_date" :maxDate="form.request_end_date || null" />
          </FormRow>
          <FormRow :label="t('leave.request_end_date')" required :errors="errors?.request_end_date">
            <DateInput v-model="form.request_end_date" :minDate="form.request_start_date || null" />
          </FormRow>
        </div>
        <FormRow :label="t('leave.duration_mode')" :errors="errors?.duration_mode">
          <SelectLabel v-model="form.duration_mode" :options="durationModeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <div v-if="form.duration_mode === 'HOURLY'" class="grid grid-cols-2 gap-3">
          <FormRow :label="t('leave.start_time')">
            <InputText v-model="form.start_time" placeholder="09:00" size="small" class="w-full" />
          </FormRow>
          <FormRow :label="t('leave.end_time')">
            <InputText v-model="form.end_time" placeholder="17:00" size="small" class="w-full" />
          </FormRow>
        </div>
        <FormRow :label="t('leave.leave_reason')" :errors="errors?.leave_reason_id">
          <SelectLabel v-model="form.leave_reason_id" :options="reasonOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('leave.leave_reason_note')" :errors="errors?.leave_reason_note">
          <TextInput v-model="form.leave_reason_note" textarea :rows="3" />
        </FormRow>
        <FormRow v-if="selectedLeaveTypeRequiresAttachment" :label="t('leave.attachment_url')" :errors="errors?.attachment_url">
          <TextInput v-model="form.attachment_url" :placeholder="'https://...'" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- Cancel confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="cancelDialogVisible"
      :title="t('leave.cancel_request')"
      :message="t('leave.confirm_cancel_request')"
      :loading="cancelling"
      :errorMsg="cancelError"
      :confirmLabel="t('leave.cancel_request')"
      @confirm="handleCancel"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const loading = ref(true)
const balances = ref([])
const balancesLoading = ref(false)
const items = ref([])
const requestsLoading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const calendarEntries = ref([])

const leaveTypes = ref([])
const reasons = ref([])

const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const cancelDialogVisible = ref(false)
const cancelling = ref(false)
const cancelError = ref('')
const cancelTarget = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

const canCreate = computed(() => hasPermission('leave.create'))
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const leaveTypeOptions = computed(() =>
  leaveTypes.value
    .filter(lt => lt.is_active)
    .map(lt => ({ label: lt.name, value: lt.id }))
)

const reasonOptions = computed(() =>
  reasons.value.map(r => ({ label: r.name, value: r.id }))
)

const selectedLeaveTypeRequiresAttachment = computed(() => {
  const lt = leaveTypes.value.find(x => x.id === form.value.leave_type_id)
  return !!lt?.requires_attachment
})

const durationModeOptions = computed(() => [
  { label: t('leave.duration_full_day'), value: 'FULL_DAY' },
  { label: t('leave.duration_half_day_am'), value: 'HALF_DAY_AM' },
  { label: t('leave.duration_half_day_pm'), value: 'HALF_DAY_PM' },
  { label: t('leave.duration_hourly'), value: 'HOURLY' }
])

function defaultForm() {
  return {
    leave_type_id: '',
    request_start_date: '',
    request_end_date: '',
    duration_mode: 'FULL_DAY',
    start_time: '',
    end_time: '',
    leave_reason_id: null,
    leave_reason_note: '',
    attachment_url: ''
  }
}

function formatDays(n) {
  const num = Number(n) || 0
  return Number.isInteger(num) ? String(num) : num.toFixed(2).replace(/0+$/, '').replace(/\.$/, '')
}

function leaveTypeName(id) {
  return leaveTypes.value.find(x => x.id === id)?.name || '—'
}

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED_FINAL': return 'success'
    case 'REJECTED_FINAL': return 'danger'
    case 'PENDING_APPROVAL': return 'warn'
    case 'SUBMITTED': return 'info'
    case 'CANCELLED': return 'secondary'
    case 'DRAFT': return 'secondary'
    default: return 'secondary'
  }
}

function canCancel(status) {
  return ['DRAFT', 'SUBMITTED', 'PENDING_APPROVAL'].includes(status)
}

function monthRange(date) {
  const from = new Date(date.getFullYear(), date.getMonth(), 1)
  const to = new Date(date.getFullYear(), date.getMonth() + 1, 0)
  const fmt = d => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
  return { from: fmt(from), to: fmt(to) }
}

async function loadReferences() {
  const [typesRes, reasonsRes] = await Promise.all([
    api.get('/api/v1/tenant/leave/types', { params: { page: 1, per_page: 100 } }),
    api.get('/api/v1/tenant/leave/reasons')
  ])
  leaveTypes.value = typesRes.data?.data || []
  reasons.value = reasonsRes.data?.data || []
}

async function loadBalances() {
  balancesLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/leave/balances', { params: { employee_id: employeeId.value, page: 1, per_page: 100 } })
    balances.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    balancesLoading.value = false
  }
}

async function loadRequests() {
  requestsLoading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value, employee_id: employeeId.value }
    const res = await api.get('/api/v1/tenant/leave/requests', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    requestsLoading.value = false
  }
}

async function loadCalendar() {
  try {
    const { from, to } = monthRange(new Date())
    const res = await api.get('/api/v1/tenant/leave/calendar', {
      params: { employee_id: employeeId.value, from, to }
    })
    calendarEntries.value = res.data?.data || []
  } catch (e) {
    calendarEntries.value = []
  }
}

async function loadAll() {
  loading.value = true
  try {
    employeeId.value = await loadMyEmployeeId()
    if (!employeeId.value) return
    await loadReferences()
    await Promise.all([loadBalances(), loadRequests(), loadCalendar()])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadRequests()
}

function openDialog() {
  errors.value = {}
  form.value = defaultForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = defaultForm()
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  const fieldErrors = {}
  if (!form.value.leave_type_id) fieldErrors.leave_type_id = t('form.required')
  if (!form.value.request_start_date) fieldErrors.request_start_date = t('form.required')
  if (!form.value.request_end_date) fieldErrors.request_end_date = t('form.required')
  if (Object.keys(fieldErrors).length > 0) { errors.value = fieldErrors; return }

  saving.value = true
  try {
    const payload = {
      employee_id: employeeId.value,
      leave_type_id: form.value.leave_type_id,
      request_start_date: form.value.request_start_date,
      request_end_date: form.value.request_end_date,
      duration_mode: form.value.duration_mode || 'FULL_DAY',
      leave_reason_id: form.value.leave_reason_id || null,
      leave_reason_note: form.value.leave_reason_note?.trim() || null,
      attachment_url: form.value.attachment_url?.trim() || null
    }
    if (form.value.duration_mode === 'HOURLY') {
      payload.start_time = form.value.start_time || null
      payload.end_time = form.value.end_time || null
    }
    await api.post('/api/v1/tenant/leave/requests', payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('leave.request_created'), life: 3000 })
    dialogVisible.value = false
    await Promise.all([loadBalances(), loadRequests(), loadCalendar()])
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

function confirmCancel(item) {
  cancelTarget.value = item
  cancelError.value = ''
  cancelDialogVisible.value = true
}

async function handleCancel() {
  cancelling.value = true
  cancelError.value = ''
  try {
    await api.put(`/api/v1/tenant/leave/requests/${cancelTarget.value.id}/status`, { status: 'CANCELLED' })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('leave.request_cancelled'), life: 3000 })
    cancelDialogVisible.value = false
    await Promise.all([loadBalances(), loadRequests(), loadCalendar()])
  } catch (e) {
    cancelError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    cancelling.value = false
  }
}

onMounted(loadAll)
</script>
