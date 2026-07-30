<script setup>
import { ref, watch, computed } from 'vue'
import PrimeDatePicker from 'primevue/datepicker'
import { useI18n } from '@/composables/useI18n'
import { getPrimeLocale } from '@/utils/primevueLocale'

const props = defineProps({
    modelValue: [String, Date, Number],
    dateFormat: { type: String, default: 'dd/mm/yy' },
    showIcon: { type: Boolean, default: true },
    showTodayButton: { type: Boolean, default: false },
    showClear: { type: Boolean, default: false },
    disabled: { type: Boolean, default: false },
    placeholder: { type: String, default: '' },
    minDate: { type: [Date, String], default: null },
    maxDate: { type: [Date, String], default: null },
    view: { type: String, default: 'date' }
})
const emit = defineEmits(['update:modelValue'])
const dpRef = ref(null)

const { locale } = useI18n()
const primeLocale = computed(() => getPrimeLocale(locale.value))

function toDate(val) {
    if (!val) return null
    if (val instanceof Date) return isNaN(val.getTime()) ? null : val
    if (typeof val === 'string') {
        if (/^\d{4}-\d{2}-\d{2}$/.test(val.trim())) {
            const [y, m, d] = val.split('-').map(Number)
            return new Date(y, m - 1, d)
        }
        const parsed = new Date(val)
        return isNaN(parsed.getTime()) ? null : parsed
    }
    return null
}
function toString(val) {
    if (!val || !(val instanceof Date) || isNaN(val.getTime())) return ''
    const y = val.getFullYear()
    const m = String(val.getMonth() + 1).padStart(2, '0')
    const d = String(val.getDate()).padStart(2, '0')
    return `${y}-${m}-${d}`
}

const parsedMinDate = computed(() => toDate(props.minDate))
const parsedMaxDate = computed(() => toDate(props.maxDate))
const localDate = ref(null)

watch(() => props.modelValue, (nv) => { localDate.value = toDate(nv) }, { immediate: true })

function onDateSelect(selectedVal) {
    const targetDate = selectedVal instanceof Date ? selectedVal : localDate.value
    localDate.value = targetDate
    emit('update:modelValue', toString(targetDate))
    if (dpRef.value) {
        if (typeof dpRef.value.hide === 'function') dpRef.value.hide()
        else if ('overlayVisible' in dpRef.value) dpRef.value.overlayVisible = false
    }
}
function onClearClick() { localDate.value = null; emit('update:modelValue', '') }
</script>
<template>
    <PrimeDatePicker ref="dpRef" v-model="localDate" size="small" :dateFormat="dateFormat" :placeholder="placeholder" :showIcon="showIcon" :showTodayButton="showTodayButton" :showClear="showClear" :disabled="disabled" :minDate="parsedMinDate" :maxDate="parsedMaxDate" :view="view" :locale="primeLocale" class="w-full" @date-select="onDateSelect" @clear-click="onClearClick" />
</template>
