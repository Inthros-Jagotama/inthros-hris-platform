<template>
  <div class="space-y-4">
    <!-- Toolbar: pilih posisi target -->
    <div class="flex items-center gap-2 flex-wrap">
      <SelectLabel
        v-model="positionId"
        :options="positionOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('internal_candidates.select_position')"
        class="!w-72"
        showClear
        @update:modelValue="onSelectPosition()"
      />
      <Button
        :label="t('common.refresh')"
        icon="pi pi-refresh"
        size="small"
        text
        class="!text-xs"
        :loading="loading"
        @click="loadData()"
      />
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('internal_candidates.candidates') }}</span>
    </div>

    <!-- Empty: belum pilih posisi -->
    <div
      v-if="!positionId && !loading"
      class="bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-sitemap text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('internal_candidates.select_hint') }}</p>
      <p class="text-xs mt-1 text-center max-w-sm">{{ t('internal_candidates.select_hint_desc') }}</p>
    </div>

    <SkeletonTable v-else-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="items"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-user-plus text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('internal_candidates.empty') }}</p>
        </div>
      </template>

      <Column :header="t('internal_candidates.employee')">
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-full bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400 flex items-center justify-center text-[11px] font-semibold shrink-0">
              {{ initials(data) }}
            </div>
            <div class="min-w-0">
              <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ data.Name || data.name }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 font-mono truncate">{{ data.EmployeeID || data.employee_id }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('internal_candidates.current_position')">
        <template #body="{ data }">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ data.CurrentPositionName || data.current_position_name || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('internal_candidates.step_sequence')" style="width: 140px">
        <template #body="{ data }">
          <Tag :value="data.SourceStepSequence ?? data.source_step_sequence ?? 0" severity="info" class="!text-xs !px-2 !py-0.5" />
        </template>
      </Column>

      <Column :header="t('internal_candidates.target_position')">
        <template #body="{ data }">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ data.TargetPositionName || data.target_position_name || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('internal_candidates.career_path')">
        <template #body="{ data }">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ data.PathName || data.path_name || '—' }}</span>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(false)
const items = ref([])
const positionId = ref(null)
const positions = ref([])

const skeletonColumns = [
  { field: 'employee', header: 'Employee', width: '30%' },
  { field: 'current', header: 'Current Position', width: '25%' },
  { field: 'step', header: 'Step', width: '12%' },
  { field: 'target', header: 'Target', width: '20%' }
]

const positionOptions = computed(() => positions.value.map(p => ({ label: p.nomenclature || p.full_code, value: p.id })))
const totalRecords = computed(() => items.value.length)

function initials(data) {
  const name = data.Name || data.name || '?'
  return name.charAt(0).toUpperCase()
}

async function loadPositions() {
  try {
    const res = await api.get('/api/v1/tenant/positions', { params: { per_page: 500 } })
    positions.value = res.data?.data || []
  } catch {
    // fail-silent — dropdown kosong
  }
}

async function loadData() {
  if (!positionId.value) {
    items.value = []
    return
  }
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/recruitment/eligible-internal-candidates', { params: { position_id: positionId.value } })
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onSelectPosition() {
  items.value = []
  loadData()
}

onMounted(() => {
  loadPositions()
})
</script>
