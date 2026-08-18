<template>
  <div class="min-h-[60vh] flex items-center justify-center p-6">
    <div class="text-center max-w-md">
      <div class="w-16 h-16 mx-auto rounded-full bg-amber-50 dark:bg-amber-500/10 border border-amber-200 dark:border-amber-500/30 flex items-center justify-center mb-4">
        <i class="pi pi-lock text-2xl text-amber-500"></i>
      </div>
      <h1 class="text-xl font-bold text-navy-800 dark:text-gray-100">403</h1>
      <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">{{ t('common.forbidden_title') }}</p>
      <p class="mt-1 text-sm text-gray-400 dark:text-gray-500">{{ t('common.forbidden_desc') }}</p>
      <Button
        :label="t('common.back_to_dashboard')"
        icon="pi pi-home"
        size="small"
        class="mt-6"
        @click="router.push('/dashboard')"
      />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'

const router = useRouter()
const { t } = useI18n()
</script>
