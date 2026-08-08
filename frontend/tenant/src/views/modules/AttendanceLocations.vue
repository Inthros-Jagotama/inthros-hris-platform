<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button v-if="hasPermission('attendance.create')" :label="t('common.add')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable
      v-else
      :value="items"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-map-marker text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('attendance.locations_empty') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('common.name')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="latitude" :header="t('attendance.latitude')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.latitude }}</span></template>
      </Column>
      <Column field="longitude" :header="t('attendance.longitude')" style="width:140px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.longitude }}</span></template>
      </Column>
      <Column field="radius_m" :header="t('attendance.radius_m')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.radius_m }} m</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button v-if="hasPermission('attendance.update')" icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button v-if="hasPermission('attendance.delete')" icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('common.edit') : t('common.add')" modal :style="{ width: '720px', maxWidth: '95vw' }" @hide="onDialogHide" @shown="onDialogShown">
      <div class="space-y-3">
        <FormRow :label="t('common.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="255" autofocus />
        </FormRow>
        <div class="relative">
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
        <div class="relative z-0 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div ref="mapEl" class="h-80 w-full bg-gray-100 dark:bg-gray-700" @click.self="mapFocus"></div>
          <div class="px-2 py-1 text-[11px] text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between gap-2">
            <span>{{ t('attendance.map_click_hint') }}</span>
            <span v-if="form.latitude !== null && form.longitude !== null" class="font-mono whitespace-nowrap">
              {{ formatCoord(form.latitude) }}, {{ formatCoord(form.longitude) }}
            </span>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('attendance.latitude')" required :errors="errors?.latitude">
            <InputNumber v-model="form.latitude" class="!w-full" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
          </FormRow>
          <FormRow :label="t('attendance.longitude')" required :errors="errors?.longitude">
            <InputNumber v-model="form.longitude" class="!w-full" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('attendance.radius_m')" :errors="errors?.radius_m">
          <InputNumber v-model="form.radius_m" class="!w-full" :min="0" size="small" suffix=" m" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="deleting"
      :error="deleteError"
      :title="t('attendance.confirm_delete_title')"
      :message="t('attendance.confirm_delete_location', { name: deleteTarget?.name || '' })"
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import L from 'leaflet'
import 'leaflet/dist/leaflet.css'
import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

// Fix default Leaflet marker icon paths for bundlers (webpack/vite)
delete L.Icon.Default.prototype._getIconUrl
L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerIcon2x,
  shadowUrl: markerShadow
})

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

// ── Leaflet map state ──
const mapEl = ref(null)
let map = null
let marker = null
let radiusCircle = null
const DEFAULT_CENTER = [-6.2088, 106.8456] // Jakarta
const DEFAULT_ZOOM = 13

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ name: '', latitude: null, longitude: null, radius_m: null })

// ── Nominatim address search state ──
const searchQuery = ref('')
const searchResults = ref([])
const searching = ref(false)
const searchOpen = ref(false)
const searchHighlight = ref(-1)
let searchTimer = null
const NOMINATIM_ENDPOINT = 'https://nominatim.openstreetmap.org/search'

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const skeletonColumns = [
  { type: 'text', width: 'w-44', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/attendance/locations', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = {
    name: item?.name || '',
    latitude: item?.latitude ?? null,
    longitude: item?.longitude ?? null,
    radius_m: item?.radius_m ?? null
  }
  dialogVisible.value = true
  nextTick(() => initMap())
}

function formatCoord(v) {
  if (v === null || v === undefined) return '-'
  return Number(v).toFixed(6)
}

function mapFocus() {
  // Keep map interactions clickable; container click (outside tiles) is ignored
}

// ── Nominatim address search ──
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
    // Guard anti race condition: buang respons jika query sudah berubah sejak fetch dimulai
    const q = searchQuery.value.trim()
    try {
      // fetch langsung (bukan axios `api`) agar token/auth header tidak ikut terkirim ke pihak ketiga
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

function resetSearch() {
  clearTimeout(searchTimer)
  searchQuery.value = ''
  searchResults.value = []
  searchOpen.value = false
  searchHighlight.value = -1
}

function onDialogShown() {
  // Dipanggil setelah transisi dialog selesai — container peta sudah punya ukuran
  // nyata, jadi invalidateSize() tidak menghitung 0 (perbaiki peta blank saat buka ulang).
  nextTick(() => {
    initMap()
    // Jeda tambahan sebagai jaring pengaman bila transisi/animasi masih berlangsung.
    setTimeout(() => map?.invalidateSize(), 150)
  })
}

function destroyMap() {
  if (map) {
    map.remove()
    map = null
    marker = null
    radiusCircle = null
  }
}

function onDialogHide() {
  resetForm()
  // Dialog PrimeVue me-unmount kontennya saat ditutup — hancurkan map agar
  // instance Leaflet lama tidak menempel ke elemen DOM yang sudah hilang
  // (jika dibiarkan, peta blank saat dialog dibuka lagi).
  destroyMap()
}

function initMap() {
  if (!mapEl.value) return
  if (map) {
    // Pengaman: bila container map lama sudah tidak terhubung ke DOM (mis. dialog
    // di-remount), buang instance lama & buat baru di container yang baru.
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

  // Click on the map → set lat/long
  map.on('click', (e) => {
    const { lat, lng } = e.latlng
    form.value.latitude = roundCoord(lat)
    form.value.longitude = roundCoord(lng)
    syncMarker()
  })

  // Drag the marker → set lat/long
  marker.on('dragend', (e) => {
    const { lat, lng } = e.target.getLatLng()
    form.value.latitude = roundCoord(lat)
    form.value.longitude = roundCoord(lng)
  })

  syncMarker()
}

function roundCoord(v) {
  return Math.round(Number(v) * 1e6) / 1e6
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
  const r = Number(form.value.radius_m)
  const hasPoint = lat !== null && lat !== undefined && lng !== null && lng !== undefined
  if (hasPoint && r > 0) {
    radiusCircle.setLatLng(L.latLng(Number(lat), Number(lng)))
    radiusCircle.setRadius(r)
    radiusCircle.setStyle({ opacity: 1, fillOpacity: 0.12 })
  } else {
    radiusCircle.setStyle({ opacity: 0, fillOpacity: 0 })
  }
}

// Keep marker + radius circle in sync when lat/long/radius typed manually
watch(
  () => [form.value.latitude, form.value.longitude, form.value.radius_m],
  () => {
    if (dialogVisible.value) syncMarker()
  }
)

function resetForm() {
  form.value = { name: '', latitude: null, longitude: null, radius_m: null }
  errors.value = {}
  editing.value = false
  editingId.value = null
  resetSearch()
}

async function handleSave() {
  errors.value = {}
  if (!form.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  if (form.value.latitude === null || form.value.latitude === undefined) { errors.value = { latitude: t('form.required') }; return }
  if (form.value.longitude === null || form.value.longitude === undefined) { errors.value = { longitude: t('form.required') }; return }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (editing.value) {
      await api.put(`/api/v1/tenant/attendance/locations/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/attendance/locations', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/attendance/locations/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)

onBeforeUnmount(() => {
  clearTimeout(searchTimer)
  destroyMap()
})
</script>
