package stock

import (
	"context"
	"fmt"
	"time"

	domaintypes "github.com/ceheng.io/stock-go/types"
)

// GetFullQuotes fetches detailed CN quotes.
func (c *Client) GetFullQuotes(ctx context.Context, codes []string) ([]FullQuote, error) {
	return c.Quotes.CN(ctx, codes)
}

// GetSimpleQuotes fetches compact CN quotes.
func (c *Client) GetSimpleQuotes(ctx context.Context, codes []string) ([]SimpleQuote, error) {
	return c.Quotes.SimpleCN(ctx, codes)
}

// GetHKQuotes fetches Hong Kong quotes.
func (c *Client) GetHKQuotes(ctx context.Context, codes []string) ([]HKQuote, error) {
	return c.Quotes.HK(ctx, codes)
}

// GetUSQuotes fetches US quotes.
func (c *Client) GetUSQuotes(ctx context.Context, codes []string) ([]USQuote, error) {
	return c.Quotes.US(ctx, codes)
}

// GetFundQuotes fetches public fund quotes.
func (c *Client) GetFundQuotes(ctx context.Context, codes []string) ([]FundQuote, error) {
	return c.Quotes.Fund(ctx, codes)
}

// GetFundFlow fetches Tencent stock fund-flow rows.
func (c *Client) GetFundFlow(ctx context.Context, codes []string) ([]FundFlow, error) {
	return c.Quotes.FundFlow(ctx, codes)
}

// GetPanelLargeOrder fetches panel large-order ratios.
func (c *Client) GetPanelLargeOrder(ctx context.Context, codes []string) ([]PanelLargeOrder, error) {
	return c.Quotes.PanelLargeOrder(ctx, codes)
}

// GetTodayTimeline fetches today's intraday timeline for one code.
func (c *Client) GetTodayTimeline(ctx context.Context, code string) (TodayTimelineResponse, error) {
	return c.Quotes.TodayTimeline(ctx, code)
}

// Search searches stock, index, fund, bond, futures, and option symbols.
func (c *Client) Search(ctx context.Context, keyword string) ([]SearchResult, error) {
	return c.Quotes.Search(ctx, keyword)
}

// GetAShareCodeList fetches A-share codes.
func (c *Client) GetAShareCodeList(ctx context.Context, options ...GetAShareCodeListOptions) ([]string, error) {
	return c.Quotes.CodesCN(ctx, firstOption(options))
}

// GetUSCodeList fetches US stock codes.
func (c *Client) GetUSCodeList(ctx context.Context, options ...GetUSCodeListOptions) ([]string, error) {
	return c.Quotes.CodesUS(ctx, firstOption(options))
}

// GetHKCodeList fetches Hong Kong stock codes.
func (c *Client) GetHKCodeList(ctx context.Context) ([]string, error) {
	return c.Quotes.CodesHK(ctx)
}

// GetFundCodeList fetches fund codes.
func (c *Client) GetFundCodeList(ctx context.Context) ([]string, error) {
	return c.Quotes.CodesFund(ctx)
}

// GetAllAShareQuotes fetches all A-share quotes matching the options.
func (c *Client) GetAllAShareQuotes(ctx context.Context, optionList ...GetAllAShareQuotesOptions) ([]FullQuote, error) {
	options := firstOption(optionList)
	return c.Quotes.AllCN(ctx, CodeListOptions{Market: options.Market}, BatchOptions{
		BatchSize:   options.BatchSize,
		Concurrency: options.Concurrency,
		OnProgress:  options.OnProgress,
	})
}

// GetAllHKShareQuotes fetches all Hong Kong stock quotes.
func (c *Client) GetAllHKShareQuotes(ctx context.Context, options ...GetAllHKQuotesOptions) ([]HKQuote, error) {
	return c.Quotes.AllHK(ctx, firstOption(options))
}

// GetAllUSShareQuotes fetches all US stock quotes matching the options.
func (c *Client) GetAllUSShareQuotes(ctx context.Context, optionList ...GetAllUSQuotesOptions) ([]USQuote, error) {
	options := firstOption(optionList)
	return c.Quotes.AllUS(ctx, USCodeListOptions{Market: options.Market}, BatchOptions{
		BatchSize:   options.BatchSize,
		Concurrency: options.Concurrency,
		OnProgress:  options.OnProgress,
	})
}

// GetAllQuotesByCodes fetches detailed CN quotes for explicit codes in batches.
func (c *Client) GetAllQuotesByCodes(ctx context.Context, codes []string, optionList ...GetAllAShareQuotesOptions) ([]FullQuote, error) {
	options := firstOption(optionList)
	return c.Quotes.BatchCN(ctx, codes, BatchOptions{
		BatchSize:   options.BatchSize,
		Concurrency: options.Concurrency,
		OnProgress:  options.OnProgress,
	})
}

// BatchRaw calls the Tencent batch quote endpoint with custom params.
func (c *Client) BatchRaw(ctx context.Context, params string) ([]domaintypes.TencentQuoteItem, error) {
	return c.Quotes.BatchRaw(ctx, params)
}

// GetTradingCalendar fetches A-share trading dates.
func (c *Client) GetTradingCalendar(ctx context.Context) ([]string, error) {
	return c.Quotes.TradingCalendar(ctx)
}

// IsTradingDay checks whether the date is an A-share trading day.
func (c *Client) IsTradingDay(ctx context.Context, date string) (bool, error) {
	return c.Calendar.IsTradingDay(ctx, date)
}

// NextTradingDay returns the next A-share trading day after date.
func (c *Client) NextTradingDay(ctx context.Context, date string) (string, error) {
	return c.Calendar.NextTradingDay(ctx, date)
}

// PrevTradingDay returns the previous A-share trading day before date.
func (c *Client) PrevTradingDay(ctx context.Context, date string) (string, error) {
	return c.Calendar.PrevTradingDay(ctx, date)
}

// GetMarketStatus returns the market session status at the supplied time.
func (c *Client) GetMarketStatus(market SupportedMarket, now time.Time) MarketStatus {
	return c.Calendar.MarketStatus(market, now)
}

// GetDividendDetail fetches stock dividend detail rows.
func (c *Client) GetDividendDetail(ctx context.Context, symbol string) ([]DividendDetail, error) {
	return c.Quotes.DividendDetail(ctx, symbol)
}

// GetIndustryList fetches industry boards.
func (c *Client) GetIndustryList(ctx context.Context) ([]Board, error) {
	return c.Board.IndustryList(ctx)
}

// GetConceptList fetches concept boards.
func (c *Client) GetConceptList(ctx context.Context) ([]Board, error) {
	return c.Board.ConceptList(ctx)
}

// GetIndustrySpot fetches industry board spot rows.
func (c *Client) GetIndustrySpot(ctx context.Context, boardCode string) ([]BoardSpot, error) {
	return c.Board.IndustrySpot(ctx, boardCode)
}

// GetConceptSpot fetches concept board spot rows.
func (c *Client) GetConceptSpot(ctx context.Context, boardCode string) ([]BoardSpot, error) {
	return c.Board.ConceptSpot(ctx, boardCode)
}

// GetIndustryConstituents fetches industry board constituents.
func (c *Client) GetIndustryConstituents(ctx context.Context, boardCode string) ([]BoardConstituent, error) {
	return c.Board.IndustryConstituents(ctx, boardCode)
}

// GetConceptConstituents fetches concept board constituents.
func (c *Client) GetConceptConstituents(ctx context.Context, boardCode string) ([]BoardConstituent, error) {
	return c.Board.ConceptConstituents(ctx, boardCode)
}

// GetIndustryKline fetches industry board historical K-line rows.
func (c *Client) GetIndustryKline(ctx context.Context, boardCode string, options ...IndustryBoardKlineOptions) ([]BoardKline, error) {
	return c.Board.IndustryKline(ctx, boardCode, firstOption(options))
}

// GetConceptKline fetches concept board historical K-line rows.
func (c *Client) GetConceptKline(ctx context.Context, boardCode string, options ...ConceptBoardKlineOptions) ([]BoardKline, error) {
	return c.Board.ConceptKline(ctx, boardCode, firstOption(options))
}

// GetIndustryMinuteKline fetches industry board timeline or minute K-line rows.
func (c *Client) GetIndustryMinuteKline(ctx context.Context, boardCode string, options ...IndustryBoardMinuteKlineOptions) (BoardMinuteKlineResult, error) {
	return c.Board.IndustryMinute(ctx, boardCode, firstOption(options))
}

// GetConceptMinuteKline fetches concept board timeline or minute K-line rows.
func (c *Client) GetConceptMinuteKline(ctx context.Context, boardCode string, options ...ConceptBoardMinuteKlineOptions) (BoardMinuteKlineResult, error) {
	return c.Board.ConceptMinute(ctx, boardCode, firstOption(options))
}

// GetHistoryKline fetches CN historical K-line rows.
func (c *Client) GetHistoryKline(ctx context.Context, symbol string, options ...HistoryKlineOptions) ([]HistoryKline, error) {
	return c.Kline.CN(ctx, symbol, firstOption(options))
}

// GetMinuteKline fetches CN timeline or minute K-line rows.
func (c *Client) GetMinuteKline(ctx context.Context, symbol string, options ...MinuteKlineOptions) (MinuteKlineResult, error) {
	return c.Kline.CNMinute(ctx, symbol, firstOption(options))
}

// GetHKHistoryKline fetches Hong Kong historical K-line rows.
func (c *Client) GetHKHistoryKline(ctx context.Context, symbol string, options ...HKKlineOptions) ([]HKHistoryKline, error) {
	return c.Kline.HK(ctx, symbol, firstOption(options))
}

// GetHKMinuteKline fetches Hong Kong timeline or minute K-line rows.
func (c *Client) GetHKMinuteKline(ctx context.Context, symbol string, options ...HKMinuteKlineOptions) (HKMinuteKlineResult, error) {
	return c.Kline.HKMinute(ctx, symbol, firstOption(options))
}

// GetUSHistoryKline fetches US historical K-line rows.
func (c *Client) GetUSHistoryKline(ctx context.Context, symbol string, options ...USKlineOptions) ([]USHistoryKline, error) {
	return c.Kline.US(ctx, symbol, firstOption(options))
}

// GetUSMinuteKline fetches US timeline or minute K-line rows.
func (c *Client) GetUSMinuteKline(ctx context.Context, symbol string, options ...USMinuteKlineOptions) (USMinuteKlineResult, error) {
	return c.Kline.USMinute(ctx, symbol, firstOption(options))
}

// GetKlineWithIndicators fetches K-line rows and appends indicators.
func (c *Client) GetKlineWithIndicators(ctx context.Context, symbol string, options ...KlineWithIndicatorsOptions) ([]KlineWithIndicators, error) {
	return c.Indicator.KlineWithIndicators(ctx, symbol, firstOption(options))
}

// GetFuturesKline fetches domestic futures historical K-line rows.
func (c *Client) GetFuturesKline(ctx context.Context, symbol string, options ...FuturesKlineOptions) ([]FuturesKline, error) {
	return c.Futures.Kline(ctx, symbol, firstOption(options))
}

// GetGlobalFuturesSpot fetches global futures spot quote rows.
func (c *Client) GetGlobalFuturesSpot(ctx context.Context, options ...GlobalFuturesSpotOptions) ([]GlobalFuturesQuote, error) {
	return c.Futures.GlobalSpot(ctx, firstOption(options))
}

// GetGlobalFuturesKline fetches global futures historical K-line rows.
func (c *Client) GetGlobalFuturesKline(ctx context.Context, symbol string, options ...GlobalFuturesKlineOptions) ([]FuturesKline, error) {
	return c.Futures.GlobalKline(ctx, symbol, firstOption(options))
}

// GetFuturesInventorySymbols fetches futures inventory symbols.
func (c *Client) GetFuturesInventorySymbols(ctx context.Context) ([]FuturesInventorySymbol, error) {
	return c.Futures.InventorySymbols(ctx)
}

// GetFuturesInventory fetches domestic futures inventory rows.
func (c *Client) GetFuturesInventory(ctx context.Context, symbol string, options ...FuturesInventoryOptions) ([]FuturesInventory, error) {
	return c.Futures.Inventory(ctx, symbol, firstOption(options))
}

// GetComexInventory fetches COMEX gold or silver inventory rows.
func (c *Client) GetComexInventory(ctx context.Context, symbol string, options ...ComexInventoryOptions) ([]ComexInventory, error) {
	return c.Futures.ComexInventory(ctx, symbol, firstOption(options))
}

// GetIndexOptionSpot fetches index option T-quotes.
func (c *Client) GetIndexOptionSpot(ctx context.Context, product IndexOptionProduct, contract string) (OptionTQuoteResult, error) {
	return c.Options.IndexOptionSpot(ctx, product, contract)
}

// GetIndexOptionKline fetches index option daily K-line rows.
func (c *Client) GetIndexOptionKline(ctx context.Context, symbol string) ([]OptionKline, error) {
	return c.Options.IndexOptionKline(ctx, symbol)
}

// GetCFFEXOptionQuotes fetches CFFEX option quote rows.
func (c *Client) GetCFFEXOptionQuotes(ctx context.Context, options ...CFFEXOptionQuotesOptions) ([]CFFEXOptionQuote, error) {
	return c.Options.CFFEXQuotes(ctx, firstOption(options))
}

// GetETFOptionMonths fetches available ETF option months.
func (c *Client) GetETFOptionMonths(ctx context.Context, cate ETFOptionCate) (ETFOptionMonth, error) {
	return c.Options.ETFOptionMonths(ctx, cate)
}

// GetETFOptionExpireDay fetches ETF option expiry information.
func (c *Client) GetETFOptionExpireDay(ctx context.Context, cate ETFOptionCate, month string) (ETFOptionExpireDay, error) {
	return c.Options.ETFOptionExpireDay(ctx, cate, month)
}

// GetETFOptionMinute fetches ETF option intraday minute rows.
func (c *Client) GetETFOptionMinute(ctx context.Context, code string) ([]OptionMinute, error) {
	return c.Options.ETFOptionMinute(ctx, code)
}

// GetETFOptionDailyKline fetches ETF option daily K-line rows.
func (c *Client) GetETFOptionDailyKline(ctx context.Context, code string) ([]OptionKline, error) {
	return c.Options.ETFOptionDailyKline(ctx, code)
}

// GetETFOption5DayMinute fetches ETF option five-day minute rows.
func (c *Client) GetETFOption5DayMinute(ctx context.Context, code string) ([]OptionMinute, error) {
	return c.Options.ETFOption5DayMinute(ctx, code)
}

// GetCommodityOptionSpot fetches commodity option T-quotes.
func (c *Client) GetCommodityOptionSpot(ctx context.Context, variety string, contract string) (OptionTQuoteResult, error) {
	return c.Options.CommodityOptionSpot(ctx, variety, contract)
}

// GetCommodityOptionKline fetches commodity option daily K-line rows.
func (c *Client) GetCommodityOptionKline(ctx context.Context, symbol string) ([]OptionKline, error) {
	return c.Options.CommodityOptionKline(ctx, symbol)
}

// GetOptionLHB fetches option dragon-tiger billboard rows.
func (c *Client) GetOptionLHB(ctx context.Context, symbol string, date string) ([]OptionLHBItem, error) {
	return c.Options.LHB(ctx, symbol, date)
}

// GetIndividualFundFlow fetches stock fund-flow history rows.
func (c *Client) GetIndividualFundFlow(ctx context.Context, symbol string, options ...FundFlowOptions) ([]StockFundFlow, error) {
	return c.FundFlow.Individual(ctx, symbol, firstOption(options))
}

// GetMarketFundFlow fetches Shanghai and Shenzhen market fund-flow rows.
func (c *Client) GetMarketFundFlow(ctx context.Context) ([]MarketFundFlow, error) {
	return c.FundFlow.Market(ctx)
}

// GetFundFlowRank fetches stock fund-flow ranking rows.
func (c *Client) GetFundFlowRank(ctx context.Context, options ...FundFlowRankOptions) ([]FundFlowRankItem, error) {
	return c.FundFlow.Rank(ctx, firstOption(options))
}

// GetSectorFundFlowRank fetches sector fund-flow ranking rows.
func (c *Client) GetSectorFundFlowRank(ctx context.Context, options ...FundFlowRankOptions) ([]SectorFundFlowItem, error) {
	return c.FundFlow.SectorRank(ctx, firstOption(options))
}

// GetSectorFundFlowHistory fetches sector fund-flow history rows.
func (c *Client) GetSectorFundFlowHistory(ctx context.Context, symbol string, options ...FundFlowOptions) ([]StockFundFlow, error) {
	return c.FundFlow.SectorHistory(ctx, symbol, firstOption(options))
}

// GetNorthboundMinute fetches northbound or southbound intraday flow rows.
func (c *Client) GetNorthboundMinute(ctx context.Context, directions ...NorthboundDirection) ([]NorthboundMinuteItem, error) {
	return c.Northbound.Minute(ctx, firstOption(directions))
}

// GetNorthboundFlowSummary fetches Shanghai/Shenzhen/HK connect flow summary rows.
func (c *Client) GetNorthboundFlowSummary(ctx context.Context) ([]NorthboundFlowSummary, error) {
	return c.Northbound.Summary(ctx)
}

// GetNorthboundHoldingRank fetches northbound holding ranking rows.
func (c *Client) GetNorthboundHoldingRank(ctx context.Context, options ...NorthboundHoldingRankOptions) ([]NorthboundHoldingRankItem, error) {
	return c.Northbound.HoldingRank(ctx, firstOption(options))
}

// GetNorthboundHistory fetches northbound or southbound daily flow history rows.
func (c *Client) GetNorthboundHistory(ctx context.Context, args ...any) ([]NorthboundHistoryItem, error) {
	direction, options, err := northboundHistoryArgs(args)
	if err != nil {
		return nil, err
	}
	return c.Northbound.History(ctx, direction, options)
}

// GetNorthboundIndividual fetches a stock's northbound holding history rows.
func (c *Client) GetNorthboundIndividual(ctx context.Context, symbol string, options ...NorthboundHistoryOptions) ([]NorthboundIndividualItem, error) {
	return c.Northbound.Individual(ctx, symbol, firstOption(options))
}

// GetZTPool fetches limit-up/down stock pool rows.
func (c *Client) GetZTPool(ctx context.Context, args ...any) ([]ZTPoolItem, error) {
	poolType, date, err := ztPoolArgs(args)
	if err != nil {
		return nil, err
	}
	return c.MarketEvent.ZTPool(ctx, poolType, date)
}

// GetStockChanges fetches stock intraday abnormal change rows.
func (c *Client) GetStockChanges(ctx context.Context, changeTypes ...StockChangeType) ([]StockChangeItem, error) {
	return c.MarketEvent.StockChanges(ctx, firstOption(changeTypes))
}

// GetBoardChanges fetches board intraday abnormal change rows.
func (c *Client) GetBoardChanges(ctx context.Context) ([]BoardChangeItem, error) {
	return c.MarketEvent.BoardChanges(ctx)
}

// GetTHSLimitUpPool fetches Tonghuashun limit-up pool rows.
func (c *Client) GetTHSLimitUpPool(ctx context.Context, options ...THSLimitUpPoolOptions) (THSLimitUpPoolResult, error) {
	return c.MarketEvent.THSLimitUpPool(ctx, firstOption(options))
}

// GetDragonTigerDetail fetches dragon-tiger billboard detail rows.
func (c *Client) GetDragonTigerDetail(ctx context.Context, options DragonTigerDateOptions) ([]DragonTigerDetailItem, error) {
	return c.DragonTiger.Detail(ctx, options)
}

// GetDragonTigerStockStats fetches stock dragon-tiger statistics rows.
func (c *Client) GetDragonTigerStockStats(ctx context.Context, periods ...DragonTigerPeriod) ([]DragonTigerStockStatItem, error) {
	return c.DragonTiger.StockStats(ctx, firstOption(periods))
}

// GetDragonTigerInstitution fetches institution trading rows.
func (c *Client) GetDragonTigerInstitution(ctx context.Context, options DragonTigerDateOptions) ([]DragonTigerInstitutionItem, error) {
	return c.DragonTiger.Institution(ctx, options)
}

// GetDragonTigerBranchRank fetches brokerage branch ranking rows.
func (c *Client) GetDragonTigerBranchRank(ctx context.Context, periods ...DragonTigerPeriod) ([]DragonTigerBranchItem, error) {
	return c.DragonTiger.BranchRank(ctx, firstOption(periods))
}

// GetDragonTigerStockSeatDetail fetches stock seat detail rows.
func (c *Client) GetDragonTigerStockSeatDetail(ctx context.Context, symbol string, date string) ([]DragonTigerSeatItem, error) {
	return c.DragonTiger.SeatDetail(ctx, symbol, date)
}

// GetBlockTradeMarketStat fetches block-trade market statistics rows.
func (c *Client) GetBlockTradeMarketStat(ctx context.Context) ([]BlockTradeMarketStatItem, error) {
	return c.BlockTrade.MarketStat(ctx)
}

// GetBlockTradeDetail fetches block-trade detail rows.
func (c *Client) GetBlockTradeDetail(ctx context.Context, options ...BlockTradeDateOptions) ([]BlockTradeDetailItem, error) {
	return c.BlockTrade.Detail(ctx, firstOption(options))
}

// GetBlockTradeDailyStat fetches block-trade daily statistics rows.
func (c *Client) GetBlockTradeDailyStat(ctx context.Context, options ...BlockTradeDateOptions) ([]BlockTradeDailyStatItem, error) {
	return c.BlockTrade.DailyStat(ctx, firstOption(options))
}

// GetMarginAccountInfo fetches daily margin account statistics rows.
func (c *Client) GetMarginAccountInfo(ctx context.Context) ([]MarginAccountItem, error) {
	return c.Margin.AccountInfo(ctx)
}

// GetMarginTargetList fetches margin target detail rows.
func (c *Client) GetMarginTargetList(ctx context.Context, dates ...string) ([]MarginTargetItem, error) {
	return c.Margin.TargetList(ctx, firstOption(dates))
}

// GetFundDividendList fetches public fund dividend distribution rows.
func (c *Client) GetFundDividendList(ctx context.Context, options ...FundDividendListOptions) (FundDividendListResult, error) {
	return c.Fund.DividendList(ctx, firstOption(options))
}

// GetFundNavHistory fetches a public fund's net-value history.
func (c *Client) GetFundNavHistory(ctx context.Context, code string) (FundNavHistory, error) {
	return c.Fund.NavHistory(ctx, code)
}

// GetFundEstimate fetches a public fund's latest net-value estimate.
func (c *Client) GetFundEstimate(ctx context.Context, code string) (FundEstimate, error) {
	return c.Fund.Estimate(ctx, code)
}

// GetFundRankHistory fetches a public fund's similar-type rank history.
func (c *Client) GetFundRankHistory(ctx context.Context, code string) (FundRankHistory, error) {
	return c.Fund.RankHistory(ctx, code)
}

func firstOption[T any](options []T) T {
	var zero T
	if len(options) == 0 {
		return zero
	}
	return options[0]
}

func northboundHistoryArgs(args []any) (NorthboundDirection, NorthboundHistoryOptions, error) {
	var direction NorthboundDirection
	var options NorthboundHistoryOptions
	for _, arg := range args {
		switch value := arg.(type) {
		case NorthboundDirection:
			direction = value
		case string:
			parsed := NorthboundDirection(value)
			if parsed != NorthboundNorth && parsed != NorthboundSouth {
				return "", NorthboundHistoryOptions{}, NewInvalidArgumentError(
					fmt.Sprintf("unsupported GetNorthboundHistory direction %q", value),
					map[string]any{"argument": arg},
				)
			}
			direction = parsed
		case NorthboundHistoryOptions:
			options = value
		default:
			return "", NorthboundHistoryOptions{}, NewInvalidArgumentError(
				fmt.Sprintf("unsupported GetNorthboundHistory argument type %T", arg),
				map[string]any{"argument": arg},
			)
		}
	}
	return direction, options, nil
}

func ztPoolArgs(args []any) (ZTPoolType, string, error) {
	var poolType ZTPoolType
	var date string
	for _, arg := range args {
		switch value := arg.(type) {
		case ZTPoolType:
			poolType = value
		case string:
			if parsed, ok := parseZTPoolType(value); ok {
				poolType = parsed
			} else if isZTPoolDateArg(value) {
				date = value
			} else {
				return "", "", NewInvalidArgumentError(
					fmt.Sprintf("unsupported GetZTPool string argument %q", value),
					map[string]any{"argument": arg},
				)
			}
		default:
			return "", "", NewInvalidArgumentError(
				fmt.Sprintf("unsupported GetZTPool argument type %T", arg),
				map[string]any{"argument": arg},
			)
		}
	}
	return poolType, date, nil
}

func parseZTPoolType(value string) (ZTPoolType, bool) {
	switch poolType := ZTPoolType(value); poolType {
	case ZTPoolZT, ZTPoolYesterday, ZTPoolStrong, ZTPoolSubNew, ZTPoolBroken, ZTPoolDT:
		return poolType, true
	default:
		return "", false
	}
}

func isZTPoolDateArg(value string) bool {
	if len(value) == len("20060102") {
		for _, ch := range value {
			if ch < '0' || ch > '9' {
				return false
			}
		}
		return true
	}
	if len(value) != len("2006-01-02") || value[4] != '-' || value[7] != '-' {
		return false
	}
	for i, ch := range value {
		if i == 4 || i == 7 {
			continue
		}
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
