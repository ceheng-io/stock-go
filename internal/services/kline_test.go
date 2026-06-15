package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng.io/stock-go/internal/providers/eastmoney"
)

type klineClientStub struct {
	lastURL string
}

func (k *klineClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	k.lastURL = requestURL
	payload := map[string]any{
		"data": map[string]any{
			"klines": []string{"2024-12-16,1500,1510,1520,1490,12345,67890000,2,1.5,22.3,0.5"},
		},
	}
	b, _ := json.Marshal(payload)
	return json.Unmarshal(b, target)
}

func TestKlineServiceCN(t *testing.T) {
	client := &klineClientStub{}
	service := NewKlineService(client, KlineURLs{CN: "https://em.test/kline"})

	rows, err := service.CN(context.Background(), "600519", eastmoney.HistoryKlineOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" {
		t.Fatalf("rows = %+v", rows)
	}
	if client.lastURL == "" {
		t.Fatal("expected service to call client")
	}
}

func TestKlineServiceCNMinute(t *testing.T) {
	client := &klineClientStub{}
	service := NewKlineService(client, KlineURLs{CN: "https://em.test/kline"})

	rows, err := service.CNMinute(context.Background(), "600519", eastmoney.MinuteKlineOptions{Period: eastmoney.MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 1 || rows.Klines[0].Code != "600519" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestKlineServiceHKUS(t *testing.T) {
	client := &klineClientStub{}
	service := NewKlineService(client, KlineURLs{
		HK: "https://em.test/hk",
		US: "https://em.test/us",
	})

	hk, err := service.HK(context.Background(), "00700", eastmoney.HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Code != "00700" || hk[0].Currency != "HKD" {
		t.Fatalf("hk = %+v", hk)
	}

	us, err := service.US(context.Background(), "105.AAPL", eastmoney.HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Code != "AAPL" || us[0].Currency != "USD" {
		t.Fatalf("us = %+v", us)
	}
}

func TestKlineServiceHKUSMinute(t *testing.T) {
	client := &klineClientStub{}
	service := NewKlineService(client, KlineURLs{
		HK:       "https://em.test/hk-kline",
		HKTrends: "https://em.test/hk-trends",
		US:       "https://em.test/us-kline",
		USTrends: "https://em.test/us-trends",
	})

	hk, err := service.HKMinute(context.Background(), "00700", eastmoney.MinuteKlineOptions{Period: eastmoney.MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk.Klines) != 1 || hk.Klines[0].Code != "00700" || hk.Klines[0].Currency != "HKD" {
		t.Fatalf("hk = %+v", hk)
	}

	us, err := service.USMinute(context.Background(), "105.AAPL", eastmoney.MinuteKlineOptions{Period: eastmoney.MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(us.Klines) != 1 || us.Klines[0].Code != "AAPL" || us.Klines[0].Currency != "USD" {
		t.Fatalf("us = %+v", us)
	}
}
