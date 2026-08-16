<template>
  <div class="space-y-4 max-w-4xl">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies/templates')" />
        <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ editing ? t('competency_360.edit_template') : t('competency_360.new_template') }}</h2>
      </div>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="space-y-3">
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="255" :placeholder="t('common.name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.code_auto_hint') }}</p>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('competency_360.scale')">
          <Select v-model="form.scale_id" :options="scaleOptions" optionLabel="name" optionValue="id" showClear filter class="w-full" :placeholder="t('common.select')" />
        </FormRow>
      </div>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center justify-between mb-2">
        <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.competencies') }}</p>
        <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('competency_360.add_competency')" @click="addCompetency" />
      </div>
      <div class="flex items-center gap-2 mb-1 px-1">
        <div class="flex-1 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.competency') }}</div>
        <div class="w-24 shrink-0 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.req_level') }}</div>
        <div class="w-20 shrink-0 text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('competency_360.weight') }}</div>
        <div class="w-8 shrink-0"></div>
      </div>
      <div v-for="(comp, idx) in form.competencies" :key="idx" class="flex items-center gap-2 mb-2">
        <div class="flex-1">
          <Select v-model="comp.competency_id" :options="competencyOptions" optionLabel="name" optionValue="id" filter class="w-full !text-xs" :placeholder="t('competency_360.select_competency')" />
        </div>
        <div class="w-24 shrink-0">
          <TextInput v-model="comp.required_level" type="number" :placeholder="t('competency_360.req_level')" class="!text-xs" />
        </div>
        <div class="w-20 shrink-0">
          <TextInput v-model="comp.weight" type="number" :placeholder="t('competency_360.weight')" class="!text-xs" />
        </div>
        <Button icon="pi pi-trash" size="small" text severity="danger" @click="form.competencies.splice(idx, 1)" />
      </div>
      <p v-if="form.competencies.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_competencies_hint') }}</p>
    </div>

    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center justify-between mb-2">
        <p class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.rater_types') }}</p>
        <Button icon="pi pi-plus" size="small" text severity="secondary" :label="t('competency_360.add_rater_type')" @click="addRaterType" />
      </div>
      <div v-for="(rt, idx) in form.rater_types" :key="idx" class="border border-gray-200 dark:border-gray-700 rounded-lg p-2 mb-2">
        <div class="grid grid-cols-2 gap-2">
          <div>
            <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.rater_type') }}</label>
            <Select v-model="rt.rater_type" :options="raterTypeOptions" optionLabel="label" optionValue="value" class="w-full !text-xs" />
          </div>
          <div>
            <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.weight') }}</label>
            <TextInput v-model="rt.weight" type="number" class="!text-xs" />
          </div>
          <div>
            <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.min_rater') }}</label>
            <TextInput v-model="rt.min_rater" type="number" class="!text-xs" />
          </div>
          <div>
            <label class="text-[11px] text-gray-500 dark:text-gray-400 block mb-1">{{ t('competency_360.max_rater') }}</label>
            <TextInput v-model="rt.max_rater" type="number" class="!text-xs" />
          </div>
        </div>
        <div class="flex items-center gap-4 mt-2">
          <label class="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300">
            <ToggleSwitch v-model="rt.required" /> {{ t('competency_360.required') }}
          </label>
          <label class="flex items-center gap-1.5 text-xs text-gray-600 dark:text-gray-300">
            <ToggleSwitch v-model="rt.anonymous" /> {{ t('competency_360.anonymous') }}
          </label>
          <Button icon="pi pi-trash" size="small" text severity="danger" class="ml-auto" @click="form.rater_types.splice(idx, 1)" />
        </div>
      </div>
      <p v-if="form.rater_types.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_rater_types_hint') }}</p>
    </div>

    <div class="flex items-center justify-end gap-2">
      <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="router.push('/competencies/templates')" />
      <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import Button from 'primevue/button'
import Select from 'primevue/select'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const editingId = route.params.id || null
const editing = computed(() => !!editingId)

const scales = ref([])
const competencies = ref([])
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const raterTypeOptions = [
  { label: t('competency_360.rater_type_self'), value: 'self' },
  { label: t('competency_360.rater_type_superior'), value: 'superior' },
  { label: t('competency_360.rater_type_peer'), value: 'peer' },
  { label: t('competency_360.rater_type_subordinate'), value: 'subordinate' },
  { label: t('competency_360.rater_type_other'), value: 'other' }
]

const scaleOptions = computed(() => scales.value.filter(s => s.status !== 'inactive' || s.id === form.value.scale_id))
const competencyOptions = computed(() => competencies.value)

function defaultForm() {
  return { name: '', description: '', status: 'active', scale_id: null, competencies: [], rater_types: [] }
}

function newCompetency() {
  return { competency_id: null, required_level: null, weight: 1, sort_order: form.value.competencies.length }
}

function newRaterType() {
  return { rater_type: 'peer', weight: 1, min_rater: 1, max_rater: null, required: false, anonymous: false }
}

function addCompetency() {
  form.value.competencies.push(newCompetency())
}

function addRaterType() {
  form.value.rater_types.push(newRaterType())
}

async function loadReferences() {
  try {
    const [scaleRes, compRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/rating-scales', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/competency/competencies', { params: { per_page: 500 } })
    ])
    scales.value = scaleRes.status === 'fulfilled' ? (scaleRes.value.data?.data || []) : []
    competencies.value = compRes.status === 'fulfilled' ? (compRes.value.data?.data || []) : []
  } catch {
    // fail-silent
  }
}

async function loadTemplate() {
  try {
    const res = await api.get(`/api/v1/tenant/competency/templates/${editingId}`)
    const item = res.data?.data
    if (!item) return
    form.value = {
      name: item.name || '',
      description: item.description || '',
      status: item.status || 'active',
      scale_id: item.scale_id || null,
      competencies: (item.competencies || []).map(c => ({ competency_id: c.competency_id, required_level: c.required_level ?? null, weight: c.weight ?? 1, sort_order: c.sort_order || 0 })),
      rater_types: (item.rater_types || []).map(rt => ({ rater_type: rt.rater_type, weight: rt.weight ?? 1, min_rater: rt.min_rater ?? 1, max_rater: rt.max_rater ?? null, required: !!rt.required, anonymous: !!rt.anonymous }))
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

async function handleSave() {
  errors.value = {}
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      name: form.value.name.trim(),
      description: form.value.description?.trim() || '',
      status: form.value.status || 'active',
      scale_id: form.value.scale_id || undefined,
      competencies: form.value.competencies
        .filter(c => c.competency_id)
        .map((c, idx) => ({
          competency_id: c.competency_id,
          required_level: c.required_level !== null && c.required_level !== '' ? Number(c.required_level) : undefined,
          weight: Number(c.weight) || 0,
          sort_order: Number(c.sort_order) || idx
        })),
      rater_types: form.value.rater_types
        .map((rt, idx) => ({
          rater_type: rt.rater_type,
          weight: Number(rt.weight) || 0,
          min_rater: Number(rt.min_rater) || 0,
          max_rater: rt.max_rater !== null && rt.max_rater !== '' ? Number(rt.max_rater) : undefined,
          required: !!rt.required,
          anonymous: !!rt.anonymous
        }))
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/competency/templates/${editingId}`, payload)
    } else {
      await api.post('/api/v1/tenant/competency/templates', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    router.push('/competencies/templates')
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadReferences()
  if (editing.value) {
    loadTemplate()
  } else {
    addCompetency()
    addRaterType()
  }
})
</script>
