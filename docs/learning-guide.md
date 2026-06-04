# Go-Vue-Admin 学习指南

> 一个从零开始的全栈后台管理系统项目，覆盖 Go + Gin + Gorm + Vue3 + Redis + MySQL 六项技术。

## 如何使用这个项目

### 前置准备

1. **安装 Go** (1.21+): https://go.dev/dl/
2. **安装 Node.js** (18+): https://nodejs.org/
3. **安装 MySQL** (8.0): https://dev.mysql.com/downloads/
4. **安装 Redis** (7.0): https://redis.io/download/
5. **创建数据库**: 在 MySQL 中执行 `CREATE DATABASE go_vue_admin CHARACTER SET utf8mb4;`

### 启动后端

```bash
cd server
go mod tidy          # 下载依赖
go run cmd/main.go   # 启动服务，监听 8080 端口
```

### 启动前端

```bash
cd web
npm install          # 安装依赖
npm run dev          # 启动开发服务器，监听 5173 端口
```

打开浏览器访问 http://localhost:5173

---

## 学习路线（6个阶段）

### 阶段 1: Go 基础 + 项目骨架

**目标**: 理解 Go 项目结构、配置管理、数据库连接

**阅读顺序**:

1. **`config/config.yaml`** — YAML 配置文件格式，了解每项配置的含义
2. **`config/config.go`** — Go 的 struct 定义，学习:
   - struct 和 struct tag（`mapstructure` 标签的作用）
   - 方法定义（`func (c *DatabaseConfig) DSN()`）
   - 大写/小写的可见性规则
3. **`config/loader.go`** — Viper 配置加载，学习:
   - Go 的 `import` 机制
   - 多返回值和 error 处理模式
   - `fmt.Errorf` 和 `%w` 错误包装
4. **`internal/pkg/database.go`** — MySQL 连接初始化，学习:
   - Gorm 的 `Open` 和连接池配置
   - `SetMaxIdleConns` / `SetMaxOpenConns` / `SetConnMaxLifetime` 的含义
5. **`cmd/main.go`** — 程序入口，学习:
   - Go 程序的启动顺序
   - 依赖注入的组装过程

**核心知识点**: struct、方法、指针、error 处理、import、包管理

---

### 阶段 2: Gin 路由 + 中间件

**目标**: 理解 HTTP 服务框架的核心机制

**阅读顺序**:

1. **`internal/router/router.go`** — 路由注册，学习:
   - RESTful API 设计（GET/POST/PUT/DELETE 的语义）
   - 路由组（Route Group）的作用
   - 公开路由 vs 受保护路由的分离
2. **`internal/middleware/jwt.go`** — 中间件，学习:
   - Gin 中间件的洋葱模型
   - 闭包的概念和用法
   - `c.Abort()` 和 `c.Next()` 的区别
   - JWT Token 从 Authorization 头提取和验证
   - CORS 跨域的原理和解决
3. **`internal/handler/handler.go`** — 请求处理层，学习:
   - `gin.Context` 的核心方法
   - `ShouldBindJSON` 参数绑定和校验
   - `c.Param()` 路径参数 vs `ShouldBindQuery()` 查询参数
4. **`internal/pkg/response.go`** — 统一响应格式，学习:
   - 为什么需要统一的响应格式
   - `interface{}` 空接口类型

**核心知识点**: Gin 路由、中间件、参数绑定、JWT、CORS、RESTful

---

### 阶段 3: Gorm + MySQL

**目标**: 理解 ORM 模型和数据库操作

**阅读顺序**:

1. **`internal/model/model.go`** — 数据模型定义，学习:
   - Gorm 模型和 struct tag（`gorm:"primaryKey"`, `gorm:"many2many"` 等）
   - 关联关系: 多对多（User ↔ Role ↔ Permission）
   - 软删除（`gorm.DeletedAt`）的原理
   - 分页查询结构体设计
2. **`internal/repository/user_repository.go`** — 数据访问层，学习:
   - Repository 模式的意义（为什么不在 Service 里直接写 SQL）
   - Gorm CRUD 操作（Create / First / Find / Save / Delete）
   - Preload 预加载关联数据
   - LIKE 模糊查询和分页（Offset + Limit）
3. **`internal/repository/role_repository.go`** — 简单参考
4. **`cmd/main.go`** 中的 `AutoMigrate` — 自动数据库迁移

**核心知识点**: Gorm 模型、关联关系、CRUD、Preload、软删除、分页

---

### 阶段 4: Redis 缓存

**目标**: 理解 Redis 在 Web 应用中的典型用法

**阅读顺序**:

1. **`internal/pkg/redis.go`** — Redis 连接初始化，学习:
   - go-redis 客户端的基本配置
   - Ping 测试连接
   - `context.Background()` 的含义
2. **`internal/pkg/jwt.go`** — JWT 工具包，学习:
   - JWT 的三部分结构（Header / Payload / Signature）
   - Claims 自定义载荷和组合（embedding）
   - Token 生成和验证的流程

**核心知识点**: Redis 基本用法、JWT 原理、context

---

### 阶段 5: Vue3 前端

**目标**: 理解 Vue3 组合式 API 和前端工程化

**阅读顺序**:

1. **`web/vite.config.js`** — Vite 配置，学习:
   - 代理配置解决跨域
2. **`web/src/utils/request.js`** — Axios 封装，学习:
   - Axios 实例创建和配置
   - 请求拦截器（自动附加 Token）
   - 响应拦截器（统一错误处理、解包数据）
3. **`web/src/stores/user.js`** — Pinia 状态管理，学习:
   - `ref()` / `reactive()` / `computed()` 响应式系统
   - Pinia Store 的定义（state / getters / actions）
   - Composition API 写法
4. **`web/src/router/index.js`** — 路由配置，学习:
   - 路由懒加载（`() => import()`）
   - 嵌套路由（Layout 包裹子页面）
   - 路由守卫（beforeEach）做权限控制
5. **`web/src/views/login/LoginView.vue`** — 登录页面，学习:
   - Vue 单文件组件（SFC）结构
   - `<script setup>` 语法糖
   - `v-model` 双向绑定
   - Element Plus 表单和校验
6. **`web/src/views/users/UsersView.vue`** — 用户管理页面，学习:
   - 表格 + 分页 + 搜索
   - 对话框新增/编辑
   - 作用域插槽自定义列内容
   - `onMounted` 生命周期

**核心知识点**: Composition API、响应式、Pinia、路由、Axios 拦截器、Element Plus

---

### 阶段 6: 前后端联调

**目标**: 理解完整的前后端数据流

**跟着做**:

1. 先启动 MySQL 和 Redis
2. 启动后端 `go run cmd/main.go`，确认 `GET /api/health` 返回成功
3. 启动前端 `npm run dev`
4. 打开 http://localhost:5173
5. 尝试完整流程:
   - 注册账号 → 登录 → 查看仪表盘 → 管理用户
6. 打开浏览器 F12 开发者工具:
   - Network 面板: 观察每个请求的请求/响应
   - Application 面板: 查看 localStorage 中的 token
7. 故意制造错误: 不带 token 直接访问 /dashboard，观察路由守卫的拦截

**思考**:

- 请求从浏览器发出到数据库，经过了哪些层？
- 如果去掉 Redis，哪些功能会受影响？
- Token 过期后，前端应该怎么处理？

---

## 架构总览

```
浏览器 (Vue3)
    ↓ Axios + JWT Token
Gin Router
    ↓ JWT中间件 → CORS中间件
Handler (参数绑定)
    ↓
Service (业务逻辑)
    ↓
Repository (Gorm操作)
    ↓
MySQL / Redis
```

## 每个文件的角色

| 层级 | 目录 | 职责 |
|------|------|------|
| 入口 | `cmd/main.go` | 启动、组装依赖、启动服务 |
| 路由 | `internal/router/` | URL → Handler 映射 |
| 中间件 | `internal/middleware/` | JWT鉴权、CORS、限流 |
| 处理器 | `internal/handler/` | HTTP 请求/响应处理 |
| 业务 | `internal/service/` | 业务逻辑编排 |
| 数据 | `internal/repository/` | 数据库/缓存操作 |
| 模型 | `internal/model/` | 数据结构定义 |
| 工具 | `internal/pkg/` | JWT、响应封装、DB/Redis初始化 |
| 配置 | `config/` | 配置文件和加载 |
