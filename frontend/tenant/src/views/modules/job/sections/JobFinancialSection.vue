<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('job_management.financials') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.financial_description') }}</p>
    </div>

    <div>
      <!-- Skeleton while loading financial data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-5" />

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <!-- Has Financial Authority — deskripsi kiri, switch kanan -->
        <div class="flex items-center justify-between gap-4 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
          <div class="min-w-0">
            <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ t('job_management.is_authorized') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('job_management.is_authorized_description') }}</p>
          </div>
          <ToggleSwitch v-model="form.is_authorized" />
        </div>

        <!-- Editable financial fields — satu kolom -->
        <div class="space-y-4 pt-4 border-t border-gray-200 dark:border-gray-700">
          <!-- Level Kas hanya muncul jika Memiliki Wewenang Keuangan aktif -->
          <FormRow v-if="form.is_authorized" :label="t('job_management.cash_level')" :errors="errors?.job_management_value_cash_id">
            <SelectLabel
              v-model="form.job_management_value_cash_id"
              :options="cashOptions"
              option-label="label"
              option-value="value"
              :placeholder="t('common.select')"
              :class="{ 'p-invalid': errors?.job_management_value_cash_id }"
              showClear
            />
          </FormRow>
          <FormRow :label="t('job_management.authority_level')" :errors="errors?.job_management_value_authority_id">
            <SelectLabel
              v-model="form.job_management_value_authority_id"
              :options="currentAuthOptions"
              option-label="label"
              option-value="value"
              :placeholder="t('common.select')"
              :class="{ 'p-invalid': errors?.job_management_value_authority_id }"
              showClear
            />
          </FormRow>
          <FormRow :label="t('job_management.impact_level')" :errors="errors?.job_management_value_impact_id">
            <SelectLabel
              v-model="form.job_management_value_impact_id"
              :options="currentImpactOptions"
              option-label="label"
              option-value="value"
              :placeholder="t('common.select')"
              :class="{ 'p-invalid': errors?.job_management_value_impact_id }"
              showClear
            />
          </FormRow>
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
import { ref, computed, watch, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import ToggleSwitch from 'primevue/toggleswitch'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
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
const loading = ref(true)
const errorMsg = ref('')
const errors = ref({})
const existingId = ref('')
const deleteVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')

const form = ref({
  is_authorized: false,
  job_management_value_cash_id: '',
  job_management_value_authority_id: '',
  job_management_value_impact_id: ''
})

const apiBase = '/api/v1/tenant/job-management/financials'

// Dua set opsi per field — dipilih berdasarkan status is_authorized
const cashOptions = ref([])
const authOptions = ref([])            // type=authority
const authUnauthorizedOptions = ref([]) // type=authority_unauthorized
const impactOptions = ref([])           // type=impact
const impactUnauthorizedOptions = ref([]) // type=impact_unauthorized

const currentAuthOptions = computed(() =>
  form.value.is_authorized ? authOptions.value : authUnauthorizedOptions.value
)
const currentImpactOptions = computed(() =>
  form.value.is_authorized ? impactOptions.value : impactUnauthorizedOptions.value
)

// Muat opsi dari job_management_values:
//   is_authorized=true  → cash, authority, impact
//   is_authorized=false → authority_unauthorized, impact_unauthorized (Level Kas disembunyikan)
async function loadOptions() {
  try {
    const [cashRes, authRes, authUnauthRes, impactRes, impactUnauthRes] = await Promise.all([
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'cash', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'authority', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'authority_unauthorized', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'impact', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'impact_unauthorized', per_page: 100 } })
    ])
    cashOptions.value = (cashRes.data?.data || []).map(v => ({ label: v.descriptions, value: v.id }))
    authOptions.value = (authRes.data?.data || []).map(v => ({ label: v.descriptions, value: v.id }))
    authUnauthorizedOptions.value = (authUnauthRes.data?.data || []).map(v => ({ label: v.descriptions, value: v.id }))
    impactOptions.value = (impactRes.data?.data || []).map(v => ({ label: v.descriptions, value: v.id }))
    impactUnauthorizedOptions.value = (impactUnauthRes.data?.data || []).map(v => ({ label: v.descriptions, value: v.id }))
  } catch { /* ignore */ }
}

// Guard: saat loadData mengisi form dari DB, watch tidak boleh mereset field.
let hydrating = false

// Saat status wewenang berubah (oleh user), kosongkan field yang tidak relevan:
//   - Level Kas hanya valid jika is_authorized
//   - Level Kewenangan & Dampak berganti set opsi (authority → authority_unauthorized, dst)
// flush: 'sync' PENTING — dengan flush default ('pre') callback berjalan setelah blok
// loadData selesai (hydrating sudah false) sehingga field hasil DB ikut ter-reset.
watch(() => form.value.is_authorized, (val, oldVal) => {
  if (hydrating || val === oldVal) return
  form.value.job_management_value_cash_id = ''
  form.value.job_management_value_authority_id = ''
  form.value.job_management_value_impact_id = ''
}, { flush: 'sync' })

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
      hydrating = true
      existingId.value = item.id
      form.value.is_authorized = !!item.is_authorized
      form.value.job_management_value_cash_id = item.job_management_value_cash_id || ''
      form.value.job_management_value_authority_id = item.job_management_value_authority_id || ''
      form.value.job_management_value_impact_id = item.job_management_value_impact_id || ''
      hydrating = false
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
    // Level Kas hanya dikirim jika Memiliki Wewenang Keuangan aktif
    const isAuthorized = !!form.value.is_authorized
    const payload = {
      nomenclature: props.orgName || '',
      full_code: props.orgCode || '',
      is_authorized: isAuthorized,
      job_management_value_cash_id: isAuthorized ? (form.value.job_management_value_cash_id || null) : null,
      job_management_value_authority_id: form.value.job_management_value_authority_id || null,
      job_management_value_impact_id: form.value.job_management_value_impact_id || null,
      organization_id: props.orgId
    }

    if (existingId.value) {
      await api.put(`${apiBase}/${existingId.value}`, {
        is_authorized: isAuthorized,
        job_management_value_cash_id: isAuthorized ? (form.value.job_management_value_cash_id || '') : '',
        job_management_value_authority_id: form.value.job_management_value_authority_id || '',
        job_management_value_impact_id: form.value.job_management_value_impact_id || ''
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
    form.value.is_authorized = false
    form.value.job_management_value_cash_id = ''
    form.value.job_management_value_authority_id = ''
    form.value.job_management_value_impact_id = ''
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
    await Promise.all([loadOptions(), loadData()])
  } finally {
    // Skeleton ditutup setelah master options & record data selesai dimuat
    loading.value = false
  }
})
</script>
