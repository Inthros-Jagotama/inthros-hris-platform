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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
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
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
          @click="goBackToApprovals"
        >
          {{ t('approval.my_tasks') }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ t('approval.flows') }}</span>
      </template>

      <!-- Breadcrumb: generic — driven by route.meta.backRoute/backLabelKey -->
      <template v-else-if="showMetaBackBreadcrumb">
        <Button
          text
          size="small"
          class="!p-0 !text-xs !text-gray-500 dark:!text-gray-400 hover:!text-teal-600 dark:hover:!text-teal-400"
          @click="goBackToMetaRoute"
        >
          {{ t(route.meta.backLabelKey) }}
        </Button>
        <i class="pi pi-chevron-right text-xs text-gray-300"></i>
        <span class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ route.meta?.titleKey ? t(route.meta.titleKey) : (route.meta?.title || '') }}</span>
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
          @click="toggleNotifications"
        />
        <Badge
          v-if="notifState.unreadCount > 0"
          :value="notifState.unreadCount > 99 ? '99+' : notifState.unreadCount"
          severity="danger"
          class="!absolute -top-0.5 -right-0.5 !text-xs !min-w-[1.1rem] !h-[1.1rem]"
        />
      </div>
      <Popover ref="notificationsPanel" class="!w-80">
        <div class="flex items-center justify-between px-1 pb-2 border-b border-gray-100 dark:border-gray-700">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('notification.title') }}</span>
          <Button
            v-if="notifState.unreadCount > 0"
            :label="t('notification.mark_all_read')"
            text
            size="small"
            class="!p-1 !text-xs"
            @click="handleMarkAllAsRead"
          />
        </div>
        <div class="max-h-96 overflow-y-auto">
          <div v-if="notifState.recentItems.length === 0" class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
            <i class="pi pi-bell-slash text-2xl mb-2 opacity-50"></i>
            <p class="text-xs">{{ t('notification.empty') }}</p>
          </div>
          <div
            v-for="item in notifState.recentItems"
            :key="item.id"
            class="px-2 py-2 border-b border-gray-50 dark:border-gray-700/50 last:border-0 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/40 rounded"
            :class="{ 'bg-orange-50/60 dark:bg-orange-500/5': !item.is_read }"
            @click="handleNotificationClick(item)"
          >
            <div class="flex items-start gap-2">
              <span v-if="!item.is_read" class="w-1.5 h-1.5 rounded-full bg-orange-500 mt-1.5 shrink-0"></span>
              <span v-else class="w-1.5 h-1.5 shrink-0"></span>
              <div class="min-w-0 flex-1">
                <p class="text-xs font-medium text-gray-700 dark:text-gray-200 truncate">{{ item.title }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2">{{ item.body }}</p>
                <p class="text-[10px] text-gray-400 dark:text-gray-500 mt-0.5">{{ relativeTime(item.created_at) }}</p>
              </div>
            </div>
          </div>
        </div>
        <div class="pt-2 border-t border-gray-100 dark:border-gray-700">
          <Button :label="t('notification.view_all')" text size="small" class="w-full !text-xs" @click="goToNotifications" />
        </div>
      </Popover>

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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLanguage } from '@/stores/language'
import { useTheme } from '@/stores/theme'
import { useAuth } from '@/stores/auth'
import { useNotifications } from '@/stores/notifications'
import { useI18n } from '@/composables/useI18n'
import { jobValueTypeLabel as jobValueTypeLabelFn } from '@/utils/jobValues'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Badge from 'primevue/badge'
import Menu from 'primevue/menu'
import Popover from 'primevue/popover'

const emit = defineEmits(['toggle-sidebar', 'logout'])

const route = useRoute()
const router = useRouter()
const userMenu = ref(null)
const langStore = useLanguage()
const themeStore = useTheme()
const { state: authState } = useAuth()
const { t } = useI18n()

// ── Notifications ──
const notificationsPanel = ref(null)
const { state: notifState, refresh: refreshNotifications, markAsRead, markAllAsRead, startPolling, stopPolling } = useNotifications()

function toggleNotifications(event) {
  notificationsPanel.value.toggle(event)
  if (!notifState.loaded) refreshNotifications()
}

async function handleMarkAllAsRead() {
  await markAllAsRead()
}

async function handleNotificationClick(item) {
  if (!item.is_read) await markAsRead(item.id)
  notificationsPanel.value.hide()
  router.push('/notifications')
}

function goToNotifications() {
  notificationsPanel.value.hide()
  router.push('/notifications')
}

// relativeTime — short "X minutes/hours/days ago" label for the dropdown
// feed, built on the existing time.* locale keys (previously unused).
function relativeTime(value) {
  if (!value) return ''
  const diffMs = Date.now() - new Date(value).getTime()
  const diffMin = Math.floor(diffMs / 60000)
  if (diffMin < 1) return t('time.just_now')
  if (diffMin < 60) return t('time.minutes_ago', diffMin)
  const diffHours = Math.floor(diffMin / 60)
  if (diffHours < 24) return t('time.hours_ago', diffHours)
  const diffDays = Math.floor(diffHours / 24)
  return t('time.days_ago', diffDays)
}

onMounted(() => {
  refreshNotifications()
  startPolling()
})
onUnmounted(() => {
  stopPolling()
})

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

/** Generic breadcrumb fallback for any route declaring backRoute/backLabelKey in its meta */
const showMetaBackBreadcrumb = computed(() => {
  return !!(route.meta?.backRoute && route.meta?.backLabelKey)
})

function goBackToMetaRoute() {
  router.push(route.meta.backRoute)
}

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
