<template>
  <div class="space-y-4">
    <div class="flex items-center gap-2 flex-wrap">
      <IconField>
        <InputIcon class="pi pi-search" />
        <InputText v-model="searchQuery" :placeholder="t('candidates.search')" size="small" class="!pl-8 !text-sm !py-1.5 !w-64" @input="onSearchInput" />
      </IconField>
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('candidates.new_candidate')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-users text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('candidates.empty') }}</p>
        </div>
      </template>

      <Column :header="t('common.name')">
        <template #body="{ data }">
          <button type="button" class="flex items-center gap-2.5 min-w-0 text-left group" @click="router.push(`/recruitment/candidates/${data.id}`)">
            <div class="w-8 h-8 rounded-full bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-[11px] font-semibold shrink-0">
              {{ initials(data) }}
            </div>
            <div class="min-w-0">
              <p class="font-medium text-navy-800 dark:text-gray-100 truncate group-hover:text-sky-600 dark:group-hover:text-sky-400">{{ data.first_name }} {{ data.last_name }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 truncate">{{ data.email }}</p>
            </div>
          </button>
        </template>
      </Column>

      <Column :header="t('candidates.candidate_number')" style="width: 170px">
        <template #body="{ data }">
          <span class="text-xs font-mono text-gray-600 dark:text-gray-300">{{ data.candidate_number || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('candidates.current_title')">
        <template #body="{ data }">
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ data.current_title || '—' }}</span>
          <span v-if="data.current_company" class="text-xs text-gray-400 dark:text-gray-500"> · {{ data.current_company }}</span>
        </template>
      </Column>

      <Column :header="t('candidates.candidate_type')" style="width: 110px">
        <template #body="{ data }">
          <Tag
            :value="t('candidates.type_' + (data.candidate_type || 'external').toLowerCase())"
            :severity="data.candidate_type === 'INTERNAL' ? 'info' : 'secondary'"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>

      <Column :header="t('candidates.source')" style="width: 140px">
        <template #body="{ data }">
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ data.source || '—' }}</span>
        </template>
      </Column>

      <Column :header="t('common.actions')" :exportable="false" style="width: 110px">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" text size="small" class="!w-7 !h-7" @click="openDialog(data)" />
            <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('common.edit') : t('candidates.new_candidate')" :modal="true" class="!w-[min(95vw,640px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('candidates.first_name')" :required="true">
          <TextInput v-model="form.first_name" class="!w-full" />
        </FormRow>
        <FormRow :label="t('candidates.last_name')" :required="true">
          <TextInput v-model="form.last_name" class="!w-full" />
        </FormRow>

        <FormRow :label="t('candidates.email')" :required="true">
          <TextInput v-model="form.email" class="!w-full" />
        </FormRow>
        <FormRow :label="t('candidates.phone')">
          <TextInput v-model="form.phone" class="!w-full" />
        </FormRow>

        <FormRow :label="t('candidates.current_company')">
          <TextInput v-model="form.current_company" class="!w-full" />
        </FormRow>
        <FormRow :label="t('candidates.current_title')">
          <TextInput v-model="form.current_title" class="!w-full" />
        </FormRow>

        <FormRow :label="t('candidates.source')">
          <TextInput v-model="form.source" :placeholder="t('candidates.source_placeholder')" class="!w-full" />
        </FormRow>
        <FormRow :label="t('candidates.linkedin_url')">
          <TextInput v-model="form.linkedin_url" class="!w-full" />
        </FormRow>

        <FormRow :label="t('candidates.resume_url')">
          <TextInput v-model="form.resume_url" class="!w-full" />
        </FormRow>
        <FormRow :label="t('candidates.portfolio_url')">
          <TextInput v-model="form.portfolio_url" class="!w-full" />
        </FormRow>

        <FormRow :label="t('candidates.address')">
          <Textarea v-model="form.address" :rows="2" class="!w-full" />
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

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('common.delete')"
      :message="t('candidates.delete_confirm', { name: pendingDelete ? `${pendingDelete.first_name} ${pendingDelete.last_name}` : '' })"
      :loading="deleting"
      @confirm="doDelete()"
    />
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
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Textarea from 'primevue/textarea'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const deleting = ref(false)
const items = ref([])
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(10)
const searchQuery = ref('')
const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const editing = ref(false)
const pendingDelete = ref(null)
let searchTimer = null

const skeletonColumns = [
  { field: 'name', header: 'Name', width: '28%' },
  { field: 'number', header: 'No.', width: '14%' },
  { field: 'title', header: 'Title', width: '22%' },
  { field: 'type', header: 'Type', width: '10%' },
  { field: 'source', header: 'Source', width: '14%' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

const emptyForm = () => ({
  id: null,
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  address: '',
  current_company: '',
  current_title: '',
  resume_url: '',
  portfolio_url: '',
  linkedin_url: '',
  source: '',
  notes: ''
})

const form = ref(emptyForm())

function initials(data) {
  const a = (data.first_name || '?').charAt(0)
  const b = (data.last_name || '').charAt(0)
  return (a + b).toUpperCase()
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const q = searchQuery.value?.trim()
    if (q) params.search = q
    const res = await api.get('/api/v1/tenant/recruitment/candidates', { params })
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

function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadData()
  }, 400)
}

function openDialog(row) {
  if (row) {
    editing.value = true
    form.value = { ...emptyForm(), ...row }
  } else {
    editing.value = false
    form.value = emptyForm()
  }
  dialogVisible.value = true
}

async function save() {
  if (!form.value.first_name.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.first_name_required'), life: 4000 })
    return
  }
  if (!form.value.last_name.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.last_name_required'), life: 4000 })
    return
  }
  if (!form.value.email.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.email_required'), life: 4000 })
    return
  }
  saving.value = true
  try {
    const payload = { ...form.value }
    delete payload.id
    Object.keys(payload).forEach(k => {
      if (payload[k] === '' || payload[k] === null || payload[k] === undefined) delete payload[k]
    })
    if (editing.value) {
      await api.put(`/api/v1/tenant/recruitment/candidates/${form.value.id}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/recruitment/candidates', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.created'), life: 3000 })
    }
    dialogVisible.value = false
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    saving.value = false
  }
}

function confirmDelete(row) {
  pendingDelete.value = row
  deleteDialogVisible.value = true
}

async function doDelete() {
  if (!pendingDelete.value) return
  deleting.value = true
  try {
    await api.delete(`/api/v1/tenant/recruitment/candidates/${pendingDelete.value.id}`)
    deleteDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.deleted'), life: 3000 })
    loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    deleting.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

onMounted(() => {
  loadData()
})
</script>
