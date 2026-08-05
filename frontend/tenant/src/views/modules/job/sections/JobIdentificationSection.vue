<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.identifications') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.identification_description') }}</p>
    </div>

    <div class="max-w-2xl">
      <!-- Skeleton while loading identification data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-5" />

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <FormRow :label="t('organization.job_family')">
          <TextInput :model-value="jobFamilyLabel" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
        </FormRow>

        <!-- Editable grading at the bottom -->
        <FormRow :label="t('organization.grading')">
          <Select
            v-model="form.grading_id"
            :options="gradingOptions"
            option-label="label"
            option-value="value"
            :placeholder="t('organization.select_grading')"
            class="w-full"
            size="small"
            :invalid="!!errors.grading_id"
          />
        </FormRow>

        <!-- Error display -->
        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>
        <!-- Save button -->
        <div class="flex justify-end pt-2">
          <Button
            :label="t('common.save')"
            icon="pi pi-check"
            size="small"
            :loading="saving"
            :disabled="!form.grading_id"
            @click="handleSave"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import Button from 'primevue/button'
import Select from 'primevue/select'
import SkeletonCard from '@/components/SkeletonCard.vue'
import api from '@/services/api'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  orgName: { type: String, default: '' },
  orgCode: { type: String, default: '' },
  orgGradingId: { type: String, default: '' },
  orgJobFamilyId: { type: String, default: '' },
  gradingOptions: { type: Array, default: () => [] },
  jobFamilyOptions: { type: Array, default: () => [] }
})

const { t } = useI18n()
const toast = useToast()

const saving = ref(false)
const loading = ref(true)
const errorMsg = ref('')
const errors = ref({})
const existingId = ref('')

const form = ref({
  grading_id: ''
})

const apiBase = '/api/v1/tenant/job-management/identifications'

const jobFamilyLabel = computed(() => {
  const found = props.jobFamilyOptions.find(j => j.value === props.orgJobFamilyId)
  return found ? found.label : (props.orgJobFamilyId || '-')
})

function getValidationErrors(err) {
  const fields = err?.response?.data?.error?.fields
  if (fields && typeof fields === 'object') {
    const map = {}
    for (const [key, msgs] of Object.entries(fields)) {
      map[key] = Array.isArray(msgs) ? msgs[0] : msgs
    }
    return map
  }
  return {}
}

async function loadData() {
  if (!props.orgId) {
    loading.value = false
    return
  }
  try {
    const res = await api.get(apiBase, { params: { organization_id: props.orgId, per_page: 1 } })
    const list = res.data?.data || []
    if (list.length > 0) {
      const item = list[0]
      existingId.value = item.id
      form.value.grading_id = item.grading_id || props.orgGradingId || ''
    } else {
      // New record: default grading from org data
      form.value.grading_id = props.orgGradingId || ''
    }
  } catch {
    // No existing record
    form.value.grading_id = props.orgGradingId || ''
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  errorMsg.value = ''
  errors.value = {}

  if (!form.value.grading_id) {
    errorMsg.value = t('job_management.grading_required')
    return
  }

  saving.value = true
  try {
    const payload = {
      nomenclature: props.orgName || '',
      full_code: props.orgCode || '',
      grading_id: form.value.grading_id,
      organization_id: props.orgId
    }

    if (existingId.value) {
      await api.put(`${apiBase}/${existingId.value}`, { grading_id: form.value.grading_id })
    } else {
      const res = await api.post(apiBase, payload)
      existingId.value = res.data?.data?.id || ''
    }

    toast.add({ severity: 'success', summary: t('message.success'), detail: t('common.saved'), life: 2000 })
    emit('saved')
  } catch (err) {
    const ve = getValidationErrors(err)
    if (Object.keys(ve).length > 0) {
      errors.value = ve
      errorMsg.value = Object.values(ve).join(', ')
    } else {
      errorMsg.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
