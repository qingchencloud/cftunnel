package cmd

import (
	"fmt"

	"github.com/qingchencloud/cftunnel/internal/selfupdate"
	"github.com/spf13/cobra"
)

var checkUpdate bool

func init() {
	versionCmd.Flags().BoolVar(&checkUpdate, "check", false, "检查是否有新版本")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("cftunnel %s\n", Version)
		if !checkUpdate {
			return nil
		}
		fmt.Println("正在检查更新...")
		latest, err := selfupdate.LatestVersion()
		if err != nil {
			return fmt.Errorf("检查更新失败: %w", err)
		}
		if latest == "v"+Version || latest == Version {
			fmt.Println("已是最新版本")
		} else {
			fmt.Printf("发现新版本: %s → %s\n", Version, latest)
			fmt.Println("运行 cftunnel update 进行更新")
		}
		return nil
	},
}
