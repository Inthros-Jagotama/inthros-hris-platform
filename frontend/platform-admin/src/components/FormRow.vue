<script setup>
import InputLabel from './InputLabel.vue';

defineProps({
    label: String,
    labelSpan: { type: String, default: 'col-span-3' },
    inputSpan: { type: String, default: 'col-span-9' },
    required: { type: Boolean, default: false },
    errors: { type: String, default: '' }
});
</script>

<template>
    <InputLabel :class="labelSpan">{{ label }} <span v-if="required" class="text-rose-500">*</span> </InputLabel>
    <div :class="inputSpan">
        <slot />
        <small v-if="errors" class="text-red-500 text-xs mt-1 block">{{ errors }}</small>
    </div>
</template>
