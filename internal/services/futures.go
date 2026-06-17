package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/types"
)

// FuturesClient is the request client interface required by FuturesService.
type FuturesClient interface {
	eastmoney.KlineClient
}

// FuturesService orchestrates futures providers.
type FuturesService struct {
	client FuturesClient
	urls   FuturesURLs
}

// FuturesURLs contains Eastmoney futures endpoints.
type FuturesURLs struct {
	Kline       string
	GlobalSpot  string
	GlobalKline string
	Datacenter  string
}

// NewFuturesService creates a FuturesService.
func NewFuturesService(client FuturesClient, urls FuturesURLs) *FuturesService {
	if urls.GlobalKline == "" {
		urls.GlobalKline = urls.Kline
	}
	return &FuturesService{client: client, urls: urls}
}

// Kline returns domestic futures historical K-line rows.
func (s *FuturesService) Kline(ctx context.Context, symbol string, options eastmoney.FuturesKlineOptions) ([]types.FuturesKline, error) {
	return eastmoney.GetFuturesHistoryKline(ctx, s.client, symbol, s.urls.Kline, options)
}

// GlobalSpot returns global futures spot quote rows.
func (s *FuturesService) GlobalSpot(ctx context.Context, options eastmoney.GlobalFuturesSpotOptions) ([]types.GlobalFuturesQuote, error) {
	return eastmoney.GetGlobalFuturesSpot(ctx, s.client, s.urls.GlobalSpot, options)
}

// GlobalKline returns global futures historical K-line rows.
func (s *FuturesService) GlobalKline(ctx context.Context, symbol string, options eastmoney.GlobalFuturesKlineOptions) ([]types.FuturesKline, error) {
	return eastmoney.GetGlobalFuturesKline(ctx, s.client, symbol, s.urls.GlobalKline, options)
}

// InventorySymbols 返回国内期货库存品种列表。
func (s *FuturesService) InventorySymbols(ctx context.Context) ([]types.FuturesInventorySymbol, error) {
	return eastmoney.GetFuturesInventorySymbols(ctx, s.client, s.urls.Datacenter)
}

// Inventory 返回国内期货库存数据。
func (s *FuturesService) Inventory(ctx context.Context, symbol string, options eastmoney.FuturesInventoryOptions) ([]types.FuturesInventory, error) {
	return eastmoney.GetFuturesInventory(ctx, s.client, s.urls.Datacenter, symbol, options)
}

// ComexInventory 返回 COMEX 黄金或白银库存数据。
func (s *FuturesService) ComexInventory(ctx context.Context, symbol string, options eastmoney.ComexInventoryOptions) ([]types.ComexInventory, error) {
	return eastmoney.GetComexInventory(ctx, s.client, s.urls.Datacenter, symbol, options)
}
