<script setup>
import { ref, watch, computed } from 'vue'
import PrimeDatePicker from 'primevue/datepicker' // Menggunakan alias agar tidak bentrok nama

const props = defineProps({
    modelValue: [String, Date, Number],
    dateFormat: { type: String, default: 'dd/mm/yy' }, // 'yy' di PrimeVue = 4 digit tahun (contoh: 2026)
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

// 1. Safe parsing string "YYYY-MM-DD" menjadi Objek Date lokal
function toDate(val) {
    if (!val) return null
    if (val instanceof Date) return isNaN(val.getTime()) ? null : val
    
    if (typeof val === 'string') {
        const cleanStr = val.trim()
        if (/^\d{4}-\d{2}-\d{2}$/.test(cleanStr)) {
            const [y, m, d] = cleanStr.split('-').map(Number)
            return new Date(y, m - 1, d)
        }
        const parsed = new Date(cleanStr)
        return isNaN(parsed.getTime()) ? null : parsed
    }
    return null
}

// 2. Convert Objek Date menjadi string "YYYY-MM-DD"
function toString(val) {
    if (!val || !(val instanceof Date) || isNaN(val.getTime())) return ''
    const y = val.getFullYear()
    const m = String(val.getMonth() + 1).padStart(2, '0')
    const d = String(val.getDate()).padStart(2, '0')
    return `${y}-${m}-${d}`
}

// Ensure min/max date always valid Date objects for PrimeVue
const parsedMinDate = computed(() => toDate(props.minDate))
const parsedMaxDate = computed(() => toDate(props.maxDate))

// Local Date state
const localDate = ref(null)

// Sync Props -> Local State (Hanya ketika parent mengubah nilai dari luar)
watch(
    () => props.modelValue,
    (nv) => {
        localDate.value = toDate(nv)
    },
    { immediate: true }
)

// Dijalankan HANYA saat pengguna mengklik/memilih tanggal di kalender
function onDateSelect(selectedVal) {
    // 1. Ambil nilai Date dari parameter @date-select
    const targetDate = selectedVal instanceof Date ? selectedVal : localDate.value
    localDate.value = targetDate
    // 2. Format menjadi string dan Emit ke Parent
    const strValue = toString(targetDate)
    emit('update:modelValue', strValue)

    // 3. Tutup Popup Kalender
    if (dpRef.value) {
        if (typeof dpRef.value.hidePicker === 'function') {
            dpRef.value.hidePicker()
        } else if (typeof dpRef.value.hide === 'function') {
            dpRef.value.hide()
        } else if ('overlayVisible' in dpRef.value) {
            dpRef.value.overlayVisible = false
        }
    }
}

// Opsional: Handler saat tombol Clear diklik
function onClearClick() {
    localDate.value = null
    emit('update:modelValue', '')
}

const date = ref();
</script>

<template>
    <PrimeDatePicker 
        ref="dpRef"
        v-model="localDate" 
        size="small"
        :dateFormat="dateFormat"
        :placeholder="placeholder"
        :showIcon="showIcon"
        :showTodayButton="showTodayButton"
        :showClear="showClear"
        :todayButton="todayButton"
        :clearButton="clearButton"
        :disabled="disabled"
        :minDate="parsedMinDate"
        :maxDate="parsedMaxDate"
        :view="view"
        class="w-full"
        @date-select="onDateSelect"
        @clear-click="onClearClick"
        v-mask="'99/99/9999'" />
</template>