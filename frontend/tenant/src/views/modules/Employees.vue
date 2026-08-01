<template>
  <div class="space-y-1">
    <!-- Search + Add -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('employee.search')" size="small" class="!pl-8 !text-sm !py-1.5 !w-64" @input="onSearchInput" />
        </IconField>
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button v-if="hasPermission('employee.create')" :label="t('common.add')" icon="pi pi-plus" size="small" @click="goToNew()" />
      </div>
    </div>

    <!-- Table -->
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
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('employee.empty_title') }}</p>
        </div>
      </template>

      <Column field="employee_id" :header="t('employee.employee_id')" sortable style="width:130px">
        <template #body="{data}">
          <Tag :value="data.employee_id" severity="info" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="name" :header="t('common.name')" sortable style="min-width:180px">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="nik" :header="t('employee.nik')" sortable style="width:140px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.nik || '-' }}</span>
        </template>
      </Column>
      <Column field="gender" :header="t('employee.gender')" sortable style="width:80px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.gender ? t('employee.gender_' + data.gender.toLowerCase()) : '-' }}</span>
        </template>
      </Column>
      <Column field="phone_number" :header="t('employee.phone')" sortable style="width:140px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.phone_number || '-' }}</span>
        </template>
      </Column>
      <Column field="email" :header="t('employee.email')" sortable style="min-width:200px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-400">{{ data.email || '-' }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="t('common_status.' + data.status)" :severity="data.status === 'active' ? 'success' : data.status === 'inactive' ? 'warn' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button v-if="hasPermission('employee.update')" icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="goToEdit(data)" />
            <Button v-if="hasPermission('employee.delete')" icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Delete Confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="deleting"
      :error="deleteError"
      :title="t('employee.confirm_delete_title')"
      :message="t('employee.confirm_delete', { name: deleteTarget?.name || '' })"
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Tag from 'primevue/tag'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

// ── Data ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const searchQuery = ref('')
let searchTimer = null

// ── Navigation ──
function goToNew() { router.push('/employees/new') }
function goToEdit(item) { router.push(`/employees/${item.id}/edit`) }

// ── Delete ──
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

// ── Computed ──
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const skeletonColumns = [
  { type: 'tag', width: 'w-28', headerWidth: 'w-20' },
  { type: 'text', width: 'w-44', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-12' },
  { type: 'text', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-28', headerWidth: 'w-16' },
  { type: 'text', width: 'w-44', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

// ── Methods ──
async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const q = searchQuery.value?.trim()
    if (q) params.search = q
    const res = await api.get('/api/v1/tenant/employees', { params })
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
  }, 400)
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
    await api.delete(`/api/v1/tenant/employees/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
