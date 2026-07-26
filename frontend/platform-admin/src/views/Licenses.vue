<template>
  <div class="space-y-1">
    <div class="flex items-center justify-end">
      <Button label="New License" icon="pi pi-plus" size="small" @click="openCreate" />
    </div>

    <DataTable :value="licenses" paginator :rows="15" size="small" class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-id-card text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">No licenses issued yet</p>
          <p class="text-sm mt-1">Click "New License" to issue the first license.</p>
        </div>
      </template>
      <Column field="company_name" header="Company" sortable />
      <Column field="plan_type" header="Plan" sortable>
        <template #body="{ data }">
          <Tag :value="data.plan_type" :severity="planSeverity(data.plan_type)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="seat_count" header="Seats" sortable />
      <Column field="start_date" header="Valid From" sortable>
        <template #body="{ data }">{{ data.start_date || '-' }}</template>
      </Column>
      <Column field="end_date" header="Valid Until" sortable>
        <template #body="{ data }">{{ data.end_date || '-' }}</template>
      </Column>
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :value="data.status || 'active'" :severity="data.status === 'active' ? 'success' : 'warn'" class="!text-xs" />
        </template>
      </Column>
      <Column header="Actions" :style="{ width: '80px' }">
        <template #body="{ data }">
          <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="'Edit'" @click="openEdit(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? 'Edit License' : 'New License'" modal :style="{ width: '480px' }">
      <div class="space-y-3">
        <div class="grid grid-cols-2 gap-3">
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Company</label>
            <Select v-model="form.company_id" :options="companies" optionLabel="name" optionValue="id" placeholder="Select company" class="!w-full !h-8 !text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Plan Type</label>
            <Select v-model="form.plan_type" :options="planOptions" optionLabel="label" optionValue="value" placeholder="Select plan" class="!w-full !h-8 !text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Start Date</label>
            <InputText v-model="form.start_date" type="date" class="!w-full !text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">End Date</label>
            <InputText v-model="form.end_date" type="date" class="!w-full !text-sm" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-600 mb-1">Seats</label>
            <InputNumber v-model="form.seat_count" class="!w-full" inputClass="!w-full !text-sm" :min="1" />
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" text size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? 'Update' : 'Create'" size="small" :loading="saving" :disabled="saving" @click="saveLicense" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'

const toast = useToast()

const licenses = ref([])
const companies = ref([])
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ company_id: null, plan_type: 'trial', seat_count: 1, start_date: '', end_date: '' })
const planOptions = [
  { label: 'Trial', value: 'trial' },
  { label: 'Basic', value: 'basic' },
  { label: 'Professional', value: 'professional' },
  { label: 'Enterprise', value: 'enterprise' }
]

onMounted(async () => {
  try {
    const [licRes, compRes] = await Promise.all([
      api.get('/api/v1/platform/licenses'),
      api.get('/api/v1/platform/companies')
    ])
    const licPayload = licRes.data
    licenses.value = Array.isArray(licPayload.data) ? licPayload.data : (Array.isArray(licPayload) ? licPayload : [])
    const compPayload = compRes.data
    companies.value = Array.isArray(compPayload.data) ? compPayload.data : (Array.isArray(compPayload) ? compPayload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load data', life: 3000 })
  }
})

function planSeverity(plan) {
  switch (plan) {
    case 'enterprise': return 'danger'
    case 'professional': return 'warn'
    case 'basic': return 'info'
    default: return 'success'
  }
}

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { company_id: null, plan_type: 'trial', seat_count: 1, start_date: '', end_date: '' }
  dialogVisible.value = true
}

function openEdit(lic) {
  isEditing.value = true
  editingId.value = lic.id
  form.value = {
    company_id: lic.company_id,
    plan_type: lic.plan_type,
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
      toast.add({ severity: 'success', summary: 'Updated', life: 2000 })
    } else {
      await api.post('/api/v1/platform/licenses', form.value)
      toast.add({ severity: 'success', summary: 'Created', life: 2000 })
    }
    dialogVisible.value = false
    const res = await api.get('/api/v1/platform/licenses')
    const payload = res.data
    licenses.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message || 'Operation failed', life: 3000 })
  } finally {
    saving.value = false
  }
}
</script>
