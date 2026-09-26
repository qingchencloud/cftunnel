package cmd

import (
	"fmt"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/presets"
	"github.com/spf13/cobra"
)

var (
	presetAuth  string
	presetRelay bool
	presetProto string
	presetShare shareFlags
)

func init() {
	presetCmd := &cobra.Command{
		Use:     "preset [名称]",
		Aliases: []string{"scenario", "template"},
		Short:   "按常用场景一键启动隧道",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || strings.EqualFold(args[0], "list") {
				fmt.Println("可用场景:")
				for _, item := range presets.All() {
					fmt.Printf("  %-14s %-10s 端口 %-5s  (%s)\n", item.ID, item.Label, item.Port, item.Hint)
				}
				fmt.Println("示例: cftunnel preset frontend --share")
				return nil
			}
			item, ok := presets.Find(args[0])
			if !ok {
				return fmt.Errorf("未知场景 %q，请执行 cftunnel preset list 查看可用场景", args[0])
			}
			if presetRelay && presetShare.enabled() {
				return fmt.Errorf("Relay 场景暂不支持二维码或浏览器分享，请使用 cftunnel share <公网地址>")
			}
			return startQuick(item.Port, presetAuth, presetRelay, presetProto, presetShare)
		},
	}
	presetCmd.Flags().StringVar(&presetAuth, "auth", "", "启用密码保护 (格式: 用户名:密码)")
	presetCmd.Flags().BoolVar(&presetRelay, "relay", false, "使用中继模式（需先 relay init）")
	presetCmd.Flags().StringVar(&presetProto, "proto", "tcp", "中继协议 (tcp/udp)，仅 --relay 时有效")
	presetShare.bind(presetCmd)
	rootCmd.AddCommand(presetCmd)
}
