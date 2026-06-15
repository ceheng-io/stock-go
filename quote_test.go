package stock

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func TestClientQuotesSimpleCN(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "s_sh600519" {
			t.Fatalf("q = %q, want s_sh600519", r.URL.Query().Get("q"))
		}
		writeGBK(t, w, `v_s_sh600519="1~贵州茅台~600519~1700.00~-1.23~-0.07~12345~67890~~25000~GP-A";`)
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	quotes, err := client.Quotes.SimpleCN(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
	if quotes[0].Name != "贵州茅台" || quotes[0].Code != "600519" {
		t.Fatalf("quote = %+v", quotes[0])
	}
}

func TestClientQuotesSimpleCNUsesRetryOption(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
			return
		}
		writeGBK(t, w, `v_s_sh600519="1~贵州茅台~600519~1700.00~-1.23~-0.07~12345~67890~~25000~GP-A";`)
	}))
	defer server.Close()

	client := New(
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithRetry(RetryOptions{MaxRetries: 2, BaseDelay: time.Nanosecond}),
	)
	quotes, err := client.Quotes.SimpleCN(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Code != "600519" || attempts != 3 {
		t.Fatalf("quotes=%+v attempts=%d", quotes, attempts)
	}
}

func TestClientQuotesSimpleCNReturnsCircuitBreakerOpenError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		http.Error(w, "temporary", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := New(
		WithBaseURL(server.URL),
		WithHTTPClient(server.Client()),
		WithRetry(RetryOptions{MaxRetries: 0}),
		WithCircuitBreaker(CircuitBreakerOptions{FailureThreshold: 1, ResetTimeout: time.Minute}),
	)
	_, err := client.Quotes.SimpleCN(context.Background(), []string{"sh600519"})
	if err == nil {
		t.Fatal("expected first request to fail")
	}
	_, err = client.Quotes.SimpleCN(context.Background(), []string{"sh600519"})
	if !errors.Is(err, ErrCircuitBreakerOpen) {
		t.Fatalf("error = %v, want ErrCircuitBreakerOpen", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1", attempts)
	}
}

func TestClientCodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var list []string
		switch r.URL.Path {
		case "/a.json":
			list = []string{"sh600000", "sz000001", "bj830799"}
		case "/us.json":
			list = []string{"105.AAPL", "106.BABA"}
		case "/hk.json":
			list = []string{"00700"}
		case "/fund.json":
			_, _ = io.WriteString(w, ",110011")
			return
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": list})
	}))
	defer server.Close()

	client := New(
		WithAShareListURL(server.URL+"/a.json"),
		WithUSListURL(server.URL+"/us.json"),
		WithHKListURL(server.URL+"/hk.json"),
		WithFundListURL(server.URL+"/fund.json"),
		WithHTTPClient(server.Client()),
	)
	cn, err := client.Quotes.CodesCN(context.Background(), CodeListOptions{Market: AShareMarketBJ})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0] != "bj830799" {
		t.Fatalf("cn = %#v", cn)
	}
	us, err := client.Quotes.CodesUS(context.Background(), USCodeListOptions{Simple: true})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 2 || us[0] != "AAPL" {
		t.Fatalf("us = %#v", us)
	}
	fund, err := client.Quotes.CodesFund(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(fund) != 1 || fund[0] != "110011" {
		t.Fatalf("fund = %#v", fund)
	}
}

func TestClientBatchCN(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "sh600519" {
			t.Fatalf("q = %q, want sh600519", r.URL.Query().Get("q"))
		}
		fields := make([]string, 80)
		fields[0] = "1"
		fields[1] = "贵州茅台"
		fields[2] = "600519"
		fields[3] = "1700.00"
		writeGBK(t, w, `v_sh600519="`+joinFields(fields)+`";`)
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	quotes, err := client.Quotes.BatchCN(context.Background(), []string{"sh600519"}, BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Code != "600519" {
		t.Fatalf("quotes = %+v", quotes)
	}
}

func TestClientBatchInvalidOptionsReturnInvalidArgument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/hk.json":
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": []string{"00700"}})
		case "/us.json":
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": []string{"105.AAPL"}})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithHKListURL(server.URL+"/hk.json"),
		WithUSListURL(server.URL+"/us.json"),
		WithHTTPClient(server.Client()),
	)
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service cn negative batch size",
			call: func() error {
				_, err := client.Quotes.BatchCN(context.Background(), []string{"sh600519"}, BatchOptions{BatchSize: -1})
				return err
			},
		},
		{
			name: "wrapper cn negative batch size",
			call: func() error {
				_, err := client.GetAllQuotesByCodes(context.Background(), []string{"sh600519"}, GetAllAShareQuotesOptions{BatchSize: -1})
				return err
			},
		},
		{
			name: "service hk negative concurrency",
			call: func() error {
				_, err := client.Quotes.BatchHK(context.Background(), []string{"00700"}, BatchOptions{Concurrency: -1})
				return err
			},
		},
		{
			name: "wrapper hk negative concurrency",
			call: func() error {
				_, err := client.GetAllHKShareQuotes(context.Background(), GetAllHKQuotesOptions{Concurrency: -1})
				return err
			},
		},
		{
			name: "service us negative batch size",
			call: func() error {
				_, err := client.Quotes.BatchUS(context.Background(), []string{"AAPL"}, BatchOptions{BatchSize: -1})
				return err
			},
		},
		{
			name: "wrapper us negative concurrency",
			call: func() error {
				_, err := client.GetAllUSShareQuotes(context.Background(), GetAllUSQuotesOptions{Concurrency: -1})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientQuotesAllCNHKUS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a.json":
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": []string{"sh600000", "sz000001"}})
		case "/hk.json":
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": []string{"00700"}})
		case "/us.json":
			_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "list": []string{"105.AAPL", "106.BABA"}})
		case "/":
			switch r.URL.Query().Get("q") {
			case "sz000001":
				fields := make([]string, 80)
				fields[0] = "1"
				fields[1] = "平安银行"
				fields[2] = "000001"
				fields[3] = "10.10"
				writeGBK(t, w, `v_sz000001="`+joinFields(fields)+`";`)
			case "hk00700":
				fields := make([]string, 50)
				fields[0] = "100"
				fields[1] = "腾讯控股"
				fields[2] = "00700"
				fields[3] = "390.20"
				fields[47] = "HKD"
				writeGBK(t, w, `v_hk00700="`+joinFields(fields)+`";`)
			case "usAAPL":
				fields := make([]string, 50)
				fields[0] = "200"
				fields[1] = "APPLE"
				fields[2] = "AAPL"
				fields[3] = "205.50"
				writeGBK(t, w, `v_usAAPL="`+joinFields(fields)+`";`)
			default:
				t.Fatalf("unexpected q = %q", r.URL.Query().Get("q"))
			}
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithBaseURL(server.URL),
		WithAShareListURL(server.URL+"/a.json"),
		WithHKListURL(server.URL+"/hk.json"),
		WithUSListURL(server.URL+"/us.json"),
		WithHTTPClient(server.Client()),
	)

	cn, err := client.Quotes.AllCN(context.Background(), CodeListOptions{Market: AShareMarketSZ}, BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(cn) != 1 || cn[0].Code != "000001" {
		t.Fatalf("cn = %+v", cn)
	}
	hk, err := client.Quotes.AllHK(context.Background(), BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Code != "00700" {
		t.Fatalf("hk = %+v", hk)
	}
	us, err := client.Quotes.AllUS(context.Background(), USCodeListOptions{Market: USMarketNASDAQ}, BatchOptions{BatchSize: 1})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Code != "AAPL" {
		t.Fatalf("us = %+v", us)
	}
}

func TestClientQuotesFundFlowAndPanelLargeOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("q") {
		case "ff_sh600519":
			fields := make([]string, 14)
			fields[0] = "600519"
			fields[1] = "100"
			fields[2] = "80"
			fields[3] = "20"
			fields[12] = "贵州茅台"
			fields[13] = "20240613"
			writeGBK(t, w, `v_ff_sh600519="`+joinFields(fields)+`";`)
		case "s_pksh600519":
			_, _ = w.Write([]byte(`v_s_pksh600519="10~20~30~40";`))
		default:
			t.Fatalf("q = %q", r.URL.Query().Get("q"))
		}
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	flow, err := client.Quotes.FundFlow(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(flow) != 1 || flow[0].Code != "600519" || flow[0].MainNet != 20 {
		t.Fatalf("flow = %+v", flow)
	}
	panel, err := client.Quotes.PanelLargeOrder(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(panel) != 1 || panel[0].BuyLargeRatio != 10 || panel[0].SellSmallRatio != 40 {
		t.Fatalf("panel = %+v", panel)
	}
}

func TestClientQuotesTodayTimeline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/minute" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("code") != "sz000001" {
			t.Fatalf("code = %q", r.URL.Query().Get("code"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": float64(0),
			"data": map[string]any{
				"sz000001": map[string]any{
					"data": map[string]any{
						"date": "20240613",
						"data": []string{"0931 10.10 100 101000.00"},
					},
					"qt": map[string]any{
						"sz000001": []string{"", "", "", "", "9.90"},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := New(WithTencentMinuteURL(server.URL+"/minute"), WithHTTPClient(server.Client()))
	timeline, err := client.Quotes.TodayTimeline(context.Background(), "sz000001")
	if err != nil {
		t.Fatal(err)
	}
	if timeline.Code != "sz000001" || timeline.PreClose == nil || *timeline.PreClose != 9.9 || len(timeline.Data) != 1 {
		t.Fatalf("timeline = %+v", timeline)
	}
	if timeline.Data[0].Volume != 10000 || timeline.Data[0].AvgPrice != 10.1 {
		t.Fatalf("tick = %+v", timeline.Data[0])
	}
}

func TestClientKlineCN(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/kline" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("secid") != "1.600519" {
			t.Fatalf("secid = %q, want 1.600519", r.URL.Query().Get("secid"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"klines": []string{"2024-12-16,1500,1510,1520,1490,12345,67890000,2,1.5,22.3,0.5"},
			},
		})
	}))
	defer server.Close()

	client := New(WithEastmoneyKlineURL(server.URL+"/kline"), WithHTTPClient(server.Client()))
	rows, err := client.Kline.CN(context.Background(), "sh600519", HistoryKlineOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "600519" || rows[0].Close == nil || *rows[0].Close != 1510 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientKlineCNUsesProviderRetryPolicy(t *testing.T) {
	transport := &sequenceRoundTripper{
		responses: []sequenceResponse{
			{status: http.StatusServiceUnavailable, body: "temporary"},
			{status: http.StatusOK, body: `{"data":{"klines":["2024-12-16,1500,1510,1520,1490,12345,67890000,2,1.5,22.3,0.5"]}}`},
		},
	}

	client := New(
		WithEastmoneyKlineURL("https://push2his.eastmoney.com/api/qt/stock/kline/get"),
		WithHTTPClient(&http.Client{Transport: transport}),
		WithRetry(RetryOptions{MaxRetries: 0}),
		WithProviderPolicy(ProviderEastmoney, ProviderPolicy{
			Retry: &RetryOptions{MaxRetries: 1, BaseDelay: time.Nanosecond},
		}),
	)
	rows, err := client.Kline.CN(context.Background(), "sh600519", HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || transport.calls != 2 {
		t.Fatalf("rows=%+v attempts=%d, want one row after two attempts", rows, transport.calls)
	}
}

func TestClientKlineCNMinute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("klt") != "5" {
			t.Fatalf("klt = %q, want 5", r.URL.Query().Get("klt"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"klines": []string{"2024-12-16 09:35,1500,1510,1520,1490,12345,67890000,2,1.5,22.3,0.5"},
			},
		})
	}))
	defer server.Close()

	client := New(WithEastmoneyKlineURL(server.URL), WithHTTPClient(server.Client()))
	rows, err := client.Kline.CNMinute(context.Background(), "sh600519", MinuteKlineOptions{Period: MinutePeriodFive})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows.Klines) != 1 || rows.Klines[0].Close == nil || *rows.Klines[0].Close != 1510 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientKlineHKUS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/hk":
			if r.URL.Query().Get("secid") != "116.00700" {
				t.Fatalf("hk secid = %q, want 116.00700", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"code":   "00700",
					"name":   "腾讯控股",
					"klines": []string{"2024-12-16,390,392,395,388,12345,67890000,1.8,0.51,2,0.03"},
				},
			})
		case "/us":
			if r.URL.Query().Get("secid") != "105.AAPL" {
				t.Fatalf("us secid = %q, want 105.AAPL", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"code":   "AAPL",
					"name":   "Apple Inc.",
					"klines": []string{"2024-12-16,250,251.5,253,249,34567,87650000,1.6,0.8,2,0.1"},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyHKKlineURL(server.URL+"/hk"),
		WithEastmoneyUSKlineURL(server.URL+"/us"),
		WithHTTPClient(server.Client()),
	)
	hk, err := client.Kline.HK(context.Background(), "00700", HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Name != "腾讯控股" || hk[0].Currency != "HKD" {
		t.Fatalf("hk = %+v", hk)
	}
	us, err := client.Kline.US(context.Background(), "105.AAPL", HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Name != "Apple Inc." || us[0].Currency != "USD" {
		t.Fatalf("us = %+v", us)
	}
}

func TestClientKlineHKUSMinute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/hk-trends":
			if r.URL.Query().Get("secid") != "116.00700" {
				t.Fatalf("hk secid = %q, want 116.00700", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"trends": []string{"2024-12-16 09:31,391,392,393,390,110,210000,391.5"},
				},
			})
		case "/us-kline":
			if r.URL.Query().Get("secid") != "105.AAPL" {
				t.Fatalf("us secid = %q, want 105.AAPL", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"klines": []string{"2024-12-16 22:35,251,252,253,250,110,210000,1.19,0.40,1,0.02"},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyHKTrendsURL(server.URL+"/hk-trends"),
		WithEastmoneyUSKlineURL(server.URL+"/us-kline"),
		WithHTTPClient(server.Client()),
	)
	hk, err := client.Kline.HKMinute(context.Background(), "00700", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk.Timeline) != 1 || hk.Timeline[0].Currency != "HKD" || hk.Timeline[0].Close == nil || *hk.Timeline[0].Close != 392 {
		t.Fatalf("hk = %+v", hk)
	}
	us, err := client.Kline.USMinute(context.Background(), "105.AAPL", MinuteKlineOptions{
		Period:    MinutePeriodFive,
		StartDate: "2024-12-16 09:35",
		EndDate:   "2024-12-16",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(us.Klines) != 1 || us.Klines[0].Time != "2024-12-16 09:35" || us.Klines[0].Currency != "USD" {
		t.Fatalf("us = %+v", us)
	}
}

func TestClientKlineInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "cn history invalid period",
			call: func() error {
				_, err := client.GetHistoryKline(context.Background(), "600519", HistoryKlineOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "cn history invalid adjust",
			call: func() error {
				_, err := client.GetHistoryKline(context.Background(), "600519", HistoryKlineOptions{Adjust: "bad"})
				return err
			},
		},
		{
			name: "cn minute invalid period",
			call: func() error {
				_, err := client.GetMinuteKline(context.Background(), "600519", MinuteKlineOptions{Period: "120"})
				return err
			},
		},
		{
			name: "cn minute invalid adjust",
			call: func() error {
				_, err := client.GetMinuteKline(context.Background(), "600519", MinuteKlineOptions{Adjust: "bad"})
				return err
			},
		},
		{
			name: "hk history invalid period",
			call: func() error {
				_, err := client.GetHKHistoryKline(context.Background(), "00700", HKKlineOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "hk minute invalid adjust",
			call: func() error {
				_, err := client.GetHKMinuteKline(context.Background(), "00700", HKMinuteKlineOptions{Adjust: "bad"})
				return err
			},
		},
		{
			name: "us history invalid adjust",
			call: func() error {
				_, err := client.GetUSHistoryKline(context.Background(), "AAPL", USKlineOptions{Adjust: "bad"})
				return err
			},
		},
		{
			name: "us minute invalid period",
			call: func() error {
				_, err := client.GetUSMinuteKline(context.Background(), "AAPL", USMinuteKlineOptions{Period: "120"})
				return err
			},
		},
		{
			name: "industry board history invalid period",
			call: func() error {
				_, err := client.GetIndustryKline(context.Background(), "BK0001", IndustryBoardKlineOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "industry board minute invalid period",
			call: func() error {
				_, err := client.GetIndustryMinuteKline(context.Background(), "BK0001", IndustryBoardMinuteKlineOptions{Period: "120"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientFuturesKline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/futures" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("secid") != "113.rb2605" {
			t.Fatalf("secid = %q, want 113.rb2605", r.URL.Query().Get("secid"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"code":   "rb2605",
				"name":   "螺纹钢2605",
				"klines": []string{"2024-12-16,3500,3520,3530,3490,12345,67890000,1.14,0.57,20,0,0,98765,0"},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFuturesKlineURL(server.URL+"/futures"),
		WithHTTPClient(server.Client()),
	)
	rows, err := client.Futures.Kline(context.Background(), "rb2605", FuturesKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "rb2605" || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 98765 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientFuturesGlobalSpotAndKline(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/global":
			if r.URL.Query().Get("pageIndex") != "0" {
				t.Fatalf("pageIndex = %q", r.URL.Query().Get("pageIndex"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"total": float64(1),
				"list": []map[string]any{
					{"dm": "GC00Y", "name": "COMEX黄金", "p": 2400.5, "zde": 10.5, "zdf": 0.44, "o": 2390, "h": 2410, "l": 2380, "zjsj": 2390, "vol": 1000, "wp": 10, "np": 20, "ccl": 3000},
				},
			})
		case "/global-kline":
			if r.URL.Query().Get("secid") != "101.HG00Y" {
				t.Fatalf("secid = %q, want 101.HG00Y", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"code":   "HG00Y",
					"name":   "COMEX铜",
					"klines": []string{"2024-12-16,4.1,4.2,4.3,4.0,123,456,2.5,1.2,0.1,0,0,789,0"},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFuturesGlobalSpotURL(server.URL+"/global"),
		WithEastmoneyFuturesGlobalKlineURL(server.URL+"/global-kline"),
		WithHTTPClient(server.Client()),
	)
	spot, err := client.Futures.GlobalSpot(context.Background(), GlobalFuturesSpotOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(spot) != 1 || spot[0].Code != "GC00Y" || spot[0].Price == nil || *spot[0].Price != 2400.5 {
		t.Fatalf("spot = %+v", spot)
	}

	rows, err := client.Futures.GlobalKline(context.Background(), "HG00Y", GlobalFuturesKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Code != "HG00Y" || rows[0].OpenInterest == nil || *rows[0].OpenInterest != 789 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientFuturesInventory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		switch r.URL.Query().Get("reportName") {
		case "RPT_FUTU_POSITIONCODE":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"data": []map[string]any{{
						"TRADE_CODE":        "CU",
						"TRADE_TYPE":        "铜",
						"TRADE_MARKET_CODE": "SHFE",
					}},
				},
			})
		case "RPT_FUTU_STOCKDATA":
			if r.URL.Query().Get("filter") != `(SECURITY_CODE="CU")(TRADE_DATE>='2024-01-01')` {
				t.Fatalf("filter = %q", r.URL.Query().Get("filter"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"data": []map[string]any{{
						"SECURITY_CODE":  "CU",
						"TRADE_DATE":     "2024-12-16",
						"ON_WARRANT_NUM": 12345.0,
						"ADDCHANGE":      -100.0,
					}},
				},
			})
		case "RPT_FUTUOPT_GOLDSIL":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"data": []map[string]any{{
						"REPORT_DATE":   "2024-12-16",
						"STORAGE_TON":   100.5,
						"STORAGE_OUNCE": 3231.0,
					}},
				},
			})
		default:
			t.Fatalf("unexpected reportName %q", r.URL.Query().Get("reportName"))
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	symbols, err := client.Futures.InventorySymbols(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(symbols) != 1 || symbols[0].Code != "CU" || symbols[0].MarketCode != "SHFE" {
		t.Fatalf("symbols = %+v", symbols)
	}
	inventory, err := client.Futures.Inventory(context.Background(), "cu", FuturesInventoryOptions{StartDate: "2024-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	if len(inventory) != 1 || inventory[0].Inventory == nil || *inventory[0].Inventory != 12345 {
		t.Fatalf("inventory = %+v", inventory)
	}
	comex, err := client.Futures.ComexInventory(context.Background(), "gold", ComexInventoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(comex) != 1 || comex[0].Name != "黄金" || comex[0].StorageOunce == nil || *comex[0].StorageOunce != 3231 {
		t.Fatalf("comex = %+v", comex)
	}
}

func TestClientFuturesInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "domestic invalid symbol",
			call: func() error {
				_, err := client.GetFuturesKline(context.Background(), "2605")
				return err
			},
		},
		{
			name: "domestic unknown variety",
			call: func() error {
				_, err := client.GetFuturesKline(context.Background(), "unknown2605")
				return err
			},
		},
		{
			name: "domestic invalid period",
			call: func() error {
				_, err := client.GetFuturesKline(context.Background(), "rb2605", FuturesKlineOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "global invalid symbol",
			call: func() error {
				_, err := client.GetGlobalFuturesKline(context.Background(), "bad2507")
				return err
			},
		},
		{
			name: "global unknown variety",
			call: func() error {
				_, err := client.GetGlobalFuturesKline(context.Background(), "ZZZ2507")
				return err
			},
		},
		{
			name: "comex invalid symbol",
			call: func() error {
				_, err := client.GetComexInventory(context.Background(), "platinum")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error")
			}
			if got := GetErrorCode(err); got != CodeInvalidArgument {
				t.Fatalf("GetErrorCode = %s, want %s; err=%v", got, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/option":
			if r.URL.Query().Get("pageIndex") != "0" {
				t.Fatalf("pageIndex = %q", r.URL.Query().Get("pageIndex"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"list": []map[string]any{{
					"dm": "io2501C4000", "name": "沪深300购2501", "p": 123.4, "ccl": 3000.0,
				}},
			})
		case "/lhb":
			if r.URL.Query().Get("filter") != `(SECURITY_CODE="510050")(TRADE_DATE='2024-12-16')` {
				t.Fatalf("filter = %q", r.URL.Query().Get("filter"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"data": []map[string]any{{"SECURITY_CODE": "510050", "TRADE_DATE": "2024-12-16", "MEMBER_RANK": 1.0}},
				},
			})
		case "/etf-months":
			if r.URL.Query().Get("cate") != "50ETF" || r.URL.Query().Get("exchange") != "null" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":{"contractMonth":["全部","2024-06"],"stockId":"510050","cateId":"50ETF","cateList":["50ETF"]}}})`))
		case "/etf-expire":
			if r.URL.Query().Get("cate") != "50ETF" || r.URL.Query().Get("date") != "2024-06" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":{"expireDay":"2024-06-26","remainderDays":8,"stockId":"510050","other":{"name":"上证50ETF期权"}}}})`))
		case "/etf-minute":
			if r.URL.Query().Get("symbol") != "CON_OP_10009633" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":[{"i":"09:31","d":"2024-06-13","p":"1.23","v":"100","t":"200","a":"1.21"}]}})`))
		case "/jsonp_v2.php/ceheng_jsonp/StockOptionDaylineService.getSymbolInfo":
			if r.URL.Query().Get("symbol") != "CON_OP_10009633" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`cb([{"d":"2024-06-13","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"}])`))
		case "/etf-5day":
			if r.URL.Query().Get("symbol") != "CON_OP_10009633" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":[[{"i":"09:31","d":"2024-06-12","p":"1.10","v":"10","t":"20","a":"1.09"}],[{"i":"09:31","d":"2024-06-13","p":"1.20","v":"30","t":"40","a":"1.19"}]]}})`))
		case "/index-spot":
			query := r.URL.Query()
			if query.Get("type") != "futures" || query.Get("product") != "io" || query.Get("exchange") != "cffex" || query.Get("pinzhong") != "io2504" {
				t.Fatalf("query = %v", query)
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","3600","io2504C3600"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","io2504P3600"]]}}})`))
		case "/jsonp.php/ceheng_jsonp/FutureOptionAllService.getOptionDayline":
			symbol := r.URL.Query().Get("symbol")
			if symbol != "io2504C3600" && symbol != "au2506C580" {
				t.Fatalf("query = %v", r.URL.Query())
			}
			_, _ = w.Write([]byte(`cb([{"d":"2024-06-13","o":"1.10","h":"1.30","l":"1.00","c":"1.20","v":"1000"}])`))
		case "/commodity-spot":
			query := r.URL.Query()
			if query.Get("type") != "futures" || query.Get("product") != "au_o" || query.Get("exchange") != "shfe" || query.Get("pinzhong") != "au2506" {
				t.Fatalf("query = %v", query)
			}
			_, _ = w.Write([]byte(`callback({"result":{"data":{"up":[["10","1.10","1.20","1.30","20","300","0.05","580","au2506C580"]],"down":[["11","2.10","2.20","2.30","21","301","-0.06","au2506P580"]]}}})`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyOptionCFFEXURL(server.URL+"/option"),
		WithEastmoneyOptionLHBURL(server.URL+"/lhb"),
		WithSinaETFOptionListURL(server.URL+"/etf-months"),
		WithSinaETFOptionExpireURL(server.URL+"/etf-expire"),
		WithSinaETFOptionMinuteURL(server.URL+"/etf-minute"),
		WithSinaETFOptionDailyURL(server.URL+"/jsonp_v2.php/{callback}/StockOptionDaylineService.getSymbolInfo"),
		WithSinaETFOption5DayURL(server.URL+"/etf-5day"),
		WithSinaIndexOptionSpotURL(server.URL+"/index-spot"),
		WithSinaIndexOptionKlineURL(server.URL+"/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline"),
		WithSinaCommodityOptionSpotURL(server.URL+"/commodity-spot"),
		WithSinaCommodityOptionKlineURL(server.URL+"/jsonp.php/{callback}/FutureOptionAllService.getOptionDayline"),
		WithHTTPClient(server.Client()),
	)
	quotes, err := client.Options.CFFEXQuotes(context.Background(), CFFEXOptionQuotesOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Code != "io2501C4000" || quotes[0].OpenInterest == nil || *quotes[0].OpenInterest != 3000 {
		t.Fatalf("quotes = %+v", quotes)
	}
	lhb, err := client.Options.LHB(context.Background(), "510050", "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(lhb) != 1 || lhb[0].Symbol != "510050" || lhb[0].Rank != 1 {
		t.Fatalf("lhb = %+v", lhb)
	}
	months, err := client.Options.ETFOptionMonths(context.Background(), ETFOptionCate50ETF)
	if err != nil {
		t.Fatal(err)
	}
	if len(months.Months) != 1 || months.Months[0] != "2024-06" || months.CateID != "50ETF" {
		t.Fatalf("months = %+v", months)
	}
	expire, err := client.Options.ETFOptionExpireDay(context.Background(), ETFOptionCate50ETF, "2024-06")
	if err != nil {
		t.Fatal(err)
	}
	if expire.ExpireDay != "2024-06-26" || expire.RemainderDays != 8 || expire.Name != "上证50ETF期权" {
		t.Fatalf("expire = %+v", expire)
	}
	minutes, err := client.Options.ETFOptionMinute(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(minutes) != 1 || minutes[0].Price == nil || *minutes[0].Price != 1.23 {
		t.Fatalf("minutes = %+v", minutes)
	}
	klines, err := client.Options.ETFOptionDailyKline(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 1.20 {
		t.Fatalf("klines = %+v", klines)
	}
	fiveDay, err := client.Options.ETFOption5DayMinute(context.Background(), "10009633")
	if err != nil {
		t.Fatal(err)
	}
	if len(fiveDay) != 2 || fiveDay[1].Date != "2024-06-13" || fiveDay[1].OpenInterest == nil || *fiveDay[1].OpenInterest != 40 {
		t.Fatalf("fiveDay = %+v", fiveDay)
	}
	spot, err := client.Options.IndexOptionSpot(context.Background(), IndexOptionProductIO, "io2504")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot.Calls) != 1 || spot.Calls[0].Symbol != "io2504C3600" || spot.Calls[0].StrikePrice == nil || *spot.Calls[0].StrikePrice != 3600 {
		t.Fatalf("spot = %+v", spot)
	}
	indexKlines, err := client.Options.IndexOptionKline(context.Background(), "io2504C3600")
	if err != nil {
		t.Fatal(err)
	}
	if len(indexKlines) != 1 || indexKlines[0].Close == nil || *indexKlines[0].Close != 1.20 {
		t.Fatalf("indexKlines = %+v", indexKlines)
	}
	commoditySpot, err := client.Options.CommodityOptionSpot(context.Background(), "au", "au2506")
	if err != nil {
		t.Fatal(err)
	}
	if len(commoditySpot.Calls) != 1 || commoditySpot.Calls[0].Symbol != "au2506C580" || commoditySpot.Calls[0].StrikePrice == nil || *commoditySpot.Calls[0].StrikePrice != 580 {
		t.Fatalf("commoditySpot = %+v", commoditySpot)
	}
	commodityKlines, err := client.Options.CommodityOptionKline(context.Background(), "au2506C580")
	if err != nil {
		t.Fatal(err)
	}
	if len(commodityKlines) != 1 || commodityKlines[0].Close == nil || *commodityKlines[0].Close != 1.20 {
		t.Fatalf("commodityKlines = %+v", commodityKlines)
	}
}

func TestClientCommodityOptionUnknownVarietyReturnsInvalidArgument(t *testing.T) {
	client := New()

	_, err := client.GetCommodityOptionSpot(context.Background(), "unknown", "unknown2506")
	if err == nil {
		t.Fatal("expected error")
	}
	if got := GetErrorCode(err); got != CodeInvalidArgument {
		t.Fatalf("GetErrorCode = %s, want %s; err=%v", got, CodeInvalidArgument, err)
	}
}

func TestClientBoardIndustryConcept(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/industry":
			if r.URL.Query().Get("fs") != "m:90 t:2 f:!50" {
				t.Fatalf("industry fs = %q", r.URL.Query().Get("fs"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"diff": []map[string]any{{"f12": "BK0001", "f14": "酿酒行业", "f3": 2.1}},
				},
			})
		case "/concept":
			if r.URL.Query().Get("fs") != "m:90 t:3 f:!50" {
				t.Fatalf("concept fs = %q", r.URL.Query().Get("fs"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"diff": []map[string]any{{"f12": "BK1001", "f14": "人工智能", "f3": 3.2}},
				},
			})
		case "/spot":
			if r.URL.Query().Get("secid") != "90.BK0001" {
				t.Fatalf("secid = %q", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"f43": 1050.0},
			})
		case "/constituents":
			if r.URL.Query().Get("fs") != "b:BK0001 f:!50" {
				t.Fatalf("constituents fs = %q", r.URL.Query().Get("fs"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"diff": []map[string]any{{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0}},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyIndustryListURL(server.URL+"/industry"),
		WithEastmoneyConceptListURL(server.URL+"/concept"),
		WithEastmoneyIndustrySpotURL(server.URL+"/spot"),
		WithEastmoneyIndustryConstituentsURL(server.URL+"/constituents"),
		WithHTTPClient(server.Client()),
	)
	industry, err := client.Board.IndustryList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(industry) != 1 || industry[0].Name != "酿酒行业" {
		t.Fatalf("industry = %+v", industry)
	}
	concept, err := client.Board.ConceptList(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(concept) != 1 || concept[0].Name != "人工智能" {
		t.Fatalf("concept = %+v", concept)
	}
	spot, err := client.Board.IndustrySpot(context.Background(), "BK0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(spot) != 10 || spot[0].Value == nil || *spot[0].Value != 10.5 {
		t.Fatalf("spot = %+v", spot)
	}
	constituents, err := client.Board.IndustryConstituents(context.Background(), "BK0001")
	if err != nil {
		t.Fatal(err)
	}
	if len(constituents) != 1 || constituents[0].Code != "600519" {
		t.Fatalf("constituents = %+v", constituents)
	}
}

func TestClientBoardKlineMinute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/kline":
			if r.URL.Query().Get("secid") != "90.BK0001" {
				t.Fatalf("secid = %q", r.URL.Query().Get("secid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"klines": []string{"2024-12-16,100,105,106,99,12345,67890000,3.5,2,2,1.2"},
				},
			})
		case "/trends":
			if r.URL.Query().Get("ndays") != "1" {
				t.Fatalf("ndays = %q", r.URL.Query().Get("ndays"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"trends": []string{"2024-12-16 09:31,100,101,102,99,1000,2000,101.5"},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyIndustryKlineURL(server.URL+"/kline"),
		WithEastmoneyIndustryTrendsURL(server.URL+"/trends"),
		WithHTTPClient(server.Client()),
	)
	klines, err := client.Board.IndustryKline(context.Background(), "BK0001", HistoryKlineOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(klines) != 1 || klines[0].Close == nil || *klines[0].Close != 105 {
		t.Fatalf("klines = %+v", klines)
	}
	minute, err := client.Board.IndustryMinute(context.Background(), "BK0001", MinuteKlineOptions{Period: MinutePeriodOne})
	if err != nil {
		t.Fatal(err)
	}
	if len(minute.Timeline) != 1 || minute.Timeline[0].Price == nil || *minute.Timeline[0].Price != 101.5 {
		t.Fatalf("minute = %+v", minute)
	}
}

func TestClientFundFlow(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/fflow":
			if r.URL.Query().Get("secid") == "" {
				t.Fatalf("missing secid in %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"klines": []string{"2024-12-16,1000,100,200,300,400,10,1,2,3,4,1500,1.5,3500,1.2"},
				},
			})
		case "/clist":
			if r.URL.Query().Get("fid") != "f62" {
				t.Fatalf("fid = %q, want f62", r.URL.Query().Get("fid"))
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"diff": []map[string]any{{"f12": "600519", "f14": "贵州茅台", "f2": 1500.0, "f3": 1.5, "f62": 1000.0}},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFundFlowURL(server.URL+"/fflow"),
		WithEastmoneyClistURL(server.URL+"/clist"),
		WithHTTPClient(server.Client()),
	)
	individual, err := client.FundFlow.Individual(context.Background(), "600519", FundFlowOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(individual) != 1 || individual[0].Close == nil || *individual[0].Close != 1500 {
		t.Fatalf("individual = %+v", individual)
	}
	rank, err := client.FundFlow.Rank(context.Background(), FundFlowRankOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(rank) != 1 || rank[0].Name != "贵州茅台" {
		t.Fatalf("rank = %+v", rank)
	}
}

func TestClientFundFlowInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service individual invalid period",
			call: func() error {
				_, err := client.FundFlow.Individual(context.Background(), "600519", FundFlowOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "wrapper individual invalid period",
			call: func() error {
				_, err := client.GetIndividualFundFlow(context.Background(), "600519", FundFlowOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "service rank invalid indicator",
			call: func() error {
				_, err := client.FundFlow.Rank(context.Background(), FundFlowRankOptions{Indicator: "yearly"})
				return err
			},
		},
		{
			name: "wrapper rank invalid indicator",
			call: func() error {
				_, err := client.GetFundFlowRank(context.Background(), FundFlowRankOptions{Indicator: "yearly"})
				return err
			},
		},
		{
			name: "service sector rank invalid sector type",
			call: func() error {
				_, err := client.FundFlow.SectorRank(context.Background(), FundFlowRankOptions{SectorType: "theme"})
				return err
			},
		},
		{
			name: "wrapper sector rank invalid sector type",
			call: func() error {
				_, err := client.GetSectorFundFlowRank(context.Background(), FundFlowRankOptions{SectorType: "theme"})
				return err
			},
		},
		{
			name: "service sector history invalid period",
			call: func() error {
				_, err := client.FundFlow.SectorHistory(context.Background(), "BK0001", FundFlowOptions{Period: "yearly"})
				return err
			},
		},
		{
			name: "wrapper sector history invalid period",
			call: func() error {
				_, err := client.GetSectorFundFlowHistory(context.Background(), "BK0001", FundFlowOptions{Period: "yearly"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientNorthbound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/minute":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"s2nDate": "20241216",
					"s2n":     []string{"09:31,100,0,200,0,300"},
				},
			})
		case "/datacenter":
			if r.URL.Query().Get("reportName") == "" {
				t.Fatalf("missing reportName")
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"data": []map[string]any{{"TRADE_DATE": "2024-12-16", "NET_DEAL_AMT": 1000.0}},
				},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := New(
		WithEastmoneyNorthboundMinuteURL(server.URL+"/minute"),
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	minute, err := client.Northbound.Minute(context.Background(), NorthboundNorth)
	if err != nil {
		t.Fatal(err)
	}
	if len(minute) != 1 || minute[0].TotalNetInflow == nil || *minute[0].TotalNetInflow != 300 {
		t.Fatalf("minute = %+v", minute)
	}
	history, err := client.Northbound.History(context.Background(), NorthboundNorth, NorthboundHistoryOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].NetBuyAmount == nil || *history[0].NetBuyAmount != 1000 {
		t.Fatalf("history = %+v", history)
	}
}

func TestClientNorthboundInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service holding rank invalid period",
			call: func() error {
				_, err := client.Northbound.HoldingRank(context.Background(), NorthboundHoldingRankOptions{Period: "2year"})
				return err
			},
		},
		{
			name: "wrapper holding rank invalid period",
			call: func() error {
				_, err := client.GetNorthboundHoldingRank(context.Background(), NorthboundHoldingRankOptions{Period: "2year"})
				return err
			},
		},
		{
			name: "service holding rank invalid market",
			call: func() error {
				_, err := client.Northbound.HoldingRank(context.Background(), NorthboundHoldingRankOptions{Market: "hongkong"})
				return err
			},
		},
		{
			name: "wrapper holding rank invalid market",
			call: func() error {
				_, err := client.GetNorthboundHoldingRank(context.Background(), NorthboundHoldingRankOptions{Market: "hongkong"})
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientDragonTiger(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("reportName") == "" {
			t.Fatal("missing reportName")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "BILLBOARD_NET_AMT": 1000.0}},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	detail, err := client.DragonTiger.Detail(context.Background(), DragonTigerDateOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].Name != "贵州茅台" {
		t.Fatalf("detail = %+v", detail)
	}
}

func TestClientDragonTigerInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service stock stats invalid period",
			call: func() error {
				_, err := client.DragonTiger.StockStats(context.Background(), DragonTigerPeriod("2year"))
				return err
			},
		},
		{
			name: "wrapper stock stats invalid period",
			call: func() error {
				_, err := client.GetDragonTigerStockStats(context.Background(), DragonTigerPeriod("2year"))
				return err
			},
		},
		{
			name: "service branch rank invalid period",
			call: func() error {
				_, err := client.DragonTiger.BranchRank(context.Background(), DragonTigerPeriod("2year"))
				return err
			},
		},
		{
			name: "wrapper branch rank invalid period",
			call: func() error {
				_, err := client.GetDragonTigerBranchRank(context.Background(), DragonTigerPeriod("2year"))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientBlockTrade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("reportName") == "" {
			t.Fatal("missing reportName")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "DEAL_PRICE": 1490.0}},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	detail, err := client.BlockTrade.Detail(context.Background(), BlockTradeDateOptions{StartDate: "20241201", EndDate: "20241231"})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].Name != "贵州茅台" || detail[0].DealPrice == nil || *detail[0].DealPrice != 1490 {
		t.Fatalf("detail = %+v", detail)
	}
}

func TestClientMargin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("reportName") == "" {
			t.Fatal("missing reportName")
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "FIN_BUY_AMT": 200.0}},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	targets, err := client.Margin.TargetList(context.Background(), "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].Name != "贵州茅台" || targets[0].FinBuyAmount == nil || *targets[0].FinBuyAmount != 200 {
		t.Fatalf("targets = %+v", targets)
	}
}

func TestClientDividend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("reportName") != "RPT_SHAREBONUS_DET" {
			t.Fatalf("reportName = %q", r.URL.Query().Get("reportName"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{
				"data": []map[string]any{{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "REPORT_DATE": "2024-12-31", "PRETAX_BONUS_RMB": 30.0}},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyDatacenterURL(server.URL+"/datacenter"),
		WithHTTPClient(server.Client()),
	)
	rows, err := client.Dividend.Detail(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Name != "贵州茅台" || rows[0].DividendPretax == nil || *rows[0].DividendPretax != 30 {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientMarketEvent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/getAllStockChanges" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("type") != "8193" {
			t.Fatalf("type = %q", r.URL.Query().Get("type"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"allstock": []map[string]any{{"tm": 93055.0, "c": "600519", "n": "贵州茅台", "t": "8193", "i": "大单买入"}},
			},
		})
	}))
	defer server.Close()

	client := New(
		WithEastmoneyTopicURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	rows, err := client.MarketEvent.StockChanges(context.Background(), StockChangeLargeBuy)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].Name != "贵州茅台" || rows[0].ChangeTypeLabel != "大笔买入" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestClientTHSLimitUpPool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/limit_up_pool" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("date") != "20250613" || r.URL.Query().Get("order_field") != "330324" {
			t.Fatalf("query = %v", r.URL.Query())
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status_code": 0,
			"data": map[string]any{
				"page": map[string]any{"limit": 1, "total": 1, "count": 1, "page": 1},
				"info": []map[string]any{
					{"code": "002190", "name": "成飞集成", "last_limit_up_time": "1749797760"},
				},
				"date": "20250613",
			},
		})
	}))
	defer server.Close()

	client := New(
		WithTHSLimitUpPoolURL(server.URL+"/limit_up_pool"),
		WithHTTPClient(server.Client()),
	)
	result, err := client.GetTHSLimitUpPool(context.Background(), THSLimitUpPoolOptions{Date: "2025-06-13", Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if result.Date != "20250613" || len(result.Items) != 1 || result.Items[0].Code != "002190" {
		t.Fatalf("result = %+v", result)
	}
}

func TestClientDefaultTHSProviderPolicyAddsBrowserHeaders(t *testing.T) {
	client := New()
	policy, ok := client.core.ProviderPolicy(ProviderTHS)
	if !ok {
		t.Fatal("missing ths provider policy")
	}
	if policy.UserAgent == "" || policy.Headers["Referer"] != "https://data.10jqka.com.cn/market/ztStock/" || policy.Headers["Cookie"] == "" {
		t.Fatalf("policy = %+v", policy)
	}
}

func TestClientMarketEventInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	client := New()
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service zt pool invalid type",
			call: func() error {
				_, err := client.MarketEvent.ZTPool(context.Background(), ZTPoolType("unknown"), "")
				return err
			},
		},
		{
			name: "wrapper zt pool invalid typed type",
			call: func() error {
				_, err := client.GetZTPool(context.Background(), ZTPoolType("unknown"))
				return err
			},
		},
		{
			name: "service stock changes invalid type",
			call: func() error {
				_, err := client.MarketEvent.StockChanges(context.Background(), StockChangeType("unknown"))
				return err
			},
		},
		{
			name: "wrapper stock changes invalid type",
			call: func() error {
				_, err := client.GetStockChanges(context.Background(), StockChangeType("unknown"))
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func TestClientFundEstimate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/110011.js" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`jsonpgz({"fundcode":"110011","name":"易方达中小盘","dwjz":"3.5000","gsz":"3.5600","gszzl":"1.71"});`))
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFundGZURL(server.URL),
		WithHTTPClient(server.Client()),
	)
	row, err := client.Fund.Estimate(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if row.Code != "110011" || row.Name == nil || *row.Name != "易方达中小盘" || row.EstimatedChangePercent == nil || *row.EstimatedChangePercent != 1.71 {
		t.Fatalf("row = %+v", row)
	}
}

func TestClientFundNavAndRankHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pingzhongdata/110011.js" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`
var fS_code = "110011";
var fS_name = "易方达中小盘";
var Data_netWorthTrend = [{"x":1702857600000,"y":3.5,"equityReturn":"1.2","unitMoney":""}];
var Data_ACWorthTrend = [[1702857600000,5.6]];
var Data_rateInSimilarType = [{"x":1702857600000,"y":"12","sc":"300"}];
var Data_rateInSimilarPersent = [[1702857600000,4.0]];
`))
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFundPingzhongURL(server.URL+"/pingzhongdata"),
		WithHTTPClient(server.Client()),
	)
	nav, err := client.Fund.NavHistory(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(nav.Items) != 1 || nav.Items[0].Nav != 3.5 {
		t.Fatalf("nav = %+v", nav)
	}
	rank, err := client.Fund.RankHistory(context.Background(), "110011")
	if err != nil {
		t.Fatal(err)
	}
	if len(rank.Items) != 1 || rank.Items[0].Percentile == nil || *rank.Items[0].Percentile != 4 {
		t.Fatalf("rank = %+v", rank)
	}
}

func TestClientFundDividendList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Data/funddataIndex_Interface.aspx" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.URL.Query().Get("year") != "2024" {
			t.Fatalf("year = %q", r.URL.Query().Get("year"))
		}
		_, _ = w.Write([]byte(`var pageinfo = [1,20,1]; var jjfh_data = [["110011","易方达中小盘","2024-12-16","2024-12-17","0.12","2024-12-18","混合型"]];`))
	}))
	defer server.Close()

	client := New(
		WithEastmoneyFundDataIndexURL(server.URL+"/Data/funddataIndex_Interface.aspx"),
		WithHTTPClient(server.Client()),
	)
	result, err := client.Fund.DividendList(context.Background(), FundDividendListOptions{Year: "2024"})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].Name != "易方达中小盘" {
		t.Fatalf("result = %+v", result)
	}
}

func TestClientQuotesCN(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "sh600519" {
			t.Fatalf("q = %q, want sh600519", r.URL.Query().Get("q"))
		}
		fields := make([]string, 80)
		fields[0] = "1"
		fields[1] = "贵州茅台"
		fields[2] = "600519"
		fields[3] = "1700.00"
		fields[31] = "-10.00"
		fields[32] = "-0.58"
		writeGBK(t, w, `v_sh600519="`+joinFields(fields)+`";`)
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	quotes, err := client.Quotes.CN(context.Background(), []string{"sh600519"})
	if err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 {
		t.Fatalf("len(quotes) = %d, want 1", len(quotes))
	}
	if quotes[0].Name != "贵州茅台" || quotes[0].Code != "600519" || quotes[0].Price != 1700 {
		t.Fatalf("quote = %+v", quotes[0])
	}
}

func TestClientQuotesHKUSFund(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("q") {
		case "hk00700":
			fields := make([]string, 50)
			fields[0] = "100"
			fields[1] = "腾讯控股"
			fields[2] = "00700"
			fields[3] = "390.20"
			fields[47] = "HKD"
			writeGBK(t, w, `v_hk00700="`+joinFields(fields)+`";`)
		case "usAAPL":
			fields := make([]string, 50)
			fields[0] = "200"
			fields[1] = "APPLE"
			fields[2] = "AAPL"
			fields[3] = "205.50"
			writeGBK(t, w, `v_usAAPL="`+joinFields(fields)+`";`)
		case "jj110011":
			fields := make([]string, 9)
			fields[0] = "110011"
			fields[1] = "易方达中小盘"
			fields[5] = "3.5000"
			fields[8] = "2024-05-10"
			writeGBK(t, w, `v_jj110011="`+joinFields(fields)+`";`)
		default:
			t.Fatalf("unexpected q = %q", r.URL.Query().Get("q"))
		}
	}))
	defer server.Close()

	client := New(WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	hk, err := client.Quotes.HK(context.Background(), []string{"00700"})
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 1 || hk[0].Name != "腾讯控股" {
		t.Fatalf("hk = %+v", hk)
	}
	us, err := client.Quotes.US(context.Background(), []string{"AAPL"})
	if err != nil {
		t.Fatal(err)
	}
	if len(us) != 1 || us[0].Name != "APPLE" {
		t.Fatalf("us = %+v", us)
	}
	fund, err := client.Quotes.Fund(context.Background(), []string{"110011"})
	if err != nil {
		t.Fatal(err)
	}
	if len(fund) != 1 || fund[0].Name != "易方达中小盘" {
		t.Fatalf("fund = %+v", fund)
	}
}

func TestClientQuotesSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/s3/" || r.URL.Query().Get("q") != "茅台" {
			t.Fatalf("unexpected URL %s", r.URL.String())
		}
		_, _ = w.Write([]byte(`v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A";`))
	}))
	defer server.Close()

	client := New(WithSearchBaseURL(server.URL+"/s3/"), WithHTTPClient(server.Client()))
	results, err := client.Quotes.Search(context.Background(), "茅台")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "sh600519" || results[0].Name != "贵州茅台" {
		t.Fatalf("results = %+v", results)
	}
}

func TestClientDataNamespace(t *testing.T) {
	searchServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/s3/" || r.URL.Query().Get("q") != "茅台" {
			t.Fatalf("unexpected search URL %s", r.URL.String())
		}
		_, _ = w.Write([]byte(`v_hint="sh~600519~\u8d35\u5dde\u8305\u53f0~GZMT~GP-A";`))
	}))
	defer searchServer.Close()

	dataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/datacenter" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		var row map[string]any
		switch r.URL.Query().Get("reportName") {
		case "RPT_BLOCK_TRADE_DETAIL":
			row = map[string]any{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "DEAL_PRICE": 1490.0}
		case "RPT_MARGIN_TRADE_DETAIL":
			row = map[string]any{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "TRADE_DATE": "2024-12-16", "FIN_BUY_AMT": 200.0}
		case "RPT_SHAREBONUS_DET":
			row = map[string]any{"SECURITY_CODE": "600519", "SECURITY_NAME_ABBR": "贵州茅台", "REPORT_DATE": "2024-12-31", "PRETAX_BONUS_RMB": 30.0}
		default:
			t.Fatalf("reportName = %q", r.URL.Query().Get("reportName"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{"data": []map[string]any{row}},
		})
	}))
	defer dataServer.Close()

	client := New(
		WithSearchBaseURL(searchServer.URL+"/s3/"),
		WithEastmoneyDatacenterURL(dataServer.URL+"/datacenter"),
		WithHTTPClient(searchServer.Client()),
	)
	results, err := client.Data.Search(context.Background(), "茅台")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Code != "sh600519" {
		t.Fatalf("results = %+v", results)
	}
	detail, err := client.Data.BlockTradeDetail(context.Background(), BlockTradeDateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(detail) != 1 || detail[0].DealPrice == nil || *detail[0].DealPrice != 1490 {
		t.Fatalf("detail = %+v", detail)
	}
	targets, err := client.Data.MarginTargetList(context.Background(), "2024-12-16")
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 || targets[0].FinBuyAmount == nil || *targets[0].FinBuyAmount != 200 {
		t.Fatalf("targets = %+v", targets)
	}
	dividends, err := client.Data.DividendDetail(context.Background(), "sh600519")
	if err != nil {
		t.Fatal(err)
	}
	if len(dividends) != 1 || dividends[0].DividendPretax == nil || *dividends[0].DividendPretax != 30 {
		t.Fatalf("dividends = %+v", dividends)
	}
}

func TestClientQuotesTradingCalendar(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendar.txt" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("1990-12-19,1990-12-20"))
	}))
	defer server.Close()

	client := New(WithCalendarURL(server.URL+"/calendar.txt"), WithHTTPClient(server.Client()))
	calendar, err := client.Quotes.TradingCalendar(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(calendar) != 2 || calendar[1] != "1990-12-20" {
		t.Fatalf("calendar = %#v", calendar)
	}
}

func TestClientCalendarTradingDays(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendar.txt" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("2024-06-12,2024-06-13,2024-06-17"))
	}))
	defer server.Close()

	client := New(WithCalendarURL(server.URL+"/calendar.txt"), WithHTTPClient(server.Client()))
	ok, err := client.Calendar.IsTradingDay(context.Background(), "20240613")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected 20240613 to be trading day")
	}
	next, err := client.Calendar.NextTradingDay(context.Background(), "2024-06-13")
	if err != nil {
		t.Fatal(err)
	}
	if next != "2024-06-17" {
		t.Fatalf("next = %q", next)
	}
	if got := client.Calendar.MarketStatus(MarketA, time.Date(2024, 6, 13, 2, 0, 0, 0, time.UTC)); got != MarketStatusOpen {
		t.Fatalf("market status = %s", got)
	}
}

func TestClientCalendarOutOfRangeReturnsInvalidArgument(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/calendar.txt" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		_, _ = w.Write([]byte("2024-06-12,2024-06-13,2024-06-17"))
	}))
	defer server.Close()

	client := New(WithCalendarURL(server.URL+"/calendar.txt"), WithHTTPClient(server.Client()))
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "service next out of range",
			call: func() error {
				_, err := client.Calendar.NextTradingDay(context.Background(), "2024-06-17")
				return err
			},
		},
		{
			name: "wrapper next out of range",
			call: func() error {
				_, err := client.NextTradingDay(context.Background(), "2024-06-17")
				return err
			},
		},
		{
			name: "service prev out of range",
			call: func() error {
				_, err := client.Calendar.PrevTradingDay(context.Background(), "2024-06-12")
				return err
			},
		},
		{
			name: "wrapper prev out of range",
			call: func() error {
				_, err := client.PrevTradingDay(context.Background(), "2024-06-12")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected invalid argument error")
			}
			if code := GetErrorCode(err); code != CodeInvalidArgument {
				t.Fatalf("GetErrorCode(err) = %s, want %s; err=%v", code, CodeInvalidArgument, err)
			}
		})
	}
}

func joinFields(fields []string) string {
	out := ""
	for i, field := range fields {
		if i > 0 {
			out += "~"
		}
		out += field
	}
	return out
}

func writeGBK(t *testing.T, w http.ResponseWriter, text string) {
	t.Helper()
	encoded, err := simplifiedchinese.GBK.NewEncoder().String(text)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = w.Write([]byte(encoded))
}

type sequenceResponse struct {
	status int
	body   string
}

type sequenceRoundTripper struct {
	calls     int
	responses []sequenceResponse
}

func (t *sequenceRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.calls >= len(t.responses) {
		return nil, errors.New("unexpected request to " + req.URL.String())
	}
	response := t.responses[t.calls]
	t.calls++
	return &http.Response{
		StatusCode: response.status,
		Status:     http.StatusText(response.status),
		Body:       io.NopCloser(strings.NewReader(response.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}
