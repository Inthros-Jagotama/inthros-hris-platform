<script setup>
import { ref, watch } from 'vue'
import DatePicker from 'primevue/datepicker'

const props = defineProps({
    modelValue: [String, Date, Number],
    dateFormat: { type: String, default: 'yy-mm-dd' },
    showIcon: { type: Boolean, default: true },
    showTodayButton: { type: Boolean, default: false },
    showClear: { type: Boolean, default: false },
    todayButton: { type: String, default: 'Today' },
    clearButton: { type: String, default: 'Clear' },
    disabled: { type: Boolean, default: false },
    placeholder: { type: String, default: '' },
    minDate: { type: [Date, String], default: null },
    maxDate: { type: [Date, String], default: null },
    view: { type: String, default: 'date' },
    selectionMode: { type: String, default: 'date' }
})

const emit = defineEmits(['update:modelValue'])

const dpRef = ref(null)

// Convert string 'yyyy-mm-dd' → Date
function toDate(val) {
    if (!val) return null
    if (val instanceof Date) return val
    if (typeof val === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(val)) {
        const [y, m, d] = val.split('-').map(Number)
        return new Date(y, m - 1, d)
    }
    return null
}

// Convert Date → string 'yyyy-mm-dd'
function toString(val) {
    if (!val) return ''
    if (val instanceof Date && !isNaN(val)) {
        const y = val.getFullYear()
        const m = String(val.getMonth() + 1).padStart(2, '0')
        const d = String(val.getDate()).padStart(2, '0')
        return `${y}-${m}-${d}`
    }
    return ''
}

// Local ref — stable reference for PrimeVue's v-model
const localDate = ref(null)

// Sync prop change → localDate (parent writes)
watch(() => props.modelValue, (nv) => {
    localDate.value = toDate(nv)
}, { immediate: true })

// Handle user date selection directly from PrimeVue event
function onSelect(value) {
    const str = toString(value)
    if (str !== props.modelValue) {
        emit('update:modelValue', str)
    }
}

defineExpose({
    focus: () => dpRef.value?.$el?.querySelector('input')?.focus()
})
</script>

<template>
    <DatePicker
        ref="dpRef"
        v-model="localDate"
        :dateFormat="dateFormat"
        :showIcon="showIcon"
        :showTodayButton="showTodayButton"
        :showClear="showClear"
        :todayButton="todayButton"
        :clearButton="clearButton"
        :disabled="disabled"
        :placeholder="placeholder"
        :minDate="minDate"
        :maxDate="maxDate"
        :view="view"
        :selectionMode="selectionMode"
        class="w-full"
        @update:modelValue="onSelect"
    />
</template>
