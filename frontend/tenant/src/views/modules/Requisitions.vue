<template>
  <div class="space-y-4">
    <!-- Toolbar: filter status + tombol buat requisition -->
    <div class="flex items-center gap-2 flex-wrap">
      <SelectLabel
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('requisitions.filter_status')"
        class="!w-44"
        showClear
        @update:modelValue="onFilterChange()"
      />
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('requisitions.new_requisition')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-briefcase text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('requisitions.empty') }}</p>
        </div>
      </template>

      <Column :header="t('requisitions.title')" sortable field="title">
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-sky-50 dark:bg-sky-500/10 text-sky-600 dark:text-sky-400 flex items-center justify-center shrink-0">
              <i class="pi pi-briefcase text-xs"></i>
            </div>
            <div class="min-w-0">
              <p class="font-medium text-gray-800 dark:text-gray-100 truncate">{{ data.title }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ data.department || '—' }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('requisitions.slots')" style="width: 120px">
        <template #body="{ data }">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ data.slots_filled }}/{{ data.slots_available }}</span>
        </template>
      </Column>

      <Column :header="t('requisitions.reason_type')" style="width: 170px">
        <template #body="{ data }">
          <Tag
            :value="t('requisitions.reason_' + (data.reason_type || 'new_position').toLowerCase())"
            :severity="reasonSeverity(data.reason_type)"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>

      <Column :header="t('common.status')" style="width: 120px">
        <template #body="{ data }">
          <Tag
            :value="t('requisitions.status_' + (data.status || '').toLowerCase())"
            :severity="statusSeverity(data.status)"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>

      <Column :header="t('requisitions.target_start_date')" style="width: 130px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.target_start_date || '—' }}</span>
        </template>
      </Column>

      <!-- G-1: aksi — submit draft ke Central Approval -->
      <Column :header="t('common.actions')" :exportable="false" style="width: 100px">
        <template #body="{ data }">
          <Button
            v-if="data.status === 'DRAFT'"
            :label="t('requisitions.submit')"
            icon="pi pi-send"
            size="small"
            severity="info"
            outlined
            class="!text-xs !px-2.5 !py-1"
            @click="openSubmitDialog(data)"
          />
          <span v-else class="text-xs text-gray-400 dark:text-gray-500 italic">—</span>
        </template>
      </Column>
    </DataTable>

    <!-- G-1: dialog konfirmasi submit ke Central Approval -->
    <ConfirmActionDialog
      v-model:visible="submitDialogVisible"
      :title="t('requisitions.submit_confirm_title')"
      :message="t('requisitions.submit_confirm_message', { title: pendingSubmit?.title || '' })"
      :confirm-label="t('requisitions.submit')"
      :loading="submitting"
      :error-msg="submitError"
      icon="pi pi-send"
      severity="info"
      @confirm="submitRequisition()"
    />

    <!-- Dialog: buat requisition (S-1/S-5 — reason_type WORKFORCE_GAP / SUCCESSION_GAP) -->
    <Dialog v-model:visible="dialogVisible" :header="t('requisitions.new_requisition')" :modal="true" class="!w-[min(95vw,640px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('requisitions.org')" :required="true">
          <SelectLabel
            v-model="form.organization_id"
            :options="organizationOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
        </FormRow>
        <FormRow :label="t('requisitions.title')" :required="true">
          <TextInput v-model="form.title" :placeholder="t('requisitions.title_placeholder')" class="!w-full" />
        </FormRow>

        <FormRow :label="t('requisitions.department')">
          <TextInput v-model="form.department" :placeholder="t('requisitions.department_placeholder')" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.employment_type')">
          <TextInput v-model="form.employment_type" :placeholder="t('requisitions.employment_type_placeholder')" class="!w-full" />
        </FormRow>

        <FormRow :label="t('requisitions.location')">
          <TextInput v-model="form.location" :placeholder="t('requisitions.location_placeholder')" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.slots_available')">
          <InputNumber v-model="form.slots_available" :min="1" class="!w-full" />
        </FormRow>

        <FormRow :label="t('requisitions.reason_type')">
          <SelectLabel
            v-model="form.reason_type"
            :options="reasonOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
        </FormRow>
        <FormRow :label="t('requisitions.target_start_date')">
          <DateInput v-model="form.target_start_date" class="!w-full" />
        </FormRow>

        <!-- S-1: reason WORKFORCE_GAP → tautkan workforce gap -->
        <div v-if="form.reason_type === 'WORKFORCE_GAP'" class="md:col-span-2 rounded-lg border border-sky-200 dark:border-sky-700/40 bg-sky-50 dark:bg-sky-500/5 px-3 py-2.5">
          <p class="text-xs text-sky-700 dark:text-sky-300 flex items-center gap-1.5">
            <i class="pi pi-info-circle"></i>{{ t('requisitions.gap_hint') }}
          </p>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
            <FormRow :label="t('requisitions.workforce_gap_id')">
              <TextInput v-model="form.workforce_gap_id" :placeholder="t('requisitions.workforce_gap_placeholder')" class="!w-full" />
            </FormRow>
            <FormRow :label="t('requisitions.workforce_plan_id')">
              <TextInput v-model="form.workforce_plan_id" :placeholder="t('requisitions.workforce_plan_placeholder')" class="!w-full" />
            </FormRow>
          </div>
        </div>

        <!-- S-5: reason SUCCESSION_GAP → tautkan posisi kunci -->
        <div v-if="form.reason_type === 'SUCCESSION_GAP'" class="md:col-span-2 rounded-lg border border-emerald-200 dark:border-emerald-700/40 bg-emerald-50 dark:bg-emerald-500/5 px-3 py-2.5">
          <p class="text-xs text-emerald-700 dark:text-emerald-300 flex items-center gap-1.5">
            <i class="pi pi-info-circle"></i>{{ t('requisitions.succession_hint') }}
          </p>
          <FormRow :label="t('requisitions.succession_position_id')" class="!mt-2">
            <SelectLabel
              v-model="form.succession_position_id"
              :options="positionOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('common.select')"
              class="!w-full"
              showClear
            />
          </FormRow>
        </div>

        <FormRow :label="t('requisitions.min_salary')">
          <InputNumber v-model="form.min_salary" :min="0" mode="currency" currency="IDR" locale="id-ID" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.max_salary')">
          <InputNumber v-model="form.max_salary" :min="0" mode="currency" currency="IDR" locale="id-ID" class="!w-full" />
        </FormRow>

        <FormRow :label="t('requisitions.description_label')" class="md:col-span-2">
          <Textarea v-model="form.description" :rows="2" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.requirements')" class="md:col-span-2">
          <Textarea v-model="form.requirements" :rows="2" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.responsibilities')" class="md:col-span-2">
          <Textarea v-model="form.responsibilities" :rows="2" class="!w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" @click="save()" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import SelectLabel from '@/components/SelectLabel.vue'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const submitting = ref(false)
const submitError = ref('')
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const statusFilter = ref(null)
const dialogVisible = ref(false)
const submitDialogVisible = ref(false)
const pendingSubmit = ref(null)

const organizations = ref([])
const positions = ref([])

const skeletonColumns = [
  { field: 'title', header: 'Title', width: '40%' },
  { field: 'slots', header: 'Slots', width: '15%' },
  { field: 'reason', header: 'Reason', width: '20%' },
  { field: 'status', header: 'Status', width: '15%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = computed(() => ['DRAFT', 'SUBMITTED', 'OPEN', 'IN_PROGRESS', 'FILLED', 'REJECTED', 'CANCELLED'].map(v => ({ label: t(`requisitions.status_${v.toLowerCase()}`), value: v })))

const reasonOptions = computed(() => ['NEW_POSITION', 'REPLACEMENT', 'EXPANSION', 'WORKFORCE_GAP', 'SUCCESSION_GAP'].map(v => ({ label: t(`requisitions.reason_${v.toLowerCase()}`), value: v })))

const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.name, value: o.id })))
const positionOptions = computed(() => positions.value.map(p => ({ label: p.nomenclature || p.full_code, value: p.id })))

const emptyForm = () => ({
  organization_id: null,
  title: '',
  department: '',
  employment_type: '',
  location: '',
  min_salary: null,
  max_salary: null,
  slots_available: null,
  reason_type: null,
  workforce_gap_id: '',
  workforce_plan_id: '',
  succession_position_id: null,
  target_start_date: null,
  description: '',
  requirements: '',
  responsibilities: ''
})

const form = ref(emptyForm())

function reasonSeverity(reason) {
  switch (reason) {
    case 'WORKFORCE_GAP': return 'info'
    case 'SUCCESSION_GAP': return 'success'
    case 'REPLACEMENT': return 'warn'
    case 'EXPANSION': return 'help'
    default: return 'secondary'
  }
}

function statusSeverity(status) {
  switch (status) {
    case 'OPEN': return 'success'
    case 'IN_PROGRESS': return 'info'
    case 'FILLED': return 'help'
    case 'SUBMITTED': return 'info'
    case 'DRAFT': return 'secondary'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}

function openSubmitDialog(row) {
  pendingSubmit.value = row
  submitError.value = ''
  submitDialogVisible.value = true
}

async function submitRequisition() {
  if (!pendingSubmit.value) return
  submitting.value = true
  submitError.value = ''
  try {
    // flow_id kosong → backend auto-resolve flow aktif modul recruitment (G-1)
    await api.post(`/api/v1/tenant/recruitment/requisitions/${pendingSubmit.value.id}/submit`, {})
    submitDialogVisible.value = false
    pendingSubmit.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('requisitions.submitted'), life: 3000 })
    loadData()
  } catch (e) {
    submitError.value = getErrorMessage(e, t('message.failed_to_save'))
  } finally {
    submitting.value = false
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/recruitment/requisitions', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadOptions() {
  try {
    const [orgRes, posRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/organizations', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/positions', { params: { per_page: 500 } })
    ])
    organizations.value = orgRes.status === 'fulfilled' ? (orgRes.value.data?.data || []) : []
    positions.value = posRes.status === 'fulfilled' ? (posRes.value.data?.data || []) : []
  } catch {
    // fail-silent — dropdown kosong, requisition tetap bisa dibuat
  }
}

function openDialog() {
  form.value = emptyForm()
  dialogVisible.value = true
}

async function save() {
  // Validasi client-side — field wajib backend (organization_id & title)
  if (!form.value.organization_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('requisitions.org_required'), life: 4000 })
    return
  }
  if (!form.value.title.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('requisitions.title_required'), life: 4000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    // Slots kosong → backend auto-resolve dari hiring need WI (S-1)
    if (payload.slots_available === null) delete payload.slots_available
    // Hapus field yang tidak relevan dengan reason_type
    if (payload.reason_type !== 'WORKFORCE_GAP') {
      delete payload.workforce_gap_id
      delete payload.workforce_plan_id
    }
    if (payload.reason_type !== 'SUCCESSION_GAP') {
      delete payload.succession_position_id
    }
    // Field kosong → jangan dikirim
    Object.keys(payload).forEach(k => {
      if (payload[k] === '' || payload[k] === null || payload[k] === undefined) delete payload[k]
    })
    await api.post('/api/v1/tenant/recruitment/requisitions', payload)
    dialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('requisitions.created'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    saving.value = false
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

onMounted(() => {
  loadData()
  loadOptions()
})
</script>
