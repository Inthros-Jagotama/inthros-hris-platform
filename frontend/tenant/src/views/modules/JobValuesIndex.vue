<template>
  <div class="space-y-6">
    <!-- Grouped Cards -->
    <div v-for="group in filteredGroups" :key="group.key" class="space-y-3">
      <div class="flex items-center gap-2">
        <i :class="group.icon" class="text-indigo-400 text-sm"></i>
        <h2 class="text-sm font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t(group.labelKey) }}</h2>
        <span class="text-xs text-gray-400 dark:text-gray-500">{{ group.items.length }}</span>
        <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
        <button
          v-for="item in group.items"
          :key="item.type"
          type="button"
          class="group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-indigo-300 dark:hover:border-indigo-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/50"
          @click="openType(item.type)"
        >
          <div
            class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors"
            :class="item.tint"
          >
            <i :class="item.icon" class="text-base"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ typeLabel(item.type) }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(`job_values.type_desc.${item.type}`) }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-indigo-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>
      </div>
    </div>

    <!-- Empty search result -->
    <div v-if="filteredGroups.length === 0" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-search text-3xl mb-2 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('settings.no_results') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { jobValueTypeLabel } from '@/utils/jobValues'

const router = useRouter()
const { t } = useI18n()

const searchQuery = ref('')

// ── Definisi semua tipe job value (sub-menu), dikelompokkan ──
const groups = computed(() => [
  {
    key: 'basic',
    icon: 'pi pi-user',
    labelKey: 'job_values.group_basic',
    items: [
      { type: 'education', icon: 'pi pi-graduation-cap', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' },
      { type: 'experience', icon: 'pi pi-briefcase', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400' },
      { type: 'subordinate', icon: 'pi pi-sitemap', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400' },
      { type: 'activity', icon: 'pi pi-directions', tint: 'bg-lime-50 dark:bg-lime-500/10 text-lime-600 dark:text-lime-400' },
      { type: 'communicating_influencing_skill', icon: 'pi pi-comments', tint: 'bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400' },
      { type: 'thinking_environment', icon: 'pi pi-wave-pulse', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400' },
      { type: 'thinking_chalenge', icon: 'pi pi-flag', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400' }
    ]
  },
  {
    key: 'psychological',
    icon: 'pi pi-heart',
    labelKey: 'job_values.group_psychological',
    items: [
      { type: 'kecerdasan', icon: 'pi pi-microchip', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' },
      { type: 'innovation_creativity', icon: 'pi pi-lightbulb', tint: 'bg-yellow-50 dark:bg-yellow-500/10 text-yellow-600 dark:text-yellow-400' },
      { type: 'self_confidence', icon: 'pi pi-user', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' },
      { type: 'flexibility', icon: 'pi pi-refresh', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400' },
      { type: 'tenacity', icon: 'pi pi-bolt', tint: 'bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400' },
      { type: 'continuous_learning', icon: 'pi pi-book', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400' }
    ]
  },
  {
    key: 'technical',
    icon: 'pi pi-cog',
    labelKey: 'job_values.group_technical',
    items: [
      { type: 'competency_based_human_resources_management', icon: 'pi pi-users', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' },
      { type: 'competency_development', icon: 'pi pi-chart-line', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400' },
      { type: 'people_development', icon: 'pi pi-user-plus', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400' },
      { type: 'career_management', icon: 'pi pi-sitemap', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400' },
      { type: 'hr_assessment', icon: 'pi pi-clipboard', tint: 'bg-fuchsia-50 dark:bg-fuchsia-500/10 text-fuchsia-600 dark:text-fuchsia-400' },
      { type: 'recruitement_selection', icon: 'pi pi-search', tint: 'bg-pink-50 dark:bg-pink-500/10 text-pink-600 dark:text-pink-400' },
      { type: 'job_analysis_evaluation', icon: 'pi pi-briefcase', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400' },
      { type: 'organizational_development', icon: 'pi pi-building', tint: 'bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400' },
      { type: 'human_resources_information_system', icon: 'pi pi-database', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' },
      { type: 'workload_analysis', icon: 'pi pi-chart-bar', tint: 'bg-yellow-50 dark:bg-yellow-500/10 text-yellow-600 dark:text-yellow-400' },
      { type: 'performance_apraisal', icon: 'pi pi-verified', tint: 'bg-lime-50 dark:bg-lime-500/10 text-lime-600 dark:text-lime-400' },
      { type: 'remuneration_manajemen', icon: 'pi pi-wallet', tint: 'bg-green-50 dark:bg-green-500/10 text-green-600 dark:text-green-400' },
      { type: 'reward_punisment_management', icon: 'pi pi-gift', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' },
      { type: 'health_safety_environment', icon: 'pi pi-heart', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400' },
      { type: 'hubungan_industrial', icon: 'pi pi-comments', tint: 'bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400' },
      { type: 'budgeting', icon: 'pi pi-dollar', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400' }
    ]
  },
  {
    key: 'managerial',
    icon: 'pi pi-briefcase',
    labelKey: 'job_values.group_managerial',
    items: [
      { type: 'integrity', icon: 'pi pi-shield', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' },
      { type: 'achievement_orientation', icon: 'pi pi-bullseye', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400' },
      { type: 'building_partnership', icon: 'pi pi-link', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400' },
      { type: 'planning_organizing', icon: 'pi pi-calendar', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400' },
      { type: 'leadership', icon: 'pi pi-star', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400' },
      { type: 'developing_others', icon: 'pi pi-users', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' }
    ]
  },
  {
    key: 'work_environment',
    icon: 'pi pi-globe',
    labelKey: 'job_values.group_environment',
    items: [
      { type: 'environment', icon: 'pi pi-globe', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' },
      { type: 'risk', icon: 'pi pi-exclamation-triangle', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400' },
      { type: 'relationship', icon: 'pi pi-users', tint: 'bg-pink-50 dark:bg-pink-500/10 text-pink-600 dark:text-pink-400' },
      { type: 'frequency', icon: 'pi pi-clock', tint: 'bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400' }
    ]
  },
  {
    key: 'authority_asset',
    icon: 'pi pi-shield',
    labelKey: 'job_values.group_authority',
    items: [
      { type: 'asset', icon: 'pi pi-box', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400' },
      { type: 'asset_authority', icon: 'pi pi-box', tint: 'bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400' },
      { type: 'authority', icon: 'pi pi-shield', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400' },
      { type: 'authority_unauthorized', icon: 'pi pi-ban', tint: 'bg-slate-100 dark:bg-slate-500/10 text-slate-600 dark:text-slate-400' }
    ]
  },
  {
    key: 'financial',
    icon: 'pi pi-dollar',
    labelKey: 'job_values.group_financial',
    items: [
      { type: 'cash', icon: 'pi pi-dollar', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400' },
      { type: 'impact', icon: 'pi pi-bolt', tint: 'bg-fuchsia-50 dark:bg-fuchsia-500/10 text-fuchsia-600 dark:text-fuchsia-400' },
      { type: 'impact_unauthorized', icon: 'pi pi-lock', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400' }
    ]
  }
])

function typeLabel(value) {
  return jobValueTypeLabel(t, value)
}

// Filter groups by search query
const filteredGroups = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return groups.value

  return groups.value
    .map(group => ({
      ...group,
      items: group.items.filter(item =>
        typeLabel(item.type).toLowerCase().includes(q) ||
        t(`job_values.type_desc.${item.type}`).toLowerCase().includes(q)
      )
    }))
    .filter(group => group.items.length > 0)
})

function openType(type) {
  router.push(`/job-management/values/${type}`)
}
</script>
