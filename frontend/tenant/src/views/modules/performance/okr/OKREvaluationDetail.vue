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
            <div class="w-14 h-14 rounded-xl bg-purple-50 dark:bg-purple-500/10 flex items-center justify-center">
              <i class="pi pi-bullseye text-2xl text-purple-600 dark:text-purple-400"></i>
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
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('okr.final_score') }}</p>
              <p class="text-3xl font-bold" :class="getScoreClass(evaluation.final_score)">
                {{ evaluation.final_score?.toFixed(1) || '0.0' }}
              </p>
            </div>
            <Tag v-if="evaluation.rating_name" :value="evaluation.rating_name" :severity="getRatingSeverity(evaluation.rating_color)" />
          </div>
        </div>
      </div>

      <!-- Objectives -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-bullseye text-purple-500"></i>
            {{ t('okr.objectives') }}
            <span class="text-xs font-normal text-gray-500">({{ objectiveGroups.length }})</span>
          </h3>
          <Button v-if="canEdit" :label="t('okr.recalculate')" icon="pi pi-refresh" size="small" outlined severity="secondary" :loading="recalculating" @click="recalculate" />
        </div>

        <div
          v-for="group in objectiveGroups"
          :key="group.key"
          class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5"
        >
          <div class="flex items-center justify-between mb-3">
            <div>
              <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ group.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('okr.weight') }}: {{ group.weight?.toFixed(1) }}%</p>
            </div>
            <div class="text-right">
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('okr.achievement') }}</p>
              <p class="text-lg font-bold" :class="getAchievementClass(group.achievement)">{{ group.achievement.toFixed(1) }}%</p>
            </div>
          </div>

          <DataTable :value="group.items" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
            <Column field="key_result_title" :header="t('okr.key_result_title')" style="min-width:200px">
              <template #body="{data}">
                <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.key_result_title }}</span>
              </template>
            </Column>

            <Column field="key_result_weight" :header="t('okr.weight')" style="width:80px">
              <template #body="{data}">
                <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.key_result_weight?.toFixed(1) }}%</span>
              </template>
            </Column>

            <Column field="target_value" :header="t('okr.target')" style="width:100px">
              <template #body="{data}">
                <span class="text-gray-600 dark:text-gray-300 font-mono">{{ formatNumber(data.target_value) }} {{ data.unit }}</span>
              </template>
            </Column>

            <Column field="actual_value" :header="t('okr.actual')" style="width:120px">
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

            <Column field="achievement" :header="t('okr.achievement')" style="width:100px">
              <template #body="{data}">
                <span class="font-mono font-semibold" :class="getAchievementClass(data.achievement)">
                  {{ data.achievement?.toFixed(1) || '0.0' }}%
                </span>
              </template>
            </Column>

            <Column field="score" :header="t('okr.score')" style="width:80px">
              <template #body="{data}">
                <span class="font-mono font-bold" :class="getScoreClass(data.score)">
                  {{ data.score?.toFixed(1) || '0.0' }}
                </span>
              </template>
            </Column>

            <Column style="width:50px">
              <template #body="{data}">
                <Button icon="pi pi-history" size="small" text severity="secondary" v-tooltip.left="t('okr.progress')" @click="openProgress(data)" />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="goBack" />
        <div class="flex items-center gap-2">
          <template v-if="evaluation.status === 'DRAFT'">
            <Button :label="t('okr.submit')" icon="pi pi-send" size="small" :loading="submitting" @click="submitEvaluation" />
          </template>
          <template v-else-if="evaluation.status === 'SUBMITTED'">
            <Button :label="t('okr.reject')" icon="pi pi-times" severity="danger" outlined size="small" :loading="rejecting" @click="rejectEvaluation" />
            <Button :label="t('okr.approve')" icon="pi pi-check" severity="success" size="small" :loading="approving" @click="approveEvaluation" />
          </template>
          <template v-else-if="evaluation.status === 'APPROVED'">
            <Button :label="t('okr.complete')" icon="pi pi-check-circle" severity="success" size="small" :loading="completing" @click="completeEvaluation" />
          </template>
        </div>
      </div>
    </template>

    <!-- Not Found -->
    <div v-else class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-exclamation-circle text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('okr.evaluation_not_found') }}</p>
      <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" class="mt-4" @click="goBack" />
    </div>

    <!-- Progress Check-in Dialog -->
    <Dialog v-model:visible="progressDialogVisible" :header="t('okr.check_in')" modal :style="{ width: '480px' }">
      <div v-if="progressTarget" class="space-y-4">
        <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ progressTarget.key_result_title }}</p>

        <!-- Add Progress Form -->
        <div v-if="canEdit" class="grid grid-cols-2 gap-3 pb-3 border-b border-gray-100 dark:border-gray-700">
          <FormRow :label="t('common.date')">
            <DatePicker v-model="progressForm.progress_date" class="w-full" size="small" />
          </FormRow>
          <FormRow :label="t('okr.actual')">
            <InputNumber v-model="progressForm.actual_value" :minFractionDigits="0" :maxFractionDigits="2" class="w-full" size="small" />
          </FormRow>
          <div class="col-span-2">
            <FormRow :label="t('common.notes')">
              <Textarea v-model="progressForm.notes" rows="2" class="w-full" size="small" />
            </FormRow>
          </div>
          <div class="col-span-2 flex justify-end">
            <Button :label="t('okr.add_progress')" icon="pi pi-plus" size="small" :loading="addingProgress" @click="addProgress" />
          </div>
        </div>

        <!-- Progress History -->
        <div class="max-h-64 overflow-y-auto space-y-2">
          <div v-if="progressLoading" class="text-center py-4">
            <i class="pi pi-spin pi-spinner text-gray-400"></i>
          </div>
          <div v-else-if="progressHistory.length === 0" class="text-center py-4 text-xs text-gray-400">
            {{ t('okr.no_progress') }}
          </div>
          <div
            v-for="p in progressHistory"
            :key="p.id"
            class="flex items-center justify-between text-xs p-2 rounded-lg bg-gray-50 dark:bg-gray-700/50"
          >
            <div>
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ p.progress_date }}</span>
              <span class="text-gray-400 mx-1">-</span>
              <span class="text-gray-600 dark:text-gray-300">{{ formatNumber(p.actual_value) }}</span>
              <span class="ml-2" :class="getAchievementClass(p.achievement)">({{ p.achievement?.toFixed(1) }}%)</span>
              <p v-if="p.notes" class="text-gray-500 dark:text-gray-400 mt-0.5">{{ p.notes }}</p>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="progressDialogVisible = false" />
      </template>
    </Dialog>
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
import Dialog from 'primevue/dialog'
import DatePicker from 'primevue/datepicker'
import Textarea from 'primevue/textarea'
import FormRow from '@/components/FormRow.vue'

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

const progressDialogVisible = ref(false)
const progressTarget = ref(null)
const progressHistory = ref([])
const progressLoading = ref(false)
const addingProgress = ref(false)
const progressForm = ref({ progress_date: new Date(), actual_value: 0, notes: '' })

const evaluationId = computed(() => route.params.id)
const canEdit = computed(() => evaluation.value?.status === 'DRAFT')

const objectiveGroups = computed(() => {
  const groups = {}
  const order = []
  for (const d of details.value) {
    const key = d.objective_id || d.objective_title
    if (!groups[key]) {
      groups[key] = {
        key,
        title: d.objective_title,
        weight: d.objective_weight,
        items: []
      }
      order.push(key)
    }
    groups[key].items.push(d)
  }
  return order.map(key => {
    const g = groups[key]
    const totalScore = g.items.reduce((sum, i) => sum + (i.score || 0), 0)
    const totalWeight = g.items.reduce((sum, i) => sum + (i.key_result_weight || 0), 0) || 1
    const weightedAchievement = g.items.reduce((sum, i) => sum + (i.achievement || 0) * (i.key_result_weight || 0), 0)
    g.achievement = weightedAchievement / totalWeight
    g.score = totalScore
    return g
  })
})

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
  router.push('/performance/okr')
}

async function loadEvaluation() {
  pageLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/details`)
    const data = res.data?.data || res.data
    evaluation.value = {
      id: data.id,
      employee_id: data.employee_id,
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
      objective_id: d.objective_id,
      objective_title: d.objective_title,
      objective_weight: d.objective_weight,
      key_result_id: d.key_result_id,
      key_result_title: d.key_result_title,
      key_result_weight: d.key_result_weight,
      target_value: d.target_value,
      target_type: d.target_type,
      unit: d.unit,
      formula_type: d.formula_type,
      actual_value: d.actual_value,
      achievement: d.achievement,
      score: d.score
    }))
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
    evaluation.value = null
  } finally {
    pageLoading.value = false
  }
}

async function updateActual(detail) {
  try {
    await api.put(`/api/v1/tenant/performance/okr/evaluation-details/${detail.id}`, {
      actual_value: detail.actual_value || 0
    })
    await recalculate()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

async function recalculate() {
  recalculating.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/recalculate`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.score_recalculated'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    recalculating.value = false
  }
}

async function submitEvaluation() {
  submitting.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/submit`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.evaluation_submitted'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    submitting.value = false
  }
}

async function approveEvaluation() {
  approving.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/approve`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.evaluation_approved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    approving.value = false
  }
}

async function rejectEvaluation() {
  rejecting.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/reject`)
    await loadEvaluation()
    toast.add({ severity: 'warn', summary: t('message.success'), detail: t('okr.evaluation_rejected'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    rejecting.value = false
  }
}

async function completeEvaluation() {
  completing.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/complete`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.evaluation_completed'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    completing.value = false
  }
}

function formatDateForApi(date) {
  const d = date instanceof Date ? date : new Date(date)
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

async function openProgress(detail) {
  progressTarget.value = detail
  progressForm.value = { progress_date: new Date(), actual_value: detail.actual_value || 0, notes: '' }
  progressDialogVisible.value = true
  await loadProgress()
}

async function loadProgress() {
  if (!progressTarget.value) return
  progressLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/performance/okr/evaluation-details/${progressTarget.value.id}/progress`)
    progressHistory.value = res.data?.data || []
  } catch {
    progressHistory.value = []
  } finally {
    progressLoading.value = false
  }
}

async function addProgress() {
  if (!progressTarget.value) return
  addingProgress.value = true
  try {
    await api.post('/api/v1/tenant/performance/okr/progress', {
      evaluation_detail_id: progressTarget.value.id,
      progress_date: formatDateForApi(progressForm.value.progress_date),
      actual_value: progressForm.value.actual_value || 0,
      notes: progressForm.value.notes || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.add_progress'), life: 3000 })
    progressForm.value.notes = ''
    await loadProgress()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    addingProgress.value = false
  }
}

onMounted(loadEvaluation)
</script>