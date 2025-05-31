<script setup lang="ts">
definePageMeta({
  title: 'User',
  layout: 'default',
  requiresAuth: true
})

const columns = [
  { id: 'id',lable: 'Id'},
  { id: 'name', label: 'Name' },
  { id: 'email', label: 'Email' }
]

const isLoadingPage = ref(true);
const users = ref("");

const { $axios } = useNuxtApp();

onMounted(() => {
    onLoad()
})

const onLoad = () => { 
    $axios.get("/user")
    .then(res => {
        users.value = res.data;
    })
    .catch(err => {
        console.log(err);
    })
    .finally(() => {
        isLoadingPage.value = false;
    })
}
</script>

<template>
  <div class="max-w-xl mx-auto p-4"
    v-if="!isLoadingPage">
    {{  users.data }}
    <UTable :rows="users.data" :columns="columns" />
  </div>

  <div class="max-w-xl mx-auto p-4"
    v-if="isLoadingPage">
  </div>
</template>
