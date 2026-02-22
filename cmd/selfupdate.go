package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/selfupdate"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新 cftunnel 到最新版本",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("正在检查更新...")
		latest, err := selfupdate.LatestVersion()
		if err != nil {
			return fmt.Errorf("检查更新失败: %w", err)
		}
		if latest == "v"+Version || latest == Version {
			fmt.Println("已是最新版本")
			return nil
		}
		fmt.Printf("发现新版本: %s → %s\n", Version, latest)
		fmt.Println("正在下载...")
		if err := selfupdate.Update(latest); err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}
		fmt.Printf("已更新到 %s\n", latest)
		return nil
	},
}
