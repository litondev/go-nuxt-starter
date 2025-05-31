
<script setup lang="ts">
definePageMeta({
  title: 'Index',
  layout: 'default',
  auth: {
    unauthenticatedOnly: false,
    navigateUnauthenticatedTo: '/login'
  }
})
</script>

<template>
  <div>
    <p>Index</p>
  </div>
</template>