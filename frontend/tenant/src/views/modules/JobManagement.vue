<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
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
          <i class="pi pi-sitemap text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('job_management.empty_title') }}</p>
          <p class="text-xs mt-1">{{ t('job_management.empty_hint') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('organization.code')" sortable style="width:100px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !font-mono !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="full_code" :header="t('organization.full_code')" sortable style="width:120px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ data.full_code }}</span></template>
      </Column>
      <Column field="nomenclature" :header="t('organization.nomenclature')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.nomenclature }}</span></template>
      </Column>
      <Column field="level" :header="t('organization.level')" sortable style="width:90px">
        <template #body="{data}">
          <Tag :value="'L' + data.level" severity="contrast" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="sort_order" :header="t('organization.sort_order')" sortable style="width:100px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:60px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button
            icon="pi pi-sliders-h"
            size="small"
            text
            severity="info"
            v-tooltip.left="t('job_management.values')"
            @click="openValuesPage(data)"
          />
        </template>
      </Column>
    </DataTable>
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
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const skeletonColumns = [
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/organizations', {
      params: { page: currentPage.value, per_page: perPage.value, active_only: true }
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

function openValuesPage(org) {
  if (!org.id) return
  // Navigate to the multi-section Job Management form page for this organization
  router.push(`/job-management/form?org_id=${org.id}`)
}

onMounted(loadData)
</script>
