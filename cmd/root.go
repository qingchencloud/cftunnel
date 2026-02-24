package cmd

import (
	"fmt"
	"os"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "cftunnel",
	Short:   "Cloudflare Tunnel 一键管理工具",
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkWindowsVersion()
		if config.Portable() {
			fmt.Printf("[便携模式] 数据目录: %s\n", config.Dir())
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
