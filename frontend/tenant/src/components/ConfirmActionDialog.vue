<template>
  <Dialog
    v-model:visible="visibleModel"
    :header="title"
    modal
    :closable="!loading"
    :dismissable-mask="!loading"
    :close-on-escape="!loading"
    :style="{ width: '420px' }"
    @hide="onHide"
  >
    <div class="flex items-start gap-3">
      <div class="flex-shrink-0 w-10 h-10 rounded-full flex items-center justify-center" :class="iconWrapClass">
        <i :class="[icon, iconColorClass]" class="text-lg"></i>
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-sm text-gray-700 dark:text-gray-300 leading-relaxed">{{ message }}</p>
        <p v-if="errorMsg" class="mt-2 text-xs text-red-500 bg-red-50 dark:bg-red-900/20 rounded px-2 py-1.5">
          <i class="pi pi-info-circle mr-1"></i>{{ errorMsg }}
        </p>
      </div>
    </div>
    <template #footer>
      <div class="flex items-center justify-end gap-2">
        <Button
          :label="cancelLabel"
          severity="secondary"
          outlined
          size="small"
          :disabled="loading"
          @click="onCancel"
        />
        <Button
          :label="confirmLabel"
          :severity="severity"
          size="small"
          :loading="loading"
          :disabled="loading"
          @click="$emit('confirm')"
        />
      </div>
    </template>
  </Dialog>
</template>

<script setup>
import { computed } from 'vue'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'

const props = defineProps({
  visible: { type: Boolean, default: false },
  title: { type: String, default: 'Confirm' },
  message: { type: String, default: 'Are you sure?' },
  loading: { type: Boolean, default: false },
  errorMsg: { type: String, default: '' },
  cancelLabel: { type: String, default: 'Cancel' },
  confirmLabel: { type: String, default: 'Confirm' },
  icon: { type: String, default: 'pi pi-question-circle' },
  severity: { type: String, default: 'primary' },
})

const emit = defineEmits(['update:visible', 'confirm', 'cancel'])

const severityClasses = {
  primary: { wrap: 'bg-blue-100 dark:bg-blue-900/30', icon: 'text-blue-500' },
  info: { wrap: 'bg-blue-100 dark:bg-blue-900/30', icon: 'text-blue-500' },
  success: { wrap: 'bg-emerald-100 dark:bg-emerald-900/30', icon: 'text-emerald-500' },
  warn: { wrap: 'bg-amber-100 dark:bg-amber-900/30', icon: 'text-amber-500' },
  danger: { wrap: 'bg-red-100 dark:bg-red-900/30', icon: 'text-red-500' },
}

const iconWrapClass = computed(() => (severityClasses[props.severity] || severityClasses.primary).wrap)
const iconColorClass = computed(() => (severityClasses[props.severity] || severityClasses.primary).icon)

const visibleModel = computed({
  get: () => props.visible,
  set: (val) => {
    if (!props.loading) {
      emit('update:visible', val)
    }
  }
})

function onHide() {
  if (props.loading) return
  emit('update:visible', false)
}

function onCancel() {
  emit('cancel')
  emit('update:visible', false)
}
</script>
