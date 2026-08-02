<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.education_experience') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.education_experience_description') }}</p>
    </div>

    <div class="max-w-2xl space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <!-- Read-only from org data -->
      <FormRow :label="t('organization.nomenclature')">
        <TextInput :model-value="orgName" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
      </FormRow>
      <FormRow :label="t('organization.full_code')">
        <TextInput :model-value="orgCode" disabled class="!bg-gray-50 dark:!bg-gray-700 !cursor-not-allowed" />
      </FormRow>

      <!-- Editable selects (relasi ke master module setting) -->
      <FormRow :label="t('job_management.education_level')" :errors="errors?.education_id">
        <Select
          v-model="form.education_id"
          :options="eduOptions"
          option-label="label"
          option-value="value"
          :placeholder="t('common.select')"
          class="w-full"
          size="small"
          showClear
          :invalid="!!errors.education_id"
        />
      </FormRow>
      <FormRow :label="t('job_management.education_major')" :errors="errors?.education_major_id">
        <Select
          v-model="form.education_major_id"
          :options="majorOptions"
          option-label="label"
          option-value="value"
          :placeholder="t('common.select')"
          class="w-full"
          size="small"
          showClear
          :invalid="!!errors.education_major_id"
        />
      </FormRow>
      <FormRow :label="t('job_management.job_family')" :errors="errors?.job_family_id">
        <Select
          v-model="form.job_family_id"
          :options="jobFamilyOptions"
          option-label="label"
          option-value="value"
          :placeholder="t('common.select')"
          class="w-full"
          size="small"
          showClear
          :invalid="!!errors.job_family_id"
        />
      </FormRow>
      <FormRow :label="t('job_management.experience_range')" :errors="errors?.experience_range">
        <Select
          v-model="form.experience_range"
          :options="expOptions"
          option-label="label"
          option-value="value"
          :placeholder="t('common.select')"
          class="w-full"
          size="small"
          showClear
          :invalid="!!errors.experience_range"
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

    <ConfirmDeleteDialog v-model:visible="deleteVisible" :loading="deleting" :error-msg="deleteError" @confirm="handleDelete" @cancel="deleteVisible=false" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import Select from 'primevue/select'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  orgName: { type: String, default: '' },
  orgCode: { type: String, default: '' },
  // Dideklarasikan agar tidak jadi fallthrough attr (parent masih pass untuk section lain)
  jobValueMap: { type: Object, default: () => ({}) }
})

const { t } = useI18n()
const toast = useToast()

const saving = ref(false)
const deleting = ref(false)
const errorMsg = ref('')
const errors = ref({})
const existingId = ref('')
const deleteVisible = ref(false)
const deleteError = ref('')
const form = ref({
  education_id: '',
  education_major_id: '',
  job_family_id: '',
  experience_range: ''
})

const apiBase = '/api/v1/tenant/job-management/education-experiences'

// Pengalaman Kerja — hardcoded dropdown (mengikuti seeder Laravel)
const expOptions = [
  { label: '0-2 Tahun', value: '0-2 Tahun' },
  { label: '3-5 Tahun', value: '3-5 Tahun' },
  { label: '6-8 Tahun', value: '6-8 Tahun' },
  { label: '9-11 Tahun', value: '9-11 Tahun' },
  { label: '> 12 Tahun', value: '> 12 Tahun' }
]

const eduOptions = ref([])
const majorOptions = ref([])
const jobFamilyOptions = ref([])

// Muat master dari module setting (educations, education-majors, job-families)
async function loadMaster() {
  try {
    const [eduRes, majorRes, jfRes] = await Promise.all([
      api.get('/api/v1/tenant/settings/educations?per_page=100'),
      api.get('/api/v1/tenant/settings/education-majors?per_page=200'),
      api.get('/api/v1/tenant/settings/job-families?per_page=100')
    ])
    eduOptions.value = (eduRes.data?.data || []).map(e => ({ label: `${e.code} - ${e.name}`, value: e.id }))
    majorOptions.value = (majorRes.data?.data || []).map(m => ({ label: `${m.code} - ${m.name}`, value: m.id }))
    jobFamilyOptions.value = (jfRes.data?.data || []).map(j => ({ label: `${j.code} - ${j.name}`, value: j.id }))
  } catch { /* ignore */ }
}

async function loadData() {
  if (!props.orgId) return
  try {
    const res = await api.get(apiBase, { params: { organization_id: props.orgId, per_page: 1 } })
    const list = res.data?.data || []
    if (list.length > 0) {
      const item = list[0]
      existingId.value = item.id
      form.value.education_id = item.education_id || ''
      form.value.education_major_id = item.education_major_id || ''
      form.value.job_family_id = item.job_family_id || ''
      form.value.experience_range = item.experience_range || ''
    }
  } catch {
    // No existing record
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
      education_id: form.value.education_id || null,
      education_major_id: form.value.education_major_id || null,
      job_family_id: form.value.job_family_id || null,
      experience_range: form.value.experience_range || null,
      organization_id: props.orgId
    }

    if (existingId.value) {
      // Kirim string kosong ('') agar backend bisa mengosongkan field saat clear
      await api.put(`${apiBase}/${existingId.value}`, {
        education_id: form.value.education_id || '',
        education_major_id: form.value.education_major_id || '',
        job_family_id: form.value.job_family_id || '',
        experience_range: form.value.experience_range || ''
      })
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
    form.value.education_id = ''
    form.value.education_major_id = ''
    form.value.job_family_id = ''
    form.value.experience_range = ''
    emit('saved')
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
  } catch (err) {
    deleteError.value = err?.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadMaster(), loadData()])
})
</script>
