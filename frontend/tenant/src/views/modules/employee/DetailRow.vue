<script setup>
defineProps({
  label: { type: String, default: '' },
  value: { type: [String, Number], default: '' }
})
</script>

<template>
  <div class="min-w-0">
    <dt class="text-xs text-gray-400 dark:text-gray-500">{{ label }}</dt>
    <dd class="text-sm text-navy-800 dark:text-gray-100 mt-0.5 truncate">{{ value || '-' }}</dd>
  </div>
</template>
