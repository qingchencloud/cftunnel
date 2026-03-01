package daemon

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/qingchencloud/cftunnel/internal/authproxy"
	"github.com/qingchencloud/cftunnel/internal/config"
)

// quickConfigPath 返回 quick 模式专用的空配置文件路径
// 防止 cloudflared 读取用户已有的 ~/.cloudflared/config.yml 导致 UUID 解析失败
func quickConfigPath() string {
	p := filepath.Join(config.Dir(), "quick-config.yml")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		os.MkdirAll(config.Dir(), 0700)
		os.WriteFile(p, []byte("# cftunnel quick mode - empty config\n"), 0600)
	}
	return p
}

// StartQuick 启动免域名模式（前台运行，Ctrl+C 退出）
func StartQuick(port string) error {
	binPath, err := EnsureCloudflared()
	if err != nil {
		return err
	}
	if Running() {
		return fmt.Errorf("cloudflared 已在运行，请先执行 cftunnel down")
	}

	// 显式指定空配置文件，防止 cloudflared 读取用户已有的 ~/.cloudflared/config.yml
	// 避免残留的 tunnel: 字段触发 UUID 解析失败 (issue #13)
	cfgPath := quickConfigPath()
	cmd := exec.Command(binPath, "tunnel", "--config", cfgPath, "--url", "http://localhost:"+port)

	// 捕获 stderr 提取随机域名
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 cloudflared 失败: %w", err)
	}

	// 后台读取 stderr，提取域名并转发输出
	go scanForURL(stderr)

	// 捕获 Ctrl+C 优雅退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-sig:
		stopChildProcess(cmd)
		<-done
	case err := <-done:
		if err != nil {
			return fmt.Errorf("cloudflared 异常退出: %w", err)
		}
	}
	return nil
}

func scanForURL(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		// cloudflared 输出格式: ... https://xxx.trycloudflare.com ...
		if strings.Contains(line, "trycloudflare.com") {
			url := extractURL(line)
			if url != "" {
				fmt.Printf("\n✔ 隧道已启动: %s\n\n", url)
			}
		}
		fmt.Fprintln(os.Stderr, line)
	}
}

func extractURL(line string) string {
	for _, part := range strings.Fields(line) {
		if strings.Contains(part, "trycloudflare.com") && strings.HasPrefix(part, "http") {
			return part
		}
	}
	return ""
}

// StartQuickWithAuth 启动带鉴权代理的免域名模式
func StartQuickWithAuth(port, username, password string) error {
	binPath, err := EnsureCloudflared()
	if err != nil {
		return err
	}
	if Running() {
		return fmt.Errorf("cloudflared 已在运行，请先执行 cftunnel down")
	}

	// 启动鉴权代理
	proxy, err := authproxy.New(authproxy.Config{
		Username:   username,
		Password:   password,
		TargetPort: port,
		SigningKey:  authproxy.RandomKey(),
		CookieTTL:  24 * time.Hour,
	})
	if err != nil {
		return fmt.Errorf("启动鉴权代理失败: %w", err)
	}
	if err := proxy.Start(); err != nil {
		return fmt.Errorf("启动鉴权代理失败: %w", err)
	}
	defer proxy.Stop()

	proxyPort := fmt.Sprintf("%d", proxy.ListenPort())
	fmt.Printf("鉴权代理已启动 127.0.0.1:%s → 127.0.0.1:%s\n", proxyPort, port)

	// cloudflared 指向代理端口（同样隔离配置文件）
	cfgPath := quickConfigPath()
	cmd := exec.Command(binPath, "tunnel", "--config", cfgPath, "--url", "http://localhost:"+proxyPort)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 cloudflared 失败: %w", err)
	}

	go scanForURL(stderr)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	select {
	case <-sig:
		stopChildProcess(cmd)
		<-done
	case err := <-done:
		if err != nil {
			return fmt.Errorf("cloudflared 异常退出: %w", err)
		}
	}
	return nil
}
