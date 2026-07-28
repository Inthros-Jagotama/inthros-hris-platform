<script setup>
defineProps({
  columns: { type: Array, default: () => [] },
  rows: { type: Number, default: 5 }
})
</script>

<template>
  <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden animate-pulse">
    <table class="w-full">
      <!-- Header -->
      <thead>
        <tr class="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
          <th v-for="(col, i) in columns" :key="i" class="px-2 py-2.5">
            <div class="h-3 bg-gray-200 dark:bg-gray-600 rounded" :class="col.headerWidth || 'w-20'"></div>
          </th>
        </tr>
      </thead>
      <!-- Body -->
      <tbody>
        <tr v-for="row in rows" :key="row" class="border-b border-gray-100 dark:border-gray-700 last:border-0">
          <td v-for="(col, i) in columns" :key="i" class="px-2 py-3">
            <!-- Compound: two stacked bars -->
            <div v-if="col.type === 'compound'" class="space-y-1.5">
              <div class="h-3 bg-gray-100 dark:bg-gray-700 rounded w-3/4"></div>
              <div class="h-2.5 bg-gray-100 dark:bg-gray-700 rounded w-full"></div>
            </div>
            <!-- Tag: rounded bar -->
            <div v-else-if="col.type === 'tag'">
              <div class="h-5 bg-gray-100 dark:bg-gray-700 rounded-full" :class="col.width || 'w-16'"></div>
            </div>
            <!-- Icons: multiple small circles -->
            <div v-else-if="col.type === 'icons'" class="flex items-center gap-1">
              <div v-for="j in (col.count || 3)" :key="j" class="w-6 h-6 bg-gray-100 dark:bg-gray-700 rounded-full"></div>
            </div>
            <!-- Checkbox: small square -->
            <div v-else-if="col.type === 'checkbox'" class="w-4 h-4 bg-gray-100 dark:bg-gray-700 rounded mx-auto"></div>
            <!-- Key + Copy: code style bar + button -->
            <div v-else-if="col.type === 'key-copy'" class="flex items-center gap-1">
              <div class="h-5 bg-gray-100 dark:bg-gray-700 rounded flex-1" :class="col.width || 'w-28'"></div>
              <div class="w-5 h-5 bg-gray-100 dark:bg-gray-700 rounded"></div>
            </div>
            <!-- Text: single bar (default) -->
            <div v-else>
              <div class="h-3 bg-gray-100 dark:bg-gray-700 rounded" :class="col.width || 'w-20'"></div>
            </div>
          </td>
        </tr>
        <!-- Spacer row for shadow effect -->
        <tr v-if="rows > 3">
          <td :colspan="columns.length" class="h-1"></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
