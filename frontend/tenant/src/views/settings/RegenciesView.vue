<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('regencies.new_regency')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      sortField="code"
      :sortOrder="1"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-map text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('regencies.empty_title') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('regencies.code')" sortable style="width:120px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('regencies.name')" sortable>
        <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="province_name" :header="t('regencies.province')" sortable style="width:200px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ getProvinceName(data.province_id) }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('regencies.edit_regency') : t('regencies.new_regency')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5">
            <i class="pi pi-map text-indigo-400 text-sm"></i>
            {{ editing ? t('regencies.edit_regency') : t('regencies.new_regency') }}
          </h3>
          <div class="space-y-2">
            <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('regencies.code') }} <span class="text-red-500">*</span></label><InputText v-model="form.code" class="!w-full" :class="{'p-invalid':errors?.code}" maxlength="10" autofocus :placeholder="t('regencies.code')" /><small v-if="errors?.code" class="text-red-500 text-xs mt-1 block">{{ errors.code }}</small></div>
            <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('regencies.name') }} <span class="text-red-500">*</span></label><InputText v-model="form.name" class="!w-full" :class="{'p-invalid':errors?.name}" maxlength="255" :placeholder="t('regencies.name')" /><small v-if="errors?.name" class="text-red-500 text-xs mt-1 block">{{ errors.name }}</small></div>
            <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('regencies.province') }} <span class="text-red-500">*</span></label>
              <Select v-model="form.province_id" :options="provinceOptions" optionValue="id" optionLabel="label" :placeholder="t('regencies.select_province')" class="!w-full" :class="{'p-invalid':errors?.province_id}" :showClear="true" />
              <small v-if="errors?.province_id" class="text-red-500 text-xs mt-1 block">{{ errors.province_id }}</small>
            </div>
          </div>
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

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import ConfirmDialog from 'primevue/confirmdialog'
import SkeletonTable from '@/components/SkeletonTable.vue'

const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ code: '', name: '', province_id: null })
const allProvinces = ref([])

const provinceOptions = computed(() => allProvinces.value.map(p => ({ id: p.id, label: `${p.code} — ${p.name}` })))

function getProvinceName(id) {
  const p = allProvinces.value.find(x => x.id === id)
  return p ? p.name : '—'
}

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-36', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/regencies', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
    if (allProvinces.value.length === 0) {
      const r = await api.get('/api/v1/tenant/settings/provinces?all=true')
      allProvinces.value = r.data?.data || r.data || []
    }
  } catch(e) {
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

function openDialog(item) {
  editing.value = !!item
  editingId.value = item?.id || null
  errors.value = {}
  form.value = { code: item?.code || '', name: item?.name || '', province_id: item?.province_id || null }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { code: '', name: '', province_id: null }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.code?.trim()) { errors.value = { code: [t('form.required')] }; return }
  if (!form.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  if (!form.value.province_id) { errors.value = { province_id: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = { code: form.value.code, name: form.value.name, province_id: form.value.province_id }
    if (editing.value) {
      await api.put(`/api/v1/tenant/settings/regencies/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('regencies.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/settings/regencies', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('regencies.created'), life: 3000 })
    }
    dialogVisible.value = false
    await loadData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  confirm.require({
    header: t('regencies.confirm_delete_title'),
    message: t('regencies.confirm_delete', { name: item.name }),
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: t('common.cancel'),
    acceptLabel: t('common.delete'),
    rejectClass: 'p-button-outlined p-button-secondary',
    acceptClass: 'p-button-danger',
    accept: async () => {
      try {
        await api.delete(`/api/v1/tenant/settings/regencies/${item.id}`)
        toast.add({ severity: 'success', summary: t('message.success'), detail: t('regencies.deleted'), life: 3000 })
        await loadData()
      } catch(e) {
        toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
      }
    }
  })
}

onMounted(loadData)
</script>
