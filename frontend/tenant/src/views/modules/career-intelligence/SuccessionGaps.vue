<template>
  <div class="space-y-4">
    <!-- Toolbar: ringkasan + refresh -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1.5">
        <i class="pi pi-info-circle text-xs"></i>{{ t('succession_gaps.subtitle') }}
      </p>
      <Button
        :label="t('common.refresh')"
        icon="pi pi-refresh"
        size="small"
        text
        class="!text-xs"
        :loading="loading"
        @click="loadData()"
      />
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <!-- Error state -->
    <div
      v-else-if="loadFailed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-rose-200 dark:border-rose-700/50 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-exclamation-triangle text-4xl mb-3 text-rose-400 opacity-70"></i>
      <p class="text-sm font-medium text-rose-600 dark:text-rose-400">{{ t('message.failed_to_load') }}</p>
      <Button :label="t('common.retry')" icon="pi pi-refresh" size="small" class="!mt-4 !text-xs" @click="loadData()" />
    </div>

    <DataTable
      v-else
      :value="items"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-check-circle text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('succession_gaps.empty') }}</p>
        </div>
      </template>

      <Column :header="t('succession_gaps.position')">
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center shrink-0">
              <i class="pi pi-shield text-xs"></i>
            </div>
            <div class="min-w-0">
              <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ data.position_title }}</p>
              <p v-if="data.position_code" class="text-xs text-gray-400 dark:text-gray-500 font-mono truncate">{{ data.position_code }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('succession_gaps.successors')" style="width: 130px">
        <template #body="{ data }">
          <Tag :value="data.successor_count" severity="info" class="!text-xs !px-2 !py-0.5" />
        </template>
      </Column>

      <Column :header="t('succession_gaps.ready_now')" style="width: 130px">
        <template #body="{ data }">
          <Tag :value="data.ready_now_count" :severity="data.ready_now_count > 0 ? 'success' : 'secondary'" class="!text-xs !px-2 !py-0.5" />
        </template>
      </Column>

      <Column :header="t('succession_gaps.has_ready')" style="width: 160px">
        <template #body="{ data }">
          <Tag
            :value="data.has_ready_successor ? t('common.yes') : t('common.no')"
            :severity="data.has_ready_successor ? 'success' : 'danger'"
            class="!text-xs !px-2 !py-0.5"
          />
        </template>
      </Column>

      <Column :header="t('succession_gaps.external')" style="width: 200px">
        <template #body="{ data }">
          <Tag
            :value="data.requires_external_recruitment ? t('succession_gaps.requires_external') : t('succession_gaps.covered')"
            :severity="data.requires_external_recruitment ? 'danger' : 'success'"
            class="!text-xs !px-2 !py-0.5"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const loadFailed = ref(false)
const items = ref([])

const skeletonColumns = [
  { field: 'position', header: 'Position', width: '30%' },
  { field: 'successors', header: 'Successors', width: '15%' },
  { field: 'ready', header: 'Ready Now', width: '15%' },
  { field: 'external', header: 'External', width: '20%' }
]

async function loadData() {
  loading.value = true
  loadFailed.value = false
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/successions/gaps')
    items.value = res.data?.data || []
  } catch (e) {
    loadFailed.value = true
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
