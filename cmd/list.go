package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有路由",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if len(cfg.Routes) == 0 {
			fmt.Println("暂无路由")
			return nil
		}
		fmt.Printf("%-12s %-30s %s\n", "名称", "域名", "服务")
		for _, r := range cfg.Routes {
			fmt.Printf("%-12s %-30s %s\n", r.Name, r.Hostname, r.Service)
		}
		return nil
	},
}
