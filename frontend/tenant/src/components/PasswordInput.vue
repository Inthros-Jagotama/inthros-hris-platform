<script setup>
import { ref, onMounted } from 'vue'
import Password from 'primevue/password'

defineProps({
  modelValue: [String, Number],
  size: {
    type: String,
    default: 'small'
  },
  placeholder: {
    type: String,
    default: ''
  },
  feedback: {
    type: Boolean,
    default: false // Set true jika ingin menampilkan indikator kekuatan password
  }
})

defineEmits(['update:modelValue'])

const inputRef = ref(null)

onMounted(() => {
  // PrimeVue Password membungkus input asli di dalam $el
  const nativeInput = inputRef.value?.$el?.querySelector('input')
  if (nativeInput?.hasAttribute('autofocus')) {
    nativeInput.focus()
  }
})

// Focus expose untuk elemen input di dalamnya
defineExpose({
  focus: () => {
    const nativeInput = inputRef.value?.$el?.querySelector('input')
    nativeInput?.focus()
  }
})
</script>

<template>
  <Password
    ref="inputRef"
    :model-value="modelValue"
    :size="size"
    :placeholder="placeholder"
    :feedback="feedback"
    toggle-mask
    class="w-full"
    input-class="w-full"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>