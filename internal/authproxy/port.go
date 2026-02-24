package authproxy

import (
	"fmt"
	"net"
	"strconv"
)

// FindAvailableListener 从 startPort 开始探测，返回第一个可用的 listener
// 直接返回 listener 而非端口号，避免 TOCTOU 竞态
func FindAvailableListener(startPort int) (net.Listener, error) {
	for p := startPort; p < startPort+100; p++ {
		ln, err := net.Listen("tcp", "127.0.0.1:"+strconv.Itoa(p))
		if err == nil {
			return ln, nil
		}
	}
	return nil, fmt.Errorf("在 %d-%d 范围内未找到可用端口", startPort, startPort+99)
}
