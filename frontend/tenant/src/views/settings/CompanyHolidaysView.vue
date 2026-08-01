<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('company_holidays.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable v-else :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" sortField="holiday_date" :sortOrder="1">
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-calendar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('company_holidays.empty_title') }}</p></div></template>
      <Column field="holiday_date" :header="t('company_holidays.holiday_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatDate(data.holiday_date, locale) }}</span></template></Column>
      <Column field="name" :header="t('company_holidays.name')" sortable><template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
      <Column field="description" :header="t('company_holidays.desc')" sortable><template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.description || '—' }}</span></template></Column>
      <Column field="is_active" :header="t('common.status')" sortable style="width:110px"><template #body="{data}"><Tag :value="data.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="data.is_active ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
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
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import ToggleSwitch from 'primevue/toggleswitch'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'
import { formatDate } from '@/utils/formatDate'
const { t, locale } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(15)
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null); const form = ref({ holiday_date: '', name: '', description: '', is_active: true })
const skeletonColumns = [{type:'text',width:'w-28',headerWidth:'w-16'},{type:'text',width:'w-32',headerWidth:'w-16'},{type:'text',width:'w-24',headerWidth:'w-16'},{type:'text',width:'w-16',headerWidth:'w-16'},{type:'icons',count:2,headerWidth:'w-16'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

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
function openDialog(item) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={holiday_date:item?.holiday_date||'',name:item?.name||'',description:item?.description||'',is_active:item?.is_active!==undefined?item.is_active:true}; dialogVisible.value=true }
function resetForm() { form.value={holiday_date:'',name:'',description:'',is_active:true}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.holiday_date){errors.value={holiday_date:[t('form.required')]};return} if(!form.value.name?.trim()){errors.value={name:[t('form.required')]};return}; saving.value=true; try { const payload={holiday_date:form.value.holiday_date,name:form.value.name,description:form.value.description||'',is_active:form.value.is_active}; if(editing.value){await api.put(`/api/v1/tenant/settings/company-holidays/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('company_holidays.updated'),life:3000})}else{await api.post('/api/v1/tenant/settings/company-holidays',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('company_holidays.created'),life:3000})}; dialogVisible.value=false; await loadData() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
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
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}
onMounted(loadData)
</script>
