<template>
  <div class="space-y-4">
    <!-- ── Toolbar ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <InputText
          v-model="period"
          :placeholder="t('risk_dashboard.period_placeholder')"
          class="!w-32 !text-sm"
          @keyup.enter="loadAll"
        />
        <Select
          v-model="levelFilter"
          :options="levelOptions"
          optionLabel="label"
          optionValue="value"
          showClear
          class="!w-40 !text-sm"
          :placeholder="t('risk_dashboard.filter_by_level')"
          @change="loadIndicators"
        />
        <Button icon="pi pi-refresh" size="small" severity="secondary" outlined v-tooltip.bottom="t('common.refresh')" @click="loadAll" />
      </div>
      <Button :label="t('risk_dashboard.add_indicator')" icon="pi pi-plus" size="small" @click="openDialog()" />
    </div>

    <!-- ── Ringkasan ── -->
    <div class="grid grid-cols-2 sm:grid-cols-3 gap-3">
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('risk_dashboard.total_risks') }}</p>
        <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ dashboard?.total_risks ?? '—' }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('risk_dashboard.high_risks') }}</p>
        <p class="text-lg font-bold text-amber-600 dark:text-amber-400">{{ dashboard?.high_risks ?? '—' }}</p>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <p class="text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('risk_dashboard.critical_risks') }}</p>
        <p class="text-lg font-bold text-rose-600 dark:text-rose-400">{{ dashboard?.critical_risks ?? '—' }}</p>
      </div>
    </div>

    <!-- ── Widget detail (data real) ── -->
    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
      <div v-for="widget in widgets" :key="widget.key" class="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm font-semibold text-navy-800 dark:text-gray-100">{{ t('risk_dashboard.widget_' + widget.key) }}</p>
          <Tag v-if="widget.data?.risk_level" :value="widget.data.risk_level" :severity="levelSeverity(widget.data.risk_level)" class="!text-xs !px-1.5 !py-0.5" />
        </div>
        <div v-if="widget.data" class="space-y-1">
          <p class="text-2xl font-bold" :class="widget.data.exceeded_by > 0 ? 'text-rose-600 dark:text-rose-400' : 'text-emerald-600 dark:text-emerald-400'">
            {{ widget.data.value?.toFixed(1) }}<span class="text-sm font-normal text-gray-400">{{ widget.suffix }}</span>
          </p>
          <p class="text-xs text-gray-400">{{ t('risk_dashboard.threshold') }}: {{ widget.data.threshold }}{{ widget.suffix }}</p>
          <div v-if="widget.data.affected_departments?.length" class="flex flex-wrap gap-1 mt-2">
            <span v-for="d in widget.data.affected_departments" :key="d.label" class="inline-flex items-center px-2 py-0.5 rounded-md text-xs border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 text-gray-600 dark:text-gray-300">
              {{ d.label }}: {{ d.value.toFixed(1) }}
            </span>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('risk_dashboard.no_department_data') }}</p>
        </div>
        <SkeletonCard v-else type="stat" :count="1" />
      </div>
    </div>

    <!-- ── Daftar Indikator ── -->
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
          <i class="pi pi-exclamation-triangle text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('risk_dashboard.empty_indicators') }}</p>
        </div>
      </template>
      <Column field="risk_code" :header="t('risk_dashboard.risk_code')" style="width:160px" />
      <Column field="risk_name" :header="t('risk_dashboard.risk_name')" />
      <Column :header="t('risk_dashboard.risk_level')" style="width:110px">
        <template #body="{data}"><Tag :value="data.risk_level" :severity="levelSeverity(data.risk_level)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="score" :header="t('risk_dashboard.score')" style="width:90px" />
      <Column field="threshold" :header="t('risk_dashboard.threshold')" style="width:90px" />
      <Column field="snapshot_at" :header="t('risk_dashboard.snapshot_at')" style="width:120px" />
      <Column :header="t('common.actions')" style="width:60px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Tambah / Ubah ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingId ? t('risk_dashboard.edit_indicator') : t('risk_dashboard.add_indicator')"
      modal
      :style="{ width: '480px' }"
      @hide="resetForm"
    >
      <div class="space-y-3">
        <FormRow :label="t('risk_dashboard.period')" required :errors="errors?.period">
          <InputText v-model="form.period" :disabled="!!editingId" :placeholder="t('risk_dashboard.period_placeholder')" class="w-full !text-sm" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.risk_code')" required :errors="errors?.risk_code">
          <InputText v-model="form.risk_code" :disabled="!!editingId" class="w-full !text-sm" placeholder="e.g. HIGH_TURNOVER" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.risk_name')" required :errors="errors?.risk_name">
          <InputText v-model="form.risk_name" :disabled="!!editingId" class="w-full !text-sm" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.risk_level')" required :errors="errors?.risk_level">
          <Select v-model="form.risk_level" :options="levelOptions" optionLabel="label" optionValue="value" class="w-full !text-sm" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.score')" :errors="errors?.score">
          <InputNumber v-model="form.score" :disabled="!!editingId" :minFractionDigits="0" :maxFractionDigits="2" inputClass="!text-sm" class="w-full" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.threshold')" :errors="errors?.threshold">
          <InputNumber v-model="form.threshold" :disabled="!!editingId" :minFractionDigits="0" :maxFractionDigits="2" inputClass="!text-sm" class="w-full" />
        </FormRow>
        <FormRow :label="t('risk_dashboard.department')" :errors="errors?.department_id">
          <Select
            v-model="form.department_id"
            :options="organizationOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            :disabled="!!editingId"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('risk_dashboard.recommendation')" :errors="errors?.recommendation">
          <TextInput v-model="form.recommendation" textarea :rows="2" />
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
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

function currentPeriod() {
  const now = new Date()
  return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`
}

// ── State ──
const period = ref(currentPeriod())
const levelFilter = ref(null)
const loading = ref(false)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const dashboard = ref(null)
const organizations = ref([])

const widgets = ref([
  { key: 'high-turnover', suffix: '%', data: null },
  { key: 'retirement', suffix: '%', data: null },
  { key: 'contract-expiry', suffix: '', data: null },
  { key: 'high-absenteeism', suffix: '%', data: null }
])

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const organizationOptions = computed(() =>
  organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id }))
)
const levelOptions = computed(() => [
  { label: t('risk_dashboard.level_LOW'), value: 'LOW' },
  { label: t('risk_dashboard.level_MEDIUM'), value: 'MEDIUM' },
  { label: t('risk_dashboard.level_HIGH'), value: 'HIGH' },
  { label: t('risk_dashboard.level_CRITICAL'), value: 'CRITICAL' }
])

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-48', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

function emptyForm() {
  return { period: period.value, risk_code: '', risk_name: '', risk_level: null, score: null, threshold: null, department_id: null, recommendation: '' }
}

function levelSeverity(level) {
  if (level === 'CRITICAL') return 'danger'
  if (level === 'HIGH') return 'warn'
  if (level === 'MEDIUM') return 'info'
  return 'success'
}

// ── Load data ──
async function loadDashboard() {
  try {
    const res = await api.get('/api/v1/tenant/workforce-intelligence/risk/dashboard', { params: { period: period.value } })
    dashboard.value = res.data?.data || null
  } catch (e) {
    dashboard.value = null
  }
}

async function loadIndicators() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value, period: period.value }
    if (levelFilter.value) params.risk_level = levelFilter.value
    const res = await api.get('/api/v1/tenant/workforce-intelligence/risk/indicators', { params })
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

async function loadWidgets() {
  await Promise.all(widgets.value.map(async w => {
    try {
      const res = await api.get(`/api/v1/tenant/workforce-intelligence/risk/${w.key}`)
      w.data = res.data?.data || null
    } catch (e) {
      w.data = null
    }
  }))
}

async function loadAll() {
  currentPage.value = 1
  await Promise.all([loadDashboard(), loadIndicators(), loadWidgets()])
}

async function loadReferences() {
  const orgRes = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
  organizations.value = orgRes.data?.data || []
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadIndicators()
}

// ── Dialog create/edit ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      period: data.period,
      risk_code: data.risk_code,
      risk_name: data.risk_name,
      risk_level: data.risk_level,
      score: data.score,
      threshold: data.threshold,
      department_id: data.department_id || null,
      recommendation: data.recommendation || ''
    }
  } else {
    form.value = emptyForm()
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = emptyForm()
  errors.value = {}
  editingId.value = null
}

function validateForm() {
  const e = {}
  if (!editingId.value) {
    if (!form.value.period?.trim()) e.period = t('career_paths.field_required')
    if (!form.value.risk_code?.trim()) e.risk_code = t('career_paths.field_required')
    if (!form.value.risk_name?.trim()) e.risk_name = t('career_paths.field_required')
  }
  if (!form.value.risk_level) e.risk_level = t('career_paths.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    if (editingId.value) {
      const payload = {
        risk_level: form.value.risk_level,
        recommendation: form.value.recommendation?.trim() || ''
      }
      await api.put(`/api/v1/tenant/workforce-intelligence/risk/indicators/${editingId.value}`, payload)
    } else {
      const payload = {
        period: form.value.period.trim(),
        risk_code: form.value.risk_code.trim(),
        risk_name: form.value.risk_name.trim(),
        risk_level: form.value.risk_level,
        score: form.value.score ?? 0,
        threshold: form.value.threshold ?? 0,
        department_id: form.value.department_id || undefined,
        recommendation: form.value.recommendation?.trim() || undefined
      }
      await api.post('/api/v1/tenant/workforce-intelligence/risk/indicators', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadAll()
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

onMounted(() => {
  loadAll()
  loadReferences()
})
</script>
