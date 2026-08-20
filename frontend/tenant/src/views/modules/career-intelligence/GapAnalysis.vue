<template>
  <div class="space-y-4">
    <!-- Filter -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 flex items-center gap-1.5">
        <i class="pi pi-info-circle text-xs"></i>{{ t('gap_analysis.subtitle') }}
      </p>
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
        <FormRow :label="t('gap_analysis.employee')" required>
          <SelectLabel
            v-model="form.employee_id"
            :options="employees"
            optionLabel="name"
            optionValue="id"
            :placeholder="t('common.select')"
          />
        </FormRow>
        <FormRow :label="t('gap_analysis.target_title')" required>
          <SelectLabel
            v-model="form.target_title_id"
            :options="jobTitles"
            optionLabel="name"
            optionValue="id"
            :placeholder="t('common.select')"
          />
        </FormRow>
      </div>
      <div class="flex justify-end mt-3">
        <Button
          :label="t('gap_analysis.analyze')"
          icon="pi pi-search"
          size="small"
          :loading="loading"
          :disabled="!form.employee_id || !form.target_title_id"
          @click="analyze"
        />
      </div>
    </div>

    <SkeletonCard v-if="loading" type="stat" :count="4" />

    <!-- Empty state (belum pernah analyze) -->
    <div
      v-else-if="!result && !loadFailed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-sitemap text-4xl mb-3 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('gap_analysis.empty') }}</p>
    </div>

    <!-- Error state -->
    <div
      v-else-if="loadFailed"
      class="bg-white dark:bg-gray-800 rounded-lg border border-rose-200 dark:border-rose-700/50 p-10 flex flex-col items-center justify-center text-gray-400 dark:text-gray-500"
    >
      <i class="pi pi-exclamation-triangle text-4xl mb-3 text-rose-400 opacity-70"></i>
      <p class="text-sm font-medium text-rose-600 dark:text-rose-400">{{ t('message.failed_to_load') }}</p>
      <Button :label="t('common.retry')" icon="pi pi-refresh" size="small" class="!mt-4 !text-xs" @click="analyze" />
    </div>

    <!-- Result -->
    <template v-else-if="result">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('gap_analysis.target_title') }}</p>
          <p class="text-lg font-bold text-navy-800 dark:text-gray-100 truncate">{{ result.target_title || '-' }}</p>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('gap_analysis.matched_skills') }}</p>
          <p class="text-lg font-bold text-emerald-600 dark:text-emerald-400">{{ result.matched_skills }} / {{ result.total_required }}</p>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('gap_analysis.gap_percentage') }}</p>
          <p class="text-lg font-bold" :class="gapColorClass">{{ result.gap_percentage.toFixed(0) }}%</p>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('gap_analysis.estimated_timeline') }}</p>
          <p class="text-lg font-bold text-navy-800 dark:text-gray-100">{{ result.estimated_timeline || '-' }}</p>
        </div>
      </div>

      <!-- Progress bar kesenjangan -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center justify-between mb-2">
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ t('gap_analysis.readiness') }}</p>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400">{{ (100 - result.gap_percentage).toFixed(0) }}%</p>
        </div>
        <div class="w-full h-2 rounded-full bg-gray-100 dark:bg-gray-700 overflow-hidden">
          <div class="h-full rounded-full transition-all" :class="gapBarClass" :style="{ width: (100 - result.gap_percentage) + '%' }"></div>
        </div>
      </div>

      <!-- Rekomendasi -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100">{{ t('gap_analysis.recommendations') }}</h3>
        </div>
        <div v-if="!result.recommendations?.length" class="p-6 text-center text-sm text-gray-400 dark:text-gray-500">
          {{ t('gap_analysis.no_recommendations') }}
        </div>
        <div v-else class="divide-y divide-gray-100 dark:divide-gray-700">
          <div v-for="(rec, idx) in result.recommendations" :key="idx" class="p-4 flex items-start gap-3">
            <div class="w-8 h-8 rounded-lg shrink-0 flex items-center justify-center" :class="categoryTint(rec.category)">
              <i :class="categoryIcon(rec.category)" class="text-xs"></i>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">{{ rec.category }}</span>
                <Tag :value="rec.priority" :severity="prioritySeverity(rec.priority)" class="!text-[10px] !px-1.5 !py-0.5" />
              </div>
              <p class="text-sm text-navy-800 dark:text-gray-100 mt-0.5">{{ rec.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(false)
const loadFailed = ref(false)
const result = ref(null)

const employees = ref([])
const jobTitles = ref([])

const form = ref({ employee_id: null, target_title_id: null })

async function loadReferences() {
  const [empRes, titleRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/job-management/titles', { params: { per_page: 500 } })
  ])
  employees.value = empRes.value?.data?.data || []
  jobTitles.value = titleRes.value?.data?.data || []
}

async function analyze() {
  if (!form.value.employee_id || !form.value.target_title_id) return
  loading.value = true
  loadFailed.value = false
  try {
    const res = await api.get('/api/v1/tenant/career-intelligence/paths/gap-analysis', {
      params: { employee_id: form.value.employee_id, target_title_id: form.value.target_title_id }
    })
    result.value = res.data?.data || null
  } catch (e) {
    loadFailed.value = true
    result.value = null
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

const gapColorClass = computed(() => {
  const g = result.value?.gap_percentage ?? 0
  if (g >= 60) return 'text-rose-600 dark:text-rose-400'
  if (g >= 30) return 'text-amber-600 dark:text-amber-400'
  return 'text-emerald-600 dark:text-emerald-400'
})

const gapBarClass = computed(() => {
  const g = result.value?.gap_percentage ?? 0
  if (g >= 60) return 'bg-rose-500'
  if (g >= 30) return 'bg-amber-500'
  return 'bg-emerald-500'
})

function prioritySeverity(priority) {
  switch (priority) {
    case 'HIGH': return 'danger'
    case 'MEDIUM': return 'warn'
    case 'LOW': return 'info'
    default: return 'secondary'
  }
}

function categoryTint(category) {
  switch (category) {
    case 'TRAINING': return 'bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400'
    case 'EXPERIENCE': return 'bg-violet-50 dark:bg-violet-500/10 text-violet-600 dark:text-violet-400'
    case 'CERTIFICATION': return 'bg-amber-50 dark:bg-amber-500/10 text-amber-600 dark:text-amber-400'
    default: return 'bg-gray-50 dark:bg-gray-700 text-gray-500 dark:text-gray-400'
  }
}

function categoryIcon(category) {
  switch (category) {
    case 'TRAINING': return 'pi pi-book'
    case 'EXPERIENCE': return 'pi pi-briefcase'
    case 'CERTIFICATION': return 'pi pi-verified'
    default: return 'pi pi-circle'
  }
}

onMounted(loadReferences)
</script>
