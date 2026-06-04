// ============================================
// auth.js - 认证相关 API
// ============================================
// 【知识点】把 API 请求按业务模块组织
// 每个文件对应后端的一个 Controller/Handler

import request from '../utils/request'

// 用户注册
export function register(data) {
  return request.post('/auth/register', data)
}

// 用户登录
// 【知识点】返回值是一个 Promise
// 调用方式: const res = await login({ username, password })
// 或者: login({ username, password }).then(res => ...)
export function login(data) {
  return request.post('/auth/login', data)
}
