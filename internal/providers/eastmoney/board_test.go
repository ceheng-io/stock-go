package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeBoardClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeBoardClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	b, _ := json.Marshal(f.payload)
	return json.Unmarshal(b, target)
}

type fakeBoardPagesClient struct {
	urls     []string
	payloads []map[string]any
}

func (f *fakeBoardPagesClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.urls = append(f.urls, requestURL)
	index := len(f.urls) - 1
	payload := map[string]any{"data": map[string]any{"total": 0, "diff": []map[string]any{}}}
	if index < len(f.payloads) {
		payload = f.payloads[index]
	}
	b, _ := json.Marshal(payload)
	return json.Unmarshal(b, target)
}

func TestGetIndustryListBuildsRequestAndParsesSortedRows(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"total": float64(2),
			"diff": []map[string]any{
				{"f12": "BK0001", "f14": "低涨幅行业", "f2": 10.5, "f3": 1.2, "f4": 0.12, "f8": 2.1, "f20": 100000.0, "f104": 12.0, "f105": 3.0, "f128": "领涨A", "f136": 4.5},
				{"f12": "BK0002", "f14": "高涨幅行业", "f2": 20.5, "f3": 3.4, "f4": 0.34, "f8": 5.6, "f20": 200000.0, "f104": 20.0, "f105": 2.0, "f128": "领涨B", "f136": 7.8},
			},
		},
	}}

	rows, err := GetIndustryList(context.Background(), client, "https://em.test/list")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	if rows[0].Code != "BK0002" || rows[0].Rank != 1 || rows[0].LeadingStock == nil || *rows[0].LeadingStock != "领涨B" {
		t.Fatalf("rows[0] = %+v", rows[0])
	}
	if rows[0].ChangePercent == nil || *rows[0].ChangePercent != 3.4 {
		t.Fatalf("change percent = %+v", rows[0])
	}
	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "fs", "m:90 t:2 f:!50")
	assertQuery(t, query, "fid", "f3")
	if !strings.Contains(query.Get("fields"), "f136") {
		t.Fatalf("fields = %q, want f136 included", query.Get("fields"))
	}
}

func TestGetIndustryListFetchesAllPages(t *testing.T) {
	client := &fakeBoardPagesClient{payloads: []map[string]any{
		boardListPayloadWithTotal(101, boardRows(100)),
		boardListPayloadWithTotal(101, []map[string]any{{"f12": "BK9999", "f14": "最后行业", "f3": 9.9}}),
	}}

	rows, err := GetIndustryList(context.Background(), client, "https://em.test/list")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 101 {
		t.Fatalf("len(rows) = %d, want 101", len(rows))
	}
	if len(client.urls) != 2 {
		t.Fatalf("len(urls) = %d, want 2 (%v)", len(client.urls), client.urls)
	}
	firstURL, _ := url.Parse(client.urls[0])
	secondURL, _ := url.Parse(client.urls[1])
	assertQuery(t, firstURL.Query(), "pn", "1")
	assertQuery(t, secondURL.Query(), "pn", "2")
	assertQuery(t, firstURL.Query(), "pz", "100")
}

func TestGetIndustryListStopsOnNonArrayDiffLikeTypeScript(t *testing.T) {
	client := &fakeBoardPagesClient{payloads: []map[string]any{
		boardListPayloadWithTotal(101, boardRows(100)),
		{"data": map[string]any{
			"total": 101,
			"diff":  map[string]any{"f12": "BK9999"},
		}},
	}}

	rows, err := GetIndustryList(context.Background(), client, "https://em.test/list")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 100 {
		t.Fatalf("len(rows) = %d, want first page only", len(rows))
	}
	if len(client.urls) != 2 {
		t.Fatalf("len(urls) = %d, want 2 (%v)", len(client.urls), client.urls)
	}
}

func TestGetIndustryListKeepsNilLeadingStock(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{
				{"f12": "BK0001", "f14": "酿酒行业", "f3": 1.2, "f128": ""},
			},
		},
	}}

	rows, err := GetIndustryList(context.Background(), client, "https://em.test/list")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].LeadingStock != nil {
		t.Fatalf("rows[0] = %+v", rows[0])
	}
}

func TestGetConceptListUsesConceptFilter(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{
				{"f12": "BK1001", "f14": "人工智能", "f3": 2.5},
			},
		},
	}}

	rows, err := GetConceptList(context.Background(), client, "https://em.test/concept")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "BK1001" || rows[0].Name != "人工智能" {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "fs", "m:90 t:3 f:!50")
	assertQuery(t, query, "fid", "f12")
}

func TestGetBoardSpotBuildsRequestAndScalesFields(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"f43":  1050.0,
			"f44":  1100.0,
			"f45":  1000.0,
			"f46":  1010.0,
			"f47":  123456.0,
			"f48":  987654.0,
			"f170": 250.0,
			"f171": 310.0,
			"f168": 120.0,
			"f169": 50.0,
		},
	}}

	rows, err := GetBoardSpot(context.Background(), client, "BK0001", "https://em.test/spot")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 10 {
		t.Fatalf("len(rows) = %d, want 10", len(rows))
	}
	if rows[0].Item != "最新" || rows[0].Value == nil || *rows[0].Value != 10.5 {
		t.Fatalf("latest = %+v", rows[0])
	}
	if rows[4].Item != "成交量" || rows[4].Value == nil || *rows[4].Value != 123456 {
		t.Fatalf("volume = %+v", rows[4])
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "90.BK0001")
}

func TestGetBoardSpotReturnsEmptyWhenDataMissing(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{}}

	rows, err := GetBoardSpot(context.Background(), client, "BK0001", "https://em.test/spot")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("len(rows) = %d, want 0: %+v", len(rows), rows)
	}
}

func TestGetBoardConstituentsBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{
				{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0, "f3": 1.1, "f4": 16.0, "f5": 100.0, "f6": 200.0, "f7": 2.2, "f15": 1510.0, "f16": 1490.0, "f17": 1501.0, "f18": 1484.0, "f8": 0.5, "f9": 30.0, "f23": 8.0},
			},
		},
	}}

	rows, err := GetBoardConstituents(context.Background(), client, "BK0001", "https://em.test/constituents")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.Rank != 1 {
		t.Fatalf("row = %+v", row)
	}
	if row.Price == nil || *row.Price != 1500 || row.PB == nil || *row.PB != 8 {
		t.Fatalf("numeric row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "fs", "b:BK0001 f:!50")
	assertQuery(t, parsed.Query(), "fid", "f3")
}

func TestGetBoardConstituentsFetchesAllPages(t *testing.T) {
	client := &fakeBoardPagesClient{payloads: []map[string]any{
		boardListPayloadWithTotal(101, constituentRows(100)),
		boardListPayloadWithTotal(101, []map[string]any{{"f12": "000002", "f14": "万科A", "f3": 1.9}}),
	}}

	rows, err := GetBoardConstituents(context.Background(), client, "BK0001", "https://em.test/constituents")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 101 {
		t.Fatalf("len(rows) = %d, want 101", len(rows))
	}
	if len(client.urls) != 2 {
		t.Fatalf("len(urls) = %d, want 2 (%v)", len(client.urls), client.urls)
	}
	firstURL, _ := url.Parse(client.urls[0])
	secondURL, _ := url.Parse(client.urls[1])
	assertQuery(t, firstURL.Query(), "pn", "1")
	assertQuery(t, secondURL.Query(), "pn", "2")
	assertQuery(t, firstURL.Query(), "pz", "100")
}

func TestGetBoardKlineBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"klines": []string{
				"2024-12-16,100,105,106,99,12345,67890000,3.5,2.0,2,1.2",
			},
		},
	}}

	rows, err := GetBoardKline(context.Background(), client, "BK0001", "https://em.test/kline", HistoryKlineOptions{
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
	if row.Date != "2024-12-16" || row.Close == nil || *row.Close != 105 || row.TurnoverRate == nil || *row.TurnoverRate != 1.2 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "secid", "90.BK0001")
	assertQuery(t, query, "klt", "102")
	assertQuery(t, query, "fqt", "2")
	assertQuery(t, query, "smplmt", "10000")
	assertQuery(t, query, "lmt", "1000000")
}

func TestGetBoardKlineDefaultsToNoAdjustment(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{"klines": []string{}},
	}}

	if _, err := GetBoardKline(context.Background(), client, "BK0001", "https://em.test/kline", HistoryKlineOptions{}); err != nil {
		t.Fatal(err)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "klt", "101")
	assertQuery(t, parsed.Query(), "fqt", "0")
}

func TestGetBoardMinuteKlinePeriodOneParsesTimeline(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"trends": []string{
				"2024-12-16 09:31,100,101,102,99,1000,2000,101.5",
			},
		},
	}}

	rows, err := GetBoardMinuteKline(context.Background(), client, "BK0001", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 1 {
		t.Fatalf("timeline = %+v", rows.Timeline)
	}
	row := rows.Timeline[0]
	if row.Time != "2024-12-16 09:31" || row.Price == nil || *row.Price != 101.5 || row.Close == nil || *row.Close != 101 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "90.BK0001")
	assertQuery(t, parsed.Query(), "ndays", "1")
}

func TestGetBoardMinuteKlinePeriodOneReturnsEmptyRowsForNonArrayTrends(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{"trends": map[string]any{"time": "2024-12-16 09:31"}},
	}}

	rows, err := GetBoardMinuteKline(context.Background(), client, "BK0001", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Timeline) != 0 {
		t.Fatalf("timeline = %+v, want empty", rows.Timeline)
	}
}

func TestGetBoardMinuteKlinePeriodFiveParsesKlines(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{
			"klines": []string{
				"2024-12-16 09:35,100,105,106,99,12345,67890000,3.5,2.0,2,1.2",
			},
		},
	}}

	rows, err := GetBoardMinuteKline(context.Background(), client, "BK0001", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 1 {
		t.Fatalf("klines = %+v", rows.Klines)
	}
	row := rows.Klines[0]
	if row.Time != "2024-12-16 09:35" || row.Close == nil || *row.Close != 105 || row.Change == nil || *row.Change != 2 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "klt", "5")
	assertQuery(t, query, "fqt", "1")
	assertQuery(t, query, "beg", "0")
	assertQuery(t, query, "end", "20500101")
}

func TestGetBoardMinuteKlinePeriodFiveReturnsEmptyRowsForNonArrayKlines(t *testing.T) {
	client := &fakeBoardClient{payload: map[string]any{
		"data": map[string]any{"klines": map[string]any{"time": "2024-12-16 09:35"}},
	}}

	rows, err := GetBoardMinuteKline(context.Background(), client, "BK0001", "https://em.test/kline", "https://em.test/trends", MinuteKlineOptions{Period: MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 0 {
		t.Fatalf("klines = %+v, want empty", rows.Klines)
	}
}

func boardListPayloadWithTotal(total int, diff []map[string]any) map[string]any {
	return map[string]any{"data": map[string]any{"total": total, "diff": diff}}
}

func boardRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, map[string]any{"f12": "BK0001", "f14": "酿酒行业", "f3": 1.2})
	}
	return rows
}

func constituentRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, map[string]any{"f12": "600519", "f14": "贵州茅台", "f3": 1.1})
	}
	return rows
}
