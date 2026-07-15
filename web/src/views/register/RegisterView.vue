<template>
  <div class="register-container">
    <div class="register-card">
      <h2>创建学习账号</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名（3-32 位）" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            show-password
            placeholder="密码（至少 6 位）"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱（可选）" size="large" />
        </el-form-item>
        <el-form-item prop="nickname">
          <el-input v-model="form.nickname" placeholder="昵称（可选）" size="large" />
        </el-form-item>
        <el-button type="primary" :loading="loading" style="width: 100%" size="large" @click="submit">
          {{ loading ? '注册中...' : '注 册' }}
        </el-button>
      </el-form>
      <el-button text style="width: 100%; margin-top: 12px" @click="router.push('/login')">
        已有账号，去登录
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)
const form = reactive({
  username: '',
  password: '',
  email: '',
  nickname: '',
})

const rules = {
  username: [{ required: true, min: 3, max: 32, message: '用户名长度为 3-32 位', trigger: 'blur' }],
  password: [{ required: true, min: 6, max: 32, message: '密码长度为 6-32 位', trigger: 'blur' }],
  email: [{ type: 'email', message: '请输入正确邮箱', trigger: 'blur' }],
}

async function submit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await userStore.register(form)
    router.push('/login')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.register-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.register-card h2 {
  margin: 0 0 30px;
  color: #333;
  font-size: 22px;
  text-align: center;
}
</style>
