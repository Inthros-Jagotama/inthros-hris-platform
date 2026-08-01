<template>
  <div class="space-y-4">
    <!-- Back navigation -->
    <div class="flex items-center gap-2">
      <Button icon="pi pi-arrow-left" severity="secondary" text size="small" v-tooltip.bottom="t('common.back')" @click="goBack" />
      <span class="text-sm text-gray-500 dark:text-gray-400 font-medium">{{ t('companies.detail_title') }}</span>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-spin pi-spinner text-2xl mb-2"></i>
      <p class="text-sm">{{ t('common.loading') }}</p>
    </div>

    <!-- Not found -->
    <div v-else-if="!company" class="flex flex-col items-center justify-center py-16 text-gray-400 dark:text-gray-500">
      <i class="pi pi-building text-3xl mb-2 opacity-50"></i>
      <p class="text-sm font-medium">{{ t('companies.detail_not_found') }}</p>
      <Button :label="t('common.back')" severity="secondary" outlined size="small" class="mt-3" @click="goBack" />
    </div>

    <template v-else>
      <!-- Header card -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4 space-y-3">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="flex items-center gap-3 min-w-0">
            <div class="w-10 h-10 rounded-lg bg-indigo-100 dark:bg-indigo-500/20 flex items-center justify-center shrink-0">
              <i class="pi pi-building text-indigo-500 dark:text-indigo-400 text-lg"></i>
            </div>
            <div class="min-w-0">
              <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100 uppercase truncate">{{ company.name }}</h2>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">{{ company.slug }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2 flex-wrap">
            <Tag :value="company.status" :severity="statusSeverity(company.status)" class="!text-xs !px-1.5 !py-0.5" />
            <CompanyActions :company="company" mode="buttons" @updated="loadDetail" />
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <!-- Basic Info -->
        <section class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-1.5">
            <i class="pi pi-info-circle text-indigo-400 text-sm"></i>
            {{ t('companies.detail_basic_info') }}
          </h3>
          <dl class="space-y-2 text-sm">
            <InfoRow :label="t('companies.company_name')" :value="company.name" />
            <InfoRow :label="t('companies.slug')" :value="company.slug" />
            <InfoRow :label="t('common.status')" :value="company.status" :tag="statusSeverity(company.status)" />
            <InfoRow :label="t('companies.email')" :value="company.email || '—'" />
            <InfoRow :label="t('companies.phone')" :value="company.phone || '—'" />
            <InfoRow :label="t('companies.address')" :value="company.address || '—'" />
            <InfoRow :label="t('companies.npwp')" :value="company.npwp || '—'" />
            <InfoRow :label="t('companies.nib')" :value="company.nib || '—'" />
          </dl>
        </section>

        <!-- License Info -->
        <section class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-1.5">
            <i class="pi pi-id-card text-indigo-400 text-sm"></i>
            {{ t('companies.license_section') }}
          </h3>
          <template v-if="company.license_info">
            <dl class="space-y-2 text-sm">
              <InfoRow :label="t('companies.license_plan')" :value="company.license_info.plan_type" :tag="planSeverity(company.license_info.plan_type)" />
              <InfoRow :label="t('companies.license_key_label')" :value="company.license_info.license_key" :mono="true" />
              <InfoRow :label="t('companies.package')" :value="company.license_info.package_id || '—'" :mono="true" />
            </dl>
          </template>
          <p v-else class="text-sm text-gray-400 dark:text-gray-500 italic">—</p>
        </section>

        <!-- Provisioning Info -->
        <section class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-1.5">
            <i class="pi pi-database text-indigo-400 text-sm"></i>
            {{ t('companies.detail_database') }}
          </h3>
          <template v-if="company.provisioning_info">
            <dl class="space-y-2 text-sm">
              <InfoRow :label="t('companies.provision_header')" :value="company.provisioning_info.db_name || '—'" :mono="true" />
              <InfoRow :label="t('companies.detail_driver')" :value="company.provisioning_info.driver || '—'" />
              <InfoRow
                :label="t('companies.detail_is_active')"
                :value="company.provisioning_info.is_active !== false ? t('common_status.yes') : t('common_status.no')"
              />
              <InfoRow :label="t('companies.detail_provisioned')" :value="company.provisioning_info.provisioned ? t('common_status.yes') : t('common_status.no')" />
            </dl>
          </template>
          <p v-else class="text-sm text-gray-400 dark:text-gray-500 italic">—</p>
        </section>

        <!-- Admin User -->
        <section class="border border-gray-200 dark:border-gray-700 rounded-lg bg-white dark:bg-gray-800 p-4">
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-3 flex items-center gap-1.5">
            <i class="pi pi-user text-indigo-400 text-sm"></i>
            {{ t('companies.admin_section_title') }}
          </h3>
          <template v-if="company.admin_user">
            <dl class="space-y-2 text-sm">
              <InfoRow :label="t('companies.admin_name')" :value="company.admin_user.name || '—'" />
              <InfoRow :label="t('companies.admin_email')" :value="company.admin_user.email || '—'" />
              <InfoRow :label="t('companies.admin_role')" :value="company.admin_user.role || '—'" />
            </dl>
          </template>
          <p v-else class="text-sm text-gray-400 dark:text-gray-500 italic">—</p>
        </section>
      </div>

      <!-- Timestamps -->
      <div class="flex flex-wrap items-center gap-4 text-xs text-gray-400 dark:text-gray-500">
        <span>{{ t('companies.created') }}: {{ company.created_at || company.createdAt || '-' }}</span>
        <span>{{ t('companies.detail_updated') }}: {{ company.updated_at || company.updatedAt || '-' }}</span>
      </div>
    </template>

  </div>
</template>

<script setup>
import { ref, onMounted, watch, defineComponent, h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import CompanyActions from '@/components/CompanyActions.vue'

// InfoRow — baris label:nilai untuk section detail.
const InfoRow = defineComponent({
  name: 'InfoRow',
  props: {
    label: { type: String, default: '' },
    value: { type: [String, Number], default: '—' },
    mono: { type: Boolean, default: false },
    tag: { type: String, default: '' }
  },
  setup(props) {
    return () => h('div', { class: 'flex items-start justify-between gap-4 py-1' }, [
      h('dt', { class: 'text-gray-500 dark:text-gray-400 shrink-0' }, props.label),
      h('dd', { class: 'text-right font-medium text-gray-700 dark:text-gray-200 min-w-0 break-words' }, [
        props.tag
          ? h(Tag, { value: props.value, severity: props.tag, class: '!text-xs !px-1.5 !py-0.5' })
          : h('span', { class: props.mono ? 'font-mono text-xs' : '' }, props.value)
      ])
    ])
  }
})

const toast = useToast()
const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const company = ref(null)
const loading = ref(true)

async function loadDetail() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/platform/companies/${route.params.id}`)
    company.value = res.data?.data || res.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.failed_to_load'), life: 3000 })
    company.value = null
  } finally {
    loading.value = false
  }
}

onMounted(loadDetail)

// Reload jika berpindah antar detail company (mis. via back/forward atau link detail→detail)
watch(() => route.params.id, () => {
  if (route.name === 'CompanyDetail') loadDetail()
})

function goBack() {
  router.push('/companies')
}

function statusSeverity(status) {
  switch (status) {
    case 'active': return 'success'
    case 'suspended': return 'warn'
    case 'terminated': return 'danger'
    default: return 'info'
  }
}

function planSeverity(plan) {
  switch (plan?.toLowerCase()) {
    case 'enterprise': return 'danger'
    case 'professional': return 'warn'
    case 'basic': return 'info'
    case 'trial': return 'success'
    default: return 'info'
  }
}

</script>
