<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.scores') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.score_description') }}</p>
    </div>

    <!-- Loading / Empty -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spin pi-spinner text-emerald-500 text-2xl"></i>
    </div>

    <template v-else-if="score">
      <!-- Component Breakdown — dari field components (JSON) di job_management_scores -->
      <div v-if="breakdown.length" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300">
          {{ t('job_management.component_breakdown') }}
        </div>

        <div class="divide-y divide-gray-100 dark:divide-gray-700">
          <!-- Header -->
          <div class="hidden md:grid grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] gap-4 px-5 py-2.5 bg-gray-50 dark:bg-gray-900/40 text-[11px] uppercase tracking-wider text-gray-400 dark:text-gray-500 font-medium">
            <span>{{ t('job_management.score_component') }}</span>
            <span>{{ t('job_management.score_points') }}</span>
            <span class="text-right">{{ t('job_management.score_score') }}</span>
          </div>

          <!-- Baris komponen -->
          <div v-for="comp in breakdown" :key="comp.key" class="px-5 py-1">
            <div class="grid grid-cols-1 md:grid-cols-[minmax(0,2fr)_minmax(0,3fr)_auto] md:items-center gap-2">
              <div class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t(comp.labelKey) }}</div>

              <!-- Detail poin -->
              <div class="flex flex-wrap gap-1.5">
                <template v-for="p in comp.points" :key="p.labelKey">
                  <span
                    class="inline-flex items-center gap-1 text-[11px] px-2 py-0.5 rounded-md border border-gray-200 dark:border-gray-600 bg-gray-50 dark:bg-gray-900/40 text-gray-600 dark:text-gray-300"
                    :class="{ 'opacity-50': p.level == null }"
                  >
                    <span class="font-medium">{{ t(p.labelKey) }}</span>
                    <span v-if="p.level != null" class="font-mono">Lv.{{ formatLevel(p.level) }}</span>
                    <template v-if="p.points != null">
                      <i class="pi pi-arrow-right text-[8px] opacity-60"></i>
                      <span class="font-bold text-emerald-600 dark:text-emerald-400">{{ formatNumber(p.points) }}</span>
                    </template>
                    <template v-else-if="p.level == null">—</template>
                  </span>
                </template>
              </div>

              <!-- Skor komponen -->
              <div class="text-right">
                <span class="text-sm font-bold text-gray-900 dark:text-gray-100">{{ formatNumber(comp.score) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Total -->
        <div class="px-5 py-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="flex items-center justify-between">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('job_management.value_with_financial') }}</span>
              <span class="text-sm font-bold text-emerald-600 dark:text-emerald-400">{{ formatNumber(score.job_value_with_financial) }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('job_management.value_without_financial') }}</span>
              <span class="text-sm font-bold text-blue-600 dark:text-blue-400">{{ formatNumber(score.job_value_without_financial) }}</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- No score yet -->
    <div v-else>
      <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
        <i class="pi pi-calculator text-4xl mb-3 opacity-50"></i>
        <p class="text-sm font-medium">{{ t('job_management.no_score') }}</p>
        <p class="text-xs mt-1">{{ t('job_management.score_hint') }}</p>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-3">
      <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" text @click="loadData" />
      <Button v-if="score" :label="t('job_management.recalculate')" icon="pi pi-calculator" size="small" severity="info" :loading="calculating" @click="recalculate" />
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import Button from 'primevue/button'

const props = defineProps({ orgId: String })
const emit = defineEmits(['saved'])
const { t } = useI18n()
const toast = useToast()
const apiBase = '/api/v1/tenant/job-management/scores/org'

const loading = ref(false)
const calculating = ref(false)
const score = ref(null)

// Urutan komponen mengikuti navigasi sub-form job management (education →
// potensi → kompetensi → keuangan → aset → bawahan → hubungan → aktivitas → risiko).
// Field `score` default 'score'; competencies memakai base_score agar tidak
// dobel dengan potentials & problem_solving (legacy: competencies.score =
// base + potential + problem_solving).
const COMPONENTS = [
  {
    key: 'education_experience',
    labelKey: 'job_management.education_experience',
    points: [
      { labelKey: 'job_management.group_education', level: 'education_level', pts: 'education_points' },
      { labelKey: 'job_management.group_experience', level: 'experience_level', pts: 'experience_points' }
    ]
  },
  {
    key: 'potentials',
    labelKey: 'job_management.score_potentials',
    points: [
      { labelKey: 'job_management.average_level', level: 'average_level', pts: null }
    ]
  },
  {
    key: 'competencies',
    labelKey: 'job_management.potency_competencies',
    score: 'base_score',
    points: [
      { labelKey: 'job_management.potency_technical_title', level: 'technical_average_level', pts: 'technical_points' },
      { labelKey: 'job_management.potency_managerial_title', level: 'managerial_average_level', pts: 'managerial_points' },
      { labelKey: 'job_management.skill_communicating_influencing', level: 'communication_level', pts: 'communication_points' }
    ]
  },
  {
    key: 'problem_solving',
    labelKey: 'job_management.problem_solving_title',
    points: [
      { labelKey: 'job_management.problem_solving_environment', level: 'environment_level', pts: 'environment_points' },
      { labelKey: 'job_management.problem_solving_challenge', level: 'challenge_level', pts: 'challenge_points' }
    ]
  },
  {
    key: 'financial_authority',
    labelKey: 'job_management.financials',
    points: [
      { labelKey: 'job_management.cash_level', level: 'money_level', pts: 'money_points' },
      { labelKey: 'job_management.authority_level', level: 'authority_level', pts: 'authority_points' },
      { labelKey: 'job_management.impact_level', level: 'impact_level', pts: 'impact_points' }
    ]
  },
  {
    key: 'asset_authority',
    labelKey: 'job_management.assets',
    points: [
      { labelKey: 'job_management.asset_type', level: 'asset_value_level', pts: 'asset_value_points' },
      { labelKey: 'job_management.authority_level', level: 'asset_authority_level', pts: 'asset_authority_points' }
    ]
  },
  {
    key: 'subordinate_control',
    labelKey: 'job_management.subordinate_controls',
    points: [
      { labelKey: 'job_management.score_level', level: 'level', pts: 'points' }
    ]
  },
  {
    key: 'work_scope',
    labelKey: 'job_management.relationships',
    points: [
      { labelKey: 'job_management.relationship_group_scope', level: 'scope_level', pts: 'scope_points' },
      { labelKey: 'job_management.frequency', level: 'frequency_level', pts: 'frequency_points' }
    ]
  },
  {
    key: 'work_activity',
    labelKey: 'job_management.activities',
    points: [
      { labelKey: 'job_management.score_level', level: 'level', pts: 'points' }
    ]
  },
  {
    key: 'work_risk',
    labelKey: 'job_management.risks',
    points: [
      { labelKey: 'job_management.environment_risk', level: 'environment_level', pts: 'environment_points' },
      { labelKey: 'job_management.hazard', level: 'hazard_level', pts: 'hazard_points' }
    ]
  }
]

// Parse field `components` (JSON nested dari job_management_scores)
const parsedComponents = computed(() => {
  if (!score.value?.components) return null
  try { return JSON.parse(score.value.components) }
  catch { return null }
})

const breakdown = computed(() => {
  if (!parsedComponents.value) return []
  return COMPONENTS.map(c => {
    const data = parsedComponents.value[c.key] || {}
    return {
      key: c.key,
      labelKey: c.labelKey,
      score: data[c.score || 'score'] ?? 0,
      points: c.points.map(p => ({
        labelKey: p.labelKey,
        level: data[p.level] ?? null,
        points: p.pts != null ? (data[p.pts] ?? 0) : null
      }))
    }
  })
})

function formatNumber(n) {
  return n?.toLocaleString?.('id-ID') ?? '-'
}

function formatLevel(v) {
  if (v == null) return '-'
  // Level bisa float (mis. rata-rata potensi) — tampilkan apa adanya
  return String(v)
}

async function loadData() {
  if (!props.orgId) return
  loading.value = true
  try {
    const res = await api.get(`${apiBase}/${props.orgId}`)
    score.value = res.data?.data || null
    emit('saved')
  } catch {
    score.value = null
  } finally {
    loading.value = false
  }
}

async function recalculate() {
  if (!props.orgId) return
  calculating.value = true
  try {
    // Recalculate: upsert with empty body triggers server-side recalculation
    const res = await api.put(`${apiBase}/${props.orgId}`, { components: null })
    score.value = res.data?.data || null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_management.score_calculated'), life: 2000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    calculating.value = false
  }
}

onMounted(loadData)
</script>
