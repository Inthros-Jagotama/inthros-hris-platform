<template>
  <div class="space-y-4">
    <!-- ── Settings list ── -->
    <div class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalRecords }} {{ t('common.items') }}</span>
        </div>
        <Button :label="t('payroll.new_bpjs_setting')" icon="pi pi-plus" size="small" @click="openSettingDialog()" />
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
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-shield text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.bpjs_settings_empty') }}</p></div></template>
        <Column field="setting_code" :header="t('payroll.bpjs_setting_code')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.setting_code }}</span></template></Column>
        <Column field="setting_name" :header="t('payroll.bpjs_setting_name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.setting_name }}</span></template></Column>
        <Column field="base_source" :header="t('payroll.base_source')" sortable style="width:180px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ baseSourceLabel(data.base_source) }}</span></template></Column>
        <Column field="default_jkk_risk_class" :header="t('payroll.default_jkk_risk_class')" sortable style="width:150px"><template #body="{data}"><Tag :value="riskLabel(data.default_jkk_risk_class)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="effective_start_date" :header="t('payroll.effective_start_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column :header="t('common.actions')" style="width:150px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-list" size="small" text severity="info" v-tooltip.left="t('payroll.rate_components')" @click="openRates(data)" /><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openSettingDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteSetting(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── Rate components (drawer dialog) ── -->
    <Dialog v-model:visible="ratesDialogVisible" :header="`${t('payroll.rate_components')} — ${selectedSetting?.setting_code || ''}`" modal :style="{ width: '900px' }" :closable="true">
      <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
        <span v-if="rateItems.length" class="text-xs text-gray-400">{{ rateItems.length }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_rate_component')" icon="pi pi-plus" size="small" @click="openRateDialog()" />
      </div>
      <DataTable :value="rateItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="ratesLoading">
        <template #empty><div class="text-center py-6 text-sm text-gray-400">{{ t('payroll.bpjs_settings_empty') }}</div></template>
        <Column field="rate_code" :header="t('payroll.rate_code')" style="width:130px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.rate_code }}</span></template></Column>
        <Column field="rate_name" :header="t('payroll.rate_name')"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200">{{ data.rate_name }}</span></template></Column>
        <Column field="bpjs_program" :header="t('payroll.bpjs_program')" style="width:100px"><template #body="{data}"><Tag :value="t('payroll.bpjs_program_' + data.bpjs_program.toLowerCase())" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="paid_by" :header="t('payroll.paid_by')" style="width:100px"><template #body="{data}"><Tag :value="t('payroll.paid_by_' + data.paid_by.toLowerCase())" :severity="data.paid_by === 'EMPLOYEE' ? 'danger' : 'warn'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="rate_percent" :header="t('payroll.rate_percent')" style="width:90px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.rate_percent }}%</span></template></Column>
        <Column field="status" :header="t('common.status')" style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openRateDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteRate(data)" /></div></template></Column>
      </DataTable>
    </Dialog>

    <!-- ── Setting dialog ── -->
    <Dialog v-model:visible="settingDialogVisible" :header="editingSetting ? t('payroll.bpjs') : t('payroll.new_bpjs_setting')" modal :style="{ width: '640px' }" @hide="resetSettingForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.bpjs_setting_name')" required :errors="errors?.setting_name">
          <TextInput v-model="settingForm.setting_name" maxlength="150" autofocus :placeholder="t('payroll.bpjs_setting_name')" :class="{'p-invalid':errors?.setting_name}" />
        </FormRow>
        <FormRow :label="t('payroll.health_max_base_amount')">
          <InputNumber v-model="settingForm.health_max_base_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.pension_max_base_amount')">
          <InputNumber v-model="settingForm.pension_max_base_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.default_jkk_risk_class')">
          <SelectLabel v-model="settingForm.default_jkk_risk_class" :options="riskOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <FormRow :label="t('payroll.rounding_mode')">
          <SelectLabel v-model="settingForm.rounding_mode" :options="roundingOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="settingForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="settingForm.effective_end_date" />
          </FormRow>
        </div>
        <!-- Sumber Dasar — card tersendiri dengan pilihan 1 kolom -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.base_source') }}</span>
          <div class="space-y-2">
            <div v-for="opt in baseSourceOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-3"
                 :class="settingForm.base_source === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="settingForm.base_source = opt.value">
              <RadioButton :modelValue="settingForm.base_source" :inputId="'base-source-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="settingForm.base_source = $event" class="mt-0.5" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'base-source-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
          <div>
            <span class="text-sm font-medium">{{ t('common.status') }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500 block">{{ t('payroll.status_desc') }}</span>
          </div>
          <ToggleSwitch v-model="settingForm.status" :true-value="'ACTIVE'" :false-value="'INACTIVE'" :label="statusLabel(settingForm.status)" class="shrink-0" />
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="settingDialogVisible=false" /><Button :label="editingSetting ? t('common.update') : t('common.save')" size="small" :loading="settingSaving" :disabled="settingSaving" @click="handleSaveSetting" /></div></template>
    </Dialog>

    <!-- ── Rate dialog ── -->
    <Dialog v-model:visible="rateDialogVisible" :header="editingRate ? t('payroll.rate_components') : t('payroll.new_rate_component')" modal :style="{ width: '560px' }" @hide="resetRateForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.rate_name')" required :errors="errors?.rate_name">
          <TextInput v-model="rateForm.rate_name" maxlength="180" autofocus :placeholder="t('payroll.rate_name')" :class="{'p-invalid':errors?.rate_name}" />
        </FormRow>
        <FormRow :label="t('payroll.bpjs_program')" required :errors="errors?.bpjs_program">
          <SelectLabel v-model="rateForm.bpjs_program" :options="programOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.bpjs_program}" />
        </FormRow>
        <FormRow :label="t('payroll.paid_by')" required :errors="errors?.paid_by">
          <SelectLabel v-model="rateForm.paid_by" :options="paidByOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.paid_by}" />
        </FormRow>
        <FormRow :label="t('payroll.rate_percent')" required :errors="errors?.rate_percent">
          <InputNumber v-model="rateForm.rate_percent" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.fixed_amount')">
          <InputNumber v-model="rateForm.fixed_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.min_base_amount')">
          <InputNumber v-model="rateForm.min_base_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.max_base_amount')">
          <InputNumber v-model="rateForm.max_base_amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.jkk_risk_class')">
          <SelectLabel v-model="rateForm.jkk_risk_class" :options="riskOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('payroll.salary_components')">
          <SelectLabel v-model="rateForm.salary_component_id" :options="componentOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
          <DateInput v-model="rateForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
        </FormRow>
        <FormRow :label="t('payroll.effective_end_date')">
          <DateInput v-model="rateForm.effective_end_date" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="rateDialogVisible=false" /><Button :label="editingRate ? t('common.update') : t('common.save')" size="small" :loading="rateSaving" :disabled="rateSaving" @click="handleSaveRate" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="deleteTitle" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import RadioButton from 'primevue/radiobutton'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
const { t, locale } = useI18n(); const toast = useToast()
const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const settingDialogVisible = ref(false); const editingSetting = ref(false); const editingSettingId = ref(null); const settingSaving = ref(false)
const ratesDialogVisible = ref(false); const selectedSetting = ref(null); const rateItems = ref([]); const ratesLoading = ref(false)
const rateDialogVisible = ref(false); const editingRate = ref(false); const editingRateId = ref(null); const rateSaving = ref(false)
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null); const deleteIsSetting = ref(true)
const components = ref([])
const errors = ref({})
const settingForm = ref({ setting_code: '', setting_name: '', base_source: 'BPJS_BASE_COMPONENTS', health_max_base_amount: null, pension_max_base_amount: null, default_jkk_risk_class: 'LOW', rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' })
const rateForm = ref({ rate_code: '', rate_name: '', bpjs_program: 'HEALTH', paid_by: 'EMPLOYEE', rate_percent: 0, fixed_amount: null, min_base_amount: null, max_base_amount: null, jkk_risk_class: null, salary_component_id: null, effective_start_date: '', effective_end_date: '' })

const baseSourceOptions = computed(() => ['BPJS_BASE_COMPONENTS', 'BASIC_SALARY', 'GROSS_EARNING'].map(v => ({ label: t(`payroll.base_source_${v.toLowerCase()}`), desc: t(`payroll.base_source_${v.toLowerCase()}_desc`), value: v })))
const riskOptions = computed(() => ['VERY_LOW', 'LOW', 'MEDIUM', 'HIGH', 'VERY_HIGH'].map(v => ({ label: t(`payroll.risk_${v.toLowerCase()}`), value: v })))
const roundingOptions = computed(() => ['NONE', 'ROUND', 'CEIL', 'FLOOR'].map(v => ({ label: t(`payroll.rounding_mode_${v.toLowerCase()}`), value: v })))
const programOptions = computed(() => ['HEALTH', 'JHT', 'JP', 'JKK', 'JKM', 'JKP'].map(v => ({ label: t(`payroll.bpjs_program_${v.toLowerCase()}`), value: v })))
const paidByOptions = computed(() => ['EMPLOYEE', 'EMPLOYER'].map(v => ({ label: t(`payroll.paid_by_${v.toLowerCase()}`), value: v })))
const componentOptions = computed(() => components.value.map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))
const skeletonColumns = [{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-40',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-20'},{type:'tag',width:'w-20',headerWidth:'w-16'},{type:'text',width:'w-20',headerWidth:'w-16'},{type:'tag',width:'w-16',headerWidth:'w-16'},{type:'icons',count:3,headerWidth:'w-24'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// Kode digenerate otomatis dari nama (inisial tiap kata + timestamp + random),
// memakai _ sebagai pemisah agar tidak konflik dengan operator minus di formula.
function nameInitials(name) {
  const words = (name || '').trim().split(/\s+/).filter(Boolean)
  return words.map(w => w[0]).join('').toUpperCase().replace(/[^A-Z0-9]/g, '')
}
function generateCode(name) {
  const prefix = nameInitials(name)
  const d = new Date()
  const pad = n => String(n).padStart(2, '0')
  const stamp = `${d.getFullYear()}${pad(d.getMonth() + 1)}${pad(d.getDate())}${pad(d.getHours())}${pad(d.getMinutes())}${pad(d.getSeconds())}`
  const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  const rand = Array.from({ length: 4 }, () => alphabet[Math.floor(Math.random() * alphabet.length)]).join('')
  return prefix ? `${prefix}_${stamp}_${rand}` : ''
}
watch(() => settingForm.value.setting_name, (nv) => { if (!editingSetting.value) settingForm.value.setting_code = generateCode(nv) })
watch(() => rateForm.value.rate_name, (nv) => { if (!editingRate.value) rateForm.value.rate_code = generateCode(nv) })
const deleteTitle = computed(() => deleteIsSetting.value ? t('payroll.bpjs') : t('payroll.rate_components'))
const deleteMessage = computed(() => {
  if (!deleteTarget.value) return t('common.no_data')
  return deleteIsSetting.value ? `${t('payroll.bpjs_setting_code')}: ${deleteTarget.value.setting_code}` : `${t('payroll.rate_code')}: ${deleteTarget.value.rate_code}`
})

function baseSourceLabel(v) { const key = `payroll.base_source_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function riskLabel(v) { const key = `payroll.risk_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function statusLabel(v) { const key = `payroll.status_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/bpjs-settings', { params: { page: currentPage.value, per_page: perPage.value } })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally { loading.value = false }
}
async function loadComponents() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { per_page: 500 } })
    components.value = res.data?.data || []
  } catch { components.value = [] }
}
function onPage(event) { currentPage.value = event.page + 1; perPage.value = event.rows; loadData() }

async function openRates(setting) {
  selectedSetting.value = setting
  ratesDialogVisible.value = true
  ratesLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/bpjs-rate-components', { params: { bpjs_setting_id: setting.id, per_page: 500 } })
    rateItems.value = res.data?.data || []
  } catch { rateItems.value = [] } finally { ratesLoading.value = false }
}

function openSettingDialog(item) {
  editingSetting.value = !!item; editingSettingId.value = item?.id || null; errors.value = {}
  settingForm.value = {
    setting_code: item?.setting_code || '', setting_name: item?.setting_name || '',
    base_source: item?.base_source || 'BPJS_BASE_COMPONENTS',
    health_max_base_amount: item?.health_max_base_amount ?? null,
    pension_max_base_amount: item?.pension_max_base_amount ?? null,
    default_jkk_risk_class: item?.default_jkk_risk_class || 'LOW',
    rounding_mode: item?.rounding_mode || 'ROUND',
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || '',
    status: item?.status || 'ACTIVE'
  }
  settingDialogVisible.value = true
}
function resetSettingForm() {
  settingForm.value = { setting_code: '', setting_name: '', base_source: 'BPJS_BASE_COMPONENTS', health_max_base_amount: null, pension_max_base_amount: null, default_jkk_risk_class: 'LOW', rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' }
  errors.value = {}; editingSetting.value = false; editingSettingId.value = null
}
async function handleSaveSetting() {
  errors.value = {}
  if (!settingForm.value.setting_name?.trim()) { errors.value = { setting_name: [t('form.required')] }; return }
  if (!settingForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  settingSaving.value = true
  try {
    const payload = {
      setting_code: settingForm.value.setting_code.trim(),
      setting_name: settingForm.value.setting_name.trim(),
      base_source: settingForm.value.base_source,
      health_max_base_amount: settingForm.value.health_max_base_amount ?? null,
      pension_max_base_amount: settingForm.value.pension_max_base_amount ?? null,
      default_jkk_risk_class: settingForm.value.default_jkk_risk_class,
      rounding_mode: settingForm.value.rounding_mode,
      effective_start_date: settingForm.value.effective_start_date,
      effective_end_date: settingForm.value.effective_end_date || null,
      status: settingForm.value.status
    }
    if (editingSetting.value) {
      await api.put(`/api/v1/tenant/payroll/bpjs-settings/${editingSettingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bpjs_setting_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/bpjs-settings', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bpjs_setting_created'), life: 3000 })
    }
    settingDialogVisible.value = false; await loadData()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { settingSaving.value = false }
}

function openRateDialog(item) {
  editingRate.value = !!item; editingRateId.value = item?.id || null; errors.value = {}
  rateForm.value = {
    rate_code: item?.rate_code || '', rate_name: item?.rate_name || '',
    bpjs_program: item?.bpjs_program || 'HEALTH', paid_by: item?.paid_by || 'EMPLOYEE',
    rate_percent: item?.rate_percent ?? 0, fixed_amount: item?.fixed_amount ?? null,
    min_base_amount: item?.min_base_amount ?? null, max_base_amount: item?.max_base_amount ?? null,
    jkk_risk_class: item?.jkk_risk_class || null, salary_component_id: item?.salary_component_id || null,
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || ''
  }
  rateDialogVisible.value = true
}
function resetRateForm() {
  rateForm.value = { rate_code: '', rate_name: '', bpjs_program: 'HEALTH', paid_by: 'EMPLOYEE', rate_percent: 0, fixed_amount: null, min_base_amount: null, max_base_amount: null, jkk_risk_class: null, salary_component_id: null, effective_start_date: '', effective_end_date: '' }
  errors.value = {}; editingRate.value = false; editingRateId.value = null
}
async function handleSaveRate() {
  errors.value = {}
  if (!rateForm.value.rate_name?.trim()) { errors.value = { rate_name: [t('form.required')] }; return }
  if (!rateForm.value.bpjs_program) { errors.value = { bpjs_program: [t('form.required')] }; return }
  if (!rateForm.value.paid_by) { errors.value = { paid_by: [t('form.required')] }; return }
  if (!rateForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  rateSaving.value = true
  try {
    const payload = {
      bpjs_setting_id: selectedSetting.value.id,
      rate_code: rateForm.value.rate_code.trim(),
      rate_name: rateForm.value.rate_name.trim(),
      bpjs_program: rateForm.value.bpjs_program,
      paid_by: rateForm.value.paid_by,
      rate_percent: rateForm.value.rate_percent || 0,
      fixed_amount: rateForm.value.fixed_amount ?? null,
      min_base_amount: rateForm.value.min_base_amount ?? null,
      max_base_amount: rateForm.value.max_base_amount ?? null,
      jkk_risk_class: rateForm.value.jkk_risk_class || null,
      salary_component_id: rateForm.value.salary_component_id || null,
      effective_start_date: rateForm.value.effective_start_date,
      effective_end_date: rateForm.value.effective_end_date || null
    }
    if (editingRate.value) {
      await api.put(`/api/v1/tenant/payroll/bpjs-rate-components/${editingRateId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/bpjs-rate-components', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_created'), life: 3000 })
    }
    rateDialogVisible.value = false
    await openRates(selectedSetting.value)
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { rateSaving.value = false }
}

function confirmDeleteSetting(item) { deleteIsSetting.value = true; deleteTarget.value = item; deleteError.value = ''; deleteDialogVisible.value = true }
function confirmDeleteRate(item) { deleteIsSetting.value = false; deleteTarget.value = item; deleteError.value = ''; deleteDialogVisible.value = true }
async function handleDelete() {
  deleting.value = true; deleteError.value = ''
  try {
    if (deleteIsSetting.value) {
      await api.delete(`/api/v1/tenant/payroll/bpjs-settings/${deleteTarget.value.id}`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bpjs_setting_deleted'), life: 3000 })
      deleteDialogVisible.value = false; await loadData()
    } else {
      await api.delete(`/api/v1/tenant/payroll/bpjs-rate-components/${deleteTarget.value.id}`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_deleted'), life: 3000 })
      deleteDialogVisible.value = false
      if (selectedSetting.value) await openRates(selectedSetting.value)
    }
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadData(); loadComponents() })
</script>
