package types

import "testing"

func TestNorthboundConstantsMatchTSUnions(t *testing.T) {
	directions := map[NorthboundDirection]string{
		NorthboundNorth: "north",
		NorthboundSouth: "south",
	}
	markets := map[NorthboundMarket]string{
		NorthboundMarketAll:      "all",
		NorthboundMarketShanghai: "shanghai",
		NorthboundMarketShenzhen: "shenzhen",
	}
	periods := map[NorthboundRankPeriod]string{
		NorthboundRankToday:    "today",
		NorthboundRankThreeDay: "3day",
		NorthboundRankFiveDay:  "5day",
		NorthboundRankTenDay:   "10day",
		NorthboundRankMonth:    "month",
		NorthboundRankQuarter:  "quarter",
		NorthboundRankYear:     "year",
	}

	assertStringEnumValues(t, directions)
	assertStringEnumValues(t, markets)
	assertStringEnumValues(t, periods)
}

func assertStringEnumValues[T ~string](t *testing.T, values map[T]string) {
	t.Helper()
	for got, want := range values {
		if string(got) != want {
			t.Fatalf("unexpected enum value: got %q want %q", got, want)
		}
	}
}
