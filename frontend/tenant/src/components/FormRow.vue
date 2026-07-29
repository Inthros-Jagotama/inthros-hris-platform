<script setup>
defineProps({
    label: String,
    required: { type: Boolean, default: false },
    errors: { type: String, default: '' }
});
</script>
<template>
    <div class="flex items-start gap-3">
        <label class="w-36 shrink-0 pt-2 text-sm font-medium text-surface-700 dark:text-surface-0/80">
            {{ label }} <span v-if="required" class="text-rose-500">*</span>
        </label>
        <div class="flex-1 min-w-0">
            <slot />
            <small v-if="errors" class="text-red-500 text-xs mt-1 block">{{ errors }}</small>
        </div>
    </div>
</template>
