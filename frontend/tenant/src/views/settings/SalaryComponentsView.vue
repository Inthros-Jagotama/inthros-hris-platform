<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <Button :label="t('payroll.new_component')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      sortField="display_order"
      :sortOrder="1"
    >
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-list text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.components_empty') }}</p></div></template>
      <Column field="code" :header="t('payroll.component_code')" sortable style="width:120px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.code }}</span></template></Column>
      <Column field="name" :header="t('payroll.component_name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
      <Column field="component_type" :header="t('payroll.component_type')" sortable style="width:160px"><template #body="{data}"><Tag :value="typeLabel(data.component_type)" :severity="typeSeverity(data.component_type)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column field="calculation_type" :header="t('payroll.calculation_type')" sortable style="width:130px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ calcLabel(data.calculation_type) }}</span></template></Column>
      <Column field="status" :header="t('common.status')" sortable style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('payroll.salary_components') : t('payroll.new_component')" modal :style="{ width: '640px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.component_code')" required :errors="errors?.code">
            <TextInput v-model="form.code" maxlength="50" autofocus :placeholder="t('payroll.component_code')" :class="{'p-invalid':errors?.code}" :disabled="editing" />
          </FormRow>
          <FormRow :label="t('payroll.component_name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="150" :placeholder="t('payroll.component_name')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.component_type')" required :errors="errors?.component_type">
            <SelectLabel v-model="form.component_type" :options="typeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.component_type}" />
          </FormRow>
          <FormRow :label="t('payroll.calculation_type')" required :errors="errors?.calculation_type">
            <SelectLabel v-model="form.calculation_type" :options="calcOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.calculation_type}" />
          </FormRow>
        </div>
        <FormRow :label="t('common.description')">
          <TextInput v-model="form.description" maxlength="1000" textarea :rows="2" />
        </FormRow>
        <div v-if="form.calculation_type === 'FORMULA' || form.calculation_type === 'PERCENTAGE'">
          <FormRow :label="t('payroll.formula')" :errors="errors?.formula">
            <div class="flex items-start gap-2">
              <TextInput v-model="form.formula" maxlength="500" :placeholder="t('payroll.formula_placeholder')" :class="{'p-invalid':errors?.formula}" />
              <Button :label="t('payroll.validate_formula')" icon="pi pi-check" size="small" severity="secondary" outlined class="!whitespace-nowrap shrink-0" :loading="validatingFormula" @click="validateFormula" />
            </div>
            <small v-if="formulaStatus" class="text-xs mt-1 block" :class="formulaValid ? 'text-emerald-500' : 'text-rose-500'">{{ formulaStatus }}</small>
          </FormRow>
        </div>
        <div v-if="form.calculation_type === 'REFERENCE'">
          <FormRow :label="t('payroll.reference_component')" required :errors="errors?.reference_component_id">
            <SelectLabel v-model="form.reference_component_id" :options="referenceOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.reference_component_id}" showClear />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.is_taxable')"><ToggleSwitch v-model="form.is_taxable" /></FormRow>
          <FormRow :label="t('payroll.is_bpjs_base')"><ToggleSwitch v-model="form.is_bpjs_base" /></FormRow>
          <FormRow :label="t('payroll.is_recurring')"><ToggleSwitch v-model="form.is_recurring" /></FormRow>
          <FormRow :label="t('payroll.is_proratable')"><ToggleSwitch v-model="form.is_proratable" /></FormRow>
          <FormRow :label="t('payroll.print_on_salary_structure')"><ToggleSwitch v-model="form.print_on_salary_structure" /></FormRow>
          <FormRow :label="t('payroll.display_order')">
            <InputNumber v-model="form.display_order" class="!w-full" :min="0" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="form.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.salary_components')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
const { t } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const form = ref({ code: '', name: '', description: '', component_type: 'EARNING', calculation_type: 'FIXED', formula: '', reference_component_id: null, is_taxable: true, is_bpjs_base: false, is_recurring: true, is_proratable: true, print_on_salary_structure: true, display_order: 100, status: 'ACTIVE' })
const allComponents = ref([])
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.component_name')}: ${deleteTarget.value.code} — ${deleteTarget.value.name}` : t('common.no_data'))
const validatingFormula = ref(false); const formulaStatus = ref(''); const formulaValid = ref(false)

const typeOptions = computed(() => ['EARNING', 'DEDUCTION', 'EMPLOYER_CONTRIBUTION', 'INFORMATION'].map(v => ({ label: t(`payroll.component_type_${v.toLowerCase()}`), value: v })))
const calcOptions = computed(() => ['FIXED', 'PERCENTAGE', 'FORMULA', 'REFERENCE', 'MANUAL'].map(v => ({ label: t(`payroll.calculation_type_${v.toLowerCase()}`), value: v })))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: t(`payroll.status_${v.toLowerCase()}`), value: v })))
const referenceOptions = computed(() => allComponents.value.filter(c => c.id !== editingId.value).map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))

const skeletonColumns = [{type:'text',width:'w-20',headerWidth:'w-16'},{type:'text',width:'w-40',headerWidth:'w-16'},{type:'tag',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-20',headerWidth:'w-16'},{type:'tag',width:'w-16',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]

function typeLabel(v) { const key = `payroll.component_type_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function typeSeverity(v) { return { EARNING: 'success', DEDUCTION: 'danger', EMPLOYER_CONTRIBUTION: 'warn', INFORMATION: 'info' }[v] || 'secondary' }
function calcLabel(v) { const key = `payroll.calculation_type_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function statusLabel(v) { const key = `payroll.status_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { page: currentPage.value, per_page: perPage.value } })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally { loading.value = false }
}
async function loadAllComponents() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { per_page: 500 } })
    allComponents.value = res.data?.data || []
  } catch { allComponents.value = [] }
}
function onPage(event) { currentPage.value = event.page + 1; perPage.value = event.rows; loadData() }
function openDialog(item) {
  editing.value = !!item; editingId.value = item?.id || null; errors.value = {}; formulaStatus.value = ''; formulaValid.value = false
  form.value = {
    code: item?.code || '', name: item?.name || '', description: item?.description || '',
    component_type: item?.component_type || 'EARNING', calculation_type: item?.calculation_type || 'FIXED',
    formula: item?.formula || '', reference_component_id: item?.reference_component_id || null,
    is_taxable: item?.is_taxable ?? true, is_bpjs_base: item?.is_bpjs_base ?? false,
    is_recurring: item?.is_recurring ?? true, is_proratable: item?.is_proratable ?? true,
    print_on_salary_structure: item?.print_on_salary_structure ?? true, display_order: item?.display_order ?? 100,
    status: item?.status || 'ACTIVE'
  }
  dialogVisible.value = true
}
function resetForm() {
  form.value = { code: '', name: '', description: '', component_type: 'EARNING', calculation_type: 'FIXED', formula: '', reference_component_id: null, is_taxable: true, is_bpjs_base: false, is_recurring: true, is_proratable: true, print_on_salary_structure: true, display_order: 100, status: 'ACTIVE' }
  errors.value = {}; editing.value = false; editingId.value = null; formulaStatus.value = ''; formulaValid.value = false
}
async function validateFormula() {
  if (!form.value.formula?.trim()) { formulaStatus.value = ''; return }
  validatingFormula.value = true
  try {
    await api.post('/api/v1/tenant/payroll/formula/validate', { formula: form.value.formula })
    formulaValid.value = true; formulaStatus.value = t('payroll.formula_valid')
  } catch (e) {
    formulaValid.value = false; formulaStatus.value = e.response?.data?.error?.message || t('payroll.formula_invalid')
  } finally { validatingFormula.value = false }
}
async function handleSave() {
  errors.value = {}
  if (!form.value.code?.trim()) { errors.value = { code: [t('form.required')] }; return }
  if (!form.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  if (!form.value.component_type) { errors.value = { component_type: [t('form.required')] }; return }
  if (!form.value.calculation_type) { errors.value = { calculation_type: [t('form.required')] }; return }
  saving.value = true
  try {
    const payload = {
      code: form.value.code.trim(), name: form.value.name.trim(),
      description: form.value.description || null,
      component_type: form.value.component_type, calculation_type: form.value.calculation_type,
      formula: (form.value.calculation_type === 'FORMULA' || form.value.calculation_type === 'PERCENTAGE') ? (form.value.formula || null) : null,
      reference_component_id: form.value.calculation_type === 'REFERENCE' ? (form.value.reference_component_id || null) : null,
      is_taxable: form.value.is_taxable, is_bpjs_base: form.value.is_bpjs_base,
      is_recurring: form.value.is_recurring, is_proratable: form.value.is_proratable,
      print_on_salary_structure: form.value.print_on_salary_structure,
      display_order: form.value.display_order, status: form.value.status
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/payroll/salary-components/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/salary-components', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_created'), life: 3000 })
    }
    dialogVisible.value = false; await loadData(); await loadAllComponents()
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
    await api.delete(`/api/v1/tenant/payroll/salary-components/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadData(); await loadAllComponents()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadData(); loadAllComponents() })
</script>
