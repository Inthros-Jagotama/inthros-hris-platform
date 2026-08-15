<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <Button :label="t('payroll.new_profile')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-users text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.profiles_empty') }}</p></div></template>
      <Column field="employee_id" :header="t('payroll.employee')" sortable style="width:200px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeLabel(data.employee_id) }}</span></template></Column>
      <Column field="payroll_group_code" :header="t('payroll.payroll_group_code')" sortable style="width:150px"><template #body="{data}"><Tag :value="data.payroll_group_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column field="payroll_frequency" :header="t('payroll.payroll_frequency')" sortable style="width:110px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ freqLabel(data.payroll_frequency) }}</span></template></Column>
      <Column field="payment_method" :header="t('payroll.payment_method')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ methodLabel(data.payment_method) }}</span></template></Column>
      <Column field="effective_start_date" :header="t('payroll.effective_start_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
      <Column field="is_payroll_active" :header="t('payroll.is_payroll_active')" sortable style="width:110px"><template #body="{data}"><Tag :value="data.is_payroll_active ? t('common.yes') : t('common.no')" :severity="data.is_payroll_active ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:80px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="t('payroll.new_profile')" modal :style="{ width: '560px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.employee')" required :errors="errors?.employee_id">
          <SelectLabel v-model="form.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{'p-invalid':errors?.employee_id}" showClear />
        </FormRow>
        <FormRow :label="t('payroll.payroll_group_code')" required :errors="errors?.payroll_group_code">
          <TextInput v-model="form.payroll_group_code" maxlength="50" :placeholder="t('payroll.payroll_group_code')" :class="{'p-invalid':errors?.payroll_group_code}" />
        </FormRow>
        <FormRow :label="t('payroll.payroll_frequency')">
          <SelectLabel v-model="form.payroll_frequency" :options="freqOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('payroll.payment_method')">
          <SelectLabel v-model="form.payment_method" :options="methodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('payroll.salary_currency')">
          <TextInput v-model="form.salary_currency" maxlength="3" placeholder="IDR" />
        </FormRow>
        <FormRow :label="t('payroll.is_payroll_active')">
          <ToggleSwitch v-model="form.is_payroll_active" />
        </FormRow>
        <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
          <DateInput v-model="form.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
        </FormRow>
        <FormRow :label="t('payroll.effective_end_date')">
          <DateInput v-model="form.effective_end_date" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.employee_profiles')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import DateInput from '@/components/DateInput.vue'
const { t, locale } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const employees = ref([])
const form = ref({ employee_id: null, payroll_group_code: '', payroll_frequency: 'MONTHLY', payment_method: 'BANK_TRANSFER', salary_currency: 'IDR', is_payroll_active: true, effective_start_date: '', effective_end_date: '' })

const freqOptions = computed(() => ['MONTHLY', 'WEEKLY', 'DAILY'].map(v => ({ label: t(`payroll.payroll_frequency_${v.toLowerCase()}`), value: v })))
const methodOptions = computed(() => ['BANK_TRANSFER', 'CASH', 'CHEQUE'].map(v => ({ label: t(`payroll.payment_method_${v.toLowerCase()}`), value: v })))
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const skeletonColumns = [{type:'text',width:'w-40',headerWidth:'w-20'},{type:'tag',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-20',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'tag',width:'w-16',headerWidth:'w-16'},{type:'icons',count:1,headerWidth:'w-16'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.employee')}: ${employeeLabel(deleteTarget.value.employee_id)}` : t('common.no_data'))

function freqLabel(v) { const key = `payroll.payroll_frequency_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function methodLabel(v) { const key = `payroll.payment_method_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function employeeLabel(id) {
  const e = employees.value.find(x => x.id === id)
  return e ? `${e.name} (${e.employee_id})` : id
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/employee-payroll-profiles', { params: { page: currentPage.value, per_page: perPage.value } })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally { loading.value = false }
}
async function loadEmployees() {
  try {
    const res = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
    employees.value = res.data?.data || []
  } catch { employees.value = [] }
}
function onPage(event) { currentPage.value = event.page + 1; perPage.value = event.rows; loadData() }
function openDialog() {
  errors.value = {}
  form.value = { employee_id: null, payroll_group_code: '', payroll_frequency: 'MONTHLY', payment_method: 'BANK_TRANSFER', salary_currency: 'IDR', is_payroll_active: true, effective_start_date: '', effective_end_date: '' }
  dialogVisible.value = true
}
function resetForm() {
  errors.value = {}
  form.value = { employee_id: null, payroll_group_code: '', payroll_frequency: 'MONTHLY', payment_method: 'BANK_TRANSFER', salary_currency: 'IDR', is_payroll_active: true, effective_start_date: '', effective_end_date: '' }
}
async function handleSave() {
  errors.value = {}
  if (!form.value.employee_id) { errors.value = { employee_id: [t('form.required')] }; return }
  if (!form.value.payroll_group_code?.trim()) { errors.value = { payroll_group_code: [t('form.required')] }; return }
  if (!form.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  saving.value = true
  try {
    await api.post('/api/v1/tenant/payroll/employee-payroll-profiles', {
      employee_id: form.value.employee_id,
      payroll_group_code: form.value.payroll_group_code.trim(),
      payroll_frequency: form.value.payroll_frequency,
      payment_method: form.value.payment_method,
      salary_currency: form.value.salary_currency || 'IDR',
      is_payroll_active: form.value.is_payroll_active,
      effective_start_date: form.value.effective_start_date,
      effective_end_date: form.value.effective_end_date || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.profile_created'), life: 3000 })
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
    await api.delete(`/api/v1/tenant/payroll/employee-payroll-profiles/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.profile_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadData()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadData(); loadEmployees() })
</script>
