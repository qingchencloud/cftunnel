package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看隧道状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.ID == "" {
			fmt.Println("未初始化，请运行 cftunnel init")
			return nil
		}
		fmt.Printf("隧道: %s (%s)\n", cfg.Tunnel.Name, cfg.Tunnel.ID)
		if daemon.Running() {
			fmt.Printf("状态: 运行中 (PID: %d)\n", daemon.PID())
		} else {
			fmt.Println("状态: 已停止")
		}
		fmt.Printf("路由: %d 条\n", len(cfg.Routes))
		for _, r := range cfg.Routes {
			fmt.Printf("  %s → %s\n", r.Hostname, r.Service)
		}
		return nil
	},
}
