<template>
  <div class="space-y-4">
    <!-- Tab switcher: Roles | Users -->
    <div class="flex items-center gap-1.5">
      <Button
        v-for="tab in tabs"
        :key="tab.value"
        :label="tab.label"
        :severity="activeTab === tab.value ? 'info' : 'secondary'"
        :outlined="activeTab !== tab.value"
        size="small"
        class="!text-xs !px-3 !py-1.5"
        @click="switchTab(tab.value)"
      />
    </div>

    <!-- ═══ ROLES TAB ═══ -->
    <template v-if="activeTab === 'roles'">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <span v-if="roles.length > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ roles.length }} {{ t('common.items') }}</span>
        </div>
        <Button :label="t('rbac.new_role')" icon="pi pi-plus" size="small" @click="openRoleDialog()" />
      </div>

      <SkeletonTable v-if="loadingRoles" :columns="roleSkeletonColumns" :rows="6" />

      <DataTable
        v-else
        :value="roles"
        size="small"
        class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        sortField="name"
        :sortOrder="1"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-shield text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('rbac.empty_roles') }}</p>
          </div>
        </template>
        <Column field="name" :header="t('rbac.role_name')" sortable>
          <template #body="{ data }">
            <div class="flex items-center gap-2">
              <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
              <Tag v-if="data.is_system" :value="t('rbac.system')" severity="secondary" class="!text-[10px] !px-1.5 !py-0.5" />
              <Tag v-if="data.is_default" :value="t('rbac.default')" severity="info" class="!text-[10px] !px-1.5 !py-0.5" />
            </div>
          </template>
        </Column>
        <Column field="description" :header="t('common.description')">
          <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400">{{ data.description || '—' }}</span></template>
        </Column>
        <Column field="user_count" :header="t('rbac.users_count')" sortable style="width:110px">
          <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400">{{ data.user_count }}</span></template>
        </Column>
        <Column field="permission_count" :header="t('rbac.permissions_count')" sortable style="width:130px">
          <template #body="{ data }">
            <span class="text-gray-500 dark:text-gray-400">{{ data.permission_ids?.length || 0 }}</span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width:130px" frozen alignFrozen="right">
          <template #body="{ data }">
            <div class="flex items-center gap-1">
              <Button icon="pi pi-shield" size="small" text severity="info" v-tooltip.left="t('rbac.assign_permissions')" @click="goPermissions(data)" />
              <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" :disabled="data.is_system" @click="openRoleDialog(data)" />
              <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" :disabled="data.is_system" @click="confirmDeleteRole(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </template>

    <!-- ═══ USERS TAB ═══ -->
    <template v-else>
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <span v-if="users.length > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ users.length }} {{ t('common.items') }}</span>
        </div>
      </div>

      <SkeletonTable v-if="loadingUsers" :columns="userSkeletonColumns" :rows="6" />

      <DataTable
        v-else
        :value="users"
        size="small"
        class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        sortField="name"
        :sortOrder="1"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('rbac.empty_users') }}</p>
          </div>
        </template>
        <Column field="name" :header="t('common.name')" sortable>
          <template #body="{ data }"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
        </Column>
        <Column field="email" :header="t('common.email')" sortable>
          <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400">{{ data.email }}</span></template>
        </Column>
        <Column field="is_active" :header="t('common.status')" sortable style="width:100px">
          <template #body="{ data }">
            <Tag :value="data.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="data.is_active ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5" />
          </template>
        </Column>
        <Column field="role_names" :header="t('rbac.roles')" style="width:220px">
          <template #body="{ data }">
            <div class="flex flex-wrap gap-1">
              <Tag v-for="r in (data.role_names || [])" :key="r" :value="r" severity="info" class="!text-[10px] !px-1.5 !py-0.5" />
              <span v-if="!data.role_names?.length" class="text-gray-400">—</span>
            </div>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
          <template #body="{ data }">
            <Button icon="pi pi-users" size="small" text severity="info" v-tooltip.left="t('rbac.assign_roles')" @click="openUserRoleDialog(data)" />
          </template>
        </Column>
      </DataTable>
    </template>

    <!-- ═══ ROLE CREATE/EDIT DIALOG ═══ -->
    <Dialog v-model:visible="roleDialogVisible" :header="editingRole ? t('rbac.edit_role') : t('rbac.new_role')" modal :style="{ width: '480px' }" :closable="true" @hide="resetRoleForm">
      <div class="space-y-3">
        <FormRow :label="t('rbac.role_name')" required :errors="roleErrors?.name">
          <TextInput v-model="roleForm.name" maxlength="255" autofocus :placeholder="t('rbac.role_name')" :class="{ 'p-invalid': roleErrors?.name }" :disabled="editingRole && roleForm.is_system" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="roleErrors?.description">
          <TextInput v-model="roleForm.description" maxlength="255" :placeholder="t('common.description')" />
        </FormRow>
        <div class="flex items-center justify-between pt-1">
          <label class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('rbac.is_default') }}</label>
          <ToggleSwitch v-model="roleForm.is_default" />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="roleDialogVisible=false" />
          <Button :label="editingRole ? t('common.update') : t('common.save')" size="small" :loading="savingRole" :disabled="savingRole" @click="handleSaveRole" />
        </div>
      </template>
    </Dialog>

    <!-- ═══ ASSIGN USER ROLES DIALOG ═══ -->
    <Dialog v-model:visible="userRoleDialogVisible" :header="`${t('rbac.assign_roles')} — ${userRoleTarget?.name || ''}`" modal :style="{ width: '480px' }" :closable="true" @hide="userRoleForm.role_ids = []">
      <div class="space-y-3">
        <FormRow :label="t('rbac.roles')" required :errors="userRoleErrors?.role_ids">
          <MultiSelect
            v-model="userRoleForm.role_ids"
            :options="roles"
            optionLabel="name"
            optionValue="id"
            :placeholder="t('rbac.select_roles')"
            :maxSelectedLabels="3"
            class="!w-full"
            display="chip"
            :class="{ 'p-invalid': userRoleErrors?.role_ids }"
          />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="userRoleDialogVisible=false" />
          <Button :label="t('common.save')" size="small" :loading="savingUserRoles" :disabled="savingUserRoles" @click="handleSaveUserRoles" />
        </div>
      </template>
    </Dialog>

    <!-- Delete confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('rbac.confirm_delete_title')"
      :message="t('rbac.confirm_delete', { name: deleteTarget?.name || '' })"
      :loading="deletingRole"
      :error-msg="deleteError"
      :confirm-label="t('common.delete')"
      :cancel-label="t('common.cancel')"
      @confirm="handleDeleteRole"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import MultiSelect from 'primevue/multiselect'
import ToggleSwitch from 'primevue/toggleswitch'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const activeTab = ref('roles')
const tabs = computed(() => [
  { label: t('rbac.tab_roles'), value: 'roles' },
  { label: t('rbac.tab_users'), value: 'users' }
])

// ── Roles state ──
const roles = ref([])
const loadingRoles = ref(false)
const roleDialogVisible = ref(false)
const editingRole = ref(false)
const editingRoleId = ref(null)
const savingRole = ref(false)
const roleErrors = ref({})
const roleForm = ref({ name: '', description: '', is_default: false, is_system: false })

const roleSkeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-56', headerWidth: 'w-24' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

// ── Users state ──
const users = ref([])
const loadingUsers = ref(false)

const userSkeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-56', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-40', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

// ── User-role state ──
const userRoleDialogVisible = ref(false)
const userRoleTarget = ref(null)
const userRoleForm = ref({ role_ids: [] })
const userRoleErrors = ref({})
const savingUserRoles = ref(false)

// ── Delete state ──
const deleteDialogVisible = ref(false)
const deletingRole = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

// ── Loaders ──
async function loadRoles() {
  loadingRoles.value = true
  try {
    const res = await api.get('/api/v1/tenant/rbac/roles')
    roles.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loadingRoles.value = false
  }
}

async function loadUsers() {
  loadingUsers.value = true
  try {
    const res = await api.get('/api/v1/tenant/rbac/users')
    users.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loadingUsers.value = false
  }
}

function switchTab(tab) {
  activeTab.value = tab
  if (tab === 'roles' && roles.value.length === 0 && !loadingRoles.value) loadRoles()
  if (tab === 'users' && users.value.length === 0 && !loadingUsers.value) loadUsers()
}

// ── Role dialog ──
function openRoleDialog(item) {
  editingRole.value = !!item
  editingRoleId.value = item?.id || null
  roleErrors.value = {}
  roleForm.value = {
    name: item?.name || '',
    description: item?.description || '',
    is_default: item?.is_default || false,
    is_system: item?.is_system || false
  }
  roleDialogVisible.value = true
}

function resetRoleForm() {
  roleForm.value = { name: '', description: '', is_default: false, is_system: false }
  roleErrors.value = {}
  editingRole.value = false
  editingRoleId.value = null
}

async function handleSaveRole() {
  roleErrors.value = {}
  if (!roleForm.value.name?.trim()) { roleErrors.value = { name: [t('form.required')] }; return }
  savingRole.value = true
  try {
    const payload = {
      name: roleForm.value.name.trim(),
      description: roleForm.value.description || undefined,
      is_default: roleForm.value.is_default
    }
    if (editingRole.value) {
      await api.put(`/api/v1/tenant/rbac/roles/${editingRoleId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('rbac.role_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/rbac/roles', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('rbac.role_created'), life: 3000 })
    }
    roleDialogVisible.value = false
    await loadRoles()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      roleErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    savingRole.value = false
  }
}

// ── Permission assignment → halaman khusus ──
function goPermissions(role) {
  router.push(`/settings/rbac/roles/${role.id}/permissions`)
}

// ── User-role assignment ──
function openUserRoleDialog(user) {
  userRoleTarget.value = user
  userRoleErrors.value = {}
  userRoleForm.value = { role_ids: [...(user.role_ids || [])] }
  if (roles.value.length === 0) loadRoles()
  userRoleDialogVisible.value = true
}

async function handleSaveUserRoles() {
  userRoleErrors.value = {}
  if (!userRoleForm.value.role_ids?.length) { userRoleErrors.value = { role_ids: [t('form.required')] }; return }
  savingUserRoles.value = true
  try {
    await api.put(`/api/v1/tenant/rbac/users/${userRoleTarget.value.id}/roles`, { role_ids: userRoleForm.value.role_ids })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('rbac.user_roles_updated'), life: 3000 })
    userRoleDialogVisible.value = false
    await loadUsers()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      userRoleErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    savingUserRoles.value = false
  }
}

// ── Delete role ──
function confirmDeleteRole(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDeleteRole() {
  deletingRole.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/rbac/roles/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('rbac.role_deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadRoles()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deletingRole.value = false
  }
}

onMounted(() => {
  loadRoles()
  loadUsers()
})
</script>
