package stock

import "testing"

func TestRootReExportsTimeUtilities(t *testing.T) {
	timestamp, ok := ParseMarketTime("2024-06-13 09:30", MarketTZ.US)
	if !ok {
		t.Fatal("ParseMarketTime ok = false")
	}
	if timestamp != 1718285400000 {
		t.Fatalf("timestamp = %d, want 1718285400000", timestamp)
	}

	meta := BuildTimeMeta("20240613", MarketTZ.CN)
	if meta.Timestamp == nil || *meta.Timestamp != 1718208000000 {
		t.Fatalf("meta timestamp = %v, want 1718208000000", meta.Timestamp)
	}
	if got := FormatInTz(meta.Timestamp, MarketTZ.CN); got != "2024-06-13 00:00" {
		t.Fatalf("FormatInTz = %q, want 2024-06-13 00:00", got)
	}
}

func TestRootReExportsTSTimeConstantName(t *testing.T) {
	if MARKET_TZ.CN != MarketTZ.CN || MARKET_TZ.HK != MarketTZ.HK || MARKET_TZ.US != MarketTZ.US {
		t.Fatalf("MARKET_TZ = %+v, want %+v", MARKET_TZ, MarketTZ)
	}
}
