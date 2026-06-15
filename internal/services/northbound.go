package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// NorthboundClient is the request client interface required by NorthboundService.
type NorthboundClient interface {
	eastmoney.NorthboundClient
}

// NorthboundURLs contains Eastmoney northbound endpoints.
type NorthboundURLs struct {
	Minute     string
	Datacenter string
}

// NorthboundService orchestrates Eastmoney northbound providers.
type NorthboundService struct {
	client NorthboundClient
	urls   NorthboundURLs
}

// NewNorthboundService creates a NorthboundService.
func NewNorthboundService(client NorthboundClient, urls NorthboundURLs) *NorthboundService {
	return &NorthboundService{client: client, urls: urls}
}

// Minute returns northbound or southbound intraday flow rows.
func (s *NorthboundService) Minute(ctx context.Context, direction eastmoney.NorthboundDirection) ([]types.NorthboundMinuteItem, error) {
	return eastmoney.GetNorthboundMinute(ctx, s.client, s.urls.Minute, direction)
}

// Summary returns Shanghai/Shenzhen/HK connect flow summary rows.
func (s *NorthboundService) Summary(ctx context.Context) ([]types.NorthboundFlowSummary, error) {
	return eastmoney.GetNorthboundFlowSummary(ctx, s.client, s.urls.Datacenter)
}

// HoldingRank returns northbound holding ranking rows.
func (s *NorthboundService) HoldingRank(ctx context.Context, options eastmoney.NorthboundHoldingRankOptions) ([]types.NorthboundHoldingRankItem, error) {
	return eastmoney.GetNorthboundHoldingRank(ctx, s.client, s.urls.Datacenter, options)
}

// History returns northbound or southbound daily flow history rows.
func (s *NorthboundService) History(ctx context.Context, direction eastmoney.NorthboundDirection, options eastmoney.NorthboundHistoryOptions) ([]types.NorthboundHistoryItem, error) {
	return eastmoney.GetNorthboundHistory(ctx, s.client, s.urls.Datacenter, direction, options)
}

// Individual returns a stock's northbound holding history rows.
func (s *NorthboundService) Individual(ctx context.Context, symbol string, options eastmoney.NorthboundHistoryOptions) ([]types.NorthboundIndividualItem, error) {
	return eastmoney.GetNorthboundIndividual(ctx, s.client, s.urls.Datacenter, symbol, options)
}
