package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
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
