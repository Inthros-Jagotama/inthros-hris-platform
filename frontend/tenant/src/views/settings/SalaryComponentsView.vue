<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <span v-if="totalCount > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ totalCount }} {{ t('common.items') }}</span>
    </div>

    <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <SkeletonTable v-for="i in 2" :key="i" :columns="skeletonColumns" :rows="5" />
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div
        v-for="group in groups"
        :key="group.type"
        class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
      >
        <div class="flex items-center justify-between gap-2 px-4 py-3 border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-2 min-w-0">
            <i :class="group.icon" class="text-sm text-gray-400 dark:text-gray-500"></i>
            <span class="text-sm font-semibold text-surface-800 dark:text-surface-0 truncate">{{ group.label }}</span>
            <Tag :value="String(group.items.length)" :severity="group.severity" class="!text-xs !px-1.5 !py-0.5" />
          </div>
          <Button :label="t('payroll.new_component')" icon="pi pi-plus" size="small" severity="secondary" outlined class="!whitespace-nowrap shrink-0" @click="openDialog(null, group.type)" />
        </div>
        <div class="divide-y divide-gray-100 dark:divide-gray-700/60">
          <div v-if="group.items.length === 0" class="px-4 py-6 text-center text-xs text-gray-400 dark:text-gray-500">
            {{ t('payroll.components_empty_group', { type: group.label }) }}
          </div>
          <div v-for="item in group.items" :key="item.id" class="flex items-center justify-between gap-2 px-4 py-2.5 hover:bg-gray-50 dark:hover:bg-gray-700/40 transition-colors">
            <div class="flex items-center gap-2 min-w-0">
              <span class="font-mono text-xs text-gray-400 dark:text-gray-500 shrink-0">{{ item.code }}</span>
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0 truncate">{{ item.name }}</span>
              <Tag :value="calcLabel(item.calculation_type)" severity="secondary" class="!text-xs !px-1.5 !py-0.5 hidden sm:inline-flex shrink-0" />
              <Tag :value="statusLabel(item.status)" :severity="item.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5 shrink-0" />
            </div>
            <div class="flex items-center gap-1 shrink-0">
              <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(item)" />
              <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(item)" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <Dialog v-model:visible="dialogVisible" :header="editing ? t('payroll.edit_component_of_type', { type: typeLabel(form.component_type) }) : t('payroll.new_component_of_type', { type: typeLabel(form.component_type) })" modal :style="{ width: 'min(900px, 95vw)' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.component_name')" required :errors="errors?.name">
          <TextInput v-model="form.name" maxlength="150" autofocus :placeholder="t('payroll.component_name')" :class="{'p-invalid':errors?.name}" />
        </FormRow>
        <FormRow :label="t('common.description')">
          <TextInput v-model="form.description" maxlength="1000" textarea :rows="2" />
        </FormRow>
        <FormRow :label="t('payroll.calculation_type')" required :errors="errors?.calculation_type">
          <div class="flex flex-wrap gap-2">
            <RadioLabel
              v-for="opt in calcOptions"
              :key="opt.value"
              v-model="form.calculation_type"
              :value="opt.value"
              :id="'calc-type-' + opt.value.toLowerCase()"
              :label="opt.label"
            />
          </div>
        </FormRow>
        <div v-if="form.calculation_type === 'FORMULA' || form.calculation_type === 'PERCENTAGE'">
          <FormRow :label="t('payroll.formula')" :errors="errors?.formula">
            <div class="flex items-start gap-2">
              <TextInput v-model="form.formula" textarea :rows="3" maxlength="500" :placeholder="t('payroll.formula_placeholder')" :class="{'p-invalid':errors?.formula}" />
              <Button :label="t('payroll.validate_formula')" icon="pi pi-check" size="small" severity="secondary" outlined class="!whitespace-nowrap shrink-0" :loading="validatingFormula" @click="validateFormula" />
            </div>
            <div v-if="mentionOptions.length" class="mt-2 border border-gray-200 dark:border-gray-700 rounded-lg p-2 flex flex-wrap gap-1.5 max-h-28 overflow-y-auto">
              <Button v-for="opt in mentionOptions" :key="opt.code" :label="opt.label" size="small" text severity="secondary" class="!px-2 !py-1 !text-xs" @click="insertMention(opt)" />
            </div>
            <span class="text-xs text-gray-400 dark:text-gray-500 mt-1 block">{{ t('payroll.formula_at_hint') }}</span>
            <small v-if="formulaStatus" class="text-xs mt-1 block" :class="formulaValid ? 'text-emerald-500' : 'text-rose-500'">{{ formulaStatus }}</small>
          </FormRow>
        </div>
        <div v-if="form.calculation_type === 'REFERENCE'">
          <FormRow :label="t('payroll.reference_component')" required :errors="errors?.reference_component_id">
            <SelectLabel v-model="form.reference_component_id" :options="referenceOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{'p-invalid':errors?.reference_component_id}" showClear />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_taxable') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_taxable_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_taxable" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_bpjs_base') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_bpjs_base_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_bpjs_base" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_recurring') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_recurring_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_recurring" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_proratable') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_proratable_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_proratable" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.print_on_salary_structure') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.print_on_salary_structure_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.print_on_salary_structure" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_pph21_component') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_pph21_component_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_pph21_component" class="shrink-0" />
          </div>
          <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.is_pph21_deductible') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.is_pph21_deductible_desc') }}</span>
            </div>
            <ToggleSwitch v-model="form.is_pph21_deductible" class="shrink-0" />
          </div>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
          <div class="flex flex-col gap-0.5 min-w-0">
            <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('common.status') }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.status_desc') }}</span>
          </div>
          <ToggleSwitch v-model="form.status" :true-value="'ACTIVE'" :false-value="'INACTIVE'" :label="statusLabel(form.status)" class="shrink-0" />
        </div>
        <FormRow :label="t('payroll.display_order')">
          <InputNumber v-model="form.display_order" class="!w-full" :min="0" size="small" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.salary_components')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import api from '@/services/api'
import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import RadioLabel from '@/components/RadioLabel.vue'
const { t } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const presetType = ref('')
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const form = ref({ code: '', name: '', description: '', component_type: 'EARNING', calculation_type: 'FIXED', formula: '', reference_component_id: null, is_taxable: true, is_bpjs_base: false, is_recurring: true, is_proratable: true, print_on_salary_structure: true, is_pph21_component: false, is_pph21_deductible: false, display_order: 100, status: 'ACTIVE' })

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
watch(() => form.value.name, (nv) => { if (!editing.value) form.value.code = generateCode(nv) })
const allComponents = ref([])
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.component_name')}: ${deleteTarget.value.code} — ${deleteTarget.value.name}` : t('common.no_data'))
const validatingFormula = ref(false); const formulaStatus = ref(''); const formulaValid = ref(false)
const mentionQuery = computed(() => {
  const text = form.value.formula || ''
  const idx = text.lastIndexOf('@')
  if (idx < 0) return null
  return { idx, q: text.slice(idx + 1) }
})
const mentionOptions = computed(() => {
  const m = mentionQuery.value
  if (!m) return []
  const q = m.q.toLowerCase()
  return allComponents.value
    .filter(c => c.id !== editingId.value && (c.name.toLowerCase().includes(q) || c.code.toLowerCase().includes(q)))
    .slice(0, 10)
    .map(c => ({ label: `${c.name} (${c.code})`, code: c.code }))
})

const typeOrder = ['EARNING', 'DEDUCTION', 'EMPLOYER_CONTRIBUTION', 'INFORMATION']
const groupMeta = {
  EARNING: { icon: 'pi pi-wallet' },
  DEDUCTION: { icon: 'pi pi-arrow-down-left' },
  EMPLOYER_CONTRIBUTION: { icon: 'pi pi-building' },
  INFORMATION: { icon: 'pi pi-info-circle' }
}
const groups = computed(() => typeOrder.map(type => ({
  type,
  label: typeLabel(type),
  severity: typeSeverity(type),
  icon: groupMeta[type].icon,
  items: items.value
    .filter(c => c.component_type === type)
    .sort((a, b) => (a.display_order ?? 0) - (b.display_order ?? 0))
})))
const totalCount = computed(() => items.value.length)

const calcOptions = computed(() => ['FIXED', 'PERCENTAGE', 'FORMULA', 'REFERENCE', 'MANUAL'].map(v => ({ label: t(`payroll.calculation_type_${v.toLowerCase()}`), value: v })))
const referenceOptions = computed(() => allComponents.value.filter(c => c.id !== editingId.value).map(c => ({ label: `${c.code} — ${c.name}`, value: c.id })))

const skeletonColumns = [{type:'text',width:'w-16',headerWidth:'w-16'},{type:'text',width:'w-40',headerWidth:'w-16'},{type:'tag',width:'w-20',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]

function typeLabel(v) { const key = `payroll.component_type_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function typeSeverity(v) { return { EARNING: 'success', DEDUCTION: 'danger', EMPLOYER_CONTRIBUTION: 'warn', INFORMATION: 'info' }[v] || 'secondary' }
function calcLabel(v) { const key = `payroll.calculation_type_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function statusLabel(v) { const key = `payroll.status_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }

async function loadData() {
  loading.value = true
  try {
    const all = []
    let page = 1
    const perPage = 100
    while (true) {
      const res = await api.get('/api/v1/tenant/payroll/salary-components', { params: { page, per_page: perPage } })
      const body = res.data
      const rows = body?.data || []
      all.push(...rows)
      if (!rows.length || all.length >= (body?.total || 0)) break
      page++
    }
    items.value = all
    allComponents.value = all
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally { loading.value = false }
}
function openDialog(item, type) {
  editing.value = !!item; editingId.value = item?.id || null; errors.value = {}; formulaStatus.value = ''; formulaValid.value = false
  presetType.value = item ? '' : (type || '')
  form.value = {
    code: item?.code || '', name: item?.name || '', description: item?.description || '',
    component_type: item?.component_type || presetType.value || 'EARNING', calculation_type: item?.calculation_type || 'FIXED',
    formula: item?.formula || '', reference_component_id: item?.reference_component_id || null,
    is_taxable: item?.is_taxable ?? true, is_bpjs_base: item?.is_bpjs_base ?? false,
    is_recurring: item?.is_recurring ?? true, is_proratable: item?.is_proratable ?? true,
    print_on_salary_structure: item?.print_on_salary_structure ?? true,
    is_pph21_component: item?.is_pph21_component ?? false, is_pph21_deductible: item?.is_pph21_deductible ?? false,
    display_order: item?.display_order ?? 100,
    status: item?.status || 'ACTIVE'
  }
  dialogVisible.value = true
}
function resetForm() {
  form.value = { code: '', name: '', description: '', component_type: 'EARNING', calculation_type: 'FIXED', formula: '', reference_component_id: null, is_taxable: true, is_bpjs_base: false, is_recurring: true, is_proratable: true, print_on_salary_structure: true, is_pph21_component: false, is_pph21_deductible: false, display_order: 100, status: 'ACTIVE' }
  errors.value = {}; editing.value = false; editingId.value = null; presetType.value = ''; formulaStatus.value = ''; formulaValid.value = false
}
function insertMention(opt) {
  const m = mentionQuery.value
  if (!m) return
  const text = form.value.formula || ''
  form.value.formula = text.slice(0, m.idx) + opt.code + text.slice(m.idx + 1 + m.q.length)
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
      is_pph21_component: form.value.is_pph21_component, is_pph21_deductible: form.value.is_pph21_deductible,
      display_order: form.value.display_order, status: form.value.status
    }
    if (editing.value) {
      await api.put(`/api/v1/tenant/payroll/salary-components/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/salary-components', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_created'), life: 3000 })
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
    await api.delete(`/api/v1/tenant/payroll/salary-components/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadData()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(loadData)
</script>
