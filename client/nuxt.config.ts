// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  
  devServer: {
    port: process.env.NUXT_PORT || 3005
  },

  runtimeConfig: {
    public : {
      apiUrl : process.env.NUXT_API_URL 
    },
  },

  experimental: {
    respectNoSSRHeader: true
  },

  compatibilityDate: '2025-05-15',

  devtools: { 
    enabled: true 
  },

  modules: [
    '@nuxt/image',
    '@nuxt/ui',
    '@nuxt/icon',
    '@pinia/nuxt',
    '@samk-dev/nuxt-vcalendar',
    '@formkit/auto-animate',
    '@sidebase/nuxt-auth',
    '@vee-validate/nuxt',
    "@vueuse/nuxt",
    '@pinia/nuxt',
  ],

  css: ['~/assets/css/main.css'],

  pinia: {
    storesDirs: ['./store/**'],
  },

  auth: {
    // globalAppMiddleware: true,
    isEnabled: true,
    disableServerSideAuth: false,
    originEnvKey: 'AUTH_ORIGIN',
    baseURL: process.env.NUXT_API_URL,
    provider: { 
      type: 'local',
      endpoints: {
        signIn: { path: '/login', method: 'post' },
        signOut: { path: '/logout', method: 'post' },
        signUp: { path: '/register', method: 'post' },
        getSession: { path: '/me', method: 'get' },
      },
      token: {
        signInResponseTokenPointer: '/token',
        type: 'Bearer',
        cookieName: 'auth.token',
        headerName: 'Authorization',
        maxAgeInSeconds: 1800,
        sameSiteAttribute: 'lax',
        cookieDomain: 'localhost',
        secureCookieAttribute: false,
        httpOnlyCookieAttribute: false,
      },
      refresh: {
        isEnabled: true,
        endpoint: { path: '/refresh-token', method: 'post' },
        refreshOnlyToken: true,
        token: {
          signInResponseRefreshTokenPointer: '/token',
          refreshResponseTokenPointer: '',
          refreshRequestTokenPointer: '/token',
          cookieName: 'auth.token',
          maxAgeInSeconds: 3600,
          sameSiteAttribute: 'lax',
          secureCookieAttribute: false,
          cookieDomain: 'localhost',
          httpOnlyCookieAttribute: false,
        }
      }
    },
    sessionRefresh: {
      enablePeriodically: false,
      enableOnWindowFocus: false,
    }
  },

  plugins: [
    {src : "@/plugins/axios.ts",mode : 'client'},
    {src : "@/plugins/init.ts",mode : 'client'},
    {src : "@/plugins/moment.ts",mode : 'client'},
    {src : "@/plugins/floating-vue.ts",mode : 'client'},
    {src : "@/plugins/sweetalert.ts",mode : 'client'},
  ],
})
