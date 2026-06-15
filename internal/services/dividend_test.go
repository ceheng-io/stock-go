package services

import (
	"context"
	"encoding/json"
	"testing"
)

type dividendClientStub struct {
	payload map[string]any
}

func (d *dividendClientStub) GetJSON(_ context.Context, _ string, target any) error {
	body, _ := json.Marshal(d.payload)
	return json.Unmarshal(body, target)
}

func TestDividendServiceDetail(t *testing.T) {
	client := &dividendClientStub{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "REPORT_DATE": "2024-12-31", "PRETAX_BONUS_RMB": 30.0}},
		},
	}}
	service := NewDividendService(client, "https://em.test/datacenter")

	rows, err := service.Detail(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" || rows[0].DividendPretax == nil || *rows[0].DividendPretax != 30 {
		t.Fatalf("rows = %+v", rows)
	}
}
