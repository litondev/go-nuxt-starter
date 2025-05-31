import { defineNuxtRouteMiddleware, navigateTo } from '#app'

export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuth() 

  if (to.meta.requiresAuth && auth.status.value !== 'authenticated') {
    return navigateTo('/login')
  }

  if (to.path === '/login' && auth.status.value === 'authenticated') {
    return navigateTo('/')
  }
})