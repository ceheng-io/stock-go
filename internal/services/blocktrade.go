package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// BlockTradeClient is the request client interface required by BlockTradeService.
type BlockTradeClient interface {
	eastmoney.BlockTradeClient
}

// BlockTradeService orchestrates Eastmoney block-trade providers.
type BlockTradeService struct {
	client     BlockTradeClient
	datacenter string
}

// NewBlockTradeService creates a BlockTradeService.
func NewBlockTradeService(client BlockTradeClient, datacenterURL string) *BlockTradeService {
	return &BlockTradeService{client: client, datacenter: datacenterURL}
}

// MarketStat returns block-trade market statistics rows.
func (s *BlockTradeService) MarketStat(ctx context.Context) ([]types.BlockTradeMarketStatItem, error) {
	return eastmoney.GetBlockTradeMarketStat(ctx, s.client, s.datacenter)
}

// Detail returns block-trade deal detail rows.
func (s *BlockTradeService) Detail(ctx context.Context, options eastmoney.BlockTradeDateOptions) ([]types.BlockTradeDetailItem, error) {
	return eastmoney.GetBlockTradeDetail(ctx, s.client, s.datacenter, options)
}

// DailyStat returns block-trade daily stock statistics rows.
func (s *BlockTradeService) DailyStat(ctx context.Context, options eastmoney.BlockTradeDateOptions) ([]types.BlockTradeDailyStatItem, error) {
	return eastmoney.GetBlockTradeDailyStat(ctx, s.client, s.datacenter, options)
}
