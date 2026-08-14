<template>
  <div class="space-y-4">
    <!-- ── Sub-tabs ── -->
    <div class="flex items-center gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
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

    <!-- ── Settings ── -->
    <div v-if="activeTab === 'settings'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <span v-if="settingTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ settingTotal }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_pph21_setting')" icon="pi pi-plus" size="small" @click="openSettingDialog()" />
      </div>
      <DataTable :value="settingItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="settingLoading">
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-percentage text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.pph21_settings_empty') }}</p></div></template>
        <Column field="setting_code" :header="t('payroll.pph21_setting_code')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.setting_code }}</span></template></Column>
        <Column field="setting_name" :header="t('payroll.pph21_setting_name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.setting_name }}</span></template></Column>
        <Column field="calculation_method" :header="t('payroll.calculation_method')" sortable style="width:200px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.calculation_method }}</span></template></Column>
        <Column field="default_tax_method" :header="t('payroll.default_tax_method')" sortable style="width:120px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.default_tax_method }}</span></template></Column>
        <Column field="effective_start_date" :header="t('payroll.effective_start_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openSettingDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteSetting(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── PTKP Rates ── -->
    <div v-if="activeTab === 'ptkp'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <span v-if="ptkpItems.length" class="text-xs text-gray-400 dark:text-gray-500">{{ ptkpItems.length }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_rate_component')" icon="pi pi-plus" size="small" @click="openPtkpDialog()" />
      </div>
      <DataTable :value="ptkpItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="ptkpLoading">
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-receipt text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.pph21_settings_empty') }}</p></div></template>
        <Column field="ptkp_status" :header="t('payroll.pph21_setting_code')" sortable style="width:120px"><template #body="{data}"><Tag :value="data.ptkp_status" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="description" :header="t('common.description')"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200">{{ data.description || '-' }}</span></template></Column>
        <Column field="annual_amount" :header="t('payroll.gross_salary')" sortable style="width:170px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ formatMoney(data.annual_amount) }}</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      </DataTable>
    </div>

    <!-- ── Tax Brackets ── -->
    <div v-if="activeTab === 'brackets'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <span v-if="bracketItems.length" class="text-xs text-gray-400 dark:text-gray-500">{{ bracketItems.length }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_rate_component')" icon="pi pi-plus" size="small" @click="openBracketDialog()" />
      </div>
      <DataTable :value="bracketItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="bracketLoading">
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-chart-bar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.pph21_settings_empty') }}</p></div></template>
        <Column field="bracket_order" :header="t('payroll.display_order')" sortable style="width:100px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.bracket_order }}</span></template></Column>
        <Column field="lower_bound" :header="t('payroll.min_base_amount')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.lower_bound) }}</span></template></Column>
        <Column field="upper_bound" :header="t('payroll.max_base_amount')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ data.upper_bound ? formatMoney(data.upper_bound) : '∞' }}</span></template></Column>
        <Column field="rate_percent" :header="t('payroll.rate_percent')" sortable style="width:100px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.rate_percent }}%</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      </DataTable>
    </div>

    <!-- ── Setting dialog ── -->
    <Dialog v-model:visible="settingDialogVisible" :header="editingSetting ? t('payroll.pph21') : t('payroll.new_pph21_setting')" modal :style="{ width: '640px' }" @hide="resetSettingForm">
      <div class="space-y-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.pph21_setting_code')" required :errors="errors?.setting_code">
            <TextInput v-model="settingForm.setting_code" maxlength="50" autofocus :placeholder="t('payroll.pph21_setting_code')" :class="{'p-invalid':errors?.setting_code}" :disabled="editingSetting" />
          </FormRow>
          <FormRow :label="t('payroll.pph21_setting_name')" required :errors="errors?.setting_name">
            <TextInput v-model="settingForm.setting_name" maxlength="150" :placeholder="t('payroll.pph21_setting_name')" :class="{'p-invalid':errors?.setting_name}" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.pph21_component')" required :errors="errors?.pph21_component_id">
            <SelectLabel v-model="settingForm.pph21_component_id" :options="componentOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.pph21_component_id}" showClear />
          </FormRow>
          <FormRow :label="t('payroll.default_tax_method')">
            <SelectLabel v-model="settingForm.default_tax_method" :options="taxMethodOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <FormRow :label="t('payroll.occupational_expense_rate')">
            <InputNumber v-model="settingForm.occupational_expense_rate_percent" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.occupational_expense_max_monthly')">
            <InputNumber v-model="settingForm.occupational_expense_max_monthly" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.occupational_expense_max_yearly')">
            <InputNumber v-model="settingForm.occupational_expense_max_yearly" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <FormRow :label="t('payroll.annualization_months')">
            <InputNumber v-model="settingForm.annualization_months" class="!w-full" :min="1" :max="12" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.pkp_rounding_unit')">
            <InputNumber v-model="settingForm.pkp_rounding_unit" class="!w-full" :min="0" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.non_npwp_multiplier_percent')">
            <InputNumber v-model="settingForm.non_npwp_multiplier_percent" class="!w-full" :min="0" size="small" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.deduct_bpjs_health_employee')"><ToggleSwitch v-model="settingForm.deduct_bpjs_health_employee" /></FormRow>
          <FormRow :label="t('payroll.deduct_bpjs_jht_employee')"><ToggleSwitch v-model="settingForm.deduct_bpjs_jht_employee" /></FormRow>
          <FormRow :label="t('payroll.deduct_bpjs_jp_employee')"><ToggleSwitch v-model="settingForm.deduct_bpjs_jp_employee" /></FormRow>
          <FormRow :label="t('payroll.rounding_mode')">
            <SelectLabel v-model="settingForm.rounding_mode" :options="roundingOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="settingForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="settingForm.effective_end_date" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="settingForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="settingDialogVisible=false" /><Button :label="editingSetting ? t('common.update') : t('common.save')" size="small" :loading="settingSaving" :disabled="settingSaving" @click="handleSaveSetting" /></div></template>
    </Dialog>

    <!-- ── PTKP dialog ── -->
    <Dialog v-model:visible="ptkpDialogVisible" :header="t('payroll.new_rate_component')" modal :style="{ width: '480px' }" @hide="resetPtkpForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.pph21_setting_code')" required :errors="errors?.ptkp_status">
          <TextInput v-model="ptkpForm.ptkp_status" maxlength="20" autofocus :placeholder="t('payroll.pph21_setting_code')" :class="{'p-invalid':errors?.ptkp_status}" />
        </FormRow>
        <FormRow :label="t('common.description')">
          <TextInput v-model="ptkpForm.description" maxlength="255" />
        </FormRow>
        <FormRow :label="t('payroll.gross_salary')" required :errors="errors?.annual_amount">
          <InputNumber v-model="ptkpForm.annual_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" :class="{'p-invalid':errors?.annual_amount}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="ptkpForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="ptkpForm.effective_end_date" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="ptkpDialogVisible=false" /><Button :label="t('common.save')" size="small" :loading="ptkpSaving" :disabled="ptkpSaving" @click="handleSavePtkp" /></div></template>
    </Dialog>

    <!-- ── Bracket dialog ── -->
    <Dialog v-model:visible="bracketDialogVisible" :header="t('payroll.new_rate_component')" modal :style="{ width: '480px' }" @hide="resetBracketForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.display_order')" required :errors="errors?.bracket_order">
          <InputNumber v-model="bracketForm.bracket_order" class="!w-full" :min="1" size="small" :class="{'p-invalid':errors?.bracket_order}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.min_base_amount')" required :errors="errors?.lower_bound">
            <InputNumber v-model="bracketForm.lower_bound" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" :class="{'p-invalid':errors?.lower_bound}" />
          </FormRow>
          <FormRow :label="t('payroll.max_base_amount')">
            <InputNumber v-model="bracketForm.upper_bound" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('payroll.rate_percent')" required :errors="errors?.rate_percent">
          <InputNumber v-model="bracketForm.rate_percent" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" :class="{'p-invalid':errors?.rate_percent}" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="bracketForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="bracketForm.effective_end_date" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="bracketDialogVisible=false" /><Button :label="t('common.save')" size="small" :loading="bracketSaving" :disabled="bracketSaving" @click="handleSaveBracket" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.pph21')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
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
const activeTab = ref('settings')
const tabs = [
  { key: 'settings', labelKey: 'payroll.pph21' },
  { key: 'ptkp', labelKey: 'payroll.report_tax' },
  { key: 'brackets', labelKey: 'payroll.rate_components' }
]
const settingItems = ref([]); const settingLoading = ref(false); const settingTotal = ref(0)
const ptkpItems = ref([]); const ptkpLoading = ref(false)
const bracketItems = ref([]); const bracketLoading = ref(false)
const settingDialogVisible = ref(false); const editingSetting = ref(false); const editingSettingId = ref(null); const settingSaving = ref(false)
const ptkpDialogVisible = ref(false); const ptkpSaving = ref(false)
const bracketDialogVisible = ref(false); const bracketSaving = ref(false)
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const components = ref([])
const errors = ref({})
const settingForm = ref({ setting_code: '', setting_name: '', pph21_component_id: null, default_tax_method: 'GROSS', occupational_expense_rate_percent: 5, occupational_expense_max_monthly: 500000, occupational_expense_max_yearly: 6000000, deduct_bpjs_health_employee: false, deduct_bpjs_jht_employee: true, deduct_bpjs_jp_employee: true, annualization_months: 12, pkp_rounding_unit: 1000, non_npwp_multiplier_percent: 100, rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' })
const ptkpForm = ref({ ptkp_status: '', description: '', annual_amount: 0, effective_start_date: '', effective_end_date: '' })
const bracketForm = ref({ bracket_order: 1, lower_bound: 0, upper_bound: null, rate_percent: 0, effective_start_date: '', effective_end_date: '' })

const taxMethodOptions = computed(() => ['GROSS', 'GROSS_UP', 'NETT'].map(v => ({ label: v, value: v })))
const roundingOptions = computed(() => ['NONE', 'ROUND', 'CEIL', 'FLOOR'].map(v => ({ label: t(`payroll.rounding_mode_${v.toLowerCase()}`), value: v })))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: t(`payroll.status_${v.toLowerCase()}`), value: v })))
const componentOptions = computed(() => components.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.pph21_setting_code')}: ${deleteTarget.value.setting_code}` : t('common.no_data'))

function statusLabel(v) { const key = `payroll.status_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function formatMoney(val) { const n = Number(val || 0); return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(n) }

async function loadSettings() {
  settingLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/pph21-settings', { params: { per_page: 200 } })
    settingItems.value = res.data?.data || []
    settingTotal.value = res.data?.total || 0
  } catch(e) { toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) }
  finally { settingLoading.value = false }
}
async function loadPtkp() {
  ptkpLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/pph21-ptkp-rates', { params: { per_page: 200 } })
    ptkpItems.value = res.data?.data || []
  } catch { ptkpItems.value = [] } finally { ptkpLoading.value = false }
}
async function loadBrackets() {
  bracketLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/pph21-tax-brackets', { params: { per_page: 200 } })
    bracketItems.value = res.data?.data || []
  } catch { bracketItems.value = [] } finally { bracketLoading.value = false }
}
async function loadComponents() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { per_page: 500 } })
    components.value = res.data?.data || []
  } catch { components.value = [] }
}

function openSettingDialog(item) {
  editingSetting.value = !!item; editingSettingId.value = item?.id || null; errors.value = {}
  settingForm.value = {
    setting_code: item?.setting_code || '', setting_name: item?.setting_name || '',
    pph21_component_id: item?.pph21_component_id || null,
    default_tax_method: item?.default_tax_method || 'GROSS',
    occupational_expense_rate_percent: item?.occupational_expense_rate_percent ?? 5,
    occupational_expense_max_monthly: item?.occupational_expense_max_monthly ?? 500000,
    occupational_expense_max_yearly: item?.occupational_expense_max_yearly ?? 6000000,
    deduct_bpjs_health_employee: item?.deduct_bpjs_health_employee ?? false,
    deduct_bpjs_jht_employee: item?.deduct_bpjs_jht_employee ?? true,
    deduct_bpjs_jp_employee: item?.deduct_bpjs_jp_employee ?? true,
    annualization_months: item?.annualization_months ?? 12,
    pkp_rounding_unit: item?.pkp_rounding_unit ?? 1000,
    non_npwp_multiplier_percent: item?.non_npwp_multiplier_percent ?? 100,
    rounding_mode: item?.rounding_mode || 'ROUND',
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || '',
    status: item?.status || 'ACTIVE'
  }
  settingDialogVisible.value = true
}
function resetSettingForm() {
  settingForm.value = { setting_code: '', setting_name: '', pph21_component_id: null, default_tax_method: 'GROSS', occupational_expense_rate_percent: 5, occupational_expense_max_monthly: 500000, occupational_expense_max_yearly: 6000000, deduct_bpjs_health_employee: false, deduct_bpjs_jht_employee: true, deduct_bpjs_jp_employee: true, annualization_months: 12, pkp_rounding_unit: 1000, non_npwp_multiplier_percent: 100, rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' }
  errors.value = {}; editingSetting.value = false; editingSettingId.value = null
}
async function handleSaveSetting() {
  errors.value = {}
  if (!settingForm.value.setting_code?.trim()) { errors.value = { setting_code: [t('form.required')] }; return }
  if (!settingForm.value.setting_name?.trim()) { errors.value = { setting_name: [t('form.required')] }; return }
  if (!settingForm.value.pph21_component_id) { errors.value = { pph21_component_id: [t('form.required')] }; return }
  if (!settingForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  settingSaving.value = true
  try {
    const payload = {
      setting_code: settingForm.value.setting_code.trim(),
      setting_name: settingForm.value.setting_name.trim(),
      pph21_component_id: settingForm.value.pph21_component_id,
      default_tax_method: settingForm.value.default_tax_method,
      occupational_expense_rate_percent: settingForm.value.occupational_expense_rate_percent,
      occupational_expense_max_monthly: settingForm.value.occupational_expense_max_monthly,
      occupational_expense_max_yearly: settingForm.value.occupational_expense_max_yearly,
      deduct_bpjs_health_employee: settingForm.value.deduct_bpjs_health_employee,
      deduct_bpjs_jht_employee: settingForm.value.deduct_bpjs_jht_employee,
      deduct_bpjs_jp_employee: settingForm.value.deduct_bpjs_jp_employee,
      annualization_months: settingForm.value.annualization_months,
      pkp_rounding_unit: settingForm.value.pkp_rounding_unit,
      non_npwp_multiplier_percent: settingForm.value.non_npwp_multiplier_percent,
      rounding_mode: settingForm.value.rounding_mode,
      effective_start_date: settingForm.value.effective_start_date,
      effective_end_date: settingForm.value.effective_end_date || null,
      status: settingForm.value.status
    }
    if (editingSetting.value) {
      await api.put(`/api/v1/tenant/payroll/pph21-settings/${editingSettingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/pph21-settings', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_created'), life: 3000 })
    }
    settingDialogVisible.value = false; await loadSettings()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { settingSaving.value = false }
}

function openPtkpDialog() {
  errors.value = {}
  ptkpForm.value = { ptkp_status: '', description: '', annual_amount: 0, effective_start_date: '', effective_end_date: '' }
  ptkpDialogVisible.value = true
}
function resetPtkpForm() {
  ptkpForm.value = { ptkp_status: '', description: '', annual_amount: 0, effective_start_date: '', effective_end_date: '' }
  errors.value = {}
}
async function handleSavePtkp() {
  errors.value = {}
  if (!ptkpForm.value.ptkp_status?.trim()) { errors.value = { ptkp_status: [t('form.required')] }; return }
  if (!ptkpForm.value.annual_amount) { errors.value = { annual_amount: [t('form.required')] }; return }
  if (!ptkpForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  ptkpSaving.value = true
  try {
    await api.post('/api/v1/tenant/payroll/pph21-ptkp-rates', {
      ptkp_status: ptkpForm.value.ptkp_status.trim(),
      description: ptkpForm.value.description || null,
      annual_amount: ptkpForm.value.annual_amount,
      effective_start_date: ptkpForm.value.effective_start_date,
      effective_end_date: ptkpForm.value.effective_end_date || null,
      status: 'ACTIVE'
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_created'), life: 3000 })
    ptkpDialogVisible.value = false; await loadPtkp()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { ptkpSaving.value = false }
}

function openBracketDialog() {
  errors.value = {}
  bracketForm.value = { bracket_order: bracketItems.value.length + 1, lower_bound: 0, upper_bound: null, rate_percent: 0, effective_start_date: '', effective_end_date: '' }
  bracketDialogVisible.value = true
}
function resetBracketForm() {
  bracketForm.value = { bracket_order: 1, lower_bound: 0, upper_bound: null, rate_percent: 0, effective_start_date: '', effective_end_date: '' }
  errors.value = {}
}
async function handleSaveBracket() {
  errors.value = {}
  if (!bracketForm.value.bracket_order) { errors.value = { bracket_order: [t('form.required')] }; return }
  if (bracketForm.value.lower_bound === null || bracketForm.value.lower_bound === undefined) { errors.value = { lower_bound: [t('form.required')] }; return }
  if (!bracketForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  bracketSaving.value = true
  try {
    await api.post('/api/v1/tenant/payroll/pph21-tax-brackets', {
      bracket_order: bracketForm.value.bracket_order,
      lower_bound: bracketForm.value.lower_bound,
      upper_bound: bracketForm.value.upper_bound || null,
      rate_percent: bracketForm.value.rate_percent || 0,
      effective_start_date: bracketForm.value.effective_start_date,
      effective_end_date: bracketForm.value.effective_end_date || null,
      status: 'ACTIVE'
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_created'), life: 3000 })
    bracketDialogVisible.value = false; await loadBrackets()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { bracketSaving.value = false }
}

function confirmDeleteSetting(item) { deleteTarget.value = item; deleteError.value = ''; deleteDialogVisible.value = true }
async function handleDelete() {
  deleting.value = true; deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/payroll/pph21-settings/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadSettings()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadSettings(); loadPtkp(); loadBrackets(); loadComponents() })
</script>
