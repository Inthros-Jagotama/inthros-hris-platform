<template>
  <div class="max-w-full mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <Tag v-if="savedEmployeeId" :value="savedEmployeeId" severity="success" class="!text-xs !px-2 !py-1" />
        <span v-if="isEdit && savedEmployeeId" class="text-sm text-gray-500 dark:text-gray-400">{{ savedEmployeeId }}</span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="flex gap-4"><div class="w-56 space-y-2"><div v-for="n in 10" :key="n" class="h-12 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div></div><div class="flex-1 grid grid-cols-1 md:grid-cols-2 gap-4"><div v-for="n in 8" :key="n" class="space-y-2"><div class="h-3 w-20 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div><div class="h-9 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div></div></div></div>
    </div>

    <!-- Two-column layout -->
    <div v-else class="flex gap-6">
      <!-- Left: Navigation Sidebar -->
      <div class="w-56 shrink-0 space-y-1">
        <div v-for="(s, i) in steps" :key="i" role="button" :tabindex="isStepEnabled(i) ? 0 : -1"
          class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none"
          :class="getNavItemClass(i)" @click="selectStep(i)" @keydown.enter="selectStep(i)">
          <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold shrink-0 transition-colors duration-150" :class="getCircleClass(i)">
            <i v-if="stepSaved[i] && i > 0" class="pi pi-check text-xs"></i><span v-else>{{ i + 1 }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="text-sm font-medium truncate" :class="getNavTextClass(i)">{{ t(s.labelKey) }}</div>
            <div v-if="stepSaved[i] && i > 0" class="text-[10px] text-emerald-500 font-medium">{{ t('employee.saved') }}</div>
            <div v-else-if="!isStepEnabled(i) && i > 0" class="text-[10px] text-gray-400">{{ t('employee.personal_first') }}</div>
          </div>
          <i v-if="!isStepEnabled(i) && i > 0" class="pi pi-lock text-xs text-gray-300 dark:text-gray-600 shrink-0"></i>
          <i v-else-if="stepSaved[i] && i > 0" class="pi pi-check-circle text-emerald-400 text-xs shrink-0"></i>
        </div>

      </div>

      <!-- Right: Form Content -->
      <div class="flex-1 min-w-0">
        <PersonalForm v-if="activeStep === 0" :form="form" :errors="stepErrors" :gender-options="genderOptions" :religion-options="religionOptions" :marital-status-options="maritalStatusOptions" :nationality-options="nationalityOptions" :saving="stepLoading" :disabled="isEdit" :employee-id="employeeId" :photo-url="profilePicture" @save="savePersonalData" @update:photo="profilePicture = $event" />

        <AddressForm v-else-if="activeStep === 1" :items="addresses" :errs="addrErrors" :address-type-options="addressTypeOptions" :province-options="provinceOptions" :regency-options="regencyOptions" :district-options="districtOptions" :village-options="villageOptions" :saving="stepLoading" :employee-id="employeeId" :on-search-village="searchVillages" @update:items="addresses = $event" @save="() => saveStep(1)" />

        <ContactForm v-else-if="activeStep === 2" :items="contacts" :errs="contactErrors" :relationship-type-options="relationshipTypeOptions" :saving="stepLoading" :employee-id="employeeId" @update:items="contacts = $event" @save="() => saveStep(2)" />

        <FamilyForm v-else-if="activeStep === 3" :items="families" :errs="famErrors" :relationship-type-options="relationshipTypeOptions" :education-options="educationOptions" :saving="stepLoading" :employee-id="employeeId" @update:items="families = $event" @save="() => saveStep(3)" />

        <EducationForm v-else-if="activeStep === 4" :items="educations" :errs="eduErrors" :education-options="educationOptions" :saving="stepLoading" :employee-id="employeeId" @update:items="educations = $event" @save="() => saveStep(4)" />

        <ExperienceForm v-else-if="activeStep === 5" :items="experiences" :errs="expErrors" :saving="stepLoading" :employee-id="employeeId" @update:items="experiences = $event" @save="() => saveStep(5)" />

        <DocumentForm v-else-if="activeStep === 6" :items="documents" :errs="docErrors" :saving="stepLoading" :employee-id="employeeId" @update:items="documents = $event" @save="() => saveStep(6)" />

        <InsuranceForm v-else-if="activeStep === 7" :items="insurances" :errs="insErrors" :insurance-category-options="insuranceCategoryOptions" :saving="stepLoading" :employee-id="employeeId" @update:items="insurances = $event" @save="() => saveStep(7)" />

        <BankProfileForm v-else-if="activeStep === 8" :items="banks" :errs="bankErrors" :bank-options="bankOptions" :saving="stepLoading" :employee-id="employeeId" @update:items="banks = $event" @save="() => saveStep(8)" />

        <EmploymentForm v-else-if="activeStep === 9" :items="employments" :errs="empErrors" :organization-options="organizationOptions" :employment-status-options="employmentStatusOptions" :saving="stepLoading" @update:items="employments = $event" @save="() => saveStep(9)" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors, getErrorCode } from '@/services/responseHandler'
import api from '@/services/api'

import Tag from 'primevue/tag'

import PersonalForm from './employee/PersonalForm.vue'
import AddressForm from './employee/AddressForm.vue'
import ContactForm from './employee/ContactForm.vue'
import FamilyForm from './employee/FamilyForm.vue'
import EducationForm from './employee/EducationForm.vue'
import ExperienceForm from './employee/ExperienceForm.vue'
import DocumentForm from './employee/DocumentForm.vue'
import InsuranceForm from './employee/InsuranceForm.vue'
import BankProfileForm from './employee/BankProfileForm.vue'
import EmploymentForm from './employee/EmploymentForm.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const employeeId = ref(route.params.id || null)
const activeStep = ref(0)
const stepLoading = ref(false)
const pageLoading = ref(true)
const savedEmployeeId = ref(null)
const profilePicture = ref('')

const stepSaved = reactive(Array(10).fill(false))
const personalDataSaved = computed(() => stepSaved[0])

const steps = [
  { labelKey: 'employee.wizard_step_personal' },
  { labelKey: 'employee.wizard_step_address' },
  { labelKey: 'employee.wizard_step_contact' },
  { labelKey: 'employee.wizard_step_family' },
  { labelKey: 'employee.wizard_step_education' },
  { labelKey: 'employee.wizard_step_experience' },
  { labelKey: 'employee.wizard_step_documents' },
  { labelKey: 'employee.wizard_step_insurance' },
  { labelKey: 'employee.wizard_step_bank' },
  { labelKey: 'employee.wizard_step_employment' }
]

// ── Form Data ──
const form = ref({ employee_id: '', name: '', nik: '', family_id: '', gender: '', mother_name: '', pob: '', dob: '', phone_number: '', email: '', linkedin: '', ig: '', religion_id: '', marital_status_id: '', nationality_type: '', nationality_id: '', passport: '' })
const addresses = ref([])
const contacts = ref([])
const families = ref([])
const educations = ref([])
const experiences = ref([])
const documents = ref([])
const insurances = ref([])
const banks = ref([])
const employments = ref([])

// ── Errors ──
const stepErrors = ref({})
const addrErrors = ref([])
const contactErrors = ref([])
const famErrors = ref([])
const eduErrors = ref([])
const expErrors = ref([])
const docErrors = ref([])
const insErrors = ref([])
const bankErrors = ref([])
const empErrors = ref([])

// ── Nav helpers ──
function isStepEnabled(i) { return i === 0 || !!personalDataSaved.value }
function selectStep(i) {
  if (isStepEnabled(i)) {
    activeStep.value = i
    router.replace({ query: { ...route.query, step: String(i) } })
  }
}

function getNavItemClass(i) {
  const active = activeStep.value === i
  const saved = stepSaved[i]
  if (active) return 'bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700'
  if (!isStepEnabled(i) && i > 0) return 'opacity-50 cursor-not-allowed hover:bg-transparent dark:hover:bg-transparent'
  if (saved && i > 0) return 'hover:bg-gray-50 dark:hover:bg-gray-800'
  return 'hover:bg-gray-50 dark:hover:bg-gray-800'
}

function getCircleClass(i) {
  const active = activeStep.value === i
  const saved = stepSaved[i]
  if (active) return 'bg-emerald-600 text-white'
  if (saved && i > 0) return 'bg-emerald-100 dark:bg-emerald-800 text-emerald-600 dark:text-emerald-300'
  if (!isStepEnabled(i) && i > 0) return 'bg-gray-100 dark:bg-gray-700 text-gray-400 dark:text-gray-500'
  return 'bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300'
}

function getNavTextClass(i) {
  const active = activeStep.value === i
  const saved = stepSaved[i]
  if (active) return 'text-emerald-700 dark:text-emerald-300'
  if (saved && i > 0) return 'text-emerald-600 dark:text-emerald-400'
  if (!isStepEnabled(i) && i > 0) return 'text-gray-400 dark:text-gray-500'
  return 'text-gray-700 dark:text-gray-300'
}

// ── Reference Data Options ──
const genderOptions = computed(() => [
  { label: t('employee.gender_m'), value: 'M' }, { label: t('employee.gender_f'), value: 'F' }
])
const addressTypeOptions = computed(() => [
  { label: t('employee.address_main'), value: 'MAIN' }, { label: t('employee.address_domicile'), value: 'DOMICILE' }
])
const insuranceCategoryOptions = computed(() => [
  { label: 'BPJS Kesehatan', value: 'BPJS Kesehatan' }, { label: 'BPJS Ketenagakerjaan', value: 'BPJS Ketenagakerjaan' }, { label: 'Non BPJS', value: 'Non BPJS' }
])

const religionOptions = ref([])
const maritalStatusOptions = ref([])
const nationalityOptions = ref([])
const relationshipTypeOptions = ref([])
const educationOptions = ref([])
const employmentStatusOptions = ref([])
const bankOptions = ref([])
const provinceOptions = ref([])
const organizationOptions = ref([])
const regencyOptions = ref([])
const districtOptions = ref([])
const villageOptions = ref([])

async function loadRefData() {
  const [relRes, msRes, natRes, rtRes, eduRes, esRes, bankRes, provRes, orgRes] = await Promise.all([
    api.get('/api/v1/tenant/settings/religions?per_page=100'),
    api.get('/api/v1/tenant/settings/marital-statuses?per_page=100'),
    api.get('/api/v1/tenant/settings/nationalities?per_page=250'),
    api.get('/api/v1/tenant/settings/relationship-types?per_page=100'),
    api.get('/api/v1/tenant/settings/educations?per_page=100'),
    api.get('/api/v1/tenant/settings/employment-statuses?per_page=100'),
    api.get('/api/v1/tenant/settings/banks?per_page=200'),
    api.get('/api/v1/tenant/settings/provinces?per_page=100'),
    api.get('/api/v1/tenant/organizations?tree=true')
  ])
  religionOptions.value = (relRes.data?.data || []).map(r => ({ label: r.name, value: r.id }))
  maritalStatusOptions.value = (msRes.data?.data || []).map(m => ({ label: m.name, value: m.id }))
  nationalityOptions.value = (natRes.data?.data || []).map(n => ({ label: `${n.code} - ${n.name}`, value: n.code }))
  relationshipTypeOptions.value = (rtRes.data?.data || []).map(r => ({ label: r.name, value: r.id }))
  educationOptions.value = (eduRes.data?.data || []).map(e => ({ label: e.name, value: e.id }))
  employmentStatusOptions.value = (esRes.data?.data || []).map(e => ({ label: e.name, value: e.id }))
  bankOptions.value = (bankRes.data?.data || []).map(b => ({ label: `${b.name} - ${b.code || ''}`, value: b.id }))
  provinceOptions.value = (provRes.data?.data || []).map(p => ({ label: `${p.code} - ${p.name}`, value: p.code }))
  const orgList = []
  function flattenOrgTree(nodes, depth) {
    if (!nodes) return
    if (Array.isArray(nodes)) { nodes.forEach(n => flattenOrgTree(n, depth)); return }
    orgList.push({ label: '—'.repeat(depth) + ' ' + nodes.nomenclature, value: nodes.id })
    if (nodes.children) flattenOrgTree(nodes.children, depth + 1)
  }
  flattenOrgTree(orgRes.data?.data || [], 0)
  organizationOptions.value = orgList
}

// ── Village AutoComplete search ──
async function searchVillages(query) {
  try {
    const res = await api.get(`/api/v1/tenant/settings/villages/search?q=${encodeURIComponent(query)}`)
    return res.data?.data || []
  } catch { return [] }
}

// ── Load village labels for existing addresses (edit mode) ──
async function loadExistingVillageLabels(addressList) {
  for (const addr of addressList) {
    if (addr.village_id && !addr._villageLabel) {
      try {
        const res = await api.get(`/api/v1/tenant/settings/villages/${addr.village_id}/detail`)
        const found = res.data?.data
        if (found) {
          addr._villageLabel = found.name
          addr._district_name = found.district_name || ''
          addr._regency_name = found.regency_name || ''
          addr._province_name = found.province_name || ''
        }
      } catch { /* ignore */ }
    }
  }
}

// ── Cascading options for addresses ──
const regencyCache = ref({})
const districtCache = ref({})
const villageCache = ref({})

function loadRegencyOptions(provinceId) {
  if (!provinceId || regencyCache.value[provinceId]) return
  api.get(`/api/v1/tenant/settings/provinces/${provinceId}/regencies?per_page=200`).then(res => {
    regencyCache.value[provinceId] = (res.data?.data || []).map(r => ({ label: `${r.code} - ${r.name}`, value: r.code }))
    regencyCache.value = { ...regencyCache.value }
  }).catch(() => {})
}
function loadDistrictOptions(regencyId) {
  if (!regencyId || districtCache.value[regencyId]) return
  api.get(`/api/v1/tenant/settings/regencies/${regencyId}/districts?per_page=200`).then(res => {
    districtCache.value[regencyId] = (res.data?.data || []).map(d => ({ label: `${d.code} - ${d.name}`, value: d.code }))
    districtCache.value = { ...districtCache.value }
  }).catch(() => {})
}
function loadVillageOptions(districtId) {
  if (!districtId || villageCache.value[districtId]) return
  api.get(`/api/v1/tenant/settings/districts/${districtId}/villages?per_page=200`).then(res => {
    villageCache.value[districtId] = (res.data?.data || []).map(v => ({ label: `${v.code} - ${v.name}`, value: v.code }))
    villageCache.value = { ...villageCache.value }
  }).catch(() => {})
}

// Watch address fields and trigger cascading loads
function watchAddresses(items) {
  items.forEach(a => {
    if (a.province_id) loadRegencyOptions(a.province_id)
    if (a.regency_id) loadDistrictOptions(a.regency_id)
    if (a.district_id) loadVillageOptions(a.district_id)
  })
}

// ── Save Personal Data ──
async function savePersonalData() {
  stepErrors.value = {}
  if (!form.value.employee_id?.trim()) { stepErrors.value = { employee_id: [t('form.required')] }; return }
  if (!form.value.name?.trim()) { stepErrors.value = { name: [t('form.required')] }; return }
  stepLoading.value = true
  try {
    const payload = { employee_id: form.value.employee_id, name: form.value.name, nik: form.value.nik || null, family_id: form.value.family_id || null, gender: form.value.gender || null, mother_name: form.value.mother_name || null, nationality_type: form.value.nationality_type || null, nationality_id: form.value.nationality_id || null, passport: form.value.passport || null, pob: form.value.pob || null, dob: form.value.dob || null, phone_number: form.value.phone_number || null, email: form.value.email || null, linkedin: form.value.linkedin || null, ig: form.value.ig || null, religion_id: form.value.religion_id || null, marital_status_id: form.value.marital_status_id || null }
    if (isEdit.value && employeeId.value) {
      await api.put(`/api/v1/tenant/employees/${employeeId.value}`, payload)
    } else {
      const res = await api.post('/api/v1/tenant/employees', payload)
      employeeId.value = res.data?.data?.id
      router.replace(`/employees/${employeeId.value}/edit`)
    }
    savedEmployeeId.value = form.value.employee_id
    stepSaved[0] = true
    toast.add({ severity: 'success', summary: t('message.success'), detail: isEdit.value ? t('employee.updated') : t('employee.created'), life: 2000 })
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { stepErrors.value = fe }
    else if (getErrorCode(e) === 'QUOTA_EXCEEDED') {
      toast.add({ severity: 'error', summary: t('message.error'), detail: t('employee.quota_exceeded'), life: 5000 })
    }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { stepLoading.value = false }
}

// ── Save Step (sub-records) — sequential per-item with per-index error mapping ──
async function saveStep(stepIndex) {
  if (!employeeId.value) { toast.add({ severity: 'warn', summary: t('employee.personal_first'), life: 3000 }); return }
  stepLoading.value = true
  const empId = employeeId.value

  // Map stepIndex to error refs — clear errors for this step      const errorRefs = [null, addrErrors, contactErrors, famErrors, eduErrors, expErrors, docErrors, insErrors, bankErrors, empErrors]
      if (errorRefs[stepIndex]) errorRefs[stepIndex].value = []

  try {
    const stepConfigs = {
      1: { items: addresses.value, ep: 'addresses', build: a => (!a.type || !a.address) ? null : { type: a.type, address: a.address, province_id: a.province_id || null, regency_id: a.regency_id || null, district_id: a.district_id || null, village_id: a.village_id || null, postal_code: a.postal_code || null } },
      2: { items: contacts.value, ep: 'emergency-contacts', build: c => (!c.name || !c.phone_number) ? null : { name: c.name, phone_number: c.phone_number, relationship_type_id: c.relationship_type_id || null, address: c.address || null } },
      3: { items: families.value, ep: 'families', build: f => (!f.name) ? null : { nik: f.nik || null, name: f.name, dob: f.dob || null, relationship_type_id: f.relationship_type_id || null, education_id: f.education_id || null } },
      4: { items: educations.value, ep: 'educations', build: e => (!e.name) ? null : { education_id: e.education_id || null, name: e.name, major: e.major || null, graduation_year: e.graduation_year ? parseInt(e.graduation_year) : null } },
      5: { items: experiences.value, ep: 'experiences', build: x => (!x.company) ? null : { company: x.company, position: x.position || null, start_year: x.start_year ? parseInt(x.start_year) : null, end_year: x.end_year ? parseInt(x.end_year) : null } },
      6: { items: documents.value, ep: 'documents', build: d => (!d.name || !d.file) ? null : { name: d.name, file: d.file, note: d.note || null } },
      7: { items: insurances.value, ep: 'insurances', build: i => (!i.number || !i.name) ? null : { category: i.category || null, number: i.number, name: i.name, type: i.type || null } },
      8: { items: banks.value, ep: 'banks', build: b => (!b.account_number || !b.account_name) ? null : { bank_id: b.bank_id || null, account_number: b.account_number, account_name: b.account_name } },
      9: { items: employments.value, ep: 'employments', build: em => (!em.decision_letter_number || !em.decision_letter_date || !em.effective_date) ? null : { organization_id: em.organization_id || null, employment_status_id: em.employment_status_id || null, decision_letter_number: em.decision_letter_number, decision_letter_date: em.decision_letter_date, effective_date: em.effective_date, effective_end_date: em.effective_end_date || null } }
    }
    const cfg = stepConfigs[stepIndex]
    if (!cfg) { stepLoading.value = false; return }

    let savedAny = false
    for (let idx = 0; idx < cfg.items.length; idx++) {
      const item = cfg.items[idx]
      if (item._saved) continue
      const payload = cfg.build(item)
      if (!payload) continue

      try {
        await api.post(`/api/v1/tenant/employees/${empId}/${cfg.ep}`, payload)
        item._saved = true
        savedAny = true
      } catch (e) {
        const fe = getValidationErrors(e)
        if (Object.keys(fe).length > 0 && errorRefs[stepIndex]) {
          const arr = [...errorRefs[stepIndex].value]
          arr[idx] = fe
          errorRefs[stepIndex].value = arr
        } else {
          throw e // re-throw non-validation errors to outer catch
        }
      }
    }

    if (savedAny) {
      stepSaved[stepIndex] = true
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.saved'), life: 2000 })
    } else {
      const hasErrors = errorRefs[stepIndex] && errorRefs[stepIndex].value.some(e => e && Object.keys(e).length > 0)
      if (!hasErrors) {
        toast.add({ severity: 'info', summary: t('employee.no_new_items'), life: 2000 })
      }
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally { stepLoading.value = false }
}

async function loadEmployee() {
  if (!isEdit.value || !employeeId.value) return
  const res = await api.get(`/api/v1/tenant/employees/${employeeId.value}`)
  const data = res.data?.data
  if (!data) return
  form.value = { employee_id: data.employee_id || '', name: data.name || '', nik: data.nik || '', family_id: data.family_id || '', gender: data.gender || '', mother_name: data.mother_name || '', nationality_type: data.nationality_type || '', nationality_id: data.nationality_id || '', passport: data.passport || '', pob: data.pob || '', dob: data.dob || '', phone_number: data.phone_number || '', email: data.email || '', linkedin: data.linkedin || '', ig: data.ig || '', religion_id: data.religion_id || '', marital_status_id: data.marital_status_id || '' }
  savedEmployeeId.value = data.employee_id
  profilePicture.value = data.profile_picture || ''
  stepSaved[0] = true
  const markSaved = item => { item._saved = true }
  if (data.addresses) {
    addresses.value = data.addresses.map(a => {
      const item = { type: a.type || '', address: a.address || '', province_id: a.province_id || '', regency_id: a.regency_id || '', district_id: a.district_id || '', village_id: a.village_id || '', postal_code: a.postal_code || '', _villageLabel: a.village_label || '', _id: a.id || '' }
      markSaved(item)
      return item
    })
    if (addresses.value.length > 0) stepSaved[1] = true
    watchAddresses(addresses.value)
    // Load village labels for existing villages
    loadExistingVillageLabels(addresses.value)
  }
  if (data.emergency_contacts) { contacts.value = data.emergency_contacts.map(c => { const item = { name: c.name || '', relationship_type_id: c.relationship_type_id || '', phone_number: c.phone_number || '', address: c.address || '', _id: c.id || '' }; markSaved(item); return item }); if (contacts.value.length > 0) stepSaved[2] = true }
  if (data.families) { families.value = data.families.map(f => { const item = { nik: f.nik || '', name: f.name || '', dob: f.dob || '', relationship_type_id: f.relationship_type_id || '', education_id: f.education_id || '', _id: f.id || '' }; markSaved(item); return item }); if (families.value.length > 0) stepSaved[3] = true }
  if (data.educations) { educations.value = data.educations.map(e => { const item = { education_id: e.education_id || '', name: e.name || '', major: e.major || '', graduation_year: e.graduation_year ? String(e.graduation_year) : '', _id: e.id || '' }; markSaved(item); return item }); if (educations.value.length > 0) stepSaved[4] = true }
  if (data.experiences) { experiences.value = data.experiences.map(x => { const item = { company: x.company || '', position: x.position || '', start_year: x.start_year ? String(x.start_year) : '', end_year: x.end_year ? String(x.end_year) : '', _id: x.id || '' }; markSaved(item); return item }); if (experiences.value.length > 0) stepSaved[5] = true }
  if (data.documents) { documents.value = data.documents.map(d => { const item = { name: d.name || '', file: d.file || '', note: d.note || '', _id: d.id || '' }; markSaved(item); return item }); if (documents.value.length > 0) stepSaved[6] = true }
  if (data.insurances) { insurances.value = data.insurances.map(i => { const item = { category: i.category || '', number: i.number || '', name: i.name || '', type: i.type || '', _id: i.id || '' }; markSaved(item); return item }); if (insurances.value.length > 0) stepSaved[7] = true }
  if (data.banks) { banks.value = data.banks.map(b => { const item = { bank_id: b.bank_id || '', account_number: b.account_number || '', account_name: b.account_name || '', _id: b.id || '' }; markSaved(item); return item }); if (banks.value.length > 0) stepSaved[8] = true }
  if (data.employments) { employments.value = data.employments.map(em => { const item = { organization_id: em.organization_id || '', employment_status_id: em.employment_status_id || '', decision_letter_number: em.decision_letter_number || '', decision_letter_date: em.decision_letter_date || '', effective_date: em.effective_date || '', effective_end_date: em.effective_end_date || '' }; markSaved(item); return item }); if (employments.value.length > 0) stepSaved[9] = true }
}

onMounted(async () => {
  try {
    await loadRefData()
    if (isEdit.value) await loadEmployee()

    // Restore active step from query param (only if step is enabled/personal data saved)
    const stepParam = parseInt(route.query.step)
    if (!isNaN(stepParam) && stepParam >= 0 && stepParam < steps.length) {
      if (isStepEnabled(stepParam)) {
        activeStep.value = stepParam
      }
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 4000 })
  } finally { pageLoading.value = false }
})
</script>
