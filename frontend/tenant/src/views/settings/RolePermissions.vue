<template>
  <div class="space-y-4">
    <!-- Card nama role + aksi — sticky, tetap terlihat saat scroll -->
    <div class="sticky top-0 z-10 bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 shadow-sm">
      <div class="flex items-center justify-between flex-wrap gap-2">
        <div>
          <p class="text-sm font-semibold text-gray-800 dark:text-gray-100">
            {{ t('rbac.assign_permissions') }} — <span class="text-indigo-600 dark:text-indigo-400">{{ role?.name || '' }}</span>
          </p>
          <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ selectedCount }} / {{ permissions.length }} {{ t('rbac.permissions_count') }}</p>
        </div>
        <div class="flex items-center gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="router.push('/settings/rbac')" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving || loading" @click="handleSave" />
        </div>
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="4" />
    <DataTable
      v-else
      :value="permissions"
      rowGroupMode="subheader"
      groupRowsBy="resource"
      sortField="resource"
      :sortOrder="1"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="text-center text-sm text-gray-400 py-8">{{ t('rbac.empty_permissions') }}</div>
      </template>
      <!-- Subheader per resource: nama resource + All/Clear + pilih semua -->
      <template #subheader="{ data }">
        <div class="flex items-center justify-between">
          <span class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ data.resource }}</span>
          <div class="flex items-center gap-1.5">
            <Button
              :label="groupAllSelected(data.resource) ? t('common.clear') : t('common.all')"
              severity="secondary"
              text
              size="small"
              class="!text-[11px] !p-0.5"
              @click="toggleResource(data.resource)"
            />
            <Checkbox :binary="true" :model-value="groupAllSelected(data.resource)" @update:model-value="toggleResource(data.resource)" />
          </div>
        </div>
      </template>
      <Column field="action" :header="t('rbac.permission')">
        <template #body="{ data }"><span class="text-gray-700 dark:text-gray-200">{{ data.action }}</span></template>
      </Column>
      <Column :header="t('common.select')" style="width:80px" frozen alignFrozen="right">
        <template #body="{ data }">
          <Checkbox :binary="true" :model-value="selected[data.id]" @update:model-value="v => selected[data.id] = v" />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import Button from 'primevue/button'
import Checkbox from 'primevue/checkbox'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import SkeletonTable from '@/components/SkeletonTable.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const roleId = route.params.id
const role = ref(null)
const permissions = ref([])
const selected = ref({})
const loading = ref(true)
const saving = ref(false)

const skeletonColumns = [
  { type: 'text', width: 'w-40', headerWidth: 'w-40' },
  { type: 'icons', count: 1, headerWidth: 'w-24' }
]

const selectedCount = computed(() => Object.values(selected.value).filter(Boolean).length)

function resourceOf(p) {
  return p.resource || 'other'
}

// Semua permission di satu resource terpilih?
function groupAllSelected(resource) {
  const items = permissions.value.filter(p => resourceOf(p) === resource)
  return items.length > 0 && items.every(i => selected.value[i.id])
}

// Toggle pilih semua / kosongkan satu resource.
function toggleResource(resource) {
  const all = groupAllSelected(resource)
  for (const p of permissions.value) {
    if (resourceOf(p) === resource) {
      selected.value[p.id] = !all
    }
  }
}

async function loadData() {
  loading.value = true
  try {
    const [rolesRes, permsRes] = await Promise.all([
      api.get('/api/v1/tenant/rbac/roles'),
      api.get('/api/v1/tenant/rbac/permissions')
    ])
    const roles = rolesRes.data?.data || []
    role.value = roles.find(r => r.id === roleId) || null
    permissions.value = permsRes.data?.data || []
    selected.value = {}
    for (const pid of (role.value?.permission_ids || [])) {
      selected.value[pid] = true
    }
    if (!role.value) {
      toast.add({ severity: 'warn', summary: t('message.error'), detail: t('rbac.empty_roles'), life: 4000 })
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const ids = Object.keys(selected.value).filter(k => selected.value[k])
    await api.put(`/api/v1/tenant/rbac/roles/${roleId}/permissions`, { permission_ids: ids })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('rbac.permissions_updated'), life: 3000 })
    router.push('/settings/rbac')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
