<!--
  UsersView.vue - 用户管理页面（CRUD）
  【知识点】这是一个完整的 CRUD 页面，展示了:
  - 表格展示 + 分页
  - 搜索过滤
  - 新增/编辑对话框
  - 删除确认
-->
<template>
  <div>
    <!-- 搜索栏 + 操作按钮 -->
    <el-card shadow="hover" style="margin-bottom: 20px">
      <el-row :gutter="16" align="middle">
        <el-col :span="8">
          <!-- 【知识点】v-model 双向绑定搜索关键词 -->
          <el-input
            v-model="query.keyword"
            placeholder="搜索用户名/邮箱/昵称"
            clearable
            @keyup.enter="loadUsers"
          />
        </el-col>
        <el-col :span="2">
          <el-button type="primary" @click="loadUsers">搜索</el-button>
        </el-col>
        <el-col :span="14" style="text-align: right">
          <el-button type="primary" @click="openDialog('add')">
            新增用户
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 用户列表表格 -->
    <el-card shadow="hover">
      <el-table :data="users" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="nickname" label="昵称" width="150" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="status" label="状态" width="100">
          <!-- 【知识点】作用域插槽: 自定义列的渲染内容
               #default="scope" 接收当前行的数据
               scope.row 是当前行的对象 -->
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button type="primary" text @click="openDialog('edit', scope.row)">
              编辑
            </el-button>
            <el-popconfirm
              title="确定要删除该用户吗？"
              @confirm="handleDelete(scope.row.id)"
            >
              <!-- 【知识点】el-popconfirm 确认弹出框
                   点击按钮弹出确认框，用户确认后才执行 -->
              <template #reference>
                <el-button type="danger" text>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div style="margin-top: 20px; text-align: right">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="loadUsers"
          @size-change="loadUsers"
        />
        <!-- 【知识点】v-model:current-page 双向绑定当前页码
             用户点击翻页 → query.page 自动更新
             也可以代码修改 query.page → 分页组件自动更新 -->
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <!-- 【知识点】el-dialog 对话框组件
         v-model 控制显示/隐藏 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="500px"
    >
      <el-form :model="formData" label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="formData.username" />
        </el-form-item>
        <el-form-item label="密码" v-if="dialogType === 'add'">
          <el-input v-model="formData.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="formData.nickname" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="formData.email" />
        </el-form-item>
        <el-form-item label="状态">
          <!-- 【知识点】el-switch 开关组件 -->
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getUsers, createUser, updateUser, deleteUser } from '../../api/user'
import { ElMessage } from 'element-plus'

// 表格数据
const users = ref([])
const total = ref(0)
const loading = ref(false)

// 搜索/分页参数
const query = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
})

// 对话框状态
const dialogVisible = ref(false)
const dialogType = ref('add') // 'add' 或 'edit'
const submitLoading = ref(false)

// 表单数据
const formData = reactive({
  id: null,
  username: '',
  password: '',
  nickname: '',
  email: '',
  status: 1,
})

// 加载用户列表
async function loadUsers() {
  loading.value = true
  try {
    const data = await getUsers(query)
    // 【知识点】响应拦截器解包后，data 就是后端返回的 data 字段
    // 即 { total, list, page }
    users.value = data.list
    total.value = data.total
  } catch (err) {
    // 错误由拦截器处理
  } finally {
    loading.value = false
  }
}

// 打开对话框
function openDialog(type, row) {
  dialogType.value = type
  if (type === 'edit' && row) {
    // 编辑: 填充已有数据
    Object.assign(formData, row)
  } else {
    // 新增: 清空表单
    formData.id = null
    formData.username = ''
    formData.password = ''
    formData.nickname = ''
    formData.email = ''
    formData.status = 1
  }
  dialogVisible.value = true
}

// 提交表单（新增或编辑）
async function handleSubmit() {
  submitLoading.value = true
  try {
    if (dialogType.value === 'add') {
      await createUser(formData)
      ElMessage.success('创建成功')
    } else {
      await updateUser(formData.id, formData)
      ElMessage.success('更新成功')
    }
    dialogVisible.value = false
    loadUsers() // 重新加载列表
  } catch (err) {
    // 错误由拦截器处理
  } finally {
    submitLoading.value = false
  }
}

// 删除用户
async function handleDelete(id) {
  try {
    await deleteUser(id)
    ElMessage.success('删除成功')
    loadUsers()
  } catch (err) {
    // 错误由拦截器处理
  }
}

// 格式化日期
// 【知识点】Go 返回的时间是 RFC3339 格式（如 2024-01-15T08:30:00+08:00）
// 这里转成更易读的格式
function formatDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 组件挂载后加载数据
onMounted(() => {
  loadUsers()
})
</script>
