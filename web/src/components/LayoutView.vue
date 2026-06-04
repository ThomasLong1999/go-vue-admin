<!--
  LayoutView.vue - 后台布局组件
  【知识点】这是后台管理系统的通用布局: 侧边栏 + 顶栏 + 内容区
  所有需要登录的页面都嵌套在 Layout 内
-->
<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside width="220px" class="aside">
      <div class="logo">
        <h3>Go-Vue-Admin</h3>
      </div>
      <!-- 【知识点】el-menu 是 Element Plus 的导航菜单
           :default-active 当前激活的菜单项
           router=true 点击菜单项自动导航 -->
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/users">
          <span>用户管理</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 右侧主体 -->
    <el-container>
      <!-- 顶栏 -->
      <el-header class="header">
        <div class="header-right">
          <span class="username">{{ userStore.username }}</span>
          <el-button type="danger" text @click="handleLogout">
            退出登录
          </el-button>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="main">
        <!-- 【知识点】<router-view /> 嵌套路由
             Layout 里的 router-view 显示子路由对应的组件
             例如: /dashboard → 显示 DashboardView, /users → 显示 UsersView -->
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 【知识点】computed 自动计算当前激活的菜单项
// route.path 就是当前 URL 路径
const activeMenu = computed(() => route.path)

function handleLogout() {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout-container {
  height: 100vh;
}
.aside {
  background-color: #304156;
  overflow-y: auto;
}
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #263445;
}
.logo h3 {
  color: #fff;
  margin: 0;
  font-size: 16px;
}
.header {
  background: #fff;
  border-bottom: 1px solid #eee;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}
.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.username {
  color: #333;
  font-size: 14px;
}
.main {
  background: #f0f2f5;
  padding: 20px;
}
</style>
