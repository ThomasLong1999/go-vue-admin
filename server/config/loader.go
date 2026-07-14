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
	"path/filepath"

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

	// 设置配置文件名和路径
	filename := filepath.Base(configPath) // 提取文件名: "config.yaml"
	ext := filepath.Ext(filename)         // 提取扩展名: ".yaml"
	dir := filepath.Dir(configPath)       // 提取目录: "config"

	// 【知识点】filepath 包处理跨平台的路径问题
	// Windows 用 \ 分隔路径，Linux/Mac 用 / 分隔
	// filepath 会自动处理这个差异

	v.SetConfigName(filename[:len(filename)-len(ext)]) // 去掉扩展名: "config"
	v.SetConfigType(ext[1:])                           // 设置格式: "yaml"
	v.AddConfigPath(dir)                               // 搜索目录: "config"

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

	// 【知识点】log.Printf 是标准库的日志函数
	// %s 替换字符串, %d 替换整数
	log.Printf("配置加载成功，文件: %s", v.ConfigFileUsed())
	log.Printf("服务器端口: %d, 数据库: %s, Redis: %s",
		cfg.Server.Port, cfg.Database.DBName, cfg.Redis.Addr())

	return &cfg, nil
}
