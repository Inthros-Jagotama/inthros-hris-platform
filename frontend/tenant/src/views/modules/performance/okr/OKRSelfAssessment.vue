<template>
  <div class="space-y-4 max-w-2xl">
    <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <div class="flex items-center gap-2 mb-2">
        <i class="pi pi-bullseye text-rose-500"></i>
        <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('okr.self_assessment_title') }}</h3>
      </div>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('okr.self_assessment_intro') }}</p>
    </div>

    <div v-if="loading" class="space-y-3">
      <div class="h-24 bg-gray-200 dark:bg-gray-700 rounded-xl animate-pulse"></div>
    </div>

    <!-- No active employment/position -->
    <div v-else-if="!context.has_position" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-8 text-center">
      <i class="pi pi-exclamation-circle text-3xl text-amber-400 mb-3"></i>
      <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('okr.no_position_title') }}</p>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ t('okr.no_position_hint') }}</p>
    </div>

    <!-- No active template for this position -->
    <div v-else-if="context.templates.length === 0" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-8 text-center">
      <i class="pi pi-file-excel text-3xl text-gray-300 dark:text-gray-600 mb-3"></i>
      <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('okr.no_template_title') }}</p>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">{{ t('okr.no_template_hint', { organization: context.organization_name }) }}</p>
    </div>

    <!-- Ready to fill -->
    <div v-else class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5 space-y-4">
      <div>
        <p class="text-xs text-gray-400">{{ t('okr.your_position') }}</p>
        <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ context.organization_name }}</p>
      </div>

      <FormRow :label="t('okr.period')" required>
        <Select v-model="selectedPeriodId" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('okr.select_period')" class="w-full" @change="checkExistingEvaluation" />
      </FormRow>

      <div v-for="tpl in context.templates" :key="tpl.id" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 flex items-center justify-between gap-3">
        <div class="min-w-0">
          <p class="text-sm font-medium text-gray-800 dark:text-gray-100 truncate">{{ tpl.name }}</p>
        </div>
        <Button
          :label="existingEvaluationId ? t('okr.continue_self_assessment') : t('okr.start_self_assessment')"
          icon="pi pi-arrow-right"
          iconPos="right"
          size="small"
          :loading="starting"
          @click="startAssessment(tpl)"
        />
      </div>

      <p v-if="existingEvaluationId" class="text-xs text-emerald-600 dark:text-emerald-400 flex items-center gap-1.5">
        <i class="pi pi-check-circle"></i> {{ t('okr.existing_evaluation_found') }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import Select from 'primevue/select'
import FormRow from '@/components/FormRow.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

const loading = ref(true)
const starting = ref(false)
const context = ref({ has_position: false, templates: [] })

const periods = ref([])
const selectedPeriodId = ref(null)
const existingEvaluationId = ref(null)

const periodOptions = computed(() => periods.value.map(p => ({ label: p.period_code, value: p.id })))

async function loadContext() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/okr/my-context')
    context.value = res.data?.data || { has_position: false, templates: [] }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadPeriods() {
  try {
    const res = await api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } })
    const list = res.data?.data || []
    periods.value = list
    if (list.length > 0) {
      selectedPeriodId.value = list[0].id
      await checkExistingEvaluation()
    }
  } catch {
    periods.value = []
  }
}

async function checkExistingEvaluation() {
  existingEvaluationId.value = null
  if (!selectedPeriodId.value || !context.value.employee_id) return
  try {
    const res = await api.get('/api/v1/tenant/performance/okr/evaluations', {
      params: { employee_id: context.value.employee_id, period_id: selectedPeriodId.value, per_page: 50 }
    })
    const list = res.data?.data || []
    const templateIds = new Set(context.value.templates.map(t => t.id))
    const match = list.find(e => templateIds.has(e.template_id))
    if (match) existingEvaluationId.value = match.id
  } catch {
    existingEvaluationId.value = null
  }
}

async function startAssessment(template) {
  if (!selectedPeriodId.value) {
    toast.add({ severity: 'warn', summary: t('message.error'), detail: t('okr.period_required'), life: 3000 })
    return
  }

  if (existingEvaluationId.value) {
    router.push(`/performance/okr/evaluation/${existingEvaluationId.value}`)
    return
  }

  starting.value = true
  try {
    const res = await api.post('/api/v1/tenant/performance/okr/evaluations', {
      employee_id: context.value.employee_id,
      organization_id: context.value.organization_id,
      period_id: selectedPeriodId.value,
      template_id: template.id
    })
    const evalId = res.data?.data?.id || res.data?.id
    if (evalId) {
      router.push(`/performance/okr/evaluation/${evalId}`)
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    starting.value = false
  }
}

onMounted(async () => {
  await loadContext()
  if (context.value.has_position && context.value.templates.length > 0) {
    await loadPeriods()
  }
})
</script>
