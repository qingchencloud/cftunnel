package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有路由和规则",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		hasCloud := len(cfg.Routes) > 0
		hasRelay := len(cfg.Relay.Rules) > 0

		if !hasCloud && !hasRelay {
			fmt.Println("暂无路由或规则")
			return nil
		}

		if hasCloud {
			fmt.Println("Cloud 路由:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "名称\t域名\t服务\t鉴权")
			fmt.Fprintln(w, "----\t----\t----\t----")
			for _, r := range cfg.Routes {
				auth := "-"
				if r.Auth != nil {
					auth = "✓"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Name, r.Hostname, r.Service, auth)
			}
			w.Flush()
		}

		if hasCloud && hasRelay {
			fmt.Println()
		}

		if hasRelay {
			fmt.Println("Relay 规则:")
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "名称\t协议\t本地端口\t远程端口\t域名")
			fmt.Fprintln(w, "----\t----\t--------\t--------\t----")
			for _, r := range cfg.Relay.Rules {
				remote := "-"
				if r.RemotePort > 0 {
					remote = fmt.Sprintf("%d", r.RemotePort)
				}
				domain := "-"
				if r.Domain != "" {
					domain = r.Domain
				}
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n",
					r.Name, r.Proto, r.LocalPort, remote, domain)
			}
			w.Flush()
		}

		return nil
	},
}
