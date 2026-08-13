<template>
  <div class="space-y-4">
    <!-- Toolbar: ringkasan + refresh -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1.5">
        <i class="pi pi-info-circle text-xs"></i>{{ t('recruitment_analytics.subtitle') }}
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

    <!-- Loading skeleton -->
    <div v-if="loading" class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
      <div v-for="i in 8" :key="i" class="h-28 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
    </div>

    <!-- Error state -->
    <div
      v-else-if="loadFailed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-rose-200 dark:border-rose-700/50 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-exclamation-triangle text-4xl mb-3 text-rose-400 opacity-70"></i>
      <p class="text-sm font-medium text-rose-600 dark:text-rose-400">{{ t('message.failed_to_load') }}</p>
      <Button :label="t('common.retry')" icon="pi pi-refresh" size="small" class="!mt-4 !text-xs" @click="loadData()" />
    </div>

    <template v-else>
      <!-- KPI cards -->
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
        <!-- Remaining gap (S-2) — kartu utama -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-sky-300 dark:border-sky-700/50 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.remaining_gap') }}</span>
            <i class="pi pi-arrows-alt text-base text-sky-500"></i>
          </div>
          <p class="text-2xl font-bold text-sky-600 dark:text-sky-400">{{ data.remaining_gap }}</p>
          <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-1">{{ t('recruitment_analytics.remaining_gap_hint') }}</p>
        </div>

        <!-- Expected hires (S-2) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.expected_hires') }}</span>
            <i class="pi pi-user-plus text-base text-emerald-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ data.expected_hires }}</p>
        </div>

        <!-- Open positions -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.open_positions') }}</span>
            <i class="pi pi-briefcase text-base text-indigo-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ data.open_positions }}</p>
        </div>

        <!-- Filled positions -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.filled_positions') }}</span>
            <i class="pi pi-check-circle text-base text-teal-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ data.filled_positions }}</p>
        </div>

        <!-- Time to hire (S-3) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.time_to_hire') }}</span>
            <i class="pi pi-clock text-base text-amber-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ fmtDays(data.time_to_hire) }}</p>
        </div>

        <!-- Time to fill (S-3) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.time_to_fill') }}</span>
            <i class="pi pi-hourglass text-base text-rose-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ fmtDays(data.time_to_fill) }}</p>
        </div>

        <!-- Offer acceptance rate (S-3) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.offer_acceptance') }}</span>
            <i class="pi pi-percentage text-base text-violet-500"></i>
          </div>
          <p class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ fmtPct(data.offer_acceptance_rate) }}</p>
        </div>

        <!-- Candidate match (placeholder) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 p-3.5 opacity-70">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('recruitment_analytics.candidate_match') }}</span>
            <i class="pi pi-filter text-base text-gray-400"></i>
          </div>
          <p class="text-2xl font-bold text-gray-400 dark:text-gray-500">{{ fmtScore(data.candidate_match_score) }}</p>
          <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-1">{{ t('recruitment_analytics.placeholder_note') }}</p>
        </div>
      </div>

      <!-- Source conversion + pipeline -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <!-- Source conversion -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div class="px-3 py-2.5 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('recruitment_analytics.source_conversion') }}</h2>
          </div>
          <div v-if="data.source_conversion.length === 0" class="p-6 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500">
            <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
            <p class="text-sm">{{ t('recruitment_analytics.no_data') }}</p>
          </div>
          <DataTable v-else :value="data.source_conversion" size="small" class="!text-sm">
            <Column :header="t('recruitment_analytics.source')">
              <template #body="{ data: row }">
                <span class="font-medium text-gray-700 dark:text-gray-200 capitalize">{{ row.source || '—' }}</span>
              </template>
            </Column>
            <Column :header="t('recruitment_analytics.candidates')" style="width: 90px">
              <template #body="{ data: row }"><span class="text-gray-600 dark:text-gray-300">{{ row.candidates }}</span></template>
            </Column>
            <Column :header="t('recruitment_analytics.hires')" style="width: 90px">
              <template #body="{ data: row }"><Tag :value="row.hires" severity="success" class="!text-xs !px-2 !py-0.5" /></template>
            </Column>
            <Column :header="t('recruitment_analytics.conversion')" style="width: 140px">
              <template #body="{ data: row }">
                <div class="flex items-center gap-2">
                  <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 w-12">{{ fmtPct(row.conversion_rate) }}</span>
                  <div class="flex-1 h-1.5 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
                    <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${pct(row.conversion_rate)}%` }"></div>
                  </div>
                </div>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- Pipeline funnel -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
          <div class="px-3 py-2.5 border-b border-gray-200 dark:border-gray-700">
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('recruitment_analytics.pipeline') }}</h2>
          </div>
          <div v-if="data.pipeline.length === 0" class="p-6 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500">
            <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
            <p class="text-sm">{{ t('recruitment_analytics.no_data') }}</p>
          </div>
          <DataTable v-else :value="data.pipeline" size="small" class="!text-sm">
            <Column :header="t('recruitment_analytics.stage')">
              <template #body="{ data: row }">
                <span class="font-medium text-gray-700 dark:text-gray-200 capitalize">{{ row.label || '—' }}</span>
              </template>
            </Column>
            <Column :header="t('recruitment_analytics.count')" style="width: 160px">
              <template #body="{ data: row }">
                <div class="flex items-center gap-2">
                  <Tag :value="row.value" severity="info" class="!text-xs !px-2 !py-0.5" />
                  <div class="flex-1 h-1.5 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
                    <div class="h-full rounded-full bg-sky-500 transition-all" :style="{ width: `${maxPct(row.value)}%` }"></div>
                  </div>
                </div>
              </template>
            </Column>
          </DataTable>
        </div>
      </div>
    </template>
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

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const loadFailed = ref(false)
const data = ref({
  time_to_hire: 0,
  time_to_fill: 0,
  offer_acceptance_rate: 0,
  candidate_match_score: 0,
  cost_per_hire: 0,
  expected_hires: 0,
  open_positions: 0,
  filled_positions: 0,
  remaining_gap: 0,
  source_conversion: [],
  pipeline: []
})

function fmtDays(v) {
  const n = Number(v) || 0
  return `${n.toFixed(1)}d`
}

function fmtPct(v) {
  const n = Number(v) || 0
  return `${n.toFixed(1)}%`
}

function fmtScore(v) {
  const n = Number(v) || 0
  return n.toFixed(1)
}

function pct(v) {
  return Math.min(100, Math.max(0, Number(v) || 0))
}

function maxPct(v) {
  const max = Math.max(...data.value.pipeline.map(r => Number(r.value) || 0), 1)
  return Math.min(100, Math.max(0, ((Number(v) || 0) / max) * 100))
}

async function loadData() {
  loading.value = true
  loadFailed.value = false
  try {
    const res = await api.get('/api/v1/tenant/workforce-intelligence/analytics/recruitment')
    const body = res.data?.data
    if (body) data.value = body
  } catch (e) {
    loadFailed.value = true
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
