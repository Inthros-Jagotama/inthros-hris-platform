<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <div><h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_competencies') }}</h2><p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_description') }}</p></div>
      <Button :label="t('common.create')" icon="pi pi-plus" size="small" @click="openCreate()" />
    </div>
    <DataTableSection :items="items" :loading="loading" :total="total" :columns="cols" entity="potency-competencies" :org-id="orgId" :on-load="loadData" @edit="openEdit" @delete="confirmDelete">
      <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400"><i class="pi pi-star text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('job_management.empty_potency') }}</p></div></template>
    </DataTableSection>
    <DialogForm v-model:visible="dialogVisible" :title="editing ? t('common.edit') : t('common.create')" :saving="saving" :errors="errors" @save="handleSave" @cancel="dialogVisible=false">
      <FormRow :label="t('job_management.competency')" required :errors="errors?.competency_id"><SelectLabel v-model="form.competency_id" :options="competencyOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear /></FormRow>
      <FormRow :label="t('job_management.value_ref')" :errors="errors?.job_management_value_id"><SelectLabel v-model="form.job_management_value_id" :options="allOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear /></FormRow>
      <FormRow :label="t('job_management.weight')" :errors="errors?.weight"><InputNumber v-model="form.weight" :min="0" :max="100" :class="{'p-invalid':errors?.weight}" class="w-full" size="small" /></FormRow>
    </DialogForm>
    <ConfirmDeleteDialog v-model:visible="deleteVisible" :loading="deleting" :error-msg="deleteError" @confirm="handleDelete" @cancel="deleteVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useI18n } from '@/composables/useI18n'; import { useToast } from 'primevue/usetoast'; import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'
import FormRow from '@/components/FormRow.vue'; import SelectLabel from '@/components/SelectLabel.vue'
import DataTableSection from './DataTableSection.vue'; import DialogForm from './DialogForm.vue'; import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
const props = defineProps({ orgId: String, jobValueMap: Object, competencyOptions: Array }); const emit = defineEmits(['saved'])
const { t } = useI18n(); const toast = useToast(); const apiBase = '/api/v1/tenant/job-management/potency-competencies'
const items=ref([]); const loading=ref(false); const total=ref(0); const dialogVisible=ref(false); const editing=ref(false); const editId=ref('')
const saving=ref(false); const errors=ref({}); const deleteVisible=ref(false); const deleting=ref(false); const deleteError=ref(''); const deleteTarget=ref(null)
const form=ref({ competency_id:'', job_management_value_id:'', weight:null })
const allOptions=computed(()=>Object.values(props.jobValueMap||{}).flat())
const cols=computed(()=>[{field:'_competency',header:t('job_management.competency')},{field:'weight',header:t('job_management.weight')}])
async function loadData(page,perPage){ loading.value=true; try{ const r=await api.get(apiBase,{params:{page,per_page:perPage,organization_id:props.orgId}}); items.value=(r.data?.data||[]).map(d=>{ const c=props.competencyOptions?.find(o=>o.value===d.competency_id); return {...d,_competency:c?.label||d.competency_id}}); total.value=r.data?.total||0 }catch(e){ toast.add({severity:'error',detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) }finally{ loading.value=false } }
function openCreate(){ editing.value=false; editId.value=''; form.value={competency_id:'',job_management_value_id:'',weight:null}; errors.value={}; dialogVisible.value=true }
function openEdit(d){ editing.value=true; editId.value=d.id; form.value={competency_id:d.competency_id||'',job_management_value_id:d.job_management_value_id||'',weight:d.weight??null}; errors.value={}; dialogVisible.value=true }
async function handleSave(){ saving.value=true; errors.value={}; try{ const p={...form.value,organization_id:props.orgId}; if(editing.value) await api.put(`${apiBase}/${editId.value}`,p); else await api.post(apiBase,p); dialogVisible.value=false; emit('saved'); toast.add({severity:'success',detail:t('message.saved'),life:2000}); loadData(1,15) }catch(e){ const fe=getValidationErrors(e); if(Object.keys(fe).length) errors.value=fe; else toast.add({severity:'error',detail:e.response?.data?.error?.message||t('message.operation_failed'),life:4000}) }finally{ saving.value=false } }
function confirmDelete(d){ deleteTarget.value=d; deleteError.value=''; deleteVisible.value=true }
async function handleDelete(){ if(!deleteTarget.value) return; deleting.value=true; deleteError.value=''; try{ await api.delete(`${apiBase}/${deleteTarget.value.id}`); deleteVisible.value=false; emit('saved'); toast.add({severity:'success',detail:t('message.deleted'),life:2000}); loadData(1,15) }catch(e){ deleteError.value=e.response?.data?.error?.message||t('message.operation_failed') }finally{ deleting.value=false } }
</script>
