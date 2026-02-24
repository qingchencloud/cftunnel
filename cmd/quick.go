package cmd

import (
	"fmt"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

var quickAuth string

func init() {
	quickCmd.Flags().StringVar(&quickAuth, "auth", "", "启用密码保护 (格式: 用户名:密码)")
	rootCmd.AddCommand(quickCmd)
}

var quickCmd = &cobra.Command{
	Use:   "quick <端口>",
	Short: "快速启动免域名隧道（生成 *.trycloudflare.com 随机域名）",
	Long:  "无需 Cloudflare 账户、API Token 或域名，一条命令生成临时公网地址。\n适合临时分享、快速调试，Ctrl+C 退出后域名自动失效。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if quickAuth != "" {
			user, pass, err := parseAuth(quickAuth)
			if err != nil {
				return err
			}
			return daemon.StartQuickWithAuth(args[0], user, pass)
		}
		return daemon.StartQuick(args[0])
	},
}

// parseAuth 解析 "用户名:密码" 格式，密码部分允许包含冒号
func parseAuth(s string) (string, string, error) {
	idx := strings.Index(s, ":")
	if idx < 0 {
		return "", "", fmt.Errorf("--auth 格式错误，应为 用户名:密码")
	}
	user := s[:idx]
	pass := s[idx+1:]
	if user == "" || pass == "" {
		return "", "", fmt.Errorf("--auth 用户名和密码不能为空")
	}
	return user, pass, nil
}
