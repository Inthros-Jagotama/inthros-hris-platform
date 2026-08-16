<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.results') }}</h2>
      </div>
    </div>

    <!-- Employee picker -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center gap-2 flex-wrap">
        <Button v-if="!employeeId" :label="t('competency_360.my_result')" icon="pi pi-user" size="small" severity="secondary" outlined @click="useMy()" />
        <Select v-model="selectedEmployee" :options="employeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-80" :placeholder="t('competency_360.select_employee')" @change="loadAll" />
        <span v-if="employeeId" class="text-xs text-gray-500 dark:text-gray-400">{{ employeeName(employeeId) }}</span>
      </div>
    </div>

    <SkeletonCard v-if="loading" :rows="6" />
    <template v-else-if="error">
      <Message severity="warn" :closable="false">{{ error }}</Message>
    </template>
    <template v-else-if="result">
      <!-- Summary cards -->
      <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-2">
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.overall_score') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ result.overall_score ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.total_gap') }}</p>
          <p :class="[ (result.total_gap ?? 0) <= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-rose-600 dark:text-rose-400']" class="text-xl font-bold">{{ result.total_gap ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.self_score') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ result.self_score ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.others_score') }}</p>
          <p class="text-xl font-bold text-gray-800 dark:text-gray-100">{{ result.others_score ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.perception_gap') }}</p>
          <p :class="[ (result.perception_gap ?? 0) <= 0 ? 'text-emerald-600 dark:text-emerald-400' : 'text-amber-600 dark:text-amber-400']" class="text-xl font-bold">{{ result.perception_gap ?? 0 }}</p>
        </div>
        <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-[11px] font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ t('competency_360.gap_report') }}</p>
          <Button icon="pi pi-chart-line" size="small" text severity="info" :label="t('common.view')" @click="showGap = !showGap" />
        </div>
      </div>

      <!-- Gap analysis -->
      <div v-if="showGap && gap" class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <h3 class="text-sm font-semibold text-emerald-600 dark:text-emerald-400 mb-3 flex items-center gap-2"><i class="pi pi-arrow-up"></i>{{ t('competency_360.strengths') }}</h3>
          <div v-if="gap.strengths?.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_data') }}</div>
          <div v-for="s in gap.strengths" :key="s.competency_id" class="flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0">
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ s.competency_name || s.competency_id }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ s.score }} / {{ s.required_level }}</span>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
          <h3 class="text-sm font-semibold text-rose-600 dark:text-rose-400 mb-3 flex items-center gap-2"><i class="pi pi-arrow-down"></i>{{ t('competency_360.development_areas') }}</h3>
          <div v-if="gap.development_areas?.length === 0" class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.no_data') }}</div>
          <div v-for="d in gap.development_areas" :key="d.competency_id" class="flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0">
            <span class="text-sm text-gray-700 dark:text-gray-200">{{ d.competency_name || d.competency_id }}</span>
            <span class="text-xs text-gray-500 dark:text-gray-400">{{ d.score }} / {{ d.required_level }} ({{ d.gap }})</span>
          </div>
        </div>
      </div>

      <!-- Per-competency table -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3">{{ t('competency_360.competency_scores') }}</h3>
        <DataTable :value="result.competencies || []" size="small" class="!text-sm p-datatable-sm">
          <template #empty>
            <div class="text-center py-6 text-gray-400 dark:text-gray-500 text-sm">{{ t('competency_360.no_data') }}</div>
          </template>
          <Column field="competency_name" :header="t('competency_360.competency')">
            <template #body="{data}"><span class="text-gray-800 dark:text-gray-100">{{ data.competency_name || data.competency_id }}</span></template>
          </Column>
          <Column field="score" :header="t('competency_360.score')" style="width:80px">
            <template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ data.score }}</span></template>
          </Column>
          <Column field="required_level" :header="t('competency_360.required_level')" style="width:110px">
            <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.required_level }}</span></template>
          </Column>
          <Column field="gap" :header="t('competency_360.gap')" style="width:80px">
            <template #body="{data}">
              <Tag :value="String(data.gap)" :severity="data.gap >= 0 ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column field="rater_scores" :header="t('competency_360.rater_scores')" style="width:220px">
            <template #body="{data}">
              <div class="flex flex-wrap gap-1">
                <Tag v-if="data.rater_scores?.self" :value="`S:${data.rater_scores.self}`" severity="info" class="!text-[10px] !px-1 !py-0.5" />
                <Tag v-if="data.rater_scores?.superior" :value="`Sup:${data.rater_scores.superior}`" severity="success" class="!text-[10px] !px-1 !py-0.5" />
                <Tag v-if="data.rater_scores?.peer" :value="`P:${data.rater_scores.peer}`" severity="warn" class="!text-[10px] !px-1 !py-0.5" />
                <Tag v-if="data.rater_scores?.subordinate" :value="`Sub:${data.rater_scores.subordinate}`" severity="secondary" class="!text-[10px] !px-1 !py-0.5" />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Message from 'primevue/message'
import SkeletonCard from '@/components/SkeletonCard.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const employees = ref([])
const selectedEmployee = ref(null)
const loading = ref(false)
const error = ref('')
const result = ref(null)
const gap = ref(null)
const showGap = ref(false)

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.id })))

function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || id?.slice(0, 8)
}

async function loadReferences() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}

async function useMy() {
  selectedEmployee.value = await loadMyEmployeeId()
  await loadAll()
}

async function loadAll() {
  error.value = ''
  result.value = null
  gap.value = null
  showGap.value = false
  if (!selectedEmployee.value) return
  loading.value = true
  try {
    const [res, gapRes] = await Promise.allSettled([
      api.get(`/api/v1/tenant/competency/employees/${selectedEmployee.value}/result`),
      api.get(`/api/v1/tenant/competency/employees/${selectedEmployee.value}/gap`)
    ])
    if (res.status === 'fulfilled') {
      result.value = res.value.data?.data || res.value.data
    } else {
      error.value = getErrorMessage(res.reason, t('competency_360.no_result'))
    }
    if (gapRes.status === 'fulfilled') {
      gap.value = gapRes.value.data?.data || gapRes.value.data
    }
  } catch (e) {
    error.value = getErrorMessage(e, t('competency_360.no_result'))
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadReferences()
  await useMy()
})
</script>
