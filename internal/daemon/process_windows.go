//go:build windows

package daemon

import (
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// processRunning 检查进程是否存活（Windows: tasklist）
func processRunning(pid int) bool {
	out, err := exec.Command("tasklist", "/FI", "PID eq "+strconv.Itoa(pid), "/NH").Output()
	if err != nil {
		return false
	}
	return !strings.Contains(string(out), "No tasks")
}

// processKill 终止进程（Windows: CTRL_BREAK_EVENT，降级为 taskkill /F）
func processKill(pid int) error {
	r, _, _ := generateConsoleCtrl.Call(uintptr(syscall.CTRL_BREAK_EVENT), uintptr(pid))
	if r == 0 {
		return exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/F").Run()
	}
	return nil
}
