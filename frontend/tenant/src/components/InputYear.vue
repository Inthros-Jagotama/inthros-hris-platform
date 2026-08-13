<script setup>
import { ref, watch, computed } from 'vue'
import PrimeDatePicker from 'primevue/datepicker'
import { useI18n } from '@/composables/useI18n'
import { getPrimeLocale } from '@/utils/primevueLocale'

const props = defineProps({
    modelValue: [Number, String],
    placeholder: { type: String, default: '' },
    showClear: { type: Boolean, default: false },
    disabled: { type: Boolean, default: false },
    minYear: { type: Number, default: 2000 },
    maxYear: { type: Number, default: 2100 }
})
const emit = defineEmits(['update:modelValue'])
const dpRef = ref(null)

const { locale } = useI18n()
const primeLocale = computed(() => getPrimeLocale(locale.value))

function toYear(val) {
    if (val === null || val === undefined || val === '') return null
    if (val instanceof Date) return isNaN(val.getTime()) ? null : val.getFullYear()
    const y = Number(val)
    return Number.isFinite(y) ? y : null
}

const localDate = ref(null)
const minDate = computed(() => (props.minYear ? new Date(props.minYear, 0, 1) : null))
const maxDate = computed(() => (props.maxYear ? new Date(props.maxYear, 11, 31) : null))

watch(() => props.modelValue, (nv) => {
    const y = toYear(nv)
    localDate.value = y ? new Date(y, 0, 1) : null
}, { immediate: true })

function onYearSelect(selectedVal) {
    const targetDate = selectedVal instanceof Date ? selectedVal : localDate.value
    localDate.value = targetDate
    emit('update:modelValue', targetDate instanceof Date ? targetDate.getFullYear() : null)
    if (dpRef.value) {
        if (typeof dpRef.value.hide === 'function') dpRef.value.hide()
        else if ('overlayVisible' in dpRef.value) dpRef.value.overlayVisible = false
    }
}
function onClearClick() {
    localDate.value = null
    emit('update:modelValue', null)
}
</script>
<template>
    <PrimeDatePicker ref="dpRef" v-model="localDate" size="small" view="year" dateFormat="yy" :placeholder="placeholder" showIcon :showClear="showClear" :disabled="disabled" :minDate="minDate" :maxDate="maxDate" :locale="primeLocale" class="w-full" @date-select="onYearSelect" @clear-click="onClearClick" />
</template>
