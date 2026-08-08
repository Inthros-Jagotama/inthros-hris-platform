<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
    <button
      v-for="item in items"
      :key="item.path"
      type="button"
      class="cursor-pointer group flex items-center gap-3 p-3.5 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-left transition-all hover:border-emerald-300 dark:hover:border-emerald-500/60 hover:shadow-md hover:-translate-y-0.5 focus:outline-none focus-visible:ring-2 focus-visible:ring-emerald-500/50"
      @click="router.push(item.path)"
    >
      <div class="w-10 h-10 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
        <i :class="item.icon" class="text-base"></i>
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ t(item.titleKey) }}</p>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">{{ t(item.descKey) }}</p>
      </div>
      <i class="pi pi-chevron-right text-xs text-gray-300 dark:text-gray-600 group-hover:text-emerald-400 group-hover:translate-x-0.5 transition-all shrink-0"></i>
    </button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'

const router = useRouter()
const { t } = useI18n()

const items = [
  { path: '/attendance/settings', icon: 'pi pi-cog', titleKey: 'attendance.settings', descKey: 'attendance.settings_description' },
  { path: '/attendance/shifts', icon: 'pi pi-clock', titleKey: 'attendance.shifts', descKey: 'attendance.shifts_description' },
  { path: '/attendance/employee-shifts', icon: 'pi pi-users', titleKey: 'attendance.employee_shifts', descKey: 'attendance.employee_shifts_description' },
  { path: '/attendance/locations', icon: 'pi pi-map-marker', titleKey: 'attendance.locations', descKey: 'attendance.locations_description' },
  { path: '/attendance/exempt-positions', icon: 'pi pi-shield', titleKey: 'attendance.exempt_positions', descKey: 'attendance.exempt_positions_description' }
]
</script>
