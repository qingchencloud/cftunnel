package config

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version     int               `yaml:"version"`
	Auth        AuthConfig        `yaml:"auth"`
	Tunnel      TunnelConfig      `yaml:"tunnel"`
	Routes      []RouteConfig     `yaml:"routes"`
	Relay       RelayConfig       `yaml:"relay,omitempty"`
	Cloudflared CloudflaredConfig `yaml:"cloudflared"`
	SelfUpdate  SelfUpdateConfig  `yaml:"self_update"`
}

type AuthConfig struct {
	APIToken  string `yaml:"api_token"`
	AccountID string `yaml:"account_id"`
}

type TunnelConfig struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type RouteConfig struct {
	Name        string     `yaml:"name"`
	Hostname    string     `yaml:"hostname"`
	Service     string     `yaml:"service"`
	ZoneID      string     `yaml:"zone_id"`
	DNSRecordID string     `yaml:"dns_record_id"`
	Auth        *AuthProxy `yaml:"auth,omitempty"`
}

// AuthProxy 鉴权代理配置
type AuthProxy struct {
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	SigningKey string `yaml:"signing_key,omitempty"`
	CookieTTL int    `yaml:"cookie_ttl,omitempty"` // 秒，默认 86400
}

// CookieTTLOrDefault 返回 Cookie 有效期（秒），默认 86400
func (a *AuthProxy) CookieTTLOrDefault() int {
	if a.CookieTTL > 0 {
		return a.CookieTTL
	}
	return 86400
}

// RelayConfig 中继模式配置
type RelayConfig struct {
	Server string      `yaml:"server,omitempty"`
	Token  string      `yaml:"token,omitempty"`
	Rules  []RelayRule `yaml:"rules,omitempty"`
}

// RelayRule 中继穿透规则
type RelayRule struct {
	Name       string `yaml:"name"`
	Proto      string `yaml:"proto"`                   // tcp/udp/http/https/stcp
	LocalIP    string `yaml:"local_ip,omitempty"`       // 默认 127.0.0.1
	LocalPort  int    `yaml:"local_port"`
	RemotePort int    `yaml:"remote_port,omitempty"`    // HTTP 模式可选
	Domain     string `yaml:"domain,omitempty"`         // HTTP 模式用
}

type CloudflaredConfig struct {
	Path       string `yaml:"path"`
	AutoUpdate bool   `yaml:"auto_update"`
}

type SelfUpdateConfig struct {
	AutoCheck bool `yaml:"auto_check"` // 启动时自动检查 cftunnel 更新
}

var (
	dirOnce    sync.Once
	dirPath    string
	isPortable bool
)

// Dir 返回配置目录路径
// 便携模式：程序同级目录存在 portable 文件时，使用程序所在目录
// 普通模式：~/.cftunnel/
func Dir() string {
	dirOnce.Do(func() {
		if exe, err := os.Executable(); err == nil {
			if real, err := filepath.EvalSymlinks(exe); err == nil {
				exeDir := filepath.Dir(real)
				if _, err := os.Stat(filepath.Join(exeDir, "portable")); err == nil {
					dirPath = exeDir
					isPortable = true
					return
				}
			}
		}
		home, _ := os.UserHomeDir()
		dirPath = filepath.Join(home, ".cftunnel")
	})
	return dirPath
}

// Portable 返回当前是否处于便携模式
func Portable() bool {
	Dir() // 确保 dirOnce 已执行
	return isPortable
}

func Path() string {
	return filepath.Join(Dir(), "config.yml")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &Config{Version: 1}
			cfg.applyEnvOverrides()
			return cfg, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cfg.applyEnvOverrides()
	return &cfg, nil
}

// applyEnvOverrides 用环境变量覆盖配置（CI/CD 和 Docker 场景）
func (c *Config) applyEnvOverrides() {
	if v := os.Getenv("CFTUNNEL_API_TOKEN"); v != "" {
		c.Auth.APIToken = v
	}
	if v := os.Getenv("CFTUNNEL_ACCOUNT_ID"); v != "" {
		c.Auth.AccountID = v
	}
	if v := os.Getenv("CFTUNNEL_RELAY_SERVER"); v != "" {
		c.Relay.Server = v
	}
	if v := os.Getenv("CFTUNNEL_RELAY_TOKEN"); v != "" {
		c.Relay.Token = v
	}
}

func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0600)
}

func (c *Config) FindRoute(name string) *RouteConfig {
	for i := range c.Routes {
		if c.Routes[i].Name == name {
			return &c.Routes[i]
		}
	}
	return nil
}

func (c *Config) RemoveRoute(name string) bool {
	for i, r := range c.Routes {
		if r.Name == name {
			c.Routes = append(c.Routes[:i], c.Routes[i+1:]...)
			return true
		}
	}
	return false
}

// FindRelayRule 查找中继规则
func (c *Config) FindRelayRule(name string) *RelayRule {
	for i := range c.Relay.Rules {
		if c.Relay.Rules[i].Name == name {
			return &c.Relay.Rules[i]
		}
	}
	return nil
}

// RemoveRelayRule 删除中继规则
func (c *Config) RemoveRelayRule(name string) bool {
	for i, r := range c.Relay.Rules {
		if r.Name == name {
			c.Relay.Rules = append(c.Relay.Rules[:i], c.Relay.Rules[i+1:]...)
			return true
		}
	}
	return false
}
