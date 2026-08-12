<template>
  <div class="space-y-4">
    <!-- Toolbar: filter status + tombol buat offer -->
    <div class="flex items-center gap-2 flex-wrap">
      <SelectLabel
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('offers.filter_status')"
        class="!w-52"
        showClear
        @update:modelValue="onFilterChange()"
      />
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('offers.new_offer')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-file-edit text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('offers.empty') }}</p>
        </div>
      </template>

      <!-- G-3: nomor offer (auto OFF-YYYYMM-XXXXXXXX) -->
      <Column :header="t('offers.offer_number')" style="width: 170px">
        <template #body="{ data }">
          <span class="text-xs font-mono text-gray-600 dark:text-gray-300">{{ data.offer_number || '—' }}</span>
        </template>
      </Column>

      <!-- Kandidat (enrich dari application → candidate) -->
      <Column :header="t('offers.candidate')" style="min-width: 200px">
        <template #body="{ data }">
          <div class="flex items-center gap-2.5 min-w-0">
            <div class="w-8 h-8 rounded-lg bg-indigo-50 dark:bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 flex items-center justify-center shrink-0">
              <i class="pi pi-user text-xs"></i>
            </div>
            <div class="min-w-0">
              <p class="font-medium text-gray-800 dark:text-gray-100 truncate">{{ candidateName(data.application_id) }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ applicationStatus(data.application_id) }}</p>
            </div>
          </div>
        </template>
      </Column>

      <Column :header="t('offers.employment_type')" style="width: 130px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.employment_type || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('offers.salary')" style="width: 140px">
        <template #body="{ data }">
          <span class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ formatCurrency(data.salary) }}</span>
        </template>
      </Column>

      <Column :header="t('offers.start_date')" style="width: 120px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.start_date || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('offers.expiry_date')" style="width: 120px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.expiry_date || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('common.status')" style="width: 150px">
        <template #body="{ data }">
          <Tag
            :value="t('offers.status_' + (data.status || '').toLowerCase())"
            :severity="statusSeverity(data.status)"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>

      <!-- G-3: aksi kontekstual per status -->
      <Column :header="t('common.actions')" :exportable="false" style="width: 220px">
        <template #body="{ data }">
          <div class="flex items-center gap-1.5 flex-wrap">
            <!-- DRAFT: edit + delete + submit (ke Central Approval) + withdraw -->
            <template v-if="data.status === 'DRAFT'">
              <Button icon="pi pi-pencil" size="small" severity="secondary" outlined class="!w-7 !h-7" :title="t('common.edit')" @click="openEditDialog(data)" />
              <Button icon="pi pi-trash" size="small" severity="danger" outlined class="!w-7 !h-7" :title="t('common.delete')" @click="openDeleteDialog(data)" />
              <Button :label="t('offers.submit')" icon="pi pi-send" size="small" severity="info" outlined class="!text-xs !px-2.5 !py-1" @click="openActionDialog(data, 'submit')" />
              <Button icon="pi pi-times" size="small" severity="secondary" text class="!w-7 !h-7" :title="t('offers.withdraw')" @click="openActionDialog(data, 'withdraw')" />
            </template>
            <!-- PENDING_APPROVAL: menunggu approval -->
            <template v-else-if="data.status === 'PENDING_APPROVAL'">
              <span class="text-xs text-gray-400 dark:text-gray-500 italic">{{ t('offers.waiting_approval') }}</span>
            </template>
            <!-- APPROVED: send ke kandidat + withdraw -->
            <template v-else-if="data.status === 'APPROVED'">
              <Button :label="t('offers.send')" icon="pi pi-send" size="small" severity="success" class="!text-xs !px-2.5 !py-1" @click="openActionDialog(data, 'send')" />
              <Button icon="pi pi-times" size="small" severity="secondary" text class="!w-7 !h-7" :title="t('offers.withdraw')" @click="openActionDialog(data, 'withdraw')" />
            </template>
            <!-- SENT: accept / reject dari kandidat -->
            <template v-else-if="data.status === 'SENT'">
              <Button :label="t('offers.accept')" icon="pi pi-check" size="small" severity="success" class="!text-xs !px-2.5 !py-1" @click="openActionDialog(data, 'accept')" />
              <Button :label="t('offers.reject')" icon="pi pi-times" size="small" severity="danger" outlined class="!text-xs !px-2.5 !py-1" @click="openActionDialog(data, 'reject')" />
            </template>
            <span v-else class="text-xs text-gray-400 dark:text-gray-500 italic">—</span>
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Dialog konfirmasi aksi workflow (G-3) -->
    <ConfirmActionDialog
      v-model:visible="actionDialogVisible"
      :title="actionDialogTitle"
      :message="actionDialogMessage"
      :confirm-label="actionConfirmLabel"
      :loading="actionLoading"
      :error-msg="actionError"
      :icon="actionIcon"
      :severity="actionSeverity"
      @confirm="runAction()"
    />

    <!-- Dialog konfirmasi delete (hanya DRAFT) -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('offers.delete_confirm_title')"
      :message="t('offers.delete_confirm_message', { number: pendingDelete?.offer_number || '' })"
      :confirm-label="t('common.delete')"
      :loading="deleting"
      :error-msg="deleteError"
      @confirm="deleteOffer()"
    />

    <!-- Dialog: buat/edit offer -->
    <Dialog v-model:visible="dialogVisible" :header="dialogHeader" :modal="true" class="!w-[min(95vw,640px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <!-- Create: pilih aplikasi pipeline. Edit: aplikasi tidak bisa diganti —
             tampilkan read-only (SelectLabel tidak meneruskan prop disabled,
             dan mengganti application_id saat edit tidak berefek di backend) -->
        <FormRow v-if="!isEdit" :label="t('offers.application')" :required="true" class="md:col-span-2">
          <SelectLabel
            v-model="form.application_id"
            :options="applicationOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
        </FormRow>
        <FormRow v-else :label="t('offers.application')" class="md:col-span-2">
          <div class="flex items-center gap-2.5 text-sm text-gray-700 dark:text-gray-200 bg-gray-50 dark:bg-gray-800/60 rounded-lg px-3 py-2">
            <i class="pi pi-user text-gray-400"></i>
            <span class="font-medium truncate">{{ candidateName(form.application_id) }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500 shrink-0">{{ applicationStatus(form.application_id) }}</span>
          </div>
        </FormRow>
        <FormRow :label="t('offers.employment_type')">
          <TextInput v-model="form.employment_type" :placeholder="t('offers.employment_type_placeholder')" class="!w-full" />
        </FormRow>
        <FormRow :label="t('offers.salary')">
          <InputNumber v-model="form.salary" :min="0" mode="currency" currency="IDR" locale="id-ID" class="!w-full" />
        </FormRow>
        <FormRow :label="t('offers.allowances')">
          <InputNumber v-model="form.allowances" :min="0" mode="currency" currency="IDR" locale="id-ID" class="!w-full" />
        </FormRow>
        <FormRow :label="t('offers.start_date')">
          <DateInput v-model="form.start_date" class="!w-full" />
        </FormRow>
        <FormRow :label="t('offers.expiry_date')">
          <DateInput v-model="form.expiry_date" class="!w-full" />
        </FormRow>
        <FormRow :label="t('offers.benefits')" class="md:col-span-2">
          <Textarea v-model="form.benefits" :rows="3" :placeholder="t('offers.benefits_placeholder')" class="!w-full" />
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
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import SelectLabel from '@/components/SelectLabel.vue'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const actionLoading = ref(false)
const actionError = ref('')
const deleting = ref(false)
const deleteError = ref('')
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const statusFilter = ref(null)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const actionDialogVisible = ref(false)
const pendingAction = ref(null)
const pendingActionType = ref('')
const deleteDialogVisible = ref(false)
const pendingDelete = ref(null)

// Enrichment: application_id → candidate name + application status
const applications = ref([])
const candidates = ref([])

const skeletonColumns = [
  { field: 'offer_number', header: 'Offer No.', width: '14%' },
  { field: 'candidate', header: 'Candidate', width: '22%' },
  { field: 'employment_type', header: 'Type', width: '11%' },
  { field: 'salary', header: 'Salary', width: '12%' },
  { field: 'start_date', header: 'Start', width: '10%' },
  { field: 'expiry_date', header: 'Expiry', width: '10%' },
  { field: 'status', header: 'Status', width: '12%' },
  { field: 'actions', header: 'Actions', width: '9%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = computed(() =>
  ['DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'SENT', 'ACCEPTED', 'REJECTED', 'EXPIRED', 'WITHDRAWN'].map(v => ({ label: t(`offers.status_${v.toLowerCase()}`), value: v }))
)

const candidateById = computed(() => {
  const m = {}
  for (const c of candidates.value) m[c.id] = c
  return m
})

const applicationById = computed(() => {
  const m = {}
  for (const a of applications.value) m[a.id] = a
  return m
})

// Dropdown form: aplikasi pipeline (belum REJECTED/WITHDRAWN/ACCEPTED) berlabel kandidat
const applicationOptions = computed(() => {
  const skip = new Set(['REJECTED', 'WITHDRAWN', 'ACCEPTED'])
  return applications.value
    .filter(a => !skip.has(a.status))
    .map(a => {
      const c = candidateById.value[a.candidate_id]
      const name = c ? `${c.first_name} ${c.last_name || ''}`.trim() : a.id.slice(0, 8)
      return { label: `${name} — ${t('offers.application_status_' + (a.status || '').toLowerCase())}`, value: a.id }
    })
})

const emptyForm = () => ({
  application_id: null,
  employment_type: '',
  salary: null,
  allowances: null,
  benefits: '',
  start_date: null,
  expiry_date: null
})

const form = ref(emptyForm())

const dialogHeader = computed(() => isEdit.value ? t('offers.edit_offer') : t('offers.new_offer'))

const actionDialogTitle = computed(() => {
  const key = `offers.action_${pendingActionType.value}_title`
  return t(key)
})

const actionDialogMessage = computed(() => {
  const number = pendingAction.value?.offer_number || ''
  return t(`offers.action_${pendingActionType.value}_message`, { number })
})

const actionConfirmLabel = computed(() => t(`offers.${pendingActionType.value}`))

const actionIcon = computed(() => {
  switch (pendingActionType.value) {
    case 'accept': return 'pi pi-check'
    case 'reject': return 'pi pi-times'
    case 'withdraw': return 'pi pi-arrow-left'
    default: return 'pi pi-send'
  }
})

const actionSeverity = computed(() => {
  switch (pendingActionType.value) {
    case 'accept': return 'success'
    case 'reject': return 'danger'
    case 'withdraw': return 'warn'
    default: return 'info'
  }
})

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'SENT': return 'info'
    case 'ACCEPTED': return 'success'
    case 'PENDING_APPROVAL': return 'warn'
    case 'DRAFT': return 'secondary'
    case 'REJECTED': return 'danger'
    case 'EXPIRED': return 'danger'
    case 'WITHDRAWN': return 'secondary'
    default: return 'secondary'
  }
}

function candidateName(applicationId) {
  const app = applicationById.value[applicationId]
  if (!app) return '—'
  const c = candidateById.value[app.candidate_id]
  return c ? `${c.first_name} ${c.last_name || ''}`.trim() : '—'
}

function applicationStatus(applicationId) {
  const app = applicationById.value[applicationId]
  return app ? t('offers.application_status_' + (app.status || '').toLowerCase()) : ''
}

function formatCurrency(value) {
  if (!value && value !== 0) return '—'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value)
}

function openActionDialog(row, type) {
  pendingAction.value = row
  pendingActionType.value = type
  actionError.value = ''
  actionDialogVisible.value = true
}

async function runAction() {
  if (!pendingAction.value) return
  actionLoading.value = true
  actionError.value = ''
  try {
    await api.post(`/api/v1/tenant/recruitment/offers/${pendingAction.value.id}/${pendingActionType.value}`, {})
    actionDialogVisible.value = false
    pendingAction.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t(`offers.action_${pendingActionType.value}_done`), life: 3000 })
    loadData()
  } catch (e) {
    actionError.value = getErrorMessage(e, t('message.failed_to_save'))
  } finally {
    actionLoading.value = false
  }
}

function openDeleteDialog(row) {
  pendingDelete.value = row
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function deleteOffer() {
  if (!pendingDelete.value) return
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/recruitment/offers/${pendingDelete.value.id}`)
    deleteDialogVisible.value = false
    pendingDelete.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('offers.deleted'), life: 3000 })
    loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.failed_to_save'))
  } finally {
    deleting.value = false
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/recruitment/offers', { params })
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

async function loadEnrichment() {
  try {
    const [appRes, candRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/recruitment/applications', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/candidates', { params: { per_page: 500 } })
    ])
    applications.value = appRes.status === 'fulfilled' ? (appRes.value.data?.data || []) : []
    candidates.value = candRes.status === 'fulfilled' ? (candRes.value.data?.data || []) : []
  } catch {
    // fail-silent — kolom kandidat jadi '—', offer tetap bisa dikelola
  }
}

function openDialog() {
  isEdit.value = false
  editId.value = null
  form.value = emptyForm()
  dialogVisible.value = true
}

function openEditDialog(row) {
  isEdit.value = true
  editId.value = row.id
  form.value = {
    application_id: row.application_id,
    employment_type: row.employment_type || '',
    salary: row.salary || null,
    allowances: row.allowances || null,
    benefits: row.benefits || '',
    start_date: row.start_date || null,
    expiry_date: row.expiry_date || null
  }
  dialogVisible.value = true
}

async function save() {
  if (!form.value.application_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('offers.application_required'), life: 4000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    // Field kosong → jangan dikirim (backend default / auto-generate)
    Object.keys(payload).forEach(k => {
      if (payload[k] === '' || payload[k] === null || payload[k] === undefined) delete payload[k]
    })
    if (isEdit.value) {
      // Aplikasi tidak bisa diganti saat edit (UpdateOfferRequest tanpa
      // application_id) — jangan kirim.
      delete payload.application_id
      await api.put(`/api/v1/tenant/recruitment/offers/${editId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/recruitment/offers', payload)
    }
    dialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: isEdit.value ? t('offers.updated') : t('offers.created'), life: 3000 })
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
  loadEnrichment()
})
</script>
