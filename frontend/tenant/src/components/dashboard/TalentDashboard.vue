<template>
  <div>
    <div v-if="assessLoading" class="space-y-4">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
        <div v-for="i in 4" :key="i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
          <div class="h-4 w-28 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
          <div class="h-8 w-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
        </div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div v-for="i in 2" :key="'d' + i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
          <div class="h-4 w-40 bg-gray-200 dark:bg-gray-700 rounded mb-3"></div>
          <div class="w-32 h-32 mx-auto rounded-full bg-gray-200 dark:bg-gray-700"></div>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
        <div class="flex items-center gap-2">
          <i class="pi pi-user text-sm text-teal-500"></i>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.view_talent') }}</h2>
          <span class="text-[11px] text-gray-400 dark:text-gray-500 hidden sm:inline">{{ t('dashboard.talent_desc') }}</span>
        </div>
        <div class="flex items-center gap-2">
          <Tag v-if="kpiPeriod" :value="kpiPeriod.period_code + ' · ' + kpiPeriod.year" severity="info" class="!text-xs !px-2 !py-0.5" />
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/performance')" />
        </div>
      </div>

      <!-- Pipeline proses assessment (KPI) -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 mb-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.assessment_pipeline') }}</h2>
          <span class="text-xs font-semibold text-teal-600 dark:text-teal-400">{{ kpiRate }}%</span>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 mb-3">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_draft') }}</p>
            <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ kpiStats.draft_count ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-2.5">
            <p class="text-[11px] text-amber-600 dark:text-amber-400 uppercase tracking-wider">{{ t('dashboard.assessment_submitted') }}</p>
            <p class="text-lg font-bold text-amber-700 dark:text-amber-300">{{ kpiStats.submitted_count ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-2.5">
            <p class="text-[11px] text-sky-600 dark:text-sky-400 uppercase tracking-wider">{{ t('dashboard.assessment_approved') }}</p>
            <p class="text-lg font-bold text-sky-700 dark:text-sky-300">{{ kpiStats.approved_count ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-emerald-300 dark:border-emerald-700/50 bg-emerald-50/50 dark:bg-emerald-900/10 p-2.5">
            <p class="text-[11px] text-emerald-600 dark:text-emerald-400 uppercase tracking-wider">{{ t('dashboard.assessment_completed') }}</p>
            <p class="text-lg font-bold text-emerald-700 dark:text-emerald-300">{{ kpiStats.completed_count ?? 0 }}</p>
          </div>
        </div>
        <div class="h-2 bg-gray-100 dark:bg-gray-700 rounded-full overflow-hidden">
          <div class="h-2 bg-teal-500 rounded-full transition-all" :style="{ width: Math.min(100, kpiRate) + '%' }"></div>
        </div>
        <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-1.5">{{ t('dashboard.assessment_rate') }}: {{ kpiRate }}%</p>
      </div>

      <!-- KPI summary -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 mb-3">
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_total') }}</span>
            <i class="pi pi-users text-lg text-teal-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ kpiStats.total_employees ?? 0 }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_in_progress') }}</span>
            <i class="pi pi-spinner text-lg text-amber-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ kpiInProgress }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_average_score') }}</span>
            <i class="pi pi-chart-line text-lg text-emerald-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ fmtScore(kpiStats.average_score) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_average_achievement') }}</span>
            <i class="pi pi-bullseye text-lg text-cyan-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ fmtPct(kpiStats.average_achievement) }}</div>
        </div>
      </div>

      <!-- OKR summary -->
      <div v-if="okrHR" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 mb-3">
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_okr_total') }}</span>
            <i class="pi pi-target text-lg text-violet-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ okrHR.total_evaluations ?? 0 }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_okr_completed') }}</span>
            <i class="pi pi-check-circle text-lg text-emerald-500"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ okrHR.completed_count ?? 0 }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_average_score') }}</span>
            <i class="pi pi-chart-line text-lg text-violet-400"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ fmtScore(okrHR.average_score) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-2">
            <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.assessment_average_achievement') }}</span>
            <i class="pi pi-bullseye text-lg text-violet-300"></i>
          </div>
          <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ fmtPct(okrHR.average_achievement) }}</div>
        </div>
      </div>

      <!-- Donut: distribusi rating KPI & OKR -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
        <div v-if="kpiRatingSegments.length" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.assessment_rating_kpi') }}</h2>
          <div class="flex items-center gap-6 flex-wrap">
            <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
              <circle
                v-for="seg in kpiRatingSegments"
                :key="seg.label"
                cx="60" cy="60" r="45" fill="none" stroke-width="18"
                :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
                transform="rotate(-90 60 60)"
              />
              <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ kpiRatingTotal }}</text>
              <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.assessment_completed') }}</text>
            </svg>
            <div class="space-y-2 flex-1 min-w-0">
              <div v-for="seg in kpiRatingSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
                <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
                <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
                <span class="font-semibold text-navy-800 dark:text-gray-100">{{ seg.value }}</span>
                <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="okrRatingSegments.length" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.assessment_rating_okr') }}</h2>
          <div class="flex items-center gap-6 flex-wrap">
            <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
              <circle
                v-for="seg in okrRatingSegments"
                :key="seg.label"
                cx="60" cy="60" r="45" fill="none" stroke-width="18"
                :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
                transform="rotate(-90 60 60)"
              />
              <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ okrRatingTotal }}</text>
              <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.assessment_completed') }}</text>
            </svg>
            <div class="space-y-2 flex-1 min-w-0">
              <div v-for="seg in okrRatingSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
                <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
                <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
                <span class="font-semibold text-navy-800 dark:text-gray-100">{{ seg.value }}</span>
                <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Performer terbaik (KPI) -->
      <div v-if="kpiTop.length" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-trophy text-sm text-amber-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.assessment_top_performers') }}</h2>
          </div>
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/performance')" />
        </div>
        <div class="space-y-2">
          <div v-for="p in kpiTop" :key="p.employee_id" class="flex items-center gap-3 text-sm">
            <span class="w-6 h-6 rounded-full bg-amber-100 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400 flex items-center justify-center text-xs font-bold shrink-0">{{ p.rank }}</span>
            <div class="min-w-0 flex-1">
              <p class="text-gray-700 dark:text-gray-200 truncate">{{ p.employee_name || '—' }}</p>
              <p class="text-[11px] text-gray-400 dark:text-gray-500 truncate">{{ p.organization_name || '' }}</p>
            </div>
            <span v-if="p.rating_name" class="text-xs font-semibold px-2 py-0.5 rounded-full shrink-0" :style="{ color: p.rating_color || '#8b5cf6', backgroundColor: (p.rating_color || '#8b5cf6') + '1a' }">{{ p.rating_name }}</span>
            <span class="font-semibold text-navy-800 dark:text-gray-100 shrink-0">{{ fmtScore(p.final_score) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'
import { fmtScore, fmtPct, buildDonutSegments } from '@/utils/dashboard'

import Button from 'primevue/button'
import Tag from 'primevue/tag'

const { t } = useI18n()
const toast = useToast()

const assessLoading = ref(false)
const kpiHR = ref(null) // /performance/kpi/dashboard/hr
const okrHR = ref(null) // /performance/okr/dashboard/hr

const kpiPeriod = computed(() => kpiHR.value?.current_period || null)
const kpiStats = computed(() => kpiHR.value?.completion_stats || {})
const kpiInProgress = computed(() => (Number(kpiStats.value.submitted_count) || 0) + (Number(kpiStats.value.approved_count) || 0))
const kpiRate = computed(() => {
  const total = Number(kpiStats.value.total_employees) || 0
  if (!total) return 0
  return Math.round(((Number(kpiStats.value.completed_count) || 0) / total) * 100)
})
const RATING_FALLBACK_COLORS = ['#3b82f6', '#ec4899', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#ef4444', '#84cc16', '#f97316', '#6366f1']
const kpiRatingSegments = computed(() => buildDonutSegments(
  (kpiHR.value?.rating_distribution || []).map((r, i) => ({ label: r.rating_name || r.rating_code, value: r.count, color: r.rating_color || RATING_FALLBACK_COLORS[i % RATING_FALLBACK_COLORS.length] }))
))
const kpiRatingTotal = computed(() => kpiRatingSegments.value.reduce((s, i) => s + i.value, 0))
const okrRatingSegments = computed(() => buildDonutSegments(
  (okrHR.value?.rating_distribution || []).map((r, i) => ({ label: r.rating_name, value: r.count, color: r.color || RATING_FALLBACK_COLORS[i % RATING_FALLBACK_COLORS.length] }))
))
const okrRatingTotal = computed(() => okrRatingSegments.value.reduce((s, i) => s + i.value, 0))
const kpiTop = computed(() => (kpiHR.value?.top_performers || []).slice(0, 5))

async function loadTalentDashboard() {
  if (assessLoading.value) return
  assessLoading.value = true
  try {
    const [kpiRes, okrRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/performance/kpi/dashboard/hr'),
      api.get('/api/v1/tenant/performance/okr/dashboard/hr')
    ])
    kpiHR.value = kpiRes.status === 'fulfilled' ? (kpiRes.value.data?.data || null) : null
    okrHR.value = okrRes.status === 'fulfilled' ? (okrRes.value.data?.data || null) : null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    assessLoading.value = false
  }
}

onMounted(loadTalentDashboard)
</script>
