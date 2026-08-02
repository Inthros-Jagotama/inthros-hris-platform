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
        <div class="flex items-center gap-1 border border-gray-200 dark:border-gray-700 rounded-lg p-0.5 bg-white dark:bg-gray-800">
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
          <Button
            :label="t('organization.tree_view')"
            icon="pi pi-share-alt"
            :severity="viewMode === 'tree' ? 'primary' : 'secondary'"
            :outlined="viewMode !== 'tree'"
            size="small"
            @click="viewMode = 'tree'"
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

      <OrgTreeTable
        v-else
        :nodes="filteredOrgs"
        v-model:expandedRows="expandedRows"
        @add-child="openCreate"
        @edit="openEdit"
        @delete="confirmDelete"
      />
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

    <!-- Organization Tree View (drag & drop untuk pindah parent) -->
    <div v-if="viewMode === 'tree'" class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 overflow-hidden">
      <div class="px-4 py-2.5 border-b border-gray-200 dark:border-gray-700 flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
        <i class="pi pi-share-alt text-xs"></i>
        <span>{{ t('organization.move_hint') }}</span>
      </div>
      <Skeleton v-if="loading" width="100%" height="300px" class="rounded-none" />
      <Tree
        v-else
        v-model:expandedKeys="expandedKeys"
        :value="treeNodes"
        :draggable-nodes="!moving"
        :droppable-nodes="!moving"
        validate-drop
        class="!text-sm !border-0 !shadow-none !p-4"
        @node-drop="onNodeDrop"
      >
        <template #default="{ node }">
          <div class="flex items-center gap-2 py-0.5 select-none">
            <i :class="node.children?.length ? 'pi pi-folder-open text-amber-500 text-xs' : 'pi pi-file text-gray-400 text-xs'"></i>
            <span class="font-medium text-gray-800 dark:text-gray-100">{{ node.data.nomenclature }}</span>
            <Tag :value="node.data.code" severity="info" class="!text-[10px] !px-1 !py-0" />
            <span class="text-xs text-gray-400 dark:text-gray-500 font-mono">{{ node.data.full_code }}</span>
          </div>
        </template>
      </Tree>
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
        <FormRow :label="t('organization.parent')">
          <SelectLabel v-model="form.parent_id" :options="parentOptions" option-value="id" option-label="label" :placeholder="t('organization.select_parent')" :showClear="false" />
        </FormRow>
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
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Skeleton from 'primevue/skeleton'
import Tree from 'primevue/tree'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import OrgChartView from './OrgChartView.vue'
import OrgTreeTable from './OrgTreeTable.vue'

const { t } = useI18n()
const toast = useToast()
const route = useRoute()

const summaryID = computed(() => route.query.summary_id || '')
const viewMode = ref('table')
const loading = ref(false)
const saving = ref(false)
const moving = ref(false)
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

// ── Tree view helpers (drag & drop) ──
function toTreeNodes(orgs) {
  return (orgs || []).map(org => ({
    key: org.id,
    label: org.nomenclature,
    data: org,
    children: org.children?.length ? toTreeNodes(org.children) : undefined
  }))
}

// Tree nodes (dari filteredOrgs agar search juga berlaku di view ini)
const treeNodes = computed(() => toTreeNodes(filteredOrgs.value))

// Cari node org di tree (berdasarkan id)
function findTreeNode(nodes, id) {
  for (const n of nodes || []) {
    if (n.id === id) return n
    const found = findTreeNode(n.children, id)
    if (found) return found
  }
  return null
}

// Kumpulkan semua id descendant dari sebuah node tree
function collectDescendantIds(node) {
  const ids = []
  node.children?.forEach(child => {
    ids.push(child.id)
    ids.push(...collectDescendantIds(child))
  })
  return ids
}

// Apakah targetId merupakan descendant (keturunan) dari orgId?
function isDescendantOf(orgId, targetId) {
  const node = findTreeNode(rootOrgs.value, orgId)
  if (!node) return false
  return collectDescendantIds(node).includes(targetId)
}

// Expanded keys Tree view — root di-expand default saat pertama dibuka
const expandedKeys = ref({})
function seedExpandedKeys() {
  const keys = {}
  rootOrgs.value.forEach(root => { keys[root.id] = true })
  expandedKeys.value = keys
}
watch(viewMode, (mode) => {
  if (mode === 'tree' && Object.keys(expandedKeys.value).length === 0 && rootOrgs.value.length) {
    seedExpandedKeys()
  }
})

// Handler drop pada Tree (PrimeVue 4.5.5): payload = { dragNode, dropNode, dropPosition, accept }.
// - dropPosition 0  → drop ON node → parent baru = dropNode
// - dropPosition ±1 → drop sebelum/sesudah node → sibling → parent baru = parent dropNode
// - dropNode null    → drop di root → parent baru = null (root)
// validate-drop aktif: tree tidak berubah sampai accept() dipanggil, jadi drop ilegal otomatis batal.
async function onNodeDrop(event) {
  if (moving.value) return
  if (!event.accept) return // event lanjutan setelah accept() dipanggil — abaikan
  const dragged = event.dragNode
  const dropNode = event.dropNode
  const pos = event.dropPosition ?? 0

  let newParentId = null
  if (dropNode) {
    newParentId = pos === 0 ? dropNode.data.id : (dropNode.data.parent_id || null)
  }

  // Guard: dilarang meletakkan node ke dirinya sendiri / ke bawah keturunannya.
  // accept() tetap dipanggil untuk me-reset state drag internal PrimeVue,
  // lalu loadTree() segera mengembalikan tampilan ke data server.
  if (newParentId === dragged.data.id || (newParentId && isDescendantOf(dragged.data.id, newParentId))) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('organization.cannot_move_into_descendant'), life: 4000 })
    event.accept()
    await loadTree()
    return
  }
  // Parent tidak berubah → settle drop + reload (urutan sibling tak dikelola via drag)
  if (newParentId === (dragged.data.parent_id || null)) {
    event.accept()
    await loadTree()
    return
  }
  // Terima perubahan visual, lalu simpan ke server
  event.accept()
  await moveOrg(dragged.data.id, newParentId)
}

// Pindahkan organisasi ke parent baru (null/'' = root)
async function moveOrg(id, parentId) {
  moving.value = true
  try {
    await api.put(`/api/v1/tenant/organizations/${id}`, { parent_id: parentId || '' })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('organization.moved'), life: 3000 })
    await loadTree()
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000
    })
    await loadTree()
  } finally {
    moving.value = false
  }
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

// Parent options for select dropdown (exclude self + descendants saat edit)
const parentOptions = computed(() => {
  const options = [
    { id: null, label: t('organization.no_parent') }
  ]
  const excluded = new Set()
  if (isEditing.value && editingId.value) {
    excluded.add(editingId.value)
    const node = findTreeNode(rootOrgs.value, editingId.value)
    if (node) collectDescendantIds(node).forEach(id => excluded.add(id))
  }
  flatAllOrgs.value.forEach(org => {
    if (!excluded.has(org.id)) {
      options.push({
        id: org.id,
        label: `${org.full_code} — ${org.nomenclature}`
      })
    }
  })
  return options
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
        parent_id: form.value.parent_id || '',
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
        parent_id: form.value.parent_id || '',
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
