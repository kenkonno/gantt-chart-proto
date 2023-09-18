import "bootstrap/dist/css/bootstrap.min.css"
import 'vue3-toastify/dist/index.css';
import '@vueform/multiselect/themes/default.css'

import {createApp} from 'vue'
import App from './App.vue'
import router from './router'

import ganttastic from '@infectoone/vue-ganttastic'
import {dateFormat, dateFormatYMD, unixTimeFormat} from "@/utils/filters";

const app = createApp(App)
app.config.globalProperties.$filters = {
    dateFormat: dateFormat,
    dateFormatYMD: dateFormatYMD,
    unixTimeFormat: unixTimeFormat,
}
app.use(router)
    .use(ganttastic)
    .mount('#app')

