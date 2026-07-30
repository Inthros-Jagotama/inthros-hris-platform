<script setup>
import RadioButton from 'primevue/radiobutton';
import { computed } from 'vue';

const props = defineProps({
    modelValue: [String, Number, Boolean],
    value: [String, Number, Boolean],
    id: {
        type: String,
        required: true
    },
    label: String
});

const emit = defineEmits(['update:modelValue']);

const isSelected = computed(() => props.modelValue === props.value);
</script>

<template>
    <div
        class="flex items-center gap-3 px-3 py-1.5 border rounded-lg cursor-pointer select-none transition-all duration-150"
        :class="isSelected
            ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
            : 'border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 hover:border-gray-300 dark:hover:border-gray-500 hover:shadow-sm'"
        @click="emit('update:modelValue', value)"
    >
        <RadioButton
            :modelValue="modelValue"
            :inputId="id"
            :value="value"
            @update:modelValue="emit('update:modelValue', $event)"
        />
        <label
            :for="id"
            class="text-sm font-medium cursor-pointer select-none"
            :class="isSelected
                ? 'text-emerald-700 dark:text-emerald-300'
                : 'text-surface-700 dark:text-surface-0/80'"
        >
            {{ label }}
        </label>
    </div>
</template>
