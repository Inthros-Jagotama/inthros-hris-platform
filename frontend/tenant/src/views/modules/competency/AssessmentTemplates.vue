<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-2 flex-wrap">
      <div class="flex items-center gap-2">
        <Button icon="pi pi-arrow-left" size="small" text severity="secondary" v-tooltip.top="t('common.back')" @click="router.push('/competencies')" />
        <span v-if="totalRecords > 0" class="text-xs text-gray-400 dark:text-gray-500 whitespace-nowrap">{{ totalRecords }} {{ t('common.items') }}</span>
      </div>
      <div class="flex items-center gap-2 ml-auto">
        <Button :label="t('competency_360.new_template')" icon="pi pi-plus" size="small" @click="router.push('/competencies/templates/new')" />
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
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-clone text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('competency_360.templates_empty') }}</p>
        </div>
      </template>
      <Column field="code" :header="t('competency_360.code')" style="width:130px">
        <template #body="{data}"><Tag :value="data.code || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column field="name" :header="t('common.name')">
        <template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template>
      </Column>
      <Column :header="t('competency_360.competencies')" style="width:130px">
        <template #body="{data}">
          <span class="text-gray-500 dark:text-gray-400">{{ data.competencies?.length || 0 }}</span>
        </template>
      </Column>
      <Column :header="t('competency_360.rater_types')" style="width:150px">
        <template #body="{data}">
          <div class="flex flex-wrap gap-1">
            <Tag v-for="rt in (data.rater_types || [])" :key="rt.id" :value="raterTypeLabel(rt.rater_type)" severity="secondary" class="!text-[10px] !px-1 !py-0.5" />
          </div>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:100px">
        <template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'active' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
      </Column>
      <Column :header="t('common.actions')" style="width:130px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1 justify-end">
            <Button icon="pi pi-list" size="small" text severity="info" v-tooltip.left="t('competency_360.indicators')" @click="openIndicators(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="router.push('/competencies/templates/' + data.id + '/edit')" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Indicators dialog -->
    <Dialog v-model:visible="indicatorsVisible" :header="t('competency_360.template_indicators')" modal :style="{ width: '720px' }" @hide="indicatorsForm = []">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('competency_360.indicators_hint') }}</p>
      <div v-if="allIndicators.length === 0" class="text-center py-6 text-gray-400 dark:text-gray-500 text-sm">
        {{ t('competency_360.no_indicators_available') }}
      </div>
      <div v-else class="space-y-2 max-h-[50vh] overflow-y-auto pr-1">
        <div v-for="ind in allIndicators" :key="ind.id" class="flex items-start gap-2 border border-gray-200 dark:border-gray-700 rounded-lg p-2">
          <ToggleSwitch v-model="selectedIndicatorIds[ind.id]" />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-navy-800 dark:text-gray-100">{{ ind.statement }}</p>
            <p class="text-[11px] text-gray-400 dark:text-gray-500 mt-0.5">{{ ind.competency_name || '-' }}</p>
          </div>
          <div class="w-20 shrink-0">
            <TextInput v-model="indicatorWeights[ind.id]" type="number" :placeholder="t('competency_360.weight')" class="!text-xs" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="indicatorsVisible = false" />
          <Button :label="t('common.save')" size="small" :loading="savingIndicators" :disabled="savingIndicators" @click="handleSaveIndicators" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('competency_360.delete_template_title')"
      :message="t('competency_360.delete_template', { name: deleteTarget?.name || '' })"
      :loading="deleting"
      :errorMsg="deleteError"
      @confirm="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import TextInput from '@/components/TextInput.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const items = ref([])
const loading = ref(false)
const totalRecords = ref(0)
const currentPage = ref(1)
const perPage = ref(15)

const allIndicators = ref([])

const indicatorsVisible = ref(false)
const savingIndicators = ref(false)
const indicatorsTemplateId = ref(null)
const selectedIndicatorIds = ref({})
const indicatorWeights = ref({})

const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)

const skeletonColumns = [
  { type: 'tag', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-16', headerWidth: 'w-16' },
  { type: 'text', width: 'w-28', headerWidth: 'w-20' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-16' },
  { type: 'icons', count: 3, headerWidth: 'w-24' }
]

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

function raterTypeLabel(type) {
  const key = `competency_360.rater_type_${type}`
  return t(key) !== key ? t(key) : type
}

function statusLabel(status) {
  const key = `common_status.${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

async function loadReferences() {
  try {
    const res = await api.get('/api/v1/tenant/competency/indicators', { params: { per_page: 500 } })
    allIndicators.value = res.data?.data || []
  } catch {
    // fail-silent
  }
}

async function loadData() {
  loading.value = true
  try {
    const params = { page: currentPage.value, per_page: perPage.value }
    const res = await api.get('/api/v1/tenant/competency/templates', { params })
    const body = res.data
    items.value = body?.data || []
    totalRecords.value = body?.total || 0
    if (body?.page) currentPage.value = body.page
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadData()
}

async function openIndicators(template) {
  indicatorsTemplateId.value = template.id
  selectedIndicatorIds.value = {}
  indicatorWeights.value = {}
  for (const ind of allIndicators.value) {
    selectedIndicatorIds.value[ind.id] = false
  }
  try {
    const res = await api.get(`/api/v1/tenant/competency/templates/${template.id}/indicators`)
    const existing = res.data?.data || []
    for (const ti of existing) {
      selectedIndicatorIds.value[ti.indicator_id] = true
      indicatorWeights.value[ti.indicator_id] = ti.weight ?? 1
    }
  } catch {
    // template tanpa indicators — tampilkan semua unchecked
  }
  indicatorsVisible.value = true
}

async function handleSaveIndicators() {
  savingIndicators.value = true
  try {
    const payload = allIndicators.value
      .filter(ind => selectedIndicatorIds.value[ind.id])
      .map((ind, idx) => ({
        indicator_id: ind.id,
        weight: Number(indicatorWeights.value[ind.id]) || 0,
        sort_order: idx
      }))
    await api.put(`/api/v1/tenant/competency/templates/${indicatorsTemplateId.value}/indicators`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    indicatorsVisible.value = false
    await loadData()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingIndicators.value = false
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
    await api.delete(`/api/v1/tenant/competency/templates/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
  } catch (e) {
    deleteError.value = getErrorMessage(e, t('message.operation_failed'))
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  loadReferences()
  loadData()
})
</script>
