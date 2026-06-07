# Go-Vue-Admin 学习会话记录

> 系统性学习 Go+Gin+Gorm+Vue3+Redis+MySQL，通过后台管理系统项目逐阶段解读。

---

## 项目架构总览

```
go-vue-admin/
├── server/                    ← Go 后端
│   ├── cmd/main.go           ← 程序入口（"总指挥"）
│   ├── config/               ← 配置
│   │   ├── config.yaml       ← YAML 配置文件
│   │   ├── config.go         ← 配置结构体
│   │   └── loader.go         ← 配置加载器
│   └── internal/             ← 内部代码（不对外暴露）
│       ├── pkg/              ← 通用工具包
│       ├── model/            ← 数据模型
│       ├── repository/       ← 数据访问层
│       ├── service/          ← 业务逻辑层
│       ├── handler/          ← 请求处理层
│       ├── middleware/       ← 中间件
│       └── router/           ← 路由
├── web/                       ← Vue3 前端
└── docs/                      ← 文档
```

**数据流方向：**
```
浏览器 → Router → Middleware → Handler → Service → Repository → MySQL/Redis
```

---

## 阶段 1：Go 基础 + 项目骨架

### 1.1 config/config.yaml — 配置文件

YAML 是一种比 JSON 更易读的配置格式，支持注释。

```yaml
server:
  port: 8080
  mode: debug    # debug(开发) / release(生产)

database:
  host: 127.0.0.1
  port: 3306
  user: longtao
  password: "longtao"
  dbname: go_vue_admin
  max_idle_conns: 10    # 连接池闲着等待的最大连接数
  max_open_conns: 100  # 连接池最多允许的连接
  conn_max_lifetime: 3600  # 连接最大存活时间(秒)
```

**连接池的三个参数：**

| 参数 | 含义 | 类比 |
|------|------|------|
| `max_idle_conns: 10` | 空闲等待的最大连接数 | 餐厅里站着待命的服务员 |
| `max_open_conns: 100` | 总共能开的最大连接数 | 餐厅最多雇佣的服务员 |
| `conn_max_lifetime: 3600` | 每个连接最多活1小时 | 服务员轮班，到点下班 |

### 1.2 config/config.go — 配置结构体

**① Go 的可见性规则（极其重要）：**

```go
type DatabaseConfig struct {
    Host   string  // 大写开头 = 公开（其他包可以访问）
    dbname string  // 小写开头 = 私有（只能本包内访问）
}
// Go 没有 public/private 关键字，靠大小写决定！
```

**② struct tag（结构体标签）：**

```go
Port int `mapstructure:"port"`
//       ↑ 这部分就是 tag，写在反引号里
//       告诉 viper 库：YAML 里的 "port" 映射到这个字段
```

**③ 方法 vs 函数：**

```go
// 方法绑定到类型上
func (c *DatabaseConfig) DSN() string {  // ← (c *DatabaseConfig) 是"接收者"
    return c.User + ":" + c.Password + "@tcp(...)..."
}
// 调用方式: cfg.DSN()  而非 DSN(cfg)
```

**DSN 连接字符串拆解：**
```
longtao:longtao@tcp(127.0.0.1:3306)/go_vue_admin?charset=utf8mb4&parseTime=True&loc=Local
│       │        │   │              │  │              │                      │
用户    密码      协议 主机:端口      数据库名         参数1:完整Unicode       参数3:本地时区
                                                       参数2:自动解析时间类型
```

### 1.3 config/loader.go — 配置加载器

**① 多返回值 + Error 处理模式：**

```go
func Load(path string) (*Config, error) {
    if err := v.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %w", err)
    }
    return &cfg, nil
}

// 调用方必须检查 error
cfg, err := Load("config.yaml")
if err != nil {
    log.Fatalf("启动失败: %v", err)
}
```

Go 的哲学是"错误就是普通返回值，别特殊对待"。好处是错误处理逻辑一目了然。

**② `%w` vs `%v` 的错误包装：**

```go
fmt.Errorf("读取配置失败: %w", err)  // %w 包装原始错误，保留链式关系
fmt.Errorf("读取配置失败: %v", err)  // %v 只是把错误信息拼进去，切断关系
```

### 1.4 internal/pkg/database.go — MySQL 连接

**Gorm 的两层架构：**

```
你的代码
   ↓ 调用 Gorm API（gorm.DB）
Gorm（ORM层）     ← 把 Go struct 翻译成 SQL
   ↓ 调用 database/sql（标准库）
database/sql      ← 管理连接池、事务
   ↓ TCP 连接
MySQL 服务器
```

**连接池工作流程：**

```
请求1来了 → 从池里拿空闲连接 → 执行SQL → 放回池里
请求2来了 → 从池里拿空闲连接 → 执行SQL → 放回池里
100个并发请求 → MaxOpenConns=100，刚好每人一个
第101个请求 → 等！！！等别人用完放回池里
```

### 1.5 internal/pkg/redis.go — Redis 连接

**context 是什么：**

```
context 就像快递单上的"备注"：
┌─────────────────────────────┐
│  快递单（context）            │
│  ┌───────────────────────┐  │
│  │ 截止时间: 5 秒          │  │  ← 超时控制
│  │ Trace ID: abc-123      │  │  ← 链路追踪
│  │ 可取消: 是             │  │  ← 取消信号
│  └───────────────────────┘  │
│  包裹: 请求数据              │
└─────────────────────────────┘

context.Background() = 空白快递单，没有截止时间，不可取消
context.WithTimeout() = 设置"如果5秒没送到就别送了"
```

### 1.6 cmd/main.go — 程序入口（依赖注入组装图）

```
              config.Load()
                    │
                    ▼
            pkg.InitDB(cfg)  ────►  *gorm.DB
                    │
      ┌─────────────┼─────────────┐
      ▼             ▼             ▼
UserRepository  RoleRepository  DashboardRepository
      │             │             │
      ▼             ▼             ▼
   UserService  (需要 Repo + JWTManager)
      │
      ▼
   UserHandler  (需要 UserService)
      │
      ▼
router.SetupRouter(r, jwtMgr, userHandler, dashboardHandler)
```

**为什么要依赖注入？**

```go
// ❌ 全局变量写法（耦合严重，难以测试）：
var db *gorm.DB
func GetUser(id uint) *User {
    return db.First(&User{}, id)  // 测试时必须连真实数据库！
}

// ✅ 依赖注入写法（松耦合，可测试）：
type UserService struct {
    userRepo *repository.UserRepository  // ← 通过构造函数传入
}
// 测试时传一个 mock 的 repository 就行了
```

---

## 阶段 2：Gin 路由 + 中间件

### 2.1 请求处理流程

```
浏览器发送:  GET /api/users?keyword=张三
                │
                ▼
┌──────────────────────────────────────────────┐
│                  Gin 引擎                     │
│  ① 全局中间件（Logger → Recovery → CORS）     │
│  ② 路由匹配: /api/users → userHandler.ListUsers │
│  ③ 路由组中间件: JWTAuth() ← 检查 Token       │
│  ④ Handler.ListUsers(c *gin.Context)         │
│  ⑤ 返回 JSON: {"code":0, "data":[...]}       │
└──────────────────────────────────────────────┘
```

### 2.2 两层状态码体系

```
HTTP 状态码（协议层）           业务状态码（应用层）
─────────────────────         ──────────────────
200 OK         请求成功        code:0  业务成功
401 Unauthorized 未认证        code:1001 用户名已存在
404 Not Found  资源不存在       code:1002 密码错误
500 Server Error 服务器崩了     code:1003 Token过期
```

### 2.3 路由对照表（RESTful 设计）

```
公开路由（不需要登录）:
  GET    /api/health              → Health()
  POST   /api/auth/register       → Register()
  POST   /api/auth/login          → Login()

受保护路由（需要 JWT）:
  GET    /api/users               → ListUsers()
  POST   /api/users               → CreateUser()
  GET    /api/users/:id           → GetUser()
  PUT    /api/users/:id           → UpdateUser()
  DELETE /api/users/:id           → DeleteUser()
  GET    /api/dashboard/stats     → DashboardStats()
```

RESTful 的核心思想：用 HTTP 方法表达操作意图，URL 只表示资源。

### 2.4 中间件三个核心机制

**① 闭包：**

```go
func JWTAuth(jwtMgr *pkg.JWTManager) gin.HandlerFunc {
    return func(c *gin.Context) {  // 内层函数"捕获"外层的 jwtMgr
        claims, err := jwtMgr.ParseToken(tokenString)
        ...
    }
}
```

**② c.Abort() vs c.Next()：**

```
请求进来 → 中间件前置代码
    ├── c.Abort()  → 🛑 直接返回，后续不执行
    └── c.Next()   → ✅ 继续执行下一个中间件/Handler
```

**③ c.Set() / c.Get()：**

```go
// 中间件里存：
c.Set("user_id", claims.UserID)

// Handler 里取：
userID, _ := c.Get("user_id")
```

### 2.5 CORS 跨域

```
前端: http://localhost:5173  (Vite)
后端: http://localhost:8080  (Go)
                        ↑
                    端口不同 → 跨域！

CORS 中间件在响应头加 "Access-Control-Allow-Origin: *"
告诉浏览器：这个接口允许任何域名的前端调用
```

OPTIONS 预检请求：浏览器在发跨域请求前，先发 OPTIONS "试探"。

### 2.6 Handler 的三种参数绑定

```go
// 方式1: ShouldBindJSON ← POST 请求体 JSON
var input service.LoginInput
c.ShouldBindJSON(&input)

// 方式2: ShouldBindQuery ← URL 查询参数
var query model.PageQuery
c.ShouldBindQuery(&query)

// 方式3: c.Param() ← 路径参数
id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
```

### 2.7 JWT 原理

```
登录流程:
  1. 用户输密码登录
  2. 服务端用密钥签名 → JWT: "eyJhbG...eyJ1c2Vy...签名"
  3. 前端存到 localStorage

后续请求:
  Authorization: Bearer eyJhbG...eyJ1c2Vy...签名
  → 拆开 Token，验签名 → 读出 user_id → 放行

JWT 结构:
  eyJhbGciOiJIUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.abc123def456
  └── Header ────┘ └── Payload ──┘ └─ Signature ──┘

Session vs JWT:
  Session: 服务端存状态 → 分布式麻烦
  JWT:     Token 自带信息 → 任何服务器都能验证 → 简单
```

---

## 阶段 3：Gorm + MySQL 数据层

### 3.1 多对多关联

```go
Roles []Role `gorm:"many2many:user_roles;"`
```

Gorm 自动创建中间表 `user_roles`：

```
┌──────────┬─────────┐
│ user_id  │ role_id │
├──────────┼─────────┤
│    1     │    1    │  ← 用户1是管理员
│    1     │    2    │  ← 用户1也是普通用户
│    2     │    2    │  ← 用户2是普通用户
└──────────┴─────────┘
```

RBAC 数据关系：
```
users ──many2many──> roles ──many2many──> permissions
 张三是管理员           管理员有"删除用户"权限
 李四是编辑             编辑有"编辑文章"权限
```

### 3.2 软删除

```sql
-- Gorm 执行: db.Delete(&User{}, 1)
-- 不是 DELETE，而是 UPDATE:
UPDATE users SET deleted_at = '2024-01-15 10:30:00' WHERE id = 1;

-- 之后查询自动过滤:
SELECT * FROM users WHERE deleted_at IS NULL;
```

好处：数据可恢复、审计追溯、关联数据安全。

### 3.3 分页查询

```
前端: GET /api/users?page=3&page_size=10
Go:   offset = (3-1) * 10 = 20, limit = 10

SQL:  SELECT * FROM users LIMIT 10 OFFSET 20;
      SELECT COUNT(*) FROM users;  -- 算总页数
```

### 3.4 Go struct 嵌入（embedding）

```go
type Claims struct {
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims  // 嵌入，自动拥有标准字段
}
// claims.ExpiresAt 来自 RegisteredClaims
// claims.UserID 自己的字段
```

### 3.5 类型断言

```go
if claims, ok := token.Claims.(*Claims); ok && token.Valid {
    return claims, nil
}
// "我赌你是 *Claims 类型"
// ok=true → 赌对了
// ok=false → 赌错了
```

---

## 阶段 5：Vue3 前端

### 5.1 Vite 代理

```
浏览器 ──→ localhost:5173 (Vite) ──proxy──→ localhost:8080 (Go)
浏览器看来: 全在 5173 端口，没有跨域

这跟 CORS 中间件是两种互补方案:
  - Vite proxy → 开发环境用
  - CORS 中间件 → 生产环境用
```

### 5.2 响应拦截器的"解包"模式

```
后端返回: { "code": 0, "message": "ok", "data": { "total": 95, "list": [...] } }

没有拦截器: res.data.data.list（嵌套）
有了拦截器: data.list（直接拿到）
```

### 5.3 路由守卫

```
访问 /dashboard → beforeEach
  → 有 token? → YES → 放行
  → 无 token? → NO → 跳转 /login

嵌套路由:
  /dashboard → LayoutView 渲染 + DashboardView 出现在 <router-view />
  /users     → LayoutView 渲染 + UsersView 出现在 <router-view />
  侧边栏和顶栏保持不变！
```

### 5.4 ref vs reactive

```js
// ref: 包装基本类型，访问需要 .value
const count = ref(0)
count.value++
// 模板里自动解包：{{ count }}

// reactive: 包装对象，直接访问
const form = reactive({ name: '' })
form.name = '张三'
```

### 5.5 onMounted

```
Vue 组件生命周期: 创建 → 挂载到DOM → 更新 → 销毁
                         ↑
                    onMounted 在这里
onMounted(() => { loadUsers() })
```

---

## 阶段 6：完整数据流

### 一次"查看用户列表"的完整旅程

```
1. 用户点击"用户管理"
   Vue Router → 匹配 /users → 渲染 UsersView.vue → onMounted → loadUsers()

2. 前端发请求
   API层: request.get('/users', { params: {page:1, size:10} })
   请求拦截器: 从 localStorage 取 token, 加 Authorization 头
   Vite代理: localhost:5173/api/users → localhost:8080/api/users

3. Go 后端接收
   全局中间件: Logger → Recovery → CORS
   路由匹配: GET /api/users → JWTAuth()
   JWT中间件: 提取 Bearer Token → ParseToken → 验证
   Handler: ShouldBindQuery(&query) → userSvc.ListUsers(query)

4. Service 层
   ListUsers → userRepo.List(query)

5. Repository 层
   db.Where(...).Offset(0).Limit(10).Order("id DESC").Find(&users)
   SQL: SELECT * FROM users WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 10;

6. 结果层层返回
   Repository → Service → Handler → pkg.Success(c, PageResult{...})
   Gin 序列化 JSON → {"code":0,"message":"ok","data":{...}}

7. 前端接收
   响应拦截器: code===0 → 解包返回 data
   users.value = data.list → Vue 响应式自动更新表格
```

### 安全性三层防护

```
第1层——密码不返回:
  User struct: Password string `json:"-"` → 序列化时跳过

第2层——JWT 鉴权:
  不带 Token / Token 伪造 → 401 拦截

第3层——SQL 注入防护:
  db.Where("username = ?", input)  // Gorm 自动转义
  ❌ 永远不要: db.Raw("..." + input)
```

---

## 文件清单

| 层级 | 文件 | 职责 |
|------|------|------|
| 配置 | `config/config.yaml` | 数据库/Redis/JWT 参数 |
| 配置 | `config/config.go` | 配置结构体 + DSN 拼装 |
| 配置 | `config/loader.go` | Viper 读取 YAML |
| 入口 | `cmd/main.go` | 启动 + 依赖注入 + 种子数据 |
| 工具 | `internal/pkg/database.go` | Gorm + MySQL 连接池 |
| 工具 | `internal/pkg/redis.go` | go-redis 客户端 |
| 工具 | `internal/pkg/response.go` | 统一 `{code, message, data}` |
| 工具 | `internal/pkg/jwt.go` | JWT 生成/解析 |
| 模型 | `internal/model/model.go` | User/Role/Permission + 分页 |
| 数据 | `internal/repository/user_repository.go` | 用户 CRUD |
| 数据 | `internal/repository/role_repository.go` | 角色查询 |
| 数据 | `internal/repository/dashboard_repository.go` | 统计查询 |
| 业务 | `internal/service/user_service.go` | 注册登录 + 密码加密 |
| 业务 | `internal/service/dashboard_service.go` | 仪表盘数据编排 |
| 路由 | `internal/router/router.go` | URL → Handler 映射 |
| 中间件 | `internal/middleware/jwt.go` | JWT 鉴权 + CORS |
| 处理 | `internal/handler/handler.go` | HTTP 参数绑定 + 响应 |
| 前端入口 | `web/src/main.js` | Vue 应用创建 + 插件注册 |
| 前端入口 | `web/src/App.vue` | 根组件 `<router-view />` |
| 路由 | `web/src/router/index.js` | 页面路由 + 守卫 |
| 状态 | `web/src/stores/user.js` | Pinia 登录态管理 |
| 网络 | `web/src/utils/request.js` | Axios 封装 + 拦截器 |
| API | `web/src/api/auth.js` | 登录/注册接口 |
| API | `web/src/api/user.js` | 用户 CRUD 接口 |
| API | `web/src/api/dashboard.js` | 仪表盘接口 |
| 页面 | `web/src/views/login/LoginView.vue` | 登录页 |
| 页面 | `web/src/views/dashboard/DashboardView.vue` | 仪表盘 |
| 页面 | `web/src/views/users/UsersView.vue` | 用户管理 CRUD |
| 布局 | `web/src/components/LayoutView.vue` | 侧边栏+顶栏布局 |
| 构建 | `web/vite.config.js` | Vite 代理配置 |

---

## 默认登录账号

| 项目 | 值 |
|------|------|
| 用户名 | `admin` |
| 密码 | `admin123` |
