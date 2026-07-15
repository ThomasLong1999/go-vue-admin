// ============================================
// index.js - Vue Router 路由配置
// ============================================
// 【知识点】Vue Router 实现单页应用(SPA)的页面切换
// 嵌套路由: Layout 作为父路由，Dashboard/Users 作为子路由
// 这样 Layout（侧边栏+顶栏）在所有子页面中共享

import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/register/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/login',
    name: 'Login',
    // 【知识点】懒加载: () => import() 动态导入，减少首屏体积
    component: () => import('../views/login/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    path: '/',
    component: () => import('../components/LayoutView.vue'),
    // 【知识点】嵌套路由: 子路由的组件显示在父组件的 <router-view> 中
    // LayoutView.vue 里的 <router-view /> 显示这里的子组件
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/dashboard/DashboardView.vue'),
        meta: { requiresAuth: true, title: '仪表盘' },
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/users/UsersView.vue'),
        meta: { requiresAuth: true, title: '用户管理' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// 全局前置守卫（路由拦截器）
router.beforeEach((to, from, next) => {
  document.title = to.meta.title ? `${to.meta.title} - Go-Vue-Admin` : 'Go-Vue-Admin'

  if (to.meta.requiresAuth !== false) {
    const token = localStorage.getItem('token')
    if (!token) {
      next('/login')
      return
    }
  }
  next()
})

export default router
