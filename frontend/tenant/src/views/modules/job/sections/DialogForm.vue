<template>
  <Dialog :visible="visible" @update:visible="$emit('update:visible', $event)" :header="title" modal :style="{ width: dialogWidth }" class="p-fluid" :closable="!saving">
    <div class="space-y-4">
      <slot />
      <div v-if="Object.keys(errors).length" class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded p-3 text-xs text-red-700 dark:text-red-300">
        <p v-for="(msgs, field) in errors" :key="field" class="mb-1">
          <strong>{{ field }}:</strong> {{ Array.isArray(msgs) ? msgs.join(', ') : msgs }}
        </p>
      </div>
    </div>
    <template #footer>
      <Button :label="t('common.cancel')" size="small" outlined severity="secondary" :disabled="saving" @click="$emit('cancel')" />
      <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" @click="$emit('save')" />
    </template>
  </Dialog>
</template>
<script setup>
import { computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'

const props = defineProps({
  visible: Boolean,
  title: String,
  saving: Boolean,
  errors: { type: Object, default: () => ({}) },
  width: { type: String, default: '480px' }
})
defineEmits(['save', 'cancel'])
const { t } = useI18n()

const dialogWidth = computed(() => props.width === 'maximize' ? '90vw' : props.width)
</script>
