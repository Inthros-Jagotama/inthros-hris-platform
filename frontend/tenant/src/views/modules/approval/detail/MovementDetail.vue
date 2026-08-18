<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('employee_movement.employee')" :value="props.documentDetail?.employee_name" :description="props.documentDetail?.employee_code" class="col-span-2" />
    <ViewLabel :label="t('employee_movement.movement_type')">
      <Tag :value="movementTypeLabel(props.documentDetail?.movement_type)" :severity="movementTypeSeverity(props.documentDetail?.movement_type)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel :label="t('common.status')">
      <Tag :value="movementStatusLabel(props.documentDetail?.status)" :severity="movementStatusSeverity(props.documentDetail?.status)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel v-if="props.documentDetail?.decision_letter_number" :label="t('employee_movement.decision_letter_number')" :value="props.documentDetail.decision_letter_number" mono break-all class="col-span-2" />
    <ViewLabel :label="t('employee_movement.decision_letter_date')" :value="formatDate(props.documentDetail?.decision_letter_date, locale.value)" />
    <ViewLabel :label="t('employee_movement.effective_date')" :value="formatDate(props.documentDetail?.effective_date, locale.value)" />

    <!-- Dari → Ke -->
    <div v-if="hasMovementFrom" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-1.5">
      <p class="text-xs uppercase tracking-wide text-gray-400 font-medium">{{ t('employee_movement.from') }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.from_organization_name || '-' }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.from_position_name || '-' }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.from_employment_status_name || '-' }}</p>
    </div>
    <div v-if="hasMovementTo" class="rounded-lg border border-emerald-200 dark:border-emerald-900/40 p-3 space-y-1.5">
      <p class="text-xs uppercase tracking-wide text-emerald-500 font-medium">{{ t('employee_movement.to') }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.to_organization_name || '-' }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.to_position_name || '-' }}</p>
      <p class="text-gray-700 dark:text-gray-200">{{ props.documentDetail?.to_employment_status_name || '-' }}</p>
    </div>

    <ViewLabel v-if="props.documentDetail?.reason" :label="t('employee_movement.reason')" :value="props.documentDetail.reason" break-all class="col-span-2" />
    <ViewLabel v-if="props.documentDetail?.notes" :label="t('employee_movement.notes')" :value="props.documentDetail.notes" break-all class="col-span-2" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'
import { formatDate } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t, locale } = useI18n()

// movementTypeLabel/Severity — bilingual via employee_movement.type_* keys
// (sama seperti halaman Movements), dengan fallback raw slug.
function movementTypeLabel(type) {
  if (!type) return '-'
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function movementTypeSeverity(type) {
  switch (type) {
    case 'promotion': return 'success'
    case 'demotion': return 'danger'
    case 'mutation': return 'info'
    case 'contract_extension': return 'warning'
    case 'status_change': return 'info'
    case 'retirement': return 'secondary'
    case 'offboarding': return 'danger'
    default: return 'secondary'
  }
}

// movementStatusLabel/Severity — bilingual via employee_movement.status_* keys.
function movementStatusLabel(status) {
  if (!status) return '-'
  const key = `employee_movement.status_${status}`
  return t(key) !== key ? t(key) : status
}

function movementStatusSeverity(status) {
  switch (status) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'rejected': return 'danger'
    case 'executed': return 'success'
    case 'cancelled': return 'secondary'
    case 'cancellation_pending': return 'warning'
    default: return 'secondary'
  }
}

const hasMovementFrom = computed(() => !!(props.documentDetail && (props.documentDetail.from_organization_name || props.documentDetail.from_position_name || props.documentDetail.from_employment_status_name)))
const hasMovementTo = computed(() => !!(props.documentDetail && (props.documentDetail.to_organization_name || props.documentDetail.to_position_name || props.documentDetail.to_employment_status_name)))
</script>
