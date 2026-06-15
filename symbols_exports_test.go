package stock

import "testing"

func TestRootReExportsSymbolUtilities(t *testing.T) {
	var ref SymbolRef
	ref.Code = "600519"
	var input SymbolInput = ref
	var _ Market = SymbolMarketCN
	var _ AssetType = AssetStock
	var _ Exchange = ExchangeSSE
	var _ FuturesExchange = FuturesExchange{Market: SymbolMarketCN, Exchange: ExchangeCFFEX}

	normalized, err := NormalizeSymbol(input, nil)
	if err != nil {
		t.Fatal(err)
	}
	if normalized.Market != SymbolMarketCN || normalized.Exchange != ExchangeSSE || normalized.Code != "600519" {
		t.Fatalf("NormalizeSymbol = %+v", normalized)
	}

	if got := ToTencentSymbol(normalized); got != "sh600519" {
		t.Fatalf("ToTencentSymbol = %q, want sh600519", got)
	}
	if got := ToEastmoneySecID(normalized); got != "1.600519" {
		t.Fatalf("ToEastmoneySecID = %q, want 1.600519", got)
	}
	if got := ToEastmoneySecid(normalized); got != "1.600519" {
		t.Fatalf("ToEastmoneySecid = %q, want 1.600519", got)
	}
	if got := ToPlainCode(normalized); got != "600519" {
		t.Fatalf("ToPlainCode = %q, want 600519", got)
	}
	if got := InferAShareExchange("920819"); got != ExchangeBSE {
		t.Fatalf("InferAShareExchange = %q, want %q", got, ExchangeBSE)
	}
	if got := ExtractVariety("rb2510"); got != "RB" {
		t.Fatalf("ExtractVariety = %q, want RB", got)
	}

	exchanges := FuturesExchanges()
	if exchanges["CFFEX"].Market != SymbolMarketCN || exchanges["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("FuturesExchanges CFFEX = %+v", exchanges["CFFEX"])
	}
	exchanges["CFFEX"] = FuturesExchange{}
	if fresh := FuturesExchanges(); fresh["CFFEX"].Exchange != ExchangeCFFEX {
		t.Fatalf("FuturesExchanges returned mutable internal map: %+v", fresh["CFFEX"])
	}
}
