package stock

import domaintypes "github.com/ceheng.io/stock-go/types"

type Quote = domaintypes.Quote
type SimpleQuote = domaintypes.SimpleQuote
type PriceLevel = domaintypes.PriceLevel
type FullQuote = domaintypes.FullQuote
type HKQuote = domaintypes.HKQuote
type USQuote = domaintypes.USQuote
type FundQuote = domaintypes.FundQuote
type FundFlow = domaintypes.FundFlow
type PanelLargeOrder = domaintypes.PanelLargeOrder
type SearchResultType = domaintypes.SearchResultType

const (
	SearchStock   = domaintypes.SearchStock
	SearchIndex   = domaintypes.SearchIndex
	SearchFund    = domaintypes.SearchFund
	SearchBond    = domaintypes.SearchBond
	SearchFutures = domaintypes.SearchFutures
	SearchOption  = domaintypes.SearchOption
	SearchOther   = domaintypes.SearchOther
)

type SearchResult = domaintypes.SearchResult
type Kline = domaintypes.Kline
type AnyHistoryKline = domaintypes.AnyHistoryKline
type HistoryKline = domaintypes.HistoryKline
type ForeignHistoryKline = domaintypes.ForeignHistoryKline
type HKHistoryKline = domaintypes.HKHistoryKline
type USHistoryKline = domaintypes.USHistoryKline
type ForeignMinuteTimeline = domaintypes.ForeignMinuteTimeline
type ForeignMinuteKline = domaintypes.ForeignMinuteKline
type HKMinuteTimeline = domaintypes.HKMinuteTimeline
type HKMinuteKline = domaintypes.HKMinuteKline
type HKMinuteKlineResult = domaintypes.HKMinuteKlineResult
type USMinuteTimeline = domaintypes.USMinuteTimeline
type USMinuteKline = domaintypes.USMinuteKline
type USMinuteKlineResult = domaintypes.USMinuteKlineResult
type MinuteTimeline = domaintypes.MinuteTimeline
type MinuteKline = domaintypes.MinuteKline
type MinuteKlineResult = domaintypes.MinuteKlineResult
type ExternalLink = domaintypes.ExternalLink

type StockIssueInfo = domaintypes.StockIssueInfo
type StockProfile = domaintypes.StockProfile
type FinancialReportPeriod = domaintypes.FinancialReportPeriod
type FinancialIndicatorOptions = domaintypes.FinancialIndicatorOptions
type FinancialIndicator = domaintypes.FinancialIndicator
type AnnouncementOptions = domaintypes.AnnouncementOptions
type StockAnnouncementCode = domaintypes.StockAnnouncementCode
type StockAnnouncementColumn = domaintypes.StockAnnouncementColumn
type StockAnnouncement = domaintypes.StockAnnouncement
type StockAnnouncementResult = domaintypes.StockAnnouncementResult
type StockAnnouncementAttachment = domaintypes.StockAnnouncementAttachment
type StockAnnouncementDetail = domaintypes.StockAnnouncementDetail

const (
	FinancialReportPeriodAll    = domaintypes.FinancialReportPeriodAll
	FinancialReportPeriodAnnual = domaintypes.FinancialReportPeriodAnnual
)

type Board = domaintypes.Board
type IndustryBoard = domaintypes.IndustryBoard
type ConceptBoard = domaintypes.ConceptBoard
type BoardSpot = domaintypes.BoardSpot
type IndustryBoardSpot = domaintypes.IndustryBoardSpot
type ConceptBoardSpot = domaintypes.ConceptBoardSpot
type BoardConstituent = domaintypes.BoardConstituent
type IndustryBoardConstituent = domaintypes.IndustryBoardConstituent
type ConceptBoardConstituent = domaintypes.ConceptBoardConstituent
type BoardKline = domaintypes.BoardKline
type IndustryBoardKline = domaintypes.IndustryBoardKline
type ConceptBoardKline = domaintypes.ConceptBoardKline
type BoardMinuteTimeline = domaintypes.BoardMinuteTimeline
type IndustryBoardMinuteTimeline = domaintypes.IndustryBoardMinuteTimeline
type ConceptBoardMinuteTimeline = domaintypes.ConceptBoardMinuteTimeline
type BoardMinuteKline = domaintypes.BoardMinuteKline
type IndustryBoardMinuteKline = domaintypes.IndustryBoardMinuteKline
type ConceptBoardMinuteKline = domaintypes.ConceptBoardMinuteKline
type BoardMinuteKlineResult = domaintypes.BoardMinuteKlineResult

type StockFundFlow = domaintypes.StockFundFlow
type StockFundFlowDaily = domaintypes.StockFundFlowDaily
type MarketFundFlow = domaintypes.MarketFundFlow
type FundFlowRankItem = domaintypes.FundFlowRankItem
type SectorFundFlowItem = domaintypes.SectorFundFlowItem

type NorthboundMinuteItem = domaintypes.NorthboundMinuteItem
type NorthboundFlowSummary = domaintypes.NorthboundFlowSummary
type NorthboundHoldingRankItem = domaintypes.NorthboundHoldingRankItem
type NorthboundHistoryItem = domaintypes.NorthboundHistoryItem
type NorthboundIndividualItem = domaintypes.NorthboundIndividualItem

type ZTPoolType = domaintypes.ZTPoolType
type ZTPoolItem = domaintypes.ZTPoolItem
type StockChangeType = domaintypes.StockChangeType
type StockChangeItem = domaintypes.StockChangeItem
type BoardChangeItem = domaintypes.BoardChangeItem
type THSLimitUpOrderField = domaintypes.THSLimitUpOrderField
type THSLimitUpOrderType = domaintypes.THSLimitUpOrderType
type THSLimitUpPoolOptions = domaintypes.THSLimitUpPoolOptions
type THSLimitUpPoolResult = domaintypes.THSLimitUpPoolResult
type THSLimitUpPage = domaintypes.THSLimitUpPage
type THSLimitStatGroup = domaintypes.THSLimitStatGroup
type THSLimitStat = domaintypes.THSLimitStat
type THSTradeStatus = domaintypes.THSTradeStatus
type THSLimitUpItem = domaintypes.THSLimitUpItem

const (
	ZTPoolZT        = domaintypes.ZTPoolZT
	ZTPoolYesterday = domaintypes.ZTPoolYesterday
	ZTPoolStrong    = domaintypes.ZTPoolStrong
	ZTPoolSubNew    = domaintypes.ZTPoolSubNew
	ZTPoolBroken    = domaintypes.ZTPoolBroken
	ZTPoolDT        = domaintypes.ZTPoolDT

	StockChangeRocketLaunch   = domaintypes.StockChangeRocketLaunch
	StockChangeQuickRebound   = domaintypes.StockChangeQuickRebound
	StockChangeLargeBuy       = domaintypes.StockChangeLargeBuy
	StockChangeLimitUpSeal    = domaintypes.StockChangeLimitUpSeal
	StockChangeLimitDownOpen  = domaintypes.StockChangeLimitDownOpen
	StockChangeBigBuyOrder    = domaintypes.StockChangeBigBuyOrder
	StockChangeAuctionUp      = domaintypes.StockChangeAuctionUp
	StockChangeHighOpen5D     = domaintypes.StockChangeHighOpen5D
	StockChangeGapUp          = domaintypes.StockChangeGapUp
	StockChangeHigh60D        = domaintypes.StockChangeHigh60D
	StockChangeSurge60D       = domaintypes.StockChangeSurge60D
	StockChangeAccelerateDown = domaintypes.StockChangeAccelerateDown
	StockChangeHighDive       = domaintypes.StockChangeHighDive
	StockChangeLargeSell      = domaintypes.StockChangeLargeSell
	StockChangeLimitDownSeal  = domaintypes.StockChangeLimitDownSeal
	StockChangeLimitUpOpen    = domaintypes.StockChangeLimitUpOpen
	StockChangeBigSellOrder   = domaintypes.StockChangeBigSellOrder
	StockChangeAuctionDown    = domaintypes.StockChangeAuctionDown
	StockChangeLowOpen5D      = domaintypes.StockChangeLowOpen5D
	StockChangeGapDown        = domaintypes.StockChangeGapDown
	StockChangeLow60D         = domaintypes.StockChangeLow60D
	StockChangeDrop60D        = domaintypes.StockChangeDrop60D

	THSLimitUpOrderFirstLimitUpTime = domaintypes.THSLimitUpOrderFirstLimitUpTime
	THSLimitUpOrderLastLimitUpTime  = domaintypes.THSLimitUpOrderLastLimitUpTime
	THSLimitUpOrderOpenNum          = domaintypes.THSLimitUpOrderOpenNum
	THSLimitUpOrderDesc             = domaintypes.THSLimitUpOrderDesc
	THSLimitUpOrderAsc              = domaintypes.THSLimitUpOrderAsc
)

type DragonTigerDetailItem = domaintypes.DragonTigerDetailItem
type DragonTigerStockStatItem = domaintypes.DragonTigerStockStatItem
type DragonTigerInstitutionItem = domaintypes.DragonTigerInstitutionItem
type DragonTigerBranchItem = domaintypes.DragonTigerBranchItem
type DragonTigerSeatItem = domaintypes.DragonTigerSeatItem

type ETFOptionMonth = domaintypes.ETFOptionMonth
type ETFOptionExpireDay = domaintypes.ETFOptionExpireDay
type ETFOptionCate = domaintypes.ETFOptionCate
type IndexOptionProduct = domaintypes.IndexOptionProduct
type OptionKline = domaintypes.OptionKline
type OptionMinute = domaintypes.OptionMinute
type OptionTQuote = domaintypes.OptionTQuote
type OptionTQuoteResult = domaintypes.OptionTQuoteResult
type CFFEXOptionQuote = domaintypes.CFFEXOptionQuote
type OptionLHBItem = domaintypes.OptionLHBItem

const (
	ETFOptionCate50ETF           = domaintypes.ETFOptionCate50ETF
	ETFOptionCate300ETF          = domaintypes.ETFOptionCate300ETF
	ETFOptionCate500ETF          = domaintypes.ETFOptionCate500ETF
	ETFOptionCateKechuang50      = domaintypes.ETFOptionCateKechuang50
	ETFOptionCateKechuangBoard50 = domaintypes.ETFOptionCateKechuangBoard50

	IndexOptionProductHO = domaintypes.IndexOptionProductHO
	IndexOptionProductIO = domaintypes.IndexOptionProductIO
	IndexOptionProductMO = domaintypes.IndexOptionProductMO
)

type FundEstimate = domaintypes.FundEstimate
type FundNavPoint = domaintypes.FundNavPoint
type FundNavHistory = domaintypes.FundNavHistory
type FundRankPoint = domaintypes.FundRankPoint
type FundRankHistory = domaintypes.FundRankHistory
type FundDividendRank = domaintypes.FundDividendRank
type FundSortDirection = domaintypes.FundSortDirection
type FundDividendListOptions = domaintypes.FundDividendListOptions
type FundDividend = domaintypes.FundDividend
type FundDividendListResult = domaintypes.FundDividendListResult

const (
	FundDividendRankCode             = domaintypes.FundDividendRankCode
	FundDividendRankName             = domaintypes.FundDividendRankName
	FundDividendRankEquityRecordDate = domaintypes.FundDividendRankEquityRecordDate
	FundDividendRankExDividendDate   = domaintypes.FundDividendRankExDividendDate
	FundDividendRankDividendPerShare = domaintypes.FundDividendRankDividendPerShare
	FundDividendRankPayDate          = domaintypes.FundDividendRankPayDate

	FundSortAsc  = domaintypes.FundSortAsc
	FundSortDesc = domaintypes.FundSortDesc
)

// FuturesExchangeCode is the domestic futures exchange code type from types.
//
// The root package already exposes symbols.FuturesExchange as FuturesExchange,
// so the TS-compatible domain type keeps a Code suffix at the root.
type FuturesExchangeCode = domaintypes.FuturesExchange
type FuturesKline = domaintypes.FuturesKline
type GlobalFuturesQuote = domaintypes.GlobalFuturesQuote
type FuturesInventorySymbol = domaintypes.FuturesInventorySymbol
type FuturesInventory = domaintypes.FuturesInventory
type ComexInventory = domaintypes.ComexInventory

const (
	FuturesExchangeCodeSHFE  = domaintypes.FuturesExchangeSHFE
	FuturesExchangeCodeDCE   = domaintypes.FuturesExchangeDCE
	FuturesExchangeCodeCZCE  = domaintypes.FuturesExchangeCZCE
	FuturesExchangeCodeINE   = domaintypes.FuturesExchangeINE
	FuturesExchangeCodeCFFEX = domaintypes.FuturesExchangeCFFEX
	FuturesExchangeCodeGFEX  = domaintypes.FuturesExchangeGFEX
)

type BlockTradeMarketStatItem = domaintypes.BlockTradeMarketStatItem
type BlockTradeDetailItem = domaintypes.BlockTradeDetailItem
type BlockTradeDailyStatItem = domaintypes.BlockTradeDailyStatItem

type MarginAccountItem = domaintypes.MarginAccountItem
type MarginTargetItem = domaintypes.MarginTargetItem
type DividendDetail = domaintypes.DividendDetail
type TodayTimeline = domaintypes.TodayTimeline
type TodayTimelineResponse = domaintypes.TodayTimelineResponse
