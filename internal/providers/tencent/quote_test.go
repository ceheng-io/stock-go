package tencent

import (
	"context"
	"testing"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/types"
)

type fakeQuoteClient struct {
	items []core.TencentQuoteItem
	text  string
}

func (f fakeQuoteClient) GetTencentQuote(context.Context, string) ([]core.TencentQuoteItem, error) {
	return f.items, nil
}

func (f fakeQuoteClient) GetText(context.Context, string) (string, error) {
	return f.text, nil
}

func (f fakeQuoteClient) TencentSearchURL(keyword string) string {
	return "https://smartbox.test/s3/?q=" + keyword
}

func (f fakeQuoteClient) CalendarURL() string {
	return "https://calendar.test"
}

func TestNormalizeSearchType(t *testing.T) {
	tests := map[string]types.SearchResultType{
		"GP-A":         types.SearchStock,
		"gp-a":         types.SearchStock,
		"STOCK":        types.SearchStock,
		"ZS":           types.SearchIndex,
		"INDEX":        types.SearchIndex,
		"ETF":          types.SearchFund,
		"LOF":          types.SearchFund,
		"QDII-ETF":     types.SearchFund,
		"KJ-HB":        types.SearchFund,
		"JJ":           types.SearchFund,
		"FUND":         types.SearchFund,
		"ZQ":           types.SearchBond,
		"BOND":         types.SearchBond,
		"QH":           types.SearchFutures,
		"FUTURE":       types.SearchFutures,
		"QZ":           types.SearchOption,
		"OPTION":       types.SearchOption,
		"UNKNOWN_TYPE": types.SearchOther,
	}
	for raw, want := range tests {
		if got := NormalizeSearchType(raw); got != want {
			t.Fatalf("NormalizeSearchType(%q) = %s, want %s", raw, got, want)
		}
	}
}

func TestSearchParsesSmartboxResponse(t *testing.T) {
	client := fakeQuoteClient{text: `v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A^hk~00700~\u817e\u8baf\u63a7\u80a1~TXKG~GP";`}

	results, err := Search(context.Background(), client, "茅台")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("len(results) = %d, want 2", len(results))
	}
	if results[0].Code != "sh600519" || results[0].Name != "贵州茅台" || results[0].Category != types.SearchStock {
		t.Fatalf("result[0] = %+v", results[0])
	}
	if results[1].Code != "hk00700" || results[1].Name != "腾讯控股" {
		t.Fatalf("result[1] = %+v", results[1])
	}
}

func TestGetTradingCalendar(t *testing.T) {
	client := fakeQuoteClient{text: "1990-12-19, 1990-12-20,,1990-12-21\n"}

	calendar, err := GetTradingCalendar(context.Background(), client)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"1990-12-19", "1990-12-20", "1990-12-21"}
	if len(calendar) != len(want) {
		t.Fatalf("calendar = %#v", calendar)
	}
	for i := range want {
		if calendar[i] != want[i] {
			t.Fatalf("calendar[%d] = %q, want %q", i, calendar[i], want[i])
		}
	}
}

func TestGetSimpleQuotesFiltersAndParses(t *testing.T) {
	client := fakeQuoteClient{items: []core.TencentQuoteItem{
		{
			Key:    "s_sh600519",
			Fields: []string{"1", "贵州茅台", "600519", "1700.00", "-1.23", "-0.07", "12345", "67890", "", "25000", "GP-A"},
		},
		{
			Key:    "pv_none_match",
			Fields: []string{"1"},
		},
	}}

	quotes, err := GetSimpleQuotes(context.Background(), client, []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
	quote := quotes[0]
	if quote.Name != "贵州茅台" || quote.Code != "600519" || quote.Price != 1700 || quote.ChangePercent != -0.07 {
		t.Fatalf("quote = %+v", quote)
	}
}

func TestGetSimpleQuotesEmptyCodes(t *testing.T) {
	quotes, err := GetSimpleQuotes(context.Background(), fakeQuoteClient{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 0 {
		t.Fatalf("len(quotes) = %d, want 0", len(quotes))
	}
}

func TestGetFullQuotesFiltersAndParses(t *testing.T) {
	fields := make([]string, 80)
	fields[0] = "1"
	fields[1] = "贵州茅台"
	fields[2] = "600519"
	fields[3] = "1700.00"
	fields[4] = "1710.00"
	fields[5] = "1698.00"
	fields[6] = "12345"
	fields[7] = "7000"
	fields[8] = "5345"
	fields[9] = "1699.00"
	fields[10] = "100"
	fields[19] = "1701.00"
	fields[20] = "200"
	fields[30] = "20240512143015"
	fields[31] = "-10.00"
	fields[32] = "-0.58"
	fields[33] = "1720.00"
	fields[34] = "1688.00"
	fields[36] = "1234500"
	fields[37] = "209876"
	fields[38] = "0.31"
	fields[39] = "30.5"
	fields[43] = "1.87"
	fields[44] = "21000"
	fields[45] = "22000"
	fields[46] = "12.3"
	fields[47] = "1881.00"
	fields[48] = "1539.00"
	fields[49] = "0.8"

	client := fakeQuoteClient{items: []core.TencentQuoteItem{
		{Key: "sh600519", Fields: fields},
		{Key: "pv_none_match", Fields: []string{"1"}},
	}}

	quotes, err := GetFullQuotes(context.Background(), client, []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
	quote := quotes[0]
	if quote.Name != "贵州茅台" || quote.Code != "600519" || quote.Price != 1700 {
		t.Fatalf("quote identity = %+v", quote)
	}
	if quote.Bid[0].Price != 1699 || quote.Bid[0].Volume != 100 {
		t.Fatalf("bid[0] = %+v", quote.Bid[0])
	}
	if quote.Ask[0].Price != 1701 || quote.Ask[0].Volume != 200 {
		t.Fatalf("ask[0] = %+v", quote.Ask[0])
	}
	if quote.Change != -10 || quote.ChangePercent != -0.58 || quote.Amount != 209876 {
		t.Fatalf("quote market fields = %+v", quote)
	}
	if quote.TurnoverRate == nil || *quote.TurnoverRate != 0.31 {
		t.Fatalf("TurnoverRate = %v, want 0.31", quote.TurnoverRate)
	}
	if quote.Timestamp == nil || *quote.Timestamp != 1715495415000 || quote.TZ != "Asia/Shanghai" {
		t.Fatalf("time meta = (%v, %q), want (1715495415000, Asia/Shanghai)", quote.Timestamp, quote.TZ)
	}
}

func TestGetHKUSAndFundQuotes(t *testing.T) {
	hkFields := make([]string, 50)
	hkFields[0] = "100"
	hkFields[1] = "腾讯控股"
	hkFields[2] = "00700"
	hkFields[3] = "390.20"
	hkFields[4] = "388.00"
	hkFields[5] = "389.00"
	hkFields[6] = "123456"
	hkFields[30] = "20240512143015"
	hkFields[31] = "2.20"
	hkFields[32] = "0.57"
	hkFields[33] = "395.00"
	hkFields[34] = "386.00"
	hkFields[37] = "480000"
	hkFields[40] = "100"
	hkFields[47] = "HKD"

	usFields := make([]string, 50)
	usFields[0] = "200"
	usFields[1] = "APPLE"
	usFields[2] = "AAPL"
	usFields[3] = "205.50"
	usFields[4] = "203.00"
	usFields[5] = "204.00"
	usFields[6] = "987654"
	usFields[30] = "20240512143015"
	usFields[31] = "2.50"
	usFields[32] = "1.23"
	usFields[33] = "206.00"
	usFields[34] = "202.00"
	usFields[37] = "1234567"
	usFields[38] = "0.12"
	usFields[39] = "32.1"
	usFields[43] = "1.96"
	usFields[45] = "30000"
	usFields[47] = "45.6"
	usFields[48] = "220.00"
	usFields[49] = "160.00"

	fundFields := make([]string, 9)
	fundFields[0] = "110011"
	fundFields[1] = "易方达中小盘"
	fundFields[5] = "3.5000"
	fundFields[6] = "5.2000"
	fundFields[7] = "0.0100"
	fundFields[8] = "2024-05-10"

	client := fakeQuoteClient{items: []core.TencentQuoteItem{
		{Key: "hk00700", Fields: hkFields},
		{Key: "usAAPL", Fields: usFields},
		{Key: "jj110011", Fields: fundFields},
	}}

	hk, err := GetHKQuotes(context.Background(), client, []string{"00700"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Name != "腾讯控股" || hk[0].Currency != "HKD" || hk[0].LotSize == nil || *hk[0].LotSize != 100 {
		t.Fatalf("hk = %+v", hk)
	}
	if hk[0].Timestamp == nil || *hk[0].Timestamp != 1715495415000 || hk[0].TZ != "Asia/Hong_Kong" {
		t.Fatalf("hk time meta = (%v, %q), want (1715495415000, Asia/Hong_Kong)", hk[0].Timestamp, hk[0].TZ)
	}

	us, err := GetUSQuotes(context.Background(), client, []string{"AAPL"})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Name != "APPLE" || us[0].PE == nil || *us[0].PE != 32.1 {
		t.Fatalf("us = %+v", us)
	}
	if us[0].Timestamp == nil || *us[0].Timestamp != 1715538615000 || us[0].TZ != "America/New_York" {
		t.Fatalf("us time meta = (%v, %q), want (1715538615000, America/New_York)", us[0].Timestamp, us[0].TZ)
	}

	fund, err := GetFundQuotes(context.Background(), client, []string{"110011"})
	if err != nil {
		t.Fatal(err)
	}
	if len(fund) != 1 || fund[0].Name != "易方达中小盘" || fund[0].NAV != 3.5 || fund[0].NavDate != "2024-05-10" {
		t.Fatalf("fund = %+v", fund)
	}
	if fund[0].Timestamp == nil || *fund[0].Timestamp != 1715270400000 || fund[0].TZ != "Asia/Shanghai" {
		t.Fatalf("fund time meta = (%v, %q), want (1715270400000, Asia/Shanghai)", fund[0].Timestamp, fund[0].TZ)
	}
}

func TestGetFundFlowFiltersAndParses(t *testing.T) {
	fields := make([]string, 14)
	fields[0] = "600519"
	fields[1] = "100"
	fields[2] = "80"
	fields[3] = "20"
	fields[4] = "1.2"
	fields[5] = "40"
	fields[6] = "50"
	fields[7] = "-10"
	fields[8] = "-0.6"
	fields[9] = "10"
	fields[12] = "贵州茅台"
	fields[13] = "20240613"
	client := fakeQuoteClient{items: []core.TencentQuoteItem{
		{Key: "ff_sh600519", Fields: fields},
		{Key: "pv_none_match", Fields: []string{"1"}},
		{Key: "ff_sz000001", Fields: []string{"too-short"}},
	}}

	rows, err := GetFundFlow(context.Background(), client, []string{"sh600519", "sz000001"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.Date != "20240613" {
		t.Fatalf("row identity = %+v", row)
	}
	if row.MainInflow != 100 || row.MainNet != 20 || row.RetailNet != -10 || row.TotalFlow != 10 {
		t.Fatalf("row numbers = %+v", row)
	}
	if row.Timestamp == nil || *row.Timestamp != 1718208000000 || row.TZ != "Asia/Shanghai" {
		t.Fatalf("fund flow time meta = (%v, %q), want (1718208000000, Asia/Shanghai)", row.Timestamp, row.TZ)
	}
}

func TestGetPanelLargeOrderFiltersAndParses(t *testing.T) {
	client := fakeQuoteClient{items: []core.TencentQuoteItem{
		{Key: "s_pksh600519", Fields: []string{"10.5", "20.5", "30.5", "40.5"}},
		{Key: "pv_none_match", Fields: []string{"1"}},
		{Key: "s_pksz000001", Fields: []string{"too-short"}},
	}}

	rows, err := GetPanelLargeOrder(context.Background(), client, []string{"sh600519", "sz000001"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.BuyLargeRatio != 10.5 || row.BuySmallRatio != 20.5 || row.SellLargeRatio != 30.5 || row.SellSmallRatio != 40.5 {
		t.Fatalf("row = %+v", row)
	}
}
