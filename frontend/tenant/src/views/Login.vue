<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-emerald-50 via-white to-emerald-100 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800 p-4">
    <div class="w-full max-w-sm">
      <!-- Login Card -->
      <Card class="!shadow-lg !rounded-xl">
        <template #content>
          <!-- Language Switcher (top-right) -->
          <div class="flex justify-end mb-2">
            <Button
              severity="secondary"
              text
              size="small"
              class="!p-1.5 !text-xs"
              v-tooltip.left="{ value: langStore.state.lang === 'en' ? 'Bahasa Indonesia' : 'English', showDelay: 300 }"
              @click="langStore.toggleLang()"
            >
              <div class="flex items-center gap-1">
                <i class="pi pi-globe text-xs"></i>
                <span class="text-xs font-semibold uppercase">{{ langStore.state.lang }}</span>
              </div>
            </Button>
          </div>
          <!-- Logo Area -->
          <div class="text-center mb-8">
            <div class="inline-flex items-center justify-center w-14 h-14 rounded-xl bg-emerald-600 text-white mb-4">
              <i class="pi pi-building text-2xl"></i>
            </div>
            <h1 class="text-xl font-bold text-gray-900 dark:text-gray-100">{{ t('auth.title') }}</h1>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">{{ t('auth.login.subtitle') }}</p>
          </div>
          <!-- Login Form -->
          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Email -->
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('auth.login.email') }}</label>
              <IconField>
                <InputIcon class="pi pi-envelope" />
                <InputText
                  v-model="email"
                  type="email"
                  :placeholder="t('auth.login.email_placeholder')"
                  class="!w-full"
                  :disabled="loading"
                  autocomplete="email"
                />
              </IconField>
            </div>
            <!-- Password -->
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('auth.login.password') }}</label>
              <IconField>
                <InputIcon class="pi pi-lock" />
                <InputText
                  v-model="password"
                  type="password"
                  :placeholder="t('auth.login.password_placeholder')"
                  class="!w-full"
                  :disabled="loading"
                  autocomplete="current-password"
                />
              </IconField>
            </div>
            <!-- Error -->
            <div v-if="error" class="text-sm text-rose-600 dark:text-rose-400 bg-rose-50 dark:bg-rose-900/20 border border-rose-200 dark:border-rose-800 rounded-md px-3 py-2">
              <i class="pi pi-exclamation-circle mr-1"></i> {{ error }}
            </div>
            <!-- Submit -->
            <Button
              type="submit"
              :label="t('auth.login.button')"
              icon="pi pi-sign-in"
              class="!w-full"
              :loading="loading"
            />
          </form>
        </template>
      </Card>
      <p class="text-center text-sm text-gray-400 dark:text-gray-500 mt-6">
        {{ t('auth.version') }}
      </p>
    </div>
  </div>
</template>
<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import { useLanguage } from '@/stores/language'
import { useI18n } from '@/composables/useI18n'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
const router = useRouter()
const { login } = useAuth()
const langStore = useLanguage()
const { t } = useI18n()
const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')
async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = t('auth.login.validation_required')
    return
  }
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.error?.message
      || e.response?.data?.message
      || t('auth.login.invalid_credentials')
  } finally {
    loading.value = false
  }
}
</script>
