<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-end gap-2">
      <!-- Toolbar -->
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" placeholder="Search companies..." size="small" />
        </IconField>
        <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" size="small" optionValue="value" placeholder="Status" />
      </div>
      <Button label="New Company" icon="pi pi-plus" size="small" @click="openCreate" />
    </div>
    <!-- DataTable -->
    <DataTable 
    :value="filteredCompanies" 
    paginator 
    :rows="15" 
    sortField="createdAt" 
    :sortOrder="-1" 
    size="small" 
    class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-building text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">No companies registered yet</p>
          <p class="text-sm mt-1">Click "New Company" to create the first company.</p>
        </div>
      </template>
      <Column field="name" header="Company Name" sortable>
        <template #body="{ data }">
          <div class="flex-row gap-2">
            <div class="uppercase font-semibold text-gray-600">{{ data.name }}</div>
            <div class="text-sm text-gray-500">{{ data.address }}</div>
          </div>
        </template>
      </Column>
      <Column field="email" header="Email" sortable />
      <Column field="phone" header="Phone" sortable />
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :value="data.status" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="createdAt" header="Created" sortable>
        <template #body="{ data }">
          <span class="text-gray-500">{{ data.createdAt || '-' }}</span>
        </template>
      </Column>
      <Column header="Actions" :style="{ width: '160px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="'Edit'" @click="openEdit(data)" />
            <Button v-if="data.status === 'active'" icon="pi pi-pause-circle" size="small" text severity="warning" v-tooltip.left="'Suspend'" @click="confirmSuspend(data)" />
            <Button v-if="data.status === 'suspended'" icon="pi pi-play-circle" size="small" text severity="info" v-tooltip.left="'Activate'" @click="confirmActivate(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="'Terminate'" @click="confirmTerminate(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? 'Edit Company' : 'New Company'" modal :style="{ width: '600px' }" :closable="true">
      <div class="space-y-4">
        <!-- Company Info -->
        <div>
          <h3 class="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-1.5">
            <i class="pi pi-building text-indigo-400 text-sm"></i>
            Company Information
          </h3>
          <div class="space-y-2">
            <div>
              <FormRow label="Company Name" :errors="errors?.name" :required="true">
                  <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
              </FormRow>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div>
                <FormRow label="Email" :errors="errors?.email" :required="true">
                    <TextInput v-model="form.email" autofocus :class="{ 'p-invalid': errors?.email }" />
                </FormRow>
              </div>
              <div>
                <FormRow label="Phone" :errors="errors?.phone" :required="true">
                    <TextInput v-model="form.phone" autofocus :class="{ 'p-invalid': errors?.phone }" />
                </FormRow>
              </div>
            </div>
            <div>
              <FormRow label="Address" :errors="errors?.address" :required="true">
                  <TextInput v-model="form.address" autofocus :class="{ 'p-invalid': errors?.address }" />
              </FormRow> 
            </div>
          </div>
        </div>

        <!-- Admin User (only for create) -->
        <div v-if="!isEditing">
          <div class="border-t border-gray-200 my-3"></div>
          <h3 class="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-1.5">
            <i class="pi pi-user text-indigo-400 text-sm"></i>
            Admin User Account
          </h3>
          <p class="text-sm text-gray-500 mb-3">A company_admin user will be created automatically.</p>
          <div>
            <FormRow label="Admin Name" :errors="errors?.admin_name" :required="true">
                <TextInput v-model="form.admin_name" autofocus :class="{ 'p-invalid': errors?.admin_name }" />
            </FormRow> 
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <FormRow label="Admin Email" :errors="errors?.admin_email" :required="true">
                  <TextInput v-model="form.admin_email" autofocus :class="{ 'p-invalid': errors?.admin_email }" />
              </FormRow> 
            </div>
            <div>
              <FormRow label="Password" :errors="errors?.admin_password" :required="true">
                  <PasswordInput v-model="form.admin_password" autofocus :class="{ 'p-invalid': errors?.admin_password }" />
              </FormRow> 
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ml-auto">
            <Button label="Cancel" severity="secondary" outlined size="small" @click="dialogVisible = false" />
            <Button :label="isEditing ? 'Update' : 'Create Company'" size="small" :loading="saving" :disabled="saving" @click="saveCompany" />
          </div>
        </div>
      </template>
    </Dialog>

    <!-- Confirm Dialog (suspend/activate/terminate) -->
    <Dialog v-model:visible="confirmVisible" :header="confirmTitle" modal :style="{ width: '400px' }">
      <p class="text-xs text-gray-600">{{ confirmMessage }}</p>
      <template #footer>
        <Button label="Cancel" severity="secondary" text size="small" @click="confirmVisible = false" />
        <Button :label="confirmActionLabel" :severity="confirmSeverity" size="small" :loading="confirming" :disabled="confirming" @click="executeConfirm" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
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

const toast = useToast()

// Data
const companies = ref([])
const searchQuery = ref('')
const statusFilter = ref(null)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const confirmVisible = ref(false)
const confirmAction = ref(null)
const confirmTarget = ref(null)
const saving = ref(false)
const confirming = ref(false)
const form = ref({ name: '', email: '', phone: '', address: '', admin_name: '', admin_email: '', admin_password: '' })

const statusOptions = [
  { label: 'All Status', value: null },
  { label: 'Active', value: 'active' },
  { label: 'Suspended', value: 'suspended' },
  { label: 'Terminated', value: 'terminated' }
]

const filteredCompanies = computed(() => {
  let result = companies.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c => c.name?.toLowerCase().includes(q) || c.email?.toLowerCase().includes(q))
  }
  if (statusFilter.value) {
    result = result.filter(c => c.status === statusFilter.value)
  }
  return result
})

// Load
onMounted(async () => {
  try {
    const res = await api.get('/api/v1/platform/companies')
    const payload = res.data
    companies.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load companies', life: 3000 })
  }
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

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', email: '', phone: '', address: '', admin_name: '', admin_email: '', admin_password: '' }
  dialogVisible.value = true
}

function openEdit(company) {
  isEditing.value = true
  editingId.value = company.id
  form.value = { name: company.name, email: company.email, phone: company.phone, address: company.address, admin_name: '', admin_email: '', admin_password: '' }
  dialogVisible.value = true
}

async function saveCompany() {
  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/platform/companies/${editingId.value}`, form.value)
      toast.add({ severity: 'success', summary: 'Updated', detail: 'Company updated', life: 2000 })
    } else {
      const res = await api.post('/api/v1/platform/companies', form.value)
      const admin = res.data?.data?.admin_user
      if (admin) {
        toast.add({
          severity: 'success',
          summary: 'Company Created',
          detail: `Company created with admin: ${admin.name} (${admin.email})`,
          life: 5000
        })
      } else {
        toast.add({ severity: 'success', summary: 'Created', detail: 'Company created', life: 2000 })
      }
    }
    dialogVisible.value = false
    const res = await api.get('/api/v1/platform/companies')
    const payload = res.data
    companies.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message || 'Operation failed', life: 3000 })
  } finally {
    saving.value = false
  }
}

// Confirm actions
const confirmTitle = computed(() => {
  switch (confirmAction.value) {
    case 'suspend': return 'Suspend Company'
    case 'activate': return 'Activate Company'
    case 'terminate': return 'Terminate Company'
    default: return ''
  }
})
const confirmMessage = computed(() => {
  if (!confirmTarget.value) return ''
  const name = confirmTarget.value.name
  switch (confirmAction.value) {
    case 'suspend': return `Deactivate database connection for "${name}". All tenant API access will be blocked.`
    case 'activate': return `Reactivate tenant database connection for "${name}". Access will be restored.`
    case 'terminate': return `⚠️ Permanently drop the tenant database for "${name}"! This action CANNOT be undone.`
    default: return ''
  }
})
const confirmActionLabel = computed(() => {
  switch (confirmAction.value) {
    case 'suspend': return 'Suspend'
    case 'activate': return 'Activate'
    case 'terminate': return 'Terminate'
    default: return ''
  }
})
const confirmSeverity = computed(() => {
  return confirmAction.value === 'terminate' ? 'danger' : 'warn'
})

function confirmSuspend(company) { confirmAction.value = 'suspend'; confirmTarget.value = company; confirmVisible.value = true }
function confirmActivate(company) { confirmAction.value = 'activate'; confirmTarget.value = company; confirmVisible.value = true }
function confirmTerminate(company) { confirmAction.value = 'terminate'; confirmTarget.value = company; confirmVisible.value = true }

async function executeConfirm() {
  if (!confirmTarget.value || !confirmAction.value) return
  confirming.value = true
  const id = confirmTarget.value.id
  try {
    await api.post(`/api/v1/platform/companies/${id}/${confirmAction.value}`)
    toast.add({ severity: 'success', summary: 'Success', detail: `Company ${confirmAction.value}ed`, life: 2000 })
    confirmVisible.value = false
    const res = await api.get('/api/v1/platform/companies')
    const payload = res.data
    companies.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message || 'Operation failed', life: 3000 })
  } finally {
    confirming.value = false
  }
}
</script>
