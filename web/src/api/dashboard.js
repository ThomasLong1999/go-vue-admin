// ============================================
// dashboard.js - 仪表盘 API
// ============================================

import request from '../utils/request'

// 获取仪表盘统计数据
export function getDashboardStats() {
  return request.get('/dashboard/stats')
}
