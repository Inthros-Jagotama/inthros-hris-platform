<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('ters.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
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
      sortField="group"
      :sortOrder="1"
    >
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-calculator text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('ters.empty_title') }}</p></div></template>
      <Column field="group" :header="t('ters.group')" sortable style="width:80px"><template #body="{data}"><Tag :value="data.group" :severity="groupSeverity(data.group)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column field="bruto_min" :header="t('ters.bruto_min')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ formatAmount(data.bruto_min) }}</span></template></Column>
      <Column field="bruto_max" :header="t('ters.bruto_max')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 font-mono text-xs">{{ formatAmount(data.bruto_max) }}</span></template></Column>
      <Column field="rate" :header="t('ters.rate')" sortable style="width:120px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.rate }}%</span></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('ters.edit') : t('ters.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4"><div class="space-y-2">
          <FormRow :label="t('ters.group')" required :errors="errors?.group">
            <SelectLabel v-model="form.group" :options="groupOptions" option-value="value" option-label="label" :placeholder="t('ters.select_group')" :class="{'p-invalid':errors?.group}" :showClear="true" />
          </FormRow>
          <div class="grid grid-cols-2 gap-3">
            <FormRow :label="t('ters.bruto_min')" :errors="errors?.bruto_min">
              <InputNumber v-model="form.bruto_min" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
            </FormRow>
            <FormRow :label="t('ters.bruto_max')" :errors="errors?.bruto_max">
              <InputNumber v-model="form.bruto_max" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
            </FormRow>
          </div>
          <FormRow :label="t('ters.rate') + ' (%)'" required :errors="errors?.rate">
            <InputNumber v-model="form.rate" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-between"><div class="flex items-center gap-2 ml-auto"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></div></template>
    </Dialog>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
const { t } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null); const form = ref({ group: '', bruto_min: null, bruto_max: null, rate: 0 })
const groupOptions = [{ value: 'A', label: 'Group A — TK/0, TK/1, TK/2, TK/3' },{ value: 'B', label: 'Group B — K/0, K/1, K/2, K/3' },{ value: 'C', label: 'Group C — K/I/0, K/I/1, K/I/2, K/I/3' }]
function groupSeverity(g) { return { A: 'info', B: 'success', C: 'warn' }[g] || 'info' }
const skeletonColumns = [{type:'tag',width:'w-16',headerWidth:'w-12'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-16',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]
function formatAmount(val) { if (val === null || val === undefined) return '∞'; return new Intl.NumberFormat('id-ID').format(val) }
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/ters', {
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
  currentPage.value = event.page + 1  // PrimeVue uses 0-indexed page
  perPage.value = event.rows
  loadData()
}
function openDialog(item) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={group:item?.group||'',bruto_min:item?.bruto_min||null,bruto_max:item?.bruto_max||null,rate:item?.rate||0}; dialogVisible.value=true }
function resetForm() { form.value={group:'',bruto_min:null,bruto_max:null,rate:0}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.group){errors.value={group:[t('form.required')]};return}; saving.value=true; try { const payload={group:form.value.group,bruto_min:form.value.bruto_min||null,bruto_max:form.value.bruto_max||null,rate:form.value.rate}; if(editing.value){await api.put(`/api/v1/tenant/settings/ters/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('ters.updated'),life:3000})}else{await api.post('/api/v1/tenant/settings/ters',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('ters.created'),life:3000})}; dialogVisible.value=false; await loadData() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/ters/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('ters.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(loadData)
</script>
