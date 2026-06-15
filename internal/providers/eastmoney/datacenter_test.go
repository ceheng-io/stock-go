package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

type fakeDatacenterClient struct {
	urls     []string
	payloads []map[string]any
}

func (f *fakeDatacenterClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.urls = append(f.urls, requestURL)
	index := len(f.urls) - 1
	payload := map[string]any{"result": map[string]any{"pages": 1, "count": 0, "data": []map[string]any{}}}
	if index < len(f.payloads) {
		payload = f.payloads[index]
	}
	body, _ := json.Marshal(payload)
	return json.Unmarshal(body, target)
}

func TestFetchDatacenterFetchesAllPagesAndBuildsParams(t *testing.T) {
	client := &fakeDatacenterClient{payloads: []map[string]any{
		{
			"result": map[string]any{
				"pages": 2,
				"count": 3,
				"data": []map[string]any{
					{"TRADE_DATE": "2024-01-15 00:00:00", "VALUE": 1},
					{"TRADE_DATE": "2024-01-16T00:00:00.000", "VALUE": 2},
				},
			},
		},
		{
			"result": map[string]any{
				"pages": 2,
				"count": 3,
				"data": []map[string]any{
					{"TRADE_DATE": "2024-01-17", "VALUE": 3},
				},
			},
		},
	}}

	result, err := FetchDatacenter(context.Background(), client, "https://em.test/datacenter", DatacenterQuery{
		ReportName:   "RPT_TEST",
		Filter:       "(TRADE_DATE>='2024-01-15')",
		SortColumns:  "TRADE_DATE",
		SortTypes:    "-1",
		PageSize:     2,
		StartPage:    1,
		QuoteColumns: "f2,f3",
		QuoteType:    "0",
		ExtraParams:  map[string]string{"extra": "1"},
	}, func(item map[string]any, index int) string {
		return ParseDCDate(item["TRADE_DATE"]) + ":" + string(rune('0'+index))
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Total != 3 || result.Pages != 2 {
		t.Fatalf("result meta = %+v, want total=3 pages=2", result)
	}
	wantData := []string{"2024-01-15:0", "2024-01-16:1", "2024-01-17:2"}
	if len(result.Data) != len(wantData) {
		t.Fatalf("len(result.Data) = %d, want %d; data=%#v", len(result.Data), len(wantData), result.Data)
	}
	for i := range wantData {
		if result.Data[i] != wantData[i] {
			t.Fatalf("result.Data[%d] = %q, want %q", i, result.Data[i], wantData[i])
		}
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2; urls=%#v", len(client.urls), client.urls)
	}
	assertDatacenterQuery(t, client.urls[0], map[string]string{
		"reportName":   "RPT_TEST",
		"columns":      "ALL",
		"pageSize":     "2",
		"pageNumber":   "1",
		"source":       "WEB",
		"client":       "WEB",
		"filter":       "(TRADE_DATE>='2024-01-15')",
		"sortColumns":  "TRADE_DATE",
		"sortTypes":    "-1",
		"quoteColumns": "f2,f3",
		"quoteType":    "0",
		"extra":        "1",
	})
	assertDatacenterQuery(t, client.urls[1], map[string]string{"pageNumber": "2"})
}

func TestFetchDatacenterCanFetchOnlyFirstPage(t *testing.T) {
	fetchAllPages := false
	client := &fakeDatacenterClient{payloads: []map[string]any{
		{"result": map[string]any{"pages": 3, "count": 3, "data": []map[string]any{{"VALUE": "first"}}}},
	}}

	result, err := FetchDatacenter(context.Background(), client, "https://em.test/datacenter", DatacenterQuery{
		ReportName:    "RPT_TEST",
		FetchAllPages: &fetchAllPages,
	}, func(item map[string]any, index int) string {
		return item["VALUE"].(string)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 1 {
		t.Fatalf("requests = %d, want 1", len(client.urls))
	}
	if len(result.Data) != 1 || result.Data[0] != "first" {
		t.Fatalf("result = %+v", result)
	}
}

func TestFetchDatacenterStopsOnNonArrayDataLikeTypeScript(t *testing.T) {
	client := &fakeDatacenterClient{payloads: []map[string]any{
		{"result": map[string]any{
			"pages": 3,
			"count": 3,
			"data":  []map[string]any{{"VALUE": "first"}},
		}},
		{"result": map[string]any{
			"pages": 3,
			"count": 3,
			"data":  map[string]any{"VALUE": "not-array"},
		}},
		{"result": map[string]any{
			"pages": 3,
			"count": 3,
			"data":  []map[string]any{{"VALUE": "third"}},
		}},
	}}

	result, err := FetchDatacenter(context.Background(), client, "https://em.test/datacenter", DatacenterQuery{
		ReportName: "RPT_TEST",
		PageSize:   1,
	}, func(item map[string]any, index int) string {
		return item["VALUE"].(string)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 2 {
		t.Fatalf("requests = %d, want 2", len(client.urls))
	}
	if result.Total != 3 || result.Pages != 3 {
		t.Fatalf("result meta = %+v, want total=3 pages=3", result)
	}
	if len(result.Data) != 1 || result.Data[0] != "first" {
		t.Fatalf("result data = %#v, want first page only", result.Data)
	}
}

func TestFetchDatacenterCompatCanFetchOnlyFirstPage(t *testing.T) {
	client := &fakeDatacenterClient{payloads: []map[string]any{
		{"result": map[string]any{"pages": 3, "count": 3, "data": []map[string]any{{"VALUE": "first"}}}},
		{"result": map[string]any{"pages": 3, "count": 3, "data": []map[string]any{{"VALUE": "second"}}}},
	}}

	rows, err := fetchDatacenter(context.Background(), client, "https://em.test/datacenter", datacenterOptions{
		reportName:    "RPT_TEST",
		pageSize:      "500",
		fetchAllPages: datacenterBool(false),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(client.urls) != 1 {
		t.Fatalf("requests = %d, want 1", len(client.urls))
	}
	if len(rows) != 1 || rows[0]["VALUE"] != "first" {
		t.Fatalf("rows = %#v", rows)
	}
}

func TestFetchDatacenterListReturnsOnlyData(t *testing.T) {
	client := &fakeDatacenterClient{payloads: []map[string]any{
		{"result": map[string]any{"pages": 1, "count": 1, "data": []map[string]any{{"VALUE": "ok"}}}},
	}}

	rows, err := FetchDatacenterList(context.Background(), client, "https://em.test/datacenter", DatacenterQuery{
		ReportName: "RPT_TEST",
	}, func(item map[string]any, index int) string {
		return item["VALUE"].(string)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0] != "ok" {
		t.Fatalf("rows = %#v", rows)
	}
}

func TestParseDCDate(t *testing.T) {
	tests := map[any]string{
		nil:                       "",
		"2024-01-15":              "2024-01-15",
		"2024-01-15 00:00:00":     "2024-01-15",
		"2024-01-15T00:00:00.000": "2024-01-15",
		"20240115":                "2024-01-15",
		"bad":                     "bad",
	}
	for input, want := range tests {
		if got := ParseDCDate(input); got != want {
			t.Fatalf("ParseDCDate(%v) = %q, want %q", input, got, want)
		}
	}
}

func assertDatacenterQuery(t *testing.T, requestURL string, want map[string]string) {
	t.Helper()
	parsed, err := url.Parse(requestURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	for key, value := range want {
		if got := query.Get(key); got != value {
			t.Fatalf("query[%s] = %q, want %q; url=%s", key, got, value, requestURL)
		}
	}
}
