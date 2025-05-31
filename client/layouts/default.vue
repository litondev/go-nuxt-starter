<script setup lang="ts">
import { useColorMode } from '@vueuse/core'

const colorMode = useColorMode()
const { signOut } = useAuth()

const left_items = ref<NavigationMenuItem[]>([
  {
    label: 'User',
    icon: 'i-lucide-users',
    name : 'user',
    children: [
      {
        label: 'Index',
        description: 'Index User',
        icon: 'i-lucide-list',
        to : '/user'
      },
      {
        label: 'Tambah',
        description: 'Tambah User',
        icon: 'i-lucide-user-plus',
        to : '/user/add'
      },
    ]
  }
])

const openMenuLeft = ref<string | null>(null)
const toggleMenuLeft = (name: string) => {
  openMenuLeft.value = openMenuLeft.value === name ? null : name
}

const right_items = ref<NavigationMenuItem[]>([
  {
    label: 'Profil',
    icon: 'i-lucide-user',
    name : 'profil',
    children: [
      {
        label: 'Profil',
        description: 'Profil',
        icon: 'i-lucide-house',
        to: '/profil'
      },
      {
        label: 'Mode',
        description: 'Dark/Light',
        icon: 'i-lucide-sun', 
        click: () => {
          isDark.value = !isDark.value
          localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
          document.documentElement.classList.toggle('dark', isDark.value)
        }
      },
      {
        label: 'Keluar',
        description: 'Keluar',
        icon: 'i-lucide-arrow-right',
        loading : false,
        click: (child) => {
          child.loading = true;

          signOut({callbackUrl: '/login',redirect : true})
          .catch(err => {
            console.log(err)
          })
          .finally(() => {
            child.loading = false;
          })
        }
      }
    ]
  }
])

const openMenuRight = ref<string | null>(null)
const toggleMenuRight = (name: string) => {
  openMenuRight.value = openMenuRight.value === name ? null : name    
}

const menuRef = ref(null);

const onClickOutside = (event) => {
  if (
    menuRef.value && 
    !menuRef.value.contains(event.target)
  ) {
    openMenuRight.value = null;
    openMenuLeft.value = null;
  }
};

onMounted(() => {
  document.addEventListener('click', onClickOutside);
});

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside);
});

const isDark = ref(false)

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'dark') {
    document.documentElement.classList.add('dark')
    isDark.value = true
  }
})
</script>

<template>
  <UApp>
    <div
      class="flex items-center border-b border-default w-full fixed top-0 left-0 right-0 z-50"
      :class="{ dark: colorMode.value === 'dark' }"
      ref="menuRef">
      <div class="p-2 w-auto">
        <img src="/favicon.ico" alt="Logo" class="h-8" />
      </div>

      <nav class="flex justify-between w-full">
        <div class="w-5/6">
          <div
            v-for="item in left_items"
            :key="item.name"
            class="relative"
          >
            <button
              @click="toggleMenuLeft(item.name!)"
              class="flex items-center gap-1 py-2 px-3 hover:text-blue-600 focus:outline-none"
            >
              <Icon :name="item.icon" class="w-5 h-5" />
              <span>{{ item.label }}</span>
            </button>

            <div
              v-if="openMenuLeft === item.name"
              class="absolute bg-white dark:bg-gray-800 border rounded mt-2 p-3 shadow-lg min-w-[200px] z-50"
            >
              <ul class="space-y-2">
                <li
                  v-for="child in item.children"
                  :key="child.label"
                  class="flex items-center gap-2 cursor-pointer hover:text-blue-600"
                >
                  <Icon :name="child.icon" class="w-4 h-4" />
                  <template v-if="child.to">
                    <NuxtLink :to="child.to">{{ child.label }}</NuxtLink>
                  </template>
                  <template v-else-if="child.click">
                    <button @click="child.click">{{ child.label }}</button>
                  </template>
                  <template v-else>
                    <span>{{ child.label }}</span>
                  </template>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <div class="w-auto flex justify-end">
          <div
            v-for="item in right_items"
            :key="item.name"
            class="relative"
          >
            <button
              @click="toggleMenuRight(item.name!)"
              class="flex items-center gap-1 py-2 px-3 hover:text-blue-600 focus:outline-none"
            >
              <Icon :name="item.icon" class="w-5 h-5" />
              <span>{{ item.label }}</span>
            </button>

            <div
              v-if="openMenuRight === item.name"
              class="absolute bg-white dark:bg-gray-800 border rounded mt-2 p-3 shadow-lg min-w-[200px] z-50 right-5">
              <ul class="space-y-2">
                <li
                  v-for="child in item.children"
                  :key="child.label"
                  class="flex items-center gap-2 cursor-pointer hover:text-blue-600"
                >
                  <Icon :name="child.loading ? 'i-lucide-loader-circle' : child.icon" :class="['w-4 h-4',{'animate-spin' : child.loading}]" />

                  <template v-if="child.to">
                    <NuxtLink :to="child.to">{{ child.label }}</NuxtLink>
                  </template>
                  <template v-else-if="child.click">
                    <button type="button" 
                      class="cursor-pointer"
                      :loading="child.loading"
                      @click="() => child.click(child)">
                      {{ child.label }}
                    </button>
                  </template>
                  <template v-else>
                    <span>{{ child.label }}</span>
                  </template>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </nav>
    </div>

    <div class="container mx-auto px-4 py-3 pt-16" :class="{ dark: colorMode.value === 'dark' }">
      <NuxtPage />
    </div>
  </UApp>
</template>
