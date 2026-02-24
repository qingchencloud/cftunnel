package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/qingchencloud/cftunnel/internal/service"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "注册为系统服务（开机自启）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.Portable() {
			return fmt.Errorf("便携模式下不支持注册系统服务（路径不固定），请使用 cftunnel up 手动启动")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.Token == "" {
			return fmt.Errorf("请先运行 cftunnel init && cftunnel create <名称>")
		}
		binPath, err := daemon.EnsureCloudflared()
		if err != nil {
			return err
		}
		svc := service.New()
		if err := svc.Install(binPath, cfg.Tunnel.Token); err != nil {
			return fmt.Errorf("注册服务失败: %w", err)
		}
		fmt.Println("系统服务已注册，隧道将开机自启")
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "卸载系统服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := service.New()
		if err := svc.Uninstall(); err != nil {
			return fmt.Errorf("卸载服务失败: %w", err)
		}
		fmt.Println("系统服务已卸载")
		return nil
	},
}
