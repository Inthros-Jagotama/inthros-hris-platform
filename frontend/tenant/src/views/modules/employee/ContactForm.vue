<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_contacts') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.contact_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="addItem" />
    </div>
    <template v-if="items.length === 0">
      <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-phone text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee.no_contacts') }}</p>
      </div>
    </template>
    <div v-for="(item, idx) in items" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3 relative">
      <Button icon="pi pi-trash" severity="danger" text size="small" class="absolute top-2 right-2" @click="removeItem(idx)" v-tooltip.left="t('common.delete')" />
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('employee.contact_name')" required :errors="errs?.[idx]?.name">
          <TextInput v-model="item.name" maxlength="255" :placeholder="t('employee.contact_name_placeholder')" :class="{'p-invalid':errs?.[idx]?.name}" />
        </FormRow>
        <FormRow :label="t('employee.contact_phone')" required :errors="errs?.[idx]?.phone_number">
          <TextInput v-model="item.phone_number" maxlength="50" :placeholder="t('employee.phone_placeholder')" :class="{'p-invalid':errs?.[idx]?.phone_number}" />
        </FormRow>
        <FormRow :label="t('employee.relationship_type')" :errors="errs?.[idx]?.relationship_type_id">
          <SelectLabel v-model="item.relationship_type_id" :options="relationshipTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_relationship')" :class="{'p-invalid':errs?.[idx]?.relationship_type_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.address')" :errors="errs?.[idx]?.address">
          <TextInput v-model="item.address" maxlength="255" :placeholder="t('employee.address_placeholder')" :class="{'p-invalid':errs?.[idx]?.address}" />
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
import SelectLabel from '@/components/SelectLabel.vue'
const { t } = useI18n()
const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  relationshipTypeOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})
const emit = defineEmits(['update:items', 'save'])
function addItem() {
  const next = [...props.items, { name: '', relationship_type_id: '', phone_number: '', address: '' }]
  emit('update:items', next)
}
function removeItem(idx) {
  emit('update:items', props.items.filter((_, i) => i !== idx))
}
</script>
