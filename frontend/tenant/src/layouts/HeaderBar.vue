<template>
  <header class="flex items-center h-12 bg-white border-b border-gray-200 px-4 shrink-0">
    <!-- Left: Toggle + Breadcrumb -->
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <Button
        icon="pi pi-bars"
        severity="secondary"
        text
        size="small"
        @click="$emit('toggle-sidebar')"
        class="!p-1.5"
      />
      <span class="text-sm text-gray-500 font-medium truncate">
        {{ route.meta?.title || 'Dashboard' }}
      </span>
    </div>

    <!-- Right: Actions -->
    <div class="flex items-center gap-2">
      <!-- Search -->
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText
          v-model="searchQuery"
          placeholder="Search modules..."
          class="!w-48 !h-8 !text-sm"
          @keyup.enter="handleSearch"
        />
      </IconField>

      <!-- Notifications -->
      <div class="relative">
        <Button
          icon="pi pi-bell"
          severity="secondary"
          text
          size="small"
          class="!p-1.5"
        />
        <Badge value="3" severity="danger" class="!absolute -top-0.5 -right-0.5 !text-xs !min-w-[1.1rem] !h-[1.1rem]" />
      </div>

      <!-- User Menu -->
      <Button
        severity="secondary"
        text
        size="small"
        class="!p-1"
        @click="userMenuVisible = !userMenuVisible"
      >
        <div class="flex items-center gap-2">
          <Avatar
            icon="pi pi-user"
            size="small"
            class="!w-7 !h-7 !bg-emerald-100 !text-emerald-700"
          />
          <span class="text-sm text-gray-700 hidden sm:inline">Admin</span>
          <i class="pi pi-chevron-down text-sm text-gray-400"></i>
        </div>
      </Button>
      <Menu ref="userMenu" :model="userMenuItems" popup />
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'
import Badge from 'primevue/badge'

defineEmits(['toggle-sidebar'])

const route = useRoute()
const router = useRouter()
const searchQuery = ref('')
const userMenu = ref(null)
const userMenuVisible = ref(false)

const userMenuItems = [
  { label: 'Profile', icon: 'pi pi-user', command: () => {} },
  { label: 'Settings', icon: 'pi pi-cog', command: () => {} },
  { separator: true },
  { label: 'Logout', icon: 'pi pi-sign-out', command: () => {} }
]

function handleSearch() {
  if (!searchQuery.value.trim()) return
  const q = searchQuery.value.toLowerCase()
  const routes = router.getRoutes()
  const found = routes.find(r =>
    r.meta?.title?.toLowerCase().includes(q) ||
    r.path?.toLowerCase().includes(q)
  )
  if (found) {
    router.push(found.path)
    searchQuery.value = ''
  }
}
</script>
