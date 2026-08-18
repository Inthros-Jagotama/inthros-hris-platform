<template>
  <div class="space-y-3 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
    <div>
      <h3 class="text-base font-semibold text-navy-800 dark:text-gray-100">{{ t('job_values.cluster_mapping_title') }}</h3>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_values.cluster_mapping_description') }}</p>
    </div>

    <SkeletonCard v-if="loading" type="detail" :count="1" :rows="2" cols="grid-cols-1" padding="p-3" />

    <template v-else>
      <SelectLabel
        v-model="selectedClusters"
        :options="clusterOptions"
        option-label="label"
        option-value="value"
        :placeholder="t('job_values.cluster_mapping_placeholder')"
        showClear
        multiple
      />

      <!-- Daftar cluster yang sudah di-mapping -->
      <div v-if="selectedClusters.length" class="flex flex-wrap gap-1.5">
        <span
          v-for="c in selectedClusters"
          :key="c"
          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-medium bg-emerald-50 dark:bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800"
        >
          <i class="pi pi-tag text-[10px]"></i>
          {{ c }}
        </span>
      </div>

      <div class="flex justify-end">
        <Button
          :label="t('job_values.cluster_mapping_save')"
          icon="pi pi-check"
          size="small"
          :loading="saving"
          :disabled="saving"
          @click="handleSave"
        />
      </div>
      <div
        v-if="errorMsg"
        class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2"
      >
        {{ errorMsg }}
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'

const props = defineProps({
  type: { type: String, required: true }
})

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const errorMsg = ref('')
const clusterOptions = ref([])
const selectedClusters = ref([])

async function loadData() {
  loading.value = true
  try {
    const [compRes, mapRes] = await Promise.all([
      api.get('/api/v1/tenant/settings/competencies', { params: { per_page: 500 } }),
      api.get(`/api/v1/tenant/job-management/values/clusters/${props.type}`)
    ])
    const all = compRes.data?.data || []
    const clusters = [...new Set(all.map(c => c.cluster).filter(Boolean))].sort()
    clusterOptions.value = clusters.map(c => ({ label: c, value: c }))
    selectedClusters.value = mapRes.data?.data?.clusters || []
  } catch {
    clusterOptions.value = []
    selectedClusters.value = []
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  errorMsg.value = ''
  try {
    await api.put(`/api/v1/tenant/job-management/values/clusters/${props.type}`, {
      clusters: selectedClusters.value
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.cluster_mapping_saved'), life: 3000 })
    await loadData()
  } catch (e) {
    errorMsg.value = e.response?.data?.error?.message || e.message || t('message.operation_failed')
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
