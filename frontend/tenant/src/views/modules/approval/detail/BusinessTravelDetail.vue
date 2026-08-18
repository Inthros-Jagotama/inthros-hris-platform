<template>
  <div class="space-y-4">
    <!-- Informasi utama perjalanan dinas (2 kolom) -->
    <div class="grid grid-cols-2 gap-3 text-sm">
      <ViewLabel :label="t('business_travel.field_title')" :value="props.documentDetail?.title" :description="props.documentDetail?.request_number" />
      <ViewLabel :label="t('business_travel.purpose')" :value="props.documentDetail?.purpose"/>
      <ViewLabel :label="t('business_travel.description_field')" :value="props.documentDetail.description"/>
      <ViewLabel :label="t('business_travel.origin')" :value="props.documentDetail?.origin"/>
      <ViewLabel :label="t('business_travel.start_date')" :value="formatDate(props.documentDetail?.start_date, locale.value)"/>
      <ViewLabel :label="t('business_travel.end_date')" :value="formatDate(props.documentDetail?.end_date, locale.value)"/>
    </div>
    <!-- Destinations & Participants: satu row, dua kolom -->
    <div class="grid grid-cols-2 gap-4">
      <!-- Destinations -->
      <div v-if="props.documentDetail?.destinations?.length">
        <ViewLabel :label="t('business_travel.destinations')">
          <div class="space-y-1">
            <div v-for="d in props.documentDetail.destinations" :key="d.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
              <div class="font-medium text-gray-700 dark:text-gray-200">{{ destinationLabel(d) }}</div>
              <div v-if="d.purpose" class="text-gray-500 dark:text-gray-400">— {{ d.purpose }}</div>
              <div v-if="d.arrival_date || d.departure_date" class="text-gray-500 dark:text-gray-400">
                {{ d.arrival_date ? formatDate(d.arrival_date, locale.value) : '' }}{{ d.arrival_date && d.departure_date ? ' \u2192 ' : '' }}{{ d.departure_date ? formatDate(d.departure_date, locale.value) : '' }}
              </div>
            </div>
          </div>
        </ViewLabel>
      </div>

      <!-- Participants -->
      <div v-if="props.documentDetail?.participants?.length">
        <ViewLabel :label="t('business_travel.participants')">
          <div class="space-y-1">
            <div v-for="p in props.documentDetail.participants" :key="p.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60 flex items-center justify-between gap-2">
              <div class="min-w-0">
                <span class="font-medium text-gray-700 dark:text-gray-200">{{ participantDisplayName(p) }}</span>
                <span v-if="p.organization" class="text-gray-500 dark:text-gray-400"> — {{ p.organization }}</span>
              </div>
              <div class="flex items-center gap-1 shrink-0">
                <Tag :value="p.role" severity="secondary" class="!text-xs !px-1.5 !py-0.5" />
                <Tag :value="p.participant_type" severity="info" class="!text-xs !px-1.5 !py-0.5" />
              </div>
            </div>
          </div>
        </ViewLabel>
      </div>
    </div>
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

// Participant EMPLOYEE hanya membawa employee_id (nama tidak disimpan di
// record peserta) — di-resolve dari daftar karyawan (pola BusinessTravelDetail).
const employees = ref([])
async function loadEmployees() {
  try {
    employees.value = (await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })).data?.data || []
  } catch {
    employees.value = []
  }
}

watch(() => props.documentDetail, (doc) => {
  if (doc?.participants?.length) loadEmployees()
}, { immediate: true })

function statusLabel(status) {
  if (!status) return '-'
  const key = `business_travel.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'secondary'
    case 'SUBMITTED': return 'info'
    case 'IN_PROGRESS': return 'warn'
    case 'COMPLETED': return 'info'
    case 'CLOSED': return 'success'
    default: return 'secondary'
  }
}

function destinationLabel(d) {
  return [d.city, d.province, d.country].filter(Boolean).join(', ') || d.location || '-'
}

function participantDisplayName(p) {
  if (p.participant_type === 'EMPLOYEE') {
    return employees.value.find(e => e.id === p.employee_id)?.name || p.employee_id?.slice(0, 8) || '-'
  }
  return p.name || '-'
}
</script>
