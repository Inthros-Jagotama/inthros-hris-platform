<template>
  <div class="space-y-4">

    <!-- Menu cards (pola sama dengan halaman Settings / Attendance) -->
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
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'

import Tag from 'primevue/tag'

const router = useRouter()
const { t } = useI18n()
const { hasPermission } = useAuth()

// Menu strategis Workforce Intelligence — halaman hub dengan card navigasi.
// Candidate Search menampilkan posisi kosong (org tanpa employment aktif di
// bawah Organization Summary active) beserta kandidat recruitment-nya.
const menuCards = computed(() => [
  { labelKey: 'candidate_search.title', descKey: 'candidate_search.description', icon: 'pi pi-user-plus', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', route: '/workforce-intelligence/candidate-search', permission: 'workforceintelligence.candidate-search.view' },
  // Recruitment Analytics — S-2/S-3: remaining gap, expected hires, time to hire/fill, OAR, source conversion
  { labelKey: 'recruitment_analytics.title', descKey: 'recruitment_analytics.description', icon: 'pi pi-chart-line', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/workforce-intelligence/recruitment-analytics', permission: 'workforceintelligence.recruitment-analytics.view' },
  // Quality of Hire — metrik agregat kualitas hire (S-6: interview/onboarding/performance/retention)
  { labelKey: 'quality_of_hire.title', descKey: 'quality_of_hire.description', icon: 'pi pi-bullseye', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/workforce-intelligence/quality-of-hire', permission: 'workforceintelligence.quality-of-hire.view' },
  // P2-FE integrasi: analisis training (completion/cost/compliance) dari modul Training
  { labelKey: 'workforce_intel.training_analysis', descKey: 'workforce_intel.training_analysis_desc', icon: 'pi pi-graduation-cap', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400', route: '/training/reports', permission: 'workforceintelligence.training-analysis.view' }
].filter(card => !card.permission || hasPermission(card.permission)))

// Fitur yang sudah terdefinisi di locale tapi halamannya belum dibangun —
// ditampilkan sebagai card "Coming soon" agar roadmap modul terlihat jelas.
const comingSoonCards = computed(() => [
  { labelKey: 'workforce_intel.headcount_planning', icon: 'pi pi-arrows-alt', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' },
  { labelKey: 'workforce_intel.risk', icon: 'pi pi-exclamation-triangle', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400' },
  { labelKey: 'workforce_intel.executive', icon: 'pi pi-desktop', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400' },
  { labelKey: 'workforce_intel.scenarios', icon: 'pi pi-clone', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' }
])
</script>
