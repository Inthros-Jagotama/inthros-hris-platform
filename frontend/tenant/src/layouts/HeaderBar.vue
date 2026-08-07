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

      <!-- Breadcrumb: Employees > New/Edit Employee -->
      <template v-else-if="showEmployeeBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToEmployees"
        >
          {{ t('nav.employees') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ route.meta?.titleKey ? t(route.meta.titleKey) : '' }}</span>
      </template>

      <!-- Breadcrumb: Job Management > Manage Job Data -->
      <template v-else-if="showJobManagementBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToJobManagement"
        >
          {{ t('nav.job_management') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ t('job_management.manage') }}</span>
      </template>

      <!-- Breadcrumb: Job Values > type (kembali ke index card) -->
      <template v-else-if="showJobValuesBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToJobValues"
        >
          {{ t('nav.job_values_mapping') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ jobValueTypeLabel }}</span>
      </template>

      <!-- Breadcrumb: Settings > sub-setting (kembali ke index card) -->
      <template v-else-if="showSettingsBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToSettings"
        >
          {{ t('nav.settings') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ route.meta?.titleKey ? t(route.meta.titleKey) : (route.meta?.title || '') }}</span>
      </template>

      <!-- Breadcrumb: My Tasks > Approval Flows (kembali ke Tugas Saya) -->
      <template v-else-if="showApprovalFlowsBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-indigo-600 dark:hover:!text-indigo-400"
          @click="goBackToApprovals"
        >
          {{ t('approval.my_tasks') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ t('approval.flows') }}</span>
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
import { jobValueTypeLabel as jobValueTypeLabelFn } from '@/utils/jobValues'
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

/** Show breadcrumb when in Employee Form pages */
const showEmployeeBreadcrumb = computed(() => {
  return ['EmployeeNew', 'EmployeeEdit'].includes(route.name)
})

/** Show breadcrumb when in Job Management Form page */
const showJobManagementBreadcrumb = computed(() => {
  return route.name === 'JobManagementForm'
})

/** Show breadcrumb when on a Job Values type page (back to index card) */
const showJobValuesBreadcrumb = computed(() => {
  return route.name === 'JobValuesType'
})

// Label tipe utk breadcrumb — dari util bersama (bilingual + fallback)
const jobValueTypeLabel = computed(() => jobValueTypeLabelFn(t, route.params.type))

/** Show breadcrumb when on a sub-setting page (not the Settings index itself) */
const showSettingsBreadcrumb = computed(() => {
  return route.name?.startsWith('Settings') && route.name !== 'SettingsIndex'
})

/** Show breadcrumb when on the Approval Flows page (back to My Tasks) */
const showApprovalFlowsBreadcrumb = computed(() => {
  return route.name === 'ApprovalFlows'
})

function goBackToSummary() {
  router.push('/organization-summary')
}

function goBackToEmployees() {
  router.push('/employees')
}

function goBackToJobManagement() {
  router.push('/job-management')
}

function goBackToJobValues() {
  router.push('/job-management/values')
}

function goBackToSettings() {
  router.push('/settings')
}

function goBackToApprovals() {
  router.push('/approvals')
}

const userMenuItems = computed(() => [
  { label: t('auth.login.profile'), icon: 'pi pi-user', command: () => router.push('/profile') },
  { label: t('company_detail.title'), icon: 'pi pi-building', command: () => router.push('/company') },
  { separator: true },
  { label: t('auth.login.logout'), icon: 'pi pi-sign-out', command: () => emit('logout') }
])
</script>
