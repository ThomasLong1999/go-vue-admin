// ============================================
// request.js - Axios 请求封装
// ============================================
// 【知识点】Axios 是最流行的 HTTP 请求库
// 它比浏览器原生 fetch 更强大:
// - 自动转换 JSON
// - 请求/响应拦截器
// - 超时设置
// - 取消请求
//
// 这里的封装做了:
// 1. 统一 baseURL（所有请求的公共前缀）
// 2. 请求拦截器: 自动附加 JWT Token
// 3. 响应拦截器: 统一处理错误

import axios from 'axios'
import { ElMessage } from 'element-plus'

// 创建 Axios 实例
// 【知识点】axios.create() 创建一个独立的 Axios 实例
// 不影响全局 axios，可以有不同的配置
const request = axios.create({
  baseURL: '/api',  // 所有请求自动加 /api 前缀
  timeout: 10000,   // 10秒超时
})

// ==========================================
// 请求拦截器
// ==========================================
// 【知识点】每个请求发出前都会经过这个拦截器
// 用途: 添加 Token、添加时间戳防缓存等
request.interceptors.request.use(
  (config) => {
    // 自动附加 JWT Token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
      // 【知识点】Bearer Token 格式是 HTTP 标准规范
      // Authorization: Bearer <token>
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// ==========================================
// 响应拦截器
// ==========================================
// 【知识点】每个响应返回后都会经过这个拦截器
// 用途: 统一处理错误码、过期登录等
request.interceptors.response.use(
  (response) => {
    const res = response.data
    // 【知识点】后端统一响应格式: { code, message, data }
    // code !== 0 表示业务错误
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      // 【知识点】Promise.reject 会触发调用方的 .catch()
      return Promise.reject(new Error(res.message))
    }
    return res.data
    // 【知识点】直接返回 data，调用方不需要再 .data.data
    // 这是"解包"模式，减少嵌套
  },
  (error) => {
    // HTTP 层面的错误（网络错误、4xx、5xx）
    if (error.response) {
      const status = error.response.status
      if (status === 401) {
        ElMessage.error('登录已过期，请重新登录')
        localStorage.removeItem('token')
        window.location.href = '/login'
      } else {
        ElMessage.error(`请求错误 (${status})`)
      }
    } else {
      ElMessage.error('网络连接失败')
    }
    return Promise.reject(error)
  }
)

export default request
