package config

import (
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version     int               `yaml:"version"`
	Auth        AuthConfig        `yaml:"auth"`
	Tunnel      TunnelConfig      `yaml:"tunnel"`
	Routes      []RouteConfig     `yaml:"routes"`
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
	Name        string `yaml:"name"`
	Hostname    string `yaml:"hostname"`
	Service     string `yaml:"service"`
	ZoneID      string `yaml:"zone_id"`
	DNSRecordID string `yaml:"dns_record_id"`
}

type CloudflaredConfig struct {
	Path       string `yaml:"path"`
	AutoUpdate bool   `yaml:"auto_update"`
}

type SelfUpdateConfig struct {
	AutoCheck bool `yaml:"auto_check"` // 启动时自动检查 cftunnel 更新
}

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cftunnel")
}

func Path() string {
	return filepath.Join(Dir(), "config.yml")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{Version: 1}, nil
		}
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
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
