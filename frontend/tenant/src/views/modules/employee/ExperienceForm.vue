<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('employee.tab_experience') }}</h3>
      <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('common.add')" @click="addItem" />
    </div>
    <template v-if="items.length === 0">
      <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-briefcase text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee.no_experience') }}</p>
      </div>
    </template>
    <div v-for="(item, idx) in items" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3 relative">
      <Button icon="pi pi-trash" severity="danger" text size="small" class="absolute top-2 right-2" @click="removeItem(idx)" v-tooltip.left="t('common.delete')" />
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('employee.company')" required :errors="errs?.[idx]?.company">
          <TextInput v-model="item.company" maxlength="255" :placeholder="t('employee.company_placeholder')" :class="{'p-invalid':errs?.[idx]?.company}" />
        </FormRow>
        <FormRow :label="t('employee.position')" :errors="errs?.[idx]?.position">
          <TextInput v-model="item.position" maxlength="255" :placeholder="t('employee.position_placeholder')" :class="{'p-invalid':errs?.[idx]?.position}" />
        </FormRow>
        <FormRow :label="t('employee.start_year')" :errors="errs?.[idx]?.start_year">
          <TextInput v-model="item.start_year" maxlength="4" :placeholder="t('employee.year_placeholder')" :class="{'p-invalid':errs?.[idx]?.start_year}" />
        </FormRow>
        <FormRow :label="t('employee.end_year')" :errors="errs?.[idx]?.end_year">
          <TextInput v-model="item.end_year" maxlength="4" :placeholder="t('employee.year_placeholder')" :class="{'p-invalid':errs?.[idx]?.end_year}" />
        </FormRow>
      </div>
      <div v-if="item._saved" class="flex items-center gap-1 text-emerald-500 text-xs">
        <i class="pi pi-check-circle"></i><span>{{ t('employee.saved') }}</span>
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
const { t } = useI18n()
const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})
const emit = defineEmits(['update:items', 'save'])
function addItem() {
  const next = [...props.items, { company: '', position: '', start_year: '', end_year: '' }]
  emit('update:items', next)
}
function removeItem(idx) {
  emit('update:items', props.items.filter((_, i) => i !== idx))
}
</script>
