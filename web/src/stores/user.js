// ============================================
// user.js - 用户状态管理（Pinia Store）
// ============================================
// 【知识点】Pinia 是 Vue3 的官方状态管理库
// 类比: Redux(Vue) / Redux(React) / Vuex(Vue2)
//
// 核心概念:
// - state: 数据（类似组件的 data）
// - getters: 计算属性（类似组件的 computed）
// - actions: 方法（类似组件的 methods）
//
// 为什么需要状态管理?
// 组件 A 修改了数据，组件 B 能自动感知变化
// 例如: 登录后用户名显示在顶栏，多个页面都需要访问用户信息

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '../api/auth'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'

// 【知识点】defineStore 定义一个 Store
// 参数1: Store 的唯一 ID（'user'）
// 参数2: Setup 函数（用 Composition API 风格写 Store）
export const useUserStore = defineStore('user', () => {
  // ==========================================
  // State（数据）
  // ==========================================
  // 【知识点】ref() 创建响应式数据
  // 响应式: 数据变化时，依赖它的 UI 自动更新
  const token = ref(localStorage.getItem('token') || '')
  const username = ref('')
  const userInfo = ref(null)

  // ==========================================
  // Getters（计算属性）
  // ==========================================
  // 【知识点】computed() 创建计算属性
  // 它会自动根据依赖变化重新计算
  // 只有 token 有值时才认为已登录
  const isLoggedIn = computed(() => !!token.value)
  // !! 双重取反: 把任意值转成 boolean
  // !! "abc" → true, !! "" → false, !! null → false

  // ==========================================
  // Actions（方法）
  // ==========================================

  // 登录
  async function login(loginForm) {
    const data = await loginApi(loginForm)
    // 【知识点】响应拦截器已经解包了 res.data
    // 所以这里 data 就是 { user, token }

    token.value = data.token
    username.value = data.user.username
    userInfo.value = data.user

    // 持久化 Token 到 localStorage
    localStorage.setItem('token', data.token)
    localStorage.setItem('username', data.user.username)
  }

  // 注册
  async function register(registerForm) {
    await registerApi(registerForm)
    ElMessage.success('注册成功，请登录')
  }

  // 登出
  function logout() {
    token.value = ''
    username.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('username')
    // 【知识点】在 Store 中使用 useRouter 需要在组件外部
    // 这里简单返回，由组件处理跳转
  }

  // 返回所有需要暴露的 state/getters/actions
  // 【知识点】只有 return 的内容才能在组件中使用
  return {
    token,
    username,
    userInfo,
    isLoggedIn,
    login,
    register,
    logout,
  }
})
