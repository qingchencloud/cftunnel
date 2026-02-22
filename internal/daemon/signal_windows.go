//go:build windows

package daemon

import (
	"os/exec"
	"syscall"
)

var (
	kernel32            = syscall.NewLazyDLL("kernel32.dll")
	generateConsoleCtrl = kernel32.NewProc("GenerateConsoleCtrlEvent")
)

// stopChildProcess 优雅终止子进程（Windows: CTRL_BREAK_EVENT）
func stopChildProcess(cmd *exec.Cmd) {
	pid := uintptr(cmd.Process.Pid)
	r, _, _ := generateConsoleCtrl.Call(uintptr(syscall.CTRL_BREAK_EVENT), pid)
	if r == 0 {
		cmd.Process.Kill()
	}
}
