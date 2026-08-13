<template>
  <div class="space-y-1">
    <!-- ── Tabs: Certifications (master) | Certificates (issued) ── -->
    <div class="flex items-center gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
        :class="activeTab === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab(tab.key)"
      >
        {{ t(tab.labelKey) }}
      </button>
    </div>

    <!-- ══ Tab 1: Certification master ══ -->
    <div v-if="activeTab === 'certifications'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2 flex-wrap">
          <span v-if="certTotal > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ certTotal }} {{ t('common.items') }}</span>
        </div>
        <Button :label="t('training.certification_new')" icon="pi pi-plus" size="small" class="ml-auto" @click="openCertificationDialog()" />
      </div>

      <SkeletonTable v-if="certLoading" :columns="skeletonColumns" :rows="6" />
      <DataTable
        v-else
        :value="certifications"
        lazy
        :totalRecords="certTotal"
        :first="certFirst"
        :rows="perPage"
        @page="onCertPage($event)"
        paginator
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        :rowsPerPageOptions="[10, 15, 25, 50]"
        size="small"
        class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-id-card text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.certifications_empty') }}</p>
          </div>
        </template>
        <Column field="code" :header="t('training.certification_code')" style="width:130px">
          <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="name" :header="t('training.certification_name')">
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
        </Column>
        <Column field="issuing_body" :header="t('training.issuing_body')" style="width:180px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.issuing_body || '-' }}</span></template>
        </Column>
        <Column field="validity_period_month" :header="t('training.validity_period')" style="width:140px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.validity_period_unit === 'forever' ? t('training.validity_unit_forever') : (data.validity_period_month ? `${data.validity_period_month} ${validityUnitLabel(data.validity_period_unit)}` : '-') }}</span></template>
        </Column>
        <Column field="renewal_required" :header="t('training.renewal_required')" style="width:110px">
          <template #body="{data}"><Tag :value="data.renewal_required ? t('common.yes') : t('common.no')" :severity="data.renewal_required ? 'warning' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="is_active" :header="t('training.is_active')" style="width:90px">
          <template #body="{data}"><Tag :value="data.is_active ? t('common.yes') : t('common.no')" :severity="data.is_active ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
          <template #body="{data}">
            <div class="flex items-center gap-1 justify-end">
              <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openCertificationDialog(data)" />
              <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteCertification(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- ══ Tab 2: Issued certificates ══ -->
    <div v-if="activeTab === 'certificates'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2 flex-wrap">
          <SelectLabel v-model="participantFilter" :options="participantOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_participants')" class="!w-64" showClear filter @update:modelValue="onCertificateFilterChange" />
          <span v-if="certListTotal > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ certListTotal }} {{ t('common.items') }}</span>
        </div>
        <Button :label="t('training.certificate_generate')" icon="pi pi-plus" size="small" class="ml-auto" @click="openGenerateDialog()" />
      </div>

      <SkeletonTable v-if="certListLoading" :columns="certificateSkeletonColumns" :rows="8" />
      <DataTable
        v-else
        :value="issuedCertificates"
        lazy
        :totalRecords="certListTotal"
        :first="certListFirst"
        :rows="perPage"
        @page="onCertificatePage($event)"
        paginator
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        :rowsPerPageOptions="[10, 15, 25, 50]"
        size="small"
        class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      >
        <template #empty>
          <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
            <i class="pi pi-id-card text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('training.certificates_empty') }}</p>
          </div>
        </template>
        <Column field="certificate_no" :header="t('training.certificate_no')" style="width:160px">
          <template #body="{data}"><Tag :value="data.certificate_no" severity="success" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column field="participant_id" :header="t('training.employee')" style="width:200px">
          <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ participantName(data.participant_id) }}</span></template>
        </Column>
        <Column field="certification_id" :header="t('training.certification_name')">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ certificationName(data.certification_id) }}</span></template>
        </Column>
        <Column field="issued_date" :header="t('training.issued_date')" style="width:130px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.issued_date || '-' }}</span></template>
        </Column>
        <Column field="expiry_date" :header="t('training.expiry_date')" style="width:130px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.expiry_date || '-' }}</span></template>
        </Column>
        <Column field="certificate_file_url" :header="t('training.certificate_file_url')" style="width:120px">
          <template #body="{data}">
            <a v-if="data.certificate_file_url" :href="data.certificate_file_url" target="_blank" class="text-emerald-600 dark:text-emerald-400 hover:underline text-xs"><i class="pi pi-external-link mr-1"></i>{{ t('common.open') }}</a>
            <span v-else class="text-gray-400">-</span>
          </template>
        </Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
          <template #body="{data}">
            <div class="flex items-center gap-1 justify-end">
              <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('training.update_file')" @click="openFileDialog(data)" />
              <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteCertificate(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </div>

    <!-- ── Dialog: certification master ── -->
    <Dialog v-model:visible="certificationDialogVisible" :header="certificationEditing ? t('training.certification_edit') : t('training.certification_new')" modal :style="{ width: '520px' }" @hide="resetCertificationForm">
      <div class="space-y-4">
        <FormRow :label="t('training.validity_period')">
          <div class="flex items-center gap-2">
            <InputNumber v-model="certificationForm.validity_period_month" class="!w-full" :min="0" :useGrouping="false" :minFractionDigits="0" :maxFractionDigits="0" size="small" :disabled="certificationForm.validity_period_unit === 'forever'" />
            <SelectLabel v-model="certificationForm.validity_period_unit" :options="validityUnitOptions" optionLabel="label" optionValue="value" class="!w-36 shrink-0" @update:modelValue="onValidityUnitChange" />
          </div>
        </FormRow>
        <FormRow :label="t('training.certification_name')" required :errors="errors?.name">
          <TextInput v-model="certificationForm.name" maxlength="200" :placeholder="t('training.certification_name')" :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <FormRow :label="t('training.issuing_body')">
          <TextInput v-model="certificationForm.issuing_body" maxlength="200" :placeholder="t('training.issuing_body')" />
        </FormRow>
        <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
          <div>
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ t('training.renewal_required') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.certification_renewal_required_desc') }}</p>
          </div>
          <ToggleSwitch v-model="certificationForm.renewal_required" />
        </div>
        <div class="flex items-center justify-between gap-3 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2.5">
          <div>
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ t('training.is_active') }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('training.certification_is_active_desc') }}</p>
          </div>
          <ToggleSwitch v-model="certificationForm.is_active" />
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="certificationDialogVisible = false" />
          <Button :label="certificationEditing ? t('common.update') : t('common.save')" size="small" :loading="certificationSaving" :disabled="certificationSaving" @click="handleSaveCertification" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: generate certificate ── -->
    <Dialog v-model:visible="generateDialogVisible" :header="t('training.certificate_generate')" modal :style="{ width: '540px' }" @hide="resetGenerateForm">
      <div class="space-y-4">
        <p class="text-xs text-gray-500 dark:text-gray-400">{{ t('training.generate_hint') }}</p>
        <FormRow :label="t('training.employee')" required :errors="errors?.participant_id">
          <SelectLabel v-model="generateForm.participant_id" :options="completableParticipantOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.participant_id }" />
        </FormRow>
        <FormRow :label="t('training.certification_name')" required :errors="errors?.certification_id">
          <SelectLabel v-model="generateForm.certification_id" :options="activeCertificationOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.certification_id }" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.expiry_date')">
            <DateInput v-model="generateForm.expiry_date" />
          </FormRow>
          <FormRow :label="t('training.certificate_file_url')">
            <TextInput v-model="generateForm.certificate_file_url" maxlength="500" :placeholder="t('training.document_file_url_placeholder')" />
          </FormRow>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="generateDialogVisible = false" />
          <Button :label="t('training.certificate_generate')" size="small" :loading="generateSaving" :disabled="generateSaving" @click="handleGenerate" />
        </div>
      </template>
    </Dialog>

    <!-- ── Dialog: update certificate file URL ── -->
    <Dialog v-model:visible="fileDialogVisible" :header="t('training.update_file')" modal :style="{ width: '480px' }" @hide="resetFileForm">
      <div class="space-y-4">
        <FormRow :label="t('training.certificate_file_url')" :errors="errors?.certificate_file_url">
          <TextInput v-model="fileForm.certificate_file_url" maxlength="500" :placeholder="t('training.document_file_url_placeholder')" :class="{ 'p-invalid': errors?.certificate_file_url }" />
        </FormRow>
        <FormRow :label="t('training.expiry_date')">
          <DateInput v-model="fileForm.expiry_date" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="fileDialogVisible = false" />
          <Button :label="t('common.update')" size="small" :loading="fileSaving" :disabled="fileSaving" @click="handleUpdateFile" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteCertificationDialogVisible"
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_certification', { name: deleteCertificationTarget?.name || '' })"
      :loading="deletingCertification"
      :errorMsg="deleteCertificationError"
      @confirm="handleDeleteCertification"
    />
    <ConfirmDeleteDialog
      v-model:visible="deleteCertificateDialogVisible"
      :title="t('training.confirm_delete_title')"
      :message="t('training.confirm_delete_certificate', { no: deleteCertificateTarget?.certificate_no || '' })"
      :loading="deletingCertificate"
      :errorMsg="deleteCertificateError"
      @confirm="handleDeleteCertificate"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'

const { t } = useI18n()
const toast = useToast()

const tabs = [
  { key: 'certifications', labelKey: 'training.tab_certifications' },
  { key: 'certificates', labelKey: 'training.tab_certificates' }
]
const activeTab = ref('certifications')

// ── Certification master ──
const certifications = ref([])
const certLoading = ref(false)
const certTotal = ref(0)
const certPage = ref(1)
const perPage = ref(15)
const certificationDialogVisible = ref(false)
const certificationEditing = ref(false)
const certificationEditingId = ref(null)
const certificationSaving = ref(false)
const certificationForm = ref(defaultCertificationForm())

// ── Issued certificates ──
const issuedCertificates = ref([])
const certListLoading = ref(false)
const certListTotal = ref(0)
const certListPage = ref(1)
const participantFilter = ref(null)
const generateDialogVisible = ref(false)
const generateSaving = ref(false)
const generateForm = ref(defaultGenerateForm())
const fileDialogVisible = ref(false)
const fileSaving = ref(false)
const fileForm = ref({ certificate_file_url: '', expiry_date: null })
const fileTarget = ref(null)

const participants = ref([])
const employees = ref([])
const errors = ref({})

const deleteCertificationDialogVisible = ref(false)
const deletingCertification = ref(false)
const deleteCertificationError = ref('')
const deleteCertificationTarget = ref(null)
const deleteCertificateDialogVisible = ref(false)
const deletingCertificate = ref(false)
const deleteCertificateError = ref('')
const deleteCertificateTarget = ref(null)

const skeletonColumns = [
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]
const certificateSkeletonColumns = [
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const certFirst = computed(() => (certPage.value - 1) * perPage.value)
const certListFirst = computed(() => (certListPage.value - 1) * perPage.value)

const participantOptions = computed(() => participants.value.map(p => ({ label: participantName(p.id), value: p.id })))
// Catatan: `issuedForParticipant` hanya mengecek halaman sertifikat saat ini (lazy pagination),
// sehingga peserta yang sudah punya sertifikat di halaman lain tetap muncul — tidak fatal karena
// endpoint GenerateCertificate bersifat idempotent (update bila sudah ada).
const completableParticipantOptions = computed(() =>
  participants.value.filter(p => p.completion_status === 'COMPLETED' && !issuedForParticipant(p.id)).map(p => ({ label: participantName(p.id), value: p.id }))
)
const activeCertificationOptions = computed(() => certifications.value.filter(c => c.is_active).map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))

const validityUnitOptions = computed(() => [
  { label: t('training.validity_unit_year'), value: 'year' },
  { label: t('training.validity_unit_month'), value: 'month' },
  { label: t('training.validity_unit_forever'), value: 'forever' }
])
function validityUnitLabel(unit) {
  if (unit === 'forever') return t('training.validity_unit_forever')
  return unit === 'year' ? t('training.validity_unit_year') : t('training.validity_unit_month')
}
function onValidityUnitChange(unit) {
  // Berlaku selamanya → tidak perlu angka masa berlaku
  if (unit === 'forever') {
    certificationForm.value.validity_period_month = null
  }
}

function defaultCertificationForm() {
  return { name: '', issuing_body: '', validity_period_month: null, validity_period_unit: 'month', renewal_required: false, is_active: true }
}
function defaultGenerateForm() {
  return { participant_id: null, certification_id: null, certificate_file_url: '', expiry_date: null }
}

function issuedForParticipant(participantId) {
  return issuedCertificates.value.some(c => c.participant_id === participantId)
}

function participantName(participantId) {
  const p = participants.value.find(x => x.id === participantId)
  if (!p) return participantId
  const e = employees.value.find(x => x.id === p.employee_id)
  return e ? `${e.name} (${e.employee_id})` : (p.employee_id || participantId)
}
function certificationName(id) {
  const c = certifications.value.find(x => x.id === id)
  return c ? `${c.code} — ${c.name}` : (id || '-')
}

function switchTab(key) {
  activeTab.value = key
  if (key === 'certificates' && !issuedCertificates.value.length && !certListLoading.value) loadCertificates()
}

// ── Certification master CRUD ──
async function loadCertifications() {
  certLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/certifications', { params: { page: certPage.value, per_page: perPage.value } })
    const body = res.data
    certifications.value = body?.data || []
    certTotal.value = body?.total || 0
    if (body?.page) certPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    certLoading.value = false
  }
}

function onCertPage(event) {
  certPage.value = event.page + 1
  perPage.value = event.rows
  loadCertifications()
}

function openCertificationDialog(item) {
  errors.value = {}
  certificationEditing.value = !!item
  certificationEditingId.value = item?.id || null
  certificationForm.value = item
    ? {
        name: item.name || '',
        issuing_body: item.issuing_body || '',
        validity_period_month: item.validity_period_month ?? null,
        validity_period_unit: item.validity_period_unit || 'month',
        renewal_required: item.renewal_required,
        is_active: item.is_active
      }
    : defaultCertificationForm()
  certificationDialogVisible.value = true
}

function resetCertificationForm() {
  certificationForm.value = defaultCertificationForm()
  errors.value = {}
  certificationEditing.value = false
  certificationEditingId.value = null
}

async function handleSaveCertification() {
  errors.value = {}
  if (!certificationForm.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  certificationSaving.value = true
  try {
    const payload = {
      name: certificationForm.value.name.trim(),
      issuing_body: certificationForm.value.issuing_body?.trim() || '',
      validity_period_month: certificationForm.value.validity_period_month ?? null,
      validity_period_unit: certificationForm.value.validity_period_unit || 'month',
      renewal_required: certificationForm.value.renewal_required,
      is_active: certificationForm.value.is_active
    }
    if (certificationEditing.value) {
      await api.put(`/api/v1/tenant/trainings/certifications/${certificationEditingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/certifications', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    certificationDialogVisible.value = false
    await loadCertifications()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    certificationSaving.value = false
  }
}

function confirmDeleteCertification(item) {
  deleteCertificationTarget.value = item
  deleteCertificationError.value = ''
  deleteCertificationDialogVisible.value = true
}

async function handleDeleteCertification() {
  deletingCertification.value = true
  deleteCertificationError.value = ''
  try {
    await api.delete(`/api/v1/tenant/trainings/certifications/${deleteCertificationTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteCertificationDialogVisible.value = false
    await loadCertifications()
  } catch (e) {
    deleteCertificationError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deletingCertification.value = false
  }
}

// ── Issued certificates ──
async function loadCertificates() {
  certListLoading.value = true
  try {
    const params = { page: certListPage.value, per_page: perPage.value }
    if (participantFilter.value) params.participant_id = participantFilter.value
    const res = await api.get('/api/v1/tenant/trainings/certificates', { params })
    const body = res.data
    issuedCertificates.value = body?.data || []
    certListTotal.value = body?.total || 0
    if (body?.page) certListPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    certListLoading.value = false
  }
}

function onCertificatePage(event) {
  certListPage.value = event.page + 1
  perPage.value = event.rows
  loadCertificates()
}

function onCertificateFilterChange() {
  certListPage.value = 1
  loadCertificates()
}

async function loadParticipants() {
  try {
    const res = await api.get('/api/v1/tenant/trainings/participants', { params: { per_page: 500 } })
    participants.value = res.data?.data || []
  } catch {
    participants.value = []
  }
}

async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch {
    employees.value = []
  }
}

function openGenerateDialog() {
  errors.value = {}
  generateForm.value = defaultGenerateForm()
  generateDialogVisible.value = true
}

function resetGenerateForm() {
  generateForm.value = defaultGenerateForm()
  errors.value = {}
}

async function handleGenerate() {
  errors.value = {}
  if (!generateForm.value.participant_id) { errors.value = { participant_id: t('form.required') }; return }
  if (!generateForm.value.certification_id) { errors.value = { certification_id: t('form.required') }; return }
  generateSaving.value = true
  try {
    await api.post(`/api/v1/tenant/trainings/participants/${generateForm.value.participant_id}/certificate`, {
      certification_id: generateForm.value.certification_id,
      certificate_file_url: generateForm.value.certificate_file_url?.trim() || '',
      expiry_date: generateForm.value.expiry_date || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    generateDialogVisible.value = false
    await loadCertificates()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    generateSaving.value = false
  }
}

function openFileDialog(item) {
  errors.value = {}
  fileTarget.value = item
  fileForm.value = { certificate_file_url: item.certificate_file_url || '', expiry_date: item.expiry_date || null }
  fileDialogVisible.value = true
}

function resetFileForm() {
  fileForm.value = { certificate_file_url: '', expiry_date: null }
  fileTarget.value = null
  errors.value = {}
}

async function handleUpdateFile() {
  errors.value = {}
  fileSaving.value = true
  try {
    await api.put(`/api/v1/tenant/trainings/certificates/${fileTarget.value.id}`, {
      certificate_file_url: fileForm.value.certificate_file_url?.trim() || '',
      expiry_date: fileForm.value.expiry_date || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    fileDialogVisible.value = false
    await loadCertificates()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    fileSaving.value = false
  }
}

function confirmDeleteCertificate(item) {
  deleteCertificateTarget.value = item
  deleteCertificateError.value = ''
  deleteCertificateDialogVisible.value = true
}

async function handleDeleteCertificate() {
  deletingCertificate.value = true
  deleteCertificateError.value = ''
  try {
    await api.delete(`/api/v1/tenant/trainings/certificates/${deleteCertificateTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteCertificateDialogVisible.value = false
    await loadCertificates()
  } catch (e) {
    deleteCertificateError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deletingCertificate.value = false
  }
}

onMounted(() => {
  loadCertifications()
  loadParticipants()
  loadEmployees()
})
</script>
