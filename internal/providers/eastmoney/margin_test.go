package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeMarginClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeMarginClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetMarginAccountInfoBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeMarginClient{payload: datacenterPayload([]map[string]any{
		{
			"STATISTICS_DATE":      "2024-12-16 00:00:00",
			"FIN_BALANCE":          1000.0,
			"LOAN_BALANCE":         2000.0,
			"FIN_BUY_AMT":          300.0,
			"LOAN_SELL_AMT":        400.0,
			"OPERATE_INVESTOR_NUM": 50.0,
			"MARGIN_INVESTOR_NUM":  60.0,
			"TOTAL_GUARANTEE":      7000.0,
			"AVG_GUARANTEE_RATIO":  250.0,
		},
	})}

	rows, err := GetMarginAccountInfo(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Date != "2024-12-16" || row.FinBalance == nil || *row.FinBalance != 1000 || row.InvestorCount == nil || *row.InvestorCount != 50 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPTA_WEB_MARGIN_DAILYTRADE")
	assertQuery(t, query, "sortColumns", "STATISTICS_DATE")
	assertQuery(t, query, "sortTypes", "-1")
	assertQuery(t, query, "pageSize", "500")
}

func TestGetMarginAccountInfoDoesNotFallbackFromBlankPrimaryInvestorCount(t *testing.T) {
	client := &fakeMarginClient{payload: datacenterPayload([]map[string]any{
		{
			"STATISTICS_DATE":      "2024-12-16",
			"OPERATE_INVESTOR_NUM": "",
			"INVESTOR_NUM":         50.0,
		},
	})}

	rows, err := GetMarginAccountInfo(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].InvestorCount != nil {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestGetMarginAccountInfoFallsBackWhenPrimaryInvestorCountMissing(t *testing.T) {
	client := &fakeMarginClient{payload: datacenterPayload([]map[string]any{
		{
			"STATISTICS_DATE": "2024-12-16",
			"INVESTOR_NUM":    50.0,
		},
	})}

	rows, err := GetMarginAccountInfo(context.Background(), client, "https://em.test/datacenter")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].InvestorCount == nil || *rows[0].InvestorCount != 50 {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestGetMarginTargetListBuildsDateFilterAndParsesRows(t *testing.T) {
	client := &fakeMarginClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"TRADE_DATE":         "20241216",
			"FIN_BALANCE":        1000.0,
			"FIN_BUY_AMT":        200.0,
			"FIN_REPAY_AMT":      300.0,
			"LOAN_BALANCE":       400.0,
			"LOAN_SELL_VOLUME":   500.0,
			"LOAN_REPAY_VOLUME":  600.0,
		},
	})}

	rows, err := GetMarginTargetList(context.Background(), client, "https://em.test/datacenter", "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.Date != "2024-12-16" || row.LoanRepayVolume == nil || *row.LoanRepayVolume != 600 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_MARGIN_TRADE_DETAIL")
	assertQuery(t, query, "sortColumns", "FIN_BALANCE")
	assertQuery(t, query, "pageSize", "5000")
	if filter := query.Get("filter"); !strings.Contains(filter, `TRADE_DATE='2024-12-16'`) {
		t.Fatalf("filter = %q", filter)
	}
}
