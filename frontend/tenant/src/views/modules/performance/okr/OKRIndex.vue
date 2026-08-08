<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-3">
        <Select v-model="filterPeriod" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('okr.filter_period')" class="w-48" size="small" showClear />
        <Select v-model="filterStatus" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('okr.filter_status')" class="w-36" size="small" showClear />
      </div>
    </div>

    <!-- DataTable -->
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
          <i class="pi pi-bullseye text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('okr.empty_evaluations') }}</p>
          <p class="text-xs mt-1">{{ t('okr.empty_evaluations_hint') }}</p>
        </div>
      </template>

      <Column field="employee_name" :header="t('okr.employee')" sortable>
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name || data.employee_id }}</span>
        </template>
      </Column>

      <Column field="organization_name" :header="t('okr.organization')" sortable style="width:180px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.organization_name || '-' }}</span>
        </template>
      </Column>

      <Column field="period_code" :header="t('okr.period')" sortable style="width:100px">
        <template #body="{data}">
          <Tag v-if="data.period_code" :value="data.period_code" severity="info" class="!text-xs" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="final_score" :header="t('okr.score')" sortable style="width:80px">
        <template #body="{data}">
          <span v-if="data.final_score > 0" class="font-bold" :class="getScoreClass(data.final_score)">{{ data.final_score.toFixed(1) }}</span>
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="rating_name" :header="t('okr.rating')" sortable style="width:100px">
        <template #body="{data}">
          <Tag v-if="data.rating_name" :value="data.rating_name" :severity="getRatingSeverity(data.rating_color)" class="!text-xs" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="status" :header="t('okr.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="data.status" :severity="getStatusSeverity(data.status)" class="!text-xs" />
        </template>
      </Column>

      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.view')" @click="viewEvaluation(data)" />
            <Button v-if="data.status === 'DRAFT'" icon="pi pi-pencil" size="small" text severity="info" v-tooltip.left="t('common.edit')" @click="editEvaluation(data)" />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'

const router = useRouter()
const toast = useToast()
const { t } = useI18n()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const filterPeriod = ref(null)
const filterStatus = ref(null)

const periodOptions = ref([])

const statusOptions = [
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Key Result Submitted', value: 'KR_SUBMITTED' },
  { label: 'Key Result Approved', value: 'KR_APPROVED' },
  { label: 'Submitted', value: 'SUBMITTED' },
  { label: 'Approved', value: 'APPROVED' },
  { label: 'Completed', value: 'COMPLETED' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function getScoreClass(score) {
  if (score >= 85) return 'text-emerald-600 dark:text-emerald-400'
  if (score >= 70) return 'text-blue-600 dark:text-blue-400'
  if (score >= 60) return 'text-amber-600 dark:text-amber-400'
  return 'text-red-600 dark:text-red-400'
}

function getRatingSeverity(color) {
  switch (color) {
    case 'success': return 'success'
    case 'primary': return 'info'
    case 'warning': return 'warn'
    case 'danger': return 'danger'
    default: return 'secondary'
  }
}

function getStatusSeverity(status) {
  switch (status) {
    case 'COMPLETED': return 'success'
    case 'APPROVED': return 'info'
    case 'SUBMITTED': return 'warn'
    case 'KR_APPROVED': return 'info'
    case 'KR_SUBMITTED': return 'warn'
    default: return 'secondary'
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: perPage.value
    }
    if (filterPeriod.value) params.period_id = filterPeriod.value
    if (filterStatus.value) params.status = filterStatus.value

    const res = await api.get('/api/v1/tenant/performance/okr/evaluations', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.meta?.total ?? body?.total ?? 0
    if (body?.meta?.page) currentPage.value = body.meta.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadReferenceData() {
  try {
    const periodRes = await api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } })
    periodOptions.value = (periodRes.data?.data || []).map(p => ({
      label: `${p.period_code} (${p.year})`,
      value: p.id
    }))
  } catch {
    // Silently fail
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

watch([filterPeriod, filterStatus], () => {
  currentPage.value = 1
  loadData()
})

function viewEvaluation(item) {
  router.push(`/performance/okr/evaluation/${item.id}`)
}

function editEvaluation(item) {
  router.push(`/performance/okr/evaluation/${item.id}`)
}

onMounted(async () => {
  await loadReferenceData()
  await loadData()
})
</script>