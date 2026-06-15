package services

import (
	"context"
	"errors"
	"testing"

	"github.com/ceheng.io/stock-go/indicators"
	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

type indicatorKlineStub struct {
	cnRows []types.HistoryKline
	hkRows []types.HKHistoryKline
	usRows []types.USHistoryKline

	cnCalls []eastmoney.HistoryKlineOptions
	hkCalls []eastmoney.HistoryKlineOptions
	usCalls []eastmoney.HistoryKlineOptions
}

func (s *indicatorKlineStub) CN(_ context.Context, _ string, options eastmoney.HistoryKlineOptions) ([]types.HistoryKline, error) {
	s.cnCalls = append(s.cnCalls, options)
	return append([]types.HistoryKline(nil), s.cnRows...), nil
}

func (s *indicatorKlineStub) HK(_ context.Context, _ string, options eastmoney.HistoryKlineOptions) ([]types.HKHistoryKline, error) {
	s.hkCalls = append(s.hkCalls, options)
	return append([]types.HKHistoryKline(nil), s.hkRows...), nil
}

func (s *indicatorKlineStub) US(_ context.Context, _ string, options eastmoney.HistoryKlineOptions) ([]types.USHistoryKline, error) {
	s.usCalls = append(s.usCalls, options)
	return append([]types.USHistoryKline(nil), s.usRows...), nil
}

type indicatorCalendarStub struct {
	days []string
	err  error
}

func (s indicatorCalendarStub) TradingCalendar(context.Context) ([]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	return append([]string(nil), s.days...), nil
}

func TestIndicatorServiceUsesCalendarLookbackAndFiltersRequestedDates(t *testing.T) {
	kline := &indicatorKlineStub{
		cnRows: []types.HistoryKline{
			historyKline("2024-06-06", 9),
			historyKline("2024-06-07", 10),
			historyKline("2024-06-10", 11),
			historyKline("2024-06-11", 12),
			historyKline("2024-06-12", 13),
			historyKline("2024-06-13", 14),
		},
	}
	service := NewIndicatorService(kline, indicatorCalendarStub{days: []string{
		"2024-06-06",
		"2024-06-07",
		"2024-06-10",
		"2024-06-11",
		"2024-06-12",
		"2024-06-13",
	}})

	got, err := service.KlineWithIndicators(context.Background(), "600519", KlineWithIndicatorsOptions{
		Market:    MarketA,
		StartDate: "2024-06-12",
		EndDate:   "2024-06-13",
		Indicators: indicators.Options{
			MA: &indicators.MAOptions{Periods: []int{3}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(kline.cnCalls) != 1 {
		t.Fatalf("CN calls = %d, want 1", len(kline.cnCalls))
	}
	if got := kline.cnCalls[0].StartDate; got != "20240606" {
		t.Fatalf("StartDate = %q, want 20240606", got)
	}
	if got := kline.cnCalls[0].EndDate; got != "20240613" {
		t.Fatalf("EndDate = %q, want 20240613", got)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2: %+v", len(got), got)
	}
	if got[0].Date != "2024-06-12" || got[1].Date != "2024-06-13" {
		t.Fatalf("dates = %q, %q", got[0].Date, got[1].Date)
	}
	assertIndicatorFloat(t, got[0].MA["ma3"], 12)
	assertIndicatorFloat(t, got[1].MA["ma3"], 13)
}

func TestIndicatorServiceFallsBackToNaturalDateLookback(t *testing.T) {
	kline := &indicatorKlineStub{
		cnRows: []types.HistoryKline{
			historyKline("2024-06-10", 11),
			historyKline("2024-06-11", 12),
			historyKline("2024-06-12", 13),
		},
	}
	service := NewIndicatorService(kline, indicatorCalendarStub{err: errors.New("calendar unavailable")})

	_, err := service.KlineWithIndicators(context.Background(), "600519", KlineWithIndicatorsOptions{
		Market:    MarketA,
		StartDate: "2024-06-12",
		Indicators: indicators.Options{
			MA: &indicators.MAOptions{Periods: []int{2}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(kline.cnCalls) != 1 {
		t.Fatalf("CN calls = %d, want 1", len(kline.cnCalls))
	}
	if got := kline.cnCalls[0].StartDate; got != "20240607" {
		t.Fatalf("StartDate = %q, want natural-day fallback 20240607", got)
	}
}

func TestIndicatorServiceDetectsForeignMarkets(t *testing.T) {
	kline := &indicatorKlineStub{
		hkRows: []types.HKHistoryKline{{
			ForeignHistoryKline: foreignHistoryKline("2024-06-12", 20),
			Currency:            "HKD",
		}},
		usRows: []types.USHistoryKline{{
			ForeignHistoryKline: foreignHistoryKline("2024-06-12", 30),
			Currency:            "USD",
		}},
	}
	service := NewIndicatorService(kline, indicatorCalendarStub{})

	if _, err := service.KlineWithIndicators(context.Background(), "00700.HK", KlineWithIndicatorsOptions{}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.KlineWithIndicators(context.Background(), "AAPL", KlineWithIndicatorsOptions{}); err != nil {
		t.Fatal(err)
	}

	if len(kline.hkCalls) != 1 {
		t.Fatalf("HK calls = %d, want 1", len(kline.hkCalls))
	}
	if len(kline.usCalls) != 1 {
		t.Fatalf("US calls = %d, want 1", len(kline.usCalls))
	}
}

func historyKline(date string, close float64) types.HistoryKline {
	return types.HistoryKline{
		Date:   date,
		Code:   "600519",
		Open:   indicators.Float(close - 1),
		High:   indicators.Float(close + 1),
		Low:    indicators.Float(close - 2),
		Close:  indicators.Float(close),
		Volume: indicators.Float(close * 100),
	}
}

func foreignHistoryKline(date string, close float64) types.ForeignHistoryKline {
	return types.ForeignHistoryKline{
		Date:   date,
		Open:   indicators.Float(close - 1),
		High:   indicators.Float(close + 1),
		Low:    indicators.Float(close - 2),
		Close:  indicators.Float(close),
		Volume: indicators.Float(close * 100),
	}
}

func assertIndicatorFloat(t *testing.T, value *float64, want float64) {
	t.Helper()
	if value == nil {
		t.Fatalf("value is nil, want %.2f", want)
	}
	if *value != want {
		t.Fatalf("value = %.2f, want %.2f", *value, want)
	}
}
