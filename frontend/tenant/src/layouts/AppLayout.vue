<template>
  <div class="flex h-screen overflow-hidden bg-gray-50">
    <!-- Sidebar -->
    <Sidebar
      :collapsed="sidebarCollapsed"
      @toggle="sidebarCollapsed = !sidebarCollapsed"
    />

    <!-- Main Content Area -->
    <div class="flex flex-col flex-1 min-w-0">
      <!-- Top Header Bar -->
      <HeaderBar @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed" />

      <!-- Page Content -->
      <main class="flex-1 overflow-auto p-4">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Sidebar from './Sidebar.vue'
import HeaderBar from './HeaderBar.vue'

const sidebarCollapsed = ref(false)
</script>
