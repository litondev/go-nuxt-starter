
<script setup lang="ts">
definePageMeta({
  title: 'Index',
  layout: 'default',
  auth : false,
  // auth: {
  //   unauthenticatedOnly: false,
  //   navigateUnauthenticatedTo: '/auth/signin'
  // }
})

// STATE
const store = useCounterStore()
console.log(store.count);

// OTHER
const { $axios } = useNuxtApp();
const { $moment } = useNuxtApp();

// AUTH
const { signIn,signOut } = useAuth()
const { data } = useAuth()
console.log(data?.value?.id);
</script>

<template>
  <div>
    <button @click="signIn({ email: 'user1@gmail.com', password: 'password' },{ callbackUrl: '/', redirect: true })"
      v-tooltip="'Halo'">
      Login
    </button>

    <UButton icon="i-lucide-search" size="xs" color="secondary" variant="outline" @click="">
      Cari
    </UButton>

    <button @click="signOut({callbackUrl: '/login',redirect : true})">
      Logout
    </button>
  </div>
</template>