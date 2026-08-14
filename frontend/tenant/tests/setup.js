import { config } from '@vue/test-utils'
import { h, cloneVNode } from 'vue'
import { vi } from 'vitest'

// Force Indonesian locale untuk assertion string (locale store membaca
// localStorage saat modul di-load; setupFiles berjalan sebelum test file).
localStorage.setItem('hris_lang', 'id')

// jsdom tidak punya navigator.clipboard — mock manual
Object.defineProperty(navigator, 'clipboard', {
  value: {
    writeText: vi.fn().mockResolvedValue(undefined),
    readText: vi.fn().mockResolvedValue(''),
  },
  configurable: true,
})

if (!Element.prototype.scrollIntoView) {
  Element.prototype.scrollIntoView = vi.fn()
}

if (!window.matchMedia) {
  window.matchMedia = vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  }))
}

// ── Stub DataTable/Column ──
// DataTable asli PrimeVue mengiterasi `value` dan merender body slot tiap
// Column dengan `{ data: row }`. Stub ini mereplikasi perilaku tersebut:
// untuk tiap row, Column di-clone dengan prop `data` = row.
const DataTable = {
  name: 'DataTable',
  props: { value: { type: Array, default: () => [] } },
  setup(props, { slots }) {
    return () => {
      const rows = props.value || []
      if (!rows.length) {
        return h('div', { class: 'dt-stub dt-empty' }, slots.empty ? slots.empty() : [])
      }
      const cols = (slots.default ? slots.default() : []).filter(Boolean)
      return h(
        'div',
        { class: 'dt-stub' },
        rows.map((row) =>
          h(
            'div',
            { class: 'dt-row' },
            cols.map((col) => cloneVNode(col, { data: row }))
          )
        )
      )
    }
  },
}

const Column = {
  name: 'Column',
  props: { data: { type: Object, default: () => ({}) } },
  template: '<div class="col-stub"><slot name="body" :data="data" /><slot /></div>',
}

// ── Stub komponen PrimeVue lain ──
const Button = {
  props: ['label', 'icon', 'loading', 'disabled'],
  emits: ['click'],
  template:
    '<button class="btn-stub" :disabled="disabled || loading" @click="$emit(\'click\')"><i v-if="icon" :class="icon"></i>{{ label }}<slot /></button>',
}

const Tag = {
  props: ['value', 'severity'],
  template: '<span class="tag-stub" :data-severity="severity">{{ value }}</span>',
}

const Select = {
  props: ['modelValue', 'options', 'optionLabel', 'optionValue', 'placeholder', 'showClear'],
  emits: ['update:modelValue', 'change'],
  template:
    '<select class="select-stub" :value="modelValue" @change="$emit(\'update:modelValue\', $event.target.value); $emit(\'change\', $event)"><option v-if="placeholder" value="">{{ placeholder }}</option><option v-for="opt in options" :key="opt[optionValue]" :value="opt[optionValue]">{{ opt[optionLabel] }}</option><slot /></select>',
}

// PrimeVue Dialog & ConfirmDeleteDialog memakai `v-model:visible` → prop-nya `visible`
const Dialog = {
  props: ['visible', 'header', 'closable', 'contentStyle'],
  emits: ['update:visible', 'hide'],
  template:
    '<div class="dialog-stub" v-if="visible"><header class="dialog-header">{{ header }}</header><div class="dialog-content" :style="contentStyle"><slot /><slot name="footer" /></div></div>',
}

const ConfirmDeleteDialog = {
  props: ['visible', 'title', 'message', 'loading', 'confirmLabel', 'errorMsg'],
  emits: ['update:visible', 'confirm', 'cancel'],
  template:
    '<div class="confirm-stub" v-if="visible"><p class="confirm-title">{{ title }}</p><p class="confirm-message">{{ message }}</p><slot /><button class="confirm-ok" @click="$emit(\'confirm\')">{{ confirmLabel || \'OK\' }}</button><button class="confirm-cancel" @click="$emit(\'cancel\')">Batal</button></div>',
}

const SkeletonTable = { template: '<div class="skeleton-stub" />' }

const TextInput = {
  props: ['modelValue'],
  emits: ['update:modelValue'],
  template:
    '<input class="textinput-stub" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
}

const FormRow = {
  props: ['label', 'required', 'errors'],
  template:
    '<div class="formrow-stub"><label>{{ label }}<span v-if="required">*</span></label><slot /><small v-if="errors" class="err-text">{{ Array.isArray(errors) ? errors[0] : errors }}</small></div>',
}

const Textarea = {
  props: ['modelValue'],
  emits: ['update:modelValue'],
  template:
    '<textarea class="textarea-stub" :value="modelValue" @input="$emit(\'update:modelValue\', $event.target.value)" />',
}

config.global.stubs = {
  ...config.global.stubs,
  DataTable,
  Column,
  Button,
  Tag,
  Select,
  Dialog,
  ConfirmDeleteDialog,
  SkeletonTable,
  TextInput,
  FormRow,
  Textarea,
}

// Directive v-tooltip dipakai di tombol action — stub kosong
config.global.directives = {
  ...config.global.directives,
  tooltip: {},
}
