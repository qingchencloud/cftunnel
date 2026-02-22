package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

var addDomain string

func init() {
	addCmd.Flags().StringVar(&addDomain, "domain", "", "完整域名 (如 webhook.example.com)")
	addCmd.MarkFlagRequired("domain")
	rootCmd.AddCommand(addCmd)
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
			return fmt.Errorf("请先运行 cftunnel init")
		}
		if cfg.FindRoute(name) != nil {
			return fmt.Errorf("路由 %s 已存在", name)
		}

		client := cfapi.New(cfg.Auth.APIToken, cfg.Auth.AccountID)
		ctx := context.Background()

		// 提取主域名查找 Zone
		parts := strings.Split(addDomain, ".")
		if len(parts) < 2 {
			return fmt.Errorf("无效域名: %s", addDomain)
		}
		mainDomain := strings.Join(parts[len(parts)-2:], ".")
		zone, err := client.FindZoneByDomain(ctx, mainDomain)
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

		// 保存路由
		cfg.Routes = append(cfg.Routes, config.RouteConfig{
			Name:        name,
			Hostname:    addDomain,
			Service:     service,
			ZoneID:      zone.ID,
			DNSRecordID: recordID,
		})
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("路由已添加: %s → %s (%s)\n", addDomain, service, name)
		return nil
	},
}
