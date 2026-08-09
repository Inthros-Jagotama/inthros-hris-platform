<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('attendance.view_in_approvals')" icon="pi pi-external-link" size="small" severity="secondary" outlined @click="router.push('/approvals')" />
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
          <i class="pi pi-pencil text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.corrections_empty') }}</p>
        </div>
      </template>
      <Column field="correction_type" :header="t('attendance.correction_type')" style="width:170px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ t('attendance.correction_type_' + data.correction_type.toLowerCase()) }}</span></template>
      </Column>
      <Column field="reason" :header="t('attendance.reason')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.reason || '-' }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:150px">
        <template #body="{data}"><Tag :value="data.status" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="t('attendance.correction_new')" modal :style="{ width: '480px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('attendance.work_date')" required :errors="errors?.attendance_session_id">
          <DateInput v-model="form.work_date" @update:modelValue="resolveSession" />
          <small v-if="sessionLookupState === 'loading'" class="text-xs text-gray-400">{{ t('attendance.session_lookup_loading') }}</small>
          <small v-else-if="sessionLookupState === 'not_found'" class="text-xs text-rose-500">{{ t('attendance.session_lookup_not_found') }}</small>
        </FormRow>
        <FormRow :label="t('attendance.correction_type')" required :errors="errors?.correction_type">
          <Select v-model="form.correction_type" :options="correctionTypeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
        <FormRow v-if="needsCheckin" :label="t('attendance.requested_checkin')" required :errors="errors?.requested_checkin">
          <TimeInput v-model="form.requested_checkin" />
        </FormRow>
        <FormRow v-if="needsCheckout" :label="t('attendance.requested_checkout')" required :errors="errors?.requested_checkout">
          <TimeInput v-model="form.requested_checkout" />
        </FormRow>
        <FormRow :label="t('attendance.reason')" required :errors="errors?.reason">
          <TextInput v-model="form.reason" textarea :rows="3" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving || !form.attendance_session_id" @click="handleSave" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { localDateTimeISOString } from '@/utils/localTime'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import TimeInput from '@/components/TimeInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const sessionLookupState = ref('') // '', 'loading', 'found', 'not_found'
const form = ref({ work_date: '', attendance_session_id: '', correction_type: 'MISSING_CHECKIN', requested_checkin: '', requested_checkout: '', reason: '' })

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' }
]
const correctionTypeOptions = computed(() => [
  { label: t('attendance.correction_type_missing_checkin'), value: 'MISSING_CHECKIN' },
  { label: t('attendance.correction_type_missing_checkout'), value: 'MISSING_CHECKOUT' },
  { label: t('attendance.correction_type_wrong_checkin'), value: 'WRONG_CHECKIN' },
  { label: t('attendance.correction_type_wrong_checkout'), value: 'WRONG_CHECKOUT' }
])
const needsCheckin = computed(() => form.value.correction_type === 'MISSING_CHECKIN' || form.value.correction_type === 'WRONG_CHECKIN')
const needsCheckout = computed(() => form.value.correction_type === 'MISSING_CHECKOUT' || form.value.correction_type === 'WRONG_CHECKOUT')

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'PENDING_APPROVAL': return 'info'
    default: return 'secondary'
  }
}

async function loadData() {
  if (!employeeId.value) return
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value, employee_id: employeeId.value }
    const res = await api.get('/api/v1/tenant/attendance/corrections', { params })
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

// resolveSession resolves the session_id for the picked date — corrections
// target an existing attendance_session, not a raw date, per the backend's
// CreateCorrectionRequest (attendance_session_id is required).
async function resolveSession(workDate) {
  form.value.attendance_session_id = ''
  if (!workDate) { sessionLookupState.value = ''; return }
  sessionLookupState.value = 'loading'
  try {
    const res = await api.get('/api/v1/tenant/attendance/sessions/detail', { params: { employee_id: employeeId.value, work_date: workDate } })
    form.value.attendance_session_id = res.data?.data?.id || ''
    sessionLookupState.value = form.value.attendance_session_id ? 'found' : 'not_found'
  } catch {
    sessionLookupState.value = 'not_found'
  }
}

function openDialog() {
  errors.value = {}
  sessionLookupState.value = ''
  form.value = { work_date: '', attendance_session_id: '', correction_type: 'MISSING_CHECKIN', requested_checkin: '', requested_checkout: '', reason: '' }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { work_date: '', attendance_session_id: '', correction_type: 'MISSING_CHECKIN', requested_checkin: '', requested_checkout: '', reason: '' }
  errors.value = {}
  sessionLookupState.value = ''
}

async function handleSave() {
  errors.value = {}
  if (!form.value.attendance_session_id) { errors.value = { attendance_session_id: t('attendance.session_lookup_not_found') }; return }
  if (needsCheckin.value && !form.value.requested_checkin?.trim()) { errors.value = { requested_checkin: t('form.required') }; return }
  if (needsCheckout.value && !form.value.requested_checkout?.trim()) { errors.value = { requested_checkout: t('form.required') }; return }
  if (!form.value.reason?.trim()) { errors.value = { reason: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      employee_id: employeeId.value,
      attendance_session_id: form.value.attendance_session_id,
      correction_type: form.value.correction_type,
      reason: form.value.reason
    }
    if (needsCheckin.value) payload.requested_checkin = localDateTimeISOString(form.value.work_date, form.value.requested_checkin)
    if (needsCheckout.value) payload.requested_checkout = localDateTimeISOString(form.value.work_date, form.value.requested_checkout)
    await api.post('/api/v1/tenant/attendance/corrections', payload)
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

onMounted(async () => {
  await loadMyEmployeeId()
  await loadData()
})
</script>
