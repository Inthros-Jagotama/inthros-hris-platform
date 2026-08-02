<template>
  <div :class="depth === 0 ? 'org-table-scroll' : ''">
    <DataTable
      :value="nodes"
      v-model:expandedRows="expandedRows"
      dataKey="id"
      :class="depth === 0
        ? '!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden'
        : '!text-xs !border-0 !shadow-none'"
      stripedRows
      responsiveLayout="scroll"
    >
      <template #empty>
        <div v-if="depth === 0" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-sitemap text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t('organization.empty_title') }}</p>
          <p class="text-sm mt-1">{{ t('organization.empty_tree') }}</p>
        </div>
      </template>
      <Column :expander="true" style="width: 40px" />
      <Column field="nomenclature" :header="depth === 0 ? t('organization.nomenclature') : ''" style="min-width: 200px">
        <template #body="{ data }">
          <div class="flex items-center gap-2">
            <i class="pi pi-folder-open text-amber-500 text-xs"></i>
            <span class="font-medium text-gray-800 dark:text-gray-100">{{ data.nomenclature }}</span>
            <Tag
              v-if="data.children?.length"
              :value="data.children.length"
              severity="info"
              class="!text-[10px] !px-1 !py-0 !min-w-[1.1rem]"
              rounded
            />
          </div>
        </template>
      </Column>
      <Column field="code" :header="depth === 0 ? t('organization.code') : ''" style="width: 120px">
        <template #body="{ data }">
          <Tag :value="data.code" severity="info" class="!text-xs" />
        </template>
      </Column>
      <Column field="full_code" :header="depth === 0 ? t('organization.full_code') : ''" style="width: 160px">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400 text-xs font-mono">{{ data.full_code }}</span>
        </template>
      </Column>
      <Column field="level" :header="depth === 0 ? t('organization.level') : ''" style="width: 80px">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400">{{ data.level }}</span>
        </template>
      </Column>
      <Column field="sort_order" :header="depth === 0 ? t('organization.sort_order') : ''" style="width: 80px">
        <template #body="{ data }">
          <span class="text-gray-500 dark:text-gray-400">{{ data.sort_order }}</span>
        </template>
      </Column>
      <Column :header="depth === 0 ? t('common.actions') : ''" style="width: 130px" frozen alignFrozen="right">
        <template #body="{ data }">
          <div class="flex items-center gap-1">
            <Button icon="pi pi-plus" v-tooltip.top="t('organization.add_child')" severity="secondary" text size="small" class="!p-1" @click="$emit('add-child', data)" />
            <Button icon="pi pi-pencil" v-tooltip.top="t('common.edit')" severity="secondary" text size="small" class="!p-1" @click="$emit('edit', data)" />
            <Button icon="pi pi-trash" v-tooltip.top="t('common.delete')" severity="danger" text size="small" class="!p-1" @click="$emit('delete', data)" />
          </div>
        </template>
      </Column>

      <!-- Expansion rekursif: komponen memanggil dirinya sendiri untuk semua level -->
      <template #expansion="{ data }">
        <div v-if="data.children?.length" class="pl-6 pr-2 py-1">
          <OrgTreeTable
            :nodes="data.children"
            v-model:expandedRows="expandedRows"
            :depth="depth + 1"
            @add-child="$emit('add-child', $event)"
            @edit="$emit('edit', $event)"
            @delete="$emit('delete', $event)"
          />
        </div>
        <div v-else class="pl-6 pr-2 py-2 text-xs text-gray-400 italic">
          {{ t('organization.no_children') }}
        </div>
      </template>
    </DataTable>
  </div>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import { useI18n } from '@/composables/useI18n'

defineProps({
  nodes: { type: Array, default: () => [] },
  depth: { type: Number, default: 0 }
})
defineEmits(['add-child', 'edit', 'delete'])
const expandedRows = defineModel('expandedRows', { type: Object, default: () => ({}) })
const { t } = useI18n()
</script>

<style scoped>
.org-table-scroll :deep(.p-datatable-wrapper) {
  max-height: calc(100vh - 260px);
}
:deep(.p-datatable .p-datatable-tbody > tr) {
  transition: background 0.15s ease;
}
:deep(.p-datatable .p-datatable-tbody > tr:hover) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-datatable .p-datatable-tbody > tr:hover) {
  background: rgba(16, 185, 129, 0.08) !important;
}
:deep(.p-datatable .p-datatable-tbody > tr.p-row-expanded) {
  background: #f0fdf4 !important;
}
:deep(.p-dark .p-datatable .p-datatable-tbody > tr.p-row-expanded) {
  background: rgba(16, 185, 129, 0.08) !important;
}
</style>
