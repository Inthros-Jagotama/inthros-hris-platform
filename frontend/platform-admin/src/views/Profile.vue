<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <!-- User Info Card -->
    <Card>
      <template #title>
        <div class="flex items-center gap-3">
          <Avatar icon="pi pi-user" size="xlarge" class="!w-14 !h-14 !bg-indigo-100 !text-indigo-700 !text-2xl" />
          <div>
            <h2 class="text-xl font-semibold text-gray-800">{{ user?.name || '—' }}</h2>
            <p class="text-sm text-gray-500">{{ user?.email || '—' }}</p>
          </div>
        </div>
      </template>
      <template #content>
        <div class="grid grid-cols-2 gap-4 mt-2">
          <div>
            <label class="block text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('profile.role') }}</label>
            <Tag :value="t(`users.${user?.role}`)" :severity="user?.role === 'super_admin' ? 'danger' : 'info'" class="!text-xs" />
          </div>
          <div>
            <label class="block text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('profile.company') }}</label>
            <span class="text-sm text-gray-700">{{ user?.company_name || '—' }}</span>
          </div>
          <div>
            <label class="block text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('profile.status') }}</label>
            <Tag :value="user?.is_active !== false ? t('common_status.active') : t('common_status.inactive')"
              :severity="user?.is_active !== false ? 'success' : 'warn'" class="!text-xs" />
          </div>
          <div>
            <label class="block text-xs text-gray-400 uppercase tracking-wide mb-1">{{ t('profile.last_login') }}</label>
            <span class="text-sm text-gray-700">{{ user?.last_login ? formatDate(user.last_login) : '—' }}</span>
          </div>
        </div>
      </template>
    </Card>

    <!-- Change Password Card -->
    <Card>
      <template #title>
        <div class="flex items-center gap-2">
          <i class="pi pi-lock text-indigo-500"></i>
          <span class="text-lg font-semibold text-gray-800">{{ t('profile.change_password') }}</span>
        </div>
      </template>
      <template #content>
        <form @submit.prevent="handleChangePassword" class="space-y-4">
          <FormRow :label="t('profile.current_password')" :errors="errors?.current_password">
            <PasswordInput
              v-model="form.current_password"
              :class="{ 'p-invalid': errors?.current_password }"
              autofocus
            />
          </FormRow>

          <FormRow :label="t('profile.new_password')" :errors="errors?.new_password">
            <PasswordInput
              v-model="form.new_password"
              :class="{ 'p-invalid': errors?.new_password }"
            />
          </FormRow>

          <FormRow :label="t('profile.confirm_password')" :errors="errors?.confirm_password">
            <PasswordInput
              v-model="form.confirm_password"
              :class="{ 'p-invalid': errors?.confirm_password }"
            />
          </FormRow>

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
import { ref, reactive, onMounted } from 'vue'
import { useAuth } from '@/stores/auth'
import { useLanguage } from '@/stores/language'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Card from 'primevue/card'
import Avatar from 'primevue/avatar'
import Tag from 'primevue/tag'
import Button from 'primevue/button'
import FormRow from '@/components/FormRow.vue'
import PasswordInput from '@/components/PasswordInput.vue'
import { useToast } from 'primevue/usetoast'

const toast = useToast()
const { state: auth } = useAuth()
const langStore = useLanguage()
const { t } = useI18n()

const user = ref(auth.user)
const submitting = ref(false)
const errors = ref({})
console.log(user.value)
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
    errors.value.current_password = [t('form.required')]
    return
  }
  if (!form.new_password || form.new_password.length < 6) {
    errors.value.new_password = [t('form.min_length', { n: 6 })]
    return
  }
  if (form.new_password !== form.confirm_password) {
    errors.value.confirm_password = [t('profile.password_mismatch')]
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
