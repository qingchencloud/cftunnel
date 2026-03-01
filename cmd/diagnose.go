package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/daemon"
	"github.com/spf13/cobra"
)

var diagnoseJSON bool

func init() {
	diagnoseCmd.Flags().BoolVar(&diagnoseJSON, "json", false, "JSON 格式输出")
	rootCmd.AddCommand(diagnoseCmd)
}

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "诊断 Cloud 模式链路连通性",
	Long:  "检测 cloudflared 状态、Cloudflare API 连通性、本地服务、DNS 解析和域名可达性。",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		var routes []daemon.RouteInput
		for _, r := range cfg.Routes {
			routes = append(routes, daemon.RouteInput{
				Name:     r.Name,
				Hostname: r.Hostname,
				Service:  r.Service,
			})
		}

		result := daemon.Diagnose(routes)

		if diagnoseJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		printDiagnose(result)
		return nil
	},
}

func printDiagnose(r daemon.DiagnoseResult) {
	fmt.Println("Cloud 链路诊断")
	fmt.Println("==============")

	// cloudflared 状态
	c := r.Cloudflared
	if c.Installed {
		fmt.Printf("cloudflared: ✓ 已安装 (%s)\n", c.Version)
		fmt.Printf("  路径: %s\n", c.Path)
		if c.Running {
			fmt.Printf("  进程: 运行中 (PID: %d)\n", c.PID)
		} else {
			fmt.Println("  进程: 未运行")
		}
	} else {
		fmt.Println("cloudflared: ✗ 未安装")
	}

	// API 连通性
	a := r.API
	if a.Reachable {
		fmt.Printf("Cloudflare API: ✓ 可达 (%dms)\n", a.LatencyMS)
	} else {
		fmt.Printf("Cloudflare API: ✗ %s\n", a.Err)
	}
	fmt.Println()

	if len(r.Routes) == 0 {
		fmt.Println("暂无路由需要检测")
		return
	}

	// 路由检测表格
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "路由\t域名\t本地服务\tDNS\tHTTPS")
	fmt.Fprintln(w, "----\t----\t--------\t---\t-----")
	for _, route := range r.Routes {
		local := "✓"
		if !route.LocalOK {
			local = "✗ " + route.LocalErr
		}
		dns := "✓"
		if !route.DNSOK {
			dns = "✗ " + route.DNSErr
		}
		https := "-"
		if route.DNSOK {
			if route.HTTPOK {
				https = "✓"
			} else {
				https = "✗ " + route.HTTPErr
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			route.Name, route.Hostname, local, dns, https)
	}
	w.Flush()

	fmt.Printf("\n结果: %d 条路由, %d 通 / %d 断\n", r.Total, r.Passed, r.Failed)
}
