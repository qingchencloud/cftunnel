package cmd

import (
	"context"
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove <名称>",
	Short: "删除路由（清理 DNS + ingress）",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		route := cfg.FindRoute(name)
		if route == nil {
			return fmt.Errorf("路由 %s 不存在", name)
		}

		client := cfapi.New(cfg.Auth.APIToken, cfg.Auth.AccountID)
		ctx := context.Background()

		// 删除 DNS 记录
		if route.DNSRecordID != "" && route.ZoneID != "" {
			fmt.Printf("正在删除 DNS 记录 %s...\n", route.Hostname)
			if err := client.DeleteDNSRecord(ctx, route.ZoneID, route.DNSRecordID); err != nil {
				fmt.Printf("警告: 删除 DNS 记录失败: %v\n", err)
			}
		}

		cfg.RemoveRoute(name)
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("路由 %s 已删除\n", name)
		return nil
	},
}
