<template>
  <div class="flex items-center justify-center h-full">
    <div class="text-center">
      <i class="pi pi-file-edit text-3xl text-emerald-500 mb-3"></i>
      <h2 class="text-lg font-semibold text-gray-700 dark:text-gray-100">{{ t('employee_movement.contracts') }}</h2>
      <p class="text-sm text-gray-400 mt-1">{{ t('employee_movement.contracts_coming_soon') }}</p>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/composables/useI18n'

const { t } = useI18n()
</script>
