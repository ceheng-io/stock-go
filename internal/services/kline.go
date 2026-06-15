package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// KlineClient is the request client interface required by KlineService.
type KlineClient interface {
	eastmoney.KlineClient
}

// KlineURLs contains Eastmoney K-line endpoints per market.
type KlineURLs struct {
	CN       string
	CNTrends string
	HK       string
	HKTrends string
	US       string
	USTrends string
}

// KlineService orchestrates K-line providers.
type KlineService struct {
	client KlineClient
	urls   KlineURLs
}

// NewKlineService creates a KlineService.
func NewKlineService(client KlineClient, urls KlineURLs) *KlineService {
	if urls.CNTrends == "" {
		urls.CNTrends = urls.CN
	}
	if urls.HKTrends == "" {
		urls.HKTrends = urls.HK
	}
	if urls.USTrends == "" {
		urls.USTrends = urls.US
	}
	return &KlineService{client: client, urls: urls}
}

// CN returns CN historical K-line rows.
func (s *KlineService) CN(ctx context.Context, symbol string, options eastmoney.HistoryKlineOptions) ([]types.HistoryKline, error) {
	return eastmoney.GetHistoryKline(ctx, s.client, symbol, s.urls.CN, options)
}

// CNMinute returns CN minute timeline or K-line rows.
func (s *KlineService) CNMinute(ctx context.Context, symbol string, options eastmoney.MinuteKlineOptions) (types.MinuteKlineResult, error) {
	return eastmoney.GetMinuteKline(ctx, s.client, symbol, s.urls.CN, s.urls.CNTrends, options)
}

// HK returns HK historical K-line rows.
func (s *KlineService) HK(ctx context.Context, symbol string, options eastmoney.HistoryKlineOptions) ([]types.HKHistoryKline, error) {
	return eastmoney.GetHKHistoryKline(ctx, s.client, symbol, s.urls.HK, options)
}

// US returns US historical K-line rows.
func (s *KlineService) US(ctx context.Context, symbol string, options eastmoney.HistoryKlineOptions) ([]types.USHistoryKline, error) {
	return eastmoney.GetUSHistoryKline(ctx, s.client, symbol, s.urls.US, options)
}

// HKMinute returns HK minute timeline or K-line rows.
func (s *KlineService) HKMinute(ctx context.Context, symbol string, options eastmoney.MinuteKlineOptions) (types.HKMinuteKlineResult, error) {
	return eastmoney.GetHKMinuteKline(ctx, s.client, symbol, s.urls.HK, s.urls.HKTrends, options)
}

// USMinute returns US minute timeline or K-line rows.
func (s *KlineService) USMinute(ctx context.Context, symbol string, options eastmoney.MinuteKlineOptions) (types.USMinuteKlineResult, error) {
	return eastmoney.GetUSMinuteKline(ctx, s.client, symbol, s.urls.US, s.urls.USTrends, options)
}
