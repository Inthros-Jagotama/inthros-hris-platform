<template>
  <div class="space-y-4">
    <!-- Profile Picture -->
    <div class="flex items-start gap-6 mb-4">
      <div class="relative shrink-0">
        <div class="w-24 h-24 rounded-full overflow-hidden border-2 border-gray-200 dark:border-gray-600 bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
          <img v-if="previewUrl" :src="previewUrl" alt="Profile" class="w-full h-full object-cover" />
          <i v-else class="pi pi-user text-3xl text-gray-400 dark:text-gray-500"></i>
        </div>
        <button type="button"
          class="absolute -bottom-1 -right-1 w-8 h-8 rounded-full bg-emerald-500 hover:bg-emerald-600 text-white shadow flex items-center justify-center transition-colors duration-150"
          @click="$refs.fileInput.click()" v-tooltip.top="t('employee.change_photo')">
          <i class="pi pi-camera text-sm"></i>
        </button>
        <button v-if="previewUrl" type="button"
          class="absolute -top-1 -right-1 w-6 h-6 rounded-full bg-red-500 hover:bg-red-600 disabled:bg-red-300 text-white shadow flex items-center justify-center transition-colors duration-150"
          :disabled="removing" @click="confirmRemove" v-tooltip.top="t('employee.remove_photo')">
          <i :class="removing ? 'pi pi-spin pi-spinner text-xs' : 'pi pi-times text-xs'"></i>
        </button>
      </div>
      <div class="flex-1 min-w-0 pt-2">
        <div class="text-sm font-medium text-gray-700 dark:text-gray-300">{{ t('employee.profile_photo') }}</div>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-0.5">{{ t('employee.photo_hint') }}</p>
        <div v-if="uploading" class="flex items-center gap-2 mt-2">
          <i class="pi pi-spin pi-spinner text-emerald-500 text-sm"></i>
          <span class="text-xs text-gray-400">{{ t('employee.uploading') }}</span>
        </div>
        <div v-else-if="photoUrl && !previewUrl" class="flex items-center gap-2 mt-2">
          <i class="pi pi-check-circle text-emerald-500 text-sm"></i>
          <span class="text-xs text-emerald-600 dark:text-emerald-400">{{ t('employee.photo_saved') }}</span>
        </div>
      </div>
      <input ref="fileInput" type="file" accept="image/jpeg,image/png,image/gif,image/webp" class="hidden" @change="onFileSelect" />
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <FormRow :label="t('employee.employee_id')" required :errors="errors?.employee_id">
        <TextInput v-model="form.employee_id" maxlength="50" :placeholder="t('employee.employee_id_placeholder')" :class="{'p-invalid':errors?.employee_id}" :disabled="disabled" />
      </FormRow>
      <FormRow :label="t('common.name')" required :errors="errors?.name">
        <TextInput v-model="form.name" maxlength="255" :placeholder="t('common.name')" :class="{'p-invalid':errors?.name}" />
      </FormRow>
      <FormRow :label="t('employee.nik')" :errors="errors?.nik">
        <TextInput v-model="form.nik" maxlength="16" :placeholder="t('employee.nik_placeholder')" :class="{'p-invalid':errors?.nik}" />
      </FormRow>
      <FormRow :label="t('employee.gender')" :errors="errors?.gender">
        <SelectLabel v-model="form.gender" :options="genderOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_gender')" :class="{'p-invalid':errors?.gender}" :showClear="true" />
      </FormRow>
      <FormRow :label="t('employee.mother_name')" :errors="errors?.mother_name">
        <TextInput v-model="form.mother_name" maxlength="255" :placeholder="t('employee.mother_name_placeholder')" :class="{'p-invalid':errors?.mother_name}" />
      </FormRow>
      <FormRow :label="t('employee.pob')" :errors="errors?.pob">
        <TextInput v-model="form.pob" maxlength="255" :placeholder="t('employee.pob_placeholder')" :class="{'p-invalid':errors?.pob}" />
      </FormRow>
      <FormRow :label="t('employee.dob')" :errors="errors?.dob">
        <DateInput v-model="form.dob" :placeholder="t('employee.dob_placeholder')" :class="{'p-invalid':errors?.dob}" />
      </FormRow>
      <FormRow :label="t('employee.phone')" :errors="errors?.phone_number">
        <TextInput v-model="form.phone_number" maxlength="255" :placeholder="t('employee.phone_placeholder')" :class="{'p-invalid':errors?.phone_number}" />
      </FormRow>
      <FormRow :label="t('employee.email')" :errors="errors?.email">
        <TextInput v-model="form.email" maxlength="255" :placeholder="t('employee.email_placeholder')" :class="{'p-invalid':errors?.email}" />
      </FormRow>
      <FormRow :label="t('employee.religion')" :errors="errors?.religion_id">
        <SelectLabel v-model="form.religion_id" :options="religionOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_religion')" :class="{'p-invalid':errors?.religion_id}" :showClear="true" />
      </FormRow>
      <FormRow :label="t('employee.marital_status')" :errors="errors?.marital_status_id">
        <SelectLabel v-model="form.marital_status_id" :options="maritalStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('employee.select_marital_status')" :class="{'p-invalid':errors?.marital_status_id}" :showClear="true" />
      </FormRow>
    </div>
    <div class="flex justify-end pt-2">
      <Button :label="t('employee.save_personal')" icon="pi pi-check" size="small" :loading="saving" :disabled="saving" @click="$emit('save')" />
    </div>

    <!-- Crop Dialog -->
    <ImageCropDialog
      v-model:visible="cropDialogVisible"
      :image-src="cropImageSrc"
      :title="t('employee.crop_photo')"
      @crop="onCropComplete"
      @cancel="onCropCancel"
    />

    <!-- Delete Photo Confirmation -->
    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :loading="removing"
      :error-msg="deleteError"
      :title="t('employee.remove_photo')"
      :message="t('employee.confirm_remove_photo')"
      :confirm-label="t('common.delete')"
      @confirm="removePhoto"
      @cancel="deleteDialogVisible = false"
    />
  </div>
</template>
<script setup>
import { ref, watch, shallowRef } from 'vue'
import { useI18n } from '@/composables/useI18n'
import Button from 'primevue/button'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import ImageCropDialog from '@/components/ImageCropDialog.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const { t } = useI18n()
const toast = useToast()
const emit = defineEmits(['save', 'update:photo'])
const removing = ref(false)
const deleteDialogVisible = ref(false)
const deleteError = ref('')

const props = defineProps({
  form: { type: Object, required: true },
  errors: { type: Object, default: () => ({}) },
  genderOptions: { type: Array, default: () => [] },
  religionOptions: { type: Array, default: () => [] },
  maritalStatusOptions: { type: Array, default: () => [] },
  saving: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  employeeId: { type: String, default: '' },
  photoUrl: { type: String, default: '' }
})

const fileInput = ref(null)
const previewUrl = ref('')
const uploading = ref(false)
const cropDialogVisible = ref(false)
const cropImageSrc = shallowRef('')

// Watch for photoUrl prop changes to set preview
watch(() => props.photoUrl, (url) => {
  if (url && !previewUrl.value) {
    previewUrl.value = url
  }
}, { immediate: true })

async function onFileSelect(e) {
  const file = e.target?.files?.[0]
  if (!file) return

  // Validate file size (2MB max)
  if (file.size > 2 * 1024 * 1024) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('employee.photo_too_large'), life: 4000 })
    return
  }

  // Validate file type
  const validTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp']
  if (!validTypes.includes(file.type)) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('employee.photo_invalid_type'), life: 4000 })
    return
  }

  // Read file as data URL and show crop dialog
  const reader = new FileReader()
  reader.onload = (ev) => {
    cropImageSrc.value = ev.target?.result || ''
    cropDialogVisible.value = true
  }
  reader.readAsDataURL(file)
}

function onCropCancel() {
  cropDialogVisible.value = false
  cropImageSrc.value = ''
}

async function onCropComplete(blob) {
  if (!props.employeeId) {
    toast.add({ severity: 'warn', summary: t('employee.personal_first'), life: 3000 })
    return
  }

  // Create a temporary preview from the blob
  const previewReader = new FileReader()
  previewReader.onload = (ev) => { previewUrl.value = ev.target?.result || '' }
  previewReader.readAsDataURL(blob)

  uploading.value = true
  try {
    const formData = new FormData()
    // Create a File from blob with a proper name
    const croppedFile = new File([blob], 'profile.jpg', { type: 'image/jpeg' })
    formData.append('photo', croppedFile)
    const res = await api.put(`/api/v1/tenant/employees/${props.employeeId}/photo`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const photoPath = res.data?.data?.profile_picture
    if (photoPath) {
      // Add cache-busting timestamp so browser doesn't show stale cached image
      const cacheBuster = `?t=${Date.now()}`
      previewUrl.value = `${photoPath}${cacheBuster}`
      emit('update:photo', `${photoPath}${cacheBuster}`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.photo_uploaded'), life: 2000 })
    }
  } catch (e) {
    const msg = e.response?.data?.error?.message || t('message.operation_failed')
    toast.add({ severity: 'error', summary: t('message.error'), detail: msg, life: 4000 })
    previewUrl.value = ''
    // Reset to original photo if exists
    if (props.photoUrl) {
      previewUrl.value = props.photoUrl
    }
  } finally {
    uploading.value = false
  }
}

function confirmRemove() {
  deleteError.value = ''
  deleteDialogVisible.value = true
}

async function removePhoto() {
  if (!props.employeeId) return
  removing.value = true
  deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/employees/${props.employeeId}/photo`)
    previewUrl.value = ''
    emit('update:photo', '')
    deleteDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('employee.photo_removed'), life: 2000 })
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    removing.value = false
  }
}
</script>
