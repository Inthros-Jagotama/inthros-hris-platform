<template>
  <header class="flex items-center h-12 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 shrink-0">
    <div class="flex items-center gap-3 flex-1 min-w-0">
      <Button
        icon="pi pi-bars"
        severity="secondary"
        text
        size="small"
        @click="$emit('toggle-sidebar')"
        class="!p-1.5"
      />

      <!-- Breadcrumb: Organization Summary > Organization -->
      <template v-if="showOrgBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToSummary"
        >
          {{ t('org_summary.title') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ t('organization.title') }}</span>
      </template>

      <!-- Normal page title -->
      <template v-else>
        <i class="pi pi-chevron-right text-sm text-gray-300"></i>
        <span class="text-sm text-gray-500 dark:text-gray-400 font-medium truncate">{{ route.meta?.titleKey ? t(route.meta.titleKey) : (route.meta?.title || '') }}</span>
      </template>
    </div>

    <div class="flex items-center gap-1">
      <!-- Theme Switcher -->
      <Button
        severity="secondary"
        text
        size="small"
        class="!p-1.5"
        v-tooltip.top="{ value: themeStore.isDark() ? t('dashboard.light_mode') : t('dashboard.dark_mode'), showDelay: 300 }"
        @click="themeStore.toggleTheme()"
      >
        <i :class="themeStore.isDark() ? 'pi pi-sun' : 'pi pi-moon'" class="text-sm"></i>
      </Button>

      <!-- Language Switcher -->
      <Button
        severity="secondary"
        text
        size="small"
        class="!p-1.5"
        v-tooltip.top="{ value: langStore.state.lang === 'en' ? 'Bahasa Indonesia' : 'English', showDelay: 300 }"
        @click="langStore.toggleLang()"
      >
        <div class="flex items-center gap-1">
          <i class="pi pi-globe text-sm"></i>
          <span class="text-xs font-semibold uppercase">{{ langStore.state.lang }}</span>
        </div>
      </Button>

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

      <div class="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400 mr-2">
        <i class="pi pi-circle-fill text-emerald-400 text-[6px]"></i>
        <span>Live</span>
      </div>

      <Button
        severity="secondary"
        text
        size="small"
        class="!p-1"
        @click="userMenu.toggle($event)"
      >
        <div class="flex items-center gap-2">
          <Avatar
            icon="pi pi-user"
            size="small"
            class="!w-7 !h-7 !bg-emerald-100 !text-emerald-700"
          />
          <span class="text-sm text-gray-700 dark:text-gray-200 hidden sm:inline">{{ authState.user?.name || 'Admin' }}</span>
          <i class="pi pi-chevron-down text-sm text-gray-400"></i>
        </div>
      </Button>
      <Menu ref="userMenu" :model="userMenuItems" popup />
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLanguage } from '@/stores/language'
import { useTheme } from '@/stores/theme'
import { useAuth } from '@/stores/auth'
import { useI18n } from '@/composables/useI18n'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Badge from 'primevue/badge'
import Menu from 'primevue/menu'

const emit = defineEmits(['toggle-sidebar', 'logout'])

const route = useRoute()
const router = useRouter()
const userMenu = ref(null)
const langStore = useLanguage()
const themeStore = useTheme()
const { state: authState } = useAuth()
const { t } = useI18n()

/** Show breadcrumb when visiting Organizations with summary_id query param */
const showOrgBreadcrumb = computed(() => {
  return route.name === 'Organizations' && route.query?.summary_id
})

function goBackToSummary() {
  router.push('/organization-summary')
}

const userMenuItems = computed(() => [
  { label: t('auth.login.profile'), icon: 'pi pi-user', command: () => router.push('/profile') },
  { separator: true },
  { label: t('auth.login.logout'), icon: 'pi pi-sign-out', command: () => emit('logout') }
])
</script>
