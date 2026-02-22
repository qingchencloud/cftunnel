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

var initToken, initAccountID, initTunnelName string

func init() {
	initCmd.Flags().StringVar(&initToken, "token", "", "API 令牌")
	initCmd.Flags().StringVar(&initAccountID, "account", "", "账户 ID")
	initCmd.Flags().StringVar(&initTunnelName, "name", "my-tunnel", "隧道名称")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "交互式初始化：输入 Token → 选域名 → 创建隧道",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("=== Cloudflare Tunnel 初始化向导 ===")
		fmt.Println()
		fmt.Println("需要以下信息（首次使用请先创建 API 令牌）：")
		fmt.Println()
		fmt.Println("  1. API 令牌获取方式:")
		fmt.Println("     登录 Cloudflare → 右上角头像 → 我的个人资料 → API 令牌 → 创建令牌")
		fmt.Println("     选择「创建自定义令牌」→「开始使用」")
		fmt.Println()
		fmt.Println("     添加 3 条权限（点「+ 添加更多」逐条添加）：")
		fmt.Println("     ┌──────────────────────────────────────────────────┐")
		fmt.Println("     │ 第 1 行: 帐户 │ Cloudflare Tunnel │ 编辑       │")
		fmt.Println("     │ 第 2 行: 区域 │ DNS              │ 编辑       │")
		fmt.Println("     │ 第 3 行: 区域 │ 区域设置          │ 读取       │")
		fmt.Println("     └──────────────────────────────────────────────────┘")
		fmt.Println("     提示: 第 2、3 行需先将左侧「帐户」切换为「区域」")
		fmt.Println()
		fmt.Println("     区域资源 → 包括 → 特定区域 → 选择你的域名")
		fmt.Println()
		fmt.Println("  2. 账户 ID 获取方式（任选其一）:")
		fmt.Println("     方式 A: 首页 → 点击域名 → 页面右下角「API」区域 →「账户 ID」")
		fmt.Println("     方式 B: 首页 → 账户名称旁「⋯」→ 复制账户 ID")
		fmt.Println()

		apiToken := strings.TrimSpace(initToken)
		accountID := strings.TrimSpace(initAccountID)
		tunnelName := strings.TrimSpace(initTunnelName)

		// 有 flag 则跳过 TUI
		if apiToken == "" || accountID == "" {
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
		}
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
