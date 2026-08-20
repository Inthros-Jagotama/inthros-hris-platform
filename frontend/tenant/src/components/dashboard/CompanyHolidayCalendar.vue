<template>
  <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
    <div class="flex items-center justify-between gap-2 mb-3">
      <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">
        <i class="pi pi-calendar mr-1.5 text-rose-500"></i>
        {{ t('dashboard.company_holidays') }}
      </h2>
      <span class="text-xs text-gray-400 dark:text-gray-500">
        {{ totalHolidays }} {{ t('common.items') }}
      </span>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="space-y-2">
      <div v-for="i in 3" :key="i" class="flex items-center gap-2 animate-pulse">
        <div class="w-5 h-5 rounded bg-gray-200 dark:bg-gray-700"></div>
        <div class="h-3 w-20 bg-gray-200 dark:bg-gray-700 rounded"></div>
        <div class="h-3 flex-1 bg-gray-200 dark:bg-gray-700 rounded"></div>
      </div>
    </div>

    <!-- Empty state -->
    <div v-else-if="holidays.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-6 text-center">
      <i class="pi pi-calendar-plus text-2xl mb-2 opacity-50 block"></i>
      <p>{{ t('company_holidays.calendar_empty_month') }}</p>
    </div>

    <!-- Calendar + Holiday list layout -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Inline DatePicker (2/3 width on large screens) -->
      <div class="lg:col-span-2 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
        <DatePicker
          v-model="calendarDate"
          inline
          :locale="primeLocale"
          class="w-full !border-0"
          @month-change="onMonthChange"
          @year-change="onMonthChange"
        >
          <template #date="{ date }">
            <div class="relative w-full h-full flex flex-col items-center justify-center">
              <span
                class="flex items-center justify-center w-7 h-7 rounded-full transition-colors"
                :class="holidayMap[toDateKey(date)] ? 'bg-rose-500 text-white font-semibold' : ''"
              >{{ date.day }}</span>
            </div>
          </template>
        </DatePicker>
      </div>

      <!-- Holiday list for current month (1/3 width on large screens) -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('company_holidays.calendar_legend') }}</h3>
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ monthHolidays.length }} {{ t('common.items') }}</span>
        </div>
        <div v-if="monthHolidays.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-6 text-center">
          <i class="pi pi-calendar-plus text-2xl mb-2 opacity-50"></i>
          <p>{{ t('company_holidays.calendar_empty_month') }}</p>
        </div>
        <ul v-else class="space-y-1.5">
          <li
            v-for="h in monthHolidays"
            :key="h.id"
            class="flex items-center gap-2 text-sm rounded-md px-2 py-1.5"
            :class="h.isPast ? 'opacity-50' : 'hover:bg-gray-50 dark:hover:bg-gray-700/50'"
          >
            <i class="pi pi-calendar text-xs text-rose-500 shrink-0"></i>
            <span class="text-gray-500 dark:text-gray-400 text-xs shrink-0 w-20">
              {{ formatDate(h.holiday_date, locale) }}
            </span>
            <span class="text-gray-700 dark:text-gray-200 font-medium truncate">{{ h.name }}</span>
            <span v-if="!h.is_active" class="text-xs text-gray-400 shrink-0">({{ t('common_status.inactive') }})</span>
          </li>
        </ul>
      </div>
    </div>

    <!-- Link to full settings page -->
    <div v-if="!loading && holidays.length > 0" class="mt-3 pt-2 border-t border-gray-100 dark:border-gray-700/50">
      <button
        type="button"
        class="text-xs text-teal-600 dark:text-teal-400 hover:underline font-medium"
        @click="$router.push('/settings/company-holidays')"
      >
        {{ t('dashboard.company_holidays_view_all') }} →
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DatePicker from 'primevue/datepicker'
import { formatDate } from '@/utils/formatDate'
import { localDateStr } from '@/utils/dashboard'
import { getPrimeLocale } from '@/utils/primevueLocale'

const { t, locale } = useI18n()

const loading = ref(false)
const holidays = ref([])

const totalHolidays = computed(() => holidays.value.length)
const primeLocale = computed(() => getPrimeLocale(locale.value))

// ── Calendar date model & view tracking ──
const calendarDate = ref(new Date())
const viewYear = ref(calendarDate.value.getFullYear())
const viewMonth = ref(calendarDate.value.getMonth()) // 0-indexed

// Normalize date from API → YYYY-MM-DD
function normDateKey(s) {
  if (!s) return ''
  return String(s).slice(0, 10)
}

// Map YYYY-MM-DD → holiday for highlight on DatePicker
const holidayMap = computed(() => {
  const map = {}
  for (const h of holidays.value) {
    const key = normDateKey(h?.holiday_date)
    if (key) map[key] = h
  }
  return map
})

// Convert Date or PrimeVue slot { day, month, year } → YYYY-MM-DD key
function toDateKey(d) {
  if (!d) return ''
  let y, m0, day
  if (typeof d.getFullYear === 'function') {
    y = d.getFullYear(); m0 = d.getMonth(); day = d.getDate()
  } else {
    y = d.year; m0 = d.month; day = d.day
  }
  if (y == null || m0 == null || day == null || Number.isNaN(y) || Number.isNaN(m0) || Number.isNaN(day)) return ''
  const m = String(m0 + 1).padStart(2, '0')
  return `${y}-${m}-${String(day).padStart(2, '0')}`
}

// Holidays for the currently viewed month in the calendar
const monthHolidays = computed(() => {
  const y = viewYear.value
  const m = viewMonth.value
  const todayKey = localDateStr(new Date())
  return holidays.value
    .filter(h => {
      const parts = normDateKey(h?.holiday_date).split('-')
      return parts.length === 3 && parseInt(parts[0], 10) === y && parseInt(parts[1], 10) === m + 1
    })
    .map(h => ({
      ...h,
      isPast: normDateKey(h.holiday_date) < todayKey
    }))
    .sort((a, b) => (normDateKey(a.holiday_date) < normDateKey(b.holiday_date) ? -1 : 1))
})

// Navigate month/year via PrimeVue arrows & dropdowns
function onMonthChange(e) {
  viewYear.value = e.year
  viewMonth.value = e.month - 1 // PrimeVue month is 1-indexed
}

async function loadHolidays() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/company-holidays', {
      params: { page: 1, per_page: 500 }
    })
    holidays.value = res.data?.data || []
  } catch (e) {
    console.error('Failed to load company holidays:', e)
  } finally {
    loading.value = false
  }
}

onMounted(loadHolidays)
</script>
