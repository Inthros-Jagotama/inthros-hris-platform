<template>
  <div class="max-w-full mx-auto">
    <!-- Loading -->
    <div v-if="pageLoading" class="space-y-4">
      <div class="flex gap-4">
        <div class="flex-1 space-y-3">
          <div v-for="n in 6" :key="n" class="h-8 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div v-else>
      <!-- Organization Info (hanya tampil bila dibuka dari daftar organisasi) -->
      <div v-if="orgId" class="mb-4">
        <div class="flex items-center gap-3 p-3 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
          <div class="w-9 h-9 rounded-lg shrink-0 flex items-center justify-center bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
            <i class="pi pi-building text-sm"></i>
          </div>
          <div class="min-w-0">
            <div class="text-[10px] text-gray-400 uppercase tracking-wider mb-0.5">{{ t('organization.title') }}</div>
            <div class="text-sm font-semibold text-navy-800 dark:text-gray-100 truncate">{{ orgName }}</div>
            <div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 font-mono">{{ orgCode }}</div>
          </div>
        </div>
      </div>

      <JobValueSection :key="type" :type="type" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { normalizeJobValueType } from '@/utils/jobValues'
import api from '@/services/api'
import JobValueSection from './JobValueSection.vue'

const { t } = useI18n()
const route = useRoute()

// ── Tipe dari URL param :type (mis. /job-management/values/education) ──
// Reactive: ikut perubahan :type saat navigasi antar tipe (back/forward)
// → :key="type" di JobValueSection memaksa remount dengan tipe yang benar
const type = computed(() => normalizeJobValueType(route.params.type))

// Org info opsional (bila dibuka dari daftar organisasi via ?org_id=)
const orgId = route.query.org_id || ''
const orgName = ref('')
const orgCode = ref('')
const pageLoading = ref(true)

// ── Fetch Organization Info ──
async function loadOrgInfo() {
  if (!orgId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${orgId}`)
    const data = res.data?.data
    if (data) {
      orgName.value = data.nomenclature
      orgCode.value = `${data.full_code}`
    }
  } catch {
    orgName.value = 'Unknown'
    orgCode.value = ''
  }
}

onMounted(async () => {
  pageLoading.value = true
  try {
    await loadOrgInfo()
  } finally {
    pageLoading.value = false
  }
})
</script>
