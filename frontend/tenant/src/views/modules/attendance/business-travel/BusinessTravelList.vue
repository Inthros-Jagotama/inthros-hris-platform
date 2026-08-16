<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Select v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" showClear :placeholder="t('common.status')" class="w-44" @change="reload" />
        <Button :label="t('common.add')" icon="pi pi-plus" size="small" @click="openCreateDialog" />
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
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-briefcase text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('business_travel.empty') }}</p>
        </div>
      </template>
      <Column field="request_number" :header="t('business_travel.request_number')" style="width:160px">
        <template #body="{data}">
          <a class="text-indigo-600 dark:text-indigo-400 hover:underline font-medium cursor-pointer" @click="router.push(`/attendance/business-travel/${data.id}`)">{{ data.request_number }}</a>
        </template>
      </Column>
      <Column field="title" :header="t('business_travel.field_title')" />
      <Column field="start_date" :header="t('business_travel.start_date')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.start_date, locale) }}</span></template>
      </Column>
      <Column field="end_date" :header="t('business_travel.end_date')" style="width:120px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatDate(data.end_date, locale) }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:140px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:80px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button icon="pi pi-eye" size="small" text @click="router.push(`/attendance/business-travel/${data.id}`)" />
        </template>
      </Column>
    </DataTable>

    <!-- ── Dialog: Ajukan Perjalanan Dinas ── -->
    <Dialog v-model:visible="dialogVisible" :header="t('business_travel.new')" modal :style="{ width: '520px' }" @hide="resetForm">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('business_travel.new_description') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('business_travel.field_title')" required :errors="errors?.title">
          <TextInput v-model="form.title" />
        </FormRow>
        <FormRow :label="t('business_travel.purpose')" :errors="errors?.purpose">
          <TextInput v-model="form.purpose" textarea :rows="2" :placeholder="t('business_travel.purpose_placeholder')" />
        </FormRow>
        <FormRow :label="t('business_travel.start_date')" required :errors="errors?.start_date">
          <DateInput v-model="form.start_date" />
        </FormRow>
        <FormRow :label="t('business_travel.end_date')" required :errors="errors?.end_date">
          <DateInput v-model="form.end_date" />
        </FormRow>
        <FormRow :label="t('business_travel.origin')" :errors="errors?.origin">
          <TextInput v-model="form.origin" :placeholder="t('business_travel.origin_placeholder')" />
        </FormRow>
        <FormRow :label="t('business_travel.description_field')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
      </div>
      <Message severity="info" :closable="false" class="!mt-1">{{ t('business_travel.destinations_hint') }}</Message>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Message from 'primevue/message'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const statusFilter = ref(null)

const dialogVisible = ref(false)
const saving = ref(false)
const errors = ref({})
const form = ref({ title: '', purpose: '', start_date: '', end_date: '', origin: '', description: '' })

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const statusOptions = [
  { label: 'DRAFT', value: 'DRAFT' },
  { label: 'SUBMITTED', value: 'SUBMITTED' },
  { label: 'APPROVED', value: 'APPROVED' },
  { label: 'REJECTED', value: 'REJECTED' },
  { label: 'IN_PROGRESS', value: 'IN_PROGRESS' },
  { label: 'COMPLETED', value: 'COMPLETED' },
  { label: 'CLOSED', value: 'CLOSED' },
  { label: 'CANCELLED', value: 'CANCELLED' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-24', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-24', headerWidth: 'w-20' },
  { type: 'icons', count: 1, headerWidth: 'w-20' }
]

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

function statusLabel(status) {
  const key = `business_travel.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    const res = await api.get('/api/v1/tenant/attendance/business-travels', { params })
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

function reload() {
  currentPage.value = 1
  loadData()
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function openCreateDialog() {
  errors.value = {}
  form.value = { title: '', purpose: '', start_date: '', end_date: '', origin: '', description: '' }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { title: '', purpose: '', start_date: '', end_date: '', origin: '', description: '' }
  errors.value = {}
}

async function handleSave() {
  errors.value = {}
  if (!form.value.title?.trim()) { errors.value = { title: t('form.required') }; return }
  if (!form.value.start_date) { errors.value = { start_date: t('form.required') }; return }
  if (!form.value.end_date) { errors.value = { end_date: t('form.required') }; return }
  saving.value = true
  try {
    const res = await api.post('/api/v1/tenant/attendance/business-travels', form.value)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    const newId = res.data?.data?.id
    if (newId) {
      router.push(`/attendance/business-travel/${newId}`)
    } else {
      await loadData()
    }
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
