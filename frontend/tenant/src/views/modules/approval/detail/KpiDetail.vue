<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('kpi.employee')" :value="props.documentDetail?.employee_name" class="col-span-2" />
    <ViewLabel :label="t('kpi.organization')" :value="props.documentDetail?.organization_name" />
    <ViewLabel :label="t('kpi.period')" :value="props.documentDetail?.period_code" />
  </div>

  <div v-if="props.documentDetail?.details?.length">
    <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.indicators') }}</p>
    <div class="space-y-1">
      <div v-for="d in props.documentDetail.details" :key="d.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
        <div class="font-medium text-gray-700 dark:text-gray-200">{{ d.indicator_name }}</div>
        <div class="text-gray-500 dark:text-gray-400">
          {{ t('kpi.target') }}: {{ d.target }} {{ d.unit_of_measurement }}
          <span v-if="d.actual"> · {{ t('kpi.actual') }}: {{ d.actual }}</span>
        </div>
      </div>
    </div>
  </div>

  <div v-if="props.documentDetail?.program_items?.length">
    <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.program_items') }}</p>
    <div class="space-y-1">
      <div v-for="p in props.documentDetail.program_items" :key="p.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
        <div class="font-medium text-gray-700 dark:text-gray-200">{{ p.title }}</div>
        <div class="text-gray-500 dark:text-gray-400">
          {{ t('kpi.target') }}: {{ p.target }} {{ p.unit_of_measurement }}
          <span v-if="p.actual"> · {{ t('kpi.actual') }}: {{ p.actual }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from '@/composables/useI18n'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t } = useI18n()
</script>
