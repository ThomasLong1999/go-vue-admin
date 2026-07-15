package config

import (
	"os"
	"path/filepath"
	"testing"
)

func testConfig() Config {
	return Config{
		Server:   ServerConfig{Port: 8080, Mode: "debug"},
		Database: DatabaseConfig{Host: "127.0.0.1", Port: 3306, User: "root", DBName: "demo"},
		Redis:    RedisConfig{Host: "127.0.0.1", Port: 6379},
		JWT:      JWTConfig{Secret: "12345678901234567890123456789012", ExpireHours: 24},
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Config)
	}{
		{name: "valid", edit: func(*Config) {}},
		{name: "invalid port", edit: func(cfg *Config) { cfg.Server.Port = 0 }},
		{name: "invalid mode", edit: func(cfg *Config) { cfg.Server.Mode = "test" }},
		{name: "short jwt secret", edit: func(cfg *Config) { cfg.JWT.Secret = "short" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := testConfig()
			tt.edit(&cfg)
			err := cfg.Validate()
			if tt.name == "valid" && err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
			if tt.name != "valid" && err == nil {
				t.Fatal("Validate() error = nil, want validation error")
			}
		})
	}
}

func TestLoadLetsEnvironmentOverrideFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte(`
server: { port: 8080, mode: debug }
database: { host: 127.0.0.1, port: 3306, user: root, password: from-file, dbname: demo, max_idle_conns: 10, max_open_conns: 100, conn_max_lifetime: 3600 }
redis: { host: 127.0.0.1, port: 6379, password: "", db: 0 }
jwt: { secret: 12345678901234567890123456789012, expire_hours: 24 }
bootstrap: { admin_username: admin, admin_password: admin123 }
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	t.Setenv("JWT_SECRET", "environment-override-secret-123456789")

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.JWT.Secret != "environment-override-secret-123456789" {
		t.Fatalf("jwt.secret = %q, want environment value", cfg.JWT.Secret)
	}
}
