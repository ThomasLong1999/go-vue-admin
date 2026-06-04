// ============================================
// user.js - 用户管理 API
// ============================================

import request from '../utils/request'

// 获取用户列表
// 【知识点】params 是 URL 查询参数
// axios 会自动把它转成: /users?page=1&page_size=10&keyword=abc
export function getUsers(params) {
  return request.get('/users', { params })
}

// 获取单个用户
// 【知识点】模板字符串: `/users/${id}` 用反引号，${} 插值
export function getUser(id) {
  return request.get(`/users/${id}`)
}

// 创建用户
export function createUser(data) {
  return request.post('/users', data)
}

// 更新用户
export function updateUser(id, data) {
  return request.put(`/users/${id}`, data)
}

// 删除用户
export function deleteUser(id) {
  return request.delete(`/users/${id}`)
}
