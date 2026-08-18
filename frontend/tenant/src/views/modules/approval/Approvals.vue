<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="pendingTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ pendingTotal }} {{ t('common.items') }}
        </span>
        <Button icon="pi pi-refresh" size="small" text severity="secondary" @click="loadTasks" />
      </div>
      <Button v-if="hasExactPermission('approval.settings.view')" :label="t('approval.flows')" icon="pi pi-sitemap" size="small" @click="router.push({ name: 'ApprovalFlows' })" />
    </div>

    <!-- Tab: Perlu Tindakan / Riwayat -->
    <div class="flex items-center gap-1 mb-2 border-b border-gray-200 dark:border-gray-700">
      <button
        class="px-3 py-1.5 text-sm font-medium border-b-2 -mb-px transition-colors"
        :class="activeTab === 'pending' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab('pending')"
      >
        {{ t('approval.tab_pending') }}
      </button>
      <button
        class="px-3 py-1.5 text-sm font-medium border-b-2 -mb-px transition-colors"
        :class="activeTab === 'done' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab('done')"
      >
        {{ t('approval.tab_history') }}
      </button>
    </div>

    <!-- Filter status & alur -->
    <div class="flex items-center gap-2 flex-wrap mb-2">
      <Select
        v-if="activeTab === 'pending'"
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('approval.filter_status')"
        class="!w-44"
        size="small"
        showClear
        @update:modelValue="onFilterChange"
      />
      <Select
        v-model="flowFilter"
        :options="flowOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('approval.filter_flow')"
        class="!w-56"
        size="small"
        showClear
        :loading="flowsLoading"
        @update:modelValue="onFilterChange"
      />
      <Button
        v-if="statusFilter || flowFilter"
        :label="t('approval.filter_reset')"
        icon="pi pi-filter-slash"
        size="small"
        text
        severity="secondary"
        @click="clearFilters"
      />
    </div>

    <SkeletonTable v-if="tasksLoading" :columns="taskSkeletonColumns" :rows="6" />

    <DataTable
      v-else
      :value="pendingTasks"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      lazy
      :totalRecords="pendingTotal"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 20, 50]"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t(activeTab === 'pending' ? 'approval.empty_tasks' : 'approval.empty_history') }}</p>
        </div>
      </template>
      <Column field="step_order" :header="t('approval.step')" style="width:80px">
        <template #body="{data}">#{{ data.step_order }}</template>
      </Column>
      <Column field="flow_name" :header="t('approval.instance')">
        <template #body="{data}">
          <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.flow_name || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_name" :header="t('approval.submitted_by')">
        <template #body="{data}">
          <span class="text-gray-700 dark:text-gray-200">{{ data.submitter_name || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_employee_code" :header="t('approval.employee_id')" style="width:120px">
        <template #body="{data}">
          <span class="text-xs text-gray-500 dark:text-gray-400 font-mono">{{ data.submitter_employee_code || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_organization_name" :header="t('approval.organization')" style="width:160px">
        <template #body="{data}">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.submitter_organization_name || '-' }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:120px">
        <template #body="{data}">
          <Tag :value="rowStatusLabel(data)" :severity="rowStatusSeverity(data)" class="!text-xs" />
        </template>
      </Column>
      <Column field="created_at" :header="t('approval.submitted_at')" style="width:160px">
        <template #body="{data}">
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(data.created_at) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button :label="t('approval.review')" icon="pi pi-eye" size="small" text @click="openTaskDetail(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- ===================================================================== -->
    <!-- Task detail / act dialog -->
    <!-- ===================================================================== -->
    <Dialog v-model:visible="taskDetailVisible" :header="t('approval.instance_detail')" modal :style="{ width: '960px' }">
      <div v-if="instanceLoading" class="space-y-5 animate-pulse">
        <!-- Skeleton: submitted data section -->
        <div class="space-y-4 pb-5">
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-gray-200 dark:bg-gray-600 rounded"></div>
            <div class="h-3 w-32 bg-gray-200 dark:bg-gray-600 rounded"></div>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>
          <div class="space-y-3">
            <div v-for="i in 4" :key="i" class="space-y-1">
              <div class="h-3 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-4 w-48 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
          </div>
        </div>
        <!-- Skeleton: approval data section -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <div class="w-4 h-4 bg-gray-200 dark:bg-gray-600 rounded"></div>
            <div class="h-3 w-28 bg-gray-200 dark:bg-gray-600 rounded"></div>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <div class="h-3 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-5 w-20 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
            <div class="space-y-1">
              <div class="h-3 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-5 w-24 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
            <div class="col-span-2 space-y-1">
              <div class="h-3 w-20 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-4 w-40 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="h-3 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
            <div v-for="i in 3" :key="i" class="flex items-center gap-2">
              <div class="w-5 h-5 bg-gray-200 dark:bg-gray-600 rounded-full"></div>
              <div class="h-3 w-32 bg-gray-200 dark:bg-gray-600 rounded"></div>
              <div class="h-4 w-16 bg-gray-200 dark:bg-gray-600 rounded"></div>
            </div>
          </div>
        </div>
      </div>
      <div v-else-if="activeInstance" class="space-y-5">
        <!-- Section: the data being submitted -->
        <div class="space-y-4 dark:border-gray-700 pb-5">
          <div v-if="!isAttendanceModule" class="flex items-center gap-2">
            <i class="pi pi-file text-teal-400 text-sm"></i>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('approval.submitted_data') }}</h2>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>

          <SkeletonCard v-if="documentLoading" type="detail" :count="1" :rows="4" cols="grid-cols-1" padding="p-0" />

          <KpiDetail v-else-if="isKPIModule" :document-detail="documentDetail" />
          <OkrDetail v-else-if="isOKRModule" :document-detail="documentDetail" />
          <AttendanceDetail v-else-if="isAttendanceModule" :document-detail="documentDetail" :submitter-name="activeTaskRef?.submitter_name" />
          <LeaveDetail v-else-if="isLeaveModule" :document-detail="documentDetail" />
          <MovementDetail v-else-if="isEmployeemovementModule" :document-detail="documentDetail" />
          <RecruitmentDetail v-else-if="isRecruitmentModule" :document-detail="documentDetail" />
          <OfferDetail v-else-if="isOfferModule" :document-detail="documentDetail" />
          <TrainingRequestDetail v-else-if="isTrainingRequestModule" :document-detail="documentDetail" />
          <BusinessTravelDetail v-else-if="isBusinessTravelModule" :document-detail="documentDetail" />

          <div v-else-if="documentFields.length" class="grid grid-cols-2 gap-3">
            <ViewLabel v-for="f in documentFields" :key="f.label" :label="f.label" :value="f.value" break-all />
          </div>

          <p v-else class="text-xs text-gray-400">{{ t('approval.submitted_data_unavailable') }}</p>
        </div>

        <!-- Right column: existing approval data -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <i class="pi pi-check-circle text-teal-400 text-sm"></i>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('approval.approval_data') }}</h2>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <ViewLabel :label="t('approval.module')">
              <Tag :value="moduleLabel(activeInstance.module)" severity="info" class="!text-xs" />
            </ViewLabel>
            <ViewLabel :label="t('common.status')">
              <Tag :value="activeInstance.status" severity="warn" class="!text-xs" />
            </ViewLabel>
            <ViewLabel :label="t('approval.flow_name')" :value="activeInstance.flow_name" class="col-span-2" />
          </div>

          <div>
            <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-2">{{ t('approval.steps') }}</p>
            <div class="space-y-1">
              <div v-for="s in activeInstance.steps" :key="s.id" class="px-2 py-1.5 rounded" :class="s.step_order === activeInstance.current_step ? 'bg-emerald-50 dark:bg-emerald-500/10' : ''">
                <div class="flex items-center gap-2 text-xs">
                  <span class="w-5 h-5 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">{{ s.step_order }}</span>
                  <span class="font-medium text-gray-700 dark:text-gray-200">{{ s.step_name }}</span>
                  <Tag :value="s.participation_type" :severity="s.participation_type === 'WATCHER' ? 'secondary' : 'success'" class="!text-xs !px-1.5 !py-0.5" />
                </div>
                <div v-if="stepApprover(s.step_order)" class="ml-7 mt-1 text-xs text-gray-500 dark:text-gray-400">
                  <i class="pi pi-check-circle text-emerald-500 mr-1"></i>
                  {{ stepApprover(s.step_order).actor_name || '-' }}
                  <span v-if="stepApprover(s.step_order).actor_employee_code">({{ stepApprover(s.step_order).actor_employee_code }})</span>
                  <span v-if="stepApprover(s.step_order).actor_organization_name"> — {{ stepApprover(s.step_order).actor_organization_name }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeInstance.actions?.length">
            <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-2">{{ t('approval.history') }}</p>
            <div class="space-y-1">
              <div v-for="a in activeInstance.actions" :key="a.id" class="text-xs text-gray-600 dark:text-gray-300 flex items-center gap-2">
                <Tag :value="a.action" :severity="a.action === 'APPROVE' ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
                <span>{{ formatDate(a.created_at) }}</span>
                <span v-if="a.note" class="text-gray-400">— {{ a.note }}</span>
              </div>
            </div>
          </div>

          <div v-if="activeTaskIsWatcher" class="space-y-2">
            <p class="text-xs text-gray-400 flex items-center gap-1.5">
              <i class="pi pi-eye"></i> {{ t('approval.watcher_note') }}
            </p>
          </div>
          <div v-else-if="isInstanceResolved" class="space-y-2">
            <p class="text-xs text-gray-400 flex items-center gap-1.5">
              <i class="pi pi-check-circle text-emerald-500"></i> {{ t('approval.instance_resolved') }}
            </p>
          </div>
          <div v-else class="space-y-2">
            <Textarea v-model="actionNote" :placeholder="t('approval.note_placeholder')" rows="2" class="w-full" :class="{ 'p-invalid': noteError }" @input="noteError = ''" />
            <small v-if="noteError" class="text-red-500 text-xs block">{{ noteError }}</small>
          </div>
        </div>
      </div>

      <template #footer>
        <div v-if="!instanceLoading && activeInstance" class="flex items-center justify-end gap-2">
          <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="taskDetailVisible = false" />
          <template v-if="!activeTaskIsWatcher && !isInstanceResolved">
            <Button :label="t('approval.reject')" severity="danger" outlined size="small" :loading="actionSubmitting" @click="submitAction('REJECT')" />
            <Button :label="t('approval.approve')" severity="success" size="small" :loading="actionSubmitting" @click="submitAction('APPROVE')" />
          </template>
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { formatDate as formatDateGlobal } from '@/utils/formatDate'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ViewLabel from '@/components/ViewLabel.vue'
import KpiDetail from './detail/KpiDetail.vue'
import OkrDetail from './detail/OkrDetail.vue'
import AttendanceDetail from './detail/AttendanceDetail.vue'
import LeaveDetail from './detail/LeaveDetail.vue'
import MovementDetail from './detail/MovementDetail.vue'
import RecruitmentDetail from './detail/RecruitmentDetail.vue'
import OfferDetail from './detail/OfferDetail.vue'
import TrainingRequestDetail from './detail/TrainingRequestDetail.vue'
import BusinessTravelDetail from './detail/BusinessTravelDetail.vue'
import { useAuth } from '@/stores/auth'

const { t, locale } = useI18n()
const toast = useToast()
const router = useRouter()
const { hasExactPermission } = useAuth()

// formatDate — date + time, built on the app-wide date formatter
// (utils/formatDate.js) so the date portion matches every other page
// (e.g. "30 July 2026"), with a time suffix appended since the global
// utility is date-only.
function formatDate(v) {
  if (!v) return '-'
  const datePart = formatDateGlobal(v, locale.value)
  if (!datePart) return '-'
  const time = new Date(v).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

function moduleLabel(slug) {
  const label = t(`approval.module_names.${slug}`)
  return label !== `approval.module_names.${slug}` ? label : slug
}

const pendingTasks = ref([])
const pendingTotal = ref(0)
const tasksLoading = ref(false)

// ── Tab, filter & pagination ──
const activeTab = ref('pending') // 'pending' | 'done'
const statusFilter = ref(null)
const flowFilter = ref(null)
const flowOptions = ref([])
const flowsLoading = ref(false)
const currentPage = ref(1)
const perPage = ref(20)

function switchTab(tab) {
  if (activeTab.value === tab) return
  activeTab.value = tab
  statusFilter.value = null // filter status hanya untuk tab pending
  currentPage.value = 1
  loadTasks()
}

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// Filter status = status instance approval (task pending selalu PENDING,
// sehingga yang bermakna adalah status keseluruhan instance-nya).
const statusOptions = computed(() => [
  { label: t('approval.instance_status_pending'), value: 'PENDING' },
  { label: t('approval.instance_status_approved'), value: 'APPROVED' },
  { label: t('approval.instance_status_rejected'), value: 'REJECTED' },
  { label: t('approval.instance_status_cancelled'), value: 'CANCELLED' }
])

async function loadFlowOptions() {
  flowsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/approval/flows', { params: { page: 1, per_page: 100 } })
    const flows = res.data?.data || []
    flowOptions.value = flows.map(f => ({ label: f.name, value: f.id }))
  } catch {
    flowOptions.value = []
  } finally {
    flowsLoading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadTasks()
}

function onFilterChange() {
  currentPage.value = 1
  loadTasks()
}

function clearFilters() {
  statusFilter.value = null
  flowFilter.value = null
  onFilterChange()
}

const taskSkeletonColumns = [
  { type: 'text', width: 'w-10', headerWidth: 'w-12' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

async function loadTasks() {
  tasksLoading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: perPage.value,
      ...(activeTab.value === 'pending' && statusFilter.value ? { status: statusFilter.value } : {}),
      ...(flowFilter.value ? { flow_id: flowFilter.value } : {})
    }
    const endpoint = activeTab.value === 'done' ? '/api/v1/tenant/approval/tasks/done' : '/api/v1/tenant/approval/tasks/pending'
    const res = await api.get(endpoint, { params })
    const body = res.data
    pendingTasks.value = body?.data || []
    pendingTotal.value = body?.total || 0
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    tasksLoading.value = false
  }
}

const taskDetailVisible = ref(false)
const instanceLoading = ref(false)
const activeInstance = ref(null)
const activeTaskRef = ref(null)
const actionNote = ref('')
const noteError = ref('')
const actionSubmitting = ref(false)
const documentDetail = ref(null)
const documentLoading = ref(false)

// rowStatusLabel/rowStatusSeverity — for a WATCHER row, the task's own
// status is always PENDING (visible-but-not-actionable, see backend fix),
// which reads as ambiguous ("pending what? there's no approve/reject here").
// Watchers show the underlying instance's actual approval status instead —
// what the approver(s) have actually decided so far.
function rowStatusLabel(row) {
  if (row.participation_type === 'WATCHER' && row.instance_status) {
    return row.instance_status
  }
  return row.status
}
function rowStatusSeverity(row) {
  if (row.participation_type === 'WATCHER' && row.instance_status) {
    switch (row.instance_status) {
      case 'APPROVED': return 'success'
      case 'REJECTED': return 'danger'
      case 'CANCELLED': return 'secondary'
      default: return 'warn'
    }
  }
  return 'warn'
}

// stepApprover — the approve action recorded for a given step, if any, so
// each step in the preview can show who actually approved it (name,
// employee code, organization) once its status has moved past PENDING.
function stepApprover(stepOrder) {
  const actions = activeInstance.value?.actions || []
  return actions.find(a => a.step_order === stepOrder && a.action === 'APPROVE') || null
}

const activeTaskIsWatcher = computed(() => {
  if (!activeInstance.value || !activeTaskRef.value) return false
  const step = activeInstance.value.steps?.find(s => s.step_order === activeTaskRef.value.step_order)
  return step?.participation_type === 'WATCHER'
})

// Instance yang sudah final (APPROVED / REJECTED / CANCELLED) tidak boleh
// menampilkan tombol approve/reject — meskipun task row masih muncul di
// tab riwayat.
const isInstanceResolved = computed(() => {
  const status = activeInstance.value?.status
  return ['APPROVED', 'REJECTED', 'CANCELLED'].includes(status)
})

const isKPIModule = computed(() => {
  return ['performance_kpi_target', 'performance_kpi_realization'].includes(activeInstance.value?.module)
})

const isOKRModule = computed(() => {
  return ['okr_key_result', 'okr_assessment'].includes(activeInstance.value?.module)
})

const isLeaveModule = computed(() => activeInstance.value?.module === 'leave')

// ── Employee Movement (mutasi) ──
const isEmployeemovementModule = computed(() => activeInstance.value?.module === 'employeemovement')

// ── Recruitment (lowongan & penawaran kandidat) ──
const isRecruitmentModule = computed(() => activeInstance.value?.module === 'recruitment')
const isOfferModule = computed(() => activeInstance.value?.module === 'recruitment_offer')
const isTrainingRequestModule = computed(() => activeInstance.value?.module === 'training_request')

// ── Business Travel (perjalanan dinas) ──
const isBusinessTravelModule = computed(() => activeInstance.value?.module === 'business_travel')

// ── Attendance (overtime) ──
const isAttendanceModule = computed(() => activeInstance.value?.module === 'attendance')

// Fields hidden from the generic "submitted data" fallback view — internal
// IDs/timestamps/relations that aren't meaningful to a reviewer, or that
// are already shown elsewhere in the dialog.
const DOCUMENT_FIELD_DENYLIST = new Set([
  'id', 'created_at', 'updated_at', 'deleted_at',
  'employee_id', 'organization_id', 'period_id', 'template_id',
  'approval_instance_id', 'target_approval_instance_id', 'realization_approval_instance_id',
  'kr_approval_instance_id', 'assessment_approval_instance_id',
  'details', 'program_items', 'items', 'documents',
  // Business Travel (raw IDs already implied by the popup's own context —
  // approval module → module/document_id — not useful to show as fields).
  'requester_id', 'business_travel_id', 'participant_id'
])

function humanizeLabel(key) {
  return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function formatFieldValue(value) {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}

// Generic key-value fallback for modules without a dedicated renderer above
// (reimbursement, payroll, business travel, ...) — the approval module is
// document-agnostic, so this covers whatever module ends up here without
// needing a bespoke view per module.
const documentFields = computed(() => {
  if (!documentDetail.value || typeof documentDetail.value !== 'object') return []
  return Object.entries(documentDetail.value)
    .filter(([key, value]) => !DOCUMENT_FIELD_DENYLIST.has(key) && typeof value !== 'object')
    .map(([key, value]) => ({ label: humanizeLabel(key), value: formatFieldValue(value) }))
})

// Maps an approval "module" slug to the endpoint that returns the
// underlying document being approved — the approval module itself only
// knows module + document_id, not the document's shape.
function documentEndpointFor(module, documentId) {
  switch (module) {
    case 'leave':
      return `/api/v1/tenant/leave/requests/${documentId}`
    case 'reimbursement':
      return `/api/v1/tenant/reimbursements/requests/${documentId}`
    case 'employeemovement':
      return `/api/v1/tenant/employee-movements/movements/${documentId}`
    case 'attendance':
      return `/api/v1/tenant/attendance/overtime-requests/${documentId}`
    case 'payroll':
      return `/api/v1/tenant/payroll/runs/${documentId}`
    case 'performance_kpi_target':
    case 'performance_kpi_realization':
      return `/api/v1/tenant/performance/kpi/evaluations/${documentId}/full`
    case 'okr_key_result':
    case 'okr_assessment':
      return `/api/v1/tenant/performance/okr/evaluations/${documentId}/details`
    case 'recruitment':
      return `/api/v1/tenant/recruitment/requisitions/${documentId}`
    case 'recruitment_offer':
      return `/api/v1/tenant/recruitment/offers/${documentId}`
    case 'training_request':
      return `/api/v1/tenant/trainings/requests/${documentId}`
    case 'business_travel':
      return `/api/v1/tenant/attendance/business-travels/${documentId}`
    case 'business_travel_settlement':
      // Settlements only have a nested detail route (business-travels/:id/
      // settlements/:settlementId) — this flat lookup exists specifically so
      // callers here, which only know the document_id, can fetch it.
      return `/api/v1/tenant/attendance/business-travel-settlements/${documentId}`
    default:
      return null
  }
}

// The detail fetch itself; per-module name resolution (leave type, overtime
// employee, requisition org, training course, ...) is owned by each detail/*
// component via a watch on its documentDetail prop.
async function loadDocumentDetail(module, documentId) {
  documentDetail.value = null
  const endpoint = documentEndpointFor(module, documentId)
  if (!endpoint) return
  documentLoading.value = true
  try {
    const res = await api.get(endpoint)
    documentDetail.value = res.data?.data || null
  } catch {
    documentDetail.value = null
  } finally {
    documentLoading.value = false
  }
}

async function openTaskDetail(task) {
  activeTaskRef.value = task
  actionNote.value = ''
  noteError.value = ''
  taskDetailVisible.value = true
  instanceLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/approval/instances/${task.instance_id}`)
    activeInstance.value = res.data?.data || null
    if (activeInstance.value) {
      await loadDocumentDetail(activeInstance.value.module, activeInstance.value.document_id)
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    instanceLoading.value = false
  }
}

async function submitAction(action) {
  if (!activeTaskRef.value) return
  if (action === 'REJECT' && !actionNote.value?.trim()) {
    noteError.value = t('approval.reject_note_required')
    return
  }
  noteError.value = ''
  actionSubmitting.value = true
  try {
    await api.post(`/api/v1/tenant/approval/instances/${activeTaskRef.value.instance_id}/actions`, {
      action,
      note: actionNote.value || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.action_submitted'), life: 3000 })
    taskDetailVisible.value = false
    await loadTasks()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    actionSubmitting.value = false
  }
}

onMounted(() => {
  loadFlowOptions()
  loadTasks()
})
</script>
