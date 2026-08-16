<script setup>
defineProps({
  title: { type: String, default: '' },
  icon: { type: String, default: '' },
  tint: { type: String, default: '' },
  empty: { type: String, default: '' },
  items: { type: Array, default: () => [] }
})
</script>

<template>
  <section class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
    <header class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex items-center gap-2">
      <i :class="[icon, tint]" class="text-sm"></i>
      <h2 class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ title }}</h2>
      <span class="text-xs text-gray-400 dark:text-gray-500">{{ items.length }}</span>
    </header>
    <div class="divide-y divide-gray-100 dark:divide-gray-800">
      <p v-if="!items.length" class="px-4 py-6 text-center text-sm text-gray-400 dark:text-gray-500">{{ empty }}</p>
      <div v-for="(item, i) in items" :key="i" class="px-4 py-3">
        <slot name="default" :item="item"></slot>
      </div>
    </div>
  </section>
</template>
