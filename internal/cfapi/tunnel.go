package cfapi

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/shared"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
)

// CreateTunnel 创建 Cloudflare Tunnel
func (c *Client) CreateTunnel(ctx context.Context, name string) (*shared.CloudflareTunnel, error) {
	tunnel, err := c.api.ZeroTrust.Tunnels.Cloudflared.New(ctx, zero_trust.TunnelCloudflaredNewParams{
		AccountID: cf.F(c.accountID),
		Name:      cf.F(name),
		ConfigSrc: cf.F(zero_trust.TunnelCloudflaredNewParamsConfigSrcCloudflare),
	})
	if err != nil {
		return nil, fmt.Errorf("创建隧道失败: %w", err)
	}
	return tunnel, nil
}

// DeleteTunnel 删除隧道
func (c *Client) DeleteTunnel(ctx context.Context, tunnelID string) error {
	_, err := c.api.ZeroTrust.Tunnels.Cloudflared.Delete(ctx, tunnelID, zero_trust.TunnelCloudflaredDeleteParams{
		AccountID: cf.F(c.accountID),
	})
	if err != nil {
		return fmt.Errorf("删除隧道失败: %w", err)
	}
	return nil
}

// ListTunnels 列出所有隧道
func (c *Client) ListTunnels(ctx context.Context) ([]shared.CloudflareTunnel, error) {
	pager := c.api.ZeroTrust.Tunnels.Cloudflared.ListAutoPaging(ctx, zero_trust.TunnelCloudflaredListParams{
		AccountID: cf.F(c.accountID),
	})
	var result []shared.CloudflareTunnel
	for pager.Next() {
		result = append(result, pager.Current())
	}
	if err := pager.Err(); err != nil {
		return nil, fmt.Errorf("列出隧道失败: %w", err)
	}
	return result, nil
}

// GetTunnelToken 获取隧道运行 Token
func (c *Client) GetTunnelToken(ctx context.Context, tunnelID string) (string, error) {
	token, err := c.api.ZeroTrust.Tunnels.Cloudflared.Token.Get(ctx, tunnelID, zero_trust.TunnelCloudflaredTokenGetParams{
		AccountID: cf.F(c.accountID),
	})
	if err != nil {
		return "", fmt.Errorf("获取隧道 Token 失败: %w", err)
	}
	return *token, nil
}
