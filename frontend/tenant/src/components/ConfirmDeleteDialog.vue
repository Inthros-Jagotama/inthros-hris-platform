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
      <div class="flex-shrink-0 w-10 h-10 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center">
        <i class="pi pi-exclamation-triangle text-red-500 text-lg"></i>
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
          severity="danger"
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
  confirmLabel: { type: String, default: 'Delete' },
})

const emit = defineEmits(['update:visible', 'confirm', 'cancel'])

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
