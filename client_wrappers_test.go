package stock

import (
	"context"
	"testing"
	"time"
)

func TestClientReExportsFlatQuoteWrappers(t *testing.T) {
	client := New()
	ctx := context.Background()

	_ = client.GetFullQuotes
	_ = client.GetSimpleQuotes
	_ = client.GetHKQuotes
	_ = client.GetUSQuotes
	_ = client.GetFundQuotes
	_ = client.GetFundFlow
	_ = client.GetPanelLargeOrder
	_ = client.GetTodayTimeline
	_ = client.Search
	_ = client.GetAShareCodeList
	_ = client.GetUSCodeList
	_ = client.GetHKCodeList
	_ = client.GetFundCodeList
	_ = client.GetAllAShareQuotes
	_ = client.GetAllHKShareQuotes
	_ = client.GetAllUSShareQuotes
	_ = client.GetAllQuotesByCodes
	_ = client.BatchRaw
	_ = client.GetTradingCalendar
	_ = client.IsTradingDay
	_ = client.NextTradingDay
	_ = client.PrevTradingDay
	_ = client.GetMarketStatus
	_ = client.GetDividendDetail

	status := client.GetMarketStatus(MarketA, time.Date(2024, 6, 13, 10, 0, 0, 0, time.UTC))
	if status == "" {
		t.Fatal("GetMarketStatus returned empty status")
	}
	_, _ = client.GetHKCodeList(ctx)
}

func TestClientReExportsFlatMarketDataWrappers(t *testing.T) {
	client := New()

	_ = client.GetIndustryList
	_ = client.GetIndustrySpot
	_ = client.GetIndustryConstituents
	_ = client.GetIndustryKline
	_ = client.GetIndustryMinuteKline
	_ = client.GetConceptList
	_ = client.GetConceptSpot
	_ = client.GetConceptConstituents
	_ = client.GetConceptKline
	_ = client.GetConceptMinuteKline
	_ = client.GetHistoryKline
	_ = client.GetMinuteKline
	_ = client.GetHKHistoryKline
	_ = client.GetHKMinuteKline
	_ = client.GetUSHistoryKline
	_ = client.GetUSMinuteKline
	_ = client.GetKlineWithIndicators
}

func TestClientReExportsFlatExtendedDataWrappers(t *testing.T) {
	client := New()

	_ = client.GetFuturesKline
	_ = client.GetGlobalFuturesSpot
	_ = client.GetGlobalFuturesKline
	_ = client.GetFuturesInventorySymbols
	_ = client.GetFuturesInventory
	_ = client.GetComexInventory
	_ = client.GetIndexOptionSpot
	_ = client.GetIndexOptionKline
	_ = client.GetCFFEXOptionQuotes
	_ = client.GetETFOptionMonths
	_ = client.GetETFOptionExpireDay
	_ = client.GetETFOptionMinute
	_ = client.GetETFOptionDailyKline
	_ = client.GetETFOption5DayMinute
	_ = client.GetCommodityOptionSpot
	_ = client.GetCommodityOptionKline
	_ = client.GetOptionLHB
	_ = client.GetIndividualFundFlow
	_ = client.GetMarketFundFlow
	_ = client.GetFundFlowRank
	_ = client.GetSectorFundFlowRank
	_ = client.GetSectorFundFlowHistory
	_ = client.GetNorthboundMinute
	_ = client.GetNorthboundFlowSummary
	_ = client.GetNorthboundHoldingRank
	_ = client.GetNorthboundHistory
	_ = client.GetNorthboundIndividual
	_ = client.GetZTPool
	_ = client.GetStockChanges
	_ = client.GetBoardChanges
	_ = client.GetDragonTigerDetail
	_ = client.GetDragonTigerStockStats
	_ = client.GetDragonTigerInstitution
	_ = client.GetDragonTigerBranchRank
	_ = client.GetDragonTigerStockSeatDetail
	_ = client.GetBlockTradeMarketStat
	_ = client.GetBlockTradeDetail
	_ = client.GetBlockTradeDailyStat
	_ = client.GetMarginAccountInfo
	_ = client.GetMarginTargetList
	_ = client.GetFundDividendList
	_ = client.GetFundNavHistory
	_ = client.GetFundEstimate
	_ = client.GetFundRankHistory
}

func compileClientOptionalWrapperCalls(client *Client, ctx context.Context) {
	_, _ = client.GetAShareCodeList(ctx)
	_, _ = client.GetUSCodeList(ctx)
	_, _ = client.GetAllAShareQuotes(ctx)
	_, _ = client.GetAllHKShareQuotes(ctx)
	_, _ = client.GetAllUSShareQuotes(ctx)
	_, _ = client.GetAllQuotesByCodes(ctx, []string{"sh600519"})
	_, _ = client.GetIndustryKline(ctx, "BK0475")
	_, _ = client.GetIndustryMinuteKline(ctx, "BK0475")
	_, _ = client.GetConceptKline(ctx, "BK0898")
	_, _ = client.GetConceptMinuteKline(ctx, "BK0898")
	_, _ = client.GetHistoryKline(ctx, "600519")
	_, _ = client.GetMinuteKline(ctx, "600519")
	_, _ = client.GetHKHistoryKline(ctx, "00700")
	_, _ = client.GetHKMinuteKline(ctx, "00700")
	_, _ = client.GetUSHistoryKline(ctx, "105.AAPL")
	_, _ = client.GetUSMinuteKline(ctx, "105.AAPL")
	_, _ = client.GetKlineWithIndicators(ctx, "600519")
	_, _ = client.GetFuturesKline(ctx, "rb2605")
	_, _ = client.GetGlobalFuturesSpot(ctx)
	_, _ = client.GetGlobalFuturesKline(ctx, "GC00Y")
	_, _ = client.GetFuturesInventory(ctx, "RB")
	_, _ = client.GetComexInventory(ctx, "gold")
	_, _ = client.GetCFFEXOptionQuotes(ctx)
	_, _ = client.GetIndividualFundFlow(ctx, "600519")
	_, _ = client.GetFundFlowRank(ctx)
	_, _ = client.GetSectorFundFlowRank(ctx)
	_, _ = client.GetSectorFundFlowHistory(ctx, "BK0475")
	_, _ = client.GetNorthboundMinute(ctx)
	_, _ = client.GetNorthboundHoldingRank(ctx)
	_, _ = client.GetNorthboundHistory(ctx)
	_, _ = client.GetNorthboundHistory(ctx, NorthboundSouth)
	_, _ = client.GetNorthboundHistory(ctx, "south")
	_, _ = client.GetNorthboundHistory(ctx, NorthboundHistoryOptions{})
	_, _ = client.GetNorthboundHistory(ctx, NorthboundSouth, NorthboundHistoryOptions{})
	_, _ = client.GetStockChanges(ctx)
	_, _ = client.GetZTPool(ctx)
	_, _ = client.GetZTPool(ctx, ZTPoolZT)
	_, _ = client.GetZTPool(ctx, "broken")
	_, _ = client.GetZTPool(ctx, "2024-12-16")
	_, _ = client.GetZTPool(ctx, ZTPoolZT, "2024-12-16")
	_, _ = client.GetDragonTigerStockStats(ctx)
	_, _ = client.GetDragonTigerBranchRank(ctx)
	_, _ = client.GetBlockTradeDetail(ctx)
	_, _ = client.GetBlockTradeDailyStat(ctx)
	_, _ = client.GetMarginTargetList(ctx)
	_, _ = client.GetFundDividendList(ctx)
}

func TestClientOptionalWrapperRejectsUnsupportedArgumentTypes(t *testing.T) {
	client := New()
	ctx := context.Background()

	if _, err := client.GetNorthboundHistory(ctx, 1); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("GetNorthboundHistory invalid argument error = %v, code = %s", err, GetErrorCode(err))
	}
	if _, err := client.GetNorthboundHistory(ctx, "east"); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("GetNorthboundHistory invalid direction error = %v, code = %s", err, GetErrorCode(err))
	}
	if _, err := client.GetZTPool(ctx, 1); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("GetZTPool invalid argument error = %v, code = %s", err, GetErrorCode(err))
	}
	if _, err := client.GetZTPool(ctx, "unknown"); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("GetZTPool invalid string argument error = %v, code = %s", err, GetErrorCode(err))
	}
}

func TestClientOptionalWrapperArgumentParsing(t *testing.T) {
	direction, historyOptions, err := northboundHistoryArgs([]any{"south", NorthboundHistoryOptions{StartDate: "2024-01-01"}})
	if err != nil {
		t.Fatalf("northboundHistoryArgs returned error: %v", err)
	}
	if direction != NorthboundSouth {
		t.Fatalf("northbound direction = %q, want %q", direction, NorthboundSouth)
	}
	if historyOptions.StartDate != "2024-01-01" {
		t.Fatalf("northbound start date = %q", historyOptions.StartDate)
	}

	poolType, date, err := ztPoolArgs([]any{"broken", "2024-12-16"})
	if err != nil {
		t.Fatalf("ztPoolArgs returned error: %v", err)
	}
	if poolType != ZTPoolBroken {
		t.Fatalf("zt pool type = %q, want %q", poolType, ZTPoolBroken)
	}
	if date != "2024-12-16" {
		t.Fatalf("zt pool date = %q", date)
	}
}
