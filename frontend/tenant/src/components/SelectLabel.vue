<script setup>
import { onMounted, ref } from 'vue';
import Select from 'primevue/select';

defineProps({
    modelValue: [String, Number, Object],
    options: { type: Array, default: () => [] },
    optionLabel: { type: String, default: 'label' },
    optionValue: { type: String, default: 'value' },
    placeholder: { type: String, default: 'Pilih opsi...' },
    size: { type: String, default: 'small' },
    fluid: { type: Boolean, default: true },
    filter: { type: Boolean, default: true },
    emptyFilterMessage: { type: String, default: 'Data tidak ditemukan' },
    showClear: { type: Boolean, default: false }
});
defineEmits(['update:modelValue', 'change']);
const selectRef = ref(null);
onMounted(() => {
    const inputEl = selectRef.value?.$el?.querySelector('input, .p-select-label');
    if (inputEl && selectRef.value?.$el?.hasAttribute('autofocus')) inputEl.focus();
});
defineExpose({ focus: () => { const el = selectRef.value?.$el?.querySelector('input, .p-select-label'); if (el) el.focus(); } });
</script>
<template>
    <Select ref="selectRef" :modelValue="modelValue" :options="options" :optionLabel="optionLabel" :optionValue="optionValue" :placeholder="placeholder" :size="size" :fluid="fluid" :filter="filter" :showClear="showClear" :emptyFilterMessage="emptyFilterMessage" class="w-full" @update:modelValue="$emit('update:modelValue', $event)" @change="$emit('change', $event)" />
</template>
