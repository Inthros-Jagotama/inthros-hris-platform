<template>
  <div class="space-y-3">
    <div class="flex items-center gap-2 mb-3">
      <i class="pi pi-calendar-clock text-sm text-teal-500"></i>
      <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.view_hr') }}</h2>
      <span class="text-[11px] text-gray-400 dark:text-gray-500 hidden sm:inline">{{ t('dashboard.hr_desc') }}</span>
    </div>

    <!-- KPI row -->
    <div v-if="hrLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
      <div v-for="i in 4" :key="i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
        <div class="h-4 w-28 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
        <div class="h-8 w-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.on_leave_today') }}</span>
          <i class="pi pi-calendar text-lg text-amber-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ hrOnLeaveToday ?? '—' }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ t('dashboard.hr_month_label') }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.hr_present') }}</span>
          <i class="pi pi-check-circle text-lg text-emerald-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.present ?? '—' }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ hrTotalWorkHours }} {{ t('dashboard.hr_total_work_hours').toLowerCase() }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.hr_late') }}</span>
          <i class="pi pi-clock text-lg text-orange-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.late ?? '—' }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ hrAttStats.total_sessions ?? 0 }} {{ t('dashboard.hr_sessions') }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.hr_missing_checkout') }}</span>
          <i class="pi pi-sign-out text-lg text-rose-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.missing_checkout ?? '—' }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ hrOvertimeHours }} {{ t('dashboard.hr_overtime_hours').toLowerCase() }}</div>
      </div>
    </div>

    <!-- Donut charts: status sesi, penggunaan cuti & tren lembur -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-3">
      <!-- Distribusi status sesi (kolom 1, baris 1) -->
      <div v-if="hrStatusTotal > 0" class="md:col-start-1 md:row-start-1 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.hr_session_status') }}</h2>
        <div class="flex flex-col items-center gap-3">
          <svg viewBox="0 0 120 120" class="w-32 h-32 shrink-0">
            <circle
              v-for="seg in hrStatusSegments"
              :key="seg.label"
              cx="60" cy="60" r="45" fill="none" stroke-width="18"
              :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
              transform="rotate(-90 60 60)"
            />
            <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ hrStatusTotal }}</text>
            <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.hr_sessions') }}</text>
          </svg>
          <div class="space-y-2 w-full">
            <div v-for="seg in hrStatusSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
              <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
              <span class="font-semibold text-navy-800 dark:text-gray-100">{{ seg.value }}</span>
              <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Penggunaan cuti per jenis (kolom 1, baris 2 — di bawah status sesi) -->
      <div v-if="hrUsageTotal > 0" class="md:col-start-1 md:row-start-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.hr_leave_usage') }}</h2>
        <div class="flex flex-col items-center gap-3">
          <svg viewBox="0 0 120 120" class="w-32 h-32 shrink-0">
            <circle
              v-for="seg in hrUsageSegments"
              :key="seg.label"
              cx="60" cy="60" r="45" fill="none" stroke-width="18"
              :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
              transform="rotate(-90 60 60)"
            />
            <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ hrUsageTotal }}</text>
            <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.hr_days') }}</text>
          </svg>
          <div class="space-y-2 w-full">
            <div v-for="seg in hrUsageSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
              <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
              <span class="font-semibold text-navy-800 dark:text-gray-100">{{ formatDays(seg.value) }}</span>
              <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Tren lembur per minggu (kolom 2–4, mengisi sisa lebar) -->
      <div v-if="hrTrendBars.length" class="md:col-start-2 md:col-span-3 md:row-start-1 md:row-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-chart-bar text-sm text-teal-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.hr_overtime_trend') }}</h2>
          </div>
          <span class="text-[11px] text-gray-400 dark:text-gray-500">{{ t('dashboard.hr_trend_hours') }}</span>
        </div>
        <div class="flex items-end gap-6 flex-wrap">
          <!-- Bar chart SVG (jam disetujui per minggu) — melebar penuh mengisi sisa kartu -->
          <svg viewBox="0 0 100 56" class="w-full flex-1 min-w-[220px] h-40 shrink-0">
            <!-- baseline -->
            <line x1="2" y1="44" x2="98" y2="44" stroke="currentColor" stroke-opacity="0.15" stroke-width="1" />
            <rect
              v-for="b in hrTrendBars"
              :key="b.week_start"
              :x="b.x"
              :y="b.y"
              width="10"
              :height="b.height"
              rx="1.5"
              class="fill-teal-500 dark:fill-teal-400"
            >
              <title>{{ b.label }} · {{ b.hours }} {{ t('dashboard.hr_trend_hours').toLowerCase() }}</title>
            </rect>
            <text
              v-for="b in hrTrendBars"
              :key="b.week_start + '-l'"
              :x="b.x + 5"
              y="52"
              text-anchor="middle"
              class="fill-gray-400"
              style="font-size:6px"
            >{{ b.label }}</text>
          </svg>
          <!-- Legend: detail per minggu -->
          <div class="space-y-2 flex-1 min-w-0">
            <div v-for="b in hrTrendBars" :key="b.week_start + '-d'" class="flex items-center gap-2 text-sm">
              <span class="w-3 h-3 rounded-sm bg-teal-500 dark:bg-teal-400 shrink-0"></span>
              <span class="text-gray-600 dark:text-gray-300 w-14 shrink-0">{{ b.label }}</span>
              <span class="text-gray-400 dark:text-gray-500 flex-1 truncate">{{ b.count }} {{ t('dashboard.hr_requests') }} · {{ b.approved }} {{ t('dashboard.hr_approved').toLowerCase() }}</span>
              <span class="font-semibold text-navy-800 dark:text-gray-100 shrink-0">{{ b.hours }} jam</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Lembur & Perjalanan Dinas (bulan berjalan) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-3">
      <!-- Lembur -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-clock text-sm text-teal-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.hr_overtime_title') }}</h2>
          </div>
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/attendance/overtime')" />
        </div>
        <div class="grid grid-cols-4 gap-2">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.hr_total') }}</p>
            <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.overtime_total ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-amber-300 dark:border-amber-700/50 bg-amber-50/50 dark:bg-amber-900/10 p-2.5">
            <p class="text-[11px] text-amber-600 dark:text-amber-400">{{ t('dashboard.hr_pending') }}</p>
            <p class="text-lg font-bold text-amber-700 dark:text-amber-300">{{ hrAttStats.overtime_pending ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-emerald-300 dark:border-emerald-700/50 bg-emerald-50/50 dark:bg-emerald-900/10 p-2.5">
            <p class="text-[11px] text-emerald-600 dark:text-emerald-400">{{ t('dashboard.hr_approved') }}</p>
            <p class="text-lg font-bold text-emerald-700 dark:text-emerald-300">{{ hrAttStats.overtime_approved ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-2.5">
            <p class="text-[11px] text-sky-600 dark:text-sky-400">{{ t('dashboard.hr_overtime_hours_approved') }}</p>
            <p class="text-lg font-bold text-sky-700 dark:text-sky-300">{{ hrApprovedOvertimeHours }}</p>
          </div>
        </div>
      </div>

      <!-- Perjalanan Dinas -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-briefcase text-sm text-purple-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.hr_travel_title') }}</h2>
          </div>
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/attendance/business-travel')" />
        </div>
        <div class="grid grid-cols-4 gap-2">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.hr_total') }}</p>
            <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.travel_total ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-emerald-300 dark:border-emerald-700/50 bg-emerald-50/50 dark:bg-emerald-900/10 p-2.5">
            <p class="text-[11px] text-emerald-600 dark:text-emerald-400">{{ t('dashboard.hr_approved') }}</p>
            <p class="text-lg font-bold text-emerald-700 dark:text-emerald-300">{{ hrAttStats.travel_approved ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-sky-300 dark:border-sky-700/50 bg-sky-50/50 dark:bg-sky-900/10 p-2.5">
            <p class="text-[11px] text-sky-600 dark:text-sky-400">{{ t('dashboard.hr_travel_in_progress') }}</p>
            <p class="text-lg font-bold text-sky-700 dark:text-sky-300">{{ hrAttStats.travel_in_progress ?? 0 }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] text-gray-500 dark:text-gray-400">{{ t('dashboard.hr_travel_completed') }}</p>
            <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ hrAttStats.travel_completed ?? 0 }}</p>
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
import { formatDays, monthRange, shortWeekLabel, buildDonutSegments } from '@/utils/dashboard'

import Button from 'primevue/button'

const { t } = useI18n()
const toast = useToast()

const hrLoading = ref(false)
const hrAttStats = ref({ total_sessions: 0, present: 0, late: 0, missing_checkin: 0, missing_checkout: 0, absent: 0, leave_days: 0, total_work_minutes: 0, total_overtime_minutes: 0, overtime_total: 0, overtime_pending: 0, overtime_approved: 0, overtime_minutes: 0, travel_total: 0, travel_approved: 0, travel_in_progress: 0, travel_completed: 0 })
const hrOnLeaveToday = ref(null)
const hrUsageRequests = ref([])
const hrLeaveTypes = ref([])
const hrOvertimeTrend = ref([])
const HR_STATUS_COLORS = ['#10b981', '#f59e0b', '#ef4444', '#f97316', '#6b7280', '#8b5cf6', '#06b6d4']
const HR_LEAVE_COLORS = ['#3b82f6', '#ec4899', '#10b981', '#f59e0b', '#8b5cf6', '#06b6d4', '#ef4444', '#84cc16', '#f97316', '#6366f1']

// Donut distribusi status sesi (agregat seluruh karyawan, bulan berjalan).
const hrStatusSegments = computed(() => buildDonutSegments([
  { label: t('dashboard.hr_present'), value: hrAttStats.value.present, color: HR_STATUS_COLORS[0] },
  { label: t('dashboard.hr_late'), value: hrAttStats.value.late, color: HR_STATUS_COLORS[1] },
  { label: t('dashboard.hr_missing_checkin'), value: hrAttStats.value.missing_checkin, color: HR_STATUS_COLORS[2] },
  { label: t('dashboard.hr_missing_checkout'), value: hrAttStats.value.missing_checkout, color: HR_STATUS_COLORS[3] },
  { label: t('dashboard.hr_absent'), value: hrAttStats.value.absent, color: HR_STATUS_COLORS[4] }
]))
const hrStatusTotal = computed(() => hrStatusSegments.value.reduce((s, i) => s + i.value, 0))

// Donut penggunaan cuti per jenis (agregasi client-side dari usage report).
const hrUsageSegments = computed(() => {
  const byType = {}
  for (const r of hrUsageRequests.value) {
    const id = r.leave_type_id
    if (!byType[id]) byType[id] = { days: 0, count: 0 }
    byType[id].days += Number(r.requested_days) || 0
    byType[id].count++
  }
  const items = Object.entries(byType).map(([id, v]) => {
    const lt = hrLeaveTypes.value.find(x => x.id === id)
    return { label: lt?.name || id, value: v.days, count: v.count }
  })
  return buildDonutSegments(items.map((i, idx) => ({ ...i, color: HR_LEAVE_COLORS[idx % HR_LEAVE_COLORS.length] })))
})
const hrUsageTotal = computed(() => hrUsageSegments.value.reduce((s, i) => s + i.value, 0))
const hrTotalWorkHours = computed(() => ((Number(hrAttStats.value.total_work_minutes) || 0) / 60).toFixed(1))
const hrOvertimeHours = computed(() => ((Number(hrAttStats.value.total_overtime_minutes) || 0) / 60).toFixed(1))
const hrApprovedOvertimeHours = computed(() => ((Number(hrAttStats.value.overtime_minutes) || 0) / 60).toFixed(1))

// Bar chart tren lembur per minggu (jam disetujui).
const hrTrendBars = computed(() => {
  const weeks = hrOvertimeTrend.value || []
  if (!weeks.length) return []
  const max = Math.max(1, ...weeks.map(w => (Number(w.minutes) || 0) / 60))
  const W = 100
  const H = 52
  const barW = 10
  const gap = (W - weeks.length * barW) / (weeks.length + 1)
  return weeks.map((w, i) => {
    const hours = (Number(w.minutes) || 0) / 60
    const h = Math.max(0, (hours / max) * (H - 14))
    return {
      ...w,
      label: shortWeekLabel(w.week_start),
      hours,
      x: gap + i * (barW + gap),
      y: H - 8 - h,
      height: h
    }
  })
})

async function loadHRDashboard() {
  if (hrLoading.value) return
  hrLoading.value = true
  try {
    const { from, to } = monthRange()
    const [attRes, leaveRes, usageRes, typesRes, trendRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/attendance/stats/summary', { params: { from, to } }),
      api.get('/api/v1/tenant/leave/reports/on-leave-today'),
      api.get('/api/v1/tenant/leave/reports/usage', { params: { from, to } }),
      api.get('/api/v1/tenant/leave/types', { params: { page: 1, per_page: 100 } }),
      api.get('/api/v1/tenant/attendance/stats/overtime-trend', { params: { from, to } })
    ])
    if (attRes.status === 'fulfilled') hrAttStats.value = attRes.value.data?.data || hrAttStats.value
    if (leaveRes.status === 'fulfilled') hrOnLeaveToday.value = leaveRes.value.data?.data?.count ?? null
    if (usageRes.status === 'fulfilled') hrUsageRequests.value = usageRes.value.data?.data || []
    if (typesRes.status === 'fulfilled') hrLeaveTypes.value = typesRes.value.data?.data || []
    if (trendRes.status === 'fulfilled') hrOvertimeTrend.value = trendRes.value.data?.data?.weeks || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    hrLoading.value = false
  }
}

onMounted(loadHRDashboard)
</script>
