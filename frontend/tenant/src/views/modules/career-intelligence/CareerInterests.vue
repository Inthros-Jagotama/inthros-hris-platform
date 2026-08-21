<template>
  <div class="space-y-1">
    <!-- ── Toolbar ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Select
          v-model="employeeFilter"
          :options="employeeOptions"
          optionLabel="label"
          optionValue="value"
          filter
          showClear
          class="!w-56 !text-sm"
          :placeholder="t('career_interests.filter_by_employee')"
          @change="onFilterChange"
        />
        <Button :label="t('career_interests.add_interest')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-compass text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('career_interests.empty_interests') }}</p>
        </div>
      </template>
      <Column :header="t('career_interests.employee')" style="width:180px">
        <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="interest_type" :header="t('career_interests.interest_type')" style="width:140px">
        <template #body="{data}">
          <Tag :value="t('career_interests.type_' + data.interest_type)" severity="info" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('career_interests.target')" style="width:200px">
        <template #body="{data}">
          <div class="flex flex-col">
            <span v-if="data.target_position" class="text-gray-700 dark:text-gray-200">{{ data.target_position }}</span>
            <span v-if="data.target_department" class="text-xs text-gray-400">{{ data.target_department }}</span>
            <span v-if="!data.target_position && !data.target_department" class="text-gray-300 dark:text-gray-600">—</span>
          </div>
        </template>
      </Column>
      <Column field="readiness_level" :header="t('career_interests.readiness_level')" style="width:130px">
        <template #body="{data}">
          <Tag v-if="data.readiness_level" :value="t('career_interests.readiness_' + data.readiness_level)" :severity="readinessSeverity(data.readiness_level)" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-gray-300 dark:text-gray-600">—</span>
        </template>
      </Column>
      <Column field="recorded_at" :header="t('career_interests.recorded_at')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.recorded_at ? formatDate(data.recorded_at, locale) : '—' }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:60px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button icon="pi pi-eye" size="small" severity="info" text v-tooltip.left="t('common.view')" @click="openDetail(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Detail ── -->
    <Dialog v-model:visible="detailVisible" :header="t('career_interests.interest_detail_title')" modal :style="{ width: '480px' }">
      <div v-if="detailItem" class="space-y-3">
        <div>
          <p class="text-base font-semibold text-navy-800 dark:text-gray-100">{{ employeeName(detailItem.employee_id) }}</p>
          <Tag :value="t('career_interests.type_' + detailItem.interest_type)" severity="info" class="mt-1 !text-xs !px-1.5 !py-0.5" />
        </div>
        <div v-if="detailItem.target_position || detailItem.target_department" class="text-sm">
          <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-1">{{ t('career_interests.target') }}</p>
          <p class="text-gray-700 dark:text-gray-200">{{ detailItem.target_position || '—' }}</p>
          <p v-if="detailItem.target_department" class="text-xs text-gray-400">{{ detailItem.target_department }}</p>
        </div>
        <div v-if="detailItem.readiness_level" class="text-sm">
          <span class="text-gray-500 dark:text-gray-400">{{ t('career_interests.readiness_level') }}: </span>
          <Tag :value="t('career_interests.readiness_' + detailItem.readiness_level)" :severity="readinessSeverity(detailItem.readiness_level)" class="!text-xs !px-1.5 !py-0.5" />
        </div>
        <div v-if="detailItem.motivation">
          <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-1">{{ t('career_interests.motivation') }}</p>
          <p class="text-sm text-gray-700 dark:text-gray-200 whitespace-pre-line">{{ detailItem.motivation }}</p>
        </div>
        <div v-if="detailItem.recorded_at" class="text-xs text-gray-400">
          {{ t('career_interests.recorded_at') }}: {{ formatDate(detailItem.recorded_at, locale) }}
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="detailVisible = false" />
      </template>
    </Dialog>

    <!-- ── Dialog: Tambah ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="t('career_interests.add_interest')"
      modal
      :style="{ width: '520px' }"
      @hide="resetForm"
    >
      <div class="space-y-3">
        <FormRow :label="t('career_interests.employee')" required :errors="errors?.employee_id">
          <Select
            v-model="form.employee_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('career_interests.interest_type')" required :errors="errors?.interest_type">
          <Select
            v-model="form.interest_type"
            :options="interestTypeOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('career_interests.target_position')" :errors="errors?.target_position">
          <TextInput v-model="form.target_position" />
        </FormRow>
        <FormRow :label="t('career_interests.target_department')" :errors="errors?.target_department">
          <TextInput v-model="form.target_department" />
        </FormRow>
        <FormRow :label="t('career_interests.readiness_level')" :errors="errors?.readiness_level">
          <Select
            v-model="form.readiness_level"
            :options="readinessOptions"
            optionLabel="label"
            optionValue="value"
            showClear
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('career_interests.motivation')" :errors="errors?.motivation">
          <TextInput v-model="form.motivation" textarea :rows="3" />
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
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t, locale } = useI18n()
const toast = useToast()

// ── Daftar ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const employeeFilter = ref(null)

// ── Referensi ──
const employees = ref([])

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref(emptyForm())

// ── Detail ──
const detailVisible = ref(false)
const detailItem = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const employeeOptions = computed(() =>
  employees.value.map(e => ({ label: e.name, value: e.id }))
)
const interestTypeOptions = computed(() => [
  { label: t('career_interests.type_LEADERSHIP'), value: 'LEADERSHIP' },
  { label: t('career_interests.type_SPECIALIST'), value: 'SPECIALIST' },
  { label: t('career_interests.type_INTERNATIONAL'), value: 'INTERNATIONAL' },
  { label: t('career_interests.type_ENTREPRENEUR'), value: 'ENTREPRENEUR' }
])
const readinessOptions = computed(() => [
  { label: t('career_interests.readiness_NOW'), value: 'NOW' },
  { label: t('career_interests.readiness_1_YEAR'), value: '1_YEAR' },
  { label: t('career_interests.readiness_2_3_YEARS'), value: '2_3_YEARS' },
  { label: t('career_interests.readiness_3_PLUS'), value: '3_PLUS' }
])

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

function emptyForm() {
  return {
    employee_id: null,
    interest_type: null,
    target_position: '',
    target_department: '',
    readiness_level: null,
    motivation: ''
  }
}

function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || '—'
}
function readinessSeverity(level) {
  if (level === 'NOW') return 'success'
  if (level === '1_YEAR') return 'info'
  if (level === '2_3_YEARS') return 'warning'
  return 'secondary'
}

// ── Load data ──
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (employeeFilter.value) params.employee_id = employeeFilter.value
    const res = await api.get('/api/v1/tenant/career-intelligence/interests', { params })
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

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

async function loadReferences() {
  const empRes = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  employees.value = empRes.data?.data || []
}

// ── Dialog create ──
function openDialog() {
  errors.value = {}
  form.value = emptyForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = emptyForm()
  errors.value = {}
}

// ── Detail ──
function openDetail(data) {
  detailItem.value = data
  detailVisible.value = true
}

// ── Simpan ──
function validateForm() {
  const e = {}
  if (!form.value.employee_id) e.employee_id = t('career_paths.field_required')
  if (!form.value.interest_type) e.interest_type = t('career_paths.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    const payload = {
      employee_id: form.value.employee_id,
      interest_type: form.value.interest_type,
      target_position: form.value.target_position?.trim() || undefined,
      target_department: form.value.target_department?.trim() || undefined,
      motivation: form.value.motivation?.trim() || undefined,
      readiness_level: form.value.readiness_level || undefined
    }
    await api.post('/api/v1/tenant/career-intelligence/interests', payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    currentPage.value = 1
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

onMounted(() => {
  loadData()
  loadReferences()
})
</script>
