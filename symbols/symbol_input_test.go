package symbols

import "testing"

func TestSymbolInputCompatibilityName(t *testing.T) {
	var input SymbolInput = SymbolRef{Code: "600519"}
	normalized, err := Normalize(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	if normalized.Code != "600519" || normalized.Exchange != ExchangeSSE {
		t.Fatalf("Normalize(SymbolInput) = %+v", normalized)
	}
}

func TestTypeScriptStyleSymbolFunctionNames(t *testing.T) {
	normalized, err := NormalizeSymbol("hk700", nil)
	if err != nil {
		t.Fatal(err)
	}
	if got := ToTencentSymbol(normalized); got != "hk00700" {
		t.Fatalf("ToTencentSymbol = %q, want hk00700", got)
	}
	if got := ToEastmoneySecid(normalized); got != "116.00700" {
		t.Fatalf("ToEastmoneySecid = %q, want 116.00700", got)
	}
}

func TestTypeScriptStyleFuturesExchangesName(t *testing.T) {
	if FUTURES_EXCHANGES["CFFEX"].Market != MarketCN || FUTURES_EXCHANGES["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("FUTURES_EXCHANGES CFFEX = %+v", FUTURES_EXCHANGES["CFFEX"])
	}
	if FUTURES_EXCHANGES["COMEX"].Market != MarketGlobal || FUTURES_EXCHANGES["COMEX"].Exchange != ExchangeCOMEX {
		t.Fatalf("FUTURES_EXCHANGES COMEX = %+v", FUTURES_EXCHANGES["COMEX"])
	}

	FUTURES_EXCHANGES["CFFEX"] = FuturesExchange{}
	if fresh := FuturesExchanges(); fresh["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("FUTURES_EXCHANGES mutated internal map: %+v", fresh["CFFEX"])
	}
}
