<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="session">
      <!-- ── Header info session ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ session.session_code }}</h2>
              <Tag :value="statusLabel(session.status)" :severity="statusSeverity(session.status)" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ courseName(session.course_id) }}</p>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <Button :label="t('training.back_to_sessions')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/training/sessions')" />
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mt-4">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.start_date') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.start_date || '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.end_date') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.end_date || '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.provider_type') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.provider_type ? typeLabel(session.provider_type) : '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.delivery_mode') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.delivery_mode ? t(`training.mode_${session.delivery_mode.toLowerCase()}`) : '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.max_quota') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.max_quota }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.trainer') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5 truncate">{{ session.trainer_name || '-' }}</p>
          </div>
        </div>
        <p v-if="session.location || session.meeting_url" class="text-xs text-gray-400 dark:text-gray-500 mt-3">
          <span v-if="session.location"><i class="pi pi-map-marker mr-1"></i>{{ session.location }}</span>
          <span v-if="session.meeting_url" class="ml-3"><i class="pi pi-video mr-1"></i>{{ session.meeting_url }}</span>
        </p>
      </div>

      <!-- ── Tabs ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div class="flex items-center gap-1 px-3 pt-2 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            type="button"
            class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
            :class="activeTab === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
            @click="activeTab = tab.key"
          >
            {{ t(tab.labelKey) }}
          </button>
        </div>

        <!-- ── Overview ── -->
        <div v-if="activeTab === 'overview'" class="p-4">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <!-- Trainers -->
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-user text-emerald-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.trainers') }}</h3>
              </div>
              <div v-if="trainersLoading" class="p-4 space-y-2">
                <div v-for="i in 2" :key="i" class="h-10 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
              </div>
              <div v-else-if="trainers.length" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="st in trainers" :key="st.id" class="flex items-center gap-3 px-3 py-2.5">
                  <div class="w-8 h-8 rounded-full bg-emerald-50 dark:bg-emerald-900/30 flex items-center justify-center shrink-0">
                    <i class="pi pi-user text-emerald-500 text-xs"></i>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ trainerName(st.trainer_id) }}</p>
                    <p class="text-xs text-gray-400">{{ t(`training.role_${(st.role || 'main').toLowerCase()}`) }}</p>
                  </div>
                </div>
              </div>
              <div v-else class="px-3 py-6 text-center text-sm text-gray-400">{{ t('training.trainers_empty') }}</div>
            </div>

            <!-- Assessments -->
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-file-check text-indigo-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.assessments') }}</h3>
              </div>
              <div v-if="assessmentsLoading" class="p-4 space-y-2">
                <div v-for="i in 2" :key="i" class="h-10 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
              </div>
              <div v-else-if="assessments.length" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="a in assessments" :key="a.id" class="flex items-center gap-3 px-3 py-2.5">
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ a.name }}</p>
                    <p class="text-xs text-gray-400">{{ t(`training.assessment_type_${(a.type || 'other').toLowerCase()}`) }} · {{ a.passing_score }}/{{ a.max_score }}</p>
                  </div>
                  <Tag :value="a.is_required ? t('training.required') : t('training.optional')" :severity="a.is_required ? 'danger' : 'secondary'" class="!text-[10px] !px-1.5 !py-0" />
                </div>
              </div>
              <div v-else class="px-3 py-6 text-center text-sm text-gray-400">{{ t('training.assessments_empty') }}</div>
            </div>
          </div>
        </div>

        <!-- ── Participants ── -->
        <div v-if="activeTab === 'participants'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="participants.length" class="text-xs text-gray-400">{{ participants.length }} {{ t('common.items') }}</span>
            <Button :label="t('training.participant_new')" icon="pi pi-user-plus" size="small" @click="openParticipantDialog()" class="ml-auto" />
          </div>
          <DataTable :value="participants" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="participantsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-users text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.participants_empty') }}</p>
              </div>
            </template>
            <Column field="employee_name" :header="t('training.employee')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
            </Column>
            <Column field="registration_status" :header="t('training.registration_status')" style="width:140px">
              <template #body="{data}"><Tag :value="regStatusLabel(data.registration_status)" :severity="regStatusSeverity(data.registration_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="attendance_status" :header="t('training.attendance_status')" style="width:120px">
              <template #body="{data}"><Tag :value="attStatusLabel(data.attendance_status)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="completion_status" :header="t('training.completion_status')" style="width:130px">
              <template #body="{data}"><Tag :value="compStatusLabel(data.completion_status)" :severity="compStatusSeverity(data.completion_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="score" :header="t('training.score')" style="width:90px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.score ?? '-' }}</span></template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeParticipant(data)" />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Materials ── -->
        <div v-if="activeTab === 'materials'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="materials.length" class="text-xs text-gray-400">{{ materials.length }} {{ t('common.items') }}</span>
            <Button :label="t('training.material_new')" icon="pi pi-plus" size="small" @click="openMaterialDialog()" class="ml-auto" />
          </div>
          <DataTable :value="materials" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="materialsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-file text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.materials_empty') }}</p>
              </div>
            </template>
            <Column field="title" :header="t('training.material_title')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.title }}</span></template>
            </Column>
            <Column field="is_required" :header="t('training.is_required')" style="width:110px">
              <template #body="{data}"><Tag :value="data.is_required ? t('common.yes') : t('common.no')" :severity="data.is_required ? 'danger' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="file_url" :header="t('training.material_file')" style="width:200px">
              <template #body="{data}">
                <a v-if="data.file_url" :href="data.file_url" target="_blank" class="text-emerald-600 dark:text-emerald-400 hover:underline text-xs"><i class="pi pi-external-link mr-1"></i>{{ t('common.open') }}</a>
                <span v-else class="text-gray-400">-</span>
              </template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeMaterial(data)" />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- ── Dialog: register participant ── -->
      <Dialog v-model:visible="participantDialogVisible" :header="t('training.participant_new')" modal :style="{ width: '480px' }">
        <div class="space-y-4">
          <FormRow :label="t('training.employee')" required :errors="errors?.employee_id">
            <SelectLabel v-model="participantForm.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.employee_id }" />
          </FormRow>
          <FormRow :label="t('training.registration_status')">
            <SelectLabel v-model="participantForm.registration_status" :options="regStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="participantDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="participantSaving" :disabled="participantSaving" @click="handleAddParticipant" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: add material ── -->
      <Dialog v-model:visible="materialDialogVisible" :header="t('training.material_new')" modal :style="{ width: '520px' }">
        <div class="space-y-4">
          <FormRow :label="t('training.material_title')" required :errors="errors?.title">
            <TextInput v-model="materialForm.title" maxlength="200" :placeholder="t('training.material_title')" :class="{ 'p-invalid': errors?.title }" />
          </FormRow>
          <FormRow :label="t('training.material_description')">
            <TextInput v-model="materialForm.description" textarea :rows="2" />
          </FormRow>
          <FormRow :label="t('training.material_file')">
            <TextInput v-model="materialForm.file_url" maxlength="500" :placeholder="t('training.material_file_placeholder')" />
          </FormRow>
          <FormRow :label="t('training.is_required')">
            <ToggleSwitch v-model="materialForm.is_required" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="materialDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="materialSaving" :disabled="materialSaving" @click="handleAddMaterial" />
          </div>
        </template>
      </Dialog>
    </template>

    <div v-else class="text-center py-12 text-gray-400">
      <i class="pi pi-exclamation-triangle text-3xl mb-2 opacity-50"></i>
      <p class="text-sm">{{ t('training.session_not_found') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import SelectLabel from '@/components/SelectLabel.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()
const route = useRoute()

const sessionId = route.params.id
const session = ref(null)
const loading = ref(true)
const activeTab = ref('overview')

const trainers = ref([])
const trainersLoading = ref(false)
const assessments = ref([])
const assessmentsLoading = ref(false)
const participants = ref([])
const participantsLoading = ref(false)
const materials = ref([])
const materialsLoading = ref(false)

const courses = ref([])
const trainersMaster = ref([])
const employees = ref([])

const participantDialogVisible = ref(false)
const participantSaving = ref(false)
const participantForm = ref({ employee_id: null, registration_status: 'REGISTERED' })
const materialDialogVisible = ref(false)
const materialSaving = ref(false)
const materialForm = ref({ title: '', description: '', file_url: '', is_required: false })
const errors = ref({})

const tabs = [
  { key: 'overview', labelKey: 'training.tab_overview' },
  { key: 'participants', labelKey: 'training.participants' },
  { key: 'materials', labelKey: 'training.materials' }
]

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))

// Status pendaftaran saat create — CANCELLED tidak masuk akal untuk enrollment baru.
const regStatusOptions = computed(() => ['NOMINATED', 'REQUESTED', 'APPROVED', 'REGISTERED', 'WAITLISTED'].map(v => ({ label: regStatusLabel(v), value: v })))

function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || id
}
function trainerName(id) {
  const tr = trainersMaster.value.find(x => x.id === id)
  return tr ? tr.name : id
}
function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || id
}

function typeLabel(type) {
  const key = `training.type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function statusLabel(status) {
  const key = `training.status_${String(status || '').toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
function statusSeverity(status) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'SCHEDULED': return 'info'
    case 'REGISTRATION_OPEN': return 'success'
    case 'FULL': return 'warning'
    case 'IN_PROGRESS': return 'info'
    case 'COMPLETED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function regStatusLabel(s) {
  const key = `training.reg_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function regStatusSeverity(s) {
  switch (s) {
    case 'REGISTERED': return 'success'
    case 'APPROVED': return 'info'
    case 'NOMINATED': return 'warning'
    case 'REQUESTED': return 'info'
    case 'WAITLISTED': return 'warning'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function attStatusLabel(s) {
  const key = `training.att_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusLabel(s) {
  const key = `training.comp_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusSeverity(s) {
  switch (s) {
    case 'COMPLETED': return 'success'
    case 'FAILED': return 'danger'
    case 'IN_PROGRESS': return 'info'
    default: return 'secondary'
  }
}

async function loadSession() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}`)
    session.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadTrainers() {
  trainersLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/trainers`)
    trainers.value = res.data?.data || []
  } catch { trainers.value = [] } finally { trainersLoading.value = false }
}

async function loadAssessments() {
  assessmentsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/assessments`)
    assessments.value = res.data?.data || []
  } catch { assessments.value = [] } finally { assessmentsLoading.value = false }
}

async function loadParticipants() {
  participantsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/participants', { params: { session_id: sessionId, per_page: 200 } })
    participants.value = res.data?.data || []
  } catch { participants.value = [] } finally { participantsLoading.value = false }
}

async function loadMaterials() {
  materialsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/materials', { params: { session_id: sessionId, per_page: 200 } })
    materials.value = res.data?.data || []
  } catch { materials.value = [] } finally { materialsLoading.value = false }
}

async function loadReferences() {
  const [cRes, tRes, eRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/trainers', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  ])
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
  trainersMaster.value = tRes.status === 'fulfilled' ? (tRes.value.data?.data || []) : []
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
}

function openParticipantDialog() {
  errors.value = {}
  participantForm.value = { employee_id: null, registration_status: 'REGISTERED' }
  participantDialogVisible.value = true
}

async function handleAddParticipant() {
  errors.value = {}
  if (!participantForm.value.employee_id) { errors.value = { employee_id: t('form.required') }; return }
  participantSaving.value = true
  try {
    await api.post('/api/v1/tenant/trainings/participants', {
      session_id: sessionId,
      employee_id: participantForm.value.employee_id,
      registration_status: participantForm.value.registration_status
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    participantDialogVisible.value = false
    await loadParticipants()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    participantSaving.value = false
  }
}

async function removeParticipant(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/participants/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadParticipants()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

function openMaterialDialog() {
  errors.value = {}
  materialForm.value = { title: '', description: '', file_url: '', is_required: false }
  materialDialogVisible.value = true
}

async function handleAddMaterial() {
  errors.value = {}
  if (!materialForm.value.title?.trim()) { errors.value = { title: t('form.required') }; return }
  materialSaving.value = true
  try {
    await api.post('/api/v1/tenant/trainings/materials', {
      session_id: sessionId,
      title: materialForm.value.title.trim(),
      description: materialForm.value.description?.trim() || '',
      file_url: materialForm.value.file_url?.trim() || '',
      is_required: materialForm.value.is_required
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    materialDialogVisible.value = false
    await loadMaterials()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    materialSaving.value = false
  }
}

async function removeMaterial(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/materials/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadMaterials()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

onMounted(() => {
  loadSession()
  loadReferences()
  loadTrainers()
  loadAssessments()
  loadParticipants()
  loadMaterials()
})
</script>
