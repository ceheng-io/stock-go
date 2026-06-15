package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

// MarketEventClient is the request client interface required by MarketEventService.
type MarketEventClient interface {
	eastmoney.TopicDataClient
}

// MarketEventService orchestrates Eastmoney market-event providers.
type MarketEventService struct {
	client   MarketEventClient
	topicURL string
}

// NewMarketEventService creates a MarketEventService.
func NewMarketEventService(client MarketEventClient, topicURL string) *MarketEventService {
	return &MarketEventService{client: client, topicURL: topicURL}
}

// ZTPool returns limit-up pool rows.
func (s *MarketEventService) ZTPool(ctx context.Context, poolType eastmoney.ZTPoolType, date string) ([]types.ZTPoolItem, error) {
	return eastmoney.GetZTPool(ctx, s.client, s.topicURL, poolType, date)
}

// StockChanges returns intraday stock change rows.
func (s *MarketEventService) StockChanges(ctx context.Context, changeType eastmoney.StockChangeType) ([]types.StockChangeItem, error) {
	return eastmoney.GetStockChanges(ctx, s.client, s.topicURL, changeType)
}

// BoardChanges returns board change rows.
func (s *MarketEventService) BoardChanges(ctx context.Context) ([]types.BoardChangeItem, error) {
	return eastmoney.GetBoardChanges(ctx, s.client, s.topicURL)
}
