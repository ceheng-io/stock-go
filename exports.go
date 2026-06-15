package stock

import (
	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/internal/providers/tencent"
	"github.com/ceheng.io/stock-go/internal/services"
	domaintypes "github.com/ceheng.io/stock-go/types"
)

type AShareMarket = tencent.AShareMarket

const (
	AShareMarketSH = tencent.AShareMarketSH
	AShareMarketSZ = tencent.AShareMarketSZ
	AShareMarketBJ = tencent.AShareMarketBJ
	AShareMarketKC = tencent.AShareMarketKC
	AShareMarketCY = tencent.AShareMarketCY
)

type CodeListOptions = tencent.CodeListOptions
type GetAShareCodeListOptions = tencent.CodeListOptions

type USMarket = tencent.USMarket

const (
	USMarketNASDAQ = tencent.USMarketNASDAQ
	USMarketNYSE   = tencent.USMarketNYSE
	USMarketAMEX   = tencent.USMarketAMEX
)

type USCodeListOptions = tencent.USCodeListOptions
type GetUSCodeListOptions = tencent.USCodeListOptions

type BatchOptions = tencent.BatchOptions

// GetAllAShareQuotesOptions preserves the TS top-level option name while Go
// services accept code-list and batch options separately.
type GetAllAShareQuotesOptions struct {
	BatchSize   int
	Concurrency int
	OnProgress  func(completed, total int)
	Market      AShareMarket
}

// GetAllUSQuotesOptions preserves the TS top-level option name while Go
// services accept code-list and batch options separately.
type GetAllUSQuotesOptions struct {
	BatchSize   int
	Concurrency int
	OnProgress  func(completed, total int)
	Market      USMarket
}

type GetAllHKQuotesOptions = tencent.BatchOptions

type MarketStatus = services.MarketStatus

const (
	MarketStatusPreMarket  = services.MarketStatusPreMarket
	MarketStatusOpen       = services.MarketStatusOpen
	MarketStatusLunchBreak = services.MarketStatusLunchBreak
	MarketStatusAfterHours = services.MarketStatusAfterHours
	MarketStatusClosed     = services.MarketStatusClosed
)

type SupportedMarket = services.SupportedMarket
type MarketType = services.SupportedMarket

const (
	MarketA  = services.MarketA
	MarketHK = services.MarketHK
	MarketUS = services.MarketUS
)

type KlinePeriod = eastmoney.KlinePeriod

const (
	KlinePeriodDaily   = eastmoney.KlinePeriodDaily
	KlinePeriodWeekly  = eastmoney.KlinePeriodWeekly
	KlinePeriodMonthly = eastmoney.KlinePeriodMonthly
)

type AdjustType = eastmoney.AdjustType

const (
	AdjustNone = eastmoney.AdjustNone
	AdjustQFQ  = eastmoney.AdjustQFQ
	AdjustHFQ  = eastmoney.AdjustHFQ
)

type HistoryKlineOptions = eastmoney.HistoryKlineOptions
type HKKlineOptions = eastmoney.HistoryKlineOptions
type USKlineOptions = eastmoney.HistoryKlineOptions

type KlineWithIndicatorsOptions = services.KlineWithIndicatorsOptions

type MinutePeriod = eastmoney.MinutePeriod

const (
	MinutePeriodOne     = eastmoney.MinutePeriodOne
	MinutePeriodFive    = eastmoney.MinutePeriodFive
	MinutePeriodFifteen = eastmoney.MinutePeriodFifteen
	MinutePeriodThirty  = eastmoney.MinutePeriodThirty
	MinutePeriodSixty   = eastmoney.MinutePeriodSixty
)

type MinuteKlineOptions = eastmoney.MinuteKlineOptions
type HKMinuteKlineOptions = eastmoney.MinuteKlineOptions
type USMinuteKlineOptions = eastmoney.MinuteKlineOptions
type IndustryBoardKlineOptions = eastmoney.HistoryKlineOptions
type IndustryBoardMinuteKlineOptions = eastmoney.MinuteKlineOptions
type ConceptBoardKlineOptions = eastmoney.HistoryKlineOptions
type ConceptBoardMinuteKlineOptions = eastmoney.MinuteKlineOptions

type FuturesKlineOptions = eastmoney.FuturesKlineOptions
type GlobalFuturesSpotOptions = eastmoney.GlobalFuturesSpotOptions
type GlobalFuturesKlineOptions = eastmoney.GlobalFuturesKlineOptions
type FuturesInventoryOptions = eastmoney.FuturesInventoryOptions
type ComexInventoryOptions = eastmoney.ComexInventoryOptions
type CFFEXOptionQuotesOptions = eastmoney.CFFEXOptionQuotesOptions

type FundFlowPeriod = eastmoney.FundFlowPeriod

const (
	FundFlowPeriodDaily   = eastmoney.FundFlowPeriodDaily
	FundFlowPeriodWeekly  = eastmoney.FundFlowPeriodWeekly
	FundFlowPeriodMonthly = eastmoney.FundFlowPeriodMonthly
)

type FundFlowOptions = eastmoney.FundFlowOptions

type FundFlowRankIndicator = eastmoney.FundFlowRankIndicator

const (
	FundFlowRankToday    = eastmoney.FundFlowRankToday
	FundFlowRankThreeDay = eastmoney.FundFlowRankThreeDay
	FundFlowRankFiveDay  = eastmoney.FundFlowRankFiveDay
	FundFlowRankTenDay   = eastmoney.FundFlowRankTenDay
)

type FundFlowSectorType = eastmoney.FundFlowSectorType

const (
	FundFlowSectorIndustry = eastmoney.FundFlowSectorIndustry
	FundFlowSectorConcept  = eastmoney.FundFlowSectorConcept
	FundFlowSectorRegion   = eastmoney.FundFlowSectorRegion
)

type FundFlowRankOptions = eastmoney.FundFlowRankOptions

type NorthboundDirection = domaintypes.NorthboundDirection

const (
	NorthboundNorth = domaintypes.NorthboundNorth
	NorthboundSouth = domaintypes.NorthboundSouth
)

type NorthboundMarket = domaintypes.NorthboundMarket

const (
	NorthboundMarketAll      = domaintypes.NorthboundMarketAll
	NorthboundMarketShanghai = domaintypes.NorthboundMarketShanghai
	NorthboundMarketShenzhen = domaintypes.NorthboundMarketShenzhen
)

type NorthboundRankPeriod = domaintypes.NorthboundRankPeriod

const (
	NorthboundRankToday    = domaintypes.NorthboundRankToday
	NorthboundRankThreeDay = domaintypes.NorthboundRankThreeDay
	NorthboundRankFiveDay  = domaintypes.NorthboundRankFiveDay
	NorthboundRankTenDay   = domaintypes.NorthboundRankTenDay
	NorthboundRankMonth    = domaintypes.NorthboundRankMonth
	NorthboundRankQuarter  = domaintypes.NorthboundRankQuarter
	NorthboundRankYear     = domaintypes.NorthboundRankYear
)

type NorthboundHoldingRankOptions = eastmoney.NorthboundHoldingRankOptions
type NorthboundHistoryOptions = eastmoney.NorthboundHistoryOptions

type DragonTigerPeriod = domaintypes.DragonTigerPeriod

const (
	DragonTigerPeriodOneMonth   = domaintypes.DragonTigerPeriodOneMonth
	DragonTigerPeriodThreeMonth = domaintypes.DragonTigerPeriodThreeMonth
	DragonTigerPeriodSixMonth   = domaintypes.DragonTigerPeriodSixMonth
	DragonTigerPeriodOneYear    = domaintypes.DragonTigerPeriodOneYear
)

type DragonTigerDateOptions = domaintypes.DragonTigerDateOptions

type BlockTradeDateOptions = domaintypes.BlockTradeDateOptions

type DatacenterQuery = eastmoney.DatacenterQuery

type DatacenterResult[T any] struct {
	Data  []T
	Total int
	Pages int
}

// ParseDCDate 提取东方财富 datacenter 常见日期字段为 YYYY-MM-DD。
func ParseDCDate(value any) string {
	return eastmoney.ParseDCDate(value)
}

// ParseDcDate 提取东方财富 datacenter 常见日期字段，保留 TS SDK 命名风格。
func ParseDcDate(value any) string {
	return ParseDCDate(value)
}
