package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/qingchencloud/cftunnel/internal/config"
)

var pidFile = config.Dir() + "/cloudflared.pid"

// Start 启动 cloudflared（token 模式）
func Start(token string) error {
	binPath, err := EnsureCloudflared()
	if err != nil {
		return err
	}
	if Running() {
		return fmt.Errorf("cloudflared 已在运行")
	}

	cmd := exec.Command(binPath, "tunnel", "run", "--token", token)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 cloudflared 失败: %w", err)
	}

	os.MkdirAll(config.Dir(), 0700)
	os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0600)
	fmt.Printf("cloudflared 已启动 (PID: %d)\n", cmd.Process.Pid)
	return nil
}

// Stop 停止 cloudflared
func Stop() error {
	pid, err := readPID()
	if err != nil {
		return fmt.Errorf("未找到运行中的 cloudflared")
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("停止 cloudflared 失败: %w", err)
	}
	os.Remove(pidFile)
	fmt.Println("cloudflared 已停止")
	return nil
}

// Running 检查 cloudflared 是否在运行
func Running() bool {
	pid, err := readPID()
	if err != nil {
		return false
	}
	// 用 kill -0 检测进程是否存在
	return exec.Command("kill", "-0", strconv.Itoa(pid)).Run() == nil
}

// PID 返回当前运行的 PID
func PID() int {
	pid, _ := readPID()
	return pid
}

func readPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}
