<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
      <Button :label="t('payroll.new_tax_profile')" icon="pi pi-plus" size="small" @click="openDialog()" />
    </div>
    <DataTable
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
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-receipt text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.tax_profiles_empty') }}</p></div></template>
      <Column :header="t('payroll.employee')" sortable style="width:180px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeLabel(data.employee_id) }}</span></template></Column>
      <Column :header="t('payroll.npwp')" style="width:150px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ data.npwp || '-' }}</span></template></Column>
      <Column :header="t('payroll.ptkp_status')" style="width:100px"><template #body="{data}"><Tag :value="data.ptkp_status || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('payroll.tax_method')" style="width:110px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ taxMethodLabel(data.tax_method) }}</span></template></Column>
      <Column :header="t('payroll.has_npwp')" style="width:90px"><template #body="{data}"><Tag :value="data.has_npwp ? t('common.yes') : t('common.no')" :severity="data.has_npwp ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
      <Column :header="t('common.status')" style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('payroll.edit_tax_profile') : t('payroll.new_tax_profile')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.employee')" required :errors="errors?.employee_id">
          <SelectLabel v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{'p-invalid':errors?.employee_id}" :disabled="editing" @update:modelValue="onEmployeeChange" />
        </FormRow>
        <FormRow :label="t('payroll.payroll_profile')" required :errors="errors?.employee_payroll_profile_id">
          <SelectLabel v-model="form.employee_payroll_profile_id" :options="profileOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.employee_payroll_profile_id}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.npwp')">
            <TextInput v-model="form.npwp" maxlength="50" :placeholder="t('payroll.npwp_placeholder')" />
          </FormRow>
          <FormRow :label="t('payroll.npwp_registered_name')">
            <TextInput v-model="form.npwp_registered_name" maxlength="255" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.ptkp_status')">
            <SelectLabel v-model="form.ptkp_status" :options="ptkpOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
          </FormRow>
          <FormRow :label="t('payroll.tax_method')">
            <SelectLabel v-model="form.tax_method" :options="taxMethodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.has_npwp')">
            <ToggleSwitch v-model="form.has_npwp" />
          </FormRow>
          <FormRow :label="t('payroll.is_taxable')">
            <ToggleSwitch v-model="form.is_taxable" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="form.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="form.effective_end_date" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.tax_profiles')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import DateInput from '@/components/DateInput.vue'
const { t, locale } = useI18n(); const toast = useToast()
const items = ref([]); const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const employees = ref([]); const payrollProfiles = ref([])
const form = ref({ employee_id: null, employee_payroll_profile_id: null, npwp: '', npwp_registered_name: '', ptkp_status: null, tax_method: 'GROSS', is_taxable: true, has_npwp: false, effective_start_date: '', effective_end_date: null, status: 'ACTIVE', notes: '' })

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const profileOptions = computed(() => payrollProfiles.value
  .filter(p => p.employee_id === form.value.employee_id)
  .map(p => ({ label: `${p.payroll_group_code} — ${p.effective_start_date}`, value: p.id })))
const ptkpOptions = computed(() => ['TK/0', 'TK/1', 'TK/2', 'TK/3', 'K/0', 'K/1', 'K/2', 'K/3'].map(v => ({ label: v, value: v })))
const taxMethodOptions = computed(() => ['GROSS', 'GROSS_UP', 'NETT'].map(v => ({ label: taxMethodLabel(v), value: v })))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: statusLabel(v), value: v })))
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.employee')}: ${employeeLabel(deleteTarget.value.employee_id)}` : t('common.no_data'))

function statusLabel(v) { const key = `payroll.status_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function taxMethodLabel(v) { const key = `payroll.tax_method_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function employeeLabel(id) { const e = employees.value.find(x => x.id === id); return e ? `${e.name} (${e.employee_id})` : (id || '-') }

async function loadData() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/employee-tax-profiles', { params: { page: currentPage.value, per_page: perPage.value } })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
  } catch(e) { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 }) }
}
async function loadLookups() {
  try { const e = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }); employees.value = e.data?.data || [] } catch { employees.value = [] }
  try {
    const all = []
    let page = 1
    while (true) {
      const res = await api.get('/api/v1/tenant/payroll/employee-payroll-profiles', { params: { page, per_page: 100 } })
      const rows = res.data?.data || []
      all.push(...rows)
      if (!rows.length || all.length >= (res.data?.total || 0)) break
      page++
    }
    payrollProfiles.value = all
  } catch { payrollProfiles.value = [] }
}
function onPage(e) { currentPage.value = e.page + 1; perPage.value = e.rows; loadData() }
function onEmployeeChange() { form.value.employee_payroll_profile_id = null }
function openDialog(item) {
  editing.value = !!item; editingId.value = item?.id || null; errors.value = {}
  form.value = {
    employee_id: item?.employee_id || null, employee_payroll_profile_id: item?.employee_payroll_profile_id || null,
    npwp: item?.npwp || '', npwp_registered_name: item?.npwp_registered_name || '', ptkp_status: item?.ptkp_status || null,
    tax_method: item?.tax_method || 'GROSS', is_taxable: item?.is_taxable ?? true, has_npwp: item?.has_npwp ?? false,
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || null,
    status: item?.status || 'ACTIVE', notes: item?.notes || ''
  }
  dialogVisible.value = true
}
function resetForm() {
  form.value = { employee_id: null, employee_payroll_profile_id: null, npwp: '', npwp_registered_name: '', ptkp_status: null, tax_method: 'GROSS', is_taxable: true, has_npwp: false, effective_start_date: '', effective_end_date: null, status: 'ACTIVE', notes: '' }
  errors.value = {}; editing.value = false; editingId.value = null
}
async function handleSave() {
  errors.value = {}
  if (!form.value.employee_id) { errors.value = { employee_id: [t('form.required')] }; return }
  if (!form.value.employee_payroll_profile_id) { errors.value = { employee_payroll_profile_id: [t('form.required')] }; return }
  if (!form.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      employee_id: form.value.employee_id, employee_payroll_profile_id: form.value.employee_payroll_profile_id,
      npwp: form.value.npwp || null, npwp_registered_name: form.value.npwp_registered_name || null,
      ptkp_status: form.value.ptkp_status || null, tax_method: form.value.tax_method,
      is_taxable: form.value.is_taxable, has_npwp: form.value.has_npwp,
      effective_start_date: form.value.effective_start_date, effective_end_date: form.value.effective_end_date || null,
      status: form.value.status, notes: form.value.notes || null
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/payroll/employee-tax-profiles/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.tax_profile_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/employee-tax-profiles', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.tax_profile_created'), life: 3000 })
    }
    dialogVisible.value = false; await loadData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { saving.value = false }
}
function confirmDelete(item) { deleteTarget.value = item; deleteError.value = ''; deleteDialogVisible.value = true }
async function handleDelete() {
  deleting.value = true; deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/payroll/employee-tax-profiles/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.tax_profile_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadData()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadData(); loadLookups() })
</script>
