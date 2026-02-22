//go:build darwin

package service

import (
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Launchd struct{}

const plistName = "com.cftunnel.cloudflared"

func (l *Launchd) plistPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "Library/LaunchAgents", plistName+".plist")
}

const plistTmpl = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Label}}</string>
    <key>ProgramArguments</key>
    <array>
        <string>{{.BinPath}}</string>
        <string>tunnel</string>
        <string>run</string>
        <string>--token</string>
        <string>{{.Token}}</string>
    </array>
    <key>KeepAlive</key>
    <true/>
    <key>RunAtLoad</key>
    <true/>
    <key>StandardOutPath</key>
    <string>{{.LogPath}}</string>
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}</string>
</dict>
</plist>
`

func (l *Launchd) Install(binPath, token string) error {
	home, _ := os.UserHomeDir()
	data := map[string]string{
		"Label":   plistName,
		"BinPath": binPath,
		"Token":   token,
		"LogPath": filepath.Join(home, "Library/Logs/cftunnel.log"),
	}
	f, err := os.Create(l.plistPath())
	if err != nil {
		return err
	}
	defer f.Close()
	if err := template.Must(template.New("").Parse(plistTmpl)).Execute(f, data); err != nil {
		return err
	}
	return exec.Command("launchctl", "load", l.plistPath()).Run()
}

func (l *Launchd) Uninstall() error {
	exec.Command("launchctl", "unload", l.plistPath()).Run()
	return os.Remove(l.plistPath())
}

func (l *Launchd) Running() bool {
	out, err := exec.Command("launchctl", "list", plistName).Output()
	return err == nil && len(out) > 0
}

func New() Service {
	return &Launchd{}
}
