package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qingchencloud/cftunnel/internal/authproxy"
	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/qingchencloud/cftunnel/internal/selfupdate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "启动隧道",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.Token == "" {
			return fmt.Errorf("请先运行 cftunnel init && cftunnel create <名称>")
		}

		// 为有鉴权配置的路由启动代理
		var proxies []*authproxy.Proxy
		for i, r := range cfg.Routes {
			if r.Auth == nil {
				continue
			}
			sigKey, err := hex.DecodeString(r.Auth.SigningKey)
			if err != nil {
				return fmt.Errorf("路由 %s 的 signing_key 无效: %w", r.Name, err)
			}
			// 从 service URL 提取端口
			port := extractPort(r.Service)
			if port == "" {
				return fmt.Errorf("路由 %s 的 service 格式无效: %s", r.Name, r.Service)
			}
			proxy, err := authproxy.New(authproxy.Config{
				Username:   r.Auth.Username,
				Password:   r.Auth.Password,
				TargetPort: port,
				SigningKey:  sigKey,
				CookieTTL:  time.Duration(r.Auth.CookieTTLOrDefault()) * time.Second,
			})
			if err != nil {
				return fmt.Errorf("路由 %s 启动鉴权代理失败: %w", r.Name, err)
			}
			if err := proxy.Start(); err != nil {
				return fmt.Errorf("路由 %s 启动鉴权代理失败: %w", r.Name, err)
			}
			proxies = append(proxies, proxy)
			proxyPort := strconv.Itoa(proxy.ListenPort())
			fmt.Printf("鉴权代理已启动: %s → 127.0.0.1:%s → 127.0.0.1:%s\n", r.Hostname, proxyPort, port)
			// 临时修改 service 指向代理端口（仅内存，不持久化）
			cfg.Routes[i].Service = "http://localhost:" + proxyPort
		}
		// 确保退出时关闭所有代理
		defer func() {
			for _, p := range proxies {
				p.Stop()
			}
		}()

		// 启动前同步 ingress 配置到远端，确保本地与远端一致
		if len(cfg.Routes) > 0 {
			client := cfapi.New(cfg.Auth.APIToken, cfg.Auth.AccountID)
			if err := pushIngress(client, context.Background(), cfg); err != nil {
				fmt.Printf("警告: 同步 ingress 失败: %v（将使用远端现有配置）\n", err)
			} else {
				fmt.Println("ingress 配置已同步")
			}
		}

		// 自动检查更新（非阻塞，仅提示）
		if cfg.SelfUpdate.AutoCheck {
			if latest, err := selfupdate.LatestVersion(); err == nil {
				if latest != "v"+Version && latest != Version {
					fmt.Printf("发现新版本: %s → %s (运行 cftunnel update 更新)\n", Version, latest)
				}
			}
		}
		return daemon.Start(cfg.Tunnel.Token)
	},
}

// extractPort 从 "http://localhost:3000" 格式中提取端口号
func extractPort(service string) string {
	idx := strings.LastIndex(service, ":")
	if idx < 0 {
		return ""
	}
	return service[idx+1:]
}
