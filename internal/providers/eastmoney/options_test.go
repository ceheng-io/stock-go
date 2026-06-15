package eastmoney

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

type fakeOptionsClient struct {
	lastURL string
	payload map[string]any
}

func (f *fakeOptionsClient) GetJSON(_ context.Context, requestURL string, target any) error {
	f.lastURL = requestURL
	body, _ := json.Marshal(f.payload)
	return json.Unmarshal(body, target)
}

func TestGetCFFEXOptionQuotesBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeOptionsClient{payload: map[string]any{
		"total": float64(1),
		"list": []map[string]any{{
			"dm": "io2501C4000", "name": "沪深300购2501",
			"p": 123.4, "zde": 1.2, "zdf": 0.98, "vol": 1000.0,
			"cje": 123400.0, "ccl": 3000.0, "xqj": 4000.0,
			"syr": 12.0, "rz": 0.5, "zjsj": 122.0, "o": 121.0,
		}},
	}}

	rows, err := GetCFFEXOptionQuotes(context.Background(), client, "https://em.test/option", CFFEXOptionQuotesOptions{PageSize: 100})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "io2501C4000" || rows[0].Name != "沪深300购2501" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Price == nil || *rows[0].Price != 123.4 || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 3000 {
		t.Fatalf("row numbers = %+v", rows[0])
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "orderBy", "zdf")
	assertQuery(t, query, "sort", "desc")
	assertQuery(t, query, "pageSize", "100")
	assertQuery(t, query, "pageIndex", "0")
	assertQuery(t, query, "token", emFuturesGlobalSpotToken)
	assertQuery(t, query, "field", "dm,sc,name,p,zsjd,zde,zdf,f152,vol,cje,ccl,xqj,syr,rz,zjsj,o")
}

func TestGetCFFEXOptionQuotesReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeOptionsClient{payload: map[string]any{
		"total": float64(1),
		"list":  map[string]any{"dm": "io2501C4000"},
	}}

	rows, err := GetCFFEXOptionQuotes(context.Background(), client, "https://em.test/option", CFFEXOptionQuotesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}

func TestGetOptionLHBBuildsRequestAndParsesRows(t *testing.T) {
	client := &fakeOptionsClient{payload: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"TRADE_TYPE": "认购", "TRADE_DATE": "2024-12-16 00:00:00",
				"SECURITY_CODE": "510050", "TARGET_NAME": "上证50ETF",
				"MEMBER_RANK": 1.0, "MEMBER_NAME_ABBR": "期货公司A",
				"SELL_VOLUME": 100.0, "SELL_VOLUME_CHANGE": 10.0,
				"NET_SELL_VOLUME": 20.0, "SELL_VOLUME_RATIO": 1.5,
				"BUY_VOLUME": 200.0, "BUY_VOLUME_CHANGE": 15.0,
				"NET_BUY_VOLUME": 30.0, "BUY_VOLUME_RATIO": 2.5,
				"SELL_POSITION": 300.0, "SELL_POSITION_CHANGE": -5.0,
				"NET_SELL_POSITION": 40.0, "SELL_POSITION_RATIO": 3.5,
				"BUY_POSITION": 400.0, "BUY_POSITION_CHANGE": 6.0,
				"NET_BUY_POSITION": 50.0, "BUY_POSITION_RATIO": 4.5,
			}},
		},
	}}

	rows, err := GetOptionLHB(context.Background(), client, "https://em.test/lhb", "510050", "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Symbol != "510050" || rows[0].Date != "2024-12-16" || rows[0].MemberName != "期货公司A" {
		t.Fatalf("rows = %+v", rows)
	}
	if rows[0].Rank != 1 || rows[0].BuyVolume == nil || *rows[0].BuyVolume != 200 || rows[0].BuyPositionRatio == nil || *rows[0].BuyPositionRatio != 4.5 {
		t.Fatalf("row numbers = %+v", rows[0])
	}

	parsed, err := url.Parse(client.lastURL)
	if err != nil {
		t.Fatal(err)
	}
	query := parsed.Query()
	assertQuery(t, query, "type", "RPT_IF_BILLBOARD_TD")
	assertQuery(t, query, "sty", "ALL")
	assertQuery(t, query, "p", "1")
	assertQuery(t, query, "ps", "200")
	assertQuery(t, query, "source", "IFBILLBOARD")
	assertQuery(t, query, "client", "WEB")
	assertQuery(t, query, "ut", emDataToken)
	assertQuery(t, query, "filter", `(SECURITY_CODE="510050")(TRADE_DATE='2024-12-16')`)
}

func TestGetOptionLHBReturnsEmptyRowsForNonArrayPayload(t *testing.T) {
	client := &fakeOptionsClient{payload: map[string]any{
		"result": map[string]any{
			"data": map[string]any{"SECURITY_CODE": "510050"},
		},
	}}

	rows, err := GetOptionLHB(context.Background(), client, "https://em.test/lhb", "510050", "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("rows = %+v, want empty", rows)
	}
}
