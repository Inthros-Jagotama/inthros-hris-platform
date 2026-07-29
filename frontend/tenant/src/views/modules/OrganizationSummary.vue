<template>
  <div class="space-y-1">
    <!-- Header -->
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-1.5">
        <Button
          v-for="chip in filterChips"
          :key="chip.value"
          :label="chip.label"
          :severity="activeFilter === chip.value ? (chip.severity || 'secondary') : 'secondary'"
          :outlined="activeFilter !== chip.value"
          size="small"
          class="!text-xs !px-2 !py-1"
          @click="activeFilter = chip.value"
        />
      </div>
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
        <Button :label="t('org_summary.new')" icon="pi pi-plus" size="small" @click="openCreateDialog()" />
      </div>
    </div>

    <!-- DataTable -->
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="clientFiltered"
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
          <i class="pi pi-building text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('org_summary.empty_title') }}</p>
          <p class="text-xs mt-1">{{ t('org_summary.empty_hint') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('org_summary.code')" sortable style="width:120px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !font-mono !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="decree_no" :header="t('org_summary.decree_no')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.decree_no }}</span></template>
      </Column>
      <Column field="decree_date" :header="t('org_summary.decree_date')" sortable style="width:140px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.decree_date }}</span></template>
      </Column>
      <Column field="status" :header="t('common.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag
            :value="data.status === 'active' ? t('common_status.active') : data.status"
            :severity="data.status === 'active' ? 'success' : 'warn'"
            class="!text-xs !px-1.5 !py-0.5"
          />
        </template>
      </Column>
      <Column field="org_count" :header="t('org_summary.orgs')" sortable style="width:100px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.org_count || 0 }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-sitemap" size="small" text severity="info" v-tooltip.left="t('org_summary.go_to_tree')" @click="goToOrgTree(data.id)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEditDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('org_summary.edit') : t('org_summary.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div class="space-y-3">
            <FormRow :label="t('org_summary.code')" required :errors="errors?.code">
              <TextInput v-model="form.code" :class="{ 'p-invalid': errors?.code }" maxlength="7" autofocus :placeholder="t('org_summary.code_placeholder')" />
            </FormRow>
            <FormRow :label="t('org_summary.decree_no')" required :errors="errors?.decree_no">
              <TextInput v-model="form.decree_no" :class="{ 'p-invalid': errors?.decree_no }" maxlength="20" :placeholder="t('org_summary.decree_no_placeholder')" />
            </FormRow>
            <FormRow :label="t('org_summary.decree_date')" required :errors="errors?.decree_date">
              <DateInput v-model="form.decree_date" :class="{ 'p-invalid': errors?.decree_date }" :placeholder="t('org_summary.decree_date_placeholder')" />
            </FormRow>
          </div>
        </div>
      <template #footer>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 ml-auto">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
            <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
          </div>
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('org_summary.confirm_delete_title')"
      :message="t('org_summary.confirm_delete', { code: deleteTarget?.code || '', decree: deleteTarget?.decree_no || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      :confirm-label="t('common.delete')"
      :cancel-label="t('common.cancel')"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import DateInput from '@/components/DateInput.vue'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const router = useRouter()
const toast = useToast()
const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const activeFilter = ref(null)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const form = ref({ code: '', decree_no: '', decree_date: null })

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-36', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-20' }
]

const filterChips = computed(() => [
  { label: t('common.all'), value: null, severity: 'info' },
  { label: t('common_status.active'), value: 'active', severity: 'success' },
  { label: t('common_status.inactive'), value: 'inactive', severity: 'warn' }
])

const clientFiltered = computed(() => {
  let result = items.value
  if (activeFilter.value === 'active') {
    result = result.filter(i => i.status === 'active')
  } else if (activeFilter.value === 'inactive') {
    result = result.filter(i => i.status !== 'active')
  }
  return result
})

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/organization-summaries', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

function openCreateDialog() {
  editing.value = false
  editingId.value = null
  errors.value = {}
  form.value = { code: '', decree_no: '', decree_date: null }
  dialogVisible.value = true
}

function openEditDialog(item) {
  editing.value = true
  editingId.value = item.id
  errors.value = {}
  form.value = {
    code: item.code,
    decree_no: item.decree_no,
    decree_date: item.decree_date ? new Date(item.decree_date + 'T00:00:00') : null
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { code: '', decree_no: '', decree_date: null }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

function goToOrgTree(summaryId) {
  if (summaryId) {
    router.push(`/organizations?summary_id=${summaryId}`)
  } else {
    router.push('/organizations')
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
    await api.delete(`/api/v1/tenant/organization-summaries/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('org_summary.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}


async function handleSave() {
  errors.value = {}
  if (!form.value.code?.trim()) { errors.value = { code: [t('form.required')] }; return }
  if (!form.value.decree_no?.trim()) { errors.value = { decree_no: [t('form.required')] }; return }
  if (!form.value.decree_date) { errors.value = { decree_date: [t('form.required')] }; return }

  saving.value = true
  try {
    const payload = {
      code: form.value.code,
      decree_no: form.value.decree_no,
      decree_date: form.value.decree_date instanceof Date
        ? form.value.decree_date.toISOString().split('T')[0]
        : form.value.decree_date
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/organization-summaries/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('org_summary.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/organization-summaries', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('org_summary.created'), life: 3000 })
    }

    dialogVisible.value = false
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
