<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_experience') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.experience_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="openDialog()" />
    </div>

    <!-- DataTable for saved items -->
    <DataTable :value="items" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-briefcase text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee.no_experience') }}</p>
        </div>
      </template>
      <Column field="company" :header="t('employee.company')" sortable>
        <template #body="{data}">
          <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.company }}</span>
        </template>
      </Column>
      <Column field="position" :header="t('employee.position')">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.position || '-' }}</span>
        </template>
      </Column>
      <Column field="start_year" :header="t('employee.start_year')" style="width:100px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.start_year || '-' }}</span>
        </template>
      </Column>
      <Column field="end_year" :header="t('employee.end_year')" style="width:100px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.end_year || '-' }}</span>
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
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('employee.edit_experience') : t('employee.new_experience')" modal :style="{ width: '560px' }" :closable="true" @hide="resetForm">
      <div class="grid grid-cols-1 gap-3 py-2">
        <FormRow :label="t('employee.company')" required :errors="errors?.company">
          <TextInput v-model="form.company" maxlength="255" :placeholder="t('employee.company_placeholder')" autofocus :class="{'p-invalid':errors?.company}" />
        </FormRow>
        <FormRow :label="t('employee.position')" :errors="errors?.position">
          <TextInput v-model="form.position" maxlength="255" :placeholder="t('employee.position_placeholder')" :class="{'p-invalid':errors?.position}" />
        </FormRow>
        <FormRow :label="t('employee.start_year')" :errors="errors?.start_year">
          <TextInput v-model="form.start_year" maxlength="4" :placeholder="t('employee.year_placeholder')" :class="{'p-invalid':errors?.start_year}" />
        </FormRow>
        <FormRow :label="t('employee.end_year')" :errors="errors?.end_year">
          <TextInput v-model="form.end_year" maxlength="4" :placeholder="t('employee.year_placeholder')" :class="{'p-invalid':errors?.end_year}" />
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
      :message="t('employee.confirm_delete_experience')"
      @confirm="confirmDeleteExperience"
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

// ── Dialog state ──
const dialogVisible = ref(false)
const editingId = ref(null)
const editingIdx = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ company: '', position: '', start_year: '', end_year: '' })
const isEditing = computed(() => editingId.value !== null)

function openDialog() {
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
  form.value = { company: '', position: '', start_year: '', end_year: '' }
  dialogVisible.value = true
}

function openEditDialog(item, idx) {
  editingId.value = item._id || item.id || null
  editingIdx.value = idx
  errors.value = {}
  form.value = {
    company: item.company || '',
    position: item.position || '',
    start_year: item.start_year ? String(item.start_year) : '',
    end_year: item.end_year ? String(item.end_year) : ''
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { company: '', position: '', start_year: '', end_year: '' }
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  if (!form.value.company?.trim()) { errors.value = { company: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      company: form.value.company,
      position: form.value.position || null,
      start_year: form.value.start_year ? parseInt(form.value.start_year) : null,
      end_year: form.value.end_year ? parseInt(form.value.end_year) : null
    }
    const updated = [...props.items]
    if (isEditing.value && editingIdx.value !== null) {
      const res = await api.put(`/api/v1/tenant/employees/${props.employeeId}/experiences/${editingId.value}`, payload)
      const updatedItem = res.data?.data || { ...payload, _saved: true, _id: editingId.value }
      updatedItem._saved = true
      updatedItem._id = updatedItem.id || editingId.value
      updated[editingIdx.value] = { ...updated[editingIdx.value], ...updatedItem }
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.updated'), life: 2000 })
    } else {
      const res = await api.post(`/api/v1/tenant/employees/${props.employeeId}/experiences`, payload)
      const newItem = res.data?.data || { ...payload }
      newItem._saved = true
      newItem._id = newItem.id || ''
      updated.push(newItem)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.saved'), life: 2000 })
    }
    emit('update:items', updated)
    emit('save')
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

async function confirmDeleteExperience() {
  const idx = deleteTargetIdx.value
  if (idx === null || idx === undefined) return
  const item = props.items[idx]
  if (props.employeeId && item._id) {
    deleteLoading.value = true
    deleteError.value = ''
    try {
      await api.delete(`/api/v1/tenant/employees/${props.employeeId}/experiences/${item._id}`)
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
