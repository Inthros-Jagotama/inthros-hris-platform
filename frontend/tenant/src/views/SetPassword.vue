<template>
  <div class="min-h-screen bg-gray-100 dark:bg-gray-900 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="bg-white dark:bg-gray-800 rounded-xl shadow-lg overflow-hidden">
        <div class="bg-emerald-600 px-6 py-4">
          <h1 class="text-lg font-semibold text-white">{{ t('set_password.title') }}</h1>
          <p class="text-sm text-emerald-100 mt-0.5">{{ t('set_password.description') }}</p>
        </div>

        <div v-if="success" class="p-6 text-center">
          <i class="pi pi-check-circle text-4xl text-emerald-500 mb-3"></i>
          <p class="text-sm text-gray-700 dark:text-gray-200 font-medium">{{ t('set_password.success') }}</p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-1 mb-4">{{ t('set_password.success_hint') }}</p>
          <Button :label="t('set_password.go_to_login')" size="small" @click="goLogin" />
        </div>

        <form v-else @submit.prevent="submit" class="p-6 space-y-4">
          <div v-if="!hasToken" class="flex items-start gap-2 text-amber-600 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3 text-xs">
            <i class="pi pi-exclamation-triangle mt-0.5"></i>
            <span>{{ t('set_password.missing_token') }}</span>
          </div>

          <template v-if="hasToken">
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('set_password.new_password') }}</label>
              <PasswordInput v-model="form.password" :placeholder="t('set_password.new_password_placeholder')" :invalid="!!errors?.password" />
              <small v-if="errors?.password" class="p-error block mt-1">{{ errors.password[0] }}</small>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('set_password.confirm_password') }}</label>
              <PasswordInput v-model="form.confirm" :placeholder="t('set_password.confirm_password_placeholder')" :invalid="!!errors?.confirm" />
              <small v-if="errors?.confirm" class="p-error block mt-1">{{ errors.confirm[0] }}</small>
            </div>
            <div v-if="errorMsg" class="text-xs text-red-500 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-3">
              {{ errorMsg }}
            </div>
            <Button type="submit" :label="t('set_password.button')" icon="pi pi-check" class="!w-full" :loading="loading" :disabled="loading || !hasToken" />
          </template>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/composables/useI18n'
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import PasswordInput from '@/components/PasswordInput.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const form = ref({ password: '', confirm: '' })
const errors = ref({})
const errorMsg = ref('')
const loading = ref(false)
const success = ref(false)

const token = computed(() => (route.query.token || '').toString())
const companyId = computed(() => (route.query.company_id || '').toString())
const hasToken = computed(() => token.value.length > 0)

async function submit() {
  errors.value = {}
  errorMsg.value = ''
  if (form.value.password.length < 8) { errors.value = { password: [t('set_password.password_min')] }; return }
  if (form.value.password !== form.value.confirm) { errors.value = { confirm: [t('set_password.password_mismatch')] }; return }
  loading.value = true
  try {
    const url = companyId.value
      ? `/api/v1/public/account/setup-password?company_id=${encodeURIComponent(companyId.value)}`
      : '/api/v1/public/account/setup-password'
    await api.post(url, { token: token.value, new_password: form.value.password })
    success.value = true
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { errorMsg.value = e.response?.data?.error?.message || t('set_password.invalid_link') }
  } finally {
    loading.value = false
  }
}

function goLogin() {
  router.push('/login')
}

onMounted(() => {
  if (!hasToken.value) {
    errorMsg.value = t('set_password.missing_token')
  }
})
</script>
