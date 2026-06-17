package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
)

type fundFlowClientStub struct {
	lastURL string
	payload map[string]any
}

func (f *fundFlowClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestFundFlowServiceHistoryAndMarket(t *testing.T) {
	client := &fundFlowClientStub{payload: map[string]any{
		"data": map[string]any{"klines": []string{"2024-12-16,1000,100,200,300,400,10,1,2,3,4,1500,1.5,3500,1.2"}},
	}}
	service := NewFundFlowService(client, FundFlowURLs{FFlow: "https://em.test/fflow"})

	individual, err := service.Individual(context.Background(), "600519", eastmoney.FundFlowOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(individual) != 1 || individual[0].MainNetInflow == nil || *individual[0].MainNetInflow != 1000 {
		t.Fatalf("individual = %+v", individual)
	}

	client.payload = map[string]any{
		"data": map[string]any{"klines": []string{"2024-12-16,1000,100,200,300,400,10,1,2,3,4,3500,1.2,11000,1.4"}},
	}
	market, err := service.Market(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(market) != 1 || market[0].SHClose == nil || *market[0].SHClose != 3500 {
		t.Fatalf("market = %+v", market)
	}
}

func TestFundFlowServiceRank(t *testing.T) {
	client := &fundFlowClientStub{payload: map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0, "f3": 1.5, "f62": 1000.0}},
		},
	}}
	service := NewFundFlowService(client, FundFlowURLs{Clist: "https://em.test/clist"})

	rows, err := service.Rank(context.Background(), eastmoney.FundFlowRankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" {
		t.Fatalf("rows = %+v", rows)
	}
}
