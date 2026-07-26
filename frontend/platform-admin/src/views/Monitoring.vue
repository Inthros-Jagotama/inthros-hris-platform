<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-end">
      <Button label="Refresh" icon="pi pi-refresh" size="small" severity="secondary" :loading="loading" :disabled="loading" @click="loadAll" />
    </div>

    <!-- Platform Health -->
    <div class="bg-white rounded-lg border border-gray-200">
      <div class="px-4 py-2.5 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-gray-700">Platform Health</h3>
        <Tag :value="health.status || 'Checking...'" :severity="health.status === 'healthy' ? 'success' : 'warn'" class="!text-xs" />
      </div>
      <div class="p-3">
        <div v-if="health.database" class="space-y-1.5">
          <div v-for="(val, key) in health.database" :key="key" class="flex items-center justify-between text-sm px-2 py-1.5 rounded-md hover:bg-gray-50">
            <span class="text-gray-600 capitalize">{{ key.replace(/_/g, ' ') }}</span>
            <div class="flex items-center gap-2">
              <span class="text-gray-700 font-medium">{{ val }}</span>
              <i class="pi pi-check-circle text-emerald-400 text-xs"></i>
            </div>
          </div>
        </div>
        <div v-else class="text-sm text-gray-400 py-2">Loading health data...</div>
      </div>
    </div>

    <!-- Tenant Connections -->
    <div class="bg-white rounded-lg border border-gray-200">
      <div class="px-4 py-2.5 border-b border-gray-100 flex items-center justify-between">
        <h3 class="text-sm font-semibold text-gray-700">Tenant Connections</h3>
        <span class="text-sm text-gray-500">{{ tenants.length }} active connections</span>
      </div>
      <div class="p-2">
        <DataTable :value="tenants" size="small" class="!text-sm">
          <template #empty>
            <div class="flex flex-col items-center justify-center py-8 text-gray-400">
              <i class="pi pi-database text-2xl mb-1 opacity-50"></i>
              <p class="text-sm">No tenant connections found.</p>
            </div>
          </template>
          <Column field="company_name" header="Company">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <i class="pi pi-building text-indigo-400 text-sm"></i>
                <span>{{ data.company_name || data.company_id }}</span>
              </div>
            </template>
          </Column>
          <Column field="status" header="Status">
            <template #body="{ data }">
              <Tag :value="data.status || 'connected'" :severity="(data.status || 'connected') === 'healthy' ? 'success' : 'warn'" class="!text-xs" />
            </template>
          </Column>
          <Column field="pool.open" header="Pool Open" />
          <Column field="pool.idle" header="Pool Idle" />
          <Column field="last_active" header="Last Active">
            <template #body="{ data }">{{ data.last_active || '-' }}</template>
          </Column>
          <Column field="driver" header="Driver" />
        </DataTable>
      </div>
    </div>

    <!-- Quick Stats Row -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
      <div v-for="stat in quickStats" :key="stat.label" class="bg-white rounded-lg border border-gray-200 p-3">
        <p class="text-sm text-gray-500">{{ stat.label }}</p>
        <p class="text-lg font-bold text-gray-800 mt-0.5">{{ stat.value }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Button from 'primevue/button'

const health = ref({})
const tenants = ref([])
const loading = ref(false)
const quickStats = ref([
  { label: 'Active Tenants', value: '-' },
  { label: 'Platform DB', value: 'Checking...' },
  { label: 'Pool Connections', value: '-' },
  { label: 'Avg Uptime', value: '-' }
])

onMounted(loadAll)

async function loadAll() {
  if (loading.value) return
  loading.value = true
  try {
    const [healthRes, tenantRes] = await Promise.allSettled([
      api.get('/api/v1/platform/monitoring/health'),
      api.get('/api/v1/platform/monitoring/tenants')
    ])
    if (healthRes.status === 'fulfilled') {
      health.value = healthRes.value.data || {}
      quickStats.value[1].value = health.value.status || 'Unknown'
    }

    if (tenantRes.status === 'fulfilled') {
      const payload = tenantRes.value.data
      const tenantData = payload?.data?.tenants || payload?.data || payload || []
      tenants.value = Array.isArray(tenantData) ? tenantData : []
      quickStats.value[0].value = tenants.value.length.toString()
      const totalOpen = tenants.value.reduce((s, t) => s + (t.pool.open || 0), 0)
      quickStats.value[2].value = totalOpen.toString()

      const healthyTenants = payload.data.tenants.filter(t => t.status === 'healthy').length;
      const healthPercentage = payload.data.tenants > 0 
        ? Math.round((healthyTenants / payload.data.tenants) * 100) 
        : 0;
      quickStats.value[3].value = healthPercentage.toString()

    }
  } catch (e) {
    // silently handle network errors
  } finally {
    loading.value = false
  }
}
</script>
