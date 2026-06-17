package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
)

type blockTradeClientStub struct {
	payload map[string]any
}

func (b *blockTradeClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(b.payload)
	return json.Unmarshal(body, target)
}

func TestBlockTradeServiceMarketStatAndDetail(t *testing.T) {
	client := &blockTradeClientStub{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"TRADE_DATE": "2024-12-16", "CLOSE_PRICE": 3500.0, "TURNOVER": 1000000.0}},
		},
	}}
	service := NewBlockTradeService(client, "https://em.test/datacenter")

	market, err := service.MarketStat(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(market) != 1 || market[0].SHClose == nil || *market[0].SHClose != 3500 {
		t.Fatalf("market = %+v", market)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "DEAL_PRICE": 1490.0}},
		},
	}
	detail, err := service.Detail(context.Background(), eastmoney.BlockTradeDateOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].Code != "600519" || detail[0].DealPrice == nil || *detail[0].DealPrice != 1490 {
		t.Fatalf("detail = %+v", detail)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "DEAL_NUM": 2.0}},
		},
	}
	daily, err := service.DailyStat(context.Background(), eastmoney.BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(daily) != 1 || daily[0].DealCount == nil || *daily[0].DealCount != 2 {
		t.Fatalf("daily = %+v", daily)
	}
}
