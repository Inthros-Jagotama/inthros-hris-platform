<template>
  <div class="space-y-2">
    <!-- Filters & Actions -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-1.5">
        <Button
          v-for="chip in roleFilterChips"
          :key="chip.value"
          :label="chip.label"
          :severity="roleFilter === chip.value ? (chip.severity || 'secondary') : 'secondary'"
          :outlined="roleFilter !== chip.value"
          size="small"
          class="!text-xs !px-2 !py-1"
          @click="roleFilter = chip.value"
        />
      </div>
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" />
        </IconField>
        <Button :label="t('users.add_user')" icon="pi pi-user-plus" size="small" @click="openCreate" />
      </div>
    </div>

    <!-- Bulk Action Toolbar -->
    <div
      v-if="selectedUsers.length > 0"
      class="flex items-center justify-between px-3 py-2 bg-indigo-50 dark:bg-indigo-900/20 border border-indigo-200 dark:border-indigo-800 rounded-lg text-sm"
    >
      <div class="flex items-center gap-2">
        <i class="pi pi-check-circle text-indigo-400 text-sm"></i>
        <span class="text-indigo-700 dark:text-indigo-300 font-medium">{{ selectedUsers.length }} {{ t('users.selected') }}</span>
        <Button
          icon="pi pi-times"
          size="small"
          text
          severity="secondary"
          class="!text-xs !w-5 !h-5"
          @click="selectedUsers = []"
          v-tooltip.right="t('common.clear_selection')"
        />
      </div>
      <div class="flex items-center gap-1.5">
        <Button
          :label="t('users.bulk_change_role')"
          icon="pi pi-lock"
          size="small"
          severity="info"
          class="!text-xs"
          @click="openBulkRoleDialog"
        />
        <Button
          :label="t('common.delete')"
          icon="pi pi-trash"
          size="small"
          severity="danger"
          class="!text-xs"
          @click="confirmBulkDelete"
        />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable
      v-else
      :value="filteredUsers"
      paginator
      :rows="15"
      v-model:selection="selectedUsers"
      dataKey="id"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <Column selectionMode="multiple" headerStyle="width: 40px" />
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('users.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('users.empty_hint') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('users.name')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>
      <Column field="email" :header="t('users.email')" sortable />
      <Column field="role" :header="t('users.role')" sortable>
        <template #body="{ data }">
          <Tag :value="data.role" :severity="data.role === 'super_admin' ? 'danger' : 'info'" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="company_name" :header="t('users.company')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <template v-if="data.company_name">
              <span class="text-gray-700 dark:text-gray-200">{{ data.company_name }}</span>
            </template>
            <span v-else class="text-gray-400 dark:text-gray-500 italic text-sm">—</span>
          </div>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable>
        <template #body="{ data }">
          <Tag :value="data.status || 'active'" severity="success" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('common.actions')" :style="{ width: '100px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEdit(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmSingleDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('users.edit_user') : t('users.add_user')" modal :style="{ width: '450px' }">
      <div class="space-y-3">
        <FormRow :label="t('users.name')" :errors="errors?.name" :required="true">
          <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('users.email')" :errors="errors?.email" :required="true">
          <TextInput v-model="form.email" autofocus :class="{ 'p-invalid': errors?.email }" />
        </FormRow>
        <FormRow :label="t('common.password')" :errors="errors?.password" :required="true">
          <PasswordInput v-model="form.password" autofocus :class="{ 'p-invalid': errors?.password }" />
        </FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('common.create')" size="small" :loading="saving" :disabled="saving" @click="saveUser" />
      </template>
    </Dialog>

    <!-- Bulk Change Role Dialog -->
    <Dialog v-model:visible="bulkRoleVisible" :header="t('users.bulk_change_role')" modal :style="{ width: '420px' }">
      <div class="space-y-3">
        <p class="text-sm text-gray-600 dark:text-gray-300">
          {{ t('users.bulk_change_role_message', { count: selectedUsers.length }) }}
        </p>
        <FormRow :label="t('users.role')" :required="true">
          <Select
            v-model="bulkRole"
            :options="roleOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('users.select_role')"
            class="!w-full"
          />
        </FormRow>
        <div class="text-xs text-gray-400 dark:text-gray-500">
          <span class="font-medium">{{ t('users.selected') }}:</span>
          <span class="ml-1">{{ selectedUserNames }}</span>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" text size="small" @click="bulkRoleVisible = false" />
        <Button :label="t('users.bulk_change_role')" severity="info" size="small" :loading="bulkRoleSaving" :disabled="!bulkRole || bulkRoleSaving" @click="executeBulkRoleChange" />
      </template>
    </Dialog>

    <!-- Confirm Delete Dialog -->
    <Dialog v-model:visible="confirmVisible" :header="confirmDeleteTitle" modal :style="{ width: '400px' }">
      <p class="text-sm text-gray-600 dark:text-gray-300">{{ confirmDeleteMessage }}</p>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" text size="small" @click="confirmVisible = false" />
        <Button :label="t('common.delete')" severity="danger" size="small" :loading="confirmDeleting" :disabled="confirmDeleting" @click="executeDelete" />
      </template>
    </Dialog>
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
import { useSkeletonPage } from '@/composables/useSkeletonPage'

const toast = useToast()
const { t } = useI18n()

const { loading, wrapLoad } = useSkeletonPage()
const users = ref([])
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ name: '', email: '', role: 'company_admin', password: '' })
const errors = ref({})

// Filters
const searchQuery = ref('')
const roleFilter = ref(null)

// Bulk selection
const selectedUsers = ref([])

// Bulk role dialog
const bulkRoleVisible = ref(false)
const bulkRole = ref(null)
const bulkRoleSaving = ref(false)

// Confirm delete
const confirmVisible = ref(false)
const confirmDeleteTarget = ref(null) // single delete target or 'bulk'
const confirmDeleting = ref(false)

const skeletonColumns = [
  { type: 'checkbox', headerWidth: 'w-6' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-32', headerWidth: 'w-28' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-14', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-14' }
]

// Role filter chips
const roleFilterChips = computed(() => [
  { label: t('common.all'), value: null, severity: 'info' },
  { label: t('users.super_admin'), value: 'super_admin', severity: 'danger' },
  { label: t('users.company_admin'), value: 'company_admin', severity: 'info' }
])

// Role options for dropdown
const roleOptions = computed(() => [
  { label: t('users.super_admin'), value: 'super_admin' },
  { label: t('users.company_admin'), value: 'company_admin' }
])

// Selected user names for display
const selectedUserNames = computed(() =>
  selectedUsers.value.map(u => u.name).join(', ')
)

// Confirm dialog texts
const confirmDeleteTitle = computed(() => {
  if (confirmDeleteTarget.value === 'bulk') return t('users.confirm_bulk_delete_title')
  if (confirmDeleteTarget.value) return t('users.confirm_delete_title')
  return ''
})

const confirmDeleteMessage = computed(() => {
  if (confirmDeleteTarget.value === 'bulk') {
    return t('users.confirm_bulk_delete_message', { count: selectedUsers.value.length })
  }
  if (confirmDeleteTarget.value) {
    return t('users.confirm_delete_message', { name: confirmDeleteTarget.value.name })
  }
  return ''
})

// Filtered users
const filteredUsers = computed(() => {
  let result = users.value
  if (roleFilter.value) {
    result = result.filter(u => u.role === roleFilter.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(u =>
      u.name?.toLowerCase().includes(q) ||
      u.email?.toLowerCase().includes(q)
    )
  }
  return result
})

onMounted(async () => {
  await loadUsers()
})

async function loadUsers() {
  try {
    await wrapLoad(async () => {
      const res = await api.get('/api/v1/platform/users')
      const payload = res.data
      users.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
    })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
  }
}

// ── CRUD ──
function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', email: '', role: 'company_admin', password: '' }
  errors.value = {}
  dialogVisible.value = true
}

function openEdit(user) {
  isEditing.value = true
  editingId.value = user.id
  form.value = { name: user.name, email: user.email, role: user.role, password: '' }
  errors.value = {}
  dialogVisible.value = true
}

async function saveUser() {
  saving.value = true
  try {
    if (isEditing.value) {
      const payload = { name: form.value.name, email: form.value.email, role: form.value.role }
      await api.put(`/api/v1/platform/users/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.updated'), life: 2000 })
    } else {
      await api.post('/api/v1/platform/users', form.value)
      toast.add({ severity: 'success', summary: t('message.created'), life: 2000 })
    }
    dialogVisible.value = false
    await loadUsers()
  } catch (e) {
    errors.value = getValidationErrors(e)
    if (Object.keys(errors.value).length === 0) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
    }
  } finally {
    saving.value = false
  }
}

// ── Bulk Role Change ──
function openBulkRoleDialog() {
  bulkRole.value = null
  bulkRoleVisible.value = true
}

async function executeBulkRoleChange() {
  if (!bulkRole.value || selectedUsers.value.length === 0) return
  bulkRoleSaving.value = true
  let success = 0
  let failed = 0
  for (const user of selectedUsers.value) {
    if (user.role === 'super_admin') {
      failed++
      continue
    }
    try {
      await api.put(`/api/v1/platform/users/${user.id}`, { role: bulkRole.value })
      success++
    } catch {
      failed++
    }
  }
  bulkRoleVisible.value = false
  bulkRoleSaving.value = false
  selectedUsers.value = []
  if (success > 0) {
    toast.add({ severity: 'success', summary: t('message.updated'), detail: `${success} user(s) role updated`, life: 3000 })
  }
  if (failed > 0) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: `${failed} user(s) could not be updated`, life: 3000 })
  }
  await loadUsers()
}

// ── Delete ──
function confirmSingleDelete(user) {
  confirmDeleteTarget.value = user
  confirmVisible.value = true
}

function confirmBulkDelete() {
  confirmDeleteTarget.value = 'bulk'
  confirmVisible.value = true
}

async function executeDelete() {
  confirmDeleting.value = true
  try {
    if (confirmDeleteTarget.value === 'bulk') {
      // Bulk delete
      let success = 0
      let failed = 0
      for (const user of selectedUsers.value) {
        if (user.role === 'super_admin') {
          failed++
          continue
        }
        try {
          await api.delete(`/api/v1/platform/users/${user.id}`)
          success++
        } catch {
          failed++
        }
      }
      selectedUsers.value = []
      if (success > 0) {
        toast.add({ severity: 'success', summary: t('message.deleted'), detail: `${success} user(s) deleted`, life: 3000 })
      }
      if (failed > 0) {
        toast.add({ severity: 'warn', summary: t('message.warning'), detail: `${failed} user(s) could not be deleted`, life: 3000 })
      }
    } else {
      // Single delete
      const user = confirmDeleteTarget.value
      await api.delete(`/api/v1/platform/users/${user.id}`)
      toast.add({ severity: 'success', summary: t('message.deleted'), detail: user.name, life: 2000 })
    }
    confirmVisible.value = false
    await loadUsers()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  } finally {
    confirmDeleting.value = false
    confirmDeleteTarget.value = null
  }
}
</script>
