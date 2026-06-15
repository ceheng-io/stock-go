package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	stock "github.com/ceheng.io/stock-go"
)

type fakeSDK struct {
	searchKeyword           string
	fullQuoteCodes          []string
	industryConstituentCode string
	historySymbol           string
	historyOptions          stock.HistoryKlineOptions
	indicatorSymbol         string
	indicatorOptions        stock.KlineWithIndicatorsOptions
	sectorFundFlowSymbol    string
	sectorFundFlowOptions   stock.FundFlowOptions
	ztPoolType              stock.ZTPoolType
	ztPoolArgCount          int
}

func (f *fakeSDK) Search(_ context.Context, keyword string) ([]stock.SearchResult, error) {
	f.searchKeyword = keyword
	return []stock.SearchResult{{Code: "600519", Name: "贵州茅台", Market: "sh", Type: "股票", Category: stock.SearchStock}}, nil
}

func (f *fakeSDK) GetFullQuotes(_ context.Context, codes []string) ([]stock.FullQuote, error) {
	f.fullQuoteCodes = codes
	return []stock.FullQuote{{Code: codes[0], Name: "贵州茅台", Price: 1688.88}}, nil
}

func (f *fakeSDK) GetIndustryList(context.Context) ([]stock.Board, error) {
	return []stock.Board{{Code: "BK0475", Name: "酿酒行业", Rank: 1}}, nil
}

func (f *fakeSDK) GetIndustryConstituents(_ context.Context, code string) ([]stock.BoardConstituent, error) {
	f.industryConstituentCode = code
	return []stock.BoardConstituent{{Code: "600519", Name: "贵州茅台", Rank: 1}}, nil
}

func (f *fakeSDK) GetHistoryKline(_ context.Context, symbol string, options ...stock.HistoryKlineOptions) ([]stock.HistoryKline, error) {
	f.historySymbol = symbol
	if len(options) > 0 {
		f.historyOptions = options[0]
	}
	return []stock.HistoryKline{{Code: symbol, Date: "2026-06-15"}}, nil
}

func (f *fakeSDK) GetKlineWithIndicators(_ context.Context, symbol string, options ...stock.KlineWithIndicatorsOptions) ([]stock.KlineWithIndicators, error) {
	f.indicatorSymbol = symbol
	if len(options) > 0 {
		f.indicatorOptions = options[0]
	}
	return []stock.KlineWithIndicators{{Date: "2026-06-15"}}, nil
}

func (f *fakeSDK) GetSectorFundFlowHistory(_ context.Context, symbol string, options ...stock.FundFlowOptions) ([]stock.StockFundFlow, error) {
	f.sectorFundFlowSymbol = symbol
	if len(options) > 0 {
		f.sectorFundFlowOptions = options[0]
	}
	return []stock.StockFundFlow{{Date: "2026-06-15"}}, nil
}

func (f *fakeSDK) GetZTPool(_ context.Context, args ...any) ([]stock.ZTPoolItem, error) {
	f.ztPoolArgCount = len(args)
	for _, arg := range args {
		if value, ok := arg.(stock.ZTPoolType); ok {
			f.ztPoolType = value
		}
	}
	return []stock.ZTPoolItem{{Code: "600000", Name: "浦发银行"}}, nil
}

func TestHealthReturnsOK(t *testing.T) {
	server := NewServer(nil)
	request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body map[string]bool
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body["ok"] {
		t.Fatalf("expected ok=true, got %#v", body)
	}
}

func TestOptionsRequestReturnsCORSHeaders(t *testing.T) {
	server := NewServer(nil)
	request := httptest.NewRequest(http.MethodOptions, "/api/search", nil)
	request.Header.Set("Origin", "http://localhost:5173")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Fatalf("expected allow-origin *, got %q", got)
	}
	if got := response.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, http.MethodGet) {
		t.Fatalf("expected allow methods to contain GET, got %q", got)
	}
}

func TestSearchRequiresKeyword(t *testing.T) {
	server := NewServer(nil)
	request := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}

	var body struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !strings.Contains(body.Error.Message, "keyword") {
		t.Fatalf("expected keyword error, got %q", body.Error.Message)
	}
}

func TestSearchCallsSDK(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/search?keyword=%E8%8C%85%E5%8F%B0", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.searchKeyword != "茅台" {
		t.Fatalf("expected keyword 茅台, got %q", fake.searchKeyword)
	}
	var body []stock.SearchResult
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body) != 1 || body[0].Code != "600519" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestFullQuotesParsesCodes(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/quotes/full?codes=sh600519,sz000001", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if got := strings.Join(fake.fullQuoteCodes, ","); got != "sh600519,sz000001" {
		t.Fatalf("unexpected codes: %q", got)
	}
}

func TestIndustryListEndpoint(t *testing.T) {
	server := NewServer(&fakeSDK{})
	request := httptest.NewRequest(http.MethodGet, "/api/boards/industry", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	var body []stock.Board
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body) != 1 || body[0].Code != "BK0475" {
		t.Fatalf("unexpected body: %#v", body)
	}
}

func TestIndustryConstituentsEndpoint(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/boards/industry/BK0475/constituents", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.industryConstituentCode != "BK0475" {
		t.Fatalf("unexpected board code: %q", fake.industryConstituentCode)
	}
}

func TestHistoryKlineEndpoint(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/kline/history?symbol=sh600519&period=daily&adjust=qfq", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.historySymbol != "sh600519" {
		t.Fatalf("unexpected symbol: %q", fake.historySymbol)
	}
	if fake.historyOptions.Period != "daily" || fake.historyOptions.Adjust != "qfq" {
		t.Fatalf("unexpected options: %#v", fake.historyOptions)
	}
}

func TestKlineIndicatorsEndpointUsesDefaultIndicatorSet(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/kline/indicators?symbol=sh600519&period=daily&adjust=qfq", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.indicatorSymbol != "sh600519" {
		t.Fatalf("unexpected symbol: %q", fake.indicatorSymbol)
	}
	options := fake.indicatorOptions
	if options.Period != stock.KlinePeriod("daily") || options.Adjust != stock.AdjustType("qfq") {
		t.Fatalf("unexpected options: %#v", options)
	}
	indicators := options.Indicators
	if indicators.MA == nil || len(indicators.MA.Periods) == 0 ||
		indicators.MACD == nil ||
		indicators.BOLL == nil ||
		indicators.KDJ == nil ||
		indicators.RSI == nil || len(indicators.RSI.Periods) == 0 ||
		indicators.OBV == nil ||
		indicators.ROC == nil ||
		indicators.DMI == nil ||
		indicators.SAR == nil ||
		indicators.KC == nil {
		t.Fatalf("expected default indicator set, got %#v", indicators)
	}
}

func TestSectorFundFlowHistoryEndpointPassesSymbolAndPeriod(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/fund-flow/sector-history?symbol=BK0475&period=daily", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.sectorFundFlowSymbol != "BK0475" {
		t.Fatalf("unexpected symbol: %q", fake.sectorFundFlowSymbol)
	}
	if fake.sectorFundFlowOptions.Period != stock.FundFlowPeriod("daily") {
		t.Fatalf("unexpected options: %#v", fake.sectorFundFlowOptions)
	}
}

func TestZTPoolEndpoint(t *testing.T) {
	fake := &fakeSDK{}
	server := NewServer(fake)
	request := httptest.NewRequest(http.MethodGet, "/api/market-event/zt-pool?type=strong", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
	if fake.ztPoolType != stock.ZTPoolStrong {
		t.Fatalf("unexpected pool type: %q", fake.ztPoolType)
	}
	if fake.ztPoolArgCount != 1 {
		t.Fatalf("expected only type argument, got %d", fake.ztPoolArgCount)
	}
}
