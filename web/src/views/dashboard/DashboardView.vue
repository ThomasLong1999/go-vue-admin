<!--
  DashboardView.vue - 仪表盘页面
  【知识点】展示统计数据、最近注册用户列表
-->
<template>
  <div>
    <!-- 统计卡片 -->
    <!-- 【知识点】el-row / el-col 是 Element Plus 的栅格布局
         :span="6" 占 1/4 宽度（24 栏系统） -->
    <el-row :gutter="20" style="margin-bottom: 20px">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-number">{{ stats.total_users }}</div>
            <div class="stat-label">总用户数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-number">{{ stats.today_new }}</div>
            <div class="stat-label">今日新增</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-number">{{ stats.active_users }}</div>
            <div class="stat-label">启用用户</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-number">{{ stats.disabled_users }}</div>
            <div class="stat-label">禁用用户</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近注册用户 -->
    <el-card shadow="hover">
      <template #header>
        <span>最近注册用户</span>
      </template>
      <!-- 【知识点】el-table 表格组件
           :data 绑定数据源，自动渲染行 -->
      <el-table :data="stats.recent_users" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="created_at" label="注册时间" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
// 【知识点】onMounted 生命周期钩子
// 组件挂载到 DOM 后执行，适合做初始化数据加载
import { ref, reactive, onMounted } from 'vue'
import { getDashboardStats } from '../../api/dashboard'

// 统计数据（用 reactive 因为是对象）
const stats = reactive({
  total_users: 0,
  today_new: 0,
  active_users: 0,
  disabled_users: 0,
  recent_users: [],
})

// 加载统计数据
async function loadStats() {
  try {
    const data = await getDashboardStats()
    // 【知识点】Object.assign 合并对象
    // 把 data 的属性覆盖到 stats 上
    Object.assign(stats, data)
  } catch (err) {
    // 错误已由拦截器处理
  }
}

// 组件挂载后自动加载数据
onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.stat-card {
  text-align: center;
  padding: 10px 0;
}
.stat-number {
  font-size: 28px;
  font-weight: bold;
  color: #409eff;
}
.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}
</style>
