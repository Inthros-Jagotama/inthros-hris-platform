<template>
  <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
    <div>
      <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t(titleKey) }}</h3>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t(descriptionKey) }}</p>
    </div>

    <SkeletonCard v-if="loading" type="detail" :count="1" :rows="skeletonRows" cols="grid-cols-1" padding="p-5" />

    <template v-else>
      <!-- Multiple select tipe dalam group (dari tree) — full width tanpa label -->
      <SelectLabel
        v-model="selectedTypes"
        :options="typeOptions"
        option-label="label"
        option-value="value"
        :placeholder="t('common.select')"
        showClear
        multiple
      />

      <div
        v-if="rows.length === 0"
        class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
      >
        {{ t(emptyKey) }}
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
              v-for="row in rows"
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

      <div v-if="rows.length > 0" class="flex justify-end gap-2 pt-1">
        <Button
          :label="t(saveLabelKey)"
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
      :title="t(deleteTitleKey)"
      :message="t(deleteMessageKey, { name: deleteTarget?.competency_name || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      @confirm="handleDelete"
      @cancel="deleteVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { usePotencyLevels } from '@/composables/usePotencyLevels'
import api from '@/services/api'
import Button from 'primevue/button'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  // Slug group di tree (mis. 'psychological', 'technical', 'managerial')
  typeGroup: { type: String, required: true },
  skeletonRows: { type: Number, default: 5 },
  titleKey: { type: String, required: true },
  descriptionKey: { type: String, required: true },
  emptyKey: { type: String, required: true },
  saveLabelKey: { type: String, required: true },
  deleteTitleKey: { type: String, required: true },
  deleteMessageKey: { type: String, required: true }
})

const { t } = useI18n()

const loading = ref(true)

// =========================================================================
// State: tree job values + record potency
// =========================================================================
const treeData = ref([])
const rows = ref([])
const selectedTypes = ref([])

// =========================================================================
// Logika bersama simpan/hapus/hydrate (composable usePotencyLevels)
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
  // Lepas tipe dari pilihan → watch membangun ulang baris (baris hilang dari tabel)
  afterDelete: (row) => {
    const types = Array.isArray(selectedTypes.value) ? selectedTypes.value : []
    selectedTypes.value = types.filter(type => type !== row.type)
  },
  onSaved: () => emit('saved')
})

// Tipe-tipe dalam group yang dipilih dari tree
const group = computed(() =>
  (treeData.value || []).find(g => g.type_group === props.typeGroup)
)

// Label tipe: utamakan locale (job_values.types.<slug>), fallback description_group
function typeLabel(item) {
  const key = `job_values.types.${item.type}`
  const localized = t(key)
  return localized !== key ? localized : (item.description_group || item.type)
}

// Options multiple select — label = nama tipe (bilingual), value = type slug
const typeOptions = computed(() =>
  (group.value?.types || []).map(item => ({
    label: typeLabel(item),
    value: item.type
  }))
)

// Bangun baris hanya untuk tipe yang dipilih (selectedTypes)
function buildRows() {
  const byType = {}
  ;(group.value?.types || []).forEach(item => { byType[item.type] = item })
  // Normalisasi: pastikan selalu array (MultiSelect mengirim array, tapi
  // defensif bila nilai tunggal/undefined diterima)
  const selected = Array.isArray(selectedTypes.value)
    ? selectedTypes.value
    : (selectedTypes.value ? [selectedTypes.value] : [])
  rows.value = selected
    .filter(type => byType[type])
    .map(type => {
      const item = byType[type]
      return {
        competency_id: '',
        competency_name: typeLabel(item),
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
    buildRows()
  } catch {
    treeData.value = []
    rows.value = []
  }
}

// Pilih otomatis tipe yang sudah punya record tersimpan (via value id)
function hydrateSelectedTypesFromRecords() {
  const valueToType = {}
  ;(group.value?.types || []).forEach(item => {
    ;(item.options || []).forEach(o => { valueToType[o.id] = item.type })
  })
  const types = []
  records.value.forEach(r => {
    const type = r.job_management_value_id && valueToType[r.job_management_value_id]
    if (type && !types.includes(type)) types.push(type)
  })
  selectedTypes.value = types
  buildRows()
  hydrateRows()
}

// Rebuild baris saat pilihan tipe berubah
watch(selectedTypes, () => {
  buildRows()
  hydrateRows()
})

onMounted(async () => {
  try {
    await Promise.all([loadTree(), loadData()])
  } finally {
    // Pilih tipe yang punya record tersimpan → data tampil saat kembali
    hydrateSelectedTypesFromRecords()
    loading.value = false
  }
})
</script>
