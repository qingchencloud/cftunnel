package selfupdate

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

const repo = "qingchencloud/cftunnel"

type release struct {
	TagName string `json:"tag_name"`
}

// LatestVersion 查询 GitHub 最新版本
func LatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("查询失败: HTTP %d", resp.StatusCode)
	}
	var r release
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}
	return r.TagName, nil
}

// Update 下载最新版本替换自身
func Update(version string) error {
	url := fmt.Sprintf("https://github.com/%s/releases/download/%s/cftunnel_%s_%s.tar.gz",
		repo, version, runtime.GOOS, runtime.GOARCH)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	// 解压 tar.gz 提取 cftunnel 二进制
	gr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("解压失败: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return fmt.Errorf("tar.gz 中未找到 cftunnel 二进制")
		}
		if err != nil {
			return fmt.Errorf("解压失败: %w", err)
		}
		if hdr.Name == "cftunnel" {
			break
		}
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	tmp := exe + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, tr); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	if err := os.Chmod(tmp, 0755); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, exe)
}
