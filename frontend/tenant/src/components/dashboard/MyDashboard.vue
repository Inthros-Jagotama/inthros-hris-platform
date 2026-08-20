<template>
  <div class="space-y-4">
    <!-- Sapaan ke user yang login -->
    <div class="flex items-start gap-3">
      <div class="w-11 h-11 rounded-full bg-emerald-100 dark:bg-emerald-900/50 flex items-center justify-center shrink-0">
        <i class="pi pi-user text-emerald-600 dark:text-emerald-400 text-lg"></i>
      </div>
      <div class="min-w-0">
        <h1 class="text-lg font-semibold text-navy-800 dark:text-gray-100 truncate">{{ greetingText }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ t('dashboard.greeting_sub') }}</p>
      </div>
    </div>

    <!-- Check-in / Check-out -->
    <div v-if="myLoading" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 animate-pulse">
      <div class="h-4 w-40 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
      <div class="h-3 w-64 bg-gray-200 dark:bg-gray-700 rounded"></div>
    </div>
    <div v-else class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center justify-between flex-wrap gap-3">
        <div>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('attendance.today') }}</h2>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ todayLabel }}</p>
          <Tag v-if="checkinStatus" :value="t('attendance.status_' + checkinStatus.toLowerCase())" :severity="statusSeverity(checkinStatus)" class="!text-xs !px-1.5 !py-0.5 mt-2" />
        </div>
        <Button
          :label="canCheckOut ? t('attendance.check_out') : t('attendance.check_in')"
          :icon="canCheckOut ? 'pi pi-sign-out' : 'pi pi-sign-in'"
          :severity="canCheckOut ? 'warn' : 'success'"
          :loading="punching"
          :disabled="!canPunch"
          @click="handlePunch"
        />
      </div>
      <Message v-if="punchError" severity="error" :closable="true" class="mt-3" @close="punchError = ''">{{ punchError }}</Message>
    </div>

    <!-- KPI ringkasan personal -->
    <div v-if="myLoading" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-3">
      <div v-for="i in 4" :key="i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
        <div class="h-4 w-28 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
        <div class="h-8 w-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    </div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-3">
      <!-- Kehadiran hari ini -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.my_attendance') }}</span>
          <i class="pi pi-clock text-lg text-emerald-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">
          {{ attendanceStatusLabel }}
        </div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">
          {{ t('dashboard.my_check_in') }}: {{ myCheckInTime || '—' }} · {{ t('dashboard.my_check_out') }}: {{ myCheckOutTime || '—' }}
        </div>
      </div>
      <!-- Saldo cuti -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.my_leave_balance') }}</span>
          <i class="pi pi-calendar text-lg text-amber-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ formatDays(totalRemainingDays) }} <span class="text-sm font-medium text-gray-400">hari</span></div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ balances.length }} {{ t('dashboard.leave_types') }}</div>
      </div>
      <!-- Persetujuan menunggu -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.my_approvals') }}</span>
          <i class="pi pi-check-square text-lg text-violet-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ approvalTasks.length }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ t('dashboard.my_approvals') }}</div>
      </div>
      <!-- Progres KPI -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.my_kpi') }}</span>
          <i class="pi pi-chart-line text-lg text-cyan-500"></i>
        </div>
        <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ myKpiProgress }}%</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ t('dashboard.overall_progress') }}</div>
      </div>
    </div>

    <!-- Detail lists -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Saldo cuti per jenis -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.my_leave_balance') }}</h2>
        <div v-if="balances.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-4 text-center">
          {{ t('dashboard.my_leave_balance_empty') }}
        </div>
        <div v-else class="space-y-2">
          <div v-for="b in balances.slice(0, 5)" :key="b.id" class="flex items-center justify-between gap-2 text-sm">
            <span class="text-gray-600 dark:text-gray-300 truncate">{{ leaveTypeName(b.leave_type_id) }}</span>
            <span class="font-semibold text-emerald-600 dark:text-emerald-400 shrink-0">{{ formatDays(b.remaining_days) }}</span>
          </div>
        </div>
      </div>
      <!-- Persetujuan menunggu saya -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.my_approvals') }}</h2>
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/approvals')" />
        </div>
        <div v-if="approvalTasks.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-4 text-center">
          {{ t('dashboard.my_approvals_empty') }}
        </div>
        <div v-else class="space-y-2">
          <div v-for="task in approvalTasks.slice(0, 5)" :key="task.id" class="flex items-start gap-2 text-sm">
            <i class="pi pi-circle text-[8px] text-violet-400 mt-1.5 shrink-0"></i>
            <div class="min-w-0">
              <p class="text-gray-700 dark:text-gray-200 truncate">{{ task.flow_name || t('dashboard.task') }}</p>
              <p class="text-[11px] text-gray-400 dark:text-gray-500">{{ task.status }}</p>
            </div>
          </div>
        </div>
      </div>
      <!-- Progres KPI saya -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <div class="flex items-center justify-between gap-2 mb-3">
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.my_kpi') }}</h2>
          <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/performance')" />
        </div>
        <div v-if="kpiItems.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-4 text-center">
          {{ t('dashboard.my_kpi_empty') }}
        </div>
        <div v-else class="space-y-2">
          <div v-for="k in kpiItems.slice(0, 5)" :key="k.detail_id" class="text-sm">
            <div class="flex items-center justify-between gap-2">
              <span class="text-gray-600 dark:text-gray-300 truncate">{{ k.indicator_name }}</span>
              <span class="font-semibold text-cyan-600 dark:text-cyan-400 shrink-0">{{ Number(k.achievement || 0).toFixed(0) }}%</span>
            </div>
            <div class="h-1.5 bg-gray-100 dark:bg-gray-700 rounded-full mt-1">
              <div class="h-1.5 bg-cyan-500 rounded-full" :style="{ width: Math.min(100, Number(k.achievement || 0)) + '%' }"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Company Holiday Calendar + Quick Access (side by side) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Left: Company Holiday Calendar -->
      <CompanyHolidayCalendar />

      <!-- Right: Quick Access (self-service) -->
      <div v-if="quickAccessCards.length" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.quick_access') }}</h2>
        <div class="space-y-2.5">
          <button
            v-for="card in quickAccessCards"
            :key="card.route"
            type="button"
            class="cursor-pointer group flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-teal-300 dark:hover:border-teal-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-teal-500/50 w-full"
            @click="$router.push(card.route)"
          >
            <div class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors" :class="card.tint">
              <i :class="card.icon" class="text-base"></i>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-navy-800 dark:text-gray-100 truncate">{{ t(card.labelKey) }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(card.descKey) }}</p>
            </div>
            <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-teal-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
          </button>
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
import { useAuth } from '@/stores/auth'
import { useMyEmployee } from '@/composables/useMyEmployee'
import api from '@/services/api'
import { toLocalISOString } from '@/utils/localTime'
import { localDateStr, formatTime, formatDays } from '@/utils/dashboard'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'
import CompanyHolidayCalendar from './CompanyHolidayCalendar.vue'

const { t } = useI18n()
const toast = useToast()
const auth = useAuth()
const { hasPermission, hasExactPermission } = auth
const { employeeId, loadMyEmployeeId } = useMyEmployee()

// ── Sapaan berdasarkan waktu (pagi/siang/sore/malam) + nama user login ──
const greetingText = computed(() => {
  const h = new Date().getHours()
  let key = 'dashboard.greeting_night'
  if (h >= 5 && h < 11) key = 'dashboard.greeting_morning'
  else if (h >= 11 && h < 15) key = 'dashboard.greeting_afternoon'
  else if (h >= 15 && h < 18) key = 'dashboard.greeting_evening'
  const name = auth.state.user?.name || ''
  const greeting = t(key)
  return name ? `${greeting}, ${name}!` : `${greeting}!`
})

const myLoading = ref(false)
const myAttendanceEvents = ref([])
const balances = ref([])
const leaveTypes = ref([])
const approvalTasks = ref([])
const kpiData = ref(null)

// ── Check-in / Check-out ──
const punching = ref(false)
const punchError = ref('')
const todaySession = ref(null)
const companySettings = ref(null)
const todayLabel = new Date().toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

// Quick Access — menu self-service, gaya card standar halaman modul lain.
// Permission dicocokkan dengan route tujuan (mirror):
//   - Overtime/Corrections/Business Travel → attendance.operations.view (strict)
//   - Leave → leave.view (route /leave tidak punya permission — pakai module-level)
//   - My Reimbursement → reimbursement.requests.view (route tidak strict)
const quickAccessCards = computed(() => {
  const cards = [
    { labelKey: 'attendance.overtime', descKey: 'attendance.overtime_description', icon: 'pi pi-clock', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400', route: '/attendance/overtime', permission: 'attendance.operations.view', exact: true },
    { labelKey: 'attendance.corrections', descKey: 'attendance.corrections_description', icon: 'pi pi-pencil', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/attendance/corrections', permission: 'attendance.operations.view', exact: true },
    { labelKey: 'business_travel.title', descKey: 'business_travel.description', icon: 'pi pi-briefcase', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400', route: '/attendance/business-travel', permission: 'attendance.operations.view', exact: true },
    { labelKey: 'leave.title', descKey: 'leave.description', icon: 'pi pi-calendar', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/leave', permission: 'leave.view', exact: false },
    { labelKey: 'reimbursement.my_requests', descKey: 'reimbursement.description', icon: 'pi pi-credit-card', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', route: '/reimbursements/my-requests', permission: 'reimbursement.requests.view', exact: false }
  ]
  return cards.filter(c => {
    if (!c.permission) return true
    return c.exact ? hasExactPermission(c.permission) : hasPermission(c.permission)
  })
})

const todayEvents = computed(() => {
  const todayStr = localDateStr(new Date())
  return myAttendanceEvents.value.filter(e => (e.event_time_local || '').slice(0, 10) === todayStr)
})

// Status check-in: sesi bisa stale (mis. DAY_OFF dari sebelum perbaikan
// recalculateSession). Jika ada event hari ini, cerminkan event aktual:
//   CHECKIN terakhir  → MISSING_CHECKOUT (belum check-out)
//   CHECKOUT terakhir → CLOSED (sudah check-in & check-out)
const allowDayOffCheckin = computed(() => companySettings.value?.allow_checkin_on_day_off ?? true)
const checkinStatus = computed(() => {
  const base = todaySession.value?.status || null
  const list = todayEvents.value
  if ((base === 'DAY_OFF' || !base) && list.length > 0) {
    return list[0].event_type === 'CHECKIN' ? 'MISSING_CHECKOUT' : 'CLOSED'
  }
  return base
})
const canCheckOut = computed(() => checkinStatus.value === 'MISSING_CHECKOUT')
const canPunch = computed(() => {
  if (checkinStatus.value === 'MISSING_CHECKOUT') return true
  if (checkinStatus.value === 'OPEN') return true
  // Hari libur: check-in hanya boleh jika setting allow_checkin_on_day_off aktif.
  if (checkinStatus.value === 'DAY_OFF' && allowDayOffCheckin.value) return true
  // Belum ada sesi & belum ada event hari ini → boleh check-in.
  if (!checkinStatus.value) return true
  return false
})

function statusSeverity(status) {
  switch (status) {
    case 'CLOSED': return 'success'
    case 'OPEN': return 'info'
    case 'MISSING_CHECKIN':
    case 'MISSING_CHECKOUT': return 'warn'
    case 'LEAVE': return 'help'
    case 'DAY_OFF': return 'secondary'
    default: return 'secondary'
  }
}

function getCurrentPosition() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error(t('attendance.geolocation_unsupported')))
      return
    }
    navigator.geolocation.getCurrentPosition(resolve, reject, { enableHighAccuracy: true, timeout: 10000 })
  })
}

async function handlePunch() {
  punchError.value = ''
  punching.value = true
  try {
    const position = await getCurrentPosition()
    const now = new Date()
    await api.post('/api/v1/tenant/attendance/events', {
      employee_id: employeeId.value,
      event_type: canCheckOut.value ? 'CHECKOUT' : 'CHECKIN',
      event_time_utc: now.toISOString(),
      event_time_local: toLocalISOString(now),
      latitude: position.coords.latitude,
      longitude: position.coords.longitude,
      accuracy_m: position.coords.accuracy ? Math.round(position.coords.accuracy) : null,
      location_provider: 'GPS'
    })
    await loadCheckinData()
  } catch (e) {
    // Self-heal: muat ulang event & sesi agar tombol otomatis membalik
    // (mis. sudah check-in dari perangkat lain / sesi stale → tombol Check Out).
    try {
      await loadCheckinData()
    } catch {
      // Abaikan kegagalan reload — error asli tetap ditampilkan.
    }
    punchError.value = e?.code === 1
      ? t('attendance.geolocation_denied')
      : getErrorMessage(e, t('attendance.punch_failed'))
  } finally {
    punching.value = false
  }
}

// Muat ulang data check-in (event hari ini + sesi hari ini + setting hari libur).
async function loadCheckinData() {
  if (!employeeId.value) return
  const todayStr = localDateStr(new Date())
  const [eventsRes, calRes, settingsRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/attendance/events', { params: { employee_id: employeeId.value, per_page: 100 } }),
    api.get('/api/v1/tenant/attendance/calendar', { params: { employee_id: employeeId.value, from: todayStr, to: todayStr } }),
    api.get('/api/v1/tenant/attendance/settings')
  ])
  if (eventsRes.status === 'fulfilled') myAttendanceEvents.value = eventsRes.value.data?.data || []
  if (calRes.status === 'fulfilled') {
    todaySession.value = (calRes.value.data?.data || []).find(s => s.work_date === todayStr) || null
  }
  if (settingsRes.status === 'fulfilled') companySettings.value = settingsRes.value.data?.data || null
}
const attendanceStatusLabel = computed(() => {
  const events = todayEvents.value
  if (events.length === 0) return t('dashboard.my_not_checked_in')
  return t('dashboard.my_present')
})
const myCheckInTime = computed(() => {
  const checkin = todayEvents.value.filter(e => e.event_type === 'CHECKIN').pop()
  return checkin ? formatTime(checkin.event_time_local) : ''
})
const myCheckOutTime = computed(() => {
  const checkout = todayEvents.value.filter(e => e.event_type === 'CHECKOUT').pop()
  return checkout ? formatTime(checkout.event_time_local) : ''
})
const totalRemainingDays = computed(() => balances.value.reduce((sum, b) => sum + (Number(b.remaining_days) || 0), 0))
const myKpiProgress = computed(() => {
  const evalSummary = kpiData.value?.evaluation
  if (evalSummary && Number(evalSummary.overall_progress) > 0) return Math.round(Number(evalSummary.overall_progress))
  const items = kpiData.value?.kpi_progress || []
  if (!items.length) return 0
  const avg = items.reduce((s, k) => s + (Number(k.achievement) || 0), 0) / items.length
  return Math.round(avg)
})
const kpiItems = computed(() => kpiData.value?.kpi_progress || [])

function leaveTypeName(id) {
  const found = leaveTypes.value.find(lt => lt.id === id)
  return found?.name || id || ''
}

async function loadMyDashboard() {
  if (myLoading.value) return
  myLoading.value = true
  try {
    const empId = await loadMyEmployeeId()
    if (!empId) return
    const [balRes, typesRes, apprRes, kpiRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/leave/balances', { params: { employee_id: empId, page: 1, per_page: 100 } }),
      api.get('/api/v1/tenant/leave/types', { params: { page: 1, per_page: 100 } }),
      api.get('/api/v1/tenant/approval/tasks/pending', { params: { page: 1, per_page: 5 } }),
      api.get(`/api/v1/tenant/performance/kpi/dashboard/employee/${empId}`)
    ])
    balances.value = balRes.status === 'fulfilled' ? (balRes.value.data?.data || []) : []
    leaveTypes.value = typesRes.status === 'fulfilled' ? (typesRes.value.data?.data || []) : []
    approvalTasks.value = apprRes.status === 'fulfilled' ? (apprRes.value.data?.data || []) : []
    kpiData.value = kpiRes.status === 'fulfilled' ? (kpiRes.value.data?.data || null) : null
    // Check-in: event hari ini + sesi + setting hari libur.
    await loadCheckinData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    myLoading.value = false
  }
}

// Dimuat saat pertama kali view dibuka (shell memakai KeepAlive, jadi
// hanya sekali per sesi).
onMounted(loadMyDashboard)
</script>
