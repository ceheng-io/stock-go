package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
	"testing"
)

type fakeDividendClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeDividendClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetDividendDetailBuildsFilterAndParsesRows(t *testing.T) {
	client := &fakeDividendClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":       "600519",
			"SECURITY_NAME_ABBR":  "贵州茅台",
			"REPORT_DATE":         "2024-12-31 00:00:00",
			"PLAN_NOTICE_DATE":    "2025-03-31T00:00:00.000",
			"PUBLISH_DATE":        "2025-04-01 00:00:00",
			"BONUS_IT_RATIO":      1.0,
			"BONUS_RATIO":         0.2,
			"IT_RATIO":            0.8,
			"PRETAX_BONUS_RMB":    30.0,
			"IMPL_PLAN_PROFILE":   "10派30元",
			"DIVIDENT_RATIO":      2.5,
			"BASIC_EPS":           50.0,
			"BVPS":                200.0,
			"PER_CAPITAL_RESERVE": 10.0,
			"PER_UNASSIGN_PROFIT": 20.0,
			"PNP_YOY_RATIO":       3.5,
			"TOTAL_SHARES":        1000000.0,
			"EQUITY_RECORD_DATE":  "2025-06-20",
			"EX_DIVIDEND_DATE":    "2025-06-23",
			"PAY_DATE":            "2025-06-24",
			"ASSIGN_PROGRESS":     "实施分配",
			"NOTICE_DATE":         "2025-06-10",
		},
	})}

	rows, err := GetDividendDetail(context.Background(), client, "https://em.test/datacenter", "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.ReportDate == nil || *row.ReportDate != "2024-12-31" || row.PlanNoticeDate == nil || *row.PlanNoticeDate != "2025-03-31" || row.DividendPretax == nil || *row.DividendPretax != 30 {
		t.Fatalf("row = %+v", row)
	}
	if row.DividendDesc == nil || *row.DividendDesc != "10派30元" || row.AssignProgress == nil || *row.AssignProgress != "实施分配" || row.PayDate == nil || *row.PayDate != "2025-06-24" {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_SHAREBONUS_DET")
	assertQuery(t, query, "columns", "ALL")
	assertQuery(t, query, "sortColumns", "REPORT_DATE")
	assertQuery(t, query, "sortTypes", "-1")
	assertQuery(t, query, "pageSize", "500")
	if filter := query.Get("filter"); !strings.Contains(filter, `SECURITY_CODE="600519"`) {
		t.Fatalf("filter = %q", filter)
	}
}

func TestGetDividendDetailUsesNilForBlankNullableDateStrings(t *testing.T) {
	client := &fakeDividendClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"REPORT_DATE":        "",
			"PLAN_NOTICE_DATE":   nil,
			"PUBLISH_DATE":       "",
			"EQUITY_RECORD_DATE": "",
			"EX_DIVIDEND_DATE":   nil,
			"PAY_DATE":           "",
			"NOTICE_DATE":        "",
			"PRETAX_BONUS_RMB":   nil,
		},
	})}

	rows, err := GetDividendDetail(context.Background(), client, "https://em.test/datacenter", "600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.ReportDate != nil || row.PlanNoticeDate != nil || row.DisclosureDate != nil || row.EquityRecordDate != nil || row.ExDividendDate != nil || row.PayDate != nil || row.NoticeDate != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetDividendDetailKeepsBlankTextFieldsLikeTS(t *testing.T) {
	client := &fakeDividendClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"IMPL_PLAN_PROFILE":  "",
			"ASSIGN_PROGRESS":    "",
		},
	})}

	rows, err := GetDividendDetail(context.Background(), client, "https://em.test/datacenter", "600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.DividendDesc == nil || *row.DividendDesc != "" || row.AssignProgress == nil || *row.AssignProgress != "" {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetDividendDetailDoesNotFallbackFromBlankPublishDate(t *testing.T) {
	client := &fakeDividendClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"PUBLISH_DATE":       "",
			"PLAN_NOTICE_DATE":   "2025-03-31",
		},
	})}

	rows, err := GetDividendDetail(context.Background(), client, "https://em.test/datacenter", "600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].DisclosureDate != nil {
		t.Fatalf("row = %+v", rows[0])
	}
}

func TestGetDividendDetailFallsBackWhenPublishDateMissing(t *testing.T) {
	client := &fakeDividendClient{payload: datacenterPayload([]map[string]any{
		{
			"SECURITY_CODE":      "600519",
			"SECURITY_NAME_ABBR": "贵州茅台",
			"PLAN_NOTICE_DATE":   "2025-03-31",
		},
	})}

	rows, err := GetDividendDetail(context.Background(), client, "https://em.test/datacenter", "600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].DisclosureDate == nil || *rows[0].DisclosureDate != "2025-03-31" {
		t.Fatalf("row = %+v", rows[0])
	}
}
