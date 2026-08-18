<template>
  <div class="flex h-screen overflow-hidden bg-gray-50 dark:bg-gray-900">
    <Sidebar
      :collapsed="sidebarCollapsed"
      @toggle="sidebarCollapsed = !sidebarCollapsed"
    />
    <div class="flex flex-col flex-1 min-w-0">
      <HeaderBar
        @toggle-sidebar="sidebarCollapsed = !sidebarCollapsed"
        @logout="handleLogout"
      />
      <main class="flex-1 overflow-auto bg-white dark:bg-gray-900">
        <!-- Main Header — Bilingual dari route name -->
        <div v-if="pageTitle" class="border-b border-gray-200 dark:border-gray-700 px-4 py-2">
          <h1 class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ pageTitle }}</h1>
          <p v-if="pageDescription" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ pageDescription }}</p>
        </div>
        <div class="p-4">
          <router-view />
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useActiveModules } from '@/stores/activeModules'
import { useI18n } from '@/composables/useI18n'
import { jobValueTypeLabel, jobValueTypeDesc } from '@/utils/jobValues'
import Sidebar from './Sidebar.vue'
import HeaderBar from './HeaderBar.vue'

const router = useRouter()
const route = useRoute()
const { logout } = useAuth()
const activeMod = useActiveModules()
const { t } = useI18n()
const sidebarCollapsed = ref(false)

// Fetch active modules on mount
onMounted(() => {
  activeMod.fetchActiveModules()
})

/**
 * Bilingual page title & description via route.meta.title / meta.description.
 * Falls back to hardcoded English meta.title / meta.description if no key.
 * Kasus khusus: halaman tipe Job Value (JobValuesType) — title = label tipe,
 * desc = deskripsi tipe (job_values.type_desc.*).
 */
const pageTitle = computed(() => {
  if (route.name === 'JobValuesType') return jobValueTypeLabel(t, route.params?.type)
  const key = route.meta?.titleKey
  if (key) return t(key)
  return route.meta?.title || ''
})

const pageDescription = computed(() => {
  if (route.name === 'JobValuesType') return jobValueTypeDesc(t, route.params?.type)
  const key = route.meta?.descKey
  if (key) {
    const desc = t(key)
    if (desc !== key) return desc
  }
  return route.meta?.description || ''
})

function handleLogout() {
  logout()
  router.push('/login')
}
</script>
