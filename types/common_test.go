package types

import "testing"

func TestMarketConstantsMatchSymbolMarkets(t *testing.T) {
	markets := []Market{MarketCN, MarketHK, MarketUS, MarketGlobal}
	want := []string{"CN", "HK", "US", "GLOBAL"}
	for index, market := range markets {
		if string(market) != want[index] {
			t.Fatalf("markets[%d] = %q, want %q", index, market, want[index])
		}
	}
}

func TestAnyHistoryKlineAcceptsConcreteHistoryRows(t *testing.T) {
	rows := []AnyHistoryKline{
		HistoryKline{Code: "600519"},
		HKHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "00700"}},
		USHistoryKline{ForeignHistoryKline: ForeignHistoryKline{Code: "AAPL"}},
	}

	if len(rows) != 3 {
		t.Fatalf("len(rows) = %d, want 3", len(rows))
	}
}
