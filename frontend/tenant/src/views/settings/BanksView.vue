<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <IconField><InputIcon class="pi pi-search" /><InputText v-model="searchQuery" size="small" /></IconField>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('banks.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable v-else :value="filteredItems" paginator :rows="15" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="sort_order" :sortOrder="1">
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-building-column text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('banks.empty_title') }}</p></div></template>
      <Column field="code" :header="t('banks.code')" sortable style="width:120px"><template #body="{data}"><Tag :value="data.code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column field="name" :header="t('banks.name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
      <Column field="sort_order" :header="t('banks.sort_order')" sortable style="width:100px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('banks.edit') : t('banks.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4"><div><h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5"><i class="pi pi-building-column text-indigo-400 text-sm"></i>{{ editing ? t('banks.edit') : t('banks.new') }}</h3>
        <div class="space-y-2">
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('banks.code') }} <span class="text-red-500">*</span></label><InputText v-model="form.code" class="!w-full" :class="{'p-invalid':errors?.code}" maxlength="20" autofocus :placeholder="t('banks.code')" /><small v-if="errors?.code" class="text-red-500 text-xs mt-1 block">{{ errors.code }}</small></div>
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('banks.name') }} <span class="text-red-500">*</span></label><InputText v-model="form.name" class="!w-full" :class="{'p-invalid':errors?.name}" maxlength="255" :placeholder="t('banks.name')" /><small v-if="errors?.name" class="text-red-500 text-xs mt-1 block">{{ errors.name }}</small></div>
          <div><label class="block text-sm font-medium text-gray-600 dark:text-gray-300 mb-1">{{ t('banks.sort_order') }}</label><InputNumber v-model="form.sort_order" class="!w-full" :min="0" /></div>
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
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputText from 'primevue/inputtext'; import InputIcon from 'primevue/inputicon'; import IconField from 'primevue/iconfield'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import ConfirmDialog from 'primevue/confirmdialog'; import SkeletonTable from '@/components/SkeletonTable.vue'
const { t } = useI18n(); const toast = useToast(); const confirm = useConfirm()
const items = ref([]); const loading = ref(false); const searchQuery = ref('')
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({}); const form = ref({ code: '', name: '', sort_order: 0 })
const skeletonColumns = [{type:'tag',width:'w-20',headerWidth:'w-16'},{type:'text',width:'w-40',headerWidth:'w-16'},{type:'text',width:'w-12',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]
const filteredItems = computed(() => { if(!searchQuery.value) return items.value; const q = searchQuery.value.toLowerCase(); return items.value.filter(i => i.code?.toLowerCase().includes(q) || i.name?.toLowerCase().includes(q)) })
async function loadData() { loading.value = true; try { const res = await api.get('/api/v1/tenant/settings/banks?per_page=200'); items.value = res.data?.data?.data || res.data?.data || [] } catch(e) { toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) } finally { loading.value = false } }
function openDialog(item) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={code:item?.code||'',name:item?.name||'',sort_order:item?.sort_order||0}; dialogVisible.value=true }
function resetForm() { form.value={code:'',name:'',sort_order:0}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.code?.trim()){errors.value={code:[t('form.required')]};return} if(!form.value.name?.trim()){errors.value={name:[t('form.required')]};return}; saving.value=true; try { const payload={code:form.value.code,name:form.value.name,sort_order:form.value.sort_order||0}; if(editing.value){await api.put(`/api/v1/tenant/settings/banks/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('banks.updated'),life:3000})}else{await api.post('/api/v1/tenant/settings/banks',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('banks.created'),life:3000})}; dialogVisible.value=false; await loadData() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
function confirmDelete(item) { confirm.require({header:t('banks.confirm_delete_title'),message:t('banks.confirm_delete',{name:item.name}),icon:'pi pi-exclamation-triangle',rejectLabel:t('common.cancel'),acceptLabel:t('common.delete'),rejectClass:'p-button-outlined p-button-secondary',acceptClass:'p-button-danger',accept:async()=>{try{await api.delete(`/api/v1/tenant/settings/banks/${item.id}`);toast.add({severity:'success',summary:t('message.success'),detail:t('banks.deleted'),life:3000});await loadData()}catch(e){toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})}}}) }
onMounted(loadData)
</script>