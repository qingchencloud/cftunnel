package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/qingchencloud/cftunnel/internal/relay"
	"github.com/spf13/cobra"
)

var statusJSON bool

func init() {
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "JSON 格式输出")
	rootCmd.AddCommand(statusCmd)
}

// StatusOutput status 命令的结构化输出
type StatusOutput struct {
	Cloud *CloudStatus `json:"cloud,omitempty"`
	Relay *RelayStatus `json:"relay,omitempty"`
}

// CloudStatus Cloud 模式状态
type CloudStatus struct {
	TunnelName string        `json:"tunnel_name"`
	TunnelID   string        `json:"tunnel_id"`
	Running    bool          `json:"running"`
	PID        int           `json:"pid,omitempty"`
	Routes     []RouteStatus `json:"routes"`
}

// RouteStatus 路由状态
type RouteStatus struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	Auth     bool   `json:"auth"`
}

// RelayStatus Relay 模式状态
type RelayStatus struct {
	Server  string       `json:"server"`
	Running bool         `json:"running"`
	PID     int          `json:"pid,omitempty"`
	Rules   []RuleStatus `json:"rules"`
}

// RuleStatus 规则状态
type RuleStatus struct {
	Name       string `json:"name"`
	Proto      string `json:"proto"`
	LocalPort  int    `json:"local_port"`
	RemotePort int    `json:"remote_port,omitempty"`
	Domain     string `json:"domain,omitempty"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看隧道状态（Cloud + Relay）",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		out := buildStatus(cfg)

		if statusJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(out)
		}

		printStatus(out)
		return nil
	},
}

func buildStatus(cfg *config.Config) StatusOutput {
	var out StatusOutput

	if cfg.Tunnel.ID != "" {
		cs := &CloudStatus{
			TunnelName: cfg.Tunnel.Name,
			TunnelID:   cfg.Tunnel.ID,
			Running:    daemon.Running(),
		}
		if cs.Running {
			cs.PID = daemon.PID()
		}
		for _, r := range cfg.Routes {
			cs.Routes = append(cs.Routes, RouteStatus{
				Name:     r.Name,
				Hostname: r.Hostname,
				Service:  r.Service,
				Auth:     r.Auth != nil,
			})
		}
		out.Cloud = cs
	}

	if cfg.Relay.Server != "" {
		rs := &RelayStatus{
			Server:  cfg.Relay.Server,
			Running: relay.Running(),
		}
		if rs.Running {
			rs.PID = relay.PID()
		}
		for _, r := range cfg.Relay.Rules {
			rs.Rules = append(rs.Rules, RuleStatus{
				Name:       r.Name,
				Proto:      r.Proto,
				LocalPort:  r.LocalPort,
				RemotePort: r.RemotePort,
				Domain:     r.Domain,
			})
		}
		out.Relay = rs
	}

	return out
}

func printStatus(out StatusOutput) {
	if out.Cloud == nil && out.Relay == nil {
		fmt.Println("未配置任何模式，请运行 cftunnel init 或 cftunnel relay init")
		return
	}

	if out.Cloud != nil {
		cs := out.Cloud
		fmt.Println("Cloud 模式")
		fmt.Printf("  隧道: %s (%s)\n", cs.TunnelName, cs.TunnelID)
		if cs.Running {
			fmt.Printf("  状态: ✓ 运行中 (PID: %d)\n", cs.PID)
		} else {
			fmt.Println("  状态: ✗ 已停止")
		}
		fmt.Printf("  路由: %d 条\n", len(cs.Routes))
		for _, r := range cs.Routes {
			auth := ""
			if r.Auth {
				auth = " [鉴权]"
			}
			fmt.Printf("    %s → %s%s\n", r.Hostname, r.Service, auth)
		}
	}

	if out.Cloud != nil && out.Relay != nil {
		fmt.Println()
	}

	if out.Relay != nil {
		rs := out.Relay
		fmt.Println("Relay 模式")
		fmt.Printf("  服务器: %s\n", rs.Server)
		if rs.Running {
			fmt.Printf("  状态:   ✓ 运行中 (PID: %d)\n", rs.PID)
		} else {
			fmt.Println("  状态:   ✗ 已停止")
		}
		fmt.Printf("  规则:   %d 条\n", len(rs.Rules))
		for _, r := range rs.Rules {
			remote := "-"
			if r.RemotePort > 0 {
				remote = fmt.Sprintf(":%d", r.RemotePort)
			}
			if r.Domain != "" {
				remote = r.Domain
			}
			fmt.Printf("    %-8s %s  :%d → %s\n", r.Name, r.Proto, r.LocalPort, remote)
		}
	}
}
