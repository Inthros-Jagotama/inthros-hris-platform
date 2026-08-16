<template>
  <div v-if="loading" class="flex items-center justify-center h-40">
    <i class="pi pi-spinner pi-spin text-2xl text-emerald-500"></i>
  </div>

  <div v-else-if="travel" class="space-y-4">
    <!-- ── Header ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
      <div class="flex items-start justify-between flex-wrap gap-3">
        <div>
          <div class="flex items-center gap-2 flex-wrap">
            <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ travel.title }}</h2>
            <Tag :value="statusLabel(travel.status)" :severity="statusSeverity(travel.status)" class="!text-xs !px-1.5 !py-0.5" />
          </div>
          <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">{{ travel.request_number }} · {{ formatDate(travel.start_date, locale) }} — {{ formatDate(travel.end_date, locale) }}</p>
        </div>
        <div class="flex items-center gap-2">
          <Button v-if="travel.status === 'DRAFT'" :label="t('business_travel.submit')" icon="pi pi-send" size="small" :loading="submitting" @click="handleSubmit" />
          <Button v-if="travel.status === 'DRAFT' || travel.status === 'SUBMITTED'" :label="t('business_travel.cancel_travel')" icon="pi pi-times" size="small" severity="danger" outlined :loading="cancelling" @click="handleCancel" />
        </div>
      </div>
      <p v-if="travel.purpose" class="text-sm text-gray-600 dark:text-gray-300 mt-3">{{ travel.purpose }}</p>
    </div>

    <!-- ── Tabs ── -->
    <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div class="flex items-center gap-1 px-3 pt-2 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
        <button v-for="tab in tabs" :key="tab.key" type="button"
          class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
          :class="activeTab === tab.key ? 'text-indigo-600 dark:text-indigo-400 border-b-2 border-indigo-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
          @click="activeTab = tab.key">
          {{ t(tab.labelKey) }}
        </button>
      </div>

      <!-- ── Info ── -->
      <div v-if="activeTab === 'info'" class="p-4 space-y-5">
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">{{ t('business_travel.participants') }}</h3>
          <div v-if="travel.participants?.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="p in travel.participants" :key="p.id" class="flex items-center justify-between py-2 text-sm">
              <span class="text-gray-700 dark:text-gray-200">{{ p.name || p.employee_id }}</span>
              <div class="flex items-center gap-2">
                <Tag :value="p.role" severity="secondary" class="!text-xs !px-1.5 !py-0.5" />
                <Tag :value="p.participant_type" severity="info" class="!text-xs !px-1.5 !py-0.5" />
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('business_travel.destinations') }}</h3>
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" text @click="openDestinationDialog" />
          </div>
          <div v-if="travel.destinations?.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="d in travel.destinations" :key="d.id" class="py-2 text-sm text-gray-700 dark:text-gray-200">
              {{ [d.city, d.province, d.country].filter(Boolean).join(', ') || d.location || '-' }}
              <span v-if="d.purpose" class="text-gray-400 dark:text-gray-500">— {{ d.purpose }}</span>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('business_travel.activities') }}</h3>
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" text @click="openActivityDialog" />
          </div>
          <div v-if="activities.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="a in activities" :key="a.id" class="py-2 text-sm">
              <span class="text-gray-700 dark:text-gray-200 font-medium">{{ a.title }}</span>
              <span class="text-gray-400 dark:text-gray-500 ml-2">{{ formatDate(a.activity_date, locale) }}</span>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('business_travel.schedules') }}</h3>
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" text @click="openScheduleDialog" />
          </div>
          <div v-if="schedules.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="s in schedules" :key="s.id" class="py-2 text-sm">
              <Tag :value="s.schedule_type" severity="secondary" class="!text-xs !px-1.5 !py-0.5 mr-2" />
              <span class="text-gray-700 dark:text-gray-200">{{ s.origin || '-' }} → {{ s.destination || '-' }}</span>
              <span class="text-gray-400 dark:text-gray-500 ml-2">({{ s.transportation_type }})</span>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>

        <div>
          <div class="flex items-center justify-between mb-2">
            <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('business_travel.documents') }}</h3>
            <Button :label="t('common.add')" icon="pi pi-upload" size="small" text :loading="uploadingTravelDoc" @click="triggerTravelDocUpload" />
          </div>
          <div v-if="travelDocuments.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="d in travelDocuments" :key="d.id" class="py-2 flex items-center justify-between text-sm">
              <a :href="d.file_path" target="_blank" class="text-emerald-600 dark:text-emerald-400 hover:underline">
                <i class="pi pi-paperclip mr-1"></i>{{ d.file_name }}
              </a>
              <div class="flex items-center gap-2">
                <Tag :value="d.document_type" severity="secondary" class="!text-xs !px-1.5 !py-0.5" />
                <Button icon="pi pi-trash" size="small" text severity="danger" @click="handleDeleteTravelDocument(d)" />
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>
      </div>

      <!-- ── Funding ── -->
      <div v-if="activeTab === 'funding'" class="p-4 space-y-3">
        <Message v-if="!canManageFunding" severity="info" :closable="false">{{ t('business_travel.funding_requires_approved') }}</Message>
        <template v-else>
          <div class="flex justify-end">
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" @click="openFundingDialog" />
          </div>
          <div v-if="fundings.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="f in fundings" :key="f.id" class="py-3">
              <div class="flex items-center justify-between text-sm">
                <div>
                  <span class="text-gray-700 dark:text-gray-200 font-medium">{{ fundingMethodName(f.funding_method_id) }}</span>
                  <span class="text-gray-400 dark:text-gray-500 ml-2">{{ formatCurrency(f.amount) }}</span>
                  <Tag :value="f.status" :severity="fundingStatusSeverity(f.status)" class="!text-xs !px-1.5 !py-0.5 ml-2" />
                </div>
                <div class="flex items-center gap-2">
                  <Button v-if="f.status === 'PENDING' || f.status === 'PROCESSING'" :label="t('business_travel.confirm_funding')" size="small" text @click="handleConfirmFunding(f)" />
                  <Button icon="pi pi-paperclip" size="small" text :loading="uploadingDocFor === f.id" @click="triggerFundingUpload(f)" v-tooltip.top="t('business_travel.upload_proof')" />
                </div>
              </div>
              <div v-if="f.documents?.length" class="flex flex-wrap gap-2 mt-1">
                <a v-for="doc in f.documents" :key="doc.id" :href="doc.file_path" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline">
                  <i class="pi pi-paperclip mr-1"></i>{{ doc.file_name }}
                </a>
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </template>
      </div>

      <!-- ── Expenses ── -->
      <div v-if="activeTab === 'expenses'" class="p-4 space-y-3">
        <Message v-if="!canManageFunding" severity="info" :closable="false">{{ t('business_travel.expense_requires_approved') }}</Message>
        <template v-else>
          <div class="flex justify-end">
            <Button :label="t('common.add')" icon="pi pi-plus" size="small" @click="openExpenseDialog" />
          </div>
          <div v-if="expenses.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="e in expenses" :key="e.id" class="py-3">
              <div class="flex items-center justify-between text-sm">
                <div>
                  <span class="text-gray-700 dark:text-gray-200 font-medium">{{ expenseCategoryName(e.expense_category_id) }}</span>
                  <span class="text-gray-400 dark:text-gray-500 ml-2">{{ formatCurrency(e.amount) }}</span>
                  <span class="text-gray-400 dark:text-gray-500 ml-2">{{ formatDate(e.expense_date, locale) }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <Button icon="pi pi-paperclip" size="small" text :loading="uploadingDocFor === e.id" @click="triggerExpenseUpload(e)" v-tooltip.top="t('business_travel.upload_receipt')" />
                  <Button icon="pi pi-trash" size="small" text severity="danger" @click="handleDeleteExpense(e)" />
                </div>
              </div>
              <div v-if="e.documents?.length" class="flex flex-wrap gap-2 mt-1">
                <a v-for="doc in e.documents" :key="doc.id" :href="doc.file_path" target="_blank" class="text-xs text-emerald-600 dark:text-emerald-400 hover:underline">
                  <i class="pi pi-paperclip mr-1"></i>{{ doc.file_name }}
                </a>
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </template>
      </div>

      <!-- ── Settlement ── -->
      <div v-if="activeTab === 'settlement'" class="p-4 space-y-3">
        <Message v-if="travel.status !== 'COMPLETED' && !settlements.length" severity="info" :closable="false">{{ t('business_travel.settlement_requires_completed') }}</Message>
        <div class="flex justify-end" v-if="travel.status === 'COMPLETED'">
          <Button :label="t('business_travel.create_settlement')" icon="pi pi-plus" size="small" @click="handleCreateSettlement" :loading="creatingSettlement" />
        </div>
        <div v-for="s in settlements" :key="s.id" class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-2">
          <div class="flex items-center justify-between">
            <Tag :value="s.status" :severity="settlementStatusSeverity(s.status)" class="!text-xs !px-1.5 !py-0.5" />
            <Button v-if="s.status === 'PENDING'" :label="t('business_travel.submit')" size="small" text @click="handleSubmitSettlement(s)" />
          </div>
          <div class="grid grid-cols-2 sm:grid-cols-3 gap-2 text-xs">
            <div><span class="text-gray-400">{{ t('business_travel.total_advance') }}</span><br><span class="font-medium text-gray-700 dark:text-gray-200">{{ formatCurrency(s.total_advance) }}</span></div>
            <div><span class="text-gray-400">{{ t('business_travel.total_actual_expense') }}</span><br><span class="font-medium text-gray-700 dark:text-gray-200">{{ formatCurrency(s.total_actual_expense) }}</span></div>
            <div><span class="text-gray-400">{{ t('business_travel.total_company_paid') }}</span><br><span class="font-medium text-gray-700 dark:text-gray-200">{{ formatCurrency(s.total_company_paid) }}</span></div>
            <div v-if="s.total_reimbursement > 0"><span class="text-gray-400">{{ t('business_travel.total_reimbursement') }}</span><br><span class="font-medium text-amber-600">{{ formatCurrency(s.total_reimbursement) }}</span></div>
            <div v-if="s.total_refund > 0"><span class="text-gray-400">{{ t('business_travel.total_refund') }}</span><br><span class="font-medium text-emerald-600">{{ formatCurrency(s.total_refund) }}</span></div>
          </div>
        </div>
        <p v-if="!settlements.length" class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
      </div>

      <!-- ── Refund & Reimbursement ── -->
      <div v-if="activeTab === 'refund_reimbursement'" class="p-4 space-y-5">
        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">{{ t('business_travel.refunds') }}</h3>
          <div v-if="refunds.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="r in refunds" :key="r.id" class="py-3 flex items-center justify-between text-sm">
              <div>
                <span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatCurrency(r.refund_amount) }}</span>
                <Tag :value="r.status" :severity="r.status === 'CONFIRMED' ? 'success' : 'warn'" class="!text-xs !px-1.5 !py-0.5 ml-2" />
              </div>
              <Button v-if="r.status === 'PENDING'" :label="t('business_travel.confirm_refund')" size="small" text @click="handleConfirmRefund(r)" />
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>

        <div>
          <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-2">{{ t('business_travel.reimbursements') }}</h3>
          <div v-if="reimbursements.length" class="divide-y divide-gray-100 dark:divide-gray-700">
            <div v-for="r in reimbursements" :key="r.id" class="py-3 flex items-center justify-between text-sm">
              <div>
                <span class="text-gray-700 dark:text-gray-200 font-medium">{{ formatCurrency(r.amount) }}</span>
                <Tag :value="r.status" :severity="reimbursementStatusSeverity(r.status)" class="!text-xs !px-1.5 !py-0.5 ml-2" />
              </div>
              <div class="flex items-center gap-2">
                <Button v-if="r.status === 'REQUESTED'" :label="t('business_travel.approve')" size="small" text @click="handleApproveReimbursement(r)" />
                <Button v-if="r.status === 'APPROVED'" :label="t('business_travel.process')" size="small" text @click="handleProcessReimbursement(r)" />
                <Button v-if="r.status === 'PROCESSING'" :label="t('business_travel.pay')" size="small" text @click="handlePayReimbursement(r)" />
              </div>
            </div>
          </div>
          <p v-else class="text-xs text-gray-400">{{ t('business_travel.empty') }}</p>
        </div>
      </div>
    </div>

    <!-- ── Dialog: Destination ── -->
    <Dialog v-model:visible="destinationDialogVisible" :header="t('business_travel.add_destination')" modal :style="{ width: '460px' }">
      <p class="text-xs text-gray-500 dark:text-gray-400 mb-3 -mt-1">{{ t('business_travel.destination_hint') }}</p>
      <div class="space-y-3">
        <FormRow :label="t('business_travel.city')" required><TextInput v-model="destinationForm.city" /></FormRow>
        <FormRow :label="t('business_travel.province')"><TextInput v-model="destinationForm.province" /></FormRow>
        <FormRow :label="t('business_travel.country')"><TextInput v-model="destinationForm.country" /></FormRow>
        <FormRow :label="t('business_travel.destination_purpose')"><TextInput v-model="destinationForm.purpose" textarea :rows="2" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="destinationDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingDestination" @click="handleSaveDestination" />
      </template>
    </Dialog>

    <!-- ── Dialog: Activity ── -->
    <Dialog v-model:visible="activityDialogVisible" :header="t('business_travel.add_activity')" modal :style="{ width: '460px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.activity_date')" required><DateInput v-model="activityForm.activity_date" /></FormRow>
        <FormRow :label="t('business_travel.field_title')" required><TextInput v-model="activityForm.title" /></FormRow>
        <FormRow :label="t('business_travel.location')"><TextInput v-model="activityForm.location" /></FormRow>
        <FormRow :label="t('business_travel.description')"><TextInput v-model="activityForm.description" textarea :rows="2" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="activityDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingActivity" @click="handleSaveActivity" />
      </template>
    </Dialog>

    <!-- ── Dialog: Schedule ── -->
    <Dialog v-model:visible="scheduleDialogVisible" :header="t('business_travel.add_schedule')" modal :style="{ width: '460px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.schedule_type')" required>
          <Select v-model="scheduleForm.schedule_type" :options="scheduleTypeOptions" class="w-full" />
        </FormRow>
        <FormRow :label="t('business_travel.transportation_type')">
          <Select v-model="scheduleForm.transportation_type" :options="transportationTypeOptions" class="w-full" />
        </FormRow>
        <FormRow :label="t('business_travel.origin')"><TextInput v-model="scheduleForm.origin" /></FormRow>
        <FormRow :label="t('business_travel.destination')"><TextInput v-model="scheduleForm.destination" /></FormRow>
        <FormRow :label="t('business_travel.booking_reference')"><TextInput v-model="scheduleForm.booking_reference" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="scheduleDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingSchedule" @click="handleSaveSchedule" />
      </template>
    </Dialog>

    <!-- ── Dialog: Funding ── -->
    <Dialog v-model:visible="fundingDialogVisible" :header="t('business_travel.add_funding')" modal :style="{ width: '460px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.funding_method')" required>
          <div class="flex items-center gap-2">
            <Select v-model="fundingForm.funding_method_id" :options="fundingMethods" optionLabel="name" optionValue="id" class="flex-1" />
            <Button icon="pi pi-plus" size="small" text @click="quickAddFundingMethodVisible = true" v-tooltip.top="t('business_travel.add_funding_method')" />
          </div>
        </FormRow>
        <FormRow :label="t('business_travel.amount')" required>
          <InputNumber v-model="fundingForm.amount" class="!w-full" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('business_travel.payment_reference')"><TextInput v-model="fundingForm.payment_reference" /></FormRow>
        <FormRow :label="t('business_travel.notes')"><TextInput v-model="fundingForm.notes" textarea :rows="2" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="fundingDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingFunding" @click="handleSaveFunding" />
      </template>
    </Dialog>

    <!-- ── Dialog: Quick-add Funding Method ── -->
    <Dialog v-model:visible="quickAddFundingMethodVisible" :header="t('business_travel.add_funding_method')" modal :style="{ width: '380px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.code')" required><TextInput v-model="fundingMethodForm.code" /></FormRow>
        <FormRow :label="t('business_travel.name')" required><TextInput v-model="fundingMethodForm.name" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="quickAddFundingMethodVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingFundingMethod" @click="handleSaveFundingMethod" />
      </template>
    </Dialog>

    <!-- ── Dialog: Expense ── -->
    <Dialog v-model:visible="expenseDialogVisible" :header="t('business_travel.add_expense')" modal :style="{ width: '460px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.expense_category')" required>
          <div class="flex items-center gap-2">
            <Select v-model="expenseForm.expense_category_id" :options="expenseCategories" optionLabel="name" optionValue="id" class="flex-1" />
            <Button icon="pi pi-plus" size="small" text @click="quickAddExpenseCategoryVisible = true" v-tooltip.top="t('business_travel.add_expense_category')" />
          </div>
        </FormRow>
        <FormRow :label="t('business_travel.expense_date')" required><DateInput v-model="expenseForm.expense_date" /></FormRow>
        <FormRow :label="t('business_travel.amount')" required>
          <InputNumber v-model="expenseForm.amount" class="!w-full" mode="currency" currency="IDR" locale="id-ID" size="small" />
        </FormRow>
        <FormRow :label="t('business_travel.funding_method')">
          <Select v-model="expenseForm.funding_method_id" :options="fundingMethods" optionLabel="name" optionValue="id" showClear class="w-full" />
        </FormRow>
        <FormRow :label="t('business_travel.vendor')"><TextInput v-model="expenseForm.vendor" /></FormRow>
        <FormRow :label="t('business_travel.description')"><TextInput v-model="expenseForm.description" textarea :rows="2" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="expenseDialogVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingExpense" @click="handleSaveExpense" />
      </template>
    </Dialog>

    <!-- ── Dialog: Quick-add Expense Category ── -->
    <Dialog v-model:visible="quickAddExpenseCategoryVisible" :header="t('business_travel.add_expense_category')" modal :style="{ width: '380px' }">
      <div class="space-y-3">
        <FormRow :label="t('business_travel.code')" required><TextInput v-model="expenseCategoryForm.code" /></FormRow>
        <FormRow :label="t('business_travel.name')" required><TextInput v-model="expenseCategoryForm.name" /></FormRow>
      </div>
      <template #footer>
        <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="quickAddExpenseCategoryVisible = false" />
        <Button :label="t('common.save')" size="small" :loading="savingExpenseCategory" @click="handleSaveExpenseCategory" />
      </template>
    </Dialog>

    <input ref="docFileInputRef" type="file" class="hidden" @change="onDocFileSelected" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import { formatDate } from '@/utils/formatDate'
import api from '@/services/api'

import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Select from 'primevue/select'
import InputNumber from 'primevue/inputnumber'
import Message from 'primevue/message'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import DateInput from '@/components/DateInput.vue'

const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()

const travelId = route.params.id
const loading = ref(true)
const travel = ref(null)
const activeTab = ref('info')
const submitting = ref(false)
const cancelling = ref(false)

const activities = ref([])
const schedules = ref([])
const travelDocuments = ref([])
const uploadingTravelDoc = ref(false)
const fundings = ref([])
const fundingMethods = ref([])
const expenses = ref([])
const expenseCategories = ref([])
const settlements = ref([])
const refunds = ref([])
const reimbursements = ref([])

const tabs = [
  { key: 'info', labelKey: 'business_travel.tab_info' },
  { key: 'funding', labelKey: 'business_travel.tab_funding' },
  { key: 'expenses', labelKey: 'business_travel.tab_expenses' },
  { key: 'settlement', labelKey: 'business_travel.tab_settlement' },
  { key: 'refund_reimbursement', labelKey: 'business_travel.tab_refund_reimbursement' }
]

const canManageFunding = computed(() => ['APPROVED', 'IN_PROGRESS', 'COMPLETED'].includes(travel.value?.status))

function statusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'CANCELLED': return 'secondary'
    case 'SUBMITTED': return 'info'
    case 'IN_PROGRESS': return 'warn'
    case 'COMPLETED': return 'info'
    case 'CLOSED': return 'success'
    default: return 'secondary'
  }
}
function statusLabel(status) {
  const key = `business_travel.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
function fundingStatusSeverity(status) {
  return status === 'FUNDED' ? 'success' : status === 'CANCELLED' || status === 'REVERSED' ? 'danger' : 'warn'
}
function settlementStatusSeverity(status) {
  if (status === 'BALANCED' || status === 'SETTLED') return 'success'
  if (status === 'REJECTED') return 'danger'
  return 'warn'
}
function reimbursementStatusSeverity(status) {
  if (status === 'PAID') return 'success'
  if (status === 'REJECTED' || status === 'CANCELLED') return 'danger'
  return 'warn'
}
function formatCurrency(v) {
  if (v === null || v === undefined) return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v)
}
function fundingMethodName(id) {
  return fundingMethods.value.find(m => m.id === id)?.name || id?.slice(0, 8) || '-'
}
function expenseCategoryName(id) {
  return expenseCategories.value.find(c => c.id === id)?.name || id?.slice(0, 8) || '-'
}

const scheduleTypeOptions = ['DEPARTURE', 'RETURN', 'TRANSFER', 'OTHER']
const transportationTypeOptions = ['AIRPLANE', 'TRAIN', 'BUS', 'CAR', 'COMPANY_CAR', 'RENTAL_CAR', 'MOTORCYCLE', 'OTHER']

async function loadTravel() {
  const res = await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}`)
  travel.value = res.data?.data || null
}
async function loadActivities() {
  try { activities.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/activities`)).data?.data || [] } catch { activities.value = [] }
}
async function loadSchedules() {
  try { schedules.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/schedules`)).data?.data || [] } catch { schedules.value = [] }
}
async function loadTravelDocuments() {
  try { travelDocuments.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/documents`)).data?.data || [] } catch { travelDocuments.value = [] }
}
async function loadFundingMethods() {
  try { fundingMethods.value = (await api.get('/api/v1/tenant/attendance/business-travel-funding-methods')).data?.data || [] } catch { fundingMethods.value = [] }
}
async function loadFundings() {
  try { fundings.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/fundings`)).data?.data || [] } catch { fundings.value = [] }
}
async function loadExpenseCategories() {
  try { expenseCategories.value = (await api.get('/api/v1/tenant/attendance/business-travel-expense-categories')).data?.data || [] } catch { expenseCategories.value = [] }
}
async function loadExpenses() {
  try { expenses.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/expenses`)).data?.data || [] } catch { expenses.value = [] }
}
async function loadSettlements() {
  try { settlements.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/settlements`)).data?.data || [] } catch { settlements.value = [] }
}
async function loadRefunds() {
  try { refunds.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/refunds`)).data?.data || [] } catch { refunds.value = [] }
}
async function loadReimbursements() {
  try { reimbursements.value = (await api.get(`/api/v1/tenant/attendance/business-travels/${travelId}/reimbursements`)).data?.data || [] } catch { reimbursements.value = [] }
}

async function loadAll() {
  loading.value = true
  try {
    await loadTravel()
    await Promise.all([
      loadActivities(), loadSchedules(), loadTravelDocuments(), loadFundingMethods(), loadFundings(),
      loadExpenseCategories(), loadExpenses(), loadSettlements(), loadRefunds(), loadReimbursements()
    ])
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

// ── Travel actions ──
async function handleSubmit() {
  submitting.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/submit`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadTravel()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    submitting.value = false
  }
}
async function handleCancel() {
  cancelling.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/cancel`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadTravel()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    cancelling.value = false
  }
}

// ── Destination ──
const destinationDialogVisible = ref(false)
const savingDestination = ref(false)
const destinationForm = ref({ city: '', province: '', country: '', purpose: '' })
function openDestinationDialog() {
  destinationForm.value = { city: '', province: '', country: '', purpose: '' }
  destinationDialogVisible.value = true
}
async function handleSaveDestination() {
  if (!destinationForm.value.city?.trim()) return
  savingDestination.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/destinations`, destinationForm.value)
    destinationDialogVisible.value = false
    await loadTravel()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingDestination.value = false
  }
}

// ── Activity ──
const activityDialogVisible = ref(false)
const savingActivity = ref(false)
const activityForm = ref({ activity_date: '', title: '', location: '', description: '' })
function openActivityDialog() {
  activityForm.value = { activity_date: '', title: '', location: '', description: '' }
  activityDialogVisible.value = true
}
async function handleSaveActivity() {
  if (!activityForm.value.activity_date || !activityForm.value.title?.trim()) return
  savingActivity.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/activities`, activityForm.value)
    activityDialogVisible.value = false
    await loadActivities()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingActivity.value = false
  }
}

// ── Schedule ──
const scheduleDialogVisible = ref(false)
const savingSchedule = ref(false)
const scheduleForm = ref({ schedule_type: 'DEPARTURE', transportation_type: 'OTHER', origin: '', destination: '', booking_reference: '' })
function openScheduleDialog() {
  scheduleForm.value = { schedule_type: 'DEPARTURE', transportation_type: 'OTHER', origin: '', destination: '', booking_reference: '' }
  scheduleDialogVisible.value = true
}
async function handleSaveSchedule() {
  savingSchedule.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/schedules`, scheduleForm.value)
    scheduleDialogVisible.value = false
    await loadSchedules()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingSchedule.value = false
  }
}

// ── Funding ──
const fundingDialogVisible = ref(false)
const savingFunding = ref(false)
const fundingForm = ref({ funding_method_id: '', amount: null, payment_reference: '', notes: '' })
function openFundingDialog() {
  fundingForm.value = { funding_method_id: '', amount: null, payment_reference: '', notes: '' }
  fundingDialogVisible.value = true
}
async function handleSaveFunding() {
  const errs = getValidationErrors
  if (!fundingForm.value.funding_method_id || !fundingForm.value.amount) return
  savingFunding.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/fundings`, fundingForm.value)
    fundingDialogVisible.value = false
    await loadFundings()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingFunding.value = false
  }
}
async function handleConfirmFunding(f) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/fundings/${f.id}/confirm`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadFundings()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

// ── Documents (transfer proof / receipt) ──
// Two-step: upload raw file to the generic upload endpoint to get a URL,
// then attach that URL to the funding/expense via its documents endpoint —
// same pattern as AttendanceOvertime.vue's attachment_url upload.
const docFileInputRef = ref(null)
const docUploadTarget = ref(null) // { type: 'funding' | 'expense', id }
const uploadingDocFor = ref(null)

function triggerFundingUpload(f) {
  docUploadTarget.value = { type: 'funding', id: f.id }
  docFileInputRef.value?.click()
}
function triggerExpenseUpload(e) {
  docUploadTarget.value = { type: 'expense', id: e.id }
  docFileInputRef.value?.click()
}
function triggerTravelDocUpload() {
  docUploadTarget.value = { type: 'travel', id: null }
  docFileInputRef.value?.click()
}

async function onDocFileSelected(event) {
  const file = event.target.files?.[0]
  const target = docUploadTarget.value
  if (!file || !target) return
  if (target.type === 'travel') uploadingTravelDoc.value = true
  else uploadingDocFor.value = target.id
  try {
    const fd = new FormData()
    fd.append('file', file)
    const uploadRes = await api.post('/api/v1/tenant/uploads', fd, { headers: { 'Content-Type': 'multipart/form-data' } })
    const filePath = uploadRes.data?.data?.url || ''
    if (!filePath) throw new Error('upload failed')

    const documentType = target.type === 'funding' ? 'TRANSFER_RECEIPT' : target.type === 'expense' ? 'RECEIPT' : 'OTHER'
    const docPayload = { document_type: documentType, file_name: file.name, file_path: filePath, mime_type: file.type, file_size: file.size }
    if (target.type === 'funding') {
      await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/fundings/${target.id}/documents`, docPayload)
      await loadFundings()
    } else if (target.type === 'expense') {
      await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/expenses/${target.id}/documents`, docPayload)
      await loadExpenses()
    } else {
      await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/documents`, docPayload)
      await loadTravelDocuments()
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    uploadingTravelDoc.value = false
    uploadingDocFor.value = null
    docUploadTarget.value = null
    if (docFileInputRef.value) docFileInputRef.value.value = ''
  }
}

const quickAddFundingMethodVisible = ref(false)
const savingFundingMethod = ref(false)
const fundingMethodForm = ref({ code: '', name: '' })
async function handleSaveFundingMethod() {
  if (!fundingMethodForm.value.code?.trim() || !fundingMethodForm.value.name?.trim()) return
  savingFundingMethod.value = true
  try {
    await api.post('/api/v1/tenant/attendance/business-travel-funding-methods', fundingMethodForm.value)
    quickAddFundingMethodVisible.value = false
    fundingMethodForm.value = { code: '', name: '' }
    await loadFundingMethods()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingFundingMethod.value = false
  }
}

// ── Expense ──
const expenseDialogVisible = ref(false)
const savingExpense = ref(false)
const expenseForm = ref({ expense_category_id: '', expense_date: '', amount: null, funding_method_id: '', vendor: '', description: '' })
function openExpenseDialog() {
  expenseForm.value = { expense_category_id: '', expense_date: '', amount: null, funding_method_id: '', vendor: '', description: '' }
  expenseDialogVisible.value = true
}
async function handleSaveExpense() {
  if (!expenseForm.value.expense_category_id || !expenseForm.value.expense_date || !expenseForm.value.amount) return
  savingExpense.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/expenses`, expenseForm.value)
    expenseDialogVisible.value = false
    await loadExpenses()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingExpense.value = false
  }
}
async function handleDeleteExpense(e) {
  try {
    await api.delete(`/api/v1/tenant/attendance/business-travels/${travelId}/expenses/${e.id}`)
    await loadExpenses()
  } catch (err) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(err, t('message.operation_failed')), life: 4000 })
  }
}

async function handleDeleteTravelDocument(d) {
  try {
    await api.delete(`/api/v1/tenant/attendance/business-travels/${travelId}/documents/${d.id}`)
    await loadTravelDocuments()
  } catch (err) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(err, t('message.operation_failed')), life: 4000 })
  }
}

const quickAddExpenseCategoryVisible = ref(false)
const savingExpenseCategory = ref(false)
const expenseCategoryForm = ref({ code: '', name: '' })
async function handleSaveExpenseCategory() {
  if (!expenseCategoryForm.value.code?.trim() || !expenseCategoryForm.value.name?.trim()) return
  savingExpenseCategory.value = true
  try {
    await api.post('/api/v1/tenant/attendance/business-travel-expense-categories', expenseCategoryForm.value)
    quickAddExpenseCategoryVisible.value = false
    expenseCategoryForm.value = { code: '', name: '' }
    await loadExpenseCategories()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    savingExpenseCategory.value = false
  }
}

// ── Settlement ──
const creatingSettlement = ref(false)
async function handleCreateSettlement() {
  creatingSettlement.value = true
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/settlements`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadSettlements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    creatingSettlement.value = false
  }
}
async function handleSubmitSettlement(s) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/settlements/${s.id}/submit`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadSettlements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

// ── Refund / Reimbursement ──
async function handleConfirmRefund(r) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/refunds/${r.id}/confirm`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadRefunds()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}
async function handleApproveReimbursement(r) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/reimbursements/${r.id}/approve`)
    await loadReimbursements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}
async function handleProcessReimbursement(r) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/reimbursements/${r.id}/process`)
    await loadReimbursements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}
async function handlePayReimbursement(r) {
  try {
    await api.post(`/api/v1/tenant/attendance/business-travels/${travelId}/reimbursements/${r.id}/pay`, {})
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    await loadReimbursements()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

onMounted(loadAll)
</script>
