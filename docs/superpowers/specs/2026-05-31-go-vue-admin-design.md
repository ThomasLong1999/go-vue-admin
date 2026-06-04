# Go-Vue-Admin 学习项目设计

## 目标

创建一个渐进式后台管理系统项目，覆盖 Go + Gin + Gorm + Vue3 + Redis + MySQL 六项技术，适合零基础开发者通过读代码 + 学习指南系统学习全栈开发。

## 技术栈

- **后端**: Go 1.21+ / Gin / Gorm / go-redis
- **前端**: Vue3 + Vite + Pinia + Vue Router + Element Plus
- **数据库**: MySQL 8.0 / Redis 7.0

## 项目结构

```
go-vue-admin/
├── server/
│   ├── cmd/
│   │   └── main.go
│   ├── config/
│   │   └── config.yaml
│   ├── internal/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── handler/
│   │   ├── service/
│   │   ├── repository/
│   │   ├── router/
│   │   └── pkg/
│   ├── go.mod
│   └── go.sum
├── web/
│   ├── src/
│   │   ├── views/
│   │   ├── components/
│   │   ├── router/
│   │   ├── stores/
│   │   ├── api/
│   │   └── utils/
│   └── package.json
├── docs/
│   └── learning-guide.md
└── README.md
```

## 分阶段计划

### 阶段 1: Go 基础 + 项目骨架

**学习目标**: Go 项目结构、包管理、配置读取、日志、MySQL 连接初始化

**产出**:
- Go 项目骨架（go.mod、目录结构）
- config.yaml 配置文件及其加载逻辑
- MySQL 连接池初始化
- 基础日志框架
- 健康检查接口 `GET /health`

**核心知识点**:
- Go module 管理依赖
- struct 嵌套与配置映射
- database/sql 基础用法
- log 标准库

### 阶段 2: Gin 路由 + 中间件

**学习目标**: Gin 框架核心用法、RESTful API 设计、JWT 鉴权

**产出**:
- Gin 引擎初始化与路由注册
- 用户注册/登录 API（POST /api/auth/register, POST /api/auth/login）
- JWT 生成与验证中间件
- 请求参数绑定与校验（validator）
- 统一错误处理与响应格式

**核心知识点**:
- Gin 路由组与中间件链
- JWT 原理与实现
- struct tag 校验
- HTTP 状态码规范

### 阶段 3: Gorm + MySQL

**学习目标**: ORM 使用、数据模型设计、CRUD、关联关系、数据库迁移

**产出**:
- 用户、角色、权限数据模型（User/Role/Permission）
- 多对多关联关系（用户-角色、角色-权限）
- Repository 层实现完整 CRUD
- 用户管理 API（GET/POST/PUT/DELETE /api/users）
- Gorm AutoMigrate 数据库迁移

**核心知识点**:
- Gorm 模型定义与 tag
- belongs to / many to many 关联
- Preload / Join 预加载
- 事务用法
- 软删除

### 阶段 4: Redis 缓存

**学习目标**: Redis 数据结构应用、缓存策略、接口限流

**产出**:
- go-redis 连接初始化
- 用户登录 Token 存入 Redis（替代纯 JWT 无状态方案，支持主动踢人）
- 用户信息缓存（读取时先查 Redis，miss 再查 MySQL）
- API 限流中间件（基于 Redis 滑动窗口）
- 缓存工具包封装

**核心知识点**:
- Redis String / Hash 数据结构
- 缓存穿透/雪崩/击穿概念
- 滑动窗口限流算法
- go-redis 基本用法

### 阶段 5: Vue3 前端

**学习目标**: Vue3 组合式 API、路由守卫、状态管理、HTTP 请求封装

**产出**:
- Vite + Vue3 项目初始化
- Element Plus 组件库集成
- 登录页面 + 登录逻辑（Pinia store）
- Vue Router 路由守卫（未登录跳转）
- Axios 请求/响应拦截器（自动附加 Token、错误提示）
- 用户管理 CRUD 页面（表格、表单、搜索、分页）

**核心知识点**:
- Composition API（ref, reactive, computed, watch）
- Pinia store 定义与使用
- Vue Router 导航守卫
- Axios 拦截器模式
- Element Plus 表格/表单/消息组件

### 阶段 6: 前后端联调 + 仪表盘

**学习目标**: 完整联调流程、数据统计、项目收尾

**产出**:
- 完整登录流程联调（前端登录 → 后端JWT → Redis存Token → 前端存储 → 请求携带）
- 仪表盘页面（用户统计、最近登录记录、接口调用统计）
- 接口调用统计（Redis 计数）
- 菜单/侧边栏布局完善
- 项目收尾与总结

**核心知识点**:
- 跨域 CORS 配置
- 前后端协同调试技巧
- 数据聚合查询
- 项目工程化（环境变量、构建配置）

## 后端分层架构

```
HTTP 请求
   ↓
Router (路由注册)
   ↓
Middleware (JWT验证/限流/日志/CORS)
   ↓
Handler (参数绑定、调用Service、返回响应)
   ↓
Service (业务逻辑编排)
   ↓
Repository (数据库/缓存操作)
   ↓
MySQL / Redis
```

## 学习材料

- 代码中每个关键函数/结构体都有中文注释，说明作用和知识点
- `docs/learning-guide.md` 分阶段引导读代码、解释原理、串联知识脉络
