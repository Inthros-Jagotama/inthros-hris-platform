<template>
  <div class="space-y-4">
    <div class="flex items-center gap-2 flex-wrap">
      <SelectLabel
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('applications.filter_status')"
        class="!w-48"
        showClear
        @update:modelValue="onFilterChange()"
      />
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('applications.new_application')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="items"
      lazy
      :totalRecords="totalRecords"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 15, 25, 50]"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-send text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('applications.empty') }}</p>
        </div>
      </template>

      <Column :header="t('applications.candidate')">
        <template #body="{ data }">
          <button type="button" class="text-left group" @click="router.push(`/recruitment/applications/${data.id}`)">
            <p class="font-medium text-gray-800 dark:text-gray-100 group-hover:text-sky-600 dark:group-hover:text-sky-400">{{ candidateName(data.candidate_id) }}</p>
            <p class="text-xs text-gray-400 dark:text-gray-500">{{ candidateEmail(data.candidate_id) }}</p>
          </button>
        </template>
      </Column>

      <Column :header="t('applications.requisition')">
        <template #body="{ data }">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ requisitionTitle(data.requisition_id) }}</span>
        </template>
      </Column>

      <Column :header="t('common.status')" style="width: 140px">
        <template #body="{ data }">
          <Tag :value="t('applications.status_' + (data.status || 'new').toLowerCase())" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>

      <Column :header="t('applications.score')" style="width: 90px">
        <template #body="{ data }">
          <Tag v-if="data.score !== null && data.score !== undefined" :value="Math.round(data.score) + '%'" :severity="scoreSeverity(data.score)" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-sm text-gray-400 dark:text-gray-500">—</span>
        </template>
      </Column>

      <Column :header="t('applications.applied_at')" style="width: 150px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ formatTimestamp(data.applied_at) }}</span>
        </template>
      </Column>

      <Column :header="t('common.actions')" :exportable="false" style="width: 90px">
        <template #body="{ data }">
          <Button icon="pi pi-arrow-right" text size="small" class="!w-7 !h-7" @click="router.push(`/recruitment/applications/${data.id}`)" />
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="t('applications.new_application')" :modal="true" class="!w-[min(95vw,520px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('requisitions.title')" :required="true">
          <SelectLabel v-model="form.requisition_id" :options="requisitionOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('candidates.title')" :required="true">
          <SelectLabel v-model="form.candidate_id" :options="candidateOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('candidates.notes')">
          <Textarea v-model="form.notes" :rows="2" class="!w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" @click="save()" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import SelectLabel from '@/components/SelectLabel.vue'
import FormRow from '@/components/FormRow.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const statusFilter = ref(null)
const dialogVisible = ref(false)

const requisitions = ref([])
const candidates = ref([])

const skeletonColumns = [
  { field: 'candidate', header: 'Candidate', width: '26%' },
  { field: 'requisition', header: 'Requisition', width: '26%' },
  { field: 'status', header: 'Status', width: '12%' },
  { field: 'score', header: 'Score', width: '8%' },
  { field: 'applied_at', header: 'Applied At', width: '14%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = computed(() => ['NEW', 'SCREENED', 'SHORTLISTED', 'INTERVIEWED', 'OFFERED', 'ACCEPTED', 'REJECTED', 'WITHDRAWN'].map(v => ({ label: t(`applications.status_${v.toLowerCase()}`), value: v })))

const requisitionOptions = computed(() => requisitions.value.map(r => ({ label: r.title, value: r.id })))
const candidateOptions = computed(() => candidates.value.map(c => ({ label: `${c.first_name} ${c.last_name}`, value: c.id })))

const emptyForm = () => ({ requisition_id: null, candidate_id: null, notes: '' })
const form = ref(emptyForm())

function candidateName(id) {
  const c = candidates.value.find(x => x.id === id)
  return c ? `${c.first_name} ${c.last_name}` : id
}
function candidateEmail(id) {
  const c = candidates.value.find(x => x.id === id)
  return c ? c.email : ''
}
function requisitionTitle(id) {
  const r = requisitions.value.find(x => x.id === id)
  return r ? r.title : id
}

function statusSeverity(status) {
  switch (status) {
    case 'NEW': return 'secondary'
    case 'SCREENED': return 'info'
    case 'SHORTLISTED': return 'info'
    case 'INTERVIEWED': return 'help'
    case 'OFFERED': return 'warn'
    case 'ACCEPTED': return 'success'
    case 'REJECTED': return 'danger'
    case 'WITHDRAWN': return 'danger'
    default: return 'secondary'
  }
}

function scoreSeverity(score) {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warn'
  return 'danger'
}

function formatTimestamp(value) {
  if (!value) return '—'
  const ms = Number(value) / 1000000
  if (!Number.isFinite(ms) || ms <= 0) return '—'
  return new Date(ms).toLocaleDateString()
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/recruitment/applications', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadOptions() {
  try {
    const [reqRes, candRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/recruitment/requisitions', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/candidates', { params: { per_page: 500 } })
    ])
    requisitions.value = reqRes.status === 'fulfilled' ? (reqRes.value.data?.data || []) : []
    candidates.value = candRes.status === 'fulfilled' ? (candRes.value.data?.data || []) : []
  } catch {
    // fail-silent — dropdown/nama kosong
  }
}

function openDialog() {
  form.value = emptyForm()
  dialogVisible.value = true
}

async function save() {
  if (!form.value.requisition_id || !form.value.candidate_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    if (!payload.notes) delete payload.notes
    await api.post('/api/v1/tenant/recruitment/applications', payload)
    dialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('applications.created'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    saving.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

onMounted(() => {
  loadData()
  loadOptions()
})
</script>
