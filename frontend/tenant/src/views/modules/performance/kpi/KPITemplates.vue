<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <div class="relative w-64 max-w-full">
          <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-gray-400 dark:text-gray-500"></i>
          <InputText
            v-model="searchQuery"
            :placeholder="t('kpi.search_templates')"
            class="w-full !pl-8 !text-sm"
            @input="onSearchInput"
          />
          <button
            v-if="searchQuery"
            class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 text-xs"
            @click="clearSearch"
          >
            <i class="pi pi-times"></i>
          </button>
        </div>
        <Button :label="t('kpi.new_template')" icon="pi pi-plus" size="small" @click="createTemplate" />
      </div>
    </div>

    <!-- DataTable -->
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="items"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-file text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('kpi.empty_templates') }}</p>
          <p class="text-xs mt-1">{{ t('kpi.empty_templates_hint') }}</p>
        </div>
      </template>

      <Column field="name" :header="t('kpi.template_name')" sortable>
        <template #body="{data}">
          <div>
            <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
            <p v-if="data.description" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-1">{{ data.description }}</p>
          </div>
        </template>
      </Column>

      <Column field="organization_name" :header="t('kpi.organization')" sortable style="width:200px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.organization_name || '-' }}</span>
        </template>
      </Column>

      <Column field="period_code" :header="t('kpi.period')" sortable style="width:120px">
        <template #body="{data}">
          <Tag v-if="data.period_code" :value="data.period_code" severity="info" class="!text-xs" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>

      <Column :header="t('kpi.indicators')" style="width:100px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300 font-mono">{{ data.indicator_count || 0 }}</span>
        </template>
      </Column>

      <Column field="status" :header="t('kpi.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="getStatusLabel(data.status)" :severity="getStatusSeverity(data.status)" class="!text-xs" />
        </template>
      </Column>

      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <div v-if="canManageTemplate(data)" class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="editTemplate(data)" />
            <Button icon="pi pi-copy" size="small" text severity="info" v-tooltip.left="t('kpi.duplicate')" @click="duplicateTemplate(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
          <span v-else v-tooltip.left="t('kpi.cannot_manage_template')" class="inline-flex items-center justify-center w-7 h-7 text-gray-300 dark:text-gray-600">
            <i class="pi pi-lock text-sm"></i>
          </span>
        </template>
      </Column>
    </DataTable>

    <!-- Delete Confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('kpi.delete_template_title')"
      :message="t('kpi.delete_template_message', { name: deleteTarget?.name })"
      :loading="deleting"
      :error="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import InputText from 'primevue/inputtext'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const router = useRouter()
const toast = useToast()
const { t } = useI18n()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const searchQuery = ref('')
let searchTimer = null

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

// Organisasi tempat user saat ini bekerja (dari my-context — posisi jabatan
// terakhir). Template hanya bisa dikelola (edit/duplicate/hapus) oleh anggota
// organisasinya — aturan yang sama dengan validasi backend.
const myOrgId = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function getStatusLabel(status) {
  switch (status) {
    case 'PUBLISHED': return 'Published'
    case 'ARCHIVED': return 'Archived'
    default: return 'Draft'
  }
}

function getStatusSeverity(status) {
  switch (status) {
    case 'PUBLISHED': return 'success'
    case 'ARCHIVED': return 'secondary'
    default: return 'warn'
  }
}

async function loadMyOrg() {
  // Fail-open: jika gagal resolve, myOrgId tetap null sehingga semua tombol
  // tetap tampil (backend tetap menolak aksi yang tidak berhak).
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/my-context')
    const ctx = res.data?.data
    if (ctx?.has_position && ctx.organization_id) {
      myOrgId.value = ctx.organization_id
    }
  } catch {
    myOrgId.value = null
  }
}

function canManageTemplate(template) {
  if (!myOrgId.value) return true
  // Otorisasi STRICT berbasis organisasi PEMBUAT (created_by_org_id, disimpan
  // saat template pertama dibuat). Template tanpa created_by_org_id (legacy)
  // tidak bisa dikelola sama sekali.
  return !!template.created_by_org_id && template.created_by_org_id === myOrgId.value
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/performance/kpi/templates', {
      params: {
        page: currentPage.value,
        per_page: perPage.value,
        search: searchQuery.value || undefined
      }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadData()
  }, 350)
}

function clearSearch() {
  searchQuery.value = ''
  clearTimeout(searchTimer)
  currentPage.value = 1
  loadData()
}

function createTemplate() {
  router.push('/performance/kpi/templates/new')
}

function editTemplate(item) {
  router.push(`/performance/kpi/templates/${item.id}/edit`)
}

async function duplicateTemplate(item) {
  try {
    const res = await api.post(`/api/v1/tenant/performance/kpi/templates/${item.id}/duplicate`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.template_duplicated'), life: 3000 })
    const newId = res.data?.data?.id || res.data?.id
    if (newId) {
      router.push(`/performance/kpi/templates/${newId}/edit`)
    } else {
      await loadData()
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/performance/kpi/templates/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('kpi.template_deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await loadMyOrg()
  await loadData()
})
</script>
