<template>
  <div class="max-w-2xl">
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="4" />

    <div v-else class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4 space-y-4">
      <div>
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100">{{ t('employee_id_format.title') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee_id_format.description') }}</p>
      </div>

      <FormRow :label="t('employee_id_format.mode')" required>
        <SelectButton
          v-model="form.generation_mode"
          :options="MODE_OPTIONS"
          optionValue="value"
          optionLabel="label"
          class="!w-full"
          :allowEmpty="false"
          @update:modelValue="refreshPreview"
        />
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-1.5">{{ modeDescription }}</p>
      </FormRow>

      <template v-if="form.generation_mode !== 'MANUAL'">
        <FormRow :label="t('numbering_settings.format_template')" required :errors="errors?.format_template">
          <TextInput
            v-model="form.format_template"
            :placeholder="t('numbering_settings.format_template')"
            :class="{ 'p-invalid': errors?.format_template }"
            @update:modelValue="debouncedPreview"
          />
        </FormRow>

        <FormRow :label="t('numbering_settings.reset_period')" :errors="errors?.reset_period">
          <Select
            v-model="form.reset_period"
            :options="RESET_PERIOD_OPTIONS"
            optionValue="value"
            optionLabel="label"
            class="!w-full"
            :class="{ 'p-invalid': errors?.reset_period }"
            @update:modelValue="refreshPreview"
          />
        </FormRow>

        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('numbering_settings.tokens_help') }}</p>

        <div class="rounded-md bg-gray-50 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 px-3 py-2">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('numbering_settings.preview_label') }}</p>
          <p class="text-sm font-mono font-semibold text-navy-800 dark:text-gray-100">{{ preview || '—' }}</p>
        </div>
      </template>

      <div class="flex justify-end">
        <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="save" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import Select from 'primevue/select'
import SelectButton from 'primevue/selectbutton'
import Button from 'primevue/button'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

const MODE_OPTIONS = computed(() => [
  { value: 'MANUAL', label: t('employee_id_format.mode_manual') },
  { value: 'HYBRID', label: t('employee_id_format.mode_hybrid') },
  { value: 'AUTO', label: t('employee_id_format.mode_auto') }
])

const RESET_PERIOD_OPTIONS = [
  { value: 'yearly', label: t('numbering_settings.reset_yearly') },
  { value: 'monthly', label: t('numbering_settings.reset_monthly') },
  { value: 'never', label: t('numbering_settings.reset_never') }
]

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-16' }
]

const loading = ref(true)
const saving = ref(false)
const errors = ref({})
const preview = ref('')
const form = reactive({ generation_mode: 'MANUAL', format_template: 'EMP-{year}-{sequence:4}', reset_period: 'yearly' })

const modeDescription = computed(() => {
  switch (form.generation_mode) {
    case 'AUTO': return t('employee_id_format.mode_auto_desc')
    case 'HYBRID': return t('employee_id_format.mode_hybrid_desc')
    default: return t('employee_id_format.mode_manual_desc')
  }
})

let previewTimer = null
function debouncedPreview() {
  clearTimeout(previewTimer)
  previewTimer = setTimeout(refreshPreview, 400)
}

async function loadSetting() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/employee-id-format')
    const data = res.data?.data
    if (data) {
      form.generation_mode = data.generation_mode
      form.format_template = data.format_template
      form.reset_period = data.reset_period
    }
    await refreshPreview()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function refreshPreview() {
  if (form.generation_mode === 'MANUAL') { preview.value = ''; return }
  try {
    const res = await api.get('/api/v1/tenant/settings/employee-id-format/preview')
    preview.value = res.data?.data?.preview || ''
  } catch {
    preview.value = ''
  }
}

async function save() {
  saving.value = true
  errors.value = {}
  try {
    await api.put('/api/v1/tenant/settings/employee-id-format', {
      generation_mode: form.generation_mode,
      format_template: form.format_template,
      reset_period: form.reset_period
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('common.saved'), life: 3000 })
    await refreshPreview()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadSetting)
</script>
