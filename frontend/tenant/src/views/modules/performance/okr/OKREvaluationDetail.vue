<template>
  <div class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
    <i class="pi pi-bullseye text-4xl mb-3 opacity-50"></i>
    <p class="text-sm font-medium">{{ t('common.coming_soon') }}</p>
  </div>
</template>

<script setup>
import { useI18n } from '@/composables/useI18n'

const { t } = useI18n()
</script>