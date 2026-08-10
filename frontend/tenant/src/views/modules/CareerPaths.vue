<template>
  <div class="space-y-1">
    <!-- ── Toolbar: search + tombol ── -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <InputText
          v-model="searchTerm"
          :placeholder="t('career_paths.search_paths_placeholder')"
          class="!w-64"
          size="small"
          @keyup.enter="onFilterChange"
        />
        <Button
          v-if="searchTerm"
          icon="pi pi-times"
          severity="secondary"
          outlined
          size="small"
          v-tooltip.left="t('common.reset')"
          @click="resetFilters"
        />
        <Button :label="t('career_paths.add_path')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-sitemap text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('career_paths.empty_paths') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('career_paths.path_name')" style="width:220px">
        <template #body="{data}">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-200 font-medium">{{ data.name }}</span>
            <span v-if="data.description" class="text-xs text-gray-400 line-clamp-1 max-w-xs">{{ data.description }}</span>
          </div>
        </template>
      </Column>
      <Column :header="t('career_paths.path_steps')" style="width:280px">
        <template #body="{data}">
          <div class="flex items-center gap-1 flex-wrap">
            <template v-for="(step, idx) in data.steps || []" :key="step.id || idx">
              <span
                class="inline-flex items-center px-2 py-0.5 rounded-md text-xs border border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 text-gray-600 dark:text-gray-300"
              >
                {{ step.position_name || '—' }}
              </span>
              <i v-if="idx < (data.steps || []).length - 1" class="pi pi-arrow-right text-[10px] text-gray-400"></i>
            </template>
            <span v-if="!(data.steps || []).length" class="text-xs text-gray-400">—</span>
          </div>
        </template>
      </Column>
      <Column :header="t('career_paths.step_count')" style="width:90px">
        <template #body="{data}">
          <Tag :value="(data.steps || []).length" severity="info" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="is_active" :header="t('common.status')" style="width:110px">
        <template #body="{data}">
          <Tag
            :value="data.is_active ? t('career_paths.status_active') : t('career_paths.status_inactive')"
            :severity="data.is_active ? 'success' : 'secondary'"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>
      <Column field="updated_at" :header="t('career_paths.updated_at')" style="width:130px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.updated_at, locale) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:140px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center justify-end gap-1">
            <Button icon="pi pi-eye" size="small" severity="info" text v-tooltip.left="t('common.view')" @click="openDetail(data)" />
            <Button icon="pi pi-pencil" size="small" severity="secondary" text v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" severity="danger" text v-tooltip.left="t('common.delete')" @click="openDeleteConfirm(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Detail / Timeline ── -->
    <Dialog v-model:visible="detailVisible" :header="t('career_paths.path_detail_title')" modal :style="{ width: '520px' }">
      <div v-if="detailItem" class="space-y-4">
        <div>
          <p class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ detailItem.name }}</p>
          <p v-if="detailItem.description" class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ detailItem.description }}</p>
          <Tag
            class="mt-2 !text-xs !px-1.5 !py-0.5"
            :value="detailItem.is_active ? t('career_paths.status_active') : t('career_paths.status_inactive')"
            :severity="detailItem.is_active ? 'success' : 'secondary'"
          />
        </div>

        <!-- Timeline vertikal -->
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <p class="text-xs uppercase tracking-wide text-gray-400 font-medium mb-3">{{ t('career_paths.career_ladder') }}</p>
          <ol class="relative space-y-0">
            <template v-for="(step, idx) in sortedSteps(detailItem.steps)" :key="step.id || idx">
              <li class="relative flex gap-3 pb-4 last:pb-0">
                <!-- connector line -->
                <div class="flex flex-col items-center">
                  <span
                    class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-semibold border-2 shrink-0"
                    :class="idx === 0 ? 'bg-emerald-100 dark:bg-emerald-900/40 border-emerald-400 text-emerald-600 dark:text-emerald-300' : 'bg-indigo-100 dark:bg-indigo-900/40 border-indigo-400 text-indigo-600 dark:text-indigo-300'"
                  >
                    {{ idx + 1 }}
                  </span>
                  <span v-if="idx < (detailItem.steps || []).length - 1" class="w-px flex-1 bg-gray-300 dark:bg-gray-600 my-1"></span>
                </div>
                <div class="flex-1 pt-1 pb-1">
                  <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ step.position_name || step.position_id }}</p>
                  <div class="flex items-center gap-2 mt-1 flex-wrap">
                    <Tag v-if="step.minimum_service_months != null" severity="warning" class="!text-[10px] !px-1 !py-0.5">
                      {{ step.minimum_service_months }} {{ t('career_paths.months_short') }}
                    </Tag>
                    <span v-if="step.requirements" class="text-xs text-gray-500 dark:text-gray-400">
                      <i class="pi pi-info-circle mr-1"></i>{{ step.requirements }}
                    </span>
                  </div>
                </div>
              </li>
            </template>
            <li v-if="!(detailItem.steps || []).length" class="text-sm text-gray-400">{{ t('career_paths.no_steps') }}</li>
          </ol>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="detailVisible = false" />
          <Button :label="t('common.edit')" icon="pi pi-pencil" size="small" @click="openDialog(detailItem); detailVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: Tambah / Ubah Career Path ── -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="editingId ? t('career_paths.edit_path') : t('career_paths.add_path')"
      modal
      :style="{ width: '680px' }"
      @hide="resetForm"
    >
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('career_paths.path_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('career_paths.path_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" />
        </FormRow>
        <FormRow :label="t('career_paths.path_description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>

        <!-- Editor steps -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
              {{ t('career_paths.path_steps') }}
              <span class="text-red-400">*</span>
            </label>
            <Button icon="pi pi-plus" :label="t('career_paths.add_step')" severity="secondary" outlined size="small" @click="addStep" />
          </div>
          <p v-if="errors?.steps" class="text-xs text-red-500 mb-1">{{ errors.steps }}</p>
          <div class="space-y-2 max-h-72 overflow-y-auto pr-1">
            <div
              v-for="(step, idx) in form.steps"
              :key="idx"
              class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-2"
            >
              <div class="flex items-center justify-between">
                <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-indigo-500">
                  <span class="w-5 h-5 rounded-full bg-indigo-100 dark:bg-indigo-900/60 flex items-center justify-center">{{ idx + 1 }}</span>
                  {{ t('career_paths.step') }} {{ idx + 1 }}
                </span>
                <div class="flex items-center gap-1">
                  <Button icon="pi pi-arrow-up" size="small" text severity="secondary" :disabled="idx === 0" @click="moveStep(idx, -1)" />
                  <Button icon="pi pi-arrow-down" size="small" text severity="secondary" :disabled="idx === form.steps.length - 1" @click="moveStep(idx, 1)" />
                  <Button icon="pi pi-trash" size="small" text severity="danger" :disabled="form.steps.length === 1" @click="removeStep(idx)" />
                </div>
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
                <div class="sm:col-span-1">
                  <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">{{ t('career_paths.step_position') }}</label>
                  <Select
                    v-model="step.position_id"
                    :options="positionOptions"
                    optionLabel="label"
                    optionValue="value"
                    filter
                    showClear
                    class="w-full !text-sm"
                    :placeholder="t('career_paths.select_position')"
                  />
                </div>
                <div>
                  <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">{{ t('career_paths.min_service_months') }}</label>
                  <InputNumber v-model="step.minimum_service_months" :min="0" :useGrouping="false" class="w-full" />
                </div>
                <div>
                  <label class="block text-xs text-gray-500 dark:text-gray-400 mb-1">{{ t('career_paths.step_requirements') }}</label>
                  <InputText v-model="step.requirements" class="w-full !text-sm" :placeholder="t('career_paths.requirements_placeholder')" />
                </div>
              </div>
            </div>
            <div v-if="!form.steps.length" class="text-sm text-gray-400 text-center py-4 border border-dashed border-gray-300 dark:border-gray-600 rounded-lg">
              {{ t('career_paths.no_steps') }}
            </div>
          </div>
        </div>

        <FormRow :label="t('career_paths.status')" :errors="errors?.is_active">
          <ToggleSwitch v-model="form.is_active" />
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
      :title="t('career_paths.confirm_delete_path_title')"
      :message="t('career_paths.confirm_delete_path_msg')"
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
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()

// ── Daftar ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const searchTerm = ref('')

// ── Referensi ──
const organizations = ref([])

// ── Dialog ──
const dialogVisible = ref(false)
const saving = ref(false)
const editingId = ref(null)
const errors = ref({})
const form = ref(emptyForm())

// ── Detail ──
const detailVisible = ref(false)
const detailItem = ref(null)

// ── Konfirmasi hapus ──
const deleteTarget = ref(null)
const actionLoading = ref(false)
const actionError = ref('')
const deleteConfirmVisible = ref(false)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const positionOptions = computed(() =>
  organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id }))
)

const skeletonColumns = [
  { type: 'compound', width: 'w-56', headerWidth: 'w-24' },
  { type: 'text', width: 'w-64', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

function emptyForm() {
  return {
    name: '',
    description: '',
    is_active: true,
    steps: [{ position_id: null, minimum_service_months: null, requirements: '' }]
  }
}

function emptyStep() {
  return { position_id: null, minimum_service_months: null, requirements: '' }
}

function sortedSteps(steps) {
  return [...(steps || [])].sort((a, b) => (a.sequence || 0) - (b.sequence || 0))
}

// ── Load data ──
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (searchTerm.value) params.keyword = searchTerm.value
    const res = await api.get('/api/v1/tenant/career-intelligence/paths', { params })
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

function resetFilters() {
  searchTerm.value = ''
  currentPage.value = 1
  loadData()
}

// Referensi: organizations aktif (konsep organization = position).
async function loadReferences() {
  const orgRes = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
  organizations.value = orgRes.data?.data || []
}

// ── Dialog create/edit ──
function openDialog(data = null) {
  errors.value = {}
  editingId.value = data?.id || null
  if (data) {
    form.value = {
      name: data.name || '',
      description: data.description || '',
      is_active: data.is_active !== undefined ? data.is_active : true,
      steps: (data.steps || []).length
        ? sortedSteps(data.steps).map(s => ({
            position_id: s.position_id,
            minimum_service_months: s.minimum_service_months ?? null,
            requirements: s.requirements || ''
          }))
        : [emptyStep()]
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

// ── Detail ──
function openDetail(data) {
  detailItem.value = data
  detailVisible.value = true
}

// ── Editor steps ──
function addStep() {
  form.value.steps.push(emptyStep())
}

function removeStep(idx) {
  if (form.value.steps.length === 1) return
  form.value.steps.splice(idx, 1)
}

function moveStep(idx, dir) {
  const target = idx + dir
  if (target < 0 || target >= form.value.steps.length) return
  const arr = form.value.steps
  ;[arr[idx], arr[target]] = [arr[target], arr[idx]]
}

// ── Simpan ──
function validateForm() {
  const e = {}
  if (!form.value.name?.trim()) e.name = t('career_paths.field_required')
  if (!form.value.steps.length) e.steps = t('career_paths.no_steps_error')
  const seen = new Set()
  let missingPosition = false
  form.value.steps.forEach(step => {
    if (!step.position_id) {
      missingPosition = true
      return
    }
    if (seen.has(step.position_id)) {
      e.steps = t('career_paths.duplicate_position')
    }
    seen.add(step.position_id)
  })
  if (missingPosition) e.steps = t('career_paths.field_required')
  return e
}

async function handleSave() {
  errors.value = validateForm()
  if (Object.keys(errors.value).length > 0) return
  saving.value = true
  try {
    // description dikirim selalu (empty string saat dikosongkan) karena backend
    // update memakai pointer-presence (if req.Description != nil) — mengirim
    // undefined membuat deskripsi lama tidak terhapus saat diedit.
    const payload = {
      name: form.value.name.trim(),
      description: form.value.description?.trim() || '',
      is_active: form.value.is_active,
      steps: form.value.steps.map((step, idx) => ({
        position_id: step.position_id,
        sequence: idx + 1,
        minimum_service_months: step.minimum_service_months ?? undefined,
        requirements: step.requirements?.trim() || undefined
      }))
    }
    if (editingId.value) {
      await api.put(`/api/v1/tenant/career-intelligence/paths/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/career-intelligence/paths', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
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
  api.delete(`/api/v1/tenant/career-intelligence/paths/${deleteTarget.value.id}`)
    .then(async () => {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
      deleteConfirmVisible.value = false
      deleteTarget.value = null
      await loadData()
    })
    .catch(e => {
      actionError.value = getErrorMessage(e, t('message.operation_failed'))
    })
    .finally(() => { actionLoading.value = false })
}

onMounted(() => {
  loadData()
  loadReferences()
})
</script>
