<template>
  <div class="space-y-2">
    <!-- Header: Filter Chips + Search + Actions -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-1.5">
        <Button
          v-for="chip in filterChips"
          :key="chip.value"
          :label="chip.label"
          :severity="moduleTypeFilter === chip.value ? (chip.severity || 'secondary') : 'secondary'"
          :outlined="moduleTypeFilter !== chip.value"
          size="small"
          class="!text-xs !px-2 !py-1"
          @click="setModuleTypeFilter(chip.value)"
        />
      </div>
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" />
        </IconField>
        <Button :label="t('modules.register_module')" icon="pi pi-plus" size="small" @click="openCreate" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable 
    v-else
    :value="filteredModules" 
    paginator 
    :rows="15" 
    size="small" 
    class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-cog text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('modules.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('modules.empty_hint') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('modules.module_name')" sortable>
        <template #body="{ data }">
          <div class="flex-row gap-2">
            <div class="uppercase font-semibold text-gray-600 dark:text-gray-300">{{ data.name }}</div>
            <div class="text-sm text-gray-500 dark:text-gray-400">{{ data.description }}</div>
          </div>
        </template>
      </Column>
      <Column field="slug" :header="t('modules.slug')" sortable />
      <Column field="module_type" :header="t('modules.module_type')" sortable>
        <template #body="{ data }">
          <Tag :value="data.module_type === 'platform' ? t('modules.filter_platform') : t('modules.filter_tenant')" :severity="data.module_type === 'platform' ? 'info' : 'success'" class="!text-xs" />
        </template>
      </Column>
      <Column field="version" :header="t('modules.version')" sortable />
      <Column field="depends_on" :header="t('modules.depends_on')" sortable>
        <template #body="{ data }">
          <span v-if="data.depends_on" class="text-xs text-gray-500" v-tooltip.top="data.depends_on">{{ data.depends_on }}</span>
          <span v-else class="text-gray-300 italic text-xs">—</span>
        </template>
      </Column>
      <Column field="created_at" :header="t('modules.registered')" sortable>
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400">{{ data.created_at || '-' }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" :style="{ width: '120px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('modules.view_companies')" @click="viewCompanies(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEdit(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('modules.edit_module') : t('modules.register_module')" modal :style="{ width: '450px' }">
      <div class="space-y-3">
        <div>
          <FormRow :label="t('modules.module_name')" :errors="errors?.name" :required="true">
              <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
          </FormRow> 
        </div>
        <div>
          <FormRow :label="t('modules.slug')" :errors="errors?.slug" :required="true">
              <div class="relative slug-wrapper" :class="{ 'slug-highlight': slugHighlighted }">
                <TextInput v-model="form.slug" autofocus :class="{ 'p-invalid': errors?.slug }" @input="slugManuallyEdited = true" />
                <i v-if="!slugManuallyEdited && form.name" class="pi pi-sync text-[10px] absolute right-2 top-1/2 -translate-y-1/2 transition-colors duration-300" :class="slugHighlighted ? 'text-indigo-400' : 'text-gray-300'" v-tooltip.left="'Auto-generated from name'"></i>
              </div>
          </FormRow> 
        </div>
        <div>
          <FormRow :label="t('modules.version')" :errors="errors?.version" :required="true">
            <TextInput v-model="form.version" autofocus :class="{ 'p-invalid': errors?.version }" />
          </FormRow> 
        </div>
        <div>
          <FormRow :label="t('modules.desc')" :errors="errors?.description" :required="true">
              <TextInput v-model="form.description" textarea :rows="4" autofocus :class="{ 'p-invalid': errors?.description }" />
          </FormRow> 
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
        <Button :label="isEditing ? t('common.update') : t('modules.register_module')" size="small" :loading="saving" :disabled="saving" @click="saveModule" />
      </template>
    </Dialog>

    <!-- Company Assignment Dialog -->
    <Dialog v-model:visible="companyDialogVisible" :header="'Module: ' + (selectedModule?.name || '')" modal :style="{ width: '500px' }">
      <div class="space-y-2">
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">{{ t('modules.companies_subtitle') }}</p>
        <div v-for="c in moduleCompanies" :key="c.company_id" class="flex items-center justify-between px-3 py-2 rounded-md bg-gray-50 dark:bg-gray-700 text-sm">
          <span>{{ c.company_name || c.company_id }}</span>
          <Tag :value="c.is_active ? 'Active' : 'Inactive'" :severity="c.is_active ? 'success' : 'warn'" class="!text-xs" />
        </div>
        <div v-if="moduleCompanies.length === 0" class="text-sm text-gray-400 text-center py-4">{{ t('modules.no_companies') }}</div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" text size="small" @click="companyDialogVisible = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import { useSkeletonPage } from '@/composables/useSkeletonPage'
import { useSlugify } from '@/composables/useSlugify'

const toast = useToast()
const { t } = useI18n()

const { loading, wrapLoad } = useSkeletonPage()
const modules = ref([])
const searchQuery = ref('')
const moduleTypeFilter = ref(null)
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const loadingCompanies = ref(false)
const form = ref({ name: '', slug: '', version: '1.0.0', description: '' })
const errors = ref({})
const { slugManuallyEdited, slugHighlighted, resetSlug, disableAutoSlug } = useSlugify(
  () => form.value.name,
  (v) => { form.value.slug = v }
)
const companyDialogVisible = ref(false)
const selectedModule = ref(null)
const moduleCompanies = ref([])

const skeletonColumns = [
  { type: 'compound', widths: ['w-24', 'w-36'], headerWidth: 'w-28' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-18' },
  { type: 'text', width: 'w-14', headerWidth: 'w-14' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-14' }
]

// Filter chips
const filterChips = computed(() => [
  { label: t('common.all'), value: null, severity: 'info' },
  { label: t('modules.filter_platform'), value: 'platform', severity: 'info' },
  { label: t('modules.filter_tenant'), value: 'tenant', severity: 'success' }
])

// Filtered modules
const filteredModules = computed(() => {
  let result = modules.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(m =>
      m.name?.toLowerCase().includes(q) ||
      m.slug?.toLowerCase().includes(q) ||
      m.description?.toLowerCase().includes(q)
    )
  }
  return result
})



// Reload modules from API dengan filter opsional
async function loadModules(type) {
  try {
    await wrapLoad(async () => {
      const params = type ? `?module_type=${type}` : ''
      const res = await api.get(`/api/v1/platform/modules${params}`)
      const payload = res.data
      modules.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
    })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
  }
}

function setModuleTypeFilter(type) {
  moduleTypeFilter.value = type
  loadModules(type)
}

onMounted(async () => {
  await loadModules()
})

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', slug: '', version: '1.0.0', description: '' }
  errors.value = {}
  resetSlug()
  dialogVisible.value = true
}

function openEdit(mod) {
  isEditing.value = true
  editingId.value = mod.id
  form.value = { name: mod.name, slug: mod.slug, version: mod.version, description: mod.description || '' }
  errors.value = {}
  disableAutoSlug()
  dialogVisible.value = true
}

async function saveModule() {
  saving.value = true
  try {
    if (isEditing.value) {
      await api.put(`/api/v1/platform/modules/${editingId.value}`, form.value)
      toast.add({ severity: 'success', summary: t('message.updated'), life: 2000 })
    } else {
      await api.post('/api/v1/platform/modules', form.value)
      toast.add({ severity: 'success', summary: t('message.created'), life: 2000 })
    }
    dialogVisible.value = false
    await loadModules(moduleTypeFilter.value)
  } catch (e) {
    errors.value = getValidationErrors(e)
    if (Object.keys(errors.value).length === 0) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
    }
  } finally {
    saving.value = false
  }
}

async function viewCompanies(mod) {
  selectedModule.value = mod
  loadingCompanies.value = true
  try {
    const res = await api.get(`/api/v1/platform/modules/${mod.id}/companies`)
    const payload = res.data
    moduleCompanies.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    moduleCompanies.value = []
  } finally {
    loadingCompanies.value = false
  }
  companyDialogVisible.value = true
}
</script>

