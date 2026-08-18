<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('okr.employee')" :value="props.documentDetail?.employee_name" class="col-span-2" />
    <ViewLabel :label="t('okr.organization')" :value="props.documentDetail?.organization_name" />
    <ViewLabel :label="t('okr.period')" :value="props.documentDetail?.period_code" />
  </div>

  <div v-if="okrObjectiveGroups.length">
    <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('okr.objectives') }}</p>
    <div v-for="g in okrObjectiveGroups" :key="g.key" class="mb-2">
      <div class="text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{{ g.title }}</div>
      <div class="space-y-1">
        <div v-for="kr in g.items" :key="kr.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
          <div class="font-medium text-gray-700 dark:text-gray-200">{{ kr.key_result_title }}</div>
          <div class="text-gray-500 dark:text-gray-400">
            {{ t('okr.target') }}: {{ kr.target_value }} {{ kr.unit }}
            <span v-if="kr.actual_value"> · {{ t('okr.actual') }}: {{ kr.actual_value }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t } = useI18n()

const okrObjectiveGroups = computed(() => {
  const details = props.documentDetail?.details || []
  const groups = {}
  const order = []
  for (const d of details) {
    const key = d.objective_id || d.objective_title
    if (!groups[key]) {
      groups[key] = { key, title: d.objective_title, items: [] }
      order.push(key)
    }
    groups[key].items.push(d)
  }
  return order.map(key => groups[key])
})
</script>
