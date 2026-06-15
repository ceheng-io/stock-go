package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

type fakeTopicDataClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeTopicDataClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetZTPoolBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{
				{
					"c":      "600519",
					"n":      "贵州茅台",
					"p":      1500000.0,
					"zdp":    10.0,
					"tp":     1510000.0,
					"amount": 1000000.0,
					"ltsz":   2000000.0,
					"tshare": 3000000.0,
					"hs":     2.5,
					"lbc":    3.0,
					"fbt":    93005.0,
					"lbt":    145959.0,
					"fund":   400000.0,
					"zbc":    1.0,
					"hybk":   "白酒",
					"zttj":   map[string]any{"days": 2.0, "ct": 3.0},
					"zf":     5.0,
					"zs":     1.2,
				},
			},
		},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Name != "贵州茅台" || row.Price == nil || *row.Price != 1500 || row.FirstBoardTime == nil || *row.FirstBoardTime != "09:30:05" || row.LastBoardTime == nil || *row.LastBoardTime != "14:59:59" || row.ZTStatistics != "2/3" {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	if parsed.Path != "/getTopicZTPool" {
		t.Fatalf("path = %q", parsed.Path)
	}
	query := parsed.Query()
	assertQuery(t, query, "dpt", "wz.ztzt")
	assertQuery(t, query, "Pageindex", "0")
	assertQuery(t, query, "date", "20241216")
	assertQuery(t, query, "sort", "fbt:asc")
}

func TestGetZTPoolDoesNotFallbackFromBlankPrimaryFields(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{
				{"c": "", "m": "600519", "n": "贵州茅台", "amount": "", "zb": 1000000.0},
			},
		},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "" || row.Amount != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetZTPoolFallsBackWhenPrimaryFieldsMissing(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{
				{"m": "600519", "n": "贵州茅台", "zb": 1000000.0},
			},
		},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Code != "600519" || row.Amount == nil || *row.Amount != 1000000 {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetZTPoolKeepsNilBoardTimesWhenFieldsAreMissingOrNull(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{
				{"c": "600519", "n": "贵州茅台", "lbt": nil},
			},
		},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.FirstBoardTime != nil || row.LastBoardTime != nil {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetZTPoolFormatsEmptyStringBoardTimeLikeTS(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"pool": []map[string]any{
				{"c": "600519", "n": "贵州茅台", "fbt": "", "lbt": ""},
			},
		},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.FirstBoardTime == nil || *row.FirstBoardTime != "00:00:00" || row.LastBoardTime == nil || *row.LastBoardTime != "00:00:00" {
		t.Fatalf("row = %+v, want board times formatted as 00:00:00", row)
	}
}

func TestGetZTPoolReturnsEmptyRowsForNonArrayPool(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{"pool": map[string]any{"c": "600519"}},
	}}

	rows, err := GetZTPool(context.Background(), client, "https://topic.test", ZTPoolZT, "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}

func TestGetStockChangesBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"allstock": []map[string]any{
				{"tm": 93055.0, "c": "600519", "n": "贵州茅台", "t": "8193", "i": "大单买入"},
			},
		},
	}}

	rows, err := GetStockChanges(context.Background(), client, "https://topic.test", StockChangeLargeBuy)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Time != "09:30:55" || row.Code != "600519" || row.ChangeType != StockChangeLargeBuy || row.ChangeTypeLabel != "大笔买入" || row.Info != "大单买入" {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	if parsed.Path != "/getAllStockChanges" {
		t.Fatalf("path = %q", parsed.Path)
	}
	query := parsed.Query()
	assertQuery(t, query, "type", "8193")
	assertQuery(t, query, "dpt", "wzchanges")
}

func TestGetStockChangesReturnsEmptyRowsForNonArrayAllStock(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{"allstock": map[string]any{"c": "600519"}},
	}}

	rows, err := GetStockChanges(context.Background(), client, "https://topic.test", StockChangeLargeBuy)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}

func TestGetBoardChangesBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{
			"allbk": []map[string]any{
				{
					"bkn":     "白酒",
					"bkz":     2.5,
					"bkj":     1000000.0,
					"bkc":     6.0,
					"ms":      map[string]any{"m": 0.0, "c": "600519", "n": "贵州茅台"},
					"bkdfdis": map[string]any{"8193": 2.0, "8194": 1.0},
				},
			},
		},
	}}

	rows, err := GetBoardChanges(context.Background(), client, "https://topic.test")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 {
		t.Fatalf("len(rows) = %d, want 1", len(rows))
	}
	row := rows[0]
	if row.Name != "白酒" || row.TopStockDirection != "大笔买入" || row.ChangeTypeDistribution["8193"] != 2 {
		t.Fatalf("row = %+v", row)
	}
	parsed, _ := url.Parse(client.lastURL)
	if parsed.Path != "/getAllBKChanges" {
		t.Fatalf("path = %q", parsed.Path)
	}
	assertQuery(t, parsed.Query(), "pagesize", "5000")
}

func TestGetBoardChangesReturnsEmptyRowsForNonArrayAllBK(t *testing.T) {
	client := &fakeTopicDataClient{payload: map[string]any{
		"data": map[string]any{"allbk": map[string]any{"bkn": "白酒"}},
	}}

	rows, err := GetBoardChanges(context.Background(), client, "https://topic.test")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}
