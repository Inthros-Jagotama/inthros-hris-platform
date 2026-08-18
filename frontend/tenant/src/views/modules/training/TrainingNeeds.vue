<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_need_status')" class="!w-44" showClear @update:modelValue="onFilterChange" />
        <SelectLabel v-model="sourceFilter" :options="sourceOptions" optionLabel="label" optionValue="value" :placeholder="t('training.need_source_type')" class="!w-48" showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.need_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-bullseye text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.needs_empty') }}</p>
        </div>
      </template>
      <Column field="employee_id" :header="t('training.need_employee')" style="width:200px">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
      </Column>
      <Column field="course_id" :header="t('training.need_course')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ courseName(data.course_id) }}</span></template>
      </Column>
      <Column field="reason" :header="t('training.need_reason')">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300 line-clamp-1">{{ data.reason || '-' }}</span></template>
      </Column>
      <Column field="priority" :header="t('training.need_priority')" style="width:110px">
        <template #body="{data}"><Tag :value="priorityLabel(data.priority)" :severity="prioritySeverity(data.priority)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="source_type" :header="t('training.need_source_type')" style="width:130px">
        <template #body="{data}"><Tag :value="sourceLabel(data.source_type)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="status" :header="t('training.need_status')" style="width:110px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.need_edit') : t('training.need_new')" modal :style="{ width: '600px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('training.need_employee')">
          <SelectLabel v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_course')">
          <SelectLabel v-model="form.course_id" :options="courseOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_organization')">
          <SelectLabel v-model="form.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_position')">
          <SelectLabel v-model="form.position_id" :options="positionOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_reason')">
          <TextInput v-model="form.reason" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('training.need_priority')">
          <SelectLabel v-model="form.priority" :options="priorityOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_source_type')">
          <SelectLabel v-model="form.source_type" :options="sourceOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('training.need_status')">
          <SelectLabel v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_need')"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
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
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const statusFilter = ref(null)
const sourceFilter = ref(null)

const employees = ref([])
const courses = ref([])
const organizations = ref([])
const positions = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-52', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const courseOptions = computed(() => courses.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))
const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code, value: o.id })))
const positionOptions = computed(() => positions.value.map(p => ({ label: p.nomenclature || p.full_code, value: p.id })))
const priorityOptions = computed(() => ['LOW', 'MEDIUM', 'HIGH', 'URGENT'].map(v => ({ label: priorityLabel(v), value: v })))
const sourceOptions = computed(() => ['MANUAL', 'PERFORMANCE', 'COMPETENCY', 'CAREER', 'SUCCESSION', 'COMPLIANCE', 'WORKFORCE', 'ONBOARDING'].map(v => ({ label: sourceLabel(v), value: v })))
const statusOptions = computed(() => ['OPEN', 'PLANNED', 'FULFILLED', 'CANCELLED'].map(v => ({ label: statusLabel(v), value: v })))

function priorityLabel(p) {
  const key = `training.priority_${String(p || '').toLowerCase()}`
  return t(key) !== key ? t(key) : p
}
function prioritySeverity(p) {
  switch (p) {
    case 'URGENT': return 'danger'
    case 'HIGH': return 'warning'
    case 'MEDIUM': return 'info'
    default: return 'secondary'
  }
}
function sourceLabel(s) {
  const key = `training.need_source_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusLabel(s) {
  const key = `training.need_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusSeverity(s) {
  switch (s) {
    case 'OPEN': return 'info'
    case 'PLANNED': return 'warning'
    case 'FULFILLED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || (id ? id : '-')
}
function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || (id ? id : '-')
}

function defaultForm() {
  return { employee_id: null, organization_id: null, position_id: null, course_id: null, reason: '', priority: 'MEDIUM', source_type: 'MANUAL', status: 'OPEN' }
}

async function loadReferences() {
  const [eRes, cRes, oRes, pRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/organizations', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/organization-summaries', { params: { per_page: 500 } })
  ])
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
  organizations.value = oRes.status === 'fulfilled' ? (oRes.value.data?.data || []) : []
  positions.value = pRes.status === 'fulfilled' ? (pRes.value.data?.data || []) : []
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    if (sourceFilter.value) params.source_type = sourceFilter.value
    const res = await api.get('/api/v1/tenant/trainings/needs', { params })
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? {
        employee_id: item.employee_id || null,
        organization_id: item.organization_id || null,
        position_id: item.position_id || null,
        course_id: item.course_id || null,
        reason: item.reason || '',
        priority: item.priority || 'MEDIUM',
        source_type: item.source_type || 'MANUAL',
        status: item.status || 'OPEN'
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
  if (!form.value.employee_id && !form.value.organization_id && !form.value.position_id) {
    errors.value = { employee_id: t('form.required') }
    return
  }
  saving.value = true
  try {
    const payload = {
      employee_id: form.value.employee_id || null,
      organization_id: form.value.organization_id || null,
      position_id: form.value.position_id || null,
      course_id: form.value.course_id || null,
      reason: form.value.reason?.trim() || '',
      priority: form.value.priority || 'MEDIUM',
      source_type: form.value.source_type || 'MANUAL',
      status: form.value.status || 'OPEN'
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/needs/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/needs', payload)
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

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/trainings/needs/${deleteTarget.value.id}`)
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
