package daemon

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/qingchencloud/cftunnel/internal/config"
)

// CloudflaredPath 返回 cloudflared 二进制路径
func CloudflaredPath() string {
	return filepath.Join(config.Dir(), "bin", "cloudflared")
}

// EnsureCloudflared 确保 cloudflared 已安装，未安装则自动下载
func EnsureCloudflared() (string, error) {
	path := CloudflaredPath()
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	// 尝试系统 PATH
	if p, err := exec.LookPath("cloudflared"); err == nil {
		return p, nil
	}
	return path, download(path)
}

func download(dest string) error {
	url, err := downloadURL()
	if err != nil {
		return err
	}
	fmt.Printf("正在下载 cloudflared...\n")

	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}
	if err := os.Chmod(dest, 0755); err != nil {
		return err
	}
	fmt.Printf("cloudflared 已下载到 %s\n", dest)
	return nil
}

func downloadURL() (string, error) {
	const base = "https://github.com/cloudflare/cloudflared/releases/latest/download/"
	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "darwin/arm64":
		return base + "cloudflared-darwin-arm64.tgz", nil
	case "darwin/amd64":
		return base + "cloudflared-darwin-amd64.tgz", nil
	case "linux/amd64":
		return base + "cloudflared-linux-amd64", nil
	case "linux/arm64":
		return base + "cloudflared-linux-arm64", nil
	default:
		return "", fmt.Errorf("不支持的平台: %s/%s", runtime.GOOS, runtime.GOARCH)
	}
}
