<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: { type: String, default: 'kpi' },
  count: { type: Number, default: null },
  cols: { type: String, default: null },
  rows: { type: Number, default: 4 },
  padding: { type: String, default: 'p-3' },
  valueWidth: { type: String, default: '' },
  labelWidth: { type: String, default: '' }
})

const resolvedCount = computed(() => {
  if (props.count !== null && props.count !== undefined) return props.count
  const defaults = { kpi: 6, stat: 4, metric: 4, alert: 3, sparkline: 4, detail: 4 }
  return defaults[props.type] || 4
})

const gridCols = computed(() => {
  if (props.cols) return props.cols
  const gridDefaults = {
    kpi: 'grid-cols-2 md:grid-cols-4 lg:grid-cols-6',
    stat: 'grid-cols-2 md:grid-cols-4',
    metric: 'grid-cols-1',
    alert: 'grid-cols-1',
    sparkline: 'grid-cols-1',
    detail: 'grid-cols-1'
  }
  return gridDefaults[props.type] || 'grid-cols-1'
})

const sparkHeights = [45, 60, 35, 70, 50, 55, 40, 65]
</script>

<template>
  <div class="grid gap-3 animate-pulse" :class="gridCols">
    <!-- KPI: icon box + value + label -->
    <template v-if="type === 'kpi'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="w-8 h-8 bg-gray-200 dark:bg-gray-600 rounded-lg mb-2"></div>
        <div class="h-5 w-3/4 bg-gray-200 dark:bg-gray-600 rounded mb-1"></div>
        <div class="h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded"></div>
      </div>
    </template>

    <!-- Stat: label bar + value bar -->
    <template v-else-if="type === 'stat'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-2"></div>
        <div class="h-6 w-2/3 bg-gray-200 dark:bg-gray-600 rounded"></div>
      </div>
    </template>

    <!-- Metric: big value + small label + trend bar -->
    <template v-else-if="type === 'metric'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="h-8 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-1"></div>
        <div class="h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded mb-2"></div>
        <div class="h-3 w-1/4 bg-gray-200 dark:bg-gray-600 rounded"></div>
      </div>
    </template>

    <!-- Alert: icon circle + two-line text -->
    <template v-else-if="type === 'alert'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 flex items-start gap-2">
        <div class="w-8 h-8 bg-gray-200 dark:bg-gray-600 rounded-full shrink-0"></div>
        <div class="flex-1 space-y-1.5">
          <div class="h-3 w-3/4 bg-gray-200 dark:bg-gray-600 rounded"></div>
          <div class="h-3 w-1/2 bg-gray-200 dark:bg-gray-600 rounded"></div>
        </div>
      </div>
    </template>

    <!-- Sparkline: mini bar chart -->
    <template v-else-if="type === 'sparkline'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="flex items-end gap-1 h-16 mb-2">
          <div v-for="h in sparkHeights" :key="h"
            class="flex-1 bg-gray-200 dark:bg-gray-600 rounded-t"
            :style="{ height: h + '%' }"></div>
        </div>
        <div class="h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded"></div>
      </div>
    </template>

    <!-- Detail: title + row pairs -->
    <template v-else-if="type === 'detail'">
      <div v-for="i in resolvedCount" :key="i" :class="padding"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
        <div class="h-4 w-1/3 bg-gray-200 dark:bg-gray-600 rounded mb-3"></div>
        <div v-for="j in rows" :key="j" class="flex items-center justify-between py-1.5">
          <div class="h-3 w-1/4 bg-gray-200 dark:bg-gray-600 rounded"></div>
          <div class="h-3 w-1/3 bg-gray-200 dark:bg-gray-600 rounded"></div>
        </div>
      </div>
    </template>
  </div>
</template>
