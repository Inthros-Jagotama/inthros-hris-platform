<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('performance_perspectives.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable v-else :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="sort_order" :sortOrder="1">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-th-large text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('performance_perspectives.empty_title') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('performance_perspectives.name')" sortable>
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="description" :header="t('performance_perspectives.description')">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400 text-xs line-clamp-2">{{ data.description || '-' }}</span>
        </template>
      </Column>
      <Column field="sort_order" :header="t('performance_perspectives.sort_order')" sortable style="width:100px">
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

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('performance_perspectives.edit') : t('performance_perspectives.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('performance_perspectives.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="100" autofocus :placeholder="t('performance_perspectives.name_placeholder')" :class="{'p-invalid':errors?.name}" />
        </FormRow>
        <FormRow :label="t('performance_perspectives.description')" :errors="errors?.description">
          <Textarea v-model="form.description" rows="3" :placeholder="t('performance_perspectives.description_placeholder')" class="w-full" />
        </FormRow>
        <FormRow :label="t('performance_perspectives.sort_order')" :errors="errors?.sort_order">
          <InputNumber v-model="form.sort_order" class="!w-full" :min="0" size="small" />
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
      :title="t('performance_perspectives.delete_title')"
      :message="t('performance_perspectives.delete_message', { name: deleteTarget?.name })"
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
import Dialog from 'primevue/dialog'
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
  name: '',
  description: '',
  sort_order: 0
})

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-48', headerWidth: 'w-24' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/perspectives', {
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
    name: item?.name || '',
    description: item?.description || '',
    sort_order: item?.sort_order || 0
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { name: '', description: '', sort_order: 0 }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.name?.trim()) {
    errors.value = { name: [t('form.required')] }
    return
  }

  saving.value = true
  try {
    const payload = {
      name: form.value.name,
      description: form.value.description || null,
      sort_order: form.value.sort_order || 0
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/performance/perspectives/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_perspectives.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/performance/perspectives', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_perspectives.created'), life: 3000 })
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
    await api.delete(`/api/v1/tenant/performance/perspectives/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_perspectives.deleted'), life: 3000 })
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
