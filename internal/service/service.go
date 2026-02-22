package service

// Service 系统服务管理接口
type Service interface {
	Install(binPath, token string) error
	Uninstall() error
	Running() bool
}
