package services

import (
	"context"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng-io/stock-go/internal/providers/sina"
	"github.com/ceheng-io/stock-go/types"
)

// OptionsClient 是 OptionsService 需要的请求客户端接口。
type OptionsClient interface {
	eastmoney.OptionClient
	sina.JSONPClient
}

// OptionsURLs 包含期权相关端点。
type OptionsURLs struct {
	CFFEXQuotes    string
	LHB            string
	ETFMonths      string
	ETFExpire      string
	ETFMinute      string
	ETFDaily       string
	ETF5Day        string
	IndexSpot      string
	IndexKline     string
	CommoditySpot  string
	CommodityKline string
}

// OptionsService 编排期权相关 provider。
type OptionsService struct {
	client OptionsClient
	urls   OptionsURLs
}

// NewOptionsService 创建 OptionsService。
func NewOptionsService(client OptionsClient, urls OptionsURLs) *OptionsService {
	return &OptionsService{client: client, urls: urls}
}

// CFFEXQuotes 返回中金所期权实时行情。
func (s *OptionsService) CFFEXQuotes(ctx context.Context, options eastmoney.CFFEXOptionQuotesOptions) ([]types.CFFEXOptionQuote, error) {
	return eastmoney.GetCFFEXOptionQuotes(ctx, s.client, s.urls.CFFEXQuotes, options)
}

// LHB 返回期权龙虎榜数据。
func (s *OptionsService) LHB(ctx context.Context, symbol string, date string) ([]types.OptionLHBItem, error) {
	return eastmoney.GetOptionLHB(ctx, s.client, s.urls.LHB, symbol, date)
}

// ETFOptionMonths 返回新浪 ETF 期权可用月份。
func (s *OptionsService) ETFOptionMonths(ctx context.Context, cate sina.ETFOptionCate) (types.ETFOptionMonth, error) {
	return sina.GetETFOptionMonths(ctx, s.client, s.urls.ETFMonths, cate)
}

// ETFOptionExpireDay 返回新浪 ETF 期权到期日信息。
func (s *OptionsService) ETFOptionExpireDay(ctx context.Context, cate sina.ETFOptionCate, month string) (types.ETFOptionExpireDay, error) {
	return sina.GetETFOptionExpireDay(ctx, s.client, s.urls.ETFExpire, cate, month)
}

// ETFOptionMinute 返回新浪 ETF 期权当日分钟行情。
func (s *OptionsService) ETFOptionMinute(ctx context.Context, code string) ([]types.OptionMinute, error) {
	return sina.GetETFOptionMinute(ctx, s.client, s.urls.ETFMinute, code)
}

// ETFOptionDailyKline 返回新浪 ETF 期权历史日 K 线。
func (s *OptionsService) ETFOptionDailyKline(ctx context.Context, code string) ([]types.OptionKline, error) {
	return sina.GetETFOptionDailyKline(ctx, s.client, s.urls.ETFDaily, code)
}

// ETFOption5DayMinute 返回新浪 ETF 期权 5 日分钟行情。
func (s *OptionsService) ETFOption5DayMinute(ctx context.Context, code string) ([]types.OptionMinute, error) {
	return sina.GetETFOption5DayMinute(ctx, s.client, s.urls.ETF5Day, code)
}

// IndexOptionSpot 返回新浪中金所股指期权 T 型报价。
func (s *OptionsService) IndexOptionSpot(ctx context.Context, product sina.IndexOptionProduct, contract string) (types.OptionTQuoteResult, error) {
	return sina.GetIndexOptionSpot(ctx, s.client, s.urls.IndexSpot, product, contract)
}

// IndexOptionKline 返回新浪中金所股指期权合约日 K 线。
func (s *OptionsService) IndexOptionKline(ctx context.Context, symbol string) ([]types.OptionKline, error) {
	return sina.GetIndexOptionKline(ctx, s.client, s.urls.IndexKline, symbol)
}

// CommodityOptionSpot 返回新浪商品期权 T 型报价。
func (s *OptionsService) CommodityOptionSpot(ctx context.Context, variety string, contract string) (types.OptionTQuoteResult, error) {
	return sina.GetCommodityOptionSpot(ctx, s.client, s.urls.CommoditySpot, variety, contract)
}

// CommodityOptionKline 返回新浪商品期权合约日 K 线。
func (s *OptionsService) CommodityOptionKline(ctx context.Context, symbol string) ([]types.OptionKline, error) {
	return sina.GetCommodityOptionKline(ctx, s.client, s.urls.CommodityKline, symbol)
}
