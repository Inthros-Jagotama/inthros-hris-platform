<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-24 text-gray-400 dark:text-gray-500">
      <i class="pi pi-spin pi-spinner text-3xl mb-3"></i>
      <p class="text-sm">{{ t('common.loading') }}</p>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex flex-col items-center justify-center py-24 text-red-400">
      <i class="pi pi-exclamation-triangle text-3xl mb-3"></i>
      <p class="text-sm font-medium">{{ t('company_detail.load_error') }}</p>
    </div>

    <template v-else>
      <!-- Header Card -->
      <Card>
        <template #title>
          <div class="flex flex-wrap items-center gap-3">
            <div class="w-12 h-12 rounded-xl bg-emerald-100 dark:bg-emerald-900/50 text-emerald-600 dark:text-emerald-400 flex items-center justify-center">
              <i class="pi pi-building text-xl"></i>
            </div>
            <div class="min-w-0">
              <h2 class="text-xl font-semibold text-navy-800 dark:text-gray-100 truncate">{{ company.name || '—' }}</h2>
              <p class="text-sm text-gray-500 dark:text-gray-400">
                {{ t('company_detail.slug') }}: <span class="font-mono text-xs">{{ company.slug || '—' }}</span>
              </p>
            </div>
            <div class="ml-auto flex items-center gap-2">
              <Tag
                :value="statusLabel"
                :severity="statusSeverity"
                class="!text-xs"
              />
              <Button
                v-if="canEditCompany"
                :label="t('company_detail.edit')"
                icon="pi pi-pencil"
                size="small"
                outlined
                @click="openEditDialog"
              />
            </div>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-1">
            <ViewLabel :label="t('company_detail.email')" :value="company.email" break-all />
            <ViewLabel :label="t('company_detail.phone')" :value="company.phone" />
            <ViewLabel :label="t('company_detail.created_at')" :value="formatDate(company.created_at, locale)" />
            <ViewLabel :label="t('company_detail.updated_at')" :value="formatDate(company.updated_at, locale)" />
          </div>
        </template>
      </Card>

      <!-- Access & Domain -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-globe text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.access_info') }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <ViewLabel :label="t('company_detail.subdomain')" :value="company.subdomain" break-all />
            <ViewLabel :label="t('company_detail.domain')" :value="company.domain" break-all />
          </div>
        </template>
      </Card>

      <!-- Contact Info -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-phone text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.contact_info') }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <ViewLabel :label="t('company_detail.address')" :value="company.address" />
          </div>
        </template>
      </Card>

      <!-- Legal Info -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-id-card text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.legal_info') }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <ViewLabel :label="t('company_detail.npwp')" :value="company.npwp" mono break-all />
            <ViewLabel :label="t('company_detail.nib')" :value="company.nib" mono break-all />
          </div>
        </template>
      </Card>

      <!-- License & Subscription -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-key text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.license_info') }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <ViewLabel :label="t('company_detail.package')">
              <span class="font-medium text-navy-800 dark:text-gray-100">{{ company.license_info?.package_name || '—' }}</span>
            </ViewLabel>
            <ViewLabel :label="t('company_detail.plan_type')">
              <span class="capitalize">{{ company.license_info?.plan_type || '—' }}</span>
            </ViewLabel>
            <ViewLabel
              :label="t('company_detail.license_key')"
              :value="company.license_info?.license_key"
              mono
              break-all
              copyable
              :copy-title="t('common.copy')"
              :copied-title="t('common.copied')"
            />
            <ViewLabel :label="t('company_detail.start_date')" :value="formatDate(company.license_info?.start_date, locale)" />
            <ViewLabel :label="t('company_detail.end_date')" :value="formatDate(company.license_info?.end_date, locale)" />
            <ViewLabel :label="t('company_detail.valid_period')">
              <Tag
                v-if="company.license_info?.start_date"
                :value="isLicenseValid ? t('common_status.active') : t('common_status.expired')"
                :severity="isLicenseValid ? 'success' : 'danger'"
                class="!text-xs"
              />
              <span v-else>—</span>
            </ViewLabel>
            <ViewLabel :label="t('company_detail.max_employees')">
              <span class="font-medium text-navy-800 dark:text-gray-100">
                <template v-if="maxEmployees === null">—</template>
                <template v-else-if="maxEmployees === 0">∞</template>
                <template v-else>{{ maxEmployees.toLocaleString(locale === 'id' ? 'id-ID' : 'en-US') }}</template>
              </span>
            </ViewLabel>
            <ViewLabel :label="t('company_detail.employee_total')">
              <template v-if="!company.license_info">—</template>
              <template v-else>{{ employeeTotal.toLocaleString(locale === 'id' ? 'id-ID' : 'en-US') }}</template>
            </ViewLabel>
          </div>

          <!-- Employee usage progress (bila kuota > 0) -->
          <div v-if="maxEmployees > 0" class="mt-4 pt-4 border-t border-gray-100 dark:border-gray-800">
            <div class="flex items-center justify-between mb-1.5">
              <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('company_detail.employee_usage') }}</span>
              <span class="text-xs font-medium text-gray-700 dark:text-gray-300">
                {{ employeeTotal }} / {{ maxEmployees }} ({{ usagePercent }}%)
              </span>
            </div>
            <div class="w-full h-2 rounded-full bg-gray-200 dark:bg-gray-700 overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500"
                :class="usagePercent >= 100 ? 'bg-red-500' : usagePercent >= 80 ? 'bg-amber-500' : 'bg-emerald-500'"
                :style="{ width: Math.min(usagePercent, 100) + '%' }"
              ></div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Regional / Timezone Settings -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-clock text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.regional_settings') }}</span>
          </div>
        </template>
        <template #content>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 items-end">
            <FormRow :label="t('company_detail.timezone')">
              <Select
                v-model="timezone"
                :options="timezoneOptions"
                optionLabel="label"
                optionValue="value"
                class="w-full"
              />
            </FormRow>
            <div>
              <Button
                v-if="canEditCompany"
                :label="t('common.save')"
                icon="pi pi-check"
                size="small"
                :loading="savingTimezone"
                :disabled="savingTimezone || timezoneLoading"
                @click="handleSaveTimezone"
              />
            </div>
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-2">{{ t('company_detail.timezone_desc') }}</p>
        </template>
      </Card>

      <!-- Provisioning Info -->
      <Card>
        <template #title>
          <div class="flex items-center gap-2">
            <i class="pi pi-database text-emerald-500"></i>
            <span class="text-lg font-semibold text-navy-800 dark:text-gray-100">{{ t('company_detail.provisioning_info') }}</span>
          </div>
        </template>
        <template #content>
          <div v-if="!company.provisioning_info?.provisioned" class="text-sm text-gray-500 dark:text-gray-400">
            {{ t('company_detail.not_provisioned') }}
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <ViewLabel :label="t('company_detail.status')">
              <Tag :value="company.provisioning_info?.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="company.provisioning_info?.is_active ? 'success' : 'warn'" class="!text-xs" />
            </ViewLabel>
            <ViewLabel :label="t('company_detail.db_name')" :value="company.provisioning_info?.db_name" mono break-all />
            <ViewLabel :label="t('company_detail.driver')" :value="company.provisioning_info?.driver" />
          </div>
        </template>
      </Card>
    </template>

    <!-- Edit Company Info Dialog -->
    <Dialog
      v-model:visible="editDialogVisible"
      :header="t('company_detail.edit_title')"
      modal
      :style="{ width: '520px' }"
      :dismissableMask="true"
    >
      <div class="space-y-4 pt-2">
        <FormRow :label="t('company_detail.email')" :errors="errors?.email">
          <TextInput v-model="form.email" :placeholder="t('company_detail.email_placeholder')" :class="{ 'p-invalid': errors?.email }" />
        </FormRow>
        <FormRow :label="t('company_detail.phone')" :errors="errors?.phone">
          <TextInput v-model="form.phone" :placeholder="t('company_detail.phone_placeholder')" :class="{ 'p-invalid': errors?.phone }" />
        </FormRow>
        <FormRow :label="t('company_detail.address')" :errors="errors?.address">
          <TextInput v-model="form.address" textarea rows="3" :placeholder="t('company_detail.address_placeholder')" :class="{ 'p-invalid': errors?.address }" />
        </FormRow>
        <FormRow :label="t('company_detail.npwp')" :errors="errors?.npwp">
          <TextInput v-model="form.npwp" :placeholder="t('company_detail.npwp_placeholder')" :class="{ 'p-invalid': errors?.npwp }" />
        </FormRow>
        <FormRow :label="t('company_detail.nib')" :errors="errors?.nib">
          <TextInput v-model="form.nib" :placeholder="t('company_detail.nib_placeholder')" :class="{ 'p-invalid': errors?.nib }" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <Button
            :label="t('common.cancel')"
            severity="secondary"
            outlined
            size="small"
            @click="closeEditDialog"
          />
          <Button
            :label="t('common.save')"
            size="small"
            icon="pi pi-check"
            :loading="saving"
            :disabled="saving"
            @click="handleSave"
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuth } from '@/stores/auth'
import { useLanguage } from '@/stores/language'
import { useI18n } from '@/composables/useI18n'
import { formatDate } from '@/utils/formatDate'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Card from 'primevue/card'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Button from 'primevue/button'
import Select from 'primevue/select'
import { useToast } from 'primevue/usetoast'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ViewLabel from '@/components/ViewLabel.vue'

const toast = useToast()
const { state: auth, setCompany } = useAuth()
const langStore = useLanguage()
const { t } = useI18n()

const loading = ref(true)
const error = ref(false)
const company = ref({})

const locale = computed(() => langStore.state.lang)

// Hanya role administratif yang boleh mengubah profil perusahaan
// (konsisten dengan gate backend pada PUT /companies/me).
const canEditCompany = computed(() => {
  const role = (auth.user?.role || '').toLowerCase()
  return ['super_admin', 'company_admin', 'admin'].includes(role)
})

// ── Edit dialog ──
const editDialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref({
  email: '',
  phone: '',
  address: '',
  npwp: '',
  nib: ''
})

function openEditDialog() {
  form.value = {
    email: company.value.email || '',
    phone: company.value.phone || '',
    address: company.value.address || '',
    npwp: company.value.npwp || '',
    nib: company.value.nib || ''
  }
  errors.value = {}
  editDialogVisible.value = true
}

function closeEditDialog() {
  editDialogVisible.value = false
  errors.value = {}
}

// ── Regional settings: company default timezone ──
const timezoneOptions = [
  { label: 'WIB (Asia/Jakarta)', value: 'Asia/Jakarta' },
  { label: 'WITA (Asia/Makassar)', value: 'Asia/Makassar' },
  { label: 'WIT (Asia/Jayapura)', value: 'Asia/Jayapura' }
]
const timezone = ref('Asia/Jakarta')
const timezoneLoading = ref(false)
const savingTimezone = ref(false)

async function loadTimezone() {
  timezoneLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/company/timezone')
    const tz = res.data?.data?.timezone || res.data?.timezone
    if (tz) timezone.value = tz
  } catch {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('company_detail.timezone_load_error'), life: 4000 })
  } finally {
    timezoneLoading.value = false
  }
}

async function handleSaveTimezone() {
  savingTimezone.value = true
  try {
    await api.put('/api/v1/tenant/settings/company/timezone', { timezone: timezone.value })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('company_detail.timezone_updated'), life: 3000 })
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000
    })
  } finally {
    savingTimezone.value = false
  }
}

async function handleSave() {
  errors.value = {}
  saving.value = true
  try {
    // Kirim string kosong (bukan null) agar backend emptyToNil mengubahnya jadi
    // NULL — field yang dikosongkan benar-benar terhapus (nil pointer = skip).
    const res = await api.put('/api/v1/tenant/companies/me', {
      email: form.value.email || '',
      phone: form.value.phone || '',
      address: form.value.address || '',
      npwp: form.value.npwp || '',
      nib: form.value.nib || ''
    })
    company.value = res.data?.data || res.data || company.value
    setCompany({ id: company.value.id, name: company.value.name, slug: company.value.slug })
    closeEditDialog()
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('company_detail.updated'), life: 3000 })
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
    saving.value = false
  }
}

const statusLabel = computed(() => {
  const s = company.value.status
  if (!s) return '—'
  const key = `common_status.${s}`
  // Fallback ke raw status bila tidak ada locale key
  return t(key) !== key ? t(key) : s
})

const statusSeverity = computed(() => {
  switch (company.value.status) {
    case 'active': return 'success'
    case 'suspended': return 'warn'
    case 'terminated': return 'danger'
    default: return 'secondary'
  }
})

// Lisensi dianggap berlaku bila tanggal berakhir belum lewat hari ini.
// Bandingkan pada granularitas hari (bukan timestamp) agar tidak expired
// lebih awal karena perbedaan zona waktu UTC vs lokal (UTC+7).
const isLicenseValid = computed(() => {
  const end = company.value.license_info?.end_date
  if (!end) return true
  const endDate = new Date(String(end).slice(0, 10) + 'T23:59:59')
  if (isNaN(endDate.getTime())) return true
  return endDate >= new Date()
})

// Kuota employee dari lisensi (0 = unlimited, tampilkan ∞).
// null bila tidak ada license_info sama sekali (mis. on-premise) → tampilkan '—'.
const maxEmployees = computed(() => {
  const lic = company.value.license_info
  if (!lic) return null
  return lic.max_employees ?? 0
})

// Total employee terhitung di tenant database.
const employeeTotal = computed(() => company.value.license_info?.employee_total ?? 0)

// Persentase pemakaian kuota (dibulatkan).
const usagePercent = computed(() => {
  const max = maxEmployees.value
  if (!max || max <= 0) return 0
  return Math.round((employeeTotal.value / max) * 100)
})

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/tenant/companies/me')
    company.value = res.data?.data || res.data || {}
    // Sync nama company ke store bila response menyediakan name
    if (company.value.id && auth.company?.id !== company.value.id) {
      setCompany({ id: company.value.id, name: company.value.name, slug: company.value.slug })
    }
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
  loadTimezone()
})
</script>
