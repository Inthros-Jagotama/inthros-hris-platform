<template>
  <div class="border border-gray-200 rounded-lg overflow-hidden animate-pulse">
    <!-- Header Row -->
    <div class="flex items-center gap-3 px-4 py-2.5 bg-gray-50 border-b border-gray-200">
      <div
        v-for="(col, idx) in columns"
        :key="'h-' + idx"
        class="h-3 bg-gray-200 rounded shrink-0"
        :class="col.headerWidth || col.width || 'w-20'"
      />
    </div>
    <!-- Data Rows -->
    <div
      v-for="i in rows"
      :key="i"
      class="flex items-center gap-3 px-4 py-3 border-b border-gray-100 last:border-b-0"
      :class="{ 'bg-gray-50/50': i % 2 === 0 }"
    >
      <template v-for="(col, idx) in columns" :key="'r-' + idx + '-' + i">
        <!-- Checkbox column -->
        <div v-if="col.type === 'checkbox'" class="h-4 w-4 bg-gray-200 rounded shrink-0" />

        <!-- Compound column (2-line cell: name+desc, name+email, etc.) -->
        <div v-else-if="col.type === 'compound'" class="flex-1 min-w-0 space-y-1.5">
          <div class="h-3.5 bg-gray-200 rounded" :class="col.widths?.[0] || 'w-24'" />
          <div class="h-3 bg-gray-100 rounded" :class="col.widths?.[1] || 'w-36'" />
        </div>

        <!-- Tag / Badge column (rounded-full) -->
        <div v-else-if="col.type === 'tag'" class="h-5 bg-gray-200 rounded-full shrink-0" :class="col.width || 'w-16'" />

        <!-- Action Icons column -->
        <div v-else-if="col.type === 'icons'" class="flex items-center gap-1 shrink-0">
          <div v-for="j in (col.count || 1)" :key="j" class="h-6 w-6 bg-gray-200 rounded shrink-0" />
        </div>

        <!-- Name+Tag compound (name bar + system tag beside it) -->
        <div v-else-if="col.type === 'name-tag'" class="flex items-center gap-2 shrink-0">
          <div class="h-3.5 bg-gray-200 rounded" :class="col.width || 'w-20'" />
          <div v-if="col.showTag" class="h-4 w-12 bg-gray-200 rounded-full shrink-0" />
        </div>

        <!-- Key+Copy compound (key bar + copy icon) -->
        <div v-else-if="col.type === 'key-copy'" class="flex items-center gap-1.5 shrink-0">
          <div class="h-4 bg-gray-200 rounded" :class="col.width || 'w-24'" />
          <div class="h-5 w-5 bg-gray-200 rounded shrink-0" />
        </div>

        <!-- Tags group (multiple small tags + overflow) -->
        <div v-else-if="col.type === 'tags'" class="flex items-center gap-1 shrink-0">
          <div v-for="j in (col.count || 2)" :key="j" class="h-4 bg-gray-200 rounded-full shrink-0" :class="col.tagWidth || 'w-12'" />
          <div v-if="col.showOverflow" class="h-3 w-6 bg-gray-200 rounded shrink-0" />
        </div>

        <!-- Plain text bar -->
        <div v-else class="h-3.5 bg-gray-200 rounded shrink-0" :class="col.width || 'w-20'" />
      </template>
    </div>
  </div>
</template>

<script setup>
defineProps({
  columns: {
    type: Array,
    required: true,
    validator: (val) => val.length > 0
  },
  rows: {
    type: Number,
    default: 6
  }
})
</script>
