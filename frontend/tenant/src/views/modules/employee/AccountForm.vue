<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h3 class="text-medium font-semibold text-gray-700 dark:text-gray-300">{{ t('employee.tab_account') }}</h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee.account_description') }}</p>
      </div>
    </div>

    <!-- Loading skeleton → endpoint status akun sedang dimuat -->
    <div v-if="statusLoading" class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 max-w-xl">
      <div class="flex items-center gap-2 mb-4">
        <div class="w-6 h-6 bg-gray-200 dark:bg-gray-700 rounded-full animate-pulse"></div>
        <div class="h-4 w-48 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
      <div class="space-y-3">
        <div class="h-4 w-24 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="h-9 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="h-9 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="h-4 w-32 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="flex justify-end pt-2">
          <div class="h-8 w-28 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- No account yet — create form -->
    <div v-else-if="!account" class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 max-w-xl">
      <p class="text-sm text-gray-500 dark:text-gray-400 mb-4">{{ t('employee.account_no_account') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('employee.account_email')" required :errors="errors?.email">
          <TextInput v-model="email" type="email" maxlength="255" :placeholder="t('employee.account_email_placeholder')" :class="{ 'p-invalid': errors?.email }" autofocus />
          <small v-if="employeeEmail" class="text-xs text-gray-400 dark:text-gray-500 mt-1 flex items-center gap-1">
            <i class="pi pi-user"></i> {{ t('employee.account_email_from_data') }}: {{ employeeEmail }}
          </small>
        </FormRow>
        <div class="flex justify-end">
          <Button :label="t('employee.account_create')" icon="pi pi-envelope" size="small" :loading="loading" :disabled="loading" @click="createAccount" />
        </div>
      </div>
    </div>

    <!-- Account exists — status card -->
    <div v-else-if="account" class="border border-gray-200 dark:border-gray-700 rounded-lg p-5 max-w-xl">
      <div class="flex items-center gap-2 mb-4">
        <i class="pi pi-check-circle text-emerald-500 text-lg"></i>
        <span class="text-sm font-medium text-gray-700 dark:text-gray-200">{{ account.email }}</span>
      </div>
      <dl class="grid grid-cols-1 gap-3 text-sm">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-gray-800 pb-2">
          <dt class="text-gray-500 dark:text-gray-400">{{ t('employee.account_role') }}</dt>
          <dd class="font-medium text-gray-700 dark:text-gray-200">{{ account.role_name || 'Employee' }}</dd>
        </div>
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-gray-800 pb-2">
          <dt class="text-gray-500 dark:text-gray-400">{{ t('employee.account_user_id') }}</dt>
          <dd class="font-mono text-xs text-gray-600 dark:text-gray-300">{{ account.user_id }}</dd>
        </div>
        <div class="flex items-center justify-between pb-1">
          <dt class="text-gray-500 dark:text-gray-400">{{ t('common.status') }}</dt>
          <dd>
            <Tag v-if="account.password_set" :value="t('employee.account_password_set')" severity="success" class="!text-xs" />
            <Tag v-else :value="t('employee.account_password_not_set')" severity="warn" class="!text-xs" />
          </dd>
        </div>
      </dl>
      <div class="flex justify-end mt-4">
        <Button :label="t('employee.account_resend')" icon="pi pi-send" size="small" severity="secondary" outlined :loading="resending" :disabled="resending" @click="resendEmail" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/composables/useI18n'
import { ref, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'

const { t } = useI18n()
const toast = useToast()

const props = defineProps({
  employeeId: { type: String, default: '' },
  employeeEmail: { type: String, default: '' }
})
const emit = defineEmits(['save'])

// Default email dari data personal employee (bisa diubah admin).
const email = ref(props.employeeEmail || '')
const account = ref(null)
const statusLoading = ref(false)
const loading = ref(false)
const resending = ref(false)
const errors = ref({})

async function loadStatus() {
  if (!props.employeeId) return
  statusLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/user-accounts/employees/${props.employeeId}`)
    account.value = res.data?.data || null
    // Employee yang sudah punya akun → tandai step 'saved' di navigasi
    if (account.value) emit('save')
  } catch {
    account.value = null // 404 → belum ada akun
  } finally {
    statusLoading.value = false
  }
}

async function createAccount() {
  errors.value = {}
  if (!email.value?.trim()) { errors.value = { email: [t('employee.account_email_required')] }; return }
  loading.value = true
  try {
    const res = await api.post(`/api/v1/tenant/user-accounts/employees/${props.employeeId}`, { email: email.value.trim() })
    account.value = res.data?.data || null
    email.value = ''
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.account_created'), life: 3000 })
    emit('save')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    loading.value = false
  }
}

async function resendEmail() {
  resending.value = true
  try {
    await api.post(`/api/v1/tenant/user-accounts/employees/${props.employeeId}/resend`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.account_resended'), life: 3000 })
    emit('save')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    resending.value = false
  }
}

// Jika email employee baru tersedia setelah load (mis. mode edit), isi otomatis
// hanya saat field masih kosong agar tidak menimpa input admin.
watch(() => props.employeeEmail, (val) => {
  if (val && !email.value) email.value = val
})

onMounted(loadStatus)
</script>
