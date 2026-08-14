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
                <p v-if="item.education_major_name?.length" class="text-sm text-gray-700 dark:text-gray-300 mt-0.5">
                  <span class="text-xs font-medium text-gray-400 uppercase tracking-wide mr-1">{{ t('requisitions.education_major') }}:</span>{{ item.education_major_name.join(', ') }}
                </p>
                <p v-if="item.job_family_name?.length" class="text-sm text-gray-700 dark:text-gray-300 mt-0.5">
                  <span class="text-xs font-medium text-gray-400 uppercase tracking-wide mr-1">{{ t('requisitions.job_family') }}:</span>{{ item.job_family_name.join(', ') }}
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
          <div class="flex items-center justify-between mb-3 flex-wrap gap-2">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('requisitions.competencies_tab') }}</h4>
            <Button :label="t('requisitions.sync_competencies')" icon="pi pi-sync" size="small" severity="secondary" outlined :loading="syncing" :disabled="!requisition?.organization_id" @click="syncFromJobManagement()" />
          </div>

          <!-- Section: Technical & Managerial (pola Job Management) -->
          <div v-for="sec in competencySections" :key="sec.type" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 mb-3 last:mb-0">
            <div class="mb-2">
              <h5 class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">{{ t(sec.titleKey) }}</h5>
            </div>

            <!-- Own items (override requisition) -->
            <template v-if="sec.items.length">
              <div class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
                <div v-for="item in sec.items" :key="item.id" class="px-3 py-2.5">
                  <p class="text-sm text-gray-800 dark:text-gray-100">{{ competencyName(item) }}</p>
                  <p v-if="item.required_level || item.weight" class="text-xs text-gray-400">
                    <span v-if="item.required_level">{{ t('requisitions.required_level') }} {{ item.required_level }}</span>
                    <span v-if="item.weight"> · {{ t('requisitions.weight') }} {{ item.weight }}%</span>
                  </p>
                </div>
              </div>
            </template>

            <!-- Fallback Job Management (read-only, sumber match score G-9) -->
            <template v-else-if="sec.fallback.length">
              <p class="text-[11px] text-amber-600 dark:text-amber-400 mb-1.5"><i class="pi pi-info-circle mr-1"></i>{{ t('requisitions.from_job_management') }}</p>
              <div class="divide-y divide-gray-100 dark:divide-gray-800 border border-dashed border-amber-200 dark:border-amber-700/40 rounded-lg">
                <div v-for="item in sec.fallback" :key="item.id" class="flex items-center gap-2 px-3 py-2.5">
                  <div class="min-w-0 flex-1">
                    <p class="text-sm text-gray-700 dark:text-gray-300">{{ competencyName(item) }}</p>
                    <p v-if="item.level || item.weight" class="text-xs text-gray-400">
                      <span v-if="item.level">{{ t('requisitions.level') }} Lv.{{ item.level }}<span v-if="item.level_description"> — {{ item.level_description }}</span></span>
                      <span v-if="item.weight"> · {{ t('requisitions.weight') }} {{ item.weight }}%</span>
                    </p>
                  </div>
                </div>
              </div>
            </template>

            <div v-else class="flex flex-col items-center justify-center py-6 text-gray-400 dark:text-gray-500">
              <p class="text-sm">{{ t('requisitions.competencies_empty') }}</p>
            </div>
          </div>

          <!-- Sisa kompetensi yang tidak masuk Teknis/Manajerial (cluster lain) -->
          <div v-if="otherItems.length" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <div class="flex items-center justify-between mb-2">
              <h5 class="text-xs font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400">{{ t('requisitions.other_competencies') }}</h5>
            </div>
            <div class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
              <div v-for="item in otherItems" :key="item.id" class="px-3 py-2.5">
                <p class="text-sm text-gray-800 dark:text-gray-100">{{ competencyName(item) }}</p>
                <p v-if="item.required_level || item.weight" class="text-xs text-gray-400">
                  <span v-if="item.required_level">{{ t('requisitions.required_level') }} {{ item.required_level }}</span>
                  <span v-if="item.weight"> · {{ t('requisitions.weight') }} {{ item.weight }}%</span>
                </p>
              </div>
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
const technicalClusters = ref([])
const managerialClusters = ref([])
const syncing = ref(false)

// item.competency_name (backend-resolved, requisition's own competencies)
// diprioritaskan; fallback ke lookup client-side dari competencyMaster
// (dipakai oleh baris Job Management yang tidak punya competency_name).
function competencyName(item) {
  if (item.competency_name) return item.competency_name
  const c = competencyMeta(item)
  return c ? c.name : (item.competency_id || '')
}

// Meta kompetensi dari master (id → { name, field, cluster })
function competencyMeta(item) {
  return competencyMaster.value.find(x => x.id === item.competency_id) || null
}

// Split kompetensi milik requisition sendiri per section (cluster dari mapping Job Management)
function splitOwnBySection() {
  const tech = new Set(technicalClusters.value)
  const man = new Set(managerialClusters.value)
  const technical = []
  const managerial = []
  const other = []
  for (const item of competencyItems.value) {
    const c = competencyMeta(item)
    const cluster = c?.cluster || ''
    if (cluster && tech.has(cluster)) technical.push(item)
    else if (cluster && man.has(cluster)) managerial.push(item)
    else other.push(item)
  }
  return { technical, managerial, other }
}

// Klasifikasi baris fallback Job Management: type dari JobValue (backend baru);
// fallback ke field kompetensi (backend lama / JobValue tanpa type).
function fallbackType(item) {
  if (item.type) return item.type
  const c = competencyMeta(item)
  if (c?.field === 'Technical Competency') return 'technical'
  if (c?.field === 'Manajerial') return 'managerial'
  return ''
}

const competencySections = computed(() => {
  const own = splitOwnBySection()
  // Fallback Job Management hanya tampil saat requisition belum punya override sendiri
  const hasOwn = competencyItems.value.length > 0
  return [
    {
      type: 'technical',
      titleKey: 'job_management.potency_technical_title',
      items: own.technical,
      fallback: hasOwn ? [] : jobManagementCompetencies.value.filter(it => fallbackType(it) === 'technical' && it.competency_id)
    },
    {
      type: 'managerial',
      titleKey: 'job_management.potency_managerial_title',
      items: own.managerial,
      fallback: hasOwn ? [] : jobManagementCompetencies.value.filter(it => fallbackType(it) === 'managerial' && it.competency_id)
    }
  ]
})

const otherItems = computed(() => (competencyItems.value.length ? splitOwnBySection().other : []))

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
    // Pakai settings/competencies (max per_page 500, tidak di-reset ke 20 seperti
    // competency/competencies) agar master kompetensi lengkap — dipakai lookup
    // nama/cluster/field dan opsi dropdown tambah.
    const res = await api.get('/api/v1/tenant/settings/competencies', { params: { per_page: 500 } })
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

// Sinkronisasi ulang: samakan kompetensi requisition dengan kompetensi potensi
// Job Management (org requisition). Kompetensi baru dibuat (POST), kompetensi
// yang sudah ada diperbarui (PUT — level & bobot ikut ter-backfill), sehingga
// level tidak hilang saat re-sync. Dedupe per competency_id.
async function syncFromJobManagement() {
  const orgId = requisition.value?.organization_id
  if (!orgId) return
  syncing.value = true
  try {
    const compRes = await api.get('/api/v1/tenant/job-management/potency-competencies', { params: { organization_id: orgId, per_page: 100 } })
    const rows = compRes.data?.data || []
    const existingByComp = new Map(competencyItems.value.filter(c => c.competency_id).map(c => [c.competency_id, c]))
    let added = 0
    let updated = 0
    for (const row of rows) {
      if (!row.competency_id) continue
      const existing = existingByComp.get(row.competency_id)
      if (existing) {
        await api.put(`/api/v1/tenant/recruitment/requisition-competencies/${existing.id}`, cleanItemPayload({ weight: row.weight, required_level: row.level }))
        updated++
      } else {
        await api.post(`/api/v1/tenant/recruitment/requisitions/${requisitionId}/competencies`, cleanItemPayload({ competency_id: row.competency_id, weight: row.weight, required_level: row.level }))
        added++
      }
    }
    let detail
    if (added > 0 && updated > 0) {
      detail = `${t('requisitions.synced_count', { count: added })} · ${t('requisitions.synced_updated', { count: updated })}`
    } else if (added > 0) {
      detail = t('requisitions.synced_count', { count: added })
    } else if (updated > 0) {
      detail = t('requisitions.synced_updated', { count: updated })
    } else {
      detail = t('requisitions.synced_none')
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail, life: 3000 })
    loadCompetencyItems()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 5000 })
  } finally {
    syncing.value = false
  }
}

// Mapping cluster technical/managerial (Job Management) — filter dropdown tambah
async function loadClusterMappings() {
  try {
    const [techRes, manRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/job-management/values/clusters/technical'),
      api.get('/api/v1/tenant/job-management/values/clusters/managerial')
    ])
    technicalClusters.value = techRes.status === 'fulfilled' ? (techRes.value.data?.data?.clusters || []) : []
    managerialClusters.value = manRes.status === 'fulfilled' ? (manRes.value.data?.data?.clusters || []) : []
  } catch {
    technicalClusters.value = []
    managerialClusters.value = []
  }
}

onMounted(async () => {
  loading.value = true
  await loadRequisition()
  loading.value = false
  loadCompetencyMaster()
  loadClusterMappings()
  loadRequirements()
  loadCompetencyItems()
  loadJobManagementFallback()
})
</script>
