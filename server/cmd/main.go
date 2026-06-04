// ============================================
// main.go - 程序入口
// ============================================
// 【知识点】Go 程序的执行入口是 main 包的 main() 函数
// 启动顺序: 导入包的 init() → main 包的 init() → main()
//
// 这是整个项目的"总指挥"，负责:
// - 加载配置
// - 初始化数据库连接
// - 依赖注入: Repository → Service → Handler → Router
// - 启动 HTTP 服务器

package main

import (
	"errors"
	"fmt"
	"log"

	"go-vue-admin/server/config"
	"go-vue-admin/server/internal/handler"
	"go-vue-admin/server/internal/middleware"
	"go-vue-admin/server/internal/model"
	"go-vue-admin/server/internal/pkg"
	"go-vue-admin/server/internal/repository"
	"go-vue-admin/server/internal/router"
	"go-vue-admin/server/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// ==========================================
	// 第 1 步: 加载配置文件
	// ==========================================
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// ==========================================
	// 第 2 步: 设置 Gin 运行模式
	// ==========================================
	// debug: 打印路由和请求日志（开发用）
	// release: 关闭调试信息（生产用）
	gin.SetMode(cfg.Server.Mode)

	// ==========================================
	// 第 3 步: 初始化 MySQL
	// ==========================================
	db, err := pkg.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("初始化 MySQL 失败: %v", err)
	}

	// 【知识点】AutoMigrate 自动创建/更新表结构
	// 表不存在 → 创建；字段有新增 → 添加列
	// 注意: 不会删除列或修改列类型（安全原则）
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// ==========================================
	// 第 4 步: 初始化 Redis（可选）
	// ==========================================
	if _, err := pkg.InitRedis(cfg.Redis); err != nil {
		log.Printf("⚠ Redis 连接失败（缓存功能暂不可用）: %v", err)
	}

	// ==========================================
	// 第 5 步: 依赖注入 —— 逐层组装
	// ==========================================
	// 【知识点】依赖注入: 每层的依赖通过构造函数参数传入
	// 优势: 便于测试（可传入 mock）、便于替换实现
	//
	// 组装顺序（从底到顶）:
	// Repository(需要DB) → Service(需要Repo) → Handler(需要Svc) → Router(需要Handler)

	// Repository 层
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// JWT 管理器
	jwtMgr := pkg.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Service 层
	userSvc := service.NewUserService(userRepo, roleRepo, jwtMgr)
	dashboardSvc := service.NewDashboardService(dashboardRepo)

	// Handler 层
	userHandler := handler.NewUserHandler(userSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

	// ==========================================
	// 第 6 步: 初始化默认数据（角色 + 管理员账号）
	// ==========================================
	initDefaultData(db)

	// ==========================================
	// 第 7 步: 创建 Gin 引擎 + 注册路由
	// ==========================================
	r := gin.New()
	// 【知识点】gin.New() 创建纯净引擎，gin.Default() 自带 Logger+Recovery
	// 用 New() 手动添加中间件，更灵活

	r.Use(gin.Logger())   // 请求日志: 打印每个请求的方法、路径、状态码、耗时
	r.Use(gin.Recovery()) // Panic 恢复: 捕获 panic 防止服务崩溃
	r.Use(middleware.CORS())

	// 注册所有路由
	router.SetupRouter(r, jwtMgr, userHandler, dashboardHandler)

	// ==========================================
	// 第 7 步: 启动 HTTP 服务器
	// ==========================================
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 服务器启动成功: http://localhost:%d", cfg.Server.Port)
	log.Printf("📖 健康检查: http://localhost:%d/api/health", cfg.Server.Port)

	// 【知识点】r.Run() 是阻塞调用，程序会停在这里直到服务关闭
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// initDefaultData 初始化默认数据
// 【知识点】种子数据（Seed Data）: 系统启动时自动创建的初始数据
// 比如默认角色、管理员账号等，确保系统能正常使用
func initDefaultData(db *gorm.DB) {
	// 1. 创建默认角色（如果不存在）
	roles := []model.Role{
		{Name: "管理员", Code: "admin", Description: "系统管理员，拥有所有权限"},
		{Name: "普通用户", Code: "user", Description: "普通用户，基础权限"},
	}
	for _, role := range roles {
		// 【知识点】FirstOrCreate: 查找匹配的记录，不存在则创建
		// 用 Code 作为唯一标识，避免重复创建
		db.Where(model.Role{Code: role.Code}).FirstOrCreate(&role)
	}

	// 2. 创建默认管理员账号（如果不存在）
	var adminUser model.User
	result := db.Where("username = ?", "admin").First(&adminUser)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 管理员不存在，创建
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("⚠ 管理员密码加密失败: %v", err)
			return
		}

		// 获取 admin 角色
		var adminRole model.Role
		db.Where("code = ?", "admin").First(&adminRole)

		adminUser = model.User{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "超级管理员",
			Email:    "admin@example.com",
			Status:   1,
			Roles:    []model.Role{adminRole},
		}
		if err := db.Create(&adminUser).Error; err != nil {
			log.Printf("⚠ 创建管理员账号失败: %v", err)
		} else {
			log.Println("✅ 默认管理员账号创建成功: admin / admin123")
		}
	}
}
