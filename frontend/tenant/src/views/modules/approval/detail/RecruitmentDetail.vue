<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('requisitions.title')" :value="props.documentDetail?.title" :description="props.documentDetail?.requisition_number" class="col-span-2" />
    <ViewLabel :label="t('common.status')">
      <Tag :value="requisitionStatusLabel(props.documentDetail?.status)" :severity="requisitionStatusSeverity(props.documentDetail?.status)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel :label="t('requisitions.priority')">
      <Tag :value="requisitionPriorityLabel(props.documentDetail?.priority)" :severity="requisitionPrioritySeverity(props.documentDetail?.priority)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel v-if="requisitionOrganizationName" :label="t('requisitions.org')" :value="requisitionOrganizationName" class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.department" :label="t('requisitions.department')" :value="props.documentDetail.department" />
    <ViewLabel v-if="props.documentDetail?.employment_type" :label="t('requisitions.employment_type')" :value="props.documentDetail.employment_type" />
    <ViewLabel v-if="props.documentDetail?.location" :label="t('requisitions.location')" :value="props.documentDetail.location" />
    <ViewLabel v-if="props.documentDetail?.reason_type" :label="t('requisitions.reason_type')" :value="requisitionReasonLabel(props.documentDetail.reason_type)" />
    <ViewLabel :label="t('requisitions.slots_available')" :value="props.documentDetail?.slots_available" />
    <ViewLabel v-if="props.documentDetail?.slots_filled != null" :label="t('requisitions.slots')" :value="props.documentDetail.slots_filled" />
    <ViewLabel v-if="props.documentDetail?.target_start_date" :label="t('requisitions.target_start_date')" :value="formatDate(props.documentDetail.target_start_date, locale.value)" />
    <ViewLabel v-if="props.documentDetail?.min_salary || props.documentDetail?.max_salary" :label="t('requisitions.min_salary') + ' \u2013 ' + t('requisitions.max_salary')">
      <span>{{ formatCurrency(props.documentDetail?.min_salary) }} \u2013 {{ formatCurrency(props.documentDetail?.max_salary) }}</span>
    </ViewLabel>
    <ViewLabel v-if="props.documentDetail?.description" :label="t('requisitions.description_label')" :value="props.documentDetail.description" break-all class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.requirements" :label="t('requisitions.requirements')" :value="props.documentDetail.requirements" break-all class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.responsibilities" :label="t('requisitions.responsibilities')" :value="props.documentDetail.responsibilities" break-all class="col-span-2" />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { formatDate, formatCurrency } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t, locale } = useI18n()

// requisition_number tidak dibawa response detail, jadi tidak perlu
// enrichment. Organization name di-resolve dari organization_id (response
// requisition hanya membawa id, sama seperti leave).
const requisitionOrganizationName = ref('')
async function loadRequisitionOrganization(organizationId) {
  requisitionOrganizationName.value = ''
  if (!organizationId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${organizationId}`)
    const org = res.data?.data
    requisitionOrganizationName.value = org?.nomenclature || org?.full_code || org?.code || ''
  } catch {}
}

watch(() => props.documentDetail, (doc) => {
  loadRequisitionOrganization(doc?.organization_id)
}, { immediate: true })

function requisitionStatusLabel(status) {
  if (!status) return '-'
  const key = `requisitions.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function requisitionStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'submitted': return 'info'
    case 'open': return 'success'
    case 'in_progress': return 'info'
    case 'filled': return 'success'
    case 'rejected': return 'danger'
    case 'cancelled': return 'secondary'
    default: return 'secondary'
  }
}

function requisitionPriorityLabel(priority) {
  if (!priority) return '-'
  const key = `requisitions.priority_${String(priority).toLowerCase()}`
  return t(key) !== key ? t(key) : priority
}

function requisitionPrioritySeverity(priority) {
  switch (String(priority).toLowerCase()) {
    case 'urgent': return 'danger'
    case 'high': return 'warn'
    case 'medium': return 'info'
    case 'low': return 'secondary'
    default: return 'secondary'
  }
}

function requisitionReasonLabel(reasonType) {
  if (!reasonType) return '-'
  const key = `requisitions.reason_${String(reasonType).toLowerCase()}`
  return t(key) !== key ? t(key) : reasonType
}
</script>
