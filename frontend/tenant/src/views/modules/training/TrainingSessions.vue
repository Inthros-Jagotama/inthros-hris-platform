<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2 flex-wrap">
        <SelectLabel v-model="statusFilter" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_status')" class="!w-52" showClear @update:modelValue="onFilterChange" />
        <SelectLabel v-model="courseFilter" :options="courseOptions" optionLabel="label" optionValue="value" :placeholder="t('training.filter_all_courses')" class="!w-56" filter showClear @update:modelValue="onFilterChange" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('training.session_new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-calendar text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('training.sessions_empty') }}</p>
        </div>
      </template>
      <Column field="session_code" :header="t('training.session_code')" style="width:120px">
        <template #body="{data}"><Tag :value="data.session_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="course_name" :header="t('training.course')">
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ courseName(data.course_id) }}</span></template>
      </Column>
      <Column field="trainer_name" :header="t('training.trainer')" style="width:160px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.trainer_name || '-' }}</span></template>
      </Column>
      <Column field="provider_type" :header="t('training.provider_type')" style="width:120px">
        <template #body="{data}">
          <Tag v-if="data.provider_type" :value="typeLabel(data.provider_type)" :severity="data.provider_type === 'EXTERNAL' ? 'warning' : 'info'" class="!text-xs !px-1.5 !py-0.5" />
          <span v-else class="text-gray-400">-</span>
        </template>
      </Column>
      <Column field="start_date" :header="t('training.start_date')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.start_date || '-' }}</span></template>
      </Column>
      <Column field="end_date" :header="t('training.end_date')" style="width:110px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.end_date || '-' }}</span></template>
      </Column>
      <Column field="max_quota" :header="t('training.max_quota')" style="width:90px">
        <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.max_quota }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:150px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="statusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:150px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('common.view')" @click="router.push(`/training/sessions/${data.id}`)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-arrow-right-arrow-left" size="small" text severity="warning" v-tooltip.left="t('training.change_status')" @click="openStatusDialog(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('training.session_edit') : t('training.session_new')" modal :style="{ width: '900px', maxWidth: '95vw' }" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('training.course')" required :errors="errors?.course_id">
          <SelectLabel v-model="form.course_id" :options="courseOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.course_id }" />
        </FormRow>
        <FormRow :label="t('training.trainer_name')" :errors="errors?.trainer_name">
          <TextInput v-model="form.trainer_name" maxlength="200" :placeholder="t('training.trainer_name_placeholder')" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.provider_type')">
            <SelectLabel v-model="form.provider_type" :options="providerTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
          </FormRow>
          <FormRow :label="t('training.delivery_mode')">
            <SelectLabel v-model="form.delivery_mode" :options="deliveryModeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
          </FormRow>
        </div>
        <template v-if="form.provider_type === 'EXTERNAL'">
          <FormRow :label="t('training.provider')" required :errors="errors?.provider_id">
            <SelectLabel v-model="form.provider_id" :options="providerOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.provider_id }" />
          </FormRow>
        </template>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.start_date')" required :errors="errors?.start_date">
            <DateInput v-model="form.start_date" :placeholder="t('training.start_date')" :class="{ 'p-invalid': errors?.start_date }" />
          </FormRow>
          <FormRow :label="t('training.end_date')" required :errors="errors?.end_date">
            <DateInput v-model="form.end_date" :placeholder="t('training.end_date')" :class="{ 'p-invalid': errors?.end_date }" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.location')">
            <TextInput v-model="form.location" maxlength="255" :placeholder="t('training.location_placeholder')" />
          </FormRow>
          <FormRow :label="t('training.max_quota')">
            <InputNumber v-model="form.max_quota" class="!w-full" :min="1" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('training.meeting_url')">
          <TextInput v-model="form.meeting_url" maxlength="500" :placeholder="t('training.meeting_url_placeholder')" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('training.registration_deadline')">
            <DateInput v-model="form.registration_deadline" :placeholder="t('training.registration_deadline')" />
          </FormRow>
          <FormRow :label="t('training.status')">
            <SelectLabel v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <Dialog v-model:visible="statusDialogVisible" :header="t('training.change_status')" modal :style="{ width: '420px' }">
      <div class="space-y-4">
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="statusForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="statusDialogVisible = false" />
          <Button :label="t('common.update')" size="small" :loading="statusSaving" :disabled="statusSaving" @click="handleStatusSave" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()
const route = useRoute()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const courses = ref([])
const providers = ref([])
const statusFilter = ref(null)
const courseFilter = ref(route.query.course_id || null)

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const statusDialogVisible = ref(false)
const statusSaving = ref(false)
const statusForm = ref({ id: null, status: '' })

const skeletonColumns = [
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const courseOptions = computed(() => courses.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))
const providerOptions = computed(() => providers.value.map(p => ({ label: p.name, value: p.id })))

function typeLabel(type) {
  const key = `training.type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
const providerTypeOptions = computed(() => ['IN_HOUSE', 'EXTERNAL'].map(v => ({ label: typeLabel(v), value: v })))
const deliveryModeOptions = computed(() => ['ONSITE', 'ONLINE', 'HYBRID', 'SELF_PACED'].map(v => ({ label: t(`training.mode_${v.toLowerCase()}`), value: v })))

function statusLabel(status) {
  const key = `training.status_${String(status || '').toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
const statusOptions = computed(() => ['DRAFT', 'SCHEDULED', 'REGISTRATION_OPEN', 'FULL', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'].map(v => ({ label: statusLabel(v), value: v })))

function statusSeverity(status) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'SCHEDULED': return 'info'
    case 'REGISTRATION_OPEN': return 'success'
    case 'FULL': return 'warning'
    case 'IN_PROGRESS': return 'info'
    case 'COMPLETED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}

function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || id
}

function defaultForm() {
  return {
    course_id: null, trainer_name: '',
    provider_type: null, delivery_mode: null, provider_id: null,
    start_date: '', end_date: '', location: '', max_quota: 30,
    meeting_url: '', registration_deadline: '', status: 'SCHEDULED'
  }
}

async function loadReferences() {
  const [cRes, pRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/providers', { params: { per_page: 500 } })
  ])
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
  providers.value = pRes.status === 'fulfilled' ? (pRes.value.data?.data || []) : []
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (statusFilter.value) params.status = statusFilter.value
    if (courseFilter.value) params.course_id = courseFilter.value
    const res = await api.get('/api/v1/tenant/trainings/sessions', { params })
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

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function onFilterChange() {
  currentPage.value = 1
  loadData()
}

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  itemStatus.value = item?.status || 'SCHEDULED'
  errors.value = {}
  form.value = item
    ? {
        course_id: item.course_id || null,
        trainer_name: item.trainer_name || '',
        provider_type: item.provider_type || null,
        delivery_mode: item.delivery_mode || null,
        provider_id: item.provider_id || null,
        start_date: item.start_date || '',
        end_date: item.end_date || '',
        location: item.location || '',
        max_quota: item.max_quota || 30,
        meeting_url: item.meeting_url || '',
        registration_deadline: item.registration_deadline ? item.registration_deadline.slice(0, 10) : '',
        status: item.status || 'SCHEDULED'
      }
    : defaultForm()
  dialogVisible.value = true
}

function resetForm() {
  form.value = defaultForm()
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.course_id) { errors.value = { course_id: t('form.required') }; return }
  if (!form.value.start_date) { errors.value = { start_date: t('form.required') }; return }
  if (!form.value.end_date) { errors.value = { end_date: t('form.required') }; return }
  if (form.value.provider_type === 'EXTERNAL' && !form.value.provider_id) { errors.value = { provider_id: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      course_id: form.value.course_id,
      trainer_name: form.value.trainer_name?.trim() || '',
      provider_type: form.value.provider_type || null,
      delivery_mode: form.value.delivery_mode || null,
      provider_id: form.value.provider_type === 'EXTERNAL' ? form.value.provider_id : '',
      start_date: form.value.start_date,
      end_date: form.value.end_date,
      location: form.value.location?.trim() || '',
      max_quota: form.value.max_quota || 30,
      meeting_url: form.value.meeting_url?.trim() || '',
      registration_deadline: form.value.registration_deadline || null
    }
    let sessionId = editingId.value
    if (editing.value) {
      await api.put(`/api/v1/tenant/trainings/sessions/${editingId.value}`, payload)
    } else {
      const created = await api.post('/api/v1/tenant/trainings/sessions', payload)
      sessionId = created.data?.data?.id
    }
    // Status dikelola terpisah via dialog transisi status — hanya diterapkan
    // saat berubah dari nilai aslinya (atau non-default saat create).
    const defaultStatus = editing.value ? itemStatus.value : 'SCHEDULED'
    if (sessionId && form.value.status && form.value.status !== defaultStatus) {
      await api.put(`/api/v1/tenant/trainings/sessions/${sessionId}/status`, { status: form.value.status })
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    dialogVisible.value = false
    await loadData()
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

// status asli record (sebelum dialog) — untuk deteksi perubahan status saat edit
const itemStatus = ref('SCHEDULED')

function openStatusDialog(item) {
  statusForm.value = { id: item.id, status: item.status || 'SCHEDULED' }
  statusDialogVisible.value = true
}

async function handleStatusSave() {
  if (!statusForm.value.id || !statusForm.value.status) return
  statusSaving.value = true
  try {
    await api.put(`/api/v1/tenant/trainings/sessions/${statusForm.value.id}/status`, { status: statusForm.value.status })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    statusDialogVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    statusSaving.value = false
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
