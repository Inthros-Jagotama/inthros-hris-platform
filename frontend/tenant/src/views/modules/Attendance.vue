<template>
  <div class="space-y-4">
    <div v-if="loading" class="flex items-center justify-center h-40">
      <i class="pi pi-spinner pi-spin text-2xl text-emerald-500"></i>
    </div>

    <template v-else-if="!employeeId">
      <Message severity="warn" :closable="false">{{ t('attendance.no_employee_linked') }}</Message>
    </template>

    <template v-else>
      <div v-if="hasPermission('attendance.update')" class="flex justify-end">
        <Button :label="t('attendance.admin')" icon="pi pi-cog" size="small" severity="secondary" outlined @click="router.push('/attendance/admin')" />
      </div>

      <!-- Check-in / Check-out -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between flex-wrap gap-3">
          <div>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('attendance.today') }}</h2>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ todayLabel }}</p>
            <Tag v-if="todaySession" :value="t('attendance.status_' + todaySession.status.toLowerCase())" :severity="statusSeverity(todaySession.status)" class="!text-xs !px-1.5 !py-0.5 mt-2" />
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
            <span class="text-gray-600 dark:text-gray-300">{{ session.work_date }}</span>
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
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Message from 'primevue/message'

const router = useRouter()
const { t } = useI18n()
const { hasPermission } = useAuth()

const loading = ref(true)
const employeeId = ref('')
const summary = ref(null)
const calendarSessions = ref([])
const punching = ref(false)
const punchError = ref('')

const today = new Date()
const todayStr = toDateOnly(today)
const todayLabel = today.toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })

const todaySession = computed(() => calendarSessions.value.find(s => s.work_date === todayStr) || null)
// OPEN artinya sudah check-in tapi belum check-out; status lain (CLOSED/DAY_OFF/LEAVE/dst.) berarti tidak ada aksi lanjutan hari ini.
const canCheckOut = computed(() => todaySession.value?.status === 'OPEN')
const canPunch = computed(() => !todaySession.value || canCheckOut.value)

function toDateOnly(date) {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

// toLocalISOString membangun ISO string dengan offset zona waktu lokal browser
// (mis. 2026-08-08T14:30:00+07:00), karena backend attendance/events butuh
// event_time_local dengan offset eksplisit, bukan hanya UTC.
function toLocalISOString(date) {
  const pad = (n) => String(n).padStart(2, '0')
  const offsetMin = -date.getTimezoneOffset()
  const sign = offsetMin >= 0 ? '+' : '-'
  const offH = pad(Math.floor(Math.abs(offsetMin) / 60))
  const offM = pad(Math.abs(offsetMin) % 60)
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}T${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}${sign}${offH}:${offM}`
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

async function loadEmployeeId() {
  const res = await api.get('/api/v1/tenant/user-accounts/me')
  employeeId.value = res.data?.data?.employee_id || ''
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
    await loadEmployeeId()
    await loadSummaryAndCalendar()
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
      location_provider: 'browser'
    })
    await loadSummaryAndCalendar()
  } catch (e) {
    punchError.value = e?.code === 1
      ? t('attendance.geolocation_denied')
      : getErrorMessage(e, t('attendance.punch_failed'))
  } finally {
    punching.value = false
  }
}

onMounted(loadAll)
</script>
