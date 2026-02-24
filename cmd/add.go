package cmd

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/authproxy"
	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

var addDomain string
var addAuth string

func init() {
	addCmd.Flags().StringVar(&addDomain, "domain", "", "完整域名 (如 webhook.example.com)")
	addCmd.MarkFlagRequired("domain")
	addCmd.Flags().StringVar(&addAuth, "auth", "", "启用密码保护 (格式: 用户名:密码)")
	rootCmd.AddCommand(addCmd)
}

// pushIngress 推送当前所有路由的 ingress 配置到远端
func pushIngress(client *cfapi.Client, ctx context.Context, cfg *config.Config) error {
	var rules []cfapi.IngressRule
	for _, r := range cfg.Routes {
		rules = append(rules, cfapi.IngressRule{Hostname: r.Hostname, Service: r.Service})
	}
	return client.PushIngressConfig(ctx, cfg.Tunnel.ID, rules)
}

// findZoneForDomain 通过遍历账户 Zone 列表匹配域名（支持多级 TLD）
func findZoneForDomain(client *cfapi.Client, ctx context.Context, domain string) (*cfapi.ZoneInfo, error) {
	zoneList, err := client.ListZones(ctx)
	if err != nil {
		return nil, err
	}
	for _, z := range zoneList {
		if domain == z.Name || strings.HasSuffix(domain, "."+z.Name) {
			return &cfapi.ZoneInfo{ID: z.ID, Name: z.Name}, nil
		}
	}
	return nil, fmt.Errorf("未找到域名 %s 对应的 Zone，请确认域名已添加到 Cloudflare", domain)
}

var addCmd = &cobra.Command{
	Use:   "add <名称> <端口>",
	Short: "添加路由（自动创建 CNAME + 更新 ingress）",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, port := args[0], args[1]
		service := "http://localhost:" + port

		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.ID == "" {
			return fmt.Errorf("请先运行 cftunnel init && cftunnel create <名称>")
		}
		if cfg.FindRoute(name) != nil {
			return fmt.Errorf("路由 %s 已存在", name)
		}

		client := cfapi.New(cfg.Auth.APIToken, cfg.Auth.AccountID)
		ctx := context.Background()

		// 查找域名对应的 Zone（支持多级 TLD）
		zone, err := findZoneForDomain(client, ctx, addDomain)
		if err != nil {
			return err
		}

		// 创建 CNAME
		target := cfg.Tunnel.ID + ".cfargotunnel.com"
		fmt.Printf("正在创建 DNS 记录 %s → %s\n", addDomain, target)
		recordID, err := client.CreateCNAME(ctx, zone.ID, addDomain, target)
		if err != nil {
			return err
		}

		// 构建路由配置
		route := config.RouteConfig{
			Name:        name,
			Hostname:    addDomain,
			Service:     service,
			ZoneID:      zone.ID,
			DNSRecordID: recordID,
		}

		// 如果指定了 --auth，填充鉴权配置
		if addAuth != "" {
			user, pass, err := parseAuth(addAuth)
			if err != nil {
				return err
			}
			route.Auth = &config.AuthProxy{
				Username:   user,
				Password:   pass,
				SigningKey:  hex.EncodeToString(authproxy.RandomKey()),
			}
			fmt.Printf("已启用密码保护: %s\n", addDomain)
		}

		// 保存路由
		cfg.Routes = append(cfg.Routes, route)
		if err := cfg.Save(); err != nil {
			return err
		}

		// 推送 ingress 配置到远端
		fmt.Println("正在同步 ingress 配置...")
		if err := pushIngress(client, ctx, cfg); err != nil {
			return fmt.Errorf("推送 ingress 失败: %w（DNS 记录已创建，请排查后重试 add 或手动删除 DNS 记录）", err)
		}

		fmt.Printf("路由已添加: %s → %s (%s)\n", addDomain, service, name)
		return nil
	},
}
