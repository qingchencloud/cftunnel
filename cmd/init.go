package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/qingchencloud/cftunnel/internal/cfapi"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "交互式初始化：输入 Token → 选域名 → 创建隧道",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== Cloudflare Tunnel 初始化向导 ===")
		fmt.Println()
		fmt.Println("需要以下信息（首次使用请先创建 API Token）：")
		fmt.Println()
		fmt.Println("  1. API Token 获取方式:")
		fmt.Println("     打开 https://dash.cloudflare.com/profile/api-tokens")
		fmt.Println("     → Create Token → Custom Token")
		fmt.Println("     → 权限: Account/Cloudflare Tunnel/Edit + Zone/Zone/Read + Zone/DNS/Edit")
		fmt.Println()
		fmt.Println("  2. Account ID 获取方式:")
		fmt.Println("     打开任意域名的 Overview 页面，右侧栏即可看到")
		fmt.Println()

		var apiToken, accountID, tunnelName string

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("API 令牌 (API Token)").Value(&apiToken).Placeholder("在 Cloudflare 控制台创建"),
				huh.NewInput().Title("账户 ID (Account ID)").Value(&accountID).Placeholder("32 位十六进制字符串"),
				huh.NewInput().Title("隧道名称").Value(&tunnelName).Placeholder("my-tunnel"),
			),
		).Run()
		if err != nil {
			return err
		}

		apiToken = strings.TrimSpace(apiToken)
		accountID = strings.TrimSpace(accountID)
		tunnelName = strings.TrimSpace(tunnelName)
		if tunnelName == "" {
			tunnelName = "my-tunnel"
		}

		client := cfapi.New(apiToken, accountID)
		ctx := context.Background()

		// 创建隧道
		fmt.Println("正在创建隧道...")
		tunnel, err := client.CreateTunnel(ctx, tunnelName)
		if err != nil {
			return err
		}
		fmt.Printf("隧道已创建: %s (%s)\n", tunnel.Name, tunnel.ID)

		// 获取 Token
		token, err := client.GetTunnelToken(ctx, tunnel.ID)
		if err != nil {
			return err
		}

		// 保存配置
		cfg := &config.Config{
			Version: 1,
			Auth:    config.AuthConfig{APIToken: apiToken, AccountID: accountID},
			Tunnel:  config.TunnelConfig{ID: tunnel.ID, Name: tunnel.Name, Token: token},
		}
		if err := cfg.Save(); err != nil {
			return err
		}
		fmt.Printf("配置已保存到 %s\n", config.Path())
		fmt.Println("\n下一步: cftunnel add <名称> <端口> --domain <域名>")
		return nil
	},
}
