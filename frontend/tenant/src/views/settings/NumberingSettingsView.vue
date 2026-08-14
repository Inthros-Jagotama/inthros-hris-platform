<template>
  <div class="space-y-4">
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="4" />

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div
        v-for="dt in DOCUMENT_TYPES"
        :key="dt.key"
        class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4 space-y-4"
      >
        <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t(dt.labelKey) }}</h3>

        <FormRow :label="t('numbering_settings.format_template')" required :errors="errors[dt.key]?.format_template">
          <TextInput
            v-model="forms[dt.key].format_template"
            :placeholder="t('numbering_settings.format_template')"
            :class="{ 'p-invalid': errors[dt.key]?.format_template }"
            @update:modelValue="debouncedPreview(dt.key)"
          />
        </FormRow>

        <FormRow :label="t('numbering_settings.reset_period')" :errors="errors[dt.key]?.reset_period">
          <Select
            v-model="forms[dt.key].reset_period"
            :options="RESET_PERIOD_OPTIONS"
            optionValue="value"
            optionLabel="label"
            class="!w-full"
            :class="{ 'p-invalid': errors[dt.key]?.reset_period }"
            @update:modelValue="refreshPreview(dt.key)"
          />
        </FormRow>

        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('numbering_settings.tokens_help') }}</p>

        <div class="rounded-md bg-gray-50 dark:bg-gray-900/40 border border-gray-200 dark:border-gray-700 px-3 py-2">
          <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('numbering_settings.preview_label') }}</p>
          <p class="text-sm font-mono font-semibold text-gray-800 dark:text-gray-100">
            {{ previews[dt.key] || '—' }}
          </p>
        </div>

        <div class="flex justify-end">
          <Button
            :label="t('common.save')"
            size="small"
            :loading="saving[dt.key]"
            :disabled="saving[dt.key]"
            @click="save(dt.key)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Select from 'primevue/select'
import Button from 'primevue/button'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

const DOCUMENT_TYPES = [
  { key: 'employee_movement', labelKey: 'numbering_settings.employee_movement' },
  { key: 'employee_contract', labelKey: 'numbering_settings.employee_contract' },
]

const RESET_PERIOD_OPTIONS = [
  { value: 'yearly', label: t('numbering_settings.reset_yearly') },
  { value: 'monthly', label: t('numbering_settings.reset_monthly') },
  { value: 'never', label: t('numbering_settings.reset_never') },
]

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-16' },
]

const loading = ref(true)
const forms = reactive({})
const previews = reactive({})
const saving = reactive({})
const errors = reactive({})

let previewTimers = {}

function debouncedPreview(documentType) {
  clearTimeout(previewTimers[documentType])
  previewTimers[documentType] = setTimeout(() => refreshPreview(documentType), 400)
}

async function loadSettings() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/document-numbering')
    const items = res.data?.data || []
    for (const item of items) {
      forms[item.document_type] = {
        format_template: item.format_template,
        reset_period: item.reset_period,
      }
    }
    await Promise.all(DOCUMENT_TYPES.map((dt) => refreshPreview(dt.key)))
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  } finally {
    loading.value = false
  }
}

async function refreshPreview(documentType) {
  try {
    const res = await api.get(`/api/v1/tenant/settings/document-numbering/${documentType}/preview`)
    previews[documentType] = res.data?.data?.preview || ''
  } catch {
    previews[documentType] = ''
  }
}

async function save(documentType) {
  saving[documentType] = true
  errors[documentType] = {}
  try {
    const payload = forms[documentType]
    await api.put(`/api/v1/tenant/settings/document-numbering/${documentType}`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('common.saved'), life: 3000 })
    await refreshPreview(documentType)
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors[documentType] = fe
    } else {
      toast.add({
        severity: 'error',
        summary: t('message.error'),
        detail: e.response?.data?.error?.message || t('message.operation_failed'),
        life: 4000,
      })
    }
  } finally {
    saving[documentType] = false
  }
}

onMounted(loadSettings)
</script>
