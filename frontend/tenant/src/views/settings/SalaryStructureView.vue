<template>
  <div class="space-y-4">
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div class="flex items-center gap-1 px-3 pt-2 border-b border-gray-200 dark:border-gray-700">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          type="button"
          class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
          :class="activeTab === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
          @click="activeTab = tab.key"
        >
          {{ t(tab.labelKey) }}
        </button>
      </div>

      <!-- ── Grade Components ── -->
      <div v-if="activeTab === 'grade'" class="p-4">
        <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
          <span v-if="gradeTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ gradeTotal }} {{ t('common.items') }}</span>
          <Button :label="t('payroll.new_grade_component')" icon="pi pi-plus" size="small" @click="openGradeDialog()" />
        </div>
        <DataTable
          :value="gradeItems"
          lazy
          :totalRecords="gradeTotal"
          :first="gradeFirst"
          :rows="perPage"
          @page="onGradePage($event)"
          paginator
          paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
          :rowsPerPageOptions="[10, 15, 25, 50]"
          size="small"
          class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        >
          <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-chart-bar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.grade_components_empty') }}</p></div></template>
          <Column :header="t('payroll.grading')" style="width:140px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ gradingLabel(data.grading_id) }}</span></template></Column>
          <Column :header="t('payroll.component')" style="width:200px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100">{{ componentLabel(data.salary_component_id) }}</span></template></Column>
          <Column :header="t('payroll.amount')" style="width:130px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ formatMoney(data.amount) }}</span></template></Column>
          <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
          <Column :header="t('common.status')" style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
          <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openGradeDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data, 'grade')" /></div></template></Column>
        </DataTable>
      </div>

      <!-- ── Employee Components ── -->
      <div v-if="activeTab === 'employee'" class="p-4">
        <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
          <span v-if="empTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ empTotal }} {{ t('common.items') }}</span>
          <Button :label="t('payroll.new_employee_component')" icon="pi pi-plus" size="small" @click="openEmpDialog()" />
        </div>
        <DataTable
          :value="empItems"
          lazy
          :totalRecords="empTotal"
          :first="empFirst"
          :rows="perPage"
          @page="onEmpPage($event)"
          paginator
          paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
          :rowsPerPageOptions="[10, 15, 25, 50]"
          size="small"
          class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        >
          <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-users text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.employee_components_empty') }}</p></div></template>
          <Column :header="t('payroll.employee')" style="width:200px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeLabel(data.employee_id) }}</span></template></Column>
          <Column :header="t('payroll.component')" style="width:200px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100">{{ componentLabel(data.salary_component_id) }}</span></template></Column>
          <Column :header="t('payroll.amount')" style="width:130px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ formatMoney(data.amount) }}</span></template></Column>
          <Column :header="t('payroll.source_type')" style="width:120px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ sourceLabel(data.source_type) }}</span></template></Column>
          <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
          <Column :header="t('common.status')" style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
          <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEmpDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data, 'employee')" /></div></template></Column>
        </DataTable>
      </div>
    </div>

    <!-- ── Grade dialog ── -->
    <Dialog v-model:visible="gradeDialogVisible" :header="gradeEditing ? t('payroll.edit_grade_component') : t('payroll.new_grade_component')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetGradeForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.grading')">
          <SelectLabel v-model="gradeForm.grading_id" :options="gradingOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.grading_id}" showClear />
        </FormRow>
        <FormRow :label="t('payroll.component')" required :errors="errors?.salary_component_id">
          <SelectLabel v-model="gradeForm.salary_component_id" :options="componentOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{'p-invalid':errors?.salary_component_id}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.amount')" required :errors="errors?.amount">
            <InputNumber v-model="gradeForm.amount" class="!w-full" :min="0" size="small" mode="currency" currency="IDR" locale="id-ID" :class="{'p-invalid':errors?.amount}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="gradeForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="gradeForm.effective_end_date" />
          </FormRow>
          <FormRow :label="t('common.status')">
            <SelectLabel v-model="gradeForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.is_mandatory')"><ToggleSwitch v-model="gradeForm.is_mandatory" /></FormRow>
          <FormRow :label="t('payroll.is_default')"><ToggleSwitch v-model="gradeForm.is_default" /></FormRow>
        </div>
        <FormRow :label="t('common.description')">
          <TextInput v-model="gradeForm.notes" maxlength="1000" textarea :rows="2" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="gradeDialogVisible=false" /><Button :label="gradeEditing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleGradeSave" /></div></template>
    </Dialog>

    <!-- ── Employee dialog ── -->
    <Dialog v-model:visible="empDialogVisible" :header="empEditing ? t('payroll.edit_employee_component') : t('payroll.new_employee_component')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetEmpForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.employee')" required :errors="errors?.employee_id">
          <SelectLabel v-model="empForm.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{'p-invalid':errors?.employee_id}" :disabled="empEditing" />
        </FormRow>
        <FormRow :label="t('payroll.component')" required :errors="errors?.salary_component_id">
          <SelectLabel v-model="empForm.salary_component_id" :options="componentOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{'p-invalid':errors?.salary_component_id}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.amount')" required :errors="errors?.amount">
            <InputNumber v-model="empForm.amount" class="!w-full" :min="0" size="small" mode="currency" currency="IDR" locale="id-ID" :class="{'p-invalid':errors?.amount}" />
          </FormRow>
          <FormRow :label="t('payroll.source_type')">
            <SelectLabel v-model="empForm.source_type" :options="sourceOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="empForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="empForm.effective_end_date" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('common.status')">
            <SelectLabel v-model="empForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
          <FormRow :label="t('common.description')">
            <TextInput v-model="empForm.notes" maxlength="1000" textarea :rows="2" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="empDialogVisible=false" /><Button :label="empEditing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleEmpSave" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.salary_structure')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import DateInput from '@/components/DateInput.vue'
const { t, locale } = useI18n(); const toast = useToast()
const activeTab = ref('grade')
const tabs = [{ key: 'grade', labelKey: 'payroll.grade_components' }, { key: 'employee', labelKey: 'payroll.employee_components' }]
const perPage = ref(15)

const gradeItems = ref([]); const gradeTotal = ref(0); const gradePage = ref(1)
const empItems = ref([]); const empTotal = ref(0); const empPage = ref(1)
const gradings = ref([]); const components = ref([]); const employees = ref([])
const gradeFirst = computed(() => (gradePage.value - 1) * perPage.value)
const empFirst = computed(() => (empPage.value - 1) * perPage.value)

const gradeDialogVisible = ref(false); const gradeEditing = ref(false); const gradeEditingId = ref(null)
const empDialogVisible = ref(false); const empEditing = ref(false); const empEditingId = ref(null)
const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null); const deleteKind = ref('grade')

const gradeForm = ref({ grading_id: null, salary_component_id: null, amount: 0, effective_start_date: '', effective_end_date: null, is_mandatory: true, is_default: true, status: 'ACTIVE', notes: '' })
const empForm = ref({ employee_id: null, salary_component_id: null, amount: 0, source_type: 'MANUAL', effective_start_date: '', effective_end_date: null, status: 'ACTIVE', notes: '' })

const gradingOptions = computed(() => gradings.value.map(g => ({ label: g.name, value: g.id })))
const componentOptions = computed(() => components.value.map(c => ({ label: `${c.name} (${c.code})`, value: c.id })))
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: statusLabel(v), value: v })))
const sourceOptions = computed(() => ['MANUAL', 'GRADE_INHERIT', 'FORMULA', 'ADJUSTMENT'].map(v => ({ label: sourceLabel(v), value: v })))
const deleteMessage = computed(() => {
  if (!deleteTarget.value) return t('common.no_data')
  return deleteKind.value === 'grade'
    ? `${componentLabel(deleteTarget.value.salary_component_id)} — ${t('payroll.grading')}: ${gradingLabel(deleteTarget.value.grading_id)}`
    : `${componentLabel(deleteTarget.value.salary_component_id)} — ${employeeLabel(deleteTarget.value.employee_id)}`
})

function formatMoney(val) { const n = Number(val || 0); return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(n) }
function statusLabel(v) { const key = `payroll.status_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function sourceLabel(v) { const key = `payroll.source_type_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function gradingLabel(id) { const g = gradings.value.find(x => x.id === id); return g ? g.name : (id || '-') }
function componentLabel(id) { const c = components.value.find(x => x.id === id); return c ? `${c.name} (${c.code})` : (id || '-') }
function employeeLabel(id) { const e = employees.value.find(x => x.id === id); return e ? `${e.name} (${e.employee_id})` : (id || '-') }

async function loadGradeData() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-grade-components', { params: { page: gradePage.value, per_page: perPage.value } })
    const body = res.data
    gradeItems.value = body?.data || []
    gradeTotal.value = body?.total || 0
  } catch(e) { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 }) }
}
async function loadEmpData() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-employee-components', { params: { page: empPage.value, per_page: perPage.value } })
    const body = res.data
    empItems.value = body?.data || []
    empTotal.value = body?.total || 0
  } catch(e) { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 }) }
}
async function loadLookups() {
  try { const g = await api.get('/api/v1/tenant/settings/gradings', { params: { per_page: 500 } }); gradings.value = g.data?.data || [] } catch { gradings.value = [] }
  try {
    const all = []
    let page = 1
    while (true) {
      const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { page, per_page: 100 } })
      const rows = res.data?.data || []
      all.push(...rows)
      if (!rows.length || all.length >= (res.data?.total || 0)) break
      page++
    }
    components.value = all
  } catch { components.value = [] }
  try { const e = await api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }); employees.value = e.data?.data || [] } catch { employees.value = [] }
}
function onGradePage(e) { gradePage.value = e.page + 1; perPage.value = e.rows; loadGradeData() }
function onEmpPage(e) { empPage.value = e.page + 1; perPage.value = e.rows; loadEmpData() }

function openGradeDialog(item) {
  gradeEditing.value = !!item; gradeEditingId.value = item?.id || null; errors.value = {}
  gradeForm.value = {
    grading_id: item?.grading_id || null, salary_component_id: item?.salary_component_id || null,
    amount: item?.amount ?? 0, effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || null,
    is_mandatory: item?.is_mandatory ?? true, is_default: item?.is_default ?? true, status: item?.status || 'ACTIVE', notes: item?.notes || ''
  }
  gradeDialogVisible.value = true
}
function resetGradeForm() {
  gradeForm.value = { grading_id: null, salary_component_id: null, amount: 0, effective_start_date: '', effective_end_date: null, is_mandatory: true, is_default: true, status: 'ACTIVE', notes: '' }
  errors.value = {}; gradeEditing.value = false; gradeEditingId.value = null
}
async function handleGradeSave() {
  errors.value = {}
  if (!gradeForm.value.salary_component_id) { errors.value = { salary_component_id: [t('form.required')] }; return }
  if (!gradeForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      grading_id: gradeForm.value.grading_id || null, salary_component_id: gradeForm.value.salary_component_id,
      amount: gradeForm.value.amount || 0, currency_code: 'IDR',
      effective_start_date: gradeForm.value.effective_start_date, effective_end_date: gradeForm.value.effective_end_date || null,
      is_mandatory: gradeForm.value.is_mandatory, is_default: gradeForm.value.is_default,
      status: gradeForm.value.status, notes: gradeForm.value.notes || null
    }
    if (gradeEditing.value) {
      await api.put(`/api/v1/tenant/payroll/salary-grade-components/${gradeEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.grade_component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/salary-grade-components', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.grade_component_created'), life: 3000 })
    }
    gradeDialogVisible.value = false; await loadGradeData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { saving.value = false }
}

function openEmpDialog(item) {
  empEditing.value = !!item; empEditingId.value = item?.id || null; errors.value = {}
  empForm.value = {
    employee_id: item?.employee_id || null, salary_component_id: item?.salary_component_id || null,
    amount: item?.amount ?? 0, source_type: item?.source_type || 'MANUAL',
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || null,
    status: item?.status || 'ACTIVE', notes: item?.notes || ''
  }
  empDialogVisible.value = true
}
function resetEmpForm() {
  empForm.value = { employee_id: null, salary_component_id: null, amount: 0, source_type: 'MANUAL', effective_start_date: '', effective_end_date: null, status: 'ACTIVE', notes: '' }
  errors.value = {}; empEditing.value = false; empEditingId.value = null
}
async function handleEmpSave() {
  errors.value = {}
  if (!empForm.value.employee_id) { errors.value = { employee_id: [t('form.required')] }; return }
  if (!empForm.value.salary_component_id) { errors.value = { salary_component_id: [t('form.required')] }; return }
  if (!empForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      employee_id: empForm.value.employee_id, salary_component_id: empForm.value.salary_component_id,
      amount: empForm.value.amount || 0, currency_code: 'IDR', source_type: empForm.value.source_type,
      effective_start_date: empForm.value.effective_start_date, effective_end_date: empForm.value.effective_end_date || null,
      status: empForm.value.status, notes: empForm.value.notes || null
    }
    if (empEditing.value) {
      await api.put(`/api/v1/tenant/payroll/salary-employee-components/${empEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.employee_component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/salary-employee-components', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.employee_component_created'), life: 3000 })
    }
    empDialogVisible.value = false; await loadEmpData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { saving.value = false }
}

function confirmDelete(item, kind) { deleteTarget.value = item; deleteKind.value = kind; deleteError.value = ''; deleteDialogVisible.value = true }
async function handleDelete() {
  deleting.value = true; deleteError.value = ''
  try {
    const base = deleteKind.value === 'grade' ? 'salary-grade-components' : 'salary-employee-components'
    await api.delete(`/api/v1/tenant/payroll/${base}/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_deleted'), life: 3000 })
    deleteDialogVisible.value = false
    if (deleteKind.value === 'grade') { await loadGradeData() } else { await loadEmpData() }
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadGradeData(); loadEmpData(); loadLookups() })
</script>
