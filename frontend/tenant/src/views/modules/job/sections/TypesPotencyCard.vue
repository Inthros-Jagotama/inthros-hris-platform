<template>
  <PotencyLevelsCard
    :org-id="orgId"
    :rows="rows"
    :options-ready="optionsReady"
    :skeleton-rows="skeletonRows"
    :title-key="titleKey"
    :description-key="descriptionKey"
    :empty-key="emptyKey"
    :save-label-key="saveLabelKey"
    :delete-title-key="deleteTitleKey"
    :delete-message-key="deleteMessageKey"
    @saved="emit('saved')"
  />
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import PotencyLevelsCard from './PotencyLevelsCard.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  // Daftar tipe tetap: [{ type, nameKey }] — nameKey = locale key nama baris
  types: { type: Array, required: true },
  skeletonRows: { type: Number, default: 2 },
  titleKey: { type: String, required: true },
  descriptionKey: { type: String, required: true },
  emptyKey: { type: String, required: true },
  saveLabelKey: { type: String, required: true },
  deleteTitleKey: { type: String, required: true },
  deleteMessageKey: { type: String, required: true }
})

const { t } = useI18n()

const rows = ref([])
const optionsReady = ref(false)

// Bangun fixed rows: satu baris per tipe dengan level options masing-masing
function buildRows(levelOptionsMap) {
  rows.value = props.types
    .filter(tp => (levelOptionsMap[tp.type] || []).length > 0)
    .map(tp => ({
      competency_id: '',
      competency_name: t(tp.nameKey),
      competency_definition: '',
      type: tp.type,
      levelOptions: levelOptionsMap[tp.type] || [],
      recordId: '',
      job_management_value_id: ''
    }))
}

// Fetch tree values sekali (tanpa pagination — ListJobValuesTree memuat SEMUA nilai),
// ambil options hanya untuk tipe yang dibutuhkan. Dipakai juga oleh
// PsychologicalPotencyCard; per_page pada /values di-clamp server (max 100).
async function loadOptions() {
  try {
    const res = await api.get('/api/v1/tenant/job-management/values/tree')
    const byType = {}
    ;(res.data?.data || []).forEach(group => {
      ;(group.types || []).forEach(tp => {
        if (!props.types.some(pt => pt.type === tp.type)) return
        byType[tp.type] = (tp.options || []).map(o => ({
          label: `Lv.${o.level} — ${o.descriptions || ''}`,
          value: o.id,
          level: o.level,
          descriptions: o.descriptions || ''
        }))
      })
    })
    buildRows(byType)
  } catch {
    rows.value = []
  }
}

onMounted(async () => {
  await loadOptions()
  // Beri sinyal ke PotencyLevelsCard → loadData + hydrate level tersimpan
  optionsReady.value = true
})
</script>
