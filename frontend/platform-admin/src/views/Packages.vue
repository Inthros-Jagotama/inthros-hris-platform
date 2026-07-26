<template>
  <div class="space-y-2">
    <!-- Header: Search + Actions -->
    <div class="flex items-center justify-end gap-2 flex-wrap">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" />
        </IconField>
        <Button :label="t('packages.new_package')" icon="pi pi-plus" size="small" @click="openCreate" />
      </div>

    <DataTable 
      :value="filteredPackages" 
      paginator 
      :rows="15" 
      size="small" 
      sortField="sort_order" 
      :sortOrder="1"
      :loading="loading"
      class="!text-sm p-datatable-sm border border-gray-200 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400">
          <i class="pi pi-box text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('packages.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('packages.empty_hint') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('packages.package_name')" sortable>
        <template #body="{ data }">
          <div class="flex-row gap-2">
            <div class="uppercase font-semibold text-gray-600">{{ data.name }}</div>
            <div class="text-xs text-gray-400">{{ data.slug }}</div>
          </div>
        </template>
      </Column>
      <Column field="price" :header="t('packages.price')" sortable>
        <template #body="{ data }">
          <span class="font-medium text-gray-700">{{ formatPrice(data.price) }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('packages.status')" sortable>
        <template #body="{ data }">
          <Tag :value="data.status === 'published' ? t('common_status.published') : t('common_status.draft')" :severity="data.status === 'published' ? 'success' : 'info'" class="!text-xs" />
        </template>
      </Column>
      <Column field="module_count" :header="t('packages.module_count')" sortable>
        <template #body="{ data }">
          <span class="text-gray-600">{{ data.module_count }}</span>
        </template>
      </Column>
      <Column field="sort_order" :header="t('packages.sort_order')" sortable />
      <Column :header="t('common.actions')" :style="{ width: '200px' }">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="{ value: t('packages.tooltip_edit'), showDelay: 300 }" @click="openEdit(data)" />
            <Button v-if="data.status === 'draft'" icon="pi pi-send" size="small" text severity="success" v-tooltip.left="{ value: t('packages.tooltip_publish'), showDelay: 300 }" @click="confirmPublish(data)" />
            <Button v-if="data.status === 'published'" icon="pi pi-undo" size="small" text severity="warn" v-tooltip.left="{ value: t('packages.tooltip_unpublish'), showDelay: 300 }" @click="confirmUnpublish(data)" />
            <Button icon="pi pi-check-circle" size="small" text severity="info" v-tooltip.left="{ value: t('packages.tooltip_validate'), showDelay: 300 }" @click="validateDependencies(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="{ value: t('packages.tooltip_delete'), showDelay: 300 }" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="isEditing ? t('packages.edit_package') : t('packages.new_package')" modal :style="{ width: '920px' }" :closable="true" @hide="() => (errors.value = {})">
      <div class="grid grid-cols-2 gap-6">
        <!-- Left Column: Package Data -->
        <div class="space-y-3 pr-4 border-r border-gray-200">
          <div class="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">{{ t('packages.package_info') }}</div>
          
          <FormRow :label="t('packages.package_name')" :errors="errors?.name" :required="true">
            <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
          </FormRow>
          
          <div class="grid grid-cols-2 gap-3">
            <FormRow :label="t('packages.slug')" :errors="errors?.slug" :required="true">
              <div class="relative slug-wrapper" :class="{ 'slug-highlight': slugHighlighted }">
                <TextInput v-model="form.slug" :class="{ 'p-invalid': errors?.slug }" @input="slugManuallyEdited = true" />
                <i v-if="!slugManuallyEdited && form.name" class="pi pi-sync text-[10px] absolute right-2 top-1/2 -translate-y-1/2 transition-colors duration-300" :class="slugHighlighted ? 'text-indigo-400' : 'text-gray-300'" v-tooltip.left="'Auto-generated from name'"></i>
              </div>
            </FormRow>
            <FormRow :label="t('packages.price')" :errors="errors?.price" :required="true">
              <InputNumber v-model="form.price" :min="0" :step="10000" inputClass="!w-full !text-sm" class="!w-full" :class="{ 'p-invalid': errors?.price }" />
            </FormRow>
          </div>
          
          <FormRow :label="t('common.description')" :errors="errors?.description">
            <TextInput v-model="form.description" textarea :rows="4" :class="{ 'p-invalid': errors?.description }" />
          </FormRow>
          
          <FormRow :label="t('packages.sort_order')" :errors="errors?.sort_order">
            <InputNumber v-model="form.sort_order" :min="0" inputClass="!w-full !text-sm" class="!w-full" />
          </FormRow>
        </div>

        <!-- Right Column: Module Selector -->
        <div class="space-y-2">
          <div class="flex items-center justify-between mb-2">
            <div class="text-xs font-semibold text-gray-500 uppercase tracking-wider">{{ t('packages.select_modules') }} <span class="text-gray-400 font-normal normal-case">({{ selectedModuleIds.length }} {{ t('common.selected') }})</span></div>
            <Button 
              v-if="availableModules.length > 0"
              :label="allModulesSelected ? t('common.deselect_all') : t('common.select_all')"
              icon="pi pi-check-square"
              size="small"
              text
              severity="secondary"
              class="!text-xs !px-2 !py-1"
              @click="toggleSelectAll"
            />
          </div>
          
          <div class="space-y-1.5 max-h-[400px] overflow-y-auto pr-1">
            <div 
              v-for="mod in availableModules" 
              :key="mod.value" 
              class="border rounded-lg transition-all duration-150"
              :class="selectedModuleIds.includes(mod.value) 
                ? 'border-indigo-200 bg-indigo-50/40 shadow-sm' 
                : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'">
              
              <!-- Module Header (always visible) -->
              <div class="flex items-center gap-2 px-2.5 py-2 cursor-pointer" @click="toggleModule(mod.value)">
                <ToggleSwitch 
                  :modelValue="selectedModuleIds.includes(mod.value)" 
                  @update:modelValue="v => toggleModule(mod.value)" 
                  @click.stop
                  class="!shrink-0"
                />
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-1.5">
                    <span class="text-sm font-medium text-gray-700 truncate">{{ mod.module_name }}</span>
                    <Tag 
                      :value="mod.module_type === 'platform' ? 'P' : 'T'" 
                      :severity="mod.module_type === 'platform' ? 'info' : 'success'" 
                      class="!text-[10px] !px-1.5 !py-0 shrink-0" 
                    />
                  </div>
                  <div class="text-xs text-gray-400 truncate">{{ mod.module_slug }}</div>
                </div>
                <i 
                  v-if="selectedModuleIds.includes(mod.value)"
                  class="pi pi-chevron-down text-xs text-gray-400 transition-transform duration-150"
                  :class="{ 'rotate-180': expandedModules[mod.value] }"
                />
              </div>

              <!-- Module Details (shown when selected and expanded) -->
              <div v-if="selectedModuleIds.includes(mod.value) && expandedModules[mod.value]" class="px-2.5 pb-2.5 space-y-1.5">
                <div class="border-t border-indigo-100 pt-1.5 space-y-1">
                  <div v-if="mod.module_description" class="flex items-start gap-1.5">
                    <i class="pi pi-align-left text-[10px] text-gray-400 mt-0.5"></i>
                    <p class="text-xs text-gray-500 leading-relaxed">{{ mod.module_description }}</p>
                  </div>
                  <div v-if="mod.depends_on" class="flex items-start gap-1.5">
                    <i class="pi pi-sitemap text-[10px] text-gray-400 mt-0.5"></i>
                    <p class="text-xs">
                      <span class="text-gray-400">{{ t('packages.depends_on') }}:</span>
                      <span class="text-amber-600 ml-1">{{ mod.depends_on }}</span>
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-3 pt-1">
                  <div class="flex items-center gap-1.5">
                    <ToggleSwitch 
                      :modelValue="moduleMandatory[mod.value] || false" 
                      @update:modelValue="v => moduleMandatory[mod.value] = v" 
                      class="!scale-75 !origin-left"
                    />
                    <span class="text-[11px] text-gray-500 whitespace-nowrap">{{ t('packages.is_mandatory') }}</span>
                  </div>
                  <div class="flex items-center gap-1.5 ml-auto">
                    <label class="text-[11px] text-gray-500 whitespace-nowrap">{{ t('packages.sort_order') }}</label>
                    <InputNumber 
                      :modelValue="moduleSortOrders[mod.value] || 0" 
                      @update:modelValue="v => moduleSortOrders[mod.value] = v" 
                      :min="0" 
                      class="!w-16" 
                      inputClass="!w-full !text-xs !h-6" 
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="availableModules.length === 0" class="text-sm text-gray-400 text-center py-6 border border-dashed border-gray-200 rounded-lg">
            <i class="pi pi-spinner pi-spin mr-1"></i>
            {{ t('packages.loading_modules') }}
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ml-auto">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
            <Button :label="isEditing ? t('common.update') : t('common.create')" size="small" :loading="saving" :disabled="saving" @click="savePackage" />
          </div>
        </div>
      </template>
    </Dialog>

    <!-- Validate Dependencies Dialog -->
    <Dialog v-model:visible="depsDialogVisible" :header="t('packages.validate_deps')" modal :style="{ width: '500px' }">
      <div class="space-y-2">
        <div v-for="dep in dependencies" :key="dep.module_id" class="flex items-center gap-3 px-3 py-2 rounded-md text-sm" :class="dep.resolved ? 'bg-emerald-50' : 'bg-rose-50'">
          <i :class="dep.resolved ? 'pi pi-check-circle text-emerald-500' : 'pi pi-exclamation-circle text-rose-500'" class="text-base"></i>
          <div class="flex-1">
            <span class="font-medium">{{ dep.module_name }}</span>
            <span class="text-gray-500 ml-2">{{ dep.depends_on }}</span>
          </div>
          <Tag :value="dep.resolved ? t('packages.deps_resolved') : t('packages.deps_unresolved')" :severity="dep.resolved ? 'success' : 'danger'" class="!text-xs" />
        </div>
        <div v-if="dependencies.length === 0" class="text-sm text-gray-400 text-center py-4">{{ t('common.no_data') }}</div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" text size="small" @click="depsDialogVisible = false" />
      </template>
    </Dialog>

    <!-- Confirm Dialog -->
    <Dialog v-model:visible="confirmVisible" :header="confirmTitle" modal :style="{ width: '400px' }">
      <p class="text-sm text-gray-600">{{ confirmMessage }}</p>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="confirmVisible = false" />
        <Button :label="confirmActionLabel" :severity="confirmSeverity" size="small" :loading="confirming" :disabled="confirming" @click="executeConfirm" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'

import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import { useSlugify } from '@/composables/useSlugify'

const toast = useToast()
const { t } = useI18n()

// Data
const packages = ref([])
const loading = ref(true)
const availableModules = ref([])
const searchQuery = ref('')
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const form = ref({ name: '', slug: '', price: 0, description: '', sort_order: 0 })
const errors = ref({})
const selectedModuleIds = ref([])
const expandedModules = reactive({})
const { slugManuallyEdited, slugHighlighted, resetSlug, disableAutoSlug } = useSlugify(
  () => form.value.name,
  (v) => { form.value.slug = v }
)

// Confirm dialog state
const confirmVisible = ref(false)
const confirmAction = ref(null)
const confirmTarget = ref(null)
const confirming = ref(false)

// Dependencies dialog
const depsDialogVisible = ref(false)
const dependencies = ref([])

// Filtered packages (client-side search only)
const filteredPackages = computed(() => {
  let result = packages.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(p =>
      p.name?.toLowerCase().includes(q) ||
      p.slug?.toLowerCase().includes(q) ||
      p.description?.toLowerCase().includes(q)
    )
  }
  return result
})

// Track sort orders and mandatory flags for selected modules
const moduleSortOrders = ref({})
const moduleMandatory = ref({})

// Check if all available modules are selected
const allModulesSelected = computed(() => {
  return availableModules.value.length > 0 && selectedModuleIds.value.length === availableModules.value.length
})

// Select or deselect all modules
function toggleSelectAll() {
  if (allModulesSelected.value) {
    // Deselect all
    selectedModuleIds.value = []
    Object.keys(expandedModules).forEach(k => delete expandedModules[k])
  } else {
    // Select all — keep existing selections, add missing ones
    const currentIds = new Set(selectedModuleIds.value)
    availableModules.value.forEach(mod => {
      if (!currentIds.has(mod.value)) {
        selectedModuleIds.value.push(mod.value)
        expandedModules[mod.value] = true
      }
    })
  }
}



// Format price to IDR
function formatPrice(price) {
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(price || 0)
}

// Load packages from API
async function loadPackages() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/platform/packages?per_page=100')
    const payload = res.data
    packages.value = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
  } finally {
    loading.value = false
  }
}

// Toggle module selection
function toggleModule(id) {
  const idx = selectedModuleIds.value.indexOf(id)
  if (idx > -1) {
    selectedModuleIds.value.splice(idx, 1)
    delete expandedModules[id]
  } else {
    selectedModuleIds.value.push(id)
    expandedModules[id] = true
  }
}

// Load available modules for the module selector
async function loadModules() {
  try {
    const res = await api.get('/api/v1/platform/modules?per_page=100')
    const payload = res.data
    const mods = Array.isArray(payload.data) ? payload.data : (Array.isArray(payload) ? payload : [])
    availableModules.value = mods.map(m => ({
      label: `${m.name} (${m.slug})`,
      value: m.id,
      module_name: m.name,
      module_slug: m.slug,
      module_type: m.module_type,
      module_description: m.description || '',
      depends_on: m.depends_on || ''
    }))
  } catch (e) {
    availableModules.value = []
  }
}

onMounted(async () => {
  await Promise.all([loadPackages(), loadModules()])
})

function openCreate() {
  isEditing.value = false
  editingId.value = null
  form.value = { name: '', slug: '', price: 0, description: '', sort_order: 0 }
  errors.value = {}
  selectedModuleIds.value = []
  moduleSortOrders.value = {}
  moduleMandatory.value = {}
  resetSlug()
  Object.keys(expandedModules).forEach(k => delete expandedModules[k])
  dialogVisible.value = true
}

function openEdit(pkg) {
  isEditing.value = true
  editingId.value = pkg.id
  form.value = { name: pkg.name, slug: pkg.slug, price: pkg.price, description: pkg.description || '', sort_order: pkg.sort_order }
  errors.value = {}
  disableAutoSlug() // disable auto-slug when editing
  
  // Pre-select modules
  if (pkg.modules && pkg.modules.length > 0) {
    selectedModuleIds.value = pkg.modules.map(m => m.module_id)
    const orders = {}
    const mandatory = {}
    pkg.modules.forEach(m => {
      orders[m.module_id] = m.sort_order
      mandatory[m.module_id] = m.is_mandatory
    })
    moduleSortOrders.value = orders
    moduleMandatory.value = mandatory
  } else {
    selectedModuleIds.value = []
    moduleSortOrders.value = {}
    moduleMandatory.value = {}
  }
  
  // Auto-expand all selected modules
  Object.keys(expandedModules).forEach(k => delete expandedModules[k])
  selectedModuleIds.value.forEach(id => { expandedModules[id] = true })
  
  dialogVisible.value = true
}

// Build package modules payload
function buildModulesPayload() {
  return selectedModuleIds.value.map(id => ({
    module_id: id,
    is_mandatory: moduleMandatory.value[id] || false,
    sort_order: moduleSortOrders.value[id] || 0
  }))
}

async function savePackage() {
  saving.value = true
  try {
    const payload = {
      ...form.value,
      modules: buildModulesPayload()
    }
    
    if (isEditing.value) {
      await api.put(`/api/v1/platform/packages/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.updated'), life: 2000 })
    } else {
      await api.post('/api/v1/platform/packages', payload)
      toast.add({ severity: 'success', summary: t('message.created'), life: 2000 })
    }
    dialogVisible.value = false
    await loadPackages()
  } catch (e) {
    errors.value = getValidationErrors(e)
    if (Object.keys(errors.value).length === 0) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
    }
  } finally {
    saving.value = false
  }
}

// Confirm actions (publish/unpublish/delete)
const confirmTitle = computed(() => {
  switch (confirmAction.value) {
    case 'publish': return t('packages.publish')
    case 'unpublish': return t('packages.unpublish')
    case 'delete': return t('common.delete')
    default: return ''
  }
})
const confirmMessage = computed(() => {
  if (!confirmTarget.value) return ''
  const name = confirmTarget.value.name
  switch (confirmAction.value) {
    case 'publish': return t('packages.confirm_publish_message', { name })
    case 'unpublish': return t('packages.confirm_unpublish_message', { name })
    case 'delete': return t('packages.confirm_delete_message', { name })
    default: return ''
  }
})
const confirmActionLabel = computed(() => {
  switch (confirmAction.value) {
    case 'publish': return t('packages.publish')
    case 'unpublish': return t('packages.unpublish')
    case 'delete': return t('common.delete')
    default: return ''
  }
})
const confirmSeverity = computed(() => {
  return confirmAction.value === 'delete' ? 'danger' : 'warn'
})

function confirmPublish(pkg) { confirmAction.value = 'publish'; confirmTarget.value = pkg; confirmVisible.value = true }
function confirmUnpublish(pkg) { confirmAction.value = 'unpublish'; confirmTarget.value = pkg; confirmVisible.value = true }
function confirmDelete(pkg) { confirmAction.value = 'delete'; confirmTarget.value = pkg; confirmVisible.value = true }

async function executeConfirm() {
  if (!confirmTarget.value || !confirmAction.value) return
  confirming.value = true
  const id = confirmTarget.value.id
  try {
    switch (confirmAction.value) {
      case 'publish':
        await api.post(`/api/v1/platform/packages/${id}/publish`)
        toast.add({ severity: 'success', summary: t('message.updated'), detail: t('packages.publish'), life: 2000 })
        break
      case 'unpublish':
        await api.post(`/api/v1/platform/packages/${id}/unpublish`)
        toast.add({ severity: 'success', summary: t('message.updated'), detail: t('packages.unpublish'), life: 2000 })
        break
      case 'delete':
        await api.delete(`/api/v1/platform/packages/${id}`)
        toast.add({ severity: 'success', summary: t('message.deleted'), life: 2000 })
        break
    }
    confirmVisible.value = false
    await loadPackages()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  } finally {
    confirming.value = false
  }
}

// Validate dependencies
async function validateDependencies(pkg) {
  try {
    const res = await api.get(`/api/v1/platform/packages/${pkg.id}/validate`)
    const deps = res.data?.data || (Array.isArray(res.data) ? res.data : [])
    dependencies.value = deps
    depsDialogVisible.value = true
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  }
}
</script>

