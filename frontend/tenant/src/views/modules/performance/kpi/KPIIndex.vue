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
        <Select v-model="filterPeriod" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('kpi.filter_period')" class="w-48" size="small" showClear />
        <Select v-model="filterStatus" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('kpi.filter_status')" class="w-36" size="small" showClear />
        <Button :label="t('kpi.new_evaluation')" icon="pi pi-plus" size="small" @click="createEvaluation" />
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
          <i class="pi pi-chart-bar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('kpi.empty_evaluations') }}</p>
          <p class="text-xs mt-1">{{ t('kpi.empty_evaluations_hint') }}</p>
        </div>
      </template>

      <Column field="employee_name" :header="t('kpi.employee')" sortable>
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name || data.employee_id }}</span>
        </template>
      </Column>

      <Column field="organization_name" :header="t('kpi.organization')" sortable style="width:180px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.organization_name || '-' }}</span>
        </template>
      </Column>

      <Column field="period_code" :header="t('kpi.period')" sortable style="width:100px">
        <template #body="{data}">
          <Tag v-if="data.period_code" :value="data.period_code" severity="info" class="!text-xs" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="final_score" :header="t('kpi.score')" sortable style="width:80px">
        <template #body="{data}">
          <span v-if="data.final_score > 0" class="font-bold" :class="getScoreClass(data.final_score)">{{ data.final_score.toFixed(1) }}</span>
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="rating_name" :header="t('kpi.rating')" sortable style="width:100px">
        <template #body="{data}">
          <Tag v-if="data.rating_name" :value="data.rating_name" :severity="getRatingSeverity(data.rating_color)" class="!text-xs" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column field="status" :header="t('kpi.status')" sortable style="width:100px">
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

    <!-- Create Evaluation Dialog -->
    <Dialog v-model:visible="createDialogVisible" :header="t('kpi.new_evaluation')" modal :style="{ width: '500px' }">
      <div class="space-y-4">
        <FormRow :label="t('kpi.employee')" required :errors="createErrors?.employee_id">
          <Select v-model="createForm.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" filter />
        </FormRow>
        <FormRow :label="t('kpi.template')" required :errors="createErrors?.template_id">
          <Select v-model="createForm.template_id" :options="templateOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
        </FormRow>
        <FormRow :label="t('kpi.period')" required :errors="createErrors?.period_id">
          <Select v-model="createForm.period_id" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="createDialogVisible = false" />
          <Button :label="t('common.create')" size="small" :loading="creating" @click="handleCreate" />
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
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'

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

const createDialogVisible = ref(false)
const creating = ref(false)
const createErrors = ref({})
const createForm = ref({
  employee_id: null,
  template_id: null,
  period_id: null
})

const periodOptions = ref([])
const templateOptions = ref([])
const employeeOptions = ref([])

const statusOptions = [
  { label: 'Draft', value: 'DRAFT' },
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

    const res = await api.get('/api/v1/tenant/performance/kpi/evaluations', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadReferenceData() {
  try {
    const [periodRes, templateRes, empRes] = await Promise.all([
      api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } }),
      api.get('/api/v1/tenant/performance/kpi/templates', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 200, status: 'active' } })
    ])

    periodOptions.value = (periodRes.data?.data || []).map(p => ({
      label: `${p.period_code} (${p.year})`,
      value: p.id
    }))

    templateOptions.value = (templateRes.data?.data || []).map(t => ({
      label: t.name,
      value: t.id
    }))

    employeeOptions.value = (empRes.data?.data || []).map(e => ({
      label: e.full_name || `${e.first_name} ${e.last_name}`,
      value: e.id
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

function createEvaluation() {
  createForm.value = { employee_id: null, template_id: null, period_id: null }
  createErrors.value = {}
  createDialogVisible.value = true
}

async function handleCreate() {
  createErrors.value = {}
  if (!createForm.value.employee_id) {
    createErrors.value = { employee_id: [t('form.required')] }
    return
  }
  if (!createForm.value.template_id) {
    createErrors.value = { template_id: [t('form.required')] }
    return
  }
  if (!createForm.value.period_id) {
    createErrors.value = { period_id: [t('form.required')] }
    return
  }

  creating.value = true
  try {
    const res = await api.post('/api/v1/tenant/performance/kpi/evaluations/snapshot', {
      employee_id: createForm.value.employee_id,
      template_id: createForm.value.template_id,
      period_id: createForm.value.period_id
    })

    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.evaluation_created'), life: 3000 })
    createDialogVisible.value = false

    const evalId = res.data?.data?.id || res.data?.id
    if (evalId) {
      router.push(`/performance/kpi/evaluation/${evalId}`)
    } else {
      await loadData()
    }
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      createErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    creating.value = false
  }
}

function viewEvaluation(item) {
  router.push(`/performance/kpi/evaluation/${item.id}`)
}

function editEvaluation(item) {
  router.push(`/performance/kpi/evaluation/${item.id}`)
}

onMounted(async () => {
  await loadReferenceData()
  await loadData()
})
</script>
