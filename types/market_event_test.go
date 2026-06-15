package types

import "testing"

func TestMarketEventTypeConstantsMatchTSUnions(t *testing.T) {
	pools := []ZTPoolType{
		ZTPoolZT,
		ZTPoolYesterday,
		ZTPoolStrong,
		ZTPoolSubNew,
		ZTPoolBroken,
		ZTPoolDT,
	}
	wantPools := []string{"zt", "yesterday", "strong", "sub_new", "broken", "dt"}
	for index, value := range pools {
		if string(value) != wantPools[index] {
			t.Fatalf("pools[%d] = %q, want %q", index, value, wantPools[index])
		}
	}

	changes := []StockChangeType{
		StockChangeRocketLaunch,
		StockChangeQuickRebound,
		StockChangeLargeBuy,
		StockChangeLimitUpSeal,
		StockChangeLimitDownOpen,
		StockChangeBigBuyOrder,
		StockChangeAuctionUp,
		StockChangeHighOpen5D,
		StockChangeGapUp,
		StockChangeHigh60D,
		StockChangeSurge60D,
		StockChangeAccelerateDown,
		StockChangeHighDive,
		StockChangeLargeSell,
		StockChangeLimitDownSeal,
		StockChangeLimitUpOpen,
		StockChangeBigSellOrder,
		StockChangeAuctionDown,
		StockChangeLowOpen5D,
		StockChangeGapDown,
		StockChangeLow60D,
		StockChangeDrop60D,
	}
	wantChanges := []string{
		"rocket_launch",
		"quick_rebound",
		"large_buy",
		"limit_up_seal",
		"limit_down_open",
		"big_buy_order",
		"auction_up",
		"high_open_5d",
		"gap_up",
		"high_60d",
		"surge_60d",
		"accelerate_down",
		"high_dive",
		"large_sell",
		"limit_down_seal",
		"limit_up_open",
		"big_sell_order",
		"auction_down",
		"low_open_5d",
		"gap_down",
		"low_60d",
		"drop_60d",
	}
	for index, value := range changes {
		if string(value) != wantChanges[index] {
			t.Fatalf("changes[%d] = %q, want %q", index, value, wantChanges[index])
		}
	}

	if string(THSLimitUpOrderFirstLimitUpTime) != "330323" ||
		string(THSLimitUpOrderLastLimitUpTime) != "330324" ||
		string(THSLimitUpOrderOpenNum) != "330325" {
		t.Fatalf("unexpected ths limit-up order fields")
	}
	if string(THSLimitUpOrderDesc) != "0" || string(THSLimitUpOrderAsc) != "1" {
		t.Fatalf("unexpected ths limit-up order types")
	}
}
