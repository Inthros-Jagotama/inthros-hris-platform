<template>
  <div class="space-y-6">
    <!-- Quick Stats -->
    <div v-if="!statsLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-blue-50 dark:bg-blue-500/10 flex items-center justify-center">
            <i class="pi pi-users text-blue-600 dark:text-blue-400"></i>
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stats.totalEmployees }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('kpi.total_employees') }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-emerald-50 dark:bg-emerald-500/10 flex items-center justify-center">
            <i class="pi pi-check-circle text-emerald-600 dark:text-emerald-400"></i>
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stats.completedCount }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('kpi.completed') }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-amber-50 dark:bg-amber-500/10 flex items-center justify-center">
            <i class="pi pi-clock text-amber-600 dark:text-amber-400"></i>
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stats.inProgressCount }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('kpi.in_progress') }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-purple-50 dark:bg-purple-500/10 flex items-center justify-center">
            <i class="pi pi-chart-line text-purple-600 dark:text-purple-400"></i>
          </div>
          <div>
            <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stats.averageScore.toFixed(1) }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('kpi.avg_score') }}</p>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="n in 4" :key="n" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="flex items-center gap-3">
          <div class="w-10 h-10 rounded-lg bg-gray-200 dark:bg-gray-700 animate-pulse"></div>
          <div class="space-y-2">
            <div class="h-6 w-16 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
            <div class="h-3 w-20 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Menu Cards -->
    <div v-for="group in menuGroups" :key="group.key" class="space-y-3">
      <div class="flex items-center gap-2">
        <h3 class="text-xs font-semibold uppercase tracking-wide" :class="group.highlight ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 dark:text-gray-500'">{{ t(group.titleKey) }}</h3>
        <i v-if="group.highlight" class="pi pi-star-fill text-emerald-500 text-[10px]"></i>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <button
          v-for="item in group.items"
          :key="item.path"
          type="button"
          :disabled="item.muted"
          class="group flex items-center gap-4 rounded-xl text-left transition-all focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
          :class="item.muted
            ? 'p-5 border border-dashed border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 grayscale opacity-60 cursor-not-allowed'
            : group.highlight
              ? 'p-5 cursor-pointer border-2 border-emerald-300 dark:border-emerald-500/60 bg-emerald-50/60 dark:bg-emerald-500/10 shadow-md shadow-emerald-100 dark:shadow-none hover:shadow-xl hover:shadow-emerald-100/80 dark:hover:shadow-none hover:-translate-y-1 hover:border-emerald-400 dark:hover:border-emerald-400'
              : 'p-5 cursor-pointer border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:border-emerald-300 dark:hover:border-emerald-500/60 hover:shadow-lg hover:-translate-y-1'"
          @click="!item.muted && navigateTo(item.path)"
        >
          <div
            class="w-14 h-14 rounded-xl shrink-0 flex items-center justify-center transition-colors"
            :class="item.tint"
          >
            <i :class="item.icon" class="text-xl"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-base font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              {{ t(item.titleKey) }}
              <span v-if="item.muted" class="text-[10px] font-normal uppercase tracking-wide text-gray-400 dark:text-gray-500 border border-gray-300 dark:border-gray-600 rounded px-1.5 py-0.5">{{ t('performance.not_available_yet') }}</span>
            </p>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1" :class="group.highlight && !item.muted ? '' : 'line-clamp-2'">{{ item.muted ? t(item.mutedDescKey) : t(item.descKey) }}</p>
          </div>
          <i class="pi pi-chevron-right text-gray-300 dark:text-gray-600 group-hover:text-emerald-500 group-hover:translate-x-1 transition-all shrink-0"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'

const router = useRouter()
const { t } = useI18n()

const statsLoading = ref(true)
const stats = ref({
  totalEmployees: 0,
  completedCount: 0,
  inProgressCount: 0,
  averageScore: 0
})

const kpiSelfAssessmentAvailable = ref(true)
const okrSelfAssessmentAvailable = ref(true)

const menuGroups = computed(() => [
  {
    key: 'self_service',
    titleKey: 'performance.group_self_service',
    highlight: true,
    items: [
      {
        path: '/performance/kpi/my-evaluation',
        icon: 'pi pi-user-edit',
        titleKey: 'kpi.my_evaluation',
        descKey: 'kpi.my_evaluation_desc',
        mutedDescKey: 'performance.self_assessment_unavailable_desc',
        tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400',
        muted: !kpiSelfAssessmentAvailable.value
      },
      {
        path: '/performance/okr/my-evaluation',
        icon: 'pi pi-user-edit',
        titleKey: 'okr.my_evaluation',
        descKey: 'okr.my_evaluation_desc',
        mutedDescKey: 'performance.self_assessment_unavailable_desc',
        tint: 'bg-pink-50 dark:bg-pink-500/10 text-pink-600 dark:text-pink-400',
        muted: !okrSelfAssessmentAvailable.value
      }
    ]
  },
  {
    key: 'kpi',
    titleKey: 'performance.group_kpi',
    items: [
      {
        path: '/performance/kpi',
        icon: 'pi pi-chart-bar',
        titleKey: 'kpi.evaluations',
        descKey: 'kpi.evaluations_desc',
        tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
      },
      {
        path: '/performance/kpi/templates',
        icon: 'pi pi-file',
        titleKey: 'kpi.templates',
        descKey: 'kpi.templates_desc',
        tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400'
      }
    ]
  },
  {
    key: 'okr',
    titleKey: 'performance.group_okr',
    items: [
      {
        path: '/performance/okr',
        icon: 'pi pi-bullseye',
        titleKey: 'performance.okr_evaluations',
        descKey: 'performance.okr_evaluations_desc',
        tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400'
      },
      {
        path: '/performance/okr/templates',
        icon: 'pi pi-file-edit',
        titleKey: 'performance.okr_templates',
        descKey: 'performance.okr_templates_desc',
        tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400'
      }
    ]
  },
  {
    key: 'shared',
    titleKey: 'performance.group_shared',
    items: [
      {
        path: '/performance/kpi/periods',
        icon: 'pi pi-calendar',
        titleKey: 'kpi.periods',
        descKey: 'performance.periods_shared_desc',
        tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400'
      }
    ]
  }
])

function navigateTo(path) {
  router.push(path)
}

async function loadStats() {
  statsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/dashboard/hr')
    const data = res.data?.data || res.data
    if (data?.completion_stats) {
      stats.value = {
        totalEmployees: data.completion_stats.total_employees || 0,
        completedCount: data.completion_stats.completed_count + data.completion_stats.approved_count || 0,
        inProgressCount: data.completion_stats.submitted_count + data.completion_stats.draft_count || 0,
        averageScore: data.completion_stats.average_score || 0
      }
    }
  } catch {
    // Silently fail - stats are optional
  } finally {
    statsLoading.value = false
  }
}

async function loadSelfAssessmentAvailability() {
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/my-context')
    const data = res.data?.data || res.data
    kpiSelfAssessmentAvailable.value = !!data?.has_position && (data?.templates?.length || 0) > 0
  } catch {
    // Silently fail - availability check is optional
  }
  try {
    const res = await api.get('/api/v1/tenant/performance/okr/my-context')
    const data = res.data?.data || res.data
    okrSelfAssessmentAvailable.value = !!data?.has_position && (data?.templates?.length || 0) > 0
  } catch {
    // Silently fail - availability check is optional
  }
}

onMounted(() => {
  loadStats()
  loadSelfAssessmentAvailability()
})
</script>
