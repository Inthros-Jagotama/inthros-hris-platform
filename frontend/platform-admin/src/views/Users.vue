<template>
  <div class="space-y-1">
    <div class="flex items-center justify-end">
      <Button label="Add User" icon="pi pi-user-plus" size="small" @click="openCreate" />
    </div>

    <DataTable :value="users" paginator :rows="15" size="small" class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">No platform users yet</p>
          <p class="text-sm mt-1">Click "Add User" to create the first platform user.</p>
        </div>
      </template>
      <Column field="name" header="Name" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>
      <Column field="email" header="Email" sortable />
      <Column field="role" header="Role" sortable>
        <template #body="{ data }">
          <Tag :value="data.role" :severity="data.role === 'super_admin' ? 'danger' : 'info'" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="company_name" header="Company" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <template v-if="data.company_name">
              <span class="text-gray-700">{{ data.company_name }}</span>
            </template>
            <span v-else class="text-gray-400 italic text-sm">—</span>
          </div>
        </template>
      </Column>
      <Column field="status" header="Status" sortable>
        <template #body="{ data }">
          <Tag :value="data.status || 'active'" severity="success" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column header="Actions" :style="{ width: '100px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="'Edit'" @click="openEdit(data)" />
            <Button icon="pi pi-lock" size="small" text severity="warning" v-tooltip.left="'Reset Password'" @click="confirmReset(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? 'Edit User' : 'Add Platform User'" modal :style="{ width: '450px' }">
      <div class="space-y-3">
        <div>
          <FormRow label="Name" :errors="errors?.name" :required="true">
              <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
          </FormRow> 
        </div>
        <div>
          <FormRow label="Email" :errors="errors?.email" :required="true">
              <TextInput v-model="form.email" autofocus :class="{ 'p-invalid': errors?.email }" />
          </FormRow> 
        </div>
        <div>
          <FormRow label="Password" :errors="errors?.password" :required="true">
              <PasswordInput v-model="form.password" autofocus :class="{ 'p-invalid': errors?.password }" />
          </FormRow> 
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? 'Update' : 'Create'" size="small" :loading="saving" :disabled="saving" @click="saveUser" />
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
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Avatar from 'primevue/avatar'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import PasswordInput from '@/components/PasswordInput.vue'

const toast = useToast()

const users = ref([])
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ name: '', email: '', role: 'company_admin', password: '' })
const roleOptions = [
  { label: 'Super Admin', value: 'super_admin' },
  { label: 'Company Admin', value: 'company_admin' }
]

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/platform/users')
    const payload = res.data
    users.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load users', life: 3000 })
  }
})

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', email: '', role: 'company_admin', password: '' }
  dialogVisible.value = true
}

function openEdit(user) {
  isEditing.value = true
  editingId.value = user.id
  form.value = { name: user.name, email: user.email, role: user.role, password: '' }
  dialogVisible.value = true
}

async function saveUser() {
  saving.value = true
  try {
    if (isEditing.value) {
      const payload = { name: form.value.name, email: form.value.email, role: form.value.role }
      await api.put(`/api/v1/platform/users/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: 'Updated', life: 2000 })
    } else {
      await api.post('/api/v1/platform/users', form.value)
      toast.add({ severity: 'success', summary: 'Created', life: 2000 })
    }
    dialogVisible.value = false
    const res = await api.get('/api/v1/platform/users')
    const payload = res.data
    users.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message || 'Operation failed', life: 3000 })
  } finally {
    saving.value = false
  }
}

function confirmReset(user) {
  // Placeholder — will implement password reset flow
  toast.add({ severity: 'info', summary: 'Coming Soon', detail: 'Password reset functionality', life: 2000 })
}
</script>
