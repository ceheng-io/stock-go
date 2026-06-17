package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/types"
)

// FundFlowClient is the request client interface required by FundFlowService.
type FundFlowClient interface {
	eastmoney.FundFlowClient
}

// FundFlowURLs contains Eastmoney fund-flow endpoints.
type FundFlowURLs struct {
	FFlow string
	Clist string
}

// FundFlowService orchestrates Eastmoney fund-flow providers.
type FundFlowService struct {
	client FundFlowClient
	urls   FundFlowURLs
}

// NewFundFlowService creates a FundFlowService.
func NewFundFlowService(client FundFlowClient, urls FundFlowURLs) *FundFlowService {
	return &FundFlowService{client: client, urls: urls}
}

// Individual returns stock fund-flow history rows.
func (s *FundFlowService) Individual(ctx context.Context, symbol string, options eastmoney.FundFlowOptions) ([]types.StockFundFlow, error) {
	return eastmoney.GetIndividualFundFlow(ctx, s.client, symbol, s.urls.FFlow, options)
}

// Market returns market fund-flow rows.
func (s *FundFlowService) Market(ctx context.Context) ([]types.MarketFundFlow, error) {
	return eastmoney.GetMarketFundFlow(ctx, s.client, s.urls.FFlow)
}

// Rank returns stock fund-flow ranking rows.
func (s *FundFlowService) Rank(ctx context.Context, options eastmoney.FundFlowRankOptions) ([]types.FundFlowRankItem, error) {
	return eastmoney.GetFundFlowRank(ctx, s.client, s.urls.Clist, options)
}

// SectorRank returns sector fund-flow ranking rows.
func (s *FundFlowService) SectorRank(ctx context.Context, options eastmoney.FundFlowRankOptions) ([]types.SectorFundFlowItem, error) {
	return eastmoney.GetSectorFundFlowRank(ctx, s.client, s.urls.Clist, options)
}

// SectorHistory returns sector fund-flow history rows.
func (s *FundFlowService) SectorHistory(ctx context.Context, symbol string, options eastmoney.FundFlowOptions) ([]types.StockFundFlow, error) {
	return eastmoney.GetSectorFundFlowHistory(ctx, s.client, symbol, s.urls.FFlow, options)
}
