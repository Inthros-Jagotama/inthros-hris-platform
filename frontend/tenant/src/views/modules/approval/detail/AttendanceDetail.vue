<template>
  <div class="space-y-4">
    <!-- Group: Informasi Lembur -->
    <div>
      <div class="flex items-center gap-2 mb-2">
        <i class="pi pi-clock text-teal-400 text-sm"></i>
        <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('attendance.overtime_info') }}</h2>
        <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
      </div>
      <div class="grid grid-cols-2 gap-3 text-sm">
        <ViewLabel :label="t('attendance.employee')" :value="overtimeEmployeeName" :description="overtimeEmployeeCode" class="col-span-2" />
        <ViewLabel :label="t('attendance.flow_type')">
          <Tag :value="props.documentDetail?.flow_type === 'ASSIGNED' ? t('attendance.flow_assigned') : t('attendance.flow_self')" :severity="props.documentDetail?.flow_type === 'ASSIGNED' ? 'info' : 'secondary'" class="!text-xs" />
        </ViewLabel>
        <ViewLabel :label="t('common.status')">
          <Tag :value="overtimeStatusLabel(props.documentDetail?.status)" :severity="overtimeStatusSeverity(props.documentDetail?.status)" class="!text-xs" />
        </ViewLabel>
        <ViewLabel v-if="isAssignedFlow" :label="t('attendance.assigned_by')" :value="overtimeAssignedByName" class="col-span-2" />
        <ViewLabel :label="t('attendance.work_date')" :value="formatDate(props.documentDetail?.work_date, locale.value)" />
        <ViewLabel :label="t('attendance.requested_minutes')">
          <span>{{ props.documentDetail?.requested_minutes ?? '-' }} min</span>
          <span v-if="isOvertimeCrossDay" class="text-xs text-amber-500 ml-1">· {{ t('attendance.overtime_cross_day') }}</span>
        </ViewLabel>
        <ViewLabel :label="t('attendance.start_time')" :value="formatTime(props.documentDetail?.start_time_local)" />
        <ViewLabel :label="t('attendance.end_time')" :value="formatTime(props.documentDetail?.end_time_local)" />
        <ViewLabel v-if="props.documentDetail?.reason" :label="t('attendance.reason')" :value="props.documentDetail.reason"/>
        <ViewLabel v-if="props.documentDetail?.approval_note" :label="t('attendance.approval_note')" :value="props.documentDetail.approval_note"/>
        <ViewLabel :label="t('approval.submitted_at')" :value="formatDateDateTime(props.documentDetail?.created_at, locale.value)"/>
      </div>
    </div>

    <!-- Group: Detail Aktual -->
    <div v-if="hasActualData">
      <div class="flex items-center gap-2 mb-2">
        <i class="pi pi-check-square text-teal-400 text-sm"></i>
        <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('attendance.actual_data') }}</h2>
        <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
      </div>
      <div class="grid grid-cols-2 gap-3 text-sm">
        <ViewLabel :label="t('attendance.actual_start_time')" :value="formatTime(props.documentDetail?.actual_start_time_local)" />
        <ViewLabel :label="t('attendance.actual_end_time')" :value="formatTime(props.documentDetail?.actual_end_time_local)" />
        <ViewLabel v-if="props.documentDetail?.actual_note" :label="t('attendance.actual_note')" :value="props.documentDetail.actual_note"/>
        <ViewLabel v-if="props.documentDetail?.calculated_minutes != null" :label="t('attendance.calculated_minutes')">
          <span class="font-medium">{{ props.documentDetail.calculated_minutes }} min</span>
        </ViewLabel>
        <ViewLabel v-if="props.documentDetail?.attachment_url" :label="t('attendance.attachment')" class="col-span-2">
          <a :href="props.documentDetail.attachment_url" target="_blank" class="text-sm text-emerald-600 dark:text-emerald-400 hover:underline">
            <i class="pi pi-paperclip mr-1"></i>{{ t('attendance.attachment') }}
          </a>
        </ViewLabel>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { formatDate, formatDateDateTime, formatTime } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null },
  submitterName: { type: String, default: null }
})

const { t, locale } = useI18n()

const overtimeEmployeeName = ref('')
const overtimeEmployeeCode = ref('')

// Overtime's detail response only carries employee_id (no name) — resolved
// here client-side, same reasoning as the leave names loader.
async function loadOvertimeNames(employeeId) {
  overtimeEmployeeName.value = ''
  overtimeEmployeeCode.value = ''
  if (!employeeId) return
  try {
    const res = await api.get(`/api/v1/tenant/employees/${employeeId}`)
    const emp = res.data?.data
    overtimeEmployeeName.value = emp?.name || ''
    overtimeEmployeeCode.value = emp?.employee_id || ''
  } catch {}
}

watch(() => props.documentDetail, (doc) => {
  loadOvertimeNames(doc?.employee_id)
}, { immediate: true })

// ASSIGNED flow: the person who submitted the approval task IS the assigner
// (manager), so "Assigned By" maps straight to the task's submitter.
const isAssignedFlow = computed(() => props.documentDetail?.flow_type === 'ASSIGNED')
const overtimeAssignedByName = computed(() => {
  if (!isAssignedFlow.value) return null
  return props.submitterName || '-'
})

function overtimeStatusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'PENDING_APPROVAL': return 'info'
    case 'WAITING_ACTUAL': return 'warning'
    case 'ACTUAL_SUBMITTED': return 'info'
    case 'CANCELLED': return 'secondary'
    default: return 'secondary'
  }
}

function overtimeStatusLabel(status) {
  if (!status) return '-'
  const key = `attendance.status_${status.toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

// Lintas hari: end di hari berikutnya tetap disimpan sebagai timestamp tanggal
// +1 hari (pola FE), jadi bandingkan jam-jam (HH:MM) — sama seperti logika
// isCrossDayOf di AttendanceOvertime.vue.
const isOvertimeCrossDay = computed(() => {
  const s = props.documentDetail?.start_time_local
  const e = props.documentDetail?.end_time_local
  if (!s || !e) return false
  const sd = new Date(s)
  const ed = new Date(e)
  if (isNaN(sd.getTime()) || isNaN(ed.getTime())) return false
  return ed.getHours() * 60 + ed.getMinutes() <= sd.getHours() * 60 + sd.getMinutes()
})

const hasActualData = computed(() => {
  const d = props.documentDetail
  return !!(d && (d.actual_start_time_local || d.actual_end_time_local || d.actual_note || d.attachment_url))
})
</script>
