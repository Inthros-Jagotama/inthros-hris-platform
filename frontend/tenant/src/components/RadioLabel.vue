<script setup>
import RadioButton from 'primevue/radiobutton';

defineProps({
    modelValue: [String, Number, Boolean],
    value: [String, Number, Boolean],
    id: {
        type: String,
        required: true
    },
    label: String
});

defineEmits(['update:modelValue']);
</script>

<template>
    <div class="flex items-center gap-2">
        <RadioButton
            :modelValue="modelValue"
            :inputId="id"
            :value="value"
            @update:modelValue="$emit('update:modelValue', $event)"
        />
        <label
            :for="id"
            class="text-sm cursor-pointer select-none text-surface-700 dark:text-surface-0/80"
        >
            {{ label }}
        </label>
    </div>
</template>
