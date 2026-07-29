<script setup>
defineProps({
    value: String
});
</script>

<template>
    <label class="col-span-3 pt-2 text-sm font-medium text-surface-700 dark:text-surface-0/80">
        <span v-if="value">{{ value }}</span>
        <span v-else><slot /></span>
    </label>
</template>
