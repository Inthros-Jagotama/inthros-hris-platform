<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.scores') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.score_description') }}</p>
    </div>

    <!-- Loading / Empty -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <i class="pi pi-spin pi-spinner text-emerald-500 text-2xl"></i>
    </div>

    <template v-else-if="score">
      <!-- Score Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
          <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.value_with_financial') }}</div>
          <div class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{{ formatNumber(score.job_value_with_financial) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
          <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.value_without_financial') }}</div>
          <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ formatNumber(score.job_value_without_financial) }}</div>
        </div>
        <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
          <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.has_financial_authority') }}</div>
          <Tag :value="score.has_financial_authority ? t('common.yes') : t('common.no')" :severity="score.has_financial_authority ? 'success' : 'danger'" class="!text-xs" />
          <div v-if="score.calculated_at" class="text-[10px] text-gray-400 mt-2">{{ t('job_management.calculated_at') }}: {{ score.calculated_at }}</div>
        </div>
      </div>

      <!-- Component Breakdown -->
      <div v-if="parsedComponents" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden">
        <div class="px-5 py-3 border-b border-gray-200 dark:border-gray-700 font-semibold text-sm text-gray-700 dark:text-gray-300">
          {{ t('job_management.component_breakdown') }}
        </div>
        <div class="p-5">
          <div v-for="(points, comp) in parsedComponents" :key="comp" class="flex items-center justify-between py-1.5 border-b border-gray-100 dark:border-gray-700 last:border-0">
            <span class="text-sm text-gray-700 dark:text-gray-300 capitalize">{{ comp.replace(/_/g, ' ') }}</span>
            <span class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{ formatNumber(points) }}</span>
          </div>
        </div>
      </div>
    </template>

    <!-- No score yet -->
    <div v-else>
      <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
        <i class="pi pi-calculator text-4xl mb-3 opacity-50"></i>
        <p class="text-sm font-medium">{{ t('job_management.no_score') }}</p>
        <p class="text-xs mt-1">{{ t('job_management.score_hint') }}</p>
      </div>
    </div>

    <!-- Actions -->
    <div class="flex justify-end gap-3">
      <Button :label="t('common.refresh')" icon="pi pi-refresh" size="small" text @click="loadData" />
      <Button v-if="score" :label="t('job_management.recalculate')" icon="pi pi-calculator" size="small" severity="info" :loading="calculating" @click="recalculate" />
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'

const props = defineProps({ orgId: String })
const emit = defineEmits(['saved'])
const { t } = useI18n()
const toast = useToast()
const apiBase = '/api/v1/tenant/job-management/scores'

const loading = ref(false)
const calculating = ref(false)
const score = ref(null)

const parsedComponents = computed(() => {
  if (!score.value?.components) return null
  try { return JSON.parse(score.value.components) }
  catch { return null }
})

function formatNumber(n) {
  return n?.toLocaleString?.('id-ID') ?? '-'
}

async function loadData() {
  if (!props.orgId) return
  loading.value = true
  try {
    const res = await api.get(`${apiBase}/${props.orgId}`)
    score.value = res.data?.data || null
    emit('saved')
  } catch {
    score.value = null
  } finally {
    loading.value = false
  }
}

async function recalculate() {
  if (!props.orgId) return
  calculating.value = true
  try {
    // Recalculate: upsert with empty body triggers server-side recalculation
    const res = await api.put(`${apiBase}/${props.orgId}`, { components: null })
    score.value = res.data?.data || null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_management.score_calculated'), life: 2000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    calculating.value = false
  }
}

onMounted(loadData)
</script>
