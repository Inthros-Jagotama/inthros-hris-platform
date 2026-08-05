<template>
  <SelectablePotencyCard
    :org-id="orgId"
    type-group="psychological"
    :skeleton-rows="5"
    title-key="job_management.potency_required_title"
    description-key="job_management.potency_required_description"
    empty-key="job_management.potency_required_empty"
    save-label-key="job_management.save_potency_levels"
    delete-title-key="job_management.potency_confirm_delete_title"
    delete-message-key="job_management.potency_confirm_delete"
    @saved="emit('saved')"
  />
</template>

<script setup>
import SelectablePotencyCard from './SelectablePotencyCard.vue'

const emit = defineEmits(['saved'])

defineProps({
  orgId: String
})
</script>
