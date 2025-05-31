<script setup lang="ts">
definePageMeta({
  title: 'Profil',
  layout: 'default',
  requiresAuth: true
})

import { Form, Field, ErrorMessage } from 'vee-validate'
import * as yup from 'yup'

const schema = yup.object({
  email: yup.string().email('Invalid email').required('Email is required'),
  password: yup.string().nullable().notRequired().test('min-if-filled', 'Min 8 characters', val => !val || val.length >= 8),
  password_confirm: yup.string().min(8, 'Min 8 characters').required('Password Confirm is required'),
})

const { $axios } = useNuxtApp();
const toast = useToast()

const { data } = useAuth()

const name = ref(data.value.name);
const email = ref(data.value.email);
const password = ref('');
const password_confirm = ref('');
const loading = ref(false);

const onSubmit = () => {
  loading.value = true

  $axios.put("/profil/update",{
    name : name.value,
    email : email.value,
    password : password.value,
    password_confirm : password_confirm.value
  })
  .then(res => {
    toast.add({
      title : 'Success',
      color : 'success'
    })
  })
  .catch(err => {
    toast.add({
      title: 'Terjadi Kesalahan',
      // description: 'There was a problem with your request.',
      color : 'error'
    })
  })
  .finally(() => {
    loading.value = false
  })
}


const config = useRuntimeConfig()
const selectedFile = ref(null)
const previewUrl = ref(config.public.apiUrl.replace("/api/v1","") + '/assets/users/' + data.value.photo)
const isUploading = ref(false)

function handleFileChange(event) {
  const file = event.target.files[0]

  if (file && file.type.startsWith('image/')) {
    selectedFile.value = file
    previewUrl.value = URL.createObjectURL(file)
  } else {
    selectedFile.value = null
    previewUrl.value = null
  }
}

function uploadImage() {
  if (!selectedFile.value) return

  const formData = new FormData()

  formData.append('photo', selectedFile.value)

  isUploading.value = true

  $axios.post('/profil/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  }).then(res => {
    toast.add({
      title : 'Success',
      color : 'success'
    })
  })
  .catch(err => {
    toast.add({
      title: 'Terjadi Kesalahan',
      // description: 'There was a problem with your request.',
      color : 'error'
    })
  })
  .finally(() =>{
    isUploading.value = false
  })
}


</script>

<template>
  <div class="flex transition-colors duration-300">
    <div class="w-full max-w-md p-8 rounded-2xl">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-gray-100">Profil</h2>
      </div>

      <Form @submit="onSubmit" :validation-schema="schema">
        <div class="mb-4">
          <label class="block mb-1 font-medium text-gray-700 dark:text-gray-300">Nama</label>
          <Field
            name="name"
            type="name"
            class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-xl bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 focus:outline-none"
            v-model="name"/>
          <ErrorMessage name="name" class="text-red-600 text-sm mt-1" />
        </div>

        <div class="mb-4">
          <label class="block mb-1 font-medium text-gray-700 dark:text-gray-300">Email</label>
          <Field
            name="email"
            type="email"
            class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-xl bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 focus:outline-none"
            v-model="email"/>
          <ErrorMessage name="email" class="text-red-600 text-sm mt-1" />
        </div>
        
        <div class="mb-4">
          <label class="block mb-1 font-medium text-gray-700 dark:text-gray-300">Password</label>
          <Field
            name="password"
            type="password"
            class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-xl bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 focus:outline-none"
            v-model="password"/>
          <ErrorMessage name="password" class="text-red-600 text-sm mt-1" />
        </div>

         <div class="mb-4">
          <label class="block mb-1 font-medium text-gray-700 dark:text-gray-300">Password Konfirmasi</label>
          <Field
            name="password_confirm"
            type="password_confirm"
            class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-xl bg-white dark:bg-gray-700 text-gray-800 dark:text-gray-100 focus:outline-none"
            v-model="password_confirm"/>
          <ErrorMessage name="password_confirm" class="text-red-600 text-sm mt-1" />
        </div>

        <UButton type="submit" 
          :disabled="loading"
          :loading="loading" 
          class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 rounded-xl transition duration-200 flex justify-center items-center">
          Kirim
        </UButton>
      </Form>
    </div>

    <div class="w-full max-w-md p-8 rounded-2xl mx-2">
      <UInput
        type="file"
        accept="image/*"
        label="Upload Photo"
        @change="handleFileChange"
      />

      <img
        v-if="previewUrl"
        :src="previewUrl"
        alt="Preview"
        class="rounded-lg w-full object-cover h-35 border my-2"
      />

      <UButton
        :loading="isUploading"
        :disabled="!selectedFile || isUploading"
        @click="uploadImage"
        class="my-2"
      >
        Upload Photo
      </UButton>
    </div>
  </div>
</template>