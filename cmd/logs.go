package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var follow bool

func init() {
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "实时跟踪日志")
	rootCmd.AddCommand(logsCmd)
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "查看隧道日志",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, _ := os.UserHomeDir()
		logFile := filepath.Join(home, "Library/Logs/cftunnel.log")
		if _, err := os.Stat(logFile); err != nil {
			return fmt.Errorf("日志文件不存在: %s", logFile)
		}
		tailArgs := []string{"-100", logFile}
		if follow {
			tailArgs = []string{"-100", "-f", logFile}
		}
		c := exec.Command("tail", tailArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}
