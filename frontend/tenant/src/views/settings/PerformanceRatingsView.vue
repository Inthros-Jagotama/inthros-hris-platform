<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('performance_ratings.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable v-else :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="sort_order" :sortOrder="1">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-star text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('performance_ratings.empty_title') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('performance_ratings.code')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="name" :header="t('performance_ratings.name')" sortable>
        <template #body="{data}">
          <div class="flex items-center gap-2">
            <div v-if="data.color" class="w-3 h-3 rounded-full" :class="getColorClass(data.color)"></div>
            <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>
      <Column field="min_score" :header="t('performance_ratings.min_score')" sortable style="width:100px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.min_score }}</span>
        </template>
      </Column>
      <Column field="max_score" :header="t('performance_ratings.max_score')" sortable style="width:100px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.max_score }}</span>
        </template>
      </Column>
      <Column field="sort_order" :header="t('performance_ratings.sort_order')" sortable style="width:80px">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('performance_ratings.edit') : t('performance_ratings.new')" modal :style="{ width: '560px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <FormRow :label="t('performance_ratings.code')" required :errors="errors?.code">
            <TextInput v-model="form.code" maxlength="20" autofocus :placeholder="t('performance_ratings.code_placeholder')" :class="{'p-invalid':errors?.code}" />
          </FormRow>
          <FormRow :label="t('performance_ratings.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="100" :placeholder="t('performance_ratings.name_placeholder')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
        </div>
        <div class="grid grid-cols-3 gap-4">
          <FormRow :label="t('performance_ratings.min_score')" required :errors="errors?.min_score">
            <InputNumber v-model="form.min_score" class="!w-full" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
          </FormRow>
          <FormRow :label="t('performance_ratings.max_score')" required :errors="errors?.max_score">
            <InputNumber v-model="form.max_score" class="!w-full" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
          </FormRow>
          <FormRow :label="t('performance_ratings.sort_order')" :errors="errors?.sort_order">
            <InputNumber v-model="form.sort_order" class="!w-full" :min="0" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('performance_ratings.color')" :errors="errors?.color">
          <Select v-model="form.color" :options="colorOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full">
            <template #value="{ value }">
              <div v-if="value" class="flex items-center gap-2">
                <div class="w-4 h-4 rounded" :class="getColorClass(value)"></div>
                <span>{{ colorOptions.find(c => c.value === value)?.label }}</span>
              </div>
              <span v-else>{{ t('common.select') }}</span>
            </template>
            <template #option="{ option }">
              <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded" :class="getColorClass(option.value)"></div>
                <span>{{ option.label }}</span>
              </div>
            </template>
          </Select>
        </FormRow>
        <FormRow :label="t('performance_ratings.description_label')" :errors="errors?.description">
          <Textarea v-model="form.description" rows="2" :placeholder="t('performance_ratings.description_placeholder')" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('performance_ratings.delete_title')"
      :message="t('performance_ratings.delete_message', { name: deleteTarget?.name })"
      :loading="deleting"
      :error="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const form = ref({
  code: '',
  name: '',
  min_score: 0,
  max_score: 100,
  color: 'info',
  description: '',
  sort_order: 0
})

const colorOptions = [
  { label: 'Success (Green)', value: 'success' },
  { label: 'Primary (Blue)', value: 'primary' },
  { label: 'Info (Cyan)', value: 'info' },
  { label: 'Warning (Yellow)', value: 'warning' },
  { label: 'Danger (Red)', value: 'danger' }
]

const skeletonColumns = [
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function getColorClass(color) {
  switch (color) {
    case 'success': return 'bg-emerald-500'
    case 'primary': return 'bg-blue-500'
    case 'info': return 'bg-cyan-500'
    case 'warning': return 'bg-amber-500'
    case 'danger': return 'bg-red-500'
    default: return 'bg-gray-400'
  }
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/ratings', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
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
  form.value = {
    code: item?.code || '',
    name: item?.name || '',
    min_score: item?.min_score ?? 0,
    max_score: item?.max_score ?? 100,
    color: item?.color || 'info',
    description: item?.description || '',
    sort_order: item?.sort_order || 0
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = {
    code: '',
    name: '',
    min_score: 0,
    max_score: 100,
    color: 'info',
    description: '',
    sort_order: 0
  }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.code?.trim()) {
    errors.value = { code: [t('form.required')] }
    return
  }
  if (!form.value.name?.trim()) {
    errors.value = { name: [t('form.required')] }
    return
  }

  saving.value = true
  try {
    const payload = {
      code: form.value.code,
      name: form.value.name,
      min_score: form.value.min_score,
      max_score: form.value.max_score,
      color: form.value.color || null,
      description: form.value.description || null,
      sort_order: form.value.sort_order || 0
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/performance/ratings/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_ratings.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/performance/ratings', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_ratings.created'), life: 3000 })
    }
    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
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
    await api.delete(`/api/v1/tenant/performance/ratings/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_ratings.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
