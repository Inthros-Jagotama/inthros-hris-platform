<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('employee.tab_addresses') }}</h3>
      <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('common.add')" @click="addItem" />
    </div>
    <template v-if="items.length === 0">
      <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-map-marker text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee.no_addresses') }}</p>
      </div>
    </template>
    <div v-for="(item, idx) in items" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3 relative">
      <Button icon="pi pi-trash" severity="danger" text size="small" class="absolute top-2 right-2" @click="removeItem(idx)" v-tooltip.left="t('common.delete')" />
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('employee.address_type')" required :errors="errs?.[idx]?.type">
          <SelectLabel v-model="item.type" :options="addressTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_address_type')" :class="{'p-invalid':errs?.[idx]?.type}" />
        </FormRow>
        <FormRow :label="t('employee.address')" required :errors="errs?.[idx]?.address">
          <TextInput v-model="item.address" maxlength="255" :placeholder="t('employee.address_placeholder')" :class="{'p-invalid':errs?.[idx]?.address}" />
        </FormRow>
        <FormRow :label="t('employee.province')" :errors="errs?.[idx]?.province_id">
          <SelectLabel v-model="item.province_id" :options="provinceOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_province')" filter :class="{'p-invalid':errs?.[idx]?.province_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.regency')" :errors="errs?.[idx]?.regency_id">
          <SelectLabel v-model="item.regency_id" :options="regencyOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_regency')" filter :class="{'p-invalid':errs?.[idx]?.regency_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.district')" :errors="errs?.[idx]?.district_id">
          <SelectLabel v-model="item.district_id" :options="districtOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_district')" filter :class="{'p-invalid':errs?.[idx]?.district_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.village')" :errors="errs?.[idx]?.village_id">
          <SelectLabel v-model="item.village_id" :options="villageOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_village')" filter :class="{'p-invalid':errs?.[idx]?.village_id}" :showClear="true" />
        </FormRow>
        <FormRow :label="t('employee.postal_code')" :errors="errs?.[idx]?.postal_code">
          <TextInput v-model="item.postal_code" maxlength="5" :placeholder="t('employee.postal_code_placeholder')" :class="{'p-invalid':errs?.[idx]?.postal_code}" />
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
  addressTypeOptions: { type: Array, default: () => [] },
  provinceOptions: { type: Array, default: () => [] },
  regencyOptions: { type: Array, default: () => [] },
  districtOptions: { type: Array, default: () => [] },
  villageOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false }
})
const emit = defineEmits(['update:items', 'save'])
function addItem() {
  const next = [...props.items, { type: '', address: '', province_id: '', regency_id: '', district_id: '', village_id: '', postal_code: '' }]
  emit('update:items', next)
}
function removeItem(idx) {
  const next = props.items.filter((_, i) => i !== idx)
  emit('update:items', next)
}
</script>
