package services

import (
	"context"
	"sync"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/internal/providers/tencent"
	"github.com/ceheng.io/stock-go/types"
)

// QuoteClient is the request client interface required by QuoteService.
type QuoteClient interface {
	tencent.QuoteClient
	tencent.SearchClient
	tencent.CalendarClient
	tencent.CodeListClient
	tencent.TimelineClient
	eastmoney.DividendClient
}

// QuoteURLs contains quote-related endpoints.
type QuoteURLs struct {
	Minute     string
	Datacenter string
}

// QuoteService orchestrates quote providers.
type QuoteService struct {
	client QuoteClient
	urls   QuoteURLs
}

var tradingCalendarCache = newCalendarCache()

type calendarCacheState struct {
	mu      sync.Mutex
	entries map[string]*calendarCacheEntry
}

type calendarCacheEntry struct {
	cond     *sync.Cond
	values   []string
	fetching bool
}

func newCalendarCache() *calendarCacheState {
	return &calendarCacheState{entries: make(map[string]*calendarCacheEntry)}
}

// NewQuoteService creates a QuoteService.
func NewQuoteService(client QuoteClient, urls ...QuoteURLs) *QuoteService {
	config := QuoteURLs{}
	if len(urls) > 0 {
		config = urls[0]
	}
	if config.Minute == "" {
		config.Minute = "https://web.ifzq.gtimg.cn/appstock/app/minute/query"
	}
	if config.Datacenter == "" {
		config.Datacenter = "https://datacenter-web.eastmoney.com/api/data/v1/get"
	}
	return &QuoteService{client: client, urls: config}
}

// SimpleCN returns CN simple quotes.
func (s *QuoteService) SimpleCN(ctx context.Context, codes []string) ([]types.SimpleQuote, error) {
	return tencent.GetSimpleQuotes(ctx, s.client, codes)
}

// CN returns detailed CN quotes.
func (s *QuoteService) CN(ctx context.Context, codes []string) ([]types.FullQuote, error) {
	return tencent.GetFullQuotes(ctx, s.client, codes)
}

// HK returns Hong Kong quotes.
func (s *QuoteService) HK(ctx context.Context, codes []string) ([]types.HKQuote, error) {
	return tencent.GetHKQuotes(ctx, s.client, codes)
}

// US returns US quotes.
func (s *QuoteService) US(ctx context.Context, codes []string) ([]types.USQuote, error) {
	return tencent.GetUSQuotes(ctx, s.client, codes)
}

// Fund returns public fund quotes.
func (s *QuoteService) Fund(ctx context.Context, codes []string) ([]types.FundQuote, error) {
	return tencent.GetFundQuotes(ctx, s.client, codes)
}

// FundFlow returns Tencent fund-flow rows.
func (s *QuoteService) FundFlow(ctx context.Context, codes []string) ([]types.FundFlow, error) {
	return tencent.GetFundFlow(ctx, s.client, codes)
}

// PanelLargeOrder returns Tencent panel large-order ratio rows.
func (s *QuoteService) PanelLargeOrder(ctx context.Context, codes []string) ([]types.PanelLargeOrder, error) {
	return tencent.GetPanelLargeOrder(ctx, s.client, codes)
}

// TodayTimeline returns Tencent intraday timeline rows for one code.
func (s *QuoteService) TodayTimeline(ctx context.Context, code string) (types.TodayTimelineResponse, error) {
	return tencent.GetTodayTimeline(ctx, s.client, s.urls.Minute, code)
}

// Search returns Tencent Smartbox search results.
func (s *QuoteService) Search(ctx context.Context, keyword string) ([]types.SearchResult, error) {
	return tencent.Search(ctx, s.client, keyword)
}

// TradingCalendar returns A-share trading dates.
func (s *QuoteService) TradingCalendar(ctx context.Context) ([]string, error) {
	tradingCalendarCache.mu.Lock()
	cacheKey := s.client.CalendarURL()
	entry := tradingCalendarCache.entries[cacheKey]
	if entry == nil {
		entry = &calendarCacheEntry{}
		entry.cond = sync.NewCond(&tradingCalendarCache.mu)
		tradingCalendarCache.entries[cacheKey] = entry
	}
	if entry.values != nil {
		cached := append([]string(nil), entry.values...)
		tradingCalendarCache.mu.Unlock()
		return cached, nil
	}
	for entry.fetching {
		entry.cond.Wait()
		if entry.values != nil {
			cached := append([]string(nil), entry.values...)
			tradingCalendarCache.mu.Unlock()
			return cached, nil
		}
	}
	entry.fetching = true
	tradingCalendarCache.mu.Unlock()

	calendar, err := tencent.GetTradingCalendar(ctx, s.client)

	tradingCalendarCache.mu.Lock()
	entry.fetching = false
	if err == nil && entry.values == nil {
		entry.values = append([]string(nil), calendar...)
	}
	entry.cond.Broadcast()
	if err != nil {
		tradingCalendarCache.mu.Unlock()
		return nil, err
	}
	cached := append([]string(nil), entry.values...)
	tradingCalendarCache.mu.Unlock()
	return cached, nil
}

// CodesCN returns A-share codes.
func (s *QuoteService) CodesCN(ctx context.Context, options tencent.CodeListOptions) ([]string, error) {
	return tencent.GetAShareCodeList(ctx, s.client, options)
}

// CodesUS returns US stock codes.
func (s *QuoteService) CodesUS(ctx context.Context, options tencent.USCodeListOptions) ([]string, error) {
	return tencent.GetUSCodeList(ctx, s.client, options)
}

// CodesHK returns HK stock codes.
func (s *QuoteService) CodesHK(ctx context.Context) ([]string, error) {
	return tencent.GetHKCodeList(ctx, s.client)
}

// CodesFund returns fund codes.
func (s *QuoteService) CodesFund(ctx context.Context) ([]string, error) {
	return tencent.GetFundCodeList(ctx, s.client)
}

// BatchCN returns detailed CN quotes by codes with batching.
func (s *QuoteService) BatchCN(ctx context.Context, codes []string, options tencent.BatchOptions) ([]types.FullQuote, error) {
	return tencent.GetAllQuotesByCodes(ctx, s.client, codes, options)
}

// BatchHK returns HK quotes by codes with batching.
func (s *QuoteService) BatchHK(ctx context.Context, codes []string, options tencent.BatchOptions) ([]types.HKQuote, error) {
	return tencent.GetAllHKQuotesByCodes(ctx, s.client, codes, options)
}

// BatchUS returns US quotes by codes with batching.
func (s *QuoteService) BatchUS(ctx context.Context, codes []string, options tencent.BatchOptions) ([]types.USQuote, error) {
	return tencent.GetAllUSQuotesByCodes(ctx, s.client, codes, options)
}

// AllCN returns detailed CN quotes for all A-share codes matching code-list options.
func (s *QuoteService) AllCN(ctx context.Context, codeOptions tencent.CodeListOptions, batchOptions tencent.BatchOptions) ([]types.FullQuote, error) {
	codes, err := s.CodesCN(ctx, codeOptions)
	if err != nil {
		return nil, err
	}
	return s.BatchCN(ctx, codes, batchOptions)
}

// AllHK returns HK quotes for all HK stock codes.
func (s *QuoteService) AllHK(ctx context.Context, batchOptions tencent.BatchOptions) ([]types.HKQuote, error) {
	codes, err := s.CodesHK(ctx)
	if err != nil {
		return nil, err
	}
	return s.BatchHK(ctx, codes, batchOptions)
}

// AllUS returns US quotes for all US stock codes matching code-list options.
func (s *QuoteService) AllUS(ctx context.Context, codeOptions tencent.USCodeListOptions, batchOptions tencent.BatchOptions) ([]types.USQuote, error) {
	codeOptions.Simple = true
	codes, err := s.CodesUS(ctx, codeOptions)
	if err != nil {
		return nil, err
	}
	return s.BatchUS(ctx, codes, batchOptions)
}

// BatchRaw returns raw Tencent quote assignments for custom query params.
func (s *QuoteService) BatchRaw(ctx context.Context, params string) ([]types.TencentQuoteItem, error) {
	items, err := s.client.GetTencentQuote(ctx, params)
	if err != nil {
		return nil, err
	}
	result := make([]types.TencentQuoteItem, len(items))
	for i, item := range items {
		result[i] = types.TencentQuoteItem{
			Key:    item.Key,
			Fields: append([]string(nil), item.Fields...),
		}
	}
	return result, nil
}

// DividendDetail returns stock dividend detail rows.
func (s *QuoteService) DividendDetail(ctx context.Context, symbol string) ([]types.DividendDetail, error) {
	return eastmoney.GetDividendDetail(ctx, s.client, s.urls.Datacenter, symbol)
}
