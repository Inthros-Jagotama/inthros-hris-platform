<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <!-- Left: Filter Chips -->
      <div class="flex items-center gap-1.5">
        <Button
          v-for="chip in filterChips"
          :key="chip.value"
          :label="chip.label"
          :severity="statusFilter === chip.value ? (chip.severity || 'secondary') : 'secondary'"
          :outlined="statusFilter !== chip.value"
          size="small"
          class="!text-xs !px-2 !py-1"
          @click="statusFilter = chip.value"
        />
      </div>
      <!-- Right: Search + Actions -->
      <div class="flex items-center gap-2">
        <Select v-if="packages.length > 0" v-model="packageFilter" :options="packageOptions" optionLabel="label" optionValue="value" :placeholder="t('companies.filter_package')" class="!w-44 !h-8 !text-xs" size="small" showClear />
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('companies.search')" size="small" />
        </IconField>
        <Button :label="t('companies.new_company')" icon="pi pi-plus" size="small" @click="openCreate" />
      </div>
    </div>
    <!-- DataTable -->
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable 
    v-else
    :value="filteredCompanies" 
    paginator 
    :rows="15" 
    sortField="createdAt" 
    :sortOrder="-1" 
    size="small" 
    class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-building text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('companies.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('companies.empty_hint') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('companies.company_name')" sortable>
        <template #body="{ data }">
          <div class="flex-row gap-2">
            <div class="uppercase font-semibold text-gray-600 dark:text-gray-300">{{ data.name }}</div>
            <div class="text-sm text-gray-500 dark:text-gray-400">{{ data.address }}</div>
          </div>
        </template>
      </Column>
      <Column field="email" :header="t('companies.email')" sortable />
      <Column field="phone" :header="t('companies.phone')" sortable />
      <!-- License Info Column -->
      <Column field="license_info.plan_type" :header="t('companies.license_plan')" sortable :sortField="data => data.license_info?.plan_type || ''">
        <template #body="{ data }">
          <template v-if="data.license_info">
            <Tag :value="data.license_info.plan_type" :severity="planSeverity(data.license_info.plan_type)" class="!text-xs !px-1.5 !py-0.5" />
          </template>
          <span v-else class="text-gray-300 dark:text-gray-600 italic text-xs">—</span>
        </template>
      </Column>
      <!-- Provisioning Status Column -->
      <Column field="provisioning_info.provisioned" :header="t('companies.provision_header')" sortable :sortField="data => data.provisioning_info?.provisioned ? 1 : 0">
        <template #body="{ data }">
          <template v-if="data.provisioning_info">
            <Tag
              v-if="data.provisioning_info.provisioned && data.provisioning_info.is_active !== false"
              :value="t('companies.provision_status_provisioned')"
              severity="success"
              class="!text-[10px] !px-1.5 !py-0"
              v-tooltip.left="provisionTooltip(data)"
            />
            <Tag
              v-else-if="data.provisioning_info.provisioned && data.provisioning_info.is_active === false"
              :value="t('companies.provision_status_deactivated')"
              severity="warn"
              class="!text-[10px] !px-1.5 !py-0"
              v-tooltip.left="provisionTooltip(data)"
            />
            <Tag
              v-else
              :value="t('companies.provision_status_not_provisioned')"
              severity="danger"
              class="!text-[10px] !px-1.5 !py-0"
            />
          </template>
          <span v-else class="text-gray-300 dark:text-gray-600 italic text-xs">—</span>
        </template>
      </Column>

      <Column field="status" :header="t('common.status')" sortable>
        <template #body="{ data }">
          <Tag :value="data.status" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="createdAt" :header="t('companies.created')" sortable>
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400">{{ data.createdAt || '-' }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" :style="{ width: '220px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('companies.view_detail')" @click="openDetail(data)" />
            <CompanyActions :company="data" mode="icons" @updated="loadData" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create Dialog (edit dipindah ke komponen CompanyActions) -->
    <Dialog v-model:visible="dialogVisible" :header="t('companies.new_company')" modal :style="{ width: '620px' }" :closable="true">
      <div class="space-y-4">
        <!-- Company Info -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5">
            <i class="pi pi-building text-indigo-400 text-sm"></i>
            {{ t('companies.new_company') }}
          </h3>
          <div class="space-y-3">
            <FormRow :label="t('companies.company_name')" :errors="errors?.name" :required="true">
              <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
            </FormRow>
            <FormRow :label="t('companies.email')" :errors="errors?.email" :required="true">
              <TextInput v-model="form.email" autofocus :class="{ 'p-invalid': errors?.email }" />
            </FormRow>
            <FormRow :label="t('companies.phone')" :errors="errors?.phone" :required="true">
              <TextInput v-model="form.phone" autofocus :class="{ 'p-invalid': errors?.phone }" />
            </FormRow>
            <FormRow :label="t('companies.address')" :errors="errors?.address" :required="true">
              <TextInput v-model="form.address" autofocus :class="{ 'p-invalid': errors?.address }" />
            </FormRow>
          </div>
        </div>

        <!-- Admin User -->
        <div>
          <div class="border-t border-gray-200 dark:border-gray-700 my-3"></div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5">
            <i class="pi pi-user text-indigo-400 text-sm"></i>
            {{ t('companies.admin_section_title') }}
          </h3>
          <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">{{ t('companies.admin_section_hint') }}</p>
          <div class="space-y-3">
            <FormRow :label="t('companies.admin_name')" :errors="errors?.admin_name" :required="true">
              <TextInput v-model="form.admin_name" autofocus :class="{ 'p-invalid': errors?.admin_name }" />
            </FormRow>
            <FormRow :label="t('companies.admin_email')" :errors="errors?.admin_email" :required="true">
              <TextInput v-model="form.admin_email" autofocus :class="{ 'p-invalid': errors?.admin_email }" />
            </FormRow>
            <FormRow :label="t('companies.admin_password')" :errors="errors?.admin_password" :required="true">
              <PasswordInput v-model="form.admin_password" autofocus :class="{ 'p-invalid': errors?.admin_password }" />
            </FormRow>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ml-auto">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
            <Button :label="t('common.create')" size="small" :loading="saving" :disabled="saving" @click="saveCompany" />
          </div>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import PasswordInput from '@/components/PasswordInput.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import CompanyActions from '@/components/CompanyActions.vue'

const toast = useToast()
const router = useRouter()
const { t } = useI18n()

// Data
const companies = ref([])
const searchQuery = ref('')
const statusFilter = ref(null)
const packageFilter = ref(null)
const packages = ref([])
const loading = ref(true)
const dialogVisible = ref(false)
const saving = ref(false)
const form = ref({ name: '', email: '', phone: '', address: '', admin_name: '', admin_email: '', admin_password: '' })
const errors = ref({})

// Package options untuk filter dropdown
const packageOptions = computed(() => {
  const opts = packages.value.map(p => ({
    label: p.name,
    value: p.id
  }))
  return opts
})

// Filter chips untuk quick visual filter
const filterChips = computed(() => [
  { label: t('companies.all_status'), value: null, severity: 'info' },
  { label: t('common_status.active'), value: 'active', severity: 'success' },
  { label: t('common_status.suspended'), value: 'suspended', severity: 'warn' },
  { label: t('common_status.terminated'), value: 'terminated', severity: 'danger' }
])

const skeletonColumns = [
  { type: 'compound', widths: ['w-28', 'w-36'], headerWidth: 'w-20' },
  { type: 'text', width: 'w-32', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-14', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-14', headerWidth: 'w-14' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 5, headerWidth: 'w-16' }
]

const filteredCompanies = computed(() => {
  let result = companies.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c =>
      c.name?.toLowerCase().includes(q) ||
      c.email?.toLowerCase().includes(q) ||
      c.slug?.toLowerCase().includes(q) ||
      c.phone?.toLowerCase().includes(q) ||
      c.address?.toLowerCase().includes(q)
    )
  }
  if (statusFilter.value) {
    result = result.filter(c => c.status === statusFilter.value)
  }
  if (packageFilter.value) {
    result = result.filter(c => c.license_info?.package_id === packageFilter.value)
  }
  return result
})

// License plan severity
function planSeverity(plan) {
  switch (plan?.toLowerCase()) {
    case 'enterprise': return 'danger'
    case 'professional': return 'warn'
    case 'basic': return 'info'
    case 'subscription': return 'info'
    case 'trial': return 'success'
    default: return 'info'
  }
}

// Load
async function loadData() {
  loading.value = true
  try {
    // Load companies (wajib)
    try {
      const compRes = await api.get('/api/v1/platform/companies')
      const compPayload = compRes.data
      companies.value = Array.isArray(compPayload.data) ? compPayload.data : (Array.isArray(compPayload) ? compPayload : [])
    } catch (e) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
    }
    // Load packages (opsional — jika gagal, filter package tidak muncul)
    try {
      const pkgRes = await api.get('/api/v1/platform/packages?per_page=100')
      const pkgPayload = pkgRes.data
      packages.value = Array.isArray(pkgPayload.data) ? pkgPayload.data : (Array.isArray(pkgPayload) ? pkgPayload : [])
    } catch (e) {
      console.warn('[Companies] Failed to load packages:', e.response?.data || e.message)
      packages.value = []
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadData()
})

// Helpers
function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'suspended': return 'warn'
    case 'terminated': return 'danger'
    default: return 'info'
  }
}

function provisionTooltip(company) {
  if (!company.provisioning_info) return ''
  const info = company.provisioning_info
  let tip = `${t('companies.provision_header')}: ${info.db_name || 'N/A'}`
  if (info.driver) tip += ` | Driver: ${info.driver}`
  tip += ` | ${t('companies.provision_active')}: ${info.is_active !== false ? t('common_status.yes') : t('common_status.no')}`
  tip += ` | ${t('companies.provision_provisioned')}: ${info.provisioned ? t('common_status.yes') : t('common_status.no')}`
  return tip
}

function openDetail(company) {
  router.push(`/companies/${company.id}`)
}

function openCreate() {
  form.value = { name: '', email: '', phone: '', address: '', admin_name: '', admin_email: '', admin_password: '' }
  errors.value = {}
  dialogVisible.value = true
}

async function saveCompany() {
  saving.value = true
  try {
    const res = await api.post('/api/v1/platform/companies', form.value)
    const admin = res.data?.data?.admin_user
    if (admin) {
      toast.add({
        severity: 'success',
        summary: t('message.created'),
        detail: `${t('companies.title')}: ${admin.name} (${admin.email})`,
        life: 5000
      })
    } else {
      toast.add({ severity: 'success', summary: t('message.created'), detail: t('companies.title'), life: 2000 })
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
