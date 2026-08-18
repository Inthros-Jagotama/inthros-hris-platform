<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('training.request_employee')" :value="trainingRequestEmployeeName" class="col-span-2" />
    <ViewLabel :label="t('training.request_course')" :value="trainingRequestCourseName" class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.requested_date" :label="t('training.request_requested_date')" :value="formatDate(props.documentDetail.requested_date, locale.value)" />
    <ViewLabel :label="t('training.request_priority')" :value="props.documentDetail?.priority" />
    <ViewLabel :label="t('common.status')">
      <Tag :value="trainingRequestStatusLabel(props.documentDetail?.status)" :severity="trainingRequestStatusSeverity(props.documentDetail?.status)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel v-if="props.documentDetail?.reason" :label="t('training.request_reason')" :value="props.documentDetail.reason" break-all class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.supervisor_note" :label="t('training.request_supervisor_note')" :value="props.documentDetail.supervisor_note" break-all class="col-span-2" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { formatDate } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t, locale } = useI18n()

// Training request detail hanya membawa employee_id/course_id — nama
// di-resolve client-side (pola loadLeaveNames).
const trainingRequestEmployeeName = ref('')
const trainingRequestCourseName = ref('')
async function loadTrainingRequestNames(employeeId, courseId) {
  trainingRequestEmployeeName.value = ''
  trainingRequestCourseName.value = ''
  const requests = []
  if (employeeId) {
    requests.push(
      api.get(`/api/v1/tenant/employees/${employeeId}`)
        .then(res => { trainingRequestEmployeeName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  if (courseId) {
    requests.push(
      api.get(`/api/v1/tenant/trainings/courses/${courseId}`)
        .then(res => { trainingRequestCourseName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  await Promise.all(requests)
}

watch(() => props.documentDetail, (doc) => {
  loadTrainingRequestNames(doc?.employee_id, doc?.course_id)
}, { immediate: true })

function trainingRequestStatusLabel(status) {
  if (!status) return '-'
  const key = `training.request_status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function trainingRequestStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'submitted': return 'info'
    case 'pending_approval': return 'info'
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'cancelled': return 'secondary'
    default: return 'secondary'
  }
}
</script>
