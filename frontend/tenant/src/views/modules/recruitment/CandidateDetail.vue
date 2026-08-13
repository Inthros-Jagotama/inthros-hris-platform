<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-24 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="candidate">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-12 h-12 rounded-full bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-sm font-semibold shrink-0">
              {{ initials }}
            </div>
            <div class="min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100 truncate">{{ candidate.first_name }} {{ candidate.last_name }}</h2>
                <Tag :value="t('candidates.type_' + (candidate.candidate_type || 'external').toLowerCase())" :severity="candidate.candidate_type === 'INTERNAL' ? 'info' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" />
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5 truncate">{{ candidate.email }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <Button :label="t('common.edit')" icon="pi pi-pencil" size="small" severity="secondary" outlined @click="openEditDialog()" />
            <Button :label="t('candidates.back_to_candidates')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/recruitment/candidates')" />
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3 mt-4">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('candidates.candidate_number') }}</p>
            <p class="text-sm font-mono font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ candidate.candidate_number || '—' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('candidates.phone') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ candidate.phone || '—' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('candidates.current_title') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5 truncate">{{ candidate.current_title || '—' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('candidates.current_company') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5 truncate">{{ candidate.current_company || '—' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('candidates.source') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5 truncate">{{ candidate.source || '—' }}</p>
          </div>
        </div>
      </div>

      <!-- Tabs -->
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

        <!-- Educations -->
        <div v-if="activeTab === 'educations'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_education')" icon="pi pi-plus" size="small" @click="openEducationDialog()" />
          </div>
          <div v-if="educations.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in educations" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.institution_name }} <span v-if="item.is_highest" class="ml-1"><Tag :value="t('candidates.is_highest')" severity="success" class="!text-[10px] !px-1.5 !py-0" /></span></p>
                <p class="text-xs text-gray-400">{{ item.major || '—' }} · {{ item.start_year || '?' }}–{{ item.end_year || t('candidates.is_current') }} <span v-if="item.gpa">· GPA {{ item.gpa }}</span></p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmItemDelete('educations', item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.educations_empty') }}</div>
        </div>

        <!-- Work Experiences -->
        <div v-if="activeTab === 'experiences'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_experience')" icon="pi pi-plus" size="small" @click="openExperienceDialog()" />
          </div>
          <div v-if="experiences.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in experiences" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.job_title }} <span class="text-gray-400 font-normal">@ {{ item.company_name }}</span></p>
                <p class="text-xs text-gray-400">{{ item.start_date }} – {{ item.is_current ? t('candidates.is_current') : (item.end_date || '—') }}</p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmItemDelete('experiences', item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.experiences_empty') }}</div>
        </div>

        <!-- Skills -->
        <div v-if="activeTab === 'skills'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_skill')" icon="pi pi-plus" size="small" @click="openSkillDialog()" />
          </div>
          <div v-if="skills.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in skills" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.competency_name || item.competency_id }}</p>
                <p v-if="item.level" class="text-xs text-gray-400">{{ t('candidates.level') }}: {{ item.level }}</p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmItemDelete('skills', item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.skills_empty') }}</div>
        </div>

        <!-- Certifications -->
        <div v-if="activeTab === 'certifications'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_certification')" icon="pi pi-plus" size="small" @click="openCertificationDialog()" />
          </div>
          <div v-if="certifications.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in certifications" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name }}</p>
                <p class="text-xs text-gray-400">{{ item.issuing_organization || '—' }} <span v-if="item.expiry_date">· {{ t('candidates.expiry_date') }}: {{ item.expiry_date }}</span></p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmItemDelete('certifications', item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.certifications_empty') }}</div>
        </div>

        <!-- Documents -->
        <div v-if="activeTab === 'documents'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_document')" icon="pi pi-plus" size="small" @click="openDocumentDialog()" />
          </div>
          <div v-if="documents.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in documents" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name }}</p>
                <p class="text-xs text-gray-400">{{ t('candidates.doc_type_' + (item.document_type || 'other').toLowerCase()) }} · <a :href="item.file_url" target="_blank" rel="noopener" class="text-sky-500 hover:underline">{{ t('candidates.file_url') }}</a></p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmItemDelete('documents', item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.documents_empty') }}</div>
        </div>

        <!-- Consents (append-only) -->
        <div v-if="activeTab === 'consents'" class="p-4">
          <p class="text-xs text-gray-400 mb-3"><i class="pi pi-info-circle mr-1"></i>{{ t('candidates.consent_hint') }}</p>
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('candidates.add_consent')" icon="pi pi-plus" size="small" @click="openConsentDialog()" />
          </div>
          <div v-if="consents.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in consents" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <Tag :value="t('candidates.action_' + (item.action || '').toLowerCase())" :severity="item.action === 'GRANTED' ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
              <div class="min-w-0 flex-1">
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ formatTimestamp(item.changed_at) }}</p>
                <p v-if="item.notes" class="text-xs text-gray-400">{{ item.notes }}</p>
              </div>
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('candidates.consents_empty') }}</div>
        </div>
      </div>
    </template>

    <!-- Edit candidate dialog -->
    <Dialog v-model:visible="editDialogVisible" :header="t('common.edit')" :modal="true" class="!w-[min(95vw,640px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('candidates.first_name')" :required="true"><TextInput v-model="editForm.first_name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.last_name')" :required="true"><TextInput v-model="editForm.last_name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.email')" :required="true"><TextInput v-model="editForm.email" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.phone')"><TextInput v-model="editForm.phone" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.current_company')"><TextInput v-model="editForm.current_company" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.current_title')"><TextInput v-model="editForm.current_title" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.source')"><TextInput v-model="editForm.source" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.linkedin_url')"><TextInput v-model="editForm.linkedin_url" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.resume_url')"><TextInput v-model="editForm.resume_url" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.portfolio_url')"><TextInput v-model="editForm.portfolio_url" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.address')" class="md:col-span-2"><Textarea v-model="editForm.address" :rows="2" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.notes')" class="md:col-span-2"><Textarea v-model="editForm.notes" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="editDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="editSaving" @click="saveCandidate()" />
        </div>
      </template>
    </Dialog>

    <!-- Add education dialog -->
    <Dialog v-model:visible="educationDialogVisible" :header="t('candidates.add_education')" :modal="true" class="!w-[min(95vw,560px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('candidates.institution_name')" :required="true" class="md:col-span-2"><TextInput v-model="educationForm.institution_name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.major')"><TextInput v-model="educationForm.major" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.gpa')"><InputNumber v-model="educationForm.gpa" :minFractionDigits="2" :maxFractionDigits="2" :min="0" :max="4" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.start_year')"><InputNumber v-model="educationForm.start_year" :useGrouping="false" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.end_year')"><InputNumber v-model="educationForm.end_year" :useGrouping="false" class="!w-full" /></FormRow>
        <div class="md:col-span-2 flex items-center gap-2">
          <ToggleSwitch v-model="educationForm.is_highest" />
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ t('candidates.is_highest') }}</span>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="educationDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveEducation()" />
        </div>
      </template>
    </Dialog>

    <!-- Add experience dialog -->
    <Dialog v-model:visible="experienceDialogVisible" :header="t('candidates.add_experience')" :modal="true" class="!w-[min(95vw,560px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('candidates.company_name')" :required="true"><TextInput v-model="experienceForm.company_name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.job_title')" :required="true"><TextInput v-model="experienceForm.job_title" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.employment_type')"><TextInput v-model="experienceForm.employment_type" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.start_date')" :required="true"><DateInput v-model="experienceForm.start_date" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.end_date')"><DateInput v-model="experienceForm.end_date" :disabled="experienceForm.is_current" class="!w-full" /></FormRow>
        <div class="flex items-center gap-2">
          <ToggleSwitch v-model="experienceForm.is_current" />
          <span class="text-sm text-gray-600 dark:text-gray-300">{{ t('candidates.is_current') }}</span>
        </div>
        <FormRow :label="t('common.description')" class="md:col-span-2"><Textarea v-model="experienceForm.description" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="experienceDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveExperience()" />
        </div>
      </template>
    </Dialog>

    <!-- Add skill dialog -->
    <Dialog v-model:visible="skillDialogVisible" :header="t('candidates.add_skill')" :modal="true" class="!w-[min(95vw,480px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('candidates.competency')" :required="true">
          <SelectLabel v-model="skillForm.competency_id" :options="competencyOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('candidates.level')"><InputNumber v-model="skillForm.level" :min="1" :max="5" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.notes')"><Textarea v-model="skillForm.notes" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="skillDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveSkill()" />
        </div>
      </template>
    </Dialog>

    <!-- Add certification dialog -->
    <Dialog v-model:visible="certificationDialogVisible" :header="t('candidates.add_certification')" :modal="true" class="!w-[min(95vw,560px)]">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <FormRow :label="t('candidates.certification_name')" :required="true" class="md:col-span-2"><TextInput v-model="certificationForm.name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.issuing_organization')"><TextInput v-model="certificationForm.issuing_organization" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.credential_url')"><TextInput v-model="certificationForm.credential_url" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.issue_date')"><DateInput v-model="certificationForm.issue_date" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.expiry_date')"><DateInput v-model="certificationForm.expiry_date" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="certificationDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveCertification()" />
        </div>
      </template>
    </Dialog>

    <!-- Add document dialog -->
    <Dialog v-model:visible="documentDialogVisible" :header="t('candidates.add_document')" :modal="true" class="!w-[min(95vw,480px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('candidates.document_type')">
          <SelectLabel v-model="documentForm.document_type" :options="documentTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('candidates.document_name')" :required="true"><TextInput v-model="documentForm.name" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.file_url')" :required="true"><TextInput v-model="documentForm.file_url" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.notes')"><Textarea v-model="documentForm.notes" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="documentDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveDocument()" />
        </div>
      </template>
    </Dialog>

    <!-- Add consent dialog -->
    <Dialog v-model:visible="consentDialogVisible" :header="t('candidates.add_consent')" :modal="true" class="!w-[min(95vw,420px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('candidates.consent_action')" :required="true">
          <SelectLabel v-model="consentForm.action" :options="consentActionOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('candidates.notes')"><Textarea v-model="consentForm.notes" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="consentDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveConsent()" />
        </div>
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="itemDeleteDialogVisible"
      :title="t('common.delete')"
      :message="t('candidates.delete_item_confirm')"
      :loading="itemDeleting"
      @confirm="doItemDelete()"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import ToggleSwitch from 'primevue/toggleswitch'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const candidateId = route.params.id

const loading = ref(true)
const candidate = ref(null)
const activeTab = ref('educations')

const educations = ref([])
const experiences = ref([])
const skills = ref([])
const certifications = ref([])
const documents = ref([])
const consents = ref([])
const competencies = ref([])

const tabs = [
  { key: 'educations', labelKey: 'candidates.tab_educations' },
  { key: 'experiences', labelKey: 'candidates.tab_experiences' },
  { key: 'skills', labelKey: 'candidates.tab_skills' },
  { key: 'certifications', labelKey: 'candidates.tab_certifications' },
  { key: 'documents', labelKey: 'candidates.tab_documents' },
  { key: 'consents', labelKey: 'candidates.tab_consents' }
]

const competencyOptions = computed(() => competencies.value.map(c => ({ label: c.name, value: c.id })))
const documentTypeOptions = computed(() => ['RESUME', 'COVER_LETTER', 'CERTIFICATE', 'PORTFOLIO', 'IDENTITY', 'OTHER'].map(v => ({ label: t('candidates.doc_type_' + v.toLowerCase()), value: v })))
const consentActionOptions = computed(() => ['GRANTED', 'REVOKED'].map(v => ({ label: t('candidates.action_' + v.toLowerCase()), value: v })))

const initials = computed(() => {
  if (!candidate.value) return '?'
  const a = (candidate.value.first_name || '?').charAt(0)
  const b = (candidate.value.last_name || '').charAt(0)
  return (a + b).toUpperCase()
})

function formatTimestamp(value) {
  if (!value) return '—'
  const ms = Number(value) / 1000000
  if (!Number.isFinite(ms) || ms <= 0) return '—'
  return new Date(ms).toLocaleString()
}

// ── Edit candidate ──
const editDialogVisible = ref(false)
const editSaving = ref(false)
const editForm = ref({})

function openEditDialog() {
  editForm.value = { ...candidate.value }
  editDialogVisible.value = true
}

async function saveCandidate() {
  editSaving.value = true
  try {
    const payload = { ...editForm.value }
    delete payload.id
    delete payload.candidate_number
    delete payload.candidate_type
    delete payload.employee_id
    delete payload.created_at
    delete payload.updated_at
    await api.put(`/api/v1/tenant/recruitment/candidates/${candidateId}`, payload)
    editDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.updated'), life: 3000 })
    loadCandidate()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    editSaving.value = false
  }
}

// ── Generic item add/delete plumbing ──
const itemSaving = ref(false)
const itemDeleteDialogVisible = ref(false)
const itemDeleting = ref(false)
const pendingDeleteType = ref(null)
const pendingDeleteItem = ref(null)

const resourcePaths = {
  educations: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/educations`, item: '/api/v1/tenant/recruitment/educations' },
  experiences: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/work-experiences`, item: '/api/v1/tenant/recruitment/work-experiences' },
  skills: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/skills`, item: '/api/v1/tenant/recruitment/skills' },
  certifications: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/certifications`, item: '/api/v1/tenant/recruitment/certifications' },
  documents: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/documents`, item: '/api/v1/tenant/recruitment/documents' },
  consents: { list: `/api/v1/tenant/recruitment/candidates/${candidateId}/consents`, item: null }
}

const listRefs = { educations, experiences, skills, certifications, documents, consents }

function confirmItemDelete(type, item) {
  pendingDeleteType.value = type
  pendingDeleteItem.value = item
  itemDeleteDialogVisible.value = true
}

async function doItemDelete() {
  const type = pendingDeleteType.value
  const item = pendingDeleteItem.value
  if (!type || !item) return
  itemDeleting.value = true
  try {
    await api.delete(`${resourcePaths[type].item}/${item.id}`)
    itemDeleteDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_deleted'), life: 3000 })
    loadSubResource(type)
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemDeleting.value = false
  }
}

async function loadSubResource(type) {
  try {
    const res = await api.get(resourcePaths[type].list)
    listRefs[type].value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

function cleanPayload(payload) {
  const out = { ...payload }
  Object.keys(out).forEach(k => {
    if (out[k] === '' || out[k] === null || out[k] === undefined) delete out[k]
  })
  return out
}

// ── Education ──
const educationDialogVisible = ref(false)
const educationForm = ref({})
function openEducationDialog() {
  educationForm.value = { institution_name: '', major: '', gpa: null, start_year: null, end_year: null, is_highest: false }
  educationDialogVisible.value = true
}
async function saveEducation() {
  if (!educationForm.value.institution_name?.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.institution_name'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(resourcePaths.educations.list, cleanPayload(educationForm.value))
    educationDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('educations')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Work Experience ──
const experienceDialogVisible = ref(false)
const experienceForm = ref({})
function openExperienceDialog() {
  experienceForm.value = { company_name: '', job_title: '', employment_type: '', start_date: '', end_date: '', is_current: false, description: '' }
  experienceDialogVisible.value = true
}
async function saveExperience() {
  if (!experienceForm.value.company_name?.trim() || !experienceForm.value.job_title?.trim() || !experienceForm.value.start_date) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    const payload = cleanPayload(experienceForm.value)
    if (payload.is_current) delete payload.end_date
    await api.post(resourcePaths.experiences.list, payload)
    experienceDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('experiences')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Skills ──
const skillDialogVisible = ref(false)
const skillForm = ref({})
function openSkillDialog() {
  skillForm.value = { competency_id: null, level: null, notes: '' }
  skillDialogVisible.value = true
}
async function saveSkill() {
  if (!skillForm.value.competency_id) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.competency'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(resourcePaths.skills.list, cleanPayload(skillForm.value))
    skillDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('skills')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Certifications ──
const certificationDialogVisible = ref(false)
const certificationForm = ref({})
function openCertificationDialog() {
  certificationForm.value = { name: '', issuing_organization: '', issue_date: '', expiry_date: '', credential_url: '' }
  certificationDialogVisible.value = true
}
async function saveCertification() {
  if (!certificationForm.value.name?.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.certification_name'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(resourcePaths.certifications.list, cleanPayload(certificationForm.value))
    certificationDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('certifications')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Documents ──
const documentDialogVisible = ref(false)
const documentForm = ref({})
function openDocumentDialog() {
  documentForm.value = { document_type: 'OTHER', name: '', file_url: '', notes: '' }
  documentDialogVisible.value = true
}
async function saveDocument() {
  if (!documentForm.value.name?.trim() || !documentForm.value.file_url?.trim()) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(resourcePaths.documents.list, cleanPayload(documentForm.value))
    documentDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('documents')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Consents (append-only) ──
const consentDialogVisible = ref(false)
const consentForm = ref({})
function openConsentDialog() {
  consentForm.value = { action: null, notes: '' }
  consentDialogVisible.value = true
}
async function saveConsent() {
  if (!consentForm.value.action) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('candidates.consent_action'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    await api.post(resourcePaths.consents.list, cleanPayload(consentForm.value))
    consentDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadSubResource('consents')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

async function loadCandidate() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/candidates/${candidateId}`)
    candidate.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}

async function loadCompetencies() {
  try {
    const res = await api.get('/api/v1/tenant/competency/competencies', { params: { per_page: 500 } })
    competencies.value = res.data?.data || []
  } catch {
    // fail-silent — dropdown kosong
  }
}

onMounted(async () => {
  loading.value = true
  await Promise.all([
    loadCandidate(),
    loadCompetencies(),
    loadSubResource('educations'),
    loadSubResource('experiences'),
    loadSubResource('skills'),
    loadSubResource('certifications'),
    loadSubResource('documents'),
    loadSubResource('consents')
  ])
  loading.value = false
})
</script>
