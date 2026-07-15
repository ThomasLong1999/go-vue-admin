// ============================================
// loader.go - 配置文件加载器
// ============================================
// 【知识点】Go 推荐把不同功能的代码放在同一包的不同文件中
// config 包有两个文件: config.go（定义结构体）和 loader.go（加载逻辑）
// 同一个包下所有文件共享包名，可以互相访问未导出（小写）的内容

package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Load 从指定路径加载配置文件
// 【知识点】函数参数类型写在变量名后面: path string
// 这和其他语言 (string path) 刚好相反
//
// 参数:
//   - configPath: 配置文件路径，如 "config/config.yaml"
//
// 返回:
//   - *Config: 加载好的配置对象
//   - error: 如果加载失败，返回错误信息
//
// 【知识点】Go 的多返回值是语言核心特性
// 习惯上: 最后一个返回值是 error 类型，调用者必须检查
func Load(configPath string) (*Config, error) {
	// 【知识点】viper 是 Go 生态最流行的配置管理库
	// 支持多种格式（JSON/YAML/TOML/ENV）和多种来源（文件/环境变量/命令行参数）

	v := viper.New()

	// SetConfigFile 直接指定文件路径，比手动拆分文件名和目录更直观。
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	// 环境变量使用大写下划线形式：jwt.secret → JWT_SECRET。
	// BindEnv 能让环境变量也参与 Unmarshal，而不仅是 v.Get()。
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	for _, key := range []string{
		"server.port", "server.mode",
		"database.host", "database.port", "database.user", "database.password", "database.dbname",
		"database.max_idle_conns", "database.max_open_conns", "database.conn_max_lifetime",
		"redis.host", "redis.port", "redis.password", "redis.db",
		"jwt.secret", "jwt.expire_hours",
		"bootstrap.admin_username", "bootstrap.admin_password",
	} {
		if err := v.BindEnv(key); err != nil {
			return nil, fmt.Errorf("绑定环境变量 %s 失败: %w", key, err)
		}
	}

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 【知识点】%w 是 Go 1.13+ 的错误包装语法
		// 它把原始错误包装成新错误，保留错误链（可以用 errors.Unwrap 解开）
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 【知识点】viper 的 Unmarshal 会根据 struct tag 自动把 YAML 映射到 Go 结构体
	// 这个过程叫"反序列化"（deserialize）
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置校验失败: %w", err)
	}

	// 【知识点】log.Printf 是标准库的日志函数
	// %s 替换字符串, %d 替换整数
	log.Printf("配置加载成功，文件: %s", v.ConfigFileUsed())
	log.Printf("服务器端口: %d, 数据库: %s, Redis: %s",
		cfg.Server.Port, cfg.Database.DBName, cfg.Redis.Addr())

	return &cfg, nil
}
