<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
        <Select
          v-model="filterDocumentType"
          :options="documentTypeOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('document_templates.filter_document_type')"
          class="!text-sm w-56"
          showClear
          @change="onFilterChange"
        />
        <Select
          v-model="filterStatus"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('document_templates.filter_status')"
          class="!text-sm w-40"
          showClear
          @change="onFilterChange"
        />
        <span class="p-input-icon-left">
          <TextInput
            v-model="searchTerm"
            :placeholder="t('common.search')"
            class="!text-sm w-56"
            @update:modelValue="onSearchInput"
          />
        </span>
      </div>
    </div>

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
          <i class="pi pi-file-edit text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('document_templates.empty_title') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('document_templates.name')">
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
            <i v-if="data.is_default" class="pi pi-star-fill text-amber-400 text-xs" v-tooltip.top="t('document_templates.default')"></i>
          </div>
        </template>
      </Column>
      <Column field="code" :header="t('document_templates.code')">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ data.code }}</span>
        </template>
      </Column>
      <Column field="document_type" :header="t('document_templates.document_type')">
        <template #body="{ data }">
          <span class="text-gray-700 dark:text-gray-200">{{ documentTypeLabel(data.document_type) }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:120px">
        <template #body="{ data }">
          <Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="updated_at" :header="t('document_templates.updated_at')" style="width:160px">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.updated_at, locale) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:140px" frozen alignFrozen="right">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button
              v-if="!data.is_default && data.status !== 'ACTIVE'"
              icon="pi pi-check-circle"
              size="small"
              text
              severity="success"
              v-tooltip.left="t('document_templates.activate')"
              :loading="actioningId === data.id"
              @click="handleActivate(data)"
            />
            <Button
              v-if="!data.is_default && data.status === 'ACTIVE'"
              icon="pi pi-times-circle"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('document_templates.deactivate')"
              :loading="actioningId === data.id"
              @click="handleDeactivate(data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'
import TextInput from '@/components/TextInput.vue'
import { formatDate } from '@/utils/formatDate'

const { t, locale } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const filterDocumentType = ref(null)
const filterStatus = ref(null)
const searchTerm = ref('')
let searchDebounce = null

const actioningId = ref(null)

const documentTypeOptions = [
  { label: t('document_templates.type_contract_agreement'), value: 'CONTRACT_AGREEMENT' },
  { label: t('document_templates.type_movement_sk'), value: 'MOVEMENT_SK' },
]

const statusOptions = [
  { label: t('document_templates.status_active'), value: 'ACTIVE' },
  { label: t('document_templates.status_inactive'), value: 'INACTIVE' },
  { label: t('document_templates.status_reference'), value: 'REFERENCE' },
]

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' },
]

function documentTypeLabel(type) {
  const found = documentTypeOptions.find((o) => o.value === type)
  return found ? found.label : type
}

function statusLabel(status) {
  const found = statusOptions.find((o) => o.value === status)
  return found ? found.label : status
}

function statusSeverity(status) {
  if (status === 'ACTIVE') return 'success'
  if (status === 'REFERENCE') return 'info'
  return 'warn'
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/document-templates', {
      params: {
        page: currentPage.value,
        per_page: perPage.value,
        document_type: filterDocumentType.value || undefined,
        status: filterStatus.value || undefined,
        search: searchTerm.value || undefined,
      },
    })
    const body = res.data?.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

function onSearchInput() {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    currentPage.value = 1
    loadData()
  }, 400)
}

async function handleActivate(item) {
  actioningId.value = item.id
  try {
    await api.post(`/api/v1/tenant/document-templates/${item.id}/activate`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.activated'), life: 3000 })
    await loadData()
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000,
    })
  } finally {
    actioningId.value = null
  }
}

async function handleDeactivate(item) {
  actioningId.value = item.id
  try {
    await api.post(`/api/v1/tenant/document-templates/${item.id}/deactivate`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.deactivated'), life: 3000 })
    await loadData()
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000,
    })
  } finally {
    actioningId.value = null
  }
}

onMounted(() => {
  loadData()
})
</script>
