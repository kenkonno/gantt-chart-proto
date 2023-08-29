import "bootstrap/dist/css/bootstrap.min.css"
import 'vue3-toastify/dist/index.css';
import '@vueform/multiselect/themes/default.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'

import ganttastic from '@infectoone/vue-ganttastic'


createApp(App)
    .use(router)
    .use(ganttastic)
    .mount('#app')
