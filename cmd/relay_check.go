package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/qingchencloud/cftunnel/internal/config"
	"github.com/qingchencloud/cftunnel/internal/relay"
	"github.com/spf13/cobra"
)

var checkJSON bool

func init() {
	relayCheckCmd.Flags().BoolVar(&checkJSON, "json", false, "JSON 格式输出")
	relayCmd.AddCommand(relayCheckCmd)
}

var relayCheckCmd = &cobra.Command{
	Use:   "check [规则名]",
	Short: "检测中继链路连通性",
	Long:  "检测 frps 服务器、本地服务、远程穿透端口的连通性和延迟。不指定规则名则检测全部规则。",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if cfg.Relay.Server == "" {
			return fmt.Errorf("未配置中继服务器，请先执行 cftunnel relay init")
		}

		ruleName := ""
		if len(args) > 0 {
			ruleName = args[0]
		}

		result := relay.Check(&cfg.Relay, ruleName)

		if checkJSON {
			return printCheckJSON(result)
		}
		printCheckTable(result)
		return nil
	},
}

func printCheckTable(r relay.CheckResult) {
	fmt.Println("中继链路检测")
	fmt.Println("============")

	// 服务器状态
	if r.ServerOK {
		fmt.Printf("服务器: %s  ✓ 可达 (%dms)\n", r.Server, r.ServerLatency)
	} else {
		fmt.Printf("服务器: %s  ✗ 不可达\n", r.Server)
	}

	// frpc 进程状态
	if r.FrpcRunning {
		fmt.Printf("frpc:   运行中 (PID: %d)\n", r.FrpcPID)
	} else {
		fmt.Println("frpc:   未运行")
	}
	fmt.Println()

	if len(r.Rules) == 0 {
		fmt.Println("暂无规则需要检测")
		return
	}

	// 规则检测结果表格
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "规则\t协议\t本地端口\t远程端口\t本地服务\t远程穿透\t延迟")
	fmt.Fprintln(w, "----\t----\t--------\t--------\t--------\t--------\t----")
	for _, rule := range r.Rules {
		local := "✓"
		if !rule.LocalOK {
			local = "✗ " + rule.LocalErr
		}
		remote := "-"
		if rule.RemotePort > 0 {
			if rule.RemoteOK {
				remote = "✓"
			} else {
				remote = "✗ " + rule.RemoteErr
			}
		}
		latency := "-"
		if rule.LatencyMS > 0 {
			latency = fmt.Sprintf("%dms", rule.LatencyMS)
		}
		remotePort := "-"
		if rule.RemotePort > 0 {
			remotePort = fmt.Sprintf("%d", rule.RemotePort)
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			rule.Name, rule.Proto, rule.LocalPort, remotePort, local, remote, latency)
	}
	w.Flush()

	fmt.Printf("\n结果: %d 条规则, %d 通 / %d 断\n", r.Total, r.Passed, r.Failed)
}

func printCheckJSON(r relay.CheckResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
