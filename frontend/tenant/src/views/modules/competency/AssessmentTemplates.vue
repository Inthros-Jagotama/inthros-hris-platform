<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('competency_360.new_template')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-clone text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.templates_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('competency_360.code')" style="width:130px">
        <template #body="{data}"><Tag :value="data.code || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column :header="t('competency_360.competencies')" style="width:130px">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.competencies?.length || 0 }}</span>
        </template>
      </Column>
      <Column :header="t('competency_360.rater_types')" style="width:150px">
        <template #body="{data}">
          <div class="flex flex-wrap gap-1">
            <Tag v-for="rt in (data.rater_types || [])" :key="rt.id" :value="raterTypeLabel(rt.rater_type)" severity="secondary" class="!text-[10px] !px-1 !py-0.5" />
          </div>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'active' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:130px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-list" size="small" text severity="info" v-tooltip.left="t('competency_360.indicators')" @click="openIndicators(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('competency_360.edit_template') : t('competency_360.new_template')" modal :style="{ width: '720px' }" @hide="resetForm">
      <div class="space-y-3 max-h-[65vh] overflow-y-auto pr-1">
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('common.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
          </FormRow>
          <FormRow :label="t('competency_360.code')" required :errors="errors?.code">
            <TextInput v-model="form.code" maxlength="50" :placeholder="t('competency_360.code_placeholder')" :class="{ 'p-invalid': errors?.code }" />
          </FormRow>
        </div>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('competency_360.scale')">
          <Select v-model="form.scale_id" :options="scaleOptions" optionLabel="name" optionValue="id" showClear filter class="w-full" :placeholder="t('common.select')" />
        </FormRow>

        <!-- Competencies -->
        <div class="border-t border-gray-200 dark:border-gray-700 pt-3">
          <div class="flex items-center justify-between mb-2">
            <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.competencies') }}</p>
            <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('competency_360.add_competency')" @click="addCompetency" />
          </div>
          <div v-for="(comp, idx) in form.competencies" :key="idx" class="grid grid-cols-[1fr_90px_70px_30px] gap-2 items-center mb-2">
            <Select v-model="comp.competency_id" :options="competencyOptions" optionLabel="name" optionValue="id" filter class="w-full !text-xs" :placeholder="t('competency_360.select_competency')" />
            <TextInput v-model="comp.required_level" type="number" :placeholder="t('competency_360.req_level')" class="!text-xs" />
            <TextInput v-model="comp.weight" type="number" :placeholder="t('competency_360.weight')" class="!text-xs" />
            <Button icon="pi pi-trash" size="small" text severity="danger" @click="form.competencies.splice(idx, 1)" />
          </div>
          <p v-if="form.competencies.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_competencies_hint') }}</p>
        </div>

        <!-- Rater types -->
        <div class="border-t border-gray-200 dark:border-gray-700 pt-3">
          <div class="flex items-center justify-between mb-2">
            <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.rater_types') }}</p>
            <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('competency_360.add_rater_type')" @click="addRaterType" />
          </div>
          <div v-for="(rt, idx) in form.rater_types" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-2 mb-2">
            <div class="grid grid-cols-2 gap-2">
              <div>
                <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.rater_type') }}</label>
                <Select v-model="rt.rater_type" :options="raterTypeOptions" optionLabel="label" optionValue="value" class="w-full !text-xs" />
              </div>
              <div>
                <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.weight') }}</label>
                <TextInput v-model="rt.weight" type="number" class="!text-xs" />
              </div>
              <div>
                <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.min_rater') }}</label>
                <TextInput v-model="rt.min_rater" type="number" class="!text-xs" />
              </div>
              <div>
                <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.max_rater') }}</label>
                <TextInput v-model="rt.max_rater" type="number" class="!text-xs" />
              </div>
            </div>
            <div class="flex items-center gap-4 mt-2">
              <label class="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300">
                <ToggleSwitch v-model="rt.required" /> {{ t('competency_360.required') }}
              </label>
              <label class="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300">
                <ToggleSwitch v-model="rt.anonymous" /> {{ t('competency_360.anonymous') }}
              </label>
              <Button icon="pi pi-trash" size="small" text severity="danger" class="ml-auto" @click="form.rater_types.splice(idx, 1)" />
            </div>
          </div>
          <p v-if="form.rater_types.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_rater_types_hint') }}</p>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- Indicators dialog -->
    <Dialog v-model:visible="indicatorsVisible" :header="t('competency_360.template_indicators')" modal :style="{ width: '720px' }" @hide="indicatorsForm = []">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('competency_360.indicators_hint') }}</p>
      <div v-if="allIndicators.length === 0" class="text-center py-6 text-gray-400 dark:text-gray-500 text-sm">
        {{ t('competency_360.no_indicators_available') }}
      </div>
      <div v-else class="space-y-2 max-h-[50vh] overflow-y-auto pr-1">
        <div v-for="ind in allIndicators" :key="ind.id" class="flex items-start gap-2 border border-gray-200 dark:border-gray-700 rounded-lg p-2">
          <ToggleSwitch v-model="selectedIndicatorIds[ind.id]" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800 dark:text-gray-100">{{ ind.statement }}</p>
            <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-0.5">{{ ind.competency_name || '-' }}</p>
          </div>
          <div class="w-20 shrink-0">
            <TextInput v-model="indicatorWeights[ind.id]" type="number" :placeholder="t('competency_360.weight')" class="!text-xs" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="indicatorsVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="savingIndicators" :disabled="savingIndicators" @click="handleSaveIndicators" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('competency_360.delete_template_title')"
      :message="t('competency_360.delete_template', { name: deleteTarget?.name || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const scales = ref([])
const competencies = ref([])
const allIndicators = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const indicatorsVisible = ref(false)
const savingIndicators = ref(false)
const indicatorsTemplateId = ref(null)
const selectedIndicatorIds = ref({})
const indicatorWeights = ref({})

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const statusOptions = [
  { label: t('common_status.active'), value: 'active' },
  { label: t('common_status.inactive'), value: 'inactive' }
]

const raterTypeOptions = [
  { label: t('competency_360.rater_type_self'), value: 'self' },
  { label: t('competency_360.rater_type_superior'), value: 'superior' },
  { label: t('competency_360.rater_type_peer'), value: 'peer' },
  { label: t('competency_360.rater_type_subordinate'), value: 'subordinate' },
  { label: t('competency_360.rater_type_other'), value: 'other' }
]

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-24' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const scaleOptions = computed(() => scales.value.filter(s => s.status !== 'inactive' || s.id === form.value.scale_id))
const competencyOptions = computed(() => competencies.value)

function defaultForm() {
  return { name: '', code: '', description: '', status: 'active', scale_id: null, competencies: [], rater_types: [] }
}

function newCompetency() {
  return { competency_id: null, required_level: null, weight: 1, sort_order: form.value.competencies.length }
}

function newRaterType() {
  return { rater_type: 'peer', weight: 1, min_rater: 1, max_rater: null, required: false, anonymous: false }
}

function addCompetency() {
  form.value.competencies.push(newCompetency())
}

function addRaterType() {
  form.value.rater_types.push(newRaterType())
}

function raterTypeLabel(type) {
  const key = `competency_360.rater_type_${type}`
  return t(key) !== key ? t(key) : type
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadReferences() {
  try {
    const [scaleRes, compRes, indRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/rating-scales', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/competency/competencies', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/competency/indicators', { params: { per_page: 500 } })
    ])
    scales.value = scaleRes.status === 'fulfilled' ? (scaleRes.value.data?.data || []) : []
    competencies.value = compRes.status === 'fulfilled' ? (compRes.value.data?.data || []) : []
    allIndicators.value = indRes.status === 'fulfilled' ? (indRes.value.data?.data || []) : []
  } catch {
    // fail-silent
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/competency/templates', { params })
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? {
        name: item.name || '',
        code: item.code || '',
        description: item.description || '',
        status: item.status || 'active',
        scale_id: item.scale_id || null,
        competencies: (item.competencies || []).map(c => ({ competency_id: c.competency_id, required_level: c.required_level ?? null, weight: c.weight ?? 1, sort_order: c.sort_order || 0 })),
        rater_types: (item.rater_types || []).map(rt => ({ rater_type: rt.rater_type, weight: rt.weight ?? 1, min_rater: rt.min_rater ?? 1, max_rater: rt.max_rater ?? null, required: !!rt.required, anonymous: !!rt.anonymous }))
      }
    : defaultForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = defaultForm()
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  if (!form.value.code?.trim()) { errors.value = { code: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      code: form.value.code.trim(),
      description: form.value.description?.trim() || '',
      status: form.value.status || 'active',
      scale_id: form.value.scale_id || undefined,
      competencies: form.value.competencies
        .filter(c => c.competency_id)
        .map((c, idx) => ({
          competency_id: c.competency_id,
          required_level: c.required_level !== null && c.required_level !== '' ? Number(c.required_level) : undefined,
          weight: Number(c.weight) || 0,
          sort_order: Number(c.sort_order) || idx
        })),
      rater_types: form.value.rater_types
        .map((rt, idx) => ({
          rater_type: rt.rater_type,
          weight: Number(rt.weight) || 0,
          min_rater: Number(rt.min_rater) || 0,
          max_rater: rt.max_rater !== null && rt.max_rater !== '' ? Number(rt.max_rater) : undefined,
          required: !!rt.required,
          anonymous: !!rt.anonymous
        }))
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/competency/templates/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/competency/templates', payload)
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

async function openIndicators(template) {
  indicatorsTemplateId.value = template.id
  selectedIndicatorIds.value = {}
  indicatorWeights.value = {}
  for (const ind of allIndicators.value) {
    selectedIndicatorIds.value[ind.id] = false
  }
  try {
    const res = await api.get(`/api/v1/tenant/competency/templates/${template.id}/indicators`)
    const existing = res.data?.data || []
    for (const ti of existing) {
      selectedIndicatorIds.value[ti.indicator_id] = true
      indicatorWeights.value[ti.indicator_id] = ti.weight ?? 1
    }
  } catch {
    // template tanpa indicators — tampilkan semua unchecked
  }
  indicatorsVisible.value = true
}

async function handleSaveIndicators() {
  savingIndicators.value = true
  try {
    const payload = allIndicators.value
      .filter(ind => selectedIndicatorIds.value[ind.id])
      .map((ind, idx) => ({
        indicator_id: ind.id,
        weight: Number(indicatorWeights.value[ind.id]) || 0,
        sort_order: idx
      }))
    await api.put(`/api/v1/tenant/competency/templates/${indicatorsTemplateId.value}/indicators`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    indicatorsVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingIndicators.value = false
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/competency/templates/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
