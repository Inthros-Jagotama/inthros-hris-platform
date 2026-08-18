<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.manager_assessments') }}</h2>
      </div>
      <Button v-if="detail" :label="t('common.back')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="closeDetail" />
    </div>

    <!-- ── Form penilaian bawahan ── -->
    <template v-if="detail">
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
        <div class="flex items-center justify-between flex-wrap gap-2">
          <div>
            <p class="text-sm font-semibold text-navy-800 dark:text-gray-100">
              {{ t('competency_360.assessing') }}: {{ detail.target?.employee_id ? subjectName(detail.target.employee_id) : '-' }}
            </p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              {{ t('competency_360.rater_type') }}: {{ raterTypeLabel(detail.rater.rater_type) }} ·
              {{ t('competency_360.status') }}: <Tag :value="raterStatusLabel(detail.rater.status)" :severity="detail.rater.status === 'submitted' ? 'success' : 'info'" class="!text-xs !px-1.5 !py-0.5" />
            </p>
          </div>
          <div class="flex items-center gap-2">
            <Button v-if="detail.rater.status !== 'submitted'" :label="t('competency_360.save_draft')" icon="pi pi-save" size="small" severity="secondary" outlined :loading="saving" @click="saveResponses(false)" />
            <Button v-if="detail.rater.status !== 'submitted'" :label="t('competency_360.submit')" icon="pi pi-check" size="small" :loading="submitting" @click="saveResponses(true)" />
          </div>
        </div>

        <div v-if="!detail.indicators?.length" class="text-center py-8 text-gray-400 dark:text-gray-500">
          <i class="pi pi-clone text-3xl mb-2 opacity-50"></i>
          <p class="text-sm">{{ t('competency_360.no_indicators_hint') }}</p>
        </div>

        <div v-else class="space-y-3">
          <div v-for="ind in detail.indicators" :key="ind.indicator_id" class="border border-gray-200 dark:border-gray-700 rounded-lg p-3">
            <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ ind.statement }}</p>
            <div class="flex items-center gap-4 mt-2 flex-wrap">
              <div class="flex items-center gap-1.5">
                <span class="text-xs text-gray-500 dark:text-gray-400">{{ t('competency_360.rating') }}:</span>
                <Select
                  v-model="responses[ind.indicator_id]"
                  :options="ratingOptions"
                  optionLabel="label"
                  optionValue="value"
                  class="w-44"
                  size="small"
                  :disabled="detail.rater.status === 'submitted'"
                  :placeholder="t('common.select')"
                />
              </div>
              <div class="flex-1 min-w-56">
                <TextInput
                  v-model="comments[ind.indicator_id]"
                  :placeholder="t('competency_360.comment_placeholder')"
                  :disabled="detail.rater.status === 'submitted'"
                />
              </div>
            </div>
          </div>
        </div>

        <Message v-if="saveError" severity="error" :closable="false" class="!text-xs">{{ saveError }}</Message>
      </div>
    </template>

    <!-- ── Daftar bawahan (dari struktur organisasi) ── -->
    <template v-else>
      <Message v-if="!employeeId" severity="warn" :closable="false">{{ t('competency_360.no_employee_linked') }}</Message>

      <div v-else class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-center gap-2 flex-wrap mb-3">
          <Select v-model="eventFilter" :options="eventOptions" optionLabel="label" optionValue="value" showClear filter class="w-80" :placeholder="t('competency_360.select_event')" @change="loadList" />
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.manager_assessments_hint') }}</span>
        </div>

        <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
        <DataTable
          v-else
          :value="items"
          size="small"
          class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        >
          <template #empty>
            <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
              <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
              <p class="text-sm font-medium">{{ t('competency_360.manager_assessments_empty') }}</p>
            </div>
          </template>
          <Column :header="t('competency_360.subordinate')">
            <template #body="{data}">
              <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.employee_name || data.employee_id?.slice(0, 8) || '-' }}</span>
            </template>
          </Column>
          <Column :header="t('competency_360.event_period')" style="width:140px">
            <template #body="{data}">
              <span class="text-gray-600 dark:text-gray-300">{{ eventName(data.competency_event_id) }}</span>
            </template>
          </Column>
          <Column :header="t('competency_360.my_rating_status')" style="width:150px">
            <template #body="{data}">
              <Tag
                v-if="data.rater_id"
                :value="raterStatusLabel(data.rater_status)"
                :severity="raterStatusSeverity(data.rater_status)"
                class="!text-xs !px-1.5 !py-0.5"
              />
              <Tag v-else :value="t('competency_360.not_assigned')" severity="secondary" class="!text-xs !px-1.5 !py-0.5" />
            </template>
          </Column>
          <Column :header="t('competency_360.assigned_at')" style="width:160px">
            <template #body="{data}">
              <span class="text-gray-500 dark:text-gray-400">{{ data.assigned_at ? formatDate(data.assigned_at) : '-' }}</span>
            </template>
          </Column>
          <Column :header="t('common.actions')" style="width:130px" frozen alignFrozen="right">
            <template #body="{data}">
              <div class="flex items-center gap-1 justify-end">
                <Button
                  v-if="!data.rater_id"
                  :label="t('competency_360.assign_fill')"
                  icon="pi pi-user-plus"
                  size="small"
                  text
                  severity="info"
                  :loading="assigningFor === data.employee_id"
                  @click="assignAndOpen(data)"
                />
                <Button
                  v-else
                  :label="data.rater_status === 'submitted' ? t('common.view') : t('competency_360.fill')"
                  icon="pi pi-pencil"
                  size="small"
                  text
                  @click="openDetail(data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useMyEmployee } from '@/composables/useMyEmployee'
import { getErrorMessage } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import Message from 'primevue/message'
import SkeletonTable from '@/components/SkeletonTable.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const { employeeId, loadMyEmployeeId } = useMyEmployee()

const events = ref([])
const eventFilter = ref(null)
const loading = ref(false)
const items = ref([])
const assigningFor = ref(null)
const detail = ref(null)
const responses = ref({})
const comments = ref({})
const saving = ref(false)
const submitting = ref(false)
const saveError = ref('')

// Opsi nilai diambil dari skala penilaian template (competency_rating_scale_items
// via scale_id) — fallback 1-5 bila template belum memakai skala.
const ratingOptions = computed(() => {
  const items = detail.value?.scale?.items
  if (Array.isArray(items) && items.length > 0) {
    return items
      .filter(i => i.label?.trim() !== '' || i.value !== undefined)
      .map(i => ({ label: i.label || String(i.value), value: i.value }))
      .sort((a, b) => a.value - b.value)
  }
  return [
    { label: '1', value: 1 },
    { label: '2', value: 2 },
    { label: '3', value: 3 },
    { label: '4', value: 4 },
    { label: '5', value: 5 }
  ]
})

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-28' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-24' }
]

const eventOptions = computed(() => events.value.map(e => ({ label: eventLabel(e), value: e.id })))

function eventLabel(e) {
  const parts = [e.period_year]
  if (e.period_number) parts.unshift(`${e.period_type === 'quarter' ? 'Q' : e.period_type === 'semester' ? 'S' : 'P'}${e.period_number}`)
  return parts.join(' ')
}

function eventName(id) {
  const ev = events.value.find(e => e.id === id)
  return ev ? eventLabel(ev) : '-'
}

function raterTypeLabel(type) {
  const key = `competency_360.rater_type_${type}`
  return t(key) !== key ? t(key) : type
}

function raterStatusLabel(status) {
  const key = `competency_360.rater_status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function raterStatusSeverity(status) {
  switch (status) {
    case 'submitted': return 'success'
    case 'started': return 'info'
    default: return 'secondary'
  }
}

function subjectName(id) {
  const item = items.value.find(i => i.employee_id === id)
  return item?.employee_name || id?.slice(0, 8)
}

async function loadReferences() {
  try {
    const res = await api.get('/api/v1/tenant/competency/events', { params: { per_page: 100 } })
    events.value = res.data?.data || []
    if (events.value.length > 0) eventFilter.value = events.value[0].id
  } catch {
    events.value = []
  }
}

async function loadList() {
  if (!employeeId.value) return
  loading.value = true
  try {
    const params = {}
    if (eventFilter.value) params.event_id = eventFilter.value
    const res = await api.get('/api/v1/tenant/competency/manager-assessments', { params })
    items.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

// Manager belum di-assign sebagai superior rater untuk bawahan tsb — assign
// diri sendiri (superior) lalu buka form penilaian.
async function assignAndOpen(data) {
  assigningFor.value = data.employee_id
  try {
    await api.post(`/api/v1/tenant/competency/event-targets/${data.target_id}/raters`, {
      raters: [{ rater_employee_id: employeeId.value, rater_type: 'superior' }]
    })
    await loadList()
    const updated = items.value.find(i => i.employee_id === data.employee_id)
    if (updated?.rater_id) {
      await openDetail(updated)
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    assigningFor.value = null
  }
}

async function openDetail(item) {
  try {
    const res = await api.get(`/api/v1/tenant/competency/my-assessments/${item.rater_id}`)
    detail.value = res.data?.data || res.data
    responses.value = {}
    comments.value = {}
    for (const r of detail.value.responses || []) {
      responses.value[r.indicator_id] = r.rating_value
      comments.value[r.indicator_id] = r.comment || ''
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

function closeDetail() {
  detail.value = null
  responses.value = {}
  comments.value = {}
  loadList()
}

async function saveResponses(submit) {
  saveError.value = ''
  const payload = Object.entries(responses.value)
    .filter(([, v]) => v !== null && v !== undefined && v !== '')
    .map(([indicator_id, rating_value]) => ({
      indicator_id,
      rating_value: Number(rating_value),
      comment: comments.value[indicator_id]?.trim() || ''
    }))
  if (payload.length === 0) {
    saveError.value = t('competency_360.no_responses_error')
    return
  }
  if (submit) {
    submitting.value = true
  } else {
    saving.value = true
  }
  try {
    await api.post(`/api/v1/tenant/competency/my-assessments/${detail.value.rater.id}/responses`, { responses: payload })
    if (submit) {
      await api.post(`/api/v1/tenant/competency/my-assessments/${detail.value.rater.id}/submit`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('competency_360.submitted'), life: 3000 })
      closeDetail()
    } else {
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    }
  } catch (e) {
    saveError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    saving.value = false
    submitting.value = false
  }
}

onMounted(async () => {
  await loadMyEmployeeId()
  await loadReferences()
  await loadList()
})
</script>
