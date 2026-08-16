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
          <IconField iconPosition="left">
            <InputIcon class="pi pi-search" />
            <InputText v-model="filterQuery" size="small" :placeholder="t('common.search')" class="!w-48" />
          </IconField>
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="router.push('/settings/rbac')" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving || loading" @click="handleSave" />
        </div>
      </div>
    </div>

    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="4" />
    <DataTable
      v-else
      :value="filteredGroups"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    >
      <template #empty>
        <div class="text-center text-sm text-gray-400 py-8">{{ t('rbac.empty_permissions') }}</div>
      </template>
      <!-- 1 baris = 1 module: nama module → pilih semua → chips permission -->
      <Column field="resource" :header="t('rbac.module')" style="width:200px">
        <template #body="{ data }"><span class="font-medium text-gray-800 dark:text-gray-100">{{ moduleLabel(data.resource) }}</span></template>
      </Column>
      <Column :header="t('common.select')" style="width:110px">
        <template #body="{ data }">
          <div class="flex items-center gap-1.5">
            <Checkbox :binary="true" :model-value="data.allSelected" @update:model-value="toggleResource(data.resource)" />
            <Button
              :label="data.allSelected ? t('common.clear') : t('common.all')"
              severity="secondary"
              text
              size="small"
              class="!text-[11px] !p-0.5"
              @click="toggleResource(data.resource)"
            />
          </div>
        </template>
      </Column>
      <!-- Satu kolom per aksi (create, update, ...) — isi berupa switch tanpa label -->
      <Column v-for="act in allActions" :key="act" :header="act" style="width:90px">
        <template #body="{ data }">
          <div class="flex justify-center">
            <ToggleSwitch
              v-if="data.byAction[act]"
              :model-value="selected[data.byAction[act].id]"
              @update:model-value="v => selected[data.byAction[act].id] = v"
            />
            <span v-else class="text-gray-300 dark:text-gray-600">—</span>
          </div>
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
import ToggleSwitch from 'primevue/toggleswitch'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import InputText from 'primevue/inputtext'
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
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

// Label bilingual nama module (rbac.modules.<resource>) — fallback ke slug.
function moduleLabel(resource) {
  const key = `rbac.modules.${resource}`
  return t(key) !== key ? t(key) : resource
}

const filterQuery = ref('')

// Filter module by nama (label bilingual atau slug).
const filteredGroups = computed(() => {
  const q = filterQuery.value.trim().toLowerCase()
  if (!q) return permissionGroups.value
  return permissionGroups.value.filter(g =>
    moduleLabel(g.resource).toLowerCase().includes(q) || g.resource.toLowerCase().includes(q)
  )
})

// Semua aksi yang ada lintas module (create, update, ...) — urutan kemunculan
// pertama — menjadi kolom-kolom tabel.
const allActions = computed(() => {
  const seen = []
  for (const p of permissions.value) {
    if (!seen.includes(p.action)) seen.push(p.action)
  }
  return seen
})

// Grup permission per module (resource) — 1 baris = 1 module.
const permissionGroups = computed(() => {
  const map = {}
  for (const p of permissions.value) {
    const res = p.resource || 'other'
    if (!map[res]) map[res] = { resource: res, items: [] }
    map[res].items.push(p)
  }
  return Object.values(map).map(g => {
    const byAction = {}
    for (const i of g.items) byAction[i.action] = i
    return {
      ...g,
      byAction,
      allSelected: g.items.length > 0 && g.items.every(i => selected.value[i.id])
    }
  })
})

// Toggle pilih semua / kosongkan satu module.
function toggleResource(resource) {
  const group = permissionGroups.value.find(g => g.resource === resource)
  if (!group) return
  const all = group.allSelected
  for (const p of group.items) {
    selected.value[p.id] = !all
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
