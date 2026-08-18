<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-24 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="application">
      <!-- Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="text-base font-semibold text-navy-800 dark:text-gray-100 truncate">{{ candidateName }}</h2>
              <Tag :value="t('applications.status_' + (application.status || 'new').toLowerCase())" :severity="statusSeverity(application.status)" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ requisitionTitle }}</p>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <SelectLabel
              v-model="statusChange"
              :options="statusOptions"
              optionLabel="label"
              optionValue="value"
              :placeholder="t('applications.change_status')"
              class="!w-44"
              @update:modelValue="onStatusChange"
            />
            <Button :label="t('applications.back_to_applications')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/recruitment/applications')" />
          </div>
        </div>
      </div>

      <!-- Tabs -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
        <div class="flex items-center gap-1 px-3 pt-2 border-b border-gray-200 dark:border-gray-700 overflow-x-auto">
          <button
            v-for="tab in tabs"
            :key="tab.key"
            type="button"
            class="px-3 py-2 text-sm font-medium rounded-t-md transition-colors whitespace-nowrap"
            :class="activeTab === tab.key ? 'text-amber-600 dark:text-amber-400 border-b-2 border-amber-500' : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200'"
            @click="activeTab = tab.key"
          >
            {{ t(tab.labelKey) }}
          </button>
        </div>

        <!-- History -->
        <div v-if="activeTab === 'history'" class="p-4">
          <div v-if="history.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in history" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <i class="pi pi-arrow-right-arrow-left text-gray-400 text-xs"></i>
              <div class="min-w-0 flex-1">
                <p class="text-sm text-navy-800 dark:text-gray-100">
                  <span v-if="item.from_stage" class="text-gray-400">{{ item.from_stage.name }} →</span>
                  <span class="font-medium">{{ item.to_stage.name }}</span>
                </p>
                <p class="text-xs text-gray-400">{{ formatTimestamp(item.changed_at) }}</p>
              </div>
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('applications.history_empty') }}</div>
        </div>

        <!-- Screening -->
        <div v-if="activeTab === 'screening'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('applications.add_screening')" icon="pi pi-plus" size="small" @click="openScreeningDialog()" />
          </div>
          <div v-if="screenings.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in screenings" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <Tag :value="t('applications.result_' + (item.result || 'hold').toLowerCase())" :severity="resultSeverity(item.result)" class="!text-xs !px-1.5 !py-0.5" />
              <div class="min-w-0 flex-1">
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ item.notes || '—' }}</p>
                <p v-if="item.score" class="text-xs text-gray-400">{{ t('applications.score') }}: {{ item.score }}</p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-7 !h-7" @click="confirmDeleteScreening(item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('applications.screening_empty') }}</div>
        </div>

        <!-- Assessment -->
        <div v-if="activeTab === 'assessment'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('applications.add_to_assessment')" icon="pi pi-plus" size="small" :disabled="!availableAssessments.length" @click="openAssessmentDialog()" />
          </div>
          <div v-if="assessmentLoading" class="space-y-2">
            <div v-for="i in 2" :key="i" class="h-10 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
          </div>
          <div v-else-if="participations.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in participations" :key="item.participant.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ item.assessmentName }}</p>
                <p class="text-xs text-gray-400">{{ t('applications.participant_status_' + (item.participant.status || 'invited').toLowerCase()) }} <span v-if="item.participant.score">· {{ t('applications.score') }}: {{ item.participant.score }}</span></p>
                <p v-if="item.participant.recommendation" class="text-xs text-gray-400 mt-0.5">{{ t('applications.assessment_recommendation') }}: {{ item.participant.recommendation }}</p>
              </div>
              <Tag v-if="item.participant.result" :value="t('applications.result_' + item.participant.result.toLowerCase())" :severity="resultSeverity(item.participant.result)" class="!text-xs !px-1.5 !py-0.5" />
              <Button :label="t('applications.assess_participant')" icon="pi pi-pencil" size="small" text class="!h-7" @click="openAssessmentResult(item.participant)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('applications.assessment_empty') }}</div>
        </div>

        <!-- Interviews -->
        <div v-if="activeTab === 'interviews'" class="p-4">
          <div class="flex items-center justify-end mb-3">
            <Button :label="t('applications.add_interview')" icon="pi pi-plus" size="small" @click="openInterviewDialog()" />
          </div>
          <div v-if="interviews.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in interviews" :key="item.id" class="flex items-center gap-3 px-3 py-2.5">
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-navy-800 dark:text-gray-100">{{ item.stage }}</p>
                <p class="text-xs text-gray-400">{{ t('applications.interviewer') }}: {{ employeeName(item.interviewer_id) }} <span v-if="item.score">· {{ t('applications.score') }}: {{ item.score }}</span></p>
              </div>
              <Tag :value="t('applications.interview_status_' + (item.status || 'scheduled').toLowerCase())" :severity="interviewStatusSeverity(item.status)" class="!text-xs !px-1.5 !py-0.5" />
              <Button :label="t('applications.manage')" size="small" severity="secondary" outlined class="!text-xs" @click="openManageInterview(item)" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-gray-400">{{ t('applications.interviews_empty') }}</div>
        </div>

        <!-- Match Score -->
        <div v-if="activeTab === 'match_score'" class="p-4">
          <div v-if="matchScoreLoading" class="h-24 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
          <template v-else-if="matchScore">
            <div v-if="matchScore.score !== null && matchScore.score !== undefined" class="mb-4">
              <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">{{ t('applications.overall_match') }}</p>
              <p class="text-3xl font-bold text-emerald-600 dark:text-emerald-400">{{ Math.round(matchScore.score) }}%</p>
              <p v-if="matchScore.note" class="text-xs text-amber-600 dark:text-amber-400 mt-1"><i class="pi pi-info-circle mr-1"></i>{{ matchScoreNote }}</p>
            </div>
            <p v-else class="text-sm text-gray-400 mb-4">{{ matchScoreNote || t('applications.match_score_empty') }}</p>
            <div v-if="matchScore.breakdown && matchScore.breakdown.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
              <div v-for="(b, i) in matchScore.breakdown" :key="i" class="flex items-center gap-3 px-3 py-2.5">
                <div class="min-w-0 flex-1">
                  <p class="text-sm text-navy-800 dark:text-gray-100">{{ b.competency_name || b.competency_id }}</p>
                  <p class="text-xs text-gray-400">{{ t('applications.required_level') }} {{ b.required_level }} · {{ t('applications.candidate_level') }} {{ b.candidate_level }} · {{ t('applications.weight') }} {{ b.weight }}</p>
                </div>
              </div>
            </div>
          </template>
        </div>

        <!-- Penilaian Kandidat -->
        <div v-if="activeTab === 'evaluation'" class="p-4">
          <div v-if="evaluationLoading" class="h-24 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
          <template v-else-if="evaluation">
            <!-- Skor -->
            <div class="mb-4 flex items-center gap-4 flex-wrap">
              <div>
                <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">{{ t('applications.overall_match') }}</p>
                <p class="text-3xl font-bold text-amber-600 dark:text-amber-400">{{ evaluation.assessment?.score !== null && evaluation.assessment?.score !== undefined ? evaluation.assessment.score : '—' }}</p>
              </div>
              <p class="text-xs text-gray-400 max-w-md">{{ t('applications.evaluation_hint') }}</p>
            </div>

            <!-- Requirement + penilaian pendidikan/pengalaman -->
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3">
                <p class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-1"><i class="pi pi-graduation-cap mr-1 text-gray-400"></i>{{ t('requisitions.education') }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">{{ evaluation.requirements.education || '—' }}</p>
                <div class="flex items-center gap-2 mb-2">
                  <RadioLabel v-model="evaluationForm.education_match" :value="true" :id="'ev-edu-yes'" :label="t('applications.evaluation_match')" class="!flex-1" />
                  <RadioLabel v-model="evaluationForm.education_match" :value="false" :id="'ev-edu-no'" :label="t('applications.evaluation_not_match')" class="!flex-1" />
                </div>
                <FormRow :label="t('candidates.notes')"><Textarea v-model="evaluationForm.education_note" :rows="2" class="!w-full" /></FormRow>
              </div>
              <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3">
                <p class="text-sm font-semibold text-gray-700 dark:text-gray-200 mb-1"><i class="pi pi-briefcase mr-1 text-gray-400"></i>{{ t('requisitions.experience') }}</p>
                <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">{{ evaluation.requirements.experience || '—' }}</p>
                <div class="flex items-center gap-2 mb-2">
                  <RadioLabel v-model="evaluationForm.experience_match" :value="true" :id="'ev-exp-yes'" :label="t('applications.evaluation_match')" class="!flex-1" />
                  <RadioLabel v-model="evaluationForm.experience_match" :value="false" :id="'ev-exp-no'" :label="t('applications.evaluation_not_match')" class="!flex-1" />
                </div>
                <FormRow :label="t('candidates.notes')"><Textarea v-model="evaluationForm.experience_note" :rows="2" class="!w-full" /></FormRow>
              </div>
            </div>

            <!-- Kompetensi: requirement vs level kandidat -->
            <div class="border border-gray-200 dark:border-gray-700 rounded-lg mb-4">
              <div class="px-3 py-2 border-b border-gray-200 dark:border-gray-700">
                <p class="text-sm font-semibold text-gray-700 dark:text-gray-200"><i class="pi pi-star mr-1 text-gray-400"></i>{{ t('applications.evaluation_competencies') }}</p>
              </div>
              <div v-if="evaluation.requirements.competencies && evaluation.requirements.competencies.length" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="comp in evaluation.requirements.competencies" :key="comp.competency_id" class="flex items-center gap-3 px-3 py-2.5 flex-wrap">
                  <div class="min-w-0 flex-1">
                    <p class="text-sm text-navy-800 dark:text-gray-100">{{ comp.competency_name || comp.competency_id }}</p>
                    <p class="text-xs text-gray-400">{{ t('applications.required_level') }} {{ comp.required_level }}<span v-if="comp.weight"> · {{ t('applications.weight') }} {{ comp.weight }}</span></p>
                  </div>
                  <div class="flex items-center gap-2 w-full lg:w-auto">
                    <span class="text-xs text-gray-400 shrink-0">{{ t('applications.candidate_level') }}</span>
                    <SelectLabel
                      v-model="evaluationForm.competency_levels[comp.competency_id]"
                      :options="levelOptions"
                      optionLabel="label"
                      optionValue="value"
                      :placeholder="t('common.select')"
                      class="!w-full lg:!w-72"
                      showClear
                    />
                  </div>
                </div>
              </div>
              <p v-else class="px-3 py-6 text-center text-sm text-gray-400">{{ t('applications.evaluation_competencies_empty') }}</p>
            </div>

            <div class="flex items-center justify-end gap-2">
              <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="evaluationSaving" @click="saveEvaluation()" />
            </div>
          </template>
        </div>
      </div>
    </template>

    <!-- Add screening dialog -->
    <Dialog v-model:visible="screeningDialogVisible" :header="t('applications.add_screening')" :modal="true" class="!w-[min(95vw,460px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('applications.result')">
          <SelectLabel v-model="screeningForm.result" :options="resultOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('applications.score')"><InputNumber v-model="screeningForm.score" :min="0" :max="100" class="!w-full" /></FormRow>
        <FormRow :label="t('candidates.notes')"><Textarea v-model="screeningForm.notes" :rows="2" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="screeningDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveScreening()" />
        </div>
      </template>
    </Dialog>

    <!-- Add to assessment dialog -->
    <Dialog v-model:visible="assessmentDialogVisible" :header="t('applications.add_to_assessment')" :modal="true" class="!w-[min(95vw,460px)]">
      <FormRow :label="t('applications.assessment')">
        <SelectLabel v-model="assessmentForm.assessment_id" :options="availableAssessmentOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
      </FormRow>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="assessmentDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="joinAssessment()" />
        </div>
      </template>
    </Dialog>

    <!-- Fill assessment result dialog (penilaian kandidat) -->
    <Dialog v-model:visible="assessmentResultVisible" :header="t('applications.assessment_result_title')" :modal="true" class="!w-[min(95vw,520px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('applications.assessment_status')">
          <SelectLabel v-model="assessmentResultForm.status" :options="participantStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('applications.score')">
          <InputNumber v-model="assessmentResultForm.score" :min="0" :max="100" class="!w-full" />
        </FormRow>
        <FormRow :label="t('applications.assessment_result')">
          <SelectLabel v-model="assessmentResultForm.result" :options="assessmentResultOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('applications.assessment_recommendation')">
          <Textarea v-model="assessmentResultForm.recommendation" :rows="3" class="!w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" :disabled="itemSaving" @click="assessmentResultVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveAssessmentResult()" />
        </div>
      </template>
    </Dialog>

    <!-- Add interview dialog -->
    <Dialog v-model:visible="interviewDialogVisible" :header="t('applications.add_interview')" :modal="true" class="!w-[min(95vw,520px)]">
      <div class="grid grid-cols-1 gap-3">
        <FormRow :label="t('applications.interviewer')" :required="true">
          <SelectLabel v-model="interviewForm.interviewer_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
        </FormRow>
        <FormRow :label="t('applications.stage')"><TextInput v-model="interviewForm.stage" class="!w-full" /></FormRow>
        <FormRow :label="t('applications.scheduled_date')" :required="true"><DateInput v-model="interviewForm.scheduled_date" class="!w-full" /></FormRow>
        <FormRow :label="t('applications.duration_minutes')"><InputNumber v-model="interviewForm.duration_minutes" :min="15" class="!w-full" /></FormRow>
        <FormRow :label="t('applications.location')"><TextInput v-model="interviewForm.location" class="!w-full" /></FormRow>
        <FormRow :label="t('applications.meeting_link')"><TextInput v-model="interviewForm.meeting_link" class="!w-full" /></FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="interviewDialogVisible = false" />
          <Button :label="t('common.save')" icon="pi pi-check" size="small" :loading="itemSaving" @click="saveInterview()" />
        </div>
      </template>
    </Dialog>

    <!-- Manage interview dialog: interviewers + scorecard + complete -->
    <Dialog v-model:visible="manageDialogVisible" :header="t('applications.manage_interview')" :modal="true" class="!w-[min(95vw,960px)]">
      <div v-if="manageInterview" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm font-semibold text-navy-800 dark:text-gray-100">{{ manageInterview.stage }}</p>
            <Tag :value="t('applications.interview_status_' + (manageInterview.status || 'scheduled').toLowerCase())" :severity="interviewStatusSeverity(manageInterview.status)" class="!text-xs !px-1.5 !py-0.5 mt-1" />
          </div>
          <Button
            v-if="manageInterview.status !== 'COMPLETED'"
            :label="t('applications.complete_interview')"
            icon="pi pi-check-circle"
            size="small"
            severity="success"
            :loading="completingInterview"
            @click="completeInterview()"
          />
        </div>

        <!-- Interviewers -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('applications.interviewers') }}</h4>
            <Button :label="t('common.add')" icon="pi pi-plus" text size="small" class="!text-xs" @click="openAddInterviewer()" />
          </div>
          <div v-if="manageInterviewers.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="iv in manageInterviewers" :key="iv.id" class="flex items-center gap-2 px-3 py-2">
              <span class="text-sm text-gray-700 dark:text-gray-300 flex-1">{{ employeeName(iv.employee_id) }}</span>
              <span v-if="iv.role" class="text-xs text-gray-400">{{ iv.role }}</span>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-6 !h-6" @click="removeInterviewer(iv)" />
            </div>
          </div>
          <p v-else class="text-xs text-gray-400 py-2">{{ t('applications.interviewers_empty') }}</p>

          <div v-if="addInterviewerVisible" class="flex items-center gap-2 mt-2">
            <SelectLabel v-model="newInterviewerId" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" class="!w-full" showClear />
            <Button icon="pi pi-check" size="small" :loading="itemSaving" @click="saveInterviewer()" />
          </div>
        </div>

        <!-- Scorecard -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('applications.scorecard') }}</h4>
            <Button :label="t('common.add')" icon="pi pi-plus" text size="small" class="!text-xs" @click="openAddScorecardItem()" />
          </div>
          <div v-if="manageScorecard.length" class="divide-y divide-gray-100 dark:divide-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg">
            <div v-for="item in manageScorecard" :key="item.id" class="flex items-center gap-2 px-3 py-2">
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-700 dark:text-gray-300">{{ item.criterion }}</p>
                <p class="text-xs text-gray-400">{{ t('candidates.level') }}: {{ item.weight }}% <span v-if="item.score">· {{ t('applications.score') }}: {{ item.score }}</span></p>
              </div>
              <Button icon="pi pi-trash" text severity="danger" size="small" class="!w-6 !h-6" @click="removeScorecardItem(item)" />
            </div>
          </div>
          <p v-else class="text-xs text-gray-400 py-2">{{ t('applications.scorecard_empty') }}</p>

          <div v-if="addScorecardVisible" class="grid grid-cols-3 gap-2 mt-2">
            <TextInput v-model="newScorecard.criterion" :placeholder="t('applications.criterion')" class="!w-full !col-span-1" />
            <InputNumber v-model="newScorecard.weight" :placeholder="t('applications.weight')" :min="0" :max="100" class="!w-full" />
            <div class="flex items-center gap-1">
              <InputNumber v-model="newScorecard.score" :placeholder="t('applications.score')" :min="0" :max="100" class="!w-full" />
              <Button icon="pi pi-check" size="small" :loading="itemSaving" @click="saveScorecardItem()" />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button :label="t('common.close')" severity="secondary" outlined size="small" @click="manageDialogVisible = false" />
      </template>
    </Dialog>

    <ConfirmDeleteDialog
      v-model:visible="deleteDialogVisible"
      :title="t('common.delete')"
      :message="t('candidates.delete_item_confirm')"
      :loading="deleting"
      @confirm="doDeleteScreening()"
    />

    <!-- Konfirmasi ubah status aplikasi -->
    <Dialog v-model:visible="statusDialogVisible" :header="t('applications.confirm_status_title')" :modal="true" class="!w-[min(95vw,460px)]">
      <div class="space-y-3">
        <div class="flex items-center gap-2 text-sm flex-wrap">
          <Tag :value="t('applications.status_' + (application?.status || 'new').toLowerCase())" :severity="statusSeverity(application?.status)" class="!text-xs !px-1.5 !py-0.5" />
          <i class="pi pi-arrow-right text-xs text-gray-400"></i>
          <Tag v-if="pendingStatus" :value="t('applications.status_' + pendingStatus.toLowerCase())" :severity="statusSeverity(pendingStatus)" class="!text-xs !px-1.5 !py-0.5" />
        </div>
        <p class="text-sm text-gray-600 dark:text-gray-300">{{ t('applications.confirm_status_message') }}</p>
        <FormRow v-if="pendingStatus === 'REJECTED'" :label="t('applications.status_rejection_reason')">
          <Textarea v-model="statusRejectionReason" :rows="2" class="!w-full" />
        </FormRow>
        <FormRow :label="t('applications.status_notes')">
          <Textarea v-model="statusNotes" :rows="2" class="!w-full" />
        </FormRow>
      </div>
      <template #footer>
        <div class="flex items-center justify-end gap-2">
          <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="cancelStatusChange" />
          <Button :label="t('common.confirm')" icon="pi pi-check" size="small" :loading="statusSaving" @click="confirmStatusChange" />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from '@/composables/useI18n'
import { useToast } from 'primevue/usetoast'
import api from '@/services/api'
import { getErrorMessage } from '@/services/responseHandler'

import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import Tag from 'primevue/tag'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import TextInput from '@/components/TextInput.vue'
import FormRow from '@/components/FormRow.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'
import RadioLabel from '@/components/RadioLabel.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const applicationId = route.params.id

const loading = ref(true)
const application = ref(null)
const activeTab = ref('history')
const statusChange = ref(null)

// ── Konfirmasi ubah status ──
const statusDialogVisible = ref(false)
const pendingStatus = ref(null)
const statusNotes = ref('')
const statusRejectionReason = ref('')
const statusSaving = ref(false)

const candidates = ref([])
const requisitions = ref([])
const employees = ref([])

const history = ref([])
const screenings = ref([])
const interviews = ref([])
const assessments = ref([])
const assessmentLoading = ref(false)
const participations = ref([])
const matchScore = ref(null)
const matchScoreLoading = ref(false)

// ── Penilaian Kandidat (G-12) ──
const evaluation = ref(null)
const evaluationLoading = ref(false)
const evaluationSaving = ref(false)
const evaluationForm = ref({ education_match: null, education_note: '', experience_match: null, experience_note: '', competency_levels: {} })
const jobValues = ref([])
// Options level dari job_management_values tipe technical/managerial (skala
// kompetensi, pola Job Management: "Lv.X — deskripsi"). Teknis didahulukan
// (skala 1-8), fallback ke managerial; fallback terakhir "Lv.X" bila data
// belum ter-load.
const levelOptions = computed(() => {
  const byLevel = {}
  ;(jobValues.value || []).forEach(v => {
    if (!v.level || (v.type !== 'technical' && v.type !== 'managerial')) return
    const existing = byLevel[v.level]
    if (!existing || v.type === 'technical') byLevel[v.level] = v.descriptions || ''
  })
  const maxLevel = Object.keys(byLevel).length ? Math.max(...Object.keys(byLevel).map(Number)) : 5
  const opts = []
  for (let lvl = 1; lvl <= maxLevel; lvl++) {
    opts.push({
      label: byLevel[lvl] ? `Lv.${lvl} — ${byLevel[lvl]}` : `Lv.${lvl}`,
      name: byLevel[lvl] || '',
      value: lvl
    })
  }
  return opts
})

const tabs = [
  { key: 'history', labelKey: 'applications.tab_history' },
  { key: 'screening', labelKey: 'applications.tab_screening' },
  { key: 'assessment', labelKey: 'applications.tab_assessment' },
  { key: 'interviews', labelKey: 'applications.tab_interviews' },
  { key: 'evaluation', labelKey: 'applications.tab_evaluation' },
  { key: 'match_score', labelKey: 'applications.tab_match_score' }
]

const statusOptions = computed(() => ['NEW', 'SCREENED', 'SHORTLISTED', 'INTERVIEWED', 'OFFERED', 'ACCEPTED', 'REJECTED', 'WITHDRAWN'].map(v => ({ label: t(`applications.status_${v.toLowerCase()}`), value: v })))
const resultOptions = computed(() => ['PASS', 'FAIL', 'HOLD'].map(v => ({ label: t(`applications.result_${v.toLowerCase()}`), value: v })))
const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))

// Terjemahkan note backend match-score (beberapa note English-only dari
// service) ke locale aktif; fallback ke teks mentah bila tidak dikenali.
const matchScoreNote = computed(() => {
  const raw = matchScore.value?.note || ''
  const map = {
    'competency levels sourced from candidate evaluation (G-12)': t('applications.match_score_note_evaluation'),
    'competencies sourced from Job Management (organization default); requisition has no override of its own': t('applications.match_score_note_job_management'),
    'no competencies defined for this requisition (neither its own override nor a Job Management default); nothing to match': t('applications.match_score_note_no_competencies')
  }
  return map[raw] || raw
})

const candidateName = computed(() => {
  const c = candidates.value.find(x => x.id === application.value?.candidate_id)
  return c ? `${c.first_name} ${c.last_name}` : (application.value?.candidate_id || '')
})
const requisitionTitle = computed(() => {
  const r = requisitions.value.find(x => x.id === application.value?.requisition_id)
  return r ? r.title : (application.value?.requisition_id || '')
})

function employeeName(id) {
  const e = employees.value.find(x => x.id === id)
  return e ? e.name : (id || '—')
}

function statusSeverity(status) {
  switch (status) {
    case 'NEW': return 'secondary'
    case 'SCREENED': case 'SHORTLISTED': return 'info'
    case 'INTERVIEWED': return 'help'
    case 'OFFERED': return 'warn'
    case 'ACCEPTED': return 'success'
    case 'REJECTED': case 'WITHDRAWN': return 'danger'
    default: return 'secondary'
  }
}
function resultSeverity(result) {
  switch (result) {
    case 'PASS': return 'success'
    case 'FAIL': return 'danger'
    default: return 'secondary'
  }
}
function interviewStatusSeverity(status) {
  switch (status) {
    case 'COMPLETED': return 'success'
    case 'CANCELLED': return 'danger'
    case 'RESCHEDULED': return 'warn'
    default: return 'info'
  }
}

function formatTimestamp(value) {
  if (!value) return '—'
  const ms = Number(value) / 1000000
  if (!Number.isFinite(ms) || ms <= 0) return '—'
  return new Date(ms).toLocaleString()
}

function cleanPayload(payload) {
  const out = { ...payload }
  Object.keys(out).forEach(k => {
    if (out[k] === '' || out[k] === null || out[k] === undefined) delete out[k]
  })
  return out
}

async function onStatusChange(newStatus) {
  if (!newStatus || newStatus === application.value.status) return
  // Tampilkan dialog konfirmasi dulu sebelum menyimpan.
  pendingStatus.value = newStatus
  statusNotes.value = ''
  statusRejectionReason.value = ''
  statusDialogVisible.value = true
}

function cancelStatusChange() {
  statusDialogVisible.value = false
  pendingStatus.value = null
  statusChange.value = null
}

async function confirmStatusChange() {
  if (!pendingStatus.value) return
  statusSaving.value = true
  try {
    const payload = { status: pendingStatus.value }
    if (pendingStatus.value === 'REJECTED' && statusRejectionReason.value) {
      payload.rejection_reason = statusRejectionReason.value
    }
    if (statusNotes.value) {
      payload.notes = statusNotes.value
    }
    await api.put(`/api/v1/tenant/recruitment/applications/${applicationId}/status`, payload)
    statusDialogVisible.value = false
    pendingStatus.value = null
    statusChange.value = null
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('applications.status_updated'), life: 3000 })
    loadApplication()
    loadHistory()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
    statusChange.value = null
  } finally {
    statusSaving.value = false
  }
}

// ── Screening ──
const screeningDialogVisible = ref(false)
const screeningForm = ref({})
const itemSaving = ref(false)
const deleteDialogVisible = ref(false)
const deleting = ref(false)
const pendingDeleteScreening = ref(null)

function openScreeningDialog() {
  screeningForm.value = { result: 'HOLD', score: null, notes: '' }
  screeningDialogVisible.value = true
}
async function saveScreening() {
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/applications/${applicationId}/screenings`, cleanPayload(screeningForm.value))
    screeningDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadScreenings()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
function confirmDeleteScreening(item) {
  pendingDeleteScreening.value = item
  deleteDialogVisible.value = true
}
async function doDeleteScreening() {
  if (!pendingDeleteScreening.value) return
  deleting.value = true
  try {
    await api.delete(`/api/v1/tenant/recruitment/screenings/${pendingDeleteScreening.value.id}`)
    deleteDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_deleted'), life: 3000 })
    loadScreenings()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    deleting.value = false
  }
}

// ── Assessment ──
const assessmentDialogVisible = ref(false)
const assessmentForm = ref({})
const assessmentResultVisible = ref(false)
const assessmentResultForm = ref({})
const assessmentResultTarget = ref(null)

const participantStatusOptions = [
  { label: t('applications.participant_status_invited'), value: 'INVITED' },
  { label: t('applications.participant_status_completed'), value: 'COMPLETED' },
  { label: t('applications.participant_status_no_show'), value: 'NO_SHOW' }
]
const assessmentResultOptions = [
  { label: t('applications.result_pass'), value: 'PASS' },
  { label: t('applications.result_fail'), value: 'FAIL' },
  { label: t('applications.result_hold'), value: 'HOLD' }
]

function openAssessmentResult(participant) {
  assessmentResultTarget.value = participant
  assessmentResultForm.value = {
    status: participant.status || 'INVITED',
    score: participant.score ?? null,
    result: participant.result || null,
    recommendation: participant.recommendation || ''
  }
  assessmentResultVisible.value = true
}
async function saveAssessmentResult() {
  const target = assessmentResultTarget.value
  if (!target) return
  itemSaving.value = true
  try {
    await api.put(`/api/v1/tenant/recruitment/assessment-participants/${target.id}`, cleanPayload(assessmentResultForm.value))
    assessmentResultVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadAssessmentParticipation()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
const availableAssessments = computed(() => {
  const joinedIds = new Set(participations.value.map(p => p.assessmentId))
  return assessments.value.filter(a => !joinedIds.has(a.id))
})
const availableAssessmentOptions = computed(() => availableAssessments.value.map(a => ({ label: a.name, value: a.id })))

function openAssessmentDialog() {
  assessmentForm.value = { assessment_id: null }
  assessmentDialogVisible.value = true
}
async function joinAssessment() {
  if (!assessmentForm.value.assessment_id) return
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/assessments/${assessmentForm.value.assessment_id}/participants`, { application_id: applicationId })
    assessmentDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadAssessmentParticipation()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

async function loadAssessmentParticipation() {
  assessmentLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/recruitment/assessments')
    assessments.value = res.data?.data || []
    const results = await Promise.allSettled(
      assessments.value.map(a => api.get(`/api/v1/tenant/recruitment/assessments/${a.id}/participants`))
    )
    const found = []
    results.forEach((r, idx) => {
      if (r.status !== 'fulfilled') return
      const list = r.value.data?.data || []
      const mine = list.find(p => p.application_id === applicationId)
      if (mine) found.push({ assessmentId: assessments.value[idx].id, assessmentName: assessments.value[idx].name, participant: mine })
    })
    participations.value = found
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    assessmentLoading.value = false
  }
}

// ── Interviews ──
const interviewDialogVisible = ref(false)
const interviewForm = ref({})

function openInterviewDialog() {
  interviewForm.value = { interviewer_id: null, stage: 'FIRST_INTERVIEW', scheduled_date: '', duration_minutes: 60, location: '', meeting_link: '' }
  interviewDialogVisible.value = true
}
async function saveInterview() {
  if (!interviewForm.value.interviewer_id || !interviewForm.value.scheduled_date) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('message.failed_to_save'), life: 4000 })
    return
  }
  itemSaving.value = true
  try {
    const scheduledAt = new Date(interviewForm.value.scheduled_date + 'T00:00:00').getTime() * 1000000
    const payload = cleanPayload({
      application_id: applicationId,
      interviewer_id: interviewForm.value.interviewer_id,
      stage: interviewForm.value.stage,
      scheduled_at: scheduledAt,
      duration_minutes: interviewForm.value.duration_minutes,
      location: interviewForm.value.location,
      meeting_link: interviewForm.value.meeting_link
    })
    await api.post('/api/v1/tenant/recruitment/interviews', payload)
    interviewDialogVisible.value = false
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('candidates.item_added'), life: 3000 })
    loadInterviews()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}

// ── Manage interview: interviewers + scorecard + complete ──
const manageDialogVisible = ref(false)
const manageInterview = ref(null)
const manageInterviewers = ref([])
const manageScorecard = ref([])
const addInterviewerVisible = ref(false)
const newInterviewerId = ref(null)
const addScorecardVisible = ref(false)
const newScorecard = ref({ criterion: '', weight: null, score: null })
const completingInterview = ref(false)

async function openManageInterview(item) {
  manageInterview.value = item
  addInterviewerVisible.value = false
  addScorecardVisible.value = false
  manageDialogVisible.value = true
  await Promise.all([loadManageInterviewers(), loadManageScorecard()])
}

async function loadManageInterviewers() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/interviews/${manageInterview.value.id}/interviewers`)
    manageInterviewers.value = res.data?.data || []
  } catch {
    manageInterviewers.value = []
  }
}
async function loadManageScorecard() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/interviews/${manageInterview.value.id}/scorecard-items`)
    manageScorecard.value = res.data?.data || []
  } catch {
    manageScorecard.value = []
  }
}

function openAddInterviewer() {
  newInterviewerId.value = null
  addInterviewerVisible.value = true
}
async function saveInterviewer() {
  if (!newInterviewerId.value) return
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/interviews/${manageInterview.value.id}/interviewers`, { employee_id: newInterviewerId.value })
    addInterviewerVisible.value = false
    loadManageInterviewers()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
async function removeInterviewer(iv) {
  try {
    await api.delete(`/api/v1/tenant/recruitment/interviewers/${iv.id}`)
    loadManageInterviewers()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  }
}

function openAddScorecardItem() {
  newScorecard.value = { criterion: '', weight: null, score: null }
  addScorecardVisible.value = true
}
async function saveScorecardItem() {
  if (!newScorecard.value.criterion?.trim()) return
  itemSaving.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/interviews/${manageInterview.value.id}/scorecard-items`, cleanPayload(newScorecard.value))
    addScorecardVisible.value = false
    loadManageScorecard()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    itemSaving.value = false
  }
}
async function removeScorecardItem(item) {
  try {
    await api.delete(`/api/v1/tenant/recruitment/scorecard-items/${item.id}`)
    loadManageScorecard()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  }
}
async function completeInterview() {
  completingInterview.value = true
  try {
    await api.post(`/api/v1/tenant/recruitment/interviews/${manageInterview.value.id}/complete`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('applications.interview_completed'), life: 3000 })
    manageDialogVisible.value = false
    loadInterviews()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    completingInterview.value = false
  }
}

// ── Loaders ──
async function loadApplication() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/applications/${applicationId}`)
    application.value = res.data?.data || null
    statusChange.value = null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  }
}
async function loadHistory() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/applications/${applicationId}/history`)
    history.value = res.data?.data || []
  } catch {
    history.value = []
  }
}
async function loadScreenings() {
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/applications/${applicationId}/screenings`)
    screenings.value = res.data?.data || []
  } catch {
    screenings.value = []
  }
}
async function loadInterviews() {
  try {
    const res = await api.get('/api/v1/tenant/recruitment/interviews', { params: { application_id: applicationId, per_page: 100 } })
    interviews.value = res.data?.data || []
  } catch {
    interviews.value = []
  }
}
async function loadMatchScore() {
  matchScoreLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/applications/${applicationId}/match-score`)
    matchScore.value = res.data?.data || null
  } catch {
    matchScore.value = null
  } finally {
    matchScoreLoading.value = false
  }
}
async function loadEvaluation() {
  evaluationLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/recruitment/applications/${applicationId}/assessment`)
    evaluation.value = res.data?.data || null
    const req = evaluation.value?.requirements
    const saved = evaluation.value?.assessment
    const levels = {}
    ;(req?.competencies || []).forEach(c => { levels[c.competency_id] = null })
    ;(saved?.competency_levels || []).forEach(l => { if (l.competency_id != null) levels[l.competency_id] = l.level })
    evaluationForm.value = {
      education_match: saved?.education_match ?? null,
      education_note: saved?.education_note || '',
      experience_match: saved?.experience_match ?? null,
      experience_note: saved?.experience_note || '',
      competency_levels: levels
    }
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    evaluationLoading.value = false
  }
}
async function saveEvaluation() {
  evaluationSaving.value = true
  try {
    const competency_levels = Object.entries(evaluationForm.value.competency_levels)
      .filter(([id, lvl]) => id && lvl != null)
      .map(([id, lvl]) => ({ competency_id: id, level: Number(lvl) }))
    const payload = cleanPayload({
      education_match: evaluationForm.value.education_match,
      education_note: evaluationForm.value.education_note,
      experience_match: evaluationForm.value.experience_match,
      experience_note: evaluationForm.value.experience_note,
      competency_levels
    })
    await api.put(`/api/v1/tenant/recruitment/applications/${applicationId}/assessment`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('applications.evaluation_saved'), life: 3000 })
    loadEvaluation()
    loadMatchScore()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_save')), life: 5000 })
  } finally {
    evaluationSaving.value = false
  }
}
async function loadOptions() {
  try {
    const [candRes, reqRes, empRes, techValRes, mangValRes] = await Promise.allSettled([
      api.get('/api/v1/tenant/recruitment/candidates', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/recruitment/requisitions', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/employees', { params: { per_page: 500 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'technical', per_page: 100 } }),
      api.get('/api/v1/tenant/job-management/values', { params: { type: 'managerial', per_page: 100 } })
    ])
    candidates.value = candRes.status === 'fulfilled' ? (candRes.value.data?.data || []) : []
    requisitions.value = reqRes.status === 'fulfilled' ? (reqRes.value.data?.data || []) : []
    employees.value = empRes.status === 'fulfilled' ? (empRes.value.data?.data || []) : []
    jobValues.value = [
      ...(techValRes.status === 'fulfilled' ? (techValRes.value.data?.data || []) : []),
      ...(mangValRes.status === 'fulfilled' ? (mangValRes.value.data?.data || []) : [])
    ]
  } catch {
    // fail-silent
  }
}

onMounted(async () => {
  loading.value = true
  await Promise.all([loadApplication(), loadOptions()])
  loading.value = false
  loadHistory()
  loadScreenings()
  loadInterviews()
  loadAssessmentParticipation()
  loadMatchScore()
  loadEvaluation()
})
</script>
