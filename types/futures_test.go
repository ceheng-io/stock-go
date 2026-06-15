package types

import "testing"

func TestFuturesExchangeConstantsMatchTSUnion(t *testing.T) {
	values := map[FuturesExchange]string{
		FuturesExchangeSHFE:  "SHFE",
		FuturesExchangeDCE:   "DCE",
		FuturesExchangeCZCE:  "CZCE",
		FuturesExchangeINE:   "INE",
		FuturesExchangeCFFEX: "CFFEX",
		FuturesExchangeGFEX:  "GFEX",
	}

	if len(values) != 6 {
		t.Fatalf("expected six domestic futures exchanges, got %d", len(values))
	}
	for got, want := range values {
		if string(got) != want {
			t.Fatalf("unexpected futures exchange value: got %q want %q", got, want)
		}
	}
}
