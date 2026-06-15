package main

import (
	"context"

	stock "github.com/ceheng.io/stock-go"
)

func (f *fakeSDK) GetAllQuotesByCodes(context.Context, []string, ...stock.GetAllAShareQuotesOptions) ([]stock.FullQuote, error) {
	return nil, nil
}

func (f *fakeSDK) GetAllAShareQuotes(context.Context, ...stock.GetAllAShareQuotesOptions) ([]stock.FullQuote, error) {
	return nil, nil
}

func (f *fakeSDK) GetFundFlow(context.Context, []string) ([]stock.FundFlow, error) {
	return nil, nil
}

func (f *fakeSDK) GetPanelLargeOrder(context.Context, []string) ([]stock.PanelLargeOrder, error) {
	return nil, nil
}

func (f *fakeSDK) GetTodayTimeline(context.Context, string) (stock.TodayTimelineResponse, error) {
	return stock.TodayTimelineResponse{}, nil
}

func (f *fakeSDK) GetMinuteKline(context.Context, string, ...stock.MinuteKlineOptions) (stock.MinuteKlineResult, error) {
	return stock.MinuteKlineResult{}, nil
}

func (f *fakeSDK) GetConceptList(context.Context) ([]stock.Board, error) {
	return nil, nil
}

func (f *fakeSDK) GetIndustrySpot(context.Context, string) ([]stock.BoardSpot, error) {
	return nil, nil
}

func (f *fakeSDK) GetConceptSpot(context.Context, string) ([]stock.BoardSpot, error) {
	return nil, nil
}

func (f *fakeSDK) GetConceptConstituents(context.Context, string) ([]stock.BoardConstituent, error) {
	return nil, nil
}

func (f *fakeSDK) GetIndustryKline(context.Context, string, ...stock.IndustryBoardKlineOptions) ([]stock.BoardKline, error) {
	return nil, nil
}

func (f *fakeSDK) GetConceptKline(context.Context, string, ...stock.ConceptBoardKlineOptions) ([]stock.BoardKline, error) {
	return nil, nil
}

func (f *fakeSDK) GetIndustryMinuteKline(context.Context, string, ...stock.IndustryBoardMinuteKlineOptions) (stock.BoardMinuteKlineResult, error) {
	return stock.BoardMinuteKlineResult{}, nil
}

func (f *fakeSDK) GetConceptMinuteKline(context.Context, string, ...stock.ConceptBoardMinuteKlineOptions) (stock.BoardMinuteKlineResult, error) {
	return stock.BoardMinuteKlineResult{}, nil
}

func (f *fakeSDK) GetIndividualFundFlow(context.Context, string, ...stock.FundFlowOptions) ([]stock.StockFundFlow, error) {
	return nil, nil
}

func (f *fakeSDK) GetMarketFundFlow(context.Context) ([]stock.MarketFundFlow, error) {
	return nil, nil
}

func (f *fakeSDK) GetFundFlowRank(context.Context, ...stock.FundFlowRankOptions) ([]stock.FundFlowRankItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetSectorFundFlowRank(context.Context, ...stock.FundFlowRankOptions) ([]stock.SectorFundFlowItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetNorthboundMinute(context.Context, ...stock.NorthboundDirection) ([]stock.NorthboundMinuteItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetNorthboundFlowSummary(context.Context) ([]stock.NorthboundFlowSummary, error) {
	return nil, nil
}

func (f *fakeSDK) GetNorthboundHoldingRank(context.Context, ...stock.NorthboundHoldingRankOptions) ([]stock.NorthboundHoldingRankItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetNorthboundHistory(context.Context, ...any) ([]stock.NorthboundHistoryItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetNorthboundIndividual(context.Context, string, ...stock.NorthboundHistoryOptions) ([]stock.NorthboundIndividualItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetStockChanges(context.Context, ...stock.StockChangeType) ([]stock.StockChangeItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetBoardChanges(context.Context) ([]stock.BoardChangeItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetDragonTigerDetail(context.Context, stock.DragonTigerDateOptions) ([]stock.DragonTigerDetailItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetBlockTradeDetail(context.Context, ...stock.BlockTradeDateOptions) ([]stock.BlockTradeDetailItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetMarginAccountInfo(context.Context) ([]stock.MarginAccountItem, error) {
	return nil, nil
}

func (f *fakeSDK) GetDividendDetail(context.Context, string) ([]stock.DividendDetail, error) {
	return nil, nil
}

func (f *fakeSDK) GetTradingCalendar(context.Context) ([]string, error) {
	return nil, nil
}
