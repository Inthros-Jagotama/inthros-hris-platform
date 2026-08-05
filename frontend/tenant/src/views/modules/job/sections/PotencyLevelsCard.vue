<template>
  <div class="space-y-4 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl p-5">
    <div>
      <h3 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ t(titleKey) }}</h3>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t(descriptionKey) }}</p>
    </div>

    <SkeletonCard v-if="loading" type="detail" :count="1" :rows="skeletonRows" cols="grid-cols-1" padding="p-5" />

    <template v-else>
      <div
        v-if="rows.length === 0"
        class="text-sm text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-700/40 border border-gray-200 dark:border-gray-700 rounded-lg px-3 py-2"
      >
        {{ t(emptyKey) }}
      </div>

      <div v-else class="overflow-x-auto rounded-lg border border-gray-200 dark:border-gray-700">
        <table class="w-full text-sm">
          <thead>
            <tr class="bg-gray-50 dark:bg-gray-700/40 text-left text-xs uppercase tracking-wide text-gray-500 dark:text-gray-400">
              <th class="px-4 py-3 font-semibold min-w-[220px]">{{ t('job_management.potency_table_name') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_level') }}</th>
              <th class="px-4 py-3 font-semibold min-w-[260px]">{{ t('job_management.potency_table_description') }}</th>
              <th class="px-4 py-3 font-semibold w-16 text-right">{{ t('common.actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in rows"
              :key="row.type"
              class="border-t border-gray-100 dark:border-gray-700 align-top"
            >
              <td class="px-4 py-3">
                <div class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ row.competency_name }}</div>
                <div v-if="row.competency_definition" class="mt-0.5 text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
                  {{ row.competency_definition }}
                </div>
              </td>
              <td class="px-4 py-3">
                <SelectLabel
                  v-model="row.job_management_value_id"
                  :options="row.levelOptions"
                  option-label="label"
                  option-value="value"
                  :placeholder="t('common.select')"
                  showClear
                />
              </td>
              <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-300">
                {{ levelDescription(row) }}
              </td>
              <td class="px-4 py-3 text-right">
                <Button
                  icon="pi pi-trash"
                  severity="danger"
                  text
                  rounded
                  size="small"
                  :disabled="savingCard"
                  :aria-label="t('common.delete')"
                  @click="askDeleteRow(row)"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="errorMsg" class="text-sm text-red-500 dark:text-red-400 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg px-3 py-2">
        {{ errorMsg }}
      </div>

      <div v-if="rows.length > 0" class="flex justify-end gap-2 pt-1">
        <Button
          :label="t(saveLabelKey)"
          icon="pi pi-check"
          size="small"
          :loading="savingCard"
          :disabled="savingCard || !orgId"
          @click="handleSave"
        />
      </div>
    </template>

    <ConfirmDeleteDialog
      v-model:visible="deleteVisible"
      :title="t(deleteTitleKey)"
      :message="t(deleteMessageKey, { name: deleteTarget?.competency_name || '' })"
      :loading="deleting"
      :error-msg="deleteError"
      @confirm="handleDelete"
      @cancel="deleteVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { usePotencyLevels } from '@/composables/usePotencyLevels'
import Button from 'primevue/button'
import SelectLabel from '@/components/SelectLabel.vue'
import SkeletonCard from '@/components/SkeletonCard.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const emit = defineEmits(['saved'])

const props = defineProps({
  orgId: String,
  rows: { type: Array, default: () => [] },
  // Parent menyalakan setelah options (level) selesai dimuat → trigger loadData + hydrate internal
  optionsReady: { type: Boolean, default: false },
  skeletonRows: { type: Number, default: 2 },
  titleKey: { type: String, required: true },
  descriptionKey: { type: String, required: true },
  emptyKey: { type: String, required: true },
  saveLabelKey: { type: String, required: true },
  deleteTitleKey: { type: String, required: true },
  deleteMessageKey: { type: String, required: true }
})

const { t } = useI18n()

const loading = ref(true)

// rows sebagai computed → composable (usePotencyLevels) membaca rows.value terkini
const rowsRef = computed(() => props.rows)

// =========================================================================
// Logika bersama simpan/hapus/hydrate (composable usePotencyLevels)
// =========================================================================
const {
  savingCard, errorMsg,
  deleteVisible, deleting, deleteError, deleteTarget,
  levelDescription, hydrateRows,
  loadData, askDeleteRow, handleDelete, handleSave
} = usePotencyLevels({
  orgId: computed(() => props.orgId),
  rows: rowsRef,
  // Delete: cukup hapus record — hydrateRows otomatis mereset level baris
  onSaved: () => emit('saved')
})

// Hydrate level tersimpan setelah options siap (sekali saja)
let hydrated = false
watch(() => props.optionsReady, async (ready) => {
  if (!ready || hydrated) return
  hydrated = true
  try {
    await loadData()
  } finally {
    hydrateRows()
    loading.value = false
  }
}, { immediate: true })
</script>
