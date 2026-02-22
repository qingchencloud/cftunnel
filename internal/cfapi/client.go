package cfapi

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

type Client struct {
	api       *cf.Client
	accountID string
}

func New(apiToken, accountID string) *Client {
	return &Client{
		api:       cf.NewClient(option.WithAPIToken(apiToken)),
		accountID: accountID,
	}
}

func (c *Client) AccountID() string  { return c.accountID }
func (c *Client) API() *cf.Client    { return c.api }

// ListZones 列出账户下所有域名
func (c *Client) ListZones(ctx context.Context) ([]zones.Zone, error) {
	pager := c.api.Zones.ListAutoPaging(ctx, zones.ZoneListParams{})
	var result []zones.Zone
	for pager.Next() {
		result = append(result, pager.Current())
	}
	if err := pager.Err(); err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %w", err)
	}
	return result, nil
}

// FindZoneByDomain 根据域名查找 Zone
func (c *Client) FindZoneByDomain(ctx context.Context, domain string) (*zones.Zone, error) {
	page, err := c.api.Zones.List(ctx, zones.ZoneListParams{
		Name: cf.F(domain),
	})
	if err != nil {
		return nil, fmt.Errorf("查找域名失败: %w", err)
	}
	if len(page.Result) == 0 {
		return nil, fmt.Errorf("未找到域名 %s", domain)
	}
	return &page.Result[0], nil
}
