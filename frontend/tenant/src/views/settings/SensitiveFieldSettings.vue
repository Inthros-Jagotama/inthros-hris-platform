<template>
  <div class="space-y-4">
    <SkeletonTable v-if="loading" :columns="skeletonColumns" :rows="8" />

    <DataTable
      v-else
      :value="settings"
      size="small"
      class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      sortField="field_key"
      :sortOrder="1"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-lock text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('common.no_data') }}</p>
        </div>
      </template>
      <Column field="field_key" :header="t('sensitive_field.field_name')" sortable>
        <template #body="{ data }">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.field_key }}</span>
        </template>
      </Column>
      <Column :header="t('sensitive_field.encryption_enabled')" style="width:160px">
        <template #body="{ data }">
          <ToggleSwitch
            :modelValue="data.is_encryption_enabled"
            :disabled="!!savingKeys[data.field_key] || !canManage"
            @update:modelValue="(val) => toggleField(data.field_key, val)"
          />
        </template>
      </Column>
    </DataTable>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ToggleSwitch from 'primevue/toggleswitch'
import SkeletonTable from '@/components/SkeletonTable.vue'
import { useAuth } from '@/stores/auth'

const { t } = useI18n()
const toast = useToast()

// Mengubah toggle enkripsi at-rest adalah aksi setingkat admin — dipisahkan
// dari permission melihat halaman ini (lihat migrasi tenant 154).
const auth = useAuth()
const canManage = computed(() => auth.hasPermission('setting.sensitive-fields.manage'))

const settings = ref([])
const loading = ref(false)
const savingKeys = ref({})

const skeletonColumns = [
  { type: 'text', width: 'w-56', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-20' }
]

async function loadSettings() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/employees/settings/sensitive-fields')
    settings.value = res.data?.data || []
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function toggleField(fieldKey, enabled) {
  const row = settings.value.find((s) => s.field_key === fieldKey)
  const previous = row?.is_encryption_enabled
  if (row) row.is_encryption_enabled = enabled
  savingKeys.value = { ...savingKeys.value, [fieldKey]: true }
  try {
    await api.put(`/api/v1/tenant/employees/settings/sensitive-fields/${fieldKey}`, {
      is_encryption_enabled: enabled
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('sensitive_field.updated_toast'), life: 3000 })
  } catch (e) {
    if (row) row.is_encryption_enabled = previous
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    const { [fieldKey]: _removed, ...rest } = savingKeys.value
    savingKeys.value = rest
  }
}

onMounted(loadSettings)
</script>
