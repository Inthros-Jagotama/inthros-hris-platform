<template>
  <div class="space-y-1">
    <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
      <div class="flex items-center gap-2">
        <span v-if="pendingTotal > 0" class="text-xs text-gray-400 dark:text-gray-500">
          {{ pendingTotal }} {{ t('common.items') }}
        </span>
        <Button icon="pi pi-refresh" size="small" text severity="secondary" @click="loadTasks" />
      </div>
      <Button v-if="hasPermission('approval.settings.view')" :label="t('approval.flows')" icon="pi pi-sitemap" size="small" @click="router.push({ name: 'ApprovalFlows' })" />
    </div>

    <!-- Tab: Perlu Tindakan / Riwayat -->
    <div class="flex items-center gap-1 mb-2 border-b border-gray-200 dark:border-gray-700">
      <button
        class="px-3 py-1.5 text-sm font-medium border-b-2 -mb-px transition-colors"
        :class="activeTab === 'pending' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab('pending')"
      >
        {{ t('approval.tab_pending') }}
      </button>
      <button
        class="px-3 py-1.5 text-sm font-medium border-b-2 -mb-px transition-colors"
        :class="activeTab === 'done' ? 'border-emerald-500 text-emerald-600 dark:text-emerald-400' : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
        @click="switchTab('done')"
      >
        {{ t('approval.tab_history') }}
      </button>
    </div>

    <!-- Filter status & alur -->
    <div class="flex items-center gap-2 flex-wrap mb-2">
      <Select
        v-if="activeTab === 'pending'"
        v-model="statusFilter"
        :options="statusOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('approval.filter_status')"
        class="!w-44"
        size="small"
        showClear
        @update:modelValue="onFilterChange"
      />
      <Select
        v-model="flowFilter"
        :options="flowOptions"
        optionLabel="label"
        optionValue="value"
        :placeholder="t('approval.filter_flow')"
        class="!w-56"
        size="small"
        showClear
        :loading="flowsLoading"
        @update:modelValue="onFilterChange"
      />
      <Button
        v-if="statusFilter || flowFilter"
        :label="t('approval.filter_reset')"
        icon="pi pi-filter-slash"
        size="small"
        text
        severity="secondary"
        @click="clearFilters"
      />
    </div>

    <SkeletonTable v-if="tasksLoading" :columns="taskSkeletonColumns" :rows="6" />

    <DataTable
      v-else
      :value="pendingTasks"
      size="small"
      class="!text-sm p-datatable-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
      lazy
      :totalRecords="pendingTotal"
      :first="firstRecord"
      :rows="perPage"
      @page="onPage($event)"
      paginator
      paginatorTemplate="FirstPageLink PrevPageLink PageLinks NextPageLink LastPageLink RowsPerPageDropdown"
      :rowsPerPageOptions="[10, 20, 50]"
    >
      <template #empty>
        <div class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
          <i class="pi pi-inbox text-3xl mb-2 opacity-50"></i>
          <p class="text-sm font-medium">{{ t(activeTab === 'pending' ? 'approval.empty_tasks' : 'approval.empty_history') }}</p>
        </div>
      </template>
      <Column field="step_order" :header="t('approval.step')" style="width:80px">
        <template #body="{data}">#{{ data.step_order }}</template>
      </Column>
      <Column field="flow_name" :header="t('approval.instance')">
        <template #body="{data}">
          <span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.flow_name || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_name" :header="t('approval.submitted_by')">
        <template #body="{data}">
          <span class="text-gray-700 dark:text-gray-200">{{ data.submitter_name || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_employee_code" :header="t('approval.employee_id')" style="width:120px">
        <template #body="{data}">
          <span class="text-xs text-gray-500 dark:text-gray-400 font-mono">{{ data.submitter_employee_code || '-' }}</span>
        </template>
      </Column>
      <Column field="submitter_organization_name" :header="t('approval.organization')" style="width:160px">
        <template #body="{data}">
          <span class="text-xs text-gray-600 dark:text-gray-300">{{ data.submitter_organization_name || '-' }}</span>
        </template>
      </Column>
      <Column field="status" :header="t('common.status')" style="width:120px">
        <template #body="{data}">
          <Tag :value="rowStatusLabel(data)" :severity="rowStatusSeverity(data)" class="!text-xs" />
        </template>
      </Column>
      <Column field="created_at" :header="t('approval.submitted_at')" style="width:160px">
        <template #body="{data}">
          <span class="text-xs text-gray-500 dark:text-gray-400">{{ formatDate(data.created_at) }}</span>
        </template>
      </Column>
      <Column :header="t('common.actions')" style="width:120px" frozen alignFrozen="right">
        <template #body="{data}">
          <Button :label="t('approval.review')" icon="pi pi-eye" size="small" text @click="openTaskDetail(data)" />
        </template>
      </Column>
    </DataTable>

    <!-- ===================================================================== -->
    <!-- Task detail / act dialog -->
    <!-- ===================================================================== -->
    <Dialog v-model:visible="taskDetailVisible" :header="t('approval.instance_detail')" modal :style="{ width: '960px' }">
      <div v-if="instanceLoading" class="py-8 text-center text-gray-400 text-sm">{{ t('common.loading') }}</div>
      <div v-else-if="activeInstance" class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <!-- Left column: the data being submitted -->
        <div class="space-y-4 md:border-r md:border-gray-200 md:dark:border-gray-700 md:pr-5">
          <div v-if="!isAttendanceModule" class="flex items-center gap-2">
            <i class="pi pi-file text-indigo-400 text-sm"></i>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('approval.submitted_data') }}</h2>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>

          <div v-if="documentLoading" class="py-6 text-center text-gray-400 text-sm">{{ t('common.loading') }}</div>

          <template v-else-if="isKPIModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('kpi.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.employee_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('kpi.organization') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.organization_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('kpi.period') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.period_code || '-' }}</p>
              </div>
            </div>

            <div v-if="documentDetail?.details?.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.indicators') }}</p>
              <div class="space-y-1">
                <div v-for="d in documentDetail.details" :key="d.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                  <div class="font-medium text-gray-700 dark:text-gray-200">{{ d.indicator_name }}</div>
                  <div class="text-gray-500 dark:text-gray-400">
                    {{ t('kpi.target') }}: {{ d.target }} {{ d.unit_of_measurement }}
                    <span v-if="d.actual"> · {{ t('kpi.actual') }}: {{ d.actual }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="documentDetail?.program_items?.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('kpi.program_items') }}</p>
              <div class="space-y-1">
                <div v-for="p in documentDetail.program_items" :key="p.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                  <div class="font-medium text-gray-700 dark:text-gray-200">{{ p.title }}</div>
                  <div class="text-gray-500 dark:text-gray-400">
                    {{ t('kpi.target') }}: {{ p.target }} {{ p.unit_of_measurement }}
                    <span v-if="p.actual"> · {{ t('kpi.actual') }}: {{ p.actual }}</span>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="isOKRModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('okr.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.employee_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('okr.organization') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.organization_name || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('okr.period') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.period_code || '-' }}</p>
              </div>
            </div>

            <div v-if="okrObjectiveGroups.length">
              <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mt-3 mb-1">{{ t('okr.objectives') }}</p>
              <div v-for="g in okrObjectiveGroups" :key="g.key" class="mb-2">
                <div class="text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{{ g.title }}</div>
                <div class="space-y-1">
                  <div v-for="kr in g.items" :key="kr.id" class="text-xs px-2 py-1.5 rounded bg-gray-50 dark:bg-gray-800/60">
                    <div class="font-medium text-gray-700 dark:text-gray-200">{{ kr.key_result_title }}</div>
                    <div class="text-gray-500 dark:text-gray-400">
                      {{ t('okr.target') }}: {{ kr.target_value }} {{ kr.unit }}
                      <span v-if="kr.actual_value"> · {{ t('okr.actual') }}: {{ kr.actual_value }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="isAttendanceModule">
            <div class="space-y-4">
              <!-- Group: Informasi Lembur -->
              <div>
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-clock text-indigo-400 text-sm"></i>
                  <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('attendance.overtime_info') }}</h2>
                  <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
                </div>
                <div class="grid grid-cols-2 gap-3 text-sm">
                  <div class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.employee') }}</p>
                    <p class="text-gray-800 dark:text-gray-100 font-medium">{{ overtimeEmployeeName || '-' }}</p>
                    <p v-if="overtimeEmployeeCode" class="text-xs text-gray-400">{{ overtimeEmployeeCode }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.flow_type') }}</p>
                    <Tag :value="documentDetail?.flow_type === 'ASSIGNED' ? t('attendance.flow_assigned') : t('attendance.flow_self')" :severity="documentDetail?.flow_type === 'ASSIGNED' ? 'info' : 'secondary'" class="!text-xs" />
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
                    <Tag :value="overtimeStatusLabel(documentDetail?.status)" :severity="overtimeStatusSeverity(documentDetail?.status)" class="!text-xs" />
                  </div>
                  <div v-if="isAssignedFlow" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.assigned_by') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ overtimeAssignedByName || '-' }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.work_date') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.work_date) }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.requested_minutes') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">
                      {{ documentDetail?.requested_minutes ?? '-' }} min
                      <span v-if="isOvertimeCrossDay" class="text-xs text-amber-500 ml-1">· {{ t('attendance.overtime_cross_day') }}</span>
                    </p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.start_time') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatTime(documentDetail?.start_time_local) }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.end_time') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatTime(documentDetail?.end_time_local) }}</p>
                  </div>
                  <div v-if="documentDetail?.reason" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.reason') }}</p>
                    <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.reason }}</p>
                  </div>
                  <div v-if="documentDetail?.approval_note" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.approval_note') }}</p>
                    <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.approval_note }}</p>
                  </div>
                  <div class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('approval.submitted_at') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatDate(documentDetail?.created_at) }}</p>
                  </div>
                </div>
              </div>

              <!-- Group: Detail Aktual -->
              <div v-if="hasActualData">
                <div class="flex items-center gap-2 mb-2">
                  <i class="pi pi-check-square text-indigo-400 text-sm"></i>
                  <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('attendance.actual_data') }}</h2>
                  <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
                </div>
                <div class="grid grid-cols-2 gap-3 text-sm">
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.actual_start_time') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatTime(documentDetail?.actual_start_time_local) }}</p>
                  </div>
                  <div>
                    <p class="text-xs text-gray-400">{{ t('attendance.actual_end_time') }}</p>
                    <p class="text-gray-700 dark:text-gray-200">{{ formatTime(documentDetail?.actual_end_time_local) }}</p>
                  </div>
                  <div v-if="documentDetail?.actual_note" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.actual_note') }}</p>
                    <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.actual_note }}</p>
                  </div>
                  <div v-if="documentDetail?.calculated_minutes != null" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.calculated_minutes') }}</p>
                    <p class="text-gray-700 dark:text-gray-200 font-medium">{{ documentDetail.calculated_minutes }} min</p>
                  </div>
                  <div v-if="documentDetail?.attachment_url" class="col-span-2">
                    <p class="text-xs text-gray-400">{{ t('attendance.attachment') }}</p>
                    <a :href="documentDetail.attachment_url" target="_blank" class="text-sm text-emerald-600 dark:text-emerald-400 hover:underline">
                      <i class="pi pi-paperclip mr-1"></i>{{ t('attendance.attachment') }}
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </template>

          <template v-else-if="isLeaveModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ leaveEmployeeName || '-' }}</p>
                <p v-if="leaveEmployeeCode" class="text-xs text-gray-400">{{ leaveEmployeeCode }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.leave_type') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ leaveTypeName || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.request_start_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.request_start_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.request_end_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.request_end_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.duration_mode') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.duration_mode || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('leave.requested_days') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.requested_days ?? '-' }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('approval.submitted_at') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDate(documentDetail?.submitted_at) }}</p>
              </div>
              <div v-if="documentDetail?.leave_reason_note" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('leave.leave_reason_note') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.leave_reason_note }}</p>
              </div>
            </div>
          </template>

          <template v-else-if="isEmployeemovementModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('employee_movement.employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.employee_name || '-' }}</p>
                <p v-if="documentDetail?.employee_code" class="text-xs text-gray-400">{{ documentDetail.employee_code }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('employee_movement.movement_type') }}</p>
                <Tag :value="movementTypeLabel(documentDetail?.movement_type)" :severity="movementTypeSeverity(documentDetail?.movement_type)" class="!text-xs" />
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
                <Tag :value="movementStatusLabel(documentDetail?.status)" :severity="movementStatusSeverity(documentDetail?.status)" class="!text-xs" />
              </div>
              <div v-if="documentDetail?.decision_letter_number" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('employee_movement.decision_letter_number') }}</p>
                <p class="text-gray-700 dark:text-gray-200 font-mono text-xs break-all">{{ documentDetail.decision_letter_number }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('employee_movement.decision_letter_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.decision_letter_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('employee_movement.effective_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail?.effective_date) }}</p>
              </div>

              <!-- Dari → Ke (enrichment G-4 / snapshot §12.5) -->
              <div v-if="hasMovementFrom" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3 space-y-1.5">
                <p class="text-xs uppercase tracking-wide text-gray-400 font-medium">{{ t('employee_movement.from') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.from_organization_name || '-' }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.from_position_name || '-' }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.from_employment_status_name || '-' }}</p>
              </div>
              <div v-if="hasMovementTo" class="rounded-lg border border-emerald-200 dark:border-emerald-900/40 p-3 space-y-1.5">
                <p class="text-xs uppercase tracking-wide text-emerald-500 font-medium">{{ t('employee_movement.to') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.to_organization_name || '-' }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.to_position_name || '-' }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.to_employment_status_name || '-' }}</p>
              </div>

              <div v-if="documentDetail?.reason" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('employee_movement.reason') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.reason }}</p>
              </div>
              <div v-if="documentDetail?.notes" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('employee_movement.notes') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.notes }}</p>
              </div>
            </div>
          </template>

          <!-- Detail lowongan (recruitment) -->
          <template v-else-if="isRecruitmentModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('requisitions.title') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.title || '-' }}</p>
                <p v-if="documentDetail?.requisition_number" class="text-xs text-gray-400 font-mono">{{ documentDetail.requisition_number }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
                <Tag :value="requisitionStatusLabel(documentDetail?.status)" :severity="requisitionStatusSeverity(documentDetail?.status)" class="!text-xs" />
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('requisitions.priority') }}</p>
                <Tag :value="requisitionPriorityLabel(documentDetail?.priority)" :severity="requisitionPrioritySeverity(documentDetail?.priority)" class="!text-xs" />
              </div>
              <div v-if="requisitionOrganizationName" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('requisitions.org') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ requisitionOrganizationName }}</p>
              </div>
              <div v-if="documentDetail?.department">
                <p class="text-xs text-gray-400">{{ t('requisitions.department') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail.department }}</p>
              </div>
              <div v-if="documentDetail?.employment_type">
                <p class="text-xs text-gray-400">{{ t('requisitions.employment_type') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail.employment_type }}</p>
              </div>
              <div v-if="documentDetail?.location">
                <p class="text-xs text-gray-400">{{ t('requisitions.location') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail.location }}</p>
              </div>
              <div v-if="documentDetail?.reason_type">
                <p class="text-xs text-gray-400">{{ t('requisitions.reason_type') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ requisitionReasonLabel(documentDetail.reason_type) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('requisitions.slots_available') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.slots_available ?? '-' }}</p>
              </div>
              <div v-if="documentDetail?.slots_filled != null">
                <p class="text-xs text-gray-400">{{ t('requisitions.slots') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail.slots_filled }}</p>
              </div>
              <div v-if="documentDetail?.target_start_date">
                <p class="text-xs text-gray-400">{{ t('requisitions.target_start_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail.target_start_date) }}</p>
              </div>
              <div v-if="documentDetail?.min_salary || documentDetail?.max_salary">
                <p class="text-xs text-gray-400">{{ t('requisitions.min_salary') }} – {{ t('requisitions.max_salary') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatCurrency(documentDetail?.min_salary) }} – {{ formatCurrency(documentDetail?.max_salary) }}</p>
              </div>
              <div v-if="documentDetail?.description" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('requisitions.description_label') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words whitespace-pre-line">{{ documentDetail.description }}</p>
              </div>
              <div v-if="documentDetail?.requirements" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('requisitions.requirements') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words whitespace-pre-line">{{ documentDetail.requirements }}</p>
              </div>
              <div v-if="documentDetail?.responsibilities" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('requisitions.responsibilities') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words whitespace-pre-line">{{ documentDetail.responsibilities }}</p>
              </div>
            </div>
          </template>

          <!-- Detail penawaran kandidat (recruitment_offer) -->
          <template v-else-if="isOfferModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('offers.offer_number') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ documentDetail?.offer_number || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
                <Tag :value="offerStatusLabel(documentDetail?.status)" :severity="offerStatusSeverity(documentDetail?.status)" class="!text-xs" />
              </div>
              <div v-if="documentDetail?.employment_type">
                <p class="text-xs text-gray-400">{{ t('offers.employment_type') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail.employment_type }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('offers.salary') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatCurrency(documentDetail?.salary) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('offers.allowances') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatCurrency(documentDetail?.allowances) }}</p>
              </div>
              <div v-if="documentDetail?.start_date">
                <p class="text-xs text-gray-400">{{ t('offers.start_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail.start_date) }}</p>
              </div>
              <div v-if="documentDetail?.expiry_date">
                <p class="text-xs text-gray-400">{{ t('offers.expiry_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail.expiry_date) }}</p>
              </div>
              <div v-if="documentDetail?.benefits" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('offers.benefits') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words whitespace-pre-line">{{ documentDetail.benefits }}</p>
              </div>
            </div>
          </template>

          <!-- Detail permintaan training (training_request) -->
          <template v-else-if="isTrainingRequestModule">
            <div class="grid grid-cols-2 gap-3 text-sm">
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('training.request_employee') }}</p>
                <p class="text-gray-800 dark:text-gray-100 font-medium">{{ trainingRequestEmployeeName || '-' }}</p>
              </div>
              <div class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('training.request_course') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ trainingRequestCourseName || '-' }}</p>
              </div>
              <div v-if="documentDetail?.requested_date">
                <p class="text-xs text-gray-400">{{ t('training.request_requested_date') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ formatDateOnly(documentDetail.requested_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('training.request_priority') }}</p>
                <p class="text-gray-700 dark:text-gray-200">{{ documentDetail?.priority || '-' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
                <Tag :value="trainingRequestStatusLabel(documentDetail?.status)" :severity="trainingRequestStatusSeverity(documentDetail?.status)" class="!text-xs" />
              </div>
              <div v-if="documentDetail?.reason" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('training.request_reason') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.reason }}</p>
              </div>
              <div v-if="documentDetail?.supervisor_note" class="col-span-2">
                <p class="text-xs text-gray-400">{{ t('training.request_supervisor_note') }}</p>
                <p class="text-gray-700 dark:text-gray-200 break-words">{{ documentDetail.supervisor_note }}</p>
              </div>
            </div>
          </template>

          <div v-else-if="documentFields.length" class="space-y-2">
            <div v-for="f in documentFields" :key="f.label" class="text-sm">
              <p class="text-xs text-gray-400">{{ f.label }}</p>
              <p class="text-gray-700 dark:text-gray-200 break-words">{{ f.value }}</p>
            </div>
          </div>

          <p v-else class="text-xs text-gray-400">{{ t('approval.submitted_data_unavailable') }}</p>
        </div>

        <!-- Right column: existing approval data -->
        <div class="space-y-4">
          <div class="flex items-center gap-2">
            <i class="pi pi-check-circle text-indigo-400 text-sm"></i>
            <h2 class="text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-gray-400">{{ t('approval.approval_data') }}</h2>
            <div class="flex-1 border-t border-gray-200 dark:border-gray-700"></div>
          </div>

          <div class="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p class="text-xs text-gray-400">{{ t('approval.module') }}</p>
              <Tag :value="moduleLabel(activeInstance.module)" severity="info" class="!text-xs" />
            </div>
            <div>
              <p class="text-xs text-gray-400">{{ t('common.status') }}</p>
              <Tag :value="activeInstance.status" severity="warn" class="!text-xs" />
            </div>
            <div class="col-span-2">
              <p class="text-xs text-gray-400">{{ t('approval.flow_name') }}</p>
              <p class="text-gray-700 dark:text-gray-200">{{ activeInstance.flow_name || '-' }}</p>
            </div>
          </div>

          <div>
            <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-2">{{ t('approval.steps') }}</p>
            <div class="space-y-1">
              <div v-for="s in activeInstance.steps" :key="s.id" class="px-2 py-1.5 rounded" :class="s.step_order === activeInstance.current_step ? 'bg-emerald-50 dark:bg-emerald-500/10' : ''">
                <div class="flex items-center gap-2 text-xs">
                  <span class="w-5 h-5 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center shrink-0">{{ s.step_order }}</span>
                  <span class="font-medium text-gray-700 dark:text-gray-200">{{ s.step_name }}</span>
                  <Tag :value="s.participation_type" :severity="s.participation_type === 'WATCHER' ? 'secondary' : 'success'" class="!text-xs !px-1.5 !py-0.5" />
                </div>
                <div v-if="stepApprover(s.step_order)" class="ml-7 mt-1 text-xs text-gray-500 dark:text-gray-400">
                  <i class="pi pi-check-circle text-emerald-500 mr-1"></i>
                  {{ stepApprover(s.step_order).actor_name || '-' }}
                  <span v-if="stepApprover(s.step_order).actor_employee_code">({{ stepApprover(s.step_order).actor_employee_code }})</span>
                  <span v-if="stepApprover(s.step_order).actor_organization_name"> — {{ stepApprover(s.step_order).actor_organization_name }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="activeInstance.actions?.length">
            <p class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-2">{{ t('approval.history') }}</p>
            <div class="space-y-1">
              <div v-for="a in activeInstance.actions" :key="a.id" class="text-xs text-gray-600 dark:text-gray-300 flex items-center gap-2">
                <Tag :value="a.action" :severity="a.action === 'APPROVE' ? 'success' : 'danger'" class="!text-xs !px-1.5 !py-0.5" />
                <span>{{ formatDate(a.created_at) }}</span>
                <span v-if="a.note" class="text-gray-400">— {{ a.note }}</span>
              </div>
            </div>
          </div>

          <div v-if="activeTaskIsWatcher" class="space-y-2">
            <p class="text-xs text-gray-400 flex items-center gap-1.5">
              <i class="pi pi-eye"></i> {{ t('approval.watcher_note') }}
            </p>
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="taskDetailVisible = false" />
            </div>
          </div>
          <div v-else class="space-y-2">
            <Textarea v-model="actionNote" :placeholder="t('approval.note_placeholder')" rows="2" class="w-full" :class="{ 'p-invalid': noteError }" @input="noteError = ''" />
            <small v-if="noteError" class="text-red-500 text-xs block">{{ noteError }}</small>
            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="taskDetailVisible = false" />
              <Button :label="t('approval.reject')" severity="danger" outlined size="small" :loading="actionSubmitting" @click="submitAction('REJECT')" />
              <Button :label="t('approval.approve')" severity="success" size="small" :loading="actionSubmitting" @click="submitAction('APPROVE')" />
            </div>
          </div>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { formatDate as formatDateGlobal } from '@/utils/formatDate'
import api from '@/services/api'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import Select from 'primevue/select'
import SkeletonTable from '@/components/SkeletonTable.vue'
import { useAuth } from '@/stores/auth'

const { t, locale } = useI18n()
const toast = useToast()
const router = useRouter()
const { hasPermission } = useAuth()

// formatDate — date + time, built on the app-wide date formatter
// (utils/formatDate.js) so the date portion matches every other page
// (e.g. "30 July 2026"), with a time suffix appended since the global
// utility is date-only.
function formatDate(v) {
  if (!v) return '-'
  const datePart = formatDateGlobal(v, locale.value)
  if (!datePart) return '-'
  const time = new Date(v).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  return `${datePart} ${time}`
}

// formatDateOnly — date only, via the same app-wide formatter (handles the
// "YYYY-MM-DD" vs ISO-timestamp parsing itself, so a bare date string never
// shifts a day across a UTC-midnight boundary).
function formatDateOnly(v) {
  return formatDateGlobal(v, locale.value) || '-'
}

function moduleLabel(slug) {
  const label = t(`approval.module_names.${slug}`)
  return label !== `approval.module_names.${slug}` ? label : slug
}

const pendingTasks = ref([])
const pendingTotal = ref(0)
const tasksLoading = ref(false)

// ── Tab, filter & pagination ──
const activeTab = ref('pending') // 'pending' | 'done'
const statusFilter = ref(null)
const flowFilter = ref(null)
const flowOptions = ref([])
const flowsLoading = ref(false)
const currentPage = ref(1)
const perPage = ref(20)

function switchTab(tab) {
  if (activeTab.value === tab) return
  activeTab.value = tab
  statusFilter.value = null // filter status hanya untuk tab pending
  currentPage.value = 1
  loadTasks()
}

const firstRecord = computed(() => (currentPage.value - 1) * perPage.value)

// Filter status = status instance approval (task pending selalu PENDING,
// sehingga yang bermakna adalah status keseluruhan instance-nya).
const statusOptions = computed(() => [
  { label: t('approval.instance_status_pending'), value: 'PENDING' },
  { label: t('approval.instance_status_approved'), value: 'APPROVED' },
  { label: t('approval.instance_status_rejected'), value: 'REJECTED' },
  { label: t('approval.instance_status_cancelled'), value: 'CANCELLED' }
])

async function loadFlowOptions() {
  flowsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/approval/flows', { params: { page: 1, per_page: 100 } })
    const flows = res.data?.data || []
    flowOptions.value = flows.map(f => ({ label: f.name, value: f.id }))
  } catch {
    flowOptions.value = []
  } finally {
    flowsLoading.value = false
  }
}

function onPage(event) {
  currentPage.value = event.page + 1
  perPage.value = event.rows
  loadTasks()
}

function onFilterChange() {
  currentPage.value = 1
  loadTasks()
}

function clearFilters() {
  statusFilter.value = null
  flowFilter.value = null
  onFilterChange()
}

const taskSkeletonColumns = [
  { type: 'text', width: 'w-10', headerWidth: 'w-12' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'text', width: 'w-20', headerWidth: 'w-20' },
  { type: 'text', width: 'w-28', headerWidth: 'w-24' },
  { type: 'tag', width: 'w-16', headerWidth: 'w-12' },
  { type: 'text', width: 'w-32', headerWidth: 'w-24' },
  { type: 'icons', count: 1, headerWidth: 'w-16' }
]

async function loadTasks() {
  tasksLoading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: perPage.value,
      ...(activeTab.value === 'pending' && statusFilter.value ? { status: statusFilter.value } : {}),
      ...(flowFilter.value ? { flow_id: flowFilter.value } : {})
    }
    const endpoint = activeTab.value === 'done' ? '/api/v1/tenant/approval/tasks/done' : '/api/v1/tenant/approval/tasks/pending'
    const res = await api.get(endpoint, { params })
    const body = res.data
    pendingTasks.value = body?.data || []
    pendingTotal.value = body?.total || 0
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    tasksLoading.value = false
  }
}

const taskDetailVisible = ref(false)
const instanceLoading = ref(false)
const activeInstance = ref(null)
const activeTaskRef = ref(null)
const actionNote = ref('')
const noteError = ref('')
const actionSubmitting = ref(false)
const documentDetail = ref(null)
const documentLoading = ref(false)

// rowStatusLabel/rowStatusSeverity — for a WATCHER row, the task's own
// status is always PENDING (visible-but-not-actionable, see backend fix),
// which reads as ambiguous ("pending what? there's no approve/reject here").
// Watchers show the underlying instance's actual approval status instead —
// what the approver(s) have actually decided so far.
function rowStatusLabel(row) {
  if (row.participation_type === 'WATCHER' && row.instance_status) {
    return row.instance_status
  }
  return row.status
}
function rowStatusSeverity(row) {
  if (row.participation_type === 'WATCHER' && row.instance_status) {
    switch (row.instance_status) {
      case 'APPROVED': return 'success'
      case 'REJECTED': return 'danger'
      case 'CANCELLED': return 'secondary'
      default: return 'warn'
    }
  }
  return 'warn'
}

// stepApprover — the approve action recorded for a given step, if any, so
// each step in the preview can show who actually approved it (name,
// employee code, organization) once its status has moved past PENDING.
function stepApprover(stepOrder) {
  const actions = activeInstance.value?.actions || []
  return actions.find(a => a.step_order === stepOrder && a.action === 'APPROVE') || null
}

const activeTaskIsWatcher = computed(() => {
  if (!activeInstance.value || !activeTaskRef.value) return false
  const step = activeInstance.value.steps?.find(s => s.step_order === activeTaskRef.value.step_order)
  return step?.participation_type === 'WATCHER'
})

const isKPIModule = computed(() => {
  return ['performance_kpi_target', 'performance_kpi_realization'].includes(activeInstance.value?.module)
})

const isOKRModule = computed(() => {
  return ['okr_key_result', 'okr_assessment'].includes(activeInstance.value?.module)
})

const isLeaveModule = computed(() => activeInstance.value?.module === 'leave')
const leaveEmployeeName = ref('')
const leaveEmployeeCode = ref('')
const leaveTypeName = ref('')

// ── Employee Movement (mutasi) ──
const isEmployeemovementModule = computed(() => activeInstance.value?.module === 'employeemovement')

// movementTypeLabel/Severity — bilingual via employee_movement.type_* keys
// (sama seperti halaman Movements), dengan fallback raw slug.
function movementTypeLabel(type) {
  if (!type) return '-'
  const key = `employee_movement.type_${type}`
  return t(key) !== key ? t(key) : type
}

function movementTypeSeverity(type) {
  switch (type) {
    case 'promotion': return 'success'
    case 'demotion': return 'danger'
    case 'mutation': return 'info'
    case 'contract_extension': return 'warning'
    case 'status_change': return 'info'
    case 'retirement': return 'secondary'
    case 'offboarding': return 'danger'
    default: return 'secondary'
  }
}

// movementStatusLabel/Severity — bilingual via employee_movement.status_* keys.
function movementStatusLabel(status) {
  if (!status) return '-'
  const key = `employee_movement.status_${status}`
  return t(key) !== key ? t(key) : status
}

function movementStatusSeverity(status) {
  switch (status) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'rejected': return 'danger'
    case 'executed': return 'success'
    case 'cancelled': return 'secondary'
    case 'cancellation_pending': return 'warning'
    default: return 'secondary'
  }
}

const hasMovementFrom = computed(() => !!(documentDetail.value && (documentDetail.value.from_organization_name || documentDetail.value.from_position_name || documentDetail.value.from_employment_status_name)))
const hasMovementTo = computed(() => !!(documentDetail.value && (documentDetail.value.to_organization_name || documentDetail.value.to_position_name || documentDetail.value.to_employment_status_name)))

// ── Recruitment (lowongan & penawaran kandidat) ──
const isRecruitmentModule = computed(() => activeInstance.value?.module === 'recruitment')
const isOfferModule = computed(() => activeInstance.value?.module === 'recruitment_offer')
const isTrainingRequestModule = computed(() => activeInstance.value?.module === 'training_request')

// requisition_number/offer_number tidak dibawa response detail, jadi tidak
// perlu enrichment. Organization name di-resolve dari organization_id
// (response requisition hanya membawa id, sama seperti leave).
const requisitionOrganizationName = ref('')
async function loadRequisitionOrganization(organizationId) {
  requisitionOrganizationName.value = ''
  if (!organizationId) return
  try {
    const res = await api.get(`/api/v1/tenant/organizations/${organizationId}`)
    const org = res.data?.data
    requisitionOrganizationName.value = org?.nomenclature || org?.full_code || org?.code || ''
  } catch {}
}

// Training request detail hanya membawa employee_id/course_id — nama
// di-resolve client-side (pola loadLeaveNames).
const trainingRequestEmployeeName = ref('')
const trainingRequestCourseName = ref('')
async function loadTrainingRequestNames(employeeId, courseId) {
  trainingRequestEmployeeName.value = ''
  trainingRequestCourseName.value = ''
  const requests = []
  if (employeeId) {
    requests.push(
      api.get(`/api/v1/tenant/employees/${employeeId}`)
        .then(res => { trainingRequestEmployeeName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  if (courseId) {
    requests.push(
      api.get(`/api/v1/tenant/trainings/courses/${courseId}`)
        .then(res => { trainingRequestCourseName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  await Promise.all(requests)
}

function requisitionStatusLabel(status) {
  if (!status) return '-'
  const key = `requisitions.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function requisitionStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'submitted': return 'info'
    case 'open': return 'success'
    case 'in_progress': return 'info'
    case 'filled': return 'success'
    case 'rejected': return 'danger'
    case 'cancelled': return 'secondary'
    default: return 'secondary'
  }
}

function requisitionPriorityLabel(priority) {
  if (!priority) return '-'
  const key = `requisitions.priority_${String(priority).toLowerCase()}`
  return t(key) !== key ? t(key) : priority
}

function requisitionPrioritySeverity(priority) {
  switch (String(priority).toLowerCase()) {
    case 'urgent': return 'danger'
    case 'high': return 'warn'
    case 'medium': return 'info'
    case 'low': return 'secondary'
    default: return 'secondary'
  }
}

function requisitionReasonLabel(reasonType) {
  if (!reasonType) return '-'
  const key = `requisitions.reason_${String(reasonType).toLowerCase()}`
  return t(key) !== key ? t(key) : reasonType
}

function offerStatusLabel(status) {
  if (!status) return '-'
  const key = `offers.status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function offerStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'pending_approval': return 'info'
    case 'approved': return 'warning'
    case 'sent': return 'info'
    case 'accepted': return 'success'
    case 'rejected': return 'danger'
    case 'expired': return 'secondary'
    case 'withdrawn': return 'secondary'
    default: return 'secondary'
  }
}

function trainingRequestStatusLabel(status) {
  if (!status) return '-'
  const key = `training.request_status_${String(status).toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

function trainingRequestStatusSeverity(status) {
  switch (String(status).toLowerCase()) {
    case 'draft': return 'secondary'
    case 'submitted': return 'info'
    case 'pending_approval': return 'info'
    case 'approved': return 'success'
    case 'rejected': return 'danger'
    case 'cancelled': return 'secondary'
    default: return 'secondary'
  }
}

// formatCurrency — mata uang lokal tanpa desimal (pola halaman recruitment).
function formatCurrency(v) {
  if (v === null || v === undefined || v === '') return '-'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(Number(v))
}

// ── Attendance (overtime) ──
const isAttendanceModule = computed(() => activeInstance.value?.module === 'attendance')
const overtimeEmployeeName = ref('')
const overtimeEmployeeCode = ref('')

// Overtime's detail response only carries employee_id (no name) — resolved
// here client-side, same reasoning as loadLeaveNames above.
async function loadOvertimeNames(employeeId) {
  overtimeEmployeeName.value = ''
  overtimeEmployeeCode.value = ''
  if (!employeeId) return
  try {
    const res = await api.get(`/api/v1/tenant/employees/${employeeId}`)
    const emp = res.data?.data
    overtimeEmployeeName.value = emp?.name || ''
    overtimeEmployeeCode.value = emp?.employee_id || ''
  } catch {}
}

// ASSIGNED flow: the person who submitted the approval task IS the assigner
// (manager), so "Assigned By" maps straight to the task's submitter.
const isAssignedFlow = computed(() => documentDetail.value?.flow_type === 'ASSIGNED')
const overtimeAssignedByName = computed(() => {
  if (!isAssignedFlow.value) return null
  return activeTaskRef.value?.submitter_name || '-'
})

// formatTime — HH:mm lokal dari timestamp RFC3339 (pola sama AttendanceOvertime).
function formatTime(v) {
  if (!v) return '-'
  const d = new Date(v)
  if (isNaN(d.getTime())) return '-'
  return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function overtimeStatusSeverity(status) {
  switch (status) {
    case 'APPROVED': return 'success'
    case 'REJECTED': return 'danger'
    case 'PENDING_APPROVAL': return 'info'
    case 'WAITING_ACTUAL': return 'warning'
    case 'ACTUAL_SUBMITTED': return 'info'
    case 'CANCELLED': return 'secondary'
    default: return 'secondary'
  }
}

function overtimeStatusLabel(status) {
  if (!status) return '-'
  const key = `attendance.status_${status.toLowerCase()}`
  return t(key) !== key ? t(key) : status
}

// Lintas hari: end di hari berikutnya tetap disimpan sebagai timestamp tanggal
// +1 hari (pola FE), jadi bandingkan jam-jam (HH:MM) — sama seperti logika
// isCrossDayOf di AttendanceOvertime.vue.
const isOvertimeCrossDay = computed(() => {
  const s = documentDetail.value?.start_time_local
  const e = documentDetail.value?.end_time_local
  if (!s || !e) return false
  const sd = new Date(s)
  const ed = new Date(e)
  if (isNaN(sd.getTime()) || isNaN(ed.getTime())) return false
  return ed.getHours() * 60 + ed.getMinutes() <= sd.getHours() * 60 + sd.getMinutes()
})

const hasActualData = computed(() => {
  const d = documentDetail.value
  return !!(d && (d.actual_start_time_local || d.actual_end_time_local || d.actual_note || d.attachment_url))
})

// Leave's own GET /requests/:id response only has employee_id/leave_type_id
// (no names) — resolved here client-side, same reasoning as why the generic
// documentFields fallback denylists raw id fields (not meaningful to a
// reviewer on their own).
async function loadLeaveNames(employeeId, leaveTypeId) {
  leaveEmployeeName.value = ''
  leaveEmployeeCode.value = ''
  leaveTypeName.value = ''
  const requests = []
  if (employeeId) {
    requests.push(
      api.get(`/api/v1/tenant/employees/${employeeId}`)
        .then(res => {
          const emp = res.data?.data
          leaveEmployeeName.value = emp?.name || ''
          leaveEmployeeCode.value = emp?.employee_id || ''
        })
        .catch(() => {})
    )
  }
  if (leaveTypeId) {
    requests.push(
      api.get(`/api/v1/tenant/leave/types/${leaveTypeId}`)
        .then(res => { leaveTypeName.value = res.data?.data?.name || '' })
        .catch(() => {})
    )
  }
  await Promise.all(requests)
}

const okrObjectiveGroups = computed(() => {
  const details = documentDetail.value?.details || []
  const groups = {}
  const order = []
  for (const d of details) {
    const key = d.objective_id || d.objective_title
    if (!groups[key]) {
      groups[key] = { key, title: d.objective_title, items: [] }
      order.push(key)
    }
    groups[key].items.push(d)
  }
  return order.map(key => groups[key])
})

// Fields hidden from the generic "submitted data" fallback view — internal
// IDs/timestamps/relations that aren't meaningful to a reviewer, or that
// are already shown elsewhere in the dialog.
const DOCUMENT_FIELD_DENYLIST = new Set([
  'id', 'created_at', 'updated_at', 'deleted_at',
  'employee_id', 'organization_id', 'period_id', 'template_id',
  'approval_instance_id', 'target_approval_instance_id', 'realization_approval_instance_id',
  'kr_approval_instance_id', 'assessment_approval_instance_id',
  'details', 'program_items', 'items', 'documents',
  // Business Travel (raw IDs already implied by the popup's own context —
  // approval module → module/document_id — not useful to show as fields).
  'requester_id', 'business_travel_id', 'participant_id'
])

function humanizeLabel(key) {
  return key.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase())
}

function formatFieldValue(value) {
  if (value === null || value === undefined || value === '') return '-'
  if (typeof value === 'object') return JSON.stringify(value)
  return String(value)
}

// Generic key-value fallback for modules without a dedicated renderer above
// (reimbursement, employeemovement, attendance, payroll, ...) — the
// approval module is document-agnostic, so this covers whatever module
// ends up here without needing a bespoke view per module.
const documentFields = computed(() => {
  if (!documentDetail.value || typeof documentDetail.value !== 'object') return []
  return Object.entries(documentDetail.value)
    .filter(([key, value]) => !DOCUMENT_FIELD_DENYLIST.has(key) && typeof value !== 'object')
    .map(([key, value]) => ({ label: humanizeLabel(key), value: formatFieldValue(value) }))
})

// Maps an approval "module" slug to the endpoint that returns the
// underlying document being approved — the approval module itself only
// knows module + document_id, not the document's shape.
function documentEndpointFor(module, documentId) {
  switch (module) {
    case 'leave':
      return `/api/v1/tenant/leave/requests/${documentId}`
    case 'reimbursement':
      return `/api/v1/tenant/reimbursements/requests/${documentId}`
    case 'employeemovement':
      return `/api/v1/tenant/employee-movements/movements/${documentId}`
    case 'attendance':
      return `/api/v1/tenant/attendance/overtime-requests/${documentId}`
    case 'payroll':
      return `/api/v1/tenant/payroll/runs/${documentId}`
    case 'performance_kpi_target':
    case 'performance_kpi_realization':
      return `/api/v1/tenant/performance/kpi/evaluations/${documentId}/full`
    case 'okr_key_result':
    case 'okr_assessment':
      return `/api/v1/tenant/performance/okr/evaluations/${documentId}/details`
    case 'recruitment':
      return `/api/v1/tenant/recruitment/requisitions/${documentId}`
    case 'recruitment_offer':
      return `/api/v1/tenant/recruitment/offers/${documentId}`
    case 'training_request':
      return `/api/v1/tenant/trainings/requests/${documentId}`
    case 'business_travel':
      return `/api/v1/tenant/attendance/business-travels/${documentId}`
    case 'business_travel_settlement':
      // Settlements only have a nested detail route (business-travels/:id/
      // settlements/:settlementId) — this flat lookup exists specifically so
      // callers here, which only know the document_id, can fetch it.
      return `/api/v1/tenant/attendance/business-travel-settlements/${documentId}`
    default:
      return null
  }
}

async function loadDocumentDetail(module, documentId) {
  documentDetail.value = null
  const endpoint = documentEndpointFor(module, documentId)
  if (!endpoint) return
  documentLoading.value = true
  try {
    const res = await api.get(endpoint)
    documentDetail.value = res.data?.data || null
    if (module === 'leave' && documentDetail.value) {
      await loadLeaveNames(documentDetail.value.employee_id, documentDetail.value.leave_type_id)
    }
    if (module === 'attendance' && documentDetail.value) {
      await loadOvertimeNames(documentDetail.value.employee_id)
    }
    if (module === 'recruitment' && documentDetail.value) {
      await loadRequisitionOrganization(documentDetail.value.organization_id)
    }
    if (module === 'training_request' && documentDetail.value) {
      await loadTrainingRequestNames(documentDetail.value.employee_id, documentDetail.value.course_id)
    }
  } catch {
    documentDetail.value = null
  } finally {
    documentLoading.value = false
  }
}

async function openTaskDetail(task) {
  activeTaskRef.value = task
  actionNote.value = ''
  noteError.value = ''
  taskDetailVisible.value = true
  instanceLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/approval/instances/${task.instance_id}`)
    activeInstance.value = res.data?.data || null
    if (activeInstance.value) {
      await loadDocumentDetail(activeInstance.value.module, activeInstance.value.document_id)
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.failed_to_load'), life: 4000 })
  } finally {
    instanceLoading.value = false
  }
}

async function submitAction(action) {
  if (!activeTaskRef.value) return
  if (action === 'REJECT' && !actionNote.value?.trim()) {
    noteError.value = t('approval.reject_note_required')
    return
  }
  noteError.value = ''
  actionSubmitting.value = true
  try {
    await api.post(`/api/v1/tenant/approval/instances/${activeTaskRef.value.instance_id}/actions`, {
      action,
      note: actionNote.value || null
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('approval.action_submitted'), life: 3000 })
    taskDetailVisible.value = false
    await loadTasks()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: e.response?.data?.error?.message || t('message.operation_failed'), life: 4000 })
  } finally {
    actionSubmitting.value = false
  }
}

onMounted(() => {
  loadFlowOptions()
  loadTasks()
})
</script>
