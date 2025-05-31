import axios from "axios";

export default defineNuxtPlugin((nuxtApp) => {
  const config = useRuntimeConfig();
  const token = useCookie("auth.token")

  axios.defaults.headers.common['X-Requested-With'] = 'XMLHttpRequest';
  axios.defaults.baseURL = config.public.apiUrl;

  axios.interceptors.request.use(function (config) {    
    if(token.value){
      config.headers.Authorization = "Bearer " + token.value;
    }

    return config;
  }, function (error) {  
    return Promise.reject(error);
  });

  nuxtApp.$axios = axios;
})