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
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ t(menu.labelKey) }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(menu.descKey) }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-violet-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
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

import Tag from 'primevue/tag'

const router = useRouter()
const { t } = useI18n()

// Menu strategis Career Intelligence — halaman hub dengan card navigasi.
// Career Paths berada di sini (bukan di sidebar) sesuai keputusan pemisahan
// transactional vs strategical: halaman ini untuk konfigurasi/perencanaan karir.
const menuCards = computed(() => [
  { labelKey: 'career_intel.paths', descKey: 'career_paths.description', icon: 'pi pi-sitemap', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400', route: '/career-intelligence/paths' },
  // S-5: Succession Gaps — posisi kunci tanpa successor siap → external recruitment
  { labelKey: 'succession_gaps.title', descKey: 'succession_gaps.description', icon: 'pi pi-users', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/career-intelligence/successions' },
  // P2-FE integrasi: riwayat pengembangan (course–competency) dari modul Training
  { labelKey: 'career_intel.development', descKey: 'career_intel.development_desc', icon: 'pi pi-graduation-cap', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/training/history' }
])

// Fitur yang sudah terdefinisi di locale tapi halamannya belum dibangun —
// ditampilkan sebagai card "Coming soon" agar roadmap modul terlihat jelas.
const comingSoonCards = computed(() => [
  { labelKey: 'career_intel.talent_map', icon: 'pi pi-th-large', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' },
  { labelKey: 'career_intel.interests', icon: 'pi pi-compass', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' }
])
</script>
