<template>
  <div class="space-y-2">
    <!-- Header -->
    <div class="flex items-center justify-end">
      <Button :label="t('rbac.new_role')" icon="pi pi-plus" size="small" @click="openCreate" />
    </div>

    <DataTable :value="roles" paginator :rows="15" size="small" :loading="loading" class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-shield text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('rbac.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('rbac.empty_hint') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('rbac.role_name')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ data.name }}</span>
            <Tag v-if="data.is_system" value="SYSTEM" severity="contrast" class="!text-[10px] !px-1 !py-0" />
          </div>
        </template>
      </Column>
      <Column field="slug" :header="t('rbac.slug')" sortable>
        <template #body="{ data }">
          <code class="text-xs bg-gray-100 px-1.5 py-0.5 rounded font-mono">{{ data.slug }}</code>
        </template>
      </Column>
      <Column field="description" :header="t('common.description')" sortable>
        <template #body="{ data }">
          <span class="text-gray-500">{{ data.description || '—' }}</span>
        </template>
      </Column>
      <Column field="permissions" :header="t('rbac.permissions')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-1 flex-wrap max-w-[260px]">
            <Tag
              v-for="perm in (data.permissions || []).slice(0, 3)"
              :key="perm.id"
              :value="`${perm.resource}.${perm.action}`"
              severity="info"
              class="!text-[10px] !px-1 !py-0"
            />
            <span v-if="(data.permissions || []).length > 3" class="text-xs text-gray-400 ml-1">
              +{{ data.permissions.length - 3 }}
            </span>
            <span v-if="!data.permissions || data.permissions.length === 0" class="text-xs text-gray-400 italic">—</span>
          </div>
        </template>
      </Column>
      <Column :header="t('common.actions')" :style="{ width: '140px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEdit(data)" />
            <Button icon="pi pi-lock" size="small" text severity="info" v-tooltip.left="t('rbac.permissions')" @click="openPermissions(data)" />
            <Button 
              v-if="!data.is_system" 
              icon="pi pi-trash" 
              size="small" 
              text 
              severity="danger" 
              v-tooltip.left="t('common.delete')" 
              @click="confirmDelete(data)" 
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Role Dialog -->
    <Dialog v-model:visible="roleDialogVisible" :header="isEditing ? t('rbac.edit_role') : t('rbac.new_role')" modal :style="{ width: '480px' }">
      <div class="space-y-3">
        <FormRow :label="t('rbac.role_name')" :errors="errors?.name" :required="true">
          <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('rbac.slug')" :errors="errors?.slug" :required="true">
          <div class="relative slug-wrapper" :class="{ 'slug-highlight': slugHighlighted }">
            <TextInput v-model="form.slug" :class="{ 'p-invalid': errors?.slug }" @input="slugManuallyEdited = true" />
            <i
              v-if="!slugManuallyEdited && form.name"
              v-tooltip.left="t('common.auto_generated')"
              class="pi pi-sync text-[10px] absolute right-2 top-1/2 -translate-y-1/2 transition-colors duration-300"
              :class="slugHighlighted ? 'text-indigo-400' : 'text-gray-300'"
            ></i>
          </div>
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="3" :class="{ 'p-invalid': errors?.description }" />
        </FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="roleDialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('common.create')" size="small" :loading="saving" :disabled="saving" @click="saveRole" />
      </template>
    </Dialog>

    <!-- Permission Assignment Dialog -->
    <Dialog v-model:visible="permDialogVisible" :header="t('rbac.assign_permissions', { name: selectedRole?.name })" modal :style="{ width: '800px' }">
      <div v-if="allPermissions.length === 0" class="text-sm text-gray-400 text-center py-6">
        {{ t('rbac.no_permissions') }}
      </div>
      <div v-else class="max-h-[520px] overflow-y-auto space-y-3 pr-1">
        <div v-for="group in groupedPermissions" :key="group.resource" class="rounded-lg border border-gray-200 overflow-hidden">
          <!-- Module Header -->
          <div class="flex items-center justify-between px-4 py-2.5 bg-gray-50 border-b border-gray-200">
            <div class="flex items-center gap-2.5">
              <Tag :value="group.resource" severity="contrast" class="!text-[11px] !px-2 !py-0.5 !font-semibold uppercase tracking-wider" />
              <div class="flex items-center gap-1.5">
                <span class="text-xs text-gray-500">{{ group.permissions.length }} permission(s)</span>
              </div>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="text-xs text-gray-400" v-if="selectedCountPerGroup[group.resource]">
                {{ selectedCountPerGroup[group.resource] }}/{{ group.permissions.length }}
              </span>
              <Button
                :label="t('common.select_all')"
                size="small"
                text
                severity="secondary"
                class="!text-[10px] !px-1.5 !py-0.5"
                @click.stop="selectGroupPermissions(group.resource)"
                v-tooltip.left="'Select all ' + group.resource + ' permissions'"
              />
            </div>
          </div>
          <!-- Permission Rows -->
          <div class="divide-y divide-gray-100">
            <div
              v-for="perm in group.permissions"
              :key="perm.id"
              class="flex items-center justify-between px-4 py-2.5 hover:bg-indigo-50/20 transition-colors"
              :class="{ 'bg-indigo-50/40': rolePermissionIds.has(perm.id) }">
              <div class="flex items-center gap-3 min-w-0">
                <Tag :value="perm.action" severity="info" class="!text-[10px] !px-1.5 !py-0 min-w-[48px] text-center" />
                <div class="min-w-0">
                  <p v-if="perm.description" class="text-xs text-gray-500 truncate">{{ perm.description }}</p>
                  <p v-else class="text-xs text-gray-300 italic">{{ t('common.no_data') }}</p>
                </div>
              </div>
              <ToggleSwitch
                :model-value="rolePermissionIds.has(perm.id)"
                @update:model-value="setPermission(perm.id, $event)"
                class="shrink-0"
              />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center gap-2 justify-between w-full">
          <div class="flex items-center gap-2">
            <Button 
              :label="t('common.select_all')" 
              size="small" 
              text 
              severity="secondary" 
              class="!text-xs" 
              @click="selectAllPermissions" 
            />
            <Button 
              :label="t('common.deselect_all')" 
              size="small" 
              text 
              severity="secondary" 
              class="!text-xs" 
              @click="deselectAllPermissions" 
            />
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs text-gray-500">{{ rolePermissionIds.size }} / {{ allPermissions.length }}</span>
            <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="permDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="savingPerms" :disabled="savingPerms" @click="savePermissions" />
          </div>
        </div>
      </template>
    </Dialog>

    <!-- Confirm Delete Dialog -->
    <Dialog v-model:visible="confirmVisible" :header="t('common.delete')" modal :style="{ width: '400px' }">
      <p class="text-sm text-gray-600">{{ t('rbac.confirm_delete_message', { name: confirmTarget?.name }) }}</p>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="confirmVisible = false" />
        <Button :label="t('common.delete')" severity="danger" size="small" :loading="confirming" :disabled="confirming" @click="executeDelete" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useSlugify } from '@/composables/useSlugify'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const toast = useToast()
const { t } = useI18n()

// Data
const roles = ref([])
const loading = ref(true)
const allPermissions = ref([])
const roleDialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ name: '', slug: '', description: '' })
const errors = ref({})

// Auto-slug
const { slugManuallyEdited, slugHighlighted, resetSlug, disableAutoSlug } = useSlugify(
  () => form.value.name,
  (v) => { form.value.slug = v }
)

// Permissions dialog
const permDialogVisible = ref(false)
const selectedRole = ref(null)
const rolePermissionIds = reactive(new Set())
const savingPerms = ref(false)

// Group permissions by resource/module
const groupedPermissions = computed(() => {
  const groups = {}
  const perms = allPermissions.value
  // Sort: rbac → company → user → module → license → ... alphabetically
  const sortOrder = ['rbac', 'company', 'user', 'module', 'license', 'monitoring', 'package',
    'organization', 'employee', 'attendance', 'leave', 'payroll',
    'competency', 'jobmanagement', 'employeemovement', 'approval']
  
  perms.forEach(p => {
    if (!groups[p.resource]) {
      groups[p.resource] = { resource: p.resource, permissions: [], description: p.description || '' }
    }
    groups[p.resource].permissions.push(p)
  })
  
  // Convert to array and sort
  return Object.values(groups).sort((a, b) => {
    const ai = sortOrder.indexOf(a.resource)
    const bi = sortOrder.indexOf(b.resource)
    return (ai === -1 ? 999 : ai) - (bi === -1 ? 999 : bi)
  })
})

// Count selected per group
const selectedCountPerGroup = computed(() => {
  const counts = {}
  allPermissions.value.forEach(p => {
    if (rolePermissionIds.has(p.id)) {
      counts[p.resource] = (counts[p.resource] || 0) + 1
    }
  })
  return counts
})

// Confirm dialog
const confirmVisible = ref(false)
const confirmTarget = ref(null)
const confirming = ref(false)

// Load data
async function loadData() {
  loading.value = true
  try {
    const [roleRes, permRes] = await Promise.all([
      api.get('/api/v1/platform/rbac/roles'),
      api.get('/api/v1/platform/rbac/permissions')
    ])
    const rolePayload = roleRes.data
    roles.value = Array.isArray(rolePayload.data) ? rolePayload.data : (Array.isArray(rolePayload) ? rolePayload : [])
    const permPayload = permRes.data
    allPermissions.value = Array.isArray(permPayload.data) ? permPayload.data : (Array.isArray(permPayload) ? permPayload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadData)

// Role CRUD
function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', slug: '', description: '' }
  errors.value = {}
  resetSlug()
  roleDialogVisible.value = true
}

function openEdit(role) {
  isEditing.value = true
  editingId.value = role.id
  form.value = { name: role.name, slug: role.slug, description: role.description || '' }
  errors.value = {}
  disableAutoSlug()
  roleDialogVisible.value = true
}

async function saveRole() {
  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/platform/rbac/roles/${editingId.value}`, form.value)
      toast.add({ severity: 'success', summary: t('message.updated'), life: 2000 })
    } else {
      await api.post('/api/v1/platform/rbac/roles', form.value)
      toast.add({ severity: 'success', summary: t('message.created'), life: 2000 })
    }
    roleDialogVisible.value = false
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

// Permission management
function openPermissions(role) {
  selectedRole.value = role
  rolePermissionIds.clear()
  if (role.permissions) {
    role.permissions.forEach(p => rolePermissionIds.add(p.id))
  }
  permDialogVisible.value = true
}

function setPermission(id, value) {
  if (value) {
    rolePermissionIds.add(id)
  } else {
    rolePermissionIds.delete(id)
  }
}

function selectAllPermissions() {
  rolePermissionIds.clear()
  allPermissions.value.forEach(p => rolePermissionIds.add(p.id))
}

function deselectAllPermissions() {
  rolePermissionIds.clear()
}

function selectGroupPermissions(resource) {
  const groupPerms = allPermissions.value.filter(p => p.resource === resource)
  const allSelected = groupPerms.every(p => rolePermissionIds.has(p.id))
  if (allSelected) {
    groupPerms.forEach(p => rolePermissionIds.delete(p.id))
  } else {
    groupPerms.forEach(p => rolePermissionIds.add(p.id))
  }
}

async function savePermissions() {
  if (!selectedRole.value) return
  savingPerms.value = true
  try {
    const roleId = selectedRole.value.id
    const currentPerms = selectedRole.value.permissions || []
    const currentIds = new Set(currentPerms.map(p => p.id))

    // Batch revoke & assign via Promise.all
    const toRevoke = currentPerms.filter(p => !rolePermissionIds.has(p.id))
    const toAssign = allPermissions.value.filter(p => rolePermissionIds.has(p.id) && !currentIds.has(p.id))

    await Promise.all([
      ...toRevoke.map(p => api.delete(`/api/v1/platform/rbac/roles/${roleId}/permissions/${p.id}`)),
      ...toAssign.map(p => api.post(`/api/v1/platform/rbac/roles/${roleId}/permissions`, { permission_id: p.id }))
    ])

    toast.add({ severity: 'success', summary: t('message.saved'), life: 2000 })
    permDialogVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  } finally {
    savingPerms.value = false
  }
}

// Delete role
function confirmDelete(role) {
  if (role.is_system) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('rbac.cannot_delete_system'), life: 3000 })
    return
  }
  confirmTarget.value = role
  confirmVisible.value = true
}

async function executeDelete() {
  if (!confirmTarget.value) return
  confirming.value = true
  try {
    await api.delete(`/api/v1/platform/rbac/roles/${confirmTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.deleted'), life: 2000 })
    confirmVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  } finally {
    confirming.value = false
  }
}
</script>
