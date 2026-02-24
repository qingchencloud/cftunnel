//go:build windows

package cmd

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

// checkWindowsVersion 检测 Windows 版本，低于 Win10 则警告退出
func checkWindowsVersion() {
	ver := windows.RtlGetVersion()
	if ver.MajorVersion < 10 {
		fmt.Fprintf(os.Stderr, "错误: 当前系统 Windows NT %d.%d 不受支持\n",
			ver.MajorVersion, ver.MinorVersion)
		fmt.Fprintln(os.Stderr, "cftunnel 依赖的 cloudflared 和 Go 运行时均要求 Windows 10 或更高版本。")
		fmt.Fprintln(os.Stderr, "Windows 7/8/8.1 用户请参考: https://github.com/qingchencloud/cftunnel/issues/8")
		os.Exit(1)
	}
}
