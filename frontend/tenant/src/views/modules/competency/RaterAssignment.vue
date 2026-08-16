<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <h2 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('competency_360.rater_assignment') }}</h2>
      </div>
    </div>

    <!-- Event selector -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-center gap-2 flex-wrap">
        <Select v-model="eventFilter" :options="eventOptions" optionLabel="label" optionValue="value" showClear filter class="w-80" :placeholder="t('competency_360.select_event')" @change="reload" />
        <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('competency_360.rater_assignment_hint') }}</span>
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable
      v-else
      :value="targets"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.targets_empty') }}</p>
        </div>
      </template>
      <Column :header="t('competency_360.employee')" style="width:200px">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span>
        </template>
      </Column>
      <Column :header="t('competency_360.event_period')" style="width:140px">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ periodLabel(data) }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:110px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('competency_360.raters')" style="width:110px">
        <template #body="{data}">
          <Tag :value="raterSummaryLabel(data)" severity="info" class="!text-xs !px-1.5 !py-0.5" v-tooltip.top="{ value: raterSummaryTooltip(data), escape: false }" />
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:110px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button :label="t('competency_360.manage_raters')" icon="pi pi-users" size="small" text severity="info" @click="openRaters(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Raters dialog -->
    <Dialog v-model:visible="ratersVisible" :header="t('competency_360.manage_raters')" modal :style="{ width: '720px' }" @hide="resetRaters">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">
        {{ t('competency_360.subject') }}: <span class="font-medium text-gray-700 dark:text-gray-200">{{ employeeName(target?.employee_id) }}</span>
      </p>

      <div class="flex items-center gap-2 mb-2 flex-wrap">
        <Select v-model="raterForm.rater_employee_id" :options="raterEmployeeOptions" optionLabel="label" optionValue="value" filter showClear class="w-72" :placeholder="t('competency_360.select_rater')" />
        <Select v-model="raterForm.rater_type" :options="raterTypeOptions" optionLabel="label" optionValue="value" class="w-40" :placeholder="t('competency_360.rater_type')" />
        <Button :label="t('competency_360.add_rater')" icon="pi pi-plus" size="small" :loading="addingRater" @click="addRater" />
      </div>
      <div class="flex items-center gap-2 mb-3 flex-wrap">
        <Button :label="t('competency_360.auto_fill_raters')" icon="pi pi-sitemap" size="small" severity="secondary" :loading="autoFilling" :disabled="autoFilling" @click="autoFillRaters" />
        <span v-if="autoSummary" class="text-xs text-gray-500 dark:text-gray-400">{{ autoSummary }}</span>
      </div>

      <SkeletonTable v-if="ratersLoading" :columns="raterSkeletonColumns" :rows="5" />
      <DataTable v-else :value="raters" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <template #empty>
          <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
            <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
            <p class="text-sm font-medium">{{ t('competency_360.raters_empty') }}</p>
          </div>
        </template>
        <Column :header="t('competency_360.rater')">
          <template #body="{data}">
            <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.rater_employee_name || data.rater_employee_id?.slice(0, 8) }}</span>
          </template>
        </Column>
        <Column :header="t('competency_360.rater_type')" style="width:130px">
          <template #body="{data}"><Tag :value="raterTypeLabel(data.rater_type)" severity="secondary" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column :header="t('competency_360.weight')" style="width:80px">
          <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.weight }}</span></template>
        </Column>
        <Column field="status" :header="t('common.status')" style="width:110px">
          <template #body="{data}"><Tag :value="raterStatusLabel(data.status)" :severity="raterStatusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
        </Column>
        <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
          <template #body="{data}">
            <div class="flex items-center gap-1 justify-end">
              <Button v-if="data.status !== 'submitted'" icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="deleteRater(data)" />
            </div>
          </template>
        </Column>
      </DataTable>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const events = ref([])
const targets = ref([])
const employees = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const eventFilter = ref(null)

const ratersVisible = ref(false)
const ratersLoading = ref(false)
const raters = ref([])
const target = ref(null)
const raterForm = ref({ rater_employee_id: null, rater_type: 'peer' })
const addingRater = ref(false)
const autoFilling = ref(false)
const autoSummary = ref('')

const raterTypeOptions = [
  { label: t('competency_360.rater_type_self'), value: 'self' },
  { label: t('competency_360.rater_type_superior'), value: 'superior' },
  { label: t('competency_360.rater_type_peer'), value: 'peer' },
  { label: t('competency_360.rater_type_subordinate'), value: 'subordinate' },
  { label: t('competency_360.rater_type_other'), value: 'other' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-28' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-24', headerWidth: 'w-20' }
]

const raterSkeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-28' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const eventOptions = computed(() => events.value.map(e => ({ label: eventLabel(e), value: e.id })))
const raterEmployeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_code || e.employee_id})`, value: e.id })))

function eventLabel(e) {
  const parts = [e.period_year]
  if (e.period_number) parts.unshift(`${e.period_type === 'quarter' ? 'Q' : e.period_type === 'semester' ? 'S' : 'P'}${e.period_number}`)
  return parts.join(' ')
}

function periodLabel(data) {
  const ev = events.value.find(e => e.id === data.competency_event_id)
  return ev ? eventLabel(ev) : '-'
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'closed': return 'secondary'
    case 'finalized': return 'success'
    case 'submitted': return 'info'
    case 'pending_approval': return 'warn'
    case 'rejected': return 'danger'
    default: return 'secondary'
  }
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

function raterTypeLabel(type) {
  const key = `competency_360.rater_type_${type}`
  return t(key) !== key ? t(key) : type
}

// raterSummaryLabel — label ringkas kolom rater: "ditugaskan/seharusnya"
// (mis. 3/4). Fallback ke rater_count untuk data lama yang belum punya
// rater_summary.
function raterSummaryLabel(target) {
  const s = target?.rater_summary
  if (s && (s.expected > 0 || s.assigned > 0)) {
    return s.expected > 0 ? `${s.assigned}/${s.expected}` : String(s.assigned)
  }
  return String(target?.rater_count ?? 0)
}

// raterSummaryTooltip — rincian per tipe rater: seharusnya (expected) vs
// ditugaskan (assigned) vs sudah diisi (submitted).
function raterSummaryTooltip(target) {
  const s = target?.rater_summary
  if (!s || !s.details) return t('competency_360.view_raters')
  const rows = Object.entries(s.details)
    .map(([type, d]) => {
      const label = raterTypeLabel(type)
      return `<div class="flex items-center justify-between gap-6"><span>${label}</span><span>${d.assigned}/${d.expected} ${t('competency_360.rater_assigned')} &middot; ${d.submitted} ${t('competency_360.rater_filled')}</span></div>`
    })
    .join('')
  return `<div class="text-xs leading-5">${rows}</div>`
}

function employeeName(id) {
  if (!id) return '-'
  return employees.value.find(e => e.id === id)?.name || id.slice(0, 8)
}

async function loadReferences() {
  try {
    const [evRes, empRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/competency/events', { params: { per_page: 100 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    ])
    events.value = evRes.status === 'fulfilled' ? (evRes.value.data?.data || []) : []
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
    // Support ?target_id= dari halaman Events
    if (route.query.target_id) {
      const all = await api.get('/api/v1/tenant/competency/event-targets', { params: { per_page: 500 } })
      const tgt = (all.data?.data || []).find(x => x.id === route.query.target_id)
      if (tgt) {
        eventFilter.value = tgt.competency_event_id
        await loadTargets()
        openRaters(tgt)
      }
    }
  } catch {
    // fail-silent
  }
}

async function loadTargets() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/competency/event-targets', { params })
    let list = res.data?.data || []
    if (eventFilter.value) list = list.filter(t => t.competency_event_id === eventFilter.value)
    targets.value = list
    totalRecords.value = res.data?.total || 0
    if (res.data?.page) currentPage.value = res.data.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function reload() {
  currentPage.value = 1
  loadTargets()
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadTargets()
}

async function openRaters(tgt) {
  target.value = tgt
  raters.value = []
  raterForm.value = { rater_employee_id: null, rater_type: 'peer' }
  ratersVisible.value = true
  await loadRaters()
}

function resetRaters() {
  target.value = null
  raters.value = []
  raterForm.value = { rater_employee_id: null, rater_type: 'peer' }
}

async function loadRaters() {
  if (!target.value) return
  ratersLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/competency/event-targets/${target.value.id}/raters`)
    raters.value = res.data?.data || []
    target.value.rater_count = raters.value.length
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    ratersLoading.value = false
  }
}

// autoFillRaters mengisi rater superior (atasan) & subordinate (bawahan) dari
// struktur organisasi subject — saran dari backend, sudah mengecualikan rater
// yang sudah di-assign, jadi aman langsung di-POST.
async function autoFillRaters() {
  autoFilling.value = true
  autoSummary.value = ''
  try {
    const res = await api.get(`/api/v1/tenant/competency/event-targets/${target.value.id}/suggested-raters`)
    const sug = res.data?.data || { self: null, superior: [], subordinates: [] }
    const payload = [
      ...(sug.self ? [{ rater_employee_id: sug.self.id, rater_type: 'self' }] : []),
      ...(sug.superior || []).map(e => ({ rater_employee_id: e.id, rater_type: 'superior' })),
      ...(sug.subordinates || []).map(e => ({ rater_employee_id: e.id, rater_type: 'subordinate' }))
    ]
    if (payload.length === 0) {
      autoSummary.value = t('competency_360.no_suggested_raters')
      return
    }
    await api.post(`/api/v1/tenant/competency/event-targets/${target.value.id}/raters`, { raters: payload })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('competency_360.raters_filled'), life: 3000 })
    await loadRaters()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    autoFilling.value = false
  }
}

async function addRater() {
  if (!raterForm.value.rater_employee_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('competency_360.select_rater'), life: 3000 })
    return
  }
  addingRater.value = true
  try {
    await api.post(`/api/v1/tenant/competency/event-targets/${target.value.id}/raters`, {
      raters: [{ rater_employee_id: raterForm.value.rater_employee_id, rater_type: raterForm.value.rater_type }]
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    raterForm.value = { rater_employee_id: null, rater_type: 'peer' }
    await loadRaters()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    addingRater.value = false
  }
}

async function deleteRater(rater) {
  try {
    await api.delete(`/api/v1/tenant/competency/raters/${rater.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadRaters()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

onMounted(() => {
  loadReferences()
  loadTargets()
})
</script>
