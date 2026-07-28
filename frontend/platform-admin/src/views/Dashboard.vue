<template>
  <div class="space-y-4">
    <!-- Page Header -->
    <div class="flex items-center justify-end">
      <div class="flex items-center gap-2 text-sm text-gray-400 dark:text-gray-500">
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
        <span v-if="lastUpdated" class="text-gray-400 dark:text-gray-500">
          {{ t('dashboard.updated') }}: {{ lastUpdated }}
        </span>
      </div>
    </div>

    <!-- Transition: skeleton ↔ KPI cards (initial load) -->
    <Transition name="fadeSkeleton" mode="out-in">
      <!-- Full Loading Skeleton (initial load only) -->
      <div v-if="showInitialSkeleton" key="skeleton" class="space-y-4">
        <!-- KPI Cards Skeleton -->
        <SkeletonCard type="kpi" />

        <!-- Chart + Pool Skeleton -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 animate-pulse">
          <div class="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
            <div class="h-4 w-32 bg-gray-200 dark:bg-gray-600 rounded mb-4"></div>
            <div class="flex items-end gap-2 h-36">
              <div v-for="h in [30, 50, 25, 65, 40, 75]" :key="h" class="flex-1 bg-gray-100 dark:bg-gray-700 rounded-t" :style="{ height: h + '%' }"></div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
            <div class="h-4 w-24 bg-gray-200 dark:bg-gray-600 rounded mb-3"></div>
            <div class="flex flex-col items-center">
              <div class="h-10 w-20 bg-gray-200 rounded mb-2"></div>
              <div class="h-3 w-32 bg-gray-200 rounded mb-4"></div>
              <div class="w-full space-y-2">
                <div class="flex items-center justify-between">
                  <div class="h-3 w-10 bg-gray-200 dark:bg-gray-600 rounded"></div>
                  <div class="h-3 w-8 bg-gray-200 dark:bg-gray-600 rounded"></div>
                </div>
                <div class="flex items-center justify-between">
                  <div class="h-3 w-12 bg-gray-200 dark:bg-gray-600 rounded"></div>
                  <div class="h-3 w-8 bg-gray-200 dark:bg-gray-600 rounded"></div>
                </div>
                <div class="flex items-center justify-between">
                  <div class="h-3 w-8 bg-gray-200 dark:bg-gray-600 rounded"></div>
                  <div class="h-3 w-8 bg-gray-200 dark:bg-gray-600 rounded"></div>
                </div>
                <div class="h-1.5 bg-gray-200 dark:bg-gray-600 rounded-full mt-2"></div>
                <div class="h-3 w-24 bg-gray-200 dark:bg-gray-600 rounded mx-auto"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- Two-column Skeleton -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 animate-pulse">
          <!-- Recent Companies -->
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
              <div class="h-4 w-28 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-3 w-12 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
            <div class="p-3 space-y-1">
              <div v-for="i in 3" :key="i" class="flex items-center justify-between px-3 py-2">
                <div class="flex items-center gap-2">
                  <div class="w-4 h-4 bg-gray-200 dark:bg-gray-600 rounded"></div>
                  <div class="h-4 w-28 bg-gray-200 dark:bg-gray-600 rounded"></div>
                </div>
                <div class="h-5 w-14 bg-gray-200 dark:bg-gray-600 rounded-full"></div>
              </div>
            </div>
          </div>
          <!-- System Health -->
          <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
              <div class="h-4 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-3 w-12 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
            <div class="p-3 space-y-3">
              <div v-for="i in 6" :key="i" class="flex items-center justify-between px-2">
                <div class="h-3.5 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
                <div class="h-3.5 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- KPI Cards (initial load complete) -->
      <div v-else key="content" class="space-y-4">
        <!-- KPI Cards -->
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-3">
      <div
        v-for="kpi in kpis"
        :key="kpi.label"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow cursor-pointer group"
      >
        <div class="flex items-center gap-2 mb-1">
          <div class="w-8 h-8 rounded-lg flex items-center justify-center transition-transform group-hover:scale-110" :class="kpi.bg">
            <i :class="kpi.icon" class="text-sm" :style="{ color: kpi.color }"></i>
          </div>
        </div>
        <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ kpi.value }}</p>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ kpi.label }}</p>
        </div>
        </div>

        <!-- Chart Row -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Company Trend Chart -->
      <div class="lg:col-span-2 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.chart_company_trend') }}</h3>
          <Tag :value="t('dashboard.chart_monthly')" severity="info" class="!text-xs" rounded />
        </div>
        <div class="p-3">
          <div v-if="chartData" class="h-48">
            <Chart type="bar" :data="chartData" :options="chartOptions" />
          </div>
          <div v-else class="h-48 flex items-center justify-center text-sm text-gray-400 dark:text-gray-500">
            <i class="pi pi-chart-bar text-2xl mr-2 opacity-50"></i>
            {{ t('common.loading') }}
          </div>
        </div>
      </div>

      <!-- Pool Wait / System Load -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.pool_wait_count') }}</h3>
          <i class="pi pi-database text-gray-400 dark:text-gray-500 text-sm"></i>
        </div>
        <div class="p-4 flex flex-col items-center justify-center h-48">
          <div class="text-4xl font-bold text-gray-800 dark:text-gray-100 mb-1">{{ totalWaitCount }}</div>
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('dashboard.pool_wait_desc') }}</p>
          <!-- Mini sparkline: wait distribution -->
          <div class="w-full mt-4 space-y-2">
            <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
              <span>Open</span>
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ poolSummary.total_open ?? 0 }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
              <span>In Use</span>
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ poolSummary.total_in_use ?? 0 }}</span>
            </div>
            <div class="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
              <span>Idle</span>
              <span class="font-medium text-gray-700 dark:text-gray-200">{{ poolSummary.total_idle ?? 0 }}</span>
            </div>
            <ProgressBar
              :value="poolUtilization"
              :class="poolUtilization > 80 ? '!bg-rose-100' : '!bg-emerald-100'"
              class="!h-1.5 !rounded-full mt-1"
            />
            <div class="text-[10px] text-gray-400 dark:text-gray-500 text-center">{{ poolUtilization }}% {{ t('dashboard.pool_wait_desc') }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Two-column layout -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <!-- Recent Companies -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.recent_companies') }}</h3>
          <router-link to="/companies" class="text-sm text-indigo-600 dark:text-indigo-400 hover:underline">{{ t('common.view_all') }}</router-link>
        </div>
        <div class="p-2">
          <div v-if="recentCompanies.length === 0 && !loading" class="text-sm text-gray-400 dark:text-gray-500 text-center py-4">
            {{ t('dashboard.no_companies') }}
          </div>
          <div
            v-for="company in recentCompanies"
            :key="company.id"
            v-show="!(loading && loaded)"
            class="flex items-center justify-between px-3 py-2 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 text-sm transition-colors"
          >
            <div class="flex items-center gap-2 min-w-0">
              <i class="pi pi-building text-gray-400 dark:text-gray-500 text-sm shrink-0"></i>
              <span class="text-gray-700 dark:text-gray-200 truncate">{{ company.name }}</span>
            </div>
            <Tag :value="company.status" :severity="statusSeverity(company.status)" class="!text-xs !px-1.5 !py-0.5 shrink-0" />
          </div>
          <div v-if="loading && loaded" class="space-y-1 p-2">
            <div v-for="i in 3" :key="i" class="h-10 bg-gray-100 dark:bg-gray-700 rounded-md animate-pulse flex items-center justify-between px-3">
              <div class="flex items-center gap-2">
                <div class="w-4 h-4 bg-gray-200 dark:bg-gray-600 rounded"></div>
                <div class="h-3 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
              </div>
              <div class="h-4 w-12 bg-gray-200 dark:bg-gray-600 rounded-full"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- System Health -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="flex items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('dashboard.system_health') }}</h3>
          <router-link to="/monitoring" class="text-sm text-indigo-600 dark:text-indigo-400 hover:underline">{{ t('common.details') }}</router-link>
        </div>
        <div class="p-3 space-y-2">
          <Transition name="fadeSkeleton" mode="out-in">
            <div v-if="showRefreshSkeleton" key="skeleton" class="space-y-3">
              <div v-for="i in 6" :key="i" class="flex items-center justify-between">
                <div class="h-3 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
                <div class="h-3 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
              </div>
            </div>
            <div v-else key="content">
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.platform_status') }}</span>
              <Tag :value="healthStatus" :severity="healthSeverity" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.platform_db') }}</span>
              <span class="text-sm" :class="platformDbHealthy ? 'text-emerald-600' : 'text-rose-600'">
                <i :class="platformDbHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-circle'" class="mr-1"></i>
                {{ platformDbHealthy ? t('common_status.connected') : t('common_status.unhealthy') }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.redis_cache') }}</span>
              <span class="text-sm" :class="cacheHealthy ? 'text-emerald-600' : 'text-amber-600'">
                <i :class="cacheHealthy ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle'" class="mr-1"></i>
                {{ cacheHealthy ? t('common_status.connected') : cacheStatus }}
              </span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.active_tenants') }}</span>
              <span class="text-sm text-gray-600 dark:text-gray-300 font-medium">{{ activeTenantCount }} {{ t('common_status.connected') }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.pool_connections') }}</span>
              <span class="text-sm text-gray-600 dark:text-gray-300 font-medium">{{ poolStatsText }}</span>
            </div>
            <div class="flex items-center justify-between text-sm">
              <span class="text-gray-600 dark:text-gray-300">{{ t('dashboard.total_users') }}</span>
              <span class="text-sm text-gray-600 dark:text-gray-300 font-medium">{{ totalUsersText }}</span>
            </div>
          </div>
          </Transition>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="flex items-center gap-3">
      <router-link to="/companies" class="flex items-center gap-2 px-3 py-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md text-sm text-gray-600 dark:text-gray-300 hover:border-indigo-200 dark:hover:border-indigo-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">
        <i class="pi pi-plus text-sm"></i> {{ t('dashboard.quick_new_company') }}
      </router-link>
      <router-link to="/users" class="flex items-center gap-2 px-3 py-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md text-sm text-gray-600 dark:text-gray-300 hover:border-indigo-200 dark:hover:border-indigo-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">
        <i class="pi pi-user-plus text-sm"></i> {{ t('dashboard.quick_add_user') }}
      </router-link>
      <router-link to="/monitoring" class="flex items-center gap-2 px-3 py-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md text-sm text-gray-600 dark:text-gray-300 hover:border-indigo-200 dark:hover:border-indigo-500 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">
        <i class="pi pi-chart-bar text-sm"></i> {{ t('dashboard.quick_view_health') }}
      </router-link>
    </div>

    <!-- Seed Data Monitoring -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
      <div class="flex flex-wrap items-center justify-between px-4 py-2.5 border-b border-gray-100 dark:border-gray-700 gap-2">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 flex items-center gap-2">
          <i class="pi pi-database text-indigo-400 text-xs"></i>
          {{ t('seed_data.title') }}
        </h3>
        <div class="flex items-center gap-2">
          <Select
            v-model="selectedCompanyId"
            :options="seedCompanyOptions"
            option-value="id"
            option-label="name"
            :placeholder="t('seed_data.select_company')"
            class="!w-56 !text-xs"
            size="small"
            showClear
            @change="onCompanyChange"
          />
          <Button
            icon="pi pi-refresh"
            size="small"
            severity="secondary"
            text
            :loading="seedDataLoading"
            @click="fetchSeedData"
            :disabled="!selectedCompanyId"
            v-tooltip.left="t('common.refresh')"
          />
        </div>
      </div>

      <div class="p-3">
        <!-- No company selected -->
        <div v-if="!selectedCompanyId" class="flex flex-col items-center justify-center py-6 text-gray-400 dark:text-gray-500">
          <i class="pi pi-building text-2xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('seed_data.no_company') }}</p>
        </div>

        <!-- Loading -->
        <div v-else-if="seedDataLoading" class="space-y-3 animate-pulse">
          <div class="flex gap-3">
            <div v-for="i in 3" :key="i" class="flex-1 bg-gray-100 dark:bg-gray-700 rounded-lg p-3">
              <div class="h-3 w-16 bg-gray-200 dark:bg-gray-600 rounded mb-2"></div>
              <div class="h-6 w-12 bg-gray-200 dark:bg-gray-600 rounded mb-1"></div>
              <div class="h-3 w-20 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
          </div>
          <div class="space-y-1">
            <div v-for="i in 5" :key="i" class="flex items-center justify-between px-3 py-2 bg-gray-50 dark:bg-gray-700 rounded">
              <div class="h-3 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-3 w-8 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
          </div>
        </div>

        <!-- Seed Data Content -->
        <div v-else-if="seedData">
          <!-- Message if connection failed -->
          <div v-if="seedData.message" class="flex items-center gap-2 px-3 py-2 mb-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-700 rounded-md text-sm text-amber-700 dark:text-amber-300">
            <i class="pi pi-exclamation-triangle text-sm"></i>
            <span>{{ seedData.message }}</span>
          </div>

          <!-- Summary Cards -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-2 mb-3">
            <div class="bg-indigo-50 dark:bg-indigo-900/20 rounded-lg p-2.5 text-center">
              <p class="text-xs text-indigo-600 dark:text-indigo-300 font-medium">{{ t('seed_data.total_records') }}</p>
              <p class="text-xl font-bold text-indigo-700 dark:text-indigo-200">{{ seedData.total_records?.toLocaleString() || 0 }}</p>
              <p class="text-[10px] text-indigo-400 dark:text-indigo-400 mt-0.5">{{ seedTables.length }} {{ t('seed_data.tables') }}</p>
            </div>
            <div :class="categoryCardClass('reference')">
              <p class="text-xs font-medium" :class="categoryCardTextClass('reference')">{{ t('seed_data.reference_data') }}</p>
              <p class="text-xl font-bold" :class="categoryCardValueClass('reference')">{{ categoryCount('reference') }}</p>
              <p class="text-[10px] mt-0.5" :class="categoryCardTextClass('reference')">{{ categoryTableCount('reference') }} tables</p>
            </div>
            <div :class="categoryCardClass('region')">
              <p class="text-xs font-medium" :class="categoryCardTextClass('region')">{{ t('seed_data.region_data') }}</p>
              <p class="text-xl font-bold" :class="categoryCardValueClass('region')">{{ categoryCount('region') }}</p>
              <p class="text-[10px] mt-0.5" :class="categoryCardTextClass('region')">{{ categoryTableCount('region') }} tables</p>
            </div>
            <div :class="categoryCardClass('payroll')">
              <p class="text-xs font-medium" :class="categoryCardTextClass('payroll')">{{ t('seed_data.payroll_data') }}</p>
              <p class="text-xl font-bold" :class="categoryCardValueClass('payroll')">{{ categoryCount('payroll') }}</p>
              <p class="text-[10px] mt-0.5" :class="categoryCardTextClass('payroll')">{{ categoryTableCount('payroll') }} tables</p>
            </div>
          </div>

          <!-- Table Details -->
          <div class="max-h-56 overflow-y-auto space-y-0.5 border border-gray-100 dark:border-gray-700 rounded-md">
            <div
              v-for="tbl in seedTables"
              :key="tbl.table"
              class="flex items-center justify-between px-3 py-1.5 text-xs hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
            >
              <div class="flex items-center gap-2 min-w-0">
                <i class="pi pi-table text-gray-300 dark:text-gray-600 text-[10px]"></i>
                <span class="text-gray-700 dark:text-gray-300 truncate">{{ tbl.label }}</span>
                <Tag
                  :value="tbl.category"
                  :severity="categoryTagSeverity(tbl.category)"
                  class="!text-[10px] !px-1 !py-0"
                  rounded
                />
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <span class="font-medium text-gray-700 dark:text-gray-200">{{ getTableCount(tbl.table) }}</span>
                <span v-if="getTableCount(tbl.table) === 0" class="text-rose-400">
                  <i class="pi pi-exclamation-circle text-[10px]"></i>
                </span>
                <span v-else class="text-emerald-400">
                  <i class="pi pi-check-circle text-[10px]"></i>
                </span>
              </div>
            </div>
          </div>

          <!-- Footer hint -->
          <div v-if="missingTableCount > 0" class="mt-2 flex items-center gap-1.5 text-[11px] text-amber-600 dark:text-amber-400">
            <i class="pi pi-info-circle"></i>
            <span>{{ missingTableCount }} empty {{ missingTableCount === 1 ? 'table' : 'tables' }} — {{ t('seed_data.seed_all_hint') }}</span>
          </div>
          <div v-else class="mt-2 flex items-center gap-1.5 text-[11px] text-emerald-600 dark:text-emerald-400">
            <i class="pi pi-check-circle"></i>
            <span>{{ t('seed_data.fully_seeded') }}</span>
          </div>
        </div>
      </div>
    </div>

      </div>
    </Transition>
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
import Select from 'primevue/select'
import SkeletonCard from '@/components/SkeletonCard.vue'
import { useSkeletonPage } from '@/composables/useSkeletonPage'

const toast = useToast()
const { t } = useI18n()
const langStore = useLanguage()

// ── State ──
const { loading, loaded, showInitialSkeleton, showRefreshSkeleton, wrapLoad } = useSkeletonPage()
const lastUpdated = ref('')
const companies = ref([])
const usersTotal = ref(0)
const modulesTotal = ref(0)
const health = ref(null)
const autoRefreshActive = ref(true)
let refreshTimer = null

// ── Seed Data State ──
const selectedCompanyId = ref('')
const seedData = ref(null)
const seedDataLoading = ref(false)

// Companies with active tenant connections for seed data selector
const seedCompanyOptions = computed(() => {
  return safeCompanies.value
    .filter(c => c?.status === 'active' || c?.provision_status === 'provisioned')
    .map(c => ({ id: c.id, name: c.name || c.company_name || c.id }))
})

// Pre-defined list of seed tables with their categories
const seedTables = [
  { table: 'religions', category: 'reference', label: 'Religions' },
  { table: 'educations', category: 'reference', label: 'Educations' },
  { table: 'marital_statuses', category: 'reference', label: 'Marital Statuses' },
  { table: 'relationship_types', category: 'reference', label: 'Relationship Types' },
  { table: 'employment_statuses', category: 'reference', label: 'Employment Statuses' },
  { table: 'banks', category: 'reference', label: 'Banks' },
  { table: 'nationalities', category: 'reference', label: 'Nationalities' },
  { table: 'job_families', category: 'reference', label: 'Job Families' },
  { table: 'countries', category: 'reference', label: 'Countries' },
  { table: 'provinces', category: 'region', label: 'Provinces' },
  { table: 'regencies', category: 'region', label: 'Regencies' },
  { table: 'districts', category: 'region', label: 'Districts' },
  { table: 'villages', category: 'region', label: 'Villages' },
  { table: 'salary_grades', category: 'payroll', label: 'Salary Grades' },
  { table: 'ptkps', category: 'payroll', label: 'PTKPs' },
  { table: 'pph21_ptkp_rates', category: 'payroll', label: 'PPh21 PTKP Rates' },
  { table: 'pph21_tax_brackets', category: 'payroll', label: 'PPh21 Tax Brackets' },
  { table: 'ters', category: 'payroll', label: 'TER Rates' },
  { table: 'bpjs_settings', category: 'payroll', label: 'BPJS Settings' },
  { table: 'bpjs_rate_components', category: 'payroll', label: 'BPJS Rate Components' },
]

// Get count for a specific table from seed data response
function getTableCount(tableName) {
  if (!seedData.value?.tables) return 0
  const tbl = seedData.value.tables.find(t => t.table === tableName)
  return tbl?.count ?? 0
}

// Aggregate count by category
function categoryCount(category) {
  if (!seedData.value?.tables) return 0
  return seedData.value.tables
    .filter(t => t.category === category)
    .reduce((sum, t) => sum + (t.count || 0), 0)
}

// Count tables in a category
function categoryTableCount(category) {
  return seedTables.filter(t => t.category === category).length
}

// Tables with zero count
const missingTableCount = computed(() => {
  if (!seedData.value?.tables) return 0
  return seedData.value.tables.filter(t => !t.count || t.count === 0).length
})

// Category card styling
function categoryCardClass(category) {
  const base = 'rounded-lg p-2.5 text-center '
  const colors = {
    reference: 'bg-sky-50 dark:bg-sky-900/20 ',
    region: 'bg-emerald-50 dark:bg-emerald-900/20 ',
    payroll: 'bg-purple-50 dark:bg-purple-900/20 '
  }
  return base + (colors[category] || 'bg-gray-50 ')
}

function categoryCardTextClass(category) {
  const colors = {
    reference: 'text-sky-600 dark:text-sky-300',
    region: 'text-emerald-600 dark:text-emerald-300',
    payroll: 'text-purple-600 dark:text-purple-300'
  }
  return colors[category] || 'text-gray-600'
}

function categoryCardValueClass(category) {
  const colors = {
    reference: 'text-sky-700 dark:text-sky-200',
    region: 'text-emerald-700 dark:text-emerald-200',
    payroll: 'text-purple-700 dark:text-purple-200'
  }
  return colors[category] || 'text-gray-700'
}

function categoryTagSeverity(category) {
  const severities = { reference: 'info', region: 'success', payroll: 'warn' }
  return severities[category] || 'info'
}

// Fetch seed data when company changes
function onCompanyChange() {
  if (selectedCompanyId.value) {
    fetchSeedData()
  } else {
    seedData.value = null
  }
}

async function fetchSeedData() {
  if (!selectedCompanyId.value) return
  seedDataLoading.value = true
  try {
    const res = await api.get(`/api/v1/platform/monitoring/seed-status?company_id=${selectedCompanyId.value}`)
    seedData.value = res.data?.data || null
  } catch (e) {
    seedData.value = null
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: t('message.failed_to_load'),
      life: 3000
    })
  } finally {
    seedDataLoading.value = false
  }
}

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
  try {
    await wrapLoad(async () => {
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

      lastUpdated.value = formatTime(new Date())
    })
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: t('message.failed_to_load'),
      life: 3000
    })
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

<style scoped>
.fadeSkeleton-enter-active,
.fadeSkeleton-leave-active {
  transition: opacity 0.3s ease;
}
.fadeSkeleton-enter-from,
.fadeSkeleton-leave-to {
  opacity: 0;
}
</style>
