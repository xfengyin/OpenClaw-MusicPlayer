import { createApp } from 'vue'
import { createPinia } from 'pinia'
import naive from 'naive-ui'
import App from './App.vue'

import './assets/main.css'

/**
 * 创建应用实例
 */
const app = createApp(App)

// Pinia 状态管理
const pinia = createPinia()

// 使用插件
app.use(pinia)
app.use(naive)
app.mount('#app')