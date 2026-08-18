<script setup>
import { ref, watch } from 'vue'

/**
 * TimeInput — input jam 3 segmen (HH : MM : SS).
 * v-model: string "HH:MM:SS" (kosong bila belum lengkap).
 * Fitur: hanya digit, auto-advance antar segmen, navigasi panah/backspace,
 * clamp nilai (HH 00-23, MM/SS 00-59), pad otomatis saat blur.
 */
const props = defineProps({
    modelValue: { type: String, default: '' },
    disabled: { type: Boolean, default: false },
    placeholder: { type: String, default: '--' }
})
const emit = defineEmits(['update:modelValue'])

const LIMITS = [23, 59, 59]
const parts = ref(['', '', ''])
const inputs = [null, null, null]

function pad(v) { return String(v).padStart(2, '0') }

function splitModel(v) {
    const m = String(v || '').trim().match(/^(\d{1,2}):(\d{1,2}):(\d{1,2})$/)
    return m ? [m[1], m[2], m[3]] : ['', '', '']
}

watch(() => props.modelValue, (v) => {
    const [h, m, s] = splitModel(v)
    parts.value = [h, m, s]
}, { immediate: true })

function setInputRef(idx, el) { inputs[idx] = el }

function emitChange() {
    const [h, m, s] = parts.value
    const complete = h !== '' && m !== '' && s !== ''
    emit('update:modelValue', complete ? `${pad(h)}:${pad(m)}:${pad(s)}` : '')
}

function onInput(idx, e) {
    const digits = e.target.value.replace(/\D/g, '').slice(0, 2)
    parts.value[idx] = digits
    if (digits.length === 2 && idx < 2) inputs[idx + 1]?.focus()
    emitChange()
}

function onKeydown(idx, e) {
    if (e.key === 'Backspace' && parts.value[idx] === '' && idx > 0) {
        e.preventDefault()
        inputs[idx - 1]?.focus()
    } else if (e.key === 'ArrowLeft' && idx > 0) {
        e.preventDefault()
        inputs[idx - 1]?.focus()
    } else if (e.key === 'ArrowRight' && idx < 2) {
        e.preventDefault()
        inputs[idx + 1]?.focus()
    } else if (e.key === ':' && idx < 2) {
        e.preventDefault()
        inputs[idx + 1]?.focus()
    }
}

function onBlur(idx) {
    const v = parts.value[idx]
    if (v === '') return
    let n = parseInt(v, 10)
    if (isNaN(n)) { parts.value[idx] = ''; emitChange(); return }
    n = Math.min(n, LIMITS[idx])
    parts.value[idx] = pad(n)
    emitChange()
}
</script>
<template>
    <div class="flex items-center gap-1 w-full" :class="{ 'opacity-60 pointer-events-none': disabled }">
        <template v-for="(part, idx) in parts" :key="idx">
            <span v-if="idx > 0" class="text-gray-400 dark:text-gray-500 font-medium select-none">:</span>
            <input
                :ref="(el) => setInputRef(idx, el)"
                :value="part"
                :disabled="disabled"
                :placeholder="placeholder"
                inputmode="numeric"
                maxlength="2"
                class="w-11 px-1 py-1.5 text-center text-sm text-navy-800 dark:text-gray-100 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md focus:outline-none focus:ring-1 focus:ring-emerald-500 focus:border-emerald-500 placeholder:text-gray-400 dark:placeholder:text-gray-500 disabled:cursor-not-allowed"
                @input="onInput(idx, $event)"
                @keydown="onKeydown(idx, $event)"
                @blur="onBlur(idx)"
                @focus="(e) => e.target.select()"
            />
        </template>
    </div>
</template>
