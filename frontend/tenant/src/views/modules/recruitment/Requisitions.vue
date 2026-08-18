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
              <p class="font-medium text-navy-800 dark:text-gray-100 truncate">{{ data.title }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ data.department || '—' }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('requisitions.requisition_number')" style="width: 160px">
        <template #body="{ data }">
          <span class="text-xs font-mono text-gray-600 dark:text-gray-300">{{ data.requisition_number || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('requisitions.priority')" style="width: 110px">
        <template #body="{ data }">
          <Tag
            :value="t('requisitions.priority_' + (data.priority || 'medium').toLowerCase())"
            :severity="prioritySeverity(data.priority)"
            class="!text-xs !px-1.5 !py-0.5"
          />
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
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ formatDate(data.target_start_date, locale) || '—' }}</span>
        </template>
      </Column>

      <!-- G-2: opened_at — kapan requisition dibuka (approval/manual) -->
      <Column :header="t('requisitions.opened_at')" style="width: 150px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ formatOpenedAt(data.opened_at) }}</span>
        </template>
      </Column>

      <!-- G-9 sub-1: requirements/competencies — kolom sendiri, bukan gabung ke Actions -->
      <Column :header="t('requisitions.requirements_competencies')" :exportable="false" style="width: 130px">
        <template #body="{ data }">
          <Button
            :label="t('requisitions.requirements_tab')"
            icon="pi pi-list-check"
            size="small"
            text
            severity="secondary"
            class="!text-xs"
            @click="router.push(`/recruitment/requisitions/${data.id}/requirements`)"
          />
        </template>
      </Column>

      <!-- G-1: aksi — draft bisa diubah/dihapus + submit ke Central Approval -->
      <Column :header="t('common.actions')" :exportable="false" style="width: 170px">
        <template #body="{ data }">
          <div v-if="data.status === 'DRAFT'" class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
            <Button
              :label="t('requisitions.submit')"
              icon="pi pi-send"
              size="small"
              severity="info"
              outlined
              class="!text-xs !px-2.5 !py-1"
              @click="openSubmitDialog(data)"
            />
          </div>
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

    <!-- Dialog: buat/ubah requisition (S-1/S-5 — reason_type WORKFORCE_GAP / SUCCESSION_GAP) -->
    <Dialog v-model:visible="dialogVisible" :header="editingId ? t('requisitions.edit_requisition') : t('requisitions.new_requisition')" :modal="true" class="!w-[min(95vw,960px)]">
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
          <InputNumber v-model="form.slots_available" :min="1" size="small" class="!w-full" />
        </FormRow>

        <FormRow :label="t('requisitions.requisition_number')">
          <TextInput v-model="form.requisition_number" :placeholder="t('requisitions.requisition_number_placeholder')" class="!w-full" />
          <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-1">{{ t('requisitions.requisition_number_hint') }}</p>
        </FormRow>
        <FormRow :label="t('requisitions.priority')">
          <SelectLabel
            v-model="form.priority"
            :options="priorityOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
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
          <InputNumber v-model="form.min_salary" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" class="!w-full" />
        </FormRow>
        <FormRow :label="t('requisitions.max_salary')">
          <InputNumber v-model="form.max_salary" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" class="!w-full" />
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
          <Button :label="editingId ? t('common.update') : t('common.save')" icon="pi pi-check" size="small" :loading="saving" @click="save()" />
        </div>
      </template>
    </Dialog>

    <!-- Konfirmasi hapus requisition draft -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('requisitions.confirm_delete_title')"
      :message="t('requisitions.confirm_delete_message', { title: deleteTarget?.title || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      @confirm="handleDelete()"
    />

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import { formatDate } from '@/utils/formatDate'
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
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const router = useRouter()
const { t, locale } = useI18n()
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
const editingId = ref(null)
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const submitDialogVisible = ref(false)
const pendingSubmit = ref(null)

const organizations = ref([])
const positions = ref([])

const skeletonColumns = [
  { field: 'title', header: 'Title', width: '25%' },
  { field: 'requisition_number', header: 'Requisition No.', width: '11%' },
  { field: 'priority', header: 'Priority', width: '9%' },
  { field: 'slots', header: 'Slots', width: '9%' },
  { field: 'reason', header: 'Reason', width: '12%' },
  { field: 'status', header: 'Status', width: '9%' },
  { field: 'opened_at', header: 'Opened At', width: '12%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = computed(() => ['DRAFT', 'SUBMITTED', 'OPEN', 'IN_PROGRESS', 'FILLED', 'REJECTED', 'CANCELLED'].map(v => ({ label: t(`requisitions.status_${v.toLowerCase()}`), value: v })))

const reasonOptions = computed(() => ['NEW_POSITION', 'REPLACEMENT', 'EXPANSION', 'WORKFORCE_GAP', 'SUCCESSION_GAP'].map(v => ({ label: t(`requisitions.reason_${v.toLowerCase()}`), value: v })))

const priorityOptions = computed(() => ['LOW', 'MEDIUM', 'HIGH', 'URGENT'].map(v => ({ label: t(`requisitions.priority_${v.toLowerCase()}`), value: v })))

const organizationOptions = computed(() => organizations.value.map(o => ({ label: o.nomenclature || o.full_code || o.code, value: o.id })))
const positionOptions = computed(() => positions.value.map(p => ({ label: p.nomenclature || p.full_code, value: p.id })))

const emptyForm = () => ({
  organization_id: null,
  title: '',
  requisition_number: '',
  priority: null,
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

function prioritySeverity(priority) {
  switch (priority) {
    case 'URGENT': return 'danger'
    case 'HIGH': return 'warn'
    case 'MEDIUM': return 'info'
    case 'LOW': return 'secondary'
    default: return 'secondary'
  }
}

function formatOpenedAt(value) {
  // opened_at dikirim backend sebagai unix NANO detik (time.Now().UnixNano())
  // → bagi 1.000.000 dulu sebelum new Date() (menerima milidetik).
  if (!value) return '—'
  const ms = Number(value) / 1000000
  if (!Number.isFinite(ms) || ms <= 0) return '—'
  return formatDate(new Date(ms), locale.value) || '—'
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

function openDialog(item) {
  editingId.value = item?.id || null
  form.value = item
    ? {
        organization_id: item.organization_id || null,
        title: item.title || '',
        requisition_number: item.requisition_number || '',
        priority: item.priority || null,
        department: item.department || '',
        employment_type: item.employment_type || '',
        location: item.location || '',
        min_salary: item.min_salary ?? null,
        max_salary: item.max_salary ?? null,
        slots_available: item.slots_available ?? null,
        reason_type: item.reason_type || null,
        workforce_gap_id: item.workforce_gap_id || '',
        workforce_plan_id: item.workforce_plan_id || '',
        succession_position_id: item.succession_position_id || null,
        target_start_date: item.target_start_date || null,
        description: item.description || '',
        requirements: item.requirements || '',
        responsibilities: item.responsibilities || ''
      }
    : emptyForm()
  dialogVisible.value = true
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
    await api.delete(`/api/v1/tenant/recruitment/requisitions/${deleteTarget.value.id}`)
    deleteDialogVisible.value = false
    deleteTarget.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
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
    if (editingId.value) {
      await api.put(`/api/v1/tenant/recruitment/requisitions/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('requisitions.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/recruitment/requisitions', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('requisitions.created'), life: 3000 })
    }
    dialogVisible.value = false
    editingId.value = null
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
