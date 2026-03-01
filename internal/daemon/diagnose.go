package daemon

import (
	"net"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const diagnoseTimeout = 5 * time.Second

// DiagnoseResult 诊断总结果
type DiagnoseResult struct {
	Cloudflared CloudflaredCheck  `json:"cloudflared"`
	API         APICheck          `json:"api"`
	Routes      []RouteDiagnose   `json:"routes"`
	Total       int               `json:"total"`
	Passed      int               `json:"passed"`
	Failed      int               `json:"failed"`
}

// CloudflaredCheck cloudflared 二进制和进程检测
type CloudflaredCheck struct {
	Installed bool   `json:"installed"`
	Path      string `json:"path,omitempty"`
	Version   string `json:"version,omitempty"`
	Running   bool   `json:"running"`
	PID       int    `json:"pid,omitempty"`
}

// APICheck Cloudflare API 连通性
type APICheck struct {
	Reachable bool   `json:"reachable"`
	LatencyMS int64  `json:"latency_ms,omitempty"`
	Err       string `json:"err,omitempty"`
}

// RouteDiagnose 单条路由诊断
type RouteDiagnose struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
	Service  string `json:"service"`
	LocalOK  bool   `json:"local_ok"`
	LocalErr string `json:"local_err,omitempty"`
	DNSOK    bool   `json:"dns_ok"`
	DNSErr   string `json:"dns_err,omitempty"`
	HTTPOK   bool   `json:"http_ok"`
	HTTPErr  string `json:"http_err,omitempty"`
}

// Diagnose 执行 Cloud 模式链路诊断
func Diagnose(routes []RouteInput) DiagnoseResult {
	var result DiagnoseResult

	// 检测 cloudflared
	result.Cloudflared = checkCloudflared()

	// 检测 Cloudflare API
	result.API = checkAPI()

	// 并行检测路由
	result.Routes = make([]RouteDiagnose, len(routes))
	var wg sync.WaitGroup
	for i, r := range routes {
		wg.Add(1)
		go func(idx int, route RouteInput) {
			defer wg.Done()
			result.Routes[idx] = diagnoseRoute(route)
		}(i, r)
	}
	wg.Wait()

	// 统计
	result.Total = len(result.Routes)
	for _, r := range result.Routes {
		if r.LocalOK && r.DNSOK {
			result.Passed++
		} else {
			result.Failed++
		}
	}
	return result
}

// RouteInput 诊断输入
type RouteInput struct {
	Name     string
	Hostname string
	Service  string
}

func checkCloudflared() CloudflaredCheck {
	var c CloudflaredCheck
	path, err := EnsureCloudflared()
	if err != nil {
		return c
	}
	c.Installed = true
	c.Path = path

	// 获取版本
	out, err := exec.Command(path, "version").CombinedOutput()
	if err == nil {
		c.Version = strings.TrimSpace(string(out))
	}

	c.Running = Running()
	if c.Running {
		c.PID = PID()
	}
	return c
}

func checkAPI() APICheck {
	var a APICheck
	client := &http.Client{Timeout: diagnoseTimeout}
	start := time.Now()
	resp, err := client.Get("https://api.cloudflare.com/client/v4/user/tokens/verify")
	if err != nil {
		a.Err = "无法连接"
		return a
	}
	defer resp.Body.Close()
	a.Reachable = true
	a.LatencyMS = time.Since(start).Milliseconds()
	return a
}

func diagnoseRoute(r RouteInput) RouteDiagnose {
	d := RouteDiagnose{
		Name:     r.Name,
		Hostname: r.Hostname,
		Service:  r.Service,
	}

	// 检测本地服务
	port := extractPort(r.Service)
	if port != "" {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, diagnoseTimeout)
		if err == nil {
			conn.Close()
			d.LocalOK = true
		} else {
			d.LocalErr = "未监听"
		}
	} else {
		d.LocalErr = "无法解析端口"
	}

	// 检测 DNS
	if r.Hostname != "" {
		_, err := net.LookupHost(r.Hostname)
		if err == nil {
			d.DNSOK = true
		} else {
			d.DNSErr = "解析失败"
		}
	}

	// 检测 HTTPS 可达性
	if r.Hostname != "" && d.DNSOK {
		client := &http.Client{Timeout: diagnoseTimeout}
		resp, err := client.Get("https://" + r.Hostname)
		if err == nil {
			resp.Body.Close()
			d.HTTPOK = true
		} else {
			d.HTTPErr = "不可达"
		}
	}

	return d
}

// extractPort 从 service 字符串提取端口号（如 http://localhost:3000 → 3000）
func extractPort(service string) string {
	// 去掉协议前缀
	s := service
	for _, prefix := range []string{"https://", "http://"} {
		s = strings.TrimPrefix(s, prefix)
	}
	// 取 host:port 中的 port
	_, port, err := net.SplitHostPort(s)
	if err != nil {
		return ""
	}
	return port
}
