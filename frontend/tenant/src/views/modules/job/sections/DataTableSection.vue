<template>
  <SkeletonTable v-if="loading" :columns="skeletonCols" :rows="8" />
  <DataTable v-else :value="items" lazy :totalRecords="total" :first="firstRecord" :rows="perPage"
    @page="onPage" paginator paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
    :rowsPerPageOptions="[10,15,25,50]" size="small"
    class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
    <template #empty><slot name="empty" /></template>
    <Column v-for="col in columns" :key="col.field" :field="col.field" :header="col.header" sortable>
      <template #body="{data}">
        <span v-if="col.field.startsWith('_')" class="text-gray-500 dark:text-gray-400 text-xs">{{ data[col.field] || '-' }}</span>
        <span v-else class="text-gray-800 dark:text-gray-100">{{ data[col.field] || '-' }}</span>
      </template>
    </Column>
    <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
      <template #body="{data}">
        <div class="flex items-center gap-1">
          <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="$emit('edit', data)" />
          <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="$emit('delete', data)" />
        </div>
      </template>
    </Column>
  </DataTable>
</template>
<script setup>
import { computed, ref, onMounted } from 'vue'
import { useI18n } from '@/composables/useI18n'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import SkeletonTable from '@/components/SkeletonTable.vue'

const props = defineProps({
  items: Array, loading: Boolean, total: Number,
  columns: { type: Array, default: () => [] },
  entity: String, orgId: String, onLoad: Function
})
defineEmits(['edit', 'delete'])
const { t } = useI18n()
const currentPage = ref(1); const perPage = ref(15)
const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)
const skeletonCols = computed(() => [
  ...props.columns.map(c => ({ type: 'text', width: 'w-24', headerWidth: 'w-20' })),
  { type: 'icons', count: 2, headerWidth: 'w-16' }
])

function onPage(event) {
  currentPage.value = event.page + 1; perPage.value = event.rows
  if (props.onLoad) props.onLoad(currentPage.value, perPage.value)
}

onMounted(() => {
  if (props.onLoad) props.onLoad(1, 15)
})
</script>
