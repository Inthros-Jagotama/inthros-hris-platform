<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('ptkps.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
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
      sortField="name"
      :sortOrder="1"
    >
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-receipt text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('ptkps.empty_title') }}</p></div></template>
      <Column field="name" :header="t('ptkps.name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
      <Column field="group" :header="t('ptkps.group')" sortable style="width:80px"><template #body="{data}"><Tag :value="data.group" :severity="groupSeverity(data.group)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column field="ptkp" :header="t('ptkps.ptkp')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium font-mono text-xs">{{ formatAmount(data.ptkp) }}</span></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('ptkps.edit') : t('ptkps.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4"><div><h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5"><i class="pi pi-receipt text-indigo-400 text-sm"></i>{{ editing ? t('ptkps.edit') : t('ptkps.new') }}</h3>
        <div class="space-y-2">
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('ptkps.name') }} <span class="text-red-500">*</span></label><InputText v-model="form.name" class="!w-full" :class="{'p-invalid':errors?.name}" maxlength="255" autofocus :placeholder="t('ptkps.name')" /><small v-if="errors?.name" class="text-red-500 text-xs mt-1 block">{{ errors.name }}</small></div>
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('ptkps.group') }} <span class="text-red-500">*</span></label>
            <Select v-model="form.group" :options="groupOptions" optionValue="value" optionLabel="label" :placeholder="t('ptkps.select_group')" class="!w-full" :class="{'p-invalid':errors?.group}" :showClear="true" />
            <small v-if="errors?.group" class="text-red-500 text-xs mt-1 block">{{ errors.group }}</small></div>
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('ptkps.ptkp') }} <span class="text-red-500">*</span></label><InputNumber v-model="form.ptkp" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" /><small v-if="errors?.ptkp" class="text-red-500 text-xs mt-1 block">{{ errors.ptkp }}</small></div>
        </div>
      </div></div>
      <template #footer><div class="flex items-center justify-between"><div class="flex items-center gap-2 ml-auto"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></div></template>
    </Dialog>
    <ConfirmDialog />
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { useConfirm } from 'primevue/useconfirm'; import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputText from 'primevue/inputtext'; import InputNumber from 'primevue/inputnumber'; import Select from 'primevue/select'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import ConfirmDialog from 'primevue/confirmdialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
const { t } = useI18n(); const toast = useToast(); const confirm = useConfirm()
const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({}); const form = ref({ name: '', group: '', ptkp: 0 })
const groupOptions = [{ value: 'A', label: 'Group A — TK/0, TK/1, TK/2, TK/3' },{ value: 'B', label: 'Group B — K/0, K/1, K/2, K/3' },{ value: 'C', label: 'Group C — K/I/0, K/I/1, K/I/2, K/I/3' }]
function groupSeverity(g) { return { A: 'info', B: 'success', C: 'warn' }[g] || 'info' }
const skeletonColumns = [{type:'text',width:'w-40',headerWidth:'w-16'},{type:'tag',width:'w-16',headerWidth:'w-12'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]
function formatAmount(val) { if (!val) return 'Rp 0'; return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(val) }
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/ptkps', {
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
function openDialog(item) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={name:item?.name||'',group:item?.group||'',ptkp:item?.ptkp||0}; dialogVisible.value=true }
function resetForm() { form.value={name:'',group:'',ptkp:0}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.name?.trim()){errors.value={name:[t('form.required')]};return} if(!form.value.group){errors.value={group:[t('form.required')]};return}; saving.value=true; try { const payload={name:form.value.name,group:form.value.group,ptkp:form.value.ptkp}; if(editing.value){await api.put(`/api/v1/tenant/settings/ptkps/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('ptkps.updated'),life:3000})}else{await api.post('/api/v1/tenant/settings/ptkps',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('ptkps.created'),life:3000})}; dialogVisible.value=false; await loadData() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
function confirmDelete(item) { confirm.require({header:t('ptkps.confirm_delete_title'),message:t('ptkps.confirm_delete',{name:item.name}),icon:'pi pi-exclamation-triangle',rejectLabel:t('common.cancel'),acceptLabel:t('common.delete'),rejectClass:'p-button-outlined p-button-secondary',acceptClass:'p-button-danger',accept:async()=>{try{await api.delete(`/api/v1/tenant/settings/ptkps/${item.id}`);toast.add({severity:'success',summary:t('message.success'),detail:t('ptkps.deleted'),life:3000});await loadData()}catch(e){toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})}}}) }
onMounted(loadData)
</script>
