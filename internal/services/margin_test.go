package services

import (
	"context"
	"encoding/json"
	"testing"
)

type marginClientStub struct {
	payload map[string]any
}

func (m *marginClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(m.payload)
	return json.Unmarshal(body, target)
}

func TestMarginServiceAccountInfoAndTargetList(t *testing.T) {
	client := &marginClientStub{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"STATISTICS_DATE": "2024-12-16", "FIN_BALANCE": 1000.0}},
		},
	}}
	service := NewMarginService(client, "https://em.test/datacenter")

	account, err := service.AccountInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(account) != 1 || account[0].FinBalance == nil || *account[0].FinBalance != 1000 {
		t.Fatalf("account = %+v", account)
	}

	client.payload = map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "FIN_BUY_AMT": 200.0}},
		},
	}
	targets, err := service.TargetList(context.Background(), "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Code != "600519" || targets[0].FinBuyAmount == nil || *targets[0].FinBuyAmount != 200 {
		t.Fatalf("targets = %+v", targets)
	}
}
