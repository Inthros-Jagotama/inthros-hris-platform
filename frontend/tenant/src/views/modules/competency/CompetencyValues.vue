<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('competency_360.new_scale')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-sliders-h text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.scales_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('competency_360.code')" style="width:130px">
        <template #body="{data}"><Tag :value="data.code || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column :header="t('competency_360.scale_items')" style="width:120px">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.items?.length || 0 }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'active' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('competency_360.edit_scale') : t('competency_360.new_scale')" modal :style="{ width: '620px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="255" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('common.status')">
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>

        <div class="border-t border-gray-200 dark:border-gray-700 pt-3">
          <div class="flex items-center justify-between mb-2">
            <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.scale_items') }}</p>
            <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('competency_360.add_item')" @click="addItem" />
          </div>
          <div class="flex items-center gap-2 mb-1 px-1">
            <div class="w-16 shrink-0 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.item_value') }}</div>
            <div class="flex-1 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.item_label') }}</div>
            <div class="w-20 shrink-0 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.weight') }}</div>
            <div class="w-8 shrink-0"></div>
          </div>
          <div v-for="(item, idx) in form.items" :key="idx" class="flex items-start gap-2 mb-2">
            <div class="w-16 shrink-0">
              <TextInput v-model="item.value" type="number" :placeholder="t('competency_360.item_value')" />
            </div>
            <div class="flex-1">
              <TextInput v-model="item.label" :placeholder="t('competency_360.item_label')" />
            </div>
            <div class="w-20 shrink-0">
              <TextInput v-model="item.weight" type="number" :placeholder="t('competency_360.weight')" />
            </div>
            <Button icon="pi pi-trash" size="small" text severity="danger" @click="form.items.splice(idx, 1)" />
          </div>
          <p v-if="form.items.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_items_hint') }}</p>
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
      :title="t('competency_360.delete_scale_title')"
      :message="t('competency_360.delete_scale', { name: deleteTarget?.name || '' })"
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

const statusOptions = [
  { label: t('common_status.active'), value: 'active' },
  { label: t('common_status.inactive'), value: 'inactive' }
]

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function defaultForm() {
  return { name: '', description: '', status: 'active', items: [] }
}

function newItem() {
  return { value: null, label: '', description: '', weight: 1, sort_order: 0 }
}

function addItem() {
  form.value.items.push(newItem())
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/competency/rating-scales', { params })
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
        description: item.description || '',
        status: item.status || 'active',
        items: (item.items || []).map(i => ({ value: i.value, label: i.label, description: i.description || '', weight: i.weight ?? 1, sort_order: i.sort_order || 0 }))
      }
    : defaultForm()
  if (form.value.items.length === 0) addItem()
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
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      description: form.value.description?.trim() || '',
      status: form.value.status || 'active',
      items: form.value.items
        .filter(i => i.label?.trim())
        .map((i, idx) => ({
          value: Number(i.value) || 0,
          label: i.label.trim(),
          description: i.description?.trim() || '',
          weight: Number(i.weight) || 0,
          sort_order: Number(i.sort_order) || idx
        }))
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/competency/rating-scales/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/competency/rating-scales', payload)
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
    await api.delete(`/api/v1/tenant/competency/rating-scales/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
