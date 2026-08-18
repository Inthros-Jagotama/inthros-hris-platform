<template>
  <div class="space-y-4">
    <!-- ── Sub-tabs ── -->
    <div class="flex items-center gap-1 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
        :class="activeTab === tab.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="activeTab = tab.key"
      >
        {{ t(tab.labelKey) }}
      </button>
    </div>

    <!-- ── Settings ── -->
    <div v-if="activeTab === 'settings'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <span v-if="settingTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ settingTotal }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_pph21_setting')" icon="pi pi-plus" size="small" @click="openSettingDialog()" />
      </div>
      <DataTable :value="settingItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="settingLoading">
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-percentage text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.pph21_settings_empty') }}</p></div></template>
        <Column field="setting_code" :header="t('payroll.pph21_setting_code')" sortable style="width:140px"><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium font-mono text-xs">{{ data.setting_code }}</span></template></Column>
        <Column field="setting_name" :header="t('payroll.pph21_setting_name')" sortable><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.setting_name }}</span></template></Column>
        <Column field="calculation_method" :header="t('payroll.calculation_method')" sortable style="width:200px"><template #body="{data}"><Tag :value="calculationMethodLabel(data.calculation_method)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="default_tax_method" :header="t('payroll.default_tax_method')" sortable style="width:120px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.default_tax_method }}</span></template></Column>
        <Column field="effective_start_date" :header="t('payroll.effective_start_date')" sortable style="width:140px"><template #body="{data}"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ formatDate(data.effective_start_date, locale) }}</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:100px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openSettingDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteSetting(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── TER Rates ── -->
    <div v-if="activeTab === 'ter'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <span v-if="terTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ terTotal }} {{ t('common.items') }}</span>
          <IconField class="w-56">
            <InputIcon class="pi pi-search" />
            <InputText v-model="terSearch" class="!w-full" size="small" :placeholder="t('common.search')" @input="onTerSearchInput" />
          </IconField>
        </div>
        <Button :label="t('ters.new')" icon="pi pi-plus" size="small" @click="openTerDialog()" />
      </div>
      <DataTable
        :value="terItems"
        lazy
        :totalRecords="terTotal"
        :first="terFirst"
        :rows="terPerPage"
        @page="onTerPage"
        paginator
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        :rowsPerPageOptions="[10, 15, 25, 50]"
        size="small"
        class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        :loading="terLoading"
      >
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-calculator text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('ters.empty_title') }}</p></div></template>
        <Column field="group" :header="t('ters.group')" sortable style="width:100px"><template #body="{data}"><Tag :value="data.group" :severity="terGroupSeverity(data.group)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="bruto_min" :header="t('ters.bruto_min')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.bruto_min) }}</span></template></Column>
        <Column field="bruto_max" :header="t('ters.bruto_max')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ data.bruto_max != null ? formatMoney(data.bruto_max) : '∞' }}</span></template></Column>
        <Column field="rate" :header="t('ters.rate')" sortable style="width:120px"><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.rate }}%</span></template></Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openTerDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteTer(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── PTKP Rates (tabel ptkps — satu sumber kebenaran) ── -->
    <div v-if="activeTab === 'ptkp'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <div class="flex items-center gap-2">
          <span v-if="ptkpTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">{{ ptkpTotal }} {{ t('common.items') }}</span>
          <IconField class="w-56">
            <InputIcon class="pi pi-search" />
            <InputText v-model="ptkpSearch" class="!w-full" size="small" :placeholder="t('common.search')" @input="onPtkpSearchInput" />
          </IconField>
        </div>
        <Button :label="t('ptkps.new')" icon="pi pi-plus" size="small" @click="openPtkpDialog()" />
      </div>
      <DataTable
        :value="ptkpItems"
        lazy
        :totalRecords="ptkpTotal"
        :first="ptkpFirst"
        :rows="ptkpPerPage"
        @page="onPtkpPage"
        paginator
        paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
        :rowsPerPageOptions="[10, 15, 25, 50]"
        size="small"
        class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
        :loading="ptkpLoading"
      >
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-receipt text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('ptkps.empty_title') }}</p></div></template>
        <Column field="name" :header="t('ptkps.name')" sortable><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.name }}</span></template></Column>
        <Column field="group" :header="t('ptkps.group')" sortable style="width:100px"><template #body="{data}"><Tag :value="data.group" :severity="terGroupSeverity(data.group)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column field="ptkp" :header="t('ptkps.ptkp')" sortable style="width:180px"><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium font-mono text-xs">{{ formatMoney(data.ptkp) }}</span></template></Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openPtkpDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeletePtkp(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── Tax Brackets ── -->
    <div v-if="activeTab === 'brackets'" class="space-y-1">
      <div class="flex items-center justify-between gap-2 flex-wrap">
        <span v-if="bracketItems.length" class="text-xs text-gray-400 dark:text-gray-500">{{ bracketItems.length }} {{ t('common.items') }}</span>
        <Button :label="t('payroll.new_rate_component')" icon="pi pi-plus" size="small" @click="openBracketDialog()" />
      </div>
      <DataTable :value="bracketItems" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="bracketLoading">
        <template #empty><div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500"><i class="pi pi-chart-bar text-3xl mb-2 opacity-50"></i><p class="text-sm font-medium">{{ t('payroll.pph21_settings_empty') }}</p></div></template>
        <Column field="bracket_order" :header="t('payroll.display_order')" sortable style="width:100px"><template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.bracket_order }}</span></template></Column>
        <Column field="lower_bound" :header="t('payroll.min_base_amount')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.lower_bound) }}</span></template></Column>
        <Column field="upper_bound" :header="t('payroll.max_base_amount')" sortable style="width:160px"><template #body="{data}"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ data.upper_bound ? formatMoney(data.upper_bound) : '∞' }}</span></template></Column>
        <Column field="rate_percent" :header="t('payroll.rate_percent')" sortable style="width:100px"><template #body="{data}"><span class="text-navy-800 dark:text-gray-100 font-medium">{{ data.rate_percent }}%</span></template></Column>
        <Column field="status" :header="t('common.status')" sortable style="width:90px"><template #body="{data}"><Tag :value="statusLabel(data.status)" :severity="data.status === 'ACTIVE' ? 'success' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
        <Column :header="t('common.actions')" style="width:100px" frozen alignFrozen="right"><template #body="{data}"><div class="flex items-center gap-1 justify-end"><Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openBracketDialog(data)" /><Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="confirmDeleteBracket(data)" /></div></template></Column>
      </DataTable>
    </div>

    <!-- ── Setting dialog ── -->
    <Dialog v-model:visible="settingDialogVisible" :header="editingSetting ? t('payroll.pph21') : t('payroll.new_pph21_setting')" modal :style="{ width: 'min(1100px, 95vw)' }" @hide="resetSettingForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.pph21_setting_name')" required :errors="errors?.setting_name">
          <TextInput v-model="settingForm.setting_name" maxlength="150" autofocus :placeholder="t('payroll.pph21_setting_name')" :class="{'p-invalid':errors?.setting_name}" />
        </FormRow>
        <!-- Metode Perhitungan — card dengan pilihan radio berdeskripsi (2 kolom) -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.calculation_method') }}</span>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
            <div v-for="opt in calculationMethodOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-3"
                 :class="settingForm.calculation_method === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="settingForm.calculation_method = opt.value">
              <RadioButton :modelValue="settingForm.calculation_method" :inputId="'calc-method-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="settingForm.calculation_method = $event" class="mt-0.5" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'calc-method-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <!-- Metode Pajak Default — card dengan pilihan radio berdeskripsi (3 kolom) -->
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
          <span class="text-sm font-medium block mb-3">{{ t('payroll.default_tax_method') }}</span>
          <div class="grid grid-cols-1 sm:grid-cols-3 gap-2">
            <div v-for="opt in taxMethodOptions" :key="opt.value"
                 class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 cursor-pointer select-none transition-all duration-150 flex items-start gap-3"
                 :class="settingForm.default_tax_method === opt.value
                   ? 'border-emerald-400 dark:border-emerald-500 bg-emerald-50 dark:bg-emerald-900/20 shadow-sm'
                   : 'hover:border-gray-300 dark:hover:border-gray-500'"
                 @click="settingForm.default_tax_method = opt.value">
              <RadioButton :modelValue="settingForm.default_tax_method" :inputId="'tax-method-' + opt.value.toLowerCase()" :value="opt.value" @update:modelValue="settingForm.default_tax_method = $event" class="mt-0.5" />
              <div class="flex flex-col gap-0.5 min-w-0">
                <label :for="'tax-method-' + opt.value.toLowerCase()" class="text-sm font-medium cursor-pointer select-none">{{ opt.label }}</label>
                <span class="text-xs text-gray-400 dark:text-gray-500">{{ opt.desc }}</span>
              </div>
            </div>
          </div>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.occupational_expense_rate')">
            <InputNumber v-model="settingForm.occupational_expense_rate_percent" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.occupational_expense_max_monthly')">
            <InputNumber v-model="settingForm.occupational_expense_max_monthly" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.occupational_expense_max_yearly')">
            <InputNumber v-model="settingForm.occupational_expense_max_yearly" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.annualization_months')">
            <InputNumber v-model="settingForm.annualization_months" class="!w-full" :min="1" :max="12" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.pkp_rounding_unit')">
            <InputNumber v-model="settingForm.pkp_rounding_unit" class="!w-full" :min="0" size="small" />
          </FormRow>
          <FormRow :label="t('payroll.non_npwp_multiplier_percent')">
            <InputNumber v-model="settingForm.non_npwp_multiplier_percent" class="!w-full" :min="0" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('payroll.rounding_mode')">
          <SelectLabel v-model="settingForm.rounding_mode" :options="roundingOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
        </FormRow>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
            <DateInput v-model="settingForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
          </FormRow>
          <FormRow :label="t('payroll.effective_end_date')">
            <DateInput v-model="settingForm.effective_end_date" />
          </FormRow>
        </div>
        <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 flex items-center justify-between gap-3">
          <div>
            <span class="text-sm font-medium block">{{ t('common.status') }}</span>
            <span class="text-xs text-gray-400 dark:text-gray-500">{{ t('payroll.status_desc') }}</span>
          </div>
          <ToggleSwitch v-model="settingForm.status" :true-value="'ACTIVE'" :false-value="'INACTIVE'" :label="statusLabel(settingForm.status)" class="shrink-0" />
        </div>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="settingDialogVisible=false" /><Button :label="editingSetting ? t('common.update') : t('common.save')" size="small" :loading="settingSaving" :disabled="settingSaving" @click="handleSaveSetting" /></div></template>
    </Dialog>

    <!-- ── TER dialog ── -->
    <Dialog v-model:visible="terDialogVisible" :header="editingTer ? t('ters.edit') : t('ters.new')" modal :style="{ width: '520px' }" @hide="resetTerForm">
      <div class="space-y-4">
        <FormRow :label="t('ters.group')" required :errors="errors?.group">
          <SelectLabel v-model="terForm.group" :options="terGroupOptions" optionLabel="label" optionValue="value" :placeholder="t('ters.select_group')" :class="{'p-invalid':errors?.group}" autofocus />
        </FormRow>
        <div class="grid grid-cols-2 gap-3">
          <FormRow :label="t('ters.bruto_min')" :errors="errors?.bruto_min">
            <InputNumber v-model="terForm.bruto_min" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
          <FormRow :label="t('ters.bruto_max')" :errors="errors?.bruto_max">
            <InputNumber v-model="terForm.bruto_max" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
        </div>
        <FormRow :label="t('ters.rate') + ' (%)'" required :errors="errors?.rate">
          <InputNumber v-model="terForm.rate" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="terDialogVisible=false" /><Button :label="editingTer ? t('common.update') : t('common.save')" size="small" :loading="terSaving" :disabled="terSaving" @click="handleSaveTer" /></div></template>
    </Dialog>

    <!-- ── PTKP dialog ── -->
    <Dialog v-model:visible="ptkpDialogVisible" :header="editingPtkp ? t('ptkps.edit') : t('ptkps.new')" modal :style="{ width: '520px' }" @hide="resetPtkpForm">
      <div class="space-y-4">
        <FormRow :label="t('ptkps.name')" required :errors="errors?.name">
          <TextInput v-model="ptkpForm.name" maxlength="255" autofocus :placeholder="t('ptkps.name')" :class="{'p-invalid':errors?.name}" />
        </FormRow>
        <FormRow :label="t('ptkps.group')" required :errors="errors?.group">
          <SelectLabel v-model="ptkpForm.group" :options="ptkpGroupOptions" optionLabel="label" optionValue="value" :placeholder="t('ptkps.select_group')" :class="{'p-invalid':errors?.group}" />
        </FormRow>
        <FormRow :label="t('ptkps.ptkp')" required :errors="errors?.ptkp">
          <InputNumber v-model="ptkpForm.ptkp" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" :class="{'p-invalid':errors?.ptkp}" />
        </FormRow>
        <p class="text-xs text-gray-400 dark:text-gray-500">{{ t('ptkps.code_hint') }}</p>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="ptkpDialogVisible=false" /><Button :label="editingPtkp ? t('common.update') : t('common.save')" size="small" :loading="ptkpSaving" :disabled="ptkpSaving" @click="handleSavePtkp" /></div></template>
    </Dialog>

    <!-- ── Bracket dialog ── -->
    <Dialog v-model:visible="bracketDialogVisible" :header="editingBracket ? t('payroll.edit_rate_component') : t('payroll.new_rate_component')" modal :style="{ width: '480px' }" @hide="resetBracketForm">
      <div class="space-y-4">
        <FormRow :label="t('payroll.display_order')" required :errors="errors?.bracket_order">
          <InputNumber v-model="bracketForm.bracket_order" class="!w-full" :min="1" size="small" :class="{'p-invalid':errors?.bracket_order}" />
        </FormRow>
        <FormRow :label="t('payroll.min_base_amount')" required :errors="errors?.lower_bound">
          <InputNumber v-model="bracketForm.lower_bound" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" :class="{'p-invalid':errors?.lower_bound}" />
        </FormRow>
        <FormRow :label="t('payroll.max_base_amount')">
          <InputNumber v-model="bracketForm.upper_bound" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('payroll.rate_percent')" required :errors="errors?.rate_percent">
          <InputNumber v-model="bracketForm.rate_percent" class="!w-full" :min="0" :max="100" :step="0.01" :minFractionDigits="2" :maxFractionDigits="2" size="small" :class="{'p-invalid':errors?.rate_percent}" />
        </FormRow>
        <FormRow :label="t('payroll.effective_start_date')" required :errors="errors?.effective_start_date">
          <DateInput v-model="bracketForm.effective_start_date" :class="{'p-invalid':errors?.effective_start_date}" />
        </FormRow>
        <FormRow :label="t('payroll.effective_end_date')">
          <DateInput v-model="bracketForm.effective_end_date" />
        </FormRow>
      </div>
      <template #footer><div class="flex items-center justify-end gap-2"><Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="bracketDialogVisible=false" /><Button :label="t('common.save')" size="small" :loading="bracketSaving" :disabled="bracketSaving" @click="handleSaveBracket" /></div></template>
    </Dialog>

    <ConfirmDeleteDialog v-model:visible="deleteDialogVisible" :title="t('payroll.pph21')" :message="deleteMessage" :loading="deleting" :errorMsg="deleteError" @confirm="handleDelete" @cancel="deleteDialogVisible=false" />
    <ConfirmDeleteDialog v-model:visible="terDeleteDialogVisible" :title="t('ters.confirm_delete_title')" :message="`${t('ters.group')}: ${terDeleteTarget?.group || ''}`" :loading="terDeleting" :errorMsg="terDeleteError" @confirm="handleDeleteTer" @cancel="terDeleteDialogVisible=false" />
    <ConfirmDeleteDialog v-model:visible="ptkpDeleteDialogVisible" :title="t('ptkps.confirm_delete_title')" :message="`${t('ptkps.name')}: ${ptkpDeleteTarget?.name || ''}`" :loading="ptkpDeleting" :errorMsg="ptkpDeleteError" @confirm="handleDeletePtkp" @cancel="ptkpDeleteDialogVisible=false" />
    <ConfirmDeleteDialog v-model:visible="bracketDeleteDialogVisible" :title="t('payroll.delete_rate_component')" :message="`${t('payroll.display_order')}: ${bracketDeleteTarget?.bracket_order || ''} (${bracketDeleteTarget?.rate_percent ?? ''}%)`" :loading="bracketDeleting" :errorMsg="bracketDeleteError" @confirm="handleDeleteBracket" @cancel="bracketDeleteDialogVisible=false" />
  </div>
</template>
<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'; import { useI18n } from '@/composables/useI18n'; import { getValidationErrors } from '@/services/responseHandler'; import { formatDate } from '@/utils/formatDate'; import api from '@/services/api'
import DataTable from 'primevue/datatable'; import Column from 'primevue/column'; import Button from 'primevue/button'; import InputNumber from 'primevue/inputnumber'; import Tag from 'primevue/tag'; import Dialog from 'primevue/dialog'; import InputText from 'primevue/inputtext'; import InputIcon from 'primevue/inputicon'; import IconField from 'primevue/iconfield'; import RadioButton from 'primevue/radiobutton'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import DateInput from '@/components/DateInput.vue'
const { t, locale } = useI18n(); const toast = useToast()
const activeTab = ref('settings')
const tabs = [
  { key: 'settings', labelKey: 'payroll.pph21' },
  { key: 'ter', labelKey: 'ters.title' },
  { key: 'ptkp', labelKey: 'payroll.ptkp_full' },
  { key: 'brackets', labelKey: 'payroll.rate_components' }
]
const settingItems = ref([]); const settingLoading = ref(false); const settingTotal = ref(0)
const terItems = ref([]); const terLoading = ref(false); const terTotal = ref(0); const terCurrentPage = ref(1); const terPerPage = ref(15); const terSearch = ref('')
const ptkpItems = ref([]); const ptkpLoading = ref(false); const ptkpTotal = ref(0); const ptkpCurrentPage = ref(1); const ptkpPerPage = ref(15); const ptkpSearch = ref('')
const bracketItems = ref([]); const bracketLoading = ref(false)
const settingDialogVisible = ref(false); const editingSetting = ref(false); const editingSettingId = ref(null); const settingSaving = ref(false)
const terDialogVisible = ref(false); const editingTer = ref(false); const editingTerId = ref(null); const terSaving = ref(false)
const terDeleteDialogVisible = ref(false); const terDeleting = ref(false); const terDeleteError = ref(''); const terDeleteTarget = ref(null)
const ptkpDialogVisible = ref(false); const editingPtkp = ref(false); const editingPtkpId = ref(null); const ptkpSaving = ref(false)
const ptkpDeleteDialogVisible = ref(false); const ptkpDeleting = ref(false); const ptkpDeleteError = ref(''); const ptkpDeleteTarget = ref(null)
const bracketDialogVisible = ref(false); const editingBracket = ref(false); const editingBracketId = ref(null); const bracketSaving = ref(false)
const bracketDeleteDialogVisible = ref(false); const bracketDeleting = ref(false); const bracketDeleteError = ref(''); const bracketDeleteTarget = ref(null)
const deleteDialogVisible = ref(false); const deleting = ref(false); const deleteError = ref(''); const deleteTarget = ref(null)
const errors = ref({})
const settingForm = ref({ setting_code: '', setting_name: '', calculation_method: 'TER', default_tax_method: 'GROSS', occupational_expense_rate_percent: 5, occupational_expense_max_monthly: 500000, occupational_expense_max_yearly: 6000000, annualization_months: 12, pkp_rounding_unit: 1000, non_npwp_multiplier_percent: 100, rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' })
const terForm = ref({ group: '', bruto_min: null, bruto_max: null, rate: 0 })
const ptkpForm = ref({ name: '', group: '', ptkp: 0 })
const bracketForm = ref({ bracket_order: 1, lower_bound: 0, upper_bound: null, rate_percent: 0, effective_start_date: '', effective_end_date: '' })

const calculationMethodOptions = computed(() => ['TER', 'REGULAR_GROSS_ANNUALIZED'].map(v => ({ label: t(`payroll.calculation_method_${v.toLowerCase()}`), desc: t(`payroll.calculation_method_${v.toLowerCase()}_desc`), value: v })))
const terFirst = computed(() => (terCurrentPage.value - 1) * terPerPage.value)
const ptkpFirst = computed(() => (ptkpCurrentPage.value - 1) * ptkpPerPage.value)
const terGroupOptions = computed(() => [
  { value: 'A', label: 'Group A — TK/0, TK/1, K/0' },
  { value: 'B', label: 'Group B — TK/2, TK/3, K/1, K/2' },
  { value: 'C', label: 'Group C — K/3' }
])
const ptkpGroupOptions = computed(() => terGroupOptions.value)
const taxMethodOptions = computed(() => ['GROSS', 'GROSS_UP', 'NETT'].map(v => ({ label: t(`payroll.tax_method_${v.toLowerCase()}`), desc: t(`payroll.tax_method_${v.toLowerCase()}_desc`), value: v })))
const roundingOptions = computed(() => ['NONE', 'ROUND', 'CEIL', 'FLOOR'].map(v => ({ label: t(`payroll.rounding_mode_${v.toLowerCase()}`), value: v })))
const statusOptions = computed(() => ['ACTIVE', 'INACTIVE'].map(v => ({ label: t(`payroll.status_${v.toLowerCase()}`), value: v })))
const deleteMessage = computed(() => deleteTarget.value ? `${t('payroll.pph21_setting_code')}: ${deleteTarget.value.setting_code}` : t('common.no_data'))

function statusLabel(v) { const key = `payroll.status_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function calculationMethodLabel(v) { const key = `payroll.calculation_method_${String(v||'').toLowerCase()}`; return t(key) !== key ? t(key) : v }
function terGroupSeverity(g) { return { A: 'info', B: 'success', C: 'warn' }[g] || 'info' }
function formatMoney(val) { const n = Number(val || 0); return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(n) }

function nameInitials(name) {
  const words = (name || '').trim().split(/\s+/).filter(Boolean)
  return words.map(w => w[0]).join('').toUpperCase().replace(/[^A-Z0-9]/g, '')
}
function generateCode(name) {
  const prefix = nameInitials(name)
  const d = new Date()
  const pad = n => String(n).padStart(2, '0')
  const stamp = `${d.getFullYear()}${pad(d.getMonth() + 1)}${pad(d.getDate())}${pad(d.getHours())}${pad(d.getMinutes())}${pad(d.getSeconds())}`
  const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'
  const rand = Array.from({ length: 4 }, () => alphabet[Math.floor(Math.random() * alphabet.length)]).join('')
  return prefix ? `${prefix}_${stamp}_${rand}` : ''
}
watch(() => settingForm.value.setting_name, (nv) => { if (!editingSetting.value) settingForm.value.setting_code = generateCode(nv) })

async function loadSettings() {
  settingLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/pph21-settings', { params: { per_page: 200 } })
    settingItems.value = res.data?.data || []
    settingTotal.value = res.data?.total || 0
  } catch(e) { toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) }
  finally { settingLoading.value = false }
}
async function loadTers() {
  terLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/ters', { params: { page: terCurrentPage.value, per_page: terPerPage.value, search: terSearch.value.trim() } })
    terItems.value = res.data?.data || []
    terTotal.value = res.data?.total || 0
    if (res.data?.page) terCurrentPage.value = res.data.page
  } catch(e) { toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) }
  finally { terLoading.value = false }
}
let terSearchTimer = null
function onTerSearchInput() {
  clearTimeout(terSearchTimer)
  terSearchTimer = setTimeout(() => { terCurrentPage.value = 1; loadTers() }, 350)
}
function onTerPage(event) {
  terCurrentPage.value = event.page + 1
  terPerPage.value = event.rows
  loadTers()
}
async function loadPtkp() {
  ptkpLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/settings/ptkps', { params: { page: ptkpCurrentPage.value, per_page: ptkpPerPage.value, search: ptkpSearch.value.trim() } })
    ptkpItems.value = res.data?.data || []
    ptkpTotal.value = res.data?.total || 0
    if (res.data?.page) ptkpCurrentPage.value = res.data.page
  } catch(e) { toast.add({severity:'error',summary:t('message.error'),detail:e.response?.data?.error?.message||t('message.failed_to_load'),life:4000}) } finally { ptkpLoading.value = false }
}
let ptkpSearchTimer = null
function onPtkpSearchInput() {
  clearTimeout(ptkpSearchTimer)
  ptkpSearchTimer = setTimeout(() => { ptkpCurrentPage.value = 1; loadPtkp() }, 350)
}
function onPtkpPage(event) {
  ptkpCurrentPage.value = event.page + 1
  ptkpPerPage.value = event.rows
  loadPtkp()
}
async function loadBrackets() {
  bracketLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/payroll/pph21-tax-brackets', { params: { per_page: 200 } })
    bracketItems.value = res.data?.data || []
  } catch { bracketItems.value = [] } finally { bracketLoading.value = false }
}

function openSettingDialog(item) {
  editingSetting.value = !!item; editingSettingId.value = item?.id || null; errors.value = {}
  settingForm.value = {
    setting_code: item?.setting_code || '', setting_name: item?.setting_name || '',
    calculation_method: item?.calculation_method || 'TER',
    default_tax_method: item?.default_tax_method || 'GROSS',
    occupational_expense_rate_percent: item?.occupational_expense_rate_percent ?? 5,
    occupational_expense_max_monthly: item?.occupational_expense_max_monthly ?? 500000,
    occupational_expense_max_yearly: item?.occupational_expense_max_yearly ?? 6000000,
    annualization_months: item?.annualization_months ?? 12,
    pkp_rounding_unit: item?.pkp_rounding_unit ?? 1000,
    non_npwp_multiplier_percent: item?.non_npwp_multiplier_percent ?? 100,
    rounding_mode: item?.rounding_mode || 'ROUND',
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || '',
    status: item?.status || 'ACTIVE'
  }
  settingDialogVisible.value = true
}
function resetSettingForm() {
  settingForm.value = { setting_code: '', setting_name: '', calculation_method: 'TER', default_tax_method: 'GROSS', occupational_expense_rate_percent: 5, occupational_expense_max_monthly: 500000, occupational_expense_max_yearly: 6000000, annualization_months: 12, pkp_rounding_unit: 1000, non_npwp_multiplier_percent: 100, rounding_mode: 'ROUND', effective_start_date: '', effective_end_date: '', status: 'ACTIVE' }
  errors.value = {}; editingSetting.value = false; editingSettingId.value = null
}
async function handleSaveSetting() {
  errors.value = {}
  if (!settingForm.value.setting_code?.trim()) settingForm.value.setting_code = generateCode(settingForm.value.setting_name)
  if (!settingForm.value.setting_code?.trim()) { errors.value = { setting_name: [t('form.required')] }; return }
  if (!settingForm.value.setting_name?.trim()) { errors.value = { setting_name: [t('form.required')] }; return }
  if (!settingForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  settingSaving.value = true
  try {
    const payload = {
      setting_code: settingForm.value.setting_code.trim(),
      setting_name: settingForm.value.setting_name.trim(),
      calculation_method: settingForm.value.calculation_method,
      default_tax_method: settingForm.value.default_tax_method,
      occupational_expense_rate_percent: settingForm.value.occupational_expense_rate_percent,
      occupational_expense_max_monthly: settingForm.value.occupational_expense_max_monthly,
      occupational_expense_max_yearly: settingForm.value.occupational_expense_max_yearly,
      annualization_months: settingForm.value.annualization_months,
      pkp_rounding_unit: settingForm.value.pkp_rounding_unit,
      non_npwp_multiplier_percent: settingForm.value.non_npwp_multiplier_percent,
      rounding_mode: settingForm.value.rounding_mode,
      effective_start_date: settingForm.value.effective_start_date,
      effective_end_date: settingForm.value.effective_end_date || null,
      status: settingForm.value.status
    }
    if (editingSetting.value) {
      await api.put(`/api/v1/tenant/payroll/pph21-settings/${editingSettingId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/pph21-settings', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_created'), life: 3000 })
    }
    settingDialogVisible.value = false; await loadSettings()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { settingSaving.value = false }
}

function openTerDialog(item) {
  editingTer.value = !!item; editingTerId.value = item?.id || null; errors.value = {}
  terForm.value = { group: item?.group || '', bruto_min: item?.bruto_min ?? null, bruto_max: item?.bruto_max ?? null, rate: item?.rate ?? 0 }
  terDialogVisible.value = true
}
function resetTerForm() {
  terForm.value = { group: '', bruto_min: null, bruto_max: null, rate: 0 }
  errors.value = {}; editingTer.value = false; editingTerId.value = null
}
async function handleSaveTer() {
  errors.value = {}
  if (!terForm.value.group) { errors.value = { group: [t('form.required')] }; return }
  if (terForm.value.bruto_min === null || terForm.value.bruto_min === undefined) { errors.value = { bruto_min: [t('form.required')] }; return }
  terSaving.value = true
  try {
    const payload = { group: terForm.value.group, bruto_min: terForm.value.bruto_min, bruto_max: terForm.value.bruto_max ?? null, rate: terForm.value.rate }
    if (editingTer.value) {
      await api.put(`/api/v1/tenant/settings/ters/${editingTerId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('ters.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/settings/ters', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('ters.created'), life: 3000 })
    }
    terDialogVisible.value = false; await loadTers()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { terSaving.value = false }
}
function confirmDeleteTer(item) { terDeleteTarget.value = item; terDeleteError.value = ''; terDeleteDialogVisible.value = true }
async function handleDeleteTer() {
  terDeleting.value = true; terDeleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/ters/${terDeleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('ters.deleted'), life: 3000 })
    terDeleteDialogVisible.value = false; await loadTers()
  } catch(e) { terDeleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { terDeleting.value = false }
}
function openPtkpDialog(item) {
  editingPtkp.value = !!item; editingPtkpId.value = item?.id || null; errors.value = {}
  ptkpForm.value = { name: item?.name || '', group: item?.group || '', ptkp: item?.ptkp ?? 0 }
  ptkpDialogVisible.value = true
}
function resetPtkpForm() {
  ptkpForm.value = { name: '', group: '', ptkp: 0 }
  errors.value = {}; editingPtkp.value = false; editingPtkpId.value = null
}
async function handleSavePtkp() {
  errors.value = {}
  if (!ptkpForm.value.name?.trim()) { errors.value = { name: [t('form.required')] }; return }
  if (!ptkpForm.value.group) { errors.value = { group: [t('form.required')] }; return }
  if (ptkpForm.value.ptkp === null || ptkpForm.value.ptkp === undefined) { errors.value = { ptkp: [t('form.required')] }; return }
  ptkpSaving.value = true
  try {
    const payload = { name: ptkpForm.value.name.trim(), group: ptkpForm.value.group, ptkp: ptkpForm.value.ptkp }
    if (editingPtkp.value) {
      await api.put(`/api/v1/tenant/settings/ptkps/${editingPtkpId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('ptkps.updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/settings/ptkps', payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('ptkps.created'), life: 3000 })
    }
    ptkpDialogVisible.value = false; await loadPtkp()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { ptkpSaving.value = false }
}
function confirmDeletePtkp(item) { ptkpDeleteTarget.value = item; ptkpDeleteError.value = ''; ptkpDeleteDialogVisible.value = true }
async function handleDeletePtkp() {
  ptkpDeleting.value = true; ptkpDeleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/settings/ptkps/${ptkpDeleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('ptkps.deleted'), life: 3000 })
    ptkpDeleteDialogVisible.value = false; await loadPtkp()
  } catch(e) { ptkpDeleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { ptkpDeleting.value = false }
}

function openBracketDialog(item) {
  editingBracket.value = !!item; editingBracketId.value = item?.id || null; errors.value = {}
  bracketForm.value = {
    bracket_order: item?.bracket_order ?? bracketItems.value.length + 1,
    lower_bound: item?.lower_bound ?? 0,
    upper_bound: item?.upper_bound ?? null,
    rate_percent: item?.rate_percent ?? 0,
    effective_start_date: item?.effective_start_date || '', effective_end_date: item?.effective_end_date || ''
  }
  bracketDialogVisible.value = true
}
function resetBracketForm() {
  bracketForm.value = { bracket_order: 1, lower_bound: 0, upper_bound: null, rate_percent: 0, effective_start_date: '', effective_end_date: '' }
  errors.value = {}; editingBracket.value = false; editingBracketId.value = null
}
async function handleSaveBracket() {
  errors.value = {}
  if (!bracketForm.value.bracket_order) { errors.value = { bracket_order: [t('form.required')] }; return }
  if (bracketForm.value.lower_bound === null || bracketForm.value.lower_bound === undefined) { errors.value = { lower_bound: [t('form.required')] }; return }
  if (!bracketForm.value.effective_start_date) { errors.value = { effective_start_date: [t('form.required')] }; return }
  bracketSaving.value = true
  try {
    const payload = {
      bracket_order: bracketForm.value.bracket_order,
      lower_bound: bracketForm.value.lower_bound,
      upper_bound: bracketForm.value.upper_bound || null,
      rate_percent: bracketForm.value.rate_percent || 0,
      effective_start_date: bracketForm.value.effective_start_date,
      effective_end_date: bracketForm.value.effective_end_date || null
    }
    if (editingBracket.value) {
      await api.put(`/api/v1/tenant/payroll/pph21-tax-brackets/${editingBracketId.value}`, payload)
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_updated'), life: 3000 })
    } else {
      await api.post('/api/v1/tenant/payroll/pph21-tax-brackets', { ...payload, status: 'ACTIVE' })
      toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_created'), life: 3000 })
    }
    bracketDialogVisible.value = false; await loadBrackets()
  } catch(e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) { errors.value = fe }
    else { toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 }) }
  } finally { bracketSaving.value = false }
}
function confirmDeleteBracket(item) { bracketDeleteTarget.value = item; bracketDeleteError.value = ''; bracketDeleteDialogVisible.value = true }
async function handleDeleteBracket() {
  bracketDeleting.value = true; bracketDeleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/payroll/pph21-tax-brackets/${bracketDeleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.rate_component_deleted'), life: 3000 })
    bracketDeleteDialogVisible.value = false; await loadBrackets()
  } catch(e) { bracketDeleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { bracketDeleting.value = false }
}

function confirmDeleteSetting(item) { deleteTarget.value = item; deleteError.value = ''; deleteDialogVisible.value = true }
async function handleDelete() {
  deleting.value = true; deleteError.value = ''
  try {
    await api.delete(`/api/v1/tenant/payroll/pph21-settings/${deleteTarget.value.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.pph21_setting_deleted'), life: 3000 })
    deleteDialogVisible.value = false; await loadSettings()
  } catch(e) { deleteError.value = e.response?.data?.error?.message || t('message.operation_failed') }
  finally { deleting.value = false }
}
onMounted(() => { loadSettings(); loadTers(); loadPtkp(); loadBrackets() })
</script>
