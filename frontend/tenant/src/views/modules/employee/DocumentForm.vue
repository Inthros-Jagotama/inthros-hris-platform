<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_documents') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.document_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="openDialog()" />
    </div>

    <!-- DataTable for saved items -->
    <DataTable :value="items" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-file text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee.no_documents') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('employee.document_name')" sortable>
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="file" :header="t('employee.document_file')">
        <template #body="{data}">
          <a v-if="data.file" :href="data.file" target="_blank" class="text-primary-600 dark:text-primary-400 text-xs underline hover:text-primary-800 inline-flex items-center gap-1">
            <i class="pi pi-download text-xs"></i>
            <span class="font-mono">{{ getFileName(data.file) }}</span>
          </a>
          <span v-else class="text-gray-400 text-xs">-</span>
        </template>
      </Column>
      <Column field="note" :header="t('employee.note')">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.note || '-' }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:110px" frozen alignFrozen="right">
        <template #body="{data, index}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEditDialog(data, index)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="onDeleteClick(index)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Add/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('employee.edit_document') : t('employee.new_document')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="grid grid-cols-1 gap-3 py-2">
        <FormRow :label="t('employee.document_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="255" :placeholder="t('employee.document_name_placeholder')" autofocus :class="{'p-invalid':errors?.name}" />
        </FormRow>
        <FormRow :label="t('employee.document_file')" required :errors="errors?.file">
          <div class="flex items-center gap-2">
            <input
              ref="fileInputRef"
              type="file"
              class="hidden"
              :accept="'.pdf,.doc,.docx,.xls,.xlsx,.jpg,.jpeg,.png,.gif,.txt'"
              @change="onFileSelected"
            />
            <Button
              icon="pi pi-upload"
              size="small"
              severity="secondary"
              outlined
              :label="selectedFile ? selectedFile.name : t('employee.choose_file')"
              @click="fileInputRef.click()"
              class="!justify-start max-w-full overflow-hidden"
            />
            <Button
              v-if="selectedFile"
              icon="pi pi-times"
              size="small"
              text
              severity="danger"
              @click="clearSelectedFile"
            />
          </div>
        </FormRow>
        <FormRow :label="t('employee.note')" :errors="errors?.note">
          <TextInput v-model="form.note" maxlength="255" :placeholder="t('employee.note_placeholder')" :class="{'p-invalid':errors?.note}" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="deleteLoading"
      :error-msg="deleteError"
      :title="t('common.confirm')"
      :message="t('employee.confirm_delete_document')"
      @confirm="confirmDeleteDocument"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>
<script setup>
import { useI18n } from '@/composables/useI18n'
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import Dialog from 'primevue/dialog'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { getValidationErrors } from '@/services/responseHandler'

const { t } = useI18n()
const toast = useToast()

const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false },
  employeeId: { type: String, default: '' }
})
const emit = defineEmits(['update:items', 'save'])

// ── File upload state ──
const fileInputRef = ref(null)
const selectedFile = ref(null)

function onFileSelected(e) {
  const file = e.target.files?.[0]
  if (file) {
    selectedFile.value = file
    errors.value.file = null
  }
}

function clearSelectedFile() {
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

// ── Dialog state ──
const dialogVisible = ref(false)
const editingId = ref(null)
const editingIdx = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ name: '', note: '' })
const isEditing = computed(() => editingId.value !== null)

function getFileName(filePath) {
  if (!filePath) return '-'
  const parts = filePath.split('/')
  return parts[parts.length - 1]
}

function openDialog() {
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
  form.value = { name: '', note: '' }
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
  dialogVisible.value = true
}

function openEditDialog(item, idx) {
  editingId.value = item._id || item.id || null
  editingIdx.value = idx
  errors.value = {}
  form.value = {
    name: item.name || '',
    note: item.note || ''
  }
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
  dialogVisible.value = true
}

function resetForm() {
  form.value = { name: '', note: '' }
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
  selectedFile.value = null
  if (fileInputRef.value) fileInputRef.value.value = ''
}

async function handleSave() {
  errors.value = {}
  if (!form.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  saving.value = true
  try {
    const updated = [...props.items]

    if (isEditing.value && editingIdx.value !== null) {
      // ── EDIT mode ──
      if (selectedFile.value) {
        // File changed → use multipart upload endpoint
        const fd = new FormData()
        fd.append('name', form.value.name)
        fd.append('note', form.value.note || '')
        fd.append('file', selectedFile.value)
        const res = await api.put(`/api/v1/tenant/employees/${props.employeeId}/documents/${editingId.value}/upload`, fd, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        const filePath = res.data?.data?.file || ''
        updated[editingIdx.value] = {
          ...updated[editingIdx.value],
          name: form.value.name,
          note: form.value.note || '',
          file: filePath,
          _saved: true,
          _id: editingId.value
        }
      } else {
        // No file change → use JSON endpoint
        const res = await api.put(`/api/v1/tenant/employees/${props.employeeId}/documents/${editingId.value}`, {
          name: form.value.name,
          note: form.value.note || null
        })
        const updatedItem = res.data?.data || { ...form.value, _saved: true, _id: editingId.value }
        updatedItem._saved = true
        updatedItem._id = updatedItem.id || editingId.value
        updated[editingIdx.value] = { ...updated[editingIdx.value], ...updatedItem }
      }
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.updated'), life: 2000 })
    } else {
      // ── CREATE mode ──
      if (!selectedFile.value) {
        errors.value = { file: [t('form.required')] }
        saving.value = false
        return
      }
      const fd = new FormData()
      fd.append('name', form.value.name)
      fd.append('note', form.value.note || '')
      fd.append('file', selectedFile.value)
      const res = await api.post(`/api/v1/tenant/employees/${props.employeeId}/documents/upload`, fd, {
        headers: { 'Content-Type': 'multipart/form-data' }
      })
      const newItem = res.data?.data || { name: form.value.name, note: form.value.note || '' }
      newItem._saved = true
      newItem._id = newItem.id || ''
      updated.push(newItem)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.saved'), life: 2000 })
    }

    emit('update:items', updated)
    dialogVisible.value = false
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    saving.value = false
  }
}

// ── Delete state ──
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const deleteError = ref('')
const deleteTargetIdx = ref(null)

function onDeleteClick(idx) {
  deleteTargetIdx.value = idx
  deleteError.value = ''
  deleteDialogVisible.value = true
}

function removeItem(idx) {
  emit('update:items', props.items.filter((_, i) => i !== idx))
}

async function confirmDeleteDocument() {
  const idx = deleteTargetIdx.value
  if (idx === null || idx === undefined) return
  const item = props.items[idx]
  if (props.employeeId && item._id) {
    deleteLoading.value = true
    deleteError.value = ''
    try {
      await api.delete(`/api/v1/tenant/employees/${props.employeeId}/documents/${item._id}`)
      removeItem(idx)
      deleteDialogVisible.value = false
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.deleted'), life: 2000 })
    } catch (e) {
      deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
    } finally {
      deleteLoading.value = false
    }
  } else {
    removeItem(idx)
    deleteDialogVisible.value = false
  }
}
</script>