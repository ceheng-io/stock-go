package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// DragonTigerClient is the request client interface required by DragonTigerService.
type DragonTigerClient interface {
	eastmoney.DragonTigerClient
}

// DragonTigerService orchestrates Eastmoney dragon-tiger providers.
type DragonTigerService struct {
	client     DragonTigerClient
	datacenter string
}

// NewDragonTigerService creates a DragonTigerService.
func NewDragonTigerService(client DragonTigerClient, datacenterURL string) *DragonTigerService {
	return &DragonTigerService{client: client, datacenter: datacenterURL}
}

// Detail returns daily dragon-tiger billboard detail rows.
func (s *DragonTigerService) Detail(ctx context.Context, options eastmoney.DragonTigerDateOptions) ([]types.DragonTigerDetailItem, error) {
	return eastmoney.GetDragonTigerDetail(ctx, s.client, s.datacenter, options)
}

// StockStats returns stock billboard statistics rows.
func (s *DragonTigerService) StockStats(ctx context.Context, period eastmoney.DragonTigerPeriod) ([]types.DragonTigerStockStatItem, error) {
	return eastmoney.GetDragonTigerStockStats(ctx, s.client, s.datacenter, period)
}

// Institution returns institution trading rows.
func (s *DragonTigerService) Institution(ctx context.Context, options eastmoney.DragonTigerDateOptions) ([]types.DragonTigerInstitutionItem, error) {
	return eastmoney.GetDragonTigerInstitution(ctx, s.client, s.datacenter, options)
}

// BranchRank returns brokerage branch ranking rows.
func (s *DragonTigerService) BranchRank(ctx context.Context, period eastmoney.DragonTigerPeriod) ([]types.DragonTigerBranchItem, error) {
	return eastmoney.GetDragonTigerBranchRank(ctx, s.client, s.datacenter, period)
}

// SeatDetail returns buy and sell seat detail rows for a stock on a date.
func (s *DragonTigerService) SeatDetail(ctx context.Context, symbol string, date string) ([]types.DragonTigerSeatItem, error) {
	return eastmoney.GetDragonTigerStockSeatDetail(ctx, s.client, s.datacenter, symbol, date)
}
