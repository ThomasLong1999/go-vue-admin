# Go-Vue-Admin

一个适合零基础学习的全栈后台管理系统项目。

## 技术栈

- **后端**: Go + Gin + Gorm + MySQL + Redis
- **前端**: Vue3 + Vite + Pinia + Element Plus

## 快速开始

```bash
# 1. 创建 MySQL 数据库
mysql -u root -p -e "CREATE DATABASE go_vue_admin CHARACTER SET utf8mb4;"

# 2. 修改配置 server/config/config.yaml 中的数据库密码

# 3. 启动后端
cd server && go mod tidy && go run cmd/main.go

# 4. 启动前端（新终端）
cd web && npm install && npm run dev
```

访问 http://localhost:5173

## 学习文档

详见 [学习指南](docs/learning-guide.md)

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/health | 健康检查 | 否 |
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/users | 用户列表 | 是 |
| POST | /api/users | 创建用户 | 是 |
| GET | /api/users/:id | 用户详情 | 是 |
| PUT | /api/users/:id | 更新用户 | 是 |
| DELETE | /api/users/:id | 删除用户 | 是 |
| GET | /api/dashboard/stats | 仪表盘统计 | 是 |
