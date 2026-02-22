package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "启动隧道",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Tunnel.Token == "" {
			return fmt.Errorf("请先运行 cftunnel init")
		}
		return daemon.Start(cfg.Tunnel.Token)
	},
}
