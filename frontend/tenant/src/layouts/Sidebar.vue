<template>
  <aside
    class="flex flex-col bg-white border-r border-gray-200 transition-all duration-200 overflow-hidden"
    :class="collapsed ? 'w-16' : 'w-60'"
  >
    <!-- Logo Area -->
    <div class="flex items-center h-12 px-4 border-b border-gray-200 shrink-0">
      <i class="pi pi-building text-emerald-600 text-lg mr-2"></i>
      <span v-if="!collapsed" class="font-semibold text-sm text-gray-800 truncate">HRIS Platform</span>
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
        class="w-9 h-9 rounded-lg flex items-center justify-center cursor-pointer hover:bg-emerald-100 transition-colors"
        :class="{ 'bg-emerald-100 text-emerald-700': isActive(item) }"
        @click="item.command?.()"
      >
        <i :class="item.icon" class="text-sm"></i>
      </div>
    </nav>

    <!-- Bottom Section -->
    <div class="border-t border-gray-200 p-3 shrink-0">
      <div class="flex items-center gap-2 text-sm text-gray-500">
        <i class="pi pi-circle-fill text-emerald-400 text-[6px]"></i>
        <span class="truncate">Tenant: PT. ABC</span>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import PanelMenu from 'primevue/panelmenu'
import Tooltip from 'primevue/tooltip'
import { useRouter, useRoute } from 'vue-router'

defineProps({
  collapsed: { type: Boolean, default: false }
})

defineEmits(['toggle'])

const router = useRouter()
const route = useRoute()

// ── Full menu items for expanded PanelMenu ──
const menuItems = [
  {
    label: 'Dashboard',
    icon: 'pi pi-home',
    command: () => router.push('/dashboard'),
    class: route.name === 'Dashboard' ? 'bg-emerald-50 text-emerald-700 rounded-md' : ''
  },
  {
    label: 'Core HR',
    icon: 'pi pi-building',
    items: [
      { label: 'Organization', icon: 'pi pi-sitemap', command: () => router.push('/organizations') },
      { label: 'Employees', icon: 'pi pi-users', command: () => router.push('/employees') },
      { label: 'Job Management', icon: 'pi pi-briefcase', command: () => router.push('/job-management') }
    ]
  },
  {
    label: 'Talent',
    icon: 'pi pi-star',
    items: [
      { label: 'Competency', icon: 'pi pi-star', command: () => router.push('/competencies') },
      { label: 'Performance', icon: 'pi pi-chart-line', command: () => router.push('/performance') },
      { label: 'Training', icon: 'pi pi-book', command: () => router.push('/training') },
      { label: 'Recruitment', icon: 'pi pi-user-plus', command: () => router.push('/recruitment') }
    ]
  },
  {
    label: 'Operations',
    icon: 'pi pi-cog',
    items: [
      { label: 'Attendance', icon: 'pi pi-clock', command: () => router.push('/attendance') },
      { label: 'Leave', icon: 'pi pi-calendar', command: () => router.push('/leave') },
      { label: 'Movement', icon: 'pi pi-arrows-alt', command: () => router.push('/employee-movements') },
      { label: 'Approval', icon: 'pi pi-check-square', command: () => router.push('/approvals') }
    ]
  },
  {
    label: 'Finance',
    icon: 'pi pi-dollar',
    items: [
      { label: 'Payroll', icon: 'pi pi-dollar', command: () => router.push('/payroll') },
      { label: 'Reimbursement', icon: 'pi pi-credit-card', command: () => router.push('/reimbursements') }
    ]
  },
  {
    label: 'Strategic',
    icon: 'pi pi-chart-bar',
    items: [
      { label: 'Workforce Intel', icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence') },
      { label: 'Career Intel', icon: 'pi pi-road', command: () => router.push('/career-intelligence') }
    ]
  }
]

// ── Flatten top-level items for collapsed sidebar ──
// Groups like "Core HR" become a single shortcut to the first child route
const topLevelMenuItems = computed(() => {
  return [
    { key: 'Dashboard', label: 'Dashboard', path: '/dashboard', icon: 'pi pi-home', command: () => router.push('/dashboard') },
    { key: 'CoreHR', label: 'Core HR', path: '/organizations', icon: 'pi pi-building', command: () => router.push('/organizations') },
    { key: 'Talent', label: 'Talent', path: '/competencies', icon: 'pi pi-star', command: () => router.push('/competencies') },
    { key: 'Operations', label: 'Operations', path: '/attendance', icon: 'pi pi-cog', command: () => router.push('/attendance') },
    { key: 'Finance', label: 'Finance', path: '/payroll', icon: 'pi pi-dollar', command: () => router.push('/payroll') },
    { key: 'Strategic', label: 'Strategic', path: '/workforce-intelligence', icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence') }
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

:deep(.p-panelmenu-content .p-menuitem-icon) {
  font-size: 0.75rem !important;
}
</style>
