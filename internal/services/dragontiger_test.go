package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
)

type dragonTigerClientStub struct {
	payload map[string]any
}

func (d *dragonTigerClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(d.payload)
	return json.Unmarshal(body, target)
}

func TestDragonTigerServiceDetailAndStats(t *testing.T) {
	client := &dragonTigerClientStub{payload: map[string]any{
		"result": map[string]any{"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "BILLBOARD_NET_AMT": 1000.0}}},
	}}
	service := NewDragonTigerService(client, "https://em.test/datacenter")

	detail, err := service.Detail(context.Background(), eastmoney.DragonTigerDateOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].Code != "600519" {
		t.Fatalf("detail = %+v", detail)
	}

	client.payload = map[string]any{
		"result": map[string]any{"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "LATEST_TDATE": "2024-12-16", "BILLBOARD_TIMES": 3.0}}},
	}
	stats, err := service.StockStats(context.Background(), eastmoney.DragonTigerPeriodOneMonth)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 1 || stats[0].Count == nil || *stats[0].Count != 3 {
		t.Fatalf("stats = %+v", stats)
	}
}
