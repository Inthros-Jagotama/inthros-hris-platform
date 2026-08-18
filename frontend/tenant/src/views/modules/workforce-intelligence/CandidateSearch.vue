<template>
  <div class="space-y-4">
    <!-- Toolbar: search + ringkasan -->
    <div class="flex flex-wrap items-center gap-2">
      <div class="relative flex-1 min-w-[220px] max-w-md">
        <i class="pi pi-search absolute left-3 top-1/2 -translate-y-1/2 text-xs text-gray-400"></i>
        <InputText
          v-model="searchTerm"
          class="!w-full !pl-8 !py-2 !text-sm"
          :placeholder="t('candidate_search.search_placeholder')"
          @keyup.enter="onSearch()"
        />
      </div>
      <Button
        :label="t('candidate_search.search')"
        icon="pi pi-search"
        size="small"
        :loading="loading"
        @click="onSearch()"
      />
      <Tag v-if="totalRecords > 0" :value="`${totalRecords} ${t('candidate_search.' + (totalRecords === 1 ? 'vacancy_one' : 'vacancy_count'))}`" severity="info" class="!text-xs !px-2 !py-1" />
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
      v-model:expandedRows="expandedRows"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-user-plus text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('candidate_search.empty_positions') }}</p>
        </div>
      </template>

      <Column expander style="width: 3rem" />

      <Column :header="t('candidate_search.position')" sortable>
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400 flex items-center justify-center shrink-0">
              <i class="pi pi-briefcase text-xs"></i>
            </div>
            <div class="min-w-0">
              <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ data.organization_name }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 font-mono">{{ data.organization_code }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('candidate_search.organization_summary')" style="width: 240px">
        <template #body="{ data }">
          <div class="text-xs text-gray-600 dark:text-gray-300">
            <p class="font-mono">{{ data.summary_code }}</p>
            <p class="text-gray-400 dark:text-gray-500 truncate" :title="data.summary_decree_no">{{ data.summary_decree_no }}</p>
          </div>
        </template>
      </Column>

      <Column :header="t('candidate_search.candidates')" style="width: 120px">
        <template #body="{ data }">
          <Tag
            :value="data.candidate_count"
            :severity="data.candidate_count > 0 ? 'success' : 'secondary'"
            class="!text-xs !px-2 !py-0.5"
          />
        </template>
      </Column>

      <!-- S-3: internal candidates eligible (Career Intelligence) -->
      <Column :header="t('candidate_search.internal_candidates')" style="width: 170px">
        <template #body="{ data }">
          <Tag
            :value="data.internal_candidate_count || 0"
            :severity="data.internal_candidate_count > 0 ? 'warn' : 'secondary'"
            class="!text-xs !px-2 !py-0.5"
            icon="pi pi-user"
          />
        </template>
      </Column>

      <template #expansion="{ data }">
        <div class="px-5 py-3 bg-gray-50 dark:bg-gray-800/50 space-y-3">
          <!-- Internal candidates (S-3/S-4) -->
          <div>
            <p class="text-xs font-semibold text-violet-600 dark:text-violet-400 uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <i class="pi pi-user text-xs"></i>{{ t('candidate_search.internal_candidates') }}
            </p>
            <div v-if="!data.internal_candidates || data.internal_candidates.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-1">
              {{ t('candidate_search.no_internal_candidates') }}
            </div>
            <DataTable
              v-else
              :value="data.internal_candidates"
              size="small"
              class="!text-sm border border-violet-200 dark:border-violet-700/50 rounded-lg overflow-hidden"
            >
              <Column :header="t('candidate_search.employee')">
                <template #body="{ data: ic }">
                  <div class="flex items-center gap-2">
                    <div class="w-7 h-7 rounded-full bg-violet-100 dark:bg-violet-900/40 text-violet-700 dark:text-violet-300 flex items-center justify-center text-[11px] font-semibold shrink-0">
                      {{ internalInitials(ic) }}
                    </div>
                    <div class="min-w-0">
                      <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ ic.name || '—' }}</p>
                      <p class="text-xs text-gray-400 dark:text-gray-500 truncate font-mono">{{ ic.employee_id }}</p>
                    </div>
                  </div>
                </template>
              </Column>
              <Column :header="t('candidate_search.current_position')">
                <template #body="{ data: ic }">
                  <span class="text-gray-600 dark:text-gray-300">{{ ic.current_position_name || '—' }}</span>
                </template>
              </Column>
              <Column :header="t('candidate_search.step_sequence')" style="width: 120px">
                <template #body="{ data: ic }">
                  <Tag :value="ic.source_step_sequence ?? 0" severity="info" class="!text-xs !px-2 !py-0.5" />
                </template>
              </Column>
            </DataTable>
          </div>

          <!-- External candidates (pool recruitment) -->
          <div>
            <p class="text-xs font-semibold text-sky-600 dark:text-sky-400 uppercase tracking-wider mb-2 flex items-center gap-1.5">
              <i class="pi pi-users text-xs"></i>{{ t('candidate_search.external_candidates') }}
            </p>
            <div v-if="data.candidates.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-1">
              {{ t('candidate_search.no_candidates') }}
            </div>
            <DataTable
              v-else
              :value="data.candidates"
              size="small"
              class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
            >
            <Column :header="t('candidate_search.candidates')">
              <template #body="{ data: c }">
                <div class="flex items-center gap-2">
                  <div class="w-7 h-7 rounded-full bg-sky-100 dark:bg-sky-900/40 text-sky-700 dark:text-sky-300 flex items-center justify-center text-[11px] font-semibold shrink-0">
                    {{ initials(c) }}
                  </div>
                  <div class="min-w-0">
                    <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ c.first_name }} {{ c.last_name }}</p>
                    <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ c.email }}</p>
                  </div>
                </div>
              </template>
            </Column>
            <Column :header="t('candidate_search.current_title')">
              <template #body="{ data: c }">
                <span class="text-gray-600 dark:text-gray-300">{{ c.current_title || '—' }}</span>
              </template>
            </Column>
            <Column :header="t('candidate_search.applied_for')">
              <template #body="{ data: c }">
                <span class="text-gray-600 dark:text-gray-300">{{ c.requisition_title || '—' }}</span>
              </template>
            </Column>
            <Column :header="t('common.status')" style="width: 140px">
              <template #body="{ data: c }">
                <Tag
                  :value="t('candidate_search.status_' + c.application_status)"
                  :severity="candidateStatusSeverity(c.application_status)"
                  class="!text-xs !px-1.5 !py-0.5"
                />
              </template>
            </Column>
          </DataTable>
          </div>
        </div>
      </template>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const expandedRows = ref([])
const searchTerm = ref('')

const skeletonColumns = [
  { field: 'position', header: 'Position', width: '40%' },
  { field: 'summary', header: 'Organization Summary', width: '30%' },
  { field: 'candidates', header: 'Candidates', width: '15%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function initials(c) {
  return ((c.first_name || '?').charAt(0) + (c.last_name || '').charAt(0)).toUpperCase()
}

function internalInitials(ic) {
  return (ic.name || '?').charAt(0).toUpperCase()
}

function candidateStatusSeverity(status) {
  switch (status) {
    case 'ACCEPTED': return 'success'
    case 'OFFERED': return 'help'
    case 'SHORTLISTED':
    case 'INTERVIEWED': return 'warn'
    case 'REJECTED':
    case 'WITHDRAWN': return 'danger'
    case 'SCREENED':
    case 'NEW':
    default: return 'info'
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (searchTerm.value) params.search = searchTerm.value
    const res = await api.get('/api/v1/tenant/workforce-intelligence/candidate-search', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  expandedRows.value = []
  loadData()
}

function onSearch() {
  currentPage.value = 1
  expandedRows.value = []
  loadData()
}

onMounted(loadData)
</script>
