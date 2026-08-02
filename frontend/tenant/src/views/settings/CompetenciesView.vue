<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="searchQuery" :placeholder="t('common.search')" size="small" class="!pl-8 !text-sm !py-1.5 !w-64" @input="onSearchInput" />
        </IconField>
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">
          {{ totalRecords }} {{ t('common.items') }}
        </span>
      </div>
      <div class="flex items-center gap-2">
        <Button :label="t('competencies.new')" icon="pi pi-plus" size="small" @click="openDialog()" />
      </div>
    </div>
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />
    <DataTable v-else :value="items" lazy :totalRecords="totalRecords" :first="firstRecord" :rows="perPage" @page="onPage($event)" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown" :rowsPerPageOptions="[10, 15, 25, 50]" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-star text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('competencies.empty_title') }}</p></div></template>
      <Column field="name" :header="t('competencies.name')" sortable><template #body="{data}"><div class="py-0.5 pr-2"><div class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</div><p v-if="data.definition" class="text-gray-500 dark:text-gray-400 text-xs mt-1 whitespace-pre-line">{{ data.definition }}</p></div></template></Column>
      <Column field="field" :header="t('competencies.field')" sortable><template #body="{data}"><Tag v-if="data.field" :value="data.field" severity="info" class="!text-xs !px-1.5 !py-0.5" /><span v-else class="text-gray-400">-</span></template></Column>
      <Column field="cluster" :header="t('competencies.cluster')" sortable><template #body="{data}"><Tag v-if="data.cluster" :value="data.cluster" severity="secondary" class="!text-xs !px-1.5 !py-0.5" /><span v-else class="text-gray-400">-</span></template></Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" /></div></template></Column>
    </DataTable>
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('competencies.edit') : t('competencies.new')" modal :style="{ width: '600px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4"><div class="space-y-2">
          <FormRow :label="t('competencies.name')" required :errors="errors?.name">
            <TextInput v-model="form.name" maxlength="255" autofocus :placeholder="t('competencies.name')" :class="{'p-invalid':errors?.name}" />
          </FormRow>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <FormRow :label="t('competencies.field')" :errors="errors?.field">
              <TextInput v-model="form.field" maxlength="255" :placeholder="t('competencies.field')" :class="{'p-invalid':errors?.field}" />
            </FormRow>
            <FormRow :label="t('competencies.cluster')" :errors="errors?.cluster">
              <TextInput v-model="form.cluster" maxlength="255" :placeholder="t('competencies.cluster')" :class="{'p-invalid':errors?.cluster}" />
            </FormRow>
          </div>
          <FormRow :label="t('competencies.definition')" :errors="errors?.definition">
            <Textarea v-model="form.definition" rows="5" class="!w-full" :placeholder="t('competencies.definition_placeholder')" :class="{'p-invalid':errors?.definition}" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-between"><div class="flex items-center gap-2 ml-auto"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" /><Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" /></div></div></template>
    </Dialog>
    <!-- Delete Confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('competencies.confirm_delete_title')"
      :message="t('competencies.confirm_delete', { name: deleteTarget?.name || '' })"
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
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import Textarea from 'primevue/textarea'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import InputText from 'primevue/inputtext'; import IconField from 'primevue/iconfield'; import InputIcon from 'primevue/inputicon'; import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
const { t } = useI18n(); const toast = useToast(); const items = ref([]); const loading = ref(false)
const totalRecords = ref(0); const currentPage = ref(1); const perPage = ref(10)
const searchQuery = ref('')
let searchTimer = null
const dialogVisible = ref(false); const editing = ref(false); const editingId = ref(null); const saving = ref(false); const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null); const form = ref({ name: '', field: '', cluster: '', definition: '' })
const skeletonColumns = [{type:'text',width:'w-72',headerWidth:'w-32'},{type:'tag',width:'w-24',headerWidth:'w-20'},{type:'tag',width:'w-24',headerWidth:'w-20'},{type:'icons',count:2,headerWidth:'w-16'}]
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const q = searchQuery.value?.trim()
    if (q) params.search = q
    const res = await api.get('/api/v1/tenant/settings/competencies', { params })
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
function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
    loadData()
  }, 400)
}
function openDialog(item) { editing.value=!!item; editingId.value=item?.id||null; errors.value={}; form.value={name:item?.name||'',field:item?.field||'',cluster:item?.cluster||'',definition:item?.definition||''}; dialogVisible.value=true }
function resetForm() { form.value={name:'',field:'',cluster:'',definition:''}; errors.value={}; editing.value=false; editingId.value=null }
async function handleSave() { errors.value={}; if(!form.value.name?.trim()){errors.value={name:[t('form.required')]};return}; saving.value=true; try { const payload={name:form.value.name,field:form.value.field||null,cluster:form.value.cluster||null,definition:form.value.definition||null}; if(editing.value){await api.put(`/api/v1/tenant/settings/competencies/${editingId.value}`,payload);toast.add({severity:'success',summary:t('message.success'),detail:t('competencies.updated'),life:3000})}else{await api.post('/api/v1/tenant/settings/competencies',payload);toast.add({severity:'success',summary:t('message.success'),detail:t('competencies.created'),life:3000})}; dialogVisible.value=false; await loadData() } catch(e) { const fe=getValidationErrors(e); if(Object.keys(fe).length>0){errors.value=fe}else{toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000})} } finally { saving.value=false } }
function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/competencies/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('competencies.deleted'), life: 3000 })
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
