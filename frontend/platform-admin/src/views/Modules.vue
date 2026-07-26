<template>
  <div class="space-y-2">
    <div class="flex items-center justify-end">
      <Button label="Register Module" icon="pi pi-plus" size="small" @click="openCreate" />
    </div>
    <DataTable 
    :value="modules" 
    paginator 
    :rows="15" 
    size="small" 
    class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-cog text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">No modules registered yet</p>
          <p class="text-sm mt-1">Click "Register Module" to add the first module.</p>
        </div>
      </template>
      <Column field="name" header="Module" sortable>
        <template #body="{ data }">
          <div class="flex-row gap-2">
            <div class="uppercase font-semibold text-gray-600">{{ data.name }}</div>
            <div class="text-sm text-gray-500">{{ data.description }}</div>
          </div>
        </template>
      </Column>
      <Column field="slug" header="Slug" sortable />
      <Column field="version" header="Version" sortable />
      <Column field="created_at" header="Registered" sortable>
        <template #body="{ data }">
          <span class="text-gray-500">{{ data.created_at || '-' }}</span>
        </template>
      </Column>
      <Column header="Actions" :style="{ width: '120px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="'View companies'" @click="viewCompanies(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="'Edit'" @click="openEdit(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? 'Edit Module' : 'Register Module'" modal :style="{ width: '450px' }">
      <div class="space-y-3">
        <div>
          <FormRow label="Module Name" :errors="errors?.name" :required="true">
              <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
          </FormRow> 
        </div>
        <div>
          <FormRow label="Slug" :errors="errors?.slug" :required="true">
              <TextInput v-model="form.slug" autofocus :class="{ 'p-invalid': errors?.slug }" />
          </FormRow> 
        </div>
        <div>
          <FormRow label="Version" :errors="errors?.version" :required="true">
            <TextInput v-model="form.version" autofocus :class="{ 'p-invalid': errors?.version }" />
          </FormRow> 
        </div>
        <div>
          <FormRow label="Description" :errors="errors?.description" :required="true">
              <TextInput v-model="form.description" autofocus :class="{ 'p-invalid': errors?.description }" />
          </FormRow> 
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? 'Update' : 'Register'" size="small" :loading="saving" :disabled="saving" @click="saveModule" />
      </template>
    </Dialog>

    <!-- Company Assignment Dialog -->
    <Dialog v-model:visible="companyDialogVisible" :header="'Module: ' + (selectedModule?.name || '')" modal :style="{ width: '500px' }">
      <div class="space-y-2">
        <p class="text-sm text-gray-500 mb-2">Companies using this module:</p>
        <div v-for="c in moduleCompanies" :key="c.company_id" class="flex items-center justify-between px-3 py-2 rounded-md bg-gray-50 text-sm">
          <span>{{ c.company_name || c.company_id }}</span>
          <Tag :value="c.is_active ? 'Active' : 'Inactive'" :severity="c.is_active ? 'success' : 'warn'" class="!text-xs" />
        </div>
        <div v-if="moduleCompanies.length === 0" class="text-sm text-gray-400 text-center py-4">No companies using this module yet.</div>
      </div>
      <template #footer>
        <Button label="Close" severity="secondary" text size="small" @click="companyDialogVisible = false" />
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
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Card from 'primevue/card'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const toast = useToast()

const modules = ref([])
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const loadingCompanies = ref(false)
const form = ref({ name: '', slug: '', version: '1.0.0', description: '' })
const companyDialogVisible = ref(false)
const selectedModule = ref(null)
const moduleCompanies = ref([])

onMounted(async () => {
  try {
    const res = await api.get('/api/v1/platform/modules')
    const payload = res.data
    modules.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load modules', life: 3000 })
  }
})

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', slug: '', version: '1.0.0', description: '' }
  dialogVisible.value = true
}

function openEdit(mod) {
  isEditing.value = true
  editingId.value = mod.id
  form.value = { name: mod.name, slug: mod.slug, version: mod.version, description: mod.description || '' }
  dialogVisible.value = true
}

async function saveModule() {
  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/platform/modules/${editingId.value}`, form.value)
      toast.add({ severity: 'success', summary: 'Updated', life: 2000 })
    } else {
      await api.post('/api/v1/platform/modules', form.value)
      toast.add({ severity: 'success', summary: 'Registered', life: 2000 })
    }
    dialogVisible.value = false
    const res = await api.get('/api/v1/platform/modules')
    const payload = res.data
    modules.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: 'Error', detail: e.response?.data?.error?.message || 'Operation failed', life: 3000 })
  } finally {
    saving.value = false
  }
}

async function viewCompanies(mod) {
  selectedModule.value = mod
  loadingCompanies.value = true
  try {
    const res = await api.get(`/api/v1/platform/modules/${mod.id}/companies`)
    const payload = res.data
    moduleCompanies.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    moduleCompanies.value = []
  } finally {
    loadingCompanies.value = false
  }
  companyDialogVisible.value = true
}
</script>
