package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/internal/providers/tencent"
)

type dataClientStub struct {
	payload map[string]any
	lists   map[string][]string
	text    string
}

func (d *dataClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	if d.payload != nil {
		body, _ := json.Marshal(d.payload)
		return json.Unmarshal(body, target)
	}
	payload := struct {
		Success bool     `json:"success"`
		List    []string `json:"list"`
	}{Success: true, List: d.lists[requestURL]}
	body, _ := json.Marshal(payload)
	return json.Unmarshal(body, target)
}

func (d *dataClientStub) GetText(context.Context, string) (string, error) {
	if d.text != "" {
		return d.text, nil
	}
	return `v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A";`, nil
}

func (d *dataClientStub) TencentSearchURL(keyword string) string {
	return "https://smartbox.test/s3/?q=" + keyword
}

func (d *dataClientStub) AShareListURL() string { return "a" }
func (d *dataClientStub) USListURL() string     { return "us" }
func (d *dataClientStub) HKListURL() string     { return "hk" }
func (d *dataClientStub) FundListURL() string   { return "fund" }

func TestDataServiceSearchAndCodes(t *testing.T) {
	client := &dataClientStub{lists: map[string][]string{
		"a":    {"sh600000", "bj830799"},
		"us":   {"105.AAPL"},
		"hk":   {"00700"},
		"fund": {"110011"},
	}}
	service := NewDataService(client, DataServiceOptions{})

	results, err := service.Search(context.Background(), "茅台")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "sh600519" {
		t.Fatalf("results = %+v", results)
	}
	cn, err := service.CodesCN(context.Background(), tencent.CodeListOptions{Market: tencent.AShareMarketBJ})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0] != "bj830799" {
		t.Fatalf("cn = %#v", cn)
	}
	us, err := service.CodesUS(context.Background(), tencent.USCodeListOptions{Simple: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0] != "AAPL" {
		t.Fatalf("us = %#v", us)
	}
}

func TestDataServiceBlockTradeMarginAndDividend(t *testing.T) {
	client := &dataClientStub{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "DEAL_PRICE": 1490.0}},
		},
	}}
	service := NewDataService(client, DataServiceOptions{DatacenterURL: "https://em.test/datacenter"})

	detail, err := service.BlockTradeDetail(context.Background(), eastmoney.BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].DealPrice == nil || *detail[0].DealPrice != 1490 {
		t.Fatalf("detail = %+v", detail)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "FIN_BUY_AMT": 200.0}},
		},
	}
	targets, err := service.MarginTargetList(context.Background(), "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].FinBuyAmount == nil || *targets[0].FinBuyAmount != 200 {
		t.Fatalf("targets = %+v", targets)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "REPORT_DATE": "2024-12-31", "PRETAX_BONUS_RMB": 30.0}},
		},
	}
	dividends, err := service.DividendDetail(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(dividends) != 1 || dividends[0].DividendPretax == nil || *dividends[0].DividendPretax != 30 {
		t.Fatalf("dividends = %+v", dividends)
	}
}
