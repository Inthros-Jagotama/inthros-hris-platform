<template>
  <div class="space-y-4">
    <div>
      <h2 class="text-lg font-semibold text-gray-800 dark:text-gray-100">{{ t('job_management.potency_competencies') }}</h2>
      <p class="text-sm text-gray-500 dark:text-gray-400">{{ t('job_management.potency_description') }}</p>
    </div>

    <!-- Card: Potensi yang harus dimiliki — komponen terpisah (psychological multi-select + tabel) -->
    <PsychologicalPotencyCard :org-id="orgId" @saved="emit('saved')" />

    <!-- Card: Technical Competencies (filter multi-select tipe → tabel level) -->
    <TechnicalPotencyCard
      :org-id="orgId"
      @saved="emit('saved')"
      @weight-saved="technicalWeight = $event"
    />

    <!-- Card: Managerial Competencies — bobot mengikuti teknis (100 − teknis) -->
    <ManagerialPotencyCard
      :org-id="orgId"
      :technical-weight="technicalWeight"
      @saved="emit('saved')"
    />

    <!-- Card: Problem Solving & Decision Making — komponen terpisah (2 baris tetap: environment × challenge) -->
    <ProblemSolvingPotencyCard :org-id="orgId" @saved="emit('saved')" />

    <!-- Card: Communication and Influencing Skills — komponen terpisah (tabel level keterampilan) -->
    <SkillPotencyCard :org-id="orgId" @saved="emit('saved')" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from '@/composables/useI18n'
import PsychologicalPotencyCard from './PsychologicalPotencyCard.vue'
import TechnicalPotencyCard from './TechnicalPotencyCard.vue'
import ManagerialPotencyCard from './ManagerialPotencyCard.vue'
import ProblemSolvingPotencyCard from './ProblemSolvingPotencyCard.vue'
import SkillPotencyCard from './SkillPotencyCard.vue'

const emit = defineEmits(['saved'])

defineProps({
  orgId: String,
  jobValueMap: { type: Object, default: () => ({}) },
  competencyOptions: { type: Array, default: () => [] }
})

const { t } = useI18n()

// Bobot kompetensi teknis terbaru (dari TechnicalPotencyCard @weight-saved) —
// diteruskan ke ManagerialPotencyCard agar bobot manajerial (100 − teknis) otomatis.
const technicalWeight = ref(null)

// =========================================================================
// Daftar tetap slug tipe (sama dengan kalkulator backend calculator.go).
// Nama baris dari locale job_values.types.* (bilingual).
// =========================================================================
</script>
