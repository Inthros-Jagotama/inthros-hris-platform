<template>
  <header class="flex items-center h-12 bg-white border-b border-gray-200 px-4 shrink-0">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <Button
        icon="pi pi-bars"
        severity="secondary"
        text
        size="small"
        @click="$emit('toggle-sidebar')"
        class="!p-1.5"
      />
      <i class="pi pi-chevron-right text-sm text-gray-300"></i>
      <span class="text-sm text-gray-500 font-medium truncate">{{ route.meta?.title || 'Dashboard' }}</span>
    </div>

    <div class="flex items-center gap-2">
      <div class="flex items-center gap-2 text-sm text-gray-500 mr-2">
        <i class="pi pi-circle-fill text-emerald-400 text-[6px]"></i>
        <span>Live</span>
      </div>

      <Button
        severity="secondary"
        text
        size="small"
        class="!p-1"
        @click="menu.toggle($event)"
      >
        <div class="flex items-center gap-2">
          <Avatar
            icon="pi pi-user"
            size="small"
            class="!w-7 !h-7 !bg-indigo-100 !text-indigo-700"
          />
          <span class="text-sm text-gray-700 hidden sm:inline">Admin</span>
          <i class="pi pi-chevron-down text-sm text-gray-400"></i>
        </div>
      </Button>
      <Menu ref="menu" :model="menuItems" popup />
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'

const emit = defineEmits(['toggle-sidebar', 'logout'])

const route = useRoute()
const menu = ref(null)

const menuItems = [
  { label: 'Profile', icon: 'pi pi-user', command: () => {} },
  { separator: true },
  { label: 'Logout', icon: 'pi pi-sign-out', command: () => emit('logout') }
]
</script>
