package eastmoney

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/ceheng-io/stock-go/internal/core"
)

type fakeKlineClient struct {
	lastURL string
}

func (f *fakeKlineClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	payload := map[string]any{
		"data": map[string]any{
			"code": "600519",
			"name": "贵州茅台",
			"klines": []string{
				"2024-12-16,1500.00,1510.00,1520.00,1490.00,12345,67890000,2.00,1.50,22.30,0.50",
				"2024-12-17,1510.00,1525.00,1530.00,1505.00,22345,77890000,1.66,0.99,15.00,0.60",
			},
		},
	}
	b, _ := json.Marshal(payload)
	return json.Unmarshal(b, target)
}

func TestGetHistoryKlineBuildsEastmoneyRequestAndParsesRows(t *testing.T) {
	client := &fakeKlineClient{}

	rows, err := GetHistoryKline(context.Background(), client, "sh600519", "https://em.test/kline", HistoryKlineOptions{
		Period:    KlinePeriodWeekly,
		Adjust:    AdjustHFQ,
		StartDate: "20241201",
		EndDate:   "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	first := rows[0]
	if first.Code != "600519" || first.Date != "2024-12-16" {
		t.Fatalf("first identity = %+v", first)
	}
	if first.Timestamp == nil || *first.Timestamp != 1734278400000 || first.TZ != "Asia/Shanghai" {
		t.Fatalf("first time meta = %+v", first)
	}
	if first.Open == nil || *first.Open != 1500 || first.Close == nil || *first.Close != 1510 {
		t.Fatalf("first prices = %+v", first)
	}
	if first.Amount == nil || *first.Amount != 67890000 || first.TurnoverRate == nil || *first.TurnoverRate != 0.5 {
		t.Fatalf("first numeric fields = %+v", first)
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "secid", "1.600519")
	assertQuery(t, query, "klt", "102")
	assertQuery(t, query, "fqt", "2")
	assertQuery(t, query, "ut", "7eea3edcaed734bea9cbfc24409ed989")
	if !strings.Contains(query.Get("fields2"), "f116") {
		t.Fatalf("fields2 = %q, want f116 included", query.Get("fields2"))
	}
	assertQuery(t, query, "beg", "20241201")
	assertQuery(t, query, "end", "20241231")
}

func TestGetHistoryKlineDefaultsAndValidation(t *testing.T) {
	client := &fakeKlineClient{}
	if _, err := GetHistoryKline(context.Background(), client, "000001", "https://em.test/kline", HistoryKlineOptions{}); err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "0.000001")
	assertQuery(t, parsed.Query(), "klt", "101")
	assertQuery(t, parsed.Query(), "fqt", "1")

	if _, err := GetHistoryKline(context.Background(), client, "000001", "https://em.test/kline", HistoryKlineOptions{Period: "yearly"}); err == nil {
		t.Fatal("invalid period expected error")
	}
	if _, err := GetHistoryKline(context.Background(), client, "000001", "https://em.test/kline", HistoryKlineOptions{Adjust: "bad"}); err == nil {
		t.Fatal("invalid adjust expected error")
	}
}

type fakeForeignHistoryKlineClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeForeignHistoryKlineClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	b, _ := json.Marshal(f.payload)
	return json.Unmarshal(b, target)
}

func TestGetHKHistoryKlineBuildsRequestAndAddsHKFields(t *testing.T) {
	client := &fakeForeignHistoryKlineClient{payload: map[string]any{
		"data": map[string]any{
			"code": "00700",
			"name": "腾讯控股",
			"klines": []string{
				"2024-12-16,390.00,392.00,395.00,388.00,12345,67890000,1.80,0.51,2.00,0.03",
			},
		},
	}}

	rows, err := GetHKHistoryKline(context.Background(), client, "hk700", "https://em.test/hk", HistoryKlineOptions{
		Period:    KlinePeriodMonthly,
		Adjust:    AdjustNone,
		StartDate: "20240101",
		EndDate:   "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "00700" || row.Name != "腾讯控股" || row.Currency != "HKD" || row.LotSize != nil {
		t.Fatalf("row meta = %+v", row)
	}
	if row.Close == nil || *row.Close != 392 || row.Change == nil || *row.Change != 2 {
		t.Fatalf("row numeric fields = %+v", row)
	}
	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "secid", "116.00700")
	assertQuery(t, query, "klt", "103")
	assertQuery(t, query, "fqt", "0")
	assertQuery(t, query, "lmt", "1000000")
}

func TestGetUSHistoryKlineBuildsRequestAndAddsUSFields(t *testing.T) {
	client := &fakeForeignHistoryKlineClient{payload: map[string]any{
		"data": map[string]any{
			"code": "AAPL",
			"name": "Apple Inc.",
			"klines": []string{
				"2024-12-16,250.00,251.50,253.00,249.00,34567,87650000,1.60,0.80,2.00,0.10",
			},
		},
	}}

	rows, err := GetUSHistoryKline(context.Background(), client, "105.AAPL", "https://em.test/us", HistoryKlineOptions{
		Period:    KlinePeriodWeekly,
		Adjust:    AdjustHFQ,
		StartDate: "20240101",
		EndDate:   "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "AAPL" || row.Name != "Apple Inc." || row.Currency != "USD" {
		t.Fatalf("row meta = %+v", row)
	}
	if row.Close == nil || *row.Close != 251.5 || row.Amount == nil || *row.Amount != 87650000 {
		t.Fatalf("row numeric fields = %+v", row)
	}
	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "secid", "105.AAPL")
	assertQuery(t, query, "klt", "102")
	assertQuery(t, query, "fqt", "2")
	assertQuery(t, query, "lmt", "1000000")
}

func TestGetUSKlineEmptySymbolReturnsInvalidSymbol(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "history",
			call: func() error {
				_, err := GetUSHistoryKline(context.Background(), &fakeForeignHistoryKlineClient{}, " ", "https://em.test/us", HistoryKlineOptions{})
				return err
			},
		},
		{
			name: "minute",
			call: func() error {
				_, err := GetUSMinuteKline(context.Background(), &fakeMinuteClient{}, " ", "https://em.test/us-kline", "https://em.test/us-trends", MinuteKlineOptions{})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid symbol error")
			}
			var coded core.CodedError
			if !errors.As(err, &coded) {
				t.Fatalf("err = %T %v, want coded invalid symbol error", err, err)
			}
			if code := coded.SDKCode(); code != "INVALID_SYMBOL" {
				t.Fatalf("error code = %q, want INVALID_SYMBOL; err=%v", code, err)
			}
		})
	}
}

type fakeMinuteClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeMinuteClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	b, _ := json.Marshal(f.payload)
	return json.Unmarshal(b, target)
}

func TestGetMinuteKlinePeriodOneParsesTimelineAndFilters(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{
			"trends": []string{
				"2024-12-16 09:30,1500,1501,1502,1499,100,200000,1500.5",
				"2024-12-16 09:31,1501,1502,1503,1500,110,210000,1501.5",
				"2024-12-17 09:31,1505,1506,1507,1504,120,220000,1505.5",
			},
		},
	}}

	rows, err := GetMinuteKline(context.Background(), client, "600519", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{
		Period:    MinutePeriodOne,
		StartDate: "2024-12-16 09:31",
		EndDate:   "2024-12-16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 1 {
		t.Fatalf("timeline = %+v", rows.Timeline)
	}
	row := rows.Timeline[0]
	if row.Time != "2024-12-16 09:31" || row.Close == nil || *row.Close != 1502 || row.AvgPrice == nil || *row.AvgPrice != 1501.5 {
		t.Fatalf("row = %+v", row)
	}
	if row.Timestamp == nil || *row.Timestamp != 1734312660000 || row.TZ != "Asia/Shanghai" {
		t.Fatalf("row time meta = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "1.600519")
	assertQuery(t, parsed.Query(), "ndays", "5")
}

func TestGetMinuteKlinePeriodOneReturnsEmptyRowsForNonArrayTrends(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"trends": map[string]any{"time": "2024-12-16 09:31"}},
	}}

	rows, err := GetMinuteKline(context.Background(), client, "600519", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 0 {
		t.Fatalf("timeline = %+v, want empty", rows.Timeline)
	}
}

func TestGetMinuteKlinePeriodFiveParsesKlineRows(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{
			"klines": []string{
				"2024-12-16 09:35,1500,1510,1520,1490,12345,67890000,2.00,1.50,22.30,0.50",
				"2024-12-16 09:40,1510,1525,1530,1505,22345,77890000,1.66,0.99,15.00,0.60",
			},
		},
	}}

	rows, err := GetMinuteKline(context.Background(), client, "sh600519", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{
		Period:    MinutePeriodFive,
		Adjust:    AdjustNone,
		StartDate: "2024-12-16 09:36",
		EndDate:   "2024-12-16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 1 {
		t.Fatalf("klines = %+v", rows.Klines)
	}
	row := rows.Klines[0]
	if row.Time != "2024-12-16 09:40" || row.Close == nil || *row.Close != 1525 || row.Change == nil || *row.Change != 15 {
		t.Fatalf("row = %+v", row)
	}
	if row.Timestamp == nil || *row.Timestamp != 1734313200000 || row.TZ != "Asia/Shanghai" {
		t.Fatalf("row time meta = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "klt", "5")
	assertQuery(t, parsed.Query(), "fqt", "0")
}

func TestGetMinuteKlinePeriodFiveReturnsEmptyRowsForNonArrayKlines(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"klines": map[string]any{"time": "2024-12-16 09:35"}},
	}}

	rows, err := GetMinuteKline(context.Background(), client, "sh600519", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 0 {
		t.Fatalf("klines = %+v, want empty", rows.Klines)
	}
}

func TestGetHKMinuteKlinePeriodOneParsesTimeline(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{
			"trends": []string{
				"2024-12-16 09:30,390,391,392,389,100,200000,390.5",
				"2024-12-16 09:31,391,392,393,390,110,210000,391.5",
			},
		},
	}}

	rows, err := GetHKMinuteKline(context.Background(), client, "hk700", "https://em.test/hk-kline", "https://em.test/hk-trends", MinuteKlineOptions{
		Period:    MinutePeriodOne,
		StartDate: "2024-12-16 09:31",
		EndDate:   "2024-12-16",
		NDays:     3,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 1 {
		t.Fatalf("timeline = %+v", rows.Timeline)
	}
	row := rows.Timeline[0]
	if row.Time != "2024-12-16 09:31" || row.Code != "00700" || row.Currency != "HKD" || row.TZ != "Asia/Hong_Kong" {
		t.Fatalf("row meta = %+v", row)
	}
	if row.Close == nil || *row.Close != 392 || row.AvgPrice == nil || *row.AvgPrice != 391.5 {
		t.Fatalf("row numeric fields = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "116.00700")
	assertQuery(t, parsed.Query(), "ndays", "3")
}

func TestGetHKMinuteKlinePeriodOneDefaultsNDaysToOne(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"trends": []string{}},
	}}

	_, err := GetHKMinuteKline(context.Background(), client, "00700", "https://em.test/hk-kline", "https://em.test/hk-trends", MinuteKlineOptions{
		Period: MinutePeriodOne,
	})
	if err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "ndays", "1")
}

func TestGetHKMinuteKlinePeriodOneKeepsExplicitNDays(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"trends": []string{}},
	}}

	_, err := GetHKMinuteKline(context.Background(), client, "00700", "https://em.test/hk-kline", "https://em.test/hk-trends", MinuteKlineOptions{
		Period: MinutePeriodOne,
		NDays:  5,
	})
	if err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "ndays", "5")
}

func TestGetHKMinuteKlinePeriodOneReturnsEmptyRowsForNonArrayTrends(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"trends": map[string]any{"time": "2024-12-16 09:31"}},
	}}

	rows, err := GetHKMinuteKline(context.Background(), client, "00700", "https://em.test/hk-kline", "https://em.test/hk-trends", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 0 {
		t.Fatalf("timeline = %+v, want empty", rows.Timeline)
	}
}

func TestGetUSMinuteKlinePeriodFiveConvertsCNTimeToUSLocalTime(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{
			"klines": []string{
				"2024-12-16 22:30,250,251,252,249,100,200000,1.20,0.40,1.00,0.01",
				"2024-12-16 22:35,251,252,253,250,110,210000,1.19,0.40,1.00,0.02",
			},
		},
	}}

	rows, err := GetUSMinuteKline(context.Background(), client, "105.AAPL", "https://em.test/us-kline", "https://em.test/us-trends", MinuteKlineOptions{
		Period:    MinutePeriodFive,
		Adjust:    AdjustHFQ,
		StartDate: "2024-12-16 09:35",
		EndDate:   "2024-12-16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 1 {
		t.Fatalf("klines = %+v", rows.Klines)
	}
	row := rows.Klines[0]
	if row.Time != "2024-12-16 09:35" || row.Code != "AAPL" || row.Currency != "USD" || row.TZ != "America/New_York" {
		t.Fatalf("row meta = %+v", row)
	}
	if row.Close == nil || *row.Close != 252 || row.Change == nil || *row.Change != 1 {
		t.Fatalf("row numeric fields = %+v", row)
	}
	if row.Timestamp == nil {
		t.Fatalf("timestamp is nil: %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "105.AAPL")
	assertQuery(t, parsed.Query(), "klt", "5")
	assertQuery(t, parsed.Query(), "fqt", "2")
}

func TestGetUSMinuteKlinePeriodFiveReturnsEmptyRowsForNonArrayKlines(t *testing.T) {
	client := &fakeMinuteClient{payload: map[string]any{
		"data": map[string]any{"klines": map[string]any{"time": "2024-12-16 22:30"}},
	}}

	rows, err := GetUSMinuteKline(context.Background(), client, "105.AAPL", "https://em.test/us-kline", "https://em.test/us-trends", MinuteKlineOptions{Period: MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 0 {
		t.Fatalf("klines = %+v, want empty", rows.Klines)
	}
}

func assertQuery(t *testing.T, values url.Values, key, want string) {
	t.Helper()
	if got := values.Get(key); got != want {
		t.Fatalf("query %s = %q, want %q (url values: %s)", key, got, want, values.Encode())
	}
}
