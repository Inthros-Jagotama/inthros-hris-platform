<template>
  <div class="space-y-4">
    <div v-if="loading" class="flex items-center justify-center h-40">
      <i class="pi pi-spinner pi-spin text-2xl text-emerald-500"></i>
    </div>

    <template v-else-if="!employeeId">
      <Message severity="warn" :closable="false">{{ t('attendance.no_employee_linked') }}</Message>
    </template>

    <template v-else>
      <!-- Menu cards (pola sama dengan menu Settings) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <button
          v-for="menu in menuCards"
          :key="menu.route"
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-indigo-300 dark:hover:border-indigo-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/50"
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
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-indigo-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>
      </div>

      <!-- Check-in / Check-out -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('attendance.today') }}</h2>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ todayLabel }}</p>
            <Tag v-if="displayStatus" :value="t('attendance.status_' + displayStatus.toLowerCase())" :severity="statusSeverity(displayStatus)" class="!text-xs !px-1.5 !py-0.5 mt-2" />
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

      <!-- Summary -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div v-for="card in summaryCards" :key="card.label" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <div class="flex items-center justify-between mb-1">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ card.label }}</span>
            <i :class="[card.icon, card.iconColor]" class="text-base"></i>
          </div>
          <div class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ card.value }}</div>
        </div>
      </div>

      <!-- Calendar (this month) -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('attendance.calendar_this_month') }}</h2>
        <div v-if="calendarSessions.length === 0" class="text-center py-8 text-gray-400 dark:text-gray-500">
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('attendance.calendar_empty') }}</p>
        </div>
        <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
          <div v-for="session in calendarSessions" :key="session.id" class="flex items-center justify-between py-2 text-sm">
            <span class="text-gray-600 dark:text-gray-300">{{ formatDate(session.work_date, locale) }}</span>
            <div class="flex items-center gap-3">
              <span v-if="session.work_minutes" class="text-gray-400 dark:text-gray-500 text-xs">{{ formatMinutes(session.work_minutes) }}</span>
              <Tag :value="t('attendance.status_' + session.status.toLowerCase())" :severity="statusSeverity(session.status)" class="!text-xs !px-1.5 !py-0.5" />
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { useMyEmployee } from '@/composables/useMyEmployee'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'
import { toLocalISOString } from '@/utils/localTime'
import { formatDate } from '@/utils/formatDate'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'

const router = useRouter()
const { t, locale } = useI18n()
const { hasPermission } = useAuth()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const loading = ref(true)
const summary = ref(null)
const calendarSessions = ref([])
const punching = ref(false)
const punchError = ref('')

// Menu absensi sebagai card (Lembur / Koreksi / Admin — setting), pola sama
// dengan card di halaman Settings (ikon dalam kotak tinted + chevron).
const menuCards = computed(() => {
  const cards = [
    { labelKey: 'attendance.overtime', descKey: 'attendance.overtime_description', icon: 'pi pi-clock', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400', route: '/attendance/overtime', permission: 'attendance.overtime.view' },
    { labelKey: 'attendance.corrections', descKey: 'attendance.corrections_description', icon: 'pi pi-pencil', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/attendance/corrections', permission: 'attendance.corrections.view' },
    { labelKey: 'business_travel.title', descKey: 'business_travel.description', icon: 'pi pi-briefcase', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400', route: '/attendance/business-travel', permission: 'attendance.business-travel.view' }
  ]
  if (hasPermission('attendance.admin.view')) {
    cards.push({ labelKey: 'attendance.admin', descKey: 'attendance.admin_description', icon: 'pi pi-cog', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/attendance/admin' })
  }
  return cards.filter(card => !card.permission || hasPermission(card.permission))
})

const today = new Date()
const todayStr = toDateOnly(today)
const todayLabel = today.toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

const todaySession = computed(() => calendarSessions.value.find(s => s.work_date === todayStr) || null)
const companySettings = ref(null)
const allowDayOffCheckin = computed(() => companySettings.value?.allow_checkin_on_day_off ?? true)

// Event hari ini (dari endpoint events, urut event_time_utc DESC — item pertama = terbaru).
// Dipakai sebagai sumber kebenaran tombol: sesi bisa saja stale (mis. dibuat sebelum
// perbaikan recalculateSession → status DAY_OFF padahal sudah ada event CHECKIN).
const todayEvents = ref([])
function eventDateStr(ev) {
  return (ev.event_time_local || '').slice(0, 10)
}
const todayEventList = computed(() => todayEvents.value.filter(e => eventDateStr(e) === todayStr))

// Status tampilan: sesi bisa stale (mis. DAY_OFF dari sebelum perbaikan
// recalculateSession). Jika ada event hari ini, cerminkan event aktual:
//   CHECKIN terakhir  → MISSING_CHECKOUT (belum check-out)
//   CHECKOUT terakhir → CLOSED (sudah check-in & check-out)
const displayStatus = computed(() => {
  const base = todaySession.value?.status || null
  const list = todayEventList.value
  if ((base === 'DAY_OFF' || !base) && list.length > 0) {
    return list[0].event_type === 'CHECKIN' ? 'MISSING_CHECKOUT' : 'CLOSED'
  }
  return base
})

const canCheckOut = computed(() => displayStatus.value === 'MISSING_CHECKOUT')
const canPunch = computed(() => {
  if (displayStatus.value === 'MISSING_CHECKOUT') return true
  if (displayStatus.value === 'OPEN') return true
  // Hari libur: check-in hanya boleh jika setting allow_checkin_on_day_off aktif.
  if (displayStatus.value === 'DAY_OFF' && allowDayOffCheckin.value) return true
  // Belum ada sesi & belum ada event hari ini → boleh check-in.
  if (!displayStatus.value) return true
  return false
})

function toDateOnly(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function monthRange(date) {
  const from = new Date(date.getFullYear(), date.getMonth(), 1)
  const to = new Date(date.getFullYear(), date.getMonth() + 1, 0)
  return { from: toDateOnly(from), to: toDateOnly(to) }
}

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

function formatMinutes(minutes) {
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return h > 0 ? `${h}h ${m}m` : `${m}m`
}

const summaryCards = computed(() => {
  const s = summary.value
  if (!s) return []
  return [
    { label: t('attendance.summary_present'), value: s.present_days, icon: 'pi pi-check-circle', iconColor: 'text-emerald-500' },
    { label: t('attendance.summary_late'), value: s.late_days, icon: 'pi pi-clock', iconColor: 'text-amber-500' },
    { label: t('attendance.summary_missing_checkout'), value: s.missing_checkout_days, icon: 'pi pi-exclamation-triangle', iconColor: 'text-rose-500' },
    { label: t('attendance.summary_leave'), value: s.leave_days, icon: 'pi pi-calendar', iconColor: 'text-blue-500' }
  ]
})

async function loadSettings() {
  try {
    const res = await api.get('/api/v1/tenant/attendance/settings')
    companySettings.value = res.data?.data || null
  } catch {
    companySettings.value = null
  }
}

async function loadTodayEvents() {
  if (!employeeId.value) return
  try {
    const res = await api.get('/api/v1/tenant/attendance/events', { params: { employee_id: employeeId.value, per_page: 100 } })
    todayEvents.value = res.data?.data || []
  } catch {
    todayEvents.value = []
  }
}

async function loadSummaryAndCalendar() {
  if (!employeeId.value) return
  const { from, to } = monthRange(today)
  const params = { employee_id: employeeId.value, from, to }
  const [summaryRes, calendarRes] = await Promise.all([
    api.get('/api/v1/tenant/attendance/summary', { params }),
    api.get('/api/v1/tenant/attendance/calendar', { params })
  ])
  summary.value = summaryRes.data?.data || null
  calendarSessions.value = calendarRes.data?.data || []
}

async function loadAll() {
  loading.value = true
  try {
    employeeId.value = await loadMyEmployeeId()
    await loadSettings()
    await Promise.all([loadTodayEvents(), loadSummaryAndCalendar()])
  } catch (e) {
    punchError.value = getErrorMessage(e, t('message.failed_to_load'))
  } finally {
    loading.value = false
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
    await loadSummaryAndCalendar()
    await loadTodayEvents()
  } catch (e) {
    // Self-heal: muat ulang event & sesi agar tombol otomatis membalik
    // (mis. sudah check-in dari perangkat lain / sesi stale → tombol Check Out).
    try {
      await Promise.all([loadTodayEvents(), loadSummaryAndCalendar()])
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

onMounted(loadAll)
</script>
