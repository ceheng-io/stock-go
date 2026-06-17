package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/internal/providers/ths"
	"github.com/ceheng-io/stock-go/types"
)

// MarketEventClient is the request client interface required by MarketEventService.
type MarketEventClient interface {
	eastmoney.TopicDataClient
	ths.LimitUpClient
}

// MarketEventService orchestrates Eastmoney market-event providers.
type MarketEventService struct {
	client   MarketEventClient
	topicURL string
	thsURL   string
}

// NewMarketEventService creates a MarketEventService.
func NewMarketEventService(client MarketEventClient, topicURL string, thsURL ...string) *MarketEventService {
	service := &MarketEventService{client: client, topicURL: topicURL}
	if len(thsURL) > 0 {
		service.thsURL = thsURL[0]
	}
	return service
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

// THSLimitUpPool returns Tonghuashun limit-up pool rows.
func (s *MarketEventService) THSLimitUpPool(ctx context.Context, options types.THSLimitUpPoolOptions) (types.THSLimitUpPoolResult, error) {
	return ths.GetLimitUpPool(ctx, s.client, s.thsURL, options)
}
