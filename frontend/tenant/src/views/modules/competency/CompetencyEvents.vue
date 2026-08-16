<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('competency_360.new_event')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.events_empty') }}</p>
        </div>
      </template>
      <Column :header="t('competency_360.event_period')" style="width:150px">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ periodLabel(data) }}</span>
        </template>
      </Column>
      <Column field="type" :header="t('competency_360.event_type')" style="width:110px">
        <template #body="{data}"><Tag :value="data.type === 'auto' ? t('competency_360.type_auto') : t('competency_360.type_manual')" :severity="data.type === 'auto' ? 'info' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="period_type" :header="t('competency_360.period_type')" style="width:110px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ periodTypeLabel(data.period_type) }}</span>
        </template>
      </Column>
      <Column :header="t('competency_360.template')" style="width:160px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ templateName(data.template_id) }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:130px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-users" size="small" text severity="info" v-tooltip.left="t('competency_360.targets')" @click="openTargets(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('competency_360.edit_event') : t('competency_360.new_event')" modal :style="{ width: '520px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('competency_360.event_type')" required :errors="errors?.type">
          <Select v-model="form.type" :options="typeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
        <FormRow :label="t('competency_360.period_type')" required :errors="errors?.period_type">
          <Select v-model="form.period_type" :options="periodTypeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
        <FormRow :label="t('competency_360.period_year')" required :errors="errors?.period_year">
          <InputNumber v-model="form.period_year" class="!w-full" :min="2000" :max="2100" size="small" />
        </FormRow>
        <FormRow :label="t('competency_360.period_number')" :errors="errors?.period_number">
          <InputNumber v-model="form.period_number" class="!w-full" :min="1" :max="12" size="small" />
        </FormRow>
        <FormRow :label="t('competency_360.template')" :errors="errors?.template_id">
          <Select v-model="form.template_id" :options="templateOptions" optionLabel="name" optionValue="id" showClear filter class="w-full" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('common.status')">
          <Select v-model="form.status" :options="eventStatusOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- Targets dialog -->
    <Dialog v-model:visible="targetsVisible" :header="t('competency_360.targets')" modal :style="{ width: '760px' }" @hide="resetTargets">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('competency_360.targets_hint') }}</p>
      <div class="flex items-center gap-2 mb-3 flex-wrap">
        <Select v-model="targetForm.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" filter class="w-64" :placeholder="t('competency_360.select_organization')" />
        <Select v-model="targetForm.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-72" :placeholder="t('competency_360.select_employee')" />
        <Button :label="t('competency_360.add_target')" icon="pi pi-plus" size="small" :loading="addingTarget" @click="addTarget" />
      </div>
      <SkeletonTable v-if="targetsLoading" :columns="targetSkeletonColumns" :rows="5" />
      <DataTable v-else :value="targets" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
            <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('competency_360.targets_empty') }}</p>
          </div>
        </template>
        <Column :header="t('competency_360.employee')">
          <template #body="{data}">
            <span class="text-gray-800 dark:text-gray-100">{{ employeeName(data.employee_id) }}</span>
          </template>
        </Column>
        <Column field="status" :header="t('common.status')" style="width:110px">
          <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column :header="t('competency_360.raters')" style="width:80px">
          <template #body="{data}">
            <Button icon="pi pi-users" size="small" text :label="String(data.rater_count ?? 0)" v-tooltip.top="t('competency_360.view_raters')" @click="goRaters(data)" />
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
          <template #body="{data}">
            <div class="flex items-center gap-1 justify-end">
              <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteTarget(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('competency_360.delete_event_title')"
      :message="t('competency_360.delete_event', { period: deleteTarget ? periodLabel(deleteTarget) : '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
    <ConfirmDeleteDialog
      v-model:visible="deleteTargetDialogVisible"
      :title="t('competency_360.delete_target_title')"
      :message="t('competency_360.delete_target')"
      :loading="deletingTarget"
      :errorMsg="deleteTargetError"
      @confirm="handleDeleteTarget"
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
import InputNumber from 'primevue/inputnumber'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const templates = ref([])
const employees = ref([])
const organizations = ref([])

const currentYear = new Date().getFullYear()

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const targetsVisible = ref(false)
const targetsLoading = ref(false)
const targets = ref([])
const targetsEvent = ref(null)
const targetForm = ref({ organization_id: null, employee_id: null })
const addingTarget = ref(false)

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const deleteTargetDialogVisible = ref(false)
const deletingTarget = ref(false)
const deleteTargetError = ref('')
const deleteTargetItem = ref(null)

const typeOptions = [
  { label: t('competency_360.type_auto'), value: 'auto' },
  { label: t('competency_360.type_manual'), value: 'manual' }
]

const periodTypeOptions = [
  { label: t('competency_360.period_annual'), value: 'annual' },
  { label: t('competency_360.period_semester'), value: 'semester' },
  { label: t('competency_360.period_quarter'), value: 'quarter' }
]

const eventStatusOptions = [
  { label: t('common_status.draft'), value: 'draft' },
  { label: t('common_status.active'), value: 'active' },
  { label: t('common_status.closed'), value: 'closed' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-24' }
]

const targetSkeletonColumns = [
  { type: 'text', width: 'w-48', headerWidth: 'w-32' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const templateOptions = computed(() => templates.value.filter(t => t.status !== 'inactive' || t.id === form.value.template_id))
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.employee_id })))
const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.name || o.nomenclature || o.code, value: o.id })))

function defaultForm() {
  return { type: 'manual', period_type: 'annual', period_year: currentYear, period_number: null, status: 'draft', template_id: null }
}

function periodLabel(e) {
  const parts = [e.period_year]
  if (e.period_number) parts.unshift(`${e.period_type === 'quarter' ? 'Q' : e.period_type === 'semester' ? 'S' : 'P'}${e.period_number}`)
  return parts.join(' ')
}

function periodTypeLabel(type) {
  const key = `competency_360.period_${type}`
  return t(key) !== key ? t(key) : type
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'closed': return 'secondary'
    case 'finalized': return 'success'
    case 'submitted': return 'info'
    case 'pending_approval': return 'warn'
    case 'rejected': return 'danger'
    default: return 'secondary'
  }
}

function templateName(id) {
  return templates.value.find(t => t.id === id)?.name || '-'
}

function employeeName(id) {
  if (!id) return '-'
  return employees.value.find(e => e.employee_id === id)?.name || id.slice(0, 8)
}

async function loadReferences() {
  try {
    const [tplRes, empRes, orgRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/templates', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
    ])
    templates.value = tplRes.status === 'fulfilled' ? (tplRes.value.data?.data || []) : []
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
    organizations.value = orgRes.status === 'fulfilled' ? (orgRes.value.data?.data || []) : []
  } catch {
    // fail-silent
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/competency/events', { params })
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
        type: item.type || 'manual',
        period_type: item.period_type || 'annual',
        period_year: item.period_year || currentYear,
        period_number: item.period_number ?? null,
        status: item.status || 'draft',
        template_id: item.template_id || null
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
  if (!form.value.type) { errors.value = { type: t('form.required') }; return }
  if (!form.value.period_type) { errors.value = { period_type: t('form.required') }; return }
  if (!form.value.period_year) { errors.value = { period_year: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      type: form.value.type,
      period_type: form.value.period_type,
      period_year: Number(form.value.period_year),
      period_number: form.value.period_number ? Number(form.value.period_number) : undefined,
      status: form.value.status || 'draft',
      template_id: form.value.template_id || undefined
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/competency/events/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/competency/events', payload)
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

async function loadTargets() {
  if (!targetsEvent.value) return
  targetsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/competency/event-targets', { params: { per_page: 500 } })
    const all = res.data?.data || []
    targets.value = all.filter(t => t.competency_event_id === targetsEvent.value.id)
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    targetsLoading.value = false
  }
}

async function openTargets(event) {
  targetsEvent.value = event
  targets.value = []
  targetForm.value = { organization_id: organizations.value[0]?.id || null, employee_id: null }
  targetsVisible.value = true
  await loadTargets()
}

function resetTargets() {
  targets.value = []
  targetsEvent.value = null
  targetForm.value = { organization_id: null, employee_id: null }
}

async function addTarget() {
  if (!targetForm.value.organization_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('competency_360.select_organization'), life: 3000 })
    return
  }
  addingTarget.value = true
  try {
    await api.post('/api/v1/tenant/competency/event-targets', {
      competency_event_id: targetsEvent.value.id,
      organization_id: targetForm.value.organization_id,
      employee_id: targetForm.value.employee_id || undefined
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    targetForm.value = { organization_id: organizations.value[0]?.id || null, employee_id: null }
    await loadTargets()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    addingTarget.value = false
  }
}

function goRaters(target) {
  targetsVisible.value = false
  router.push({ path: '/competencies/raters', query: { target_id: target.id } })
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
    await api.delete(`/api/v1/tenant/competency/events/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

function confirmDeleteTarget(item) {
  deleteTargetItem.value = item
  deleteTargetError.value = ''
  deleteTargetDialogVisible.value = true
}

async function handleDeleteTarget() {
  deletingTarget.value = true
  deleteTargetError.value = ''
  try {
    await api.delete(`/api/v1/tenant/competency/event-targets/${deleteTargetItem.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteTargetDialogVisible.value = false
    await loadTargets()
  } catch (e) {
    deleteTargetError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deletingTarget.value = false
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
