<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-end">
      <div class="flex items-center gap-2 text-sm text-gray-400">
        <Button
          icon="pi pi-refresh"
          size="small"
          severity="secondary"
          text
          :loading="loading"
          @click="loadData"
          v-tooltip.left="'Refresh data'"
        />
        <span v-if="lastUpdated" class="text-gray-400">
          Updated: {{ lastUpdated }}
        </span>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading && !loaded" class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
      <div v-for="i in 6" :key="i" class="bg-white rounded-lg border border-gray-200 p-3 animate-pulse">
        <div class="w-8 h-8 rounded-lg bg-gray-200 mb-2"></div>
        <div class="h-5 w-16 bg-gray-200 rounded mb-1"></div>
        <div class="h-3 w-20 bg-gray-200 rounded"></div>
      </div>
    </div>

    <!-- KPI Cards -->
    <div v-else class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
      <div
        v-for="kpi in kpis"
        :key="kpi.label"
        class="bg-white rounded-lg border border-gray-200 p-3 hover:shadow-sm transition-shadow cursor-pointer group"
      >
        <div class="flex items-center gap-2 mb-1">
          <div class="w-8 h-8 rounded-lg flex items-center justify-center transition-transform group-hover:scale-110" :class="kpi.bg">
            <i :class="kpi.icon" class="text-sm" :style="{ color: kpi.color }"></i>
          </div>
        </div>
        <p class="text-lg font-bold text-gray-800">{{ kpi.value }}</p>
        <p class="text-sm text-gray-500">{{ kpi.label }}</p>
      </div>
    </div>

    <!-- Two-column layout -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Recent Companies -->
      <div class="bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">Recent Companies</h3>
          <router-link to="/companies" class="text-sm text-indigo-600 hover:underline">View all</router-link>
        </div>
        <div class="p-2">
          <div v-if="recentCompanies.length === 0 && !loading" class="text-sm text-gray-400 text-center py-4">
            No companies registered yet.
          </div>
          <div
            v-for="company in recentCompanies"
            :key="company.id"
            class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-gray-50 text-sm transition-colors"
          >
            <div class="flex items-center gap-2 min-w-0">
              <i class="pi pi-building text-gray-400 text-sm shrink-0"></i>
              <span class="text-gray-700 truncate">{{ company.name }}</span>
            </div>
            <Tag :value="company.status" :severity="statusSeverity(company.status)" class="!text-xs !px-1.5 !py-0.5 shrink-0" />
          </div>
          <div v-if="loading" class="space-y-1 p-2">
            <div v-for="i in 3" :key="i" class="h-8 bg-gray-100 rounded-md animate-pulse"></div>
          </div>
        </div>
      </div>

      <!-- System Health -->
      <div class="bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">System Health</h3>
          <router-link to="/monitoring" class="text-sm text-indigo-600 hover:underline">Details</router-link>
        </div>
        <div class="p-3 space-y-2">
          <div v-if="loading" class="space-y-2">
            <div v-for="i in 4" :key="i" class="h-6 bg-gray-100 rounded animate-pulse"></div>
          </div>
          <div v-else>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Platform Status</span>
              <Tag :value="healthStatus" :severity="healthSeverity" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Platform Database</span>
              <span class="text-sm" :class="platformDbHealthy ? 'text-emerald-600' : 'text-rose-600'">
                <i :class="platformDbHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-circle'" class="mr-1"></i>
                {{ platformDbHealthy ? 'Connected' : 'Unhealthy' }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Redis Cache</span>
              <span class="text-sm" :class="cacheHealthy ? 'text-emerald-600' : 'text-amber-600'">
                <i :class="cacheHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle'" class="mr-1"></i>
                {{ cacheHealthy ? 'Connected' : cacheStatus }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Active Tenants</span>
              <span class="text-sm text-gray-600 font-medium">{{ activeTenantCount }} companies</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Pool Connections</span>
              <span class="text-sm text-gray-600 font-medium">{{ poolStatsText }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">Total Users</span>
              <span class="text-sm text-gray-600 font-medium">{{ totalUsersText }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="flex items-center gap-3">
      <router-link to="/companies" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-plus text-sm"></i> New Company
      </router-link>
      <router-link to="/users" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-user-plus text-sm"></i> Add User
      </router-link>
      <router-link to="/monitoring" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-chart-bar text-sm"></i> View Health
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'

const toast = useToast()

// State
const loading = ref(false)
const loaded = ref(false)
const lastUpdated = ref('')
const companies = ref([])
const usersTotal = ref(0)
const modulesTotal = ref(0)
const health = ref(null)

// Helper aman untuk memastikan array tidak null/undefined
const safeCompanies = computed(() => Array.isArray(companies.value) ? companies.value : [])

// Derived KPIs dengan perlindungan Array.isArray
const activeCompanyCount = computed(() =>
  safeCompanies.value.filter(c => c?.status === 'active').length
)

const activeTenantCount = computed(() =>
  activeCompanyCount.value
)

const healthStatus = computed(() =>
  health.value?.status || 'Unknown'
)

const healthSeverity = computed(() => {
  if (health.value?.status === 'healthy') return 'success'
  if (health.value?.status === 'degraded') return 'warn'
  return 'info'
})

const platformDbHealthy = computed(() => {
  const db = health.value?.database
  if (!db) return false
  return Object.values(db).every(v => v === 'connected')
})

const cacheHealthy = computed(() =>
  health.value?.cache === 'connected'
)

const cacheStatus = computed(() =>
  health.value?.cache || 'Unknown'
)

const poolStatsText = computed(() => {
  const ps = health.value?.pool_stats
  if (!ps) return '-'
  return `${ps.total_open ?? 0} open / ${ps.total_idle ?? 0} idle`
})

const totalUsersText = computed(() =>
  `${usersTotal.value ?? 0} platform admins`
)

const poolStatsTotalOpen = computed(() => {
  const ps = health.value?.pool_stats
  return ps?.total_open ?? 0
})

const healthPercent = computed(() => {
  if (!health.value) return '-'
  const db = health.value?.database || {}
  const entries = Object.keys(db).filter(k => !k.startsWith('tenant:'))
  if (entries.length === 0) return '0%'
  const healthyCount = entries.filter(k => db[k] === 'connected').length
  return `${Math.round((healthyCount / entries.length) * 100)}%`
})

// KPIs for cards (sekarang fully safe dari null/undefined)
const kpis = computed(() => [
  {
    label: 'Total Companies',
    value: safeCompanies.value.length.toString(),
    icon: 'pi pi-building',
    bg: 'bg-indigo-50',
    color: '#4f46e5'
  },
  {
    label: 'Active Tenants',
    value: activeCompanyCount.value.toString(),
    icon: 'pi pi-check-circle',
    bg: 'bg-emerald-50',
    color: '#059669'
  },
  {
    label: 'Platform Users',
    value: (usersTotal.value ?? 0).toString(),
    icon: 'pi pi-users',
    bg: 'bg-sky-50',
    color: '#0284c7'
  },
  {
    label: 'Modules',
    value: (modulesTotal.value ?? 0).toString(),
    icon: 'pi pi-cog',
    bg: 'bg-amber-50',
    color: '#d97706'
  },
  {
    label: 'Active Connections',
    value: `${poolStatsTotalOpen.value}`,
    icon: 'pi pi-database',
    bg: 'bg-purple-50',
    color: '#7c3aed'
  },
  {
    label: 'System Health',
    value: healthPercent.value,
    icon: 'pi pi-heart',
    bg: healthPercent.value === '100%' ? 'bg-emerald-50' : 'bg-amber-50',
    color: healthPercent.value === '100%' ? '#059669' : '#d97706'
  }
])

// Recent companies (first 3)
const recentCompanies = computed(() =>
  safeCompanies.value.slice(0, 3)
)

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'suspended': return 'warn'
    case 'terminated': return 'danger'
    default: return 'info'
  }
}

function formatTime(date) {
  const d = new Date(date)
  return d.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

async function loadData() {
  if (loading.value) return
  loading.value = true

  try {
    const [companiesRes, usersRes, modulesRes, healthRes] = await Promise.allSettled([
      api.get('/api/v1/platform/companies?per_page=100'),
      api.get('/api/v1/platform/users?per_page=1'),
      api.get('/api/v1/platform/modules?per_page=100'),
      api.get('/api/v1/platform/monitoring/health')
    ])

    // Companies validation
    if (companiesRes.status === 'fulfilled') {
      const d = companiesRes.value.data
      const rawCompanies = d?.data || d
      companies.value = Array.isArray(rawCompanies) ? rawCompanies : []
    } else {
      companies.value = []
    }

    // Users total validation
    if (usersRes.status === 'fulfilled') {
      const d = usersRes.value.data
      usersTotal.value = d?.total ?? (Array.isArray(d?.data) ? d.data.length : 0)
    } else {
      usersTotal.value = 0
    }

    // Modules total validation
    if (modulesRes.status === 'fulfilled') {
      const d = modulesRes.value.data
      modulesTotal.value = d?.total ?? (Array.isArray(d?.data) ? d.data.length : 0)
    } else {
      modulesTotal.value = 0
    }

    // Health
    if (healthRes.status === 'fulfilled') {
      health.value = healthRes.value.data || null
    } else {
      health.value = null
    }

    loaded.value = true
    lastUpdated.value = formatTime(new Date())
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load dashboard data',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>