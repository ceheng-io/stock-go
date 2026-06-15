package eastmoney

import (
	"context"
	"net/url"
	"strings"
	"testing"
)

func TestGetFuturesInventorySymbolsBuildsRequest(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"result": map[string]any{
			"pages": float64(3),
			"data": []map[string]any{{
				"TRADE_CODE":        "CU",
				"TRADE_TYPE":        "铜",
				"TRADE_MARKET_CODE": "SHFE",
			}},
		},
	}}

	rows, err := GetFuturesInventorySymbols(context.Background(), client, "https://em.test/data")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "CU" || rows[0].Name != "铜" || rows[0].MarketCode != "SHFE" {
		t.Fatalf("rows = %+v", rows)
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_FUTU_POSITIONCODE")
	assertQuery(t, query, "columns", "TRADE_MARKET_CODE,TRADE_CODE,TRADE_TYPE")
	assertQuery(t, query, "filter", `(IS_MAINCODE="1")`)
	assertQuery(t, query, "pageSize", "500")
	assertQuery(t, query, "pageNumber", "1")
	if strings.Contains(client.lastURL, "pageNumber=2") {
		t.Fatalf("symbols request fetched more than first page: %s", client.lastURL)
	}
}

func TestGetFuturesInventoryBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"SECURITY_CODE":  "CU",
				"TRADE_DATE":     "2024-12-16 00:00:00",
				"ON_WARRANT_NUM": 12345.0,
				"ADDCHANGE":      -100.0,
			}},
		},
	}}

	rows, err := GetFuturesInventory(context.Background(), client, "https://em.test/data", "cu", FuturesInventoryOptions{
		StartDate: "2024-01-01",
		PageSize:  100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "CU" || rows[0].Date != "2024-12-16" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Inventory == nil || *rows[0].Inventory != 12345 || rows[0].Change == nil || *rows[0].Change != -100 {
		t.Fatalf("row numbers = %+v", rows[0])
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_FUTU_STOCKDATA")
	assertQuery(t, query, "columns", "SECURITY_CODE,TRADE_DATE,ON_WARRANT_NUM,ADDCHANGE")
	assertQuery(t, query, "filter", `(SECURITY_CODE="CU")(TRADE_DATE>='2024-01-01')`)
	assertQuery(t, query, "sortColumns", "TRADE_DATE")
	assertQuery(t, query, "sortTypes", "-1")
	assertQuery(t, query, "pageSize", "100")
}

func TestGetFuturesInventoryDoesNotFallbackFromBlankSecurityCode(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"SECURITY_CODE": "", "TRADE_DATE": "2024-12-16"}},
		},
	}}

	rows, err := GetFuturesInventory(context.Background(), client, "https://em.test/data", "cu", FuturesInventoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].Code != "" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestGetFuturesInventoryFallsBackWhenSecurityCodeMissing(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{"TRADE_DATE": "2024-12-16"}},
		},
	}}

	rows, err := GetFuturesInventory(context.Background(), client, "https://em.test/data", "cu", FuturesInventoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	if rows[0].Code != "CU" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestGetComexInventoryBuildsRequestParsesRowsAndValidatesSymbol(t *testing.T) {
	client := &fakeFuturesClient{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"REPORT_DATE":   "2024-12-16T00:00:00",
				"STORAGE_TON":   100.5,
				"STORAGE_OUNCE": 3231.0,
			}},
		},
	}}

	rows, err := GetComexInventory(context.Background(), client, "https://em.test/data", "gold", ComexInventoryOptions{
		PageSize: 100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Date != "2024-12-16" || rows[0].Name != "黄金" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].StorageTon == nil || *rows[0].StorageTon != 100.5 || rows[0].StorageOunce == nil || *rows[0].StorageOunce != 3231 {
		t.Fatalf("row numbers = %+v", rows[0])
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "reportName", "RPT_FUTUOPT_GOLDSIL")
	assertQuery(t, query, "filter", `(INDICATOR_ID1="EMI00069026")(@STORAGE_TON<>"NULL")`)
	assertQuery(t, query, "sortColumns", "REPORT_DATE")
	assertQuery(t, query, "sortTypes", "-1")
	assertQuery(t, query, "pageSize", "100")

	if _, err := GetComexInventory(context.Background(), client, "https://em.test/data", "platinum", ComexInventoryOptions{}); err == nil {
		t.Fatal("expected invalid COMEX symbol error")
	} else if !strings.Contains(err.Error(), "gold") || !strings.Contains(err.Error(), "silver") {
		t.Fatalf("error = %v, want gold/silver guidance", err)
	}
}
