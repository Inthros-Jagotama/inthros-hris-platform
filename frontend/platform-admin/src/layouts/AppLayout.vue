<template>
  <div class="flex h-screen overflow-hidden bg-gray-50">
    <Sidebar
      :collapsed="sidebarCollapsed"
      @toggle="sidebarCollapsed = !sidebarCollapsed"
    />
    <div class="flex flex-col flex-1 min-w-0">
      <HeaderBar
        @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
        @logout="handleLogout"
      />
      <main class="flex-1 overflow-auto bg-white">
        <!-- Main Header yang membaca Meta dari Route -->
        <div v-if="pageTitle" class="border-b border-gray-200 gap-1 flex-row px-4 py-1">
          <h1 class="text-lg font-semibold text-gray-800">{{ pageTitle }}</h1>
          <p v-if="pageDescription" class="text-xs text-gray-500">{{ pageDescription }}</p>
        </div>
        <div class="p-4">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/stores/auth'
import Sidebar from './Sidebar.vue'
import HeaderBar from './HeaderBar.vue'

const router = useRouter()
const route = useRoute()
const { logout } = useAuth()
const sidebarCollapsed = ref(false)
const pageTitle = computed(() => route.meta?.title || '')
const pageDescription = computed(() => route.meta?.description || '')

function handleLogout() {
  logout()
  router.push('/login')
}
</script>
