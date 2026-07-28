<template>
  <aside
    class="flex flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transition-all duration-200 overflow-hidden"
    :class="collapsed ? 'w-16' : 'w-60'"
  >
    <!-- Logo Area -->
    <div class="flex items-center h-12 px-4 border-b border-gray-200 dark:border-gray-700 shrink-0">
      <i class="pi pi-building text-emerald-600 text-lg mr-2"></i>
      <span v-if="!collapsed" class="font-semibold text-sm text-gray-800 dark:text-gray-100 truncate">HRIS Platform</span>
    </div>

    <!-- Navigation: Expanded -->
    <nav v-if="!collapsed" class="flex-1 overflow-y-auto py-2 px-2">
      <PanelMenu :model="menuItems" class="border-none !bg-transparent" />
    </nav>

    <!-- Navigation: Collapsed (icon-only) -->
    <nav v-else class="flex-1 overflow-y-auto py-3 px-1 flex flex-col items-center gap-1">
      <div
        v-for="item in topLevelMenuItems"
        :key="item.key || item.label"
        v-tooltip.left="item.tooltip || item.label"
        class="w-9 h-9 rounded-lg flex items-center justify-center cursor-pointer hover:bg-emerald-100 dark:hover:bg-emerald-900/30 transition-colors"
        :class="{ 'bg-emerald-100 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300': isActive(item) }"
        @click="item.command?.()"
      >
        <i :class="item.icon" class="text-sm"></i>
      </div>
    </nav>

    <!-- Bottom Section: Company Name -->
    <div class="border-t border-gray-200 dark:border-gray-700 p-3 shrink-0">
      <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
        <i class="pi pi-circle-fill text-emerald-400 text-[6px]"></i>
        <span class="truncate" :title="companyLabel">{{ companyLabel }}</span>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import PanelMenu from 'primevue/panelmenu'
import Tooltip from 'primevue/tooltip'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'

defineProps({
  collapsed: { type: Boolean, default: false }
})

defineEmits(['toggle'])

const router = useRouter()
const route = useRoute()
const { state: authState } = useAuth()
const { t } = useI18n()
const companyName = ref('')

// Fetch enriched user data (includes company_name)
onMounted(async () => {
  if (!authState.user?.id) return
  try {
    const res = await api.get(`/api/v1/platform/users/${authState.user.id}`)
    const data = res.data?.data || res.data
    companyName.value = data?.company_name || ''
  } catch {
    // Silently fail — fallback to login data
  }
})

// Company label for bottom section
const companyLabel = computed(() => {
  if (companyName.value) return companyName.value
  if (authState.user?.company_id) return 'Company'
  return 'HRIS Platform'
})

// ── Full menu items for expanded PanelMenu ──
const menuItems = computed(() => [
  {
    label: t('nav.dashboard'),
    icon: 'pi pi-home',
    command: () => router.push('/dashboard'),
    class: route.name === 'Dashboard' ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300 rounded-md' : ''
  },
  {
    label: t('nav.core_hr'),
    icon: 'pi pi-building',
    items: [
      { label: t('nav.organization'), icon: 'pi pi-sitemap', command: () => router.push('/organizations') },
      { label: t('nav.employees'), icon: 'pi pi-users', command: () => router.push('/employees') },
      { label: t('nav.job_management'), icon: 'pi pi-briefcase', command: () => router.push('/job-management') }
    ]
  },
  {
    label: t('nav.talent'),
    icon: 'pi pi-star',
    items: [
      { label: t('nav.competency'), icon: 'pi pi-star', command: () => router.push('/competencies') },
      { label: t('nav.performance'), icon: 'pi pi-chart-line', command: () => router.push('/performance') },
      { label: t('nav.training'), icon: 'pi pi-book', command: () => router.push('/training') },
      { label: t('nav.recruitment'), icon: 'pi pi-user-plus', command: () => router.push('/recruitment') }
    ]
  },
  {
    label: t('nav.operations'),
    icon: 'pi pi-cog',
    items: [
      { label: t('nav.attendance'), icon: 'pi pi-clock', command: () => router.push('/attendance') },
      { label: t('nav.leave'), icon: 'pi pi-calendar', command: () => router.push('/leave') },
      { label: t('nav.movement'), icon: 'pi pi-arrows-alt', command: () => router.push('/employee-movements') },
      { label: t('nav.approval'), icon: 'pi pi-check-square', command: () => router.push('/approvals') }
    ]
  },
  {
    label: t('nav.finance'),
    icon: 'pi pi-dollar',
    items: [
      { label: t('nav.payroll'), icon: 'pi pi-dollar', command: () => router.push('/payroll') },
      { label: t('nav.reimbursement'), icon: 'pi pi-credit-card', command: () => router.push('/reimbursements') }
    ]
  },
  {
    label: t('nav.strategic'),
    icon: 'pi pi-chart-bar',
    items: [
      { label: t('nav.workforce_intel'), icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence') },
      { label: t('nav.career_intel'), icon: 'pi pi-road', command: () => router.push('/career-intelligence') }
    ]
  },
  {
    label: t('nav.settings'),
    icon: 'pi pi-cog',
    items: [
      { label: t('settings.zones'), icon: 'pi pi-map-marker', command: () => router.push('/settings/zones') },
      { label: t('settings.provinces'), icon: 'pi pi-globe', command: () => router.push('/settings/provinces') },
      { label: t('settings.regencies'), icon: 'pi pi-map', command: () => router.push('/settings/regencies') },
      { label: t('settings.districts'), icon: 'pi pi-building', command: () => router.push('/settings/districts') },
      { label: t('settings.villages'), icon: 'pi pi-home', command: () => router.push('/settings/villages') },
      { label: t('settings.educations'), icon: 'pi pi-graduation-cap', command: () => router.push('/settings/educations') },
      { label: t('settings.religions'), icon: 'pi pi-globe', command: () => router.push('/settings/religions') },
      { label: t('settings.marital_statuses'), icon: 'pi pi-heart', command: () => router.push('/settings/marital-statuses') },
      { label: t('settings.relationship_types'), icon: 'pi pi-users', command: () => router.push('/settings/relationship-types') },
      { label: t('settings.banks'), icon: 'pi pi-building-column', command: () => router.push('/settings/banks') },
      { label: t('settings.employment_statuses'), icon: 'pi pi-briefcase', command: () => router.push('/settings/employment-statuses') },
      { label: t('settings.nationalities'), icon: 'pi pi-globe', command: () => router.push('/settings/nationalities') },
      { label: t('settings.job_families'), icon: 'pi pi-briefcase', command: () => router.push('/settings/job-families') },
      { label: t('settings.salary_grades'), icon: 'pi pi-chart-bar', command: () => router.push('/settings/salary-grades') }
    ]
  }
])

// ── Flatten top-level items for collapsed sidebar ──
// Groups like "Core HR" become a single shortcut to the first child route
const topLevelMenuItems = computed(() => {
  return [
    { key: 'Dashboard', label: t('nav.dashboard'), path: '/dashboard', icon: 'pi pi-home', command: () => router.push('/dashboard') },
    { key: 'CoreHR', label: t('nav.core_hr'), path: '/organizations', icon: 'pi pi-building', command: () => router.push('/organizations') },
    { key: 'Talent', label: t('nav.talent'), path: '/competencies', icon: 'pi pi-star', command: () => router.push('/competencies') },
    { key: 'Operations', label: t('nav.operations'), path: '/attendance', icon: 'pi pi-cog', command: () => router.push('/attendance') },
    { key: 'Finance', label: t('nav.finance'), path: '/payroll', icon: 'pi pi-dollar', command: () => router.push('/payroll') },
    { key: 'Strategic', label: t('nav.strategic'), path: '/workforce-intelligence', icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence') },
    { key: 'Settings', label: t('nav.settings'), path: '/settings/zones', icon: 'pi pi-cog', command: () => router.push('/settings/zones') }
  ]
})

// ── Active state check ──
function isActive(item) {
  if (!item.path) return false
  return route.path.startsWith(item.path)
}
</script>

<style scoped>
:deep(.p-panelmenu-panel) {
  background: transparent !important;
  border: none !important;
  margin-bottom: 1px;
}

:deep(.p-panelmenu-header-link) {
  padding: 0.5rem 0.625rem !important;
  font-size: 0.8125rem !important;
  border-radius: 0.375rem !important;
  transition: all 0.15s ease;
}

:deep(.p-panelmenu-header-link:hover) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-panelmenu-header-link:hover) {
  background: rgba(16, 185, 129, 0.1) !important;
}

:deep(.p-panelmenu-header-link .p-menuitem-text) {
  font-weight: 500;
}

:deep(.p-panelmenu-content) {
  background: transparent !important;
  border: none !important;
  padding: 0 !important;
}

:deep(.p-panelmenu-content .p-menuitem-link) {
  padding: 0.375rem 0.5rem 0.375rem 2rem !important;
  font-size: 0.8125rem !important;
  border-radius: 0.25rem !important;
}

:deep(.p-panelmenu-content .p-menuitem-link:hover) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-panelmenu-content .p-menuitem-link:hover) {
  background: rgba(16, 185, 129, 0.1) !important;
}

:deep(.p-panelmenu-content .p-menuitem-icon) {
  font-size: 0.75rem !important;
}
</style>
