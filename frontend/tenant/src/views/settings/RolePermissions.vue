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
    <div v-else-if="permissionGroups.length === 0" class="text-center text-sm text-gray-400 py-8">
      {{ t('rbac.empty_permissions') }}
    </div>
    <!-- Setiap grup permission berdiri sendiri sebagai card, di luar card utama -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-3">
      <div v-for="group in permissionGroups" :key="group.resource" class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden bg-white dark:bg-gray-800">
        <div class="flex items-center justify-between px-3 py-2 bg-gray-50 dark:bg-gray-800/60 border-b border-gray-200 dark:border-gray-700">
          <span class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ group.resource }}</span>
          <div class="flex items-center gap-1.5">
            <Button
              :label="group.allSelected ? t('common.clear') : t('common.all')"
              severity="secondary"
              text
              size="small"
              class="!text-[11px] !p-0.5"
              @click="toggleGroup(group)"
            />
            <Checkbox :binary="true" :model-value="group.allSelected" @update:model-value="toggleGroup(group)" />
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-1 px-3 py-2">
          <label v-for="p in group.items" :key="p.id" class="flex items-center gap-2 py-0.5 cursor-pointer hover:text-indigo-600 dark:hover:text-indigo-400">
            <Checkbox :binary="true" :model-value="selected[p.id]" @update:model-value="v => selected[p.id] = v" />
            <span class="text-sm">{{ p.action }}</span>
          </label>
        </div>
      </div>
    </div>
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

const permissionGroups = computed(() => {
  const map = {}
  for (const p of permissions.value) {
    const res = p.resource || 'other'
    if (!map[res]) map[res] = { resource: res, items: [] }
    map[res].items.push(p)
  }
  return Object.values(map).map(g => ({
    ...g,
    allSelected: g.items.length > 0 && g.items.every(i => selected.value[i.id])
  }))
})

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

function toggleGroup(group) {
  const allSelected = group.allSelected
  for (const p of group.items) {
    selected.value[p.id] = !allSelected
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
