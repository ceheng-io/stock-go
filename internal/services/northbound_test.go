package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
)

type northboundClientStub struct {
	payload map[string]any
}

func (n *northboundClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(n.payload)
	return json.Unmarshal(body, target)
}

func TestNorthboundServiceMinuteAndHistory(t *testing.T) {
	client := &northboundClientStub{payload: map[string]any{
		"data": map[string]any{
			"s2nDate": "20241216",
			"s2n":     []string{"09:31,100,0,200,0,300"},
		},
	}}
	service := NewNorthboundService(client, NorthboundURLs{
		Minute:     "https://em.test/minute",
		Datacenter: "https://em.test/datacenter",
	})

	minute, err := service.Minute(context.Background(), eastmoney.NorthboundNorth)
	if err != nil {
		t.Fatal(err)
	}
	if len(minute) != 1 || minute[0].Date != "2024-12-16" {
		t.Fatalf("minute = %+v", minute)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"TRADE_DATE": "2024-12-16", "NET_DEAL_AMT": 1000.0}},
		},
	}
	history, err := service.History(context.Background(), eastmoney.NorthboundNorth, eastmoney.NorthboundHistoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].NetBuyAmount == nil || *history[0].NetBuyAmount != 1000 {
		t.Fatalf("history = %+v", history)
	}
}
