<template>
  <div class="space-y-6">
    <!-- ── Skeleton loading ── -->
    <div v-if="loading" class="space-y-6">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div v-for="n in 3" :key="n" class="flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 animate-pulse">
          <div class="w-11 h-11 rounded-lg bg-gray-200 dark:bg-gray-700 shrink-0"></div>
          <div class="flex-1 space-y-2">
            <div class="h-3.5 bg-gray-200 dark:bg-gray-700 rounded w-2/3"></div>
            <div class="h-3 bg-gray-100 dark:bg-gray-700 rounded w-full"></div>
          </div>
        </div>
      </div>
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 animate-pulse">
        <div class="h-4 bg-gray-200 dark:bg-gray-700 rounded w-40 mb-3"></div>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
          <div v-for="n in 6" :key="n" class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50"></div>
        </div>
      </div>
    </div>

    <template v-else-if="!employeeId && !isAdmin">
      <Message severity="warn" :closable="false">{{ t('reimbursement.no_employee_linked') }}</Message>
    </template>

    <template v-else>
      <!-- ── Menu Cards ── -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <button
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-indigo-300 dark:hover:border-indigo-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500/50"
          @click="router.push('/reimbursements/all')"
        >
          <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400">
            <i class="pi pi-briefcase text-lg"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('reimbursement.card_all') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('reimbursement.card_all_desc') }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-indigo-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>

        <button
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-emerald-300 dark:hover:border-emerald-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
          @click="router.push('/reimbursements/my-requests')"
        >
          <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
            <i class="pi pi-credit-card text-lg"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('reimbursement.my_requests') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('reimbursement.card_my_desc') }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>

        <button
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-4 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-violet-300 dark:hover:border-violet-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-violet-500/50"
          @click="router.push('/reimbursements/types')"
        >
          <div class="w-11 h-11 rounded-lg shrink-0 flex items-center justify-center bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400">
            <i class="pi pi-tags text-lg"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('reimbursement.types') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t('reimbursement.card_types_desc') }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-violet-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>
      </div>

      <!-- ── Dashboard Summary ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
          <div class="flex items-center gap-2">
            <i class="pi pi-chart-bar text-sm text-indigo-500"></i>
            <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('reimbursement.summary') }}</h2>
            <span class="text-xs text-gray-400 dark:text-gray-500">{{ isAdmin ? t('reimbursement.card_all') : t('reimbursement.my_requests') }}</span>
          </div>
        </div>

        <div v-if="summaryLoading" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
          <div v-for="n in 6" :key="n" class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        </div>
        <div v-else-if="summaryCards.length === 0" class="text-center py-8 text-gray-400 dark:text-gray-500">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('reimbursement.summary_empty') }}</p>
        </div>
        <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
          <div
            v-for="card in summaryCards"
            :key="card.key"
            class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5 flex items-center justify-between gap-2 hover:shadow-sm dark:hover:shadow-gray-900/50 transition-shadow cursor-pointer"
            @click="router.push(card.path)"
          >
            <div class="min-w-0">
              <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ card.label }}</p>
              <p class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ card.value }}</p>
            </div>
            <i :class="[card.icon, card.iconColor]" class="text-base shrink-0"></i>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { useMyEmployee } from '@/composables/useMyEmployee'
import api from '@/services/api'

import Message from 'primevue/message'

const router = useRouter()
const { t } = useI18n()
const { hasPermission } = useAuth()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const loading = ref(true)
const summaryLoading = ref(false)
const summaryCounts = ref({})

// Admin/HR (who can approve) see all requests; employees default to their own.
const isAdmin = computed(() => hasPermission('reimbursement.approve'))

const summaryCards = computed(() => {
  const c = summaryCounts.value
  const basePath = isAdmin.value ? '/reimbursements/all' : '/reimbursements/my-requests'
  return [
    { key: 'total', label: t('reimbursement.summary_total'), value: c.total || 0, icon: 'pi pi-inbox', iconColor: 'text-indigo-500', path: basePath },
    { key: 'DRAFT', label: statusLabel('DRAFT'), value: c.DRAFT || 0, icon: 'pi pi-pencil', iconColor: 'text-gray-400', path: `${basePath}?status=DRAFT` },
    { key: 'PENDING', label: t('reimbursement.summary_pending'), value: (c.SUBMITTED || 0) + (c.PENDING_APPROVAL || 0), icon: 'pi pi-clock', iconColor: 'text-amber-500', path: `${basePath}?status=PENDING_APPROVAL` },
    { key: 'APPROVED', label: statusLabel('APPROVED'), value: c.APPROVED || 0, icon: 'pi pi-check-circle', iconColor: 'text-emerald-500', path: `${basePath}?status=APPROVED` },
    { key: 'PAID', label: statusLabel('PAID'), value: c.PAID || 0, icon: 'pi pi-dollar', iconColor: 'text-teal-500', path: `${basePath}?status=PAID` },
    { key: 'REJECTED', label: statusLabel('REJECTED'), value: c.REJECTED || 0, icon: 'pi pi-times-circle', iconColor: 'text-rose-500', path: `${basePath}?status=REJECTED` }
  ]
})

function statusLabel(status) {
  const key = `reimbursement.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadSummary() {
  summaryLoading.value = true
  try {
    const statuses = ['', 'DRAFT', 'SUBMITTED', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED', 'PAID', 'CANCELLED']
    const baseParams = { per_page: 1 }
    if (!isAdmin.value) {
      if (!employeeId.value) return
      baseParams.employee_id = employeeId.value
    }
    const results = await Promise.all(
      statuses.map(status => {
        const params = { ...baseParams }
        if (status) params.status = status
        return api.get('/api/v1/tenant/reimbursements/requests', { params })
          .then(res => ({ status, total: res.data?.total || 0 }))
          .catch(() => ({ status, total: 0 }))
      })
    )
    const counts = {}
    for (const r of results) {
      counts[r.status || 'total'] = r.total
    }
    summaryCounts.value = counts
  } finally {
    summaryLoading.value = false
  }
}

async function loadAll() {
  loading.value = true
  try {
    employeeId.value = await loadMyEmployeeId()
    await loadSummary()
  } catch {
    // summary gagal — halaman index tetap tampil
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>
