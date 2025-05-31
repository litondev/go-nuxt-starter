<script setup lang="ts">
import { useColorMode } from '@vueuse/core'

const colorMode = useColorMode()
</script>

<template>
  <UApp>
    <div :class="{ dark: colorMode === 'dark' }">
      <NuxtPage />
    </div>
  </UApp>
</template>