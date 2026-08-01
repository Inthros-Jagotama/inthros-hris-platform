<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_insurance') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.insurance_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="openDialog()" />
    </div>

    <!-- DataTable for saved items -->
    <DataTable :value="items" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-shield text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee.no_insurance') }}</p>
        </div>
      </template>
      <Column field="category" :header="t('employee.insurance_category')" sortable>
        <template #body="{data}">
          <Tag v-if="data.category" :value="data.category" :severity="getCategorySeverity(data.category)" class="!text-xs !px-2 !py-0.5" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column field="number" :header="t('employee.insurance_number')" sortable>
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.number }}</span>
        </template>
      </Column>
      <Column field="name" :header="t('employee.insurance_name')" sortable>
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="type" :header="t('employee.insurance_type')">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.type || '-' }}</span>
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
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('employee.edit_insurance') : t('employee.new_insurance')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="grid grid-cols-1 gap-3 py-2">
        <FormRow :label="t('employee.insurance_category')" :errors="errors?.category">
          <SelectLabel v-model="form.category" :options="insuranceCategoryOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_insurance_category')" :class="{'p-invalid':errors?.category}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.insurance_number')" required :errors="errors?.number">
          <TextInput v-model="form.number" maxlength="100" :placeholder="t('employee.insurance_number_placeholder')" :class="{'p-invalid':errors?.number}" />
        </FormRow>
        <FormRow :label="t('employee.insurance_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="100" :placeholder="t('employee.insurance_name_placeholder')" :class="{'p-invalid':errors?.name}" />
        </FormRow>
        <FormRow :label="t('employee.insurance_type')" :errors="errors?.type">
          <TextInput v-model="form.type" maxlength="100" :placeholder="t('employee.insurance_type_placeholder')" :class="{'p-invalid':errors?.type}" />
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
      :message="t('employee.confirm_delete_insurance')"
      @confirm="confirmDeleteInsurance"
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
import Tag from 'primevue/tag'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import FormRow from '@/components/FormRow.vue'
import Dialog from 'primevue/dialog'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { getValidationErrors } from '@/services/responseHandler'

const { t } = useI18n()
const toast = useToast()

const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  insuranceCategoryOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false },
  employeeId: { type: String, default: '' }
})
const emit = defineEmits(['update:items', 'save'])

function getCategorySeverity(category) {
  if (!category) return 'info'
  if (category.includes('BPJS')) return 'success'
  return 'warn'
}

// ── Dialog state ──
const dialogVisible = ref(false)
const editingId = ref(null)
const editingIdx = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ category: '', number: '', name: '', type: '' })
const isEditing = computed(() => editingId.value !== null)

function openDialog() {
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
  form.value = { category: '', number: '', name: '', type: '' }
  dialogVisible.value = true
}

function openEditDialog(item, idx) {
  editingId.value = item._id || item.id || null
  editingIdx.value = idx
  errors.value = {}
  form.value = {
    category: item.category || '',
    number: item.number || '',
    name: item.name || '',
    type: item.type || ''
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { category: '', number: '', name: '', type: '' }
  editingId.value = null
  editingIdx.value = null
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  if (!form.value.number?.trim()) { errors.value = { number: [t('form.required')] }; return }
  if (!form.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      category: form.value.category || null,
      number: form.value.number,
      name: form.value.name,
      type: form.value.type || null
    }
    const updated = [...props.items]
    if (isEditing.value && editingIdx.value !== null) {
      const res = await api.put(`/api/v1/tenant/employees/${props.employeeId}/insurances/${editingId.value}`, payload)
      const updatedItem = res.data?.data || { ...payload, _saved: true, _id: editingId.value }
      updatedItem._saved = true
      updatedItem._id = updatedItem.id || editingId.value
      updated[editingIdx.value] = { ...updated[editingIdx.value], ...updatedItem }
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.updated'), life: 2000 })
    } else {
      const res = await api.post(`/api/v1/tenant/employees/${props.employeeId}/insurances`, payload)
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

async function confirmDeleteInsurance() {
  const idx = deleteTargetIdx.value
  if (idx === null || idx === undefined) return
  const item = props.items[idx]
  if (props.employeeId && item._id) {
    deleteLoading.value = true
    deleteError.value = ''
    try {
      await api.delete(`/api/v1/tenant/employees/${props.employeeId}/insurances/${item._id}`)
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