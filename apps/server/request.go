package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	stock "github.com/ceheng.io/stock-go"
)

type SDK interface {
	Search(ctx context.Context, keyword string) ([]stock.SearchResult, error)
	GetFullQuotes(ctx context.Context, codes []string) ([]stock.FullQuote, error)
	GetAllQuotesByCodes(ctx context.Context, codes []string, options ...stock.GetAllAShareQuotesOptions) ([]stock.FullQuote, error)
	GetAllAShareQuotes(ctx context.Context, options ...stock.GetAllAShareQuotesOptions) ([]stock.FullQuote, error)
	GetFundFlow(ctx context.Context, codes []string) ([]stock.FundFlow, error)
	GetPanelLargeOrder(ctx context.Context, codes []string) ([]stock.PanelLargeOrder, error)
	GetTodayTimeline(ctx context.Context, code string) (stock.TodayTimelineResponse, error)
	GetHistoryKline(ctx context.Context, symbol string, options ...stock.HistoryKlineOptions) ([]stock.HistoryKline, error)
	GetMinuteKline(ctx context.Context, symbol string, options ...stock.MinuteKlineOptions) (stock.MinuteKlineResult, error)
	GetKlineWithIndicators(ctx context.Context, symbol string, options ...stock.KlineWithIndicatorsOptions) ([]stock.KlineWithIndicators, error)
	GetIndustryList(ctx context.Context) ([]stock.Board, error)
	GetConceptList(ctx context.Context) ([]stock.Board, error)
	GetIndustrySpot(ctx context.Context, boardCode string) ([]stock.BoardSpot, error)
	GetConceptSpot(ctx context.Context, boardCode string) ([]stock.BoardSpot, error)
	GetIndustryConstituents(ctx context.Context, boardCode string) ([]stock.BoardConstituent, error)
	GetConceptConstituents(ctx context.Context, boardCode string) ([]stock.BoardConstituent, error)
	GetIndustryKline(ctx context.Context, boardCode string, options ...stock.IndustryBoardKlineOptions) ([]stock.BoardKline, error)
	GetConceptKline(ctx context.Context, boardCode string, options ...stock.ConceptBoardKlineOptions) ([]stock.BoardKline, error)
	GetIndustryMinuteKline(ctx context.Context, boardCode string, options ...stock.IndustryBoardMinuteKlineOptions) (stock.BoardMinuteKlineResult, error)
	GetConceptMinuteKline(ctx context.Context, boardCode string, options ...stock.ConceptBoardMinuteKlineOptions) (stock.BoardMinuteKlineResult, error)
	GetIndividualFundFlow(ctx context.Context, symbol string, options ...stock.FundFlowOptions) ([]stock.StockFundFlow, error)
	GetMarketFundFlow(ctx context.Context) ([]stock.MarketFundFlow, error)
	GetFundFlowRank(ctx context.Context, options ...stock.FundFlowRankOptions) ([]stock.FundFlowRankItem, error)
	GetSectorFundFlowRank(ctx context.Context, options ...stock.FundFlowRankOptions) ([]stock.SectorFundFlowItem, error)
	GetSectorFundFlowHistory(ctx context.Context, symbol string, options ...stock.FundFlowOptions) ([]stock.StockFundFlow, error)
	GetNorthboundMinute(ctx context.Context, directions ...stock.NorthboundDirection) ([]stock.NorthboundMinuteItem, error)
	GetNorthboundFlowSummary(ctx context.Context) ([]stock.NorthboundFlowSummary, error)
	GetNorthboundHoldingRank(ctx context.Context, options ...stock.NorthboundHoldingRankOptions) ([]stock.NorthboundHoldingRankItem, error)
	GetNorthboundHistory(ctx context.Context, args ...any) ([]stock.NorthboundHistoryItem, error)
	GetNorthboundIndividual(ctx context.Context, symbol string, options ...stock.NorthboundHistoryOptions) ([]stock.NorthboundIndividualItem, error)
	GetZTPool(ctx context.Context, args ...any) ([]stock.ZTPoolItem, error)
	GetStockChanges(ctx context.Context, changeTypes ...stock.StockChangeType) ([]stock.StockChangeItem, error)
	GetBoardChanges(ctx context.Context) ([]stock.BoardChangeItem, error)
	GetDragonTigerDetail(ctx context.Context, options stock.DragonTigerDateOptions) ([]stock.DragonTigerDetailItem, error)
	GetBlockTradeDetail(ctx context.Context, options ...stock.BlockTradeDateOptions) ([]stock.BlockTradeDetailItem, error)
	GetMarginAccountInfo(ctx context.Context) ([]stock.MarginAccountItem, error)
	GetDividendDetail(ctx context.Context, symbol string) ([]stock.DividendDetail, error)
	GetTradingCalendar(ctx context.Context) ([]string, error)
}

func requiredQuery(r *http.Request, key string) (string, bool) {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	return value, value != ""
}

func queryString(r *http.Request, key string, fallback string) string {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return fallback
	}
	return value
}

func queryInt(r *http.Request, key string, fallback int) int {
	value := strings.TrimSpace(r.URL.Query().Get(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func queryCodes(r *http.Request) ([]string, bool) {
	raw, ok := requiredQuery(r, "codes")
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	codes := make([]string, 0, len(parts))
	for _, part := range parts {
		code := strings.TrimSpace(part)
		if code != "" {
			codes = append(codes, code)
		}
	}
	return codes, len(codes) > 0
}

func asBatchOptions(r *http.Request) stock.GetAllAShareQuotesOptions {
	return stock.GetAllAShareQuotesOptions{
		BatchSize:   queryInt(r, "batchSize", 0),
		Concurrency: queryInt(r, "concurrency", 0),
	}
}

func asHistoryOptions(r *http.Request) stock.HistoryKlineOptions {
	return stock.HistoryKlineOptions{
		Period:    stock.KlinePeriod(queryString(r, "period", "")),
		Adjust:    stock.AdjustType(queryString(r, "adjust", "")),
		StartDate: queryString(r, "startDate", ""),
		EndDate:   queryString(r, "endDate", ""),
	}
}

func asMinuteOptions(r *http.Request) stock.MinuteKlineOptions {
	return stock.MinuteKlineOptions{
		Period:    stock.MinutePeriod(queryString(r, "period", "")),
		Adjust:    stock.AdjustType(queryString(r, "adjust", "")),
		StartDate: queryString(r, "startDate", ""),
		EndDate:   queryString(r, "endDate", ""),
		NDays:     queryInt(r, "nDays", 0),
	}
}

func asFundFlowOptions(r *http.Request) stock.FundFlowOptions {
	return stock.FundFlowOptions{
		Period: stock.FundFlowPeriod(queryString(r, "period", "")),
	}
}

func asFundFlowRankOptions(r *http.Request) stock.FundFlowRankOptions {
	return stock.FundFlowRankOptions{
		Indicator:  stock.FundFlowRankIndicator(queryString(r, "indicator", "")),
		SectorType: stock.FundFlowSectorType(queryString(r, "sectorType", "")),
	}
}

func asNorthboundHistoryOptions(r *http.Request) stock.NorthboundHistoryOptions {
	return stock.NorthboundHistoryOptions{
		StartDate: queryString(r, "startDate", ""),
		EndDate:   queryString(r, "endDate", ""),
	}
}
