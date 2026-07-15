# Go-Vue-Admin

一个适合零基础学习的全栈后台管理系统项目。

## 技术栈

- **后端**: Go + Gin + Gorm + MySQL + Redis
- **前端**: Vue3 + Vite + Pinia + Element Plus

## 快速开始

```bash
# 1. 创建 MySQL 数据库
mysql -u root -p -e "CREATE DATABASE go_vue_admin CHARACTER SET utf8mb4;"

# 2. 复制示例配置并修改本机数据库密码
cp server/config/config.example.yaml server/config/config.yaml

# 3. 启动后端
cd server && go mod tidy && go run cmd/main.go

# 4. 启动前端（新终端）
cd web && npm install && npm run dev
```

访问 http://localhost:5173

也可以用环境变量覆盖 YAML 中的值，例如在 PowerShell 中执行：

```powershell
$env:JWT_SECRET = "至少 32 个字符的本机开发密钥"
$env:DATABASE_PASSWORD = "你的本机数据库密码"
```

## 建议的学习验证顺序

1. 从 `config.example.yaml` 复制本地配置，并尝试用 `JWT_SECRET` 覆盖其中的值。
2. 从注册页创建普通用户，观察其可以查看仪表盘但不能访问用户管理接口。
3. 使用 debug 配置中的管理员登录，再练习用户的创建、编辑和删除。
4. 连续访问两次仪表盘，并在 Redis 中查看 `dashboard:stats` 缓存键。
5. 修改 Go 代码后执行 `go test ./...`，确认配置、JWT、权限和分页行为没有回归。

## 学习文档

详见 [学习指南](docs/learning-guide.md)

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | /api/health | 健康检查 | 否 |
| POST | /api/auth/register | 用户注册 | 否 |
| POST | /api/auth/login | 用户登录 | 否 |
| GET | /api/users | 用户列表 | 管理员 |
| POST | /api/users | 创建用户 | 管理员 |
| GET | /api/users/:id | 用户详情 | 管理员 |
| PUT | /api/users/:id | 更新用户 | 管理员 |
| DELETE | /api/users/:id | 删除用户 | 管理员 |
| GET | /api/dashboard/stats | 仪表盘统计 | 登录用户 |
