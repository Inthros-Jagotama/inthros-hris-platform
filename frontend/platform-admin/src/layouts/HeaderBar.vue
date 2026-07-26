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
      <span class="text-sm text-gray-500 font-medium truncate">{{ route.meta?.titleKey ? t(route.meta.titleKey) : (route.meta?.title || '') }}</span>
    </div>

    <div class="flex items-center gap-1">
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
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useLanguage } from '@/stores/language'
import { useI18n } from '@/composables/useI18n'
import Button from 'primevue/button'
import Avatar from 'primevue/avatar'
import Menu from 'primevue/menu'

const emit = defineEmits(['toggle-sidebar', 'logout'])

const route = useRoute()
const menu = ref(null)
const langStore = useLanguage()
const { t } = useI18n()

const router = useRouter()

const menuItems = computed(() => [
  { label: t('auth.login.profile'), icon: 'pi pi-user', command: () => router.push('/profile') },
  { separator: true },
  { label: t('auth.login.logout'), icon: 'pi pi-sign-out', command: () => emit('logout') }
])
</script>
