package symbols

import (
	"errors"
	"testing"

	"github.com/ceheng.io/stock-go/internal/core"
)

func TestNormalizeAShare(t *testing.T) {
	tests := []struct {
		input    string
		exchange Exchange
		code     string
	}{
		{"sh600519", ExchangeSSE, "600519"},
		{"sz000001", ExchangeSZSE, "000001"},
		{"bj920819", ExchangeBSE, "920819"},
		{"600519", ExchangeSSE, "600519"},
		{"000001", ExchangeSZSE, "000001"},
		{"300750", ExchangeSZSE, "300750"},
		{"688981", ExchangeSSE, "688981"},
		{"510050", ExchangeSSE, "510050"},
		{"600519.SH", ExchangeSSE, "600519"},
		{"000001.SZ", ExchangeSZSE, "000001"},
		{"1.600519", ExchangeSSE, "600519"},
		{"0.000001", ExchangeSZSE, "000001"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Normalize(tt.input, nil)
			if err != nil {
				t.Fatal(err)
			}
			if got.Market != MarketCN || got.Exchange != tt.exchange || got.Code != tt.code || got.AssetType != AssetStock {
				t.Fatalf("Normalize(%q) = %+v", tt.input, got)
			}
		})
	}
}

func TestNormalizeHK(t *testing.T) {
	tests := map[string]string{
		"hk00700":   "00700",
		"hk700":     "00700",
		"00700":     "00700",
		"0700":      "00700",
		"00700.HK":  "00700",
		"116.00700": "00700",
	}

	for input, code := range tests {
		got, err := Normalize(input, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got.Market != MarketHK || got.Exchange != ExchangeHKEX || got.Code != code {
			t.Fatalf("Normalize(%q) = %+v", input, got)
		}
	}
}

func TestNormalizeUS(t *testing.T) {
	tests := []struct {
		input    string
		exchange Exchange
		code     string
	}{
		{"AAPL", ExchangeUS, "AAPL"},
		{"usAAPL", ExchangeUS, "AAPL"},
		{"105.AAPL", ExchangeNASDAQ, "AAPL"},
		{"106.BABA", ExchangeNYSE, "BABA"},
		{"SHW", ExchangeUS, "SHW"},
	}

	for _, tt := range tests {
		got, err := Normalize(tt.input, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got.Market != MarketUS || got.Exchange != tt.exchange || got.Code != tt.code {
			t.Fatalf("Normalize(%q) = %+v", tt.input, got)
		}
	}
}

func TestNormalizeFuturesAndBoard(t *testing.T) {
	cffex, err := Normalize("CFFEX.IF2412", nil)
	if err != nil {
		t.Fatal(err)
	}
	if cffex.Market != MarketCN || cffex.Exchange != ExchangeCFFEX || cffex.AssetType != AssetFutures || cffex.Code != "IF2412" || cffex.Variety != "IF" {
		t.Fatalf("CFFEX.IF2412 = %+v", cffex)
	}

	comex, err := Normalize("COMEX.GC", nil)
	if err != nil {
		t.Fatal(err)
	}
	if comex.Market != MarketGlobal || comex.Exchange != ExchangeCOMEX || comex.AssetType != AssetFutures {
		t.Fatalf("COMEX.GC = %+v", comex)
	}

	rb, err := Normalize("rb2510", &Hint{AssetType: AssetFutures})
	if err != nil {
		t.Fatal(err)
	}
	if rb.AssetType != AssetFutures || rb.Code != "RB2510" || rb.Variety != "RB" {
		t.Fatalf("rb2510 = %+v", rb)
	}

	board, err := Normalize("90.BK0475", nil)
	if err != nil {
		t.Fatal(err)
	}
	if board.AssetType != AssetBoard {
		t.Fatalf("90.BK0475 asset = %s, want board", board.AssetType)
	}
}

func TestNormalizeHintsAndSymbolRef(t *testing.T) {
	index, err := Normalize("600519", &Hint{AssetType: AssetIndex})
	if err != nil {
		t.Fatal(err)
	}
	if index.AssetType != AssetIndex {
		t.Fatalf("hint asset = %s, want index", index.AssetType)
	}

	ref, err := Normalize(SymbolRef{Code: "600519", AssetType: AssetIndex}, &Hint{AssetType: AssetFund})
	if err != nil {
		t.Fatal(err)
	}
	if ref.AssetType != AssetIndex {
		t.Fatalf("SymbolRef asset = %s, want index", ref.AssetType)
	}

	prefixed, err := Normalize(SymbolRef{Code: "sh600519"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if prefixed.Market != MarketCN || prefixed.Exchange != ExchangeSSE || prefixed.Code != "600519" {
		t.Fatalf("SymbolRef sh600519 = %+v", prefixed)
	}
}

func TestNormalizeInvalidInput(t *testing.T) {
	for _, input := range []string{"", "   ", "!!!", "@#$"} {
		t.Run(input, func(t *testing.T) {
			_, err := Normalize(input, nil)
			if err == nil {
				t.Fatalf("Normalize(%q) expected error", input)
			}
			var coded core.CodedError
			if !errors.As(err, &coded) {
				t.Fatalf("Normalize(%q) error = %T %v, want coded invalid symbol error", input, err, err)
			}
			if code := coded.SDKCode(); code != "INVALID_SYMBOL" {
				t.Fatalf("Normalize(%q) error code = %q, want INVALID_SYMBOL; err=%v", input, code, err)
			}
		})
	}
}

func TestNormalizeInvalidGlobalFuturesHintReturnsInvalidSymbol(t *testing.T) {
	_, err := Normalize("GC", &Hint{Market: MarketGlobal, AssetType: AssetFutures})
	if err == nil {
		t.Fatal("expected invalid symbol error")
	}
	var coded core.CodedError
	if !errors.As(err, &coded) {
		t.Fatalf("Normalize GLOBAL futures without exchange error = %T %v, want coded invalid symbol error", err, err)
	}
	if code := coded.SDKCode(); code != "INVALID_SYMBOL" {
		t.Fatalf("error code = %q, want INVALID_SYMBOL; err=%v", code, err)
	}
}

func TestProviderAdapters(t *testing.T) {
	cn, err := Normalize("600519", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := ToTencent(cn); got != "sh600519" {
		t.Fatalf("ToTencent = %q, want sh600519", got)
	}
	if got := ToEastmoneySecID(cn); got != "1.600519" {
		t.Fatalf("ToEastmoneySecID = %q, want 1.600519", got)
	}

	hk, err := Normalize("hk700", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := ToTencent(hk); got != "hk00700" {
		t.Fatalf("ToTencent HK = %q, want hk00700", got)
	}
	if got := ToEastmoneySecID(hk); got != "116.00700" {
		t.Fatalf("ToEastmoneySecID HK = %q, want 116.00700", got)
	}

	board, err := Normalize("90.BK0475", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := ToEastmoneySecID(board); got != "90.BK0475" {
		t.Fatalf("ToEastmoneySecID board = %q, want 90.BK0475", got)
	}
	if _, err := ToTencentE(board); err == nil {
		t.Fatal("ToTencentE(board) expected error")
	}
}

func TestProviderAdaptersInvalidArgumentsReturnInvalidArgument(t *testing.T) {
	tests := []struct {
		name string
		call func() error
	}{
		{
			name: "tencent unsupported asset type",
			call: func() error {
				_, err := ToTencentE(Normalized{Market: MarketCN, Exchange: ExchangeSSE, AssetType: AssetBoard, Code: "BK0475"})
				return err
			},
		},
		{
			name: "tencent unsupported cn exchange",
			call: func() error {
				_, err := ToTencentE(Normalized{Market: MarketCN, Exchange: ExchangeNASDAQ, AssetType: AssetStock, Code: "AAPL"})
				return err
			},
		},
		{
			name: "tencent unsupported market",
			call: func() error {
				_, err := ToTencentE(Normalized{Market: MarketGlobal, Exchange: ExchangeCOMEX, AssetType: AssetStock, Code: "GC"})
				return err
			},
		},
		{
			name: "eastmoney unsupported asset type",
			call: func() error {
				_, err := ToEastmoneySecIDE(Normalized{Market: MarketCN, Exchange: ExchangeSHFE, AssetType: AssetFutures, Code: "rb2510"})
				return err
			},
		},
		{
			name: "eastmoney unsupported exchange",
			call: func() error {
				_, err := ToEastmoneySecIDE(Normalized{Market: MarketGlobal, Exchange: ExchangeLME, AssetType: AssetStock, Code: "CAD"})
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
			var coded core.CodedError
			if !errors.As(err, &coded) {
				t.Fatalf("error = %T %v, want coded invalid argument error", err, err)
			}
			if code := coded.SDKCode(); code != "INVALID_ARGUMENT" {
				t.Fatalf("error code = %q, want INVALID_ARGUMENT; err=%v", code, err)
			}
		})
	}
}

func TestFuturesExchangesReturnsCopy(t *testing.T) {
	exchanges := FuturesExchanges()
	if exchanges["CFFEX"].Market != MarketCN || exchanges["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("CFFEX exchange = %+v", exchanges["CFFEX"])
	}
	if exchanges["COMEX"].Market != MarketGlobal || exchanges["COMEX"].Exchange != ExchangeCOMEX {
		t.Fatalf("COMEX exchange = %+v", exchanges["COMEX"])
	}

	exchanges["CFFEX"] = FuturesExchange{}
	fresh := FuturesExchanges()
	if fresh["CFFEX"].Market != MarketCN || fresh["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("FuturesExchanges returned mutable internal map: %+v", fresh["CFFEX"])
	}
}
