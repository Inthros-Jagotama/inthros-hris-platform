<template>
  <div class="space-y-3">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
        <Select
          v-model="filterDocumentType"
          :options="documentTypeOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('document_templates.filter_document_type')"
          class="!text-sm w-56"
          showClear
          @change="onFilterChange"
        />
        <Select
          v-model="filterStatus"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          :placeholder="t('document_templates.filter_status')"
          class="!text-sm w-40"
          showClear
          @change="onFilterChange"
        />
        <span class="p-input-icon-left">
          <TextInput
            v-model="searchTerm"
            :placeholder="t('common.search')"
            class="!text-sm w-56"
            @update:modelValue="onSearchInput"
          />
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('document_templates.new')" icon="pi pi-plus" size="small" @click="goCreate" />
      </div>
    </div>

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
          <i class="pi pi-file-edit text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('document_templates.empty_title') }}</p>
        </div>
      </template>
      <Column field="name" :header="t('document_templates.name')">
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
          </div>
        </template>
      </Column>
      <Column field="code" :header="t('document_templates.code')">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ data.code }}</span>
        </template>
      </Column>
      <Column field="document_type" :header="t('document_templates.document_type')">
        <template #body="{ data }">
          <div class="flex flex-col">
            <span class="text-gray-700 dark:text-gray-200">{{ documentTypeLabel(data.document_type) }}</span>
            <span v-if="data.movement_type" class="text-xs text-teal-500 dark:text-teal-400">{{ movementTypeLabel(data.movement_type) }}</span>
          </div>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:120px">
        <template #body="{ data }">
          <Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="updated_at" :header="t('document_templates.updated_at')" style="width:160px">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.updated_at, locale) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:230px" frozen alignFrozen="right">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.view')" @click="openDetail(data)" />
            <Button
              icon="pi pi-file-pdf"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('document_templates.preview')"
              :loading="previewingId === data.id"
              @click="openPreview(data)"
            />
            <Button
              icon="pi pi-pencil"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('common.edit')"
              @click="goEdit(data)"
            />
            <Button
              icon="pi pi-history"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('document_templates.versions')"
              @click="openVersions(data)"
            />
            <Button
              v-if="data.status !== 'ACTIVE'"
              icon="pi pi-check-circle"
              size="small"
              text
              severity="success"
              v-tooltip.left="t('document_templates.activate')"
              :loading="actioningId === data.id"
              @click="confirmActivate(data)"
            />
            <Button
              v-if="data.status === 'ACTIVE'"
              icon="pi pi-times-circle"
              size="small"
              text
              severity="secondary"
              v-tooltip.left="t('document_templates.deactivate')"
              :loading="actioningId === data.id"
              @click="confirmDeactivate(data)"
            />
            <Button
              icon="pi pi-trash"
              size="small"
              text
              severity="danger"
              v-tooltip.left="t('common.delete')"
              @click="confirmDelete(data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ── Detail Template ── -->
    <Dialog v-model:visible="detailVisible" :header="t('document_templates.detail')" modal :style="{ width: '720px' }" :closable="true">
      <div v-if="detail" class="space-y-4">
        <div class="flex items-center justify-between gap-2">
          <h3 class="text-base font-semibold text-navy-800 dark:text-gray-100">{{ detail.name }}</h3>
          <Tag :value="statusLabel(detail.status)" :severity="statusSeverity(detail.status)" class="!text-xs !px-1.5 !py-0.5" />
        </div>
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.code') }}</p>
            <p class="font-mono text-gray-700 dark:text-gray-200">{{ detail.code }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.document_type') }}</p>
            <p class="text-gray-700 dark:text-gray-200">{{ documentTypeLabel(detail.document_type) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.active_version') }}</p>
            <p class="text-gray-700 dark:text-gray-200">
              {{ detail.active_version_id ? activeVersionLabel : '—' }}
            </p>
          </div>
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.updated_at') }}</p>
            <p class="text-gray-700 dark:text-gray-200">{{ formatDate(detail.updated_at, locale) }}</p>
          </div>
          <div v-if="detail.description" class="col-span-2">
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.description_label') }}</p>
            <p class="text-gray-700 dark:text-gray-200 whitespace-pre-wrap">{{ detail.description }}</p>
          </div>
        </div>
        <div>
          <p class="text-xs font-medium text-gray-500 dark:text-gray-400 mb-1.5">{{ t('document_templates.template_file') }}</p>
          <div v-if="detailFile" class="flex items-center justify-between gap-3 rounded-lg border border-gray-200 dark:border-gray-700 px-4 py-3">
            <div class="flex items-center gap-3 min-w-0">
              <i class="pi pi-file-word text-xl text-emerald-500"></i>
              <div class="min-w-0">
                <p class="text-sm font-medium text-navy-800 dark:text-gray-100 truncate">{{ detailFile.name }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.active_version') }} · v{{ activeVersionNumber }}</p>
              </div>
            </div>
            <a :href="detailFile.url" :download="detailFile.name" class="text-sm text-teal-500 hover:text-teal-600 dark:text-teal-400 font-medium whitespace-nowrap">
              <i class="pi pi-download mr-1"></i>{{ t('common.download') }}
            </a>
          </div>
          <div v-else class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800 max-h-72 overflow-auto">
            <div v-if="detailContent" class="editor-content text-sm" v-html="detailContent"></div>
            <p v-else class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.no_content') }}</p>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-between">
          <Button :label="t('document_templates.versions')" icon="pi pi-history" size="small" severity="secondary" outlined @click="openVersions(detail)" />
          <div class="flex items-center gap-2 ml-auto">
            <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="detailVisible = false" />
          </div>
        </div>
      </template>
    </Dialog>

    <!-- ── Versions Management ── -->
    <Dialog v-model:visible="versionsVisible" :header="t('document_templates.versions')" modal :style="{ width: '760px' }" :closable="!savingVersion">
      <div v-if="versionsTarget" class="space-y-3">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-600 dark:text-gray-300">
            <span class="font-medium">{{ versionsTarget.name }}</span>
            <span class="text-gray-400"> · {{ documentTypeLabel(versionsTarget.document_type) }}</span>
          </div>
          <Button :label="t('document_templates.new_version')" icon="pi pi-plus" size="small" @click="openCreateVersion()" />
        </div>
        <DataTable :value="versions" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
          <template #empty>
            <div class="py-8 text-center text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.no_versions') }}</div>
          </template>
          <Column field="version" :header="t('document_templates.version')" style="width:90px">
            <template #body="{ data }">
              <Tag :value="`v${data.version}`" severity="info" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column field="paper_size" :header="t('document_templates.paper_size')" style="width:110px">
            <template #body="{ data }">
              <span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.paper_size }} / {{ data.orientation }}</span>
            </template>
          </Column>
          <Column field="created_at" :header="t('document_templates.created_at')" style="width:150px">
            <template #body="{ data }">
              <span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.created_at, locale) }}</span>
            </template>
          </Column>
          <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
            <template #body="{ data }">
              <div class="flex items-center gap-1">
                <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.view')" @click="openVersionDetail(data)" />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
      <template #footer>
        <div class="flex items-center justify-end">
          <Button :label="t('common.close')" severity="secondary" outlined size="small" :disabled="savingVersion" @click="versionsVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── New Version ── -->
    <Dialog v-model:visible="createVersionVisible" :header="t('document_templates.new_version')" modal :style="{ width: '640px' }" :closable="!savingVersion" @hide="resetCreateVersion">
      <div class="space-y-3">
        <FormRow :label="t('document_templates.template_file')" required :errors="errors?.file">
          <div class="flex items-center gap-3">
            <label class="flex items-center gap-2 cursor-pointer text-sm px-4 py-2 rounded-lg border border-dashed border-gray-300 dark:border-gray-600 hover:border-teal-400 hover:bg-teal-50/50 dark:hover:bg-teal-500/5 transition-colors" :class="{ 'p-invalid': errors?.file }">
              <i class="pi pi-upload text-teal-500"></i>
              <span class="text-gray-600 dark:text-gray-300">{{ t('document_templates.choose_file') }}</span>
              <input type="file" accept=".docx,application/vnd.openxmlformats-officedocument.wordprocessingml.document" class="hidden" @change="onVersionFileChange" />
            </label>
            <span v-if="versionForm.file" class="text-sm text-gray-600 dark:text-gray-300 font-medium truncate">{{ versionForm.file.name }}</span>
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ t('document_templates.docx_hint') }}</p>
        </FormRow>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('document_templates.paper_size')">
            <Select v-model="versionForm.paper_size" :options="paperOptions" optionLabel="label" optionValue="value" class="!w-full" />
          </FormRow>
          <FormRow :label="t('document_templates.orientation')">
            <Select v-model="versionForm.orientation" :options="orientationOptions" optionLabel="label" optionValue="value" class="!w-full" />
          </FormRow>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" :disabled="savingVersion" @click="createVersionVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="savingVersion" :disabled="savingVersion" @click="handleCreateVersion" />
        </div>
      </template>
    </Dialog>

    <!-- ── Preview PDF ── -->
    <Dialog v-model:visible="previewVisible" :header="t('document_templates.preview')" modal :style="{ width: '900px' }" :closable="!previewLoading" :contentStyle="{ height: '75vh' }">
      <div v-if="previewLoading" class="flex flex-col items-center justify-center h-full gap-3 text-gray-400 dark:text-gray-500">
        <i class="pi pi-spin pi-spinner text-2xl"></i>
        <p class="text-sm">{{ t('document_templates.preview_generating') }}</p>
      </div>
      <div v-else-if="previewError" class="flex flex-col items-center justify-center h-full gap-2 text-center px-6">
        <i class="pi pi-exclamation-triangle text-2xl text-amber-500"></i>
        <p class="text-sm text-gray-600 dark:text-gray-300">{{ previewError }}</p>
      </div>
      <iframe
        v-else-if="previewUrl"
        :src="previewUrl"
        id="preview-iframe"
        class="w-full h-full rounded-lg border border-gray-200 dark:border-gray-700"
        title="Preview PDF"
      ></iframe>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button
            :label="t('common.print')"
            icon="pi pi-print"
            size="small"
            severity="secondary"
            outlined
            :disabled="!previewUrl || previewLoading"
            @click="printPreview"
          />
          <a
            :href="previewUrl"
            :download="previewFileName"
            class="inline-flex items-center gap-1 text-sm px-3 py-1.5 rounded-md bg-teal-600 hover:bg-teal-700 text-white font-medium disabled:opacity-50"
            :class="{ 'pointer-events-none opacity-50': !previewUrl || previewLoading }"
          >
            <i class="pi pi-download"></i>{{ t('common.download') }}
          </a>
          <Button :label="t('common.close')" severity="secondary" outlined size="small" :disabled="previewLoading" @click="previewVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── Version Detail ── -->
    <Dialog v-model:visible="versionDetailVisible" :header="`${t('document_templates.version')} v${versionDetail?.version || ''}`" modal :style="{ width: '720px' }" :closable="true">
      <div v-if="versionDetail" class="space-y-3">
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.paper_size') }}</p>
            <p class="text-gray-700 dark:text-gray-200">{{ versionDetail.paper_size }} / {{ versionDetail.orientation }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.created_at') }}</p>
            <p class="text-gray-700 dark:text-gray-200">{{ formatDate(versionDetail.created_at, locale) }}</p>
          </div>
          <div class="col-span-2">
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('document_templates.template_file') }}</p>
            <div v-if="versionDetail.file_url" class="flex items-center justify-between gap-3 rounded-lg border border-gray-200 dark:border-gray-700 px-4 py-3">
              <div class="flex items-center gap-3 min-w-0">
                <i class="pi pi-file-word text-xl text-emerald-500"></i>
                <p class="text-sm font-medium text-navy-800 dark:text-gray-100 truncate">{{ versionDetail.file_name || 'template.docx' }}</p>
              </div>
              <a :href="versionDetail.file_url" :download="versionDetail.file_name || 'template.docx'" class="text-sm text-teal-500 hover:text-teal-600 dark:text-teal-400 font-medium whitespace-nowrap">
                <i class="pi pi-download mr-1"></i>{{ t('common.download') }}
              </a>
            </div>
            <div v-else class="border border-gray-200 dark:border-gray-700 rounded-lg p-4 bg-white dark:bg-gray-800 max-h-80 overflow-auto">
              <div class="editor-content text-sm" v-html="versionDetail.content"></div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end">
          <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="versionDetailVisible = false" />
        </div>
      </template>
    </Dialog>

    <!-- ── Confirmations ── -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('document_templates.delete_confirm_title')"
      :message="deleteConfirmMessage"
      :loading="deleting"
      :error-msg="deleteError"
      confirm-label="Delete"
      @confirm="handleDelete"
      @cancel="deleteDialogVisible = false"
    />
    <ConfirmDeleteDialog
      v-model:visible="activateDialogVisible"
      :title="t('document_templates.activate_confirm_title')"
      :message="activateConfirmMessage"
      :loading="actioningId === activateTarget?.id"
      confirm-label="Activate"
      confirm-severity="success"
      @confirm="handleActivate"
      @cancel="activateDialogVisible = false"
    />
    <ConfirmDeleteDialog
      v-model:visible="deactivateDialogVisible"
      :title="t('document_templates.deactivate_confirm_title')"
      :message="deactivateConfirmMessage"
      :loading="actioningId === deactivateTarget?.id"
      confirm-label="Deactivate"
      confirm-severity="secondary"
      @confirm="handleDeactivate"
      @cancel="deactivateDialogVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import { formatDate } from '@/utils/formatDate'

const { t, locale } = useI18n()
const router = useRouter()
const toast = useToast()

// ── List state ──
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const filterDocumentType = ref(null)
const filterStatus = ref(null)
const searchTerm = ref('')
let searchDebounce = null

const actioningId = ref(null)

const documentTypeOptions = [
  { label: t('document_templates.type_contract_agreement'), value: 'CONTRACT_AGREEMENT' },
  { label: t('document_templates.type_movement_sk'), value: 'MOVEMENT_SK' },
]

const statusOptions = [
  { label: t('document_templates.status_active'), value: 'ACTIVE' },
  { label: t('document_templates.status_inactive'), value: 'INACTIVE' },
]

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

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'icons', count: 5, headerWidth: 'w-16' },
]

// ── Detail ──
const detailVisible = ref(false)
const detail = ref(null)
const detailVersions = ref([])

// ── Preview ──
const previewVisible = ref(false)
const previewingId = ref(null)
const previewLoading = ref(false)
const previewUrl = ref('')
const previewFileName = ref('preview.pdf')
const previewError = ref('')

// ── Versions ──
const versionsVisible = ref(false)
const versionsTarget = ref(null)
const versions = ref([])
const createVersionVisible = ref(false)
const savingVersion = ref(false)
const versionForm = ref({ file: null, paper_size: 'A4', orientation: 'portrait' })
const versionDetailVisible = ref(false)
const versionDetail = ref(null)

// ── Delete / Activate / Deactivate confirms ──
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const activateDialogVisible = ref(false)
const activateTarget = ref(null)
const deactivateDialogVisible = ref(false)
const deactivateTarget = ref(null)

const deleteConfirmMessage = computed(() => deleteTarget.value ? t('document_templates.delete_confirm', { name: deleteTarget.value.name }) : '')
const activateConfirmMessage = computed(() => activateTarget.value ? t('document_templates.activate_confirm', { name: activateTarget.value.name }) : '')
const deactivateConfirmMessage = computed(() => deactivateTarget.value ? t('document_templates.deactivate_confirm', { name: deactivateTarget.value.name }) : '')

const activeVersionLabel = computed(() => {
  const v = detailVersions.value.find((x) => x.id === detail.value?.active_version_id)
  return v ? `v${v.version}` : '—'
})

const activeVersionNumber = computed(() => {
  const v = detailVersions.value.find((x) => x.id === detail.value?.active_version_id)
  return v ? v.version : ''
})

const detailFile = computed(() => {
  if (!detail.value) return null
  const v = detailVersions.value.find((x) => x.id === detail.value.active_version_id)
  const fileInfo = v || detail.value
  if (fileInfo.file_url) {
    return { name: fileInfo.file_name || 'template.docx', url: fileInfo.file_url }
  }
  return null
})

const detailContent = computed(() => {
  if (!detail.value) return ''
  const v = detailVersions.value.find((x) => x.id === detail.value.active_version_id)
  if (v) return v.content
  return detail.value.content || ''
})

function documentTypeLabel(type) {
  const found = documentTypeOptions.find((o) => o.value === type)
  return found ? found.label : type
}

function movementTypeLabel(type) {
  const key = `employee_movement.type_${type}`
  const label = t(key)
  return label !== key ? label : type
}

function statusLabel(status) {
  const found = statusOptions.find((o) => o.value === status)
  return found ? found.label : status
}

function statusSeverity(status) {
  if (status === 'ACTIVE') return 'success'
  if (status === 'REFERENCE') return 'info'
  return 'warn'
}

// ── List ──
async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/document-templates', {
      params: {
        page: currentPage.value,
        per_page: perPage.value,
        document_type: filterDocumentType.value || undefined,
        status: filterStatus.value || undefined,
        search: searchTerm.value || undefined,
      },
    })
    const body = res.data?.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

function onSearchInput() {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    currentPage.value = 1
    loadData()
  }, 400)
}

// ── Create / Edit (halaman terpisah) ──
function goCreate() {
  router.push('/settings/document-templates/new')
}

function goEdit(item) {
  router.push(`/settings/document-templates/${item.id}/edit`)
}

// ── Detail ──
async function openDetail(item) {
  detail.value = null
  detailVersions.value = []
  detailVisible.value = true
  try {
    const [detailRes, versionsRes] = await Promise.all([
      api.get(`/api/v1/tenant/settings/document-templates/${item.id}`),
      api.get(`/api/v1/tenant/settings/document-templates/${item.id}/versions`),
    ])
    detail.value = detailRes.data?.data
    detailVersions.value = versionsRes.data?.data || []
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  }
}

// ── Versions ──
async function openVersions(item) {
  versionsTarget.value = item
  versions.value = []
  versionsVisible.value = true
  try {
    const res = await api.get(`/api/v1/tenant/settings/document-templates/${item.id}/versions`)
    versions.value = res.data?.data || []
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  }
}

function openCreateVersion() {
  errors.value = {}
  versionForm.value = { file: null, paper_size: 'A4', orientation: 'portrait' }
  createVersionVisible.value = true
}

function resetCreateVersion() {
  versionForm.value = { file: null, paper_size: 'A4', orientation: 'portrait' }
  errors.value = {}
}

function onVersionFileChange(event) {
  const file = event.target.files?.[0] || null
  event.target.value = ''
  errors.value = { ...errors.value, file: undefined }
  if (!file) {
    versionForm.value.file = null
    return
  }
  const name = file.name.toLowerCase()
  if (!name.endsWith('.docx')) {
    versionForm.value.file = null
    errors.value = { ...errors.value, file: [t('document_templates.file_invalid_type')] }
    return
  }
  if (file.size > 10 * 1024 * 1024) {
    versionForm.value.file = null
    errors.value = { ...errors.value, file: [t('document_templates.file_too_large')] }
    return
  }
  versionForm.value.file = file
}

async function handleCreateVersion() {
  if (!versionsTarget.value) return
  errors.value = {}
  if (!versionForm.value.file) { errors.value = { ...errors.value, file: [t('document_templates.file_required')] }; return }
  savingVersion.value = true
  try {
    const fd = new FormData()
    fd.append('file', versionForm.value.file)
    fd.append('paper_size', versionForm.value.paper_size)
    fd.append('orientation', versionForm.value.orientation)
    await api.post(`/api/v1/tenant/settings/document-templates/${versionsTarget.value.id}/versions`, fd, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.version_created'), life: 3000 })
    createVersionVisible.value = false
    await openVersions(versionsTarget.value)
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
    savingVersion.value = false
  }
}

async function openVersionDetail(v) {
  if (!versionsTarget.value) return
  try {
    const res = await api.get(`/api/v1/tenant/settings/document-templates/${versionsTarget.value.id}/versions/${v.id}`)
    versionDetail.value = res.data?.data
    versionDetailVisible.value = true
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.failed_to_load'),
      life: 4000,
    })
  }
}

// ── Preview ──
async function openPreview(item) {
  previewingId.value = item.id
  previewVisible.value = true
  previewLoading.value = true
  previewError.value = ''
  previewUrl.value = ''
  try {
    const res = await api.post(`/api/v1/tenant/settings/document-templates/${item.id}/preview`)
    const data = res.data?.data
    previewUrl.value = data?.pdf_url || ''
    previewFileName.value = data?.file_name || `preview_${item.code}.pdf`
    if (!previewUrl.value) {
      previewError.value = t('document_templates.preview_failed')
    }
  } catch (e) {
    previewError.value = e.response?.data?.error?.message || t('document_templates.preview_failed')
  } finally {
    previewLoading.value = false
    previewingId.value = null
  }
}

function printPreview() {
  const frame = document.querySelector('#preview-iframe')
  if (frame?.contentWindow) {
    frame.contentWindow.print()
  }
}

// ── Delete ──
function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/document-templates/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

// ── Activate / Deactivate ──
function confirmActivate(item) {
  activateTarget.value = item
  activateDialogVisible.value = true
}

function confirmDeactivate(item) {
  deactivateTarget.value = item
  deactivateDialogVisible.value = true
}

async function handleActivate() {
  if (!activateTarget.value) return
  actioningId.value = activateTarget.value.id
  try {
    await api.post(`/api/v1/tenant/settings/document-templates/${activateTarget.value.id}/activate`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.activated'), life: 3000 })
    activateDialogVisible.value = false
    activateTarget.value = null
    await loadData()
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000,
    })
  } finally {
    actioningId.value = null
  }
}

async function handleDeactivate() {
  if (!deactivateTarget.value) return
  actioningId.value = deactivateTarget.value.id
  try {
    await api.post(`/api/v1/tenant/settings/document-templates/${deactivateTarget.value.id}/deactivate`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('document_templates.deactivated'), life: 3000 })
    deactivateDialogVisible.value = false
    deactivateTarget.value = null
    await loadData()
  } catch (e) {
    toast.add({
      severity: 'error',
      summary: t('message.error'),
      detail: e.response?.data?.error?.message || t('message.operation_failed'),
      life: 4000,
    })
  } finally {
    actioningId.value = null
  }
}

onMounted(() => {
  loadData()
})
</script>
