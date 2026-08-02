<template>
  <div class="max-w-full mx-auto">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="flex gap-4">
        <div class="w-56 space-y-2">
          <div v-for="n in 10" :key="n" class="h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
        <div class="flex-1 space-y-3">
          <div v-for="n in 6" :key="n" class="h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- Two-column layout -->
    <div v-else class="flex gap-6">
      <!-- Left: Type Navigation (setiap tipe = sub-form) -->
      <div class="w-56 shrink-0 space-y-1">
        <!-- Organization Info (hanya tampil bila dibuka dari daftar organisasi) -->
        <div v-if="orgId" class="px-3 py-2">
          <div class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">{{ t('organization.title') }}</div>
          <div class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ orgName }}</div>
          <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 font-mono">{{ orgCode }}</div>
        </div>

        <!-- Divider -->
        <div v-if="orgId" class="border-t border-gray-200 dark:border-gray-700 my-2"></div>

        <!-- Section Title -->
        <div class="px-3 py-2">
          <div class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">{{ t('job_values.title') }}</div>
        </div>

        <!-- Type list (sub-form navigation) -->
        <div
          v-for="type in typeOptions"
          :key="type.value"
          role="button"
          tabindex="0"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer"
          :class="selectedType === type.value
            ? 'bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700'
            : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
          @click="selectType(type.value)"
          @keydown.enter="selectType(type.value)"
        >
          <div
            class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0"
            :class="selectedType === type.value ? 'bg-emerald-600 text-white' : 'bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300'"
          >
            <i :class="type.icon" class="text-xs"></i>
          </div>
          <div class="flex-1 min-w-0">
            <div
              class="text-sm font-medium truncate"
              :class="selectedType === type.value ? 'text-emerald-700 dark:text-emerald-300' : 'text-gray-700 dark:text-gray-300'"
            >
              {{ typeLabel(type.value) }}
            </div>
            <div class="text-[10px] text-gray-400">{{ typeCounts[type.value] || 0 }} {{ t('common.items') }}</div>
          </div>
        </div>
      </div>

      <!-- Right: Section content untuk tipe terpilih (komponen generik per tipe) -->
      <div class="flex-1 min-w-0">
        <JobValueSection :key="selectedType" :type="selectedType" @saved="loadTypeCounts" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import JobValueSection from './jobvalues/JobValueSection.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const orgId = route.query.org_id || ''
const orgName = ref('')
const orgCode = ref('')
const pageLoading = ref(true)

// ── Tipe job value — setiap tipe menjadi sub-form di navigasi kiri ──
const typeOptions = [
  { value: 'education', icon: 'pi pi-graduation-cap' },
  { value: 'experience', icon: 'pi pi-briefcase' },
  { value: 'environment', icon: 'pi pi-globe' },
  { value: 'hazard', icon: 'pi pi-exclamation-triangle' },
  { value: 'relationship', icon: 'pi pi-users' },
  { value: 'frequency', icon: 'pi pi-clock' },
  { value: 'asset', icon: 'pi pi-box' },
  { value: 'authority', icon: 'pi pi-shield' },
  { value: 'cash', icon: 'pi pi-dollar' },
  { value: 'impact', icon: 'pi pi-bolt' }
]

// Tipe terpilih — dari query param ?type= agar state bertahan saat refresh
const selectedType = ref(route.query.type && typeOptions.some(x => x.value === route.query.type) ? route.query.type : 'education')
const typeCounts = ref({})

// Fallback label lokal bila key locale belum ada (t() bisa return key string saat missing)
const TYPE_LABEL_FALLBACK = {
  education: 'Education',
  experience: 'Experience',
  environment: 'Environment',
  hazard: 'Hazard',
  relationship: 'Relationship',
  frequency: 'Frequency',
  asset: 'Asset',
  authority: 'Authority',
  cash: 'Cash',
  impact: 'Impact'
}

function typeLabel(value) {
  const label = t(`job_values.types.${value}`)
  // Jika t() mengembalikan key mentah (bukan string terjemahan), pakai fallback lokal
  if (label && !label.startsWith('job_values.types.')) return label
  return TYPE_LABEL_FALLBACK[value] || value
}

// ── Fetch Organization Info ──
async function loadOrgInfo() {
  if (!orgId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${orgId}`)
    const data = res.data?.data
    if (data) {
      orgName.value = data.nomenclature
      orgCode.value = `${data.full_code}`
    }
  } catch {
    orgName.value = 'Unknown'
    orgCode.value = ''
  }
}

// ── Load semua value untuk hitung jumlah per tipe (badge di navigasi) ──
async function loadTypeCounts() {
  try {
    const res = await api.get('/api/v1/tenant/job-management/values', { params: { page: 1, per_page: 200 } })
    const rows = res.data?.data || []
    const counts = {}
    for (const r of rows) {
      if (r.type) counts[r.type] = (counts[r.type] || 0) + 1
    }
    typeCounts.value = counts
  } catch {
    typeCounts.value = {}
  }
}

function selectType(type) {
  if (type === selectedType.value) return
  selectedType.value = type
  // Sync ke URL agar sub-form bertahan saat refresh
  router.replace({ path: route.path, query: { ...route.query, type } })
}

onMounted(async () => {
  pageLoading.value = true
  try {
    await loadOrgInfo()
    await loadTypeCounts()
  } finally {
    pageLoading.value = false
  }
})
</script>
