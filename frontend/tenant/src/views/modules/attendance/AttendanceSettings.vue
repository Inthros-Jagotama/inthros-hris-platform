<template>
  <div class="w-full space-y-4">
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 7" :key="n" class="h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
    </div>

    <template v-else>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 items-start">
      <!-- Kolom 1: Geofence & Lokasi -->
      <!-- Card: Geofence & Lokasi -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-4 flex items-center gap-2">
          <i class="pi pi-map-marker text-emerald-500"></i>
          {{ t('attendance.settings_geofence') }}
        </h3>

        <div class="grid grid-cols-2 gap-3">
          <div class="flex flex-col gap-1 min-w-0">
            <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.latitude') }}</span>
            <InputNumber v-model="form.latitude" class="!w-full num-input" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
          </div>
          <div class="flex flex-col gap-1 min-w-0">
            <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.longitude') }}</span>
            <InputNumber v-model="form.longitude" class="!w-full num-input" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
          </div>
          <div class="col-span-2 flex flex-col gap-1 min-w-0">
            <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.max_distance_meter') }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.max_distance_meter_description') }}</span>
            <InputNumber v-model="form.max_distance_meter" class="!w-48 num-input mt-1" :min="0" size="small" suffix=" m" />
          </div>
        </div>

        <!-- Pencarian lokasi/alamat (Nominatim) -->
        <div class="relative mt-4">
          <div class="relative">
            <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-gray-400 dark:text-gray-500 pointer-events-none"></i>
            <InputText
              v-model="searchQuery"
              :placeholder="t('attendance.search_place_placeholder')"
              size="small"
              class="w-full !pl-9"
              @input="onSearchInput"
              @focus="searchOpen = true"
              @keydown.esc="searchOpen = false"
              @keydown.down.prevent="moveSearchSelection(1)"
              @keydown.up.prevent="moveSearchSelection(-1)"
              @keydown.enter.prevent="selectHighlighted"
            />
          </div>
          <div
            v-if="searchOpen && searchQuery.trim().length >= 3"
            class="absolute z-20 left-0 right-0 mt-1 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-lg max-h-52 overflow-y-auto"
          >
            <div v-if="searching" class="flex items-center gap-2 px-3 py-2.5 text-xs text-gray-400">
              <i class="pi pi-spin pi-spinner"></i>
              <span>{{ t('attendance.searching') }}</span>
            </div>
            <div v-else-if="searchResults.length === 0" class="px-3 py-2.5 text-xs text-gray-400">
              {{ t('attendance.search_no_results') }}
            </div>
            <button
              v-for="(r, idx) in searchResults"
              :key="r.place_id"
              type="button"
              class="flex w-full items-center gap-2 px-3 py-2 text-left text-xs text-gray-700 dark:text-gray-200 hover:bg-blue-50 dark:hover:bg-gray-700 transition-colors"
              :class="idx === searchHighlight ? 'bg-blue-50 dark:bg-gray-700' : ''"
              @mousedown.prevent
              @click="selectSearchResult(r)"
            >
              <i class="pi pi-map-marker text-blue-500 shrink-0"></i>
              <span class="line-clamp-2">{{ r.display_name }}</span>
            </button>
          </div>
        </div>

        <!-- Peta: klik/drag marker untuk memilih titik geofence -->
        <div class="relative z-0 mt-2 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div ref="mapEl" class="h-72 w-full bg-gray-100 dark:bg-gray-700" @click.self="mapFocus"></div>
          <div class="px-2 py-1 text-[11px] text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between gap-2">
            <span>{{ t('attendance.map_click_hint') }}</span>
            <span v-if="form.latitude !== null && form.longitude !== null" class="font-mono whitespace-nowrap">
              {{ formatCoord(form.latitude) }}, {{ formatCoord(form.longitude) }}
            </span>
          </div>
        </div>
      </div>

      <!-- Kolom 2: card lainnya -->
      <div class="space-y-4">
      <!-- Card: Toleransi & Lembur -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-4 flex items-center gap-2">
          <i class="pi pi-clock text-emerald-500"></i>
          {{ t('attendance.settings_tolerance_overtime') }}
        </h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0 flex-1">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.late_tolerance_minutes') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.late_tolerance_minutes_description') }}</span>
            </div>
            <InputNumber v-model="form.late_tolerance_minutes" class="!w-48 shrink-0 num-input" :min="0" size="small" suffix=" min" />
          </div>
          <div class="flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0 flex-1">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.overtime_min_minutes') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.overtime_min_minutes_description') }}</span>
            </div>
            <InputNumber v-model="form.overtime_min_minutes" class="!w-48 shrink-0 num-input" :min="0" size="small" suffix=" min" />
          </div>
          <div class="flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.is_overtime_enabled') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.is_overtime_enabled_description') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_overtime_enabled" class="shrink-0" />
          </div>
        </div>
      </div>

      <!-- Card: Aturan Absensi -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-4 flex items-center gap-2">
          <i class="pi pi-check-circle text-emerald-500"></i>
          {{ t('attendance.settings_attendance_rules') }}
        </h3>
        <div class="space-y-4">
          <div class="flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.is_location_required') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.is_location_required_description') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_location_required" class="shrink-0" />
          </div>
          <div class="flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.is_face_required') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.is_face_required_description') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_face_required" class="shrink-0" />
          </div>
        </div>
      </div>

      <!-- Card: Hari Libur -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-4 flex items-center gap-2">
          <i class="pi pi-calendar-plus text-emerald-500"></i>
          {{ t('attendance.settings_day_off') }}
        </h3>
        <div class="flex items-center justify-between gap-3">
          <div class="flex flex-col gap-0.5 min-w-0">
            <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('attendance.allow_checkin_on_day_off') }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('attendance.allow_checkin_on_day_off_description') }}</span>
          </div>
          <ToggleSwitch v-model="form.allow_checkin_on_day_off" class="shrink-0" />
        </div>
      </div>

          <div class="flex justify-end">
            <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving || !canUpdate" @click="handleSave" />
          </div>
        </div>
        </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'

import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'

// Fix default Leaflet marker icon paths untuk bundler (vite)
delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerIcon2x,
  shadowUrl: markerShadow
})

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const canUpdate = computed(() => hasPermission('attendance.update'))
const loading = ref(true)
const saving = ref(false)
const form = ref({
  latitude: null,
  longitude: null,
  max_distance_meter: null,
  late_tolerance_minutes: null,
  overtime_min_minutes: null,
  is_location_required: false,
  is_face_required: false,
  is_overtime_enabled: false,
  allow_checkin_on_day_off: true
})

// ── Leaflet map — klik/drag untuk memilih titik tengah geofence ──
const mapEl = ref(null)
let map = null
let marker = null
let radiusCircle = null
const DEFAULT_CENTER = [-6.2088, 106.8456] // Jakarta
const DEFAULT_ZOOM = 13

// ── Nominatim address search ──
const searchQuery = ref('')
const searchResults = ref([])
const searching = ref(false)
const searchOpen = ref(false)
const searchHighlight = ref(-1)
let searchTimer = null
const NOMINATIM_ENDPOINT = 'https://nominatim.openstreetmap.org/search'

function formatCoord(v) {
  if (v === null || v === undefined) return '-'
  return Number(v).toFixed(6)
}

function roundCoord(v) {
  return Math.round(Number(v) * 1e6) / 1e6
}

function mapFocus() {
  // container click (di luar tile) diabaikan — interaksi peta tetap normal
}

function onSearchInput() {
  clearTimeout(searchTimer)
  const q = searchQuery.value.trim()
  searchHighlight.value = -1
  if (q.length < 3) {
    searchResults.value = []
    searchOpen.value = false
    searching.value = false
    return
  }
  searchOpen.value = true
  searching.value = true
  searchTimer = setTimeout(async () => {
    const q = searchQuery.value.trim()
    try {
      // fetch langsung (bukan axios `api`) agar token/auth tidak ikut terkirim
      const res = await fetch(
        `${NOMINATIM_ENDPOINT}?format=jsonv2&limit=6&addressdetails=1&q=${encodeURIComponent(q)}`
      )
      if (!res.ok) throw new Error(res.statusText)
      const data = await res.json()
      if (q !== searchQuery.value.trim()) return
      searchResults.value = data
    } catch {
      if (q === searchQuery.value.trim()) searchResults.value = []
    } finally {
      if (q === searchQuery.value.trim()) searching.value = false
    }
  }, 500)
}

function moveSearchSelection(dir) {
  const n = searchResults.value.length
  if (!n) return
  searchHighlight.value = (searchHighlight.value + dir + n) % n
}

function selectHighlighted() {
  if (searchHighlight.value >= 0 && searchResults.value[searchHighlight.value]) {
    selectSearchResult(searchResults.value[searchHighlight.value])
  }
}

function selectSearchResult(r) {
  const lat = Number(r.lat)
  const lng = Number(r.lon)
  form.value.latitude = roundCoord(lat)
  form.value.longitude = roundCoord(lng)
  searchQuery.value = ''
  searchResults.value = []
  searchOpen.value = false
  searchHighlight.value = -1
  if (map) map.setView([lat, lng], Math.max(map.getZoom(), 16))
  syncMarker()
}

function initMap() {
  if (!mapEl.value) return
  if (map) {
    // Pengaman: bila container map lama sudah tidak terhubung ke DOM, buat ulang
    if (!map.getContainer().isConnected) {
      destroyMap()
    } else {
      map.invalidateSize()
      syncMarker()
      return
    }
  }

  map = L.map(mapEl.value, { attributionControl: false }).setView(DEFAULT_CENTER, DEFAULT_ZOOM)

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>'
  }).addTo(map)

  marker = L.marker(DEFAULT_CENTER, { draggable: true }).addTo(map)

  radiusCircle = L.circle(DEFAULT_CENTER, {
    radius: 100,
    color: '#3b82f6',
    weight: 1.5,
    fillColor: '#3b82f6',
    fillOpacity: 0.12
  }).addTo(map)

  // Klik peta → set lat/long
  map.on('click', (e) => {
    const { lat, lng } = e.latlng
    form.value.latitude = roundCoord(lat)
    form.value.longitude = roundCoord(lng)
    syncMarker()
  })

  // Drag marker → set lat/long
  marker.on('dragend', (e) => {
    const { lat, lng } = e.target.getLatLng()
    form.value.latitude = roundCoord(lat)
    form.value.longitude = roundCoord(lng)
  })

  syncMarker()
}

function destroyMap() {
  if (map) {
    map.remove()
    map = null
    marker = null
    radiusCircle = null
  }
}

function syncMarker() {
  if (!map || !marker) return
  const lat = form.value.latitude
  const lng = form.value.longitude
  if (lat !== null && lat !== undefined && lng !== null && lng !== undefined) {
    const pos = L.latLng(Number(lat), Number(lng))
    marker.setLatLng(pos)
    if (radiusCircle) radiusCircle.setLatLng(pos)
    if (!map.getBounds().contains(pos)) {
      map.setView(pos, Math.max(map.getZoom(), DEFAULT_ZOOM))
    }
  }
  syncRadius()
}

function syncRadius() {
  if (!radiusCircle) return
  const lat = form.value.latitude
  const lng = form.value.longitude
  const r = Number(form.value.max_distance_meter)
  const hasPoint = lat !== null && lat !== undefined && lng !== null && lng !== undefined
  if (hasPoint && r > 0) {
    radiusCircle.setLatLng(L.latLng(Number(lat), Number(lng)))
    radiusCircle.setRadius(r)
    radiusCircle.setStyle({ opacity: 1, fillOpacity: 0.12 })
  } else {
    radiusCircle.setStyle({ opacity: 0, fillOpacity: 0 })
  }
}

// Sinkronkan marker & lingkaran radius saat lat/long/jarak diketik manual
watch(
  () => [form.value.latitude, form.value.longitude, form.value.max_distance_meter],
  () => syncMarker()
)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/attendance/settings')
    const s = res.data?.data
    if (s) {
      form.value = {
        latitude: s.latitude ?? null,
        longitude: s.longitude ?? null,
        max_distance_meter: s.max_distance_meter ?? null,
        late_tolerance_minutes: s.late_tolerance_minutes ?? null,
        overtime_min_minutes: s.overtime_min_minutes ?? null,
        is_location_required: !!s.is_location_required,
        is_face_required: !!s.is_face_required,
        is_overtime_enabled: !!s.is_overtime_enabled,
        allow_checkin_on_day_off: !!s.allow_checkin_on_day_off
      }
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
    // Inisialisasi peta setelah card Geofence ter-render (v-else)
    nextTick(initMap)
  }
}

async function handleSave() {
  saving.value = true
  try {
    await api.put('/api/v1/tenant/attendance/settings', form.value)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadData)

onBeforeUnmount(() => {
  clearTimeout(searchTimer)
  destroyMap()
})
</script>

<style scoped>
/* Input dalam mengisi penuh root yang berlebar tetap, sehingga semua kotak input
   berukuran sama dan tepi kanannya rata (default .p-inputnumber-input tidak berwidth). */
.num-input :deep(.p-inputnumber-input) {
  width: 100%;
}
</style>

