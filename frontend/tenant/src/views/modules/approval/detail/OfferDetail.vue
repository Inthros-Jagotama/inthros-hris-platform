<template>
  <div class="grid grid-cols-2 gap-3 text-sm">
    <ViewLabel :label="t('offers.offer_number')" :value="props.documentDetail?.offer_number" class="col-span-2" />
    <ViewLabel :label="t('common.status')">
      <Tag :value="offerStatusLabel(props.documentDetail?.status)" :severity="offerStatusSeverity(props.documentDetail?.status)" class="!text-xs" />
    </ViewLabel>
    <ViewLabel v-if="props.documentDetail?.employment_type" :label="t('offers.employment_type')" :value="props.documentDetail.employment_type" />
    <ViewLabel :label="t('offers.salary')" :value="formatCurrency(props.documentDetail?.salary)" />
    <ViewLabel :label="t('offers.allowances')" :value="formatCurrency(props.documentDetail?.allowances)" />
    <ViewLabel v-if="props.documentDetail?.start_date" :label="t('offers.start_date')" :value="formatDate(props.documentDetail.start_date, locale.value)" />
    <ViewLabel v-if="props.documentDetail?.expiry_date" :label="t('offers.expiry_date')" :value="formatDate(props.documentDetail.expiry_date, locale.value)" />
    <ViewLabel v-if="props.documentDetail?.benefits" :label="t('offers.benefits')" :value="props.documentDetail.benefits" break-all class="col-span-2" />
  </div>
</template>

<script setup>
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'
import { formatDate, formatCurrency } from '@/utils/formatDate'
import ViewLabel from '@/components/ViewLabel.vue'

const props = defineProps({
  documentDetail: { type: Object, default: null }
})

const { t, locale } = useI18n()

function offerStatusLabel(status) {
  if (!status) return '-'
  const key = `offers.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function offerStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'sent': return 'info'
    case 'accepted': return 'success'
    case 'rejected': return 'danger'
    case 'expired': return 'secondary'
    case 'withdrawn': return 'secondary'
    default: return 'secondary'
  }
}
</script>
