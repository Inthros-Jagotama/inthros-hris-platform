<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('leave.employee')" :value="leaveEmployeeName" :description="leaveEmployeeCode"/>
    <ViewLabel :label="t('leave.leave_type')" :value="leaveTypeName"/>
    <ViewLabel :label="t('leave.request_start_date')" :value="formatDate(props.documentDetail?.request_start_date, locale.value)" />
    <ViewLabel :label="t('leave.request_end_date')" :value="formatDate(props.documentDetail?.request_end_date, locale.value)" />
    <ViewLabel :label="t('leave.duration_mode')" :value="props.documentDetail?.duration_mode" />
    <ViewLabel :label="t('leave.requested_days')" :value="props.documentDetail?.requested_days" />
    <ViewLabel :label="t('approval.submitted_at')" :value="formatDateDateTime(props.documentDetail?.submitted_at, locale.value)" class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.leave_reason_note" :label="t('leave.leave_reason_note')" :value="props.documentDetail.leave_reason_note" break-all class="col-span-2" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { formatDate, formatDateDateTime } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t, locale } = useI18n()

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

watch(() => props.documentDetail, (doc) => {
  loadLeaveNames(doc?.employee_id, doc?.leave_type_id)
}, { immediate: true })
</script>
