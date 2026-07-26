<template>
  <router-view />
  <Toast position="bottom-right" />
  <ConfirmDialog />
</template>

<script setup>
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
</script>
