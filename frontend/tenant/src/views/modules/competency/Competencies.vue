<template>
  <div class="space-y-6">
    <!-- ── Menu Cards ── -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      <button
        v-if="hasPermission('competency.settings.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-indigo-300 dark:hover:border-indigo-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/50"
        @click="router.push('/competencies/values')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400">
          <i class="pi pi-sliders-h text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.rating_scales') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.rating_scales_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-indigo-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.settings.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-fuchsia-300 dark:hover:border-fuchsia-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-fuchsia-500/50"
        @click="router.push('/competencies/indicators')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-fuchsia-50 dark:bg-fuchsia-500/10 text-fuchsia-600 dark:text-fuchsia-400">
          <i class="pi pi-list text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.indicators') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.indicators_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-fuchsia-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.settings.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-violet-300 dark:hover:border-violet-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/50"
        @click="router.push('/competencies/templates')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400">
          <i class="pi pi-clone text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.templates') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.templates_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-violet-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.settings.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-emerald-300 dark:hover:border-emerald-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
        @click="router.push('/competencies/events')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
          <i class="pi pi-calendar text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.events') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.events_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.settings.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-amber-300 dark:hover:border-amber-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-amber-500/50"
        @click="router.push('/competencies/raters')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400">
          <i class="pi pi-users text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.rater_assignment') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.rater_assignment_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-amber-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.assessment.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-sky-300 dark:hover:border-sky-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-sky-500/50"
        @click="router.push('/competencies/my-assessments')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400">
          <i class="pi pi-list-check text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.my_assessments') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.my_assessments_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-sky-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.assessment.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-cyan-300 dark:hover:border-cyan-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-500/50"
        @click="router.push('/competencies/manager-assessments')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400">
          <i class="pi pi-user-edit text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.manager_assessments') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.manager_assessments_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-cyan-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.report.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-rose-300 dark:hover:border-rose-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-rose-500/50"
        @click="router.push('/competencies/results')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400">
          <i class="pi pi-chart-bar text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.results') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.results_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-rose-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>

      <button
        v-if="hasPermission('competency.report.view')"
        type="button"
        class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-teal-300 dark:hover:border-teal-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-teal-500/50"
        @click="router.push('/competencies/reports')"
      >
        <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400">
          <i class="pi pi-chart-line text-lg"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('competency_360.reports') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('competency_360.reports_desc') }}</p>
        </div>
        <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-teal-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
      </button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'

const router = useRouter()
const { t } = useI18n()
const { hasPermission } = useAuth()
</script>
