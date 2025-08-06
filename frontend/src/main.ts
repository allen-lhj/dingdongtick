import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'

import App from './App.vue'
import router from './router'

// 创建应用实例
const app = createApp(App)

// 创建Pinia实例
const pinia = createPinia()

// 安装插件
app.use(pinia)
app.use(router)
app.use(naive)

// 挂载应用
app.mount('#app')
