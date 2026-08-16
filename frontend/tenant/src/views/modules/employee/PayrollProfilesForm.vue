<template>
  <div class="space-y-6">
    <!-- Loading -->
    <div v-if="loading" class="space-y-4">
      <div v-for="n in 4" :key="n" class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
        <div class="h-4 w-48 bg-gray-200 dark:bg-gray-700 rounded animate-pulse mb-4"></div>
        <div class="h-9 w-full bg-gray-200 dark:bg-gray-700 rounded animate-pulse"></div>
      </div>
    </div>

    <template v-else>
      <!-- ═══════════════════ Payroll Profile ═══════════════════ -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <div class="flex items-center justify-between gap-2 flex-wrap px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              <i class="pi pi-users text-sky-500 text-sm"></i> {{ t('payroll.employee_profiles') }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('payroll.employee_profiles_desc') }}</p>
          </div>
          <Button :label="t('payroll.new_profile')" icon="pi pi-plus" size="small" @click="openPayrollDialog()" />
        </div>
        <div class="p-4">
          <DataTable :value="payrollProfiles" size="small" class="!text-sm p-datatable-sm" :rows="10" :paginator="payrollProfiles.length > 10">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-users text-2xl mb-2 opacity-50"></i>
                <p class="text-sm font-medium">{{ t('payroll.profiles_empty') }}</p>
              </div>
            </template>
            <Column field="payroll_group_code" :header="t('payroll.payroll_group_code')" style="width:180px"><template #body="{data}"><Tag :value="data.payroll_group_code" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column field="payroll_frequency" :header="t('payroll.payroll_frequency')" style="width:120px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ freqLabel(data.payroll_frequency) }}</span></template></Column>
            <Column field="payment_method" :header="t('payroll.payment_method')" style="width:140px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ methodLabel(data.payment_method) }}</span></template></Column>
            <Column field="salary_currency" :header="t('payroll.salary_currency')" style="width:90px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.salary_currency || '-' }}</span></template></Column>
            <Column field="is_payroll_active" :header="t('payroll.is_payroll_active')" style="width:110px"><template #body="{data}"><Tag :value="data.is_payroll_active ? t('common.yes') : t('common.no')" :severity="data.is_payroll_active ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column field="effective_start_date" :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
            <Column :header="t('common.actions')" style="width:70px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center justify-end"><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete('payroll', data)" /></div></template></Column>
          </DataTable>
        </div>
      </div>

      <!-- ═══════════════════ Bank Profile ═══════════════════ -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <div class="flex items-center justify-between gap-2 flex-wrap px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              <i class="pi pi-building text-sky-500 text-sm"></i> {{ t('payroll.bank_profiles') }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('payroll.bank_profiles_desc') }}</p>
          </div>
          <Button :label="t('payroll.new_bank_profile')" icon="pi pi-plus" size="small" @click="openBankDialog()" />
        </div>
        <div class="p-4">
          <p v-if="payrollProfiles.length === 0" class="text-xs text-amber-500 dark:text-amber-400 mb-3 flex items-center gap-1.5"><i class="pi pi-info-circle"></i> {{ t('payroll.profile_first_hint') }}</p>
          <DataTable :value="bankProfiles" size="small" class="!text-sm p-datatable-sm" :rows="10" :paginator="bankProfiles.length > 10">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-building text-2xl mb-2 opacity-50"></i>
                <p class="text-sm font-medium">{{ t('payroll.bank_profiles_empty') }}</p>
              </div>
            </template>
            <Column :header="t('payroll.bank_name')" style="width:130px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.bank_name }}</span></template></Column>
            <Column :header="t('payroll.bank_account_number')" style="width:160px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ data.bank_account_number }}</span></template></Column>
            <Column :header="t('payroll.bank_account_holder_name')" style="width:180px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.bank_account_holder_name }}</span></template></Column>
            <Column :header="t('payroll.is_primary')" style="width:90px"><template #body="{data}"><Tag :value="data.is_primary ? t('common.yes') : t('common.no')" :severity="data.is_primary ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
            <Column :header="t('common.status')" style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openBankDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete('bank', data)" /></div></template></Column>
          </DataTable>
        </div>
      </div>

      <!-- ═══════════════════ BPJS Profile ═══════════════════ -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <div class="flex items-center justify-between gap-2 flex-wrap px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              <i class="pi pi-shield text-teal-500 text-sm"></i> {{ t('payroll.bpjs_profiles') }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('payroll.bpjs_profiles_desc') }}</p>
          </div>
          <Button :label="t('payroll.new_bpjs_profile')" icon="pi pi-plus" size="small" @click="openBpjsDialog()" />
        </div>
        <div class="p-4">
          <p v-if="payrollProfiles.length === 0" class="text-xs text-amber-500 dark:text-amber-400 mb-3 flex items-center gap-1.5"><i class="pi pi-info-circle"></i> {{ t('payroll.profile_first_hint') }}</p>
          <DataTable :value="bpjsProfiles" size="small" class="!text-sm p-datatable-sm" :rows="10" :paginator="bpjsProfiles.length > 10">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-shield text-2xl mb-2 opacity-50"></i>
                <p class="text-sm font-medium">{{ t('payroll.bpjs_profiles_empty') }}</p>
              </div>
            </template>
            <Column :header="t('payroll.bpjs_health_no')" style="width:150px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ data.bpjs_health_no || '-' }}</span></template></Column>
            <Column :header="t('payroll.bpjs_tk_no')" style="width:150px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ data.bpjs_tk_no || '-' }}</span></template></Column>
            <Column :header="t('payroll.jkk_risk_class')" style="width:120px"><template #body="{data}"><Tag :value="riskLabel(data.jkk_risk_class)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
            <Column :header="t('common.status')" style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openBpjsDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete('bpjs', data)" /></div></template></Column>
          </DataTable>
        </div>
      </div>

      <!-- ═══════════════════ Tax Profile ═══════════════════ -->
      <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
        <div class="flex items-center justify-between gap-2 flex-wrap px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
          <div>
            <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100 flex items-center gap-2">
              <i class="pi pi-receipt text-purple-500 text-sm"></i> {{ t('payroll.tax_profiles') }}
            </h3>
            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ t('payroll.tax_profiles_desc') }}</p>
          </div>
          <Button :label="t('payroll.new_tax_profile')" icon="pi pi-plus" size="small" @click="openTaxDialog()" />
        </div>
        <div class="p-4">
          <p v-if="payrollProfiles.length === 0" class="text-xs text-amber-500 dark:text-amber-400 mb-3 flex items-center gap-1.5"><i class="pi pi-info-circle"></i> {{ t('payroll.profile_first_hint') }}</p>
          <DataTable :value="taxProfiles" size="small" class="!text-sm p-datatable-sm" :rows="10" :paginator="taxProfiles.length > 10">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-receipt text-2xl mb-2 opacity-50"></i>
                <p class="text-sm font-medium">{{ t('payroll.tax_profiles_empty') }}</p>
              </div>
            </template>
            <Column :header="t('payroll.npwp')" style="width:160px"><template #body="{data}"><span class="font-mono text-gray-600 dark:text-gray-300 text-xs">{{ data.npwp || '-' }}</span></template></Column>
            <Column :header="t('payroll.ptkp_status')" style="width:100px"><template #body="{data}"><Tag :value="data.ptkp_status || '-'" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('payroll.tax_method')" style="width:110px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ taxMethodLabel(data.tax_method) }}</span></template></Column>
            <Column :header="t('payroll.has_npwp')" style="width:90px"><template #body="{data}"><Tag :value="data.has_npwp ? t('common.yes') : t('common.no')" :severity="data.has_npwp ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('payroll.effective_start_date')" style="width:130px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
            <Column :header="t('common.status')" style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openTaxDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDelete('tax', data)" /></div></template></Column>
          </DataTable>
        </div>
      </div>
    </template>

    <!-- ── Payroll Profile Dialog ── -->
    <Dialog v-model:visible="payrollDialogVisible" :header="t('payroll.new_profile')" modal :style="{ width: 'min(680px, 95vw)' }" :closable="true" @hide="resetPayrollForm">
      <div class="space-y-4">
        <!-- Frekuensi — card + radio, 3 kolom -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.payroll_frequency') }}</span>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <div v-for="opt in freqOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-2"
                 :class="payrollForm.payroll_frequency === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="payrollForm.payroll_frequency = opt.value">
              <RadioButton :modelValue="payrollForm.payroll_frequency" :inputId="'pf-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="payrollForm.payroll_frequency = $event" class="mt-0.5 shrink-0" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'pf-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500 leading-snug">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <!-- Metode Pembayaran — card + radio, 3 kolom -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.payment_method') }}</span>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <div v-for="opt in methodOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-2"
                 :class="payrollForm.payment_method === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="payrollForm.payment_method = opt.value">
              <RadioButton :modelValue="payrollForm.payment_method" :inputId="'pm-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="payrollForm.payment_method = $event" class="mt-0.5 shrink-0" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'pm-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500 leading-snug">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.salary_currency')">
            <TextInput v-model="payrollForm.salary_currency" maxlength="3" placeholder="IDR" />
          </FormRow>
          <FormRow :label="t('payroll.is_payroll_active')">
            <ToggleSwitch v-model="payrollForm.is_payroll_active" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="payrollErrors?.effective_start_date">
            <DateInput v-model="payrollForm.effective_start_date" :class="{ 'p-invalid': payrollErrors?.effective_start_date }" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="payrollForm.effective_end_date" />
          </FormRow>
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="payrollDialogVisible = false" /><Button :label="t('common.save')" size="small" :loading="payrollSaving" :disabled="payrollSaving" @click="savePayrollProfile" /></div></template>
    </Dialog>

    <!-- ── Bank Profile Dialog ── -->
    <Dialog v-model:visible="bankDialogVisible" :header="bankEditing ? t('payroll.edit_bank_profile') : t('payroll.new_bank_profile')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetBankForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.payroll_profile')" required :errors="bankErrors?.employee_payroll_profile_id">
          <SelectLabel v-model="bankForm.employee_payroll_profile_id" :options="payrollProfileOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': bankErrors?.employee_payroll_profile_id }" />
        </FormRow>
        <FormRow :label="t('payroll.bank_name')" required :errors="bankErrors?.bank_id">
          <SelectLabel v-model="bankForm.bank_id" :options="bankOptions" optionLabel="label" optionValue="value" filter showClear :placeholder="t('common.select')" :class="{ 'p-invalid': bankErrors?.bank_id }" @update:modelValue="onBankChange" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.bank_account_number')" required :errors="bankErrors?.bank_account_number">
            <TextInput v-model="bankForm.bank_account_number" maxlength="100" :class="{ 'p-invalid': bankErrors?.bank_account_number }" />
          </FormRow>
          <FormRow :label="t('payroll.bank_account_holder_name')" required :errors="bankErrors?.bank_account_holder_name">
            <TextInput v-model="bankForm.bank_account_holder_name" maxlength="255" :class="{ 'p-invalid': bankErrors?.bank_account_holder_name }" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.bank_branch')">
            <TextInput v-model="bankForm.bank_branch" maxlength="150" />
          </FormRow>
          <FormRow :label="t('payroll.is_primary')">
            <ToggleSwitch v-model="bankForm.is_primary" />
          </FormRow>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="bankErrors?.effective_start_date">
            <DateInput v-model="bankForm.effective_start_date" :class="{ 'p-invalid': bankErrors?.effective_start_date }" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="bankForm.effective_end_date" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="bankForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="bankDialogVisible = false" /><Button :label="bankEditing ? t('common.update') : t('common.save')" size="small" :loading="bankSaving" :disabled="bankSaving" @click="saveBankProfile" /></div></template>
    </Dialog>

    <!-- ── BPJS Profile Dialog ── -->
    <Dialog v-model:visible="bpjsDialogVisible" :header="bpjsEditing ? t('payroll.edit_bpjs_profile') : t('payroll.new_bpjs_profile')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetBpjsForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.payroll_profile')" required :errors="bpjsErrors?.employee_payroll_profile_id">
          <SelectLabel v-model="bpjsForm.employee_payroll_profile_id" :options="payrollProfileOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': bpjsErrors?.employee_payroll_profile_id }" />
        </FormRow>
        <!-- Card BPJS Kesehatan -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
          <div class="px-4 py-3 flex items-center justify-between gap-3 bg-gray-50 dark:bg-gray-800/50">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.bpjs_health_active') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.bpjs_health_active_desc') }}</span>
            </div>
            <ToggleSwitch v-model="bpjsForm.bpjs_health_active" class="shrink-0" />
          </div>
          <div v-if="bpjsForm.bpjs_health_active" class="p-4 grid grid-cols-1 sm:grid-cols-2 gap-4 border-t border-gray-200 dark:border-gray-700">
            <FormRow :label="t('payroll.bpjs_health_no')">
              <TextInput v-model="bpjsForm.bpjs_health_no" maxlength="50" />
            </FormRow>
            <FormRow :label="t('payroll.bpjs_health_registered_name')">
              <TextInput v-model="bpjsForm.bpjs_health_registered_name" maxlength="255" />
            </FormRow>
          </div>
        </div>
        <!-- Card BPJS Ketenagakerjaan -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
          <div class="px-4 py-3 flex items-center justify-between gap-3 bg-gray-50 dark:bg-gray-800/50">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.bpjs_tk_active') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.bpjs_tk_active_desc') }}</span>
            </div>
            <ToggleSwitch v-model="bpjsForm.bpjs_tk_active" class="shrink-0" />
          </div>
          <div v-if="bpjsForm.bpjs_tk_active" class="p-4 space-y-4 border-t border-gray-200 dark:border-gray-700">
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <FormRow :label="t('payroll.bpjs_tk_no')">
                <TextInput v-model="bpjsForm.bpjs_tk_no" maxlength="50" />
              </FormRow>
              <FormRow :label="t('payroll.bpjs_tk_registered_name')">
                <TextInput v-model="bpjsForm.bpjs_tk_registered_name" maxlength="255" />
              </FormRow>
            </div>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <FormRow :label="t('payroll.jkk_risk_class')">
                <SelectLabel v-model="bpjsForm.jkk_risk_class" :options="riskOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
              </FormRow>
              <FormRow :label="t('payroll.pension_active')">
                <ToggleSwitch v-model="bpjsForm.pension_active" />
              </FormRow>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="bpjsErrors?.effective_start_date">
            <DateInput v-model="bpjsForm.effective_start_date" :class="{ 'p-invalid': bpjsErrors?.effective_start_date }" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="bpjsForm.effective_end_date" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="bpjsForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="bpjsDialogVisible = false" /><Button :label="bpjsEditing ? t('common.update') : t('common.save')" size="small" :loading="bpjsSaving" :disabled="bpjsSaving" @click="saveBpjsProfile" /></div></template>
    </Dialog>

    <!-- ── Tax Profile Dialog ── -->
    <Dialog v-model:visible="taxDialogVisible" :header="taxEditing ? t('payroll.edit_tax_profile') : t('payroll.new_tax_profile')" modal :style="{ width: 'min(700px, 95vw)' }" :closable="true" @hide="resetTaxForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.payroll_profile')" required :errors="taxErrors?.employee_payroll_profile_id">
          <SelectLabel v-model="taxForm.employee_payroll_profile_id" :options="payrollProfileOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': taxErrors?.employee_payroll_profile_id }" />
        </FormRow>
        <!-- Card NPWP -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
          <div class="px-4 py-3 flex items-center justify-between gap-3 bg-gray-50 dark:bg-gray-800/50">
            <div class="flex flex-col gap-0.5 min-w-0">
              <span class="text-sm font-medium text-surface-700 dark:text-surface-0/80">{{ t('payroll.has_npwp') }}</span>
              <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.has_npwp_desc') }}</span>
            </div>
            <ToggleSwitch v-model="taxForm.has_npwp" class="shrink-0" />
          </div>
          <div v-if="taxForm.has_npwp" class="p-4 grid grid-cols-1 sm:grid-cols-2 gap-4 border-t border-gray-200 dark:border-gray-700">
            <FormRow :label="t('payroll.npwp')" required :errors="taxErrors?.npwp">
              <TextInput v-model="taxForm.npwp" maxlength="50" :placeholder="t('payroll.npwp_placeholder')" :class="{ 'p-invalid': taxErrors?.npwp }" />
            </FormRow>
            <FormRow :label="t('payroll.npwp_registered_name')" required :errors="taxErrors?.npwp_registered_name">
              <TextInput v-model="taxForm.npwp_registered_name" maxlength="255" :class="{ 'p-invalid': taxErrors?.npwp_registered_name }" />
            </FormRow>
          </div>
        </div>
        <FormRow :label="t('payroll.ptkp_status')">
          <SelectLabel v-model="taxForm.ptkp_status" :options="ptkpOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" showClear />
        </FormRow>
        <!-- Card Metode Pajak — card + radio + deskripsi -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.tax_method') }}</span>
          <div class="space-y-2">
            <div v-for="opt in taxMethodOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-3"
                 :class="taxForm.tax_method === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="taxForm.tax_method = opt.value">
              <RadioButton :modelValue="taxForm.tax_method" :inputId="'tm-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="taxForm.tax_method = $event" class="mt-0.5 shrink-0" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'tm-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500 leading-snug">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <FormRow :label="t('payroll.is_taxable')">
          <ToggleSwitch v-model="taxForm.is_taxable" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="taxErrors?.effective_start_date">
            <DateInput v-model="taxForm.effective_start_date" :class="{ 'p-invalid': taxErrors?.effective_start_date }" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="taxForm.effective_end_date" />
          </FormRow>
        </div>
        <FormRow :label="t('common.status')">
          <SelectLabel v-model="taxForm.status" :options="statusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="taxDialogVisible = false" /><Button :label="taxEditing ? t('common.update') : t('common.save')" size="small" :loading="taxSaving" :disabled="taxSaving" @click="saveTaxProfile" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="deleteTitle" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible = false" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import DateInput from '@/components/DateInput.vue'
import RadioButton from 'primevue/radiobutton'

const props = defineProps({
  employeeId: { type: String, default: '' }
})
const emit = defineEmits(['save'])

const { t, locale } = useI18n()
const toast = useToast()

const loading = ref(false)
const payrollProfiles = ref([])
const bankProfiles = ref([])
const bpjsProfiles = ref([])
const taxProfiles = ref([])
const banks = ref([])

// ── Options ──
const freqOptions = computed(() => ['MONTHLY', 'WEEKLY', 'DAILY'].map(v => ({
  label: t(`payroll.payroll_frequency_${v.toLowerCase()}`),
  desc: t(`payroll.payroll_frequency_${v.toLowerCase()}_desc`),
  value: v
})))
const methodOptions = computed(() => ['BANK_TRANSFER', 'CASH', 'CHEQUE'].map(v => ({
  label: t(`payroll.payment_method_${v.toLowerCase()}`),
  desc: t(`payroll.payment_method_${v.toLowerCase()}_desc`),
  value: v
})))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: statusLabel(v), value: v })))
const riskOptions = computed(() => ['VERY_LOW', 'LOW', 'MEDIUM', 'HIGH', 'VERY_HIGH'].map(v => ({ label: riskLabel(v), value: v })))
const ptkpOptions = computed(() => ['TK/0', 'TK/1', 'TK/2', 'TK/3', 'K/0', 'K/1', 'K/2', 'K/3'].map(v => ({ label: v, value: v })))
const taxMethodOptions = computed(() => ['GROSS', 'GROSS_UP', 'NETT'].map(v => ({
  label: taxMethodLabel(v),
  desc: t(`payroll.tax_method_${v.toLowerCase()}_desc`),
  value: v
})))
const payrollProfileOptions = computed(() => payrollProfiles.value.map(p => ({ label: `${p.payroll_group_code} - ${freqLabel(p.payroll_frequency)}`, value: p.id })))
const bankOptions = computed(() => {
  const opts = banks.value.map(b => ({ label: b.name, value: b.id, code: b.code || '' }))
  // Fallback untuk profil lama yang bank-nya tidak ada di master banks
  if (bankForm.value.bank_name && !opts.some(o => o.label === bankForm.value.bank_name)) {
    opts.unshift({ label: bankForm.value.bank_name, value: bankForm.value.bank_name, code: bankForm.value.bank_code || '' })
  }
  return opts
})

function freqLabel(v) { const key = `payroll.payroll_frequency_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function methodLabel(v) { const key = `payroll.payment_method_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function statusLabel(v) { const key = `payroll.status_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function riskLabel(v) { const key = `payroll.jkk_risk_class_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function taxMethodLabel(v) { const key = `payroll.tax_method_${String(v || '').toLowerCase()}`; return t(key) !== key ? t(key) : v }

// ── Fetch (paginate all, filter by current employee) ──
async function fetchAll(endpoint) {
  const all = []
  let page = 1
  while (true) {
    const res = await api.get(endpoint, { params: { page, per_page: 100 } })
    const rows = res.data?.data || []
    all.push(...rows)
    if (!rows.length || all.length >= (res.data?.total || 0)) break
    page++
  }
  return all
}

async function loadBanks() {
  try {
    const res = await api.get('/api/v1/tenant/settings/banks', { params: { per_page: 200 } })
    banks.value = res.data?.data || []
  } catch { banks.value = [] }
}

async function loadData() {
  if (!props.employeeId) return
  loading.value = true
  try {
    const [pps, banks, bpjs, taxes] = await Promise.all([
      fetchAll('/api/v1/tenant/payroll/employee-payroll-profiles'),
      fetchAll('/api/v1/tenant/payroll/employee-bank-profiles'),
      fetchAll('/api/v1/tenant/payroll/employee-bpjs-profiles'),
      fetchAll('/api/v1/tenant/payroll/employee-tax-profiles')
    ])
    payrollProfiles.value = pps.filter(p => p.employee_id === props.employeeId)
    const profileIds = new Set(payrollProfiles.value.map(p => p.id))
    bankProfiles.value = banks.filter(b => profileIds.has(b.employee_payroll_profile_id))
    bpjsProfiles.value = bpjs.filter(b => profileIds.has(b.employee_payroll_profile_id))
    taxProfiles.value = taxes.filter(x => profileIds.has(x.employee_payroll_profile_id))
    if (payrollProfiles.value.length || bankProfiles.value.length || bpjsProfiles.value.length || taxProfiles.value.length) {
      emit('save')
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    loading.value = false
  }
}

// ── Payroll Profile (create + delete only, matches backend) ──
const payrollDialogVisible = ref(false)
const payrollSaving = ref(false)
const payrollErrors = ref({})
const payrollForm = ref({ payroll_frequency: 'MONTHLY', payment_method: 'BANK_TRANSFER', salary_currency: 'IDR', is_payroll_active: true, effective_start_date: '', effective_end_date: '' })

function resetPayrollForm() {
  payrollErrors.value = {}
  payrollForm.value = { payroll_frequency: 'MONTHLY', payment_method: 'BANK_TRANSFER', salary_currency: 'IDR', is_payroll_active: true, effective_start_date: '', effective_end_date: '' }
}
function openPayrollDialog() {
  resetPayrollForm()
  payrollDialogVisible.value = true
}
async function savePayrollProfile() {
  payrollErrors.value = {}
  if (!payrollForm.value.effective_start_date) { payrollErrors.value = { effective_start_date: [t('form.required')] }; return }
  payrollSaving.value = true
  try {
    await api.post('/api/v1/tenant/payroll/employee-payroll-profiles', {
      employee_id: props.employeeId,
      payroll_frequency: payrollForm.value.payroll_frequency,
      payment_method: payrollForm.value.payment_method,
      salary_currency: payrollForm.value.salary_currency || 'IDR',
      is_payroll_active: payrollForm.value.is_payroll_active,
      effective_start_date: payrollForm.value.effective_start_date,
      effective_end_date: payrollForm.value.effective_end_date || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.profile_created'), life: 3000 })
    payrollDialogVisible.value = false
    await loadData()
    emit('save')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { payrollErrors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    payrollSaving.value = false
  }
}

// ── Bank Profile ──
const bankDialogVisible = ref(false)
const bankSaving = ref(false)
const bankEditing = ref(false)
const bankEditingId = ref(null)
const bankErrors = ref({})
const bankForm = ref({ employee_payroll_profile_id: null, bank_id: '', bank_code: '', bank_name: '', bank_branch: '', bank_account_number: '', bank_account_holder_name: '', is_primary: true, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' })

function resetBankForm() {
  bankErrors.value = {}
  bankEditing.value = false
  bankEditingId.value = null
  bankForm.value = { employee_payroll_profile_id: null, bank_id: '', bank_code: '', bank_name: '', bank_branch: '', bank_account_number: '', bank_account_holder_name: '', is_primary: true, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' }
}
function openBankDialog(item) {
  resetBankForm()
  if (item) {
    bankEditing.value = true
    bankEditingId.value = item.id
    const matched = banks.value.find(b => b.name === (item.bank_name || ''))
    bankForm.value = {
      employee_payroll_profile_id: item.employee_payroll_profile_id || null,
      bank_id: matched ? matched.id : (item.bank_name || ''),
      bank_code: item.bank_code || (matched?.code || ''), bank_name: item.bank_name || '', bank_branch: item.bank_branch || '',
      bank_account_number: item.bank_account_number || '', bank_account_holder_name: item.bank_account_holder_name || '',
      is_primary: item.is_primary ?? true, effective_start_date: item.effective_start_date || '', effective_end_date: item.effective_end_date || null, status: item.status || 'ACTIVE'
    }
  }
  bankDialogVisible.value = true
}
function onBankChange(val) {
  if (!val) return
  const b = banks.value.find(x => x.id === val)
  if (b) {
    bankForm.value.bank_name = b.name
    bankForm.value.bank_code = b.code || ''
  } else {
    bankForm.value.bank_name = val
  }
}
async function saveBankProfile() {
  bankErrors.value = {}
  if (!bankForm.value.employee_payroll_profile_id) { bankErrors.value = { employee_payroll_profile_id: [t('form.required')] }; return }
  if (!bankForm.value.bank_id || !bankForm.value.bank_name?.trim()) { bankErrors.value = { bank_id: [t('form.required')] }; return }
  if (!bankForm.value.bank_account_number?.trim()) { bankErrors.value = { bank_account_number: [t('form.required')] }; return }
  if (!bankForm.value.bank_account_holder_name?.trim()) { bankErrors.value = { bank_account_holder_name: [t('form.required')] }; return }
  if (!bankForm.value.effective_start_date) { bankErrors.value = { effective_start_date: [t('form.required')] }; return }
  bankSaving.value = true
  try {
    const payload = {
      employee_id: props.employeeId, employee_payroll_profile_id: bankForm.value.employee_payroll_profile_id,
      bank_code: bankForm.value.bank_code || null, bank_name: bankForm.value.bank_name.trim(), bank_branch: bankForm.value.bank_branch || null,
      bank_account_number: bankForm.value.bank_account_number.trim(), bank_account_holder_name: bankForm.value.bank_account_holder_name.trim(),
      is_primary: bankForm.value.is_primary, effective_start_date: bankForm.value.effective_start_date,
      effective_end_date: bankForm.value.effective_end_date || null, status: bankForm.value.status
    }
    if (bankEditing.value) {
      await api.put(`/api/v1/tenant/payroll/employee-bank-profiles/${bankEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bank_profile_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/employee-bank-profiles', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bank_profile_created'), life: 3000 })
    }
    bankDialogVisible.value = false
    await loadData()
    emit('save')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { bankErrors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    bankSaving.value = false
  }
}

// ── BPJS Profile ──
const bpjsDialogVisible = ref(false)
const bpjsSaving = ref(false)
const bpjsEditing = ref(false)
const bpjsEditingId = ref(null)
const bpjsErrors = ref({})
const bpjsForm = ref({ employee_payroll_profile_id: null, bpjs_health_active: false, bpjs_health_no: '', bpjs_health_registered_name: '', bpjs_tk_active: false, bpjs_tk_no: '', bpjs_tk_registered_name: '', jkk_risk_class: 'LOW', pension_active: true, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' })

function resetBpjsForm() {
  bpjsErrors.value = {}
  bpjsEditing.value = false
  bpjsEditingId.value = null
  bpjsForm.value = { employee_payroll_profile_id: null, bpjs_health_active: false, bpjs_health_no: '', bpjs_health_registered_name: '', bpjs_tk_active: false, bpjs_tk_no: '', bpjs_tk_registered_name: '', jkk_risk_class: 'LOW', pension_active: true, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' }
}
function openBpjsDialog(item) {
  resetBpjsForm()
  if (item) {
    bpjsEditing.value = true
    bpjsEditingId.value = item.id
    bpjsForm.value = {
      employee_payroll_profile_id: item.employee_payroll_profile_id || null,
      bpjs_health_active: item.bpjs_health_active ?? false, bpjs_health_no: item.bpjs_health_no || '', bpjs_health_registered_name: item.bpjs_health_registered_name || '',
      bpjs_tk_active: item.bpjs_tk_active ?? false, bpjs_tk_no: item.bpjs_tk_no || '', bpjs_tk_registered_name: item.bpjs_tk_registered_name || '',
      jkk_risk_class: item.jkk_risk_class || 'LOW', pension_active: item.pension_active ?? true,
      effective_start_date: item.effective_start_date || '', effective_end_date: item.effective_end_date || null,
      status: item.status || 'ACTIVE'
    }
  }
  bpjsDialogVisible.value = true
}
async function saveBpjsProfile() {
  bpjsErrors.value = {}
  if (!bpjsForm.value.employee_payroll_profile_id) { bpjsErrors.value = { employee_payroll_profile_id: [t('form.required')] }; return }
  if (!bpjsForm.value.effective_start_date) { bpjsErrors.value = { effective_start_date: [t('form.required')] }; return }
  bpjsSaving.value = true
  try {
    const payload = {
      employee_id: props.employeeId, employee_payroll_profile_id: bpjsForm.value.employee_payroll_profile_id,
      bpjs_health_active: bpjsForm.value.bpjs_health_active, bpjs_health_no: bpjsForm.value.bpjs_health_no || null, bpjs_health_registered_name: bpjsForm.value.bpjs_health_registered_name || null,
      bpjs_tk_active: bpjsForm.value.bpjs_tk_active, bpjs_tk_no: bpjsForm.value.bpjs_tk_no || null, bpjs_tk_registered_name: bpjsForm.value.bpjs_tk_registered_name || null,
      jkk_risk_class: bpjsForm.value.jkk_risk_class, pension_active: bpjsForm.value.pension_active,
      effective_start_date: bpjsForm.value.effective_start_date, effective_end_date: bpjsForm.value.effective_end_date || null,
      status: bpjsForm.value.status
    }
    if (bpjsEditing.value) {
      await api.put(`/api/v1/tenant/payroll/employee-bpjs-profiles/${bpjsEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bpjs_profile_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/employee-bpjs-profiles', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.bpjs_profile_created'), life: 3000 })
    }
    bpjsDialogVisible.value = false
    await loadData()
    emit('save')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { bpjsErrors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    bpjsSaving.value = false
  }
}

// ── Tax Profile ──
const taxDialogVisible = ref(false)
const taxSaving = ref(false)
const taxEditing = ref(false)
const taxEditingId = ref(null)
const taxErrors = ref({})
const taxForm = ref({ employee_payroll_profile_id: null, npwp: '', npwp_registered_name: '', ptkp_status: null, tax_method: 'GROSS', is_taxable: true, has_npwp: false, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' })

function resetTaxForm() {
  taxErrors.value = {}
  taxEditing.value = false
  taxEditingId.value = null
  taxForm.value = { employee_payroll_profile_id: null, npwp: '', npwp_registered_name: '', ptkp_status: null, tax_method: 'GROSS', is_taxable: true, has_npwp: false, effective_start_date: '', effective_end_date: null, status: 'ACTIVE' }
}
function openTaxDialog(item) {
  resetTaxForm()
  if (item) {
    taxEditing.value = true
    taxEditingId.value = item.id
    taxForm.value = {
      employee_payroll_profile_id: item.employee_payroll_profile_id || null,
      npwp: item.npwp || '', npwp_registered_name: item.npwp_registered_name || '', ptkp_status: item.ptkp_status || null,
      tax_method: item.tax_method || 'GROSS', is_taxable: item.is_taxable ?? true, has_npwp: item.has_npwp ?? false,
      effective_start_date: item.effective_start_date || '', effective_end_date: item.effective_end_date || null,
      status: item.status || 'ACTIVE'
    }
  }
  taxDialogVisible.value = true
}
async function saveTaxProfile() {
  taxErrors.value = {}
  if (!taxForm.value.employee_payroll_profile_id) { taxErrors.value = { employee_payroll_profile_id: [t('form.required')] }; return }
  if (taxForm.value.has_npwp) {
    if (!taxForm.value.npwp?.trim()) { taxErrors.value = { npwp: [t('form.required')] }; return }
    if (!taxForm.value.npwp_registered_name?.trim()) { taxErrors.value = { npwp_registered_name: [t('form.required')] }; return }
  }
  if (!taxForm.value.effective_start_date) { taxErrors.value = { effective_start_date: [t('form.required')] }; return }
  taxSaving.value = true
  try {
    const payload = {
      employee_id: props.employeeId, employee_payroll_profile_id: taxForm.value.employee_payroll_profile_id,
      npwp: taxForm.value.npwp || null, npwp_registered_name: taxForm.value.npwp_registered_name || null,
      ptkp_status: taxForm.value.ptkp_status || null, tax_method: taxForm.value.tax_method,
      is_taxable: taxForm.value.is_taxable, has_npwp: taxForm.value.has_npwp,
      effective_start_date: taxForm.value.effective_start_date, effective_end_date: taxForm.value.effective_end_date || null,
      status: taxForm.value.status
    }
    if (taxEditing.value) {
      await api.put(`/api/v1/tenant/payroll/employee-tax-profiles/${taxEditingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.tax_profile_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/employee-tax-profiles', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.tax_profile_created'), life: 3000 })
    }
    taxDialogVisible.value = false
    await loadData()
    emit('save')
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { taxErrors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally {
    taxSaving.value = false
  }
}

// ── Delete (shared) ──
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const deleteError = ref('')
const deleteTarget = ref(null)
const deleteType = ref('')
const deleteTitle = computed(() => {
  const keys = { payroll: 'payroll.employee_profiles', bank: 'payroll.bank_profiles', bpjs: 'payroll.bpjs_profiles', tax: 'payroll.tax_profiles' }
  return t(keys[deleteType.value] || 'payroll.employee_profiles')
})
const deleteMessage = computed(() => {
  if (!deleteTarget.value) return t('common.no_data')
  const d = deleteTarget.value
  if (deleteType.value === 'bank') return `${d.bank_name} ${d.bank_account_number}`
  if (deleteType.value === 'bpjs') return `${d.bpjs_health_no || d.bpjs_tk_no || '-'}`
  if (deleteType.value === 'tax') return `${d.npwp || '-'}`
  return d.payroll_group_code || ''
})

function confirmDelete(type, item) {
  deleteType.value = type
  deleteTarget.value = item
  deleteError.value = ''
  deleteDialogVisible.value = true
}
async function handleDelete() {
  deleting.value = true
  deleteError.value = ''
  try {
    const endpoints = { payroll: 'employee-payroll-profiles', bank: 'employee-bank-profiles', bpjs: 'employee-bpjs-profiles', tax: 'employee-tax-profiles' }
    await api.delete(`/api/v1/tenant/payroll/${endpoints[deleteType.value]}/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.profile_deleted'), life: 3000 })
    deleteDialogVisible.value = false
    await loadData()
    emit('save')
  } catch (e) {
    deleteError.value = e.response?.data?.error?.message || t('message.operation_failed')
  } finally {
    deleting.value = false
  }
}

onMounted(() => { loadBanks(); loadData() })
</script>
