<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="pendingTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ pendingTotal }} {{ t('common.items') }}
        </span>
        <Button icon="pi pi-refresh" size="small" text severity="secondary" @click="loadPendingTasks" />
      </div>
      <Button :label="t('approval.flows')" icon="pi pi-sitemap" size="small" @click="router.push({ name: 'ApprovalFlows' })" />
    </div>

    <SkeletonTable v-if="tasksLoading" :columns="taskSkeletonColumns" :rows="6" />

    <DataTable v-else :value="pendingTasks" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('approval.empty_tasks') }}</p>
        </div>
      </template>
      <Column field="step_order" :header="t('approval.step')" style="width:80px">
        <template #body="{data}">#{{ data.step_order }}</template>
      </Column>
      <Column field="flow_name" :header="t('approval.instance')">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.flow_name || '-' }}</span>
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
          <Tag :value="data.status" severity="warn" class="!text-xs" />
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
      <div v-if="instanceLoading" class="py-8 text-center text-gray-400 text-sm">{{ t('common.loading') }}</div>
      <div v-else-if="activeInstance" class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <!-- Left column: the data being submitted -->
        <div class="space-y-3 md:border-r md:border-gray-200 md:dark:border-gray-700 md:pr-5">
          <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">{{ t('approval.submitted_data') }}</p>

          <div v-if="documentLoading" class="py-6 text-center text-gray-400 text-sm">{{ t('common.loading') }}</div>

          <template v-else-if="isKPIModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('kpi.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.employee_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('kpi.organization') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.organization_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('kpi.period') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.period_code || '-' }}</p>
              </div>
            </div>

            <div v-if="documentDetail?.details?.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.indicators') }}</p>
              <div class="space-y-1">
                <div v-for="d in documentDetail.details" :key="d.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                  <div class="font-medium text-gray-700 dark:text-gray-200">{{ d.indicator_name }}</div>
                  <div class="text-gray-500 dark:text-gray-400">
                    {{ t('kpi.target') }}: {{ d.target }} {{ d.unit_of_measurement }}
                    <span v-if="d.actual"> · {{ t('kpi.actual') }}: {{ d.actual }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="documentDetail?.program_items?.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.program_items') }}</p>
              <div class="space-y-1">
                <div v-for="p in documentDetail.program_items" :key="p.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                  <div class="font-medium text-gray-700 dark:text-gray-200">{{ p.title }}</div>
                  <div class="text-gray-500 dark:text-gray-400">
                    {{ t('kpi.target') }}: {{ p.target }} {{ p.unit_of_measurement }}
                    <span v-if="p.actual"> · {{ t('kpi.actual') }}: {{ p.actual }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="isOKRModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('okr.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.employee_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('okr.organization') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.organization_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('okr.period') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.period_code || '-' }}</p>
              </div>
            </div>

            <div v-if="okrObjectiveGroups.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('okr.objectives') }}</p>
              <div v-for="g in okrObjectiveGroups" :key="g.key" class="mb-2">
                <div class="text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{{ g.title }}</div>
                <div class="space-y-1">
                  <div v-for="kr in g.items" :key="kr.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                    <div class="font-medium text-gray-700 dark:text-gray-200">{{ kr.key_result_title }}</div>
                    <div class="text-gray-500 dark:text-gray-400">
                      {{ t('okr.target') }}: {{ kr.target_value }} {{ kr.unit }}
                      <span v-if="kr.actual_value"> · {{ t('okr.actual') }}: {{ kr.actual_value }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="isLeaveModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ leaveEmployeeName || '-' }}</p>
                <p v-if="leaveEmployeeCode" class="text-xs text-gray-400">{{ leaveEmployeeCode }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.leave_type') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ leaveTypeName || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.request_start_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.request_start_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.request_end_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.request_end_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.duration_mode') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.duration_mode || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.requested_days') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.requested_days ?? '-' }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('approval.submitted_at') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDate(documentDetail?.submitted_at) }}</p>
              </div>
              <div v-if="documentDetail?.leave_reason_note" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.leave_reason_note') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.leave_reason_note }}</p>
              </div>
            </div>
          </template>

          <div v-else-if="documentFields.length" class="space-y-2">
            <div v-for="f in documentFields" :key="f.label" class="text-sm">
              <p class="text-xs text-gray-400">{{ f.label }}</p>
              <p class="text-gray-700 dark:text-gray-200 break-words">{{ f.value }}</p>
            </div>
          </div>

          <p v-else class="text-xs text-gray-400">{{ t('approval.submitted_data_unavailable') }}</p>
        </div>

        <!-- Right column: existing approval data -->
        <div class="space-y-4">
          <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">{{ t('approval.approval_data') }}</p>

          <div class="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p class="text-xs text-gray-400">{{ t('approval.module') }}</p>
              <Tag :value="activeInstance.module" severity="info" class="!text-xs" />
            </div>
            <div>
              <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
              <Tag :value="activeInstance.status" severity="warn" class="!text-xs" />
            </div>
            <div class="col-span-2">
              <p class="text-xs text-gray-400">{{ t('approval.flow_name') }}</p>
              <p class="text-gray-700 dark:text-gray-200">{{ activeInstance.flow_name || '-' }}</p>
            </div>
          </div>

          <div>
            <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-2">{{ t('approval.steps') }}</p>
            <div class="space-y-1">
              <div v-for="s in activeInstance.steps" :key="s.id" class="flex items-center gap-2 text-xs px-2 py-1.5 rounded" :class="s.step_order === activeInstance.current_step ? 'bg-emerald-50 dark:bg-emerald-500/10' : ''">
                <span class="w-5 h-5 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">{{ s.step_order }}</span>
                <span class="font-medium text-gray-700 dark:text-gray-200">{{ s.step_name }}</span>
                <Tag :value="s.participation_type" :severity="s.participation_type === 'WATCHER' ? 'secondary' : 'success'" class="!text-xs !px-1.5 !py-0.5" />
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

          <div v-if="activeTaskIsWatcher" class="text-xs text-gray-400 flex items-center gap-1.5">
            <i class="pi pi-eye"></i> {{ t('approval.watcher_note') }}
          </div>
          <div v-else class="space-y-2">
            <Textarea v-model="actionNote" :placeholder="t('approval.note_placeholder')" rows="2" class="w-full" />
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('approval.reject')" severity="danger" outlined size="small" :loading="actionSubmitting" @click="submitAction('REJECT')" />
              <Button :label="t('approval.approve')" severity="success" size="small" :loading="actionSubmitting" @click="submitAction('APPROVE')" />
            </div>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()

function formatDate(v) {
  if (!v) return '-'
  return new Date(v).toLocaleString()
}

// formatDateOnly renders date-only strings (e.g. "2026-01-15") without a
// time component — parsed as local calendar date rather than via `new
// Date(v)` (which treats a bare "YYYY-MM-DD" as UTC midnight and can shift
// the displayed day depending on the viewer's timezone).
function formatDateOnly(v) {
  if (!v) return '-'
  const match = String(v).match(/^(\d{4})-(\d{2})-(\d{2})/)
  if (!match) return String(v)
  const [, y, m, d] = match
  return new Date(Number(y), Number(m) - 1, Number(d)).toLocaleDateString()
}

const pendingTasks = ref([])
const pendingTotal = ref(0)
const tasksLoading = ref(false)

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

async function loadPendingTasks() {
  tasksLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/approval/tasks/pending', { params: { page: 1, per_page: 50 } })
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
const actionSubmitting = ref(false)
const documentDetail = ref(null)
const documentLoading = ref(false)

const activeTaskIsWatcher = computed(() => {
  if (!activeInstance.value || !activeTaskRef.value) return false
  const step = activeInstance.value.steps?.find(s => s.step_order === activeTaskRef.value.step_order)
  return step?.participation_type === 'WATCHER'
})

const isKPIModule = computed(() => {
  return ['performance_kpi_target', 'performance_kpi_realization'].includes(activeInstance.value?.module)
})

const isOKRModule = computed(() => {
  return ['okr_key_result', 'okr_assessment'].includes(activeInstance.value?.module)
})

const isLeaveModule = computed(() => activeInstance.value?.module === 'leave')
const leaveEmployeeName = ref('')
const leaveEmployeeCode = ref('')
const leaveTypeName = ref('')

// Leave's own GET /requests/:id response only has employee_id/leave_type_id
// (no names) — resolved here client-side, same reasoning as why the generic
// documentFields fallback denylists raw id fields (not meaningful to a
// reviewer on their own).
async function loadLeaveNames(employeeId, leaveTypeId) {
  leaveEmployeeName.value = ''
  leaveEmployeeCode.value = ''
  leaveTypeName.value = ''
  const requests = []
  if (employeeId) {
    requests.push(
      api.get(`/api/v1/tenant/employees/${employeeId}`)
        .then(res => {
          const emp = res.data?.data
          leaveEmployeeName.value = emp?.name || ''
          leaveEmployeeCode.value = emp?.employee_id || ''
        })
        .catch(() => {})
    )
  }
  if (leaveTypeId) {
    requests.push(
      api.get(`/api/v1/tenant/leave/types/${leaveTypeId}`)
        .then(res => { leaveTypeName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  await Promise.all(requests)
}

const okrObjectiveGroups = computed(() => {
  const details = documentDetail.value?.details || []
  const groups = {}
  const order = []
  for (const d of details) {
    const key = d.objective_id || d.objective_title
    if (!groups[key]) {
      groups[key] = { key, title: d.objective_title, items: [] }
      order.push(key)
    }
    groups[key].items.push(d)
  }
  return order.map(key => groups[key])
})

// Fields hidden from the generic "submitted data" fallback view — internal
// IDs/timestamps/relations that aren't meaningful to a reviewer, or that
// are already shown elsewhere in the dialog.
const DOCUMENT_FIELD_DENYLIST = new Set([
  'id', 'created_at', 'updated_at', 'deleted_at',
  'employee_id', 'organization_id', 'period_id', 'template_id',
  'approval_instance_id', 'target_approval_instance_id', 'realization_approval_instance_id',
  'kr_approval_instance_id', 'assessment_approval_instance_id',
  'details', 'program_items', 'items', 'documents'
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
// (reimbursement, employeemovement, attendance, payroll, ...) — the
// approval module is document-agnostic, so this covers whatever module
// ends up here without needing a bespoke view per module.
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
    default:
      return null
  }
}

async function loadDocumentDetail(module, documentId) {
  documentDetail.value = null
  const endpoint = documentEndpointFor(module, documentId)
  if (!endpoint) return
  documentLoading.value = true
  try {
    const res = await api.get(endpoint)
    documentDetail.value = res.data?.data || null
    if (module === 'leave' && documentDetail.value) {
      await loadLeaveNames(documentDetail.value.employee_id, documentDetail.value.leave_type_id)
    }
  } catch {
    documentDetail.value = null
  } finally {
    documentLoading.value = false
  }
}

async function openTaskDetail(task) {
  activeTaskRef.value = task
  actionNote.value = ''
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
  actionSubmitting.value = true
  try {
    await api.post(`/api/v1/tenant/approval/instances/${activeTaskRef.value.instance_id}/actions`, {
      action,
      note: actionNote.value || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.action_submitted'), life: 3000 })
    taskDetailVisible.value = false
    await loadPendingTasks()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    actionSubmitting.value = false
  }
}

onMounted(() => {
  loadPendingTasks()
})
</script>
