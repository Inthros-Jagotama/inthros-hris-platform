<script setup>
import { onMounted, ref } from 'vue';
import Select from 'primevue/select';
import MultiSelect from 'primevue/multiselect';

defineProps({
    modelValue: [String, Number, Object, Array],
    options: { type: Array, default: () => [] },
    optionLabel: { type: String, default: 'label' },
    optionValue: { type: String, default: 'value' },
    optionGroupLabel: { type: String, default: '' },
    optionGroupChildren: { type: String, default: '' },
    placeholder: { type: String, default: 'Pilih opsi...' },
    size: { type: String, default: 'small' },
    fluid: { type: Boolean, default: true },
    filter: { type: Boolean, default: true },
    emptyFilterMessage: { type: String, default: 'Data tidak ditemukan' },
    showClear: { type: Boolean, default: false },
    multiple: { type: Boolean, default: false }
});
defineEmits(['update:modelValue', 'change']);
const selectRef = ref(null);
onMounted(() => {
    const el = selectRef.value?.$el;
    const inputEl = el?.querySelector('input, .p-select-label, .p-multiselect-label');
    if (inputEl && el?.hasAttribute('autofocus')) inputEl.focus();
});
defineExpose({ focus: () => { const el = selectRef.value?.$el?.querySelector('input, .p-select-label, .p-multiselect-label'); if (el) el.focus(); } });
</script>
<template>
    <!-- PrimeVue Select (single) vs MultiSelect (multiple) — Select v4 tidak mendukung prop multiple -->
    <MultiSelect
        v-if="multiple"
        ref="selectRef"
        :modelValue="modelValue"
        :options="options"
        :optionLabel="optionLabel"
        :optionValue="optionValue"
        :optionGroupLabel="optionGroupLabel || undefined"
        :optionGroupChildren="optionGroupChildren || undefined"
        :placeholder="placeholder"
        :size="size"
        :fluid="fluid"
        :filter="filter"
        :showClear="showClear"
        :emptyFilterMessage="emptyFilterMessage"
        display="chip"
        :maxSelectedLabels="3"
        class="w-full"
        @update:modelValue="$emit('update:modelValue', $event)"
        @change="$emit('change', $event)"
    />
    <Select
        v-else
        ref="selectRef"
        :modelValue="modelValue"
        :options="options"
        :optionLabel="optionLabel"
        :optionValue="optionValue"
        :optionGroupLabel="optionGroupLabel || undefined"
        :optionGroupChildren="optionGroupChildren || undefined"
        :placeholder="placeholder"
        :size="size"
        :fluid="fluid"
        :filter="filter"
        :showClear="showClear"
        :emptyFilterMessage="emptyFilterMessage"
        class="w-full"
        @update:modelValue="$emit('update:modelValue', $event)"
        @change="$emit('change', $event)"
    />
</template>
