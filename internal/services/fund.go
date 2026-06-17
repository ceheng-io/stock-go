package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/types"
)

// FundClient is the request client interface required by FundService.
type FundClient interface {
	eastmoney.FundClient
}

// FundURLs contains Eastmoney public fund endpoints.
type FundURLs struct {
	GZ        string
	Pingzhong string
	DataIndex string
}

// FundService orchestrates Eastmoney public fund providers.
type FundService struct {
	client FundClient
	urls   FundURLs
}

// NewFundService creates a FundService.
func NewFundService(client FundClient, urls FundURLs) *FundService {
	return &FundService{client: client, urls: urls}
}

// Estimate returns the latest fund net-value estimate.
func (s *FundService) Estimate(ctx context.Context, code string) (types.FundEstimate, error) {
	return eastmoney.GetFundEstimate(ctx, s.client, s.urls.GZ, code)
}

// NavHistory returns the full net-value history.
func (s *FundService) NavHistory(ctx context.Context, code string) (types.FundNavHistory, error) {
	return eastmoney.GetFundNavHistory(ctx, s.client, s.urls.Pingzhong, code)
}

// RankHistory returns the similar-type rank history.
func (s *FundService) RankHistory(ctx context.Context, code string) (types.FundRankHistory, error) {
	return eastmoney.GetFundRankHistory(ctx, s.client, s.urls.Pingzhong, code)
}

// DividendList returns public fund dividend distribution rows.
func (s *FundService) DividendList(ctx context.Context, options eastmoney.FundDividendListOptions) (types.FundDividendListResult, error) {
	return eastmoney.GetFundDividendList(ctx, s.client, s.urls.DataIndex, options)
}
