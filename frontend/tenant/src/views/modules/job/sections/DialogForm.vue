<template>
  <Dialog :visible="visible" @update:visible="$emit('update:visible', $event)" :header="title" modal :style="{ width: '480px' }" class="p-fluid" :closable="!saving">
    <div class="space-y-4">
      <slot />
      <div v-if="Object.keys(errors).length" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300">
        <p v-for="(msgs, field) in errors" :key="field" class="mb-1">
          <strong>{{ field }}:</strong> {{ Array.isArray(msgs) ? msgs.join(', ') : msgs }}
        </p>
      </div>
    </div>
    <template #footer>
      <Button :label="t('common.cancel')" size="small" text :disabled="saving" @click="$emit('cancel')" />
      <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" @click="$emit('save')" />
    </template>
  </Dialog>
</template>
<script setup>
import { useI18n } from '@/composables/useI18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
defineProps({ visible: Boolean, title: String, saving: Boolean, errors: { type: Object, default: () => ({}) } })
defineEmits(['save', 'cancel'])
const { t } = useI18n()
</script>
