package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeNorthboundClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeNorthboundClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetNorthboundMinuteParsesNorthAndSouth(t *testing.T) {
	client := &fakeNorthboundClient{payload: map[string]any{
		"data": map[string]any{
			"s2nDate": "20241216",
			"s2n":     []string{"09:31,100,0,200,0,300"},
			"n2sDate": "2024-12-17",
			"n2s":     []string{"09:32,400,0,500,0,900"},
		},
	}}

	rows, err := GetNorthboundMinute(context.Background(), client, "https://em.test/minute", NorthboundSouth)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Date != "2024-12-17" || row.Time != "09:32" || row.TotalNetInflow == nil || *row.TotalNetInflow != 900 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	if !strings.Contains(parsed.Query().Get("fields2"), "f62") {
		t.Fatalf("fields2 = %q", parsed.Query().Get("fields2"))
	}
}

func TestGetNorthboundFlowSummaryBuildsDatacenterRequestAndParsesRows(t *testing.T) {
	client := &fakeNorthboundClient{payload: map[string]any{
		"result": map[string]any{
			"pages": float64(3),
			"count": float64(3),
			"data": []map[string]any{
				{
					"TRADE_DATE":       "2024-12-16 00:00:00",
					"MUTUAL_TYPE":      "001",
					"MUTUAL_TYPE_NAME": "沪股通",
					"FUNDS_DIRECTION":  "北向资金",
					"status":           "1",
					"netBuyAmt":        1000.0,
					"dayNetAmtIn":      2000.0,
					"dayAmtRemain":     3000.0,
					"f104":             10.0,
					"f106":             2.0,
					"f105":             5.0,
					"INDEX_CODE":       "000001",
					"INDEX_NAME":       "上证指数",
					"INDEX_f3":         1.2,
				},
			},
		},
	}}

	rows, err := GetNorthboundFlowSummary(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Date != "2024-12-16" || row.BoardName != "沪股通" || row.NetBuyAmount == nil || *row.NetBuyAmount != 1000 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_MUTUAL_QUOTA")
	assertQuery(t, query, "sortColumns", "MUTUAL_TYPE")
	assertQuery(t, query, "pageNumber", "1")
}

func TestGetNorthboundHoldingRankBuildsFiltersAndParsesRows(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":                "2024-12-16",
			"SECURITY_CODE":             "600519",
			"SECURITY_NAME":             "贵州茅台",
			"CLOSE_PRICE":               1500.0,
			"CHANGE_RATE":               1.5,
			"HOLD_SHARES":               100.0,
			"HOLD_MARKET_CAP":           200.0,
			"HOLD_RATIO":                3.0,
			"A_SHARES_RATIO":            4.0,
			"ADD_SHARES":                5.0,
			"ADD_MARKET_CAP":            6.0,
			"ADD_MARKET_CAP_PROPORTION": 7.0,
			"BOARD_NAME":                "白酒",
		},
	})}

	rows, err := GetNorthboundHoldingRank(context.Background(), client, "https://em.test/datacenter", NorthboundHoldingRankOptions{
		Market: NorthboundMarketShanghai,
		Period: NorthboundRankThreeDay,
		Date:   "2024-12-16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" || rows[0].AddMarketValuePercent == nil || *rows[0].AddMarketValuePercent != 7 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	filter := parsed.Query().Get("filter")
	if !strings.Contains(filter, `INTERVAL_TYPE="3"`) || !strings.Contains(filter, `MUTUAL_TYPE="001"`) || !strings.Contains(filter, `TRADE_DATE='2024-12-16'`) {
		t.Fatalf("filter = %q", filter)
	}
}

func TestGetNorthboundHoldingRankDoesNotFallbackFromBlankSecurityName(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":         "2024-12-16",
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME":      "",
			"SECURITY_NAME_ABBR": "贵州茅台",
		},
	})}

	rows, err := GetNorthboundHoldingRank(context.Background(), client, "https://em.test/datacenter", NorthboundHoldingRankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].Name != "" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestGetNorthboundHoldingRankFallsBackWhenSecurityNameMissing(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":         "2024-12-16",
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
		},
	})}

	rows, err := GetNorthboundHoldingRank(context.Background(), client, "https://em.test/datacenter", NorthboundHoldingRankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].Name != "贵州茅台" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestGetNorthboundHistoryBuildsFiltersAndParsesRows(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":       "2024-12-16",
			"NET_DEAL_AMT":     1000.0,
			"BUY_AMT":          2000.0,
			"SELL_AMT":         3000.0,
			"ACCUM_DEAL_AMT":   4000.0,
			"FUND_INFLOW":      5000.0,
			"QUOTA_BALANCE":    6000.0,
			"LEAD_STOCKS_CODE": "600519",
			"LEAD_STOCKS_NAME": "贵州茅台",
			"LS_CHANGE_RATE":   1.5,
		},
	})}

	rows, err := GetNorthboundHistory(context.Background(), client, "https://em.test/datacenter", NorthboundSouth, NorthboundHistoryOptions{
		StartDate: "2024-12-01",
		EndDate:   "2024-12-31",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].TopStockCode == nil || *rows[0].TopStockCode != "600519" || rows[0].TopStockName == nil || *rows[0].TopStockName != "贵州茅台" || rows[0].NetInflow == nil || *rows[0].NetInflow != 5000 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	filter := parsed.Query().Get("filter")
	if !strings.Contains(filter, `BOARD_TYPE="0"`) || !strings.Contains(filter, `TRADE_DATE>='2024-12-01'`) || !strings.Contains(filter, `TRADE_DATE<='2024-12-31'`) {
		t.Fatalf("filter = %q", filter)
	}
}

func TestGetNorthboundHistoryKeepsNilLeadStock(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":       "2024-12-16",
			"LEAD_STOCKS_CODE": "",
			"LEAD_STOCKS_NAME": nil,
			"FUND_INFLOW":      5000.0,
		},
	})}

	rows, err := GetNorthboundHistory(context.Background(), client, "https://em.test/datacenter", NorthboundSouth, NorthboundHistoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.TopStockCode != nil || row.TopStockName != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetNorthboundIndividualStripsMarketPrefix(t *testing.T) {
	client := &fakeNorthboundClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":      "2024-12-16",
			"HOLD_SHARES":     100.0,
			"HOLD_MARKET_CAP": 200.0,
			"HOLD_RATIO":      3.0,
			"A_SHARES_RATIO":  4.0,
			"CLOSE_PRICE":     1500.0,
			"CHANGE_RATE":     1.5,
		},
	})}

	rows, err := GetNorthboundIndividual(context.Background(), client, "https://em.test/datacenter", "sh600519", NorthboundHistoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].HoldShares == nil || *rows[0].HoldShares != 100 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL)
	if !strings.Contains(parsed.Query().Get("filter"), `SECURITY_CODE="600519"`) {
		t.Fatalf("filter = %q", parsed.Query().Get("filter"))
	}
}

func datacenterPayload(rows []map[string]any) map[string]any {
	return map[string]any{"result": map[string]any{"data": rows}}
}
