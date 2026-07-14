// ============================================
// config.go - 配置结构体定义
// ============================================
// 【知识点】Go 的 struct（结构体）类似其他语言的 class
// - 用 type 关键字定义新类型
// - 字段名首字母大写表示公开（public），小写表示私有（private）
// - struct tag（反引号包裹的部分）用于给字段附加元数据

package config

// 【知识点】import 导入外部包，用 "包路径" 引入，代码中用最后一段名称调用
// 例如 fmt.Sprintf 调用 fmt 包的 Sprintf 函数
import "fmt"

// Config 是整个应用的配置根结构体
// 它把所有子配置组织在一起，方便统一管理
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig HTTP 服务器相关配置
type ServerConfig struct {
	Port int    `mapstructure:"port"` // 服务监听端口
	Mode string `mapstructure:"mode"` // Gin 运行模式
}

// DatabaseConfig MySQL 数据库配置
// 【知识点】mapstructure tag 告诉 viper 库：从 YAML 读取时用哪个字段名映射
// 例如 YAML 里的 server.port 会自动映射到 ServerConfig.Port
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// JWTConfig JWT Token 配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// DSN 返回 MySQL 的连接字符串（Data Source Name）
// 【知识点】Go 的方法定义在 struct 外部，通过 (c *DatabaseConfig) 绑定
// 这种写法叫"方法接收者"，*表示指针接收者，可以修改原对象
func (c *DatabaseConfig) DSN() string {
	// 拼接   出 MySQL 连接字符串:
	// root:123456@tcp(127.0.0.1:3306)/go_vue_admin?charset=utf8mb4&parseTime=True&loc=Local
	//
	// 【知识点】各个参数含义:
	// - charset=utf8mb4: 支持完整的 Unicode（包括 emoji）
	// - parseTime=True: 让驱动自动把时间解析成 time.Time 类型
	// - loc=Local: 使用系统本地时区
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" +
		fmt.Sprintf("%d", c.Port) + ")/" + c.DBName +
		"?charset=utf8mb4&parseTime=True&loc=Local"
}

// Addr 返回 Redis 连接地址，如 "127.0.0.1:6379"
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
