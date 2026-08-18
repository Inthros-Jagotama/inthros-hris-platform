<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Select
          v-model="filterCompetencyId"
          :options="competencyOptions"
          optionLabel="label"
          optionValue="value"
          showClear
          filter
          class="w-56"
          panelClass="!whitespace-nowrap"
          :placeholder="t('competency_360.select_competency')"
          @change="reload"
        />
        <Button :label="t('competency_360.new_indicator')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
          <i class="pi pi-list text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.indicators_empty') }}</p>
        </div>
      </template>
      <Column field="competency_name" :header="t('competency_360.competency')" style="width:260px">
        <template #body="{data}">
          <span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.competency_name || '-' }}</span>
        </template>
      </Column>
      <Column field="code" :header="t('competency_360.code')" style="width:120px">
        <template #body="{data}"><Tag :value="data.code || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="statement" :header="t('competency_360.statement')">
        <template #body="{data}">
          <span class="text-gray-600 dark:text-gray-300">{{ data.statement }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'active' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('competency_360.edit_indicator') : t('competency_360.new_indicator')" modal :style="{ width: '620px' }" @hide="resetForm">
      <div class="space-y-3">
        <FormRow :label="t('competency_360.competency')" required :errors="errors?.competency_id">
          <Select
            v-model="form.competency_id"
            :options="competencyOptions"
            optionLabel="label"
            optionValue="value"
            filter
            class="w-full"
            panelClass="!whitespace-nowrap"
            :placeholder="t('competency_360.select_competency')"
            :class="{ 'p-invalid': errors?.competency_id }"
          />
        </FormRow>
        <FormRow :label="t('competency_360.statement')" required :errors="errors?.statement">
          <TextInput v-model="form.statement" textarea :rows="3" maxlength="1000" :placeholder="t('competency_360.statement')" :class="{ 'p-invalid': errors?.statement }" />
        </FormRow>
        <FormRow :label="t('common.description')" :errors="errors?.description">
          <TextInput v-model="form.description" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('common.status')">
          <Select v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible = false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('competency_360.delete_indicator_title')"
      :message="t('competency_360.delete_indicator', { statement: deleteTarget?.statement || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const filterCompetencyId = ref(null)

const competencyOptions = ref([])

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref(defaultForm())

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const statusOptions = [
  { label: t('common_status.active'), value: 'active' },
  { label: t('common_status.inactive'), value: 'inactive' }
]

const skeletonColumns = [
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-56', headerWidth: 'w-32' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function defaultForm() {
  return { competency_id: null, statement: '', description: '', status: 'active' }
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadCompetencies() {
  try {
    const res = await api.get('/api/v1/tenant/competency/competencies', { params: { per_page: 200 } })
    competencyOptions.value = (res.data?.data || []).map(c => ({ label: c.name, value: c.id }))
  } catch {
    competencyOptions.value = []
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    if (filterCompetencyId.value) params.competency_id = filterCompetencyId.value
    const res = await api.get('/api/v1/tenant/competency/indicators', { params })
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = item
    ? {
        competency_id: item.competency_id || null,
        statement: item.statement || '',
        description: item.description || '',
        status: item.status || 'active'
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
  if (!form.value.competency_id) { errors.value = { competency_id: t('form.required') }; return }
  if (!form.value.statement?.trim()) { errors.value = { statement: t('form.required') }; return }
  saving.value = true
  try {
    const payload = {
      competency_id: form.value.competency_id,
      statement: form.value.statement.trim(),
      description: form.value.description?.trim() || '',
      status: form.value.status || 'active'
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/competency/indicators/${editingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/competency/indicators', payload)
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

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/competency/indicators/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await loadCompetencies()
  await loadData()
})
</script>
