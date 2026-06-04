// ============================================
// database.go - MySQL 数据库初始化
// ============================================
// 【知识点】Go 的 init() 函数
// 每个包可以有一个 init() 函数，它在 main() 之前自动执行
// 但这里我们不用 init()，而是用显式初始化函数，因为我们需要传入配置参数
// init() 的缺点是它无法接收参数，也不返回错误，难以测试和控制

package pkg

import (
	"fmt"
	"log"
	"time"

	"go-vue-admin/server/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 是全局数据库连接对象
// 【知识点】var 声明全局变量，这个变量可以在整个包中被访问
// 但在其他包中要用 GetDB() 来获取，因为 DB 是小写的（包私有）
var DB *gorm.DB

// InitDB 初始化 MySQL 数据库连接
// 【知识点】Go 的错误处理哲学：没有 try-catch，通过返回值检查
// 这是 Go 的核心设计：显式处理错误，不忽略任何一个可能失败的操作
//
// 参数:
//   - cfg: 从配置文件加载的数据库配置
//
// 返回:
//   - *gorm.DB: Gorm 数据库连接实例
//   - error: 如果连接失败返回错误
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()

	// 【知识点】Gorm 的 Open 函数连接数据库
	// 第一个参数是数据库驱动（driver.Dialector），决定连接哪种数据库
	// 第二个参数是 Gorm 配置选项，用链式调用设置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger 配置：开发模式打印所有 SQL，方便调试
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	// 【知识点】获取底层 sql.DB 对象来配置连接池
	// Gorm 是 ORM 层，底层使用 database/sql 包的连接池
	// 连接池的作用：复用数据库连接，避免每次请求都创建新连接（太慢）
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 DB 失败: %w", err)
	}

	// 配置连接池参数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// 【知识点】连接池的核心概念:
	// - MaxIdleConns（最大空闲连接）: 没有请求时保持的最小连接数
	//   太少：突发请求时要重新建连，延迟高
	//   太多：浪费 MySQL 资源
	// - MaxOpenConns（最大打开连接）: 同时最多有多少连接
	//   受 MySQL 的 max_connections 限制
	// - ConnMaxLifetime: 连接存活时间
	//   防止连接被 MySQL 服务端关闭后客户端还在用

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	log.Printf("MySQL 连接成功: %s", cfg.Host)

	DB = db
	return db, nil
}
