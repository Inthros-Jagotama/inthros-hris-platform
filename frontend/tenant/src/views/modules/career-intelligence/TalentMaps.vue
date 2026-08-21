<template>
  <div class="space-y-4">
    <!-- ── Toolbar: period + tombol ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <InputText
          v-model="period"
          :placeholder="t('talent_maps.period_placeholder')"
          class="!w-36 !text-sm"
          @keyup.enter="loadAll"
        />
        <Button icon="pi pi-refresh" size="small" severity="secondary" outlined v-tooltip.bottom="t('common.refresh')" @click="loadAll" />
      </div>
      <div class="flex items-center gap-2">
        <Button icon="pi pi-cog" size="small" severity="secondary" outlined v-tooltip.bottom="t('talent_maps.settings_title')" @click="openSettings" />
        <Button :label="t('talent_maps.generate_assessment')" icon="pi pi-bolt" size="small" @click="openGenerateDialog" />
      </div>
    </div>

    <!-- ── Grid 9-box ── -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
      <div
        v-for="q in quadrants"
        :key="q.position"
        class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 cursor-pointer transition-colors hover:border-violet-300 dark:hover:border-violet-500/60"
        :class="selectedPosition === q.position ? 'ring-2 ring-violet-400 dark:ring-violet-500' : ''"
        @click="selectedPosition = selectedPosition === q.position ? null : q.position"
      >
        <div class="flex items-center justify-between mb-1">
          <span class="text-xs font-semibold text-gray-500 dark:text-gray-400">{{ q.position }}</span>
          <span class="inline-flex items-center justify-center w-6 h-6 rounded-full bg-violet-100 dark:bg-violet-900/40 text-violet-600 dark:text-violet-300 text-xs font-bold">{{ q.count }}</span>
        </div>
        <p class="text-sm font-medium text-navy-800 dark:text-gray-100 leading-tight">{{ q.label }}</p>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-1 line-clamp-2">{{ q.description }}</p>
      </div>
    </div>

    <!-- ── Daftar penilaian ── -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs text-gray-400 dark:text-gray-500">
          {{ selectedPosition ? t('talent_maps.filtered_by', { position: selectedPosition }) : t('talent_maps.all_period', { period: period || '—' }) }}
        </span>
        <Button v-if="selectedPosition" :label="t('common.reset')" icon="pi pi-times" size="small" text severity="secondary" @click="selectedPosition = null" />
      </div>
      <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />
      <DataTable v-else :value="filteredItems" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
            <i class="pi pi-th-large text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('talent_maps.empty_maps') }}</p>
          </div>
        </template>
        <Column :header="t('talent_maps.employee')" style="width:200px">
          <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ employeeName(data.employee_id) }}</span></template>
        </Column>
        <Column field="performance" :header="t('talent_maps.performance')" style="width:110px">
          <template #body="{data}"><Tag :value="t('talent_maps.level_' + data.performance)" :severity="levelSeverity(data.performance)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="potential" :header="t('talent_maps.potential')" style="width:110px">
          <template #body="{data}"><Tag :value="t('talent_maps.level_' + data.potential)" :severity="levelSeverity(data.potential)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="grid_position" :header="t('talent_maps.grid_position')" style="width:100px">
          <template #body="{data}"><span class="text-xs text-gray-500 dark:text-gray-400 font-mono">{{ data.grid_position }}</span></template>
        </Column>
        <Column field="notes" :header="t('talent_maps.notes')">
          <template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs line-clamp-1">{{ data.notes || '—' }}</span></template>
        </Column>
        <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
          <template #body="{data}">
            <div class="flex items-center justify-end gap-1">
              <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openEditDialog(data)" />
              <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- ── Dialog: Generate (otomatis dari Performance + Competency) ── -->
    <Dialog
      v-model:visible="generateDialogVisible"
      :header="t('talent_maps.generate_assessment')"
      modal
      :style="{ width: '480px' }"
      @hide="resetGenerateForm"
    >
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('talent_maps.generate_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('talent_maps.employee')" required :errors="generateErrors?.employee_id">
          <Select
            v-model="generateForm.employee_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.period')" required :errors="generateErrors?.period">
          <InputText v-model="generateForm.period" :placeholder="t('talent_maps.period_placeholder')" class="w-full !text-sm" />
        </FormRow>
        <FormRow :label="t('talent_maps.notes')" :errors="generateErrors?.notes">
          <TextInput v-model="generateForm.notes" textarea :rows="2" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="generateDialogVisible = false" />
          <Button :label="t('talent_maps.generate')" icon="pi pi-bolt" size="small" :loading="generating" :disabled="generating" @click="handleGenerate" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Koreksi manual (performance/potential/notes) ── -->
    <Dialog
      v-model:visible="editDialogVisible"
      :header="t('talent_maps.edit_assessment')"
      modal
      :style="{ width: '440px' }"
      @hide="resetEditForm"
    >
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('talent_maps.edit_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('talent_maps.performance')" required :errors="editErrors?.performance">
          <Select
            v-model="editForm.performance"
            :options="levelOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.potential')" required :errors="editErrors?.potential">
          <Select
            v-model="editForm.potential"
            :options="levelOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.notes')" :errors="editErrors?.notes">
          <TextInput v-model="editForm.notes" textarea :rows="3" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="editDialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="editSaving" :disabled="editSaving" @click="handleEditSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Pengaturan ambang batas ── -->
    <Dialog
      v-model:visible="settingsVisible"
      :header="t('talent_maps.settings_title')"
      modal
      :style="{ width: '460px' }"
    >
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('talent_maps.settings_hint') }}</p>
      <div class="space-y-3">
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-400">{{ t('talent_maps.performance') }}</p>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('talent_maps.low_max')" :errors="settingsErrors?.performance_low_max">
            <InputNumber v-model="settingsForm.performance_low_max" :min="0" :max="100" suffix="%" class="!w-28" />
          </FormRow>
          <FormRow :label="t('talent_maps.high_min')" :errors="settingsErrors?.performance_high_min">
            <InputNumber v-model="settingsForm.performance_high_min" :min="0" :max="100" suffix="%" class="!w-28" />
          </FormRow>
        </div>
        <p class="text-xs font-semibold uppercase tracking-wide text-gray-400 mt-2">{{ t('talent_maps.potential') }}</p>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('talent_maps.low_max')" :errors="settingsErrors?.potential_low_max">
            <InputNumber v-model="settingsForm.potential_low_max" :min="0" :max="100" suffix="%" class="!w-28" />
          </FormRow>
          <FormRow :label="t('talent_maps.high_min')" :errors="settingsErrors?.potential_high_min">
            <InputNumber v-model="settingsForm.potential_high_min" :min="0" :max="100" suffix="%" class="!w-28" />
          </FormRow>
        </div>
        <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('talent_maps.settings_band_hint') }}</p>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="settingsVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="settingsSaving" :disabled="settingsSaving" @click="handleSettingsSave" />
        </div>
      </template>
    </Dialog>

    <!-- ── Konfirmasi hapus ── -->
    <ConfirmDeleteDialog
      v-model:visible="deleteConfirmVisible"
      :title="t('talent_maps.confirm_delete_title')"
      :message="t('talent_maps.confirm_delete_msg')"
      :loading="actionLoading"
      :error-msg="actionError"
      :cancel-label="t('common.no')"
      :confirm-label="t('common.delete')"
      @confirm="handleDeleteConfirm"
      @cancel="deleteConfirmVisible = false"
    />
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
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t } = useI18n()
const toast = useToast()

function currentQuarterPeriod() {
  const now = new Date()
  const q = Math.floor(now.getMonth() / 3) + 1
  return `${now.getFullYear()}-Q${q}`
}

// ── State ──
const period = ref(currentQuarterPeriod())
const loading = ref(false)
const items = ref([])
const quadrants = ref(emptyQuadrants())
const selectedPosition = ref(null)
const employees = ref([])

// ── Dialog: generate ──
const generateDialogVisible = ref(false)
const generating = ref(false)
const generateErrors = ref({})
const generateForm = ref(emptyGenerateForm())

// ── Dialog: edit (koreksi manual) ──
const editDialogVisible = ref(false)
const editSaving = ref(false)
const editingId = ref(null)
const editErrors = ref({})
const editForm = ref({ performance: null, potential: null, notes: '' })

// ── Dialog: settings ──
const settingsVisible = ref(false)
const settingsSaving = ref(false)
const settingsErrors = ref({})
const settingsForm = ref({ performance_low_max: 50, performance_high_min: 80, potential_low_max: 50, potential_high_min: 80 })

// ── Konfirmasi hapus ──
const deleteTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const deleteConfirmVisible = ref(false)

const employeeOptions = computed(() =>
  employees.value.map(e => ({ label: e.name, value: e.id }))
)
const levelOptions = computed(() => [
  { label: t('talent_maps.level_LOW'), value: 'LOW' },
  { label: t('talent_maps.level_MEDIUM'), value: 'MEDIUM' },
  { label: t('talent_maps.level_HIGH'), value: 'HIGH' }
])
const filteredItems = computed(() =>
  selectedPosition.value ? items.value.filter(i => i.grid_position === selectedPosition.value) : items.value
)

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

function emptyQuadrants() {
  return Array.from({ length: 9 }, (_, i) => ({ position: `9-BOX-${i + 1}`, label: '', description: '', count: 0 }))
}

function emptyGenerateForm() {
  return { employee_id: null, period: period.value, notes: '' }
}

function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || '—'
}
function levelSeverity(level) {
  if (level === 'HIGH') return 'success'
  if (level === 'MEDIUM') return 'info'
  return 'warning'
}

// ── Load data ──
async function loadGrid() {
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/talent-maps/grid', { params: { period: period.value } })
    quadrants.value = res.data?.data?.quadrants || emptyQuadrants()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

async function loadList() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/talent-maps', { params: { period: period.value, page: 1, per_page: 200 } })
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadAll() {
  selectedPosition.value = null
  await Promise.all([loadGrid(), loadList()])
}

async function loadReferences() {
  const empRes = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  employees.value = empRes.data?.data || []
}

// ── Generate (otomatis) ──
function openGenerateDialog() {
  generateErrors.value = {}
  generateForm.value = emptyGenerateForm()
  generateDialogVisible.value = true
}

function resetGenerateForm() {
  generateForm.value = emptyGenerateForm()
  generateErrors.value = {}
}

function validateGenerateForm() {
  const e = {}
  if (!generateForm.value.employee_id) e.employee_id = t('career_paths.field_required')
  if (!generateForm.value.period?.trim()) e.period = t('career_paths.field_required')
  return e
}

async function handleGenerate() {
  generateErrors.value = validateGenerateForm()
  if (Object.keys(generateErrors.value).length > 0) return
  generating.value = true
  try {
    const payload = {
      employee_id: generateForm.value.employee_id,
      period: generateForm.value.period.trim(),
      notes: generateForm.value.notes?.trim() || undefined
    }
    await api.post('/api/v1/tenant/career-intelligence/talent-maps/generate', payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    generateDialogVisible.value = false
    await loadAll()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      generateErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 6000 })
    }
  } finally {
    generating.value = false
  }
}

// ── Edit (koreksi manual) ──
function openEditDialog(data) {
  editErrors.value = {}
  editingId.value = data.id
  editForm.value = {
    performance: data.performance,
    potential: data.potential,
    notes: data.notes || ''
  }
  editDialogVisible.value = true
}

function resetEditForm() {
  editForm.value = { performance: null, potential: null, notes: '' }
  editErrors.value = {}
  editingId.value = null
}

function validateEditForm() {
  const e = {}
  if (!editForm.value.performance) e.performance = t('career_paths.field_required')
  if (!editForm.value.potential) e.potential = t('career_paths.field_required')
  return e
}

async function handleEditSave() {
  editErrors.value = validateEditForm()
  if (Object.keys(editErrors.value).length > 0) return
  editSaving.value = true
  try {
    const payload = {
      performance: editForm.value.performance,
      potential: editForm.value.potential,
      notes: editForm.value.notes?.trim() || ''
    }
    await api.put(`/api/v1/tenant/career-intelligence/talent-maps/${editingId.value}`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    editDialogVisible.value = false
    await loadAll()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      editErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    editSaving.value = false
  }
}

// ── Settings ──
async function openSettings() {
  settingsErrors.value = {}
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/talent-maps/settings')
    const s = res.data?.data
    if (s) settingsForm.value = { ...s }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
  settingsVisible.value = true
}

function validateSettingsForm() {
  const e = {}
  if (settingsForm.value.performance_low_max >= settingsForm.value.performance_high_min) {
    e.performance_high_min = t('talent_maps.low_lt_high_error')
  }
  if (settingsForm.value.potential_low_max >= settingsForm.value.potential_high_min) {
    e.potential_high_min = t('talent_maps.low_lt_high_error')
  }
  return e
}

async function handleSettingsSave() {
  settingsErrors.value = validateSettingsForm()
  if (Object.keys(settingsErrors.value).length > 0) return
  settingsSaving.value = true
  try {
    await api.put('/api/v1/tenant/career-intelligence/talent-maps/settings', settingsForm.value)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    settingsVisible.value = false
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      settingsErrors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    settingsSaving.value = false
  }
}

// ── Hapus ──
function openDeleteConfirm(data) {
  deleteTarget.value = data
  actionError.value = ''
  deleteConfirmVisible.value = true
}

function handleDeleteConfirm() {
  if (!deleteTarget.value?.id) return
  actionError.value = ''
  actionLoading.value = true
  api.delete(`/api/v1/tenant/career-intelligence/talent-maps/${deleteTarget.value.id}`)
    .then(async () => {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      deleteTarget.value = null
      await loadAll()
    })
    .catch(e => {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    })
    .finally(() => { actionLoading.value = false })
}

onMounted(() => {
  loadAll()
  loadReferences()
})
</script>
