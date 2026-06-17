package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ceheng-io/stock-go/internal/core"
	"github.com/ceheng-io/stock-go/internal/providers/eastmoney"
)

type boardClientStub struct {
	lastURL string
	payload map[string]any
	handler func(requestURL string, target any) error
}

func (b *boardClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	b.lastURL = requestURL
	if b.handler != nil {
		return b.handler(requestURL, target)
	}
	if b.payload == nil {
		b.payload = map[string]any{
			"data": map[string]any{
				"diff": []map[string]any{
					{"f12": "BK0001", "f14": "测试板块", "f3": 1.2},
				},
			},
		}
	}
	body, _ := json.Marshal(b.payload)
	return json.Unmarshal(body, target)
}

func TestBoardServiceIndustryConceptList(t *testing.T) {
	client := &boardClientStub{}
	service := NewBoardService(client, BoardURLs{
		IndustryList: "https://em.test/industry",
		ConceptList:  "https://em.test/concept",
	})

	industry, err := service.IndustryList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(industry) != 1 || industry[0].Code != "BK0001" {
		t.Fatalf("industry = %+v", industry)
	}

	concept, err := service.ConceptList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(concept) != 1 || concept[0].Name != "测试板块" {
		t.Fatalf("concept = %+v", concept)
	}
}

func TestBoardServiceResolvesIndustryNameToCode(t *testing.T) {
	client := &boardClientStub{}
	client.handler = func(requestURL string, target any) error {
		var payload map[string]any
		switch {
		case strings.HasPrefix(requestURL, "https://em.test/industry?"):
			payload = map[string]any{
				"data": map[string]any{
					"diff": []map[string]any{{"f12": "BK0001", "f14": "酿酒行业", "f3": 1.2}},
				},
			}
		case strings.HasPrefix(requestURL, "https://em.test/spot?"):
			payload = map[string]any{"data": map[string]any{"f43": 1050.0}}
		default:
			t.Fatalf("unexpected request URL %q", requestURL)
		}
		body, _ := json.Marshal(payload)
		return json.Unmarshal(body, target)
	}
	service := NewBoardService(client, BoardURLs{
		IndustryList: "https://em.test/industry",
		IndustrySpot: "https://em.test/spot",
	})

	spot, err := service.IndustrySpot(context.Background(), "酿酒行业")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot) != 10 || spot[0].Value == nil || *spot[0].Value != 10.5 {
		t.Fatalf("spot = %+v", spot)
	}
	if client.lastURL == "" || !strings.Contains(client.lastURL, "secid=90.BK0001") {
		t.Fatalf("lastURL = %q, want resolved BK0001 secid", client.lastURL)
	}
}

func TestBoardServiceUnknownIndustryNameReturnsNotFound(t *testing.T) {
	client := &boardClientStub{payload: map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{{"f12": "BK0001", "f14": "酿酒行业", "f3": 1.2}},
		},
	}}
	service := NewBoardService(client, BoardURLs{
		IndustryList: "https://em.test/industry",
		IndustrySpot: "https://em.test/spot",
	})

	_, err := service.IndustrySpot(context.Background(), "不存在行业")
	if err == nil {
		t.Fatal("expected not found error")
	}
	var coded core.CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("err = %T %v, want coded not found error", err, err)
	}
	if code := coded.SDKCode(); code != "NOT_FOUND" {
		t.Fatalf("error code = %q, want NOT_FOUND; err=%v", code, err)
	}
}

func TestBoardServiceSpotAndConstituents(t *testing.T) {
	client := &boardClientStub{payload: map[string]any{
		"data": map[string]any{
			"f43": 1050.0,
		},
	}}
	service := NewBoardService(client, BoardURLs{IndustrySpot: "https://em.test/spot"})

	spot, err := service.IndustrySpot(context.Background(), "BK0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot) != 10 || spot[0].Value == nil || *spot[0].Value != 10.5 {
		t.Fatalf("spot = %+v", spot)
	}

	client.payload = map[string]any{
		"data": map[string]any{
			"diff": []map[string]any{
				{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0},
			},
		},
	}
	service = NewBoardService(client, BoardURLs{IndustryConstituents: "https://em.test/constituents"})
	constituents, err := service.IndustryConstituents(context.Background(), "BK0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(constituents) != 1 || constituents[0].Code != "600519" {
		t.Fatalf("constituents = %+v", constituents)
	}
}

func TestBoardServiceKlineAndMinute(t *testing.T) {
	client := &boardClientStub{payload: map[string]any{
		"data": map[string]any{
			"klines": []string{"2024-12-16,100,105,106,99,12345,67890000,3.5,2,2,1.2"},
		},
	}}
	service := NewBoardService(client, BoardURLs{
		IndustryKline: "https://em.test/kline",
	})

	klines, err := service.IndustryKline(context.Background(), "BK0001", eastmoney.HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 105 {
		t.Fatalf("klines = %+v", klines)
	}

	minute, err := service.IndustryMinute(context.Background(), "BK0001", eastmoney.MinuteKlineOptions{Period: eastmoney.MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(minute.Klines) != 1 || minute.Klines[0].Close == nil || *minute.Klines[0].Close != 105 {
		t.Fatalf("minute = %+v", minute)
	}
}
