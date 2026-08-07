<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="h-24 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
      <div class="h-64 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
    </div>

    <template v-else-if="evaluation">
      <!-- Header Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-start justify-between gap-4 flex-wrap">
          <div class="flex items-center gap-4">
            <div class="w-14 h-14 rounded-xl bg-emerald-50 dark:bg-emerald-500/10 flex items-center justify-center">
              <i class="pi pi-user text-2xl text-emerald-600 dark:text-emerald-400"></i>
            </div>
            <div>
              <h2 class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ evaluation.employee_name }}</h2>
              <p class="text-sm text-gray-500 dark:text-gray-400">{{ evaluation.organization_name || '-' }}</p>
              <div class="flex items-center gap-2 mt-1">
                <Tag :value="evaluation.period_code" severity="info" class="!text-xs" />
                <Tag :value="evaluation.status" :severity="getStatusSeverity(evaluation.status)" class="!text-xs" />
              </div>
            </div>
          </div>
          <div class="flex flex-col items-end gap-2">
            <div class="text-right">
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('kpi.final_score') }}</p>
              <p class="text-3xl font-bold" :class="getScoreClass(evaluation.final_score)">
                {{ evaluation.final_score?.toFixed(1) || '0.0' }}
              </p>
            </div>
            <Tag v-if="evaluation.rating_name" :value="evaluation.rating_name" :severity="getRatingSeverity(evaluation.rating_color)" />
          </div>
        </div>
      </div>

      <!-- Indicators Table -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-list text-blue-500"></i>
            {{ t('kpi.indicators') }}
            <span class="text-xs font-normal text-gray-500">({{ details.length }})</span>
          </h3>
          <div class="flex items-center gap-2">
            <Button v-if="canEdit" :label="t('kpi.recalculate')" icon="pi pi-refresh" size="small" outlined severity="secondary" :loading="recalculating" @click="recalculate" />
          </div>
        </div>

        <DataTable :value="details" size="small" class="!text-sm p-datatable-sm" :rowHover="true" rowGroupMode="subheader" groupRowsBy="perspective_name" sortField="perspective_name" :sortOrder="1">
          <template #groupheader="{data}">
            <div class="flex items-center gap-2 py-1">
              <i class="pi pi-th-large text-indigo-500"></i>
              <span class="font-semibold text-gray-700 dark:text-gray-200">{{ data.perspective_name || 'Other' }}</span>
            </div>
          </template>

          <Column field="indicator_title" :header="t('kpi.indicator_title')" style="min-width:200px">
            <template #body="{data}">
              <div>
                <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.indicator_title }}</span>
                <p v-if="data.indicator_description" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-1">{{ data.indicator_description }}</p>
              </div>
            </template>
          </Column>

          <Column field="weight" :header="t('kpi.weight')" style="width:80px">
            <template #body="{data}">
              <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.weight?.toFixed(1) }}%</span>
            </template>
          </Column>

          <Column field="target_value" :header="t('kpi.target')" style="width:100px">
            <template #body="{data}">
              <span class="text-gray-600 dark:text-gray-300 font-mono">{{ formatNumber(data.target_value) }} {{ data.unit_of_measurement }}</span>
            </template>
          </Column>

          <Column field="actual_value" :header="t('kpi.actual')" style="width:120px">
            <template #body="{data}">
              <InputNumber
                v-if="canEdit"
                v-model="data.actual_value"
                :minFractionDigits="0"
                :maxFractionDigits="2"
                class="w-full !text-xs"
                size="small"
                @blur="updateActual(data)"
              />
              <span v-else class="text-gray-800 dark:text-gray-100 font-mono font-semibold">{{ formatNumber(data.actual_value) }}</span>
            </template>
          </Column>

          <Column field="achievement" :header="t('kpi.achievement')" style="width:100px">
            <template #body="{data}">
              <span class="font-mono font-semibold" :class="getAchievementClass(data.achievement)">
                {{ data.achievement?.toFixed(1) || '0.0' }}%
              </span>
            </template>
          </Column>

          <Column field="score" :header="t('kpi.score')" style="width:80px">
            <template #body="{data}">
              <span class="font-mono font-bold" :class="getScoreClass(data.score)">
                {{ data.score?.toFixed(1) || '0.0' }}
              </span>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- Scoring Components Breakdown (Phase 5) -->
      <div v-if="components.length > 0 || hasScoringConfig" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-sliders-h text-indigo-500"></i>
            {{ t('performance_scoring.components') }}
          </h3>
          <Button v-if="canEdit" :label="t('performance_scoring.calculate')" icon="pi pi-calculator" size="small" outlined severity="secondary" :loading="calculatingScoring" @click="calculateScoring" />
        </div>

        <p v-if="components.length === 0" class="text-xs text-gray-400 dark:text-gray-500 py-2">
          {{ t('performance_scoring.not_calculated_yet') }}
        </p>

        <DataTable v-else :value="components" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
          <Column field="component_name" :header="t('performance_components.name')" style="min-width:160px">
            <template #body="{data}">
              <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.component_name }}</span>
            </template>
          </Column>
          <Column field="score" :header="t('performance_scoring.raw_score')" style="width:120px">
            <template #body="{data}">
              <InputNumber
                v-if="canEdit && isManualComponent(data)"
                v-model="data.score"
                :min="0"
                :max="100"
                :minFractionDigits="0"
                :maxFractionDigits="2"
                class="w-full !text-xs"
                size="small"
                @blur="updateComponentScore(data)"
              />
              <span v-else class="text-gray-600 dark:text-gray-300 font-mono">{{ data.score?.toFixed(1) ?? '0.0' }}</span>
            </template>
          </Column>
          <Column field="weight" :header="t('okr.weight')" style="width:100px">
            <template #body="{data}">
              <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.weight?.toFixed(1) ?? '0.0' }}%</span>
            </template>
          </Column>
          <Column field="final_score" :header="t('performance_scoring.weighted_score')" style="width:120px">
            <template #body="{data}">
              <span class="font-mono font-bold" :class="getScoreClass(data.final_score)">{{ data.final_score?.toFixed(1) ?? '0.0' }}</span>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="goBack" />
        <div class="flex items-center gap-2">
          <template v-if="evaluation.status === 'DRAFT'">
            <Button :label="t('kpi.submit')" icon="pi pi-send" size="small" :loading="submitting" @click="submitEvaluation" />
          </template>
          <template v-else-if="evaluation.status === 'SUBMITTED'">
            <Button :label="t('kpi.reject')" icon="pi pi-times" severity="danger" outlined size="small" :loading="rejecting" @click="rejectEvaluation" />
            <Button :label="t('kpi.approve')" icon="pi pi-check" severity="success" size="small" :loading="approving" @click="approveEvaluation" />
          </template>
          <template v-else-if="evaluation.status === 'APPROVED'">
            <Button :label="t('kpi.complete')" icon="pi pi-check-circle" severity="success" size="small" :loading="completing" @click="completeEvaluation" />
          </template>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-exclamation-circle text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('kpi.evaluation_not_found') }}</p>
      <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" class="mt-4" @click="goBack" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputNumber from 'primevue/inputnumber'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const { t } = useI18n()

const pageLoading = ref(true)
const evaluation = ref(null)
const details = ref([])
const recalculating = ref(false)
const submitting = ref(false)
const approving = ref(false)
const rejecting = ref(false)
const completing = ref(false)

const components = ref([])
const hasScoringConfig = ref(false)
const calculatingScoring = ref(false)
const componentCodeById = ref({})

const evaluationId = computed(() => route.params.id)
const canEdit = computed(() => evaluation.value?.status === 'DRAFT')

function isManualComponent(row) {
  const code = componentCodeById.value[row.component_id]
  return code !== 'KPI' && code !== 'SUBORDINATE'
}

function formatNumber(val) {
  if (val == null) return '-'
  return Number(val).toLocaleString('id-ID', { maximumFractionDigits: 2 })
}

function getStatusSeverity(status) {
  switch (status) {
    case 'COMPLETED': return 'success'
    case 'APPROVED': return 'info'
    case 'SUBMITTED': return 'warn'
    default: return 'secondary'
  }
}

function getRatingSeverity(color) {
  switch (color) {
    case 'success': return 'success'
    case 'primary': return 'info'
    case 'warning': return 'warn'
    case 'danger': return 'danger'
    default: return 'secondary'
  }
}

function getScoreClass(score) {
  if (!score) return 'text-gray-400'
  if (score >= 85) return 'text-emerald-600 dark:text-emerald-400'
  if (score >= 70) return 'text-blue-600 dark:text-blue-400'
  if (score >= 60) return 'text-amber-600 dark:text-amber-400'
  return 'text-red-600 dark:text-red-400'
}

function getAchievementClass(achievement) {
  if (!achievement) return 'text-gray-400'
  if (achievement >= 100) return 'text-emerald-600 dark:text-emerald-400'
  if (achievement >= 80) return 'text-blue-600 dark:text-blue-400'
  if (achievement >= 60) return 'text-amber-600 dark:text-amber-400'
  return 'text-red-600 dark:text-red-400'
}

function goBack() {
  router.push('/performance/kpi')
}

async function loadEvaluation() {
  pageLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/full`)
    const data = res.data?.data || res.data
    evaluation.value = {
      id: data.id,
      employee_id: data.employee_id,
      organization_id: data.organization_id,
      employee_name: data.employee_name || data.employee?.full_name || '-',
      organization_name: data.organization_name || data.organization?.name || '-',
      period_code: data.period_code || data.period?.period_code || '-',
      status: data.status,
      final_score: data.final_score,
      rating_name: data.rating_name || data.rating?.name,
      rating_color: data.rating_color || data.rating?.color
    }
    details.value = (data.details || []).map(d => ({
      id: d.id,
      indicator_title: d.indicator_name || d.indicator_title || d.indicator?.title,
      indicator_description: d.description || d.indicator_description || d.indicator?.description,
      perspective_name: d.perspective_name || d.perspective?.name || 'Other',
      weight: d.weight,
      target_value: d.target ?? d.target_value,
      actual_value: d.actual ?? d.actual_value,
      achievement: d.achievement,
      score: d.score,
      unit_of_measurement: d.unit_of_measurement || d.indicator?.unit_of_measurement || '',
      formula_type: d.formula_type || d.indicator?.formula_type
    }))
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
    evaluation.value = null
  } finally {
    pageLoading.value = false
  }

  if (evaluation.value) {
    await Promise.all([loadComponents(), loadMasterComponents(), checkScoringConfig()])
  }
}

async function loadComponents() {
  try {
    const res = await api.get(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/components`)
    components.value = res.data?.data || []
  } catch {
    components.value = []
  }
}

async function loadMasterComponents() {
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/components', { params: { per_page: 100 } })
    const map = {}
    for (const c of res.data?.data || []) {
      map[c.id] = c.code
    }
    componentCodeById.value = map
  } catch {
    componentCodeById.value = {}
  }
}

async function checkScoringConfig() {
  if (!evaluation.value?.organization_id) {
    hasScoringConfig.value = false
    return
  }
  try {
    const res = await api.get(`/api/v1/tenant/performance/kpi/organizations/${evaluation.value.organization_id}/components`)
    hasScoringConfig.value = (res.data?.data || []).some(c => c.is_enabled)
  } catch {
    hasScoringConfig.value = false
  }
}

async function calculateScoring() {
  calculatingScoring.value = true
  try {
    const res = await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/calculate-scoring`)
    components.value = res.data?.data || []
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_scoring.calculated'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    calculatingScoring.value = false
  }
}

async function updateComponentScore(row) {
  try {
    const res = await api.put(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/components/${row.component_id}`, {
      score: row.score || 0
    })
    components.value = res.data?.data || components.value
    await loadEvaluation()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

async function updateActual(detail) {
  try {
    await api.put(`/api/v1/tenant/performance/kpi/evaluation-details/${detail.id}/actual`, {
      actual: detail.actual_value || 0
    })
    await recalculate()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

async function recalculate() {
  recalculating.value = true
  try {
    await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/recalculate`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.score_recalculated'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    recalculating.value = false
  }
}

async function submitEvaluation() {
  submitting.value = true
  try {
    await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/submit`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.evaluation_submitted'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    submitting.value = false
  }
}

async function approveEvaluation() {
  approving.value = true
  try {
    await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/approve`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.evaluation_approved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    approving.value = false
  }
}

async function rejectEvaluation() {
  rejecting.value = true
  try {
    await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/reject`)
    await loadEvaluation()
    toast.add({ severity: 'warn', summary: t('message.success'), detail: t('kpi.evaluation_rejected'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    rejecting.value = false
  }
}

async function completeEvaluation() {
  completing.value = true
  try {
    await api.post(`/api/v1/tenant/performance/kpi/evaluations/${evaluationId.value}/complete`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.evaluation_completed'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    completing.value = false
  }
}

onMounted(loadEvaluation)
</script>
