<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <!-- User Info Card -->
    <Card>
      <template #title>
        <div class="flex items-center gap-3">
          <Avatar icon="pi pi-user" size="xlarge" class="!w-14 !h-14 !bg-emerald-100 dark:!bg-emerald-900/50 !text-emerald-700 dark:!text-emerald-300 !text-2xl" />
          <div>
            <h2 class="text-xl font-semibold text-gray-800 dark:text-gray-100">{{ user?.name || '—' }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ user?.email || '—' }}</p>
          </div>
        </div>
      </template>
      <template #content>
        <div class="grid grid-cols-2 gap-4 mt-2">
          <div>
            <label class="block text-xs text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-1">{{ t('profile.role') }}</label>
            <Tag :value="roleLabel" severity="info" class="!text-xs" />
          </div>
          <div>
            <label class="block text-xs text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-1">{{ t('profile.company') }}</label>
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ auth.company?.name || user?.company_name || '—' }}</span>
          </div>
          <div>
            <label class="block text-xs text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-1">{{ t('profile.status') }}</label>
            <Tag
              :value="user?.is_active !== false ? t('common_status.active') : t('common_status.inactive')"
              :severity="user?.is_active !== false ? 'success' : 'warn'"
              class="!text-xs"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-400 dark:text-gray-500 uppercase tracking-wide mb-1">{{ t('profile.last_login') }}</label>
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ user?.last_login ? formatDate(user.last_login) : '—' }}</span>
          </div>
        </div>
      </template>
    </Card>
    <!-- Change Password Card (platform user only — employees use the email setup-password link) -->
    <Card v-if="isPlatformUser">
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-lock text-emerald-500"></i>
          <span class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('profile.change_password') }}</span>
        </div>
      </template>
      <template #content>
        <form @submit.prevent="handleChangePassword" class="space-y-4">
          <!-- Current Password -->
          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('profile.current_password') }}</label>
            <Password
              v-model="form.current_password"
              :feedback="false"
              toggleMask
              inputClass="!w-full"
              class="!w-full"
              :class="{ 'p-invalid': errors?.current_password }"
              autofocus
            />
            <small v-if="errors?.current_password" class="text-red-500 text-xs mt-1 block">
              {{ Array.isArray(errors.current_password) ? errors.current_password.join(', ') : errors.current_password }}
            </small>
          </div>
          <!-- New Password -->
          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('profile.new_password') }}</label>
            <Password
              v-model="form.new_password"
              :feedback="true"
              toggleMask
              inputClass="!w-full"
              class="!w-full"
              :class="{ 'p-invalid': errors?.new_password }"
            />
            <small v-if="errors?.new_password" class="text-red-500 text-xs mt-1 block">
              {{ Array.isArray(errors.new_password) ? errors.new_password.join(', ') : errors.new_password }}
            </small>
          </div>
          <!-- Confirm Password -->
          <div>
            <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('profile.confirm_password') }}</label>
            <Password
              v-model="form.confirm_password"
              :feedback="false"
              toggleMask
              inputClass="!w-full"
              class="!w-full"
              :class="{ 'p-invalid': errors?.confirm_password }"
            />
            <small v-if="errors?.confirm_password" class="text-red-500 text-xs mt-1 block">
              {{ Array.isArray(errors.confirm_password) ? errors.confirm_password.join(', ') : errors.confirm_password }}
            </small>
          </div>
          <div class="flex justify-end gap-2 pt-2">
            <Button
              :label="t('common.cancel')"
              severity="secondary"
              outlined
              size="small"
              @click="resetForm"
            />
            <Button
              :label="t('profile.update_password')"
              size="small"
              :loading="submitting"
              :disabled="submitting"
              type="submit"
            />
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>
<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuth } from '@/stores/auth'
import { useLanguage } from '@/stores/language'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Card from 'primevue/card'
import Avatar from 'primevue/avatar'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import Password from 'primevue/password'
import { useToast } from 'primevue/usetoast'
const toast = useToast()
const { state: auth } = useAuth()
const langStore = useLanguage()
const { t } = useI18n()
const user = ref(auth.user)
const submitting = ref(false)
const errors = ref({})
// Endpoint /api/v1/platform/users hanya berlaku utk platform user
// (company_admin/super_admin). Untuk login employee, data sudah tersedia dari
// response login (auth.user + company_name) — skip fetch agar tidak gagal.
const isPlatformUser = computed(() => ['company_admin', 'super_admin', 'admin'].includes(auth.user?.role))
onMounted(async () => {
  if (!isPlatformUser.value) return
  try {
    const res = await api.get(`/api/v1/platform/users/${auth.user.id}`)
    const data = res.data?.data || res.data
    user.value = { ...auth.user, ...data }
  } catch {
    // Fallback to auth store data if API fails
    user.value = auth.user
  }
})
const roleLabel = computed(() => {
  if (user.value?.role === 'company_admin') return t('profile.role_company_admin')
  if (user.value?.role === 'super_admin') return t('profile.role_super_admin')
  return user.value?.role || '—'
})
const form = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})
function resetForm() {
  form.current_password = ''
  form.new_password = ''
  form.confirm_password = ''
  errors.value = {}
}
function formatDate(dateStr) {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString(langStore.state.lang === 'id' ? 'id-ID' : 'en-US')
}
async function handleChangePassword() {
  errors.value = {}
  // Client-side validation
  if (!form.current_password) {
    errors.value = { current_password: [t('form.required')] }
    return
  }
  if (!form.new_password || form.new_password.length < 6) {
    errors.value = { new_password: [t('form.min_length', { n: 6 })] }
    return
  }
  if (form.new_password !== form.confirm_password) {
    errors.value = { confirm_password: [t('profile.password_mismatch')] }
    return
  }
  submitting.value = true
  try {
    await api.put(`/api/v1/platform/users/${user.value.id}/password`, {
      current_password: form.current_password,
      new_password: form.new_password
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('profile.password_updated'), life: 3000 })
    resetForm()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({
        severity: 'error',
        summary: t('message.error'),
        detail: e.response?.data?.error?.message || t('message.operation_failed'),
        life: 4000
      })
    }
  } finally {
    submitting.value = false
  }
}
</script>
