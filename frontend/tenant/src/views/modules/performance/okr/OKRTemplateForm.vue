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

    <!-- Blocked: not yet eligible to create an Objective -->
    <div v-else-if="!isEditing && scopeLoaded && !scope?.eligible" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-lock text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('okr.objective_scope_ineligible_title') }}</p>
      <p class="text-xs mt-1 max-w-md text-center">{{ t(scope?.ineligible_reason_key || 'okr.objective_scope_ineligible_no_own_objective') }}</p>
      <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" class="mt-4" @click="goBack" />
    </div>

    <!-- Blocked: eligible but no subordinate organizations to create for -->
    <div v-else-if="!isEditing && scopeLoaded && scope?.eligible && organizationOptions.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-sitemap text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ t('okr.no_subordinate_organizations') }}</p>
      <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" class="mt-4" @click="goBack" />
    </div>

    <template v-else>
      <!-- Template Info Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 mb-4 flex items-center gap-2">
          <i class="pi pi-file text-emerald-500"></i>
          {{ t('okr.template_info') }}
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormRow :label="t('okr.template_name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" :placeholder="t('okr.template_name_placeholder')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
          <FormRow :label="t('okr.organization')" :errors="errors?.organization_id">
            <Select v-model="form.organization_id" small :options="organizationOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" filter showClear />
          </FormRow>
          <FormRow :label="t('okr.period')" :errors="errors?.period_id">
            <Select v-model="form.period_id" small :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" showClear />
          </FormRow>
          <FormRow :label="t('okr.status')" :errors="errors?.status">
            <Select v-model="form.status" small :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
          </FormRow>
        </div>
        <div class="mt-4">
          <FormRow :label="t('okr.description_label')" :errors="errors?.description">
            <Textarea v-model="form.description" small rows="2" :placeholder="t('okr.description_placeholder')" class="w-full" />
          </FormRow>
        </div>
      </div>

      <!-- Objectives Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 flex items-center gap-2">
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
        <p class="text-xs text-gray-400 dark:text-gray-500 -mt-2 mb-4">{{ t('okr.key_results_filled_at_evaluation_hint') }}</p>

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
            <div class="flex items-start gap-3">
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
const scope = ref(null)
const scopeLoaded = ref(false)

const statusOptions = [
  { label: 'Draft', value: 0 },
  { label: 'Active', value: 1 },
  { label: 'Inactive', value: 2 }
]

const objectiveTotalWeight = computed(() => {
  return objectives.value.reduce((sum, o) => sum + (o.weight || 0), 0)
})

function addObjective() {
  objectives.value.push({
    id: null,
    _key: `new-${++keyCounter}`,
    title: '',
    description: '',
    weight: 0,
    sort_order: objectives.value.length
  })
}

function removeObjective(index) {
  objectives.value.splice(index, 1)
}

function goBack() {
  router.push('/performance/okr/templates')
}

async function loadReferenceData() {
  try {
    const periodRes = await api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } })
    periodOptions.value = (periodRes.data?.data || []).map(p => ({
      label: `${p.period_code} (${p.year})`,
      value: p.id
    }))

    if (isEditing.value) {
      // Editing an existing Objective/Template keeps the full organization
      // list — the cascading-eligibility gate only governs creating new ones.
      const orgRes = await api.get('/api/v1/tenant/organizations', { params: { per_page: 200, active_only: true } })
      organizationOptions.value = (orgRes.data?.data || []).map(o => ({
        label: o.nomenclature || o.name || o.code,
        value: o.id
      }))
    } else {
      const scopeRes = await api.get('/api/v1/tenant/performance/okr/templates/objective-scope')
      scope.value = scopeRes.data?.data || scopeRes.data
      organizationOptions.value = (scope.value?.subordinate_organizations || []).map(o => ({
        label: o.name,
        value: o.id
      }))
      scopeLoaded.value = true
    }
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
      sort_order: obj.sort_order || 0
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

    // Save objectives (Key Results are filled in by the employee at
    // evaluation time, not authored here — see the two-phase OKR flow).
    for (let i = 0; i < objectives.value.length; i++) {
      const objective = objectives.value[i]
      const objPayload = {
        template_id: savedTemplateId,
        title: objective.title,
        description: objective.description || null,
        weight: objective.weight,
        sort_order: i
      }

      if (objective.id) {
        await api.put(`/api/v1/tenant/performance/okr/objectives/${objective.id}`, objPayload)
      } else {
        await api.post('/api/v1/tenant/performance/okr/objectives', objPayload)
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