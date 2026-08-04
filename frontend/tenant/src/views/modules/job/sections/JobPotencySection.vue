<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_competencies') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_description') }}</p>
    </div>

    <!-- Card: Potensi yang harus dimiliki — komponen terpisah (psychological multi-select + tabel) -->
    <PsychologicalPotencyCard :org-id="orgId" @saved="emit('saved')" />

    <!-- Card: Technical Competencies (16 kompetensi tetap) — lebar full -->
    <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <div>
        <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_technical_title') }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_technical_description') }}</p>
      </div>

      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="8" cols="grid-cols-1" padding="p-5" />

      <template v-else>
        <div
          v-if="technicalRows.length === 0"
          class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
        >
          {{ t('job_management.potency_technical_empty') }}
        </div>

        <div
          v-for="row in technicalRows"
          :key="row.type"
          class="flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-2.5 border-b border-gray-100 dark:border-gray-700 last:border-b-0"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
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

        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <div v-if="technicalRows.length > 0" class="flex justify-end gap-2 pt-1">
          <Button
            :label="t('job_management.save_technical')"
            icon="pi pi-check"
            size="small"
            :loading="savingCard"
            :disabled="savingCard || !orgId"
            @click="handleSaveTechnical"
          />
        </div>
      </template>
    </div>

    <!-- Card: Managerial Competencies (6 kompetensi tetap) — lebar full -->
    <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <div>
        <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_managerial_title') }}</h3>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_managerial_description') }}</p>
      </div>

      <SkeletonCard v-if="loading" type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-5" />

      <template v-else>
        <div
          v-if="managerialRows.length === 0"
          class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
        >
          {{ t('job_management.potency_managerial_empty') }}
        </div>

        <div
          v-for="row in managerialRows"
          :key="row.type"
          class="flex flex-col md:flex-row md:items-center gap-3 md:gap-6 py-2.5 border-b border-gray-100 dark:border-gray-700 last:border-b-0"
        >
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
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

        <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
          {{ errorMsg }}
        </div>

        <div v-if="managerialRows.length > 0" class="flex justify-end gap-2 pt-1">
          <Button
            :label="t('job_management.save_managerial')"
            icon="pi pi-check"
            size="small"
            :loading="savingCard"
            :disabled="savingCard || !orgId"
            @click="handleSaveManagerial"
          />
        </div>
      </template>
    </div>

    <!-- Card: Communication and Influencing Skills — komponen terpisah (tabel level keterampilan) -->
    <SkillPotencyCard :org-id="orgId" @saved="emit('saved')" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import PsychologicalPotencyCard from './PsychologicalPotencyCard.vue'
import SkillPotencyCard from './SkillPotencyCard.vue'

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
const savingCard = ref(false)
const errorMsg = ref('')
// Raw record potency dari API — dasar hydration semua card
let allRecords = []

const apiBase = '/api/v1/tenant/job-management/potency-competencies'

// =========================================================================
// Daftar tetap slug tipe (sama dengan kalkulator backend calculator.go)
// =========================================================================

const technicalTypes = [
  { type: 'competency_based_human_resources_management', name: 'Competency Based Human Resources Management' },
  { type: 'competency_development', name: 'Competency Development' },
  { type: 'people_development', name: 'People Development' },
  { type: 'career_management', name: 'Career Management' },
  { type: 'hr_assessment', name: 'HR Assessment' },
  { type: 'recruitement_selection', name: 'Recruitment & Selection' },
  { type: 'job_analysis_evaluation', name: 'Job Analysis & Evaluation' },
  { type: 'organizational_development', name: 'Organizational Development' },
  { type: 'human_resources_information_system', name: 'HR Information System' },
  { type: 'workload_analysis', name: 'Workload Analysis' },
  { type: 'performance_apraisal', name: 'Performance Appraisal' },
  { type: 'remuneration_manajemen', name: 'Remuneration Management' },
  { type: 'reward_punisment_management', name: 'Reward & Punishment Management' },
  { type: 'health_safety_environment', name: 'Health, Safety & Environment' },
  { type: 'hubungan_industrial', name: 'Industrial Relations' },
  { type: 'budgeting', name: 'Budgeting' }
]

const managerialTypes = [
  { type: 'integrity', name: 'Integrity' },
  { type: 'achievement_orientation', name: 'Achievement Orientation' },
  { type: 'building_partnership', name: 'Building Partnership' },
  { type: 'planning_organizing', name: 'Planning & Organizing' },
  { type: 'leadership', name: 'Leadership' },
  { type: 'developing_others', name: 'Developing Others' }
]

// =========================================================================
// Card: Technical & Managerial (kompetensi tetap + level)
// =========================================================================
const technicalRows = ref([])
const managerialRows = ref([])

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

// Simpan level card (POST/PUT/DELETE per baris) — dipakai semua card
async function saveCardRows(targetRows) {
  errorMsg.value = ''
  savingCard.value = true
  try {
    for (const row of targetRows) {
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
    hydrateFixedRows(technicalRows.value)
    hydrateFixedRows(managerialRows.value)
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

function handleSaveTechnical() {
  saveCardRows(technicalRows.value)
}

function handleSaveManagerial() {
  saveCardRows(managerialRows.value)
}

onMounted(async () => {
  buildTechnicalRows()
  buildManagerialRows()
  try {
    await loadData()
  } finally {
    hydrateFixedRows(technicalRows.value)
    hydrateFixedRows(managerialRows.value)
    loading.value = false
  }
})
</script>
