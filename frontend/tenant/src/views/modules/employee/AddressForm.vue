<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_addresses') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.address_description') }}</p>
      </div>
      <Button icon="pi pi-plus" size="small" severity="primary" :label="t('common.add')" @click="addItem" />
    </div>
    <template v-if="items.length === 0">
      <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
        <i class="pi pi-map-marker text-3xl mb-2 opacity-50"></i>
        <p class="text-sm">{{ t('employee.no_addresses') }}</p>
      </div>
    </template>
    <div v-for="(item, idx) in items" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('employee.address_type')" required :errors="errs?.[idx]?.type" class="md:col-span-2">
          <div class="flex gap-3">
            <RadioLabel v-for="opt in addressTypeOptions" :key="opt.value" :model-value="item.type" :value="opt.value" :id="`addr-type-${idx}-${opt.value}`" :label="opt.label" :class="{'p-invalid':errs?.[idx]?.type}" @update:model-value="item.type = $event" />
          </div>
        </FormRow>
        <FormRow :label="t('employee.address')" required :errors="errs?.[idx]?.address">
          <TextInput v-model="item.address" maxlength="255" :placeholder="t('employee.address_placeholder')" :class="{'p-invalid':errs?.[idx]?.address}" />
        </FormRow>
        <FormRow :label="t('employee.village')" required :errors="errs?.[idx]?.village_id">
          <AutoComplete
            :model-value="getVillageDisplay(item)"
            :suggestions="villageSuggestions[idx] || []"
            @complete="onVillageSearch($event, idx)"
            @item-select="onVillageSelect($event, idx)"
            optionLabel="label"
            :placeholder="t('employee.search_village')"
            :class="{'p-invalid':errs?.[idx]?.village_id}"
            size="small"
            class="w-full"
            forceSelection
            :dropdown="true"
            :loading="villageLoading[idx] || false"
          >
            <template #option="{option}">
              <div class="flex flex-col py-1">
                <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ option.label }}</span>
                <span class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                  <template v-if="option.district_name">Kec. {{ option.district_name }}</template>
                  <template v-if="option.regency_name"><span v-if="option.district_name">, </span>Kab. {{ option.regency_name }}</template>
                  <template v-if="option.province_name"><span v-if="option.district_name || option.regency_name">, </span>{{ option.province_name }}</template>
                </span>
              </div>
            </template>
          </AutoComplete>
        </FormRow>
        <FormRow :label="t('employee.district')" :errors="errs?.[idx]?.district_id">
          <TextInput :model-value="item._district_name || item.district_id" maxlength="255" :placeholder="t('employee.district_auto')" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
        </FormRow>
        <FormRow :label="t('employee.regency')" :errors="errs?.[idx]?.regency_id">
          <TextInput :model-value="item._regency_name || item.regency_id" maxlength="255" :placeholder="t('employee.regency_auto')" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
        </FormRow>
        <FormRow :label="t('employee.province')" :errors="errs?.[idx]?.province_id">
          <TextInput :model-value="item._province_name || item.province_id" maxlength="255" :placeholder="t('employee.province_auto')" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
        </FormRow>
        <FormRow :label="t('employee.postal_code')" :errors="errs?.[idx]?.postal_code">
          <TextInput v-model="item.postal_code" maxlength="5" :placeholder="t('employee.postal_code_placeholder')" :class="{'p-invalid':errs?.[idx]?.postal_code}" />
        </FormRow>
      </div>
      <div class="flex items-center justify-between pt-1">
        <div v-if="item._saved" class="flex items-center gap-1 text-emerald-500 text-xs">
          <i class="pi pi-check-circle"></i><span>{{ t('employee.saved') }}</span>
        </div>
        <div v-else></div>
        <Button icon="pi pi-trash" severity="danger" outlined size="small" :label="t('common.delete')" @click="onDeleteClick(idx)" />
      </div>
    </div>
    <div v-if="items.length > 0" class="flex justify-end pt-2">
      <Button :label="t('employee.save_step')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="$emit('save')" />
    </div>

    <!-- Delete Confirmation Dialog -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="deleteLoading"
      :error-msg="deleteError"
      :title="t('common.confirm')"
      :message="t('employee.confirm_delete_address')"
      @confirm="confirmDeleteAddress"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>
<script setup>
import { useI18n } from '@/composables/useI18n'
import { ref, reactive } from 'vue'
import Button from 'primevue/button'
import AutoComplete from 'primevue/autocomplete'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import RadioLabel from '@/components/RadioLabel.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import api from '@/services/api'
import { useToast } from 'primevue/usetoast'
const { t } = useI18n()
const toast = useToast()
const props = defineProps({
  items: { type: Array, required: true },
  errs: { type: Array, default: () => [] },
  addressTypeOptions: { type: Array, default: () => [] },
  provinceOptions: { type: Array, default: () => [] },
  regencyOptions: { type: Array, default: () => [] },
  districtOptions: { type: Array, default: () => [] },
  villageOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false },
  employeeId: { type: String, default: '' },
  onSearchVillage: { type: Function, default: null }
})
const emit = defineEmits(['update:items', 'save'])

const villageSuggestions = reactive({})
const villageLoading = reactive({})

function getVillageDisplay(item) {
  return item._villageLabel || ''
}

async function onVillageSearch(event, idx) {
  const query = event.query?.trim() || ''
  if (!query || query.length < 1) { villageSuggestions[idx] = []; return }
  if (!props.onSearchVillage) return
  villageLoading[idx] = true
  try {
    const results = await props.onSearchVillage(query)
    villageSuggestions[idx] = (results || []).map(r => ({
      label: `${r.name}`,
      sublabel: `${r.district_name ? 'Kec. '+r.district_name : ''}${r.regency_name ? ', Kab. '+r.regency_name : ''}${r.province_name ? ', '+r.province_name : ''}`,
      id: r.id,
      district_id: r.district_id,
      district_name: r.district_name || '',
      regency_id: r.regency_id,
      regency_name: r.regency_name || '',
      province_id: r.province_id,
      province_name: r.province_name || ''
    }))
  } catch { villageSuggestions[idx] = [] }
  finally { villageLoading[idx] = false }
}

function onVillageSelect(event, idx) {
  const sel = event.value
  const items = [...props.items]
  items[idx] = {
    ...items[idx],
    village_id: sel.id,
    district_id: sel.district_id || items[idx].district_id,
    regency_id: sel.regency_id || items[idx].regency_id,
    province_id: sel.province_id || items[idx].province_id,
    _district_name: sel.district_name || '',
    _regency_name: sel.regency_name || '',
    _province_name: sel.province_name || '',
    _villageLabel: sel.label
  }
  emit('update:items', items)
}

function addItem() {
  const next = [...props.items, { type: '', address: '', province_id: '', regency_id: '', district_id: '', village_id: '', postal_code: '', _district_name: '', _regency_name: '', _province_name: '' }]
  emit('update:items', next)
}
// ── Delete state ──
const deleteDialogVisible = ref(false)
const deleteLoading = ref(false)
const deleteError = ref('')
const deleteTargetIdx = ref(null)

function onDeleteClick(idx) {
  const item = props.items[idx]
  if (item._saved && props.employeeId) {
    deleteTargetIdx.value = idx
    deleteError.value = ''
    deleteDialogVisible.value = true
  } else {
    removeItem(idx)
  }
}

function removeItem(idx) {
  const next = props.items.filter((_, i) => i !== idx)
  emit('update:items', next)
}

async function confirmDeleteAddress() {
  const idx = deleteTargetIdx.value
  if (idx === null || idx === undefined) return
  const item = props.items[idx]
  if (!props.employeeId) return
  deleteLoading.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/employees/${props.employeeId}/addresses/${item._id}`)
    removeItem(idx)
    deleteDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.deleted'), life: 2000 })
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleteLoading.value = false
  }
}
</script>
