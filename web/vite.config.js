import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  // 【知识点】Vite 配置文件使用 ES Module 语法
  // defineConfig 提供类型提示和自动补全
  plugins: [vue()],
  server: {
    port: 5173,
    // 【知识点】Vite 开发服务器代理配置
    // 前端（localhost:5173）和后端（localhost:8080）端口不同
    // 代理: 把 /api 开头的请求转发到后端，解决跨域问题
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true, // 修改请求头中的 Origin
      },
    },
  },
})
