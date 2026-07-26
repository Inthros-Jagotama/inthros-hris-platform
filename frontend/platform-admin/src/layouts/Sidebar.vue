<template>
  <aside
    class="flex flex-col bg-white border-r border-gray-200 transition-all duration-200 overflow-hidden"
    :class="collapsed ? 'w-16' : 'w-56'"
  >
    <!-- Logo -->
    <div class="flex items-center h-12 px-4 border-b border-gray-200 shrink-0 gap-2">
      <i class="pi pi-shield text-indigo-600 text-lg"></i>
      <span v-if="!collapsed" class="font-semibold text-sm text-gray-800 truncate">Platform Admin</span>
    </div>

    <!-- Nav: Expanded -->
    <nav v-if="!collapsed" class="flex-1 overflow-y-auto py-3 px-2 space-y-1">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        class="flex items-center gap-2.5 px-2.5 py-1.5 rounded-md text-sm transition-colors"
        :class="isActive(item.path) ? 'bg-indigo-50 text-indigo-700 font-medium' : 'text-gray-600 hover:bg-gray-100'"
      >
        <i :class="item.icon" class="text-sm w-4 text-center"></i>
        <span>{{ item.label }}</span>
      </router-link>
    </nav>

    <!-- Nav: Collapsed (icon-only) -->
    <nav v-else class="flex-1 overflow-y-auto py-3 px-1 flex flex-col items-center gap-1">
      <router-link
        v-for="item in navItems"
        :key="item.path"
        :to="item.path"
        v-tooltip.left="item.label"
        class="w-9 h-9 rounded-lg flex items-center justify-center transition-colors"
        :class="isActive(item.path) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-500 hover:bg-gray-100'"
      >
        <i :class="item.icon" class="text-sm"></i>
      </router-link>
    </nav>
  </aside>
</template>

<script setup>
import { useRoute } from 'vue-router'

defineProps({
  collapsed: { type: Boolean, default: false }
})

defineEmits(['toggle'])

const route = useRoute()

const navItems = [
  { label: 'Dashboard', path: '/dashboard', icon: 'pi pi-home' },
  { label: 'Companies', path: '/companies', icon: 'pi pi-building' },
  { label: 'Users', path: '/users', icon: 'pi pi-users' },
  { label: 'Modules', path: '/modules', icon: 'pi pi-cog' },
  { label: 'Licenses', path: '/licenses', icon: 'pi pi-id-card' },
  { label: 'Monitoring', path: '/monitoring', icon: 'pi pi-chart-bar' }
]

function isActive(path) {
  return route.path.startsWith(path)
}
</script>
