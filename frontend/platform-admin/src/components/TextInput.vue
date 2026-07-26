<script setup>
import { onMounted, ref } from 'vue';
import InputText from 'primevue/inputtext';

defineProps({
    modelValue: [String, Number],
    // Meneruskan semua prop bawaan PrimeVue jika diperlukan
    size: {
        type: String,
        default: 'small'
    }
});

defineEmits(['update:modelValue']);

const inputRef = ref(null);

onMounted(() => {
    // Memastikan elemen input di dalam komponen PrimeVue mendapatkan focus jika ada atribut autofocus
    if (inputRef.value?.$el?.hasAttribute('autofocus')) {
        inputRef.value.$el.focus();
    }
});

// Mengekspos method focus ke komponen induk
defineExpose({
    focus: () => inputRef.value?.$el?.focus()
});
</script>

<template>
    <InputText
        ref="inputRef"
        :value="modelValue"
        :size="size"
        class="w-full"
        @input="$emit('update:modelValue', $event.target.value)"
    />
</template>
