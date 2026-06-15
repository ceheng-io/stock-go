package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeBlockTradeClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeBlockTradeClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetBlockTradeMarketStatBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":        "2024-12-16 00:00:00",
			"CLOSE_PRICE":       3500.0,
			"CHANGE_RATE":       1.5,
			"TURNOVER":          1000000.0,
			"PREMIUM_TURNOVER":  600000.0,
			"PREMIUM_RATIO":     60.0,
			"DISCOUNT_TURNOVER": 400000.0,
			"DISCOUNT_RATIO":    40.0,
		},
	})}

	rows, err := GetBlockTradeMarketStat(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Date != "2024-12-16" || row.SHClose == nil || *row.SHClose != 3500 || row.TotalAmount == nil || *row.TotalAmount != 1000000 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "PRT_BLOCKTRADE_MARKET_STA")
	assertQuery(t, query, "sortColumns", "TRADE_DATE")
	assertQuery(t, query, "pageSize", "500")
}

func TestGetBlockTradeMarketStatDoesNotFallbackFromBlankPrimaryFields(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":       "2024-12-16",
			"CLOSE_PRICE":      "",
			"SH_CLOSE_PRICE":   3500.0,
			"CHANGE_RATE":      "-",
			"SH_CHANGE_RATE":   1.5,
			"TURNOVER":         "",
			"TOTAL_AMOUNT":     1000000.0,
			"PREMIUM_TURNOVER": "-",
			"PREMIUM_AMOUNT":   600000.0,
		},
	})}

	rows, err := GetBlockTradeMarketStat(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.SHClose != nil || row.SHChangePercent != nil || row.TotalAmount != nil || row.PremiumAmount != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetBlockTradeMarketStatFallsBackWhenPrimaryFieldsMissing(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"TRADE_DATE":      "2024-12-16",
			"SH_CLOSE_PRICE":  3500.0,
			"SH_CHANGE_RATE":  1.5,
			"TOTAL_AMOUNT":    1000000.0,
			"PREMIUM_AMOUNT":  600000.0,
			"DISCOUNT_AMOUNT": 400000.0,
		},
	})}

	rows, err := GetBlockTradeMarketStat(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.SHClose == nil || *row.SHClose != 3500 ||
		row.SHChangePercent == nil || *row.SHChangePercent != 1.5 ||
		row.TotalAmount == nil || *row.TotalAmount != 1000000 ||
		row.PremiumAmount == nil || *row.PremiumAmount != 600000 ||
		row.DiscountAmount == nil || *row.DiscountAmount != 400000 {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetBlockTradeDetailBuildsDateFilterAndParsesRows(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":           "600519",
			"SECURITY_NAME_ABBR":      "贵州茅台",
			"TRADE_DATE":              "20241216",
			"CLOSE_PRICE":             1500.0,
			"CHANGE_RATE":             1.5,
			"PRICE":                   1490.0,
			"VOLUME":                  100.0,
			"TURNOVER":                149000.0,
			"PREMIUM_RATE":            -0.67,
			"BUYER_OPERATEDEPT_NAME":  "买方营业部",
			"SELLER_OPERATEDEPT_NAME": "卖方营业部",
		},
	})}

	rows, err := GetBlockTradeDetail(context.Background(), client, "https://em.test/datacenter", BlockTradeDateOptions{
		StartDate: "20241201",
		EndDate:   "2024-12-31",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.Date != "2024-12-16" || row.DealPrice == nil || *row.DealPrice != 1490 || row.BuyBranch != "买方营业部" {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_BLOCK_TRADE_DETAIL")
	assertQuery(t, query, "sortColumns", "TRADE_DATE,SECURITY_CODE")
	filter := query.Get("filter")
	if !strings.Contains(filter, `TRADE_DATE>='2024-12-01'`) || !strings.Contains(filter, `TRADE_DATE<='2024-12-31'`) {
		t.Fatalf("filter = %q", filter)
	}
}

func TestGetBlockTradeDetailDoesNotFallbackFromBlankPrimaryFields(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":           "600519",
			"SECURITY_NAME_ABBR":      "贵州茅台",
			"TRADE_DATE":              "20241216",
			"DEAL_PRICE":              "",
			"PRICE":                   1490.0,
			"DEAL_VOLUME":             "-",
			"VOLUME":                  100.0,
			"BUYER_DEPT":              "",
			"BUYER_OPERATEDEPT_NAME":  "买方营业部",
			"SELLER_DEPT":             "",
			"SELLER_OPERATEDEPT_NAME": "卖方营业部",
		},
	})}

	rows, err := GetBlockTradeDetail(context.Background(), client, "https://em.test/datacenter", BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.DealPrice != nil || row.DealVolume != nil || row.BuyBranch != "" || row.SellBranch != "" {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetBlockTradeDetailFallsBackWhenPrimaryFieldsMissing(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":           "600519",
			"SECURITY_NAME_ABBR":      "贵州茅台",
			"TRADE_DATE":              "20241216",
			"PRICE":                   1490.0,
			"VOLUME":                  100.0,
			"TURNOVER":                149000.0,
			"PREMIUM_RATE":            -0.67,
			"BUYER_OPERATEDEPT_NAME":  "买方营业部",
			"SELLER_OPERATEDEPT_NAME": "卖方营业部",
		},
	})}

	rows, err := GetBlockTradeDetail(context.Background(), client, "https://em.test/datacenter", BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.DealPrice == nil || *row.DealPrice != 1490 ||
		row.DealVolume == nil || *row.DealVolume != 100 ||
		row.DealAmount == nil || *row.DealAmount != 149000 ||
		row.PremiumRate == nil || *row.PremiumRate != -0.67 ||
		row.BuyBranch != "买方营业部" || row.SellBranch != "卖方营业部" {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetBlockTradeDailyStatBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeBlockTradeClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"TRADE_DATE":         "2024-12-16",
			"CHANGE_RATE":        1.5,
			"CLOSE_PRICE":        1500.0,
			"DEAL_COUNT":         2.0,
			"TOTAL_AMOUNT":       300000.0,
			"TOTAL_VOLUME":       200.0,
			"PREMIUM_AMOUNT":     100000.0,
			"DISCOUNT_AMOUNT":    200000.0,
		},
	})}

	rows, err := GetBlockTradeDailyStat(context.Background(), client, "https://em.test/datacenter", BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.DealCount == nil || *row.DealCount != 2 || row.DiscountAmount == nil || *row.DiscountAmount != 200000 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_BLOCK_TRADE_STA")
	assertQuery(t, query, "sortColumns", "TRADE_DATE,DEAL_AMT")
	assertQuery(t, query, "sortTypes", "-1,-1")
	if filter := query.Get("filter"); filter != "" {
		t.Fatalf("filter = %q, want empty", filter)
	}
}
