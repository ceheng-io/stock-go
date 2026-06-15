package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// DividendClient is the request client interface required by DividendService.
type DividendClient interface {
	eastmoney.DividendClient
}

// DividendService orchestrates Eastmoney dividend providers.
type DividendService struct {
	client     DividendClient
	datacenter string
}

// NewDividendService creates a DividendService.
func NewDividendService(client DividendClient, datacenterURL string) *DividendService {
	return &DividendService{client: client, datacenter: datacenterURL}
}

// Detail returns stock dividend detail rows.
func (s *DividendService) Detail(ctx context.Context, symbol string) ([]types.DividendDetail, error) {
	return eastmoney.GetDividendDetail(ctx, s.client, s.datacenter, symbol)
}
