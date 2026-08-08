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
          <i class="pi pi-clock text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.overtime_empty') }}</p>
        </div>
      </template>
      <Column field="work_date" :header="t('attendance.work_date')" style="width:130px" />
      <Column field="requested_minutes" :header="t('attendance.requested_minutes')" style="width:130px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.requested_minutes }} min</span></template>
      </Column>
      <Column field="calculated_minutes" :header="t('attendance.calculated_minutes')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.calculated_minutes ?? '-' }}</span></template>
      </Column>
      <Column field="reason" :header="t('attendance.reason')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.reason || '-' }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:150px">
        <template #body="{data}"><Tag :value="data.status" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
    </DataTable>

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
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
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
const form = ref({ work_date: '', start_time: '', end_time: '', requested_minutes: null, reason: '' })

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// ── Perhitungan menit otomatis + dukungan lintas hari ──
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

// end_time <= start_time → lembur lintas tengah malam (hari berikutnya)
const isCrossDay = computed(() => {
  const s = timeToMinutes(form.value.start_time)
  const e = timeToMinutes(form.value.end_time)
  return s !== null && e !== null && e <= s
})

const autoMinutes = computed(() => {
  const s = timeToMinutes(form.value.start_time)
  const e = timeToMinutes(form.value.end_time)
  if (s === null || e === null) return null
  let diff = e - s
  if (diff <= 0) diff += 1440
  return diff
})

// Isi otomatis saat start/end lengkap (tetap bisa diedit manual)
watch(autoMinutes, (m) => {
  if (m !== null && dialogVisible.value) form.value.requested_minutes = m
})

const skeletonColumns = [
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' }
]

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

onMounted(async () => {
  await loadMyEmployeeId()
  await loadData()
})
</script>
