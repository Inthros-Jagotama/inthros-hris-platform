<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="yearFilter" :options="yearOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_years')" class="!w-44" showClear @update:modelValue="onFilterChange" />
        <SelectLabel v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_status')" class="!w-44" showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.plan_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      v-model:expandedRows="expandedRows"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-calendar-plus text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.plans_empty') }}</p>
        </div>
      </template>
      <Column :expander="true" style="width: 40px" />
      <Column field="code" :header="t('training.plan_code')" style="width:130px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('training.plan_name')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="year" :header="t('training.plan_year')" style="width:90px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.year }}</span></template>
      </Column>
      <Column field="status" :header="t('training.plan_status')" style="width:120px">
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
      <template #expansion="{data: planRow}">
        <div class="p-4 bg-gray-50 dark:bg-gray-900/40">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <div class="flex items-center gap-2">
              <i class="pi pi-list text-emerald-500 text-sm"></i>
              <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.plan_items') }}</h3>
              <span v-if="planItems[data.id]?.length" class="text-xs text-gray-400">{{ planItems[data.id].length }} {{ t('common.items') }}</span>
            </div>
            <Button :label="t('training.plan_item_new')" icon="pi pi-plus" size="small" @click="openItemDialog(planRow)" />
          </div>
          <DataTable :value="planItems[data.id] || []" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="itemsLoading[data.id]">
            <template #empty>
              <div class="text-center py-6 text-sm text-gray-400">{{ t('training.plan_items_empty') }}</div>
            </template>
            <Column field="course_id" :header="t('training.item_course')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ courseName(data.course_id) }}</span></template>
            </Column>
            <Column field="target_date" :header="t('training.item_target_date')" style="width:130px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.target_date || '-' }}</span></template>
            </Column>
            <Column field="target_participants" :header="t('training.item_target_participants')" style="width:140px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.target_participants ?? '-' }}</span></template>
            </Column>
            <Column field="estimated_cost" :header="t('training.item_estimated_cost')" style="width:140px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.estimated_cost ? formatMoney(data.estimated_cost) : '-' }}</span></template>
            </Column>
            <Column field="priority" :header="t('training.item_priority')" style="width:110px">
              <template #body="{data}"><Tag :value="priorityLabel(data.priority)" :severity="prioritySeverity(data.priority)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
              <template #body="{data}">
                <div class="flex items-center gap-1 justify-end">
                  <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openItemDialog(planRow, data)" />
                  <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteItem(data)" />
                </div>
              </template>
            </Column>
          </DataTable>
        </div>
      </template>
    </DataTable>

    <!-- Dialog: plan -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.plan_edit') : t('training.plan_new')" modal :style="{ width: '520px' }" @hide="resetForm">
      <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.plan_code')" required :errors="errors?.code">
            <TextInput v-model="form.code" maxlength="30" :placeholder="t('training.plan_code_placeholder')" :class="{ 'p-invalid': errors?.code }" />
          </FormRow>
          <FormRow :label="t('training.plan_year')" required :errors="errors?.year">
            <InputNumber v-model="form.year" class="!w-full" :min="2000" :max="2100" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('training.plan_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="200" :placeholder="t('training.plan_name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('training.plan_status')">
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

    <!-- Dialog: plan item -->
    <Dialog v-model:visible="itemDialogVisible" :header="itemEditing ? t('training.plan_item_edit') : t('training.plan_item_new')" modal :style="{ width: '540px' }" @hide="resetItemForm">
      <div class="space-y-4">
        <FormRow :label="t('training.item_course')" required :errors="errors?.course_id">
          <SelectLabel v-model="itemForm.course_id" :options="courseOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.course_id }" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.item_target_date')">
            <DateInput v-model="itemForm.target_date" />
          </FormRow>
          <FormRow :label="t('training.item_priority')">
            <SelectLabel v-model="itemForm.priority" :options="priorityOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.item_target_participants')">
            <InputNumber v-model="itemForm.target_participants" class="!w-full" :min="0" size="small" />
          </FormRow>
          <FormRow :label="t('training.item_estimated_cost')">
            <InputNumber v-model="itemForm.estimated_cost" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="itemDialogVisible = false" />
          <Button :label="itemEditing ? t('common.update') : t('common.save')" size="small" :loading="itemSaving" :disabled="itemSaving" @click="handleSaveItem" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_plan', { name: deleteTarget?.name || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
    <ConfirmDeleteDialog
      v-model:visible="deleteItemDialogVisible"
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_plan_item')"
      :loading="deletingItem"
      :errorMsg="deleteItemError"
      @confirm="handleDeleteItem"
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
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const yearFilter = ref(null)
const statusFilter = ref(null)
const currentYear = new Date().getFullYear()

const expandedRows = ref({})
const planItems = ref({})
const itemsLoading = ref({})

const courses = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const itemDialogVisible = ref(false)
const itemEditing = ref(false)
const itemEditingId = ref(null)
const itemSaving = ref(false)
const itemForm = ref(defaultItemForm())
const currentPlan = ref(null)

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const deleteItemDialogVisible = ref(false)
const deletingItem = ref(false)
const deleteItemError = ref('')
const deleteItemTarget = ref(null)

const skeletonColumns = [
  { type: 'tag', width: 'w-14', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const yearOptions = computed(() => {
  const years = new Set()
  years.add(currentYear)
  items.value.forEach(i => years.add(i.year))
  return [...years].sort().map(y => ({ label: String(y), value: y }))
})
const statusOptions = computed(() => ['DRAFT', 'ACTIVE', 'ARCHIVED'].map(v => ({ label: statusLabel(v), value: v })))
const priorityOptions = computed(() => ['LOW', 'MEDIUM', 'HIGH', 'URGENT'].map(v => ({ label: priorityLabel(v), value: v })))
const courseOptions = computed(() => courses.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))

function statusLabel(s) {
  const key = `training.plan_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function statusSeverity(s) {
  switch (s) {
    case 'ACTIVE': return 'success'
    case 'DRAFT': return 'secondary'
    case 'ARCHIVED': return 'warning'
    default: return 'secondary'
  }
}
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
function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || id
}
function formatMoney(v) {
  try { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v) } catch { return v }
}

function defaultForm() {
  return { code: '', name: '', year: currentYear, description: '', status: 'DRAFT' }
}
function defaultItemForm() {
  return { course_id: null, target_date: null, target_participants: null, estimated_cost: null, priority: 'MEDIUM' }
}

async function loadCourses() {
  try {
    const res = await api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } })
    courses.value = res.data?.data || []
  } catch {
    courses.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (yearFilter.value) params.year = yearFilter.value
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/trainings/plans', { params })
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

async function loadPlanItems(planId) {
  itemsLoading.value = { ...itemsLoading.value, [planId]: true }
  try {
    const res = await api.get(`/api/v1/tenant/trainings/plans/${planId}/items`)
    planItems.value = { ...planItems.value, [planId]: res.data?.data || [] }
  } catch {
    planItems.value = { ...planItems.value, [planId]: [] }
  } finally {
    itemsLoading.value = { ...itemsLoading.value, [planId]: false }
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

function onRowExpand(event) {
  if (!planItems.value[event.data.id]) loadPlanItems(event.data.id)
}

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? { code: item.code || '', name: item.name || '', year: item.year || currentYear, description: item.description || '', status: item.status || 'DRAFT' }
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
  if (!form.value.code?.trim()) { errors.value = { code: t('form.required') }; return }
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  if (!form.value.year) { errors.value = { year: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      code: form.value.code.trim(),
      name: form.value.name.trim(),
      year: form.value.year,
      description: form.value.description?.trim() || '',
      status: form.value.status || 'DRAFT'
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/plans/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/plans', payload)
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

function openItemDialog(plan, item) {
  currentPlan.value = plan
  itemEditing.value = !!item
  itemEditingId.value = item?.id || null
  errors.value = {}
  itemForm.value = item
    ? {
        course_id: item.course_id || null,
        target_date: item.target_date || null,
        target_participants: item.target_participants ?? null,
        estimated_cost: item.estimated_cost ?? null,
        priority: item.priority || 'MEDIUM'
      }
    : defaultItemForm()
  itemDialogVisible.value = true
}

function resetItemForm() {
  itemForm.value = defaultItemForm()
  errors.value = {}
  itemEditing.value = false
  itemEditingId.value = null
  currentPlan.value = null
}

async function handleSaveItem() {
  errors.value = {}
  if (!itemForm.value.course_id) { errors.value = { course_id: t('form.required') }; return }
  itemSaving.value = true
  try {
    const payload = {
      course_id: itemForm.value.course_id,
      target_date: itemForm.value.target_date || null,
      target_participants: itemForm.value.target_participants ?? null,
      estimated_cost: itemForm.value.estimated_cost ?? null,
      priority: itemForm.value.priority || 'MEDIUM'
    }
    if (itemEditing.value) {
      await api.put(`/api/v1/tenant/trainings/plan-items/${itemEditingId.value}`, payload)
    } else {
      await api.post(`/api/v1/tenant/trainings/plans/${currentPlan.value.id}/items`, payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    itemDialogVisible.value = false
    if (currentPlan.value) await loadPlanItems(currentPlan.value.id)
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    itemSaving.value = false
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
    await api.delete(`/api/v1/tenant/trainings/plans/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

function confirmDeleteItem(item) {
  deleteItemTarget.value = item
  deleteItemError.value = ''
  deleteItemDialogVisible.value = true
}

async function handleDeleteItem() {
  deletingItem.value = true
  deleteItemError.value = ''
  try {
    await api.delete(`/api/v1/tenant/trainings/plan-items/${deleteItemTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteItemDialogVisible.value = false
    if (currentPlan.value) await loadPlanItems(currentPlan.value.id)
  } catch (e) {
    deleteItemError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deletingItem.value = false
  }
}

onMounted(() => {
  loadCourses()
  loadData()
})
</script>
