<template>
  <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
    <div class="flex items-start justify-between gap-4">
      <div>
        <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_technical_title') }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_technical_description') }}</p>
      </div>

      <!-- Bobot kompetensi teknikal (per organisasi) — tersimpan di title card -->
      <div class="flex flex-col items-end gap-1 shrink-0">
        <div class="flex items-center gap-2">
          <label class="text-xs font-medium text-gray-500 dark:text-gray-400 whitespace-nowrap">{{ t('job_management.potency_technical_weight_label') }}</label>
          <div class="w-24 shrink-0">
            <InputNumber
              v-model="groupWeight"
              fluid
              :min="0"
              :max="100"
              suffix="%"
              size="small"
              :disabled="savingWeight || !orgId"
              @blur="handleSaveWeight"
            />
          </div>
        </div>
        <div v-if="weightError" class="text-xs text-red-500 dark:text-red-400">{{ weightError }}</div>
      </div>
    </div>

    <SkeletonCard v-if="loading" type="detail" :count="1" :rows="8" cols="grid-cols-1" padding="p-5" />

    <template v-else>
      <!-- MultiSelect kompetensi dari cluster hasil mapping (diatur di Mapping Job Value → Technical) -->
      <SelectLabel
        v-model="selectedCompetencies"
        :options="competencyGroupOptions"
        option-label="name"
        option-value="id"
        option-group-label="label"
        option-group-children="items"
        :placeholder="t('job_management.potency_technical_placeholder')"
        showClear
        multiple
      />

      <div
        v-if="!hasMapping"
        class="text-sm text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg px-3 py-2"
      >
        {{ t('job_management.potency_technical_no_mapping') }}
      </div>

      <div
        v-else-if="rows.length === 0"
        class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
      >
        {{ t('job_management.potency_technical_empty') }}
      </div>

      <div v-else class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
              <th class="px-4 py-3 font-semibold min-w-[220px]">{{ t('job_management.potency_table_name') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_level') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[130px]">{{ t('job_management.potency_table_weight') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_description') }}</th>
              <th class="px-4 py-3 font-semibold w-16 text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in rows"
              :key="row.competency_id"
              class="border-t border-gray-100 dark:border-gray-700 align-top"
            >
              <td class="px-4 py-3">
                <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
                <div v-if="row.cluster" class="text-xs text-gray-400 dark:text-gray-500">{{ row.cluster }}</div>
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
              <td class="px-4 py-3">
                <InputNumber
                  v-model="row.weight"
                  class="!w-full"
                  :min="0"
                  :max="100"
                  suffix="%"
                  size="small"
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

      <div v-if="rows.length > 0" class="flex justify-end gap-2 pt-1">
        <Button
          :label="t('job_management.save_technical')"
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
import { usePotencyLevels } from '@/composables/usePotencyLevels'
import api from '@/services/api'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved', 'weight-saved'])

const props = defineProps({
  orgId: String
})

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)

// =========================================================================
// State: kompetensi + level technical + record potency
// =========================================================================
const competencies = ref([])
const values = ref([])
const rows = ref([])
const selectedCompetencies = ref([])
const mappedClusters = ref([])
const hasMapping = computed(() => mappedClusters.value.length > 0)

// =========================================================================
// State: bobot kompetensi teknikal (job_management_competency_groups, category=technical)
// =========================================================================
const groupWeight = ref('')
const groupRecordId = ref('')
const lastSavedWeight = ref('')
const savingWeight = ref(false)
const weightError = ref('')

// =========================================================================
// Logika bersama simpan/hapus/hydrate — record dicocokkan via competency_id
// =========================================================================
const {
  savingCard, errorMsg,
  deleteVisible, deleting, deleteError, deleteTarget,
  records,
  levelDescription, hydrateRows,
  loadData, askDeleteRow, handleDelete, handleSave
} = usePotencyLevels({
  orgId: computed(() => props.orgId),
  rows,
  matchBy: 'competency',
  // Kolom deskripsi level menampilkan note (deskripsi Indonesia), bukan nama level
  descriptionField: 'note',
  // Hapus baris → lepas kompetensi dari filter (baris hilang dari tabel)
  afterDelete: (row) => {
    const ids = Array.isArray(selectedCompetencies.value) ? selectedCompetencies.value : []
    selectedCompetencies.value = ids.filter(id => id !== row.competency_id)
  },
  onSaved: () => emit('saved')
})

// Kompetensi dari cluster hasil mapping (sebagai filter)
const competencyOptions = computed(() =>
  (competencies.value || []).map(c => ({
    id: c.id,
    name: c.name,
    cluster: c.cluster || ''
  }))
)

// Filter options dikelompokkan per cluster (label group = cluster, items = kompetensi)
const competencyGroupOptions = computed(() => {
  const groups = {}
  ;(competencyOptions.value || []).forEach(c => {
    ;(groups[c.cluster] = groups[c.cluster] || []).push(c)
  })
  return Object.keys(groups)
    .sort()
    .map(cluster => ({
      label: cluster,
      items: groups[cluster].sort((a, b) => a.name.localeCompare(b.name))
    }))
})

// Options level dari job_management_values type=technical
const levelOptions = computed(() =>
  (values.value || []).map(v => ({
    label: `Lv.${v.level} — ${v.descriptions || ''}`,
    value: v.id,
    level: v.level,
    descriptions: v.descriptions || '',
    note: v.note || ''
  }))
)

// Bangun baris hanya untuk kompetensi yang dipilih (selectedCompetencies)
function buildRows() {
  const byId = {}
  ;(competencyOptions.value || []).forEach(c => { byId[c.id] = c })
  const selected = Array.isArray(selectedCompetencies.value)
    ? selectedCompetencies.value
    : (selectedCompetencies.value ? [selectedCompetencies.value] : [])
  const ids = selected.filter(id => byId[id])
  // Bobot default: 100 dibagi jumlah kompetensi yang dipilih (2 desimal)
  const defaultWeight = ids.length > 0 ? Math.round((100 / ids.length) * 100) / 100 : 0
  rows.value = ids.map(id => {
    const c = byId[id]
    return {
      competency_id: id,
      competency_name: c.name,
      cluster: c.cluster,
      levelOptions: levelOptions.value,
      recordId: '',
      job_management_value_id: '',
      weight: defaultWeight
    }
  })
}

async function loadCompetencies() {
  try {
    const [compRes, mapRes] = await Promise.all([
      api.get('/api/v1/tenant/settings/competencies', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/job-management/values/clusters/technical')
    ])
    mappedClusters.value = mapRes.data?.data?.clusters || []
    const allowed = new Set(mappedClusters.value)
    const all = compRes.data?.data || []
    // Hanya kompetensi dari cluster yang di-mapping di halaman Mapping Job Value
    // (mapping kosong → tidak ada kompetensi yang tampil)
    competencies.value = all.filter(c => c.cluster && allowed.has(c.cluster))
  } catch {
    competencies.value = []
  }
}

async function loadValues() {
  try {
    const res = await api.get('/api/v1/tenant/job-management/values', { params: { type: 'technical', per_page: 100 } })
    values.value = res.data?.data || []
  } catch {
    values.value = []
  }
}

// Muat bobot kompetensi teknikal per organisasi (category=technical)
async function loadCompetencyGroup() {
  if (!props.orgId) {
    groupWeight.value = ''
    groupRecordId.value = ''
    return
  }
  try {
    const res = await api.get('/api/v1/tenant/job-management/competency-groups', {
      params: { organization_id: props.orgId }
    })
    const list = res.data?.data || []
    const rec = list.find(g => g.category === 'technical')
    groupRecordId.value = rec ? rec.id : ''
    groupWeight.value = rec ? rec.weight : ''
    lastSavedWeight.value = groupWeight.value
  } catch {
    groupWeight.value = ''
    groupRecordId.value = ''
    lastSavedWeight.value = ''
  }
}

// Simpan bobot kompetensi teknikal (create bila belum ada, update bila sudah)
async function handleSaveWeight() {
  if (groupWeight.value === '' || groupWeight.value === null || groupWeight.value === undefined) {
    weightError.value = t('job_management.potency_technical_weight_required')
    return
  }
  // Nilai tidak berubah sejak simpan terakhir → lewati (hindari PUT berulang tiap blur)
  if (lastSavedWeight.value !== '' && groupWeight.value === lastSavedWeight.value) {
    return
  }
  savingWeight.value = true
  weightError.value = ''
  try {
    // Re-check record terbaru sebelum POST (unique index org+category menolak duplikat)
    let recId = groupRecordId.value
    if (!recId) {
      try {
        const res = await api.get('/api/v1/tenant/job-management/competency-groups', {
          params: { organization_id: props.orgId }
        })
        const existing = (res.data?.data || []).find(g => g.category === 'technical')
        if (existing) recId = existing.id
      } catch { /* tetap lanjut POST */ }
    }
    const payload = { weight: groupWeight.value }
    if (recId) {
      await api.put(`/api/v1/tenant/job-management/competency-groups/${recId}`, payload)
    } else {
      await api.post('/api/v1/tenant/job-management/competency-groups', {
        organization_id: props.orgId,
        category: 'technical',
        weight: groupWeight.value
      })
    }
    lastSavedWeight.value = groupWeight.value
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_management.potency_technical_weight_saved'), life: 2000 })
    emit('saved')
    // Beri tahu parent nilai bobot teknis → card manajerial ikut menghitung & auto-save
    emit('weight-saved', groupWeight.value)
    await loadCompetencyGroup()
  } catch (err) {
    weightError.value = err?.response?.data?.error?.message || err.message || t('message.operation_failed')
  } finally {
    savingWeight.value = false
  }
}

// Pilih otomatis kompetensi yang sudah punya record tersimpan → data tampil saat kembali
function hydrateSelectedFromRecords() {
  const byId = {}
  ;(competencyOptions.value || []).forEach(c => { byId[c.id] = c })
  const ids = []
  records.value.forEach(r => {
    // Hanya kompetensi yang masih ada di filter (hindari chip phantom utk record
    // yang menunjuk ke cluster yang tidak di-mapping atau kompetensi yang sudah terhapus)
    if (r.competency_id && byId[r.competency_id] && !ids.includes(r.competency_id)) ids.push(r.competency_id)
  })
  selectedCompetencies.value = ids
  buildRows()
  hydrateRows()
}

// Rebuild baris saat pilihan kompetensi berubah
watch(selectedCompetencies, () => {
  buildRows()
  hydrateRows()
})

onMounted(async () => {
  try {
    await Promise.all([loadCompetencies(), loadValues(), loadData(), loadCompetencyGroup()])
  } finally {
    hydrateSelectedFromRecords()
    loading.value = false
  }
})
</script>
