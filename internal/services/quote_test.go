package services

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"testing"

	"github.com/ceheng.io/stock-go/internal/core"
	"github.com/ceheng.io/stock-go/internal/providers/tencent"
)

type quoteClientStub struct {
	mu          sync.Mutex
	params      string
	text        string
	lists       map[string][]string
	json        map[string]any
	url         string
	calendarURL string
	textCalls   int
	textEntered chan struct{}
	releaseText chan struct{}
}

func (q *quoteClientStub) GetTencentQuote(_ context.Context, params string) ([]core.TencentQuoteItem, error) {
	q.params = params
	if strings.HasPrefix(params, "s_pk") {
		return []core.TencentQuoteItem{{Key: "s_pksh600519", Fields: []string{"10", "20", "30", "40"}}}, nil
	}
	if strings.HasPrefix(params, "s_") {
		return []core.TencentQuoteItem{
			{
				Key:    "s_sh600519",
				Fields: []string{"1", "贵州茅台", "600519", "1700.00", "-1.23", "-0.07", "12345", "67890", "", "25000", "GP-A"},
			},
		}, nil
	}
	if strings.HasPrefix(params, "hk") {
		fields := make([]string, 50)
		fields[0] = "100"
		fields[1] = "腾讯控股"
		fields[2] = "00700"
		fields[3] = "390.20"
		fields[47] = "HKD"
		return []core.TencentQuoteItem{{Key: "hk00700", Fields: fields}}, nil
	}
	if strings.HasPrefix(params, "us") {
		fields := make([]string, 50)
		fields[0] = "200"
		fields[1] = "APPLE"
		fields[2] = "AAPL"
		fields[3] = "205.50"
		return []core.TencentQuoteItem{{Key: "usAAPL", Fields: fields}}, nil
	}
	if strings.HasPrefix(params, "jj") {
		fields := make([]string, 9)
		fields[0] = "110011"
		fields[1] = "易方达中小盘"
		fields[5] = "3.5000"
		fields[8] = "2024-05-10"
		return []core.TencentQuoteItem{{Key: "jj110011", Fields: fields}}, nil
	}
	if strings.HasPrefix(params, "ff_") {
		fields := make([]string, 14)
		fields[0] = "600519"
		fields[1] = "100"
		fields[2] = "80"
		fields[3] = "20"
		fields[12] = "贵州茅台"
		fields[13] = "20240613"
		return []core.TencentQuoteItem{{Key: "ff_sh600519", Fields: fields}}, nil
	}
	if strings.Contains(params, "sh600000") || strings.Contains(params, "sz000001") {
		return []core.TencentQuoteItem{
			fullQuoteItem("sh600000", "浦发银行", "600000", "8.88"),
			fullQuoteItem("sz000001", "平安银行", "000001", "10.10"),
		}, nil
	}
	fields := make([]string, 80)
	fields[0] = "1"
	fields[1] = "贵州茅台"
	fields[2] = "600519"
	fields[3] = "1700.00"
	return []core.TencentQuoteItem{
		{
			Key:    "sh600519",
			Fields: fields,
		},
	}, nil
}

func fullQuoteItem(key string, name string, code string, price string) core.TencentQuoteItem {
	fields := make([]string, 80)
	fields[0] = "1"
	fields[1] = name
	fields[2] = code
	fields[3] = price
	return core.TencentQuoteItem{Key: key, Fields: fields}
}

func (q *quoteClientStub) GetText(_ context.Context, requestURL string) (string, error) {
	q.mu.Lock()
	q.url = requestURL
	q.textCalls++
	if q.textEntered != nil {
		select {
		case q.textEntered <- struct{}{}:
		default:
		}
	}
	q.mu.Unlock()
	if q.releaseText != nil {
		<-q.releaseText
	}
	if q.text != "" {
		return q.text, nil
	}
	return `v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A";`, nil
}

func (q *quoteClientStub) TencentSearchURL(keyword string) string {
	return "https://smartbox.test/s3/?q=" + keyword
}

func (q *quoteClientStub) CalendarURL() string {
	if q.calendarURL != "" {
		return q.calendarURL
	}
	return "https://calendar.test"
}

func (q *quoteClientStub) GetJSON(_ context.Context, requestURL string, target any) error {
	q.url = requestURL
	if q.json != nil {
		body, _ := json.Marshal(q.json)
		return json.Unmarshal(body, target)
	}
	payload := struct {
		Success bool     `json:"success"`
		List    []string `json:"list"`
	}{Success: true, List: q.lists[requestURL]}
	b, _ := json.Marshal(payload)
	return json.Unmarshal(b, target)
}

func (q *quoteClientStub) AShareListURL() string { return "a" }
func (q *quoteClientStub) USListURL() string     { return "us" }
func (q *quoteClientStub) HKListURL() string     { return "hk" }
func (q *quoteClientStub) FundListURL() string   { return "fund" }

func resetTradingCalendarCacheForTest() {
	tradingCalendarCache.mu.Lock()
	tradingCalendarCache.entries = make(map[string]*calendarCacheEntry)
	tradingCalendarCache.mu.Unlock()
}

func TestQuoteServiceSimpleCN(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	quotes, err := service.SimpleCN(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if client.params != "s_sh600519" {
		t.Fatalf("params = %q, want s_sh600519", client.params)
	}
	if len(quotes) != 1 || quotes[0].Code != "600519" {
		t.Fatalf("quotes = %+v", quotes)
	}
}

func TestQuoteServiceFullCN(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	quotes, err := service.CN(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if client.params != "sh600519" {
		t.Fatalf("params = %q, want sh600519", client.params)
	}
	if len(quotes) != 1 || quotes[0].Code != "600519" {
		t.Fatalf("quotes = %+v", quotes)
	}
}

func TestQuoteServiceHKUSFund(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	if _, err := service.HK(context.Background(), []string{"00700"}); err != nil {
		t.Fatal(err)
	}
	if client.params != "hk00700" {
		t.Fatalf("HK params = %q, want hk00700", client.params)
	}

	if _, err := service.US(context.Background(), []string{"AAPL"}); err != nil {
		t.Fatal(err)
	}
	if client.params != "usAAPL" {
		t.Fatalf("US params = %q, want usAAPL", client.params)
	}

	if _, err := service.Fund(context.Background(), []string{"110011"}); err != nil {
		t.Fatal(err)
	}
	if client.params != "jj110011" {
		t.Fatalf("Fund params = %q, want jj110011", client.params)
	}
}

func TestQuoteServiceSearch(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	results, err := service.Search(context.Background(), "茅台")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "sh600519" || results[0].Name != "贵州茅台" {
		t.Fatalf("results = %+v", results)
	}
}

func TestQuoteServiceTradingCalendar(t *testing.T) {
	resetTradingCalendarCacheForTest()
	client := &quoteClientStub{text: "1990-12-19,1990-12-20"}
	service := NewQuoteService(client)

	calendar, err := service.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(calendar) != 2 || calendar[0] != "1990-12-19" {
		t.Fatalf("calendar = %#v", calendar)
	}
}

func TestQuoteServiceTradingCalendarCachesAndReturnsCopy(t *testing.T) {
	resetTradingCalendarCacheForTest()
	client := &quoteClientStub{text: "1990-12-19,1990-12-20"}
	service := NewQuoteService(client)

	first, err := service.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	first[0] = "mutated"

	second, err := service.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if client.textCalls != 1 {
		t.Fatalf("textCalls = %d, want 1", client.textCalls)
	}
	if second[0] != "1990-12-19" {
		t.Fatalf("second calendar = %#v, want cached copy unaffected by caller mutation", second)
	}
}

func TestQuoteServiceTradingCalendarCoalescesConcurrentColdRequests(t *testing.T) {
	resetTradingCalendarCacheForTest()
	client := &quoteClientStub{
		text:        "1990-12-19,1990-12-20",
		textEntered: make(chan struct{}, 1),
		releaseText: make(chan struct{}),
	}
	service := NewQuoteService(client)

	const workers = 8
	results := make(chan []string, workers)
	errs := make(chan error, workers)
	for range workers {
		go func() {
			calendar, err := service.TradingCalendar(context.Background())
			if err != nil {
				errs <- err
				return
			}
			results <- calendar
		}()
	}

	<-client.textEntered
	close(client.releaseText)

	for range workers {
		select {
		case err := <-errs:
			t.Fatal(err)
		case calendar := <-results:
			if len(calendar) != 2 || calendar[0] != "1990-12-19" {
				t.Fatalf("calendar = %#v", calendar)
			}
		}
	}
	client.mu.Lock()
	textCalls := client.textCalls
	client.mu.Unlock()
	if textCalls != 1 {
		t.Fatalf("textCalls = %d, want 1", textCalls)
	}
}

func TestQuoteServiceTradingCalendarUsesSharedCacheAcrossServices(t *testing.T) {
	resetTradingCalendarCacheForTest()
	client := &quoteClientStub{text: "1990-12-19,1990-12-20"}
	firstService := NewQuoteService(client)
	secondService := NewQuoteService(client)

	first, err := firstService.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := secondService.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if client.textCalls != 1 {
		t.Fatalf("textCalls = %d, want 1", client.textCalls)
	}
	if len(first) != 2 || len(second) != 2 || second[0] != "1990-12-19" {
		t.Fatalf("first=%#v second=%#v", first, second)
	}
}

func TestQuoteServiceTradingCalendarCacheIsScopedByCalendarURL(t *testing.T) {
	resetTradingCalendarCacheForTest()
	firstClient := &quoteClientStub{
		text:        "1990-12-19",
		calendarURL: "https://calendar.test/one",
	}
	secondClient := &quoteClientStub{
		text:        "2024-06-13",
		calendarURL: "https://calendar.test/two",
	}
	firstService := NewQuoteService(firstClient)
	secondService := NewQuoteService(secondClient)

	first, err := firstService.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	second, err := secondService.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if firstClient.textCalls != 1 || secondClient.textCalls != 1 {
		t.Fatalf("textCalls first=%d second=%d, want 1 each", firstClient.textCalls, secondClient.textCalls)
	}
	if first[0] != "1990-12-19" || second[0] != "2024-06-13" {
		t.Fatalf("first=%#v second=%#v", first, second)
	}
}

func TestQuoteServiceCodes(t *testing.T) {
	client := &quoteClientStub{lists: map[string][]string{
		"a":  {"sh600000", "sz000001", "bj830799"},
		"us": {"105.AAPL", "106.BABA"},
		"hk": {"00700"},
	}, text: ",110011"}
	service := NewQuoteService(client)

	cn, err := service.CodesCN(context.Background(), tencent.CodeListOptions{Market: tencent.AShareMarketBJ})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0] != "bj830799" {
		t.Fatalf("cn = %#v", cn)
	}
	us, err := service.CodesUS(context.Background(), tencent.USCodeListOptions{Simple: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 2 || us[0] != "AAPL" {
		t.Fatalf("us = %#v", us)
	}
	hk, err := service.CodesHK(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0] != "00700" {
		t.Fatalf("hk = %#v", hk)
	}
	fund, err := service.CodesFund(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(fund) != 1 || fund[0] != "110011" {
		t.Fatalf("fund = %#v", fund)
	}
}

func TestQuoteServiceBatchByCodes(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	cn, err := service.BatchCN(context.Background(), []string{"sh600519"}, tencent.BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0].Code != "600519" {
		t.Fatalf("cn = %+v", cn)
	}
}

func TestQuoteServiceAllQuotesFetchesCodeListsThenBatches(t *testing.T) {
	client := &quoteClientStub{lists: map[string][]string{
		"a":  {"sh600000", "sz000001", "bj830799"},
		"hk": {"00700"},
		"us": {"105.AAPL", "106.BABA"},
	}}
	service := NewQuoteService(client)

	var progress []int
	cn, err := service.AllCN(context.Background(), tencent.CodeListOptions{Market: tencent.AShareMarketSZ}, tencent.BatchOptions{
		BatchSize:  1,
		OnProgress: func(completed, total int) { progress = append(progress, completed*10+total) },
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0].Code != "000001" {
		t.Fatalf("cn = %+v", cn)
	}
	if client.url != "a" {
		t.Fatalf("last code-list url = %q, want a", client.url)
	}
	if len(progress) != 1 || progress[0] != 11 {
		t.Fatalf("progress = %#v, want [11]", progress)
	}

	hk, err := service.AllHK(context.Background(), tencent.BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Code != "00700" {
		t.Fatalf("hk = %+v", hk)
	}

	us, err := service.AllUS(context.Background(), tencent.USCodeListOptions{Market: tencent.USMarketNASDAQ}, tencent.BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Code != "AAPL" {
		t.Fatalf("us = %+v", us)
	}
}

func TestQuoteServiceBatchRaw(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	items, err := service.BatchRaw(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if client.params != "sh600519" {
		t.Fatalf("params = %q, want sh600519", client.params)
	}
	if len(items) != 1 || items[0].Key != "sh600519" || items[0].Fields[1] != "贵州茅台" {
		t.Fatalf("items = %+v", items)
	}
}

func TestQuoteServiceDividendDetail(t *testing.T) {
	client := &quoteClientStub{json: map[string]any{
		"result": map[string]any{
			"data": []map[string]any{{
				"SECURITY_CODE":      "600519",
				"SECURITY_NAME_ABBR": "贵州茅台",
				"REPORT_DATE":        "2024-12-31",
				"PRETAX_BONUS_RMB":   30.0,
			}},
		},
	}}
	service := NewQuoteService(client, QuoteURLs{Datacenter: "https://em.test/datacenter"})

	rows, err := service.DividendDetail(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if client.url == "" {
		t.Fatal("expected DividendDetail to call datacenter endpoint")
	}
	if len(rows) != 1 || rows[0].Code != "600519" || rows[0].DividendPretax == nil || *rows[0].DividendPretax != 30 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestQuoteServiceFundFlowAndPanelLargeOrder(t *testing.T) {
	client := &quoteClientStub{}
	service := NewQuoteService(client)

	flow, err := service.FundFlow(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if client.params != "ff_sh600519" {
		t.Fatalf("FundFlow params = %q, want ff_sh600519", client.params)
	}
	if len(flow) != 1 || flow[0].Code != "600519" || flow[0].MainNet != 20 {
		t.Fatalf("flow = %+v", flow)
	}

	panel, err := service.PanelLargeOrder(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if client.params != "s_pksh600519" {
		t.Fatalf("PanelLargeOrder params = %q, want s_pksh600519", client.params)
	}
	if len(panel) != 1 || panel[0].SellSmallRatio != 40 {
		t.Fatalf("panel = %+v", panel)
	}
}

func TestQuoteServiceTodayTimeline(t *testing.T) {
	client := &quoteClientStub{text: `{
		"code": 0,
		"data": {
			"sz000001": {
				"data": {"date": "20240613", "data": ["0931 10.10 100 101000.00"]},
				"qt": {"sz000001": ["", "", "", "", "9.90"]}
			}
		}
	}`}
	service := NewQuoteService(client, QuoteURLs{Minute: "https://minute.test/query"})

	row, err := service.TodayTimeline(context.Background(), "sz000001")
	if err != nil {
		t.Fatal(err)
	}
	if client.url != "https://minute.test/query?code=sz000001" {
		t.Fatalf("url = %q", client.url)
	}
	if row.Code != "sz000001" || row.PreClose == nil || *row.PreClose != 9.9 || len(row.Data) != 1 {
		t.Fatalf("row = %+v", row)
	}
}
