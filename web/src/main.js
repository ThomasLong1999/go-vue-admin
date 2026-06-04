// ============================================
// main.js - Vue 应用入口
// ============================================
// 【知识点】这是整个 Vue 前端的启动文件
// 职责: 创建 Vue 实例、注册插件、挂载到 DOM
//
// Vue3 相比 Vue2 的核心变化:
// 1. 用 createApp() 代替 new Vue()
// 2. 用 Composition API 代替 Options API
// 3. 用 Pinia 代替 Vuex（状态管理）

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

import App from './App.vue'
import router from './router'

// 【知识点】createApp(App) 创建一个 Vue 应用实例
// App.vue 是根组件，所有页面组件都是它的子组件
const app = createApp(App)

// 【知识点】use() 注册全局插件
// 插件会给 Vue 添加全局功能
app.use(createPinia())     // Pinia: 状态管理（类似 Redux/Vuex）
app.use(router)            // Vue Router: 路由管理（页面跳转）
app.use(ElementPlus)       // Element Plus: UI 组件库

// 【知识点】mount('#app') 把 Vue 应用挂载到 HTML 中的 <div id="app">
// 从这一刻起，Vue 接管了这个 DOM 元素及其子元素
app.mount('#app')
