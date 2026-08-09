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
    <nav v-if="!collapsed" class="flex-1 overflow-y-auto py-2 px-2 space-y-0.5">
      <!-- Dashboard -->
      <router-link
        to="/dashboard"
        class="flex items-center gap-2 px-2.5 py-2 text-sm rounded-md transition-colors"
        :class="route.path === '/dashboard' ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-600 dark:text-gray-300 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
      >
        <i class="pi pi-home" style="font-size:0.875rem"></i>
        <span>{{ t('nav.dashboard') }}</span>
      </router-link>

      <!-- Core HR Section Title -->
      <div v-if="coreHRItems.length > 0" class="pt-2">
        <div class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
          <i class="pi pi-building" style="font-size:0.7rem"></i>
          <span>{{ t('nav.core_hr') }}</span>
        </div>
        <template v-for="item in coreHRItems" :key="item.key || item.label">
          <div
            class="ml-2 flex items-center justify-between gap-2 px-2.5 py-2 text-sm rounded-md cursor-pointer transition-colors"
            :class="isItemActive(item) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-600 dark:text-gray-300 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
            @click="item.children?.length ? toggleMenu(item.key) : item.command()"
          >
            <div class="flex items-center gap-2 min-w-0">
              <i :class="item.icon" class="text-xs"></i>
              <span class="truncate">{{ item.label }}</span>
            </div>
            <i
              v-if="item.children?.length"
              class="pi text-[11px] transition-transform duration-150"
              :class="isMenuOpen(item) ? 'pi-chevron-down' : 'pi-chevron-right'"
            ></i>
          </div>
          <!-- Dropdown children -->
          <template v-if="item.children?.length && isMenuOpen(item)">
            <div
              v-for="child in item.children"
              :key="child.key || child.label"
              class="ml-6 flex items-center gap-2 px-2.5 py-1.5 text-[13px] rounded-md cursor-pointer transition-colors"
              :class="isItemActive(child) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-500 dark:text-gray-400 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
              @click="child.command()"
            >
              <i :class="child.icon" class="text-[11px]"></i>
              <span>{{ child.label }}</span>
            </div>
          </template>
        </template>
      </div>

      <!-- Talent Section Title -->
      <div v-if="talentItems.length > 0" class="pt-2">
        <div class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
          <i class="pi pi-star" style="font-size:0.7rem"></i>
          <span>{{ t('nav.talent') }}</span>
        </div>
        <template v-for="item in talentItems" :key="item.key || item.label">
          <div
            class="ml-2 flex items-center justify-between gap-2 px-2.5 py-2 text-sm rounded-md cursor-pointer transition-colors"
            :class="isItemActive(item) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-600 dark:text-gray-300 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
            @click="item.children?.length ? toggleMenu(item.key) : item.command()"
          >
            <div class="flex items-center gap-2 min-w-0">
              <i :class="item.icon" class="text-xs"></i>
              <span class="truncate">{{ item.label }}</span>
            </div>
            <i
              v-if="item.children?.length"
              class="pi text-[11px] transition-transform duration-150"
              :class="isMenuOpen(item) ? 'pi-chevron-down' : 'pi-chevron-right'"
            ></i>
          </div>
          <!-- Dropdown children -->
          <template v-if="item.children?.length && isMenuOpen(item)">
            <div
              v-for="child in item.children"
              :key="child.key || child.label"
              class="ml-6 flex items-center gap-2 px-2.5 py-1.5 text-[13px] rounded-md cursor-pointer transition-colors"
              :class="isItemActive(child) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-500 dark:text-gray-400 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
              @click="child.command()"
            >
              <i :class="child.icon" class="text-[11px]"></i>
              <span>{{ child.label }}</span>
            </div>
          </template>
        </template>
      </div>

      <!-- Operations Section Title -->
      <div v-if="operationsItems.length > 0" class="pt-2">
        <div class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
          <i class="pi pi-cog" style="font-size:0.7rem"></i>
          <span>{{ t('nav.operations') }}</span>
        </div>
        <div
          v-for="item in operationsItems"
          :key="item.key || item.label"
          class="ml-2 flex items-center gap-2 px-2.5 py-2 text-sm rounded-md cursor-pointer transition-colors"
          :class="isItemActive(item) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-600 dark:text-gray-300 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
          @click="item.command()"
        >
          <i :class="item.icon" class="text-xs"></i>
          <span>{{ item.label }}</span>
        </div>
      </div>

      <!-- Settings Section Title — masuk ke halaman index (card sub-menu) -->
      <div v-if="settingsItems.length > 0" class="pt-2">
        <div class="flex items-center gap-1.5 px-2.5 py-1.5 text-[11px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500">
          <i class="pi pi-cog" style="font-size:0.7rem"></i>
          <span>{{ t('nav.settings') }}</span>
        </div>
        <div
          v-for="item in settingsItems"
          :key="item.key || item.label"
          class="ml-2 flex items-center gap-2 px-2.5 py-2 text-sm rounded-md cursor-pointer transition-colors"
          :class="isItemActive(item) ? 'bg-emerald-50 dark:bg-emerald-900/30 text-emerald-700 dark:text-emerald-300' : 'text-gray-600 dark:text-gray-300 hover:bg-emerald-50 dark:hover:bg-emerald-900/10'"
          @click="item.command()"
        >
          <i :class="item.icon" class="text-xs"></i>
          <span>{{ item.label }}</span>
        </div>
      </div>

      <!-- Other groups as PanelMenu -->
      <div class="pt-1">
        <PanelMenu :model="panelMenuItems" class="border-none !bg-transparent" />
      </div>
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

    <!-- Bottom Section: Company Name → klik buka halaman Detail Perusahaan -->
    <div class="border-t border-gray-200 dark:border-gray-700 p-3 shrink-0">
      <button
        type="button"
        class="flex items-center gap-2 w-full text-sm text-gray-500 dark:text-gray-400 hover:text-emerald-600 dark:hover:text-emerald-400 transition-colors rounded-md px-1 py-0.5"
        :title="companyLabel"
        @click="router.push('/company')"
      >
        <i class="pi pi-building text-emerald-400 text-xs"></i>
        <span class="truncate">{{ companyLabel }}</span>
      </button>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import PanelMenu from 'primevue/panelmenu'
import Tooltip from 'primevue/tooltip'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useActiveModules } from '@/stores/activeModules'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'

defineProps({
  collapsed: { type: Boolean, default: false }
})

defineEmits(['toggle'])

const router = useRouter()
const route = useRoute()
const { state: authState, hasPermission } = useAuth()
const activeMod = useActiveModules()
const { t } = useI18n()
const companyName = ref('')

// Fetch enriched user data (includes company_name).
// Catatan: endpoint /api/v1/platform/users hanya berlaku utk platform user
// (company_admin/super_admin) — utk login employee, nama company sudah
// tersedia dari response login (user.company_name) dan disimpan di
// authState.company.name, jadi fetch platform di-skip bila sudah terisi.
onMounted(async () => {
  if (!authState.user?.id) return
  if (!authState.company?.name) {
    try {
      const res = await api.get(`/api/v1/platform/users/${authState.user.id}`)
      const data = res.data?.data || res.data
      companyName.value = data?.company_name || ''
    } catch {
      // Silently fail — fallback to login data
    }
  }

  // Fetch active modules
  await activeMod.fetchActiveModules()
})

// Company label for bottom section — prioritas: state company (sync dari
// X-Tenant-ID header) → enriched user data → fallback.
const companyLabel = computed(() => {
  if (authState.company?.name) return authState.company.name
  if (companyName.value) return companyName.value
  if (authState.company?.id) return 'Company'
  if (authState.user?.company_id) return 'Company'
  return 'HRIS Platform'
})

// ── Helper: filter items based on active modules (Level 1) + permission (Level 2) ──
function filterByModule(items) {
  return items.filter(item => {
    if (item.moduleSlug && !activeMod.hasModule(item.moduleSlug)) return false
    if (item.permission && !hasPermission(item.permission)) return false
    return true
  })
}

// ── Dropdown state untuk item menu yang punya children (mis. Job Management) ──
const openMenus = ref({})

function toggleMenu(key) {
  // Jika belum pernah di-toggle manual, state saat ini didapat dari auto-open
  // (child aktif). Toggle pertama harus membalik dari state efektif itu, bukan
  // dari `undefined`, supaya menu dengan child aktif tetap bisa ditutup.
  if (openMenus.value[key] === undefined) {
    const item = allDropdownItems.value.find(i => i.key === key)
    const autoOpen = !!item?.children?.some(c => isItemActive(c))
    openMenus.value[key] = !autoOpen
    return
  }
  openMenus.value[key] = !openMenus.value[key]
}

function isMenuOpen(item) {
  if (!item?.children?.length) return false
  // Jika user sudah pernah toggle manual, state itu yang menang.
  if (openMenus.value[item.key] !== undefined) return openMenus.value[item.key]
  // Auto-open jika salah satu child sedang aktif (mis. sedang di /job-management/values)
  return item.children.some(c => isItemActive(c))
}

// ── Core HR items (flat, not dropdown) ──
const coreHRItems = computed(() => {
  return filterByModule([
    { key: 'organization', label: t('nav.organization'), icon: 'pi pi-sitemap', command: () => router.push('/organization-summary'), path: '/organization-summary', moduleSlug: 'organization', permission: 'organization.view' },
    { key: 'employees', label: t('nav.employees'), icon: 'pi pi-users', command: () => router.push('/employees'), path: '/employees', moduleSlug: 'employee', permission: 'employee.view' },
    {
      key: 'job_management',
      label: t('nav.job_management'),
      icon: 'pi pi-briefcase',
      command: () => router.push('/job-management'),
      moduleSlug: 'jobmanagement',
      permission: 'jobmanagement.view',
      children: [
        {
          key: 'job_values_mapping',
          label: t('nav.job_values_mapping'),
          icon: 'pi pi-sliders-h',
          command: () => router.push('/job-management/values'),
          path: '/job-management/values',
          moduleSlug: 'jobmanagement',
          permission: 'jobmanagement.view'
        },
        {
          key: 'job_management_main',
          label: t('nav.job_management'),
          icon: 'pi pi-briefcase',
          command: () => router.push('/job-management'),
          path: '/job-management',
          // Jangan ter-highlight saat berada di sub-halaman Job Value Mapping
          excludePaths: ['/job-management/values'],
          moduleSlug: 'jobmanagement',
          permission: 'jobmanagement.view'
        }
      ]
    }
  ])
})

// ── Talent items (with dropdown support) ──
const talentItems = computed(() => {
  return filterByModule([
    { key: 'competency', label: t('nav.competency'), icon: 'pi pi-star', command: () => router.push('/competencies'), path: '/competencies', moduleSlug: 'competency', permission: 'competency.view' },
    {
      key: 'performance',
      label: t('nav.performance'),
      icon: 'pi pi-chart-line',
      command: () => router.push('/performance'),
      moduleSlug: 'performance',
      permission: 'performance.view',
      children: [
        {
          key: 'performance_main',
          label: t('nav.performance_dashboard'),
          icon: 'pi pi-th-large',
          command: () => router.push('/performance'),
          path: '/performance',
          excludePaths: ['/performance/kpi', '/performance/okr'],
          moduleSlug: 'performance',
          permission: 'performance.view'
        },
        {
          key: 'performance_kpi',
          label: t('nav.kpi'),
          icon: 'pi pi-chart-bar',
          command: () => router.push('/performance/kpi'),
          path: '/performance/kpi',
          moduleSlug: 'performance',
          permission: 'performance.view'
        },
        {
          key: 'performance_okr',
          label: t('nav.okr'),
          icon: 'pi pi-bullseye',
          command: () => router.push('/performance/okr'),
          path: '/performance/okr',
          moduleSlug: 'performance',
          permission: 'performance.view'
        }
      ]
    },
    { key: 'training', label: t('nav.training'), icon: 'pi pi-book', command: () => router.push('/training'), path: '/training', moduleSlug: 'training', permission: 'training.view' },
    { key: 'recruitment', label: t('nav.recruitment'), icon: 'pi pi-user-plus', command: () => router.push('/recruitment'), path: '/recruitment', moduleSlug: 'recruitment', permission: 'recruitment.view' }
  ])
})

// ── Operations items (flat, not dropdown) ──
const operationsItems = computed(() => {
  return filterByModule([
    { key: 'attendance', label: t('nav.attendance'), icon: 'pi pi-clock', command: () => router.push('/attendance'), path: '/attendance', moduleSlug: 'attendance', permission: 'attendance.view' },
    { key: 'leave', label: t('nav.leave'), icon: 'pi pi-calendar', command: () => router.push('/leave'), path: '/leave', moduleSlug: 'leave', permission: 'leave.view' },
    { key: 'movement', label: t('employee_movement.movements'), icon: 'pi pi-arrows-alt', command: () => router.push('/admin/career/movements'), path: '/admin/career/movements', moduleSlug: 'employeemovement', permission: 'employeemovement.view' },
    { key: 'contracts', label: t('employee_movement.contracts'), icon: 'pi pi-file-edit', command: () => router.push('/admin/career/contracts'), path: '/admin/career/contracts', moduleSlug: 'employeemovement', permission: 'employeemovement.view' },
    { key: 'approval', label: t('nav.approval'), icon: 'pi pi-check-square', command: () => router.push('/approvals'), path: '/approvals', moduleSlug: 'approval', permission: 'approval.view' },
    { key: 'notification', label: t('nav.notification'), icon: 'pi pi-bell', command: () => router.push('/notifications'), path: '/notifications', moduleSlug: 'notification', permission: 'notification.view' }
  ])
})

// ── Gabungan semua item dropdown (untuk lookup by key di toggleMenu) ──
const allDropdownItems = computed(() => [...coreHRItems.value, ...talentItems.value])

// ── Settings items (flat, masuk ke halaman index card) ──
  const settingsItems = computed(() => {
    return filterByModule([
      { key: 'settings', label: t('nav.settings'), icon: 'pi pi-cog', command: () => router.push('/settings'), path: '/settings', moduleSlug: 'setting', permission: 'setting.view' }
    ])
  })

  // ── Panel menu groups for other sections ──
  const panelGroups = computed(() => [
  // Finance
  {
    key: 'finance',
    label: t('nav.finance'),
    icon: 'pi pi-dollar',
    items: [
      { label: t('nav.payroll'), icon: 'pi pi-dollar', command: () => router.push('/payroll'), moduleSlug: 'payroll', permission: 'payroll.view' },
      { label: t('nav.reimbursement'), icon: 'pi pi-credit-card', command: () => router.push('/reimbursements'), moduleSlug: 'reimbursement', permission: 'reimbursement.view' }
    ]
  },
  // Strategic
  {
    key: 'strategic',
    label: t('nav.strategic'),
    icon: 'pi pi-chart-bar',
    items: [
      { label: t('nav.workforce_intel'), icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence'), moduleSlug: 'workforce-intelligence', permission: 'workforceintelligence.view' },
      { label: t('nav.career_intel'), icon: 'pi pi-chart-bar', command: () => router.push('/career-intelligence'), moduleSlug: 'career-intelligence', permission: 'careerintelligence.view' }
    ]
  }
])  // ── Filtered PanelMenu items — only show visible groups (Settings sudah flat section) ──
  const panelMenuItems = computed(() => {
    return panelGroups.value
      .map(group => {
        const visibleItems = filterByModule(group.items)
        if (visibleItems.length === 0) return null
        return { ...group, items: visibleItems }
      })
      .filter(Boolean)
  })

// ── Flatten top-level items for collapsed sidebar ──
const topLevelMenuItems = computed(() => {
  const items = []
  // Dashboard — always visible
  items.push({ key: 'Dashboard', label: t('nav.dashboard'), path: '/dashboard', icon: 'pi pi-home', command: () => router.push('/dashboard') })
  
  // Core HR items
  coreHRItems.value.forEach(item => {
    // Item yang punya children adalah pure toggle di expanded mode — di collapsed
    // cukup push children (navigasi parent sudah diwakili child, hindari ikon duplikat).
    if (item.children?.length) {
      item.children.forEach(child => {
        items.push({ ...child, key: 'CoreHR-' + (child.key || child.label) })
      })
    } else {
      items.push({ ...item, key: 'CoreHR-' + item.key })
    }
  })
  
  // Talent items
  talentItems.value.forEach(item => {
    items.push({ ...item, key: 'Talent-' + item.key })
  })

  // Operations items
  operationsItems.value.forEach(item => {
    items.push({ ...item, key: 'Operations-' + item.key })
  })

  // Finance — hanya jika ada modul aktif DAN user punya permission view
  if ((activeMod.hasModule('payroll') && hasPermission('payroll.view')) || (activeMod.hasModule('reimbursement') && hasPermission('reimbursement.view'))) {
    items.push({ key: 'Finance', label: t('nav.finance'), path: '/payroll', icon: 'pi pi-dollar', command: () => router.push('/payroll') })
  }
  // Strategic
  if ((activeMod.hasModule('workforce-intelligence') && hasPermission('workforceintelligence.view')) || (activeMod.hasModule('career-intelligence') && hasPermission('careerintelligence.view'))) {
    items.push({ key: 'Strategic', label: t('nav.strategic'), path: '/workforce-intelligence', icon: 'pi pi-chart-bar', command: () => router.push('/workforce-intelligence') })
  }
  // Settings
  if (activeMod.hasModule('setting') && hasPermission('setting.view')) {
    items.push({ key: 'Settings', label: t('nav.settings'), path: '/settings', icon: 'pi pi-cog', command: () => router.push('/settings') })
  }
  
  return items
})

// ── Active state check ──
// excludePaths: daftar path yang TIDAK boleh memicu highlight item ini
// (mis. child 'Job Management' tidak boleh aktif saat di /job-management/values).
function isExcludedPath(item) {
  if (!item.excludePaths?.length) return false
  return item.excludePaths.some(p => route.path === p || route.path.startsWith(p + '/') || route.path.startsWith(p + '?'))
}

function isActive(item) {
  if (!item.path) return false
  if (isExcludedPath(item)) return false
  return route.path.startsWith(item.path)
}

function isItemActive(item) {
  if (!item.path) return false
  if (isExcludedPath(item)) return false
  // Exact match or starts with path
  if (route.path === item.path) return true
  return route.path.startsWith(item.path + '/') || route.path.startsWith(item.path + '?')
}
</script>

<style scoped>
:deep(.p-panelmenu){
  gap: 0px !important;
}
/* :deep(.p-panelmenu-header-label) {
  font-size: 1.2em !important;
} */
:deep(.p-panelmenu-panel) {
  background: transparent !important;
  border: none !important;
  margin-bottom: 1px;
}
:deep(.p-panelmenu-item-label) {
  font-size: 0.8125rem !important;
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
