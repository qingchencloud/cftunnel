package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/qingchencloud/cftunnel/internal/quickshare"
	"github.com/qingchencloud/cftunnel/internal/relay"
	"github.com/spf13/cobra"
)

var (
	quickAuth  string
	quickRelay bool
	quickProto string
	quickShare shareFlags
)

func init() {
	quickCmd.Flags().StringVar(&quickAuth, "auth", "", "启用密码保护 (格式: 用户名:密码)")
	quickCmd.Flags().BoolVar(&quickRelay, "relay", false, "使用中继模式穿透（需先 relay init）")
	quickCmd.Flags().StringVar(&quickProto, "proto", "tcp", "中继协议 (tcp/udp)，仅 --relay 时有效")
	quickShare.bind(quickCmd)
	rootCmd.AddCommand(quickCmd)
}

var quickCmd = &cobra.Command{
	Use:   "quick <端口>",
	Short: "快速启动免域名隧道（生成 *.trycloudflare.com 随机域名）",
	Long:  "无需 Cloudflare 账户、API Token 或域名，一条命令生成临时公网地址。\n适合临时分享、快速调试，Ctrl+C 退出后域名自动失效。",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if quickRelay && quickShare.enabled() {
			return fmt.Errorf("--qr/--telegram/--copy/--open/--share 目前仅支持 Cloud quick 模式；Relay 模式请使用 cftunnel share <公网地址>")
		}
		return startQuick(args[0], quickAuth, quickRelay, quickProto, quickShare)
	},
}

func startQuick(port, auth string, relayMode bool, proto string, flags shareFlags) error {
	if relayMode && auth != "" {
		return fmt.Errorf("--auth 仅支持 Cloud quick 模式；Relay 模式请在 relay init 中配置 token")
	}
	var username, password string
	if auth != "" {
		var err error
		username, password, err = parseAuth(auth)
		if err != nil {
			return err
		}
	}
	mode := "cloud"
	if relayMode {
		mode = "relay"
	}
	if err := quickshare.RecordPort(port, mode); err != nil {
		fmt.Fprintf(os.Stderr, "记录最近端口失败（不影响启动）: %v\n", err)
	}

	if relayMode {
		return relay.StartQuick(port, proto)
	}
	onURL := func(publicURL string) {
		if flags.enabled() {
			presentShare(publicURL, flags)
		}
	}
	if auth != "" {
		return daemon.StartQuickWithAuthCallback(port, username, password, onURL)
	}
	return daemon.StartQuickWithCallback(port, onURL)
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
