<template>
  <div class="space-y-4">
    <template v-if="loading">
      <div class="space-y-3">
        <div class="h-20 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
        <div class="h-64 rounded-lg bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
      </div>
    </template>

    <template v-else-if="session">
      <!-- ── Header info session ── -->
      <div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4">
        <div class="flex items-start justify-between gap-3 flex-wrap">
          <div class="min-w-0">
            <div class="flex items-center gap-2 flex-wrap">
              <h2 class="text-base font-semibold text-gray-800 dark:text-gray-100">{{ session.session_code }}</h2>
              <Tag :value="statusLabel(session.status)" :severity="statusSeverity(session.status)" class="!text-xs !px-1.5 !py-0.5" />
            </div>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">{{ courseName(session.course_id) }}</p>
          </div>
          <div class="flex items-center gap-2 shrink-0">
            <Button :label="t('training.back_to_sessions')" icon="pi pi-arrow-left" size="small" severity="secondary" outlined @click="router.push('/training/sessions')" />
          </div>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mt-4">
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.start_date') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.start_date || '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.end_date') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.end_date || '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.provider_type') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.provider_type ? typeLabel(session.provider_type) : '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.delivery_mode') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.delivery_mode ? t(`training.mode_${session.delivery_mode.toLowerCase()}`) : '-' }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.max_quota') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5">{{ session.max_quota }}</p>
          </div>
          <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-2.5">
            <p class="text-[11px] font-medium text-gray-400 uppercase tracking-wider">{{ t('training.trainer') }}</p>
            <p class="text-sm font-semibold text-gray-800 dark:text-gray-100 mt-0.5 truncate">{{ session.trainer_name || '-' }}</p>
          </div>
        </div>
        <p v-if="session.location || session.meeting_url" class="text-xs text-gray-400 dark:text-gray-500 mt-3">
          <span v-if="session.location"><i class="pi pi-map-marker mr-1"></i>{{ session.location }}</span>
          <span v-if="session.meeting_url" class="ml-3"><i class="pi pi-video mr-1"></i>{{ session.meeting_url }}</span>
        </p>
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

        <!-- ── Overview ── -->
        <div v-if="activeTab === 'overview'" class="p-4">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <!-- Trainers -->
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-user text-emerald-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.trainers') }}</h3>
              </div>
              <div v-if="trainersLoading" class="p-4 space-y-2">
                <div v-for="i in 2" :key="i" class="h-10 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
              </div>
              <div v-else-if="trainers.length" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="st in trainers" :key="st.id" class="flex items-center gap-3 px-3 py-2.5">
                  <div class="w-8 h-8 rounded-full bg-emerald-50 dark:bg-emerald-900/30 flex items-center justify-center shrink-0">
                    <i class="pi pi-user text-emerald-500 text-xs"></i>
                  </div>
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ trainerName(st.trainer_id) }}</p>
                    <p class="text-xs text-gray-400">{{ t(`training.role_${(st.role || 'main').toLowerCase()}`) }}</p>
                  </div>
                </div>
              </div>
              <div v-else class="px-3 py-6 text-center text-sm text-gray-400">{{ t('training.trainers_empty') }}</div>
            </div>

            <!-- Assessments -->
            <div class="rounded-lg border border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                <i class="pi pi-file-check text-indigo-500 text-sm"></i>
                <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.assessments') }}</h3>
              </div>
              <div v-if="assessmentsLoading" class="p-4 space-y-2">
                <div v-for="i in 2" :key="i" class="h-10 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
              </div>
              <div v-else-if="assessments.length" class="divide-y divide-gray-100 dark:divide-gray-800">
                <div v-for="a in assessments" :key="a.id" class="flex items-center gap-3 px-3 py-2.5">
                  <div class="min-w-0 flex-1">
                    <p class="text-sm font-medium text-gray-800 dark:text-gray-100">{{ a.name }}</p>
                    <p class="text-xs text-gray-400">{{ t(`training.assessment_type_${(a.type || 'other').toLowerCase()}`) }} · {{ a.passing_score }}/{{ a.max_score }}</p>
                  </div>
                  <Tag :value="a.is_required ? t('training.required') : t('training.optional')" :severity="a.is_required ? 'danger' : 'secondary'" class="!text-[10px] !px-1.5 !py-0" />
                </div>
              </div>
              <div v-else class="px-3 py-6 text-center text-sm text-gray-400">{{ t('training.assessments_empty') }}</div>
            </div>
          </div>
        </div>

        <!-- ── Participants ── -->
        <div v-if="activeTab === 'participants'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="participants.length" class="text-xs text-gray-400">{{ participants.length }} {{ t('common.items') }}</span>
            <Button :label="t('training.participant_new')" icon="pi pi-user-plus" size="small" @click="openParticipantDialog()" class="ml-auto" />
          </div>
          <DataTable :value="participants" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="participantsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-users text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.participants_empty') }}</p>
              </div>
            </template>
            <Column field="employee_name" :header="t('training.employee')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ employeeName(data.employee_id) }}</span></template>
            </Column>
            <Column field="registration_status" :header="t('training.registration_status')" style="width:140px">
              <template #body="{data}"><Tag :value="regStatusLabel(data.registration_status)" :severity="regStatusSeverity(data.registration_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="attendance_status" :header="t('training.attendance_status')" style="width:120px">
              <template #body="{data}"><Tag :value="attStatusLabel(data.attendance_status)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="completion_status" :header="t('training.completion_status')" style="width:130px">
              <template #body="{data}"><Tag :value="compStatusLabel(data.completion_status)" :severity="compStatusSeverity(data.completion_status)" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="score" :header="t('training.score')" style="width:90px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.score ?? '-' }}</span></template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeParticipant(data)" />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Costs ── -->
        <div v-if="activeTab === 'costs'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <div class="flex items-center gap-2">
              <span v-if="costs.length" class="text-xs text-gray-400">{{ costs.length }} {{ t('common.items') }}</span>
              <span v-if="totalCost > 0" class="text-xs font-medium text-emerald-600 dark:text-emerald-400">{{ t('training.total_estimated') }}: {{ formatMoney(totalCost) }}</span>
            </div>
            <Button :label="t('training.cost_new')" icon="pi pi-plus" size="small" @click="openCostDialog()" class="ml-auto" />
          </div>
          <DataTable :value="costs" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="costsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-dollar text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.costs_empty') }}</p>
              </div>
            </template>
            <Column field="cost_type" :header="t('training.cost_type')" style="width:160px">
              <template #body="{data}"><Tag :value="costTypeLabel(data.cost_type)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="description" :header="t('training.cost_description')">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.description || '-' }}</span></template>
            </Column>
            <Column field="amount" :header="t('training.cost_amount')" style="width:150px">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.amount ? formatMoney(data.amount) : '-' }}</span></template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <div class="flex items-center gap-1 justify-end">
                  <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openCostDialog(data)" />
                  <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeCost(data)" />
                </div>
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Documents ── -->
        <div v-if="activeTab === 'documents'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="documents.length" class="text-xs text-gray-400">{{ documents.length }} {{ t('common.items') }}</span>
            <Button :label="t('training.document_new')" icon="pi pi-plus" size="small" @click="openDocumentDialog()" class="ml-auto" />
          </div>
          <DataTable :value="documents" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="documentsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-file-pdf text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.documents_empty') }}</p>
              </div>
            </template>
            <Column field="document_type" :header="t('training.document_type')" style="width:170px">
              <template #body="{data}"><Tag :value="docTypeLabel(data.document_type)" severity="warning" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="file_name" :header="t('training.document_file_name')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.file_name || '-' }}</span></template>
            </Column>
            <Column field="file_url" :header="t('training.document_file_url')" style="width:220px">
              <template #body="{data}">
                <a v-if="data.file_url" :href="data.file_url" target="_blank" class="text-emerald-600 dark:text-emerald-400 hover:underline text-xs"><i class="pi pi-external-link mr-1"></i>{{ t('common.open') }}</a>
                <span v-else class="text-gray-400">-</span>
              </template>
            </Column>
            <Column :header="t('common.actions')" style="width:70px" frozen alignFrozen="right">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeDocument(data)" />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Materials ── -->
        <div v-if="activeTab === 'materials'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <span v-if="materials.length" class="text-xs text-gray-400">{{ materials.length }} {{ t('common.items') }}</span>
            <Button :label="t('training.material_new')" icon="pi pi-plus" size="small" @click="openMaterialDialog()" class="ml-auto" />
          </div>
          <DataTable :value="materials" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="materialsLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-file text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.materials_empty') }}</p>
              </div>
            </template>
            <Column field="title" :header="t('training.material_title')">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.title }}</span></template>
            </Column>
            <Column field="is_required" :header="t('training.is_required')" style="width:110px">
              <template #body="{data}"><Tag :value="data.is_required ? t('common.yes') : t('common.no')" :severity="data.is_required ? 'danger' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
            </Column>
            <Column field="file_url" :header="t('training.material_file')" style="width:200px">
              <template #body="{data}">
                <a v-if="data.file_url" :href="data.file_url" target="_blank" class="text-emerald-600 dark:text-emerald-400 hover:underline text-xs"><i class="pi pi-external-link mr-1"></i>{{ t('common.open') }}</a>
                <span v-else class="text-gray-400">-</span>
              </template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeMaterial(data)" />
              </template>
            </Column>
          </DataTable>
        </div>

        <!-- ── Evaluation (P2-FE) ── -->
        <div v-if="activeTab === 'evaluation'" class="p-4">
          <div v-if="evalFormLoading" class="p-4 space-y-2">
            <div v-for="i in 3" :key="i" class="h-12 rounded bg-gray-100 dark:bg-gray-700/50 animate-pulse"></div>
          </div>
          <template v-else>
            <!-- Tidak ada form evaluasi untuk session ini -->
            <div v-if="!evaluationForm" class="flex flex-col items-center justify-center py-10 text-gray-400 dark:text-gray-500">
              <i class="pi pi-clipboard text-3xl mb-2 opacity-50"></i>
              <p class="text-sm font-medium">{{ t('training.eval_form_not_found') }}</p>
              <p class="text-xs mt-1 mb-3">{{ t('training.eval_form_not_found_hint') }}</p>
              <Button :label="t('training.eval_form_create')" icon="pi pi-plus" size="small" @click="openEvalFormDialog()" />
            </div>

            <!-- Form evaluasi + pertanyaan -->
            <div v-else class="space-y-4">
              <div class="rounded-lg border border-gray-200 dark:border-gray-700 p-4">
                <div class="flex items-start justify-between gap-3 flex-wrap">
                  <div class="min-w-0">
                    <div class="flex items-center gap-2 flex-wrap">
                      <h3 class="text-sm font-semibold text-gray-800 dark:text-gray-100">{{ evaluationForm.form.name }}</h3>
                      <Tag :value="evaluationForm.form.is_active ? t('training.is_active') : t('training.inactive')" :severity="evaluationForm.form.is_active ? 'success' : 'secondary'" class="!text-[10px] !px-1.5 !py-0" />
                    </div>
                    <p class="text-xs text-gray-400 mt-0.5">{{ t('training.eval_questions_count', { count: evaluationForm.questions.length }) }}</p>
                  </div>
                  <div class="flex items-center gap-2 shrink-0">
                    <Button :label="t('training.eval_question_new')" icon="pi pi-plus" size="small" severity="secondary" outlined @click="openEvalQuestionDialog()" />
                  </div>
                </div>
              </div>

              <!-- Daftar pertanyaan -->
              <DataTable :value="evaluationForm.questions || []" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                <template #empty>
                  <div class="text-center py-6 text-sm text-gray-400">{{ t('training.eval_questions_empty') }}</div>
                </template>
                <Column field="sort_order" :header="t('training.sort_order')" style="width:80px">
                  <template #body="{data}"><span class="text-gray-400">{{ data.sort_order }}</span></template>
                </Column>
                <Column field="question" :header="t('training.question')">
                  <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ data.question }}</span></template>
                </Column>
                <Column field="question_type" :header="t('training.question_type')" style="width:160px">
                  <template #body="{data}"><Tag :value="qTypeLabel(data.question_type)" severity="info" class="!text-xs !px-1.5 !py-0.5" /></template>
                </Column>
                <Column field="is_required" :header="t('training.is_required')" style="width:90px">
                  <template #body="{data}"><Tag :value="data.is_required ? t('common.yes') : t('common.no')" :severity="data.is_required ? 'danger' : 'secondary'" class="!text-xs !px-1.5 !py-0.5" /></template>
                </Column>
                <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
                  <template #body="{data}">
                    <div class="flex items-center gap-1 justify-end">
                      <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEvalQuestionDialog(data)" />
                      <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeEvalQuestion(data)" />
                    </div>
                  </template>
                </Column>
              </DataTable>

              <!-- Submit jawaban per peserta -->
              <div class="rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div class="flex items-center gap-2 px-3 py-2.5 border-b border-gray-100 dark:border-gray-800">
                  <i class="pi pi-pen-to-square text-emerald-500 text-sm"></i>
                  <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-200">{{ t('training.eval_submit_answers') }}</h3>
                </div>
                <div class="p-4 space-y-3">
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <FormRow :label="t('training.participant')">
                      <SelectLabel v-model="answerParticipantId" :options="participantOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" showClear />
                    </FormRow>
                  </div>
                  <div v-if="answerParticipantId && (evaluationForm.questions || []).length" class="space-y-3">
                    <div v-for="q in evaluationForm.questions" :key="q.id" class="rounded-lg border border-gray-200 dark:border-gray-700 p-3">
                      <p class="text-sm font-medium text-gray-800 dark:text-gray-100">
                        {{ q.question }}
                        <span v-if="q.is_required" class="text-rose-500">*</span>
                      </p>
                      <!-- RATING -->
                      <div v-if="q.question_type === 'RATING'" class="mt-2 flex items-center gap-1">
                        <button
                          v-for="n in 5"
                          :key="n"
                          type="button"
                          class="w-8 h-8 rounded-lg flex items-center justify-center transition-colors"
                          :class="Number(answers[q.id]) >= n ? 'bg-amber-100 dark:bg-amber-500/20 text-amber-500' : 'bg-gray-100 dark:bg-gray-700/50 text-gray-300 dark:text-gray-600 hover:text-gray-400'"
                          @click="answers[q.id] = String(n)"
                        >
                          <i class="pi pi-star-fill text-sm"></i>
                        </button>
                      </div>
                      <!-- TEXT -->
                      <div v-else-if="q.question_type === 'TEXT'" class="mt-2">
                        <TextInput v-model="answers[q.id]" textarea :rows="2" :placeholder="t('training.eval_answer_placeholder')" />
                      </div>
                      <!-- SINGLE_CHOICE / MULTIPLE_CHOICE -->
                      <div v-else class="mt-2">
                        <SelectLabel
                          v-if="q.question_type === 'SINGLE_CHOICE'"
                          v-model="answers[q.id]"
                          :options="choiceOptions"
                          optionLabel="label"
                          optionValue="value"
                          :placeholder="t('common.select')"
                          showClear
                          class="!max-w-sm"
                        />
                        <div v-else class="flex items-center gap-4 flex-wrap">
                          <label v-for="opt in choiceOptions" :key="opt.value" class="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-200 cursor-pointer">
                            <input type="checkbox" :checked="(answers[q.id] || '').split(',').includes(opt.value)" class="accent-emerald-500" @change="toggleChoice(q.id, opt.value)" />
                            {{ opt.label }}
                          </label>
                        </div>
                      </div>
                    </div>
                    <div class="flex items-center justify-end">
                      <Button :label="t('training.eval_submit')" icon="pi pi-check" size="small" :loading="answerSaving" :disabled="answerSaving" @click="submitAnswers" />
                    </div>
                  </div>
                  <p v-else-if="answerParticipantId" class="text-xs text-gray-400">{{ t('training.eval_no_questions_hint') }}</p>
                </div>
              </div>
            </div>
          </template>
        </div>

        <!-- ── Effectiveness (P2-FE) ── -->
        <div v-if="activeTab === 'effectiveness'" class="p-4">
          <div class="flex items-center justify-between gap-2 flex-wrap mb-3">
            <div class="flex items-center gap-2 flex-wrap">
              <SelectLabel v-model="effectivenessParticipantId" :options="participantOptions" optionLabel="label" optionValue="value" filter :placeholder="t('training.filter_all_participants')" class="!w-64" showClear @update:modelValue="onEffectivenessFilterChange" />
              <span v-if="effectivenessList.length" class="text-xs text-gray-400">{{ effectivenessList.length }} {{ t('common.items') }}</span>
            </div>
            <Button :label="t('training.effectiveness_new')" icon="pi pi-plus" size="small" class="ml-auto" @click="openEffectivenessDialog()" />
          </div>
          <DataTable :value="effectivenessList" size="small" class="!text-sm border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden" :loading="effectivenessLoading">
            <template #empty>
              <div class="flex flex-col items-center justify-center py-8 text-gray-400 dark:text-gray-500">
                <i class="pi pi-chart-line text-2xl mb-2 opacity-50"></i>
                <p class="text-sm">{{ t('training.effectiveness_empty') }}</p>
              </div>
            </template>
            <Column field="participant_id" :header="t('training.participant')" style="width:200px">
              <template #body="{data}"><span class="text-gray-800 dark:text-gray-100 font-medium">{{ participantLabel(data.participant_id) }}</span></template>
            </Column>
            <Column field="assessment_date" :header="t('training.assessment_date')" style="width:130px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.assessment_date || '-' }}</span></template>
            </Column>
            <Column field="before_score" :header="t('training.before_score')" style="width:110px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.before_score ?? '-' }}</span></template>
            </Column>
            <Column field="after_score" :header="t('training.after_score')" style="width:110px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300">{{ data.after_score ?? '-' }}</span></template>
            </Column>
            <Column field="effectiveness_score" :header="t('training.effectiveness_score')" style="width:150px">
              <template #body="{data}">
                <Tag v-if="data.effectiveness_score != null" :value="data.effectiveness_score" :severity="effectivenessSeverity(data.effectiveness_score)" class="!text-xs !px-1.5 !py-0.5" />
                <span v-else class="text-gray-400">-</span>
              </template>
            </Column>
            <Column field="remarks" :header="t('training.remarks')" style="width:220px">
              <template #body="{data}"><span class="text-gray-600 dark:text-gray-300 line-clamp-1">{{ data.remarks || '-' }}</span></template>
            </Column>
            <Column :header="t('common.actions')" style="width:90px" frozen alignFrozen="right">
              <template #body="{data}">
                <div class="flex items-center gap-1 justify-end">
                  <Button icon="pi pi-pencil" size="small" text severity="secondary" v-tooltip.left="t('common.edit')" @click="openEffectivenessDialog(data)" />
                  <Button icon="pi pi-trash" size="small" text severity="danger" v-tooltip.left="t('common.delete')" @click="removeEffectiveness(data)" />
                </div>
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- ── Dialog: create evaluation form ── -->
      <Dialog v-model:visible="evalFormDialogVisible" :header="t('training.eval_form_create')" modal :style="{ width: '480px' }" @hide="resetEvalForm">
        <div class="space-y-4">
          <FormRow :label="t('training.eval_form_name')" required :errors="errors?.name">
            <TextInput v-model="evalForm.name" maxlength="200" :placeholder="t('training.eval_form_name')" :class="{ 'p-invalid': errors?.name }" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="evalFormDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="evalFormSaving" :disabled="evalFormSaving" @click="handleCreateEvalForm" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: evaluation question ── -->
      <Dialog v-model:visible="evalQuestionDialogVisible" :header="evalQuestionEditing ? t('training.eval_question_edit') : t('training.eval_question_new')" modal :style="{ width: '520px' }" @hide="resetEvalQuestionForm">
        <div class="space-y-4">
          <FormRow :label="t('training.question')" required :errors="errors?.question">
            <TextInput v-model="evalQuestionForm.question" textarea :rows="2" :placeholder="t('training.question')" :class="{ 'p-invalid': errors?.question }" />
          </FormRow>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FormRow :label="t('training.question_type')" required :errors="errors?.question_type">
              <SelectLabel v-model="evalQuestionForm.question_type" :options="qTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.question_type }" />
            </FormRow>
            <FormRow :label="t('training.sort_order')">
              <InputNumber v-model="evalQuestionForm.sort_order" class="!w-full" :min="0" size="small" />
            </FormRow>
          </div>
          <FormRow :label="t('training.is_required')">
            <ToggleSwitch v-model="evalQuestionForm.is_required" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="evalQuestionDialogVisible = false" />
            <Button :label="evalQuestionEditing ? t('common.update') : t('common.save')" size="small" :loading="evalQuestionSaving" :disabled="evalQuestionSaving" @click="handleSaveEvalQuestion" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: effectiveness assessment ── -->
      <Dialog v-model:visible="effectivenessDialogVisible" :header="effectivenessEditing ? t('training.effectiveness_edit') : t('training.effectiveness_new')" modal :style="{ width: '540px' }" @hide="resetEffectivenessForm">
        <div class="space-y-4">
          <FormRow :label="t('training.participant')" required :errors="errors?.participant_id">
            <SelectLabel v-model="effectivenessForm.participant_id" :options="participantOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.participant_id }" />
          </FormRow>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FormRow :label="t('training.assessment_date')" required :errors="errors?.assessment_date">
              <DateInput v-model="effectivenessForm.assessment_date" :class="{ 'p-invalid': errors?.assessment_date }" />
            </FormRow>
            <FormRow :label="t('training.effectiveness_score')">
              <InputNumber v-model="effectivenessForm.effectiveness_score" class="!w-full" :min="0" :max="100" size="small" />
            </FormRow>
          </div>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <FormRow :label="t('training.before_score')">
              <InputNumber v-model="effectivenessForm.before_score" class="!w-full" :min="0" :max="100" size="small" />
            </FormRow>
            <FormRow :label="t('training.after_score')">
              <InputNumber v-model="effectivenessForm.after_score" class="!w-full" :min="0" :max="100" size="small" />
            </FormRow>
          </div>
          <FormRow :label="t('training.remarks')">
            <TextInput v-model="effectivenessForm.remarks" textarea :rows="2" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="effectivenessDialogVisible = false" />
            <Button :label="effectivenessEditing ? t('common.update') : t('common.save')" size="small" :loading="effectivenessSaving" :disabled="effectivenessSaving" @click="handleSaveEffectiveness" />
          </div>
        </template>
      </Dialog>
      <Dialog v-model:visible="participantDialogVisible" :header="t('training.participant_new')" modal :style="{ width: '480px' }">
        <div class="space-y-4">
          <FormRow :label="t('training.employee')" required :errors="errors?.employee_id">
            <SelectLabel v-model="participantForm.employee_id" :options="employeeOptions" optionLabel="label" optionValue="value" filter :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.employee_id }" />
          </FormRow>
          <FormRow :label="t('training.registration_status')">
            <SelectLabel v-model="participantForm.registration_status" :options="regStatusOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="participantDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="participantSaving" :disabled="participantSaving" @click="handleAddParticipant" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: cost ── -->
      <Dialog v-model:visible="costDialogVisible" :header="costEditing ? t('training.cost_edit') : t('training.cost_new')" modal :style="{ width: '480px' }" @hide="resetCostForm">
        <div class="space-y-4">
          <FormRow :label="t('training.cost_type')" required :errors="errors?.cost_type">
            <SelectLabel v-model="costForm.cost_type" :options="costTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.cost_type }" />
          </FormRow>
          <FormRow :label="t('training.cost_description')">
            <TextInput v-model="costForm.description" textarea :rows="2" />
          </FormRow>
          <FormRow :label="t('training.cost_amount')">
            <InputNumber v-model="costForm.amount" class="!w-full" :min="0" mode="currency" currency="IDR" locale="id-ID" size="small" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="costDialogVisible = false" />
            <Button :label="costEditing ? t('common.update') : t('common.save')" size="small" :loading="costSaving" :disabled="costSaving" @click="handleSaveCost" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: document ── -->
      <Dialog v-model:visible="documentDialogVisible" :header="t('training.document_new')" modal :style="{ width: '520px' }" @hide="resetDocumentForm">
        <div class="space-y-4">
          <FormRow :label="t('training.document_type')" required :errors="errors?.document_type">
            <SelectLabel v-model="documentForm.document_type" :options="documentTypeOptions" optionLabel="label" optionValue="value" :placeholder="t('common.select')" :class="{ 'p-invalid': errors?.document_type }" />
          </FormRow>
          <FormRow :label="t('training.document_file_name')">
            <TextInput v-model="documentForm.file_name" maxlength="255" :placeholder="t('training.document_file_name')" />
          </FormRow>
          <FormRow :label="t('training.document_file_url')" required :errors="errors?.file_url">
            <TextInput v-model="documentForm.file_url" maxlength="500" :placeholder="t('training.document_file_url_placeholder')" :class="{ 'p-invalid': errors?.file_url }" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="documentDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="documentSaving" :disabled="documentSaving" @click="handleAddDocument" />
          </div>
        </template>
      </Dialog>

      <!-- ── Dialog: add material ── -->
      <Dialog v-model:visible="materialDialogVisible" :header="t('training.material_new')" modal :style="{ width: '520px' }">
        <div class="space-y-4">
          <FormRow :label="t('training.material_title')" required :errors="errors?.title">
            <TextInput v-model="materialForm.title" maxlength="200" :placeholder="t('training.material_title')" :class="{ 'p-invalid': errors?.title }" />
          </FormRow>
          <FormRow :label="t('training.material_description')">
            <TextInput v-model="materialForm.description" textarea :rows="2" />
          </FormRow>
          <FormRow :label="t('training.material_file')">
            <TextInput v-model="materialForm.file_url" maxlength="500" :placeholder="t('training.material_file_placeholder')" />
          </FormRow>
          <FormRow :label="t('training.is_required')">
            <ToggleSwitch v-model="materialForm.is_required" />
          </FormRow>
        </div>
        <template #footer>
          <div class="flex items-center justify-end gap-2">
            <Button :label="t('common.cancel')" severity="secondary" outlined size="small" @click="materialDialogVisible = false" />
            <Button :label="t('common.save')" size="small" :loading="materialSaving" :disabled="materialSaving" @click="handleAddMaterial" />
          </div>
        </template>
      </Dialog>
    </template>

    <div v-else class="text-center py-12 text-gray-400">
      <i class="pi pi-exclamation-triangle text-3xl mb-2 opacity-50"></i>
      <p class="text-sm">{{ t('training.session_not_found') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from '@/composables/useI18n'
import { getErrorMessage, getValidationErrors } from '@/services/responseHandler'
import api from '@/services/api'

import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import InputNumber from 'primevue/inputnumber'
import Tag from 'primevue/tag'
import Dialog from 'primevue/dialog'
import FormRow from '@/components/FormRow.vue'
import TextInput from '@/components/TextInput.vue'
import ToggleSwitch from '@/components/ToggleSwitch.vue'
import SelectLabel from '@/components/SelectLabel.vue'
import DateInput from '@/components/DateInput.vue'

const { t } = useI18n()
const toast = useToast()
const router = useRouter()
const route = useRoute()

const sessionId = route.params.id
const session = ref(null)
const loading = ref(true)
const activeTab = ref('overview')

const trainers = ref([])
const trainersLoading = ref(false)
const assessments = ref([])
const assessmentsLoading = ref(false)
const participants = ref([])
const participantsLoading = ref(false)
const materials = ref([])
const materialsLoading = ref(false)
const costs = ref([])
const costsLoading = ref(false)
const documents = ref([])
const documentsLoading = ref(false)

const costDialogVisible = ref(false)
const costEditing = ref(false)
const costEditingId = ref(null)
const costSaving = ref(false)
const costForm = ref({ cost_type: null, description: '', amount: null })
const documentDialogVisible = ref(false)
const documentSaving = ref(false)
const documentForm = ref({ document_type: null, file_name: '', file_url: '' })

const courses = ref([])
const trainersMaster = ref([])
const employees = ref([])

const participantDialogVisible = ref(false)
const participantSaving = ref(false)
const participantForm = ref({ employee_id: null, registration_status: 'REGISTERED' })
const materialDialogVisible = ref(false)
const materialSaving = ref(false)
const materialForm = ref({ title: '', description: '', file_url: '', is_required: false })
const errors = ref({})

const tabs = [
  { key: 'overview', labelKey: 'training.tab_overview' },
  { key: 'participants', labelKey: 'training.participants' },
  { key: 'materials', labelKey: 'training.materials' },
  { key: 'costs', labelKey: 'training.tab_costs' },
  { key: 'documents', labelKey: 'training.tab_documents' },
  { key: 'evaluation', labelKey: 'training.tab_evaluation' },
  { key: 'effectiveness', labelKey: 'training.tab_effectiveness' }
]

const employeeOptions = computed(() => employees.value.map(e => ({ label: `${e.name} (${e.employee_id})`, value: e.id })))

// Status pendaftaran saat create — CANCELLED tidak masuk akal untuk enrollment baru.
const regStatusOptions = computed(() => ['NOMINATED', 'REQUESTED', 'APPROVED', 'REGISTERED', 'WAITLISTED'].map(v => ({ label: regStatusLabel(v), value: v })))

function courseName(id) {
  return courses.value.find(c => c.id === id)?.name || id
}
function trainerName(id) {
  const tr = trainersMaster.value.find(x => x.id === id)
  return tr ? tr.name : id
}
function employeeName(id) {
  return employees.value.find(e => e.id === id)?.name || id
}

function typeLabel(type) {
  const key = `training.type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function statusLabel(status) {
  const key = `training.status_${String(status || '').toLowerCase()}`
  return t(key) !== key ? t(key) : status
}
function statusSeverity(status) {
  switch (status) {
    case 'DRAFT': return 'secondary'
    case 'SCHEDULED': return 'info'
    case 'REGISTRATION_OPEN': return 'success'
    case 'FULL': return 'warning'
    case 'IN_PROGRESS': return 'info'
    case 'COMPLETED': return 'success'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function regStatusLabel(s) {
  const key = `training.reg_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function regStatusSeverity(s) {
  switch (s) {
    case 'REGISTERED': return 'success'
    case 'APPROVED': return 'info'
    case 'NOMINATED': return 'warning'
    case 'REQUESTED': return 'info'
    case 'WAITLISTED': return 'warning'
    case 'CANCELLED': return 'danger'
    default: return 'secondary'
  }
}
function attStatusLabel(s) {
  const key = `training.att_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusLabel(s) {
  const key = `training.comp_status_${String(s || '').toLowerCase()}`
  return t(key) !== key ? t(key) : s
}
function compStatusSeverity(s) {
  switch (s) {
    case 'COMPLETED': return 'success'
    case 'FAILED': return 'danger'
    case 'IN_PROGRESS': return 'info'
    default: return 'secondary'
  }
}

async function loadSession() {
  loading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}`)
    session.value = res.data?.data || null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.failed_to_load')), life: 4000 })
  } finally {
    loading.value = false
  }
}

async function loadTrainers() {
  trainersLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/trainers`)
    trainers.value = res.data?.data || []
  } catch { trainers.value = [] } finally { trainersLoading.value = false }
}

async function loadAssessments() {
  assessmentsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/assessments`)
    assessments.value = res.data?.data || []
  } catch { assessments.value = [] } finally { assessmentsLoading.value = false }
}

async function loadParticipants() {
  participantsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/participants', { params: { session_id: sessionId, per_page: 200 } })
    participants.value = res.data?.data || []
  } catch { participants.value = [] } finally { participantsLoading.value = false }
}

async function loadMaterials() {
  materialsLoading.value = true
  try {
    const res = await api.get('/api/v1/tenant/trainings/materials', { params: { session_id: sessionId, per_page: 200 } })
    materials.value = res.data?.data || []
  } catch { materials.value = [] } finally { materialsLoading.value = false }
}

async function loadCosts() {
  costsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/costs`)
    costs.value = res.data?.data || []
  } catch { costs.value = [] } finally { costsLoading.value = false }
}

async function loadDocuments() {
  documentsLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/documents`)
    documents.value = res.data?.data || []
  } catch { documents.value = [] } finally { documentsLoading.value = false }
}

async function loadReferences() {
  const [cRes, tRes, eRes] = await Promise.allSettled([
    api.get('/api/v1/tenant/trainings/courses', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/trainings/trainers', { params: { per_page: 500 } }),
    api.get('/api/v1/tenant/employees', { params: { per_page: 500 } })
  ])
  courses.value = cRes.status === 'fulfilled' ? (cRes.value.data?.data || []) : []
  trainersMaster.value = tRes.status === 'fulfilled' ? (tRes.value.data?.data || []) : []
  employees.value = eRes.status === 'fulfilled' ? (eRes.value.data?.data || []) : []
}

function openParticipantDialog() {
  errors.value = {}
  participantForm.value = { employee_id: null, registration_status: 'REGISTERED' }
  participantDialogVisible.value = true
}

async function handleAddParticipant() {
  errors.value = {}
  if (!participantForm.value.employee_id) { errors.value = { employee_id: t('form.required') }; return }
  participantSaving.value = true
  try {
    await api.post('/api/v1/tenant/trainings/participants', {
      session_id: sessionId,
      employee_id: participantForm.value.employee_id,
      registration_status: participantForm.value.registration_status
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    participantDialogVisible.value = false
    await loadParticipants()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    participantSaving.value = false
  }
}

async function removeParticipant(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/participants/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadParticipants()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

function openMaterialDialog() {
  errors.value = {}
  materialForm.value = { title: '', description: '', file_url: '', is_required: false }
  materialDialogVisible.value = true
}

async function handleAddMaterial() {
  errors.value = {}
  if (!materialForm.value.title?.trim()) { errors.value = { title: t('form.required') }; return }
  materialSaving.value = true
  try {
    await api.post('/api/v1/tenant/trainings/materials', {
      session_id: sessionId,
      title: materialForm.value.title.trim(),
      description: materialForm.value.description?.trim() || '',
      file_url: materialForm.value.file_url?.trim() || '',
      is_required: materialForm.value.is_required
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    materialDialogVisible.value = false
    await loadMaterials()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    materialSaving.value = false
  }
}

async function removeMaterial(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/materials/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadMaterials()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

const costTypeOptions = computed(() => ['TRAINER', 'PROVIDER', 'VENUE', 'MATERIAL', 'CERTIFICATION', 'TRAVEL', 'ACCOMMODATION', 'OTHER'].map(v => ({ label: costTypeLabel(v), value: v })))
const documentTypeOptions = computed(() => ['PROPOSAL', 'QUOTATION', 'ATTENDANCE_SHEET', 'INVOICE', 'CONTRACT', 'TRAINING_REPORT', 'OTHER'].map(v => ({ label: docTypeLabel(v), value: v })))

const totalCost = computed(() => costs.value.reduce((sum, c) => sum + (Number(c.amount) || 0), 0))

function costTypeLabel(type) {
  const key = `training.cost_type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function docTypeLabel(type) {
  const key = `training.document_type_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function formatMoney(v) {
  try { return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(v) } catch { return v }
}

function openCostDialog(item) {
  errors.value = {}
  costEditing.value = !!item
  costEditingId.value = item?.id || null
  costForm.value = item
    ? { cost_type: item.cost_type || null, description: item.description || '', amount: item.amount ?? null }
    : { cost_type: null, description: '', amount: null }
  costDialogVisible.value = true
}

function resetCostForm() {
  costForm.value = { cost_type: null, description: '', amount: null }
  errors.value = {}
  costEditing.value = false
  costEditingId.value = null
}

async function handleSaveCost() {
  errors.value = {}
  if (!costForm.value.cost_type) { errors.value = { cost_type: t('form.required') }; return }
  costSaving.value = true
  try {
    const payload = {
      cost_type: costForm.value.cost_type,
      description: costForm.value.description?.trim() || '',
      amount: costForm.value.amount ?? null
    }
    if (costEditing.value) {
      await api.put(`/api/v1/tenant/trainings/session-costs/${costEditingId.value}`, payload)
    } else {
      await api.post(`/api/v1/tenant/trainings/sessions/${sessionId}/costs`, payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    costDialogVisible.value = false
    await loadCosts()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    costSaving.value = false
  }
}

async function removeCost(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/session-costs/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadCosts()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

function openDocumentDialog() {
  errors.value = {}
  documentForm.value = { document_type: null, file_name: '', file_url: '' }
  documentDialogVisible.value = true
}

function resetDocumentForm() {
  documentForm.value = { document_type: null, file_name: '', file_url: '' }
  errors.value = {}
}

async function handleAddDocument() {
  errors.value = {}
  if (!documentForm.value.document_type) { errors.value = { document_type: t('form.required') }; return }
  if (!documentForm.value.file_url?.trim()) { errors.value = { file_url: t('form.required') }; return }
  documentSaving.value = true
  try {
    await api.post(`/api/v1/tenant/trainings/sessions/${sessionId}/documents`, {
      document_type: documentForm.value.document_type,
      file_name: documentForm.value.file_name?.trim() || '',
      file_url: documentForm.value.file_url.trim()
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    documentDialogVisible.value = false
    await loadDocuments()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    documentSaving.value = false
  }
}

async function removeDocument(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/documents/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadDocuments()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

// ════ P2-FE: Evaluation form & answers ════
const evalFormLoading = ref(false)
const evaluationForm = ref(null)
const evalFormDialogVisible = ref(false)
const evalFormSaving = ref(false)
const evalForm = ref({ name: '' })
const evalQuestionDialogVisible = ref(false)
const evalQuestionEditing = ref(false)
const evalQuestionEditingId = ref(null)
const evalQuestionSaving = ref(false)
const evalQuestionForm = ref({ question: '', question_type: 'RATING', sort_order: null, is_required: false })
const answerParticipantId = ref(null)
const answers = ref({})
const answerSaving = ref(false)

const qTypeOptions = computed(() => ['RATING', 'TEXT', 'SINGLE_CHOICE', 'MULTIPLE_CHOICE'].map(v => ({ label: qTypeLabel(v), value: v })))
const choiceOptions = computed(() => [
  { label: t('training.choice_excellent'), value: 'EXCELLENT' },
  { label: t('training.choice_good'), value: 'GOOD' },
  { label: t('training.choice_fair'), value: 'FAIR' },
  { label: t('training.choice_poor'), value: 'POOR' }
])
const participantOptions = computed(() => participants.value.map(p => ({ label: participantLabel(p.id), value: p.id })))

function qTypeLabel(type) {
  const key = `training.qtype_${String(type || '').toLowerCase()}`
  return t(key) !== key ? t(key) : type
}
function participantLabel(id) {
  const p = participants.value.find(x => x.id === id)
  if (!p) return id
  return employeeName(p.employee_id)
}

async function loadEvaluationForm() {
  evalFormLoading.value = true
  try {
    const res = await api.get(`/api/v1/tenant/trainings/sessions/${sessionId}/evaluation-form`)
    evaluationForm.value = res.data?.data || null
  } catch {
    evaluationForm.value = null
  } finally {
    evalFormLoading.value = false
  }
}

function openEvalFormDialog() {
  errors.value = {}
  evalForm.value = { name: '' }
  evalFormDialogVisible.value = true
}

function resetEvalForm() {
  evalForm.value = { name: '' }
  errors.value = {}
}

async function handleCreateEvalForm() {
  errors.value = {}
  if (!evalForm.value.name?.trim()) { errors.value = { name: t('form.required') }; return }
  evalFormSaving.value = true
  try {
    await api.post('/api/v1/tenant/trainings/evaluation-forms', {
      session_id: sessionId,
      name: evalForm.value.name.trim(),
      is_active: true
    })
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    evalFormDialogVisible.value = false
    await loadEvaluationForm()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    evalFormSaving.value = false
  }
}

function openEvalQuestionDialog(item) {
  errors.value = {}
  evalQuestionEditing.value = !!item
  evalQuestionEditingId.value = item?.id || null
  evalQuestionForm.value = item
    ? { question: item.question || '', question_type: item.question_type || 'RATING', sort_order: item.sort_order ?? null, is_required: item.is_required }
    : { question: '', question_type: 'RATING', sort_order: (evaluationForm.value?.questions?.length || 0) + 1, is_required: false }
  evalQuestionDialogVisible.value = true
}

function resetEvalQuestionForm() {
  evalQuestionForm.value = { question: '', question_type: 'RATING', sort_order: null, is_required: false }
  errors.value = {}
  evalQuestionEditing.value = false
  evalQuestionEditingId.value = null
}

async function handleSaveEvalQuestion() {
  errors.value = {}
  if (!evalQuestionForm.value.question?.trim()) { errors.value = { question: t('form.required') }; return }
  if (!evalQuestionForm.value.question_type) { errors.value = { question_type: t('form.required') }; return }
  evalQuestionSaving.value = true
  try {
    const payload = {
      question: evalQuestionForm.value.question.trim(),
      question_type: evalQuestionForm.value.question_type,
      sort_order: evalQuestionForm.value.sort_order ?? null,
      is_required: evalQuestionForm.value.is_required
    }
    const formId = evaluationForm.value?.form?.id
    if (evalQuestionEditing.value) {
      await api.put(`/api/v1/tenant/trainings/evaluation-questions/${evalQuestionEditingId.value}`, payload)
    } else if (formId) {
      await api.post(`/api/v1/tenant/trainings/evaluation-forms/${formId}/questions`, payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    evalQuestionDialogVisible.value = false
    await loadEvaluationForm()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    evalQuestionSaving.value = false
  }
}

async function removeEvalQuestion(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/evaluation-questions/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadEvaluationForm()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

function toggleChoice(questionId, value) {
  const current = (answers.value[questionId] || '').split(',').filter(Boolean)
  const idx = current.indexOf(value)
  if (idx >= 0) current.splice(idx, 1)
  else current.push(value)
  answers.value[questionId] = current.join(',')
}

async function submitAnswers() {
  if (!answerParticipantId.value || !evaluationForm.value) return
  const formId = evaluationForm.value.form.id
  const payload = {
    answers: (evaluationForm.value.questions || [])
      .filter(q => answers.value[q.id] != null && answers.value[q.id] !== '')
      .map(q => ({ question_id: q.id, answer: String(answers.value[q.id]) }))
  }
  if (!payload.answers.length) {
    toast.add({ severity: 'warn', summary: t('message.warning'), detail: t('training.eval_answers_empty'), life: 3000 })
    return
  }
  answerSaving.value = true
  try {
    await api.post(`/api/v1/tenant/trainings/evaluation-forms/${formId}/participants/${answerParticipantId.value}/answers`, payload)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    answers.value = {}
    answerParticipantId.value = null
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  } finally {
    answerSaving.value = false
  }
}

// ════ P2-FE: Effectiveness assessment ════
const effectivenessParticipantId = ref(null)
const effectivenessList = ref([])
const effectivenessLoading = ref(false)
const effectivenessDialogVisible = ref(false)
const effectivenessEditing = ref(false)
const effectivenessEditingId = ref(null)
const effectivenessSaving = ref(false)
const effectivenessForm = ref(defaultEffectivenessForm())

function defaultEffectivenessForm() {
  return { participant_id: null, assessment_date: null, before_score: null, after_score: null, effectiveness_score: null, remarks: '' }
}

function effectivenessSeverity(score) {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

async function loadEffectiveness() {
  effectivenessLoading.value = true
  try {
    const params = { per_page: 200 }
    if (effectivenessParticipantId.value) params.participant_id = effectivenessParticipantId.value
    const res = await api.get('/api/v1/tenant/trainings/effectiveness', { params })
    let rows = res.data?.data || []
    if (effectivenessParticipantId.value) {
      effectivenessList.value = rows
    } else {
      // Filter ke peserta session ini (list global tidak punya filter session di BE)
      const sessionPartIds = new Set(participants.value.map(p => p.id))
      effectivenessList.value = rows.filter(r => sessionPartIds.has(r.participant_id))
    }
  } catch {
    effectivenessList.value = []
  } finally {
    effectivenessLoading.value = false
  }
}

function onEffectivenessFilterChange() {
  loadEffectiveness()
}

function openEffectivenessDialog(item) {
  errors.value = {}
  effectivenessEditing.value = !!item
  effectivenessEditingId.value = item?.id || null
  effectivenessForm.value = item
    ? {
        participant_id: item.participant_id || null,
        assessment_date: item.assessment_date || null,
        before_score: item.before_score ?? null,
        after_score: item.after_score ?? null,
        effectiveness_score: item.effectiveness_score ?? null,
        remarks: item.remarks || ''
      }
    : defaultEffectivenessForm()
  effectivenessDialogVisible.value = true
}

function resetEffectivenessForm() {
  effectivenessForm.value = defaultEffectivenessForm()
  errors.value = {}
  effectivenessEditing.value = false
  effectivenessEditingId.value = null
}

async function handleSaveEffectiveness() {
  errors.value = {}
  if (!effectivenessForm.value.participant_id) { errors.value = { participant_id: t('form.required') }; return }
  if (!effectivenessForm.value.assessment_date) { errors.value = { assessment_date: t('form.required') }; return }
  effectivenessSaving.value = true
  try {
    const payload = {
      participant_id: effectivenessForm.value.participant_id,
      assessment_date: effectivenessForm.value.assessment_date,
      before_score: effectivenessForm.value.before_score ?? null,
      after_score: effectivenessForm.value.after_score ?? null,
      effectiveness_score: effectivenessForm.value.effectiveness_score ?? null,
      remarks: effectivenessForm.value.remarks?.trim() || ''
    }
    if (effectivenessEditing.value) {
      await api.put(`/api/v1/tenant/trainings/effectiveness/${effectivenessEditingId.value}`, payload)
    } else {
      await api.post('/api/v1/tenant/trainings/effectiveness', payload)
    }
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.saved'), life: 3000 })
    effectivenessDialogVisible.value = false
    await loadEffectiveness()
  } catch (e) {
    const fieldErrors = getValidationErrors(e)
    if (Object.keys(fieldErrors).length > 0) {
      errors.value = fieldErrors
    } else {
      toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
    }
  } finally {
    effectivenessSaving.value = false
  }
}

async function removeEffectiveness(item) {
  try {
    await api.delete(`/api/v1/tenant/trainings/effectiveness/${item.id}`)
    toast.add({ severity: 'success', summary: t('message.success'), detail: t('message.deleted'), life: 3000 })
    await loadEffectiveness()
  } catch (e) {
    toast.add({ severity: 'error', summary: t('message.error'), detail: getErrorMessage(e, t('message.operation_failed')), life: 4000 })
  }
}

onMounted(() => {
  loadSession()
  loadReferences()
  loadTrainers()
  loadAssessments()
  loadMaterials()
  loadCosts()
  loadDocuments()
  loadEvaluationForm()
  // Effectiveness membutuhkan daftar peserta session untuk filter — jalankan setelah participants termuat
  loadParticipants().finally(() => loadEffectiveness())
})
</script>
