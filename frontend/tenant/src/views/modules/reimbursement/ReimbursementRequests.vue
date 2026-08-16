<template>
  <div class="space-y-4">
    <!-- ── Skeleton loading ── -->
    <div v-if="loading" class="space-y-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between flex-wrap gap-2 mb-3 animate-pulse">
          <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-44"></div>
          <div class="h-8 bg-gray-200 dark:bg-gray-700 rounded w-28"></div>
        </div>
        <SkeletonTable :columns="skeletonColumns" :rows="8" />
      </div>
    </div>

    <template v-else-if="!employeeId && !isAdmin">
      <Message severity="warn" :closable="false">{{ t('reimbursement.no_employee_linked') }}</Message>
    </template>

    <template v-else>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between flex-wrap gap-2 mb-3">
          <div class="flex items-center gap-2">
            <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('reimbursement.back_to_index')" @click="router.push('/reimbursements')" />
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">
              {{ isAll ? t('reimbursement.requests') : t('reimbursement.my_requests') }}
            </h2>
          </div>
          <div class="flex items-center gap-2">
            <Button v-if="!isAll && canCreate" :label="t('reimbursement.new_request')" icon="pi pi-plus" size="small" @click="openDialog()" />
          </div>
        </div>

        <div class="flex items-center gap-2 flex-wrap mb-3">
          <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" showClear :placeholder="t('common.status')" class="w-44" @change="reload" />
          <Select
            v-if="isAll && isAdmin"
            v-model="employeeFilter"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            showClear
            filter
            :placeholder="t('reimbursement.all_employees')"
            class="w-64"
            @change="reload"
          />
        </div>

        <SkeletonTable v-if="listLoading" :columns="skeletonColumns" :rows="8" />
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
              <i class="pi pi-credit-card text-3xl mb-2 opacity-50"></i>
              <p class="text-sm font-medium">{{ t('reimbursement.requests_empty') }}</p>
            </div>
          </template>
          <Column field="title" :header="t('reimbursement.field_title')">
            <template #body="{data}">
              <a class="text-indigo-600 dark:text-indigo-400 hover:underline font-medium cursor-pointer" @click="router.push(`/reimbursements/${data.id}`)">{{ data.title }}</a>
            </template>
          </Column>
          <Column :header="t('reimbursement.request_type')" style="width:160px">
            <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ requestTypeName(data.request_type_id) }}</span></template>
          </Column>
          <Column :header="t('reimbursement.total_amount')" style="width:150px">
            <template #body="{data}">
              <span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatCurrency(data.total_amount, data.currency) }}</span>
              <span v-if="data.status === 'PAID' && data.paid_amount != null" class="block text-[11px] text-emerald-600 dark:text-emerald-400">{{ t('reimbursement.status_paid') }}: {{ formatCurrency(data.paid_amount, data.currency) }}</span>
            </template>
          </Column>
          <Column field="status" :header="t('common.status')" style="width:150px">
            <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
          </Column>
          <Column :header="t('reimbursement.submitted_at')" style="width:150px">
            <template #body="{data}">
              <span class="text-gray-500 dark:text-gray-400">{{ data.submitted_at ? formatDate(data.submitted_at, locale) : '-' }}</span>
            </template>
          </Column>
          <Column :header="t('common.actions')" style="width:110px" frozen alignFrozen="right">
            <template #body="{data}">
              <div class="flex items-center gap-0.5 justify-end">
                <Button v-if="data.status === 'APPROVED' && canPay" icon="pi pi-dollar" size="small" text severity="success" v-tooltip.top="t('reimbursement.pay')" :loading="payingFor === data.id" @click="confirmPay(data)" />
                <Button icon="pi pi-eye" size="small" text @click="router.push(`/reimbursements/${data.id}`)" />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
    </template>

    <!-- ── Dialog: New Reimbursement Request ── -->
    <Dialog v-model:visible="dialogVisible" :header="t('reimbursement.new_request')" modal :style="{ width: '520px' }" @hide="resetForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('reimbursement.new_request_description') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('reimbursement.request_type')" required :errors="errors?.request_type_id">
          <Select v-model="form.request_type_id" :options="requestTypeOptions" optionLabel="name" optionValue="id" filter class="w-full" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('reimbursement.field_title')" required :errors="errors?.title">
          <TextInput v-model="form.title" :placeholder="t('reimbursement.title_placeholder')" />
        </FormRow>
        <FormRow :label="t('reimbursement.description_field')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('reimbursement.currency')">
          <Select v-model="form.currency" :options="currencyOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Payment ── -->
    <Dialog v-model:visible="payDialogVisible" :header="t('reimbursement.pay_title')" modal :style="{ width: '480px' }" @hide="resetPayForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('reimbursement.pay_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('reimbursement.total_amount')" required :errors="payErrors?.amount">
          <InputNumber v-model="payForm.amount" class="!w-full" :min="0" size="small" />
        </FormRow>
        <FormRow :label="t('reimbursement.payment_method')" :errors="payErrors?.payment_method">
          <Select v-model="payForm.payment_method" :options="paymentMethodOptions" optionLabel="label" optionValue="value" showClear class="w-full" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('reimbursement.payment_reference')" :errors="payErrors?.payment_reference">
          <TextInput v-model="payForm.payment_reference" :placeholder="t('reimbursement.payment_reference_placeholder')" />
        </FormRow>
        <FormRow :label="t('reimbursement.payment_note')" :errors="payErrors?.payment_note">
          <TextInput v-model="payForm.payment_note" textarea :rows="2" />
        </FormRow>
        <Message v-if="payError" severity="error" :closable="false" class="!text-xs">{{ payError }}</Message>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="payDialogVisible = false" />
          <Button :label="t('reimbursement.pay')" icon="pi pi-dollar" severity="success" size="small" :loading="paying" :disabled="paying" @click="handlePay" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import InputNumber from 'primevue/inputnumber'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const loading = ref(true)
const listLoading = ref(false)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const statusFilter = ref(null)
const employeeFilter = ref(null)

const requestTypes = ref([])
const employees = ref([])

const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const payDialogVisible = ref(false)
const paying = ref(false)
const payingFor = ref(null)
const payError = ref('')
const payErrors = ref({})
const payForm = ref(defaultPayForm())
const payTarget = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
// Route mode: /reimbursements/all → semua; /reimbursements/my-requests → milik sendiri.
const isAll = computed(() => route.name === 'ReimbursementAllRequests')
// Admin/HR (who can approve) see all requests; employees default to their own.
const isAdmin = computed(() => hasPermission('reimbursement.approve'))
const canCreate = computed(() => hasPermission('reimbursement.create'))
const canPay = computed(() => hasPermission('reimbursement.approve'))
// Non-admin yang kebetulan membuka mode "all" tetap di-scope ke miliknya sendiri.
const scopeMine = computed(() => !isAll.value || !isAdmin.value)

const statusOptions = [
  { label: 'DRAFT', value: 'DRAFT' },
  { label: 'SUBMITTED', value: 'SUBMITTED' },
  { label: 'PENDING_APPROVAL', value: 'PENDING_APPROVAL' },
  { label: 'APPROVED', value: 'APPROVED' },
  { label: 'REJECTED', value: 'REJECTED' },
  { label: 'PAID', value: 'PAID' },
  { label: 'CANCELLED', value: 'CANCELLED' }
]

const currencyOptions = [
  { label: 'IDR', value: 'IDR' },
  { label: 'USD', value: 'USD' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-20' }
]

const requestTypeOptions = computed(() =>
  requestTypes.value.filter(rt => rt.is_active).map(rt => ({ id: rt.id, name: rt.name }))
)

const employeeOptions = computed(() =>
  employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.employee_id }))
)

function defaultForm() {
  return { request_type_id: '', title: '', description: '', currency: 'IDR' }
}

const paymentMethodOptions = [
  { label: t('reimbursement.payment_method_bank_transfer'), value: 'BANK_TRANSFER' },
  { label: t('reimbursement.payment_method_cash'), value: 'CASH' },
  { label: t('reimbursement.payment_method_cheque'), value: 'CHEQUE' }
]

function defaultPayForm() {
  return { amount: null, payment_method: 'BANK_TRANSFER', payment_reference: '', payment_note: '' }
}

function resetPayForm() {
  payForm.value = defaultPayForm()
  payErrors.value = {}
  payError.value = ''
  payTarget.value = null
}

function confirmPay(data) {
  payTarget.value = data
  payErrors.value = {}
  payError.value = ''
  payForm.value = {
    amount: data.total_amount ?? null,
    payment_method: 'BANK_TRANSFER',
    payment_reference: '',
    payment_note: ''
  }
  payDialogVisible.value = true
}

async function handlePay() {
  payErrors.value = {}
  payError.value = ''
  if (!payTarget.value) return
  const amount = Number(payForm.value.amount)
  if (!amount || amount <= 0) {
    payErrors.value = { amount: t('form.required') }
    return
  }
  paying.value = true
  payingFor.value = payTarget.value.id
  try {
    const payload = {
      status: 'PAID',
      amount,
      payment_method: payForm.value.payment_method || 'BANK_TRANSFER',
      payment_reference: payForm.value.payment_reference?.trim() || '',
      payment_note: payForm.value.payment_note?.trim() || ''
    }
    await api.put(`/api/v1/tenant/reimbursements/requests/${payTarget.value.id}/status`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    payDialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      payErrors.value = fieldErrors
    } else {
      payError.value = getErrorMessage(e, t('message.operation_failed'))
    }
  } finally {
    paying.value = false
    payingFor.value = null
  }
}

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'PAID': return 'success'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'secondary'
    case 'SUBMITTED': return 'info'
    case 'PENDING_APPROVAL': return 'warn'
    case 'DRAFT': return 'secondary'
    default: return 'secondary'
  }
}

function statusLabel(status) {
  const key = `reimbursement.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function requestTypeName(id) {
  return requestTypes.value.find(rt => rt.id === id)?.name || id?.slice(0, 8) || '-'
}

function formatCurrency(v, currency = 'IDR') {
  if (v === null || v === undefined) return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: currency || 'IDR', maximumFractionDigits: 0 }).format(v)
}

async function loadReferences() {
  try {
    const [typesRes, empRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/reimbursements/types', { params: { page: 1, per_page: 100 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    ])
    requestTypes.value = typesRes.status === 'fulfilled' ? (typesRes.value.data?.data || []) : []
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
  } catch {
    // fail-silent — dropdowns kosong, list tetap bisa dimuat
  }
}

async function loadData() {
  listLoading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (scopeMine.value) {
      if (!employeeId.value) return
      params.employee_id = employeeId.value
    } else if (employeeFilter.value) {
      params.employee_id = employeeFilter.value
    }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/reimbursements/requests', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    listLoading.value = false
  }
}

function reload() {
  currentPage.value = 1
  loadData()
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
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
  if (!form.value.request_type_id) { errors.value = { request_type_id: t('form.required') }; return }
  if (!form.value.title?.trim()) { errors.value = { title: t('form.required') }; return }
  saving.value = true
  try {
    const res = await api.post('/api/v1/tenant/reimbursements/requests', {
      request_type_id: form.value.request_type_id,
      title: form.value.title.trim(),
      description: form.value.description?.trim() || '',
      currency: form.value.currency || 'IDR'
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    const newId = res.data?.data?.id
    if (newId) {
      router.push(`/reimbursements/${newId}`)
    } else {
      await loadData()
    }
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

async function loadAll() {
  loading.value = true
  try {
    employeeId.value = await loadMyEmployeeId()
    await loadReferences()
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>
