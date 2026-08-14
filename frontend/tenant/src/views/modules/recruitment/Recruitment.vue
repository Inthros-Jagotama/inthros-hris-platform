<template>
  <div class="space-y-4">
    <!-- G-11/G-12 sub-2: summary cards — fail-silent kalau endpoint analytics error -->
    <div v-if="statsLoading" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <div v-for="i in 4" :key="i" class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
    </div>
    <div v-else-if="stats" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <div v-for="card in statCards" :key="card.labelKey" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-1.5">
          <span class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t(card.labelKey) }}</span>
          <i :class="card.icon" class="text-sm" :style="{ color: card.color }"></i>
        </div>
        <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ card.value }}</p>
      </div>
    </div>

    <!-- Menu cards (pola sama dengan halaman Workforce Intelligence / Career Intelligence) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      <button
        v-for="menu in menuCards"
        :key="menu.route"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-sky-300 dark:hover:border-sky-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-500/50"
        @click="router.push(menu.route)"
      >
        <div
          class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors"
          :class="menu.tint"
        >
          <i :class="menu.icon" class="text-base"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ t(menu.labelKey) }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(menu.descKey) }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-sky-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <!-- Fitur strategis berikutnya (belum ada halaman) — card "Coming soon" -->
      <div
        v-for="soon in comingSoonCards"
        :key="soon.labelKey"
        class="flex items-center gap-3 p-3.5 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800/50 text-left opacity-70 cursor-not-allowed"
      >
        <div
          class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center"
          :class="soon.tint"
        >
          <i :class="soon.icon" class="text-base"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-700 dark:text-gray-300 truncate">{{ t(soon.labelKey) }}</p>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
            <Tag :value="t('common.coming_soon')" severity="secondary" class="!text-[10px] !px-1.5 !py-0" />
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'

import Tag from 'primevue/tag'

const router = useRouter()
const { t } = useI18n()

// Hub Recruitment (operasional). Halaman strategis di bawah (S-1/S-4/S-5)
// dikelola di module-recruitment-strategic-layer-plan.md — Recruitment hanya
// menyediakan data operasional; logika strategis tetap di WI/CI.
const menuCards = computed(() => [
  // S-1/S-5: Job Requisitions — reason_type WORKFORCE_GAP / SUCCESSION_GAP
  { labelKey: 'recruitment.requisitions', descKey: 'requisitions.description', icon: 'pi pi-briefcase', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', route: '/recruitment/requisitions' },
  // G-6/G-12 sub-1: Candidates — profile terstruktur (education/experience/skills/cert/document/consent)
  { labelKey: 'recruitment.candidates', descKey: 'candidates.description', icon: 'pi pi-users', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/recruitment/candidates' },
  // G-12 sub-3: Applications — pipeline (history/screening/assessment/interviews/match score)
  { labelKey: 'recruitment.applications', descKey: 'applications.description', icon: 'pi pi-send', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/recruitment/applications' },
  // S-4: Internal Candidates — eligible via career path (Career Intelligence)
  { labelKey: 'internal_candidates.title', descKey: 'internal_candidates.description', icon: 'pi pi-user-plus', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400', route: '/recruitment/internal-candidates' },
  // G-3: Job Offers — offer management + approval workflow
  { labelKey: 'recruitment.offers', descKey: 'offers.description', icon: 'pi pi-file-edit', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400', route: '/recruitment/offers' },
  // G-4: Onboarding — employee hasil offer + status (COMPLETED → training handoff S-7)
  { labelKey: 'recruitment.onboarding', descKey: 'onboarding.description', icon: 'pi pi-sign-in', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400', route: '/recruitment/onboarding' },
  // G-7 sub-2: Assessments — batch session + peserta kandidat
  { labelKey: 'recruitment.assessments', descKey: 'assessments.description', icon: 'pi pi-clipboard', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400', route: '/recruitment/assessments' }
])

const comingSoonCards = computed(() => [])

// G-11: recruitment analytics summary — fail-silent (kartu tidak tampil kalau error)
const statsLoading = ref(true)
const stats = ref(null)

const statCards = computed(() => {
  if (!stats.value) return []
  const cards = [
    { labelKey: 'recruitment_hub.open_requisitions', icon: 'pi pi-briefcase', color: '#0284c7', value: stats.value.open_requisitions ?? 0 },
    { labelKey: 'recruitment_hub.candidates', icon: 'pi pi-users', color: '#059669', value: stats.value.candidates ?? 0 },
    { labelKey: 'recruitment_hub.applications', icon: 'pi pi-send', color: '#d97706', value: stats.value.applications ?? 0 },
    { labelKey: 'recruitment_hub.interviews', icon: 'pi pi-comments', color: '#4f46e5', value: stats.value.interviews ?? 0 },
    { labelKey: 'recruitment_hub.offers', icon: 'pi pi-file-edit', color: '#4338ca', value: stats.value.offers ?? 0 },
    { labelKey: 'recruitment_hub.hires', icon: 'pi pi-verified', color: '#16a34a', value: stats.value.hires ?? 0 }
  ]
  if (stats.value.time_to_hire_days !== null && stats.value.time_to_hire_days !== undefined) {
    cards.push({ labelKey: 'recruitment_hub.time_to_hire', icon: 'pi pi-clock', color: '#db2777', value: `${Math.round(stats.value.time_to_hire_days)}d` })
  }
  return cards
})

async function loadStats() {
  statsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/recruitment/analytics/summary')
    stats.value = res.data?.data || null
  } catch {
    // fail-silent — hub tetap tampil tanpa summary cards
    stats.value = null
  } finally {
    statsLoading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>
