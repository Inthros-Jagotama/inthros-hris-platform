<script setup>
import { ref, onMounted } from 'vue';
import PrimeToggleSwitch from 'primevue/toggleswitch';

const props = defineProps({
    modelValue: {
        type: [Boolean, Number, String],
        default: false
    },
    disabled: {
        type: Boolean,
        default: false
    },
    label: {
        type: String,
        default: ''
    },
    trueValue: {
        type: [Boolean, Number, String],
        default: true
    },
    falseValue: {
        type: [Boolean, Number, String],
        default: false
    }
});

const emit = defineEmits(['update:modelValue']);

const switchRef = ref(null);

// Menggunakan perbandingan non-strict (==) agar '1' (string) dan 1 (number) diperlakukan sama
const toggleValue = () => {
    const nextValue = props.modelValue == props.trueValue ? props.falseValue : props.trueValue;
    emit('update:modelValue', nextValue);
};

onMounted(() => {
    if (switchRef.value?.$el?.hasAttribute('autofocus')) {
        switchRef.value.$el.focus();
    }
});

defineExpose({
    focus: () => switchRef.value?.$el?.focus()
});
</script>

<template>
    <div class="inline-flex items-center gap-3 select-none">
        <PrimeToggleSwitch
            ref="switchRef"
            :model-value="modelValue"
            :disabled="disabled"
            :true-value="trueValue"
            :false-value="falseValue"
            @update:model-value="$emit('update:modelValue', $event)"
        />

        <label
            v-if="label"
            class="text-sm font-medium text-gray-700"
            :class="disabled ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'"
            @click="!disabled && toggleValue()"
        >
            {{ label }}
        </label>
    </div>
</template>
