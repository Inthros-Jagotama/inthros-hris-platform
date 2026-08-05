<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Skeleton hanya saat pertama load; saat refresh nilai lama tetap tampil -->
    <div v-if="loading && !score" v-for="n in 3" :key="n" class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4 animate-pulse">
      <div class="h-3 w-24 bg-gray-200 dark:bg-gray-700 rounded mb-2"></div>
      <div class="h-7 w-16 bg-gray-200 dark:bg-gray-700 rounded"></div>
    </div>

    <template v-else>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.value_with_financial') }}</div>
        <div class="text-2xl font-bold text-emerald-600 dark:text-emerald-400">{{ formatNumber(score?.job_value_with_financial) }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.value_without_financial') }}</div>
        <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ formatNumber(score?.job_value_without_financial) }}</div>
      </div>
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-4">
        <div class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-1">{{ t('job_management.has_financial_authority') }}</div>
        <Tag
          :value="score ? (score.has_financial_authority ? t('common.yes') : t('common.no')) : '-'"
          :severity="score?.has_financial_authority ? 'success' : 'danger'"
          class="!text-xs"
        />
        <!-- Badge status kelengkapan skor (is_complete & completed_at dari API) -->
        <div class="mt-3 flex flex-wrap items-center gap-2">
          <Tag
            v-if="score"
            :value="score.is_complete ? t('job_management.score_complete') : t('job_management.score_incomplete')"
            :severity="score.is_complete ? 'success' : 'warning'"
            :icon="score.is_complete ? 'pi pi-check-circle' : 'pi pi-exclamation-triangle'"
            class="!text-xs"
          />
          <span v-if="score?.is_complete && score.completed_at" class="text-[10px] text-gray-400">{{ t('job_management.completed_at') }}: {{ score.completed_at }}</span>
        </div>
        <div v-if="score?.calculated_at" class="text-[10px] text-gray-400 mt-2">{{ t('job_management.calculated_at') }}: {{ score.calculated_at }}</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Tag from 'primevue/tag'

const props = defineProps({ orgId: String })
const { t } = useI18n()
const apiBase = '/api/v1/tenant/job-management/scores/org'

const loading = ref(true)
const score = ref(null)

function formatNumber(n) {
  return n?.toLocaleString?.('id-ID') ?? '-'
}

async function refresh() {
  if (!props.orgId) return
  loading.value = true
  try {
    const res = await api.get(`${apiBase}/${props.orgId}`)
    score.value = res.data?.data || null
  } catch {
    score.value = null
  } finally {
    loading.value = false
  }
}

defineExpose({ refresh })

watch(() => props.orgId, refresh)
onMounted(refresh)
</script>
