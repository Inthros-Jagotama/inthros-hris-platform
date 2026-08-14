<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="run">
      <!-- ── Header ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ run.run_code }}</h2>
              <Tag :value="statusLabel(run.status)" :severity="statusSeverity(run.status)" class="!text-xs !px-1.5 !py-0.5" />
              <Tag :value="t('payroll.run_type_' + run.run_type.toLowerCase())" severity="info" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
              {{ t('payroll.period') }}: {{ periodLabel(run.payroll_period_id) }} · {{ t('payroll.proration_method') }}: {{ prorationLabel(run.proration_method) }}
            </p>
          </div>
          <div class="flex items-center gap-2 shrink-0 flex-wrap">
            <Button :label="t('common.back')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/payroll')" />
            <Button
              v-if="run.status === 'DRAFT'"
              :label="t('payroll.calculate')"
              icon="pi pi-calculator"
              size="small"
              :loading="actionLoading"
              @click="confirmCalculate"
            />
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mt-4">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_employees') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ run.total_employees }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_earning') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ formatMoney(run.total_earning) }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_deduction') }}</p>
            <p class="text-sm font-semibold text-rose-600 dark:text-rose-400 mt-0.5">{{ formatMoney(run.total_deduction) }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_employer_contribution') }}</p>
            <p class="text-sm font-semibold text-amber-600 dark:text-amber-400 mt-0.5">{{ formatMoney(run.total_employer_contribution) }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_net') }}</p>
            <p class="text-sm font-semibold text-emerald-600 dark:text-emerald-400 mt-0.5">{{ formatMoney(run.total_net) }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('payroll.total_company_cost') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ formatMoney(run.total_company_cost) }}</p>
          </div>
        </div>
      </div>

      <!-- ── Tabs ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div class="flex items-center gap-1 px-3 pt-2 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
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

        <!-- ── Overview / Dashboard ── -->
        <div v-if="activeTab === 'overview'" class="p-4">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-chart-bar text-emerald-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('payroll.report_summary') }}</h3>
              </div>
              <div v-if="summaryLoading" class="p-4 space-y-2">
                <div v-for="i in 5" :key="i" class="h-8 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
              </div>
              <div v-else-if="summary" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="row in summaryRows" :key="row.key" class="flex items-center justify-between px-3 py-2.5">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t(row.labelKey) }}</span>
                  <span class="text-sm font-semibold text-gray-800 dark:text-gray-100 font-mono text-xs">{{ row.value }}</span>
                </div>
              </div>
            </div>
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-clock text-indigo-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('common.status') }}</h3>
              </div>
              <div class="p-4 space-y-3">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('common.status') }}</span>
                  <Tag :value="statusLabel(run.status)" :severity="statusSeverity(run.status)" class="!text-xs !px-1.5 !py-0.5" />
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('payroll.calculated') }}</span>
                  <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ run.calculated_at ? formatDateTime(run.calculated_at) : '-' }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('payroll.status_reviewed') }}</span>
                  <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ run.reviewed_at ? formatDateTime(run.reviewed_at) : '-' }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('payroll.status_approved') }}</span>
                  <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ run.approved_at ? formatDateTime(run.approved_at) : '-' }}</span>
                </div>
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-500 dark:text-gray-400">{{ t('payroll.status_locked') }}</span>
                  <span class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ run.locked_at ? formatDateTime(run.locked_at) : '-' }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- ── Employees ── -->
        <div v-if="activeTab === 'employees'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="employees.length" class="text-xs text-gray-400">{{ employees.length }} {{ t('common.items') }}</span>
          </div>
          <DataTable :value="employees" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="employeesLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-users text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('payroll.employees_empty') }}</p>
              </div>
            </template>
            <Column field="employee_code" :header="t('payroll.employee_code')" style="width:120px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.employee_code }}</span></template>
            </Column>
            <Column field="employee_name" :header="t('payroll.employee_name')">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template>
            </Column>
            <Column field="position_title" :header="t('payroll.position_title')" style="width:180px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.position_title || '-' }}</span></template>
            </Column>
            <Column field="total_earning" :header="t('payroll.total_earning')" style="width:140px">
              <template #body="{ data }"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.total_earning) }}</span></template>
            </Column>
            <Column field="total_deduction" :header="t('payroll.total_deduction')" style="width:140px">
              <template #body="{ data }"><span class="text-rose-600 dark:text-rose-400 font-mono text-xs">{{ formatMoney(data.total_deduction) }}</span></template>
            </Column>
            <Column field="net_amount" :header="t('payroll.net_amount')" style="width:140px">
              <template #body="{ data }"><span class="text-emerald-600 dark:text-emerald-400 font-semibold font-mono text-xs">{{ formatMoney(data.net_amount) }}</span></template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Items ── -->
        <div v-if="activeTab === 'items'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="items.length" class="text-xs text-gray-400">{{ items.length }} {{ t('common.items') }}</span>
          </div>
          <DataTable :value="items" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="itemsLoading" :globalFilterFields="['employee_name', 'component_name', 'component_code']">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-list text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('payroll.items_empty') }}</p>
              </div>
            </template>
            <Column field="employee_name" :header="t('payroll.employee_name')" style="width:180px">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeNameFor(data.employee_id) }}</span></template>
            </Column>
            <Column field="component_name" :header="t('payroll.component_name')">
              <template #body="{ data }">
                <span class="text-gray-700 dark:text-gray-200">{{ data.component_name }}</span>
                <span class="text-xs text-gray-400 ml-1">{{ data.component_code }}</span>
              </template>
            </Column>
            <Column field="item_category" :header="t('payroll.item_category')" style="width:130px">
              <template #body="{ data }"><Tag :value="itemCategoryLabel(data.item_category)" :severity="itemCategorySeverity(data.item_category)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="source_group" :header="t('payroll.source_group')" style="width:110px">
              <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ data.source_group }}</span></template>
            </Column>
            <Column field="base_amount" :header="t('payroll.base_amount')" style="width:130px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ formatMoney(data.base_amount) }}</span></template>
            </Column>
            <Column field="rate" :header="t('payroll.rate')" style="width:80px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 text-xs">{{ data.rate != null ? (data.rate * 100).toFixed(2) + '%' : '-' }}</span></template>
            </Column>
            <Column field="amount" :header="t('payroll.amount')" style="width:140px">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-semibold font-mono text-xs">{{ formatMoney(data.amount) }}</span></template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Payslips ── -->
        <div v-if="activeTab === 'payslips'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="payslips.length" class="text-xs text-gray-400">{{ payslips.length }} {{ t('common.items') }}</span>
            <Button :label="t('payroll.generate_payslips')" icon="pi pi-file" size="small" :loading="payslipActionLoading" @click="confirmGeneratePayslips" />
          </div>
          <DataTable :value="payslips" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="payslipsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-file text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('payroll.payslips_empty') }}</p>
              </div>
            </template>
            <Column field="payslip_number" :header="t('payroll.payslip_number')" style="width:160px">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.payslip_number }}</span></template>
            </Column>
            <Column field="employee_name" :header="t('payroll.employee_name')">
              <template #body="{ data }"><span class="text-gray-700 dark:text-gray-200">{{ data.employee_name }}</span></template>
            </Column>
            <Column field="total_earning" :header="t('payroll.total_earning')" style="width:130px">
              <template #body="{ data }"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.total_earning) }}</span></template>
            </Column>
            <Column field="total_deduction" :header="t('payroll.total_deduction')" style="width:130px">
              <template #body="{ data }"><span class="text-rose-600 dark:text-rose-400 font-mono text-xs">{{ formatMoney(data.total_deduction) }}</span></template>
            </Column>
            <Column field="net_amount" :header="t('payroll.net_amount')" style="width:130px">
              <template #body="{ data }"><span class="text-emerald-600 dark:text-emerald-400 font-semibold font-mono text-xs">{{ formatMoney(data.net_amount) }}</span></template>
            </Column>
            <Column field="status" :header="t('common.status')" style="width:110px">
              <template #body="{ data }"><Tag :value="payslipStatusLabel(data.status)" :severity="payslipStatusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
              <template #body="{ data }">
                <div class="flex items-center gap-1 justify-end">
                  <Button icon="pi pi-eye" size="small" text severity="secondary" v-tooltip.left="t('payroll.view')" @click="viewPayslipHTML(data)" />
                  <Button v-if="data.status === 'DRAFT'" icon="pi pi-send" size="small" text severity="success" v-tooltip.left="t('payroll.publish')" @click="publishPayslip(data)" />
                  <Button v-if="data.status === 'PUBLISHED'" icon="pi pi-ban" size="small" text severity="danger" v-tooltip.left="t('payroll.cancel')" @click="cancelPayslip(data)" />
                </div>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Payments ── -->
        <div v-if="activeTab === 'payments'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <div class="flex items-center gap-2 flex-wrap">
              <span v-if="payments.length" class="text-xs text-gray-400">{{ payments.length }} {{ t('common.items') }}</span>
              <span v-if="payments.length" class="text-xs font-medium text-emerald-600 dark:text-emerald-400">{{ formatMoney(sum(payments, 'amount')) }}</span>
            </div>
            <div class="flex items-center gap-2">
              <Button v-if="payments.length" :label="t('payroll.export_csv')" icon="pi pi-download" size="small" severity="secondary" outlined @click="exportPaymentsCSV" />
              <Button :label="t('payroll.create_payment_batch')" icon="pi pi-credit-card" size="small" :loading="paymentActionLoading" @click="confirmCreatePaymentBatch" />
            </div>
          </div>
          <DataTable :value="payments" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="paymentsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-credit-card text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('payroll.payments_empty') }}</p>
              </div>
            </template>
            <Column field="employee_name" :header="t('payroll.employee_name')">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template>
            </Column>
            <Column field="bank_name" :header="t('payroll.bank_name')" style="width:130px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.bank_name || '-' }}</span></template>
            </Column>
            <Column field="bank_account_number" :header="t('payroll.account_number')" style="width:150px">
              <template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ data.bank_account_number }}</span></template>
            </Column>
            <Column field="amount" :header="t('payroll.amount')" style="width:140px">
              <template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-semibold font-mono text-xs">{{ formatMoney(data.amount) }}</span></template>
            </Column>
            <Column field="status" :header="t('common.status')" style="width:110px">
              <template #body="{ data }"><Tag :value="paymentStatusLabel(data.status)" :severity="paymentStatusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="reference" :header="t('payroll.payment_reference')" style="width:150px">
              <template #body="{ data }"><span class="text-gray-500 dark:text-gray-400 text-xs">{{ data.reference || '-' }}</span></template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{ data }">
                <Button
                  v-if="data.status === 'PENDING' || data.status === 'PROCESSING'"
                  :label="t('payroll.update_status')"
                  icon="pi pi-pencil"
                  size="small"
                  text
                  severity="secondary"
                  @click="openPaymentStatusDialog(data)"
                />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Reports ── -->
        <div v-if="activeTab === 'reports'" class="p-4">
          <div class="flex items-center gap-1 mb-3 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
            <button
              v-for="r in reportTabs"
              :key="r.key"
              type="button"
              class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
              :class="activeReport === r.key ? 'text-emerald-600 dark:text-emerald-400 border-b-2 border-emerald-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
              @click="loadReport(r.key)"
            >
              {{ t(r.labelKey) }}
            </button>
          </div>

          <!-- Summary -->
          <div v-if="activeReport === 'summary'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
            <div v-for="row in summaryRows" :key="row.key" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
              <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t(row.labelKey) }}</p>
              <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-1 font-mono">{{ row.value }}</p>
            </div>
          </div>

          <!-- Detail -->
          <DataTable v-if="activeReport === 'detail'" :value="detailRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="reportLoading">
            <template #empty><div class="text-center py-6 text-sm text-gray-400">{{ t('payroll.items_empty') }}</div></template>
            <Column field="employee_name" :header="t('payroll.employee_name')" style="width:180px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template></Column>
            <Column field="component_name" :header="t('payroll.component_name')"><template #body="{ data }"><span class="text-gray-700 dark:text-gray-200">{{ data.component_name }}</span></template></Column>
            <Column field="item_category" :header="t('payroll.item_category')" style="width:120px"><template #body="{ data }"><Tag :value="itemCategoryLabel(data.item_category)" :severity="itemCategorySeverity(data.item_category)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
            <Column field="base_amount" :header="t('payroll.base_amount')" style="width:120px"><template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ formatMoney(data.base_amount) }}</span></template></Column>
            <Column field="amount" :header="t('payroll.amount')" style="width:130px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-semibold font-mono text-xs">{{ formatMoney(data.amount) }}</span></template></Column>
          </DataTable>

          <!-- BPJS -->
          <DataTable v-if="activeReport === 'bpjs'" :value="bpjsRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="reportLoading">
            <template #empty><div class="text-center py-6 text-sm text-gray-400">{{ t('payroll.bpjs_settings_empty') }}</div></template>
            <Column field="employee_name" :header="t('payroll.employee_name')" style="width:180px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template></Column>
            <Column field="bpjs_number" :header="t('payroll.bpjs_number')" style="width:160px"><template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ data.bpjs_number || '-' }}</span></template></Column>
            <Column field="wage_basis" :header="t('payroll.wage_basis')" style="width:130px"><template #body="{ data }"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.wage_basis) }}</span></template></Column>
            <Column field="employee_contribution" :header="t('payroll.employee_contribution')" style="width:130px"><template #body="{ data }"><span class="text-rose-600 dark:text-rose-400 font-mono text-xs">{{ formatMoney(data.employee_contribution) }}</span></template></Column>
            <Column field="employer_contribution" :header="t('payroll.employer_contribution')" style="width:130px"><template #body="{ data }"><span class="text-amber-600 dark:text-amber-400 font-mono text-xs">{{ formatMoney(data.employer_contribution) }}</span></template></Column>
            <Column field="total_contribution" :header="t('payroll.total_contribution')" style="width:130px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-semibold font-mono text-xs">{{ formatMoney(data.total_contribution) }}</span></template></Column>
          </DataTable>

          <!-- Tax -->
          <DataTable v-if="activeReport === 'tax'" :value="taxRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="reportLoading">
            <template #empty><div class="text-center py-6 text-sm text-gray-400">{{ t('payroll.pph21_settings_empty') }}</div></template>
            <Column field="employee_name" :header="t('payroll.employee_name')" style="width:180px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template></Column>
            <Column field="taxable_income" :header="t('payroll.taxable_income')" style="width:160px"><template #body="{ data }"><span class="text-gray-700 dark:text-gray-200 font-mono text-xs">{{ formatMoney(data.taxable_income) }}</span></template></Column>
            <Column field="pph21" :header="t('payroll.pph21_amount')" style="width:140px"><template #body="{ data }"><span class="text-rose-600 dark:text-rose-400 font-semibold font-mono text-xs">{{ formatMoney(data.pph21) }}</span></template></Column>
          </DataTable>

          <!-- Bank -->
          <DataTable v-if="activeReport === 'bank'" :value="bankRows" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="reportLoading">
            <template #empty><div class="text-center py-6 text-sm text-gray-400">{{ t('payroll.payments_empty') }}</div></template>
            <Column field="employee_name" :header="t('payroll.employee_name')" style="width:180px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.employee_name }}</span></template></Column>
            <Column field="bank_name" :header="t('payroll.bank_name')" style="width:130px"><template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.bank_name || '-' }}</span></template></Column>
            <Column field="account_number" :header="t('payroll.account_number')" style="width:150px"><template #body="{ data }"><span class="text-gray-600 dark:text-gray-300 font-mono text-xs">{{ data.account_number }}</span></template></Column>
            <Column field="account_holder_name" :header="t('payroll.account_holder_name')" style="width:170px"><template #body="{ data }"><span class="text-gray-600 dark:text-gray-300">{{ data.account_holder_name }}</span></template></Column>
            <Column field="amount" :header="t('payroll.amount')" style="width:140px"><template #body="{ data }"><span class="text-gray-800 dark:text-gray-100 font-semibold font-mono text-xs">{{ formatMoney(data.amount) }}</span></template></Column>
            <Column field="status" :header="t('common.status')" style="width:110px"><template #body="{ data }"><Tag :value="paymentStatusLabel(data.status)" :severity="paymentStatusSeverity(data.status)" class="!text-xs !px-1.5 !py-0.5" /></template></Column>
          </DataTable>
        </div>
      </div>

      <!-- ── Payment status dialog ── -->
      <Dialog v-model:visible="paymentStatusDialogVisible" :header="t('payroll.update_status')" modal :style="{ width: '440px' }" @hide="resetPaymentStatusForm">
        <div class="space-y-4">
          <FormRow :label="t('common.status')" required :errors="errors?.status">
            <SelectLabel v-model="paymentStatusForm.status" :options="paymentStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.status }" />
          </FormRow>
          <FormRow :label="t('payroll.payment_reference')">
            <TextInput v-model="paymentStatusForm.reference" maxlength="100" :placeholder="t('payroll.payment_reference')" />
          </FormRow>
          <FormRow :label="t('payroll.failed_reason')">
            <TextInput v-model="paymentStatusForm.reason" textarea :rows="2" :placeholder="t('payroll.failed_reason')" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="paymentStatusDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="paymentStatusSaving" :disabled="paymentStatusSaving" @click="handlePaymentStatus" />
          </div>
        </template>
      </Dialog>
    </template>

    <div v-else class="text-center py-12 text-gray-400">
      <i class="pi pi-exclamation-triangle text-3xl mb-2 opacity-50"></i>
      <p class="text-sm">{{ t('common.no_data') }}</p>
    </div>

    <!-- ── Confirmations ── -->
    <ConfirmActionDialog
      v-model:visible="confirmVisible"
      :title="confirmTitle"
      :message="confirmMessage"
      :loading="confirmLoading"
      :confirmLabel="confirmLabel"
      @confirm="onConfirm"
      @cancel="confirmVisible = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import ConfirmActionDialog from '@/components/ConfirmActionDialog.vue'

const { t, locale } = useI18n()
const toast = useToast()
const router = useRouter()
const route = useRoute()

const runId = route.params.id
const run = ref(null)
const loading = ref(true)
const activeTab = ref('overview')
const actionLoading = ref(false)

const periods = ref([])
const employees = ref([])
const employeesLoading = ref(false)
const items = ref([])
const itemsLoading = ref(false)
const payslips = ref([])
const payslipsLoading = ref(false)
const payslipActionLoading = ref(false)
const payments = ref([])
const paymentsLoading = ref(false)
const paymentActionLoading = ref(false)

const summary = ref(null)
const summaryLoading = ref(false)
const detailRows = ref([])
const bpjsRows = ref([])
const taxRows = ref([])
const bankRows = ref([])
const reportLoading = ref(false)
const activeReport = ref('summary')

const paymentStatusDialogVisible = ref(false)
const paymentStatusSaving = ref(false)
const errors = ref({})
const paymentStatusForm = ref({ status: 'PROCESSING', reference: '', reason: '' })
let paymentStatusTarget = null

const confirmVisible = ref(false)
const confirmLoading = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmLabel = ref('')
let confirmAction = null

const tabs = [
  { key: 'overview', labelKey: 'payroll.tab_overview' },
  { key: 'employees', labelKey: 'payroll.tab_employees' },
  { key: 'items', labelKey: 'payroll.tab_items' },
  { key: 'payslips', labelKey: 'payroll.tab_payslips' },
  { key: 'payments', labelKey: 'payroll.tab_payments' },
  { key: 'reports', labelKey: 'payroll.tab_reports' }
]

const reportTabs = [
  { key: 'summary', labelKey: 'payroll.report_summary' },
  { key: 'detail', labelKey: 'payroll.report_detail' },
  { key: 'bpjs', labelKey: 'payroll.report_bpjs' },
  { key: 'tax', labelKey: 'payroll.report_tax' },
  { key: 'bank', labelKey: 'payroll.report_bank' }
]

const summaryRows = computed(() => {
  if (!summary.value) return []
  return [
    { key: 'total_employees', labelKey: 'payroll.total_employees', value: String(summary.value.total_employees) },
    { key: 'gross_salary', labelKey: 'payroll.gross_salary', value: formatMoney(summary.value.gross_salary) },
    { key: 'employee_deduction', labelKey: 'payroll.employee_deduction', value: formatMoney(summary.value.employee_deduction) },
    { key: 'employer_contribution', labelKey: 'payroll.employer_contribution', value: formatMoney(summary.value.employer_contribution) },
    { key: 'net_salary', labelKey: 'payroll.net_salary', value: formatMoney(summary.value.net_salary) },
    { key: 'total_company_cost', labelKey: 'payroll.total_company_cost', value: formatMoney(summary.value.total_company_cost) }
  ]
})

const paymentStatusOptions = computed(() =>
  ['PENDING', 'PROCESSING', 'PAID', 'FAILED', 'REVERSED'].map(v => ({ label: paymentStatusLabel(v), value: v }))
)

function formatMoney(val) {
  const n = Number(val || 0)
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0, maximumFractionDigits: 0 }).format(n)
}
function formatDateTime(val) {
  const d = new Date(val)
  return `${formatDate(d, locale.value)} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}
function sum(list, key) {
  return list.reduce((acc, x) => acc + Number(x[key] || 0), 0)
}
function statusLabel(status) {
  const key = `payroll.status_${String(status || '').toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
function statusSeverity(status) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'CALCULATED': return 'info'
    case 'REVIEWED': return 'warn'
    case 'APPROVED': return 'success'
    case 'LOCKED': return 'contrast'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function itemCategoryLabel(cat) {
  const key = `payroll.item_category_${String(cat || '').toLowerCase()}`
  return t(key) !== key ? t(key) : cat
}
function itemCategorySeverity(cat) {
  switch (cat) {
    case 'EARNING': return 'success'
    case 'DEDUCTION': return 'danger'
    case 'EMPLOYER_CONTRIBUTION': return 'warn'
    default: return 'secondary'
  }
}
function payslipStatusLabel(s) {
  const key = `payroll.payslip_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function payslipStatusSeverity(s) {
  switch (s) {
    case 'DRAFT': return 'secondary'
    case 'PUBLISHED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function paymentStatusLabel(s) {
  const key = `payroll.payment_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function paymentStatusSeverity(s) {
  switch (s) {
    case 'PENDING': return 'secondary'
    case 'PROCESSING': return 'info'
    case 'PAID': return 'success'
    case 'FAILED': return 'danger'
    case 'REVERSED': return 'warn'
    default: return 'secondary'
  }
}
function prorationLabel(m) {
  const key = `payroll.proration_${String(m || '').toLowerCase()}`
  return t(key) !== key ? t(key) : m
}
function periodLabel(id) {
  const p = periods.value.find(x => x.id === id)
  return p ? p.period_code : id
}
function employeeNameFor(id) {
  const e = employees.value.find(x => x.employee_id === id)
  return e ? e.employee_name : id
}

async function loadRun() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}`)
    run.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadReferences() {
  try {
    const res = await api.get('/api/v1/tenant/payroll/periods', { params: { per_page: 200 } })
    periods.value = res.data?.data || []
  } catch { periods.value = [] }
}

async function loadEmployees() {
  employeesLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/employees`)
    employees.value = res.data?.data || []
  } catch { employees.value = [] } finally { employeesLoading.value = false }
}

async function loadItems() {
  itemsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/items`)
    items.value = res.data?.data || []
  } catch { items.value = [] } finally { itemsLoading.value = false }
}

async function loadPayslips() {
  payslipsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/payslips`)
    payslips.value = res.data?.data || []
  } catch { payslips.value = [] } finally { payslipsLoading.value = false }
}

async function loadPayments() {
  paymentsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/payments`)
    payments.value = res.data?.data || []
  } catch { payments.value = [] } finally { paymentsLoading.value = false }
}

async function loadSummary() {
  summaryLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/reports/summary`)
    summary.value = res.data?.data || null
  } catch { summary.value = null } finally { summaryLoading.value = false }
}

async function loadReport(key) {
  activeReport.value = key
  if (key === 'summary') { await loadSummary(); return }
  reportLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/reports/${key}`)
    const data = res.data?.data || []
    if (key === 'detail') detailRows.value = data
    else if (key === 'bpjs') bpjsRows.value = data
    else if (key === 'tax') taxRows.value = data
    else if (key === 'bank') bankRows.value = data
  } catch { } finally { reportLoading.value = false }
}

async function onTabChange(key) {
  activeTab.value = key
  if (key === 'employees' && !employees.value.length) loadEmployees()
  if (key === 'items' && !items.value.length) loadItems()
  if (key === 'payslips') loadPayslips()
  if (key === 'payments') loadPayments()
  if (key === 'reports') loadReport('summary')
}
watch(activeTab, (v) => onTabChange(v))

function confirmCalculate() {
  confirmTitle.value = t('payroll.calculate')
  confirmMessage.value = t('payroll.confirm_calculate')
  confirmLabel.value = t('payroll.calculate')
  confirmAction = calculateRun
  confirmVisible.value = true
}
async function calculateRun() {
  await api.post(`/api/v1/tenant/payroll/runs/${runId}/calculate`)
  toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.run_calculated'), life: 3000 })
  await loadRun()
  loadEmployees()
  loadItems()
  loadSummary()
}

function confirmGeneratePayslips() {
  confirmTitle.value = t('payroll.generate_payslips')
  confirmMessage.value = t('payroll.confirm_generate_payslips')
  confirmLabel.value = t('payroll.generate_payslips')
  confirmAction = generatePayslips
  confirmVisible.value = true
}
async function generatePayslips() {
  payslipActionLoading.value = true
  try {
    await api.post(`/api/v1/tenant/payroll/runs/${runId}/payslips`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.payslips_generated'), life: 3000 })
    await loadPayslips()
  } finally {
    payslipActionLoading.value = false
  }
}

function confirmCreatePaymentBatch() {
  confirmTitle.value = t('payroll.create_payment_batch')
  confirmMessage.value = t('payroll.confirm_create_payment_batch')
  confirmLabel.value = t('payroll.create_payment_batch')
  confirmAction = createPaymentBatch
  confirmVisible.value = true
}
async function createPaymentBatch() {
  paymentActionLoading.value = true
  try {
    const res = await api.post(`/api/v1/tenant/payroll/runs/${runId}/payments`)
    const d = res.data?.data || {}
    toast.add({
      severity: 'success',
      summary: t('message.success'),
      detail: t('payroll.payment_batch_created', {
        total: d.total || 0,
        amount: formatMoney(d.total_amount),
        skipped: d.skipped_no_bank_profile || 0
      }),
      life: 5000
    })
    await loadPayments()
  } finally {
    paymentActionLoading.value = false
  }
}

async function viewPayslipHTML(payslip) {
  try {
    const res = await api.get(`/api/v1/tenant/payroll/payslips/${payslip.id}/html`, { responseType: 'text' })
    const html = typeof res.data === 'string' ? res.data : (res.data?.data || '')
    const win = window.open('', '_blank')
    if (win) {
      win.document.write(html)
      win.document.close()
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

async function publishPayslip(payslip) {
  try {
    await api.post(`/api/v1/tenant/payroll/payslips/${payslip.id}/publish`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.payslip_published'), life: 3000 })
    await loadPayslips()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

async function cancelPayslip(payslip) {
  try {
    await api.post(`/api/v1/tenant/payroll/payslips/${payslip.id}/cancel`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.payslip_cancelled'), life: 3000 })
    await loadPayslips()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

async function exportPaymentsCSV() {
  try {
    const res = await api.get(`/api/v1/tenant/payroll/runs/${runId}/payments/export`, { responseType: 'blob' })
    const url = URL.createObjectURL(new Blob([res.data]))
    const a = document.createElement('a')
    a.href = url
    a.download = 'payments.csv'
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

function openPaymentStatusDialog(payment) {
  paymentStatusTarget = payment
  errors.value = {}
  paymentStatusForm.value = { status: 'PROCESSING', reference: payment.reference || '', reason: '' }
  paymentStatusDialogVisible.value = true
}
function resetPaymentStatusForm() {
  paymentStatusTarget = null
  errors.value = {}
  paymentStatusForm.value = { status: 'PROCESSING', reference: '', reason: '' }
}

async function handlePaymentStatus() {
  errors.value = {}
  if (!paymentStatusForm.value.status) { errors.value = { status: [t('form.required')] }; return }
  paymentStatusSaving.value = true
  try {
    await api.post(`/api/v1/tenant/payroll/payments/${paymentStatusTarget.id}/status`, {
      status: paymentStatusForm.value.status,
      reference: paymentStatusForm.value.reference || '',
      reason: paymentStatusForm.value.reason || ''
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('payroll.run_status_updated'), life: 3000 })
    paymentStatusDialogVisible.value = false
    await loadPayments()
  } catch (e) {
    const fe = getValidationErrors(e)
    if (Object.keys(fe).length > 0) {
      errors.value = fe
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    paymentStatusSaving.value = false
  }
}

async function onConfirm() {
  if (!confirmAction) return
  confirmLoading.value = true
  try {
    await confirmAction()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    confirmLoading.value = false
    confirmVisible.value = false
    confirmAction = null
  }
}

onMounted(() => {
  loadRun()
  loadReferences()
  loadSummary()
})
</script>
