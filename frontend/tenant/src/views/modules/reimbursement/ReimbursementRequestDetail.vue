<template>
  <div v-if="loading" class="flex items-center justify-center h-40">
    <i class="pi pi-spinner pi-spin text-2xl text-emerald-500"></i>
  </div>

  <div v-else-if="request" class="space-y-4">
    <!-- ── Header ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-start justify-between flex-wrap gap-3">
        <div>
          <div class="flex items-center gap-2 flex-wrap">
            <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ request.title }}</h2>
            <Tag :value="statusLabel(request.status)" :severity="statusSeverity(request.status)" class="!text-xs !px-1.5 !py-0.5" />
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
            {{ requestTypeName(request.request_type_id) }} · {{ request.currency || 'IDR' }}
            <template v-if="request.submitted_at"> · {{ t('reimbursement.submitted_at') }} {{ formatDate(request.submitted_at, locale) }}</template>
          </p>
        </div>
        <div class="flex items-center gap-2">
          <Button v-if="request.status === 'DRAFT' && canManage" :label="t('reimbursement.submit')" icon="pi pi-send" size="small" :loading="submitting" @click="handleSubmit" />
          <Button v-if="request.status === 'DRAFT' && canManage" :label="t('reimbursement.edit_request')" icon="pi pi-pencil" size="small" severity="secondary" outlined @click="openEditDialog" />
          <Button v-if="['DRAFT', 'SUBMITTED', 'PENDING_APPROVAL'].includes(request.status) && canManage" :label="t('reimbursement.cancel_request')" icon="pi pi-times" size="small" severity="danger" outlined :loading="cancelling" @click="confirmCancel" />
          <Button v-if="request.status === 'APPROVED' && canPay" :label="t('reimbursement.pay')" icon="pi pi-dollar" size="small" severity="success" :loading="paying" @click="confirmPay" />
        </div>
      </div>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mt-4">
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
          <p class="text-xs text-gray-400">{{ t('reimbursement.total_amount') }}</p>
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ formatCurrency(request.total_amount, request.currency) }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
          <p class="text-xs text-gray-400">{{ t('reimbursement.request_type') }}</p>
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ requestTypeName(request.request_type_id) }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
          <p class="text-xs text-gray-400">{{ t('reimbursement.currency') }}</p>
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ request.currency || 'IDR' }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 px-3 py-2.5">
          <p class="text-xs text-gray-400">{{ t('reimbursement.paid_amount') }}</p>
          <p class="text-sm font-semibold text-emerald-600 dark:text-emerald-400">{{ request.paid_amount != null ? formatCurrency(request.paid_amount, request.currency) : '-' }}</p>
        </div>
      </div>
    </div>

    <!-- ── Info ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-3">
      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('reimbursement.request_info') }}</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-2 text-sm">
        <div><span class="text-gray-400">{{ t('reimbursement.field_title') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ request.title }}</span></div>
        <div><span class="text-gray-400">{{ t('reimbursement.description_field') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ request.description || '-' }}</span></div>
        <div><span class="text-gray-400">{{ t('reimbursement.request_type') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ requestTypeName(request.request_type_id) }}</span></div>
        <div><span class="text-gray-400">{{ t('reimbursement.total_amount') }}</span><br><span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatCurrency(request.total_amount, request.currency) }}</span></div>
        <div v-if="request.submitted_at"><span class="text-gray-400">{{ t('reimbursement.submitted_at') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ formatDate(request.submitted_at, locale) }}</span></div>
        <div v-if="request.paid_at"><span class="text-gray-400">{{ t('reimbursement.paid_at') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ formatDate(request.paid_at, locale) }}</span></div>
        <div v-if="request.supervisor_note || request.hr_note"><span class="text-gray-400">{{ t('reimbursement.approval_note') }}</span><br><span class="text-gray-700 dark:text-gray-200">{{ request.supervisor_note || request.hr_note }}</span></div>
      </div>
      <Message v-if="request.status === 'APPROVED'" severity="info" :closable="false" class="!mt-1">{{ t('reimbursement.manual_pay_hint') }}</Message>
      <Message v-if="request.status !== 'DRAFT'" severity="info" :closable="false" class="!mt-1">{{ t('reimbursement.draft_edit_hint') }}</Message>
    </div>

    <!-- ── Items ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-gray-800">
        <div class="flex items-center gap-2">
          <i class="pi pi-receipt text-indigo-500 text-sm"></i>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('reimbursement.items') }}</h3>
          <span v-if="items.length" class="text-xs text-gray-400">{{ formatCurrency(totalItemsAmount, request.currency) }}</span>
        </div>
        <Button v-if="request.status === 'DRAFT' && canManage" :label="t('reimbursement.add_item')" icon="pi pi-plus" size="small" @click="openItemDialog()" />
      </div>
      <div v-if="items.length" class="divide-y divide-gray-100 dark:divide-gray-800">
        <div v-for="item in items" :key="item.id" class="px-4 py-3 flex items-center justify-between text-sm gap-3">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <span class="text-gray-700 dark:text-gray-200 font-medium">{{ item.expense_type }}</span>
              <span class="text-gray-400 dark:text-gray-500">{{ formatDate(item.expense_date, locale) }}</span>
              <span class="text-gray-400 dark:text-gray-500">{{ formatCurrency(item.amount, request.currency) }}</span>
            </div>
            <p v-if="item.description" class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ item.description }}</p>
            <a v-if="item.receipt_url" :href="item.receipt_url" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline mt-0.5 inline-block">
              <i class="pi pi-paperclip mr-1"></i>{{ t('reimbursement.receipt') }}
            </a>
          </div>
          <div v-if="request.status === 'DRAFT' && canManage" class="flex items-center gap-1 shrink-0">
            <Button icon="pi pi-paperclip" size="small" text :loading="uploadingFor === item.id" v-tooltip.top="t('reimbursement.upload_receipt')" @click="triggerReceiptUpload(item)" />
            <Button icon="pi pi-pencil" size="small" text @click="openItemDialog(item)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" @click="handleDeleteItem(item)" />
          </div>
        </div>
      </div>
      <p v-else class="text-xs text-gray-400 px-4 py-4">{{ t('reimbursement.items_empty') }}</p>
    </div>

    <!-- ── Dialog: Edit Request ── -->
    <Dialog v-model:visible="editDialogVisible" :header="t('reimbursement.edit_request')" modal :style="{ width: '520px' }">
      <div class="space-y-3">
        <FormRow :label="t('reimbursement.field_title')" required :errors="editErrors?.title">
          <TextInput v-model="editForm.title" />
        </FormRow>
        <FormRow :label="t('reimbursement.description_field')" :errors="editErrors?.description">
          <TextInput v-model="editForm.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('reimbursement.currency')">
          <Select v-model="editForm.currency" :options="currencyOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="editDialogVisible = false" />
        <Button :label="t('reimbursement.save_changes')" size="small" :loading="savingEdit" @click="handleSaveEdit" />
      </template>
    </Dialog>

    <!-- ── Dialog: Item ── -->
    <Dialog v-model:visible="itemDialogVisible" :header="editingItemId ? t('reimbursement.edit_item') : t('reimbursement.add_item')" modal :style="{ width: '480px' }">
      <div class="space-y-3">
        <FormRow :label="t('reimbursement.expense_date')" required :errors="itemErrors?.expense_date">
          <DateInput v-model="itemForm.expense_date" />
        </FormRow>
        <FormRow :label="t('reimbursement.expense_type')" required :errors="itemErrors?.expense_type">
          <TextInput v-model="itemForm.expense_type" :placeholder="t('reimbursement.expense_type_placeholder')" />
        </FormRow>
        <FormRow :label="t('reimbursement.amount')" required :errors="itemErrors?.amount">
          <InputNumber v-model="itemForm.amount" class="!w-full" :min="0" size="small" />
        </FormRow>
        <FormRow :label="t('reimbursement.description_field')" :errors="itemErrors?.description">
          <TextInput v-model="itemForm.description" textarea :rows="2" />
        </FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="itemDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingItem" @click="handleSaveItem" />
      </template>
    </Dialog>

    <!-- ── Confirm: Cancel ── -->
    <ConfirmActionDialog
      v-model:visible="cancelDialogVisible"
      :title="t('reimbursement.confirm_cancel_title')"
      :message="t('reimbursement.confirm_cancel')"
      :loading="cancelling"
      :errorMsg="cancelError"
      :confirm-label="t('reimbursement.cancel_request')"
      severity="danger"
      icon="pi pi-exclamation-triangle"
      @confirm="handleCancel"
    />

    <!-- ── Confirm: Pay ── -->
    <ConfirmActionDialog
      v-model:visible="payDialogVisible"
      :title="t('reimbursement.confirm_pay_title')"
      :message="t('reimbursement.confirm_pay')"
      :loading="paying"
      :errorMsg="payError"
      :confirm-label="t('reimbursement.pay')"
      severity="success"
      icon="pi pi-dollar"
      @confirm="handlePay"
    />

    <input ref="receiptFileInputRef" type="file" class="hidden" @change="onReceiptFileSelected" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import InputNumber from 'primevue/inputnumber'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'

const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const requestId = route.params.id
const loading = ref(true)
const request = ref(null)
const items = ref([])
const requestTypes = ref([])

const submitting = ref(false)
const cancelling = ref(false)
const cancelDialogVisible = ref(false)
const cancelError = ref('')
const paying = ref(false)
const payDialogVisible = ref(false)
const payError = ref('')

const editDialogVisible = ref(false)
const savingEdit = ref(false)
const editErrors = ref({})
const editForm = ref({ title: '', description: '', currency: 'IDR' })

const itemDialogVisible = ref(false)
const savingItem = ref(false)
const editingItemId = ref(null)
const itemErrors = ref({})
const itemForm = ref(defaultItemForm())

const receiptFileInputRef = ref(null)
const receiptUploadTarget = ref(null)
const uploadingFor = ref(null)

const canManage = computed(() => hasPermission('reimbursement.update') || hasPermission('reimbursement.create'))
const canPay = computed(() => hasPermission('reimbursement.approve'))

const currencyOptions = [
  { label: 'IDR', value: 'IDR' },
  { label: 'USD', value: 'USD' }
]

const totalItemsAmount = computed(() => items.value.reduce((sum, i) => sum + (Number(i.amount) || 0), 0))

function defaultItemForm() {
  return { expense_date: '', expense_type: '', description: '', amount: null }
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

async function loadRequest() {
  const res = await api.get(`/api/v1/tenant/reimbursements/requests/${requestId}`)
  request.value = res.data?.data || null
}
async function loadItems() {
  try {
    const res = await api.get(`/api/v1/tenant/reimbursements/requests/${requestId}/items`)
    items.value = res.data?.data || []
  } catch {
    items.value = []
  }
}
async function loadRequestTypes() {
  try {
    const res = await api.get('/api/v1/tenant/reimbursements/types', { params: { page: 1, per_page: 100 } })
    requestTypes.value = res.data?.data || []
  } catch {
    requestTypes.value = []
  }
}

async function loadAll() {
  loading.value = true
  try {
    await loadRequest()
    await Promise.all([loadItems(), loadRequestTypes()])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

// ── Status actions ──
async function handleSubmit() {
  submitting.value = true
  try {
    await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}/status`, { status: 'SUBMITTED' })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadRequest()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    submitting.value = false
  }
}

function confirmCancel() {
  cancelError.value = ''
  cancelDialogVisible.value = true
}
async function handleCancel() {
  cancelling.value = true
  cancelError.value = ''
  try {
    await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}/status`, { status: 'CANCELLED' })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    cancelDialogVisible.value = false
    await loadRequest()
  } catch (e) {
    cancelError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    cancelling.value = false
  }
}

function confirmPay() {
  payError.value = ''
  payDialogVisible.value = true
}
async function handlePay() {
  paying.value = true
  payError.value = ''
  try {
    await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}/status`, { status: 'PAID' })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    payDialogVisible.value = false
    await loadRequest()
  } catch (e) {
    payError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    paying.value = false
  }
}

// ── Edit request (DRAFT only) ──
function openEditDialog() {
  editErrors.value = {}
  editForm.value = {
    title: request.value.title || '',
    description: request.value.description || '',
    currency: request.value.currency || 'IDR'
  }
  editDialogVisible.value = true
}
async function handleSaveEdit() {
  editErrors.value = {}
  if (!editForm.value.title?.trim()) { editErrors.value = { title: t('form.required') }; return }
  savingEdit.value = true
  try {
    await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}`, {
      title: editForm.value.title.trim(),
      description: editForm.value.description?.trim() || '',
      currency: editForm.value.currency || 'IDR'
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    editDialogVisible.value = false
    await loadRequest()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      editErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    savingEdit.value = false
  }
}

// ── Items ──
function openItemDialog(existing) {
  editingItemId.value = existing?.id || null
  itemErrors.value = {}
  itemForm.value = existing
    ? { expense_date: existing.expense_date || '', expense_type: existing.expense_type || '', description: existing.description || '', amount: existing.amount }
    : defaultItemForm()
  itemDialogVisible.value = true
}
async function handleSaveItem() {
  itemErrors.value = {}
  if (!itemForm.value.expense_date) { itemErrors.value = { expense_date: t('form.required') }; return }
  if (!itemForm.value.expense_type?.trim()) { itemErrors.value = { expense_type: t('form.required') }; return }
  if (itemForm.value.amount === null || itemForm.value.amount === undefined) { itemErrors.value = { amount: t('form.required') }; return }
  savingItem.value = true
  try {
    const payload = {
      expense_date: itemForm.value.expense_date,
      expense_type: itemForm.value.expense_type.trim(),
      description: itemForm.value.description?.trim() || '',
      amount: itemForm.value.amount
    }
    if (editingItemId.value) {
      await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}/items/${editingItemId.value}`, payload)
    } else {
      await api.post(`/api/v1/tenant/reimbursements/requests/${requestId}/items`, payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    itemDialogVisible.value = false
    await Promise.all([loadItems(), loadRequest()])
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      itemErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    savingItem.value = false
  }
}
async function handleDeleteItem(item) {
  try {
    await api.delete(`/api/v1/tenant/reimbursements/requests/${requestId}/items/${item.id}`)
    await Promise.all([loadItems(), loadRequest()])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

// ── Receipt upload (two-step: generic upload → attach URL to item) ──
function triggerReceiptUpload(item) {
  receiptUploadTarget.value = item
  receiptFileInputRef.value?.click()
}
async function onReceiptFileSelected(event) {
  const file = event.target.files?.[0]
  const item = receiptUploadTarget.value
  if (!file || !item) return
  uploadingFor.value = item.id
  try {
    const fd = new FormData()
    fd.append('file', file)
    const uploadRes = await api.post('/api/v1/tenant/uploads', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
    const filePath = uploadRes.data?.data?.url || ''
    if (!filePath) throw new Error('upload failed')
    await api.put(`/api/v1/tenant/reimbursements/requests/${requestId}/items/${item.id}`, { receipt_url: filePath })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadItems()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    uploadingFor.value = null
    receiptUploadTarget.value = null
    if (receiptFileInputRef.value) receiptFileInputRef.value.value = ''
  }
}

onMounted(loadAll)
</script>
