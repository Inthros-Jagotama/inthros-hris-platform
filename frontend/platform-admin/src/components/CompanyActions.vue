<template>
  <div class="inline-flex">
    <!-- Icons mode (compact — untuk row actions di DataTable) -->
    <template v-if="mode === 'icons'">
      <div class="flex items-center gap-1">
        <Button v-if="canEdit" icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEdit" />
        <Button v-if="company.status === 'active'" icon="pi pi-pause-circle" size="small" text severity="warning" v-tooltip.left="t('companies.action_suspend')" @click="confirmAction('suspend')" />
        <Button v-if="company.status === 'suspended'" icon="pi pi-play-circle" size="small" text severity="info" v-tooltip.left="t('companies.action_activate')" @click="confirmAction('activate')" />
        <Button v-if="canRotate" icon="pi pi-key" size="small" text severity="info" v-tooltip.left="t('companies.action_rotate')" @click="confirmAction('rotate')" />
        <Button v-if="company.status !== 'terminated'" icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('companies.action_terminate')" @click="confirmAction('terminate')" />
      </div>
    </template>

    <!-- Buttons mode (berlabel — untuk header halaman detail) -->
    <template v-else>
      <div class="flex items-center gap-2 flex-wrap">
        <Button v-if="canEdit" :label="t('common.edit')" icon="pi pi-pencil" severity="secondary" outlined size="small" @click="openEdit" />
        <Button v-if="company.status === 'active'" :label="t('companies.action_suspend')" icon="pi pi-pause-circle" severity="warning" outlined size="small" @click="confirmAction('suspend')" />
        <Button v-if="company.status === 'suspended'" :label="t('companies.action_activate')" icon="pi pi-play-circle" severity="info" outlined size="small" @click="confirmAction('activate')" />
        <Button v-if="canRotate" :label="t('companies.action_rotate')" icon="pi pi-key" severity="info" outlined size="small" @click="confirmAction('rotate')" />
        <Button v-if="company.status !== 'terminated'" :label="t('companies.action_terminate')" icon="pi pi-trash" severity="danger" outlined size="small" @click="confirmAction('terminate')" />
      </div>
    </template>

    <!-- Edit Dialog -->
    <Dialog v-model:visible="editDialogVisible" :header="t('companies.edit_company')" modal :style="{ width: '620px' }" :closable="true">
      <div class="space-y-4">
        <FormRow :label="t('companies.company_name')" :errors="errors?.name" :required="true">
          <TextInput v-model="form.name" autofocus :class="{ 'p-invalid': errors?.name }" />
        </FormRow>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('companies.subdomain')" :errors="errors?.subdomain">
            <TextInput v-model="form.subdomain" :placeholder="t('companies.subdomain_placeholder')" :class="{ 'p-invalid': errors?.subdomain }" />
          </FormRow>
          <FormRow :label="t('companies.domain')" :errors="errors?.domain">
            <TextInput v-model="form.domain" :placeholder="t('companies.domain_placeholder')" :class="{ 'p-invalid': errors?.domain }" />
          </FormRow>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('companies.email')" :errors="errors?.email">
            <TextInput v-model="form.email" :class="{ 'p-invalid': errors?.email }" />
          </FormRow>
          <FormRow :label="t('companies.phone')" :errors="errors?.phone">
            <TextInput v-model="form.phone" :class="{ 'p-invalid': errors?.phone }" />
          </FormRow>
        </div>
        <FormRow :label="t('companies.address')" :errors="errors?.address">
          <TextInput v-model="form.address" :class="{ 'p-invalid': errors?.address }" />
        </FormRow>

        <!-- Current License Info -->
        <div v-if="company.license_info">
          <div class="border-t border-gray-200 dark:border-gray-700 my-3"></div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2 flex items-center gap-1.5">
            <i class="pi pi-id-card text-indigo-400 text-sm"></i>
            {{ t('companies.license_section') }}
          </h3>
          <div class="flex items-center gap-3 text-sm">
            <Tag :value="company.license_info.plan_type" :severity="planSeverity(company.license_info.plan_type)" class="!text-xs" />
            <span class="text-gray-500 dark:text-gray-400">{{ t('companies.license_key_label') }}: {{ company.license_info.license_key?.substring(0, 12) }}...</span>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="editDialogVisible = false" />
          <Button :label="t('common.update')" size="small" :loading="saving" :disabled="saving" @click="saveCompany" />
        </div>
      </template>
    </Dialog>

    <!-- Confirm Dialog (suspend/activate/terminate/rotate) -->
    <Dialog v-model:visible="confirmVisible" :header="confirmTitle" modal :style="{ width: '400px' }">
      <p class="text-xs text-gray-600 dark:text-gray-300">{{ confirmMessage }}</p>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" text size="small" @click="confirmVisible = false" />
        <Button :label="confirmActionLabel" :severity="confirmSeverity" size="small" :loading="confirming" :disabled="confirming" @click="executeConfirm" />
      </template>
    </Dialog>

    <!-- Rotate Credentials: show auto-generated password once -->
    <Dialog v-model:visible="passwordDialogVisible" :header="t('companies.rotate_password_title')" modal :style="{ width: '520px' }" :closable="true">
      <div class="space-y-3">
        <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('companies.rotate_password_message') }}</p>
        <div class="flex items-center gap-2">
          <InputText :model-value="rotatedPassword" readonly class="!w-full !font-mono !text-sm" />
          <Button icon="pi pi-copy" severity="secondary" outlined size="small" v-tooltip.top="t('companies.copy_password')" @click="copyRotatedPassword" />
        </div>
        <div class="flex items-start gap-2 bg-amber-50 dark:bg-amber-500/10 border border-amber-200 dark:border-amber-500/30 rounded p-2.5 text-xs text-amber-700 dark:text-amber-400">
          <i class="pi pi-exclamation-triangle mt-0.5"></i>
          <span>{{ t('companies.rotate_password_warning') }}</span>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="passwordDialogVisible = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import api from '@/services/api'
import { getValidationErrors } from '@/services/responseHandler'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'

const props = defineProps({
  company: { type: Object, required: true },
  // 'icons' = tombol icon saja (untuk row di DataTable); 'buttons' = tombol berlabel (untuk header detail)
  mode: { type: String, default: 'icons' }
})

const emit = defineEmits(['updated'])

const toast = useToast()
const { t } = useI18n()

const confirmVisible = ref(false)
const confirmActionType = ref(null)
const confirming = ref(false)
const passwordDialogVisible = ref(false)
const rotatedPassword = ref('')
const editDialogVisible = ref(false)
const saving = ref(false)
const form = ref({ name: '', subdomain: '', domain: '', email: '', phone: '', address: '' })
const errors = ref({})

const canEdit = computed(() => props.company.status !== 'terminated')
const canRotate = computed(() => props.company.provisioning_info?.provisioned && props.company.status !== 'terminated')

const confirmTitle = computed(() => {
  switch (confirmActionType.value) {
    case 'suspend': return t('companies.confirm_suspend_title')
    case 'activate': return t('companies.confirm_activate_title')
    case 'terminate': return t('companies.confirm_terminate_title')
    case 'rotate': return t('companies.confirm_rotate_title')
    default: return ''
  }
})

const confirmMessage = computed(() => {
  const name = props.company.name
  switch (confirmActionType.value) {
    case 'suspend': return t('companies.confirm_suspend_message', { name })
    case 'activate': return t('companies.confirm_activate_message', { name })
    case 'terminate': return t('companies.confirm_terminate_message', { name })
    case 'rotate': return t('companies.confirm_rotate_message', { name })
    default: return ''
  }
})

const confirmActionLabel = computed(() => {
  switch (confirmActionType.value) {
    case 'suspend': return t('companies.action_suspend')
    case 'activate': return t('companies.action_activate')
    case 'terminate': return t('companies.action_terminate')
    case 'rotate': return t('companies.action_rotate')
    default: return ''
  }
})

const confirmSeverity = computed(() => {
  return confirmActionType.value === 'terminate' ? 'danger' : 'warn'
})

function planSeverity(plan) {
  switch (plan?.toLowerCase()) {
    case 'enterprise': return 'danger'
    case 'professional': return 'warn'
    case 'basic': return 'info'
    case 'subscription': return 'info'
    case 'trial': return 'success'
    default: return 'info'
  }
}

function confirmAction(action) {
  confirmActionType.value = action
  confirmVisible.value = true
}

async function executeConfirm() {
  if (!confirmActionType.value) return
  confirming.value = true
  const id = props.company.id
  try {
    if (confirmActionType.value === 'rotate') {
      const res = await api.post(`/api/v1/platform/companies/${id}/rotate-credentials`, {})
      const newPassword = res.data?.data?.new_password
      confirmVisible.value = false
      if (newPassword) {
        rotatedPassword.value = newPassword
        passwordDialogVisible.value = true
      }
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('companies.rotate_success'), life: 3000 })
    } else {
      await api.post(`/api/v1/platform/companies/${id}/${confirmActionType.value}`)
      toast.add({ severity: 'success', summary: t('message.success'), detail: `${t('companies.title')} ${confirmActionType.value}ed`, life: 2000 })
      confirmVisible.value = false
    }
    emit('updated')
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
  } finally {
    confirming.value = false
  }
}

async function copyRotatedPassword() {
  if (!rotatedPassword.value) return
  try {
    await navigator.clipboard.writeText(rotatedPassword.value)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('companies.password_copied'), life: 2000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: t('message.operation_failed'), life: 2000 })
  }
}

function openEdit() {
  form.value = {
    name: props.company.name || '',
    subdomain: props.company.subdomain || '',
    domain: props.company.domain || '',
    email: props.company.email || '',
    phone: props.company.phone || '',
    address: props.company.address || ''
  }
  errors.value = {}
  editDialogVisible.value = true
}

async function saveCompany() {
  saving.value = true
  try {
    await api.put(`/api/v1/platform/companies/${props.company.id}`, form.value)
    toast.add({ severity: 'success', summary: t('message.updated'), detail: t('companies.title'), life: 2000 })
    editDialogVisible.value = false
    emit('updated')
  } catch (e) {
    errors.value = getValidationErrors(e)
    if (Object.keys(errors.value).length === 0) {
      toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 3000 })
    }
  } finally {
    saving.value = false
  }
}
</script>
