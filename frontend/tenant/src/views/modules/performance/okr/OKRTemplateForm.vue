<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/3"></div>
      <div class="grid grid-cols-2 gap-4">
        <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
      <div class="h-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
    </div>

    <template v-else>
      <!-- Template Info Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 mb-4 flex items-center gap-2">
          <i class="pi pi-file text-emerald-500"></i>
          {{ t('okr.template_info') }}
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormRow :label="t('okr.template_name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" :placeholder="t('okr.template_name_placeholder')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
          <FormRow :label="t('okr.organization')" :errors="errors?.organization_id">
            <Select v-model="form.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" filter showClear />
          </FormRow>
          <FormRow :label="t('okr.period')" :errors="errors?.period_id">
            <Select v-model="form.period_id" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" showClear />
          </FormRow>
          <FormRow :label="t('okr.status')" :errors="errors?.status">
            <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
          </FormRow>
        </div>
        <div class="mt-4">
          <FormRow :label="t('okr.description_label')" :errors="errors?.description">
            <Textarea v-model="form.description" rows="2" :placeholder="t('okr.description_placeholder')" class="w-full" />
          </FormRow>
        </div>
      </div>

      <!-- Objectives Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-bullseye text-purple-500"></i>
            {{ t('okr.objectives') }}
            <span class="text-xs font-normal text-gray-500">({{ objectives.length }})</span>
          </h3>
          <div class="flex items-center gap-2">
            <span class="text-xs" :class="objectiveTotalWeight === 100 ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'">
              {{ t('okr.total_weight') }}: {{ objectiveTotalWeight.toFixed(2) }}%
            </span>
            <Button :label="t('okr.add_objective')" icon="pi pi-plus" size="small" outlined @click="addObjective" />
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="objectives.length === 0" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-lg">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('okr.no_objectives') }}</p>
          <p class="text-xs mt-1">{{ t('okr.no_objectives_hint') }}</p>
          <Button :label="t('okr.add_first_objective')" icon="pi pi-plus" size="small" class="mt-3" @click="addObjective" />
        </div>

        <!-- Objectives List -->
        <div v-else class="space-y-4">
          <div
            v-for="(objective, oIndex) in objectives"
            :key="objective._key"
            class="border border-gray-200 dark:border-gray-700 rounded-lg p-4"
          >
            <div class="flex items-start gap-3 mb-3">
              <div class="flex-1 grid grid-cols-1 md:grid-cols-12 gap-3">
                <div class="md:col-span-7">
                  <InputText v-model="objective.title" :placeholder="t('okr.objective_title_placeholder')" class="w-full !text-sm" />
                </div>
                <div class="md:col-span-3">
                  <InputNumber v-model="objective.weight" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" suffix="%" class="w-full !text-sm" :placeholder="t('okr.weight')" />
                </div>
                <div class="md:col-span-2 flex justify-end">
                  <Button icon="pi pi-trash" size="small" text severity="danger" @click="removeObjective(oIndex)" />
                </div>
              </div>
            </div>

            <!-- Key Results Table -->
            <div class="pl-4 border-l-2 border-gray-100 dark:border-gray-700">
              <div class="flex items-center justify-between mb-2">
                <span class="text-xs font-medium text-gray-600 dark:text-gray-300 flex items-center gap-1">
                  <i class="pi pi-list-check text-xs"></i>
                  {{ t('okr.key_results') }}
                  <span class="text-gray-400">({{ objective.key_results.length }})</span>
                </span>
                <div class="flex items-center gap-2">
                  <span
                    v-if="objective.key_results.length > 0"
                    class="text-xs"
                    :class="keyResultTotalWeight(objective) === 100 ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'"
                  >
                    {{ t('okr.total_weight') }}: {{ keyResultTotalWeight(objective).toFixed(2) }}%
                  </span>
                  <Button :label="t('okr.add_key_result')" icon="pi pi-plus" size="small" text @click="addKeyResult(objective)" />
                </div>
              </div>

              <p v-if="objective.key_results.length === 0" class="text-xs text-gray-400 dark:text-gray-500 py-2">
                {{ t('okr.no_key_results') }}
              </p>

              <DataTable v-else :value="objective.key_results" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
                <Column :header="t('okr.key_result_title')" style="min-width:180px">
                  <template #body="{data}">
                    <InputText v-model="data.title" :placeholder="t('okr.key_result_title_placeholder')" class="w-full !text-xs" size="small" />
                  </template>
                </Column>
                <Column :header="t('okr.target_type')" style="width:130px">
                  <template #body="{data}">
                    <Select v-model="data.target_type" :options="targetTypeOptions" optionLabel="label" optionValue="value" class="w-full !text-xs" size="small" />
                  </template>
                </Column>
                <Column :header="t('okr.target')" style="width:100px">
                  <template #body="{data}">
                    <InputNumber v-model="data.target_value" :minFractionDigits="0" :maxFractionDigits="2" class="w-full !text-xs" size="small" />
                  </template>
                </Column>
                <Column :header="t('okr.unit')" style="width:80px">
                  <template #body="{data}">
                    <InputText v-model="data.unit" :placeholder="t('okr.unit_placeholder')" class="w-full !text-xs" size="small" maxlength="50" />
                  </template>
                </Column>
                <Column :header="t('okr.formula_type')" style="width:130px">
                  <template #body="{data}">
                    <Select v-model="data.formula_type" :options="formulaOptions" optionLabel="label" optionValue="value" class="w-full !text-xs" size="small" />
                  </template>
                </Column>
                <Column :header="t('okr.weight')" style="width:90px">
                  <template #body="{data}">
                    <InputNumber v-model="data.weight" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" suffix="%" class="w-full !text-xs" size="small" />
                  </template>
                </Column>
                <Column style="width:50px">
                  <template #body="{index}">
                    <Button icon="pi pi-trash" size="small" text severity="danger" @click="removeKeyResult(objective, index)" />
                  </template>
                </Column>
              </DataTable>
            </div>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="goBack" />
        <div class="flex items-center gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="goBack" />
          <Button :label="isEditing ? t('common.update') : t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const { t } = useI18n()

const pageLoading = ref(true)
const saving = ref(false)
const errors = ref({})

const templateId = computed(() => route.params.id)
const isEditing = computed(() => !!templateId.value && templateId.value !== 'new')

const form = ref({
  name: '',
  description: '',
  organization_id: null,
  period_id: null,
  status: 0
})

const objectives = ref([])
let keyCounter = 0

const organizationOptions = ref([])
const periodOptions = ref([])

const statusOptions = [
  { label: 'Draft', value: 0 },
  { label: 'Active', value: 1 },
  { label: 'Inactive', value: 2 }
]

const targetTypeOptions = [
  { label: 'Number', value: 'NUMBER' },
  { label: 'Currency', value: 'CURRENCY' },
  { label: 'Percentage', value: 'PERCENTAGE' },
  { label: 'Duration', value: 'DURATION' },
  { label: 'Boolean', value: 'BOOLEAN' }
]

const formulaOptions = [
  { label: 'Higher Better', value: 'HIGHER_BETTER' },
  { label: 'Lower Better', value: 'LOWER_BETTER' },
  { label: 'Manual', value: 'MANUAL' },
  { label: 'Boolean', value: 'BOOLEAN' },
  { label: 'Percentage', value: 'PERCENTAGE' }
]

const objectiveTotalWeight = computed(() => {
  return objectives.value.reduce((sum, o) => sum + (o.weight || 0), 0)
})

function keyResultTotalWeight(objective) {
  return objective.key_results.reduce((sum, kr) => sum + (kr.weight || 0), 0)
}

function addObjective() {
  objectives.value.push({
    id: null,
    _key: `new-${++keyCounter}`,
    title: '',
    description: '',
    weight: 0,
    sort_order: objectives.value.length,
    key_results: []
  })
}

function removeObjective(index) {
  objectives.value.splice(index, 1)
}

function addKeyResult(objective) {
  objective.key_results.push({
    id: null,
    title: '',
    description: '',
    target_type: 'NUMBER',
    target_value: 0,
    unit: '',
    formula_type: 'HIGHER_BETTER',
    weight: 0,
    minimum_score: 0,
    maximum_score: 100,
    is_required: true,
    sort_order: objective.key_results.length
  })
}

function removeKeyResult(objective, index) {
  objective.key_results.splice(index, 1)
}

function goBack() {
  router.push('/performance/okr/templates')
}

async function loadReferenceData() {
  try {
    const [orgRes, periodRes] = await Promise.all([
      api.get('/api/v1/tenant/organizations', { params: { per_page: 200, active_only: true } }),
      api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } })
    ])

    organizationOptions.value = (orgRes.data?.data || []).map(o => ({
      label: o.nomenclature || o.name || o.code,
      value: o.id
    }))

    periodOptions.value = (periodRes.data?.data || []).map(p => ({
      label: `${p.period_code} (${p.year})`,
      value: p.id
    }))
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 4000 })
  }
}

async function loadTemplate() {
  if (!isEditing.value) return

  try {
    const res = await api.get(`/api/v1/tenant/performance/okr/templates/${templateId.value}`)
    const data = res.data?.data || res.data

    form.value = {
      name: data.name || '',
      description: data.description || '',
      organization_id: data.organization_id || null,
      period_id: data.period_id || null,
      status: data.status ?? 0
    }

    objectives.value = (data.objectives || []).map(obj => ({
      id: obj.id,
      _key: `existing-${obj.id}`,
      title: obj.title || '',
      description: obj.description || '',
      weight: obj.weight || 0,
      sort_order: obj.sort_order || 0,
      key_results: (obj.key_results || []).map(kr => ({
        id: kr.id,
        title: kr.title || '',
        description: kr.description || '',
        target_type: kr.target_type || 'NUMBER',
        target_value: kr.target_value || 0,
        unit: kr.unit || '',
        formula_type: kr.formula_type || 'HIGHER_BETTER',
        weight: kr.weight || 0,
        minimum_score: kr.minimum_score || 0,
        maximum_score: kr.maximum_score || 100,
        is_required: kr.is_required ?? true,
        sort_order: kr.sort_order || 0
      }))
    }))
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  }
}

async function handleSave() {
  errors.value = {}

  if (!form.value.name?.trim()) {
    errors.value = { name: [t('form.required')] }
    return
  }

  if (objectiveTotalWeight.value !== 100 && objectives.value.length > 0) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('okr.objective_weight_total'), life: 4000 })
    return
  }

  for (const objective of objectives.value) {
    const total = keyResultTotalWeight(objective)
    if (objective.key_results.length > 0 && total !== 100) {
      toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('okr.key_result_weight_total'), life: 4000 })
      return
    }
  }

  saving.value = true
  try {
    const payload = {
      name: form.value.name,
      description: form.value.description || null,
      organization_id: form.value.organization_id || null,
      period_id: form.value.period_id || null,
      status: form.value.status
    }

    let savedTemplateId = templateId.value

    if (isEditing.value) {
      await api.put(`/api/v1/tenant/performance/okr/templates/${templateId.value}`, payload)
    } else {
      const res = await api.post('/api/v1/tenant/performance/okr/templates', payload)
      savedTemplateId = res.data?.data?.id || res.data?.id
    }

    // Save objectives and their key results
    for (let i = 0; i < objectives.value.length; i++) {
      const objective = objectives.value[i]
      const objPayload = {
        template_id: savedTemplateId,
        title: objective.title,
        description: objective.description || null,
        weight: objective.weight,
        sort_order: i
      }

      let savedObjectiveId = objective.id
      if (objective.id) {
        await api.put(`/api/v1/tenant/performance/okr/objectives/${objective.id}`, objPayload)
      } else {
        const objRes = await api.post('/api/v1/tenant/performance/okr/objectives', objPayload)
        savedObjectiveId = objRes.data?.data?.id || objRes.data?.id
      }

      for (let j = 0; j < objective.key_results.length; j++) {
        const kr = objective.key_results[j]
        const krPayload = {
          objective_id: savedObjectiveId,
          title: kr.title,
          description: kr.description || null,
          target_type: kr.target_type,
          target_value: kr.target_value,
          unit: kr.unit || null,
          formula_type: kr.formula_type,
          weight: kr.weight,
          minimum_score: kr.minimum_score,
          maximum_score: kr.maximum_score,
          is_required: kr.is_required ?? true,
          sort_order: j
        }

        if (kr.id) {
          await api.put(`/api/v1/tenant/performance/okr/key-results/${kr.id}`, krPayload)
        } else {
          await api.post('/api/v1/tenant/performance/okr/key-results', krPayload)
        }
      }
    }

    toast.add({
      severity: 'success',
      summary: t('message.success'),
      detail: isEditing.value ? t('okr.template_updated') : t('okr.template_created'),
      life: 3000
    })

    router.push('/performance/okr/templates')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    await loadReferenceData()
    await loadTemplate()
  } finally {
    pageLoading.value = false
  }
})
</script>