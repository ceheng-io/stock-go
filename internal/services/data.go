package services

import (
	"context"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/internal/providers/tencent"
	"github.com/ceheng.io/stock-go/types"
)

// DataClient 是 DataService 需要的请求客户端接口。
type DataClient interface {
	tencent.SearchClient
	tencent.CodeListClient
	eastmoney.BlockTradeClient
	eastmoney.MarginClient
	eastmoney.DividendClient
}

// DataServiceOptions 包含数据类服务端点。
type DataServiceOptions struct {
	DatacenterURL string
}

// DataService 聚合数据类 provider。
type DataService struct {
	client     DataClient
	datacenter string
}

// NewDataService 创建 DataService。
func NewDataService(client DataClient, options DataServiceOptions) *DataService {
	return &DataService{client: client, datacenter: options.DatacenterURL}
}

// Search 返回腾讯 Smartbox 搜索结果。
func (s *DataService) Search(ctx context.Context, keyword string) ([]types.SearchResult, error) {
	return tencent.Search(ctx, s.client, keyword)
}

// CodesCN 返回 A 股代码列表。
func (s *DataService) CodesCN(ctx context.Context, options tencent.CodeListOptions) ([]string, error) {
	return tencent.GetAShareCodeList(ctx, s.client, options)
}

// CodesUS 返回美股代码列表。
func (s *DataService) CodesUS(ctx context.Context, options tencent.USCodeListOptions) ([]string, error) {
	return tencent.GetUSCodeList(ctx, s.client, options)
}

// CodesHK 返回港股代码列表。
func (s *DataService) CodesHK(ctx context.Context) ([]string, error) {
	return tencent.GetHKCodeList(ctx, s.client)
}

// CodesFund 返回基金代码列表。
func (s *DataService) CodesFund(ctx context.Context) ([]string, error) {
	return tencent.GetFundCodeList(ctx, s.client)
}

// BlockTradeMarketStat 返回大宗交易市场统计。
func (s *DataService) BlockTradeMarketStat(ctx context.Context) ([]types.BlockTradeMarketStatItem, error) {
	return eastmoney.GetBlockTradeMarketStat(ctx, s.client, s.datacenter)
}

// BlockTradeDetail 返回大宗交易明细。
func (s *DataService) BlockTradeDetail(ctx context.Context, options eastmoney.BlockTradeDateOptions) ([]types.BlockTradeDetailItem, error) {
	return eastmoney.GetBlockTradeDetail(ctx, s.client, s.datacenter, options)
}

// BlockTradeDailyStat 返回大宗交易每日个股统计。
func (s *DataService) BlockTradeDailyStat(ctx context.Context, options eastmoney.BlockTradeDateOptions) ([]types.BlockTradeDailyStatItem, error) {
	return eastmoney.GetBlockTradeDailyStat(ctx, s.client, s.datacenter, options)
}

// MarginAccountInfo 返回融资融券账户统计。
func (s *DataService) MarginAccountInfo(ctx context.Context) ([]types.MarginAccountItem, error) {
	return eastmoney.GetMarginAccountInfo(ctx, s.client, s.datacenter)
}

// MarginTargetList 返回融资融券标的明细。
func (s *DataService) MarginTargetList(ctx context.Context, date string) ([]types.MarginTargetItem, error) {
	return eastmoney.GetMarginTargetList(ctx, s.client, s.datacenter, date)
}

// DividendDetail 返回个股分红派送详情。
func (s *DataService) DividendDetail(ctx context.Context, symbol string) ([]types.DividendDetail, error) {
	return eastmoney.GetDividendDetail(ctx, s.client, s.datacenter, symbol)
}
