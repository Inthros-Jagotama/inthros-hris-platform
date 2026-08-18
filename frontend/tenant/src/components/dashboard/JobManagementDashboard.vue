<template>
  <div>
    <!-- Skeleton loading -->
    <div v-if="loading" class="space-y-4">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 animate-pulse">
        <div class="h-4 w-48 bg-gray-200 dark:bg-gray-700 rounded mb-3"></div>
        <div class="grid grid-cols-3 gap-2">
          <div v-for="i in 3" :key="i" class="h-8 rounded bg-gray-100 dark:bg-gray-700/50"></div>
        </div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div v-for="i in 2" :key="'d' + i" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 animate-pulse">
          <div class="h-4 w-40 bg-gray-200 dark:bg-gray-700 rounded mb-3"></div>
          <div class="w-32 h-32 mx-auto rounded-full bg-gray-200 dark:bg-gray-700"></div>
        </div>
      </div>
    </div>
    <div v-else>
      <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
        <div class="flex items-center gap-2">
          <i class="pi pi-briefcase text-sm text-teal-500"></i>
          <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.view_job') }}</h2>
          <span class="text-[11px] text-gray-400 dark:text-gray-500 hidden sm:inline">{{ t('dashboard.job_management_desc') }}</span>
        </div>
        <Button :label="t('dashboard.my_kpi_view_all')" icon="pi pi-arrow-right" size="small" text class="!text-xs" @click="$router.push('/job-management')" />
      </div>

      <!-- Belum ada summary aktif -->
      <div v-if="!data?.summary" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-8 text-center">
        <i class="pi pi-sitemap text-2xl text-gray-300 dark:text-gray-600 mb-2 block"></i>
        <p class="text-sm text-gray-400 dark:text-gray-500">{{ t('dashboard.job_no_active_summary') }}</p>
      </div>

      <template v-else>
        <!-- Summary organisasi aktif -->
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 mb-3">
          <div class="flex items-center justify-between gap-2 mb-3">
            <div class="flex items-center gap-2">
              <i class="pi pi-sitemap text-sm text-teal-500"></i>
              <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.job_summary_title') }}</h2>
            </div>
            <Tag :value="data.summary.code" severity="info" class="!text-xs !px-2 !py-0.5" />
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
              <p class="text-[11px] text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_summary_code') }}</p>
              <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ data.summary.code }}</p>
            </div>
            <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
              <p class="text-[11px] text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_summary_decree') }}</p>
              <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ data.summary.decree_no }}</p>
            </div>
            <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
              <p class="text-[11px] text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_summary_decree_date') }}</p>
              <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ formatDate(data.summary.decree_date, locale) }}</p>
            </div>
          </div>
        </div>

        <!-- KPI: jumlah organisasi & progres value -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 mb-3">
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_total_organizations') }}</span>
              <i class="pi pi-building text-lg text-teal-500"></i>
            </div>
            <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ data.total_organizations ?? 0 }}</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_value_not_started') }}</span>
              <i class="pi pi-circle text-lg text-gray-400"></i>
            </div>
            <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ data.value_not_started ?? 0 }}</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_value_on_progress') }}</span>
              <i class="pi pi-spinner text-lg text-amber-500"></i>
            </div>
            <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ data.value_on_progress ?? 0 }}</div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('dashboard.job_value_completed') }}</span>
              <i class="pi pi-check-circle text-lg text-emerald-500"></i>
            </div>
            <div class="text-xl font-bold text-navy-800 dark:text-gray-100">{{ data.value_completed ?? 0 }}</div>
          </div>
        </div>

        <!-- Donut: terisi karyawan & wewenang keuangan -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <!-- Terisi karyawan vs belum -->
          <div v-if="employeeTotal > 0" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.job_employee_fill') }}</h2>
            <div class="flex items-center gap-6 flex-wrap">
              <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
                <circle
                  v-for="seg in employeeSegments"
                  :key="seg.label"
                  cx="60" cy="60" r="45" fill="none" stroke-width="18"
                  :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
                  transform="rotate(-90 60 60)"
                />
                <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ employeeTotal }}</text>
                <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.job_total_organizations') }}</text>
              </svg>
              <div class="space-y-2 flex-1 min-w-0">
                <div v-for="seg in employeeSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
                  <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
                  <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
                  <span class="font-semibold text-navy-800 dark:text-gray-100">{{ seg.value }}</span>
                  <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Wewenang keuangan -->
          <div v-if="financialTotal > 0" class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('dashboard.job_financial_title') }}</h2>
            <div class="flex items-center gap-6 flex-wrap">
              <svg viewBox="0 0 120 120" class="w-40 h-40 shrink-0">
                <circle
                  v-for="seg in financialSegments"
                  :key="seg.label"
                  cx="60" cy="60" r="45" fill="none" stroke-width="18"
                  :stroke="seg.color" :stroke-dasharray="seg.dash" :stroke-dashoffset="seg.offset"
                  transform="rotate(-90 60 60)"
                />
                <text x="60" y="57" text-anchor="middle" class="fill-gray-800 dark:fill-gray-100" style="font-size:20px;font-weight:700">{{ financialTotal }}</text>
                <text x="60" y="74" text-anchor="middle" class="fill-gray-400" style="font-size:9px">{{ t('dashboard.job_total_organizations') }}</text>
              </svg>
              <div class="space-y-2 flex-1 min-w-0">
                <div v-for="seg in financialSegments" :key="seg.label" class="flex items-center gap-2 text-sm">
                  <span class="w-3 h-3 rounded-full shrink-0" :style="{ backgroundColor: seg.color }"></span>
                  <span class="text-gray-600 dark:text-gray-300 flex-1 truncate">{{ seg.label }}</span>
                  <span class="font-semibold text-navy-800 dark:text-gray-100">{{ seg.value }}</span>
                  <span class="text-gray-400 dark:text-gray-500 w-10 text-right shrink-0">{{ seg.pct }}%</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'
import { buildDonutSegments } from '@/utils/dashboard'
import { formatDate } from '@/utils/formatDate'

import Button from 'primevue/button'
import Tag from 'primevue/tag'

const { t, locale } = useI18n()
const toast = useToast()

const loading = ref(false)
const data = ref(null) // GET /job-management/dashboard

// Donut: organisasi terisi karyawan (employment berjalan) vs belum.
const employeeSegments = computed(() => buildDonutSegments([
  { label: t('dashboard.job_with_employees'), value: data.value?.with_employees || 0, color: '#10b981' },
  { label: t('dashboard.job_without_employees'), value: data.value?.without_employees || 0, color: '#9ca3af' }
]))
const employeeTotal = computed(() => employeeSegments.value.reduce((s, i) => s + i.value, 0))

// Donut: organisasi dengan wewenang keuangan vs tidak.
const financialSegments = computed(() => buildDonutSegments([
  { label: t('dashboard.job_with_financial'), value: data.value?.with_financial_authority || 0, color: '#f59e0b' },
  { label: t('dashboard.job_without_financial'), value: data.value?.without_financial_authority || 0, color: '#9ca3af' }
]))
const financialTotal = computed(() => financialSegments.value.reduce((s, i) => s + i.value, 0))

async function loadDashboard() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/job-management/dashboard')
    data.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>
