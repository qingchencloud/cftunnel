package cmd

import (
	"fmt"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/quickshare"
	"github.com/spf13/cobra"
)

func init() {
	historyCmd := &cobra.Command{
		Use:     "history [clear]",
		Aliases: []string{"recent"},
		Short:   "查看或清空最近使用的端口",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if strings.EqualFold(args[0], "clear") {
					if err := quickshare.ClearRecentPorts(); err != nil {
						return err
					}
					fmt.Println("最近端口记录已清空")
					return nil
				}
				return fmt.Errorf("未知操作 %q，仅支持 clear", args[0])
			}
			items, err := quickshare.RecentPorts()
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Println("暂无最近端口记录")
				return nil
			}
			fmt.Println("最近使用的端口（执行 cftunnel quick <端口> 可再次启动）:")
			for i, item := range items {
				fmt.Printf("%d. %-6s %-5s %s\n", i+1, item.Port, item.Mode, item.LastUsed.Local().Format("2006-01-02 15:04"))
			}
			return nil
		},
	}
	rootCmd.AddCommand(historyCmd)
}
