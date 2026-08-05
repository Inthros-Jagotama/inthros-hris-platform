<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.relationships') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.relationship_description') }}</p>
    </div>

    <div class="max-w-2xl">
      <!-- Skeleton while loading relationship data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-5" />

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">

        <!-- ── Group: Ruang Lingkup (Scope) — relation type & frequency ── -->
        <div class="pt-1">
          <div class="flex items-center gap-2 mb-3">
            <div class="w-8 h-8 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400">
              <i class="pi pi-compass text-sm"></i>
            </div>
            <h3 class="text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('job_management.relationship_group_scope') }}</h3>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>
          <div class="space-y-4">
            <FormRow :label="t('job_management.relationship_type')" :errors="errors?.job_management_value_relationship_id">
              <SelectLabel
                v-model="form.job_management_value_relationship_id"
                :options="relOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                :class="{ 'p-invalid': errors?.job_management_value_relationship_id }"
                showClear
              />
            </FormRow>
            <FormRow :label="t('job_management.frequency')" :errors="errors?.job_management_value_frequency_id">
              <SelectLabel
                v-model="form.job_management_value_frequency_id"
                :options="freqOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                :class="{ 'p-invalid': errors?.job_management_value_frequency_id }"
                showClear
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

    <!-- =============================================================
         Work Relations — table view
         Kolom 1 (Organization): select organisasi se-summary (tanpa diri sendiri)
         Kolom 2 (Activity In Connection): inputtext aktivitas
         ============================================================= -->
    <div class="max-w-3xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div>
          <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.relationship_details') }}</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.relationship_details_description') }}</p>
        </div>
        <Button
          :label="t('job_management.add_relationship_detail')"
          icon="pi pi-plus"
          size="small"
          outlined
          :disabled="!existingId || savingDetails"
          @click="addDetailRow"
        />
      </div>

      <!-- Hint jika relationship utama belum disimpan -->
      <div
        v-if="!existingId"
        class="text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2"
      >
        {{ t('job_management.save_relationship_first') }}
      </div>

      <!-- Empty state -->
      <div
        v-else-if="detailRows.length === 0"
        class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
      >
        {{ t('job_management.no_relationship_details') }}
      </div>

      <!-- Detail rows — tampilan tabel -->
      <div v-if="detailRows.length > 0" class="overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700/40 text-left">
              <th class="px-3 py-2 w-10 font-semibold text-gray-600 dark:text-gray-300">#</th>
              <th class="px-3 py-2 font-semibold text-gray-600 dark:text-gray-300">{{ t('job_management.relationship_organization') }}</th>
              <th class="px-3 py-2 font-semibold text-gray-600 dark:text-gray-300">{{ t('job_management.relationship_activity') }}</th>
              <th class="px-3 py-2 w-12"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, i) in detailRows" :key="row._key" class="border-t border-gray-200 dark:border-gray-700">
              <td class="px-3 py-2 align-top text-gray-500 dark:text-gray-400">{{ i + 1 }}</td>
              <td class="px-3 py-2">
                <SelectLabel
                  v-model="row.organization_id"
                  :options="orgOptions"
                  option-label="label"
                  option-value="value"
                  :placeholder="t('common.select')"
                  showClear
                />
              </td>
              <td class="px-3 py-2">
                <TextInput
                  v-model="row.activity"
                  :placeholder="t('job_management.relationship_activity')"
                />
              </td>
              <td class="px-3 py-2 align-top">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  size="small"
                  text
                  rounded
                  aria-label="Remove"
                  @click="removeDetailRow(i)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Error display -->
      <div v-if="detailErrorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
        {{ detailErrorMsg }}
      </div>

      <!-- Save details -->
      <div v-if="detailRows.length > 0" class="flex justify-end gap-2 pt-2">
        <Button
          :label="t('job_management.save_relationship_details')"
          icon="pi pi-save"
          size="small"
          :loading="savingDetails"
          :disabled="savingDetails || !existingId"
          @click="handleSaveDetails"
        />
      </div>
    </div>

    <ConfirmDeleteDialog v-model:visible="deleteVisible" :loading="deleting" :error-msg="deleteError" @confirm="handleDelete" @cancel="deleteVisible=false" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  orgName: { type: String, default: '' },
  orgCode: { type: String, default: '' },
  orgSummaryId: { type: String, default: '' },
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
  job_management_value_relationship_id: '',
  job_management_value_frequency_id: ''
})

const apiBase = '/api/v1/tenant/job-management/relationships'

const relOptions = ref([])
const freqOptions = ref([])
const orgOptions = ref([])

// =============================================================
// Relationship Details state
// =============================================================
const detailRows = ref([])
const savingDetails = ref(false)
const detailErrorMsg = ref('')

// =============================================================
// Options loading
// =============================================================

// Muat opsi: Tipe Hubungan (type=relationship), Frekuensi (type=frequency),
// dan organisasi se-summary (untuk group Hubungan Kerja pada detail)
async function loadOptions() {
  try {
    const [relRes, freqRes, orgRes] = await Promise.all([
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'relationship', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'frequency', per_page: 100 } }),
      props.orgSummaryId
        ? api.get('/api/v1/tenant/organizations', { params: { summary_id: props.orgSummaryId, per_page: 100 } })
        : Promise.resolve({ data: { data: [] } })
    ])
    // Label menampilkan level (mis. 'Lv.1 — ...') untuk memudahkan identifikasi
    relOptions.value = (relRes.data?.data || []).map(v => ({ label: `Lv.${v.level} — ${v.descriptions}`, value: v.id }))
    freqOptions.value = (freqRes.data?.data || []).map(v => ({ label: `Lv.${v.level} — ${v.descriptions}`, value: v.id }))
    // Organisasi se-summary, tanpa organisasi yang sedang dinilai (tidak boleh pilih diri sendiri)
    orgOptions.value = (orgRes.data?.data || [])
      .filter(o => o.id !== props.orgId)
      .map(o => ({
        label: o.full_code ? `${o.full_code} - ${o.nomenclature}` : o.nomenclature,
        value: o.id
      }))
  } catch { /* ignore */ }
}

// =============================================================
// Main relationship (single record per org)
// =============================================================

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
      form.value.job_management_value_relationship_id = item.job_management_value_relationship_id || ''
      form.value.job_management_value_frequency_id = item.job_management_value_frequency_id || ''
      await loadDetails()
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
      job_management_value_relationship_id: form.value.job_management_value_relationship_id || null,
      job_management_value_frequency_id: form.value.job_management_value_frequency_id || null,
      organization_id: props.orgId
    }

    if (existingId.value) {
      await api.put(`${apiBase}/${existingId.value}`, {
        job_management_value_relationship_id: form.value.job_management_value_relationship_id || '',
        job_management_value_frequency_id: form.value.job_management_value_frequency_id || ''
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
    form.value.job_management_value_relationship_id = ''
    form.value.job_management_value_frequency_id = ''
    detailRows.value = []
    emit('saved')
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
  } catch (err) {
    deleteError.value = err?.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

// =============================================================
// Relationship Details — CRUD
// =============================================================

let detailKeySeq = 0

function addDetailRow() {
  if (!existingId.value) return
  detailRows.value.push({ _key: `new-${++detailKeySeq}`, id: '', organization_id: '', activity: '' })
}

function removeDetailRow(i) {
  const row = detailRows.value[i]
  if (!row) return
  // Row yang sudah tersimpan di DB dihapus langsung
  if (row.id) {
    deleteDetailRow(row.id, i)
  } else {
    detailRows.value.splice(i, 1)
  }
}

async function deleteDetailRow(detailId, index) {
  try {
    await api.delete(`${apiBase}/${existingId.value}/details/${detailId}`)
    detailRows.value.splice(index, 1)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
  } catch (err) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: err?.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000
    })
  }
}

async function loadDetails() {
  if (!existingId.value) return
  try {
    const res = await api.get(`${apiBase}/${existingId.value}/details`)
    detailRows.value = (res.data?.data || []).map(d => ({
      _key: `db-${++detailKeySeq}`,
      id: d.id,
      organization_id: d.organization_id || '',
      activity: d.activity || ''
    }))
  } catch { /* ignore */ }
}

async function handleSaveDetails() {
  if (!existingId.value || savingDetails.value) return
  detailErrorMsg.value = ''
  savingDetails.value = true

  try {
    for (const row of detailRows.value) {
      const payload = {
        organization_id: row.organization_id || '',
        activity: row.activity || ''
      }
      if (row.id) {
        await api.put(`${apiBase}/${existingId.value}/details/${row.id}`, payload)
      } else {
        const res = await api.post(`${apiBase}/${existingId.value}/details`, payload)
        row.id = res.data?.data?.id || ''
      }
    }
    await loadDetails()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_management.relationship_details_saved'), life: 2000 })
  } catch (err) {
    const ve = getValidationErrors(err)
    if (Object.keys(ve).length > 0) {
      detailErrorMsg.value = Object.values(ve).join(', ')
    } else {
      detailErrorMsg.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
    }
  } finally {
    savingDetails.value = false
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
