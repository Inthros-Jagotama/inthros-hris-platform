<template>
  <div class="max-w-2xl">
    <div v-if="loading" class="space-y-3">
      <div v-for="n in 7" :key="n" class="h-9 bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
    </div>
    <div v-else class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 space-y-4">
      <FormRow :label="t('attendance.latitude')">
        <InputNumber v-model="form.latitude" class="!w-full" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
      </FormRow>
      <FormRow :label="t('attendance.longitude')">
        <InputNumber v-model="form.longitude" class="!w-full" :minFractionDigits="0" :maxFractionDigits="8" size="small" />
      </FormRow>
      <FormRow :label="t('attendance.max_distance_meter')">
        <InputNumber v-model="form.max_distance_meter" class="!w-full" :min="0" size="small" suffix=" m" />
      </FormRow>
      <FormRow :label="t('attendance.late_tolerance_minutes')">
        <InputNumber v-model="form.late_tolerance_minutes" class="!w-full" :min="0" size="small" suffix=" min" />
      </FormRow>
      <FormRow :label="t('attendance.overtime_min_minutes')">
        <InputNumber v-model="form.overtime_min_minutes" class="!w-full" :min="0" size="small" suffix=" min" />
      </FormRow>
      <FormRow :label="t('attendance.is_location_required')">
        <ToggleSwitch v-model="form.is_location_required" />
      </FormRow>
      <FormRow :label="t('attendance.is_face_required')">
        <ToggleSwitch v-model="form.is_face_required" />
      </FormRow>
      <FormRow :label="t('attendance.is_overtime_enabled')">
        <ToggleSwitch v-model="form.is_overtime_enabled" />
      </FormRow>
      <div class="flex justify-end pt-2">
        <Button :label="t('common.save')" size="small" :loading="saving" :disabled="saving || !canUpdate" @click="handleSave" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { useAuth } from '@/stores/auth'
import { getErrorMessage } from '@/services/responseHandler'
import api from '@/services/api'

import InputNumber from 'primevue/inputnumber'
import ToggleSwitch from 'primevue/toggleswitch'
import Button from 'primevue/button'
import FormRow from '@/components/FormRow.vue'

const { t } = useI18n()
const toast = useToast()
const { hasPermission } = useAuth()

const canUpdate = computed(() => hasPermission('attendance.update'))
const loading = ref(true)
const saving = ref(false)
const form = ref({
  latitude: null,
  longitude: null,
  max_distance_meter: null,
  late_tolerance_minutes: null,
  overtime_min_minutes: null,
  is_location_required: false,
  is_face_required: false,
  is_overtime_enabled: false
})

async function loadData() {
  loading.value = true
  try {
    const res = await api.get('/api/v1/tenant/attendance/settings')
    const s = res.data?.data
    if (s) {
      form.value = {
        latitude: s.latitude ?? null,
        longitude: s.longitude ?? null,
        max_distance_meter: s.max_distance_meter ?? null,
        late_tolerance_minutes: s.late_tolerance_minutes ?? null,
        overtime_min_minutes: s.overtime_min_minutes ?? null,
        is_location_required: !!s.is_location_required,
        is_face_required: !!s.is_face_required,
        is_overtime_enabled: !!s.is_overtime_enabled
      }
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    await api.put('/api/v1/tenant/attendance/settings', form.value)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
