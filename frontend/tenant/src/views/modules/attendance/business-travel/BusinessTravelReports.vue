<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <DateInput v-model="fromDate" :placeholder="t('attendance.from_date')" class="!w-36" />
        <span class="text-gray-400">-</span>
        <DateInput v-model="toDate" :placeholder="t('attendance.to_date')" class="!w-36" />
        <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" showClear :placeholder="t('common.status')" class="w-44" />
        <Button :label="t('common.refresh')" icon="pi pi-search" size="small" :loading="loading" class="!whitespace-nowrap shrink-0" @click="loadData" />
      </div>
      <span v-if="currentItems.length > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ currentItems.length }} {{ t('common.items') }}</span>
    </div>

    <div class="flex items-center gap-1 flex-wrap border-b border-gray-200 dark:border-gray-700">
      <button
        v-for="tab in reportTabs"
        :key="tab.key"
        type="button"
        class="px-3 py-2 text-sm font-medium border-b-2 -mb-px transition-colors"
        :class="activeTab === tab.key
          ? 'border-indigo-600 text-indigo-600 dark:text-indigo-400'
          : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="activeTab = tab.key"
      >
        {{ t(tab.labelKey) }}
      </button>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <template v-else>
      <!-- Travel List -->
      <DataTable v-if="activeTab === 'travel'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="title" :header="t('business_travel.field_title')" sortable />
        <Column field="start_date" :header="t('business_travel.start_date')" sortable style="width:130px">
          <template #body="{data}">{{ formatDate(data.start_date, locale) }}</template>
        </Column>
        <Column field="end_date" :header="t('business_travel.end_date')" sortable style="width:130px">
          <template #body="{data}">{{ formatDate(data.end_date, locale) }}</template>
        </Column>
        <Column field="status" :header="t('common.status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="approval_status" :header="t('business_travel.approval_status')" sortable style="width:140px" />
        <Column field="destination_city" :header="t('business_travel.destination_city')" style="width:140px">
          <template #body="{data}">{{ data.destination_city || '-' }}</template>
        </Column>
        <Column field="destination_count" :header="t('business_travel.destination_count')" sortable style="width:100px" />
      </DataTable>

      <!-- Funding -->
      <DataTable v-else-if="activeTab === 'funding'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="funding_method" :header="t('business_travel.funding_method')" sortable />
        <Column field="amount" :header="t('business_travel.amount')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.amount) }}</template>
        </Column>
        <Column field="funding_date" :header="t('business_travel.funding_date')" sortable style="width:130px">
          <template #body="{data}">{{ data.funding_date ? formatDate(data.funding_date, locale) : '-' }}</template>
        </Column>
        <Column field="funded_by" :header="t('business_travel.funded_by')" style="width:140px">
          <template #body="{data}">{{ data.funded_by || '-' }}</template>
        </Column>
        <Column field="status" :header="t('common.status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.status" :severity="fundingSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>

      <!-- Advance vs Actual -->
      <DataTable v-else-if="activeTab === 'advance'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="total_advance" :header="t('business_travel.total_advance')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.total_advance) }}</template>
        </Column>
        <Column field="total_actual_expense" :header="t('business_travel.total_actual_expense_col')" sortable style="width:170px">
          <template #body="{data}">{{ formatCurrency(data.total_actual_expense) }}</template>
        </Column>
        <Column field="remaining" :header="t('business_travel.remaining')" sortable style="width:150px">
          <template #body="{data}"><span :class="data.remaining < 0 ? 'text-rose-600 dark:text-rose-400' : 'text-emerald-600 dark:text-emerald-400'">{{ formatCurrency(data.remaining) }}</span></template>
        </Column>
        <Column field="settlement_status" :header="t('business_travel.settlement_status')" sortable style="width:150px" />
      </DataTable>

      <!-- Reimbursement -->
      <DataTable v-else-if="activeTab === 'reimbursement'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="amount" :header="t('business_travel.amount')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.amount) }}</template>
        </Column>
        <Column field="approved_at" :header="t('business_travel.approved_at')" sortable style="width:150px">
          <template #body="{data}">{{ data.approved_at ? formatDate(data.approved_at, locale) : '-' }}</template>
        </Column>
        <Column field="paid_at" :header="t('business_travel.paid_at')" sortable style="width:150px">
          <template #body="{data}">{{ data.paid_at ? formatDate(data.paid_at, locale) : '-' }}</template>
        </Column>
        <Column field="payment_reference" :header="t('business_travel.payment_reference')" style="width:170px">
          <template #body="{data}">{{ data.payment_reference || '-' }}</template>
        </Column>
        <Column field="status" :header="t('common.status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.status" :severity="fundingSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>

      <!-- Refund -->
      <DataTable v-else-if="activeTab === 'refund'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="advance" :header="t('business_travel.advance')" sortable style="width:140px">
          <template #body="{data}">{{ formatCurrency(data.advance) }}</template>
        </Column>
        <Column field="actual" :header="t('business_travel.actual')" sortable style="width:140px">
          <template #body="{data}">{{ formatCurrency(data.actual) }}</template>
        </Column>
        <Column field="refund_amount" :header="t('business_travel.refund_amount')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.refund_amount) }}</template>
        </Column>
        <Column field="refund_date" :header="t('business_travel.refund_date')" sortable style="width:140px">
          <template #body="{data}">{{ data.refund_date ? formatDate(data.refund_date, locale) : '-' }}</template>
        </Column>
        <Column field="status" :header="t('common.status')" sortable style="width:130px">
          <template #body="{data}"><Tag :value="data.status" :severity="fundingSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
      </DataTable>

      <!-- Travel Cost -->
      <DataTable v-else-if="activeTab === 'travel-cost'" :value="currentItems" paginator :rows="20" :rowsPerPageOptions="[10, 20, 50, 100]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty><ReportEmpty /></template>
        <Column field="request_number" :header="t('business_travel.request_number')" sortable style="width:150px" />
        <Column field="title" :header="t('business_travel.field_title')" sortable />
        <Column field="total_company_paid" :header="t('business_travel.total_company_paid')" sortable style="width:160px">
          <template #body="{data}">{{ formatCurrency(data.total_company_paid) }}</template>
        </Column>
        <Column field="total_advance" :header="t('business_travel.total_advance')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.total_advance) }}</template>
        </Column>
        <Column field="total_reimbursement" :header="t('business_travel.total_reimbursement')" sortable style="width:150px">
          <template #body="{data}">{{ formatCurrency(data.total_reimbursement) }}</template>
        </Column>
        <Column field="total_actual_cost" :header="t('business_travel.total_actual_cost')" sortable style="width:160px">
          <template #body="{data}"><span class="font-semibold">{{ formatCurrency(data.total_actual_cost) }}</span></template>
        </Column>
      </DataTable>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, h } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { formatDate } from '@/utils/formatDate'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import DateInput from '@/components/DateInput.vue'

const { t, locale } = useI18n()
const toast = useToast()

function toDateOnly(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const today = new Date()
const fromDate = ref(toDateOnly(new Date(today.getFullYear(), today.getMonth(), 1)))
const toDate = ref(toDateOnly(today))
const statusFilter = ref(null)

const statusOptions = [
  { label: 'DRAFT', value: 'DRAFT' },
  { label: 'SUBMITTED', value: 'SUBMITTED' },
  { label: 'APPROVED', value: 'APPROVED' },
  { label: 'REJECTED', value: 'REJECTED' },
  { label: 'IN_PROGRESS', value: 'IN_PROGRESS' },
  { label: 'COMPLETED', value: 'COMPLETED' },
  { label: 'CLOSED', value: 'CLOSED' },
  { label: 'CANCELLED', value: 'CANCELLED' }
]

const reportTabs = [
  { key: 'travel', labelKey: 'business_travel.report_travel', endpoint: '/api/v1/tenant/attendance/business-travels/reports/travel' },
  { key: 'funding', labelKey: 'business_travel.report_funding', endpoint: '/api/v1/tenant/attendance/business-travels/reports/funding' },
  { key: 'advance', labelKey: 'business_travel.report_advance', endpoint: '/api/v1/tenant/attendance/business-travels/reports/advance' },
  { key: 'reimbursement', labelKey: 'business_travel.report_reimbursement', endpoint: '/api/v1/tenant/attendance/business-travels/reports/reimbursement' },
  { key: 'refund', labelKey: 'business_travel.report_refund', endpoint: '/api/v1/tenant/attendance/business-travels/reports/refund' },
  { key: 'travel-cost', labelKey: 'business_travel.report_travel_cost', endpoint: '/api/v1/tenant/attendance/business-travels/reports/travel-cost' }
]

const activeTab = ref('travel')
const loading = ref(false)
// cache of results per tab key
const dataByTab = ref({})

const currentItems = computed(() => dataByTab.value[activeTab.value] || [])

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' }
]

function formatCurrency(v) {
  if (v === null || v === undefined) return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'secondary'
    case 'SUBMITTED': return 'info'
    case 'IN_PROGRESS': return 'warn'
    case 'COMPLETED': return 'info'
    case 'CLOSED': return 'success'
    default: return 'secondary'
  }
}

function statusLabel(status) {
  const key = `business_travel.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function fundingSeverity(status) {
  switch (String(status).toUpperCase()) {
    case 'CONFIRMED':
    case 'PAID':
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'PENDING': return 'warn'
    default: return 'secondary'
  }
}

async function loadData() {
  if (!fromDate.value || !toDate.value) return
  loading.value = true
  try {
    const tab = reportTabs.find(x => x.key === activeTab.value)
    const params = { from: fromDate.value, to: toDate.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get(tab.endpoint, { params })
    dataByTab.value = { ...dataByTab.value, [tab.key]: res.data?.data || [] }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

// Simple inline "empty" component to avoid duplicating the empty-state markup per table
const ReportEmpty = {
  setup() {
    return () => h('div', { class: 'flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500' }, [
      h('i', { class: 'pi pi-chart-bar text-3xl mb-2 opacity-50' }),
      h('p', { class: 'text-sm font-medium' }, t('business_travel.reports_empty'))
    ])
  }
}

watch(activeTab, () => {
  if (!dataByTab.value[activeTab.value]) loadData()
})

onMounted(loadData)
</script>
