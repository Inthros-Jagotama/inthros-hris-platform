<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse w-1/3"></div>
      <div class="grid grid-cols-2 gap-4">
        <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        <div class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
      <div class="h-40 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
    </div>

    <template v-else>
      <!-- Template Info Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 mb-4 flex items-center gap-2">
          <i class="pi pi-file-edit text-teal-500"></i>
          {{ t('document_templates.template_info') }}
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormRow :label="t('document_templates.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" autofocus :placeholder="t('document_templates.name')" :class="{ 'p-invalid': errors?.name }" />
          </FormRow>
          <FormRow :label="t('document_templates.document_type')" required :errors="errors?.document_type">
            <Select
              v-if="!isEditing"
              v-model="form.document_type"
              :options="documentTypeOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('document_templates.filter_document_type')"
              class="w-full"
              :class="{ 'p-invalid': errors?.document_type }"
            />
            <Select v-else :modelValue="form.document_type" :options="documentTypeOptions" optionLabel="label" optionValue="value" disabled class="w-full opacity-60" />
          </FormRow>
          <FormRow v-if="form.document_type === 'MOVEMENT_SK'" :label="t('document_templates.movement_type')" :errors="errors?.movement_type">
            <Select
              v-model="form.movement_type"
              :options="movementTypeOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('document_templates.movement_type_generic')"
              :disabled="isEditing"
              :class="{ 'opacity-60': isEditing }"
              class="w-full"
            />
          </FormRow>
          <FormRow :label="t('document_templates.description_label')" :errors="errors?.description">
            <Textarea v-model="form.description" rows="2" :placeholder="t('document_templates.description_label')" class="w-full" />
          </FormRow>
        </div>
      </div>

      <!-- Template File + Document Configuration (2 kolom) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Template File Card (DOCX) -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 mb-4 flex items-center gap-2">
          <i class="pi pi-file-word text-emerald-500"></i>
          {{ t('document_templates.template_file') }}
        </h3>

        <div v-if="isEditing && currentFile" class="flex items-center justify-between gap-3 rounded-lg border border-gray-200 dark:border-gray-700 px-4 py-3 mb-3">
          <div class="flex items-center gap-3 min-w-0">
            <i class="pi pi-file-word text-xl text-emerald-500"></i>
            <div class="min-w-0">
              <p class="text-sm font-medium text-navy-800 dark:text-gray-100 truncate">{{ currentFile.name }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.current_template_file') }}</p>
            </div>
          </div>
          <a v-if="currentFile.url" :href="currentFile.url" :download="currentFile.name" class="text-sm text-teal-500 hover:text-teal-600 dark:text-teal-400 font-medium whitespace-nowrap">
            <i class="pi pi-download mr-1"></i>{{ t('common.download') }}
          </a>
        </div>

        <div class="flex items-center gap-3 flex-wrap">
          <label
            class="flex items-center gap-2 cursor-pointer text-sm px-4 py-2 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 hover:border-teal-400 hover:bg-teal-50/50 dark:hover:bg-teal-500/5 transition-colors"
          >
            <i class="pi pi-upload text-teal-500"></i>
            <span class="text-gray-600 dark:text-gray-300">
              {{ isEditing && currentFile ? t('document_templates.replace_file') : t('document_templates.choose_file') }}
            </span>
            <input type="file" accept=".docx,application/vnd.openxmlformats-officedocument.wordprocessingml.document" class="hidden" @change="onFileChange" />
          </label>
          <span v-if="templateFile" class="text-sm text-gray-600 dark:text-gray-300 font-medium">
            {{ templateFile.name }}
            <span class="text-xs text-gray-400">({{ formatFileSize(templateFile.size) }})</span>
          </span>
        </div>
        <p v-if="errors?.file" class="text-xs text-red-500 mt-2">{{ errors.file[0] }}</p>
        <div class="mt-3 text-xs text-gray-400 dark:text-gray-500 space-y-1">
          <p>{{ t('document_templates.docx_hint') }}</p>
          <p>{{ t('document_templates.docx_layout_hint') }}</p>
        </div>
      </div>

      <!-- Document Configuration Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 mb-1 flex items-center gap-2">
          <i class="pi pi-sliders-h text-sky-500"></i>
          {{ t('document_templates.document_configuration') }}
        </h3>
        <p class="text-xs text-gray-400 dark:text-gray-500 mb-4">{{ t('document_templates.document_config_hint') }}</p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormRow :label="t('document_templates.paper_size')">
            <Select v-model="form.paper_size" :options="paperOptions" optionLabel="label" optionValue="value" class="w-full" />
          </FormRow>
          <FormRow :label="t('document_templates.orientation')">
            <Select v-model="form.orientation" :options="orientationOptions" optionLabel="label" optionValue="value" class="w-full" />
          </FormRow>
        </div>
      </div>
      </div> <!-- /grid 2 kolom -->

      <!-- Actions -->
      <div class="flex items-center justify-between">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="goBack" />
        <div class="flex items-center gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="goBack" />
          <Button :label="isEditing ? t('common.update') : t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </div>

      <!-- Variable Reference Card -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
        <h3 class="text-sm font-semibold text-navy-800 dark:text-gray-100 mb-1 flex items-center gap-2">
          <i class="pi pi-database text-violet-500"></i>
          {{ t('document_templates.variable_reference') }}
        </h3>
        <p class="text-xs text-gray-400 dark:text-gray-500 mb-4">{{ t('document_templates.variable_reference_hint') }}</p>
        <div class="mb-4 flex items-center gap-2">
          <IconField class="w-full max-w-xs">
            <InputIcon class="pi pi-search" />
            <InputText v-model="variableSearch" :placeholder="t('document_templates.variable_search_placeholder')" size="small" class="!pl-8 !text-sm !py-1.5" />
          </IconField>
          <Button v-if="variableSearch" icon="pi pi-times" severity="secondary" text rounded size="small" class="!p-1 shrink-0" @click="variableSearch = ''" />
        </div>
        <div v-if="variableGroups.length === 0" class="text-xs text-gray-400 dark:text-gray-500">
          <i class="pi pi-spin pi-spinner mr-1"></i>{{ t('common.loading') }}
        </div>
        <div v-else-if="variableColumns.length === 0" class="text-xs text-gray-400 dark:text-gray-500">
          <i class="pi pi-search mr-1"></i>{{ t('document_templates.variable_search_no_results') }}
        </div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6 items-start">
          <!-- Kolom kiri: employee; kolom kanan: contract, movement, company -->
          <div v-for="(colGroups, colIdx) in variableColumns" :key="colIdx" :class="{ 'space-y-6': colIdx > 0 }">
            <div v-for="group in colGroups" :key="group.category">
              <p class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400 mb-2">{{ variableCategoryLabel(group.category) }}</p>
              <div class="space-y-1.5">
                <button
                  v-for="v in group.variables"
                  :key="v.key"
                  type="button"
                  class="w-full text-left flex items-center justify-between gap-2 px-3 py-2 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-teal-400 hover:bg-teal-50/50 dark:hover:bg-teal-500/5 transition-colors"
                  @click="copyVariable(v)"
                >
                  <span class="text-sm text-gray-700 dark:text-gray-200">{{ variableLabel(v) }}</span>
                  <code class="text-[11px] text-teal-500 dark:text-teal-400 font-mono">{{ placeholderText(v.key) }}</code>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const route = useRoute()
const toast = useToast()
const { t } = useI18n()

const pageLoading = ref(true)
const saving = ref(false)
const errors = ref({})

const templateId = computed(() => route.params.id)
const isEditing = computed(() => !!templateId.value && templateId.value !== 'new')

// File template .docx
const templateFile = ref(null)
// File tersimpan pada versi aktif (mode edit) — { name, url }
const currentFile = ref(null)

const form = ref({
  name: '',
  document_type: '',
  movement_type: '',
  description: '',
  paper_size: 'A4',
  orientation: 'portrait',
})

const documentTypeOptions = [
  { label: t('document_templates.type_contract_agreement'), value: 'CONTRACT_AGREEMENT' },
  { label: t('document_templates.type_movement_sk'), value: 'MOVEMENT_SK' },
]

// Daftar jenis movement dari backend (GET /movement-types) + opsi umum.
const movementTypeOptions = ref([])

async function loadMovementTypes() {
  try {
    const res = await api.get('/api/v1/tenant/settings/document-templates/movement-types')
    const list = res.data?.data || []
    movementTypeOptions.value = [
      { label: t('document_templates.movement_type_generic'), value: '' },
      ...list.map((m) => {
        const key = `employee_movement.type_${m.value}`
        const label = t(key)
        return { label: label !== key ? label : m.label, value: m.value }
      }),
    ]
  } catch (e) {
    movementTypeOptions.value = [{ label: t('document_templates.movement_type_generic'), value: '' }]
  }
}

const paperOptions = [
  { label: 'A4', value: 'A4' },
  { label: 'A5', value: 'A5' },
  { label: 'Letter', value: 'Letter' },
  { label: 'Legal', value: 'Legal' },
]

const orientationOptions = [
  { label: 'Portrait', value: 'portrait' },
  { label: 'Landscape', value: 'landscape' },
]

function goBack() {
  router.push('/settings/document-templates')
}

function formatFileSize(bytes) {
  if (!bytes) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

// ── File upload ──
function onFileChange(event) {
  const file = event.target.files?.[0] || null
  event.target.value = ''
  errors.value = { ...errors.value, file: undefined }
  if (!file) {
    templateFile.value = null
    return
  }
  const name = file.name.toLowerCase()
  if (!name.endsWith('.docx')) {
    templateFile.value = null
    errors.value = { ...errors.value, file: [t('document_templates.file_invalid_type')] }
    return
  }
  if (file.size > 10 * 1024 * 1024) {
    templateFile.value = null
    errors.value = { ...errors.value, file: [t('document_templates.file_too_large')] }
    return
  }
  templateFile.value = file
}

// ── Variable reference ──
const variableGroups = ref([])
const variableSearch = ref('')

// variableLabel menampilkan label variabel sesuai bahasa aktif (bilingual),
// dengan fallback ke label bawaan registry backend bila belum ada terjemahan.
function variableLabel(v) {
  const key = `document_templates.var_label.${v.key}`
  const label = t(key)
  return label !== key ? label : v.label
}

// Grup variabel difilter pencarian (cocok pada label terjemahan, label
// bawaan, ATAU key — case-insensitive) lalu dibagi ke 2 kolom: employee di
// kolom kiri, sisanya (contract, movement, company) ditumpuk di kolom kanan —
// urutan registry dipertahankan. Grup tanpa variabel yang cocok dihilangkan;
// hasil kosong ditampilkan sebagai "tidak ada yang cocok" di template.
const variableColumns = computed(() => {
  const q = variableSearch.value.trim().toLowerCase()
  const groups = variableGroups.value
    .map((g) => ({
      ...g,
      variables: q
        ? g.variables.filter((v) => {
            const translated = (variableLabel(v) || '').toLowerCase()
            const raw = (v.label || '').toLowerCase()
            const key = (v.key || '').toLowerCase()
            return translated.includes(q) || raw.includes(q) || key.includes(q)
          })
        : g.variables,
    }))
    .filter((g) => g.variables.length > 0)
  const primary = groups.filter((g) => g.category === 'employee')
  const rest = groups.filter((g) => g.category !== 'employee')
  return [primary, rest].filter((col) => col.length > 0)
})

function variableCategoryLabel(category) {
  const key = `document_templates.var_${category}`
  const label = t(key)
  return label !== key ? label : category
}

function placeholderText(key) {
  return `{{${key}}}`
}

async function copyVariable(v) {
  const token = `{{${v.key}}}`
  try {
    await navigator.clipboard.writeText(token)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.variable_copied', { variable: token }), life: 2000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.operation_failed'), life: 3000 })
  }
}

async function loadVariables() {
  try {
    const res = await api.get('/api/v1/tenant/settings/document-templates/variables')
    variableGroups.value = res.data?.data || []
  } catch (e) {
    variableGroups.value = []
  }
}

// ── Load ──
async function loadTemplate() {
  if (!isEditing.value) return

  try {
    const res = await api.get(`/api/v1/tenant/settings/document-templates/${templateId.value}`)
    const data = res.data?.data || res.data
    form.value.name = data.name || ''
    form.value.document_type = data.document_type || ''
    form.value.movement_type = data.movement_type || ''
    // Pastikan nilai movement_type yang tersimpan tetap muncul di dropdown
    // meski fetch movement-types gagal / daftarnya berubah.
    if (data.movement_type && !movementTypeOptions.value.some((o) => o.value === data.movement_type)) {
      const key = `employee_movement.type_${data.movement_type}`
      const label = t(key)
      movementTypeOptions.value.push({ label: label !== key ? label : data.movement_type, value: data.movement_type })
    }
    form.value.description = data.description || ''

    // Muat info file + konfigurasi dari versi aktif (atau fallback versi pertama).
    let config = {}
    try {
      const vRes = await api.get(`/api/v1/tenant/settings/document-templates/${templateId.value}/versions`)
      const versions = vRes.data?.data || []
      const active = versions.find((v) => v.id === data.active_version_id) || versions[0]
      if (active) {
        config = active
        if (active.file_url) {
          currentFile.value = { name: active.file_name || 'template.docx', url: active.file_url }
        }
      }
    } catch (e) {
      // versions mungkin kosong — pakai fallback di atas
    }
    form.value.paper_size = config.paper_size || 'A4'
    form.value.orientation = config.orientation || 'portrait'
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  }
}

// ── Save ──
function buildVersionFormData() {
  const fd = new FormData()
  fd.append('file', templateFile.value)
  fd.append('paper_size', form.value.paper_size)
  fd.append('orientation', form.value.orientation)
  return fd
}

async function handleSave() {
  errors.value = {}

  if (!form.value.name?.trim()) { errors.value = { ...errors.value, name: [t('form.required')] }; return }
  if (!isEditing.value && !form.value.document_type) { errors.value = { ...errors.value, document_type: [t('form.required')] }; return }
  if (!templateFile.value) {
    errors.value = { ...errors.value, file: [t('document_templates.file_required')] }
    return
  }

  saving.value = true
  try {
    let targetId = templateId.value

    if (isEditing.value) {
      await api.put(`/api/v1/tenant/settings/document-templates/${templateId.value}`, {
        name: form.value.name.trim(),
        description: form.value.description || undefined,
      })
    } else {
      const res = await api.post('/api/v1/tenant/settings/document-templates', {
        name: form.value.name.trim(),
        document_type: form.value.document_type,
        movement_type: form.value.movement_type || undefined,
        description: form.value.description || undefined,
      })
      targetId = res.data?.data?.id || res.data?.id
    }

    // Simpan file .docx sebagai versi baru (v1 untuk template baru, versi
    // berikutnya saat file diganti di edit mode). Backend menyimpan file ke
    // /uploads/document_templates/ dan content = path file, lalu memvalidasi
    // placeholder (variable validation) — unknown variable akan ditolak 400.
    const versionRes = await api.post(`/api/v1/tenant/settings/document-templates/${targetId}/versions`, buildVersionFormData(), {
      headers: { 'Content-Type': 'multipart/form-data' },
    })

    const placeholders = versionRes.data?.data?.placeholders
    if (Array.isArray(placeholders) && placeholders.length > 0) {
      toast.add({
        severity: 'info',
        summary: t('document_templates.variables_found_title'),
        detail: t('document_templates.variables_found', { count: String(placeholders.length) }),
        life: 4000,
      })
    }
    toast.add({
      severity: 'success',
      summary: t('message.success'),
      detail: isEditing.value ? t('document_templates.updated') : t('document_templates.created'),
      life: 3000,
    })
    goBack()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({
        severity: 'error',
        summary: t('message.error'),
        detail: e.response?.data?.error?.message || t('message.operation_failed'),
        life: 4000,
      })
    }
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    await Promise.all([loadVariables(), loadMovementTypes(), loadTemplate()])
  } finally {
    pageLoading.value = false
  }
})
</script>
