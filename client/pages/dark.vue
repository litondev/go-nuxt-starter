<script setup>
import { ref } from 'vue'

const isDark = ref(false)

function toggleDarkMode() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}
</script>

<template>
  <div
    :class="[
      ' border-b w-full backdrop-blur-md transition-colors duration-300',
      isDark ? 'bg-gray-900 bg-opacity-90 border-gray-700 text-white' : 'bg-white bg-opacity-70 border-gray-300 text-gray-900'
    ]"
  >
    <div class="p-4 flex-1">
      <h1 class="font-bold">My Navbar</h1>
    </div>

    <button
      @click="toggleDarkMode"
      class="p-2 mr-4 rounded bg-gray-200 dark:bg-gray-700 text-gray-900 dark:text-white"
      aria-label="Toggle Dark Mode"
    >
      {{ isDark ? 'Light Mode' : 'Dark Mode' }}
    </button>
  </div>

  <div class="pt-16">
    <!-- page content -->
    <p>Page content goes here...</p>
  </div>
</template>