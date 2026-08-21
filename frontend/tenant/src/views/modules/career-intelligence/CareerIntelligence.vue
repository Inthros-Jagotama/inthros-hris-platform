<template>
  <div class="space-y-4">
    <!-- Menu cards (pola sama dengan menu Settings / Attendance) -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      <button
        v-for="menu in menuCards"
        :key="menu.route"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-violet-300 dark:hover:border-violet-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/50"
        @click="router.push(menu.route)"
      >
        <div
          class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors"
          :class="menu.tint"
        >
          <i :class="menu.icon" class="text-base"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-navy-800 dark:text-gray-100 truncate">{{ t(menu.labelKey) }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(menu.descKey) }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-violet-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'

const router = useRouter()
const { t } = useI18n()
const { hasPermission } = useAuth()

// Menu strategis Career Intelligence — halaman hub dengan card navigasi.
// Career Paths berada di sini (bukan di sidebar) sesuai keputusan pemisahan
// transactional vs strategical: halaman ini untuk konfigurasi/perencanaan karir.
const menuCards = computed(() => [
  { labelKey: 'career_intel.talent_map', descKey: 'talent_maps.description', icon: 'pi pi-th-large', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400', route: '/career-intelligence/talent-maps', permission: 'careerintelligence.view' },
  { labelKey: 'career_intel.paths', descKey: 'career_paths.description', icon: 'pi pi-sitemap', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400', route: '/career-intelligence/paths', permission: 'careerintelligence.paths.view' },
  { labelKey: 'career_intel.interests', descKey: 'career_interests.description', icon: 'pi pi-compass', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/career-intelligence/interests', permission: 'careerintelligence.view' },
  // S-5: Succession Gaps — posisi kunci tanpa successor siap → external recruitment
  { labelKey: 'career_intel.succession_plans', descKey: 'succession_plans.description', icon: 'pi pi-user-plus', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400', route: '/career-intelligence/succession-plans', permission: 'careerintelligence.view' },
  { labelKey: 'succession_gaps.title', descKey: 'succession_gaps.description', icon: 'pi pi-users', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/career-intelligence/successions', permission: 'careerintelligence.successions.view' },
  { labelKey: 'gap_analysis.title', descKey: 'gap_analysis.description', icon: 'pi pi-sitemap', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', route: '/career-intelligence/gap-analysis', permission: 'careerintelligence.view' },
  // P2-FE integrasi: riwayat pengembangan (course–competency) dari modul Training
  { labelKey: 'career_intel.development', descKey: 'career_intel.development_desc', icon: 'pi pi-graduation-cap', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/training/history', permission: 'careerintelligence.development.view' }
].filter(card => !card.permission || hasPermission(card.permission)))
</script>
