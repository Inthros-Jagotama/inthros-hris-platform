<template>
  <aside
    class="flex flex-col bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 transition-all duration-200 overflow-hidden"
    :class="collapsed ? 'w-16' : 'w-56'"
  >
    <!-- Logo -->
    <div class="flex items-center h-12 px-4 border-b border-gray-200 dark:border-gray-700 shrink-0 gap-2">
      <i class="pi pi-shield text-indigo-600 text-lg"></i>
      <span v-if="!collapsed" class="font-semibold text-sm text-gray-800 dark:text-gray-100 truncate">{{ t('nav.platform_admin') }}</span>
    </div>

    <!-- Nav: Expanded -->
    <nav v-if="!collapsed" class="flex-1 overflow-y-auto py-3 px-2 space-y-1">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-sm transition-colors"
        :class="isActive(item.path) ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300 font-medium' : 'text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'"
      >
        <i :class="item.icon" class="text-sm w-4 text-center"></i>
        <span>{{ t(item.labelKey) }}</span>
      </router-link>
    </nav>

    <!-- Nav: Collapsed (icon-only) -->
    <nav v-else class="flex-1 overflow-y-auto py-3 px-1 flex flex-col items-center gap-1">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        v-tooltip.left="t(item.labelKey)"
        class="w-9 h-9 rounded-lg flex items-center justify-center transition-colors"
        :class="isActive(item.path) ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 'text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700'"
      >
        <i :class="item.icon" class="text-sm"></i>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useI18n } from '@/composables/useI18n'

defineProps({
  collapsed: { type: Boolean, default: false }
})

defineEmits(['toggle'])

const route = useRoute()
const { t } = useI18n()

const navItems = [
  { labelKey: 'nav.dashboard', path: '/dashboard', icon: 'pi pi-home' },
  { labelKey: 'nav.companies', path: '/companies', icon: 'pi pi-building' },
  { labelKey: 'nav.users', path: '/users', icon: 'pi pi-users' },
  { labelKey: 'nav.modules', path: '/modules', icon: 'pi pi-cog' },
  { labelKey: 'nav.licenses', path: '/licenses', icon: 'pi pi-id-card' },
  { labelKey: 'nav.packages', path: '/packages', icon: 'pi pi-box' },
  { labelKey: 'nav.rbac', path: '/rbac', icon: 'pi pi-shield' },
  { labelKey: 'nav.monitoring', path: '/monitoring', icon: 'pi pi-chart-bar' }
]

function isActive(path) {
  return route.path.startsWith(path)
}
</script>
