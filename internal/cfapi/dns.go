package cfapi

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
)

// CreateCNAME 创建 CNAME 记录指向隧道
func (c *Client) CreateCNAME(ctx context.Context, zoneID, name, target string) (string, error) {
	record, err := c.api.DNS.Records.New(ctx, dns.RecordNewParams{
		ZoneID: cf.F(zoneID),
		Body: dns.CNAMERecordParam{
			Name:    cf.F(name),
			Content: cf.F(target),
			Type:    cf.F(dns.CNAMERecordTypeCNAME),
			TTL:     cf.F(dns.TTL(1)),
			Proxied: cf.F(true),
		},
	})
	if err != nil {
		return "", fmt.Errorf("创建 CNAME 记录失败: %w", err)
	}
	return record.ID, nil
}

// DeleteDNSRecord 删除 DNS 记录
func (c *Client) DeleteDNSRecord(ctx context.Context, zoneID, recordID string) error {
	_, err := c.api.DNS.Records.Delete(ctx, recordID, dns.RecordDeleteParams{
		ZoneID: cf.F(zoneID),
	})
	if err != nil {
		return fmt.Errorf("删除 DNS 记录失败: %w", err)
	}
	return nil
}
