<template>
  <div class="space-y-2">
    <!-- Filters: Status Chips + Package Filter + Search + Actions -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-1.5">
        <Button
          v-for="chip in statusFilterChips"
          :key="chip.value"
          :label="chip.label"
          :severity="statusFilter === chip.value ? (chip.severity || 'secondary') : 'secondary'"
          :outlined="statusFilter !== chip.value"
          size="small"
          class="!text-xs !px-2 !py-1"
          @click="setStatusFilter(chip.value)"
        />
      </div>
      <div class="flex items-center gap-2">
        <SelectLabel
          v-if="packages.length > 0"
          v-model="packageFilter"
          :options="packageOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('licenses.filter_package')"
          class="!w-48 !h-8 !text-sm"
          showClear
          :filter="false"
        />
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" />
        </IconField>
        <Button :label="t('licenses.new_license')" icon="pi pi-plus" size="small" @click="openCreate" />
      </div>
    </div>

    <DataTable :value="filteredLicenses" paginator :rows="15" size="small" :loading="loading" class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-id-card text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('licenses.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('licenses.empty_hint') }}</p>
        </div>
      </template>
      <Column field="company_name" :header="t('licenses.company')" sortable />
      <Column field="package_name" :header="t('licenses.package')" sortable>
        <template #body="{ data }">
          <span v-if="data.package_name" class="text-gray-700">{{ data.package_name }}</span>
          <span v-else class="text-gray-400 italic">—</span>
        </template>
      </Column>
      <Column field="plan_type" :header="t('licenses.plan')" sortable>
        <template #body="{ data }">
          <Tag :value="data.plan_type" :severity="planSeverity(data.plan_type)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="seat_count" :header="t('licenses.seats')" sortable />
      <Column field="start_date" :header="t('licenses.valid_from')" sortable>
        <template #body="{ data }">{{ data.start_date || '-' }}</template>
      </Column>
      <Column field="end_date" :header="t('licenses.valid_until')" sortable>
        <template #body="{ data }">
          <span class="mr-1.5">{{ data.end_date || '-' }}</span>
          <Tag
            v-if="isExpiringSoon(data.end_date)"
            :value="t('licenses.expiring_soon')"
            severity="warn"
            class="!text-[10px] !px-1 !py-0"
          />
          <Tag
            v-else-if="isExpired(data.end_date)"
            :value="t('common_status.expired')"
            severity="danger"
            class="!text-[10px] !px-1 !py-0"
          />
        </template>
      </Column>
      <Column field="license_key" :header="t('licenses.key')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <code class="text-xs bg-gray-100 px-1.5 py-0.5 rounded font-mono truncate max-w-[160px] block">{{ data.license_key }}</code>
            <Button
              icon="pi pi-copy"
              size="small"
              text
              severity="secondary"
              class="!w-6 !h-6 !min-w-0"
              v-tooltip.left="{ value: t('licenses.tooltip_copy_key'), showDelay: 300 }"
              @click="copyKey(data.license_key)"
            />
          </div>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable>
        <template #body="{ data }">
          <Tag :value="data.status || 'active'" :severity="statusSeverity(data.status)" class="!text-xs" />
        </template>
      </Column>
      <Column :header="t('common.actions')" :style="{ width: '80px' }">
        <template #body="{ data }">
          <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEdit(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('licenses.edit_license') : t('licenses.new_license')" modal :style="{ width: '540px' }">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-x-4 gap-y-3">
          <div class="col-span-2">
            <FormRow :label="t('licenses.company')" :errors="errors?.company_id" :required="true">
              <SelectLabel v-model="form.company_id" :options="companies" optionLabel="name" optionValue="id" :placeholder="t('licenses.select_company')" :class="{ 'p-invalid': errors?.company_id }" />
            </FormRow>
          </div>
          <div>
            <FormRow :label="t('licenses.plan_type')" :errors="errors?.plan_type" :required="true">
              <SelectLabel v-model="form.plan_type" :options="planOptions" optionLabel="label" optionValue="value" :placeholder="t('licenses.select_plan')" :class="{ 'p-invalid': errors?.plan_type }" />
            </FormRow>
          </div>
          <div>
            <FormRow :label="t('licenses.package')" :errors="errors?.package_id">
              <SelectLabel v-model="form.package_id" :options="packageOptions" optionLabel="label" optionValue="value" :placeholder="t('licenses.select_package')" showClear :class="{ 'p-invalid': errors?.package_id }" :filter="false" />
            </FormRow>
          </div>
          <div>
            <FormRow :label="t('licenses.start_date')" :errors="errors?.start_date">
              <DateInput v-model="form.start_date" showClear :class="{ 'p-invalid': errors?.start_date }" />
            </FormRow>
          </div>
          <div>
            <FormRow :label="t('licenses.end_date')" :errors="errors?.end_date">
              <DateInput v-model="form.end_date" showClear :class="{ 'p-invalid': errors?.end_date }" />
            </FormRow>
          </div>
          <div>
            <FormRow :label="t('licenses.seats')" :errors="errors?.seat_count" :required="true">
              <InputNumber v-model="form.seat_count" class="!w-full" inputClass="!w-full !text-sm" :min="1" :class="{ 'p-invalid': errors?.seat_count }" />
            </FormRow>
          </div>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('common.create')" size="small" :loading="saving" :disabled="saving" @click="saveLicense" />
      </template>
    </Dialog>

    <!-- Toast for copy confirmation (handled via toast directly) -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import InputText from 'primevue/inputtext'

const toast = useToast()
const { t } = useI18n()

// Data
const licenses = ref([])
const companies = ref([])
const packages = ref([])
const loading = ref(true)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ company_id: null, plan_type: 'trial', package_id: null, seat_count: 1, start_date: '', end_date: '' })
const errors = ref({})

// Filters
const searchQuery = ref('')
const statusFilter = ref(null)
const packageFilter = ref(null)

// Plan options
const planOptions = computed(() => [
  { label: t('licenses.plan_trial'), value: 'trial' },
  { label: t('licenses.plan_basic'), value: 'basic' },
  { label: t('licenses.plan_professional'), value: 'professional' },
  { label: t('licenses.plan_enterprise'), value: 'enterprise' }
])

// Status filter chips
const statusFilterChips = computed(() => [
  { label: t('common.all'), value: null, severity: 'info' },
  { label: t('common_status.active'), value: 'active', severity: 'success' },
  { label: t('common_status.expired'), value: 'expired', severity: 'danger' },
  { label: t('common_status.suspended'), value: 'suspended', severity: 'warn' }
])

// Package filter options
const packageOptions = computed(() =>
  packages.value.map(p => ({
    label: p.name,
    value: p.id
  }))
)

// Filtered licenses
const filteredLicenses = computed(() => {
  let result = licenses.value

  // Status filter
  if (statusFilter.value) {
    result = result.filter(l => l.status === statusFilter.value)
  }

  // Package filter
  if (packageFilter.value) {
    result = result.filter(l => l.package_id === packageFilter.value)
  }

  // Search
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(l =>
      l.company_name?.toLowerCase().includes(q) ||
      l.plan_type?.toLowerCase().includes(q) ||
      l.license_key?.toLowerCase().includes(q) ||
      l.package_name?.toLowerCase().includes(q)
    )
  }

  return result
})

// Expiration helpers
function isExpiringSoon(dateStr) {
  if (!dateStr) return false
  const end = new Date(dateStr)
  const now = new Date()
  const diffMs = end.getTime() - now.getTime()
  const diffDays = diffMs / (1000 * 60 * 60 * 24)
  return diffDays > 0 && diffDays <= 30
}

function isExpired(dateStr) {
  if (!dateStr) return false
  const end = new Date(dateStr)
  const now = new Date()
  return end.getTime() < now.getTime()
}

// Load data
async function loadData() {
  loading.value = true
  try {
    const [licRes, compRes, pkgRes] = await Promise.all([
      api.get('/api/v1/platform/licenses'),
      api.get('/api/v1/platform/companies'),
      api.get('/api/v1/platform/packages?per_page=100')
    ])
    const licPayload = licRes.data
    licenses.value = Array.isArray(licPayload.data) ? licPayload.data : (Array.isArray(licPayload) ? licPayload : [])
    const compPayload = compRes.data
    companies.value = Array.isArray(compPayload.data) ? compPayload.data : (Array.isArray(compPayload) ? compPayload : [])
    const pkgPayload = pkgRes.data
    packages.value = Array.isArray(pkgPayload.data) ? pkgPayload.data : (Array.isArray(pkgPayload) ? pkgPayload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

// Severity helpers
function planSeverity(plan) {
  switch (plan) {
    case 'enterprise': return 'danger'
    case 'professional': return 'warn'
    case 'basic': return 'info'
    default: return 'success'
  }
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'expired': return 'danger'
    case 'suspended': return 'warn'
    default: return 'info'
  }
}

// Filter handlers
function setStatusFilter(value) {
  statusFilter.value = value
}

// Copy license key
async function copyKey(key) {
  try {
    await navigator.clipboard.writeText(key)
    toast.add({ severity: 'success', summary: t('licenses.key_copied'), life: 1500 })
  } catch {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.operation_failed'), life: 2000 })
  }
}

// Dialog handlers
function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { company_id: null, plan_type: 'trial', package_id: null, seat_count: 1, start_date: '', end_date: '' }
  errors.value = {}
  dialogVisible.value = true
}

function openEdit(lic) {
  isEditing.value = true
  editingId.value = lic.id
  errors.value = {}
  form.value = {
    company_id: lic.company_id,
    plan_type: lic.plan_type,
    package_id: lic.package_id || null,
    seat_count: lic.seat_count,
    start_date: lic.start_date,
    end_date: lic.end_date
  }
  dialogVisible.value = true
}

async function saveLicense() {
  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/platform/licenses/${editingId.value}`, form.value)
      toast.add({ severity: 'success', summary: t('message.updated'), life: 2000 })
    } else {
      await api.post('/api/v1/platform/licenses', form.value)
      toast.add({ severity: 'success', summary: t('message.created'), life: 2000 })
    }
    dialogVisible.value = false
    await loadData()
  } catch (e) {
    errors.value = getValidationErrors(e)
    if (Object.keys(errors.value).length === 0) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
    }
  } finally {
    saving.value = false
  }
}
</script>
