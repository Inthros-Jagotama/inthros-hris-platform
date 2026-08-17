<template>
  <div class="space-y-6">
    <!-- ── Card navigasi, dikelompokkan per kategori: Pengaturan / Assessment / Laporan ── -->
    <div v-for="group in menuGroups" :key="group.titleKey" class="space-y-2">
      <div class="md:col-span-2">
        <div class="flex items-center gap-2 pt-2">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase">{{ t(group.titleKey) }}</span>
          <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
        </div>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <button
          v-for="menu in group.items"
          :key="menu.route"
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-indigo-300 dark:hover:border-indigo-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/50"
          @click="router.push(menu.route)"
        >
          <div class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors" :class="menu.tint">
            <i :class="menu.icon" class="text-base"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ t(menu.titleKey) }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(menu.descKey) }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-indigo-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>
      </div>
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
// Semua card dicek dengan hasExactPermission (tanpa fallback module-level
// competency.view) — sama seperti halaman attendance/training: mencabut
// competency.settings.view / assessment.view / report.view dari sebuah role
// otomatis menyembunyikan grup terkait.
const { hasExactPermission } = useAuth()

const menuGroups = computed(() => {
  const groups = [
    {
      titleKey: 'competency_360.group_settings',
      items: [
        { route: '/competencies/values', titleKey: 'competency_360.rating_scales', descKey: 'competency_360.rating_scales_desc', icon: 'pi pi-sliders-h', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400', permission: 'competency.settings.view' },
        { route: '/competencies/indicators', titleKey: 'competency_360.indicators', descKey: 'competency_360.indicators_desc', icon: 'pi pi-list', tint: 'bg-fuchsia-50 dark:bg-fuchsia-500/10 text-fuchsia-600 dark:text-fuchsia-400', permission: 'competency.settings.view' },
        { route: '/competencies/templates', titleKey: 'competency_360.templates', descKey: 'competency_360.templates_desc', icon: 'pi pi-clone', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400', permission: 'competency.settings.view' },
        { route: '/competencies/events', titleKey: 'competency_360.events', descKey: 'competency_360.events_desc', icon: 'pi pi-calendar', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', permission: 'competency.settings.view' },
        { route: '/competencies/raters', titleKey: 'competency_360.rater_assignment', descKey: 'competency_360.rater_assignment_desc', icon: 'pi pi-users', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', permission: 'competency.settings.view' }
      ]
    },
    {
      titleKey: 'competency_360.group_assessment',
      items: [
        { route: '/competencies/my-assessments', titleKey: 'competency_360.my_assessments', descKey: 'competency_360.my_assessments_desc', icon: 'pi pi-list-check', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', permission: 'competency.assessment.view' },
        { route: '/competencies/manager-assessments', titleKey: 'competency_360.manager_assessments', descKey: 'competency_360.manager_assessments_desc', icon: 'pi pi-user-edit', tint: 'bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400', permission: 'competency.assessment.view' }
      ]
    },
    {
      titleKey: 'competency_360.group_reports',
      items: [
        { route: '/competencies/results', titleKey: 'competency_360.results', descKey: 'competency_360.results_desc', icon: 'pi pi-chart-bar', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400', permission: 'competency.report.view' },
        { route: '/competencies/reports', titleKey: 'competency_360.reports', descKey: 'competency_360.reports_desc', icon: 'pi pi-chart-line', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400', permission: 'competency.report.view' }
      ]
    }
  ]
  return groups
    .map(g => ({ ...g, items: g.items.filter(card => hasExactPermission(card.permission)) }))
    .filter(g => g.items.length > 0)
})
</script>
