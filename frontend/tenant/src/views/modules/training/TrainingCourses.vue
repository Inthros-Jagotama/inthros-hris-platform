<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="categoryFilter" :options="categoryOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_categories')" class="!w-56" showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.course_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-book text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.courses_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('training.code')" style="width:110px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="category_name" :header="t('training.category')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ categoryName(data.category_id) }}</span></template>
      </Column>
      <Column field="course_type" :header="t('training.course_type')" style="width:130px">
        <template #body="{data}">
          <Tag v-if="data.course_type" :value="courseTypeLabel(data.course_type)" severity="info" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column field="delivery_type" :header="t('training.delivery_type')" style="width:120px">
        <template #body="{data}">
          <Tag v-if="data.delivery_type" :value="deliveryTypeLabel(data.delivery_type)" severity="warning" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column field="duration_hour" :header="t('training.duration_hour')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.duration_hour ? data.duration_hour + ' h' : '-' }}</span></template>
      </Column>
      <Column field="cost" :header="t('training.cost')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.cost ? formatMoney(data.cost) : '-' }}</span></template>
      </Column>
      <Column field="is_mandatory" :header="t('training.is_mandatory')" style="width:100px">
        <template #body="{data}"><Tag :value="data.is_mandatory ? t('common.yes') : t('common.no')" :severity="data.is_mandatory ? 'danger' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="is_active" :header="t('training.is_active')" style="width:90px">
        <template #body="{data}"><Tag :value="data.is_active ? t('common.yes') : t('common.no')" :severity="data.is_active ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:140px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-calendar" size="small" text severity="secondary" v-tooltip.left="t('training.view_sessions')" @click="router.push(`/training/sessions?course_id=${data.id}`)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.course_edit') : t('training.course_new')" modal :style="{ width: '640px' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('training.category')" required :errors="errors?.category_id">
          <SelectLabel v-model="form.category_id" :options="categoryOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.category_id }" />
        </FormRow>
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="200" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
          <FormRow :label="t('training.course_type')">
            <SelectLabel v-model="form.course_type" :options="courseTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
          </FormRow>
          <FormRow :label="t('training.delivery_type')">
            <div class="flex flex-wrap gap-3">
              <RadioLabel v-for="opt in deliveryTypeOptions" :key="opt.value" v-model="form.delivery_type" :value="opt.value" :label="opt.label" :id="'course-delivery-' + opt.value" />
            </div>
          </FormRow>
          <FormRow :label="t('training.duration_hour')" :errors="errors?.duration_hour">
            <InputNumber v-model="form.duration_hour" class="!w-full" :min="0" :maxFractionDigits="1" size="small" />
          </FormRow>
          <FormRow :label="t('training.min_score')" :errors="errors?.min_score">
            <InputNumber v-model="form.min_score" class="!w-full" :min="0" :max="100" :maxFractionDigits="2" size="small" />
          </FormRow>
          <FormRow :label="t('training.cost')" :errors="errors?.cost">
            <InputNumber v-model="form.cost" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
          <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
            <div>
              <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ t('training.is_certified') }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.is_certified_desc') }}</p>
            </div>
            <ToggleSwitch v-model="form.is_certified" />
          </div>
          <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
            <div>
              <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ t('training.is_mandatory') }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.is_mandatory_desc') }}</p>
            </div>
            <ToggleSwitch v-model="form.is_mandatory" />
          </div>
          <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
            <div>
              <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ t('training.is_active') }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.is_active_desc') }}</p>
            </div>
            <ToggleSwitch v-model="form.is_active" />
          </div>
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
      :message="t('training.confirm_delete_course', { name: deleteTarget?.name || '' })"
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
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import RadioLabel from '@/components/RadioLabel.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const categories = ref([])
const categoryFilter = ref(null)

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
  { type: 'tag', width: 'w-14', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const categoryOptions = computed(() => categories.value.map(c => ({ label: c.name, value: c.id })))

const courseTypeOptions = computed(() => ['TECHNICAL', 'SOFT_SKILL', 'COMPLIANCE', 'MANAGEMENT', 'CERTIFICATION', 'OTHER'].map(v => ({ label: courseTypeLabel(v), value: v })))
const deliveryTypeOptions = computed(() => ['IN_HOUSE', 'EXTERNAL', 'BOTH'].map(v => ({ label: deliveryTypeLabel(v), value: v })))

function courseTypeLabel(type) {
  const key = `training.type_${String(type).toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function deliveryTypeLabel(type) {
  const key = `training.delivery_${String(type).toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function categoryName(id) {
  return categories.value.find(c => c.id === id)?.name || id
}
function formatMoney(v) {
  try { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v) } catch { return v }
}

function defaultForm() {
  return {
    category_id: null, name: '', description: '',
    course_type: null, delivery_type: null,
    duration_hour: null, min_score: null, cost: null,
    is_certified: false, is_mandatory: false, is_active: true
  }
}

async function loadCategories() {
  try {
    const res = await api.get('/api/v1/tenant/trainings/categories', { params: { per_page: 500 } })
    categories.value = res.data?.data || []
  } catch {
    categories.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (categoryFilter.value) params.category_id = categoryFilter.value
    const res = await api.get('/api/v1/tenant/trainings/courses', { params })
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
        category_id: item.category_id || null,
        name: item.name || '',
        description: item.description || '',
        course_type: item.course_type || null,
        delivery_type: item.delivery_type || null,
        duration_hour: item.duration_hour || null,
        min_score: item.min_score || null,
        cost: item.cost || null,
        is_certified: !!item.is_certified,
        is_mandatory: !!item.is_mandatory,
        is_active: !!item.is_active
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
  if (!form.value.category_id) { errors.value = { category_id: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      category_id: form.value.category_id,
      name: form.value.name.trim(),
      description: form.value.description?.trim() || '',
      course_type: form.value.course_type || null,
      delivery_type: form.value.delivery_type || null,
      duration_hour: form.value.duration_hour,
      min_score: form.value.min_score,
      cost: form.value.cost,
      is_certified: form.value.is_certified,
      is_mandatory: form.value.is_mandatory,
      is_active: form.value.is_active
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/courses/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/courses', payload)
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
    await api.delete(`/api/v1/tenant/trainings/courses/${deleteTarget.value.id}`)
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
  loadCategories()
  loadData()
})
</script>
