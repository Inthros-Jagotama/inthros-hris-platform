<template>
  <div class="space-y-4">
    <!-- ── Header actions ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
        <Button
          :label="t('common.refresh')"
          icon="pi pi-refresh"
          size="small"
          severity="secondary"
          outlined
          :loading="loading"
          class="!whitespace-nowrap shrink-0"
          @click="loadData"
        />
      </div>
      <Button :label="t('payroll.new_run')" icon="pi pi-plus" size="small" @click="openCreateDialog" />
    </div>

    <!-- ── Summary cards ── -->
    <div v-if="items.length" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 bg-white dark:bg-gray-800">
        <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_earning') }}</p>
        <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-1">{{ formatMoney(sum(items, 'total_earning')) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 bg-white dark:bg-gray-800">
        <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_deduction') }}</p>
        <p class="text-sm font-semibold text-rose-600 dark:text-rose-400 mt-1">{{ formatMoney(sum(items, 'total_deduction')) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 bg-white dark:bg-gray-800">
        <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_employer_contribution') }}</p>
        <p class="text-sm font-semibold text-amber-600 dark:text-amber-400 mt-1">{{ formatMoney(sum(items, 'total_employer_contribution')) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 bg-white dark:bg-gray-800">
        <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_net') }}</p>
        <p class="text-sm font-semibold text-emerald-600 dark:text-emerald-400 mt-1">{{ formatMoney(sum(items, 'total_net')) }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 bg-white dark:bg-gray-800">
        <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_company_cost') }}</p>
        <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-1">{{ formatMoney(sum(items, 'total_company_cost')) }}</p>
      </div>
    </div>

    <!-- ── Runs table ── -->
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
      sortField="created_at"
      :sortOrder="-1"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-dollar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('payroll.runs_empty') }}</p>
        </div>
      </template>
      <Column field="run_code" :header="t('payroll.run_code')" sortable style="width:150px">
        <template #body="{ data }">
          <button class="text-emerald-600 dark:text-emerald-400 font-medium hover:underline text-left" @click="openRun(data)">{{ data.run_code }}</button>
        </template>
      </Column>
      <Column field="payroll_period_id" :header="t('payroll.period')" sortable style="width:120px">
        <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ periodLabel(data.payroll_period_id) }}</span></template>
      </Column>
      <Column field="run_type" :header="t('payroll.run_type')" sortable style="width:110px">
        <template #body="{ data }"><Tag :value="t('payroll.run_type_' + data.run_type.toLowerCase())" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable style="width:120px">
        <template #body="{ data }"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="total_employees" :header="t('payroll.total_employees')" sortable style="width:90px">
        <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.total_employees }}</span></template>
      </Column>
      <Column field="total_net" :header="t('payroll.total_net')" sortable style="width:150px">
        <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ formatMoney(data.total_net) }}</span></template>
      </Column>
      <Column field="calculated_at" :header="t('payroll.calculated')" sortable style="width:130px">
        <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ data.calculated_at ? formatDateTime(data.calculated_at) : '-' }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:180px" frozen alignFrozen="right">
        <template #body="{ data }">
          <div class="flex items-center gap-1 justify-end flex-wrap">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.details')" @click="openRun(data)" />
            <Button
              v-if="data.status === 'DRAFT'"
              icon="pi pi-calculator"
              size="small"
              text
              severity="info"
              v-tooltip.left="t('payroll.calculate')"
              @click="confirmCalculate(data)"
            />
            <Button
              v-if="data.status === 'CALCULATED' || data.status === 'APPROVED'"
              icon="pi pi-check"
              size="small"
              text
              severity="success"
              v-tooltip.left="t('payroll.status_reviewed')"
              @click="confirmStatus(data, 'REVIEWED')"
            />
            <Button
              v-if="data.status === 'REVIEWED'"
              icon="pi pi-check-circle"
              size="small"
              text
              severity="success"
              v-tooltip.left="t('payroll.status_approved')"
              @click="confirmStatus(data, 'APPROVED')"
            />
            <Button
              v-if="data.status === 'APPROVED'"
              icon="pi pi-lock"
              size="small"
              text
              severity="warn"
              v-tooltip.left="t('payroll.status_locked')"
              @click="confirmStatus(data, 'LOCKED')"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Create run dialog ── -->
    <Dialog v-model:visible="createDialogVisible" :header="t('payroll.new_run')" modal :style="{ width: '560px' }" @hide="resetCreateForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.select_period')" required :errors="errors?.payroll_period_id">
          <SelectLabel v-model="createForm.payroll_period_id" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.payroll_period_id }" showClear />
        </FormRow>
        <FormRow :label="t('payroll.run_type')">
          <SelectLabel v-model="createForm.run_type" :options="runTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          <small v-if="selectedRunTypeDesc" class="text-xs text-gray-400 mt-1 block"><i class="pi pi-info-circle mr-1"></i>{{ selectedRunTypeDesc }}</small>
        </FormRow>
        <FormRow :label="t('payroll.proration_method')">
          <SelectLabel v-model="createForm.proration_method" :options="prorationOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          <small v-if="selectedProrationDesc" class="text-xs text-gray-400 mt-1 block"><i class="pi pi-info-circle mr-1"></i>{{ selectedProrationDesc }}</small>
        </FormRow>
        <FormRow :label="t('payroll.select_employee')" :errors="errors?.employee_ids">
          <SelectLabel v-model="createForm.employee_ids" :options="employeeOptions" optionLabel="label" optionValue="value" multiple filter :placeholder="t('common.select')" />
          <small class="text-xs text-gray-400 mt-1 block">{{ t('payroll.employee_ids_hint') }}</small>
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="createDialogVisible = false" />
          <Button :label="t('common.create')" size="small" :loading="creating" :disabled="creating" @click="handleCreate" />
        </div>
      </template>
    </Dialog>

    <!-- ── Confirmation dialogs ── -->
    <ConfirmActionDialog
      v-model:visible="confirmVisible"
      :title="confirmTitle"
      :message="confirmMessage"
      :loading="confirmLoading"
      :confirmLabel="confirmLabel"
      @confirm="onConfirm"
      @cancel="confirmVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()
const router = useRouter()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const periods = ref([])
const employees = ref([])

const createDialogVisible = ref(false)
const creating = ref(false)
const errors = ref({})
const createForm = ref({ payroll_period_id: null, run_type: 'REGULAR', proration_method: 'CALENDAR_DAYS', employee_ids: [] })

const confirmVisible = ref(false)
const confirmLoading = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmLabel = ref('')
let confirmAction = null

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'icons', count: 4, headerWidth: 'w-24' }
]

const periodOptions = computed(() =>
  periods.value.map(p => ({ label: `${p.period_code} (${p.period_year}-${String(p.period_month).padStart(2, '0')})`, value: p.id }))
)
const employeeOptions = computed(() =>
  employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id }))
)
const runTypeOptions = computed(() =>
  ['REGULAR', 'OFF_CYCLE', 'THR', 'BONUS'].map(v => ({
    label: t(`payroll.run_type_${v.toLowerCase()}`),
    desc: t(`payroll.run_type_${v.toLowerCase()}_desc`),
    value: v
  }))
)
const prorationOptions = computed(() =>
  ['CALENDAR_DAYS', 'WORKING_DAYS', 'FIXED_30_DAYS', 'ATTENDANCE_DAYS'].map(v => ({
    label: t(`payroll.proration_${v.toLowerCase()}`),
    desc: t(`payroll.proration_${v.toLowerCase()}_desc`),
    value: v
  }))
)
const selectedRunTypeDesc = computed(() =>
  runTypeOptions.value.find(o => o.value === createForm.value.run_type)?.desc || ''
)
const selectedProrationDesc = computed(() =>
  prorationOptions.value.find(o => o.value === createForm.value.proration_method)?.desc || ''
)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function formatMoney(val) {
  const n = Number(val || 0)
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(n)
}
function formatDateTime(val) {
  const d = new Date(val)
  return `${formatDate(d, locale.value)} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}
function sum(list, key) {
  return list.reduce((acc, x) => acc + Number(x[key] || 0), 0)
}
function statusLabel(status) {
  const key = `payroll.status_${String(status || '').toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
function statusSeverity(status) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'CALCULATED': return 'info'
    case 'REVIEWED': return 'warn'
    case 'APPROVED': return 'success'
    case 'LOCKED': return 'contrast'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function periodLabel(id) {
  const p = periods.value.find(x => x.id === id)
  return p ? p.period_code : id
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/runs', { params: { page: currentPage.value, per_page: perPage.value } })
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

async function loadReferences() {
  const [pRes, eRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/payroll/periods', { params: { per_page: 200 } }),
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  ])
  periods.value = pRes.status === 'fulfilled' ? (pRes.value.data?.data || []) : []
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function openRun(run) {
  router.push(`/payroll/runs/${run.id}`)
}

function openCreateDialog() {
  errors.value = {}
  createForm.value = { payroll_period_id: null, run_type: 'REGULAR', proration_method: 'CALENDAR_DAYS', employee_ids: [] }
  createDialogVisible.value = true
}
function resetCreateForm() {
  errors.value = {}
  createForm.value = { payroll_period_id: null, run_type: 'REGULAR', proration_method: 'CALENDAR_DAYS', employee_ids: [] }
}

async function handleCreate() {
  errors.value = {}
  if (!createForm.value.payroll_period_id) { errors.value = { payroll_period_id: [t('form.required')] }; return }
  creating.value = true
  try {
    const payload = {
      payroll_period_id: createForm.value.payroll_period_id,
      run_type: createForm.value.run_type,
      proration_method: createForm.value.proration_method,
      employee_ids: createForm.value.employee_ids || []
    }
    const res = await api.post('/api/v1/tenant/payroll/runs', payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.run_created'), life: 3000 })
    createDialogVisible.value = false
    const run = res.data?.data
    if (run?.id) {
      router.push(`/payroll/runs/${run.id}`)
    } else {
      loadData()
    }
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    creating.value = false
  }
}

function confirmCalculate(run) {
  confirmTitle.value = t('payroll.calculate')
  confirmMessage.value = t('payroll.confirm_calculate')
  confirmLabel.value = t('payroll.calculate')
  confirmAction = () => calculateRun(run)
  confirmVisible.value = true
}
async function calculateRun(run) {
  await api.post(`/api/v1/tenant/payroll/runs/${run.id}/calculate`)
  toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.run_calculated'), life: 3000 })
  loadData()
}

function confirmStatus(run, status) {
  confirmTitle.value = t('payroll.update_status')
  confirmMessage.value = t('payroll.confirm_status_change', { status: statusLabel(status) })
  confirmLabel.value = statusLabel(status)
  confirmAction = () => updateStatus(run, status)
  confirmVisible.value = true
}
async function updateStatus(run, status) {
  await api.put(`/api/v1/tenant/payroll/runs/${run.id}/status`, { status })
  toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.run_status_updated'), life: 3000 })
  loadData()
}

async function onConfirm() {
  if (!confirmAction) return
  confirmLoading.value = true
  try {
    await confirmAction()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    confirmLoading.value = false
    confirmVisible.value = false
    confirmAction = null
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
