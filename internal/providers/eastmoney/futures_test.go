package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

type fakeFuturesClient struct {
	lastURL  string
	urls     []string
	payload  map[string]any
	payloads []map[string]any
}

func (f *fakeFuturesClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	f.urls = append(f.urls, requestURL)
	if len(f.payloads) > 0 {
		payload := f.payloads[0]
		f.payloads = f.payloads[1:]
		body, _ := json.Marshal(payload)
		return json.Unmarshal(body, target)
	}
	if f.payload == nil {
		f.payload = map[string]any{
			"data": map[string]any{
				"code": "rb2605",
				"name": "螺纹钢2605",
				"klines": []string{
					"2024-12-16,3500,3520,3530,3490,12345,67890000,1.14,0.57,20,0,0,98765,0",
				},
			},
		}
	}
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetFuturesHistoryKlineBuildsRequestAndParsesOpenInterest(t *testing.T) {
	client := &fakeFuturesClient{}

	rows, err := GetFuturesHistoryKline(context.Background(), client, "rb2605", "https://em.test/futures", FuturesKlineOptions{
		Period:    KlinePeriodWeekly,
		StartDate: "20241201",
		EndDate:   "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "rb2605" || row.Name != "螺纹钢2605" || row.Date != "2024-12-16" {
		t.Fatalf("row meta = %+v", row)
	}
	if row.Close == nil || *row.Close != 3520 || row.OpenInterest == nil || *row.OpenInterest != 98765 {
		t.Fatalf("row numeric = %+v", row)
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "secid", "113.rb2605")
	assertQuery(t, query, "klt", "102")
	assertQuery(t, query, "fqt", "0")
	assertQuery(t, query, "beg", "20241201")
	assertQuery(t, query, "end", "20241231")
}

func TestGetFuturesHistoryKlineSupportsMainContractAndValidation(t *testing.T) {
	client := &fakeFuturesClient{}

	if _, err := GetFuturesHistoryKline(context.Background(), client, "RBM", "https://em.test/futures", FuturesKlineOptions{}); err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "113.RBM")
	assertQuery(t, parsed.Query(), "klt", "101")

	if _, err := GetFuturesHistoryKline(context.Background(), client, "???", "https://em.test/futures", FuturesKlineOptions{}); err == nil {
		t.Fatal("expected invalid symbol error")
	}
	if _, err := GetFuturesHistoryKline(context.Background(), client, "unknown2605", "https://em.test/futures", FuturesKlineOptions{}); err == nil {
		t.Fatal("expected unknown variety error")
	}
	if _, err := GetFuturesHistoryKline(context.Background(), client, "rb2605", "https://em.test/futures", FuturesKlineOptions{Period: "yearly"}); err == nil {
		t.Fatal("expected invalid period error")
	}
}

func TestGetGlobalFuturesSpotFetchesAllPages(t *testing.T) {
	client := &fakeFuturesClient{payloads: []map[string]any{
		{
			"total": float64(2),
			"list": []map[string]any{
				{"dm": "GC00Y", "name": "COMEX黄金", "p": 2400.5, "zde": 10.5, "zdf": 0.44, "o": 2390, "h": 2410, "l": 2380, "zjsj": 2390, "vol": 1000, "wp": 10, "np": 20, "ccl": 3000},
			},
		},
		{
			"total": float64(2),
			"list": []map[string]any{
				{"dm": "SI00Y", "name": "COMEX白银", "p": 30.5, "zde": -0.5, "zdf": -1.61, "o": 31, "h": 31.2, "l": 30, "zjsj": 31, "vol": 2000, "wp": 30, "np": 40, "ccl": 4000},
			},
		},
	}}

	rows, err := GetGlobalFuturesSpot(context.Background(), client, "https://em.test/global", GlobalFuturesSpotOptions{PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	if rows[0].Code != "GC00Y" || rows[0].Price == nil || *rows[0].Price != 2400.5 || rows[1].OpenInterest == nil || *rows[1].OpenInterest != 4000 {
		t.Fatalf("rows = %+v", rows)
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.urls))
	}
	firstURL, _ := url.Parse(client.urls[0])
	secondURL, _ := url.Parse(client.urls[1])
	assertQuery(t, firstURL.Query(), "pageSize", "1")
	assertQuery(t, firstURL.Query(), "pageIndex", "0")
	assertQuery(t, secondURL.Query(), "pageIndex", "1")
}

func TestGetGlobalFuturesSpotStopsOnNonArrayListLikeTypeScript(t *testing.T) {
	client := &fakeFuturesClient{payloads: []map[string]any{
		{
			"total": float64(2),
			"list": []map[string]any{
				{"dm": "GC00Y", "name": "COMEX黄金", "p": 2400.5},
			},
		},
		{
			"total": float64(2),
			"list":  map[string]any{"dm": "SI00Y"},
		},
	}}

	rows, err := GetGlobalFuturesSpot(context.Background(), client, "https://em.test/global", GlobalFuturesSpotOptions{PageSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "GC00Y" {
		t.Fatalf("rows = %+v, want first page only", rows)
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.urls))
	}
}

func TestGetGlobalFuturesKlineBuildsRequestAndSupportsMarketCodeOverride(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"data": map[string]any{
			"code": "HG00Y",
			"name": "COMEX铜",
			"klines": []string{
				"2024-12-16,4.1,4.2,4.3,4.0,123,456,2.5,1.2,0.1,0,0,789,0",
			},
		},
	}}

	rows, err := GetGlobalFuturesKline(context.Background(), client, "HG00Y", "https://em.test/kline", GlobalFuturesKlineOptions{
		Period:    KlinePeriodMonthly,
		StartDate: "20240101",
		EndDate:   "20241231",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "HG00Y" || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 789 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "101.HG00Y")
	assertQuery(t, parsed.Query(), "klt", "103")

	_, err = GetGlobalFuturesKline(context.Background(), client, "ZZZ2507", "https://em.test/kline", GlobalFuturesKlineOptions{MarketCode: 999})
	if err != nil {
		t.Fatal(err)
	}
	parsed, _ = url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "999.ZZZ2507")

	if _, err := GetGlobalFuturesKline(context.Background(), client, "bad2507", "https://em.test/kline", GlobalFuturesKlineOptions{}); err == nil {
		t.Fatal("expected invalid global futures symbol error")
	}
	if _, err := GetGlobalFuturesKline(context.Background(), client, "ZZZ2507", "https://em.test/kline", GlobalFuturesKlineOptions{}); err == nil {
		t.Fatal("expected unknown global futures variety error")
	}
}
