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
          {{ t('kpi.template_info') }}
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormRow :label="t('kpi.template_name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" :placeholder="t('kpi.template_name_placeholder')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
          <FormRow :label="t('kpi.organization')" :errors="errors?.organization_id">
            <Select v-model="form.organization_id" :options="organizationOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" filter showClear />
          </FormRow>
          <FormRow :label="t('kpi.period')" :errors="errors?.period_id">
            <Select v-model="form.period_id" :options="periodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" showClear />
          </FormRow>
          <FormRow :label="t('kpi.status')" :errors="errors?.status">
            <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full" />
          </FormRow>
        </div>
        <div class="mt-4">
          <FormRow :label="t('kpi.description_label')" :errors="errors?.description">
            <Textarea v-model="form.description" rows="2" :placeholder="t('kpi.description_placeholder')" class="w-full" />
          </FormRow>
        </div>
      </div>

      <!-- Indicators Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-list text-blue-500"></i>
            {{ t('kpi.indicators') }}
            <span class="text-xs font-normal text-gray-500">({{ indicators.length }})</span>
          </h3>
          <div class="flex items-center gap-2">
            <span class="text-xs" :class="totalWeight === 100 ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'">
              {{ t('kpi.total_weight') }}: {{ totalWeight.toFixed(2) }}%
            </span>
            <Button :label="t('kpi.add_indicator')" icon="pi pi-plus" size="small" outlined @click="addIndicator" />
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="indicators.length === 0" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-lg">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('kpi.no_indicators') }}</p>
          <p class="text-xs mt-1">{{ t('kpi.no_indicators_hint') }}</p>
          <Button :label="t('kpi.add_first_indicator')" icon="pi pi-plus" size="small" class="mt-3" @click="addIndicator" />
        </div>

        <!-- Indicators Table -->
        <DataTable v-else :value="indicators" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
          <Column :header="t('kpi.perspective')" style="width:150px">
            <template #body="{data, index}">
              <Select v-model="data.perspective_id" :options="perspectiveOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="w-full !text-xs" size="small" />
            </template>
          </Column>
          <Column :header="t('kpi.indicator_title')" style="min-width:200px">
            <template #body="{data}">
              <InputText v-model="data.title" :placeholder="t('kpi.indicator_title_placeholder')" class="w-full !text-xs" size="small" />
            </template>
          </Column>
          <Column :header="t('kpi.weight')" style="width:90px">
            <template #body="{data}">
              <InputNumber v-model="data.weight" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" suffix="%" class="w-full !text-xs" size="small" />
            </template>
          </Column>
          <Column :header="t('kpi.formula')" style="width:130px">
            <template #body="{data}">
              <Select v-model="data.formula_type" :options="formulaOptions" optionLabel="label" optionValue="value" class="w-full !text-xs" size="small" />
            </template>
          </Column>
          <Column style="width:50px">
            <template #body="{data, index}">
              <Button icon="pi pi-trash" size="small" text severity="danger" @click="removeIndicator(index)" />
            </template>
          </Column>
        </DataTable>
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
  status: 'DRAFT'
})

const indicators = ref([])

const organizationOptions = ref([])
const periodOptions = ref([])
const perspectiveOptions = ref([])
const formulaOptions = ref([])

const statusOptions = [
  { label: 'Draft', value: 'DRAFT' },
  { label: 'Published', value: 'PUBLISHED' },
  { label: 'Archived', value: 'ARCHIVED' }
]

const totalWeight = computed(() => {
  return indicators.value.reduce((sum, ind) => sum + (ind.weight || 0), 0)
})

function addIndicator() {
  indicators.value.push({
    id: null,
    perspective_id: perspectiveOptions.value[0]?.value || null,
    title: '',
    description: '',
    weight: 0,
    formula_type: 'HIGHER_BETTER',
    target_type: 'NUMBER',
    indicator_type: 'MAXIMIZATION',
    is_required: true,
    sort_order: indicators.value.length
  })
}

function removeIndicator(index) {
  indicators.value.splice(index, 1)
}

function goBack() {
  router.push('/performance/kpi/templates')
}

async function loadReferenceData() {
  try {
    const [orgRes, periodRes, perspRes, formulaRes] = await Promise.all([
      api.get('/api/v1/tenant/organizations', { params: { per_page: 200, active_only: true } }),
      api.get('/api/v1/tenant/performance/periods', { params: { per_page: 50 } }),
      api.get('/api/v1/tenant/performance/kpi/perspectives', { params: { per_page: 20 } }),
      api.get('/api/v1/tenant/performance/indicator-formulas', { params: { per_page: 20 } })
    ])

    organizationOptions.value = (orgRes.data?.data || []).map(o => ({
      label: o.nomenclature || o.name || o.code,
      value: o.id
    }))

    periodOptions.value = (periodRes.data?.data || []).map(p => ({
      label: `${p.period_code} (${p.year})`,
      value: p.id
    }))

    perspectiveOptions.value = (perspRes.data?.data || []).map(p => ({
      label: p.name,
      value: p.id
    }))

    formulaOptions.value = (formulaRes.data?.data || []).map(f => ({
      label: f.name,
      value: f.formula_type || f.code
    }))

    // Add default formula options if empty
    if (formulaOptions.value.length === 0) {
      formulaOptions.value = [
        { label: 'Higher Better', value: 'HIGHER_BETTER' },
        { label: 'Lower Better', value: 'LOWER_BETTER' },
        { label: 'Range', value: 'RANGE' },
        { label: 'Manual', value: 'MANUAL' }
      ]
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 4000 })
  }
}

async function loadTemplate() {
  if (!isEditing.value) return

  try {
    const res = await api.get(`/api/v1/tenant/performance/kpi/templates/${templateId.value}`)
    const data = res.data?.data || res.data

    form.value = {
      name: data.name || '',
      description: data.description || '',
      organization_id: data.organization_id || null,
      period_id: data.period_id || null,
      status: data.status || 'DRAFT'
    }

    // Load indicators for this template
    const indRes = await api.get('/api/v1/tenant/performance/kpi/indicators', {
      params: { template_id: templateId.value, per_page: 100 }
    })
    indicators.value = (indRes.data?.data || []).map(ind => ({
      id: ind.id,
      perspective_id: ind.perspective_id,
      title: ind.title || '',
      description: ind.description || '',
      weight: ind.weight || 0,
      formula_type: ind.formula_type || 'HIGHER_BETTER',
      target_type: ind.target_type || 'NUMBER',
      indicator_type: ind.indicator_type || 'MAXIMIZATION',
      is_required: ind.is_required ?? true,
      sort_order: ind.sort_order || 0
    }))
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  }
}

async function handleSave() {
  errors.value = {}

  // Validation
  if (!form.value.name?.trim()) {
    errors.value = { name: [t('form.required')] }
    return
  }

  if (totalWeight.value !== 100 && indicators.value.length > 0) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('kpi.weight_must_100'), life: 4000 })
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
      await api.put(`/api/v1/tenant/performance/kpi/templates/${templateId.value}`, payload)
    } else {
      const res = await api.post('/api/v1/tenant/performance/kpi/templates', payload)
      savedTemplateId = res.data?.data?.id || res.data?.id
    }

    // Save indicators
    for (let i = 0; i < indicators.value.length; i++) {
      const ind = indicators.value[i]
      const indPayload = {
        performance_template_id: savedTemplateId,
        perspective_id: ind.perspective_id,
        title: ind.title,
        description: ind.description || null,
        weight: ind.weight,
        formula_type: ind.formula_type,
        target_type: ind.target_type || 'NUMBER',
        indicator_type: ind.indicator_type || 'MAXIMIZATION',
        is_required: ind.is_required ?? true,
        sort_order: i,
        minimum_score: 0,
        maximum_score: 100
      }

      if (ind.id) {
        await api.put(`/api/v1/tenant/performance/kpi/indicators/${ind.id}`, indPayload)
      } else {
        await api.post('/api/v1/tenant/performance/kpi/indicators', indPayload)
      }
    }

    toast.add({
      severity: 'success',
      summary: t('message.success'),
      detail: isEditing.value ? t('kpi.template_updated') : t('kpi.template_created'),
      life: 3000
    })

    router.push('/performance/kpi/templates')
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
