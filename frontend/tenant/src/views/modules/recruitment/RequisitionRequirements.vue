<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-16 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="requisition">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 flex items-center justify-between gap-3 flex-wrap">
        <div class="min-w-0">
          <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100 truncate">{{ requisition.title }}</h2>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ t('requisitions.requirements_competencies') }}</p>
        </div>
        <Button :label="t('common.back')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/recruitment/requisitions')" />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <!-- Requirements -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('requisitions.requirements_tab') }}</h4>
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" @click="openAddRequirement()" />
          </div>
          <div v-if="requirements.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in requirements" :key="item.id" class="flex items-center gap-2 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm text-gray-800 dark:text-gray-100">{{ item.name }}</p>
                <p class="text-xs text-gray-400">{{ item.requirement_type }}<span v-if="item.minimum_value || item.maximum_value"> · {{ item.minimum_value ?? '?' }}–{{ item.maximum_value ?? '?' }}</span></p>
              </div>
              <Tag v-if="item.is_required" :value="t('requisitions.required')" severity="danger" class="!text-[10px] !px-1.5 !py-0" />
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="removeRequirement(item)" />
            </div>
          </div>
          <!-- Job Management default (fallback read-only, hanya tampil kalau requisition belum punya override sendiri) -->
          <template v-else-if="jobManagementEducationExperiences.length">
            <p class="text-[11px] text-amber-600 dark:text-amber-400 mb-1.5"><i class="pi pi-info-circle mr-1"></i>{{ t('requisitions.from_job_management') }}</p>
            <div class="divide-y divide-gray-100 dark:divide-gray-800 border border-dashed border-amber-200 dark:border-amber-700/40 rounded-lg">
              <div v-for="item in jobManagementEducationExperiences" :key="item.id" class="px-3 py-2.5">
                <p v-if="!item.education_name && !item.experience_name" class="text-sm text-gray-700 dark:text-gray-300">{{ item.nomenclature }}</p>
                <p v-if="item.education_name" class="text-sm text-gray-700 dark:text-gray-300">
                  <span class="text-xs font-medium text-gray-400 uppercase tracking-wide mr-1">{{ t('requisitions.education') }}:</span>{{ item.education_name }}
                </p>
                <p v-if="item.experience_name" class="text-sm text-gray-700 dark:text-gray-300 mt-0.5">
                  <span class="text-xs font-medium text-gray-400 uppercase tracking-wide mr-1">{{ t('requisitions.experience') }}:</span>{{ item.experience_name }}
                </p>
              </div>
            </div>
          </template>
          <div v-else class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-list-check text-2xl mb-2 opacity-50"></i>
            <p class="text-sm">{{ t('requisitions.requirements_empty') }}</p>
          </div>

          <div v-if="addRequirementVisible" class="mt-3 space-y-2 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
            <TextInput v-model="newRequirement.requirement_type" :placeholder="t('requisitions.requirement_type')" class="!w-full" />
            <TextInput v-model="newRequirement.name" :placeholder="t('requisitions.requirement_name')" class="!w-full" />
            <div class="grid grid-cols-2 gap-2">
              <InputNumber v-model="newRequirement.minimum_value" :placeholder="t('requisitions.minimum_value')" class="!w-full" />
              <InputNumber v-model="newRequirement.maximum_value" :placeholder="t('requisitions.maximum_value')" class="!w-full" />
            </div>
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="addRequirementVisible = false" />
              <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveRequirement()" />
            </div>
          </div>
        </div>

        <!-- Competencies -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <div class="flex items-center justify-between mb-3">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('requisitions.competencies_tab') }}</h4>
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" @click="openAddCompetency()" />
          </div>
          <div v-if="competencyItems.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in competencyItems" :key="item.id" class="flex items-center gap-2 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm text-gray-800 dark:text-gray-100">{{ competencyName(item.competency_id) }}</p>
                <p v-if="item.required_level || item.weight" class="text-xs text-gray-400">
                  <span v-if="item.required_level">{{ t('requisitions.required_level') }} {{ item.required_level }}</span>
                  <span v-if="item.weight"> · {{ t('requisitions.weight') }} {{ item.weight }}%</span>
                </p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="removeCompetency(item)" />
            </div>
          </div>
          <!-- Job Management default (fallback read-only — sumber yang dipakai match score G-9 bila requisition ini belum punya override) -->
          <template v-else-if="jobManagementCompetencies.length">
            <p class="text-[11px] text-amber-600 dark:text-amber-400 mb-1.5"><i class="pi pi-info-circle mr-1"></i>{{ t('requisitions.from_job_management') }}</p>
            <div class="divide-y divide-gray-100 dark:divide-gray-800 border border-dashed border-amber-200 dark:border-amber-700/40 rounded-lg">
              <div v-for="item in jobManagementCompetencies" :key="item.id" class="flex items-center gap-2 px-3 py-2.5">
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-gray-700 dark:text-gray-300">{{ competencyName(item.competency_id) }}</p>
                  <p v-if="item.weight" class="text-xs text-gray-400">{{ t('requisitions.weight') }} {{ item.weight }}%</p>
                </div>
              </div>
            </div>
          </template>
          <div v-else class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-sparkles text-2xl mb-2 opacity-50"></i>
            <p class="text-sm">{{ t('requisitions.competencies_empty') }}</p>
          </div>

          <div v-if="addCompetencyVisible" class="mt-3 space-y-2 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
            <SelectLabel v-model="newCompetency.competency_id" :options="competencyOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
            <div class="grid grid-cols-2 gap-2">
              <InputNumber v-model="newCompetency.required_level" :placeholder="t('requisitions.required_level')" :min="1" :max="5" class="!w-full" />
              <InputNumber v-model="newCompetency.weight" :placeholder="t('requisitions.weight')" :min="0" :max="100" class="!w-full" />
            </div>
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="addCompetencyVisible = false" />
              <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveCompetency()" />
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import InputNumber from 'primevue/inputnumber'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const requisitionId = route.params.id

const loading = ref(true)
const requisition = ref(null)
const requirements = ref([])
const competencyItems = ref([])
const competencyMaster = ref([])
const jobManagementEducationExperiences = ref([])
const jobManagementCompetencies = ref([])
const itemSaving = ref(false)
const addRequirementVisible = ref(false)
const newRequirement = ref({})
const addCompetencyVisible = ref(false)
const newCompetency = ref({})

const competencyOptions = computed(() => competencyMaster.value.map(c => ({ label: c.name, value: c.id })))

function competencyName(id) {
  const c = competencyMaster.value.find(x => x.id === id)
  return c ? c.name : id
}

function cleanItemPayload(payload) {
  const out = { ...payload }
  Object.keys(out).forEach(k => {
    if (out[k] === '' || out[k] === null || out[k] === undefined) delete out[k]
  })
  return out
}

async function loadRequisition() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/requisitions/${requisitionId}`)
    requisition.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

async function loadCompetencyMaster() {
  try {
    const res = await api.get('/api/v1/tenant/competency/competencies', { params: { per_page: 500 } })
    competencyMaster.value = res.data?.data || []
  } catch {
    competencyMaster.value = []
  }
}

// Fallback read-only (keputusan user: "Job Management jadi default, override
// tetap di Recruitment") — dipanggil hanya untuk ditampilkan saat requisition
// belum punya requirement/competency sendiri; sumber sesungguhnya untuk match
// score dihitung backend (GetCandidateMatchScore, G-9).
async function loadJobManagementFallback() {
  const orgId = requisition.value?.organization_id
  if (!orgId) {
    jobManagementEducationExperiences.value = []
    jobManagementCompetencies.value = []
    return
  }
  try {
    const [eduRes, compRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/job-management/education-experiences', { params: { organization_id: orgId, per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/potency-competencies', { params: { organization_id: orgId, per_page: 100 } })
    ])
    jobManagementEducationExperiences.value = eduRes.status === 'fulfilled' ? (eduRes.value.data?.data || []) : []
    jobManagementCompetencies.value = compRes.status === 'fulfilled' ? (compRes.value.data?.data || []) : []
  } catch {
    jobManagementEducationExperiences.value = []
    jobManagementCompetencies.value = []
  }
}

async function loadRequirements() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/requisitions/${requisitionId}/requirements`)
    requirements.value = res.data?.data || []
  } catch {
    requirements.value = []
  }
}
async function loadCompetencyItems() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/requisitions/${requisitionId}/competencies`)
    competencyItems.value = res.data?.data || []
  } catch {
    competencyItems.value = []
  }
}

function openAddRequirement() {
  newRequirement.value = { requirement_type: '', name: '', minimum_value: null, maximum_value: null }
  addRequirementVisible.value = true
}
async function saveRequirement() {
  if (!newRequirement.value.requirement_type?.trim() || !newRequirement.value.name?.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/requisitions/${requisitionId}/requirements`, cleanItemPayload(newRequirement.value))
    addRequirementVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadRequirements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
async function removeRequirement(item) {
  try {
    await api.delete(`/api/v1/tenant/recruitment/requirements/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_deleted'), life: 3000 })
    loadRequirements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  }
}

function openAddCompetency() {
  newCompetency.value = { competency_id: null, required_level: null, weight: null }
  addCompetencyVisible.value = true
}
async function saveCompetency() {
  if (!newCompetency.value.competency_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/requisitions/${requisitionId}/competencies`, cleanItemPayload(newCompetency.value))
    addCompetencyVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadCompetencyItems()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
async function removeCompetency(item) {
  try {
    await api.delete(`/api/v1/tenant/recruitment/requisition-competencies/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_deleted'), life: 3000 })
    loadCompetencyItems()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  }
}

onMounted(async () => {
  loading.value = true
  await loadRequisition()
  loading.value = false
  loadCompetencyMaster()
  loadRequirements()
  loadCompetencyItems()
  loadJobManagementFallback()
})
</script>
