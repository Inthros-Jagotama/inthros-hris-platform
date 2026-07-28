<template>
  <div class="space-y-4">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ activeTab === 0 ? t('organization.title') : t('zones.title') }}</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ activeTab === 0 ? t('organization.description') : t('zones.description') }}</p>
      </div>
      <div class="flex items-center gap-2">
        <Button
          v-if="activeTab === 0"
          :label="t('organization.add_root')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          outlined
          @click="openCreate(null)"
        />
        <Button
          v-if="activeTab === 1"
          :label="t('zones.new_zone')"
          icon="pi pi-plus"
          size="small"
          severity="secondary"
          outlined
          @click="openZoneDialog()"
        />
        <Button
          :label="t('common.refresh')"
          icon="pi pi-refresh"
          size="small"
          severity="secondary"
          text
          @click="activeTab === 0 ? loadTree() : loadZones()"
          :loading="loading"
        />
      </div>
    </div>

    <!-- Tabs -->
    <TabView v-model:activeIndex="activeTab" class="!text-sm">
      <!-- Tab 1: Organization Tree -->
      <TabPanel :header="t('organization.tree_view')">
        <div class="pt-2">
          <!-- Skeleton -->
          <div v-if="loading" class="space-y-2">
            <div v-for="i in 5" :key="i" class="flex items-center gap-3 py-1">
              <Skeleton shape="rectangle" width="1.25rem" height="1.25rem" class="!rounded" />
              <Skeleton width="8rem" height="1rem" />
              <Skeleton width="12rem" height="1rem" />
            </div>
          </div>

          <!-- Empty State -->
          <div
            v-else-if="treeData.length === 0"
            class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500"
          >
            <i class="pi pi-sitemap text-4xl mb-3 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('organization.empty_title') }}</p>
            <p class="text-sm mt-1 mb-4">{{ t('organization.empty_tree') }}</p>
            <Button
              :label="t('organization.add_root')"
              icon="pi pi-plus"
              size="small"
              @click="openCreate(null)"
            />
          </div>

          <!-- TreeTable -->
          <TreeTable
            v-else
            :value="treeData"
            class="!text-sm !border-0"
            :scrollable="true"
            scrollHeight="flex"
            stripedRows
            selectionMode="single"
            v-model:selectionKeys="selectedNodeKey"
          >
            <Column field="nomenclature" :header="t('organization.nomenclature')" :expander="true">
              <template #body="{ node }">
                <div class="flex items-center gap-2">
                  <i class="pi pi-folder-open text-amber-500 text-xs"></i>
                  <span class="font-medium text-gray-800 dark:text-gray-100">{{ node.data.nomenclature }}</span>
                </div>
              </template>
            </Column>
            <Column field="code" :header="t('organization.code')" style="width: 120px">
              <template #body="{ node }">
                <Tag :value="node.data.code" severity="info" class="!text-xs" />
              </template>
            </Column>
            <Column field="full_code" :header="t('organization.full_code')" style="width: 160px">
              <template #body="{ node }">
                <span class="text-gray-500 dark:text-gray-400 text-xs font-mono">{{ node.data.full_code }}</span>
              </template>
            </Column>
            <Column field="level" :header="t('organization.level')" style="width: 80px">
              <template #body="{ node }">
                <span class="text-gray-500 dark:text-gray-400">{{ node.data.level }}</span>
              </template>
            </Column>
            <Column field="sort_order" :header="t('organization.sort_order')" style="width: 90px">
              <template #body="{ node }">
                <span class="text-gray-500 dark:text-gray-400">{{ node.data.sort_order }}</span>
              </template>
            </Column>
            <Column :header="t('common.actions')" style="width: 140px" frozen alignFrozen="right">
              <template #body="{ node }">
                <div class="flex items-center gap-1">
                  <Button
                    icon="pi pi-plus"
                    v-tooltip.top="t('organization.add_child')"
                    severity="secondary"
                    text
                    size="small"
                    class="!p-1"
                    @click="openCreate(node.data)"
                  />
                  <Button
                    icon="pi pi-pencil"
                    v-tooltip.top="t('common.edit')"
                    severity="secondary"
                    text
                    size="small"
                    class="!p-1"
                    @click="openEdit(node.data)"
                  />
                  <Button
                    icon="pi pi-trash"
                    v-tooltip.top="t('common.delete')"
                    severity="danger"
                    text
                    size="small"
                    class="!p-1"
                    @click="confirmDelete(node.data)"
                  />
                </div>
              </template>
            </Column>
          </TreeTable>
        </div>
      </TabPanel>

      <!-- Tab 2: Zones -->
      <TabPanel :header="t('zones.title')">
        <div class="pt-2">
          <!-- Skeleton -->
          <div v-if="zonesLoading" class="space-y-2">
            <div v-for="i in 4" :key="i" class="flex items-center gap-4 py-2">
              <Skeleton width="5rem" height="1rem" />
              <Skeleton width="10rem" height="1rem" />
              <Skeleton width="6rem" height="1rem" />
              <Skeleton width="4rem" height="1.25rem" />
            </div>
          </div>

          <!-- Empty State -->
          <div
            v-else-if="zones.length === 0"
            class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500"
          >
            <i class="pi pi-map-marker text-4xl mb-3 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('zones.empty_title') }}</p>
            <p class="text-sm mt-1 mb-4">{{ t('zones.description') }}</p>
            <Button
              :label="t('zones.new_zone')"
              icon="pi pi-plus"
              size="small"
              @click="openZoneDialog()"
            />
          </div>

          <!-- Zones DataTable -->
          <DataTable
            v-else
            :value="zones"
            class="!text-sm"
            stripedRows
            :loading="zonesLoading"
            paginator
            :rows="20"
            :totalRecords="zonesTotal"
            :lazy="true"
            @page="onZonePage"
          >
            <Column field="code" :header="t('zones.code')" style="width: 120px">
              <template #body="{ data }">
                <Tag :value="data.code" severity="info" class="!text-xs" />
              </template>
            </Column>
            <Column field="name" :header="t('zones.name')">
              <template #body="{ data }">
                <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
              </template>
            </Column>
            <Column field="region" :header="t('zones.region')" style="width: 150px">
              <template #body="{ data }">
                <span class="text-gray-500 dark:text-gray-400">{{ data.region || '—' }}</span>
              </template>
            </Column>
            <Column field="is_active" :header="t('zones.is_active')" style="width: 100px">
              <template #body="{ data }">
                <Tag
                  :value="data.is_active ? t('common_status.active') : t('common_status.inactive')"
                  :severity="data.is_active ? 'success' : 'warn'"
                  class="!text-xs"
                />
              </template>
            </Column>
            <Column field="sort_order" :header="t('zones.sort_order')" style="width: 100px">
              <template #body="{ data }">
                <span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span>
              </template>
            </Column>
            <Column :header="t('common.actions')" style="width: 100px" frozen alignFrozen="right">
              <template #body="{ data }">
                <div class="flex items-center gap-1">
                  <Button
                    icon="pi pi-pencil"
                    v-tooltip.top="t('common.edit')"
                    severity="secondary"
                    text
                    size="small"
                    class="!p-1"
                    @click="openZoneDialog(data)"
                  />
                  <Button
                    icon="pi pi-trash"
                    v-tooltip.top="t('common.delete')"
                    severity="danger"
                    text
                    size="small"
                    class="!p-1"
                    @click="confirmDeleteZone(data)"
                  />
                </div>
              </template>
            </Column>
          </DataTable>
        </div>
      </TabPanel>
    </TabView>

    <!-- Organization: Create / Edit Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :header="isEditing ? t('organization.edit') : t('organization.create')"
      :modal="true"
      :closable="true"
      class="!w-full !max-w-lg"
      @hide="resetForm"
    >
      <div class="space-y-4">
        <div v-if="form.parent_id" class="bg-emerald-50 dark:bg-emerald-900/20 border border-emerald-200 dark:border-emerald-800 rounded-md px-3 py-2 text-sm text-emerald-700 dark:text-emerald-300">
          <i class="pi pi-arrow-right mr-1"></i>
          {{ t('organization.parent') }}: <strong>{{ parentLabel }}</strong>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('organization.code') }} <span class="text-red-500">*</span></label>
          <InputText v-model="form.code" class="!w-full" :class="{ 'p-invalid': errors?.code }" maxlength="10" :placeholder="t('organization.code')" />
          <small v-if="errors?.code" class="text-red-500 text-xs mt-1 block">{{ errors.code }}</small>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('organization.nomenclature') }} <span class="text-red-500">*</span></label>
          <InputText v-model="form.nomenclature" class="!w-full" :class="{ 'p-invalid': errors?.nomenclature }" maxlength="255" :placeholder="t('organization.nomenclature')" />
          <small v-if="errors?.nomenclature" class="text-red-500 text-xs mt-1 block">{{ errors.nomenclature }}</small>
        </div>
        <div v-if="!isEditing">
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('organization.parent') }}</label>
          <Select v-model="form.parent_id" :options="parentOptions" optionValue="id" optionLabel="label" :placeholder="t('organization.select_parent')" class="!w-full" :showClear="true" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('organization.sort_order') }}</label>
          <InputNumber v-model="form.sort_order" class="!w-full" :min="0" />
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
      </template>
    </Dialog>

    <!-- Zones: Create / Edit Dialog -->
    <Dialog
      v-model:visible="zoneDialogVisible"
      :header="zoneEditing ? t('zones.edit_zone') : t('zones.new_zone')"
      :modal="true"
      :closable="true"
      class="!w-full !max-w-md"
      @hide="resetZoneForm"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('zones.code') }} <span class="text-red-500">*</span></label>
          <InputText v-model="zoneForm.code" class="!w-full" :class="{ 'p-invalid': zoneErrors?.code }" maxlength="20" :placeholder="t('zones.code')" />
          <small v-if="zoneErrors?.code" class="text-red-500 text-xs mt-1 block">{{ zoneErrors.code }}</small>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('zones.name') }} <span class="text-red-500">*</span></label>
          <InputText v-model="zoneForm.name" class="!w-full" :class="{ 'p-invalid': zoneErrors?.name }" maxlength="255" :placeholder="t('zones.name')" />
          <small v-if="zoneErrors?.name" class="text-red-500 text-xs mt-1 block">{{ zoneErrors.name }}</small>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('zones.region') }}</label>
          <InputText v-model="zoneForm.region" class="!w-full" maxlength="100" :placeholder="t('zones.region')" />
        </div>
        <div class="flex items-center justify-between">
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('zones.is_active') }}</label>
          <ToggleSwitch v-model="zoneForm.is_active" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('zones.sort_order') }}</label>
          <InputNumber v-model="zoneForm.sort_order" class="!w-full" :min="0" />
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="zoneDialogVisible = false" />
        <Button :label="zoneEditing ? t('common.update') : t('common.save')" size="small" :loading="zoneSaving" :disabled="zoneSaving" @click="handleZoneSave" />
      </template>
    </Dialog>

    <!-- Delete Confirm -->
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import TreeTable from 'primevue/treetable'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import ToggleSwitch from 'primevue/toggleswitch'
import Skeleton from 'primevue/skeleton'
import ConfirmDialog from 'primevue/confirmdialog'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

// ── Organization Tree ──
const loading = ref(false)
const saving = ref(false)
const treeData = ref([])
const selectedNodeKey = ref(null)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const errors = ref({})
const flatList = ref([])
const activeTab = ref(0)

const form = ref({
  code: '',
  nomenclature: '',
  parent_id: null,
  sort_order: 0
})

// ── Zones CRUD ──
const zones = ref([])
const zonesTotal = ref(0)
const zonesPage = ref(1)
const zonesLoading = ref(false)
const zoneDialogVisible = ref(false)
const zoneEditing = ref(false)
const zoneEditingId = ref(null)
const zoneSaving = ref(false)
const zoneErrors = ref({})
const zoneForm = ref({
  code: '',
  name: '',
  region: '',
  is_active: true,
  sort_order: 0
})

// Helper functions to convert tree data to TreeNode format for TreeTable
function buildTreeNodes(orgs) {
  return orgs.map(org => ({
    key: org.id,
    data: org,
    children: org.children ? buildTreeNodes(org.children) : []
  }))
}

// Load tree data
async function loadTree() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/organizations?tree=true&per_page=200')
    const data = res.data?.data || res.data || []
    treeData.value = buildTreeNodes(data)

    // Also load flat list for parent selector options
    const flatRes = await api.get('/api/v1/tenant/organizations?per_page=200')
    const flatData = flatRes.data?.data || []
    flatList.value = flatData
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000
    })
  } finally {
    loading.value = false
  }
}

// Parent options for select dropdown
const parentOptions = computed(() => {
  const options = [
    { id: null, label: t('organization.no_parent') }
  ]
  flatList.value.forEach(org => {
    if (!isEditing.value || org.id !== editingId.value) {
      options.push({
        id: org.id,
        label: `${org.full_code} — ${org.nomenclature}`
      })
    }
  })
  return options
})

// Parent label for dialog header
const parentLabel = computed(() => {
  if (!form.value.parent_id) return t('organization.no_parent')
  const parent = flatList.value.find(o => o.id === form.value.parent_id)
  return parent ? `${parent.full_code} — ${parent.nomenclature}` : ''
})

// Open create dialog
function openCreate(parent) {
  isEditing.value = false
  editingId.value = null
  errors.value = {}
  form.value = {
    code: '',
    nomenclature: '',
    parent_id: parent?.id || null,
    sort_order: 0
  }
  dialogVisible.value = true
}

// Open edit dialog
function openEdit(org) {
  isEditing.value = true
  editingId.value = org.id
  errors.value = {}
  form.value = {
    code: org.code,
    nomenclature: org.nomenclature,
    parent_id: org.parent_id || null,
    sort_order: org.sort_order || 0
  }
  dialogVisible.value = true
}

// Reset form
function resetForm() {
  form.value = { code: '', nomenclature: '', parent_id: null, sort_order: 0 }
  errors.value = {}
  isEditing.value = false
  editingId.value = null
}

// Handle save (create or update)
async function handleSave() {
  errors.value = {}

  // Client-side validation
  if (!form.value.code?.trim()) {
    errors.value.code = [t('form.required')]
    return
  }
  if (!form.value.nomenclature?.trim()) {
    errors.value.nomenclature = [t('form.required')]
    return
  }

  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/tenant/organizations/${editingId.value}`, {
        code: form.value.code,
        nomenclature: form.value.nomenclature,
        sort_order: form.value.sort_order || 0
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/organizations', {
        code: form.value.code,
        nomenclature: form.value.nomenclature,
        parent_id: form.value.parent_id,
        sort_order: form.value.sort_order || 0
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.created'), life: 3000 })
    }
    dialogVisible.value = false
    await loadTree()
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

// Confirm delete
function confirmDelete(org) {
  confirm.require({
    header: t('organization.confirm_delete_title'),
    message: t('organization.confirm_delete', { name: org.nomenclature }),
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: t('common.cancel'),
    acceptLabel: t('common.delete'),
    rejectClass: 'p-button-outlined p-button-secondary',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.delete(`/api/v1/tenant/organizations/${org.id}`)
        toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.deleted'), life: 3000 })
        await loadTree()
      } catch (e) {
        toast.add({
          severity: 'error',
          summary: t('message.error'),
          detail: e.response?.data?.error?.message || t('message.operation_failed'),
          life: 4000
        })
      }
    }
  })
}

// ── Zones CRUD Functions ──

async function loadZones() {
  zonesLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/settings/zones?page=${zonesPage.value}&per_page=20`)
    zones.value = res.data?.data?.data || res.data?.data || []
    zonesTotal.value = res.data?.data?.total || res.data?.total || 0
    zonesPage.value = res.data?.data?.page || res.data?.page || 1
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000
    })
  } finally {
    zonesLoading.value = false
  }
}

function openZoneDialog(zone) {
  zoneEditing.value = !!zone
  zoneEditingId.value = zone?.id || null
  zoneErrors.value = {}
  zoneForm.value = {
    code: zone?.code || '',
    name: zone?.name || '',
    region: zone?.region || '',
    is_active: zone?.is_active !== undefined ? zone.is_active : true,
    sort_order: zone?.sort_order || 0
  }
  zoneDialogVisible.value = true
}

function resetZoneForm() {
  zoneForm.value = { code: '', name: '', region: '', is_active: true, sort_order: 0 }
  zoneErrors.value = {}
  zoneEditing.value = false
  zoneEditingId.value = null
}

async function handleZoneSave() {
  zoneErrors.value = {}

  if (!zoneForm.value.code?.trim()) {
    zoneErrors.value = { code: [t('form.required')] }
    return
  }
  if (!zoneForm.value.name?.trim()) {
    zoneErrors.value = { name: [t('form.required')] }
    return
  }

  zoneSaving.value = true
  try {
    if (zoneEditing.value) {
      await api.put(`/api/v1/tenant/settings/zones/${zoneEditingId.value}`, {
        code: zoneForm.value.code,
        name: zoneForm.value.name,
        region: zoneForm.value.region || undefined,
        is_active: zoneForm.value.is_active,
        sort_order: zoneForm.value.sort_order || 0
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/settings/zones', {
        code: zoneForm.value.code,
        name: zoneForm.value.name,
        region: zoneForm.value.region || undefined,
        is_active: zoneForm.value.is_active,
        sort_order: zoneForm.value.sort_order || 0
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.created'), life: 3000 })
    }
    zoneDialogVisible.value = false
    await loadZones()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      zoneErrors.value = fieldErrors
    } else {
      toast.add({
        severity: 'error',
        summary: t('message.error'),
        detail: e.response?.data?.error?.message || t('message.operation_failed'),
        life: 4000
      })
    }
  } finally {
    zoneSaving.value = false
  }
}

function confirmDeleteZone(zone) {
  confirm.require({
    header: t('zones.confirm_delete_title'),
    message: t('zones.confirm_delete', { name: zone.name }),
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: t('common.cancel'),
    acceptLabel: t('common.delete'),
    rejectClass: 'p-button-outlined p-button-secondary',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.delete(`/api/v1/tenant/settings/zones/${zone.id}`)
        toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.deleted'), life: 3000 })
        await loadZones()
      } catch (e) {
        toast.add({
          severity: 'error',
          summary: t('message.error'),
          detail: e.response?.data?.error?.message || t('message.operation_failed'),
          life: 4000
        })
      }
    }
  })
}

function onZonePage(event) {
  zonesPage.value = event.page + 1
  loadZones()
}

onMounted(() => {
  loadTree()
  loadZones()
})
</script>

<style scoped>
:deep(.p-treetable-wrapper) {
  max-height: calc(100vh - 260px);
}
:deep(.p-treetable .p-treetable-tbody > tr) {
  transition: background 0.15s ease;
}
:deep(.p-treetable .p-treetable-tbody > tr:hover) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-treetable .p-treetable-tbody > tr:hover) {
  background: rgba(16, 185, 129, 0.08) !important;
}
</style>
