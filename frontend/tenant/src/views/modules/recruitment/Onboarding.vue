<template>
  <div class="space-y-4">
    <!-- Toolbar: filter status + tombol buat onboarding -->
    <div class="flex items-center gap-2 flex-wrap">
      <SelectLabel
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('onboarding.filter_status')"
        class="!w-48"
        showClear
        @update:modelValue="onFilterChange()"
      />
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('onboarding.new_onboarding')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-rocket text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('onboarding.empty') }}</p>
        </div>
      </template>

      <!-- Employee (hasil hire — G-4: recruited_from_application_id) -->
      <Column :header="t('onboarding.employee')" style="min-width: 220px">
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-teal-50 dark:bg-teal-500/10 text-teal-600 dark:text-teal-400 flex items-center justify-center shrink-0">
              <i class="pi pi-user text-xs"></i>
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-1.5">
                <p class="font-medium text-gray-800 dark:text-gray-100 truncate">{{ employeeName(data.employee_id) }}</p>
                <!-- G-4: badge employee yang dibuat dari offer recruitment -->
                <Tag
                  v-if="isHiredFromOffer(data.employee_id)"
                  :value="t('onboarding.from_offer')"
                  severity="success"
                  class="!text-[10px] !px-1.5 !py-0"
                />
              </div>
              <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ employeeCode(data.employee_id) }}</p>
            </div>
          </div>
        </template>
      </Column>

      <!-- Kandidat (via application → candidate) -->
      <Column :header="t('onboarding.candidate')" style="min-width: 180px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ candidateName(data.application_id) }}</span>
        </template>
      </Column>

      <Column :header="t('onboarding.start_date')" style="width: 130px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.start_date || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('common.status')" style="width: 130px">
        <template #body="{ data }">
          <Tag
            :value="t('onboarding.status_' + (data.status || '').toLowerCase())"
            :severity="statusSeverity(data.status)"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>

      <!-- Aksi: mulai / selesaikan onboarding (COMPLETED → training handoff S-7) -->
      <Column :header="t('common.actions')" :exportable="false" style="width: 200px">
        <template #body="{ data }">
          <div class="flex items-center gap-1.5 flex-wrap">
            <template v-if="data.status === 'PENDING'">
              <Button :label="t('onboarding.start')" icon="pi pi-play" size="small" severity="info" outlined class="!text-xs !px-2.5 !py-1" @click="setStatus(data, 'IN_PROGRESS')" />
            </template>
            <template v-else-if="data.status === 'IN_PROGRESS'">
              <Button :label="t('onboarding.complete')" icon="pi pi-check" size="small" severity="success" class="!text-xs !px-2.5 !py-1" @click="openCompleteDialog(data)" />
            </template>
            <span v-else class="text-xs text-gray-400 dark:text-gray-500 italic">—</span>
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Dialog konfirmasi selesaikan onboarding (COMPLETED → handoff training S-7) -->
    <ConfirmActionDialog
      v-model:visible="completeDialogVisible"
      :title="t('onboarding.complete_confirm_title')"
      :message="t('onboarding.complete_confirm_message', { name: employeeName(pendingComplete?.employee_id) })"
      :confirm-label="t('onboarding.complete')"
      :loading="completing"
      :error-msg="completeError"
      icon="pi pi-check"
      severity="success"
      @confirm="completeOnboarding()"
    />

    <!-- Dialog: buat onboarding untuk employee hasil offer / kandidat diterima -->
    <Dialog v-model:visible="dialogVisible" :header="t('onboarding.new_onboarding')" :modal="true" class="!w-[min(95vw,560px)]">
      <div class="grid grid-cols-1 gap-3">
        <!-- Pilih aplikasi ACCEPTED (hasil hire) — auto-suggest employee dari offer -->
        <FormRow :label="t('onboarding.application')" :required="true">
          <SelectLabel
            v-model="form.application_id"
            :options="applicationOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('onboarding.application_placeholder')"
            class="!w-full"
            showClear
            @update:modelValue="onApplicationChange()"
          />
        </FormRow>

        <FormRow :label="t('onboarding.employee')" :required="true">
          <SelectLabel
            v-model="form.employee_id"
            :options="employeeOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
          <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-1">{{ t('onboarding.employee_hint') }}</p>
        </FormRow>

        <FormRow :label="t('onboarding.start_date')" :required="true">
          <DateInput v-model="form.start_date" class="!w-full" />
        </FormRow>

        <FormRow :label="t('onboarding.notes')">
          <Textarea v-model="form.notes" :rows="2" :placeholder="t('onboarding.notes_placeholder')" class="!w-full" />
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
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'
import Textarea from 'primevue/textarea'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const completing = ref(false)
const completeError = ref('')
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const statusFilter = ref(null)
const dialogVisible = ref(false)
const completeDialogVisible = ref(false)
const pendingComplete = ref(null)

// Enrichment: employee (dengan recruited_from_application_id G-4) + application + candidate
const employees = ref([])
const applications = ref([])
const candidates = ref([])
const onboardings = ref([])

const skeletonColumns = [
  { field: 'employee', header: 'Employee', width: '26%' },
  { field: 'candidate', header: 'Candidate', width: '20%' },
  { field: 'start_date', header: 'Start', width: '14%' },
  { field: 'status', header: 'Status', width: '13%' },
  { field: 'actions', header: 'Actions', width: '27%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = computed(() =>
  ['PENDING', 'IN_PROGRESS', 'COMPLETED'].map(v => ({ label: t(`onboarding.status_${v.toLowerCase()}`), value: v }))
)

const employeeById = computed(() => {
  const m = {}
  for (const e of employees.value) m[e.id] = e
  return m
})

const applicationById = computed(() => {
  const m = {}
  for (const a of applications.value) m[a.id] = a
  return m
})

const candidateById = computed(() => {
  const m = {}
  for (const c of candidates.value) m[c.id] = c
  return m
})

// Dropdown aplikasi: hanya yang ACCEPTED (hasil hire) & belum punya onboarding
// (exclusion set dari daftar onboardings — cegah duplikat untuk hire yang sama).
const onboardedApplicationIds = computed(() => new Set(onboardings.value.map(o => o.application_id)))

const applicationOptions = computed(() =>
  applications.value
    .filter(a => a.status === 'ACCEPTED' && !onboardedApplicationIds.value.has(a.id))
    .map(a => ({
      label: `${candidateName(a.id)} — ${t('onboarding.application_accepted')}`,
      value: a.id
    }))
)

// Dropdown employee: label nama + kode; employee hasil offer ditandai
const employeeOptions = computed(() =>
  employees.value.map(e => ({
    label: isHiredFromOffer(e.id)
      ? `${e.name} ★ ${t('onboarding.from_offer')}`
      : e.name,
    value: e.id
  }))
)

const emptyForm = () => ({
  application_id: null,
  employee_id: null,
  start_date: null,
  notes: ''
})

const form = ref(emptyForm())

function statusSeverity(status) {
  switch (status) {
    case 'COMPLETED': return 'success'
    case 'IN_PROGRESS': return 'info'
    case 'PENDING': return 'warn'
    default: return 'secondary'
  }
}

// G-4: employee yang dibuat dari offer → recruited_from_application_id terisi
function isHiredFromOffer(employeeId) {
  const e = employeeById.value[employeeId]
  return !!e && !!e.recruited_from_application_id
}

function employeeName(employeeId) {
  return employeeById.value[employeeId]?.name || '—'
}

function employeeCode(employeeId) {
  return employeeById.value[employeeId]?.employee_id || ''
}

function candidateName(applicationId) {
  const app = applicationById.value[applicationId]
  if (!app) return '—'
  const c = candidateById.value[app.candidate_id]
  return c ? `${c.first_name} ${c.last_name || ''}`.trim() : '—'
}

// Auto-suggest: aplikasi dipilih → employee yang dibuat dari offer aplikasi itu
function onApplicationChange() {
  if (!form.value.application_id) return
  const matched = employees.value.find(e => e.recruited_from_application_id === form.value.application_id)
  if (matched) {
    form.value.employee_id = matched.id
    toast.add({ severity: 'info', summary: t('message.info'), detail: t('onboarding.auto_suggest'), life: 3000 })
  }
}

async function setStatus(row, status) {
  try {
    await api.put(`/api/v1/tenant/recruitment/employee-onboardings/${row.id}`, { status })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('onboarding.status_updated'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  }
}

function openCompleteDialog(row) {
  pendingComplete.value = row
  completeError.value = ''
  completeDialogVisible.value = true
}

async function completeOnboarding() {
  if (!pendingComplete.value) return
  completing.value = true
  completeError.value = ''
  try {
    // COMPLETED → backend memicu training handoff (S-7) untuk employee tsb.
    await api.put(`/api/v1/tenant/recruitment/employee-onboardings/${pendingComplete.value.id}`, { status: 'COMPLETED' })
    completeDialogVisible.value = false
    pendingComplete.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('onboarding.completed_done'), life: 3000 })
    loadData()
  } catch (e) {
    completeError.value = getErrorMessage(e, t('message.failed_to_save'))
  } finally {
    completing.value = false
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/recruitment/employee-onboardings', { params })
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

async function loadEnrichment() {
  try {
    const [empRes, appRes, candRes, onbRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/applications', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/candidates', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/employee-onboardings', { params: { per_page: 500 } })
    ])
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
    applications.value = appRes.status === 'fulfilled' ? (appRes.value.data?.data || []) : []
    candidates.value = candRes.status === 'fulfilled' ? (candRes.value.data?.data || []) : []
    onboardings.value = onbRes.status === 'fulfilled' ? (onbRes.value.data?.data || []) : []
  } catch {
    // fail-silent — kolom enrichment jadi '—', onboarding tetap bisa dikelola
  }
}

function openDialog() {
  form.value = emptyForm()
  dialogVisible.value = true
}

async function save() {
  if (!form.value.application_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('onboarding.application_required'), life: 4000 })
    return
  }
  if (!form.value.employee_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('onboarding.employee_required'), life: 4000 })
    return
  }
  if (!form.value.start_date) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('onboarding.start_date_required'), life: 4000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    Object.keys(payload).forEach(k => {
      if (payload[k] === '' || payload[k] === null || payload[k] === undefined) delete payload[k]
    })
    await api.post('/api/v1/tenant/recruitment/employee-onboardings', payload)
    dialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('onboarding.created'), life: 3000 })
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
  loadEnrichment()
})
</script>
