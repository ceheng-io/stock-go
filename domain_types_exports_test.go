package stock

import "testing"

func TestRootReExportsQuoteAndKlineTypes(t *testing.T) {
	var _ SimpleQuote = SimpleQuote{Code: "600519"}
	var _ PriceLevel = PriceLevel{Price: 10, Volume: 100}
	var _ FullQuote = FullQuote{Code: "600519", Bid: []PriceLevel{{Price: 10}}}
	var _ HKQuote = HKQuote{Code: "00700", Currency: "HKD"}
	var _ USQuote = USQuote{Code: "AAPL"}
	var _ FundQuote = FundQuote{Code: "110011"}
	var _ Quote = FullQuote{Code: "600519"}
	var _ Quote = HKQuote{Code: "00700"}
	var _ Quote = USQuote{Code: "AAPL"}
	var _ Quote = FundQuote{Code: "110011"}
	var _ FundFlow = FundFlow{Code: "600519"}
	var _ PanelLargeOrder = PanelLargeOrder{BuyLargeRatio: 10}
	var _ SearchResultType = SearchStock
	var _ SearchResult = SearchResult{Code: "sh600519", Category: SearchStock}

	var _ Kline = Kline{Code: "600519"}
	var _ HistoryKline = HistoryKline{Code: "600519"}
	var _ ForeignHistoryKline = ForeignHistoryKline{Code: "AAPL"}
	var _ HKHistoryKline = HKHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "00700"}, Currency: "HKD"}
	var _ USHistoryKline = USHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "AAPL"}, Currency: "USD"}
	var _ AnyHistoryKline = HistoryKline{Code: "600519"}
	var _ AnyHistoryKline = HKHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "00700"}}
	var _ AnyHistoryKline = USHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "AAPL"}}
	var _ ForeignMinuteTimeline = ForeignMinuteTimeline{Code: "AAPL"}
	var _ ForeignMinuteKline = ForeignMinuteKline{Code: "AAPL"}
	var _ HKMinuteTimeline = HKMinuteTimeline{ForeignMinuteTimeline: ForeignMinuteTimeline{Code: "00700"}}
	var _ HKMinuteKline = HKMinuteKline{ForeignMinuteKline: ForeignMinuteKline{Code: "00700"}}
	var _ HKMinuteKlineResult = HKMinuteKlineResult{Klines: []HKMinuteKline{{Currency: "HKD"}}}
	var _ USMinuteTimeline = USMinuteTimeline{ForeignMinuteTimeline: ForeignMinuteTimeline{Code: "AAPL"}}
	var _ USMinuteKline = USMinuteKline{ForeignMinuteKline: ForeignMinuteKline{Code: "AAPL"}}
	var _ USMinuteKlineResult = USMinuteKlineResult{Klines: []USMinuteKline{{Currency: "USD"}}}
	var _ MinuteTimeline = MinuteTimeline{Code: "600519"}
	var _ MinuteKline = MinuteKline{Code: "600519"}
	var _ MinuteKlineResult = MinuteKlineResult{Klines: []MinuteKline{{Code: "600519"}}}
}

func TestRootReExportsDomainDataTypes(t *testing.T) {
	var _ ExternalLink = ExternalLink{Name: "东方财富"}
	var _ Board = Board{Code: "BK0001"}
	var _ IndustryBoard = IndustryBoard{Code: "BK0001", LeadingStock: stringPtrForDomainTest("贵州茅台")}
	var _ ConceptBoard = ConceptBoard{Code: "BK1001", LeadingStock: nil}
	var _ BoardSpot = BoardSpot{Item: "最新"}
	var _ IndustryBoardSpot = IndustryBoardSpot{Item: "最新"}
	var _ ConceptBoardSpot = ConceptBoardSpot{Item: "最新"}
	var _ BoardConstituent = BoardConstituent{Code: "600519"}
	var _ IndustryBoardConstituent = IndustryBoardConstituent{Code: "600519"}
	var _ ConceptBoardConstituent = ConceptBoardConstituent{Code: "600519"}
	var _ BoardKline = BoardKline{Date: "2024-06-13"}
	var _ IndustryBoardKline = IndustryBoardKline{Date: "2024-06-13"}
	var _ ConceptBoardKline = ConceptBoardKline{Date: "2024-06-13"}
	var _ BoardMinuteTimeline = BoardMinuteTimeline{Time: "09:30"}
	var _ IndustryBoardMinuteTimeline = IndustryBoardMinuteTimeline{Time: "09:30"}
	var _ ConceptBoardMinuteTimeline = ConceptBoardMinuteTimeline{Time: "09:30"}
	var _ BoardMinuteKline = BoardMinuteKline{Time: "09:35"}
	var _ IndustryBoardMinuteKline = IndustryBoardMinuteKline{Time: "09:35"}
	var _ ConceptBoardMinuteKline = ConceptBoardMinuteKline{Time: "09:35"}
	var _ BoardMinuteKlineResult = BoardMinuteKlineResult{Klines: []BoardMinuteKline{{Time: "09:35"}}}

	var _ StockFundFlow = StockFundFlow{Date: "2024-06-13"}
	var _ StockFundFlowDaily = StockFundFlowDaily{Date: "2024-06-13"}
	var _ MarketFundFlow = MarketFundFlow{Date: "2024-06-13"}
	var _ FundFlowRankItem = FundFlowRankItem{Code: "600519"}
	var _ SectorFundFlowItem = SectorFundFlowItem{Code: "BK0001", TopStockCode: stringPtrForDomainTest("600519")}
	var _ NorthboundMinuteItem = NorthboundMinuteItem{Date: "2024-06-13"}
	var _ NorthboundFlowSummary = NorthboundFlowSummary{Date: "2024-06-13"}
	var _ NorthboundHoldingRankItem = NorthboundHoldingRankItem{Code: "600519"}
	var _ NorthboundHistoryItem = NorthboundHistoryItem{Date: "2024-06-13", TopStockCode: stringPtrForDomainTest("600519")}
	var _ NorthboundIndividualItem = NorthboundIndividualItem{Date: "2024-06-13"}

	var _ ZTPoolItem = ZTPoolItem{Code: "600519", FirstBoardTime: stringPtrForDomainTest("09:30:05")}
	var _ StockChangeItem = StockChangeItem{Code: "600519", ChangeType: StockChangeLargeBuy}
	var _ BoardChangeItem = BoardChangeItem{Name: "白酒"}
	var _ THSLimitUpPoolOptions = THSLimitUpPoolOptions{Date: "20250613", OrderField: THSLimitUpOrderLastLimitUpTime, OrderType: THSLimitUpOrderDesc}
	var _ THSLimitUpPoolResult = THSLimitUpPoolResult{Items: []THSLimitUpItem{{Code: "002190"}}}
	var _ DragonTigerDetailItem = DragonTigerDetailItem{Code: "600519"}
	var _ DragonTigerStockStatItem = DragonTigerStockStatItem{Code: "600519"}
	var _ DragonTigerInstitutionItem = DragonTigerInstitutionItem{Code: "600519"}
	var _ DragonTigerBranchItem = DragonTigerBranchItem{Name: "营业部"}
	var _ DragonTigerSeatItem = DragonTigerSeatItem{Side: "buy"}
	var _ DragonTigerPeriod = DragonTigerPeriodOneMonth
	var _ DragonTigerDateOptions = DragonTigerDateOptions{StartDate: "20241201", EndDate: "20241231"}
	var _ BlockTradeDateOptions = BlockTradeDateOptions{StartDate: "20241201", EndDate: "20241231"}

	var _ ETFOptionMonth = ETFOptionMonth{StockID: "510050"}
	var _ ETFOptionExpireDay = ETFOptionExpireDay{ExpireDay: "2024-06-26"}
	var _ ETFOptionCate = ETFOptionCate50ETF
	var _ IndexOptionProduct = IndexOptionProductIO
	var _ OptionKline = OptionKline{Date: "2024-06-13"}
	var _ OptionMinute = OptionMinute{Time: "09:30"}
	var _ OptionTQuote = OptionTQuote{Symbol: "io2501C4000"}
	var _ OptionTQuoteResult = OptionTQuoteResult{Calls: []OptionTQuote{{Symbol: "call"}}}
	var _ CFFEXOptionQuote = CFFEXOptionQuote{Code: "io2501C4000"}
	var _ OptionLHBItem = OptionLHBItem{Symbol: "510050"}

	var _ FundEstimate = FundEstimate{Code: "110011", Name: stringPtrForDomainTest("易方达中小盘")}
	var _ FundNavPoint = FundNavPoint{Date: "2024-06-13", Timestamp: int64PtrForDomainTest(1718208000000)}
	var _ FundNavHistory = FundNavHistory{Code: "110011", Name: stringPtrForDomainTest("易方达中小盘")}
	var _ FundRankPoint = FundRankPoint{Date: "2024-06-13", Timestamp: int64PtrForDomainTest(1718208000000)}
	var _ FundRankHistory = FundRankHistory{Code: "110011", Name: stringPtrForDomainTest("易方达中小盘")}
	var _ FundDividendRank = FundDividendRankExDividendDate
	var _ FundSortDirection = FundSortDesc
	var _ FundDividendListOptions = FundDividendListOptions{Year: "2024", Page: "all", Rank: FundDividendRankExDividendDate, Sort: FundSortDesc}
	var _ FundDividend = FundDividend{Code: "110011", EquityRecordDate: stringPtrForDomainTest("2024-12-16")}
	var _ FundDividendListResult = FundDividendListResult{Items: []FundDividend{{Code: "110011"}}}

	var _ FuturesExchangeCode = FuturesExchangeCodeCFFEX
	var _ FuturesKline = FuturesKline{Code: "rb2605"}
	var _ GlobalFuturesQuote = GlobalFuturesQuote{Code: "GC00Y"}
	var _ FuturesInventorySymbol = FuturesInventorySymbol{Code: "RB"}
	var _ FuturesInventory = FuturesInventory{Code: "RB"}
	var _ ComexInventory = ComexInventory{Name: "COMEX黄金"}
	var _ BlockTradeMarketStatItem = BlockTradeMarketStatItem{Date: "2024-06-13"}
	var _ BlockTradeDetailItem = BlockTradeDetailItem{Code: "600519"}
	var _ BlockTradeDailyStatItem = BlockTradeDailyStatItem{Code: "600519"}
	var _ MarginAccountItem = MarginAccountItem{Date: "2024-06-13"}
	var _ MarginTargetItem = MarginTargetItem{Code: "600519"}
	var _ DividendDetail = DividendDetail{
		Code:             "600519",
		ReportDate:       stringPtrForDomainTest("2024-12-31"),
		DividendDesc:     stringPtrForDomainTest("10派30元"),
		EquityRecordDate: stringPtrForDomainTest("2025-06-20"),
	}
	var _ TodayTimeline = TodayTimeline{Time: "09:30"}
	var _ TodayTimelineResponse = TodayTimelineResponse{Code: "600519", PreClose: float64PtrForDomainTest(10), Data: []TodayTimeline{{Time: "09:30"}}}
}

func TestRootSearchResultTypeConstants(t *testing.T) {
	want := []SearchResultType{SearchStock, SearchIndex, SearchFund, SearchBond, SearchFutures, SearchOption, SearchOther}
	for _, value := range want {
		if value == "" {
			t.Fatalf("empty search result type in %#v", want)
		}
	}
}

func stringPtrForDomainTest(value string) *string {
	return &value
}

func int64PtrForDomainTest(value int64) *int64 {
	return &value
}

func float64PtrForDomainTest(value float64) *float64 {
	return &value
}
