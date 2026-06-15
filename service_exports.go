package stock

import "github.com/ceheng.io/stock-go/internal/services"

// QuoteService is the root quote service type.
type QuoteService = services.QuoteService

// KlineService is the root K-line service type.
type KlineService = services.KlineService

// IndicatorService is the root indicator service type.
type IndicatorService = services.IndicatorService

// BoardService is the root board service type.
type BoardService = services.BoardService

// TradingCalendarService preserves the TypeScript SDK top-level service type name.
type TradingCalendarService = services.CalendarService

// CalendarService is the root trading-calendar service type.
type CalendarService = services.CalendarService

// FundFlowService is the root fund-flow service type.
type FundFlowService = services.FundFlowService

// NorthboundService is the root northbound service type.
type NorthboundService = services.NorthboundService

// DragonTigerService is the root dragon-tiger service type.
type DragonTigerService = services.DragonTigerService

// BlockTradeService is the root block-trade service type.
type BlockTradeService = services.BlockTradeService

// MarginService is the root margin service type.
type MarginService = services.MarginService

// DividendService is the root dividend service type.
type DividendService = services.DividendService

// DataService is the root data aggregation service type.
type DataService = services.DataService

// MarketEventService is the root market-event service type.
type MarketEventService = services.MarketEventService

// FundService preserves the TypeScript SDK top-level service type name.
type FundService = services.FundService

// FuturesService is the root futures service type.
type FuturesService = services.FuturesService

// OptionsService is the root options service type.
type OptionsService = services.OptionsService
