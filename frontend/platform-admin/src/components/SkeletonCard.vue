<template>
  <div class="grid gap-3 animate-pulse" :class="gridCols">
    <div
      v-for="i in resolvedCount"
      :key="i"
      class="bg-white rounded-lg border border-gray-200"
      :class="padding"
    >
      <!-- KPI Card: icon box + value bar + label bar -->
      <template v-if="type === 'kpi'">
        <div class="w-8 h-8 rounded-lg bg-gray-200 mb-2"></div>
        <div class="h-5 rounded bg-gray-200 mb-1" :class="valueWidth"></div>
        <div class="h-3 rounded bg-gray-200" :class="labelWidth"></div>
      </template>

      <!-- Stat Card: label bar + value bar (no icon) -->
      <template v-else-if="type === 'stat'">
        <div class="h-3 rounded bg-gray-200 mb-2" :class="labelWidth"></div>
        <div class="h-5 rounded bg-gray-200" :class="valueWidth"></div>
      </template>

      <!-- Metric Card: big value bar + small label bar + mini trend indicator -->
      <template v-else-if="type === 'metric'">
        <div class="h-4 rounded bg-gray-200 mb-3" :class="valueWidth || 'w-24'"></div>
        <div class="h-6 rounded bg-gray-200 mb-1" :class="labelWidth || 'w-16'"></div>
        <div class="h-2 w-12 rounded bg-gray-100"></div>
      </template>

      <!-- Alert Card: icon + two-line content -->
      <template v-else-if="type === 'alert'">
        <div class="flex items-center gap-2">
          <div class="w-5 h-5 rounded-full bg-gray-200 shrink-0"></div>
          <div class="flex-1 space-y-1.5">
            <div class="h-3 rounded bg-gray-200" :class="labelWidth || 'w-full'"></div>
            <div class="h-3 rounded bg-gray-100 w-3/4"></div>
          </div>
        </div>
      </template>

      <!-- Sparkline Card: mini chart skeleton (dashboard pool, monitoring utilization) -->
      <template v-else-if="type === 'sparkline'">
        <div class="flex items-end gap-1 h-20 mb-2">
          <div
            v-for="h in sparkHeights"
            :key="h"
            class="flex-1 bg-gray-100 rounded-t"
            :style="{ height: h + '%' }"
          ></div>
        </div>
        <div class="h-3 rounded bg-gray-200" :class="labelWidth || 'w-16'"></div>
      </template>

      <!-- Detail Card: title bar + multiple row bars -->
      <template v-else-if="type === 'detail'">
        <div class="h-4 rounded bg-gray-200 mb-3" :class="valueWidth || 'w-32'"></div>
        <div class="space-y-2.5">
          <div v-for="j in (rows || 4)" :key="j" class="flex items-center justify-between">
            <div class="h-3 rounded bg-gray-200 w-20"></div>
            <div class="h-3 rounded bg-gray-200 w-12"></div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'kpi',
    validator: (val) => ['kpi', 'stat', 'metric', 'alert', 'sparkline', 'detail'].includes(val)
  },
  count: {
    type: Number,
    default: null
  },
  cols: {
    type: String,
    default: null
  },
  rows: {
    type: Number,
    default: 4
  },
  padding: {
    type: String,
    default: 'p-3'
  },
  valueWidth: {
    type: String,
    default: null
  },
  labelWidth: {
    type: String,
    default: null
  },
  sparkHeights: {
    type: Array,
    default: () => [40, 60, 30, 70, 45, 55, 35, 65]
  }
})

const resolvedCount = computed(() => {
  if (props.count !== null) return props.count
  if (props.type === 'kpi') return 6
  if (props.type === 'stat') return 4
  if (props.type === 'alert') return 3
  return 4
})

const gridCols = computed(() => {
  if (props.cols !== null) return props.cols
  if (props.type === 'kpi') return 'grid-cols-2 md:grid-cols-4 lg:grid-cols-6'
  if (props.type === 'stat') return 'grid-cols-2 md:grid-cols-4'
  return 'grid-cols-1'
})
</script>
