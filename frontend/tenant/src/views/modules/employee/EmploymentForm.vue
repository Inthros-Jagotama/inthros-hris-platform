<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_employment') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.employment_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="addItem" />
    </div>
    <template v-if="items.length === 0">
      <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-briefcase text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee.no_employment') }}</p>
      </div>
    </template>
    <div v-for="(item, idx) in items" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3 relative">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('employee.organization')" :errors="errs?.[idx]?.organization_id">
          <SelectLabel v-model="item.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_organization')" filter :class="{'p-invalid':errs?.[idx]?.organization_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.employment_status')" :errors="errs?.[idx]?.employment_status_id">
          <SelectLabel v-model="item.employment_status_id" :options="employmentStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_employment_status')" :class="{'p-invalid':errs?.[idx]?.employment_status_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.decision_letter')" required :errors="errs?.[idx]?.decision_letter_number">
          <TextInput v-model="item.decision_letter_number" maxlength="50" :placeholder="t('employee.decision_letter_placeholder')" :class="{'p-invalid':errs?.[idx]?.decision_letter_number}" />
        </FormRow>
        <FormRow :label="t('employee.decision_letter_date')" required :errors="errs?.[idx]?.decision_letter_date">
          <DateInput v-model="item.decision_letter_date" :placeholder="t('employee.date_placeholder')" :class="{'p-invalid':errs?.[idx]?.decision_letter_date}" />
        </FormRow>
        <FormRow :label="t('employee.effective_date')" required :errors="errs?.[idx]?.effective_date">
          <DateInput v-model="item.effective_date" :placeholder="t('employee.date_placeholder')" :class="{'p-invalid':errs?.[idx]?.effective_date}" />
        </FormRow>
        <FormRow :label="t('employee.effective_end_date')" :errors="errs?.[idx]?.effective_end_date">
          <DateInput v-model="item.effective_end_date" :placeholder="t('employee.date_placeholder')" :class="{'p-invalid':errs?.[idx]?.effective_end_date}" />
        </FormRow>
      </div>
      <div v-if="item._saved" class="justify-between flex items-center mt-2">
        <div class="flex items-center gap-2 text-sm text-green-600 dark:text-green-400">
          <i class="pi pi-check-circle"></i><span>{{ t('employee.saved') }}</span>
        </div>
        <Button icon="pi pi-trash" severity="danger" size="small" class="absolute top-2 right-2" @click="removeItem(idx)" v-tooltip.left="t('common.delete')" :label="t('common.delete')" />
      </div>
    </div>
    <div v-if="items.length > 0" class="flex justify-end pt-2">
      <Button :label="t('employee.save_step')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="$emit('save')" />
    </div>
  </div>
</template>
<script setup>
import { useI18n } from '@/composables/useI18n'
import Button from 'primevue/button'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
const { t } = useI18n()
const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  organizationOptions: { type: Array, default: () => [] },
  employmentStatusOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})
const emit = defineEmits(['update:items', 'save'])
function addItem() {
  const next = [...props.items, { organization_id: '', employment_status_id: '', decision_letter_number: '', decision_letter_date: '', effective_date: '', effective_end_date: '' }]
  emit('update:items', next)
}
function removeItem(idx) {
  emit('update:items', props.items.filter((_, i) => i !== idx))
}
</script>
