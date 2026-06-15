package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
)

type futuresClientStub struct {
	lastURL string
	payload map[string]any
}

func (f *futuresClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	if f.payload == nil {
		f.payload = map[string]any{
			"data": map[string]any{
				"code":   "rb2605",
				"name":   "螺纹钢2605",
				"klines": []string{"2024-12-16,3500,3520,3530,3490,12345,67890000,1.14,0.57,20,0,0,98765,0"},
			},
		}
	}
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestFuturesServiceKline(t *testing.T) {
	client := &futuresClientStub{}
	service := NewFuturesService(client, FuturesURLs{Kline: "https://em.test/futures"})

	rows, err := service.Kline(context.Background(), "rb2605", eastmoney.FuturesKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "rb2605" || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 98765 {
		t.Fatalf("rows = %+v", rows)
	}
	if client.lastURL == "" {
		t.Fatal("expected service to call client")
	}
}

func TestFuturesServiceGlobalSpotAndKline(t *testing.T) {
	client := &futuresClientStub{payload: map[string]any{
		"total": float64(1),
		"list": []map[string]any{
			{"dm": "GC00Y", "name": "COMEX黄金", "p": 2400.5, "zde": 10.5, "zdf": 0.44, "o": 2390, "h": 2410, "l": 2380, "zjsj": 2390, "vol": 1000, "wp": 10, "np": 20, "ccl": 3000},
		},
	}}
	service := NewFuturesService(client, FuturesURLs{
		Kline:      "https://em.test/kline",
		GlobalSpot: "https://em.test/global",
	})

	spot, err := service.GlobalSpot(context.Background(), eastmoney.GlobalFuturesSpotOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(spot) != 1 || spot[0].Code != "GC00Y" {
		t.Fatalf("spot = %+v", spot)
	}

	client.payload = map[string]any{
		"data": map[string]any{
			"code":   "HG00Y",
			"name":   "COMEX铜",
			"klines": []string{"2024-12-16,4.1,4.2,4.3,4.0,123,456,2.5,1.2,0.1,0,0,789,0"},
		},
	}
	rows, err := service.GlobalKline(context.Background(), "HG00Y", eastmoney.GlobalFuturesKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "HG00Y" || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 789 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestFuturesServiceInventory(t *testing.T) {
	client := &futuresClientStub{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"TRADE_CODE":        "CU",
				"TRADE_TYPE":        "铜",
				"TRADE_MARKET_CODE": "SHFE",
			}},
		},
	}}
	service := NewFuturesService(client, FuturesURLs{Datacenter: "https://em.test/data"})

	symbols, err := service.InventorySymbols(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(symbols) != 1 || symbols[0].Code != "CU" || symbols[0].Name != "铜" {
		t.Fatalf("symbols = %+v", symbols)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"SECURITY_CODE":  "CU",
				"TRADE_DATE":     "2024-12-16",
				"ON_WARRANT_NUM": 12345.0,
				"ADDCHANGE":      -100.0,
			}},
		},
	}
	inventory, err := service.Inventory(context.Background(), "cu", eastmoney.FuturesInventoryOptions{StartDate: "2024-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 1 || inventory[0].Code != "CU" || inventory[0].Inventory == nil || *inventory[0].Inventory != 12345 {
		t.Fatalf("inventory = %+v", inventory)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"REPORT_DATE":   "2024-12-16",
				"STORAGE_TON":   100.5,
				"STORAGE_OUNCE": 3231.0,
			}},
		},
	}
	comex, err := service.ComexInventory(context.Background(), "gold", eastmoney.ComexInventoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(comex) != 1 || comex[0].Name != "黄金" || comex[0].StorageTon == nil || *comex[0].StorageTon != 100.5 {
		t.Fatalf("comex = %+v", comex)
	}
}
