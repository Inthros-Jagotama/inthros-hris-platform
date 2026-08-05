<template>
  <TypesPotencyCard
    :org-id="orgId"
    :types="types"
    :skeleton-rows="2"
    title-key="job_management.skill_communicating_influencing_title"
    description-key="job_management.skill_communicating_influencing_description"
    empty-key="job_management.skill_communicating_influencing_empty"
    save-label-key="job_management.save_skill"
    delete-title-key="job_management.skill_confirm_delete_title"
    delete-message-key="job_management.skill_confirm_delete"
    @saved="emit('saved')"
  />
</template>

<script setup>
import TypesPotencyCard from './TypesPotencyCard.vue'

const emit = defineEmits(['saved'])

defineProps({
  orgId: String
})

// Satu tipe di bawah grup Communication & Influencing Skill (migration 039)
const types = [
  { type: 'communicating_influencing_skill', nameKey: 'job_management.skill_communicating_influencing' }
]
</script>
