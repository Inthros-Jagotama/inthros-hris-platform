<template>
  <div class="space-y-4">
    <!-- Mapping type ↔ cluster kompetensi (tipe yang memakai filter cluster: technical & managerial) -->
    <JobValueClusterCard v-if="['technical', 'managerial'].includes(type)" :type="type" />

    <!-- Header — title/desc dipindah ke page title layout (AppLayout) -->
    <div class="flex items-center justify-end mb-4">
      <Button :label="t('common.create')" icon="pi pi-plus" size="small" @click="openCreateDialog()" />
    </div>

    <!-- DataTable -->
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
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-sliders-h text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('job_values.empty_title') }}</p>
          <p class="text-xs mt-1">{{ t('job_values.empty_hint') }}</p>
        </div>
      </template>
      <Column field="level" :header="t('job_values.level')" sortable style="width:80px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.level || '-' }}</span></template>
      </Column>
      <!-- Kolom description menampilkan nama pendidikan (auto-fill) untuk tipe education;
           fallback ke nama dari ref_id utk data lama yg description-nya kosong -->
      <Column field="descriptions" :header="t('job_values.descriptions')" sortable>
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100">{{ data.descriptions || educationName(data.ref_id) || '-' }}</span></template>
      </Column>
      <Column field="note" :header="t('job_values.note')" sortable>
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.note || '-' }}</span></template>
      </Column>
      <Column field="sort" :header="t('job_values.sort')" sortable style="width:80px">
        <template #body="{data}"><span class="text-gray-500 dark:text-gray-400">{{ data.sort || 0 }}</span></template>
      </Column>
      <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEditDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Create/Edit Dialog — type terkunci ke tipe section ini -->
    <Dialog v-model:visible="dialogVisible" :header="editing ? t('job_values.edit') : t('job_values.new')" modal :style="{ width: '520px' }" :closable="true" @hide="resetForm">
      <div class="space-y-4">
        <FormRow :label="t('job_values.type')" required>
          <InputText :model-value="typeLabel(type)" disabled class="!w-full" size="small" />
        </FormRow>
        <FormRow v-if="type === 'education'" :label="t('job_values.education')" required :errors="errors?.ref_id">
          <SelectLabel
            v-model="form.ref_id"
            :options="educationOptions"
            :placeholder="t('job_values.select_education')"
            :class="{ 'p-invalid': errors?.ref_id }"
          />
        </FormRow>
        <FormRow :label="t('job_values.level')" :errors="errors?.level">
          <InputNumber v-model="form.level" class="!w-full" :min="0" size="small" :class="{ 'p-invalid': errors?.level }" />
        </FormRow>
        <!-- Untuk tipe education, description tidak ditampilkan di form —
             otomatis terisi teks dari pendidikan yang dipilih saat simpan -->
        <FormRow v-if="type !== 'education'" :label="t('job_values.descriptions')" :errors="errors?.descriptions">
          <Textarea v-model="form.descriptions" class="!w-full" :class="{ 'p-invalid': errors?.descriptions }" :placeholder="t('job_values.descriptions_placeholder')" rows="3" />
        </FormRow>
        <FormRow :label="t('job_values.note')" :errors="errors?.note">
          <Textarea v-model="form.note" class="!w-full" :class="{ 'p-invalid': errors?.note }" :placeholder="t('job_values.note_placeholder')" rows="2" />
        </FormRow>
        <FormRow :label="t('job_values.sort')" :errors="errors?.sort">
          <InputNumber v-model="form.sort" class="!w-full" :min="0" size="small" :class="{ 'p-invalid': errors?.sort }" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="dialogVisible=false" />
          <Button :label="editing ? t('common.update') : t('common.save')" size="small" :loading="saving" :disabled="saving" @click="handleSave" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('job_values.confirm_delete_title')"
      :message="t('job_values.confirm_delete', { type: typeLabel(deleteTarget?.type || type) })"
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
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import { jobValueTypeLabel } from '@/utils/jobValues'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import SkeletonTable from '@/components/SkeletonTable.vue'
import JobValueClusterCard from './JobValueClusterCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'

const props = defineProps({
  type: { type: String, required: true }
})
const emit = defineEmits(['saved'])

const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref(null)
const saving = ref(false)
const errors = ref({})
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const form = ref({ level: null, descriptions: '', note: '', sort: null, ref_id: null })

// ── Data pendidikan (master settings) untuk tipe 'education' ──
const educations = ref([])
const educationOptions = computed(() =>
  educations.value.map(e => ({ label: e.name, value: e.id }))
)
const educationNameMap = computed(() => {
  const m = {}
  for (const e of educations.value) m[e.id] = e.name
  return m
})
function educationName(id) {
  return id ? educationNameMap.value[id] || '' : ''
}

async function loadEducations() {
  try {
    const res = await api.get('/api/v1/tenant/settings/educations', { params: { page: 1, per_page: 200 } })
    educations.value = res.data?.data || []
  } catch {
    educations.value = []
  }
}

function typeLabel(value) {
  return jobValueTypeLabel(t, value)
}

const skeletonColumns = [
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'text', width: 'w-48', headerWidth: 'w-24' },
  { type: 'text', width: 'w-32', headerWidth: 'w-20' },
  { type: 'text', width: 'w-12', headerWidth: 'w-12' },
  { type: 'icons', count: 2, headerWidth: 'w-16' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// ── Load Job Values untuk tipe section ini (server-side filter ?type=) ──
async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/job-management/values', {
      params: { page: currentPage.value, per_page: perPage.value, type: props.type }
    })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

// Container merender komponen dgn :key="selectedType" → remount otomatis
// saat tipe berubah, jadi tidak perlu watch props.type.

function resetForm() {
  form.value = { level: null, descriptions: '', note: '', sort: null, ref_id: null }
  errors.value = {}
  editing.value = false
  editingId.value = null
}

function openCreateDialog() {
  resetForm()
  dialogVisible.value = true
}

function openEditDialog(item) {
  editing.value = true
  editingId.value = item.id
  errors.value = {}
  form.value = {
    level: item.level ?? null,
    descriptions: item.descriptions || '',
    note: item.note || '',
    sort: item.sort ?? null,
    ref_id: item.ref_id || null
  }
  dialogVisible.value = true
}

async function handleSave() {
  errors.value = {}
  saving.value = true
  try {
    const payload = {
      type: props.type,
      level: form.value.level ?? null,
      descriptions: form.value.descriptions || '',
      note: form.value.note || '',
      sort: form.value.sort ?? null
    }

    // Untuk tipe 'education': relasikan ke master pendidikan
    // ref_id = education_id, ref_type = 'educations' (polymorphic relation)
    // Hanya dikirim bila pendidikan terpilih — hindari record ref_type tanpa ref_id.
    if (props.type === 'education' && form.value.ref_id) {
      payload.ref_id = form.value.ref_id
      payload.ref_type = 'educations'
      // Description otomatis terisi teks dari pendidikan yang dipilih
      payload.descriptions = educationName(form.value.ref_id)
    } else if (props.type === 'education' && editing.value) {
      // Edit: clear relasi + description bila pendidikan dikosongkan (hindari data stale)
      payload.ref_id = ''
      payload.ref_type = ''
      payload.descriptions = ''
    }

    if (editing.value) {
      await api.put(`/api/v1/tenant/job-management/values/${editingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/job-management/values', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.created'), life: 3000 })
    }

    dialogVisible.value = false
    emit('saved')
    await loadData()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    saving.value = false
  }
}

function confirmDelete(item) {
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/job-management/values/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('job_values.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    emit('saved')
    await loadData()
  } catch(e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(async () => {
  await loadData()
  if (props.type === 'education') await loadEducations()
})
</script>
