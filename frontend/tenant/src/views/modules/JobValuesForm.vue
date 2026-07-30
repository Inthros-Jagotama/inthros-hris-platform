<template>
  <div class="max-w-full mx-auto">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="flex gap-4">
        <div class="w-56 space-y-2">
          <div v-for="n in 5" :key="n" class="h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
        <div class="flex-1 space-y-3">
          <div v-for="n in 6" :key="n" class="h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- Two-column layout -->
    <div v-else class="flex gap-6">
      <!-- Left: Navigation Sidebar -->
      <div class="w-56 shrink-0 space-y-1">
        <!-- Back to Job Management -->
        <div
          role="button"
          tabindex="0"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
          @click="goBack"
          @keydown.enter="goBack"
        >
          <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300">
            <i class="pi pi-arrow-left text-xs"></i>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-gray-700 dark:text-gray-300 truncate">{{ t('common.back') }}</div>
            <div class="text-[10px] text-gray-400">{{ t('nav.job_management') }}</div>
          </div>
        </div>

        <!-- Divider -->
        <div class="border-t border-gray-200 dark:border-gray-700 my-2"></div>

        <!-- Organization Info -->
        <div class="px-3 py-2">
          <div class="text-[10px] text-gray-400 uppercase tracking-wider mb-1">{{ t('organization.title') }}</div>
          <div class="text-sm font-semibold text-gray-800 dark:text-gray-100 truncate">{{ orgName }}</div>
          <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 font-mono">{{ orgCode }}</div>
        </div>

        <!-- Divider -->
        <div class="border-t border-gray-200 dark:border-gray-700 my-2"></div>

        <!-- Active Section: Nilai Jabatan -->
        <div
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700"
        >
          <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 bg-emerald-600 text-white">
            <i class="pi pi-sliders-h text-xs"></i>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium text-emerald-700 dark:text-emerald-300 truncate">{{ t('job_management.values') }}</div>
            <div class="text-[10px] text-emerald-500">{{ items.length }} {{ t('common.items') }}</div>
          </div>
        </div>
      </div>

      <!-- Right: Form Content -->
      <div class="flex-1 min-w-0">
        <!-- Header -->
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.values') }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_values.description') }}</p>
          </div>
          <Button :label="t('common.create')" icon="pi pi-plus" size="small" @click="openCreateDialog()" />
        </div>

        <!-- DataTable -->
        <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

        <DataTable
          v-else
          :value="items"
          lazy
          :totalRecords="totalRecords"
          :first="firstRecord"
          :rows="perPage"
          @page="onPage($event)"
          paginator
          paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
          :rowsPerPageOptions="[10, 15, 25, 50]"
          size="small"
          class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        >
          <template #empty>
            <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
              <i class="pi pi-sliders-h text-3xl mb-2 opacity-50"></i>
              <p class="text-sm font-medium">{{ t('job_values.empty_title') }}</p>
              <p class="text-xs mt-1">{{ t('job_values.empty_hint') }}</p>
            </div>
          </template>
          <Column field="type" :header="t('job_values.type')" sortable style="width:150px">
            <template #body="{data}"><Tag :value="data.type" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
          </Column>
          <Column field="level" :header="t('job_values.level')" sortable style="width:80px">
            <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.level || '-' }}</span></template>
          </Column>
          <Column field="descriptions" :header="t('job_values.descriptions')" sortable>
            <template #body="{data}"><span class="text-gray-800 dark:text-gray-100">{{ data.descriptions || '-' }}</span></template>
          </Column>
          <Column field="sort" :header="t('job_values.sort')" sortable style="width:80px">
            <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort || 0 }}</span></template>
          </Column>
          <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
            <template #body="{data}">
              <div class="flex items-center gap-1">
                <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEditDialog(data)" />
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
              </div>
            </template>
          </Column>
        </DataTable>

        <!-- Create/Edit Dialog -->
        <Dialog v-model:visible="dialogVisible" :header="editing ? t('job_values.edit') : t('job_values.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
          <div class="space-y-4">
            <FormRow :label="t('job_values.type')" required :errors="errors?.type">
              <SelectLabel v-model="form.type" :options="typeOptions" :class="{ 'p-invalid': errors?.type }" :placeholder="t('job_values.select_type')" />
            </FormRow>
            <FormRow :label="t('job_values.level')" :errors="errors?.level">
              <InputNumber v-model="form.level" class="!w-full" :min="0" size="small" :class="{ 'p-invalid': errors?.level }" />
            </FormRow>
            <FormRow :label="t('job_values.descriptions')" :errors="errors?.descriptions">
              <Textarea v-model="form.descriptions" :class="{ 'p-invalid': errors?.descriptions }" :placeholder="t('job_values.descriptions_placeholder')" rows="3" />
            </FormRow>
            <FormRow :label="t('job_values.note')" :errors="errors?.note">
              <Textarea v-model="form.note" :class="{ 'p-invalid': errors?.note }" :placeholder="t('job_values.note_placeholder')" rows="2" />
            </FormRow>
            <FormRow :label="t('job_values.sort')" :errors="errors?.sort">
              <InputNumber v-model="form.sort" class="!w-full" :min="0" size="small" :class="{ 'p-invalid': errors?.sort }" />
            </FormRow>
          </div>
          <template #footer>
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
              <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
            </div>
          </template>
        </Dialog>

        <ConfirmDeleteDialog
          v-model:visible="deleteDialogVisible"
          :title="t('job_values.confirm_delete_title')"
          :message="t('job_values.confirm_delete', { type: deleteTarget?.type || '' })"
          :loading="deleting"
          :error-msg="deleteError"
          :confirm-label="t('common.delete')"
          :cancel-label="t('common.cancel')"
          @confirm="handleDelete"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const toast = useToast()

const orgId = route.query.org_id || ''
const orgName = ref('')
const orgCode = ref('')
const pageLoading = ref(true)

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const form = ref({ type: '', level: null, descriptions: '', note: '', sort: null })

const typeOptions = [
  { label: 'Education', value: 'education' },
  { label: 'Experience', value: 'experience' },
  { label: 'Environment', value: 'environment' },
  { label: 'Hazard', value: 'hazard' },
  { label: 'Relationship', value: 'relationship' },
  { label: 'Frequency', value: 'frequency' },
  { label: 'Asset', value: 'asset' },
  { label: 'Authority', value: 'authority' },
  { label: 'Cash', value: 'cash' },
  { label: 'Impact', value: 'impact' }
]

const skeletonColumns = [
  { type: 'tag', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-48', headerWidth: 'w-24' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// ── Fetch Organization Info ──
async function loadOrgInfo() {
  if (!orgId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${orgId}`)
    const data = res.data?.data
    if (data) {
      orgName.value = data.nomenclature
      orgCode.value = `${data.full_code}`
    }
  } catch {
    orgName.value = 'Unknown'
    orgCode.value = ''
  }
}

// ── Load Job Values ──
async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/job-management/values', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function goBack() {
  router.push('/job-management')
}

function resetForm() {
  form.value = { type: '', level: null, descriptions: '', note: '', sort: null }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

function openCreateDialog() {
  resetForm()
  dialogVisible.value = true
}

function openEditDialog(item) {
  editing.value = true
  editingId.value = item.id
  errors.value = {}
  form.value = {
    type: item.type || '',
    level: item.level ?? null,
    descriptions: item.descriptions || '',
    note: item.note || '',
    sort: item.sort ?? null
  }
  dialogVisible.value = true
}

async function handleSave() {
  errors.value = {}
  if (!form.value.type?.trim()) { errors.value = { type: [t('form.required')] }; return }

  saving.value = true
  try {
    const payload = {
      type: form.value.type,
      level: form.value.level ?? null,
      descriptions: form.value.descriptions || '',
      note: form.value.note || '',
      sort: form.value.sort ?? null
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/job-management/values/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/job-management/values', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.created'), life: 3000 })
    }

    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/job-management/values/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  pageLoading.value = true
  try {
    await loadOrgInfo()
    await loadData()
  } finally {
    pageLoading.value = false
  }
})
</script>
