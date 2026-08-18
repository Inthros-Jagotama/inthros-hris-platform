<template>
  <div class="space-y-4">
    <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
      <FormRow :label="t('performance_scoring.organization')">
        <Select
          v-model="selectedOrgId"
          :options="organizationOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('performance_scoring.select_organization')"
          class="w-full md:w-96"
          filter
          showClear
          @change="loadConfig"
        />
      </FormRow>
    </div>

    <div v-if="!selectedOrgId" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-sliders-h text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('performance_scoring.select_organization_hint') }}</p>
    </div>

    <template v-else>
      <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="4" />

      <div v-else class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 flex items-center gap-2">
            <i class="pi pi-sliders-h text-purple-500"></i>
            {{ t('performance_scoring.components') }}
          </h3>
          <span class="text-xs" :class="totalWeight === 100 ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400'">
            {{ t('okr.total_weight') }}: {{ totalWeight.toFixed(2) }}%
          </span>
        </div>

        <div v-if="rows.length === 0" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-lg">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('performance_scoring.no_master_components') }}</p>
        </div>

        <DataTable v-else :value="rows" size="small" class="!text-sm p-datatable-sm" :rowHover="true">
          <Column :header="t('performance_components.code')" style="width:160px">
            <template #body="{data}">
              <Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column :header="t('performance_components.name')" style="min-width:180px">
            <template #body="{data}">
              <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
            </template>
          </Column>
          <Column :header="t('performance_scoring.enabled')" style="width:100px">
            <template #body="{data}">
              <ToggleSwitch v-model="data.is_enabled" />
            </template>
          </Column>
          <Column :header="t('okr.weight')" style="width:140px">
            <template #body="{data}">
              <InputNumber v-model="data.weight" :min="0" :max="100" :minFractionDigits="2" :maxFractionDigits="2" suffix="%" class="w-full !text-xs" size="small" :disabled="!data.is_enabled" />
            </template>
          </Column>
          <Column :header="t('performance_components.sort_order')" style="width:80px">
            <template #body="{data}">
              <InputNumber v-model="data.sort_order" :min="0" class="w-full !text-xs" size="small" />
            </template>
          </Column>
        </DataTable>

        <div class="flex items-center justify-end gap-2 mt-4">
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'

const { t } = useI18n()
const toast = useToast()

const organizationOptions = ref([])
const selectedOrgId = ref(null)
const loading = ref(false)
const saving = ref(false)
const masterComponents = ref([])
const rows = ref([])

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' }
]

const totalWeight = computed(() => {
  return rows.value.filter(r => r.is_enabled).reduce((sum, r) => sum + (r.weight || 0), 0)
})

async function loadOrganizations() {
  try {
    const res = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
    organizationOptions.value = (res.data?.data || []).map(o => ({
      label: o.nomenclature || o.name || o.code,
      value: o.id
    }))
  } catch {
    // silently fail
  }
}

async function loadMasterComponents() {
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/components', { params: { per_page: 100 } })
    masterComponents.value = (res.data?.data || []).filter(c => c.is_active)
  } catch {
    masterComponents.value = []
  }
}

async function loadConfig() {
  if (!selectedOrgId.value) {
    rows.value = []
    return
  }
  loading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/performance/kpi/organizations/${selectedOrgId.value}/components`)
    const existing = res.data?.data || []
    const existingByComponentId = new Map(existing.map(e => [e.component_id, e]))

    rows.value = masterComponents.value.map(c => {
      const ex = existingByComponentId.get(c.id)
      return {
        component_id: c.id,
        code: c.code,
        name: c.name,
        weight: ex?.weight ?? 0,
        is_enabled: ex?.is_enabled ?? false,
        sort_order: ex?.sort_order ?? c.sort_order ?? 0
      }
    })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
    rows.value = []
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!selectedOrgId.value) return

  const enabledRows = rows.value.filter(r => r.is_enabled)
  if (enabledRows.length > 0) {
    const total = enabledRows.reduce((sum, r) => sum + (r.weight || 0), 0)
    if (Math.round(total * 100) / 100 !== 100) {
      toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('performance_scoring.weight_must_100'), life: 4000 })
      return
    }
  }

  saving.value = true
  try {
    for (const row of rows.value) {
      await api.post('/api/v1/tenant/performance/kpi/organization-components', {
        organization_id: selectedOrgId.value,
        component_id: row.component_id,
        weight: row.weight || 0,
        is_enabled: row.is_enabled,
        sort_order: row.sort_order || 0
      })
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('performance_scoring.saved'), life: 3000 })
    await loadConfig()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  await Promise.all([loadOrganizations(), loadMasterComponents()])
})
</script>
