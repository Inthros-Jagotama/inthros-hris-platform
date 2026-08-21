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
      <Button :label="t('talent_maps.add_assessment')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
              <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
              <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- ── Dialog: Tambah / Ubah ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingId ? t('talent_maps.edit_assessment') : t('talent_maps.add_assessment')"
      modal
      :style="{ width: '480px' }"
      @hide="resetForm"
    >
      <div class="space-y-3">
        <FormRow :label="t('talent_maps.employee')" required :errors="errors?.employee_id">
          <Select
            v-model="form.employee_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            filter
            showClear
            :disabled="!!editingId"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.period')" required :errors="errors?.period">
          <InputText v-model="form.period" :disabled="!!editingId" :placeholder="t('talent_maps.period_placeholder')" class="w-full !text-sm" />
        </FormRow>
        <FormRow :label="t('talent_maps.performance')" required :errors="errors?.performance">
          <Select
            v-model="form.performance"
            :options="levelOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.potential')" required :errors="errors?.potential">
          <Select
            v-model="form.potential"
            :options="levelOptions"
            optionLabel="label"
            optionValue="value"
            class="w-full !text-sm"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('talent_maps.notes')" :errors="errors?.notes">
          <TextInput v-model="form.notes" textarea :rows="3" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
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

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())

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

function emptyForm() {
  return { employee_id: null, period: period.value, performance: null, potential: null, notes: '' }
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

// ── Dialog create/edit ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      employee_id: data.employee_id,
      period: data.period,
      performance: data.performance,
      potential: data.potential,
      notes: data.notes || ''
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

// ── Simpan ──
function validateForm() {
  const e = {}
  if (!editingId.value && !form.value.employee_id) e.employee_id = t('career_paths.field_required')
  if (!editingId.value && !form.value.period?.trim()) e.period = t('career_paths.field_required')
  if (!form.value.performance) e.performance = t('career_paths.field_required')
  if (!form.value.potential) e.potential = t('career_paths.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    if (editingId.value) {
      const payload = {
        performance: form.value.performance,
        potential: form.value.potential,
        notes: form.value.notes?.trim() || ''
      }
      await api.put(`/api/v1/tenant/career-intelligence/talent-maps/${editingId.value}`, payload)
    } else {
      const payload = {
        employee_id: form.value.employee_id,
        period: form.value.period.trim(),
        performance: form.value.performance,
        potential: form.value.potential,
        notes: form.value.notes?.trim() || undefined
      }
      await api.post('/api/v1/tenant/career-intelligence/talent-maps', payload)
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
