package relay

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/qingchencloud/cftunnel/internal/config"
)

const checkTimeout = 3 * time.Second

// CheckResult 链路检测总结果
type CheckResult struct {
	Server        string            `json:"server"`
	ServerOK      bool              `json:"server_ok"`
	ServerLatency int64             `json:"server_latency_ms"`
	FrpcRunning   bool              `json:"frpc_running"`
	FrpcPID       int               `json:"frpc_pid"`
	Rules         []RuleCheckResult `json:"rules"`
	Total         int               `json:"total"`
	Passed        int               `json:"passed"`
	Failed        int               `json:"failed"`
}

// RuleCheckResult 单条规则检测结果
type RuleCheckResult struct {
	Name       string `json:"name"`
	Proto      string `json:"proto"`
	LocalPort  int    `json:"local_port"`
	RemotePort int    `json:"remote_port"`
	LocalOK    bool   `json:"local_ok"`
	RemoteOK   bool   `json:"remote_ok"`
	LatencyMS  int64  `json:"latency_ms"`
	LocalErr   string `json:"local_err,omitempty"`
	RemoteErr  string `json:"remote_err,omitempty"`
}

// Check 执行链路检测
func Check(cfg *config.RelayConfig, ruleName string) CheckResult {
	result := CheckResult{Server: cfg.Server}

	// 检测 frps 服务器连通性
	if cfg.Server != "" {
		start := time.Now()
		conn, err := net.DialTimeout("tcp", cfg.Server, checkTimeout)
		if err == nil {
			conn.Close()
			result.ServerOK = true
			result.ServerLatency = time.Since(start).Milliseconds()
		}
	}

	// 检测 frpc 进程
	result.FrpcRunning = Running()
	result.FrpcPID = PID()

	// 筛选要检测的规则
	rules := cfg.Rules
	if ruleName != "" {
		rules = nil
		for _, r := range cfg.Rules {
			if r.Name == ruleName {
				rules = append(rules, r)
				break
			}
		}
	}

	// 并行检测所有规则
	result.Rules = make([]RuleCheckResult, len(rules))
	var wg sync.WaitGroup
	for i, rule := range rules {
		wg.Add(1)
		go func(idx int, r config.RelayRule) {
			defer wg.Done()
			result.Rules[idx] = checkRule(r, cfg.Server)
		}(i, rule)
	}
	wg.Wait()

	// 统计
	result.Total = len(result.Rules)
	for _, r := range result.Rules {
		if r.LocalOK && (r.RemoteOK || r.RemotePort == 0) {
			result.Passed++
		} else {
			result.Failed++
		}
	}
	return result
}

// checkRule 检测单条规则
func checkRule(r config.RelayRule, server string) RuleCheckResult {
	rc := RuleCheckResult{
		Name:       r.Name,
		Proto:      r.Proto,
		LocalPort:  r.LocalPort,
		RemotePort: r.RemotePort,
	}

	localIP := r.LocalIP
	if localIP == "" {
		localIP = "127.0.0.1"
	}

	// 检测本地服务
	localAddr := fmt.Sprintf("%s:%d", localIP, r.LocalPort)
	conn, err := net.DialTimeout("tcp", localAddr, checkTimeout)
	if err == nil {
		conn.Close()
		rc.LocalOK = true
	} else {
		rc.LocalErr = "未监听"
	}

	// 检测远程穿透端口
	if r.RemotePort > 0 && server != "" {
		host, _, _ := net.SplitHostPort(server)
		remoteAddr := fmt.Sprintf("%s:%d", host, r.RemotePort)
		start := time.Now()
		conn, err := net.DialTimeout("tcp", remoteAddr, checkTimeout)
		if err == nil {
			conn.Close()
			rc.RemoteOK = true
			rc.LatencyMS = time.Since(start).Milliseconds()
		} else {
			rc.RemoteErr = "超时"
		}
	}

	return rc
}
