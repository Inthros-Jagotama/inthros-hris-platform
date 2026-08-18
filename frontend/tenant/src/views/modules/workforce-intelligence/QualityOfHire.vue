<template>
  <div class="space-y-4">
    <!-- Toolbar: ringkasan + refresh -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <p class="text-xs text-gray-500 dark:text-gray-400 flex items-center gap-1.5">
        <i class="pi pi-info-circle text-xs"></i>{{ t('quality_of_hire.subtitle') }}
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

    <!-- Error state (gagal load — beda dari empty data) -->
    <div
      v-else-if="loadFailed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-rose-200 dark:border-rose-700/50 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-exclamation-triangle text-4xl mb-3 text-rose-400 opacity-70"></i>
      <p class="text-sm font-medium text-rose-600 dark:text-rose-400">{{ t('message.failed_to_load') }}</p>
      <Button :label="t('common.retry')" icon="pi pi-refresh" size="small" class="!mt-4 !text-xs" @click="loadData()" />
    </div>

    <!-- Empty state -->
    <div
      v-else-if="!data.hires_analyzed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-chart-line text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('quality_of_hire.empty_title') }}</p>
      <p class="text-xs mt-1 text-center max-w-sm">{{ t('quality_of_hire.empty_desc') }}</p>
    </div>

    <template v-else>
      <!-- KPI cards -->
      <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-3">
        <!-- Overall score — kartu utama -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-emerald-300 dark:border-emerald-700/50 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.overall_score') }}</span>
            <i class="pi pi-bullseye text-base text-emerald-500"></i>
          </div>
          <p class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{{ fmtScore(data.overall_score) }}</p>
          <div class="mt-2 h-1.5 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
            <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${scorePct(data.overall_score)}%` }"></div>
          </div>
        </div>

        <!-- Hires analyzed -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.hires_analyzed') }}</span>
            <i class="pi pi-user-plus text-base text-sky-500"></i>
          </div>
          <p class="text-2xl font-bold text-navy-800 dark:text-gray-100">{{ data.hires_analyzed }}</p>
        </div>

        <!-- Interview score -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.interview_score') }}</span>
            <i class="pi pi-comments text-base text-teal-500"></i>
          </div>
          <p class="text-2xl font-bold text-navy-800 dark:text-gray-100">{{ fmtScore(data.interview_score) }}</p>
        </div>

        <!-- Performance score -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.performance_score') }}</span>
            <i class="pi pi-star text-base text-violet-500"></i>
          </div>
          <p class="text-2xl font-bold text-navy-800 dark:text-gray-100">{{ fmtScore(data.performance_score) }}</p>
        </div>

        <!-- Onboarding completion rate -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.onboarding_completion') }}</span>
            <i class="pi pi-check-circle text-base text-emerald-500"></i>
          </div>
          <p class="text-2xl font-bold text-navy-800 dark:text-gray-100">{{ fmtPct(data.onboarding_completion_rate) }}</p>
        </div>

        <!-- Retention rate -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3.5">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.retention_rate') }}</span>
            <i class="pi pi-users text-base text-teal-500"></i>
          </div>
          <p class="text-2xl font-bold text-navy-800 dark:text-gray-100">{{ fmtPct(data.retention_rate) }}</p>
        </div>

        <!-- Recruitment match (placeholder) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 p-3.5 opacity-70">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.recruitment_match') }}</span>
            <i class="pi pi-filter text-base text-gray-400"></i>
          </div>
          <p class="text-2xl font-bold text-gray-400 dark:text-gray-500">{{ fmtScore(data.recruitment_match_score) }}</p>
          <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-1">{{ t('quality_of_hire.placeholder_note') }}</p>
        </div>

        <!-- Assessment (placeholder) -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 p-3.5 opacity-70">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('quality_of_hire.assessment_score') }}</span>
            <i class="pi pi-file-check text-base text-gray-400"></i>
          </div>
          <p class="text-2xl font-bold text-gray-400 dark:text-gray-500">{{ fmtScore(data.assessment_score) }}</p>
          <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-1">{{ t('quality_of_hire.placeholder_note') }}</p>
        </div>
      </div>

      <!-- Breakdown tabs -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div class="flex items-center gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto px-2">
          <button
            v-for="tab in breakdownTabs"
            :key="tab.key"
            type="button"
            class="px-3 py-2.5 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
            :class="activeBreakdown === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
            @click="activeBreakdown = tab.key"
          >
            {{ t(tab.labelKey) }}
          </button>
        </div>

        <div v-if="breakdownRows.length === 0" class="p-8 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500">
          <i class="pi pi-table text-3xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('quality_of_hire.no_breakdown') }}</p>
        </div>

        <DataTable
          v-else
          :value="breakdownRows"
          size="small"
          class="!text-sm"
        >
          <Column :header="t('quality_of_hire.breakdown_key')">
            <template #body="{ data: row }">
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ row.name || row.key }}</span>
            </template>
          </Column>
          <Column :header="t('quality_of_hire.breakdown_hires')" style="width: 160px">
            <template #body="{ data: row }">
              <Tag :value="row.hires" :severity="row.hires > 0 ? 'info' : 'secondary'" class="!text-xs !px-2 !py-0.5" />
            </template>
          </Column>
          <Column :header="t('quality_of_hire.breakdown_score')" style="width: 240px">
            <template #body="{ data: row }">
              <div class="flex items-center gap-2">
                <span class="text-sm font-semibold text-gray-700 dark:text-gray-200 w-10">{{ fmtScore(row.score) }}</span>
                <div class="flex-1 h-1.5 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
                  <div class="h-full rounded-full bg-emerald-500 transition-all" :style="{ width: `${scorePct(row.score)}%` }"></div>
                </div>
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const loadFailed = ref(false)
const data = ref({
  overall_score: 0,
  hires_analyzed: 0,
  recruitment_match_score: 0,
  interview_score: 0,
  assessment_score: 0,
  onboarding_completion_rate: 0,
  performance_score: 0,
  retention_rate: 0,
  by_source: [],
  by_requisition: [],
  by_organization: []
})
const activeBreakdown = ref('source')

const breakdownTabs = computed(() => [
  { key: 'source', labelKey: 'quality_of_hire.tab_source' },
  { key: 'requisition', labelKey: 'quality_of_hire.tab_requisition' },
  { key: 'organization', labelKey: 'quality_of_hire.tab_organization' }
])

const breakdownRows = computed(() => {
  const map = {
    source: data.value.by_source,
    requisition: data.value.by_requisition,
    organization: data.value.by_organization
  }
  return map[activeBreakdown.value] || []
})

// Skor & rate backend berskala 0–100.
function fmtScore(v) {
  const n = Number(v) || 0
  return n.toFixed(1)
}

function fmtPct(v) {
  const n = Number(v) || 0
  return `${n.toFixed(1)}%`
}

function scorePct(v) {
  const n = Number(v) || 0
  return Math.min(100, Math.max(0, n))
}

async function loadData() {
  loading.value = true
  loadFailed.value = false
  try {
    const res = await api.get('/api/v1/tenant/workforce-intelligence/analytics/quality-of-hire')
    data.value = res.data?.data || data.value
  } catch (e) {
    loadFailed.value = true
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>
