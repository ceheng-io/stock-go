package services

import (
	"context"

	"github.com/ceheng.io/stock-go/indicators"
	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/symbols"
	"github.com/ceheng.io/stock-go/types"
)

// IndicatorKlineService 是 IndicatorService 依赖的 K 线接口。
type IndicatorKlineService interface {
	CN(context.Context, string, eastmoney.HistoryKlineOptions) ([]types.HistoryKline, error)
	HK(context.Context, string, eastmoney.HistoryKlineOptions) ([]types.HKHistoryKline, error)
	US(context.Context, string, eastmoney.HistoryKlineOptions) ([]types.USHistoryKline, error)
}

// KlineWithIndicatorsOptions 配置 K 线指标聚合请求。
type KlineWithIndicatorsOptions struct {
	Market     SupportedMarket
	Period     eastmoney.KlinePeriod
	Adjust     eastmoney.AdjustType
	StartDate  string
	EndDate    string
	Indicators indicators.Options
}

// IndicatorService 编排 K 线和技术指标能力。
type IndicatorService struct {
	kline  IndicatorKlineService
	quotes CalendarQuoteService
}

// NewIndicatorService 创建 IndicatorService。
func NewIndicatorService(kline IndicatorKlineService, quotes CalendarQuoteService) *IndicatorService {
	return &IndicatorService{kline: kline, quotes: quotes}
}

// KlineWithIndicators 返回附带技术指标的历史 K 线。
func (s *IndicatorService) KlineWithIndicators(ctx context.Context, symbol string, options KlineWithIndicatorsOptions) ([]indicators.KlineWithIndicators, error) {
	market := options.Market
	if market == "" {
		market = detectIndicatorMarket(symbol)
	}
	lookback := indicators.EstimateLookback(options.Indicators)
	klineOptions := eastmoney.HistoryKlineOptions{
		Period:    options.Period,
		Adjust:    options.Adjust,
		StartDate: indicatorCompactDate(options.StartDate),
		EndDate:   indicatorCompactDate(options.EndDate),
	}
	if options.StartDate != "" {
		klineOptions.StartDate = s.actualStartDate(ctx, market, options.StartDate, lookback.RequiredBars)
	}

	rows, err := s.fetchKlines(ctx, market, symbol, klineOptions)
	if err != nil {
		return nil, err
	}
	if options.StartDate != "" && len(rows) < lookback.RequiredBars {
		klineOptions.StartDate = ""
		rows, err = s.fetchKlines(ctx, market, symbol, klineOptions)
		if err != nil {
			return nil, err
		}
	}

	withIndicators := indicators.AddIndicators(rows, options.Indicators)
	return filterIndicatorRows(withIndicators, options.StartDate, options.EndDate), nil
}

func (s *IndicatorService) fetchKlines(ctx context.Context, market SupportedMarket, symbol string, options eastmoney.HistoryKlineOptions) ([]indicators.KlineInput, error) {
	switch market {
	case MarketHK:
		rows, err := s.kline.HK(ctx, symbol, options)
		if err != nil {
			return nil, err
		}
		return hkKlineInputs(rows), nil
	case MarketUS:
		rows, err := s.kline.US(ctx, symbol, options)
		if err != nil {
			return nil, err
		}
		return usKlineInputs(rows), nil
	default:
		rows, err := s.kline.CN(ctx, symbol, options)
		if err != nil {
			return nil, err
		}
		return cnKlineInputs(rows), nil
	}
}

func (s *IndicatorService) actualStartDate(ctx context.Context, market SupportedMarket, startDate string, requiredBars int) string {
	if market == MarketA && s.quotes != nil {
		calendar, err := s.quotes.TradingCalendar(ctx)
		if err == nil {
			if actualStart := actualStartDateByCalendar(startDate, requiredBars, calendar); actualStart != "" {
				return actualStart
			}
		}
	}
	return actualStartDateByNaturalDays(startDate, requiredBars, indicatorRatio(market))
}

func cnKlineInputs(rows []types.HistoryKline) []indicators.KlineInput {
	result := make([]indicators.KlineInput, len(rows))
	for i, row := range rows {
		result[i] = indicators.KlineInput{
			Date:   row.Date,
			Open:   row.Open,
			High:   row.High,
			Low:    row.Low,
			Close:  row.Close,
			Volume: row.Volume,
		}
	}
	return result
}

func hkKlineInputs(rows []types.HKHistoryKline) []indicators.KlineInput {
	result := make([]indicators.KlineInput, len(rows))
	for i, row := range rows {
		result[i] = foreignKlineInput(row.ForeignHistoryKline)
	}
	return result
}

func usKlineInputs(rows []types.USHistoryKline) []indicators.KlineInput {
	result := make([]indicators.KlineInput, len(rows))
	for i, row := range rows {
		result[i] = foreignKlineInput(row.ForeignHistoryKline)
	}
	return result
}

func foreignKlineInput(row types.ForeignHistoryKline) indicators.KlineInput {
	return indicators.KlineInput{
		Date:   row.Date,
		Open:   row.Open,
		High:   row.High,
		Low:    row.Low,
		Close:  row.Close,
		Volume: row.Volume,
	}
}

func filterIndicatorRows(rows []indicators.KlineWithIndicators, startDate string, endDate string) []indicators.KlineWithIndicators {
	if startDate == "" && endDate == "" {
		return rows
	}
	start := indicatorDateKey(startDate)
	end := indicatorDateKey(endDate)
	result := make([]indicators.KlineWithIndicators, 0, len(rows))
	for _, row := range rows {
		date := indicatorDateKey(row.Date)
		if start != "" && date < start {
			continue
		}
		if end != "" && date > end {
			continue
		}
		result = append(result, row)
	}
	return result
}

func detectIndicatorMarket(symbol string) SupportedMarket {
	normalized, err := symbols.Normalize(symbol, nil)
	if err != nil {
		return MarketA
	}
	switch normalized.Market {
	case symbols.MarketHK:
		return MarketHK
	case symbols.MarketUS:
		return MarketUS
	default:
		return MarketA
	}
}
