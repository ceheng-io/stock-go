package tencent

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ceheng.io/stock-go/internal/core"
)

type fakeTimelineClient struct {
	requestURL string
	text       string
}

func (f *fakeTimelineClient) GetText(_ context.Context, requestURL string) (string, error) {
	f.requestURL = requestURL
	return f.text, nil
}

func TestGetTodayTimelineBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeTimelineClient{text: `{
		"code": 0,
		"data": {
			"sz000001": {
				"data": {
					"date": "20240613",
					"data": ["0930 10.00 0 0.00", "0931 10.10 100 101000.00"]
				},
				"qt": {"sz000001": ["", "", "", "", "9.90"]}
			}
		}
	}`}

	row, err := GetTodayTimeline(context.Background(), client, "https://minute.test/query", "sz000001")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(client.requestURL, "https://minute.test/query?code=sz000001") {
		t.Fatalf("requestURL = %q", client.requestURL)
	}
	if row.Code != "sz000001" || row.Date != "20240613" || row.PreClose == nil || *row.PreClose != 9.9 {
		t.Fatalf("meta = %+v", row)
	}
	if len(row.Data) != 2 {
		t.Fatalf("len(data) = %d, want 2", len(row.Data))
	}
	if row.Data[1].Time != "09:31" || row.Data[1].Price != 10.1 {
		t.Fatalf("tick = %+v", row.Data[1])
	}
	if row.Data[1].Volume != 10000 || row.Data[1].Amount != 101000 || row.Data[1].AvgPrice != 10.1 {
		t.Fatalf("tick numbers = %+v", row.Data[1])
	}
}

func TestGetTodayTimelineReturnsEmptyForMissingStockData(t *testing.T) {
	client := &fakeTimelineClient{text: `{"code":0,"data":{}}`}

	row, err := GetTodayTimeline(context.Background(), client, "https://minute.test/query", "sz000001")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "sz000001" || row.Date != "" || row.PreClose == nil || *row.PreClose != 0 || len(row.Data) != 0 {
		t.Fatalf("row = %+v", row)
	}
}

func TestGetTodayTimelineUsesZeroPreCloseWhenQuoteFieldMissing(t *testing.T) {
	client := &fakeTimelineClient{text: `{
		"code": 0,
		"data": {
			"sz000001": {
				"data": {"date": "20240613", "data": ["0931 10.10 100 101000.00"]},
				"qt": {"sz000001": ["", "", "", ""]}
			}
		}
	}`}

	row, err := GetTodayTimeline(context.Background(), client, "https://minute.test/query", "sz000001")
	if err != nil {
		t.Fatal(err)
	}
	if row.PreClose == nil || *row.PreClose != 0 {
		t.Fatalf("PreClose = %v, want 0", row.PreClose)
	}
}

func TestGetTodayTimelineReturnsUpstreamError(t *testing.T) {
	client := &fakeTimelineClient{text: `{"code":1,"msg":"bad request"}`}

	if _, err := GetTodayTimeline(context.Background(), client, "https://minute.test/query", "sz000001"); err == nil {
		t.Fatal("expected upstream error")
	} else if !strings.Contains(err.Error(), "bad request") {
		t.Fatalf("err = %v", err)
	} else {
		var coded core.CodedError
		if !errors.As(err, &coded) {
			t.Fatalf("err = %T %v, want coded upstream error", err, err)
		}
		if code := coded.SDKCode(); code != "UPSTREAM_ERROR" {
			t.Fatalf("error code = %q, want UPSTREAM_ERROR; err=%v", code, err)
		}
	}
}
