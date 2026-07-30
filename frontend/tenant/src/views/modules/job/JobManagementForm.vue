<template>
  <div class="max-w-full mx-auto">
    <!-- Loading -->
    <div v-if="pageLoading" class="flex gap-6">
      <div class="w-56 space-y-2">
        <div v-for="n in 8" :key="n" class="h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
      <div class="flex-1 space-y-3">
        <div v-for="n in 6" :key="n" class="h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
    </div>

    <!-- Two-column layout -->
    <div v-else class="flex gap-6">
      <!-- Left: Navigation Sidebar -->
      <div class="w-56 shrink-0 space-y-1">
        <!-- Section Navigation -->
        <div
          v-for="(s, i) in sections" :key="i"
          role="button" :tabindex="0"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none"
          :class="getNavItemClass(i)"
          @click="selectSection(i)" @keydown.enter="selectSection(i)"
        >
          <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150" :class="getCircleClass(i)">
            <i v-if="sectionSaved[i]" class="pi pi-check text-xs"></i>
            <i v-else :class="s.icon"></i>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate" :class="getNavTextClass(i)">{{ t(s.labelKey) }}</div>
          </div>
          <i v-if="sectionSaved[i]" class="pi pi-check-circle text-emerald-400 text-xs shrink-0"></i>
        </div>
      </div>

      <!-- Right: Form Content -->
      <div class="flex-1 min-w-0">
        <component
          :is="currentSectionComp"
          :key="activeSection"
          :org-id="orgId"
          :org-name="orgName"
          :org-code="orgCode"
          :org-grading-id="orgGradingId"
          :org-job-family-id="orgJobFamilyId"
          :grading-options="gradingOptions"
          :job-family-options="jobFamilyOptions"
          :job-value-options="jobValueOptions"
          :competency-options="competencyOptions"
          :job-value-map="jobValueMap"
          
          @saved="onSectionSaved"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'

// Section components
import JobIdentificationSection from './sections/JobIdentificationSection.vue'
import JobObjectiveSection from './sections/JobObjectiveSection.vue'
import JobEduExpSection from './sections/JobEduExpSection.vue'
import JobResponsibilitySection from './sections/JobResponsibilitySection.vue'
import JobHRAuthoritySection from './sections/JobHRAuthoritySection.vue'
import JobOpAuthoritySection from './sections/JobOpAuthoritySection.vue'
import JobActivitySection from './sections/JobActivitySection.vue'
import JobRiskSection from './sections/JobRiskSection.vue'
import JobRelationshipSection from './sections/JobRelationshipSection.vue'
import JobSubordinateSection from './sections/JobSubordinateSection.vue'
import JobAssetSection from './sections/JobAssetSection.vue'
import JobFinancialSection from './sections/JobFinancialSection.vue'
import JobPotencySection from './sections/JobPotencySection.vue'
import JobCompetencyGroupSection from './sections/JobCompetencyGroupSection.vue'
import JobScoreSection from './sections/JobScoreSection.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const orgId = route.query.org_id || ''
const activeSection = ref(0)
const pageLoading = ref(true)
const sectionSaved = ref(Array(15).fill(false))

const orgName = ref('')
const orgCode = ref('')
const orgGradingId = ref('')
const orgJobFamilyId = ref('')
const gradingOptions = ref([])
const jobFamilyOptions = ref([])
const jobValueOptions = ref([])
const jobValueMap = ref({})
const competencyOptions = ref([])

const sections = [
  { labelKey: 'job_management.identifications', icon: 'pi pi-id-card', comp: JobIdentificationSection },
  { labelKey: 'job_management.objectives', icon: 'pi pi-bullseye', comp: JobObjectiveSection },
  { labelKey: 'job_management.education_experience', icon: 'pi pi-graduation-cap', comp: JobEduExpSection },
  { labelKey: 'job_management.responsibilities_title', icon: 'pi pi-list-check', comp: JobResponsibilitySection },
  { labelKey: 'job_management.hr_authorities', icon: 'pi pi-users', comp: JobHRAuthoritySection },
  { labelKey: 'job_management.op_authorities', icon: 'pi pi-cog', comp: JobOpAuthoritySection },
  { labelKey: 'job_management.activities', icon: 'pi pi-bolt', comp: JobActivitySection },
  { labelKey: 'job_management.risks', icon: 'pi pi-exclamation-triangle', comp: JobRiskSection },
  { labelKey: 'job_management.relationships', icon: 'pi pi-share-alt', comp: JobRelationshipSection },
  { labelKey: 'job_management.subordinate_controls', icon: 'pi pi-sitemap', comp: JobSubordinateSection },
  { labelKey: 'job_management.assets', icon: 'pi pi-box', comp: JobAssetSection },
  { labelKey: 'job_management.financials', icon: 'pi pi-money-bill', comp: JobFinancialSection },
  { labelKey: 'job_management.potency_competencies', icon: 'pi pi-star', comp: JobPotencySection },
  { labelKey: 'job_management.competency_groups', icon: 'pi pi-chart-pie', comp: JobCompetencyGroupSection },
  { labelKey: 'job_management.scores', icon: 'pi pi-calculator', comp: JobScoreSection },
]

const currentSectionComp = computed(() => sections[activeSection.value]?.comp || null)

function getNavItemClass(i) {
  const active = activeSection.value === i
  if (active) return 'bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700'
  if (sectionSaved.value[i]) return 'hover:bg-gray-50 dark:hover:bg-gray-800'
  return 'hover:bg-gray-50 dark:hover:bg-gray-800'
}

function getCircleClass(i) {
  const active = activeSection.value === i
  if (active) return 'bg-emerald-600 text-white'
  if (sectionSaved.value[i]) return 'bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300'
  return 'bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300'
}

function getNavTextClass(i) {
  const active = activeSection.value === i
  if (active) return 'text-emerald-700 dark:text-emerald-300'
  if (sectionSaved.value[i]) return 'text-emerald-600 dark:text-emerald-400'
  return 'text-gray-700 dark:text-gray-300'
}

function selectSection(i) {
  activeSection.value = i
  router.replace({ query: { ...route.query, section: String(i) } })
}

function onSectionSaved(i) {
  if (typeof i === 'number') sectionSaved.value[i] = true
}

function goBack() {
  router.push('/job-management')
}

async function loadOrgInfo() {
  if (!orgId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${orgId}`)
    const data = res.data?.data
    if (data) {
      orgName.value = data.nomenclature || ''
      orgCode.value = data.full_code || data.code || ''
      orgGradingId.value = data.grading_id || ''
      orgJobFamilyId.value = data.job_family_id || ''
    }
  } catch { /* ignore */ }
}

async function loadRefData() {
  try {
    const [gradingRes, valuesRes, compRes, jfRes] = await Promise.all([
      api.get('/api/v1/tenant/settings/gradings?per_page=100'),
      api.get('/api/v1/tenant/job-management/values?per_page=200'),
      api.get('/api/v1/tenant/competencies?per_page=200').catch(() => ({ data: { data: [] } })),
      api.get('/api/v1/tenant/settings/job-families?per_page=100')
    ])
    gradingOptions.value = (gradingRes.data?.data || []).map(g => ({ label: `${g.code} - ${g.name}`, value: g.id }))
    jobFamilyOptions.value = (jfRes.data?.data || []).map(j => ({ label: `${j.code} - ${j.name}`, value: j.id }))
    const vals = valuesRes.data?.data || []
    jobValueOptions.value = vals.map(v => ({
      label: `${v.type}${v.level ? ' Lv.'+v.level : ''}${v.descriptions ? ' — '+v.descriptions : ''}`,
      value: v.id, type: v.type, level: v.level, descriptions: v.descriptions
    }))
    // Build map: type -> [options]
    const map = {}
    vals.forEach(v => {
      if (!map[v.type]) map[v.type] = []
      map[v.type].push({ label: `Lv.${v.level} — ${v.descriptions || ''}`, value: v.id, level: v.level })
    })
    jobValueMap.value = map
    competencyOptions.value = (compRes.data?.data || []).map(c => ({ label: c.name || c.code, value: c.id }))
  } catch { /* ignore */ }
}

onMounted(async () => {
  try {
    await Promise.all([loadOrgInfo(), loadRefData()])
    const sectionParam = parseInt(route.query.section)
    if (!isNaN(sectionParam) && sectionParam >= 0 && sectionParam < sections.length) {
      activeSection.value = sectionParam
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally { pageLoading.value = false }
})
</script>
