<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <span v-if="flowsTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">
        {{ flowsTotal }} {{ t('common.items') }}
      </span>
      <Button :label="t('approval.new_flow')" icon="pi pi-plus" size="small" @click="openFlowDialog()" />
    </div>

    <SkeletonTable v-if="flowsLoading" :columns="flowSkeletonColumns" :rows="6" />

    <DataTable v-else :value="flows" size="small" class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-check-square text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('approval.empty_flows') }}</p>
        </div>
      </template>
      <Column field="module" :header="t('approval.module')" style="width:160px">
        <template #body="{data}">
          <Tag :value="data.module" severity="info" class="!text-xs !px-1.5 !py-0.5" />
        </template>
      </Column>
      <Column field="name" :header="t('approval.flow_name')">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.name }}</span>
        </template>
      </Column>
      <Column field="version" :header="t('approval.version')" style="width:90px">
        <template #body="{data}">v{{ data.version }}</template>
      </Column>
      <Column field="steps" :header="t('approval.steps')" style="width:90px">
        <template #body="{data}">{{ data.steps?.length || 0 }}</template>
      </Column>
      <Column field="is_active" :header="t('common.status')" style="width:100px">
        <template #body="{data}">
          <Tag :value="data.is_active ? t('common_status.active') : t('common_status.inactive')" :severity="data.is_active ? 'success' : 'secondary'" class="!text-xs" />
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:140px" frozen alignFrozen="right">
        <template #body="{data}">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-sitemap" size="small" text severity="secondary" v-tooltip.left="t('approval.manage_steps')" @click="openStepsDialog(data)" />
            <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openFlowDialog(data)" />
            <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteFlow(data)" />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- ===================================================================== -->
    <!-- Flow create/edit dialog -->
    <!-- ===================================================================== -->
    <Dialog v-model:visible="flowDialogVisible" :header="flowEditing ? t('approval.edit_flow') : t('approval.new_flow')" modal :style="{ width: '480px' }" @hide="resetFlowForm">
      <div class="space-y-4">
        <FormRow :label="t('approval.module')" required :errors="flowErrors?.module?.[0]">
          <Select
            v-model="flowForm.module"
            :options="availableModules"
            :placeholder="t('approval.select_module')"
            class="w-full"
            :disabled="flowEditing"
            :class="{'p-invalid': flowErrors?.module}"
          />
        </FormRow>
        <FormRow :label="t('approval.flow_name')" required :errors="flowErrors?.name?.[0]">
          <TextInput v-model="flowForm.name" maxlength="150" :class="{'p-invalid': flowErrors?.name}" />
        </FormRow>
        <FormRow :label="t('approval.version')">
          <InputNumber v-model="flowForm.version" class="!w-full" :min="1" size="small" />
        </FormRow>
        <FormRow v-if="flowEditing" :label="t('common.status')">
          <div class="flex items-center h-full pt-1">
            <ToggleSwitch v-model="flowForm.is_active" />
            <span class="ml-2 text-sm text-gray-600 dark:text-gray-300">{{ flowForm.is_active ? t('common_status.active') : t('common_status.inactive') }}</span>
          </div>
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="flowDialogVisible=false" />
          <Button :label="flowEditing ? t('common.update') : t('common.save')" size="small" :loading="flowSaving" :disabled="flowSaving" @click="handleSaveFlow" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="flowDeleteDialogVisible"
      :title="t('approval.delete_flow_title')"
      :message="t('approval.delete_flow_message', { name: flowDeleteTarget?.name })"
      :loading="flowDeleting"
      :error="flowDeleteError"
      @confirm="handleDeleteFlow"
    />

    <!-- ===================================================================== -->
    <!-- Steps management dialog -->
    <!-- ===================================================================== -->
    <Dialog v-model:visible="stepsDialogVisible" :header="t('approval.manage_steps') + (activeFlow ? ' — ' + activeFlow.name : '')" modal :style="{ width: '820px' }" :closable="true">
      <div class="space-y-3">
        <div class="flex items-center justify-end">
          <Button :label="t('approval.new_step')" icon="pi pi-plus" size="small" @click="openStepDialog()" />
        </div>

        <div v-if="stepsLoading" class="py-8 text-center text-gray-400 text-sm">{{ t('common.loading') }}</div>

        <div v-else-if="steps.length === 0" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-lg">
          <i class="pi pi-sitemap text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('approval.empty_steps') }}</p>
        </div>

        <div v-else class="space-y-2">
          <div v-for="(step, idx) in steps" :key="step.id" class="flex items-center gap-3 border border-gray-200 dark:border-gray-700 rounded-lg p-3">
            <div class="w-8 h-8 rounded-full bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-xs font-semibold shrink-0">
              {{ step.step_order }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ step.step_name }}</span>
                <Tag :value="step.approver_type" severity="info" class="!text-xs !px-1.5 !py-0.5" />
                <Tag :value="step.participation_type" :severity="step.participation_type === 'WATCHER' ? 'secondary' : 'success'" class="!text-xs !px-1.5 !py-0.5" />
                <Tag :value="step.approval_mode" severity="secondary" class="!text-xs !px-1.5 !py-0.5" />
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                <span v-if="step.hierarchy_level">{{ t('approval.hierarchy_level') }}: {{ step.hierarchy_level }}</span>
                <span v-if="step.organization_ids?.length"> · {{ step.organization_ids.length }} {{ t('approval.organizations') }}</span>
                <span v-if="step.required_approvals"> · {{ t('approval.required_approvals') }}: {{ step.required_approvals }}</span>
              </p>
            </div>
            <div class="flex items-center gap-1 shrink-0">
              <Button icon="pi pi-arrow-up" size="small" text severity="secondary" :disabled="idx === 0" @click="moveStep(idx, -1)" />
              <Button icon="pi pi-arrow-down" size="small" text severity="secondary" :disabled="idx === steps.length - 1" @click="moveStep(idx, 1)" />
              <Button icon="pi pi-pencil" size="small" text severity="secondary" @click="openStepDialog(step)" />
              <Button icon="pi pi-trash" size="small" text severity="danger" @click="confirmDeleteStep(step)" />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="stepsDialogVisible=false" />
      </template>
    </Dialog>

    <!-- Step create/edit dialog -->
    <Dialog v-model:visible="stepDialogVisible" :header="stepEditing ? t('approval.edit_step') : t('approval.new_step')" modal :style="{ width: '520px' }" @hide="resetStepForm">
      <div class="space-y-4">
        <FormRow :label="t('approval.step_name')" required :errors="stepErrors?.step_name?.[0]">
          <TextInput v-model="stepForm.step_name" maxlength="150" :class="{'p-invalid': stepErrors?.step_name}" />
        </FormRow>

        <FormRow :label="t('approval.approver_type')" required>
          <Select v-model="stepForm.approver_type" :options="approverTypeOptions" optionLabel="label" optionValue="value" class="w-full" />
        </FormRow>

        <FormRow v-if="['SUPERVISOR','BOTH'].includes(stepForm.approver_type)" :label="t('approval.hierarchy_level')">
          <InputNumber v-model="stepForm.hierarchy_level" class="!w-full" :min="1" :max="10" size="small" />
          <p class="text-xs text-gray-400 mt-1">{{ t('approval.hierarchy_level_hint') }}</p>
        </FormRow>

        <FormRow v-if="['ORGANIZATION','BOTH'].includes(stepForm.approver_type)" :label="t('approval.organizations')">
          <MultiSelect v-model="stepForm.organization_ids" :options="organizationOptions" optionLabel="label" optionValue="value" filter display="chip" class="w-full" :placeholder="t('approval.select_organizations')" />
        </FormRow>

        <FormRow v-if="stepForm.approver_type === 'ROLE'" :label="t('approval.role')">
          <Select v-model="stepForm.role_id" :options="roleOptions" optionLabel="label" optionValue="value" filter showClear class="w-full" :placeholder="t('approval.select_role')" />
        </FormRow>

        <FormRow v-if="stepForm.approver_type === 'USER'" :label="t('approval.user_id')">
          <TextInput v-model="stepForm.approver_user_id" :placeholder="t('approval.user_id_placeholder')" />
        </FormRow>

        <FormRow :label="t('approval.participation_type')">
          <Select v-model="stepForm.participation_type" :options="participationTypeOptions" optionLabel="label" optionValue="value" class="w-full" />
          <p class="text-xs text-gray-400 mt-1">{{ t('approval.participation_type_hint') }}</p>
        </FormRow>

        <template v-if="stepForm.participation_type === 'APPROVER'">
          <FormRow :label="t('approval.approval_mode')">
            <Select v-model="stepForm.approval_mode" :options="approvalModeOptions" optionLabel="label" optionValue="value" class="w-full" />
            <p class="text-xs text-gray-400 mt-1">{{ t(`approval.approval_mode_hint_${stepForm.approval_mode}`) }}</p>
          </FormRow>
          <FormRow v-if="stepForm.approval_mode === 'N_OF_M'" :label="t('approval.required_approvals')">
            <InputNumber v-model="stepForm.required_approvals" class="!w-full" :min="1" size="small" />
          </FormRow>
          <FormRow :label="t('approval.allow_reject')">
            <ToggleSwitch v-model="stepForm.allow_reject" />
          </FormRow>
        </template>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="stepDialogVisible=false" />
          <Button :label="stepEditing ? t('common.update') : t('common.save')" size="small" :loading="stepSaving" :disabled="stepSaving" @click="handleSaveStep" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="stepDeleteDialogVisible"
      :title="t('approval.delete_step_title')"
      :message="t('approval.delete_step_message', { name: stepDeleteTarget?.step_name })"
      :loading="stepDeleting"
      :error="stepDeleteError"
      @confirm="handleDeleteStep"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Select from 'primevue/select'
import MultiSelect from 'primevue/multiselect'
import Dialog from 'primevue/dialog'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const { t } = useI18n()
const toast = useToast()

// =========================================================================
// Flows
// =========================================================================

const flows = ref([])
const flowsLoading = ref(false)
const flowsTotal = ref(0)
const availableModules = ref([])

const flowSkeletonColumns = [
  { type: 'tag', width: 'w-24', headerWidth: 'w-16' },
  { type: 'text', width: 'w-40', headerWidth: 'w-24' },
  { type: 'text', width: 'w-10', headerWidth: 'w-12' },
  { type: 'text', width: 'w-10', headerWidth: 'w-12' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'icons', count: 3, headerWidth: 'w-16' }
]

async function loadFlows() {
  flowsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/approval/flows', { params: { page: 1, per_page: 100 } })
    const body = res.data
    flows.value = body?.data || []
    flowsTotal.value = body?.total || 0
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    flowsLoading.value = false
  }
}

async function loadAvailableModules() {
  try {
    const res = await api.get('/api/v1/tenant/approval/available-modules')
    availableModules.value = res.data?.data || []
  } catch {
    availableModules.value = []
  }
}

const flowDialogVisible = ref(false)
const flowEditing = ref(false)
const flowEditingId = ref(null)
const flowSaving = ref(false)
const flowErrors = ref({})
const flowForm = ref({ module: '', name: '', version: 1, is_active: true })

function openFlowDialog(item) {
  flowEditing.value = !!item
  flowEditingId.value = item?.id || null
  flowErrors.value = {}
  flowForm.value = {
    module: item?.module || '',
    name: item?.name || '',
    version: item?.version || 1,
    is_active: item?.is_active ?? true
  }
  flowDialogVisible.value = true
}

function resetFlowForm() {
  flowForm.value = { module: '', name: '', version: 1, is_active: true }
  flowErrors.value = {}
  flowEditing.value = false
  flowEditingId.value = null
}

async function handleSaveFlow() {
  flowErrors.value = {}
  if (!flowForm.value.module) {
    flowErrors.value = { module: [t('form.required')] }
    return
  }
  if (!flowForm.value.name?.trim()) {
    flowErrors.value = { name: [t('form.required')] }
    return
  }

  flowSaving.value = true
  try {
    if (flowEditing.value) {
      await api.put(`/api/v1/tenant/approval/flows/${flowEditingId.value}`, {
        name: flowForm.value.name,
        version: flowForm.value.version,
        is_active: flowForm.value.is_active
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.flow_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/approval/flows', {
        module: flowForm.value.module,
        name: flowForm.value.name,
        version: flowForm.value.version
      })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.flow_created'), life: 3000 })
    }
    flowDialogVisible.value = false
    await loadFlows()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      flowErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    flowSaving.value = false
  }
}

const flowDeleteDialogVisible = ref(false)
const flowDeleting = ref(false)
const flowDeleteError = ref('')
const flowDeleteTarget = ref(null)

function confirmDeleteFlow(item) {
  flowDeleteTarget.value = item
  flowDeleteError.value = ''
  flowDeleteDialogVisible.value = true
}

async function handleDeleteFlow() {
  flowDeleting.value = true
  flowDeleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/approval/flows/${flowDeleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.flow_deleted'), life: 3000 })
    flowDeleteDialogVisible.value = false
    await loadFlows()
  } catch (e) {
    flowDeleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    flowDeleting.value = false
  }
}

// =========================================================================
// Steps
// =========================================================================

const stepsDialogVisible = ref(false)
const stepsLoading = ref(false)
const activeFlow = ref(null)
const steps = ref([])

const organizationOptions = ref([])
const roleOptions = ref([])

async function loadOrganizationOptions() {
  try {
    const res = await api.get('/api/v1/tenant/organizations', { params: { per_page: 500, active_only: true } })
    const list = res.data?.data || []
    organizationOptions.value = list.map(o => ({ label: o.name, value: o.id }))
  } catch {
    organizationOptions.value = []
  }
}

async function loadRoleOptions() {
  try {
    const res = await api.get('/api/v1/tenant/rbac/roles')
    const list = res.data?.data || []
    roleOptions.value = list.map(r => ({ label: r.name, value: r.id }))
  } catch {
    roleOptions.value = []
  }
}

async function openStepsDialog(flow) {
  activeFlow.value = flow
  stepsDialogVisible.value = true
  await loadSteps()
}

async function loadSteps() {
  if (!activeFlow.value) return
  stepsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps`)
    steps.value = (res.data?.data || []).sort((a, b) => a.step_order - b.step_order)
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    stepsLoading.value = false
  }
}

const approverTypeOptions = [
  { label: 'SUPERVISOR', value: 'SUPERVISOR' },
  { label: 'ORGANIZATION', value: 'ORGANIZATION' },
  { label: 'BOTH', value: 'BOTH' },
  { label: 'ROLE', value: 'ROLE' },
  { label: 'USER', value: 'USER' }
]
const participationTypeOptions = [
  { label: 'APPROVER', value: 'APPROVER' },
  { label: 'WATCHER', value: 'WATCHER' }
]
const approvalModeOptions = [
  { label: 'ANY_ONE', value: 'ANY_ONE' },
  { label: 'ALL', value: 'ALL' },
  { label: 'N_OF_M', value: 'N_OF_M' }
]

const stepDialogVisible = ref(false)
const stepEditing = ref(false)
const stepEditingId = ref(null)
const stepSaving = ref(false)
const stepErrors = ref({})
const stepForm = ref(defaultStepForm())

function defaultStepForm() {
  return {
    step_name: '',
    approver_type: 'SUPERVISOR',
    hierarchy_level: 1,
    organization_ids: [],
    role_id: null,
    approver_user_id: '',
    participation_type: 'APPROVER',
    approval_mode: 'ANY_ONE',
    required_approvals: 1,
    allow_reject: true
  }
}

function openStepDialog(step) {
  stepEditing.value = !!step
  stepEditingId.value = step?.id || null
  stepErrors.value = {}
  if (step) {
    stepForm.value = {
      step_name: step.step_name,
      approver_type: step.approver_type,
      hierarchy_level: step.hierarchy_level || 1,
      organization_ids: step.organization_ids || [],
      role_id: step.role_id || null,
      approver_user_id: step.approver_user_id || '',
      participation_type: step.participation_type,
      approval_mode: step.approval_mode,
      required_approvals: step.required_approvals || 1,
      allow_reject: step.allow_reject
    }
  } else {
    stepForm.value = defaultStepForm()
  }
  if (organizationOptions.value.length === 0) loadOrganizationOptions()
  if (roleOptions.value.length === 0) loadRoleOptions()
  stepDialogVisible.value = true
}

function resetStepForm() {
  stepForm.value = defaultStepForm()
  stepErrors.value = {}
  stepEditing.value = false
  stepEditingId.value = null
}

async function handleSaveStep() {
  stepErrors.value = {}
  if (!stepForm.value.step_name?.trim()) {
    stepErrors.value = { step_name: [t('form.required')] }
    return
  }

  const payload = {
    step_name: stepForm.value.step_name,
    approver_type: stepForm.value.approver_type,
    participation_type: stepForm.value.participation_type
  }
  if (['SUPERVISOR', 'BOTH'].includes(stepForm.value.approver_type)) {
    payload.hierarchy_level = stepForm.value.hierarchy_level || 1
  }
  if (['ORGANIZATION', 'BOTH'].includes(stepForm.value.approver_type)) {
    payload.organization_ids = stepForm.value.organization_ids
  }
  if (stepForm.value.approver_type === 'ROLE') {
    payload.role_id = stepForm.value.role_id
  }
  if (stepForm.value.approver_type === 'USER') {
    payload.approver_user_id = stepForm.value.approver_user_id
  }
  if (stepForm.value.participation_type === 'APPROVER') {
    payload.approval_mode = stepForm.value.approval_mode
    payload.allow_reject = stepForm.value.allow_reject
    if (stepForm.value.approval_mode === 'N_OF_M') {
      payload.required_approvals = stepForm.value.required_approvals
    }
  }

  stepSaving.value = true
  try {
    if (stepEditing.value) {
      await api.put(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps/${stepEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.step_updated'), life: 3000 })
    } else {
      await api.post(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.step_created'), life: 3000 })
    }
    stepDialogVisible.value = false
    await loadSteps()
    await loadFlows()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      stepErrors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
    }
  } finally {
    stepSaving.value = false
  }
}

const stepDeleteDialogVisible = ref(false)
const stepDeleting = ref(false)
const stepDeleteError = ref('')
const stepDeleteTarget = ref(null)

function confirmDeleteStep(step) {
  stepDeleteTarget.value = step
  stepDeleteError.value = ''
  stepDeleteDialogVisible.value = true
}

async function handleDeleteStep() {
  stepDeleting.value = true
  stepDeleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps/${stepDeleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.step_deleted'), life: 3000 })
    stepDeleteDialogVisible.value = false
    await loadSteps()
    await loadFlows()
  } catch (e) {
    stepDeleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    stepDeleting.value = false
  }
}

async function moveStep(idx, direction) {
  const other = idx + direction
  if (other < 0 || other >= steps.value.length) return
  const a = steps.value[idx]
  const b = steps.value[other]
  try {
    await Promise.all([
      api.put(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps/${a.id}`, { step_order: b.step_order }),
      api.put(`/api/v1/tenant/approval/flows/${activeFlow.value.id}/steps/${b.id}`, { step_order: a.step_order })
    ])
    await loadSteps()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  }
}

onMounted(() => {
  loadFlows()
  loadAvailableModules()
})
</script>
