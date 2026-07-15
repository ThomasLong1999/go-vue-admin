<!--
  LoginView.vue - 登录页面
  【知识点】这是一个完整的 Vue3 页面组件
  展示了: Composition API、表单处理、Element Plus 组件使用、Pinia Store 调用
-->
<template>
  <!-- 【知识点】Vue 模板中可以用 v-if、v-for 等指令控制渲染 -->
  <div class="login-container">
    <div class="login-card">
      <h2>Go-Vue-Admin 后台管理系统</h2>

      <!-- 【知识点】Element Plus 的表单组件
           :model 绑定表单数据对象
           :rules 绑定校验规则
           ref 用于获取表单实例（调用 validate 方法） -->
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="rules"
        label-width="0"
      >
        <el-form-item prop="username">
          <!-- 【知识点】v-model 双向数据绑定
               输入框的值变化 → loginForm.username 更新
               loginForm.username 更新 → 输入框的值变化 -->
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>

        <el-form-item prop="password">
          <!-- 【知识点】show-password 控制密码可见性切换 -->
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            show-password
            @keyup.enter="handleLogin"
          />
          <!-- 【知识点】@keyup.enter 监听回车键，按下回车就触发登录 -->
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            style="width: 100%"
            @click="handleLogin"
          >
            <!-- 【知识点】:loading 绑定加载状态
                 loading 为 true 时按钮显示加载动画且不可点击 -->
            {{ loading ? '登录中...' : '登 录' }}
          </el-button>
        </el-form-item>
      </el-form>
      <el-button text style="width: 100%" @click="router.push('/register')">
        没有账号？创建学习账号
      </el-button>
    </div>
  </div>
</template>

<!-- 【知识点】<script setup> 中:
     - import 的组件可以直接在模板中使用（不需要注册）
     - 定义的变量和函数自动在模板中可用
     - 不需要 return -->
<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { ElMessage } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

// 表单引用（用于调用 validate 方法）
const formRef = ref(null)

// 加载状态
const loading = ref(false)

// 表单数据
// 【知识点】reactive() 创建响应式对象
// 与 ref 的区别: ref 用于基本类型（string, number），reactive 用于对象
// 但 ref 也能用于对象，ref(obj).value 访问
// 一般规则: 基本类型用 ref，对象用 reactive
const loginForm = reactive({
  username: '',
  password: '',
})

// 表单校验规则
// 【知识点】Element Plus 的表单校验
// required: 必填
// message: 校验失败时的提示
// trigger: 触发时机（blur=失焦, change=值变化）
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
  ],
}

// 登录处理
async function handleLogin() {
  // 1. 先校验表单
  // 【知识点】formRef.value.validate() 返回 Promise
  // 校验通过 → resolve，校验失败 → reject
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  // 2. 发起登录请求
  loading.value = true
  try {
    await userStore.login(loginForm)
    ElMessage.success('登录成功')
    // 【知识点】router.push() 导航到指定路由
    router.push('/dashboard')
  } catch (err) {
    // 错误已由拦截器处理，这里不需要额外提示
  } finally {
    // 【知识点】finally 无论成功或失败都会执行
    // 适合做"关闭 loading"这类清理操作
    loading.value = false
  }
}
</script>

<!-- 【知识点】scoped 样式只作用于当前组件
     Vue 会给组件的 DOM 元素添加唯一属性（data-v-xxx）
     CSS 选择器会自动加上这个属性，实现样式隔离 -->
<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.login-card h2 {
  text-align: center;
  margin-bottom: 30px;
  color: #333;
  font-size: 22px;
}
</style>
