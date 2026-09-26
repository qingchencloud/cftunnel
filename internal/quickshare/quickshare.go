// Package quickshare contains the small, local-only helpers shared by the
// terminal quick workflow and the desktop client's share experience.
package quickshare

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/skip2/go-qrcode"
)

const (
	TelegramGroupURL = "https://t.me/+-53et5QXFh0xYzhk"
	InviteURL        = "https://cftunnel.qt.cool/?utm_source=cli&utm_medium=share"
	maxRecentPorts   = 8
)

// RecentPort is a local history entry. It never contains credentials or
// public tunnel URLs, only the port and the mode used to start it.
type RecentPort struct {
	Port     string    `json:"port"`
	Mode     string    `json:"mode"`
	LastUsed time.Time `json:"last_used"`
}

func historyPath() string { return filepath.Join(config.Dir(), "recent-ports.json") }

// RecordPort stores a port in the local recent list. Failures are returned so
// callers can report them without preventing the tunnel from starting.
func RecordPort(port, mode string) error {
	port = strings.TrimSpace(port)
	if port == "" {
		return nil
	}
	items, err := RecentPorts()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	next := make([]RecentPort, 0, maxRecentPorts)
	next = append(next, RecentPort{Port: port, Mode: mode, LastUsed: now})
	for _, item := range items {
		if item.Port == port && item.Mode == mode {
			continue
		}
		next = append(next, item)
		if len(next) == maxRecentPorts {
			break
		}
	}
	if err := os.MkdirAll(config.Dir(), 0700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(historyPath(), append(b, '\n'), 0600)
}

// RecentPorts returns newest-first history. A missing history file is empty.
func RecentPorts() ([]RecentPort, error) {
	b, err := os.ReadFile(historyPath())
	if os.IsNotExist(err) {
		return []RecentPort{}, nil
	}
	if err != nil {
		return nil, err
	}
	var items []RecentPort
	if len(strings.TrimSpace(string(b))) == 0 {
		return []RecentPort{}, nil
	}
	if err := json.Unmarshal(b, &items); err != nil {
		return nil, fmt.Errorf("读取最近端口记录失败: %w", err)
	}
	return items, nil
}

// ClearRecentPorts deletes the local recent-port history.
func ClearRecentPorts() error {
	err := os.Remove(historyPath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// TelegramShareURL builds the standard Telegram share link without making a
// network request.
func TelegramShareURL(target, text string) string {
	values := url.Values{}
	values.Set("url", target)
	if strings.TrimSpace(text) != "" {
		values.Set("text", text)
	}
	return "https://t.me/share/url?" + values.Encode()
}

// QRText returns a terminal-friendly QR image for a URL.
func QRText(target string) (string, error) {
	qr, err := qrcode.New(target, qrcode.Medium)
	if err != nil {
		return "", err
	}
	return qr.ToSmallString(false), nil
}

// Copy copies text to the system clipboard when one is available.
func Copy(text string) error { return clipboard.WriteAll(text) }

// OpenURL asks the native desktop opener to open a URL. It is intentionally
// best-effort and is only called when the user explicitly asks for it.
func OpenURL(target string) error {
	var command string
	var args []string
	switch runtime.GOOS {
	case "windows":
		command, args = "rundll32", []string{"url.dll,FileProtocolHandler", target}
	case "darwin":
		command, args = "open", []string{target}
	default:
		command, args = "xdg-open", []string{target}
	}
	return exec.Command(command, args...).Start()
}
