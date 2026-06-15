package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeDragonTigerClient struct {
	urls    []string
	payload map[string]any
}

func (f *fakeDragonTigerClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.urls = append(f.urls, requestURL)
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func (f *fakeDragonTigerClient) lastURL() string {
	if len(f.urls) == 0 {
		return ""
	}
	return f.urls[len(f.urls)-1]
}

func TestGetDragonTigerDetailBuildsFilterAndParsesRows(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":       "600519",
			"SECURITY_NAME_ABBR":  "贵州茅台",
			"TRADE_DATE":          "2024-12-16 00:00:00",
			"CLOSE_PRICE":         1500.0,
			"CHANGE_RATE":         1.5,
			"BILLBOARD_NET_AMT":   1000.0,
			"BILLBOARD_BUY_AMT":   2000.0,
			"BILLBOARD_SELL_AMT":  1000.0,
			"BILLBOARD_DEAL_AMT":  3000.0,
			"ACCUM_AMOUNT":        4000.0,
			"DEAL_NET_RATIO":      5.0,
			"DEAL_AMOUNT_RATIO":   6.0,
			"TURNOVERRATE":        7.0,
			"FREE_MARKET_CAP":     8000.0,
			"EXPLANATION":         "日涨幅偏离值达7%",
			"D1_CLOSE_ADJCHRATE":  1.0,
			"D2_CLOSE_ADJCHRATE":  2.0,
			"D5_CLOSE_ADJCHRATE":  5.0,
			"D10_CLOSE_ADJCHRATE": 10.0,
		},
	})}

	rows, err := GetDragonTigerDetail(context.Background(), client, "https://em.test/datacenter", DragonTigerDateOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" || rows[0].Reason != "日涨幅偏离值达7%" {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL())
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_DAILYBILLBOARD_DETAILSNEW")
	filter := query.Get("filter")
	if !strings.Contains(filter, `TRADE_DATE<='2024-12-31'`) || !strings.Contains(filter, `TRADE_DATE>='2024-12-01'`) {
		t.Fatalf("filter = %q", filter)
	}
}

func TestGetDragonTigerStockStatsUsesPeriod(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "LATEST_TDATE": "2024-12-16", "BILLBOARD_TIMES": 3.0},
	})}

	rows, err := GetDragonTigerStockStats(context.Background(), client, "https://em.test/datacenter", DragonTigerPeriodThreeMonth)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Count == nil || *rows[0].Count != 3 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL())
	assertQuery(t, parsed.Query(), "reportName", "RPT_BILLBOARD_TRADEALL")
	if !strings.Contains(parsed.Query().Get("filter"), `STATISTICS_CYCLE="02"`) {
		t.Fatalf("filter = %q", parsed.Query().Get("filter"))
	}
}

func TestGetDragonTigerInstitutionBuildsFilterAndParsesRows(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "BUY_TIMES": 2.0, "SELL_TIMES": 1.0, "BUY_AMT": 1000.0, "SELL_AMT": 500.0, "NET_AMT": 500.0},
	})}

	rows, err := GetDragonTigerInstitution(context.Background(), client, "https://em.test/datacenter", DragonTigerDateOptions{StartDate: "2024-12-01", EndDate: "2024-12-31"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].OrgNetAmount == nil || *rows[0].OrgNetAmount != 500 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL())
	assertQuery(t, parsed.Query(), "reportName", "RPT_ORGANIZATION_TRADE_DETAILS")
}

func TestGetDragonTigerBranchRankUsesPeriodAndParsesRows(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{"OPERATEDEPT_CODE": "B001", "OPERATEDEPT_NAME": "测试营业部", "TOTAL_BUYAMT": 1000.0, "TOTAL_SELLAMT": 500.0, "TOTAL_BUYER_SALESTIMES": 4.0, "TOTAL_SELLER_SALESTIMES": 3.0, "TOTAL_TIMES": 7.0},
	})}

	rows, err := GetDragonTigerBranchRank(context.Background(), client, "https://em.test/datacenter", DragonTigerPeriodOneYear)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "B001" || rows[0].TotalCount == nil || *rows[0].TotalCount != 7 {
		t.Fatalf("rows = %+v", rows)
	}
	parsed, _ := url.Parse(client.lastURL())
	if !strings.Contains(parsed.Query().Get("filter"), `STATISTICS_CYCLE="04"`) {
		t.Fatalf("filter = %q", parsed.Query().Get("filter"))
	}
}

func TestGetDragonTigerBranchRankDoesNotFallbackFromBlankPrimaryFields(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{
			"OPERATEDEPT_CODE":        "B001",
			"OPERATEDEPT_NAME":        "测试营业部",
			"TOTAL_BUYAMT":            "",
			"BUY_AMT":                 1000.0,
			"TOTAL_SELLAMT":           "-",
			"SELL_AMT":                500.0,
			"TOTAL_BUYER_SALESTIMES":  "",
			"BUY_TIMES":               4.0,
			"TOTAL_SELLER_SALESTIMES": "-",
			"SELL_TIMES":              3.0,
			"TOTAL_TIMES":             7.0,
		},
	})}

	rows, err := GetDragonTigerBranchRank(context.Background(), client, "https://em.test/datacenter", DragonTigerPeriodOneYear)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].TotalBuyAmount != nil || rows[0].TotalSellAmount != nil || rows[0].BuyCount != nil || rows[0].SellCount != nil {
		t.Fatalf("fallback fields = %+v", rows[0])
	}
}

func TestGetDragonTigerBranchRankFallsBackWhenPrimaryFieldsMissing(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{
			"OPERATEDEPT_CODE": "B001",
			"OPERATEDEPT_NAME": "测试营业部",
			"BUY_AMT":          1000.0,
			"SELL_AMT":         500.0,
			"BUY_TIMES":        4.0,
			"SELL_TIMES":       3.0,
		},
	})}

	rows, err := GetDragonTigerBranchRank(context.Background(), client, "https://em.test/datacenter", DragonTigerPeriodOneYear)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].TotalBuyAmount == nil || *rows[0].TotalBuyAmount != 1000 ||
		rows[0].TotalSellAmount == nil || *rows[0].TotalSellAmount != 500 ||
		rows[0].BuyCount == nil || *rows[0].BuyCount != 4 ||
		rows[0].SellCount == nil || *rows[0].SellCount != 3 {
		t.Fatalf("fallback fields = %+v", rows[0])
	}
}

func TestGetDragonTigerStockSeatDetailFetchesBuyAndSell(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{"RANK": 1.0, "OPERATEDEPT_NAME": "测试营业部", "BUY_AMT_REAL": 1000.0, "BUY_RATIO_TOTAL": 10.0, "SELL_AMT_REAL": 100.0, "SELL_RATIO_TOTAL": 1.0, "NET_AMT": 900.0},
	})}

	rows, err := GetDragonTigerStockSeatDetail(context.Background(), client, "https://em.test/datacenter", "sh600519", "20241216")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 || rows[0].Side != "buy" || rows[1].Side != "sell" {
		t.Fatalf("rows = %+v", rows)
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.urls))
	}
	first, _ := url.Parse(client.urls[0])
	second, _ := url.Parse(client.urls[1])
	assertQuery(t, first.Query(), "reportName", "RPT_BILLBOARD_DAILYDETAILSBUY")
	assertQuery(t, second.Query(), "reportName", "RPT_BILLBOARD_DAILYDETAILSSELL")
	if !strings.Contains(first.Query().Get("filter"), `SECURITY_CODE="600519"`) || !strings.Contains(first.Query().Get("filter"), `TRADE_DATE='2024-12-16'`) {
		t.Fatalf("filter = %q", first.Query().Get("filter"))
	}
}

func TestGetDragonTigerStockSeatDetailDoesNotFallbackFromBlankPrimaryFields(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{
			"RANK":             1.0,
			"OPERATEDEPT_NAME": "测试营业部",
			"BUY_AMT_REAL":     "",
			"BUY_AMT":          1000.0,
			"BUY_RATIO_TOTAL":  "-",
			"BUY_AMT_RATIO":    10.0,
			"SELL_AMT_REAL":    "",
			"SELL_AMT":         100.0,
			"SELL_RATIO_TOTAL": "-",
			"SELL_AMT_RATIO":   1.0,
			"NET_AMT":          900.0,
		},
	})}

	rows, err := GetDragonTigerStockSeatDetail(context.Background(), client, "https://em.test/datacenter", "sh600519", "20241216")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %+v", rows)
	}
	first := rows[0]
	if first.BuyAmount != nil || first.BuyAmountRatio != nil || first.SellAmount != nil || first.SellAmountRatio != nil {
		t.Fatalf("fallback fields = %+v", first)
	}
}

func TestGetDragonTigerStockSeatDetailFallsBackWhenPrimaryFieldsMissing(t *testing.T) {
	client := &fakeDragonTigerClient{payload: datacenterPayload([]map[string]any{
		{
			"RANK":             1.0,
			"OPERATEDEPT_NAME": "测试营业部",
			"BUY_AMT":          1000.0,
			"BUY_AMT_RATIO":    10.0,
			"SELL_AMT":         100.0,
			"SELL_AMT_RATIO":   1.0,
			"NET_AMT":          900.0,
		},
	})}

	rows, err := GetDragonTigerStockSeatDetail(context.Background(), client, "https://em.test/datacenter", "sh600519", "20241216")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %+v", rows)
	}
	first := rows[0]
	if first.BuyAmount == nil || *first.BuyAmount != 1000 ||
		first.BuyAmountRatio == nil || *first.BuyAmountRatio != 10 ||
		first.SellAmount == nil || *first.SellAmount != 100 ||
		first.SellAmountRatio == nil || *first.SellAmountRatio != 1 {
		t.Fatalf("fallback fields = %+v", first)
	}
}
