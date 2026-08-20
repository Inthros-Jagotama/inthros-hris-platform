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
        <Button :label="t('zones.new_zone')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      sortField="sort_order"
      :sortOrder="1"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-map-marker text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('zones.empty_title') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('zones.code')" sortable style="width:120px">
        <template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('zones.name')" sortable>
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column field="region" :header="t('zones.region')" sortable style="width:150px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.region || '—' }}</span></template>
      </Column>
      <Column field="is_active" :header="t('common.status')" sortable style="width:100px">
        <template #body="{data}">
          <Tag :value="data.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="data.is_active ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="sort_order" :header="t('zones.sort_order')" sortable style="width:100px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span></template>
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

    <!-- Create/Edit Dialog -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('zones.edit_zone') : t('zones.new_zone')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div class="space-y-2">
            <FormRow :label="t('zones.code')" required :errors="errors?.code">
            <TextInput v-model="form.code" maxlength="20" autofocus :placeholder="t('zones.code')" :class="{'p-invalid':errors?.code}" />
          </FormRow>
            <FormRow :label="t('zones.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" :placeholder="t('zones.name')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
            <FormRow :label="t('zones.region')" :errors="errors?.region">
            <TextInput v-model="form.region" maxlength="100" :placeholder="t('zones.region')" />
          </FormRow>
            <FormRow :label="t('zones.timezone')" :errors="errors?.timezone">
            <Select v-model="form.timezone" :options="timezoneOverrideOptions" optionLabel="label" optionValue="value" class="w-full" />
          </FormRow>
            <div class="flex items-center justify-between pt-2"><label class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('zones.is_active') }}</label><ToggleSwitch v-model="form.is_active" /></div>
            <FormRow :label="t('zones.sort_order')" :errors="errors?.sort_order">
            <InputNumber v-model="form.sort_order" class="!w-full" :min="0" size="small" />
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
      :title="t('zones.confirm_delete_title')"
      :message="t('zones.confirm_delete', { name: deleteTarget?.name || '' })"
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
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import ToggleSwitch from 'primevue/toggleswitch'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const activeFilter = ref(null)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const form = ref({ code: '', name: '', region: '', timezone: null, is_active: true, sort_order: 0 })

const timezoneOverrideOptions = computed(() => [
  { label: t('zones.timezone_inherit'), value: null },
  { label: 'WIB (Asia/Jakarta)', value: 'Asia/Jakarta' },
  { label: 'WITA (Asia/Makassar)', value: 'Asia/Makassar' },
  { label: 'WIT (Asia/Jayapura)', value: 'Asia/Jayapura' }
])
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-16' },
  { type: 'text', width: 'w-32', headerWidth: 'w-16' },
  { type: 'text', width: 'w-24', headerWidth: 'w-16' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-12', headerWidth: 'w-16' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const filterChips = computed(() => [
  { label: t('common.all'), value: null, severity: 'info' },
  { label: t('common_status.active'), value: 'active', severity: 'success' },
  { label: t('common_status.inactive'), value: 'inactive', severity: 'warn' }
])

const clientFiltered = computed(() => {
  let result = items.value
  if (activeFilter.value === 'active') {
    result = result.filter(i => i.is_active === true)
  } else if (activeFilter.value === 'inactive') {
    result = result.filter(i => i.is_active === false)
  }
  return result
})

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/zones', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
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
  form.value = {
    code: item?.code || '',
    name: item?.name || '',
    region: item?.region || '',
    timezone: item?.timezone ?? null,
    is_active: item?.is_active !== undefined ? item.is_active : true,
    sort_order: item?.sort_order || 0
  }
  dialogVisible.value = true
}

function resetForm() {
  form.value = { code: '', name: '', region: '', timezone: null, is_active: true, sort_order: 0 }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

async function handleSave() {
  errors.value = {}
  if (!form.value.code?.trim()) { errors.value = { code: [t('form.required')] }; return }
  if (!form.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      code: form.value.code,
      name: form.value.name,
      region: form.value.region || undefined,
      is_active: form.value.is_active,
      sort_order: form.value.sort_order || 0
    }
    if (editing.value) {
      // Update: selalu kirim "timezone" secara eksplisit. String kosong ("")
      // berarti "ikut default perusahaan" (hapus override) — backend
      // membedakan "" (clear) dari field yang di-omit (no-op, JSON null pun
      // unmarshal ke nil pointer yang sama). Jangan kirim null di sini.
      payload.timezone = form.value.timezone || ''
      await api.put(`/api/v1/tenant/settings/zones/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.updated'), life: 3000 })
    } else {
      // Create: zona baru tidak punya state sebelumnya untuk "di-clear" —
      // cukup omit field bila kosong agar backend menyimpan override sebagai
      // nil (ikut default perusahaan) dari awal.
      if (form.value.timezone) payload.timezone = form.value.timezone
      await api.post('/api/v1/tenant/settings/zones', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.created'), life: 3000 })
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
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/zones/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('zones.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
