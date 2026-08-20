<template>
  <div class="flex flex-col items-end shrink-0 text-right" data-testid="live-clock">
    <span class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ dateLabel }}</span>
    <span class="text-xs text-gray-500 dark:text-gray-400 tabular-nums">{{ timeLabel }}</span>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useLanguage } from '@/stores/language'
import api from '@/services/api'

const { state } = useLanguage()

const timezone = ref('Asia/Jakarta')
const now = ref(new Date())
let intervalId = null

const locale = computed(() => (state.lang === 'id' ? 'id-ID' : 'en-US'))

const dateLabel = computed(() =>
  new Intl.DateTimeFormat(locale.value, {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
    year: 'numeric',
    timeZone: timezone.value,
  }).format(now.value)
)

const timeLabel = computed(() =>
  new Intl.DateTimeFormat(locale.value, {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
    timeZone: timezone.value,
  }).format(now.value)
)

async function loadTimezone() {
  try {
    const res = await api.get('/api/v1/tenant/attendance/timezone/me')
    const tz = res.data?.data?.timezone || res.data?.timezone
    if (tz) timezone.value = tz
  } catch {
    // Elemen dekoratif — biarkan fallback Asia/Jakarta, tidak perlu toast error.
  }
}

onMounted(() => {
  loadTimezone()
  intervalId = setInterval(() => {
    now.value = new Date()
  }, 1000)
})

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId)
})
</script>
