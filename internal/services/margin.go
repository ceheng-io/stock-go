package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/types"
)

// MarginClient is the request client interface required by MarginService.
type MarginClient interface {
	eastmoney.MarginClient
}

// MarginService orchestrates Eastmoney margin providers.
type MarginService struct {
	client     MarginClient
	datacenter string
}

// NewMarginService creates a MarginService.
func NewMarginService(client MarginClient, datacenterURL string) *MarginService {
	return &MarginService{client: client, datacenter: datacenterURL}
}

// AccountInfo returns daily margin account statistics rows.
func (s *MarginService) AccountInfo(ctx context.Context) ([]types.MarginAccountItem, error) {
	return eastmoney.GetMarginAccountInfo(ctx, s.client, s.datacenter)
}

// TargetList returns stock margin target detail rows.
func (s *MarginService) TargetList(ctx context.Context, date string) ([]types.MarginTargetItem, error) {
	return eastmoney.GetMarginTargetList(ctx, s.client, s.datacenter, date)
}
