<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-end">
      <div class="flex items-center gap-2 text-sm text-gray-400">
        <span v-if="autoRefreshActive" class="text-xs text-emerald-500 flex items-center gap-1">
          <i class="pi pi-sync text-[10px] animate-spin"></i> {{ t('dashboard.auto_refresh') }}
        </span>
        <Button
          icon="pi pi-refresh"
          size="small"
          severity="secondary"
          text
          :loading="loading"
          @click="loadData"
          v-tooltip.left="t('dashboard.refresh_tooltip')"
        />
        <span v-if="lastUpdated" class="text-gray-400">
          {{ t('dashboard.updated') }}: {{ lastUpdated }}
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

    <!-- Chart Row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Company Trend Chart -->
      <div class="lg:col-span-2 bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">{{ t('dashboard.chart_company_trend') }}</h3>
          <Tag :value="t('dashboard.chart_monthly')" severity="info" class="!text-xs" rounded />
        </div>
        <div class="p-3">
          <div v-if="chartData" class="h-48">
            <Chart type="bar" :data="chartData" :options="chartOptions" />
          </div>
          <div v-else class="h-48 flex items-center justify-center text-sm text-gray-400">
            <i class="pi pi-chart-bar text-2xl mr-2 opacity-50"></i>
            {{ t('common.loading') }}
          </div>
        </div>
      </div>

      <!-- Pool Wait / System Load -->
      <div class="bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">{{ t('dashboard.pool_wait_count') }}</h3>
          <i class="pi pi-database text-gray-400 text-sm"></i>
        </div>
        <div class="p-4 flex flex-col items-center justify-center h-48">
          <div class="text-4xl font-bold text-gray-800 mb-1">{{ totalWaitCount }}</div>
          <p class="text-sm text-gray-500">{{ t('dashboard.pool_wait_desc') }}</p>
          <!-- Mini sparkline: wait distribution -->
          <div class="w-full mt-4 space-y-2">
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>Open</span>
              <span class="font-medium text-gray-700">{{ poolSummary.total_open ?? 0 }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>In Use</span>
              <span class="font-medium text-gray-700">{{ poolSummary.total_in_use ?? 0 }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500">
              <span>Idle</span>
              <span class="font-medium text-gray-700">{{ poolSummary.total_idle ?? 0 }}</span>
            </div>
            <ProgressBar
              :value="poolUtilization"
              :class="poolUtilization > 80 ? '!bg-rose-100' : '!bg-emerald-100'"
              class="!h-1.5 !rounded-full mt-1"
            />
            <div class="text-[10px] text-gray-400 text-center">{{ poolUtilization }}% {{ t('dashboard.pool_wait_desc') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Two-column layout -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Recent Companies -->
      <div class="bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">{{ t('dashboard.recent_companies') }}</h3>
          <router-link to="/companies" class="text-sm text-indigo-600 hover:underline">{{ t('common.view_all') }}</router-link>
        </div>
        <div class="p-2">
          <div v-if="recentCompanies.length === 0 && !loading" class="text-sm text-gray-400 text-center py-4">
            {{ t('dashboard.no_companies') }}
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
          <h3 class="text-sm font-semibold text-gray-700">{{ t('dashboard.system_health') }}</h3>
          <router-link to="/monitoring" class="text-sm text-indigo-600 hover:underline">{{ t('common.details') }}</router-link>
        </div>
        <div class="p-3 space-y-2">
          <div v-if="loading" class="space-y-2">
            <div v-for="i in 4" :key="i" class="h-6 bg-gray-100 rounded animate-pulse"></div>
          </div>
          <div v-else>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.platform_status') }}</span>
              <Tag :value="healthStatus" :severity="healthSeverity" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.platform_db') }}</span>
              <span class="text-sm" :class="platformDbHealthy ? 'text-emerald-600' : 'text-rose-600'">
                <i :class="platformDbHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-circle'" class="mr-1"></i>
                {{ platformDbHealthy ? t('common_status.connected') : t('common_status.unhealthy') }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.redis_cache') }}</span>
              <span class="text-sm" :class="cacheHealthy ? 'text-emerald-600' : 'text-amber-600'">
                <i :class="cacheHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle'" class="mr-1"></i>
                {{ cacheHealthy ? t('common_status.connected') : cacheStatus }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.active_tenants') }}</span>
              <span class="text-sm text-gray-600 font-medium">{{ activeTenantCount }} {{ t('common_status.connected') }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.pool_connections') }}</span>
              <span class="text-sm text-gray-600 font-medium">{{ poolStatsText }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600">{{ t('dashboard.total_users') }}</span>
              <span class="text-sm text-gray-600 font-medium">{{ totalUsersText }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="flex items-center gap-3">
      <router-link to="/companies" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-plus text-sm"></i> {{ t('dashboard.quick_new_company') }}
      </router-link>
      <router-link to="/users" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-user-plus text-sm"></i> {{ t('dashboard.quick_add_user') }}
      </router-link>
      <router-link to="/monitoring" class="flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-200 rounded-md text-sm text-gray-600 hover:border-indigo-200 hover:text-indigo-600 transition-colors">
        <i class="pi pi-chart-bar text-sm"></i> {{ t('dashboard.quick_view_health') }}
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useLanguage } from '@/stores/language'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Chart from 'primevue/chart'
import ProgressBar from 'primevue/progressbar'

const toast = useToast()
const { t } = useI18n()
const langStore = useLanguage()

// ── State ──
const loading = ref(false)
const loaded = ref(false)
const lastUpdated = ref('')
const companies = ref([])
const usersTotal = ref(0)
const modulesTotal = ref(0)
const health = ref(null)
const autoRefreshActive = ref(true)
let refreshTimer = null

// ── Derived ──
const safeCompanies = computed(() => Array.isArray(companies.value) ? companies.value : [])

const activeCompanyCount = computed(() =>
  safeCompanies.value.filter(c => c?.status === 'active').length
)

const activeTenantCount = computed(() => activeCompanyCount.value)

const healthStatus = computed(() => health.value?.status || t('common_status.checking'))

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

const cacheHealthy = computed(() => health.value?.cache === 'connected')

const cacheStatus = computed(() => health.value?.cache || t('common_status.checking'))

const poolStatsText = computed(() => {
  const ps = health.value?.pool_stats
  if (!ps) return '-'
  return `${ps.total_open ?? 0} open / ${ps.total_idle ?? 0} idle`
})

const totalUsersText = computed(() => `${usersTotal.value ?? 0} platform admins`)

const poolStatsTotalOpen = computed(() => health.value?.pool_stats?.total_open ?? 0)

const healthPercent = computed(() => {
  if (!health.value) return '-'
  const db = health.value?.database || {}
  const entries = Object.keys(db).filter(k => !k.startsWith('tenant:'))
  if (entries.length === 0) return '0%'
  const healthyCount = entries.filter(k => db[k] === 'connected').length
  return `${Math.round((healthyCount / entries.length) * 100)}%`
})

// ── Pool Wait & Utilization ──
const totalWaitCount = computed(() => health.value?.pool_stats?.total_wait_count ?? 0)

const poolSummary = computed(() => health.value?.pool_stats || {})

const poolUtilization = computed(() => {
  const open = health.value?.pool_stats?.total_open ?? 0
  const inUse = health.value?.pool_stats?.total_in_use ?? 0
  if (open === 0) return 0
  return Math.round((inUse / open) * 100)
})

// ── KPI Cards ──
const kpis = computed(() => [
  {
    label: t('dashboard.kpi_total_companies'),
    value: safeCompanies.value.length.toString(),
    icon: 'pi pi-building', bg: 'bg-indigo-50', color: '#4f46e5'
  },
  {
    label: t('dashboard.kpi_active_tenants'),
    value: activeCompanyCount.value.toString(),
    icon: 'pi pi-check-circle', bg: 'bg-emerald-50', color: '#059669'
  },
  {
    label: t('dashboard.kpi_platform_users'),
    value: (usersTotal.value ?? 0).toString(),
    icon: 'pi pi-users', bg: 'bg-sky-50', color: '#0284c7'
  },
  {
    label: t('dashboard.kpi_modules'),
    value: (modulesTotal.value ?? 0).toString(),
    icon: 'pi pi-cog', bg: 'bg-amber-50', color: '#d97706'
  },
  {
    label: t('dashboard.kpi_active_connections'),
    value: `${poolStatsTotalOpen.value}`,
    icon: 'pi pi-database', bg: 'bg-purple-50', color: '#7c3aed'
  },
  {
    label: t('dashboard.kpi_system_health'),
    value: healthPercent.value,
    icon: 'pi pi-heart',
    bg: healthPercent.value === '100%' ? 'bg-emerald-50' : 'bg-amber-50',
    color: healthPercent.value === '100%' ? '#059669' : '#d97706'
  }
])

// ── Chart: Company Trend per Month ──
const chartData = computed(() => {
  const c = safeCompanies.value
  if (c.length === 0) return null

  // Group by month (last 6 months)
  const now = new Date()
  const months = []
  for (let i = 5; i >= 0; i--) {
    const d = new Date(now.getFullYear(), now.getMonth() - i, 1)
      const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
      const lang = langStore.state.lang === 'id' ? 'id-ID' : 'en-US'
      months.push({ key, label: d.toLocaleDateString(lang, { month: 'short', year: '2-digit' }), count: 0 })
  }

  c.forEach(company => {
    const created = company.createdAt || company.created_at
    if (!created) return
    const d = new Date(created)
    if (isNaN(d.getTime())) return
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
    const month = months.find(m => m.key === key)
    if (month) month.count++
  })

  return {
    labels: months.map(m => m.label),
    datasets: [{
      label: t('dashboard.chart_new_companies'),
      data: months.map(m => m.count),
      backgroundColor: ['#4f46e5', '#6366f1', '#818cf8', '#a5b4fc', '#c7d2fe', '#e0e7ff'],
      borderRadius: 4,
      borderSkipped: false
    }]
  }
})

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: { legend: { display: false } },
  scales: {
    x: {
      grid: { display: false },
      ticks: { font: { size: 11 } }
    },
    y: {
      beginAtZero: true,
      ticks: {
        font: { size: 11 },
        stepSize: 1
      },
      grid: { color: '#f3f4f6' }
    }
  }
}

// ── Recent Companies ──
const recentCompanies = computed(() => safeCompanies.value.slice(0, 3))

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

// ── Load Data ──
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

    if (companiesRes.status === 'fulfilled') {
      const d = companiesRes.value.data
      const rawCompanies = d?.data || d
      companies.value = Array.isArray(rawCompanies) ? rawCompanies : []
    } else {
      companies.value = []
    }

    if (usersRes.status === 'fulfilled') {
      const d = usersRes.value.data
      usersTotal.value = d?.total ?? (Array.isArray(d?.data) ? d.data.length : 0)
    } else {
      usersTotal.value = 0
    }

    if (modulesRes.status === 'fulfilled') {
      const d = modulesRes.value.data
      modulesTotal.value = d?.total ?? (Array.isArray(d?.data) ? d.data.length : 0)
    } else {
      modulesTotal.value = 0
    }

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
      summary: t('message.error'),
      detail: t('message.failed_to_load'),
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

// ── Auto-refresh polling (30s) ──
function startAutoRefresh() {
  stopAutoRefresh()
  refreshTimer = setInterval(() => {
    if (!autoRefreshActive.value) return
    loadData()
  }, 30000)
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  loadData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>
