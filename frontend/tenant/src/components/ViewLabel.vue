<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
    label: { type: String, default: '' },
    value: { type: [String, Number], default: '' },
    mono: { type: Boolean, default: false },
    breakAll: { type: Boolean, default: false },
    copyable: { type: Boolean, default: false },
    // Bilingual: halaman boleh mengirim t() untuk judul tombol copy.
    copyTitle: { type: String, default: 'Copy' },
    copiedTitle: { type: String, default: 'Copied' }
})

// Nil kosong ('' / null / undefined) ditampilkan sebagai '—' (bukan 0).
const hasValue = computed(() => props.value !== '' && props.value !== null && props.value !== undefined)

const copied = ref(false)
async function copyValue() {
    if (!hasValue.value) return
    try {
        await navigator.clipboard.writeText(String(props.value))
        copied.value = true
        setTimeout(() => { copied.value = false }, 1500)
    } catch {
        // Clipboard tidak tersedia — abaikan.
    }
}
</script>

<template>
    <div class="flex flex-col gap-1">
        <!-- Label: font lebih kecil dari value (text-xs < text-sm) -->
        <label class="text-xs text-gray-400 dark:text-gray-500 uppercase tracking-wide">
            {{ label }}
        </label>

        <!-- Slot untuk konten custom (Tag, link, dll) -->
        <span v-if="$slots.default" class="text-sm text-gray-700 dark:text-gray-200">
            <slot />
        </span>

        <!-- Render value standar + tombol copy opsional -->
        <span
            v-else
            class="inline-flex items-center gap-1.5 text-sm text-gray-700 dark:text-gray-200"
            :class="{ 'font-mono': mono, 'break-all': breakAll }"
        >
            <template v-if="hasValue">{{ value }}</template>
            <template v-else>—</template>
            <button
                v-if="copyable && hasValue"
                type="button"
                class="text-gray-400 hover:text-emerald-500 transition-colors"
                :title="copied ? props.copiedTitle : props.copyTitle"
                :aria-label="copied ? props.copiedTitle : props.copyTitle"
                @click="copyValue"
            >
                <i :class="copied ? 'pi pi-check' : 'pi pi-copy'" style="font-size: 0.75rem"></i>
            </button>
        </span>
    </div>
</template>
