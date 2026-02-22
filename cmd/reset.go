package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

var resetForce bool

func init() {
	resetCmd.Flags().BoolVar(&resetForce, "force", false, "跳过确认")
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "重置全部（删除隧道 + 清除本地配置）",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !resetForce {
			fmt.Println("即将删除隧道并清除所有本地配置，此操作不可恢复！")
			fmt.Print("确认重置？(y/N): ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(input)) != "y" {
				fmt.Println("已取消")
				return nil
			}
		}

		// 先执行 destroy 逻辑（如果有隧道）
		cfg, _ := config.Load()
		if cfg != nil && cfg.Tunnel.ID != "" {
			destroyForce = true
			if err := destroyCmd.RunE(cmd, nil); err != nil {
				fmt.Printf("警告: 删除隧道失败: %v\n", err)
			}
		}

		// 删除整个配置目录
		dir := config.Dir()
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("清除配置目录失败: %w", err)
		}

		fmt.Printf("已清除 %s，回到全新状态\n", dir)
		fmt.Println("重新开始: cftunnel init")
		return nil
	},
}
