package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeFundFlowClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeFundFlowClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

type fakeFundFlowPagesClient struct {
	urls     []string
	payloads []map[string]any
}

func (f *fakeFundFlowPagesClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.urls = append(f.urls, requestURL)
	index := len(f.urls) - 1
	payload := map[string]any{"data": map[string]any{"total": 0, "diff": []map[string]any{}}}
	if index < len(f.payloads) {
		payload = f.payloads[index]
	}
	body, _ := json.Marshal(payload)
	return json.Unmarshal(body, target)
}

func TestGetIndividualFundFlowBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFundFlowClient{payload: fundFlowKlinesPayload([]string{
		"2024-12-16,1000,100,200,300,400,10,1,2,3,4,1500,1.5",
	})}

	rows, err := GetIndividualFundFlow(context.Background(), client, "sh600519", "https://em.test/fflow", FundFlowOptions{Period: FundFlowPeriodWeekly})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Date != "2024-12-16" || row.MainNetInflow == nil || *row.MainNetInflow != 1000 || row.Close == nil || *row.Close != 1500 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "secid", "1.600519")
	assertQuery(t, query, "klt", "102")
	assertQuery(t, query, "lmt", "0")
}

func TestGetIndividualFundFlowReturnsEmptyRowsForNonArrayKlines(t *testing.T) {
	client := &fakeFundFlowClient{payload: map[string]any{
		"data": map[string]any{"klines": map[string]any{"date": "2024-12-16"}},
	}}

	rows, err := GetIndividualFundFlow(context.Background(), client, "sh600519", "https://em.test/fflow", FundFlowOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}

func TestGetMarketFundFlowBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFundFlowClient{payload: fundFlowKlinesPayload([]string{
		"2024-12-16,1000,100,200,300,400,10,1,2,3,4,3500,1.2,11000,1.4",
	})}

	rows, err := GetMarketFundFlow(context.Background(), client, "https://em.test/fflow")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.SHClose == nil || *row.SHClose != 3500 || row.SZChangePercent == nil || *row.SZChangePercent != 1.4 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "secid", "1.000001")
	assertQuery(t, query, "secid2", "0.399001")
}

func TestGetFundFlowRankBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFundFlowClient{payload: clistPayload([]map[string]any{
		{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0, "f109": 5.5, "f164": 1000.0, "f165": 10.0, "f166": 400.0, "f167": 4.0, "f168": 300.0, "f169": 3.0, "f170": 200.0, "f171": 2.0, "f172": 100.0, "f173": 1.0},
	})}

	rows, err := GetFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{Indicator: FundFlowRankFiveDay})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.MainNetInflow == nil || *row.MainNetInflow != 1000 || row.ChangePercent == nil || *row.ChangePercent != 5.5 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "fid", "f164")
	if query.Get("fs") == "" {
		t.Fatalf("expected stock fs in query: %s", query.Encode())
	}
}

func TestGetFundFlowRankFetchesAllPages(t *testing.T) {
	client := &fakeFundFlowPagesClient{payloads: []map[string]any{
		clistPayloadWithTotal(101, stockFundFlowRankRows(100)),
		clistPayloadWithTotal(101, []map[string]any{
			{"f12": "000002", "f14": "万科A", "f2": 8.0, "f62": 300.0},
		}),
	}}

	rows, err := GetFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{})
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

func TestGetFundFlowRankStopsOnNonArrayDiffLikeTypeScript(t *testing.T) {
	client := &fakeFundFlowPagesClient{payloads: []map[string]any{
		clistPayloadWithTotal(101, stockFundFlowRankRows(100)),
		{"data": map[string]any{
			"total": 101,
			"diff":  map[string]any{"f12": "000002"},
		}},
	}}

	rows, err := GetFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{})
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

func TestGetFundFlowRankStopsAfterFirstPageWhenTotalMissing(t *testing.T) {
	client := &fakeFundFlowPagesClient{payloads: []map[string]any{
		clistPayloadWithTotal(0, stockFundFlowRankRows(100)),
	}}

	rows, err := GetFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 100 {
		t.Fatalf("len(rows) = %d, want 100", len(rows))
	}
	if len(client.urls) != 1 {
		t.Fatalf("len(urls) = %d, want 1 (%v)", len(client.urls), client.urls)
	}
}

func TestGetSectorFundFlowRankBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFundFlowClient{payload: clistPayload([]map[string]any{
		{"f12": "BK1001", "f14": "人工智能", "f3": 3.2, "f62": 2000.0, "f184": 12.0, "f66": 800.0, "f72": 600.0, "f78": 400.0, "f84": 200.0, "f204": "600519", "f205": "贵州茅台"},
	})}

	rows, err := GetSectorFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{SectorType: FundFlowSectorConcept})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "BK1001" || row.TopStockCode == nil || *row.TopStockCode != "600519" || row.MainNetInflow == nil || *row.MainNetInflow != 2000 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "fs", "m:90+t:3")
	if !strings.HasSuffix(query.Get("fields"), "f204,f205") {
		t.Fatalf("fields = %q, want f204,f205 suffix", query.Get("fields"))
	}
}

func TestGetSectorFundFlowRankFetchesAllPages(t *testing.T) {
	client := &fakeFundFlowPagesClient{payloads: []map[string]any{
		clistPayloadWithTotal(101, sectorFundFlowRankRows(100)),
		clistPayloadWithTotal(101, []map[string]any{
			{"f12": "BK1003", "f14": "半导体", "f3": 1.2, "f62": 500.0, "f204": "000002", "f205": "万科A"},
		}),
	}}

	rows, err := GetSectorFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{SectorType: FundFlowSectorConcept})
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

func TestGetSectorFundFlowRankKeepsNilTopStockWhenFieldsMissing(t *testing.T) {
	client := &fakeFundFlowClient{payload: clistPayload([]map[string]any{
		{"f12": "BK1001", "f14": "人工智能", "f3": 3.2, "f62": 2000.0},
		{"f12": "BK1002", "f14": "机器人", "f3": 2.2, "f62": 1000.0, "f204": "", "f205": ""},
	})}

	rows, err := GetSectorFundFlowRank(context.Background(), client, "https://em.test/clist", FundFlowRankOptions{SectorType: FundFlowSectorConcept})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("len(rows) = %d, want 2", len(rows))
	}
	for _, row := range rows {
		if row.TopStockCode != nil || row.TopStockName != nil {
			t.Fatalf("row = %+v, want nil top stock fields", row)
		}
	}
}

func TestGetSectorFundFlowHistoryBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFundFlowClient{payload: fundFlowKlinesPayload([]string{
		"2024-12-16,1000,100,200,300,400,10,1,2,3,4,1500,1.5",
	})}

	rows, err := GetSectorFundFlowHistory(context.Background(), client, "BK1001", "https://em.test/fflow", FundFlowOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Date != "2024-12-16" {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	assertQuery(t, parsed.Query(), "secid", "90.BK1001")
}

func fundFlowKlinesPayload(klines []string) map[string]any {
	return map[string]any{"data": map[string]any{"klines": klines}}
}

func clistPayload(diff []map[string]any) map[string]any {
	return map[string]any{"data": map[string]any{"total": len(diff), "diff": diff}}
}

func clistPayloadWithTotal(total int, diff []map[string]any) map[string]any {
	return map[string]any{"data": map[string]any{"total": total, "diff": diff}}
}

func stockFundFlowRankRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, map[string]any{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0, "f62": 1000.0})
	}
	return rows
}

func sectorFundFlowRankRows(count int) []map[string]any {
	rows := make([]map[string]any, 0, count)
	for i := 0; i < count; i++ {
		rows = append(rows, map[string]any{"f12": "BK1001", "f14": "人工智能", "f3": 3.2, "f62": 2000.0, "f204": "600519", "f205": "贵州茅台"})
	}
	return rows
}
