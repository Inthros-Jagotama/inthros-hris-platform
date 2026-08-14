<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <Button :label="t('payroll.new_period')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      sortField="period_year"
      :sortOrder="-1"
    >
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-calendar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.periods_empty') }}</p></div></template>
      <Column field="period_code" :header="t('payroll.period_code')" sortable style="width:120px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.period_code }}</span></template></Column>
      <Column field="period_year" :header="t('payroll.period_year')" sortable style="width:90px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.period_year }}</span></template></Column>
      <Column field="period_month" :header="t('payroll.period_month')" sortable style="width:90px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ formatMonth(data.period_month) }}</span></template></Column>
      <Column field="start_date" :header="t('payroll.start_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.start_date, locale) }}</span></template></Column>
      <Column field="end_date" :header="t('payroll.end_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.end_date, locale) }}</span></template></Column>
      <Column field="as_of_date" :header="t('payroll.as_of_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.as_of_date, locale) }}</span></template></Column>
      <Column field="status" :header="t('common.status')" sortable style="width:100px"><template #body="{data}"><Tag :value="t('payroll.status_' + data.status.toLowerCase())" :severity="data.status === 'OPEN' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /></div></template></Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('payroll.payroll_periods') : t('payroll.new_period')" modal :style="{ width: '560px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div v-if="!editing" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.period_year')" required :errors="errors?.period_year">
            <InputNumber v-model="form.period_year" class="!w-full" :min="2000" :max="2100" size="small" :class="{'p-invalid':errors?.period_year}" />
          </FormRow>
          <FormRow :label="t('payroll.period_month')" required :errors="errors?.period_month">
            <SelectLabel v-model="form.period_month" :options="monthOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.period_month}" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.start_date')" required :errors="errors?.start_date">
            <DateInput v-model="form.start_date" :class="{'p-invalid':errors?.start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.end_date')" required :errors="errors?.end_date">
            <DateInput v-model="form.end_date" :class="{'p-invalid':errors?.end_date}" />
          </FormRow>
        </div>
        <FormRow :label="t('payroll.as_of_date')" required :errors="errors?.as_of_date">
          <DateInput v-model="form.as_of_date" :class="{'p-invalid':errors?.as_of_date}" />
        </FormRow>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></template>
    </Dialog>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate, formatMonth } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
const { t, locale } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const form = ref({ period_year: new Date().getFullYear(), period_month: null, start_date: '', end_date: '', as_of_date: '', status: 'OPEN' })

const statusOptions = computed(() => ['OPEN', 'CLOSED'].map(v => ({ label: t(`payroll.status_${v.toLowerCase()}`), value: v })))
const monthOptions = computed(() => Array.from({ length: 12 }, (_, i) => ({ label: formatMonth(i + 1, locale.value), value: i + 1 })))
const skeletonColumns = [{type:'text',width:'w-20',headerWidth:'w-16'},{type:'text',width:'w-12',headerWidth:'w-12'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'tag',width:'w-16',headerWidth:'w-16'},{type:'icons',count:1,headerWidth:'w-16'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/periods', { params: { page: currentPage.value, per_page: perPage.value } })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally { loading.value = false }
}
function onPage(event) { currentPage.value = event.page + 1; perPage.value = event.rows; loadData() }
function openDialog(item) {
  editing.value = !!item; editingId.value = item?.id || null; errors.value = {}
  form.value = {
    period_year: item?.period_year || new Date().getFullYear(),
    period_month: item?.period_month || null,
    start_date: item?.start_date || '', end_date: item?.end_date || '', as_of_date: item?.as_of_date || '',
    status: item?.status || 'OPEN'
  }
  dialogVisible.value = true
}
function resetForm() {
  form.value = { period_year: new Date().getFullYear(), period_month: null, start_date: '', end_date: '', as_of_date: '', status: 'OPEN' }
  errors.value = {}; editing.value = false; editingId.value = null
}
async function handleSave() {
  errors.value = {}
  if (!form.value.period_year) { errors.value = { period_year: [t('form.required')] }; return }
  if (!form.value.period_month) { errors.value = { period_month: [t('form.required')] }; return }
  if (!form.value.start_date) { errors.value = { start_date: [t('form.required')] }; return }
  if (!form.value.end_date) { errors.value = { end_date: [t('form.required')] }; return }
  if (!form.value.as_of_date) { errors.value = { as_of_date: [t('form.required')] }; return }
  saving.value = true
  try {
    if (editing.value) {
      await api.put(`/api/v1/tenant/payroll/periods/${editingId.value}`, {
        start_date: form.value.start_date, end_date: form.value.end_date,
        as_of_date: form.value.as_of_date, status: form.value.status
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.period_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/periods', {
        period_year: form.value.period_year, period_month: form.value.period_month,
        start_date: form.value.start_date, end_date: form.value.end_date,
        as_of_date: form.value.as_of_date, status: form.value.status
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.period_created'), life: 3000 })
    }
    dialogVisible.value = false; await loadData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { saving.value = false }
}
onMounted(loadData)
</script>
