<template>
  <div class="space-y-4">
    <!-- Page Header: navigasi view dashboard -->
    <div class="flex items-center justify-end">
      <SelectButton
        v-model="view"
        :options="viewOptions"
        optionLabel="label"
        optionValue="value"
        size="small"
      />
    </div>

    <!-- Sub-halaman dashboard. KeepAlive membuat tiap view hanya memuat data
         sekali saat pertama dibuka; pindah view lalu kembali tidak reload. -->
    <KeepAlive>
      <MyDashboard v-if="view === 'my'" />
      <JobManagementDashboard v-else-if="view === 'job'" />
      <TalentDashboard v-else-if="view === 'talent'" />
      <EmploymentDashboard v-else-if="view === 'employment'" />
      <HRAttendanceLeaveDashboard v-else-if="view === 'hr' && hrViewAllowed" />
    </KeepAlive>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'

import SelectButton from 'primevue/selectbutton'

import MyDashboard from '@/components/dashboard/MyDashboard.vue'
import JobManagementDashboard from '@/components/dashboard/JobManagementDashboard.vue'
import TalentDashboard from '@/components/dashboard/TalentDashboard.vue'
import EmploymentDashboard from '@/components/dashboard/EmploymentDashboard.vue'
import HRAttendanceLeaveDashboard from '@/components/dashboard/HRAttendanceLeaveDashboard.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuth()
const { hasExactPermission } = auth

// ── Navigasi view — di-persist ke query ?view= agar refresh tetap di view yang sama ──
// View HR (Absensi & Cuti) menampilkan data org-wide, jadi hanya muncul bila
// role punya attendance.report.view (strict — sama seperti gating Reports).
const baseViews = ['my', 'job', 'talent', 'employment']
const hrViewAllowed = computed(() => hasExactPermission('attendance.report.view'))
const VIEWS = computed(() => (hrViewAllowed.value ? [...baseViews, 'hr'] : baseViews))
const initialView = route.query.view && VIEWS.value.includes(route.query.view) ? route.query.view : 'my'
const view = ref(initialView)

const viewOptions = computed(() => {
  const opts = [
    { label: t('dashboard.view_my'), value: 'my' },
    { label: t('dashboard.view_job'), value: 'job' },
    { label: t('dashboard.view_talent'), value: 'talent' },
    { label: t('dashboard.view_employment'), value: 'employment' }
  ]
  if (hrViewAllowed.value) opts.push({ label: t('dashboard.view_hr'), value: 'hr' })
  return opts
})

// Data tiap view dimuat sendiri oleh komponennya (onMounted + KeepAlive);
// watch di sini hanya mem-persist view aktif ke URL.
watch(view, (v) => {
  router.replace({ query: { ...route.query, view: v } })
})
</script>
