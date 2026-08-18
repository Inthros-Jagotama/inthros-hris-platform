<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <span v-if="items.length > 0" class="text-xs text-gray-400 dark:text-gray-500">
        {{ items.length }} {{ t('common.items') }}
      </span>
      <span v-else></span>
      <Button :label="t('assessments.new_assessment')" icon="pi pi-plus" size="small" @click="openDialog()" />
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="6" />

    <DataTable
      v-else
      :value="items"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-12 text-gray-400 dark:text-gray-500">
          <i class="pi pi-clipboard text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('assessments.empty') }}</p>
        </div>
      </template>

      <Column field="name" :header="t('assessments.name')">
        <template #body="{ data }">
          <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
        </template>
      </Column>
      <Column :header="t('assessments.type')" style="width:150px">
        <template #body="{ data }">
          <Tag :value="typeLabel(data.type)" :severity="typeSeverity(data.type)" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column :header="t('assessments.requisition')" style="width:200px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ requisitionTitle(data.requisition_id) }}</span>
        </template>
      </Column>
      <Column :header="t('assessments.scheduled_at')" style="width:150px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ formatScheduledAt(data.scheduled_at) }}</span>
        </template>
      </Column>
      <Column :header="t('assessments.participants')" style="width:90px">
        <template #body="{ data }">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ participantCounts[data.id] ?? '-' }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{ data }">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Form assessment -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('assessments.edit_assessment') : t('assessments.new_assessment')" :modal="true" class="!w-[min(95vw,520px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('assessments.name')" required :errors="errors?.name">
          <TextInput v-model="form.name" class="!w-full" :placeholder="t('assessments.name_placeholder')" />
        </FormRow>
        <FormRow :label="t('assessments.type')">
          <SelectLabel
            v-model="form.type"
            :options="typeOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('common.select')"
            class="!w-full"
            showClear
          />
        </FormRow>
        <FormRow :label="t('assessments.requisition')">
          <SelectLabel
            v-model="form.requisition_id"
            :options="requisitionOptions"
            optionLabel="label"
            optionValue="value"
            :placeholder="t('assessments.requisition_placeholder')"
            class="!w-full"
            filter
            showClear
          />
        </FormRow>
        <FormRow :label="t('assessments.scheduled_date')">
          <DateInput v-model="form.scheduled_date" class="!w-full" />
        </FormRow>
        <FormRow :label="t('assessments.scheduled_time')">
          <TimeInput v-model="form.scheduled_time" class="!w-full" />
        </FormRow>
        <FormRow :label="t('assessments.location')">
          <TextInput v-model="form.location" class="!w-full" />
        </FormRow>
        <FormRow :label="t('assessments.meeting_link')">
          <TextInput v-model="form.meeting_link" class="!w-full" />
        </FormRow>
        <FormRow :label="t('assessments.notes')">
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

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('assessments.confirm_delete_title')"
      :message="t('assessments.confirm_delete_message')"
      :loading="deleting"
      @confirm="doDelete()"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'
import { formatDate } from '@/utils/formatDate'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
import TimeInput from '@/components/TimeInput.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)
const items = ref([])
const requisitions = ref([])
const participantCounts = ref({})
const dialogVisible = ref(false)
const editingId = ref(null)
const errors = ref(null)
const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-28', headerWidth: 'w-16' },
  { type: 'text', width: 'w-36', headerWidth: 'w-24' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-20' }
]

const typeOptions = computed(() =>
  ['TECHNICAL', 'PSYCHOLOGICAL', 'COGNITIVE', 'PERSONALITY', 'CASE_STUDY', 'CODING', 'LANGUAGE', 'OTHER'].map(v => ({
    label: t(`assessments.type_${v.toLowerCase()}`),
    value: v
  }))
)

const requisitionOptions = computed(() =>
  requisitions.value.map(r => ({ label: r.title, value: r.id }))
)

function typeLabel(type) {
  if (!type) return '-'
  const key = `assessments.type_${type.toLowerCase()}`
  return t(key) !== key ? t(key) : type
}

function typeSeverity(type) {
  switch (type) {
    case 'TECHNICAL': return 'info'
    case 'PSYCHOLOGICAL': case 'PERSONALITY': return 'help'
    case 'COGNITIVE': case 'CODING': case 'LANGUAGE': return 'warn'
    case 'CASE_STUDY': return 'secondary'
    default: return 'secondary'
  }
}

function formatScheduledAt(value) {
  if (!value) return '-'
  const ms = Number(value) / 1000000
  if (!Number.isFinite(ms) || ms <= 0) return '-'
  const d = new Date(ms)
  const datePart = formatDate(d, locale.value)
  const time = d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

function requisitionTitle(id) {
  if (!id) return '-'
  const r = requisitions.value.find(x => x.id === id)
  return r ? r.title : id
}

const defaultForm = () => ({
  name: '',
  type: null,
  requisition_id: null,
  scheduled_date: '',
  scheduled_time: '',
  location: '',
  meeting_link: '',
  notes: ''
})

const form = ref(defaultForm())
const editing = computed(() => !!editingId.value)

function openDialog(row) {
  editingId.value = row?.id || null
  errors.value = null
  if (row) {
    form.value = {
      name: row.name || '',
      type: row.type || null,
      requisition_id: row.requisition_id || null,
      scheduled_date: row.scheduled_at ? new Date(Number(row.scheduled_at) / 1000000).toISOString().slice(0, 10) : '',
      scheduled_time: row.scheduled_at ? new Date(Number(row.scheduled_at) / 1000000).toTimeString().slice(0, 8) : '',
      location: row.location || '',
      meeting_link: row.meeting_link || '',
      notes: row.notes || ''
    }
  } else {
    form.value = defaultForm()
  }
  dialogVisible.value = true
}

async function loadData() {
  loading.value = true
  try {
    const [assessRes, reqRes] = await Promise.all([
      api.get('/api/v1/tenant/recruitment/assessments'),
      api.get('/api/v1/tenant/recruitment/requisitions', { params: { per_page: 500 } })
    ])
    items.value = assessRes.data?.data || []
    requisitions.value = reqRes.data?.data || []
    await loadParticipantCounts()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadParticipantCounts() {
  const counts = {}
  const results = await Promise.allSettled(
    items.value.map(a => api.get(`/api/v1/tenant/recruitment/assessments/${a.id}/participants`))
  )
  results.forEach((res, i) => {
    if (res.status === 'fulfilled') {
      counts[items.value[i].id] = (res.value.data?.data || []).length
    }
  })
  participantCounts.value = counts
}

function buildPayload() {
  let scheduledAt = 0
  if (form.value.scheduled_date) {
    const dateStr = form.value.scheduled_time
      ? `${form.value.scheduled_date}T${form.value.scheduled_time}`
      : `${form.value.scheduled_date}T00:00:00`
    scheduledAt = new Date(dateStr).getTime() * 1000000
  }
  const payload = {
    name: form.value.name,
    type: form.value.type || undefined,
    requisition_id: form.value.requisition_id || undefined,
    scheduled_at: scheduledAt || undefined,
    location: form.value.location || undefined,
    meeting_link: form.value.meeting_link || undefined,
    notes: form.value.notes || undefined
  }
  Object.keys(payload).forEach(k => {
    if (payload[k] === undefined || payload[k] === null || payload[k] === '') delete payload[k]
  })
  return payload
}

async function save() {
  errors.value = null
  if (!form.value.name?.trim()) {
    errors.value = { name: [t('form.required')] }
    return
  }
  saving.value = true
  try {
    const payload = buildPayload()
    if (editing.value) {
      await api.put(`/api/v1/tenant/recruitment/assessments/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/recruitment/assessments', payload)
    }
    dialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: editing.value ? t('assessments.updated') : t('assessments.created'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 4000 })
  } finally {
    saving.value = false
  }
}

function confirmDelete(row) {
  deleteTarget.value = row
  deleteDialogVisible.value = true
}

async function doDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.delete(`/api/v1/tenant/recruitment/assessments/${deleteTarget.value.id}`)
    deleteDialogVisible.value = false
    deleteTarget.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('assessments.deleted'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
