package ths

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"testing"

	"github.com/ceheng-io/stock-go/internal/core"
)

type fakeLimitUpClient struct {
	lastURL string
	payload map[string]any
}

func TestGetLimitUpPoolIntegration(t *testing.T) {
	if os.Getenv("CEHENG_INTEGRATION") != "1" {
		t.Skip("set CEHENG_INTEGRATION=1 to run live Tonghuashun request")
	}
	client := core.NewClient(core.Config{
		ProviderPolicies: map[core.ProviderName]core.ProviderPolicy{
			core.ProviderTHS: {
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
				Headers: map[string]string{
					"Accept":  "application/json, text/plain, */*",
					"Referer": "https://data.10jqka.com.cn/market/ztStock/",
					"Cookie":  "v=A0aSl97zW6psJw9OiWEn2CdlkTfNp4vvXOm-xTBvMghEJ-jpmDfacSx7DtgD",
				},
			},
		},
	})
	result, err := GetLimitUpPool(context.Background(), client, "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool", LimitUpPoolOptions{Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Date == "" {
		t.Fatalf("result = %+v", result)
	}
}

func (f *fakeLimitUpClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetLimitUpPoolBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeLimitUpClient{payload: map[string]any{
		"status_code": 0.0,
		"status_msg":  "success",
		"data": map[string]any{
			"page": map[string]any{"limit": 2.0, "total": 56.0, "count": 28.0, "page": 1.0},
			"info": []map[string]any{
				{
					"code":                "002190",
					"name":                "成飞集成",
					"latest":              36.72,
					"change_rate":         10.006,
					"first_limit_up_time": "1749797760",
					"last_limit_up_time":  "1749797760",
					"open_num":            nil,
					"limit_up_type":       "换手板",
					"order_volume":        2095399.0,
					"order_amount":        7.6943051e7,
					"turnover_rate":       25.1498,
					"currency_value":      1.31725415e10,
					"reason_type":         "成飞概念+军工",
					"high_days":           "首板",
					"high_days_value":     65537.0,
					"change_tag":          "FIRST_LIMIT",
					"market_type":         "HS",
					"market_id":           33.0,
					"is_new":              0.0,
					"is_again_limit":      0.0,
					"limit_up_suc_rate":   0.7777777777777778,
					"time_preview":        []any{2.3068, 10.006},
				},
			},
			"limit_up_count": map[string]any{
				"today":     map[string]any{"num": 56.0, "history_num": 67.0, "rate": 0.835820895522388, "open_num": 11.0},
				"yesterday": map[string]any{"num": 64.0, "history_num": 79.0, "rate": 0.810126582278481, "open_num": 15.0},
			},
			"limit_down_count": map[string]any{
				"today":     map[string]any{"num": 5.0, "history_num": 11.0, "rate": 0.45454545454545453, "open_num": 6.0},
				"yesterday": map[string]any{"num": 0.0, "history_num": 0.0, "rate": nil, "open_num": 0.0},
			},
			"date": "20250613",
			"trade_status": map[string]any{
				"id":         "closed",
				"name":       "已收盘",
				"start_time": "15:30",
				"end_time":   "23:59:59.999999999",
			},
		},
	}}

	result, err := GetLimitUpPool(context.Background(), client, "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool", LimitUpPoolOptions{
		Date:       "2025-06-13",
		Page:       1,
		Limit:      2,
		Filter:     "HS,GEM2STAR",
		OrderField: LimitUpOrderLastLimitUpTime,
		OrderType:  LimitUpOrderDesc,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Date != "20250613" || result.Page.Total != 56 || result.LimitUpCount.Today.Num != 56 || result.LimitDownCount.Yesterday.Rate != nil {
		t.Fatalf("result metadata = %+v", result)
	}
	if result.TradeStatus.ID != "closed" || result.TradeStatus.Name != "已收盘" {
		t.Fatalf("trade status = %+v", result.TradeStatus)
	}
	if len(result.Items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(result.Items))
	}
	row := result.Items[0]
	if row.Code != "002190" || row.Name != "成飞集成" || row.OpenNum != nil || row.FirstLimitUpTimeText != "14:56:00" || row.LastLimitUpTimeText != "14:56:00" {
		t.Fatalf("row = %+v", row)
	}
	if row.Latest == nil || *row.Latest != 36.72 || len(row.TimePreview) != 2 || row.TimePreview[1] != 10.006 {
		t.Fatalf("row values = %+v", row)
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("date") != "20250613" || query.Get("page") != "1" || query.Get("limit") != "2" {
		t.Fatalf("query = %v", query)
	}
	if query.Get("field") != defaultLimitUpFields || query.Get("order_field") != "330324" || query.Get("order_type") != "0" || query.Get("filter") != "HS,GEM2STAR" {
		t.Fatalf("query = %v", query)
	}
	if query.Get("_") == "" {
		t.Fatalf("query missing cache buster: %v", query)
	}
}

func TestGetLimitUpPoolDefaultsOptions(t *testing.T) {
	client := &fakeLimitUpClient{payload: map[string]any{
		"status_code": 0.0,
		"data": map[string]any{
			"page": map[string]any{"limit": 50.0, "total": 0.0, "count": 0.0, "page": 1.0},
			"info": []map[string]any{},
		},
	}}

	result, err := GetLimitUpPool(context.Background(), client, "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool", LimitUpPoolOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if result.Items == nil {
		t.Fatalf("Items = nil, want empty slice")
	}
	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	if query.Get("page") != "1" || query.Get("limit") != "50" || query.Get("date") != "" || query.Get("filter") != defaultLimitUpFilter {
		t.Fatalf("query = %v", query)
	}
}

func TestGetLimitUpPoolReturnsUpstreamErrorForNonZeroStatus(t *testing.T) {
	client := &fakeLimitUpClient{payload: map[string]any{
		"status_code": 102.0,
		"status_msg":  "forbidden",
	}}

	_, err := GetLimitUpPool(context.Background(), client, "https://data.10jqka.com.cn/dataapi/limit_up/limit_up_pool", LimitUpPoolOptions{})
	if err == nil {
		t.Fatal("expected upstream error")
	}
	var coded core.CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("err = %T %v, want coded error", err, err)
	}
	if code := coded.SDKCode(); code != "UPSTREAM_ERROR" {
		t.Fatalf("error code = %q, want UPSTREAM_ERROR", code)
	}
}
