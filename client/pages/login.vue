<script setup lang="ts">
definePageMeta({
  title: 'Login',
  layout: 'empty',
  requiresAuth: false
})

import { Form, Field, ErrorMessage } from 'vee-validate'
import * as yup from 'yup'

const schema = yup.object({
  email: yup.string().email('Invalid email').required('Email is required'),
  password: yup.string().min(8, 'Min 8 characters').required('Password is required'),
})

const { signIn } = useAuth()
const toast = useToast()
const route = useRoute();

const email = ref('');
const password = ref('');
const loading = ref(false);

const onSubmit = () => {
  loading.value = true

  signIn({ email: email.value, password: password.value })
  .then(res => {
    route.push('/')
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
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 dark:bg-gray-900 transition-colors duration-300">
    <div class="w-full max-w-md p-8 bg-white dark:bg-gray-800 rounded-2xl shadow dark:shadow-lg">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-800 dark:text-gray-100">Login</h2>
      </div>

      <Form @submit="onSubmit" :validation-schema="schema">
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

        <UButton type="submit" 
          :disabled="loading"
          :loading="loading" 
          class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 rounded-xl transition duration-200 flex justify-center items-center">
          Kirim
        </UButton>
      </Form>
    </div>
  </div>
</template>