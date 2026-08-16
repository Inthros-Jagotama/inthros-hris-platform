<template>
  <div class="space-y-4">
    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div v-for="n in 2" :key="n" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div class="h-4 w-48 bg-gray-200 dark:bg-gray-700 rounded animate-pulse mb-4"></div>
        <div class="space-y-3">
          <div v-for="m in 5" :key="m" class="h-9 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
      </div>
    </div>

    <template v-else>
      <!-- Effective date for new overrides -->
      <div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg px-4 py-3">
        <FormRow :label="t('payroll.effective_start_date')" class="!mb-0 w-80">
          <DateInput v-model="defaultEffectiveDate" class="!w-full" />
        </FormRow>
      </div>

      <!-- All master components grouped by type -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div
          v-for="group in groups"
          :key="group.type"
          class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden"
        >
          <div class="flex items-center justify-between gap-2 px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
            <div class="flex items-center gap-2 min-w-0">
              <i :class="group.icon" class="text-sm text-gray-400 dark:text-gray-500"></i>
              <span class="text-sm font-semibold text-surface-800 dark:text-surface-0 truncate">{{ group.label }}</span>
              <Tag :value="String(group.items.length)" :severity="group.severity" class="!text-xs !px-1.5 !py-0.5" />
            </div>
          </div>
          <div class="p-4 space-y-4">
            <div v-if="group.items.length === 0" class="py-6 text-center text-xs text-gray-400 dark:text-gray-500">
              {{ t('payroll.components_empty_group', { type: group.label }) }}
            </div>
            <FormRow v-for="item in group.items" :key="item.id" :label="item.name">
              <div class="flex items-center gap-2 min-w-0">
                <InputNumber v-model="amounts[item.id]" class="!flex-1" :min="0" size="small" mode="currency" currency="IDR" locale="id-ID" :input-id="'salary-amount-' + item.id" />
                <Button
                  v-if="hasOverride(item)"
                  icon="pi pi-trash"
                  size="small"
                  severity="danger"
                  text
                  v-tooltip.top="t('common.delete')"
                  class="shrink-0"
                  @click="confirmDelete(item)"
                />
              </div>
            </FormRow>
          </div>
        </div>
      </div>

      <!-- Save all changes -->
      <div class="flex items-center justify-end gap-3 pt-2">
        <span v-if="hasChanges" class="text-xs text-gray-400 dark:text-gray-500">
          <i class="pi pi-info-circle mr-1"></i>{{ t('payroll.salary_structure_inline_hint') }}
        </span>
        <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="saveAll" />
      </div>
    </template>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.employee_components')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import DateInput from '@/components/DateInput.vue'

const props = defineProps({
  employeeId: { type: String, default: null }
})
const emit = defineEmits(['save'])

const { t, locale } = useI18n()
const toast = useToast()

const loading = ref(true)
const components = ref([])        // semua komponen master
const overrides = ref([])         // employee components milik karyawan ini
const amounts = reactive({})      // component_id -> nilai input
const defaultEffectiveDate = ref(today())
const saving = ref(false)
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

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
  items: components.value
    .filter(c => c.component_type === type)
    .sort((a, b) => (a.display_order ?? 0) - (b.display_order ?? 0))
})))

function overrideOf(component) {
  return overrides.value.find(o => o.salary_component_id === component.id) || null
}
function hasOverride(component) {
  return !!overrideOf(component)
}
// Ada nilai yang berubah / belum disimpan (menentukan tampilnya hint & status tombol)
const hasChanges = computed(() => components.value.some(c => {
  const val = amounts[c.id] ?? null
  const existing = overrideOf(c)
  if (existing) {
    return val === null || Number(val) !== Number(existing.amount)
  }
  return val !== null && val !== undefined && Number(val) > 0
}))
const deleteMessage = computed(() => deleteTarget.value
  ? `${deleteTarget.value.code || ''} ${deleteTarget.value.name || ''}`
  : t('common.no_data'))

function today() {
  const d = new Date()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${d.getFullYear()}-${m}-${day}`
}

function typeLabel(v) { const key = `payroll.component_type_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function typeSeverity(v) { return { EARNING: 'success', DEDUCTION: 'danger', EMPLOYER_CONTRIBUTION: 'warn', INFORMATION: 'info' }[v] || 'secondary' }

async function loadComponents() {
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
}

async function loadOverrides() {
  const all = []
  let page = 1
  while (true) {
    const res = await api.get('/api/v1/tenant/payroll/salary-employee-components', { params: { page, per_page: 100 } })
    const rows = res.data?.data || []
    all.push(...rows)
    if (!rows.length || all.length >= (res.data?.total || 0)) break
    page++
  }
  overrides.value = props.employeeId ? all.filter(x => x.employee_id === props.employeeId) : all
  // Inisialisasi nilai input dari override yang ada
  for (const o of overrides.value) {
    amounts[o.salary_component_id] = o.amount
  }
}

async function loadData() {
  loading.value = true
  try {
    await Promise.all([loadComponents(), loadOverrides()])
    // Pastikan setiap komponen punya slot nilai (null = belum ada override)
    for (const c of components.value) {
      if (!(c.id in amounts)) amounts[c.id] = null
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

// Simpan semua perubahan sekaligus: buat yang baru, perbarui yang berubah.
async function saveAll() {
  saving.value = true
  const created = []
  const updated = []
  try {
    for (const c of components.value) {
      const val = amounts[c.id] ?? null
      const existing = overrideOf(c)
      const amount = Number(val || 0)
      if (existing) {
        if (val !== null && amount !== Number(existing.amount)) {
          await api.put(`/api/v1/tenant/payroll/salary-employee-components/${existing.id}`, {
            employee_id: props.employeeId,
            salary_component_id: c.id,
            amount,
            currency_code: 'IDR',
            source_type: existing.source_type || 'MANUAL',
            effective_start_date: existing.effective_start_date,
            effective_end_date: existing.effective_end_date || null,
            status: existing.status || 'ACTIVE',
            notes: existing.notes || null
          })
          updated.push(c.name)
        }
      } else if (val !== null && val !== undefined && amount > 0) {
        await api.post('/api/v1/tenant/payroll/salary-employee-components', {
          employee_id: props.employeeId,
          salary_component_id: c.id,
          amount,
          currency_code: 'IDR',
          source_type: 'MANUAL',
          effective_start_date: defaultEffectiveDate.value,
          effective_end_date: null,
          status: 'ACTIVE',
          notes: null
        })
        created.push(c.name)
      }
    }
    if (created.length || updated.length) {
      toast.add({
        severity: 'success',
        summary: t('message.success'),
        detail: t('payroll.salary_structure_saved', { created: created.length, updated: updated.length }),
        life: 3000
      })
    } else {
      toast.add({ severity: 'info', summary: t('message.warning'), detail: t('payroll.no_changes'), life: 3000 })
    }
    await loadOverrides()
    emit('save')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    saving.value = false
  }
}

function confirmDelete(component) { deleteTarget.value = component; deleteError.value = ''; deleteDialogVisible.value = true }

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    const existing = overrideOf(deleteTarget.value)
    await api.delete(`/api/v1/tenant/payroll/salary-employee-components/${existing.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.component_deleted'), life: 3000 })
    deleteDialogVisible.value = false
    amounts[deleteTarget.value.id] = null
    await loadOverrides()
    emit('save')
  } catch (e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}

onMounted(() => { loadData() })
</script>
