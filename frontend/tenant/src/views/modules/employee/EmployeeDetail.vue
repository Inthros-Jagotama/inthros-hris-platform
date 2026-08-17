<template>
  <div class="max-w-full mx-auto">
    <!-- Header -->
    <div class="flex items-start justify-between gap-4 flex-wrap mb-6">
      <div class="flex items-center gap-4 min-w-0">
        <div class="w-16 h-16 rounded-full overflow-hidden shrink-0 bg-gray-100 dark:bg-gray-800 flex items-center justify-center border border-gray-200 dark:border-gray-700">
          <img v-if="emp.profile_picture" :src="emp.profile_picture" alt="" class="w-full h-full object-cover" />
          <i v-else class="pi pi-user text-2xl text-gray-400"></i>
        </div>
        <div class="min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <h1 class="text-lg font-semibold text-gray-800 dark:text-gray-100 truncate">{{ emp.name || '-' }}</h1>
            <Tag v-if="emp.employee_id" :value="emp.employee_id" severity="info" class="!text-xs !px-1.5 !py-0.5" />
            <Tag v-if="emp.status" :value="t('common_status.' + emp.status)" :severity="emp.status === 'active' ? 'success' : emp.status === 'inactive' ? 'warn' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
            <Tag v-if="emp.recruited_from_application_id" :value="t('employee.from_offer')" severity="success" class="!text-xs !px-1.5 !py-0.5" />
          </div>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 flex items-center gap-2 flex-wrap">
            <span v-if="emp.nik" class="font-mono">{{ emp.nik }}</span>
            <span v-if="emp.email" class="flex items-center gap-1"><i class="pi pi-envelope text-xs"></i>{{ emp.email }}</span>
            <span v-if="emp.phone_number" class="flex items-center gap-1"><i class="pi pi-phone text-xs"></i>{{ emp.phone_number }}</span>
          </p>
        </div>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        <Button :label="t('common.back')" icon="pi pi-arrow-left" severity="secondary" outlined size="small" @click="router.push('/employees')" />
        <Button v-if="canEdit" :label="t('employee.edit')" icon="pi pi-pencil" size="small" @click="router.push(`/employees/${emp.id}/edit`)" />
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <div class="flex gap-6">
        <div class="w-56 space-y-2"><div v-for="n in 10" :key="n" class="h-10 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div></div>
        <div class="flex-1 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div v-for="n in 8" :key="n" class="space-y-2">
              <div class="h-3 w-24 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
              <div class="h-6 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Two-column layout: sub-menu + content -->
    <div v-else class="flex gap-6">
      <!-- Left: Sub-menu navigasi (sama seperti wizard form) -->
      <div class="w-56 shrink-0 space-y-1">
        <div v-for="(s, i) in sections" :key="i" role="button" :tabindex="0"
             class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150 cursor-pointer select-none"
             :class="active === i ? 'bg-emerald-50 dark:bg-emerald-900/20 ring-1 ring-emerald-300 dark:ring-emerald-700' : 'hover:bg-gray-50 dark:hover:bg-gray-800'"
             @click="selectSection(i)" @keydown.enter="selectSection(i)">
          <i :class="[s.icon, 'text-xs shrink-0', active === i ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 dark:text-gray-500']"></i>
          <span class="flex-1 min-w-0 truncate" :class="active === i ? 'text-emerald-700 dark:text-emerald-300 font-medium' : 'text-gray-700 dark:text-gray-300'">{{ t(s.labelKey) }}</span>
          <span v-if="s.count" class="text-[10px] font-semibold px-1.5 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400">{{ s.count }}</span>
        </div>
      </div>

      <!-- Right: Konten section aktif -->
      <div class="flex-1 min-w-0">
        <!-- Personal Data -->
        <section v-if="active === 0" class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
          <header class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex items-center gap-2">
            <i class="pi pi-user text-indigo-500 text-sm"></i>
            <h2 class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ t('employee.tab_profile') }}</h2>
          </header>
          <div class="p-4 grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-x-6 gap-y-3">
            <DetailRow :label="t('employee.employee_id')" :value="emp.employee_id" />
            <DetailRow :label="t('employee.nik')" :value="emp.nik" />
            <DetailRow :label="t('employee.gender')" :value="genderLabel(emp.gender)" />
            <DetailRow :label="t('employee.mother_name')" :value="emp.mother_name" />
            <DetailRow :label="t('employee.religion')" :value="emp.religion_name || labelOf(religionMap, emp.religion_id)" />
            <DetailRow :label="t('employee.marital_status')" :value="emp.marital_status_name || labelOf(maritalMap, emp.marital_status_id)" />
            <DetailRow :label="t('employee.nationality_type')" :value="emp.nationality_type ? t('employee.' + (emp.nationality_type === 'WNI' ? 'wni' : 'wna')) : ''" />
            <DetailRow :label="t('employee.nationality')" :value="emp.nationality_name || labelOf(nationalityMap, emp.nationality_id)" />
            <DetailRow :label="t('employee.passport')" :value="emp.passport" />
            <DetailRow :label="t('employee.pob')" :value="emp.pob" />
            <DetailRow :label="t('employee.dob')" :value="formatDate(emp.dob, locale)" />
            <DetailRow :label="t('employee.family_id')" :value="emp.family_id" />
            <DetailRow :label="t('employee.phone')" :value="emp.phone_number" />
            <DetailRow :label="t('employee.email')" :value="emp.email" />
            <DetailRow :label="t('employee.linkedin')" :value="emp.linkedin" />
            <DetailRow :label="t('employee.instagram')" :value="emp.ig" />
          </div>
        </section>

        <!-- Addresses -->
        <ListCard v-else-if="active === 1" :title="t('employee.tab_addresses')" icon="pi pi-map-marker" tint="text-sky-500" :empty="t('employee.no_addresses')" :items="emp.addresses">
          <template #default="{ item }">
            <p class="text-sm text-gray-700 dark:text-gray-200">{{ item.address || '-' }}</p>
            <p v-if="item._villageLabel || item._district_name || item._regency_name || item._province_name" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              {{ [item._villageLabel, item._district_name, item._regency_name, item._province_name].filter(Boolean).join(', ') }}
            </p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ item.postal_code ? `${t('employee.postal_code')}: ${item.postal_code}` : '' }}</p>
            <Tag :value="item.type" severity="info" class="!text-[10px] !px-1.5 !py-0.5 mt-1" />
          </template>
        </ListCard>

        <!-- Emergency Contacts -->
        <ListCard v-else-if="active === 2" :title="t('employee.tab_contacts')" icon="pi pi-phone" tint="text-emerald-500" :empty="t('employee.no_contacts')" :items="emp.emergency_contacts">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ item.phone_number }}<template v-if="item.relationship_type_id"> · {{ item.relationship_type_name || labelOf(relationshipMap, item.relationship_type_id) }}</template></p>
            <p v-if="item.address" class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ item.address }}</p>
          </template>
        </ListCard>

        <!-- Family -->
        <ListCard v-else-if="active === 3" :title="t('employee.tab_family')" icon="pi pi-users" tint="text-rose-500" :empty="t('employee.no_family')" :items="emp.families">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              <template v-if="item.relationship_type_id">{{ item.relationship_type_name || labelOf(relationshipMap, item.relationship_type_id) }}</template>
              <template v-if="item.education_name || item.education_id"> · {{ item.education_name || labelOf(educationMap, item.education_id) }}</template>
              <template v-if="item.dob"> · {{ formatDate(item.dob, locale) }}</template>
            </p>
            <p v-if="item.nik" class="text-xs text-gray-400 dark:text-gray-500 mt-0.5 font-mono">{{ item.nik }}</p>
          </template>
        </ListCard>

        <!-- Education -->
        <ListCard v-else-if="active === 4" :title="t('employee.tab_education')" icon="pi pi-graduation-cap" tint="text-amber-500" :empty="t('employee.no_education')" :items="emp.educations">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name || item.education_name || labelOf(educationMap, item.education_id) || '-' }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ item.major || item.major_name || '' }}<template v-if="item.graduation_year"> · {{ item.graduation_year }}</template></p>
          </template>
        </ListCard>

        <!-- Work Experience -->
        <ListCard v-else-if="active === 5" :title="t('employee.tab_experience')" icon="pi pi-briefcase" tint="text-violet-500" :empty="t('employee.no_experience')" :items="emp.experiences">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.company }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ item.position || '' }}<template v-if="item.start_year"> · {{ item.start_year }}<template v-if="item.end_year"> – {{ item.end_year }}</template></template></p>
          </template>
        </ListCard>

        <!-- Documents -->
        <ListCard v-else-if="active === 6" :title="t('employee.tab_documents')" icon="pi pi-file" tint="text-indigo-500" :empty="t('employee.no_documents')" :items="emp.documents">
          <template #default="{ item }">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.name }}</p>
                <p v-if="item.note" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ item.note }}</p>
                <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5 truncate">{{ item.file }}</p>
              </div>
              <a v-if="item.file" :href="item.file" :download="downloadName(item)"
                 class="shrink-0 inline-flex items-center gap-1.5 text-xs font-medium text-indigo-600 dark:text-indigo-400 hover:underline select-none">
                <i class="pi pi-download"></i>
                {{ t('common.download') }}
              </a>
            </div>
          </template>
        </ListCard>

        <!-- Insurance (non-BPJS) -->
        <ListCard v-else-if="active === 7" :title="t('employee.tab_insurance')" icon="pi pi-shield" tint="text-teal-500" :empty="t('employee.no_insurance')" :items="emp.insurances">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.insurance_name || labelOf(insuranceMap, item.insurance_id) || '-' }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ item.number }}<template v-if="item.type"> · {{ item.type }}</template></p>
          </template>
        </ListCard>

        <!-- Bank Accounts -->
        <ListCard v-else-if="active === 8" :title="t('employee.tab_bank')" icon="pi pi-building-columns" tint="text-sky-500" :empty="t('employee.no_bank')" :items="emp.banks">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.bank_name || labelOf(bankMap, item.bank_id) || '-' }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 font-mono">{{ item.account_number }}</p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ item.account_name }}</p>
          </template>
        </ListCard>

        <!-- Employment Record -->
        <ListCard v-else-if="active === 9" :title="t('employee.tab_employment')" icon="pi pi-briefcase" tint="text-cyan-500" :empty="t('employee.no_employment')" :items="emp.employments">
          <template #default="{ item }">
            <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.organization_name || labelOf(organizationMap, item.organization_id) || '-' }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
              <template v-if="item.employment_status_id">{{ item.employment_status_name || labelOf(employmentStatusMap, item.employment_status_id) }}</template>
              <template v-if="item.decision_letter_number"> · {{ t('employee.decision_letter') }}: {{ item.decision_letter_number }}</template>
            </p>
            <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">
              <template v-if="item.effective_date">{{ t('employee.effective_date') }}: {{ item.effective_date }}</template>
              <template v-if="item.effective_end_date"> · {{ t('employee.effective_end_date') }}: {{ item.effective_end_date }}</template>
            </p>
          </template>
        </ListCard>

        <!-- Career Timeline -->
        <div v-else-if="active === 10" class="space-y-4">
          <!-- Current position summary -->
          <div
            v-if="careerData?.current_position"
            class="rounded-lg border border-emerald-200 dark:border-emerald-800/60 bg-emerald-50/50 dark:bg-emerald-900/10 p-3 flex items-center gap-3 flex-wrap"
          >
            <div class="w-9 h-9 rounded-lg bg-emerald-500 flex items-center justify-center shrink-0">
              <i class="pi pi-briefcase text-white text-sm"></i>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs font-medium text-emerald-600 dark:text-emerald-400 uppercase tracking-wider">{{ t('employee_movement.timeline_current_position') }}</p>
              <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">
                {{ careerData.current_position.position_name || '-' }}
                <span v-if="careerData.current_position.organization_name" class="text-gray-400 font-normal"> · {{ careerData.current_position.organization_name }}</span>
                <span v-if="careerData.current_position.employment_status_name" class="text-gray-400 font-normal"> · {{ careerData.current_position.employment_status_name }}</span>
              </p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('employee_movement.effective_date') }}: {{ formatDate(careerData.current_position.effective_date, locale) }}</p>
            </div>
          </div>

          <!-- Timeline events — grouped by year, ASC kronologis (plan §12.8) -->
          <ol v-if="careerData?.timeline?.length" class="relative border-l border-gray-200 dark:border-gray-700 ml-3 space-y-4">
            <template v-for="(group, gi) in careerByYear" :key="gi">
              <li class="relative pl-6">
                <p class="text-[11px] font-bold uppercase tracking-wider text-gray-400 dark:text-gray-500">{{ group.year }}</p>
              </li>
              <li v-for="(ev, i) in group.items" :key="`${gi}-${i}`" class="relative pl-6">
                <span
                  class="absolute -left-[13px] top-0 w-6 h-6 rounded-full flex items-center justify-center ring-4 ring-white dark:ring-gray-800"
                  :class="careerDotClass(ev)"
                >
                  <i :class="careerEventIcon(ev)" class="text-[10px] text-white"></i>
                </span>
                <div class="flex items-center gap-2 flex-wrap">
                  <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ careerEventLabel(ev) }}</p>
                  <Tag v-if="ev.movement_type" :value="careerTypeLabel(ev.movement_type)" :severity="careerTypeSeverity(ev.movement_type)" class="!text-[10px] !px-1.5 !py-0" />
                  <Tag v-if="ev.contract_type" :value="careerTypeLabel(ev.contract_type)" severity="info" class="!text-[10px] !px-1.5 !py-0" />
                </div>
                <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ formatDate(ev.date, locale) }}</p>
                <p v-if="careerEventTitle(ev)" class="text-sm text-gray-600 dark:text-gray-300 mt-0.5">{{ careerEventTitle(ev) }}</p>
                <p v-if="ev.description" class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ ev.description }}</p>
              </li>
            </template>
          </ol>
          <div v-else class="border border-gray-200 dark:border-gray-700 rounded-lg text-center py-8 text-gray-400 dark:text-gray-500">
            <i class="pi pi-history text-2xl mb-2 opacity-50"></i>
            <p class="text-sm">{{ t('employee_movement.timeline_empty') }}</p>
          </div>
        </div>

        <!-- Payroll Profiles (read-only) -->
        <div v-else-if="active === 11" class="space-y-4">
          <ListCard :title="t('payroll.employee_profiles')" icon="pi pi-users" tint="text-sky-500" :empty="t('payroll.profiles_empty')" :items="payroll.profiles">
            <template #default="{ item }">
              <div class="flex items-center gap-2 flex-wrap">
                <Tag :value="item.payroll_group_code" severity="info" class="!text-xs !px-1.5 !py-0.5" />
                <span class="text-sm text-gray-700 dark:text-gray-200">{{ trans(item.payroll_frequency, 'payroll.payroll_frequency_') }}</span>
                <span class="text-sm text-gray-700 dark:text-gray-200">{{ trans(item.payment_method, 'payroll.payment_method_') }}</span>
                <Tag :value="item.is_payroll_active ? t('common.yes') : t('common.no')" :severity="item.is_payroll_active ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" />
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {{ t('payroll.salary_currency') }}: {{ item.salary_currency || 'IDR' }} ·
                {{ t('payroll.effective_start_date') }}: {{ formatDate(item.effective_start_date, locale) }}<template v-if="item.effective_end_date"> – {{ formatDate(item.effective_end_date, locale) }}</template>
              </p>
            </template>
          </ListCard>

          <ListCard :title="t('payroll.bank_profiles')" icon="pi pi-building" tint="text-sky-500" :empty="t('payroll.bank_profiles_empty')" :items="payroll.banks">
            <template #default="{ item }">
              <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.bank_name }} <Tag :value="item.is_primary ? t('common.yes') : t('common.no')" :severity="item.is_primary ? 'success' : 'secondary'" class="!text-[10px] !px-1.5 !py-0.5 ml-1" /></p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 font-mono">{{ item.bank_account_number }}</p>
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ item.bank_account_holder_name }} · {{ t('payroll.effective_start_date') }}: {{ formatDate(item.effective_start_date, locale) }}</p>
            </template>
          </ListCard>

          <ListCard :title="t('payroll.bpjs_profiles')" icon="pi pi-shield" tint="text-teal-500" :empty="t('payroll.bpjs_profiles_empty')" :items="payroll.bpjs">
            <template #default="{ item }">
              <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ trans(item.jkk_risk_class, 'payroll.jkk_risk_class_') }} <Tag :value="trans(item.status, 'payroll.status_')" :severity="item.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-[10px] !px-1.5 !py-0.5 ml-1" /></p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                <template v-if="item.bpjs_health_no">{{ t('payroll.bpjs_health_no') }}: {{ item.bpjs_health_no }}</template>
                <template v-if="item.bpjs_tk_no"> · {{ t('payroll.bpjs_tk_no') }}: {{ item.bpjs_tk_no }}</template>
              </p>
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ t('payroll.effective_start_date') }}: {{ formatDate(item.effective_start_date, locale) }}</p>
            </template>
          </ListCard>

          <ListCard :title="t('payroll.tax_profiles')" icon="pi pi-receipt" tint="text-purple-500" :empty="t('payroll.tax_profiles_empty')" :items="payroll.taxes">
            <template #default="{ item }">
              <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ item.npwp || '-' }} <Tag :value="item.ptkp_status || '-'" severity="info" class="!text-[10px] !px-1.5 !py-0.5 ml-1" /></p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ trans(item.tax_method, 'payroll.tax_method_') }}<template v-if="item.has_npwp"> · {{ t('payroll.has_npwp') }}: {{ t('common.yes') }}</template></p>
              <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ t('payroll.effective_start_date') }}: {{ formatDate(item.effective_start_date, locale) }}</p>
            </template>
          </ListCard>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import DetailRow from './DetailRow.vue'
import ListCard from './ListCard.vue'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const emp = ref({})
const loading = ref(true)
// Section aktif dipakai dari query ?tab= agar refresh tetap di section yang sama
// (bukan selalu kembali ke Profile). Nilai invalid → fallback ke 0 (Profile).
const initialTab = parseInt(route.query.tab, 10)
const active = ref(Number.isInteger(initialTab) && initialTab >= 0 ? initialTab : 0)
const canEdit = computed(() => hasPermission('employee.update'))

function selectSection(i) {
  active.value = i
  router.replace({ query: { ...route.query, tab: String(i) } })
}
const payroll = ref({ profiles: [], banks: [], bpjs: [], taxes: [] })
const payrollCount = computed(() => payroll.value.profiles.length + payroll.value.banks.length + payroll.value.bpjs.length + payroll.value.taxes.length)
const careerData = ref(null)

// ── Career timeline helpers (pola sama EmployeeMovementReports) ──
const careerByYear = computed(() => {
  const entries = careerData.value?.timeline || []
  const groups = []
  let currentYear = null
  for (const ev of entries) {
    const year = String(ev.date || '').slice(0, 4)
    if (year !== currentYear) {
      currentYear = year
      groups.push({ year, items: [] })
    }
    groups[groups.length - 1].items.push(ev)
  }
  return groups
})

function careerTypeLabel(type) {
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : (type || '-')
}

function careerTypeIcon(type) {
  switch (type) {
    case 'promotion': return 'pi pi-arrow-up'
    case 'demotion': return 'pi pi-arrow-down'
    case 'mutation': return 'pi pi-shuffle'
    case 'contract_extension': return 'pi pi-file-edit'
    case 'status_change': return 'pi pi-id-card'
    case 'retirement': return 'pi pi-sun'
    case 'offboarding': return 'pi pi-sign-out'
    default: return 'pi pi-circle'
  }
}

function careerTypeSeverity(type) {
  switch (type) {
    case 'promotion': return 'success'
    case 'demotion': return 'danger'
    case 'mutation': return 'info'
    case 'contract_extension': return 'warning'
    case 'status_change': return 'info'
    case 'retirement': return 'secondary'
    case 'offboarding': return 'danger'
    default: return 'secondary'
  }
}

function careerEventLabel(ev) {
  if (ev.event_type === 'JOINED') return t('employee_movement.timeline_joined')
  if (ev.event_type === 'CONTRACT') return t('employee_movement.timeline_contract')
  if (ev.movement_type) return careerTypeLabel(ev.movement_type)
  return ev.event_type || ''
}

function careerEventTitle(ev) {
  if (ev.event_type === 'MOVEMENT') return ''
  return ev.title || ''
}

function careerEventIcon(ev) {
  if (ev.event_type === 'JOINED') return 'pi pi-user-plus'
  if (ev.event_type === 'CONTRACT') return 'pi pi-file-edit'
  return careerTypeIcon(ev.movement_type)
}

function careerDotClass(ev) {
  if (ev.event_type === 'JOINED') return 'bg-emerald-500'
  if (ev.event_type === 'CONTRACT') return 'bg-amber-500'
  switch (ev.movement_type) {
    case 'promotion': return 'bg-emerald-500'
    case 'demotion': return 'bg-red-500'
    case 'mutation': return 'bg-sky-500'
    case 'contract_extension': return 'bg-amber-500'
    case 'status_change': return 'bg-indigo-500'
    case 'retirement': return 'bg-gray-400'
    case 'offboarding': return 'bg-rose-500'
    default: return 'bg-slate-400'
  }
}

async function loadCareerHistory() {
  if (!emp.value.id) return
  try {
    const res = await api.get(`/api/v1/tenant/employee-movements/employees/${emp.value.id}/career-history`)
    careerData.value = res.data?.data || null
  } catch { /* section timeline tampil kosong bila endpoint gagal */ }
}

// ── Sub-menu navigasi (urutan sama dengan wizard form) ──
const sections = computed(() => [
  { labelKey: 'employee.tab_profile', icon: 'pi pi-user' },
  { labelKey: 'employee.tab_addresses', icon: 'pi pi-map-marker', count: emp.value.addresses?.length || 0 },
  { labelKey: 'employee.tab_contacts', icon: 'pi pi-phone', count: emp.value.emergency_contacts?.length || 0 },
  { labelKey: 'employee.tab_family', icon: 'pi pi-users', count: emp.value.families?.length || 0 },
  { labelKey: 'employee.tab_education', icon: 'pi pi-graduation-cap', count: emp.value.educations?.length || 0 },
  { labelKey: 'employee.tab_experience', icon: 'pi pi-briefcase', count: emp.value.experiences?.length || 0 },
  { labelKey: 'employee.tab_documents', icon: 'pi pi-file', count: emp.value.documents?.length || 0 },
  { labelKey: 'employee.tab_insurance', icon: 'pi pi-shield', count: emp.value.insurances?.length || 0 },
  { labelKey: 'employee.tab_bank', icon: 'pi pi-building-columns', count: emp.value.banks?.length || 0 },
  { labelKey: 'employee.tab_employment', icon: 'pi pi-briefcase', count: emp.value.employments?.length || 0 },
  { labelKey: 'employee_movement.career_timeline', icon: 'pi pi-history', count: careerData.value?.timeline?.length || 0 },
  { labelKey: 'employee.wizard_step_payroll', icon: 'pi pi-dollar', count: payrollCount }
])

// ── Label maps (dari ref data) ──
const religionMap = ref({})
const maritalMap = ref({})
const nationalityMap = ref({})
const relationshipMap = ref({})
const educationMap = ref({})
const insuranceMap = ref({})
const bankMap = ref({})
const organizationMap = ref({})
const employmentStatusMap = ref({})

// Dipanggil dari template — ref sudah di-unwrap otomatis, jadi map berupa objek polos.
function labelOf(map, id) {
  if (!id || !map) return ''
  return map[id] || id
}

function genderLabel(v) {
  if (!v) return ''
  return t('employee.gender_' + v.toLowerCase())
}

// Nama file unduhan: nama dokumen + ekstensi dari path file.
function downloadName(item) {
  const file = item?.file || ''
  const ext = file.includes('.') ? file.substring(file.lastIndexOf('.')) : ''
  const base = (item?.name || 'document').replace(/\.[^.]+$/, '')
  return base + ext
}

// Terjemahan label enum payroll (dengan fallback ke nilai mentah).
function trans(v, prefix) {
  const key = prefix + String(v || '').toLowerCase()
  return t(key) !== key ? t(key) : (v || '-')
}

async function fetchAll(endpoint) {
  const all = []
  let page = 1
  while (true) {
    const res = await api.get(endpoint, { params: { page, per_page: 100 } })
    const rows = res.data?.data || []
    all.push(...rows)
    if (!rows.length || all.length >= (res.data?.total || 0)) break
    page++
  }
  return all
}

async function loadPayrollProfiles() {
  if (!emp.value.id) return
  try {
    await doLoadPayrollProfiles()
  } catch { /* section payroll tampil kosong bila endpoint gagal */ }
}

async function doLoadPayrollProfiles() {
  const [pps, banks, bpjs, taxes] = await Promise.all([
    fetchAll('/api/v1/tenant/payroll/employee-payroll-profiles'),
    fetchAll('/api/v1/tenant/payroll/employee-bank-profiles'),
    fetchAll('/api/v1/tenant/payroll/employee-bpjs-profiles'),
    fetchAll('/api/v1/tenant/payroll/employee-tax-profiles')
  ])
  const profiles = pps.filter(p => p.employee_id === emp.value.id)
  const profileIds = new Set(profiles.map(p => p.id))
  payroll.value = {
    profiles,
    banks: banks.filter(b => profileIds.has(b.employee_payroll_profile_id)),
    bpjs: bpjs.filter(b => profileIds.has(b.employee_payroll_profile_id)),
    taxes: taxes.filter(x => profileIds.has(x.employee_payroll_profile_id))
  }
}

async function loadRefData() {
  try {
    await doLoadRefData()
  } catch { /* label map tetap kosong — label jatuh ke ID asli */ }
}

async function doLoadRefData() {
  // allSettled: satu endpoint gagal (mis. 403 untuk role terbatas) tidak
  // menggagalkan map label lainnya — tiap map diisi independen.
  const results = await Promise.allSettled([
    api.get('/api/v1/tenant/settings/religions?per_page=100'),
    api.get('/api/v1/tenant/settings/marital-statuses?per_page=100'),
    api.get('/api/v1/tenant/settings/nationalities?per_page=250'),
    api.get('/api/v1/tenant/settings/relationship-types?per_page=100'),
    api.get('/api/v1/tenant/settings/educations?per_page=100'),
    api.get('/api/v1/tenant/settings/insurances?per_page=100'),
    api.get('/api/v1/tenant/settings/banks?per_page=200'),
    api.get('/api/v1/tenant/organizations?tree=true'),
    api.get('/api/v1/tenant/settings/employment-statuses?per_page=100')
  ])
  const rows = i => (results[i]?.status === 'fulfilled' ? (results[i].value?.data?.data || []) : [])
  religionMap.value = Object.fromEntries(rows(0).map(r => [r.id, r.name]))
  maritalMap.value = Object.fromEntries(rows(1).map(m => [m.id, m.name]))
  nationalityMap.value = Object.fromEntries(rows(2).map(n => [n.code, n.name]))
  relationshipMap.value = Object.fromEntries(rows(3).map(r => [r.id, r.name]))
  educationMap.value = Object.fromEntries(rows(4).map(e => [e.id, e.name]))
  insuranceMap.value = Object.fromEntries(rows(5).map(x => [x.id, x.name]))
  bankMap.value = Object.fromEntries(rows(6).map(b => [b.id, b.name]))
  employmentStatusMap.value = Object.fromEntries(rows(8).map(x => [x.id, x.name]))
  const orgList = {}
  function flattenOrgTree(nodes) {
    if (!nodes) return
    if (Array.isArray(nodes)) { nodes.forEach(n => flattenOrgTree(n)); return }
    orgList[nodes.id] = nodes.nomenclature
    if (nodes.children) flattenOrgTree(nodes.children)
  }
  flattenOrgTree(rows(7))
  organizationMap.value = orgList
}

async function loadVillageLabels(addressList) {
  for (const addr of addressList) {
    if (!addr.village_id || addr._villageLabel) continue
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

onMounted(async () => {
  try {
    const res = await api.get(`/api/v1/tenant/employees/${route.params.id}`)
    emp.value = res.data?.data || {}
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
    router.replace('/employees')
    return
  } finally {
    // Halaman langsung tampil setelah data employee — label ref & payroll dimuat di background
    loading.value = false
  }
  // Peningkatan (tidak memblokir halaman & tidak boleh menggagalkan render)
  if (emp.value.addresses) loadVillageLabels(emp.value.addresses)
  loadRefData()
  loadPayrollProfiles()
  loadCareerHistory()
})
</script>
