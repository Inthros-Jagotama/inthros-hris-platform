<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-between gap-3 flex-wrap">
      <div class="flex items-center gap-3">
        <!-- Auto-refresh Toggle -->
        <div class="flex items-center gap-2 text-sm">
          <ToggleSwitch v-model="autoRefreshActive" size="small" />
          <label class="text-xs text-gray-500 cursor-pointer select-none" @click="autoRefreshActive = !autoRefreshActive">
            <i class="pi pi-sync mr-1" :class="{ 'text-emerald-500': autoRefreshActive, 'text-gray-300': !autoRefreshActive }"></i>
            {{ t('monitoring.auto_refresh') }}
          </label>
        </div>
        <!-- Live Indicator -->
        <span v-if="autoRefreshActive" class="flex items-center gap-1 text-[10px] text-emerald-500 font-medium">
          <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
          {{ t('dashboard.live_indicator') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <!-- Alert Badges -->
        <Tag
          v-if="alerts.length > 0"
          :value="`${alerts.length} Alert${alerts.length > 1 ? 's' : ''}`"
          severity="danger"
          class="!text-[10px] !px-1.5 !py-0"
          v-tooltip.left="alertTooltipText"
        />
        <span v-if="lastUpdated" class="text-xs text-gray-400">{{ t('dashboard.updated') }}: {{ lastUpdated }}</span>
        <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" severity="secondary" :loading="loading" :disabled="loading" @click="loadAll" />
      </div>
    </div>

    <!-- Alert Threshold Cards -->
    <div v-if="alerts.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2">
      <div
        v-for="alert in alerts"
        :key="alert.message"
        class="flex items-center gap-2 px-3 py-2 rounded-lg text-xs"
        :class="alert.severity === 'danger' ? 'bg-rose-50 text-rose-700 border border-rose-200' : 'bg-amber-50 text-amber-700 border border-amber-200'"
      >
        <i :class="alert.severity === 'danger' ? 'pi pi-exclamation-circle text-rose-400' : 'pi pi-exclamation-triangle text-amber-400'"></i>
        <span>{{ alert.message }}</span>
      </div>
    </div>

    <!-- Loading Skeleton -->
    <div v-if="loading && !loaded" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="lg:col-span-2 bg-white rounded-lg border border-gray-200 p-4 animate-pulse">
        <div class="h-4 w-32 bg-gray-200 rounded mb-3"></div>
        <div class="h-32 bg-gray-100 rounded"></div>
      </div>
      <div class="bg-white rounded-lg border border-gray-200 p-4 animate-pulse">
        <div class="h-4 w-24 bg-gray-200 rounded mb-3"></div>
        <div class="space-y-2">
          <div v-for="i in 4" :key="i" class="h-6 bg-gray-100 rounded"></div>
        </div>
      </div>
    </div>

    <!-- Charts & Platform Health -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Pool Utilization Chart -->
      <div class="lg:col-span-2 bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">{{ t('monitoring.pool_utilization_chart') }}</h3>
          <div class="flex items-center gap-3">
            <span class="text-[10px] text-gray-400">{{ poolHistory.length }} samples</span>
            <Button icon="pi pi-trash" size="small" text severity="secondary" class="!text-[10px] !w-5 !h-5" v-tooltip.left="'Clear history'" @click="poolHistory = []" />
          </div>
        </div>
        <div class="p-3">
          <div v-if="poolHistory.length > 1" class="h-52">
            <Chart type="line" :data="poolChartData" :options="poolChartOptions" />
          </div>
          <div v-else class="h-52 flex items-center justify-center text-sm text-gray-400">
            <i class="pi pi-chart-line text-2xl mr-2 opacity-50"></i>
            {{ t('monitoring.collecting_data') }}
          </div>
        </div>
      </div>

      <!-- Platform Health Card -->
      <div class="bg-white rounded-lg border border-gray-200">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100">
          <h3 class="text-sm font-semibold text-gray-700">{{ t('monitoring.platform_health') }}</h3>
          <Tag :value="overallHealthStatus" :severity="overallHealthSeverity" class="!text-xs" />
        </div>
        <div class="p-3 space-y-2">
          <div v-if="!health.database && !health.cache" class="text-sm text-gray-400 py-2">{{ t('monitoring.loading_health') }}</div>
          <template v-else>
            <!-- Database Status -->
            <div v-for="(val, key) in health.database" :key="key" class="flex items-center justify-between text-sm px-2 py-1.5 rounded-md hover:bg-gray-50">
              <span class="text-gray-600 capitalize text-xs">{{ key.replace(/_/g, ' ') }}</span>
              <div class="flex items-center gap-1.5">
                <span class="text-xs" :class="val === 'connected' ? 'text-emerald-600' : 'text-rose-600'">{{ val }}</span>
                <i v-if="val === 'connected'" class="pi pi-check-circle text-emerald-400 text-[10px]"></i>
                <i v-else class="pi pi-exclamation-circle text-rose-400 text-[10px]"></i>
              </div>
            </div>
            <!-- Cache Status -->
            <div class="flex items-center justify-between text-sm px-2 py-1.5 rounded-md hover:bg-gray-50">
              <span class="text-gray-600 text-xs">Redis Cache</span>
              <div class="flex items-center gap-1.5">
                <span class="text-xs" :class="cacheIsHealthy ? 'text-emerald-600' : 'text-amber-600'">{{ health.cache || 'N/A' }}</span>
                <i v-if="cacheIsHealthy" class="pi pi-check-circle text-emerald-400 text-[10px]"></i>
                <i v-else class="pi pi-exclamation-triangle text-amber-400 text-[10px]"></i>
              </div>
            </div>
          </template>
          <!-- Quick Pool Summary -->
          <div class="border-t border-gray-100 pt-2 mt-2">
            <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
              <span>{{ t('monitoring.pool_open') }}</span>
              <span class="font-medium text-gray-700">{{ poolStatValue('total_open') }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
              <span>{{ t('monitoring.pool_idle') }}</span>
              <span class="font-medium text-gray-700">{{ poolStatValue('total_idle') }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
              <span>In Use</span>
              <span class="font-medium text-gray-700">{{ poolStatValue('total_in_use') }}</span>
            </div>
            <ProgressBar :value="utilizationPercent" :class="utilizationPercent > 80 ? '!bg-rose-100' : utilizationPercent > 50 ? '!bg-amber-100' : '!bg-emerald-100'" class="!h-1.5 !rounded-full mt-2" />
            <div class="text-[10px] text-gray-400 text-center mt-1">{{ utilizationPercent }}% {{ t('monitoring.utilization') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tenant Connections -->
    <div class="bg-white rounded-lg border border-gray-200">
      <div class="px-4 py-2.5 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-gray-700">{{ t('monitoring.tenant_connections') }}</h3>
        <div class="flex items-center gap-2">
          <span v-if="unhealthyTenantCount > 0" class="text-xs text-rose-500 font-medium flex items-center gap-1">
            <i class="pi pi-exclamation-circle"></i>
            {{ unhealthyTenantCount }} unhealthy
          </span>
          <span class="text-sm text-gray-500">{{ tenants.length }} {{ t('monitoring.active_connections') }}</span>
        </div>
      </div>
      <div class="p-2">
        <DataTable :value="tenants" size="small" class="!text-sm">
          <template #empty>
            <div class="flex flex-col items-center justify-center py-8 text-gray-400">
              <i class="pi pi-database text-2xl mb-1 opacity-50"></i>
              <p class="text-sm">{{ t('monitoring.no_tenants') }}</p>
            </div>
          </template>
          <Column field="company_name" :header="t('monitoring.company')">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-indigo-400 text-sm"></i>
                <span>{{ data.company_name || data.company_id }}</span>
              </div>
            </template>
          </Column>
          <Column field="status" :header="t('common.status')">
            <template #body="{ data }">
              <Tag :value="data.status || 'connected'" :severity="(data.status || 'connected') === 'healthy' ? 'success' : 'danger'" class="!text-xs" />
            </template>
          </Column>
          <Column field="pool.open" :header="t('monitoring.pool_open')">
            <template #body="{ data }">
              <span :class="poolCellClass(data, 'open')">{{ data.pool?.open ?? '-' }}</span>
            </template>
          </Column>
          <Column field="pool.in_use" header="In Use">
            <template #body="{ data }">
              <div class="flex items-center gap-1">
                <span :class="poolCellClass(data, 'in_use')">{{ data.pool?.in_use ?? '-' }}</span>
                <i v-if="data.pool?.max_open && data.pool?.in_use > data.pool?.max_open * 0.8" class="pi pi-exclamation-triangle text-amber-400 text-[10px]" v-tooltip.left="'High utilization'"></i>
              </div>
            </template>
          </Column>
          <Column field="pool.wait_count" header="Waits">
            <template #body="{ data }">
              <div class="flex items-center gap-1">
                <span :class="data.pool?.wait_count > 0 ? 'text-amber-600 font-medium' : 'text-gray-500'">{{ data.pool?.wait_count ?? '-' }}</span>
                <i v-if="data.pool?.wait_count > 0" class="pi pi-exclamation-triangle text-amber-400 text-[10px]" v-tooltip.left="data.pool?.wait_duration ? `Wait: ${data.pool.wait_duration}` : 'Has waiters'"></i>
              </div>
            </template>
          </Column>
          <Column field="pool.idle" :header="t('monitoring.pool_idle')" />
          <Column field="driver" :header="t('monitoring.driver')" />
          <Column field="last_active" :header="t('monitoring.last_active')">
            <template #body="{ data }">{{ data.last_active || '-' }}</template>
          </Column>
        </DataTable>
      </div>
    </div>

    <!-- Quick Stats Row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <div v-for="stat in quickStats" :key="stat.label" class="bg-white rounded-lg border border-gray-200 p-3" :class="stat.highlight ? 'ring-2 ring-amber-200' : ''">
        <p class="text-sm text-gray-500">{{ stat.label }}</p>
        <p class="text-lg font-bold mt-0.5" :class="stat.color || 'text-gray-800'">{{ stat.value }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import Chart from 'primevue/chart'
import ProgressBar from 'primevue/progressbar'

const { t } = useI18n()

const health = ref({})
const tenants = ref([])
const loading = ref(false)
const loaded = ref(false)
const lastUpdated = ref('')
const autoRefreshActive = ref(true)

// Pool history for chart (rolling buffer of last 20 samples)
const poolHistory = ref([])
const MAX_HISTORY = 20

let refreshTimer = null

const overallHealthStatus = computed(() => {
  if (health.value?.status === 'healthy') return t('common_status.healthy')
  if (health.value?.status === 'degraded') return t('common_status.degraded')
  return t('common_status.checking')
})

const overallHealthSeverity = computed(() => {
  if (health.value?.status === 'healthy') return 'success'
  if (health.value?.status === 'degraded') return 'warn'
  return 'info'
})

const cacheIsHealthy = computed(() => health.value?.cache === 'connected')

// Pool stat helpers
function poolStatValue(key) {
  return health.value?.pool_stats?.[key] ?? 0
}

const utilizationPercent = computed(() => {
  const open = poolStatValue('total_open')
  const inUse = poolStatValue('total_in_use')
  if (open === 0) return 0
  return Math.round((inUse / open) * 100)
})

// ── Alert Thresholds ──
const unhealthyTenantCount = computed(() =>
  tenants.value.filter(t => t.status === 'unhealthy').length
)

const totalWaitCount = computed(() => poolStatValue('total_wait_count'))

const alerts = computed(() => {
  const list = []
  // 1. Connection pressure: wait count > 0
  if (totalWaitCount.value > 0) {
    list.push({
      severity: 'warn',
      message: `Connection pressure: ${totalWaitCount.value} pending ${totalWaitCount.value === 1 ? 'wait' : 'waits'}`
    })
  }
  // 2. High utilization: in_use > 80% of open connections
  if (utilizationPercent.value > 80) {
    list.push({
      severity: 'danger',
      message: `Pool utilization at ${utilizationPercent.value}% — consider scaling`
    })
  } else if (utilizationPercent.value > 50) {
    list.push({
      severity: 'warn',
      message: `Pool utilization at ${utilizationPercent.value}%`
    })
  }
  // 3. Unhealthy tenants
  if (unhealthyTenantCount.value > 0) {
    list.push({
      severity: 'danger',
      message: `${unhealthyTenantCount.value} unhealthy tenant connection${unhealthyTenantCount.value > 1 ? 's' : ''}`
    })
  }
  // 4. Cache unhealthy
  if (health.value?.cache && !cacheIsHealthy.value) {
    list.push({
      severity: 'warn',
      message: `Redis cache is ${health.value.cache}`
    })
  }
  return list
})

const alertTooltipText = computed(() =>
  alerts.value.map(a => `• ${a.message}`).join('\n')
)

// ── Pool Chart ──
const poolChartData = computed(() => {
  const history = poolHistory.value
  if (history.length < 2) return null

  const labels = history.map((_, i) => `#${i + 1}`)

  return {
    labels,
    datasets: [
      {
        label: 'Open',
        data: history.map(h => h.total_open),
        borderColor: '#6366f1',
        backgroundColor: 'rgba(99, 102, 241, 0.1)',
        fill: true,
        tension: 0.3,
        pointRadius: 2,
        borderWidth: 1.5
      },
      {
        label: 'In Use',
        data: history.map(h => h.total_in_use),
        borderColor: '#f59e0b',
        backgroundColor: 'rgba(245, 158, 11, 0.1)',
        fill: true,
        tension: 0.3,
        pointRadius: 2,
        borderWidth: 1.5
      },
      {
        label: 'Idle',
        data: history.map(h => h.total_idle),
        borderColor: '#10b981',
        backgroundColor: 'rgba(16, 185, 129, 0.1)',
        fill: true,
        tension: 0.3,
        pointRadius: 2,
        borderWidth: 1.5
      }
    ]
  }
})

const poolChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  interaction: {
    intersect: false,
    mode: 'index'
  },
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        usePointStyle: true,
        boxWidth: 6,
        font: { size: 10 }
      }
    },
    tooltip: {
      bodyFont: { size: 11 },
      titleFont: { size: 11 }
    }
  },
  scales: {
    x: {
      grid: { display: false },
      ticks: { font: { size: 10 }, maxTicksLimit: 10 }
    },
    y: {
      beginAtZero: true,
      ticks: {
        font: { size: 10 },
        stepSize: 1
      },
      grid: { color: '#f3f4f6' }
    }
  }
}

// ── Quick Stats ──
const quickStats = computed(() => [
  { label: t('monitoring.stat_active_tenants'), value: tenants.value.length.toString() || '-' },
  { label: t('monitoring.stat_platform_db'), value: overallHealthStatus.value },
  { label: t('monitoring.stat_pool_connections'), value: poolStatValue('total_open').toString() || '-', highlight: utilizationPercent.value > 80, color: utilizationPercent.value > 80 ? 'text-rose-600' : '' },
  { label: t('monitoring.stat_avg_uptime'), value: '-' }
])

// ── Pool cell coloring ──
function poolCellClass(data, field) {
  const val = data.pool?.[field]
  const max = data.pool?.max_open
  if (!val || !max) return 'text-gray-500'
  if (field === 'in_use' && val > max * 0.8) return 'text-rose-600 font-medium'
  if (field === 'open' && val > max * 0.9) return 'text-amber-600 font-medium'
  return 'text-gray-700'
}

// ── Load Data ──
async function loadAll() {
  if (loading.value) return
  loading.value = true
  try {
    const [healthRes, tenantRes] = await Promise.allSettled([
      api.get('/api/v1/platform/monitoring/health'),
      api.get('/api/v1/platform/monitoring/tenants')
    ])

    if (healthRes.status === 'fulfilled') {
      const data = healthRes.value.data || {}
      health.value = data
      // Record pool history for chart
      if (data.pool_stats) {
        poolHistory.value.push({ ...data.pool_stats })
        if (poolHistory.value.length > MAX_HISTORY) {
          poolHistory.value.shift()
        }
      }
    } else {
      health.value = {}
    }

    if (tenantRes.status === 'fulfilled') {
      const payload = tenantRes.value.data
      const tenantData = payload?.data?.tenants || payload?.data || payload || []
      tenants.value = Array.isArray(tenantData) ? tenantData : []
    } else {
      tenants.value = []
    }

    loaded.value = true
    lastUpdated.value = new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
  } catch (e) {
    // silently handle
  } finally {
    loading.value = false
  }
}

// ── Auto-refresh ──
function startAutoRefresh() {
  stopAutoRefresh()
  refreshTimer = setInterval(() => {
    if (!autoRefreshActive.value) return
    loadAll()
  }, 30000)
}

function stopAutoRefresh() {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
}

// Restart timer when toggle changes
watch(autoRefreshActive, (val) => {
  if (val) startAutoRefresh()
  else stopAutoRefresh()
})

onMounted(() => {
  loadAll()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>
