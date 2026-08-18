<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <div class="flex items-center gap-1 border border-gray-200 dark:border-gray-700 rounded-lg p-0.5 bg-white dark:bg-gray-800">
          <Button
            :label="t('company_holidays.view_table')"
            icon="pi pi-table"
            :severity="viewMode === 'table' ? 'primary' : 'secondary'"
            :outlined="viewMode !== 'table'"
            size="small"
            @click="viewMode = 'table'"
          />
          <Button
            :label="t('company_holidays.view_calendar')"
            icon="pi pi-calendar"
            :severity="viewMode === 'calendar' ? 'primary' : 'secondary'"
            :outlined="viewMode !== 'calendar'"
            size="small"
            @click="viewMode = 'calendar'"
          />
        </div>
        <Button :label="t('company_holidays.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>
    <SkeletonTable v-if="loading && viewMode === 'table'" :columns="skeletonColumns" :rows="8" />
    <!-- Table View -->
    <DataTable v-else-if="viewMode === 'table'" :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="holiday_date" :sortOrder="1">
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-calendar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('company_holidays.empty_title') }}</p></div></template>
      <Column field="holiday_date" :header="t('company_holidays.holiday_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatDate(data.holiday_date, locale) }}</span></template></Column>
      <Column field="name" :header="t('company_holidays.name')" sortable><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
      <Column field="description" :header="t('company_holidays.desc')" sortable><template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.description || '—' }}</span></template></Column>
      <Column field="is_active" :header="t('common.status')" sortable style="width:110px"><template #body="{data}"><Tag :value="data.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="data.is_active ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
    <!-- Calendar View -->
    <div v-else-if="viewMode === 'calendar'" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <div class="lg:col-span-2 border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-3">
        <DatePicker v-model="calendarDate" inline :locale="primeLocale" class="w-full !border-0" @date-select="onCalendarDateSelect" @month-change="onMonthChange" @year-change="onMonthChange">
          <template #date="{ date }">
            <div class="relative w-full h-full flex flex-col items-center justify-center">
              <span
                class="flex items-center justify-center w-7 h-7 rounded-full transition-colors"
                :class="holidayMap[toDateKey(date)] ? 'bg-rose-500 text-white font-semibold' : ''"
              >{{ date.day }}</span>
            </div>
          </template>
        </DatePicker>
      </div>
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('company_holidays.calendar_legend') }}</h3>
          <span class="text-xs text-gray-400 dark:text-gray-500">{{ monthHolidays.length }} {{ t('common.items') }}</span>
        </div>
        <div v-if="monthHolidays.length === 0" class="text-sm text-gray-400 dark:text-gray-500 py-6 text-center">
          <i class="pi pi-calendar-plus text-2xl mb-2 opacity-50"></i>
          <p>{{ t('company_holidays.calendar_empty_month') }}</p>
        </div>
        <ul v-else class="space-y-2">
          <li
            v-for="h in monthHolidays"
            :key="h.id"
            class="flex items-center gap-2 text-sm cursor-pointer hover:bg-rose-50 dark:hover:bg-rose-900/20 rounded-md px-2 py-1.5 transition-colors"
            @click="openDialog(h)"
          >
            <i class="pi pi-calendar text-xs text-rose-500 shrink-0"></i>
            <span class="text-gray-500 dark:text-gray-400 text-xs shrink-0 w-24">{{ formatDate(h.holiday_date, locale) }}</span>
            <span class="text-navy-800 dark:text-gray-100 font-medium truncate">{{ h.name }}</span>
            <span v-if="!h.is_active" class="text-xs text-gray-400 shrink-0">({{ t('common_status.inactive') }})</span>
          </li>
        </ul>
      </div>
    </div>
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('company_holidays.edit') : t('company_holidays.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <div class="space-y-2">
          <FormRow :label="t('company_holidays.holiday_date')" required :errors="errors?.holiday_date">
            <DateInput v-model="form.holiday_date" :placeholder="t('company_holidays.holiday_date_placeholder')" :class="{'p-invalid':errors?.holiday_date}" />
          </FormRow>
          <FormRow :label="t('company_holidays.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="200" autofocus :placeholder="t('company_holidays.name')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
          <FormRow :label="t('company_holidays.desc')" :errors="errors?.description">
            <TextInput v-model="form.description" maxlength="500" :placeholder="t('company_holidays.desc_placeholder')" />
          </FormRow>
          <div class="flex items-center justify-between pt-2"><label class="text-sm font-medium text-gray-600 dark:text-gray-300">{{ t('company_holidays.is_active') }}</label><ToggleSwitch v-model="form.is_active" /></div>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-between"><div class="flex items-center gap-2 ml-auto"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></div></template>
    </Dialog>
    <!-- Delete Confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('company_holidays.confirm_delete_title')"
      :message="t('company_holidays.confirm_delete', { name: deleteTarget?.name || '' })"
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
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import ToggleSwitch from 'primevue/toggleswitch'; import DatePicker from 'primevue/datepicker'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import { formatDate } from '@/utils/formatDate'
import { getPrimeLocale } from '@/utils/primevueLocale'
const { t, locale } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null); const form = ref({ holiday_date: '', name: '', description: '', is_active: true })
const skeletonColumns = [{type:'text',width:'w-28',headerWidth:'w-16'},{type:'text',width:'w-32',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-16',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// ── View mode: table | calendar ──
const viewMode = ref('table')
const calendarDate = ref(new Date())
const primeLocale = computed(() => getPrimeLocale(locale.value))
const allHolidays = ref([])

// Normalisasi tanggal dari API → YYYY-MM-DD (tahan format apa pun: dengan waktu/TZ)
function normDateKey(s) {
  if (!s) return ''
  return String(s).slice(0, 10)
}

// Map tanggal (YYYY-MM-DD) → holiday untuk highlight di calendar
const holidayMap = computed(() => {
  const map = {}
  for (const h of allHolidays.value) {
    const key = normDateKey(h?.holiday_date)
    if (key) map[key] = h
  }
  return map
})

// Terima instance Date ATAU object slot PrimeVue #date ({ day, month 0-indexed, year })
function toDateKey(d) {
  if (!d) return ''
  let y, m0, day
  if (typeof d.getFullYear === 'function') {
    y = d.getFullYear(); m0 = d.getMonth(); day = d.getDate()
  } else {
    y = d.year; m0 = d.month; day = d.day
  }
  if (y == null || m0 == null || day == null || Number.isNaN(y) || Number.isNaN(m0) || Number.isNaN(day)) return ''
  const m = String(m0 + 1).padStart(2, '0')
  return `${y}-${m}-${String(day).padStart(2, '0')}`
}

// Bulan yang sedang dilihat di kalender — ikut panah navigasi (month-change),
// bukan hanya modelValue, agar panel legend tidak basi saat ganti bulan.
const viewYear = ref(calendarDate.value.getFullYear())
const viewMonth = ref(calendarDate.value.getMonth()) // 0-indexed

const monthHolidays = computed(() => {
  const y = viewYear.value
  const m = viewMonth.value
  return allHolidays.value
    .filter(h => {
      const parts = normDateKey(h?.holiday_date).split('-')
      return parts.length === 3 && parseInt(parts[0], 10) === y && parseInt(parts[1], 10) === m + 1
    })
    .sort((a, b) => (normDateKey(a.holiday_date) < normDateKey(b.holiday_date) ? -1 : 1))
})

async function loadCalendar() {
  try {
    const res = await api.get('/api/v1/tenant/settings/company-holidays', {
      params: { page: 1, per_page: 500 }
    })
    allHolidays.value = res.data?.data || []
  } catch(e) {
    console.error('Failed to load holidays for calendar:', e)
  }
}

// Navigasi bulan (panah) & tahun (dropdown) di kalender —
// PrimeVue month-change/year-change: { month: 1-indexed, year }
function onMonthChange(e) {
  viewYear.value = e.year
  viewMonth.value = e.month - 1
}

// Klik tanggal di calendar → ada libur: edit; belum ada: modal tambah dengan tanggal terisi
function onCalendarDateSelect(date) {
  if (!date) return
  const d = date instanceof Date ? date : new Date(date)
  viewYear.value = d.getFullYear()
  viewMonth.value = d.getMonth()
  const key = toDateKey(d)
  const h = holidayMap.value[key]
  if (h) openDialog(h)
  else openDialog(null, key)
}

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/company-holidays', {
      params: { page: currentPage.value, per_page: perPage.value }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch(e) {
    toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000})
  } finally {
    loading.value = false
  }
}
function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}
function openDialog(item, prefillDate) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={holiday_date:prefillDate||item?.holiday_date||'',name:item?.name||'',description:item?.description||'',is_active:item?.is_active!==undefined?item.is_active:true}; dialogVisible.value=true }
function resetForm() { form.value={holiday_date:'',name:'',description:'',is_active:true}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.holiday_date){errors.value={holiday_date:[t('form.required')]};return} if(!form.value.name?.trim()){errors.value={name:[t('form.required')]};return}; saving.value=true; try { const payload={holiday_date:form.value.holiday_date,name:form.value.name,description:form.value.description||'',is_active:form.value.is_active}; if(editing.value){await api.put(`/api/v1/tenant/settings/company-holidays/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('company_holidays.updated'),life:3000})      }else{await api.post('/api/v1/tenant/settings/company-holidays',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('company_holidays.created'),life:3000})}; dialogVisible.value=false; await loadData(); await loadCalendar() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/company-holidays/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('company_holidays.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
    await loadCalendar()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}
onMounted(() => {
  loadData()
  loadCalendar()
})
</script>
