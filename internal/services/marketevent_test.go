package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
	"github.com/ceheng.io/stock-go/types"
)

type marketEventClientStub struct {
	payload map[string]any
}

func (m *marketEventClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(m.payload)
	return json.Unmarshal(body, target)
}

func TestMarketEventServiceZTPoolStockChangesAndBoardChanges(t *testing.T) {
	client := &marketEventClientStub{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{{"c": "600519", "n": "贵州茅台", "p": 1500000.0}},
		},
	}}
	service := NewMarketEventService(client, "https://topic.test")

	pool, err := service.ZTPool(context.Background(), eastmoney.ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(pool) != 1 || pool[0].Code != "600519" || pool[0].Price == nil || *pool[0].Price != 1500 {
		t.Fatalf("pool = %+v", pool)
	}

	client.payload = map[string]any{
		"data": map[string]any{
			"allstock": []map[string]any{{"tm": 93055.0, "c": "600519", "n": "贵州茅台", "t": "8193", "i": "大单买入"}},
		},
	}
	changes, err := service.StockChanges(context.Background(), eastmoney.StockChangeLargeBuy)
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 || changes[0].ChangeTypeLabel != "大笔买入" {
		t.Fatalf("changes = %+v", changes)
	}

	client.payload = map[string]any{
		"data": map[string]any{
			"allbk": []map[string]any{{"bkn": "白酒", "ms": map[string]any{"m": 1.0, "c": "600519", "n": "贵州茅台"}}},
		},
	}
	boards, err := service.BoardChanges(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(boards) != 1 || boards[0].TopStockDirection != "大笔卖出" {
		t.Fatalf("boards = %+v", boards)
	}
}

func TestMarketEventServiceTHSLimitUpPool(t *testing.T) {
	client := &marketEventClientStub{payload: map[string]any{
		"status_code": 0.0,
		"data": map[string]any{
			"page": map[string]any{"limit": 1.0, "total": 1.0, "count": 1.0, "page": 1.0},
			"info": []map[string]any{
				{"code": "002190", "name": "成飞集成", "last_limit_up_time": "1749797760"},
			},
			"date": "20250613",
		},
	}}
	service := NewMarketEventService(client, "https://topic.test", "https://ths.test/limit_up_pool")

	result, err := service.THSLimitUpPool(context.Background(), types.THSLimitUpPoolOptions{Date: "2025-06-13"})
	if err != nil {
		t.Fatal(err)
	}
	if result.Date != "20250613" || len(result.Items) != 1 || result.Items[0].Code != "002190" {
		t.Fatalf("result = %+v", result)
	}
}
