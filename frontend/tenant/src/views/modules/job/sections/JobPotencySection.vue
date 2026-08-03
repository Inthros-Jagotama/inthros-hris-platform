<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_competencies') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_description') }}</p>
    </div>

    <!-- Card paling atas: Potensi yang harus dimiliki (kompetensi psikologi + level) — lebar full -->
    <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <div>
        <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_required_title') }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_required_description') }}</p>
      </div>

      <!-- Skeleton while loading potency data -->
      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="5" cols="grid-cols-1" padding="p-5" />

      <template v-else>
        <!-- Empty state -->
        <div
          v-if="psychRows.length === 0"
          class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
        >
          {{ t('job_management.potency_required_empty') }}
        </div>

        <!-- Baris per kompetensi psikologi: nama + definisi (kiri), level select (kanan) -->
        <div
          v-for="row in psychRows"
          :key="row.competency_id || row.type"
          class="flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-3 border-b border-gray-100 dark:border-gray-700 last:border-b-0"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
            <div v-if="row.competency_definition" class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
              {{ row.competency_definition }}
            </div>
          </div>
          <div class="w-full md:w-80 shrink-0">
            <SelectLabel
              v-model="row.job_management_value_id"
              :options="row.levelOptions"
              option-label="label"
              option-value="value"
              :placeholder="t('common.select')"
              showClear
            />
          </div>
        </div>

        <!-- Error display -->
        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <!-- Actions -->
        <div v-if="psychRows.length > 0" class="flex justify-end gap-2 pt-1">
          <Button
            :label="t('job_management.save_potency_levels')"
            icon="pi pi-check"
            size="small"
            :loading="savingCard"
            :disabled="savingCard || !orgId"
            @click="handleSavePsych"
          />
        </div>
      </template>
    </div>

    <!-- Card: Communication and Influencing Skills — lebar full -->
    <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <div>
        <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.skill_communicating_influencing_title') }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.skill_communicating_influencing_description') }}</p>
      </div>

      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="2" cols="grid-cols-1" padding="p-5" />

      <template v-else>
        <!-- Empty state -->
        <div
          v-if="skillRows.length === 0"
          class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
        >
          {{ t('job_management.skill_communicating_influencing_empty') }}
        </div>

        <!-- Baris keterampilan: nama (kiri), level select (kanan) -->
        <div
          v-for="row in skillRows"
          :key="row.type"
          class="flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-3 border-b border-gray-100 dark:border-gray-700 last:border-b-0"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
            <div v-if="row.competency_definition" class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
              {{ row.competency_definition }}
            </div>
          </div>
          <div class="w-full md:w-80 shrink-0">
            <SelectLabel
              v-model="row.job_management_value_id"
              :options="row.levelOptions"
              option-label="label"
              option-value="value"
              :placeholder="t('common.select')"
              showClear
            />
          </div>
        </div>

        <!-- Error display -->
        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <!-- Actions -->
        <div v-if="skillRows.length > 0" class="flex justify-end gap-2 pt-1">
          <Button
            :label="t('job_management.save_skill')"
            icon="pi pi-check"
            size="small"
            :loading="savingCard"
            :disabled="savingCard || !orgId"
            @click="handleSaveSkill"
          />
        </div>
      </template>
    </div>

    <!-- Form kompetensi tambahan (kompetensi lain dengan bobot) -->
    <div class="max-w-2xl">
      <div v-if="loading" class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <SkeletonCard type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-5" />
      </div>

      <div v-else class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between gap-2 flex-wrap">
          <div>
            <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_competencies') }}</h3>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_description') }}</p>
          </div>
          <Button
            :label="t('common.add')"
            icon="pi pi-plus"
            size="small"
            outlined
            :disabled="savingRows"
            @click="addRow"
          />
        </div>

        <!-- Empty state -->
        <div
          v-if="rows.length === 0"
          class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
        >
          {{ t('job_management.empty_potency') }}
        </div>

        <!-- Multi-row form -->
        <div v-for="(row, i) in rows" :key="row._key" class="space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">{{ t('job_management.potency_item') }} {{ i + 1 }}</span>
            <Button
              icon="pi pi-trash"
              severity="danger"
              size="small"
              text
              rounded
              aria-label="Remove"
              @click="removeRow(i)"
            />
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <FormRow :label="t('job_management.competency')">
              <SelectLabel
                v-model="row.competency_id"
                :options="competencyOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                showClear
              />
            </FormRow>
            <FormRow :label="t('job_management.value_ref')">
              <SelectLabel
                v-model="row.job_management_value_id"
                :options="allOptions"
                option-label="label"
                option-value="value"
                :placeholder="t('common.select')"
                showClear
              />
            </FormRow>
            <FormRow :label="t('job_management.weight')">
              <InputNumber
                v-model="row.weight"
                :min="0"
                :max="100"
                size="small"
                class="w-full"
              />
            </FormRow>
          </div>
        </div>

        <!-- Error display -->
        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <!-- Actions -->
        <div v-if="rows.length > 0" class="flex justify-end gap-2 pt-2">
          <Button
            :label="t('job_management.save_potency')"
            icon="pi pi-check"
            size="small"
            :loading="savingRows"
            :disabled="savingRows"
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
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  orgName: { type: String, default: '' },
  orgCode: { type: String, default: '' },
  jobValueMap: { type: Object, default: () => ({}) },
  competencyOptions: { type: Array, default: () => [] }
})

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const savingRows = ref(false)
const errorMsg = ref('')
const rows = ref([])
// Raw record potency dari API (termasuk kompetensi psikologi) — dasar hydration card atas
let allRecords = []

// =========================================================================
// Card atas: Potensi yang harus dimiliki (kompetensi psikologi + level)
// =========================================================================
const psychRows = ref([])
const savingCard = ref(false)

// Id job_management_values tipe 'kecerdasan' (untuk match record tanpa competency_id)
const kecerdasanValueIds = computed(() => new Set(((props.jobValueMap && props.jobValueMap['kecerdasan']) || []).map(o => o.value)))

// Pemetaan nama kompetensi psikologi (field 'Potential') → tipe job value psikologi
const psychTypeMap = {
  'tenacity': 'tenacity',
  'creativy & innovation': 'innovation_creativity',
  'creativity & innovation': 'innovation_creativity',
  'self confidence': 'self_confidence',
  'flexibility': 'flexibility',
  'continuous learning': 'continuous_learning'
}

function psychTypeForName(name) {
  const n = (name || '').toLowerCase().replace(/\(.*?\)/g, '').replace(/\s+/g, ' ').trim()
  return psychTypeMap[n] || ''
}

// Bangun daftar kompetensi psikologi + opsi level dari jobValueMap
function buildPsychRows() {
  const rows = (props.competencyOptions || [])
    .filter(c => psychTypeForName(c.label))
    .map(c => {
      const type = psychTypeForName(c.label)
      return {
        competency_id: c.value,
        competency_name: c.label,
        competency_definition: c.definition || '',
        type,
        levelOptions: (props.jobValueMap && props.jobValueMap[type]) || [],
        recordId: '',
        job_management_value_id: ''
      }
    })

  // Kecerdasan (Intelligence) — dimensi psikologi tanpa kompetensi master,
  // disimpan dengan competency_id kosong (hanya job_management_value_id = level)
  const kecerdasanOptions = (props.jobValueMap && props.jobValueMap['kecerdasan']) || []
  if (kecerdasanOptions.length > 0) {
    rows.unshift({
      competency_id: '',
      competency_name: t('job_management.potency_kecerdasan'),
      competency_definition: '',
      type: 'kecerdasan',
      levelOptions: kecerdasanOptions,
      recordId: '',
      job_management_value_id: ''
    })
  }

  psychRows.value = rows
}

// Cocokkan record tanpa competency_id untuk baris tetap (kecerdasan/skill)
function findFixedRecord(valueIds) {
  return allRecords.find(r => !r.competency_id && valueIds.value.has(r.job_management_value_id)) || null
}

// Isi level terpilih dari record potency yang sudah tersimpan (per competency)
function hydratePsychRows() {
  const byComp = {}
  allRecords.forEach(r => { if (r.competency_id) byComp[r.competency_id] = r })
  psychRows.value.forEach(row => {
    // Kecerdasan disimpan tanpa competency_id → cocokkan via value id tipe 'kecerdasan'
    const rec = row.competency_id
      ? (byComp[row.competency_id] || null)
      : findFixedRecord(kecerdasanValueIds)
    row.recordId = rec ? rec.id : ''
    row.job_management_value_id = rec ? (rec.job_management_value_id || '') : ''
  })
}

// Simpan level card (POST/PUT/DELETE per baris) — dipakai card potensi & card skill
async function saveCardRows(targetRows) {
  errorMsg.value = ''
  savingCard.value = true
  try {
    for (const row of targetRows) {
      if (row.job_management_value_id) {
        // Kecerdasan tanpa competency_id → kirim hanya level (job_management_value_id)
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
    hydrateSkillRows()
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

function handleSavePsych() {
  saveCardRows(psychRows.value)
}

function handleSaveSkill() {
  saveCardRows(skillRows.value)
}

// =========================================================================
// Card: Communication and Influencing Skills (level keterampilan)
// =========================================================================
const skillRows = ref([])

// Id job_management_values tipe 'communicating_influencing_skill' (untuk match record tanpa competency_id)
const skillValueIds = computed(() => new Set(((props.jobValueMap && props.jobValueMap['communicating_influencing_skill']) || []).map(o => o.value)))

// Bangun baris keterampilan — tanpa competency_id (tidak ada kompetensi master untuk grup ini)
function buildSkillRows() {
  const options = (props.jobValueMap && props.jobValueMap['communicating_influencing_skill']) || []
  skillRows.value = options.length > 0
    ? [{
        competency_id: '',
        competency_name: t('job_management.skill_communicating_influencing'),
        competency_definition: '',
        type: 'communicating_influencing_skill',
        levelOptions: options,
        recordId: '',
        job_management_value_id: ''
      }]
    : []
}

// Isi level terpilih dari record potency (cocokkan via value id tipe skill)
function hydrateSkillRows() {
  skillRows.value.forEach(row => {
    const rec = findFixedRecord(skillValueIds)
    row.recordId = rec ? rec.id : ''
    row.job_management_value_id = rec ? (rec.job_management_value_id || '') : ''
  })
}

// =========================================================================
// Form kompetensi tambahan (multi-row dengan bobot)
// =========================================================================
const apiBase = '/api/v1/tenant/job-management/potency-competencies'

// Semua nilai jabatan (dari jobValueMap parent) untuk kolom Nilai Referensi
const allOptions = computed(() => Object.values(props.jobValueMap || {}).flat())

let keySeq = 0

function addRow() {
  rows.value.push({ _key: `new-${++keySeq}`, id: '', competency_id: '', job_management_value_id: '', weight: null })
}

function removeRow(i) {
  const row = rows.value[i]
  if (!row) return
  if (row.id) {
    deleteRow(row.id, i)
  } else {
    rows.value.splice(i, 1)
  }
}

async function deleteRow(rowId, index) {
  try {
    await api.delete(`${apiBase}/${rowId}`)
    rows.value.splice(index, 1)
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

async function loadData() {
  if (!props.orgId) {
    allRecords = []
    rows.value = []
    return
  }
  try {
    const res = await api.get(apiBase, { params: { organization_id: props.orgId, per_page: 100 } })
    allRecords = res.data?.data || []
    // Form bawah hanya menampilkan kompetensi NON-psikologi (psikologi dikelola di card atas)
    // — exclude: kompetensi psikologi (by competency_id) & kecerdasan (by value id tipe 'kecerdasan')
    const psychIds = new Set(psychRows.value.filter(r => r.competency_id).map(r => r.competency_id))
    rows.value = allRecords
      .filter(d => !(d.competency_id && psychIds.has(d.competency_id)))
      .filter(d => !(!d.competency_id && kecerdasanValueIds.value.has(d.job_management_value_id)))
      .filter(d => !(!d.competency_id && skillValueIds.value.has(d.job_management_value_id)))
      .map(d => ({
        _key: `db-${++keySeq}`,
        id: d.id,
        competency_id: d.competency_id || '',
        job_management_value_id: d.job_management_value_id || '',
        weight: d.weight ?? null
      }))
  } catch {
    allRecords = []
    rows.value = []
  }
}

async function handleSave() {
  errorMsg.value = ''
  savingRows.value = true
  try {
    for (const row of rows.value) {
      if (row.id) {
        await api.put(`${apiBase}/${row.id}`, {
          competency_id: row.competency_id || '',
          job_management_value_id: row.job_management_value_id || '',
          weight: row.weight ?? null
        })
      } else {
        const res = await api.post(apiBase, {
          competency_id: row.competency_id || '',
          job_management_value_id: row.job_management_value_id || '',
          weight: row.weight ?? null,
          organization_id: props.orgId
        })
        row.id = res.data?.data?.id || ''
      }
    }
    await loadData()
    hydratePsychRows()
    hydrateSkillRows()
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
    savingRows.value = false
  }
}

onMounted(async () => {
  buildPsychRows()
  buildSkillRows()
  try {
    await loadData()
  } finally {
    hydratePsychRows()
    hydrateSkillRows()
    loading.value = false
  }
})
</script>
