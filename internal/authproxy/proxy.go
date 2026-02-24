package authproxy

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//go:embed login.html
var loginHTML []byte

const cookieName = "__cftunnel_auth"
const loginPath = "/___auth/login"

// RandomKey 生成 32 字节随机签名密钥
func RandomKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

// Config 鉴权代理配置
type Config struct {
	Username   string
	Password   string
	TargetPort string
	SigningKey  []byte
	CookieTTL  time.Duration
}

// Proxy 鉴权反向代理
type Proxy struct {
	cfg      Config
	listener net.Listener
	server   *http.Server
	reverse  *httputil.ReverseProxy
}

// New 创建鉴权代理实例，自动探测可用端口
func New(cfg Config) (*Proxy, error) {
	port, _ := strconv.Atoi(cfg.TargetPort)
	ln, err := FindAvailableListener(port + 1)
	if err != nil {
		return nil, err
	}

	target, _ := url.Parse("http://127.0.0.1:" + cfg.TargetPort)
	rp := httputil.NewSingleHostReverseProxy(target)

	if cfg.CookieTTL == 0 {
		cfg.CookieTTL = 24 * time.Hour
	}

	p := &Proxy{
		cfg:      cfg,
		listener: ln,
		reverse:  rp,
	}
	p.server = &http.Server{Handler: p}
	return p, nil
}

// ListenPort 返回代理实际监听的端口
func (p *Proxy) ListenPort() int {
	return p.listener.Addr().(*net.TCPAddr).Port
}

// Start 非阻塞启动代理
func (p *Proxy) Start() error {
	go p.server.Serve(p.listener)
	return nil
}

// Stop 优雅关闭代理
func (p *Proxy) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return p.server.Shutdown(ctx)
}

// ServeHTTP 核心路由逻辑
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// WebSocket 升级请求直接透传
	if isWebSocket(r) {
		p.reverse.ServeHTTP(w, r)
		return
	}

	// 登录表单提交
	if r.Method == http.MethodPost && r.URL.Path == loginPath {
		p.handleLogin(w, r)
		return
	}

	// 检查 Cookie 鉴权
	if p.checkAuth(r) {
		p.reverse.ServeHTTP(w, r)
		return
	}

	// 未认证，返回登录页
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(loginHTML)
}

// handleLogin 处理登录表单提交
func (p *Proxy) handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != p.cfg.Username || password != p.cfg.Password {
		http.Redirect(w, r, "/?error=1", http.StatusSeeOther)
		return
	}

	// 签发 Cookie
	expiry := time.Now().Add(p.cfg.CookieTTL).Unix()
	payload := fmt.Sprintf("%s:%x", username, expiry)
	sig := signPayload(p.cfg.SigningKey, payload)
	value := payload + "." + sig

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int(p.cfg.CookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// checkAuth 校验请求中的鉴权 Cookie
func (p *Proxy) checkAuth(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}

	// 格式：username:expiry_hex.hmac_hex
	dotIdx := strings.LastIndex(cookie.Value, ".")
	if dotIdx < 0 {
		return false
	}
	payload := cookie.Value[:dotIdx]
	sig := cookie.Value[dotIdx+1:]

	// 验证签名
	if signPayload(p.cfg.SigningKey, payload) != sig {
		return false
	}

	// 验证过期时间
	colonIdx := strings.LastIndex(payload, ":")
	if colonIdx < 0 {
		return false
	}
	expiryHex := payload[colonIdx+1:]
	expiry, err := strconv.ParseInt(expiryHex, 16, 64)
	if err != nil {
		return false
	}
	return time.Now().Unix() < expiry
}

// signPayload 使用 HMAC-SHA256 签名
func signPayload(key []byte, payload string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// isWebSocket 检测是否为 WebSocket 升级请求
func isWebSocket(r *http.Request) bool {
	return strings.EqualFold(r.Header.Get("Upgrade"), "websocket")
}
