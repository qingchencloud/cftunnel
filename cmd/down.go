package cmd

import (
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "停止隧道",
	RunE: func(cmd *cobra.Command, args []string) error {
		return daemon.Stop()
	},
}
