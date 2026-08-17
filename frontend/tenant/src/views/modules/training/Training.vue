<template>
  <div class="space-y-4">
    <!-- ── Stat ringkas (P0-FE: dihitung dari endpoint list existing) ── -->
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_upcoming') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.upcoming }}</p>
        </div>
        <i class="pi pi-calendar-plus text-lg text-emerald-500 shrink-0"></i>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_ongoing') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.ongoing }}</p>
        </div>
        <i class="pi pi-play-circle text-lg text-sky-500 shrink-0"></i>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_completed') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.completed }}</p>
        </div>
        <i class="pi pi-check-circle text-lg text-indigo-500 shrink-0"></i>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_total_courses') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.courses }}</p>
        </div>
        <i class="pi pi-book text-lg text-amber-500 shrink-0"></i>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_total_participants') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.participants }}</p>
        </div>
        <i class="pi pi-users text-lg text-violet-500 shrink-0"></i>
      </div>
      <div class="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-3.5 flex items-center justify-between hover:shadow-md transition-shadow">
        <div class="min-w-0">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider truncate">{{ t('training.stat_total_providers') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ stats.providers }}</p>
        </div>
        <i class="pi pi-building text-lg text-rose-500 shrink-0"></i>
      </div>
    </div>

    <!-- ── Card navigasi, dikelompokkan per kategori: Pengaturan / Operasional / Rekam Data / Laporan ── -->
    <div v-for="group in menuGroups" :key="group.titleKey" class="space-y-2">
      <div class="md:col-span-2">
        <div class="flex items-center gap-2 pt-2">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase">{{ t(group.titleKey) }}</span>
          <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
        </div>
      </div>
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
        <button
          v-for="menu in group.items"
          :key="menu.route"
          type="button"
          class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-emerald-300 dark:hover:border-emerald-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
          @click="router.push(menu.route)"
        >
          <div class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center transition-colors" :class="menu.tint">
            <i :class="menu.icon" class="text-base"></i>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ t(menu.labelKey) }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(menu.descKey) }}</p>
          </div>
          <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
        </button>
      </div>
    </div>

    <!-- Fitur phase berikutnya (P1/P2) — card "Coming soon" -->
    <div v-if="comingSoonCards.length" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
      <div
        v-for="soon in comingSoonCards"
        :key="soon.labelKey"
        class="flex items-center gap-3 p-3.5 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-800/50 text-left opacity-70 cursor-not-allowed"
      >
        <div class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center" :class="soon.tint">
          <i :class="soon.icon" class="text-base"></i>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-700 dark:text-gray-300 truncate">{{ t(soon.labelKey) }}</p>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
            <Tag :value="t('common.coming_soon')" severity="secondary" class="!text-[10px] !px-1.5 !py-0" />
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import Tag from 'primevue/tag'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const { hasExactPermission } = useAuth()

const stats = ref({ upcoming: 0, ongoing: 0, completed: 0, courses: 0, participants: 0, providers: 0 })

// P0-FE: stat awal dihitung dari endpoint list existing (count sessions by status,
// total courses/participants/providers). Dashboard analitik penuh di P2-FE.
async function loadStats() {
  try {
    const [sessRes, courseRes, partRes, provRes] = await Promise.all([
      api.get('/api/v1/tenant/trainings/sessions', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 1 } }),
      api.get('/api/v1/tenant/trainings/participants', { params: { per_page: 1 } }),
      api.get('/api/v1/tenant/trainings/providers', { params: { per_page: 1 } })
    ])
    const sessions = sessRes.data?.data || []
    stats.value = {
      upcoming: sessions.filter(s => ['SCHEDULED', 'REGISTRATION_OPEN', 'DRAFT', 'FULL'].includes(s.status)).length,
      ongoing: sessions.filter(s => s.status === 'IN_PROGRESS').length,
      completed: sessions.filter(s => s.status === 'COMPLETED').length,
      courses: courseRes.data?.total || 0,
      participants: partRes.data?.total || 0,
      providers: provRes.data?.total || 0
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

// Semua card menu dicek dengan hasExactPermission (tanpa fallback module-level
// training.view) — sama seperti attendance/leave: "training.view" dimiliki
// hampir semua role (termasuk Employee default), sehingga tidak boleh otomatis
// mencakup submenu settings/operations/records/reports. Key permission
// mengikuti submenu RBAC yang sebenarnya (bukan key per-halaman).
const menuGroups = computed(() => {
  const groups = [
    {
      titleKey: 'training.group_settings',
      items: [
        { labelKey: 'training.courses', descKey: 'training.courses_desc', icon: 'pi pi-book', tint: 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400', route: '/training/courses', permission: 'training.settings.view' },
        { labelKey: 'training.categories', descKey: 'training.categories_desc', icon: 'pi pi-tags', tint: 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400', route: '/training/categories', permission: 'training.settings.view' },
        { labelKey: 'training.providers', descKey: 'training.providers_desc', icon: 'pi pi-building', tint: 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400', route: '/training/providers', permission: 'training.settings.view' },
        { labelKey: 'training.trainers', descKey: 'training.trainers_desc', icon: 'pi pi-user', tint: 'bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400', route: '/training/trainers', permission: 'training.settings.view' }
      ]
    },
    {
      titleKey: 'training.group_operations',
      items: [
        { labelKey: 'training.sessions', descKey: 'training.sessions_desc', icon: 'pi pi-calendar', tint: 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400', route: '/training/sessions', permission: 'training.operations.view' },
        { labelKey: 'training.participants', descKey: 'training.participants_desc', icon: 'pi pi-users', tint: 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400', route: '/training/participants', permission: 'training.operations.view' },
        { labelKey: 'training.planning', descKey: 'training.planning_desc', icon: 'pi pi-calendar-plus', tint: 'bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400', route: '/training/plans', permission: 'training.operations.view' },
        { labelKey: 'training.requests', descKey: 'training.requests_desc', icon: 'pi pi-send', tint: 'bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400', route: '/training/requests', permission: 'training.operations.view' },
        { labelKey: 'training.needs', descKey: 'training.needs_desc', icon: 'pi pi-bullseye', tint: 'bg-cyan-50 dark:bg-cyan-500/10 text-cyan-600 dark:text-cyan-400', route: '/training/needs', permission: 'training.operations.view' }
      ]
    },
    {
      titleKey: 'training.group_records',
      items: [
        { labelKey: 'training.certificates', descKey: 'training.certificates_desc', icon: 'pi pi-id-card', tint: 'bg-orange-50 dark:bg-orange-500/10 text-orange-600 dark:text-orange-400', route: '/training/certificates', permission: 'training.records.view' },
        { labelKey: 'training.history', descKey: 'training.history_desc', icon: 'pi pi-history', tint: 'bg-purple-50 dark:bg-purple-500/10 text-purple-600 dark:text-purple-400', route: '/training/history', permission: 'training.records.view' }
      ]
    },
    {
      titleKey: 'training.group_reports',
      items: [
        { labelKey: 'training.reports', descKey: 'training.reports_desc', icon: 'pi pi-chart-bar', tint: 'bg-blue-50 dark:bg-blue-500/10 text-blue-600 dark:text-blue-400', route: '/training/reports', permission: 'training.reports.view' }
      ]
    }
  ]
  return groups
    .map(g => ({ ...g, items: g.items.filter(card => !card.permission || hasExactPermission(card.permission)) }))
    .filter(g => g.items.length > 0)
})

// Semua fitur modul sudah memiliki halaman — tidak ada card coming-soon.
const comingSoonCards = computed(() => [])

onMounted(loadStats)
</script>
