<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('job_management.objectives') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.objective_description') }}</p>
    </div>

    <div class="max-w-2xl">
      <!-- Skeleton while loading objective data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="3" cols="grid-cols-1" padding="p-5" />

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">

        <!-- Editable objective -->
        <FormRow :label="t('job_management.objective')">
          <Textarea
            v-model="form.objective"
            rows="3"
            class="w-full"
            :class="{ 'p-invalid': errors.objective }"
            :placeholder="t('job_management.objective') + '...'"
          />
        </FormRow>

        <!-- Error display -->
        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <!-- Actions -->
        <div class="flex justify-end gap-2 pt-2">
          <Button
            v-if="existingId"
            :label="t('common.delete')"
            icon="pi pi-trash"
            severity="danger"
            size="small"
            outlined
            @click="deleteVisible = true"
          />
          <Button
            :label="existingId ? t('common.update') : t('common.save')"
            icon="pi pi-check"
            size="small"
            :loading="saving"
            :disabled="saving"
            @click="handleSave"
          />
        </div>
      </div>
    </div>

    <ConfirmDeleteDialog v-model:visible="deleteVisible" :loading="deleting" :error-msg="deleteError" @confirm="handleDelete" @cancel="deleteVisible=false" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import Textarea from 'primevue/textarea'
import Button from 'primevue/button'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import api from '@/services/api'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  orgName: { type: String, default: '' },
  orgCode: { type: String, default: '' }
})

const { t } = useI18n()
const toast = useToast()

const saving = ref(false)
const loading = ref(true)
const deleting = ref(false)
const errorMsg = ref('')
const errors = ref({})
const existingId = ref('')
const deleteVisible = ref(false)
const deleteError = ref('')

const form = ref({
  objective: ''
})

const apiBase = '/api/v1/tenant/job-management/objectives'

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
      form.value.objective = item.objective || ''
    }
  } catch {
    // No existing record
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  errorMsg.value = ''
  errors.value = {}

  saving.value = true
  try {
    const payload = {
      nomenclature: props.orgName || '',
      full_code: props.orgCode || '',
      objective: form.value.objective || '',
      organization_id: props.orgId
    }

    if (existingId.value) {
      await api.put(`${apiBase}/${existingId.value}`, { objective: form.value.objective || '' })
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

async function handleDelete() {
  if (!existingId.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`${apiBase}/${existingId.value}`)
    deleteVisible.value = false
    existingId.value = ''
    form.value.objective = ''
    emit('saved')
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
  } catch (err) {
    deleteError.value = err?.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
