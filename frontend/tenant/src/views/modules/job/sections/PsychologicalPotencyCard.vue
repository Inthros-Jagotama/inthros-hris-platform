<template>
  <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
    <div>
      <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_required_title') }}</h3>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_required_description') }}</p>
    </div>

    <SkeletonCard v-if="loading" type="detail" :count="1" :rows="5" cols="grid-cols-1" padding="p-5" />

    <template v-else>
      <!-- Multiple select tipe psikologis (group psychological) dari tree — full width tanpa label -->
      <SelectLabel
        v-model="selectedTypes"
        :options="psychTypeOptions"
        option-label="label"
        option-value="value"
        :placeholder="t('common.select')"
        showClear
        multiple
      />

      <div
        v-if="psychRows.length === 0"
        class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
      >
        {{ t('job_management.potency_required_empty') }}
      </div>

      <div v-else class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
              <th class="px-4 py-3 font-semibold min-w-[220px]">{{ t('job_management.potency_table_name') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_level') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_description') }}</th>
              <th class="px-4 py-3 font-semibold w-16 text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in psychRows"
              :key="row.type"
              class="border-t border-gray-100 dark:border-gray-700 align-top"
            >
              <td class="px-4 py-3">
                <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
              </td>
              <td class="px-4 py-3">
                <SelectLabel
                  v-model="row.job_management_value_id"
                  :options="row.levelOptions"
                  option-label="label"
                  option-value="value"
                  :placeholder="t('common.select')"
                  showClear
                />
              </td>
              <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                {{ levelDescription(row) }}
              </td>
              <td class="px-4 py-3 text-right">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  size="small"
                  :disabled="savingCard"
                  :aria-label="t('common.delete')"
                  @click="askDeleteRow(row)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
        {{ errorMsg }}
      </div>

      <div v-if="psychRows.length > 0" class="flex justify-end gap-2 pt-1">
        <Button
          :label="t('job_management.save_potency_levels')"
          icon="pi pi-check"
          size="small"
          :loading="savingCard"
          :disabled="savingCard || !orgId"
          @click="handleSave"
        />
      </div>
    </template>

    <ConfirmDeleteDialog
      v-model:visible="deleteVisible"
      :title="t('job_management.potency_confirm_delete_title')"
      :message="t('job_management.potency_confirm_delete', { name: deleteTarget?.competency_name || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      @confirm="handleDelete"
      @cancel="deleteVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String
})

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const savingCard = ref(false)
const errorMsg = ref('')
const deleteVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const apiBase = '/api/v1/tenant/job-management/potency-competencies'

// =========================================================================
// State: tree job values + record potency
// =========================================================================
const treeData = ref([])
const psychRows = ref([])
const selectedTypes = ref([])

// Raw record potency dari API — dasar hydration
let allRecords = []

// Tipe-tipe dalam group psychological dari tree
const psychGroup = computed(() =>
  (treeData.value || []).find(g => g.type_group === 'psychological')
)

// Options multiple select — label = description_group tipe, value = type slug
const psychTypeOptions = computed(() =>
  (psychGroup.value?.types || []).map(item => ({
    label: item.description_group || item.type,
    value: item.type
  }))
)

// Bangun baris hanya untuk tipe yang dipilih (selectedTypes)
function buildPsychRows() {
  const byType = {}
  ;(psychGroup.value?.types || []).forEach(item => { byType[item.type] = item })
  // Normalisasi: pastikan selalu array (MultiSelect mengirim array, tapi
  // defensif bila nilai tunggal/undefined diterima)
  const selected = Array.isArray(selectedTypes.value)
    ? selectedTypes.value
    : (selectedTypes.value ? [selectedTypes.value] : [])
  psychRows.value = selected
    .filter(type => byType[type])
    .map(type => {
      const item = byType[type]
      return {
        competency_id: '',
        competency_name: item.description_group || item.type,
        competency_definition: '',
        type: item.type,
        // Konversi option tree ({id, level, descriptions}) → format SelectLabel
        levelOptions: (item.options || []).map(o => ({
          label: `Lv.${o.level} — ${o.descriptions || ''}`,
          value: o.id,
          level: o.level,
          descriptions: o.descriptions || ''
        })),
        recordId: '',
        job_management_value_id: ''
      }
    })
}

async function loadTree() {
  try {
    const res = await api.get('/api/v1/tenant/job-management/values/tree')
    treeData.value = res.data?.data || []
    buildPsychRows()
  } catch {
    treeData.value = []
    psychRows.value = []
  }
}

// Pilih otomatis tipe psikologis yang sudah punya record tersimpan (via value id)
function hydrateSelectedTypesFromRecords() {
  const valueToType = {}
  ;(psychGroup.value?.types || []).forEach(item => {
    ;(item.options || []).forEach(o => { valueToType[o.id] = item.type })
  })
  const types = []
  allRecords.forEach(r => {
    const type = r.job_management_value_id && valueToType[r.job_management_value_id]
    if (type && !types.includes(type)) types.push(type)
  })
  selectedTypes.value = types
  buildPsychRows()
  hydratePsychRows()
}

// Cocokkan record tanpa competency_id (baris tetap) via value id
function findFixedRecord(valueIds) {
  return allRecords.find(r => r.job_management_value_id && valueIds.has(r.job_management_value_id)) || null
}

// Deskripsi level dari option yang dipilih
function levelDescription(row) {
  const opt = (row.levelOptions || []).find(o => o.value === row.job_management_value_id)
  return opt ? (opt.descriptions || '') : ''
}

// Isi level terpilih dari record potency yang sudah tersimpan (via value id per tipe)
function hydratePsychRows() {
  psychRows.value.forEach(row => {
    const valueIds = new Set(row.levelOptions.map(o => o.value))
    const rec = findFixedRecord(valueIds)
    row.recordId = rec ? rec.id : ''
    row.job_management_value_id = rec ? (rec.job_management_value_id || '') : ''
  })
}

// =========================================================================
// Delete per row (hapus record tersimpan + lepas tipe dari pilihan)
// =========================================================================

function askDeleteRow(row) {
  deleteTarget.value = row
  deleteError.value = ''
  deleteVisible.value = true
}

async function handleDelete() {
  const row = deleteTarget.value
  if (!row) return
  deleting.value = true
  deleteError.value = ''
  try {
    if (row.recordId) {
      await api.delete(`${apiBase}/${row.recordId}`)
    }
    // Lepas tipe dari pilihan → watch membangun ulang baris (row hilang dari tabel)
    const types = Array.isArray(selectedTypes.value) ? selectedTypes.value : []
    selectedTypes.value = types.filter(type => type !== row.type)
    deleteVisible.value = false
    await loadData()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 2000 })
    emit('saved')
  } catch (err) {
    deleteError.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
  } finally {
    deleting.value = false
    deleteTarget.value = null
  }
}

// =========================================================================
// Data & simpan (POST/PUT/DELETE per baris)
// =========================================================================

async function loadData() {
  if (!props.orgId) {
    allRecords = []
    return
  }
  try {
    const res = await api.get(apiBase, { params: { organization_id: props.orgId, per_page: 100 } })
    allRecords = res.data?.data || []
  } catch {
    allRecords = []
  }
}

async function handleSave() {
  errorMsg.value = ''
  savingCard.value = true
  try {
    for (const row of psychRows.value) {
      if (row.job_management_value_id) {
        // Baris tanpa competency_id → kirim hanya level (job_management_value_id)
        const payload = row.competency_id
          ? { competency_id: row.competency_id, job_management_value_id: row.job_management_value_id }
          : { job_management_value_id: row.job_management_value_id }
        if (row.recordId) {
          await api.put(`${apiBase}/${row.recordId}`, payload)
        } else {
          const res = await api.post(apiBase, { organization_id: props.orgId, ...payload })
          row.recordId = res.data?.data?.id || ''
        }
      } else if (row.recordId) {
        // Level dikosongkan → hapus record potency
        await api.delete(`${apiBase}/${row.recordId}`)
        row.recordId = ''
      }
    }
    await loadData()
    hydratePsychRows()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('common.saved'), life: 2000 })
    emit('saved')
  } catch (err) {
    const ve = getValidationErrors(err)
    if (Object.keys(ve).length > 0) {
      errorMsg.value = Object.values(ve).join(', ')
    } else {
      errorMsg.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
    }
  } finally {
    savingCard.value = false
  }
}

// Rebuild baris saat pilihan tipe psikologis berubah
watch(selectedTypes, () => {
  buildPsychRows()
  hydratePsychRows()
})

onMounted(async () => {
  try {
    await Promise.all([loadTree(), loadData()])
  } finally {
    // Pilih tipe psikologis yang punya record tersimpan → data tampil saat kembali
    hydrateSelectedTypesFromRecords()
    loading.value = false
  }
})
</script>
