<template>
  <div class="space-y-1">
    <!-- Search & Actions — satu baris -->
    <div class="flex items-center justify-between gap-3 dark:border-gray-700">
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" />
        </IconField>
        <Button
          v-if="searchQuery"
          icon="pi pi-times"
          severity="secondary"
          text
          rounded
          size="small"
          class="!p-1"
          @click="searchQuery = ''"
        />
      </div>
      <div class="flex items-center gap-2">
        <Button
          :label="t('organization.add_root')"
          icon="pi pi-plus"
          size="small"
          severity="primary"
          @click="openCreate(null)"
        />
        <Button
          icon="pi pi-refresh"
          v-tooltip.top="t('common.refresh')"
          severity="secondary"
          text
          size="small"
          @click="loadTree()"
          :loading="loading"
        />
        <div class="flex items-center gap-1 ml-2 pl-2 border-l border-gray-200 dark:border-gray-700">
          <Button
            :label="t('organization.table_view')"
            icon="pi pi-table"
            :severity="viewMode === 'table' ? 'primary' : 'secondary'"
            :outlined="viewMode !== 'table'"
            size="small"
            @click="viewMode = 'table'"
          />
          <Button
            :label="t('organization.chart_view')"
            icon="pi pi-sitemap"
            :severity="viewMode === 'chart' ? 'primary' : 'secondary'"
            :outlined="viewMode !== 'chart'"
            size="small"
            @click="viewMode = 'chart'"
          />
        </div>
      </div>
    </div>

    <!-- Organization Drilldown Table -->
    <div v-show="viewMode === 'table'" class="">
      <!-- Skeleton -->
      <div v-if="loading" class="space-y-2">
        <div v-for="i in 5" :key="i" class="flex items-center gap-3 py-1">
          <Skeleton width="4rem" height="1rem" />
          <Skeleton width="12rem" height="1rem" />
          <Skeleton width="8rem" height="1rem" />
        </div>
      </div>

      <!-- DataTable with drilldown -->
      <DataTable v-if="!loading"
        :value="searchQuery ? filteredOrgs : rootOrgs"
        v-model:expandedRows="expandedRows"
        dataKey="id"
        class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        stripedRows
        responsiveLayout="scroll"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-sitemap text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('organization.empty_title') }}</p>
            <p class="text-sm mt-1">{{ t('organization.empty_tree') }}</p>
          </div>
        </template>
        <Column :expander="true" style="width: 40px" />
        <Column field="nomenclature" :header="t('organization.nomenclature')">
          <template #body="{ data }">
            <div class="flex items-center gap-2">
              <i class="pi pi-folder-open text-amber-500 text-xs"></i>
              <span class="font-medium text-gray-800 dark:text-gray-100">{{ data.nomenclature }}</span>
              <Tag
                v-if="data.children?.length"
                :value="data.children.length"
                severity="info"
                class="!text-[10px] !px-1 !py-0 !min-w-[1.1rem]"
                rounded
              />
            </div>
          </template>
        </Column>
        <Column field="code" :header="t('organization.code')" style="width: 120px">
          <template #body="{ data }">
            <Tag :value="data.code" severity="info" class="!text-xs" />
          </template>
        </Column>
        <Column field="full_code" :header="t('organization.full_code')" style="width: 160px">
          <template #body="{ data }">
            <span class="text-gray-500 dark:text-gray-400 text-xs font-mono">{{ data.full_code }}</span>
          </template>
        </Column>
        <Column field="level" :header="t('organization.level')" style="width: 80px">
          <template #body="{ data }">
            <span class="text-gray-500 dark:text-gray-400">{{ data.level }}</span>
          </template>
        </Column>
        <Column field="sort_order" :header="t('organization.sort_order')" style="width: 80px">
          <template #body="{ data }">
            <span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width: 130px" frozen alignFrozen="right">
          <template #body="{ data }">
            <div class="flex items-center gap-1">
              <Button icon="pi pi-plus" v-tooltip.top="t('organization.add_child')" severity="secondary" text size="small" class="!p-1" @click="openCreate(data)" />
              <Button icon="pi pi-pencil" v-tooltip.top="t('common.edit')" severity="secondary" text size="small" class="!p-1" @click="openEdit(data)" />
              <Button icon="pi pi-trash" v-tooltip.top="t('common.delete')" severity="danger" text size="small" class="!p-1" @click="confirmDelete(data)" />
            </div>
          </template>
        </Column>

        <!-- Root expansion slot (PrimeVue 4 — at DataTable level) -->
        <template #expansion="{ data }">
          <div v-if="data.children?.length" class="pl-6 pr-2 py-1">
            <DataTable
              :value="data.children"
              v-model:expandedRows="expandedRows"
              dataKey="id"
              class="!text-xs !border-0 !shadow-none"
              stripedRows
              responsiveLayout="scroll"
            >
              <Column :expander="true" style="width: 40px" />
              <Column field="nomenclature">
                <template #body="{ data: child }">
                  <div class="flex items-center gap-2">
                    <i class="pi pi-folder-open text-amber-500 text-xs"></i>
                    <span class="font-medium text-gray-800 dark:text-gray-100">{{ child.nomenclature }}</span>
                    <Tag
                      v-if="child.children?.length"
                      :value="child.children.length"
                      severity="info"
                      class="!text-[10px] !px-1 !py-0 !min-w-[1.1rem]"
                      rounded
                    />
                  </div>
                </template>
              </Column>
              <Column field="code" :header="t('organization.code')" style="width: 120px">
                <template #body="{ data: child }">
                  <Tag :value="child.code" severity="info" class="!text-xs" />
                </template>
              </Column>
              <Column field="full_code" :header="t('organization.full_code')" style="width: 160px">
                <template #body="{ data: child }">
                  <span class="text-gray-500 dark:text-gray-400 text-xs font-mono">{{ child.full_code }}</span>
                </template>
              </Column>
              <Column field="level" :header="t('organization.level')" style="width: 80px">
                <template #body="{ data: child }">
                  <span class="text-gray-500 dark:text-gray-400">{{ child.level }}</span>
                </template>
              </Column>
              <Column field="sort_order" :header="t('organization.sort_order')" style="width: 80px">
                <template #body="{ data: child }">
                  <span class="text-gray-500 dark:text-gray-400">{{ child.sort_order }}</span>
                </template>
              </Column>
              <Column :header="t('common.actions')" style="width: 120px" frozen alignFrozen="right">
                <template #body="{ data: child }">
                  <div class="flex items-center gap-1">
                    <Button icon="pi pi-plus" v-tooltip.top="t('organization.add_child')" severity="secondary" text size="small" class="!p-1" @click="openCreate(child)" />
                    <Button icon="pi pi-pencil" v-tooltip.top="t('common.edit')" severity="secondary" text size="small" class="!p-1" @click="openEdit(child)" />
                    <Button icon="pi pi-trash" v-tooltip.top="t('common.delete')" severity="danger" text size="small" class="!p-1" @click="confirmDelete(child)" />
                  </div>
                </template>
              </Column>
              <!-- Child expansion slot (PrimeVue 4 — at DataTable level) -->
              <template #expansion="{ data: child }">
                <div v-if="child.children?.length" class="pl-6 pr-2 py-1">
                  <DataTable
                    :value="child.children"
                    v-model:expandedRows="expandedRows"
                    dataKey="id"
                    class="!text-xs !border-0 !shadow-none"
                    stripedRows
                    responsiveLayout="scroll"
                  >
                    <Column :expander="true" style="width: 40px" />
                    <Column field="nomenclature">
                      <template #body="{ data: grandchild }">
                        <div class="flex items-center gap-2">
                          <i class="pi pi-folder-open text-amber-500 text-xs"></i>
                          <span class="font-medium text-gray-800 dark:text-gray-100">{{ grandchild.nomenclature }}</span>
                        </div>
                      </template>
                    </Column>
                    <Column field="code" :header="t('organization.code')" style="width: 120px">
                      <template #body="{ data: grandchild }">
                        <Tag :value="grandchild.code" severity="info" class="!text-xs" />
                      </template>
                    </Column>
                    <Column field="full_code" :header="t('organization.full_code')" style="width: 160px">
                      <template #body="{ data: grandchild }">
                        <span class="text-gray-500 dark:text-gray-400 text-xs font-mono">{{ grandchild.full_code }}</span>
                      </template>
                    </Column>
                    <Column field="sort_order" :header="t('organization.sort_order')" style="width: 80px">
                      <span class="text-gray-500 dark:text-gray-400">{{ grandchild.sort_order }}</span>
                    </Column>
                    <Column :header="t('common.actions')" style="width: 120px" frozen alignFrozen="right">
                      <template #body="{ data: grandchild }">
                        <div class="flex items-center gap-1">
                          <Button icon="pi pi-plus" v-tooltip.top="t('organization.add_child')" severity="secondary" text size="small" class="!p-1" @click="openCreate(grandchild)" />
                          <Button icon="pi pi-pencil" v-tooltip.top="t('common.edit')" severity="secondary" text size="small" class="!p-1" @click="openEdit(grandchild)" />
                          <Button icon="pi pi-trash" v-tooltip.top="t('common.delete')" severity="danger" text size="small" class="!p-1" @click="confirmDelete(grandchild)" />
                        </div>
                      </template>
                    </Column>
                  </DataTable>
                </div>
                <div v-else class="pl-6 pr-2 py-2 text-xs text-gray-400 italic">
                  {{ t('organization.no_children') }}
                </div>
              </template>
            </DataTable>
          </div>
          <div v-else class="pl-6 pr-2 py-2 text-xs text-gray-400 italic">
            {{ t('organization.no_children') }}
          </div>
        </template>
      </DataTable>
    </div>

    <!-- Organization Chart View (v-if agar chart init dgn dimensi benar) -->
    <div v-if="viewMode === 'chart'" class="">
      <Skeleton v-if="loading" width="100%" height="400px" class="rounded-lg" />
      <OrgChartView
        v-else
        :data="rootOrgs"
        :loading="loading"
        @node-click="handleChartNodeClick"
      />
    </div>

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
        <FormRow :label="t('organization.code')" required :errors="errors?.code">
            <TextInput v-model="form.code" :class="{ 'p-invalid': errors?.code }" maxlength="10" :placeholder="t('organization.code')" />
          </FormRow>
          <FormRow :label="t('organization.nomenclature')" required :errors="errors?.nomenclature">
            <TextInput v-model="form.nomenclature" :class="{ 'p-invalid': errors?.nomenclature }" maxlength="255" :placeholder="t('organization.nomenclature')" />
          </FormRow>
          <FormRow :label="t('organization.zone')">
            <SelectLabel v-model="form.zone_id" :options="zoneOptions" option-value="id" option-label="label" :placeholder="t('organization.select_zone')" :showClear="true" />
          </FormRow>
          <FormRow :label="t('organization.job_family')">
            <SelectLabel v-model="form.job_family_id" :options="jobFamilyOptions" option-value="id" option-label="label" :placeholder="t('organization.select_job_family')" :showClear="true" />
          </FormRow>
          <FormRow :label="t('organization.sort_order')">
            <InputNumber v-model="form.sort_order" class="!w-full" :min="0" size="small" />
          </FormRow>
          <FormRow :label="t('organization.grading')">
            <SelectLabel v-model="form.grading_id" :options="gradingOptions" option-value="id" option-label="label" :placeholder="t('organization.select_grading')" :showClear="true" />
          </FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
      </template>
    </Dialog>

    <!-- Delete Confirm -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('organization.confirm_delete_title')"
      :message="t('organization.confirm_delete', { name: deleteTarget?.nomenclature || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      :confirm-label="t('common.delete')"
      :cancel-label="t('common.cancel')"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Skeleton from 'primevue/skeleton'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import OrgChartView from './OrgChartView.vue'

const { t } = useI18n()
const toast = useToast()
const route = useRoute()

const summaryID = computed(() => route.query.summary_id || '')
const viewMode = ref('table')
const loading = ref(false)
const saving = ref(false)
const rootOrgs = ref([])
const searchQuery = ref('')
const expandedRows = ref({})

// Reference data for dropdowns
const zones = ref([])
const jobFamilies = ref([])
const gradings = ref([])
const refLoading = ref(false)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const form = ref({
  code: '',
  nomenclature: '',
  parent_id: null,
  zone_id: null,
  job_family_id: null,
  grading_id: null,
  sort_order: 0,
  organization_summary_id: ''
})

// Filter tree nodes recursively by search query
function filterTree(nodes, query) {
  if (!query) return nodes
  const q = query.toLowerCase()
  const result = []
  nodes.forEach(node => {
    const match = node.nomenclature?.toLowerCase().includes(q)
      || node.code?.toLowerCase().includes(q)
      || node.full_code?.toLowerCase().includes(q)
    const filteredChildren = node.children?.length ? filterTree(node.children, q) : []
    if (match || filteredChildren.length > 0) {
      result.push({
        ...node,
        children: filteredChildren
      })
    }
  })
  return result
}

// Filtered org tree (by search query)
const filteredOrgs = computed(() => filterTree(rootOrgs.value, searchQuery.value))

// Flatten tree recursively for parent options dropdown
function flattenTree(nodes) {
  const result = []
  nodes.forEach(node => {
    result.push(node)
    if (node.children?.length) {
      result.push(...flattenTree(node.children))
    }
  })
  return result
}

// Computed options for reference dropdowns
const zoneOptions = computed(() => zones.value.map(z => ({ id: z.id, label: `${z.code} — ${z.name}` })))
const jobFamilyOptions = computed(() => jobFamilies.value.map(jf => ({ id: jf.id, label: `${jf.code} — ${jf.name}` })))
const gradingOptions = computed(() => gradings.value.map(g => ({ id: g.id, label: `${g.code} — ${g.name}` })))


// Load reference data (zones, job families, gradings) from setting module
async function loadRefData() {
  refLoading.value = true
  const promises = [
    api.get('/api/v1/tenant/settings/zones?per_page=200').then(r => { zones.value = r.data?.data || [] }).catch(() => { zones.value = [] }),
    api.get('/api/v1/tenant/settings/job-families?per_page=200').then(r => { jobFamilies.value = r.data?.data || [] }).catch(() => { jobFamilies.value = [] }),
    api.get('/api/v1/tenant/settings/gradings?per_page=200').then(r => { gradings.value = r.data?.data || [] }).catch(() => { gradings.value = [] })
  ]
  await Promise.all(promises)
  refLoading.value = false
}

// Load organization tree data
async function loadTree() {
  loading.value = true
  try {
    const params = { tree: 'true' }
    if (summaryID.value) params.summary_id = summaryID.value
    const res = await api.get('/api/v1/tenant/organizations', { params })
    const data = res.data?.data || []
    rootOrgs.value = data
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

// Flat list of all orgs from tree (for dropdowns)
const flatAllOrgs = computed(() => flattenTree(rootOrgs.value))

// Parent options for select dropdown
const parentOptions = computed(() => {
  const options = [
    { id: null, label: t('organization.no_parent') }
  ]
  flatAllOrgs.value.forEach(org => {
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
  const parent = flatAllOrgs.value.find(o => o.id === form.value.parent_id)
  return parent ? `${parent.full_code} — ${parent.nomenclature}` : ''
})

// Open create dialog
function openCreate(parent) {
  // Wajib ada summary_id dari URL
  if (!summaryID.value) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('organization.require_summary'), life: 4000 })
    return
  }

  isEditing.value = false
  editingId.value = null
  errors.value = {}
  form.value = {
    code: '',
    nomenclature: '',
    parent_id: parent?.id || null,
    sort_order: 0,
    organization_summary_id: summaryID.value
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
    zone_id: org.zone_id || null,
    job_family_id: org.job_family_id || null,
    grading_id: org.grading_id || null,
    sort_order: org.sort_order || 0,
    organization_summary_id: summaryID.value
  }
  dialogVisible.value = true
}

// Reset form
function resetForm() {
  form.value = { code: '', nomenclature: '', parent_id: null, zone_id: null, job_family_id: null, grading_id: null, sort_order: 0, organization_summary_id: '' }
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
  if (!form.value.organization_summary_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('organization.require_summary'), life: 4000 })
    return
  }

  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/tenant/organizations/${editingId.value}`, {
        code: form.value.code,
        nomenclature: form.value.nomenclature,
        zone_id: form.value.zone_id || null,
        job_family_id: form.value.job_family_id || null,
        grading_id: form.value.grading_id || null,
        sort_order: form.value.sort_order || 0
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/organizations', {
        code: form.value.code,
        nomenclature: form.value.nomenclature,
        parent_id: form.value.parent_id,
        zone_id: form.value.zone_id || null,
        job_family_id: form.value.job_family_id || null,
        grading_id: form.value.grading_id || null,
        sort_order: form.value.sort_order || 0,
        organization_summary_id: form.value.organization_summary_id
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
  deleteTarget.value = org
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/organizations/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadTree()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

// Handle node click from org chart
function handleChartNodeClick(nodeData) {
  // Find the actual org object from rootOrgs tree
  const flat = flattenTree(rootOrgs.value)
  const org = flat.find(o => o.id === nodeData.id)
  if (org) {
    openEdit(org)
  }
}

onMounted(async () => {
  await Promise.all([loadTree(), loadRefData()])
})
</script>

<style scoped>
:deep(.p-datatable-wrapper) {
  max-height: calc(100vh - 260px);
}
:deep(.p-datatable .p-datatable-tbody > tr) {
  transition: background 0.15s ease;
}
:deep(.p-datatable .p-datatable-tbody > tr:hover) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-datatable .p-datatable-tbody > tr:hover) {
  background: rgba(16, 185, 129, 0.08) !important;
}
:deep(.p-datatable .p-datatable-tbody > tr.p-row-expanded) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-datatable .p-datatable-tbody > tr.p-row-expanded) {
  background: rgba(16, 185, 129, 0.08) !important;
}
</style>
