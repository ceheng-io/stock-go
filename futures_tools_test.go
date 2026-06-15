package stock

import "testing"

func TestFuturesToolsExtractVarietyAndMarketCode(t *testing.T) {
	tests := map[string]string{
		"rb2605": "rb",
		"RBM":    "RBM",
		"IF2604": "IF",
		"T2609":  "T",
	}
	for symbol, want := range tests {
		got, err := ExtractFuturesVariety(symbol)
		if err != nil {
			t.Fatalf("ExtractFuturesVariety(%q) returned error: %v", symbol, err)
		}
		if got != want {
			t.Fatalf("ExtractFuturesVariety(%q) = %q, want %q", symbol, got, want)
		}
	}

	marketTests := map[string]int{
		"rb":  113,
		"RBM": 113,
		"IF":  220,
		"IFM": 220,
		"TA":  115,
		"si":  225,
	}
	for variety, want := range marketTests {
		got, err := FuturesMarketCode(variety)
		if err != nil {
			t.Fatalf("FuturesMarketCode(%q) returned error: %v", variety, err)
		}
		if got != want {
			t.Fatalf("FuturesMarketCode(%q) = %d, want %d", variety, got, want)
		}
	}
}

func TestFuturesToolsReturnInvalidArgumentErrors(t *testing.T) {
	if _, err := ExtractFuturesVariety("2605"); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("ExtractFuturesVariety error = %v, code = %s", err, GetErrorCode(err))
	}
	if _, err := FuturesMarketCode("UNKNOWN"); GetErrorCode(err) != CodeInvalidArgument {
		t.Fatalf("FuturesMarketCode error = %v, code = %s", err, GetErrorCode(err))
	}
}
