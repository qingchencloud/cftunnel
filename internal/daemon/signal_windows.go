//go:build windows

package daemon

import (
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	generateConsoleCtrl  = kernel32.NewProc("GenerateConsoleCtrlEvent")
)

// stopChildProcess 优雅终止子进程（Windows: CTRL_BREAK_EVENT）
func stopChildProcess(cmd *exec.Cmd) {
	pid := cmd.Process.Pid
	// 发送 CTRL_BREAK_EVENT，cloudflared 可捕获并优雅退出
	r, _, _ := generateConsoleCtrl.Call(uintptr(syscall.CTRL_BREAK_EVENT), uintptr(unsafe.Pointer(nil)))
	if r == 0 {
		// 降级为强制终止
		cmd.Process.Kill()
		return
	}
	_ = pid
}
