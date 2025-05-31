
<script setup lang="ts">
definePageMeta({
  title: 'Test',
  layout: 'default'
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

    <button @click="signOut({callbackUrl: '/login',redirect : true})">
      Logout
    </button>
  </div>
</template>