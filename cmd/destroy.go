package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

var destroyForce bool

func init() {
	destroyCmd.Flags().BoolVar(&destroyForce, "force", false, "跳过确认")
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "删除隧道（清理所有 DNS 记录 + 删除 CF 隧道）",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.ID == "" {
			return fmt.Errorf("未初始化，无隧道可删除")
		}

		if !destroyForce {
			fmt.Printf("即将删除隧道 %s (%s) 及其 %d 条路由，此操作不可恢复！\n", cfg.Tunnel.Name, cfg.Tunnel.ID, len(cfg.Routes))
			fmt.Print("确认删除？(y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) != "y" {
				fmt.Println("已取消")
				return nil
			}
		}

		// 停止运行中的进程
		if daemon.Running() {
			fmt.Println("正在停止隧道...")
			daemon.Stop()
		}

		client := cfapi.New(cfg.Auth.APIToken, cfg.Auth.AccountID)
		ctx := context.Background()

		// 删除所有 DNS 记录
		for _, r := range cfg.Routes {
			if r.DNSRecordID != "" && r.ZoneID != "" {
				fmt.Printf("删除 DNS: %s\n", r.Hostname)
				if err := client.DeleteDNSRecord(ctx, r.ZoneID, r.DNSRecordID); err != nil {
					fmt.Printf("  警告: %v\n", err)
				}
			}
		}

		// 删除隧道
		fmt.Println("删除隧道...")
		if err := client.DeleteTunnel(ctx, cfg.Tunnel.ID); err != nil {
			fmt.Printf("警告: %v\n", err)
		}

		// 清空配置
		cfg.Tunnel = config.TunnelConfig{}
		cfg.Routes = nil
		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Println("隧道已删除，配置已清空")
		return nil
	},
}
