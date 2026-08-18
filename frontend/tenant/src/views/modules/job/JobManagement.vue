<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="relative w-64 max-w-full">
        <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-gray-400 dark:text-gray-500"></i>
        <InputText
          v-model="searchQuery"
          :placeholder="t('job_management.search_positions')"
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
      <Column field="full_code" :header="t('organization.full_code')" sortable style="width:120px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ data.full_code }}</span></template>
      </Column>
      <Column field="nomenclature" :header="t('organization.nomenclature')" sortable>
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.nomenclature }}</span></template>
      </Column>
      <Column :header="t('job_management.score')" style="width:110px" field="score">
        <template #body="{data}">
          <span v-if="data.score != null" class="text-emerald-600 dark:text-emerald-400 font-bold">{{ formatNumber(data.score) }}</span>
          <span v-else class="text-gray-300 dark:text-gray-600">—</span>
        </template>
      </Column>
      <Column :header="t('job_management.score_with_financial')" style="width:120px" field="score_has_financial">
        <template #body="{data}">
          <Tag
            v-if="data.score_has_financial != null"
            :value="data.score_has_financial ? t('common.yes') : t('common.no')"
            :severity="data.score_has_financial ? 'success' : 'danger'"
            class="!text-xs"
          />
          <span v-else class="text-gray-300 dark:text-gray-600">—</span>
        </template>
      </Column>
      <Column :header="t('job_management.score_complete_status')" style="width:130px" field="score_complete">
        <template #body="{data}">
          <Tag
            v-if="data.score_complete != null"
            :value="data.score_complete ? t('job_management.score_complete') : t('job_management.score_incomplete')"
            :severity="data.score_complete ? 'success' : 'warning'"
            :icon="data.score_complete ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle'"
            class="!text-xs"
          />
          <span v-else class="text-gray-300 dark:text-gray-600">—</span>
        </template>
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
import InputText from 'primevue/inputtext'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const searchQuery = ref('')
let searchTimer = null

const skeletonColumns = [
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/organizations', {
      params: {
        page: currentPage.value,
        per_page: perPage.value,
        active_only: true,
        search: searchQuery.value || undefined
      }
    })
    const body = res.data
    const orgs = body?.data || []
    items.value = orgs
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page

    // Ambil skor jabatan untuk semua organisasi (map by organization_id).
    // Service membatasi per_page maks 100 → loop pagination sampai semua terambil.
    const scores = []
    let scorePage = 1
    let scoreTotal = 0
    let pageData = []
    do {
      const scoresRes = await api.get('/api/v1/tenant/job-management/scores', { params: { page: scorePage, per_page: 100 } }).catch(() => null)
      pageData = scoresRes?.data?.data || []
      scores.push(...pageData)
      scoreTotal = scoresRes?.data?.total || scores.length
      scorePage++
    } while (scores.length < scoreTotal && pageData.length > 0)
    const scoreMap = {}
    scores.forEach(s => {
      if (s.organization_id) scoreMap[s.organization_id] = s
    })
    items.value = orgs.map(o => {
      const sc = scoreMap[o.id]
      return {
        ...o,
        score: sc ? sc.job_value_without_financial : null,
        score_has_financial: sc ? (sc.has_financial_authority ? 1 : 0) : null,
        score_complete: sc ? sc.is_complete : null
      }
    })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

function formatNumber(n) {
  if (n == null) return ''
  return Number(n).toLocaleString?.('id-ID') ?? String(n)
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

function openValuesPage(org) {
  if (!org.id) return
  // Navigate to the multi-section Job Management form page for this organization
  router.push(`/job-management/form?org_id=${org.id}`)
}

onMounted(loadData)
</script>
