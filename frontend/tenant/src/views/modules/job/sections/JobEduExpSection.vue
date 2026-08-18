<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('job_management.education_experience') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.education_experience_description') }}</p>
    </div>

    <div class="max-w-2xl">
      <!-- Skeleton while loading education & experience data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="6" cols="grid-cols-1" padding="p-5" />

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">

        <!-- ── Group 1: Pendidikan (Education) ── -->
        <div class="pt-1">
          <div class="flex items-center gap-2 mb-3">
            <div class="w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400">
              <i class="pi pi-graduation-cap text-sm"></i>
            </div>
            <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('job_management.group_education') }}</h3>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>
          <div class="space-y-4">
            <FormRow :label="t('job_management.education_level')" :errors="errors?.education_id">
              <SelectLabel
                v-model="form.education_id"
                :options="eduOptions"
                :placeholder="t('job_values.select_education')"
                :class="{ 'p-invalid': errors?.education_id }"
              />
            </FormRow>
            <FormRow :label="t('job_management.education_major')" :errors="errors?.education_major_id">
              <MultiSelect
                v-model="form.education_major_id"
                :options="majorOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                class="w-full"
                size="small"
                filter
                showClear
                display="chip"
                :maxSelectedLabels="2"
                :invalid="!!errors.education_major_id"
              />
            </FormRow>
          </div>
        </div>

        <!-- ── Group 2: Pengalaman (Experience) ── -->
        <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-2 mb-3">
            <div class="w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400">
              <i class="pi pi-briefcase text-sm"></i>
            </div>
            <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('job_management.group_experience') }}</h3>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>
          <div class="space-y-4">
            <FormRow :label="t('job_management.experience_range')" :errors="errors?.experience_id">
              <SelectLabel
                v-model="form.experience_id"
                :options="expOptions"
                :placeholder="t('common.select')"
                :class="{ 'p-invalid': errors?.experience_id }"
              />
            </FormRow>
            <FormRow :label="t('job_management.job_family')" :errors="errors?.job_family_id">
              <MultiSelect
                v-model="form.job_family_id"
                :options="jobFamilyOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                class="w-full"
                size="small"
                filter
                showClear
                display="chip"
                :maxSelectedLabels="2"
                :invalid="!!errors.job_family_id"
              />
            </FormRow>
          </div>
        </div>

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
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import MultiSelect from 'primevue/multiselect'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import SelectLabel from '@/components/SelectLabel.vue'

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
const loading = ref(true)
const deleting = ref(false)
const errorMsg = ref('')
const errors = ref({})
const existingId = ref('')
const deleteVisible = ref(false)
const deleteError = ref('')
const form = ref({
  education_id: '',
  education_major_id: [],
  job_family_id: [],
  experience_id: ''
})

const apiBase = '/api/v1/tenant/job-management/education-experiences'

const expOptions = ref([])
const eduOptions = ref([])
const majorOptions = ref([])
const jobFamilyOptions = ref([])

// Muat master: Pendidikan & Pengalaman dari job_management_values (type=education/experience),
// jurusan & bidang pekerjaan dari module setting (education-majors, job-families)
async function loadMaster() {
  try {
    const [eduRes, expRes, majorRes, jfRes] = await Promise.all([
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'education', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'experience', per_page: 100 } }),
      api.get('/api/v1/tenant/settings/education-majors?per_page=200'),
      api.get('/api/v1/tenant/settings/job-families?per_page=100')
    ])
    // Pendidikan diambil dari tabel job_management_values type=education
    // (level 1-5: Sekolah Menengah Pertama → Strata 3) — id = job_management_values.id
    eduOptions.value = (eduRes.data?.data || []).map(e => ({ label: `Lv.${e.level} — ${e.descriptions}`, value: e.id }))
    // Pengalaman Kerja dari job_management_values type=experience (0-2, 3-5, ... Tahun)
    expOptions.value = (expRes.data?.data || []).map(x => ({ label: `Lv.${x.level} — ${x.descriptions}`, value: x.id }))
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
      form.value.education_major_id = Array.isArray(item.education_major_id) ? item.education_major_id : []
      form.value.job_family_id = Array.isArray(item.job_family_id) ? item.job_family_id : []
      form.value.experience_id = item.experience_id || ''
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
      education_major_id: form.value.education_major_id || [],
      job_family_id: form.value.job_family_id || [],
      experience_id: form.value.experience_id || null,
      organization_id: props.orgId
    }

    if (existingId.value) {
      // Array kosong ([]) agar backend mengosongkan jurusan/bidang pekerjaan saat clear
      await api.put(`${apiBase}/${existingId.value}`, {
        education_id: form.value.education_id || '',
        education_major_id: form.value.education_major_id || [],
        job_family_id: form.value.job_family_id || [],
        experience_id: form.value.experience_id || ''
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
    form.value.education_major_id = []
    form.value.job_family_id = []
    form.value.experience_id = ''
    emit('saved')
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
  } catch (err) {
    deleteError.value = err?.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([loadMaster(), loadData()])
  } finally {
    // Skeleton ditutup setelah master options & record data selesai dimuat
    loading.value = false
  }
})
</script>
