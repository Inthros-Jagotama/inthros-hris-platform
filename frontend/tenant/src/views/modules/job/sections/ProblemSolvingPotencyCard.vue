<template>
  <TypesPotencyCard
    :org-id="orgId"
    :types="types"
    :skeleton-rows="2"
    title-key="job_management.problem_solving_title"
    description-key="job_management.problem_solving_description"
    empty-key="job_management.problem_solving_empty"
    save-label-key="job_management.save_problem_solving"
    delete-title-key="job_management.problem_solving_confirm_delete_title"
    delete-message-key="job_management.problem_solving_confirm_delete"
    @saved="emit('saved')"
  />
</template>

<script setup>
import TypesPotencyCard from './TypesPotencyCard.vue'

const emit = defineEmits(['saved'])

defineProps({
  orgId: String
})

// Dua tipe di bawah grup Problem Solving & Decision Making (migration 040):
// thinking_environment (8 level) × thinking_chalenge (5 level)
const types = [
  { type: 'thinking_environment', nameKey: 'job_management.problem_solving_environment' },
  { type: 'thinking_chalenge', nameKey: 'job_management.problem_solving_challenge' }
]
</script>
