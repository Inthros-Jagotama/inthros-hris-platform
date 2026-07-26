<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-50 via-white to-indigo-100 p-4">
    <div class="w-full max-w-sm">
      
      <!-- Login Card -->
      <Card class="!shadow-lg !rounded-xl">
        <template #content>
          <!-- Logo Area -->
          <div class="text-center mb-8">
            <div class="inline-flex items-center justify-center w-14 h-14 rounded-xl bg-indigo-600 text-white mb-4">
              <i class="pi pi-shield text-2xl"></i>
            </div>
            <h1 class="text-xl font-bold text-gray-900">HRIS Platform Admin</h1>
            <p class="text-sm text-gray-500 mt-1">Sign in to manage your platform</p>
          </div>
          <!-- Login Form -->
          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Email -->
            <div>
              <label class="block text-sm font-medium text-gray-600 mb-1">Email</label>
              <IconField>
                <InputIcon class="pi pi-envelope" />
                <InputText
                  v-model="email"
                  type="email"
                  placeholder="admin@company.com"
                  class="!w-full"
                  :disabled="loading"
                  autocomplete="email"
                />
              </IconField>
            </div>

            <!-- Password -->
            <div>
              <label class="block text-sm font-medium text-gray-600 mb-1">Password</label>
              <IconField>
                <InputIcon class="pi pi-lock" />
                <InputText
                  v-model="password"
                  type="password"
                  placeholder="••••••••"
                  class="!w-full"
                  :disabled="loading"
                  autocomplete="current-password"
                />
              </IconField>
            </div>

            <!-- Error -->
            <div v-if="error" class="text-sm text-rose-600 bg-rose-50 border border-rose-200 rounded-md px-3 py-2">
              <i class="pi pi-exclamation-circle mr-1"></i> {{ error }}
            </div>

            <!-- Submit -->
            <Button
              type="submit"
              label="Sign In"
              icon="pi pi-sign-in"
              class="!w-full"
              :loading="loading"
            />
          </form>
        </template>
      </Card>

      <p class="text-center text-sm text-gray-400 mt-6">
        HRIS Platform v1.6.3 &mdash; Enterprise Edition
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/stores/auth'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'

const router = useRouter()
const { login } = useAuth()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  if (!email.value || !password.value) {
    error.value = 'Please enter email and password'
    return
  }
  loading.value = true
  error.value = ''
  try {
    await login(email.value, password.value)
    router.push('/dashboard')
  } catch (e) {
    error.value = e.response?.data?.message || 'Invalid credentials. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
