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
              <h2 class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ evaluation.employee_name }}</h2>
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

      <!-- Phase hint -->
      <div v-if="phaseHintKey" class="bg-blue-50 dark:bg-blue-500/10 border border-blue-200 dark:border-blue-500/30 rounded-xl p-3 flex items-center gap-2">
        <i class="pi pi-info-circle text-blue-500"></i>
        <span class="text-xs text-blue-700 dark:text-blue-300">{{ t(phaseHintKey) }}</span>
      </div>

      <!-- Objectives / Key Results -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-bullseye text-purple-500"></i>
            {{ t('okr.objectives') }}
            <span class="text-xs font-normal text-gray-500">({{ objectiveGroups.length }})</span>
          </h3>
          <Button v-if="canEditActual" :label="t('okr.recalculate')" icon="pi pi-refresh" size="small" outlined severity="secondary" :loading="recalculating" @click="recalculate" />
        </div>

        <div
          v-for="group in objectiveGroups"
          :key="group.objective_id"
          class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5"
        >
          <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
            <div>
              <p class="text-sm font-semibold text-navy-800 dark:text-gray-100">{{ group.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('okr.weight') }}: {{ group.weight?.toFixed(1) }}%</p>
            </div>
            <div class="flex items-center gap-3">
              <span v-if="canProposeKR" class="text-xs" :class="group.totalKRWeight > 100 ? 'text-red-600 dark:text-red-400' : 'text-gray-500 dark:text-gray-400'">
                {{ t('okr.total_weight') }}: {{ group.totalKRWeight.toFixed(2) }}% / 100%
              </span>
              <div class="text-right" v-if="showActualColumn">
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('okr.achievement') }}</p>
                <p class="text-lg font-bold" :class="getAchievementClass(group.achievement)">{{ group.achievement.toFixed(1) }}%</p>
              </div>
              <Button v-if="canProposeKR" :label="t('okr.add_key_result')" icon="pi pi-plus" size="small" outlined @click="openKeyResultDialog(group)" />
            </div>
          </div>

          <p v-if="group.items.length === 0" class="text-xs text-gray-400 dark:text-gray-500 py-2">
            {{ t('okr.no_key_results') }}
          </p>

          <DataTable v-else :value="group.items" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
            <Column field="key_result_title" :header="t('okr.key_result_title')" style="min-width:200px">
              <template #body="{data}">
                <InputText
                  v-if="canProposeKR"
                  v-model="data.key_result_title"
                  class="w-full !text-xs"
                  size="small"
                />
                <span v-else class="text-navy-800 dark:text-gray-100 font-medium">{{ data.key_result_title }}</span>
              </template>
            </Column>

            <Column field="formula_type" :header="t('okr.formula_type')" style="width:140px">
              <template #body="{data}">
                <Select
                  v-if="canProposeKR"
                  v-model="data.formula_type"
                  :options="formulaOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-full !text-xs"
                  size="small"
                />
                <Tag v-else :value="data.formula_type" severity="secondary" class="!text-xs" />
              </template>
            </Column>

            <Column field="key_result_weight" :header="t('okr.weight')" style="width:100px">
              <template #body="{data}">
                <InputNumber
                  v-if="canProposeKR"
                  v-model="data.key_result_weight"
                  :min="0"
                  :max="100"
                  :minFractionDigits="0"
                  :maxFractionDigits="2"
                  suffix="%"
                  class="w-full !text-xs"
                  size="small"
                />
                <span v-else class="text-gray-600 dark:text-gray-300 font-mono">{{ data.key_result_weight?.toFixed(1) }}%</span>
              </template>
            </Column>

            <Column field="target_value" :header="t('okr.target')" style="width:110px">
              <template #body="{data}">
                <InputNumber
                  v-if="canProposeKR"
                  v-model="data.target_value"
                  :minFractionDigits="0"
                  :maxFractionDigits="2"
                  class="w-full !text-xs"
                  size="small"
                />
                <span v-else class="text-gray-600 dark:text-gray-300 font-mono">{{ formatNumber(data.target_value) }} {{ data.unit }}</span>
              </template>
            </Column>

            <Column field="unit" :header="t('okr.unit')" style="width:90px">
              <template #body="{data}">
                <InputText
                  v-if="canProposeKR"
                  v-model="data.unit"
                  class="w-full !text-xs"
                  size="small"
                  maxlength="50"
                />
                <span v-else class="text-gray-500 dark:text-gray-400 text-xs">{{ data.unit || '-' }}</span>
              </template>
            </Column>

            <Column field="actual_value" :header="t('okr.actual')" style="width:120px">
              <template #body="{data}">
                <InputNumber
                  v-if="canEditActual"
                  v-model="data.actual_value"
                  :minFractionDigits="0"
                  :maxFractionDigits="2"
                  class="w-full !text-xs"
                  size="small"
                />
                <span v-else-if="showActualColumn" class="text-navy-800 dark:text-gray-100 font-mono font-semibold">{{ formatNumber(data.actual_value) }}</span>
                <span v-else class="text-gray-300 dark:text-gray-600">—</span>
              </template>
            </Column>

            <Column field="achievement" :header="t('okr.achievement')" style="width:90px">
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

            <Column v-if="canProposeKR" style="width:50px">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" @click="removeKeyResult(data)" />
              </template>
            </Column>
            <Column v-else style="width:50px">
              <template #body="{data}">
                <Button icon="pi pi-history" size="small" text severity="secondary" v-tooltip.left="t('okr.progress')" @click="openProgress(data)" />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- Save Target -->
      <div v-if="canProposeKR" class="flex justify-end">
        <Button :label="t('okr.save_target')" icon="pi pi-save" size="small" :loading="savingTargets" @click="saveAllKeyResultTargets" />
      </div>

      <!-- Save Actual -->
      <div v-if="canEditActual" class="flex justify-end">
        <Button :label="t('okr.save_actual')" icon="pi pi-save" size="small" :loading="savingActuals" @click="saveAllActuals" />
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="goBack" />
        <div class="flex items-center gap-2">
          <template v-if="evaluation.status === 'DRAFT'">
            <Button :label="t('okr.submit_key_results')" icon="pi pi-send" size="small" :loading="submittingKR" @click="submitKeyResults" />
          </template>
          <template v-else-if="evaluation.status === 'KR_SUBMITTED'">
            <span v-if="evaluation.kr_approval_instance_id" class="text-xs text-gray-400 dark:text-gray-500 flex items-center gap-1.5">
              <i class="pi pi-clock"></i> {{ t('okr.awaiting_central_approval') }}
            </span>
            <template v-else>
              <Button :label="t('okr.reject_key_results')" icon="pi pi-times" severity="danger" outlined size="small" :loading="rejectingKR" @click="rejectKeyResults" />
              <Button :label="t('okr.approve_key_results')" icon="pi pi-check" severity="success" size="small" :loading="approvingKR" @click="approveKeyResults" />
            </template>
          </template>
          <template v-else-if="evaluation.status === 'KR_APPROVED'">
            <Button :label="t('okr.submit')" icon="pi pi-send" size="small" :loading="submitting" @click="submitEvaluation" />
          </template>
          <template v-else-if="evaluation.status === 'SUBMITTED'">
            <span v-if="evaluation.assessment_approval_instance_id" class="text-xs text-gray-400 dark:text-gray-500 flex items-center gap-1.5">
              <i class="pi pi-clock"></i> {{ t('okr.awaiting_central_approval') }}
            </span>
            <template v-else>
              <Button :label="t('okr.reject')" icon="pi pi-times" severity="danger" outlined size="small" :loading="rejecting" @click="rejectEvaluation" />
              <Button :label="t('okr.approve')" icon="pi pi-check" severity="success" size="small" :loading="approving" @click="approveEvaluation" />
            </template>
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

    <!-- Add Key Result Dialog -->
    <Dialog v-model:visible="keyResultDialogVisible" :header="t('okr.add_key_result')" modal :style="{ width: '480px' }">
      <div v-if="keyResultTargetObjective" class="space-y-4">
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ keyResultTargetObjective.title }}</p>
        <FormRow :label="t('okr.key_result_title')" required :errors="keyResultErrors?.title">
          <TextInput v-model="keyResultForm.title" maxlength="255" autofocus :placeholder="t('okr.key_result_title_placeholder')" :class="{'p-invalid':keyResultErrors?.title}" />
        </FormRow>
        <FormRow :label="t('okr.weight')" required :errors="keyResultErrors?.weight">
          <InputNumber v-model="keyResultForm.weight" class="!w-full" :min="0" :max="100" :minFractionDigits="0" :maxFractionDigits="2" suffix="%" size="small" />
        </FormRow>
        <FormRow :label="t('okr.target')" required :errors="keyResultErrors?.target_value">
          <InputNumber v-model="keyResultForm.target_value" class="!w-full" :minFractionDigits="0" :maxFractionDigits="2" size="small" />
        </FormRow>
        <FormRow :label="t('okr.unit')">
          <TextInput v-model="keyResultForm.unit" />
        </FormRow>
        <FormRow :label="t('okr.formula_type')">
          <Select v-model="keyResultForm.formula_type" :options="formulaOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="keyResultDialogVisible=false" />
          <Button :label="t('common.save')" size="small" :loading="savingKeyResult" :disabled="savingKeyResult" @click="saveKeyResult" />
        </div>
      </template>
    </Dialog>

    <!-- Progress Check-in Dialog -->
    <Dialog v-model:visible="progressDialogVisible" :header="t('okr.check_in')" modal :style="{ width: '480px' }">
      <div v-if="progressTarget" class="space-y-4">
        <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ progressTarget.key_result_title }}</p>

        <div class="grid grid-cols-2 gap-3 pb-3 border-b border-gray-100 dark:border-gray-700">
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
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ formatDate(p.progress_date, locale) }}</span>
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
import { getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import DatePicker from 'primevue/datepicker'
import Textarea from 'primevue/textarea'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const { t, locale } = useI18n()

const pageLoading = ref(true)
const evaluation = ref(null)
const details = ref([])
const objectives = ref([])
const recalculating = ref(false)
const submittingKR = ref(false)
const approvingKR = ref(false)
const rejectingKR = ref(false)
const submitting = ref(false)
const approving = ref(false)
const rejecting = ref(false)
const completing = ref(false)
const savingTargets = ref(false)
const savingActuals = ref(false)

const keyResultDialogVisible = ref(false)
const keyResultTargetObjective = ref(null)
const savingKeyResult = ref(false)
const keyResultErrors = ref({})
const keyResultForm = ref({ title: '', weight: 0, target_value: 0, unit: '', formula_type: 'HIGHER_BETTER' })

const progressDialogVisible = ref(false)
const progressTarget = ref(null)
const progressHistory = ref([])
const progressLoading = ref(false)
const addingProgress = ref(false)
const progressForm = ref({ progress_date: new Date(), actual_value: 0, notes: '' })

const formulaOptions = [
  { label: 'Higher Better', value: 'HIGHER_BETTER' },
  { label: 'Lower Better', value: 'LOWER_BETTER' }
]

const evaluationId = computed(() => route.params.id)
// Two-phase gating: Key Results proposed only in DRAFT ("Ajukan Key Result"),
// actual editable only once Key Results have been approved ("OKR Active").
const canProposeKR = computed(() => evaluation.value?.status === 'DRAFT')
const canEditActual = computed(() => evaluation.value?.status === 'KR_APPROVED')
const showActualColumn = computed(() => !['DRAFT', 'KR_SUBMITTED'].includes(evaluation.value?.status))

const phaseHintKey = computed(() => {
  switch (evaluation.value?.status) {
    case 'DRAFT': return 'okr.phase_hint_draft'
    case 'KR_SUBMITTED': return 'okr.phase_hint_kr_submitted'
    case 'KR_APPROVED': return 'okr.phase_hint_kr_approved'
    default: return null
  }
})

const objectiveGroups = computed(() => {
  return objectives.value.map(obj => {
    const items = details.value.filter(d => d.objective_id === obj.id)
    const totalScore = items.reduce((sum, i) => sum + (i.score || 0), 0)
    const totalKRWeight = items.reduce((sum, i) => sum + (i.key_result_weight || 0), 0)
    const totalItemWeight = items.reduce((sum, i) => sum + (i.key_result_weight || 0), 0) || 1
    const weightedAchievement = items.reduce((sum, i) => sum + (i.achievement || 0) * (i.key_result_weight || 0), 0)
    return {
      objective_id: obj.id,
      title: obj.title,
      weight: obj.weight,
      items,
      score: totalScore,
      totalKRWeight,
      achievement: weightedAchievement / totalItemWeight
    }
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
    case 'KR_APPROVED': return 'info'
    case 'KR_SUBMITTED': return 'warn'
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
      template_id: data.template_id,
      employee_name: data.employee_name || data.employee?.full_name || '-',
      organization_name: data.organization_name || data.organization?.name || '-',
      period_code: data.period_code || data.period?.period_code || '-',
      status: data.status,
      final_score: data.final_score,
      rating_name: data.rating_name || data.rating?.name,
      rating_color: data.rating_color || data.rating?.color,
      kr_approval_instance_id: data.kr_approval_instance_id || null,
      assessment_approval_instance_id: data.assessment_approval_instance_id || null
    }
    details.value = (data.details || []).map(d => ({
      id: d.id,
      objective_id: d.objective_id,
      objective_title: d.objective_title,
      objective_weight: d.objective_weight,
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

  if (evaluation.value?.template_id) {
    await loadObjectives()
  }
}

async function loadObjectives() {
  try {
    const res = await api.get(`/api/v1/tenant/performance/okr/templates/${evaluation.value.template_id}/objectives`)
    objectives.value = res.data?.data || []
  } catch {
    objectives.value = []
  }
}

function openKeyResultDialog(group) {
  keyResultTargetObjective.value = group
  keyResultForm.value = { title: '', weight: 0, target_value: 0, unit: '', formula_type: 'HIGHER_BETTER' }
  keyResultErrors.value = {}
  keyResultDialogVisible.value = true
}

async function saveKeyResult() {
  keyResultErrors.value = {}
  if (!keyResultForm.value.title?.trim()) {
    keyResultErrors.value = { title: [t('form.required')] }
    return
  }
  if (!keyResultForm.value.weight) {
    keyResultErrors.value = { weight: [t('form.required')] }
    return
  }
  if (!keyResultForm.value.target_value) {
    keyResultErrors.value = { target_value: [t('form.required')] }
    return
  }
  if (keyResultTargetObjective.value.totalKRWeight + keyResultForm.value.weight > 100) {
    keyResultErrors.value = { weight: [t('okr.key_result_weight_exceeds_100')] }
    return
  }

  savingKeyResult.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/key-results`, {
      evaluation_id: evaluationId.value,
      objective_id: keyResultTargetObjective.value.objective_id,
      objective_title: keyResultTargetObjective.value.title,
      objective_weight: keyResultTargetObjective.value.weight,
      title: keyResultForm.value.title,
      weight: keyResultForm.value.weight,
      target_value: keyResultForm.value.target_value,
      unit: keyResultForm.value.unit || null,
      formula_type: keyResultForm.value.formula_type
    })
    keyResultDialogVisible.value = false
    await loadEvaluation()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      keyResultErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    savingKeyResult.value = false
  }
}

async function removeKeyResult(item) {
  try {
    await api.delete(`/api/v1/tenant/performance/okr/evaluation-key-results/${item.id}`)
    await loadEvaluation()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

async function saveAllKeyResultTargets() {
  const overWeight = objectiveGroups.value.some(g => g.totalKRWeight > 100)
  if (overWeight) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('okr.key_result_weight_exceeds_100'), life: 4000 })
    return
  }
  savingTargets.value = true
  try {
    await Promise.all(details.value.map(d => api.put(`/api/v1/tenant/performance/okr/evaluation-key-results/${d.id}/target`, {
      title: d.key_result_title,
      target_type: d.target_type,
      target_value: d.target_value || 0,
      unit: d.unit || null,
      formula_type: d.formula_type,
      weight: d.key_result_weight || 0
    })))
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.target_saved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    savingTargets.value = false
  }
}

async function saveAllActuals() {
  savingActuals.value = true
  try {
    await Promise.all(details.value.map(d => api.put(`/api/v1/tenant/performance/okr/evaluation-details/${d.id}`, {
      actual_value: d.actual_value || 0
    })))
    await recalculate()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.actual_saved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    savingActuals.value = false
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

async function submitKeyResults() {
  submittingKR.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/submit-key-results`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.key_results_submitted'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    submittingKR.value = false
  }
}

async function approveKeyResults() {
  approvingKR.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/approve-key-results`)
    await loadEvaluation()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('okr.key_results_approved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    approvingKR.value = false
  }
}

async function rejectKeyResults() {
  rejectingKR.value = true
  try {
    await api.post(`/api/v1/tenant/performance/okr/evaluations/${evaluationId.value}/reject-key-results`)
    await loadEvaluation()
    toast.add({ severity: 'warn', summary: t('message.success'), detail: t('okr.key_results_rejected'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    rejectingKR.value = false
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
